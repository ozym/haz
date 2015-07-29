package main

import (
	"github.com/GeoNet/web"
	"github.com/GeoNet/web/api/apidoc"
	"html/template"
	"net/http"
)

var volcanoDoc = apidoc.Endpoint{Title: "Volcano",
	Description: "Look up volcano information.",
	Queries: []*apidoc.Query{
		alertLevelD,
		alertBulletinD,
	},
}

const alertBulletinURL = `http://info.geonet.org.nz/createrssfeed.action?types=blogpost&spaces=volc&title=GeoNet+Volcano+RSS+Feed&labelString=vab&excludedSpaceKeys%3D&sort=created&maxResults=10&timeSpan=500&showContent=true&publicFeed=true&confirm=Create+RSS+Feed`

var alertLevelD = &apidoc.Query{
	Accept:      web.V1GeoJSON,
	Title:       "Volcanic Alert Level",
	Description: `Volcanic Alert Level information for all volcanoes.`,
	Discussion:  `<p>Volcanic Alert Level information for all volcanoes.  Please refer to <a href="http://info.geonet.org.nz/x/PYAO">Volcanic Alert Levels</a> for additional information.</p>`,
	Example:     "/volcano/alert/level",
	ExampleHost: exHost,
	URI:         "/volcano/alert/level",
	Required: map[string]template.HTML{
		"none": `no query parameters are required.`,
	},
	Props: map[string]template.HTML{
		`volcanoID`:    `a unique identifier for the volcano.`,
		`volcanoTitle`: `the volcano title.`,
		`level`:        `volcanic alert level.`,
		`activity`:     `volcanic activity.`,
		`hazards`:      `most likely hazards.`,
	},
}

func alertLevel(w http.ResponseWriter, r *http.Request) {
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
	web.Ok(w, r, &b)
}

var alertBulletinD = &apidoc.Query{
	Accept:      web.V1JSON,
	Title:       "Volcanic Alert Bulletins",
	Description: " Returns a simple JSON version of the GeoNet Volcanic Alert Bulletin RSS feed.",
	Example:     "/volcano/alert/bulletin",
	ExampleHost: exHost,
	URI:         "/volcano/alert/bulletin",
	Required: map[string]template.HTML{
		"none": `no query parameters are required.`,
	},
	Props: map[string]template.HTML{
		"mlink":     "a link to a mobile version of the bulletin.",
		"link":      "a link to the bulletin.",
		"published": "the date the bulletin was published",
		"title":     "the title of the bulletin.",
	},
}

func alertBulletin(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query()) != 0 {
		web.BadRequest(w, r, "incorrect number of query parameters.")
		return
	}

	j, err := fetchRSS(alertBulletinURL)
	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return
	}

	w.Header().Set("Surrogate-Control", web.MaxAge300)

	web.Ok(w, r, &j)
}
