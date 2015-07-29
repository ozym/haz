package main

import (
	"database/sql"
	"github.com/GeoNet/web"
	"github.com/GeoNet/web/api/apidoc"
	"html/template"
	"net/http"
	"regexp"
	"time"
)

var impactDoc = apidoc.Endpoint{Title: "Impact",
	Description: `Look up impact information`,
	Queries: []*apidoc.Query{
		// intensityReportedD,
		// intensityReportedLatestD,
		intensityMeasuredLatestD,
	},
}

var zoomRe = regexp.MustCompile(`^(5|6)$`)

var intensityMeasuredLatestD = &apidoc.Query{
	Accept:      web.V1GeoJSON,
	Title:       "Measured Intensity - Latest",
	Description: "Retrieve measured intensity information in the last sixty minutes.",
	Example:     "/intensity?type=measured",
	ExampleHost: exHost,
	URI:         "/intensity?type",
	Required: map[string]template.HTML{
		"type": `<code>measured</code> is the only allowed value.`,
	},
	Props: map[string]template.HTML{
		"max_mmi": `the maximum <a href="http://info.geonet.org.nz/x/w4IO">Modified Mercalli Intensity (MMI)</a> measured at the point in the last sixty minutes.`,
	},
}

func intensityMeasuredLatest(w http.ResponseWriter, r *http.Request) {
	if err := intensityMeasuredLatestD.CheckParams(r.URL.Query()); err != nil {
		web.BadRequest(w, r, err.Error())
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
	web.Ok(w, r, &b)
}

// latest reported intensity

var intensityReportedLatestD = &apidoc.Query{
	Accept:      web.V1GeoJSON,
	Title:       "Reported Intensity - Latest",
	Description: "Retrieve reported intensity information in the last sixty minutes.",
	Example:     "/intensity?type=reported&zoom=5",
	ExampleHost: exHost,
	URI:         "/intensity?type=reported&zoom=(int)",
	Required: map[string]template.HTML{
		"type": `<code>reported</code> is the only allowed value.`,
		"zoom": `The zoom level to aggregate values at.  This controls the size of the area that values are aggregated at.  The point returned
				will be the center of each area.  Allowed values are one of <code>5, 6</code>.`,
	},
	Props: map[string]template.HTML{
		"max_mmi": `the maximum <a href="http://info.geonet.org.nz/x/w4IO">Modified Mercalli Intensity (MMI)</a> 
					in the area around the point in the last sixty minutes.`,
		"min_mmi": `the minimum <a href="http://info.geonet.org.nz/x/w4IO">Modified Mercalli Intensity (MMI)</a> 
					in the area of around the point in the last sixty minutes.`,
		"count": `the count of <a href="http://info.geonet.org.nz/x/w4IO">Modified Mercalli Intensity (MMI)</a> 
					values reported in the area of around the point in the last sixty minutes.`,
	},
}

func intensityReportedLatest(w http.ResponseWriter, r *http.Request) {
	if err := intensityReportedLatestD.CheckParams(r.URL.Query()); err != nil {
		web.BadRequest(w, r, err.Error())
		return
	}

	if r.URL.Query().Get("type") != "reported" {
		web.BadRequest(w, r, "type must be reported.")
		return
	}

	zoom := r.URL.Query().Get("zoom")

	if !zoomRe.MatchString(r.URL.Query().Get("zoom")) {
		web.BadRequest(w, r, "Invalid zoom")
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
										select max_mmi,
										min_mmi,
										count
										) as l )) 
					as properties from (select st_pointfromgeohash(geohash` + zoom + `) as location, 
						min(mmi) as min_mmi, 
						max(mmi) as max_mmi, 
						count(mmi) as count 
						FROM impact.intensity_reported  
						WHERE time >= (now() - interval '60 minutes')
						group by (geohash` + zoom + `)) as s
					) As f )  as fc`).Scan(&d)
	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return
	}

	b := []byte(d)
	web.Ok(w, r, &b)
}

// reported intensity

var intensityReportedD = &apidoc.Query{
	Accept:      web.V1GeoJSON,
	Title:       "Reported Intensity",
	Description: "Retrieve reported intensity information in a 15 minute time window after an event.",
	Example:     "/intensity?type=reported&zoom=5&publicID=2013p407387",
	ExampleHost: exHost,
	URI:         "/intensity?type=reported&zoom=(int)&publicID=(publicID)",
	Required: map[string]template.HTML{
		"type": `<code>reported</code> is the only allowed value.`,
		"zoom": `The zoom level to aggregate values at.  This controls the size of the area that values are aggregated at.  The point returned
						will be the center of each area.  Allowed values are one of <code>5, 6</code>.`,
		"publicID": `a valid quake ID e.g., <code>2014p715167</code>`,
	},
	Props: map[string]template.HTML{
		"max_mmi": `the maximum <a href="http://info.geonet.org.nz/x/w4IO">Modified Mercalli Intensity (MMI)</a> 
							in the area around the point in the last sixty minutes.`,
		"min_mmi": `the minimum <a href="http://info.geonet.org.nz/x/w4IO">Modified Mercalli Intensity (MMI)</a> 
							in the area of around the point in the last sixty minutes.`,
		"count": `the count of <a href="http://info.geonet.org.nz/x/w4IO">Modified Mercalli Intensity (MMI)</a> 
							values reported in the area of around the point in the last sixty minutes.`,
	},
}

func intensityReported(w http.ResponseWriter, r *http.Request) {
	if err := intensityReportedD.CheckParams(r.URL.Query()); err != nil {
		web.BadRequest(w, r, err.Error())
		return
	}

	if r.URL.Query().Get("type") != "reported" {
		web.BadRequest(w, r, "type must be reported.")
		return
	}

	zoom := r.URL.Query().Get("zoom")

	if !zoomRe.MatchString(r.URL.Query().Get("zoom")) {
		web.BadRequest(w, r, "Invalid zoom")
		return
	}

	// Check that the publicid exists in the DB.
	// If it does we keep the origintime - we need it later on.
	publicID := r.URL.Query().Get("publicID")

	originTime := time.Time{}

	err := db.QueryRow("select origintime FROM haz.quake where publicid = $1", publicID).Scan(&originTime)
	if err == sql.ErrNoRows {
		web.NotFound(w, r, "invalid publicID: "+publicID)
		return
	}
	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return
	}

	query := `SELECT row_to_json(fc)
				FROM ( SELECT 'FeatureCollection' as type, COALESCE(array_to_json(array_agg(f)), '[]') as features
					FROM (SELECT 'Feature' as type,
						ST_AsGeoJSON(s.location)::json as geometry,
						row_to_json(( select l from 
							( 
							select max_mmi,
							min_mmi,
							count
							) as l )) 
							as properties from (select st_pointfromgeohash(geohash` + zoom + `) as location, 
							min(mmi) as min_mmi, 
							max(mmi) as max_mmi, 
							count(mmi) as count 
							FROM impact.intensity_reported 
							WHERE time >= $1
							AND time <= $2
							group by (geohash` + zoom + `)) as s
			) As f )  as fc`

	var d string

	err = db.QueryRow(query, originTime.Add(time.Duration(-1*time.Minute)), originTime.Add(time.Duration(15*time.Minute))).Scan(&d)
	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return
	}

	b := []byte(d)
	web.Ok(w, r, &b)
}
