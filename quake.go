package msg

import (
	"bytes"
	"fmt"
	"github.com/GeoNet/Golang-Ellipsoid/ellipsoid"
	"log"
	"math"
	"text/template"
	"time"
)

var (
	geo      ellipsoid.Ellipsoid
	alertAge = time.Duration(-60) * time.Minute
	nz       *time.Location
	t        = template.Must(template.New("eqNews").Parse(eqNews))
)

const (
	dutyTime    = "3:04 PM, 02/01/2006 MST"
	eqNewsNow   = "Mon 2 Jan 2006 at 3:04 pm"
	eqNewsUTC   = "2006/01/02 at 15:04:05"
	eqNewsLocal = "(MST):      Monday 2 Jan 2006 at 3:04 pm"
)

const eqNews = `                PRELIMINARY EARTHQUAKE REPORT

                      GeoNet Data Centre
                         GNS Science
                   Lower Hutt, New Zealand
                   http://www.geonet.org.nz

        Report Issued at: {{.Now}}


A likely felt earthquake has been detected by GeoNet; this is PRELIMINARY information only:

        Public ID:              {{.Q.PublicID}}
        Universal Time:         {{.UT}}
        Local Time {{.LT}}
        Latitude, Longitude:    {{.LL}}
        Location:               {{.Location}}
        Intensity:              {{.Intensity}} (MM{{.MMI}})
        Depth:                  {{ printf "%.f"  .Q.Depth}} km
        Magnitude:              {{ printf "%.1f"  .Q.Magnitude}}

Check for the LATEST information at http://www.geonet.org.nz/quakes/{{.Q.PublicID}}
`

type eqNewsD struct {
	Q         *Quake
	MMI       int
	Location  string
	Now       string
	TZ        string // timezone for the quake.
	UT        string // quake time in UTC
	LT        string // quake in local time
	LL        string // lon lat string
	Intensity string // word version of MMI
}

func init() {
	geo = ellipsoid.Init("WGS84", ellipsoid.Degrees, ellipsoid.Kilometer, ellipsoid.LongitudeIsSymmetric, ellipsoid.BearingNotSymmetric)
	var err error
	nz, err = time.LoadLocation("Pacific/Auckland")
	if err != nil {
		log.Println("Error loading TZNZ carrying on with UTC")
		nz = time.UTC
	}
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
	DepthType             string
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
	Site                  string
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

// Publish returns true if the quake is suitable for publishing.
// site is either 'primary' or 'backup'.
func (q *Quake) Publish() bool {
	if q.err != nil {
		return false
	}

	p := true
	switch q.Site {
	case "primary", "":
		if q.Status() == "automatic" && !(q.Depth >= 0.1 && q.AzimuthalGap <= 320.0 && q.MinimumDistance <= 2.5) {
			p = false
			q.SetErr(fmt.Errorf("Not publising automatic quake %s with poor quality from primary site.", q.PublicID))
		}
	case "backup":
		if q.Status() == "automatic" {
			p = false
			q.SetErr(fmt.Errorf("Not publising unreviewed quake %s from backup site.", q.PublicID))
		}
	}
	return p
}

// AlertDuty returns alert = true and message formated if the quake is suitable for alerting the
// duty people, alert = false and empty message if not.
func (q *Quake) AlertDuty() (alert bool, message string) {
	if q.Err() != nil {
		return
	}

	if !q.AlertQuality() {
		return
	}

	mmi := q.MMI()

	if mmi >= 6 || q.Magnitude >= 4.5 {
		alert = true

		c, d, b, err := q.Closest()
		if err != nil {
			q.SetErr(err)
			return
		}

		// Eq Rpt: MAG 5.0, MM7, DEP 10, LOC 105 km N of White Island, TIME 08:33 AM, 26/02/2015
		message = fmt.Sprintf("Eq Rpt: MAG %.1f, MM%d, DEP %.f, LOC %s %s of %s, TIME %s",
			q.Magnitude,
			int(mmi),
			q.Depth,
			Distance(d),
			Compass(b),
			c.Name,
			q.Time.In(nz).Format(dutyTime))
	}

	return
}

// AlertPIM returns alert = true and message formated if the quake is suitable for alerting the
// Pubilc Information people, alert = false and empty message if not.
func (q *Quake) AlertPIM() (alert bool, message string) {
	if q.Err() != nil {
		return
	}

	if !q.AlertQuality() {
		return
	}

	if q.Magnitude >= 6.0 {
		alert = true

		mmi := q.MMI()

		c, d, b, err := q.Closest()
		if err != nil {
			q.SetErr(err)
			return
		}

		// Eq Rpt: MAG 5.0, MM7, DEP 10, LOC 105 km N of White Island, TIME 08:33 AM, 26/02/2015
		message = fmt.Sprintf("Eq Rpt: MAG %.1f, MM%d, DEP %.f, LOC %s %s of %s, TIME %s",
			q.Magnitude,
			int(mmi),
			q.Depth,
			Distance(d),
			Compass(b),
			c.Name,
			q.Time.In(nz).Format(dutyTime))
	}

	return
}

func (q *Quake) AlertEqNews() (alert bool, subject, body string) {
	if q.Err() != nil {
		return
	}

	if !q.AlertQuality() {
		return
	}

	mmi := q.MMI()

	c, d, b, err := q.Closest()
	if err != nil {
		q.SetErr(err)
		return
	}

	mmid := MMIDistance(d, q.Depth, mmi)

	if mmi >= 7.0 || mmid >= 3.5 {
		alert = true

		// NZ EQ: M3.5, weak intensity, 5km deep, 20 km N of Reefton
		subject = fmt.Sprintf("NZ EQ: M%.1f, %s intensity, %.fkm deep, %s %s of %s",
			q.Magnitude,
			MMIIntensity(mmi),
			q.Depth,
			Distance(d),
			Compass(b),
			c.Name)

	}

	buf := new(bytes.Buffer)

	err = t.ExecuteTemplate(buf, "eqNews", &eqNewsD{
		Q:         q,
		MMI:       int(mmi),
		Location:  fmt.Sprintf("%s %s of %s", Distance(d), Compass(b), c.Name),
		Now:       time.Now().In(nz).Format(eqNewsNow),
		UT:        q.Time.Format(eqNewsUTC),
		LT:        q.Time.In(nz).Format(eqNewsLocal),
		LL:        q.eqNewsLonLat(),
		Intensity: MMIIntensity(mmi),
	})
	if err != nil {
		q.SetErr(err)
		alert = false
		return
	}

	body = buf.String()

	return

}

func Distance(km float64) string {
	s := "Within 5 km of"

	d := math.Floor(km / 5.0)
	if d > 0 {
		s = fmt.Sprintf("%.f km", d*5)
	}
	return s
}

func (q *Quake) eqNewsLonLat() string {
	var lon, lat string

	switch q.Longitude < 0.0 {
	case true:
		lon = fmt.Sprintf("%.2fW", q.Longitude*-1.0)
	case false:
		lon = fmt.Sprintf("%.2fE", q.Longitude)
	}

	switch q.Latitude < 0.0 {
	case true:
		lat = fmt.Sprintf("%.2fS", q.Latitude*-1.0)
	case false:
		lat = fmt.Sprintf("%.2fN", q.Latitude)
	}

	// 41.94S, 171.86E
	return lat + ", " + lon
}
