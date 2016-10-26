/*
Package sc3ml is for converting SeisComPML to other formats.
*/
package sc3ml

import (
	"encoding/xml"
	"sort"
	"time"
)

type seiscomp struct {
	EventParameters eventParameters `xml:"EventParameters"`
}

type eventParameters struct {
	Events     []event     `xml:"event"`
	Picks      []pick      `xml:"pick"`
	Amplitudes []amplitude `xml:"amplitude"`
	Origins    []origin    `xml:"origin"`
}

type event struct {
	PublicID             string       `xml:"publicID,attr"`
	PreferredOriginID    string       `xml:"preferredOriginID"`
	PreferredMagnitudeID string       `xml:"preferredMagnitudeID"`
	Type                 string       `xml:"type"`
	CreationInfo         creationInfo `xml:"creationInfo"`
	PreferredOrigin      origin
	PreferredMagnitude   magnitude
}

type creationInfo struct {
	AgencyID         string    `xml:"agencyID"`
	CreationTime     time.Time `xml:"creationTime"`
	ModificationTime time.Time `xml:"modificationTime"`
}

type origin struct {
	PublicID          string             `xml:"publicID,attr"`
	Time              timeValue          `xml:"time"`
	Latitude          realQuantity       `xml:"latitude"`
	Longitude         realQuantity       `xml:"longitude"`
	Depth             realQuantity       `xml:"depth"`
	DepthType         string             `xml:"depthType"`
	MethodID          string             `xml:"methodID"`
	EarthModelID      string             `xml:"earthModelID"`
	Quality           quality            `xml:"quality"`
	EvaluationMode    string             `xml:"evaluationMode"`
	EvaluationStatus  string             `xml:"evaluationStatus"`
	Arrivals          []arrival          `xml:"arrival"`
	StationMagnitudes []stationMagnitude `xml:"stationMagnitude"`
	Magnitudes        []magnitude        `xml:"magnitude"`
	CreationInfo      creationInfo       `xml:"creationInfo"`
}

type quality struct {
	UsedPhaseCount   int64   `xml:"usedPhaseCount"`
	UsedStationCount int64   `xml:"usedStationCount"`
	StandardError    float64 `xml:"standardError"`
	AzimuthalGap     float64 `xml:"azimuthalGap"`
	MinimumDistance  float64 `xml:"minimumDistance"`
}

type arrival struct {
	PickID       string  `xml:"pickID"`
	Phase        string  `xml:"phase"`
	Azimuth      float64 `xml:"azimuth"`
	Distance     float64 `xml:"distance"`
	TimeResidual float64 `xml:"timeResidual"`
	Weight       float64 `xml:"weight"`
	Pick         pick
}

type pick struct {
	PublicID         string     `xml:"publicID,attr"`
	Time             timeValue  `xml:"time"`
	WaveformID       waveformID `xml:"waveformID"`
	EvaluationMode   string     `xml:"evaluationMode"`
	EvaluationStatus string     `xml:"evaluationStatus"`
}

type waveformID struct {
	NetworkCode  string `xml:"networkCode,attr"`
	StationCode  string `xml:"stationCode,attr"`
	LocationCode string `xml:"locationCode,attr"`
	ChannelCode  string `xml:"channelCode,attr"`
}

type realQuantity struct {
	Value       float64 `xml:"value"`
	Uncertainty float64 `xml:"uncertainty"`
}

type timeValue struct {
	Value time.Time `xml:"value"`
}

type magnitude struct {
	PublicID                      string                         `xml:"publicID,attr"`
	Magnitude                     realQuantity                   `xml:"magnitude"`
	Type                          string                         `xml:"type"`
	MethodID                      string                         `xml:"methodID"`
	StationCount                  int64                          `xml:"stationCount"`
	StationMagnitudeContributions []stationMagnitudeContribution `xml:"stationMagnitudeContribution"`
	CreationInfo                  creationInfo                   `xml:"creationInfo"`
}

type stationMagnitudeContribution struct {
	StationMagnitudeID string  `xml:"stationMagnitudeID"`
	Weight             float64 `xml:"weight"`
	Residual           float64 `xml:"residual"`
	StationMagnitude   stationMagnitude
}

