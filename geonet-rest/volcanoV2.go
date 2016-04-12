package main

import (
	"bytes"
	"net/http"
)

func valV2(r *http.Request, h http.Header, b *bytes.Buffer) *result {
	if res := checkQuery(r, []string{}, []string{}); !res.ok {
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
		return serviceUnavailableError(err)
	}

	b.WriteString(d)
	h.Set("Content-Type", V2GeoJSON)
	return &statusOK
}
