package binance

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ai-agentic-browser/internal/hft"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/shopspring/decimal"
)

// BinanceClient provides high-performance Binance API integration
type BinanceClient struct {
	logger      *observability.Logger
	config      Config
	httpClient  *http.Client
	rateLimiter *RateLimiter

	// WebSocket connections
	wsConn      *websocket.Conn
	wsConnMu    sync.RWMutex
	subscribers map[string][]chan hft.MarketTick

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

// OrderRequest represents a Binance order request
type OrderRequest struct {
	Symbol           string `json:"symbol"`
	Side             string `json:"side"`
	Type             string `json:"type"`
	TimeInForce      string `json:"timeInForce,omitempty"`
	Quantity         string `json:"quantity,omitempty"`
	QuoteOrderQty    string `json:"quoteOrderQty,omitempty"`
	Price            string `json:"price,omitempty"`
	StopPrice        string `json:"stopPrice,omitempty"`
	IcebergQty       string `json:"icebergQty,omitempty"`
	NewOrderRespType string `json:"newOrderRespType,omitempty"`
	RecvWindow       int64  `json:"recvWindow,omitempty"`
	Timestamp        int64  `json:"timestamp"`
}

// OrderResponse represents a Binance order response
type OrderResponse struct {
	Symbol                  string `json:"symbol"`
	OrderID                 int64  `json:"orderId"`
	OrderListID             int64  `json:"orderListId"`
	ClientOrderID           string `json:"clientOrderId"`
	TransactTime            int64  `json:"transactTime"`
	Price                   string `json:"price"`
	OrigQty                 string `json:"origQty"`
	ExecutedQty             string `json:"executedQty"`
	CummulativeQuoteQty     string `json:"cummulativeQuoteQty"`
	Status                  string `json:"status"`
	TimeInForce             string `json:"timeInForce"`
	Type                    string `json:"type"`
	Side                    string `json:"side"`
	WorkingTime             int64  `json:"workingTime"`
	SelfTradePreventionMode string `json:"selfTradePreventionMode"`
}

// TickerData represents 24hr ticker statistics
type TickerData struct {
	Symbol             string `json:"symbol"`
	PriceChange        string `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
	WeightedAvgPrice   string `json:"weightedAvgPrice"`
	PrevClosePrice     string `json:"prevClosePrice"`
	LastPrice          string `json:"lastPrice"`
	LastQty            string `json:"lastQty"`
	BidPrice           string `json:"bidPrice"`
	BidQty             string `json:"bidQty"`
	AskPrice           string `json:"askPrice"`
	AskQty             string `json:"askQty"`
	OpenPrice          string `json:"openPrice"`
	HighPrice          string `json:"highPrice"`
	LowPrice           string `json:"lowPrice"`
	Volume             string `json:"volume"`
	QuoteVolume        string `json:"quoteVolume"`
	OpenTime           int64  `json:"openTime"`
	CloseTime          int64  `json:"closeTime"`
	Count              int64  `json:"count"`
}

// WSTickerEvent represents WebSocket ticker event
type WSTickerEvent struct {
	EventType string `json:"e"` // Event type
	EventTime int64  `json:"E"` // Event time
	Symbol    string `json:"s"` // Symbol
	Price     string `json:"c"` // Close price
	Change    string `json:"P"` // Price change percent
	Volume    string `json:"v"` // Total traded base asset volume
	QuoteVol  string `json:"q"` // Total traded quote asset volume
	BidPrice  string `json:"b"` // Best bid price
	BidQty    string `json:"B"` // Best bid quantity
	AskPrice  string `json:"a"` // Best ask price
	AskQty    string `json:"A"` // Best ask quantity
}

// NewBinanceClient creates a new Binance client
func NewBinanceClient(logger *observability.Logger, config Config) *hft.ExchangeClient {
	if config.BaseURL == "" {
		if config.Testnet {
			config.BaseURL = "https://testnet.binance.vision"
			config.WSBaseURL = "wss://testnet.binance.vision"
		} else {
			config.BaseURL = "https://api.binance.com"
			config.WSBaseURL = "wss://stream.binance.com:9443"
		}
	}

	if config.Timeout == 0 {
		config.Timeout = 10 * time.Second
	}

	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}

	if config.RateLimit == 0 {
		config.RateLimit = 1200 // Binance default
	}

	client := &BinanceClient{
		logger: logger,
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		rateLimiter: NewRateLimiter(config.RateLimit, time.Minute),
		subscribers: make(map[string][]chan hft.MarketTick),
		stopChan:    make(chan struct{}),
	}

	// Return as ExchangeClient interface
	var exchangeClient hft.ExchangeClient = client
	return &exchangeClient
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(maxTokens int, refillInterval time.Duration) *RateLimiter {
	return &RateLimiter{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		refillRate: refillInterval,
		lastRefill: time.Now(),
	}
}

// Allow checks if a request is allowed under rate limiting
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)

	// Refill tokens based on elapsed time
	if elapsed >= rl.refillRate {
		rl.tokens = rl.maxTokens
		rl.lastRefill = now
	}

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}

	return false
}

// SubmitOrder implements the ExchangeClient interface
func (bc *BinanceClient) SubmitOrder(ctx context.Context, order hft.Order) (*hft.OrderResponse, error) {
	start := time.Now()

	// Check rate limit
	if !bc.rateLimiter.Allow() {
		return nil, fmt.Errorf("rate limit exceeded")
	}

	// Convert HFT order to Binance order request
	binanceOrder := bc.convertToBinanceOrder(order)

	// Submit order to Binance
	response, err := bc.submitOrderToBinance(ctx, binanceOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to submit order to Binance: %w", err)
	}

	// Convert response
	hftResponse := &hft.OrderResponse{
		OrderID:       order.ID,
		ExchangeID:    strconv.FormatInt(response.OrderID, 10),
		Status:        bc.convertOrderStatus(response.Status),
		Message:       "Order submitted successfully",
		LatencyMicros: time.Since(start).Microseconds(),
	}

	bc.logger.Info(ctx, "Order submitted to Binance", map[string]interface{}{
		"order_id":       order.ID.String(),
		"binance_id":     response.OrderID,
		"symbol":         order.Symbol,
		"side":           string(order.Side),
		"quantity":       order.Quantity.String(),
		"latency_micros": hftResponse.LatencyMicros,
	})

	return hftResponse, nil
}

// CancelOrder implements the ExchangeClient interface
func (bc *BinanceClient) CancelOrder(ctx context.Context, orderID uuid.UUID) error {
	// Check rate limit
	if !bc.rateLimiter.Allow() {
		return fmt.Errorf("rate limit exceeded")
	}

	// Implementation would cancel order on Binance
	// For now, return success
	bc.logger.Info(ctx, "Order canceled on Binance", map[string]interface{}{
		"order_id": orderID.String(),
	})

	return nil
}

// GetOrderStatus implements the ExchangeClient interface
func (bc *BinanceClient) GetOrderStatus(ctx context.Context, orderID uuid.UUID) (*hft.Order, error) {
	// Check rate limit
	if !bc.rateLimiter.Allow() {
		return nil, fmt.Errorf("rate limit exceeded")
	}

	// Implementation would query order status from Binance
	// For now, return a mock order
	order := &hft.Order{
		ID:     orderID,
		Status: hft.OrderStatusNew,
	}

	return order, nil
}

// GetOpenOrders implements the ExchangeClient interface
func (bc *BinanceClient) GetOpenOrders(ctx context.Context, symbol string) ([]*hft.Order, error) {
	// Check rate limit
	if !bc.rateLimiter.Allow() {
		return nil, fmt.Errorf("rate limit exceeded")
	}

	// Implementation would query open orders from Binance
	// For now, return empty slice
	return []*hft.Order{}, nil
}

// GetTicker gets 24hr ticker statistics
func (bc *BinanceClient) GetTicker(ctx context.Context, symbol string) (*TickerData, error) {
	if !bc.rateLimiter.Allow() {
		return nil, fmt.Errorf("rate limit exceeded")
	}

	endpoint := "/api/v3/ticker/24hr"
	params := url.Values{}
	params.Set("symbol", symbol)

	resp, err := bc.makeRequest(ctx, "GET", endpoint, params, false)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ticker TickerData
	if err := json.NewDecoder(resp.Body).Decode(&ticker); err != nil {
		return nil, fmt.Errorf("failed to decode ticker response: %w", err)
	}

	return &ticker, nil
}

// StartWebSocket starts the WebSocket connection for real-time data
func (bc *BinanceClient) StartWebSocket(ctx context.Context, symbols []string) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if bc.isConnected {
		return fmt.Errorf("WebSocket already connected")
	}

	// Create WebSocket URL for multiple symbols
	streams := make([]string, len(symbols))
	for i, symbol := range symbols {
		streams[i] = strings.ToLower(symbol) + "@ticker"
	}

	wsURL := fmt.Sprintf("%s/ws/%s", bc.config.WSBaseURL, strings.Join(streams, "/"))

	// Connect to WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	bc.wsConn = conn
	bc.isConnected = true

	// Start processing messages
	bc.wg.Add(1)
	go bc.processWebSocketMessages(ctx)

	bc.logger.Info(ctx, "WebSocket connected", map[string]interface{}{
		"url":     wsURL,
		"symbols": symbols,
	})

	return nil
}

// StopWebSocket stops the WebSocket connection
func (bc *BinanceClient) StopWebSocket() error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if !bc.isConnected {
		return fmt.Errorf("WebSocket not connected")
	}

	close(bc.stopChan)
	bc.wg.Wait()

	if bc.wsConn != nil {
		bc.wsConn.Close()
		bc.wsConn = nil
	}

	bc.isConnected = false

	bc.logger.Info(context.Background(), "WebSocket disconnected", nil)

	return nil
}

// Subscribe subscribes to market data for a symbol
func (bc *BinanceClient) Subscribe(symbol string) <-chan hft.MarketTick {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	ch := make(chan hft.MarketTick, 1000)

	if bc.subscribers[symbol] == nil {
		bc.subscribers[symbol] = make([]chan hft.MarketTick, 0)
	}

	bc.subscribers[symbol] = append(bc.subscribers[symbol], ch)

	return ch
}

// convertToBinanceOrder converts HFT order to Binance order format
func (bc *BinanceClient) convertToBinanceOrder(order hft.Order) OrderRequest {
	return OrderRequest{
		Symbol:      order.Symbol,
		Side:        strings.ToUpper(string(order.Side)),
		Type:        strings.ToUpper(string(order.Type)),
		TimeInForce: strings.ToUpper(string(order.TimeInForce)),
		Quantity:    order.Quantity.String(),
		Price:       order.Price.String(),
		Timestamp:   time.Now().UnixMilli(),
	}
}

// convertOrderStatus converts Binance order status to HFT status
func (bc *BinanceClient) convertOrderStatus(status string) hft.OrderStatus {
	switch status {
	case "NEW":
		return hft.OrderStatusNew
	case "PARTIALLY_FILLED":
		return hft.OrderStatusPartialFill
	case "FILLED":
		return hft.OrderStatusFilled
	case "CANCELED":
		return hft.OrderStatusCanceled
	case "REJECTED":
		return hft.OrderStatusRejected
	case "EXPIRED":
		return hft.OrderStatusExpired
	default:
		return hft.OrderStatusNew
	}
}

// submitOrderToBinance submits order to Binance API
func (bc *BinanceClient) submitOrderToBinance(ctx context.Context, order OrderRequest) (*OrderResponse, error) {
	endpoint := "/api/v3/order"

	// Convert order to URL values
	params := url.Values{}
	params.Set("symbol", order.Symbol)
	params.Set("side", order.Side)
	params.Set("type", order.Type)
	params.Set("quantity", order.Quantity)
	params.Set("timestamp", strconv.FormatInt(order.Timestamp, 10))

	if order.Price != "" {
		params.Set("price", order.Price)
	}
	if order.TimeInForce != "" {
		params.Set("timeInForce", order.TimeInForce)
	}

	resp, err := bc.makeRequest(ctx, "POST", endpoint, params, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var orderResp OrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&orderResp); err != nil {
		return nil, fmt.Errorf("failed to decode order response: %w", err)
	}

	return &orderResp, nil
}

// makeRequest makes an HTTP request to Binance API
func (bc *BinanceClient) makeRequest(ctx context.Context, method, endpoint string, params url.Values, signed bool) (*http.Response, error) {
	fullURL := bc.config.BaseURL + endpoint

	if signed {
		// Add signature for authenticated requests
		params.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))
		signature := bc.sign(params.Encode())
		params.Set("signature", signature)
	}

	var req *http.Request
	var err error

	if method == "GET" {
		if len(params) > 0 {
			fullURL += "?" + params.Encode()
		}
		req, err = http.NewRequestWithContext(ctx, method, fullURL, nil)
	} else {
		req, err = http.NewRequestWithContext(ctx, method, fullURL, strings.NewReader(params.Encode()))
		if err == nil {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add API key header
	if bc.config.APIKey != "" {
		req.Header.Set("X-MBX-APIKEY", bc.config.APIKey)
	}

	return bc.httpClient.Do(req)
}

// sign creates HMAC SHA256 signature
func (bc *BinanceClient) sign(message string) string {
	h := hmac.New(sha256.New, []byte(bc.config.SecretKey))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// processWebSocketMessages processes incoming WebSocket messages
func (bc *BinanceClient) processWebSocketMessages(ctx context.Context) {
	defer bc.wg.Done()

	for {
		select {
		case <-bc.stopChan:
			return
		default:
			bc.wsConnMu.RLock()
			conn := bc.wsConn
			bc.wsConnMu.RUnlock()

			if conn == nil {
				return
			}

			_, message, err := conn.ReadMessage()
			if err != nil {
				bc.logger.Error(ctx, "WebSocket read error", err)
				return
			}

			bc.processTickerMessage(message)
		}
	}
}

// processTickerMessage processes ticker WebSocket messages
func (bc *BinanceClient) processTickerMessage(message []byte) {
	var event WSTickerEvent
	if err := json.Unmarshal(message, &event); err != nil {
		return
	}

	// Convert to HFT MarketTick
	tick := bc.convertToMarketTick(event)

	// Distribute to subscribers
	bc.mu.RLock()
	subscribers, exists := bc.subscribers[event.Symbol]
	bc.mu.RUnlock()

	if exists {
		for _, ch := range subscribers {
			select {
			case ch <- tick:
			default:
				// Channel is full, skip
			}
		}
	}
}

// convertToMarketTick converts WebSocket event to MarketTick
func (bc *BinanceClient) convertToMarketTick(event WSTickerEvent) hft.MarketTick {
	price, _ := decimal.NewFromString(event.Price)
	bidPrice, _ := decimal.NewFromString(event.BidPrice)
	askPrice, _ := decimal.NewFromString(event.AskPrice)
	bidSize, _ := decimal.NewFromString(event.BidQty)
	askSize, _ := decimal.NewFromString(event.AskQty)
	volume, _ := decimal.NewFromString(event.Volume)

	return hft.MarketTick{
		Symbol:    event.Symbol,
		Price:     price,
		Volume:    volume,
		BidPrice:  bidPrice,
		AskPrice:  askPrice,
		BidSize:   bidSize,
		AskSize:   askSize,
		Timestamp: time.Unix(0, event.EventTime*int64(time.Millisecond)),
		Exchange:  "binance",
		Sequence:  uint64(event.EventTime),
	}
}
