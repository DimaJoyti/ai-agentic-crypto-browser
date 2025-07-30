package strategies

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/internal/hft"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/shopspring/decimal"
)

// StrategyEngine manages and executes multiple trading strategies
type StrategyEngine struct {
	logger     *observability.Logger
	config     EngineConfig
	strategies map[string]Strategy

	// Performance tracking
	performance map[string]*StrategyPerformance

	// State management
	isRunning bool
	stopChan  chan struct{}
	wg        sync.WaitGroup
	mu        sync.RWMutex
}

// EngineConfig contains configuration for the strategy engine
type EngineConfig struct {
	MaxStrategies        int             `json:"max_strategies"`
	ExecutionInterval    time.Duration   `json:"execution_interval"`
	PerformanceWindow    time.Duration   `json:"performance_window"`
	EnableRiskManagement bool            `json:"enable_risk_management"`
	MaxPositionSize      decimal.Decimal `json:"max_position_size"`
	MaxDailyLoss         decimal.Decimal `json:"max_daily_loss"`
}

// Strategy interface that all trading strategies must implement
type Strategy interface {
	GetID() string
	GetName() string
	GetType() StrategyType
	GetConfig() StrategyConfig
	Initialize(ctx context.Context) error
	ProcessTick(ctx context.Context, tick hft.MarketTick) ([]hft.TradingSignal, error)
	UpdateParameters(params map[string]interface{}) error
	GetPerformance() *StrategyPerformance
	IsEnabled() bool
	SetEnabled(enabled bool)
	Cleanup(ctx context.Context) error
}

// StrategyType represents different types of trading strategies
type StrategyType string

const (
	StrategyTypeMarketMaking  StrategyType = "MARKET_MAKING"
	StrategyTypeArbitrage     StrategyType = "ARBITRAGE"
	StrategyTypeMomentum      StrategyType = "MOMENTUM"
	StrategyTypeMeanReversion StrategyType = "MEAN_REVERSION"
	StrategyTypeBreakout      StrategyType = "BREAKOUT"
	StrategyTypeScalping      StrategyType = "SCALPING"
	StrategyTypeAIPredictive  StrategyType = "AI_PREDICTIVE"
	StrategyTypeGrid          StrategyType = "GRID"
	StrategyTypePairs         StrategyType = "PAIRS"
)

