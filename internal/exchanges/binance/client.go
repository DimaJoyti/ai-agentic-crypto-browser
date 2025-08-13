package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ai-agentic-browser/internal/exchanges/common"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/shopspring/decimal"
)

// Client implements the unified ExchangeClient interface for Binance
type Client struct {
	logger      *observability.Logger
	config      Config
	httpClient  *http.Client
	rateLimiter *RateLimiter

	// WebSocket connections
	wsManager *WebSocketManager

	// Performance tracking
	latencyStats    *common.LatencyStats
	connectionStats *common.ConnectionStats

	// State management
	isConnected bool
	stopChan    chan struct{}
	wg          sync.WaitGroup
	mu          sync.RWMutex
}

// Config contains Binance API configuration
type Config struct {
	APIKey     string        `json:"api_key"`
	SecretKey  string        `json:"secret_key"`
	BaseURL    string        `json:"base_url"`
	WSBaseURL  string        `json:"ws_base_url"`
	Testnet    bool          `json:"testnet"`
	Timeout    time.Duration `json:"timeout"`
	MaxRetries int           `json:"max_retries"`
	RateLimit  int           `json:"rate_limit"` // requests per minute
}

// RateLimiter implements token bucket rate limiting
type RateLimiter struct {
	tokens     int
	maxTokens  int
	refillRate time.Duration
	lastRefill time.Time
	mu         sync.Mutex
}

// NewClient creates a new Binance client
func NewClient(logger *observability.Logger, config Config) *Client {
	if config.BaseURL == "" {
		if config.Testnet {
			config.BaseURL = "https://testnet.binance.vision"
		} else {
			config.BaseURL = "https://api.binance.com"
		}
	}

	if config.WSBaseURL == "" {
		if config.Testnet {
			config.WSBaseURL = "wss://testnet.binance.vision/ws"
		} else {
			config.WSBaseURL = "wss://stream.binance.com:9443/ws"
		}
	}

	if config.Timeout == 0 {
		config.Timeout = 10 * time.Second
	}

	if config.RateLimit == 0 {
		config.RateLimit = 1200 // Default Binance rate limit
	}

	client := &Client{
		logger: logger,
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		rateLimiter: &RateLimiter{
			tokens:     config.RateLimit,
			maxTokens:  config.RateLimit,
			refillRate: time.Minute,
			lastRefill: time.Now(),
		},
		latencyStats: &common.LatencyStats{
			LastUpdated: time.Now(),
		},
		connectionStats: &common.ConnectionStats{},
		stopChan:        make(chan struct{}),
	}

	// Initialize WebSocket manager
	client.wsManager = NewWebSocketManager(logger, config)

	return client
}

// Connection Management

// Connect establishes connection to Binance
func (c *Client) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isConnected {
		return nil
	}

	// Test API connectivity
	if err := c.testConnectivity(ctx); err != nil {
		return fmt.Errorf("failed to test connectivity: %w", err)
	}

	// Start WebSocket manager
	if err := c.wsManager.Start(ctx); err != nil {
		return fmt.Errorf("failed to start WebSocket manager: %w", err)
	}

	c.isConnected = true
	c.connectionStats.IsConnected = true
	c.connectionStats.ConnectedSince = time.Now()

	c.logger.Info(ctx, "Connected to Binance", map[string]interface{}{
		"testnet":  c.config.Testnet,
		"base_url": c.config.BaseURL,
	})

	return nil
}

// Disconnect closes connection to Binance
func (c *Client) Disconnect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isConnected {
		return nil
	}

	// Stop WebSocket manager
	if err := c.wsManager.Stop(ctx); err != nil {
		c.logger.Error(ctx, "Failed to stop WebSocket manager", err)
	}

	close(c.stopChan)
	c.wg.Wait()

	c.isConnected = false
	c.connectionStats.IsConnected = false

	c.logger.Info(ctx, "Disconnected from Binance", nil)

	return nil
}

// IsConnected returns connection status
func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isConnected
}

// GetExchangeName returns the exchange name
func (c *Client) GetExchangeName() string {
	return "binance"
}

// Market Data

