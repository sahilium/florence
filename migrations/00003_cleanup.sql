-- +goose Up
DROP TABLE IF EXISTS seen_events;

-- +goose Down
CREATE TABLE seen_events (
    event_id TEXT PRIMARY KEY,
    source TEXT NOT NULL,
    seen_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
