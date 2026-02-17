package client

import (
	"backend/internal/domain"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) FetchStocks(page *string) (*domain.StocksPage, error) {
	if c.ApiURL == "" {
		return nil, errors.New("provider client: empty ApiURL (check PROVIDER_URL/API_ENDPOINT)")
	}

	req, err := http.NewRequest(http.MethodGet, c.ApiURL, nil)
	if err != nil {
		return nil, err
	}

	if page != nil {
		q := req.URL.Query()
		q.Add("next_page", *page)
		req.URL.RawQuery = q.Encode()
	}

	req.Header.Set("Authorization", c.Autorization)

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("provider: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var result domain.StocksPage

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil

}
