package trading

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// PortfolioOptimizer provides advanced portfolio optimization capabilities
type PortfolioOptimizer struct {
	logger            *observability.Logger
	config            *OptimizerConfig
	portfolios        map[string]*OptimizedPortfolio
	optimizers        map[string]Optimizer
	constraints       map[string]*OptimizationConstraints
	objectives        map[string]*OptimizationObjective
	marketData        map[string]*MarketData
	correlationMatrix map[string]map[string]decimal.Decimal
	mu                sync.RWMutex
	isRunning         bool
	stopChan          chan struct{}
}

// OptimizerConfig contains portfolio optimization configuration
type OptimizerConfig struct {
	OptimizationMethod   OptimizationMethod `json:"optimization_method"`
	RebalanceFrequency   time.Duration      `json:"rebalance_frequency"`
	LookbackPeriod       time.Duration      `json:"lookback_period"`
	MinWeight            decimal.Decimal    `json:"min_weight"`
	MaxWeight            decimal.Decimal    `json:"max_weight"`
	TransactionCosts     decimal.Decimal    `json:"transaction_costs"`
	RiskFreeRate         decimal.Decimal    `json:"risk_free_rate"`
	TargetReturn         decimal.Decimal    `json:"target_return"`
	MaxVolatility        decimal.Decimal    `json:"max_volatility"`
	MaxDrawdown          decimal.Decimal    `json:"max_drawdown"`
	EnableRebalancing    bool               `json:"enable_rebalancing"`
	EnableRiskBudgeting  bool               `json:"enable_risk_budgeting"`
	EnableBlackLitterman bool               `json:"enable_black_litterman"`
}

// OptimizationMethod defines optimization methods
type OptimizationMethod string

const (
	OptimizationMethodMeanVariance   OptimizationMethod = "mean_variance"
	OptimizationMethodMinVariance    OptimizationMethod = "min_variance"
	OptimizationMethodMaxSharpe      OptimizationMethod = "max_sharpe"
	OptimizationMethodRiskParity     OptimizationMethod = "risk_parity"
	OptimizationMethodBlackLitterman OptimizationMethod = "black_litterman"
	OptimizationMethodHierarchicalRP OptimizationMethod = "hierarchical_rp"
	OptimizationMethodCriticalLine   OptimizationMethod = "critical_line"
	OptimizationMethodMonteCarlo     OptimizationMethod = "monte_carlo"
)

// OptimizedPortfolio represents an optimized portfolio
type OptimizedPortfolio struct {
	ID                 string                     `json:"id"`
	Name               string                     `json:"name"`
	Method             OptimizationMethod         `json:"method"`
	Weights            map[string]decimal.Decimal `json:"weights"`
	ExpectedReturn     decimal.Decimal            `json:"expected_return"`
	ExpectedVolatility decimal.Decimal            `json:"expected_volatility"`
	SharpeRatio        decimal.Decimal            `json:"sharpe_ratio"`
	MaxDrawdown        decimal.Decimal            `json:"max_drawdown"`
	VaR95              decimal.Decimal            `json:"var_95"`
	CVaR95             decimal.Decimal            `json:"cvar_95"`
	Beta               decimal.Decimal            `json:"beta"`
	Alpha              decimal.Decimal            `json:"alpha"`
	TrackingError      decimal.Decimal            `json:"tracking_error"`
	InformationRatio   decimal.Decimal            `json:"information_ratio"`
	TurnoverRate       decimal.Decimal            `json:"turnover_rate"`
	TransactionCosts   decimal.Decimal            `json:"transaction_costs"`
	Constraints        *OptimizationConstraints   `json:"constraints"`
	Objective          *OptimizationObjective     `json:"objective"`
	Performance        *PortfolioPerformance      `json:"performance"`
	LastOptimized      time.Time                  `json:"last_optimized"`
	LastRebalanced     time.Time                  `json:"last_rebalanced"`
	IsActive           bool                       `json:"is_active"`
}

