package ml

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// PortfolioOptimizer implements modern portfolio theory and optimization algorithms
type PortfolioOptimizer struct {
	logger *observability.Logger
	config OptimizerConfig
}

// OptimizerConfig contains portfolio optimizer configuration
type OptimizerConfig struct {
	RiskFreeRate       decimal.Decimal `json:"risk_free_rate"`       // Risk-free rate for Sharpe ratio
	MaxIterations      int             `json:"max_iterations"`       // Maximum optimization iterations
	Tolerance          decimal.Decimal `json:"tolerance"`            // Convergence tolerance
	MinWeight          decimal.Decimal `json:"min_weight"`           // Minimum asset weight
	MaxWeight          decimal.Decimal `json:"max_weight"`           // Maximum asset weight
	EnableShortSelling bool            `json:"enable_short_selling"` // Allow negative weights
	RebalanceThreshold decimal.Decimal `json:"rebalance_threshold"`  // Threshold for rebalancing
	LookbackPeriod     time.Duration   `json:"lookback_period"`      // Historical data lookback
}

// OptimizationObjective represents the optimization objective
type OptimizationObjective string

const (
	ObjectiveMaxSharpe      OptimizationObjective = "max_sharpe"
	ObjectiveMinVariance    OptimizationObjective = "min_variance"
	ObjectiveMaxReturn      OptimizationObjective = "max_return"
	ObjectiveRiskParity     OptimizationObjective = "risk_parity"
	ObjectiveBlackLitterman OptimizationObjective = "black_litterman"
	ObjectiveCustom         OptimizationObjective = "custom"
)

// OptimizationRequest represents a portfolio optimization request
type OptimizationRequest struct {
	ID               uuid.UUID                `json:"id"`
	Objective        OptimizationObjective    `json:"objective"`
	Assets           []string                 `json:"assets"`
	Returns          [][]decimal.Decimal      `json:"returns"`           // Historical returns matrix
	ExpectedReturns  []decimal.Decimal        `json:"expected_returns"`  // Expected returns vector
	CovarianceMatrix [][]decimal.Decimal      `json:"covariance_matrix"` // Covariance matrix
	CurrentWeights   []decimal.Decimal        `json:"current_weights"`   // Current portfolio weights
	Constraints      *OptimizationConstraints `json:"constraints"`
	RiskAversion     decimal.Decimal          `json:"risk_aversion"` // Risk aversion parameter
	TargetReturn     decimal.Decimal          `json:"target_return"` // Target return for efficient frontier
	TargetRisk       decimal.Decimal          `json:"target_risk"`   // Target risk level
	Views            []*MarketView            `json:"views"`         // Black-Litterman views
	Timestamp        time.Time                `json:"timestamp"`
}

// OptimizationConstraints defines portfolio constraints
type OptimizationConstraints struct {
	MinWeights       []decimal.Decimal          `json:"min_weights"`       // Minimum weights per asset
	MaxWeights       []decimal.Decimal          `json:"max_weights"`       // Maximum weights per asset
	GroupConstraints []*GroupConstraint         `json:"group_constraints"` // Asset group constraints
	TurnoverLimit    decimal.Decimal            `json:"turnover_limit"`    // Maximum portfolio turnover
	TransactionCosts []decimal.Decimal          `json:"transaction_costs"` // Transaction costs per asset
	LiquidityLimits  []decimal.Decimal          `json:"liquidity_limits"`  // Liquidity constraints
	SectorLimits     map[string]decimal.Decimal `json:"sector_limits"`     // Sector exposure limits
	ESGConstraints   *ESGConstraints            `json:"esg_constraints"`   // ESG constraints
}

// GroupConstraint defines constraints on asset groups
type GroupConstraint struct {
	Name      string          `json:"name"`
	Assets    []string        `json:"assets"`
	MinWeight decimal.Decimal `json:"min_weight"`
	MaxWeight decimal.Decimal `json:"max_weight"`
}

// ESGConstraints defines ESG (Environmental, Social, Governance) constraints
type ESGConstraints struct {
	MinESGScore        decimal.Decimal `json:"min_esg_score"`
	MaxCarbonFootprint decimal.Decimal `json:"max_carbon_footprint"`
	ExcludedSectors    []string        `json:"excluded_sectors"`
}

// MarketView represents a Black-Litterman market view
type MarketView struct {
	Assets     []string          `json:"assets"`
	Weights    []decimal.Decimal `json:"weights"`
	Return     decimal.Decimal   `json:"return"`
	Confidence decimal.Decimal   `json:"confidence"`
}

