package metar

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type FlightCategory int

const (
	FlightCategoryUnknown FlightCategory = iota
	FlightCategoryLIFR
	FlightCategoryIFR
	FlightCategoryMVFR
	FlightCategoryVFR
)

func (c FlightCategory) String() string {
	return c.Name()
}

func (c FlightCategory) Name() string {
	switch c {
	case FlightCategoryVFR:
		return "VFR"
	case FlightCategoryMVFR:
		return "Marginal VFR"
	case FlightCategoryIFR:
		return "IFR"
	case FlightCategoryLIFR:
		return "Low IFR"
	}
	return "Unknown"
}

func (c FlightCategory) IsWorseThan(oth FlightCategory) bool {
	return c < oth
}

type METAR struct {
	ID                    int64         `json:"metar_id"`
	ICAOID                string        `json:"icaoId"`
	ReceiptTime           Time          `json:"receiptTime"`
	ObservationTime       Time          `json:"obsTime"`
	ReportTime            Time          `json:"reportTime"`
	Temperature           float64       `json:"temp"`
	Dewpoint              float64       `json:"dewp"`
	WindDirection         WindDirection `json:"wdir"`
	WindSpeed             float64       `json:"wspd"`
	WindGust              float64       `json:"wgst"`
	Visibility            *Visibility   `json:"visib"`
	Altimeter             float64       `json:"altim"`
	SeaLevelPressure      float64       `json:"slp"`
	QCField               float64       `json:"qcField"`
	WxString              string        `json:"wxString"`
	PressureTendency      *float64      `json:"presTend"`
	MaxTemperature        *float64      `json:"maxT"`
	MinTemperature        *float64      `json:"minT"`
	MaxTemperature24Hours *float64      `json:"maxT24"`
	MinTemperature24Hours *float64      `json:"minT24"`
	Precipitation         *float64      `json:"precip"`
	Precipitation3Hour    *float64      `json:"pcp3hr"`
	Precipitation6Hour    *float64      `json:"pcp6hr"`
	Precipitation24Hour   *float64      `json:"pcp24hr"`
	Snow                  *float64      `json:"snow"`
	VerticalVisibility    *float64      `json:"vertVis"`
	MetarType             METARType     `json:"metarType"`
	RawObservation        string        `json:"rawOb"`
	MostRecent            float64       `json:"mostRecent"`
	Latitude              float64       `json:"lat"`
	Longitude             float64       `json:"lon"`
	Elevation             float64       `json:"elev"`
	Prior                 float64       `json:"prior"`
	Name                  string        `json:"name"`
	Clouds                []CloudLayer  `json:"clouds"`
}

func (m METAR) FlightCategory() FlightCategory {
	vis := m.Visibility
	if vis == nil {
		return FlightCategoryUnknown
	}
	out := vis.FlightCategory()
	for _, lyr := range m.Clouds {
		c := lyr.FlightCategory()
		if c.IsWorseThan(out) {
			out = c
		}
	}
	return out
}

type CloudCover string

const (
	CloudCoverClear     = "CLR"
	CloudCoverFew       = "FEW"
	CloudCoverScattered = "SCT"
	CloudCoverBroken    = "BKN"
	CloudCoverOvercast  = "OVC"
	CloudCoverObscured  = "OVX"
)

func (m CloudCover) String() string {
	return string(m)
}

func (m CloudCover) IsCeiling() bool {
	switch m {
	case CloudCoverBroken, CloudCoverOvercast, CloudCoverObscured:
		return true
	}
	return false
}

func (m CloudCover) Name() string {
	switch m {
	case CloudCoverClear:
		return "Clear"
	case CloudCoverFew:
		return "Few"
	case CloudCoverScattered:
		return "Scattered"
	case CloudCoverBroken:
		return "Broken"
	case CloudCoverOvercast:
		return "Overcast"
	case CloudCoverObscured:
		return "Obscured"
	}
	return "Unknown"
}

type CloudLayer struct {
	Cover CloudCover `json:"cover"`
	Base  *float64   `json:"base"`
}

func (lyr CloudLayer) FlightCategory() FlightCategory {
	if !lyr.Cover.IsCeiling() {
		return FlightCategoryVFR
	}
	if lyr.Base == nil {
		return FlightCategoryUnknown
	}
	base := *lyr.Base
	switch {
	case base > 3000:
		return FlightCategoryVFR
	case base > 1000:
		return FlightCategoryMVFR
	case base > 500:
		return FlightCategoryIFR
	}
	return FlightCategoryLIFR
}

type METARType string

const (
	METARTypeMETAR   = "METAR"
	METARTypeSpecial = "SPECI"
)

type Visibility struct {
	Visibility  float64 `json:"vis"`
	GreaterThan bool    `json:"gt"`
}

func (v *Visibility) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	if b[0] != '"' {
		return json.Unmarshal(b, &v.Visibility)
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	if strings.HasSuffix(s, "+") {
		v.GreaterThan = true
		s = strings.TrimSuffix(s, "+")
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return &json.UnmarshalTypeError{
			Value: fmt.Sprintf(`string "%s"`, s),
			Type:  reflect.TypeOf(v),
		}
	}
	v.Visibility = f
	return nil
}

func (v Visibility) MarshalJSON() ([]byte, error) {
	s := strconv.FormatFloat(v.Visibility, 'f', -1, 64)
	if v.GreaterThan {
		return []byte(`"` + s + "+" + `"`), nil
	}
	return []byte(s), nil
}

func (v Visibility) String() string {
	vis := strconv.FormatFloat(v.Visibility, 'f', -1, 64)
	if v.GreaterThan {
		return vis + "+"
	}
	return vis
}

func (v Visibility) FlightCategory() FlightCategory {
	if v.Visibility > 5 {
		return FlightCategoryVFR
	} else if v.Visibility >= 3 {
		return FlightCategoryMVFR
	} else if v.Visibility >= 1 {
		return FlightCategoryIFR
	}
	return FlightCategoryLIFR
}

type Time time.Time

func (t *Time) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	if b[0] != '"' {
		var unix int64
		if err := json.Unmarshal(b, &unix); err != nil {
			return err
		}
		if unix == 0 {
			return nil
		}
		*t = Time(time.Unix(unix, 0).UTC())
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	tt, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return err
	}
	*t = Time(tt.UTC())
	return nil
}

type WindDirection struct {
	From     int32 `json:"fm"`
	Variable bool  `json:"vrb"`
}

func (wd WindDirection) String() string {
	if wd.Variable {
		return "VRB"
	}
	return strconv.Itoa(int(wd.From))
}

func (wd *WindDirection) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	if b[0] != '"' {
		return json.Unmarshal(b, &wd.From)
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	if s == "VRB" {
		wd.Variable = true
		return nil
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	wd.From = int32(i)
	return nil
}

func (wd WindDirection) MarshalJSON() ([]byte, error) {
	if wd.Variable {
		return []byte(`"VRB"`), nil
	}
	return json.Marshal(wd.From)
}
