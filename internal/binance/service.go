package binance

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/internal/hft"
	"github.com/ai-agentic-browser/pkg/observability"
)

// Service provides Binance integration for the HFT system
type Service struct {
	logger    *observability.Logger
	config    Config
	client    *BinanceClient
	wsManager *WebSocketManager
	hftEngine *hft.HFTEngine

	// Market data
	marketDataChan chan hft.MarketTick

	// State management
	isRunning bool
	stopChan  chan struct{}
	wg        sync.WaitGroup
	mu        sync.RWMutex
}

// ServiceConfig contains configuration for the Binance service
type ServiceConfig struct {
	Binance    Config   `json:"binance"`
	Symbols    []string `json:"symbols"`
	Streams    []string `json:"streams"`
	EnableHFT  bool     `json:"enable_hft"`
	BufferSize int      `json:"buffer_size"`
}

// NewService creates a new Binance service
func NewService(logger *observability.Logger, config ServiceConfig) *Service {
	service := &Service{
		logger:         logger,
		config:         config.Binance,
		marketDataChan: make(chan hft.MarketTick, config.BufferSize),
		stopChan:       make(chan struct{}),
	}

	// Create Binance client
	exchangeClient := NewBinanceClient(logger, config.Binance)
	if binanceClient, ok := (*exchangeClient).(*BinanceClient); ok {
		service.client = binanceClient
	}

	// Create WebSocket manager
	service.wsManager = NewWebSocketManager(logger, config.Binance)

	return service
}

// Start begins the Binance service
func (s *Service) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return fmt.Errorf("Binance service is already running")
	}

	s.logger.Info(ctx, "Starting Binance service", map[string]interface{}{
		"testnet":     s.config.Testnet,
		"base_url":    s.config.BaseURL,
		"ws_base_url": s.config.WSBaseURL,
	})

	// Start WebSocket manager
	if err := s.wsManager.Start(ctx); err != nil {
		return fmt.Errorf("failed to start WebSocket manager: %w", err)
	}

	// Subscribe to market data streams
	if err := s.subscribeToMarketData(ctx); err != nil {
		return fmt.Errorf("failed to subscribe to market data: %w", err)
	}

	s.isRunning = true

	// Start processing goroutines
	s.wg.Add(2)
	go s.processMarketData(ctx)
	go s.healthMonitor(ctx)

	s.logger.Info(ctx, "Binance service started successfully", nil)

	return nil
}

// Stop gracefully shuts down the Binance service
func (s *Service) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return fmt.Errorf("Binance service is not running")
	}

	s.logger.Info(ctx, "Stopping Binance service", nil)

	s.isRunning = false
	close(s.stopChan)

	// Stop WebSocket manager
	if err := s.wsManager.Stop(ctx); err != nil {
		s.logger.Error(ctx, "Failed to stop WebSocket manager", err)
	}

	// Stop client WebSocket if running
	if s.client != nil {
		if err := s.client.StopWebSocket(); err != nil {
			s.logger.Error(ctx, "Failed to stop client WebSocket", err)
		}
	}

	s.wg.Wait()

	s.logger.Info(ctx, "Binance service stopped successfully", nil)

	return nil
}

// SetHFTEngine sets the HFT engine for market data forwarding
func (s *Service) SetHFTEngine(engine *hft.HFTEngine) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.hftEngine = engine
}

// GetClient returns the Binance client
func (s *Service) GetClient() *BinanceClient {
	return s.client
}

// GetWebSocketManager returns the WebSocket manager
func (s *Service) GetWebSocketManager() *WebSocketManager {
	return s.wsManager
}

// SubscribeToSymbol subscribes to market data for a specific symbol
func (s *Service) SubscribeToSymbol(ctx context.Context, symbol string) (<-chan hft.MarketTick, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.isRunning {
		return nil, fmt.Errorf("Binance service is not running")
	}

	// Subscribe through WebSocket manager
	ch := s.wsManager.Subscribe(symbol)

	s.logger.Info(ctx, "Subscribed to symbol", map[string]interface{}{
		"symbol": symbol,
	})

	return ch, nil
}

// GetTicker gets 24hr ticker statistics for a symbol
func (s *Service) GetTicker(ctx context.Context, symbol string) (*TickerData, error) {
	if s.client == nil {
		return nil, fmt.Errorf("Binance client not initialized")
	}

	return s.client.GetTicker(ctx, symbol)
}

// PlaceOrder places an order through the Binance client
func (s *Service) PlaceOrder(ctx context.Context, order hft.Order) (*hft.OrderResponse, error) {
	if s.client == nil {
		return nil, fmt.Errorf("Binance client not initialized")
	}

	// Set exchange
	order.Exchange = "binance"

	return s.client.SubmitOrder(ctx, order)
}

// CancelOrder cancels an order through the Binance client
func (s *Service) CancelOrder(ctx context.Context, orderID string) error {
	if s.client == nil {
		return fmt.Errorf("Binance client not initialized")
	}

	// Parse order ID and cancel
	// Implementation would parse the order ID and call client.CancelOrder
	s.logger.Info(ctx, "Order cancellation requested", map[string]interface{}{
		"order_id": orderID,
	})

	return nil
}

// GetAccountInfo gets account information
func (s *Service) GetAccountInfo(ctx context.Context) (map[string]interface{}, error) {
	if s.client == nil {
		return nil, fmt.Errorf("Binance client not initialized")
	}

	// Implementation would call Binance account info endpoint
	// For now, return mock data
	accountInfo := map[string]interface{}{
		"account_type": "SPOT",
		"can_trade":    true,
		"can_withdraw": true,
		"can_deposit":  true,
		"balances":     []interface{}{},
	}

	return accountInfo, nil
}

