package framework

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/internal/exchanges"
	"github.com/ai-agentic-browser/internal/exchanges/common"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// StrategyEngine manages multiple trading strategies
type StrategyEngine struct {
	logger          *observability.Logger
	exchangeManager *exchanges.Manager
	orderManager    *exchanges.OrderManager
	strategies      map[uuid.UUID]Strategy
	signalChannel   chan *Signal
	config          EngineConfig

	// Performance tracking
	totalSignals    int64
	executedSignals int64
	rejectedSignals int64

	// State management
	isRunning bool
	stopChan  chan struct{}
	wg        sync.WaitGroup
	mu        sync.RWMutex
}

// EngineConfig contains strategy engine configuration
type EngineConfig struct {
	MaxStrategies      int           `json:"max_strategies"`
	SignalBufferSize   int           `json:"signal_buffer_size"`
	ProcessingInterval time.Duration `json:"processing_interval"`
	EnableRiskChecks   bool          `json:"enable_risk_checks"`
	EnablePaperTrading bool          `json:"enable_paper_trading"`
	LogLevel           string        `json:"log_level"`
}

// NewStrategyEngine creates a new strategy engine
func NewStrategyEngine(
	logger *observability.Logger,
	exchangeManager *exchanges.Manager,
	orderManager *exchanges.OrderManager,
	config EngineConfig,
) *StrategyEngine {
	if config.MaxStrategies == 0 {
		config.MaxStrategies = 100
	}
	if config.SignalBufferSize == 0 {
		config.SignalBufferSize = 10000
	}
	if config.ProcessingInterval == 0 {
		config.ProcessingInterval = 100 * time.Millisecond
	}

	return &StrategyEngine{
		logger:          logger,
		exchangeManager: exchangeManager,
		orderManager:    orderManager,
		strategies:      make(map[uuid.UUID]Strategy),
		signalChannel:   make(chan *Signal, config.SignalBufferSize),
		config:          config,
		stopChan:        make(chan struct{}),
	}
}

// Start begins the strategy engine
func (se *StrategyEngine) Start(ctx context.Context) error {
	se.mu.Lock()
	defer se.mu.Unlock()

	if se.isRunning {
		return fmt.Errorf("strategy engine is already running")
	}

	se.logger.Info(ctx, "Starting strategy engine", map[string]interface{}{
		"max_strategies":       se.config.MaxStrategies,
		"signal_buffer_size":   se.config.SignalBufferSize,
		"enable_paper_trading": se.config.EnablePaperTrading,
		"enable_risk_checks":   se.config.EnableRiskChecks,
	})

	se.isRunning = true

	// Start signal processing
	se.wg.Add(1)
	go se.processSignals(ctx)

	// Start market data distribution
	se.wg.Add(1)
	go se.distributeMarketData(ctx)

	se.logger.Info(ctx, "Strategy engine started successfully", map[string]interface{}{
		"active_strategies": len(se.strategies),
	})

	return nil
}

// Stop gracefully shuts down the strategy engine
func (se *StrategyEngine) Stop(ctx context.Context) error {
	se.mu.Lock()
	defer se.mu.Unlock()

	if !se.isRunning {
		return fmt.Errorf("strategy engine is not running")
	}

	se.logger.Info(ctx, "Stopping strategy engine", nil)

	// Stop all strategies
	for id, strategy := range se.strategies {
		if strategy.IsRunning() {
			if err := strategy.Stop(ctx); err != nil {
				se.logger.Error(ctx, "Failed to stop strategy", err, map[string]interface{}{
					"strategy_id": id.String(),
				})
			}
		}
	}

	close(se.stopChan)
	se.wg.Wait()

	se.isRunning = false

	se.logger.Info(ctx, "Strategy engine stopped successfully", nil)

	return nil
}

