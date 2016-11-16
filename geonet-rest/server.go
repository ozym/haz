package main

import (
	"github.com/GeoNet/haz/database"
	"github.com/GeoNet/weft"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
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
	V2GeoJSON = "application/vnd.geo+json;version=2"
	V2JSON    = "application/json;version=2"
	protobuf  = "application/x-protobuf"
)

// These are for CAP format and Atom which is not versioned by Accept.
const (
	CAP  = "application/cap+xml"
	Atom = "application/xml"
)

const (
	ErrContent  = "text/plain; charset=utf-8"
	HtmlContent = "text/html; charset=utf-8"
)

var (
	db     database.DB
	client *http.Client
)

// main connects to the database, sets up request routing, and starts the http server.
func main() {
	var err error
	db, err = database.InitPG()
	if err != nil {
		log.Println("Problem with DB config.")
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Println("ERROR: problem pinging DB - is it up and contactable? 500s will be served")
	}

	// create an http client to share.
	timeout := time.Duration(5 * time.Second)
	client = &http.Client{
		Timeout: timeout,
	}

	log.Println("starting server")
	http.Handle("/", handler())
	log.Fatal(http.ListenAndServe(":"+os.Getenv("WEB_SERVER_PORT"), nil))
}

// handler creates a mux and wraps it with default handlers.  Separate function to enable testing.
func handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", router)
	return inbound(mux)
}

func inbound(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO this is a browser cache directive and does not make a lot
		// of sense for an API.
		w.Header().Set("Cache-Control", maxAge10)
		switch r.Method {
		case "GET":
			// Routing is based on Accept query parameters
			// e.g., version=1 in application/json;version=1
			// so caching must Vary based on Accept.
			w.Header().Set("Vary", "Accept")

			h.ServeHTTP(w, r)
		default:
			weft.Write(w, r, &weft.MethodNotAllowed)
			weft.MethodNotAllowed.Count()
			return
		}
	})
}
