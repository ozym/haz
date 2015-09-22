package main

import (
	"github.com/GeoNet/web"
	"net/http"
	"time"
)

func intensityV2(w http.ResponseWriter, r *http.Request) {
	if badQuery(w, r, []string{"type"}, []string{"publicID"}) {
		return
	}

	var t string
	var ok bool

	if t, ok = getIntensityType(w, r); !ok {
		return
	}

	var d string
	var err error

	switch t {
	case "measured":
		d, err = measuredLatest()
	case "reported":
		switch r.URL.Query().Get("publicID") {
		case "":
			d, err = reportedLatest()
		default:
			var t time.Time
			var ok bool
			if t, ok = getQuakeTime(w, r); !ok {
				return
			}
			d, err = reportedWindow(t)
		}
	}

	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return
	}

	b := []byte(d)
	w.Header().Set("Content-Type", web.V2GeoJSON)
	web.Ok(w, r, &b)
}

func measuredLatest() (string, error) {
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
	return d, err
}

func reportedLatest() (string, error) {
	var d string
	err := db.QueryRow(
		`WITH features as (
	select COALESCE(array_to_json(array_agg(fs)), '[]') as features from (SELECT 'Feature' as type,
						ST_AsGeoJSON(s.location)::json as geometry,
						row_to_json(( select l from 
							( 
							select mmi,
							count
							) as l )) 
							as properties from (select st_pointfromgeohash(geohash6) as location, 
							max(mmi) as mmi, 
							count(mmi) as count 
							FROM impact.intensity_reported 
								WHERE time >= (now() - interval '60 minutes')
							group by (geohash6)) as s) as fs
		), summary as (
			select COALESCE(json_object_agg(summ.mmi, summ.count), '{}') as count_mmi, COALESCE(sum(count), 0) as count
			from (select mmi as mmi, count(*) as count from impact.intensity_reported 
			WHERE time >= (now() - interval '60 minutes')
			group by mmi
			) as summ
		)
		SELECT row_to_json(fc)
			FROM ( SELECT 'FeatureCollection' as type, 
					features.features, 
					summary.count_mmi,
					summary.count
				FROM features, summary )  as fc`).Scan(&d)
	return d, err
}

func reportedWindow(t time.Time) (string, error) {
	var d string
	var err error

	err = db.QueryRow(`WITH features as (
	select COALESCE(array_to_json(array_agg(fs)), '[]') as features from (SELECT 'Feature' as type,
						ST_AsGeoJSON(s.location)::json as geometry,
						row_to_json(( select l from 
							( 
							select mmi,
							count
							) as l )) 
							as properties from (select st_pointfromgeohash(geohash6) as location, 
							max(mmi) as mmi, 
							count(mmi) as count 
							FROM impact.intensity_reported 
								WHERE time >= $1
								AND time <= $2
							group by (geohash6)) as s) as fs
		), summary as (
			select COALESCE(json_object_agg(summ.mmi, summ.count), '{}') as count_mmi, COALESCE(sum(count), 0) as count
			from (select mmi as mmi, count(*) as count from impact.intensity_reported 
			WHERE time >= $1
			AND time <= $2
			group by mmi
			) as summ
		)
		SELECT row_to_json(fc)
			FROM ( SELECT 'FeatureCollection' as type, 
					features.features, 
					summary.count_mmi,
					summary.count
				FROM features, summary )  as fc`, t.Add(time.Duration(-1*time.Minute)), t.Add(time.Duration(15*time.Minute))).Scan(&d)
	return d, err
}
