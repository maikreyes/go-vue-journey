package cockroachdb

import (
	"database/sql"
	"fmt"
	"go-vue-journey/internal/stock"
	"strings"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Upsert(s stock.Stock) error {
	return r.UpsertMany([]stock.Stock{s})
}

func (r *Repository) UpsertMany(items []stock.Stock) error {
	if len(items) == 0 {
		return nil
	}

	var b strings.Builder
	b.WriteString("INSERT INTO stocks (ticker, target_from, target_to, company, action, brokerage, rating_from, rating_to, time) VALUES ")

	args := make([]any, 0, len(items)*9)
	for i, s := range items {
		if i > 0 {
			b.WriteString(",")
		}

		base := i*9 + 1
		b.WriteString("(")
		b.WriteString(fmt.Sprintf("$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d", base, base+1, base+2, base+3, base+4, base+5, base+6, base+7, base+8))
		b.WriteString(")")

		args = append(args,
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
	}

	b.WriteString(" ON CONFLICT (ticker) DO UPDATE SET ")
	b.WriteString("target_from = EXCLUDED.target_from, ")
	b.WriteString("target_to = EXCLUDED.target_to, ")
	b.WriteString("company = EXCLUDED.company, ")
	b.WriteString("action = EXCLUDED.action, ")
	b.WriteString("brokerage = EXCLUDED.brokerage, ")
	b.WriteString("rating_from = EXCLUDED.rating_from, ")
	b.WriteString("rating_to = EXCLUDED.rating_to, ")
	b.WriteString("time = EXCLUDED.time")

	_, err := r.db.Exec(b.String(), args...)
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

func (r *Repository) GetStocksByTicker(tickerPrefix string, limit int, cursorTicker *string) ([]stock.Stock, error) {
	var (
		rows *sql.Rows
		err  error
	)

	like := tickerPrefix + "%"

	if limit <= 0 {
		limit = 10
	}

	if cursorTicker == nil || *cursorTicker == "" {
		rows, err = r.db.Query(`
			SELECT *
			FROM stocks
			WHERE ticker LIKE $1
			ORDER BY ticker ASC
			LIMIT $2;
			`,
			like,
			limit,
		)
	} else {
		rows, err = r.db.Query(`
			SELECT *
			FROM stocks
			WHERE ticker LIKE $1
				AND ticker > $2
			ORDER BY ticker ASC
			LIMIT $3;
			`,
			like,
			*cursorTicker,
			limit,
		)
	}
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	stocks := make([]stock.Stock, 0, limit)

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

func (r *Repository) CountStocksByTicker(tickerPrefix string) (int, error) {
	// Deprecated: mantenido temporalmente para evitar romper callers antiguos.
	stats, err := r.GetStocksStatsByTicker(tickerPrefix)
	if err != nil {
		return 0, err
	}
	return stats.Total, nil
}

func (r *Repository) GetStocksStatsByTicker(tickerPrefix string) (stock.StocksStats, error) {
	like := tickerPrefix + "%"
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
		FROM stocks
		WHERE ticker LIKE $1;
		`,
		like,
	).Scan(&stats.Total, &stats.Up, &stats.Down)
	if err != nil {
		return stock.StocksStats{}, err
	}
	return stats, nil
}
