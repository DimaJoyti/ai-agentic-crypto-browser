package testing

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/shopspring/decimal"
)

// MockMarketDataProvider simulates market data for testing
type MockMarketDataProvider struct {
	logger *observability.Logger

	// Market data
	klines     map[string][]*Kline       // symbol -> klines
	tickers    map[string]*Ticker        // symbol -> ticker
	orderBooks map[string]*OrderBookData // symbol -> order book

	// Data generation
	generators map[string]*DataGenerator // symbol -> data generator

	// Configuration
	config *MarketDataConfig

	// Synchronization
	mu        sync.RWMutex
	isRunning bool
	stopChan  chan struct{}
}

// MarketDataConfig holds configuration for market data simulation
type MarketDataConfig struct {
	UpdateInterval   time.Duration   `json:"update_interval"`
	KlineInterval    string          `json:"kline_interval"`
	HistoryLength    int             `json:"history_length"`
	PriceVolatility  decimal.Decimal `json:"price_volatility"`
	VolumeVariation  decimal.Decimal `json:"volume_variation"`
	SupportedSymbols []string        `json:"supported_symbols"`
}

// Kline represents a candlestick/kline data point
type Kline struct {
	Symbol      string          `json:"symbol"`
	OpenTime    time.Time       `json:"open_time"`
	CloseTime   time.Time       `json:"close_time"`
	Open        decimal.Decimal `json:"open"`
	High        decimal.Decimal `json:"high"`
	Low         decimal.Decimal `json:"low"`
	Close       decimal.Decimal `json:"close"`
	Volume      decimal.Decimal `json:"volume"`
	QuoteVolume decimal.Decimal `json:"quote_volume"`
	TradeCount  int             `json:"trade_count"`
	Interval    string          `json:"interval"`
}

// Ticker represents current market ticker information
type Ticker struct {
	Symbol             string          `json:"symbol"`
	Price              decimal.Decimal `json:"price"`
	PriceChange        decimal.Decimal `json:"price_change"`
	PriceChangePercent decimal.Decimal `json:"price_change_percent"`
	High24h            decimal.Decimal `json:"high_24h"`
	Low24h             decimal.Decimal `json:"low_24h"`
	Volume24h          decimal.Decimal `json:"volume_24h"`
	QuoteVolume24h     decimal.Decimal `json:"quote_volume_24h"`
	Timestamp          time.Time       `json:"timestamp"`
}

// OrderBookData represents order book data
type OrderBookData struct {
	Symbol    string           `json:"symbol"`
	Bids      []OrderBookLevel `json:"bids"`
	Asks      []OrderBookLevel `json:"asks"`
	Timestamp time.Time        `json:"timestamp"`
}

// OrderBookLevel represents a price level in the order book
type OrderBookLevel struct {
	Price  decimal.Decimal `json:"price"`
	Amount decimal.Decimal `json:"amount"`
}

// DataGenerator generates realistic market data
type DataGenerator struct {
	Symbol          string          `json:"symbol"`
	BasePrice       decimal.Decimal `json:"base_price"`
	CurrentPrice    decimal.Decimal `json:"current_price"`
	Volatility      decimal.Decimal `json:"volatility"`
	Trend           decimal.Decimal `json:"trend"`
	Volume          decimal.Decimal `json:"volume"`
	LastUpdate      time.Time       `json:"last_update"`
	MarketCondition string          `json:"market_condition"`
}

// NewMockMarketDataProvider creates a new mock market data provider
func NewMockMarketDataProvider(logger *observability.Logger) *MockMarketDataProvider {
	config := &MarketDataConfig{
		UpdateInterval:   1 * time.Second,
		KlineInterval:    "1m",
		HistoryLength:    1000,
		PriceVolatility:  decimal.NewFromFloat(0.02), // 2%
		VolumeVariation:  decimal.NewFromFloat(0.30), // 30%
		SupportedSymbols: []string{"BTC/USDT", "ETH/USDT", "BNB/USDT", "ADA/USDT", "DOT/USDT"},
	}

	provider := &MockMarketDataProvider{
		logger:     logger,
		klines:     make(map[string][]*Kline),
		tickers:    make(map[string]*Ticker),
		orderBooks: make(map[string]*OrderBookData),
		generators: make(map[string]*DataGenerator),
		config:     config,
		stopChan:   make(chan struct{}),
	}

	// Initialize data generators
	provider.initializeGenerators()

	return provider
}

