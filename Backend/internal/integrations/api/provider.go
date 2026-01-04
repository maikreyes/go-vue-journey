package api

import (
	"go-vue-journey/internal/stock"
)

type Provider struct {
	client *Client
}

func NewProvider(client *Client) *Provider {
	return &Provider{client: client}
}

func (s *Provider) GetStocks(page *string) (*stock.Page, error) {

	resp, err := s.client.GetStockResponse(page)
	if err != nil {
		return nil, err
	}

	return &stock.Page{
		Items:    resp.Items,
		NextPage: resp.NextPage,
	}, nil
}
