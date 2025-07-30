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

// RiskManager provides advanced risk management for HFT
type RiskManager struct {
	logger     *observability.Logger
	config     HFTConfig
	riskLimits map[string]*RiskLimit
	violations []RiskViolation

	// Risk metrics
	dailyPnL      decimal.Decimal
	maxDrawdown   decimal.Decimal
	positionSizes map[string]decimal.Decimal
	orderCounts   map[string]int64

	// Circuit breakers
	tradingHalted int32
	emergencyStop int32
	lastResetTime time.Time

	// State management
	isRunning int32
	stopChan  chan struct{}
	wg        sync.WaitGroup
	mu        sync.RWMutex
}

// RiskLimit defines risk limits for trading
type RiskLimit struct {
	ID              uuid.UUID       `json:"id"`
	Name            string          `json:"name"`
	Type            RiskLimitType   `json:"type"`
	Symbol          string          `json:"symbol,omitempty"`
	MaxPositionSize decimal.Decimal `json:"max_position_size,omitempty"`
	MaxDailyLoss    decimal.Decimal `json:"max_daily_loss,omitempty"`
	MaxOrderSize    decimal.Decimal `json:"max_order_size,omitempty"`
	MaxOrdersPerMin int64           `json:"max_orders_per_min,omitempty"`
	MaxDrawdown     decimal.Decimal `json:"max_drawdown,omitempty"`
	VaRLimit        decimal.Decimal `json:"var_limit,omitempty"`
	Enabled         bool            `json:"enabled"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// RiskLimitType represents different types of risk limits
type RiskLimitType string

const (
	RiskLimitTypePosition      RiskLimitType = "POSITION"
	RiskLimitTypeDailyLoss     RiskLimitType = "DAILY_LOSS"
	RiskLimitTypeOrderSize     RiskLimitType = "ORDER_SIZE"
	RiskLimitTypeOrderRate     RiskLimitType = "ORDER_RATE"
	RiskLimitTypeDrawdown      RiskLimitType = "DRAWDOWN"
	RiskLimitTypeVaR           RiskLimitType = "VAR"
	RiskLimitTypeConcentration RiskLimitType = "CONCENTRATION"
)

// RiskViolation represents a risk limit violation
type RiskViolation struct {
	ID            uuid.UUID         `json:"id"`
	LimitID       uuid.UUID         `json:"limit_id"`
	LimitName     string            `json:"limit_name"`
	Symbol        string            `json:"symbol,omitempty"`
	ViolationType RiskLimitType     `json:"violation_type"`
	CurrentValue  decimal.Decimal   `json:"current_value"`
	LimitValue    decimal.Decimal   `json:"limit_value"`
	Severity      ViolationSeverity `json:"severity"`
	Action        ViolationAction   `json:"action"`
	Message       string            `json:"message"`
	Timestamp     time.Time         `json:"timestamp"`
	Resolved      bool              `json:"resolved"`
	ResolvedAt    *time.Time        `json:"resolved_at,omitempty"`
}

// ViolationSeverity represents the severity of a risk violation
type ViolationSeverity string

const (
	ViolationSeverityLow      ViolationSeverity = "LOW"
	ViolationSeverityMedium   ViolationSeverity = "MEDIUM"
	ViolationSeverityHigh     ViolationSeverity = "HIGH"
	ViolationSeverityCritical ViolationSeverity = "CRITICAL"
)

// ViolationAction represents the action taken for a violation
type ViolationAction string

const (
	ViolationActionWarning       ViolationAction = "WARNING"
	ViolationActionRejectOrder   ViolationAction = "REJECT_ORDER"
	ViolationActionHaltTrading   ViolationAction = "HALT_TRADING"
	ViolationActionEmergencyStop ViolationAction = "EMERGENCY_STOP"
)

// NewRiskManager creates a new risk manager
func NewRiskManager(logger *observability.Logger, config HFTConfig) *RiskManager {
	rm := &RiskManager{
		logger:        logger,
		config:        config,
		riskLimits:    make(map[string]*RiskLimit),
		violations:    make([]RiskViolation, 0),
		positionSizes: make(map[string]decimal.Decimal),
		orderCounts:   make(map[string]int64),
		lastResetTime: time.Now(),
		stopChan:      make(chan struct{}),
	}

	// Initialize default risk limits
	rm.initializeDefaultLimits()

	return rm
}

// Start begins the risk manager
func (rm *RiskManager) Start(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&rm.isRunning, 0, 1) {
		return fmt.Errorf("risk manager is already running")
	}

	rm.logger.Info(ctx, "Starting risk manager", map[string]interface{}{
		"risk_limits":       len(rm.riskLimits),
		"max_daily_loss":    rm.config.MaxDailyLoss.String(),
		"max_position_size": rm.config.MaxPositionSize.String(),
	})

	// Start monitoring goroutines
	rm.wg.Add(2)
	go rm.monitorRiskLimits(ctx)
	go rm.resetDailyCounters(ctx)

	return nil
}

// Stop gracefully shuts down the risk manager
func (rm *RiskManager) Stop(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&rm.isRunning, 1, 0) {
		return fmt.Errorf("risk manager is not running")
	}

	rm.logger.Info(ctx, "Stopping risk manager", nil)

	close(rm.stopChan)
	rm.wg.Wait()

	rm.logger.Info(ctx, "Risk manager stopped", map[string]interface{}{
		"total_violations": len(rm.violations),
		"trading_halted":   atomic.LoadInt32(&rm.tradingHalted) == 1,
	})

	return nil
}

// ValidateSignal validates a trading signal against risk limits
func (rm *RiskManager) ValidateSignal(signal TradingSignal) error {
	if atomic.LoadInt32(&rm.emergencyStop) == 1 {
		return fmt.Errorf("emergency stop activated - all trading halted")
	}

	if atomic.LoadInt32(&rm.tradingHalted) == 1 {
		return fmt.Errorf("trading halted due to risk violation")
	}

	rm.mu.RLock()
	defer rm.mu.RUnlock()

	// Check position size limit
	if err := rm.checkPositionSizeLimit(signal); err != nil {
		return err
	}

	// Check order size limit
	if err := rm.checkOrderSizeLimit(signal); err != nil {
		return err
	}

	// Check order rate limit
	if err := rm.checkOrderRateLimit(signal); err != nil {
		return err
	}

	// Check daily loss limit
	if err := rm.checkDailyLossLimit(); err != nil {
		return err
	}

	// Check concentration limits
	if err := rm.checkConcentrationLimit(signal); err != nil {
		return err
	}

	return nil
}

// CheckLimits checks if current market conditions violate risk limits
func (rm *RiskManager) CheckLimits(symbol string, price decimal.Decimal) bool {
	if atomic.LoadInt32(&rm.emergencyStop) == 1 {
		return false
	}

	if atomic.LoadInt32(&rm.tradingHalted) == 1 {
		return false
	}

	// Check drawdown limits
	if rm.checkDrawdownLimit() {
		return false
	}

	return true
}

// UpdatePosition updates position information for risk monitoring
func (rm *RiskManager) UpdatePosition(update PositionUpdate) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rm.positionSizes[update.Symbol] = update.Size

	// Check if position update triggers any violations
	rm.checkPositionViolations(update)
}

// AddRiskLimit adds a new risk limit
func (rm *RiskManager) AddRiskLimit(limit *RiskLimit) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if limit.ID == uuid.Nil {
		limit.ID = uuid.New()
	}

	limit.CreatedAt = time.Now()
	limit.UpdatedAt = time.Now()

	rm.riskLimits[limit.ID.String()] = limit

	rm.logger.Info(context.Background(), "Risk limit added", map[string]interface{}{
		"limit_id":   limit.ID.String(),
		"limit_name": limit.Name,
		"limit_type": string(limit.Type),
	})
}

// RemoveRiskLimit removes a risk limit
func (rm *RiskManager) RemoveRiskLimit(limitID uuid.UUID) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if _, exists := rm.riskLimits[limitID.String()]; !exists {
		return fmt.Errorf("risk limit not found: %s", limitID.String())
	}

	delete(rm.riskLimits, limitID.String())

	rm.logger.Info(context.Background(), "Risk limit removed", map[string]interface{}{
		"limit_id": limitID.String(),
	})

	return nil
}

// HaltTrading halts all trading due to risk violation
func (rm *RiskManager) HaltTrading(reason string) {
	atomic.StoreInt32(&rm.tradingHalted, 1)

	rm.logger.Error(context.Background(), "Trading halted", fmt.Errorf(reason), map[string]interface{}{
		"reason":    reason,
		"timestamp": time.Now(),
	})
}

// ResumeTrading resumes trading after halt
func (rm *RiskManager) ResumeTrading() {
	atomic.StoreInt32(&rm.tradingHalted, 0)

	rm.logger.Info(context.Background(), "Trading resumed", map[string]interface{}{
		"timestamp": time.Now(),
	})
}

// EmergencyStop triggers emergency stop
func (rm *RiskManager) EmergencyStop(reason string) {
	atomic.StoreInt32(&rm.emergencyStop, 1)
	atomic.StoreInt32(&rm.tradingHalted, 1)

	rm.logger.Error(context.Background(), "Emergency stop activated", fmt.Errorf(reason), map[string]interface{}{
		"reason":    reason,
		"timestamp": time.Now(),
	})
}

// initializeDefaultLimits sets up default risk limits
func (rm *RiskManager) initializeDefaultLimits() {
	// Daily loss limit
	rm.AddRiskLimit(&RiskLimit{
		Name:         "Daily Loss Limit",
		Type:         RiskLimitTypeDailyLoss,
		MaxDailyLoss: rm.config.MaxDailyLoss,
		Enabled:      true,
	})

	// Position size limit
	rm.AddRiskLimit(&RiskLimit{
		Name:            "Position Size Limit",
		Type:            RiskLimitTypePosition,
		MaxPositionSize: rm.config.MaxPositionSize,
		Enabled:         true,
	})

	// Order rate limit
	rm.AddRiskLimit(&RiskLimit{
		Name:            "Order Rate Limit",
		Type:            RiskLimitTypeOrderRate,
		MaxOrdersPerMin: int64(rm.config.MaxOrdersPerSecond * 60),
		Enabled:         true,
	})
}

// checkPositionSizeLimit checks position size against limits
func (rm *RiskManager) checkPositionSizeLimit(signal TradingSignal) error {
	currentPosition := rm.positionSizes[signal.Symbol]
	newPosition := currentPosition

	if signal.Side == OrderSideBuy {
		newPosition = currentPosition.Add(signal.Quantity)
	} else {
		newPosition = currentPosition.Sub(signal.Quantity)
	}

	for _, limit := range rm.riskLimits {
		if !limit.Enabled || limit.Type != RiskLimitTypePosition {
			continue
		}

		if limit.Symbol != "" && limit.Symbol != signal.Symbol {
			continue
		}

		if newPosition.Abs().GreaterThan(limit.MaxPositionSize) {
			return fmt.Errorf("position size limit exceeded: %s > %s",
				newPosition.Abs().String(), limit.MaxPositionSize.String())
		}
	}

	return nil
}

// checkOrderSizeLimit checks order size against limits
func (rm *RiskManager) checkOrderSizeLimit(signal TradingSignal) error {
	for _, limit := range rm.riskLimits {
		if !limit.Enabled || limit.Type != RiskLimitTypeOrderSize {
			continue
		}

		if limit.Symbol != "" && limit.Symbol != signal.Symbol {
			continue
		}

		if signal.Quantity.GreaterThan(limit.MaxOrderSize) {
			return fmt.Errorf("order size limit exceeded: %s > %s",
				signal.Quantity.String(), limit.MaxOrderSize.String())
		}
	}

	return nil
}

// checkOrderRateLimit checks order rate against limits
func (rm *RiskManager) checkOrderRateLimit(signal TradingSignal) error {
	now := time.Now()
	minuteKey := fmt.Sprintf("%s_%d", signal.Symbol, now.Unix()/60)

	currentCount := rm.orderCounts[minuteKey]

	for _, limit := range rm.riskLimits {
		if !limit.Enabled || limit.Type != RiskLimitTypeOrderRate {
			continue
		}

		if limit.Symbol != "" && limit.Symbol != signal.Symbol {
			continue
		}

		if currentCount >= limit.MaxOrdersPerMin {
			return fmt.Errorf("order rate limit exceeded: %d >= %d orders per minute",
				currentCount, limit.MaxOrdersPerMin)
		}
	}

	// Increment counter
	rm.orderCounts[minuteKey] = currentCount + 1

	return nil
}

// checkDailyLossLimit checks daily loss against limits
func (rm *RiskManager) checkDailyLossLimit() error {
	if rm.dailyPnL.LessThan(decimal.Zero) && rm.dailyPnL.Abs().GreaterThan(rm.config.MaxDailyLoss) {
		return fmt.Errorf("daily loss limit exceeded: %s > %s",
			rm.dailyPnL.Abs().String(), rm.config.MaxDailyLoss.String())
	}

	return nil
}

// checkConcentrationLimit checks position concentration
func (rm *RiskManager) checkConcentrationLimit(signal TradingSignal) error {
	// Simple concentration check - ensure no single position exceeds 20% of portfolio
	totalValue := decimal.Zero
	for _, size := range rm.positionSizes {
		totalValue = totalValue.Add(size.Abs())
	}

	if totalValue.IsZero() {
		return nil
	}

	currentPosition := rm.positionSizes[signal.Symbol]
	newPosition := currentPosition

	if signal.Side == OrderSideBuy {
		newPosition = currentPosition.Add(signal.Quantity)
	} else {
		newPosition = currentPosition.Sub(signal.Quantity)
	}

	concentration := newPosition.Abs().Div(totalValue)
	maxConcentration := decimal.NewFromFloat(0.2) // 20%

	if concentration.GreaterThan(maxConcentration) {
		return fmt.Errorf("concentration limit exceeded: %s > %s",
			concentration.String(), maxConcentration.String())
	}

	return nil
}

// checkDrawdownLimit checks drawdown against limits
func (rm *RiskManager) checkDrawdownLimit() bool {
	// This would be implemented with actual portfolio data
	// For now, return false (no violation)
	return false
}

// checkPositionViolations checks for position-related violations
func (rm *RiskManager) checkPositionViolations(update PositionUpdate) {
	// Check for violations and create violation records
	// This is a simplified implementation
}

// monitorRiskLimits continuously monitors risk limits
func (rm *RiskManager) monitorRiskLimits(ctx context.Context) {
	defer rm.wg.Done()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-rm.stopChan:
			return
		case <-ticker.C:
			rm.performRiskChecks(ctx)
		}
	}
}

// resetDailyCounters resets daily counters at midnight
func (rm *RiskManager) resetDailyCounters(ctx context.Context) {
	defer rm.wg.Done()

	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-rm.stopChan:
			return
		case now := <-ticker.C:
			if now.Day() != rm.lastResetTime.Day() {
				rm.mu.Lock()
				rm.dailyPnL = decimal.Zero
				rm.orderCounts = make(map[string]int64)
				rm.lastResetTime = now
				rm.mu.Unlock()

				rm.logger.Info(ctx, "Daily risk counters reset", map[string]interface{}{
					"date": now.Format("2006-01-02"),
				})
			}
		}
	}
}

// performRiskChecks performs periodic risk checks
func (rm *RiskManager) performRiskChecks(ctx context.Context) {
	// Implement periodic risk monitoring
	// Check for limit violations, update metrics, etc.
}

// GetViolations returns all risk violations
func (rm *RiskManager) GetViolations() []RiskViolation {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	return rm.violations
}

// GetRiskLimits returns all risk limits
func (rm *RiskManager) GetRiskLimits() map[string]*RiskLimit {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	limits := make(map[string]*RiskLimit)
	for id, limit := range rm.riskLimits {
		limits[id] = limit
	}

	return limits
}

// IsTrading returns whether trading is currently allowed
func (rm *RiskManager) IsTrading() bool {
	return atomic.LoadInt32(&rm.tradingHalted) == 0 && atomic.LoadInt32(&rm.emergencyStop) == 0
}

// GetMetrics returns risk manager metrics
func (rm *RiskManager) GetMetrics() RiskMetrics {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	return RiskMetrics{
		VaR95:             decimal.NewFromFloat(-1000.0), // Mock VaR calculation
		VaR99:             decimal.NewFromFloat(-2000.0), // Mock VaR calculation
		ExpectedShortfall: decimal.NewFromFloat(-2500.0), // Mock ES calculation
		SharpeRatio:       1.5,
		SortinoRatio:      2.0,
		MaxDrawdown:       rm.maxDrawdown,
		Beta:              0.8,
		Alpha:             0.1,
		Volatility:        0.2,
		Correlation:       map[string]float64{"BTC": 0.7, "ETH": 0.6},
	}
}
