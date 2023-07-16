package models

import "time"

type ClosingType string
type ClosingReason string

const (
	TwoWay ClosingType = "two_way"
	OneWay ClosingType = "one_way"
)

const (
	BoatCrossing ClosingReason = "boat_crossing"
	Maintenance  ClosingReason = "maintenance"
)

type Forecast struct {
	ID                       string          `json:"id" bson:"_id" example:"63a6430fc07ff1d895c9555ef2ef6e41c1e3b1f5"`
	ClosingType              ClosingType     `json:"closing_type" bson:"closing_type"`
	ClosingDuration          time.Duration   `json:"closing_duration_ns" bson:"closing_duration_ns" swaggertype:"primitive,integer" example:"4980000000000"`
	CirculationClosingDate   time.Time       `json:"circulation_closing_date" bson:"circulation_closing_date" example:"2021-05-25T00:53:16.535668Z" format:"date-time"`
	CirculationReopeningDate time.Time       `json:"circulation_reopening_date" bson:"circulation_reopening_date" example:"2021-05-25T00:53:16.535668Z" format:"date-time"`
	ClosingReason            ClosingReason   `json:"closing_reason" bson:"closing_reason"`
	Boats                    []Boat          `json:"boats,omitempty" bson:"boats,omitempty"`
	Link                     OpenAPISelfLink `json:"_links" bson:"_links" swaggerignore:"true"`
}
