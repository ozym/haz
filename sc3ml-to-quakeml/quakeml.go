package main

import (
	"bytes"
	"fmt"
	"github.com/GeoNet/weft"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

const s3 = "http://seiscompml07.s3-website-ap-southeast-2.amazonaws.com/"

func init() {
	client = &http.Client{
		Timeout: timeout,
	}
}

func quakeml12(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	by, res := getBytes(s3 + strings.TrimPrefix(r.URL.Path, "/quakeml/1.2/") + ".xml", "")
	if !res.Ok {
		return res
	}

	return xslt(by, b, "sc3ml_0.7__quakeml_1.2.xsl")
}

func quakeml12RT(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	by, res := getBytes(s3 + strings.TrimPrefix(r.URL.Path, "/quakeml-rt/1.2/") + ".xml", "")
	if !res.Ok {
		return res
	}

	return xslt(by, b, "sc3ml_0.7__quakeml_1.2-RT.xsl")
}

func csv(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	p := strings.Split(strings.TrimPrefix(r.URL.Path, "/csv/1.0.0/"), "/")

	if len(p) < 2 {
		return &weft.NotFound
	}

	by, res := getBytes(s3 + p[0] + ".xml", "")
	if !res.Ok {
		return res
	}

	// remove the publicID (p[0])
	p = append(p[:0], p[1:]...)

	switch strings.Join(p, "/") {
	case "picks":
		return xslt(by, b, "sc3ml_0.7_to_csv.xsl", "--stringparam", "picks", "true")
	case "event":
		return xslt(by, b, "sc3ml_0.7_to_csv.xsl", "--stringparam", "event", "true")
	case "event/picks":
		return xslt(by, b, "sc3ml_0.7_to_csv.xsl", "--stringparam", "picks", "true", "--stringparam", "event", "true")
	default:
		return &weft.NotFound
	}
}

func xslt(src []byte, b *bytes.Buffer, xsl string, args ...string) *weft.Result {
	cmd := exec.Command("/usr/bin/xsltproc")
	cmd.Args = append(cmd.Args, args...)
	cmd.Args = append(cmd.Args, "assets/" + xsl)
	cmd.Args = append(cmd.Args, "-")

	var err error
	var in io.WriteCloser
	var out io.ReadCloser

	if in, err = cmd.StdinPipe(); err != nil {
		return weft.InternalServerError(err)
	}
	defer in.Close()

	if out, err = cmd.StdoutPipe(); err != nil {
		return weft.InternalServerError(err)
	}
	defer out.Close()

	if err := cmd.Start(); err != nil {
		return weft.InternalServerError(err)
	}
	if _, err = in.Write(src); err != nil {
		return weft.InternalServerError(err)
	}

	if err = in.Close(); err != nil {
		return weft.InternalServerError(err)
	}

	if _, err = b.ReadFrom(out); err != nil {
		return weft.InternalServerError(err)
	}

	if err = cmd.Wait(); err != nil {
		return weft.InternalServerError(err)
	}

	return &weft.StatusOK
}

/*
getBytes fetches bytes for the requested url.  accept
may be left as the empty string.
*/
func getBytes(url, accept string) ([]byte, *weft.Result) {
	var r *http.Response
	var req *http.Request
	var err error
	var b []byte

	if accept == "" {
		r, err = client.Get(url)
		if err != nil {
			return b, weft.InternalServerError(err)
		}
	} else {
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			return b, weft.InternalServerError(err)
		}

		req.Header.Add("Accept", accept)

		r, err = client.Do(req)
		if err != nil {
			return b, weft.InternalServerError(err)
		}
	}
	defer r.Body.Close()

	b, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return b, weft.InternalServerError(err)
	}

	switch r.StatusCode {
	case http.StatusOK:
		return b, &weft.StatusOK
	case http.StatusNotFound:
		return b, &weft.NotFound
	default:
		// TODO do we need to handle more errors here?
		return b, weft.InternalServerError(fmt.Errorf("server error"))

	}

	return b, &weft.StatusOK
}
