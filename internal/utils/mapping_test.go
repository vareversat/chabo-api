package utils

import (
	"reflect"
	"testing"
	"time"

	"github.com/vareversat/chabo-api/internal/models"
)

func TestMapClosingType_OUI(t *testing.T) {
	want := models.TwoWay
	value := MapClosingType("oui")
	if want != value {
		t.Fatalf(`MapClosingType("oui") = %q, want match for %#q`, value, want)
	}
}

func TestMapClosingType_NON(t *testing.T) {
	want := models.OneWay
	value := MapClosingType("non")
	if want != value {
		t.Fatalf(`MapClosingType("non") = %q, want match for %#q`, value, want)
	}
}

func TestMapClosingReason_MAINTENANCE(t *testing.T) {
	want := models.Maintenance
	value := MapClosingReason("MAINTENANCE")
	if want != value {
		t.Fatalf(`MapClosingReason("MAINTENANCE") = %q, want match for %#q`, value, want)
	}
}

func TestMapClosingReason_BOAT(t *testing.T) {
	want := models.Maintenance
	value := MapClosingReason("MAINTENANCE")
	if want != value {
		t.Fatalf(`MapClosingReason("MAINTENANCE") = %q, want match for %#q`, value, want)
	}
}

func TestMapBoats(t *testing.T) {
	var alreadySeenBoatNames []string
	var duration time.Duration
	localTime := time.Now()
	duration = 10000000000
	crossingTime := localTime.Add(duration / 2)

	want := []models.Boat{{Name: "MY_BOAT", Maneuver: models.Entering, ApproximativeCrossingDate: crossingTime}}
	value := MapBoats(
		models.BoatReason,
		"MY_BOAT",
		duration,
		localTime,
		&alreadySeenBoatNames,
		"BOAT_ID")

	if !reflect.DeepEqual(want, value) {
		t.Fatalf(`MapBoats("...") = %q, want match for %#q`, value, want)
	}
}

func TestMapBoatsMultiBoats(t *testing.T) {
	var alreadySeenBoatNames []string
	var duration time.Duration
	localTime := time.Now()
	duration = 10000000000
	crossingTime := localTime.Add(duration / 2)

	want := []models.Boat{{Name: "MY_BOAT", Maneuver: models.Entering, ApproximativeCrossingDate: crossingTime}, {Name: "MY_SECOND_BOAT", Maneuver: models.Entering, ApproximativeCrossingDate: crossingTime}}
	value := MapBoats(
		models.BoatReason,
		"MY_BOAT /MY_SECOND_BOAT",
		duration,
		localTime,
		&alreadySeenBoatNames,
		"BOAT_ID")

	if !reflect.DeepEqual(want, value) {
		t.Fatalf(`MapBoats("...") = %q, want match for %#q`, value, want)
	}
}
