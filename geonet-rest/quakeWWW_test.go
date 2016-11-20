package main

import (
	"encoding/json"
	wt "github.com/GeoNet/weft/wefttest"
	"math"
	"testing"
	"time"
)

type quakeWWWFeatures struct {
	Features []quakeWWWFeature
}

type quakeWWWFeature struct {
	Properties quakeWWWProperties
	Geometry   geometry
}

type quakeWWWProperties struct {
	PublicID         string
	OriginTime       string
	Depth, Magnitude float64
	Status           string
	Intensity        string
	Agency           string
	UpdateTime       string
}

const (
	timeFmt = "2006-01-02 15:04:05"
)

func TestQuakesWWW(t *testing.T) {
	setup()
	defer teardown()

	b, err := wt.Request{Accept: JSON, URL: "/quakes/services/all.json"}.Do(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	var quakes quakeWWWFeatures

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
			if q.Properties.Intensity != "moderate" {
				t.Error("incorrect intensity")
			}
			if q.Properties.Agency != "WEL(GNS_Primary)" {
				t.Error("incorrect agency")
			}
			if q.Properties.Status != "reviewed" {
				t.Error("incorrect status")
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

			tm, _ := time.Parse(timeFmt, q.Properties.OriginTime)
			if math.Abs(float64(time.Now().UTC().Sub(tm)/time.Second)) > 10 {
				t.Error("time should be closer to now.  The time in the test data is reset when it is loaded. ")
			}
		}
	}

	if !found {
		t.Error("didn't find quake 2013p407387 in the list of Features.")
	}

	b, err = wt.Request{Accept: JSON, URL: "/quakes/services/felt.json"}.Do(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	quakes = quakeWWWFeatures{}

	err = json.Unmarshal(b, &quakes)
	if err != nil {
		t.Fatal(err)
	}

	if len(quakes.Features) != 2 {
		t.Error("found wrong number of quakes")
	}

	b, err = wt.Request{Accept: JSON, URL: "/quakes/services/quakes/newzealand/5/100.json"}.Do(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	quakes = quakeWWWFeatures{}

	err = json.Unmarshal(b, &quakes)
	if err != nil {
		t.Fatal(err)
	}

	if len(quakes.Features) != 1 {
		t.Error("found wrong number of quakes")
	}
}

func TestQuakeWWW(t *testing.T) {
	setup()
	defer teardown()

	b, err := wt.Request{Accept: JSON, URL: "/quake/services/quake/2013p407387.json"}.Do(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	var quakes quakeWWWFeatures

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
			if q.Properties.Intensity != "moderate" {
				t.Error("incorrect intensity")
			}
			if q.Properties.Agency != "WEL(GNS_Primary)" {
				t.Error("incorrect agency")
			}
			if q.Properties.Status != "reviewed" {
				t.Error("incorrect status")
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

			tm, _ := time.Parse(timeFmt, q.Properties.OriginTime)
			if math.Abs(float64(time.Now().UTC().Sub(tm)/time.Second)) > 10 {
				t.Error("time should be closer to now.  The time in the test data is reset when it is loaded. ")
			}
		}
	}

	if !found {
		t.Error("didn't find quake 2013p407387 in the list of Features.")
	}
}
