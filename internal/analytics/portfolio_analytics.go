package analytics

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/ai-agentic-browser/internal/web3"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// PortfolioAnalytics provides comprehensive portfolio analysis and tracking
type PortfolioAnalytics struct {
	logger         *observability.Logger
	tradingEngine  *web3.TradingEngine
	dataRetention  time.Duration
	updateInterval time.Duration
	cache          map[uuid.UUID]*PortfolioMetrics
}

// PortfolioMetrics contains comprehensive portfolio performance metrics
type PortfolioMetrics struct {
	PortfolioID     uuid.UUID          `json:"portfolio_id"`
	UserID          uuid.UUID          `json:"user_id"`
	Name            string             `json:"name"`
	TotalValue      decimal.Decimal    `json:"total_value"`
	TotalPnL        decimal.Decimal    `json:"total_pnl"`
	TotalPnLPercent decimal.Decimal    `json:"total_pnl_percent"`
	DailyPnL        decimal.Decimal    `json:"daily_pnl"`
	WeeklyPnL       decimal.Decimal    `json:"weekly_pnl"`
	MonthlyPnL      decimal.Decimal    `json:"monthly_pnl"`
	MaxDrawdown     decimal.Decimal    `json:"max_drawdown"`
	SharpeRatio     decimal.Decimal    `json:"sharpe_ratio"`
	SortinoRatio    decimal.Decimal    `json:"sortino_ratio"`
	Volatility      decimal.Decimal    `json:"volatility"`
	Beta            decimal.Decimal    `json:"beta"`
	Alpha           decimal.Decimal    `json:"alpha"`
	Holdings        []HoldingMetrics   `json:"holdings"`
	Positions       []PositionMetrics  `json:"positions"`
	Performance     PerformanceHistory `json:"performance"`
	RiskMetrics     RiskAnalysis       `json:"risk_metrics"`
	LastUpdated     time.Time          `json:"last_updated"`
}

// HoldingMetrics represents metrics for individual holdings
type HoldingMetrics struct {
	Symbol       string          `json:"symbol"`
	TokenAddress string          `json:"token_address"`
	Amount       decimal.Decimal `json:"amount"`
	AveragePrice decimal.Decimal `json:"average_price"`
	CurrentPrice decimal.Decimal `json:"current_price"`
	Value        decimal.Decimal `json:"value"`
	PnL          decimal.Decimal `json:"pnl"`
	PnLPercent   decimal.Decimal `json:"pnl_percent"`
	Weight       decimal.Decimal `json:"weight"`
	DayChange    decimal.Decimal `json:"day_change"`
	WeekChange   decimal.Decimal `json:"week_change"`
	MonthChange  decimal.Decimal `json:"month_change"`
	LastUpdated  time.Time       `json:"last_updated"`
}

// PositionMetrics represents metrics for trading positions
type PositionMetrics struct {
	ID            uuid.UUID       `json:"id"`
	Symbol        string          `json:"symbol"`
	Strategy      string          `json:"strategy"`
	Side          string          `json:"side"`
	Amount        decimal.Decimal `json:"amount"`
	EntryPrice    decimal.Decimal `json:"entry_price"`
	CurrentPrice  decimal.Decimal `json:"current_price"`
	UnrealizedPnL decimal.Decimal `json:"unrealized_pnl"`
	RealizedPnL   decimal.Decimal `json:"realized_pnl"`
	StopLoss      decimal.Decimal `json:"stop_loss,omitempty"`
	TakeProfit    decimal.Decimal `json:"take_profit,omitempty"`
	Duration      time.Duration   `json:"duration"`
	Status        string          `json:"status"`
	OpenedAt      time.Time       `json:"opened_at"`
	LastUpdated   time.Time       `json:"last_updated"`
}

// PerformanceHistory tracks portfolio performance over time
type PerformanceHistory struct {
	Daily   []PerformancePoint `json:"daily"`
	Weekly  []PerformancePoint `json:"weekly"`
	Monthly []PerformancePoint `json:"monthly"`
}

