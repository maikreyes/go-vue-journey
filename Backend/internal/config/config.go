package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port           string
	ApiEndpoint    string
	Authentication string
	Dsn            string
	FrontendURL    string
	SyncBatchSize  int
}

func Load() *Config {
	batchSize := 10
	if v := os.Getenv("SYNC_BATCH_SIZE"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			batchSize = parsed
		}
	}

	return &Config{
		Port:           os.Getenv("PORT"),
		ApiEndpoint:    os.Getenv("API_ENDPOINT"),
		Authentication: os.Getenv("AUTHENTICATION"),
		Dsn:            os.Getenv("CONNECTION_STRING"),
		FrontendURL:    os.Getenv("FRONTEND_URL"),
		SyncBatchSize:  batchSize,
	}
}
