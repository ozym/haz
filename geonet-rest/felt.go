package main

import (
	"database/sql"
	"errors"
	"github.com/GeoNet/web"
	"github.com/GeoNet/web/api/apidoc"
	"html/template"
	"io/ioutil"
	"net/http"
)

const (
	feltURL = "http://felt.geonet.org.nz/services/reports/"
)

var feltDoc = apidoc.Endpoint{
	Title:       "Felt",
	Description: `Look up Felt Report information.`,
	Queries: []*apidoc.Query{
		feltD,
	},
}

var feltD = &apidoc.Query{
	Accept:      web.V1GeoJSON,
	Title:       "Felt",
	Description: "Look up Felt Report information about earthquakes",
	Example:     "/felt/report?publicID=2013p407387",
	ExampleHost: exHost,
	URI:         "/felt/report?publicID=(publicID)",
	Required: map[string]template.HTML{
		"publicID": `a valid quake ID e.g., <code>2014p715167</code>`,
	},
	Props: map[string]template.HTML{
		"todo": `todo`,
	},
}

func felt(w http.ResponseWriter, r *http.Request) {
	if err := feltD.CheckParams(r.URL.Query()); err != nil {
		web.BadRequest(w, r, err.Error())
		return
	}

	publicID := r.URL.Query().Get("publicID")

	var d string

	err := db.QueryRow("select publicid FROM haz.quake where publicid = $1", publicID).Scan(&d)
	if err == sql.ErrNoRows {
		web.NotFound(w, r, "invalid publicID: "+publicID)
		return
	}
	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return
	}

	res, err := client.Get(feltURL + publicID + ".geojson")
	defer res.Body.Close()
	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return
	}

	// Felt returns a 400 when it should probably be a 404.  Tapestry quirk?
	switch {
	case 200 == res.StatusCode:
		web.Ok(w, r, &b)
		return
	case 4 == res.StatusCode/100:
		web.NotFound(w, r, string(b))
		return
	case 5 == res.StatusCode/500:
		web.ServiceUnavailable(w, r, errors.New("error proxying felt resports.  Shrug."))
		return
	}

	web.ServiceUnavailable(w, r, errors.New("unknown response from felt."))
}
