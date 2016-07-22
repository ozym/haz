package main

import (
	"log"
	"net/http"
	"time"
	_ "github.com/GeoNet/log/logentries"
)

var (
	client  *http.Client
	timeout = time.Duration(30 * time.Second)
)

func init() {
	mux.HandleFunc("/health", health)
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
