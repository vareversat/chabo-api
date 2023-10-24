package domains

import (
	"testing"
	"time"
)

func TestBoatIsEqualOK(t *testing.T) {
	approximativeCrossingDate, _ := time.Parse(time.RFC3339, "2023-02-26T21:00:00+01:00")
	boat := Boat{
		Name:                      "MY_BOAT",
		Maneuver:                  Entering,
		ApproximativeCrossingDate: approximativeCrossingDate,
	}
	otherBoat := Boat{
		Name:                      "MY_BOAT",
		Maneuver:                  Entering,
		ApproximativeCrossingDate: approximativeCrossingDate,
	}
	result := boat.IsEqual(otherBoat)
	want := true

	if want != result {
		t.Fatalf(`IsEqual("otherBoat") = %v, want match for %v`, result, want)
	}
}

func TestBoatIsEqualNOK(t *testing.T) {
	approximativeCrossingDate, _ := time.Parse(time.RFC3339, "2023-02-26T21:00:00+01:00")
	boat := Boat{
		Name:                      "MY_BOAT",
		Maneuver:                  Entering,
		ApproximativeCrossingDate: approximativeCrossingDate,
	}
	otherBoat := Boat{
		Name:                      "MY_BOAT",
		Maneuver:                  Leaving,
		ApproximativeCrossingDate: approximativeCrossingDate,
	}
	result := boat.IsEqual(otherBoat)
	want := false

	if want != result {
		t.Fatalf(`IsEqual("otherBoat") = %v, want match for %v`, result, want)
	}
}

func TestBoatsAreEqualOK(t *testing.T) {
	approximativeCrossingDate, _ := time.Parse(time.RFC3339, "2023-02-26T21:00:00+01:00")
	boats := Boats{
		Boat{
			Name:                      "MY_BOAT",
			Maneuver:                  Entering,
			ApproximativeCrossingDate: approximativeCrossingDate,
		},
	}
	otherBoats := Boats{
		Boat{
			Name:                      "MY_BOAT",
			Maneuver:                  Entering,
			ApproximativeCrossingDate: approximativeCrossingDate,
		},
	}
	result := boats.AreEqual(otherBoats)
	want := true

	if want != result {
		t.Fatalf(`AreEqual("otherBoats") = %v, want match for %v`, result, want)
	}
}

func TestBoatsAreEqualNOK(t *testing.T) {
	approximativeCrossingDate, _ := time.Parse(time.RFC3339, "2023-02-26T21:00:00+01:00")
	boats := Boats{
		Boat{
			Name:                      "MY_BOAT",
			Maneuver:                  Entering,
			ApproximativeCrossingDate: approximativeCrossingDate,
		},
	}
	otherBoats := Boats{
		Boat{
			Name:                      "MY_BOAT2",
			Maneuver:                  Entering,
			ApproximativeCrossingDate: approximativeCrossingDate,
		},
	}
	result := boats.AreEqual(otherBoats)
	want := false

	if want != result {
		t.Fatalf(`AreEqual("otherBoats") = %v, want match for %v`, result, want)
	}
}
