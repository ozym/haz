package main

import (
	"database/sql"
	"github.com/GeoNet/haz/msg"
	_ "github.com/lib/pq"
	"log"
	"testing"
	"time"
)

func setup() {
	var err error
	db, err = sql.Open("postgres", config.DataBase.Postgres())
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}
}

func teardown() {
	db.Close()
}

func TestInsert(t *testing.T) {
	setup()
	defer teardown()

	q := msg.Quake{
		PublicID:              "2015p278423",
		Time:                  time.Now().UTC(),
		Latitude:              -37.92257397,
		Longitude:             178.3544071,
		Depth:                 9.62890625,
		EvaluationStatus:      "automatic",
		UsedPhaseCount:        25,
		AzimuthalGap:          180,
		MinimumDistance:       2.4,
		Magnitude:             6.0,
		MagnitudeStationCount: 12,
		Site: "primary",
	}

	h := msg.Haz{
		Quake: &q,
	}

	m := message{
		h,
	}

	m.saveQuake()

	if m.Err() != nil {
		t.Errorf("TestInsert failed to insert database:", m.Err().Error())
	}

}
