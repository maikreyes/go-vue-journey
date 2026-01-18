package stocks

import (
	"backend/internal/domain"
	"fmt"
	"time"
)

func (r *Repository) Upsert(stocks []domain.Stock) error {

	start := time.Now()

	err := r.Repository.Upsert(stocks)

	if err != nil {

		elapsed := time.Since(start)

		fmt.Printf("[LOGGER][UPSERT] Upserted %d stocks in %s\n", len(stocks), elapsed)

		return err
	}

	elapsed := time.Since(start)

	fmt.Printf("[LOGGER][UPSERT] Upserted %d stocks in %s\n", len(stocks), elapsed)

	return nil
}