// OptimizationConstraints defines optimization constraints
type OptimizationConstraints struct {
	MinWeights           map[string]decimal.Decimal `json:"min_weights"`
	MaxWeights           map[string]decimal.Decimal `json:"max_weights"`
	SectorLimits         map[string]decimal.Decimal `json:"sector_limits"`
	CountryLimits        map[string]decimal.Decimal `json:"country_limits"`
	MaxVolatility        decimal.Decimal            `json:"max_volatility"`
	MinReturn            decimal.Decimal            `json:"min_return"`
	MaxDrawdown          decimal.Decimal            `json:"max_drawdown"`
	MaxTurnover          decimal.Decimal            `json:"max_turnover"`
	MaxConcentration     decimal.Decimal            `json:"max_concentration"`
	LiquidityRequirement decimal.Decimal            `json:"liquidity_requirement"`
	ESGScore             decimal.Decimal            `json:"esg_score"`
}

// OptimizationObjective defines optimization objective
type OptimizationObjective struct {
	Type             ObjectiveType              `json:"type"`
	TargetReturn     decimal.Decimal            `json:"target_return"`
	RiskAversion     decimal.Decimal            `json:"risk_aversion"`
	UtilityFunction  string                     `json:"utility_function"`
	BenchmarkWeights map[string]decimal.Decimal `json:"benchmark_weights"`
	RiskBudgets      map[string]decimal.Decimal `json:"risk_budgets"`
	Views            []*MarketView              `json:"views"`
	Confidence       decimal.Decimal            `json:"confidence"`
}

// ObjectiveType defines optimization objective types
type ObjectiveType string

const (
	ObjectiveTypeMaxReturn     ObjectiveType = "max_return"
	ObjectiveTypeMinRisk       ObjectiveType = "min_risk"
	ObjectiveTypeMaxSharpe     ObjectiveType = "max_sharpe"
	ObjectiveTypeMaxUtility    ObjectiveType = "max_utility"
	ObjectiveTypeRiskParity    ObjectiveType = "risk_parity"
	ObjectiveTypeTrackingError ObjectiveType = "tracking_error"
)

// MarketView represents a market view for Black-Litterman
type MarketView struct {
	Assets     []string        `json:"assets"`
	Return     decimal.Decimal `json:"return"`
	Confidence decimal.Decimal `json:"confidence"`
}

// PortfolioPerformance tracks portfolio performance
type PortfolioPerformance struct {
	TotalReturn      decimal.Decimal `json:"total_return"`
	AnnualizedReturn decimal.Decimal `json:"annualized_return"`
	Volatility       decimal.Decimal `json:"volatility"`
	SharpeRatio      decimal.Decimal `json:"sharpe_ratio"`
	SortinoRatio     decimal.Decimal `json:"sortino_ratio"`
	CalmarRatio      decimal.Decimal `json:"calmar_ratio"`
	MaxDrawdown      decimal.Decimal `json:"max_drawdown"`
	VaR95            decimal.Decimal `json:"var_95"`
	CVaR95           decimal.Decimal `json:"cvar_95"`
	Beta             decimal.Decimal `json:"beta"`
	Alpha            decimal.Decimal `json:"alpha"`
	TrackingError    decimal.Decimal `json:"tracking_error"`
	InformationRatio decimal.Decimal `json:"information_ratio"`
	WinRate          decimal.Decimal `json:"win_rate"`
	ProfitFactor     decimal.Decimal `json:"profit_factor"`
	LastUpdated      time.Time       `json:"last_updated"`
}

// MarketData represents market data for optimization
type MarketData struct {
	Symbol      string            `json:"symbol"`
	Price       decimal.Decimal   `json:"price"`
	Returns     []decimal.Decimal `json:"returns"`
	Volatility  decimal.Decimal   `json:"volatility"`
	Beta        decimal.Decimal   `json:"beta"`
	MarketCap   decimal.Decimal   `json:"market_cap"`
	Volume      decimal.Decimal   `json:"volume"`
	Liquidity   decimal.Decimal   `json:"liquidity"`
	ESGScore    decimal.Decimal   `json:"esg_score"`
	LastUpdated time.Time         `json:"last_updated"`
}

// Optimizer interface for different optimization algorithms
type Optimizer interface {
	Optimize(ctx context.Context, data *OptimizationData) (*OptimizedPortfolio, error)
	GetName() string
	GetMethod() OptimizationMethod
}

