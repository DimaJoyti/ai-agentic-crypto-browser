package analytics

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/shopspring/decimal"
)

// PortfolioPerformanceAnalyzer analyzes portfolio performance metrics
type PortfolioPerformanceAnalyzer struct {
	logger    *observability.Logger
	config    PerformanceConfig
	metrics   PortfolioPerformanceMetrics
	history   []PortfolioSnapshot
	positions map[string]*Position
	mu        sync.RWMutex
	isRunning int32
}

// PortfolioSnapshot represents a point-in-time portfolio snapshot
type PortfolioSnapshot struct {
	Timestamp   time.Time       `json:"timestamp"`
	TotalValue  decimal.Decimal `json:"total_value"`
	Cash        decimal.Decimal `json:"cash"`
	Positions   []Position      `json:"positions"`
	DailyReturn decimal.Decimal `json:"daily_return"`
	TotalReturn decimal.Decimal `json:"total_return"`
	Benchmark   decimal.Decimal `json:"benchmark,omitempty"`
}

// Position represents a portfolio position
type Position struct {
	Symbol        string          `json:"symbol"`
	Quantity      decimal.Decimal `json:"quantity"`
	AveragePrice  decimal.Decimal `json:"average_price"`
	CurrentPrice  decimal.Decimal `json:"current_price"`
	MarketValue   decimal.Decimal `json:"market_value"`
	UnrealizedPnL decimal.Decimal `json:"unrealized_pnl"`
	Weight        float64         `json:"weight"`
	LastUpdated   time.Time       `json:"last_updated"`
}

// RiskMetrics represents portfolio risk metrics
type RiskMetrics struct {
	VaR95             decimal.Decimal `json:"var_95"`
	VaR99             decimal.Decimal `json:"var_99"`
	ExpectedShortfall decimal.Decimal `json:"expected_shortfall"`
	Beta              float64         `json:"beta"`
	Alpha             float64         `json:"alpha"`
	TrackingError     float64         `json:"tracking_error"`
	InformationRatio  float64         `json:"information_ratio"`
	MaxDrawdown       decimal.Decimal `json:"max_drawdown"`
	CurrentDrawdown   decimal.Decimal `json:"current_drawdown"`
	VolatilityAnnual  float64         `json:"volatility_annual"`
	SharpeRatio       float64         `json:"sharpe_ratio"`
	SortinoRatio      float64         `json:"sortino_ratio"`
	CalmarRatio       float64         `json:"calmar_ratio"`
}

// NewPortfolioPerformanceAnalyzer creates a new portfolio performance analyzer
func NewPortfolioPerformanceAnalyzer(logger *observability.Logger, config PerformanceConfig) *PortfolioPerformanceAnalyzer {
	return &PortfolioPerformanceAnalyzer{
		logger:    logger,
		config:    config,
		history:   make([]PortfolioSnapshot, 0, config.MetricsBufferSize),
		positions: make(map[string]*Position),
	}
}

// Start starts the portfolio performance analyzer
func (ppa *PortfolioPerformanceAnalyzer) Start(ctx context.Context) error {
	ppa.logger.Info(ctx, "Starting portfolio performance analyzer", nil)
	ppa.isRunning = 1
	return nil
}

// Stop stops the portfolio performance analyzer
func (ppa *PortfolioPerformanceAnalyzer) Stop(ctx context.Context) error {
	ppa.logger.Info(ctx, "Stopping portfolio performance analyzer", nil)
	ppa.isRunning = 0
	return nil
}

// UpdatePosition updates a position in the portfolio
func (ppa *PortfolioPerformanceAnalyzer) UpdatePosition(position Position) {
	ppa.mu.Lock()
	defer ppa.mu.Unlock()

	position.LastUpdated = time.Now()
	ppa.positions[position.Symbol] = &position

	// Recalculate metrics
	ppa.calculateMetrics()
}

// AddSnapshot adds a portfolio snapshot for analysis
func (ppa *PortfolioPerformanceAnalyzer) AddSnapshot(snapshot PortfolioSnapshot) {
	ppa.mu.Lock()
	defer ppa.mu.Unlock()

	// Add snapshot to history
	ppa.history = append(ppa.history, snapshot)

	// Maintain buffer size
	if len(ppa.history) > ppa.config.MetricsBufferSize {
		ppa.history = ppa.history[1:]
	}

	// Recalculate metrics
	ppa.calculateMetrics()
}

