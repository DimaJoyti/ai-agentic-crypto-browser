package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ai-agentic-browser/internal/exchanges/common"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/gorilla/websocket"
)

// WebSocketManager manages WebSocket connections for Binance
type WebSocketManager struct {
	logger      *observability.Logger
	config      Config
	connections map[string]*WSConnection
	subscribers map[string]*SubscriberGroup

	// State management
	isRunning bool
	stopChan  chan struct{}
	wg        sync.WaitGroup
	mu        sync.RWMutex
}

// WSConnection represents a WebSocket connection
type WSConnection struct {
	conn     *websocket.Conn
	url      string
	symbols  []string
	streams  []string
	lastPing time.Time
	isActive bool
	mu       sync.RWMutex
}

// SubscriberGroup manages subscribers for a specific data type
type SubscriberGroup struct {
	tickerSubs    map[string][]chan *common.TickerData
	orderBookSubs map[string][]chan *common.OrderBookData
	tradeSubs     map[string][]chan *common.TradeData
	userDataSubs  []chan *common.UserDataUpdate
	mu            sync.RWMutex
}

// NewWebSocketManager creates a new WebSocket manager
func NewWebSocketManager(logger *observability.Logger, config Config) *WebSocketManager {
	return &WebSocketManager{
		logger:      logger,
		config:      config,
		connections: make(map[string]*WSConnection),
		subscribers: make(map[string]*SubscriberGroup),
		stopChan:    make(chan struct{}),
	}
}

// Start starts the WebSocket manager
func (wsm *WebSocketManager) Start(ctx context.Context) error {
	wsm.mu.Lock()
	defer wsm.mu.Unlock()

	if wsm.isRunning {
		return nil
	}

	wsm.isRunning = true

	wsm.logger.Info(ctx, "WebSocket manager started", map[string]interface{}{
		"ws_base_url": wsm.config.WSBaseURL,
	})

	return nil
}

