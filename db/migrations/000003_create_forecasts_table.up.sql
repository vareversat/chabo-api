CREATE TABLE
    IF NOT EXISTS forecasts (
        forecast_id TEXT,
        closing_event_name TEXT,
        closing_duration_min INTEGER NOT NULL,
        circulation_closing_date TIMESTAMPTZ NOT NULL,
        circulation_reopening_date TIMESTAMPTZ NOT NULL,
        is_traffic_fully_closed BOOLEAN NOT NULL,

        PRIMARY KEY(forecast_id)
    );