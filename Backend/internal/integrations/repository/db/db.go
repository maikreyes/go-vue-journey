package db

import (
	"database/sql"
	"fmt"
)

func Connect(conectString string) (*sql.DB, error) {

	db, err := sql.Open("pgx", conectString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db unreachable: %w", err)
	}

	return db, nil
}
