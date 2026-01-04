package cockroachdb

import (
	"database/sql"
	"go-vue-journey/internal/stock"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Upsert(s stock.Stock) error {
	_, err := r.db.Exec(`
		INSERT INTO stocks (
			ticker, target_from, target_to, company,
			action, brokerage, rating_from, rating_to, time
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9
		)
		ON CONFLICT (ticker) DO UPDATE SET
			target_from = EXCLUDED.target_from,
			target_to = EXCLUDED.target_to,
			company = EXCLUDED.company,
			action = EXCLUDED.action,
			brokerage = EXCLUDED.brokerage,
			rating_from = EXCLUDED.rating_from,
			rating_to = EXCLUDED.rating_to,
			time = EXCLUDED.time
	`,
		s.Ticker,
		s.TargetFrom,
		s.TargetTo,
		s.Company,
		s.Action,
		s.Brokerage,
		s.RatingFrom,
		s.RatingTo,
		s.Time,
	)

	return err
}
