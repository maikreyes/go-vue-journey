package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	baseUrl        string
	authentication string
	httpClient     *http.Client
}

func NewClient(baseUrl string, authentication string) *Client {
	return &Client{
		baseUrl:        baseUrl,
		authentication: authentication,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) GetStockResponse(page *string) (*StocksResponse, error) {

	req, err := http.NewRequest(http.MethodGet, c.baseUrl, nil)
	if err != nil {
		return nil, err
	}

	if page != nil {
		q := req.URL.Query()
		q.Set("next_page", *page)
		req.URL.RawQuery = q.Encode()
	}

	req.Header.Set("Authorization", "Bearer "+c.authentication)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var result StocksResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
