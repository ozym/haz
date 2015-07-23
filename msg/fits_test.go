package msg

import (
	"io/ioutil"
	"testing"
)

func TestFits1(t *testing.T) {
	b, err := ioutil.ReadFile("etc/fits-1.json")
	if err != nil {
		t.Fatal(err)
	}

	o := Observation{}

	o.Decode(b)
	if o.Err() != nil {
		t.Fatal(o.Err())
	}

	if o.SiteID != "VGT2" {
		t.Error("wrong site for obs")
	}
}