// OptimizationData contains data for optimization
type OptimizationData struct {
	Assets           []string                              `json:"assets"`
	Returns          map[string][]decimal.Decimal          `json:"returns"`
	Covariance       map[string]map[string]decimal.Decimal `json:"covariance"`
	ExpectedReturns  map[string]decimal.Decimal            `json:"expected_returns"`
	MarketCaps       map[string]decimal.Decimal            `json:"market_caps"`
	Constraints      *OptimizationConstraints              `json:"constraints"`
	Objective        *OptimizationObjective                `json:"objective"`
	CurrentWeights   map[string]decimal.Decimal            `json:"current_weights"`
	TransactionCosts decimal.Decimal                       `json:"transaction_costs"`
	RiskFreeRate     decimal.Decimal                       `json:"risk_free_rate"`
}

// NewPortfolioOptimizer creates a new portfolio optimizer
func NewPortfolioOptimizer(logger *observability.Logger) *PortfolioOptimizer {
	config := &OptimizerConfig{
		OptimizationMethod:   OptimizationMethodMeanVariance,
		RebalanceFrequency:   24 * time.Hour,
		LookbackPeriod:       252 * 24 * time.Hour,        // 1 year
		MinWeight:            decimal.NewFromFloat(0.01),  // 1%
		MaxWeight:            decimal.NewFromFloat(0.20),  // 20%
		TransactionCosts:     decimal.NewFromFloat(0.001), // 0.1%
		RiskFreeRate:         decimal.NewFromFloat(0.02),  // 2%
		TargetReturn:         decimal.NewFromFloat(0.10),  // 10%
		MaxVolatility:        decimal.NewFromFloat(0.15),  // 15%
		MaxDrawdown:          decimal.NewFromFloat(0.10),  // 10%
		EnableRebalancing:    true,
		EnableRiskBudgeting:  true,
		EnableBlackLitterman: false,
	}

	return &PortfolioOptimizer{
		logger:            logger,
		config:            config,
		portfolios:        make(map[string]*OptimizedPortfolio),
		optimizers:        make(map[string]Optimizer),
		constraints:       make(map[string]*OptimizationConstraints),
		objectives:        make(map[string]*OptimizationObjective),
		marketData:        make(map[string]*MarketData),
		correlationMatrix: make(map[string]map[string]decimal.Decimal),
		stopChan:          make(chan struct{}),
	}
}

// Start starts the portfolio optimizer
func (po *PortfolioOptimizer) Start(ctx context.Context) error {
	po.mu.Lock()
	defer po.mu.Unlock()

	if po.isRunning {
		return fmt.Errorf("portfolio optimizer is already running")
	}

	po.isRunning = true

	// Initialize optimizers
	po.initializeOptimizers()

	// Start background processes
	go po.optimizationLoop(ctx)
	go po.rebalancingLoop(ctx)
	go po.performanceTrackingLoop(ctx)

	po.logger.Info(ctx, "Portfolio optimizer started", map[string]interface{}{
		"optimization_method": po.config.OptimizationMethod,
		"rebalance_frequency": po.config.RebalanceFrequency,
		"optimizers_count":    len(po.optimizers),
	})

	return nil
}

// Stop stops the portfolio optimizer
func (po *PortfolioOptimizer) Stop(ctx context.Context) error {
	po.mu.Lock()
	defer po.mu.Unlock()

	if !po.isRunning {
		return nil
	}

	po.isRunning = false
	close(po.stopChan)

	po.logger.Info(ctx, "Portfolio optimizer stopped", nil)
	return nil
}

// OptimizePortfolio optimizes a portfolio using the specified method
func (po *PortfolioOptimizer) OptimizePortfolio(ctx context.Context, name string, assets []string, method OptimizationMethod, constraints *OptimizationConstraints, objective *OptimizationObjective) (*OptimizedPortfolio, error) {
	po.mu.Lock()
	defer po.mu.Unlock()

	// Get optimizer
	optimizer, exists := po.optimizers[string(method)]
	if !exists {
		return nil, fmt.Errorf("optimizer not found for method: %s", method)
	}

	// Prepare optimization data
	data, err := po.prepareOptimizationData(assets, constraints, objective)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare optimization data: %w", err)
	}

	// Perform optimization
	portfolio, err := optimizer.Optimize(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("optimization failed: %w", err)
	}

	// Set portfolio properties
	portfolio.ID = uuid.New().String()
	portfolio.Name = name
	portfolio.Method = method
	portfolio.Constraints = constraints
	portfolio.Objective = objective
	portfolio.LastOptimized = time.Now()
	portfolio.IsActive = true

	// Calculate performance metrics
	po.calculatePerformanceMetrics(portfolio, data)

	// Store portfolio
	po.portfolios[portfolio.ID] = portfolio

	po.logger.Info(ctx, "Portfolio optimized", map[string]interface{}{
		"portfolio_id":        portfolio.ID,
		"portfolio_name":      portfolio.Name,
		"optimization_method": portfolio.Method,
		"expected_return":     portfolio.ExpectedReturn.String(),
		"expected_volatility": portfolio.ExpectedVolatility.String(),
		"sharpe_ratio":        portfolio.SharpeRatio.String(),
	})

	return portfolio, nil
}

