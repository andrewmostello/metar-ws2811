package metar_test

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/andrewmostello/metar-ws281x/metar"
)

func floatPtr(f float64) *float64 {
	return &f
}

func TestUnmarshalMETAR(t *testing.T) {
	type fixture struct {
		name string
		json string
		exp  metar.METAR
	}

	fixtures := []fixture{
		{
			name: "KFIT",
			json: `{"metar_id":529708844,"icaoId":"KFIT","receiptTime":"2024-04-14 13:56:21","obsTime":1713102720,"reportTime":"2024-04-14 14:00:00","temp":11.7,"dewp":-0.6,"wdir":260,"wspd":10,"wgst":15,"visib":"10+","altim":1010.9,"slp":1010.9,"qcField":6,"wxString":null,"presTend":null,"maxT":null,"minT":null,"maxT24":null,"minT24":null,"precip":null,"pcp3hr":null,"pcp6hr":null,"pcp24hr":null,"snow":null,"vertVis":null,"metarType":"METAR","rawOb":"KFIT 141352Z AUTO 26010G15KT 10SM CLR 12/M01 A2985 RMK AO2 SLP109 T01171006","mostRecent":1,"lat":42.5549,"lon":-71.757,"elev":102,"prior":5,"name":"Fitchburg Muni, MA, US","clouds":[{"cover":"CLR","base":null}]}`,
			exp: metar.METAR{
				ID:              529708844,
				ICAOID:          "KFIT",
				ReceiptTime:     metar.Time(time.Date(2024, 4, 14, 13, 56, 21, 0, time.UTC)),
				ObservationTime: metar.Time(time.Unix(1713102720, 0)),
				ReportTime:      metar.Time(time.Date(2024, 4, 14, 14, 0, 0, 0, time.UTC)),
				Temperature:     11.7,
				Dewpoint:        -0.6,
				WindDirection:   260,
				WindSpeed:       10,
				WindGust:        15,
				Visibility: metar.Visibility{
					Visibility:  10,
					GreaterThan: true,
				},
				Altimeter:             1010.9,
				SeaLevelPressure:      1010.9,
				QCField:               6,
				WxString:              "",
				PressureTendency:      nil,
				MaxTemperature:        nil,
				MinTemperature:        nil,
				MaxTemperature24Hours: nil,
				MinTemperature24Hours: nil,
				Precipitation:         nil,
				Precipitation3Hour:    nil,
				Precipitation6Hour:    nil,
				Precipitation24Hour:   nil,
				Snow:                  nil,
				VerticalVisibility:    nil,
				MetarType:             "METAR",
				RawObservation:        "KFIT 141352Z AUTO 26010G15KT 10SM CLR 12/M01 A2985 RMK AO2 SLP109 T01171006",
				MostRecent:            1,
				Latitude:              42.5549,
				Longitude:             -71.757,
				Elevation:             102,
				Prior:                 5,
				Name:                  "Fitchburg Muni, MA, US",
				Clouds: []metar.CloudLayer{
					{
						Cover: "CLR",
						Base:  nil,
					},
				},
			},
		},
		{
			name: "KECU",
			json: `{"metar_id":529727590,"icaoId":"KECU","receiptTime":"2024-04-14 14:42:06","obsTime":1713105300,"reportTime":"2024-04-14 14:35:00","temp":17,"dewp":15,"wdir":160,"wspd":9,"wgst":16,"visib":"10+","altim":1020.1,"slp":null,"qcField":6,"wxString":null,"presTend":null,"maxT":null,"minT":null,"maxT24":null,"minT24":null,"precip":null,"pcp3hr":null,"pcp6hr":null,"pcp24hr":null,"snow":null,"vertVis":null,"metarType":"METAR","rawOb":"KECU 141435Z AUTO 16009G16KT 10SM OVC007 17/15 A3012 RMK AO2","mostRecent":1,"lat":29.948,"lon":-100.174,"elev":725,"prior":4,"name":"Rocksprings/Edwards Cnty, TX, US","clouds":[{"cover":"OVC","base":700}]}`,
			exp: metar.METAR{
				ID:              529727590,
				ICAOID:          "KECU",
				ReceiptTime:     metar.Time(time.Date(2024, 4, 14, 14, 42, 6, 0, time.UTC)),
				ObservationTime: metar.Time(time.Unix(1713105300, 0)),
				ReportTime:      metar.Time(time.Date(2024, 4, 14, 14, 35, 0, 0, time.UTC)),
				Temperature:     17,
				Dewpoint:        15,
				WindDirection:   160,
				WindSpeed:       9,
				WindGust:        16,
				Visibility: metar.Visibility{
					Visibility:  10,
					GreaterThan: true,
				},
				Altimeter:             1020.1,
				SeaLevelPressure:      0,
				QCField:               6,
				WxString:              "",
				PressureTendency:      nil,
				MaxTemperature:        nil,
				MinTemperature:        nil,
				MaxTemperature24Hours: nil,
				MinTemperature24Hours: nil,
				Precipitation:         nil,
				Precipitation3Hour:    nil,
				Precipitation6Hour:    nil,
				Precipitation24Hour:   nil,
				Snow:                  nil,
				VerticalVisibility:    nil,
				MetarType:             "METAR",
				RawObservation:        "KECU 141435Z AUTO 16009G16KT 10SM OVC007 17/15 A3012 RMK AO2",
				MostRecent:            1,
				Latitude:              29.948,
				Longitude:             -100.174,
				Elevation:             725,
				Prior:                 4,
				Name:                  "Rocksprings/Edwards Cnty, TX, US",
				Clouds: []metar.CloudLayer{
					{
						Cover: "OVC",
						Base:  floatPtr(700),
					},
				},
			},
		},
		{
			name: "KRNM",
			json: `{"metar_id":529728418,"icaoId":"KRNM","receiptTime":"2024-04-14 14:50:07","obsTime":1713106080,"reportTime":"2024-04-14 14:48:00","temp":7,"dewp":7,"wdir":0,"wspd":0,"wgst":null,"visib":0.25,"altim":1021.1,"slp":null,"qcField":6,"wxString":"FG","presTend":null,"maxT":null,"minT":null,"maxT24":null,"minT24":null,"precip":null,"pcp3hr":null,"pcp6hr":null,"pcp24hr":null,"snow":null,"vertVis":200,"metarType":"SPECI","rawOb":"KRNM 141448Z AUTO 00000KT 1/4SM FG VV002 07/07 A3015 RMK AO2","mostRecent":1,"lat":33.038,"lon":-116.916,"elev":423,"prior":5,"name":"Ramona Arpt, CA, US","clouds":[{"cover":"OVX","base":0}]}`,
			exp: metar.METAR{
				ID:              529728418,
				ICAOID:          "KRNM",
				ReceiptTime:     metar.Time(time.Date(2024, 4, 14, 14, 50, 7, 0, time.UTC)),
				ObservationTime: metar.Time(time.Unix(1713106080, 0)),
				ReportTime:      metar.Time(time.Date(2024, 4, 14, 14, 48, 0, 0, time.UTC)),
				Temperature:     7,
				Dewpoint:        7,
				WindDirection:   0,
				WindSpeed:       0,
				WindGust:        0,
				Visibility: metar.Visibility{
					Visibility:  0.25,
					GreaterThan: false,
				},
				Altimeter:             1021.1,
				SeaLevelPressure:      0,
				QCField:               6,
				WxString:              "FG",
				PressureTendency:      nil,
				MaxTemperature:        nil,
				MinTemperature:        nil,
				MaxTemperature24Hours: nil,
				MinTemperature24Hours: nil,
				Precipitation:         nil,
				Precipitation3Hour:    nil,
				Precipitation6Hour:    nil,
				Precipitation24Hour:   nil,
				Snow:                  nil,
				VerticalVisibility:    floatPtr(200),
				MetarType:             "SPECI",
				RawObservation:        "KRNM 141448Z AUTO 00000KT 1/4SM FG VV002 07/07 A3015 RMK AO2",
				MostRecent:            1,
				Latitude:              33.038,
				Longitude:             -116.916,
				Elevation:             423,
				Prior:                 5,
				Name:                  "Ramona Arpt, CA, US",
				Clouds: []metar.CloudLayer{
					{
						Cover: "OVX",
						Base:  floatPtr(0),
					},
				},
			},
		},
		{
			name: "KJCT",
			json: `{"metar_id":529708362,"icaoId":"KJCT","receiptTime":"2024-04-14 13:56:11","obsTime":1713102660,"reportTime":"2024-04-14 14:00:00","temp":21.1,"dewp":16.1,"wdir":200,"wspd":7,"wgst":null,"visib":"10+","altim":1017.7,"slp":1014.2,"qcField":14,"wxString":null,"presTend":null,"maxT":null,"minT":null,"maxT24":null,"minT24":null,"precip":null,"pcp3hr":null,"pcp6hr":null,"pcp24hr":null,"snow":null,"vertVis":null,"metarType":"METAR","rawOb":"KJCT 141351Z AUTO 20007KT 10SM OVC019 21/16 A3005 RMK AO2 SLP142 T02110161 $","mostRecent":1,"lat":30.5105,"lon":-99.7665,"elev":522,"prior":3,"name":"Junction/Kimble Cnty, TX, US","clouds":[{"cover":"OVC","base":1900}]}`,
			exp: metar.METAR{
				ID:              529708362,
				ICAOID:          "KJCT",
				ReceiptTime:     metar.Time(time.Date(2024, 4, 14, 13, 56, 11, 0, time.UTC)),
				ObservationTime: metar.Time(time.Unix(1713102660, 0)),
				ReportTime:      metar.Time(time.Date(2024, 4, 14, 14, 0, 0, 0, time.UTC)),
				Temperature:     21.1,
				Dewpoint:        16.1,
				WindDirection:   200,
				WindSpeed:       7,
				WindGust:        0,
				Visibility: metar.Visibility{
					Visibility:  10,
					GreaterThan: true,
				},
				Altimeter:             1017.7,
				SeaLevelPressure:      1014.2,
				QCField:               14,
				WxString:              "",
				PressureTendency:      nil,
				MaxTemperature:        nil,
				MinTemperature:        nil,
				MaxTemperature24Hours: nil,
				MinTemperature24Hours: nil,
				Precipitation:         nil,
				Precipitation3Hour:    nil,
				Precipitation6Hour:    nil,
				Precipitation24Hour:   nil,
				Snow:                  nil,
				VerticalVisibility:    nil,
				MetarType:             "METAR",
				RawObservation:        "KJCT 141351Z AUTO 20007KT 10SM OVC019 21/16 A3005 RMK AO2 SLP142 T02110161 $",
				MostRecent:            1,
				Latitude:              30.5105,
				Longitude:             -99.7665,
				Elevation:             522,
				Prior:                 3,
				Name:                  "Junction/Kimble Cnty, TX, US",
				Clouds: []metar.CloudLayer{
					{
						Cover: "OVC",
						Base:  floatPtr(1900),
					},
				},
			},
		},
		{
			name: "PATQ",
			json: `{"metar_id":529741071,"icaoId":"PATQ","receiptTime":"2024-04-14 15:14:05","obsTime":1713107520,"reportTime":"2024-04-14 15:12:00","temp":-24,"dewp":-26,"wdir":230,"wspd":15,"wgst":null,"visib":3,"altim":1018.4,"slp":null,"qcField":102,"wxString":"BR","presTend":null,"maxT":null,"minT":null,"maxT24":null,"minT24":null,"precip":0.005,"pcp3hr":null,"pcp6hr":null,"pcp24hr":null,"snow":null,"vertVis":null,"metarType":"SPECI","rawOb":"PATQ 141512Z AUTO 23015KT 3SM BR OVC014 M24/M26 A3007 RMK AO2 SNE02 P0000 FZRANO TSNO","mostRecent":1,"lat":70.469,"lon":-157.428,"elev":28,"prior":2,"name":"Atqasuk/Burnell Mem, AK, US","clouds":[{"cover":"OVC","base":1400}]}`,
			exp: metar.METAR{
				ID:              529741071,
				ICAOID:          "PATQ",
				ReceiptTime:     metar.Time(time.Date(2024, 4, 14, 15, 14, 5, 0, time.UTC)),
				ObservationTime: metar.Time(time.Unix(1713107520, 0)),
				ReportTime:      metar.Time(time.Date(2024, 4, 14, 15, 12, 0, 0, time.UTC)),
				Temperature:     -24,
				Dewpoint:        -26,
				WindDirection:   230,
				WindSpeed:       15,
				WindGust:        0,
				Visibility: metar.Visibility{
					Visibility:  3,
					GreaterThan: false,
				},
				Altimeter:             1018.4,
				SeaLevelPressure:      0,
				QCField:               102,
				WxString:              "BR",
				PressureTendency:      nil,
				MaxTemperature:        nil,
				MinTemperature:        nil,
				MaxTemperature24Hours: nil,
				MinTemperature24Hours: nil,
				Precipitation:         floatPtr(0.005),
				Precipitation3Hour:    nil,
				Precipitation6Hour:    nil,
				Precipitation24Hour:   nil,
				Snow:                  nil,
				VerticalVisibility:    nil,
				MetarType:             "SPECI",
				RawObservation:        "PATQ 141512Z AUTO 23015KT 3SM BR OVC014 M24/M26 A3007 RMK AO2 SNE02 P0000 FZRANO TSNO",
				MostRecent:            1,
				Latitude:              70.469,
				Longitude:             -157.428,
				Elevation:             28,
				Prior:                 2,
				Name:                  "Atqasuk/Burnell Mem, AK, US",
				Clouds: []metar.CloudLayer{
					{
						Cover: "OVC",
						Base:  floatPtr(1400),
					},
				},
			},
		},
		{
			name: "PASM",
			json: `{"metar_id":529744489,"icaoId":"PASM","receiptTime":"2024-04-14 15:24:05","obsTime":1713108120,"reportTime":"2024-04-14 15:22:00","temp":-3,"dewp":-3,"wdir":260,"wspd":4,"wgst":null,"visib":0.25,"altim":1024.1,"slp":null,"qcField":70,"wxString":"FZFG","presTend":null,"maxT":null,"minT":null,"maxT24":null,"minT24":null,"precip":null,"pcp3hr":null,"pcp6hr":null,"pcp24hr":null,"snow":null,"vertVis":null,"metarType":"SPECI","rawOb":"PASM 141522Z AUTO 26004KT 1/4SM FZFG OVC003 M03/M03 A3024 RMK AO2 FZRANO","mostRecent":1,"lat":62.057,"lon":-163.298,"elev":108,"prior":2,"name":"St Marys Arpt, AK, US","clouds":[{"cover":"OVC","base":300}]}`,
			exp: metar.METAR{
				ID:              529744489,
				ICAOID:          "PASM",
				ReceiptTime:     metar.Time(time.Date(2024, 4, 14, 15, 24, 5, 0, time.UTC)),
				ObservationTime: metar.Time(time.Unix(1713108120, 0)),
				ReportTime:      metar.Time(time.Date(2024, 4, 14, 15, 22, 0, 0, time.UTC)),
				Temperature:     -3,
				Dewpoint:        -3,
				WindDirection:   260,
				WindSpeed:       4,
				WindGust:        0,
				Visibility: metar.Visibility{
					Visibility:  0.25,
					GreaterThan: false,
				},
				Altimeter:             1024.1,
				SeaLevelPressure:      0,
				QCField:               70,
				WxString:              "FZFG",
				PressureTendency:      nil,
				MaxTemperature:        nil,
				MinTemperature:        nil,
				MaxTemperature24Hours: nil,
				MinTemperature24Hours: nil,
				Precipitation:         nil,
				Precipitation3Hour:    nil,
				Precipitation6Hour:    nil,
				Precipitation24Hour:   nil,
				Snow:                  nil,
				VerticalVisibility:    nil,
				MetarType:             "SPECI",
				RawObservation:        "PASM 141522Z AUTO 26004KT 1/4SM FZFG OVC003 M03/M03 A3024 RMK AO2 FZRANO",
				MostRecent:            1,
				Latitude:              62.057,
				Longitude:             -163.298,
				Elevation:             108,
				Prior:                 2,
				Name:                  "St Marys Arpt, AK, US",
				Clouds: []metar.CloudLayer{
					{
						Cover: "OVC",
						Base:  floatPtr(300),
					},
				},
			},
		},
		{
			name: "KBUF",
			json: `{"metar_id":529733443,"icaoId":"KBUF","receiptTime":"2024-04-14 14:58:16","obsTime":1713106440,"reportTime":"2024-04-14 15:00:00","temp":11.7,"dewp":6.7,"wdir":180,"wspd":9,"wgst":null,"visib":"10+","altim":1004.5,"slp":1004.6,"qcField":4,"wxString":null,"presTend":-4.3,"maxT":null,"minT":null,"maxT24":null,"minT24":null,"precip":null,"pcp3hr":0.1,"pcp6hr":null,"pcp24hr":null,"snow":null,"vertVis":null,"metarType":"METAR","rawOb":"KBUF 141454Z 18009KT 10SM FEW050 SCT070 BKN140 BKN200 12/07 A2966 RMK AO2 SLP046 60010 T01170067 58043","mostRecent":1,"lat":42.94,"lon":-78.7361,"elev":217,"prior":1,"name":"Buffalo-Niagara Intl, NY, US","clouds":[{"cover":"FEW","base":5000},{"cover":"SCT","base":7000},{"cover":"BKN","base":14000},{"cover":"BKN","base":20000}]}`,
			exp: metar.METAR{
				ID:              529733443,
				ICAOID:          "KBUF",
				ReceiptTime:     metar.Time(time.Date(2024, 4, 14, 14, 58, 16, 0, time.UTC)),
				ObservationTime: metar.Time(time.Unix(1713106440, 0)),
				ReportTime:      metar.Time(time.Date(2024, 4, 14, 15, 0, 0, 0, time.UTC)),
				Temperature:     11.7,
				Dewpoint:        6.7,
				WindDirection:   180,
				WindSpeed:       9,
				WindGust:        0,
				Visibility: metar.Visibility{
					Visibility:  10,
					GreaterThan: true,
				},
				Altimeter:             1004.5,
				SeaLevelPressure:      1004.6,
				QCField:               4,
				WxString:              "",
				PressureTendency:      floatPtr(-4.3),
				MaxTemperature:        nil,
				MinTemperature:        nil,
				MaxTemperature24Hours: nil,
				MinTemperature24Hours: nil,
				Precipitation:         nil,
				Precipitation3Hour:    floatPtr(0.1),
				Precipitation6Hour:    nil,
				Precipitation24Hour:   nil,
				Snow:                  nil,
				VerticalVisibility:    nil,
				MetarType:             "METAR",
				RawObservation:        "KBUF 141454Z 18009KT 10SM FEW050 SCT070 BKN140 BKN200 12/07 A2966 RMK AO2 SLP046 60010 T01170067 58043",
				MostRecent:            1,
				Latitude:              42.94,
				Longitude:             -78.7361,
				Elevation:             217,
				Prior:                 1,
				Name:                  "Buffalo-Niagara Intl, NY, US",
				Clouds: []metar.CloudLayer{
					{
						Cover: "FEW",
						Base:  floatPtr(5000),
					},
					{
						Cover: "SCT",
						Base:  floatPtr(7000),
					},
					{
						Cover: "BKN",
						Base:  floatPtr(14000),
					},
					{
						Cover: "BKN",
						Base:  floatPtr(20000),
					},
				},
			},
		},
	}

	for _, f := range fixtures {
		t.Run(f.name, func(t *testing.T) {
			var m metar.METAR

			if err := json.Unmarshal([]byte(f.json), &m); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if reflect.DeepEqual(m, f.exp) {
				t.Fatalf("expected %v, got %v", f.exp, m)
			}
		})
	}
}

