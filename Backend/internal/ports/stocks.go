package ports

import "backend/internal/domain"

type StockProvider interface {
	FetchStocks(page *string) (*domain.StocksPage, error)
}

type StocksRepository interface {
	Upsert(stocks []domain.Stock) error
	GetStocks(page *string, limit int) (*domain.StocksPage, error)
	GetTopStocks(limit int) (*[]domain.Stock, error)
	GetFilterStocks(page *string, limit int, filter *string) (*domain.StocksPage, error)
	GetStats(limit int, filter *string, ticker *string) (*domain.StocksStats, error)
	GetStockByTicker(ticker string, limit int, page *string, filter *string) (*domain.StocksPage, error)
}
