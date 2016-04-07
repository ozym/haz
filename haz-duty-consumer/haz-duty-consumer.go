package main

import (
	"github.com/GeoNet/haz/msg"
	"github.com/GeoNet/haz/pagerduty"
	"github.com/GeoNet/haz/sqs"
	"github.com/GeoNet/log/logentries"
	"log"
	"os"
)

//go:generate configer haz-duty-consumer.json
var (
	idp = msg.IdpQuake{}
	pd  *pagerduty.Client
)

func init() {
	logentries.Init(os.Getenv("LOGENTRIES_TOKEN"))
	msg.InitLibrato(os.Getenv("LIBRATO_USER"), os.Getenv("LIBRATO_KEY"), os.Getenv("LIBRATO_SOURCE"))
	pd = pagerduty.Init()
	sqs.MaxNumberOfMessages = 1
	sqs.VisibilityTimeout = 600
	sqs.WaitTimeSeconds = 20
}

type message struct {
	msg.Haz
}

func main() {
	rx, dx, err := sqs.InitRx()
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

		if idp.Seen(*m.Quake) {
			log.Printf("Already sent notification for %s", m.Quake.PublicID)
			return false
		}

		alert, message := m.Quake.AlertDuty()
		if alert {
			log.Printf("Notifying the duty officer for quake %s", m.Quake.PublicID)
			err := pd.Trigger(message, m.Quake.PublicID, 3)
			if err != nil {
				m.SetErr(err)
				return true
			}

			idp.Add(*m.Quake)
		}
	}

	return false
}
