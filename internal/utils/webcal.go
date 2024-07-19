package utils

import (
	"time"

	"github.com/vareversat/chabo-api/internal/domains"
	"github.com/vareversat/gics"
	"github.com/vareversat/gics/components"
	"github.com/vareversat/gics/parameters"
	"github.com/vareversat/gics/properties"
	"github.com/vareversat/gics/types"
)

// ComputeCalendar take all the fetched forecasts and the requested timezone
// Return a gics.Calendar (webcal)
func ComputeCalendar(forecasts domains.Forecasts, timezone string) (gics.Calendar, error) {
	calendarComponents := components.CalendarComponents{}
	for _, forecast := range forecasts {
		calendarComponents = append(calendarComponents, components.NewEventCalendarComponent(
			properties.NewUidProperty(
				forecast.ID,
			),
			properties.NewDateTimeStampProperty(time.Now().UTC()),
			[]components.AlarmCalendarComponent{},
			properties.NewDateTimeStartProperty(
				forecast.CirculationClosingDate,
				types.WithLocalTime,
				parameters.NewTimeZoneIdentifierParam(timezone),
			),
			properties.NewDateTimeEndProperty(
				forecast.CirculationReopeningDate,
				types.WithLocalTime,
				parameters.NewTimeZoneIdentifierParam(timezone),
			),
			properties.NewDescriptionProperty(forecast.GetSummary()),
			properties.NewGeographicPositionProperty(44.858339101606994, -0.551626089048817),
			properties.NewSummaryProperty("Fermeture du pont Chaban"),
			properties.NewLocationProperty("Pont Jacques Chaban Delmas - Bordeaux"),
		))
	}

	return gics.NewCalendar(
		calendarComponents,
		"-//Valentin REVERSAT//https://github.com/vareversat/gics//FR",
		"PUBLISH",
		"2.0",
	)

}
