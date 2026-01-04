package sync_test

import (
	"errors"
	"go-vue-journey/internal/integrations/sync"
	"go-vue-journey/internal/stock"
	"testing"
	"time"
)

// MockStockProvider es un mock del StockProvider para testing
type MockStockProvider struct {
	callCount int
	responses []*stock.Page
	errors    []error
}

func (m *MockStockProvider) GetStocks(page *string) (*stock.Page, error) {
	if m.callCount >= len(m.responses) {
		return nil, errors.New("unexpected call to GetStocks")
	}

	err := m.errors[m.callCount]
	response := m.responses[m.callCount]
	m.callCount++

	return response, err
}

// MockStockRepository es un mock del StockRepository para testing
type MockStockRepository struct {
	UpsertedStocks []stock.Stock
	UpsertError    error
}

func (m *MockStockRepository) Upsert(s stock.Stock) error {
	if m.UpsertError != nil {
		return m.UpsertError
	}
	m.UpsertedStocks = append(m.UpsertedStocks, s)
	return nil
}

func TestNewService(t *testing.T) {
	provider := &MockStockProvider{}
	repository := &MockStockRepository{}
	workers := 5

	svc := sync.NewService(provider, repository, workers)

	if svc == nil {
		t.Error("NewService debería retornar un servicio no nil")
	}
}

func TestRun_Success_SinglePage(t *testing.T) {
	// Arrange
	stocks := []stock.Stock{
		{
			Ticker:     "AAPL",
			Company:    "Apple",
			Action:     "BUY",
			Brokerage:  "Goldman",
			RatingFrom: "HOLD",
			RatingTo:   "BUY",
			Time:       time.Now(),
		},
		{
			Ticker:     "GOOGL",
			Company:    "Google",
			Action:     "SELL",
			Brokerage:  "Morgan",
			RatingFrom: "BUY",
			RatingTo:   "HOLD",
			Time:       time.Now(),
		},
	}

	provider := &MockStockProvider{
		responses: []*stock.Page{
			{
				Items:    stocks,
				NextPage: nil,
			},
		},
		errors: []error{nil},
	}

	repository := &MockStockRepository{}

	svc := sync.NewService(provider, repository, 2)

	// Act
	err := svc.Run()

	// Assert
	if err != nil {
		t.Fatalf("Run() no debería retornar error, pero obtuvo: %v", err)
	}

	if len(repository.UpsertedStocks) != len(stocks) {
		t.Errorf("Se esperaban %d stocks, pero se obtuvieron %d", len(stocks), len(repository.UpsertedStocks))
	}

	for i, stock := range stocks {
		if repository.UpsertedStocks[i].Ticker != stock.Ticker {
			t.Errorf("Stock %d: ticker esperado %s, obtuvo %s", i, stock.Ticker, repository.UpsertedStocks[i].Ticker)
		}
	}
}

func TestRun_Success_MultiplePages(t *testing.T) {
	// Arrange
	page2 := "page2"

	stocks1 := []stock.Stock{
		{
			Ticker:    "AAPL",
			Company:   "Apple",
			Action:    "BUY",
			Brokerage: "Goldman",
			Time:      time.Now(),
		},
	}

	stocks2 := []stock.Stock{
		{
			Ticker:    "GOOGL",
			Company:   "Google",
			Action:    "SELL",
			Brokerage: "Morgan",
			Time:      time.Now(),
		},
	}

	provider := &MockStockProvider{
		responses: []*stock.Page{
			{
				Items:    stocks1,
				NextPage: &page2,
			},
			{
				Items:    stocks2,
				NextPage: nil,
			},
		},
		errors: []error{nil, nil},
	}

	repository := &MockStockRepository{}

	svc := sync.NewService(provider, repository, 1)

	// Act
	err := svc.Run()

	// Assert
	if err != nil {
		t.Fatalf("Run() no debería retornar error, pero obtuvo: %v", err)
	}

	if len(repository.UpsertedStocks) != 2 {
		t.Errorf("Se esperaban 2 stocks, pero se obtuvieron %d", len(repository.UpsertedStocks))
	}

	if provider.callCount != 2 {
		t.Errorf("Provider debería haber sido llamado 2 veces, pero fue llamado %d veces", provider.callCount)
	}
}