// OptimizationResult contains the optimization results
type OptimizationResult struct {
	ID                uuid.UUID              `json:"id"`
	RequestID         uuid.UUID              `json:"request_id"`
	Objective         OptimizationObjective  `json:"objective"`
	OptimalWeights    []decimal.Decimal      `json:"optimal_weights"`
	ExpectedReturn    decimal.Decimal        `json:"expected_return"`
	ExpectedRisk      decimal.Decimal        `json:"expected_risk"`
	SharpeRatio       decimal.Decimal        `json:"sharpe_ratio"`
	Turnover          decimal.Decimal        `json:"turnover"`
	TransactionCosts  decimal.Decimal        `json:"transaction_costs"`
	Converged         bool                   `json:"converged"`
	Iterations        int                    `json:"iterations"`
	OptimizationTime  time.Duration          `json:"optimization_time"`
	EfficientFrontier []*EfficientPoint      `json:"efficient_frontier,omitempty"`
	RiskContribution  []decimal.Decimal      `json:"risk_contribution"`
	Metadata          map[string]interface{} `json:"metadata"`
	Timestamp         time.Time              `json:"timestamp"`
}

// EfficientPoint represents a point on the efficient frontier
type EfficientPoint struct {
	Return  decimal.Decimal   `json:"return"`
	Risk    decimal.Decimal   `json:"risk"`
	Weights []decimal.Decimal `json:"weights"`
}

// PortfolioMetrics contains portfolio performance metrics
type PortfolioMetrics struct {
	Return           decimal.Decimal `json:"return"`
	Volatility       decimal.Decimal `json:"volatility"`
	SharpeRatio      decimal.Decimal `json:"sharpe_ratio"`
	SortinoRatio     decimal.Decimal `json:"sortino_ratio"`
	MaxDrawdown      decimal.Decimal `json:"max_drawdown"`
	VaR95            decimal.Decimal `json:"var_95"`
	CVaR95           decimal.Decimal `json:"cvar_95"`
	Beta             decimal.Decimal `json:"beta"`
	Alpha            decimal.Decimal `json:"alpha"`
	InformationRatio decimal.Decimal `json:"information_ratio"`
	TrackingError    decimal.Decimal `json:"tracking_error"`
	TreynorRatio     decimal.Decimal `json:"treynor_ratio"`
}

// NewPortfolioOptimizer creates a new portfolio optimizer
func NewPortfolioOptimizer(logger *observability.Logger, config OptimizerConfig) *PortfolioOptimizer {
	// Set defaults
	if config.MaxIterations == 0 {
		config.MaxIterations = 1000
	}
	if config.Tolerance.IsZero() {
		config.Tolerance = decimal.NewFromFloat(1e-6)
	}
	if config.MinWeight.IsZero() {
		config.MinWeight = decimal.NewFromFloat(0.0)
	}
	if config.MaxWeight.IsZero() {
		config.MaxWeight = decimal.NewFromFloat(1.0)
	}
	if config.RebalanceThreshold.IsZero() {
		config.RebalanceThreshold = decimal.NewFromFloat(0.05) // 5%
	}
	if config.LookbackPeriod == 0 {
		config.LookbackPeriod = 252 * 24 * time.Hour // 252 trading days
	}

	return &PortfolioOptimizer{
		logger: logger,
		config: config,
	}
}

// Optimize performs portfolio optimization
func (po *PortfolioOptimizer) Optimize(ctx context.Context, request *OptimizationRequest) (*OptimizationResult, error) {
	startTime := time.Now()

	po.logger.Info(ctx, "Starting portfolio optimization", map[string]interface{}{
		"request_id": request.ID.String(),
		"objective":  string(request.Objective),
		"assets":     len(request.Assets),
	})

	// Validate request
	if err := po.validateRequest(request); err != nil {
		return nil, fmt.Errorf("invalid optimization request: %w", err)
	}

	// Prepare data
	if err := po.prepareData(request); err != nil {
		return nil, fmt.Errorf("failed to prepare data: %w", err)
	}

	var result *OptimizationResult
	var err error

	// Perform optimization based on objective
	switch request.Objective {
	case ObjectiveMaxSharpe:
		result, err = po.optimizeMaxSharpe(ctx, request)
	case ObjectiveMinVariance:
		result, err = po.optimizeMinVariance(ctx, request)
	case ObjectiveMaxReturn:
		result, err = po.optimizeMaxReturn(ctx, request)
	case ObjectiveRiskParity:
		result, err = po.optimizeRiskParity(ctx, request)
	case ObjectiveBlackLitterman:
		result, err = po.optimizeBlackLitterman(ctx, request)
	default:
		return nil, fmt.Errorf("unsupported optimization objective: %s", request.Objective)
	}

	if err != nil {
		return nil, fmt.Errorf("optimization failed: %w", err)
	}

	// Calculate additional metrics
	po.calculateMetrics(result, request)

	result.OptimizationTime = time.Since(startTime)
	result.Timestamp = time.Now()

	po.logger.Info(ctx, "Portfolio optimization completed", map[string]interface{}{
		"request_id":        request.ID.String(),
		"expected_return":   result.ExpectedReturn.String(),
		"expected_risk":     result.ExpectedRisk.String(),
		"sharpe_ratio":      result.SharpeRatio.String(),
		"converged":         result.Converged,
		"iterations":        result.Iterations,
		"optimization_time": result.OptimizationTime,
	})

	return result, nil
}

