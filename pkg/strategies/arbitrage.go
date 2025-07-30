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

// ArbitrageStrategy implements an arbitrage trading strategy
type ArbitrageStrategy struct {
	logger      *observability.Logger
	config      StrategyConfig
	parameters  ArbitrageParams
	performance *StrategyPerformance

	// Arbitrage state
	exchangePrices map[string]map[string]*ExchangePrice // [symbol][exchange] -> price
	opportunities  map[string]*ArbitrageOpportunity
	activeArbs     map[string]*ActiveArbitrage

	// State management
	enabled bool
	mu      sync.RWMutex
}

// ArbitrageParams contains parameters for arbitrage strategy
type ArbitrageParams struct {
	MinProfitBps        int             `json:"min_profit_bps"`        // Minimum profit in basis points
	MaxPositionSize     decimal.Decimal `json:"max_position_size"`     // Maximum position size
	ExecutionTimeoutMs  int64           `json:"execution_timeout_ms"`  // Execution timeout
	TransactionCostBps  int             `json:"transaction_cost_bps"`  // Transaction cost in bps
	MinVolumeUSD        decimal.Decimal `json:"min_volume_usd"`        // Minimum volume in USD
	MaxLatencyMs        int64           `json:"max_latency_ms"`        // Maximum acceptable latency
	EnableTriangular    bool            `json:"enable_triangular"`     // Enable triangular arbitrage
	EnableCrossExchange bool            `json:"enable_cross_exchange"` // Enable cross-exchange arbitrage
	RiskAdjustment      decimal.Decimal `json:"risk_adjustment"`       // Risk adjustment factor
}

// ExchangePrice represents price data from an exchange
type ExchangePrice struct {
	Exchange  string          `json:"exchange"`
	Symbol    string          `json:"symbol"`
	BidPrice  decimal.Decimal `json:"bid_price"`
	AskPrice  decimal.Decimal `json:"ask_price"`
	BidSize   decimal.Decimal `json:"bid_size"`
	AskSize   decimal.Decimal `json:"ask_size"`
	Timestamp time.Time       `json:"timestamp"`
	Latency   time.Duration   `json:"latency"`
}

// ArbitrageOpportunity represents a detected arbitrage opportunity
type ArbitrageOpportunity struct {
	ID               uuid.UUID       `json:"id"`
	Type             ArbitrageType   `json:"type"`
	Symbol           string          `json:"symbol"`
	BuyExchange      string          `json:"buy_exchange"`
	SellExchange     string          `json:"sell_exchange"`
	BuyPrice         decimal.Decimal `json:"buy_price"`
	SellPrice        decimal.Decimal `json:"sell_price"`
	ProfitBps        int             `json:"profit_bps"`
	ProfitUSD        decimal.Decimal `json:"profit_usd"`
	MaxQuantity      decimal.Decimal `json:"max_quantity"`
	EstimatedLatency time.Duration   `json:"estimated_latency"`
	Confidence       float64         `json:"confidence"`
	DetectedAt       time.Time       `json:"detected_at"`
	ExpiresAt        time.Time       `json:"expires_at"`
}

// ActiveArbitrage represents an active arbitrage trade
type ActiveArbitrage struct {
	ID             uuid.UUID          `json:"id"`
	OpportunityID  uuid.UUID          `json:"opportunity_id"`
	Symbol         string             `json:"symbol"`
	Quantity       decimal.Decimal    `json:"quantity"`
	BuyOrder       *hft.TradingSignal `json:"buy_order"`
	SellOrder      *hft.TradingSignal `json:"sell_order"`
	Status         ArbitrageStatus    `json:"status"`
	ExpectedProfit decimal.Decimal    `json:"expected_profit"`
	ActualProfit   decimal.Decimal    `json:"actual_profit"`
	CreatedAt      time.Time          `json:"created_at"`
	CompletedAt    *time.Time         `json:"completed_at,omitempty"`
}

