package hft

import (
	"context"
	"fmt"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// OrderRequest represents an order request for routing
type OrderRequest struct {
	ID          uuid.UUID
	Symbol      string
	Side        OrderSide
	Quantity    decimal.Decimal
	Type        OrderType
	Price       decimal.Decimal
	TimeInForce TimeInForce
	ClientID    string
	StrategyID  string
	Metadata    map[string]interface{}
}

// VenueSelector selects optimal venues for order execution
type VenueSelector struct {
	logger *observability.Logger
	config SORConfig
}

// LiquidityAggregator aggregates liquidity across venues
type LiquidityAggregator struct {
	logger *observability.Logger
	config SORConfig
}

// OrderSplitter splits orders across multiple venues
type OrderSplitter struct {
	logger *observability.Logger
	config SORConfig
}

// ExecutionEngine manages order execution across venues
type ExecutionEngine struct {
	logger *observability.Logger
	config SORConfig
}

// PreTradeRiskEngine performs pre-trade risk checks
type PreTradeRiskEngine struct {
	logger *observability.Logger
	config SORConfig
}

// ComplianceEngine performs compliance checks
type ComplianceEngine struct {
	logger *observability.Logger
	config SORConfig
}

// NewVenueSelector creates a new venue selector
func NewVenueSelector(logger *observability.Logger, config SORConfig) *VenueSelector {
	return &VenueSelector{
		logger: logger,
		config: config,
	}
}

// NewLiquidityAggregator creates a new liquidity aggregator
func NewLiquidityAggregator(logger *observability.Logger, config SORConfig) *LiquidityAggregator {
	return &LiquidityAggregator{
		logger: logger,
		config: config,
	}
}

// NewOrderSplitter creates a new order splitter
func NewOrderSplitter(logger *observability.Logger, config SORConfig) *OrderSplitter {
	return &OrderSplitter{
		logger: logger,
		config: config,
	}
}

// NewExecutionEngine creates a new execution engine
func NewExecutionEngine(logger *observability.Logger, config SORConfig) *ExecutionEngine {
	return &ExecutionEngine{
		logger: logger,
		config: config,
	}
}

// NewPreTradeRiskEngine creates a new pre-trade risk engine
func NewPreTradeRiskEngine(logger *observability.Logger, config SORConfig) *PreTradeRiskEngine {
	return &PreTradeRiskEngine{
		logger: logger,
		config: config,
	}
}

// NewComplianceEngine creates a new compliance engine
func NewComplianceEngine(logger *observability.Logger, config SORConfig) *ComplianceEngine {
	return &ComplianceEngine{
		logger: logger,
		config: config,
	}
}

// SelectVenues selects optimal venues for an order
func (vs *VenueSelector) SelectVenues(ctx context.Context, order *OrderRequest, marketData *BestPriceAggregation) ([]*TradingVenue, error) {
	// This is a simplified implementation
	// In production, this would use sophisticated venue selection algorithms

	var selectedVenues []*TradingVenue

	// Mock venue creation for demonstration
	venues := []*TradingVenue{
		{
			ID:           "binance",
			Name:         "Binance",
			Type:         VenueTypeExchange,
			Exchange:     "binance",
			IsConnected:  true,
			LastPing:     time.Millisecond,
			Reliability:  0.99,
			BestBid:      decimal.NewFromFloat(45000.0),
			BestAsk:      decimal.NewFromFloat(45001.0),
			BidSize:      decimal.NewFromFloat(10.0),
			AskSize:      decimal.NewFromFloat(8.0),
			LastUpdate:   time.Now(),
			MinOrderSize: decimal.NewFromFloat(0.001),
			MaxOrderSize: decimal.NewFromFloat(1000.0),
			TickSize:     decimal.NewFromFloat(0.01),
			AvgFillRate:  0.95,
			AvgLatency:   2 * time.Millisecond,
			RejectRate:   0.01,
		},
		{
			ID:           "coinbase",
			Name:         "Coinbase Pro",
			Type:         VenueTypeExchange,
			Exchange:     "coinbase",
			IsConnected:  true,
			LastPing:     2 * time.Millisecond,
			Reliability:  0.98,
			BestBid:      decimal.NewFromFloat(44999.5),
			BestAsk:      decimal.NewFromFloat(45001.5),
			BidSize:      decimal.NewFromFloat(5.0),
			AskSize:      decimal.NewFromFloat(12.0),
			LastUpdate:   time.Now(),
			MinOrderSize: decimal.NewFromFloat(0.001),
			MaxOrderSize: decimal.NewFromFloat(500.0),
			TickSize:     decimal.NewFromFloat(0.01),
			AvgFillRate:  0.93,
			AvgLatency:   3 * time.Millisecond,
			RejectRate:   0.02,
		},
		{
			ID:              "dark_pool_1",
			Name:            "Institutional Dark Pool",
			Type:            VenueTypeDarkPool,
			Exchange:        "dark_pool",
			IsConnected:     true,
			LastPing:        time.Millisecond,
			Reliability:     0.97,
			BestBid:         decimal.NewFromFloat(45000.0),
			BestAsk:         decimal.NewFromFloat(45001.0),
			BidSize:         decimal.NewFromFloat(0.0), // Hidden
			AskSize:         decimal.NewFromFloat(0.0), // Hidden
			LastUpdate:      time.Now(),
			MinOrderSize:    decimal.NewFromFloat(1.0),
			MaxOrderSize:    decimal.NewFromFloat(100.0),
			TickSize:        decimal.NewFromFloat(0.01),
			AvgFillRate:     0.85,
			AvgLatency:      5 * time.Millisecond,
			RejectRate:      0.05,
			IsDarkPool:      true,
			HiddenLiquidity: decimal.NewFromFloat(50.0),
		},
	}

	// Filter venues based on order requirements
	for _, venue := range venues {
		if vs.isVenueSuitable(venue, order) {
			selectedVenues = append(selectedVenues, venue)
		}
	}

	return selectedVenues, nil
}

// isVenueSuitable checks if a venue is suitable for an order
func (vs *VenueSelector) isVenueSuitable(venue *TradingVenue, order *OrderRequest) bool {
	// Check connectivity
	if !venue.IsConnected {
		return false
	}

	// Check order size constraints
	if order.Quantity.LessThan(venue.MinOrderSize) || order.Quantity.GreaterThan(venue.MaxOrderSize) {
		return false
	}

	// Check reliability threshold
	if venue.Reliability < 0.9 {
		return false
	}

	// Check latency requirements
	if venue.AvgLatency > 10*time.Millisecond {
		return false
	}

	return true
}

// ValidateOrder performs pre-trade risk validation
func (pre *PreTradeRiskEngine) ValidateOrder(ctx context.Context, order *OrderRequest) error {
	// Check order size limits
	if order.Quantity.GreaterThan(decimal.NewFromFloat(1000.0)) {
		return fmt.Errorf("order quantity exceeds maximum limit")
	}

	// Check price reasonableness (simplified)
	if order.Type == OrderTypeLimit && !order.Price.IsZero() {
		if order.Price.LessThan(decimal.NewFromFloat(1.0)) || order.Price.GreaterThan(decimal.NewFromFloat(1000000.0)) {
			return fmt.Errorf("order price is outside reasonable range")
		}
	}

	return nil
}

// ValidateOrder performs compliance checks
func (ce *ComplianceEngine) ValidateOrder(ctx context.Context, order *OrderRequest) error {
	// Check for wash trading (simplified)
	if order.ClientID == "" {
		return fmt.Errorf("client ID is required for compliance")
	}

	// Check for market manipulation patterns (simplified)
	// In production, this would check against historical patterns

	return nil
}

// SplitOrder splits an order into child orders based on allocations
func (os *OrderSplitter) SplitOrder(ctx context.Context, order *OrderRequest, allocations []VenueAllocation) ([]ChildOrder, error) {
	var childOrders []ChildOrder

	for i, allocation := range allocations {
		childOrder := ChildOrder{
			ID:          uuid.New(),
			ParentID:    order.ID,
			VenueID:     allocation.VenueID,
			Symbol:      order.Symbol,
			Side:        order.Side,
			Quantity:    allocation.Quantity,
			Price:       allocation.ExpectedPrice,
			OrderType:   allocation.OrderType,
			TimeInForce: allocation.TimeInForce,
			Status:      OrderStatusNew,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		childOrders = append(childOrders, childOrder)

		os.logger.Info(ctx, "Created child order", map[string]interface{}{
			"child_id":  childOrder.ID.String(),
			"parent_id": order.ID.String(),
			"venue_id":  allocation.VenueID,
			"quantity":  allocation.Quantity.String(),
			"price":     allocation.ExpectedPrice.String(),
			"priority":  allocation.Priority,
			"index":     i,
		})
	}

	return childOrders, nil
}
