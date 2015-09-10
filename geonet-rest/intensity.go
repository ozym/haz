package main

import (
	"github.com/GeoNet/web"
	"net/http"
	"regexp"
)

var zoomRe = regexp.MustCompile(`^(5|6)$`)

func intensityMeasuredLatestV1(w http.ResponseWriter, r *http.Request) {
	if badQuery(w, r, []string{"type"}, []string{}) {
		return
	}

	if r.URL.Query().Get("type") != "measured" {
		web.BadRequest(w, r, "type must be measured.")
		return
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
		web.ServiceUnavailable(w, r, err)
		return
	}

	b := []byte(d)
	w.Header().Set("Content-Type", web.V1GeoJSON)
	web.Ok(w, r, &b)
}
