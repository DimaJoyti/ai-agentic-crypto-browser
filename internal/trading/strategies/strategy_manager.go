package strategies

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/shopspring/decimal"
)

// StrategyManager manages all 7 trading bot strategies
type StrategyManager struct {
	logger     *observability.Logger
	strategies map[string]TradingStrategy
	mu         sync.RWMutex
}

// TradingStrategy interface that all strategies must implement
type TradingStrategy interface {
	GetName() string
	GetType() string
	Execute(ctx context.Context, marketData *MarketData) (interface{}, error)
	Validate() error
	Reset()
	GetPerformance() *StrategyPerformance
}

// StrategyType represents the type of trading strategy
type StrategyType string

const (
	StrategyTypeDCA           StrategyType = "dca"
	StrategyTypeGrid          StrategyType = "grid"
	StrategyTypeMomentum      StrategyType = "momentum"
	StrategyTypeMeanReversion StrategyType = "mean_reversion"
	StrategyTypeArbitrage     StrategyType = "arbitrage"
	StrategyTypeScalping      StrategyType = "scalping"
	StrategyTypeSwing         StrategyType = "swing"
)

// NewStrategyManager creates a new strategy manager
func NewStrategyManager(logger *observability.Logger) *StrategyManager {
	return &StrategyManager{
		logger:     logger,
		strategies: make(map[string]TradingStrategy),
	}
}

// RegisterStrategy registers a new trading strategy
func (sm *StrategyManager) RegisterStrategy(id string, strategy TradingStrategy) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if err := strategy.Validate(); err != nil {
		return fmt.Errorf("strategy validation failed: %w", err)
	}

	sm.strategies[id] = strategy

	sm.logger.Info(context.Background(), "Strategy registered", map[string]interface{}{
		"strategy_id":   id,
		"strategy_name": strategy.GetName(),
		"strategy_type": strategy.GetType(),
	})

	return nil
}

// GetStrategy retrieves a strategy by ID
func (sm *StrategyManager) GetStrategy(id string) (TradingStrategy, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	strategy, exists := sm.strategies[id]
	if !exists {
		return nil, fmt.Errorf("strategy not found: %s", id)
	}

	return strategy, nil
}

// ListStrategies returns all registered strategies
func (sm *StrategyManager) ListStrategies() map[string]TradingStrategy {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	strategies := make(map[string]TradingStrategy)
	for id, strategy := range sm.strategies {
		strategies[id] = strategy
	}

	return strategies
}

// ExecuteStrategy executes a specific strategy
func (sm *StrategyManager) ExecuteStrategy(ctx context.Context, strategyID string, marketData *MarketData) (interface{}, error) {
	strategy, err := sm.GetStrategy(strategyID)
	if err != nil {
		return nil, err
	}

	return strategy.Execute(ctx, marketData)
}

// ExecuteAllStrategies executes all registered strategies
func (sm *StrategyManager) ExecuteAllStrategies(ctx context.Context, marketData *MarketData) map[string]interface{} {
	sm.mu.RLock()
	strategies := make(map[string]TradingStrategy)
	for id, strategy := range sm.strategies {
		strategies[id] = strategy
	}
	sm.mu.RUnlock()

	results := make(map[string]interface{})
	var wg sync.WaitGroup

	for id, strategy := range strategies {
		wg.Add(1)
		go func(strategyID string, strat TradingStrategy) {
			defer wg.Done()
			
			result, err := strat.Execute(ctx, marketData)
			if err != nil {
				sm.logger.Error(ctx, "Strategy execution failed", err, map[string]interface{}{
					"strategy_id": strategyID,
					"symbol":      marketData.Symbol,
				})
				results[strategyID] = map[string]interface{}{
					"error": err.Error(),
				}
			} else {
				results[strategyID] = result
			}
		}(id, strategy)
	}

	wg.Wait()
	return results
}

// GetPerformanceMetrics returns performance metrics for all strategies
func (sm *StrategyManager) GetPerformanceMetrics() map[string]*StrategyPerformance {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	metrics := make(map[string]*StrategyPerformance)
	for id, strategy := range sm.strategies {
		metrics[id] = strategy.GetPerformance()
	}

	return metrics
}

// ResetStrategy resets a specific strategy
func (sm *StrategyManager) ResetStrategy(strategyID string) error {
	strategy, err := sm.GetStrategy(strategyID)
	if err != nil {
		return err
	}

	strategy.Reset()
	return nil
}

// ResetAllStrategies resets all strategies
func (sm *StrategyManager) ResetAllStrategies() {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	for _, strategy := range sm.strategies {
		strategy.Reset()
	}
}

