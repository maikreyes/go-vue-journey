package stock_test

import (
	"errors"
	"go-vue-journey/internal/stock"
	"testing"
)

type fakeProvider struct {
	called bool
}

func (f *fakeProvider) GetStocks(page *string) (*stock.Page, error) {
	f.called = true

	next := "page2"
	return &stock.Page{
		Items: []stock.Stock{
			{Ticker: "AAA"},
			{Ticker: "BBB"},
			{Ticker: "CCC"},
		},
		NextPage: &next,
	}, nil
}

type fakeRepo struct {
	saved []stock.Stock
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{
		saved: []stock.Stock{},
	}
}

func (f *fakeRepo) Upsert(s stock.Stock) error {
	f.saved = append(f.saved, s)
	return nil
}

func (f *fakeRepo) GetStocks() ([]stock.Stock, error) {
	return f.saved, nil
}

func (f *fakeRepo) GetTopStocks(n int) ([]stock.Stock, error) {
	return nil, nil
}

type failingRepo struct{}

func (f failingRepo) Upsert(stock.Stock) error {
	return errors.New("db down")
}

func (f failingRepo) GetStocks() ([]stock.Stock, error) {
	return nil, errors.New("db down")
}

func (f failingRepo) GetTopStocks(n int) ([]stock.Stock, error) {
	return nil, errors.New("db down")
}

func TestService_ListStocks_FetchAndStore(t *testing.T) {
	provider := &fakeProvider{}
	repo := newFakeRepo()

	service := stock.NewService(provider, repo)

	result, err := service.ListStocks(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(*result) != 3 {
		t.Fatalf("expected 3 items, got %d", len(*result))
	}

	if !provider.called {
		t.Fatal("expected provider to be called")
	}

	// Verify all stocks were saved to repo
	if len(repo.saved) != 3 {
		t.Fatalf("expected 3 stocks to be saved, got %d", len(repo.saved))
	}

	saved := map[string]bool{}
	for _, stock := range repo.saved {
		saved[stock.Ticker] = true
	}

	expected := []string{"AAA", "BBB", "CCC"}
	for _, tkr := range expected {
		if !saved[tkr] {
			t.Errorf("expected stock %s to be saved", tkr)
		}
	}
}

func TestService_ListStocks_RepoFailureDoesNotBreakResponse(t *testing.T) {
	provider := &fakeProvider{}
	repo := failingRepo{}

	service := stock.NewService(provider, repo)

	result, err := service.ListStocks(nil)
	if err != nil {
		t.Fatal("expected no error even if repo fails")
	}

	// Should fallback to API results when repo fails
	if len(*result) != 3 {
		t.Fatalf("expected 3 items from API fallback, got %d", len(*result))
	}
}
