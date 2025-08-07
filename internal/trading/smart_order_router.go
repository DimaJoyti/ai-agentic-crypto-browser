package trading

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/shopspring/decimal"
)

// SmartOrderRouter provides intelligent order routing capabilities
type SmartOrderRouter struct {
	logger        *observability.Logger
	config        *RouterConfig
	venues        map[string]*VenueInfo
	routingRules  []*RoutingRule
	venueSelector *VenueSelector
	metrics       *RouterMetrics
	mu            sync.RWMutex
	isRunning     bool
	stopChan      chan struct{}
}

// RouterConfig contains smart order router configuration
type RouterConfig struct {
	EnableSmartRouting     bool               `json:"enable_smart_routing"`
	EnableVenueAggregation bool               `json:"enable_venue_aggregation"`
	MaxVenuesPerOrder      int                `json:"max_venues_per_order"`
	MinFillSize            decimal.Decimal    `json:"min_fill_size"`
	MaxLatencyMs           int                `json:"max_latency_ms"`
	MaxSlippageBps         int                `json:"max_slippage_bps"`
	PreferredVenues        []string           `json:"preferred_venues"`
	VenueWeights           map[string]float64 `json:"venue_weights"`
	RoutingStrategy        RoutingStrategy    `json:"routing_strategy"`
	RebalanceInterval      time.Duration      `json:"rebalance_interval"`
}

// RoutingStrategy defines routing strategies
type RoutingStrategy string

const (
	RoutingStrategyBestPrice     RoutingStrategy = "best_price"
	RoutingStrategyLowestLatency RoutingStrategy = "lowest_latency"
	RoutingStrategyHighestFill   RoutingStrategy = "highest_fill"
	RoutingStrategyLowestCost    RoutingStrategy = "lowest_cost"
	RoutingStrategyBalanced      RoutingStrategy = "balanced"
	RoutingStrategyLiquidity     RoutingStrategy = "liquidity"
	RoutingStrategyDarkPool      RoutingStrategy = "dark_pool"
)

// VenueInfo contains information about a trading venue
type VenueInfo struct {
	ID                 string                     `json:"id"`
	Name               string                     `json:"name"`
	Type               VenueType                  `json:"type"`
	IsActive           bool                       `json:"is_active"`
	Latency            time.Duration              `json:"latency"`
	FeeRate            decimal.Decimal            `json:"fee_rate"`
	MinOrderSize       decimal.Decimal            `json:"min_order_size"`
	MaxOrderSize       decimal.Decimal            `json:"max_order_size"`
	SupportedSymbols   []string                   `json:"supported_symbols"`
	TradingHours       *TradingHours              `json:"trading_hours"`
	Liquidity          map[string]decimal.Decimal `json:"liquidity"`
	HistoricalFillRate decimal.Decimal            `json:"historical_fill_rate"`
	AverageSlippage    decimal.Decimal            `json:"average_slippage"`
	ReliabilityScore   decimal.Decimal            `json:"reliability_score"`
	LastUpdated        time.Time                  `json:"last_updated"`
	Metadata           map[string]interface{}     `json:"metadata"`
}

// VenueType defines types of trading venues
type VenueType string

const (
	VenueTypeExchange    VenueType = "exchange"
	VenueTypeDarkPool    VenueType = "dark_pool"
	VenueTypeECN         VenueType = "ecn"
	VenueTypeMarketMaker VenueType = "market_maker"
	VenueTypeCrossing    VenueType = "crossing"
)

// TradingHours defines trading hours for a venue
type TradingHours struct {
	OpenTime   time.Time `json:"open_time"`
	CloseTime  time.Time `json:"close_time"`
	TimeZone   string    `json:"time_zone"`
	IsOpen24x7 bool      `json:"is_open_24x7"`
	Holidays   []string  `json:"holidays"`
}

