package main

import (
	"bytes"
	"github.com/GeoNet/haz"
	"github.com/GeoNet/weft"
	"github.com/golang/protobuf/proto"
	"net/http"
	"time"
)

func intensityProto(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	if res := weft.CheckQuery(r, []string{"type"}, []string{"publicID"}); !res.Ok {
		return res
	}

	var ts string
	var err error

	if ts, err = getIntensityType(r); err != nil {
		return weft.BadRequest(err.Error())
	}

	var shaking *haz.Shaking

	switch ts {
	case "measured":
		if shaking, err = intensityMeasuredLatest(); err != nil {
			return weft.ServiceUnavailableError(err)
		}
	case "reported":
		publicID := r.URL.Query().Get("publicID")
		switch publicID {
		case "":
			end := time.Now().UTC()
			start := end.Add(-60 * time.Minute)
			if shaking, err = intensityReported(start, end); err != nil {
				return weft.ServiceUnavailableError(err)
			}
		default:
			var t time.Time
			var res *weft.Result
			if t, res = getQuakeTime(r); !res.Ok {
				return res
			}
			start := t.Add(-1 * time.Minute)
			end := t.Add(15 * time.Minute)
			if shaking, err = intensityReported(start, end); err != nil {
				return weft.ServiceUnavailableError(err)
			}
		}
	}

	var by []byte
	if by, err = proto.Marshal(shaking); err != nil {
		return weft.ServiceUnavailableError(err)
	}

	b.Write(by)
	h.Set("Content-Type", protobuf)

	return &weft.StatusOK
}

func intensityMeasuredLatest() (*haz.Shaking, error) {
	// there is only one mmi at each point for measured so no need to handle count anywhere
	rows, err := db.Query(`SELECT ST_X(location::geometry), ST_Y(location::geometry), mmi
					FROM impact.intensity_measured`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var s haz.Shaking
	s.MmiSummary = make(map[int32]int32)

	for rows.Next() {
		var m haz.MMI
		if err = rows.Scan(&m.Longitude, &m.Latitude, &m.Mmi); err != nil {
			return nil, err
		}
		s.Mmi = append(s.Mmi, &m)
		s.MmiTotal++
		s.MmiSummary[m.Mmi]++
	}

	return &s, nil
}

func intensityReported(start, end time.Time) (*haz.Shaking, error) {

	// query the DB twice.  MMI is counted both times, the first time
	// in the geohash, the second time as totals.
	// This is a race. When there are lots of reports being submitted
	// the total counts could be different between the two queries.
	// Seems like an acceptable trade off between a more complex query
	// (creating JSON and then unmarshallaing it) or doing a lot of
	// aggregation in code (which requires iterating many rows outside the db).

	rows, err := db.Query(`SELECT
  				ST_X(st_pointfromgeohash(geohash6)),
  				ST_Y(st_pointfromgeohash(geohash6)),
  				max(mmi),
  				count(mmi)
			       FROM impact.intensity_reported
				WHERE time >= $1 AND time <= $2
				GROUP BY geohash6`, start, end)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var s haz.Shaking

	for rows.Next() {
		var m haz.MMI
		if err = rows.Scan(&m.Longitude, &m.Latitude, &m.Mmi, &m.Count); err != nil {
			return nil, err
		}
		s.Mmi = append(s.Mmi, &m)
	}

	rows.Close()

	rows, err = db.Query(`SELECT mmi, count(mmi) FROM impact.intensity_reported
				WHERE time >= $1 AND time <= $2
				GROUP BY mmi`, start, end)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	s.MmiSummary = make(map[int32]int32)

	for rows.Next() {
		var k, v int32
		if err = rows.Scan(&k, &v); err != nil {
			return nil, err
		}
		s.MmiSummary[k] = v
		s.MmiTotal += v
	}

	return &s, nil
}
