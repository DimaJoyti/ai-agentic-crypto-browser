package risk

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/shopspring/decimal"
)

// VaRCalculator calculates Value at Risk using various methods
type VaRCalculator struct {
	logger *observability.Logger
	config VaRConfig
}

// VaRConfig contains VaR calculation configuration
type VaRConfig struct {
	Method           VaRMethod       `json:"method"`
	ConfidenceLevel  decimal.Decimal `json:"confidence_level"` // e.g., 0.95 for 95%
	TimeHorizon      time.Duration   `json:"time_horizon"`     // e.g., 24h
	HistoryPeriod    time.Duration   `json:"history_period"`   // e.g., 252 days
	DecayFactor      decimal.Decimal `json:"decay_factor"`     // for EWMA
	LambdaHalfLife   time.Duration   `json:"lambda_half_life"` // for EWMA
	MonteCarloSims   int             `json:"monte_carlo_sims"` // for Monte Carlo
	EnableStressTest bool            `json:"enable_stress_test"`
}

// VaRMethod represents different VaR calculation methods
type VaRMethod string

const (
	VaRMethodHistorical VaRMethod = "historical"
	VaRMethodParametric VaRMethod = "parametric"
	VaRMethodMonteCarlo VaRMethod = "monte_carlo"
	VaRMethodEWMA       VaRMethod = "ewma"
)

// VaRResult contains VaR calculation results
type VaRResult struct {
	Method            VaRMethod              `json:"method"`
	ConfidenceLevel   decimal.Decimal        `json:"confidence_level"`
	TimeHorizon       time.Duration          `json:"time_horizon"`
	VaR               decimal.Decimal        `json:"var"`
	ExpectedShortfall decimal.Decimal        `json:"expected_shortfall"` // Conditional VaR
	PortfolioValue    decimal.Decimal        `json:"portfolio_value"`
	Volatility        decimal.Decimal        `json:"volatility"`
	CalculatedAt      time.Time              `json:"calculated_at"`
	DataPoints        int                    `json:"data_points"`
	Metadata          map[string]interface{} `json:"metadata"`
}

// PortfolioData represents portfolio data for VaR calculation
type PortfolioData struct {
	Timestamp      time.Time               `json:"timestamp"`
	PortfolioValue decimal.Decimal         `json:"portfolio_value"`
	Returns        decimal.Decimal         `json:"returns"`
	Positions      map[string]PositionData `json:"positions"`
}

// PositionData represents position data
type PositionData struct {
	Symbol     string          `json:"symbol"`
	Quantity   decimal.Decimal `json:"quantity"`
	Price      decimal.Decimal `json:"price"`
	Value      decimal.Decimal `json:"value"`
	Weight     decimal.Decimal `json:"weight"`
	Volatility decimal.Decimal `json:"volatility"`
	Beta       decimal.Decimal `json:"beta"`
}

// StressTestScenario represents a stress test scenario
type StressTestScenario struct {
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	Shocks      map[string]decimal.Decimal `json:"shocks"` // symbol -> shock percentage
	MarketShock decimal.Decimal            `json:"market_shock"`
}

// StressTestResult contains stress test results
type StressTestResult struct {
	Scenario       StressTestScenario         `json:"scenario"`
	PortfolioLoss  decimal.Decimal            `json:"portfolio_loss"`
	PositionLosses map[string]decimal.Decimal `json:"position_losses"`
	CalculatedAt   time.Time                  `json:"calculated_at"`
}

// NewVaRCalculator creates a new VaR calculator
func NewVaRCalculator(logger *observability.Logger, config VaRConfig) *VaRCalculator {
	// Set defaults
	if config.ConfidenceLevel.IsZero() {
		config.ConfidenceLevel = decimal.NewFromFloat(0.95)
	}
	if config.TimeHorizon == 0 {
		config.TimeHorizon = 24 * time.Hour
	}
	if config.HistoryPeriod == 0 {
		config.HistoryPeriod = 252 * 24 * time.Hour // 252 trading days
	}
	if config.DecayFactor.IsZero() {
		config.DecayFactor = decimal.NewFromFloat(0.94)
	}
	if config.MonteCarloSims == 0 {
		config.MonteCarloSims = 10000
	}

	return &VaRCalculator{
		logger: logger,
		config: config,
	}
}

