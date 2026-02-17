package sync

import (
	"backend/internal/domain"
	"errors"
	"sync"
)

func (s *Service) Run() error {
	batchSize := s.BatchSize
	workers := s.Workers

	if workers <= 0 {
		return errors.New("sync: WORKERS must be > 0")
	}
	if batchSize <= 0 {
		return errors.New("sync: BATCH_SIZE must be > 0")
	}

	batchesCh := make(chan []domain.Stock)
	errCh := make(chan error, 1)
	stopCh := make(chan struct{})
	var stopOnce sync.Once
	fail := func(err error) {
		if err == nil {
			return
		}
		stopOnce.Do(func() {
			select {
			case errCh <- err:
			default:
			}
			close(stopCh)
		})
	}

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for batch := range batchesCh {
				if err := s.Repository.Upsert(batch); err != nil {
					fail(err)
					return
				}
			}
		}()
	}

	producerDone := make(chan struct{})
	go func() {
		defer close(producerDone)
		defer close(batchesCh)

		var page *string
		var buffer []domain.Stock
		seenPages := make(map[string]bool)

		for {
			select {
			case <-stopCh:
				return
			default:
			}

			stocksPage, err := s.Provider.FetchStocks(page)
			if err != nil {
				fail(err)
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
					select {
					case <-stopCh:
						return
					case batchesCh <- buffer:
					}
					buffer = nil
				}
			}

			if page == nil {
				break
			}
		}

		if len(buffer) > 0 {
			select {
			case <-stopCh:
				return
			case batchesCh <- buffer:
			}
		}
	}()

	// Wait for the producer to finish closing the channel,
	// then for all workers to drain it.
	<-producerDone
	wg.Wait()

	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}
