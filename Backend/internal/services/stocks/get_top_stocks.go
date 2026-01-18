package stocks

import "backend/internal/domain"

func (s *Service) GetTopStocks() (*[]domain.Stock, error) {

	limit := 5

	stocks, err := s.Repository.GetTopStocks(limit)

	if err != nil {
		return nil, err
	}

	return stocks, nil
}
