package sync_test

import (
	"errors"
	syncsvc "go-vue-journey/internal/integrations/sync"
	"go-vue-journey/internal/stock"
	"sync"
	"testing"
	"time"
)

// MockStockProvider es un mock del StockProvider para testing
type MockStockProvider struct {
	//Cuantas veces se llama a la función GetStocks
	callCount int
	//Que páginas debe devolver con cada llamada
	responses []*stock.Page
	//Que error debe de devolver con cada llamada
	errors []error
}

// Método que reemplaza al método real que se requiere probar
func (m *MockStockProvider) GetStocks(page *string) (*stock.Page, error) {
	//Esto se hace con el fin de frenar llamadas inesperasas, si el test se llama más de lo esperado
	if m.callCount >= len(m.responses) {
		//Termina la ejecucion y falla el test
		return nil, errors.New("unexpected call to GetStocks")
	}

	/*

		Esto es un retorno controlado
		en la primera ejecución la funcion envia
		m.errors[0]
		m.response[0]
		...
		asi sucesivamente

	*/
	err := m.errors[m.callCount]
	response := m.responses[m.callCount]
	//añade 1 al contador de llamadas de la función
	m.callCount++

	return response, err
}

// MockStockRepository es un mock del StockRepository para testing
type MockStockRepository struct {
	//Guarda los stocks que se procesaron
	UpsertedStocks []stock.Stock
	//Simula un error en la base de datos
	UpsertError error
	mu          sync.Mutex
}

// Método que prueba el método original del repositorio de guardar archivos
func (m *MockStockRepository) Upsert(s stock.Stock) error {
	//Si la db falla entonces envia el error
	if m.UpsertError != nil {
		return m.UpsertError
	}

	//Si la db envia un ok entonces guarda el dato
	m.mu.Lock()
	defer m.mu.Unlock()
	m.UpsertedStocks = append(m.UpsertedStocks, s)
	return nil
}

/*
Estos métodos no se testea aquí debido a que este test esta
diseñado principalmente para la sincronozación de la base de datos
con la API externa
*/

func (m *MockStockRepository) GetStocks(limit int, cursorTicker *string, filter stock.StockFilter) ([]stock.Stock, error) {
	return nil, nil
}

func (m *MockStockRepository) GetStocksStats() (stock.StocksStats, error) {
	return stock.StocksStats{}, nil
}

func (m *MockStockRepository) GetTopStocks(n int) ([]stock.Stock, error) {
	return nil, nil
}

// Este test prueba principalmente la creación de los nuevos servicions
func TestNewService(t *testing.T) {
	provider := &MockStockProvider{}
	repository := &MockStockRepository{}
	workers := 5

	svc := syncsvc.NewService(provider, repository, workers)

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

	svc := syncsvc.NewService(provider, repository, 2)

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

	svc := syncsvc.NewService(provider, repository, 1)

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

	svc := syncsvc.NewService(provider, repository, 1)

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

	svc := syncsvc.NewService(provider, repository, 1)

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

	svc := syncsvc.NewService(provider, repository, 5)

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

	svc := syncsvc.NewService(provider, repository, 1)

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
