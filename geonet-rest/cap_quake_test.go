package main

import (
	"encoding/xml"
	"github.com/GeoNet/web/webtest"
	"testing"
)

type atomFeed struct {
	Entries []atomEntry `xml:"entry"`
}

type atomEntry struct {
	Links []atomLink `xml:"link"`
}

type atomLink struct {
	LinkType string `xml:"type,attr"`
	Href     string `xml:"href,attr"`
}

func TestCapQuakeFeed(t *testing.T) {
	setup()
	defer teardown()

	c := webtest.Content{
		Accept: "application/xml",
		URI:    "/cap/1.2/GPA1.0/feed/atom1.0/quake",
	}

	b, err := c.Get(ts)
	if err != nil {
		t.Fatal(err)
	}

	var feed atomFeed

	err = xml.Unmarshal(b, &feed)

	if len(feed.Entries) != 1 {
		t.Error("should find 1 quake")
	}

	if len(feed.Entries[0].Links) != 2 {
		t.Error("should find 2 links")
	}

	var found bool

	for _, l := range feed.Entries[0].Links {
		if l.LinkType == "application/cap+xml" {
			found = true
			if l.Href != "https://localhost/cap/1.2/GPA1.0/quake/2015p012816.1420493554884741" {
				t.Error("Didn't find correct CAP link")
			}
		}
	}

	if !found {
		t.Error("didn't find application/cap+xml in links.")
	}
}

type capAlert struct {
	Identifier string  `xml:"identifier"`
	Info       capInfo `xml:"info"`
}

type capInfo struct {
	Severity string `xml:"severity"`
}

func TestCapQuake(t *testing.T) {
	setup()
	defer teardown()

	c := webtest.Content{
		Accept: "application/cap+xml",
		URI:    "/cap/1.2/GPA1.0/quake/2015p012816.1420493554884741",
	}

	b, err := c.Get(ts)
	if err != nil {
		t.Fatal(err)
	}

	var a capAlert

	err = xml.Unmarshal(b, &a)

	// The CAP response can be validated with the Google Public Alerts CAP v1.0 profile at
	//  https://cap-validator.appspot.com/
	if a.Identifier != "2015p012816.1420493554884741" {
		t.Error("incorrect identifier")
	}
	if a.Info.Severity != "Moderate" {
		t.Error("incorrect severity")
	}
}
