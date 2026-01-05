package router

import (
	"go-vue-journey/internal/middleware"
	"go-vue-journey/internal/stock"
	"net/http"

	"github.com/gorilla/mux"
)

func NewServerMux(stockHandler stock.Handler) http.Handler {
	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/stocks", stockHandler.List).
		Methods(http.MethodGet, http.MethodOptions)

	api.HandleFunc("/stocks/top", stockHandler.GetTopStocks).
		Methods(http.MethodGet, http.MethodOptions)

	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World"))
	})

	routerWithCors := middleware.CORS(r)

	return routerWithCors
}
