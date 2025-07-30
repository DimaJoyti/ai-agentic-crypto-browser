package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/ai-agentic-browser/internal/hft"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/gorilla/websocket"
	"github.com/shopspring/decimal"
)

// WebSocketManager manages multiple WebSocket connections for high-frequency data
type WebSocketManager struct {
	logger      *observability.Logger
	config      Config
	connections map[string]*WSConnection
	subscribers map[string][]chan hft.MarketTick

	// State management
	isRunning bool
	stopChan  chan struct{}
	wg        sync.WaitGroup
	mu        sync.RWMutex
}

// WSConnection represents a single WebSocket connection
type WSConnection struct {
	conn       *websocket.Conn
	url        string
	symbols    []string
	lastPing   time.Time
	isActive   bool
	reconnects int
	mu         sync.RWMutex
}

// StreamType represents different types of WebSocket streams
type StreamType string

const (
	StreamTypeTicker     StreamType = "ticker"
	StreamTypeDepth      StreamType = "depth"
	StreamTypeTrade      StreamType = "trade"
	StreamTypeKline      StreamType = "kline"
	StreamTypeBookTicker StreamType = "bookTicker"
)

// WSDepthEvent represents depth/orderbook WebSocket event
type WSDepthEvent struct {
	EventType     string     `json:"e"` // Event type
	EventTime     int64      `json:"E"` // Event time
	Symbol        string     `json:"s"` // Symbol
	FirstUpdateID int64      `json:"U"` // First update ID in event
	FinalUpdateID int64      `json:"u"` // Final update ID in event
	Bids          [][]string `json:"b"` // Bids to be updated
	Asks          [][]string `json:"a"` // Asks to be updated
}

// WSTradeEvent represents trade WebSocket event
type WSTradeEvent struct {
	EventType     string `json:"e"` // Event type
	EventTime     int64  `json:"E"` // Event time
	Symbol        string `json:"s"` // Symbol
	TradeID       int64  `json:"t"` // Trade ID
	Price         string `json:"p"` // Price
	Quantity      string `json:"q"` // Quantity
	BuyerOrderID  int64  `json:"b"` // Buyer order ID
	SellerOrderID int64  `json:"a"` // Seller order ID
	TradeTime     int64  `json:"T"` // Trade time
	IsBuyerMaker  bool   `json:"m"` // Is the buyer the market maker?
}

// WSKlineEvent represents kline/candlestick WebSocket event
type WSKlineEvent struct {
	EventType string `json:"e"` // Event type
	EventTime int64  `json:"E"` // Event time
	Symbol    string `json:"s"` // Symbol
	Kline     struct {
		StartTime           int64  `json:"t"` // Kline start time
		CloseTime           int64  `json:"T"` // Kline close time
		Symbol              string `json:"s"` // Symbol
		Interval            string `json:"i"` // Interval
		FirstTradeID        int64  `json:"f"` // First trade ID
		LastTradeID         int64  `json:"L"` // Last trade ID
		OpenPrice           string `json:"o"` // Open price
		ClosePrice          string `json:"c"` // Close price
		HighPrice           string `json:"h"` // High price
		LowPrice            string `json:"l"` // Low price
		BaseAssetVolume     string `json:"v"` // Base asset volume
		NumberOfTrades      int64  `json:"n"` // Number of trades
		IsKlineClosed       bool   `json:"x"` // Is this kline closed?
		QuoteAssetVolume    string `json:"q"` // Quote asset volume
		TakerBuyBaseVolume  string `json:"V"` // Taker buy base asset volume
		TakerBuyQuoteVolume string `json:"Q"` // Taker buy quote asset volume
	} `json:"k"`
}

// WSBookTickerEvent represents book ticker WebSocket event
type WSBookTickerEvent struct {
	UpdateID     int64  `json:"u"` // Order book updateId
	Symbol       string `json:"s"` // Symbol
	BestBidPrice string `json:"b"` // Best bid price
	BestBidQty   string `json:"B"` // Best bid qty
	BestAskPrice string `json:"a"` // Best ask price
	BestAskQty   string `json:"A"` // Best ask qty
}

// NewWebSocketManager creates a new WebSocket manager
func NewWebSocketManager(logger *observability.Logger, config Config) *WebSocketManager {
	return &WebSocketManager{
		logger:      logger,
		config:      config,
		connections: make(map[string]*WSConnection),
		subscribers: make(map[string][]chan hft.MarketTick),
		stopChan:    make(chan struct{}),
	}
}

// Start begins the WebSocket manager
func (wsm *WebSocketManager) Start(ctx context.Context) error {
	wsm.mu.Lock()
	defer wsm.mu.Unlock()

	if wsm.isRunning {
		return fmt.Errorf("WebSocket manager is already running")
	}

	wsm.isRunning = true

	// Start connection monitor
	wsm.wg.Add(1)
	go wsm.connectionMonitor(ctx)

	wsm.logger.Info(ctx, "WebSocket manager started", nil)

	return nil
}

