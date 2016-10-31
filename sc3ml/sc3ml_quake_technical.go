package sc3ml

import (
	"fmt"
	"github.com/GeoNet/haz"
)

// QuakeTechnical converts the SC3ML in b to a haz.QuakeTechnical.
//
// The current supported SC3ML version is 0.7
//e.g., http://geofon.gfz-potsdam.de/schema/0.7/sc3ml_0.7.xsd
//
//The schema allows many elements to be 0..1 or 0..* The elements
//mapped here are assumed to be present and will have zero values
//if not.
//
// It is an error if b does not contain exactly one event.
func QuakeTechnical(b []byte) (haz.QuakeTechnical, error) {
	ep, err := unmarshal(b)
	if err != nil {
		return haz.QuakeTechnical{}, err
	}

	if len(ep.Events) != 1 {
		return haz.QuakeTechnical{}, fmt.Errorf("did not find exactly 1 event: got %d", len(ep.Events))
	}

	e := ep.Events[0]

	mt := e.modificationTime()

	h := haz.QuakeTechnical{
		PublicID: e.PublicID,
		Type:     e.Type,
		Time: &haz.Timestamp{
			Sec:  e.PreferredOrigin.Time.Value.Unix(),
			Nsec: int64(e.PreferredOrigin.Time.Value.Nanosecond()),
		},
		ModificationTime:      &haz.Timestamp{
			Sec: mt.Unix(),
			Nsec: int64(mt.Nanosecond()),
		},
		Latitude: &haz.RealQuantity{
			Value: e.PreferredOrigin.Latitude.Value,
			Uncertainty:e.PreferredOrigin.Latitude.Uncertainty,
		},
		Longitude: &haz.RealQuantity{
			Value: e.PreferredOrigin.Longitude.Value,
			Uncertainty:e.PreferredOrigin.Longitude.Uncertainty,
		},
		Depth: &haz.RealQuantity{
			Value: e.PreferredOrigin.Depth.Value,
			Uncertainty:e.PreferredOrigin.Depth.Uncertainty,
		},
		DepthType:             e.PreferredOrigin.DepthType,
		Method:                e.PreferredOrigin.MethodID,
		EarthModel:            e.PreferredOrigin.EarthModelID,
		EvaluationMode:        e.PreferredOrigin.EvaluationMode,
		EvaluationStatus:      e.PreferredOrigin.EvaluationStatus,
		UsedPhaseCount:        e.PreferredOrigin.Quality.UsedPhaseCount,
		UsedStationCount:      e.PreferredOrigin.Quality.UsedStationCount,
		StandardError:         e.PreferredOrigin.Quality.StandardError,
		AzimuthalGap:          e.PreferredOrigin.Quality.AzimuthalGap,
		MinimumDistance:       e.PreferredOrigin.Quality.MinimumDistance,
		Magnitude: &haz.RealQuantity{
			Value: e.PreferredMagnitude.Magnitude.Value,
			Uncertainty:e.PreferredMagnitude.Magnitude.Uncertainty,
		},
		MagnitudeType:         e.PreferredMagnitude.Type,
		Agency:                e.CreationInfo.AgencyID,
	}

	for _, v := range e.PreferredOrigin.Arrivals {
		p := haz.Pick{
			Waveform: &haz.Waveform{
				Network:      v.Pick.WaveformID.NetworkCode,
				Station:      v.Pick.WaveformID.StationCode,
				Location:     v.Pick.WaveformID.LocationCode,
				Channel:      v.Pick.WaveformID.ChannelCode,
			},

			Phase:            v.Phase,
			Time:             &haz.Timestamp{Sec: v.Pick.Time.Value.Unix(), Nsec: int64(v.Pick.Time.Value.Nanosecond())},
			Residual:         v.TimeResidual,
			Weight:           v.Weight,
			Azimuth:          v.Azimuth,
			Distance:         v.Distance,
			EvaluationMode:   v.Pick.EvaluationMode,
			EvaluationStatus: v.Pick.EvaluationStatus,
		}

		h.Pick = append(h.Pick, &p)
	}

	for _, v := range e.PreferredOrigin.Magnitudes {
		m := haz.Magnitude{
			Magnitude: &haz.RealQuantity{
				Value:            v.Magnitude.Value,
				Uncertainty: v.Magnitude.Uncertainty,
			},
			Type:                 v.Type,
			StationCount:         v.StationCount,
		}

		for _, vv := range v.StationMagnitudeContributions {
			sm := haz.StationMagnitude{
				Waveform: &haz.Waveform{
					Network:  vv.StationMagnitude.WaveformID.NetworkCode,
					Station:  vv.StationMagnitude.WaveformID.StationCode,
					Location: vv.StationMagnitude.WaveformID.LocationCode,
					Channel:  vv.StationMagnitude.WaveformID.ChannelCode,
				},
				Magnitude:    &haz.RealQuantity{
					Value: vv.StationMagnitude.Magnitude.Value,
				},
				Type:         vv.StationMagnitude.Type,
				Azimuth: vv.StationMagnitude.Amplitude.Azimuth,
				Distance: vv.StationMagnitude.Amplitude.Distance,
				Residual:     vv.Residual,
				Weight:       vv.Weight,
				Amplitude:    &haz.RealQuantity{
					Value:vv.StationMagnitude.Amplitude.Amplitude.Value,
				},
			}
			m.StationMagnitude = append(m.StationMagnitude, &sm)
		}

		h.Magnitudes = append(h.Magnitudes, &m)
	}

	return h, nil
}