// Start starts the market data provider
func (mdp *MockMarketDataProvider) Start(ctx context.Context) error {
	mdp.mu.Lock()
	defer mdp.mu.Unlock()

	if mdp.isRunning {
		return fmt.Errorf("market data provider is already running")
	}

	mdp.isRunning = true

	// Generate initial historical data
	mdp.generateHistoricalData()

	// Start data generation loop
	go mdp.dataGenerationLoop(ctx)

	mdp.logger.Info(ctx, "Mock market data provider started", map[string]interface{}{
		"supported_symbols": len(mdp.config.SupportedSymbols),
		"update_interval":   mdp.config.UpdateInterval.String(),
		"history_length":    mdp.config.HistoryLength,
	})

	return nil
}

// Stop stops the market data provider
func (mdp *MockMarketDataProvider) Stop(ctx context.Context) error {
	mdp.mu.Lock()
	defer mdp.mu.Unlock()

	if !mdp.isRunning {
		return nil
	}

	mdp.isRunning = false
	close(mdp.stopChan)

	mdp.logger.Info(ctx, "Mock market data provider stopped", nil)
	return nil
}

// GetKlines returns historical kline data
func (mdp *MockMarketDataProvider) GetKlines(symbol, interval string, limit int, startTime, endTime *time.Time) ([]*Kline, error) {
	mdp.mu.RLock()
	defer mdp.mu.RUnlock()

	klines, exists := mdp.klines[symbol]
	if !exists {
		return nil, fmt.Errorf("klines not available for symbol: %s", symbol)
	}

	// Filter by time range if specified
	filteredKlines := make([]*Kline, 0)
	for _, kline := range klines {
		if startTime != nil && kline.OpenTime.Before(*startTime) {
			continue
		}
		if endTime != nil && kline.OpenTime.After(*endTime) {
			continue
		}
		filteredKlines = append(filteredKlines, kline)
	}

	// Apply limit
	if limit > 0 && len(filteredKlines) > limit {
		start := len(filteredKlines) - limit
		filteredKlines = filteredKlines[start:]
	}

	return filteredKlines, nil
}

// GetTicker returns current ticker information
func (mdp *MockMarketDataProvider) GetTicker(symbol string) (*Ticker, error) {
	mdp.mu.RLock()
	defer mdp.mu.RUnlock()

	ticker, exists := mdp.tickers[symbol]
	if !exists {
		return nil, fmt.Errorf("ticker not available for symbol: %s", symbol)
	}

	return ticker, nil
}

// GetAllTickers returns all ticker information
func (mdp *MockMarketDataProvider) GetAllTickers() (map[string]*Ticker, error) {
	mdp.mu.RLock()
	defer mdp.mu.RUnlock()

	// Return a copy to avoid concurrent access issues
	result := make(map[string]*Ticker)
	for symbol, ticker := range mdp.tickers {
		result[symbol] = ticker
	}

	return result, nil
}

// GetOrderBook returns order book data
func (mdp *MockMarketDataProvider) GetOrderBook(symbol string, limit int) (*OrderBookData, error) {
	mdp.mu.RLock()
	defer mdp.mu.RUnlock()

	orderBook, exists := mdp.orderBooks[symbol]
	if !exists {
		return nil, fmt.Errorf("order book not available for symbol: %s", symbol)
	}

	// Apply limit if specified
	if limit > 0 {
		limitedOrderBook := &OrderBookData{
			Symbol:    orderBook.Symbol,
			Timestamp: orderBook.Timestamp,
			Bids:      orderBook.Bids,
			Asks:      orderBook.Asks,
		}

		if len(limitedOrderBook.Bids) > limit {
			limitedOrderBook.Bids = limitedOrderBook.Bids[:limit]
		}
		if len(limitedOrderBook.Asks) > limit {
			limitedOrderBook.Asks = limitedOrderBook.Asks[:limit]
		}

		return limitedOrderBook, nil
	}

	return orderBook, nil
}