// Stop gracefully shuts down the WebSocket manager
func (wsm *WebSocketManager) Stop(ctx context.Context) error {
	wsm.mu.Lock()
	defer wsm.mu.Unlock()

	if !wsm.isRunning {
		return fmt.Errorf("WebSocket manager is not running")
	}

	wsm.isRunning = false
	close(wsm.stopChan)

	// Close all connections
	for _, conn := range wsm.connections {
		conn.Close()
	}

	wsm.wg.Wait()

	wsm.logger.Info(ctx, "WebSocket manager stopped", nil)

	return nil
}

// SubscribeToStreams subscribes to multiple streams for symbols
func (wsm *WebSocketManager) SubscribeToStreams(ctx context.Context, symbols []string, streams []StreamType) error {
	wsm.mu.Lock()
	defer wsm.mu.Unlock()

	// Create stream names
	streamNames := make([]string, 0, len(symbols)*len(streams))
	for _, symbol := range symbols {
		for _, stream := range streams {
			streamName := wsm.buildStreamName(symbol, stream)
			streamNames = append(streamNames, streamName)
		}
	}

	// Create connection key
	connKey := strings.Join(streamNames, ",")

	// Check if connection already exists
	if _, exists := wsm.connections[connKey]; exists {
		return fmt.Errorf("connection already exists for streams: %s", connKey)
	}

	// Create WebSocket URL
	wsURL := wsm.buildWebSocketURL(streamNames)

	// Create connection
	conn, err := wsm.createConnection(ctx, wsURL, symbols)
	if err != nil {
		return fmt.Errorf("failed to create WebSocket connection: %w", err)
	}

	wsm.connections[connKey] = conn

	// Start processing messages for this connection
	wsm.wg.Add(1)
	go wsm.processConnection(ctx, conn)

	wsm.logger.Info(ctx, "Subscribed to WebSocket streams", map[string]interface{}{
		"symbols": symbols,
		"streams": streams,
		"url":     wsURL,
	})

	return nil
}

// Subscribe subscribes to market data for a symbol
func (wsm *WebSocketManager) Subscribe(symbol string) <-chan hft.MarketTick {
	wsm.mu.Lock()
	defer wsm.mu.Unlock()

	ch := make(chan hft.MarketTick, 10000) // Large buffer for HFT

	if wsm.subscribers[symbol] == nil {
		wsm.subscribers[symbol] = make([]chan hft.MarketTick, 0)
	}

	wsm.subscribers[symbol] = append(wsm.subscribers[symbol], ch)

	return ch
}

// buildStreamName builds a stream name for Binance WebSocket
func (wsm *WebSocketManager) buildStreamName(symbol string, streamType StreamType) string {
	symbol = strings.ToLower(symbol)

	switch streamType {
	case StreamTypeTicker:
		return symbol + "@ticker"
	case StreamTypeDepth:
		return symbol + "@depth@100ms" // 100ms depth updates
	case StreamTypeTrade:
		return symbol + "@trade"
	case StreamTypeKline:
		return symbol + "@kline_1m" // 1-minute klines
	case StreamTypeBookTicker:
		return symbol + "@bookTicker"
	default:
		return symbol + "@ticker"
	}
}

// buildWebSocketURL builds the WebSocket URL for multiple streams
func (wsm *WebSocketManager) buildWebSocketURL(streams []string) string {
	baseURL := wsm.config.WSBaseURL
	if baseURL == "" {
		if wsm.config.Testnet {
			baseURL = "wss://testnet.binance.vision"
		} else {
			baseURL = "wss://stream.binance.com:9443"
		}
	}

	// For multiple streams, use the combined streams endpoint
	if len(streams) > 1 {
		streamParam := url.QueryEscape(strings.Join(streams, "/"))
		return fmt.Sprintf("%s/stream?streams=%s", baseURL, streamParam)
	}

	// For single stream, use the simple endpoint
	return fmt.Sprintf("%s/ws/%s", baseURL, streams[0])
}

// createConnection creates a new WebSocket connection
func (wsm *WebSocketManager) createConnection(ctx context.Context, wsURL string, symbols []string) (*WSConnection, error) {
	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = 10 * time.Second

	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to dial WebSocket: %w", err)
	}

	// Set read deadline
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	// Set pong handler
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	wsConn := &WSConnection{
		conn:     conn,
		url:      wsURL,
		symbols:  symbols,
		lastPing: time.Now(),
		isActive: true,
	}

	return wsConn, nil
}

