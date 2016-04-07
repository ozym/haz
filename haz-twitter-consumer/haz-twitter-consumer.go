// haz-twitter-consumer listens to an AWS SQS queue for Haz JSON messages and
// post it to twitter if it passed the given threshold.
package main

import (
	"github.com/GeoNet/haz/msg"
	"github.com/GeoNet/haz/sqs"
	"github.com/GeoNet/haz/twitter"
	"github.com/GeoNet/log/logentries"
	"log"
	"os"
	"strconv"
)

//go:generate configer haz-twitter-consumer.json
var (
	idp       = msg.IdpQuake{}
	ttr       twitter.Twitter
	threshold float64
)

func init() {
	logentries.Init(os.Getenv("LOGENTRIES_TOKEN"))
	msg.InitLibrato(os.Getenv("LIBRATO_USER"), os.Getenv("LIBRATO_KEY"), os.Getenv("LIBRATO_SOURCE"))
	sqs.MaxNumberOfMessages = 1
	sqs.VisibilityTimeout = 600
	sqs.WaitTimeSeconds = 20
	var err error
	if threshold, err = strconv.ParseFloat(os.Getenv("TWITTER_THRESHOLD"), 64); err != nil {
		log.Fatalln("TWITTER_THRESHOLD format error: ", err.Error())
	}
}

type message struct {
	msg.Haz
}

func main() {
	rx, dx, err := sqs.InitRx()
	if err != nil {
		log.Fatalf("ERROR - problem creating SQS from config: %s", err)
	}

	ttr, err = twitter.Init()
	if err != nil {
		log.Fatalf("ERROR: Twitter init error: %s", err.Error())
	}

	log.Printf("Twitter magnitude threshold %.1f", threshold)

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
