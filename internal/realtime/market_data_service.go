package realtime

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/gorilla/websocket"
	"github.com/shopspring/decimal"
)

// MarketDataService provides real-time market data from multiple exchanges
type MarketDataService struct {
	logger      *observability.Logger
	connections map[string]*ExchangeConnection
	subscribers map[string][]chan MarketUpdate
	config      MarketDataConfig
	mu          sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
}

// MarketDataConfig holds configuration for market data service
type MarketDataConfig struct {
	Exchanges       []ExchangeConfig `json:"exchanges"`
	ReconnectDelay  time.Duration    `json:"reconnect_delay"`
	PingInterval    time.Duration    `json:"ping_interval"`
	MaxReconnects   int              `json:"max_reconnects"`
	BufferSize      int              `json:"buffer_size"`
	EnableHeartbeat bool             `json:"enable_heartbeat"`
}

// ExchangeConfig holds configuration for a specific exchange
type ExchangeConfig struct {
	Name      string            `json:"name"`
	WSUrl     string            `json:"ws_url"`
	APIKey    string            `json:"api_key,omitempty"`
	APISecret string            `json:"api_secret,omitempty"`
	Symbols   []string          `json:"symbols"`
	Channels  []string          `json:"channels"`
	Headers   map[string]string `json:"headers,omitempty"`
	Enabled   bool              `json:"enabled"`
}

// ExchangeConnection represents a WebSocket connection to an exchange
type ExchangeConnection struct {
	Name         string
	Config       ExchangeConfig
	Conn         *websocket.Conn
	LastPing     time.Time
	LastPong     time.Time
	Reconnects   int
	IsConnected  bool
	MessageCount int64
	ErrorCount   int64
	mu           sync.RWMutex
}

