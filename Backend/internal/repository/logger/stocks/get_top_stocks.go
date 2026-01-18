package stocks

import (
	"backend/internal/domain"
	"fmt"
	"time"
)

func (r *Repository) GetTopStocks(limit int) (*[]domain.Stock, error) {

	start := time.Now()
	stocks, err := r.Repository.GetTopStocks(limit)
	if err != nil {
		elapsed := time.Since(start)
		fmt.Printf("[LOGGER][GET_TOP_STOCKS] Fetched top %d stocks in %s\n", limit, elapsed)
		return nil, err
	}

	elapsed := time.Since(start)
	fmt.Printf("[LOGGER][GET_TOP_STOCKS] Fetched top %d stocks in %s\n", limit, elapsed)
	return stocks, nil
}
