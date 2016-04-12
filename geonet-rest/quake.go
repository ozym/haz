package main

import (
	"bytes"
	"github.com/GeoNet/haz/msg"
	"net/http"
)

func quakeV1(r *http.Request, h http.Header, b *bytes.Buffer) *result {
	if res := checkQuery(r, []string{}, []string{}); !res.ok {
		return res
	}

	var publicID string
	var res *result

	if publicID, res = getPublicIDPath(r); !res.ok {
		return res
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
		return serviceUnavailableError(err)
	}

	b.WriteString(d)
	h.Set("Content-Type", V1GeoJSON)
	return &statusOK
}

func quakesRegionV1(r *http.Request, h http.Header, b *bytes.Buffer) *result {
	if res := checkQuery(r, []string{"regionID", "regionIntensity", "number", "quality"}, []string{}); !res.ok {
		return res
	}

	var err error
	if _, err = getRegionID(r); err != nil {
		return badRequest(err.Error())
	}

	if _, err = getQuality(r); err != nil {
		return badRequest(err.Error())
	}

	var regionIntensity string

	if regionIntensity, err = getRegionIntensity(r); err != nil {
		return badRequest(err.Error())
	}

	var n int
	if n, err = getNumberQuakes(r); err != nil {
		return badRequest(err.Error())
	}

	var d string
	err = db.QueryRow(
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
                         ORDER BY time DESC  limit $2 ) as f ) as fc`, int(msg.IntensityMMI(regionIntensity)), n).Scan(&d)
	if err != nil {
		return serviceUnavailableError(err)
	}

	b.WriteString(d)
	h.Set("Content-Type", V1GeoJSON)
	return &statusOK
}
