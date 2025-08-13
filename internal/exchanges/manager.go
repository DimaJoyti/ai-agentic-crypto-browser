package exchanges

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/internal/exchanges/binance"
	"github.com/ai-agentic-browser/internal/exchanges/common"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/shopspring/decimal"
)

// Manager manages multiple exchange connections and provides unified access
type Manager struct {
	logger    *observability.Logger
	exchanges map[string]common.ExchangeClient
	config    Config

	// State management
	isRunning bool
	stopChan  chan struct{}
	wg        sync.WaitGroup
	mu        sync.RWMutex
}

// Config contains configuration for the exchange manager
type Config struct {
	Binance  binance.Config `json:"binance"`
	// Future exchanges will be added here
	// Coinbase coinbase.Config `json:"coinbase"`
	// Kraken   kraken.Config   `json:"kraken"`
	
	DefaultExchange string `json:"default_exchange"`
	EnabledExchanges []string `json:"enabled_exchanges"`
}

// NewManager creates a new exchange manager
func NewManager(logger *observability.Logger, config Config) *Manager {
	return &Manager{
		logger:    logger,
		exchanges: make(map[string]common.ExchangeClient),
		config:    config,
		stopChan:  make(chan struct{}),
	}
}

// Start initializes and starts all configured exchanges
func (m *Manager) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.isRunning {
		return nil
	}

	m.logger.Info(ctx, "Starting exchange manager", map[string]interface{}{
		"enabled_exchanges": m.config.EnabledExchanges,
		"default_exchange":  m.config.DefaultExchange,
	})

	// Initialize enabled exchanges
	for _, exchangeName := range m.config.EnabledExchanges {
		if err := m.initializeExchange(ctx, exchangeName); err != nil {
			m.logger.Error(ctx, "Failed to initialize exchange", err, map[string]interface{}{
				"exchange": exchangeName,
			})
			continue
		}
	}

	m.isRunning = true

	m.logger.Info(ctx, "Exchange manager started successfully", map[string]interface{}{
		"active_exchanges": len(m.exchanges),
	})

	return nil
}

// Stop gracefully shuts down all exchanges
func (m *Manager) Stop(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.isRunning {
		return nil
	}

	m.logger.Info(ctx, "Stopping exchange manager", nil)

	close(m.stopChan)

	// Disconnect all exchanges
	for name, exchange := range m.exchanges {
		if err := exchange.Disconnect(ctx); err != nil {
			m.logger.Error(ctx, "Failed to disconnect exchange", err, map[string]interface{}{
				"exchange": name,
			})
		}
	}

	m.wg.Wait()
	m.isRunning = false

	m.logger.Info(ctx, "Exchange manager stopped", nil)

	return nil
}

// initializeExchange initializes a specific exchange
func (m *Manager) initializeExchange(ctx context.Context, exchangeName string) error {
	switch exchangeName {
	case "binance":
		client := binance.NewClient(m.logger, m.config.Binance)
		if err := client.Connect(ctx); err != nil {
			return fmt.Errorf("failed to connect to Binance: %w", err)
		}
		m.exchanges["binance"] = client

	// Future exchanges will be added here
	// case "coinbase":
	//     client := coinbase.NewClient(m.logger, m.config.Coinbase)
	//     if err := client.Connect(ctx); err != nil {
	//         return fmt.Errorf("failed to connect to Coinbase: %w", err)
	//     }
	//     m.exchanges["coinbase"] = client

	// case "kraken":
	//     client := kraken.NewClient(m.logger, m.config.Kraken)
	//     if err := client.Connect(ctx); err != nil {
	//         return fmt.Errorf("failed to connect to Kraken: %w", err)
	//     }
	//     m.exchanges["kraken"] = client

	default:
		return fmt.Errorf("unsupported exchange: %s", exchangeName)
	}

	m.logger.Info(ctx, "Exchange initialized successfully", map[string]interface{}{
		"exchange": exchangeName,
	})

	return nil
}

// GetExchange returns a specific exchange client
func (m *Manager) GetExchange(name string) (common.ExchangeClient, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if name == "" {
		name = m.config.DefaultExchange
	}

	exchange, exists := m.exchanges[name]
	if !exists {
		return nil, fmt.Errorf("exchange not found: %s", name)
	}

	return exchange, nil
}

// GetDefaultExchange returns the default exchange client
func (m *Manager) GetDefaultExchange() (common.ExchangeClient, error) {
	return m.GetExchange(m.config.DefaultExchange)
}

