package stocks

import "backend/internal/domain"

func (s *Service) GetStats(filter *string, ticker *string) (*domain.StocksStats, error) {

	limit := 10
	stats, err := s.Repository.GetStats(limit, filter, ticker)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
