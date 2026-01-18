package domain

type StocksPage struct {
	Items    []Stock `json:"items"`
	NextPage string  `json:"next_page"`
}
