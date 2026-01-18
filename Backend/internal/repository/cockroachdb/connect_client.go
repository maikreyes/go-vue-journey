package cockroachdb

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(DSN *string) (*pgxpool.Pool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(*DSN)
	if err != nil {
		return nil, err
	}

	// Valor conservador para manejar concurrencia (HTTP + workers).
	// Ajustable vía DSN (pool_max_conns) si lo necesitas.
	if poolConfig.MaxConns == 0 {
		poolConfig.MaxConns = 20
	}

	db, err := pgxpool.NewWithConfig(ctx, poolConfig)

	if err != nil {
		return nil, err
	}

	return db, nil
}
