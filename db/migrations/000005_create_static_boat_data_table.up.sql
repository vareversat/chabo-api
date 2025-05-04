CREATE TABLE
    IF NOT EXISTS static_boat_data (
        boat_id INTEGER GENERATED ALWAYS AS IDENTITY,
        name TEXT UNIQUE NOT NULL,
        imo INTEGER,
        mmsi INTEGER,

        PRIMARY KEY(boat_id)
    );
