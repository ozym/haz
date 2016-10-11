package main

import (
	wt "github.com/GeoNet/weft/wefttest"
	"math"
	"testing"
	"github.com/GeoNet/haz"
	"github.com/golang/protobuf/proto"
)

func TestValProto(t *testing.T) {
	setup()
	defer teardown()

	b, err := wt.Request{Accept: protobuf, URL: "/volcano/val"}.Do(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	var vol haz.Volcanoes

	if err = proto.Unmarshal(b, &vol); err != nil {
		t.Error(err)
	}

	if len(vol.Volcanoes) != 12 {
		t.Error("found wrong number of val")
	}

	var found bool

	for _, v := range vol.Volcanoes {
		if v.VolcanoID == "ruapehu" {
			found = true

			if v.Title != "Ruapehu" {
				t.Error("incorrect title")
			}
			if v.Val.Activity != "Minor volcanic unrest." {
				t.Error("incorrect activity")
			}
			if v.Val.Hazards != "Volcanic unrest hazards." {
				t.Error("incorrect hazards")
			}
			if v.Val.Level != 1 {
				t.Error("incorrect level")
			}
			if math.Abs(v.Longitude -175.563) > tolerance {
				t.Error("incorrect Longitude")
			}
			if math.Abs(v.Latitude +39.281) > tolerance {
				t.Error("incorrect Latitude")
			}
		}
	}

	if !found {
		t.Error("didn't find Ruapehu")
	}
}