// RegisterStrategy registers a new strategy
func (se *StrategyEngine) RegisterStrategy(ctx context.Context, strategy Strategy) error {
	se.mu.Lock()
	defer se.mu.Unlock()

	if len(se.strategies) >= se.config.MaxStrategies {
		return fmt.Errorf("maximum number of strategies reached: %d", se.config.MaxStrategies)
	}

	strategyID := strategy.GetID()
	if _, exists := se.strategies[strategyID]; exists {
		return fmt.Errorf("strategy already registered: %s", strategyID.String())
	}

	se.strategies[strategyID] = strategy

	se.logger.Info(ctx, "Strategy registered", map[string]interface{}{
		"strategy_id":      strategyID.String(),
		"strategy_name":    strategy.GetName(),
		"total_strategies": len(se.strategies),
	})

	return nil
}

// UnregisterStrategy removes a strategy
func (se *StrategyEngine) UnregisterStrategy(ctx context.Context, strategyID uuid.UUID) error {
	se.mu.Lock()
	defer se.mu.Unlock()

	strategy, exists := se.strategies[strategyID]
	if !exists {
		return fmt.Errorf("strategy not found: %s", strategyID.String())
	}

	// Stop strategy if running
	if strategy.IsRunning() {
		if err := strategy.Stop(ctx); err != nil {
			se.logger.Error(ctx, "Failed to stop strategy during unregistration", err, map[string]interface{}{
				"strategy_id": strategyID.String(),
			})
		}
	}

	delete(se.strategies, strategyID)

	se.logger.Info(ctx, "Strategy unregistered", map[string]interface{}{
		"strategy_id":          strategyID.String(),
		"remaining_strategies": len(se.strategies),
	})

	return nil
}

// StartStrategy starts a specific strategy
func (se *StrategyEngine) StartStrategy(ctx context.Context, strategyID uuid.UUID) error {
	se.mu.RLock()
	strategy, exists := se.strategies[strategyID]
	se.mu.RUnlock()

	if !exists {
		return fmt.Errorf("strategy not found: %s", strategyID.String())
	}

	if strategy.IsRunning() {
		return fmt.Errorf("strategy is already running: %s", strategyID.String())
	}

	if err := strategy.Start(ctx); err != nil {
		return fmt.Errorf("failed to start strategy: %w", err)
	}

	se.logger.Info(ctx, "Strategy started", map[string]interface{}{
		"strategy_id":   strategyID.String(),
		"strategy_name": strategy.GetName(),
	})

	return nil
}

// StopStrategy stops a specific strategy
func (se *StrategyEngine) StopStrategy(ctx context.Context, strategyID uuid.UUID) error {
	se.mu.RLock()
	strategy, exists := se.strategies[strategyID]
	se.mu.RUnlock()

	if !exists {
		return fmt.Errorf("strategy not found: %s", strategyID.String())
	}

	if !strategy.IsRunning() {
		return fmt.Errorf("strategy is not running: %s", strategyID.String())
	}

	if err := strategy.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop strategy: %w", err)
	}

	se.logger.Info(ctx, "Strategy stopped", map[string]interface{}{
		"strategy_id":   strategyID.String(),
		"strategy_name": strategy.GetName(),
	})

	return nil
}

// GetStrategy returns a strategy by ID
func (se *StrategyEngine) GetStrategy(strategyID uuid.UUID) (Strategy, error) {
	se.mu.RLock()
	defer se.mu.RUnlock()

	strategy, exists := se.strategies[strategyID]
	if !exists {
		return nil, fmt.Errorf("strategy not found: %s", strategyID.String())
	}

	return strategy, nil
}

// GetAllStrategies returns all registered strategies
func (se *StrategyEngine) GetAllStrategies() []Strategy {
	se.mu.RLock()
	defer se.mu.RUnlock()

	strategies := make([]Strategy, 0, len(se.strategies))
	for _, strategy := range se.strategies {
		strategies = append(strategies, strategy)
	}

	return strategies
}

// SubmitSignal submits a trading signal for processing
func (se *StrategyEngine) SubmitSignal(signal *Signal) error {
	select {
	case se.signalChannel <- signal:
		se.totalSignals++
		return nil
	default:
		se.rejectedSignals++
		return fmt.Errorf("signal channel is full")
	}
}

