package stock_test

import (
	"encoding/json"
	"errors"
	"go-vue-journey/internal/stock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeService struct {
	receivedPage *string
	stocks       []stock.Stock
	err          error
}

type fakerepository struct {
	upsertedStocks []stock.Stock
	err            error
}

func (f *fakerepository) Upsert(s stock.Stock) error {
	if f.err != nil {
		return f.err
	}
	f.upsertedStocks = append(f.upsertedStocks, s)
	return nil
}

func (f *fakerepository) GetStocks(limit int, cursorTicker *string, filter stock.StockFilter) ([]stock.Stock, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.upsertedStocks, nil
}

func (f *fakerepository) GetStocksStats() (stock.StocksStats, error) {
	if f.err != nil {
		return stock.StocksStats{}, f.err
	}
	stats := stock.StocksStats{Total: len(f.upsertedStocks)}
	for _, s := range f.upsertedStocks {
		if s.TargetTo > s.TargetFrom {
			stats.Up++
		} else if s.TargetTo < s.TargetFrom {
			stats.Down++
		}
	}
	return stats, nil
}

func (f *fakerepository) GetTopStocks(n int) ([]stock.Stock, error) {
	return nil, nil
}

// GetStocks implements stock.StockProvider.
func (f *fakeService) GetStocks(page *string) (*stock.Page, error) {
	f.receivedPage = page
	if f.err != nil {
		return nil, f.err
	}
	return &stock.Page{Items: f.stocks, NextPage: nil}, nil
}

func TestHandler_List_NoNextPage(t *testing.T) {
	service := &fakeService{
		stocks: []stock.Stock{
			{Ticker: "AAA"},
			{Ticker: "BBB"},
		},
	}

	repository := &fakerepository{
		upsertedStocks: []stock.Stock{
			{Ticker: "AAA", TargetFrom: "10", TargetTo: "20"},
			{Ticker: "BBB", TargetFrom: "20", TargetTo: "10"},
		},
	}
	handler := stock.NewHandler(*stock.NewService(service, repository))

	req := httptest.NewRequest(http.MethodGet, "/stocks", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	res := w.Result()
	defer res.Body.Close()

	// Status
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	var body stock.StocksResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(body.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(body.Items))
	}
	if body.Stats.Total != 2 || body.Stats.Up != 1 || body.Stats.Down != 1 {
		t.Fatalf("unexpected stats: %+v", body.Stats)
	}
	if body.TotalPages != 1 {
		t.Fatalf("expected total_pages=1, got %d", body.TotalPages)
	}

	// Como el repo respondió, NO debería llamar al provider.
	if service.receivedPage != nil {
		t.Fatalf("expected provider not to be called, but got page %v", *service.receivedPage)
	}
}

func TestHandler_List_ServiceError(t *testing.T) {
	service := &fakeService{
		err: errors.New("boom"),
	}
	repository := &fakerepository{
		err: errors.New("db down"),
	}
	handler := stock.NewHandler(*stock.NewService(service, repository))

	req := httptest.NewRequest(http.MethodGet, "/stocks", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", res.StatusCode)
	}
}