// RoutingRule defines a routing rule
type RoutingRule struct {
	ID         string              `json:"id"`
	Name       string              `json:"name"`
	Priority   int                 `json:"priority"`
	Conditions []*RoutingCondition `json:"conditions"`
	Actions    []*RoutingAction    `json:"actions"`
	IsActive   bool                `json:"is_active"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
}

// RoutingCondition defines a routing condition
type RoutingCondition struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

// RoutingAction defines a routing action
type RoutingAction struct {
	Type       ActionType             `json:"type"`
	VenueID    string                 `json:"venue_id,omitempty"`
	Percentage decimal.Decimal        `json:"percentage,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// ActionType defines routing action types
type ActionType string

const (
	ActionTypeRoute  ActionType = "route"
	ActionTypeSplit  ActionType = "split"
	ActionTypeReject ActionType = "reject"
	ActionTypeDelay  ActionType = "delay"
	ActionTypeModify ActionType = "modify"
)

// VenueSelector selects optimal venues for order execution
type VenueSelector struct {
	strategy RoutingStrategy
	weights  map[string]float64
}

// RouterMetrics tracks routing performance
type RouterMetrics struct {
	TotalOrders      int64                       `json:"total_orders"`
	RoutedOrders     int64                       `json:"routed_orders"`
	SplitOrders      int64                       `json:"split_orders"`
	RejectedOrders   int64                       `json:"rejected_orders"`
	AverageLatency   time.Duration               `json:"average_latency"`
	AverageFillRate  decimal.Decimal             `json:"average_fill_rate"`
	AverageSlippage  decimal.Decimal             `json:"average_slippage"`
	VenuePerformance map[string]*VenueMetrics    `json:"venue_performance"`
	StrategyMetrics  map[string]*StrategyMetrics `json:"strategy_metrics"`
	LastUpdated      time.Time                   `json:"last_updated"`
}

// VenueMetrics tracks venue-specific metrics
type VenueMetrics struct {
	OrderCount       int64           `json:"order_count"`
	FillRate         decimal.Decimal `json:"fill_rate"`
	AverageLatency   time.Duration   `json:"average_latency"`
	AverageSlippage  decimal.Decimal `json:"average_slippage"`
	ReliabilityScore decimal.Decimal `json:"reliability_score"`
	LastUpdated      time.Time       `json:"last_updated"`
}

// StrategyMetrics tracks strategy-specific metrics
type StrategyMetrics struct {
	OrderCount       int64           `json:"order_count"`
	SuccessRate      decimal.Decimal `json:"success_rate"`
	AverageExecution time.Duration   `json:"average_execution"`
	CostSavings      decimal.Decimal `json:"cost_savings"`
	LastUpdated      time.Time       `json:"last_updated"`
}

// RoutingDecision represents a routing decision
type RoutingDecision struct {
	OrderID          string             `json:"order_id"`
	Strategy         RoutingStrategy    `json:"strategy"`
	SelectedVenues   []*VenueAllocation `json:"selected_venues"`
	RejectionReason  string             `json:"rejection_reason,omitempty"`
	EstimatedCost    decimal.Decimal    `json:"estimated_cost"`
	EstimatedLatency time.Duration      `json:"estimated_latency"`
	ConfidenceScore  decimal.Decimal    `json:"confidence_score"`
	DecisionTime     time.Time          `json:"decision_time"`
}

// VenueAllocation represents venue allocation for an order
type VenueAllocation struct {
	VenueID    string          `json:"venue_id"`
	VenueName  string          `json:"venue_name"`
	Quantity   decimal.Decimal `json:"quantity"`
	Percentage decimal.Decimal `json:"percentage"`
	Priority   int             `json:"priority"`
	Reason     string          `json:"reason"`
}

// NewSmartOrderRouter creates a new smart order router
func NewSmartOrderRouter(logger *observability.Logger) *SmartOrderRouter {
	config := &RouterConfig{
		EnableSmartRouting:     true,
		EnableVenueAggregation: true,
		MaxVenuesPerOrder:      3,
		MinFillSize:            decimal.NewFromFloat(0.01),
		MaxLatencyMs:           100,
		MaxSlippageBps:         10,
		PreferredVenues:        []string{"binance", "coinbase", "kraken"},
		VenueWeights: map[string]float64{
			"binance":  0.4,
			"coinbase": 0.3,
			"kraken":   0.3,
		},
		RoutingStrategy:   RoutingStrategyBalanced,
		RebalanceInterval: 5 * time.Minute,
	}

	return &SmartOrderRouter{
		logger:       logger,
		config:       config,
		venues:       make(map[string]*VenueInfo),
		routingRules: make([]*RoutingRule, 0),
		venueSelector: &VenueSelector{
			strategy: config.RoutingStrategy,
			weights:  config.VenueWeights,
		},
		metrics: &RouterMetrics{
			VenuePerformance: make(map[string]*VenueMetrics),
			StrategyMetrics:  make(map[string]*StrategyMetrics),
			LastUpdated:      time.Now(),
		},
		stopChan: make(chan struct{}),
	}
}

// Start starts the smart order router
func (sor *SmartOrderRouter) Start(ctx context.Context) error {
	sor.mu.Lock()
	defer sor.mu.Unlock()

	if sor.isRunning {
		return fmt.Errorf("smart order router is already running")
	}

	sor.isRunning = true

	// Initialize default venues
	sor.initializeDefaultVenues()

	// Initialize default routing rules
	sor.initializeDefaultRules()

	// Start background processes
	go sor.venueMonitoringLoop(ctx)
	go sor.metricsUpdateLoop(ctx)

	sor.logger.Info(ctx, "Smart order router started", map[string]interface{}{
		"routing_strategy":     sor.config.RoutingStrategy,
		"venues_count":         len(sor.venues),
		"routing_rules_count":  len(sor.routingRules),
		"max_venues_per_order": sor.config.MaxVenuesPerOrder,
	})

	return nil
}

// Stop stops the smart order router
func (sor *SmartOrderRouter) Stop(ctx context.Context) error {
	sor.mu.Lock()
	defer sor.mu.Unlock()

	if !sor.isRunning {
		return nil
	}

	sor.isRunning = false
	close(sor.stopChan)

	sor.logger.Info(ctx, "Smart order router stopped", nil)
	return nil
}

// RouteOrder routes an order to optimal venues
func (sor *SmartOrderRouter) RouteOrder(ctx context.Context, order *ExecutionOrder) (*RoutingDecision, error) {
	sor.mu.RLock()
	defer sor.mu.RUnlock()

	if !sor.config.EnableSmartRouting {
		// Default routing to first available venue
		return sor.defaultRouting(order)
	}

	// Apply routing rules
	decision, err := sor.applyRoutingRules(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to apply routing rules: %w", err)
	}

	// Select optimal venues
	venues, err := sor.selectOptimalVenues(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to select optimal venues: %w", err)
	}

	decision.SelectedVenues = venues
	decision.DecisionTime = time.Now()

	// Calculate estimates
	sor.calculateEstimates(decision, order)

	// Update metrics
	sor.updateRoutingMetrics(decision)

	sor.logger.Info(ctx, "Order routed", map[string]interface{}{
		"order_id":          order.ID,
		"strategy":          decision.Strategy,
		"venues_count":      len(decision.SelectedVenues),
		"estimated_cost":    decision.EstimatedCost.String(),
		"estimated_latency": decision.EstimatedLatency,
		"confidence_score":  decision.ConfidenceScore.String(),
	})

	return decision, nil
}

// RegisterVenue registers a new trading venue
func (sor *SmartOrderRouter) RegisterVenue(venue *VenueInfo) error {
	sor.mu.Lock()
	defer sor.mu.Unlock()

	venue.LastUpdated = time.Now()
	sor.venues[venue.ID] = venue

	// Initialize venue metrics
	sor.metrics.VenuePerformance[venue.ID] = &VenueMetrics{
		LastUpdated: time.Now(),
	}

	sor.logger.Info(context.Background(), "Venue registered", map[string]interface{}{
		"venue_id":   venue.ID,
		"venue_name": venue.Name,
		"venue_type": venue.Type,
		"fee_rate":   venue.FeeRate.String(),
		"latency":    venue.Latency,
	})

	return nil
}

// GetVenueInfo retrieves venue information
func (sor *SmartOrderRouter) GetVenueInfo(venueID string) (*VenueInfo, error) {
	sor.mu.RLock()
	defer sor.mu.RUnlock()

	venue, exists := sor.venues[venueID]
	if !exists {
		return nil, fmt.Errorf("venue not found: %s", venueID)
	}

	return venue, nil
}

// GetMetrics returns routing metrics
func (sor *SmartOrderRouter) GetMetrics() *RouterMetrics {
	sor.mu.RLock()
	defer sor.mu.RUnlock()

	return sor.metrics
}

// defaultRouting provides default routing logic
func (sor *SmartOrderRouter) defaultRouting(order *ExecutionOrder) (*RoutingDecision, error) {
	// Find first available venue
	for venueID, venue := range sor.venues {
		if venue.IsActive && sor.supportsSymbol(venue, order.Symbol) {
			allocation := &VenueAllocation{
				VenueID:    venueID,
				VenueName:  venue.Name,
				Quantity:   order.Quantity,
				Percentage: decimal.NewFromFloat(1.0),
				Priority:   1,
				Reason:     "default_routing",
			}

			return &RoutingDecision{
				OrderID:        order.ID,
				Strategy:       RoutingStrategyBestPrice,
				SelectedVenues: []*VenueAllocation{allocation},
				DecisionTime:   time.Now(),
			}, nil
		}
	}

	return nil, fmt.Errorf("no available venues for symbol: %s", order.Symbol)
}

// applyRoutingRules applies routing rules to an order
func (sor *SmartOrderRouter) applyRoutingRules(ctx context.Context, order *ExecutionOrder) (*RoutingDecision, error) {
	decision := &RoutingDecision{
		OrderID:  order.ID,
		Strategy: sor.config.RoutingStrategy,
	}

	// Sort rules by priority
	rules := make([]*RoutingRule, len(sor.routingRules))
	copy(rules, sor.routingRules)
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Priority < rules[j].Priority
	})

	// Apply rules in priority order
	for _, rule := range rules {
		if !rule.IsActive {
			continue
		}

		if sor.evaluateConditions(rule.Conditions, order) {
			sor.executeActions(rule.Actions, decision, order)
			break // First matching rule wins
		}
	}

	return decision, nil
}

