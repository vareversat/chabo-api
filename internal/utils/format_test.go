package utils

import (
	"testing"
	"time"
)

func TestFormatDateTime(t *testing.T) {
	want, _ := time.Parse(time.RFC3339, "2023-02-26T21:00:00+01:00")
	loc, _ := time.LoadLocation("Europe/Paris")
	if value, err := FormatDataTime("21:00", "2023-02-26", 3600, *loc); err == nil {
		if !want.Equal(value) {
			t.Fatalf(`TestFormatDataTime("...") = %q, want match for %#q`, value, want)
		}
	} else {
		t.Fatalf(err.Error())
	}
}
