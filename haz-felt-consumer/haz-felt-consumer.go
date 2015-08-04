// haz-felt-client listens to an AWS SQS queue for Haz JSON messages and
// saves the messages into a DB.
package main

import (
	"database/sql"
	"github.com/GeoNet/cfg"
	"github.com/GeoNet/haz/msg"
	"github.com/GeoNet/haz/sqs"
	"github.com/GeoNet/log/logentries"
	_ "github.com/lib/pq"
	"log"
	"time"
)

//go:generate configer haz-felt-consumer.json
var (
	config = cfg.Load()
	db     *sql.DB
	retry  = time.Duration(30) * time.Second
)

func init() {
	logentries.Init(config.Logentries.Token)
	msg.InitLibrato(config.Librato.User, config.Librato.Key, config.Librato.Source)
	config.SQS.MaxNumberOfMessages = 1
	config.SQS.VisibilityTimeout = 600
	config.SQS.WaitTimeSeconds = 20
}

type message struct {
	msg.Haz
}

func main() {
	var err error

	db, err = sql.Open("postgres", config.DataBase.Postgres())
	if err != nil {
		log.Fatalf("ERROR: problem with DB config: %s", err)
	}
	defer db.Close()

	db.SetMaxIdleConns(config.DataBase.MaxIdleConns)
	db.SetMaxOpenConns(config.DataBase.MaxOpenConns)

	dbPing()

	log.Println("starting message listener.")
	listen()
}

// dbPing does not return until it has successfully pinged the DB.
func dbPing() {
	for {
		if err := db.Ping(); err != nil {
			log.Printf("WARN - pinging DB: %s", err)
			log.Println("WARN - waiting then trying DB again.")
			time.Sleep(retry)
			continue
		}
		break
	}
}

// listen for haz messages and saves them to the DB.
func listen() {
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
	// case statements in switch have no fallthrough to next
	// statement e.g., we're assuming that Haz messages only hold one type.
	switch {
	case m.Err() != nil:
		log.Println("WARN received errored message: " + m.Err().Error())
		return false
	case m.HeartBeat != nil:
		m.saveHeartBeat()
	case m.Quake != nil:
		m.saveQuake()
	}

	// Block processing here if we can't contact the DB (the most likely source of
	// errors at this point). This leaves all the messages except the currrent one visible on the queue.
	// Then ask for the message to be redelivered
	if m.Err() != nil {
		dbPing()
		return true
	}

	return false
}

// saveQuake saves quakes to the DB.
func (m *message) saveQuake() {
	m.Quake.RxLog()
	q := m.Quake
	_, err := db.Exec("SELECT qrt.add_event($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
		q.PublicID, q.AgencyID, q.Latitude, q.Longitude, q.Time, q.ModificationTime, q.Depth,
		q.UsedPhaseCount, q.Magnitude, q.MagnitudeType, q.Status(), q.Type)

	if err != nil {
		m.SetErr(err)
		return
	}
}

// saveHeartBeat saves heart beat messages to the DB.
func (m *message) saveHeartBeat() {
	if m.HeartBeat.RxLog() {
		return
	}
	_, err := db.Exec(`SELECT qrt.add_soh($1, $2)`, m.HeartBeat.ServiceID, m.HeartBeat.SentTime)
	if err != nil {
		m.SetErr(err)
		return
	}
}