// ArbitrageType represents different types of arbitrage
type ArbitrageType string

const (
	ArbitrageTypeSimple      ArbitrageType = "SIMPLE"      // Simple price difference
	ArbitrageTypeTriangular  ArbitrageType = "TRIANGULAR"  // Triangular arbitrage
	ArbitrageTypeCross       ArbitrageType = "CROSS"       // Cross-exchange arbitrage
	ArbitrageTypeStatistical ArbitrageType = "STATISTICAL" // Statistical arbitrage
)

// ArbitrageStatus represents arbitrage execution status
type ArbitrageStatus string

const (
	ArbitrageStatusPending   ArbitrageStatus = "PENDING"
	ArbitrageStatusExecuting ArbitrageStatus = "EXECUTING"
	ArbitrageStatusCompleted ArbitrageStatus = "COMPLETED"
	ArbitrageStatusFailed    ArbitrageStatus = "FAILED"
	ArbitrageStatusExpired   ArbitrageStatus = "EXPIRED"
)

// NewArbitrageStrategy creates a new arbitrage strategy
func NewArbitrageStrategy(logger *observability.Logger, config StrategyConfig) *ArbitrageStrategy {
	// Set default parameters
	params := ArbitrageParams{
		MinProfitBps:        20, // 20 bps minimum profit
		MaxPositionSize:     decimal.NewFromFloat(1.0),
		ExecutionTimeoutMs:  5000, // 5 seconds
		TransactionCostBps:  5,    // 5 bps transaction cost
		MinVolumeUSD:        decimal.NewFromFloat(1000),
		MaxLatencyMs:        100, // 100ms max latency
		EnableTriangular:    true,
		EnableCrossExchange: true,
		RiskAdjustment:      decimal.NewFromFloat(0.8),
	}

	// Override with config parameters if provided
	if config.Parameters != nil {
		if minProfitBps, ok := config.Parameters["min_profit_bps"].(int); ok {
			params.MinProfitBps = minProfitBps
		}
		if maxPositionSize, ok := config.Parameters["max_position_size"].(float64); ok {
			params.MaxPositionSize = decimal.NewFromFloat(maxPositionSize)
		}
		if transactionCostBps, ok := config.Parameters["transaction_cost_bps"].(int); ok {
			params.TransactionCostBps = transactionCostBps
		}
	}

	strategy := &ArbitrageStrategy{
		logger:         logger,
		config:         config,
		parameters:     params,
		exchangePrices: make(map[string]map[string]*ExchangePrice),
		opportunities:  make(map[string]*ArbitrageOpportunity),
		activeArbs:     make(map[string]*ActiveArbitrage),
		enabled:        config.Enabled,
		performance: &StrategyPerformance{
			StrategyID: config.ID,
			LastUpdate: time.Now(),
		},
	}

	return strategy
}

// GetID returns the strategy ID
func (as *ArbitrageStrategy) GetID() string {
	return as.config.ID
}

// GetName returns the strategy name
func (as *ArbitrageStrategy) GetName() string {
	return as.config.Name
}

// GetType returns the strategy type
func (as *ArbitrageStrategy) GetType() StrategyType {
	return StrategyTypeArbitrage
}

// GetConfig returns the strategy configuration
func (as *ArbitrageStrategy) GetConfig() StrategyConfig {
	return as.config
}

// Initialize initializes the strategy
func (as *ArbitrageStrategy) Initialize(ctx context.Context) error {
	as.logger.Info(ctx, "Initializing arbitrage strategy", map[string]interface{}{
		"strategy_id":           as.config.ID,
		"symbols":               as.config.Symbols,
		"min_profit_bps":        as.parameters.MinProfitBps,
		"enable_triangular":     as.parameters.EnableTriangular,
		"enable_cross_exchange": as.parameters.EnableCrossExchange,
	})

	// Initialize price tracking for all symbols
	for _, symbol := range as.config.Symbols {
		as.exchangePrices[symbol] = make(map[string]*ExchangePrice)
	}

	return nil
}

