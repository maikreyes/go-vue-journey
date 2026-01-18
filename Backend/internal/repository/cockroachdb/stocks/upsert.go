package stocks

import (
	"backend/internal/domain"
	"context"
	"fmt"
	"strings"
	"time"
)

func (r *Repository) Upsert(stocks []domain.Stock) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	var (
		values []string
		args   []any
	)

	for i, s := range stocks {
		start := i*9 + 1

		values = append(values,
			fmt.Sprintf(
				"($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)",
				start, start+1, start+2, start+3,
				start+4, start+5, start+6, start+7, start+8,
			),
		)

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

	query := `
		INSERT INTO stocks (
			ticker,
			target_from,
			target_to,
			company,
			action,
			brokerage,
			rating_from,
			rating_to,
			time
		) VALUES ` + strings.Join(values, ",") + `
		 ON CONFLICT (ticker) DO UPDATE SET
			target_from = EXCLUDED.target_from,
			target_to = EXCLUDED.target_to,
			company = EXCLUDED.company,
			action = EXCLUDED.action,
			rating_from = EXCLUDED.rating_from,
			rating_to = EXCLUDED.rating_to;
	`

	_, err := r.db.Exec(ctx, query, args...)

	if err != nil {

		return err

	}

	return nil

}
