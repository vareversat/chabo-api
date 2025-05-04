package domains

import (
	"context"
	"time"

	"github.com/vareversat/chabo-api/internal/errors"
)

type ForecastsBoats struct {
	ForecastID                string
	BoatID                    int
	IsLeavingDock             bool
	ApproximativeCrossingDate time.Time
}

type ForecastsBoatsRepository interface {
	InsertOne(ctx context.Context, forecast Forecast, boat Boat, boatId int) (err error, id string)
}

type ForecastsUseCase interface {
	InsertOne(ctx context.Context, fB ForecastsBoats) (err errors.CustomError)
}