// ProcessTick processes a market tick and generates arbitrage signals
func (as *ArbitrageStrategy) ProcessTick(ctx context.Context, tick hft.MarketTick) ([]hft.TradingSignal, error) {
	as.mu.Lock()
	defer as.mu.Unlock()

	if !as.enabled {
		return nil, nil
	}

	// Update exchange price data
	as.updateExchangePrice(tick)

	// Detect arbitrage opportunities
	opportunities := as.detectArbitrageOpportunities(ctx, tick.Symbol)

	// Generate signals for valid opportunities
	var signals []hft.TradingSignal
	for _, opportunity := range opportunities {
		if as.isOpportunityValid(opportunity) {
			arbSignals := as.generateArbitrageSignals(ctx, opportunity)
			signals = append(signals, arbSignals...)
		}
	}

	return signals, nil
}

// UpdateParameters updates strategy parameters
func (as *ArbitrageStrategy) UpdateParameters(params map[string]interface{}) error {
	as.mu.Lock()
	defer as.mu.Unlock()

	if minProfitBps, ok := params["min_profit_bps"].(int); ok {
		as.parameters.MinProfitBps = minProfitBps
	}

	if maxPositionSize, ok := params["max_position_size"].(float64); ok {
		as.parameters.MaxPositionSize = decimal.NewFromFloat(maxPositionSize)
	}

	if transactionCostBps, ok := params["transaction_cost_bps"].(int); ok {
		as.parameters.TransactionCostBps = transactionCostBps
	}

	as.config.UpdatedAt = time.Now()

	return nil
}

// GetPerformance returns strategy performance metrics
func (as *ArbitrageStrategy) GetPerformance() *StrategyPerformance {
	as.mu.RLock()
	defer as.mu.RUnlock()

	return as.performance
}

// IsEnabled returns whether the strategy is enabled
func (as *ArbitrageStrategy) IsEnabled() bool {
	as.mu.RLock()
	defer as.mu.RUnlock()

	return as.enabled
}

// SetEnabled sets the strategy enabled state
func (as *ArbitrageStrategy) SetEnabled(enabled bool) {
	as.mu.Lock()
	defer as.mu.Unlock()

	as.enabled = enabled
	as.config.Enabled = enabled
}

// Cleanup cleans up strategy resources
func (as *ArbitrageStrategy) Cleanup(ctx context.Context) error {
	as.mu.Lock()
	defer as.mu.Unlock()

	as.logger.Info(ctx, "Cleaning up arbitrage strategy", map[string]interface{}{
		"strategy_id":       as.config.ID,
		"active_arbitrages": len(as.activeArbs),
		"opportunities":     len(as.opportunities),
	})

	// Clear active arbitrages and opportunities
	as.activeArbs = make(map[string]*ActiveArbitrage)
	as.opportunities = make(map[string]*ArbitrageOpportunity)

	return nil
}

// updateExchangePrice updates price data for an exchange
func (as *ArbitrageStrategy) updateExchangePrice(tick hft.MarketTick) {
	if as.exchangePrices[tick.Symbol] == nil {
		as.exchangePrices[tick.Symbol] = make(map[string]*ExchangePrice)
	}

	exchangePrice := &ExchangePrice{
		Exchange:  tick.Exchange,
		Symbol:    tick.Symbol,
		BidPrice:  tick.BidPrice,
		AskPrice:  tick.AskPrice,
		BidSize:   tick.BidSize,
		AskSize:   tick.AskSize,
		Timestamp: tick.Timestamp,
		Latency:   time.Since(tick.Timestamp),
	}

	as.exchangePrices[tick.Symbol][tick.Exchange] = exchangePrice
}

