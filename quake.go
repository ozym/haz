package msg

import (
	"log"
	"time"
)

type Quake struct {
	PublicID              string
	Type                  string
	AgencyID              string
	ModificationTime      time.Time
	Time                  time.Time
	Latitude              float64
	Longitude             float64
	Depth                 float64
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

// Status returns the public status for the Quake referred to by q.
// Returns 'error' if q.Err() is not nil.
func (q *Quake) Status() string {
	if q.err != nil {
		return "error"
	}

	switch {
	case q.Type == "not existing":
		return "deleted"
	case q.Type == "duplicate":
		return "duplicate"
	case q.EvaluationMode == "manual":
		return "reviewed"
	case q.EvaluationStatus == "confirmed":
		return "reviewed"
	default:
		return "automatic"
	}
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
