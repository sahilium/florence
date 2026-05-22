package storage

import (
	"fmt"
	"time"
)

func (db *DB) RecordDelivery(eventID, sink string, delivered bool, errMsg string) error {
	var errPtr *string
	if errMsg != "" {
		errPtr = &errMsg
	}

	var deliveredAt *time.Time
	if delivered {
		t := time.Now().UTC()
		deliveredAt = &t
	}

	_, err := db.Exec(
		`INSERT INTO deliveries (event_id, sink, delivered, delivered_at, error) VALUES (?, ?, ?, ?, ?)`,
		eventID, sink, delivered, deliveredAt, errPtr,
	)
	if err != nil {
		return fmt.Errorf("record delivery: %w", err)
	}

	return nil
}
