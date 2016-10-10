package main

import (
	wt "github.com/GeoNet/weft/wefttest"
	"math"
	"testing"
	"github.com/GeoNet/haz"
	"github.com/golang/protobuf/proto"
	"time"
)

func TestQuakeProto(t *testing.T) {
	setup()
	defer teardown()

	b, err := wt.Request{Accept: protobuf, URL: "/quake/2013p407387"}.Do(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	var q haz.Quake

	if err = proto.Unmarshal(b, &q); err != nil {
		t.Fatal(err)
	}

	if q.Locality != "25 km south-east of Amberley" {
		t.Error("incorrect locality")
	}
	if q.PublicID != "2013p407387" {
		t.Error("incorrect publicID")
	}
	if q.Quality != "best" {
		t.Error("incorrect quality")
	}
	if q.Mmi != 4 {
		t.Error("incorrect MMI")
	}
	if math.Abs(q.Depth - 20.334389) > tolerance {
		t.Error("incorrect depth")
	}
	if math.Abs(q.Magnitude - 4.027879) > tolerance {
		t.Error("incorrect magnitude")
	}
	if math.Abs(q.Longitude - 172.94479) > tolerance {
		t.Error("incorrect Longitude")
	}
	if math.Abs(q.Latitude + 43.359699) > tolerance {
		t.Error("incorrect Latitude")
	}

	if q.Time == nil {
		t.Fatal("nil time")
	}

	tm := time.Unix(q.Time.Sec, q.Time.Nsec)

	if math.Abs(float64(time.Now().UTC().Sub(tm) / time.Second)) > 10 {
		t.Error("time should be closer to now.  The time in the test data is reset when it is loaded. ")
	}
}

func TestQuakesProto(t *testing.T) {
	setup()
	defer teardown()

	b, err := wt.Request{Accept: protobuf, URL: "/quake?MMI=3"}.Do(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	var quakes haz.Quakes

	if err = proto.Unmarshal(b, &quakes); err != nil {
		t.Fatal(err)
	}

	var found bool

	for _, q := range quakes.Quakes {

		if q.PublicID == "2013p407387" {

			found = true

			if q.Locality != "25 km south-east of Amberley" {
				t.Error("incorrect locality")
			}
			if q.PublicID != "2013p407387" {
				t.Error("incorrect publicID")
			}
			if q.Quality != "best" {
				t.Error("incorrect quality")
			}
			if q.Mmi != 4 {
				t.Error("incorrect MMI")
			}
			if math.Abs(q.Depth - 20.334389) > tolerance {
				t.Error("incorrect depth")
			}
			if math.Abs(q.Magnitude - 4.027879) > tolerance {
				t.Error("incorrect magnitude")
			}
			if math.Abs(q.Longitude - 172.94479) > tolerance {
				t.Error("incorrect Longitude")
			}
			if math.Abs(q.Latitude + 43.359699) > tolerance {
				t.Error("incorrect Latitude")
			}

			if q.Time == nil {
				t.Fatal("nil time")
			}

			tm := time.Unix(q.Time.Sec, q.Time.Nsec)

			if math.Abs(float64(time.Now().UTC().Sub(tm) / time.Second)) > 10 {
				t.Error("time should be closer to now.  The time in the test data is reset when it is loaded. ")
			}
		}
	}

	if !found {
		t.Error("didn't find quake")
	}
}

func TestQuakeHistoryProto(t *testing.T) {
	setup()
	defer teardown()

	b, err := wt.Request{Accept: protobuf, URL: "/quake/history/2013p407387"}.Do(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	var quakes haz.Quakes

	if err = proto.Unmarshal(b, &quakes); err != nil {
		t.Fatal(err)
	}

	if len(quakes.Quakes) != 1 {
		t.Error("expected 1 quake")
	}

	var found bool

	for _, q := range quakes.Quakes {

		if q.PublicID == "2013p407387" {

			found = true

			if q.Locality != "25 km south-east of Amberley" {
				t.Error("incorrect locality")
			}
			if q.PublicID != "2013p407387" {
				t.Error("incorrect publicID")
			}
			if q.Quality != "best" {
				t.Error("incorrect quality")
			}
			if q.Mmi != 4 {
				t.Error("incorrect MMI")
			}
			if math.Abs(q.Depth - 20.334389) > tolerance {
				t.Error("incorrect depth")
			}
			if math.Abs(q.Magnitude - 4.027879) > tolerance {
				t.Error("incorrect magnitude")
			}
			if math.Abs(q.Longitude - 172.94479) > tolerance {
				t.Error("incorrect Longitude")
			}
			if math.Abs(q.Latitude + 43.359699) > tolerance {
				t.Error("incorrect Latitude")
			}

			if q.Time == nil {
				t.Fatal("nil time")
			}

			tm := time.Unix(q.Time.Sec, q.Time.Nsec)

			if math.Abs(float64(time.Now().UTC().Sub(tm) / time.Second)) > 10 {
				t.Error("time should be closer to now.  The time in the test data is reset when it is loaded. ")
			}
		}
	}

	if !found {
		t.Error("didn't find quake")
	}
}