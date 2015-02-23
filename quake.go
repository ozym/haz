package msg

import (
	"fmt"
	"github.com/StefanSchroeder/Golang-Ellipsoid/ellipsoid"
	"log"
	"math"
	"time"
)

var geo ellipsoid.Ellipsoid
var alertAge = time.Duration(-60) * time.Minute

func init() {
	geo = ellipsoid.Init("WGS84", ellipsoid.Degrees, ellipsoid.Kilometer, ellipsoid.LongitudeIsSymmetric, ellipsoid.BearingNotSymmetric)
}

type Quake struct {
	PublicID              string
	Type                  string
	AgencyID              string
	ModificationTime      time.Time
	Time                  time.Time
	Latitude              float64
	Longitude             float64
	Depth                 float64
	MethodID              string
	EarthModelID          string
	EvaluationMode        string
	EvaluationStatus      string
	UsedPhaseCount        int
	UsedStationCount      int
	StandardError         float64
	AzimuthalGap          float64
	MinimumDistance       float64
	Magnitude             float64
	MagnitudeUncertainty  float64
	MagnitudeType         string
	MagnitudeStationCount int
	err                   error
}

// Status returns the public status for the Quake referred to by q.
// Returns 'error' if q.Err() is not nil.
func (q *Quake) Status() string {
	if q.err != nil {
		return "error"
	}

	switch {
	case q.Type == "not existing":
		return "deleted"
	case q.Type == "duplicate":
		return "duplicate"
	case q.EvaluationMode == "manual":
		return "reviewed"
	case q.EvaluationStatus == "confirmed":
		return "reviewed"
	default:
		return "automatic"
	}
}

func (q *Quake) Err() error {
	return q.err
}

func (q *Quake) SetErr(err error) {
	q.err = err
}

func (q *Quake) RxLog() {
	if q.err != nil {
		return
	}

	log.Printf("Received quake %s", q.PublicID)
}

func (q *Quake) TxLog() {
	if q.err != nil {
		return
	}

	log.Printf("Sending quake %s", q.PublicID)
}

// MMI calculates the maximum Modificed Mercalli Intensity for the quake.
func (q *Quake) MMI() float64 {
	if q.err != nil {
		return -1. - 0
	}

	var w, m float64
	d := math.Abs(q.Depth)
	rupture := d

	if d < 100 {
		w = math.Min(0.5*math.Pow(10, q.Magnitude-5.39), 30.0)
		rupture = math.Max(d-0.5*w*0.85, 0.0)
	}

	if d < 70.0 {
		m = 4.40 + 1.26*q.Magnitude - 3.67*math.Log10(rupture*rupture*rupture+1634.691752)/3.0 + 0.012*d + 0.409
	} else {
		m = 3.76 + 1.48*q.Magnitude - 3.50*math.Log10(rupture*rupture*rupture)/3.0 + 0.0031*d
	}

	if m < 3.0 {
		m = -1.0
	}

	return m
}

// MMIDistance calculates the MMI at distance for New Zealand.  Distance and depth are in km.
func MMIDistance(distance, depth, mmi float64) float64 {
	// Minimum depth of 5 for numerical instability.
	d := math.Max(math.Abs(depth), 5.0)
	s := math.Hypot(d, distance)

	return math.Max(mmi-1.18*math.Log(s/d)-0.0044*(s-d), -1.0)
}

// MMIIntensity returns the string describing mmi.
func MMIIntensity(mmi float64) string {
	switch {
	case mmi >= 7:
		return "severe"
	case mmi >= 6:
		return "strong"
	case mmi >= 5:
		return "moderate"
	case mmi >= 4:
		return "light"
	case mmi >= 3:
		return "weak"
	default:
		return "unnoticeable"
	}
}

// Closest returns the New Zealand Locality closest to the quake. Distance is from the
// Quake to the Locality in km.  Bearing is from the Locality to the Quake in degrees.
func (q *Quake) Closest() (locality Locality, distance float64, bearing float64, err error) {
	if q.err != nil {
		err = q.err
		return
	}

	distance = 20000.0

	for _, l := range localities {
		d, b := geo.To(l.Latitude, l.Longitude, q.Latitude, q.Longitude)
		if d < distance {
			distance = d
			locality = l
			bearing = b
		}
	}

	// ensure larger locality when distant quake.
	if distance > 300 && locality.size >= 2 {
		distance = 20000

		for _, l := range localities {
			if l.size == 0 || l.size == 1 {
				d, b := geo.To(l.Latitude, l.Longitude, q.Latitude, q.Longitude)
				if d < distance {
					distance = d
					locality = l
					bearing = b
				}
			}
		}
	}

	return locality, distance, bearing, nil
}

// Returns true of the Quake is of high enough quality to consider for alerting.
//  false if not.  If false Quake.Err() is also set.
func (q *Quake) AlertQuality() bool {
	if q.err != nil {
		return false
	}

	switch {
	case q.Status() == "deleted":
		q.err = fmt.Errorf("%s status deleted not suitable for alerting.", q.PublicID)
		return false
	case q.Status() == "duplicate":
		q.err = fmt.Errorf("%s status duplicate not suitable for alerting.", q.PublicID)
		return false
	case q.Status() == "automatic" && (q.UsedPhaseCount < 20 || q.MagnitudeStationCount < 10):
		q.err = fmt.Errorf("%s unreviewed with %d phases and %d magnitudes not suitable for alerting.", q.PublicID, q.UsedPhaseCount, q.MagnitudeStationCount)
		return false
	case q.Time.Before(time.Now().UTC().Add(alertAge)):
		q.err = fmt.Errorf("%s to old for alerting", q.PublicID)
		return false
	}

	return true
}
