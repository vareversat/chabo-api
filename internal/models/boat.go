package models

import "time"

type BoatManeuver string

const (
	Leaving  BoatManeuver = "leaving_bordeaux"
	Entering BoatManeuver = "entering_in_bordeaux"
)

type Boat struct {
	Name                      string       `json:"name" bson:"name" example:"EUROPA 2"`
	Maneuver                  BoatManeuver `json:"maneuver" bson:"maneuver"`
	ApproximativeCrossingDate time.Time    `json:"approximative_crossing_date" bson:"approximative_crossing_date" example:"2021-05-25T00:53:16.535668Z" format:"date-time"`
}