// GetTicker gets 24hr ticker statistics for a symbol
func (c *Client) GetTicker(ctx context.Context, symbol string) (*common.TickerData, error) {
	start := time.Now()

	if !c.rateLimiter.Allow() {
		return nil, fmt.Errorf("rate limit exceeded")
	}

	endpoint := "/api/v3/ticker/24hr"
	params := url.Values{}
	params.Set("symbol", strings.ToUpper(symbol))

	resp, err := c.makeRequest(ctx, "GET", endpoint, params, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticker: %w", err)
	}

	var binanceTicker BinanceTickerResponse
	if err := json.Unmarshal(resp, &binanceTicker); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ticker response: %w", err)
	}

	ticker := c.convertTickerData(&binanceTicker)

	// Update latency stats
	latency := time.Since(start).Microseconds()
	c.updateLatencyStats(latency)

	return ticker, nil
}

// GetOrderBook gets order book data for a symbol
func (c *Client) GetOrderBook(ctx context.Context, symbol string, limit int) (*common.OrderBookData, error) {
	start := time.Now()

	if !c.rateLimiter.Allow() {
		return nil, fmt.Errorf("rate limit exceeded")
	}

	endpoint := "/api/v3/depth"
	params := url.Values{}
	params.Set("symbol", strings.ToUpper(symbol))
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}

	resp, err := c.makeRequest(ctx, "GET", endpoint, params, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get order book: %w", err)
	}

	var binanceOrderBook BinanceOrderBookResponse
	if err := json.Unmarshal(resp, &binanceOrderBook); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order book response: %w", err)
	}

	orderBook := c.convertOrderBookData(&binanceOrderBook, symbol)

	// Update latency stats
	latency := time.Since(start).Microseconds()
	c.updateLatencyStats(latency)

	return orderBook, nil
}

// GetRecentTrades gets recent trades for a symbol
func (c *Client) GetRecentTrades(ctx context.Context, symbol string, limit int) ([]*common.TradeData, error) {
	start := time.Now()

	if !c.rateLimiter.Allow() {
		return nil, fmt.Errorf("rate limit exceeded")
	}

	endpoint := "/api/v3/trades"
	params := url.Values{}
	params.Set("symbol", strings.ToUpper(symbol))
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}

	resp, err := c.makeRequest(ctx, "GET", endpoint, params, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent trades: %w", err)
	}

	var binanceTrades []BinanceTradeResponse
	if err := json.Unmarshal(resp, &binanceTrades); err != nil {
		return nil, fmt.Errorf("failed to unmarshal trades response: %w", err)
	}

	trades := make([]*common.TradeData, len(binanceTrades))
	for i, trade := range binanceTrades {
		trades[i] = c.convertTradeData(&trade, symbol)
	}

	// Update latency stats
	latency := time.Since(start).Microseconds()
	c.updateLatencyStats(latency)

	return trades, nil
}

// GetKlines gets candlestick data for a symbol
func (c *Client) GetKlines(ctx context.Context, symbol string, interval string, limit int) ([]*common.KlineData, error) {
	start := time.Now()

	if !c.rateLimiter.Allow() {
		return nil, fmt.Errorf("rate limit exceeded")
	}

	endpoint := "/api/v3/klines"
	params := url.Values{}
	params.Set("symbol", strings.ToUpper(symbol))
	params.Set("interval", interval)
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}

	resp, err := c.makeRequest(ctx, "GET", endpoint, params, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get klines: %w", err)
	}

	var binanceKlines [][]interface{}
	if err := json.Unmarshal(resp, &binanceKlines); err != nil {
		return nil, fmt.Errorf("failed to unmarshal klines response: %w", err)
	}

	klines := make([]*common.KlineData, len(binanceKlines))
	for i, kline := range binanceKlines {
		klines[i] = c.convertKlineData(kline, symbol)
	}

	// Update latency stats
	latency := time.Since(start).Microseconds()
	c.updateLatencyStats(latency)

	return klines, nil
}

// WebSocket Streaming

// SubscribeToTicker subscribes to ticker updates for a symbol
func (c *Client) SubscribeToTicker(ctx context.Context, symbol string) (<-chan *common.TickerData, error) {
	return c.wsManager.SubscribeToTicker(ctx, symbol)
}

// SubscribeToOrderBook subscribes to order book updates for a symbol
func (c *Client) SubscribeToOrderBook(ctx context.Context, symbol string) (<-chan *common.OrderBookData, error) {
	return c.wsManager.SubscribeToOrderBook(ctx, symbol)
}

