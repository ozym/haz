package main

import (
	"database/sql"
	"github.com/GeoNet/haz/msg"
	"github.com/GeoNet/web"
	"github.com/GeoNet/web/api/apidoc"
	"html/template"
	"net/http"
	"regexp"
	"strings"
)

// These constants are the length of parts of the URI and are used for
// extracting query params embedded in the URI.
const (
	quakeLen = 7 //  len("/quake/")
)

var quakeDoc = apidoc.Endpoint{Title: "Quake",
	Description: `Look up quake information.`,
	Queries: []*apidoc.Query{
		quakeD,
		quakesRegionD,
	},
}

var intensityRe = regexp.MustCompile(`^(unnoticeable|weak|light|moderate|strong|severe)$`)
var numberRe = regexp.MustCompile(`^(3|30|100|500|1000|1500)$`)
var qualityRe = regexp.MustCompile(`^(best|caution|deleted|good)$`)
var publicIDRe = regexp.MustCompile(`^[0-9a-z]+$`)

// all requests have the same properties in the return.
// this is a map for all apidoc.Query{} structs.
var propsD = map[string]template.HTML{
	`publicID`:        `the unique public identifier for this quake.`,
	`time`:            `the origin time of the quake.`,
	`depth`:           `the depth of the quake in km.`,
	`magnitude`:       `the summary magnitude for the quake.  This is <b>not</b> Richter magnitude.`,
	`locality`:        `distance and direction to the nearest locality.`,
	`intensity`:       `the calculated <a href="http://info.geonet.org.nz/x/b4Ih">intensity</a> at the surface above the quake (epicenter) e.g., strong.`,
	`regionIntensity`: `the calculated intensity at the closest locality in the region for the request. `,
	`quality`:         `the quality of this information; <code>best</code>, <code>good</code>, <code>caution</code>, <code>deleted</code>.`,
}

// /quake/2013p407387

var quakeD = &apidoc.Query{
	Accept:      web.V1GeoJSON,
	Title:       "Quake",
	Description: "Information for a single quake.",
	Example:     "/quake/2013p407387",
	ExampleHost: exHost,
	URI:         "/quake/(publicID)",
	Params: map[string]template.HTML{
		"publicID": `a valid quake ID e.g., <code>2014p715167</code>`,
	},
	Props: propsD,
}

func quake(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query()) != 0 {
		web.BadRequest(w, r, "incorrect number of query parameters.")
		return
	}

	publicID := r.URL.Path[quakeLen:]

	// TODO bother with this?
	if !publicIDRe.MatchString(publicID) {
		web.BadRequest(w, r, "invalid publicID: "+publicID)
		return
	}

	var d string

	// Check that the publicid exists in the DB.  This is needed as the handle method will return empty
	// JSON for an invalid publicID.
	err := db.QueryRow("select publicid FROM haz.quake where publicid = $1", publicID).Scan(&d)
	if err == sql.ErrNoRows {
		web.NotFound(w, r, "invalid publicID: "+publicID)
		return
	}
	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return
	}

	err = db.QueryRow(
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
		web.ServiceUnavailable(w, r, err)
		return
	}

	b := []byte(d)
	web.Ok(w, r, &b)
}

// /quake?regionID=newzealand&regionIntensity=unnoticeable&number=30&quality=best,caution,good
var quakesRegionD = &apidoc.Query{
	Accept:      web.V1GeoJSON,
	Title:       "Quakes Possibly Felt in a Region",
	Description: "quakes possibly felt in a region during the last 365 days.",
	Example:     "/quake?regionID=newzealand&regionIntensity=weak&number=3&quality=best,caution,good",
	ExampleHost: exHost,
	URI:         "/quake?regionID=(region)&regionIntensity=(intensity)&number=(n)&quality=(quality)",
	Required: map[string]template.HTML{
		`regionID`: `a valid quake region identifier the only permissable value is <code>newzealand</code>.`,
		`regionIntensity`: `the minimum intensity in the region e.g., <code>weak</code>.  
		Must be one of <code>unnoticeable</code>, <code>weak</code>, <code>light</code>, 
		<code>moderate</code>, <code>strong</code>, <code>severe</code>.`,
		`number`: `the maximum number of quakes to return.  Must be one of 
		<code>3</code>, <code>30</code>, <code>100</code>, <code>500</code>, <code>1000</code>, <code>1500</code>.`,
		`quality`: `a comma separated list of quality values to be included in the response.  The only allowable option is:  
		<code>best</code>,<code>caution</code>,<code>deleted</code>,<code>good</code>.`,
	},
	Props: propsD,
}

func quakesRegion(w http.ResponseWriter, r *http.Request) {
	if err := quakesRegionD.CheckParams(r.URL.Query()); err != nil {
		web.BadRequest(w, r, err.Error())
		return
	}

	v := r.URL.Query()

	number := v.Get("number")
	regionID := v.Get("regionID")
	regionIntensity := v.Get("regionIntensity")
	quality := strings.Split(v.Get("quality"), ",")

	if regionID != "newzealand" {
		web.BadRequest(w, r, "Invalid query parameter regionID: "+regionID)
		return
	}

	if !numberRe.MatchString(number) {
		web.BadRequest(w, r, "Invalid query parameter number: "+number)
		return
	}

	if !intensityRe.MatchString(regionIntensity) {
		web.BadRequest(w, r, "Invalid regionIntensity: "+regionIntensity)
		return
	}

	for _, q := range quality {
		if !qualityRe.MatchString(q) {
			web.BadRequest(w, r, "Invalid quality: "+q)
			return
		}
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
                                intensity,
                                intensity_newzealand as "regionIntensity",
                                quality
                           ) as l
                         )) as properties FROM haz.quakeapi as q where mmid_newzealand >= $1
                         ORDER BY time DESC  limit $2 ) as f ) as fc`, msg.IntensityMMI(regionIntensity), number).Scan(&d)
	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return
	}

	b := []byte(d)
	web.Ok(w, r, &b)
}