// detectArbitrageOpportunities detects arbitrage opportunities for a symbol
func (as *ArbitrageStrategy) detectArbitrageOpportunities(ctx context.Context, symbol string) []*ArbitrageOpportunity {
	var opportunities []*ArbitrageOpportunity

	exchangePrices := as.exchangePrices[symbol]
	if len(exchangePrices) < 2 {
		return opportunities // Need at least 2 exchanges
	}

	// Simple arbitrage detection
	if as.parameters.EnableCrossExchange {
		simpleOpps := as.detectSimpleArbitrage(symbol, exchangePrices)
		opportunities = append(opportunities, simpleOpps...)
	}

	// Triangular arbitrage detection
	if as.parameters.EnableTriangular {
		triangularOpps := as.detectTriangularArbitrage(symbol, exchangePrices)
		opportunities = append(opportunities, triangularOpps...)
	}

	return opportunities
}

// detectSimpleArbitrage detects simple price difference arbitrage
func (as *ArbitrageStrategy) detectSimpleArbitrage(symbol string, exchangePrices map[string]*ExchangePrice) []*ArbitrageOpportunity {
	var opportunities []*ArbitrageOpportunity

	// Compare all exchange pairs
	exchanges := make([]*ExchangePrice, 0, len(exchangePrices))
	for _, price := range exchangePrices {
		if as.isPriceDataValid(price) {
			exchanges = append(exchanges, price)
		}
	}

	for i := 0; i < len(exchanges); i++ {
		for j := i + 1; j < len(exchanges); j++ {
			exchange1 := exchanges[i]
			exchange2 := exchanges[j]

			// Check if we can buy on exchange1 and sell on exchange2
			if exchange2.BidPrice.GreaterThan(exchange1.AskPrice) {
				profit := exchange2.BidPrice.Sub(exchange1.AskPrice)
				profitBps := profit.Div(exchange1.AskPrice).Mul(decimal.NewFromInt(10000))

				if profitBps.IntPart() >= int64(as.parameters.MinProfitBps+as.parameters.TransactionCostBps) {
					opportunity := &ArbitrageOpportunity{
						ID:           uuid.New(),
						Type:         ArbitrageTypeSimple,
						Symbol:       symbol,
						BuyExchange:  exchange1.Exchange,
						SellExchange: exchange2.Exchange,
						BuyPrice:     exchange1.AskPrice,
						SellPrice:    exchange2.BidPrice,
						ProfitBps:    int(profitBps.IntPart()),
						MaxQuantity:  as.calculateMaxQuantity(exchange1.AskSize, exchange2.BidSize),
						DetectedAt:   time.Now(),
						ExpiresAt:    time.Now().Add(time.Duration(as.parameters.ExecutionTimeoutMs) * time.Millisecond),
						Confidence:   as.calculateConfidence(exchange1, exchange2),
					}

					opportunity.ProfitUSD = profit.Mul(opportunity.MaxQuantity)
					opportunities = append(opportunities, opportunity)
				}
			}

			// Check if we can buy on exchange2 and sell on exchange1
			if exchange1.BidPrice.GreaterThan(exchange2.AskPrice) {
				profit := exchange1.BidPrice.Sub(exchange2.AskPrice)
				profitBps := profit.Div(exchange2.AskPrice).Mul(decimal.NewFromInt(10000))

				if profitBps.IntPart() >= int64(as.parameters.MinProfitBps+as.parameters.TransactionCostBps) {
					opportunity := &ArbitrageOpportunity{
						ID:           uuid.New(),
						Type:         ArbitrageTypeSimple,
						Symbol:       symbol,
						BuyExchange:  exchange2.Exchange,
						SellExchange: exchange1.Exchange,
						BuyPrice:     exchange2.AskPrice,
						SellPrice:    exchange1.BidPrice,
						ProfitBps:    int(profitBps.IntPart()),
						MaxQuantity:  as.calculateMaxQuantity(exchange2.AskSize, exchange1.BidSize),
						DetectedAt:   time.Now(),
						ExpiresAt:    time.Now().Add(time.Duration(as.parameters.ExecutionTimeoutMs) * time.Millisecond),
						Confidence:   as.calculateConfidence(exchange2, exchange1),
					}

					opportunity.ProfitUSD = profit.Mul(opportunity.MaxQuantity)
					opportunities = append(opportunities, opportunity)
				}
			}
		}
	}

	return opportunities
}

