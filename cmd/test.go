package cmd

import (
	"context"
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

func testLEDs(logger *slog.Logger, ctrl *ws2811.Controller, cfg config.LED) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dur := time.Second

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
			if dur <= 0 {
				dur = time.Second
			}

			vec := make(map[int]metar.FlightCategory)
			nxt := metar.FlightCategoryUnknown
			for i := 0; i < cfg.Count; i++ {
				vec[i] = nxt
				nxt = next(nxt)
			}

			tick := time.NewTicker(dur)

			for {
				select {
				case <-tick.C:
					logger.Info("rendering next", "vec", vec)
					src <- metar.FlightCategoryToRGB(nil, vec)
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

func nextVec(last map[int]metar.FlightCategory) map[int]metar.FlightCategory {
	nxt := make(map[int]metar.FlightCategory)
	for i, cat := range last {
		nxt[i] = next(cat)
	}
	return nxt
}

func next(last metar.FlightCategory) metar.FlightCategory {
	switch last {
	case metar.FlightCategoryUnknown:
		return metar.FlightCategoryVFR
	case metar.FlightCategoryVFR:
		return metar.FlightCategoryMVFR
	case metar.FlightCategoryMVFR:
		return metar.FlightCategoryIFR
	case metar.FlightCategoryIFR:
		return metar.FlightCategoryLIFR
	case metar.FlightCategoryLIFR:
		return metar.FlightCategoryUnknown
	}
	return metar.FlightCategoryUnknown
}
