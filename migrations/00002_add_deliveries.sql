-- +goose Up
CREATE TABLE deliveries (
    event_id TEXT NOT NULL,
    sink TEXT NOT NULL,
    delivered BOOLEAN NOT NULL DEFAULT 0,
    delivered_at DATETIME,
    error TEXT
);

-- +goose Down
DROP TABLE IF EXISTS deliveries;
