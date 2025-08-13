package hft

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// SmartOrderRouter provides intelligent order routing with liquidity aggregation,
// venue selection algorithms, and pre-trade risk checks across multiple exchanges
type SmartOrderRouter struct {
	logger *observability.Logger
	config SORConfig

	// Exchange connections and venue management
	venues        map[string]*TradingVenue
	venueSelector *VenueSelector
	liquidityAggr *LiquidityAggregator

	// Order management
	orderSplitter   *OrderSplitter
	executionEngine *ExecutionEngine

	// Risk and compliance
	preTradeRisk    *PreTradeRiskEngine
	complianceCheck *ComplianceEngine

	// Performance tracking
	ordersRouted      int64
	venuesUsed        int64
	avgExecutionTime  int64
	slippageReduction float64

	// State management
	isRunning int32
	stopChan  chan struct{}
	wg        sync.WaitGroup
	mu        sync.RWMutex

	// Real-time market data
	marketDataFeeds map[string]chan *VenueMarketData
	bestPrices      map[string]*BestPriceAggregation
}

// SORConfig contains configuration for Smart Order Routing
type SORConfig struct {
	// Venue configuration
	EnabledVenues     []string        `json:"enabled_venues"`
	MaxVenuesPerOrder int             `json:"max_venues_per_order"`
	MinVenueSize      decimal.Decimal `json:"min_venue_size"`

	// Routing algorithms
	RoutingAlgorithm string        `json:"routing_algorithm"` // TWAP, VWAP, Implementation Shortfall
	SlippageTarget   float64       `json:"slippage_target"`   // Target slippage in basis points
	TimeHorizon      time.Duration `json:"time_horizon"`      // Execution time horizon

	// Liquidity aggregation
	EnableDarkPools       bool            `json:"enable_dark_pools"`
	DarkPoolParticipation float64         `json:"dark_pool_participation"` // 0.0 to 1.0
	MinLiquidityThreshold decimal.Decimal `json:"min_liquidity_threshold"`

	// Risk controls
	MaxOrderSize    decimal.Decimal `json:"max_order_size"`
	MaxMarketImpact float64         `json:"max_market_impact"` // Maximum allowed market impact
	EnablePreTrade  bool            `json:"enable_pre_trade"`  // Enable pre-trade risk checks

	// Performance optimization
	LatencyTarget   time.Duration `json:"latency_target"`   // Target routing latency
	CacheTimeout    time.Duration `json:"cache_timeout"`    // Price cache timeout
	RefreshInterval time.Duration `json:"refresh_interval"` // Market data refresh
}

// TradingVenue represents a trading venue (exchange, dark pool, etc.)
type TradingVenue struct {
	ID       string
	Name     string
	Type     VenueType
	Exchange string

	// Connectivity
	IsConnected bool
	LastPing    time.Duration
	Reliability float64 // 0.0 to 1.0

	// Market data
	BestBid    decimal.Decimal
	BestAsk    decimal.Decimal
	BidSize    decimal.Decimal
	AskSize    decimal.Decimal
	LastUpdate time.Time

	// Venue characteristics
	MinOrderSize decimal.Decimal
	MaxOrderSize decimal.Decimal
	TickSize     decimal.Decimal
	Fees         VenueFees

	// Performance metrics
	AvgFillRate float64 // Historical fill rate
	AvgLatency  time.Duration
	RejectRate  float64

	// Dark pool specific
	IsDarkPool      bool
	HiddenLiquidity decimal.Decimal

	mu sync.RWMutex
}

// VenueType represents different types of trading venues
type VenueType string

const (
	VenueTypeExchange VenueType = "EXCHANGE"
	VenueTypeDarkPool VenueType = "DARK_POOL"
	VenueTypeECN      VenueType = "ECN"
	VenueTypeCrossing VenueType = "CROSSING_NETWORK"
	VenueTypeRetail   VenueType = "RETAIL_BROKER"
)

// VenueFees represents fee structure for a venue
type VenueFees struct {
	MakerFee decimal.Decimal `json:"maker_fee"`
	TakerFee decimal.Decimal `json:"taker_fee"`
	MinFee   decimal.Decimal `json:"min_fee"`
	MaxFee   decimal.Decimal `json:"max_fee"`
}

// VenueMarketData represents real-time market data from a venue
type VenueMarketData struct {
	VenueID   string
	Symbol    string
	BidPrice  decimal.Decimal
	AskPrice  decimal.Decimal
	BidSize   decimal.Decimal
	AskSize   decimal.Decimal
	LastPrice decimal.Decimal
	Volume    decimal.Decimal
	Timestamp time.Time
	Sequence  uint64
}

// BestPriceAggregation represents aggregated best prices across venues
type BestPriceAggregation struct {
	Symbol    string
	BestBid   decimal.Decimal
	BestAsk   decimal.Decimal
	BidVenue  string
	AskVenue  string
	BidSize   decimal.Decimal
	AskSize   decimal.Decimal
	Spread    decimal.Decimal
	UpdatedAt time.Time

	// Liquidity depth
	BidDepth []SORPriceLevel
	AskDepth []SORPriceLevel
}