// selectOptimalVenues selects optimal venues based on strategy
func (sor *SmartOrderRouter) selectOptimalVenues(ctx context.Context, order *ExecutionOrder) ([]*VenueAllocation, error) {
	availableVenues := sor.getAvailableVenues(order.Symbol)
	if len(availableVenues) == 0 {
		return nil, fmt.Errorf("no available venues for symbol: %s", order.Symbol)
	}

	switch sor.config.RoutingStrategy {
	case RoutingStrategyBestPrice:
		return sor.selectByBestPrice(availableVenues, order)
	case RoutingStrategyLowestLatency:
		return sor.selectByLowestLatency(availableVenues, order)
	case RoutingStrategyHighestFill:
		return sor.selectByHighestFill(availableVenues, order)
	case RoutingStrategyLowestCost:
		return sor.selectByLowestCost(availableVenues, order)
	case RoutingStrategyBalanced:
		return sor.selectByBalanced(availableVenues, order)
	case RoutingStrategyLiquidity:
		return sor.selectByLiquidity(availableVenues, order)
	default:
		return sor.selectByBalanced(availableVenues, order)
	}
}

// getAvailableVenues returns venues available for a symbol
func (sor *SmartOrderRouter) getAvailableVenues(symbol string) []*VenueInfo {
	var available []*VenueInfo
	for _, venue := range sor.venues {
		if venue.IsActive && sor.supportsSymbol(venue, symbol) {
			available = append(available, venue)
		}
	}
	return available
}

