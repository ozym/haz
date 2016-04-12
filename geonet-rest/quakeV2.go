package main

import (
	"bytes"
	"net/http"
)

func quakeV2(r *http.Request, h http.Header, b *bytes.Buffer) *result {
	if res := checkQuery(r, []string{}, []string{}); !res.ok {
		return res
	}

	if len(r.URL.Query()) != 0 {
		return badRequest("incorrect number of query parameters.")
	}

	var publicID string
	var res *result

	if publicID, res = getPublicIDPath(r); !res.ok {
		return res
	}

	var d string
	err := db.QueryRow(quakeV2SQL, publicID).Scan(&d)
	if err != nil {
		return serviceUnavailableError(err)
	}

	b.WriteString(d)
	h.Set("Content-Type", V2GeoJSON)
	return &statusOK
}

func quakesV2(r *http.Request, h http.Header, b *bytes.Buffer) *result {
	if res := checkQuery(r, []string{"MMI"}, []string{}); !res.ok {
		return res
	}

	var mmi int
	var err error

	if mmi, err = getMMI(r); err != nil {
		return badRequest(err.Error())
	}

	var d string
	err = db.QueryRow(quakesV2SQL, mmi).Scan(&d)
	if err != nil {
		return serviceUnavailableError(err)
	}

	b.WriteString(d)
	h.Set("Content-Type", V2GeoJSON)
	return &statusOK
}

func quakeHistoryV2(r *http.Request, h http.Header, b *bytes.Buffer) *result {
	if res := checkQuery(r, []string{}, []string{}); !res.ok {
		return res
	}

	if len(r.URL.Query()) != 0 {
		return badRequest("incorrect number of query parameters.")
	}

	var publicID string
	var res *result

	if publicID, res = getPublicIDHistoryPath(r); !res.ok {
		return res
	}

	var d string
	err := db.QueryRow(quakeHistoryV2SQL, publicID).Scan(&d)
	if err != nil {
		return serviceUnavailableError(err)
	}

	b.WriteString(d)
	h.Set("Content-Type", V2GeoJSON)
	return &statusOK
}

func quakeStatsV2(r *http.Request, h http.Header, b *bytes.Buffer) *result {
	if res := checkQuery(r, []string{}, []string{}); !res.ok {
		return res
	}

	if len(r.URL.Query()) != 0 {
		return badRequest("incorrect number of query parameters.")
	}

	var d string
	err := db.QueryRow(quakeStatsV2SQL).Scan(&d)
	if err != nil {
		return serviceUnavailableError(err)
	}

	b.WriteString(d)
	h.Set("Surrogate-Control", maxAge300)
	h.Set("Content-Type", V2JSON)
	return &statusOK
}
