package stock

import "backend/internal/provider/stock/client"

type Provider struct {
	Client *client.Client
}

func NewProvider(client *client.Client) *Provider {
	return &Provider{
		Client: client,
	}
}