// CalculateEfficientFrontier calculates the efficient frontier
func (po *PortfolioOptimizer) CalculateEfficientFrontier(ctx context.Context, request *OptimizationRequest, numPoints int) ([]*EfficientPoint, error) {
	if numPoints <= 0 {
		numPoints = 50
	}

	po.logger.Info(ctx, "Calculating efficient frontier", map[string]interface{}{
		"assets":     len(request.Assets),
		"num_points": numPoints,
	})

	// Calculate minimum variance portfolio
	minVarRequest := *request
	minVarRequest.Objective = ObjectiveMinVariance
	minVarResult, err := po.optimizeMinVariance(ctx, &minVarRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate minimum variance portfolio: %w", err)
	}

	// Calculate maximum return portfolio
	maxRetRequest := *request
	maxRetRequest.Objective = ObjectiveMaxReturn
	maxRetResult, err := po.optimizeMaxReturn(ctx, &maxRetRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate maximum return portfolio: %w", err)
	}

	// Generate points along the efficient frontier
	frontier := make([]*EfficientPoint, numPoints)
	minReturn := minVarResult.ExpectedReturn
	maxReturn := maxRetResult.ExpectedReturn
	returnStep := maxReturn.Sub(minReturn).Div(decimal.NewFromInt(int64(numPoints - 1)))

	for i := 0; i < numPoints; i++ {
		targetReturn := minReturn.Add(returnStep.Mul(decimal.NewFromInt(int64(i))))

		// Optimize for target return
		targetRequest := *request
		targetRequest.TargetReturn = targetReturn

		result, err := po.optimizeTargetReturn(ctx, &targetRequest)
		if err != nil {
			continue // Skip points that fail to optimize
		}

		frontier[i] = &EfficientPoint{
			Return:  result.ExpectedReturn,
			Risk:    result.ExpectedRisk,
			Weights: result.OptimalWeights,
		}
	}

	// Remove nil points
	validFrontier := make([]*EfficientPoint, 0, numPoints)
	for _, point := range frontier {
		if point != nil {
			validFrontier = append(validFrontier, point)
		}
	}

	po.logger.Info(ctx, "Efficient frontier calculated", map[string]interface{}{
		"valid_points": len(validFrontier),
		"total_points": numPoints,
	})

	return validFrontier, nil
}

// CalculatePortfolioMetrics calculates comprehensive portfolio metrics
func (po *PortfolioOptimizer) CalculatePortfolioMetrics(ctx context.Context, weights []decimal.Decimal, returns [][]decimal.Decimal, benchmark []decimal.Decimal) (*PortfolioMetrics, error) {
	if len(weights) != len(returns[0]) {
		return nil, fmt.Errorf("weights and returns dimensions mismatch")
	}

	// Calculate portfolio returns
	portfolioReturns := make([]decimal.Decimal, len(returns))
	for i, periodReturns := range returns {
		portfolioReturn := decimal.NewFromInt(0)
		for j, assetReturn := range periodReturns {
			portfolioReturn = portfolioReturn.Add(weights[j].Mul(assetReturn))
		}
		portfolioReturns[i] = portfolioReturn
	}

	// Calculate basic metrics
	avgReturn := po.calculateMean(portfolioReturns)
	volatility := po.calculateStandardDeviation(portfolioReturns, avgReturn)

	// Annualize metrics (assuming daily returns)
	annualizedReturn := avgReturn.Mul(decimal.NewFromInt(252))
	annualizedVolatility := volatility.Mul(decimal.NewFromFloat(math.Sqrt(252)))

	// Calculate Sharpe ratio
	excessReturn := annualizedReturn.Sub(po.config.RiskFreeRate)
	sharpeRatio := decimal.NewFromInt(0)
	if annualizedVolatility.GreaterThan(decimal.NewFromInt(0)) {
		sharpeRatio = excessReturn.Div(annualizedVolatility)
	}

	// Calculate downside deviation for Sortino ratio
	downsideDeviation := po.calculateDownsideDeviation(portfolioReturns, po.config.RiskFreeRate.Div(decimal.NewFromInt(252)))
	annualizedDownsideDeviation := downsideDeviation.Mul(decimal.NewFromFloat(math.Sqrt(252)))

	sortinoRatio := decimal.NewFromInt(0)
	if annualizedDownsideDeviation.GreaterThan(decimal.NewFromInt(0)) {
		sortinoRatio = excessReturn.Div(annualizedDownsideDeviation)
	}

	// Calculate maximum drawdown
	maxDrawdown := po.calculateMaxDrawdown(portfolioReturns)

	// Calculate VaR and CVaR
	var95, cvar95 := po.calculateVaRAndCVaR(portfolioReturns, decimal.NewFromFloat(0.95))

	metrics := &PortfolioMetrics{
		Return:       annualizedReturn,
		Volatility:   annualizedVolatility,
		SharpeRatio:  sharpeRatio,
		SortinoRatio: sortinoRatio,
		MaxDrawdown:  maxDrawdown,
		VaR95:        var95,
		CVaR95:       cvar95,
	}

	// Calculate benchmark-relative metrics if benchmark provided
	if len(benchmark) > 0 && len(benchmark) == len(portfolioReturns) {
		po.calculateBenchmarkMetrics(metrics, portfolioReturns, benchmark)
	}

	return metrics, nil
}

