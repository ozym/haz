package main

import (
	"bytes"
	"github.com/GeoNet/weft"
	"net/http"
)

func quakeV2(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
		return res
	}

	if len(r.URL.Query()) != 0 {
		return weft.BadRequest("incorrect number of query parameters.")
	}

	var publicID string
	var res *weft.Result

	if publicID, res = getPublicIDPath(r); !res.Ok {
		return res
	}

	var d string
	err := db.QueryRow(quakeV2SQL, publicID).Scan(&d)
	if err != nil {
		return weft.ServiceUnavailableError(err)
	}

	b.WriteString(d)
	h.Set("Content-Type", V2GeoJSON)
	return &weft.StatusOK
}

func quakesV2(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	if res := weft.CheckQuery(r, []string{"MMI"}, []string{}); !res.Ok {
		return res
	}

	var mmi int
	var err error

	if mmi, err = getMMI(r); err != nil {
		return weft.BadRequest(err.Error())
	}

	var d string
	err = db.QueryRow(quakesV2SQL, mmi).Scan(&d)
	if err != nil {
		return weft.ServiceUnavailableError(err)
	}

	b.WriteString(d)
	h.Set("Content-Type", V2GeoJSON)
	return &weft.StatusOK
}

func quakeHistoryV2(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
		return res
	}

	if len(r.URL.Query()) != 0 {
		return weft.BadRequest("incorrect number of query parameters.")
	}

	var publicID string
	var res *weft.Result

	if publicID, res = getPublicIDHistoryPath(r); !res.Ok {
		return res
	}

	var d string
	err := db.QueryRow(quakeHistoryV2SQL, publicID).Scan(&d)
	if err != nil {
		return weft.ServiceUnavailableError(err)
	}

	b.WriteString(d)
	h.Set("Content-Type", V2GeoJSON)
	return &weft.StatusOK
}

func quakeStatsV2(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
		return res
	}

	if len(r.URL.Query()) != 0 {
		return weft.BadRequest("incorrect number of query parameters.")
	}

	var d string
	err := db.QueryRow(quakeStatsV2SQL).Scan(&d)
	if err != nil {
		return weft.ServiceUnavailableError(err)
	}

	b.WriteString(d)
	h.Set("Surrogate-Control", maxAge300)
	h.Set("Content-Type", V2JSON)
	return &weft.StatusOK
}
