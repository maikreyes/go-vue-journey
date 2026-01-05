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

func (r *Repository) GetStocks() ([]stock.Stock, error) {
	start := time.Now()

	stocks, err := r.next.GetStocks()

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
