package client

import (
	"backend/internal/domain"
	"encoding/json"
	"net/http"
)

func (c *Client) FetchStocks(page *string) (*domain.StocksPage, error) {

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
		return nil, err
	}

	var result domain.StocksPage

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil

}
