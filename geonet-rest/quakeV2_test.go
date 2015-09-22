package main

import (
	"encoding/json"
	"github.com/GeoNet/web"
	"github.com/GeoNet/web/webtest"
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

	c := webtest.Content{
		Accept: web.V2GeoJSON,
		URI:    "/quake?MMI=3",
	}

	b, err := c.Get(ts)
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

	c = webtest.Content{
		Accept: web.V2GeoJSON,
		URI:    "/quake?MMI=6",
	}

	b, err = c.Get(ts)
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

	c := webtest.Content{
		Accept: web.V2GeoJSON,
		URI:    "/quake/2013p407387",
	}

	b, err := c.Get(ts)
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
