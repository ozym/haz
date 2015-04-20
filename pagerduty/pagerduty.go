package pagerduty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/GeoNet/cfg"
	"log"
	"net/http"
	"time"
)

const api = "https://events.pagerduty.com/generic/2010-04-15/create_event.json"

var retry = time.Duration(30) * time.Second

type trigger struct {
	ServiceKey  string `json:"service_key"`
	EventType   string `json:"event_type"`
	Description string `json:"description"`
	IncidentKey string `json:"incident_key"`
}

type Client struct {
	apiToken string
	h        *http.Client
}

func Init(c *cfg.PagerDuty) *Client {
	a := &Client{
		apiToken: c.ApiToken,
		h:        &http.Client{},
	}

	return a
}

// Trigger an event via the PagerDuty api.  IncidentKey is used for de-duping.
// https://developer.pagerduty.com/documentation/integration/events/trigger
// If an error is encountered then creating the incident is attempted retries more times with
// a pause of 30s between each attempt.  retries can be 0 to attempt publishing only once.
// Anything other than a 200 response from the API is treated as an error.
func (a *Client) Trigger(c *cfg.PagerDuty, Description, IncidentKey string, retries int) (err error) {
	t := trigger{
		ServiceKey:  c.ServiceKey,
		Description: Description,
		IncidentKey: IncidentKey,
		EventType:   "trigger",
	}

	b, err := json.Marshal(t)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", api, bytes.NewBuffer(b))
	if err != nil {
		return
	}

	req.Header.Add("Authorization", "Token token="+a.apiToken)
	req.Header.Add("Content-Type", "application/json")

	cnt := 0
	var res *http.Response
	for {
		res, err = a.h.Do(req)
		if err == nil {
			defer res.Body.Close()

			// Return for success or unrecoverable errors.  Other errors (including network) we can back
			// off and try again for.
			// https://developer.pagerduty.com/documentation/rest/errors
			switch res.StatusCode {
			case 200:
				return // success, we're out
			case 400:
				err = fmt.Errorf("400 from PagerDuty: Caller provided invalid arguments.")
				return
			case 401:
				err = fmt.Errorf("401 forbidden  from PagerDuty")
				return
			case 404:
				err = fmt.Errorf("404 not found from PagerDuty")
				return
			}
		}
		if cnt >= retries {
			if err == nil {
				err = fmt.Errorf("ran out of retries sending to PagerDuty.")
			}
			break
		}
		cnt++

		log.Println("WARN error sending to PagerDuty.  Sleeping and trying again")
		time.Sleep(retry)
	}

	return
}
