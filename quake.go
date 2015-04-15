package msg

import (
	"log"
	"time"
)

var alertAge = time.Duration(-60) * time.Minute

type Quake struct {
	PublicID              string
	Type                  string
	AgencyID              string
	ModificationTime      time.Time
	Time                  time.Time
	Latitude              float64
	Longitude             float64
	Depth                 float64
	DepthType             string
	MethodID              string
	EarthModelID          string
	EvaluationMode        string
	EvaluationStatus      string
	UsedPhaseCount        int
	UsedStationCount      int
	StandardError         float64
	AzimuthalGap          float64
	MinimumDistance       float64
	Magnitude             float64
	MagnitudeUncertainty  float64
	MagnitudeType         string
	MagnitudeStationCount int
	err                   error
}

func (q *Quake) Err() error {
	return q.err
}

func (q *Quake) SetErr(err error) {
	q.err = err
}

func (q *Quake) RxLog() {
	if q.err != nil {
		return
	}

	log.Printf("Received quake %s", q.PublicID)
}

func (q *Quake) TxLog() {
	if q.err != nil {
		return
	}

	log.Printf("Sending quake %s", q.PublicID)
}