// supportsSymbol checks if venue supports a symbol
func (sor *SmartOrderRouter) supportsSymbol(venue *VenueInfo, symbol string) bool {
	for _, supportedSymbol := range venue.SupportedSymbols {
		if supportedSymbol == symbol {
			return true
		}
	}
	return false
}

// selectByBestPrice selects venues by best price
func (sor *SmartOrderRouter) selectByBestPrice(venues []*VenueInfo, order *ExecutionOrder) ([]*VenueAllocation, error) {
	// Sort by fee rate (ascending)
	sort.Slice(venues, func(i, j int) bool {
		return venues[i].FeeRate.LessThan(venues[j].FeeRate)
	})

	// Select top venue(s)
	maxVenues := min(sor.config.MaxVenuesPerOrder, len(venues))
	allocations := make([]*VenueAllocation, 0, maxVenues)

	for i := 0; i < maxVenues; i++ {
		venue := venues[i]
		percentage := decimal.NewFromFloat(1.0).Div(decimal.NewFromInt(int64(maxVenues)))
		quantity := order.Quantity.Mul(percentage)

		allocation := &VenueAllocation{
			VenueID:    venue.ID,
			VenueName:  venue.Name,
			Quantity:   quantity,
			Percentage: percentage,
			Priority:   i + 1,
			Reason:     "best_price",
		}
		allocations = append(allocations, allocation)
	}

	return allocations, nil
}

