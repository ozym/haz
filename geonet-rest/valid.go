package main

import (
	"database/sql"
	"github.com/GeoNet/web"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const quakeLen = 7 //  len("/quake/")

var publicIDRe = regexp.MustCompile(`^[0-9a-z]+$`)
var intensityRe = regexp.MustCompile(`^(unnoticeable|weak|light|moderate|strong|severe)$`)
var qualityRe = regexp.MustCompile(`^(best|caution|deleted|good)$`)

func getPublicIDPath(w http.ResponseWriter, r *http.Request) (string, bool) {
	publicID := r.URL.Path[quakeLen:]

	if !publicIDRe.MatchString(publicID) {
		web.BadRequest(w, r, "invalid publicID: "+publicID)
		return publicID, false
	}

	var d string

	err := db.QueryRow("select publicid FROM haz.quake where publicid = $1", publicID).Scan(&d)
	if err == sql.ErrNoRows {
		web.NotFound(w, r, "invalid publicID: "+publicID)
		return publicID, false
	}
	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return publicID, false
	}

	return publicID, true
}

func getPublicID(w http.ResponseWriter, r *http.Request) (string, bool) {
	publicID := r.URL.Query().Get("publicID")

	if !publicIDRe.MatchString(publicID) {
		web.BadRequest(w, r, "invalid publicID: "+publicID)
		return publicID, false
	}

	var d string

	err := db.QueryRow("select publicid FROM haz.quake where publicid = $1", publicID).Scan(&d)
	if err == sql.ErrNoRows {
		web.NotFound(w, r, "invalid publicID: "+publicID)
		return publicID, false
	}
	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return publicID, false
	}

	return publicID, true
}

func getMMI(w http.ResponseWriter, r *http.Request) (int, bool) {
	mmi, err := strconv.Atoi(r.URL.Query().Get("MMI"))
	if err != nil || mmi > 8 || mmi < -1 {
		web.BadRequest(w, r, "Invalid MMI query param.")
		return 0, false
	}

	if mmi <= 2 {
		mmi = -9
	}

	return mmi, true
}

func getIntensityType(w http.ResponseWriter, r *http.Request) (string, bool) {
	t := r.URL.Query().Get("type")
	switch t {
	case `measured`, `reported`:
		return t, true
	default:
		web.BadRequest(w, r, "invalid intensity type")
		return ``, false
	}
}

func getQuakeTime(w http.ResponseWriter, r *http.Request) (time.Time, bool) {
	publicID := r.URL.Query().Get("publicID")
	originTime := time.Time{}

	if !publicIDRe.MatchString(publicID) {
		web.BadRequest(w, r, "invalid publicID: "+publicID)
		return originTime, false
	}

	var err error

	err = db.QueryRow("select time FROM haz.quake where publicid = $1", publicID).Scan(&originTime)
	if err == sql.ErrNoRows {
		web.NotFound(w, r, "invalid publicID: "+publicID)
		return originTime, false
	}
	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return originTime, false
	}

	return originTime, true
}

func getRegionID(w http.ResponseWriter, r *http.Request) (string, bool) {
	regionID := r.URL.Query().Get("regionID")

	if regionID != "newzealand" {
		web.BadRequest(w, r, "Invalid query parameter regionID: "+regionID)
		return regionID, false
	}

	return regionID, true
}

func getQuality(w http.ResponseWriter, r *http.Request) ([]string, bool) {
	quality := strings.Split(r.URL.Query().Get("quality"), ",")

	for _, q := range quality {
		if !qualityRe.MatchString(q) {
			web.BadRequest(w, r, "Invalid quality: "+q)
			return quality, false
		}
	}

	return quality, true
}

func getRegionIntensity(w http.ResponseWriter, r *http.Request) (string, bool) {
	regionIntensity := r.URL.Query().Get("regionIntensity")

	if !intensityRe.MatchString(regionIntensity) {
		web.BadRequest(w, r, "Invalid regionIntensity: "+regionIntensity)
		return regionIntensity, false
	}

	return regionIntensity, true
}

func getNumberQuakes(w http.ResponseWriter, r *http.Request) (int, bool) {
	n, err := strconv.Atoi(r.URL.Query().Get("number"))
	if err != nil {
		web.BadRequest(w, r, "Invalid query parameter number")
		return n, false
	}

	switch n {
	case 3, 30, 100, 500, 1000, 1500:
		return n, true
	default:
		web.BadRequest(w, r, "Invalid query parameter number")
		return n, false
	}
}
