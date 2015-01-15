package impact

import (
	"testing"
	"time"
)

func TestImpact(t *testing.T) {
	i := Intensity{}

	if i.Valid() {
		t.Error("zero intensity should not be valid.")
	}
	if i.Err() == nil {
		t.Error("not valid should have set i.err.")
	}

	i = Intensity{}

	if !i.Old() {
		t.Error("zero intensity should be old.")
	}
	if i.Err() == nil {
		t.Error("old should have set i.err.")
	}

	i = Intensity{}

	i.Source = "test.test"
	i.Quality = "measured"
	i.MMI = 4

	if !i.Valid() {
		t.Error("should be valid")
	}
	if i.Err() != nil {
		t.Error("valid should have nil i.err.")
	}

	i.Time = time.Now().UTC().Add(time.Duration(-30) * time.Minute)

	if i.Old() {
		t.Error("should not be old")
	}
	if i.Err() != nil {
		t.Error("not old should have nil i.err.")
	}

	i.Time = time.Now().UTC().Add(time.Duration(-61) * time.Minute)
	if !i.Old() {
		t.Error("should be old.")
	}
	if i.Err() == nil {
		t.Error("old should have set i.err.")
	}
}
