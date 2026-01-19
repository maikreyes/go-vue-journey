package stocks

import (
	"backend/internal/domain"
	"encoding/json"
	"net/http"
)

func (h *Handler) Write(stockPage domain.StocksPage, stats domain.StocksStats, w http.ResponseWriter) error {
	response := domain.ApiResponse{
		Items:      stockPage.Items,
		Stats:      stats,
		NextCursor: &stockPage.NextPage,
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
