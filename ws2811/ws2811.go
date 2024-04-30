package ws2811

import (
	"context"
	"fmt"
	"log/slog"

	ws281x "github.com/rpi-ws281x/rpi-ws281x-go"
)

const (
	DefaultBrightness = 128
	DefaultLEDCount   = 50
	DefaultGPIOPin    = 18
)

var (
	Off = RGB{
		Red:   0,
		Green: 0,
		Blue:  0,
	}
)

type RGB struct {
	Red   int
	Green int
	Blue  int
}

func (rgb RGB) ToColor() uint32 {
	return uint32(uint32(rgb.Green)<<16 | uint32(rgb.Red)<<8 | uint32(rgb.Blue))
}

type Option func(*ws281x.ChannelOption)

type Controller struct {
	Logger  *slog.Logger
	Options []Option
}

func RGBToColor(r int, g int, b int) uint32 {
	return uint32(uint32(r)<<16 | uint32(g)<<8 | uint32(b))
}

func (ctrl *Controller) Render(drv *ws281x.WS2811, cats map[int]RGB) error {

	leds := drv.Leds(0)

	for i := 0; i < len(leds); i++ {
		rgb := cats[i]
		leds[i] = rgb.ToColor()

		if l := ctrl.Logger; l != nil {
			l.Debug("set color", "index", i, "color", leds[i])
		}
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
			opt.LedCount = DefaultLEDCount
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

func (ctrl *Controller) Serve(ctx context.Context, src chan (map[int]RGB)) error {

	drvopts := ws281x.DefaultOptions
	ctrl.applyOptions(&drvopts, ctrl.DefaultOptions()...)
	ctrl.applyOptions(&drvopts, ctrl.Options...)

	if l := ctrl.Logger; l != nil {
		l.Info("serving", "brightness", drvopts.Channels[0].Brightness, "ledCount", drvopts.Channels[0].LedCount, "gpioPin", drvopts.Channels[0].GpioPin)
	}

	drv, err := ws281x.MakeWS2811(&drvopts)
	if err != nil {
		return fmt.Errorf("failed to create driver: %w", err)
	}

	if err := drv.Init(); err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	defer func() {

		if l := ctrl.Logger; l != nil {
			l.Info("stopping")
		}

		if err := ctrl.setAllLEDs(ctx, drv, Off); err != nil {
			if l := ctrl.Logger; l != nil {
				l.Error("failed turning off LEDs", "error", err)
			}
		}

		drv.Fini()

		if l := ctrl.Logger; l != nil {
			l.Info("stopped serving")
		}
	}()

	for {
		select {
		case colors := <-src:
			if l := ctrl.Logger; l != nil {
				l.Debug("render", "categories", colors)
			}
			if err := ctrl.Render(drv, colors); err != nil {
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

func (ctrl *Controller) setAllLEDs(ctx context.Context, drv *ws281x.WS2811, color RGB) error {
	leds := drv.Leds(0)

	for i := 0; i < len(leds); i++ {
		leds[i] = color.ToColor()

		if l := ctrl.Logger; l != nil {
			l.Debug("set color", "index", i, "color", leds[i])
		}
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

func (ctrl *Controller) SetAllLEDs(ctx context.Context, color RGB) error {

	drvopts := ws281x.DefaultOptions
	ctrl.applyOptions(&drvopts, ctrl.DefaultOptions()...)
	ctrl.applyOptions(&drvopts, ctrl.Options...)

	if l := ctrl.Logger; l != nil {
		l.Info("setting all LEDs", "brightness", drvopts.Channels[0].Brightness, "ledCount", drvopts.Channels[0].LedCount, "gpioPin", drvopts.Channels[0].GpioPin)
	}

	drv, err := ws281x.MakeWS2811(&drvopts)
	if err != nil {
		return fmt.Errorf("failed to create driver: %w", err)
	}

	if err := drv.Init(); err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	defer func() {
		drv.Fini()
	}()

	return ctrl.setAllLEDs(ctx, drv, color)
}
