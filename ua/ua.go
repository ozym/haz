package ua

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/GeoNet/cfg"
	"io/ioutil"
	"net/http"
)

const api = "https://go.urbanairship.com/api/push"

// UA's push object JSON
type pushData struct {
	Audience     audience     `json:"audience,omitempty"`
	Notification notification `json:"notification,omitempty"`
	DeviceTypes  []string     `json:"device_types,omitempty"`
}

type audience struct {
	DeviceToken []string `json:"device_token,omitempty"`
	Tag         []string `json:"tag,omitempty"`
	ApID        []string `json:"apid,omitempty"`
}

type notification struct {
	Ios     ios     `json:"ios,omitempty"`
	Android android `json:"android,omitempty"`
}

type ios struct {
	Badge            string `json:"badge,omitempty"`
	ContentAvailable bool   `json:"content-available,omitempty"`
	Extra            extra  `json:"extra"`
	Alert            string `json:"alert"`
	Sound            string `json:"sound,omitempty"`
}

type android struct {
	Alert string `json:"alert"`
	Extra extra  `json:"extra"`
}

type extra struct {
	PublicID string `json:"publicid"`
}

// end UA's push object JSON

type Client struct {
	h *http.Client
	c *cfg.UA
}

func Init(c *cfg.UA) *Client {
	a := &Client{
		h: &http.Client{},
		c: c,
	}

	return a
}

func (a *Client) Push(publicID string, message string, tags []string) (err error) {
	e := extra{
		PublicID: publicID,
	}
	i := ios{
		Badge:            "+1",
		ContentAvailable: true,
		Extra:            e,
		Alert:            message,
		Sound:            "default",
	}
	an := android{
		Alert: message,
		Extra: e,
	}
	n := notification{
		Ios:     i,
		Android: an,
	}
	au := audience{
		Tag: tags,
	}
	data := pushData{
		Audience:     au,
		Notification: n,
		DeviceTypes:  []string{"ios", "android"},
	}

	b, err := json.Marshal(data)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", api, bytes.NewBuffer(b))
	if err != nil {
		return
	}

	req.SetBasicAuth(a.c.AppKey, a.c.AppMasterSecret)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/vnd.urbanairship+json; version=3;")

	var res *http.Response
	res, err = a.h.Do(req)
	if err == nil {
		defer res.Body.Close()

		switch res.StatusCode {
		case 200, 202:
			return // success, we're out
		default:
			respString, _ := ioutil.ReadAll(res.Body)
			err = fmt.Errorf("Response error message from UA: %s", string(respString))
			return
		}
	}

	return
}