// PerformancePoint represents a single performance data point
type PerformancePoint struct {
	Timestamp  time.Time       `json:"timestamp"`
	Value      decimal.Decimal `json:"value"`
	PnL        decimal.Decimal `json:"pnl"`
	PnLPercent decimal.Decimal `json:"pnl_percent"`
	Drawdown   decimal.Decimal `json:"drawdown"`
	Volume     decimal.Decimal `json:"volume"`
	Trades     int             `json:"trades"`
}

// RiskAnalysis provides comprehensive risk metrics
type RiskAnalysis struct {
	VaR95           decimal.Decimal `json:"var_95"`
	VaR99           decimal.Decimal `json:"var_99"`
	CVaR95          decimal.Decimal `json:"cvar_95"`
	MaxDrawdown     decimal.Decimal `json:"max_drawdown"`
	DownsideRisk    decimal.Decimal `json:"downside_risk"`
	UpsideCapture   decimal.Decimal `json:"upside_capture"`
	DownsideCapture decimal.Decimal `json:"downside_capture"`
	Correlation     decimal.Decimal `json:"correlation"`
	RiskScore       decimal.Decimal `json:"risk_score"`
	RiskGrade       string          `json:"risk_grade"`
}

// NewPortfolioAnalytics creates a new portfolio analytics service
func NewPortfolioAnalytics(logger *observability.Logger, tradingEngine *web3.TradingEngine) *PortfolioAnalytics {
	return &PortfolioAnalytics{
		logger:         logger,
		tradingEngine:  tradingEngine,
		dataRetention:  365 * 24 * time.Hour, // 1 year
		updateInterval: 5 * time.Minute,
		cache:          make(map[uuid.UUID]*PortfolioMetrics),
	}
}

// GetPortfolioMetrics returns comprehensive metrics for a portfolio
func (p *PortfolioAnalytics) GetPortfolioMetrics(ctx context.Context, portfolioID uuid.UUID) (*PortfolioMetrics, error) {
	// Check cache first
	if cached, exists := p.cache[portfolioID]; exists {
		if time.Since(cached.LastUpdated) < p.updateInterval {
			return cached, nil
		}
	}

	// Get portfolio from trading engine
	portfolio, err := p.tradingEngine.GetPortfolio(portfolioID)
	if err != nil {
		return nil, fmt.Errorf("failed to get portfolio: %w", err)
	}

	// Calculate comprehensive metrics
	metrics := &PortfolioMetrics{
		PortfolioID: portfolioID,
		UserID:      portfolio.UserID,
		Name:        portfolio.Name,
		TotalValue:  portfolio.TotalValue,
		TotalPnL:    portfolio.TotalPnL,
		LastUpdated: time.Now(),
	}

	// Calculate percentage P&L
	if portfolio.InvestedAmount.IsPositive() {
		metrics.TotalPnLPercent = portfolio.TotalPnL.Div(portfolio.InvestedAmount).Mul(decimal.NewFromInt(100))
	}

	// Calculate holdings metrics
	metrics.Holdings = p.calculateHoldingMetrics(portfolio)

	// Calculate position metrics
	metrics.Positions = p.calculatePositionMetrics(portfolio)

	// Calculate performance history
	metrics.Performance = p.calculatePerformanceHistory(portfolioID)

	// Calculate risk metrics
	metrics.RiskMetrics = p.calculateRiskMetrics(portfolioID, metrics.Performance)

	// Calculate advanced metrics
	p.calculateAdvancedMetrics(metrics)

	// Cache the results
	p.cache[portfolioID] = metrics

	p.logger.Info(ctx, "Portfolio metrics calculated", map[string]interface{}{
		"portfolio_id": portfolioID.String(),
		"total_value":  metrics.TotalValue.String(),
		"total_pnl":    metrics.TotalPnL.String(),
		"holdings":     len(metrics.Holdings),
		"positions":    len(metrics.Positions),
	})

	return metrics, nil
}

