package main

import (
	"github.com/GeoNet/haz/database"
	"github.com/GeoNet/haz/msg"
	"net/http/httptest"
	"testing"
)

var (
	ts *httptest.Server
)

// setup starts a db connection and test server then inits an http client.
func setup(t *testing.T) {

	var err error
	database.DBUser = "hazard_w"
	tdb, err := database.InitPG()
	if err != nil {
		t.Fatal(err)
	}

	tdb.Check()

	_, err = tdb.Exec("delete from haz.quake")
	if err != nil {
		t.Fatal(err)
	}

	q := msg.ReadSC3ML07("etc/test/files/1542894.xml")
	if err != nil {
		t.Fatal(err)
	}

	err = tdb.SaveQuake(q)
	if err != nil {
		t.Fatal(err)
	}

	q = msg.ReadSC3ML07("etc/test/files/2190619.xml")
	if err != nil {
		t.Fatal(err)
	}

	err = tdb.SaveQuake(q)
	if err != nil {
		t.Fatal(err)
	}

	q = msg.ReadSC3ML07("etc/test/files/3366146.xml")
	if err != nil {
		t.Fatal(err)
	}

	err = tdb.SaveQuake(q)
	if err != nil {
		t.Fatal(err)
	}

	tdb.Close()

	// switch back to the correct user for the tests.
	// hazard_r can read haz and impact.
	database.DBUser = "hazard_r"
	db, err = database.InitPG()
	if err != nil {
		t.Fatal(err)
	}

	err = db.Ping()

	if err != nil {
		t.Fatal(err)
	}

	ts = httptest.NewServer(handler())

}

// teardown closes the db connection and  test server.  Defer this after setup() e.g.,
// ...
// setup()
// defer teardown()
func teardown() {
	ts.Close()
	db.Close()
}