// calculateMetrics calculates comprehensive portfolio performance metrics
func (ppa *PortfolioPerformanceAnalyzer) calculateMetrics() {
	if len(ppa.history) == 0 {
		return
	}

	latest := ppa.history[len(ppa.history)-1]

	// Calculate basic metrics
	totalValue := latest.TotalValue
	totalReturn := ppa.calculateTotalReturn()
	dailyReturn := ppa.calculateDailyReturn()
	weeklyReturn := ppa.calculatePeriodReturn(7 * 24 * time.Hour)
	monthlyReturn := ppa.calculatePeriodReturn(30 * 24 * time.Hour)
	yearlyReturn := ppa.calculatePeriodReturn(365 * 24 * time.Hour)
	annualizedReturn := ppa.calculateAnnualizedReturn()

	// Calculate risk metrics
	returns := ppa.extractReturns()
	annualizedVolatility := ppa.calculateVolatility(returns)
	sharpeRatio := ppa.calculateSharpeRatio(returns, annualizedVolatility)
	sortinoRatio := ppa.calculateSortinoRatio(returns)
	maxDrawdown := ppa.calculateMaxDrawdown()
	currentDrawdown := ppa.calculateCurrentDrawdown()
	drawdownDuration := ppa.calculateDrawdownDuration()
	recoveryTime := ppa.calculateRecoveryTime()

	// Calculate VaR and Expected Shortfall
	var95 := ppa.calculateVaR(returns, 0.95)
	var99 := ppa.calculateVaR(returns, 0.99)
	expectedShortfall := ppa.calculateExpectedShortfall(returns, 0.95)

	// Calculate benchmark-relative metrics
	beta := ppa.calculateBeta(returns)
	alpha := ppa.calculateAlpha(returns, beta)
	trackingError := ppa.calculateTrackingError(returns)
	informationRatio := ppa.calculateInformationRatio(returns, trackingError)
	calmarRatio := ppa.calculateCalmarRatio(annualizedReturn, maxDrawdown)

	// Calculate additional risk metrics
	ulcerIndex := ppa.calculateUlcerIndex()
	sterlingRatio := ppa.calculateSterlingRatio(annualizedReturn, maxDrawdown)
	burkeRatio := ppa.calculateBurkeRatio(returns)

	// Update metrics
	ppa.metrics = PortfolioPerformanceMetrics{
		TotalValue:           totalValue,
		TotalReturn:          totalReturn,
		DailyReturn:          dailyReturn,
		WeeklyReturn:         weeklyReturn,
		MonthlyReturn:        monthlyReturn,
		YearlyReturn:         yearlyReturn,
		AnnualizedReturn:     annualizedReturn,
		AnnualizedVolatility: annualizedVolatility,
		SharpeRatio:          sharpeRatio,
		SortinoRatio:         sortinoRatio,
		MaxDrawdown:          maxDrawdown,
		CurrentDrawdown:      currentDrawdown,
		DrawdownDuration:     drawdownDuration,
		RecoveryTime:         recoveryTime,
		VaR95:                var95,
		VaR99:                var99,
		ExpectedShortfall:    expectedShortfall,
		Beta:                 beta,
		Alpha:                alpha,
		TrackingError:        trackingError,
		InformationRatio:     informationRatio,
		CalmarRatio:          calmarRatio,
		UlcerIndex:           ulcerIndex,
		SterlingRatio:        sterlingRatio,
		BurkeRatio:           burkeRatio,
	}
}

// calculateTotalReturn calculates total return since inception
func (ppa *PortfolioPerformanceAnalyzer) calculateTotalReturn() decimal.Decimal {
	if len(ppa.history) < 2 {
		return decimal.Zero
	}

	initial := ppa.history[0].TotalValue
	current := ppa.history[len(ppa.history)-1].TotalValue

	if initial.IsZero() {
		return decimal.Zero
	}

	return current.Sub(initial).Div(initial)
}

// calculateDailyReturn calculates daily return
func (ppa *PortfolioPerformanceAnalyzer) calculateDailyReturn() decimal.Decimal {
	if len(ppa.history) < 2 {
		return decimal.Zero
	}

	current := ppa.history[len(ppa.history)-1].TotalValue
	previous := ppa.history[len(ppa.history)-2].TotalValue

	if previous.IsZero() {
		return decimal.Zero
	}

	return current.Sub(previous).Div(previous)
}

