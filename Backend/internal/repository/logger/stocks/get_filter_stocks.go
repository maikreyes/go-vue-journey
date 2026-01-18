package stocks

import (
	"backend/internal/domain"
	"fmt"
	"time"
)

func (r *Repository) GetFilterStocks(page *string, limit int, filter *string) (*domain.StocksPage, error) {

	start := time.Now()
	stocksPage, err := r.Repository.GetFilterStocks(page, limit, filter)
	if err != nil {
		elapsed := time.Since(start)
		fmt.Printf("[LOGGER][GET_UP_STOCKS] Fetched up stocks page in %s\n", elapsed)
		return nil, err
	}

	elapsed := time.Since(start)
	fmt.Printf("[LOGGER][GET_UP_STOCKS] Fetched up stocks page in %s\n", elapsed)

	return stocksPage, nil

}
