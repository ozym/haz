package main

import (
	"database/sql"
	"github.com/GeoNet/cfg"
	"github.com/GeoNet/log/logentries"
	"github.com/GeoNet/web"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

//go:generate configer geonet-rest.json
var (
	config = cfg.Load()
	db     *sql.DB
	client *http.Client
)

var header = web.Header{
	Cache:     web.MaxAge10,
	Surrogate: web.MaxAge10,
	Vary:      "Accept",
}

func init() {
	logentries.Init(config.Logentries.Token)
	web.InitLibrato(config.Librato.User, config.Librato.Key, config.Librato.Source)
}

// main connects to the database, sets up request routing, and starts the http server.
func main() {
	var err error
	db, err = sql.Open("postgres", config.DataBase.Postgres())
	if err != nil {
		log.Println("Problem with DB config.")
		log.Fatal(err)
	}
	defer db.Close()

	db.SetMaxIdleConns(config.DataBase.MaxIdleConns)
	db.SetMaxOpenConns(config.DataBase.MaxOpenConns)

	if err = db.Ping(); err != nil {
		log.Println("ERROR: problem pinging DB - is it up and contactable? 500s will be served")
	}

	// create an http client to share.
	timeout := time.Duration(5 * time.Second)
	client = &http.Client{
		Timeout: timeout,
	}

	http.Handle("/", handler())
	log.Fatal(http.ListenAndServe(":"+config.WebServer.Port, nil))
}

// handler creates a mux and wraps it with default handlers.  Seperate function to enable testing.
func handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", router)
	return header.GetGzip(mux)
}