// Stop stops the WebSocket manager
func (wsm *WebSocketManager) Stop(ctx context.Context) error {
	wsm.mu.Lock()
	defer wsm.mu.Unlock()

	if !wsm.isRunning {
		return nil
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

// SubscribeToTicker subscribes to ticker updates for a symbol
func (wsm *WebSocketManager) SubscribeToTicker(ctx context.Context, symbol string) (<-chan *common.TickerData, error) {
	wsm.mu.Lock()
	defer wsm.mu.Unlock()

	symbol = strings.ToLower(symbol)
	ch := make(chan *common.TickerData, 1000)

	// Get or create subscriber group
	if wsm.subscribers[symbol] == nil {
		wsm.subscribers[symbol] = &SubscriberGroup{
			tickerSubs:    make(map[string][]chan *common.TickerData),
			orderBookSubs: make(map[string][]chan *common.OrderBookData),
			tradeSubs:     make(map[string][]chan *common.TradeData),
		}
	}

	group := wsm.subscribers[symbol]
	group.mu.Lock()
	if group.tickerSubs[symbol] == nil {
		group.tickerSubs[symbol] = make([]chan *common.TickerData, 0)
	}
	group.tickerSubs[symbol] = append(group.tickerSubs[symbol], ch)
	group.mu.Unlock()

	// Create WebSocket connection if needed
	if err := wsm.ensureConnection(ctx, symbol, "ticker"); err != nil {
		return nil, fmt.Errorf("failed to ensure connection: %w", err)
	}

	return ch, nil
}

// SubscribeToOrderBook subscribes to order book updates for a symbol
func (wsm *WebSocketManager) SubscribeToOrderBook(ctx context.Context, symbol string) (<-chan *common.OrderBookData, error) {
	wsm.mu.Lock()
	defer wsm.mu.Unlock()

	symbol = strings.ToLower(symbol)
	ch := make(chan *common.OrderBookData, 1000)

	// Get or create subscriber group
	if wsm.subscribers[symbol] == nil {
		wsm.subscribers[symbol] = &SubscriberGroup{
			tickerSubs:    make(map[string][]chan *common.TickerData),
			orderBookSubs: make(map[string][]chan *common.OrderBookData),
			tradeSubs:     make(map[string][]chan *common.TradeData),
		}
	}

	group := wsm.subscribers[symbol]
	group.mu.Lock()
	if group.orderBookSubs[symbol] == nil {
		group.orderBookSubs[symbol] = make([]chan *common.OrderBookData, 0)
	}
	group.orderBookSubs[symbol] = append(group.orderBookSubs[symbol], ch)
	group.mu.Unlock()

	// Create WebSocket connection if needed
	if err := wsm.ensureConnection(ctx, symbol, "depth"); err != nil {
		return nil, fmt.Errorf("failed to ensure connection: %w", err)
	}

	return ch, nil
}

// SubscribeToTrades subscribes to trade updates for a symbol
func (wsm *WebSocketManager) SubscribeToTrades(ctx context.Context, symbol string) (<-chan *common.TradeData, error) {
	wsm.mu.Lock()
	defer wsm.mu.Unlock()

	symbol = strings.ToLower(symbol)
	ch := make(chan *common.TradeData, 1000)

	// Get or create subscriber group
	if wsm.subscribers[symbol] == nil {
		wsm.subscribers[symbol] = &SubscriberGroup{
			tickerSubs:    make(map[string][]chan *common.TickerData),
			orderBookSubs: make(map[string][]chan *common.OrderBookData),
			tradeSubs:     make(map[string][]chan *common.TradeData),
		}
	}

	group := wsm.subscribers[symbol]
	group.mu.Lock()
	if group.tradeSubs[symbol] == nil {
		group.tradeSubs[symbol] = make([]chan *common.TradeData, 0)
	}
	group.tradeSubs[symbol] = append(group.tradeSubs[symbol], ch)
	group.mu.Unlock()

	// Create WebSocket connection if needed
	if err := wsm.ensureConnection(ctx, symbol, "trade"); err != nil {
		return nil, fmt.Errorf("failed to ensure connection: %w", err)
	}

	return ch, nil
}

// SubscribeToUserData subscribes to user data updates
func (wsm *WebSocketManager) SubscribeToUserData(ctx context.Context) (<-chan *common.UserDataUpdate, error) {
	wsm.mu.Lock()
	defer wsm.mu.Unlock()

	ch := make(chan *common.UserDataUpdate, 1000)

	// Get or create subscriber group for user data
	if wsm.subscribers["userdata"] == nil {
		wsm.subscribers["userdata"] = &SubscriberGroup{
			tickerSubs:    make(map[string][]chan *common.TickerData),
			orderBookSubs: make(map[string][]chan *common.OrderBookData),
			tradeSubs:     make(map[string][]chan *common.TradeData),
			userDataSubs:  make([]chan *common.UserDataUpdate, 0),
		}
	}

	group := wsm.subscribers["userdata"]
	group.mu.Lock()
	group.userDataSubs = append(group.userDataSubs, ch)
	group.mu.Unlock()

	// Create user data stream connection if needed
	if err := wsm.ensureUserDataConnection(ctx); err != nil {
		return nil, fmt.Errorf("failed to ensure user data connection: %w", err)
	}

	return ch, nil
}

// UnsubscribeAll unsubscribes from all streams
func (wsm *WebSocketManager) UnsubscribeAll(ctx context.Context) error {
	wsm.mu.Lock()
	defer wsm.mu.Unlock()

	// Close all subscriber channels
	for _, group := range wsm.subscribers {
		group.mu.Lock()
		for _, subs := range group.tickerSubs {
			for _, ch := range subs {
				close(ch)
			}
		}
		for _, subs := range group.orderBookSubs {
			for _, ch := range subs {
				close(ch)
			}
		}
		for _, subs := range group.tradeSubs {
			for _, ch := range subs {
				close(ch)
			}
		}
		for _, ch := range group.userDataSubs {
			close(ch)
		}
		group.mu.Unlock()
	}

	// Clear subscribers
	wsm.subscribers = make(map[string]*SubscriberGroup)

	// Close all connections
	for _, conn := range wsm.connections {
		conn.Close()
	}
	wsm.connections = make(map[string]*WSConnection)

	return nil
}

// ensureConnection ensures a WebSocket connection exists for the given symbol and stream
func (wsm *WebSocketManager) ensureConnection(ctx context.Context, symbol, streamType string) error {
	connKey := fmt.Sprintf("%s_%s", symbol, streamType)

	if conn, exists := wsm.connections[connKey]; exists && conn.isActive {
		return nil
	}

	// Build stream name
	streamName := fmt.Sprintf("%s@%s", symbol, streamType)
	if streamType == "depth" {
		streamName = fmt.Sprintf("%s@depth20@100ms", symbol)
	}

	// Create WebSocket URL
	wsURL := fmt.Sprintf("%s/%s", wsm.config.WSBaseURL, streamName)

	// Create connection
	conn, err := wsm.createConnection(ctx, wsURL, []string{symbol}, []string{streamName})
	if err != nil {
		return fmt.Errorf("failed to create WebSocket connection: %w", err)
	}

	wsm.connections[connKey] = conn

	// Start processing messages for this connection
	wsm.wg.Add(1)
	go wsm.processConnection(ctx, conn, streamType)

	wsm.logger.Info(ctx, "WebSocket connection established", map[string]interface{}{
		"symbol":      symbol,
		"stream_type": streamType,
		"url":         wsURL,
	})

	return nil
}

// ensureUserDataConnection ensures a user data stream connection exists
func (wsm *WebSocketManager) ensureUserDataConnection(ctx context.Context) error {
	connKey := "userdata"

	if conn, exists := wsm.connections[connKey]; exists && conn.isActive {
		return nil
	}

	// For user data stream, we would need to create a listen key first
	// This is a simplified implementation
	wsURL := wsm.config.WSBaseURL + "/userdata"

	// Create connection
	conn, err := wsm.createConnection(ctx, wsURL, []string{}, []string{"userdata"})
	if err != nil {
		return fmt.Errorf("failed to create user data connection: %w", err)
	}

	wsm.connections[connKey] = conn

	// Start processing messages for this connection
	wsm.wg.Add(1)
	go wsm.processConnection(ctx, conn, "userdata")

	return nil
}

// createConnection creates a new WebSocket connection
func (wsm *WebSocketManager) createConnection(ctx context.Context, wsURL string, symbols, streams []string) (*WSConnection, error) {
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
		streams:  streams,
		lastPing: time.Now(),
		isActive: true,
	}

	return wsConn, nil
}

// Close closes the WebSocket connection
func (conn *WSConnection) Close() {
	conn.mu.Lock()
	defer conn.mu.Unlock()

	if conn.conn != nil {
		conn.conn.Close()
		conn.conn = nil
	}
	conn.isActive = false
}

// processConnection processes messages from a WebSocket connection
func (wsm *WebSocketManager) processConnection(ctx context.Context, wsConn *WSConnection, streamType string) {
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
					"url":         wsConn.url,
					"stream_type": streamType,
				})

				// Attempt reconnection
				if wsm.isRunning {
					wsm.reconnectConnection(ctx, wsConn, streamType)
				}
				return
			}

			if messageType == websocket.TextMessage {
				wsm.processMessage(ctx, message, streamType)
			}
		}
	}
}

