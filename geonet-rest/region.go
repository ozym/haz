package main

import (
	"database/sql"
	"github.com/GeoNet/web"
	"github.com/GeoNet/web/api/apidoc"
	"html/template"
	"net/http"
)

// These constants are the length of parts of the URI and are used for
// extracting query params embedded in the URI.
const (
	regionLen = 8 // len("/region/")
)

var regionDoc = apidoc.Endpoint{
	Title:       "Region",
	Description: `Look up region information.`,
	Queries: []*apidoc.Query{
		regionsD,
		regionD,
	},
}

var regionsD = &apidoc.Query{
	Accept:      web.V1GeoJSON,
	Title:       "Regions",
	Description: "Retrieve regions.",
	Example:     "/region?type=quake",
	ExampleHost: exHost,
	URI:         "/region?type=(type)",
	Required: map[string]template.HTML{
		"type": `the region type.  The only allowable value is <code>quake</code>.`,
	},
	Props: map[string]template.HTML{
		`regionID`: `a unique indentifier for the region.`,
		`title`:    `the region title.`,
		`group`:    `the region group.`,
	},
}

// just quake regions at the moment.
func regions(w http.ResponseWriter, r *http.Request) {
	if err := regionsD.CheckParams(r.URL.Query()); err != nil {
		web.BadRequest(w, r, err.Error())
		return
	}

	if r.URL.Query().Get("type") != "quake" {
		web.BadRequest(w, r, "type must be quake.")
		return
	}

	var d string

	err := db.QueryRow(`SELECT row_to_json(fc)
                         FROM ( SELECT 'FeatureCollection' as type, array_to_json(array_agg(f)) as features
                         FROM (SELECT 'Feature' as type,
                         ST_AsGeoJSON(q.geom)::json as geometry,
                         row_to_json((SELECT l FROM
                         	(
                         		SELECT
                         		regionname as "regionID",
                         		title,
                         		groupname as group
                           ) as l
                         )) as properties FROM haz.quakeregion as q where groupname in ('region', 'north', 'south')) as f ) as fc`).Scan(&d)

	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return
	}

	w.Header().Set("Surrogate-Control", web.MaxAge86400)
	b := []byte(d)
	web.Ok(w, r, &b)
}

// /region/wellington

var regionD = &apidoc.Query{
	Accept:      web.V1GeoJSON,
	Title:       "Region",
	Description: "Retrieve a single region.",
	Example:     "/region/wellington",
	ExampleHost: exHost,
	URI:         "/region/(regionID)",
	Required: map[string]template.HTML{
		"regionID": `A region ID e.g., <code>wellington</code>.`,
	},
	Props: map[string]template.HTML{
		`regionID`: `a unique indentifier for the region.`,
		`title`:    `the region title.`,
		`group`:    `the region group.`,
	},
}

func region(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query()) != 0 {
		web.BadRequest(w, r, "incorrect number of query parameters.")
		return
	}

	regionID := r.URL.Path[regionLen:]

	var d string

	err := db.QueryRow("select regionname FROM haz.quakeregion where regionname = $1", regionID).Scan(&d)
	if err == sql.ErrNoRows {
		web.BadRequest(w, r, "invalid regionID: "+regionID)
		return
	}
	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return
	}

	err = db.QueryRow(`SELECT row_to_json(fc)
                         FROM ( SELECT 'FeatureCollection' as type, array_to_json(array_agg(f)) as features
                         FROM (SELECT 'Feature' as type,
                         ST_AsGeoJSON(q.geom)::json as geometry,
                         row_to_json((SELECT l FROM 
                         	(
                         		SELECT 
                         		regionname as "regionID",
                         		title, 
                         		groupname as group
                           ) as l
                         )) as properties FROM haz.quakeregion as q where regionname = $1 ) as f ) as fc`, regionID).Scan(&d)
	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return
	}

	w.Header().Set("Surrogate-Control", web.MaxAge86400)
	b := []byte(d)
	web.Ok(w, r, &b)
}
