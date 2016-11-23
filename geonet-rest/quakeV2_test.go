package main

import (
	"encoding/json"
	wt "github.com/GeoNet/weft/wefttest"
	"math"
	"testing"
	"time"
)

type quakeV2Features struct {
	Features []quakeV2Feature
}

type quakeV2Feature struct {
	Properties quakeV2Properties
	Geometry   geometry
}

type geometry struct {
	Coordinates []float64
}

type quakeV2Properties struct {
	PublicID         string
	Time             time.Time
	ModificationTime time.Time // in quake history only.
	Depth, Magnitude float64
	Locality         string
	MMI              int
	Quality          string
}

func (g *geometry) Longitude() float64 {
	return g.Coordinates[0]
}

func (g *geometry) Latitude() float64 {
	return g.Coordinates[1]
}

func TestQuakesV2(t *testing.T) {
	setup()
	defer teardown()

	b, err := wt.Request{Accept: V2GeoJSON, URL: "/quake?MMI=3"}.Do(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	var quakes quakeV2Features

	err = json.Unmarshal(b, &quakes)
	if err != nil {
		t.Fatal(err)
	}

	if len(quakes.Features) != 2 {
		t.Error("found wrong number of quakes")
	}

	var found bool

	for _, q := range quakes.Features {
		if q.Properties.PublicID == "2013p407387" {
			found = true
			if q.Properties.Locality != "25 km south-east of Amberley" {
				t.Error("incorrect locality")
			}
			if q.Properties.Quality != "best" {
				t.Error("incorrect quality")
			}
			if q.Properties.MMI != 4 {
				t.Error("incorrect MMI")
			}
			if math.Abs(q.Properties.Depth-20.334389) > tolerance {
				t.Error("incorrect depth")
			}
			if math.Abs(q.Properties.Magnitude-4.027879) > tolerance {
				t.Error("incorrect magnitude")
			}
			if math.Abs(q.Geometry.Longitude()-172.94479) > tolerance {
				t.Error("incorrect Longitude")
			}
			if math.Abs(q.Geometry.Latitude()+43.359699) > tolerance {
				t.Error("incorrect Latitude")
			}
			if math.Abs(float64(time.Now().UTC().Sub(q.Properties.Time)/time.Second)) > 10 {
				t.Error("time should be closer to now.  The time in the test data is reset when it is loaded. ")
			}
		}
	}

	if !found {
		t.Error("didn't find quake 2013p407387 in the list of Features.")
	}

	b, err = wt.Request{Accept: V2GeoJSON, URL: "/quake?MMI=6"}.Do(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	quakes = quakeV2Features{}

	err = json.Unmarshal(b, &quakes)
	if err != nil {
		t.Fatal(err)
	}

	if len(quakes.Features) != 1 {
		t.Error("found wrong number of quakes")
	}
}

func TestQuakeV2(t *testing.T) {
	setup()
	defer teardown()

	b, err := wt.Request{Accept: V2GeoJSON, URL: "/quake/2013p407387"}.Do(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	var quakes quakeV2Features

	err = json.Unmarshal(b, &quakes)
	if err != nil {
		t.Fatal(err)
	}

	if len(quakes.Features) != 1 {
		t.Error("found more than 1 quake")
	}

	var found bool

	for _, q := range quakes.Features {
		if q.Properties.PublicID == "2013p407387" {
			found = true
			if q.Properties.Locality != "25 km south-east of Amberley" {
				t.Error("incorrect locality")
			}
			if q.Properties.Quality != "best" {
				t.Error("incorrect quality")
			}
			if q.Properties.MMI != 4 {
				t.Error("incorrect MMI")
			}
			if math.Abs(q.Properties.Depth-20.334389) > tolerance {
				t.Error("incorrect depth")
			}
			if math.Abs(q.Properties.Magnitude-4.027879) > tolerance {
				t.Error("incorrect magnitude")
			}
			if math.Abs(q.Geometry.Longitude()-172.94479) > tolerance {
				t.Error("incorrect Longitude")
			}
			if math.Abs(q.Geometry.Latitude()+43.359699) > tolerance {
				t.Error("incorrect Latitude")
			}
			if math.Abs(float64(time.Now().UTC().Sub(q.Properties.Time)/time.Second)) > 10 {
				t.Error("time should be closer to now.  The time in the test data is reset when it is loaded. ")
			}
		}
	}

	if !found {
		t.Error("didn't find quake 2013p407387 in the list of Features.")
	}
}

func TestQuakeHistoryV2(t *testing.T) {
	setup()
	defer teardown()

	b, err := wt.Request{Accept: V2GeoJSON, URL: "/quake/history/2013p407387"}.Do(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	var quakes quakeV2Features

	err = json.Unmarshal(b, &quakes)
	if err != nil {
		t.Fatal(err)
	}

	if len(quakes.Features) != 1 {
		t.Error("found more than 1 quake")
	}

	var found bool

	for _, q := range quakes.Features {
		if q.Properties.PublicID == "2013p407387" {
			found = true
			if q.Properties.Locality != "25 km south-east of Amberley" {
				t.Error("incorrect locality")
			}
			if q.Properties.Quality != "best" {
				t.Error("incorrect quality")
			}
			if q.Properties.MMI != 4 {
				t.Error("incorrect MMI")
			}
			if math.Abs(q.Properties.Depth-20.334389) > tolerance {
				t.Error("incorrect depth")
			}
			if math.Abs(q.Properties.Magnitude-4.027879) > tolerance {
				t.Error("incorrect magnitude")
			}
			if math.Abs(q.Geometry.Longitude()-172.94479) > tolerance {
				t.Error("incorrect Longitude")
			}
			if math.Abs(q.Geometry.Latitude()+43.359699) > tolerance {
				t.Error("incorrect Latitude")
			}
			if math.Abs(float64(time.Now().UTC().Sub(q.Properties.Time)/time.Second)) > 10 {
				t.Error("time should be closer to now.  The time in the test data is reset when it is loaded. ")
			}
			mt, err := time.Parse(time.RFC3339, "2013-05-31T21:37:41.549Z")
			if err != nil {
				t.Fatal(err)
			}
			if !q.Properties.ModificationTime.Equal(mt) {
				t.Error("incorrect modification time")
			}
		}
	}

	if !found {
		t.Error("didn't find quake 2013p407387 in the list of Features.")
	}
}

type magCountV2 struct {
	MagnitudeCount struct {
		Days7   map[string]int
		Days28  map[string]int
		Days365 map[string]int
	}
	Rate struct {
		PerDay map[string]int
	}
}

func TestQuakeStatsV2(t *testing.T) {
	setup()
	defer teardown()

	b, err := wt.Request{Accept: V2JSON, URL: "/quake/stats"}.Do(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	var m magCountV2

	err = json.Unmarshal(b, &m)
	if err != nil {
		t.Fatal(err)
	}

	if m.MagnitudeCount.Days28["4"] != 1 {
		t.Errorf("exptected 1 got %d", m.MagnitudeCount.Days28["4"])
	}

	if m.MagnitudeCount.Days28["5"] != 1 {
		t.Errorf("exptected 1 got %d", m.MagnitudeCount.Days28["5"])
	}

	if m.MagnitudeCount.Days7["4"] != 1 {
		t.Errorf("exptected 1 got %d", m.MagnitudeCount.Days7["4"])
	}

	if m.MagnitudeCount.Days7["5"] != 1 {
		t.Errorf("exptected 1 got %d", m.MagnitudeCount.Days7["5"])
	}

	if m.MagnitudeCount.Days365["4"] != 1 {
		t.Errorf("exptected 1 got %d", m.MagnitudeCount.Days365["4"])
	}

	if m.MagnitudeCount.Days365["5"] != 1 {
		t.Errorf("exptected 1 got %d", m.MagnitudeCount.Days365["5"])
	}

	if m.MagnitudeCount.Days28["6"] != 0 {
		t.Errorf("exptected 0 got %d", m.MagnitudeCount.Days28["6"])
	}

	if len(m.Rate.PerDay) != 1 {
		t.Errorf("expected 1 for daily rate, got %d", len(m.Rate.PerDay))
	}

	for k := range m.Rate.PerDay {
		if m.Rate.PerDay[k] != 2 {
			t.Errorf("expected 2 quakes for day %s", k)
		}
	}
}