// GetPortfolio retrieves an optimized portfolio
func (po *PortfolioOptimizer) GetPortfolio(portfolioID string) (*OptimizedPortfolio, error) {
	po.mu.RLock()
	defer po.mu.RUnlock()

	portfolio, exists := po.portfolios[portfolioID]
	if !exists {
		return nil, fmt.Errorf("portfolio not found: %s", portfolioID)
	}

	return portfolio, nil
}

// GetActivePortfolios returns all active portfolios
func (po *PortfolioOptimizer) GetActivePortfolios() []*OptimizedPortfolio {
	po.mu.RLock()
	defer po.mu.RUnlock()

	var activePortfolios []*OptimizedPortfolio
	for _, portfolio := range po.portfolios {
		if portfolio.IsActive {
			activePortfolios = append(activePortfolios, portfolio)
		}
	}

	return activePortfolios
}

// RebalancePortfolio rebalances a portfolio
func (po *PortfolioOptimizer) RebalancePortfolio(ctx context.Context, portfolioID string) error {
	po.mu.Lock()
	defer po.mu.Unlock()

	portfolio, exists := po.portfolios[portfolioID]
	if !exists {
		return fmt.Errorf("portfolio not found: %s", portfolioID)
	}

	// Get current market data
	assets := make([]string, 0, len(portfolio.Weights))
	for asset := range portfolio.Weights {
		assets = append(assets, asset)
	}

	// Prepare optimization data
	data, err := po.prepareOptimizationData(assets, portfolio.Constraints, portfolio.Objective)
	if err != nil {
		return fmt.Errorf("failed to prepare rebalancing data: %w", err)
	}

	// Set current weights
	data.CurrentWeights = portfolio.Weights

	// Get optimizer
	optimizer, exists := po.optimizers[string(portfolio.Method)]
	if !exists {
		return fmt.Errorf("optimizer not found for method: %s", portfolio.Method)
	}

	// Perform rebalancing optimization
	rebalancedPortfolio, err := optimizer.Optimize(ctx, data)
	if err != nil {
		return fmt.Errorf("rebalancing optimization failed: %w", err)
	}

	// Calculate turnover
	turnover := po.calculateTurnover(portfolio.Weights, rebalancedPortfolio.Weights)

	// Update portfolio
	portfolio.Weights = rebalancedPortfolio.Weights
	portfolio.ExpectedReturn = rebalancedPortfolio.ExpectedReturn
	portfolio.ExpectedVolatility = rebalancedPortfolio.ExpectedVolatility
	portfolio.SharpeRatio = rebalancedPortfolio.SharpeRatio
	portfolio.TurnoverRate = turnover
	portfolio.LastRebalanced = time.Now()

	// Calculate transaction costs
	portfolio.TransactionCosts = turnover.Mul(po.config.TransactionCosts)

	po.logger.Info(ctx, "Portfolio rebalanced", map[string]interface{}{
		"portfolio_id":      portfolioID,
		"turnover_rate":     turnover.String(),
		"transaction_costs": portfolio.TransactionCosts.String(),
		"new_sharpe_ratio":  portfolio.SharpeRatio.String(),
	})

	return nil
}

// initializeOptimizers initializes optimization algorithms
func (po *PortfolioOptimizer) initializeOptimizers() {
	// Initialize mean-variance optimizer
	po.optimizers[string(OptimizationMethodMeanVariance)] = &MeanVarianceOptimizer{}

	// Initialize minimum variance optimizer
	po.optimizers[string(OptimizationMethodMinVariance)] = &MinVarianceOptimizer{}

	// Initialize maximum Sharpe optimizer
	po.optimizers[string(OptimizationMethodMaxSharpe)] = &MaxSharpeOptimizer{}

	// Initialize risk parity optimizer
	po.optimizers[string(OptimizationMethodRiskParity)] = &RiskParityOptimizer{}
}