// SubscribeToTrades subscribes to trade updates for a symbol
func (c *Client) SubscribeToTrades(ctx context.Context, symbol string) (<-chan *common.TradeData, error) {
	return c.wsManager.SubscribeToTrades(ctx, symbol)
}

// SubscribeToUserData subscribes to user data updates
func (c *Client) SubscribeToUserData(ctx context.Context) (<-chan *common.UserDataUpdate, error) {
	return c.wsManager.SubscribeToUserData(ctx)
}

// UnsubscribeAll unsubscribes from all streams
func (c *Client) UnsubscribeAll(ctx context.Context) error {
	return c.wsManager.UnsubscribeAll(ctx)
}

// Order Management

// PlaceOrder places an order
func (c *Client) PlaceOrder(ctx context.Context, order *common.OrderRequest) (*common.OrderResponse, error) {
	start := time.Now()

	if !c.rateLimiter.Allow() {
		return nil, fmt.Errorf("rate limit exceeded")
	}

	endpoint := "/api/v3/order"
	params := c.buildOrderParams(order)

	resp, err := c.makeRequest(ctx, "POST", endpoint, params, true)
	if err != nil {
		return nil, fmt.Errorf("failed to place order: %w", err)
	}

	var binanceOrder BinanceOrderResponse
	if err := json.Unmarshal(resp, &binanceOrder); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order response: %w", err)
	}

	orderResp := c.convertOrderResponse(&binanceOrder)
	orderResp.LatencyMicros = time.Since(start).Microseconds()

	// Update latency stats
	c.updateLatencyStats(orderResp.LatencyMicros)

	c.logger.Info(ctx, "Order placed successfully", map[string]interface{}{
		"order_id":       orderResp.OrderID,
		"symbol":         orderResp.Symbol,
		"side":           string(orderResp.Side),
		"quantity":       orderResp.Quantity.String(),
		"latency_micros": orderResp.LatencyMicros,
	})

	return orderResp, nil
}

// CancelOrder cancels an order
func (c *Client) CancelOrder(ctx context.Context, symbol string, orderID string) (*common.OrderResponse, error) {
	start := time.Now()

	if !c.rateLimiter.Allow() {
		return nil, fmt.Errorf("rate limit exceeded")
	}

	endpoint := "/api/v3/order"
	params := url.Values{}
	params.Set("symbol", strings.ToUpper(symbol))
	params.Set("orderId", orderID)
	params.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))

	resp, err := c.makeRequest(ctx, "DELETE", endpoint, params, true)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel order: %w", err)
	}

	var binanceOrder BinanceOrderResponse
	if err := json.Unmarshal(resp, &binanceOrder); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cancel order response: %w", err)
	}

	orderResp := c.convertOrderResponse(&binanceOrder)
	orderResp.LatencyMicros = time.Since(start).Microseconds()

	// Update latency stats
	c.updateLatencyStats(orderResp.LatencyMicros)

	return orderResp, nil
}

// CancelAllOrders cancels all open orders for a symbol
func (c *Client) CancelAllOrders(ctx context.Context, symbol string) ([]*common.OrderResponse, error) {
	start := time.Now()

	if !c.rateLimiter.Allow() {
		return nil, fmt.Errorf("rate limit exceeded")
	}

	endpoint := "/api/v3/openOrders"
	params := url.Values{}
	params.Set("symbol", strings.ToUpper(symbol))
	params.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))

	resp, err := c.makeRequest(ctx, "DELETE", endpoint, params, true)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel all orders: %w", err)
	}

	var binanceOrders []BinanceOrderResponse
	if err := json.Unmarshal(resp, &binanceOrders); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cancel all orders response: %w", err)
	}

	orders := make([]*common.OrderResponse, len(binanceOrders))
	for i, order := range binanceOrders {
		orders[i] = c.convertOrderResponse(&order)
		orders[i].LatencyMicros = time.Since(start).Microseconds()
	}

	// Update latency stats
	c.updateLatencyStats(time.Since(start).Microseconds())

	return orders, nil
}

