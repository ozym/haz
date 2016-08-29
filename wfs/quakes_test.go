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

/**
* Test the wfs web service working
 */
const (
	HtmlContent = "text/html; charset=utf-8"
)

var (
	client *http.Client

	url_wfs = "/geonet/ows?service=WFS&version=1.0.0&request=GetFeature&typeName=geonet:quake_search_v1&maxFeatures=100"
	url_kml = "/geonet/wms/kml?layers=geonet:quake_search_v1&maxFeatures=100"

	cql1 = "cql_filter=BBOX(origin_geom,163.60840,-49.18170,182.98828,-32.28713)"
	cql2 = "cql_filter=BBOX(origin_geom,163.60840,-49.18170,182.98828,-32.28713)+AND+origintime>='2010-01-01'+AND+origintime<'2015-01-01'"
	cql3 = "cql_filter=BBOX(origin_geom,163.60840,-49.18170,182.98828,-32.28713)+AND+magnitude>=3+AND+magnitude<10"
	cql4 = "cql_filter=DWITHIN(origin_geom,Point+(174.201+-40.589),500,meters)+AND+depth>=10+AND+depth<200"

	cql1b = "cql_filter=BBOX(origin_geom,163.60840,-49.18170,182.98828)"
	cql2b = "cql_filter=BBOX(origin_geom,163.6o840,-49.18170,182.g8828,-32.28713)+AND+origintime>='2010-q1-01'+AND+origintime<'2@15-01-01'"
	cql3b = "cql_filter=BBOX(origin_geom,163.60840,-49.18170,182.98828,-32.28713)+AND+magnitude>=e+AND+magnitude<1o0"
	cql4b = "cql_filter=DWITHIN(geom,Peoint+(174.201+-40.589),50o0,km)+AND+depth>=10+AND+depth<2o00"

	routes = wt.Requests{
		//1. geojson routes
		{ID: wt.L(), URL: url_wfs + "&outputFormat=json&" + cql1, Content: CONTENT_TYPE_GeoJSON, Accept: CONTENT_TYPE_GeoJSON},
		{ID: wt.L(), URL: url_wfs + "&outputFormat=json&" + cql2, Content: CONTENT_TYPE_GeoJSON, Accept: CONTENT_TYPE_GeoJSON},
		{ID: wt.L(), URL: url_wfs + "&outputFormat=json&" + cql3, Content: CONTENT_TYPE_GeoJSON, Accept: CONTENT_TYPE_GeoJSON},
		{ID: wt.L(), URL: url_wfs + "&outputFormat=json&" + cql4, Content: CONTENT_TYPE_GeoJSON, Accept: CONTENT_TYPE_GeoJSON},

		//some invalid requests
		{ID: wt.L(), URL: url_wfs + "&outputFormat=json&" + cql1b, Content: CONTENT_TYPE_GeoJSON, Accept: CONTENT_TYPE_GeoJSON, Status: http.StatusBadRequest},
		{ID: wt.L(), URL: url_wfs + "&outputFormat=json&" + cql2b, Content: CONTENT_TYPE_GeoJSON, Accept: CONTENT_TYPE_GeoJSON, Status: http.StatusBadRequest},
		{ID: wt.L(), URL: url_wfs + "&outputFormat=json&" + cql3b, Content: CONTENT_TYPE_GeoJSON, Accept: CONTENT_TYPE_GeoJSON, Status: http.StatusBadRequest},
		{ID: wt.L(), URL: url_wfs + "&outputFormat=json&" + cql4b, Content: CONTENT_TYPE_GeoJSON, Accept: CONTENT_TYPE_GeoJSON, Status: http.StatusBadRequest},

		//2. csv routes
		{ID: wt.L(), URL: url_wfs + "&outputFormat=csv&" + cql1, Content: CONTENT_TYPE_CSV, Accept: CONTENT_TYPE_CSV},
		{ID: wt.L(), URL: url_wfs + "&outputFormat=csv&" + cql2, Content: CONTENT_TYPE_CSV, Accept: CONTENT_TYPE_CSV},
		{ID: wt.L(), URL: url_wfs + "&outputFormat=csv&" + cql3, Content: CONTENT_TYPE_CSV, Accept: CONTENT_TYPE_CSV},
		{ID: wt.L(), URL: url_wfs + "&outputFormat=csv&" + cql4, Content: CONTENT_TYPE_CSV, Accept: CONTENT_TYPE_CSV},

		//3. gml2 routes
		{ID: wt.L(), URL: url_wfs + "&outputFormat=GML2&" + cql1, Content: CONTENT_TYPE_XML, Accept: CONTENT_TYPE_XML},
		{ID: wt.L(), URL: url_wfs + "&outputFormat=GML2&" + cql2, Content: CONTENT_TYPE_XML, Accept: CONTENT_TYPE_XML},
		{ID: wt.L(), URL: url_wfs + "&outputFormat=GML2&" + cql3, Content: CONTENT_TYPE_XML, Accept: CONTENT_TYPE_XML},
		{ID: wt.L(), URL: url_wfs + "&outputFormat=GML2&" + cql4, Content: CONTENT_TYPE_XML, Accept: CONTENT_TYPE_XML},

		//3. gml3 routes
		{ID: wt.L(), URL: url_wfs + "&outputFormat=text/xml;subtype=gml/3.2&" + cql1, Content: CONTENT_TYPE_XML, Accept: CONTENT_TYPE_XML},
		{ID: wt.L(), URL: url_wfs + "&outputFormat=text/xml;subtype=gml/3.2&" + cql2, Content: CONTENT_TYPE_XML, Accept: CONTENT_TYPE_XML},
		{ID: wt.L(), URL: url_wfs + "&outputFormat=text/xml;subtype=gml/3.2&" + cql3, Content: CONTENT_TYPE_XML, Accept: CONTENT_TYPE_XML},
		{ID: wt.L(), URL: url_wfs + "&outputFormat=text/xml;subtype=gml/3.2&" + cql4, Content: CONTENT_TYPE_XML, Accept: CONTENT_TYPE_XML},

		//4. kml routes
		{ID: wt.L(), URL: url_kml + "&" + cql1, Content: CONTENT_TYPE_KML, Accept: CONTENT_TYPE_KML},
		{ID: wt.L(), URL: url_kml + "&" + cql2, Content: CONTENT_TYPE_KML, Accept: CONTENT_TYPE_KML},
		{ID: wt.L(), URL: url_kml + "&" + cql3, Content: CONTENT_TYPE_KML, Accept: CONTENT_TYPE_KML},
		{ID: wt.L(), URL: url_kml + "&" + cql4, Content: CONTENT_TYPE_KML, Accept: CONTENT_TYPE_KML},
	}
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

func TestQuakesGeoJson(t *testing.T) {
	setup(t)
	defer teardown()
	//1. get all quakes cql_filter=BBOX(origin_geom,163.60840,-49.18170,182.98828,-32.28713)
	r := wt.Request{URL: url_wfs + "&outputFormat=json&cql_filter=BBOX(origin_geom,163.60840,-49.18170,182.98828,-32.28713)", Accept: CONTENT_TYPE_GeoJSON}

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
	r = wt.Request{URL: url_wfs + "&outputFormat=json&cql_filter=BBOX(origin_geom,163.60840,-49.18170,182.98828,-32.28713)+AND+origintime>='2010-01-01'+AND+origintime<'2015-01-01'", Accept: CONTENT_TYPE_GeoJSON}

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
		URL:     url_wfs + "&outputFormat=json&cql_filter=BBOX(origin_geom,163.60840,-49.18170,182.98828,-32.28713)+AND+origintime>='2000-01-01'+AND+origintime<'2015-01-01'",
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
