package stocks

import (
	"backend/internal/domain"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (h *Handler) GetStocks(w http.ResponseWriter, r *http.Request) {

	var page *string
	var filter *string
	var ticker *string

	queryValues := r.URL.Query()

	if nextPage := queryValues.Get("next_page"); nextPage != "" {
		page = &nextPage
	}

	if filterParam := queryValues.Get("filter"); filterParam != "" {
		filter = &filterParam
	}

	if tickerParam := queryValues.Get("ticker"); tickerParam != "" {
		tickerUpper := strings.ToUpper(tickerParam)
		ticker = &tickerUpper
	}

	stats, err := h.Service.GetStats(filter, ticker)

	if err != nil {
		http.Error(w, "Failed to fetch stocks stats", http.StatusInternalServerError)
		log.Println("Error fetching stocks stats:", err)
		return
	}

	if ticker != nil {
		stocks, err := h.Service.GetStockByTicker(*ticker, page, filter)

		if err != nil {
			http.Error(w, "Failed to get stock by ticker", http.StatusInternalServerError)
			log.Println("Error fetching stock by ticker:", err)
			return
		}

		response := domain.ApiResponse{
			Items:      stocks.Items,
			Stats:      *stats,
			NextCursor: &stocks.NextPage,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	if filter != nil {
		stocks, err := h.Service.GetFilterStocks(page, filter)

		if err != nil {
			http.Error(w, "Failed to get up stocks", http.StatusInternalServerError)
			log.Println("Error fetching up stocks:", err)
			return
		}

		response := domain.ApiResponse{
			Items:      stocks.Items,
			Stats:      *stats,
			NextCursor: &stocks.NextPage,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	stocks, err := h.Service.GetStocks(page)

	if err != nil {
		http.Error(w, "Failed to get stocks", http.StatusInternalServerError)
		log.Println("Error fetching stocks:", err)
		return
	}

	response := domain.ApiResponse{
		Items:      stocks.Items,
		Stats:      *stats,
		NextCursor: &stocks.NextPage,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
