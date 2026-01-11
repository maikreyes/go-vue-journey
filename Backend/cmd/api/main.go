package main

import (
	"fmt"
	"go-vue-journey/internal/config"
	"go-vue-journey/internal/integrations/api"
	"go-vue-journey/internal/integrations/repository/cockroachdb"
	"go-vue-journey/internal/integrations/repository/db"
	"go-vue-journey/internal/integrations/repository/logging"
	"go-vue-journey/internal/integrations/sync"
	"go-vue-journey/internal/router"
	"go-vue-journey/internal/stock"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		fmt.Println(http.StatusNotFound)
	}

	ctg := config.Load()

	apiclient := api.NewClient(ctg.ApiEndpoint, ctg.Authentication)
	apiService := api.NewProvider(apiclient)

	DB, err := db.Connect(ctg.Dsn)

	if err != nil {
		log.Fatal(err)
	}

	cockroachdb.Migrate(DB)

	repo := cockroachdb.New(DB)
	repoWithLogging := logging.New(repo)

	syncService := sync.NewService(apiService, repoWithLogging, 5, ctg.SyncBatchSize)

	go syncService.Run()

	stockService := stock.NewService(apiService, repoWithLogging)
	stockHandler := stock.NewHandler(*stockService)

	r := router.NewServerMux(*stockHandler)

	port := ":" + ctg.Port

	fmt.Printf("Servidor montado en localhost%s\n", port)

	server := http.ListenAndServe(port, r)

	if server != nil {
		panic(server)
	}

}
