package hft

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// AdvancedRiskManager provides real-time risk controls with position limits,
// exposure monitoring, circuit breakers, and automated risk mitigation
type AdvancedRiskManager struct {
	logger *observability.Logger
	config RiskConfig

	// Risk monitoring components
	positionMonitor *PositionMonitor
	exposureTracker *ExposureTracker
	circuitBreakers *CircuitBreakerManager
	riskCalculator  *RiskCalculator
	alertManager    *AlertManager

	// Real-time risk metrics
	currentRisk      *RiskMetrics
	riskLimits       *RiskLimits
	violationHistory []RiskViolation

	// Performance tracking
	checksPerformed    int64
	violationsDetected int64
	ordersBlocked      int64
	avgCheckTimeNanos  int64

	// State management
	isRunning int32
	stopChan  chan struct{}
	wg        sync.WaitGroup
	mu        sync.RWMutex

	// Risk event subscribers
	subscribers map[string][]chan *RiskEvent
}

// RiskConfig contains configuration for advanced risk management
type RiskConfig struct {
	// Position limits
	MaxPositionSize   decimal.Decimal `json:"max_position_size"`
	MaxPortfolioValue decimal.Decimal `json:"max_portfolio_value"`
	MaxConcentration  float64         `json:"max_concentration"` // Max % in single position
	MaxLeverage       float64         `json:"max_leverage"`      // Maximum leverage ratio

	// Loss limits
	MaxDailyLoss   decimal.Decimal `json:"max_daily_loss"`
	MaxWeeklyLoss  decimal.Decimal `json:"max_weekly_loss"`
	MaxMonthlyLoss decimal.Decimal `json:"max_monthly_loss"`
	MaxDrawdown    float64         `json:"max_drawdown"` // Max drawdown %

	// Order controls
	MaxOrderSize       decimal.Decimal `json:"max_order_size"`
	MaxOrderValue      decimal.Decimal `json:"max_order_value"`
	MaxOrdersPerSecond int             `json:"max_orders_per_second"`
	MaxOrdersPerMinute int             `json:"max_orders_per_minute"`

	// Market risk
	MaxMarketImpact       float64         `json:"max_market_impact"`       // Max market impact %
	MaxVolatilityExposure float64         `json:"max_volatility_exposure"` // Max volatility exposure
	VaRLimit              decimal.Decimal `json:"var_limit"`               // Value at Risk limit

	// Circuit breaker settings
	EnableCircuitBreakers bool    `json:"enable_circuit_breakers"`
	PriceChangeThreshold  float64 `json:"price_change_threshold"` // % price change
	VolumeThreshold       float64 `json:"volume_threshold"`       // Volume spike threshold
	VolatilityThreshold   float64 `json:"volatility_threshold"`   // Volatility threshold

	// Response settings
	AutoStopOnViolation bool          `json:"auto_stop_on_violation"`
	AlertThresholds     []float64     `json:"alert_thresholds"`     // Risk level thresholds
	CheckInterval       time.Duration `json:"check_interval"`       // Risk check frequency
	ResponseTimeTarget  time.Duration `json:"response_time_target"` // Target response time
}

// RiskMetrics contains current risk measurements
type RiskMetrics struct {
	// Portfolio metrics
	TotalValue    decimal.Decimal `json:"total_value"`
	TotalExposure decimal.Decimal `json:"total_exposure"`
	NetPosition   decimal.Decimal `json:"net_position"`
	Leverage      float64         `json:"leverage"`
	Concentration float64         `json:"concentration"`

	// P&L metrics
	DailyPnL      decimal.Decimal `json:"daily_pnl"`
	WeeklyPnL     decimal.Decimal `json:"weekly_pnl"`
	MonthlyPnL    decimal.Decimal `json:"monthly_pnl"`
	UnrealizedPnL decimal.Decimal `json:"unrealized_pnl"`
	MaxDrawdown   float64         `json:"max_drawdown"`

	// Risk measures
	VaR95              decimal.Decimal `json:"var_95"`              // 95% Value at Risk
	VaR99              decimal.Decimal `json:"var_99"`              // 99% Value at Risk
	ExpectedShortfall  decimal.Decimal `json:"expected_shortfall"`  // Conditional VaR
	BetaExposure       float64         `json:"beta_exposure"`       // Market beta exposure
	VolatilityExposure float64         `json:"volatility_exposure"` // Volatility exposure

	// Order flow metrics
	OrdersPerSecond int             `json:"orders_per_second"`
	OrdersPerMinute int             `json:"orders_per_minute"`
	AvgOrderSize    decimal.Decimal `json:"avg_order_size"`
	MarketImpact    float64         `json:"market_impact"`

	// Timestamps
	LastUpdate      time.Time     `json:"last_update"`
	CalculationTime time.Duration `json:"calculation_time"`
}

