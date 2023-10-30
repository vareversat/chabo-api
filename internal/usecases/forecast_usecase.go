package usecases

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vareversat/chabo-api/internal/domains"
	"github.com/vareversat/chabo-api/internal/utils"
)

type forecastUsecase struct {
	forecastRepository domains.ForecastRepository
	refreshRepository  domains.RefreshRepository
	contextTimeout     time.Duration
}

func NewForecastUsecase(
	forecastRepository domains.ForecastRepository,
	refreshRepository domains.RefreshRepository,
	timeout time.Duration,
) domains.ForecastUsecase {
	return &forecastUsecase{
		forecastRepository: forecastRepository,
		refreshRepository:  refreshRepository,
		contextTimeout:     timeout,
	}
}

func (fU *forecastUsecase) GetByID(
	ctx context.Context,
	id string,
	forecast *domains.Forecast,
	location *time.Location,
) error {
	ctx, cancel := context.WithTimeout(ctx, fU.contextTimeout)
	defer cancel()

	// Do a refresh attempt in case of the data are too old
	fU.RefreshAll(ctx)

	err := fU.forecastRepository.GetByID(ctx, id, forecast)
	if err != nil {
		return err
	}

	forecast.ChangeLocation(location)
	return nil
}

func (fU *forecastUsecase) GetTodayForecasts(
	ctx context.Context,
	forecasts *domains.Forecasts,
	offset int,
	limit int,
	location *time.Location,
	totalItemCount *int,
) error {
	ctx, cancel := context.WithTimeout(ctx, fU.contextTimeout)
	// Get the current time
	from := time.Now()
	// Convert the local time into the requested TZ
	fromLocal := from.In(location)
	// Get the first second of the current day
	fromRounded := time.Date(
		fromLocal.Year(),
		fromLocal.Month(),
		fromLocal.Day(),
		0,
		0,
		0,
		0,
		location,
	)
	// The 'to' time is computed by adding one day to the 'fromRounded' time
	to := fromRounded.AddDate(0, 0, 1).In(location)
	// Get the first second of the next day
	toRounded := time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, location)
	defer cancel()

	// Do a refresh attempt in case of the data are too old
	fU.RefreshAll(ctx)

	err := fU.forecastRepository.GetAllBetweenTwoDates(
		ctx,
		offset,
		limit,
		fromRounded,
		toRounded,
		forecasts,
		totalItemCount,
	)
	if err != nil {
		return err
	}

	forecasts.ChangeLocations(location)
	return nil
}

func (fU *forecastUsecase) GetAllFiltered(
	ctx context.Context,
	location *time.Location,
	offset int,
	limit int,
	from time.Time,
	reason string,
	maneuver string,
	boat string,
	forecasts *domains.Forecasts,
	totalItemCount *int,
) error {
	ctx, cancel := context.WithTimeout(ctx, fU.contextTimeout)
	defer cancel()

	// Do a refresh attempt in case of the data are too old
	fU.RefreshAll(ctx)

	err := fU.forecastRepository.GetAllFiltered(
		ctx,
		offset,
		limit,
		from,
		reason,
		maneuver,
		boat,
		forecasts,
		totalItemCount,
	)
	if err != nil {
		return err
	}

	forecasts.ChangeLocations(location)
	return nil
}

