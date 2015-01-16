// Package msg provides interfaces and methods for processing messages.
//
// Models a message flow that looks like:
//   [transport] receive -> decode -> process -> encode -> send [transport]
// There is no requirement for a flow to have both send and receive.
//
// Message processing is tracked into expvar.  The following values are added:
//  * messages.process - the number of messages sent to process.  Updated every 20s.
//  * messages.error - the number of messages that resulted in an error in 20s.  Updates every 20s.
//  * messages.time - the average processing time per message in 20s.  Updates every 20s.
//
package msg

import (
	"expvar"
	"github.com/GeoNet/app/metrics"
	"log"
	"time"
)

var (
	pt metrics.Timer
	p  metrics.Rate
	e  metrics.Rate
	v  = expvar.NewMap("messages")
)

func init() {
	pt.Init(time.Duration(20) * time.Second)
	p.Init(time.Duration(20)*time.Second, time.Duration(20)*time.Second)
	e.Init(time.Duration(20)*time.Second, time.Duration(20)*time.Second)
	v.Init()
	v.Set("time", &pt)
	v.Set("process", &p)
	v.Set("error", &e)
}

// Message defines an interface the allows for message processing.
// Types that implement Message should use a no-op approach in Process()
// based on checking Err() (or it's private val).  See for example
// github.com/GeoNet/msg/impact.Intensity
type Message interface {
	Decode(b []byte)
	Encode() ([]byte, error)
	Process() (reprocess bool) // a hint to try to reprocess the message.
	Err() error
}

// Process excutes m.Process with logging and metrics.
// Logs errors.  Increments metrics.
func Process(m Message) bool {
	p.Inc()
	start := time.Now()
	s := m.Process()
	pt.Inc(start)
	if m.Err() != nil {
		log.Printf("WARN %s", m.Err())
		e.Inc()
	}
	return s
}
