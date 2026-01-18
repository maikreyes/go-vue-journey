package stocks

import (
	"backend/internal/domain"
)

func (s *Service) GetStocks(page *string) (*domain.StocksPage, error) {

	limit := 10

	stocksPage, err := s.Repository.GetStocks(page, limit)

	if err != nil {
		return nil, err
	}

	return stocksPage, nil

}
