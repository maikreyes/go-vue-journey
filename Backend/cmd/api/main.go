package main

import (
	"fmt"
	"go-vue-journey/internal/config"
	"go-vue-journey/internal/integrations/api"
	"go-vue-journey/internal/router"
	"go-vue-journey/internal/stock"
	"net/http"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		fmt.Println(http.StatusNotFound)
	}

	ctg := config.Load()

	apiclient := api.NewClient(ctg.ApiEndpoint, ctg.Authentication)
	apiService := api.NewProvider(apiclient)

	stockService := stock.NewService(apiService)
	stockHandler := stock.NewHandler(*stockService)

	r := router.NewServerMux(*stockHandler)

	port := ":" + ctg.Port

	fmt.Printf("Servidor montado en localhost%s\n", port)

	server := http.ListenAndServe(port, r)

	if server != nil {
		panic(server)
	}

}
