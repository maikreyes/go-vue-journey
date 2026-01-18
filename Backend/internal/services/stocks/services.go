package stocks

import "backend/internal/ports"

type Service struct {
	Provider   ports.StockProvider
	Repository ports.StocksRepository
}

func NewService(provider ports.StockProvider, repository ports.StocksRepository) *Service {
	return &Service{
		Provider:   provider,
		Repository: repository,
	}
}
