package risk

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/internal/compliance"
	"github.com/ai-agentic-browser/internal/exchanges"
	"github.com/ai-agentic-browser/internal/strategies/framework"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// RiskManagementService provides comprehensive risk management capabilities
type RiskManagementService struct {
	logger            *observability.Logger
	config            RiskServiceConfig
	riskEngine        *RiskEngine
	varCalculator     *VaRCalculator
	complianceManager *compliance.ComplianceManager
	exchangeManager   *exchanges.Manager

	// Portfolio tracking
	portfolios    map[uuid.UUID]*PortfolioRisk
	globalMetrics *GlobalRiskMetrics

	// Circuit breakers
	circuitBreakers map[string]*CircuitBreaker

	// Performance tracking
	totalRiskChecks     int64
	riskViolations      int64
	circuitBreakerTrips int64

	// State management
	isRunning bool
	stopChan  chan struct{}
	wg        sync.WaitGroup
	mu        sync.RWMutex
}

// RiskServiceConfig contains risk service configuration
type RiskServiceConfig struct {
	EnableRiskEngine      bool              `json:"enable_risk_engine"`
	EnableVaRCalculation  bool              `json:"enable_var_calculation"`
	EnableCompliance      bool              `json:"enable_compliance"`
	EnableCircuitBreakers bool              `json:"enable_circuit_breakers"`
	UpdateInterval        time.Duration     `json:"update_interval"`
	VaRUpdateInterval     time.Duration     `json:"var_update_interval"`
	MaxPortfolios         int               `json:"max_portfolios"`
	GlobalLimits          *GlobalRiskLimits `json:"global_limits"`
	VaRConfig             VaRConfig         `json:"var_config"`
	RiskConfig            RiskConfig        `json:"risk_config"`
}

// PortfolioRisk tracks risk metrics for a portfolio
type PortfolioRisk struct {
	ID               uuid.UUID                `json:"id"`
	StrategyID       uuid.UUID                `json:"strategy_id"`
	Name             string                   `json:"name"`
	TotalValue       decimal.Decimal          `json:"total_value"`
	Cash             decimal.Decimal          `json:"cash"`
	Positions        map[string]*PositionRisk `json:"positions"`
	Metrics          *RiskMetrics             `json:"metrics"`
	VaRResult        *VaRResult               `json:"var_result,omitempty"`
	LastUpdated      time.Time                `json:"last_updated"`
	RiskScore        decimal.Decimal          `json:"risk_score"`
	ComplianceStatus string                   `json:"compliance_status"`
}

// PositionRisk tracks risk metrics for a position
type PositionRisk struct {
	Symbol           string          `json:"symbol"`
	Quantity         decimal.Decimal `json:"quantity"`
	AvgPrice         decimal.Decimal `json:"avg_price"`
	CurrentPrice     decimal.Decimal `json:"current_price"`
	MarketValue      decimal.Decimal `json:"market_value"`
	UnrealizedPnL    decimal.Decimal `json:"unrealized_pnl"`
	RealizedPnL      decimal.Decimal `json:"realized_pnl"`
	DailyPnL         decimal.Decimal `json:"daily_pnl"`
	Volatility       decimal.Decimal `json:"volatility"`
	Beta             decimal.Decimal `json:"beta"`
	VaR              decimal.Decimal `json:"var"`
	RiskContribution decimal.Decimal `json:"risk_contribution"`
	LastUpdated      time.Time       `json:"last_updated"`
}

// GlobalRiskMetrics tracks system-wide risk metrics
type GlobalRiskMetrics struct {
	TotalPortfolioValue  decimal.Decimal `json:"total_portfolio_value"`
	TotalExposure        decimal.Decimal `json:"total_exposure"`
	TotalVaR             decimal.Decimal `json:"total_var"`
	MaxDrawdown          decimal.Decimal `json:"max_drawdown"`
	AverageRiskScore     decimal.Decimal `json:"average_risk_score"`
	ActivePortfolios     int             `json:"active_portfolios"`
	ActiveStrategies     int             `json:"active_strategies"`
	TotalPositions       int             `json:"total_positions"`
	ComplianceViolations int             `json:"compliance_violations"`
	CircuitBreakerTrips  int             `json:"circuit_breaker_trips"`
	LastUpdated          time.Time       `json:"last_updated"`
}

