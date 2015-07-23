package main

import (
	"github.com/GeoNet/cfg"
	"github.com/GeoNet/haz/msg"
	"github.com/GeoNet/haz/sqs"
	"github.com/GeoNet/log/logentries"
	"log"
	"net/smtp"
	"strconv"
)

//go:generate configer haz-eqnews-consumer.json
var (
	config   = cfg.Load()
	idp      = msg.IdpQuake{}
	mailHost string
	auth     smtp.Auth
)

func init() {
	logentries.Init(config.Logentries.Token)
	msg.InitLibrato(config.Librato.User, config.Librato.Key, config.Librato.Source)
	config.SQS.MaxNumberOfMessages = 1
	config.SQS.VisibilityTimeout = 600
	config.SQS.WaitTimeSeconds = 20

	mailHost = config.SMTP.Host + ":" + strconv.Itoa(config.SMTP.Port)

	auth = smtp.PlainAuth("",
		config.SMTP.UserName,
		config.SMTP.Password,
		config.SMTP.Host)
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

		if idp.Seen(*m.Quake) {
			log.Printf("Already sent email for %s", m.Quake.PublicID)
			return false
		}

		alert, subject, body := m.Quake.AlertEqNews()

		if alert {
			log.Printf("Sending email for quake %s", m.Quake.PublicID)

			mail := "Subject: " + subject + "\r\n\r\n" + body

			err := smtp.SendMail(mailHost, auth,
				config.SMTP.From, []string{config.SMTP.To},
				[]byte(mail))

			if err != nil {
				m.SetErr(err)
				return true
			}

			idp.Add(*m.Quake)
		}
	}

	return false
}
