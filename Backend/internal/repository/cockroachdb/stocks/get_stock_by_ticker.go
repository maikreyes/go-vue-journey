package stocks

import (
	"backend/internal/domain"
	"context"
	"time"
)

func (r *Repository) GetStockByTicker(ticker string, limit int, page *string, filter *string) (*domain.StocksPage, error) {

	operator := ""

	switch {
	case filter != nil && *filter == "up":
		operator = ">"
	case filter != nil && *filter == "down":
		operator = "<"
	case filter != nil && *filter == "equal":
		operator = "="
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT *
			FROM stocks
			WHERE ticker LIKE $1
			AND ($3::TEXT IS NULL OR ticker > $3::TEXT)
			`

	if operator != "" {
		query += `
			AND (
				NULLIF(REPLACE(REPLACE(target_to, '$', ''), ',', ''), '')::FLOAT
				` + operator + `
				NULLIF(REPLACE(REPLACE(target_from, '$', ''), ',', ''), '')::FLOAT
			)
			`
	}
	query += `
	ORDER BY ticker ASC
	LIMIT $2;
	`

	rows, err := r.db.Query(ctx, query, ticker+"%", limit, page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []domain.Stock

	for rows.Next() {
		var stock domain.Stock
		if err := rows.Scan(
			&stock.Ticker,
			&stock.TargetFrom,
			&stock.TargetTo,
			&stock.Company,
			&stock.Action,
			&stock.Brokerage,
			&stock.RatingFrom,
			&stock.RatingTo,
			&stock.Time,
		); err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var nextPage string

	if len(stocks) > 0 {
		nextPage = stocks[len(stocks)-1].Ticker
	}

	stocksPage := domain.StocksPage{
		Items:    stocks,
		NextPage: nextPage,
	}

	return &stocksPage, nil
}
