package main

import (
	"github.com/GeoNet/haz/database"
	"github.com/GeoNet/haz/msg"
	"github.com/GeoNet/log/logentries"
	"github.com/GeoNet/web"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

//go:generate configer geonet-rest.json
var (
	db     database.DB
	client *http.Client
)

var header = web.Header{
	Cache:     web.MaxAge10,
	Surrogate: web.MaxAge10,
	Vary:      "Accept",
}

func init() {
	logentries.Init(os.Getenv("LOGENTRIES_TOKEN"))
	msg.InitLibrato(os.Getenv("LIBRATO_USER"), os.Getenv("LIBRATO_KEY"), os.Getenv("LIBRATO_SOURCE"))
}

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

	http.Handle("/", handler())
	log.Fatal(http.ListenAndServe(":"+os.Getenv("WEB_SERVER_PORT"), nil))
}

// handler creates a mux and wraps it with default handlers.  Seperate function to enable testing.
func handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", router)
	return header.GetGzip(mux)
}