// processMessage processes a WebSocket message
func (wsm *WebSocketManager) processMessage(ctx context.Context, message []byte, streamType string) {
	switch streamType {
	case "ticker":
		wsm.processTickerMessage(ctx, message)
	case "depth":
		wsm.processDepthMessage(ctx, message)
	case "trade":
		wsm.processTradeMessage(ctx, message)
	case "userdata":
		wsm.processUserDataMessage(ctx, message)
	}
}

// processTickerMessage processes ticker messages
func (wsm *WebSocketManager) processTickerMessage(ctx context.Context, message []byte) {
	var tickerData BinanceWSTickerData
	if err := json.Unmarshal(message, &tickerData); err != nil {
		wsm.logger.Error(ctx, "Failed to unmarshal ticker data", err)
		return
	}

	// Convert to common ticker data
	commonTicker := &common.TickerData{
		Symbol:             tickerData.Symbol,
		PriceChange:        ParseDecimal(tickerData.PriceChange),
		PriceChangePercent: ParseDecimal(tickerData.PriceChangePercent),
		WeightedAvgPrice:   ParseDecimal(tickerData.WeightedAvgPrice),
		PrevClosePrice:     ParseDecimal(tickerData.FirstTradePrice),
		LastPrice:          ParseDecimal(tickerData.LastPrice),
		LastQty:            ParseDecimal(tickerData.LastQty),
		BidPrice:           ParseDecimal(tickerData.BestBidPrice),
		BidQty:             ParseDecimal(tickerData.BestBidQty),
		AskPrice:           ParseDecimal(tickerData.BestAskPrice),
		AskQty:             ParseDecimal(tickerData.BestAskQty),
		OpenPrice:          ParseDecimal(tickerData.OpenPrice),
		HighPrice:          ParseDecimal(tickerData.HighPrice),
		LowPrice:           ParseDecimal(tickerData.LowPrice),
		Volume:             ParseDecimal(tickerData.Volume),
		QuoteVolume:        ParseDecimal(tickerData.QuoteVolume),
		OpenTime:           ParseTime(tickerData.OpenTime),
		CloseTime:          ParseTime(tickerData.CloseTime),
		Count:              tickerData.TradeCount,
		Exchange:           "binance",
		Timestamp:          time.Now(),
	}

	// Send to subscribers
	wsm.sendToTickerSubscribers(strings.ToLower(tickerData.Symbol), commonTicker)
}

