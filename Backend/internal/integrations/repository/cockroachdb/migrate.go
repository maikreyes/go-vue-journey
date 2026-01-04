package cockroachdb

import (
	"database/sql"
)

func Migrate(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS stocks (
		ticker TEXT PRIMARY KEY,
		target_from TEXT,
		target_to TEXT,
		company TEXT,
		action TEXT,
		brokerage TEXT,
		rating_from TEXT,
		rating_to TEXT,
		time TIMESTAMPTZ
	)
`)
	if err != nil {
		return err
	}
	return nil
}
