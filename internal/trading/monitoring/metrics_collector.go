package monitoring

import (
	"context"
	"runtime"
	"time"

	"github.com/ai-agentic-browser/internal/trading"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/shopspring/decimal"
)

// MetricsCollector collects various metrics for trading bots
type MetricsCollector struct {
	logger *observability.Logger
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector(logger *observability.Logger) *MetricsCollector {
	return &MetricsCollector{
		logger: logger,
	}
}

// CollectPortfolioMetrics collects portfolio-level metrics
func (mc *MetricsCollector) CollectPortfolioMetrics(ctx context.Context, bots []*trading.TradingBot) *PortfolioMetrics {
	totalValue := decimal.Zero
	totalPnL := decimal.Zero
	dailyPnL := decimal.Zero
	activeBots := 0
	strategyBreakdown := make(map[string]decimal.Decimal)
	assetAllocation := make(map[string]decimal.Decimal)

	for _, bot := range bots {
		if bot.State == trading.StateRunning {
			activeBots++
		}

		if bot.Performance != nil {
			totalPnL = totalPnL.Add(bot.Performance.NetProfit)
			// Simplified daily PnL calculation
			dailyPnL = dailyPnL.Add(bot.Performance.NetProfit.Div(decimal.NewFromInt(30))) // Approximate daily

			// Strategy breakdown
			strategy := string(bot.Strategy)
			if existing, exists := strategyBreakdown[strategy]; exists {
				strategyBreakdown[strategy] = existing.Add(bot.Performance.NetProfit)
			} else {
				strategyBreakdown[strategy] = bot.Performance.NetProfit
			}
		}

		// Asset allocation (simplified)
		for _, pair := range bot.Config.TradingPairs {
			if existing, exists := assetAllocation[pair]; exists {
				assetAllocation[pair] = existing.Add(decimal.NewFromFloat(1))
			} else {
				assetAllocation[pair] = decimal.NewFromFloat(1)
			}
		}
	}

	// Calculate total portfolio value (simplified)
	totalValue = decimal.NewFromFloat(100000).Add(totalPnL) // Base capital + PnL

	// Calculate portfolio metrics
	totalReturn := decimal.Zero
	if !totalValue.IsZero() {
		totalReturn = totalPnL.Div(totalValue.Sub(totalPnL)).Mul(decimal.NewFromFloat(100))
	}

	return &PortfolioMetrics{
		Timestamp:         time.Now(),
		TotalValue:        totalValue,
		TotalPnL:          totalPnL,
		DailyPnL:          dailyPnL,
		TotalReturn:       totalReturn,
		SharpeRatio:       mc.calculatePortfolioSharpe(bots),
		MaxDrawdown:       mc.calculatePortfolioMaxDrawdown(bots),
		VaR95:             mc.calculatePortfolioVaR(bots),
		ActiveBots:        activeBots,
		TotalBots:         len(bots),
		StrategyBreakdown: strategyBreakdown,
		AssetAllocation:   assetAllocation,
	}
}

// CollectSystemMetrics collects system-level metrics
func (mc *MetricsCollector) CollectSystemMetrics(ctx context.Context) *SystemMetrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return &SystemMetrics{
		Timestamp:      time.Now(),
		CPUUsage:       mc.getCPUUsage(),
		MemoryUsage:    int64(m.Alloc),
		DiskUsage:      mc.getDiskUsage(),
		NetworkIO:      mc.getNetworkIO(),
		GoroutineCount: runtime.NumGoroutine(),
		APIRequestRate: mc.getAPIRequestRate(),
		ErrorRate:      mc.getErrorRate(),
		ResponseTime:   mc.getResponseTime(),
	}
}

// CollectTradingMetrics collects trading execution metrics for a bot
func (mc *MetricsCollector) CollectTradingMetrics(ctx context.Context, bot *trading.TradingBot) *TradingMetrics {
	// Simplified trading metrics - in production this would come from actual trade data
	return &TradingMetrics{
		OrdersPlaced:      100,
		OrdersFilled:      95,
		OrdersCancelled:   3,
		OrdersFailed:      2,
		FillRate:          decimal.NewFromFloat(0.95),
		AvgExecutionTime:  500 * time.Millisecond,
		AvgSlippage:       decimal.NewFromFloat(0.001),
		TotalVolume:       decimal.NewFromFloat(50000),
		TotalFees:         decimal.NewFromFloat(25),
		LastTradeTime:     time.Now().Add(-5 * time.Minute),
	}
}

// CollectBotSystemMetrics collects system metrics for a specific bot
func (mc *MetricsCollector) CollectBotSystemMetrics(ctx context.Context, bot *trading.TradingBot) *BotSystemMetrics {
	return &BotSystemMetrics{
		CPUUsage:          mc.getBotCPUUsage(bot.ID),
		MemoryUsage:       mc.getBotMemoryUsage(bot.ID),
		GoroutineCount:    mc.getBotGoroutineCount(bot.ID),
		APICallsPerMinute: mc.getBotAPICallsPerMinute(bot.ID),
		ErrorCount:        mc.getBotErrorCount(bot.ID),
		LastErrorTime:     time.Now().Add(-1 * time.Hour),
		Uptime:            mc.getBotUptime(bot.ID),
	}
}

