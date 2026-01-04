package api

import "go-vue-journey/internal/stock"

type StocksResponse struct {
	Items    []stock.Stock `json:"items"`
	NextPage *string       `json:"next_page"`
}