// RiskCheckRequest represents a request to check risk for a trading signal
type RiskCheckRequest struct {
	Signal      *framework.Signal `json:"signal"`
	PortfolioID uuid.UUID         `json:"portfolio_id"`
	StrategyID  uuid.UUID         `json:"strategy_id"`
	UserID      string            `json:"user_id,omitempty"`
	SessionID   string            `json:"session_id,omitempty"`
}

// RiskCheckResponse represents the response from a risk check
type RiskCheckResponse struct {
	Approved         bool                   `json:"approved"`
	RiskScore        decimal.Decimal        `json:"risk_score"`
	Violations       []*RiskViolation       `json:"violations"`
	Warnings         []string               `json:"warnings"`
	CircuitBreakers  []string               `json:"circuit_breakers"`
	ComplianceStatus string                 `json:"compliance_status"`
	Details          map[string]interface{} `json:"details"`
	ProcessedAt      time.Time              `json:"processed_at"`
}

// NewRiskManagementService creates a new risk management service
func NewRiskManagementService(
	logger *observability.Logger,
	config RiskServiceConfig,
	exchangeManager *exchanges.Manager,
	complianceManager *compliance.ComplianceManager,
) *RiskManagementService {
	// Set defaults
	if config.UpdateInterval == 0 {
		config.UpdateInterval = 1 * time.Second
	}
	if config.VaRUpdateInterval == 0 {
		config.VaRUpdateInterval = 5 * time.Minute
	}
	if config.MaxPortfolios == 0 {
		config.MaxPortfolios = 1000
	}

	// Create risk engine
	var riskEngine *RiskEngine
	if config.EnableRiskEngine {
		riskEngine = NewRiskEngine(logger, config.RiskConfig)
	}

	// Create VaR calculator
	var varCalculator *VaRCalculator
	if config.EnableVaRCalculation {
		varCalculator = NewVaRCalculator(logger, config.VaRConfig)
	}

	return &RiskManagementService{
		logger:            logger,
		config:            config,
		riskEngine:        riskEngine,
		varCalculator:     varCalculator,
		complianceManager: complianceManager,
		exchangeManager:   exchangeManager,
		portfolios:        make(map[uuid.UUID]*PortfolioRisk),
		globalMetrics:     &GlobalRiskMetrics{},
		circuitBreakers:   make(map[string]*CircuitBreaker),
		stopChan:          make(chan struct{}),
	}
}

// Start starts the risk management service
func (rms *RiskManagementService) Start(ctx context.Context) error {
	rms.mu.Lock()
	defer rms.mu.Unlock()

	if rms.isRunning {
		return fmt.Errorf("risk management service is already running")
	}

	rms.logger.Info(ctx, "Starting risk management service", map[string]interface{}{
		"enable_risk_engine":      rms.config.EnableRiskEngine,
		"enable_var_calculation":  rms.config.EnableVaRCalculation,
		"enable_compliance":       rms.config.EnableCompliance,
		"enable_circuit_breakers": rms.config.EnableCircuitBreakers,
		"update_interval":         rms.config.UpdateInterval,
	})

	rms.isRunning = true

	// Start risk engine
	if rms.config.EnableRiskEngine && rms.riskEngine != nil {
		if err := rms.riskEngine.Start(ctx); err != nil {
			return fmt.Errorf("failed to start risk engine: %w", err)
		}

		// Set global limits
		if rms.config.GlobalLimits != nil {
			rms.riskEngine.SetGlobalLimits(rms.config.GlobalLimits)
		}
	}

	// Start monitoring loops
	rms.wg.Add(1)
	go rms.monitorRisk(ctx)

	if rms.config.EnableVaRCalculation {
		rms.wg.Add(1)
		go rms.updateVaR(ctx)
	}

	rms.wg.Add(1)
	go rms.updateGlobalMetrics(ctx)

	rms.logger.Info(ctx, "Risk management service started", map[string]interface{}{
		"active_portfolios": len(rms.portfolios),
	})

	return nil
}

// Stop stops the risk management service
func (rms *RiskManagementService) Stop(ctx context.Context) error {
	rms.mu.Lock()
	defer rms.mu.Unlock()

	if !rms.isRunning {
		return fmt.Errorf("risk management service is not running")
	}

	rms.logger.Info(ctx, "Stopping risk management service", nil)

	// Stop risk engine
	if rms.riskEngine != nil {
		if err := rms.riskEngine.Stop(ctx); err != nil {
			rms.logger.Error(ctx, "Failed to stop risk engine", err, nil)
		}
	}

	close(rms.stopChan)
	rms.wg.Wait()

	rms.isRunning = false

	rms.logger.Info(ctx, "Risk management service stopped", nil)

	return nil
}

