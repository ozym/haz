package main

import (
	"github.com/GeoNet/web"
	"net/http"
	"time"
)

func intensityV2(w http.ResponseWriter, r *http.Request) {
	if badQuery(w, r, []string{"type"}, []string{"publicID"}) {
		return
	}

	var t string
	var ok bool

	if t, ok = getIntensityType(w, r); !ok {
		return
	}

	var d string
	var err error

	switch t {
	case "measured":
		err = db.QueryRow(intensityMeasuredLatestV2SQL).Scan(&d)
	case "reported":
		switch r.URL.Query().Get("publicID") {
		case "":
			err = db.QueryRow(intenstityReportedLatestV2SQL).Scan(&d)
		default:
			var t time.Time
			var ok bool
			if t, ok = getQuakeTime(w, r); !ok {
				return
			}
			err = db.QueryRow(intenstityReportedWindowV2SQL, t.Add(time.Duration(-1*time.Minute)), t.Add(time.Duration(15*time.Minute))).Scan(&d)
		}
	}

	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return
	}

	b := []byte(d)
	w.Header().Set("Content-Type", web.V2GeoJSON)
	web.Ok(w, r, &b)
}
