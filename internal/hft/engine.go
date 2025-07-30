package hft

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// HFTEngine provides high-frequency trading capabilities with microsecond latency
type HFTEngine struct {
	logger         *observability.Logger
	orderManager   *OrderManager
	portfolioMgr   *PortfolioManager
	riskManager    *RiskManager
	strategyEngine *StrategyEngine
	latencyTracker *LatencyTracker

	// Performance metrics
	ordersPerSecond int64
	latencyMicros   int64
	successRate     float64

	// Configuration
	config    HFTConfig
	isRunning int32
	stopChan  chan struct{}
	wg        sync.WaitGroup
	mu        sync.RWMutex

	// Real-time data
	marketData      chan MarketTick
	orderUpdates    chan OrderUpdate
	positionUpdates chan PositionUpdate
}

// HFTConfig contains configuration for the HFT engine
type HFTConfig struct {
	MaxOrdersPerSecond   int             `json:"max_orders_per_second"`
	LatencyTargetMicros  int64           `json:"latency_target_micros"`
	MaxPositionSize      decimal.Decimal `json:"max_position_size"`
	MaxDailyLoss         decimal.Decimal `json:"max_daily_loss"`
	RiskLimitPercent     float64         `json:"risk_limit_percent"`
	EnableMarketMaking   bool            `json:"enable_market_making"`
	EnableArbitrage      bool            `json:"enable_arbitrage"`
	TickerUpdateInterval time.Duration   `json:"ticker_update_interval"`
	OrderTimeoutMs       int64           `json:"order_timeout_ms"`
	MaxSlippageBps       int             `json:"max_slippage_bps"`
	MinProfitBps         int             `json:"min_profit_bps"`
}

// MarketTick represents real-time market data
type MarketTick struct {
	Symbol    string          `json:"symbol"`
	Price     decimal.Decimal `json:"price"`
	Volume    decimal.Decimal `json:"volume"`
	BidPrice  decimal.Decimal `json:"bid_price"`
	AskPrice  decimal.Decimal `json:"ask_price"`
	BidSize   decimal.Decimal `json:"bid_size"`
	AskSize   decimal.Decimal `json:"ask_size"`
	Timestamp time.Time       `json:"timestamp"`
	Exchange  string          `json:"exchange"`
	Sequence  uint64          `json:"sequence"`
}

// OrderUpdate represents order status changes
type OrderUpdate struct {
	OrderID       uuid.UUID       `json:"order_id"`
	Symbol        string          `json:"symbol"`
	Side          OrderSide       `json:"side"`
	Type          OrderType       `json:"type"`
	Status        OrderStatus     `json:"status"`
	Quantity      decimal.Decimal `json:"quantity"`
	Price         decimal.Decimal `json:"price"`
	FilledQty     decimal.Decimal `json:"filled_qty"`
	AvgFillPrice  decimal.Decimal `json:"avg_fill_price"`
	Timestamp     time.Time       `json:"timestamp"`
	LatencyMicros int64           `json:"latency_micros"`
}

// PositionUpdate represents position changes
type PositionUpdate struct {
	Symbol        string          `json:"symbol"`
	Size          decimal.Decimal `json:"size"`
	AvgPrice      decimal.Decimal `json:"avg_price"`
	UnrealizedPnL decimal.Decimal `json:"unrealized_pnl"`
	RealizedPnL   decimal.Decimal `json:"realized_pnl"`
	Timestamp     time.Time       `json:"timestamp"`
}

// OrderSide represents buy or sell
type OrderSide string

const (
	OrderSideBuy  OrderSide = "BUY"
	OrderSideSell OrderSide = "SELL"
)

// OrderType represents different order types
type OrderType string

const (
	OrderTypeMarket     OrderType = "MARKET"
	OrderTypeLimit      OrderType = "LIMIT"
	OrderTypeStopLoss   OrderType = "STOP_LOSS"
	OrderTypeTakeProfit OrderType = "TAKE_PROFIT"
	OrderTypeOCO        OrderType = "OCO"
)

// OrderStatus represents order status
type OrderStatus string

const (
	OrderStatusNew         OrderStatus = "NEW"
	OrderStatusPartialFill OrderStatus = "PARTIALLY_FILLED"
	OrderStatusFilled      OrderStatus = "FILLED"
	OrderStatusCanceled    OrderStatus = "CANCELED"
	OrderStatusRejected    OrderStatus = "REJECTED"
	OrderStatusExpired     OrderStatus = "EXPIRED"
)

