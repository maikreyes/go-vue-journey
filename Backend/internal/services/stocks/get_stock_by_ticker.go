package stocks

import "backend/internal/domain"

func (s *Service) GetStockByTicker(ticker string, page *string, filter *string) (*domain.StocksPage, error) {

	limit := 10

	stocks, err := s.Repository.GetStockByTicker(ticker, limit, page, filter)
	if err != nil {
		return nil, err
	}

	return stocks, nil
}
