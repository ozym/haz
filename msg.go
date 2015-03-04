// Package msg provides interfaces and methods for processing messages.
//
// Models a message flow that looks like:
//   [transport] receive -> decode -> process -> encode -> send [transport]
package msg

import (
	"github.com/GeoNet/app/metrics"
	"github.com/GeoNet/app/metrics/librato"
	"log"
	"os"
	"strings"
	"time"
)

// Message defines an interface the allows for message processing.
type Message interface {
	Process() (reprocess bool) // a hint to try to reprocess the message.
	Err() error
}

// Process excutes m.Process with logging and metrics.
func Process(m Message) bool {
	mtr.r.Inc()
	start := time.Now()
	s := m.Process()
	mtr.pt.Inc(start)
	mtr.p.Inc()
	if m.Err() != nil {
		log.Printf("%s", m.Err().Error())
		mtr.e.Inc()
	}
	return s
}

// metrics gathering

type metric struct {
	interval                time.Duration // Rates calculated over interval.
	period                  time.Duration // Metrics updated every period.
	libratoUser, libratoKey string
	r, p, e                 metrics.Rate
	pt                      metrics.Timer
}

var (
	mtr metric
)

// InitLibrato initialises gathering and sending metrics to Librato metrics.
// Call from an init func.  Use empty strings to send metrics to the logs only.
func InitLibrato(user, key string) {
	mtr = metric{
		interval:    time.Duration(60) * time.Second,
		period:      time.Duration(60) * time.Second,
		libratoUser: user,
		libratoKey:  key,
	}

	mtr.pt.Init(mtr.period)
	mtr.r.Init(mtr.period, mtr.interval)
	mtr.p.Init(mtr.period, mtr.interval)
	mtr.e.Init(mtr.period, mtr.interval)

	if mtr.libratoUser != "" && mtr.libratoKey != "" {
		log.Println("Sending metrics to Librato Metrics.")
		go mtr.libratoMetrics()
	} else {
		log.Println("Sending metrics to logger only.")
		go mtr.logMetrics()
	}
}

// logMetrics and libratoMetrics could be combined with the use of a little more logic.  Keep them
// separated so it's easier to remove Librato or add other collectors.

func (m *metric) logMetrics() {
	rate := m.interval.String()
	for {
		select {
		case v := <-m.pt.Avg:
			log.Printf("Metric: Messages.AverageProcessingTime=%fs", v)
		case v := <-m.p.Avg:
			log.Printf("Metric: Messages.Processed=%f per %s", v, rate)
		case v := <-m.e.Avg:
			log.Printf("Metric: Messages.Errors=%f per %s", v, rate)
		case v := <-m.r.Avg:
			log.Printf("Metric: Messages.Received=%f per %s", v, rate)
		}
	}
}

func (m *metric) libratoMetrics() {
	lbr := make(chan []librato.Gauge, 1)

	librato.Init(m.libratoUser, m.libratoKey, lbr)

	rate := m.interval.String()

	host, err := os.Hostname()
	if err != nil {
		host = "unknown"
	}

	a := strings.Split(os.Args[0], "/")
	source := a[len(a)-1]

	ptg := &librato.Gauge{Source: host, Name: source + ".Messages.ProcessingTime"}
	rg := &librato.Gauge{Source: source, Name: "Messages.Received"}
	pg := &librato.Gauge{Source: source, Name: "Messages.Processed"}
	eg := &librato.Gauge{Source: source, Name: "Messages.Error"}

	var g []librato.Gauge

	for {
		select {
		case v := <-m.pt.Avg:
			ptg.SetValue(v)
			g = append(g, *ptg)
			log.Printf("Metric: Messages.AverageProcessingTime=%fs", ptg.Value)
		case v := <-m.p.Avg:
			pg.SetValue(v)
			g = append(g, *pg)
			log.Printf("Metric: Messages.Processed=%f per %s", pg.Value, rate)
		case v := <-m.e.Avg:
			eg.SetValue(v)
			g = append(g, *eg)
			log.Printf("Metric: Messages.Errors=%f per %s", eg.Value, rate)
		case v := <-m.r.Avg:
			rg.SetValue(v)
			g = append(g, *rg)
			log.Printf("Metric: Messages.Received=%f per %s", rg.Value, rate)
		}
		if len(g) == 4 {
			if len(lbr) < cap(lbr) { // the lbr chan shouldn't be blocked but would rather drop metrics and keep operating.
				lbr <- g
			}
			g = nil
		}
	}
}
