package stocks

import (
	"backend/internal/domain"
	"fmt"
	"time"
)

func (r *Repository) GetStocks(page *string, limit int) (*domain.StocksPage, error) {

	start := time.Now()
	stocksPage, err := r.Repository.GetStocks(page, limit)

	if err != nil {
		elapsed := time.Since(start)
		fmt.Printf("[LOGGER][GET_STOCKS] Fetched stocks page in %s\n", elapsed)
		return nil, err
	}

	elapsed := time.Since(start)
	fmt.Printf("[LOGGER][GET_STOCKS] Fetched stocks page in %s\n", elapsed)

	return stocksPage, nil
}