// Private optimization methods

// validateRequest validates the optimization request
func (po *PortfolioOptimizer) validateRequest(request *OptimizationRequest) error {
	if len(request.Assets) == 0 {
		return fmt.Errorf("no assets provided")
	}

	if len(request.ExpectedReturns) != len(request.Assets) {
		return fmt.Errorf("expected returns length mismatch")
	}

	if len(request.CovarianceMatrix) != len(request.Assets) {
		return fmt.Errorf("covariance matrix dimension mismatch")
	}

	for i, row := range request.CovarianceMatrix {
		if len(row) != len(request.Assets) {
			return fmt.Errorf("covariance matrix row %d dimension mismatch", i)
		}
	}

	return nil
}

// prepareData prepares optimization data
func (po *PortfolioOptimizer) prepareData(request *OptimizationRequest) error {
	// Calculate expected returns if not provided
	if len(request.ExpectedReturns) == 0 && len(request.Returns) > 0 {
		request.ExpectedReturns = make([]decimal.Decimal, len(request.Assets))
		for i := 0; i < len(request.Assets); i++ {
			assetReturns := make([]decimal.Decimal, len(request.Returns))
			for j, periodReturns := range request.Returns {
				assetReturns[j] = periodReturns[i]
			}
			request.ExpectedReturns[i] = po.calculateMean(assetReturns)
		}
	}

	// Calculate covariance matrix if not provided
	if len(request.CovarianceMatrix) == 0 && len(request.Returns) > 0 {
		request.CovarianceMatrix = po.calculateCovarianceMatrix(request.Returns)
	}

	return nil
}

// optimizeMaxSharpe optimizes for maximum Sharpe ratio
func (po *PortfolioOptimizer) optimizeMaxSharpe(ctx context.Context, request *OptimizationRequest) (*OptimizationResult, error) {
	// This is a simplified implementation
	// In practice, this would use quadratic programming or other optimization algorithms

	numAssets := len(request.Assets)
	weights := make([]decimal.Decimal, numAssets)

	// Equal weight as starting point
	equalWeight := decimal.NewFromInt(1).Div(decimal.NewFromInt(int64(numAssets)))
	for i := range weights {
		weights[i] = equalWeight
	}

	// Simple iterative optimization (placeholder)
	bestSharpe := decimal.NewFromFloat(-999)
	bestWeights := make([]decimal.Decimal, numAssets)
	copy(bestWeights, weights)

	for iter := 0; iter < po.config.MaxIterations; iter++ {
		// Calculate portfolio metrics
		portfolioReturn := po.calculatePortfolioReturn(weights, request.ExpectedReturns)
		portfolioRisk := po.calculatePortfolioRisk(weights, request.CovarianceMatrix)

		if portfolioRisk.GreaterThan(decimal.NewFromInt(0)) {
			sharpe := portfolioReturn.Sub(po.config.RiskFreeRate).Div(portfolioRisk)
			if sharpe.GreaterThan(bestSharpe) {
				bestSharpe = sharpe
				copy(bestWeights, weights)
			}
		}

		// Simple gradient-like update (placeholder)
		// In practice, use proper optimization algorithms
		if iter < po.config.MaxIterations-1 {
			po.updateWeights(weights, request)
		}
	}

	result := &OptimizationResult{
		ID:             uuid.New(),
		RequestID:      request.ID,
		Objective:      request.Objective,
		OptimalWeights: bestWeights,
		ExpectedReturn: po.calculatePortfolioReturn(bestWeights, request.ExpectedReturns),
		ExpectedRisk:   po.calculatePortfolioRisk(bestWeights, request.CovarianceMatrix),
		SharpeRatio:    bestSharpe,
		Converged:      true,
		Iterations:     po.config.MaxIterations,
		Metadata:       make(map[string]interface{}),
	}

	return result, nil
}

