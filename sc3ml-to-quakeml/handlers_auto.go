package main

// This file is auto generated - do not edit.
// It was created with weftgenapi from github.com/GeoNet/weft/weftgenapi

import (
	"bytes"
	"github.com/GeoNet/weft"
	"io/ioutil"
	"net/http"
)

var mux = http.NewServeMux()

func init() {
	mux.HandleFunc("/api-docs", weft.MakeHandlerPage(docHandler))
	mux.HandleFunc("/csv/1.0.0/", weft.MakeHandlerAPI(csv1_0_0sHandler))
	mux.HandleFunc("/quakeml/1.2/", weft.MakeHandlerAPI(quakeml1_2sHandler))
	mux.HandleFunc("/quakeml-rt/1.2/", weft.MakeHandlerAPI(quakeml_rt1_2sHandler))
}

func docHandler(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	switch r.Method {
	case "GET":
		by, err := ioutil.ReadFile("assets/api-docs/index.html")
		if err != nil {
			return weft.InternalServerError(err)
		}
		b.Write(by)
		return &weft.StatusOK
	default:
		return &weft.MethodNotAllowed
	}
}
func csv1_0_0sHandler(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	switch r.Method {
	case "GET":
		switch r.Header.Get("Accept") {
		case "text/csv":
			if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
				return res
			}
			h.Set("Content-Type", "text/csv")
			return csv(r, h, b)
		default:
			if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
				return res
			}
			h.Set("Content-Type", "text/csv")
			return csv(r, h, b)
		}
	default:
		return &weft.MethodNotAllowed
	}
}

func quakeml1_2sHandler(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	switch r.Method {
	case "GET":
		switch r.Header.Get("Accept") {
		case "text/xml":
			if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
				return res
			}
			h.Set("Content-Type", "text/xml")
			return quakeml12(r, h, b)
		default:
			if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
				return res
			}
			h.Set("Content-Type", "text/xml")
			return quakeml12(r, h, b)
		}
	default:
		return &weft.MethodNotAllowed
	}
}

func quakeml_rt1_2sHandler(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	switch r.Method {
	case "GET":
		switch r.Header.Get("Accept") {
		case "text/xml":
			if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
				return res
			}
			h.Set("Content-Type", "text/xml")
			return quakeml12RT(r, h, b)
		default:
			if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
				return res
			}
			h.Set("Content-Type", "text/xml")
			return quakeml12RT(r, h, b)
		}
	default:
		return &weft.MethodNotAllowed
	}
}
