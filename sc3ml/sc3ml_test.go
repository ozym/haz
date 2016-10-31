package sc3ml

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestUnmarshal(t *testing.T) {
	var err error
	var ep eventParameters
	var f *os.File
	var b []byte

	if f, err = os.Open("etc/2015p768477.xml"); err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if b, err = ioutil.ReadAll(f); err != nil {
		t.Fatal(err)
	}

	if ep, err = unmarshal(b); err != nil {
		t.Fatal(err)
	}

	if len(ep.Events) != 1 {
		t.Fatal("should have found 1 event.")
	}

	e := ep.Events[0]

	if e.PublicID != "2015p768477" {
		t.Errorf("expected publicID 2015p768477 got %s", e.PublicID)
	}

	if e.Type != "earthquake" {
		t.Errorf("expected type earthquake got %s", e.Type)
	}

	if e.CreationInfo.AgencyID != "WEL(GNS_Primary)" {
		t.Errorf("AgencyID expected WEL(GNS_Primary) got %s", e.CreationInfo.AgencyID)
	}

	if e.CreationInfo.CreationTime.Format(time.RFC3339Nano) != "2015-10-12T08:05:25.610839Z" {
		t.Errorf("CreationTime expected 2015-10-12T08:05:25.610839Z got %s", e.CreationInfo.CreationTime.Format(time.RFC3339Nano))
	}

	if e.CreationInfo.ModificationTime.Format(time.RFC3339Nano) != "2015-10-12T22:46:41.228824Z" {
		t.Errorf("ModificationTime expected 2015-10-12T22:46:41.228824Z got %s", e.CreationInfo.ModificationTime.Format(time.RFC3339Nano))
	}

	if e.PreferredOriginID != "NLL.20151012224503.620592.155845" {
		t.Errorf("expected preferredOriginID NLL.20151012224503.620592.155845 got %s", e.PreferredOriginID)
	}

	if e.PreferredOrigin.PublicID != "NLL.20151012224503.620592.155845" {
		t.Errorf("expected NLL.20151012224503.620592.155845 got %s", e.PreferredOrigin.PublicID)
	}

	if e.PreferredOrigin.Time.Value.Format(time.RFC3339Nano) != "2015-10-12T08:05:01.717692Z" {
		t.Errorf("expected 2015-10-12T08:05:01.717692Z, got %s", e.PreferredOrigin.Time.Value.Format(time.RFC3339Nano))
	}

	if e.PreferredOrigin.Latitude.Value != -40.57806609 {
		t.Errorf("Latitude expected -40.57806609 got %f", e.PreferredOrigin.Latitude.Value)
	}
	if e.PreferredOrigin.Latitude.Uncertainty != 1.922480006 {
		t.Errorf("Latitude uncertainty expected 1.922480006 got %f", e.PreferredOrigin.Latitude.Uncertainty)
	}

	if e.PreferredOrigin.Longitude.Value != 176.3257242 {
		t.Errorf("Longitude expected 176.3257242 got %f", e.PreferredOrigin.Longitude.Value)
	}
	if e.PreferredOrigin.Longitude.Uncertainty != 3.435738791 {
		t.Errorf("Longitude uncertainty expected 3.435738791 got %f", e.PreferredOrigin.Longitude.Uncertainty)
	}

	if e.PreferredOrigin.Depth.Value != 23.28125 {
		t.Errorf("Depth expected 23.28125 got %f", e.PreferredOrigin.Depth.Value)
	}
	if e.PreferredOrigin.Depth.Uncertainty != 3.575079654 {
		t.Errorf("Depth uncertainty expected 3.575079654 got %f", e.PreferredOrigin.Depth.Uncertainty)
	}

	if e.PreferredOrigin.MethodID != "NonLinLoc" {
		t.Errorf("MethodID expected NonLinLoc got %s", e.PreferredOrigin.MethodID)
	}

	if e.PreferredOrigin.EarthModelID != "nz3drx" {
		t.Errorf("EarthModelID expected NonLinLoc got %s", e.PreferredOrigin.EarthModelID)
	}

	if e.PreferredOrigin.Quality.StandardError != 0.5592857863 {
		t.Errorf("StandardError expected 0.5592857863 got %f", e.PreferredOrigin.Quality.StandardError)
	}

	if e.PreferredOrigin.Quality.AzimuthalGap != 166.4674465 {
		t.Errorf("AzimuthalGap expected 166.4674465 got %f", e.PreferredOrigin.Quality.AzimuthalGap)
	}

	if e.PreferredOrigin.Quality.MinimumDistance != 0.1217162272 {
		t.Errorf("MinimumDistance expected 0.1217162272 got %f", e.PreferredOrigin.Quality.MinimumDistance)
	}

	if e.PreferredOrigin.Quality.UsedPhaseCount != 44 {
		t.Errorf("UsedPhaseCount expected 44 got %d", e.PreferredOrigin.Quality.UsedPhaseCount)
	}

	if e.PreferredOrigin.Quality.UsedStationCount != 32 {
		t.Errorf("UsedStationCount expected 32 got %d", e.PreferredOrigin.Quality.UsedStationCount)
	}

	var found bool
	for _, v := range e.PreferredOrigin.Arrivals {
		if v.PickID == "Pick#20151012081200.115203.26387" {
			found = true
			if v.Phase != "P" {
				t.Errorf("expected P got %s", v.Phase)
			}

			if v.Azimuth != 211.917806 {
				t.Errorf("azimuth expected 211.917806 got %f", v.Azimuth)
			}

			if v.Distance != 0.1217162272 {
				t.Errorf("distance expected 0.1217162272 got %f", v.Distance)
			}

			if v.Weight != 1.406866218 {
				t.Errorf("weight expected 1.406866218 got %f", v.Weight)
			}

			if v.TimeResidual != -0.01664948232 {
				t.Errorf("time residual expected -0.01664948232 got %f", v.TimeResidual)
			}

			if v.Pick.WaveformID.NetworkCode != "NZ" {
				t.Errorf("Pick.WaveformID.NetworkCode expected NZ, got %s", v.Pick.WaveformID.NetworkCode)
			}

			if v.Pick.WaveformID.StationCode != "BFZ" {
				t.Errorf("Pick.WaveformID.StationCode expected BFZ, got %s", v.Pick.WaveformID.StationCode)
			}

			if v.Pick.WaveformID.LocationCode != "10" {
				t.Errorf("Pick.WaveformID.LocationCode expected 10, got %s", v.Pick.WaveformID.LocationCode)
			}

			if v.Pick.WaveformID.ChannelCode != "HHN" {
				t.Errorf("Pick.WaveformID.ChannelCode expected HHN, got %s", v.Pick.WaveformID.ChannelCode)
			}

			if v.Pick.EvaluationMode != "manual" {
				t.Errorf("Pick.WaveformID.EvaluationMode expected manual got %s", v.Pick.EvaluationMode)
			}

			if v.Pick.EvaluationStatus != "" {
				t.Errorf("Pick.WaveformID.EvaluationStatus expected empty string got %s", v.Pick.EvaluationStatus)
			}

			if v.Pick.Time.Value.Format(time.RFC3339Nano) != "2015-10-12T08:05:06.792207Z" {
				t.Errorf("Pick.Time expected 2015-10-12T08:05:06.792207Z got %s", v.Pick.Time.Value.Format(time.RFC3339Nano))
			}
		}

	}
	if !found {
		t.Error("didn't find PickID Pick#20151012081200.115203.26387")
	}

	if e.PreferredMagnitude.Type != "M" {
		t.Errorf("e.PreferredMagnitude.Type expected M got %s", e.PreferredMagnitude.Type)
	}
	if e.PreferredMagnitude.Magnitude.Value != 5.691131913 {
		t.Errorf("magnitude expected 5.691131913 got %f", e.PreferredMagnitude.Magnitude.Value)
	}
	if e.PreferredMagnitude.Magnitude.Uncertainty != 0 {
		t.Errorf("uncertainty expected 0 got %f", e.PreferredMagnitude.Magnitude.Uncertainty)
	}
	if e.PreferredMagnitude.StationCount != 171 {
		t.Errorf("e.PreferredMagnitude.StationCount expected 171 got %d", e.PreferredMagnitude.StationCount)
	}
	if e.PreferredMagnitude.MethodID != "weighted average" {
		t.Errorf("MethodID expected weighted average got %s", e.PreferredMagnitude.MethodID)
	}

	if e.PreferredMagnitude.CreationInfo.AgencyID != "WEL(GNS_Primary)" {
		t.Errorf("AgencyID expected WEL(GNS_Primary) got %s", e.PreferredMagnitude.CreationInfo.AgencyID)
	}

	if e.PreferredMagnitude.CreationInfo.CreationTime.Format(time.RFC3339Nano) != "2015-10-12T22:46:41.218145Z" {
		t.Errorf("CreationTime expected 2015-10-12T22:46:41.218145Z got %s", e.PreferredMagnitude.CreationInfo.CreationTime.Format(time.RFC3339Nano))
	}

	found = false

	for _, m := range e.PreferredOrigin.Magnitudes {
		if m.PublicID == "Magnitude#20151012224509.743338.156745" {
			found = true

			if m.Type != "ML" {
				t.Error("m.Type expected ML, got ", m.Type)
			}
			if m.Magnitude.Value != 6.057227661 {
				t.Errorf("magnitude expected 6.057227661 got %f", m.Magnitude.Value)
			}
			if m.Magnitude.Uncertainty != 0.2576927171 {
				t.Errorf("Uncertainty expected 0.2576927171 got %f", m.Magnitude.Uncertainty)
			}
			if m.StationCount != 23 {
				t.Errorf("m.StationCount expected 23 got %d", m.StationCount)
			}
			if m.MethodID != "trimmed mean" {
				t.Errorf("m.MethodID expected trimmed mean got %s", m.MethodID)
			}

			if !(len(m.StationMagnitudeContributions) > 1) {
				t.Error("expected more than 1 StationMagnitudeContribution")
			}

			var foundSM bool

			for _, s := range m.StationMagnitudeContributions {
				if s.StationMagnitudeID == "StationMagnitude#20151012224509.743511.156746" {
					foundSM = true

					if s.Weight != 1.0 {
						t.Errorf("Weight expected 1.0 got %f", s.Weight)
					}

					if s.StationMagnitude.Magnitude.Value != 6.096018735 {
						t.Errorf("StationMagnitude.Magnitude.Value expected 6.096018735 got %f", s.StationMagnitude.Magnitude.Value)
					}

					if s.StationMagnitude.Type != "ML" {
						t.Errorf("StationMagnitude.Type expected ML got %s", s.StationMagnitude.Type)
					}

					if s.StationMagnitude.WaveformID.NetworkCode != "NZ" {
						t.Errorf("Pick.WaveformID.NetworkCode expected NZ, got %s", s.StationMagnitude.WaveformID.NetworkCode)
					}

					if s.StationMagnitude.WaveformID.StationCode != "ANWZ" {
						t.Errorf("Pick.WaveformID.StationCode expected ANWZ, got %s", s.StationMagnitude.WaveformID.StationCode)
					}

					if s.StationMagnitude.WaveformID.LocationCode != "10" {
						t.Errorf("Pick.WaveformID.LocationCode expected 10, got %s", s.StationMagnitude.WaveformID.LocationCode)
					}

					if s.StationMagnitude.WaveformID.ChannelCode != "EH" {
						t.Errorf("Pick.WaveformID.ChannelCode expected EH, got %s", s.StationMagnitude.WaveformID.ChannelCode)
					}

					if s.StationMagnitude.Amplitude.Amplitude.Value != 21899.94892 {
						t.Errorf("Amplitude.Value expected 21899.94892 got %f", s.StationMagnitude.Amplitude.Amplitude.Value)
					}
				}
			}
			if !foundSM {
				t.Error("did not find StationMagnitudeContrib StationMagnitude#20151012224509.743511.156746")
			}
		}
	}

	if !found {
		t.Error("did not find magnitude smi:scs/0.7/Origin#20131202033820.196288.25287#netMag.MLv")
	}

	if e.modificationTime().Format(time.RFC3339Nano) != "2015-10-12T22:46:41.228824Z" {
		t.Errorf("Modification time expected 2015-10-12T22:46:41.228824Z got %s", e.modificationTime().Format(time.RFC3339Nano))
	}
}

