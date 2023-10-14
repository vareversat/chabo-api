package domains

import (
	"testing"
	"time"
)

func TestForecastIsEqualOK(t *testing.T) {
	approximativeCrossingDate, _ := time.Parse(time.RFC3339, "2023-02-26T21:00:00+01:00")
	circulationClosingDate, _ := time.Parse(time.RFC3339, "2023-02-26T21:00:00Z")
	circulationReopeningDate, _ := time.Parse(time.RFC3339, "2023-02-26T23:00:00Z")
	forecast := Forecast{
		ID:                       "recordid",
		ClosingType:              TwoWay,
		ClosingDuration:          7200000000000,
		CirculationClosingDate:   circulationClosingDate,
		CirculationReopeningDate: circulationReopeningDate,
		ClosingReason:            BoatReason,
		Boats: []Boat{
			{
				Name:                      "MY_BOAT",
				Maneuver:                  Entering,
				ApproximativeCrossingDate: approximativeCrossingDate,
			},
		},
		Link: OpenAPISelfLink{Self: OpenAPILink{Link: "/v1/forecasts/recordid"}},
	}
	otherForecast := Forecast{
		ID:                       "recordid",
		ClosingType:              TwoWay,
		ClosingDuration:          7200000000000,
		CirculationClosingDate:   circulationClosingDate,
		CirculationReopeningDate: circulationReopeningDate,
		ClosingReason:            BoatReason,
		Boats: []Boat{
			{
				Name:                      "MY_BOAT",
				Maneuver:                  Entering,
				ApproximativeCrossingDate: approximativeCrossingDate,
			},
		},
		Link: OpenAPISelfLink{Self: OpenAPILink{Link: "/v1/forecasts/recordid"}},
	}
	result := forecast.IsEqual(otherForecast)
	want := true

	if want != result {
		t.Fatalf(`IsEqual("otherForecast") = %v, want match for %v`, result, want)
	}
}

func TestForecastIsEqualNOK(t *testing.T) {
	approximativeCrossingDate, _ := time.Parse(time.RFC3339, "2023-02-26T21:00:00+01:00")
	circulationClosingDate, _ := time.Parse(time.RFC3339, "2023-02-26T21:00:00Z")
	circulationReopeningDate, _ := time.Parse(time.RFC3339, "2023-02-26T23:00:00Z")
	forecast := Forecast{
		ID:                       "recordid",
		ClosingType:              TwoWay,
		ClosingDuration:          7200000000000,
		CirculationClosingDate:   circulationClosingDate,
		CirculationReopeningDate: circulationReopeningDate,
		ClosingReason:            BoatReason,
		Boats: []Boat{
			{
				Name:                      "MY_BOAT2",
				Maneuver:                  Entering,
				ApproximativeCrossingDate: approximativeCrossingDate,
			},
		},
		Link: OpenAPISelfLink{Self: OpenAPILink{Link: "/v1/forecasts/recordid"}},
	}
	otherForecast := Forecast{
		ID:                       "recordid",
		ClosingType:              TwoWay,
		ClosingDuration:          7200000000000,
		CirculationClosingDate:   circulationClosingDate,
		CirculationReopeningDate: circulationReopeningDate,
		ClosingReason:            BoatReason,
		Boats: []Boat{
			{
				Name:                      "MY_BOAT",
				Maneuver:                  Entering,
				ApproximativeCrossingDate: approximativeCrossingDate,
			},
		},
		Link: OpenAPISelfLink{Self: OpenAPILink{Link: "/v1/forecasts/recordid"}},
	}
	result := forecast.IsEqual(otherForecast)
	want := false

	if want != result {
		t.Fatalf(`IsEqual("otherForecast") = %v, want match for %v`, result, want)
	}
}

func TestForecastsAreEqualOK(t *testing.T) {
	approximativeCrossingDate, _ := time.Parse(time.RFC3339, "2023-02-26T21:00:00+01:00")
	circulationClosingDate, _ := time.Parse(time.RFC3339, "2023-02-26T21:00:00Z")
	circulationReopeningDate, _ := time.Parse(time.RFC3339, "2023-02-26T23:00:00Z")
	forecasts := Forecasts{Forecast{
		ID:                       "recordid",
		ClosingType:              TwoWay,
		ClosingDuration:          7200000000000,
		CirculationClosingDate:   circulationClosingDate,
		CirculationReopeningDate: circulationReopeningDate,
		ClosingReason:            BoatReason,
		Boats: []Boat{
			{
				Name:                      "MY_BOAT",
				Maneuver:                  Entering,
				ApproximativeCrossingDate: approximativeCrossingDate,
			},
		},
		Link: OpenAPISelfLink{Self: OpenAPILink{Link: "/v1/forecasts/recordid"}},
	}}
	otherForecasts := Forecasts{Forecast{
		ID:                       "recordid",
		ClosingType:              TwoWay,
		ClosingDuration:          7200000000000,
		CirculationClosingDate:   circulationClosingDate,
		CirculationReopeningDate: circulationReopeningDate,
		ClosingReason:            BoatReason,
		Boats: []Boat{
			{
				Name:                      "MY_BOAT",
				Maneuver:                  Entering,
				ApproximativeCrossingDate: approximativeCrossingDate,
			},
		},
		Link: OpenAPISelfLink{Self: OpenAPILink{Link: "/v1/forecasts/recordid"}},
	}}
	result := forecasts.AreEqual(otherForecasts)
	want := true

	if want != result {
		t.Fatalf(`AreEqual("otherForecasts") = %v, want match for %v`, result, want)
	}
}

func TestForecastsAreEqualNOK(t *testing.T) {
	approximativeCrossingDate, _ := time.Parse(time.RFC3339, "2023-02-26T21:00:00+01:00")
	circulationClosingDate, _ := time.Parse(time.RFC3339, "2023-02-26T21:00:00Z")
	circulationReopeningDate, _ := time.Parse(time.RFC3339, "2023-02-26T23:00:00Z")
	forecasts := Forecasts{Forecast{
		ID:                       "recordid",
		ClosingType:              TwoWay,
		ClosingDuration:          7200000000000,
		CirculationClosingDate:   circulationClosingDate,
		CirculationReopeningDate: circulationReopeningDate,
		ClosingReason:            BoatReason,
		Boats: []Boat{
			{
				Name:                      "MY_BOAT2",
				Maneuver:                  Entering,
				ApproximativeCrossingDate: approximativeCrossingDate,
			},
		},
		Link: OpenAPISelfLink{Self: OpenAPILink{Link: "/v1/forecasts/recordid"}},
	}}
	otherForecasts := Forecasts{Forecast{
		ID:                       "recordid2",
		ClosingType:              TwoWay,
		ClosingDuration:          7200000000000,
		CirculationClosingDate:   circulationClosingDate,
		CirculationReopeningDate: circulationReopeningDate,
		ClosingReason:            BoatReason,
		Boats: []Boat{
			{
				Name:                      "MY_BOAT2",
				Maneuver:                  Entering,
				ApproximativeCrossingDate: approximativeCrossingDate,
			},
		},
		Link: OpenAPISelfLink{Self: OpenAPILink{Link: "/v1/forecasts/recordid"}},
	}}
	result := forecasts.AreEqual(otherForecasts)
	want := false

	if want != result {
		t.Fatalf(`AreEqual("otherForecasts") = %v, want match for %v`, result, want)
	}
}
