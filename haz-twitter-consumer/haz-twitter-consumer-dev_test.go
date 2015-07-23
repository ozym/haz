// +build devtest

package main

import (
	"github.com/GeoNet/haz/msg"
	"github.com/GeoNet/haz/twitter"
	"log"
	"testing"
	"time"
)

func setup() {
	var err error
	ttr, err = twitter.Init(config.Twitter)
	if err != nil {
		log.Fatalf("ERROR: Twitter init error: %s", err.Error())
	}
}

func teardown() {
}

func TestTweet(t *testing.T) {
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

	m := message{
		msg.Haz{
			Quake: &q,
		},
	}

	if false != m.processTweet() {
		t.Errorf("TestTweet failed")
	}

	// test "Palmerston North", this will cause trucate happen
	m.Quake.PublicID = "2015p278424"
	m.Quake.Longitude = 175.62
	m.Quake.Latitude = -40.37

	if false != m.processTweet() {
		t.Errorf("TestTweet failed")
	}
}
