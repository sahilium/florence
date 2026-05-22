package storage

import (
	"database/sql"
	"fmt"
)

func (db *DB) GetSourceState(source string) (string, error) {
	var cursor sql.NullString
	err := db.QueryRow(
		"SELECT last_cursor FROM source_state WHERE source = ?",
		source,
	).Scan(&cursor)
	if err != nil {
		return "", fmt.Errorf("get source state: %w", err)
	}

	return cursor.String, nil
}

func (db *DB) SetSourceState(source, cursor string) error {
	_, err := db.Exec(
		`INSERT INTO source_state (source, last_cursor) VALUES (?, ?)
		 ON CONFLICT(source) DO UPDATE SET last_cursor = excluded.last_cursor,
		                                   updated_at = CURRENT_TIMESTAMP`,
		source, cursor,
	)
	if err != nil {
		return fmt.Errorf("set source state: %w", err)
	}

	return nil
}
