package utils

import (
	"strings"
	"time"

	"github.com/vareversat/chabo-api/internal/domains"
	"github.com/vareversat/chabo-api/internal/models"
)

// Return the corresponding models.ClosingType according of the string value
func MapClosingType(stringClosingType string) domains.ClosingType {
	if stringClosingType == "oui" {
		return domains.TwoWay
	} else {
		return domains.OneWay
	}
}

// Return the corresponding models.ClosingReason according of the string value
func MapClosingReason(stringClosingReason string) domains.ClosingReason {
	if stringClosingReason == "MAINTENANCE" {
		return domains.Maintenance
	} else {
		return domains.BoatReason
	}
}

// Return a []models.Boat of a boat crossing forecast.
// closingReason : If it is a maintenance forecast, no computation
// boatNames : The raw string containing the boat name(s)
// closingDuration : Used to compute the approximated crossing time
// circulationClosingDate : Used to compute the approximated crossing time
// alreadySeenBoatNames : Array pointer to keep track of the boats. Used to compute the boat Maneuver
// forecastID : Used to compute the "self" link
func MapBoats(
	closingReason domains.ClosingReason,
	boatNames string,
	closingDuration time.Duration,
	circulationClosingDate time.Time,
	alreadySeenBoatNames *[]string,
	forecastID string,
) []models.Boat {
	var boats []models.Boat
	if closingReason == domains.BoatReason {
		// The string may contains multiple boat name separated by a "/"
		boatNamesSlice := strings.Split(boatNames, "/")
		for index, boat := range boatNamesSlice {
			boatName := strings.TrimSpace(boat)
			var action models.BoatManeuver
			if contains(*alreadySeenBoatNames, boatName) {
				// If the boat is already in the list, that means that it is docked in Bordeaux
				*alreadySeenBoatNames = remove(*alreadySeenBoatNames, boatName)
				action = models.Leaving
			} else {
				// If not, that means that it is entering in Bordeaux
				action = models.Entering
				*alreadySeenBoatNames = append(*alreadySeenBoatNames, boatName)
			}
			boats = append(boats, models.Boat{
				Name:     boatName,
				Maneuver: action,
				ApproximativeCrossingDate: computeApproximativeCrossingDate(
					circulationClosingDate,
					closingDuration,
					len(boatNamesSlice),
					index,
				),
			})
		}
	}
	return boats

}

// Return a time.Time representing the boat may cross the Chaban bridge.
// circulationClosingDate : The moment the bridge will close
// closingDuration : The duration of the closing
// boatCount : How many boats will cross the bridge
// boatIndex : The place of the boat
func computeApproximativeCrossingDate(
	circulationClosingDate time.Time,
	closingDuration time.Duration,
	boatCount int,
	boatIndex int,
) time.Time {
	// Get the fraction by how much the duration will be split.
	// Example : 1 boat => 1/2 | 2 boats => 1st = 1/3 & 2nd = 2/3
	durationFraction := (float64(boatIndex+1) / float64(boatCount+1))
	return circulationClosingDate.Add(time.Duration(float64(closingDuration) * durationFraction))
}
