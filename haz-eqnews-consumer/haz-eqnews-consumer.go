package main

import (
	"github.com/GeoNet/haz/msg"
	"github.com/GeoNet/haz/sqs"
	"github.com/GeoNet/log/logentries"
	"log"
	"net/smtp"
	"os"
)

//go:generate configer haz-eqnews-consumer.json
var (
	idp      = msg.IdpQuake{}
	mailHost string
	auth     smtp.Auth
	smtpFrom = os.Getenv("SMTP_FROM")
	smtpTo   = os.Getenv("SMTP_TO")
)

func init() {
	logentries.Init(os.Getenv("LOGENTRIES_TOKEN"))
	msg.InitLibrato(os.Getenv("LIBRATO_USER"), os.Getenv("LIBRATO_KEY"), os.Getenv("LIBRATO_SOURCE"))
	sqs.MaxNumberOfMessages = 1
	sqs.VisibilityTimeout = 600
	sqs.WaitTimeSeconds = 20

	mailHost = os.Getenv("SMTP_HOST") + ":" + os.Getenv("SMTP_PORT")

	auth = smtp.PlainAuth("",
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASSWORD"),
		os.Getenv("SMTP_HOST"))
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
			log.Printf("Already sent email for %s", m.Quake.PublicID)
			return false
		}

		alert, subject, body := m.Quake.AlertEqNews()

		if alert {
			log.Printf("Sending email for quake %s", m.Quake.PublicID)

			mail := "Subject: " + subject + "\r\n\r\n" + body

			err := smtp.SendMail(mailHost, auth,
				smtpFrom, []string{smtpTo},
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
