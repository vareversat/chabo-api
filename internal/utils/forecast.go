package utils

import (
	"log"
	"os"
	"time"

	"github.com/vareversat/chabo-api/internal/models"
)

var (
	ErrorLogger = log.New(os.Stdout, "ERROR: ", log.LUTC|log.Ltime|log.Lshortfile)
)

// Populate a []models.Forecast pointer with the OpenAPI data
func ComputeForecasts(forecasts *[]models.Forecast, openDataForecasts models.OpenDataAPIResponse) {
	// alreadySeenBoatNames is used to compute the maneuver of each boats
	var alreadySeenBoatNames []string

	for _, openAPIForecast := range openDataForecasts.Records {
		_, offset := openAPIForecast.RecordTimestamp.Zone()
		closingReason := MapClosingReason(openAPIForecast.Fields.Boat)
		circulationClosingDate, errClosingDate := FormatDataTime(
			openAPIForecast.Fields.ClosingTime,
			openAPIForecast.Fields.ClosingDate,
			offset,
			*time.UTC,
		)
		circulationReopeningDate, errReopeningDate := FormatDataTime(
			openAPIForecast.Fields.OpeningTime,
			openAPIForecast.Fields.ClosingDate,
			offset,
			*time.UTC,
		)

		if errClosingDate != nil || errReopeningDate != nil {
			ErrorLogger.Printf(errClosingDate.Error(), '\n', errReopeningDate.Error())
		}

		// Check if the forecast is during 2 days
		if circulationReopeningDate.Compare(circulationClosingDate) == -1 {
			// On day is added because the closing date is after the reopening date
			circulationReopeningDate = circulationReopeningDate.AddDate(0, 0, 1)
		}
		closingDuration := circulationReopeningDate.Sub(circulationClosingDate)
		*forecasts = append(*forecasts, models.Forecast{
			ID:                       openAPIForecast.RecordID,
			ClosingDuration:          closingDuration,
			CirculationClosingDate:   circulationClosingDate,
			CirculationReopeningDate: circulationReopeningDate,
			ClosingType:              MapClosingType(openAPIForecast.Fields.TotalClosing),
			ClosingReason:            closingReason,
			Boats: MapBoats(
				closingReason,
				openAPIForecast.Fields.Boat,
				closingDuration,
				circulationClosingDate,
				&alreadySeenBoatNames,
				openAPIForecast.RecordID,
			),
			Link: models.OpenAPISelfLink{
				Self: models.OpenAPILink{Link: "/v1/forecasts/" + openAPIForecast.RecordID},
			},
		})
	}

}
