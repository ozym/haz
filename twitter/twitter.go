package twitter

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/GeoNet/cfg"
	"net/url"
	"strconv"
)

type Twitter struct {
	api *anaconda.TwitterApi
}

func Init(c *cfg.Twitter) (Twitter, error) {
	anaconda.SetConsumerKey(c.ConsumerKey)
	anaconda.SetConsumerSecret(c.ConsumerSecret)
	api := anaconda.NewTwitterApi(c.OAuthToken, c.OAuthSecret)

	var err error

	t := Twitter{
		api: api,
	}

	if api.Credentials == nil {
		err = fmt.Errorf("Credentials are invalid")
	} else {

		err = nil
	}
	return t, err
}

// Note: Twitter will rejects messages longer than 140 chars.
func (a *Twitter) PostTweet(message string, longitude float64, latitude float64) (err error) {
	if a.api.Credentials == nil {
		return fmt.Errorf("Credentials are invalid, cannot post.")
	}

	v := url.Values{}
	v.Set("long", strconv.FormatFloat(longitude, 'f', -1, 64))
	v.Set("lat", strconv.FormatFloat(latitude, 'f', -1, 64))

	_, err = a.api.PostTweet(message, v)

	return err
}
