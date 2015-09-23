package msg

import (
	"testing"
	"time"
)

func TestImpact(t *testing.T) {
	i := Intensity{}

	i.Valid()
	if i.Err() == nil {
		t.Error("not valid should have set i.err.")
	}

	i = Intensity{}

	i.Old()
	if i.Err() == nil {
		t.Error("old should have set i.err.")
	}

	i = Intensity{}

	i.Source = "test.test"
	i.Quality = "measured"
	i.MMI = 4

	i.Valid()
	if i.Err() != nil {
		t.Error("valid should have nil i.err.")
	}

	i.Time = time.Now().UTC().Add(time.Duration(-30) * time.Minute)

	i.Old()
	if i.Err() != nil {
		t.Error("not old should have nil i.err.")
	}

	i.Time = time.Now().UTC().Add(time.Duration(-61) * time.Minute)
	i.Old()
	if i.Err() == nil {
		t.Error("old should have set i.err.")
	}

	i = Intensity{}

	i.Source = "test.test"
	i.Quality = "measured"
	i.MMI = 4

	i.Time = time.Now().UTC().Add(time.Duration(-30) * time.Minute)

	i.Future()
	if i.Err() != nil {
		t.Error("not future should have nil i.err.")
	}

	i.Time = time.Now().UTC().Add(time.Duration(30) * time.Second)

	i.Future()
	if i.Err() == nil {
		t.Error("future should have non nil i.err.")
	}

	i = Intensity{}

	i.Source = "test.test"
	i.Quality = "measured"
	i.MMI = 4

	i.Time = time.Now().UTC().Add(time.Duration(9) * time.Second)

	i.Future()
	if i.Err() != nil {
		t.Error("less than 10s in the future should have nil i.err.")
	}
}