type stationMagnitude struct {
	PublicID    string       `xml:"publicID,attr"`
	Magnitude   realQuantity `xml:"magnitude"`
	Type        string       `xml:"type"`
	AmplitudeID string       `xml:"amplitudeID"`
	WaveformID  waveformID   `xml:"waveformID"`
	Amplitude   amplitude
}

type amplitude struct {
	PublicID  string       `xml:"publicID,attr"`
	Amplitude realQuantity `xml:"amplitude"`
	PickID    string  `xml:"pickID"`
	Azimuth   float64 // not in the SC3ML - will be mapped from arrival using PickID
	Distance  float64 // not in the SC3ML - will be mapped from arrival using PickID
}

// unmarshal unmarshals the SeisComPML in b and initialises all
// the objects referenced by ID in the SeisComPML e.g., PreferredOrigin,
// PreferredMagnitude etc.
func unmarshal(b []byte) (eventParameters, error) {
	var q seiscomp

	if err := xml.Unmarshal(b, &q); err != nil {
		return q.EventParameters, err
	}

	var picks = make(map[string]pick)
	for k, v := range q.EventParameters.Picks {
		picks[v.PublicID] = q.EventParameters.Picks[k]
	}

	var arrivals = make(map[string]arrival)
	for i := range q.EventParameters.Origins {
		for _, v := range q.EventParameters.Origins[i].Arrivals {
			arrivals[v.PickID] = v
		}
	}

	var amplitudes = make(map[string]amplitude)
	for k, v := range q.EventParameters.Amplitudes {
		a := q.EventParameters.Amplitudes[k]

		// add distance and azimuth from the arrival with the matching PickID.
		pk := arrivals[v.PickID]

		a.Distance = pk.Distance
		a.Azimuth = pk.Azimuth

		amplitudes[v.PublicID] = a
	}

	for i := range q.EventParameters.Origins {
		for k, v := range q.EventParameters.Origins[i].Arrivals {
			q.EventParameters.Origins[i].Arrivals[k].Pick = picks[v.PickID]
		}

		var stationMagnitudes = make(map[string]stationMagnitude)

		for k, v := range q.EventParameters.Origins[i].StationMagnitudes {
			q.EventParameters.Origins[i].StationMagnitudes[k].Amplitude = amplitudes[v.AmplitudeID]
			stationMagnitudes[v.PublicID] = q.EventParameters.Origins[i].StationMagnitudes[k]
		}

		for j := range q.EventParameters.Origins[i].Magnitudes {
			for k, v := range q.EventParameters.Origins[i].Magnitudes[j].StationMagnitudeContributions {
				q.EventParameters.Origins[i].Magnitudes[j].StationMagnitudeContributions[k].StationMagnitude = stationMagnitudes[v.StationMagnitudeID]
			}
		}
	}

	// set the preferred origin.
	// set the preferred mag which can come from any origin
	for i := range q.EventParameters.Events {
		for k, v := range q.EventParameters.Origins {
			if v.PublicID == q.EventParameters.Events[i].PreferredOriginID {
				q.EventParameters.Events[i].PreferredOrigin = q.EventParameters.Origins[k]
			}
			for _, mag := range v.Magnitudes {
				if mag.PublicID == q.EventParameters.Events[i].PreferredMagnitudeID {
					q.EventParameters.Events[i].PreferredMagnitude = mag
				}
			}
		}
	}

	return q.EventParameters, nil
}

// modificationTime returns the most recent creation or modification time
// for the Event, PreferredOrigin, or PreferredMagnitude.
func (e *event) modificationTime() time.Time {
	var t []string

	t = append(t, e.CreationInfo.CreationTime.Format(time.RFC3339Nano))
	t = append(t, e.CreationInfo.ModificationTime.Format(time.RFC3339Nano))
	t = append(t, e.PreferredOrigin.CreationInfo.CreationTime.Format(time.RFC3339Nano))
	t = append(t, e.PreferredOrigin.CreationInfo.ModificationTime.Format(time.RFC3339Nano))
	t = append(t, e.PreferredMagnitude.CreationInfo.CreationTime.Format(time.RFC3339Nano))
	t = append(t, e.PreferredMagnitude.CreationInfo.ModificationTime.Format(time.RFC3339Nano))

	sort.Sort(sort.Reverse(sort.StringSlice(t)))

	tm, _ := time.Parse(time.RFC3339Nano, t[0])
	return tm
}
