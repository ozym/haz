// haz-ua-consumer listens to an AWS SQS queue for Haz JSON messages and
// generate tags for push message subscribers, then send it to UA to push message out.
package main

import (
	"github.com/GeoNet/cfg"
	"github.com/GeoNet/haz/msg"
	"github.com/GeoNet/haz/sqs"
	"github.com/GeoNet/haz/ua"
	"github.com/GeoNet/log/logentries"
	"log"
)

//go:generate configer haz-ua-consumer.json
var (
	config = cfg.Load()
	idp    = msg.IdpQuake{}
	uac    *ua.Client
)

func init() {
	logentries.Init(config.Logentries.Token)
	msg.InitLibrato(config.Librato.User, config.Librato.Key, config.Librato.Source)
	uac = ua.Init(config.UA)
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
		return m.processPush()
	}

	return false
}

func (m *message) processPush() bool {
	if idp.Seen(*m.Quake) {
		log.Printf("%s already pushed.", m.Quake.PublicID)
		return false
	}

	message, tags := m.Quake.AlertUAPush()
	if tags == nil {
		log.Printf("Quake %s didn't produce any tag.", m.Quake.PublicID)
		return false
	}

	log.Printf("Sending quake %s with %d tags to UA.", m.Quake.PublicID, len(tags))
	err := uac.Push(m.Quake.PublicID, message, tags)

	if err != nil {
		log.Printf(err.Error())
		m.SetErr(err)
		return true
	}

	idp.Add(*m.Quake)

	return false
}