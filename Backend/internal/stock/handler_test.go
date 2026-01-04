package stock_test

import (
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

	repository := &fakerepository{}
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

	// next_page NO enviado
	if service.receivedPage != nil {
		t.Fatalf("expected nil page, got %v", *service.receivedPage)
	}
}

func TestHandler_List_ServiceError(t *testing.T) {
	service := &fakeService{
		err: errors.New("boom"),
	}
	repository := &fakerepository{
		err: errors.New("boom"),
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
