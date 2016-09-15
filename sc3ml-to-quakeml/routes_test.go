package main

import (
	wt "github.com/GeoNet/weft/wefttest"
	"testing"
	"net/http/httptest"
)

var testServer *httptest.Server

var routes = wt.Requests{
	{ID: wt.L(), URL: "/quakeml/1.2/2016p500086"},
	{ID: wt.L(), URL: "/quakeml-rt/1.2/2016p500086"},
	{ID: wt.L(), URL: "/csv/1.0.0/2016p500086/event"},
	{ID: wt.L(), URL: "/csv/1.0.0/2016p500086/picks"},
	{ID: wt.L(), URL: "/csv/1.0.0/2016p500086/event/picks"},

	// soh routes
	{ID: wt.L(), URL: "/soh"},
	{ID: wt.L(), URL: "/soh/up"},
}

func TestRoutes(t *testing.T) {
	setup(t)
	defer teardown()

	for _, r := range routes {
		if b, err := r.Do(testServer.URL); err != nil {
			t.Error(err)
			t.Error(string(b))
		}
	}

	if err := routes.DoCheckQuery(testServer.URL); err != nil {
		t.Error(err)
	}
}

func setup(t *testing.T) {
	testServer = httptest.NewServer(mux)
}

func teardown() {
	testServer.Close()
}