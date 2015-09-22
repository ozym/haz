package main

import (
	"encoding/json"
	"github.com/GeoNet/web"
	"github.com/GeoNet/web/webtest"
	"math"
	"testing"
	"time"
)

type quakeV1Features struct {
	Features []quakeV1Feature
}

type quakeV1Feature struct {
	Properties quakeV1Properties
	Geometry   geometry
}

type quakeV1Properties struct {
	Depth           float64
	Intensity       string
	Locality        string
	Magnitude       float64
	PublicID        string
	Quality         string
	RegionIntensity string
	Time            time.Time
}

func TestQuakesV1(t *testing.T) {
	setup()
	defer teardown()

	c := webtest.Content{
		Accept: web.V1GeoJSON,
		URI:    "/quake?regionID=newzealand&regionIntensity=weak&number=100&quality=best,caution,good,deleted",
	}

	b, err := c.Get(ts)
	if err != nil {
		t.Fatal(err)
	}

	var quakes quakeV1Features

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
			if q.Properties.RegionIntensity != "light" {
				t.Error("wrong region intensity")
			}
			if q.Properties.Intensity != "moderate" {
				t.Error("wrong intensity")
			}
			if q.Properties.Locality != "25 km south-east of Amberley" {
				t.Error("incorrect locality")
			}
			if q.Properties.Quality != "best" {
				t.Error("incorrect quality")
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
		Accept: web.V1GeoJSON,
		URI:    "/quake?regionID=newzealand&regionIntensity=moderate&number=100&quality=best,caution,good,deleted",
	}

	b, err = c.Get(ts)
	if err != nil {
		t.Fatal(err)
	}

	quakes = quakeV1Features{}

	err = json.Unmarshal(b, &quakes)
	if err != nil {
		t.Fatal(err)
	}

	if len(quakes.Features) != 1 {
		t.Error("found wrong number of quakes")
	}
}

func TestQuakeV1(t *testing.T) {
	setup()
	defer teardown()

	c := webtest.Content{
		Accept: web.V1GeoJSON,
		URI:    "/quake/2013p407387",
	}

	b, err := c.Get(ts)
	if err != nil {
		t.Fatal(err)
	}

	var quakes quakeV1Features

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
			if q.Properties.RegionIntensity != "light" {
				t.Error("wrong region intensity")
			}
			if q.Properties.Intensity != "moderate" {
				t.Error("wrong intensity")
			}
			if q.Properties.Locality != "25 km south-east of Amberley" {
				t.Error("incorrect locality")
			}
			if q.Properties.Quality != "best" {
				t.Error("incorrect quality")
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