// CalculateVaR calculates Value at Risk for a portfolio
func (vc *VaRCalculator) CalculateVaR(ctx context.Context, portfolioData []*PortfolioData) (*VaRResult, error) {
	if len(portfolioData) == 0 {
		return nil, fmt.Errorf("no portfolio data provided")
	}

	// Sort data by timestamp
	sort.Slice(portfolioData, func(i, j int) bool {
		return portfolioData[i].Timestamp.Before(portfolioData[j].Timestamp)
	})

	// Calculate returns if not provided
	if err := vc.calculateReturns(portfolioData); err != nil {
		return nil, fmt.Errorf("failed to calculate returns: %w", err)
	}

	// Extract returns for calculation
	returns := make([]decimal.Decimal, len(portfolioData)-1)
	for i := 1; i < len(portfolioData); i++ {
		returns[i-1] = portfolioData[i].Returns
	}

	if len(returns) == 0 {
		return nil, fmt.Errorf("insufficient data for VaR calculation")
	}

	var result *VaRResult
	var err error

	// Calculate VaR using specified method
	switch vc.config.Method {
	case VaRMethodHistorical:
		result, err = vc.calculateHistoricalVaR(ctx, returns, portfolioData[len(portfolioData)-1].PortfolioValue)
	case VaRMethodParametric:
		result, err = vc.calculateParametricVaR(ctx, returns, portfolioData[len(portfolioData)-1].PortfolioValue)
	case VaRMethodMonteCarlo:
		result, err = vc.calculateMonteCarloVaR(ctx, returns, portfolioData[len(portfolioData)-1].PortfolioValue)
	case VaRMethodEWMA:
		result, err = vc.calculateEWMAVaR(ctx, returns, portfolioData[len(portfolioData)-1].PortfolioValue)
	default:
		return nil, fmt.Errorf("unsupported VaR method: %s", vc.config.Method)
	}

	if err != nil {
		return nil, fmt.Errorf("VaR calculation failed: %w", err)
	}

	result.CalculatedAt = time.Now()
	result.DataPoints = len(returns)

	vc.logger.Info(ctx, "VaR calculated", map[string]interface{}{
		"method":           string(result.Method),
		"confidence_level": result.ConfidenceLevel.String(),
		"var":              result.VaR.String(),
		"portfolio_value":  result.PortfolioValue.String(),
		"data_points":      result.DataPoints,
	})

	return result, nil
}

// calculateReturns calculates portfolio returns
func (vc *VaRCalculator) calculateReturns(portfolioData []*PortfolioData) error {
	for i := 1; i < len(portfolioData); i++ {
		prevValue := portfolioData[i-1].PortfolioValue
		currentValue := portfolioData[i].PortfolioValue

		if prevValue.IsZero() {
			return fmt.Errorf("zero portfolio value at index %d", i-1)
		}

		// Calculate log return
		returnValue := currentValue.Div(prevValue)
		if returnValue.LessThanOrEqual(decimal.NewFromInt(0)) {
			return fmt.Errorf("invalid return calculation at index %d", i)
		}

		// Use simple return for now (could use log return)
		portfolioData[i].Returns = currentValue.Sub(prevValue).Div(prevValue)
	}

	return nil
}

