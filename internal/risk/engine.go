package risk

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/internal/strategies/framework"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// RiskEngine provides comprehensive risk management and monitoring
type RiskEngine struct {
	logger     *observability.Logger
	config     RiskConfig
	limits     *GlobalRiskLimits
	monitors   map[string]*RiskMonitor
	alerts     chan *RiskAlert
	violations []*RiskViolation

	// Circuit breakers
	circuitBreakers map[string]*CircuitBreaker

	// Performance tracking
	totalChecks    int64
	violationCount int64
	alertCount     int64

	// State management
	isRunning bool
	stopChan  chan struct{}
	wg        sync.WaitGroup
	mu        sync.RWMutex
}

// RiskConfig contains risk engine configuration
type RiskConfig struct {
	CheckInterval        time.Duration   `json:"check_interval"`
	AlertBufferSize      int             `json:"alert_buffer_size"`
	ViolationRetention   time.Duration   `json:"violation_retention"`
	EnableCircuitBreaker bool            `json:"enable_circuit_breaker"`
	EnableStressTest     bool            `json:"enable_stress_test"`
	EnableVaRCalculation bool            `json:"enable_var_calculation"`
	VaRConfidenceLevel   decimal.Decimal `json:"var_confidence_level"`
	VaRTimeHorizon       time.Duration   `json:"var_time_horizon"`
	MaxConcurrentChecks  int             `json:"max_concurrent_checks"`
}

// GlobalRiskLimits defines system-wide risk limits
type GlobalRiskLimits struct {
	MaxTotalExposure    decimal.Decimal `json:"max_total_exposure"`
	MaxDailyLoss        decimal.Decimal `json:"max_daily_loss"`
	MaxDrawdown         decimal.Decimal `json:"max_drawdown"`
	MaxPositionSize     decimal.Decimal `json:"max_position_size"`
	MaxOrdersPerSecond  int             `json:"max_orders_per_second"`
	MaxOrdersPerMinute  int             `json:"max_orders_per_minute"`
	MaxOrdersPerHour    int             `json:"max_orders_per_hour"`
	MaxOrdersPerDay     int             `json:"max_orders_per_day"`
	MaxOpenPositions    int             `json:"max_open_positions"`
	MaxStrategies       int             `json:"max_strategies"`
	MaxLeverage         decimal.Decimal `json:"max_leverage"`
	MinCashReserve      decimal.Decimal `json:"min_cash_reserve"`
	VaRLimit            decimal.Decimal `json:"var_limit"`
	StressTestThreshold decimal.Decimal `json:"stress_test_threshold"`
	AllowedSymbols      []string        `json:"allowed_symbols"`
	BlockedSymbols      []string        `json:"blocked_symbols"`
	TradingHours        *TradingHours   `json:"trading_hours,omitempty"`
}

