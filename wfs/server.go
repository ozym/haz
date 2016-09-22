package main

import (
	"github.com/GeoNet/haz/database"
	"github.com/GeoNet/weft"
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
	mux.HandleFunc("/", router)
	return inbound(mux)
}

func inbound(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			h.ServeHTTP(w, r)
		default:
			weft.Write(w, r, &weft.MethodNotAllowed)
			weft.MethodNotAllowed.Count()
			return
		}
	})
}
