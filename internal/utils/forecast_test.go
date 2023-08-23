package utils

import (
	"testing"
	"time"

	"github.com/vareversat/chabo-api/internal/models"
)

func TestComputeForecasts(t *testing.T) {
	var forecasts models.Forecasts
	recordTimestamp, _ := time.Parse(time.RFC3339, "2023-02-26T21:00:00Z")
	openDataForecasts := models.OpenDataAPIResponse{
		Hits: 1,
		Parameters: models.OpenDataAPIResponseParameters{
			Dataset:  "dataset",
			Row:      1,
			Start:    0,
			Format:   "format",
			Timezone: "UTC",
		},
		Records: []models.OpenDataAPIResponseForecast{
			{
				DatasetID:       "datasetid",
				RecordID:        "recordid",
				RecordTimestamp: recordTimestamp,
				Fields: models.OpenDataAPIResponseForecastField{
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
	want := models.Forecasts{
		{
			ID:                       "recordid",
			ClosingType:              models.TwoWay,
			ClosingDuration:          7200000000000,
			CirculationClosingDate:   circulationClosingDate,
			CirculationReopeningDate: circulationReopeningDate,
			ClosingReason:            models.BoatReason,
			Boats: []models.Boat{
				{
					Name:                      "MY_BOAT",
					Maneuver:                  models.Entering,
					ApproximativeCrossingDate: approximativeCrossingDate,
				},
			},
			Link: models.OpenAPISelfLink{Self: models.OpenAPILink{Link: "/v1/forecasts/recordid"}},
		},
	}

	ComputeForecasts(&forecasts, openDataForecasts)
	if !(want.AreEqual(forecasts)) {
		t.Fatalf(`ComputeForecasts("...") = %q, want match for %#q`, forecasts, want)
	}
}
