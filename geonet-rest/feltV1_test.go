//# Felt Reports
//
//##/felt
//
// Look up felt report information.
//
package main

import (
	"testing"
)

//## Felt Reports for a Quake.
//
// **GET /felt/report?publicID=(publicID)**
//
// Get Felt Reports for a quake.
//
//### Parameters
//
// * `publicID` - a valid quake identfier.
//
//### Example request:
//
// `/felt/report?publicID=2013p407387`
//
func TestReportsV1(t *testing.T) {
	// tests are done in routes and geojson test.  This is just a handle for the docs.
}
