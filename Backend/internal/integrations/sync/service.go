package sync

import (
	"go-vue-journey/internal/stock"
	"sync"
)

type Service struct {
	provider   stock.StockProvider
	repository stock.StockRepository
	workers    int
	batchSize  int
}

func NewService(provider stock.StockProvider, repository stock.StockRepository, workers int, batchSize int) *Service {
	if batchSize <= 0 {
		batchSize = 10
	}
	return &Service{
		provider:   provider,
		repository: repository,
		workers:    workers,
		batchSize:  batchSize,
	}
}

func (s *Service) Run() error {
	stocksCh := make(chan []stock.Stock)
	var wg sync.WaitGroup

	batchSize := s.batchSize

	// workers
	for i := 0; i < s.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for batch := range stocksCh {
				if batchRepo, ok := s.repository.(interface{ UpsertMany([]stock.Stock) error }); ok {
					_ = batchRepo.UpsertMany(batch)
					continue
				}
				for _, item := range batch {
					_ = s.repository.Upsert(item)
				}
			}
		}()
	}

	seenPages := make(map[string]struct{})

	var page *string
	buffer := make([]stock.Stock, 0, batchSize)

	for {
		result, err := s.provider.GetStocks(page)
		if err != nil {
			return err
		}

		for _, item := range result.Items {
			buffer = append(buffer, item)
			if len(buffer) == batchSize {
				batch := make([]stock.Stock, batchSize)
				copy(batch, buffer)
				stocksCh <- batch
				buffer = buffer[:0]
			}
		}

		if result.NextPage == nil {
			break
		}

		if _, exists := seenPages[*result.NextPage]; exists {
			break
		}

		seenPages[*result.NextPage] = struct{}{}
		page = result.NextPage
	}

	if len(buffer) > 0 {
		batch := make([]stock.Stock, len(buffer))
		copy(batch, buffer)
		stocksCh <- batch
	}

	close(stocksCh)
	wg.Wait()
	return nil
}
