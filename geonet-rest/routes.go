package main

import (
	"github.com/GeoNet/web"
	"github.com/GeoNet/web/api/apidoc"
	"net/http"
	"strings"
)

var docs = apidoc.Docs{
	Production: config.WebServer.Production,
	APIHost:    config.WebServer.CNAME,
	Title:      `GeoNet API`,
	Description: `<p>The data provided here is used for the GeoNet web site and other similar services. 
			If you are looking for data for research or other purposes then please check the 
			<a href="http://info.geonet.org.nz/x/DYAO">full range of data</a> available from GeoNet. </p>`,
	RepoURL:          `https://github.com/GeoNet/geonet-rest`,
	StrictVersioning: false,
}

func init() {
	docs.AddEndpoint("quake", &quakeDoc)
	docs.AddEndpoint("region", &regionDoc)
	docs.AddEndpoint("felt", &feltDoc)
	docs.AddEndpoint("news", &newsDoc)
	docs.AddEndpoint("impact", &impactDoc)
	docs.AddEndpoint("volcano", &volcanoDoc)
	docs.AddEndpoint("cap", &capDoc)
}

var exHost = "http://localhost:" + config.WebServer.Port

func router(w http.ResponseWriter, r *http.Request) {
	// requests that don't have a specific version header are routed to the latest version.
	var latest bool
	accept := r.Header.Get("Accept")
	switch accept {
	case web.V1GeoJSON, web.V1JSON:
	default:
		latest = true
	}

	switch {
	case strings.HasPrefix(r.URL.Path, "/quake") && (accept == web.V1GeoJSON || latest):
		w.Header().Set("Content-Type", web.V1GeoJSON)
		switch {
		case r.URL.Query().Get("intensity") != "":
			quakes(w, r)
		case r.URL.Query().Get("regionIntensity") != "":
			quakesRegion(w, r)
		case strings.HasPrefix(r.URL.Path, "/quake/"):
			quake(w, r)
		default:
			web.BadRequest(w, r, "Can't find a route for this request. Please refer to /api-docs")
		}
	case r.URL.Path == "/intensity" && (accept == web.V1GeoJSON || latest):
		w.Header().Set("Content-Type", web.V1GeoJSON)
		switch {
		case r.URL.Query().Get("type") == "measured":
			intensityMeasuredLatest(w, r)
		// case r.URL.Query().Get("type") == "reported" && r.URL.Query().Get("publicID") == "":
		// 	intensityReportedLatest(w, r)
		// case r.URL.Query().Get("type") == "reported" && r.URL.Query().Get("publicID") != "":
		// 	intensityReported(w, r)
		default:
			web.BadRequest(w, r, "Can't find a route for this request. Please refer to /api-docs")
		}
	case r.URL.Path == "/felt/report" && (accept == web.V1GeoJSON || latest):
		w.Header().Set("Content-Type", web.V1GeoJSON)
		felt(w, r)
	case r.URL.Path == "/volcano/alert/level" && (accept == web.V1GeoJSON || latest):
		w.Header().Set("Content-Type", web.V1GeoJSON)
		alertLevel(w, r)
	case r.URL.Path == "/volcano/alert/bulletin" && (accept == web.V1JSON || latest):
		w.Header().Set("Content-Type", web.V1JSON)
		alertBulletin(w, r)
	case strings.HasPrefix(r.URL.Path, "/region/") && (accept == web.V1GeoJSON || latest):
		w.Header().Set("Content-Type", web.V1GeoJSON)
		region(w, r)
	case r.URL.Path == "/region" && (accept == web.V1GeoJSON || latest):
		w.Header().Set("Content-Type", web.V1GeoJSON)
		regions(w, r)
	case r.URL.Path == "/news/geonet" && (accept == web.V1JSON || latest):
		w.Header().Set("Content-Type", web.V1JSON)
		news(w, r)
	case strings.HasPrefix(r.URL.Path, "/cap/1.2/GPA1.0/quake"):
		w.Header().Set("Content-Type", web.CAP)
		capQuake(w, r)
	case r.URL.Path == "/cap/1.2/GPA1.0/feed/atom1.0/quake":
		w.Header().Set("Content-Type", web.Atom)
		capQuakeFeed(w, r)
	case strings.HasPrefix(r.URL.Path, apidoc.Path):
		docs.Serve(w, r)
	case r.URL.Path == "/soh":
		soh(w, r)
	case r.URL.Path == "/soh/impact":
		impactSOH(w, r)
	default:
		web.BadRequest(w, r, "Can't find a route for this request. Please refer to /api-docs")
	}
}