func TestMETARConditions(t *testing.T) {
	type fixture struct {
		name string
		wx   metar.METAR
		exp  metar.Conditions
	}

	fixtures := []fixture{
		{
			name: "KFIT",
			exp:  metar.ConditionsVFR,
			wx: metar.METAR{
				ID:              529708844,
				ICAOID:          "KFIT",
				ReceiptTime:     metar.Time(time.Date(2024, 4, 14, 13, 56, 21, 0, time.UTC)),
				ObservationTime: metar.Time(time.Unix(1713102720, 0)),
				ReportTime:      metar.Time(time.Date(2024, 4, 14, 14, 0, 0, 0, time.UTC)),
				Temperature:     11.7,
				Dewpoint:        -0.6,
				WindDirection:   260,
				WindSpeed:       10,
				WindGust:        15,
				Visibility: metar.Visibility{
					Visibility:  10,
					GreaterThan: true,
				},
				Altimeter:             1010.9,
				SeaLevelPressure:      1010.9,
				QCField:               6,
				WxString:              "",
				PressureTendency:      nil,
				MaxTemperature:        nil,
				MinTemperature:        nil,
				MaxTemperature24Hours: nil,
				MinTemperature24Hours: nil,
				Precipitation:         nil,
				Precipitation3Hour:    nil,
				Precipitation6Hour:    nil,
				Precipitation24Hour:   nil,
				Snow:                  nil,
				VerticalVisibility:    nil,
				MetarType:             "METAR",
				RawObservation:        "KFIT 141352Z AUTO 26010G15KT 10SM CLR 12/M01 A2985 RMK AO2 SLP109 T01171006",
				MostRecent:            1,
				Latitude:              42.5549,
				Longitude:             -71.757,
				Elevation:             102,
				Prior:                 5,
				Name:                  "Fitchburg Muni, MA, US",
				Clouds: []metar.CloudLayer{
					{
						Cover: "CLR",
						Base:  nil,
					},
				},
			},
		},
		{
			name: "KECU",
			exp:  metar.ConditionsIFR,
			wx: metar.METAR{
				ID:              529727590,
				ICAOID:          "KECU",
				ReceiptTime:     metar.Time(time.Date(2024, 4, 14, 14, 42, 6, 0, time.UTC)),
				ObservationTime: metar.Time(time.Unix(1713105300, 0)),
				ReportTime:      metar.Time(time.Date(2024, 4, 14, 14, 35, 0, 0, time.UTC)),
				Temperature:     17,
				Dewpoint:        15,
				WindDirection:   160,
				WindSpeed:       9,
				WindGust:        16,
				Visibility: metar.Visibility{
					Visibility:  10,
					GreaterThan: true,
				},
				Altimeter:             1020.1,
				SeaLevelPressure:      0,
				QCField:               6,
				WxString:              "",
				PressureTendency:      nil,
				MaxTemperature:        nil,
				MinTemperature:        nil,
				MaxTemperature24Hours: nil,
				MinTemperature24Hours: nil,
				Precipitation:         nil,
				Precipitation3Hour:    nil,
				Precipitation6Hour:    nil,
				Precipitation24Hour:   nil,
				Snow:                  nil,
				VerticalVisibility:    nil,
				MetarType:             "METAR",
				RawObservation:        "KECU 141435Z AUTO 16009G16KT 10SM OVC007 17/15 A3012 RMK AO2",
				MostRecent:            1,
				Latitude:              29.948,
				Longitude:             -100.174,
				Elevation:             725,
				Prior:                 4,
				Name:                  "Rocksprings/Edwards Cnty, TX, US",
				Clouds: []metar.CloudLayer{
					{
						Cover: "OVC",
						Base:  floatPtr(700),
					},
				},
			},
		},
		{
			name: "KRNM",
			exp:  metar.ConditionsLIFR,
			wx: metar.METAR{
				ID:              529728418,
				ICAOID:          "KRNM",
				ReceiptTime:     metar.Time(time.Date(2024, 4, 14, 14, 50, 7, 0, time.UTC)),
				ObservationTime: metar.Time(time.Unix(1713106080, 0)),
				ReportTime:      metar.Time(time.Date(2024, 4, 14, 14, 48, 0, 0, time.UTC)),
				Temperature:     7,
				Dewpoint:        7,
				WindDirection:   0,
				WindSpeed:       0,
				WindGust:        0,
				Visibility: metar.Visibility{
					Visibility:  0.25,
					GreaterThan: false,
				},
				Altimeter:             1021.1,
				SeaLevelPressure:      0,
				QCField:               6,
				WxString:              "FG",
				PressureTendency:      nil,
				MaxTemperature:        nil,
				MinTemperature:        nil,
				MaxTemperature24Hours: nil,
				MinTemperature24Hours: nil,
				Precipitation:         nil,
				Precipitation3Hour:    nil,
				Precipitation6Hour:    nil,
				Precipitation24Hour:   nil,
				Snow:                  nil,
				VerticalVisibility:    floatPtr(200),
				MetarType:             "SPECI",
				RawObservation:        "KRNM 141448Z AUTO 00000KT 1/4SM FG VV002 07/07 A3015 RMK AO2",
				MostRecent:            1,
				Latitude:              33.038,
				Longitude:             -116.916,
				Elevation:             423,
				Prior:                 5,
				Name:                  "Ramona Arpt, CA, US",
				Clouds: []metar.CloudLayer{
					{
						Cover: "OVX",
						Base:  floatPtr(0),
					},
				},
			},
		},
		{
			name: "KJCT",
			exp:  metar.ConditionsMVFR,
			wx: metar.METAR{
				ID:              529708362,
				ICAOID:          "KJCT",
				ReceiptTime:     metar.Time(time.Date(2024, 4, 14, 13, 56, 11, 0, time.UTC)),
				ObservationTime: metar.Time(time.Unix(1713102660, 0)),
				ReportTime:      metar.Time(time.Date(2024, 4, 14, 14, 0, 0, 0, time.UTC)),
				Temperature:     21.1,
				Dewpoint:        16.1,
				WindDirection:   200,
				WindSpeed:       7,
				WindGust:        0,
				Visibility: metar.Visibility{
					Visibility:  10,
					GreaterThan: true,
				},
				Altimeter:             1017.7,
				SeaLevelPressure:      1014.2,
				QCField:               14,
				WxString:              "",
				PressureTendency:      nil,
				MaxTemperature:        nil,
				MinTemperature:        nil,
				MaxTemperature24Hours: nil,
				MinTemperature24Hours: nil,
				Precipitation:         nil,
				Precipitation3Hour:    nil,
				Precipitation6Hour:    nil,
				Precipitation24Hour:   nil,
				Snow:                  nil,
				VerticalVisibility:    nil,
				MetarType:             "METAR",
				RawObservation:        "KJCT 141351Z AUTO 20007KT 10SM OVC019 21/16 A3005 RMK AO2 SLP142 T02110161 $",
				MostRecent:            1,
				Latitude:              30.5105,
				Longitude:             -99.7665,
				Elevation:             522,
				Prior:                 3,
				Name:                  "Junction/Kimble Cnty, TX, US",
				Clouds: []metar.CloudLayer{
					{
						Cover: "OVC",
						Base:  floatPtr(1900),
					},
				},
			},
		},
		{
			name: "PATQ",
			exp:  metar.ConditionsMVFR,
			wx: metar.METAR{
				ID:              529741071,
				ICAOID:          "PATQ",
				ReceiptTime:     metar.Time(time.Date(2024, 4, 14, 15, 14, 5, 0, time.UTC)),
				ObservationTime: metar.Time(time.Unix(1713107520, 0)),
				ReportTime:      metar.Time(time.Date(2024, 4, 14, 15, 12, 0, 0, time.UTC)),
				Temperature:     -24,
				Dewpoint:        -26,
				WindDirection:   230,
				WindSpeed:       15,
				WindGust:        0,
				Visibility: metar.Visibility{
					Visibility:  3,
					GreaterThan: false,
				},
				Altimeter:             1018.4,
				SeaLevelPressure:      0,
				QCField:               102,
				WxString:              "BR",
				PressureTendency:      nil,
				MaxTemperature:        nil,
				MinTemperature:        nil,
				MaxTemperature24Hours: nil,
				MinTemperature24Hours: nil,
				Precipitation:         floatPtr(0.005),
				Precipitation3Hour:    nil,
				Precipitation6Hour:    nil,
				Precipitation24Hour:   nil,
				Snow:                  nil,
				VerticalVisibility:    nil,
				MetarType:             "SPECI",
				RawObservation:        "PATQ 141512Z AUTO 23015KT 3SM BR OVC014 M24/M26 A3007 RMK AO2 SNE02 P0000 FZRANO TSNO",
				MostRecent:            1,
				Latitude:              70.469,
				Longitude:             -157.428,
				Elevation:             28,
				Prior:                 2,
				Name:                  "Atqasuk/Burnell Mem, AK, US",
				Clouds: []metar.CloudLayer{
					{
						Cover: "OVC",
						Base:  floatPtr(1400),
					},
				},
			},
		},
		{
			name: "PASM",
			exp:  metar.ConditionsLIFR,
			wx: metar.METAR{
				ID:              529744489,
				ICAOID:          "PASM",
				ReceiptTime:     metar.Time(time.Date(2024, 4, 14, 15, 24, 5, 0, time.UTC)),
				ObservationTime: metar.Time(time.Unix(1713108120, 0)),
				ReportTime:      metar.Time(time.Date(2024, 4, 14, 15, 22, 0, 0, time.UTC)),
				Temperature:     -3,
				Dewpoint:        -3,
				WindDirection:   260,
				WindSpeed:       4,
				WindGust:        0,
				Visibility: metar.Visibility{
					Visibility:  0.25,
					GreaterThan: false,
				},
				Altimeter:             1024.1,
				SeaLevelPressure:      0,
				QCField:               70,
				WxString:              "FZFG",
				PressureTendency:      nil,
				MaxTemperature:        nil,
				MinTemperature:        nil,
				MaxTemperature24Hours: nil,
				MinTemperature24Hours: nil,
				Precipitation:         nil,
				Precipitation3Hour:    nil,
				Precipitation6Hour:    nil,
				Precipitation24Hour:   nil,
				Snow:                  nil,
				VerticalVisibility:    nil,
				MetarType:             "SPECI",
				RawObservation:        "PASM 141522Z AUTO 26004KT 1/4SM FZFG OVC003 M03/M03 A3024 RMK AO2 FZRANO",
				MostRecent:            1,
				Latitude:              62.057,
				Longitude:             -163.298,
				Elevation:             108,
				Prior:                 2,
				Name:                  "St Marys Arpt, AK, US",
				Clouds: []metar.CloudLayer{
					{
						Cover: "OVC",
						Base:  floatPtr(300),
					},
				},
			},
		},
		{
			name: "KBUF",
			exp:  metar.ConditionsVFR,
			wx: metar.METAR{
				ID:              529733443,
				ICAOID:          "KBUF",
				ReceiptTime:     metar.Time(time.Date(2024, 4, 14, 14, 58, 16, 0, time.UTC)),
				ObservationTime: metar.Time(time.Unix(1713106440, 0)),
				ReportTime:      metar.Time(time.Date(2024, 4, 14, 15, 0, 0, 0, time.UTC)),
				Temperature:     11.7,
				Dewpoint:        6.7,
				WindDirection:   180,
				WindSpeed:       9,
				WindGust:        0,
				Visibility: metar.Visibility{
					Visibility:  10,
					GreaterThan: true,
				},
				Altimeter:             1004.5,
				SeaLevelPressure:      1004.6,
				QCField:               4,
				WxString:              "",
				PressureTendency:      floatPtr(-4.3),
				MaxTemperature:        nil,
				MinTemperature:        nil,
				MaxTemperature24Hours: nil,
				MinTemperature24Hours: nil,
				Precipitation:         nil,
				Precipitation3Hour:    floatPtr(0.1),
				Precipitation6Hour:    nil,
				Precipitation24Hour:   nil,
				Snow:                  nil,
				VerticalVisibility:    nil,
				MetarType:             "METAR",
				RawObservation:        "KBUF 141454Z 18009KT 10SM FEW050 SCT070 BKN140 BKN200 12/07 A2966 RMK AO2 SLP046 60010 T01170067 58043",
				MostRecent:            1,
				Latitude:              42.94,
				Longitude:             -78.7361,
				Elevation:             217,
				Prior:                 1,
				Name:                  "Buffalo-Niagara Intl, NY, US",
				Clouds: []metar.CloudLayer{
					{
						Cover: "FEW",
						Base:  floatPtr(5000),
					},
					{
						Cover: "SCT",
						Base:  floatPtr(7000),
					},
					{
						Cover: "BKN",
						Base:  floatPtr(14000),
					},
					{
						Cover: "BKN",
						Base:  floatPtr(20000),
					},
				},
			},
		},
	}

	for _, f := range fixtures {
		t.Run(f.name, func(t *testing.T) {
			cnd := f.wx.Conditions()
			if cnd != f.exp {
				t.Fatalf("expected %v, got %v", f.exp, cnd)
			}
		})
	}
}
