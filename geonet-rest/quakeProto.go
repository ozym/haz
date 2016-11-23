package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/GeoNet/haz"
	"github.com/GeoNet/haz/sc3ml"
	"github.com/GeoNet/weft"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const s3 = "http://seiscompml07.s3-website-ap-southeast-2.amazonaws.com/"

func quakeStatsProto(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
		return res
	}

	var q haz.QuakeStats

	var rows *sql.Rows
	var err error

	if rows, err = db.Query(quakesPerDaySQL); err != nil {
		return weft.ServiceUnavailableError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var t time.Time
		var r haz.Rate

		if err = rows.Scan(&t, &r.Count); err != nil {
			return weft.ServiceUnavailableError(err)
		}

		r.Time = &haz.Timestamp{Sec: t.Unix(), Nsec: int64(t.Nanosecond())}

		q.PerDay = append(q.PerDay, &r)
	}
	rows.Close()

	q.Year = make(map[int32]int32)

	if rows, err = db.Query(fmt.Sprintf(sumMagsSQL, 365)); err != nil {
		return weft.ServiceUnavailableError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var k, v int32
		if err = rows.Scan(&k, &v); err != nil {
			return weft.ServiceUnavailableError(err)
		}
		q.Year[k] = v
	}
	rows.Close()

	q.Month = make(map[int32]int32)

	if rows, err = db.Query(fmt.Sprintf(sumMagsSQL, 28)); err != nil {
		return weft.ServiceUnavailableError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var k, v int32
		if err = rows.Scan(&k, &v); err != nil {
			return weft.ServiceUnavailableError(err)
		}
		q.Month[k] = v
	}
	rows.Close()

	q.Week = make(map[int32]int32)

	if rows, err = db.Query(fmt.Sprintf(sumMagsSQL, 7)); err != nil {
		return weft.ServiceUnavailableError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var k, v int32
		if err = rows.Scan(&k, &v); err != nil {
			return weft.ServiceUnavailableError(err)
		}
		q.Week[k] = v
	}
	rows.Close()

	var by []byte

	if by, err = proto.Marshal(&q); err != nil {
		return weft.ServiceUnavailableError(err)
	}

	b.Write(by)

	h.Set("Content-Type", protobuf)
	h.Set("Surrogate-Control", maxAge300)

	return &weft.StatusOK
}

func quakeProto(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
		return res
	}

	var q haz.Quake
	var res *weft.Result

	if q.PublicID, res = getPublicIDPath(r); !res.Ok {
		return res
	}

	var t time.Time
	var mt time.Time
	var err error

	if err = db.QueryRow(quakeProtoSQL, q.PublicID).Scan(&t, &mt,
		&q.Depth, &q.Magnitude, &q.Locality, &q.Mmi, &q.Quality,
		&q.Longitude, &q.Latitude); err != nil {
		return weft.ServiceUnavailableError(err)
	}

	q.Time = &haz.Timestamp{Sec: t.Unix(), Nsec: int64(t.Nanosecond())}
	q.ModificationTime = &haz.Timestamp{Sec: mt.Unix(), Nsec: int64(mt.Nanosecond())}

	var by []byte

	if by, err = proto.Marshal(&q); err != nil {
		return weft.ServiceUnavailableError(err)
	}

	b.Write(by)
	h.Set("Content-Type", protobuf)
	return &weft.StatusOK
}

