package main

import (
	"bufio"
	"bytes"
	"github.com/GeoNet/weft"
	"html/template"
	"net/http"
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

	w := bufio.NewWriter(b)
	err := indexTemp.Execute(w, nil)

	if err != nil {
		return weft.InternalServerError(err)
	}

	return &weft.StatusOK
}
