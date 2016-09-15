package main

import (

	"bytes"
	"github.com/GeoNet/weft"
	"html/template"
	"net/http"
	"log"
)

var (
	indexTemp *template.Template
	serveMux  *http.ServeMux
)

func init() {
	indexTemp = template.Must(template.ParseFiles("assets/tmpl/index.html"))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	/**
	  Request example:
	  1. ows
	   http://wfs.geonet.org.nz/geonet/ows?service=WFS&version=1.0.0&request=GetFeature&typeName=geonet:quake_search_v1
	    &outputFormat=json&cql_filter=origintime>='2013-06-01'+AND+origintime<'2014-01-01'+AND+usedphasecount>60

	  2. wms
	   http://wfs.geonet.org.nz/geonet/wms/kml?layers=geonet:quake_search_v1&maxFeatures=50
	*/
	serveMux = http.NewServeMux()
	serveMux.HandleFunc("/geonet/wms/kml", weft.MakeHandlerAPI(getQuakesKml))
	serveMux.HandleFunc("/geonet/ows", weft.MakeHandlerAPI(getQuakesWfs))
	serveMux.HandleFunc("/", weft.MakeHandlerPage(indexPage))

	// routes for balancers and probes.
	serveMux.HandleFunc("/soh/up", http.HandlerFunc(up))
	serveMux.HandleFunc("/soh", http.HandlerFunc(soh))
}

func router(w http.ResponseWriter, r *http.Request) {
	serveMux.ServeHTTP(w, r)
}

func indexPage(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
		return res
	}
	if r.URL.Path != "/" {
		return weft.BadRequest("invalid path")
	}

	err := indexTemp.Execute(b, nil)

	if err != nil {
		return weft.InternalServerError(err)
	}

	return &weft.StatusOK
}

// up is for testing that the app has started e.g., for with load balancers.
// It indicates the app is started.  It may still be serving errors.
// Not useful for inclusion in app metrics so weft not used.
func up(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
		w.Header().Set("Surrogate-Control", "max-age=86400")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write([]byte("<html><head></head><body>up</body></html>"))
	log.Print("up ok")
}

// soh is for external service probes.
// writes a service unavailable error to w if the service is not working.
// Not useful for inclusion in app metrics so weft not used.
func soh(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
		w.Header().Set("Surrogate-Control", "max-age=86400")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var c int

	if err := db.QueryRow("SELECT 1").Scan(&c); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("<html><head></head><body>service error</body></html>"))
		log.Printf("ERROR: soh service error %s", err)
		return
	}

	w.Write([]byte("<html><head></head><body>ok</body></html>"))
	log.Print("soh ok")
}
