package hft

import (
	"math"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ExecutionAnalyzer analyzes execution quality metrics
type ExecutionAnalyzer struct {
	logger *observability.Logger
	config PostTradeConfig
}

// SlippageCalculator calculates various types of slippage
type SlippageCalculator struct {
	logger *observability.Logger
	config PostTradeConfig
}

// MarketImpactTracker tracks and measures market impact
type MarketImpactTracker struct {
	logger *observability.Logger
	config PostTradeConfig
}

// PerformanceEngine calculates performance metrics and attribution
type PerformanceEngine struct {
	logger *observability.Logger
	config PostTradeConfig
}

// ReportGenerator generates various types of reports
type ReportGenerator struct {
	logger *observability.Logger
	config PostTradeConfig
}

// NewExecutionAnalyzer creates a new execution analyzer
func NewExecutionAnalyzer(logger *observability.Logger, config PostTradeConfig) *ExecutionAnalyzer {
	return &ExecutionAnalyzer{
		logger: logger,
		config: config,
	}
}

// NewSlippageCalculator creates a new slippage calculator
func NewSlippageCalculator(logger *observability.Logger, config PostTradeConfig) *SlippageCalculator {
	return &SlippageCalculator{
		logger: logger,
		config: config,
	}
}

// NewMarketImpactTracker creates a new market impact tracker
func NewMarketImpactTracker(logger *observability.Logger, config PostTradeConfig) *MarketImpactTracker {
	return &MarketImpactTracker{
		logger: logger,
		config: config,
	}
}

// NewPerformanceEngine creates a new performance engine
func NewPerformanceEngine(logger *observability.Logger, config PostTradeConfig) *PerformanceEngine {
	return &PerformanceEngine{
		logger: logger,
		config: config,
	}
}

// NewReportGenerator creates a new report generator
func NewReportGenerator(logger *observability.Logger, config PostTradeConfig) *ReportGenerator {
	return &ReportGenerator{
		logger: logger,
		config: config,
	}
}

// CalculateSlippage calculates slippage for a trade
func (sc *SlippageCalculator) CalculateSlippage(trade *TradeRecord) float64 {
	// Implementation slippage vs arrival price
	if trade.ArrivalPrice.IsZero() {
		return 0.0
	}

	priceDiff := trade.Price.Sub(trade.ArrivalPrice)
	if trade.Side == OrderSideSell {
		priceDiff = priceDiff.Neg()
	}

	slippage := priceDiff.Div(trade.ArrivalPrice).InexactFloat64() * 10000 // Convert to basis points
	return slippage
}

// CalculateVWAPSlippage calculates slippage vs VWAP
func (sc *SlippageCalculator) CalculateVWAPSlippage(trade *TradeRecord) float64 {
	if trade.VWAPPrice.IsZero() {
		return 0.0
	}

	priceDiff := trade.Price.Sub(trade.VWAPPrice)
	if trade.Side == OrderSideSell {
		priceDiff = priceDiff.Neg()
	}

	slippage := priceDiff.Div(trade.VWAPPrice).InexactFloat64() * 10000 // Convert to basis points
	return slippage
}

// CalculateTWAPSlippage calculates slippage vs TWAP
func (sc *SlippageCalculator) CalculateTWAPSlippage(trade *TradeRecord) float64 {
	if trade.TWAPPrice.IsZero() {
		return 0.0
	}

	priceDiff := trade.Price.Sub(trade.TWAPPrice)
	if trade.Side == OrderSideSell {
		priceDiff = priceDiff.Neg()
	}

	slippage := priceDiff.Div(trade.TWAPPrice).InexactFloat64() * 10000 // Convert to basis points
	return slippage
}

// CalculateImpact calculates market impact for a trade
func (mit *MarketImpactTracker) CalculateImpact(trade *TradeRecord) float64 {
	// Simplified market impact calculation
	// In production, this would use sophisticated models considering:
	// - Order size relative to average daily volume
	// - Market volatility
	// - Liquidity conditions
	// - Time of day effects

	// Base impact based on trade size (simplified)
	baseImpact := trade.Quantity.InexactFloat64() / 1000.0 * 2.0 // 2 bps per 1000 units

	// Adjust for market conditions (simplified)
	volatilityMultiplier := 1.0 // Would use actual volatility
	liquidityMultiplier := 1.0  // Would use actual liquidity

	impact := baseImpact * volatilityMultiplier * liquidityMultiplier

	// Cap at reasonable maximum
	if impact > 50.0 {
		impact = 50.0
	}

	return impact
}

// CalculateTemporaryImpact calculates temporary market impact
func (mit *MarketImpactTracker) CalculateTemporaryImpact(trade *TradeRecord) float64 {
	// Temporary impact that reverts quickly
	permanentImpact := mit.CalculateImpact(trade)
	temporaryImpact := permanentImpact * 0.6 // Assume 60% is temporary

	return temporaryImpact
}

// CalculatePermanentImpact calculates permanent market impact
func (mit *MarketImpactTracker) CalculatePermanentImpact(trade *TradeRecord) float64 {
	// Permanent impact that doesn't revert
	totalImpact := mit.CalculateImpact(trade)
	permanentImpact := totalImpact * 0.4 // Assume 40% is permanent

	return permanentImpact
}

// AnalyzeExecution analyzes execution quality for a trade
func (ea *ExecutionAnalyzer) AnalyzeExecution(trade *TradeRecord) *ExecutionAnalysis {
	analysis := &ExecutionAnalysis{
		TradeID:       trade.TradeID,
		Symbol:        trade.Symbol,
		ExecutionTime: trade.ExecutionTime,
	}

	// Calculate execution scores
	analysis.SlippageScore = ea.calculateSlippageScore(trade.Slippage)
	analysis.TimingScore = ea.calculateTimingScore(trade.TimingCost)
	analysis.CostScore = ea.calculateCostScore(trade.TotalCosts, trade.Value)
	analysis.ImpactScore = ea.calculateImpactScore(trade.MarketImpact)

	// Overall execution score (weighted average)
	analysis.OverallScore = (analysis.SlippageScore*0.3 +
		analysis.TimingScore*0.25 +
		analysis.CostScore*0.25 +
		analysis.ImpactScore*0.2)

	// Determine execution quality
	if analysis.OverallScore >= 90 {
		analysis.Quality = "EXCELLENT"
	} else if analysis.OverallScore >= 75 {
		analysis.Quality = "GOOD"
	} else if analysis.OverallScore >= 60 {
		analysis.Quality = "FAIR"
	} else {
		analysis.Quality = "POOR"
	}

	return analysis
}

// ExecutionAnalysis contains execution analysis results
type ExecutionAnalysis struct {
	TradeID       uuid.UUID `json:"trade_id"`
	Symbol        string    `json:"symbol"`
	ExecutionTime time.Time `json:"execution_time"`

	// Scores (0-100)
	SlippageScore float64 `json:"slippage_score"`
	TimingScore   float64 `json:"timing_score"`
	CostScore     float64 `json:"cost_score"`
	ImpactScore   float64 `json:"impact_score"`
	OverallScore  float64 `json:"overall_score"`

	// Quality assessment
	Quality         string   `json:"quality"`
	Recommendations []string `json:"recommendations,omitempty"`
}

// calculateSlippageScore calculates a score based on slippage
func (ea *ExecutionAnalyzer) calculateSlippageScore(slippage float64) float64 {
	// Score decreases as slippage increases
	// 0 bps = 100 score, 20 bps = 0 score
	score := 100.0 - (math.Abs(slippage) * 5.0)
	if score < 0 {
		score = 0
	}
	return score
}

// calculateTimingScore calculates a score based on timing cost
func (ea *ExecutionAnalyzer) calculateTimingScore(timingCost float64) float64 {
	// Score decreases as timing cost increases
	score := 100.0 - (math.Abs(timingCost) * 3.0)
	if score < 0 {
		score = 0
	}
	return score
}

// calculateCostScore calculates a score based on transaction costs
func (ea *ExecutionAnalyzer) calculateCostScore(totalCosts, tradeValue decimal.Decimal) float64 {
	if tradeValue.IsZero() {
		return 100.0
	}

	costRatio := totalCosts.Div(tradeValue).InexactFloat64() * 10000 // Convert to basis points
	score := 100.0 - (costRatio * 10.0)                              // 10 bps = 0 score
	if score < 0 {
		score = 0
	}
	return score
}

// calculateImpactScore calculates a score based on market impact
func (ea *ExecutionAnalyzer) calculateImpactScore(marketImpact float64) float64 {
	// Score decreases as market impact increases
	score := 100.0 - (math.Abs(marketImpact) * 4.0)
	if score < 0 {
		score = 0
	}
	return score
}

// CalculatePerformance calculates performance metrics for a strategy
func (pe *PerformanceEngine) CalculatePerformance(trades []TradeRecord, strategyID string) *PerformanceData {
	if len(trades) == 0 {
		return &PerformanceData{StrategyID: strategyID}
	}

	// Filter trades for this strategy
	strategyTrades := make([]TradeRecord, 0)
	for _, trade := range trades {
		if trade.StrategyID == strategyID {
			strategyTrades = append(strategyTrades, trade)
		}
	}

	if len(strategyTrades) == 0 {
		return &PerformanceData{StrategyID: strategyID}
	}

	perf := &PerformanceData{
		StrategyID:  strategyID,
		Period:      time.Now(),
		TotalTrades: len(strategyTrades),
	}

	// Calculate P&L metrics
	var grossPnL, totalCosts decimal.Decimal
	var returns []float64
	var wins, losses int
	var totalWin, totalLoss decimal.Decimal

	for _, trade := range strategyTrades {
		// Calculate trade P&L (simplified)
		tradePnL := trade.Value.Sub(trade.TotalCosts)
		grossPnL = grossPnL.Add(tradePnL)
		totalCosts = totalCosts.Add(trade.TotalCosts)

		// Track returns for volatility calculation
		if !trade.ArrivalPrice.IsZero() {
			ret := trade.Price.Sub(trade.ArrivalPrice).Div(trade.ArrivalPrice).InexactFloat64()
			returns = append(returns, ret)

			if ret > 0 {
				wins++
				totalWin = totalWin.Add(decimal.NewFromFloat(ret))
			} else if ret < 0 {
				losses++
				totalLoss = totalLoss.Add(decimal.NewFromFloat(math.Abs(ret)))
			}
		}
	}

	perf.GrossPnL = grossPnL
	perf.TotalCosts = totalCosts
	perf.NetPnL = grossPnL.Sub(totalCosts)

	// Calculate risk metrics
	if len(returns) > 1 {
		perf.Volatility = pe.calculateVolatility(returns)
		perf.SharpeRatio = pe.calculateSharpeRatio(returns)
		perf.SortinoRatio = pe.calculateSortinoRatio(returns)
	}

	// Calculate win rate and profit factor
	if perf.TotalTrades > 0 {
		perf.WinRate = float64(wins) / float64(perf.TotalTrades) * 100.0
	}

	if wins > 0 {
		perf.AvgWin = totalWin.Div(decimal.NewFromInt(int64(wins)))
	}
	if losses > 0 {
		perf.AvgLoss = totalLoss.Div(decimal.NewFromInt(int64(losses)))
	}

	if !perf.AvgLoss.IsZero() {
		perf.ProfitFactor = perf.AvgWin.Div(perf.AvgLoss).InexactFloat64()
	}

	return perf
}

// calculateVolatility calculates the volatility of returns
func (pe *PerformanceEngine) calculateVolatility(returns []float64) float64 {
	if len(returns) < 2 {
		return 0.0
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

	// Return standard deviation (volatility)
	return math.Sqrt(variance)
}

// calculateSharpeRatio calculates the Sharpe ratio
func (pe *PerformanceEngine) calculateSharpeRatio(returns []float64) float64 {
	if len(returns) < 2 {
		return 0.0
	}

	volatility := pe.calculateVolatility(returns)
	if volatility == 0 {
		return 0.0
	}

	// Calculate mean return
	var sum float64
	for _, ret := range returns {
		sum += ret
	}
	meanReturn := sum / float64(len(returns))

	// Assume risk-free rate of 0 for simplicity
	return meanReturn / volatility
}

// calculateSortinoRatio calculates the Sortino ratio
func (pe *PerformanceEngine) calculateSortinoRatio(returns []float64) float64 {
	if len(returns) < 2 {
		return 0.0
	}

	// Calculate mean return
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
			downsideVariance += math.Pow(ret, 2)
			downsideCount++
		}
	}

	if downsideCount == 0 {
		return math.Inf(1) // No downside risk
	}

	downsideDeviation := math.Sqrt(downsideVariance / float64(downsideCount))
	if downsideDeviation == 0 {
		return 0.0
	}

	return meanReturn / downsideDeviation
}