// SORPriceLevel represents a price level with aggregated liquidity for SOR
type SORPriceLevel struct {
	Price  decimal.Decimal
	Size   decimal.Decimal
	Venues []string
}

// RoutingDecision represents the result of smart order routing
type RoutingDecision struct {
	OrderID       uuid.UUID
	Symbol        string
	TotalQuantity decimal.Decimal

	// Routing strategy
	Algorithm    string
	Venues       []VenueAllocation
	ExpectedCost decimal.Decimal
	ExpectedTime time.Duration

	// Risk assessment
	MarketImpact float64
	RiskScore    float64
	Approved     bool

	// Execution plan
	ChildOrders []ChildOrder
	Timestamp   time.Time
}

// VenueAllocation represents allocation to a specific venue
type VenueAllocation struct {
	VenueID       string
	Quantity      decimal.Decimal
	ExpectedPrice decimal.Decimal
	Priority      int
	OrderType     OrderType
	TimeInForce   TimeInForce
}

// ChildOrder represents a child order sent to a venue
type ChildOrder struct {
	ID          uuid.UUID
	ParentID    uuid.UUID
	VenueID     string
	Symbol      string
	Side        OrderSide
	Quantity    decimal.Decimal
	Price       decimal.Decimal
	OrderType   OrderType
	TimeInForce TimeInForce
	Status      OrderStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewSmartOrderRouter creates a new smart order routing system
func NewSmartOrderRouter(logger *observability.Logger, config SORConfig) *SmartOrderRouter {
	// Set default values
	if config.MaxVenuesPerOrder == 0 {
		config.MaxVenuesPerOrder = 5
	}
	if config.RoutingAlgorithm == "" {
		config.RoutingAlgorithm = "TWAP"
	}
	if config.SlippageTarget == 0 {
		config.SlippageTarget = 10.0 // 10 basis points
	}
	if config.TimeHorizon == 0 {
		config.TimeHorizon = 30 * time.Second
	}
	if config.LatencyTarget == 0 {
		config.LatencyTarget = 5 * time.Millisecond
	}

	sor := &SmartOrderRouter{
		logger:          logger,
		config:          config,
		venues:          make(map[string]*TradingVenue),
		marketDataFeeds: make(map[string]chan *VenueMarketData),
		bestPrices:      make(map[string]*BestPriceAggregation),
		stopChan:        make(chan struct{}),
	}

	// Initialize components
	sor.venueSelector = NewVenueSelector(logger, config)
	sor.liquidityAggr = NewLiquidityAggregator(logger, config)
	sor.orderSplitter = NewOrderSplitter(logger, config)
	sor.executionEngine = NewExecutionEngine(logger, config)
	sor.preTradeRisk = NewPreTradeRiskEngine(logger, config)
	sor.complianceCheck = NewComplianceEngine(logger, config)

	return sor
}

// Start begins the smart order routing system
func (sor *SmartOrderRouter) Start(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&sor.isRunning, 0, 1) {
		return fmt.Errorf("smart order router is already running")
	}

	sor.logger.Info(ctx, "Starting smart order routing system", map[string]interface{}{
		"enabled_venues":    len(sor.config.EnabledVenues),
		"routing_algorithm": sor.config.RoutingAlgorithm,
		"slippage_target":   sor.config.SlippageTarget,
		"enable_dark_pools": sor.config.EnableDarkPools,
	})

	// Initialize trading venues
	if err := sor.initializeVenues(ctx); err != nil {
		return fmt.Errorf("failed to initialize venues: %w", err)
	}

	// Start market data aggregation
	if err := sor.startMarketDataAggregation(ctx); err != nil {
		return fmt.Errorf("failed to start market data aggregation: %w", err)
	}

	// Start processing threads
	sor.wg.Add(4)
	go sor.processMarketData(ctx)
	go sor.monitorVenues(ctx)
	go sor.updateBestPrices(ctx)
	go sor.performanceMonitor(ctx)

	sor.logger.Info(ctx, "Smart order routing system started successfully", nil)
	return nil
}

// Stop gracefully shuts down the smart order routing system
func (sor *SmartOrderRouter) Stop(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&sor.isRunning, 1, 0) {
		return fmt.Errorf("smart order router is not running")
	}

	sor.logger.Info(ctx, "Stopping smart order routing system", nil)

	close(sor.stopChan)
	sor.wg.Wait()

	sor.logger.Info(ctx, "Smart order routing system stopped", map[string]interface{}{
		"orders_routed":      atomic.LoadInt64(&sor.ordersRouted),
		"venues_used":        atomic.LoadInt64(&sor.venuesUsed),
		"avg_execution_time": atomic.LoadInt64(&sor.avgExecutionTime),
		"slippage_reduction": sor.slippageReduction,
	})

	return nil
}

