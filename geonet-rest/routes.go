package main

import (
	"fmt"
	"github.com/GeoNet/web"
	"net/http"
	"strings"
)

var (
	muxV1GeoJSON *http.ServeMux
	muxV1JSON    *http.ServeMux
	muxV2GeoJSON *http.ServeMux
	muxV2JSON    *http.ServeMux
	muxDefault   *http.ServeMux
)

func init() {
	muxV1GeoJSON = http.NewServeMux()
	muxV1GeoJSON.HandleFunc("/quake", toHandler(quakesRegionV1))
	muxV1GeoJSON.HandleFunc("/quake/", toHandler(quakeV1))
	muxV1GeoJSON.HandleFunc("/intensity", toHandler(intensityMeasuredLatestV1))
	muxV1GeoJSON.HandleFunc("/felt/report", toHandler(feltV1))

	muxV1JSON = http.NewServeMux()
	muxV1JSON.HandleFunc("/news/geonet", toHandler(newsV1))

	muxV2GeoJSON = http.NewServeMux()
	muxV2GeoJSON.HandleFunc("/intensity", toHandler(intensityV2))
	muxV2GeoJSON.HandleFunc("/quake", toHandler(quakesV2))
	muxV2GeoJSON.HandleFunc("/quake/", toHandler(quakeV2))
	muxV2GeoJSON.HandleFunc("/quake/history/", toHandler(quakeHistoryV2))
	muxV2GeoJSON.HandleFunc("/volcano/val", toHandler(valV2))

	muxV2JSON = http.NewServeMux()
	muxV2JSON.HandleFunc("/news/geonet", toHandler(newsV2))
	muxV2JSON.HandleFunc("/quake/stats", toHandler(quakeStatsV2))

	// muxDefault handles routes with no Accept version.
	muxDefault = http.NewServeMux()
	muxDefault.HandleFunc("/soh", toHandler(soh))
	muxDefault.HandleFunc("/soh/impact", toHandler(impactSOH))
	muxDefault.HandleFunc("/cap/1.2/GPA1.0/quake/", toHandler(capQuake))
	muxDefault.HandleFunc("/cap/1.2/GPA1.0/feed/atom1.0/quake", toHandler(capQuakeFeed))
	// The 'latest' version of the API for unversioned requests.
	muxDefault.HandleFunc("/quake/", toHandler(quakeV2))
	muxDefault.HandleFunc("/quake", toHandler(quakesV2))
	muxDefault.HandleFunc("/quake/history/", toHandler(quakeHistoryV2))
	muxDefault.HandleFunc("/quake/stats", toHandler(quakeStatsV2))
	muxDefault.HandleFunc("/intensity", toHandler(intensityV2))
	muxDefault.HandleFunc("/news/geonet", toHandler(newsV2))
	muxDefault.HandleFunc("/volcano/val", toHandler(valV2))

	for _, v := range []*http.ServeMux{muxV1JSON, muxV2JSON, muxV1GeoJSON, muxV2GeoJSON, muxDefault} {
		v.HandleFunc("/", toHandler(docs))
	}

}

func router(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get("Accept") {
	case web.V2GeoJSON:
		muxV2GeoJSON.ServeHTTP(w, r)
	case web.V1GeoJSON:
		muxV1GeoJSON.ServeHTTP(w, r)
	case web.V1JSON:
		muxV1JSON.ServeHTTP(w, r)
	case web.V2JSON:
		muxV2JSON.ServeHTTP(w, r)
	default:
		muxDefault.ServeHTTP(w, r)
	}
}

func noRoute(w http.ResponseWriter, r *http.Request) {
	badRequest("Can't find a route for this request. Please refer to /api-docs")
}
