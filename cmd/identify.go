package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/andrewmostello/metar-ws2811/config"
	"github.com/andrewmostello/metar-ws2811/ws2811"
	"github.com/oklog/oklog/pkg/group"
	"github.com/spf13/cobra"
)

var identifyCmd = &cobra.Command{
	Use:   "identify",
	Short: "Flash an LED to identify it",
	Long:  `Flash an LED to identify it`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Fprintf(os.Stderr, "requires index argument")
			os.Exit(1)
			return
		}
		index, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid index: %v", err)
			os.Exit(1)
			return
		}
		execOp(flashLED(index))
	},
}

func init() {
	rootCmd.AddCommand(identifyCmd)
}

func flashLED(index int) func(logger *slog.Logger, ctrl *ws2811.Controller, cfg config.LED) error {

	return func(logger *slog.Logger, ctrl *ws2811.Controller, cfg config.LED) error {

		if index < 0 {
			return fmt.Errorf("index must be positive : %d", index)
		}

		if index >= cfg.Count {
			return fmt.Errorf("index must be below max LED count: %d >= %d", index, cfg.Count)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		dur := 1 * time.Second

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

		src := make(chan (map[int]ws2811.RGB))

		g.Add(
			func() error {
				tick := time.NewTicker(dur)

				vec := make(map[int]ws2811.RGB, cfg.Count)
				on := true

				nxt := func() {
					if on {
						vec[index] = ws2811.RGB{Red: 0, Green: 255, Blue: 0}
					} else {
						vec[index] = ws2811.RGB{Red: 0, Green: 0, Blue: 0}
					}
					on = !on
					src <- vec
				}

				nxt()

				for {
					select {
					case <-tick.C:
						nxt()
					case <-ctx.Done():
						tick.Stop()
						close(src)
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
				return ctrl.Serve(ctx, src)
			},
			func(err error) {
				cancel()
			},
		)

		logger.Info("flashing LED", "index", index)

		defer func() {
			logger.Info("stopping")
		}()

		return g.Run()
	}
}