// RiskLimits defines all risk limits and thresholds
type RiskLimits struct {
	// Position limits
	PositionLimits map[string]decimal.Decimal `json:"position_limits"` // Per symbol
	SectorLimits   map[string]decimal.Decimal `json:"sector_limits"`   // Per sector
	CountryLimits  map[string]decimal.Decimal `json:"country_limits"`  // Per country

	// Dynamic limits based on volatility
	VolatilityAdjusted   bool            `json:"volatility_adjusted"`
	BaseLimit            decimal.Decimal `json:"base_limit"`
	VolatilityMultiplier float64         `json:"volatility_multiplier"`

	// Time-based limits
	IntradayLimits  map[string]decimal.Decimal `json:"intraday_limits"`
	OvernightLimits map[string]decimal.Decimal `json:"overnight_limits"`

	// Strategy-specific limits
	StrategyLimits map[string]RiskLimit `json:"strategy_limits"`

	// Emergency limits
	EmergencyMode   bool                       `json:"emergency_mode"`
	EmergencyLimits map[string]decimal.Decimal `json:"emergency_limits"`
}

// RiskLimit represents a specific risk limit
type RiskLimit struct {
	Symbol         string          `json:"symbol"`
	MaxPosition    decimal.Decimal `json:"max_position"`
	MaxOrderSize   decimal.Decimal `json:"max_order_size"`
	MaxDailyVolume decimal.Decimal `json:"max_daily_volume"`
	MaxLoss        decimal.Decimal `json:"max_loss"`
	IsActive       bool            `json:"is_active"`
	LastUpdated    time.Time       `json:"last_updated"`
}

// RiskViolation represents a risk limit violation
type RiskViolation struct {
	ID           uuid.UUID       `json:"id"`
	Type         ViolationType   `json:"type"`
	Severity     Severity        `json:"severity"`
	Symbol       string          `json:"symbol"`
	Description  string          `json:"description"`
	CurrentValue decimal.Decimal `json:"current_value"`
	LimitValue   decimal.Decimal `json:"limit_value"`
	ExcessAmount decimal.Decimal `json:"excess_amount"`
	Timestamp    time.Time       `json:"timestamp"`
	ActionTaken  string          `json:"action_taken"`
	Resolved     bool            `json:"resolved"`
	ResolvedAt   time.Time       `json:"resolved_at,omitempty"`
}

// ViolationType represents different types of risk violations
type ViolationType string

const (
	ViolationTypePosition      ViolationType = "POSITION_LIMIT"
	ViolationTypeExposure      ViolationType = "EXPOSURE_LIMIT"
	ViolationTypeLoss          ViolationType = "LOSS_LIMIT"
	ViolationTypeConcentration ViolationType = "CONCENTRATION_LIMIT"
	ViolationTypeLeverage      ViolationType = "LEVERAGE_LIMIT"
	ViolationTypeOrderSize     ViolationType = "ORDER_SIZE_LIMIT"
	ViolationTypeOrderRate     ViolationType = "ORDER_RATE_LIMIT"
	ViolationTypeMarketImpact  ViolationType = "MARKET_IMPACT_LIMIT"
	ViolationTypeVaR           ViolationType = "VAR_LIMIT"
	ViolationTypeDrawdown      ViolationType = "DRAWDOWN_LIMIT"
)

// Severity represents the severity level of a risk violation
type Severity string

const (
	SeverityLow      Severity = "LOW"
	SeverityMedium   Severity = "MEDIUM"
	SeverityHigh     Severity = "HIGH"
	SeverityCritical Severity = "CRITICAL"
)

