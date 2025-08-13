package market_making

import (
	"context"
	"fmt"
	"time"

	"github.com/ai-agentic-browser/internal/exchanges/common"
	"github.com/ai-agentic-browser/internal/strategies/framework"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// MarketMakingStrategy implements a market making trading strategy
type MarketMakingStrategy struct {
	*framework.BaseStrategy
	
	// Strategy-specific fields
	targetSymbol    string
	spreadPercent   decimal.Decimal
	orderSize       decimal.Decimal
	maxPosition     decimal.Decimal
	inventoryTarget decimal.Decimal
	
	// Current state
	currentBid      decimal.Decimal
	currentAsk      decimal.Decimal
	currentSpread   decimal.Decimal
	inventory       decimal.Decimal
	lastOrderTime   time.Time
	
	// Configuration
	config MarketMakingConfig
}

// MarketMakingConfig contains market making strategy configuration
type MarketMakingConfig struct {
	Symbol              string          `json:"symbol"`
	SpreadPercent       decimal.Decimal `json:"spread_percent"`       // Target spread as percentage
	OrderSize           decimal.Decimal `json:"order_size"`           // Size of each order
	MaxPosition         decimal.Decimal `json:"max_position"`         // Maximum position size
	InventoryTarget     decimal.Decimal `json:"inventory_target"`     // Target inventory level
	MinSpread           decimal.Decimal `json:"min_spread"`           // Minimum spread to maintain
	MaxSpread           decimal.Decimal `json:"max_spread"`           // Maximum spread allowed
	OrderRefreshTime    time.Duration   `json:"order_refresh_time"`   // How often to refresh orders
	RiskAdjustment      bool            `json:"risk_adjustment"`      // Enable inventory risk adjustment
	VolatilityAdjustment bool           `json:"volatility_adjustment"` // Enable volatility-based adjustment
}

// NewMarketMakingStrategy creates a new market making strategy
func NewMarketMakingStrategy(config MarketMakingConfig) *MarketMakingStrategy {
	strategy := &MarketMakingStrategy{
		BaseStrategy: framework.NewBaseStrategy(
			"Market Making",
			"Provides liquidity by placing bid and ask orders with a target spread",
			"1.0.0",
		),
		config: config,
	}

	// Initialize strategy parameters
	strategy.initializeParameters()
	
	// Set default risk limits
	strategy.setDefaultRiskLimits()

	return strategy
}

// Initialize initializes the strategy with configuration
func (mms *MarketMakingStrategy) Initialize(ctx context.Context, config framework.StrategyConfig) error {
	if err := mms.BaseStrategy.Initialize(ctx, config); err != nil {
		return fmt.Errorf("failed to initialize base strategy: %w", err)
	}

	// Extract strategy-specific configuration
	if symbol, ok := config.Parameters["symbol"].(string); ok {
		mms.targetSymbol = symbol
		mms.config.Symbol = symbol
	}

	if spreadPercent, ok := config.Parameters["spread_percent"].(decimal.Decimal); ok {
		mms.spreadPercent = spreadPercent
		mms.config.SpreadPercent = spreadPercent
	}

	if orderSize, ok := config.Parameters["order_size"].(decimal.Decimal); ok {
		mms.orderSize = orderSize
		mms.config.OrderSize = orderSize
	}

	if maxPosition, ok := config.Parameters["max_position"].(decimal.Decimal); ok {
		mms.maxPosition = maxPosition
		mms.config.MaxPosition = maxPosition
	}

	return nil
}

// OnMarketData processes market data and generates trading signals
func (mms *MarketMakingStrategy) OnMarketData(ctx context.Context, data *framework.MarketData) ([]*framework.Signal, error) {
	if !mms.IsRunning() {
		return nil, nil
	}

	// Only process ticker and order book data for our target symbol
	if data.Symbol != mms.targetSymbol {
		return nil, nil
	}

	var signals []*framework.Signal

	switch data.Type {
	case framework.MarketDataTypeTicker:
		tickerData, ok := data.Data.(*common.TickerData)
		if !ok {
			return nil, fmt.Errorf("invalid ticker data type")
		}
		signals = mms.processTickerData(ctx, tickerData)

	case framework.MarketDataTypeOrderBook:
		orderBookData, ok := data.Data.(*common.OrderBookData)
		if !ok {
			return nil, fmt.Errorf("invalid order book data type")
		}
		signals = mms.processOrderBookData(ctx, orderBookData)
	}

	return signals, nil
}

// OnOrderUpdate handles order status updates
func (mms *MarketMakingStrategy) OnOrderUpdate(ctx context.Context, update *framework.OrderUpdate) error {
	if !mms.IsRunning() {
		return nil
	}

	// Update inventory based on filled orders
	if update.Status == common.OrderStatusFilled || update.Status == common.OrderStatusPartiallyFilled {
		if update.Side == common.OrderSideBuy {
			mms.inventory = mms.inventory.Add(update.FilledQty)
		} else {
			mms.inventory = mms.inventory.Sub(update.FilledQty)
		}

		// Update metrics
		trade := &framework.Trade{
			Symbol:     update.Symbol,
			Side:       string(update.Side),
			Quantity:   update.FilledQty,
			Price:      update.AvgPrice,
			Commission: update.Commission,
			Timestamp:  update.Timestamp,
		}
		mms.UpdateMetrics(trade)
	}

	return nil
}

// OnPositionUpdate handles position updates
func (mms *MarketMakingStrategy) OnPositionUpdate(ctx context.Context, update *framework.PositionUpdate) error {
	if !mms.IsRunning() {
		return nil
	}

	// Update inventory from position data
	mms.inventory = update.Size

	return nil
}

// Private methods

// initializeParameters sets up strategy parameters
func (mms *MarketMakingStrategy) initializeParameters() {
	mms.AddParameter(framework.Parameter{
		Name:         "symbol",
		Type:         framework.ParamTypeString,
		DefaultValue: "BTCUSDT",
		Description:  "Trading symbol",
		Required:     true,
	})

	mms.AddParameter(framework.Parameter{
		Name:         "spread_percent",
		Type:         framework.ParamTypeDecimal,
		DefaultValue: decimal.NewFromFloat(0.001), // 0.1%
		Description:  "Target spread as percentage",
		Required:     true,
	})

	mms.AddParameter(framework.Parameter{
		Name:         "order_size",
		Type:         framework.ParamTypeDecimal,
		DefaultValue: decimal.NewFromFloat(0.01),
		Description:  "Size of each order",
		Required:     true,
	})

	mms.AddParameter(framework.Parameter{
		Name:         "max_position",
		Type:         framework.ParamTypeDecimal,
		DefaultValue: decimal.NewFromFloat(1.0),
		Description:  "Maximum position size",
		Required:     true,
	})

	mms.AddParameter(framework.Parameter{
		Name:         "inventory_target",
		Type:         framework.ParamTypeDecimal,
		DefaultValue: decimal.NewFromInt(0),
		Description:  "Target inventory level",
		Required:     false,
	})
}

// setDefaultRiskLimits sets default risk management limits
func (mms *MarketMakingStrategy) setDefaultRiskLimits() {
	limits := &framework.RiskLimits{
		MaxPositionSize:    decimal.NewFromFloat(10.0),
		MaxDailyLoss:       decimal.NewFromFloat(1000.0),
		MaxDrawdown:        decimal.NewFromFloat(500.0),
		MaxOpenPositions:   5,
		MaxOrdersPerSecond: 10,
		MaxOrdersPerMinute: 100,
		MaxOrdersPerHour:   1000,
		MaxOrdersPerDay:    10000,
		StopLossRequired:   false,
		TakeProfitRequired: false,
		AllowedSymbols:     []string{"BTCUSDT", "ETHUSDT", "ADAUSDT"},
		BlockedSymbols:     []string{},
	}

	mms.SetRiskLimits(limits)
}

// processTickerData processes ticker data and generates signals
func (mms *MarketMakingStrategy) processTickerData(ctx context.Context, ticker *common.TickerData) []*framework.Signal {
	var signals []*framework.Signal

	// Update current market prices
	mms.currentBid = ticker.BidPrice
	mms.currentAsk = ticker.AskPrice
	mms.currentSpread = ticker.AskPrice.Sub(ticker.BidPrice)

	// Calculate target bid and ask prices
	midPrice := ticker.BidPrice.Add(ticker.AskPrice).Div(decimal.NewFromInt(2))
	targetSpread := midPrice.Mul(mms.spreadPercent)

	// Adjust for inventory risk
	inventoryAdjustment := mms.calculateInventoryAdjustment()
	
	targetBid := midPrice.Sub(targetSpread.Div(decimal.NewFromInt(2))).Add(inventoryAdjustment)
	targetAsk := midPrice.Add(targetSpread.Div(decimal.NewFromInt(2))).Add(inventoryAdjustment)

	// Check if we need to place new orders
	if mms.shouldPlaceOrders() {
		// Generate buy signal if we're not at max long position
		if mms.inventory.Add(mms.orderSize).LessThanOrEqual(mms.maxPosition) {
			buySignal := &framework.Signal{
				ID:          uuid.New(),
				StrategyID:  mms.GetID(),
				Type:        framework.SignalTypeBuy,
				Symbol:      mms.targetSymbol,
				Side:        common.OrderSideBuy,
				Strength:    decimal.NewFromFloat(0.8),
				Confidence:  decimal.NewFromFloat(0.9),
				Price:       targetBid,
				Quantity:    mms.orderSize,
				TimeInForce: common.TimeInForceGTC,
				CreatedAt:   time.Now(),
				Metadata: map[string]interface{}{
					"strategy_type": "market_making",
					"order_type":    "bid",
				},
			}
			signals = append(signals, buySignal)
		}

		// Generate sell signal if we're not at max short position
		if mms.inventory.Sub(mms.orderSize).GreaterThanOrEqual(mms.maxPosition.Neg()) {
			sellSignal := &framework.Signal{
				ID:          uuid.New(),
				StrategyID:  mms.GetID(),
				Type:        framework.SignalTypeSell,
				Symbol:      mms.targetSymbol,
				Side:        common.OrderSideSell,
				Strength:    decimal.NewFromFloat(0.8),
				Confidence:  decimal.NewFromFloat(0.9),
				Price:       targetAsk,
				Quantity:    mms.orderSize,
				TimeInForce: common.TimeInForceGTC,
				CreatedAt:   time.Now(),
				Metadata: map[string]interface{}{
					"strategy_type": "market_making",
					"order_type":    "ask",
				},
			}
			signals = append(signals, sellSignal)
		}

		mms.lastOrderTime = time.Now()
	}

	return signals
}

// processOrderBookData processes order book data
func (mms *MarketMakingStrategy) processOrderBookData(ctx context.Context, orderBook *common.OrderBookData) []*framework.Signal {
	// For now, we primarily use ticker data for market making
	// Order book data could be used for more sophisticated strategies
	return nil
}

// calculateInventoryAdjustment calculates price adjustment based on current inventory
func (mms *MarketMakingStrategy) calculateInventoryAdjustment() decimal.Decimal {
	if !mms.config.RiskAdjustment {
		return decimal.NewFromInt(0)
	}

	// Simple inventory risk adjustment
	// If we have too much inventory, adjust prices to encourage selling
	// If we have too little inventory, adjust prices to encourage buying
	
	inventoryRatio := mms.inventory.Div(mms.maxPosition)
	maxAdjustment := mms.currentSpread.Mul(decimal.NewFromFloat(0.1)) // 10% of spread
	
	return inventoryRatio.Mul(maxAdjustment)
}

// shouldPlaceOrders determines if new orders should be placed
func (mms *MarketMakingStrategy) shouldPlaceOrders() bool {
	// Check if enough time has passed since last order
	if time.Since(mms.lastOrderTime) < mms.config.OrderRefreshTime {
		return false
	}

	// Check if current spread is within acceptable range
	if mms.currentSpread.LessThan(mms.config.MinSpread) {
		return false
	}

	if mms.currentSpread.GreaterThan(mms.config.MaxSpread) {
		return false
	}

	return true
}
