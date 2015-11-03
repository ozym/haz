package main

import (
	"github.com/GeoNet/web"
	"net/http"
)

func quakeV2(w http.ResponseWriter, r *http.Request) {
	if badQuery(w, r, []string{}, []string{}) {
		return
	}

	if len(r.URL.Query()) != 0 {
		web.BadRequest(w, r, "incorrect number of query parameters.")
		return
	}

	var publicID string
	var ok bool

	if publicID, ok = getPublicIDPath(w, r); !ok {
		return
	}

	var d string
	err := db.QueryRow(quakeV2SQL, publicID).Scan(&d)
	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return
	}

	b := []byte(d)
	w.Header().Set("Content-Type", web.V2GeoJSON)
	web.Ok(w, r, &b)
}

func quakesV2(w http.ResponseWriter, r *http.Request) {
	if badQuery(w, r, []string{"MMI"}, []string{}) {
		return
	}

	var mmi int
	var ok bool

	if mmi, ok = getMMI(w, r); !ok {
		return
	}

	var d string
	err := db.QueryRow(quakesV2SQL, mmi).Scan(&d)
	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return
	}

	b := []byte(d)
	w.Header().Set("Content-Type", web.V2GeoJSON)
	web.Ok(w, r, &b)
}

func quakeHistoryV2(w http.ResponseWriter, r *http.Request) {
	if badQuery(w, r, []string{}, []string{}) {
		return
	}

	if len(r.URL.Query()) != 0 {
		web.BadRequest(w, r, "incorrect number of query parameters.")
		return
	}

	var publicID string
	var ok bool

	if publicID, ok = getPublicIDHistoryPath(w, r); !ok {
		return
	}

	var d string
	err := db.QueryRow(quakeHistoryV2SQL, publicID).Scan(&d)
	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return
	}

	b := []byte(d)
	w.Header().Set("Content-Type", web.V2GeoJSON)
	web.Ok(w, r, &b)
}

func quakeStatsV2(w http.ResponseWriter, r *http.Request) {
	if badQuery(w, r, []string{}, []string{}) {
		return
	}

	if len(r.URL.Query()) != 0 {
		web.BadRequest(w, r, "incorrect number of query parameters.")
		return
	}

	var d string
	err := db.QueryRow(quakeStatsV2SQL).Scan(&d)
	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return
	}

	b := []byte(d)
	w.Header().Set("Surrogate-Control", web.MaxAge300)
	w.Header().Set("Content-Type", web.V2JSON)
	web.Ok(w, r, &b)
}