// optimizeMinVariance optimizes for minimum variance
func (po *PortfolioOptimizer) optimizeMinVariance(ctx context.Context, request *OptimizationRequest) (*OptimizationResult, error) {
	numAssets := len(request.Assets)

	// For minimum variance, we want to minimize w^T * Σ * w subject to sum(w) = 1
	// This is a quadratic programming problem

	// Simplified implementation: equal weight
	weights := make([]decimal.Decimal, numAssets)
	equalWeight := decimal.NewFromInt(1).Div(decimal.NewFromInt(int64(numAssets)))
	for i := range weights {
		weights[i] = equalWeight
	}

	result := &OptimizationResult{
		ID:             uuid.New(),
		RequestID:      request.ID,
		Objective:      request.Objective,
		OptimalWeights: weights,
		ExpectedReturn: po.calculatePortfolioReturn(weights, request.ExpectedReturns),
		ExpectedRisk:   po.calculatePortfolioRisk(weights, request.CovarianceMatrix),
		Converged:      true,
		Iterations:     1,
		Metadata:       make(map[string]interface{}),
	}

	// Calculate Sharpe ratio
	if result.ExpectedRisk.GreaterThan(decimal.NewFromInt(0)) {
		result.SharpeRatio = result.ExpectedReturn.Sub(po.config.RiskFreeRate).Div(result.ExpectedRisk)
	}

	return result, nil
}

// optimizeMaxReturn optimizes for maximum return
func (po *PortfolioOptimizer) optimizeMaxReturn(ctx context.Context, request *OptimizationRequest) (*OptimizationResult, error) {
	numAssets := len(request.Assets)
	weights := make([]decimal.Decimal, numAssets)

	// Find asset with highest expected return
	maxReturnIndex := 0
	maxReturn := request.ExpectedReturns[0]
	for i, ret := range request.ExpectedReturns {
		if ret.GreaterThan(maxReturn) {
			maxReturn = ret
			maxReturnIndex = i
		}
	}

	// Allocate 100% to highest return asset (simplified)
	weights[maxReturnIndex] = decimal.NewFromInt(1)

	result := &OptimizationResult{
		ID:             uuid.New(),
		RequestID:      request.ID,
		Objective:      request.Objective,
		OptimalWeights: weights,
		ExpectedReturn: po.calculatePortfolioReturn(weights, request.ExpectedReturns),
		ExpectedRisk:   po.calculatePortfolioRisk(weights, request.CovarianceMatrix),
		Converged:      true,
		Iterations:     1,
		Metadata:       make(map[string]interface{}),
	}

	// Calculate Sharpe ratio
	if result.ExpectedRisk.GreaterThan(decimal.NewFromInt(0)) {
		result.SharpeRatio = result.ExpectedReturn.Sub(po.config.RiskFreeRate).Div(result.ExpectedRisk)
	}

	return result, nil
}

// optimizeRiskParity optimizes for risk parity
func (po *PortfolioOptimizer) optimizeRiskParity(ctx context.Context, request *OptimizationRequest) (*OptimizationResult, error) {
	numAssets := len(request.Assets)

	// Risk parity: equal risk contribution from each asset
	// This requires iterative optimization

	weights := make([]decimal.Decimal, numAssets)
	equalWeight := decimal.NewFromInt(1).Div(decimal.NewFromInt(int64(numAssets)))
	for i := range weights {
		weights[i] = equalWeight
	}

	// Simplified risk parity (equal weights as approximation)
	result := &OptimizationResult{
		ID:             uuid.New(),
		RequestID:      request.ID,
		Objective:      request.Objective,
		OptimalWeights: weights,
		ExpectedReturn: po.calculatePortfolioReturn(weights, request.ExpectedReturns),
		ExpectedRisk:   po.calculatePortfolioRisk(weights, request.CovarianceMatrix),
		Converged:      true,
		Iterations:     1,
		Metadata:       make(map[string]interface{}),
	}

	// Calculate Sharpe ratio
	if result.ExpectedRisk.GreaterThan(decimal.NewFromInt(0)) {
		result.SharpeRatio = result.ExpectedReturn.Sub(po.config.RiskFreeRate).Div(result.ExpectedRisk)
	}

	return result, nil
}

