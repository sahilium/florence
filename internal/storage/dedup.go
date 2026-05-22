package storage

import (
	"database/sql"
	"fmt"
)

func (db *DB) EventExists(eventID string) (bool, error) {
	var exists int
	err := db.QueryRow("SELECT 1 FROM events WHERE id = ?", eventID).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("event exists: %w", err)
	}
	return true, nil
}
