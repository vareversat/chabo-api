package utils

import (
	"strings"
	"time"

	"github.com/vareversat/chabo-api/internal/models"
)

// Return the corresponding models.ClosingType according of the string value
func MapClosingType(stringClosingType string) models.ClosingType {
	if stringClosingType == "oui" {
		return models.TwoWay
	} else {
		return models.OneWay
	}
}

// Return the corresponding models.ClosingReason according of the string value
func MapClosingReason(stringClosingReason string) models.ClosingReason {
	if stringClosingReason == "MAINTENANCE" {
		return models.Maintenance
	} else {
		return models.BoatReason
	}
}

// Return a []models.Boat of a boat crossing forecast.
// closingReason : If it is a maintenance forecast, no computation
// boatNames : The raw string containing the boat name(s)
// closingDuration : Used to compute the approximated crossing time
// circulationClosingDate : Used to compute the approximated crossing time
// alreadySeenBoatNames : Array pointer to keep track of the boats. Used to compute the boat Maneuver
// forecastID : Used to compute teh "self" link
func MapBoats(
	closingReason models.ClosingReason,
	boatNames string,
	closingDuration time.Duration,
	circulationClosingDate time.Time,
	alreadySeenBoatNames *[]string,
	forecastID string,
) []models.Boat {
	var boats []models.Boat
	if closingReason == models.BoatReason {
		// The string may contains multiple boat name separated by a "/"
		boatNamesSlice := strings.Split(boatNames, "/")
		for _, boat := range boatNamesSlice {
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
				Name:                      boatName,
				Maneuver:                  action,
				ApproximativeCrossingDate: circulationClosingDate.Add(closingDuration / 2),
			})
		}
	}
	return boats

}