// processDepthMessage processes depth/order book messages
func (wsm *WebSocketManager) processDepthMessage(ctx context.Context, message []byte) {
	var depthData BinanceWSDepthData
	if err := json.Unmarshal(message, &depthData); err != nil {
		wsm.logger.Error(ctx, "Failed to unmarshal depth data", err)
		return
	}

	// Convert to common order book data
	bids := make([]common.PriceLevel, len(depthData.Bids))
	for i, bid := range depthData.Bids {
		if len(bid) >= 2 {
			bids[i] = common.PriceLevel{
				Price:    ParseDecimal(bid[0]),
				Quantity: ParseDecimal(bid[1]),
			}
		}
	}

	asks := make([]common.PriceLevel, len(depthData.Asks))
	for i, ask := range depthData.Asks {
		if len(ask) >= 2 {
			asks[i] = common.PriceLevel{
				Price:    ParseDecimal(ask[0]),
				Quantity: ParseDecimal(ask[1]),
			}
		}
	}

	commonOrderBook := &common.OrderBookData{
		Symbol:       depthData.Symbol,
		Bids:         bids,
		Asks:         asks,
		Timestamp:    ParseTime(depthData.EventTime),
		Exchange:     "binance",
		LastUpdateID: depthData.FinalUpdateId,
	}

	// Send to subscribers
	wsm.sendToOrderBookSubscribers(strings.ToLower(depthData.Symbol), commonOrderBook)
}

// processTradeMessage processes trade messages
func (wsm *WebSocketManager) processTradeMessage(ctx context.Context, message []byte) {
	var tradeData BinanceWSTradeData
	if err := json.Unmarshal(message, &tradeData); err != nil {
		wsm.logger.Error(ctx, "Failed to unmarshal trade data", err)
		return
	}

	side := "buy"
	if tradeData.IsBuyerMaker {
		side = "sell"
	}

	// Convert to common trade data
	commonTrade := &common.TradeData{
		ID:        fmt.Sprintf("%d", tradeData.TradeId),
		Symbol:    tradeData.Symbol,
		Price:     ParseDecimal(tradeData.Price),
		Quantity:  ParseDecimal(tradeData.Quantity),
		Side:      side,
		Timestamp: ParseTime(tradeData.TradeTime),
		Exchange:  "binance",
		IsBuyer:   !tradeData.IsBuyerMaker,
	}

	// Send to subscribers
	wsm.sendToTradeSubscribers(strings.ToLower(tradeData.Symbol), commonTrade)
}

