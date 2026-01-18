package cockroachdb

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Migrate(db *pgxpool.Pool) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	_, err := db.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS stocks (
		ticker TEXT PRIMARY KEY,
		target_from TEXT,
		target_to TEXT,
		company TEXT,
		action TEXT,
		brokerage TEXT,
		rating_from TEXT,
		rating_to TEXT,
		created_at TIMESTAMPTZ
	);
`)
	if err != nil {
		fmt.Printf("[MIGRATE][ERROR] Migrate Failed: %v\n", err)
		return err
	}

	fmt.Printf("[MIGRATE] Migrate Successfully\n")
	return nil
}
