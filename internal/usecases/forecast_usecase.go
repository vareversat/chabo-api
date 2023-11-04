package usecases

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vareversat/chabo-api/internal/domains"
	"github.com/vareversat/chabo-api/internal/errors"
	"github.com/vareversat/chabo-api/internal/utils"
)

type forecastUsecase struct {
	forecastRepository domains.ForecastRepository
	syncRepository     domains.SyncRepository
	contextTimeout     time.Duration
}

func NewForecastUsecase(
	forecastRepository domains.ForecastRepository,
	syncRepository domains.SyncRepository,
	timeout time.Duration,
) domains.ForecastUsecase {
	return &forecastUsecase{
		forecastRepository: forecastRepository,
		syncRepository:     syncRepository,
		contextTimeout:     timeout,
	}
}

func (fU *forecastUsecase) GetByID(
	ctx context.Context,
	id string,
	forecast *domains.Forecast,
	location *time.Location,
) errors.CustomError {
	ctx, cancel := context.WithTimeout(ctx, fU.contextTimeout)
	defer cancel()

	// Do a sync attempt in case of the data are too old
	fU.SyncAll(ctx)

	err := fU.forecastRepository.GetByID(ctx, id, forecast)
	if err != nil {
		return errors.NewNotFoundError("Forecast not found in the database")
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
) errors.CustomError {
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

	// Do a sync attempt in case of the data are too old
	fU.SyncAll(ctx)

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
		return errors.NewInternalServerError(err.Error())
	}

	if len(*forecasts) == 0 {
		return errors.NewNotFoundError("No closing scheduled for today")
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
) errors.CustomError {
	ctx, cancel := context.WithTimeout(ctx, fU.contextTimeout)
	defer cancel()

	// Do a sync attempt in case of the data are too old
	fU.SyncAll(ctx)

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
		return errors.NewInternalServerError(err.Error())
	}

	if len(*forecasts) == 0 {
		return errors.NewNotFoundError("Forecasts not found in the database")
	}

	forecasts.ChangeLocations(location)
	return nil
}

func (fU *forecastUsecase) SyncAll(ctx context.Context) (domains.Sync, errors.CustomError) {
	ctx, cancel := context.WithTimeout(ctx, fU.contextTimeout)
	defer cancel()

	if fU.SyncIsNeeded(ctx) {
		var openDataForecasts domains.BordeauxAPIResponse
		var forecasts domains.Forecasts

		// Start the timer
		start := time.Now()
		// Fetch the fresh data
		errGet := utils.GetOpenAPIData(&openDataForecasts)
		if errGet != nil {
			return domains.Sync{}, errors.NewInternalServerError(errGet.Error())
		}
		// Compute all forecasts
		customError := fU.ComputeBordeauxAPIResponse(&forecasts, openDataForecasts)
		if customError != nil {
			return domains.Sync{}, customError
		}
		// Delete all forecasts
		_, err := fU.forecastRepository.DeleteAll(ctx)
		if err != nil {
			return domains.Sync{}, errors.NewInternalServerError(err.Error())
		}
		// Insert all forecasts
		insertCount, err := fU.forecastRepository.InsertAll(ctx, forecasts)
		if err != nil {
			return domains.Sync{}, errors.NewInternalServerError(err.Error())
		}
		// STOP the timer
		elapsed := time.Since(start)
		// Insert a sync proof
		sync := domains.Sync{
			ItemCount: insertCount,
			Duration:  time.Duration(elapsed.Milliseconds()),
			Timestamp: start,
		}
		err = fU.syncRepository.InsertOne(ctx, sync)
		if err != nil {
			return domains.Sync{}, errors.NewInternalServerError(err.Error())
		}
		return sync, nil

	}

	return domains.Sync{}, errors.NewSyncAttemptTooEarlyError()
}

func (fU *forecastUsecase) ComputeBordeauxAPIResponse(
	forecasts *domains.Forecasts,
	boredeauxAPIResponse domains.BordeauxAPIResponse,
) errors.CustomError {
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
			return errors.NewInternalServerError(errClosingDate.Error())
		}
		circulationReopeningDate, errReopeningDate := utils.FormatDataTime(
			openAPIForecast.Fields.OpeningTime,
			openAPIForecast.Fields.ClosingDate,
			offset,
			*time.UTC,
		)
		if errReopeningDate != nil {
			logrus.Fatalf(errReopeningDate.Error())
			return errors.NewInternalServerError(errReopeningDate.Error())
		}

		// Check if the forecast is during 2 days
		if circulationReopeningDate.Compare(circulationClosingDate) == -1 {
			// On day is added because the closing date is after the reopening date
			circulationReopeningDate = circulationReopeningDate.AddDate(0, 0, 1)
		}
		closingDuration := circulationReopeningDate.Sub(circulationClosingDate)
		*forecasts = append(*forecasts, domains.Forecast{
			ID: openAPIForecast.RecordID,
			ClosingDuration: time.Duration(
				closingDuration.Minutes(),
			), // Convert into minutes
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
	return nil

}

// Check if it's possible to perform a data sync
func (fU *forecastUsecase) SyncIsNeeded(ctx context.Context) bool {

	var lastSync domains.Sync

	// Get the last sync to be sure this is not too early
	err := fU.syncRepository.GetLast(ctx, &lastSync)

	if err != nil {
		// An error here means that the collection is empty
		return true
	} else {
		currentTime := time.Now()
		diff := currentTime.Sub(lastSync.Timestamp)

		cooldown, _ := strconv.Atoi(os.Getenv("SYNC_COOLDOWN_SECONDS"))

		return int(diff.Seconds()) >= cooldown
	}

}