// calculatePeriodReturn calculates return for a specific period
func (ppa *PortfolioPerformanceAnalyzer) calculatePeriodReturn(period time.Duration) decimal.Decimal {
	if len(ppa.history) == 0 {
		return decimal.Zero
	}

	current := ppa.history[len(ppa.history)-1]
	cutoff := current.Timestamp.Add(-period)

	// Find snapshot closest to cutoff time
	var baseline PortfolioSnapshot
	found := false
	for i := len(ppa.history) - 1; i >= 0; i-- {
		if ppa.history[i].Timestamp.Before(cutoff) {
			baseline = ppa.history[i]
			found = true
			break
		}
	}

	if !found || baseline.TotalValue.IsZero() {
		return decimal.Zero
	}

	return current.TotalValue.Sub(baseline.TotalValue).Div(baseline.TotalValue)
}

// calculateAnnualizedReturn calculates annualized return
func (ppa *PortfolioPerformanceAnalyzer) calculateAnnualizedReturn() decimal.Decimal {
	if len(ppa.history) < 2 {
		return decimal.Zero
	}

	first := ppa.history[0]
	last := ppa.history[len(ppa.history)-1]

	daysDiff := last.Timestamp.Sub(first.Timestamp).Hours() / 24
	if daysDiff <= 0 || first.TotalValue.IsZero() {
		return decimal.Zero
	}

	totalReturn := last.TotalValue.Sub(first.TotalValue).Div(first.TotalValue)
	totalReturnFloat, _ := totalReturn.Float64()

	annualizedReturn := math.Pow(1+totalReturnFloat, 365/daysDiff) - 1
	return decimal.NewFromFloat(annualizedReturn)
}

// extractReturns extracts daily returns from history
func (ppa *PortfolioPerformanceAnalyzer) extractReturns() []float64 {
	if len(ppa.history) < 2 {
		return []float64{}
	}

	returns := make([]float64, 0, len(ppa.history)-1)
	for i := 1; i < len(ppa.history); i++ {
		prev := ppa.history[i-1].TotalValue
		curr := ppa.history[i].TotalValue

		if !prev.IsZero() {
			ret, _ := curr.Sub(prev).Div(prev).Float64()
			returns = append(returns, ret)
		}
	}

	return returns
}

// calculateVolatility calculates annualized volatility
func (ppa *PortfolioPerformanceAnalyzer) calculateVolatility(returns []float64) float64 {
	if len(returns) < 2 {
		return 0
	}

	// Calculate mean
	var sum float64
	for _, ret := range returns {
		sum += ret
	}
	mean := sum / float64(len(returns))

	// Calculate variance
	var variance float64
	for _, ret := range returns {
		variance += math.Pow(ret-mean, 2)
	}
	variance /= float64(len(returns) - 1)

	// Annualize (assuming daily returns)
	return math.Sqrt(variance) * math.Sqrt(252)
}

// calculateSharpeRatio calculates Sharpe ratio
func (ppa *PortfolioPerformanceAnalyzer) calculateSharpeRatio(returns []float64, volatility float64) float64 {
	if len(returns) == 0 || volatility == 0 {
		return 0
	}

	var sum float64
	for _, ret := range returns {
		sum += ret
	}
	meanReturn := sum / float64(len(returns))

	// Assume risk-free rate of 2% annually
	riskFreeRate := 0.02 / 252 // Daily risk-free rate

	return (meanReturn*252 - riskFreeRate*252) / volatility
}

// calculateSortinoRatio calculates Sortino ratio
func (ppa *PortfolioPerformanceAnalyzer) calculateSortinoRatio(returns []float64) float64 {
	if len(returns) == 0 {
		return 0
	}

	var sum float64
	for _, ret := range returns {
		sum += ret
	}
	meanReturn := sum / float64(len(returns))

	// Calculate downside deviation
	var downsideVariance float64
	downsideCount := 0
	for _, ret := range returns {
		if ret < 0 {
			downsideVariance += ret * ret
			downsideCount++
		}
	}

	if downsideCount == 0 {
		return math.Inf(1)
	}

	downsideVariance /= float64(downsideCount)
	downsideStdDev := math.Sqrt(downsideVariance) * math.Sqrt(252)

	if downsideStdDev == 0 {
		return 0
	}

	riskFreeRate := 0.02
	return (meanReturn*252 - riskFreeRate) / downsideStdDev
}

// calculateMaxDrawdown calculates maximum drawdown
func (ppa *PortfolioPerformanceAnalyzer) calculateMaxDrawdown() decimal.Decimal {
	if len(ppa.history) == 0 {
		return decimal.Zero
	}

	var peak, maxDrawdown decimal.Decimal
	peak = ppa.history[0].TotalValue

	for _, snapshot := range ppa.history {
		if snapshot.TotalValue.GreaterThan(peak) {
			peak = snapshot.TotalValue
		}

		drawdown := peak.Sub(snapshot.TotalValue).Div(peak)
		if drawdown.GreaterThan(maxDrawdown) {
			maxDrawdown = drawdown
		}
	}

	return maxDrawdown
}

