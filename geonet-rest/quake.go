package main

import (
	"github.com/GeoNet/haz/msg"
	"github.com/GeoNet/web"
	"net/http"
)

func quakeV1(w http.ResponseWriter, r *http.Request) {
	if badQuery(w, r, []string{}, []string{}) {
		return
	}

	var publicID string
	var ok bool

	if publicID, ok = getPublicIDPath(w, r); !ok {
		return
	}

	var d string
	err := db.QueryRow(
		`SELECT row_to_json(fc)
                         FROM ( SELECT 'FeatureCollection' as type, array_to_json(array_agg(f)) as features
                         FROM (SELECT 'Feature' as type,
                         ST_AsGeoJSON(q.geom)::json as geometry,
                         row_to_json((SELECT l FROM 
                         	(
                         		SELECT 
                         		publicid AS "publicID",
                                to_char(time, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"') as "time",
                                depth, 
                                magnitude, 
                                locality,
                                intensity,
                                intensity_newzealand as "regionIntensity",
                                quality
                           ) as l
                         )) as properties FROM haz.quake as q where publicid = $1 ) As f )  as fc`, publicID).Scan(&d)
	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return
	}

	b := []byte(d)
	w.Header().Set("Content-Type", web.V1GeoJSON)
	web.Ok(w, r, &b)
}

func quakesRegionV1(w http.ResponseWriter, r *http.Request) {
	if badQuery(w, r, []string{"regionID", "regionIntensity", "number", "quality"}, []string{}) {
		return
	}

	var ok bool

	if _, ok = getRegionID(w, r); !ok {
		return
	}

	if _, ok = getQuality(w, r); !ok {
		return
	}

	var regionIntensity string

	if regionIntensity, ok = getRegionIntensity(w, r); !ok {
		return
	}

	var n int
	if n, ok = getNumberQuakes(w, r); !ok {
		return
	}

	var d string
	err := db.QueryRow(
		`SELECT row_to_json(fc)
                         FROM ( SELECT 'FeatureCollection' as type, COALESCE(array_to_json(array_agg(f)), '[]') as features
                         FROM (SELECT 'Feature' as type,
                         ST_AsGeoJSON(q.geom)::json as geometry,
                         row_to_json((SELECT l FROM
                         	(
                         		SELECT
                         		publicid AS "publicID",
                                to_char(time, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"') as "time",
                                depth,
                                magnitude,
                                locality,
                                intensity,
                                intensity_newzealand as "regionIntensity",
                                quality
                           ) as l
                         )) as properties FROM haz.quakeapi as q where mmid_newzealand >= $1
                         ORDER BY time DESC  limit $2 ) as f ) as fc`, msg.IntensityMMI(regionIntensity), n).Scan(&d)
	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return
	}

	b := []byte(d)
	w.Header().Set("Content-Type", web.V1GeoJSON)
	web.Ok(w, r, &b)
}