// NewHFTEngine creates a new high-frequency trading engine
func NewHFTEngine(logger *observability.Logger, config HFTConfig) *HFTEngine {
	engine := &HFTEngine{
		logger:          logger,
		config:          config,
		stopChan:        make(chan struct{}),
		marketData:      make(chan MarketTick, 10000),
		orderUpdates:    make(chan OrderUpdate, 10000),
		positionUpdates: make(chan PositionUpdate, 1000),
	}

	// Initialize components
	engine.orderManager = NewOrderManager(logger, config)
	engine.portfolioMgr = NewPortfolioManager(logger, config)
	engine.riskManager = NewRiskManager(logger, config)
	engine.strategyEngine = NewStrategyEngine(logger, config)
	engine.latencyTracker = NewLatencyTracker(logger)

	return engine
}

// Start begins the HFT engine
func (e *HFTEngine) Start(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&e.isRunning, 0, 1) {
		return fmt.Errorf("HFT engine is already running")
	}

	e.logger.Info(ctx, "Starting HFT engine", map[string]interface{}{
		"max_orders_per_second": e.config.MaxOrdersPerSecond,
		"latency_target_micros": e.config.LatencyTargetMicros,
		"market_making_enabled": e.config.EnableMarketMaking,
		"arbitrage_enabled":     e.config.EnableArbitrage,
	})

	// Start core components
	if err := e.orderManager.Start(ctx); err != nil {
		return fmt.Errorf("failed to start order manager: %w", err)
	}

	if err := e.portfolioMgr.Start(ctx); err != nil {
		return fmt.Errorf("failed to start portfolio manager: %w", err)
	}

	if err := e.riskManager.Start(ctx); err != nil {
		return fmt.Errorf("failed to start risk manager: %w", err)
	}

	if err := e.strategyEngine.Start(ctx); err != nil {
		return fmt.Errorf("failed to start strategy engine: %w", err)
	}

	// Start processing loops
	e.wg.Add(4)
	go e.processMarketData(ctx)
	go e.processOrderUpdates(ctx)
	go e.processPositionUpdates(ctx)
	go e.performanceMonitor(ctx)

	e.logger.Info(ctx, "HFT engine started successfully", nil)
	return nil
}

// Stop gracefully shuts down the HFT engine
func (e *HFTEngine) Stop(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&e.isRunning, 1, 0) {
		return fmt.Errorf("HFT engine is not running")
	}

	e.logger.Info(ctx, "Stopping HFT engine", nil)

	close(e.stopChan)
	e.wg.Wait()

	// Stop components
	e.strategyEngine.Stop(ctx)
	e.riskManager.Stop(ctx)
	e.portfolioMgr.Stop(ctx)
	e.orderManager.Stop(ctx)

	e.logger.Info(ctx, "HFT engine stopped successfully", nil)
	return nil
}

// SubmitMarketData submits real-time market data to the engine
func (e *HFTEngine) SubmitMarketData(tick MarketTick) {
	select {
	case e.marketData <- tick:
	default:
		// Channel is full, drop the tick
		atomic.AddInt64(&e.latencyMicros, 1000) // Penalty for dropped data
	}
}

// processMarketData processes incoming market data with ultra-low latency
func (e *HFTEngine) processMarketData(ctx context.Context) {
	defer e.wg.Done()

	for {
		select {
		case <-e.stopChan:
			return
		case tick := <-e.marketData:
			start := time.Now()

			// Update portfolio with latest prices
			e.portfolioMgr.UpdatePrice(tick.Symbol, tick.Price)

			// Check risk limits
			if !e.riskManager.CheckLimits(tick.Symbol, tick.Price) {
				continue
			}

			// Generate trading signals
			signals := e.strategyEngine.ProcessTick(tick)

			// Execute signals
			for _, signal := range signals {
				if err := e.executeSignal(ctx, signal); err != nil {
					e.logger.Error(ctx, "Failed to execute signal", err)
				}
			}

			// Track latency
			latency := time.Since(start).Microseconds()
			e.latencyTracker.Record(latency)
			atomic.StoreInt64(&e.latencyMicros, latency)
		}
	}
}

// processOrderUpdates handles order status updates
func (e *HFTEngine) processOrderUpdates(ctx context.Context) {
	defer e.wg.Done()

	for {
		select {
		case <-e.stopChan:
			return
		case update := <-e.orderUpdates:
			e.orderManager.HandleUpdate(update)
			e.portfolioMgr.HandleOrderUpdate(update)

			// Update performance metrics
			atomic.AddInt64(&e.ordersPerSecond, 1)
		}
	}
}

// processPositionUpdates handles position changes
func (e *HFTEngine) processPositionUpdates(ctx context.Context) {
	defer e.wg.Done()

	for {
		select {
		case <-e.stopChan:
			return
		case update := <-e.positionUpdates:
			e.portfolioMgr.HandlePositionUpdate(update)
			e.riskManager.UpdatePosition(update)
		}
	}
}