// MarketUpdate represents a real-time market data update
type MarketUpdate struct {
	Exchange    string          `json:"exchange"`
	Symbol      string          `json:"symbol"`
	Type        UpdateType      `json:"type"`
	Price       decimal.Decimal `json:"price"`
	Volume      decimal.Decimal `json:"volume"`
	Bid         decimal.Decimal `json:"bid,omitempty"`
	Ask         decimal.Decimal `json:"ask,omitempty"`
	High24h     decimal.Decimal `json:"high_24h,omitempty"`
	Low24h      decimal.Decimal `json:"low_24h,omitempty"`
	Change24h   decimal.Decimal `json:"change_24h,omitempty"`
	Timestamp   time.Time       `json:"timestamp"`
	Sequence    int64           `json:"sequence,omitempty"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
}

// UpdateType represents the type of market data update
type UpdateType string

const (
	UpdateTypeTicker    UpdateType = "ticker"
	UpdateTypeTrade     UpdateType = "trade"
	UpdateTypeOrderBook UpdateType = "orderbook"
	UpdateTypeKline     UpdateType = "kline"
	UpdateTypeVolume    UpdateType = "volume"
)

// NewMarketDataService creates a new market data service
func NewMarketDataService(logger *observability.Logger, config MarketDataConfig) *MarketDataService {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &MarketDataService{
		logger:      logger,
		connections: make(map[string]*ExchangeConnection),
		subscribers: make(map[string][]chan MarketUpdate),
		config:      config,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Start begins the market data service
func (m *MarketDataService) Start() error {
	m.logger.Info(m.ctx, "Starting market data service", map[string]interface{}{
		"exchanges": len(m.config.Exchanges),
	})

	for _, exchangeConfig := range m.config.Exchanges {
		if !exchangeConfig.Enabled {
			continue
		}

		if err := m.connectToExchange(exchangeConfig); err != nil {
			m.logger.Error(m.ctx, "Failed to connect to exchange", err, map[string]interface{}{
				"exchange": exchangeConfig.Name,
			})
			continue
		}
	}

	// Start heartbeat monitoring
	if m.config.EnableHeartbeat {
		go m.heartbeatMonitor()
	}

	return nil
}

// Stop stops the market data service
func (m *MarketDataService) Stop() error {
	m.logger.Info(m.ctx, "Stopping market data service")
	
	m.cancel()
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Close all connections
	for name, conn := range m.connections {
		if conn.Conn != nil {
			conn.Conn.Close()
		}
		m.logger.Info(m.ctx, "Closed exchange connection", map[string]interface{}{
			"exchange": name,
		})
	}
	
	// Close all subscriber channels
	for symbol, channels := range m.subscribers {
		for _, ch := range channels {
			close(ch)
		}
		m.logger.Info(m.ctx, "Closed subscriber channels", map[string]interface{}{
			"symbol": symbol,
			"count":  len(channels),
		})
	}
	
	return nil
}

// Subscribe subscribes to market data updates for a symbol
func (m *MarketDataService) Subscribe(symbol string) <-chan MarketUpdate {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	ch := make(chan MarketUpdate, m.config.BufferSize)
	
	if m.subscribers[symbol] == nil {
		m.subscribers[symbol] = make([]chan MarketUpdate, 0)
	}
	
	m.subscribers[symbol] = append(m.subscribers[symbol], ch)
	
	m.logger.Info(m.ctx, "New subscriber added", map[string]interface{}{
		"symbol":      symbol,
		"subscribers": len(m.subscribers[symbol]),
	})
	
	return ch
}

// Unsubscribe removes a subscription for market data updates
func (m *MarketDataService) Unsubscribe(symbol string, ch <-chan MarketUpdate) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if channels, exists := m.subscribers[symbol]; exists {
		for i, subscriber := range channels {
			if subscriber == ch {
				// Remove the channel from the slice
				m.subscribers[symbol] = append(channels[:i], channels[i+1:]...)
				close(subscriber)
				break
			}
		}
		
		// Clean up empty subscriber lists
		if len(m.subscribers[symbol]) == 0 {
			delete(m.subscribers, symbol)
		}
	}
}

// GetConnectionStatus returns the status of all exchange connections
func (m *MarketDataService) GetConnectionStatus() map[string]ConnectionStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	status := make(map[string]ConnectionStatus)
	
	for name, conn := range m.connections {
		conn.mu.RLock()
		status[name] = ConnectionStatus{
			Exchange:     name,
			IsConnected:  conn.IsConnected,
			LastPing:     conn.LastPing,
			LastPong:     conn.LastPong,
			Reconnects:   conn.Reconnects,
			MessageCount: conn.MessageCount,
			ErrorCount:   conn.ErrorCount,
		}
		conn.mu.RUnlock()
	}
	
	return status
}

// ConnectionStatus represents the status of an exchange connection
type ConnectionStatus struct {
	Exchange     string    `json:"exchange"`
	IsConnected  bool      `json:"is_connected"`
	LastPing     time.Time `json:"last_ping"`
	LastPong     time.Time `json:"last_pong"`
	Reconnects   int       `json:"reconnects"`
	MessageCount int64     `json:"message_count"`
	ErrorCount   int64     `json:"error_count"`
}

// connectToExchange establishes a WebSocket connection to an exchange
func (m *MarketDataService) connectToExchange(config ExchangeConfig) error {
	u, err := url.Parse(config.WSUrl)
	if err != nil {
		return fmt.Errorf("invalid WebSocket URL: %w", err)
	}

	// Set up WebSocket headers
	headers := make(map[string][]string)
	for key, value := range config.Headers {
		headers[key] = []string{value}
	}

	// Create WebSocket connection
	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = 10 * time.Second
	
	conn, _, err := dialer.Dial(u.String(), headers)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", config.Name, err)
	}

	// Create exchange connection
	exchangeConn := &ExchangeConnection{
		Name:        config.Name,
		Config:      config,
		Conn:        conn,
		LastPing:    time.Now(),
		IsConnected: true,
	}

	m.mu.Lock()
	m.connections[config.Name] = exchangeConn
	m.mu.Unlock()

	// Start message handler
	go m.handleMessages(exchangeConn)

	// Send subscription messages
	if err := m.subscribeToChannels(exchangeConn); err != nil {
		m.logger.Error(m.ctx, "Failed to subscribe to channels", err, map[string]interface{}{
			"exchange": config.Name,
		})
	}

	m.logger.Info(m.ctx, "Connected to exchange", map[string]interface{}{
		"exchange": config.Name,
		"url":      config.WSUrl,
		"symbols":  config.Symbols,
		"channels": config.Channels,
	})

	return nil
}

// subscribeToChannels sends subscription messages to the exchange
func (m *MarketDataService) subscribeToChannels(conn *ExchangeConnection) error {
	// This is a simplified subscription - in reality, each exchange has different formats
	for _, symbol := range conn.Config.Symbols {
		for _, channel := range conn.Config.Channels {
			subscribeMsg := map[string]interface{}{
				"method": "SUBSCRIBE",
				"params": []string{fmt.Sprintf("%s@%s", symbol, channel)},
				"id":     time.Now().Unix(),
			}
			
			if err := conn.Conn.WriteJSON(subscribeMsg); err != nil {
				return fmt.Errorf("failed to subscribe to %s@%s: %w", symbol, channel, err)
			}
		}
	}
	
	return nil
}

// handleMessages processes incoming WebSocket messages
func (m *MarketDataService) handleMessages(conn *ExchangeConnection) {
	defer func() {
		conn.mu.Lock()
		conn.IsConnected = false
		conn.mu.Unlock()
		
		if conn.Conn != nil {
			conn.Conn.Close()
		}
		
		// Attempt reconnection if not cancelled
		if m.ctx.Err() == nil && conn.Reconnects < m.config.MaxReconnects {
			time.Sleep(m.config.ReconnectDelay)
			m.reconnectExchange(conn)
		}
	}()

	for {
		select {
		case <-m.ctx.Done():
			return
		default:
			var rawMessage json.RawMessage
			if err := conn.Conn.ReadJSON(&rawMessage); err != nil {
				conn.mu.Lock()
				conn.ErrorCount++
				conn.mu.Unlock()
				
				m.logger.Error(m.ctx, "WebSocket read error", err, map[string]interface{}{
					"exchange": conn.Name,
				})
				return
			}

			conn.mu.Lock()
			conn.MessageCount++
			conn.mu.Unlock()

			// Parse and distribute the message
			if update, err := m.parseMessage(conn.Name, rawMessage); err == nil {
				m.distributeUpdate(update)
			}
		}
	}
}

// parseMessage parses a raw WebSocket message into a MarketUpdate
func (m *MarketDataService) parseMessage(exchange string, rawMessage json.RawMessage) (MarketUpdate, error) {
	// This is a simplified parser - in reality, each exchange has different message formats
	var genericMessage map[string]interface{}
	if err := json.Unmarshal(rawMessage, &genericMessage); err != nil {
		return MarketUpdate{}, err
	}

	// Extract common fields (this would be exchange-specific in reality)
	update := MarketUpdate{
		Exchange:  exchange,
		Timestamp: time.Now(),
		Metadata:  rawMessage,
	}

	// Parse based on message structure (simplified)
	if symbol, ok := genericMessage["s"].(string); ok {
		update.Symbol = symbol
	}
	if price, ok := genericMessage["p"].(string); ok {
		if p, err := decimal.NewFromString(price); err == nil {
			update.Price = p
		}
	}
	if volume, ok := genericMessage["v"].(string); ok {
		if v, err := decimal.NewFromString(volume); err == nil {
			update.Volume = v
		}
	}

	update.Type = UpdateTypeTicker // Default type

	return update, nil
}

// distributeUpdate sends a market update to all subscribers
func (m *MarketDataService) distributeUpdate(update MarketUpdate) {
	m.mu.RLock()
	subscribers, exists := m.subscribers[update.Symbol]
	m.mu.RUnlock()

	if !exists {
		return
	}

	for _, ch := range subscribers {
		select {
		case ch <- update:
		default:
			// Channel is full, skip this update
		}
	}
}

// reconnectExchange attempts to reconnect to an exchange
func (m *MarketDataService) reconnectExchange(conn *ExchangeConnection) {
	conn.mu.Lock()
	conn.Reconnects++
	conn.mu.Unlock()

	m.logger.Info(m.ctx, "Attempting to reconnect to exchange", map[string]interface{}{
		"exchange":   conn.Name,
		"reconnects": conn.Reconnects,
	})

	if err := m.connectToExchange(conn.Config); err != nil {
		m.logger.Error(m.ctx, "Reconnection failed", err, map[string]interface{}{
			"exchange": conn.Name,
		})
	}
}

// heartbeatMonitor monitors connection health and sends pings
func (m *MarketDataService) heartbeatMonitor() {
	ticker := time.NewTicker(m.config.PingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.mu.RLock()
			for _, conn := range m.connections {
				if conn.IsConnected {
					go m.sendPing(conn)
				}
			}
			m.mu.RUnlock()
		}
	}
}

// sendPing sends a ping message to maintain connection
func (m *MarketDataService) sendPing(conn *ExchangeConnection) {
	conn.mu.Lock()
	defer conn.mu.Unlock()

	if err := conn.Conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
		m.logger.Error(m.ctx, "Failed to send ping", err, map[string]interface{}{
			"exchange": conn.Name,
		})
		conn.ErrorCount++
		return
	}

	conn.LastPing = time.Now()
}
