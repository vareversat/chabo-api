package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vareversat/chabo-api/internal/domains"
)

type forecastsBoatsRepository struct {
	connectionPool *pgxpool.Pool
}

func NewForecastsBoatsRepository(connectionPool *pgxpool.Pool) domains.ForecastsBoatsRepository {
	return &forecastsBoatsRepository{
		connectionPool: connectionPool,
	}
}

func (fbR *forecastsBoatsRepository) InsertOne(ctx context.Context, forecast domains.Forecast, boat domains.Boat, boatId int) (err error, id string) {
	query := `
        INSERT INTO forecasts_boats (forecast_id, boat_id, is_leaving_dock, approximative_crossing_date)
        VALUES (@forecast_id, @boat_id, @is_leaving_dock, @approximative_crossing_date)
        RETURNING forecast_id
    `

	args := pgx.NamedArgs{
		"forecast_id":                 forecast.ID,
		"boat_id":                     boatId,
		"is_leaving_dock":             boat.IsLeavingDock,
		"approximative_crossing_date": boat.ApproximativeCrossingDate,
	}
	err = fbR.connectionPool.QueryRow(ctx, query, args).Scan(&id)

	return err, id
}
