package main

import (
	"bytes"
	"net/http"
	"strconv"
	"time"
	"github.com/GeoNet/weft"
	"log"
)

const head = `<html xmlns="http://www.w3.org/1999/xhtml"><head><title>GeoNet - SOH</title><style type="text/css">
table {border-collapse: collapse; margin: 0px; padding: 2px;}
table th {background-color: black; color: white;}
table td {border: 1px solid silver; margin: 0px;}
table tr {background-color: #99ff99;}
table tr.error {background-color: #FF0000;}
</style></head><h2>State of Health</h2>`
const foot = "</body></html>"

var (
	old time.Duration
)

func init() {
	old = time.Duration(-1) * time.Minute
}

// returns a simple state of health page.  If heartbeat times in the
// DB are old then it also returns an http status of 500.
// Not useful for inclusion in app metrics so weft not used.
func sohEsb(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", HtmlContent)

	if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
		w.Header().Set("Surrogate-Control", "max-age=86400")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var b bytes.Buffer

	b.Write([]byte(head))
	b.Write([]byte(`<p>Current time is: ` + time.Now().UTC().String() + `</p>`))
	b.Write([]byte(`<h3>Messaging</h3>`))

	var bad bool
	var s string
	var t time.Time

	b.Write([]byte(`<table><tr><th>Service</th><th>Time Received</th></tr>`))

	rows, err := db.Query("select serverid, timereceived from haz.soh")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&s, &t)
			if err == nil {
				if t.Before(time.Now().UTC().Add(old)) {
					bad = true
					b.Write([]byte(`<tr class="tr error">`))
				} else {
					b.Write([]byte(`<tr>`))
				}
				b.Write([]byte(`<td>` + s + `</td><td>` + t.String() + `</td></tr>`))
			} else {
				bad = true
				b.Write([]byte(`<tr class="tr error"><td>DB error</td><td>` + err.Error() + `</td></tr>`))
			}
		}
		rows.Close()
	} else {
		log.Printf("ERROR: %v", err)
		bad = true
		b.Write([]byte(`<tr class="tr error"><td>DB error</td><td>` + err.Error() + `</td></tr>`))
	}
	b.Write([]byte(`</table>`))

	b.Write([]byte(foot))

	if bad {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	b.WriteTo(w)
}

// returns a simple state of health page.  If the count of measured intensities falls below 50 this it also returns an http status of 500.
//func impactSOH(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
// Not useful for inclusion in app metrics so weft not used.
func impactSOH(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", HtmlContent)

	if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
		w.Header().Set("Surrogate-Control", "max-age=86400")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var meas int
	err := db.QueryRow("select count(*) from impact.intensity_measured").Scan(&meas)
	if err != nil  {
		w.WriteHeader(http.StatusServiceUnavailable)
		log.Printf("ERROR: %v", err)
	}

	if meas < 50 {
		w.WriteHeader(http.StatusServiceUnavailable)
		log.Printf("ERROR: less than 50 shaking stations %d", meas)
	}

	w.Write([]byte(head))
	w.Write([]byte(`<p>Current time is: ` + time.Now().UTC().String() + `</p>`))
	w.Write([]byte(`<h3>Impact</h3>`))

	w.Write([]byte(`<table><tr><th>Impact</th><th>Count</th></tr>`))


	if err == nil {
		if meas < 50 {
			w.Write([]byte(`<tr class="tr error"><td>shaking measured</td><td>` + strconv.Itoa(meas) + ` < 50</td></tr>`))
		} else {
			w.Write([]byte(`<tr><td>shaking measured</td><td>` + strconv.Itoa(meas) + ` >= 50</td></tr>`))
		}
	} else {
		w.Write([]byte(`<tr class="tr error"><td>DB error</td><td>` + err.Error() + `</td></tr>`))
	}
	w.Write([]byte(`</table>`))

	w.Write([]byte(foot))
}

// up is for testing that the app has started e.g., for with load balancers.
// It indicates the app is started.  It may still be serving errors.
// Not useful for inclusion in app metrics so weft not used.
func up(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
		w.Header().Set("Surrogate-Control", "max-age=86400")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write([]byte("<html><head></head><body>up</body></html>"))
}

// soh is for external service probes.
// writes a service unavailable error to w if the service is not working.
// Not useful for inclusion in app metrics so weft not used.
func soh(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
		w.Header().Set("Surrogate-Control", "max-age=86400")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var c int

	if err := db.QueryRow("SELECT 1").Scan(&c); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("<html><head></head><body>service error</body></html>"))
		log.Printf("ERROR: soh service error %s", err)
		return
	}

	w.Write([]byte("<html><head></head><body>ok</body></html>"))
}