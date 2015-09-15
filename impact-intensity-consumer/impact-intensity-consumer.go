package main

// impact-intensity-consumer to take intensity messages from SQS and store them in the impact DB.

import (
	"database/sql"
	"fmt"
	"github.com/GeoNet/cfg"
	"github.com/GeoNet/haz/msg"
	"github.com/GeoNet/haz/sqs"
	"github.com/GeoNet/log/logentries"
	_ "github.com/lib/pq"
	"log"
	"time"
)

//go:generate configer impact-intensity-consumer.json
var (
	config         = cfg.Load()
	db             *sql.DB
	retry          = time.Duration(30) * time.Second
	expireInterval = time.Duration(10) * time.Second
)

type message struct {
	msg.Intensity
}

func init() {
	logentries.Init(config.Logentries.Token)
	msg.InitLibrato(config.Librato.User, config.Librato.Key, config.Librato.Source)
	config.SQS.MaxNumberOfMessages = 1
	config.SQS.VisibilityTimeout = 600
	config.SQS.WaitTimeSeconds = 20
}

// main sets up the DB connection and then runs listen to process intensity messages from SQS.
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
	go deleteExpired()

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

// listen listens to the SQS queue for intensity messages and saves them to the DB.
func listen() {
	rx, dx, err := sqs.InitRx(config.SQS)
	if err != nil {
		log.Fatalf("ERROR - problem creating SQS from config: %s", err)
	}

	for {
		r := <-rx
		m := message{}
		m.Decode([]byte(r.Body))
		if !msg.Process(&m) {
			dx <- r.ReceiptHandle
		}
	}
}

func (m *message) Process() bool {
	// To be stored to the DB messages must be valid and in the last 60 minutes.
	m.Valid()
	m.Old()
	m.Future()

	if m.Err() != nil {
		return false
	}

	switch m.Quality {
	case "measured":
		m.saveMeasured()
	case "reported":
		m.saveReported()
	default:
		m.SetErr(fmt.Errorf("no method to save intensity message with quality: %s", m.Quality))
		return false
	}

	if m.Err() != nil {
		dbPing()
		return true
	}

	return false
}

func (m *message) saveMeasured() {
	_, err := db.Exec("select impact.add_intensity_measured($1, $2, $3, $4, $5)", m.Source, m.Longitude, m.Latitude, m.Time, m.MMI)

	if err != nil {
		m.SetErr(err)
	}
}

func (m *message) saveReported() {
	_, err := db.Exec("select impact.add_intensity_reported($1, $2, $3, $4, $5, $6)", m.Source, m.Longitude, m.Latitude, m.Time, m.MMI, m.Comment)

	if err != nil {
		m.SetErr(err)
	}
}

func deleteExpired() {
	for {
		_, err := db.Exec("delete from impact.intensity_measured where time < $1", time.Now().UTC().Add(msg.MeasuredAge))
		if err != nil {
			log.Printf("WARN expiring measured instensity values: %s", err)
		}
		time.Sleep(expireInterval)
	}
}
