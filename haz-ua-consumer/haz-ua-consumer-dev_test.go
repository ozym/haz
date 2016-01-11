// +build devtest

// Note: This test sends a push notification by tags.
// Carefully configure UA's key/secret to make sure that only testers will receive the notification.
package main

import (
	"github.com/GeoNet/haz/msg"
	"testing"
	"time"
)

func setup() {
}

func teardown() {
}

func TestUAPush(t *testing.T) {
	setup()
	defer teardown()

	// NOTICE: Always test with a small earthquake!
	// If possible, set the location to Arthur's Pass.
	q := msg.Quake{
		PublicID:              "2015p999999",
		Time:                  time.Now().UTC(),
		Latitude:              -42.850020,
		Longitude:             171.696670,
		Depth:                 9.62890625,
		EvaluationStatus:      "automatic",
		UsedPhaseCount:        25,
		AzimuthalGap:          180,
		MinimumDistance:       2.4,
		Magnitude:             5.0,
		MagnitudeStationCount: 12,
		Site: "primary",
	}

	m := message{
		msg.Haz{
			Quake: &q,
		},
	}

	if false != m.processPush() {
		t.Errorf("TestUAPush failed")
	}

}
