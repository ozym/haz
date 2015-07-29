package main

import (
	"github.com/GeoNet/web"
	"github.com/GeoNet/web/webtest"
	"net/http"
	"testing"
)

func TestRoutes(t *testing.T) {
	setup()
	defer teardown()

	// GeoJSON routes
	r := webtest.Route{
		Accept:     web.V1GeoJSON,
		Content:    web.V1GeoJSON,
		Cache:      web.MaxAge10,
		Surrogate:  web.MaxAge10,
		Response:   http.StatusOK,
		Vary:       "Accept",
		TestAccept: false,
	}
	r.Add("/quake/2013p407387")
	r.Add("/felt/report?publicID=2013p407387")
	r.Add("/quake?regionID=newzealand&regionIntensity=unnoticeable&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&regionIntensity=weak&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&regionIntensity=light&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&regionIntensity=moderate&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&regionIntensity=strong&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&regionIntensity=severe&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&regionIntensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&regionIntensity=unnoticeable&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&regionIntensity=unnoticeable&number=100&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&regionIntensity=unnoticeable&number=500&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&regionIntensity=unnoticeable&number=1000&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&regionIntensity=unnoticeable&number=1500&quality=best,caution,good")
	r.Add("/quake?regionID=aucklandnorthland&regionIntensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=tongagrirobayofplenty&regionIntensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=gisborne&regionIntensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=hawkesbay&regionIntensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=taranaki&regionIntensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=wellington&regionIntensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=nelsonwestcoast&regionIntensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=canterbury&regionIntensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=fiordland&regionIntensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=otagosouthland&regionIntensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&intensity=unnoticeable&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&intensity=weak&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&intensity=light&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&intensity=moderate&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&intensity=strong&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&intensity=severe&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&intensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&intensity=unnoticeable&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&intensity=unnoticeable&number=100&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&intensity=unnoticeable&number=500&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&intensity=unnoticeable&number=1000&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&intensity=unnoticeable&number=1500&quality=best,caution,good")
	r.Add("/quake?regionID=aucklandnorthland&intensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=tongagrirobayofplenty&intensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=gisborne&intensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=hawkesbay&intensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=taranaki&intensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=wellington&intensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=nelsonwestcoast&intensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=canterbury&intensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=fiordland&intensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=otagosouthland&intensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/intensity?type=measured")
	// r.Add("/intensity?type=reported&zoom=5")
	// r.Add("/intensity?type=reported&zoom=5&publicID=2012p673624")
	r.Add("/volcano/alert/level")

	r.Test(ts, t)

	// GeoJSON routes without explicit accept should route to latest version
	r = webtest.Route{
		Accept:     "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		Content:    web.V1GeoJSON,
		Cache:      web.MaxAge10,
		Surrogate:  web.MaxAge10,
		Response:   http.StatusOK,
		Vary:       "Accept",
		TestAccept: false,
	}
	r.Add("/quake/2013p407387")
	r.Add("/felt/report?publicID=2013p407387")
	r.Add("/quake?regionID=newzealand&regionIntensity=unnoticeable&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&regionIntensity=weak&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&regionIntensity=light&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&regionIntensity=moderate&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&regionIntensity=strong&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&regionIntensity=severe&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&regionIntensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&regionIntensity=unnoticeable&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&regionIntensity=unnoticeable&number=100&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&regionIntensity=unnoticeable&number=500&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&regionIntensity=unnoticeable&number=1000&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&regionIntensity=unnoticeable&number=1500&quality=best,caution,good")
	r.Add("/quake?regionID=aucklandnorthland&regionIntensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=tongagrirobayofplenty&regionIntensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=gisborne&regionIntensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=hawkesbay&regionIntensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=taranaki&regionIntensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=wellington&regionIntensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=nelsonwestcoast&regionIntensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=canterbury&regionIntensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=fiordland&regionIntensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=otagosouthland&regionIntensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&intensity=unnoticeable&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&intensity=weak&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&intensity=light&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&intensity=moderate&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&intensity=strong&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&intensity=severe&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&intensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&intensity=unnoticeable&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&intensity=unnoticeable&number=100&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&intensity=unnoticeable&number=500&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&intensity=unnoticeable&number=1000&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&intensity=unnoticeable&number=1500&quality=best,caution,good")
	r.Add("/quake?regionID=aucklandnorthland&intensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=tongagrirobayofplenty&intensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=gisborne&intensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=hawkesbay&intensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=taranaki&intensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=wellington&intensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=nelsonwestcoast&intensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=canterbury&intensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=fiordland&intensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=otagosouthland&intensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/intensity?type=measured")
	// r.Add("/intensity?type=reported&zoom=5")
	// r.Add("/intensity?type=reported&zoom=5&publicID=2012p673624")
	r.Add("/volcano/alert/level")

	r.Test(ts, t)

	// GeoJSON routes with long cache times
	r = webtest.Route{
		Accept:     web.V1GeoJSON,
		Content:    web.V1GeoJSON,
		Cache:      web.MaxAge10,
		Surrogate:  web.MaxAge86400,
		Response:   http.StatusOK,
		Vary:       "Accept",
		TestAccept: false,
	}
	r.Add("/region/newzealand")
	r.Add("/region/aucklandnorthland")
	r.Add("/region/tongagrirobayofplenty")
	r.Add("/region/gisborne")
	r.Add("/region/hawkesbay")
	r.Add("/region/taranaki")
	r.Add("/region/wellington")
	r.Add("/region/nelsonwestcoast")
	r.Add("/region/canterbury")
	r.Add("/region/fiordland")
	r.Add("/region/otagosouthland")

	r.Test(ts, t)

	// Routes that should 404
	r = webtest.Route{
		Accept:     web.V1GeoJSON,
		Content:    web.ErrContent,
		Cache:      web.MaxAge10,
		Surrogate:  web.MaxAge10,
		Response:   http.StatusNotFound,
		Vary:       "Accept",
		TestAccept: false,
	}
	r.Add("/quake/2013p407399")
	r.Add("/felt/report?publicID=2013p407399")

	r.Test(ts, t)

	// JSON routes
	r = webtest.Route{
		Accept:     web.V1JSON,
		Content:    web.V1JSON,
		Cache:      web.MaxAge10,
		Surrogate:  web.MaxAge300,
		Response:   http.StatusOK,
		Vary:       "Accept",
		TestAccept: false,
	}
	r.Add("/news/geonet")
	r.Add("/volcano/alert/bulletin")

	r.Test(ts, t)

	// CAP routes - not versioned by Accept
	r = webtest.Route{
		Content:    web.CAP,
		Cache:      web.MaxAge10,
		Surrogate:  web.MaxAge10,
		Response:   http.StatusOK,
		Vary:       "Accept",
		TestAccept: false,
	}
	r.Add("/cap/1.2/GPA1.0/quake/2013p407387.1370036261549894")

	r.Test(ts, t)

	// Atom feed routes - not versioned by Accept
	r = webtest.Route{
		Content:    web.Atom,
		Cache:      web.MaxAge10,
		Surrogate:  web.MaxAge10,
		Response:   http.StatusOK,
		Vary:       "Accept",
		TestAccept: false,
	}
	r.Add("/cap/1.2/GPA1.0/feed/atom1.0/quake")

	r.Test(ts, t)

	// GeoJSON routes that should bad request
	r = webtest.Route{
		Accept:     web.V1GeoJSON,
		Content:    web.ErrContent,
		Cache:      web.MaxAge10,
		Surrogate:  web.MaxAge86400,
		Response:   http.StatusBadRequest,
		Vary:       "Accept",
		TestAccept: false,
	}
	r.Add("/quake?regionID=newzealand&regionIntensity=bad&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&regionIntensity=unnoticeable&number=30&quality=best,caution,bad")
	r.Add("/quake?regionID=newzealand&regionIntensity=unnoticeable&number=999&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&regionIntensity=unnoticeable&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&regionIntensity=unnoticeable")
	r.Add("/quake?regionID=newzealand")
	r.Add("/quake")
	r.Add("/quake?regionID=ruapehu&regionIntensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=bad&regionIntensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&intensity=bad&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&intensity=unnoticeable&number=30&quality=best,caution,bad")
	r.Add("/quake?regionID=newzealand&intensity=unnoticeable&number=999&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&intensity=unnoticeable&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&intensity=unnoticeable")
	r.Add("/quake?regionID=newzealand")
	r.Add("/quake")
	r.Add("/quake?regionID=ruapehu&intensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/quake?regionID=bad&intensity=unnoticeable&number=3&quality=best,caution,good")
	r.Add("/region/bad")
	r.Add("/region?type=badQuery")
	r.Add("/")
	r.Add("/felt/report?quakeID=2012p498491")
	r.Test(ts, t)

}

func TestGeoJSON(t *testing.T) {
	setup()
	defer teardown()

	// GeoJSON routes
	r := webtest.Route{
		Accept:     web.V1GeoJSON,
		Content:    web.V1GeoJSON,
		Cache:      web.MaxAge10,
		Surrogate:  web.MaxAge10,
		Response:   http.StatusOK,
		Vary:       "Accept",
		TestAccept: false,
	}
	r.Add("/quake/2013p407387")
	r.Add("/quake?regionID=newzealand&regionIntensity=unnoticeable&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=wellington&regionIntensity=severe&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=newzealand&intensity=unnoticeable&number=30&quality=best,caution,good")
	r.Add("/quake?regionID=wellington&intensity=severe&number=30&quality=best,caution,good")
	r.Add("/region/tongagrirobayofplenty")
	r.Add("/region?type=quake")
	r.Add("/felt/report?publicID=2013p407387")
	r.Add("/intensity?type=measured")
	// r.Add("/intensity?type=reported&zoom=5")
	// r.Add("/intensity?type=reported&zoom=5&publicID=2012p673624")
	r.Add("/volcano/alert/level")

	r.GeoJSON(ts, t)
}