// prepareOptimizationData prepares data for optimization
func (po *PortfolioOptimizer) prepareOptimizationData(assets []string, constraints *OptimizationConstraints, objective *OptimizationObjective) (*OptimizationData, error) {
	data := &OptimizationData{
		Assets:           assets,
		Returns:          make(map[string][]decimal.Decimal),
		Covariance:       make(map[string]map[string]decimal.Decimal),
		ExpectedReturns:  make(map[string]decimal.Decimal),
		MarketCaps:       make(map[string]decimal.Decimal),
		Constraints:      constraints,
		Objective:        objective,
		CurrentWeights:   make(map[string]decimal.Decimal),
		TransactionCosts: po.config.TransactionCosts,
		RiskFreeRate:     po.config.RiskFreeRate,
	}

	// Generate sample data (in production, this would come from market data)
	for _, asset := range assets {
		// Generate sample returns
		returns := make([]decimal.Decimal, 252) // 1 year of daily returns
		for i := range returns {
			returns[i] = decimal.NewFromFloat(0.001 * (math.Sin(float64(i)*0.1) + 0.5*math.Cos(float64(i)*0.05)))
		}
		data.Returns[asset] = returns

		// Calculate expected return (mean of returns)
		sum := decimal.Zero
		for _, ret := range returns {
			sum = sum.Add(ret)
		}
		data.ExpectedReturns[asset] = sum.Div(decimal.NewFromInt(int64(len(returns))))

		// Set market cap (sample data)
		data.MarketCaps[asset] = decimal.NewFromFloat(1000000000) // $1B
	}

	// Calculate covariance matrix
	po.calculateCovarianceMatrix(data)

	return data, nil
}

// calculateCovarianceMatrix calculates the covariance matrix
func (po *PortfolioOptimizer) calculateCovarianceMatrix(data *OptimizationData) {
	for _, asset1 := range data.Assets {
		data.Covariance[asset1] = make(map[string]decimal.Decimal)
		for _, asset2 := range data.Assets {
			if asset1 == asset2 {
				// Variance
				variance := po.calculateVariance(data.Returns[asset1])
				data.Covariance[asset1][asset2] = variance
			} else {
				// Covariance
				covariance := po.calculateCovariance(data.Returns[asset1], data.Returns[asset2])
				data.Covariance[asset1][asset2] = covariance
			}
		}
	}
}

// calculateVariance calculates variance of returns
func (po *PortfolioOptimizer) calculateVariance(returns []decimal.Decimal) decimal.Decimal {
	if len(returns) == 0 {
		return decimal.Zero
	}

	// Calculate mean
	sum := decimal.Zero
	for _, ret := range returns {
		sum = sum.Add(ret)
	}
	mean := sum.Div(decimal.NewFromInt(int64(len(returns))))

	// Calculate variance
	sumSquaredDiff := decimal.Zero
	for _, ret := range returns {
		diff := ret.Sub(mean)
		sumSquaredDiff = sumSquaredDiff.Add(diff.Mul(diff))
	}

	return sumSquaredDiff.Div(decimal.NewFromInt(int64(len(returns) - 1)))
}

// calculateCovariance calculates covariance between two return series
func (po *PortfolioOptimizer) calculateCovariance(returns1, returns2 []decimal.Decimal) decimal.Decimal {
	if len(returns1) != len(returns2) || len(returns1) == 0 {
		return decimal.Zero
	}

	// Calculate means
	sum1, sum2 := decimal.Zero, decimal.Zero
	for i := range returns1 {
		sum1 = sum1.Add(returns1[i])
		sum2 = sum2.Add(returns2[i])
	}
	mean1 := sum1.Div(decimal.NewFromInt(int64(len(returns1))))
	mean2 := sum2.Div(decimal.NewFromInt(int64(len(returns2))))

	// Calculate covariance
	sumProduct := decimal.Zero
	for i := range returns1 {
		diff1 := returns1[i].Sub(mean1)
		diff2 := returns2[i].Sub(mean2)
		sumProduct = sumProduct.Add(diff1.Mul(diff2))
	}

	return sumProduct.Div(decimal.NewFromInt(int64(len(returns1) - 1)))
}

