package stock

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: &service,
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	var page *string
	var cursorTicker *string
	limit := 10
	filter := StockFilterAll

	if p := r.URL.Query().Get("next_page"); p != "" {
		page = &p
	}

	if c := r.URL.Query().Get("cursor"); c != "" {
		cursorTicker = &c
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if f := r.URL.Query().Get("filter"); f != "" {
		switch StockFilter(f) {
		case StockFilterAll, StockFilterUp, StockFilterDown:
			filter = StockFilter(f)
		}
	}

	result, err := h.service.ListStocksWithMeta(page, limit, cursorTicker, filter)
	if err != nil {
		http.Error(w, "failed to load stocks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *Handler) GetTopStocks(w http.ResponseWriter, r *http.Request) {
	limit := 5

	stocks, err := h.service.repository.GetTopStocks(limit)
	if err != nil {
		log.Println("GetTopStocks error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stocks)

}