// calculateHistoricalVaR calculates VaR using historical simulation
func (vc *VaRCalculator) calculateHistoricalVaR(ctx context.Context, returns []decimal.Decimal, portfolioValue decimal.Decimal) (*VaRResult, error) {
	if len(returns) == 0 {
		return nil, fmt.Errorf("no returns data")
	}

	// Sort returns in ascending order
	sortedReturns := make([]decimal.Decimal, len(returns))
	copy(sortedReturns, returns)
	sort.Slice(sortedReturns, func(i, j int) bool {
		return sortedReturns[i].LessThan(sortedReturns[j])
	})

	// Calculate percentile index
	alpha := decimal.NewFromInt(1).Sub(vc.config.ConfidenceLevel)
	percentileIndex := alpha.Mul(decimal.NewFromInt(int64(len(sortedReturns))))
	index := int(percentileIndex.IntPart())

	if index >= len(sortedReturns) {
		index = len(sortedReturns) - 1
	}

	// VaR is the negative of the percentile return
	varReturn := sortedReturns[index].Neg()
	var_ := portfolioValue.Mul(varReturn)

	// Calculate Expected Shortfall (Conditional VaR)
	var tailSum decimal.Decimal
	tailCount := 0
	for i := 0; i <= index; i++ {
		tailSum = tailSum.Add(sortedReturns[i])
		tailCount++
	}

	var expectedShortfall decimal.Decimal
	if tailCount > 0 {
		avgTailReturn := tailSum.Div(decimal.NewFromInt(int64(tailCount)))
		expectedShortfall = portfolioValue.Mul(avgTailReturn.Neg())
	}

	// Calculate volatility
	volatility := vc.calculateVolatility(returns)

	return &VaRResult{
		Method:            VaRMethodHistorical,
		ConfidenceLevel:   vc.config.ConfidenceLevel,
		TimeHorizon:       vc.config.TimeHorizon,
		VaR:               var_,
		ExpectedShortfall: expectedShortfall,
		PortfolioValue:    portfolioValue,
		Volatility:        volatility,
		Metadata: map[string]interface{}{
			"percentile_index": index,
			"sorted_returns":   len(sortedReturns),
		},
	}, nil
}

// calculateParametricVaR calculates VaR using parametric method (normal distribution)
func (vc *VaRCalculator) calculateParametricVaR(ctx context.Context, returns []decimal.Decimal, portfolioValue decimal.Decimal) (*VaRResult, error) {
	if len(returns) == 0 {
		return nil, fmt.Errorf("no returns data")
	}

	// Calculate mean and standard deviation
	mean := vc.calculateMean(returns)
	stdDev := vc.calculateStandardDeviation(returns, mean)

	// Calculate z-score for confidence level
	alpha := decimal.NewFromInt(1).Sub(vc.config.ConfidenceLevel)
	zScore := vc.calculateZScore(alpha)

	// VaR = -(mean + z * stdDev) * portfolio_value
	varReturn := mean.Add(zScore.Mul(stdDev)).Neg()
	var_ := portfolioValue.Mul(varReturn)

	// Expected Shortfall for normal distribution
	// ES = -(mean - stdDev * φ(z) / α) * portfolio_value
	phi := vc.calculateStandardNormalPDF(zScore)
	expectedShortfall := portfolioValue.Mul(
		mean.Sub(stdDev.Mul(phi).Div(alpha)).Neg(),
	)

	return &VaRResult{
		Method:            VaRMethodParametric,
		ConfidenceLevel:   vc.config.ConfidenceLevel,
		TimeHorizon:       vc.config.TimeHorizon,
		VaR:               var_,
		ExpectedShortfall: expectedShortfall,
		PortfolioValue:    portfolioValue,
		Volatility:        stdDev,
		Metadata: map[string]interface{}{
			"mean":    mean.String(),
			"std_dev": stdDev.String(),
			"z_score": zScore.String(),
		},
	}, nil
}

// Helper methods for statistical calculations

// calculateMean calculates the mean of returns
func (vc *VaRCalculator) calculateMean(returns []decimal.Decimal) decimal.Decimal {
	if len(returns) == 0 {
		return decimal.NewFromInt(0)
	}

	sum := decimal.NewFromInt(0)
	for _, ret := range returns {
		sum = sum.Add(ret)
	}

	return sum.Div(decimal.NewFromInt(int64(len(returns))))
}

// calculateStandardDeviation calculates the standard deviation of returns
func (vc *VaRCalculator) calculateStandardDeviation(returns []decimal.Decimal, mean decimal.Decimal) decimal.Decimal {
	if len(returns) <= 1 {
		return decimal.NewFromInt(0)
	}

	sumSquaredDiffs := decimal.NewFromInt(0)
	for _, ret := range returns {
		diff := ret.Sub(mean)
		sumSquaredDiffs = sumSquaredDiffs.Add(diff.Mul(diff))
	}

	variance := sumSquaredDiffs.Div(decimal.NewFromInt(int64(len(returns) - 1)))
	return decimal.NewFromFloat(math.Sqrt(variance.InexactFloat64()))
}