// CheckRisk performs comprehensive risk checks on a trading signal
func (rms *RiskManagementService) CheckRisk(ctx context.Context, request *RiskCheckRequest) (*RiskCheckResponse, error) {
	rms.totalRiskChecks++

	response := &RiskCheckResponse{
		Approved:         true,
		RiskScore:        decimal.NewFromInt(0),
		Violations:       make([]*RiskViolation, 0),
		Warnings:         make([]string, 0),
		CircuitBreakers:  make([]string, 0),
		ComplianceStatus: "compliant",
		Details:          make(map[string]interface{}),
		ProcessedAt:      time.Now(),
	}

	// Check risk engine limits
	if rms.config.EnableRiskEngine && rms.riskEngine != nil {
		if err := rms.riskEngine.CheckSignal(ctx, request.Signal); err != nil {
			response.Approved = false
			response.Violations = append(response.Violations, &RiskViolation{
				ID:            uuid.New(),
				ViolationType: RiskViolationTypeExposure,
				Message:       err.Error(),
				Timestamp:     time.Now(),
			})
			rms.riskViolations++
		}
	}

	// Check circuit breakers
	if rms.config.EnableCircuitBreakers {
		if breakers := rms.checkCircuitBreakers(ctx, request.Signal); len(breakers) > 0 {
			response.Approved = false
			response.CircuitBreakers = breakers
		}
	}

	// Check compliance rules
	if rms.config.EnableCompliance && rms.complianceManager != nil {
		// TODO: Integrate with compliance manager
		// complianceResult := rms.complianceManager.CheckCompliance(ctx, request)
		// response.ComplianceStatus = complianceResult.Status
	}

	// Calculate risk score
	response.RiskScore = rms.calculateRiskScore(ctx, request)

	// Add portfolio context
	if portfolio, exists := rms.portfolios[request.PortfolioID]; exists {
		response.Details["portfolio_value"] = portfolio.TotalValue.String()
		response.Details["portfolio_risk_score"] = portfolio.RiskScore.String()
	}

	rms.logger.Debug(ctx, "Risk check completed", map[string]interface{}{
		"signal_id":         request.Signal.ID.String(),
		"approved":          response.Approved,
		"risk_score":        response.RiskScore.String(),
		"violations":        len(response.Violations),
		"circuit_breakers":  len(response.CircuitBreakers),
		"compliance_status": response.ComplianceStatus,
	})

	return response, nil
}

// RegisterPortfolio registers a portfolio for risk monitoring
func (rms *RiskManagementService) RegisterPortfolio(ctx context.Context, portfolioID, strategyID uuid.UUID, name string) error {
	rms.mu.Lock()
	defer rms.mu.Unlock()

	if len(rms.portfolios) >= rms.config.MaxPortfolios {
		return fmt.Errorf("maximum number of portfolios reached: %d", rms.config.MaxPortfolios)
	}

	portfolio := &PortfolioRisk{
		ID:               portfolioID,
		StrategyID:       strategyID,
		Name:             name,
		TotalValue:       decimal.NewFromInt(0),
		Cash:             decimal.NewFromInt(0),
		Positions:        make(map[string]*PositionRisk),
		Metrics:          &RiskMetrics{},
		LastUpdated:      time.Now(),
		RiskScore:        decimal.NewFromInt(0),
		ComplianceStatus: "compliant",
	}

	rms.portfolios[portfolioID] = portfolio

	// Register with risk engine
	if rms.riskEngine != nil {
		monitor := &RiskMonitor{
			ID:       fmt.Sprintf("portfolio_%s", portfolioID.String()),
			Type:     RiskMonitorTypePortfolio,
			EntityID: portfolioID,
			Metrics:  portfolio.Metrics,
			Limits:   &RiskLimits{}, // TODO: Set appropriate limits
		}

		if err := rms.riskEngine.RegisterMonitor(monitor); err != nil {
			rms.logger.Error(ctx, "Failed to register portfolio monitor", err, map[string]interface{}{
				"portfolio_id": portfolioID.String(),
			})
		}
	}

	rms.logger.Info(ctx, "Portfolio registered for risk monitoring", map[string]interface{}{
		"portfolio_id": portfolioID.String(),
		"strategy_id":  strategyID.String(),
		"name":         name,
	})

	return nil
}

