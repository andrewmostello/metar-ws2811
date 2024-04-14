package cmd

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andrewmostello/metar-ws2811/metar"
	"github.com/andrewmostello/metar-ws2811/ws2811"
	"github.com/oklog/oklog/pkg/group"
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
			if dur <= 0 {
				dur = time.Second
			}

			nxt := map[int]metar.FlightCategory{
				0:  metar.FlightCategoryUnknown,
				1:  metar.FlightCategoryVFR,
				2:  metar.FlightCategoryMVFR,
				3:  metar.FlightCategoryIFR,
				4:  metar.FlightCategoryLIFR,
				5:  metar.FlightCategoryUnknown,
				6:  metar.FlightCategoryVFR,
				7:  metar.FlightCategoryMVFR,
				8:  metar.FlightCategoryIFR,
				9:  metar.FlightCategoryLIFR,
				10: metar.FlightCategoryUnknown,
				11: metar.FlightCategoryVFR,
				12: metar.FlightCategoryMVFR,
				13: metar.FlightCategoryIFR,
				14: metar.FlightCategoryLIFR,
				15: metar.FlightCategoryUnknown,
			}

			tick := time.NewTicker(dur)

			for {
				select {
				case <-tick.C:
					logger.Info("rendering next", "next", nxt)
					src <- nxt
					nxt = next(nxt)
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

	logger.Info("starting test", "duration", dur)

	defer func() {
		logger.Info("stopping test")
	}()

	return g.Run()
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
