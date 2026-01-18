package client

import (
	"net/http"
	"time"
)

type Client struct {
	ApiURL       string
	Autorization string
	httpClient   *http.Client
}

func NewClient(apiURL, autorization string) *Client {
	return &Client{
		ApiURL:       apiURL,
		Autorization: autorization,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}
