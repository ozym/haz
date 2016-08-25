package main

import (
	"testing"
	"time"
)

func TestSeenReported(t *testing.T) {
	m := message{}
	m.Source = "test"
	m.Time = time.Now().UTC()

	m.seenReported()

	if m.Err() != nil {
		t.Errorf("should get nil error %v", m.Err())
	}

	// add this source to the rate limiter, usually happens in saveReported
	reportedLimiter.Lock()
	reportedLimiter.seen[m.Source] = m.Time
	reportedLimiter.Unlock()

	m.seenReported()

	if m.Err() == nil {
		t.Error("should get non nil error for rate limiter")
	}

	// 30 seconds difference should still trigger the limiter
	m.Time = m.Time.Add(time.Second * -30)
	m.SetErr(nil)

	m.seenReported()

	if m.Err() == nil {
		t.Error("should get non nil error for rate limiter")
	}

	// more than 60s difference from the first message should not trigger the limiter
	m.Time = m.Time.Add(time.Second * -90)
	m.SetErr(nil)

	m.seenReported()

	if m.Err() != nil {
		t.Errorf("should get nil error %v", m.Err())
	}
}