// GenerateDailyReport generates a daily execution report
func (rg *ReportGenerator) GenerateDailyReport(trades []TradeRecord, metrics map[string]*ExecutionMetrics) *DailyReport {
	report := &DailyReport{
		ID:          uuid.New(),
		Date:        time.Now(),
		GeneratedAt: time.Now(),
		Type:        "DAILY_EXECUTION",
	}

	// Calculate summary metrics
	report.Summary = rg.calculateDailySummary(trades, metrics)

	// Generate symbol breakdown
	report.SymbolBreakdown = rg.generateSymbolBreakdown(metrics)

	// Generate execution quality analysis
	report.ExecutionQuality = rg.analyzeExecutionQuality(trades)

	return report
}

// DailyReport represents a daily execution report
type DailyReport struct {
	ID               uuid.UUID                    `json:"id"`
	Date             time.Time                    `json:"date"`
	GeneratedAt      time.Time                    `json:"generated_at"`
	Type             string                       `json:"type"`
	Summary          *DailySummary                `json:"summary"`
	SymbolBreakdown  map[string]*ExecutionMetrics `json:"symbol_breakdown"`
	ExecutionQuality *ExecutionQualitySummary     `json:"execution_quality"`
}

// DailySummary contains daily summary metrics
type DailySummary struct {
	TotalTrades      int             `json:"total_trades"`
	TotalVolume      decimal.Decimal `json:"total_volume"`
	AvgSlippage      float64         `json:"avg_slippage"`
	AvgMarketImpact  float64         `json:"avg_market_impact"`
	TotalCosts       decimal.Decimal `json:"total_costs"`
	SymbolsTraded    int             `json:"symbols_traded"`
	StrategiesActive int             `json:"strategies_active"`
}

