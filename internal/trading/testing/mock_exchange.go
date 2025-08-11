package testing

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// MockExchange simulates a cryptocurrency exchange for testing
type MockExchange struct {
	logger *observability.Logger

	// Exchange state
	balances   map[string]map[string]decimal.Decimal // userID -> asset -> balance
	orders     map[string]*MockOrder                 // orderID -> order
	trades     map[string][]*MockTrade               // userID -> trades
	orderBooks map[string]*MockOrderBook             // symbol -> order book

	// Market simulation
	prices     map[string]decimal.Decimal // symbol -> current price
	priceFeeds map[string]*PriceFeed      // symbol -> price feed

	// Exchange configuration
	config *MockExchangeConfig

	// Synchronization
	mu        sync.RWMutex
	isRunning bool
	stopChan  chan struct{}
}

// MockExchangeConfig holds configuration for the mock exchange
type MockExchangeConfig struct {
	// Trading fees
	MakerFee decimal.Decimal `json:"maker_fee"`
	TakerFee decimal.Decimal `json:"taker_fee"`

	// Order execution
	ExecutionDelay time.Duration   `json:"execution_delay"`
	SlippageRate   decimal.Decimal `json:"slippage_rate"`
	FailureRate    float64         `json:"failure_rate"`

	// Market simulation
	PriceVolatility decimal.Decimal `json:"price_volatility"`
	UpdateInterval  time.Duration   `json:"update_interval"`

	// Supported assets
	SupportedPairs []string `json:"supported_pairs"`
	BaseCurrency   string   `json:"base_currency"`
}

// MockOrder represents an order in the mock exchange
type MockOrder struct {
	ID              string          `json:"id"`
	UserID          string          `json:"user_id"`
	Symbol          string          `json:"symbol"`
	Side            string          `json:"side"` // "buy" or "sell"
	Type            string          `json:"type"` // "market", "limit", "stop"
	Amount          decimal.Decimal `json:"amount"`
	Price           decimal.Decimal `json:"price"`
	FilledAmount    decimal.Decimal `json:"filled_amount"`
	RemainingAmount decimal.Decimal `json:"remaining_amount"`
	Status          string          `json:"status"` // "pending", "filled", "cancelled", "failed"
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	ExecutedAt      *time.Time      `json:"executed_at,omitempty"`
}

// MockTrade represents a completed trade
type MockTrade struct {
	ID        string          `json:"id"`
	OrderID   string          `json:"order_id"`
	UserID    string          `json:"user_id"`
	Symbol    string          `json:"symbol"`
	Side      string          `json:"side"`
	Amount    decimal.Decimal `json:"amount"`
	Price     decimal.Decimal `json:"price"`
	Fee       decimal.Decimal `json:"fee"`
	Timestamp time.Time       `json:"timestamp"`
}

