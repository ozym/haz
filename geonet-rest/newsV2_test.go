package main

import (
	"encoding/json"
	"github.com/GeoNet/web/webtest"
	"testing"
	"time"
)

type feedJSON struct {
	Feed []newsJSON `json:"feed"`
}

type newsJSON struct {
	Title, Link, Mlink string
	Published          time.Time
}

func TestNewsV2(t *testing.T) {
	setup()
	defer teardown()

	c := webtest.Content{
		Accept: V2JSON,
		URI:    "/news/geonet",
	}

	b, err := c.Get(ts)
	if err != nil {
		t.Fatal(err)
	}

	var n feedJSON

	err = json.Unmarshal(b, &n)
	if err != nil {
		t.Fatal(err)
	}

	if len(n.Feed) == 0 {
		t.Error("empty news feed.")
	}

	if n.Feed[0].Title == "" {
		t.Error("empty title for news feed.")
	}

	if n.Feed[0].Link == "" {
		t.Error("empty link for news feed.")
	}

	if n.Feed[0].Mlink == "" {
		t.Error("empty mlink for news feed.")
	}

	var tm time.Time
	if n.Feed[0].Published == tm {
		t.Error("incorrect time.")
	}
}
