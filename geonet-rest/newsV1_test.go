package main

import (
	"encoding/json"
	wt "github.com/GeoNet/weft/wefttest"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestUnmarshalNews(t *testing.T) {
	xmlFile, err := os.Open("etc/test/files/geonet-news.xml")
	if err != nil {
		t.Error("problem opening etc/test/files/geonet-news.xml")
	}
	defer xmlFile.Close()

	b, _ := ioutil.ReadAll(xmlFile)

	f, err := unmarshalNews(b)
	if err != nil {
		t.Error("Problem unmarshaling XML")
	}

	if f.Entries[0].Title != "Strong quake wakes southern North Island" {
		t.Error("wrong title")
	}

	if f.Entries[0].Link.Href != "http://info.geonet.org.nz/display/quake/2014/09/23/Strong+quake+wakes+southern+North+Island" {
		t.Error("wrong link")
	}

	if f.Entries[0].Href != "http://info.geonet.org.nz/display/quake/2014/09/23/Strong+quake+wakes+southern+North+Island" {
		t.Error("wrong link")
	}

	if f.Entries[0].MHref != "http://info.geonet.org.nz/m/view-rendered-page.action?abstractPageId=11567177" {
		t.Error("wrong mobile link")
	}

	if f.Entries[0].Published != "2014-09-22T16:11:48Z" {
		t.Error("wrong published")
	}

	if f.Entries[0].Id != "tag:info.geonet.org.nz,2009:blogpost-11567177-5" {
		t.Error("wrong id")
	}
}

func TestNewsV1(t *testing.T) {
	setup()
	defer teardown()

	b, err := wt.Request{Accept: V1JSON, URL: "/news/geonet"}.Do(ts.URL)
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
