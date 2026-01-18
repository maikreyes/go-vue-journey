package stocks

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) GetTopStocks(w http.ResponseWriter, r *http.Request) {

	stocks, err := h.Service.GetTopStocks()

	if err != nil {
		http.Error(w, "Failed to fetch top stocks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stocks)
}
