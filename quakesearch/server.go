package main

import (
	"github.com/GeoNet/haz/database"
	"github.com/GeoNet/haz/msg"
	"github.com/GeoNet/log/logentries"
	"github.com/GeoNet/weft"
	"log"
	"net/http"
	"os"
)

var db database.DB

func init() {
	logentries.Init(os.Getenv("LOGENTRIES_TOKEN"))
	// TODO - we're not using LIBRATO anymore - please remove this and the env var.
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
			// TODO you're not versioning on accept that I can find so remove this.
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