// RouteOrder performs intelligent order routing with liquidity aggregation
func (sor *SmartOrderRouter) RouteOrder(ctx context.Context, order *OrderRequest) (*RoutingDecision, error) {
	if atomic.LoadInt32(&sor.isRunning) != 1 {
		return nil, fmt.Errorf("smart order router is not running")
	}

	start := time.Now()

	sor.logger.Info(ctx, "Routing order", map[string]interface{}{
		"symbol":   order.Symbol,
		"side":     string(order.Side),
		"quantity": order.Quantity.String(),
		"type":     string(order.Type),
	})

	// Step 1: Pre-trade risk checks
	if sor.config.EnablePreTrade {
		if err := sor.preTradeRisk.ValidateOrder(ctx, order); err != nil {
			return nil, fmt.Errorf("pre-trade risk check failed: %w", err)
		}
	}

	// Step 2: Compliance checks
	if err := sor.complianceCheck.ValidateOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("compliance check failed: %w", err)
	}

	// Step 3: Get current market data and liquidity
	marketData := sor.getMarketData(order.Symbol)
	if marketData == nil {
		return nil, fmt.Errorf("no market data available for symbol %s", order.Symbol)
	}

	// Step 4: Select optimal venues
	venues, err := sor.venueSelector.SelectVenues(ctx, order, marketData)
	if err != nil {
		return nil, fmt.Errorf("venue selection failed: %w", err)
	}

	// Step 5: Calculate optimal allocation
	allocations, err := sor.calculateOptimalAllocation(ctx, order, venues, marketData)
	if err != nil {
		return nil, fmt.Errorf("allocation calculation failed: %w", err)
	}

	// Step 6: Create child orders
	childOrders, err := sor.orderSplitter.SplitOrder(ctx, order, allocations)
	if err != nil {
		return nil, fmt.Errorf("order splitting failed: %w", err)
	}

	// Step 7: Calculate expected costs and market impact
	expectedCost, marketImpact := sor.calculateExpectedCost(allocations, marketData)

	// Step 8: Create routing decision
	decision := &RoutingDecision{
		OrderID:       uuid.New(),
		Symbol:        order.Symbol,
		TotalQuantity: order.Quantity,
		Algorithm:     sor.config.RoutingAlgorithm,
		Venues:        allocations,
		ExpectedCost:  expectedCost,
		ExpectedTime:  sor.config.TimeHorizon,
		MarketImpact:  marketImpact,
		RiskScore:     sor.calculateRiskScore(order, allocations),
		Approved:      marketImpact <= sor.config.MaxMarketImpact,
		ChildOrders:   childOrders,
		Timestamp:     time.Now(),
	}

	// Update performance metrics
	executionTime := time.Since(start).Nanoseconds()
	atomic.StoreInt64(&sor.avgExecutionTime, executionTime)
	atomic.AddInt64(&sor.ordersRouted, 1)
	atomic.AddInt64(&sor.venuesUsed, int64(len(allocations)))

	sor.logger.Info(ctx, "Order routing completed", map[string]interface{}{
		"order_id":       decision.OrderID.String(),
		"venues_used":    len(allocations),
		"expected_cost":  expectedCost.String(),
		"market_impact":  marketImpact,
		"execution_time": executionTime,
		"approved":       decision.Approved,
	})

	return decision, nil
}

// calculateOptimalAllocation calculates the optimal allocation across venues
func (sor *SmartOrderRouter) calculateOptimalAllocation(ctx context.Context, order *OrderRequest, venues []*TradingVenue, marketData *BestPriceAggregation) ([]VenueAllocation, error) {
	switch sor.config.RoutingAlgorithm {
	case "TWAP":
		return sor.calculateTWAPAllocation(order, venues, marketData)
	case "VWAP":
		return sor.calculateVWAPAllocation(order, venues, marketData)
	case "IMPLEMENTATION_SHORTFALL":
		return sor.calculateImplementationShortfallAllocation(order, venues, marketData)
	case "LIQUIDITY_SEEKING":
		return sor.calculateLiquiditySeekingAllocation(order, venues, marketData)
	default:
		return sor.calculateTWAPAllocation(order, venues, marketData)
	}
}

