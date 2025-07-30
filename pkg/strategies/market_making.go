package strategies

import (
	"context"
	"sync"
	"time"

	"github.com/ai-agentic-browser/internal/hft"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// MarketMakingStrategy implements a market making trading strategy
type MarketMakingStrategy struct {
	logger      *observability.Logger
	config      StrategyConfig
	parameters  MarketMakingParams
	performance *StrategyPerformance

	// Market making state
	activeOrders map[string]*ActiveOrder
	inventory    map[string]decimal.Decimal
	lastQuotes   map[string]*Quote

	// Risk management
	maxInventory decimal.Decimal
	maxSpread    decimal.Decimal
	minSpread    decimal.Decimal

	// State management
	enabled bool
	mu      sync.RWMutex
}

// MarketMakingParams contains parameters for market making strategy
type MarketMakingParams struct {
	SpreadBps            int             `json:"spread_bps"`            // Spread in basis points
	OrderSize            decimal.Decimal `json:"order_size"`            // Size of each order
	MaxInventory         decimal.Decimal `json:"max_inventory"`         // Maximum inventory per symbol
	InventoryTarget      decimal.Decimal `json:"inventory_target"`      // Target inventory (usually 0)
	SkewFactor           decimal.Decimal `json:"skew_factor"`           // Inventory skew factor
	MinSpreadBps         int             `json:"min_spread_bps"`        // Minimum spread
	MaxSpreadBps         int             `json:"max_spread_bps"`        // Maximum spread
	OrderRefreshMs       int64           `json:"order_refresh_ms"`      // Order refresh interval
	RiskAdjustment       bool            `json:"risk_adjustment"`       // Enable risk-based adjustments
	VolatilityAdjustment bool            `json:"volatility_adjustment"` // Enable volatility-based adjustments
}

// ActiveOrder represents an active market making order
type ActiveOrder struct {
	ID         uuid.UUID       `json:"id"`
	Symbol     string          `json:"symbol"`
	Side       hft.OrderSide   `json:"side"`
	Price      decimal.Decimal `json:"price"`
	Quantity   decimal.Decimal `json:"quantity"`
	CreatedAt  time.Time       `json:"created_at"`
	LastUpdate time.Time       `json:"last_update"`
}

// Quote represents a bid/ask quote
type Quote struct {
	Symbol    string          `json:"symbol"`
	BidPrice  decimal.Decimal `json:"bid_price"`
	AskPrice  decimal.Decimal `json:"ask_price"`
	BidSize   decimal.Decimal `json:"bid_size"`
	AskSize   decimal.Decimal `json:"ask_size"`
	MidPrice  decimal.Decimal `json:"mid_price"`
	Spread    decimal.Decimal `json:"spread"`
	Timestamp time.Time       `json:"timestamp"`
}

// NewMarketMakingStrategy creates a new market making strategy
func NewMarketMakingStrategy(logger *observability.Logger, config StrategyConfig) *MarketMakingStrategy {
	// Set default parameters
	params := MarketMakingParams{
		SpreadBps:            10, // 10 bps default spread
		OrderSize:            decimal.NewFromFloat(0.1),
		MaxInventory:         decimal.NewFromFloat(1.0),
		InventoryTarget:      decimal.Zero,
		SkewFactor:           decimal.NewFromFloat(0.1),
		MinSpreadBps:         5,
		MaxSpreadBps:         50,
		OrderRefreshMs:       1000, // 1 second
		RiskAdjustment:       true,
		VolatilityAdjustment: true,
	}

	// Override with config parameters if provided
	if config.Parameters != nil {
		if spreadBps, ok := config.Parameters["spread_bps"].(int); ok {
			params.SpreadBps = spreadBps
		}
		if orderSize, ok := config.Parameters["order_size"].(float64); ok {
			params.OrderSize = decimal.NewFromFloat(orderSize)
		}
		if maxInventory, ok := config.Parameters["max_inventory"].(float64); ok {
			params.MaxInventory = decimal.NewFromFloat(maxInventory)
		}
	}

	strategy := &MarketMakingStrategy{
		logger:       logger,
		config:       config,
		parameters:   params,
		activeOrders: make(map[string]*ActiveOrder),
		inventory:    make(map[string]decimal.Decimal),
		lastQuotes:   make(map[string]*Quote),
		maxInventory: params.MaxInventory,
		maxSpread:    decimal.NewFromFloat(float64(params.MaxSpreadBps) / 10000),
		minSpread:    decimal.NewFromFloat(float64(params.MinSpreadBps) / 10000),
		enabled:      config.Enabled,
		performance: &StrategyPerformance{
			StrategyID: config.ID,
			LastUpdate: time.Now(),
		},
	}

	return strategy
}

// GetID returns the strategy ID
func (mms *MarketMakingStrategy) GetID() string {
	return mms.config.ID
}

// GetName returns the strategy name
func (mms *MarketMakingStrategy) GetName() string {
	return mms.config.Name
}

// GetType returns the strategy type
func (mms *MarketMakingStrategy) GetType() StrategyType {
	return StrategyTypeMarketMaking
}

// GetConfig returns the strategy configuration
func (mms *MarketMakingStrategy) GetConfig() StrategyConfig {
	return mms.config
}

// Initialize initializes the strategy
func (mms *MarketMakingStrategy) Initialize(ctx context.Context) error {
	mms.logger.Info(ctx, "Initializing market making strategy", map[string]interface{}{
		"strategy_id":   mms.config.ID,
		"symbols":       mms.config.Symbols,
		"spread_bps":    mms.parameters.SpreadBps,
		"order_size":    mms.parameters.OrderSize.String(),
		"max_inventory": mms.parameters.MaxInventory.String(),
	})

	// Initialize inventory for all symbols
	for _, symbol := range mms.config.Symbols {
		mms.inventory[symbol] = decimal.Zero
	}

	return nil
}

// ProcessTick processes a market tick and generates trading signals
func (mms *MarketMakingStrategy) ProcessTick(ctx context.Context, tick hft.MarketTick) ([]hft.TradingSignal, error) {
	mms.mu.Lock()
	defer mms.mu.Unlock()

	if !mms.enabled {
		return nil, nil
	}

	// Update quote information
	quote := mms.updateQuote(tick)
	if quote == nil {
		return nil, nil
	}

	// Generate market making signals
	signals := mms.generateMarketMakingSignals(ctx, quote)

	return signals, nil
}

// UpdateParameters updates strategy parameters
func (mms *MarketMakingStrategy) UpdateParameters(params map[string]interface{}) error {
	mms.mu.Lock()
	defer mms.mu.Unlock()

	if spreadBps, ok := params["spread_bps"].(int); ok {
		mms.parameters.SpreadBps = spreadBps
	}

	if orderSize, ok := params["order_size"].(float64); ok {
		mms.parameters.OrderSize = decimal.NewFromFloat(orderSize)
	}

	if maxInventory, ok := params["max_inventory"].(float64); ok {
		mms.parameters.MaxInventory = decimal.NewFromFloat(maxInventory)
		mms.maxInventory = mms.parameters.MaxInventory
	}

	mms.config.UpdatedAt = time.Now()

	return nil
}

// GetPerformance returns strategy performance metrics
func (mms *MarketMakingStrategy) GetPerformance() *StrategyPerformance {
	mms.mu.RLock()
	defer mms.mu.RUnlock()

	return mms.performance
}

// IsEnabled returns whether the strategy is enabled
func (mms *MarketMakingStrategy) IsEnabled() bool {
	mms.mu.RLock()
	defer mms.mu.RUnlock()

	return mms.enabled
}

// SetEnabled sets the strategy enabled state
func (mms *MarketMakingStrategy) SetEnabled(enabled bool) {
	mms.mu.Lock()
	defer mms.mu.Unlock()

	mms.enabled = enabled
	mms.config.Enabled = enabled
}

// Cleanup cleans up strategy resources
func (mms *MarketMakingStrategy) Cleanup(ctx context.Context) error {
	mms.mu.Lock()
	defer mms.mu.Unlock()

	mms.logger.Info(ctx, "Cleaning up market making strategy", map[string]interface{}{
		"strategy_id":   mms.config.ID,
		"active_orders": len(mms.activeOrders),
	})

	// Clear active orders
	mms.activeOrders = make(map[string]*ActiveOrder)

	return nil
}

// updateQuote updates quote information from market tick
func (mms *MarketMakingStrategy) updateQuote(tick hft.MarketTick) *Quote {
	if tick.BidPrice.IsZero() || tick.AskPrice.IsZero() {
		return nil
	}

	midPrice := tick.BidPrice.Add(tick.AskPrice).Div(decimal.NewFromInt(2))
	spread := tick.AskPrice.Sub(tick.BidPrice)

	quote := &Quote{
		Symbol:    tick.Symbol,
		BidPrice:  tick.BidPrice,
		AskPrice:  tick.AskPrice,
		BidSize:   tick.BidSize,
		AskSize:   tick.AskSize,
		MidPrice:  midPrice,
		Spread:    spread,
		Timestamp: tick.Timestamp,
	}

	mms.lastQuotes[tick.Symbol] = quote

	return quote
}

// generateMarketMakingSignals generates buy and sell signals for market making
func (mms *MarketMakingStrategy) generateMarketMakingSignals(ctx context.Context, quote *Quote) []hft.TradingSignal {
	var signals []hft.TradingSignal

	// Calculate target spread
	targetSpread := mms.calculateTargetSpread(quote)

	// Calculate inventory skew
	inventorySkew := mms.calculateInventorySkew(quote.Symbol)

	// Calculate bid and ask prices
	bidPrice, askPrice := mms.calculateQuotePrices(quote, targetSpread, inventorySkew)

	// Generate bid signal
	if mms.shouldPlaceBid(quote.Symbol, bidPrice) {
		bidSignal := hft.TradingSignal{
			ID:         uuid.New(),
			Symbol:     quote.Symbol,
			Side:       hft.OrderSideBuy,
			OrderType:  hft.OrderTypeLimit,
			Quantity:   mms.calculateOrderSize(quote.Symbol, hft.OrderSideBuy),
			Price:      bidPrice,
			Confidence: 0.8,
			StrategyID: mms.config.ID,
			Timestamp:  time.Now(),
			Metadata: map[string]interface{}{
				"strategy_type":  "market_making",
				"side":           "bid",
				"target_spread":  targetSpread.String(),
				"inventory_skew": inventorySkew.String(),
			},
		}
		signals = append(signals, bidSignal)
	}

	// Generate ask signal
	if mms.shouldPlaceAsk(quote.Symbol, askPrice) {
		askSignal := hft.TradingSignal{
			ID:         uuid.New(),
			Symbol:     quote.Symbol,
			Side:       hft.OrderSideSell,
			OrderType:  hft.OrderTypeLimit,
			Quantity:   mms.calculateOrderSize(quote.Symbol, hft.OrderSideSell),
			Price:      askPrice,
			Confidence: 0.8,
			StrategyID: mms.config.ID,
			Timestamp:  time.Now(),
			Metadata: map[string]interface{}{
				"strategy_type":  "market_making",
				"side":           "ask",
				"target_spread":  targetSpread.String(),
				"inventory_skew": inventorySkew.String(),
			},
		}
		signals = append(signals, askSignal)
	}

	return signals
}

// calculateTargetSpread calculates the target spread based on market conditions
func (mms *MarketMakingStrategy) calculateTargetSpread(quote *Quote) decimal.Decimal {
	baseSpread := decimal.NewFromFloat(float64(mms.parameters.SpreadBps) / 10000)

	// Adjust for volatility if enabled
	if mms.parameters.VolatilityAdjustment {
		// Simple volatility adjustment based on current spread
		currentSpreadRatio := quote.Spread.Div(quote.MidPrice)
		if currentSpreadRatio.GreaterThan(mms.maxSpread) {
			baseSpread = baseSpread.Mul(decimal.NewFromFloat(1.5))
		} else if currentSpreadRatio.LessThan(mms.minSpread) {
			baseSpread = baseSpread.Mul(decimal.NewFromFloat(0.8))
		}
	}

	// Ensure spread is within bounds
	if baseSpread.LessThan(mms.minSpread) {
		baseSpread = mms.minSpread
	}
	if baseSpread.GreaterThan(mms.maxSpread) {
		baseSpread = mms.maxSpread
	}

	return baseSpread
}

// calculateInventorySkew calculates inventory skew adjustment
func (mms *MarketMakingStrategy) calculateInventorySkew(symbol string) decimal.Decimal {
	if !mms.parameters.RiskAdjustment {
		return decimal.Zero
	}

	currentInventory := mms.inventory[symbol]
	targetInventory := mms.parameters.InventoryTarget

	inventoryDiff := currentInventory.Sub(targetInventory)
	maxInventory := mms.parameters.MaxInventory

	if maxInventory.IsZero() {
		return decimal.Zero
	}

	// Calculate skew as percentage of max inventory
	skewRatio := inventoryDiff.Div(maxInventory)
	skew := skewRatio.Mul(mms.parameters.SkewFactor)

	return skew
}

// calculateQuotePrices calculates bid and ask prices
func (mms *MarketMakingStrategy) calculateQuotePrices(quote *Quote, targetSpread, inventorySkew decimal.Decimal) (decimal.Decimal, decimal.Decimal) {
	halfSpread := targetSpread.Div(decimal.NewFromInt(2))

	// Apply inventory skew
	bidAdjustment := inventorySkew
	askAdjustment := inventorySkew.Neg()

	bidPrice := quote.MidPrice.Sub(halfSpread).Add(bidAdjustment.Mul(quote.MidPrice))
	askPrice := quote.MidPrice.Add(halfSpread).Add(askAdjustment.Mul(quote.MidPrice))

	return bidPrice, askPrice
}

// shouldPlaceBid determines if a bid order should be placed
func (mms *MarketMakingStrategy) shouldPlaceBid(symbol string, bidPrice decimal.Decimal) bool {
	// Check inventory limits
	currentInventory := mms.inventory[symbol]
	if currentInventory.GreaterThanOrEqual(mms.maxInventory) {
		return false
	}

	// Check if we already have a similar bid order
	for _, order := range mms.activeOrders {
		if order.Symbol == symbol && order.Side == hft.OrderSideBuy {
			priceDiff := order.Price.Sub(bidPrice).Abs()
			if priceDiff.LessThan(bidPrice.Mul(decimal.NewFromFloat(0.001))) { // Within 0.1%
				return false
			}
		}
	}

	return true
}

// shouldPlaceAsk determines if an ask order should be placed
func (mms *MarketMakingStrategy) shouldPlaceAsk(symbol string, askPrice decimal.Decimal) bool {
	// Check inventory limits
	currentInventory := mms.inventory[symbol]
	if currentInventory.LessThanOrEqual(mms.maxInventory.Neg()) {
		return false
	}

	// Check if we already have a similar ask order
	for _, order := range mms.activeOrders {
		if order.Symbol == symbol && order.Side == hft.OrderSideSell {
			priceDiff := order.Price.Sub(askPrice).Abs()
			if priceDiff.LessThan(askPrice.Mul(decimal.NewFromFloat(0.001))) { // Within 0.1%
				return false
			}
		}
	}

	return true
}

// calculateOrderSize calculates the order size for a given side
func (mms *MarketMakingStrategy) calculateOrderSize(symbol string, side hft.OrderSide) decimal.Decimal {
	baseSize := mms.parameters.OrderSize

	// Adjust size based on inventory if risk adjustment is enabled
	if mms.parameters.RiskAdjustment {
		currentInventory := mms.inventory[symbol]
		inventoryRatio := currentInventory.Div(mms.maxInventory).Abs()

		// Reduce size as inventory approaches limits
		if inventoryRatio.GreaterThan(decimal.NewFromFloat(0.8)) {
			sizeReduction := inventoryRatio.Sub(decimal.NewFromFloat(0.8)).Mul(decimal.NewFromFloat(0.5))
			baseSize = baseSize.Mul(decimal.NewFromInt(1).Sub(sizeReduction))
		}
	}

	// Ensure minimum size
	minSize := decimal.NewFromFloat(0.001)
	if baseSize.LessThan(minSize) {
		baseSize = minSize
	}

	return baseSize
}

// UpdateInventory updates inventory for a symbol (called when orders are filled)
func (mms *MarketMakingStrategy) UpdateInventory(symbol string, quantity decimal.Decimal, side hft.OrderSide) {
	mms.mu.Lock()
	defer mms.mu.Unlock()

	if side == hft.OrderSideBuy {
		mms.inventory[symbol] = mms.inventory[symbol].Add(quantity)
	} else {
		mms.inventory[symbol] = mms.inventory[symbol].Sub(quantity)
	}

	mms.logger.Debug(context.Background(), "Inventory updated", map[string]interface{}{
		"strategy_id":   mms.config.ID,
		"symbol":        symbol,
		"quantity":      quantity.String(),
		"side":          string(side),
		"new_inventory": mms.inventory[symbol].String(),
	})
}

// GetInventory returns current inventory for all symbols
func (mms *MarketMakingStrategy) GetInventory() map[string]decimal.Decimal {
	mms.mu.RLock()
	defer mms.mu.RUnlock()

	inventory := make(map[string]decimal.Decimal)
	for symbol, qty := range mms.inventory {
		inventory[symbol] = qty
	}

	return inventory
}

// GetActiveOrders returns current active orders
func (mms *MarketMakingStrategy) GetActiveOrders() map[string]*ActiveOrder {
	mms.mu.RLock()
	defer mms.mu.RUnlock()

	orders := make(map[string]*ActiveOrder)
	for id, order := range mms.activeOrders {
		orders[id] = order
	}

	return orders
}
