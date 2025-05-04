CREATE TABLE
    IF NOT EXISTS forecasts_boats (
        forecast_id TEXT,
        boat_id INTEGER,
        is_leaving_dock BOOLEAN NOT NULL,
        approximative_crossing_date TIMESTAMPTZ NOT NULL,

        PRIMARY KEY(forecast_id, boat_id),
        CONSTRAINT fk_forecast
            FOREIGN KEY(forecast_id)
                REFERENCES forecasts(forecast_id),
        CONSTRAINT fk_boat
            FOREIGN KEY(boat_id)
                REFERENCES boats(boat_id)
    );
