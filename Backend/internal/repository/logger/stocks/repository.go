package stocks

import "backend/internal/repository/cockroachdb/stocks"

type Repository struct {
	Repository *stocks.Repository
}

func NewLoggerRepository(repo *stocks.Repository) *Repository {
	return &Repository{
		Repository: repo,
	}
}
