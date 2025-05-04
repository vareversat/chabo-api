package utils

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vareversat/chabo-api/internal/domains"
)

func TestMapClosingReason_MAINTENANCE(t *testing.T) {
	expected := domains.Maintenance
	value := MapClosingReason("MAINTENANCE")

	assert.Equal(t, expected, value)
}

func TestMapClosingReason_BOAT(t *testing.T) {
	expected := domains.BoatPassage
	value := MapClosingReason("SILVER DAWN")

	assert.Equal(t, expected, value)
}

func TestMapClosingReason_WINE_FESTIVAL(t *testing.T) {
	expected := domains.WineFestivalBoats
	value := MapClosingReason("Bateaux fete du vin")

	assert.Equal(t, expected, value)
}

func TestMapClosingReason_SPECIAL_EVENT(t *testing.T) {
	expected := domains.SpecialEvent
	value := MapClosingReason("Les étoiles filantes")

	assert.Equal(t, expected, value)
}

func TestMapBoats(t *testing.T) {
	var alreadySeenBoatNames []string
	var duration time.Duration
	localTime := time.Now()
	duration = 10000000000
	crossingTime := localTime.Add(duration / 2)

	expected := []domains.Boat{
		{Name: "MY_BOAT", IsLeavingDock: false, ApproximativeCrossingDate: crossingTime},
	}
	value := MapBoats(domains.BoatPassage, "MY_BOAT", duration, localTime, &alreadySeenBoatNames)

	assert.True(t, reflect.DeepEqual(expected, value))
}

func TestMapBoatsMultiBoats(t *testing.T) {
	var alreadySeenBoatNames []string
	var duration time.Duration
	localTime := time.Now()
	duration = 10000000000
	crossingTimeBoat1 := localTime.Add(duration / 3)
	crossingTimeBoat2 := localTime.Add(time.Duration(float64(duration) * (float64(2) / float64(3))))

	expected := []domains.Boat{
		{Name: "MY_BOAT", IsLeavingDock: false, ApproximativeCrossingDate: crossingTimeBoat1},
		{
			Name:                      "MY_SECOND_BOAT",
			IsLeavingDock:             false,
			ApproximativeCrossingDate: crossingTimeBoat2,
		},
	}
	value := MapBoats(
		domains.BoatPassage,
		"MY_BOAT /MY_SECOND_BOAT",
		duration,
		localTime,
		&alreadySeenBoatNames,
	)

	assert.True(t, reflect.DeepEqual(expected, value))
}

func TestMapBoatsExistingBoats(t *testing.T) {
	var alreadySeenBoatNames []string
	alreadySeenBoatNames = append(alreadySeenBoatNames, "MY_BOAT")
	var duration time.Duration
	localTime := time.Now()
	duration = 10000000000
	crossingTime := localTime.Add(duration / 2)

	expected := []domains.Boat{
		{Name: "MY_BOAT", IsLeavingDock: true, ApproximativeCrossingDate: crossingTime},
	}
	value := MapBoats(domains.BoatPassage, "MY_BOAT", duration, localTime, &alreadySeenBoatNames)

	assert.True(t, reflect.DeepEqual(expected, value))
}

func TestComputeCrossingDateApproximationBoat(t *testing.T) {
	var closingDuration time.Duration
	localTime := time.Now()
	closingDuration = 10000000000

	expected := localTime.Add(closingDuration / 2)

	value := computeCrossingDateApproximation(
		localTime,
		closingDuration,
		1, 0)

	assert.True(t, reflect.DeepEqual(expected, value))
}

func TestComputeCrossingDateApproximationTwoBoat_First(t *testing.T) {
	var closingDuration time.Duration
	localTime := time.Now()
	closingDuration = 10000000000

	expected := localTime.Add(closingDuration / 3)

	value := computeCrossingDateApproximation(
		localTime,
		closingDuration,
		2, 0)

	assert.True(t, reflect.DeepEqual(expected, value))
}

func TestComputeCrossingDateApproximationTwoBoat_Second(t *testing.T) {
	var closingDuration time.Duration
	localTime := time.Now()
	closingDuration = 10000000000

	expected := localTime.Add(time.Duration(float64(closingDuration) * (float64(2) / float64(3))))

	value := computeCrossingDateApproximation(
		localTime,
		closingDuration,
		2, 1)

	assert.True(t, reflect.DeepEqual(expected, value))
}
