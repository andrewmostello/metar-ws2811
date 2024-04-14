package metar

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

const (
	DefaultBaseURL = "https://aviationweather.gov/api/data"
)

type Client struct {
	HTTPClient *http.Client
	BaseURL    string
}

func (c Client) Route(pth string) (*url.URL, error) {
	base := c.BaseURL
	if c.BaseURL == "" {
		base = DefaultBaseURL
	}

	u, err := url.Parse(base)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	u.Path = path.Join(u.Path, pth)

	return u, nil
}

func (c Client) Do(r *http.Request) (*http.Response, error) {
	hc := c.HTTPClient
	if hc == nil {
		hc = http.DefaultClient
	}
	return hc.Do(r)
}

func (c Client) GetMETARs(ctx context.Context, airportIDs ...string) (map[string]METAR, error) {

	if len(airportIDs) == 0 {
		return nil, fmt.Errorf("no airport identifiers specified")
	}

	u, err := c.Route("/metar")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("ids", strings.Join(airportIDs, ","))
	q.Set("format", "json")
	u.RawQuery = q.Encode()

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	r.Header.Set("accept", "application/json")

	resp, err := c.Do(r)
	if err != nil {
		return nil, fmt.Errorf("failed retrieving METAR(s): %w", err)
	}
	defer resp.Body.Close()

	var bdy []METAR

	bts, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if err := json.Unmarshal(bts, &bdy); err != nil {
		fmt.Println(string(bts))
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	out := make(map[string]METAR)

	for _, m := range bdy {
		out[m.ICAOID] = m
	}

	return out, nil
}
