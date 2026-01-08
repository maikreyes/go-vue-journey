package db

import (
	"database/sql"
	"log"
)

func Connect(conectString string) (*sql.DB, error) {

	db, err := sql.Open("pgx", conectString)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("db unreachable:", err)
	}

	return db, nil
}