func (fU *forecastUsecase) RefreshAll(ctx context.Context) (domains.Refresh, error) {
	ctx, cancel := context.WithTimeout(ctx, fU.contextTimeout)
	defer cancel()

	if fU.RefreshIsNeeded(ctx) {
		var openDataForecasts domains.BordeauxAPIResponse
		var forecasts domains.Forecasts

		// Start the timer
		start := time.Now()
		// Fetch the fresh data
		errGet := utils.GetOpenAPIData(&openDataForecasts)
		if errGet != nil {
			return domains.Refresh{}, fmt.Errorf(errGet.Error())
		}
		// Compute all forecasts
		err := fU.ComputeBordeauxAPIResponse(&forecasts, openDataForecasts)
		if err != nil {
			return domains.Refresh{}, err
		}
		// Delete all forecasts
		_, err = fU.forecastRepository.DeleteAll(ctx)
		if err != nil {
			return domains.Refresh{}, err
		}
		// Insert all forecasts
		insertCount, err := fU.forecastRepository.InsertAll(ctx, forecasts)
		if err != nil {
			return domains.Refresh{}, err
		}
		// STOP the timer
		elapsed := time.Since(start)
		// Insert a refresh proof
		refresh := domains.Refresh{
			ItemCount: insertCount,
			Duration:  elapsed,
			Timestamp: start,
		}
		err = fU.refreshRepository.InsertOne(ctx, refresh)
		if err != nil {
			return domains.Refresh{}, err
		}
		return refresh, nil

	}

	return domains.Refresh{}, fmt.Errorf("data does not need to be refresh (aborting the refresh)")
}

func (fU *forecastUsecase) ComputeBordeauxAPIResponse(
	forecasts *domains.Forecasts,
	boredeauxAPIResponse domains.BordeauxAPIResponse,
) error {
	// alreadySeenBoatNames is used to compute the maneuver of each boats
	var alreadySeenBoatNames []string

	for _, openAPIForecast := range boredeauxAPIResponse.Records {
		_, offset := openAPIForecast.RecordTimestamp.Zone()
		closingReason := utils.MapClosingReason(openAPIForecast.Fields.Boat)
		circulationClosingDate, errClosingDate := utils.FormatDataTime(
			openAPIForecast.Fields.ClosingTime,
			openAPIForecast.Fields.ClosingDate,
			offset,
			*time.UTC,
		)
		if errClosingDate != nil {
			logrus.Fatalf(errClosingDate.Error())
			return errClosingDate
		}
		circulationReopeningDate, errReopeningDate := utils.FormatDataTime(
			openAPIForecast.Fields.OpeningTime,
			openAPIForecast.Fields.ClosingDate,
			offset,
			*time.UTC,
		)
		if errReopeningDate != nil {
			logrus.Fatalf(errReopeningDate.Error())
			return errReopeningDate
		}

		// Check if the forecast is during 2 days
		if circulationReopeningDate.Compare(circulationClosingDate) == -1 {
			// On day is added because the closing date is after the reopening date
			circulationReopeningDate = circulationReopeningDate.AddDate(0, 0, 1)
		}
		closingDuration := circulationReopeningDate.Sub(circulationClosingDate)
		*forecasts = append(*forecasts, domains.Forecast{
			ID:                       openAPIForecast.RecordID,
			ClosingDuration:          closingDuration,
			CirculationClosingDate:   circulationClosingDate,
			CirculationReopeningDate: circulationReopeningDate,
			ClosingType:              utils.MapClosingType(openAPIForecast.Fields.TotalClosing),
			ClosingReason:            closingReason,
			Boats: utils.MapBoats(
				closingReason,
				openAPIForecast.Fields.Boat,
				closingDuration,
				circulationClosingDate,
				&alreadySeenBoatNames,
				openAPIForecast.RecordID,
			),
			Link: domains.APIResponseSelfLink{
				Self: domains.APIResponseLink{Link: "/v1/forecasts/" + openAPIForecast.RecordID},
			},
		})
	}
	logrus.Infof("all %d forecasts computed with success", len(*forecasts))
	return nil

}

// Check if it's possible to perform a data refresh
func (fU *forecastUsecase) RefreshIsNeeded(ctx context.Context) bool {

	var lastRefresh domains.Refresh

	// Get the last refresh to be sure this is not too early
	err := fU.refreshRepository.GetLast(ctx, &lastRefresh)

	if err != nil {
		// An error here means that the collection is empty
		return true
	} else {
		currentTime := time.Now()
		diff := currentTime.Sub(lastRefresh.Timestamp)

		cooldown, _ := strconv.Atoi(os.Getenv("REFRESH_COOLDOWN_SECONDS"))

		return int(diff.Seconds()) >= cooldown
	}

}
