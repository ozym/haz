package main

import (
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

	/**
	  Request example:
	  1. ows
	   http://wfs.geonet.org.nz/geonet/ows?service=WFS&version=1.0.0&request=GetFeature&typeName=geonet:quake_search_v1
	    &outputFormat=json&cql_filter=origintime>='2013-06-01'+AND+origintime<'2014-01-01'+AND+usedphasecount>60

	  2. wms
	   http://wfs.geonet.org.nz/geonet/wms/kml?layers=geonet:quake_search_v1&maxFeatures=50
	*/
	serveMux = http.NewServeMux()
	serveMux.HandleFunc("/wms/kml", weft.MakeHandlerAPI(getQuakesKml))
	serveMux.HandleFunc("/ows", weft.MakeHandlerAPI(getQuakesWfs))
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

	err := indexTemp.Execute(b, nil)

	if err != nil {
		return weft.InternalServerError(err)
	}

	return &weft.StatusOK
}
