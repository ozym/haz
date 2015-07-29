package main

import (
	"encoding/json"
	"encoding/xml"
	"github.com/GeoNet/web"
	"github.com/GeoNet/web/api/apidoc"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	mlink   = "http://info.geonet.org.nz/m/view-rendered-page.action?abstractPageId="
	newsURL = "http://info.geonet.org.nz/createrssfeed.action?types=blogpost&spaces=conf_all&title=GeoNet+News+RSS+Feed&labelString%3D&excludedSpaceKeys%3D&sort=created&maxResults=10&timeSpan=500&showContent=true&publicFeed=true&confirm=Create+RSS+Feed"
)

var newsDoc = apidoc.Endpoint{
	Title:       "News",
	Description: `GeoNet news stories.`,
	Queries: []*apidoc.Query{
		newsD,
	},
}

// Feed is used for unmarshaling XML (from the GeoNet RSS news feed)
// and marshaling JSON
type Feed struct {
	Entries []Entry `xml:"entry" json:"feed"`
}

// Entry is used for unmarshaling XML and marshaling JSON.
// JSON tags with a - will not be include in the output.
type Entry struct {
	Title     string `xml:"title" json:"title"`
	Published string `xml:"published" json:"published"`
	Link      Link   `xml:"link" json:"-"`
	Id        string `xml:"id" json:"-"`
	Href      string `json:"link"`
	MHref     string `json:"mlink"`
}

// Link is used for unmarshaling XML.
type Link struct {
	Href string `xml:"href,attr"`
}

// unmarshalNews unmarshals the GeoNet News RSS XML.
func unmarshalNews(b []byte) (f Feed, err error) {
	err = xml.Unmarshal(b, &f)
	if err != nil {
		return f, err
	}

	// Copy the story link and make the link to the
	// mobile friendly version of the story.
	for i, _ := range f.Entries {
		f.Entries[i].Href = f.Entries[i].Link.Href
		f.Entries[i].MHref = mlink + strings.Split(f.Entries[i].Id, "-")[1]
	}

	return f, err
}

// /news/geonet

var newsD = &apidoc.Query{
	Accept:      web.V1JSON,
	Title:       "News",
	Description: " Returns a simple JSON version of the GeoNet News RSS feed.",
	Example:     "/news/geonet",
	ExampleHost: exHost,
	URI:         "/news/geonet",
	Params: map[string]template.HTML{
		"none": `no query parameters are required.`,
	},
	Props: map[string]template.HTML{
		"mlink":     "a link to a mobile version of the news story.",
		"link":      "a link to the news story.",
		"published": "the date the story was published",
		"title":     "the title of the story.",
	},
}

func news(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query()) != 0 {
		web.BadRequest(w, r, "incorrect number of query parameters.")
		return
	}

	j, err := fetchRSS(newsURL)
	if err != nil {
		web.ServiceUnavailable(w, r, err)
		return
	}

	w.Header().Set("Surrogate-Control", web.MaxAge300)

	web.Ok(w, r, &j)
}

func fetchRSS(url string) (b []byte, err error) {
	res, err := client.Get(url)
	defer res.Body.Close()
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	rss, err := unmarshalNews(body)
	if err != nil {
		return
	}

	b, err = json.Marshal(rss)

	return
}
