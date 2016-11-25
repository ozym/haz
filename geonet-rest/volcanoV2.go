package main

import (
	"bytes"
	"github.com/GeoNet/weft"
	"net/http"
)

func valV2(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
		return res
	}

	var d string

	err := db.QueryRow(`SELECT row_to_json(fc)
                         FROM ( SELECT 'FeatureCollection' as type, array_to_json(array_agg(f)) as features
                         FROM (SELECT 'Feature' as type,
                         ST_AsGeoJSON(v.location)::json as geometry,
                         row_to_json((SELECT l FROM 
                         	(
                         		SELECT 
                                id AS "volcanoID",
                                title AS "volcanoTitle",
                                alert_level as "level",
                                activity,
                                hazards 
                           ) as l
                         )) as properties FROM (haz.volcano JOIN haz.volcanic_alert_level using (alert_level)) as v ) As f )  as fc`).Scan(&d)
	if err != nil {
		return weft.ServiceUnavailableError(err)
	}

	b.WriteString(d)
	h.Set("Content-Type", V2GeoJSON)
	return &weft.StatusOK
}


func quakesVolcanoRegionV2(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {

	if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
		return res
	}

	var err error
	var days int = 90
	var volcanoId string


	if volcanoId, err = getVolcanoIDQuake(r); err != nil {
		return weft.BadRequest(err.Error())
	}


	var d string
	err = db.QueryRow(`SELECT row_to_json(fc)
			FROM ( SELECT 'FeatureCollection' as type, COALESCE(array_to_json(array_agg(f)), '[]') as features
			FROM (SELECT 'Feature' as type,
			ST_AsGeoJSON(q.geom)::json as geometry,
			row_to_json((SELECT l FROM
				(
				SELECT	publicid AS "publicID",
					to_char(time, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"') as "time",
					depth,
					magnitude,
					locality,
					intensity,
					intensity_newzealand as "regionIntensity",
					quality
				) as l
			)) as properties
			FROM haz.quakeapi as q
			where st_covers((select region from haz.volcano hv where id = $1 and q.depth < hv.depth and DATE_PART('day', now() - q.time) < $2), q.geom)
			ORDER BY time DESC ) as f ) as fc;`, volcanoId, days).Scan(&d)
	if err != nil {
		return weft.ServiceUnavailableError(err)
	}

	b.WriteString(d)
	h.Set("Content-Type", V2GeoJSON)
	return &weft.StatusOK
}

func volcanoRegionV2(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {

	if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
		return res
	}

	var err error
	var volcanoId string

	if volcanoId, err = getVolcanoIDRegion(r); err != nil {
		return weft.BadRequest(err.Error())
	}

	var d string
	err = db.QueryRow(`SELECT row_to_json(fc)
			FROM ( SELECT 'FeatureCollection' as type, COALESCE(array_to_json(array_agg(f)), '[]') as features
			FROM (SELECT 'Feature' as type,
			ST_AsGeoJSON(r.region)::json as geometry,
			row_to_json((SELECT l FROM
			(
			SELECT	id,
			title
			) as l
			)) as properties
			FROM haz.volcano as r
			where id = $1) as f ) as fc;`, volcanoId).Scan(&d)
	if err != nil {
		return weft.ServiceUnavailableError(err)
	}

	b.WriteString(d)
	h.Set("Content-Type", V2GeoJSON)
	return &weft.StatusOK
}
