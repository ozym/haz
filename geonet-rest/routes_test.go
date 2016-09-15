package main

import (
	wt "github.com/GeoNet/weft/wefttest"
	"net/http"
	"testing"
)

var routes = wt.Requests{
	// GeoJSON routes
	{ID: wt.L(), Accept: V1GeoJSON, Content: V1GeoJSON, Surrogate: maxAge10, URL: "/quake/2013p407387"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: V1GeoJSON, Surrogate: maxAge10, URL: "/felt/report?publicID=2013p407387"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: V1GeoJSON, Surrogate: maxAge10, URL: "/quake?regionID=newzealand&regionIntensity=unnoticeable&number=30&quality=best,caution,good"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: V1GeoJSON, Surrogate: maxAge10, URL: "/quake?regionID=newzealand&regionIntensity=weak&number=30&quality=best,caution,good"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: V1GeoJSON, Surrogate: maxAge10, URL: "/quake?regionID=newzealand&regionIntensity=light&number=30&quality=best,caution,good"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: V1GeoJSON, Surrogate: maxAge10, URL: "/quake?regionID=newzealand&regionIntensity=moderate&number=30&quality=best,caution,good"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: V1GeoJSON, Surrogate: maxAge10, URL: "/quake?regionID=newzealand&regionIntensity=strong&number=30&quality=best,caution,good"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: V1GeoJSON, Surrogate: maxAge10, URL: "/quake?regionID=newzealand&regionIntensity=severe&number=30&quality=best,caution,good"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: V1GeoJSON, Surrogate: maxAge10, URL: "/quake?regionID=newzealand&regionIntensity=unnoticeable&number=3&quality=best,caution,good"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: V1GeoJSON, Surrogate: maxAge10, URL: "/quake?regionID=newzealand&regionIntensity=unnoticeable&number=30&quality=best,caution,good"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: V1GeoJSON, Surrogate: maxAge10, URL: "/quake?regionID=newzealand&regionIntensity=unnoticeable&number=100&quality=best,caution,good"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: V1GeoJSON, Surrogate: maxAge10, URL: "/quake?regionID=newzealand&regionIntensity=unnoticeable&number=500&quality=best,caution,good"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: V1GeoJSON, Surrogate: maxAge10, URL: "/quake?regionID=newzealand&regionIntensity=unnoticeable&number=1000&quality=best,caution,good"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: V1GeoJSON, Surrogate: maxAge10, URL: "/quake?regionID=newzealand&regionIntensity=unnoticeable&number=1500&quality=best,caution,good"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: V1GeoJSON, Surrogate: maxAge10, URL: "/intensity?type=measured"},

	// GeoJSON V2 routes
	{ID: wt.L(), Accept: V2GeoJSON, Content: V2GeoJSON, Surrogate: maxAge10, URL: "/quake/2013p407387"},
	{ID: wt.L(), Accept: V2GeoJSON, Content: V2GeoJSON, Surrogate: maxAge10, URL: "/quake/history/2013p407387"},
	{ID: wt.L(), Accept: V2GeoJSON, Content: V2GeoJSON, Surrogate: maxAge10, URL: "/quake?MMI=-1"},
	{ID: wt.L(), Accept: V2GeoJSON, Content: V2GeoJSON, Surrogate: maxAge10, URL: "/quake?MMI=0"},
	{ID: wt.L(), Accept: V2GeoJSON, Content: V2GeoJSON, Surrogate: maxAge10, URL: "/quake?MMI=1"},
	{ID: wt.L(), Accept: V2GeoJSON, Content: V2GeoJSON, Surrogate: maxAge10, URL: "/quake?MMI=2"},
	{ID: wt.L(), Accept: V2GeoJSON, Content: V2GeoJSON, Surrogate: maxAge10, URL: "/quake?MMI=3"},
	{ID: wt.L(), Accept: V2GeoJSON, Content: V2GeoJSON, Surrogate: maxAge10, URL: "/quake?MMI=4"},
	{ID: wt.L(), Accept: V2GeoJSON, Content: V2GeoJSON, Surrogate: maxAge10, URL: "/quake?MMI=5"},
	{ID: wt.L(), Accept: V2GeoJSON, Content: V2GeoJSON, Surrogate: maxAge10, URL: "/quake?MMI=6"},
	{ID: wt.L(), Accept: V2GeoJSON, Content: V2GeoJSON, Surrogate: maxAge10, URL: "/quake?MMI=7"},
	{ID: wt.L(), Accept: V2GeoJSON, Content: V2GeoJSON, Surrogate: maxAge10, URL: "/quake?MMI=8"},
	{ID: wt.L(), Accept: V2GeoJSON, Content: V2GeoJSON, Surrogate: maxAge10, URL: "/intensity?type=measured"},
	{ID: wt.L(), Accept: V2GeoJSON, Content: V2GeoJSON, Surrogate: maxAge10, URL: "/intensity?type=reported"},
	{ID: wt.L(), Accept: V2GeoJSON, Content: V2GeoJSON, Surrogate: maxAge10, URL: "/intensity?type=reported&publicID=2013p407387"},
	{ID: wt.L(), Accept: V2GeoJSON, Content: V2GeoJSON, Surrogate: maxAge10, URL: "/volcano/val"},

	// GeoJSON routes without explicit accept should route to latest version
	{ID: wt.L(), Accept: "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8", Content: V2GeoJSON, Surrogate: maxAge10, URL: "/quake/2013p407387"},
	{ID: wt.L(), Accept: "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8", Content: V2GeoJSON, Surrogate: maxAge10, URL: "/quake?MMI=3"},
	{ID: wt.L(), Accept: "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8", Content: V2GeoJSON, Surrogate: maxAge10, URL: "/intensity?type=measured"},
	{ID: wt.L(), Accept: "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8", Content: V2GeoJSON, Surrogate: maxAge10, URL: "/intensity?type=reported"},
	{ID: wt.L(), Accept: "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8", Content: V2GeoJSON, Surrogate: maxAge10, URL: "/volcano/val"},

	// Routes that should 404
	{ID: wt.L(), Accept: V1GeoJSON, Content: ErrContent, Status: http.StatusNotFound, Surrogate: maxAge10, URL: "/quake/2013p407399"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: ErrContent, Status: http.StatusNotFound, Surrogate: maxAge10, URL: "/felt/report?publicID=2013p407399"},

	// JSON routes
	{ID: wt.L(), Accept: V1JSON, Content: V1JSON, Surrogate: maxAge300, URL: "/news/geonet"},

	// V2 JSON routes
	{ID: wt.L(), Accept: V2JSON, Content: V2JSON, Surrogate: maxAge300, URL: "/news/geonet"},
	{ID: wt.L(), Accept: V2JSON, Content: V2JSON, Surrogate: maxAge300, URL: "/quake/stats"},

	// JSON routes without explicit accept should route to latest version
	{ID: wt.L(), Accept: "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8", Content: V2JSON, Surrogate: maxAge300, URL: "/news/geonet"},

	// CAP routes - not versioned by Accept
	{ID: wt.L(), Content: CAP, Surrogate: maxAge10, URL: "/cap/1.2/GPA1.0/quake/2013p407387.1370036261549894"},

	// Atom feed routes - not versioned by Accept
	{ID: wt.L(), Content: Atom, Surrogate: maxAge10, URL: "/cap/1.2/GPA1.0/feed/atom1.0/quake"},

	// GeoJSON routes that should bad request
	{ID: wt.L(), Accept: V1GeoJSON, Content: ErrContent, Surrogate: maxAge86400, Status: http.StatusBadRequest, URL: "/quake?regionID=newzealand&regionIntensity=bad&number=30&quality=best,caution,good"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: ErrContent, Surrogate: maxAge86400, Status: http.StatusBadRequest, URL: "/quake?regionID=newzealand&regionIntensity=unnoticeable&number=30&quality=best,caution,bad"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: ErrContent, Surrogate: maxAge86400, Status: http.StatusBadRequest, URL: "/quake?regionID=newzealand&regionIntensity=unnoticeable&number=999&quality=best,caution,good"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: ErrContent, Surrogate: maxAge86400, Status: http.StatusBadRequest, URL: "/quake?regionID=newzealand&regionIntensity=unnoticeable&quality=best,caution,good"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: ErrContent, Surrogate: maxAge86400, Status: http.StatusBadRequest, URL: "/quake?regionID=newzealand&regionIntensity=unnoticeable"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: ErrContent, Surrogate: maxAge86400, Status: http.StatusBadRequest, URL: "/quake?regionID=newzealand"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: ErrContent, Surrogate: maxAge86400, Status: http.StatusBadRequest, URL: "/quake"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: ErrContent, Surrogate: maxAge86400, Status: http.StatusBadRequest, URL: "/quake?regionID=ruapehu&regionIntensity=unnoticeable&number=3&quality=best,caution,good"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: ErrContent, Surrogate: maxAge86400, Status: http.StatusBadRequest, URL: "/quake?regionID=bad&regionIntensity=unnoticeable&number=3&quality=best,caution,good"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: ErrContent, Surrogate: maxAge86400, Status: http.StatusBadRequest, URL: "/quake?regionID=newzealand&intensity=bad&number=30&quality=best,caution,good"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: ErrContent, Surrogate: maxAge86400, Status: http.StatusBadRequest, URL: "/quake?regionID=newzealand&intensity=unnoticeable&number=30&quality=best,caution,bad"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: ErrContent, Surrogate: maxAge86400, Status: http.StatusBadRequest, URL: "/quake?regionID=newzealand&intensity=unnoticeable&number=999&quality=best,caution,good"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: ErrContent, Surrogate: maxAge86400, Status: http.StatusBadRequest, URL: "/quake?regionID=newzealand&intensity=unnoticeable&quality=best,caution,good"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: ErrContent, Surrogate: maxAge86400, Status: http.StatusBadRequest, URL: "/quake?regionID=newzealand&intensity=unnoticeable"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: ErrContent, Surrogate: maxAge86400, Status: http.StatusBadRequest, URL: "/quake?regionID=newzealand"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: ErrContent, Surrogate: maxAge86400, Status: http.StatusBadRequest, URL: "/quake"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: ErrContent, Surrogate: maxAge86400, Status: http.StatusBadRequest, URL: "/quake?regionID=ruapehu&intensity=unnoticeable&number=3&quality=best,caution,good"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: ErrContent, Surrogate: maxAge86400, Status: http.StatusBadRequest, URL: "/quake?regionID=bad&intensity=unnoticeable&number=3&quality=best,caution,good"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: ErrContent, Surrogate: maxAge86400, Status: http.StatusBadRequest, URL: "/region/bad"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: ErrContent, Surrogate: maxAge86400, Status: http.StatusBadRequest, URL: "/region?type=badQuery"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: ErrContent, Surrogate: maxAge86400, Status: http.StatusBadRequest, URL: "/fred"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: ErrContent, Surrogate: maxAge86400, Status: http.StatusBadRequest, URL: "/felt/report?quakeID=2012p498491"},
	{ID: wt.L(), Accept: V1GeoJSON, Content: ErrContent, Surrogate: maxAge86400, Status: http.StatusBadRequest, URL: "/intensity?type=reported"}, // no reported at V1

	// V2 GeoJSON routes that should bad request
	{ID: wt.L(), Accept: V2GeoJSON, Content: ErrContent, Surrogate: maxAge86400, Status: http.StatusBadRequest, URL: "/quake?MMI=9"},
	{ID: wt.L(), Accept: V2GeoJSON, Content: ErrContent, Surrogate: maxAge86400, Status: http.StatusBadRequest, URL: "/quake?MMI=-2"},

	// soh routes
	{ID: wt.L(), URL: "/soh"},
	{ID: wt.L(), URL: "/soh/up"},
	{ID: wt.L(), URL: "/soh/esb"},
	{ID: wt.L(), Status: http.StatusServiceUnavailable, URL: "/soh/impact"}, // not enough data so gets an error
}

// Test all routes give the expected response.  Also check with
// cache busters and extra query parameters.
func TestRoutes(t *testing.T) {
	setup()
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
