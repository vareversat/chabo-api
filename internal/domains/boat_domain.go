package domains

import (
	"time"
)

type BoatManeuver string

const (
	Leaving  BoatManeuver = "leaving_bordeaux"
	Entering BoatManeuver = "entering_in_bordeaux"
)

type Boat struct {
	Name                      string       `json:"name"                        bson:"name"                        example:"EUROPA 2"`
	Maneuver                  BoatManeuver `json:"maneuver"                    bson:"maneuver"`
	CrossingDateApproximation time.Time    `json:"crossing_date_approximation" bson:"crossing_date_approximation" example:"2021-05-25T00:53:16.535668Z" format:"date-time"`
}

type Boats []Boat

func (b Boat) IsEqual(other Boat) bool {
	return b.Name == other.Name &&
		b.Maneuver == other.Maneuver &&
		b.CrossingDateApproximation.Equal(other.CrossingDateApproximation)
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
