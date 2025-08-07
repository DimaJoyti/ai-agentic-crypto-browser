package trading

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// AlgorithmManager manages trading algorithms and their execution
type AlgorithmManager struct {
	logger     *observability.Logger
	algorithms map[string]*TradingAlgorithm
	strategies map[string]*TradingStrategy
	mu         sync.RWMutex
	isRunning  bool
	stopChan   chan struct{}
}

// TradingStrategy represents a complete trading strategy
type TradingStrategy struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	Algorithm       *TradingAlgorithm      `json:"algorithm"`
	Parameters      map[string]interface{} `json:"parameters"`
	RiskProfile     RiskProfile            `json:"risk_profile"`
	Performance     *StrategyPerformance   `json:"performance"`
	IsActive        bool                   `json:"is_active"`
	CreatedAt       time.Time              `json:"created_at"`
	LastExecuted    time.Time              `json:"last_executed"`
	ExecutionCount  int64                  `json:"execution_count"`
	TotalPnL        decimal.Decimal        `json:"total_pnl"`
	WinRate         float64                `json:"win_rate"`
	SharpeRatio     float64                `json:"sharpe_ratio"`
	MaxDrawdown     decimal.Decimal        `json:"max_drawdown"`
	AverageReturn   decimal.Decimal        `json:"average_return"`
	VolatilityScore float64                `json:"volatility_score"`
}

// StrategyPerformance tracks strategy performance metrics
type StrategyPerformance struct {
	TotalTrades     int64           `json:"total_trades"`
	WinningTrades   int64           `json:"winning_trades"`
	LosingTrades    int64           `json:"losing_trades"`
	TotalPnL        decimal.Decimal `json:"total_pnl"`
	TotalReturn     decimal.Decimal `json:"total_return"`
	WinRate         float64         `json:"win_rate"`
	ProfitFactor    float64         `json:"profit_factor"`
	SharpeRatio     float64         `json:"sharpe_ratio"`
	SortinoRatio    float64         `json:"sortino_ratio"`
	MaxDrawdown     decimal.Decimal `json:"max_drawdown"`
	MaxDrawdownPct  float64         `json:"max_drawdown_pct"`
	AverageWin      decimal.Decimal `json:"average_win"`
	AverageLoss     decimal.Decimal `json:"average_loss"`
	LargestWin      decimal.Decimal `json:"largest_win"`
	LargestLoss     decimal.Decimal `json:"largest_loss"`
	ConsecutiveWins int             `json:"consecutive_wins"`
	ConsecutiveLoss int             `json:"consecutive_loss"`
	CalmarRatio     float64         `json:"calmar_ratio"`
	VaR95           decimal.Decimal `json:"var_95"`
	CVaR95          decimal.Decimal `json:"cvar_95"`
	LastUpdated     time.Time       `json:"last_updated"`
}

// AlgorithmPerformance tracks algorithm-specific performance
type AlgorithmPerformance struct {
	ExecutionCount   int64           `json:"execution_count"`
	SuccessfulTrades int64           `json:"successful_trades"`
	FailedTrades     int64           `json:"failed_trades"`
	AverageLatency   time.Duration   `json:"average_latency"`
	TotalPnL         decimal.Decimal `json:"total_pnl"`
	SuccessRate      float64         `json:"success_rate"`
	AverageSlippage  decimal.Decimal `json:"average_slippage"`
	FillRate         float64         `json:"fill_rate"`
	RejectionRate    float64         `json:"rejection_rate"`
	LastExecuted     time.Time       `json:"last_executed"`
	PerformanceScore float64         `json:"performance_score"`
}

// RiskProfile defines risk parameters for algorithms and strategies
type RiskProfile struct {
	MaxPositionSize    decimal.Decimal `json:"max_position_size"`
	MaxDailyLoss       decimal.Decimal `json:"max_daily_loss"`
	MaxDrawdown        decimal.Decimal `json:"max_drawdown"`
	RiskPerTrade       decimal.Decimal `json:"risk_per_trade"`
	StopLossPercent    decimal.Decimal `json:"stop_loss_percent"`
	TakeProfitPercent  decimal.Decimal `json:"take_profit_percent"`
	MaxLeverage        decimal.Decimal `json:"max_leverage"`
	VolatilityLimit    decimal.Decimal `json:"volatility_limit"`
	CorrelationLimit   decimal.Decimal `json:"correlation_limit"`
	ConcentrationLimit decimal.Decimal `json:"concentration_limit"`
}

// RiskLevel defines risk tolerance levels
type RiskLevel string