// SetMarketCondition sets the market condition for data generation
func (mdp *MockMarketDataProvider) SetMarketCondition(condition string) {
	mdp.mu.Lock()
	defer mdp.mu.Unlock()

	for _, generator := range mdp.generators {
		generator.MarketCondition = condition

		switch condition {
		case "bull":
			generator.Trend = decimal.NewFromFloat(0.001) // 0.1% upward trend
			generator.Volatility = mdp.config.PriceVolatility
		case "bear":
			generator.Trend = decimal.NewFromFloat(-0.001) // 0.1% downward trend
			generator.Volatility = mdp.config.PriceVolatility
		case "sideways":
			generator.Trend = decimal.Zero
			generator.Volatility = mdp.config.PriceVolatility.Div(decimal.NewFromFloat(2))
		case "volatile":
			generator.Trend = decimal.Zero
			generator.Volatility = mdp.config.PriceVolatility.Mul(decimal.NewFromFloat(3))
		default:
			generator.Trend = decimal.Zero
			generator.Volatility = mdp.config.PriceVolatility
		}
	}

	mdp.logger.Info(context.Background(), "Market condition set", map[string]interface{}{
		"condition": condition,
	})
}

// initializeGenerators sets up data generators for all symbols
func (mdp *MockMarketDataProvider) initializeGenerators() {
	basePrices := map[string]decimal.Decimal{
		"BTC/USDT": decimal.NewFromFloat(50000),
		"ETH/USDT": decimal.NewFromFloat(3000),
		"BNB/USDT": decimal.NewFromFloat(400),
		"ADA/USDT": decimal.NewFromFloat(1.5),
		"DOT/USDT": decimal.NewFromFloat(25),
	}

	for _, symbol := range mdp.config.SupportedSymbols {
		basePrice := basePrices[symbol]
		if basePrice.IsZero() {
			basePrice = decimal.NewFromFloat(100) // Default price
		}

		mdp.generators[symbol] = &DataGenerator{
			Symbol:          symbol,
			BasePrice:       basePrice,
			CurrentPrice:    basePrice,
			Volatility:      mdp.config.PriceVolatility,
			Trend:           decimal.Zero,
			Volume:          decimal.NewFromFloat(1000),
			LastUpdate:      time.Now(),
			MarketCondition: "normal",
		}

		// Initialize empty data structures
		mdp.klines[symbol] = make([]*Kline, 0)
		mdp.tickers[symbol] = &Ticker{
			Symbol:    symbol,
			Price:     basePrice,
			Timestamp: time.Now(),
		}
		mdp.orderBooks[symbol] = mdp.generateOrderBook(symbol, basePrice)
	}
}

// generateHistoricalData generates initial historical data
func (mdp *MockMarketDataProvider) generateHistoricalData() {
	now := time.Now()

	for symbol, generator := range mdp.generators {
		klines := make([]*Kline, 0, mdp.config.HistoryLength)

		currentPrice := generator.BasePrice
		currentTime := now.Add(-time.Duration(mdp.config.HistoryLength) * time.Minute)

		for i := 0; i < mdp.config.HistoryLength; i++ {
			// Generate realistic OHLCV data
			open := currentPrice

			// Generate price movement
			change := mdp.generatePriceChange(generator)
			close := open.Add(change)

			// Ensure close price is positive
			if close.LessThan(decimal.NewFromFloat(0.01)) {
				close = decimal.NewFromFloat(0.01)
			}

			// Generate high and low
			high := decimal.Max(open, close).Add(open.Mul(decimal.NewFromFloat(rand.Float64() * 0.01)))
			low := decimal.Min(open, close).Sub(open.Mul(decimal.NewFromFloat(rand.Float64() * 0.01)))

			// Ensure low is positive
			if low.LessThan(decimal.NewFromFloat(0.01)) {
				low = decimal.NewFromFloat(0.01)
			}

			// Generate volume
			baseVolume := generator.Volume
			volumeVariation := baseVolume.Mul(mdp.config.VolumeVariation).Mul(decimal.NewFromFloat(rand.Float64() - 0.5))
			volume := baseVolume.Add(volumeVariation)
			if volume.LessThan(decimal.Zero) {
				volume = baseVolume.Div(decimal.NewFromFloat(2))
			}

			kline := &Kline{
				Symbol:      symbol,
				OpenTime:    currentTime,
				CloseTime:   currentTime.Add(time.Minute),
				Open:        open,
				High:        high,
				Low:         low,
				Close:       close,
				Volume:      volume,
				QuoteVolume: volume.Mul(close),
				TradeCount:  int(rand.Float64()*100) + 10,
				Interval:    mdp.config.KlineInterval,
			}

			klines = append(klines, kline)

			// Update for next iteration
			currentPrice = close
			currentTime = currentTime.Add(time.Minute)
		}

		mdp.klines[symbol] = klines

		// Update generator current price
		generator.CurrentPrice = currentPrice

		// Update ticker with latest data
		mdp.updateTicker(symbol, klines[len(klines)-1])
	}
}