// calculatePerformanceMetrics calculates portfolio performance metrics
func (po *PortfolioOptimizer) calculatePerformanceMetrics(portfolio *OptimizedPortfolio, data *OptimizationData) {
	// Calculate portfolio variance
	portfolioVariance := decimal.Zero
	for asset1, weight1 := range portfolio.Weights {
		for asset2, weight2 := range portfolio.Weights {
			covariance := data.Covariance[asset1][asset2]
			portfolioVariance = portfolioVariance.Add(weight1.Mul(weight2).Mul(covariance))
		}
	}

	// Calculate portfolio volatility (annualized)
	portfolio.ExpectedVolatility = portfolioVariance.Pow(decimal.NewFromFloat(0.5)).Mul(decimal.NewFromFloat(math.Sqrt(252)))

	// Calculate Sharpe ratio
	if portfolio.ExpectedVolatility.GreaterThan(decimal.Zero) {
		excessReturn := portfolio.ExpectedReturn.Sub(data.RiskFreeRate)
		portfolio.SharpeRatio = excessReturn.Div(portfolio.ExpectedVolatility)
	}

	// Initialize performance tracking
	portfolio.Performance = &PortfolioPerformance{
		LastUpdated: time.Now(),
	}
}

// calculateTurnover calculates portfolio turnover
func (po *PortfolioOptimizer) calculateTurnover(oldWeights, newWeights map[string]decimal.Decimal) decimal.Decimal {
	turnover := decimal.Zero

	// Get all assets
	allAssets := make(map[string]bool)
	for asset := range oldWeights {
		allAssets[asset] = true
	}
	for asset := range newWeights {
		allAssets[asset] = true
	}

	// Calculate turnover
	for asset := range allAssets {
		oldWeight := oldWeights[asset]
		newWeight := newWeights[asset]
		diff := newWeight.Sub(oldWeight).Abs()
		turnover = turnover.Add(diff)
	}

	return turnover.Div(decimal.NewFromInt(2)) // Divide by 2 for one-way turnover
}

// optimizationLoop performs periodic optimization
func (po *PortfolioOptimizer) optimizationLoop(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-po.stopChan:
			return
		case <-ticker.C:
			po.performPeriodicOptimization(ctx)
		}
	}
}

// rebalancingLoop performs periodic rebalancing
func (po *PortfolioOptimizer) rebalancingLoop(ctx context.Context) {
	ticker := time.NewTicker(po.config.RebalanceFrequency)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-po.stopChan:
			return
		case <-ticker.C:
			if po.config.EnableRebalancing {
				po.performPeriodicRebalancing(ctx)
			}
		}
	}
}

// performanceTrackingLoop tracks portfolio performance
func (po *PortfolioOptimizer) performanceTrackingLoop(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-po.stopChan:
			return
		case <-ticker.C:
			po.updatePerformanceMetrics(ctx)
		}
	}
}

// performPeriodicOptimization performs periodic optimization for all portfolios
func (po *PortfolioOptimizer) performPeriodicOptimization(ctx context.Context) {
	po.mu.RLock()
	portfolios := make([]*OptimizedPortfolio, 0, len(po.portfolios))
	for _, portfolio := range po.portfolios {
		if portfolio.IsActive {
			portfolios = append(portfolios, portfolio)
		}
	}
	po.mu.RUnlock()

	for _, portfolio := range portfolios {
		if time.Since(portfolio.LastOptimized) > 24*time.Hour {
			// Re-optimize portfolio
			assets := make([]string, 0, len(portfolio.Weights))
			for asset := range portfolio.Weights {
				assets = append(assets, asset)
			}

			_, err := po.OptimizePortfolio(ctx, portfolio.Name, assets, portfolio.Method, portfolio.Constraints, portfolio.Objective)
			if err != nil {
				po.logger.Error(ctx, "Failed to re-optimize portfolio", err)
			}
		}
	}
}

// performPeriodicRebalancing performs periodic rebalancing for all portfolios
func (po *PortfolioOptimizer) performPeriodicRebalancing(ctx context.Context) {
	po.mu.RLock()
	portfolios := make([]*OptimizedPortfolio, 0, len(po.portfolios))
	for _, portfolio := range po.portfolios {
		if portfolio.IsActive {
			portfolios = append(portfolios, portfolio)
		}
	}
	po.mu.RUnlock()

	for _, portfolio := range portfolios {
		if time.Since(portfolio.LastRebalanced) > po.config.RebalanceFrequency {
			err := po.RebalancePortfolio(ctx, portfolio.ID)
			if err != nil {
				po.logger.Error(ctx, "Failed to rebalance portfolio", err)
			}
		}
	}
}