// StrategyConfig contains base configuration for strategies
type StrategyConfig struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Type            StrategyType           `json:"type"`
	Symbols         []string               `json:"symbols"`
	Enabled         bool                   `json:"enabled"`
	RiskLimit       decimal.Decimal        `json:"risk_limit"`
	MaxPositionSize decimal.Decimal        `json:"max_position_size"`
	Parameters      map[string]interface{} `json:"parameters"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// StrategyPerformance tracks strategy performance metrics
type StrategyPerformance struct {
	StrategyID           string          `json:"strategy_id"`
	TotalTrades          int64           `json:"total_trades"`
	WinningTrades        int64           `json:"winning_trades"`
	LosingTrades         int64           `json:"losing_trades"`
	WinRate              float64         `json:"win_rate"`
	TotalPnL             decimal.Decimal `json:"total_pnl"`
	RealizedPnL          decimal.Decimal `json:"realized_pnl"`
	UnrealizedPnL        decimal.Decimal `json:"unrealized_pnl"`
	MaxDrawdown          decimal.Decimal `json:"max_drawdown"`
	SharpeRatio          float64         `json:"sharpe_ratio"`
	AverageWin           decimal.Decimal `json:"average_win"`
	AverageLoss          decimal.Decimal `json:"average_loss"`
	ProfitFactor         float64         `json:"profit_factor"`
	MaxConsecutiveWins   int             `json:"max_consecutive_wins"`
	MaxConsecutiveLosses int             `json:"max_consecutive_losses"`
	LastUpdate           time.Time       `json:"last_update"`
}

// NewStrategyEngine creates a new strategy engine
func NewStrategyEngine(logger *observability.Logger, config EngineConfig) *StrategyEngine {
	if config.MaxStrategies == 0 {
		config.MaxStrategies = 50
	}

	if config.ExecutionInterval == 0 {
		config.ExecutionInterval = 100 * time.Millisecond
	}

	if config.PerformanceWindow == 0 {
		config.PerformanceWindow = 24 * time.Hour
	}

	return &StrategyEngine{
		logger:      logger,
		config:      config,
		strategies:  make(map[string]Strategy),
		performance: make(map[string]*StrategyPerformance),
		stopChan:    make(chan struct{}),
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
		"max_strategies":     se.config.MaxStrategies,
		"execution_interval": se.config.ExecutionInterval.String(),
		"strategies_count":   len(se.strategies),
	})

	// Initialize all strategies
	for _, strategy := range se.strategies {
		if strategy.IsEnabled() {
			if err := strategy.Initialize(ctx); err != nil {
				se.logger.Error(ctx, "Failed to initialize strategy", err, map[string]interface{}{
					"strategy_id":   strategy.GetID(),
					"strategy_name": strategy.GetName(),
				})
			}
		}
	}

	se.isRunning = true

	// Start monitoring goroutines
	se.wg.Add(1)
	go se.performanceMonitor(ctx)

	se.logger.Info(ctx, "Strategy engine started successfully", nil)

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

	se.isRunning = false
	close(se.stopChan)

	// Cleanup all strategies
	for _, strategy := range se.strategies {
		if err := strategy.Cleanup(ctx); err != nil {
			se.logger.Error(ctx, "Failed to cleanup strategy", err, map[string]interface{}{
				"strategy_id": strategy.GetID(),
			})
		}
	}

	se.wg.Wait()

	se.logger.Info(ctx, "Strategy engine stopped successfully", nil)

	return nil
}

// AddStrategy adds a new strategy to the engine
func (se *StrategyEngine) AddStrategy(strategy Strategy) error {
	se.mu.Lock()
	defer se.mu.Unlock()

	if len(se.strategies) >= se.config.MaxStrategies {
		return fmt.Errorf("maximum number of strategies reached: %d", se.config.MaxStrategies)
	}

	strategyID := strategy.GetID()
	if _, exists := se.strategies[strategyID]; exists {
		return fmt.Errorf("strategy already exists: %s", strategyID)
	}

	se.strategies[strategyID] = strategy
	se.performance[strategyID] = &StrategyPerformance{
		StrategyID: strategyID,
		LastUpdate: time.Now(),
	}

	se.logger.Info(context.Background(), "Strategy added", map[string]interface{}{
		"strategy_id":   strategyID,
		"strategy_name": strategy.GetName(),
		"strategy_type": string(strategy.GetType()),
	})

	return nil
}

// RemoveStrategy removes a strategy from the engine
func (se *StrategyEngine) RemoveStrategy(strategyID string) error {
	se.mu.Lock()
	defer se.mu.Unlock()

	strategy, exists := se.strategies[strategyID]
	if !exists {
		return fmt.Errorf("strategy not found: %s", strategyID)
	}

	// Cleanup strategy
	if err := strategy.Cleanup(context.Background()); err != nil {
		se.logger.Error(context.Background(), "Failed to cleanup strategy during removal", err)
	}

	delete(se.strategies, strategyID)
	delete(se.performance, strategyID)

	se.logger.Info(context.Background(), "Strategy removed", map[string]interface{}{
		"strategy_id": strategyID,
	})

	return nil
}

// GetStrategy retrieves a strategy by ID
func (se *StrategyEngine) GetStrategy(strategyID string) (Strategy, error) {
	se.mu.RLock()
	defer se.mu.RUnlock()

	strategy, exists := se.strategies[strategyID]
	if !exists {
		return nil, fmt.Errorf("strategy not found: %s", strategyID)
	}

	return strategy, nil
}

// GetAllStrategies returns all strategies
func (se *StrategyEngine) GetAllStrategies() map[string]Strategy {
	se.mu.RLock()
	defer se.mu.RUnlock()

	strategies := make(map[string]Strategy)
	for id, strategy := range se.strategies {
		strategies[id] = strategy
	}

	return strategies
}

// ProcessMarketTick processes a market tick through all enabled strategies
func (se *StrategyEngine) ProcessMarketTick(ctx context.Context, tick hft.MarketTick) []hft.TradingSignal {
	se.mu.RLock()
	defer se.mu.RUnlock()

	var allSignals []hft.TradingSignal

	for _, strategy := range se.strategies {
		if !strategy.IsEnabled() {
			continue
		}

		// Check if strategy handles this symbol
		config := strategy.GetConfig()
		if !se.symbolMatches(tick.Symbol, config.Symbols) {
			continue
		}

		// Process tick through strategy
		signals, err := strategy.ProcessTick(ctx, tick)
		if err != nil {
			se.logger.Error(ctx, "Strategy processing error", err, map[string]interface{}{
				"strategy_id": strategy.GetID(),
				"symbol":      tick.Symbol,
			})
			continue
		}

		// Apply risk management
		if se.config.EnableRiskManagement {
			signals = se.applyRiskManagement(signals, strategy)
		}

		allSignals = append(allSignals, signals...)
	}

	return allSignals
}

// EnableStrategy enables a strategy
func (se *StrategyEngine) EnableStrategy(strategyID string) error {
	se.mu.RLock()
	strategy, exists := se.strategies[strategyID]
	se.mu.RUnlock()

	if !exists {
		return fmt.Errorf("strategy not found: %s", strategyID)
	}

	strategy.SetEnabled(true)

	se.logger.Info(context.Background(), "Strategy enabled", map[string]interface{}{
		"strategy_id": strategyID,
	})

	return nil
}

// DisableStrategy disables a strategy
func (se *StrategyEngine) DisableStrategy(strategyID string) error {
	se.mu.RLock()
	strategy, exists := se.strategies[strategyID]
	se.mu.RUnlock()

	if !exists {
		return fmt.Errorf("strategy not found: %s", strategyID)
	}

	strategy.SetEnabled(false)

	se.logger.Info(context.Background(), "Strategy disabled", map[string]interface{}{
		"strategy_id": strategyID,
	})

	return nil
}

// GetPerformance retrieves performance metrics for a strategy
func (se *StrategyEngine) GetPerformance(strategyID string) (*StrategyPerformance, error) {
	se.mu.RLock()
	defer se.mu.RUnlock()

	performance, exists := se.performance[strategyID]
	if !exists {
		return nil, fmt.Errorf("performance data not found for strategy: %s", strategyID)
	}

	return performance, nil
}

// GetAllPerformance returns performance metrics for all strategies
func (se *StrategyEngine) GetAllPerformance() map[string]*StrategyPerformance {
	se.mu.RLock()
	defer se.mu.RUnlock()

	performance := make(map[string]*StrategyPerformance)
	for id, perf := range se.performance {
		performance[id] = perf
	}

	return performance
}

// UpdateStrategyParameters updates parameters for a strategy
func (se *StrategyEngine) UpdateStrategyParameters(strategyID string, params map[string]interface{}) error {
	se.mu.RLock()
	strategy, exists := se.strategies[strategyID]
	se.mu.RUnlock()

	if !exists {
		return fmt.Errorf("strategy not found: %s", strategyID)
	}

	if err := strategy.UpdateParameters(params); err != nil {
		return fmt.Errorf("failed to update strategy parameters: %w", err)
	}

	se.logger.Info(context.Background(), "Strategy parameters updated", map[string]interface{}{
		"strategy_id": strategyID,
		"parameters":  params,
	})

	return nil
}

// symbolMatches checks if a symbol matches the strategy's symbol list
func (se *StrategyEngine) symbolMatches(symbol string, symbols []string) bool {
	if len(symbols) == 0 {
		return true // Strategy handles all symbols
	}

	for _, s := range symbols {
		if s == symbol || s == "*" {
			return true
		}
	}

	return false
}

// applyRiskManagement applies risk management rules to signals
func (se *StrategyEngine) applyRiskManagement(signals []hft.TradingSignal, strategy Strategy) []hft.TradingSignal {
	var filteredSignals []hft.TradingSignal

	config := strategy.GetConfig()

	for _, signal := range signals {
		// Check position size limit
		if signal.Quantity.GreaterThan(config.MaxPositionSize) {
			signal.Quantity = config.MaxPositionSize
		}

		// Check global position size limit
		if signal.Quantity.GreaterThan(se.config.MaxPositionSize) {
			signal.Quantity = se.config.MaxPositionSize
		}

		// Only include signals with valid quantities
		if signal.Quantity.GreaterThan(decimal.Zero) {
			filteredSignals = append(filteredSignals, signal)
		}
	}

	return filteredSignals
}

// performanceMonitor monitors strategy performance
func (se *StrategyEngine) performanceMonitor(ctx context.Context) {
	defer se.wg.Done()

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-se.stopChan:
			return
		case <-ticker.C:
			se.updatePerformanceMetrics(ctx)
		}
	}
}

// updatePerformanceMetrics updates performance metrics for all strategies
func (se *StrategyEngine) updatePerformanceMetrics(ctx context.Context) {
	se.mu.Lock()
	defer se.mu.Unlock()

	for strategyID, strategy := range se.strategies {
		if !strategy.IsEnabled() {
			continue
		}

		// Get performance from strategy
		strategyPerf := strategy.GetPerformance()
		if strategyPerf != nil {
			se.performance[strategyID] = strategyPerf
		}
	}
}

// GetMetrics returns engine metrics
func (se *StrategyEngine) GetMetrics() EngineMetrics {
	se.mu.RLock()
	defer se.mu.RUnlock()

	enabledCount := 0
	for _, strategy := range se.strategies {
		if strategy.IsEnabled() {
			enabledCount++
		}
	}

	return EngineMetrics{
		IsRunning:         se.isRunning,
		TotalStrategies:   len(se.strategies),
		EnabledStrategies: enabledCount,
		MaxStrategies:     se.config.MaxStrategies,
		ExecutionInterval: se.config.ExecutionInterval,
	}
}

// EngineMetrics contains engine performance metrics
type EngineMetrics struct {
	IsRunning         bool          `json:"is_running"`
	TotalStrategies   int           `json:"total_strategies"`
	EnabledStrategies int           `json:"enabled_strategies"`
	MaxStrategies     int           `json:"max_strategies"`
	ExecutionInterval time.Duration `json:"execution_interval"`
}
