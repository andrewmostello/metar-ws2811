package cmd

import (
	"context"
	"log/slog"
	"time"

	"github.com/andrewmostello/metar-ws2811/metar"
	"github.com/andrewmostello/metar-ws2811/ws2811"
	"github.com/spf13/cobra"
)

var offCmd = &cobra.Command{
	Use:   "off",
	Short: "Turn off LED lights",
	Long:  `Turn off all LEDs`,
	Run: func(cmd *cobra.Command, args []string) {
		execOp(offLEDs)
	},
}

func init() {
	rootCmd.AddCommand(offCmd)
}

func offLEDs(logger *slog.Logger) error {
	ctrl := &ws2811.Controller{
		Logger: logger,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dur := 5 * time.Second

	src := make(chan (map[int]metar.FlightCategory))
	go func() {
		src <- map[int]metar.FlightCategory{}
		cancel()
	}()

	logger.Info("starting rand", "duration", dur)

	defer func() {
		logger.Info("stopping rand")
	}()

	return ctrl.Serve(ctx, src)
}
