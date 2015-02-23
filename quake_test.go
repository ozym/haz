package msg

import (
	"fmt"
	"math"
	"runtime"
	"strconv"
	"testing"
	"time"
)

func TestMMI(t *testing.T) {
	// Max MMI
	q := Quake{}

	// Christchurch 2011
	q.Depth = 5.0
	q.Magnitude = 6.3
	delta(t, 8.86, q.MMI(), 0.005)

	// Gisbourne 2007
	q.Depth = 40.0
	q.Magnitude = 6.8
	delta(t, 8.19, q.MMI(), 0.005)

	// Darfield 2010
	q.Depth = 11.0
	q.Magnitude = 7.1
	delta(t, 9.96, q.MMI(), 0.005)

	// small deep event
	q.Depth = 150.0
	q.Magnitude = 1.5
	delta(t, -1.0, q.MMI(), 0.005)

	// large deep event
	q.Depth = 150
	q.Magnitude = 6.5
	delta(t, 6.23, q.MMI(), 0.005)

	// moderate shallow event
	q.Depth = 7.0
	q.Magnitude = 4.4
	delta(t, 6.41, q.MMI(), 0.005)

	// errored quake.
	q.SetErr(fmt.Errorf("errored quake"))
	delta(t, -1, q.MMI(), 0.005)

	// MMI @ distance
	q.SetErr(nil)
	q.Depth = 27.4
	q.Magnitude = 3.9
	delta(t, 2.65, MMIDistance(110, 27.4, q.MMI()), 0.1)

	q.Depth = 22.2
	q.Magnitude = 4.2
	delta(t, 5.27, MMIDistance(5, 22.2, q.MMI()), 0.1)
	delta(t, 5.27, MMIDistance(0, 22.2, q.MMI()), 0.1)

}

func TestClosest(t *testing.T) {
	q := Quake{}
	q.Longitude = 171.29
	q.Latitude = -43.06

	l, d, b, err := q.Closest()
	if err != nil {
		t.Error("non nil error for closest locality.")
	}

	if l.Name != `Arthur's Pass` {
		t.Errorf("expected name Arthur's Pass got %s", l.Name)
	}

	delta(t, 25.89, d, 0.05)
	delta(t, 241.74, b, 0.05)
}

func TestAlertQuality(t *testing.T) {
	q := Quake{}
	q.Time = time.Now().UTC()
	q.Type = "earthquake"
	q.EvaluationMode = "automatic"
	q.UsedPhaseCount = 8
	q.MagnitudeStationCount = 8

	if q.Status() != "automatic" {
		t.Error("quake should be automatic.")
	}

	if q.AlertQuality() {
		t.Error("should not alert for automatic with 8 phases and mags.")
	}

	if q.Err() == nil {
		t.Errorf("err  should be set for low phases and mags.")
	}

	q.SetErr(nil)

	q.UsedPhaseCount = 22
	q.MagnitudeStationCount = 8

	if q.AlertQuality() {
		t.Error("should not alert for automatic with 22 phases and 8 mags.")
	}

	if q.Err() == nil {
		t.Errorf("err  should be set for 22 phases and 8 mags.")
	}

	q.SetErr(nil)

	q.UsedPhaseCount = 22
	q.MagnitudeStationCount = 12

	if !q.AlertQuality() {
		t.Error("should be true for 22 phases and 12 mags.")
	}

	q.Type = "not existing"

	if q.AlertQuality() {
		t.Error("should not alert for deleted quake.")
	}

	if q.Err() == nil {
		t.Errorf("err  should be set for deleted quake.")
	}

	q.Type = "duplicate"
	q.SetErr(nil)

	if q.AlertQuality() {
		t.Error("should not alert for duplicate quake.")
	}

	if q.Err() == nil {
		t.Errorf("err  should be set for duplicate quake.")
	}

	q = Quake{}
	q.Time = time.Now().UTC()
	q.Type = "earthquake"
	q.EvaluationMode = "manual"
	q.UsedPhaseCount = 8
	q.MagnitudeStationCount = 8

	if !q.AlertQuality() {
		t.Error("should be true for manual review")
	}

	q.EvaluationMode = ""
	q.EvaluationStatus = "confirmed"

	if !q.AlertQuality() {
		t.Error("should be true for confirmed quake")
	}

	q.UsedPhaseCount = 28
	q.MagnitudeStationCount = 28

	if !q.AlertQuality() {
		t.Error("should be true for confirmed quake high quality")
	}

	q = Quake{}
	q.Time = time.Now().UTC().Add(time.Duration(-61) * time.Minute)
	q.Type = "earthquake"
	q.EvaluationMode = "manual"
	q.UsedPhaseCount = 8
	q.MagnitudeStationCount = 8

	if q.AlertQuality() {
		t.Error("should not alert for old reviewed quake.")
	}

	if q.Err() == nil {
		t.Errorf("err should be set for old reviewed quake.")
	}

	q = Quake{}
	q.Time = time.Now().UTC().Add(time.Duration(-61) * time.Minute)
	q.Type = "earthquake"
	q.UsedPhaseCount = 28
	q.MagnitudeStationCount = 28

	if q.AlertQuality() {
		t.Error("should not alert for old high quality quake.")
	}

	if q.Err() == nil {
		t.Errorf("err should be set for old high quality quake.")
	}

}

func delta(t *testing.T, expected, actual, delta float64) {
	if math.Abs(expected-actual) > delta {
		t.Errorf("%s expected %f got %f diff = %f", loc(), expected, actual, math.Abs(expected-actual))
	}
}

func loc() string {
	_, _, l, _ := runtime.Caller(2)
	return "L" + strconv.Itoa(l)
}
