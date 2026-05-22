package storage

import (
	"fmt"

	"github.com/sahilium/florence/internal/event"
)

func (db *DB) StoreEvent(ev event.Event) error {
	_, err := db.Exec(
		`INSERT OR IGNORE INTO events (id, source, title, body, url, severity) VALUES (?, ?, ?, ?, ?, ?)`,
		ev.ID, ev.Source, ev.Title, ev.Body, ev.URL, ev.Severity,
	)
	if err != nil {
		return fmt.Errorf("store event: %w", err)
	}

	return nil
}
