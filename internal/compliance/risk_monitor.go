package compliance

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// RiskMonitor provides real-time risk monitoring and alerting
type RiskMonitor struct {
	logger      *observability.Logger
	config      ComplianceConfig
	riskMetrics *RiskMetrics
	alerts      map[string]*RiskAlert
	thresholds  map[string]*RiskThreshold
	positions   map[string]*Position
	mu          sync.RWMutex
	isRunning   int32
	stopChan    chan struct{}
}

// RiskMetrics contains current risk metrics
type RiskMetrics struct {
	TotalExposure     decimal.Decimal `json:"total_exposure"`
	MaxDrawdown       decimal.Decimal `json:"max_drawdown"`
	CurrentDrawdown   decimal.Decimal `json:"current_drawdown"`
	VaR95             decimal.Decimal `json:"var_95"`
	VaR99             decimal.Decimal `json:"var_99"`
	PortfolioValue    decimal.Decimal `json:"portfolio_value"`
	LeverageRatio     float64         `json:"leverage_ratio"`
	ConcentrationRisk float64         `json:"concentration_risk"`
	CorrelationRisk   float64         `json:"correlation_risk"`
	LiquidityRisk     float64         `json:"liquidity_risk"`
	VolatilityRisk    float64         `json:"volatility_risk"`
	LastUpdated       time.Time       `json:"last_updated"`
}

// RiskAlert represents a risk alert
type RiskAlert struct {
	ID             uuid.UUID       `json:"id"`
	Type           RiskAlertType   `json:"type"`
	Severity       AlertSeverity   `json:"severity"`
	Title          string          `json:"title"`
	Description    string          `json:"description"`
	Symbol         string          `json:"symbol,omitempty"`
	Value          decimal.Decimal `json:"value"`
	Threshold      decimal.Decimal `json:"threshold"`
	Timestamp      time.Time       `json:"timestamp"`
	Acknowledged   bool            `json:"acknowledged"`
	AcknowledgedBy string          `json:"acknowledged_by,omitempty"`
	AcknowledgedAt *time.Time      `json:"acknowledged_at,omitempty"`
	Resolved       bool            `json:"resolved"`
	ResolvedAt     *time.Time      `json:"resolved_at,omitempty"`
}

// RiskAlertType defines types of risk alerts
type RiskAlertType string

const (
	RiskAlertTypePositionLimit RiskAlertType = "POSITION_LIMIT"
	RiskAlertTypeDailyLoss     RiskAlertType = "DAILY_LOSS"
	RiskAlertTypeDrawdown      RiskAlertType = "DRAWDOWN"
	RiskAlertTypeVaR           RiskAlertType = "VAR"
	RiskAlertTypeLeverage      RiskAlertType = "LEVERAGE"
	RiskAlertTypeConcentration RiskAlertType = "CONCENTRATION"
	RiskAlertTypeCorrelation   RiskAlertType = "CORRELATION"
	RiskAlertTypeLiquidity     RiskAlertType = "LIQUIDITY"
	RiskAlertTypeVolatility    RiskAlertType = "VOLATILITY"
	RiskAlertTypeOrderSize     RiskAlertType = "ORDER_SIZE"
	RiskAlertTypeOrderRate     RiskAlertType = "ORDER_RATE"
)

// AlertSeverity defines alert severity levels
type AlertSeverity string

const (
	AlertSeverityInfo     AlertSeverity = "INFO"
	AlertSeverityWarning  AlertSeverity = "WARNING"
	AlertSeverityError    AlertSeverity = "ERROR"
	AlertSeverityCritical AlertSeverity = "CRITICAL"
)