// TradingHours defines allowed trading hours
type TradingHours struct {
	Enabled   bool      `json:"enabled"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Timezone  string    `json:"timezone"`
	Weekdays  []int     `json:"weekdays"` // 0=Sunday, 1=Monday, etc.
}

// RiskMonitor tracks risk metrics for a specific entity
type RiskMonitor struct {
	ID             string                 `json:"id"`
	Type           RiskMonitorType        `json:"type"`
	EntityID       uuid.UUID              `json:"entity_id"`
	Metrics        *RiskMetrics           `json:"metrics"`
	Limits         *RiskLimits            `json:"limits"`
	LastCheck      time.Time              `json:"last_check"`
	ViolationCount int                    `json:"violation_count"`
	IsActive       bool                   `json:"is_active"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// RiskMonitorType represents the type of risk monitor
type RiskMonitorType string

const (
	RiskMonitorTypeStrategy  RiskMonitorType = "strategy"
	RiskMonitorTypePortfolio RiskMonitorType = "portfolio"
	RiskMonitorTypeSymbol    RiskMonitorType = "symbol"
	RiskMonitorTypeExchange  RiskMonitorType = "exchange"
	RiskMonitorTypeGlobal    RiskMonitorType = "global"
)

// RiskMetrics contains current risk metrics
type RiskMetrics struct {
	TotalExposure   decimal.Decimal `json:"total_exposure"`
	NetExposure     decimal.Decimal `json:"net_exposure"`
	GrossExposure   decimal.Decimal `json:"gross_exposure"`
	DailyPnL        decimal.Decimal `json:"daily_pnl"`
	UnrealizedPnL   decimal.Decimal `json:"unrealized_pnl"`
	RealizedPnL     decimal.Decimal `json:"realized_pnl"`
	MaxDrawdown     decimal.Decimal `json:"max_drawdown"`
	CurrentDrawdown decimal.Decimal `json:"current_drawdown"`
	VaR             decimal.Decimal `json:"var"`
	Leverage        decimal.Decimal `json:"leverage"`
	CashBalance     decimal.Decimal `json:"cash_balance"`
	PositionCount   int             `json:"position_count"`
	OrderCount      int             `json:"order_count"`
	OrderRate       decimal.Decimal `json:"order_rate"`
	LastUpdated     time.Time       `json:"last_updated"`
}

// RiskLimits defines risk limits for an entity
type RiskLimits struct {
	MaxExposure     decimal.Decimal `json:"max_exposure"`
	MaxDailyLoss    decimal.Decimal `json:"max_daily_loss"`
	MaxDrawdown     decimal.Decimal `json:"max_drawdown"`
	MaxPositionSize decimal.Decimal `json:"max_position_size"`
	MaxPositions    int             `json:"max_positions"`
	MaxOrderRate    int             `json:"max_order_rate"`
	MaxLeverage     decimal.Decimal `json:"max_leverage"`
	MinCashReserve  decimal.Decimal `json:"min_cash_reserve"`
	VaRLimit        decimal.Decimal `json:"var_limit"`
	AllowedSymbols  []string        `json:"allowed_symbols"`
	BlockedSymbols  []string        `json:"blocked_symbols"`
}

// RiskAlert represents a risk alert
type RiskAlert struct {
	ID             uuid.UUID              `json:"id"`
	Type           RiskAlertType          `json:"type"`
	Severity       AlertSeverity          `json:"severity"`
	MonitorID      string                 `json:"monitor_id"`
	EntityID       uuid.UUID              `json:"entity_id"`
	Message        string                 `json:"message"`
	Details        map[string]interface{} `json:"details"`
	Timestamp      time.Time              `json:"timestamp"`
	Acknowledged   bool                   `json:"acknowledged"`
	AcknowledgedBy string                 `json:"acknowledged_by,omitempty"`
	AcknowledgedAt time.Time              `json:"acknowledged_at,omitempty"`
}

// RiskAlertType represents the type of risk alert
type RiskAlertType string

const (
	RiskAlertTypeExposure       RiskAlertType = "exposure"
	RiskAlertTypeLoss           RiskAlertType = "loss"
	RiskAlertTypeDrawdown       RiskAlertType = "drawdown"
	RiskAlertTypeVaR            RiskAlertType = "var"
	RiskAlertTypeOrderRate      RiskAlertType = "order_rate"
	RiskAlertTypePosition       RiskAlertType = "position"
	RiskAlertTypeLeverage       RiskAlertType = "leverage"
	RiskAlertTypeCash           RiskAlertType = "cash"
	RiskAlertTypeCompliance     RiskAlertType = "compliance"
	RiskAlertTypeCircuitBreaker RiskAlertType = "circuit_breaker"
)

// AlertSeverity represents the severity of an alert
type AlertSeverity string

const (
	AlertSeverityInfo     AlertSeverity = "info"
	AlertSeverityWarning  AlertSeverity = "warning"
	AlertSeverityError    AlertSeverity = "error"
	AlertSeverityCritical AlertSeverity = "critical"
)

// RiskViolation represents a risk limit violation
type RiskViolation struct {
	ID            uuid.UUID              `json:"id"`
	MonitorID     string                 `json:"monitor_id"`
	EntityID      uuid.UUID              `json:"entity_id"`
	ViolationType RiskViolationType      `json:"violation_type"`
	LimitValue    decimal.Decimal        `json:"limit_value"`
	ActualValue   decimal.Decimal        `json:"actual_value"`
	Severity      AlertSeverity          `json:"severity"`
	Action        ViolationAction        `json:"action"`
	Message       string                 `json:"message"`
	Details       map[string]interface{} `json:"details"`
	Timestamp     time.Time              `json:"timestamp"`
	Resolved      bool                   `json:"resolved"`
	ResolvedAt    time.Time              `json:"resolved_at,omitempty"`
}

// RiskViolationType represents the type of risk violation
type RiskViolationType string

const (
	RiskViolationTypeExposure  RiskViolationType = "exposure"
	RiskViolationTypeLoss      RiskViolationType = "loss"
	RiskViolationTypeDrawdown  RiskViolationType = "drawdown"
	RiskViolationTypeVaR       RiskViolationType = "var"
	RiskViolationTypeOrderRate RiskViolationType = "order_rate"
	RiskViolationTypePosition  RiskViolationType = "position"
	RiskViolationTypeLeverage  RiskViolationType = "leverage"
	RiskViolationTypeCash      RiskViolationType = "cash"
)

// ViolationAction represents the action taken for a violation
type ViolationAction string

const (
	ViolationActionAlert        ViolationAction = "alert"
	ViolationActionBlock        ViolationAction = "block"
	ViolationActionLiquidate    ViolationAction = "liquidate"
	ViolationActionSuspend      ViolationAction = "suspend"
	ViolationActionCircuitBreak ViolationAction = "circuit_break"
)

// CircuitBreaker implements circuit breaker pattern for risk management
type CircuitBreaker struct {
	ID              string          `json:"id"`
	Type            string          `json:"type"`
	Threshold       decimal.Decimal `json:"threshold"`
	TimeWindow      time.Duration   `json:"time_window"`
	CooldownPeriod  time.Duration   `json:"cooldown_period"`
	State           CircuitState    `json:"state"`
	FailureCount    int             `json:"failure_count"`
	LastFailure     time.Time       `json:"last_failure"`
	LastStateChange time.Time       `json:"last_state_change"`
	TotalTrips      int             `json:"total_trips"`
}

// CircuitState represents circuit breaker state
type CircuitState string

const (
	CircuitStateClosed   CircuitState = "closed"
	CircuitStateOpen     CircuitState = "open"
	CircuitStateHalfOpen CircuitState = "half_open"
)

// NewRiskEngine creates a new risk management engine
func NewRiskEngine(logger *observability.Logger, config RiskConfig) *RiskEngine {
	if config.CheckInterval == 0 {
		config.CheckInterval = 1 * time.Second
	}
	if config.AlertBufferSize == 0 {
		config.AlertBufferSize = 1000
	}
	if config.ViolationRetention == 0 {
		config.ViolationRetention = 24 * time.Hour
	}
	if config.MaxConcurrentChecks == 0 {
		config.MaxConcurrentChecks = 10
	}

	return &RiskEngine{
		logger:          logger,
		config:          config,
		monitors:        make(map[string]*RiskMonitor),
		alerts:          make(chan *RiskAlert, config.AlertBufferSize),
		violations:      make([]*RiskViolation, 0),
		circuitBreakers: make(map[string]*CircuitBreaker),
		stopChan:        make(chan struct{}),
	}
}

// Start starts the risk management engine
func (re *RiskEngine) Start(ctx context.Context) error {
	re.mu.Lock()
	defer re.mu.Unlock()

	if re.isRunning {
		return fmt.Errorf("risk engine is already running")
	}

	re.logger.Info(ctx, "Starting risk management engine", map[string]interface{}{
		"check_interval":         re.config.CheckInterval,
		"enable_circuit_breaker": re.config.EnableCircuitBreaker,
		"enable_var_calculation": re.config.EnableVaRCalculation,
		"max_concurrent_checks":  re.config.MaxConcurrentChecks,
	})

	re.isRunning = true

	// Start risk monitoring loop
	re.wg.Add(1)
	go re.monitorRisk(ctx)

	// Start alert processing loop
	re.wg.Add(1)
	go re.processAlerts(ctx)

	// Start violation cleanup loop
	re.wg.Add(1)
	go re.cleanupViolations(ctx)

	re.logger.Info(ctx, "Risk management engine started", map[string]interface{}{
		"active_monitors": len(re.monitors),
	})

	return nil
}

// Stop stops the risk management engine
func (re *RiskEngine) Stop(ctx context.Context) error {
	re.mu.Lock()
	defer re.mu.Unlock()

	if !re.isRunning {
		return fmt.Errorf("risk engine is not running")
	}

	re.logger.Info(ctx, "Stopping risk management engine", nil)

	close(re.stopChan)
	re.wg.Wait()

	re.isRunning = false

	re.logger.Info(ctx, "Risk management engine stopped", nil)

	return nil
}

// SetGlobalLimits sets global risk limits
func (re *RiskEngine) SetGlobalLimits(limits *GlobalRiskLimits) {
	re.mu.Lock()
	defer re.mu.Unlock()
	re.limits = limits
}

// GetGlobalLimits returns current global risk limits
func (re *RiskEngine) GetGlobalLimits() *GlobalRiskLimits {
	re.mu.RLock()
	defer re.mu.RUnlock()
	if re.limits == nil {
		return nil
	}
	// Return a copy
	limits := *re.limits
	return &limits
}

// RegisterMonitor registers a new risk monitor
func (re *RiskEngine) RegisterMonitor(monitor *RiskMonitor) error {
	re.mu.Lock()
	defer re.mu.Unlock()

	if _, exists := re.monitors[monitor.ID]; exists {
		return fmt.Errorf("monitor already exists: %s", monitor.ID)
	}

	monitor.IsActive = true
	monitor.LastCheck = time.Now()
	re.monitors[monitor.ID] = monitor

	re.logger.Info(context.Background(), "Risk monitor registered", map[string]interface{}{
		"monitor_id":   monitor.ID,
		"monitor_type": string(monitor.Type),
		"entity_id":    monitor.EntityID.String(),
	})

	return nil
}

// UnregisterMonitor removes a risk monitor
func (re *RiskEngine) UnregisterMonitor(monitorID string) error {
	re.mu.Lock()
	defer re.mu.Unlock()

	monitor, exists := re.monitors[monitorID]
	if !exists {
		return fmt.Errorf("monitor not found: %s", monitorID)
	}

	monitor.IsActive = false
	delete(re.monitors, monitorID)

	re.logger.Info(context.Background(), "Risk monitor unregistered", map[string]interface{}{
		"monitor_id": monitorID,
	})

	return nil
}

// CheckSignal validates a trading signal against risk limits
func (re *RiskEngine) CheckSignal(ctx context.Context, signal *framework.Signal) error {
	re.totalChecks++

	// Check global limits first
	if err := re.checkGlobalLimits(ctx, signal); err != nil {
		return re.handleViolation(ctx, "global", signal.StrategyID, err)
	}

	// Check strategy-specific limits
	strategyMonitorID := fmt.Sprintf("strategy_%s", signal.StrategyID.String())
	if monitor, exists := re.monitors[strategyMonitorID]; exists {
		if err := re.checkStrategyLimits(ctx, monitor, signal); err != nil {
			return re.handleViolation(ctx, strategyMonitorID, signal.StrategyID, err)
		}
	}

	// Check symbol-specific limits
	symbolMonitorID := fmt.Sprintf("symbol_%s", signal.Symbol)
	if monitor, exists := re.monitors[symbolMonitorID]; exists {
		if err := re.checkSymbolLimits(ctx, monitor, signal); err != nil {
			return re.handleViolation(ctx, symbolMonitorID, signal.StrategyID, err)
		}
	}

	// Check circuit breakers
	if re.config.EnableCircuitBreaker {
		if err := re.checkCircuitBreakers(ctx, signal); err != nil {
			return err
		}
	}

	return nil
}

// UpdateMetrics updates risk metrics for a monitor
func (re *RiskEngine) UpdateMetrics(monitorID string, metrics *RiskMetrics) error {
	re.mu.Lock()
	defer re.mu.Unlock()

	monitor, exists := re.monitors[monitorID]
	if !exists {
		return fmt.Errorf("monitor not found: %s", monitorID)
	}

	monitor.Metrics = metrics
	monitor.LastCheck = time.Now()

	return nil
}

// GetMetrics returns current risk metrics for a monitor
func (re *RiskEngine) GetMetrics(monitorID string) (*RiskMetrics, error) {
	re.mu.RLock()
	defer re.mu.RUnlock()

	monitor, exists := re.monitors[monitorID]
	if !exists {
		return nil, fmt.Errorf("monitor not found: %s", monitorID)
	}

	if monitor.Metrics == nil {
		return nil, fmt.Errorf("no metrics available for monitor: %s", monitorID)
	}

	// Return a copy
	metrics := *monitor.Metrics
	return &metrics, nil
}

// GetAlerts returns recent risk alerts
func (re *RiskEngine) GetAlerts(limit int) []*RiskAlert {
	alerts := make([]*RiskAlert, 0, limit)

	// Drain alerts from channel up to limit
	for i := 0; i < limit; i++ {
		select {
		case alert := <-re.alerts:
			alerts = append(alerts, alert)
		default:
			break
		}
	}

	return alerts
}

// GetViolations returns recent risk violations
func (re *RiskEngine) GetViolations(limit int) []*RiskViolation {
	re.mu.RLock()
	defer re.mu.RUnlock()

	if limit <= 0 || limit > len(re.violations) {
		limit = len(re.violations)
	}

	violations := make([]*RiskViolation, limit)
	copy(violations, re.violations[len(re.violations)-limit:])

	return violations
}

// GetEngineMetrics returns risk engine performance metrics
func (re *RiskEngine) GetEngineMetrics() *RiskEngineMetrics {
	re.mu.RLock()
	defer re.mu.RUnlock()

	return &RiskEngineMetrics{
		TotalChecks:    re.totalChecks,
		ViolationCount: re.violationCount,
		AlertCount:     re.alertCount,
		ActiveMonitors: len(re.monitors),
		IsRunning:      re.isRunning,
	}
}

// RiskEngineMetrics contains risk engine performance metrics
type RiskEngineMetrics struct {
	TotalChecks    int64 `json:"total_checks"`
	ViolationCount int64 `json:"violation_count"`
	AlertCount     int64 `json:"alert_count"`
	ActiveMonitors int   `json:"active_monitors"`
	IsRunning      bool  `json:"is_running"`
}

// Private methods

// monitorRisk continuously monitors risk metrics
func (re *RiskEngine) monitorRisk(ctx context.Context) {
	defer re.wg.Done()

	ticker := time.NewTicker(re.config.CheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-re.stopChan:
			return
		case <-ticker.C:
			re.performRiskChecks(ctx)
		}
	}
}