// performanceMonitor tracks engine performance metrics
func (e *HFTEngine) performanceMonitor(ctx context.Context) {
	defer e.wg.Done()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var lastOrderCount int64

	for {
		select {
		case <-e.stopChan:
			return
		case <-ticker.C:
			currentOrders := atomic.LoadInt64(&e.ordersPerSecond)
			ordersThisSecond := currentOrders - lastOrderCount
			lastOrderCount = currentOrders

			latency := atomic.LoadInt64(&e.latencyMicros)

			e.logger.Info(ctx, "HFT performance metrics", map[string]interface{}{
				"orders_per_second":  ordersThisSecond,
				"avg_latency_micros": latency,
				"target_latency":     e.config.LatencyTargetMicros,
			})
		}
	}
}

// executeSignal executes a trading signal
func (e *HFTEngine) executeSignal(ctx context.Context, signal TradingSignal) error {
	// Validate signal
	if err := e.riskManager.ValidateSignal(signal); err != nil {
		return fmt.Errorf("signal validation failed: %w", err)
	}

	// Create order
	order := Order{
		ID:         uuid.New(),
		Symbol:     signal.Symbol,
		Side:       signal.Side,
		Type:       signal.OrderType,
		Quantity:   signal.Quantity,
		Price:      signal.Price,
		Status:     OrderStatusNew,
		CreatedAt:  time.Now(),
		StrategyID: signal.StrategyID,
	}

	// Submit order
	return e.orderManager.SubmitOrder(ctx, order)
}

// GetMetrics returns current engine metrics
func (e *HFTEngine) GetMetrics() HFTMetrics {
	return HFTMetrics{
		OrdersPerSecond:  atomic.LoadInt64(&e.ordersPerSecond),
		AvgLatencyMicros: atomic.LoadInt64(&e.latencyMicros),
		SuccessRate:      e.successRate,
		IsRunning:        atomic.LoadInt32(&e.isRunning) == 1,
		ActiveStrategies: e.strategyEngine.GetActiveCount(),
		OpenPositions:    e.portfolioMgr.GetPositionCount(),
		TotalPnL:         e.portfolioMgr.GetTotalPnL(),
	}
}

// HFTMetrics contains engine performance metrics
type HFTMetrics struct {
	OrdersPerSecond  int64           `json:"orders_per_second"`
	AvgLatencyMicros int64           `json:"avg_latency_micros"`
	SuccessRate      float64         `json:"success_rate"`
	IsRunning        bool            `json:"is_running"`
	ActiveStrategies int             `json:"active_strategies"`
	OpenPositions    int             `json:"open_positions"`
	TotalPnL         decimal.Decimal `json:"total_pnl"`
}

// StrategyEngine placeholder - will be implemented in strategy_engine.go
type StrategyEngine struct {
	logger *observability.Logger
	config HFTConfig
}

func NewStrategyEngine(logger *observability.Logger, config HFTConfig) *StrategyEngine {
	return &StrategyEngine{logger: logger, config: config}
}

func (se *StrategyEngine) Start(ctx context.Context) error { return nil }
func (se *StrategyEngine) Stop(ctx context.Context) error  { return nil }
func (se *StrategyEngine) ProcessTick(tick MarketTick) []TradingSignal {
	// Simple momentum strategy placeholder
	var signals []TradingSignal

	// Generate a simple momentum signal based on price movement
	signal := TradingSignal{
		ID:         uuid.New(),
		Symbol:     tick.Symbol,
		Side:       OrderSideBuy, // Default to buy
		OrderType:  OrderTypeMarket,
		Quantity:   decimal.NewFromFloat(0.01),
		Price:      tick.Price,
		Confidence: 0.6,
		StrategyID: "hft_momentum",
		Timestamp:  time.Now(),
		Metadata: map[string]interface{}{
			"strategy": "momentum",
			"price":    tick.Price.String(),
		},
	}

	signals = append(signals, signal)
	return signals
}
func (se *StrategyEngine) GetActiveCount() int { return 1 }

// LatencyTracker placeholder - will be implemented in latency_tracker.go
type LatencyTracker struct {
	logger *observability.Logger
}

func NewLatencyTracker(logger *observability.Logger) *LatencyTracker {
	return &LatencyTracker{logger: logger}
}

func (lt *LatencyTracker) Record(latencyMicros int64) {}

// IsRunning returns whether the engine is currently running
func (e *HFTEngine) IsRunning() bool {
	return atomic.LoadInt32(&e.isRunning) == 1
}

// GetUptime returns the engine uptime
func (e *HFTEngine) GetUptime() time.Duration {
	// Mock implementation - in real system would track start time
	return 2*time.Hour + 30*time.Minute
}

// GetConfig returns the current engine configuration
func (e *HFTEngine) GetConfig() HFTConfig {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.config
}

// UpdateConfig updates the engine configuration
func (e *HFTEngine) UpdateConfig(config HFTConfig) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.config = config
	return nil
}
