// Package msg provides interfaces and methods for processing messages.
//
// Models a message flow that looks like:
//   [transport] receive -> decode -> process -> encode -> send [transport]
// There is no requirement for a flow to have both send and receive.
//
// Message processing is tracked into expvar.  The following values are added:
//  * messages.inbound - the number of messages inbound for processing in 20s.  Updated every 20s.
//  * messages.outbound - the number of messages outbound from processing in 20s.  Updated every 20s.
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
	pt  metrics.Timer
	in  metrics.Rate
	out metrics.Rate
	e   metrics.Rate
	v   = expvar.NewMap("messages")
)

func init() {
	pt.Init(time.Duration(20) * time.Second)
	in.Init(time.Duration(20)*time.Second, time.Duration(20)*time.Second)
	out.Init(time.Duration(20)*time.Second, time.Duration(20)*time.Second)
	e.Init(time.Duration(20)*time.Second, time.Duration(20)*time.Second)
	v.Init()
	v.Set("time", &pt)
	v.Set("inbound", &in)   // use 'inbound' instead of received incase we ever want to track that as well.
	v.Set("outbound", &out) // use 'outbound' instead of sent incase we ever want to track that as well.
	v.Set("error", &e)
}

type Message interface {
	Unmarshal(b []byte) error
	Marshal() ([]byte, error)
	// Process returns a bool that can be used as a hint for the
	// transport.  If false it may be appropriate for the transport
	// to try to redeliver the message.
	Process() bool
	Err() error
}

// Process excutes m.Process with logging and metrics.
// Logs errors.  Increments metrics.
func Process(m Message) bool {
	start := time.Now()
	s := m.Process()
	pt.Inc(start)
	if m.Err() != nil {
		log.Printf("WARN %s", m.Err())
		e.Inc()
	}
	return s
}

// Decode uses m.Unmarshal to decode an inbound message.  Returns false
// if the message cannot be unmarshaled in which case
// it may be appropriate for the transport to discard the message.
// Logs errors.  Increments metrics.
func Decode(m Message, b []byte) bool {
	in.Inc()
	if err := m.Unmarshal(b); err != nil {
		log.Printf("WARN: %s", err)
		e.Inc()
		return false
	}
	return true
}

// Encode uses m.Marshal to encode m into the byte slice pointed to by b.
// Returns false if the message cannot be marshaled in which case it is
// probably inappropriate to send the message to the outbound transport.
// Logs errors.  Increments metrics.
func Encode(m Message, b *[]byte) bool {
	out.Inc()
	var err error
	*b, err = m.Marshal()
	if err != nil {
		log.Printf("WARN: %s", err)
		e.Inc()
		return false
	}
	return true
}
