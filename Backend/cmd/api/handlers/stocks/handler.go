package stocks

import "backend/internal/services/stocks"

type Handler struct {
	Service *stocks.Service
}

func NewHandler(service *stocks.Service) *Handler {
	return &Handler{Service: service}
}
