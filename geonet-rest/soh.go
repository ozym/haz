package main

import (
	"bytes"
	"github.com/GeoNet/web"
	"net/http"
	"strconv"
	"time"
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

// returns a simple state of health page.  If heartbeat times in the DB are old then it also returns an http status of 500.
func soh(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", web.HtmlContent)
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
		bad = true
		b.Write([]byte(`<tr class="tr error"><td>DB error</td><td>` + err.Error() + `</td></tr>`))
	}
	b.Write([]byte(`</table>`))

	b.Write([]byte(foot))

	if bad {
		web.ServiceInternalServerErrorBuf(w, r, &b)
		return
	}

	web.OkBuf(w, r, &b)
}

// returns a simple state of health page.  If the count of measured intensities falls below 50 this it also returns an http status of 500.
func impactSOH(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", web.HtmlContent)
	var b bytes.Buffer

	b.Write([]byte(head))
	b.Write([]byte(`<p>Current time is: ` + time.Now().UTC().String() + `</p>`))
	b.Write([]byte(`<h3>Impact</h3>`))

	var bad bool

	b.Write([]byte(`<table><tr><th>Impact</th><th>Count</th></tr>`))

	var meas int
	err := db.QueryRow("select count(*) from impact.intensity_measured").Scan(&meas)
	if err == nil {
		if meas < 50 {
			bad = true
			b.Write([]byte(`<tr class="tr error"><td>shaking measured</td><td>` + strconv.Itoa(meas) + ` < 50</td></tr>`))
		} else {
			b.Write([]byte(`<tr><td>shaking measured</td><td>` + strconv.Itoa(meas) + ` >= 50</td></tr>`))
		}
	} else {
		bad = true
		b.Write([]byte(`<tr class="tr error"><td>DB error</td><td>` + err.Error() + `</td></tr>`))
	}
	b.Write([]byte(`</table>`))

	b.Write([]byte(foot))

	if bad {
		web.ServiceInternalServerErrorBuf(w, r, &b)
		return
	}

	web.OkBuf(w, r, &b)
}