// GetOrder gets order information
func (c *Client) GetOrder(ctx context.Context, symbol string, orderID string) (*common.OrderResponse, error) {
	start := time.Now()

	if !c.rateLimiter.Allow() {
		return nil, fmt.Errorf("rate limit exceeded")
	}

	endpoint := "/api/v3/order"
	params := url.Values{}
	params.Set("symbol", strings.ToUpper(symbol))
	params.Set("orderId", orderID)
	params.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))

	resp, err := c.makeRequest(ctx, "GET", endpoint, params, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	var binanceOrder BinanceOrderResponse
	if err := json.Unmarshal(resp, &binanceOrder); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order response: %w", err)
	}

	orderResp := c.convertOrderResponse(&binanceOrder)
	orderResp.LatencyMicros = time.Since(start).Microseconds()

	// Update latency stats
	c.updateLatencyStats(orderResp.LatencyMicros)

	return orderResp, nil
}

// GetOpenOrders gets open orders for a symbol
func (c *Client) GetOpenOrders(ctx context.Context, symbol string) ([]*common.OrderResponse, error) {
	start := time.Now()

	if !c.rateLimiter.Allow() {
		return nil, fmt.Errorf("rate limit exceeded")
	}

	endpoint := "/api/v3/openOrders"
	params := url.Values{}
	if symbol != "" {
		params.Set("symbol", strings.ToUpper(symbol))
	}
	params.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))

	resp, err := c.makeRequest(ctx, "GET", endpoint, params, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get open orders: %w", err)
	}

	var binanceOrders []BinanceOrderResponse
	if err := json.Unmarshal(resp, &binanceOrders); err != nil {
		return nil, fmt.Errorf("failed to unmarshal open orders response: %w", err)
	}

	orders := make([]*common.OrderResponse, len(binanceOrders))
	for i, order := range binanceOrders {
		orders[i] = c.convertOrderResponse(&order)
		orders[i].LatencyMicros = time.Since(start).Microseconds()
	}

	// Update latency stats
	c.updateLatencyStats(time.Since(start).Microseconds())

	return orders, nil
}

// GetOrderHistory gets order history for a symbol
func (c *Client) GetOrderHistory(ctx context.Context, symbol string, limit int) ([]*common.OrderResponse, error) {
	start := time.Now()

	if !c.rateLimiter.Allow() {
		return nil, fmt.Errorf("rate limit exceeded")
	}

	endpoint := "/api/v3/allOrders"
	params := url.Values{}
	params.Set("symbol", strings.ToUpper(symbol))
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	params.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))

	resp, err := c.makeRequest(ctx, "GET", endpoint, params, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get order history: %w", err)
	}

	var binanceOrders []BinanceOrderResponse
	if err := json.Unmarshal(resp, &binanceOrders); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order history response: %w", err)
	}

	orders := make([]*common.OrderResponse, len(binanceOrders))
	for i, order := range binanceOrders {
		orders[i] = c.convertOrderResponse(&order)
		orders[i].LatencyMicros = time.Since(start).Microseconds()
	}

	// Update latency stats
	c.updateLatencyStats(time.Since(start).Microseconds())

	return orders, nil
}

// Account Management

// GetAccountInfo gets account information
func (c *Client) GetAccountInfo(ctx context.Context) (*common.AccountInfo, error) {
	start := time.Now()

	if !c.rateLimiter.Allow() {
		return nil, fmt.Errorf("rate limit exceeded")
	}

	endpoint := "/api/v3/account"
	params := url.Values{}
	params.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))

	resp, err := c.makeRequest(ctx, "GET", endpoint, params, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get account info: %w", err)
	}

	var binanceAccount BinanceAccountResponse
	if err := json.Unmarshal(resp, &binanceAccount); err != nil {
		return nil, fmt.Errorf("failed to unmarshal account response: %w", err)
	}

	accountInfo := &common.AccountInfo{
		AccountType: binanceAccount.AccountType,
		CanTrade:    binanceAccount.CanTrade,
		CanWithdraw: binanceAccount.CanWithdraw,
		CanDeposit:  binanceAccount.CanDeposit,
		UpdateTime:  ParseTime(binanceAccount.UpdateTime),
		Exchange:    "binance",
	}

	// Calculate total balances
	totalBalance := decimal.NewFromInt(0)
	for _, balance := range binanceAccount.Balances {
		free := ParseDecimal(balance.Free)
		locked := ParseDecimal(balance.Locked)
		totalBalance = totalBalance.Add(free).Add(locked)
	}
	accountInfo.TotalWalletBalance = totalBalance

	// Update latency stats
	c.updateLatencyStats(time.Since(start).Microseconds())

	return accountInfo, nil
}