// calculateVolatility calculates volatility (same as standard deviation for now)
func (vc *VaRCalculator) calculateVolatility(returns []decimal.Decimal) decimal.Decimal {
	mean := vc.calculateMean(returns)
	return vc.calculateStandardDeviation(returns, mean)
}

// calculateZScore calculates the z-score for a given alpha (for normal distribution)
func (vc *VaRCalculator) calculateZScore(alpha decimal.Decimal) decimal.Decimal {
	// Simplified z-score calculation for common confidence levels
	// In practice, you would use inverse normal CDF
	alphaFloat := alpha.InexactFloat64()

	switch {
	case alphaFloat <= 0.01: // 99% confidence
		return decimal.NewFromFloat(-2.326)
	case alphaFloat <= 0.025: // 97.5% confidence
		return decimal.NewFromFloat(-1.96)
	case alphaFloat <= 0.05: // 95% confidence
		return decimal.NewFromFloat(-1.645)
	case alphaFloat <= 0.1: // 90% confidence
		return decimal.NewFromFloat(-1.282)
	default:
		return decimal.NewFromFloat(-1.645) // Default to 95%
	}
}

// calculateStandardNormalPDF calculates the standard normal probability density function
func (vc *VaRCalculator) calculateStandardNormalPDF(z decimal.Decimal) decimal.Decimal {
	zFloat := z.InexactFloat64()
	pdf := (1.0 / math.Sqrt(2*math.Pi)) * math.Exp(-0.5*zFloat*zFloat)
	return decimal.NewFromFloat(pdf)
}

// generateNormalRandom generates a random number from normal distribution
func (vc *VaRCalculator) generateNormalRandom(mean, stdDev decimal.Decimal) decimal.Decimal {
	// Simplified random generation - in practice use proper random number generator
	// This is a placeholder implementation
	return mean.Add(stdDev.Mul(decimal.NewFromFloat(0.0))) // Returns mean for now
}

// RunStressTest runs stress tests on the portfolio
func (vc *VaRCalculator) RunStressTest(ctx context.Context, portfolioData *PortfolioData, scenarios []StressTestScenario) ([]*StressTestResult, error) {
	if !vc.config.EnableStressTest {
		return nil, fmt.Errorf("stress testing is disabled")
	}

	results := make([]*StressTestResult, len(scenarios))

	for i, scenario := range scenarios {
		result, err := vc.runSingleStressTest(ctx, portfolioData, scenario)
		if err != nil {
			vc.logger.Error(ctx, "Stress test failed", err, map[string]interface{}{
				"scenario": scenario.Name,
			})
			continue
		}
		results[i] = result
	}

	return results, nil
}

// runSingleStressTest runs a single stress test scenario
func (vc *VaRCalculator) runSingleStressTest(ctx context.Context, portfolioData *PortfolioData, scenario StressTestScenario) (*StressTestResult, error) {
	positionLosses := make(map[string]decimal.Decimal)
	totalLoss := decimal.NewFromInt(0)

	// Apply shocks to each position
	for symbol, position := range portfolioData.Positions {
		var shock decimal.Decimal

		// Check for symbol-specific shock
		if symbolShock, exists := scenario.Shocks[symbol]; exists {
			shock = symbolShock
		} else {
			// Use market shock as default
			shock = scenario.MarketShock
		}

		// Calculate position loss
		positionLoss := position.Value.Mul(shock)
		positionLosses[symbol] = positionLoss
		totalLoss = totalLoss.Add(positionLoss)
	}

	return &StressTestResult{
		Scenario:       scenario,
		PortfolioLoss:  totalLoss,
		PositionLosses: positionLosses,
		CalculatedAt:   time.Now(),
	}, nil
}

