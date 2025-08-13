package hft

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/shopspring/decimal"
)

// PositionMonitor tracks positions across all symbols
type PositionMonitor struct {
	logger    *observability.Logger
	config    RiskConfig
	positions map[string]decimal.Decimal
	mu        sync.RWMutex
}

// ExposureTracker tracks total portfolio exposure
type ExposureTracker struct {
	logger           *observability.Logger
	config           RiskConfig
	totalExposure    decimal.Decimal
	exposureBySymbol map[string]decimal.Decimal
	mu               sync.RWMutex
}

// CircuitBreakerManager manages circuit breakers for market conditions
type CircuitBreakerManager struct {
	logger   *observability.Logger
	config   RiskConfig
	breakers map[string]*CircuitBreaker
	mu       sync.RWMutex
}

// CircuitBreaker represents a single circuit breaker
type CircuitBreaker struct {
	Symbol      string
	IsTripped   bool
	TripTime    time.Time
	ResetTime   time.Time
	TripCount   int
	LastPrice   decimal.Decimal
	PriceChange float64
	VolumeSpike float64
}

// RiskCalculator performs risk calculations
type RiskCalculator struct {
	logger *observability.Logger
	config RiskConfig
}

// AlertManager manages risk alerts and notifications
type AlertManager struct {
	logger *observability.Logger
	config RiskConfig
	alerts []RiskAlert
	mu     sync.RWMutex
}

// RiskAlert represents a risk alert
type RiskAlert struct {
	ID        string
	Level     Severity
	Message   string
	Timestamp time.Time
	Resolved  bool
}

// NewPositionMonitor creates a new position monitor
func NewPositionMonitor(logger *observability.Logger, config RiskConfig) *PositionMonitor {
	return &PositionMonitor{
		logger:    logger,
		config:    config,
		positions: make(map[string]decimal.Decimal),
	}
}

// NewExposureTracker creates a new exposure tracker
func NewExposureTracker(logger *observability.Logger, config RiskConfig) *ExposureTracker {
	return &ExposureTracker{
		logger:           logger,
		config:           config,
		exposureBySymbol: make(map[string]decimal.Decimal),
	}
}

// NewCircuitBreakerManager creates a new circuit breaker manager
func NewCircuitBreakerManager(logger *observability.Logger, config RiskConfig) *CircuitBreakerManager {
	return &CircuitBreakerManager{
		logger:   logger,
		config:   config,
		breakers: make(map[string]*CircuitBreaker),
	}
}

// NewRiskCalculator creates a new risk calculator
func NewRiskCalculator(logger *observability.Logger, config RiskConfig) *RiskCalculator {
	return &RiskCalculator{
		logger: logger,
		config: config,
	}
}

// NewAlertManager creates a new alert manager
func NewAlertManager(logger *observability.Logger, config RiskConfig) *AlertManager {
	return &AlertManager{
		logger: logger,
		config: config,
		alerts: make([]RiskAlert, 0),
	}
}

// GetPosition returns the current position for a symbol
func (pm *PositionMonitor) GetPosition(symbol string) decimal.Decimal {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	if position, exists := pm.positions[symbol]; exists {
		return position
	}
	return decimal.Zero
}

// UpdatePosition updates the position for a symbol
func (pm *PositionMonitor) UpdatePosition(symbol string, quantity decimal.Decimal, side OrderSide) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	currentPosition := pm.positions[symbol]

	if side == OrderSideBuy {
		pm.positions[symbol] = currentPosition.Add(quantity)
	} else {
		pm.positions[symbol] = currentPosition.Sub(quantity)
	}
}

// GetAllPositions returns all current positions
func (pm *PositionMonitor) GetAllPositions() map[string]decimal.Decimal {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	positions := make(map[string]decimal.Decimal)
	for symbol, position := range pm.positions {
		positions[symbol] = position
	}
	return positions
}

// GetTotalExposure returns the total portfolio exposure
func (et *ExposureTracker) GetTotalExposure() decimal.Decimal {
	et.mu.RLock()
	defer et.mu.RUnlock()
	return et.totalExposure
}

// UpdateExposure updates exposure for a symbol
func (et *ExposureTracker) UpdateExposure(symbol string, exposure decimal.Decimal) {
	et.mu.Lock()
	defer et.mu.Unlock()

	oldExposure := et.exposureBySymbol[symbol]
	et.exposureBySymbol[symbol] = exposure

	// Update total exposure
	et.totalExposure = et.totalExposure.Sub(oldExposure).Add(exposure)
}

