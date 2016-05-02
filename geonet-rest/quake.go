package main

import (
	"bytes"
	"github.com/GeoNet/haz/msg"
	"net/http"
	"github.com/GeoNet/weft"
)

func quakeV1(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
		return res
	}

	var publicID string
	var res *weft.Result

	if publicID, res = getPublicIDPath(r); !res.Ok {
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
		return weft.ServiceUnavailableError(err)
	}

	b.WriteString(d)
	h.Set("Content-Type", V1GeoJSON)
	return &weft.StatusOK
}

func quakesRegionV1(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	if res := weft.CheckQuery(r, []string{"regionID", "regionIntensity", "number", "quality"}, []string{}); !res.Ok {
		return res
	}

	var err error
	if _, err = getRegionID(r); err != nil {
		return weft.BadRequest(err.Error())
	}

	if _, err = getQuality(r); err != nil {
		return weft.BadRequest(err.Error())
	}

	var regionIntensity string

	if regionIntensity, err = getRegionIntensity(r); err != nil {
		return weft.BadRequest(err.Error())
	}

	var n int
	if n, err = getNumberQuakes(r); err != nil {
		return weft.BadRequest(err.Error())
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
		return weft.ServiceUnavailableError(err)
	}

	b.WriteString(d)
	h.Set("Content-Type", V1GeoJSON)
	return &weft.StatusOK
}