func quakeHistoryProto(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
		return res
	}

	var publicID string
	var res *weft.Result

	if publicID, res = getPublicIDHistoryPath(r); !res.Ok {
		return res
	}

	var rows *sql.Rows
	var err error

	if rows, err = db.Query(quakeHistoryProtoSQL, publicID); err != nil {
		return weft.ServiceUnavailableError(err)
	}

	var quakes []*haz.Quake

	for rows.Next() {
		var t time.Time
		var mt time.Time
		q := haz.Quake{PublicID: publicID}

		if err = rows.Scan(&t, &mt, &q.Depth,
			&q.Magnitude, &q.Locality, &q.Mmi, &q.Quality,
			&q.Longitude, &q.Latitude); err != nil {
			return weft.ServiceUnavailableError(err)
		}

		q.Time = &haz.Timestamp{Sec: t.Unix(), Nsec: int64(t.Nanosecond())}
		q.ModificationTime = &haz.Timestamp{Sec: mt.Unix(), Nsec: int64(mt.Nanosecond())}

		quakes = append(quakes, &q)
	}

	qs := haz.Quakes{Quakes: quakes}

	var by []byte

	if by, err = proto.Marshal(&qs); err != nil {
		return weft.ServiceUnavailableError(err)
	}

	b.Write(by)
	h.Set("Content-Type", protobuf)
	return &weft.StatusOK
}

func quakesProto(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	if res := weft.CheckQuery(r, []string{"MMI"}, []string{}); !res.Ok {
		return res
	}

	var mmi int
	var err error

	if mmi, err = getMMI(r); err != nil {
		return weft.BadRequest(err.Error())
	}

	var rows *sql.Rows

	if rows, err = db.Query(quakesProtoSQL, mmi); err != nil {
		return weft.ServiceUnavailableError(err)
	}

	var quakes haz.Quakes

	for rows.Next() {
		var t time.Time
		var mt time.Time
		var q haz.Quake

		if err = rows.Scan(&q.PublicID, &t, &mt, &q.Depth,
			&q.Magnitude, &q.Locality, &q.Mmi, &q.Quality,
			&q.Longitude, &q.Latitude); err != nil {
			return weft.ServiceUnavailableError(err)
		}

		q.Time = &haz.Timestamp{Sec: t.Unix(), Nsec: int64(t.Nanosecond())}
		q.ModificationTime = &haz.Timestamp{Sec: mt.Unix(), Nsec: int64(mt.Nanosecond())}

		quakes.Quakes = append(quakes.Quakes, &q)
	}

	var by []byte

	if by, err = proto.Marshal(&quakes); err != nil {
		return weft.ServiceUnavailableError(err)
	}

	b.Write(by)
	h.Set("Content-Type", protobuf)
	return &weft.StatusOK
}

// fetches SC3ML and turns it into a protobuf.
func quakeTechnicalProto(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
		return res
	}

	by, res := getBytes(s3+strings.TrimPrefix(r.URL.Path, "/quake/technical/")+".xml", "")
	if !res.Ok {
		return res
	}

	q, err := sc3ml.QuakeTechnical(by)
	if err != nil {
		return weft.ServiceUnavailableError(err)
	}

	m, err := proto.Marshal(&q)
	if err != nil {
		return weft.ServiceUnavailableError(err)
	}

	b.Write(m)

	h.Set("Content-Type", protobuf)
	return &weft.StatusOK
}

/*
getBytes fetches bytes for the requested url.  accept
may be left as the empty string.
*/
func getBytes(url, accept string) ([]byte, *weft.Result) {
	var r *http.Response
	var req *http.Request
	var err error
	var b []byte

	if accept == "" {
		r, err = client.Get(url)
		if err != nil {
			return b, weft.InternalServerError(err)
		}
	} else {
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			return b, weft.InternalServerError(err)
		}

		req.Header.Add("Accept", accept)

		r, err = client.Do(req)
		if err != nil {
			return b, weft.InternalServerError(err)
		}
	}
	defer r.Body.Close()

	b, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return b, weft.InternalServerError(err)
	}

	switch r.StatusCode {
	case http.StatusOK:
		return b, &weft.StatusOK
	case http.StatusNotFound:
		return b, &weft.NotFound
	default:
		// TODO do we need to handle more errors here?
		return b, weft.InternalServerError(fmt.Errorf("server error"))

	}

	return b, &weft.StatusOK
}