// calculateCurrentDrawdown calculates current drawdown
func (ppa *PortfolioPerformanceAnalyzer) calculateCurrentDrawdown() decimal.Decimal {
	if len(ppa.history) == 0 {
		return decimal.Zero
	}

	var peak decimal.Decimal
	current := ppa.history[len(ppa.history)-1].TotalValue

	for _, snapshot := range ppa.history {
		if snapshot.TotalValue.GreaterThan(peak) {
			peak = snapshot.TotalValue
		}
	}

	if peak.IsZero() {
		return decimal.Zero
	}

	return peak.Sub(current).Div(peak)
}

// calculateDrawdownDuration calculates current drawdown duration
func (ppa *PortfolioPerformanceAnalyzer) calculateDrawdownDuration() time.Duration {
	if len(ppa.history) < 2 {
		return 0
	}

	current := ppa.history[len(ppa.history)-1]
	var peak decimal.Decimal
	var peakTime time.Time

	// Find the most recent peak
	for i := len(ppa.history) - 1; i >= 0; i-- {
		if ppa.history[i].TotalValue.GreaterThan(peak) {
			peak = ppa.history[i].TotalValue
			peakTime = ppa.history[i].Timestamp
		}
	}

	if current.TotalValue.GreaterThanOrEqual(peak) {
		return 0 // No current drawdown
	}

	return current.Timestamp.Sub(peakTime)
}

// calculateRecoveryTime calculates average recovery time from drawdowns
func (ppa *PortfolioPerformanceAnalyzer) calculateRecoveryTime() time.Duration {
	// Simplified implementation - would need more sophisticated drawdown analysis
	return 30 * 24 * time.Hour // Mock: 30 days average recovery
}

// calculateVaR calculates Value at Risk
func (ppa *PortfolioPerformanceAnalyzer) calculateVaR(returns []float64, confidence float64) decimal.Decimal {
	if len(returns) == 0 {
		return decimal.Zero
	}

	// Sort returns
	sortedReturns := make([]float64, len(returns))
	copy(sortedReturns, returns)

	// Simple bubble sort for small arrays
	for i := 0; i < len(sortedReturns); i++ {
		for j := i + 1; j < len(sortedReturns); j++ {
			if sortedReturns[i] > sortedReturns[j] {
				sortedReturns[i], sortedReturns[j] = sortedReturns[j], sortedReturns[i]
			}
		}
	}

	// Find percentile
	index := int(float64(len(sortedReturns)) * (1 - confidence))
	if index >= len(sortedReturns) {
		index = len(sortedReturns) - 1
	}

	var95Return := sortedReturns[index]

	// Convert to portfolio value terms
	if len(ppa.history) > 0 {
		currentValue := ppa.history[len(ppa.history)-1].TotalValue
		return currentValue.Mul(decimal.NewFromFloat(-var95Return))
	}

	return decimal.Zero
}

// calculateExpectedShortfall calculates Expected Shortfall (Conditional VaR)
func (ppa *PortfolioPerformanceAnalyzer) calculateExpectedShortfall(returns []float64, confidence float64) decimal.Decimal {
	if len(returns) == 0 {
		return decimal.Zero
	}

	// Sort returns
	sortedReturns := make([]float64, len(returns))
	copy(sortedReturns, returns)

	for i := 0; i < len(sortedReturns); i++ {
		for j := i + 1; j < len(sortedReturns); j++ {
			if sortedReturns[i] > sortedReturns[j] {
				sortedReturns[i], sortedReturns[j] = sortedReturns[j], sortedReturns[i]
			}
		}
	}

	// Calculate average of worst returns
	cutoff := int(float64(len(sortedReturns)) * (1 - confidence))
	if cutoff == 0 {
		cutoff = 1
	}

	var sum float64
	for i := 0; i < cutoff; i++ {
		sum += sortedReturns[i]
	}
	avgWorstReturn := sum / float64(cutoff)

	// Convert to portfolio value terms
	if len(ppa.history) > 0 {
		currentValue := ppa.history[len(ppa.history)-1].TotalValue
		return currentValue.Mul(decimal.NewFromFloat(-avgWorstReturn))
	}

	return decimal.Zero
}

