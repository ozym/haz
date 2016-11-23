package main

import (
	"bytes"
	"database/sql"
	"github.com/GeoNet/haz"
	"github.com/GeoNet/weft"
	"github.com/golang/protobuf/proto"
	"net/http"
)

func valProto(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
		return res
	}

	var err error
	var rows *sql.Rows

	if rows, err = db.Query(`SELECT id, title, alert_level, activity, hazards,
				ST_X(location::geometry), ST_Y(location::geometry)
				FROM haz.volcano JOIN haz.volcanic_alert_level using (alert_level)
				ORDER BY alert_level DESC, title ASC`); err != nil {
		return weft.ServiceUnavailableError(err)
	}

	var vol haz.Volcanoes

	for rows.Next() {
		v := haz.Volcano{Val: &haz.VAL{}}
		if err = rows.Scan(&v.VolcanoID, &v.Title, &v.Val.Level, &v.Val.Activity, &v.Val.Hazards,
			&v.Longitude, &v.Latitude); err != nil {
			return weft.ServiceUnavailableError(err)
		}

		vol.Volcanoes = append(vol.Volcanoes, &v)
	}

	var by []byte

	if by, err = proto.Marshal(&vol); err != nil {
		return weft.ServiceUnavailableError(err)
	}

	b.Write(by)

	h.Set("Content-Type", protobuf)
	return &weft.StatusOK
}