// ExecutionQualitySummary contains execution quality analysis
type ExecutionQualitySummary struct {
	ExcellentTrades int      `json:"excellent_trades"`
	GoodTrades      int      `json:"good_trades"`
	FairTrades      int      `json:"fair_trades"`
	PoorTrades      int      `json:"poor_trades"`
	OverallScore    float64  `json:"overall_score"`
	Recommendations []string `json:"recommendations"`
}

// calculateDailySummary calculates daily summary metrics
func (rg *ReportGenerator) calculateDailySummary(trades []TradeRecord, metrics map[string]*ExecutionMetrics) *DailySummary {
	summary := &DailySummary{
		TotalTrades:   len(trades),
		SymbolsTraded: len(metrics),
	}

	var totalVolume, totalCosts decimal.Decimal
	var totalSlippage, totalMarketImpact float64
	strategiesMap := make(map[string]bool)

	for _, trade := range trades {
		totalVolume = totalVolume.Add(trade.Value)
		totalCosts = totalCosts.Add(trade.TotalCosts)
		totalSlippage += trade.Slippage
		totalMarketImpact += trade.MarketImpact

		if trade.StrategyID != "" {
			strategiesMap[trade.StrategyID] = true
		}
	}

	summary.TotalVolume = totalVolume
	summary.TotalCosts = totalCosts
	summary.StrategiesActive = len(strategiesMap)

	if len(trades) > 0 {
		summary.AvgSlippage = totalSlippage / float64(len(trades))
		summary.AvgMarketImpact = totalMarketImpact / float64(len(trades))
	}

	return summary
}

