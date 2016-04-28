package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/GeoNet/weft"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

const minMMID float64 = 5.0

var capIDRe = regexp.MustCompile(`^[0-9a-z]+\.[0-9]+$`)
var serverCName = os.Getenv("WEB_SERVER_CNAME")

func capQuake(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
		return res
	}

	id := r.URL.Path[22:]

	if !capIDRe.MatchString(id) {
		return weft.BadRequest("invalid ID: " + id)
	}

	p := strings.Split(id, `.`)
	if len(p) != 2 {
		return weft.BadRequest("invalid ID: " + id)
	}

	c := capQuakeT{ID: id}
	c.Quake.PublicID = p[0]

	rows, err := db.Query(`select modificationTimeUnixMicro, modificationtime from haz.quakehistory
		where publicid = $1 AND modificationTimeUnixMicro < $2 AND status in ('reviewed','deleted')`, p[0], p[1])
	if err != nil {
		return weft.ServiceUnavailableError(err)
	}
	defer rows.Close()

	c.References = make([]string, 0)

	for rows.Next() {
		var i int
		var t time.Time
		err := rows.Scan(&i, &t)
		if err != nil {
			return weft.ServiceUnavailableError(err)
		}
		c.References = append(c.References, fmt.Sprintf("%s.%d,%s", c.Quake.PublicID, i, t.In(nz).Format(time.RFC3339)))
	}
	rows.Close()

	err = db.QueryRow(`select depth, 
		magnitude, 
		status, 
		usedPhaseCount,
		magnitudestationcount,
		longitude,
		latitude,
		time,
		modificationTime,
		intensity
	 FROM haz.quakehistory where publicid = $1 AND modificationTimeUnixMicro = $2`,
		p[0], p[1]).Scan(
		&c.Quake.Depth,
		&c.Quake.Magnitude,
		&c.Status,
		&c.Quake.UsedPhaseCount,
		&c.Quake.MagnitudeStationCount,
		&c.Quake.Longitude,
		&c.Quake.Latitude,
		&c.Quake.Time,
		&c.Quake.ModificationTime,
		&c.Intensity,
	)
	if err == sql.ErrNoRows {
		return &weft.NotFound
	}
	if err != nil {
		return weft.ServiceUnavailableError(err)
	}

	cl, err := c.Quake.Closest()
	if err != nil {
		return weft.ServiceUnavailableError(err)
	}

	c.Localities = c.Quake.Localities(minMMID)

	if len(c.Localities) == 0 {
		c.Localities = append(c.Localities, cl)
	}

	c.Closest = cl

	err = capTemplates.ExecuteTemplate(b, "capQuake", c)
	if err != nil {
		return weft.ServiceUnavailableError(err)
	}

	h.Set("Content-Type", CAP)
	return &weft.StatusOK
}

func capQuakeFeed(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
		return res
	}

	// we are only serving /cap/1.2/GPAv1.0/feed/atom1.0/quake at the moment and the router
	// matches that so no need for any further validation here yet.

	atom := capAtomFeed{
		Title: `CAP quakes`,
		ID:    fmt.Sprintf("https://%s/cap/1.2/GPA1.0/feed/atom1.0/quake", serverCName),
		Link:  fmt.Sprintf("https://%s/cap/1.2/GPA1.0/feed/atom1.0/quake", serverCName),
	}

	rows, err := db.Query(capQuakeFeedSQL, int(minMMID))
	if err != nil {
		return weft.ServiceUnavailableError(err)
	}
	defer rows.Close()

	tLatest := time.Time{}
	for rows.Next() {

		var p string
		var i int
		t := time.Time{}

		err := rows.Scan(&p, &i, &t)
		if err != nil {
			return weft.ServiceUnavailableError(err)
		}

		entry := capAtomEntry{
			ID:       fmt.Sprintf("http://geonet.org/nz/quakes/%s.%d", p, i),
			Title:    fmt.Sprintf("Quake CAP Message %s.%d", p, i),
			Updated:  t,
			Summary:  fmt.Sprintf("Quake CAP Message %s.%d", p, i),
			HrefCAP:  fmt.Sprintf("https://%s/cap/1.2/GPA1.0/quake/%s.%d", serverCName, p, i),
			HrefHTML: fmt.Sprintf("http://geonet.org.nz/quakes/%s", p),
		}

		atom.Entries = append(atom.Entries, entry)

		if t.After(tLatest) {
			tLatest = t
		}
	}
	rows.Close()

	if tLatest.Equal(time.Time{}) {
		tLatest = time.Now().UTC()
	}

	atom.Updated = tLatest
	err = capTemplates.ExecuteTemplate(b, "capAtom", atom)
	if err != nil {
		return weft.ServiceUnavailableError(err)
	}

	h.Set("Content-Type", Atom)
	return &weft.StatusOK
}
