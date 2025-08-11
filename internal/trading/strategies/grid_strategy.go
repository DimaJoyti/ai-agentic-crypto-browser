package strategies

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/shopspring/decimal"
)

// GridStrategy implements Grid Trading strategy
type GridStrategy struct {
	logger        *observability.Logger
	config        *GridConfig
	gridLevels    []*GridLevel
	activeOrders  map[string]*GridOrder
	basePrice     decimal.Decimal
	isInitialized bool
}

// GridConfig holds configuration for Grid strategy
type GridConfig struct {
	GridLevels   int             `yaml:"grid_levels"`
	GridSpacing  decimal.Decimal `yaml:"grid_spacing"`
	UpperBound   decimal.Decimal `yaml:"upper_bound"`
	LowerBound   decimal.Decimal `yaml:"lower_bound"`
	OrderAmount  decimal.Decimal `yaml:"order_amount"`
	TradingPairs []string        `yaml:"trading_pairs"`
	Exchange     string          `yaml:"exchange"`
}

// GridLevel represents a single level in the grid
type GridLevel struct {
	ID        int             `json:"id"`
	Price     decimal.Decimal `json:"price"`
	BuyOrder  *GridOrder      `json:"buy_order,omitempty"`
	SellOrder *GridOrder      `json:"sell_order,omitempty"`
	IsActive  bool            `json:"is_active"`
	Filled    bool            `json:"filled"`
}

// GridOrder represents an order in the grid
type GridOrder struct {
	ID        string          `json:"id"`
	Level     int             `json:"level"`
	Symbol    string          `json:"symbol"`
	Side      OrderSide       `json:"side"`
	Amount    decimal.Decimal `json:"amount"`
	Price     decimal.Decimal `json:"price"`
	Status    OrderStatus     `json:"status"`
	Timestamp time.Time       `json:"timestamp"`
	FilledAt  *time.Time      `json:"filled_at,omitempty"`
}

// OrderSide represents the side of an order
type OrderSide string

const (
	OrderSideBuy  OrderSide = "buy"
	OrderSideSell OrderSide = "sell"
)

// NewGridStrategy creates a new Grid strategy instance
func NewGridStrategy(logger *observability.Logger, config *GridConfig) *GridStrategy {
	return &GridStrategy{
		logger:       logger,
		config:       config,
		gridLevels:   make([]*GridLevel, 0),
		activeOrders: make(map[string]*GridOrder),
	}
}

// Initialize initializes the grid based on current market price
func (gs *GridStrategy) Initialize(ctx context.Context, currentPrice decimal.Decimal) error {
	if gs.isInitialized {
		return fmt.Errorf("grid strategy already initialized")
	}

	gs.basePrice = currentPrice

	// Calculate grid levels
	upperPrice := currentPrice.Mul(gs.config.UpperBound)
	lowerPrice := currentPrice.Mul(gs.config.LowerBound)

	priceRange := upperPrice.Sub(lowerPrice)
	levelSpacing := priceRange.Div(decimal.NewFromInt(int64(gs.config.GridLevels - 1)))

	// Create grid levels
	for i := 0; i < gs.config.GridLevels; i++ {
		levelPrice := lowerPrice.Add(levelSpacing.Mul(decimal.NewFromInt(int64(i))))

		level := &GridLevel{
			ID:       i,
			Price:    levelPrice,
			IsActive: true,
			Filled:   false,
		}

		gs.gridLevels = append(gs.gridLevels, level)
	}

	// Sort levels by price
	sort.Slice(gs.gridLevels, func(i, j int) bool {
		return gs.gridLevels[i].Price.LessThan(gs.gridLevels[j].Price)
	})

	gs.isInitialized = true

	gs.logger.Info(ctx, "Grid strategy initialized", map[string]interface{}{
		"base_price":    currentPrice.String(),
		"upper_price":   upperPrice.String(),
		"lower_price":   lowerPrice.String(),
		"grid_levels":   len(gs.gridLevels),
		"level_spacing": levelSpacing.String(),
	})

	return nil
}

