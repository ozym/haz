package main

import (
	"fmt"
	"github.com/GeoNet/haz/msg"
	"log"
	"math"
	"sort"
	"strings"
	"text/template"
	"time"
)

type capQuakeT struct {
	References []string
	Quake      msg.Quake
	Intensity  string
	Closest    msg.LocalityQuake
	Status     string // quake status
	Localities []msg.LocalityQuake
	ID         string // CAP message ID
}

type capAtomEntry struct {
	ID       string // unique and perm URI for the entry
	Title    string
	Updated  time.Time
	Summary  string
	HrefCAP  string // link to the CAP formated info
	HrefHTML string // link to web page with futher info.
}

type capAtomFeed struct {
	ID      string    // univerally unique URI for the feed.
	Title   string    //
	Updated time.Time // last time feed content modified.
	Link    string    // link to the feed itself.
	Entries []capAtomEntry
}

const (
	displayTime = "Mon Jan 2 2006 3:04 PM (MST)"
)

var (
	capTemplates = template.Must(template.New("").Funcs(funcMap).ParseGlob("assets/tmpl/cap*.tmpl"))
	nz           *time.Location
	expire       = time.Duration(48) * time.Hour
)

func init() {
	var err error
	nz, err = time.LoadLocation("Pacific/Auckland")
	if err != nil {
		log.Println("Error loading TZNZ carrying on with UTC")
		nz = time.UTC
	}
}

var funcMap = template.FuncMap{
	"capTime": func(t time.Time) string {
		return t.In(nz).Format(time.RFC3339)
	},
	"atomTime": func(t time.Time) string {
		return t.Format(time.RFC3339)
	},
	"nzTime": func(t time.Time) string {
		return t.In(nz).Format(displayTime)
	},
	"capID": func(t time.Time) string {
		return fmt.Sprintf("%d", t.Unix())
	},
	"expires": func(t time.Time) string {
		return t.Add(expire).In(nz).Format(time.RFC3339)
	},
	"feltIn": func(localities []msg.LocalityQuake) string {
		var s string

		for _, l := range localities {
			s = s + fmt.Sprintf("%s, ", l.Locality.Name)
		}

		return strings.Trim(s, ", ")
	},
	"radius": func(localities []msg.LocalityQuake) (r float64) {
		for _, l := range localities {
			if l.Distance > r {
				r = l.Distance
			}
		}
		return
	},
	"area": func(localities []msg.LocalityQuake) string {
		var area string

		sort.Sort(msg.ByDistance(localities))

		for _, l := range localities {
			if l.Distance < 5 {
				area = area + fmt.Sprintf("Within 5 km of %s, ", l.Locality.Name)
			} else {
				area = area + fmt.Sprintf("%s %s of %s, ", msg.Distance(l.Distance), compass(l.Bearing), l.Locality.Name)
			}
		}

		return strings.Trim(area, ", ")
	},
	"distance": msg.Distance,
	"location": func(d, b float64, name string) string {
		if d < 5 {
			return "Within 5 km of " + name
		}

		return fmt.Sprintf("%.f km %s of %s", math.Floor(d/5.0)*5, compass(b), name)
	},
	"compass": compass,
	"severity": func(mmi float64) string {
		switch {
		case mmi >= 8:
			return `Extreme`
		case mmi >= 7:
			return `Severe`
		case mmi >= 6:
			return `Moderate`
		case mmi < 6:
			return "Minor"
		default:
			return "Unknown"
		}
	},
	"certainty": func(q msg.Quake, status string) string {
		switch {
		case status == `reviewed`:
			return `Observed`
		case status == `deleted`:
			return `Unlikely`
		case q.UsedPhaseCount >= 20 && q.MagnitudeStationCount >= 10:
			return `Likely`
		case status == "automatic" && (q.UsedPhaseCount < 20 || q.MagnitudeStationCount < 10):
			return `Possible`
		default:
			return `Unknown`
		}
	},
	"msgType": func(status string, references []string) string {
		switch {
		case status == `deleted`:
			return `Cancel`
		case len(references) == 0:
			return `Alert`
		default:
			return `Update`
		}
	},
}

func compass(b float64) string {
	switch {
	case b >= 337.5 && b <= 360:
		return "N"
	case b >= 0 && b <= 22.5:
		return "N"
	case b > 22.5 && b < 67.5:
		return "NE"
	case b >= 67.5 && b <= 112.5:
		return "E"
	case b > 112.5 && b < 157.5:
		return "SE"
	case b >= 157.5 && b <= 202.5:
		return "S"
	case b > 202.5 && b < 247.5:
		return "SW"
	case b >= 247.5 && b <= 292.5:
		return "W"
	case b > 292.5 && b < 337.5:
		return "NW"
	default:
		return "N"
	}
}

func (c capQuakeT) MsgType() string {
	switch {
	case c.Status == `deleted`:
		return `Cancel`
	case len(c.References) == 0:
		return `Alert`
	default:
		return `Update`
	}
}
