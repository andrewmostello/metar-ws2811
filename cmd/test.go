package cmd

import (
	"context"
	"log/slog"
	"time"

	"github.com/andrewmostello/metar-ws2811/metar"
	"github.com/andrewmostello/metar-ws2811/ws2811"
	"github.com/spf13/cobra"
)

// testCmd represents the version command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test LED strip",
	Long:  `Cycle through METAR flight categories on the LED strip.`,
	Run: func(cmd *cobra.Command, args []string) {
		execOp(testLEDs)
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}

func testLEDs(logger *slog.Logger) error {
	ctrl := &ws2811.Controller{
		Logger: logger,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dur := time.Second

	src := make(chan (map[int]metar.FlightCategory))
	go testSource(ctx, logger, dur, src)

	logger.Info("starting test", "duration", dur)

	defer func() {
		logger.Info("stopping test")
	}()

	return ctrl.Serve(ctx, src)
}

func next(last map[int]metar.FlightCategory) map[int]metar.FlightCategory {
	next := make(map[int]metar.FlightCategory)
	for i, cat := range last {
		switch cat {
		case metar.FlightCategoryUnknown:
			next[i] = metar.FlightCategoryVFR
		case metar.FlightCategoryVFR:
			next[i] = metar.FlightCategoryMVFR
		case metar.FlightCategoryMVFR:
			next[i] = metar.FlightCategoryIFR
		case metar.FlightCategoryIFR:
			next[i] = metar.FlightCategoryLIFR
		case metar.FlightCategoryLIFR:
			next[i] = metar.FlightCategoryUnknown
		}
	}
	return next
}

func testSource(ctx context.Context, logger *slog.Logger, delay time.Duration, src chan (map[int]metar.FlightCategory)) {

	if delay <= 0 {
		delay = time.Second
	}

	nxt := map[int]metar.FlightCategory{
		0:  metar.FlightCategoryUnknown,
		2:  metar.FlightCategoryVFR,
		4:  metar.FlightCategoryMVFR,
		6:  metar.FlightCategoryIFR,
		8:  metar.FlightCategoryLIFR,
		10: metar.FlightCategoryUnknown,
		12: metar.FlightCategoryVFR,
		14: metar.FlightCategoryMVFR,
		16: metar.FlightCategoryIFR,
		18: metar.FlightCategoryLIFR,
		20: metar.FlightCategoryUnknown,
		22: metar.FlightCategoryVFR,
		24: metar.FlightCategoryMVFR,
		26: metar.FlightCategoryIFR,
		28: metar.FlightCategoryLIFR,
	}

	tick := time.NewTicker(delay)

	for {
		select {
		case <-tick.C:
			logger.Info("rendering next", "next", nxt)
			src <- nxt
			nxt = next(nxt)
		case <-ctx.Done():
			return
		}
	}
}
