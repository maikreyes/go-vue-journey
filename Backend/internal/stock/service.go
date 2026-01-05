package stock

type StockProvider interface {
	GetStocks(*string) (*Page, error)
}

type StockRepository interface {
	Upsert(Stock) error
	GetStocks() ([]Stock, error)
	GetTopStocks(int) ([]Stock, error)
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

func (s *Service) ListStocks(page *string) (*[]Stock, error) {
	apiResult, err := s.provider.GetStocks(page)
	if err != nil {
		return nil, err
	}

	// Store API results in DB
	for _, item := range apiResult.Items {
		s.repository.Upsert(item)
	}

	// Get stocks from DB
	repoResult, err := s.repository.GetStocks()
	if err != nil {
		// If DB fails, return API results as fallback
		return &apiResult.Items, nil
	}

	return &repoResult, nil
}