// generateSymbolBreakdown generates symbol-level breakdown
func (rg *ReportGenerator) generateSymbolBreakdown(metrics map[string]*ExecutionMetrics) map[string]*ExecutionMetrics {
	// Return a copy of the metrics
	breakdown := make(map[string]*ExecutionMetrics)
	for symbol, metric := range metrics {
		breakdown[symbol] = metric
	}
	return breakdown
}

// analyzeExecutionQuality analyzes overall execution quality
func (rg *ReportGenerator) analyzeExecutionQuality(trades []TradeRecord) *ExecutionQualitySummary {
	quality := &ExecutionQualitySummary{}

	var totalScore float64

	for _, trade := range trades {
		// Simplified quality scoring
		score := 100.0 - math.Abs(trade.Slippage)*5.0 - math.Abs(trade.MarketImpact)*3.0
		if score < 0 {
			score = 0
		}

		totalScore += score

		if score >= 90 {
			quality.ExcellentTrades++
		} else if score >= 75 {
			quality.GoodTrades++
		} else if score >= 60 {
			quality.FairTrades++
		} else {
			quality.PoorTrades++
		}
	}

	if len(trades) > 0 {
		quality.OverallScore = totalScore / float64(len(trades))
	}

	// Generate recommendations
	if quality.OverallScore < 70 {
		quality.Recommendations = append(quality.Recommendations, "Consider reviewing order routing algorithms")
		quality.Recommendations = append(quality.Recommendations, "Analyze market timing strategies")
	}
	if quality.PoorTrades > len(trades)/10 {
		quality.Recommendations = append(quality.Recommendations, "Investigate high-slippage trades")
	}

	return quality
}