// Helper methods for metric calculations

func (mc *MetricsCollector) calculatePortfolioSharpe(bots []*trading.TradingBot) decimal.Decimal {
	// Simplified Sharpe ratio calculation
	totalSharpe := decimal.Zero
	count := 0

	for _, bot := range bots {
		if bot.Performance != nil && !bot.Performance.SharpeRatio.IsZero() {
			totalSharpe = totalSharpe.Add(bot.Performance.SharpeRatio)
			count++
		}
	}

	if count == 0 {
		return decimal.Zero
	}

	return totalSharpe.Div(decimal.NewFromInt(int64(count)))
}

func (mc *MetricsCollector) calculatePortfolioMaxDrawdown(bots []*trading.TradingBot) decimal.Decimal {
	maxDrawdown := decimal.Zero

	for _, bot := range bots {
		if bot.Performance != nil && bot.Performance.MaxDrawdown.GreaterThan(maxDrawdown) {
			maxDrawdown = bot.Performance.MaxDrawdown
		}
	}

	return maxDrawdown
}

func (mc *MetricsCollector) calculatePortfolioVaR(bots []*trading.TradingBot) decimal.Decimal {
	// Simplified VaR calculation
	totalVaR := decimal.Zero

	for _, bot := range bots {
		// Estimate VaR as 2% of current exposure
		if bot.Performance != nil {
			botVaR := bot.Performance.NetProfit.Abs().Mul(decimal.NewFromFloat(0.02))
			totalVaR = totalVaR.Add(botVaR)
		}
	}

	return totalVaR
}

// System metric helper methods (simplified implementations)

func (mc *MetricsCollector) getCPUUsage() float64 {
	// Simplified CPU usage - in production would use proper system monitoring
	return 25.5
}

func (mc *MetricsCollector) getDiskUsage() int64 {
	// Simplified disk usage
	return 1024 * 1024 * 1024 // 1GB
}

func (mc *MetricsCollector) getNetworkIO() int64 {
	// Simplified network I/O
	return 1024 * 1024 // 1MB
}

func (mc *MetricsCollector) getAPIRequestRate() float64 {
	// Simplified API request rate
	return 50.0 // 50 requests per second
}

func (mc *MetricsCollector) getErrorRate() float64 {
	// Simplified error rate
	return 1.5 // 1.5% error rate
}

func (mc *MetricsCollector) getResponseTime() time.Duration {
	// Simplified response time
	return 150 * time.Millisecond
}

// Bot-specific metric helper methods

func (mc *MetricsCollector) getBotCPUUsage(botID string) float64 {
	// Simplified bot CPU usage
	return 5.0
}

func (mc *MetricsCollector) getBotMemoryUsage(botID string) int64 {
	// Simplified bot memory usage
	return 50 * 1024 * 1024 // 50MB
}

func (mc *MetricsCollector) getBotGoroutineCount(botID string) int {
	// Simplified goroutine count
	return 10
}

func (mc *MetricsCollector) getBotAPICallsPerMinute(botID string) int {
	// Simplified API calls per minute
	return 60
}

func (mc *MetricsCollector) getBotErrorCount(botID string) int {
	// Simplified error count
	return 2
}

func (mc *MetricsCollector) getBotUptime(botID string) time.Duration {
	// Simplified uptime
	return 24 * time.Hour
}

// ConvertRiskMetrics converts risk manager metrics to monitoring metrics
func (mc *MetricsCollector) ConvertRiskMetrics(riskMetrics *trading.BotRiskMetrics) *BotRiskMetrics {
	return &BotRiskMetrics{
		VaR95:             riskMetrics.VaR95,
		VaR99:             riskMetrics.VaR95.Mul(decimal.NewFromFloat(1.3)), // Approximate VaR99
		ExpectedShortfall: riskMetrics.VaR95.Mul(decimal.NewFromFloat(1.2)), // Approximate ES
		RiskScore:         riskMetrics.RiskScore,
		ExposureRatio:     mc.calculateExposureRatio(riskMetrics),
		LeverageRatio:     decimal.NewFromFloat(1.0), // Simplified
		ConcentrationRisk: mc.calculateConcentrationRisk(riskMetrics),
		CorrelationRisk:   decimal.NewFromFloat(0.1), // Simplified
		LiquidityRisk:     decimal.NewFromFloat(0.05), // Simplified
	}
}

func (mc *MetricsCollector) calculateExposureRatio(riskMetrics *trading.BotRiskMetrics) decimal.Decimal {
	// Simplified exposure ratio calculation
	if riskMetrics.CurrentExposure.IsZero() {
		return decimal.Zero
	}
	return riskMetrics.CurrentExposure.Div(decimal.NewFromFloat(10000)) // Assume 10k base
}

func (mc *MetricsCollector) calculateConcentrationRisk(riskMetrics *trading.BotRiskMetrics) decimal.Decimal {
	// Simplified concentration risk
	return decimal.NewFromFloat(0.15) // 15% concentration risk
}
