package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type DB struct {
	*sql.DB
}

func Open(databaseURL, authToken string) (*DB, error) {
	dsn := databaseURL
	if authToken != "" {
		dsn = fmt.Sprintf("%s?authToken=%s", databaseURL, authToken)
	}

	db, err := sql.Open("libsql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	return &DB{db}, nil
}