// MockOrderBook represents an order book for a trading pair
type MockOrderBook struct {
	Symbol    string            `json:"symbol"`
	Bids      []*OrderBookEntry `json:"bids"`
	Asks      []*OrderBookEntry `json:"asks"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// OrderBookEntry represents an entry in the order book
type OrderBookEntry struct {
	Price  decimal.Decimal `json:"price"`
	Amount decimal.Decimal `json:"amount"`
}

// PriceFeed simulates price movements for a trading pair
type PriceFeed struct {
	Symbol       string          `json:"symbol"`
	CurrentPrice decimal.Decimal `json:"current_price"`
	BasePrice    decimal.Decimal `json:"base_price"`
	Volatility   decimal.Decimal `json:"volatility"`
	Trend        decimal.Decimal `json:"trend"`
	LastUpdate   time.Time       `json:"last_update"`
	PriceHistory []PricePoint    `json:"price_history"`
}

// PricePoint represents a price at a specific time
type PricePoint struct {
	Price     decimal.Decimal `json:"price"`
	Volume    decimal.Decimal `json:"volume"`
	Timestamp time.Time       `json:"timestamp"`
}

// NewMockExchange creates a new mock exchange
func NewMockExchange(logger *observability.Logger) *MockExchange {
	config := &MockExchangeConfig{
		MakerFee:        decimal.NewFromFloat(0.001), // 0.1%
		TakerFee:        decimal.NewFromFloat(0.001), // 0.1%
		ExecutionDelay:  100 * time.Millisecond,
		SlippageRate:    decimal.NewFromFloat(0.0005), // 0.05%
		FailureRate:     0.02,                         // 2% failure rate
		PriceVolatility: decimal.NewFromFloat(0.02),   // 2% volatility
		UpdateInterval:  1 * time.Second,
		SupportedPairs:  []string{"BTC/USDT", "ETH/USDT", "BNB/USDT", "ADA/USDT", "DOT/USDT"},
		BaseCurrency:    "USDT",
	}

	exchange := &MockExchange{
		logger:     logger,
		balances:   make(map[string]map[string]decimal.Decimal),
		orders:     make(map[string]*MockOrder),
		trades:     make(map[string][]*MockTrade),
		orderBooks: make(map[string]*MockOrderBook),
		prices:     make(map[string]decimal.Decimal),
		priceFeeds: make(map[string]*PriceFeed),
		config:     config,
		stopChan:   make(chan struct{}),
	}

	// Initialize price feeds
	exchange.initializePriceFeeds()

	return exchange
}

// Start starts the mock exchange
func (me *MockExchange) Start(ctx context.Context) error {
	me.mu.Lock()
	defer me.mu.Unlock()

	if me.isRunning {
		return fmt.Errorf("mock exchange is already running")
	}

	me.isRunning = true

	// Start price simulation
	go me.priceSimulationLoop(ctx)

	me.logger.Info(ctx, "Mock exchange started", map[string]interface{}{
		"supported_pairs": len(me.config.SupportedPairs),
		"maker_fee":       me.config.MakerFee.String(),
		"taker_fee":       me.config.TakerFee.String(),
	})

	return nil
}

// Stop stops the mock exchange
func (me *MockExchange) Stop(ctx context.Context) error {
	me.mu.Lock()
	defer me.mu.Unlock()

	if !me.isRunning {
		return nil
	}

	me.isRunning = false
	close(me.stopChan)

	me.logger.Info(ctx, "Mock exchange stopped", nil)
	return nil
}

// CreateAccount creates a new trading account with initial balances
func (me *MockExchange) CreateAccount(userID string, initialBalances map[string]decimal.Decimal) error {
	me.mu.Lock()
	defer me.mu.Unlock()

	if me.balances[userID] != nil {
		return fmt.Errorf("account already exists for user %s", userID)
	}

	me.balances[userID] = make(map[string]decimal.Decimal)
	for asset, balance := range initialBalances {
		me.balances[userID][asset] = balance
	}

	me.trades[userID] = make([]*MockTrade, 0)

	me.logger.Info(context.Background(), "Account created", map[string]interface{}{
		"user_id":  userID,
		"balances": initialBalances,
	})

	return nil
}

// PlaceOrder places a new order
func (me *MockExchange) PlaceOrder(ctx context.Context, userID, symbol, side, orderType string, amount, price decimal.Decimal) (*MockOrder, error) {
	me.mu.Lock()
	defer me.mu.Unlock()

	// Validate user account
	if me.balances[userID] == nil {
		return nil, fmt.Errorf("account not found for user %s", userID)
	}

	// Create order
	order := &MockOrder{
		ID:              uuid.New().String(),
		UserID:          userID,
		Symbol:          symbol,
		Side:            side,
		Type:            orderType,
		Amount:          amount,
		Price:           price,
		FilledAmount:    decimal.Zero,
		RemainingAmount: amount,
		Status:          "pending",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Store order
	me.orders[order.ID] = order

	// Simulate order execution
	go me.executeOrder(ctx, order)

	me.logger.Info(ctx, "Order placed", map[string]interface{}{
		"order_id": order.ID,
		"user_id":  userID,
		"symbol":   symbol,
		"side":     side,
		"amount":   amount.String(),
		"price":    price.String(),
	})

	return order, nil
}

// GetOrder retrieves an order by ID
func (me *MockExchange) GetOrder(orderID string) (*MockOrder, error) {
	me.mu.RLock()
	defer me.mu.RUnlock()

	order, exists := me.orders[orderID]
	if !exists {
		return nil, fmt.Errorf("order not found: %s", orderID)
	}

	return order, nil
}

// CancelOrder cancels an order
func (me *MockExchange) CancelOrder(ctx context.Context, orderID string) error {
	me.mu.Lock()
	defer me.mu.Unlock()

	order, exists := me.orders[orderID]
	if !exists {
		return fmt.Errorf("order not found: %s", orderID)
	}

	if order.Status != "pending" {
		return fmt.Errorf("cannot cancel order with status: %s", order.Status)
	}

	order.Status = "cancelled"
	order.UpdatedAt = time.Now()

	me.logger.Info(ctx, "Order cancelled", map[string]interface{}{
		"order_id": orderID,
		"user_id":  order.UserID,
	})

	return nil
}

// GetBalance returns the balance for a specific asset
func (me *MockExchange) GetBalance(userID, asset string) (decimal.Decimal, error) {
	me.mu.RLock()
	defer me.mu.RUnlock()

	userBalances, exists := me.balances[userID]
	if !exists {
		return decimal.Zero, fmt.Errorf("account not found for user %s", userID)
	}

	balance, exists := userBalances[asset]
	if !exists {
		return decimal.Zero, nil
	}

	return balance, nil
}

// GetAllBalances returns all balances for a user
func (me *MockExchange) GetAllBalances(userID string) (map[string]decimal.Decimal, error) {
	me.mu.RLock()
	defer me.mu.RUnlock()

	userBalances, exists := me.balances[userID]
	if !exists {
		return nil, fmt.Errorf("account not found for user %s", userID)
	}

	// Return a copy to avoid concurrent access issues
	result := make(map[string]decimal.Decimal)
	for asset, balance := range userBalances {
		result[asset] = balance
	}

	return result, nil
}

// GetCurrentPrice returns the current price for a symbol
func (me *MockExchange) GetCurrentPrice(symbol string) (decimal.Decimal, error) {
	me.mu.RLock()
	defer me.mu.RUnlock()

	price, exists := me.prices[symbol]
	if !exists {
		return decimal.Zero, fmt.Errorf("price not available for symbol: %s", symbol)
	}

	return price, nil
}

// GetOrderBook returns the order book for a symbol
func (me *MockExchange) GetOrderBook(symbol string) (*MockOrderBook, error) {
	me.mu.RLock()
	defer me.mu.RUnlock()

	orderBook, exists := me.orderBooks[symbol]
	if !exists {
		return nil, fmt.Errorf("order book not available for symbol: %s", symbol)
	}

	return orderBook, nil
}

// GetTrades returns trade history for a user
func (me *MockExchange) GetTrades(userID string, limit int) ([]*MockTrade, error) {
	me.mu.RLock()
	defer me.mu.RUnlock()

	userTrades, exists := me.trades[userID]
	if !exists {
		return []*MockTrade{}, nil
	}

	// Return the most recent trades
	start := 0
	if len(userTrades) > limit {
		start = len(userTrades) - limit
	}

	return userTrades[start:], nil
}

// initializePriceFeeds sets up initial price feeds
func (me *MockExchange) initializePriceFeeds() {
	basePrices := map[string]decimal.Decimal{
		"BTC/USDT": decimal.NewFromFloat(50000),
		"ETH/USDT": decimal.NewFromFloat(3000),
		"BNB/USDT": decimal.NewFromFloat(400),
		"ADA/USDT": decimal.NewFromFloat(1.5),
		"DOT/USDT": decimal.NewFromFloat(25),
	}

	for _, symbol := range me.config.SupportedPairs {
		basePrice := basePrices[symbol]
		if basePrice.IsZero() {
			basePrice = decimal.NewFromFloat(100) // Default price
		}

		me.prices[symbol] = basePrice
		me.priceFeeds[symbol] = &PriceFeed{
			Symbol:       symbol,
			CurrentPrice: basePrice,
			BasePrice:    basePrice,
			Volatility:   me.config.PriceVolatility,
			Trend:        decimal.Zero,
			LastUpdate:   time.Now(),
			PriceHistory: make([]PricePoint, 0),
		}

		// Initialize order book
		me.orderBooks[symbol] = me.generateOrderBook(symbol, basePrice)
	}
}

// priceSimulationLoop simulates price movements
func (me *MockExchange) priceSimulationLoop(ctx context.Context) {
	ticker := time.NewTicker(me.config.UpdateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-me.stopChan:
			return
		case <-ticker.C:
			me.updatePrices()
		}
	}
}

// updatePrices updates all price feeds
func (me *MockExchange) updatePrices() {
	me.mu.Lock()
	defer me.mu.Unlock()

	for symbol, feed := range me.priceFeeds {
		// Generate price movement
		change := me.generatePriceChange(feed)
		newPrice := feed.CurrentPrice.Add(change)

		// Ensure price doesn't go negative
		if newPrice.LessThan(decimal.NewFromFloat(0.01)) {
			newPrice = decimal.NewFromFloat(0.01)
		}

		// Update price feed
		feed.CurrentPrice = newPrice
		feed.LastUpdate = time.Now()

		// Add to price history
		pricePoint := PricePoint{
			Price:     newPrice,
			Volume:    decimal.NewFromFloat(rand.Float64() * 1000),
			Timestamp: time.Now(),
		}
		feed.PriceHistory = append(feed.PriceHistory, pricePoint)

		// Maintain history size
		if len(feed.PriceHistory) > 1000 {
			feed.PriceHistory = feed.PriceHistory[1:]
		}

		// Update current price
		me.prices[symbol] = newPrice

		// Update order book
		me.orderBooks[symbol] = me.generateOrderBook(symbol, newPrice)
	}
}

// generatePriceChange generates a realistic price change
func (me *MockExchange) generatePriceChange(feed *PriceFeed) decimal.Decimal {
	// Random walk with trend and volatility
	randomComponent := decimal.NewFromFloat((rand.Float64() - 0.5) * 2) // -1 to 1
	volatilityComponent := randomComponent.Mul(feed.Volatility).Mul(feed.CurrentPrice)
	trendComponent := feed.Trend.Mul(feed.CurrentPrice)

	return volatilityComponent.Add(trendComponent)
}

// generateOrderBook generates a realistic order book
func (me *MockExchange) generateOrderBook(symbol string, currentPrice decimal.Decimal) *MockOrderBook {
	orderBook := &MockOrderBook{
		Symbol:    symbol,
		Bids:      make([]*OrderBookEntry, 0),
		Asks:      make([]*OrderBookEntry, 0),
		UpdatedAt: time.Now(),
	}

	// Generate bids (buy orders)
	for i := 0; i < 10; i++ {
		priceOffset := decimal.NewFromFloat(float64(i+1) * 0.001) // 0.1% increments
		price := currentPrice.Sub(currentPrice.Mul(priceOffset))
		amount := decimal.NewFromFloat(rand.Float64() * 10)

		orderBook.Bids = append(orderBook.Bids, &OrderBookEntry{
			Price:  price,
			Amount: amount,
		})
	}

	// Generate asks (sell orders)
	for i := 0; i < 10; i++ {
		priceOffset := decimal.NewFromFloat(float64(i+1) * 0.001) // 0.1% increments
		price := currentPrice.Add(currentPrice.Mul(priceOffset))
		amount := decimal.NewFromFloat(rand.Float64() * 10)

		orderBook.Asks = append(orderBook.Asks, &OrderBookEntry{
			Price:  price,
			Amount: amount,
		})
	}

	return orderBook
}

// executeOrder simulates order execution
func (me *MockExchange) executeOrder(ctx context.Context, order *MockOrder) {
	// Simulate execution delay
	time.Sleep(me.config.ExecutionDelay)

	me.mu.Lock()
	defer me.mu.Unlock()

	// Check if order was cancelled
	if order.Status == "cancelled" {
		return
	}

	// Simulate order failure
	if rand.Float64() < me.config.FailureRate {
		order.Status = "failed"
		order.UpdatedAt = time.Now()
		me.logger.Warn(ctx, "Order execution failed", map[string]interface{}{
			"order_id": order.ID,
			"user_id":  order.UserID,
		})
		return
	}

	// Get current price
	currentPrice, exists := me.prices[order.Symbol]
	if !exists {
		order.Status = "failed"
		order.UpdatedAt = time.Now()
		return
	}

	// Calculate execution price with slippage
	executionPrice := me.calculateExecutionPrice(order, currentPrice)

	// Check if user has sufficient balance
	if !me.checkSufficientBalance(order, executionPrice) {
		order.Status = "failed"
		order.UpdatedAt = time.Now()
		me.logger.Warn(ctx, "Insufficient balance for order", map[string]interface{}{
			"order_id": order.ID,
			"user_id":  order.UserID,
		})
		return
	}

	// Execute the trade
	me.executeTrade(ctx, order, executionPrice)
}

// calculateExecutionPrice calculates the execution price with slippage
func (me *MockExchange) calculateExecutionPrice(order *MockOrder, currentPrice decimal.Decimal) decimal.Decimal {
	slippage := currentPrice.Mul(me.config.SlippageRate)

	if order.Type == "market" {
		// Market orders experience slippage
		if order.Side == "buy" {
			return currentPrice.Add(slippage)
		} else {
			return currentPrice.Sub(slippage)
		}
	} else {
		// Limit orders execute at specified price (if possible)
		return order.Price
	}
}

// checkSufficientBalance checks if user has sufficient balance for the order
func (me *MockExchange) checkSufficientBalance(order *MockOrder, executionPrice decimal.Decimal) bool {
	userBalances := me.balances[order.UserID]
	if userBalances == nil {
		return false
	}

	if order.Side == "buy" {
		// Check base currency balance (e.g., USDT)
		requiredAmount := order.Amount.Mul(executionPrice)
		fee := requiredAmount.Mul(me.config.TakerFee)
		totalRequired := requiredAmount.Add(fee)

		baseCurrency := me.getBaseCurrency(order.Symbol)
		balance := userBalances[baseCurrency]
		return balance.GreaterThanOrEqual(totalRequired)
	} else {
		// Check asset balance
		asset := me.getAsset(order.Symbol)
		balance := userBalances[asset]
		return balance.GreaterThanOrEqual(order.Amount)
	}
}

// executeTrade executes the trade and updates balances
func (me *MockExchange) executeTrade(ctx context.Context, order *MockOrder, executionPrice decimal.Decimal) {
	// Calculate fee
	tradeValue := order.Amount.Mul(executionPrice)
	fee := tradeValue.Mul(me.config.TakerFee)

	// Update balances
	userBalances := me.balances[order.UserID]
	asset := me.getAsset(order.Symbol)
	baseCurrency := me.getBaseCurrency(order.Symbol)

	if order.Side == "buy" {
		// Deduct base currency and fee
		totalCost := tradeValue.Add(fee)
		userBalances[baseCurrency] = userBalances[baseCurrency].Sub(totalCost)

		// Add asset
		if userBalances[asset].IsZero() {
			userBalances[asset] = order.Amount
		} else {
			userBalances[asset] = userBalances[asset].Add(order.Amount)
		}
	} else {
		// Deduct asset
		userBalances[asset] = userBalances[asset].Sub(order.Amount)

		// Add base currency minus fee
		netAmount := tradeValue.Sub(fee)
		if userBalances[baseCurrency].IsZero() {
			userBalances[baseCurrency] = netAmount
		} else {
			userBalances[baseCurrency] = userBalances[baseCurrency].Add(netAmount)
		}
	}

	// Update order status
	order.Status = "filled"
	order.FilledAmount = order.Amount
	order.RemainingAmount = decimal.Zero
	order.UpdatedAt = time.Now()
	executedAt := time.Now()
	order.ExecutedAt = &executedAt

	// Create trade record
	trade := &MockTrade{
		ID:        uuid.New().String(),
		OrderID:   order.ID,
		UserID:    order.UserID,
		Symbol:    order.Symbol,
		Side:      order.Side,
		Amount:    order.Amount,
		Price:     executionPrice,
		Fee:       fee,
		Timestamp: time.Now(),
	}

	// Store trade
	me.trades[order.UserID] = append(me.trades[order.UserID], trade)

	me.logger.Info(ctx, "Trade executed", map[string]interface{}{
		"trade_id":        trade.ID,
		"order_id":        order.ID,
		"user_id":         order.UserID,
		"symbol":          order.Symbol,
		"side":            order.Side,
		"amount":          order.Amount.String(),
		"execution_price": executionPrice.String(),
		"fee":             fee.String(),
	})
}

// getAsset extracts the asset from a trading pair symbol
func (me *MockExchange) getAsset(symbol string) string {
	// For symbols like "BTC/USDT", return "BTC"
	parts := []rune(symbol)
	for i, char := range parts {
		if char == '/' {
			return string(parts[:i])
		}
	}
	return symbol
}

// getBaseCurrency extracts the base currency from a trading pair symbol
func (me *MockExchange) getBaseCurrency(symbol string) string {
	// For symbols like "BTC/USDT", return "USDT"
	parts := []rune(symbol)
	for i, char := range parts {
		if char == '/' {
			return string(parts[i+1:])
		}
	}
	return me.config.BaseCurrency
}

// SetMarketCondition sets the market condition for simulation
func (me *MockExchange) SetMarketCondition(condition string) {
	me.mu.Lock()
	defer me.mu.Unlock()

	for _, feed := range me.priceFeeds {
		switch condition {
		case "bull":
			feed.Trend = decimal.NewFromFloat(0.001) // 0.1% upward trend
		case "bear":
			feed.Trend = decimal.NewFromFloat(-0.001) // 0.1% downward trend
		case "sideways":
			feed.Trend = decimal.Zero
		case "volatile":
			feed.Volatility = me.config.PriceVolatility.Mul(decimal.NewFromFloat(2))
		default:
			feed.Trend = decimal.Zero
			feed.Volatility = me.config.PriceVolatility
		}
	}

	me.logger.Info(context.Background(), "Market condition set", map[string]interface{}{
		"condition": condition,
	})
}

// GetExchangeInfo returns exchange information
func (me *MockExchange) GetExchangeInfo() map[string]interface{} {
	me.mu.RLock()
	defer me.mu.RUnlock()

	return map[string]interface{}{
		"name":            "MockExchange",
		"supported_pairs": me.config.SupportedPairs,
		"maker_fee":       me.config.MakerFee.String(),
		"taker_fee":       me.config.TakerFee.String(),
		"base_currency":   me.config.BaseCurrency,
		"is_running":      me.isRunning,
		"total_orders":    len(me.orders),
		"total_accounts":  len(me.balances),
	}
}
