package main

import (
	stocksHanlder "backend/cmd/api/handlers/stocks"
	"backend/cmd/api/router"
	"backend/internal/config"
	"backend/internal/provider/stock"
	"backend/internal/provider/stock/client"
	"backend/internal/repository/cockroachdb"
	StocksRepository "backend/internal/repository/cockroachdb/stocks"
	LoggerRepository "backend/internal/repository/logger/stocks"
	stockService "backend/internal/services/stocks"
	"backend/internal/services/sync"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	ctg := config.Load()

	db, err := cockroachdb.ConnectDB(&ctg.DSN)

	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	defer db.Close()

	err = cockroachdb.Migrate(db)

	if err != nil {
		log.Fatalf("Error migrating the database: %v", err)
	}

	stockRepo := StocksRepository.NewRepository(db)
	logRepo := LoggerRepository.NewLoggerRepository(stockRepo)

	providerClient := client.NewClient(ctg.ProviderURL, ctg.Autorization)
	provider := stock.NewProvider(providerClient)

	syncService := sync.NewService(provider, logRepo, ctg.Workers, ctg.BatchSize)
	service := stockService.NewService(provider, logRepo)
	hanlder := stocksHanlder.NewHandler(service)

	router := router.NewRouter(hanlder)

	go syncService.Run()

	port := ":" + ctg.Port

	fmt.Printf("Servidor montado en localhost%s\n", port)

	server := http.ListenAndServe(port, router)

	if server != nil {
		panic(server)
	}

}
