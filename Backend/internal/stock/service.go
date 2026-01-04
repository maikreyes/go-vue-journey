package stock

type StockProvider interface {
	GetStocks(*string) (*Page, error)
}

type Service struct {
	provider StockProvider
}

func NewService(provider StockProvider) *Service {
	return &Service{provider: provider}
}

func (s *Service) ListStocks(page *string) (*Page, error) {
	return s.provider.GetStocks(page)
}
