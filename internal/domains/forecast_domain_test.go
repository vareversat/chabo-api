package domains

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestForecastIsEqualOK(t *testing.T) {
	crossingDateApproximation, _ := time.Parse(time.RFC3339, "2023-02-26T21:00:00+01:00")
	circulationClosingDate, _ := time.Parse(time.RFC3339, "2023-02-26T21:00:00Z")
	circulationReopeningDate, _ := time.Parse(time.RFC3339, "2023-02-26T23:00:00Z")
	forecast := Forecast{
		ID:                       "recordId",
		IsTrafficFullyClosed:     true,
		ClosingDuration:          7200000000000,
		CirculationClosingDate:   circulationClosingDate,
		CirculationReopeningDate: circulationReopeningDate,
		ClosingReason:            BoatPassage,
		Boats: []Boat{
			{
				Name:                      "MY_BOAT",
				IsLeavingDock:             false,
				ApproximativeCrossingDate: crossingDateApproximation,
			},
		},
		Link: APIResponseSelfLink{Self: APIResponseLink{Link: "/v1/forecasts/recordId"}},
	}
	otherForecast := Forecast{
		ID:                       "recordId",
		IsTrafficFullyClosed:     true,
		ClosingDuration:          7200000000000,
		CirculationClosingDate:   circulationClosingDate,
		CirculationReopeningDate: circulationReopeningDate,
		ClosingReason:            BoatPassage,
		Boats: []Boat{
			{
				Name:                      "MY_BOAT",
				IsLeavingDock:             false,
				ApproximativeCrossingDate: crossingDateApproximation,
			},
		},
		Link: APIResponseSelfLink{Self: APIResponseLink{Link: "/v1/forecasts/recordId"}},
	}
	result := forecast.IsEqual(otherForecast)

	assert.True(t, result)
}

func TestForecastIsEqualNOK(t *testing.T) {
	crossingDateApproximation, _ := time.Parse(time.RFC3339, "2023-02-26T21:00:00+01:00")
	circulationClosingDate, _ := time.Parse(time.RFC3339, "2023-02-26T21:00:00Z")
	circulationReopeningDate, _ := time.Parse(time.RFC3339, "2023-02-26T23:00:00Z")
	forecast := Forecast{
		ID:                       "recordId",
		IsTrafficFullyClosed:     true,
		ClosingDuration:          7200000000000,
		CirculationClosingDate:   circulationClosingDate,
		CirculationReopeningDate: circulationReopeningDate,
		ClosingReason:            BoatPassage,
		Boats: []Boat{
			{
				Name:                      "MY_BOAT2",
				IsLeavingDock:             false,
				ApproximativeCrossingDate: crossingDateApproximation,
			},
		},
		Link: APIResponseSelfLink{Self: APIResponseLink{Link: "/v1/forecasts/recordId"}},
	}
	otherForecast := Forecast{
		ID:                       "recordId",
		IsTrafficFullyClosed:     true,
		ClosingDuration:          7200000000000,
		CirculationClosingDate:   circulationClosingDate,
		CirculationReopeningDate: circulationReopeningDate,
		ClosingReason:            BoatPassage,
		Boats: []Boat{
			{
				Name:                      "MY_BOAT",
				IsLeavingDock:             false,
				ApproximativeCrossingDate: crossingDateApproximation,
			},
		},
		Link: APIResponseSelfLink{Self: APIResponseLink{Link: "/v1/forecasts/recordId"}},
	}
	result := forecast.IsEqual(otherForecast)

	assert.False(t, result)
}

func TestForecastsAreEqualOK(t *testing.T) {
	crossingDateApproximation, _ := time.Parse(time.RFC3339, "2023-02-26T21:00:00+01:00")
	circulationClosingDate, _ := time.Parse(time.RFC3339, "2023-02-26T21:00:00Z")
	circulationReopeningDate, _ := time.Parse(time.RFC3339, "2023-02-26T23:00:00Z")
	forecasts := Forecasts{Forecast{
		ID:                       "recordId",
		IsTrafficFullyClosed:     true,
		ClosingDuration:          7200000000000,
		CirculationClosingDate:   circulationClosingDate,
		CirculationReopeningDate: circulationReopeningDate,
		ClosingReason:            BoatPassage,
		Boats: []Boat{
			{
				Name:                      "MY_BOAT",
				IsLeavingDock:             false,
				ApproximativeCrossingDate: crossingDateApproximation,
			},
		},
		Link: APIResponseSelfLink{Self: APIResponseLink{Link: "/v1/forecasts/recordId"}},
	}}
	otherForecasts := Forecasts{Forecast{
		ID:                       "recordId",
		IsTrafficFullyClosed:     true,
		ClosingDuration:          7200000000000,
		CirculationClosingDate:   circulationClosingDate,
		CirculationReopeningDate: circulationReopeningDate,
		ClosingReason:            BoatPassage,
		Boats: []Boat{
			{
				Name:                      "MY_BOAT",
				IsLeavingDock:             false,
				ApproximativeCrossingDate: crossingDateApproximation,
			},
		},
		Link: APIResponseSelfLink{Self: APIResponseLink{Link: "/v1/forecasts/recordId"}},
	}}
	result := forecasts.AreEqual(otherForecasts)

	assert.True(t, result)
}

func TestForecastsAreEqualNOK(t *testing.T) {
	crossingDateApproximation, _ := time.Parse(time.RFC3339, "2023-02-26T21:00:00+01:00")
	circulationClosingDate, _ := time.Parse(time.RFC3339, "2023-02-26T21:00:00Z")
	circulationReopeningDate, _ := time.Parse(time.RFC3339, "2023-02-26T23:00:00Z")
	forecasts := Forecasts{Forecast{
		ID:                       "recordId",
		IsTrafficFullyClosed:     true,
		ClosingDuration:          7200000000000,
		CirculationClosingDate:   circulationClosingDate,
		CirculationReopeningDate: circulationReopeningDate,
		ClosingReason:            BoatPassage,
		Boats: []Boat{
			{
				Name:                      "MY_BOAT2",
				IsLeavingDock:             false,
				ApproximativeCrossingDate: crossingDateApproximation,
			},
		},
		Link: APIResponseSelfLink{Self: APIResponseLink{Link: "/v1/forecasts/recordId"}},
	}}
	otherForecasts := Forecasts{Forecast{
		ID:                       "recordId2",
		IsTrafficFullyClosed:     true,
		ClosingDuration:          7200000000000,
		CirculationClosingDate:   circulationClosingDate,
		CirculationReopeningDate: circulationReopeningDate,
		ClosingReason:            BoatPassage,
		Boats: []Boat{
			{
				Name:                      "MY_BOAT2",
				IsLeavingDock:             false,
				ApproximativeCrossingDate: crossingDateApproximation,
			},
		},
		Link: APIResponseSelfLink{Self: APIResponseLink{Link: "/v1/forecasts/recordId"}},
	}}
	result := forecasts.AreEqual(otherForecasts)

	assert.False(t, result)
}