// updatePerformanceMetrics updates performance metrics for all portfolios
func (po *PortfolioOptimizer) updatePerformanceMetrics(ctx context.Context) {
	po.mu.RLock()
	defer po.mu.RUnlock()

	for _, portfolio := range po.portfolios {
		if portfolio.IsActive && portfolio.Performance != nil {
			// Update performance metrics (simplified)
			portfolio.Performance.LastUpdated = time.Now()
		}
	}
}

// MeanVarianceOptimizer implements mean-variance optimization
type MeanVarianceOptimizer struct{}

func (mvo *MeanVarianceOptimizer) GetName() string {
	return "Mean-Variance Optimizer"
}

func (mvo *MeanVarianceOptimizer) GetMethod() OptimizationMethod {
	return OptimizationMethodMeanVariance
}

func (mvo *MeanVarianceOptimizer) Optimize(ctx context.Context, data *OptimizationData) (*OptimizedPortfolio, error) {
	// Simplified mean-variance optimization
	weights := make(map[string]decimal.Decimal)
	numAssets := len(data.Assets)
	equalWeight := decimal.NewFromFloat(1.0).Div(decimal.NewFromInt(int64(numAssets)))

	// Start with equal weights and adjust based on expected returns
	for _, asset := range data.Assets {
		expectedReturn := data.ExpectedReturns[asset]
		// Simple adjustment based on expected return
		adjustment := expectedReturn.Mul(decimal.NewFromFloat(0.5))
		weight := equalWeight.Add(adjustment)

		// Apply constraints
		if data.Constraints != nil {
			if minWeight, exists := data.Constraints.MinWeights[asset]; exists && weight.LessThan(minWeight) {
				weight = minWeight
			}
			if maxWeight, exists := data.Constraints.MaxWeights[asset]; exists && weight.GreaterThan(maxWeight) {
				weight = maxWeight
			}
		}

		weights[asset] = weight
	}

	// Normalize weights to sum to 1
	totalWeight := decimal.Zero
	for _, weight := range weights {
		totalWeight = totalWeight.Add(weight)
	}
	for asset, weight := range weights {
		weights[asset] = weight.Div(totalWeight)
	}

	// Calculate expected portfolio return
	expectedReturn := decimal.Zero
	for asset, weight := range weights {
		expectedReturn = expectedReturn.Add(weight.Mul(data.ExpectedReturns[asset]))
	}

	return &OptimizedPortfolio{
		Weights:        weights,
		ExpectedReturn: expectedReturn.Mul(decimal.NewFromFloat(252)), // Annualized
	}, nil
}

// MinVarianceOptimizer implements minimum variance optimization
type MinVarianceOptimizer struct{}

func (mvo *MinVarianceOptimizer) GetName() string {
	return "Minimum Variance Optimizer"
}

func (mvo *MinVarianceOptimizer) GetMethod() OptimizationMethod {
	return OptimizationMethodMinVariance
}

func (mvo *MinVarianceOptimizer) Optimize(ctx context.Context, data *OptimizationData) (*OptimizedPortfolio, error) {
	// Simplified minimum variance optimization
	weights := make(map[string]decimal.Decimal)

	// Inverse volatility weighting (simplified approach)
	totalInverseVol := decimal.Zero
	inverseVols := make(map[string]decimal.Decimal)

	for _, asset := range data.Assets {
		variance := data.Covariance[asset][asset]
		if variance.GreaterThan(decimal.Zero) {
			volatility := variance.Pow(decimal.NewFromFloat(0.5))
			inverseVol := decimal.NewFromFloat(1.0).Div(volatility)
			inverseVols[asset] = inverseVol
			totalInverseVol = totalInverseVol.Add(inverseVol)
		} else {
			// Equal weight if no variance data
			inverseVols[asset] = decimal.NewFromFloat(1.0)
			totalInverseVol = totalInverseVol.Add(decimal.NewFromFloat(1.0))
		}
	}

	// Calculate weights
	for _, asset := range data.Assets {
		weight := inverseVols[asset].Div(totalInverseVol)
		weights[asset] = weight
	}

	// Calculate expected portfolio return
	expectedReturn := decimal.Zero
	for asset, weight := range weights {
		expectedReturn = expectedReturn.Add(weight.Mul(data.ExpectedReturns[asset]))
	}

	return &OptimizedPortfolio{
		Weights:        weights,
		ExpectedReturn: expectedReturn.Mul(decimal.NewFromFloat(252)), // Annualized
	}, nil
}