// calculateTWAPAllocation implements Time-Weighted Average Price allocation
func (sor *SmartOrderRouter) calculateTWAPAllocation(order *OrderRequest, venues []*TradingVenue, marketData *BestPriceAggregation) ([]VenueAllocation, error) {
	var allocations []VenueAllocation
	remainingQty := order.Quantity

	// Sort venues by best price for the order side
	sortedVenues := sor.sortVenuesByPrice(venues, order.Side)

	// Allocate quantity across venues based on available liquidity
	for i, venue := range sortedVenues {
		if remainingQty.IsZero() {
			break
		}

		// Calculate available liquidity at this venue
		var availableLiquidity decimal.Decimal
		var expectedPrice decimal.Decimal

		if order.Side == OrderSideBuy {
			availableLiquidity = venue.AskSize
			expectedPrice = venue.BestAsk
		} else {
			availableLiquidity = venue.BidSize
			expectedPrice = venue.BestBid
		}

		// Don't allocate if no liquidity or price is stale
		if availableLiquidity.IsZero() || time.Since(venue.LastUpdate) > sor.config.CacheTimeout {
			continue
		}

		// Calculate allocation quantity
		allocQty := decimal.Min(remainingQty, availableLiquidity)

		// Apply minimum venue size constraint
		if allocQty.LessThan(sor.config.MinVenueSize) && i < len(sortedVenues)-1 {
			continue
		}

		allocation := VenueAllocation{
			VenueID:       venue.ID,
			Quantity:      allocQty,
			ExpectedPrice: expectedPrice,
			Priority:      i + 1,
			OrderType:     order.Type,
			TimeInForce:   order.TimeInForce,
		}

		allocations = append(allocations, allocation)
		remainingQty = remainingQty.Sub(allocQty)

		// Limit number of venues per order
		if len(allocations) >= sor.config.MaxVenuesPerOrder {
			break
		}
	}

	if len(allocations) == 0 {
		return nil, fmt.Errorf("no suitable venues found for order")
	}

	return allocations, nil
}

// calculateVWAPAllocation implements Volume-Weighted Average Price allocation
func (sor *SmartOrderRouter) calculateVWAPAllocation(order *OrderRequest, venues []*TradingVenue, marketData *BestPriceAggregation) ([]VenueAllocation, error) {
	// Calculate total available volume across venues
	var totalVolume decimal.Decimal
	venueVolumes := make(map[string]decimal.Decimal)

	for _, venue := range venues {
		var volume decimal.Decimal
		if order.Side == OrderSideBuy {
			volume = venue.AskSize
		} else {
			volume = venue.BidSize
		}

		if volume.GreaterThan(decimal.Zero) {
			venueVolumes[venue.ID] = volume
			totalVolume = totalVolume.Add(volume)
		}
	}

	if totalVolume.IsZero() {
		return nil, fmt.Errorf("no liquidity available across venues")
	}

	// Allocate based on volume proportions
	var allocations []VenueAllocation
	remainingQty := order.Quantity

	for _, venue := range venues {
		volume, exists := venueVolumes[venue.ID]
		if !exists || remainingQty.IsZero() {
			continue
		}

		// Calculate proportional allocation
		proportion := volume.Div(totalVolume)
		allocQty := order.Quantity.Mul(proportion)

		// Don't exceed remaining quantity
		allocQty = decimal.Min(allocQty, remainingQty)

		// Apply minimum size constraint
		if allocQty.LessThan(sor.config.MinVenueSize) {
			continue
		}

		var expectedPrice decimal.Decimal
		if order.Side == OrderSideBuy {
			expectedPrice = venue.BestAsk
		} else {
			expectedPrice = venue.BestBid
		}

		allocation := VenueAllocation{
			VenueID:       venue.ID,
			Quantity:      allocQty,
			ExpectedPrice: expectedPrice,
			Priority:      1, // VWAP allocations have equal priority
			OrderType:     order.Type,
			TimeInForce:   order.TimeInForce,
		}

		allocations = append(allocations, allocation)
		remainingQty = remainingQty.Sub(allocQty)
	}

	return allocations, nil
}

// calculateImplementationShortfallAllocation implements Implementation Shortfall allocation
func (sor *SmartOrderRouter) calculateImplementationShortfallAllocation(order *OrderRequest, venues []*TradingVenue, marketData *BestPriceAggregation) ([]VenueAllocation, error) {
	// Implementation Shortfall minimizes the cost of trading including market impact
	// This is a simplified version - production would use more sophisticated models

	var allocations []VenueAllocation
	remainingQty := order.Quantity

	// Sort venues by execution quality (combination of price and reliability)
	sortedVenues := sor.sortVenuesByExecutionQuality(venues, order.Side)

	for i, venue := range sortedVenues {
		if remainingQty.IsZero() {
			break
		}

		// Calculate market impact for this venue
		marketImpact := sor.calculateVenueMarketImpact(venue, order)

		// Skip venues with high market impact
		if marketImpact > sor.config.MaxMarketImpact {
			continue
		}

		var availableLiquidity decimal.Decimal
		var expectedPrice decimal.Decimal

		if order.Side == OrderSideBuy {
			availableLiquidity = venue.AskSize
			expectedPrice = venue.BestAsk
		} else {
			availableLiquidity = venue.BidSize
			expectedPrice = venue.BestBid
		}

		// Calculate optimal allocation considering market impact
		allocQty := sor.calculateOptimalQuantity(remainingQty, availableLiquidity, marketImpact)

		if allocQty.LessThan(sor.config.MinVenueSize) {
			continue
		}

		allocation := VenueAllocation{
			VenueID:       venue.ID,
			Quantity:      allocQty,
			ExpectedPrice: expectedPrice,
			Priority:      i + 1,
			OrderType:     order.Type,
			TimeInForce:   order.TimeInForce,
		}

		allocations = append(allocations, allocation)
		remainingQty = remainingQty.Sub(allocQty)

		if len(allocations) >= sor.config.MaxVenuesPerOrder {
			break
		}
	}

	return allocations, nil
}

