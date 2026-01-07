package noop

import "go-vue-journey/internal/stock"

type Repository struct{}

func New() *Repository {
	return &Repository{}
}

func (r *Repository) Upsert(stock.Stock) error {
	return nil
}

func (r *Repository) GetStocks(limit int, cursorTicker *string, filter stock.StockFilter) ([]stock.Stock, error) {
	return []stock.Stock{}, nil
}

func (r *Repository) GetStocksStats() (stock.StocksStats, error) {
	return stock.StocksStats{}, nil
}

func (r *Repository) GetTopStocks(n int) ([]stock.Stock, error) {
	return []stock.Stock{}, nil
}
