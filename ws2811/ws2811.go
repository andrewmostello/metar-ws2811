package light

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/andrewmostello/metar-ws2811/metar"
	ws281x "github.com/rpi-ws281x/rpi-ws281x-go"
)

const (
	DefaultBrightness = 128
	DefaultLEDCounts  = 300
	DefaultGPIOPin    = 18
)

var (
	DefaultColors = map[metar.FlightCategory]RGB{
		metar.FlightCategoryUnknown: RGB{
			Red:   0,
			Green: 0,
			Blue:  0,
		},
		metar.FlightCategoryVFR: RGB{
			Red:   0,
			Green: 255,
			Blue:  0,
		},
		metar.FlightCategoryMVFR: RGB{
			Red:   0,
			Green: 0,
			Blue:  255,
		},
		metar.FlightCategoryIFR: RGB{
			Red:   255,
			Green: 0,
			Blue:  0,
		},
		metar.FlightCategoryLIFR: RGB{
			Red:   125,
			Green: 0,
			Blue:  125,
		},
	}
)

type RGB struct {
	Red   int
	Green int
	Blue  int
}

func (rgb RGB) ToColor() uint32 {
	return uint32(uint32(rgb.Red)<<16 | uint32(rgb.Green)<<8 | uint32(rgb.Blue))
}

type Option func(*ws281x.ChannelOption)

type Controller struct {
	Logger *slog.Logger
	Colors map[metar.FlightCategory]RGB
}

func RGBToColor(r int, g int, b int) uint32 {
	return uint32(uint32(r)<<16 | uint32(g)<<8 | uint32(b))
}

func (ctrl *Controller) Render(drv *ws281x.WS2811, cats map[int]metar.FlightCategory) error {
	leds := drv.Leds(0)
	for i := 0; i < len(leds); i++ {
		cat := cats[i]
		rgb := ctrl.Colors[cat]
		leds[i] = rgb.ToColor()
	}

	if err := drv.Render(); err != nil {
		return err
	}

	if err := drv.Wait(); err != nil {
		if l := ctrl.Logger; l != nil {
			l.Warn("wait failure", "error", err)
		}
	}

	return nil
}

func (ctrl *Controller) DefaultOptions() []Option {
	return []Option{
		func(opt *ws281x.ChannelOption) {
			opt.Brightness = DefaultBrightness
			opt.LedCount = DefaultLEDCounts
			opt.GpioPin = DefaultGPIOPin
		},
	}
}

func (ctrl *Controller) applyOptions(drvopt *ws281x.Option, opts ...Option) {
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		chopt := drvopt.Channels[0]
		opt(&chopt)
		drvopt.Channels[0] = chopt
	}
}

func (ctrl *Controller) Serve(ctx context.Context, rnd chan (map[int]metar.FlightCategory), opts ...Option) error {

	drvopts := ws281x.DefaultOptions
	ctrl.applyOptions(&drvopts, ctrl.DefaultOptions()...)
	ctrl.applyOptions(&drvopts, opts...)

	drv, err := ws281x.MakeWS2811(&drvopts)
	if err != nil {
		return fmt.Errorf("failed to create driver: %w", err)
	}

	if err := drv.Init(); err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	defer drv.Fini()

	for {
		select {
		case cats := <-rnd:
			if err := ctrl.Render(drv, cats); err != nil {
				if l := ctrl.Logger; l != nil {
					l.Error("render failure", "error", err)
				}
			}
		case <-ctx.Done():
			if l := ctrl.Logger; l != nil {
				l.Debug("context done", "error", ctx.Err())
			}
			return nil
		}
	}
}