// calculateHoldingMetrics calculates metrics for individual holdings
func (p *PortfolioAnalytics) calculateHoldingMetrics(portfolio *web3.Portfolio) []HoldingMetrics {
	holdings := make([]HoldingMetrics, 0, len(portfolio.Holdings))

	for symbol, holding := range portfolio.Holdings {
		// Calculate current value
		currentValue := holding.Amount.Mul(holding.CurrentPrice)

		// Calculate P&L
		costBasis := holding.Amount.Mul(holding.AveragePrice)
		pnl := currentValue.Sub(costBasis)
		pnlPercent := decimal.Zero
		if costBasis.IsPositive() {
			pnlPercent = pnl.Div(costBasis).Mul(decimal.NewFromInt(100))
		}

		// Calculate weight in portfolio
		weight := decimal.Zero
		if portfolio.TotalValue.IsPositive() {
			weight = currentValue.Div(portfolio.TotalValue).Mul(decimal.NewFromInt(100))
		}

		holdingMetrics := HoldingMetrics{
			Symbol:       symbol,
			TokenAddress: holding.TokenAddress,
			Amount:       holding.Amount,
			AveragePrice: holding.AveragePrice,
			CurrentPrice: holding.CurrentPrice,
			Value:        currentValue,
			PnL:          pnl,
			PnLPercent:   pnlPercent,
			Weight:       weight,
			LastUpdated:  holding.LastUpdated,
		}

		// Calculate time-based changes (simplified - would need historical data)
		holdingMetrics.DayChange = p.calculatePriceChange(symbol, 24*time.Hour)
		holdingMetrics.WeekChange = p.calculatePriceChange(symbol, 7*24*time.Hour)
		holdingMetrics.MonthChange = p.calculatePriceChange(symbol, 30*24*time.Hour)

		holdings = append(holdings, holdingMetrics)
	}

	return holdings
}

// calculatePositionMetrics calculates metrics for trading positions
func (p *PortfolioAnalytics) calculatePositionMetrics(portfolio *web3.Portfolio) []PositionMetrics {
	positions := make([]PositionMetrics, 0, len(portfolio.ActivePositions))

	// For now, we'll create simplified position metrics since we only have position IDs
	// In a real implementation, we would fetch the actual Position objects from the trading engine
	for _, positionID := range portfolio.ActivePositions {
		// Create simplified position metrics
		positionMetrics := PositionMetrics{
			ID:            positionID,
			Symbol:        "ETH",      // Simplified - would get from actual position
			Strategy:      "momentum", // Simplified - would get from actual position
			Side:          "long",
			Amount:        decimal.NewFromFloat(1.0),    // Simplified
			EntryPrice:    decimal.NewFromFloat(2400.0), // Simplified
			CurrentPrice:  decimal.NewFromFloat(2450.0), // Simplified
			UnrealizedPnL: decimal.NewFromFloat(50.0),   // Simplified
			RealizedPnL:   decimal.Zero,
			StopLoss:      decimal.NewFromFloat(2160.0), // Simplified
			TakeProfit:    decimal.NewFromFloat(2880.0), // Simplified
			Duration:      time.Hour * 24,               // Simplified
			Status:        "open",
			OpenedAt:      time.Now().Add(-24 * time.Hour), // Simplified
			LastUpdated:   time.Now(),
		}

		positions = append(positions, positionMetrics)
	}

	return positions
}

