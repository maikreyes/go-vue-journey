package sync

import (
	"backend/internal/domain"
	"sync"
)

func (s *Service) Run() error {
	batchSize := s.BatchSize
	workers := s.Workers

	batchesCh := make(chan []domain.Stock)
	errCh := make(chan error, 1)

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for batch := range batchesCh {
				if err := s.Repository.Upsert(batch); err != nil {
					select {
					case errCh <- err:
					default:
					}
					return
				}
			}
		}()
	}

	go func() {
		defer close(batchesCh)

		var page *string
		var buffer []domain.Stock
		seenPages := make(map[string]bool)

		for {
			stocksPage, err := s.Provider.FetchStocks(page)
			if err != nil {
				errCh <- err
				return
			}

			if stocksPage.NextPage != "" {
				if seenPages[stocksPage.NextPage] {
					break
				}
				seenPages[stocksPage.NextPage] = true
				page = &stocksPage.NextPage
			} else {
				page = nil
			}

			for _, stock := range stocksPage.Items {
				buffer = append(buffer, stock)

				if len(buffer) == batchSize {
					batchesCh <- buffer
					buffer = nil
				}
			}

			if page == nil {
				break
			}
		}

		if len(buffer) > 0 {
			batchesCh <- buffer
		}
	}()

	wg.Wait()

	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}