// GetMetrics returns engine performance metrics
func (se *StrategyEngine) GetMetrics() *EngineMetrics {
	se.mu.RLock()
	defer se.mu.RUnlock()

	activeStrategies := 0
	for _, strategy := range se.strategies {
		if strategy.IsRunning() {
			activeStrategies++
		}
	}

	return &EngineMetrics{
		TotalStrategies:   len(se.strategies),
		ActiveStrategies:  activeStrategies,
		TotalSignals:      se.totalSignals,
		ExecutedSignals:   se.executedSignals,
		RejectedSignals:   se.rejectedSignals,
		SignalBufferUsage: len(se.signalChannel),
		IsRunning:         se.isRunning,
	}
}

// EngineMetrics contains strategy engine performance metrics
type EngineMetrics struct {
	TotalStrategies   int   `json:"total_strategies"`
	ActiveStrategies  int   `json:"active_strategies"`
	TotalSignals      int64 `json:"total_signals"`
	ExecutedSignals   int64 `json:"executed_signals"`
	RejectedSignals   int64 `json:"rejected_signals"`
	SignalBufferUsage int   `json:"signal_buffer_usage"`
	IsRunning         bool  `json:"is_running"`
}

// Private methods

// processSignals processes trading signals from strategies
func (se *StrategyEngine) processSignals(ctx context.Context) {
	defer se.wg.Done()

	ticker := time.NewTicker(se.config.ProcessingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-se.stopChan:
			return
		case signal := <-se.signalChannel:
			se.handleSignal(ctx, signal)
		case <-ticker.C:
			// Periodic processing if needed
		}
	}
}

// handleSignal processes a single trading signal
func (se *StrategyEngine) handleSignal(ctx context.Context, signal *Signal) {
	// Validate signal
	if err := se.validateSignal(signal); err != nil {
		se.logger.Error(ctx, "Invalid signal", err, map[string]interface{}{
			"signal_id":   signal.ID.String(),
			"strategy_id": signal.StrategyID.String(),
		})
		se.rejectedSignals++
		return
	}

	// Apply risk checks if enabled
	if se.config.EnableRiskChecks {
		if err := se.checkRiskLimits(ctx, signal); err != nil {
			se.logger.Warn(ctx, "Signal rejected by risk management", map[string]interface{}{
				"signal_id":   signal.ID.String(),
				"strategy_id": signal.StrategyID.String(),
				"reason":      err.Error(),
			})
			se.rejectedSignals++
			return
		}
	}

	// Execute signal
	if err := se.executeSignal(ctx, signal); err != nil {
		se.logger.Error(ctx, "Failed to execute signal", err, map[string]interface{}{
			"signal_id":   signal.ID.String(),
			"strategy_id": signal.StrategyID.String(),
		})
		se.rejectedSignals++
		return
	}

	se.executedSignals++
}

// validateSignal validates a trading signal
func (se *StrategyEngine) validateSignal(signal *Signal) error {
	if signal.Symbol == "" {
		return fmt.Errorf("symbol is required")
	}
	if signal.Quantity.IsZero() || signal.Quantity.IsNegative() {
		return fmt.Errorf("quantity must be positive")
	}
	if signal.Strength.LessThan(decimal.NewFromFloat(0)) || signal.Strength.GreaterThan(decimal.NewFromFloat(1)) {
		return fmt.Errorf("strength must be between 0 and 1")
	}
	if signal.Confidence.LessThan(decimal.NewFromFloat(0)) || signal.Confidence.GreaterThan(decimal.NewFromFloat(1)) {
		return fmt.Errorf("confidence must be between 0 and 1")
	}
	return nil
}

