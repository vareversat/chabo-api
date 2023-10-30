package usecases

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vareversat/chabo-api/internal/domains"
)

func TestComputeForecasts(t *testing.T) {
	forecastRepository := new(domains.ForecastRepository)
	refreshRepository := new(domains.RefreshRepository)

	var forecasts domains.Forecasts
	recordTimestamp, _ := time.Parse(time.RFC3339, "2023-02-26T21:00:00Z")
	bordeauxAPIForecasts := domains.BordeauxAPIResponse{
		Hits: 1,
		Parameters: domains.BordeauxAPIResponseParameters{
			Dataset:  "dataset",
			Row:      1,
			Start:    0,
			Format:   "format",
			Timezone: "UTC",
		},
		Records: []domains.BordeauxAPIResponseForecast{
			{
				DatasetID:       "datasetid",
				RecordID:        "recordid",
				RecordTimestamp: recordTimestamp,
				Fields: domains.BordeauxAPIResponseForecastField{
					ClosingDate:  "2023-02-26",
					ClosingTime:  "21:00",
					OpeningTime:  "23:00",
					TotalClosing: "oui",
					Boat:         "MY_BOAT",
					ClosingType:  "oui",
				},
			},
		},
	}
	circulationClosingDate, _ := time.Parse(time.RFC3339, "2023-02-26T21:00:00Z")
	circulationReopeningDate, _ := time.Parse(time.RFC3339, "2023-02-26T23:00:00Z")
	approximativeCrossingDate, _ := time.Parse(time.RFC3339, "2023-02-26T22:00:00Z")
	expectedForecasts := domains.Forecasts{
		{
			ID:                       "recordid",
			ClosingType:              domains.TwoWay,
			ClosingDuration:          7200000000000,
			CirculationClosingDate:   circulationClosingDate,
			CirculationReopeningDate: circulationReopeningDate,
			ClosingReason:            domains.BoatReason,
			Boats: []domains.Boat{
				{
					Name:                      "MY_BOAT",
					Maneuver:                  domains.Entering,
					ApproximativeCrossingDate: approximativeCrossingDate,
				},
			},
			Link: domains.APIResponseSelfLink{
				Self: domains.APIResponseLink{Link: "/v1/forecasts/recordid"},
			},
		},
	}

	u := NewForecastUsecase(*forecastRepository, *refreshRepository, time.Second*2)
	u.ComputeBordeauxAPIResponse(&forecasts, bordeauxAPIForecasts)

	assert.True(t, true, expectedForecasts.AreEqual(forecasts))
}