// CreateDefaultStrategies creates and registers all 7 default trading strategies
func (sm *StrategyManager) CreateDefaultStrategies() error {
	// 1. DCA Strategy
	dcaConfig := &DCAConfig{
		InvestmentAmount:   decimal.NewFromFloat(100),
		Interval:           time.Hour,
		MaxDeviation:       decimal.NewFromFloat(0.05),
		AccumulationPeriod: time.Hour * 24,
		TradingPairs:       []string{"BTC/USDT", "ETH/USDT"},
		Exchange:           "binance",
	}
	dcaStrategy := NewDCAStrategy(sm.logger, dcaConfig)
	if err := sm.RegisterStrategy("dca_bot", dcaStrategy); err != nil {
		return fmt.Errorf("failed to register DCA strategy: %w", err)
	}

	// 2. Grid Strategy
	gridConfig := &GridConfig{
		GridLevels:   20,
		GridSpacing:  decimal.NewFromFloat(0.02),
		UpperBound:   decimal.NewFromFloat(1.20),
		LowerBound:   decimal.NewFromFloat(0.80),
		OrderAmount:  decimal.NewFromFloat(50),
		TradingPairs: []string{"BNB/USDT", "ADA/USDT"},
		Exchange:     "binance",
	}
	gridStrategy := NewGridStrategy(sm.logger, gridConfig)
	if err := sm.RegisterStrategy("grid_bot", gridStrategy); err != nil {
		return fmt.Errorf("failed to register Grid strategy: %w", err)
	}

	// 3. Momentum Strategy
	momentumConfig := &MomentumConfig{
		MomentumPeriod:     14,
		RSIThresholdBuy:    decimal.NewFromFloat(30),
		RSIThresholdSell:   decimal.NewFromFloat(70),
		VolumeThreshold:    decimal.NewFromFloat(1.5),
		BreakoutThreshold:  decimal.NewFromFloat(0.03),
		TradingPairs:       []string{"SOL/USDT", "AVAX/USDT"},
		Exchange:           "coinbase",
		PositionSize:       decimal.NewFromFloat(100),
		StopLoss:           decimal.NewFromFloat(0.08),
		TakeProfit:         decimal.NewFromFloat(0.15),
	}
	momentumStrategy := NewMomentumStrategy(sm.logger, momentumConfig)
	if err := sm.RegisterStrategy("momentum_bot", momentumStrategy); err != nil {
		return fmt.Errorf("failed to register Momentum strategy: %w", err)
	}

	// 4. Mean Reversion Strategy
	meanReversionConfig := &MeanReversionConfig{
		LookbackPeriod:    20,
		StdDevMultiplier:  decimal.NewFromFloat(2.0),
		RSIOversold:       decimal.NewFromFloat(25),
		RSIOverbought:     decimal.NewFromFloat(75),
		BollingerPeriod:   20,
		TradingPairs:      []string{"DOT/USDT", "LINK/USDT"},
		Exchange:          "kraken",
		PositionSize:      decimal.NewFromFloat(100),
		StopLoss:          decimal.NewFromFloat(0.10),
		TakeProfit:        decimal.NewFromFloat(0.20),
	}
	meanReversionStrategy := NewMeanReversionStrategy(sm.logger, meanReversionConfig)
	if err := sm.RegisterStrategy("mean_reversion_bot", meanReversionStrategy); err != nil {
		return fmt.Errorf("failed to register Mean Reversion strategy: %w", err)
	}

	// 5. Arbitrage Strategy
	arbitrageConfig := &ArbitrageConfig{
		MinProfitThreshold: decimal.NewFromFloat(0.005),
		MaxExecutionTime:   30 * time.Second,
		SlippageTolerance:  decimal.NewFromFloat(0.002),
		BalanceThreshold:   decimal.NewFromFloat(0.1),
		TradingPairs:       []string{"BTC/USDT"},
		Exchanges:          []string{"binance", "coinbase", "kraken"},
		PositionSize:       decimal.NewFromFloat(1000),
	}
	arbitrageStrategy := NewArbitrageStrategy(sm.logger, arbitrageConfig)
	if err := sm.RegisterStrategy("arbitrage_bot", arbitrageStrategy); err != nil {
		return fmt.Errorf("failed to register Arbitrage strategy: %w", err)
	}

	// 6. Scalping Strategy
	scalpingConfig := &ScalpingConfig{
		Timeframe:            "1m",
		ProfitTarget:         decimal.NewFromFloat(0.002),
		MaxHoldingTime:       5 * time.Minute,
		VolumeSpikeThreshold: decimal.NewFromFloat(2.0),
		SpreadThreshold:      decimal.NewFromFloat(0.001),
		TradingPairs:         []string{"ETH/USDT"},
		Exchange:             "binance",
		PositionSize:         decimal.NewFromFloat(200),
		StopLoss:             decimal.NewFromFloat(0.003),
	}
	scalpingStrategy := NewScalpingStrategy(sm.logger, scalpingConfig)
	if err := sm.RegisterStrategy("scalping_bot", scalpingStrategy); err != nil {
		return fmt.Errorf("failed to register Scalping strategy: %w", err)
	}

	// 7. Swing Trading Strategy
	swingConfig := &SwingConfig{
		Timeframe:    "4h",
		MAFast:       12,
		MASlow:       26,
		SignalLine:   9,
		RSIPeriod:    14,
		TradingPairs: []string{"MATIC/USDT", "UNI/USDT", "AAVE/USDT", "SUSHI/USDT"},
		Exchange:     "binance",
		PositionSize: decimal.NewFromFloat(100),
		StopLoss:     decimal.NewFromFloat(0.12),
		TakeProfit:   decimal.NewFromFloat(0.25),
	}
	swingStrategy := NewSwingStrategy(sm.logger, swingConfig)
	if err := sm.RegisterStrategy("swing_bot", swingStrategy); err != nil {
		return fmt.Errorf("failed to register Swing strategy: %w", err)
	}

	sm.logger.Info(context.Background(), "All 7 default trading strategies created", map[string]interface{}{
		"total_strategies": len(sm.strategies),
	})

	return nil
}