// selectByLowestLatency selects venues by lowest latency
func (sor *SmartOrderRouter) selectByLowestLatency(venues []*VenueInfo, order *ExecutionOrder) ([]*VenueAllocation, error) {
	// Sort by latency (ascending)
	sort.Slice(venues, func(i, j int) bool {
		return venues[i].Latency < venues[j].Latency
	})

	// Select venue with lowest latency
	venue := venues[0]
	allocation := &VenueAllocation{
		VenueID:    venue.ID,
		VenueName:  venue.Name,
		Quantity:   order.Quantity,
		Percentage: decimal.NewFromFloat(1.0),
		Priority:   1,
		Reason:     "lowest_latency",
	}

	return []*VenueAllocation{allocation}, nil
}

// selectByHighestFill selects venues by highest fill rate
func (sor *SmartOrderRouter) selectByHighestFill(venues []*VenueInfo, order *ExecutionOrder) ([]*VenueAllocation, error) {
	// Sort by fill rate (descending)
	sort.Slice(venues, func(i, j int) bool {
		return venues[i].HistoricalFillRate.GreaterThan(venues[j].HistoricalFillRate)
	})

	// Select venue with highest fill rate
	venue := venues[0]
	allocation := &VenueAllocation{
		VenueID:    venue.ID,
		VenueName:  venue.Name,
		Quantity:   order.Quantity,
		Percentage: decimal.NewFromFloat(1.0),
		Priority:   1,
		Reason:     "highest_fill",
	}

	return []*VenueAllocation{allocation}, nil
}

// selectByLowestCost selects venues by lowest total cost
func (sor *SmartOrderRouter) selectByLowestCost(venues []*VenueInfo, order *ExecutionOrder) ([]*VenueAllocation, error) {
	// Calculate total cost (fees + estimated slippage)
	type venueScore struct {
		venue *VenueInfo
		cost  decimal.Decimal
	}

	scores := make([]venueScore, 0, len(venues))
	for _, venue := range venues {
		totalCost := venue.FeeRate.Add(venue.AverageSlippage)
		scores = append(scores, venueScore{venue: venue, cost: totalCost})
	}

	// Sort by total cost (ascending)
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].cost.LessThan(scores[j].cost)
	})

	// Select venue with lowest cost
	venue := scores[0].venue
	allocation := &VenueAllocation{
		VenueID:    venue.ID,
		VenueName:  venue.Name,
		Quantity:   order.Quantity,
		Percentage: decimal.NewFromFloat(1.0),
		Priority:   1,
		Reason:     "lowest_cost",
	}

	return []*VenueAllocation{allocation}, nil
}

// selectByBalanced selects venues using balanced approach
func (sor *SmartOrderRouter) selectByBalanced(venues []*VenueInfo, order *ExecutionOrder) ([]*VenueAllocation, error) {
	// Calculate composite score for each venue
	type venueScore struct {
		venue *VenueInfo
		score decimal.Decimal
	}

	scores := make([]venueScore, 0, len(venues))
	for _, venue := range venues {
		// Composite score: reliability * (1 - fee_rate) * fill_rate * (1 / latency_factor)
		latencyFactor := decimal.NewFromFloat(float64(venue.Latency.Milliseconds())).Div(decimal.NewFromFloat(1000))
		if latencyFactor.IsZero() {
			latencyFactor = decimal.NewFromFloat(0.001)
		}

		score := venue.ReliabilityScore.
			Mul(decimal.NewFromFloat(1.0).Sub(venue.FeeRate)).
			Mul(venue.HistoricalFillRate).
			Mul(decimal.NewFromFloat(1.0).Div(latencyFactor))

		scores = append(scores, venueScore{venue: venue, score: score})
	}

	// Sort by score (descending)
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score.GreaterThan(scores[j].score)
	})

	// Select top venues based on weights
	maxVenues := min(sor.config.MaxVenuesPerOrder, len(scores))
	allocations := make([]*VenueAllocation, 0, maxVenues)

	totalWeight := decimal.Zero
	for i := 0; i < maxVenues; i++ {
		totalWeight = totalWeight.Add(scores[i].score)
	}

	for i := 0; i < maxVenues; i++ {
		venue := scores[i].venue
		percentage := scores[i].score.Div(totalWeight)
		quantity := order.Quantity.Mul(percentage)

		allocation := &VenueAllocation{
			VenueID:    venue.ID,
			VenueName:  venue.Name,
			Quantity:   quantity,
			Percentage: percentage,
			Priority:   i + 1,
			Reason:     "balanced_score",
		}
		allocations = append(allocations, allocation)
	}

	return allocations, nil
}

