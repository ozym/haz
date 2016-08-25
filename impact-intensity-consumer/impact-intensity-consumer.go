package main

// impact-intensity-consumer to take intensity messages from SQS and store them in the impact DB.

import (
	"fmt"
	"github.com/GeoNet/haz/database"
	"github.com/GeoNet/haz/msg"
	"github.com/GeoNet/haz/sqs"
	_ "github.com/lib/pq"
	"log"
	"time"
	"sync"
	"math"
)

var (
	db database.DB
	retry = time.Duration(30) * time.Second
	expireInterval = time.Duration(10) * time.Second
)

type message struct {
	msg.Intensity
}

// reportedLimiter is used to limit the rate at which each source
// can send reports.
var reportedLimiter = struct {
	sync.RWMutex
	seen map[string]time.Time
}{seen: make(map[string]time.Time)}

func init() {
	sqs.MaxNumberOfMessages = 1
	sqs.VisibilityTimeout = 600
	sqs.WaitTimeSeconds = 20

	// once per minute remove any seen sources from more than 1 hour ago.
	go func() {
		ticker := time.NewTicker(time.Minute).C

		for {
			select {
			case <-ticker:
				reportedLimiter.Lock()

				now := time.Now().UTC().Add(time.Hour * -1)

				for k, t := range reportedLimiter.seen {
					if t.Before(now) {
						delete(reportedLimiter.seen, k)
					}
				}

				reportedLimiter.Unlock()
			}
		}
	}()
}

// main sets up the DB connection and then runs listen to process intensity messages from SQS.
func main() {
	var err error

	db, err = database.InitPG()
	if err != nil {
		log.Fatalf("ERROR: problem with DB config: %s", err)
	}
	defer db.Close()

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
	rx, dx, err := sqs.InitRx()
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
		m.seenReported()
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
	if m.Err() != nil {
		return
	}

	_, err := db.Exec("select impact.add_intensity_reported($1, $2, $3, $4, $5, $6)", m.Source, m.Longitude, m.Latitude, m.Time, m.MMI, m.Comment)

	if err != nil {
		m.SetErr(err)
		return
	}

	reportedLimiter.Lock()

	t := reportedLimiter.seen[m.Source]

	// there are no guarantees that messages arrive in time order.  In
	// general they are oldest first.
	if m.Time.After(t) {
		reportedLimiter.seen[m.Source] = m.Time
	}

	reportedLimiter.Unlock()
}

// seenReported sets i.err if a source sends reports that are within 60s.
func (m *message) seenReported() {
	if m.Err() != nil {
		return
	}

	reportedLimiter.RLock()
	t, ok := reportedLimiter.seen[m.Source]
	reportedLimiter.RUnlock()

	// if there is no entry for this source in reportedLimiter there is no more work to do.
	if !ok {
		return
	}

	if math.Abs(t.Sub(m.Time).Seconds()) < 60.0 {
		m.SetErr(fmt.Errorf("message from source %s already seen within 60s", m.Source))
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
