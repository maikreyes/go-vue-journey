package stocks

import (
	"backend/internal/domain"
	"context"
	"time"
)

func (r *Repository) GetFilterStocks(page *string, limit int, filter *string) (*domain.StocksPage, error) {

	operator := ""

	switch {
	case filter != nil && *filter == "up":
		operator = ">"
	case filter != nil && *filter == "down":
		operator = "<"
	case filter != nil && *filter == "equal":
		operator = "="
	}

	query := `
	SELECT
		ticker,
		target_from,
		target_to,
		company,
		action,
		brokerage,
		rating_from,
		rating_to,
		time
	FROM stocks
	WHERE ($1::TEXT IS NULL OR ticker > $1::TEXT)
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.Query(ctx, query, page, limit)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var stocks []domain.Stock

	for rows.Next() {
		var stock domain.Stock
		err := rows.Scan(
			&stock.Ticker,
			&stock.TargetFrom,
			&stock.TargetTo,
			&stock.Company,
			&stock.Action,
			&stock.Brokerage,
			&stock.RatingFrom,
			&stock.RatingTo,
			&stock.Time,
		)

		if err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}

	var nextPage string

	if len(stocks) > 0 {
		nextPage = stocks[len(stocks)-1].Ticker
	}

	return &domain.StocksPage{
		Items:    stocks,
		NextPage: nextPage,
	}, nil

}
