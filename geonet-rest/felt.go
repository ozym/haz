package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"github.com/GeoNet/weft"
)

const feltURL = "http://felt.geonet.org.nz/services/reports/"

func feltV1(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	if res := weft.CheckQuery(r, []string{"publicID"}, []string{}); !res.Ok {
		return res
	}

	var publicID string
	var res *weft.Result

	if publicID, res = getPublicID(r); !res.Ok {
		return res
	}

	var rs *http.Response
	rs, err := client.Get(feltURL + publicID + ".geojson")
	defer rs.Body.Close()
	if err != nil {
		return weft.ServiceUnavailableError(err)
	}

	bt, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		return weft.ServiceUnavailableError(err)
	}

	// Felt returns a 400 when it should probably be a 404.  Tapestry quirk?
	switch {
	case http.StatusOK == rs.StatusCode:
		h.Set("Content-Type", V1GeoJSON)
		b.Write(bt)
		return &weft.StatusOK
	case 4 == rs.StatusCode/100:
		//res := &notFound
		//res.msg = string(bt)
		return &weft.NotFound
	case 5 == rs.StatusCode/500:
		return weft.ServiceUnavailableError(errors.New("error proxying felt resports.  Shrug."))
	}

	return weft.ServiceUnavailableError(errors.New("unknown response from felt."))
}
