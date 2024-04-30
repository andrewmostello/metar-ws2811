package metar

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/andrewmostello/metar-ws2811/ws2811"
	"github.com/robfig/cron/v3"
)

// Colors  map[metar.FlightCategory]RGB

var (
	DefaultColors = map[FlightCategory]ws2811.RGB{
		FlightCategoryUnknown: ws2811.RGB{
			Red:   0,
			Green: 0,
			Blue:  0,
		},
		FlightCategoryVFR: ws2811.RGB{
			Red:   0,
			Green: 255,
			Blue:  0,
		},
		FlightCategoryMVFR: ws2811.RGB{
			Red:   0,
			Green: 0,
			Blue:  255,
		},
		FlightCategoryIFR: ws2811.RGB{
			Red:   255,
			Green: 0,
			Blue:  0,
		},
		FlightCategoryLIFR: ws2811.RGB{
			Red:   255,
			Green: 0,
			Blue:  255,
		},
	}
)

func FlightCategoryToRGB(colors map[FlightCategory]ws2811.RGB, cats map[int]FlightCategory) map[int]ws2811.RGB {
	if colors == nil {
		colors = DefaultColors
	}
	out := make(map[int]ws2811.RGB, len(cats))
	for i, cat := range cats {
		out[i] = colors[cat]
	}
	return out
}

type ColorServer struct {
	Logger              *slog.Logger
	Colors              map[FlightCategory]ws2811.RGB
	AirportIDs          []string
	LEDIndexByAirportID map[string]int
	Timeout             time.Duration
	Client              Client
}

func (srv *ColorServer) log(f func(l *slog.Logger)) {
	if srv.Logger == nil {
		return
	}
	f(srv.Logger.With("pkg", "metar", "func", "ColorServer.GetMETARs"))
}

func (srv *ColorServer) FlightCategoryToRGB(cats map[int]FlightCategory) map[int]ws2811.RGB {
	return FlightCategoryToRGB(srv.Colors, cats)
}

func (srv *ColorServer) timeout() time.Duration {
	if srv.Timeout > 0 {
		return srv.Timeout
	}
	return 15 * time.Second
}

func (srv *ColorServer) GetMETARs(ctx context.Context) (map[int]FlightCategory, error) {

	to := srv.timeout()

	ctx, cancel := context.WithTimeout(ctx, to)
	defer cancel()

	srv.log(func(l *slog.Logger) {
		l.Info("getting METARs", "timeout", to)
	})

	metars, err := srv.Client.GetMETARs(ctx, srv.AirportIDs...)
	if err != nil {
		srv.log(func(l *slog.Logger) {
			l.Error("failed to get METARs", "error", err)
		})
		return nil, fmt.Errorf("failed to get METARs: %w", err)
	}

	idxs := srv.LEDIndexByAirportID

	fcs := make(map[int]FlightCategory, len(idxs))

	for id, wx := range metars {
		idx, ok := idxs[id]
		if !ok {
			srv.log(func(l *slog.Logger) {
				l.Warn("no LED index for airport", "airport", id)
			})
			continue
		}
		fc := wx.FlightCategory()
		srv.log(func(l *slog.Logger) {
			l.Info("METAR", "airport", id, "index", idx, "flightCategory", fc.Name(), "weather", wx.RawObservation)
		})
		fcs[idx] = fc
	}

	return fcs, nil
}

func (srv *ColorServer) Serve(ctx context.Context, scd cron.Schedule, output chan (map[int]ws2811.RGB)) error {

	srv.log(func(l *slog.Logger) {
		l.Info("serving METARs", "airports", srv.AirportIDs)
	})

	defer func() {
		srv.log(func(l *slog.Logger) {
			l.Info("stopped serving METARs")
		})
	}()

	for {
		nxt := scd.Next(time.Now())

		srv.log(func(l *slog.Logger) {
			l.Info("scheduled next refresh", "next", nxt)
		})

		t := time.NewTimer(time.Until(nxt))

		select {
		case <-t.C:

			fcs, err := srv.GetMETARs(ctx)
			if err != nil {
				srv.log(func(l *slog.Logger) {
					l.Error("failed refresh", "error", err)
				})
			}

			output <- srv.FlightCategoryToRGB(fcs)

		case <-ctx.Done():
			srv.log(func(l *slog.Logger) {
				l.Info("stopping")
			})
			if !t.Stop() {
				<-t.C
			}
			return nil
		}
	}
}
