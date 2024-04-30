package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

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

func serve(logger *slog.Logger, ctrl *ws2811.Controller, ledcfg config.LED) error {

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
		signal.Notify(term, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
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

	srv := &metar.ColorServer{
		Logger:              logger,
		AirportIDs:          cfg.AirportIDs,
		LEDIndexByAirportID: cfg.LEDIndexes,
		Timeout:             mcfg.Timeout,
		Client: metar.Client{
			BaseURL: mcfg.BaseURL,
		},
	}

	leds := make(chan (map[int]ws2811.RGB))

	g.Add(
		func() error {
			return srv.Serve(ctx, cfg.RefreshCron, leds)
		},
		func(err error) {
			cancel()
			close(leds)
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

	logger.Info("serving LEDs")

	defer func() {
		logger.Info("shutting down")
	}()

	return g.Run()
}
