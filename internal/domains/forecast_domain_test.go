package domains

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
		Link: APIResponseSelfLink{Self: APIResponseLink{Link: "/v1/forecasts/recordid"}},
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
		Link: APIResponseSelfLink{Self: APIResponseLink{Link: "/v1/forecasts/recordid"}},
	}
	result := forecast.IsEqual(otherForecast)

	assert.True(t, result)
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
		Link: APIResponseSelfLink{Self: APIResponseLink{Link: "/v1/forecasts/recordid"}},
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
		Link: APIResponseSelfLink{Self: APIResponseLink{Link: "/v1/forecasts/recordid"}},
	}
	result := forecast.IsEqual(otherForecast)

	assert.False(t, result)
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
		Link: APIResponseSelfLink{Self: APIResponseLink{Link: "/v1/forecasts/recordid"}},
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
		Link: APIResponseSelfLink{Self: APIResponseLink{Link: "/v1/forecasts/recordid"}},
	}}
	result := forecasts.AreEqual(otherForecasts)

	assert.True(t, result)
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
		Link: APIResponseSelfLink{Self: APIResponseLink{Link: "/v1/forecasts/recordid"}},
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
		Link: APIResponseSelfLink{Self: APIResponseLink{Link: "/v1/forecasts/recordid"}},
	}}
	result := forecasts.AreEqual(otherForecasts)

	assert.False(t, result)
}
