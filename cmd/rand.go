package cmd

import (
	"context"
	"log/slog"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andrewmostello/metar-ws2811/metar"
	"github.com/andrewmostello/metar-ws2811/ws2811"
	"github.com/oklog/oklog/pkg/group"
	"github.com/spf13/cobra"
)

var randCmd = &cobra.Command{
	Use:   "rand",
	Short: "Random LED colors",
	Long:  `Run random colors across the LEDs`,
	Run: func(cmd *cobra.Command, args []string) {
		execOp(testLEDs)
	},
}

func init() {
	rootCmd.AddCommand(randCmd)
}

func randLEDs(logger *slog.Logger) error {
	ctrl := &ws2811.Controller{
		Logger: logger,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dur := 5 * time.Second

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

	src := make(chan (map[int]metar.FlightCategory))

	g.Add(
		func() error {
			tick := time.NewTicker(dur)

			rand.New(rand.NewSource(time.Now().UnixNano()))

			vec := make(map[int]metar.FlightCategory, ws2811.DefaultLEDCount)
			for i := 0; i < ws2811.DefaultLEDCount; i++ {
				vec[i] = metar.FlightCategoryVFR
			}

			for {
				select {
				case <-tick.C:
					rgb := ws2811.RGB{
						Red:   rand.Intn(256),
						Green: rand.Intn(256),
						Blue:  rand.Intn(256),
					}
					logger.Info("rendering", "color", rgb)
					ctrl.Colors = map[metar.FlightCategory]ws2811.RGB{
						metar.FlightCategoryUnknown: ws2811.RGB{
							Red:   0,
							Green: 0,
							Blue:  0,
						},
						metar.FlightCategoryVFR: rgb,
					}
					src <- vec
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

	logger.Info("starting rand", "duration", dur)

	defer func() {
		logger.Info("stopping rand")
	}()

	return g.Run()
}
