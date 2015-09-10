package main

import (
	"github.com/GeoNet/web"
	"github.com/GeoNet/web/api/apidoc"
	"html/template"
	"net/http"
	"regexp"
)

var impactDoc = apidoc.Endpoint{Title: "Impact",
	Description: `Look up impact information`,
	Queries: []*apidoc.Query{
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
