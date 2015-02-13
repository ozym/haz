package msg

import (
	"reflect"
	"testing"
	"time"
)

func TestDecodeSC3ML07(t *testing.T) {
	ev := Quake{}

	ev.PublicID = "2012p070732"
	ev.Type = ""
	ev.AgencyID = "WEL(GNS_Primary)"

	mt, err := time.Parse(time.RFC3339Nano, "2014-01-09T11:01:05.67583Z")
	if err != nil {
		t.Fatal(err)
	}

	ev.ModificationTime = mt

	otm, err := time.Parse(time.RFC3339Nano, "2012-01-27T04:06:25.369465Z")
	if err != nil {
		t.Fatal(err)
	}

	ev.Time = otm

	ev.Latitude = -43.15704211
	ev.Longitude = 170.9096047
	ev.Depth = 5.234375
	ev.MethodID = "NonLinLoc"
	ev.EarthModelID = "nz3drx"
	ev.EvaluationMode = "automatic"
	ev.EvaluationStatus = ""
	ev.UsedPhaseCount = 8
	ev.UsedStationCount = 8
	ev.StandardError = 0.5944258933
	ev.AzimuthalGap = 98.73193969
	ev.MinimumDistance = 0.1515531055
	ev.Magnitude = 2.652616042
	ev.MagnitudeUncertainty = 0
	ev.MagnitudeType = "M"
	ev.MagnitudeStationCount = 7

	es := ReadSC3ML07("etc/2012p070732-sc3.xml")

	if es.Err() != nil {
		t.Fatalf("es.Err non nil: %s", es.Err().Error())
	}

	if !reflect.DeepEqual(ev, es) {
		t.Error("events ev and es not equal")
	}

	if ev.Status() != "automatic" {
		t.Errorf("incorrect status for ev expected automatic got %s", ev.Status())
	}
}

func TestUnmarshalBadSC3ML07(t *testing.T) {
	es := ReadSC3ML07("etc/2012p070732-missing-sc3.xml")

	if es.Err() == nil {
		t.Fatal("es.Err not set")
	}
}

func TestUnmarshalEmptySC3ML07(t *testing.T) {
	es := ReadSC3ML07("etc/999.xml")

	if es.Err() == nil {
		t.Fatal("es.Err not set")
	}
}