// checkRiskLimits applies risk management checks to a signal
func (se *StrategyEngine) checkRiskLimits(ctx context.Context, signal *Signal) error {
	strategy, err := se.GetStrategy(signal.StrategyID)
	if err != nil {
		return fmt.Errorf("strategy not found: %w", err)
	}

	limits := strategy.GetRiskLimits()
	if limits == nil {
		return nil // No limits defined
	}

	// Check position size limit
	if signal.Quantity.GreaterThan(limits.MaxPositionSize) {
		return fmt.Errorf("quantity exceeds max position size")
	}

	// Check allowed symbols
	if len(limits.AllowedSymbols) > 0 {
		allowed := false
		for _, allowedSymbol := range limits.AllowedSymbols {
			if signal.Symbol == allowedSymbol {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("symbol not in allowed list")
		}
	}

	// Check blocked symbols
	for _, blockedSymbol := range limits.BlockedSymbols {
		if signal.Symbol == blockedSymbol {
			return fmt.Errorf("symbol is blocked")
		}
	}

	// Check trading hours
	if limits.TradingHours != nil && limits.TradingHours.Enabled {
		if !se.isWithinTradingHours(limits.TradingHours) {
			return fmt.Errorf("outside trading hours")
		}
	}

	return nil
}

// executeSignal executes a trading signal
func (se *StrategyEngine) executeSignal(ctx context.Context, signal *Signal) error {
	if se.config.EnablePaperTrading {
		return se.executePaperSignal(ctx, signal)
	}

	return se.executeLiveSignal(ctx, signal)
}

// executePaperSignal executes a signal in paper trading mode
func (se *StrategyEngine) executePaperSignal(ctx context.Context, signal *Signal) error {
	// TODO: Implement paper trading execution
	se.logger.Info(ctx, "Paper trading signal executed", map[string]interface{}{
		"signal_id":   signal.ID.String(),
		"strategy_id": signal.StrategyID.String(),
		"symbol":      signal.Symbol,
		"side":        string(signal.Side),
		"quantity":    signal.Quantity.String(),
	})
	return nil
}

// executeLiveSignal executes a signal in live trading mode
func (se *StrategyEngine) executeLiveSignal(ctx context.Context, signal *Signal) error {
	// Convert signal to order request
	orderReq := se.signalToOrderRequest(signal)

	// Submit order through order manager
	_, err := se.orderManager.SubmitOrder(ctx, orderReq)
	if err != nil {
		return fmt.Errorf("failed to submit order: %w", err)
	}

	se.logger.Info(ctx, "Live trading signal executed", map[string]interface{}{
		"signal_id":   signal.ID.String(),
		"strategy_id": signal.StrategyID.String(),
		"symbol":      signal.Symbol,
		"side":        string(signal.Side),
		"quantity":    signal.Quantity.String(),
	})

	return nil
}

// signalToOrderRequest converts a signal to an order request
func (se *StrategyEngine) signalToOrderRequest(signal *Signal) *common.OrderRequest {
	return &common.OrderRequest{
		Symbol:      signal.Symbol,
		Side:        signal.Side,
		Type:        common.OrderTypeLimit, // Default to limit orders
		Quantity:    signal.Quantity,
		Price:       signal.Price,
		TimeInForce: signal.TimeInForce,
		Metadata: map[string]interface{}{
			"strategy_id": signal.StrategyID.String(),
			"signal_id":   signal.ID.String(),
		},
	}
}

// distributeMarketData distributes market data to all running strategies
func (se *StrategyEngine) distributeMarketData(ctx context.Context) {
	defer se.wg.Done()

	// TODO: Subscribe to market data from exchange manager
	// and distribute to all running strategies

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-se.stopChan:
			return
		case <-ticker.C:
			// Placeholder for market data distribution
		}
	}
}

// isWithinTradingHours checks if current time is within trading hours
func (se *StrategyEngine) isWithinTradingHours(hours *TradingHours) bool {
	now := time.Now()

	// Check weekday
	weekday := int(now.Weekday())
	allowed := false
	for _, allowedDay := range hours.Weekdays {
		if weekday == allowedDay {
			allowed = true
			break
		}
	}
	if !allowed {
		return false
	}

	// Check time range
	currentTime := now.Format("15:04:05")
	startTime := hours.StartTime.Format("15:04:05")
	endTime := hours.EndTime.Format("15:04:05")

	return currentTime >= startTime && currentTime <= endTime
}