// calculateLiquiditySeekingAllocation implements Liquidity Seeking allocation
func (sor *SmartOrderRouter) calculateLiquiditySeekingAllocation(order *OrderRequest, venues []*TradingVenue, marketData *BestPriceAggregation) ([]VenueAllocation, error) {
	var allocations []VenueAllocation
	remainingQty := order.Quantity

	// Include dark pools if enabled
	if sor.config.EnableDarkPools {
		darkPoolAllocation := order.Quantity.Mul(decimal.NewFromFloat(sor.config.DarkPoolParticipation))

		for _, venue := range venues {
			if venue.IsDarkPool && remainingQty.GreaterThan(decimal.Zero) {
				allocQty := decimal.Min(darkPoolAllocation, remainingQty)
				allocQty = decimal.Min(allocQty, venue.HiddenLiquidity)

				if allocQty.GreaterThanOrEqual(sor.config.MinVenueSize) {
					allocation := VenueAllocation{
						VenueID:       venue.ID,
						Quantity:      allocQty,
						ExpectedPrice: marketData.BestBid.Add(marketData.BestAsk).Div(decimal.NewFromInt(2)), // Mid price
						Priority:      1,                                                                     // High priority for dark pools
						OrderType:     OrderTypeLimit,
						TimeInForce:   TimeInForceIOC,
					}

					allocations = append(allocations, allocation)
					remainingQty = remainingQty.Sub(allocQty)
				}
			}
		}
	}

	// Allocate remaining quantity to lit venues
	if remainingQty.GreaterThan(decimal.Zero) {
		litVenues := sor.filterLitVenues(venues)
		twapAllocations, err := sor.calculateTWAPAllocation(&OrderRequest{
			Symbol:      order.Symbol,
			Side:        order.Side,
			Quantity:    remainingQty,
			Type:        order.Type,
			TimeInForce: order.TimeInForce,
		}, litVenues, marketData)

		if err != nil {
			return allocations, err // Return dark pool allocations even if lit allocation fails
		}

		allocations = append(allocations, twapAllocations...)
	}

	return allocations, nil
}

// Helper methods for Smart Order Routing

// sortVenuesByPrice sorts venues by best price for the given order side
func (sor *SmartOrderRouter) sortVenuesByPrice(venues []*TradingVenue, side OrderSide) []*TradingVenue {
	sortedVenues := make([]*TradingVenue, len(venues))
	copy(sortedVenues, venues)

	sort.Slice(sortedVenues, func(i, j int) bool {
		if side == OrderSideBuy {
			// For buy orders, prefer lower ask prices
			return sortedVenues[i].BestAsk.LessThan(sortedVenues[j].BestAsk)
		} else {
			// For sell orders, prefer higher bid prices
			return sortedVenues[i].BestBid.GreaterThan(sortedVenues[j].BestBid)
		}
	})

	return sortedVenues
}

// sortVenuesByExecutionQuality sorts venues by execution quality
func (sor *SmartOrderRouter) sortVenuesByExecutionQuality(venues []*TradingVenue, side OrderSide) []*TradingVenue {
	sortedVenues := make([]*TradingVenue, len(venues))
	copy(sortedVenues, venues)

	sort.Slice(sortedVenues, func(i, j int) bool {
		scoreI := sor.calculateExecutionQualityScore(sortedVenues[i], side)
		scoreJ := sor.calculateExecutionQualityScore(sortedVenues[j], side)
		return scoreI > scoreJ // Higher score is better
	})

	return sortedVenues
}

// calculateExecutionQualityScore calculates a quality score for a venue
func (sor *SmartOrderRouter) calculateExecutionQualityScore(venue *TradingVenue, side OrderSide) float64 {
	// Combine multiple factors: price, reliability, fill rate, latency
	priceScore := 1.0
	if side == OrderSideBuy && !venue.BestAsk.IsZero() {
		priceScore = 1.0 / venue.BestAsk.InexactFloat64() * 100000 // Normalize
	} else if side == OrderSideSell && !venue.BestBid.IsZero() {
		priceScore = venue.BestBid.InexactFloat64() / 100000 // Normalize
	}

	reliabilityScore := venue.Reliability * 100
	fillRateScore := venue.AvgFillRate * 100
	latencyScore := 100.0 / (float64(venue.AvgLatency.Nanoseconds()) / 1000000.0) // Lower latency is better

	// Weighted combination
	return priceScore*0.4 + reliabilityScore*0.3 + fillRateScore*0.2 + latencyScore*0.1
}

