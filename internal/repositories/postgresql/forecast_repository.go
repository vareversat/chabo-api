package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vareversat/chabo-api/internal/domains"
)

type forecastRepository struct {
	connectionPool           *pgxpool.Pool
	boatRepository           domains.BoatRepository
	forecastsBoatsRepository domains.ForecastsBoatsRepository
}

func NewForecastRepository(connectionPool *pgxpool.Pool) domains.ForecastRepository {
	return &forecastRepository{
		connectionPool:           connectionPool,
		boatRepository:           NewBoatRepository(connectionPool),
		forecastsBoatsRepository: NewForecastsBoatsRepository(connectionPool),
	}
}

// InsertAll implements domains.ForecastRepository.
func (rR *forecastRepository) InsertAll(ctx context.Context, forecasts domains.Forecasts) (int, error) {
	forecastQuery := `
        INSERT INTO forecasts (forecast_id, 
				closing_event_name, 
				closing_duration_min, 
				circulation_closing_date, 
				circulation_reopening_date, 
				is_traffic_fully_closed)
        VALUES (@forecast_id, 
			@closing_event_name, 
			@closing_duration_min, 
			@circulation_closing_date, 
			@circulation_reopening_date, 
			@is_traffic_fully_closed)
        RETURNING forecast_id
    `
	for _, forecast := range forecasts {
		var forecastId string
		args := pgx.NamedArgs{
			"forecast_id":                forecast.ID,
			"closing_event_name":         forecast.ClosingReason,
			"closing_duration_min":       forecast.ClosingDuration,
			"circulation_closing_date":   forecast.CirculationClosingDate,
			"circulation_reopening_date": forecast.CirculationReopeningDate,
			"is_traffic_fully_closed":    forecast.IsTrafficFullyClosed,
		}
		forecastErr := rR.connectionPool.QueryRow(ctx, forecastQuery, args).Scan(&forecastId)
		if forecastErr != nil {
			return 0, fmt.Errorf("Error while inserting %s into forecast: %s\n", forecast.ID, forecastErr.Error())
		}
		for _, boat := range forecast.Boats {

			boatErr, boatId := rR.boatRepository.InsertOne(ctx, boat)
			if boatErr != nil {
				return 0, fmt.Errorf("Error while inserting %s into boat: %s\n", boat.Name, boatErr.Error())
			}
			forecastsBoatErr, forecastsBoatId := rR.forecastsBoatsRepository.InsertOne(ctx, forecast, boat, boatId)
			if boatErr != nil {
				return 0, fmt.Errorf("Error while inserting %s into forecasts_boats: %s\n", forecastsBoatId, forecastsBoatErr.Error())
			}
		}
	}

	return len(forecasts), nil
}

// GetAllFiltered implements domains.ForecastRepository.
func (rR *forecastRepository) GetAllFiltered(ctx context.Context, offset int, limit int, from time.Time, reason string, maneuver string, boat string, forecasts *domains.Forecasts, totalItemCount *int) error {
	panic("unimplemented")
}

// DeleteAll implements domains.ForecastRepository.
func (rR *forecastRepository) DeleteAll(ctx context.Context) (int64, error) {
	query := `
		TRUNCATE TABLE forecasts_boats RESTART IDENTITY;
        TRUNCATE TABLE boats RESTART IDENTITY CASCADE;
        TRUNCATE TABLE forecasts CASCADE;
    `
	tag, err := rR.connectionPool.Exec(ctx, query)
	return tag.RowsAffected(), err
}

// GetAllBetweenTwoDates implements domains.ForecastRepository.
func (rR *forecastRepository) GetAllBetweenTwoDates(ctx context.Context, offset int, limit int, from time.Time, to time.Time, forecasts *domains.Forecasts, totalItemCount *int) error {
	panic("unimplemented")
}

// GetByID implements domains.ForecastRepository.
func (rR *forecastRepository) GetByID(ctx context.Context, id string, forecast *domains.Forecast) error {
	var closingEventName string

	query := `
		SELECT * FROM forecasts
		WHERE forecast_id = @forecast_id
	`
	args := pgx.NamedArgs{
		"forecast_id": id,
	}
	err := rR.connectionPool.QueryRow(ctx, query, args).Scan(
		&forecast.ID,
		&closingEventName,
		&forecast.ClosingDuration,
		&forecast.CirculationClosingDate,
		&forecast.CirculationReopeningDate,
		&forecast.IsTrafficFullyClosed)

	if err != nil {
		return fmt.Errorf("Error while scanning forecast with ID %s: %s", id, err.Error())
	}
	forecast.ClosingReason = domains.ClosingReason(closingEventName)

	if forecast.ClosingReason == domains.BoatPassage || forecast.ClosingReason == domains.WineFestivalBoats {
		var boats domains.Boats
		err := rR.boatRepository.GetByForecastId(ctx, id, &boats)
		if err != nil {
			return fmt.Errorf("Error while getting boats for forecast with ID %s: %s", id, err.Error())
		}
		forecast.Boats = boats
	}

	return nil
}

// GetCurrentForecast implements domains.ForecastRepository.
func (rR *forecastRepository) GetCurrentForecast(ctx context.Context, forecast *domains.Forecast) error {
	panic("unimplemented")
}

// GetNextForecast implements domains.ForecastRepository.
func (rR *forecastRepository) GetNextForecast(ctx context.Context, forecast *domains.Forecast, now time.Time) error {
	panic("unimplemented")
}
