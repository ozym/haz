package main

import (
	"bytes"
	"net/http"
	"time"
)

func intensityV2(r *http.Request, h http.Header, b *bytes.Buffer) *result {
	if res := checkQuery(r, []string{"type"}, []string{"publicID"}); !res.ok {
		return res
	}

	var ts string
	var err error

	if ts, err = getIntensityType(r); err != nil {
		return badRequest(err.Error())
	}

	var d string

	switch ts {
	case "measured":
		err = db.QueryRow(intensityMeasuredLatestV2SQL).Scan(&d)
	case "reported":
		publicID := r.URL.Query().Get("publicID")
		switch publicID {
		case "":
			err = db.QueryRow(intenstityReportedLatestV2SQL).Scan(&d)
		default:
			var t time.Time
			var res *result
			if t, res = getQuakeTime(r); !res.ok {
				return res
			}
			err = db.QueryRow(intenstityReportedWindowV2SQL, t.Add(time.Duration(-1*time.Minute)), t.Add(time.Duration(15*time.Minute))).Scan(&d)
		}
	}

	if err != nil {
		return serviceUnavailableError(err)
	}

	b.WriteString(d)
	h.Set("Content-Type", V2GeoJSON)

	return &statusOK
}
