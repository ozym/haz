package msg

import (
	"log"
	"time"
)

type HeartBeat struct {
	ServiceID string
	SentTime  time.Time
	err       error
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
