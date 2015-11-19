package main

import (
	"github.com/GeoNet/web"
	"net/http"
)

func valV2(w http.ResponseWriter, r *http.Request) {
	if badQuery(w, r, []string{}, []string{}) {
		return
	}

	if len(r.URL.Query()) != 0 {
		web.BadRequest(w, r, "incorrect number of query parameters.")
		return
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
		web.ServiceUnavailable(w, r, err)
		return
	}

	b := []byte(d)
	w.Header().Set("Content-Type", web.V2GeoJSON)
	web.Ok(w, r, &b)
}
