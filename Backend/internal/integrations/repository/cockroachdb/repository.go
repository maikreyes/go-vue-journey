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

func (r *Repository) GetStocks(limit int, cursorTicker *string, filter stock.StockFilter) ([]stock.Stock, error) {
	var (
		rows *sql.Rows
		err  error
	)

	whereFilter := ""
	switch filter {
	case stock.StockFilterUp:
		whereFilter = `
			AND (
				NULLIF(REPLACE(REPLACE(target_to, '$', ''), ',', ''), '')::numeric
				>
				NULLIF(REPLACE(REPLACE(target_from, '$', ''), ',', ''), '')::numeric
			)
		`
	case stock.StockFilterDown:
		whereFilter = `
			AND (
				NULLIF(REPLACE(REPLACE(target_to, '$', ''), ',', ''), '')::numeric
				<
				NULLIF(REPLACE(REPLACE(target_from, '$', ''), ',', ''), '')::numeric
			)
		`
	}

	if cursorTicker == nil || *cursorTicker == "" {
		query := "SELECT * FROM stocks WHERE 1=1 " + whereFilter + " ORDER BY ticker ASC LIMIT $1;"
		rows, err = r.db.Query(query, limit)
	} else {
		query := "SELECT * FROM stocks WHERE ticker > $1 " + whereFilter + " ORDER BY ticker ASC LIMIT $2;"
		rows, err = r.db.Query(query, *cursorTicker, limit)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var stocks []stock.Stock

	for rows.Next() {
		var s stock.Stock
		err := rows.Scan(
			&s.Ticker,
			&s.TargetFrom,
			&s.TargetTo,
			&s.Company,
			&s.Action,
			&s.Brokerage,
			&s.RatingFrom,
			&s.RatingTo,
			&s.Time,
		)
		if err != nil {
			return nil, err
		}
		stocks = append(stocks, s)
	}
	return stocks, nil
}

func (r *Repository) GetStocksStats() (stock.StocksStats, error) {
	var stats stock.StocksStats

	err := r.db.QueryRow(`
		SELECT
			COUNT(*)::INT AS total,
			SUM(
				CASE
					WHEN NULLIF(REPLACE(REPLACE(target_to, '$', ''), ',', ''), '')::numeric
						>
						NULLIF(REPLACE(REPLACE(target_from, '$', ''), ',', ''), '')::numeric
					THEN 1
					ELSE 0
				END
			)::INT AS up,
			SUM(
				CASE
					WHEN NULLIF(REPLACE(REPLACE(target_to, '$', ''), ',', ''), '')::numeric
						<
						NULLIF(REPLACE(REPLACE(target_from, '$', ''), ',', ''), '')::numeric
					THEN 1
					ELSE 0
				END
			)::INT AS down
		FROM stocks;
	`).Scan(&stats.Total, &stats.Up, &stats.Down)

	if err != nil {
		return stock.StocksStats{}, err
	}

	return stats, nil
}

func (r *Repository) GetTopStocks(limit int) ([]stock.Stock, error) {

	rows, err := r.db.Query(`
		SELECT *
		FROM stocks
		WHERE rating_to = $1
		ORDER BY
		(
		(
			REPLACE(REPLACE(target_to, '$', ''), ',', '')::numeric
		- REPLACE(REPLACE(target_from, '$', ''), ',', '')::numeric
		)
		/
		NULLIF(
			REPLACE(REPLACE(target_from, '$', ''), ',', '')::numeric,
			0
		)
		) DESC
		LIMIT $2;
		`,
		"Buy",
		limit,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var stocks []stock.Stock
	for rows.Next() {
		var s stock.Stock
		err := rows.Scan(
			&s.Ticker,
			&s.TargetFrom,
			&s.TargetTo,
			&s.Company,
			&s.Action,
			&s.Brokerage,
			&s.RatingFrom,
			&s.RatingTo,
			&s.Time,
		)
		if err != nil {
			return nil, err
		}
		stocks = append(stocks, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stocks, nil
}