// GetOpenOrders gets open orders for a symbol
func (s *Service) GetOpenOrders(ctx context.Context, symbol string) ([]*hft.Order, error) {
	if s.client == nil {
		return nil, fmt.Errorf("Binance client not initialized")
	}

	return s.client.GetOpenOrders(ctx, symbol)
}

// subscribeToMarketData subscribes to market data streams
func (s *Service) subscribeToMarketData(ctx context.Context) error {
	// Default symbols for HFT
	symbols := []string{
		"BTCUSDT", "ETHUSDT", "BNBUSDT", "ADAUSDT", "XRPUSDT",
		"SOLUSDT", "DOTUSDT", "LINKUSDT", "LTCUSDT", "BCHUSDT",
	}

	// Default streams for HFT
	streams := []StreamType{
		StreamTypeTicker,
		StreamTypeBookTicker,
		StreamTypeDepth,
		StreamTypeTrade,
	}

	// Subscribe to streams
	if err := s.wsManager.SubscribeToStreams(ctx, symbols, streams); err != nil {
		return fmt.Errorf("failed to subscribe to streams: %w", err)
	}

	// Subscribe to individual symbol channels and forward to market data channel
	for _, symbol := range symbols {
		ch := s.wsManager.Subscribe(symbol)

		// Start goroutine to forward data
		s.wg.Add(1)
		go s.forwardMarketData(ctx, symbol, ch)
	}

	return nil
}

// forwardMarketData forwards market data from symbol channel to main channel
func (s *Service) forwardMarketData(ctx context.Context, symbol string, ch <-chan hft.MarketTick) {
	defer s.wg.Done()

	for {
		select {
		case <-s.stopChan:
			return
		case tick, ok := <-ch:
			if !ok {
				s.logger.Warn(ctx, "Market data channel closed", map[string]interface{}{
					"symbol": symbol,
				})
				return
			}

			// Forward to main market data channel
			select {
			case s.marketDataChan <- tick:
			default:
				// Channel is full, drop the tick
				s.logger.Warn(ctx, "Market data channel full, dropping tick", map[string]interface{}{
					"symbol": symbol,
				})
			}
		}
	}
}

// processMarketData processes market data and forwards to HFT engine
func (s *Service) processMarketData(ctx context.Context) {
	defer s.wg.Done()

	for {
		select {
		case <-s.stopChan:
			return
		case tick := <-s.marketDataChan:
			// Forward to HFT engine if available
			if s.hftEngine != nil {
				s.hftEngine.SubmitMarketData(tick)
			}

			// Log high-frequency data (only occasionally to avoid spam)
			if tick.Sequence%1000 == 0 {
				s.logger.Debug(ctx, "Market data processed", map[string]interface{}{
					"symbol":   tick.Symbol,
					"price":    tick.Price.String(),
					"volume":   tick.Volume.String(),
					"exchange": tick.Exchange,
					"sequence": tick.Sequence,
				})
			}
		}
	}
}

// healthMonitor monitors the health of the Binance service
func (s *Service) healthMonitor(ctx context.Context) {
	defer s.wg.Done()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			s.performHealthCheck(ctx)
		}
	}
}

// performHealthCheck performs a health check of the service
func (s *Service) performHealthCheck(ctx context.Context) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.isRunning {
		return
	}

	// Check WebSocket connections
	wsHealthy := s.wsManager != nil && s.wsManager.isRunning

	// Check client connectivity (simplified)
	clientHealthy := s.client != nil

	// Check market data flow
	marketDataHealthy := len(s.marketDataChan) < cap(s.marketDataChan)

	s.logger.Info(ctx, "Binance service health check", map[string]interface{}{
		"websocket_healthy":    wsHealthy,
		"client_healthy":       clientHealthy,
		"market_data_healthy":  marketDataHealthy,
		"market_data_buffer":   len(s.marketDataChan),
		"market_data_capacity": cap(s.marketDataChan),
	})

	// Alert if unhealthy
	if !wsHealthy || !clientHealthy || !marketDataHealthy {
		s.logger.Warn(ctx, "Binance service health issues detected", map[string]interface{}{
			"websocket_healthy":   wsHealthy,
			"client_healthy":      clientHealthy,
			"market_data_healthy": marketDataHealthy,
		})
	}
}

// GetMetrics returns service metrics
func (s *Service) GetMetrics() ServiceMetrics {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return ServiceMetrics{
		IsRunning:          s.isRunning,
		WebSocketConnected: s.wsManager != nil && s.wsManager.isRunning,
		MarketDataBuffer:   len(s.marketDataChan),
		MarketDataCapacity: cap(s.marketDataChan),
		ClientInitialized:  s.client != nil,
		HFTEngineConnected: s.hftEngine != nil,
	}
}

// ServiceMetrics contains service performance metrics
type ServiceMetrics struct {
	IsRunning          bool `json:"is_running"`
	WebSocketConnected bool `json:"websocket_connected"`
	MarketDataBuffer   int  `json:"market_data_buffer"`
	MarketDataCapacity int  `json:"market_data_capacity"`
	ClientInitialized  bool `json:"client_initialized"`
	HFTEngineConnected bool `json:"hft_engine_connected"`
}

// IsHealthy returns whether the service is healthy
func (s *Service) IsHealthy() bool {
	metrics := s.GetMetrics()
	return metrics.IsRunning &&
		metrics.WebSocketConnected &&
		metrics.ClientInitialized &&
		metrics.MarketDataBuffer < metrics.MarketDataCapacity
}