// calculatePerformanceHistory calculates historical performance data
func (p *PortfolioAnalytics) calculatePerformanceHistory(portfolioID uuid.UUID) PerformanceHistory {
	// This would typically fetch from a time-series database
	// For now, we'll generate sample data

	now := time.Now()
	history := PerformanceHistory{
		Daily:   make([]PerformancePoint, 0, 30),
		Weekly:  make([]PerformancePoint, 0, 12),
		Monthly: make([]PerformancePoint, 0, 12),
	}

	// Generate daily data for last 30 days
	for i := 29; i >= 0; i-- {
		timestamp := now.AddDate(0, 0, -i)
		point := PerformancePoint{
			Timestamp:  timestamp,
			Value:      decimal.NewFromFloat(10000 + float64(i)*100 + float64(i%7)*50),
			PnL:        decimal.NewFromFloat(float64(i) * 10),
			PnLPercent: decimal.NewFromFloat(float64(i) * 0.1),
			Drawdown:   decimal.NewFromFloat(math.Max(0, float64(i%10)*-0.5)),
			Volume:     decimal.NewFromFloat(float64(i) * 1000),
			Trades:     i % 5,
		}
		history.Daily = append(history.Daily, point)
	}

	// Generate weekly data for last 12 weeks
	for i := 11; i >= 0; i-- {
		timestamp := now.AddDate(0, 0, -i*7)
		point := PerformancePoint{
			Timestamp:  timestamp,
			Value:      decimal.NewFromFloat(10000 + float64(i)*500),
			PnL:        decimal.NewFromFloat(float64(i) * 50),
			PnLPercent: decimal.NewFromFloat(float64(i) * 0.5),
			Drawdown:   decimal.NewFromFloat(math.Max(0, float64(i%4)*-2)),
			Volume:     decimal.NewFromFloat(float64(i) * 5000),
			Trades:     i * 3,
		}
		history.Weekly = append(history.Weekly, point)
	}

	// Generate monthly data for last 12 months
	for i := 11; i >= 0; i-- {
		timestamp := now.AddDate(0, -i, 0)
		point := PerformancePoint{
			Timestamp:  timestamp,
			Value:      decimal.NewFromFloat(10000 + float64(i)*1000),
			PnL:        decimal.NewFromFloat(float64(i) * 100),
			PnLPercent: decimal.NewFromFloat(float64(i) * 1.0),
			Drawdown:   decimal.NewFromFloat(math.Max(0, float64(i%3)*-5)),
			Volume:     decimal.NewFromFloat(float64(i) * 20000),
			Trades:     i * 10,
		}
		history.Monthly = append(history.Monthly, point)
	}

	return history
}

// calculateRiskMetrics calculates comprehensive risk analysis
func (p *PortfolioAnalytics) calculateRiskMetrics(portfolioID uuid.UUID, performance PerformanceHistory) RiskAnalysis {
	// Extract returns from daily performance
	returns := make([]float64, 0, len(performance.Daily))
	for i := 1; i < len(performance.Daily); i++ {
		if performance.Daily[i-1].Value.IsPositive() {
			returnPct := performance.Daily[i].Value.Sub(performance.Daily[i-1].Value).
				Div(performance.Daily[i-1].Value).InexactFloat64()
			returns = append(returns, returnPct)
		}
	}

	if len(returns) == 0 {
		return RiskAnalysis{}
	}

	// Calculate VaR (Value at Risk)
	sortedReturns := make([]float64, len(returns))
	copy(sortedReturns, returns)
	sort.Float64s(sortedReturns)

	var95Index := int(float64(len(sortedReturns)) * 0.05)
	var99Index := int(float64(len(sortedReturns)) * 0.01)

	var95 := decimal.NewFromFloat(sortedReturns[var95Index])
	var99 := decimal.NewFromFloat(sortedReturns[var99Index])

	// Calculate CVaR (Conditional VaR)
	cvar95Sum := 0.0
	cvar95Count := 0
	for i := 0; i <= var95Index; i++ {
		cvar95Sum += sortedReturns[i]
		cvar95Count++
	}
	cvar95 := decimal.NewFromFloat(cvar95Sum / float64(cvar95Count))

	// Calculate max drawdown
	maxDrawdown := decimal.Zero
	for _, point := range performance.Daily {
		if point.Drawdown.LessThan(maxDrawdown) {
			maxDrawdown = point.Drawdown
		}
	}

	// Calculate downside risk (simplified)
	downsideSum := 0.0
	downsideCount := 0
	for _, ret := range returns {
		if ret < 0 {
			downsideSum += ret * ret
			downsideCount++
		}
	}
	downsideRisk := decimal.Zero
	if downsideCount > 0 {
		downsideRisk = decimal.NewFromFloat(math.Sqrt(downsideSum / float64(downsideCount)))
	}

	// Calculate risk score (0-100)
	riskScore := decimal.NewFromFloat(50) // Default medium risk
	if var95.LessThan(decimal.NewFromFloat(-0.05)) {
		riskScore = decimal.NewFromFloat(80) // High risk
	} else if var95.GreaterThan(decimal.NewFromFloat(-0.02)) {
		riskScore = decimal.NewFromFloat(30) // Low risk
	}

	// Determine risk grade
	riskGrade := "B"
	if riskScore.LessThan(decimal.NewFromFloat(30)) {
		riskGrade = "A"
	} else if riskScore.GreaterThan(decimal.NewFromFloat(70)) {
		riskGrade = "D"
	} else if riskScore.GreaterThan(decimal.NewFromFloat(50)) {
		riskGrade = "C"
	}

	return RiskAnalysis{
		VaR95:           var95,
		VaR99:           var99,
		CVaR95:          cvar95,
		MaxDrawdown:     maxDrawdown,
		DownsideRisk:    downsideRisk,
		UpsideCapture:   decimal.NewFromFloat(1.05), // Simplified
		DownsideCapture: decimal.NewFromFloat(0.95), // Simplified
		Correlation:     decimal.NewFromFloat(0.75), // Simplified
		RiskScore:       riskScore,
		RiskGrade:       riskGrade,
	}
}