// GetBalances gets account balances
func (c *Client) GetBalances(ctx context.Context) ([]*common.Balance, error) {
	start := time.Now()

	if !c.rateLimiter.Allow() {
		return nil, fmt.Errorf("rate limit exceeded")
	}

	endpoint := "/api/v3/account"
	params := url.Values{}
	params.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))

	resp, err := c.makeRequest(ctx, "GET", endpoint, params, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get balances: %w", err)
	}

	var binanceAccount BinanceAccountResponse
	if err := json.Unmarshal(resp, &binanceAccount); err != nil {
		return nil, fmt.Errorf("failed to unmarshal account response: %w", err)
	}

	balances := make([]*common.Balance, 0, len(binanceAccount.Balances))
	for _, balance := range binanceAccount.Balances {
		free := ParseDecimal(balance.Free)
		locked := ParseDecimal(balance.Locked)
		total := free.Add(locked)

		// Only include non-zero balances
		if !total.IsZero() {
			balances = append(balances, &common.Balance{
				Asset:  balance.Asset,
				Free:   free,
				Locked: locked,
				Total:  total,
			})
		}
	}

	// Update latency stats
	c.updateLatencyStats(time.Since(start).Microseconds())

	return balances, nil
}

// GetTradingFees gets trading fees for a symbol
func (c *Client) GetTradingFees(ctx context.Context, symbol string) (*common.TradingFees, error) {
	start := time.Now()

	if !c.rateLimiter.Allow() {
		return nil, fmt.Errorf("rate limit exceeded")
	}

	endpoint := "/api/v3/account"
	params := url.Values{}
	params.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))

	resp, err := c.makeRequest(ctx, "GET", endpoint, params, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get trading fees: %w", err)
	}

	var binanceAccount BinanceAccountResponse
	if err := json.Unmarshal(resp, &binanceAccount); err != nil {
		return nil, fmt.Errorf("failed to unmarshal account response: %w", err)
	}

	// Convert commission rates from basis points to decimal
	makerCommission := decimal.NewFromInt(binanceAccount.MakerCommission).Div(decimal.NewFromInt(10000))
	takerCommission := decimal.NewFromInt(binanceAccount.TakerCommission).Div(decimal.NewFromInt(10000))

	fees := &common.TradingFees{
		Symbol:          symbol,
		MakerCommission: makerCommission,
		TakerCommission: takerCommission,
	}

	// Update latency stats
	c.updateLatencyStats(time.Since(start).Microseconds())

	return fees, nil
}

// Advanced Order Types

// PlaceStopLossOrder places a stop-loss order
func (c *Client) PlaceStopLossOrder(ctx context.Context, req *common.StopLossOrderRequest) (*common.OrderResponse, error) {
	orderReq := &common.OrderRequest{
		Symbol:        req.Symbol,
		Side:          req.Side,
		Type:          common.OrderTypeStopLoss,
		Quantity:      req.Quantity,
		StopPrice:     req.StopPrice,
		Price:         req.Price,
		TimeInForce:   req.TimeInForce,
		ClientOrderID: req.ClientOrderID,
	}

	if !req.Price.IsZero() {
		orderReq.Type = common.OrderTypeStopLossLimit
	}

	return c.PlaceOrder(ctx, orderReq)
}

// PlaceTakeProfitOrder places a take-profit order
func (c *Client) PlaceTakeProfitOrder(ctx context.Context, req *common.TakeProfitOrderRequest) (*common.OrderResponse, error) {
	orderReq := &common.OrderRequest{
		Symbol:        req.Symbol,
		Side:          req.Side,
		Type:          common.OrderTypeTakeProfit,
		Quantity:      req.Quantity,
		StopPrice:     req.StopPrice,
		Price:         req.Price,
		TimeInForce:   req.TimeInForce,
		ClientOrderID: req.ClientOrderID,
	}

	if !req.Price.IsZero() {
		orderReq.Type = common.OrderTypeTakeProfitLimit
	}

	return c.PlaceOrder(ctx, orderReq)
}

