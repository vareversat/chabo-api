SELECT b.name, f.circulation_closing_date AT TIME ZONE 'Europe/Paris', f.circulation_reopening_date, fb.is_leaving
FROM boats AS b
INNER JOIN forecasts_boats AS fb ON b.boat_id = fb.boat_id
INNER JOIN forecasts AS f ON f.forecast_id = fb.forecast_id
ORDER BY f.circulation_closing_date;