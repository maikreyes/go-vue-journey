package config

import (
	"os"
	"strconv"
)

type Config struct {
	DSN          string
	ProviderURL  string
	Autorization string
	Port         string
	Workers      int
	BatchSize    int
	FrontendURL  string
}

func getenvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func Load() *Config {
	return &Config{
		DSN:          os.Getenv("CONNECTION_STRING"),
		ProviderURL:  os.Getenv("PROVIDER_URL"),
		Autorization: "Bearer " + os.Getenv("AUTORIZATION_TOKEN"),
		Port:         os.Getenv("PORT"),
		Workers:      getenvInt("WORKERS", 0),
		BatchSize:    getenvInt("BATCH_SIZE", 0),
		FrontendURL:  os.Getenv("FRONTEND_URL"),
	}
}
