package domain

type ApiResponse struct {
	Items      []Stock     `json:"items"`
	Stats      StocksStats `json:"stats"`
	NextCursor *string     `json:"next_cursor,omitempty"`
}
