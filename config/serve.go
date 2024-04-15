package config

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	cfgKeyServeRefreshCron = "serve.refresh_cron"
	cfgKeyServeAirportIDs  = "serve.airport_ids"
	cfgKeyServeLEDIndexes  = "serve.led_indexes"
	cfgKeyMETARBaseURL     = "metar.base_url"
	cfgKeyMETARTimeout     = "metar.timeout_seconds"
)

func durationInSeconds(dur int64) time.Duration {
	if dur == 0 {
		return 0
	}
	return time.Duration(dur) * time.Second
}

func expandCommaSeparatedList(s []string) []string {
	expanded := make([]string, 0, len(s))
	for _, v := range s {
		expanded = append(expanded, strings.Split(v, ",")...)
	}
	return expanded
}

type Serve struct {
	RefreshCron cron.Schedule
	AirportIDs  []string
	LEDIndexes  map[string]int
}

func GetServe() (Serve, error) {

	refreshCron := viper.GetString(cfgKeyServeRefreshCron)

	refreshSchedule, err := cron.ParseStandard(refreshCron)
	if err != nil {
		return Serve{}, fmt.Errorf("unable to parse refresh cron schedule: %w", err)
	}

	ledIndexes := expandCommaSeparatedList(viper.GetStringSlice(cfgKeyServeLEDIndexes))
	ledIndexMap := make(map[string]int, len(ledIndexes))
	for _, ledIndex := range ledIndexes {
		parts := strings.Split(ledIndex, "=")
		if len(parts) != 2 {
			return Serve{}, fmt.Errorf("invalid LED index format: %s", ledIndex)
		}
		idx, err := strconv.Atoi(parts[1])
		if err != nil {
			return Serve{}, fmt.Errorf("invalid LED index value: %s", parts[1])
		}
		ledIndexMap[parts[0]] = idx
	}

	ids := expandCommaSeparatedList(viper.GetStringSlice(cfgKeyServeAirportIDs))
	last := -1
	for _, id := range ids {
		if idx, ok := ledIndexMap[id]; ok {
			last = idx
			continue
		}
		last++
		ledIndexMap[id] = last
	}

	return Serve{
		RefreshCron: refreshSchedule,
		AirportIDs:  ids,
		LEDIndexes:  ledIndexMap,
	}, nil
}

func AddServeFlags(cmd *cobra.Command) {
	flag := "serve-refresh-cron"
	cmd.PersistentFlags().String(flag, "*/15 * * * *", "Cron format schedule on which to refresh METAR data.")
	viper.BindPFlag(cfgKeyServeRefreshCron, cmd.PersistentFlags().Lookup(flag))

	flag = "serve-airport-ids"
	cmd.PersistentFlags().StringSlice(flag, []string{"KBOS,KJFK,KSFO,KORD"}, "Airport IDs to retrieve METAR data for. Accepts multiple arguments and will explode any comma separated lists.")
	viper.BindPFlag(cfgKeyServeAirportIDs, cmd.PersistentFlags().Lookup(flag))

	flag = "serve-led-indexes"
	cmd.PersistentFlags().StringSlice(flag, []string{}, "Index of LED for a specified airport ID. Arguments should be in the format of 'airport_id=led_index', e.g. \"KBOS=15\". Accepts multiple arguments and will explode any comma separated lists.")
	viper.BindPFlag(cfgKeyServeLEDIndexes, cmd.PersistentFlags().Lookup(flag))
}

type METAR struct {
	BaseURL string
	Timeout time.Duration
}

func GetMETAR() METAR {
	return METAR{
		BaseURL: viper.GetString(cfgKeyMETARBaseURL),
		Timeout: durationInSeconds(viper.GetInt64(cfgKeyMETARTimeout)),
	}
}

func AddMetarFlags(cmd *cobra.Command) {
	flag := "metar-timeout-seconds"
	cmd.PersistentFlags().Int(flag, 15, "Seconds to wait for a METAR request to complete.")
	viper.BindPFlag(cfgKeyMETARTimeout, cmd.PersistentFlags().Lookup(flag))

	flag = "metar-base-url"
	cmd.PersistentFlags().String(flag, "https://aviationweather.gov/api/data", "Base URL for Aviation Weather Center data API for METAR requests.")
	viper.BindPFlag(cfgKeyMETARBaseURL, cmd.PersistentFlags().Lookup(flag))
}