// calculateVenueMarketImpact calculates market impact for a venue
func (sor *SmartOrderRouter) calculateVenueMarketImpact(venue *TradingVenue, order *OrderRequest) float64 {
	var availableLiquidity decimal.Decimal

	if order.Side == OrderSideBuy {
		availableLiquidity = venue.AskSize
	} else {
		availableLiquidity = venue.BidSize
	}

	if availableLiquidity.IsZero() {
		return 100.0 // Maximum impact if no liquidity
	}

	// Simple market impact model: impact increases with order size relative to available liquidity
	ratio := order.Quantity.Div(availableLiquidity).InexactFloat64()

	// Square root model for market impact
	impact := ratio * ratio * 100.0 // Convert to basis points

	if impact > 100.0 {
		impact = 100.0 // Cap at 100%
	}

	return impact
}

// calculateOptimalQuantity calculates optimal quantity considering market impact
func (sor *SmartOrderRouter) calculateOptimalQuantity(remainingQty, availableLiquidity decimal.Decimal, marketImpact float64) decimal.Decimal {
	// Don't take more than 50% of available liquidity to minimize impact
	maxQty := availableLiquidity.Mul(decimal.NewFromFloat(0.5))

	// Reduce quantity if market impact is high
	if marketImpact > 50.0 {
		impactReduction := decimal.NewFromFloat(1.0 - (marketImpact-50.0)/100.0)
		maxQty = maxQty.Mul(impactReduction)
	}

	return decimal.Min(remainingQty, maxQty)
}

// filterLitVenues filters out dark pools to get only lit venues
func (sor *SmartOrderRouter) filterLitVenues(venues []*TradingVenue) []*TradingVenue {
	var litVenues []*TradingVenue

	for _, venue := range venues {
		if !venue.IsDarkPool {
			litVenues = append(litVenues, venue)
		}
	}

	return litVenues
}

// calculateExpectedCost calculates expected execution cost and market impact
func (sor *SmartOrderRouter) calculateExpectedCost(allocations []VenueAllocation, marketData *BestPriceAggregation) (decimal.Decimal, float64) {
	var totalCost decimal.Decimal
	var totalQuantity decimal.Decimal
	var weightedMarketImpact float64

	for _, allocation := range allocations {
		cost := allocation.Quantity.Mul(allocation.ExpectedPrice)
		totalCost = totalCost.Add(cost)
		totalQuantity = totalQuantity.Add(allocation.Quantity)

		// Calculate market impact for this allocation (simplified)
		impact := allocation.Quantity.Div(totalQuantity).InexactFloat64() * 10.0 // Basis points
		weightedMarketImpact += impact
	}

	return totalCost, weightedMarketImpact
}

// calculateRiskScore calculates overall risk score for the routing decision
func (sor *SmartOrderRouter) calculateRiskScore(order *OrderRequest, allocations []VenueAllocation) float64 {
	// Risk factors: number of venues, total quantity, market impact
	venueRisk := float64(len(allocations)) * 5.0            // More venues = more risk
	quantityRisk := order.Quantity.InexactFloat64() / 100.0 // Larger orders = more risk

	// Venue concentration risk
	var maxAllocation decimal.Decimal
	for _, allocation := range allocations {
		if allocation.Quantity.GreaterThan(maxAllocation) {
			maxAllocation = allocation.Quantity
		}
	}

	concentrationRisk := maxAllocation.Div(order.Quantity).InexactFloat64() * 50.0

	totalRisk := venueRisk + quantityRisk + concentrationRisk

	// Cap at 100
	if totalRisk > 100.0 {
		totalRisk = 100.0
	}

	return totalRisk
}

// getMarketData retrieves current market data for a symbol
func (sor *SmartOrderRouter) getMarketData(symbol string) *BestPriceAggregation {
	sor.mu.RLock()
	defer sor.mu.RUnlock()

	return sor.bestPrices[symbol]
}

// initializeVenues initializes trading venues
func (sor *SmartOrderRouter) initializeVenues(ctx context.Context) error {
	sor.logger.Info(ctx, "Initializing trading venues", map[string]interface{}{
		"enabled_venues": sor.config.EnabledVenues,
	})

	// Initialize venues based on configuration
	for _, venueID := range sor.config.EnabledVenues {
		venue := &TradingVenue{
			ID:           venueID,
			Name:         venueID,
			Type:         VenueTypeExchange,
			Exchange:     venueID,
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
		}

		sor.venues[venueID] = venue
	}

	sor.logger.Info(ctx, "Trading venues initialized", map[string]interface{}{
		"venue_count": len(sor.venues),
	})

	return nil
}

