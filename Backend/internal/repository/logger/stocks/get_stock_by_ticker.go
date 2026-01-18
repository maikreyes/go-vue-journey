package stocks

import (
	"backend/internal/domain"
	"fmt"
	"time"
)

func (r *Repository) GetStockByTicker(ticker string, limit int, page *string, filter *string) (*domain.StocksPage, error) {

	start := time.Now()

	stocks, err := r.Repository.GetStockByTicker(ticker, limit, page, filter)
	if err != nil {
		elapsed := time.Since(start)
		fmt.Printf("[LOGGER][GET_STOCK_BY_TICKER] Fetched stock by ticker %s in %s\n", ticker, elapsed)
		return nil, err
	}

	elapsed := time.Since(start)
	fmt.Printf("[LOGGER][GET_STOCK_BY_TICKER] Fetched stock by ticker %s in %s\n", ticker, elapsed)

	return stocks, nil
}
