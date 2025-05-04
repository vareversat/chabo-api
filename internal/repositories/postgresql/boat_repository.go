package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vareversat/chabo-api/internal/domains"
)

var TABLE_NAME = "boats"

type boatRepository struct {
	connectionPool *pgxpool.Pool
}

func NewBoatRepository(connectionPool *pgxpool.Pool) domains.BoatRepository {
	return &boatRepository{
		connectionPool: connectionPool,
	}
}

func (bR *boatRepository) InsertOne(ctx context.Context, boat domains.Boat) (err error, id int) {
	query := `
        INSERT INTO ` + TABLE_NAME + `(name)
        VALUES (@name)
		ON CONFLICT (name)
		DO UPDATE 
		SET name = EXCLUDED.name
        RETURNING boat_id
    `

	args := pgx.NamedArgs{
		"name": boat.Name,
	}
	err = bR.connectionPool.QueryRow(ctx, query, args).Scan(&id)

	return err, id
}

func (bR *boatRepository) GetByForecastId(ctx context.Context, forecastId string, boats *domains.Boats) error {
	query := `
		SELECT b.name,
			fb.is_leaving_dock,
			fb.approximative_crossing_date,
			sbd.imo,
			sbd.mmsi
		FROM ` + TABLE_NAME + ` AS b
		INNER JOIN forecasts_boats AS fb
			ON b.boat_id = fb.boat_id
			AND fb.forecast_id = '@forecast_id'
		LEFT JOIN static_boat_data AS sbd
			ON b.name = sbd.name
	`

	args := pgx.NamedArgs{
		"forecast_id": forecastId,
	}
	rows, err := bR.connectionPool.Query(ctx, query, args)

	if err != nil {
		return fmt.Errorf("Error while querying boat for forecasts %s: %s", forecastId, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var boat domains.Boat
		err := rows.Scan(
			&boat.Name,
			&boat.IsLeavingDock,
			&boat.ApproximativeCrossingDate,
			&boat.IMO,
			&boat.MMSI,
		)
		if err != nil {
			return fmt.Errorf("Error while scanning row on boat: %s", err)
		}
		*boats = append(*boats, boat)
	}

	return nil
}
