package main

import (

	"bytes"
	"github.com/GeoNet/weft"
	"html/template"
	"net/http"
	"os"
	"log"
)

var (
	indexTemp *template.Template
	serveMux  *http.ServeMux
)

func init() {
	indexTemp = template.Must(template.ParseFiles("assets/tmpl/index.html"))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	serveMux = http.NewServeMux()
	serveMux.HandleFunc("/geojson", weft.MakeHandlerAPI(getQuakesGeoJson))
	serveMux.HandleFunc("/count", weft.MakeHandlerAPI(getQuakesCount))
	serveMux.HandleFunc("/csv", weft.MakeHandlerAPI(getQuakesCsv))
	serveMux.HandleFunc("/gml", weft.MakeHandlerAPI(getQuakesGml))
	serveMux.HandleFunc("/kml", weft.MakeHandlerAPI(getQuakesKml))
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

	var p searchPage
	p.ApiKey = os.Getenv("BING_API_KEY")

	err := indexTemp.ExecuteTemplate(b, "base", p)

	if err != nil {
		return weft.InternalServerError(err)
	}

	return &weft.StatusOK
}

type searchPage struct {
	ApiKey string
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
}