// UpdatePortfolio updates portfolio risk metrics
func (rms *RiskManagementService) UpdatePortfolio(ctx context.Context, portfolioID uuid.UUID, positions map[string]*PositionRisk, cash decimal.Decimal) error {
	rms.mu.Lock()
	defer rms.mu.Unlock()

	portfolio, exists := rms.portfolios[portfolioID]
	if !exists {
		return fmt.Errorf("portfolio not found: %s", portfolioID.String())
	}

	// Update positions
	portfolio.Positions = positions
	portfolio.Cash = cash

	// Calculate total value
	totalValue := cash
	for _, position := range positions {
		totalValue = totalValue.Add(position.MarketValue)
	}
	portfolio.TotalValue = totalValue

	// Update risk metrics
	rms.updatePortfolioMetrics(ctx, portfolio)

	portfolio.LastUpdated = time.Now()

	// Update risk engine metrics
	if rms.riskEngine != nil {
		monitorID := fmt.Sprintf("portfolio_%s", portfolioID.String())
		if err := rms.riskEngine.UpdateMetrics(monitorID, portfolio.Metrics); err != nil {
			rms.logger.Error(ctx, "Failed to update risk engine metrics", err, map[string]interface{}{
				"portfolio_id": portfolioID.String(),
			})
		}
	}

	return nil
}

// GetPortfolioRisk returns risk metrics for a portfolio
func (rms *RiskManagementService) GetPortfolioRisk(portfolioID uuid.UUID) (*PortfolioRisk, error) {
	rms.mu.RLock()
	defer rms.mu.RUnlock()

	portfolio, exists := rms.portfolios[portfolioID]
	if !exists {
		return nil, fmt.Errorf("portfolio not found: %s", portfolioID.String())
	}

	// Return a copy
	portfolioCopy := *portfolio
	return &portfolioCopy, nil
}

// GetGlobalMetrics returns global risk metrics
func (rms *RiskManagementService) GetGlobalMetrics() *GlobalRiskMetrics {
	rms.mu.RLock()
	defer rms.mu.RUnlock()

	// Return a copy
	metrics := *rms.globalMetrics
	return &metrics
}

// GetServiceMetrics returns service performance metrics
func (rms *RiskManagementService) GetServiceMetrics() *RiskServiceMetrics {
	rms.mu.RLock()
	defer rms.mu.RUnlock()

	return &RiskServiceMetrics{
		TotalRiskChecks:     rms.totalRiskChecks,
		RiskViolations:      rms.riskViolations,
		CircuitBreakerTrips: rms.circuitBreakerTrips,
		ActivePortfolios:    len(rms.portfolios),
		IsRunning:           rms.isRunning,
	}
}

// RiskServiceMetrics contains service performance metrics
type RiskServiceMetrics struct {
	TotalRiskChecks     int64 `json:"total_risk_checks"`
	RiskViolations      int64 `json:"risk_violations"`
	CircuitBreakerTrips int64 `json:"circuit_breaker_trips"`
	ActivePortfolios    int   `json:"active_portfolios"`
	IsRunning           bool  `json:"is_running"`
}

// Private methods

// monitorRisk continuously monitors portfolio risk
func (rms *RiskManagementService) monitorRisk(ctx context.Context) {
	defer rms.wg.Done()

	ticker := time.NewTicker(rms.config.UpdateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-rms.stopChan:
			return
		case <-ticker.C:
			rms.performRiskMonitoring(ctx)
		}
	}
}

// updateVaR periodically updates VaR calculations
func (rms *RiskManagementService) updateVaR(ctx context.Context) {
	defer rms.wg.Done()

	ticker := time.NewTicker(rms.config.VaRUpdateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-rms.stopChan:
			return
		case <-ticker.C:
			rms.performVaRCalculations(ctx)
		}
	}
}

// updateGlobalMetrics periodically updates global risk metrics
func (rms *RiskManagementService) updateGlobalMetrics(ctx context.Context) {
	defer rms.wg.Done()

	ticker := time.NewTicker(30 * time.Second) // Update every 30 seconds
	defer ticker.Stop()

	for {
		select {
		case <-rms.stopChan:
			return
		case <-ticker.C:
			rms.calculateGlobalMetrics(ctx)
		}
	}
}

