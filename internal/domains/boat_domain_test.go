package domains

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBoatIsEqualOK(t *testing.T) {
	crossingDateApproximation, _ := time.Parse(time.RFC3339, "2023-02-26T21:00:00+01:00")
	boat := Boat{
		Name:                      "MY_BOAT",
		IsLeavingDock:             false,
		ApproximativeCrossingDate: crossingDateApproximation,
	}
	otherBoat := Boat{
		Name:                      "MY_BOAT",
		IsLeavingDock:             false,
		ApproximativeCrossingDate: crossingDateApproximation,
	}
	result := boat.IsEqual(otherBoat)

	assert.True(t, result)
}

func TestBoatIsEqualNOK(t *testing.T) {
	crossingDateApproximation, _ := time.Parse(time.RFC3339, "2023-02-26T21:00:00+01:00")
	boat := Boat{
		Name:                      "MY_BOAT",
		IsLeavingDock:             false,
		ApproximativeCrossingDate: crossingDateApproximation,
	}
	otherBoat := Boat{
		Name:                      "MY_BOAT",
		IsLeavingDock:             false,
		ApproximativeCrossingDate: crossingDateApproximation,
	}
	result := boat.IsEqual(otherBoat)

	assert.False(t, result)
}

func TestBoatsAreEqualOK(t *testing.T) {
	crossingDateApproximation, _ := time.Parse(time.RFC3339, "2023-02-26T21:00:00+01:00")
	boats := Boats{
		Boat{
			Name:                      "MY_BOAT",
			IsLeavingDock:             false,
			ApproximativeCrossingDate: crossingDateApproximation,
		},
	}
	otherBoats := Boats{
		Boat{
			Name:                      "MY_BOAT",
			IsLeavingDock:             false,
			ApproximativeCrossingDate: crossingDateApproximation,
		},
	}
	result := boats.AreEqual(otherBoats)

	assert.True(t, result)
}

func TestBoatsAreEqualNOK(t *testing.T) {
	crossingDateApproximation, _ := time.Parse(time.RFC3339, "2023-02-26T21:00:00+01:00")
	boats := Boats{
		Boat{
			Name:                      "MY_BOAT",
			IsLeavingDock:             false,
			ApproximativeCrossingDate: crossingDateApproximation,
		},
	}
	otherBoats := Boats{
		Boat{
			Name:                      "MY_BOAT2",
			IsLeavingDock:             false,
			ApproximativeCrossingDate: crossingDateApproximation,
		},
	}
	result := boats.AreEqual(otherBoats)

	assert.False(t, result)
}