func TestDecodeSC3ML07CMT(t *testing.T) {
	var err error
	var ep eventParameters
	var f *os.File
	var b []byte

	if f, err = os.Open("etc/2016p408314-201606010431276083.xml"); err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if b, err = ioutil.ReadAll(f); err != nil {
		t.Fatal(err)
	}

	if ep, err = unmarshal(b); err != nil {
		t.Fatal(err)
	}

	if len(ep.Events) != 1 {
		t.Fatal("should have found 1 event.")
	}

	e := ep.Events[0]

	if e.PublicID != "2016p408314" {
		t.Errorf("expected publicID 2016p408314 got %s", e.PublicID)
	}

	if e.Type != "earthquake" {
		t.Errorf("expected type earthquake got %s", e.Type)
	}

	if e.CreationInfo.AgencyID != "WEL(GNS_Test)" {
		t.Errorf("AgencyID expected WEL(GNS_Test) got %s", e.CreationInfo.AgencyID)
	}

	if e.CreationInfo.ModificationTime.Format(time.RFC3339Nano) != "2016-06-01T04:31:27.60558Z" {
		t.Errorf("ModificationTime expected 2016-06-01T04:31:27.60558Z got %s", e.CreationInfo.ModificationTime.Format(time.RFC3339Nano))
	}

	if e.PreferredOrigin.Time.Value.Format(time.RFC3339Nano) != "2016-05-31T01:50:12.062388Z" {
		t.Errorf("expected 2016-05-31T01:50:12.062388Z, got %s", e.PreferredOrigin.Time.Value.Format(time.RFC3339Nano))
	}

	if e.PreferredOrigin.Latitude.Value != -45.19537735 {
		t.Errorf("Latitude expected -45.19537735 got %f", e.PreferredOrigin.Latitude.Value)
	}

	if e.PreferredOrigin.Longitude.Value != 167.3780823 {
		t.Errorf("Longitude expected 167.3780823 got %f", e.PreferredOrigin.Longitude.Value)
	}

	if e.PreferredOrigin.Depth.Value != 100.126976 {
		t.Errorf("Depth expected 100.126976 got %f", e.PreferredOrigin.Depth.Value)
	}

	if e.PreferredOrigin.MethodID != "LOCSAT" {
		t.Errorf("MethodID expected LOCSAT got %s", e.PreferredOrigin.MethodID)
	}

	if e.PreferredOrigin.EarthModelID != "iasp91" {
		t.Errorf("EarthModelID expected iasp91 got %s", e.PreferredOrigin.EarthModelID)
	}

	if e.PreferredOrigin.Quality.AzimuthalGap != 186.5389404 {
		t.Errorf("AzimuthalGap expected 186.5389404 got %f", e.PreferredOrigin.Quality.AzimuthalGap)
	}

	if e.PreferredOrigin.Quality.MinimumDistance != 0.3124738038 {
		t.Errorf("MinimumDistance expected 0.3124738038 got %f", e.PreferredOrigin.Quality.MinimumDistance)
	}

	if e.PreferredOrigin.Quality.UsedPhaseCount != 18 {
		t.Errorf("UsedPhaseCount expected 44 got %d", e.PreferredOrigin.Quality.UsedPhaseCount)
	}

	if e.PreferredOrigin.Quality.UsedStationCount != 14 {
		t.Errorf("UsedStationCount expected 32 got %d", e.PreferredOrigin.Quality.UsedStationCount)
	}

	if e.PreferredMagnitude.Magnitude.Value != 4.452756951 {
		t.Errorf("Magnitude expected 4.452756951 got %f", e.PreferredMagnitude.Magnitude.Value)
	}

	if e.PreferredMagnitude.Type != "Mw" {
		t.Errorf("Magnitude type expected Mw got %s", e.PreferredMagnitude.Type)
	}

	if e.PreferredMagnitude.StationCount != 19 {
		t.Errorf("Expected StationCount 19 gor %d", e.PreferredMagnitude.StationCount)
	}
}

func BenchmarkUnmarshalSeiscompml(b *testing.B) {
	var err error
	var ep eventParameters
	var f *os.File
	var by []byte

	if f, err = os.Open("etc/2015p768477.xml"); err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	if by, err = ioutil.ReadAll(f); err != nil {
		b.Fatal(err)
	}

	if ep, err = unmarshal(by); err != nil {
		b.Fatal(err)
	}

	for n := 0; n < b.N; n++ {
		// ignore errors
		ep, _ = unmarshal(by)
	}

	_ = ep

}