func TestRun_ProviderError(t *testing.T) {
	// Arrange
	provider := &MockStockProvider{
		responses: []*stock.Page{nil},
		errors:    []error{errors.New("provider error")},
	}

	repository := &MockStockRepository{}

	svc := sync.NewService(provider, repository, 1)

	// Act
	err := svc.Run()

	// Assert
	if err == nil {
		t.Error("Run() debería retornar un error del provider")
	}

	if err.Error() != "provider error" {
		t.Errorf("Mensaje de error esperado 'provider error', obtuvo: %v", err)
	}

	if len(repository.UpsertedStocks) != 0 {
		t.Errorf("Repository no debería tener stocks, pero tiene %d", len(repository.UpsertedStocks))
	}
}

func TestRun_RepositoryError(t *testing.T) {
	// Arrange
	stocks := []stock.Stock{
		{
			Ticker:    "AAPL",
			Company:   "Apple",
			Action:    "BUY",
			Brokerage: "Goldman",
			Time:      time.Now(),
		},
	}

	provider := &MockStockProvider{
		responses: []*stock.Page{
			{
				Items:    stocks,
				NextPage: nil,
			},
		},
		errors: []error{nil},
	}

	repository := &MockStockRepository{
		UpsertError: errors.New("database error"),
	}

	svc := sync.NewService(provider, repository, 1)

	// Act
	err := svc.Run()

	// Assert
	// Los errores del repository se ignoran en los workers, así que Run() no retorna error
	if err != nil {
		t.Errorf("Run() no debería retornar error (errores del repository se ignoran): %v", err)
	}
}

func TestRun_MultipleWorkers(t *testing.T) {
	// Arrange
	stocks := []stock.Stock{
		{Ticker: "AAPL", Company: "Apple", Action: "BUY", Brokerage: "Goldman", Time: time.Now()},
		{Ticker: "GOOGL", Company: "Google", Action: "SELL", Brokerage: "Morgan", Time: time.Now()},
		{Ticker: "MSFT", Company: "Microsoft", Action: "HOLD", Brokerage: "JPM", Time: time.Now()},
		{Ticker: "AMZN", Company: "Amazon", Action: "BUY", Brokerage: "GS", Time: time.Now()},
		{Ticker: "TSLA", Company: "Tesla", Action: "SELL", Brokerage: "Citi", Time: time.Now()},
	}

	provider := &MockStockProvider{
		responses: []*stock.Page{
			{
				Items:    stocks,
				NextPage: nil,
			},
		},
		errors: []error{nil},
	}

	repository := &MockStockRepository{}

	svc := sync.NewService(provider, repository, 5)

	// Act
	err := svc.Run()

	// Assert
	if err != nil {
		t.Fatalf("Run() no debería retornar error, pero obtuvo: %v", err)
	}

	if len(repository.UpsertedStocks) != len(stocks) {
		t.Errorf("Se esperaban %d stocks procesados, pero se obtuvieron %d", len(stocks), len(repository.UpsertedStocks))
	}

	// Verificar que todos los stocks fueron procesados
	tickersProcessed := make(map[string]bool)
	for _, s := range repository.UpsertedStocks {
		tickersProcessed[s.Ticker] = true
	}

	for _, stock := range stocks {
		if !tickersProcessed[stock.Ticker] {
			t.Errorf("Stock con ticker %s no fue procesado", stock.Ticker)
		}
	}
}

func TestRun_EmptyPage(t *testing.T) {
	// Arrange
	provider := &MockStockProvider{
		responses: []*stock.Page{
			{
				Items:    []stock.Stock{},
				NextPage: nil,
			},
		},
		errors: []error{nil},
	}

	repository := &MockStockRepository{}

	svc := sync.NewService(provider, repository, 1)

	// Act
	err := svc.Run()

	// Assert
	if err != nil {
		t.Fatalf("Run() no debería retornar error para página vacía, pero obtuvo: %v", err)
	}

	if len(repository.UpsertedStocks) != 0 {
		t.Errorf("Se esperaban 0 stocks, pero se obtuvieron %d", len(repository.UpsertedStocks))
	}
}
