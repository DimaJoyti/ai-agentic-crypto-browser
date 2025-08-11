package strategies

import (
	"context"
	"fmt"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/shopspring/decimal"
)

// DCAStrategy implements Dollar Cost Averaging trading strategy
type DCAStrategy struct {
	logger         *observability.Logger
	config         *DCAConfig
	lastExecution  time.Time
	totalInvested  decimal.Decimal
	totalTokens    decimal.Decimal
	averagePrice   decimal.Decimal
	executionCount int
}

// DCAConfig holds configuration for DCA strategy
type DCAConfig struct {
	InvestmentAmount   decimal.Decimal `yaml:"investment_amount"`
	Interval           time.Duration   `yaml:"interval"`
	MaxDeviation       decimal.Decimal `yaml:"max_deviation"`
	AccumulationPeriod time.Duration   `yaml:"accumulation_period"`
	TradingPairs       []string        `yaml:"trading_pairs"`
	Exchange           string          `yaml:"exchange"`
}

// DCAOrder represents a DCA buy order
type DCAOrder struct {
	ID             string          `json:"id"`
	Symbol         string          `json:"symbol"`
	Amount         decimal.Decimal `json:"amount"`
	Price          decimal.Decimal `json:"price"`
	Timestamp      time.Time       `json:"timestamp"`
	Status         OrderStatus     `json:"status"`
	ExecutedAmount decimal.Decimal `json:"executed_amount"`
	ExecutedPrice  decimal.Decimal `json:"executed_price"`
}

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusExecuted  OrderStatus = "executed"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusFailed    OrderStatus = "failed"
)

// NewDCAStrategy creates a new DCA strategy instance
func NewDCAStrategy(logger *observability.Logger, config *DCAConfig) *DCAStrategy {
	return &DCAStrategy{
		logger:        logger,
		config:        config,
		totalInvested: decimal.Zero,
		totalTokens:   decimal.Zero,
		averagePrice:  decimal.Zero,
	}
}

// Execute executes the DCA strategy
func (dca *DCAStrategy) Execute(ctx context.Context, marketData *MarketData) (interface{}, error) {
	// Check if it's time to execute
	if !dca.shouldExecute() {
		return nil, nil
	}

	// Get current price
	currentPrice := marketData.Price
	if currentPrice.IsZero() {
		return nil, fmt.Errorf("invalid market price")
	}

	// Check price deviation if we have previous executions
	if dca.executionCount > 0 && !dca.isPriceWithinDeviation(currentPrice) {
		dca.logger.Info(ctx, "Skipping DCA execution due to price deviation", map[string]interface{}{
			"current_price": currentPrice.String(),
			"average_price": dca.averagePrice.String(),
			"max_deviation": dca.config.MaxDeviation.String(),
		})
		return nil, nil
	}

	// Calculate order amount
	orderAmount := dca.config.InvestmentAmount
	tokenAmount := orderAmount.Div(currentPrice)

	// Create trading signal
	signal := &TradingSignal{
		ID:        generateOrderID(),
		Symbol:    marketData.Symbol,
		Action:    ActionBuy,
		Amount:    tokenAmount,
		Price:     currentPrice,
		OrderType: OrderTypeMarket,
		Timestamp: time.Now(),
		Strategy:  "DCA",
		Metadata: map[string]interface{}{
			"investment_amount": orderAmount.String(),
			"execution_count":   dca.executionCount + 1,
			"average_price":     dca.averagePrice.String(),
		},
	}

	// Update internal state
	dca.updateState(orderAmount, tokenAmount, currentPrice)

	dca.logger.Info(ctx, "DCA order created", map[string]interface{}{
		"symbol":            signal.Symbol,
		"amount":            signal.Amount.String(),
		"price":             signal.Price.String(),
		"investment_amount": orderAmount.String(),
		"execution_count":   dca.executionCount,
	})

	return signal, nil
}

// shouldExecute checks if it's time to execute the DCA strategy
func (dca *DCAStrategy) shouldExecute() bool {
	if dca.lastExecution.IsZero() {
		return true
	}

	return time.Since(dca.lastExecution) >= dca.config.Interval
}

// isPriceWithinDeviation checks if current price is within acceptable deviation
func (dca *DCAStrategy) isPriceWithinDeviation(currentPrice decimal.Decimal) bool {
	if dca.averagePrice.IsZero() {
		return true
	}

	deviation := currentPrice.Sub(dca.averagePrice).Div(dca.averagePrice).Abs()
	return deviation.LessThanOrEqual(dca.config.MaxDeviation)
}

// updateState updates the internal state after execution
func (dca *DCAStrategy) updateState(investmentAmount, tokenAmount, price decimal.Decimal) {
	dca.totalInvested = dca.totalInvested.Add(investmentAmount)
	dca.totalTokens = dca.totalTokens.Add(tokenAmount)
	dca.averagePrice = dca.totalInvested.Div(dca.totalTokens)
	dca.executionCount++
	dca.lastExecution = time.Now()
}

