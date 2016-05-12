package main

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

/**
* Test the wfs web service working
 */
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

// test number of quakes
func TestQuakesGeoJson(t *testing.T) {
	setup(t)
	defer teardown()
	//1. get all quakes
	c := Content{
		Accept: CONTENT_TYPE_GeoJSON,
		URI:    "/ows?service=WFS&version=1.0.0&request=GetFeature&typeName=geonet:quake_search_v1&outputFormat=json&maxFeatures=100&&cql_filter=BBOX(origin_geom,163.60840,-49.18170,182.98828,-32.28713)",
	}
	b, err := c.getContent(ts)
	if err != nil {
		t.Fatal(err)
	}
	var f GeoJsonFeatureCollection
	err = json.Unmarshal(b, &f)
	if err != nil {
		t.Fatal(err)
	}
	if len(f.Features) != 3 {
		t.Errorf("Found wrong number of features: %d", len(f.Features))
	}

	//2. get only one quake
	c = Content{
		Accept: CONTENT_TYPE_GeoJSON,
		URI:    "/ows?service=WFS&version=1.0.0&request=GetFeature&typeName=geonet:quake_search_v1&outputFormat=json&maxFeatures=100&cql_filter=BBOX(origin_geom,163.60840,-49.18170,182.98828,-32.28713)+AND+origintime>='2010-01-01'+AND+origintime<'2015-01-01'",
	}
	b, err = c.getContent(ts)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(b, &f)
	if err != nil {
		t.Fatal(err)
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
	r.addRoute("/ows?service=WFS&version=1.0.0&request=GetFeature&typeName=geonet:quake_search_v1&outputFormat=json&maxFeatures=100&cql_filter=BBOX(origin_geom,163.60840,-49.18170,182.98828,-32.28713)")
	r.addRoute("/ows?service=WFS&version=1.0.0&request=GetFeature&typeName=geonet:quake_search_v1&outputFormat=json&maxFeatures=100&cql_filter=BBOX(origin_geom,163.60840,-49.18170,182.98828,-32.28713)+AND+origintime>='2010-01-01'+AND+origintime<'2015-01-01'")
	r.addRoute("/ows?service=WFS&version=1.0.0&request=GetFeature&typeName=geonet:quake_search_v1&outputFormat=json&maxFeatures=100&cql_filter=BBOX(origin_geom,163.60840,-49.18170,182.98828,-32.28713)+AND+magnitude>=3+AND+magnitude<10")
	r.addRoute("/ows?service=WFS&version=1.0.0&request=GetFeature&typeName=geonet:quake_search_v1&outputFormat=json&maxFeatures=100&cql_filter=DWITHIN(origin_geom,Point+(174.201+-40.589),500,meters)+AND+depth>=10+AND+depth<200")
	r.testRoutes("GeoJSON routes", ts, t)

	//2 CSV routes
	r = Route{
		Accept:   CONTENT_TYPE_CSV,
		Content:  CONTENT_TYPE_CSV,
		Response: http.StatusOK,
	}
	r.addRoute("/ows?service=WFS&version=1.0.0&request=GetFeature&typeName=geonet:quake_search_v1&outputFormat=csv&maxFeatures=100&cql_filter=BBOX(origin_geom,163.60840,-49.18170,182.98828,-32.28713)")
	r.addRoute("/ows?service=WFS&version=1.0.0&request=GetFeature&typeName=geonet:quake_search_v1&outputFormat=csv&maxFeatures=100&cql_filter=BBOX(origin_geom,163.60840,-49.18170,182.98828,-32.28713)+AND+origintime>='2010-01-01'+AND+origintime<'2015-01-01'")
	r.addRoute("/ows?service=WFS&version=1.0.0&request=GetFeature&typeName=geonet:quake_search_v1&outputFormat=csv&maxFeatures=100&cql_filter=WITHIN(origin_geom,POLYGON((172.951+-41.767,+172.001+-42.832,+169.564+-44.341,+172.312+-45.412,+175.748+-42.908,+172.951+-41.767)))+AND+magnitude>=3+AND+magnitude<10")
	r.addRoute("/ows?service=WFS&version=1.0.0&request=GetFeature&typeName=geonet:quake_search_v1&outputFormat=csv&maxFeatures=100&cql_filter=DWITHIN(origin_geom,Point+(174.201+-40.589),500,meters)+AND+depth>=10+AND+depth<200")
	r.testRoutes("V1CSV", ts, t)

	//3 GML2 routes
	r = Route{
		Accept:   CONTENT_TYPE_XML,
		Content:  CONTENT_TYPE_XML,
		Response: http.StatusOK,
	}
	r.addRoute("/ows?service=WFS&version=1.0.0&request=GetFeature&typeName=geonet:quake_search_v1&outputFormat=GML2&maxFeatures=100&cql_filter=BBOX(origin_geom,163.60840,-49.18170,182.98828,-32.28713)")
	r.addRoute("/ows?service=WFS&version=1.0.0&request=GetFeature&typeName=geonet:quake_search_v1&outputFormat=GML2&maxFeatures=100&cql_filter=BBOX(origin_geom,163.60840,-49.18170,182.98828,-32.28713)+AND+origintime>='2010-01-01'+AND+origintime<'2015-01-01'")
	r.addRoute("/ows?service=WFS&version=1.0.0&request=GetFeature&typeName=geonet:quake_search_v1&outputFormat=GML2&maxFeatures=100&cql_filter=WITHIN(origin_geom,POLYGON((172.951+-41.767,+172.001+-42.832,+169.564+-44.341,+172.312+-45.412,+175.748+-42.908,+172.951+-41.767)))+AND+magnitude>=3+AND+magnitude<10")
	r.addRoute("/ows?service=WFS&version=1.0.0&request=GetFeature&typeName=geonet:quake_search_v1&outputFormat=GML2&maxFeatures=100&cql_filter=DWITHIN(origin_geom,Point+(174.201+-40.589),500,meters)+AND+depth>=10+AND+depth<200")
	r.testRoutes("GML2 routes", ts, t)

	//4. GML3 text/xml;subtype=gml/3.2
	r = Route{
		Accept:   CONTENT_TYPE_XML,
		Content:  CONTENT_TYPE_XML,
		Response: http.StatusOK,
	}
	r.addRoute("/ows?service=WFS&version=1.0.0&request=GetFeature&typeName=geonet:quake_search_v1&outputFormat=text/xml;subtype=gml/3.2&maxFeatures=100&cql_filter=BBOX(origin_geom,163.60840,-49.18170,182.98828,-32.28713)")
	r.addRoute("/ows?service=WFS&version=1.0.0&request=GetFeature&typeName=geonet:quake_search_v1&outputFormat=text/xml;subtype=gml/3.2&maxFeatures=100&cql_filter=BBOX(origin_geom,163.60840,-49.18170,182.98828,-32.28713)+AND+origintime>='2010-01-01'+AND+origintime<'2015-01-01'")
	r.addRoute("/ows?service=WFS&version=1.0.0&request=GetFeature&typeName=geonet:quake_search_v1&outputFormat=text/xml;subtype=gml/3.2&maxFeatures=100&cql_filter=WITHIN(origin_geom,POLYGON((172.951+-41.767,+172.001+-42.832,+169.564+-44.341,+172.312+-45.412,+175.748+-42.908,+172.951+-41.767)))+AND+magnitude>=3+AND+magnitude<10")
	r.addRoute("/ows?service=WFS&version=1.0.0&request=GetFeature&typeName=geonet:quake_search_v1&outputFormat=text/xml;subtype=gml/3.2&maxFeatures=100&cql_filter=DWITHIN(origin_geom,Point+(174.201+-40.589),500,meters)+AND+depth>=10+AND+depth<200")
	r.testRoutes("GML3 routes", ts, t)

	//5 KML routes
	r = Route{
		Accept:   CONTENT_TYPE_KML,
		Content:  CONTENT_TYPE_KML,
		Response: http.StatusOK,
	}
	r.addRoute("/wms/kml?layers=geonet:quake_search_v1&maxFeatures=100&cql_filter=BBOX(origin_geom,163.60840,-49.18170,182.98828,-32.28713)")
	r.addRoute("/wms/kml?layers=geonet:quake_search_v1&maxFeatures=100&cql_filter=BBOX(origin_geom,163.60840,-49.18170,182.98828,-32.28713)+AND+origintime>='2010-01-01'+AND+origintime<'2015-01-01'")
	r.addRoute("/wms/kml?layers=geonet:quake_search_v1&maxFeatures=100&cql_filter=WITHIN(origin_geom,POLYGON((172.951+-41.767,+172.001+-42.832,+169.564+-44.341,+172.312+-45.412,+175.748+-42.908,+172.951+-41.767)))+AND+magnitude>=3+AND+magnitude<10")
	r.addRoute("/wms/kml?layers=geonet:quake_search_v1&maxFeatures=100&cql_filter=DWITHIN(origin_geom,Point+(174.201+-40.589),500,meters)+AND+depth>=10+AND+depth<200")
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
	r.addRoute("/ows?service=WFS&version=1.0.0&request=GetFeature&typeName=geonet:quake_search_v1&outputFormat=json&maxFeatures=100&cql_filter=BBOX(origin_geom,163.60840,-49.18170,182.98828,-32.28713)+AND+origintime>='2000-01-01'+AND+origintime<'2015-01-01'")

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
