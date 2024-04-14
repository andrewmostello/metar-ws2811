package metar_test

import (
	"context"
	"testing"
	"time"

	"github.com/andrewmostello/metar-ws2811/metar"
)

func TestGetMETARs(t *testing.T) {

	c := metar.Client{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ids := []string{
		"KFIT",
		"KECU",
		"KRNM",
		"KJCT",
		"PATQ",
		"PASM",
		"KBUF",
		"CZSJ",
		"KFIT",
		"KECU",
		"KRNM",
		"KJCT",
		"PATQ",
		"PASM",
		"KBUF",
	}

	out, err := c.GetMETARs(ctx, ids...)
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range out {
		t.Logf("%s: %+v", k, v.FlightCategory().Name())
	}
}