// performRiskMonitoring performs risk monitoring on all portfolios
func (rms *RiskManagementService) performRiskMonitoring(ctx context.Context) {
	rms.mu.RLock()
	portfolios := make([]*PortfolioRisk, 0, len(rms.portfolios))
	for _, portfolio := range rms.portfolios {
		portfolios = append(portfolios, portfolio)
	}
	rms.mu.RUnlock()

	for _, portfolio := range portfolios {
		rms.updatePortfolioMetrics(ctx, portfolio)
	}
}

// performVaRCalculations performs VaR calculations for all portfolios
func (rms *RiskManagementService) performVaRCalculations(ctx context.Context) {
	if rms.varCalculator == nil {
		return
	}

	rms.mu.RLock()
	portfolios := make([]*PortfolioRisk, 0, len(rms.portfolios))
	for _, portfolio := range rms.portfolios {
		portfolios = append(portfolios, portfolio)
	}
	rms.mu.RUnlock()

	for _, portfolio := range portfolios {
		rms.calculatePortfolioVaR(ctx, portfolio)
	}
}

// updatePortfolioMetrics updates risk metrics for a portfolio
func (rms *RiskManagementService) updatePortfolioMetrics(ctx context.Context, portfolio *PortfolioRisk) {
	metrics := portfolio.Metrics

	// Calculate total exposure
	totalExposure := decimal.NewFromInt(0)
	for _, position := range portfolio.Positions {
		totalExposure = totalExposure.Add(position.MarketValue.Abs())
	}
	metrics.TotalExposure = totalExposure

	// Calculate net exposure
	netExposure := decimal.NewFromInt(0)
	for _, position := range portfolio.Positions {
		netExposure = netExposure.Add(position.MarketValue)
	}
	metrics.NetExposure = netExposure

	// Calculate unrealized P&L
	unrealizedPnL := decimal.NewFromInt(0)
	for _, position := range portfolio.Positions {
		unrealizedPnL = unrealizedPnL.Add(position.UnrealizedPnL)
	}
	metrics.UnrealizedPnL = unrealizedPnL

	// Calculate daily P&L
	dailyPnL := decimal.NewFromInt(0)
	for _, position := range portfolio.Positions {
		dailyPnL = dailyPnL.Add(position.DailyPnL)
	}
	metrics.DailyPnL = dailyPnL

	// Update position count
	metrics.PositionCount = len(portfolio.Positions)

	// Calculate leverage
	if portfolio.TotalValue.GreaterThan(decimal.NewFromInt(0)) {
		metrics.Leverage = totalExposure.Div(portfolio.TotalValue)
	}

	// Update cash balance
	metrics.CashBalance = portfolio.Cash

	// Calculate risk score
	portfolio.RiskScore = rms.calculatePortfolioRiskScore(portfolio)

	metrics.LastUpdated = time.Now()
}

// calculatePortfolioVaR calculates VaR for a portfolio
func (rms *RiskManagementService) calculatePortfolioVaR(ctx context.Context, portfolio *PortfolioRisk) {
	// TODO: Implement VaR calculation using historical portfolio data
	// This would require collecting historical portfolio values

	// For now, set a placeholder VaR based on portfolio value and volatility
	portfolioVaR := portfolio.TotalValue.Mul(decimal.NewFromFloat(0.05)) // 5% of portfolio value

	portfolio.VaRResult = &VaRResult{
		Method:          VaRMethodHistorical,
		ConfidenceLevel: decimal.NewFromFloat(0.95),
		TimeHorizon:     24 * time.Hour,
		VaR:             portfolioVaR,
		PortfolioValue:  portfolio.TotalValue,
		CalculatedAt:    time.Now(),
	}

	// Update metrics VaR
	portfolio.Metrics.VaR = portfolioVaR
}