const (
	RiskLevelConservative RiskLevel = "conservative"
	RiskLevelModerate     RiskLevel = "moderate"
	RiskLevelAggressive   RiskLevel = "aggressive"
	RiskLevelSpeculative  RiskLevel = "speculative"
)

// NewAlgorithmManager creates a new algorithm manager
func NewAlgorithmManager(logger *observability.Logger) *AlgorithmManager {
	return &AlgorithmManager{
		logger:     logger,
		algorithms: make(map[string]*TradingAlgorithm),
		strategies: make(map[string]*TradingStrategy),
		stopChan:   make(chan struct{}),
	}
}

// Start starts the algorithm manager
func (am *AlgorithmManager) Start(ctx context.Context) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	if am.isRunning {
		return fmt.Errorf("algorithm manager is already running")
	}

	am.isRunning = true

	// Initialize default algorithms
	am.initializeDefaultAlgorithms()

	// Start background processes
	go am.performanceMonitoringLoop(ctx)

	// Logger info call removed to prevent hanging

	return nil
}

// Stop stops the algorithm manager
func (am *AlgorithmManager) Stop(ctx context.Context) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	if !am.isRunning {
		return nil
	}

	am.isRunning = false
	close(am.stopChan)

	am.logger.Info(ctx, "Algorithm manager stopped", nil)
	return nil
}

// RegisterAlgorithm registers a new trading algorithm
func (am *AlgorithmManager) RegisterAlgorithm(algorithm *TradingAlgorithm) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	if algorithm.ID == "" {
		algorithm.ID = uuid.New().String()
	}

	algorithm.CreatedAt = time.Now()
	algorithm.Performance = &AlgorithmPerformance{
		LastExecuted: time.Now(),
	}

	am.algorithms[algorithm.ID] = algorithm

	// Logger info call removed to prevent hanging

	return nil
}

// CreateStrategy creates a new trading strategy
func (am *AlgorithmManager) CreateStrategy(name, description string, algorithmType AlgorithmType, parameters map[string]interface{}, riskProfile RiskProfile) (*TradingStrategy, error) {
	am.mu.Lock()
	defer am.mu.Unlock()

	// Find algorithm
	var algorithm *TradingAlgorithm
	for _, alg := range am.algorithms {
		if alg.Type == algorithmType && alg.IsActive {
			algorithm = alg
			break
		}
	}

	if algorithm == nil {
		return nil, fmt.Errorf("no active algorithm found for type: %s", algorithmType)
	}

	strategy := &TradingStrategy{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Algorithm:   algorithm,
		Parameters:  parameters,
		RiskProfile: riskProfile,
		Performance: &StrategyPerformance{
			LastUpdated: time.Now(),
		},
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	am.strategies[strategy.ID] = strategy

	// Logger info call removed to prevent hanging

	return strategy, nil
}

// GetAlgorithm retrieves an algorithm by ID
func (am *AlgorithmManager) GetAlgorithm(algorithmID string) (*TradingAlgorithm, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	algorithm, exists := am.algorithms[algorithmID]
	if !exists {
		return nil, fmt.Errorf("algorithm not found: %s", algorithmID)
	}

	return algorithm, nil
}

// GetStrategy retrieves a strategy by ID
func (am *AlgorithmManager) GetStrategy(strategyID string) (*TradingStrategy, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	strategy, exists := am.strategies[strategyID]
	if !exists {
		return nil, fmt.Errorf("strategy not found: %s", strategyID)
	}

	return strategy, nil
}

// GetActiveStrategies returns all active strategies
func (am *AlgorithmManager) GetActiveStrategies() []*TradingStrategy {
	am.mu.RLock()
	defer am.mu.RUnlock()

	var activeStrategies []*TradingStrategy
	for _, strategy := range am.strategies {
		if strategy.IsActive {
			activeStrategies = append(activeStrategies, strategy)
		}
	}

	return activeStrategies
}

// UpdateStrategyPerformance updates strategy performance metrics
func (am *AlgorithmManager) UpdateStrategyPerformance(strategyID string, trade *TradeResult) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	strategy, exists := am.strategies[strategyID]
	if !exists {
		return fmt.Errorf("strategy not found: %s", strategyID)
	}

	perf := strategy.Performance
	perf.TotalTrades++
	perf.TotalPnL = perf.TotalPnL.Add(trade.PnL)

	if trade.PnL.GreaterThan(decimal.Zero) {
		perf.WinningTrades++
	} else {
		perf.LosingTrades++
	}

	// Update win rate
	perf.WinRate = float64(perf.WinningTrades) / float64(perf.TotalTrades)

	// Update other metrics
	am.calculateAdvancedMetrics(perf)

	perf.LastUpdated = time.Now()
	strategy.LastExecuted = time.Now()
	strategy.ExecutionCount++

	return nil
}

