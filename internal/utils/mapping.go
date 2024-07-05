package utils

import (
	"strings"
	"time"

	"github.com/vareversat/chabo-api/internal/domains"
)

// MapClosingType Return the corresponding domains.ClosingType according to the string value
func MapClosingType(stringClosingType string) domains.ClosingType {
	if stringClosingType == "oui" {
		return domains.TwoWay
	} else {
		return domains.OneWay
	}
}

// MapClosingReason Return the corresponding domains.ClosingReason according to the string value
func MapClosingReason(stringClosingReason string) domains.ClosingReason {
	switch {
	case strings.Contains(strings.ToLower(stringClosingReason), "vin"):
		return domains.WineFestivalBoats
	// Special events are returned as lower strings by the API
	case StartWithUpperCase(stringClosingReason) && EndWithLowerCase(stringClosingReason):
		return domains.SpecialEvent
	case stringClosingReason == "MAINTENANCE":
		return domains.Maintenance
	default:
		return domains.BoatReason
	}
}

// GetSpecialEventName Return the special name event (if this is a special event)
func GetSpecialEventName(reason domains.ClosingReason, eventName string) string {
	switch {
	// Return the name of the 'boat' (not an actual boat)
	case reason == domains.SpecialEvent || reason == domains.WineFestivalBoats:
		return eventName
	default:
		return ""
	}
}

// MapBoats Return a []domains.Boat of a boat crossing forecast.
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
) []domains.Boat {
	var boats []domains.Boat
	if closingReason == domains.BoatReason {
		// The string may contain multiple boat name separated by a "/"
		boatNamesSlice := strings.Split(boatNames, "/")
		for index, boat := range boatNamesSlice {
			boatName := strings.TrimSpace(boat)
			var action domains.BoatManeuver
			if contains(*alreadySeenBoatNames, boatName) {
				// If the boat is already in the list, that means that it is docked in Bordeaux
				*alreadySeenBoatNames = remove(*alreadySeenBoatNames, boatName)
				action = domains.Leaving
			} else {
				// If not, that means that it is entering in Bordeaux
				action = domains.Entering
				*alreadySeenBoatNames = append(*alreadySeenBoatNames, boatName)
			}
			boats = append(boats, domains.Boat{
				Name:     boatName,
				Maneuver: action,
				CrossingDateApproximation: computeCrossingDateApproximation(
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
func computeCrossingDateApproximation(
	circulationClosingDate time.Time,
	closingDuration time.Duration,
	boatCount int,
	boatIndex int,
) time.Time {
	// Get the fraction by how much the duration will be split.
	// Example : 1 boat => 1/2 | 2 boats => 1st = 1/3 & 2nd = 2/3
	durationFraction := float64(boatIndex+1) / float64(boatCount+1)
	return circulationClosingDate.Add(time.Duration(float64(closingDuration) * durationFraction))
}
