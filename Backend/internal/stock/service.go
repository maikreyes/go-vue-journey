package stock

type StockProvider interface {
	GetStocks(*string) (*Page, error)
}

type StockRepository interface {
	Upsert(Stock) error
	GetStocks(limit int, cursorTicker *string, filter StockFilter) ([]Stock, error)
	GetStocksStats() (StocksStats, error)
	GetTopStocks(int) ([]Stock, error)
	GetStocksByTicker(tickerPrefix string, limit int, cursorTicker *string) ([]Stock, error)
	GetStocksStatsByTicker(tickerPrefix string) (StocksStats, error)
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

func (s *Service) ListStocks(page *string, limit int, cursorTicker *string, filter StockFilter) (*[]Stock, error) {
	repoResult, err := s.repository.GetStocks(limit, cursorTicker, filter)
	if err == nil {
		return &repoResult, nil
	}

	apiResult, apiErr := s.provider.GetStocks(page)
	if apiErr != nil {
		return nil, apiErr
	}

	return &apiResult.Items, nil
}

func (s *Service) ListStocksWithMeta(page *string, limit int, cursorTicker *string, filter StockFilter) (*StocksResponse, error) {
	itemsPtr, err := s.ListStocks(page, limit, cursorTicker, filter)
	if err != nil {
		return nil, err
	}

	stats, statsErr := s.repository.GetStocksStats()
	if statsErr != nil {
		return nil, statsErr
	}

	viewTotal := stats.Total
	switch filter {
	case StockFilterUp:
		viewTotal = stats.Up
	case StockFilterDown:
		viewTotal = stats.Down
	}

	totalPages := 0
	if limit > 0 {
		totalPages = (viewTotal + limit - 1) / limit
	}

	var nextCursor *string
	if len(*itemsPtr) > 0 {
		last := (*itemsPtr)[len(*itemsPtr)-1].Ticker
		nextCursor = &last
	}

	return &StocksResponse{
		Items:      *itemsPtr,
		Stats:      stats,
		TotalPages: totalPages,
		NextCursor: nextCursor,
	}, nil
}

func (s *Service) ListStocksByTicker(ticker string, limit int, cursorTicker *string) (*StocksResponse, error) {
	items, err := s.repository.GetStocksByTicker(ticker, limit, cursorTicker)
	if err != nil {
		return nil, err
	}

	if limit > 0 && len(items) > limit {
		items = items[:limit]
	}

	stats, err := s.repository.GetStocksStatsByTicker(ticker)
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if limit > 0 {
		totalPages = (stats.Total + limit - 1) / limit
	}

	var nextCursor *string
	if len(items) > 0 {
		last := items[len(items)-1].Ticker
		nextCursor = &last
	}

	return &StocksResponse{
		Items:      items,
		Stats:      stats,
		TotalPages: totalPages,
		NextCursor: nextCursor,
	}, nil
}
