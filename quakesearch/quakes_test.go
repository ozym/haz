package main

import (
	"bytes"
	"encoding/json"
	wt "github.com/GeoNet/weft/wefttest"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

const (
	HtmlContent = "text/html; charset=utf-8"
)

var (
	client *http.Client
)

type valid struct {
	Status string
}

func init() {
	timeout := time.Duration(5 * time.Second)
	client = &http.Client{
		Timeout: timeout,
	}
}

var routes = wt.Requests{
	//1. geojson routes
	{ID: wt.L(), URL: "/geojson?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713", Content: CONTENT_TYPE_GeoJSON, Accept: CONTENT_TYPE_GeoJSON},
	{ID: wt.L(), URL: "/geojson?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&startdate=2010-1-1T00:00:00&enddate=2015-1-1T00:00:00", Content: CONTENT_TYPE_GeoJSON, Accept: CONTENT_TYPE_GeoJSON},
	{ID: wt.L(), URL: "/geojson?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&minmag=3&maxmag=10", Content: CONTENT_TYPE_GeoJSON, Accept: CONTENT_TYPE_GeoJSON},
	{ID: wt.L(), URL: "/geojson?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&mindepth=10&maxdepth=200", Content: CONTENT_TYPE_GeoJSON, Accept: CONTENT_TYPE_GeoJSON},
	{ID: wt.L(), URL: "/geojson?limit=100&region=canterbury&minmag=3&maxmag=7&mindepth=1&maxdepth=200", Content: CONTENT_TYPE_GeoJSON, Accept: CONTENT_TYPE_GeoJSON},
	//some invalid requests
	{ID: wt.L(), URL: "/geojson?limit=100&bbox=163.60840,-49.18170,182.98828", Content: CONTENT_TYPE_GeoJSON, Accept: CONTENT_TYPE_GeoJSON, Status: http.StatusBadRequest},
	{ID: wt.L(), URL: "/geojson?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&startdate=2o10-1-1T00:00:00&enddate=2@15-1-1T00:00:00", Content: CONTENT_TYPE_GeoJSON, Accept: CONTENT_TYPE_GeoJSON, Status: http.StatusBadRequest},
	{ID: wt.L(), URL: "/geojson?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&minmag=3&maxmag=10a", Content: CONTENT_TYPE_GeoJSON, Accept: CONTENT_TYPE_GeoJSON, Status: http.StatusBadRequest},
	{ID: wt.L(), URL: "/geojson?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&mindepth=10&maxdepth=200a", Content: CONTENT_TYPE_GeoJSON, Accept: CONTENT_TYPE_GeoJSON, Status: http.StatusBadRequest},
	{ID: wt.L(), URL: "/geojson?limit=100&region=c@nterbury&minmag=3a&maxmag=7a&mindepth=1&maxdepth=200a", Content: CONTENT_TYPE_GeoJSON, Accept: CONTENT_TYPE_GeoJSON, Status: http.StatusBadRequest},

	//2. counts routes
	{ID: wt.L(), URL: "/count?bbox=163.60840,-49.18170,182.98828,-32.28713", Content: CONTENT_TYPE_JSON, Accept: CONTENT_TYPE_JSON},
	{ID: wt.L(), URL: "/count?bbox=163.60840,-49.18170,182.98828,-32.28713&startdate=2010-1-1T00:00:00&enddate=2015-1-1T00:00:00", Content: CONTENT_TYPE_JSON, Accept: CONTENT_TYPE_JSON},
	{ID: wt.L(), URL: "/count?bbox=163.60840,-49.18170,182.98828,-32.28713&minmag=3&maxmag=10", Content: CONTENT_TYPE_JSON, Accept: CONTENT_TYPE_JSON},
	{ID: wt.L(), URL: "/count?bbox=163.60840,-49.18170,182.98828,-32.28713&mindepth=10&maxdepth=200", Content: CONTENT_TYPE_JSON, Accept: CONTENT_TYPE_JSON},
	{ID: wt.L(), URL: "/count?region=canterbury&minmag=3&maxmag=7&mindepth=1&maxdepth=200", Content: CONTENT_TYPE_JSON, Accept: CONTENT_TYPE_JSON},

	//3. csv routes
	{ID: wt.L(), URL: "/csv?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713", Content: CONTENT_TYPE_CSV, Accept: CONTENT_TYPE_CSV},
	{ID: wt.L(), URL: "/csv?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&startdate=2010-1-1T00:00:00&enddate=2015-1-1T00:00:00", Content: CONTENT_TYPE_CSV, Accept: CONTENT_TYPE_CSV},
	{ID: wt.L(), URL: "/csv?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&minmag=3&maxmag=10", Content: CONTENT_TYPE_CSV, Accept: CONTENT_TYPE_CSV},
	{ID: wt.L(), URL: "/csv?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&mindepth=10&maxdepth=200", Content: CONTENT_TYPE_CSV, Accept: CONTENT_TYPE_CSV},
	{ID: wt.L(), URL: "/csv?limit=100&region=canterbury&minmag=3&maxmag=7&mindepth=1&maxdepth=200", Content: CONTENT_TYPE_CSV, Accept: CONTENT_TYPE_CSV},

	//4. gml routes
	{ID: wt.L(), URL: "/gml?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713", Content: CONTENT_TYPE_XML, Accept: CONTENT_TYPE_XML},
	{ID: wt.L(), URL: "/gml?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&startdate=2010-1-1T00:00:00&enddate=2015-1-1T00:00:00", Content: CONTENT_TYPE_XML, Accept: CONTENT_TYPE_XML},
	{ID: wt.L(), URL: "/gml?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&minmag=3&maxmag=10", Content: CONTENT_TYPE_XML, Accept: CONTENT_TYPE_XML},
	{ID: wt.L(), URL: "/gml?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&mindepth=10&maxdepth=200", Content: CONTENT_TYPE_XML, Accept: CONTENT_TYPE_XML},
	{ID: wt.L(), URL: "/gml?limit=100&region=canterbury&minmag=3&maxmag=7&mindepth=1&maxdepth=200", Content: CONTENT_TYPE_XML, Accept: CONTENT_TYPE_XML},

	//kml routes
	{ID: wt.L(), URL: "/kml?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713", Content: CONTENT_TYPE_KML, Accept: CONTENT_TYPE_KML},
	{ID: wt.L(), URL: "/kml?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&startdate=2010-1-1T00:00:00&enddate=2015-1-1T00:00:00", Content: CONTENT_TYPE_KML, Accept: CONTENT_TYPE_KML},
	{ID: wt.L(), URL: "/kml?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&minmag=3&maxmag=10", Content: CONTENT_TYPE_KML, Accept: CONTENT_TYPE_KML},
	{ID: wt.L(), URL: "/kml?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&mindepth=10&maxdepth=200", Content: CONTENT_TYPE_KML, Accept: CONTENT_TYPE_KML},
	{ID: wt.L(), URL: "/kml?limit=100&region=canterbury&minmag=3&maxmag=7&mindepth=1&maxdepth=200", Content: CONTENT_TYPE_KML, Accept: CONTENT_TYPE_KML},

	// soh routes
	{ID: wt.L(), URL: "/soh"},
	{ID: wt.L(), URL: "/soh/up"},
}

func TestRoutes(t *testing.T) {
	setup(t)
	defer teardown()

	for _, r := range routes {
		if b, err := r.Do(ts.URL); err != nil {
			t.Error(err)
			t.Error(string(b))
		}
	}

	if err := routes.DoCheckQuery(ts.URL); err != nil {
		t.Error(err)
	}
}

func TestQuakesCount(t *testing.T) {
	setup(t)
	defer teardown()
	//1. get all quakes
	r := wt.Request{
		Accept: CONTENT_TYPE_JSON,
		URL:    "/count?bbox=163.60840,-49.18170,182.98828,-32.28713",
	}
	b, err := r.Do(ts.URL)
	if err != nil {
		t.Error(err)
		t.Error(string(b))
	}
	var qc QuakesCount
	err = json.Unmarshal(b, &qc)
	if err != nil {
		log.Fatal(err)
	}
	if qc.Count != 3 {
		t.Errorf("Found wrong number of quakes: %d", qc.Count)
	}

	//2. get only one quake
	r = wt.Request{
		Accept: CONTENT_TYPE_JSON,
		URL:    "/count?bbox=163.60840,-49.18170,182.98828,-32.28713&startdate=2010-1-1T00:00:00&enddate=2015-1-1T00:00:00",
	}
	b, err = r.Do(ts.URL)
	if err != nil {
		t.Error(err)
		t.Error(string(b))
	}

	err = json.Unmarshal(b, &qc)
	if err != nil {
		t.Error(err)
	}
	if qc.Count != 1 {
		t.Errorf("Found wrong number of quakes: %d", qc.Count)
	}

	//3. get 2 quakes
	r = wt.Request{
		Accept: CONTENT_TYPE_JSON,
		URL:    "/count?bbox=163.60840,-49.18170,182.98828,-32.28713&minmag=5",
	}
	b, err = r.Do(ts.URL)
	if err != nil {
		t.Error(err)
		t.Error(string(b))
	}

	err = json.Unmarshal(b, &qc)
	if err != nil {
		t.Error(err)
	}
	if qc.Count != 2 {
		t.Errorf("Found wrong number of quakes: %d", qc.Count)
	}
}

func TestQuakesGeoJson(t *testing.T) {
	setup(t)
	defer teardown()
	//1. get all quakes
	r := wt.Request{URL: "/geojson?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713", Accept: CONTENT_TYPE_GeoJSON}

	b, err := r.Do(ts.URL)
	if err != nil {
		t.Error(err)
		t.Error(string(b))
	}

	var f GeoJsonFeatureCollection
	err = json.Unmarshal(b, &f)
	if err != nil {
		log.Fatal(err)
	}
	if len(f.Features) != 3 {
		t.Errorf("Found wrong number of features: %d", len(f.Features))
	}

	//2. get only one quake
	r = wt.Request{URL: "/geojson?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&startdate=2010-1-1T00:00:00&enddate=2015-1-1T00:00:00", Accept: CONTENT_TYPE_GeoJSON}

	b, err = r.Do(ts.URL)
	if err != nil {
		t.Error(err)
		t.Error(string(b))
	}

	err = json.Unmarshal(b, &f)
	if err != nil {
		log.Fatal(err)
	}
	if len(f.Features) != 1 {
		t.Errorf("Found wrong number of features: %d", len(f.Features))
	}
	if f.Features[0].Properties.Publicid != "3366146" {
		t.Errorf("Found wrong publicid: %s", f.Features[0].Properties.Publicid)
	}

}

func TestGeoJSONFormat(t *testing.T) {
	setup(t)
	defer teardown()

	r := wt.Request{
		ID:      wt.L(),
		Accept:  CONTENT_TYPE_GeoJSON,
		Content: CONTENT_TYPE_GeoJSON,
		URL:     "/geojson?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&startdate=2000-1-1T00:00:00&enddate=2015-1-1T00:00:00",
	}

	b, err := r.Do(ts.URL)
	if err != nil {
		t.Error(err)
		t.Error(string(b))
	}
	body := bytes.NewBuffer(b)

	res, err := client.Post("http://geojsonlint.com/validate", "application/vnd.geo+json", body)
	defer res.Body.Close()
	if err != nil {
		t.Errorf("Problem contacting geojsonlint for test %s", r.ID)
	}

	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Problem reading body from geojsonlint for test %s", r.ID)
	}

	var v valid

	err = json.Unmarshal(b, &v)
	if err != nil {
		t.Errorf("Problem unmarshalling body from geojsonlint for test %s", r.ID)
	}

	if v.Status != "ok" {
		t.Errorf("invalid geoJSON for test %s" + r.ID)
	}
}

type QuakesCount struct {
	Count int      `json:"count"`
	Dates []string `json:"dates"`
}
