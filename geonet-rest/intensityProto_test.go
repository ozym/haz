package main

import (
	wt "github.com/GeoNet/weft/wefttest"
	"math"
	"testing"
	"github.com/GeoNet/haz"
	"github.com/golang/protobuf/proto"
)

func TestIntensityMeasuredProto(t *testing.T) {
	setup()
	defer teardown()

	b, err := wt.Request{Accept: protobuf, URL: "/intensity?type=measured"}.Do(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	var i haz.Shaking

	if err = proto.Unmarshal(b, &i); err != nil {
		t.Fatal(err)
	}

	if len(i.Mmi) != 1 {
		t.Error("found wrong number of intensities.")
	}
	if math.Abs(i.Mmi[0].Longitude - 175.49) > tolerance {
		t.Error("incorrect Longitude")
	}
	if math.Abs(i.Mmi[0].Latitude + 40.2) > tolerance {
		t.Error("incorrect Latitude")
	}
	if i.Mmi[0].Mmi != 4 {
		t.Errorf("incorrect MMI expected 4 got %d", i.Mmi[0].Mmi)
	}
}

func TestIntensityReportedLatestProto(t *testing.T) {
	setup()
	defer teardown()

	b, err := wt.Request{Accept: protobuf, URL: "/intensity?type=reported"}.Do(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	var i haz.Shaking

	if err = proto.Unmarshal(b, &i); err != nil {
		t.Fatal(err)
	}

	// all of the mmi are at the same geohash so no need for us to worry about order.
	if len(i.Mmi) != 1 {
		t.Error("found wrong number of intensities.")
	}
	if math.Abs(i.Mmi[0].Longitude - 176.489868) > tolerance {
		t.Error("incorrect Longitude")
	}
	if math.Abs(i.Mmi[0].Latitude + 40.201721) > tolerance {
		t.Error("incorrect Latitude")
	}
	if i.Mmi[0].Mmi != 6 {
		t.Error("incorrect MMI")
	}
	if i.Mmi[0].Count != 3 {
		t.Error("incorrect count")
	}
	if i.MmiTotal != 3 {
		t.Error("incorrect total count")
	}
	if len(i.MmiSummary) != 3 {
		t.Errorf("wrong mmi count length expected 3 got %d", len(i.MmiSummary))
		t.Errorf("%+v", i.MmiSummary)
	}

	for _, v := range []int32{4, 5, 6} {
		count, ok := i.MmiSummary[v]
		if !ok {
			t.Errorf("missing count for %d", v)
		}
		if count != 1 {
			t.Errorf("count for %d should be 1", v)
		}
	}
}

func TestIntensityReportedProto(t *testing.T) {
	setup()
	defer teardown()

	b, err := wt.Request{Accept: protobuf, URL: "/intensity?type=reported&publicID=2013p407387"}.Do(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	var i haz.Shaking

	if err = proto.Unmarshal(b, &i); err != nil {
		t.Fatal(err)
	}

	if len(i.Mmi) != 1 {
		t.Error("found wrong number of intensities.")
	}
	if math.Abs(i.Mmi[0].Longitude - 176.489868) > tolerance {
		t.Error("incorrect Longitude")
	}
	if math.Abs(i.Mmi[0].Latitude + 40.201721) > tolerance {
		t.Error("incorrect Latitude")
	}
	if i.Mmi[0].Mmi != 6 {
		t.Error("incorrect MMI")
	}
	if i.Mmi[0].Count != 3 {
		t.Error("incorrect count")
	}
	if i.MmiTotal != 3 {
		t.Error("incorrect total count")
	}
	if len(i.MmiSummary) != 3 {
		t.Error("wrong mmi count length")
	}

	for _, v := range []int32{4, 5, 6} {
		count, ok := i.MmiSummary[v]
		if !ok {
			t.Errorf("missing count for %d", v)
		}
		if count != 1 {
			t.Errorf("count for %d should be 1", v)
		}
	}
}