// selectByLiquidity selects venues by liquidity
func (sor *SmartOrderRouter) selectByLiquidity(venues []*VenueInfo, order *ExecutionOrder) ([]*VenueAllocation, error) {
	// Sort by liquidity for the symbol (descending)
	sort.Slice(venues, func(i, j int) bool {
		liquidity1 := venues[i].Liquidity[order.Symbol]
		liquidity2 := venues[j].Liquidity[order.Symbol]
		return liquidity1.GreaterThan(liquidity2)
	})

	// Select venue with highest liquidity
	venue := venues[0]
	allocation := &VenueAllocation{
		VenueID:    venue.ID,
		VenueName:  venue.Name,
		Quantity:   order.Quantity,
		Percentage: decimal.NewFromFloat(1.0),
		Priority:   1,
		Reason:     "highest_liquidity",
	}

	return []*VenueAllocation{allocation}, nil
}

// evaluateConditions evaluates routing rule conditions
func (sor *SmartOrderRouter) evaluateConditions(conditions []*RoutingCondition, order *ExecutionOrder) bool {
	for _, condition := range conditions {
		if !sor.evaluateCondition(condition, order) {
			return false
		}
	}
	return true
}

// evaluateCondition evaluates a single routing condition
func (sor *SmartOrderRouter) evaluateCondition(condition *RoutingCondition, order *ExecutionOrder) bool {
	switch condition.Field {
	case "symbol":
		return sor.compareValues(order.Symbol, condition.Operator, condition.Value)
	case "quantity":
		return sor.compareValues(order.Quantity.String(), condition.Operator, condition.Value)
	case "side":
		return sor.compareValues(string(order.Side), condition.Operator, condition.Value)
	case "order_type":
		return sor.compareValues(string(order.OrderType), condition.Operator, condition.Value)
	default:
		return false
	}
}

// compareValues compares values based on operator
func (sor *SmartOrderRouter) compareValues(actual string, operator string, expected interface{}) bool {
	expectedStr := fmt.Sprintf("%v", expected)

	switch operator {
	case "equals":
		return actual == expectedStr
	case "not_equals":
		return actual != expectedStr
	case "contains":
		return fmt.Sprintf("%s", actual) == expectedStr
	default:
		return false
	}
}

// executeActions executes routing rule actions
func (sor *SmartOrderRouter) executeActions(actions []*RoutingAction, decision *RoutingDecision, order *ExecutionOrder) {
	for _, action := range actions {
		switch action.Type {
		case ActionTypeRoute:
			if action.VenueID != "" {
				allocation := &VenueAllocation{
					VenueID:    action.VenueID,
					Quantity:   order.Quantity,
					Percentage: decimal.NewFromFloat(1.0),
					Priority:   1,
					Reason:     "routing_rule",
				}
				decision.SelectedVenues = append(decision.SelectedVenues, allocation)
			}
		case ActionTypeReject:
			decision.RejectionReason = "rejected_by_rule"
		}
	}
}

// calculateEstimates calculates cost and latency estimates
func (sor *SmartOrderRouter) calculateEstimates(decision *RoutingDecision, order *ExecutionOrder) {
	var totalCost decimal.Decimal
	var maxLatency time.Duration

	for _, allocation := range decision.SelectedVenues {
		venue, exists := sor.venues[allocation.VenueID]
		if !exists {
			continue
		}

		// Calculate cost for this allocation
		notionalValue := allocation.Quantity.Mul(order.Price)
		allocationCost := notionalValue.Mul(venue.FeeRate)
		totalCost = totalCost.Add(allocationCost)

		// Track maximum latency
		if venue.Latency > maxLatency {
			maxLatency = venue.Latency
		}
	}

	decision.EstimatedCost = totalCost
	decision.EstimatedLatency = maxLatency
	decision.ConfidenceScore = decimal.NewFromFloat(0.85) // Simplified confidence score
}

