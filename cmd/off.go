package cmd

import (
	"context"
	"log/slog"

	"github.com/andrewmostello/metar-ws2811/config"
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

func offLEDs(logger *slog.Logger, ctrl *ws2811.Controller, cfg config.LED) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Info("shutting off LEDs")

	if err := ctrl.SetAllLEDs(ctx, ws2811.Off); err != nil {
		logger.With("error", err).Error("failed to set LEDs")
		return err
	}

	return nil
}
