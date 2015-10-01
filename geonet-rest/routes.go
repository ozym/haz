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
	muxV1GeoJSON.HandleFunc("/quake", quakesRegionV1)
	muxV1GeoJSON.HandleFunc("/quake/", quakeV1)
	muxV1GeoJSON.HandleFunc("/intensity", intensityMeasuredLatestV1)
	muxV1GeoJSON.HandleFunc("/felt/report", feltV1)

	muxV1JSON = http.NewServeMux()
	muxV1JSON.HandleFunc("/news/geonet", newsV1)

	muxV2GeoJSON = http.NewServeMux()
	muxV2GeoJSON.HandleFunc("/intensity", intensityV2)
	muxV2GeoJSON.HandleFunc("/quake", quakesV2)
	muxV2GeoJSON.HandleFunc("/quake/", quakeV2)
	muxV2GeoJSON.HandleFunc("/quake/history/", quakeHistoryV2)

	muxV2JSON = http.NewServeMux()
	muxV2JSON.HandleFunc("/news/geonet", newsV2)

	// muxDefault handles routes with no Accept version.
	muxDefault = http.NewServeMux()
	muxDefault.HandleFunc("/soh", soh)
	muxDefault.HandleFunc("/soh/impact", impactSOH)
	muxDefault.HandleFunc("/api-docs", docs)
	muxDefault.HandleFunc("/cap/1.2/GPA1.0/quake/", capQuake)
	muxDefault.HandleFunc("/cap/1.2/GPA1.0/feed/atom1.0/quake", capQuakeFeed)
	// The 'latest' version of the API for unversioned requests.
	muxDefault.HandleFunc("/quake/", quakeV2)
	muxDefault.HandleFunc("/quake", quakesV2)
	muxDefault.HandleFunc("/quake/history/", quakeHistoryV2)
	muxDefault.HandleFunc("/intensity", intensityV2)
	muxDefault.HandleFunc("/news/geonet", newsV2)

	for _, v := range []*http.ServeMux{muxV1JSON, muxV2JSON, muxV1GeoJSON, muxV2GeoJSON, muxDefault} {
		v.HandleFunc("/", noRoute)
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
	web.BadRequest(w, r, "Can't find a route for this request. Please refer to /api-docs")
}

func badQuery(w http.ResponseWriter, r *http.Request, required, optional []string) bool {
	v := r.URL.Query()

	var err error
	var missing []string

	for _, k := range required {
		if v.Get(k) == "" {
			missing = append(missing, k)

		}
	}

	switch len(missing) {
	case 0:
	case 1:
		err = fmt.Errorf("missing query parameter: " + missing[0])
	default:
		err = fmt.Errorf("missing query parameters: " + strings.Join(missing, ", "))
	}

	for _, k := range required {
		v.Del(k)
	}

	for _, k := range optional {
		v.Del(k)
	}

	if len(v) > 0 {
		web.BadRequest(w, r, "found additional query parameters")
		return true
	}

	if err != nil {
		web.BadRequest(w, r, err.Error())
		return true
	}

	return false
}