// RiskThreshold defines risk thresholds
type RiskThreshold struct {
	ID            uuid.UUID       `json:"id"`
	Type          RiskAlertType   `json:"type"`
	Symbol        string          `json:"symbol,omitempty"`
	WarningLevel  decimal.Decimal `json:"warning_level"`
	ErrorLevel    decimal.Decimal `json:"error_level"`
	CriticalLevel decimal.Decimal `json:"critical_level"`
	Enabled       bool            `json:"enabled"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// Position represents a trading position for risk calculation
type Position struct {
	Symbol        string          `json:"symbol"`
	Size          decimal.Decimal `json:"size"`
	EntryPrice    decimal.Decimal `json:"entry_price"`
	CurrentPrice  decimal.Decimal `json:"current_price"`
	UnrealizedPnL decimal.Decimal `json:"unrealized_pnl"`
	MarketValue   decimal.Decimal `json:"market_value"`
	LastUpdated   time.Time       `json:"last_updated"`
}

// NewRiskMonitor creates a new risk monitor
func NewRiskMonitor(logger *observability.Logger, config ComplianceConfig) *RiskMonitor {
	rm := &RiskMonitor{
		logger:     logger,
		config:     config,
		alerts:     make(map[string]*RiskAlert),
		thresholds: make(map[string]*RiskThreshold),
		positions:  make(map[string]*Position),
		stopChan:   make(chan struct{}),
		riskMetrics: &RiskMetrics{
			LastUpdated: time.Now(),
		},
	}

	// Initialize default thresholds
	rm.initializeDefaultThresholds()

	return rm
}

// Start starts the risk monitor
func (rm *RiskMonitor) Start(ctx context.Context) error {
	rm.logger.Info(ctx, "Starting risk monitor", nil)
	rm.isRunning = 1

	// Start monitoring goroutine
	go rm.monitoringLoop(ctx)

	return nil
}

// Stop stops the risk monitor
func (rm *RiskMonitor) Stop(ctx context.Context) error {
	rm.logger.Info(ctx, "Stopping risk monitor", nil)
	rm.isRunning = 0
	close(rm.stopChan)
	return nil
}

// monitoringLoop runs the main monitoring loop
func (rm *RiskMonitor) monitoringLoop(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second) // Monitor every second
	defer ticker.Stop()

	for {
		select {
		case <-rm.stopChan:
			return
		case <-ticker.C:
			rm.updateRiskMetrics(ctx)
			rm.checkThresholds(ctx)
		}
	}
}

// UpdatePosition updates a position for risk monitoring
func (rm *RiskMonitor) UpdatePosition(ctx context.Context, position *Position) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	position.LastUpdated = time.Now()
	rm.positions[position.Symbol] = position

	rm.logger.Debug(ctx, "Position updated for risk monitoring", map[string]interface{}{
		"symbol":         position.Symbol,
		"size":           position.Size,
		"current_price":  position.CurrentPrice,
		"unrealized_pnl": position.UnrealizedPnL,
	})

	return nil
}

// updateRiskMetrics calculates and updates current risk metrics
func (rm *RiskMonitor) updateRiskMetrics(ctx context.Context) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	metrics := &RiskMetrics{
		LastUpdated: time.Now(),
	}

	// Calculate total exposure and portfolio value
	var totalExposure, portfolioValue decimal.Decimal
	for _, position := range rm.positions {
		exposure := position.Size.Abs().Mul(position.CurrentPrice)
		totalExposure = totalExposure.Add(exposure)
		portfolioValue = portfolioValue.Add(position.MarketValue)
	}

	metrics.TotalExposure = totalExposure
	metrics.PortfolioValue = portfolioValue

	// Calculate leverage ratio
	if !portfolioValue.IsZero() {
		leverageFloat, _ := totalExposure.Div(portfolioValue).Float64()
		metrics.LeverageRatio = leverageFloat
	}

	// Calculate concentration risk (largest position as % of portfolio)
	if !portfolioValue.IsZero() {
		var maxPositionValue decimal.Decimal
		for _, position := range rm.positions {
			positionValue := position.MarketValue.Abs()
			if positionValue.GreaterThan(maxPositionValue) {
				maxPositionValue = positionValue
			}
		}
		concentrationFloat, _ := maxPositionValue.Div(portfolioValue).Float64()
		metrics.ConcentrationRisk = concentrationFloat
	}

	// Calculate VaR (simplified Monte Carlo simulation)
	metrics.VaR95 = rm.calculateVaR(0.95)
	metrics.VaR99 = rm.calculateVaR(0.99)

	// Calculate other risk metrics
	metrics.CorrelationRisk = rm.calculateCorrelationRisk()
	metrics.LiquidityRisk = rm.calculateLiquidityRisk()
	metrics.VolatilityRisk = rm.calculateVolatilityRisk()

	rm.riskMetrics = metrics
}

// calculateVaR calculates Value at Risk using simplified method
func (rm *RiskMonitor) calculateVaR(confidence float64) decimal.Decimal {
	if len(rm.positions) == 0 {
		return decimal.Zero
	}

	// Simplified VaR calculation - in production, use historical simulation or Monte Carlo
	var totalValue decimal.Decimal
	var weightedVolatility float64

	for _, position := range rm.positions {
		value := position.MarketValue.Abs()
		totalValue = totalValue.Add(value)

		// Assume 2% daily volatility for simplification
		volatility := 0.02
		weight, _ := value.Div(totalValue).Float64()
		weightedVolatility += weight * volatility
	}

	// Z-score for confidence level
	var zScore float64
	switch confidence {
	case 0.95:
		zScore = 1.645
	case 0.99:
		zScore = 2.326
	default:
		zScore = 1.645
	}

	varValue := totalValue.Mul(decimal.NewFromFloat(weightedVolatility * zScore))
	return varValue
}

// calculateCorrelationRisk calculates correlation risk between positions
func (rm *RiskMonitor) calculateCorrelationRisk() float64 {
	// Simplified correlation risk - assume higher correlation for similar assets
	if len(rm.positions) < 2 {
		return 0.0
	}

	// In production, calculate actual correlations using historical data
	// For now, return a mock value based on position diversity
	return 0.3 // 30% correlation risk
}

// calculateLiquidityRisk calculates liquidity risk
func (rm *RiskMonitor) calculateLiquidityRisk() float64 {
	// Simplified liquidity risk based on position sizes
	var totalValue, largePositionValue decimal.Decimal

	for _, position := range rm.positions {
		value := position.MarketValue.Abs()
		totalValue = totalValue.Add(value)

		// Consider positions > $100k as potentially illiquid
		if value.GreaterThan(decimal.NewFromInt(100000)) {
			largePositionValue = largePositionValue.Add(value)
		}
	}

	if totalValue.IsZero() {
		return 0.0
	}

	liquidityRisk, _ := largePositionValue.Div(totalValue).Float64()
	return liquidityRisk
}

// calculateVolatilityRisk calculates volatility risk
func (rm *RiskMonitor) calculateVolatilityRisk() float64 {
	// Simplified volatility risk calculation
	// In production, use actual volatility calculations
	return 0.15 // 15% volatility risk
}

// checkThresholds checks all risk thresholds and generates alerts
func (rm *RiskMonitor) checkThresholds(ctx context.Context) {
	for _, threshold := range rm.thresholds {
		if !threshold.Enabled {
			continue
		}

		rm.checkThreshold(ctx, threshold)
	}
}

// checkThreshold checks a specific threshold
func (rm *RiskMonitor) checkThreshold(ctx context.Context, threshold *RiskThreshold) {
	var currentValue decimal.Decimal
	var description string

	switch threshold.Type {
	case RiskAlertTypePositionLimit:
		if threshold.Symbol != "" {
			if position, exists := rm.positions[threshold.Symbol]; exists {
				currentValue = position.MarketValue.Abs()
				description = fmt.Sprintf("Position size for %s", threshold.Symbol)
			}
		}
	case RiskAlertTypeDailyLoss:
		// Calculate daily P&L
		var dailyPnL decimal.Decimal
		for _, position := range rm.positions {
			dailyPnL = dailyPnL.Add(position.UnrealizedPnL)
		}
		if dailyPnL.IsNegative() {
			currentValue = dailyPnL.Abs()
			description = "Daily loss limit"
		}
	case RiskAlertTypeDrawdown:
		currentValue = rm.riskMetrics.CurrentDrawdown
		description = "Portfolio drawdown"
	case RiskAlertTypeVaR:
		currentValue = rm.riskMetrics.VaR95
		description = "Value at Risk (95%)"
	case RiskAlertTypeLeverage:
		currentValue = decimal.NewFromFloat(rm.riskMetrics.LeverageRatio)
		description = "Leverage ratio"
	case RiskAlertTypeConcentration:
		currentValue = decimal.NewFromFloat(rm.riskMetrics.ConcentrationRisk)
		description = "Concentration risk"
	}

	// Check threshold levels
	var severity AlertSeverity
	var exceeded bool

	if currentValue.GreaterThanOrEqual(threshold.CriticalLevel) {
		severity = AlertSeverityCritical
		exceeded = true
	} else if currentValue.GreaterThanOrEqual(threshold.ErrorLevel) {
		severity = AlertSeverityError
		exceeded = true
	} else if currentValue.GreaterThanOrEqual(threshold.WarningLevel) {
		severity = AlertSeverityWarning
		exceeded = true
	}

	if exceeded {
		rm.generateAlert(ctx, threshold, currentValue, severity, description)
	}
}

// generateAlert generates a risk alert
func (rm *RiskMonitor) generateAlert(ctx context.Context, threshold *RiskThreshold, currentValue decimal.Decimal, severity AlertSeverity, description string) {
	alert := &RiskAlert{
		ID:          uuid.New(),
		Type:        threshold.Type,
		Severity:    severity,
		Title:       fmt.Sprintf("%s threshold exceeded", threshold.Type),
		Description: description,
		Symbol:      threshold.Symbol,
		Value:       currentValue,
		Threshold:   rm.getThresholdForSeverity(threshold, severity),
		Timestamp:   time.Now(),
	}

	rm.mu.Lock()
	rm.alerts[alert.ID.String()] = alert
	rm.mu.Unlock()

	rm.logger.Warn(ctx, "Risk alert generated", map[string]interface{}{
		"alert_id":  alert.ID,
		"type":      alert.Type,
		"severity":  alert.Severity,
		"symbol":    alert.Symbol,
		"value":     alert.Value,
		"threshold": alert.Threshold,
	})
}

// getThresholdForSeverity returns the threshold value for a given severity
func (rm *RiskMonitor) getThresholdForSeverity(threshold *RiskThreshold, severity AlertSeverity) decimal.Decimal {
	switch severity {
	case AlertSeverityCritical:
		return threshold.CriticalLevel
	case AlertSeverityError:
		return threshold.ErrorLevel
	case AlertSeverityWarning:
		return threshold.WarningLevel
	default:
		return threshold.WarningLevel
	}
}

// initializeDefaultThresholds sets up default risk thresholds
func (rm *RiskMonitor) initializeDefaultThresholds() {
	thresholds := []*RiskThreshold{
		{
			ID:            uuid.New(),
			Type:          RiskAlertTypeDailyLoss,
			WarningLevel:  decimal.NewFromInt(10000), // $10k
			ErrorLevel:    decimal.NewFromInt(25000), // $25k
			CriticalLevel: decimal.NewFromInt(50000), // $50k
			Enabled:       true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			ID:            uuid.New(),
			Type:          RiskAlertTypeLeverage,
			WarningLevel:  decimal.NewFromFloat(2.0),  // 2:1
			ErrorLevel:    decimal.NewFromFloat(5.0),  // 5:1
			CriticalLevel: decimal.NewFromFloat(10.0), // 10:1
			Enabled:       true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			ID:            uuid.New(),
			Type:          RiskAlertTypeConcentration,
			WarningLevel:  decimal.NewFromFloat(0.3), // 30%
			ErrorLevel:    decimal.NewFromFloat(0.5), // 50%
			CriticalLevel: decimal.NewFromFloat(0.7), // 70%
			Enabled:       true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			ID:            uuid.New(),
			Type:          RiskAlertTypeVaR,
			WarningLevel:  decimal.NewFromInt(50000),  // $50k
			ErrorLevel:    decimal.NewFromInt(100000), // $100k
			CriticalLevel: decimal.NewFromInt(200000), // $200k
			Enabled:       true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	for _, threshold := range thresholds {
		rm.thresholds[threshold.ID.String()] = threshold
	}
}

// GetRiskMetrics returns current risk metrics
func (rm *RiskMonitor) GetRiskMetrics() *RiskMetrics {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.riskMetrics
}

// GetAlerts returns current alerts with optional filtering
func (rm *RiskMonitor) GetAlerts(filter RiskAlertFilter) []*RiskAlert {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	var alerts []*RiskAlert
	for _, alert := range rm.alerts {
		if rm.matchesAlertFilter(alert, filter) {
			alerts = append(alerts, alert)
		}
	}

	return alerts
}

// RiskAlertFilter defines filtering criteria for risk alerts
type RiskAlertFilter struct {
	Type         RiskAlertType `json:"type,omitempty"`
	Severity     AlertSeverity `json:"severity,omitempty"`
	Symbol       string        `json:"symbol,omitempty"`
	Acknowledged *bool         `json:"acknowledged,omitempty"`
	Resolved     *bool         `json:"resolved,omitempty"`
	Limit        int           `json:"limit,omitempty"`
}

// matchesAlertFilter checks if an alert matches the filter criteria
func (rm *RiskMonitor) matchesAlertFilter(alert *RiskAlert, filter RiskAlertFilter) bool {
	if filter.Type != "" && alert.Type != filter.Type {
		return false
	}
	if filter.Severity != "" && alert.Severity != filter.Severity {
		return false
	}
	if filter.Symbol != "" && alert.Symbol != filter.Symbol {
		return false
	}
	if filter.Acknowledged != nil && alert.Acknowledged != *filter.Acknowledged {
		return false
	}
	if filter.Resolved != nil && alert.Resolved != *filter.Resolved {
		return false
	}
	return true
}