// detectTriangularArbitrage detects triangular arbitrage opportunities
func (as *ArbitrageStrategy) detectTriangularArbitrage(symbol string, exchangePrices map[string]*ExchangePrice) []*ArbitrageOpportunity {
	// Simplified triangular arbitrage detection
	// In a real implementation, this would be more sophisticated
	var opportunities []*ArbitrageOpportunity

	// For now, return empty slice
	// TODO: Implement triangular arbitrage detection

	return opportunities
}

// isPriceDataValid checks if price data is valid and recent
func (as *ArbitrageStrategy) isPriceDataValid(price *ExchangePrice) bool {
	if price == nil {
		return false
	}

	// Check if data is recent
	if time.Since(price.Timestamp) > time.Duration(as.parameters.MaxLatencyMs)*time.Millisecond {
		return false
	}

	// Check if prices are valid
	if price.BidPrice.IsZero() || price.AskPrice.IsZero() {
		return false
	}

	// Check if spread is reasonable
	spread := price.AskPrice.Sub(price.BidPrice)
	if spread.LessThanOrEqual(decimal.Zero) {
		return false
	}

	return true
}

// calculateMaxQuantity calculates the maximum quantity for arbitrage
func (as *ArbitrageStrategy) calculateMaxQuantity(buySize, sellSize decimal.Decimal) decimal.Decimal {
	maxQuantity := buySize
	if sellSize.LessThan(maxQuantity) {
		maxQuantity = sellSize
	}

	// Apply position size limit
	if maxQuantity.GreaterThan(as.parameters.MaxPositionSize) {
		maxQuantity = as.parameters.MaxPositionSize
	}

	return maxQuantity
}

// calculateConfidence calculates confidence score for an opportunity
func (as *ArbitrageStrategy) calculateConfidence(buyExchange, sellExchange *ExchangePrice) float64 {
	baseConfidence := 0.8

	// Adjust based on latency
	avgLatency := (buyExchange.Latency + sellExchange.Latency) / 2
	if avgLatency > time.Duration(as.parameters.MaxLatencyMs/2)*time.Millisecond {
		baseConfidence -= 0.2
	}

	// Adjust based on size
	minSize := buyExchange.AskSize
	if sellExchange.BidSize.LessThan(minSize) {
		minSize = sellExchange.BidSize
	}

	if minSize.LessThan(as.parameters.MinVolumeUSD.Div(buyExchange.AskPrice)) {
		baseConfidence -= 0.1
	}

	if baseConfidence < 0.1 {
		baseConfidence = 0.1
	}

	return baseConfidence
}

// isOpportunityValid validates an arbitrage opportunity
func (as *ArbitrageStrategy) isOpportunityValid(opportunity *ArbitrageOpportunity) bool {
	// Check if opportunity has expired
	if time.Now().After(opportunity.ExpiresAt) {
		return false
	}

	// Check minimum profit
	if opportunity.ProfitBps < as.parameters.MinProfitBps {
		return false
	}

	// Check minimum volume
	if opportunity.ProfitUSD.LessThan(as.parameters.MinVolumeUSD) {
		return false
	}

	// Check confidence
	if opportunity.Confidence < 0.5 {
		return false
	}

	return true
}

