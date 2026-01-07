package stock

type StocksStats struct {
	Total int `json:"total"`
	Up    int `json:"up"`
	Down  int `json:"down"`
}

type StocksResponse struct {
	Items      []Stock     `json:"items"`
	Stats      StocksStats `json:"stats"`
	TotalPages int         `json:"total_pages"`
	NextCursor *string     `json:"next_cursor,omitempty"`
}
