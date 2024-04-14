package metar

import "time"

type METAR struct {
	METARID               string       `json:"metar_id"`
	ICAOID                string       `json:"icaoId"`
	ReceiptTime           time.Time    `json:"receiptTime"`
	ObservationTime       int64        `json:"obsTime"`
	ReportTime            time.Time    `json:"reportTime"`
	Temperature           string       `json:"temp"`
	Dewpoint              float64      `json:"dewp"`
	WindDirection         float64      `json:"wdir"`
	WindSpeed             float64      `json:"wspd"`
	WindGust              float64      `json:"wgst"`
	Visibility            string       `json:"visib"`
	Altimeter             float64      `json:"altim"`
	SeaLevelPressure      float64      `json:"slp"`
	QCField               string       `json:"qcField"`
	WxString              string       `json:"wxString"`
	PressureTendency      string       `json:"presTend"`
	MaxTemperature        string       `json:"maxT"`
	MinTemperature        string       `json:"minT"`
	MaxTemperature24Hours string       `json:"maxT24"`
	MinTemperature24Hours string       `json:"minT24"`
	Precipitation         string       `json:"precip"`
	Precipitation3Hour    string       `json:"pcp3hr"`
	Precipitation6Hour    string       `json:"pcp6hr"`
	Precipitation24Hour   string       `json:"pcp24hr"`
	Snow                  string       `json:"snow"`
	VerticalVisibility    string       `json:"vertVis"`
	MetarType             string       `json:"metarType"`
	RawObservation        string       `json:"rawOb"`
	MostRecent            string       `json:"mostRecent"`
	Latitude              float64      `json:"lat"`
	Longitude             float64      `json:"lon"`
	Elevation             float64      `json:"elev"`
	Prior                 float64      `json:"prior"`
	Name                  string       `json:"name"`
	Clouds                []METARCloud `json:"clouds"`
}

type METARCloudCover string

const (
	METARCloudCoverClear     = "CLR"
	METARCloudCoverFew       = "FEW"
	METARCloudCoverScattered = "SCT"
	METARCloudCoverBroken    = "BKN"
	METARCloudCoverOvercast  = "OVC"
)

func (m METARCloudCover) String() string {
	return string(m)
}

func (m METARCloudCover) IsCeiling() bool {
	switch m {
	case METARCloudCoverBroken, METARCloudCoverOvercast:
		return true
	}
	return false
}

func (m METARCloudCover) Name() string {
	switch m {
	case METARCloudCoverClear:
		return "Clear"
	case METARCloudCoverFew:
		return "Few"
	case METARCloudCoverScattered:
		return "Scattered"
	case METARCloudCoverBroken:
		return "Broken"
	case METARCloudCoverOvercast:
		return "Overcast"
	}
	return "Unknown"
}

type METARCloud struct {
	Cover string  `json:"cover"`
	Base  float64 `json:"base"`
}