// processConnection processes messages from a WebSocket connection
func (wsm *WebSocketManager) processConnection(ctx context.Context, wsConn *WSConnection) {
	defer wsm.wg.Done()
	defer wsConn.Close()

	for {
		select {
		case <-wsm.stopChan:
			return
		default:
			wsConn.mu.RLock()
			conn := wsConn.conn
			wsConn.mu.RUnlock()

			if conn == nil {
				return
			}

			messageType, message, err := conn.ReadMessage()
			if err != nil {
				wsm.logger.Error(ctx, "WebSocket read error", err, map[string]interface{}{
					"url": wsConn.url,
				})

				// Attempt reconnection
				if wsm.isRunning {
					wsm.reconnectConnection(ctx, wsConn)
				}
				return
			}

			if messageType == websocket.TextMessage {
				wsm.processMessage(ctx, message)
			}
		}
	}
}

// processMessage processes a WebSocket message
func (wsm *WebSocketManager) processMessage(ctx context.Context, message []byte) {
	// Try to determine message type and process accordingly
	var baseEvent struct {
		EventType string          `json:"e"`
		Stream    string          `json:"stream,omitempty"`
		Data      json.RawMessage `json:"data,omitempty"`
	}

	if err := json.Unmarshal(message, &baseEvent); err != nil {
		return
	}

	// Handle combined stream format
	if baseEvent.Stream != "" && baseEvent.Data != nil {
		wsm.processStreamData(ctx, baseEvent.Stream, baseEvent.Data)
		return
	}

	// Handle direct stream format
	switch baseEvent.EventType {
	case "24hrTicker":
		wsm.processTickerEvent(ctx, message)
	case "depthUpdate":
		wsm.processDepthEvent(ctx, message)
	case "trade":
		wsm.processTradeEvent(ctx, message)
	case "kline":
		wsm.processKlineEvent(ctx, message)
	}
}

// processStreamData processes data from combined streams
func (wsm *WebSocketManager) processStreamData(ctx context.Context, stream string, data json.RawMessage) {
	// Extract stream type from stream name
	if strings.Contains(stream, "@ticker") {
		wsm.processTickerEvent(ctx, data)
	} else if strings.Contains(stream, "@depth") {
		wsm.processDepthEvent(ctx, data)
	} else if strings.Contains(stream, "@trade") {
		wsm.processTradeEvent(ctx, data)
	} else if strings.Contains(stream, "@kline") {
		wsm.processKlineEvent(ctx, data)
	} else if strings.Contains(stream, "@bookTicker") {
		wsm.processBookTickerEvent(ctx, data)
	}
}

// processTickerEvent processes ticker events
func (wsm *WebSocketManager) processTickerEvent(ctx context.Context, data []byte) {
	var event WSTickerEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return
	}

	tick := wsm.convertTickerToMarketTick(event)
	wsm.distributeMarketTick(tick)
}

// processDepthEvent processes depth/orderbook events
func (wsm *WebSocketManager) processDepthEvent(ctx context.Context, data []byte) {
	var event WSDepthEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return
	}

	// Convert depth event to market tick (simplified)
	// In a real implementation, you'd maintain an orderbook
	tick := hft.MarketTick{
		Symbol:    event.Symbol,
		Timestamp: time.Unix(0, event.EventTime*int64(time.Millisecond)),
		Exchange:  "binance",
		Sequence:  uint64(event.FinalUpdateID),
	}

	// Extract best bid/ask from depth data
	if len(event.Bids) > 0 {
		if price, err := decimal.NewFromString(event.Bids[0][0]); err == nil {
			tick.BidPrice = price
		}
		if size, err := decimal.NewFromString(event.Bids[0][1]); err == nil {
			tick.BidSize = size
		}
	}

	if len(event.Asks) > 0 {
		if price, err := decimal.NewFromString(event.Asks[0][0]); err == nil {
			tick.AskPrice = price
		}
		if size, err := decimal.NewFromString(event.Asks[0][1]); err == nil {
			tick.AskSize = size
		}
	}

	wsm.distributeMarketTick(tick)
}

// processTradeEvent processes trade events
func (wsm *WebSocketManager) processTradeEvent(ctx context.Context, data []byte) {
	var event WSTradeEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return
	}

	price, _ := decimal.NewFromString(event.Price)
	quantity, _ := decimal.NewFromString(event.Quantity)

	tick := hft.MarketTick{
		Symbol:    event.Symbol,
		Price:     price,
		Volume:    quantity,
		Timestamp: time.Unix(0, event.TradeTime*int64(time.Millisecond)),
		Exchange:  "binance",
		Sequence:  uint64(event.TradeID),
	}

	wsm.distributeMarketTick(tick)
}