// GetDefaultStressScenarios returns default stress test scenarios
func (vc *VaRCalculator) GetDefaultStressScenarios() []StressTestScenario {
	return []StressTestScenario{
		{
			Name:        "Market Crash",
			Description: "Severe market downturn similar to 2008 financial crisis",
			Shocks:      make(map[string]decimal.Decimal),
			MarketShock: decimal.NewFromFloat(-0.30), // 30% decline
		},
		{
			Name:        "Flash Crash",
			Description: "Sudden market drop similar to May 2010 flash crash",
			Shocks:      make(map[string]decimal.Decimal),
			MarketShock: decimal.NewFromFloat(-0.10), // 10% decline
		},
		{
			Name:        "Crypto Winter",
			Description: "Severe cryptocurrency market decline",
			Shocks: map[string]decimal.Decimal{
				"BTCUSDT": decimal.NewFromFloat(-0.50), // 50% decline
				"ETHUSDT": decimal.NewFromFloat(-0.60), // 60% decline
				"ADAUSDT": decimal.NewFromFloat(-0.70), // 70% decline
			},
			MarketShock: decimal.NewFromFloat(-0.40), // 40% decline for others
		},
		{
			Name:        "Interest Rate Shock",
			Description: "Sudden interest rate increase affecting all assets",
			Shocks:      make(map[string]decimal.Decimal),
			MarketShock: decimal.NewFromFloat(-0.15), // 15% decline
		},
		{
			Name:        "Liquidity Crisis",
			Description: "Market liquidity dries up causing price gaps",
			Shocks:      make(map[string]decimal.Decimal),
			MarketShock: decimal.NewFromFloat(-0.25), // 25% decline
		},
	}
}

// calculateMonteCarloVaR calculates VaR using Monte Carlo simulation
func (vc *VaRCalculator) calculateMonteCarloVaR(ctx context.Context, returns []decimal.Decimal, portfolioValue decimal.Decimal) (*VaRResult, error) {
	if len(returns) == 0 {
		return nil, fmt.Errorf("no returns data")
	}

	// Calculate historical mean and standard deviation
	mean := vc.calculateMean(returns)
	stdDev := vc.calculateStandardDeviation(returns, mean)

	// Generate Monte Carlo simulations
	simulatedReturns := make([]decimal.Decimal, vc.config.MonteCarloSims)
	for i := 0; i < vc.config.MonteCarloSims; i++ {
		// Generate random normal return
		randomReturn := vc.generateNormalRandom(mean, stdDev)
		simulatedReturns[i] = randomReturn
	}

	// Calculate VaR from simulated returns
	return vc.calculateHistoricalVaR(ctx, simulatedReturns, portfolioValue)
}

// calculateEWMAVaR calculates VaR using Exponentially Weighted Moving Average
func (vc *VaRCalculator) calculateEWMAVaR(ctx context.Context, returns []decimal.Decimal, portfolioValue decimal.Decimal) (*VaRResult, error) {
	if len(returns) == 0 {
		return nil, fmt.Errorf("no returns data")
	}

	// Calculate EWMA variance
	lambda := vc.config.DecayFactor
	variance := decimal.NewFromInt(0)

	// Initialize with first return squared
	if len(returns) > 0 {
		variance = returns[0].Mul(returns[0])
	}

	// Calculate EWMA variance
	for i := 1; i < len(returns); i++ {
		returnSquared := returns[i].Mul(returns[i])
		variance = lambda.Mul(variance).Add(
			decimal.NewFromInt(1).Sub(lambda).Mul(returnSquared),
		)
	}

	// Standard deviation is square root of variance
	stdDev := decimal.NewFromFloat(math.Sqrt(variance.InexactFloat64()))

	// Calculate z-score for confidence level
	alpha := decimal.NewFromInt(1).Sub(vc.config.ConfidenceLevel)
	zScore := vc.calculateZScore(alpha)

	// VaR = -z * stdDev * portfolio_value
	var_ := portfolioValue.Mul(zScore.Mul(stdDev).Neg())

	return &VaRResult{
		Method:          VaRMethodEWMA,
		ConfidenceLevel: vc.config.ConfidenceLevel,
		TimeHorizon:     vc.config.TimeHorizon,
		VaR:             var_,
		PortfolioValue:  portfolioValue,
		Volatility:      stdDev,
		Metadata: map[string]interface{}{
			"lambda":   lambda.String(),
			"variance": variance.String(),
		},
	}, nil
}
