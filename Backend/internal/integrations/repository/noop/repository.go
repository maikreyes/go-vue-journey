package noop

import "go-vue-journey/internal/stock"

type Repository struct{}

func New() *Repository {
	return &Repository{}
}

func (r *Repository) Upsert(stock.Stock) error {
	return nil
}