// optimizeBlackLitterman optimizes using Black-Litterman model
func (po *PortfolioOptimizer) optimizeBlackLitterman(ctx context.Context, request *OptimizationRequest) (*OptimizationResult, error) {
	// Black-Litterman model combines market equilibrium with investor views
	// This is a complex implementation that would require matrix operations

	// Simplified implementation: fall back to max Sharpe
	return po.optimizeMaxSharpe(ctx, request)
}

// optimizeTargetReturn optimizes for a target return level
func (po *PortfolioOptimizer) optimizeTargetReturn(ctx context.Context, request *OptimizationRequest) (*OptimizationResult, error) {
	// This would minimize risk subject to achieving target return
	// Simplified implementation
	return po.optimizeMinVariance(ctx, request)
}

// Helper calculation methods

// calculateMean calculates the mean of a slice of decimals
func (po *PortfolioOptimizer) calculateMean(values []decimal.Decimal) decimal.Decimal {
	if len(values) == 0 {
		return decimal.NewFromInt(0)
	}

	sum := decimal.NewFromInt(0)
	for _, value := range values {
		sum = sum.Add(value)
	}

	return sum.Div(decimal.NewFromInt(int64(len(values))))
}

// calculateStandardDeviation calculates the standard deviation
func (po *PortfolioOptimizer) calculateStandardDeviation(values []decimal.Decimal, mean decimal.Decimal) decimal.Decimal {
	if len(values) <= 1 {
		return decimal.NewFromInt(0)
	}

	sumSquaredDiffs := decimal.NewFromInt(0)
	for _, value := range values {
		diff := value.Sub(mean)
		sumSquaredDiffs = sumSquaredDiffs.Add(diff.Mul(diff))
	}

	variance := sumSquaredDiffs.Div(decimal.NewFromInt(int64(len(values) - 1)))
	return decimal.NewFromFloat(math.Sqrt(variance.InexactFloat64()))
}

// calculateDownsideDeviation calculates downside deviation
func (po *PortfolioOptimizer) calculateDownsideDeviation(returns []decimal.Decimal, threshold decimal.Decimal) decimal.Decimal {
	if len(returns) == 0 {
		return decimal.NewFromInt(0)
	}

	sumSquaredDownside := decimal.NewFromInt(0)
	count := 0

	for _, ret := range returns {
		if ret.LessThan(threshold) {
			diff := ret.Sub(threshold)
			sumSquaredDownside = sumSquaredDownside.Add(diff.Mul(diff))
			count++
		}
	}

	if count == 0 {
		return decimal.NewFromInt(0)
	}

	variance := sumSquaredDownside.Div(decimal.NewFromInt(int64(count)))
	return decimal.NewFromFloat(math.Sqrt(variance.InexactFloat64()))
}

// calculateMaxDrawdown calculates maximum drawdown
func (po *PortfolioOptimizer) calculateMaxDrawdown(returns []decimal.Decimal) decimal.Decimal {
	if len(returns) == 0 {
		return decimal.NewFromInt(0)
	}

	// Calculate cumulative returns
	cumulative := decimal.NewFromInt(1)
	peak := decimal.NewFromInt(1)
	maxDrawdown := decimal.NewFromInt(0)

	for _, ret := range returns {
		cumulative = cumulative.Mul(decimal.NewFromInt(1).Add(ret))

		if cumulative.GreaterThan(peak) {
			peak = cumulative
		}

		drawdown := peak.Sub(cumulative).Div(peak)
		if drawdown.GreaterThan(maxDrawdown) {
			maxDrawdown = drawdown
		}
	}

	return maxDrawdown
}

// calculateVaRAndCVaR calculates Value at Risk and Conditional VaR
func (po *PortfolioOptimizer) calculateVaRAndCVaR(returns []decimal.Decimal, confidence decimal.Decimal) (decimal.Decimal, decimal.Decimal) {
	if len(returns) == 0 {
		return decimal.NewFromInt(0), decimal.NewFromInt(0)
	}

	// Sort returns in ascending order
	sortedReturns := make([]decimal.Decimal, len(returns))
	copy(sortedReturns, returns)
	sort.Slice(sortedReturns, func(i, j int) bool {
		return sortedReturns[i].LessThan(sortedReturns[j])
	})

	// Calculate VaR
	alpha := decimal.NewFromInt(1).Sub(confidence)
	varIndex := int(alpha.Mul(decimal.NewFromInt(int64(len(sortedReturns)))).IntPart())
	if varIndex >= len(sortedReturns) {
		varIndex = len(sortedReturns) - 1
	}
	var_ := sortedReturns[varIndex].Neg()

	// Calculate CVaR (Expected Shortfall)
	tailSum := decimal.NewFromInt(0)
	tailCount := 0
	for i := 0; i <= varIndex; i++ {
		tailSum = tailSum.Add(sortedReturns[i])
		tailCount++
	}

	cvar := decimal.NewFromInt(0)
	if tailCount > 0 {
		cvar = tailSum.Div(decimal.NewFromInt(int64(tailCount))).Neg()
	}

	return var_, cvar
}

