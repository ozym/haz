package main

import (
	"encoding/json"
	"github.com/GeoNet/web/webtest"
	"math"
	"testing"
)

// Measured intensity.

type intensityMeasuredV2Features struct {
	Features []intensityMeasuredV2Feature
}

type intensityMeasuredV2Feature struct {
	Properties intensityMeasuredV2Properties
	Geometry   geometry
}

type intensityMeasuredV2Properties struct {
	MMI int
}

func TestIntensityMeasuredV2(t *testing.T) {
	setup()
	defer teardown()

	c := webtest.Content{
		Accept: V2GeoJSON,
		URI:    "/intensity?type=measured",
	}

	b, err := c.Get(ts)
	if err != nil {
		t.Fatal(err)
	}

	var i intensityMeasuredV2Features

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

// Reported intensity.

type intensityReportedV2Features struct {
	Features []intensityReportedV2Feature
	Count    int
	MMICount map[string]int `json:"count_mmi"`
}

type intensityReportedV2Feature struct {
	Properties intensityReportedV2Properties
	Geometry   geometry
}

type intensityReportedV2Properties struct {
	MMI   int
	Count int
}

func TestIntensityReportedLatestV2(t *testing.T) {
	setup()
	defer teardown()

	c := webtest.Content{
		Accept: V2GeoJSON,
		URI:    "/intensity?type=reported",
	}

	b, err := c.Get(ts)
	if err != nil {
		t.Fatal(err)
	}

	var i intensityReportedV2Features

	err = json.Unmarshal(b, &i)
	if err != nil {
		t.Fatal(err)
	}

	if len(i.Features) != 1 {
		t.Error("found wrong number of intensities.")
	}
	if math.Abs(i.Features[0].Geometry.Longitude()-176.489868) > tolerance {
		t.Error("incorrect Longitude")
	}
	if math.Abs(i.Features[0].Geometry.Latitude()+40.201721) > tolerance {
		t.Error("incorrect Latitude")
	}
	if i.Features[0].Properties.MMI != 6 {
		t.Error("incorrect MMI")
	}
	if i.Features[0].Properties.Count != 3 {
		t.Error("incorrect count")
	}
	if i.Count != 3 {
		t.Error("incorrect total count")
	}
	if len(i.MMICount) != 3 {
		t.Error("wrong mmi count length")
	}

	for _, v := range []string{"4", "5", "6"} {
		count, ok := i.MMICount[v]
		if !ok {
			t.Errorf("missing count for %s", v)
		}
		if count != 1 {
			t.Errorf("count for %s should be 1", v)
		}
	}
}

func TestIntensityReportedV2(t *testing.T) {
	setup()
	defer teardown()

	c := webtest.Content{
		Accept: V2GeoJSON,
		URI:    "/intensity?type=reported&publicID=2013p407387",
	}

	b, err := c.Get(ts)
	if err != nil {
		t.Fatal(err)
	}

	var i intensityReportedV2Features

	err = json.Unmarshal(b, &i)
	if err != nil {
		t.Fatal(err)
	}

	if len(i.Features) != 1 {
		t.Error("found wrong number of intensities.")
	}
	if math.Abs(i.Features[0].Geometry.Longitude()-176.489868) > tolerance {
		t.Error("incorrect Longitude")
	}
	if math.Abs(i.Features[0].Geometry.Latitude()+40.201721) > tolerance {
		t.Error("incorrect Latitude")
	}
	if i.Features[0].Properties.MMI != 6 {
		t.Error("incorrect MMI")
	}
	if i.Features[0].Properties.Count != 3 {
		t.Error("incorrect count")
	}
	if i.Count != 3 {
		t.Error("incorrect total count")
	}
	if len(i.MMICount) != 3 {
		t.Error("wrong mmi count length")
	}

	for _, v := range []string{"4", "5", "6"} {
		count, ok := i.MMICount[v]
		if !ok {
			t.Errorf("missing count for %s", v)
		}
		if count != 1 {
			t.Errorf("count for %s should be 1", v)
		}
	}
}
