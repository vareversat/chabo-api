package domains

import (
	"context"
	"os"
	"reflect"
	"time"

	"github.com/vareversat/chabo-api/internal/errors"
)

var (
	ForecastCollection = os.Getenv("MONGO_FORECASTS_COLLECTION_NAME")
)

type ClosingType string
type ClosingReason string

const (
	TwoWay ClosingType = "two_way"
	OneWay ClosingType = "one_way"
)

const (
	BoatReason        ClosingReason = "boat"
	Maintenance       ClosingReason = "maintenance"
	WineFestivalBoats ClosingReason = "wine_festival_boats"
	SpecialEvent      ClosingReason = "special_event"
)

type Forecasts []Forecast

type Forecast struct {
	ID                       string              `json:"id"                         bson:"_id"                          example:"63a6430fc07ff1d895c9555ef2ef6e41c1e3b1f5"`
	ClosingType              ClosingType         `json:"closing_type"               bson:"closing_type"`
	ClosingDuration          time.Duration       `json:"closing_duration_min"       bson:"closing_duration_min"         example:"83"                                       swaggertype:"primitive,integer"`
	CirculationClosingDate   time.Time           `json:"circulation_closing_date"   bson:"circulation_closing_date"     example:"2021-05-25T00:53:16.535668Z"                                              format:"date-time"`
	CirculationReopeningDate time.Time           `json:"circulation_reopening_date" bson:"circulation_reopening_date"   example:"2021-05-25T00:53:16.535668Z"                                              format:"date-time"`
	ClosingReason            ClosingReason       `json:"closing_reason"             bson:"closing_reason"`
	Boats                    Boats               `json:"boats,omitempty"            bson:"boats,omitempty"`
	Link                     APIResponseSelfLink `json:"_links"                     bson:"_links,omitempty"                                                                                                                   swaggerignore:"true"`
	SpecialEventName         string              `json:"special_event_name"         bson:"special_event_name,omitempty" example:"Les étoiles filantes"                                                                        swaggerignore:"true"`
}

type ForecastsResponse struct {
	Hits      int           `json:"hits"`
	Limit     int           `json:"limit"`
	Offset    int           `json:"offset"`
	Timezone  string        `json:"timezone"         example:"UTC"`
	Links     []interface{} `json:"_links,omitempty"`
	Forecasts Forecasts     `json:"forecasts"`
}

type ForecastResponse struct {
	Timezone string   `json:"timezone" example:"UTC"`
	Forecast Forecast `json:"forecast"`
}

type ForecastMongoResponse struct {
	Results Forecasts                    `json:"results" bson:"results"`
	Count   []ForecastMongoCountResponse `json:"count"   bson:"count"`
}

type ForecastMongoCountResponse struct {
	ItemCount int `json:"itemCount" bson:"itemCount"`
}

func (f *Forecast) IsEqual(other Forecast) bool {
	return f.ID == other.ID &&
		f.ClosingType == other.ClosingType &&
		f.ClosingDuration == other.ClosingDuration &&
		f.CirculationClosingDate.Equal(other.CirculationClosingDate) &&
		f.CirculationReopeningDate.Equal(other.CirculationReopeningDate) &&
		f.ClosingReason == other.ClosingReason &&
		f.Boats.AreEqual(other.Boats) &&
		reflect.DeepEqual(f.Link, other.Link)
}

func (forecasts *Forecasts) AreEqual(other Forecasts) bool {
	if len(*forecasts) != len(other) {
		return false
	}
	for i, b := range *forecasts {
		if !b.IsEqual(other[i]) {
			return false
		}
	}

	return true
}

func (f *Forecast) ChangeLocation(location *time.Location) {
	f.CirculationClosingDate = f.CirculationClosingDate.In(location)
	f.CirculationReopeningDate = f.CirculationReopeningDate.In(location)
	for index, boat := range f.Boats {
		f.Boats[index].CrossingDateApproximation = boat.CrossingDateApproximation.In(
			location,
		)
	}
}

func (forecasts *Forecasts) ChangeLocations(location *time.Location) {

	for index := range *forecasts {
		(*forecasts)[index].ChangeLocation(location)
	}
}

type ForecastRepository interface {
	GetByID(ctx context.Context, id string, forecast *Forecast) error
	GetAllBetweenTwoDates(
		ctx context.Context,
		offset int,
		limit int,
		from time.Time,
		to time.Time,
		forecasts *Forecasts,
		totalItemCount *int,
	) error
	GetAllFiltered(
		ctx context.Context,
		offset int,
		limit int,
		from time.Time,
		reason string,
		maneuver string,
		boat string,
		forecasts *Forecasts,
		totalItemCount *int,
	) error
	DeleteAll(ctx context.Context) (int64, error)
	InsertAll(ctx context.Context, forecasts Forecasts) (int, error)
	GetNextForecast(
		ctx context.Context,
		forecast *Forecast,
		now time.Time,
	) error
	GetCurrentForecast(
		ctx context.Context,
		forecast *Forecast,
	) error
}

type ForecastUseCase interface {
	GetByID(
		ctx context.Context,
		id string,
		forecast *Forecast,
		location *time.Location,
	) errors.CustomError
	GetCurrentForecast(
		ctx context.Context,
		forecast *Forecast,
		location *time.Location,
	) errors.CustomError
	GetNextForecast(
		ctx context.Context,
		forecast *Forecast,
		location *time.Location,
	) errors.CustomError
	GetTodayForecasts(
		ctx context.Context,
		forecasts *Forecasts,
		offset int,
		limit int,
		location *time.Location,
		totalItemCount *int,
	) errors.CustomError
	GetAllFiltered(
		ctx context.Context,
		location *time.Location,
		offset int,
		limit int,
		from time.Time,
		reason string,
		maneuver string,
		boat string,
		forecasts *Forecasts,
		totalItemCount *int,
	) errors.CustomError
	TryToSyncAll(ctx context.Context) (Sync, errors.CustomError)
	ComputeBordeauxAPIResponse(
		forecasts *Forecasts,
		boredeauxAPIResponse BordeauxAPIResponse,
	) errors.CustomError
	SyncIsNeeded(ctx context.Context) bool
}
