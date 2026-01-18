package domain

type StocksStats struct {
	AllStocks  int `json:"all_stocks"`
	UpStocks   int `json:"up_stocks"`
	DownStocks int `json:"down_stocks"`
	NoChange   int `json:"no_change"`
	Pages      int `json:"pages"`
}
