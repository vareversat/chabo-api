CREATE TABLE
    IF NOT EXISTS syncs (
        sync_id INTEGER GENERATED ALWAYS AS IDENTITY,
        item_count INTEGER NOT NULL,
        duration INTEGER NOT NULL,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

        PRIMARY KEY(sync_id)
    );