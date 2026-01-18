package router

import (
	"backend/cmd/api/handlers/stocks"
	"backend/internal/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(handler *stocks.Handler) *mux.Router {

	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()

	v1 := api.PathPrefix("/v1").Subrouter()

	v1.HandleFunc("/stocks", handler.GetStocks).Methods(http.MethodGet, http.MethodOptions)
	v1.HandleFunc("/stocks/top", handler.GetTopStocks).Methods(http.MethodGet, http.MethodOptions)

	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World"))
	})

	r.Use(middleware.CORS)

	return r

}
