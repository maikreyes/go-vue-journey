package stocks

import (
	"backend/internal/domain"
	"context"
	"time"
)

func (r *Repository) GetTopStocks(limit int) (*[]domain.Stock, error) {

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
	WHERE rating_to = $1
		AND target_from IS NOT NULL
		AND REPLACE(REPLACE(target_from, '$', ''), ',', '')::FLOAT > 0
	ORDER BY
		(
			(
			REPLACE(REPLACE(target_to, '$', ''), ',', '')::FLOAT
			-
			REPLACE(REPLACE(target_from, '$', ''), ',', '')::FLOAT
			)
			/
			REPLACE(REPLACE(target_from, '$', ''), ',', '')::FLOAT
		) DESC
	LIMIT $2;
	`, "Buy", limit)

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

	return &stocks, nil

}
