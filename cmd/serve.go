package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andrewmostello/metar-ws2811/config"
	"github.com/andrewmostello/metar-ws2811/metar"
	"github.com/andrewmostello/metar-ws2811/ws2811"
	"github.com/oklog/oklog/pkg/group"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Update LED strip with METAR data",
	Long:  `Retrieve METAR data and update the LED strip.`,
	Run: func(cmd *cobra.Command, args []string) {
		execOp(serve)
	},
}

func init() {
	config.AddServeFlags(serveCmd)
	config.AddMetarFlags(serveCmd)

	rootCmd.AddCommand(serveCmd)
}

func serve(logger *slog.Logger, ctrl *ws2811.Controller) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.GetServe()
	if err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	mcfg := config.GetMETAR()

	var g group.Group
	{
		term := make(chan os.Signal, 1)
		signal.Notify(term, os.Interrupt, syscall.SIGTERM)
		cancel := make(chan struct{})
		g.Add(
			func() error {
				select {
				case <-term:
					break
				case <-cancel:
					break
				}
				return nil
			},
			func(err error) {
				close(cancel)
			},
		)
	}

	leds := make(chan (map[int]metar.FlightCategory))

	g.Add(
		func() error {

			{
				fcs, err := refreshMETARs(ctx, logger, cfg, mcfg)
				if err != nil {
					return err
				}

				leds <- fcs
			}

			for {
				nxt := cfg.RefreshCron.Next(time.Now())

				logger.Info("scheduled next refresh", "next", nxt)

				t := time.NewTimer(time.Until(nxt))

				select {
				case <-t.C:
					logger.Debug("refreshing")

					fcs, err := refreshMETARs(ctx, logger, cfg, mcfg)
					if err != nil {
						logger.Error("failed refresh", "error", err)
					}

					leds <- fcs

				case <-ctx.Done():
					if !t.Stop() {
						<-t.C
					}
					close(leds)
					return nil
				}
			}
		},
		func(err error) {
			cancel()
		},
	)

	g.Add(
		func() error {
			return ctrl.Serve(ctx, leds)
		},
		func(err error) {
			cancel()
		},
	)

	logger.Info("updating LEDs")

	defer func() {
		logger.Info("shutting down")
	}()

	return g.Run()
}

func refreshMETARs(ctx context.Context, logger *slog.Logger, cfg config.Serve, mcfg config.METAR) (map[int]metar.FlightCategory, error) {

	to := 15 * time.Second
	if v := mcfg.Timeout; v > 0 {
		to = v
	}

	ctx, cancel := context.WithTimeout(ctx, to)
	defer cancel()

	clnt := metar.Client{
		BaseURL: mcfg.BaseURL,
	}

	metars, err := clnt.GetMETARs(ctx, cfg.AirportIDs...)
	if err != nil {
		return nil, fmt.Errorf("failed to get METARs: %w", err)
	}

	idxs := cfg.LEDIndexes

	fcs := make(map[int]metar.FlightCategory, len(idxs))

	for id, wx := range metars {
		idx, ok := idxs[id]
		if !ok {
			logger.Warn("no LED index for airport", "airport", id)
			continue
		}
		fc := wx.FlightCategory()
		logger.Info("METAR", "airport", id, "flightCategory", fc.Name(), "weather", wx.RawObservation)
		fcs[idx] = fc
	}

	return fcs, nil
}
