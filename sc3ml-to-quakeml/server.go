package main

import (
	"log"
	"net/http"
	"time"
	"github.com/GeoNet/weft"
	"bytes"
)

var (
	client  *http.Client
	timeout = time.Duration(30 * time.Second)
)

func init() {
	mux.HandleFunc("/", weft.MakeHandlerPage(home))
	mux.HandleFunc("/health", health)

	// routes for balancers and probes.
	mux.HandleFunc("/soh/up", http.HandlerFunc(up))
	mux.HandleFunc("/soh", http.HandlerFunc(soh))

}

func main() {
	log.Println("starting server")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

/*
health does not require auth - for use with AWS EB load balancer checks.
*/
func health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

// up is for testing that the app has started e.g., for with load balancers.
// It indicates the app is started.  It may still be serving errors.
// Not useful for inclusion in app metrics so weft not used.
func up(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
		w.Header().Set("Surrogate-Control", "max-age=86400")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write([]byte("<html><head></head><body>up</body></html>"))
	log.Print("up ok")
}

// soh is for external service probes.
// writes a service unavailable error to w if the service is not working.
// Not useful for inclusion in app metrics so weft not used.
func soh(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
		w.Header().Set("Surrogate-Control", "max-age=86400")
		w.WriteHeader(http.StatusBadRequest)
		log.Print("out")
		return
	}

	// TODO - is there a resource to test?
	// should we fetch SC3ML from S3 here?
	// would be a lot of overhead for a probe.

	w.Write([]byte("<html><head></head><body>ok</body></html>"))
	log.Print("soh ok")
}

func home(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
		return res
	}

	return &weft.NotFound
}

