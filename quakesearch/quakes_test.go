package main

// TODO would be great to update this to use wefttest for testing the routes.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strconv"
	"testing"
	"time"
)

// Test helper functions and constants have been decoupled from github.com/GeoNet/web/webtest.
// These constants are for error and other pages.  They can be changed.
const (
	HtmlContent = "text/html; charset=utf-8"
)

var (
	client *http.Client
)

type Content struct {
	Accept string
	URI    string
}

type Route struct {
	Accept, Content string
	Response        int
	routes          []route
}

type route struct {
	id, uri string
}

type valid struct {
	Status string
}

func init() {
	timeout := time.Duration(5 * time.Second)
	client = &http.Client{
		Timeout: timeout,
	}
}

func TestQuakesCount(t *testing.T) {
	setup(t)
	defer teardown()
	//1. get all quakes
	c := Content{
		Accept: CONTENT_TYPE_JSON,
		URI:    "/count?bbox=163.60840,-49.18170,182.98828,-32.28713",
	}
	b, err := c.getContent(ts)
	if err != nil {
		t.Fatal(err)
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
	c = Content{
		Accept: CONTENT_TYPE_JSON,
		URI:    "/count?bbox=163.60840,-49.18170,182.98828,-32.28713&startdate=2010-1-1T00:00:00&enddate=2015-1-1T00:00:00",
	}
	b, err = c.getContent(ts)
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(b, &qc)
	if err != nil {
		log.Fatal(err)
	}
	if qc.Count != 1 {
		t.Errorf("Found wrong number of quakes: %d", qc.Count)
	}

	//3. get 2 quakes
	c = Content{
		Accept: CONTENT_TYPE_JSON,
		URI:    "/count?bbox=163.60840,-49.18170,182.98828,-32.28713&minmag=5",
	}
	b, err = c.getContent(ts)
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(b, &qc)
	if err != nil {
		log.Fatal(err)
	}
	if qc.Count != 2 {
		t.Errorf("Found wrong number of quakes: %d", qc.Count)
	}
}

func TestQuakesGeoJson(t *testing.T) {
	setup(t)
	defer teardown()
	//1. get all quakes
	c := Content{
		Accept: CONTENT_TYPE_GeoJSON,
		URI:    "/geojson?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713",
	}
	b, err := c.getContent(ts)
	if err != nil {
		t.Fatal(err)
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
	c = Content{
		Accept: CONTENT_TYPE_GeoJSON,
		URI:    "/geojson?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&startdate=2010-1-1T00:00:00&enddate=2015-1-1T00:00:00",
	}
	b, err = c.getContent(ts)
	if err != nil {
		t.Fatal(err)
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

func TestRoutes(t *testing.T) {
	setup(t)
	defer teardown()

	//1 GeoJSON routes
	r := Route{
		Accept:   CONTENT_TYPE_GeoJSON,
		Content:  CONTENT_TYPE_GeoJSON,
		Response: http.StatusOK,
	}
	r.addRoute("/geojson?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713")
	r.addRoute("/geojson?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&startdate=2010-1-1T00:00:00&enddate=2015-1-1T00:00:00")
	r.addRoute("/geojson?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&minmag=3&maxmag=10")
	r.addRoute("/geojson?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&mindepth=10&maxdepth=200")
	r.addRoute("/geojson?limit=100&region=canterbury&minmag=3&maxmag=7&mindepth=1&maxdepth=200")
	r.testRoutes("GeoJSON routes", ts, t)

	//2. Count
	r = Route{
		Accept:   CONTENT_TYPE_JSON,
		Content:  CONTENT_TYPE_JSON,
		Response: http.StatusOK,
	}
	r.addRoute("/count?bbox=163.60840,-49.18170,182.98828,-32.28713")
	r.addRoute("/count?bbox=163.60840,-49.18170,182.98828,-32.28713&startdate=2010-1-1T00:00:00&enddate=2015-1-1T00:00:00")
	r.addRoute("/count?bbox=163.60840,-49.18170,182.98828,-32.28713&minmag=3&maxmag=10")
	r.addRoute("/count?bbox=163.60840,-49.18170,182.98828,-32.28713&mindepth=10&maxdepth=200")
	r.addRoute("/count?region=canterbury&minmag=3&maxmag=7&mindepth=1&maxdepth=200")
	r.testRoutes("V1JSON routes", ts, t)

	//3 CSV routes
	r = Route{
		Accept:   CONTENT_TYPE_CSV,
		Content:  CONTENT_TYPE_CSV,
		Response: http.StatusOK,
	}
	r.addRoute("/csv?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713")
	r.addRoute("/csv?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&startdate=2010-1-1T00:00:00&enddate=2015-1-1T00:00:00")
	r.addRoute("/csv?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&minmag=3&maxmag=10")
	r.addRoute("/csv?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&mindepth=10&maxdepth=200")
	r.addRoute("/csv?limit=100&region=canterbury&minmag=3&maxmag=7&mindepth=1&maxdepth=200")
	r.testRoutes("V1CSV", ts, t)

	//4 GML routes
	r = Route{
		Accept:   CONTENT_TYPE_XML,
		Content:  CONTENT_TYPE_XML,
		Response: http.StatusOK,
	}
	r.addRoute("/gml?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713")
	r.addRoute("/gml?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&startdate=2010-1-1T00:00:00&enddate=2015-1-1T00:00:00")
	r.addRoute("/gml?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&minmag=3&maxmag=10")
	r.addRoute("/gml?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&mindepth=10&maxdepth=200")
	r.addRoute("/gml?limit=100&region=canterbury&minmag=3&maxmag=7&mindepth=1&maxdepth=200")
	r.testRoutes("GML routes", ts, t)

	//5 KML routes
	r = Route{
		Accept:   CONTENT_TYPE_KML,
		Content:  CONTENT_TYPE_KML,
		Response: http.StatusOK,
	}
	r.addRoute("/kml?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713")
	r.addRoute("/kml?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&startdate=2010-1-1T00:00:00&enddate=2015-1-1T00:00:00")
	r.addRoute("/kml?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&minmag=3&maxmag=10")
	r.addRoute("/kml?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&mindepth=10&maxdepth=200")
	r.addRoute("/kml?limit=100&region=canterbury&minmag=3&maxmag=7&mindepth=1&maxdepth=200")
	r.testRoutes("KML routes", ts, t)
}

func TestGeoJSON(t *testing.T) {
	setup(t)
	defer teardown()

	// GeoJSON routes
	r := Route{
		Accept:   CONTENT_TYPE_GeoJSON,
		Content:  CONTENT_TYPE_GeoJSON,
		Response: http.StatusOK,
	}
	r.addRoute("/geojson?limit=100&bbox=163.60840,-49.18170,182.98828,-32.28713&startdate=2000-1-1T00:00:00&enddate=2015-1-1T00:00:00")

	r.testGeoJSONRoutes(ts, t)
}

type QuakesCount struct {
	Count int      `json:"count"`
	Dates []string `json:"dates"`
}

// Get returns the Content from the test server.
// If the tests are not being run verbose (go test -v) then silences the application logging.
func (c *Content) getContent(s *httptest.Server) (b []byte, err error) {
	if !testing.Verbose() {
		log.SetOutput(ioutil.Discard)
	}

	req, _ := http.NewRequest("GET", s.URL+c.URI, nil)
	req.Header.Add("Accept", c.Accept)
	res, _ := client.Do(req)
	defer res.Body.Close()

	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	if res.StatusCode != 200 {
		err = fmt.Errorf("Non 200 error code: %d", res.StatusCode)
		return
	}

	if res.Header.Get("Content-Type") != c.Accept {
		err = fmt.Errorf("incorrect Content-Type: %s", res.Header.Get("Content-Type"))
		return
	}

	return
}

// Add a URI to be tested for the Route.
// The line that this function is called from will be included in test failure messages.
func (r *Route) addRoute(uri string) {
	r.routes = append(r.routes, route{loc(), uri})
}

func loc() (loc string) {
	_, _, l, _ := runtime.Caller(2)
	return "L" + strconv.Itoa(l)
}

func (rt *Route) testRoutes(m string, s *httptest.Server, t *testing.T) {
	if !testing.Verbose() {
		log.SetOutput(ioutil.Discard)
	}

	for _, r := range rt.routes {
		req, err := http.NewRequest("GET", s.URL+r.uri, nil)
		if err != nil {
			t.Error(err)
			continue
		}

		req.Header.Add("Accept", rt.Accept)
		res, err := client.Do(req)
		if err != nil {
			t.Error(err)
			continue
		}
		defer res.Body.Close()

		if res.StatusCode != rt.Response {
			t.Errorf("%s: wrong response code for test %s: got %d expected %d", m, r.id, res.StatusCode, rt.Response)
		}

		// Allow for error pages with type web.HtmlContent
		if res.Header.Get("Content-Type") != rt.Content {
			if res.Header.Get("Content-Type") != HtmlContent {
				t.Errorf("%s: incorrect Content-Type for test %s: got %s expected %s", m, r.id, res.Header.Get("Content-Type"), rt.Content)
			}
		}
	}
}

// GeoJSON test the content returned from the Route is valid GeoJSON using http://geojsonlint.com/
func (rt *Route) testGeoJSONRoutes(s *httptest.Server, t *testing.T) {
	if !testing.Verbose() {
		log.SetOutput(ioutil.Discard)
	}

	for _, r := range rt.routes {
		req, _ := http.NewRequest("GET", s.URL+r.uri, nil)
		req.Header.Add("Accept", rt.Accept)
		res, _ := client.Do(req)

		if res.StatusCode != rt.Response {
			t.Errorf("Wrong response code for test %s: got %d expected %d", r.id, res.StatusCode, rt.Response)
		}

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Problem reading body for test %s", r.id)
		}

		body := bytes.NewBuffer(b)

		res, err = client.Post("http://geojsonlint.com/validate", "application/vnd.geo+json", body)
		defer res.Body.Close()
		if err != nil {
			t.Errorf("Problem contacting geojsonlint for test %s", r.id)
		}

		b, err = ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Problem reading body from geojsonlint for test %s", r.id)
		}

		var v valid

		err = json.Unmarshal(b, &v)
		if err != nil {
			t.Errorf("Problem unmarshalling body from geojsonlint for test %s", r.id)
		}

		if v.Status != "ok" {
			t.Errorf("invalid geoJSON for test %s" + r.id)
		}
	}
}