// Strategy configuration structs for the remaining strategies
type MeanReversionConfig struct {
	LookbackPeriod   int             `yaml:"lookback_period"`
	StdDevMultiplier decimal.Decimal `yaml:"std_dev_multiplier"`
	RSIOversold      decimal.Decimal `yaml:"rsi_oversold"`
	RSIOverbought    decimal.Decimal `yaml:"rsi_overbought"`
	BollingerPeriod  int             `yaml:"bollinger_period"`
	TradingPairs     []string        `yaml:"trading_pairs"`
	Exchange         string          `yaml:"exchange"`
	PositionSize     decimal.Decimal `yaml:"position_size"`
	StopLoss         decimal.Decimal `yaml:"stop_loss"`
	TakeProfit       decimal.Decimal `yaml:"take_profit"`
}

type ArbitrageConfig struct {
	MinProfitThreshold decimal.Decimal `yaml:"min_profit_threshold"`
	MaxExecutionTime   time.Duration   `yaml:"max_execution_time"`
	SlippageTolerance  decimal.Decimal `yaml:"slippage_tolerance"`
	BalanceThreshold   decimal.Decimal `yaml:"balance_threshold"`
	TradingPairs       []string        `yaml:"trading_pairs"`
	Exchanges          []string        `yaml:"exchanges"`
	PositionSize       decimal.Decimal `yaml:"position_size"`
}

type ScalpingConfig struct {
	Timeframe            string          `yaml:"timeframe"`
	ProfitTarget         decimal.Decimal `yaml:"profit_target"`
	MaxHoldingTime       time.Duration   `yaml:"max_holding_time"`
	VolumeSpikeThreshold decimal.Decimal `yaml:"volume_spike_threshold"`
	SpreadThreshold      decimal.Decimal `yaml:"spread_threshold"`
	TradingPairs         []string        `yaml:"trading_pairs"`
	Exchange             string          `yaml:"exchange"`
	PositionSize         decimal.Decimal `yaml:"position_size"`
	StopLoss             decimal.Decimal `yaml:"stop_loss"`
}

type SwingConfig struct {
	Timeframe    string          `yaml:"timeframe"`
	MAFast       int             `yaml:"ma_fast"`
	MASlow       int             `yaml:"ma_slow"`
	SignalLine   int             `yaml:"signal_line"`
	RSIPeriod    int             `yaml:"rsi_period"`
	TradingPairs []string        `yaml:"trading_pairs"`
	Exchange     string          `yaml:"exchange"`
	PositionSize decimal.Decimal `yaml:"position_size"`
	StopLoss     decimal.Decimal `yaml:"stop_loss"`
	TakeProfit   decimal.Decimal `yaml:"take_profit"`
}

// Placeholder strategy constructors (to be implemented)
func NewMeanReversionStrategy(logger *observability.Logger, config *MeanReversionConfig) TradingStrategy {
	// Implementation will be added
	return &DCAStrategy{logger: logger} // Temporary placeholder
}

func NewArbitrageStrategy(logger *observability.Logger, config *ArbitrageConfig) TradingStrategy {
	// Implementation will be added
	return &DCAStrategy{logger: logger} // Temporary placeholder
}

func NewScalpingStrategy(logger *observability.Logger, config *ScalpingConfig) TradingStrategy {
	// Implementation will be added
	return &DCAStrategy{logger: logger} // Temporary placeholder
}

func NewSwingStrategy(logger *observability.Logger, config *SwingConfig) TradingStrategy {
	// Implementation will be added
	return &DCAStrategy{logger: logger} // Temporary placeholder
}
