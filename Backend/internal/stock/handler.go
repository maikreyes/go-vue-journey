package stock

import (
	"encoding/json"
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
