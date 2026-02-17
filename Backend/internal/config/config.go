package config

import (
	"os"
	"strconv"
	"strings"
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
	providerURL := os.Getenv("API_ENDPOINT")

	token := os.Getenv("AUTHENTICATION")

	authorization := strings.TrimSpace(token)
	if authorization != "" && !strings.HasPrefix(strings.ToLower(authorization), "bearer ") {
		authorization = "Bearer " + authorization
	}

	return &Config{
		DSN:          os.Getenv("CONNECTION_STRING"),
		ProviderURL:  providerURL,
		Autorization: authorization,
		Port:         os.Getenv("PORT"),
		Workers:      getenvInt("WORKERS", 5),
		BatchSize:    getenvInt("BATCH_SIZE", 200),
		FrontendURL:  os.Getenv("FRONTEND_URL"),
	}
}
