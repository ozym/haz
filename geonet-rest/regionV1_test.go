//# Region Information
//
//##/region
//
// Look up region information.  All calls return [GeoJSON](http://geojson.org/) with Polygon features.
//
//### Properties
//
// * regionID - a unique indentifier for the region.
// * title - the region title.
// * group - the region group.
//
package main

import (
	"encoding/json"
	"github.com/GeoNet/web"
	"github.com/GeoNet/web/webtest"
	"log"
	"testing"
)

type RegionFeatures struct {
	Features []RegionFeature
}

type RegionFeature struct {
	Properties RegionProperties
	Geometry   RegionGeometry
}

type RegionGeometry struct {
	Type string
}

type RegionProperties struct {
	RegionID, Title, Group string
}

//## Quake Regions
//
// **GET /region?type=quake**
//
// Get all quake regions.
//
//### Example request:
//
// `/region?type=quake`
//
func TestRegionsV1(t *testing.T) {
	setup()
	defer teardown()

	c := webtest.Content{
		Accept: web.V1GeoJSON,
		URI:    "/region?type=quake",
	}

	b, err := c.Get(ts)
	if err != nil {
		t.Fatal(err)
	}

	var f RegionFeatures

	err = json.Unmarshal(b, &f)
	if err != nil {
		log.Fatal(err)
	}

	if !(len(f.Features) >= 1) {
		t.Error("Found wrong number of features")
	}

	for _, feat := range f.Features {
		var g = feat.Properties.Group
		if !(g == "region" || g == "north" || g == "south") {
			t.Error("Found non quake region")
		}
	}
}

//## Single Region
//
// **GET /region/(regionID)**
//
// Get a single region.
//
//### Example request:
//
// `/region/wellington`
//
func TestRegionV1(t *testing.T) {
	setup()
	defer teardown()

	c := webtest.Content{
		Accept: web.V1GeoJSON,
		URI:    "/region/wellington",
	}

	b, err := c.Get(ts)
	if err != nil {
		t.Fatal(err)
	}

	var f RegionFeatures

	err = json.Unmarshal(b, &f)
	if err != nil {
		log.Fatal(err)
	}

	if !(len(f.Features) == 1) {
		t.Error("Found wrong number of features")
	}

	if f.Features[0].Properties.RegionID != "wellington" {
		t.Errorf("wrong region: %s", f.Features[0].Properties.RegionID)
	}
}