// GetAllExchanges returns all active exchange clients
func (m *Manager) GetAllExchanges() map[string]common.ExchangeClient {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy to avoid race conditions
	exchanges := make(map[string]common.ExchangeClient)
	for name, exchange := range m.exchanges {
		exchanges[name] = exchange
	}

	return exchanges
}

// GetActiveExchanges returns names of all active exchanges
func (m *Manager) GetActiveExchanges() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.exchanges))
	for name := range m.exchanges {
		names = append(names, name)
	}

	return names
}

// Unified Market Data Methods

// GetBestPrice gets the best price across all exchanges for a symbol
func (m *Manager) GetBestPrice(ctx context.Context, symbol string, side common.OrderSide) (*BestPriceResult, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var bestPrice decimal.Decimal
	var bestExchange string
	var bestTicker *common.TickerData

	for name, exchange := range m.exchanges {
		ticker, err := exchange.GetTicker(ctx, symbol)
		if err != nil {
			m.logger.Error(ctx, "Failed to get ticker", err, map[string]interface{}{
				"exchange": name,
				"symbol":   symbol,
			})
			continue
		}

		var price decimal.Decimal
		if side == common.OrderSideBuy {
			price = ticker.AskPrice // Best ask for buying
		} else {
			price = ticker.BidPrice // Best bid for selling
		}

		if bestPrice.IsZero() || 
			(side == common.OrderSideBuy && price.LessThan(bestPrice)) ||
			(side == common.OrderSideSell && price.GreaterThan(bestPrice)) {
			bestPrice = price
			bestExchange = name
			bestTicker = ticker
		}
	}

	if bestPrice.IsZero() {
		return nil, fmt.Errorf("no price found for symbol %s", symbol)
	}

	return &BestPriceResult{
		Price:    bestPrice,
		Exchange: bestExchange,
		Ticker:   bestTicker,
		Side:     side,
		Symbol:   symbol,
	}, nil
}

// GetAggregatedOrderBook gets aggregated order book across all exchanges
func (m *Manager) GetAggregatedOrderBook(ctx context.Context, symbol string, limit int) (*AggregatedOrderBook, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var allBids []common.PriceLevel
	var allAsks []common.PriceLevel
	exchangeBooks := make(map[string]*common.OrderBookData)

	for name, exchange := range m.exchanges {
		orderBook, err := exchange.GetOrderBook(ctx, symbol, limit)
		if err != nil {
			m.logger.Error(ctx, "Failed to get order book", err, map[string]interface{}{
				"exchange": name,
				"symbol":   symbol,
			})
			continue
		}

		exchangeBooks[name] = orderBook
		allBids = append(allBids, orderBook.Bids...)
		allAsks = append(allAsks, orderBook.Asks...)
	}

	// Sort bids (highest first) and asks (lowest first)
	// This would need proper sorting implementation

	return &AggregatedOrderBook{
		Symbol:        symbol,
		Bids:          allBids,
		Asks:          allAsks,
		ExchangeBooks: exchangeBooks,
		Timestamp:     time.Now(),
	}, nil
}

// Smart Order Routing

// ExecuteSmartOrder executes an order using smart routing across exchanges
func (m *Manager) ExecuteSmartOrder(ctx context.Context, order *common.OrderRequest) ([]*common.OrderResponse, error) {
	// For now, route to the exchange with the best price
	bestPrice, err := m.GetBestPrice(ctx, order.Symbol, order.Side)
	if err != nil {
		return nil, fmt.Errorf("failed to get best price: %w", err)
	}

	exchange, err := m.GetExchange(bestPrice.Exchange)
	if err != nil {
		return nil, fmt.Errorf("failed to get exchange: %w", err)
	}

	response, err := exchange.PlaceOrder(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to place order: %w", err)
	}

	return []*common.OrderResponse{response}, nil
}

// Result types

// BestPriceResult represents the best price result across exchanges
type BestPriceResult struct {
	Price    decimal.Decimal      `json:"price"`
	Exchange string               `json:"exchange"`
	Ticker   *common.TickerData   `json:"ticker"`
	Side     common.OrderSide     `json:"side"`
	Symbol   string               `json:"symbol"`
}

// AggregatedOrderBook represents aggregated order book data
type AggregatedOrderBook struct {
	Symbol        string                             `json:"symbol"`
	Bids          []common.PriceLevel               `json:"bids"`
	Asks          []common.PriceLevel               `json:"asks"`
	ExchangeBooks map[string]*common.OrderBookData  `json:"exchange_books"`
	Timestamp     time.Time                         `json:"timestamp"`
}
