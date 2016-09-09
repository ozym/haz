package main

import (
	"encoding/json"
	wt "github.com/GeoNet/weft/wefttest"
	"math"
	"testing"
)

type valV2Features struct {
	Features []valV2Feature
}

type valV2Feature struct {
	Properties valV2Properties
	Geometry   geometry
}

type valV2Properties struct {
	VolcanoID, VolcanoTitle, Activity, Hazards string
	Level                                      int
}

func TestValV2(t *testing.T) {
	setup()
	defer teardown()

	b, err := wt.Request{Accept: V2GeoJSON, URL: "/volcano/val"}.Do(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	var val valV2Features

	err = json.Unmarshal(b, &val)
	if err != nil {
		t.Fatal(err)
	}

	if len(val.Features) != 12 {
		t.Error("found wrong number of val")
	}

	var found bool

	for _, q := range val.Features {
		if q.Properties.VolcanoID == "ruapehu" {
			found = true

			if q.Properties.Activity != "Minor volcanic unrest." {
				t.Error("incorrect activity")
			}
			if q.Properties.Hazards != "Volcanic unrest hazards." {
				t.Error("incorrect hazards")
			}
			if q.Properties.Level != 1 {
				t.Error("incorrect level")
			}
			if math.Abs(q.Geometry.Longitude()-175.563) > tolerance {
				t.Error("incorrect Longitude")
			}
			if math.Abs(q.Geometry.Latitude()+39.281) > tolerance {
				t.Error("incorrect Latitude")
			}
		}
	}

	if !found {
		t.Error("didn't find Ruapehu")
	}
}