// processKlineEvent processes kline/candlestick events
func (wsm *WebSocketManager) processKlineEvent(ctx context.Context, data []byte) {
	var event WSKlineEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return
	}

	// Only process closed klines for consistency
	if !event.Kline.IsKlineClosed {
		return
	}

	price, _ := decimal.NewFromString(event.Kline.ClosePrice)
	volume, _ := decimal.NewFromString(event.Kline.BaseAssetVolume)

	tick := hft.MarketTick{
		Symbol:    event.Symbol,
		Price:     price,
		Volume:    volume,
		Timestamp: time.Unix(0, event.Kline.CloseTime*int64(time.Millisecond)),
		Exchange:  "binance",
		Sequence:  uint64(event.Kline.LastTradeID),
	}

	wsm.distributeMarketTick(tick)
}

// processBookTickerEvent processes book ticker events
func (wsm *WebSocketManager) processBookTickerEvent(ctx context.Context, data []byte) {
	var event WSBookTickerEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return
	}

	bidPrice, _ := decimal.NewFromString(event.BestBidPrice)
	askPrice, _ := decimal.NewFromString(event.BestAskPrice)
	bidSize, _ := decimal.NewFromString(event.BestBidQty)
	askSize, _ := decimal.NewFromString(event.BestAskQty)

	tick := hft.MarketTick{
		Symbol:    event.Symbol,
		BidPrice:  bidPrice,
		AskPrice:  askPrice,
		BidSize:   bidSize,
		AskSize:   askSize,
		Timestamp: time.Now(), // Book ticker doesn't have timestamp
		Exchange:  "binance",
		Sequence:  uint64(event.UpdateID),
	}

	wsm.distributeMarketTick(tick)
}

// convertTickerToMarketTick converts ticker event to market tick
func (wsm *WebSocketManager) convertTickerToMarketTick(event WSTickerEvent) hft.MarketTick {
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

// distributeMarketTick distributes market tick to subscribers
func (wsm *WebSocketManager) distributeMarketTick(tick hft.MarketTick) {
	wsm.mu.RLock()
	subscribers, exists := wsm.subscribers[tick.Symbol]
	wsm.mu.RUnlock()

	if !exists {
		return
	}

	for _, ch := range subscribers {
		select {
		case ch <- tick:
		default:
			// Channel is full, skip this update
		}
	}
}

// reconnectConnection attempts to reconnect a WebSocket connection
func (wsm *WebSocketManager) reconnectConnection(ctx context.Context, wsConn *WSConnection) {
	wsConn.mu.Lock()
	defer wsConn.mu.Unlock()

	wsConn.reconnects++

	if wsConn.reconnects > 10 {
		wsm.logger.Error(ctx, "Max reconnection attempts reached", fmt.Errorf("giving up"), map[string]interface{}{
			"url":        wsConn.url,
			"reconnects": wsConn.reconnects,
		})
		return
	}

	// Wait before reconnecting
	time.Sleep(time.Duration(wsConn.reconnects) * time.Second)

	// Create new connection
	newConn, err := wsm.createConnection(ctx, wsConn.url, wsConn.symbols)
	if err != nil {
		wsm.logger.Error(ctx, "Failed to reconnect WebSocket", err, map[string]interface{}{
			"url":     wsConn.url,
			"attempt": wsConn.reconnects,
		})
		return
	}

	// Replace connection
	if wsConn.conn != nil {
		wsConn.conn.Close()
	}
	wsConn.conn = newConn.conn
	wsConn.isActive = true

	wsm.logger.Info(ctx, "WebSocket reconnected", map[string]interface{}{
		"url":     wsConn.url,
		"attempt": wsConn.reconnects,
	})
}

// connectionMonitor monitors WebSocket connections
func (wsm *WebSocketManager) connectionMonitor(ctx context.Context) {
	defer wsm.wg.Done()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-wsm.stopChan:
			return
		case <-ticker.C:
			wsm.pingConnections(ctx)
		}
	}
}

// pingConnections sends ping to all connections
func (wsm *WebSocketManager) pingConnections(ctx context.Context) {
	wsm.mu.RLock()
	connections := make([]*WSConnection, 0, len(wsm.connections))
	for _, conn := range wsm.connections {
		connections = append(connections, conn)
	}
	wsm.mu.RUnlock()

	for _, conn := range connections {
		conn.mu.Lock()
		if conn.conn != nil && conn.isActive {
			if err := conn.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				wsm.logger.Error(ctx, "Failed to ping WebSocket", err, map[string]interface{}{
					"url": conn.url,
				})
				conn.isActive = false
			} else {
				conn.lastPing = time.Now()
			}
		}
		conn.mu.Unlock()
	}
}

// Close closes a WebSocket connection
func (conn *WSConnection) Close() {
	conn.mu.Lock()
	defer conn.mu.Unlock()

	if conn.conn != nil {
		conn.conn.Close()
		conn.conn = nil
	}
	conn.isActive = false
}
