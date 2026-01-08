package config

import "os"

type Config struct {
	Port           string
	ApiEndpoint    string
	Authentication string
	Dsn            string
	FrontendURL    string
}

func Load() *Config {
	return &Config{
		Port:           os.Getenv("PORT"),
		ApiEndpoint:    os.Getenv("API_ENDPOINT"),
		Authentication: os.Getenv("AUTHENTICATION"),
		Dsn:            os.Getenv("CONNECTION_STRING"),
		FrontendURL:    os.Getenv("FRONTEND_URL"),
	}
}
