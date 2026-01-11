package logging

import (
	"go-vue-journey/internal/stock"
	"log"
	"time"
)

type Repository struct {
	next stock.StockRepository
}

func New(next stock.StockRepository) *Repository {
	return &Repository{next: next}
}

func (r *Repository) Upsert(s stock.Stock) error {
	start := time.Now()

	err := r.next.Upsert(s)

	elapsed := time.Since(start)

	if err != nil {
		log.Printf(
			"[STOCK][ERROR] ticker=%s duration=%s err=%v",
			s.Ticker, elapsed, err,
		)
		return err
	}

	log.Printf(
		"[STOCK][UPSERT] ticker=%s duration=%s",
		s.Ticker, elapsed,
	)

	return nil
}

func (r *Repository) UpsertMany(items []stock.Stock) error {
	start := time.Now()

	if batchRepo, ok := r.next.(interface{ UpsertMany([]stock.Stock) error }); ok {
		err := batchRepo.UpsertMany(items)
		elapsed := time.Since(start)

		if err != nil {
			log.Printf(
				"[STOCK][ERROR] UpsertMany count=%d duration=%s err=%v",
				len(items), elapsed, err,
			)
			return err
		}

		log.Printf(
			"[STOCK][UPSERT_MANY] count=%d duration=%s",
			len(items), elapsed,
		)
		return nil
	}

	for _, s := range items {
		if err := r.next.Upsert(s); err != nil {
			elapsed := time.Since(start)
			log.Printf(
				"[STOCK][ERROR] UpsertMany(fallback) count=%d duration=%s err=%v",
				len(items), elapsed, err,
			)
			return err
		}
	}

	elapsed := time.Since(start)
	log.Printf(
		"[STOCK][UPSERT_MANY(fallback)] count=%d duration=%s",
		len(items), elapsed,
	)
	return nil
}

func (r *Repository) GetStocks(limit int, cursorTicker *string, filter stock.StockFilter) ([]stock.Stock, error) {
	start := time.Now()

	stocks, err := r.next.GetStocks(limit, cursorTicker, filter)

	elapsed := time.Since(start)

	if err != nil {
		log.Printf(
			"[STOCK][ERROR] GetStocks duration=%s error=%v",
			elapsed, err,
		)
		return nil, err
	}

	log.Printf(
		"[STOCK][GET_STOCKS] duration=%s count=%d",
		elapsed, len(stocks),
	)
	return stocks, nil
}

func (r *Repository) GetStocksStats() (stock.StocksStats, error) {
	start := time.Now()

	stats, err := r.next.GetStocksStats()

	elapsed := time.Since(start)

	if err != nil {
		log.Printf(
			"[STOCK][ERROR] GetStocksStats duration=%s error=%v",
			elapsed, err,
		)
		return stock.StocksStats{}, err
	}

	log.Printf(
		"[STOCK][GET_STATS] duration=%s total=%d up=%d down=%d",
		elapsed, stats.Total, stats.Up, stats.Down,
	)

	return stats, nil
}

func (r *Repository) GetTopStocks(limit int) ([]stock.Stock, error) {
	start := time.Now()

	stocks, err := r.next.GetTopStocks(limit)

	elapsed := time.Since(start)

	if err != nil {
		log.Printf(
			"[STOCK][ERROR] GetTopStocks limit=%d duration=%s err=%v",
			limit, elapsed, err,
		)
		return nil, err
	}

	log.Printf(
		"[STOCK][GET_TOP] limit=%d duration=%s count=%d",
		limit, elapsed, len(stocks),
	)

	return stocks, nil
}

func (r *Repository) GetStocksByTicker(tickerPrefix string, limit int, cursorTicker *string) ([]stock.Stock, error) {
	start := time.Now()
	items, err := r.next.GetStocksByTicker(tickerPrefix, limit, cursorTicker)
	elapsed := time.Since(start)

	if err != nil {
		log.Printf(
			"[STOCK][ERROR] GetStocksByTicker prefix=%s limit=%d duration=%s err=%v",
			tickerPrefix, limit, elapsed, err,
		)
		return nil, err
	}
	log.Printf(
		"[STOCK][GET_BY_TICKER] prefix=%s limit=%d duration=%s count=%d",
		tickerPrefix, limit, elapsed, len(items),
	)
	return items, nil
}

func (r *Repository) GetStocksStatsByTicker(tickerPrefix string) (stock.StocksStats, error) {
	start := time.Now()
	stats, err := r.next.GetStocksStatsByTicker(tickerPrefix)
	elapsed := time.Since(start)
	if err != nil {
		log.Printf(
			"[STOCK][ERROR] GetStocksStatsByTicker prefix=%s duration=%s err=%v",
			tickerPrefix, elapsed, err,
		)
		return stock.StocksStats{}, err
	}
	log.Printf(
		"[STOCK][STATS_BY_TICKER] prefix=%s duration=%s total=%d up=%d down=%d",
		tickerPrefix, elapsed, stats.Total, stats.Up, stats.Down,
	)
	return stats, nil
}