// updateRoutingMetrics updates routing metrics
func (sor *SmartOrderRouter) updateRoutingMetrics(decision *RoutingDecision) {
	sor.metrics.TotalOrders++

	if len(decision.SelectedVenues) > 0 {
		sor.metrics.RoutedOrders++

		if len(decision.SelectedVenues) > 1 {
			sor.metrics.SplitOrders++
		}
	} else {
		sor.metrics.RejectedOrders++
	}

	// Update strategy metrics
	strategyKey := string(decision.Strategy)
	if sor.metrics.StrategyMetrics[strategyKey] == nil {
		sor.metrics.StrategyMetrics[strategyKey] = &StrategyMetrics{
			LastUpdated: time.Now(),
		}
	}
	sor.metrics.StrategyMetrics[strategyKey].OrderCount++

	sor.metrics.LastUpdated = time.Now()
}

// initializeDefaultVenues initializes default trading venues
func (sor *SmartOrderRouter) initializeDefaultVenues() {
	defaultVenues := []*VenueInfo{
		{
			ID:                 "binance",
			Name:               "Binance",
			Type:               VenueTypeExchange,
			IsActive:           true,
			Latency:            50 * time.Millisecond,
			FeeRate:            decimal.NewFromFloat(0.001), // 0.1%
			MinOrderSize:       decimal.NewFromFloat(0.001),
			MaxOrderSize:       decimal.NewFromFloat(1000000),
			SupportedSymbols:   []string{"BTC/USD", "ETH/USD", "BNB/USD"},
			HistoricalFillRate: decimal.NewFromFloat(0.95),
			AverageSlippage:    decimal.NewFromFloat(0.0005),
			ReliabilityScore:   decimal.NewFromFloat(0.98),
			Liquidity: map[string]decimal.Decimal{
				"BTC/USD": decimal.NewFromFloat(1000000),
				"ETH/USD": decimal.NewFromFloat(500000),
				"BNB/USD": decimal.NewFromFloat(200000),
			},
		},
		{
			ID:                 "coinbase",
			Name:               "Coinbase Pro",
			Type:               VenueTypeExchange,
			IsActive:           true,
			Latency:            75 * time.Millisecond,
			FeeRate:            decimal.NewFromFloat(0.005), // 0.5%
			MinOrderSize:       decimal.NewFromFloat(0.001),
			MaxOrderSize:       decimal.NewFromFloat(500000),
			SupportedSymbols:   []string{"BTC/USD", "ETH/USD"},
			HistoricalFillRate: decimal.NewFromFloat(0.92),
			AverageSlippage:    decimal.NewFromFloat(0.0008),
			ReliabilityScore:   decimal.NewFromFloat(0.96),
			Liquidity: map[string]decimal.Decimal{
				"BTC/USD": decimal.NewFromFloat(800000),
				"ETH/USD": decimal.NewFromFloat(400000),
			},
		},
		{
			ID:                 "kraken",
			Name:               "Kraken",
			Type:               VenueTypeExchange,
			IsActive:           true,
			Latency:            100 * time.Millisecond,
			FeeRate:            decimal.NewFromFloat(0.0026), // 0.26%
			MinOrderSize:       decimal.NewFromFloat(0.001),
			MaxOrderSize:       decimal.NewFromFloat(200000),
			SupportedSymbols:   []string{"BTC/USD", "ETH/USD"},
			HistoricalFillRate: decimal.NewFromFloat(0.90),
			AverageSlippage:    decimal.NewFromFloat(0.001),
			ReliabilityScore:   decimal.NewFromFloat(0.94),
			Liquidity: map[string]decimal.Decimal{
				"BTC/USD": decimal.NewFromFloat(600000),
				"ETH/USD": decimal.NewFromFloat(300000),
			},
		},
	}

	for _, venue := range defaultVenues {
		sor.RegisterVenue(venue)
	}
}

// initializeDefaultRules initializes default routing rules
func (sor *SmartOrderRouter) initializeDefaultRules() {
	// Default rules are already initialized in initializeDefaultVenues
}

// venueMonitoringLoop monitors venue performance
func (sor *SmartOrderRouter) venueMonitoringLoop(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-sor.stopChan:
			return
		case <-ticker.C:
			// Update venue metrics (simplified)
		}
	}
}

// metricsUpdateLoop updates routing metrics
func (sor *SmartOrderRouter) metricsUpdateLoop(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-sor.stopChan:
			return
		case <-ticker.C:
			// Update aggregate metrics (simplified)
		}
	}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