// PlaceIcebergOrder places an iceberg order
func (c *Client) PlaceIcebergOrder(ctx context.Context, req *common.IcebergOrderRequest) (*common.OrderResponse, error) {
	orderReq := &common.OrderRequest{
		Symbol:        req.Symbol,
		Side:          req.Side,
		Type:          common.OrderTypeLimit,
		Quantity:      req.Quantity,
		Price:         req.Price,
		TimeInForce:   req.TimeInForce,
		ClientOrderID: req.ClientOrderID,
		IcebergQty:    req.IcebergQty,
	}

	return c.PlaceOrder(ctx, orderReq)
}

// PlaceTWAPOrder places a TWAP order (simplified implementation)
func (c *Client) PlaceTWAPOrder(ctx context.Context, req *common.TWAPOrderRequest) (*common.OrderResponse, error) {
	// For TWAP orders, we would typically break the order into smaller chunks
	// and execute them over time. This is a simplified implementation.
	orderReq := &common.OrderRequest{
		Symbol:        req.Symbol,
		Side:          req.Side,
		Type:          common.OrderTypeLimit,
		Quantity:      req.Quantity,
		Price:         req.PriceLimit,
		TimeInForce:   common.TimeInForceGTC,
		ClientOrderID: req.ClientOrderID,
		Metadata: map[string]interface{}{
			"twap_duration": req.Duration.String(),
			"twap_interval": req.Interval.String(),
		},
	}

	return c.PlaceOrder(ctx, orderReq)
}

// Risk Management

// GetPositionRisk gets position risk information
func (c *Client) GetPositionRisk(ctx context.Context, symbol string) (*common.PositionRisk, error) {
	// Binance Spot doesn't have positions like futures
	// This is a placeholder implementation
	return &common.PositionRisk{
		Symbol:           symbol,
		PositionAmt:      decimal.NewFromInt(0),
		EntryPrice:       decimal.NewFromInt(0),
		MarkPrice:        decimal.NewFromInt(0),
		UnrealizedPnL:    decimal.NewFromInt(0),
		LiquidationPrice: decimal.NewFromInt(0),
		Leverage:         decimal.NewFromInt(1),
		MaxNotionalValue: decimal.NewFromInt(1000000),
		MarginType:       "ISOLATED",
		IsolatedMargin:   decimal.NewFromInt(0),
		IsAutoAddMargin:  false,
		PositionSide:     "BOTH",
		UpdateTime:       time.Now(),
	}, nil
}

// GetMaxOrderSize gets maximum order size for a symbol
func (c *Client) GetMaxOrderSize(ctx context.Context, symbol string) (decimal.Decimal, error) {
	// This would typically query exchange info for symbol filters
	// For now, return a reasonable default
	return decimal.NewFromInt(1000), nil
}

// ValidateOrder validates an order request
func (c *Client) ValidateOrder(ctx context.Context, order *common.OrderRequest) error {
	if order.Symbol == "" {
		return fmt.Errorf("symbol is required")
	}

	if order.Quantity.IsZero() || order.Quantity.IsNegative() {
		return fmt.Errorf("quantity must be positive")
	}

	if order.Type == common.OrderTypeLimit && (order.Price.IsZero() || order.Price.IsNegative()) {
		return fmt.Errorf("price must be positive for limit orders")
	}

	if (order.Type == common.OrderTypeStopLoss || order.Type == common.OrderTypeStopLossLimit ||
		order.Type == common.OrderTypeTakeProfit || order.Type == common.OrderTypeTakeProfitLimit) &&
		(order.StopPrice.IsZero() || order.StopPrice.IsNegative()) {
		return fmt.Errorf("stop price must be positive for stop orders")
	}

	return nil
}

// Performance Metrics

// GetLatencyStats returns latency statistics
func (c *Client) GetLatencyStats() *common.LatencyStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Create a copy to avoid race conditions
	stats := *c.latencyStats
	return &stats
}

// GetConnectionStats returns connection statistics
func (c *Client) GetConnectionStats() *common.ConnectionStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Create a copy to avoid race conditions
	stats := *c.connectionStats
	return &stats
}
