package domains

import (
	"context"
	"time"

	"github.com/vareversat/chabo-api/internal/errors"
)

type BoatManeuver string

const (
	Leaving  BoatManeuver = "leaving_bordeaux"
	Entering BoatManeuver = "entering_in_bordeaux"
)

type Boat struct {
	Name                      string    `json:"name"                        bson:"name"                        example:"EUROPA 2"`
	IsLeavingDock             bool      `json:"is_leaving_dock"             bson:"is_leaving_dock"`
	ApproximativeCrossingDate time.Time `json:"approximative_crossing_date" bson:"approximative_crossing_date" example:"2021-05-25T00:53:16.535668Z" format:"date-time"`
	IMO                       int64     `json:"imo"                         bson:"imo"                         example:"9616230"`
	MMSI                      int64     `json:"mmsi"                        bson:"mmsi"                        example:"229378000"`
}

type Boats []Boat

func (b Boat) IsEqual(other Boat) bool {
	return b.Name == other.Name &&
		b.IsLeavingDock == other.IsLeavingDock &&
		b.ApproximativeCrossingDate.Equal(other.ApproximativeCrossingDate)
}

func (boats Boats) AreEqual(other Boats) bool {
	if len(boats) != len(other) {
		return false
	}
	for i, b := range boats {
		if !b.IsEqual(other[i]) {
			return false
		}
	}

	return true
}

type BoatRepository interface {
	InsertOne(ctx context.Context, boat Boat) (error, int)
	GetByForecastId(ctx context.Context, forecastId string, boats *Boats) (err error)
}

type BoatUseCase interface {
	InsertOne(ctx context.Context, boat Boat) errors.CustomError
}
