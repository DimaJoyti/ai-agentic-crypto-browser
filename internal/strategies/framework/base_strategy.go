package framework

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// BaseStrategy provides common functionality for all strategies
type BaseStrategy struct {
	id          uuid.UUID
	name        string
	description string
	version     string
	config      StrategyConfig
	parameters  map[string]Parameter
	riskLimits  *RiskLimits
	metrics     *StrategyMetrics
	positions   map[string]*Position
	orders      map[uuid.UUID]*Order

	// State management
	isRunning bool
	startTime time.Time
	stopChan  chan struct{}
	wg        sync.WaitGroup
	mu        sync.RWMutex
}

// NewBaseStrategy creates a new base strategy
func NewBaseStrategy(name, description, version string) *BaseStrategy {
	return &BaseStrategy{
		id:          uuid.New(),
		name:        name,
		description: description,
		version:     version,
		parameters:  make(map[string]Parameter),
		positions:   make(map[string]*Position),
		orders:      make(map[uuid.UUID]*Order),
		metrics: &StrategyMetrics{
			TotalTrades:     0,
			WinningTrades:   0,
			LosingTrades:    0,
			WinRate:         decimal.NewFromInt(0),
			TotalPnL:        decimal.NewFromInt(0),
			RealizedPnL:     decimal.NewFromInt(0),
			UnrealizedPnL:   decimal.NewFromInt(0),
			MaxDrawdown:     decimal.NewFromInt(0),
			SharpeRatio:     decimal.NewFromInt(0),
			SortinoRatio:    decimal.NewFromInt(0),
			ProfitFactor:    decimal.NewFromInt(0),
			AvgWin:          decimal.NewFromInt(0),
			AvgLoss:         decimal.NewFromInt(0),
			MaxWin:          decimal.NewFromInt(0),
			MaxLoss:         decimal.NewFromInt(0),
			TotalVolume:     decimal.NewFromInt(0),
			TotalCommission: decimal.NewFromInt(0),
		},
		stopChan: make(chan struct{}),
	}
}

// Lifecycle management

// Initialize initializes the strategy with configuration
func (bs *BaseStrategy) Initialize(ctx context.Context, config StrategyConfig) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	bs.config = config
	bs.id = config.ID

	// Set risk limits if provided
	if config.RiskLimits != nil {
		bs.riskLimits = config.RiskLimits
	}

	// Set parameters
	for name, value := range config.Parameters {
		if param, exists := bs.parameters[name]; exists {
			param.Value = value
			bs.parameters[name] = param
		}
	}

	return bs.ValidateParameters()
}

// Start starts the strategy
func (bs *BaseStrategy) Start(ctx context.Context) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if bs.isRunning {
		return fmt.Errorf("strategy is already running")
	}

	bs.isRunning = true
	bs.startTime = time.Now()
	bs.metrics.StartTime = bs.startTime

	return nil
}

// Stop stops the strategy
func (bs *BaseStrategy) Stop(ctx context.Context) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if !bs.isRunning {
		return fmt.Errorf("strategy is not running")
	}

	close(bs.stopChan)
	bs.wg.Wait()

	bs.isRunning = false

	return nil
}

// IsRunning returns whether the strategy is running
func (bs *BaseStrategy) IsRunning() bool {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	return bs.isRunning
}

// Strategy identification

// GetID returns the strategy ID
func (bs *BaseStrategy) GetID() uuid.UUID {
	return bs.id
}

// GetName returns the strategy name
func (bs *BaseStrategy) GetName() string {
	return bs.name
}

// GetDescription returns the strategy description
func (bs *BaseStrategy) GetDescription() string {
	return bs.description
}

// GetVersion returns the strategy version
func (bs *BaseStrategy) GetVersion() string {
	return bs.version
}

// Configuration and parameters

// GetParameters returns strategy parameters
func (bs *BaseStrategy) GetParameters() map[string]Parameter {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	// Return a copy to avoid race conditions
	params := make(map[string]Parameter)
	for name, param := range bs.parameters {
		params[name] = param
	}
	return params
}

// SetParameter sets a strategy parameter
func (bs *BaseStrategy) SetParameter(name string, value interface{}) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	param, exists := bs.parameters[name]
	if !exists {
		return fmt.Errorf("parameter not found: %s", name)
	}

	// Validate parameter value
	if err := bs.validateParameterValue(param, value); err != nil {
		return fmt.Errorf("invalid parameter value: %w", err)
	}

	param.Value = value
	bs.parameters[name] = param

	return nil
}

// ValidateParameters validates all strategy parameters
func (bs *BaseStrategy) ValidateParameters() error {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	for name, param := range bs.parameters {
		if param.Required && param.Value == nil {
			return fmt.Errorf("required parameter missing: %s", name)
		}

		if param.Value != nil {
			if err := bs.validateParameterValue(param, param.Value); err != nil {
				return fmt.Errorf("invalid parameter %s: %w", name, err)
			}
		}
	}

	return nil
}

