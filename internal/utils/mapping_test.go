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
	want := models.BoatReason
	value := MapClosingReason("BOAT")
	if want != value {
		t.Fatalf(`MapClosingReason("BOAT") = %q, want match for %#q`, value, want)
	}
}

func TestMapBoats(t *testing.T) {
	var alreadySeenBoatNames []string
	var duration time.Duration
	localTime := time.Now()
	duration = 10000000000
	crossingTime := localTime.Add(duration / 2)

	want := []models.Boat{
		{Name: "MY_BOAT", Maneuver: models.Entering, ApproximativeCrossingDate: crossingTime},
	}
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
	crossingTimeBoat1 := localTime.Add(duration / 3)
	crossingTimeBoat2 := localTime.Add(time.Duration(float64(duration) * (float64(2) / float64(3))))

	want := []models.Boat{
		{Name: "MY_BOAT", Maneuver: models.Entering, ApproximativeCrossingDate: crossingTimeBoat1},
		{
			Name:                      "MY_SECOND_BOAT",
			Maneuver:                  models.Entering,
			ApproximativeCrossingDate: crossingTimeBoat2,
		},
	}
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

func TestMapBoatsExistingBoats(t *testing.T) {
	var alreadySeenBoatNames []string
	alreadySeenBoatNames = append(alreadySeenBoatNames, "MY_BOAT")
	var duration time.Duration
	localTime := time.Now()
	duration = 10000000000
	crossingTime := localTime.Add(duration / 2)

	want := []models.Boat{
		{Name: "MY_BOAT", Maneuver: models.Leaving, ApproximativeCrossingDate: crossingTime},
	}
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

func TestComputeApproximativeCrossingDateOneBoat(t *testing.T) {
	var closingDuration time.Duration
	localTime := time.Now()
	closingDuration = 10000000000

	want := localTime.Add(closingDuration / 2)

	value := computeApproximativeCrossingDate(
		localTime,
		closingDuration,
		1, 0)

	if !reflect.DeepEqual(want, value) {
		t.Fatalf(`computeApproximativeCrossingDate("...") = %q, want match for %#q`, value, want)
	}
}

func TestComputeApproximativeCrossingDateTwoBoat_First(t *testing.T) {
	var closingDuration time.Duration
	localTime := time.Now()
	closingDuration = 10000000000

	want := localTime.Add(closingDuration / 3)

	value := computeApproximativeCrossingDate(
		localTime,
		closingDuration,
		2, 0)

	if !reflect.DeepEqual(want, value) {
		t.Fatalf(`computeApproximativeCrossingDate("...") = %q, want match for %#q`, value, want)
	}
}

func TestComputeApproximativeCrossingDateTwoBoat_Second(t *testing.T) {
	var closingDuration time.Duration
	localTime := time.Now()
	closingDuration = 10000000000

	want := localTime.Add(time.Duration(float64(closingDuration) * (float64(2) / float64(3))))

	value := computeApproximativeCrossingDate(
		localTime,
		closingDuration,
		2, 1)

	if !reflect.DeepEqual(want, value) {
		t.Fatalf(`computeApproximativeCrossingDate("...") = %q, want match for %#q`, value, want)
	}
}
