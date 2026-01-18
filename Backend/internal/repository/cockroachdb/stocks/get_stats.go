package stocks

import (
	"backend/internal/domain"
	"context"
	"time"
)

func (r *Repository) GetStats(limit int, filter *string, ticker *string) (*domain.StocksStats, error) {

	var stats domain.StocksStats

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	var tickerFilter any
	if ticker != nil && *ticker != "" {
		tickerFilter = *ticker
	}

	err := r.db.QueryRow(ctx, `
	SELECT
		COUNT(*) AS all_stocks,
		COUNT(CASE WHEN NULLIF(REPLACE(REPLACE(target_to, '$', ''), ',', ''), '')::FLOAT >
			NULLIF(REPLACE(REPLACE(target_from, '$', ''), ',', ''), '')::FLOAT THEN 1 END) AS up_stocks,
		COUNT(CASE WHEN NULLIF(REPLACE(REPLACE(target_to, '$', ''), ',', ''), '')::FLOAT <
			NULLIF(REPLACE(REPLACE(target_from, '$', ''), ',', ''), '')::FLOAT THEN 1 END) AS down_stocks,
		COUNT(CASE WHEN NULLIF(REPLACE(REPLACE(target_to, '$', ''), ',', ''), '')::FLOAT =
			NULLIF(REPLACE(REPLACE(target_from, '$', ''), ',', ''), '')::FLOAT THEN 1 END) AS equal_stocks
	FROM stocks
	WHERE ($1::TEXT IS NULL OR ticker LIKE ($1::TEXT || '%'));
	`, tickerFilter).Scan(
		&stats.AllStocks,
		&stats.UpStocks,
		&stats.DownStocks,
		&stats.NoChange,
	)

	if err != nil {
		return nil, err
	}

	if limit > 0 {
		switch {
		case filter != nil && *filter == "up":
			stats.Pages = (stats.UpStocks + limit - 1) / limit
		case filter != nil && *filter == "down":
			stats.Pages = (stats.DownStocks + limit - 1) / limit
		case filter != nil && *filter == "equal":
			stats.Pages = (stats.NoChange + limit - 1) / limit
		default:
			stats.Pages = (stats.AllStocks + limit - 1) / limit
		}
	}

	return &stats, nil

}
