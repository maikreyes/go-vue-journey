package stock

type Page struct {
	Items    []Stock `json:"items"`
	NextPage *string `json:"next_page"`
}
