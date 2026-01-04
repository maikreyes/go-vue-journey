package sync

import (
	"go-vue-journey/internal/stock"
	"sync"
)

type Service struct {
	provider   stock.StockProvider
	repository stock.StockRepository
	workers    int
}

func NewService(provider stock.StockProvider, repository stock.StockRepository, workers int) *Service {
	return &Service{
		provider:   provider,
		repository: repository,
		workers:    workers,
	}
}

func (s *Service) Run() error {
	stocksCh := make(chan stock.Stock)
	var wg sync.WaitGroup

	// workers
	for i := 0; i < s.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for stock := range stocksCh {
				_ = s.repository.Upsert(stock)
			}
		}()
	}

	seenPages := make(map[string]struct{})

	var page *string

	for {
		result, err := s.provider.GetStocks(page)
		if err != nil {
			return err
		}

		// Guardar stocks
		for _, item := range result.Items {
			stocksCh <- item
		}

		// Si no hay siguiente → fin
		if result.NextPage == nil {
			break
		}

		// Detectar loop
		if _, exists := seenPages[*result.NextPage]; exists {
			break
		}

		seenPages[*result.NextPage] = struct{}{}
		page = result.NextPage
	}

	close(stocksCh)
	wg.Wait()
	return nil
}