// calculatePortfolioReturn calculates expected portfolio return
func (po *PortfolioOptimizer) calculatePortfolioReturn(weights []decimal.Decimal, expectedReturns []decimal.Decimal) decimal.Decimal {
	portfolioReturn := decimal.NewFromInt(0)
	for i, weight := range weights {
		portfolioReturn = portfolioReturn.Add(weight.Mul(expectedReturns[i]))
	}
	return portfolioReturn
}

// calculatePortfolioRisk calculates portfolio risk (volatility)
func (po *PortfolioOptimizer) calculatePortfolioRisk(weights []decimal.Decimal, covMatrix [][]decimal.Decimal) decimal.Decimal {
	// Portfolio variance = w^T * Σ * w
	variance := decimal.NewFromInt(0)

	for i := 0; i < len(weights); i++ {
		for j := 0; j < len(weights); j++ {
			variance = variance.Add(weights[i].Mul(weights[j]).Mul(covMatrix[i][j]))
		}
	}

	// Return standard deviation (square root of variance)
	if variance.LessThanOrEqual(decimal.NewFromInt(0)) {
		return decimal.NewFromInt(0)
	}

	return decimal.NewFromFloat(math.Sqrt(variance.InexactFloat64()))
}

// calculateCovarianceMatrix calculates covariance matrix from returns
func (po *PortfolioOptimizer) calculateCovarianceMatrix(returns [][]decimal.Decimal) [][]decimal.Decimal {
	numAssets := len(returns[0])
	numPeriods := len(returns)

	// Calculate means for each asset
	means := make([]decimal.Decimal, numAssets)
	for i := 0; i < numAssets; i++ {
		assetReturns := make([]decimal.Decimal, numPeriods)
		for j := 0; j < numPeriods; j++ {
			assetReturns[j] = returns[j][i]
		}
		means[i] = po.calculateMean(assetReturns)
	}

	// Calculate covariance matrix
	covMatrix := make([][]decimal.Decimal, numAssets)
	for i := 0; i < numAssets; i++ {
		covMatrix[i] = make([]decimal.Decimal, numAssets)
		for j := 0; j < numAssets; j++ {
			covariance := decimal.NewFromInt(0)
			for k := 0; k < numPeriods; k++ {
				diff_i := returns[k][i].Sub(means[i])
				diff_j := returns[k][j].Sub(means[j])
				covariance = covariance.Add(diff_i.Mul(diff_j))
			}
			covMatrix[i][j] = covariance.Div(decimal.NewFromInt(int64(numPeriods - 1)))
		}
	}

	return covMatrix
}

// updateWeights updates weights during optimization (placeholder)
func (po *PortfolioOptimizer) updateWeights(weights []decimal.Decimal, request *OptimizationRequest) {
	// This is a placeholder for weight updates during optimization
	// In practice, this would implement gradient descent or other optimization algorithms

	// Simple random perturbation for demonstration
	for i := range weights {
		perturbation := decimal.NewFromFloat(0.01 * (0.5 - float64(i%2)))
		weights[i] = weights[i].Add(perturbation)

		// Ensure weights stay within bounds
		if weights[i].LessThan(po.config.MinWeight) {
			weights[i] = po.config.MinWeight
		}
		if weights[i].GreaterThan(po.config.MaxWeight) {
			weights[i] = po.config.MaxWeight
		}
	}

	// Normalize weights to sum to 1
	po.normalizeWeights(weights)
}

// normalizeWeights normalizes weights to sum to 1
func (po *PortfolioOptimizer) normalizeWeights(weights []decimal.Decimal) {
	sum := decimal.NewFromInt(0)
	for _, weight := range weights {
		sum = sum.Add(weight)
	}

	if sum.GreaterThan(decimal.NewFromInt(0)) {
		for i := range weights {
			weights[i] = weights[i].Div(sum)
		}
	}
}

