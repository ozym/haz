package main

import (
	"bytes"
	_ "github.com/GeoNet/log/logentries"
	"github.com/GeoNet/mtr/mtrapp"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

// For setting Cache-Control and Surrogate-Control headers.
const (
	maxAge10    = "max-age=10"
	maxAge300   = "max-age=300"
	maxAge86400 = "max-age=86400"
)

const (
	V1GeoJSON = "application/vnd.geo+json;version=1"
	V1JSON    = "application/json;version=1"
	V1CSV     = "text/csv;version=1"
	V2GeoJSON = "application/vnd.geo+json;version=2"
	V2JSON    = "application/json;version=2"
	V2CSV     = "text/csv;version=2"
)

// These are for CAP format and Atom which is not versioned by Accept.
const (
	CAP  = "application/cap+xml"
	Atom = "application/xml"
)

// These constants are for error and other pages.  They can be changed.
const (
	ErrContent  = "text/plain; charset=utf-8"
	HtmlContent = "text/html; charset=utf-8"
)

type Header struct {
	Cache, Surrogate string // Set as the default in the response header - can override in handler funcs.
	Vary             string // This is added to the response header (which may already Vary on gzip).
}
type result struct {
	ok   bool   // set true to indicated success
	code int    // http status code for writing back the client e.g., http.StatusOK for success.
	msg  string // any error message for logging or to send to the client.
}

/*
requestHandler for handling http requests.  The response for the request
should be written into.  Any header values for the client can be set in h
e.g., Content-Type.
*/
type requestHandler func(r *http.Request, h http.Header, b *bytes.Buffer) *result

var (
	statusOK         = result{ok: true, code: http.StatusOK, msg: ""}
	methodNotAllowed = result{ok: false, code: http.StatusMethodNotAllowed, msg: "method not allowed"}
	notFound         = result{ok: false, code: http.StatusNotFound, msg: ""}
	notAcceptable    = result{ok: false, code: http.StatusNotAcceptable, msg: "specify accept"}
)

func internalServerError(err error) *result {
	return &result{ok: false, code: http.StatusInternalServerError, msg: err.Error()}
}

func serviceUnavailableError(err error) *result {
	return &result{ok: false, code: http.StatusServiceUnavailable, msg: err.Error()}
}

func badRequest(message string) *result {
	return &result{ok: false, code: http.StatusBadRequest, msg: message}
}

/*
checkQuery inspects r and makes sure all required query parameters
are present and that no more than the required and optional parameters
are present.
*/
func checkQuery(r *http.Request, required, optional []string) *result {
	if strings.Contains(r.URL.Path, ";") {
		return badRequest("cache buster")
	}

	v := r.URL.Query()

	if len(required) == 0 && len(optional) == 0 {
		if len(v) == 0 {
			return &statusOK
		} else {
			return badRequest("found unexpected query parameters")
		}
	}

	var missing []string

	for _, k := range required {
		if v.Get(k) == "" {
			missing = append(missing, k)
		} else {
			v.Del(k)
		}
	}

	switch len(missing) {
	case 0:
	case 1:
		return badRequest("missing required query parameter: " + missing[0])
	default:
		return badRequest("missing required query parameters: " + strings.Join(missing, ", "))
	}

	for _, k := range optional {
		v.Del(k)
	}

	if len(v) > 0 {
		return badRequest("found additional query parameters")
	}

	return &statusOK
}

/*
toHandler adds basic auth to f and returns a handler.
*/
func toHandler(f requestHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Find the name of the function f to use as the timer id
		id := r.Method
		fn := runtime.FuncForPC(reflect.ValueOf(f).Pointer())
		if fn != nil {
			id = fn.Name() + "." + id
		}

		t := mtrapp.Start(id)

		mtrapp.Requests.Inc()

		switch r.Method {
		case "GET":
			var b bytes.Buffer
			w.Header().Set("Cache-Control", maxAge10)
			w.Header().Set("Surrogate-Control", maxAge10)
			res := f(r, w.Header(), &b)
			t.Stop()

			switch res.code {
			case http.StatusOK:
				b.WriteTo(w)
				t.Track()
				if t.Taken() > 500 {
					log.Printf("%s took %d ms to handle %s", id, t.Taken(), r.URL.Path)
				}

				mtrapp.StatusOK.Inc()
			case http.StatusBadRequest:
				w.Header().Set("Cache-Control", maxAge10)
				w.Header().Set("Surrogate-Control", maxAge86400)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(res.msg))
				mtrapp.StatusBadRequest.Inc()
			case http.StatusInternalServerError:
				w.Header().Set("Cache-Control", maxAge10)
				w.Header().Set("Surrogate-Control", maxAge86400)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(res.msg))
				mtrapp.StatusInternalServerError.Inc()
				log.Printf("500 serving GET %s %s", r.URL, res.msg)
			case http.StatusNotFound:
				w.Header().Set("Cache-Control", maxAge10)
				w.Header().Set("Surrogate-Control", maxAge10)
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(res.msg))
				mtrapp.StatusNotFound.Inc()
			default:
				w.Header().Set("Cache-Control", maxAge10)
				w.Header().Set("Surrogate-Control", maxAge86400)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(res.msg))
			}
		default:
			w.Header().Set("Cache-Control", maxAge10)
			w.Header().Set("Surrogate-Control", maxAge86400)
			w.Write([]byte("method not allowed"))
			return
		}
	}
}