// Execute executes the Grid strategy
func (gs *GridStrategy) Execute(ctx context.Context, marketData *MarketData) (interface{}, error) {
	if !gs.isInitialized {
		if err := gs.Initialize(ctx, marketData.Price); err != nil {
			return nil, fmt.Errorf("failed to initialize grid: %w", err)
		}
	}

	var signals []*TradingSignal
	currentPrice := marketData.Price

	// Check each grid level
	for _, level := range gs.gridLevels {
		if !level.IsActive {
			continue
		}

		// Check for buy opportunities (price at or below level)
		if currentPrice.LessThanOrEqual(level.Price) && level.BuyOrder == nil {
			buySignal := gs.createBuyOrder(level, marketData)
			if buySignal != nil {
				signals = append(signals, buySignal)
			}
		}

		// Check for sell opportunities (price at or above level and we have inventory)
		if currentPrice.GreaterThanOrEqual(level.Price) && level.SellOrder == nil && gs.hasInventoryAtLevel(level) {
			sellSignal := gs.createSellOrder(level, marketData)
			if sellSignal != nil {
				signals = append(signals, sellSignal)
			}
		}
	}

	// Process filled orders
	gs.processFilledOrders(ctx, currentPrice)

	return signals, nil
}

// createBuyOrder creates a buy order for a grid level
func (gs *GridStrategy) createBuyOrder(level *GridLevel, marketData *MarketData) *TradingSignal {
	orderID := fmt.Sprintf("grid_buy_%d_%d", level.ID, time.Now().UnixNano())

	order := &GridOrder{
		ID:        orderID,
		Level:     level.ID,
		Symbol:    marketData.Symbol,
		Side:      OrderSideBuy,
		Amount:    gs.config.OrderAmount.Div(level.Price), // Calculate token amount
		Price:     level.Price,
		Status:    OrderStatusPending,
		Timestamp: time.Now(),
	}

	level.BuyOrder = order
	gs.activeOrders[orderID] = order

	return &TradingSignal{
		ID:        orderID,
		Symbol:    marketData.Symbol,
		Action:    ActionBuy,
		Amount:    order.Amount,
		Price:     order.Price,
		OrderType: OrderTypeLimit,
		Timestamp: time.Now(),
		Strategy:  "Grid",
		Metadata: map[string]interface{}{
			"grid_level":   level.ID,
			"grid_price":   level.Price.String(),
			"order_amount": gs.config.OrderAmount.String(),
		},
	}
}

// createSellOrder creates a sell order for a grid level
func (gs *GridStrategy) createSellOrder(level *GridLevel, marketData *MarketData) *TradingSignal {
	orderID := fmt.Sprintf("grid_sell_%d_%d", level.ID, time.Now().UnixNano())

	order := &GridOrder{
		ID:        orderID,
		Level:     level.ID,
		Symbol:    marketData.Symbol,
		Side:      OrderSideSell,
		Amount:    gs.config.OrderAmount.Div(level.Price), // Calculate token amount
		Price:     level.Price,
		Status:    OrderStatusPending,
		Timestamp: time.Now(),
	}

	level.SellOrder = order
	gs.activeOrders[orderID] = order

	return &TradingSignal{
		ID:        orderID,
		Symbol:    marketData.Symbol,
		Action:    ActionSell,
		Amount:    order.Amount,
		Price:     order.Price,
		OrderType: OrderTypeLimit,
		Timestamp: time.Now(),
		Strategy:  "Grid",
		Metadata: map[string]interface{}{
			"grid_level":   level.ID,
			"grid_price":   level.Price.String(),
			"order_amount": gs.config.OrderAmount.String(),
		},
	}
}

// hasInventoryAtLevel checks if we have inventory to sell at this level
func (gs *GridStrategy) hasInventoryAtLevel(level *GridLevel) bool {
	// Check if there was a previous buy order that was filled
	return level.BuyOrder != nil && level.BuyOrder.Status == OrderStatusExecuted
}

// processFilledOrders processes orders that have been filled
func (gs *GridStrategy) processFilledOrders(ctx context.Context, currentPrice decimal.Decimal) {
	for orderID, order := range gs.activeOrders {
		if order.Status == OrderStatusExecuted {
			gs.handleFilledOrder(ctx, order)
			delete(gs.activeOrders, orderID)
		}
	}
}

// handleFilledOrder handles a filled order
func (gs *GridStrategy) handleFilledOrder(ctx context.Context, order *GridOrder) {
	level := gs.gridLevels[order.Level]

	if order.Side == OrderSideBuy {
		level.BuyOrder = order
		level.Filled = true

		gs.logger.Info(ctx, "Grid buy order filled", map[string]interface{}{
			"order_id": order.ID,
			"level":    order.Level,
			"price":    order.Price.String(),
			"amount":   order.Amount.String(),
		})
	} else {
		level.SellOrder = order
		level.Filled = false // Reset for next cycle
		level.BuyOrder = nil // Clear buy order

		gs.logger.Info(ctx, "Grid sell order filled", map[string]interface{}{
			"order_id": order.ID,
			"level":    order.Level,
			"price":    order.Price.String(),
			"amount":   order.Amount.String(),
		})
	}
}

