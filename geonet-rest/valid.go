package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const quakeLen = 7         //  len("/quake/")
const quakeHistoryLen = 15 //  len("/quake/history/")

var publicIDRe = regexp.MustCompile(`^[0-9a-z]+$`)
var intensityRe = regexp.MustCompile(`^(unnoticeable|weak|light|moderate|strong|severe)$`)
var qualityRe = regexp.MustCompile(`^(best|caution|deleted|good)$`)

func getPublicIDPath(r *http.Request) (string, *result) {
	publicID := r.URL.Path[quakeLen:]

	if !publicIDRe.MatchString(publicID) {
		return publicID, badRequest("invalid publicID: " + publicID)
	}

	var d string
	err := db.QueryRow("select publicid FROM haz.quake where publicid = $1", publicID).Scan(&d)
	if err == sql.ErrNoRows {
		return publicID, &notFound
	}
	if err != nil {
		return publicID, serviceUnavailableError(err)
	}

	return publicID, &statusOK
}

func getPublicIDHistoryPath(r *http.Request) (string, *result) {
	publicID := r.URL.Path[quakeHistoryLen:]

	if !publicIDRe.MatchString(publicID) {
		return publicID, badRequest("invalid publicID: " + publicID)
	}

	var d string

	err := db.QueryRow("select publicid FROM haz.quake where publicid = $1", publicID).Scan(&d)
	if err == sql.ErrNoRows {
		return publicID, &notFound
	}
	if err != nil {
		return publicID, serviceUnavailableError(err)
	}

	return publicID, &statusOK
}

func getPublicID(r *http.Request) (string, *result) {
	publicID := r.URL.Query().Get("publicID")

	if !publicIDRe.MatchString(publicID) {
		return publicID, badRequest(fmt.Sprintf("invalid publicID " + publicID))
	}

	var d string

	err := db.QueryRow("select publicid FROM haz.quake where publicid = $1", publicID).Scan(&d)

	if err == sql.ErrNoRows {
		return publicID, &notFound
	}
	if err != nil {
		return publicID, serviceUnavailableError(err)
	}

	return publicID, &statusOK
}

func getMMI(r *http.Request) (int, error) {
	mmi, err := strconv.Atoi(r.URL.Query().Get("MMI"))
	if err != nil || mmi > 8 || mmi < -1 {
		return 0, fmt.Errorf("Invalid MMI query param.")
	}

	if mmi <= 2 {
		mmi = -9
	}

	return mmi, nil
}

func getIntensityType(r *http.Request) (string, error) {
	t := r.URL.Query().Get("type")
	switch t {
	case `measured`, `reported`:
		return t, nil
	default:
		return ``, fmt.Errorf("invalid intensity type")
	}
}

func getQuakeTime(r *http.Request) (time.Time, *result) {
	publicID := r.URL.Query().Get("publicID")
	originTime := time.Time{}

	if !publicIDRe.MatchString(publicID) {
		return originTime, badRequest(fmt.Sprintf("invalid publicID " + publicID))
	}

	err := db.QueryRow("select time FROM haz.quake where publicid = $1", publicID).Scan(&originTime)
	if err == sql.ErrNoRows {
		return originTime, &notFound
	}
	if err != nil {
		return originTime, serviceUnavailableError(err)
	}

	return originTime, &statusOK
}

func getRegionID(r *http.Request) (string, error) {
	regionID := r.URL.Query().Get("regionID")

	if regionID != "newzealand" {
		return regionID, fmt.Errorf("Invalid query parameter regionID: " + regionID)
	}

	return regionID, nil
}

func getQuality(r *http.Request) ([]string, error) {
	quality := strings.Split(r.URL.Query().Get("quality"), ",")

	for _, q := range quality {
		if !qualityRe.MatchString(q) {
			return quality, fmt.Errorf("Invalid quality: " + q)
		}
	}

	return quality, nil
}

func getRegionIntensity(r *http.Request) (string, error) {
	regionIntensity := r.URL.Query().Get("regionIntensity")

	if !intensityRe.MatchString(regionIntensity) {
		return regionIntensity, fmt.Errorf("Invalid regionIntensity: " + regionIntensity)
	}

	return regionIntensity, nil
}

func getNumberQuakes(r *http.Request) (int, error) {
	n, err := strconv.Atoi(r.URL.Query().Get("number"))
	if err != nil {
		return n, fmt.Errorf("Invalid query parameter number")
	}

	switch n {
	case 3, 30, 100, 500, 1000, 1500:
		return n, nil
	default:
		return n, fmt.Errorf("Invalid query parameter number")
	}
}
