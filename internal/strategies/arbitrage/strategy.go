package arbitrage

import (
	"context"
	"fmt"
	"time"

	"github.com/ai-agentic-browser/internal/exchanges/common"
	"github.com/ai-agentic-browser/internal/strategies/framework"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ArbitrageStrategy implements a cross-exchange arbitrage trading strategy
type ArbitrageStrategy struct {
	*framework.BaseStrategy
	
	// Strategy-specific fields
	targetSymbol     string
	minProfitPercent decimal.Decimal
	maxOrderSize     decimal.Decimal
	
	// Market data tracking
	exchangePrices   map[string]*ExchangePrice
	lastUpdateTime   map[string]time.Time
	
	// Configuration
	config ArbitrageConfig
}

// ArbitrageConfig contains arbitrage strategy configuration
type ArbitrageConfig struct {
	Symbol           string          `json:"symbol"`
	MinProfitPercent decimal.Decimal `json:"min_profit_percent"` // Minimum profit percentage to execute
	MaxOrderSize     decimal.Decimal `json:"max_order_size"`     // Maximum order size
	MaxLatency       time.Duration   `json:"max_latency"`        // Maximum acceptable latency
	MinVolume        decimal.Decimal `json:"min_volume"`         // Minimum volume required
	Exchanges        []string        `json:"exchanges"`          // Exchanges to monitor
	CooldownPeriod   time.Duration   `json:"cooldown_period"`    // Time to wait between trades
}

// ExchangePrice tracks price data for an exchange
type ExchangePrice struct {
	Exchange    string          `json:"exchange"`
	BidPrice    decimal.Decimal `json:"bid_price"`
	AskPrice    decimal.Decimal `json:"ask_price"`
	BidVolume   decimal.Decimal `json:"bid_volume"`
	AskVolume   decimal.Decimal `json:"ask_volume"`
	Timestamp   time.Time       `json:"timestamp"`
	Latency     time.Duration   `json:"latency"`
}

// ArbitrageOpportunity represents a detected arbitrage opportunity
type ArbitrageOpportunity struct {
	BuyExchange   string          `json:"buy_exchange"`
	SellExchange  string          `json:"sell_exchange"`
	BuyPrice      decimal.Decimal `json:"buy_price"`
	SellPrice     decimal.Decimal `json:"sell_price"`
	Quantity      decimal.Decimal `json:"quantity"`
	ProfitPercent decimal.Decimal `json:"profit_percent"`
	ProfitAmount  decimal.Decimal `json:"profit_amount"`
	Timestamp     time.Time       `json:"timestamp"`
}

// NewArbitrageStrategy creates a new arbitrage strategy
func NewArbitrageStrategy(config ArbitrageConfig) *ArbitrageStrategy {
	strategy := &ArbitrageStrategy{
		BaseStrategy: framework.NewBaseStrategy(
			"Cross-Exchange Arbitrage",
			"Exploits price differences between exchanges for risk-free profit",
			"1.0.0",
		),
		config:         config,
		exchangePrices: make(map[string]*ExchangePrice),
		lastUpdateTime: make(map[string]time.Time),
	}

	// Initialize strategy parameters
	strategy.initializeParameters()
	
	// Set default risk limits
	strategy.setDefaultRiskLimits()

	return strategy
}

// Initialize initializes the strategy with configuration
func (as *ArbitrageStrategy) Initialize(ctx context.Context, config framework.StrategyConfig) error {
	if err := as.BaseStrategy.Initialize(ctx, config); err != nil {
		return fmt.Errorf("failed to initialize base strategy: %w", err)
	}

	// Extract strategy-specific configuration
	if symbol, ok := config.Parameters["symbol"].(string); ok {
		as.targetSymbol = symbol
		as.config.Symbol = symbol
	}

	if minProfitPercent, ok := config.Parameters["min_profit_percent"].(decimal.Decimal); ok {
		as.minProfitPercent = minProfitPercent
		as.config.MinProfitPercent = minProfitPercent
	}

	if maxOrderSize, ok := config.Parameters["max_order_size"].(decimal.Decimal); ok {
		as.maxOrderSize = maxOrderSize
		as.config.MaxOrderSize = maxOrderSize
	}

	return nil
}

// OnMarketData processes market data and generates trading signals
func (as *ArbitrageStrategy) OnMarketData(ctx context.Context, data *framework.MarketData) ([]*framework.Signal, error) {
	if !as.IsRunning() {
		return nil, nil
	}

	// Only process ticker data for our target symbol
	if data.Symbol != as.targetSymbol || data.Type != framework.MarketDataTypeTicker {
		return nil, nil
	}

	tickerData, ok := data.Data.(*common.TickerData)
	if !ok {
		return nil, fmt.Errorf("invalid ticker data type")
	}

	// Update exchange price data
	as.updateExchangePrice(data.Exchange, tickerData)

	// Look for arbitrage opportunities
	opportunities := as.findArbitrageOpportunities()

	// Generate signals for profitable opportunities
	var signals []*framework.Signal
	for _, opportunity := range opportunities {
		if opportunity.ProfitPercent.GreaterThanOrEqual(as.minProfitPercent) {
			signals = append(signals, as.createArbitrageSignals(opportunity)...)
		}
	}

	return signals, nil
}

// OnOrderUpdate handles order status updates
func (as *ArbitrageStrategy) OnOrderUpdate(ctx context.Context, update *framework.OrderUpdate) error {
	if !as.IsRunning() {
		return nil
	}

	// Update metrics for completed trades
	if update.Status == common.OrderStatusFilled {
		trade := &framework.Trade{
			Symbol:     update.Symbol,
			Side:       string(update.Side),
			Quantity:   update.FilledQty,
			Price:      update.AvgPrice,
			Commission: update.Commission,
			Timestamp:  update.Timestamp,
		}
		as.UpdateMetrics(trade)
	}

	return nil
}

// OnPositionUpdate handles position updates
func (as *ArbitrageStrategy) OnPositionUpdate(ctx context.Context, update *framework.PositionUpdate) error {
	// Arbitrage strategies typically don't hold positions long-term
	return nil
}

// Private methods

// initializeParameters sets up strategy parameters
func (as *ArbitrageStrategy) initializeParameters() {
	as.AddParameter(framework.Parameter{
		Name:         "symbol",
		Type:         framework.ParamTypeString,
		DefaultValue: "BTCUSDT",
		Description:  "Trading symbol",
		Required:     true,
	})

	as.AddParameter(framework.Parameter{
		Name:         "min_profit_percent",
		Type:         framework.ParamTypeDecimal,
		DefaultValue: decimal.NewFromFloat(0.002), // 0.2%
		Description:  "Minimum profit percentage to execute trade",
		Required:     true,
	})

	as.AddParameter(framework.Parameter{
		Name:         "max_order_size",
		Type:         framework.ParamTypeDecimal,
		DefaultValue: decimal.NewFromFloat(1.0),
		Description:  "Maximum order size",
		Required:     true,
	})

	as.AddParameter(framework.Parameter{
		Name:         "max_latency",
		Type:         framework.ParamTypeString,
		DefaultValue: "100ms",
		Description:  "Maximum acceptable latency",
		Required:     false,
	})
}

// setDefaultRiskLimits sets default risk management limits
func (as *ArbitrageStrategy) setDefaultRiskLimits() {
	limits := &framework.RiskLimits{
		MaxPositionSize:    decimal.NewFromFloat(5.0),
		MaxDailyLoss:       decimal.NewFromFloat(500.0),
		MaxDrawdown:        decimal.NewFromFloat(200.0),
		MaxOpenPositions:   10,
		MaxOrdersPerSecond: 5,
		MaxOrdersPerMinute: 50,
		MaxOrdersPerHour:   500,
		MaxOrdersPerDay:    2000,
		StopLossRequired:   false,
		TakeProfitRequired: false,
		AllowedSymbols:     []string{"BTCUSDT", "ETHUSDT", "ADAUSDT"},
		BlockedSymbols:     []string{},
	}

	as.SetRiskLimits(limits)
}

// updateExchangePrice updates price data for an exchange
func (as *ArbitrageStrategy) updateExchangePrice(exchange string, ticker *common.TickerData) {
	as.exchangePrices[exchange] = &ExchangePrice{
		Exchange:  exchange,
		BidPrice:  ticker.BidPrice,
		AskPrice:  ticker.AskPrice,
		BidVolume: ticker.BidQty,
		AskVolume: ticker.AskQty,
		Timestamp: ticker.Timestamp,
	}
	as.lastUpdateTime[exchange] = time.Now()
}

// findArbitrageOpportunities scans for arbitrage opportunities across exchanges
func (as *ArbitrageStrategy) findArbitrageOpportunities() []*ArbitrageOpportunity {
	var opportunities []*ArbitrageOpportunity

	// Need at least 2 exchanges to find arbitrage
	if len(as.exchangePrices) < 2 {
		return opportunities
	}

	// Check all exchange pairs for arbitrage opportunities
	exchanges := make([]string, 0, len(as.exchangePrices))
	for exchange := range as.exchangePrices {
		exchanges = append(exchanges, exchange)
	}

	for i := 0; i < len(exchanges); i++ {
		for j := i + 1; j < len(exchanges); j++ {
			exchange1 := exchanges[i]
			exchange2 := exchanges[j]

			price1 := as.exchangePrices[exchange1]
			price2 := as.exchangePrices[exchange2]

			// Check if data is recent enough
			if time.Since(price1.Timestamp) > as.config.MaxLatency ||
				time.Since(price2.Timestamp) > as.config.MaxLatency {
				continue
			}

			// Check opportunity: buy on exchange1, sell on exchange2
			opportunity1 := as.calculateOpportunity(exchange1, exchange2, price1, price2)
			if opportunity1 != nil {
				opportunities = append(opportunities, opportunity1)
			}

			// Check opportunity: buy on exchange2, sell on exchange1
			opportunity2 := as.calculateOpportunity(exchange2, exchange1, price2, price1)
			if opportunity2 != nil {
				opportunities = append(opportunities, opportunity2)
			}
		}
	}

	return opportunities
}

// calculateOpportunity calculates arbitrage opportunity between two exchanges
func (as *ArbitrageStrategy) calculateOpportunity(buyExchange, sellExchange string, buyPrice, sellPrice *ExchangePrice) *ArbitrageOpportunity {
	// Buy at ask price on buy exchange, sell at bid price on sell exchange
	buyPx := buyPrice.AskPrice
	sellPx := sellPrice.BidPrice

	// Check if there's a profit opportunity
	if sellPx.LessThanOrEqual(buyPx) {
		return nil
	}

	// Calculate profit
	profitPerUnit := sellPx.Sub(buyPx)
	profitPercent := profitPerUnit.Div(buyPx).Mul(decimal.NewFromInt(100))

	// Determine maximum quantity based on available volume
	maxQuantity := decimal.Min(buyPrice.AskVolume, sellPrice.BidVolume)
	if maxQuantity.GreaterThan(as.maxOrderSize) {
		maxQuantity = as.maxOrderSize
	}

	// Check minimum volume requirement
	if maxQuantity.LessThan(as.config.MinVolume) {
		return nil
	}

	totalProfit := profitPerUnit.Mul(maxQuantity)

	return &ArbitrageOpportunity{
		BuyExchange:   buyExchange,
		SellExchange:  sellExchange,
		BuyPrice:      buyPx,
		SellPrice:     sellPx,
		Quantity:      maxQuantity,
		ProfitPercent: profitPercent,
		ProfitAmount:  totalProfit,
		Timestamp:     time.Now(),
	}
}

// createArbitrageSignals creates trading signals for an arbitrage opportunity
func (as *ArbitrageStrategy) createArbitrageSignals(opportunity *ArbitrageOpportunity) []*framework.Signal {
	var signals []*framework.Signal

	// Create buy signal for the cheaper exchange
	buySignal := &framework.Signal{
		ID:          uuid.New(),
		StrategyID:  as.GetID(),
		Type:        framework.SignalTypeBuy,
		Symbol:      as.targetSymbol,
		Side:        common.OrderSideBuy,
		Strength:    decimal.NewFromFloat(1.0), // High strength for arbitrage
		Confidence:  decimal.NewFromFloat(0.95),
		Price:       opportunity.BuyPrice,
		Quantity:    opportunity.Quantity,
		TimeInForce: common.TimeInForceIOC, // Immediate or cancel for arbitrage
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(5 * time.Second), // Short expiry for arbitrage
		Metadata: map[string]interface{}{
			"strategy_type":    "arbitrage",
			"exchange":         opportunity.BuyExchange,
			"counterpart_exchange": opportunity.SellExchange,
			"profit_percent":   opportunity.ProfitPercent.String(),
			"profit_amount":    opportunity.ProfitAmount.String(),
		},
	}

	// Create sell signal for the more expensive exchange
	sellSignal := &framework.Signal{
		ID:          uuid.New(),
		StrategyID:  as.GetID(),
		Type:        framework.SignalTypeSell,
		Symbol:      as.targetSymbol,
		Side:        common.OrderSideSell,
		Strength:    decimal.NewFromFloat(1.0), // High strength for arbitrage
		Confidence:  decimal.NewFromFloat(0.95),
		Price:       opportunity.SellPrice,
		Quantity:    opportunity.Quantity,
		TimeInForce: common.TimeInForceIOC, // Immediate or cancel for arbitrage
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(5 * time.Second), // Short expiry for arbitrage
		Metadata: map[string]interface{}{
			"strategy_type":    "arbitrage",
			"exchange":         opportunity.SellExchange,
			"counterpart_exchange": opportunity.BuyExchange,
			"profit_percent":   opportunity.ProfitPercent.String(),
			"profit_amount":    opportunity.ProfitAmount.String(),
		},
	}

	signals = append(signals, buySignal, sellSignal)
	return signals
}
