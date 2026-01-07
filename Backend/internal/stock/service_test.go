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

func (f *fakeRepo) GetStocks(limit int, cursorTicker *string, filter stock.StockFilter) ([]stock.Stock, error) {
	return f.saved, nil
}

func (f *fakeRepo) GetStocksStats() (stock.StocksStats, error) {
	stats := stock.StocksStats{Total: len(f.saved)}
	for _, s := range f.saved {
		if s.TargetTo > s.TargetFrom {
			stats.Up++
		} else if s.TargetTo < s.TargetFrom {
			stats.Down++
		}
	}
	return stats, nil
}

func (f *fakeRepo) GetTopStocks(n int) ([]stock.Stock, error) {
	return nil, nil
}

type failingRepo struct{}

func (f failingRepo) Upsert(stock.Stock) error {
	return errors.New("db down")
}

func (f failingRepo) GetStocks(limit int, cursorTicker *string, filter stock.StockFilter) ([]stock.Stock, error) {
	return nil, errors.New("db down")
}

func (f failingRepo) GetStocksStats() (stock.StocksStats, error) {
	return stock.StocksStats{}, errors.New("db down")
}

func (f failingRepo) GetTopStocks(n int) ([]stock.Stock, error) {
	return nil, errors.New("db down")
}

func TestService_ListStocks_FetchAndStore(t *testing.T) {
	provider := &fakeProvider{}
	repo := newFakeRepo()
	repo.saved = []stock.Stock{
		{Ticker: "AAA", TargetFrom: "10", TargetTo: "20"},
		{Ticker: "BBB", TargetFrom: "20", TargetTo: "10"},
		{Ticker: "CCC", TargetFrom: "10", TargetTo: "20"},
	}

	service := stock.NewService(provider, repo)

	result, err := service.ListStocks(nil, 10, nil, stock.StockFilterAll)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(*result) != 3 {
		t.Fatalf("expected 3 items, got %d", len(*result))
	}

	if provider.called {
		t.Fatal("expected provider NOT to be called when repo works")
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

func TestService_ListStocksWithMeta_ReturnsStatsAndPages(t *testing.T) {
	provider := &fakeProvider{}
	repo := newFakeRepo()
	repo.saved = []stock.Stock{
		{Ticker: "AAA", TargetFrom: "10", TargetTo: "20"},
		{Ticker: "BBB", TargetFrom: "20", TargetTo: "10"},
		{Ticker: "CCC", TargetFrom: "10", TargetTo: "20"},
	}

	service := stock.NewService(provider, repo)

	resp, err := service.ListStocksWithMeta(nil, 2, nil, stock.StockFilterAll)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Stats.Total != 3 {
		t.Fatalf("expected total=3, got %d", resp.Stats.Total)
	}
	if resp.Stats.Up != 2 {
		t.Fatalf("expected up=2, got %d", resp.Stats.Up)
	}
	if resp.Stats.Down != 1 {
		t.Fatalf("expected down=1, got %d", resp.Stats.Down)
	}
	if resp.TotalPages != 2 {
		t.Fatalf("expected total_pages=2, got %d", resp.TotalPages)
	}
	if resp.NextCursor == nil || *resp.NextCursor != "CCC" {
		if resp.NextCursor == nil {
			t.Fatalf("expected next_cursor to be set")
		}
		t.Fatalf("expected next_cursor=CCC, got %s", *resp.NextCursor)
	}
}

func TestService_ListStocks_RepoFailureDoesNotBreakResponse(t *testing.T) {
	provider := &fakeProvider{}
	repo := failingRepo{}

	service := stock.NewService(provider, repo)

	result, err := service.ListStocks(nil, 10, nil, stock.StockFilterAll)
	if err != nil {
		t.Fatal("expected no error even if repo fails")
	}

	// Should fallback to API results when repo fails
	if len(*result) != 3 {
		t.Fatalf("expected 3 items from API fallback, got %d", len(*result))
	}

	if !provider.called {
		t.Fatal("expected provider to be called on repo failure")
	}
}
