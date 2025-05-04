CREATE TABLE
    IF NOT EXISTS boats (
        boat_id INTEGER GENERATED ALWAYS AS IDENTITY,
        name TEXT UNIQUE NOT NULL,

        PRIMARY KEY(boat_id)
    );
