package router

import (
	"go-vue-journey/internal/stock"
	"net/http"

	"github.com/gorilla/mux"
)

func NewServerMux(stockHandler stock.Handler) *mux.Router {
	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/stocks", stockHandler.List).Methods(http.MethodGet)
	api.HandleFunc("/stocks/sync", stockHandler.Run).Methods(http.MethodPost)

	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World"))
	})

	return r
}