// Performance and metrics

// GetMetrics returns strategy performance metrics
func (bs *BaseStrategy) GetMetrics() *StrategyMetrics {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	// Update uptime
	if bs.isRunning {
		bs.metrics.Uptime = time.Since(bs.startTime)
	}

	// Update active positions and orders
	bs.metrics.ActivePositions = len(bs.positions)
	bs.metrics.OpenOrders = len(bs.orders)

	// Return a copy
	metrics := *bs.metrics
	return &metrics
}

// GetPositions returns current positions
func (bs *BaseStrategy) GetPositions() []*Position {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	positions := make([]*Position, 0, len(bs.positions))
	for _, position := range bs.positions {
		positions = append(positions, position)
	}
	return positions
}

// GetOrders returns current orders
func (bs *BaseStrategy) GetOrders() []*Order {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	orders := make([]*Order, 0, len(bs.orders))
	for _, order := range bs.orders {
		orders = append(orders, order)
	}
	return orders
}

// Risk management

// GetRiskLimits returns risk limits
func (bs *BaseStrategy) GetRiskLimits() *RiskLimits {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	if bs.riskLimits == nil {
		return nil
	}

	// Return a copy
	limits := *bs.riskLimits
	return &limits
}

// SetRiskLimits sets risk limits
func (bs *BaseStrategy) SetRiskLimits(limits *RiskLimits) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	bs.riskLimits = limits
	return nil
}

// Protected methods for derived strategies

// AddParameter adds a parameter definition
func (bs *BaseStrategy) AddParameter(param Parameter) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bs.parameters[param.Name] = param
}

// GetParameterValue gets a parameter value with type assertion
func (bs *BaseStrategy) GetParameterValue(name string) (interface{}, error) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	param, exists := bs.parameters[name]
	if !exists {
		return nil, fmt.Errorf("parameter not found: %s", name)
	}

	return param.Value, nil
}

// UpdateMetrics updates strategy metrics
func (bs *BaseStrategy) UpdateMetrics(trade *Trade) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	bs.metrics.TotalTrades++
	bs.metrics.TotalVolume = bs.metrics.TotalVolume.Add(trade.Quantity)
	bs.metrics.TotalCommission = bs.metrics.TotalCommission.Add(trade.Commission)
	bs.metrics.LastTradeTime = trade.Timestamp

	if trade.PnL.GreaterThan(decimal.NewFromInt(0)) {
		bs.metrics.WinningTrades++
		if trade.PnL.GreaterThan(bs.metrics.MaxWin) {
			bs.metrics.MaxWin = trade.PnL
		}
	} else if trade.PnL.LessThan(decimal.NewFromInt(0)) {
		bs.metrics.LosingTrades++
		if trade.PnL.LessThan(bs.metrics.MaxLoss) {
			bs.metrics.MaxLoss = trade.PnL
		}
	}

	// Update win rate
	if bs.metrics.TotalTrades > 0 {
		bs.metrics.WinRate = decimal.NewFromInt(bs.metrics.WinningTrades).Div(decimal.NewFromInt(bs.metrics.TotalTrades))
	}

	// Update total PnL
	bs.metrics.TotalPnL = bs.metrics.TotalPnL.Add(trade.PnL)
	bs.metrics.RealizedPnL = bs.metrics.RealizedPnL.Add(trade.PnL)
}

// Trade represents a completed trade
type Trade struct {
	Symbol     string          `json:"symbol"`
	Side       string          `json:"side"`
	Quantity   decimal.Decimal `json:"quantity"`
	Price      decimal.Decimal `json:"price"`
	PnL        decimal.Decimal `json:"pnl"`
	Commission decimal.Decimal `json:"commission"`
	Timestamp  time.Time       `json:"timestamp"`
}

// Private methods

// validateParameterValue validates a parameter value
func (bs *BaseStrategy) validateParameterValue(param Parameter, value interface{}) error {
	switch param.Type {
	case ParamTypeString:
		if _, ok := value.(string); !ok {
			return fmt.Errorf("expected string, got %T", value)
		}
	case ParamTypeInt:
		if _, ok := value.(int); !ok {
			return fmt.Errorf("expected int, got %T", value)
		}
	case ParamTypeFloat:
		if _, ok := value.(float64); !ok {
			return fmt.Errorf("expected float64, got %T", value)
		}
	case ParamTypeBool:
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("expected bool, got %T", value)
		}
	case ParamTypeDecimal:
		if _, ok := value.(decimal.Decimal); !ok {
			return fmt.Errorf("expected decimal.Decimal, got %T", value)
		}
	}

	// Check min/max values if specified
	if param.MinValue != nil || param.MaxValue != nil {
		// TODO: Implement min/max validation for different types
	}

	return nil
}