// generateArbitrageSignals generates trading signals for an arbitrage opportunity
func (as *ArbitrageStrategy) generateArbitrageSignals(ctx context.Context, opportunity *ArbitrageOpportunity) []hft.TradingSignal {
	var signals []hft.TradingSignal

	// Calculate execution quantity
	quantity := opportunity.MaxQuantity.Mul(as.parameters.RiskAdjustment)

	// Create buy signal
	buySignal := hft.TradingSignal{
		ID:         uuid.New(),
		Symbol:     opportunity.Symbol,
		Side:       hft.OrderSideBuy,
		OrderType:  hft.OrderTypeMarket, // Use market orders for speed
		Quantity:   quantity,
		Price:      opportunity.BuyPrice,
		Confidence: opportunity.Confidence,
		StrategyID: as.config.ID,
		Timestamp:  time.Now(),
		Metadata: map[string]interface{}{
			"strategy_type":       "arbitrage",
			"arbitrage_type":      string(opportunity.Type),
			"opportunity_id":      opportunity.ID.String(),
			"buy_exchange":        opportunity.BuyExchange,
			"sell_exchange":       opportunity.SellExchange,
			"expected_profit_bps": opportunity.ProfitBps,
			"expected_profit_usd": opportunity.ProfitUSD.String(),
		},
	}

	// Create sell signal
	sellSignal := hft.TradingSignal{
		ID:         uuid.New(),
		Symbol:     opportunity.Symbol,
		Side:       hft.OrderSideSell,
		OrderType:  hft.OrderTypeMarket,
		Quantity:   quantity,
		Price:      opportunity.SellPrice,
		Confidence: opportunity.Confidence,
		StrategyID: as.config.ID,
		Timestamp:  time.Now(),
		Metadata: map[string]interface{}{
			"strategy_type":       "arbitrage",
			"arbitrage_type":      string(opportunity.Type),
			"opportunity_id":      opportunity.ID.String(),
			"buy_exchange":        opportunity.BuyExchange,
			"sell_exchange":       opportunity.SellExchange,
			"expected_profit_bps": opportunity.ProfitBps,
			"expected_profit_usd": opportunity.ProfitUSD.String(),
		},
	}

	signals = append(signals, buySignal, sellSignal)

	// Track active arbitrage
	activeArb := &ActiveArbitrage{
		ID:             uuid.New(),
		OpportunityID:  opportunity.ID,
		Symbol:         opportunity.Symbol,
		Quantity:       quantity,
		BuyOrder:       &buySignal,
		SellOrder:      &sellSignal,
		Status:         ArbitrageStatusPending,
		ExpectedProfit: opportunity.ProfitUSD,
		CreatedAt:      time.Now(),
	}

	as.activeArbs[activeArb.ID.String()] = activeArb
	as.opportunities[opportunity.ID.String()] = opportunity

	as.logger.Info(ctx, "Arbitrage opportunity detected", map[string]interface{}{
		"opportunity_id": opportunity.ID.String(),
		"symbol":         opportunity.Symbol,
		"type":           string(opportunity.Type),
		"profit_bps":     opportunity.ProfitBps,
		"profit_usd":     opportunity.ProfitUSD.String(),
		"buy_exchange":   opportunity.BuyExchange,
		"sell_exchange":  opportunity.SellExchange,
		"quantity":       quantity.String(),
	})

	return signals
}

// GetOpportunities returns current arbitrage opportunities
func (as *ArbitrageStrategy) GetOpportunities() map[string]*ArbitrageOpportunity {
	as.mu.RLock()
	defer as.mu.RUnlock()

	opportunities := make(map[string]*ArbitrageOpportunity)
	for id, opp := range as.opportunities {
		opportunities[id] = opp
	}

	return opportunities
}

// GetActiveArbitrages returns current active arbitrages
func (as *ArbitrageStrategy) GetActiveArbitrages() map[string]*ActiveArbitrage {
	as.mu.RLock()
	defer as.mu.RUnlock()

	arbitrages := make(map[string]*ActiveArbitrage)
	for id, arb := range as.activeArbs {
		arbitrages[id] = arb
	}

	return arbitrages
}
