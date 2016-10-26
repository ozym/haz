package sc3ml

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestEarthquake(t *testing.T) {
	var err error
	var f *os.File
	var b []byte

	if f, err = os.Open("etc/2015p768477.xml"); err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if b, err = ioutil.ReadAll(f); err != nil {
		t.Fatal(err)
	}

	e, err := QuakeTechnical(b)
	if err != nil {
		t.Fatal(err)
	}

	if e.PublicID != "2015p768477" {
		t.Errorf("expected publicID 2015p768477 got %s", e.PublicID)
	}

	if e.Type != "earthquake" {
		t.Errorf("expected type earthquake got %s", e.Type)
	}

	if e.Agency != "WEL(GNS_Primary)" {
		t.Errorf("AgencyID expected WEL(GNS_Primary) got %s", e.Agency)
	}

	if time.Unix(e.ModificationTime.Sec, e.ModificationTime.Nsec).UTC().Format(time.RFC3339Nano) != "2015-10-12T22:46:41.228824Z" {
		t.Errorf("ModificationTime expected 2015-10-12T22:46:41.228824Z got %s", time.Unix(e.ModificationTime.Sec, e.ModificationTime.Nsec).Format(time.RFC3339Nano))
	}

	if time.Unix(e.Time.Sec, e.Time.Nsec).UTC().Format(time.RFC3339Nano) != "2015-10-12T08:05:01.717692Z" {
		t.Errorf("expected 2015-10-12T08:05:01.717692Z, got %s", time.Unix(e.Time.Sec, e.Time.Nsec).UTC().Format(time.RFC3339Nano))
	}

	if e.Latitude.Value != -40.57806609 {
		t.Errorf("Latitude expected -40.57806609 got %f", e.Latitude)
	}
	if e.Latitude.Uncertainty != 1.922480006 {
		t.Errorf("Latitude uncertainty expected 1.922480006 got %f", e.Latitude.Uncertainty)
	}

	if e.Longitude.Value != 176.3257242 {
		t.Errorf("Longitude expected 176.3257242 got %f", e.Longitude.Value)
	}
	if e.Longitude.Uncertainty != 3.435738791 {
		t.Errorf("Longitude uncertainty expected 3.435738791 got %f", e.Longitude.Uncertainty)
	}

	if e.Depth.Value != 23.28125 {
		t.Errorf("Depth expected 23.28125 got %f", e.Depth.Value)
	}
	if e.Depth.Uncertainty != 3.575079654 {
		t.Errorf("Depth uncertainty expected 3.575079654 got %f", e.Depth.Uncertainty)
	}

	if e.Method != "NonLinLoc" {
		t.Errorf("MethodID expected NonLinLoc got %s", e.Method)
	}

	if e.EarthModel != "nz3drx" {
		t.Errorf("EarthModelID expected NonLinLoc got %s", e.EarthModel)
	}

	if e.StandardError != 0.5592857863 {
		t.Errorf("StandardError expected 0.5592857863 got %f", e.StandardError)
	}

	if e.AzimuthalGap != 166.4674465 {
		t.Errorf("AzimuthalGap expected 166.4674465 got %f", e.AzimuthalGap)
	}

	if e.MinimumDistance != 0.1217162272 {
		t.Errorf("MinimumDistance expected 0.1217162272 got %f", e.MinimumDistance)
	}

	if e.UsedPhaseCount != 44 {
		t.Errorf("UsedPhaseCount expected 44 got %d", e.UsedPhaseCount)
	}

	if e.UsedStationCount != 32 {
		t.Errorf("UsedStationCount expected 32 got %d", e.UsedStationCount)
	}

	if e.MagnitudeType != "M" {
		t.Errorf("e.PreferredMagnitude.Type expected M got %s", e.MagnitudeType)
	}

	if e.Magnitude.Value != 5.691131913 {
		t.Errorf("magnitude expected 5.691131913 got %f", e.Magnitude.Value)
	}

	if e.Magnitude.Uncertainty != 0 {
		t.Errorf("uncertainty expected 0 got %f", e.Magnitude.Uncertainty)
	}

	if len(e.Pick) != 190 {
		t.Errorf("expected 190 Pick got %d", len(e.Pick))
	}

	var found bool

	for _, v := range e.Pick {
		if v.Waveform.Network == "NZ" && v.Waveform.Station == "BFZ" && v.Waveform.Location == "10" && v.Waveform.Channel == "HHN" {
			found = true
			if v.Time.Sec != 1444637106 {
				t.Errorf("pick sec expected 1444637106 got %d", v.Time.Sec)
			}

			if v.Time.Nsec != 792207000 {
				t.Errorf("pick nsec expected 792207000 got %d", v.Time.Nsec)
			}

			if v.Phase != "P" {
				t.Errorf("expected P phase got: %s", v.Phase)
			}

			if v.Azimuth != 211.917806 {
				t.Errorf("azimuth expected 211.917806 got %f", v.Azimuth)
			}

			if v.Distance != 0.1217162272 {
				t.Errorf("distance exected 0.1217162272 got %f", v.Distance)
			}

			if v.Residual != -0.01664948232 {
				t.Errorf("expected residual -0.01664948232 got %f", v.Residual)
			}

			if v.Weight != 1.406866218 {
				t.Errorf("expected weight 1.406866218 got %f", v.Weight)
			}

			if v.EvaluationMode != "manual" {
				t.Errorf("expected evaluation mode manual got %s", v.EvaluationMode)
			}
		}
	}

	if !found {
		t.Error("didn't find pick waveform NZ BFZ 10 HHZ.")
	}

	if len(e.Magnitudes) != 3 {
		t.Errorf("expected 3 magnitudes got %d", len(e.Magnitudes))
	}

	var foundMag bool

	for _, m := range e.Magnitudes {
		if m.Type == "MLv" {
			foundMag = true

			if m.Magnitude.Value != 5.691131913 {
				t.Errorf("expected magnitude 5.691131913 got %f", m.Magnitude.Value)
			}

			if len(m.StationMagnitude) != 171 {
				t.Errorf("expected 171 station magnitudes got %d", len(m.StationMagnitude))
			}

			var found bool

			for _, v := range m.StationMagnitude {
				if v.Waveform.Network == "NZ" && v.Waveform.Station == "BFZ" && v.Waveform.Location == "10" && v.Waveform.Channel == "HHZ" {
					found = true

					if v.Type != "MLv" {
						t.Errorf("expected type MLv got %s", v.Type)
					}

					if v.Magnitude.Value != 5.210250362 {
						t.Errorf("expected magnitude 5.210250362 got %f", v.Magnitude.Value)
					}

					if v.Azimuth != 211.917806 {
						t.Errorf("expected azimuth 211.917806 got %f", v.Azimuth)
					}

					if v.Distance != 0.1217162272 {
						t.Errorf("expected distance 0.1217162272 got %f", v.Distance)
					}

					if v.Weight != 1 {
						t.Errorf("expected weight 1 got %f", v.Weight)
					}

					if v.Residual != 0.000000 {
						t.Errorf("expected residual 0.000000 got %f", v.Residual)
					}

					if v.Amplitude.Value != 3819.596931 {
						t.Errorf("expected amplitude 3819.596931 got %f", v.Amplitude.Value)
					}
				}
			}

			if !found {
				t.Error("didn't find station magnitude waveform NZ BFZ 10 HHZ.")
			}
		}
	}

	if !foundMag {
		t.Error("didn't find mag MLv.")
	}
}
