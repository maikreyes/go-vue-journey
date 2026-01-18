package stocks

import (
	"backend/internal/domain"
	"context"
	"time"
)

func (r *Repository) GetStocks(page *string, limit int) (*domain.StocksPage, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	rows, err := r.db.Query(ctx, `
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
	LIMIT $2
	`, page, limit)

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