// UpdateOrderStatus updates the status of an order
func (gs *GridStrategy) UpdateOrderStatus(orderID string, status OrderStatus) error {
	order, exists := gs.activeOrders[orderID]
	if !exists {
		return fmt.Errorf("order not found: %s", orderID)
	}

	order.Status = status
	if status == OrderStatusExecuted {
		now := time.Now()
		order.FilledAt = &now
	}

	return nil
}

// GetPerformance returns performance metrics for the Grid strategy
func (gs *GridStrategy) GetPerformance() *StrategyPerformance {
	totalBuyOrders := 0
	totalSellOrders := 0
	totalProfit := decimal.Zero
	activeGridLevels := 0

	for _, level := range gs.gridLevels {
		if level.IsActive {
			activeGridLevels++
		}

		if level.BuyOrder != nil && level.BuyOrder.Status == OrderStatusExecuted {
			totalBuyOrders++
		}

		if level.SellOrder != nil && level.SellOrder.Status == OrderStatusExecuted {
			totalSellOrders++
			// Calculate profit from this sell order
			if level.BuyOrder != nil {
				profit := level.SellOrder.Price.Sub(level.BuyOrder.Price).Mul(level.SellOrder.Amount)
				totalProfit = totalProfit.Add(profit)
			}
		}
	}

	// Use base price as current price estimate
	currentPrice := gs.basePrice
	if currentPrice.IsZero() {
		currentPrice = decimal.NewFromFloat(1.0) // Default fallback
	}

	// Calculate total capital from grid configuration
	totalCapital := gs.config.OrderAmount.Mul(decimal.NewFromInt(int64(gs.config.GridLevels)))

	// Convert to StrategyPerformance
	return &StrategyPerformance{
		TotalInvested:  totalCapital,
		CurrentValue:   totalCapital.Add(totalProfit),
		UnrealizedPnL:  totalProfit,
		ROI:            totalProfit.Div(totalCapital).Mul(decimal.NewFromInt(100)),
		AveragePrice:   currentPrice,
		TotalTokens:    decimal.NewFromInt(int64(totalBuyOrders + totalSellOrders)),
		ExecutionCount: totalBuyOrders + totalSellOrders,
		LastExecution:  time.Now(),
	}
}

// GridPerformance represents performance metrics for Grid strategy
type GridPerformance struct {
	ActiveLevels    int             `json:"active_levels"`
	TotalBuyOrders  int             `json:"total_buy_orders"`
	TotalSellOrders int             `json:"total_sell_orders"`
	TotalProfit     decimal.Decimal `json:"total_profit"`
	CurrentPrice    decimal.Decimal `json:"current_price"`
	BasePrice       decimal.Decimal `json:"base_price"`
	GridRange       string          `json:"grid_range"`
}

// Reset resets the Grid strategy state
func (gs *GridStrategy) Reset() {
	gs.gridLevels = make([]*GridLevel, 0)
	gs.activeOrders = make(map[string]*GridOrder)
	gs.basePrice = decimal.Zero
	gs.isInitialized = false
}

// GetConfig returns the strategy configuration
func (gs *GridStrategy) GetConfig() *GridConfig {
	return gs.config
}

// GetName returns the strategy name
func (gs *GridStrategy) GetName() string {
	return "Grid Trading"
}

// GetType returns the strategy type
func (gs *GridStrategy) GetType() string {
	return "Grid"
}

// Validate validates the strategy configuration
func (gs *GridStrategy) Validate() error {
	if gs.config.GridLevels < 3 {
		return fmt.Errorf("grid levels must be at least 3")
	}

	if gs.config.GridSpacing.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("grid spacing must be positive")
	}

	if gs.config.UpperBound.LessThanOrEqual(decimal.NewFromFloat(1.0)) {
		return fmt.Errorf("upper bound must be greater than 1.0")
	}

	if gs.config.LowerBound.GreaterThanOrEqual(decimal.NewFromFloat(1.0)) {
		return fmt.Errorf("lower bound must be less than 1.0")
	}

	if gs.config.OrderAmount.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("order amount must be positive")
	}

	return nil
}