// processAlerts processes risk alerts
func (re *RiskEngine) processAlerts(ctx context.Context) {
	defer re.wg.Done()

	for {
		select {
		case <-re.stopChan:
			return
		case alert := <-re.alerts:
			re.handleAlert(ctx, alert)
		}
	}
}

// cleanupViolations cleans up old violations
func (re *RiskEngine) cleanupViolations(ctx context.Context) {
	defer re.wg.Done()

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-re.stopChan:
			return
		case <-ticker.C:
			re.cleanupOldViolations()
		}
	}
}

// checkGlobalLimits checks signal against global risk limits
func (re *RiskEngine) checkGlobalLimits(ctx context.Context, signal *framework.Signal) error {
	if re.limits == nil {
		return nil
	}

	// Check allowed symbols
	if len(re.limits.AllowedSymbols) > 0 {
		allowed := false
		for _, allowedSymbol := range re.limits.AllowedSymbols {
			if signal.Symbol == allowedSymbol {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("symbol not in allowed list: %s", signal.Symbol)
		}
	}

	// Check blocked symbols
	for _, blockedSymbol := range re.limits.BlockedSymbols {
		if signal.Symbol == blockedSymbol {
			return fmt.Errorf("symbol is blocked: %s", signal.Symbol)
		}
	}

	// Check position size
	if signal.Quantity.GreaterThan(re.limits.MaxPositionSize) {
		return fmt.Errorf("quantity exceeds max position size: %s > %s",
			signal.Quantity.String(), re.limits.MaxPositionSize.String())
	}

	// Check trading hours
	if re.limits.TradingHours != nil && re.limits.TradingHours.Enabled {
		if !re.isWithinTradingHours(re.limits.TradingHours) {
			return fmt.Errorf("outside trading hours")
		}
	}

	return nil
}

// checkStrategyLimits checks signal against strategy-specific limits
func (re *RiskEngine) checkStrategyLimits(ctx context.Context, monitor *RiskMonitor, signal *framework.Signal) error {
	if monitor.Limits == nil {
		return nil
	}

	limits := monitor.Limits

	// Check position size
	if signal.Quantity.GreaterThan(limits.MaxPositionSize) {
		return fmt.Errorf("quantity exceeds strategy max position size")
	}

	// Check exposure limits
	if monitor.Metrics != nil {
		notionalValue := signal.Quantity.Mul(signal.Price)
		newExposure := monitor.Metrics.TotalExposure.Add(notionalValue)

		if newExposure.GreaterThan(limits.MaxExposure) {
			return fmt.Errorf("would exceed strategy max exposure")
		}
	}

	return nil
}

// checkSymbolLimits checks signal against symbol-specific limits
func (re *RiskEngine) checkSymbolLimits(ctx context.Context, monitor *RiskMonitor, signal *framework.Signal) error {
	if monitor.Limits == nil {
		return nil
	}

	// Check symbol-specific position size
	if signal.Quantity.GreaterThan(monitor.Limits.MaxPositionSize) {
		return fmt.Errorf("quantity exceeds symbol max position size")
	}

	return nil
}

// checkCircuitBreakers checks if any circuit breakers should trip
func (re *RiskEngine) checkCircuitBreakers(ctx context.Context, signal *framework.Signal) error {
	for _, cb := range re.circuitBreakers {
		if cb.State == CircuitStateOpen {
			return fmt.Errorf("circuit breaker %s is open", cb.ID)
		}
	}
	return nil
}

// handleViolation handles a risk violation
func (re *RiskEngine) handleViolation(ctx context.Context, monitorID string, entityID uuid.UUID, err error) error {
	re.violationCount++

	violation := &RiskViolation{
		ID:            uuid.New(),
		MonitorID:     monitorID,
		EntityID:      entityID,
		ViolationType: RiskViolationTypeExposure, // Default type
		Severity:      AlertSeverityError,
		Action:        ViolationActionBlock,
		Message:       err.Error(),
		Details:       make(map[string]interface{}),
		Timestamp:     time.Now(),
		Resolved:      false,
	}

	// Store violation
	re.mu.Lock()
	re.violations = append(re.violations, violation)
	re.mu.Unlock()

	// Create alert
	alert := &RiskAlert{
		ID:        uuid.New(),
		Type:      RiskAlertTypeCompliance,
		Severity:  AlertSeverityError,
		MonitorID: monitorID,
		EntityID:  entityID,
		Message:   fmt.Sprintf("Risk violation: %s", err.Error()),
		Details: map[string]interface{}{
			"violation_id": violation.ID.String(),
		},
		Timestamp: time.Now(),
	}

	// Send alert
	select {
	case re.alerts <- alert:
		re.alertCount++
	default:
		re.logger.Warn(ctx, "Alert channel full, dropping alert", map[string]interface{}{
			"alert_id": alert.ID.String(),
		})
	}

	return err
}

// performRiskChecks performs periodic risk checks on all monitors
func (re *RiskEngine) performRiskChecks(ctx context.Context) {
	re.mu.RLock()
	monitors := make([]*RiskMonitor, 0, len(re.monitors))
	for _, monitor := range re.monitors {
		if monitor.IsActive {
			monitors = append(monitors, monitor)
		}
	}
	re.mu.RUnlock()

	// Check each monitor
	for _, monitor := range monitors {
		re.checkMonitor(ctx, monitor)
	}
}

// checkMonitor performs risk checks on a specific monitor
func (re *RiskEngine) checkMonitor(ctx context.Context, monitor *RiskMonitor) {
	if monitor.Metrics == nil {
		return
	}

	metrics := monitor.Metrics
	limits := monitor.Limits

	if limits == nil {
		return
	}

	// Check exposure limits
	if !limits.MaxExposure.IsZero() && metrics.TotalExposure.GreaterThan(limits.MaxExposure) {
		re.createAlert(ctx, monitor, RiskAlertTypeExposure, AlertSeverityWarning,
			fmt.Sprintf("Exposure limit exceeded: %s > %s",
				metrics.TotalExposure.String(), limits.MaxExposure.String()))
	}

	// Check daily loss limits
	if !limits.MaxDailyLoss.IsZero() && metrics.DailyPnL.LessThan(limits.MaxDailyLoss.Neg()) {
		re.createAlert(ctx, monitor, RiskAlertTypeLoss, AlertSeverityError,
			fmt.Sprintf("Daily loss limit exceeded: %s < %s",
				metrics.DailyPnL.String(), limits.MaxDailyLoss.Neg().String()))
	}

	// Check drawdown limits
	if !limits.MaxDrawdown.IsZero() && metrics.CurrentDrawdown.GreaterThan(limits.MaxDrawdown) {
		re.createAlert(ctx, monitor, RiskAlertTypeDrawdown, AlertSeverityCritical,
			fmt.Sprintf("Drawdown limit exceeded: %s > %s",
				metrics.CurrentDrawdown.String(), limits.MaxDrawdown.String()))
	}

	// Check VaR limits
	if re.config.EnableVaRCalculation && !limits.VaRLimit.IsZero() && metrics.VaR.GreaterThan(limits.VaRLimit) {
		re.createAlert(ctx, monitor, RiskAlertTypeVaR, AlertSeverityWarning,
			fmt.Sprintf("VaR limit exceeded: %s > %s",
				metrics.VaR.String(), limits.VaRLimit.String()))
	}

	// Check leverage limits
	if !limits.MaxLeverage.IsZero() && metrics.Leverage.GreaterThan(limits.MaxLeverage) {
		re.createAlert(ctx, monitor, RiskAlertTypeLeverage, AlertSeverityWarning,
			fmt.Sprintf("Leverage limit exceeded: %s > %s",
				metrics.Leverage.String(), limits.MaxLeverage.String()))
	}

	// Check cash reserve limits
	if !limits.MinCashReserve.IsZero() && metrics.CashBalance.LessThan(limits.MinCashReserve) {
		re.createAlert(ctx, monitor, RiskAlertTypeCash, AlertSeverityWarning,
			fmt.Sprintf("Cash reserve below minimum: %s < %s",
				metrics.CashBalance.String(), limits.MinCashReserve.String()))
	}
}

// createAlert creates and sends a risk alert
func (re *RiskEngine) createAlert(ctx context.Context, monitor *RiskMonitor, alertType RiskAlertType, severity AlertSeverity, message string) {
	alert := &RiskAlert{
		ID:        uuid.New(),
		Type:      alertType,
		Severity:  severity,
		MonitorID: monitor.ID,
		EntityID:  monitor.EntityID,
		Message:   message,
		Details: map[string]interface{}{
			"monitor_type": string(monitor.Type),
		},
		Timestamp: time.Now(),
	}

	select {
	case re.alerts <- alert:
		re.alertCount++
	default:
		re.logger.Warn(ctx, "Alert channel full, dropping alert", map[string]interface{}{
			"alert_id": alert.ID.String(),
		})
	}
}

// handleAlert processes a risk alert
func (re *RiskEngine) handleAlert(ctx context.Context, alert *RiskAlert) {
	re.logger.Info(ctx, "Processing risk alert", map[string]interface{}{
		"alert_id":   alert.ID.String(),
		"alert_type": string(alert.Type),
		"severity":   string(alert.Severity),
		"message":    alert.Message,
		"monitor_id": alert.MonitorID,
		"entity_id":  alert.EntityID.String(),
	})

	// Handle different alert types
	switch alert.Type {
	case RiskAlertTypeCircuitBreaker:
		re.handleCircuitBreakerAlert(ctx, alert)
	case RiskAlertTypeDrawdown:
		if alert.Severity == AlertSeverityCritical {
			re.handleCriticalDrawdown(ctx, alert)
		}
	case RiskAlertTypeLoss:
		if alert.Severity == AlertSeverityError || alert.Severity == AlertSeverityCritical {
			re.handleExcessiveLoss(ctx, alert)
		}
	}
}

// handleCircuitBreakerAlert handles circuit breaker alerts
func (re *RiskEngine) handleCircuitBreakerAlert(ctx context.Context, alert *RiskAlert) {
	re.logger.Warn(ctx, "Circuit breaker triggered", map[string]interface{}{
		"alert_id": alert.ID.String(),
		"details":  alert.Details,
	})
}

// handleCriticalDrawdown handles critical drawdown alerts
func (re *RiskEngine) handleCriticalDrawdown(ctx context.Context, alert *RiskAlert) {
	re.logger.Error(ctx, "Critical drawdown detected", nil, map[string]interface{}{
		"alert_id":  alert.ID.String(),
		"entity_id": alert.EntityID.String(),
	})

	// Could implement automatic position liquidation here
}

// handleExcessiveLoss handles excessive loss alerts
func (re *RiskEngine) handleExcessiveLoss(ctx context.Context, alert *RiskAlert) {
	re.logger.Error(ctx, "Excessive loss detected", nil, map[string]interface{}{
		"alert_id":  alert.ID.String(),
		"entity_id": alert.EntityID.String(),
	})

	// Could implement automatic trading suspension here
}

// cleanupOldViolations removes old violations
func (re *RiskEngine) cleanupOldViolations() {
	re.mu.Lock()
	defer re.mu.Unlock()

	cutoff := time.Now().Add(-re.config.ViolationRetention)
	newViolations := make([]*RiskViolation, 0)

	for _, violation := range re.violations {
		if violation.Timestamp.After(cutoff) {
			newViolations = append(newViolations, violation)
		}
	}

	re.violations = newViolations
}

// isWithinTradingHours checks if current time is within trading hours
func (re *RiskEngine) isWithinTradingHours(hours *TradingHours) bool {
	now := time.Now()

	// Check weekday
	weekday := int(now.Weekday())
	allowed := false
	for _, allowedDay := range hours.Weekdays {
		if weekday == allowedDay {
			allowed = true
			break
		}
	}
	if !allowed {
		return false
	}

	// Check time range
	currentTime := now.Format("15:04:05")
	startTime := hours.StartTime.Format("15:04:05")
	endTime := hours.EndTime.Format("15:04:05")

	return currentTime >= startTime && currentTime <= endTime
}