// calculateAdvancedMetrics calculates Sharpe ratio, Sortino ratio, etc.
func (p *PortfolioAnalytics) calculateAdvancedMetrics(metrics *PortfolioMetrics) {
	// Extract returns from daily performance
	returns := make([]float64, 0, len(metrics.Performance.Daily))
	for i := 1; i < len(metrics.Performance.Daily); i++ {
		if metrics.Performance.Daily[i-1].Value.IsPositive() {
			returnPct := metrics.Performance.Daily[i].Value.Sub(metrics.Performance.Daily[i-1].Value).
				Div(metrics.Performance.Daily[i-1].Value).InexactFloat64()
			returns = append(returns, returnPct)
		}
	}

	if len(returns) == 0 {
		return
	}

	// Calculate mean return
	meanReturn := 0.0
	for _, ret := range returns {
		meanReturn += ret
	}
	meanReturn /= float64(len(returns))

	// Calculate volatility (standard deviation)
	variance := 0.0
	for _, ret := range returns {
		variance += (ret - meanReturn) * (ret - meanReturn)
	}
	volatility := math.Sqrt(variance / float64(len(returns)))
	metrics.Volatility = decimal.NewFromFloat(volatility)

	// Calculate Sharpe ratio (assuming risk-free rate of 2% annually)
	riskFreeRate := 0.02 / 365 // Daily risk-free rate
	if volatility > 0 {
		sharpeRatio := (meanReturn - riskFreeRate) / volatility
		metrics.SharpeRatio = decimal.NewFromFloat(sharpeRatio)
	}

	// Calculate Sortino ratio (downside deviation)
	downsideVariance := 0.0
	downsideCount := 0
	for _, ret := range returns {
		if ret < riskFreeRate {
			downsideVariance += (ret - riskFreeRate) * (ret - riskFreeRate)
			downsideCount++
		}
	}
	if downsideCount > 0 {
		downsideDeviation := math.Sqrt(downsideVariance / float64(downsideCount))
		if downsideDeviation > 0 {
			sortinoRatio := (meanReturn - riskFreeRate) / downsideDeviation
			metrics.SortinoRatio = decimal.NewFromFloat(sortinoRatio)
		}
	}

	// Calculate Beta and Alpha (simplified - would need market benchmark)
	metrics.Beta = decimal.NewFromFloat(1.0)   // Simplified
	metrics.Alpha = decimal.NewFromFloat(0.02) // Simplified
}

// calculatePriceChange calculates price change over a time period (simplified)
func (p *PortfolioAnalytics) calculatePriceChange(symbol string, duration time.Duration) decimal.Decimal {
	// This would typically fetch historical price data
	// For now, return a random change
	change := float64((time.Now().Unix()+int64(len(symbol)))%20 - 10)
	return decimal.NewFromFloat(change)
}

// GetPortfolioComparison compares multiple portfolios
func (p *PortfolioAnalytics) GetPortfolioComparison(ctx context.Context, portfolioIDs []uuid.UUID) (map[uuid.UUID]*PortfolioMetrics, error) {
	comparison := make(map[uuid.UUID]*PortfolioMetrics)

	for _, portfolioID := range portfolioIDs {
		metrics, err := p.GetPortfolioMetrics(ctx, portfolioID)
		if err != nil {
			p.logger.Error(ctx, "Failed to get portfolio metrics for comparison", err, map[string]interface{}{
				"portfolio_id": portfolioID.String(),
			})
			continue
		}
		comparison[portfolioID] = metrics
	}

	return comparison, nil
}
