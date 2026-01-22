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
	WHERE rating_to IN ('Buy', 'Strong-Buy', 'Outperform', 'Overweight')
		AND target_from IS NOT NULL
		AND REGEXP_REPLACE(target_from, '[^0-9.]', '', 'g')::FLOAT > 0
	ORDER BY
		(
          CASE
            WHEN rating_to = 'Strong-Buy' THEN 2
            WHEN rating_to = 'Buy' THEN 1
            ELSE 0
          END) DESC,
  		((
			REGEXP_REPLACE(target_to, '[^0-9.]', '', 'g')::FLOAT
			-
			REGEXP_REPLACE(target_from, '[^0-9.]', '', 'g')::FLOAT
			)
			/
			REGEXP_REPLACE(target_from, '[^0-9.]', '', 'g')::FLOAT
		) DESC
	LIMIT $1;
	`, limit)

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