// dataGenerationLoop continuously generates new market data
func (mdp *MockMarketDataProvider) dataGenerationLoop(ctx context.Context) {
	ticker := time.NewTicker(mdp.config.UpdateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-mdp.stopChan:
			return
		case <-ticker.C:
			mdp.generateNewData()
		}
	}
}

// generateNewData generates new market data for all symbols
func (mdp *MockMarketDataProvider) generateNewData() {
	mdp.mu.Lock()
	defer mdp.mu.Unlock()

	now := time.Now()

	for symbol, generator := range mdp.generators {
		// Generate new kline
		lastKline := mdp.getLastKline(symbol)
		newKline := mdp.generateNextKline(generator, lastKline, now)

		// Add to klines
		mdp.klines[symbol] = append(mdp.klines[symbol], newKline)

		// Maintain history length
		if len(mdp.klines[symbol]) > mdp.config.HistoryLength {
			mdp.klines[symbol] = mdp.klines[symbol][1:]
		}

		// Update ticker
		mdp.updateTicker(symbol, newKline)

		// Update order book
		mdp.orderBooks[symbol] = mdp.generateOrderBook(symbol, newKline.Close)

		// Update generator
		generator.CurrentPrice = newKline.Close
		generator.LastUpdate = now
	}
}

// generatePriceChange generates a realistic price change
func (mdp *MockMarketDataProvider) generatePriceChange(generator *DataGenerator) decimal.Decimal {
	// Random walk with trend and volatility
	randomComponent := decimal.NewFromFloat((rand.Float64() - 0.5) * 2) // -1 to 1
	volatilityComponent := randomComponent.Mul(generator.Volatility).Mul(generator.CurrentPrice)
	trendComponent := generator.Trend.Mul(generator.CurrentPrice)

	return volatilityComponent.Add(trendComponent)
}

// generateNextKline generates the next kline based on current state
func (mdp *MockMarketDataProvider) generateNextKline(generator *DataGenerator, lastKline *Kline, timestamp time.Time) *Kline {
	open := generator.CurrentPrice
	if lastKline != nil {
		open = lastKline.Close
	}

	// Generate price movement
	change := mdp.generatePriceChange(generator)
	close := open.Add(change)

	// Ensure close price is positive
	if close.LessThan(decimal.NewFromFloat(0.01)) {
		close = decimal.NewFromFloat(0.01)
	}

	// Generate high and low with some randomness
	highVariation := open.Mul(decimal.NewFromFloat(rand.Float64() * 0.01))
	lowVariation := open.Mul(decimal.NewFromFloat(rand.Float64() * 0.01))

	high := decimal.Max(open, close).Add(highVariation)
	low := decimal.Min(open, close).Sub(lowVariation)

	// Ensure low is positive
	if low.LessThan(decimal.NewFromFloat(0.01)) {
		low = decimal.NewFromFloat(0.01)
	}

	// Generate volume with variation
	baseVolume := generator.Volume
	volumeVariation := baseVolume.Mul(mdp.config.VolumeVariation).Mul(decimal.NewFromFloat(rand.Float64() - 0.5))
	volume := baseVolume.Add(volumeVariation)
	if volume.LessThan(decimal.Zero) {
		volume = baseVolume.Div(decimal.NewFromFloat(2))
	}

	return &Kline{
		Symbol:      generator.Symbol,
		OpenTime:    timestamp,
		CloseTime:   timestamp.Add(time.Minute),
		Open:        open,
		High:        high,
		Low:         low,
		Close:       close,
		Volume:      volume,
		QuoteVolume: volume.Mul(close),
		TradeCount:  int(rand.Float64()*100) + 10,
		Interval:    mdp.config.KlineInterval,
	}
}

