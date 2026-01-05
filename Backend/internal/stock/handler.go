package stock

import (
	"encoding/json"
	"log"
	"net/http"
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

	if p := r.URL.Query().Get("next_page"); p != "" {
		page = &p
	}

	result, err := h.service.ListStocks(page)
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
