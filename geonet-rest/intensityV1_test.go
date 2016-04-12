package main

import (
	"encoding/json"
	"github.com/GeoNet/web/webtest"
	"math"
	"testing"
)

// Measured intensity.

type intensityMeasuredV1Features struct {
	Features []intensityMeasuredV1Feature
}

type intensityMeasuredV1Feature struct {
	Properties intensityMeasuredV1Properties
	Geometry   geometry
}

type intensityMeasuredV1Properties struct {
	MMI int
}

func TestIntensityMeasuredV1(t *testing.T) {
	setup()
	defer teardown()

	c := webtest.Content{
		Accept: V1GeoJSON,
		URI:    "/intensity?type=measured",
	}

	b, err := c.Get(ts)
	if err != nil {
		t.Fatal(err)
	}

	var i intensityMeasuredV1Features

	err = json.Unmarshal(b, &i)
	if err != nil {
		t.Fatal(err)
	}

	if len(i.Features) != 1 {
		t.Error("found wrong number of intensities.")
	}
	if math.Abs(i.Features[0].Geometry.Longitude()-175.49) > tolerance {
		t.Error("incorrect Longitude")
	}
	if math.Abs(i.Features[0].Geometry.Latitude()+40.2) > tolerance {
		t.Error("incorrect Latitude")
	}
	if i.Features[0].Properties.MMI != 4 {
		t.Error("incorrect MMI")
	}
}
