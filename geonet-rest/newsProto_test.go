package main

import (
	"github.com/GeoNet/haz"
	wt "github.com/GeoNet/weft/wefttest"
	"github.com/golang/protobuf/proto"
	"testing"
)

func TestNewsProto(t *testing.T) {
	setup()
	defer teardown()

	b, err := wt.Request{Accept: protobuf, URL: "/news/geonet"}.Do(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	var n haz.News

	err = proto.Unmarshal(b, &n)
	if err != nil {
		t.Fatal(err)
	}

	if len(n.Stories) == 0 {
		t.Error("empty news feed.")
	}

	if n.Stories[0].Title == "" {
		t.Error("empty title for news feed.")
	}

	if n.Stories[0].Link == "" {
		t.Error("empty link for news feed.")
	}

	if n.Stories[0].Published.Sec == 0 {
		t.Error("incorrect time.")
	}
}
