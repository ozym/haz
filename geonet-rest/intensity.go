package main

import (
	"bytes"
	"net/http"
	"github.com/GeoNet/weft"
)

func intensityMeasuredLatestV1(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	if res := weft.CheckQuery(r, []string{"type"}, []string{}); !res.Ok {
		return res
	}

	if r.URL.Query().Get("type") != "measured" {
		return weft.BadRequest("type must be measured.")
	}

	var d string

	err := db.QueryRow(
		`SELECT row_to_json(fc)
				FROM ( SELECT 'FeatureCollection' as type, COALESCE(array_to_json(array_agg(f)), '[]') as features
					FROM (SELECT 'Feature' as type,
						ST_AsGeoJSON(s.location)::json as geometry,
						row_to_json(( select l from 
							( 
								select mmi
								) as l )) 
			as properties from (select location, mmi 
				FROM impact.intensity_measured) as s 
			) As f )  as fc`).Scan(&d)
	if err != nil {
		return weft.ServiceUnavailableError(err)
	}

	b.WriteString(d)
	h.Set("Content-Type", V1GeoJSON)
	return &weft.StatusOK
}