// MaxSharpeOptimizer implements maximum Sharpe ratio optimization
type MaxSharpeOptimizer struct{}

func (mso *MaxSharpeOptimizer) GetName() string {
	return "Maximum Sharpe Optimizer"
}

func (mso *MaxSharpeOptimizer) GetMethod() OptimizationMethod {
	return OptimizationMethodMaxSharpe
}

func (mso *MaxSharpeOptimizer) Optimize(ctx context.Context, data *OptimizationData) (*OptimizedPortfolio, error) {
	// Simplified maximum Sharpe ratio optimization
	weights := make(map[string]decimal.Decimal)

	// Calculate excess returns
	totalExcessReturn := decimal.Zero
	excessReturns := make(map[string]decimal.Decimal)

	for _, asset := range data.Assets {
		excessReturn := data.ExpectedReturns[asset].Sub(data.RiskFreeRate)
		variance := data.Covariance[asset][asset]

		if variance.GreaterThan(decimal.Zero) {
			// Weight by excess return / variance (simplified)
			weight := excessReturn.Div(variance)
			if weight.GreaterThan(decimal.Zero) {
				excessReturns[asset] = weight
				totalExcessReturn = totalExcessReturn.Add(weight)
			}
		}
	}

	// Normalize weights
	if totalExcessReturn.GreaterThan(decimal.Zero) {
		for _, asset := range data.Assets {
			if excessReturn, exists := excessReturns[asset]; exists {
				weights[asset] = excessReturn.Div(totalExcessReturn)
			} else {
				weights[asset] = decimal.Zero
			}
		}
	} else {
		// Fall back to equal weights
		equalWeight := decimal.NewFromFloat(1.0).Div(decimal.NewFromInt(int64(len(data.Assets))))
		for _, asset := range data.Assets {
			weights[asset] = equalWeight
		}
	}

	// Calculate expected portfolio return
	expectedReturn := decimal.Zero
	for asset, weight := range weights {
		expectedReturn = expectedReturn.Add(weight.Mul(data.ExpectedReturns[asset]))
	}

	return &OptimizedPortfolio{
		Weights:        weights,
		ExpectedReturn: expectedReturn.Mul(decimal.NewFromFloat(252)), // Annualized
	}, nil
}

// RiskParityOptimizer implements risk parity optimization
type RiskParityOptimizer struct{}

func (rpo *RiskParityOptimizer) GetName() string {
	return "Risk Parity Optimizer"
}

func (rpo *RiskParityOptimizer) GetMethod() OptimizationMethod {
	return OptimizationMethodRiskParity
}

func (rpo *RiskParityOptimizer) Optimize(ctx context.Context, data *OptimizationData) (*OptimizedPortfolio, error) {
	// Simplified risk parity optimization
	weights := make(map[string]decimal.Decimal)

	// Inverse volatility weighting (equal risk contribution)
	totalInverseVol := decimal.Zero
	inverseVols := make(map[string]decimal.Decimal)

	for _, asset := range data.Assets {
		variance := data.Covariance[asset][asset]
		if variance.GreaterThan(decimal.Zero) {
			volatility := variance.Pow(decimal.NewFromFloat(0.5))
			inverseVol := decimal.NewFromFloat(1.0).Div(volatility)
			inverseVols[asset] = inverseVol
			totalInverseVol = totalInverseVol.Add(inverseVol)
		} else {
			inverseVols[asset] = decimal.NewFromFloat(1.0)
			totalInverseVol = totalInverseVol.Add(decimal.NewFromFloat(1.0))
		}
	}

	// Calculate risk parity weights
	for _, asset := range data.Assets {
		weight := inverseVols[asset].Div(totalInverseVol)
		weights[asset] = weight
	}

	// Calculate expected portfolio return
	expectedReturn := decimal.Zero
	for asset, weight := range weights {
		expectedReturn = expectedReturn.Add(weight.Mul(data.ExpectedReturns[asset]))
	}

	return &OptimizedPortfolio{
		Weights:        weights,
		ExpectedReturn: expectedReturn.Mul(decimal.NewFromFloat(252)), // Annualized
	}, nil
}
