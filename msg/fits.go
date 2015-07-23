package msg

import (
	"encoding/json"
	"log"
	"time"
)

/*
Observation is a useful wire format for Fits data.
*/
type Observation struct {
	NetworkID string
	SiteID    string
	TypeID    string
	MethodID  string
	DateTime  time.Time
	Value     float64
	Error     float64 // Error for the observation value.  0.0 for no or unknown error.
	err       error
}

func (o *Observation) SetErr(err error) {
	o.err = err
}

func (o *Observation) Err() error {
	return o.err
}

func (o *Observation) Decode(b []byte) {
	o.err = json.Unmarshal(b, o)
}

func (o *Observation) Encode() ([]byte, error) {
	if o.err != nil {
		return nil, o.err
	}

	return json.Marshal(o)
}

func (o *Observation) RxLog() {
	if o.err != nil {
		return
	}

	log.Printf("Received observation %s.%s %s %s", o.NetworkID, o.SiteID, o.TypeID, o.MethodID)
}
