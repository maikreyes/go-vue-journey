package stocks

import (
	"backend/internal/domain"
)

func (s *Service) GetFilterStocks(page *string, filter *string) (*domain.StocksPage, error) {
	limit := 10

	stocksPage, err := s.Repository.GetFilterStocks(page, limit, filter)

	if err != nil {
		return nil, err
	}

	return stocksPage, nil
}
