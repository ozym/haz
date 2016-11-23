package main

import (
	"bytes"
	"github.com/GeoNet/weft"
	"net/http"
	"strconv"
	"strings"
)

const (
	allMmi             = -1
	feltMmi            = 3
	defaultRecordCount = 30
	quakesNZServiceLen = 35 // len("/quakes/services/quakes/newzealand/")
)

func quakesWWWall(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
		return res
	}
	var d string
	err := db.QueryRow(quakesWWWSQL, allMmi, defaultRecordCount).Scan(&d)
	if err != nil {
		return weft.ServiceUnavailableError(err)
	}

	b.WriteString(d)
	h.Set("Content-Type", JSON)
	return &weft.StatusOK
}

func quakesWWWfelt(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
		return res
	}
	var d string
	err := db.QueryRow(quakesNZWWWSQL, feltMmi, defaultRecordCount).Scan(&d)
	if err != nil {
		return weft.ServiceUnavailableError(err)
	}

	b.WriteString(d)
	h.Set("Content-Type", JSON)
	return &weft.StatusOK
}

func quakesWWWnz(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
		return res
	}

	path := r.URL.Path[quakesNZServiceLen:] // ..."3/100.json"
	tokens := strings.Split(path, "/")
	if len(tokens) != 2 {
		return weft.BadRequest("Bad URL path.")
	}

	var mmi int
	var err error
	var count int
	if mmi, err = strconv.Atoi(tokens[0]); err != nil {
		return weft.BadRequest("Bad URL path. Invalid mmi.")
	}

	if count, err = strconv.Atoi(tokens[1][:len(tokens[1])-5]); err != nil { // len(".json")
		return weft.BadRequest("Bad URL path. Invalid count.")
	}

	var d string
	err = db.QueryRow(quakesNZWWWSQL, mmi, count).Scan(&d)
	if err != nil {
		return weft.ServiceUnavailableError(err)
	}

	b.WriteString(d)
	h.Set("Content-Type", JSON)
	return &weft.StatusOK
}

func quakeWWW(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
		return res
	}

	publicID := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/quake/services/quake/"), ".json")
	if publicID=="" {
		return weft.BadRequest("invalid publicID path: " + r.URL.Path)
	}


	var d string
	err := db.QueryRow(quakeWWWSQL, publicID).Scan(&d)
	if err != nil {
		return weft.ServiceUnavailableError(err)
	}

	b.WriteString(d)
	h.Set("Content-Type", JSON)
	return &weft.StatusOK
}
