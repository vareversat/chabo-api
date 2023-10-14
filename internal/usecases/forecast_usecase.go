package usecases

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

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
	_, err := fU.RefreshAll(ctx)
	fmt.Println(err)

	err = fU.forecastRepository.GetByID(ctx, id, forecast)
	if err != nil {
		return err
	}

	forecast.ChangeLocation(location)
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
	_, err := fU.RefreshAll(ctx)
	fmt.Println(err)

	err = fU.forecastRepository.GetAllFiltered(
		ctx,
		location,
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

	if refreshIsNeeded(ctx, fU.refreshRepository) {
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
		utils.ComputeForecasts(&forecasts, openDataForecasts)
		// Delete all forecasts
		_, err := fU.forecastRepository.DeleteAll(ctx)
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

// Check if it's possible to perform a data refresh
func refreshIsNeeded(ctx context.Context, refreshRepository domains.RefreshRepository) bool {

	var lastRefresh domains.Refresh

	// Get the last refresh to be sure this is not too early
	err := refreshRepository.GetLast(ctx, &lastRefresh)

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
