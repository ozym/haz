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
	r.Add("/intensity?type=measured")

	r.Test(ts, t)

	// GeoJSON V2 routes
	r = webtest.Route{
		Accept:     web.V2GeoJSON,
		Content:    web.V2GeoJSON,
		Cache:      web.MaxAge10,
		Surrogate:  web.MaxAge10,
		Response:   http.StatusOK,
		Vary:       "Accept",
		TestAccept: false,
	}
	r.Add("/quake/2013p407387")
	r.Add("/quake?MMI=-1")
	r.Add("/quake?MMI=0")
	r.Add("/quake?MMI=1")
	r.Add("/quake?MMI=2")
	r.Add("/quake?MMI=3")
	r.Add("/quake?MMI=4")
	r.Add("/quake?MMI=5")
	r.Add("/quake?MMI=6")
	r.Add("/quake?MMI=7")
	r.Add("/quake?MMI=8")
	r.Add("/intensity?type=measured")
	r.Add("/intensity?type=reported")
	r.Add("/intensity?type=reported&publicID=2013p407387")

	r.Test(ts, t)

	// GeoJSON routes without explicit accept should route to latest version
	r = webtest.Route{
		Accept:     "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		Content:    web.V2GeoJSON,
		Cache:      web.MaxAge10,
		Surrogate:  web.MaxAge10,
		Response:   http.StatusOK,
		Vary:       "Accept",
		TestAccept: false,
	}
	r.Add("/quake/2013p407387")
	r.Add("/quake?MMI=3")
	r.Add("/intensity?type=measured")
	r.Add("/intensity?type=reported")

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

	r.Test(ts, t)

	// V2 JSON routes
	r = webtest.Route{
		Accept:     web.V2JSON,
		Content:    web.V2JSON,
		Cache:      web.MaxAge10,
		Surrogate:  web.MaxAge300,
		Response:   http.StatusOK,
		Vary:       "Accept",
		TestAccept: false,
	}
	r.Add("/news/geonet")

	r.Test(ts, t)

	// JSON routes without explicit accept should route to latest version
	r = webtest.Route{
		Accept:     "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		Content:    web.V2JSON,
		Cache:      web.MaxAge10,
		Surrogate:  web.MaxAge300,
		Response:   http.StatusOK,
		Vary:       "Accept",
		TestAccept: false,
	}
	r.Add("/news/geonet")
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
	r.Add("/intensity?type=reported") // no reported at V1
	r.Test(ts, t)

	// V2 GeoJSON routes that should bad request
	r = webtest.Route{
		Accept:     web.V2GeoJSON,
		Content:    web.ErrContent,
		Cache:      web.MaxAge10,
		Surrogate:  web.MaxAge86400,
		Response:   http.StatusBadRequest,
		Vary:       "Accept",
		TestAccept: false,
	}
	r.Add("/quake?MMI=9")
	r.Add("/quake?MMI=-2")

	r.Test(ts, t)

}
