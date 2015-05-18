// Package msg provides interfaces and methods for processing messages.
//
// Models a message flow that looks like:
//   [transport] receive -> decode -> process -> encode -> send [transport]
package msg

import (
	"github.com/GeoNet/metrics"
	"github.com/GeoNet/metrics/librato"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type Raw struct {
	Subject       string
	Body          string
	ReceiptHandle string
}

// Message defines an interface that allows for message processing.
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
	interval                        time.Duration // Rates calculated over interval.
	period                          time.Duration // Metrics updated every period.
	libratoUser, libratoKey, source string
	r, p, e                         metrics.Rate
	pt                              metrics.Timer
}

var (
	mtr metric
)

// InitLibrato initialises gathering and sending metrics to Librato metrics.
// Source can be used to further identify the source of the metrics.  If it is
// non empty it is appended to the hostname to form the metric source.
// Call from an init func.  Use empty strings to send metrics to the logs only.
func InitLibrato(user, key, source string) {
	mtr = metric{
		interval:    time.Duration(60) * time.Second,
		period:      time.Duration(60) * time.Second,
		libratoUser: user,
		libratoKey:  key,
		source:      source,
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
		go logRuntimeStats()
	}
}

// logMetrics and libratoMetrics could be combined with the use of a little more logic.  Keep them
// separated so it's easier to remove Librato or add other collectors.

func (m *metric) logMetrics() {
	rate := m.interval.String()
	for {
		select {
		case v := <-m.pt.Avg:
			log.Printf("Metric: Messages.ProcessingTime=%fs", v)
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

	host, err := os.Hostname()
	if err != nil {
		host = "unknown"
	}

	if m.source != "" {
		host = host + "-" + m.source
	}

	a := strings.Split(os.Args[0], "/")
	exe := a[len(a)-1]

	ptg := &librato.Gauge{Source: host, Name: exe + ".Messages.ProcessingTime"}
	rg := &librato.Gauge{Source: host, Name: exe + ".Messages.Received"}
	pg := &librato.Gauge{Source: host, Name: exe + ".Messages.Processed"}
	eg := &librato.Gauge{Source: host, Name: exe + ".Messages.Error"}

	// run time stats.
	mem := &runtime.MemStats{}
	memSys := &librato.Gauge{Source: host, Name: exe + ".MemStats.Sys"}
	memHeap := &librato.Gauge{Source: host, Name: exe + ".MemStats.HeapAlloc"}
	memStack := &librato.Gauge{Source: host, Name: exe + ".MemStats.StackInuse"}
	goRoutines := &librato.Gauge{Source: host, Name: exe + ".NumGoroutine"}

	rate := m.interval.String()

	var g []librato.Gauge

	for {
		select {
		case v := <-m.pt.Avg:
			ptg.SetValue(v)
			g = append(g, *ptg)
		case v := <-m.p.Avg:
			pg.SetValue(v)
			g = append(g, *pg)
		case v := <-m.e.Avg:
			eg.SetValue(v)
			g = append(g, *eg)
		case v := <-m.r.Avg:
			rg.SetValue(v)
			g = append(g, *rg)
			log.Printf("Messages received=%f per %s", v, rate)
		}
		if len(g) == 4 {
			runtime.ReadMemStats(mem)
			memSys.SetValue(float64(mem.Sys))
			memHeap.SetValue(float64(mem.HeapAlloc))
			memStack.SetValue(float64(mem.StackInuse))
			goRoutines.SetValue(float64(runtime.NumGoroutine()))

			g = append(g, *memSys)
			g = append(g, *memHeap)
			g = append(g, *memStack)
			g = append(g, *goRoutines)

			if len(lbr) < cap(lbr) { // the lbr chan shouldn't be blocked but would rather drop metrics and keep operating.
				lbr <- g
			}
			g = nil
		}
	}
}

func logRuntimeStats() {
	m := &runtime.MemStats{}
	s := time.Duration(60) * time.Second

	for {
		runtime.ReadMemStats(m)
		log.Printf("MemStats.Sys %d", m.Sys)
		log.Printf("MemStats.HeapAlloc %d", m.HeapAlloc)
		log.Printf("MemStats.StackInuse %d", m.StackInuse)
		log.Printf("NumGoroutine %d", runtime.NumGoroutine())
		time.Sleep(s)
	}
}
