package main

import (
	"github.com/GeoNet/web"
	"net/http"
)

func quakeV2(w http.ResponseWriter, r *http.Request) {
	if badQuery(w, r, []string{}, []string{}) {
		return
	}

	if len(r.URL.Query()) != 0 {
		web.BadRequest(w, r, "incorrect number of query parameters.")
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
                                floor(mmid_newzealand) as "mmi",
                                quality
                           ) as l
                         )) as properties FROM haz.quake as q where publicid = $1 ) As f )  as fc`, publicID).Scan(&d)
	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return
	}

	b := []byte(d)
	w.Header().Set("Content-Type", web.V2GeoJSON)
	web.Ok(w, r, &b)
}

func quakesV2(w http.ResponseWriter, r *http.Request) {
	if badQuery(w, r, []string{"MMI"}, []string{}) {
		return
	}

	var mmi int
	var ok bool

	if mmi, ok = getMMI(w, r); !ok {
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
                                floor(mmid_newzealand) as "mmi",
                                quality
                           ) as l
                         )) as properties FROM haz.quakeapi as q where mmid_newzealand >= $1
		AND In_newzealand = true
                         ORDER BY time DESC  limit 100 ) as f ) as fc`, mmi).Scan(&d)
	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return
	}

	b := []byte(d)
	w.Header().Set("Content-Type", web.V2GeoJSON)
	web.Ok(w, r, &b)
}

func quakeHistoryV2(w http.ResponseWriter, r *http.Request) {
	if badQuery(w, r, []string{}, []string{}) {
		return
	}

	if len(r.URL.Query()) != 0 {
		web.BadRequest(w, r, "incorrect number of query parameters.")
		return
	}

	var publicID string
	var ok bool

	if publicID, ok = getPublicIDHistoryPath(w, r); !ok {
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
                                to_char(modificationtime, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"') as "modificationTime",
                                depth, 
                                magnitude, 
                                locality,
                                floor(mmid_newzealand) as "mmi",
                                quality
                           ) as l
                         )) as properties FROM haz.quakehistory as q where publicid = $1 order by modificationtime desc ) As f )  as fc`, publicID).Scan(&d)
	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return
	}

	b := []byte(d)
	w.Header().Set("Content-Type", web.V2GeoJSON)
	web.Ok(w, r, &b)
}
