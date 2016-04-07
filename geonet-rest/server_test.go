package main

import (
	"github.com/GeoNet/haz/database"
	"github.com/GeoNet/haz/msg"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"net/http/httptest"
	"time"
)

const tolerance float64 = 0.0001

var measuredTest = []msg.Intensity{
	msg.Intensity{Source: "NZ.WEL",
		Longitude: 175.49,
		Latitude:  -40.2,
		Time:      time.Now().UTC(),
		MMI:       4,
	},
}

var reportedTest = []msg.Intensity{
	msg.Intensity{Source: "id1",
		Longitude: 176.49,
		Latitude:  -40.2,
		Time:      time.Now().UTC(),
		MMI:       4,
	},
	msg.Intensity{Source: "id2",
		Longitude: 176.49,
		Latitude:  -40.2,
		Time:      time.Now().UTC(),
		MMI:       5,
	},
	msg.Intensity{Source: "id3",
		Longitude: 176.49,
		Latitude:  -40.2,
		Time:      time.Now().UTC(),
		MMI:       6,
	},
}

var ts *httptest.Server

// setup starts a db connection and test server then inits an http client.
func setup() {
	// load some test data.  Needs a write user.
	var err error
	database.DBUser = "hazard_w"
	tdb, err := database.InitPG()
	if err != nil {
		log.Fatal(err)
	}

	tdb.Check()

	_, err = tdb.Exec("delete from haz.quake")
	if err != nil {
		log.Fatal(err)
	}

	_, err = tdb.Exec("delete from haz.quakeapi")
	if err != nil {
		log.Fatal(err)
	}

	_, err = tdb.Exec("delete from haz.quakehistory")
	if err != nil {
		log.Fatal(err)
	}

	q := msg.ReadSC3ML07("etc/test/files/2013p407387.xml")
	if err != nil {
		log.Fatal(err)
	}

	// stop the quake being deleted from haz.quakehistory and haz.quakeapi
	q.Time = time.Now().UTC()

	err = tdb.SaveQuake(q)
	if err != nil {
		log.Fatal(err)
	}

	q = msg.ReadSC3ML07("etc/test/files/2015p012816.xml")
	if err != nil {
		log.Fatal(err)
	}

	// stop the quake being deleted from haz.quakehistory and haz.quakeapi
	q.Time = time.Now().UTC()

	err = tdb.SaveQuake(q)
	if err != nil {
		log.Fatal(err)
	}

	tdb.Close()

	database.DBUser = "impact_w"
	tdb, err = database.InitPG()
	if err != nil {
		log.Fatal(err)
	}

	tdb.Check()

	_, err = tdb.Exec("delete from impact.intensity_measured")
	if err != nil {
		log.Fatal(err)
	}

	_, err = tdb.Exec("delete from impact.intensity_reported")
	if err != nil {
		log.Fatal(err)
	}

	for _, m := range measuredTest {
		_, err = tdb.Exec("select impact.add_intensity_measured($1, $2, $3, $4, $5)", m.Source, m.Longitude, m.Latitude, m.Time, m.MMI)
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, m := range reportedTest {
		_, err = tdb.Exec("select impact.add_intensity_reported($1, $2, $3, $4, $5, $6)", m.Source, m.Longitude, m.Latitude, m.Time, m.MMI, m.Comment)
		if err != nil {
			log.Fatal(err)
		}
	}

	tdb.Close()

	// switch back to the correct user for the tests.
	// hazard_r can read haz and impact.
	database.DBUser = "hazard_r"
	db, err = database.InitPG()
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	ts = httptest.NewServer(handler())

	client = &http.Client{}
}

// teardown closes the db connection and  test server.  Defer this after setup() e.g.,
// ...
// setup()
// defer teardown()
func teardown() {
	ts.Close()
	db.Close()
}

// Valid is used to hold the response from GeoJSON validation.
type Valid struct {
	Status string
}
