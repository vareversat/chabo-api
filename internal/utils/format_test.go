package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatDateTime(t *testing.T) {
	want, _ := time.Parse(time.RFC3339, "2023-02-26T21:00:00+01:00")
	loc, _ := time.LoadLocation("Europe/Paris")
	if value, err := FormatDataTime("21:00", "2023-02-26", 3600, *loc); err == nil {
		assert.True(t, want.Equal(value))
	} else {
		t.Fatalf(err.Error())
	}
}