// startMarketDataAggregation starts market data aggregation
func (sor *SmartOrderRouter) startMarketDataAggregation(ctx context.Context) error {
	sor.logger.Info(ctx, "Starting market data aggregation", nil)

	// Initialize market data feeds for each venue
	for venueID := range sor.venues {
		feed := make(chan *VenueMarketData, 1000)
		sor.marketDataFeeds[venueID] = feed
	}

	// Initialize best prices aggregation
	sor.bestPrices["BTCUSDT"] = &BestPriceAggregation{
		Symbol:    "BTCUSDT",
		BestBid:   decimal.NewFromFloat(45000.0),
		BestAsk:   decimal.NewFromFloat(45001.0),
		BidVenue:  "binance",
		AskVenue:  "binance",
		BidSize:   decimal.NewFromFloat(10.0),
		AskSize:   decimal.NewFromFloat(8.0),
		Spread:    decimal.NewFromFloat(1.0),
		UpdatedAt: time.Now(),
	}

	return nil
}

// processMarketData processes incoming market data
func (sor *SmartOrderRouter) processMarketData(ctx context.Context) {
	defer sor.wg.Done()

	sor.logger.Info(ctx, "Starting market data processor", nil)

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-sor.stopChan:
			return
		case <-ticker.C:
			// Process market data from all venues
			for venueID, feed := range sor.marketDataFeeds {
				select {
				case data := <-feed:
					sor.processVenueMarketData(ctx, venueID, data)
				default:
					// No data available
				}
			}
		}
	}
}

// processVenueMarketData processes market data from a specific venue
func (sor *SmartOrderRouter) processVenueMarketData(ctx context.Context, venueID string, data *VenueMarketData) {
	sor.mu.Lock()
	defer sor.mu.Unlock()

	// Update venue data
	if venue, exists := sor.venues[venueID]; exists {
		venue.BestBid = data.BidPrice
		venue.BestAsk = data.AskPrice
		venue.BidSize = data.BidSize
		venue.AskSize = data.AskSize
		venue.LastUpdate = data.Timestamp
	}
}

// monitorVenues monitors venue health and connectivity
func (sor *SmartOrderRouter) monitorVenues(ctx context.Context) {
	defer sor.wg.Done()

	sor.logger.Info(ctx, "Starting venue monitor", nil)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-sor.stopChan:
			return
		case <-ticker.C:
			sor.checkVenueHealth(ctx)
		}
	}
}

// checkVenueHealth checks the health of all venues
func (sor *SmartOrderRouter) checkVenueHealth(ctx context.Context) {
	sor.mu.RLock()
	defer sor.mu.RUnlock()

	for venueID, venue := range sor.venues {
		// Check if venue data is stale
		if time.Since(venue.LastUpdate) > 30*time.Second {
			venue.IsConnected = false
			sor.logger.Warn(ctx, "Venue appears disconnected", map[string]interface{}{
				"venue_id":    venueID,
				"last_update": venue.LastUpdate,
			})
		} else {
			venue.IsConnected = true
		}
	}
}

// updateBestPrices updates aggregated best prices
func (sor *SmartOrderRouter) updateBestPrices(ctx context.Context) {
	defer sor.wg.Done()

	sor.logger.Info(ctx, "Starting best price updater", nil)

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-sor.stopChan:
			return
		case <-ticker.C:
			sor.aggregateBestPrices(ctx)
		}
	}
}

// aggregateBestPrices aggregates best prices across all venues
func (sor *SmartOrderRouter) aggregateBestPrices(ctx context.Context) {
	sor.mu.Lock()
	defer sor.mu.Unlock()

	// For each symbol, find best bid and ask across venues
	symbols := []string{"BTCUSDT"} // In production, this would be dynamic

	for _, symbol := range symbols {
		var bestBid, bestAsk decimal.Decimal
		var bidVenue, askVenue string
		var bidSize, askSize decimal.Decimal

		for venueID, venue := range sor.venues {
			if !venue.IsConnected {
				continue
			}

			// Check best bid
			if bestBid.IsZero() || venue.BestBid.GreaterThan(bestBid) {
				bestBid = venue.BestBid
				bidVenue = venueID
				bidSize = venue.BidSize
			}

			// Check best ask
			if bestAsk.IsZero() || venue.BestAsk.LessThan(bestAsk) {
				bestAsk = venue.BestAsk
				askVenue = venueID
				askSize = venue.AskSize
			}
		}

		// Update aggregated prices
		if !bestBid.IsZero() && !bestAsk.IsZero() {
			sor.bestPrices[symbol] = &BestPriceAggregation{
				Symbol:    symbol,
				BestBid:   bestBid,
				BestAsk:   bestAsk,
				BidVenue:  bidVenue,
				AskVenue:  askVenue,
				BidSize:   bidSize,
				AskSize:   askSize,
				Spread:    bestAsk.Sub(bestBid),
				UpdatedAt: time.Now(),
			}
		}
	}
}

