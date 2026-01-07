package stock

type StockFilter string

const (
	StockFilterAll  StockFilter = "all"
	StockFilterUp   StockFilter = "up"
	StockFilterDown StockFilter = "down"
)
