// haz-twitter-consumer listens to an AWS SQS queue for Haz JSON messages and
// post it to twitter if it passed the given threshold.
package main

import (
	"github.com/GeoNet/cfg"
	"github.com/GeoNet/haz/msg"
	"github.com/GeoNet/haz/sqs"
	"github.com/GeoNet/haz/twitter"
	"github.com/GeoNet/log/logentries"
	"log"
)

//go:generate configer haz-twitter-consumer.json
var (
	config    = cfg.Load()
	idp       = msg.IdpQuake{}
	threshold float64
	ttr       twitter.Twitter
)

func init() {
	logentries.Init(config.Logentries.Token)
	msg.InitLibrato(config.Librato.User, config.Librato.Key, config.Librato.Source)
	threshold = config.Twitter.MinMagnitude
	config.SQS.MaxNumberOfMessages = 1
	config.SQS.VisibilityTimeout = 600
	config.SQS.WaitTimeSeconds = 20
}

type message struct {
	msg.Haz
}

func main() {
	rx, dx, err := sqs.InitRx(config.SQS)
	if err != nil {
		log.Fatalf("ERROR - problem creating SQS from config: %s", err)
	}

	ttr, err = twitter.Init(config.Twitter)
	if err != nil {
		log.Fatalf("ERROR: Twitter init error: %s", err.Error())
	}

	for {
		r := <-rx
		h := message{}
		h.Decode([]byte(r.Body))
		if !msg.Process(&h) {
			dx <- r.ReceiptHandle
		}
	}
}

func (m *message) Process() bool {
	switch {
	case m.Err() != nil:
		log.Println("WARN received errored message: " + m.Err().Error())
	case m.HeartBeat != nil:
		m.HeartBeat.RxLog()
	case m.Quake != nil:
		m.Quake.RxLog()
		return m.processTweet()
	}

	return false
}

func (m *message) processTweet() bool {
	if idp.Seen(*m.Quake) {
		log.Printf("%s already tweeted.", m.Quake.PublicID)
		return false
	}

	alert, message := m.Quake.AlertTwitter(threshold)
	if alert {
		log.Printf("Tweeting quake %s.", m.Quake.PublicID)
		err := ttr.PostTweet(message, m.Quake.Longitude, m.Quake.Latitude)

		if err != nil {
			m.SetErr(err)
			return true
		}

		idp.Add(*m.Quake)
	} else {
		log.Printf("quake %s not suitable for tweeting.", m.Quake.PublicID)
	}

	return false
}
