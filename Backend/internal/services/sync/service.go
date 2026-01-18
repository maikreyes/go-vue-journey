package sync

import "backend/internal/ports"

type Service struct {
	Provider   ports.StockProvider
	Repository ports.StocksRepository
	Workers    int
	BatchSize  int
}

func NewService(provider ports.StockProvider, repository ports.StocksRepository, workers int, batchSize int) *Service {
	return &Service{
		Provider:   provider,
		Repository: repository,
		Workers:    workers,
		BatchSize:  batchSize,
	}
}