// performanceMonitor tracks and reports SOR performance metrics
func (sor *SmartOrderRouter) performanceMonitor(ctx context.Context) {
	defer sor.wg.Done()

	sor.logger.Info(ctx, "Starting SOR performance monitor", nil)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var lastOrderCount int64

	for {
		select {
		case <-sor.stopChan:
			return
		case <-ticker.C:
			currentOrders := atomic.LoadInt64(&sor.ordersRouted)
			ordersPerSecond := currentOrders - lastOrderCount
			lastOrderCount = currentOrders

			avgExecutionTime := atomic.LoadInt64(&sor.avgExecutionTime)
			venuesUsed := atomic.LoadInt64(&sor.venuesUsed)

			sor.logger.Info(ctx, "SOR performance metrics", map[string]interface{}{
				"orders_per_second":  ordersPerSecond,
				"avg_execution_time": avgExecutionTime,
				"avg_execution_ms":   avgExecutionTime / 1000000,
				"total_orders":       currentOrders,
				"total_venues_used":  venuesUsed,
				"slippage_reduction": sor.slippageReduction,
				"active_venues":      len(sor.venues),
			})
		}
	}
}

// GetMetrics returns current SOR performance metrics
func (sor *SmartOrderRouter) GetMetrics() SORMetrics {
	sor.mu.RLock()
	defer sor.mu.RUnlock()

	connectedVenues := 0
	for _, venue := range sor.venues {
		if venue.IsConnected {
			connectedVenues++
		}
	}

	return SORMetrics{
		OrdersRouted:      atomic.LoadInt64(&sor.ordersRouted),
		VenuesUsed:        atomic.LoadInt64(&sor.venuesUsed),
		AvgExecutionTime:  atomic.LoadInt64(&sor.avgExecutionTime),
		SlippageReduction: sor.slippageReduction,
		IsRunning:         atomic.LoadInt32(&sor.isRunning) == 1,
		ActiveVenues:      len(sor.venues),
		ConnectedVenues:   connectedVenues,
		RoutingAlgorithm:  sor.config.RoutingAlgorithm,
		DarkPoolsEnabled:  sor.config.EnableDarkPools,
	}
}

// SORMetrics contains performance metrics for the Smart Order Router
type SORMetrics struct {
	OrdersRouted      int64   `json:"orders_routed"`
	VenuesUsed        int64   `json:"venues_used"`
	AvgExecutionTime  int64   `json:"avg_execution_time_nanos"`
	SlippageReduction float64 `json:"slippage_reduction"`
	IsRunning         bool    `json:"is_running"`
	ActiveVenues      int     `json:"active_venues"`
	ConnectedVenues   int     `json:"connected_venues"`
	RoutingAlgorithm  string  `json:"routing_algorithm"`
	DarkPoolsEnabled  bool    `json:"dark_pools_enabled"`
}

// ExecuteOrder executes a routing decision by sending child orders to venues
func (sor *SmartOrderRouter) ExecuteOrder(ctx context.Context, decision *RoutingDecision) error {
	if !decision.Approved {
		return fmt.Errorf("routing decision not approved due to risk constraints")
	}

	sor.logger.Info(ctx, "Executing routing decision", map[string]interface{}{
		"order_id":     decision.OrderID.String(),
		"child_orders": len(decision.ChildOrders),
		"venues":       len(decision.Venues),
	})

	// Execute child orders in parallel
	for _, childOrder := range decision.ChildOrders {
		go sor.executeChildOrder(ctx, &childOrder)
	}

	return nil
}

// executeChildOrder executes a single child order at a venue
func (sor *SmartOrderRouter) executeChildOrder(ctx context.Context, childOrder *ChildOrder) {
	sor.logger.Info(ctx, "Executing child order", map[string]interface{}{
		"child_id": childOrder.ID.String(),
		"venue_id": childOrder.VenueID,
		"symbol":   childOrder.Symbol,
		"side":     string(childOrder.Side),
		"quantity": childOrder.Quantity.String(),
		"price":    childOrder.Price.String(),
	})

	// In production, this would send the order to the actual venue
	// For now, we'll simulate execution
	time.Sleep(time.Millisecond) // Simulate network latency

	// Update order status
	childOrder.Status = OrderStatusFilled
	childOrder.UpdatedAt = time.Now()

	sor.logger.Info(ctx, "Child order executed", map[string]interface{}{
		"child_id": childOrder.ID.String(),
		"status":   string(childOrder.Status),
	})
}

// GetBestPrices returns current best prices for a symbol
func (sor *SmartOrderRouter) GetBestPrices(symbol string) *BestPriceAggregation {
	sor.mu.RLock()
	defer sor.mu.RUnlock()

	return sor.bestPrices[symbol]
}

// GetVenues returns all configured venues
func (sor *SmartOrderRouter) GetVenues() map[string]*TradingVenue {
	sor.mu.RLock()
	defer sor.mu.RUnlock()

	venues := make(map[string]*TradingVenue)
	for id, venue := range sor.venues {
		venues[id] = venue
	}
	return venues
}
