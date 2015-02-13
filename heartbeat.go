package msg

import (
	"fmt"
	"log"
	"time"
)

var heartBeatAge = time.Duration(-5) * time.Minute

type HeartBeat struct {
	ServiceID string
	SentTime  time.Time
	err       error
}

// Old sets Error if the HeartBeat pointed to by h is old.
func (h *HeartBeat) Old() {
	if h.err != nil {
		return
	}
	if h.SentTime.Before(time.Now().UTC().Add(heartBeatAge)) {
		h.err = fmt.Errorf("old HeartBeat message from %s", h.ServiceID)
	}
	return
}

func (h *HeartBeat) RxLog() {
	if h.err != nil {
		return
	}

	log.Printf("Received heartbeat for %s", h.ServiceID)
}

func (h *HeartBeat) TxLog() {
	if h.err != nil {
		return
	}

	log.Printf("Sending heartbeat for %s", h.ServiceID)
}

func (h *HeartBeat) Err() error {
	return h.err
}

func (h *HeartBeat) SetErr(err error) {
	h.err = err
}
