package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
)

const feltURL = "http://felt.geonet.org.nz/services/reports/"

func feltV1(r *http.Request, h http.Header, b *bytes.Buffer) *result {
	if res := checkQuery(r, []string{"publicID"}, []string{}); !res.ok {
		return res
	}

	var publicID string
	var res *result

	if publicID, res = getPublicID(r); !res.ok {
		return res
	}

	var rs *http.Response
	rs, err := client.Get(feltURL + publicID + ".geojson")
	defer rs.Body.Close()
	if err != nil {
		return serviceUnavailableError(err)
	}

	bt, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		return serviceUnavailableError(err)
	}

	// Felt returns a 400 when it should probably be a 404.  Tapestry quirk?
	switch {
	case http.StatusOK == rs.StatusCode:
		h.Set("Content-Type", V1GeoJSON)
		b.Write(bt)
		return &statusOK
	case 4 == rs.StatusCode/100:
		res := &notFound
		res.msg = string(bt)
		return res
	case 5 == rs.StatusCode/500:
		return serviceUnavailableError(errors.New("error proxying felt resports.  Shrug."))
	}

	return serviceUnavailableError(errors.New("unknown response from felt."))
}
