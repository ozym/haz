package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/GeoNet/web"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const minMMID float64 = 5.0

var capIDRe = regexp.MustCompile(`^[0-9a-z]+\.[0-9]+$`)

func capQuake(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query()) != 0 {
		web.BadRequest(w, r, "incorrect number of query parameters.")
		return
	}

	id := r.URL.Path[22:]

	if !capIDRe.MatchString(id) {
		web.BadRequest(w, r, "invalid ID: "+id)
		return
	}

	p := strings.Split(id, `.`)
	if len(p) != 2 {
		web.BadRequest(w, r, "invalid ID: "+id)
		return
	}

	c := capQuakeT{ID: id}
	c.Quake.PublicID = p[0]

	rows, err := db.Query(`select modificationTimeUnixMicro, modificationtime from haz.quakehistory
		where publicid = $1 AND modificationTimeUnixMicro < $2 AND status in ('reviewed','deleted')`, p[0], p[1])
	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return
	}
	defer rows.Close()

	c.References = make([]string, 0)

	for rows.Next() {
		var i int
		var t time.Time
		err := rows.Scan(&i, &t)
		if err != nil {
			web.ServiceUnavailable(w, r, err)
			return
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
		web.NotFound(w, r, "invalid ID: "+id)
		return
	}
	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return
	}

	cl, err := c.Quake.Closest()
	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return
	}

	c.Localities = c.Quake.Localities(minMMID)

	if len(c.Localities) == 0 {
		c.Localities = append(c.Localities, cl)
	}

	c.Closest = cl

	b := new(bytes.Buffer)

	err = capTemplates.ExecuteTemplate(b, "capQuake", c)
	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return
	}

	w.Header().Set("Content-Type", web.CAP)
	web.OkBuf(w, r, b)
}

func capQuakeFeed(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query()) != 0 {
		web.BadRequest(w, r, "incorrect number of query parameters.")
		return
	}

	// we are only serving /cap/1.2/GPAv1.0/feed/atom1.0/quake at the moment and the router
	// matches that so no need for any further validation here yet.

	atom := capAtomFeed{
		Title: `CAP quakes`,
		ID:    fmt.Sprintf("https://%s/cap/1.2/GPA1.0/feed/atom1.0/quake", config.WebServer.CNAME),
		Link:  fmt.Sprintf("https://%s/cap/1.2/GPA1.0/feed/atom1.0/quake", config.WebServer.CNAME),
	}

	/*
		There needs to be an Atom feed entity for every CAP message.
		A CAP message ID is unique for each CAP message (and is not just the quake PublicID).
		1. Find the first time any quake was reviewed within an hour of the quake
		and was strong enough for a CAP message.  This is the first CAP message.
		2. Select the first CAP message and any subsequent reviews or deletes that happened with an hour
		of the quake.  Each of this is a CAP message and gets an entity in the feed.
	*/
	rows, err := db.Query(`with first_review as 
	(select publicid, min(modificationtimeunixmicro) as modificationtimeunixmicro 
		from haz.quakehistory 
		where status = 'reviewed' 
		AND modificationTime - time < interval '1 hour' 
		AND MMID_newzealand >= $1 
		AND now() - time < interval '48 hours' group by publicid)
	select h.publicid, h.modificationtimeunixmicro, h.modificationTime 
	from haz.quakehistory h, first_review 
	where h.publicid = first_review.publicid 
	and h.modificationtimeunixmicro >= first_review.modificationtimeunixmicro 
	and status in ('reviewed','deleted') 
	AND modificationTime - time < interval '1 hour' ORDER BY time DESC, modificationTime DESC`, int(minMMID))
	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return
	}
	defer rows.Close()

	tLatest := time.Time{}
	for rows.Next() {

		var p string
		var i int
		t := time.Time{}

		err := rows.Scan(&p, &i, &t)
		if err != nil {
			web.ServiceUnavailable(w, r, err)
			return
		}

		entry := capAtomEntry{
			ID:       fmt.Sprintf("http://geonet.org/nz/quakes/%s.%d", p, i),
			Title:    fmt.Sprintf("Quake CAP Message %s.%d", p, i),
			Updated:  t,
			Summary:  fmt.Sprintf("Quake CAP Message %s.%d", p, i),
			HrefCAP:  fmt.Sprintf("https://%s/cap/1.2/GPA1.0/quake/%s.%d", config.WebServer.CNAME, p, i),
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

	b := new(bytes.Buffer)

	err = capTemplates.ExecuteTemplate(b, "capAtom", atom)
	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return
	}

	w.Header().Set("Content-Type", web.Atom)
	web.OkBuf(w, r, b)
}