// processUserDataMessage processes user data messages
func (wsm *WebSocketManager) processUserDataMessage(ctx context.Context, message []byte) {
	var userData BinanceWSUserData
	if err := json.Unmarshal(message, &userData); err != nil {
		wsm.logger.Error(ctx, "Failed to unmarshal user data", err)
		return
	}

	// Convert to common user data update
	commonUserData := &common.UserDataUpdate{
		Type:      userData.EventType,
		Data:      userData.Data,
		Timestamp: ParseTime(userData.EventTime),
		Exchange:  "binance",
	}

	// Send to subscribers
	wsm.sendToUserDataSubscribers(commonUserData)
}

// Subscriber notification methods

// sendToTickerSubscribers sends ticker data to all subscribers
func (wsm *WebSocketManager) sendToTickerSubscribers(symbol string, data *common.TickerData) {
	wsm.mu.RLock()
	group, exists := wsm.subscribers[symbol]
	wsm.mu.RUnlock()

	if !exists {
		return
	}

	group.mu.RLock()
	subs, exists := group.tickerSubs[symbol]
	group.mu.RUnlock()

	if !exists {
		return
	}

	for _, ch := range subs {
		select {
		case ch <- data:
		default:
			// Channel is full, skip
		}
	}
}

// sendToOrderBookSubscribers sends order book data to all subscribers
func (wsm *WebSocketManager) sendToOrderBookSubscribers(symbol string, data *common.OrderBookData) {
	wsm.mu.RLock()
	group, exists := wsm.subscribers[symbol]
	wsm.mu.RUnlock()

	if !exists {
		return
	}

	group.mu.RLock()
	subs, exists := group.orderBookSubs[symbol]
	group.mu.RUnlock()

	if !exists {
		return
	}

	for _, ch := range subs {
		select {
		case ch <- data:
		default:
			// Channel is full, skip
		}
	}
}

// sendToTradeSubscribers sends trade data to all subscribers
func (wsm *WebSocketManager) sendToTradeSubscribers(symbol string, data *common.TradeData) {
	wsm.mu.RLock()
	group, exists := wsm.subscribers[symbol]
	wsm.mu.RUnlock()

	if !exists {
		return
	}

	group.mu.RLock()
	subs, exists := group.tradeSubs[symbol]
	group.mu.RUnlock()

	if !exists {
		return
	}

	for _, ch := range subs {
		select {
		case ch <- data:
		default:
			// Channel is full, skip
		}
	}
}

// sendToUserDataSubscribers sends user data to all subscribers
func (wsm *WebSocketManager) sendToUserDataSubscribers(data *common.UserDataUpdate) {
	wsm.mu.RLock()
	group, exists := wsm.subscribers["userdata"]
	wsm.mu.RUnlock()

	if !exists {
		return
	}

	group.mu.RLock()
	subs := group.userDataSubs
	group.mu.RUnlock()

	for _, ch := range subs {
		select {
		case ch <- data:
		default:
			// Channel is full, skip
		}
	}
}

// reconnectConnection attempts to reconnect a WebSocket connection
func (wsm *WebSocketManager) reconnectConnection(ctx context.Context, wsConn *WSConnection, streamType string) {
	wsm.logger.Info(ctx, "Attempting to reconnect WebSocket", map[string]interface{}{
		"url":         wsConn.url,
		"stream_type": streamType,
	})

	// Wait before reconnecting
	time.Sleep(5 * time.Second)

	// Create new connection
	newConn, err := wsm.createConnection(ctx, wsConn.url, wsConn.symbols, wsConn.streams)
	if err != nil {
		wsm.logger.Error(ctx, "Failed to reconnect WebSocket", err)
		return
	}

	// Replace old connection
	wsm.mu.Lock()
	for key, conn := range wsm.connections {
		if conn == wsConn {
			wsm.connections[key] = newConn
			break
		}
	}
	wsm.mu.Unlock()

	// Start processing new connection
	wsm.wg.Add(1)
	go wsm.processConnection(ctx, newConn, streamType)

	wsm.logger.Info(ctx, "WebSocket reconnected successfully", map[string]interface{}{
		"url":         wsConn.url,
		"stream_type": streamType,
	})
}
