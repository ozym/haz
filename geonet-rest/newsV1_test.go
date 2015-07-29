//# News
//
//## /news
//
package main

import (
	"io/ioutil"
	"os"
	"testing"
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

//## GeoNet News
//
// **/news/geonet**
//
// Returns a simple JSON version of the GeoNet News RSS feed.
//
//### Example request:
//
// `/news/geonet`
//
func TestGeoNetNewsV1(t *testing.T) {
	// tested in routes.  This is a handle for the docs.
}
