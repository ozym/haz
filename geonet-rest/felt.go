package main

import (
	"errors"
	"github.com/GeoNet/web"
	"io/ioutil"
	"net/http"
)

const feltURL = "http://felt.geonet.org.nz/services/reports/"

func feltV1(w http.ResponseWriter, r *http.Request) {
	if badQuery(w, r, []string{"publicID"}, []string{}) {
		return
	}

	var publicID string
	var ok bool

	if publicID, ok = getPublicID(w, r); !ok {
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
		w.Header().Set("Content-Type", web.V1GeoJSON)
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
