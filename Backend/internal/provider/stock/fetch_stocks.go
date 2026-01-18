package stock

import "backend/internal/domain"

func (p *Provider) FetchStocks(page *string) (*domain.StocksPage, error) {

	resp, err := p.Client.FetchStocks(page)

	if err != nil {
		return nil, err
	}

	return &domain.StocksPage{
		Items:    resp.Items,
		NextPage: resp.NextPage,
	}, nil
}
