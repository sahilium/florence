-- +goose Up
CREATE TABLE seen_events (
    event_id TEXT PRIMARY KEY,
    source TEXT NOT NULL,
    seen_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE source_state (
    source TEXT PRIMARY KEY,
    last_cursor TEXT,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE events (
    id TEXT PRIMARY KEY,
    source TEXT NOT NULL,
    title TEXT,
    body TEXT,
    url TEXT,
    severity TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS source_state;
DROP TABLE IF EXISTS seen_events;
