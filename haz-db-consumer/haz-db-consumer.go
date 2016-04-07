// haz-db-consumer listens to an AWS SQS queue for Haz JSON messages and
// saves the messages into a DB.
package main

import (
	"github.com/GeoNet/haz/database"
	"github.com/GeoNet/haz/msg"
	"github.com/GeoNet/haz/sqs"
	"github.com/GeoNet/log/logentries"
	_ "github.com/lib/pq"
	"log"
	"os"
)

//go:generate configer haz-db-consumer.json
var (
	db database.DB
)

func init() {
	logentries.Init(os.Getenv("LOGENTRIES_TOKEN"))
	msg.InitLibrato(os.Getenv("LIBRATO_USER"), os.Getenv("LIBRATO_KEY"), os.Getenv("LIBRATO_SOURCE"))
	sqs.MaxNumberOfMessages = 1
	sqs.VisibilityTimeout = 600
	sqs.WaitTimeSeconds = 20
}

type message struct {
	msg.Haz
}

func main() {
	var err error

	db, err = database.InitPG()
	if err != nil {
		log.Fatalf("ERROR: problem with DB config: %s", err)
	}
	defer db.Close()

	db.Check()

	log.Println("starting message listener.")
	listen()
}

// listen for haz messages and saves them to the DB.
func listen() {
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
	// case statements in switch have no fallthrough to next
	// statement e.g., we're assuming that Haz messages only hold one type.
	switch {
	case m.Err() != nil:
		log.Println("WARN received errored message: " + m.Err().Error())
		return false
	case m.HeartBeat != nil:
		m.HeartBeat.RxLog()
		m.SetErr(db.SaveHeartBeat(*m.HeartBeat))
	case m.Quake != nil:
		m.Quake.RxLog()
		m.SetErr(db.SaveQuake(*m.Quake))
	}

	// Block processing here if we can't contact the DB (the most likely source of
	// errors at this point). This leaves all the messages except the currrent one visible on the queue.
	// Then ask for the message to be redelivered
	if m.Err() != nil {
		db.Check()
		return true
	}

	return false
}