// calculateGlobalMetrics calculates global risk metrics
func (rms *RiskManagementService) calculateGlobalMetrics(ctx context.Context) {
	rms.mu.Lock()
	defer rms.mu.Unlock()

	metrics := rms.globalMetrics

	// Reset counters
	totalValue := decimal.NewFromInt(0)
	totalExposure := decimal.NewFromInt(0)
	totalVaR := decimal.NewFromInt(0)
	totalRiskScore := decimal.NewFromInt(0)
	totalPositions := 0
	activePortfolios := 0

	// Aggregate metrics from all portfolios
	for _, portfolio := range rms.portfolios {
		if portfolio.TotalValue.GreaterThan(decimal.NewFromInt(0)) {
			activePortfolios++
		}

		totalValue = totalValue.Add(portfolio.TotalValue)
		totalExposure = totalExposure.Add(portfolio.Metrics.TotalExposure)
		totalVaR = totalVaR.Add(portfolio.Metrics.VaR)
		totalRiskScore = totalRiskScore.Add(portfolio.RiskScore)
		totalPositions += len(portfolio.Positions)
	}

	metrics.TotalPortfolioValue = totalValue
	metrics.TotalExposure = totalExposure
	metrics.TotalVaR = totalVaR
	metrics.ActivePortfolios = activePortfolios
	metrics.TotalPositions = totalPositions

	// Calculate average risk score
	if activePortfolios > 0 {
		metrics.AverageRiskScore = totalRiskScore.Div(decimal.NewFromInt(int64(activePortfolios)))
	}

	// Get additional metrics from risk engine
	if rms.riskEngine != nil {
		engineMetrics := rms.riskEngine.GetEngineMetrics()
		metrics.ComplianceViolations = int(engineMetrics.ViolationCount)
	}

	metrics.CircuitBreakerTrips = int(rms.circuitBreakerTrips)
	metrics.LastUpdated = time.Now()
}

// checkCircuitBreakers checks if any circuit breakers should trip
func (rms *RiskManagementService) checkCircuitBreakers(ctx context.Context, signal *framework.Signal) []string {
	var trippedBreakers []string

	for id, breaker := range rms.circuitBreakers {
		if breaker.State == CircuitStateOpen {
			trippedBreakers = append(trippedBreakers, id)
		}
	}

	return trippedBreakers
}

// calculateRiskScore calculates a risk score for a trading signal
func (rms *RiskManagementService) calculateRiskScore(ctx context.Context, request *RiskCheckRequest) decimal.Decimal {
	score := decimal.NewFromInt(0)

	// Base score from signal strength and confidence
	signalRisk := decimal.NewFromInt(1).Sub(request.Signal.Confidence)
	score = score.Add(signalRisk.Mul(decimal.NewFromFloat(30))) // Max 30 points

	// Add portfolio risk contribution
	if portfolio, exists := rms.portfolios[request.PortfolioID]; exists {
		portfolioRisk := portfolio.RiskScore.Mul(decimal.NewFromFloat(0.3))
		score = score.Add(portfolioRisk) // Max 30 points from portfolio
	}

	// Add position size risk
	notionalValue := request.Signal.Quantity.Mul(request.Signal.Price)
	if notionalValue.GreaterThan(decimal.NewFromFloat(10000)) {
		sizeRisk := decimal.NewFromFloat(20) // 20 points for large positions
		score = score.Add(sizeRisk)
	}

	// Add symbol volatility risk (placeholder)
	volatilityRisk := decimal.NewFromFloat(10) // 10 points base volatility risk
	score = score.Add(volatilityRisk)

	// Cap at 100
	if score.GreaterThan(decimal.NewFromInt(100)) {
		score = decimal.NewFromInt(100)
	}

	return score
}

// calculatePortfolioRiskScore calculates a risk score for a portfolio
func (rms *RiskManagementService) calculatePortfolioRiskScore(portfolio *PortfolioRisk) decimal.Decimal {
	score := decimal.NewFromInt(0)

	// Leverage risk
	if portfolio.Metrics.Leverage.GreaterThan(decimal.NewFromFloat(2)) {
		leverageRisk := portfolio.Metrics.Leverage.Mul(decimal.NewFromFloat(10))
		score = score.Add(leverageRisk)
	}

	// Concentration risk
	if len(portfolio.Positions) < 5 {
		concentrationRisk := decimal.NewFromFloat(20)
		score = score.Add(concentrationRisk)
	}

	// Drawdown risk
	if portfolio.Metrics.CurrentDrawdown.GreaterThan(decimal.NewFromFloat(0.1)) {
		drawdownRisk := portfolio.Metrics.CurrentDrawdown.Mul(decimal.NewFromFloat(100))
		score = score.Add(drawdownRisk)
	}

	// VaR risk
	if portfolio.VaRResult != nil && portfolio.TotalValue.GreaterThan(decimal.NewFromInt(0)) {
		varRatio := portfolio.VaRResult.VaR.Div(portfolio.TotalValue)
		if varRatio.GreaterThan(decimal.NewFromFloat(0.05)) {
			varRisk := varRatio.Mul(decimal.NewFromFloat(200))
			score = score.Add(varRisk)
		}
	}

	// Cap at 100
	if score.GreaterThan(decimal.NewFromInt(100)) {
		score = decimal.NewFromInt(100)
	}

	return score
}
