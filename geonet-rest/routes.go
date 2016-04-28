package main

import (
	"net/http"
	"github.com/GeoNet/weft"
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
	muxV1GeoJSON.HandleFunc("/quake", weft.MakeHandlerAPI(quakesRegionV1))
	muxV1GeoJSON.HandleFunc("/quake/", weft.MakeHandlerAPI(quakeV1))
	muxV1GeoJSON.HandleFunc("/intensity", weft.MakeHandlerAPI(intensityMeasuredLatestV1))
	muxV1GeoJSON.HandleFunc("/felt/report", weft.MakeHandlerAPI(feltV1))

	muxV1JSON = http.NewServeMux()
	muxV1JSON.HandleFunc("/news/geonet", weft.MakeHandlerAPI(newsV1))

	muxV2GeoJSON = http.NewServeMux()
	muxV2GeoJSON.HandleFunc("/intensity", weft.MakeHandlerAPI(intensityV2))
	muxV2GeoJSON.HandleFunc("/quake", weft.MakeHandlerAPI(quakesV2))
	muxV2GeoJSON.HandleFunc("/quake/", weft.MakeHandlerAPI(quakeV2))
	muxV2GeoJSON.HandleFunc("/quake/history/", weft.MakeHandlerAPI(quakeHistoryV2))
	muxV2GeoJSON.HandleFunc("/volcano/val", weft.MakeHandlerAPI(valV2))

	muxV2JSON = http.NewServeMux()
	muxV2JSON.HandleFunc("/news/geonet", weft.MakeHandlerAPI(newsV2))
	muxV2JSON.HandleFunc("/quake/stats", weft.MakeHandlerAPI(quakeStatsV2))

	// muxDefault handles routes with no Accept version.
	muxDefault = http.NewServeMux()
	muxDefault.HandleFunc("/soh", weft.MakeHandlerAPI(soh))
	muxDefault.HandleFunc("/soh/impact", weft.MakeHandlerAPI(impactSOH))
	muxDefault.HandleFunc("/cap/1.2/GPA1.0/quake/", weft.MakeHandlerAPI(capQuake))
	muxDefault.HandleFunc("/cap/1.2/GPA1.0/feed/atom1.0/quake", weft.MakeHandlerAPI(capQuakeFeed))
	// The 'latest' version of the API for unversioned requests.
	muxDefault.HandleFunc("/quake/", weft.MakeHandlerAPI(quakeV2))
	muxDefault.HandleFunc("/quake", weft.MakeHandlerAPI(quakesV2))
	muxDefault.HandleFunc("/quake/history/", weft.MakeHandlerAPI(quakeHistoryV2))
	muxDefault.HandleFunc("/quake/stats", weft.MakeHandlerAPI(quakeStatsV2))
	muxDefault.HandleFunc("/intensity", weft.MakeHandlerAPI(intensityV2))
	muxDefault.HandleFunc("/news/geonet", weft.MakeHandlerAPI(newsV2))
	muxDefault.HandleFunc("/volcano/val", weft.MakeHandlerAPI(valV2))

	for _, v := range []*http.ServeMux{muxV1JSON, muxV2JSON, muxV1GeoJSON, muxV2GeoJSON, muxDefault} {
		v.HandleFunc("/", weft.MakeHandlerPage(docs))
	}

}

func router(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get("Accept") {
	case V2GeoJSON:
		muxV2GeoJSON.ServeHTTP(w, r)
	case V1GeoJSON:
		muxV1GeoJSON.ServeHTTP(w, r)
	case V1JSON:
		muxV1JSON.ServeHTTP(w, r)
	case V2JSON:
		muxV2JSON.ServeHTTP(w, r)
	default:
		muxDefault.ServeHTTP(w, r)
	}
}