// initializeDefaultAlgorithms creates default trading algorithms
func (am *AlgorithmManager) initializeDefaultAlgorithms() {
	defaultAlgorithms := []*TradingAlgorithm{
		{
			Type:        AlgorithmTypeTWAP,
			Name:        "Time-Weighted Average Price",
			Description: "Executes orders over time to achieve average price",
			Parameters: map[string]interface{}{
				"duration_minutes": 60,
				"slice_count":      10,
				"randomization":    0.1,
			},
			RiskProfile: RiskProfile{
				MaxPositionSize: decimal.NewFromFloat(0.1),
				RiskPerTrade:    decimal.NewFromFloat(0.02),
			},
			IsActive: true,
		},
		{
			Type:        AlgorithmTypeVWAP,
			Name:        "Volume-Weighted Average Price",
			Description: "Executes orders based on historical volume patterns",
			Parameters: map[string]interface{}{
				"lookback_days":    30,
				"volume_threshold": 0.05,
				"participation":    0.1,
			},
			RiskProfile: RiskProfile{
				MaxPositionSize: decimal.NewFromFloat(0.15),
				RiskPerTrade:    decimal.NewFromFloat(0.025),
			},
			IsActive: true,
		},
		{
			Type:        AlgorithmTypeIceberg,
			Name:        "Iceberg Order",
			Description: "Hides large orders by showing only small portions",
			Parameters: map[string]interface{}{
				"visible_size":    0.05,
				"randomization":   0.2,
				"refresh_time_ms": 1000,
			},
			RiskProfile: RiskProfile{
				MaxPositionSize: decimal.NewFromFloat(0.2),
				RiskPerTrade:    decimal.NewFromFloat(0.03),
			},
			IsActive: true,
		},
	}

	for _, algorithm := range defaultAlgorithms {
		if algorithm.ID == "" {
			algorithm.ID = uuid.New().String()
		}
		algorithm.CreatedAt = time.Now()
		algorithm.Performance = &AlgorithmPerformance{
			LastExecuted: time.Now(),
		}
		am.algorithms[algorithm.ID] = algorithm
	}
}

// performanceMonitoringLoop monitors algorithm and strategy performance
func (am *AlgorithmManager) performanceMonitoringLoop(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-am.stopChan:
			return
		case <-ticker.C:
			am.updatePerformanceMetrics()
		}
	}
}

// updatePerformanceMetrics updates performance metrics for all strategies
func (am *AlgorithmManager) updatePerformanceMetrics() {
	am.mu.RLock()
	defer am.mu.RUnlock()

	for _, strategy := range am.strategies {
		if strategy.IsActive {
			am.calculateAdvancedMetrics(strategy.Performance)
		}
	}
}

// calculateAdvancedMetrics calculates advanced performance metrics
func (am *AlgorithmManager) calculateAdvancedMetrics(perf *StrategyPerformance) {
	if perf.TotalTrades == 0 {
		return
	}

	// Calculate profit factor
	if perf.LosingTrades > 0 {
		avgWin := perf.TotalPnL.Div(decimal.NewFromInt(perf.WinningTrades))
		avgLoss := perf.TotalPnL.Div(decimal.NewFromInt(perf.LosingTrades)).Abs()
		if avgLoss.GreaterThan(decimal.Zero) {
			profitFactor, _ := avgWin.Div(avgLoss).Float64()
			perf.ProfitFactor = profitFactor
		}
	}

	// Calculate Sharpe ratio (simplified)
	if perf.TotalTrades > 10 {
		avgReturn, _ := perf.TotalPnL.Div(decimal.NewFromInt(perf.TotalTrades)).Float64()
		perf.SharpeRatio = avgReturn * 0.1 // Simplified calculation
	}
}

// TradeResult represents the result of a trade execution
type TradeResult struct {
	TradeID    string          `json:"trade_id"`
	StrategyID string          `json:"strategy_id"`
	Symbol     string          `json:"symbol"`
	Side       string          `json:"side"`
	Quantity   decimal.Decimal `json:"quantity"`
	Price      decimal.Decimal `json:"price"`
	PnL        decimal.Decimal `json:"pnl"`
	Commission decimal.Decimal `json:"commission"`
	Slippage   decimal.Decimal `json:"slippage"`
	ExecutedAt time.Time       `json:"executed_at"`
	Success    bool            `json:"success"`
	ErrorMsg   string          `json:"error_msg,omitempty"`
}
