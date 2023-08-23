package models

import (
	"reflect"
	"time"
)

type ClosingType string
type ClosingReason string

const (
	TwoWay ClosingType = "two_way"
	OneWay ClosingType = "one_way"
)

const (
	BoatReason  ClosingReason = "boat"
	Maintenance ClosingReason = "maintenance"
)

type Forecast struct {
	ID                       string          `json:"id"                         bson:"_id"                        example:"63a6430fc07ff1d895c9555ef2ef6e41c1e3b1f5"`
	ClosingType              ClosingType     `json:"closing_type"               bson:"closing_type"`
	ClosingDuration          time.Duration   `json:"closing_duration_ns"        bson:"closing_duration_ns"        example:"4980000000000"                            swaggertype:"primitive,integer"`
	CirculationClosingDate   time.Time       `json:"circulation_closing_date"   bson:"circulation_closing_date"   example:"2021-05-25T00:53:16.535668Z"                                              format:"date-time"`
	CirculationReopeningDate time.Time       `json:"circulation_reopening_date" bson:"circulation_reopening_date" example:"2021-05-25T00:53:16.535668Z"                                              format:"date-time"`
	ClosingReason            ClosingReason   `json:"closing_reason"             bson:"closing_reason"`
	Boats                    Boats           `json:"boats,omitempty"            bson:"boats,omitempty"`
	Link                     OpenAPISelfLink `json:"_links"                     bson:"_links"                                                                                                                           swaggerignore:"true"`
}

type Forecasts []Forecast

func (f Forecast) IsEqual(other Forecast) bool {
	return f.ID == other.ID &&
		f.ClosingType == other.ClosingType &&
		f.ClosingDuration == other.ClosingDuration &&
		f.CirculationClosingDate.Equal(other.CirculationClosingDate) &&
		f.CirculationReopeningDate.Equal(other.CirculationReopeningDate) &&
		f.ClosingReason == other.ClosingReason &&
		f.Boats.AreEqual(other.Boats) &&
		reflect.DeepEqual(f.Link, other.Link)
}

func (forecats Forecasts) AreEqual(other Forecasts) bool {
	if len(forecats) != len(other) {
		return false
	}
	for i, b := range forecats {
		if !b.IsEqual(other[i]) {
			return false
		}
	}
	return true
}