// getLastKline returns the last kline for a symbol
func (mdp *MockMarketDataProvider) getLastKline(symbol string) *Kline {
	klines := mdp.klines[symbol]
	if len(klines) == 0 {
		return nil
	}
	return klines[len(klines)-1]
}

// updateTicker updates ticker information based on latest kline
func (mdp *MockMarketDataProvider) updateTicker(symbol string, kline *Kline) {
	ticker := mdp.tickers[symbol]

	// Calculate 24h change (simplified)
	priceChange := decimal.Zero
	priceChangePercent := decimal.Zero

	klines := mdp.klines[symbol]
	if len(klines) > 24*60 { // 24 hours of minute data
		oldPrice := klines[len(klines)-24*60].Close
		priceChange = kline.Close.Sub(oldPrice)
		if !oldPrice.IsZero() {
			priceChangePercent = priceChange.Div(oldPrice).Mul(decimal.NewFromFloat(100))
		}
	}

	// Calculate 24h high/low (simplified)
	high24h := kline.High
	low24h := kline.Low
	volume24h := decimal.Zero

	if len(klines) > 24*60 {
		for i := len(klines) - 24*60; i < len(klines); i++ {
			if klines[i].High.GreaterThan(high24h) {
				high24h = klines[i].High
			}
			if klines[i].Low.LessThan(low24h) {
				low24h = klines[i].Low
			}
			volume24h = volume24h.Add(klines[i].Volume)
		}
	}

	ticker.Price = kline.Close
	ticker.PriceChange = priceChange
	ticker.PriceChangePercent = priceChangePercent
	ticker.High24h = high24h
	ticker.Low24h = low24h
	ticker.Volume24h = volume24h
	ticker.QuoteVolume24h = volume24h.Mul(kline.Close)
	ticker.Timestamp = kline.CloseTime
}

// generateOrderBook generates a realistic order book
func (mdp *MockMarketDataProvider) generateOrderBook(symbol string, currentPrice decimal.Decimal) *OrderBookData {
	orderBook := &OrderBookData{
		Symbol:    symbol,
		Bids:      make([]OrderBookLevel, 0),
		Asks:      make([]OrderBookLevel, 0),
		Timestamp: time.Now(),
	}

	// Generate bids (buy orders)
	for i := 0; i < 20; i++ {
		priceOffset := decimal.NewFromFloat(float64(i+1) * 0.001) // 0.1% increments
		price := currentPrice.Sub(currentPrice.Mul(priceOffset))
		amount := decimal.NewFromFloat(rand.Float64() * 10)

		orderBook.Bids = append(orderBook.Bids, OrderBookLevel{
			Price:  price,
			Amount: amount,
		})
	}

	// Generate asks (sell orders)
	for i := 0; i < 20; i++ {
		priceOffset := decimal.NewFromFloat(float64(i+1) * 0.001) // 0.1% increments
		price := currentPrice.Add(currentPrice.Mul(priceOffset))
		amount := decimal.NewFromFloat(rand.Float64() * 10)

		orderBook.Asks = append(orderBook.Asks, OrderBookLevel{
			Price:  price,
			Amount: amount,
		})
	}

	return orderBook
}

// GetMarketDataInfo returns market data provider information
func (mdp *MockMarketDataProvider) GetMarketDataInfo() map[string]interface{} {
	mdp.mu.RLock()
	defer mdp.mu.RUnlock()

	return map[string]interface{}{
		"name":              "MockMarketDataProvider",
		"supported_symbols": mdp.config.SupportedSymbols,
		"update_interval":   mdp.config.UpdateInterval.String(),
		"kline_interval":    mdp.config.KlineInterval,
		"history_length":    mdp.config.HistoryLength,
		"is_running":        mdp.isRunning,
		"total_symbols":     len(mdp.generators),
	}
}