// calculateMetrics calculates additional portfolio metrics
func (po *PortfolioOptimizer) calculateMetrics(result *OptimizationResult, request *OptimizationRequest) {
	// Calculate turnover if current weights provided
	if len(request.CurrentWeights) == len(result.OptimalWeights) {
		turnover := decimal.NewFromInt(0)
		for i, newWeight := range result.OptimalWeights {
			turnover = turnover.Add(newWeight.Sub(request.CurrentWeights[i]).Abs())
		}
		result.Turnover = turnover.Div(decimal.NewFromInt(2)) // Divide by 2 for one-way turnover
	}

	// Calculate risk contribution
	result.RiskContribution = po.calculateRiskContribution(result.OptimalWeights, request.CovarianceMatrix)
}

// calculateRiskContribution calculates risk contribution of each asset
func (po *PortfolioOptimizer) calculateRiskContribution(weights []decimal.Decimal, covMatrix [][]decimal.Decimal) []decimal.Decimal {
	numAssets := len(weights)
	riskContrib := make([]decimal.Decimal, numAssets)

	portfolioVariance := decimal.NewFromInt(0)
	for i := 0; i < numAssets; i++ {
		for j := 0; j < numAssets; j++ {
			portfolioVariance = portfolioVariance.Add(weights[i].Mul(weights[j]).Mul(covMatrix[i][j]))
		}
	}

	if portfolioVariance.GreaterThan(decimal.NewFromInt(0)) {
		for i := 0; i < numAssets; i++ {
			marginalContrib := decimal.NewFromInt(0)
			for j := 0; j < numAssets; j++ {
				marginalContrib = marginalContrib.Add(weights[j].Mul(covMatrix[i][j]))
			}
			riskContrib[i] = weights[i].Mul(marginalContrib).Div(portfolioVariance)
		}
	}

	return riskContrib
}

// calculateBenchmarkMetrics calculates benchmark-relative metrics
func (po *PortfolioOptimizer) calculateBenchmarkMetrics(metrics *PortfolioMetrics, portfolioReturns, benchmarkReturns []decimal.Decimal) {
	if len(portfolioReturns) != len(benchmarkReturns) {
		return
	}

	// Calculate excess returns
	excessReturns := make([]decimal.Decimal, len(portfolioReturns))
	for i := range portfolioReturns {
		excessReturns[i] = portfolioReturns[i].Sub(benchmarkReturns[i])
	}

	// Calculate tracking error
	excessMean := po.calculateMean(excessReturns)
	trackingError := po.calculateStandardDeviation(excessReturns, excessMean)
	metrics.TrackingError = trackingError.Mul(decimal.NewFromFloat(math.Sqrt(252))) // Annualized

	// Calculate information ratio
	annualizedExcessReturn := excessMean.Mul(decimal.NewFromInt(252))
	if metrics.TrackingError.GreaterThan(decimal.NewFromInt(0)) {
		metrics.InformationRatio = annualizedExcessReturn.Div(metrics.TrackingError)
	}

	// Calculate beta and alpha
	benchmarkMean := po.calculateMean(benchmarkReturns)
	benchmarkStdDev := po.calculateStandardDeviation(benchmarkReturns, benchmarkMean)

	if benchmarkStdDev.GreaterThan(decimal.NewFromInt(0)) {
		// Calculate covariance between portfolio and benchmark
		covariance := decimal.NewFromInt(0)
		portfolioMean := po.calculateMean(portfolioReturns)

		for i := range portfolioReturns {
			portfolioDiff := portfolioReturns[i].Sub(portfolioMean)
			benchmarkDiff := benchmarkReturns[i].Sub(benchmarkMean)
			covariance = covariance.Add(portfolioDiff.Mul(benchmarkDiff))
		}
		covariance = covariance.Div(decimal.NewFromInt(int64(len(portfolioReturns) - 1)))

		// Beta = Cov(portfolio, benchmark) / Var(benchmark)
		benchmarkVariance := benchmarkStdDev.Mul(benchmarkStdDev)
		metrics.Beta = covariance.Div(benchmarkVariance)

		// Alpha = Portfolio Return - (Risk-free Rate + Beta * (Benchmark Return - Risk-free Rate))
		benchmarkReturn := benchmarkMean.Mul(decimal.NewFromInt(252)) // Annualized
		expectedReturn := po.config.RiskFreeRate.Add(metrics.Beta.Mul(benchmarkReturn.Sub(po.config.RiskFreeRate)))
		metrics.Alpha = metrics.Return.Sub(expectedReturn)

		// Treynor ratio
		if metrics.Beta.GreaterThan(decimal.NewFromInt(0)) {
			excessPortfolioReturn := metrics.Return.Sub(po.config.RiskFreeRate)
			metrics.TreynorRatio = excessPortfolioReturn.Div(metrics.Beta)
		}
	}
}
