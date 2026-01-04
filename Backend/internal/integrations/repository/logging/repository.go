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
