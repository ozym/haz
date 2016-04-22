package main

import (
	_ "github.com/GeoNet/log/logentries"
	"github.com/NYTimes/gziphandler"

	"github.com/GeoNet/haz/database"
	"log"
	"net/http"
	"os"
)

var db database.DB

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

	log.Println("starting server")
	http.Handle("/", handler())
	log.Fatal(http.ListenAndServe(":"+os.Getenv("WEB_SERVER_PORT"), nil))
}

func handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", toHandler(router))
	return gziphandler.GzipHandler(mux)
}
