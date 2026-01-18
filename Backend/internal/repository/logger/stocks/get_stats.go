package stocks

import (
	"backend/internal/domain"
	"fmt"
	"time"
)

func (r *Repository) GetStats(limit int, filter *string, ticker *string) (*domain.StocksStats, error) {

	startTime := time.Now()

	stats, err := r.Repository.GetStats(limit, filter, ticker)
	if err != nil {
		elapsed := time.Since(startTime)
		fmt.Printf("[LOGGER][GET_STATS] Fetched stocks stats in %s\n", elapsed)
		return nil, err
	}

	elapsed := time.Since(startTime)
	fmt.Printf("[LOGGER][GET_STATS] Fetched stocks stats in %s\n", elapsed)
	return stats, nil

}