// RiskEvent represents a risk-related event
type RiskEvent struct {
	ID        uuid.UUID              `json:"id"`
	Type      RiskEventType          `json:"type"`
	Severity  Severity               `json:"severity"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
	Source    string                 `json:"source"`
}

// RiskEventType represents different types of risk events
type RiskEventType string

const (
	RiskEventViolation      RiskEventType = "VIOLATION"
	RiskEventAlert          RiskEventType = "ALERT"
	RiskEventCircuitBreaker RiskEventType = "CIRCUIT_BREAKER"
	RiskEventLimitUpdate    RiskEventType = "LIMIT_UPDATE"
	RiskEventEmergencyStop  RiskEventType = "EMERGENCY_STOP"
)

// NewAdvancedRiskManager creates a new advanced risk management system
func NewAdvancedRiskManager(logger *observability.Logger, config RiskConfig) *AdvancedRiskManager {
	// Set default values
	if config.MaxPositionSize.IsZero() {
		config.MaxPositionSize = decimal.NewFromFloat(1000000) // $1M default
	}
	if config.MaxDailyLoss.IsZero() {
		config.MaxDailyLoss = decimal.NewFromFloat(50000) // $50K default
	}
	if config.MaxOrderSize.IsZero() {
		config.MaxOrderSize = decimal.NewFromFloat(100000) // $100K default
	}
	if config.CheckInterval == 0 {
		config.CheckInterval = 100 * time.Millisecond
	}
	if config.ResponseTimeTarget == 0 {
		config.ResponseTimeTarget = 1 * time.Millisecond
	}

	arm := &AdvancedRiskManager{
		logger:           logger,
		config:           config,
		currentRisk:      &RiskMetrics{},
		riskLimits:       &RiskLimits{},
		violationHistory: make([]RiskViolation, 0),
		subscribers:      make(map[string][]chan *RiskEvent),
		stopChan:         make(chan struct{}),
	}

	// Initialize components
	arm.positionMonitor = NewPositionMonitor(logger, config)
	arm.exposureTracker = NewExposureTracker(logger, config)
	arm.circuitBreakers = NewCircuitBreakerManager(logger, config)
	arm.riskCalculator = NewRiskCalculator(logger, config)
	arm.alertManager = NewAlertManager(logger, config)

	// Initialize default limits
	arm.initializeDefaultLimits()

	return arm
}

// Start begins the advanced risk management system
func (arm *AdvancedRiskManager) Start(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&arm.isRunning, 0, 1) {
		return fmt.Errorf("advanced risk manager is already running")
	}

	arm.logger.Info(ctx, "Starting advanced risk management system", map[string]interface{}{
		"max_position_size":    arm.config.MaxPositionSize.String(),
		"max_daily_loss":       arm.config.MaxDailyLoss.String(),
		"check_interval":       arm.config.CheckInterval.String(),
		"circuit_breakers":     arm.config.EnableCircuitBreakers,
		"response_time_target": arm.config.ResponseTimeTarget.String(),
	})

	// Start monitoring components
	if err := arm.startComponents(ctx); err != nil {
		return fmt.Errorf("failed to start risk components: %w", err)
	}

	// Start processing threads
	arm.wg.Add(4)
	go arm.monitorRisk(ctx)
	go arm.processViolations(ctx)
	go arm.updateMetrics(ctx)
	go arm.performanceMonitor(ctx)

	arm.logger.Info(ctx, "Advanced risk management system started successfully", nil)
	return nil
}

// Stop gracefully shuts down the advanced risk management system
func (arm *AdvancedRiskManager) Stop(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&arm.isRunning, 1, 0) {
		return fmt.Errorf("advanced risk manager is not running")
	}

	arm.logger.Info(ctx, "Stopping advanced risk management system", nil)

	close(arm.stopChan)
	arm.wg.Wait()

	// Stop components
	arm.stopComponents(ctx)

	arm.logger.Info(ctx, "Advanced risk management system stopped", map[string]interface{}{
		"checks_performed":     atomic.LoadInt64(&arm.checksPerformed),
		"violations_detected":  atomic.LoadInt64(&arm.violationsDetected),
		"orders_blocked":       atomic.LoadInt64(&arm.ordersBlocked),
		"avg_check_time_nanos": atomic.LoadInt64(&arm.avgCheckTimeNanos),
	})

	return nil
}

// ValidateOrder performs comprehensive pre-trade risk validation
func (arm *AdvancedRiskManager) ValidateOrder(ctx context.Context, order *OrderRequest) error {
	if atomic.LoadInt32(&arm.isRunning) != 1 {
		return fmt.Errorf("risk manager is not running")
	}

	start := time.Now()
	defer func() {
		checkTime := time.Since(start).Nanoseconds()
		atomic.StoreInt64(&arm.avgCheckTimeNanos, checkTime)
		atomic.AddInt64(&arm.checksPerformed, 1)
	}()

	arm.logger.Debug(ctx, "Validating order", map[string]interface{}{
		"symbol":   order.Symbol,
		"side":     string(order.Side),
		"quantity": order.Quantity.String(),
		"type":     string(order.Type),
	})

	// Check order size limits
	if err := arm.validateOrderSize(ctx, order); err != nil {
		atomic.AddInt64(&arm.ordersBlocked, 1)
		arm.recordViolation(ctx, ViolationTypeOrderSize, SeverityMedium, order.Symbol, err.Error(), order.Quantity, arm.config.MaxOrderSize)
		return err
	}

	// Check position limits
	if err := arm.validatePositionLimits(ctx, order); err != nil {
		atomic.AddInt64(&arm.ordersBlocked, 1)
		arm.recordViolation(ctx, ViolationTypePosition, SeverityHigh, order.Symbol, err.Error(), decimal.Zero, decimal.Zero)
		return err
	}

	// Check exposure limits
	if err := arm.validateExposureLimits(ctx, order); err != nil {
		atomic.AddInt64(&arm.ordersBlocked, 1)
		arm.recordViolation(ctx, ViolationTypeExposure, SeverityHigh, order.Symbol, err.Error(), decimal.Zero, decimal.Zero)
		return err
	}

	// Check concentration limits
	if err := arm.validateConcentrationLimits(ctx, order); err != nil {
		atomic.AddInt64(&arm.ordersBlocked, 1)
		arm.recordViolation(ctx, ViolationTypeConcentration, SeverityMedium, order.Symbol, err.Error(), decimal.Zero, decimal.Zero)
		return err
	}

	// Check order rate limits
	if err := arm.validateOrderRateLimits(ctx, order); err != nil {
		atomic.AddInt64(&arm.ordersBlocked, 1)
		arm.recordViolation(ctx, ViolationTypeOrderRate, SeverityLow, order.Symbol, err.Error(), decimal.Zero, decimal.Zero)
		return err
	}

	// Check market impact
	if err := arm.validateMarketImpact(ctx, order); err != nil {
		atomic.AddInt64(&arm.ordersBlocked, 1)
		arm.recordViolation(ctx, ViolationTypeMarketImpact, SeverityMedium, order.Symbol, err.Error(), decimal.Zero, decimal.Zero)
		return err
	}

	// Check circuit breakers
	if arm.config.EnableCircuitBreakers {
		if err := arm.circuitBreakers.CheckOrder(ctx, order); err != nil {
			atomic.AddInt64(&arm.ordersBlocked, 1)
			arm.publishRiskEvent(ctx, RiskEventCircuitBreaker, SeverityHigh, "Circuit breaker triggered", map[string]interface{}{
				"symbol": order.Symbol,
				"reason": err.Error(),
			})
			return err
		}
	}

	arm.logger.Debug(ctx, "Order validation passed", map[string]interface{}{
		"symbol":     order.Symbol,
		"check_time": time.Since(start).String(),
	})

	return nil
}

// validateOrderSize validates order size against limits
func (arm *AdvancedRiskManager) validateOrderSize(ctx context.Context, order *OrderRequest) error {
	// Check absolute order size
	if order.Quantity.GreaterThan(arm.config.MaxOrderSize) {
		return fmt.Errorf("order size %s exceeds maximum limit %s",
			order.Quantity.String(), arm.config.MaxOrderSize.String())
	}

	// Check order value (for limit orders)
	if order.Type == OrderTypeLimit && !order.Price.IsZero() {
		orderValue := order.Quantity.Mul(order.Price)
		if orderValue.GreaterThan(arm.config.MaxOrderValue) {
			return fmt.Errorf("order value %s exceeds maximum limit %s",
				orderValue.String(), arm.config.MaxOrderValue.String())
		}
	}

	// Check symbol-specific limits
	arm.mu.RLock()
	if symbolLimit, exists := arm.riskLimits.PositionLimits[order.Symbol]; exists {
		if order.Quantity.GreaterThan(symbolLimit) {
			arm.mu.RUnlock()
			return fmt.Errorf("order size %s exceeds symbol limit %s for %s",
				order.Quantity.String(), symbolLimit.String(), order.Symbol)
		}
	}
	arm.mu.RUnlock()

	return nil
}

// validatePositionLimits validates position limits
func (arm *AdvancedRiskManager) validatePositionLimits(ctx context.Context, order *OrderRequest) error {
	// Get current position
	currentPosition := arm.positionMonitor.GetPosition(order.Symbol)

	// Calculate new position after order
	var newPosition decimal.Decimal
	if order.Side == OrderSideBuy {
		newPosition = currentPosition.Add(order.Quantity)
	} else {
		newPosition = currentPosition.Sub(order.Quantity)
	}

	// Check against maximum position size
	if newPosition.Abs().GreaterThan(arm.config.MaxPositionSize) {
		return fmt.Errorf("new position %s would exceed maximum position size %s",
			newPosition.String(), arm.config.MaxPositionSize.String())
	}

	// Check symbol-specific position limits
	arm.mu.RLock()
	if symbolLimit, exists := arm.riskLimits.PositionLimits[order.Symbol]; exists {
		if newPosition.Abs().GreaterThan(symbolLimit) {
			arm.mu.RUnlock()
			return fmt.Errorf("new position %s would exceed symbol position limit %s for %s",
				newPosition.String(), symbolLimit.String(), order.Symbol)
		}
	}
	arm.mu.RUnlock()

	return nil
}

// validateExposureLimits validates exposure limits
func (arm *AdvancedRiskManager) validateExposureLimits(ctx context.Context, order *OrderRequest) error {
	// Calculate order exposure
	var orderExposure decimal.Decimal
	if order.Type == OrderTypeLimit && !order.Price.IsZero() {
		orderExposure = order.Quantity.Mul(order.Price)
	} else {
		// Use current market price for market orders
		marketPrice := arm.getMarketPrice(order.Symbol)
		orderExposure = order.Quantity.Mul(marketPrice)
	}

	// Get current total exposure
	currentExposure := arm.exposureTracker.GetTotalExposure()
	newExposure := currentExposure.Add(orderExposure)

	// Check against maximum portfolio value
	if newExposure.GreaterThan(arm.config.MaxPortfolioValue) {
		return fmt.Errorf("new exposure %s would exceed maximum portfolio value %s",
			newExposure.String(), arm.config.MaxPortfolioValue.String())
	}

	return nil
}

// validateConcentrationLimits validates concentration limits
func (arm *AdvancedRiskManager) validateConcentrationLimits(ctx context.Context, order *OrderRequest) error {
	// Get all current positions
	positions := arm.positionMonitor.GetAllPositions()

	// Calculate current concentration
	concentration := arm.riskCalculator.CalculateConcentration(positions)

	// Check if concentration exceeds limit
	if concentration > arm.config.MaxConcentration {
		return fmt.Errorf("portfolio concentration %.2f%% exceeds maximum limit %.2f%%",
			concentration, arm.config.MaxConcentration)
	}

	return nil
}

// validateOrderRateLimits validates order rate limits
func (arm *AdvancedRiskManager) validateOrderRateLimits(ctx context.Context, order *OrderRequest) error {
	// This would track order rates in production
	// For now, return nil (no rate limiting)
	return nil
}

// validateMarketImpact validates market impact limits
func (arm *AdvancedRiskManager) validateMarketImpact(ctx context.Context, order *OrderRequest) error {
	// Calculate estimated market impact
	marketImpact := arm.calculateOrderMarketImpact(order)

	// Check against maximum allowed market impact
	if marketImpact > arm.config.MaxMarketImpact {
		return fmt.Errorf("estimated market impact %.2f%% exceeds maximum limit %.2f%%",
			marketImpact, arm.config.MaxMarketImpact)
	}

	return nil
}

// calculateOrderMarketImpact calculates estimated market impact for an order
func (arm *AdvancedRiskManager) calculateOrderMarketImpact(order *OrderRequest) float64 {
	// Simplified market impact calculation
	// In production, this would use sophisticated models

	// Base impact based on order size (larger orders = higher impact)
	baseImpact := order.Quantity.InexactFloat64() / 1000.0 // 0.1% per 1000 units

	// Adjust for market conditions (simplified)
	volatilityMultiplier := 1.0 // Would use actual volatility
	liquidityMultiplier := 1.0  // Would use actual liquidity

	impact := baseImpact * volatilityMultiplier * liquidityMultiplier

	// Cap at reasonable maximum
	if impact > 10.0 {
		impact = 10.0
	}

	return impact
}

// getMarketPrice gets current market price for a symbol
func (arm *AdvancedRiskManager) getMarketPrice(symbol string) decimal.Decimal {
	// Mock price - in production, get from market data feed
	return decimal.NewFromFloat(45000.0)
}

// recordViolation records a risk violation
func (arm *AdvancedRiskManager) recordViolation(ctx context.Context, violationType ViolationType, severity Severity, symbol, description string, currentValue, limitValue decimal.Decimal) {
	violation := RiskViolation{
		ID:           uuid.New(),
		Type:         violationType,
		Severity:     severity,
		Symbol:       symbol,
		Description:  description,
		CurrentValue: currentValue,
		LimitValue:   limitValue,
		ExcessAmount: currentValue.Sub(limitValue),
		Timestamp:    time.Now(),
		ActionTaken:  "ORDER_BLOCKED",
		Resolved:     false,
	}

	arm.mu.Lock()
	arm.violationHistory = append(arm.violationHistory, violation)
	atomic.AddInt64(&arm.violationsDetected, 1)
	arm.mu.Unlock()

	// Publish risk event
	arm.publishRiskEvent(ctx, RiskEventViolation, severity, description, map[string]interface{}{
		"violation_id":   violation.ID.String(),
		"violation_type": string(violationType),
		"symbol":         symbol,
		"current_value":  currentValue.String(),
		"limit_value":    limitValue.String(),
	})

	arm.logger.Warn(ctx, "Risk violation recorded", map[string]interface{}{
		"violation_id": violation.ID.String(),
		"type":         string(violationType),
		"severity":     string(severity),
		"symbol":       symbol,
		"description":  description,
	})
}

// publishRiskEvent publishes a risk event to subscribers
func (arm *AdvancedRiskManager) publishRiskEvent(ctx context.Context, eventType RiskEventType, severity Severity, message string, data map[string]interface{}) {
	event := &RiskEvent{
		ID:        uuid.New(),
		Type:      eventType,
		Severity:  severity,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
		Source:    "AdvancedRiskManager",
	}

	arm.mu.RLock()
	defer arm.mu.RUnlock()

	// Send to event type subscribers
	if subscribers, exists := arm.subscribers[string(eventType)]; exists {
		for _, ch := range subscribers {
			select {
			case ch <- event:
			default:
				// Subscriber channel is full, skip
			}
		}
	}

	// Send to wildcard subscribers
	if subscribers, exists := arm.subscribers["*"]; exists {
		for _, ch := range subscribers {
			select {
			case ch <- event:
			default:
				// Subscriber channel is full, skip
			}
		}
	}
}

// startComponents starts all risk management components
func (arm *AdvancedRiskManager) startComponents(ctx context.Context) error {
	arm.logger.Info(ctx, "Starting risk management components", nil)

	// Components are already initialized, no additional startup needed
	// In production, components might need to connect to data sources

	return nil
}

// stopComponents stops all risk management components
func (arm *AdvancedRiskManager) stopComponents(ctx context.Context) {
	arm.logger.Info(ctx, "Stopping risk management components", nil)

	// Clean shutdown of components
	// In production, this would close connections and clean up resources
}

// initializeDefaultLimits initializes default risk limits
func (arm *AdvancedRiskManager) initializeDefaultLimits() {
	arm.riskLimits.PositionLimits = make(map[string]decimal.Decimal)
	arm.riskLimits.SectorLimits = make(map[string]decimal.Decimal)
	arm.riskLimits.CountryLimits = make(map[string]decimal.Decimal)
	arm.riskLimits.IntradayLimits = make(map[string]decimal.Decimal)
	arm.riskLimits.OvernightLimits = make(map[string]decimal.Decimal)
	arm.riskLimits.StrategyLimits = make(map[string]RiskLimit)
	arm.riskLimits.EmergencyLimits = make(map[string]decimal.Decimal)

	// Set default position limits
	arm.riskLimits.PositionLimits["BTCUSDT"] = decimal.NewFromFloat(100.0)
	arm.riskLimits.PositionLimits["ETHUSDT"] = decimal.NewFromFloat(1000.0)

	// Set default sector limits
	arm.riskLimits.SectorLimits["CRYPTO"] = decimal.NewFromFloat(5000000.0) // $5M

	// Set volatility adjustment parameters
	arm.riskLimits.VolatilityAdjusted = true
	arm.riskLimits.BaseLimit = decimal.NewFromFloat(1000000.0) // $1M base
	arm.riskLimits.VolatilityMultiplier = 0.5                  // Reduce limits by 50% in high volatility
}

// monitorRisk continuously monitors risk metrics
func (arm *AdvancedRiskManager) monitorRisk(ctx context.Context) {
	defer arm.wg.Done()

	arm.logger.Info(ctx, "Starting risk monitor", nil)

	ticker := time.NewTicker(arm.config.CheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-arm.stopChan:
			return
		case <-ticker.C:
			arm.updateRiskMetrics(ctx)
			arm.checkRiskLimits(ctx)
		}
	}
}

// updateRiskMetrics updates current risk metrics
func (arm *AdvancedRiskManager) updateRiskMetrics(ctx context.Context) {
	start := time.Now()

	// Get current positions
	positions := arm.positionMonitor.GetAllPositions()

	// Calculate portfolio metrics
	var totalValue decimal.Decimal
	for _, position := range positions {
		price := arm.getMarketPrice("BTCUSDT") // Simplified
		value := position.Abs().Mul(price)
		totalValue = totalValue.Add(value)
	}

	// Calculate risk metrics
	var95 := arm.riskCalculator.CalculateVaR(positions, 0.95)
	var99 := arm.riskCalculator.CalculateVaR(positions, 0.99)
	concentration := arm.riskCalculator.CalculateConcentration(positions)

	// Update current risk metrics
	arm.mu.Lock()
	arm.currentRisk.TotalValue = totalValue
	arm.currentRisk.TotalExposure = arm.exposureTracker.GetTotalExposure()
	arm.currentRisk.Concentration = concentration
	arm.currentRisk.VaR95 = var95
	arm.currentRisk.VaR99 = var99
	arm.currentRisk.LastUpdate = time.Now()
	arm.currentRisk.CalculationTime = time.Since(start)
	arm.mu.Unlock()
}

// checkRiskLimits checks all risk limits and triggers alerts
func (arm *AdvancedRiskManager) checkRiskLimits(ctx context.Context) {
	arm.mu.RLock()
	currentRisk := *arm.currentRisk
	arm.mu.RUnlock()

	// Check VaR limits
	if !arm.config.VaRLimit.IsZero() && currentRisk.VaR95.GreaterThan(arm.config.VaRLimit) {
		arm.alertManager.CreateAlert(SeverityHigh, fmt.Sprintf("VaR95 %s exceeds limit %s",
			currentRisk.VaR95.String(), arm.config.VaRLimit.String()))
	}

	// Check concentration limits
	if currentRisk.Concentration > arm.config.MaxConcentration {
		arm.alertManager.CreateAlert(SeverityMedium, fmt.Sprintf("Portfolio concentration %.2f%% exceeds limit %.2f%%",
			currentRisk.Concentration, arm.config.MaxConcentration))
	}

	// Check drawdown limits
	if currentRisk.MaxDrawdown > arm.config.MaxDrawdown {
		arm.alertManager.CreateAlert(SeverityCritical, fmt.Sprintf("Drawdown %.2f%% exceeds limit %.2f%%",
			currentRisk.MaxDrawdown, arm.config.MaxDrawdown))

		// Auto-stop if configured
		if arm.config.AutoStopOnViolation {
			arm.publishRiskEvent(ctx, RiskEventEmergencyStop, SeverityCritical, "Emergency stop triggered due to drawdown", map[string]interface{}{
				"drawdown": currentRisk.MaxDrawdown,
				"limit":    arm.config.MaxDrawdown,
			})
		}
	}
}

// processViolations processes risk violations
func (arm *AdvancedRiskManager) processViolations(ctx context.Context) {
	defer arm.wg.Done()

	arm.logger.Info(ctx, "Starting violation processor", nil)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-arm.stopChan:
			return
		case <-ticker.C:
			arm.processActiveViolations(ctx)
		}
	}
}

// processActiveViolations processes active violations
func (arm *AdvancedRiskManager) processActiveViolations(ctx context.Context) {
	arm.mu.RLock()
	violations := make([]RiskViolation, len(arm.violationHistory))
	copy(violations, arm.violationHistory)
	arm.mu.RUnlock()

	for _, violation := range violations {
		if !violation.Resolved {
			// Check if violation should be auto-resolved
			if time.Since(violation.Timestamp) > 5*time.Minute {
				arm.resolveViolation(ctx, violation.ID)
			}
		}
	}
}

// resolveViolation resolves a violation by ID
func (arm *AdvancedRiskManager) resolveViolation(ctx context.Context, violationID uuid.UUID) {
	arm.mu.Lock()
	defer arm.mu.Unlock()

	for i := range arm.violationHistory {
		if arm.violationHistory[i].ID == violationID {
			arm.violationHistory[i].Resolved = true
			arm.violationHistory[i].ResolvedAt = time.Now()

			arm.logger.Info(ctx, "Risk violation resolved", map[string]interface{}{
				"violation_id": violationID.String(),
			})
			break
		}
	}
}

// updateMetrics updates performance metrics
func (arm *AdvancedRiskManager) updateMetrics(ctx context.Context) {
	defer arm.wg.Done()

	arm.logger.Info(ctx, "Starting metrics updater", nil)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-arm.stopChan:
			return
		case <-ticker.C:
			// Update real-time metrics
			arm.updateRiskMetrics(ctx)
		}
	}
}

// performanceMonitor tracks and reports performance metrics
func (arm *AdvancedRiskManager) performanceMonitor(ctx context.Context) {
	defer arm.wg.Done()

	arm.logger.Info(ctx, "Starting risk performance monitor", nil)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	var lastCheckCount int64

	for {
		select {
		case <-arm.stopChan:
			return
		case <-ticker.C:
			currentChecks := atomic.LoadInt64(&arm.checksPerformed)
			checksPerSecond := (currentChecks - lastCheckCount) / 10
			lastCheckCount = currentChecks

			violations := atomic.LoadInt64(&arm.violationsDetected)
			blocked := atomic.LoadInt64(&arm.ordersBlocked)
			avgCheckTime := atomic.LoadInt64(&arm.avgCheckTimeNanos)

			arm.logger.Info(ctx, "Risk management performance", map[string]interface{}{
				"checks_per_second":    checksPerSecond,
				"total_checks":         currentChecks,
				"violations_detected":  violations,
				"orders_blocked":       blocked,
				"avg_check_time_nanos": avgCheckTime,
				"avg_check_time_us":    avgCheckTime / 1000,
			})
		}
	}
}

// Subscribe registers a subscriber for risk events
func (arm *AdvancedRiskManager) Subscribe(eventType string) <-chan *RiskEvent {
	arm.mu.Lock()
	defer arm.mu.Unlock()

	ch := make(chan *RiskEvent, 1000) // Buffered channel
	if arm.subscribers[eventType] == nil {
		arm.subscribers[eventType] = make([]chan *RiskEvent, 0)
	}
	arm.subscribers[eventType] = append(arm.subscribers[eventType], ch)

	return ch
}

// GetMetrics returns current risk metrics
func (arm *AdvancedRiskManager) GetMetrics() *RiskMetrics {
	arm.mu.RLock()
	defer arm.mu.RUnlock()

	// Return a copy of current metrics
	metrics := *arm.currentRisk
	return &metrics
}

// GetViolations returns recent risk violations
func (arm *AdvancedRiskManager) GetViolations(limit int) []RiskViolation {
	arm.mu.RLock()
	defer arm.mu.RUnlock()

	if limit <= 0 || limit > len(arm.violationHistory) {
		limit = len(arm.violationHistory)
	}

	// Return most recent violations
	start := len(arm.violationHistory) - limit
	if start < 0 {
		start = 0
	}

	violations := make([]RiskViolation, limit)
	copy(violations, arm.violationHistory[start:])
	return violations
}

// UpdateLimits updates risk limits
func (arm *AdvancedRiskManager) UpdateLimits(ctx context.Context, limits *RiskLimits) error {
	arm.mu.Lock()
	defer arm.mu.Unlock()

	arm.riskLimits = limits

	arm.publishRiskEvent(ctx, RiskEventLimitUpdate, SeverityLow, "Risk limits updated", map[string]interface{}{
		"updated_at": time.Now(),
	})

	arm.logger.Info(ctx, "Risk limits updated", nil)
	return nil
}

// GetLimits returns current risk limits
func (arm *AdvancedRiskManager) GetLimits() *RiskLimits {
	arm.mu.RLock()
	defer arm.mu.RUnlock()

	// Return a copy of current limits
	limits := *arm.riskLimits
	return &limits
}

// EmergencyStop triggers an emergency stop of all trading
func (arm *AdvancedRiskManager) EmergencyStop(ctx context.Context, reason string) error {
	arm.logger.Error(ctx, "Emergency stop triggered", fmt.Errorf("reason: %s", reason))

	// Set emergency mode
	arm.mu.Lock()
	arm.riskLimits.EmergencyMode = true
	arm.mu.Unlock()

	// Publish emergency stop event
	arm.publishRiskEvent(ctx, RiskEventEmergencyStop, SeverityCritical, "Emergency stop activated", map[string]interface{}{
		"reason":    reason,
		"timestamp": time.Now(),
	})

	return nil
}

// IsEmergencyMode returns whether emergency mode is active
func (arm *AdvancedRiskManager) IsEmergencyMode() bool {
	arm.mu.RLock()
	defer arm.mu.RUnlock()
	return arm.riskLimits.EmergencyMode
}