// GetPerformance returns performance metrics for the DCA strategy
func (dca *DCAStrategy) GetPerformance() *StrategyPerformance {
	if dca.totalTokens.IsZero() {
		return &StrategyPerformance{
			TotalInvested:  decimal.Zero,
			CurrentValue:   decimal.Zero,
			UnrealizedPnL:  decimal.Zero,
			ROI:            decimal.Zero,
			ExecutionCount: 0,
		}
	}

	// Use average price as current price estimate if no current price available
	currentPrice := dca.averagePrice
	if currentPrice.IsZero() {
		currentPrice = decimal.NewFromFloat(1.0) // Default fallback
	}

	currentValue := dca.totalTokens.Mul(currentPrice)
	unrealizedPnL := currentValue.Sub(dca.totalInvested)
	roi := decimal.Zero
	if !dca.totalInvested.IsZero() {
		roi = unrealizedPnL.Div(dca.totalInvested).Mul(decimal.NewFromInt(100))
	}

	return &StrategyPerformance{
		TotalInvested:  dca.totalInvested,
		CurrentValue:   currentValue,
		UnrealizedPnL:  unrealizedPnL,
		ROI:            roi,
		AveragePrice:   dca.averagePrice,
		TotalTokens:    dca.totalTokens,
		ExecutionCount: dca.executionCount,
		LastExecution:  dca.lastExecution,
	}
}

// Reset resets the DCA strategy state
func (dca *DCAStrategy) Reset() {
	dca.totalInvested = decimal.Zero
	dca.totalTokens = decimal.Zero
	dca.averagePrice = decimal.Zero
	dca.executionCount = 0
	dca.lastExecution = time.Time{}
}

// GetConfig returns the strategy configuration
func (dca *DCAStrategy) GetConfig() *DCAConfig {
	return dca.config
}

// UpdateConfig updates the strategy configuration
func (dca *DCAStrategy) UpdateConfig(config *DCAConfig) {
	dca.config = config
}

// GetName returns the strategy name
func (dca *DCAStrategy) GetName() string {
	return "Dollar Cost Averaging"
}

// GetType returns the strategy type
func (dca *DCAStrategy) GetType() string {
	return "DCA"
}

// Validate validates the strategy configuration
func (dca *DCAStrategy) Validate() error {
	if dca.config.InvestmentAmount.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("investment amount must be positive")
	}

	if dca.config.Interval <= 0 {
		return fmt.Errorf("interval must be positive")
	}

	if dca.config.MaxDeviation.LessThan(decimal.Zero) || dca.config.MaxDeviation.GreaterThan(decimal.NewFromFloat(1.0)) {
		return fmt.Errorf("max deviation must be between 0 and 1")
	}

	if len(dca.config.TradingPairs) == 0 {
		return fmt.Errorf("at least one trading pair must be specified")
	}

	return nil
}

// MarketData represents market data for a trading pair
type MarketData struct {
	Symbol    string          `json:"symbol"`
	Price     decimal.Decimal `json:"price"`
	Volume    decimal.Decimal `json:"volume"`
	Timestamp time.Time       `json:"timestamp"`
	High24h   decimal.Decimal `json:"high_24h"`
	Low24h    decimal.Decimal `json:"low_24h"`
	Change24h decimal.Decimal `json:"change_24h"`
}

// TradingSignal represents a trading signal generated by a strategy
type TradingSignal struct {
	ID        string                 `json:"id"`
	Symbol    string                 `json:"symbol"`
	Action    TradingAction          `json:"action"`
	Amount    decimal.Decimal        `json:"amount"`
	Price     decimal.Decimal        `json:"price"`
	OrderType OrderType              `json:"order_type"`
	Timestamp time.Time              `json:"timestamp"`
	Strategy  string                 `json:"strategy"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// TradingAction represents the type of trading action
type TradingAction string

const (
	ActionBuy  TradingAction = "buy"
	ActionSell TradingAction = "sell"
	ActionHold TradingAction = "hold"
)

// OrderType represents the type of order
type OrderType string

const (
	OrderTypeMarket OrderType = "market"
	OrderTypeLimit  OrderType = "limit"
	OrderTypeStop   OrderType = "stop"
)

// StrategyPerformance represents performance metrics for a strategy
type StrategyPerformance struct {
	TotalInvested  decimal.Decimal `json:"total_invested"`
	CurrentValue   decimal.Decimal `json:"current_value"`
	UnrealizedPnL  decimal.Decimal `json:"unrealized_pnl"`
	ROI            decimal.Decimal `json:"roi"`
	AveragePrice   decimal.Decimal `json:"average_price"`
	TotalTokens    decimal.Decimal `json:"total_tokens"`
	ExecutionCount int             `json:"execution_count"`
	LastExecution  time.Time       `json:"last_execution"`
}

// generateOrderID generates a unique order ID
func generateOrderID() string {
	return fmt.Sprintf("dca_%d", time.Now().UnixNano())
}