// GetExposureBySymbol returns exposure for a specific symbol
func (et *ExposureTracker) GetExposureBySymbol(symbol string) decimal.Decimal {
	et.mu.RLock()
	defer et.mu.RUnlock()

	if exposure, exists := et.exposureBySymbol[symbol]; exists {
		return exposure
	}
	return decimal.Zero
}

// CheckOrder checks if an order should trigger circuit breakers
func (cbm *CircuitBreakerManager) CheckOrder(ctx context.Context, order *OrderRequest) error {
	cbm.mu.RLock()
	defer cbm.mu.RUnlock()

	if breaker, exists := cbm.breakers[order.Symbol]; exists {
		if breaker.IsTripped {
			if time.Since(breaker.TripTime) < 5*time.Minute { // 5 minute cooldown
				return fmt.Errorf("circuit breaker is active for %s", order.Symbol)
			} else {
				// Reset circuit breaker
				breaker.IsTripped = false
				breaker.ResetTime = time.Now()
			}
		}
	}

	return nil
}

// TripBreaker trips a circuit breaker for a symbol
func (cbm *CircuitBreakerManager) TripBreaker(symbol string, reason string) {
	cbm.mu.Lock()
	defer cbm.mu.Unlock()

	if breaker, exists := cbm.breakers[symbol]; exists {
		breaker.IsTripped = true
		breaker.TripTime = time.Now()
		breaker.TripCount++
	} else {
		cbm.breakers[symbol] = &CircuitBreaker{
			Symbol:    symbol,
			IsTripped: true,
			TripTime:  time.Now(),
			TripCount: 1,
		}
	}
}

// CalculateVaR calculates Value at Risk
func (rc *RiskCalculator) CalculateVaR(positions map[string]decimal.Decimal, confidence float64) decimal.Decimal {
	// Simplified VaR calculation - in production, use historical simulation or Monte Carlo
	var totalValue decimal.Decimal

	for _, position := range positions {
		// Mock price for calculation
		price := decimal.NewFromFloat(45000.0) // Would use real market prices
		value := position.Abs().Mul(price)
		totalValue = totalValue.Add(value)
	}

	// Simple VaR estimate: assume 2% daily volatility
	volatility := 0.02
	if confidence == 0.95 {
		return totalValue.Mul(decimal.NewFromFloat(volatility * 1.645)) // 95% VaR
	} else if confidence == 0.99 {
		return totalValue.Mul(decimal.NewFromFloat(volatility * 2.326)) // 99% VaR
	}

	return totalValue.Mul(decimal.NewFromFloat(volatility * 1.645))
}

// CalculateDrawdown calculates current drawdown
func (rc *RiskCalculator) CalculateDrawdown(currentValue, peakValue decimal.Decimal) float64 {
	if peakValue.IsZero() {
		return 0.0
	}

	drawdown := peakValue.Sub(currentValue).Div(peakValue)
	return drawdown.InexactFloat64() * 100.0 // Return as percentage
}

// CalculateConcentration calculates position concentration
func (rc *RiskCalculator) CalculateConcentration(positions map[string]decimal.Decimal) float64 {
	if len(positions) == 0 {
		return 0.0
	}

	var totalValue decimal.Decimal
	var maxPosition decimal.Decimal

	for _, position := range positions {
		value := position.Abs()
		totalValue = totalValue.Add(value)
		if value.GreaterThan(maxPosition) {
			maxPosition = value
		}
	}

	if totalValue.IsZero() {
		return 0.0
	}

	concentration := maxPosition.Div(totalValue)
	return concentration.InexactFloat64() * 100.0 // Return as percentage
}

// CreateAlert creates a new risk alert
func (am *AlertManager) CreateAlert(level Severity, message string) {
	am.mu.Lock()
	defer am.mu.Unlock()

	alert := RiskAlert{
		ID:        fmt.Sprintf("alert_%d", time.Now().UnixNano()),
		Level:     level,
		Message:   message,
		Timestamp: time.Now(),
		Resolved:  false,
	}

	am.alerts = append(am.alerts, alert)
}

// GetActiveAlerts returns all active (unresolved) alerts
func (am *AlertManager) GetActiveAlerts() []RiskAlert {
	am.mu.RLock()
	defer am.mu.RUnlock()

	var activeAlerts []RiskAlert
	for _, alert := range am.alerts {
		if !alert.Resolved {
			activeAlerts = append(activeAlerts, alert)
		}
	}

	return activeAlerts
}

// ResolveAlert resolves an alert by ID
func (am *AlertManager) ResolveAlert(alertID string) {
	am.mu.Lock()
	defer am.mu.Unlock()

	for i := range am.alerts {
		if am.alerts[i].ID == alertID {
			am.alerts[i].Resolved = true
			break
		}
	}
}
