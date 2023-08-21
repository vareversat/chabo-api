package utils

import (
	"fmt"
	"strconv"
	"time"
)

// Take the date, the hours and the timezone offset and return the related time.Time object
func FormatDataTime(
	stringTime string,
	closingDate string,
	offsetInSec int,
	location time.Location,
) (time.Time, error) {
	offsetInHours := offsetInSec / 3600
	stringDate := closingDate + "T" + stringTime + ":00+0" + strconv.Itoa(offsetInHours) + ":00"
	time, err := time.Parse(time.RFC3339, stringDate)
	if err != nil {
		return time.UTC(), fmt.Errorf(err.Error())
	}
	return time.In(&location), nil
}
