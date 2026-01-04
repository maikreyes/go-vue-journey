package stock

type StockProvider interface {
	GetStocks(*string) (*Page, error)
}

type StockRepository interface {
	Upsert(Stock) error
}

type Service struct {
	provider   StockProvider
	repository StockRepository
}

func NewService(provider StockProvider, repository StockRepository) *Service {
	return &Service{
		provider:   provider,
		repository: repository,
	}
}

func (s *Service) ListStocks(page *string) (*Page, error) {
	result, err := s.provider.GetStocks(page)
	if err != nil {
		return nil, err
	}
	go func(items []Stock) {
		for _, item := range items {
			s.repository.Upsert(item)
		}
	}(result.Items)

	return result, nil
}