// calculateBeta calculates portfolio beta (simplified)
func (ppa *PortfolioPerformanceAnalyzer) calculateBeta(returns []float64) float64 {
	// Simplified beta calculation - would need benchmark returns
	// For now, return a mock value
	return 1.2
}

// calculateAlpha calculates portfolio alpha
func (ppa *PortfolioPerformanceAnalyzer) calculateAlpha(returns []float64, beta float64) float64 {
	// Simplified alpha calculation - would need benchmark returns
	// For now, return a mock value
	return 0.03 // 3% annual alpha
}

// calculateTrackingError calculates tracking error vs benchmark
func (ppa *PortfolioPerformanceAnalyzer) calculateTrackingError(returns []float64) float64 {
	// Simplified tracking error - would need benchmark returns
	// For now, return a mock value
	return 0.05 // 5% tracking error
}

// calculateInformationRatio calculates information ratio
func (ppa *PortfolioPerformanceAnalyzer) calculateInformationRatio(returns []float64, trackingError float64) float64 {
	if trackingError == 0 {
		return 0
	}

	// Simplified calculation - would need actual excess returns
	excessReturn := 0.02 // Mock 2% excess return
	return excessReturn / trackingError
}

// calculateCalmarRatio calculates Calmar ratio
func (ppa *PortfolioPerformanceAnalyzer) calculateCalmarRatio(annualizedReturn decimal.Decimal, maxDrawdown decimal.Decimal) float64 {
	if maxDrawdown.IsZero() {
		return math.Inf(1)
	}

	annualizedReturnFloat, _ := annualizedReturn.Float64()
	maxDrawdownFloat, _ := maxDrawdown.Float64()

	return annualizedReturnFloat / maxDrawdownFloat
}

// calculateUlcerIndex calculates Ulcer Index
func (ppa *PortfolioPerformanceAnalyzer) calculateUlcerIndex() float64 {
	// Simplified Ulcer Index calculation
	return 5.0 // Mock value
}

// calculateSterlingRatio calculates Sterling ratio
func (ppa *PortfolioPerformanceAnalyzer) calculateSterlingRatio(annualizedReturn decimal.Decimal, maxDrawdown decimal.Decimal) float64 {
	if maxDrawdown.IsZero() {
		return math.Inf(1)
	}

	annualizedReturnFloat, _ := annualizedReturn.Float64()
	maxDrawdownFloat, _ := maxDrawdown.Float64()

	// Sterling ratio uses average drawdown, but we'll use max drawdown for simplification
	return annualizedReturnFloat / (maxDrawdownFloat + 0.1) // Add 10% penalty
}

// calculateBurkeRatio calculates Burke ratio
func (ppa *PortfolioPerformanceAnalyzer) calculateBurkeRatio(returns []float64) float64 {
	// Simplified Burke ratio calculation
	return 1.5 // Mock value
}

// GetMetrics returns current portfolio performance metrics
func (ppa *PortfolioPerformanceAnalyzer) GetMetrics() PortfolioPerformanceMetrics {
	ppa.mu.RLock()
	defer ppa.mu.RUnlock()
	return ppa.metrics
}

// GetPositions returns current portfolio positions
func (ppa *PortfolioPerformanceAnalyzer) GetPositions() map[string]*Position {
	ppa.mu.RLock()
	defer ppa.mu.RUnlock()

	positions := make(map[string]*Position)
	for k, v := range ppa.positions {
		positions[k] = v
	}
	return positions
}

// GetRiskMetrics returns comprehensive risk metrics
func (ppa *PortfolioPerformanceAnalyzer) GetRiskMetrics() RiskMetrics {
	ppa.mu.RLock()
	defer ppa.mu.RUnlock()

	return RiskMetrics{
		VaR95:             ppa.metrics.VaR95,
		VaR99:             ppa.metrics.VaR99,
		ExpectedShortfall: ppa.metrics.ExpectedShortfall,
		Beta:              ppa.metrics.Beta,
		Alpha:             ppa.metrics.Alpha,
		TrackingError:     ppa.metrics.TrackingError,
		InformationRatio:  ppa.metrics.InformationRatio,
		MaxDrawdown:       ppa.metrics.MaxDrawdown,
		CurrentDrawdown:   ppa.metrics.CurrentDrawdown,
		VolatilityAnnual:  ppa.metrics.AnnualizedVolatility,
		SharpeRatio:       ppa.metrics.SharpeRatio,
		SortinoRatio:      ppa.metrics.SortinoRatio,
		CalmarRatio:       ppa.metrics.CalmarRatio,
	}
}
