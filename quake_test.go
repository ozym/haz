package msg

import (
	"fmt"
	"math"
	"runtime"
	"strconv"
	"strings"
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

	l, err := q.Closest()
	if err != nil {
		t.Error("non nil error for closest locality.")
	}

	if l.Locality.Name != `Arthur's Pass` {
		t.Errorf("expected name Arthur's Pass got %s", l.Locality.Name)
	}

	delta(t, 25.89, l.Distance, 0.05)
	delta(t, 241.74, l.Bearing, 0.05)
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

	q.UsedPhaseCount = 22
	q.MagnitudeStationCount = 8

	if q.AlertQuality() {
		t.Error("should not alert for automatic with 22 phases and 8 mags.")
	}

	q.UsedPhaseCount = 22
	q.MagnitudeStationCount = 12

	if !q.AlertQuality() {
		t.Error("should be true for 22 phases and 12 mags.")
	}

	q.Type = "not existing"

	if q.AlertQuality() {
		t.Error("should not alert for deleted quake.")
	}

	q.Type = "duplicate"
	q.SetErr(nil)

	if q.AlertQuality() {
		t.Error("should not alert for duplicate quake.")
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

	q = Quake{}
	q.Time = time.Now().UTC().Add(time.Duration(-61) * time.Minute)
	q.Type = "earthquake"
	q.UsedPhaseCount = 28
	q.MagnitudeStationCount = 28

	if q.AlertQuality() {
		t.Error("should not alert for old high quality quake.")
	}
}

func TestPublish(t *testing.T) {
	q := Quake{}

	q.Site = "primary"
	eq(t, false, q.Publish())

	q.SetErr(nil)

	q.Site = "backup"
	eq(t, false, q.Publish())

	q.Type = ""
	q.EvaluationStatus = "automatic"
	q.SetErr(nil)

	q.Site = "primary"
	eq(t, false, q.Publish())

	q.SetErr(nil)

	q.Site = "backup"
	eq(t, false, q.Publish())

	q.EvaluationStatus = "confirmed"
	q.SetErr(nil)

	q.Site = "primary"
	eq(t, true, q.Publish())
	q.Site = "backup"
	eq(t, true, q.Publish())

	q.EvaluationStatus = "automatic"
	q.Depth = 0.01
	q.AzimuthalGap = 321.0
	q.MinimumDistance = 3.0
	q.Site = "primary"
	q.SetErr(nil)

	eq(t, false, q.Publish())

	q.SetErr(nil)
	q.Site = "backup"
	eq(t, false, q.Publish())

	q.EvaluationStatus = "automatic"
	q.Depth = 0.2
	q.AzimuthalGap = 319.0
	q.MinimumDistance = 2.4
	q.SetErr(nil)

	q.Site = "primary"
	eq(t, true, q.Publish())
	q.Site = "backup"
	eq(t, false, q.Publish())

	q.SetErr(nil)
	q.EvaluationStatus = "confirmed"

	q.Site = "primary"
	eq(t, true, q.Publish())
	q.Site = "backup"
	eq(t, true, q.Publish())

	q.SetErr(fmt.Errorf("errored quake"))

	q.Site = "primary"
	eq(t, false, q.Publish())
	q.Site = "backup"
	eq(t, false, q.Publish())
}

func TestAlertDuty(t *testing.T) {
	q := Quake{
		PublicID:              "2015p278423",
		Time:                  time.Now().UTC(),
		Latitude:              -37.92257397,
		Longitude:             178.3544071,
		Depth:                 9.62890625,
		EvaluationStatus:      "automatic",
		UsedPhaseCount:        25,
		AzimuthalGap:          180,
		MinimumDistance:       2.4,
		Magnitude:             5.0,
		MagnitudeStationCount: 12,
	}

	ab, am := q.AlertDuty()
	eq(t, true, ab)
	eq(t, true, strings.HasPrefix(am, "Eq Rpt: MAG 5.0, MM7, DEP 10, LOC 5 km south-east of Ruatoria,"))

	q.Longitude = 179.8 // quake off shore should still alert the Duty people.
	ab, am = q.AlertDuty()
	eq(t, true, ab)
	eq(t, true, strings.HasPrefix(am, "Eq Rpt: MAG 5.0, MM7, DEP 10, LOC 130 km east of Te Araroa,"))

	q.Magnitude = 3.0 // small quake off shore no alert.
	ab, _ = q.AlertDuty()
	eq(t, false, ab)

	q.Magnitude = 5.5
	q.SetErr(fmt.Errorf("errored quake no alert"))
	ab, _ = q.AlertDuty()
	eq(t, false, ab)
}

func TestAlertPIM(t *testing.T) {
	q := Quake{
		PublicID:              "2015p278423",
		Time:                  time.Now().UTC(),
		Latitude:              -37.92257397,
		Longitude:             178.3544071,
		Depth:                 9.62890625,
		EvaluationStatus:      "automatic",
		UsedPhaseCount:        25,
		AzimuthalGap:          180,
		MinimumDistance:       2.4,
		Magnitude:             6.0,
		MagnitudeStationCount: 12,
	}

	ab, am := q.AlertPIM()
	eq(t, true, ab)
	eq(t, true, strings.HasPrefix(am, "Eq Rpt: MAG 6.0, MM8, DEP 10, LOC 5 km south-east of Ruatoria,"))

	q.Magnitude = 3.0 // small quake off shore no alert.
	ab, _ = q.AlertPIM()
	eq(t, false, ab)
}

func TestAlertEqNews(t *testing.T) {
	q := Quake{
		PublicID:              "2015p278423",
		Time:                  time.Now().UTC(),
		Latitude:              -37.92257397,
		Longitude:             178.3544071,
		Depth:                 9.62890625,
		EvaluationStatus:      "automatic",
		UsedPhaseCount:        25,
		AzimuthalGap:          180,
		MinimumDistance:       2.4,
		Magnitude:             6.0,
		MagnitudeStationCount: 12,
	}

	ab, subject, body := q.AlertEqNews()
	eq(t, true, ab)

	eq(t, "NZ EQ: M6.0, severe intensity, 10km deep, 5 km south-east of Ruatoria", subject)

	// there is no sensible way to test the body (we have to change the quake time to now) so just eyeball it.
	fmt.Println(body)
}

func TestAlertTwitter(t *testing.T) {
	q := Quake{
		PublicID:              "2015p278423",
		Time:                  time.Now().UTC(),
		Latitude:              -37.92257397,
		Longitude:             178.3544071,
		Depth:                 9.62890625,
		EvaluationStatus:      "automatic",
		UsedPhaseCount:        25,
		AzimuthalGap:          180,
		MinimumDistance:       2.4,
		Magnitude:             6.0,
		MagnitudeStationCount: 12,
	}

	a, m := q.AlertTwitter(0)

	eq(t, true, a)
	eq(t, true, strings.HasPrefix(m, `Quake 5 km south-east of Ruatoria, intensity severe, approx. M6.0, depth 10 km http://geonet.org.nz/quakes/2015p278423`))

	a, m = q.AlertTwitter(7)
	eq(t, false, a)

}

func TestAlertUAPush(t *testing.T) {
	q := Quake{
		PublicID:              "2015p278423",
		Time:                  time.Now().UTC(),
		Latitude:              -37.92257397,
		Longitude:             178.3544071,
		Depth:                 9.62890625,
		EvaluationStatus:      "automatic",
		UsedPhaseCount:        25,
		AzimuthalGap:          180,
		MinimumDistance:       2.4,
		Magnitude:             6.0,
		MagnitudeStationCount: 12,
	}

	m, _ := q.AlertUAPush()

	eq(t, true, m == `Severe quake 5 km south-east of Ruatoria`)

}

func delta(t *testing.T, expected, actual, delta float64) {
	if math.Abs(expected-actual) > delta {
		t.Errorf("%s expected %f got %f diff = %f", loc(), expected, actual, math.Abs(expected-actual))
	}
}

func eq(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		t.Error("%s not equal", loc())
	}
}

func loc() string {
	_, _, l, _ := runtime.Caller(2)
	return "L" + strconv.Itoa(l)
}
