package metar_test

import (
	"context"
	"testing"
	"time"

	"github.com/andrewmostello/metar-ws281x/metar"
)

func TestGetMETARs(t *testing.T) {

	c := metar.Client{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	out, err := c.GetMETARs(ctx, []string{"KJFK", "KBOS", "KFIT"})
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range out {
		t.Logf("%s: %+v", k, v.FlightCategory().Name())
	}

	t.Fail()
}
