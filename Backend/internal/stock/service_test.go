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
	saved chan stock.Stock
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{
		saved: make(chan stock.Stock, 10),
	}
}

func (f *fakeRepo) Upsert(s stock.Stock) error {
	f.saved <- s
	return nil
}

type failingRepo struct{}

func (f failingRepo) Upsert(stock.Stock) error {
	return errors.New("db down")
}

func TestService_ListStocks_FetchAndStore(t *testing.T) {
	provider := &fakeProvider{}
	repo := newFakeRepo()

	service := stock.NewService(provider, repo)

	result, err := service.ListStocks(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(result.Items))
	}

	if !provider.called {
		t.Fatal("expected provider to be called")
	}

	saved := map[string]bool{}

	for i := 0; i < 3; i++ {
		stock := <-repo.saved
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

	_, err := service.ListStocks(nil)
	if err != nil {
		t.Fatal("expected no error even if repo fails")
	}
}
