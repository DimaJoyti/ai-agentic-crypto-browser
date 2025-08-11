package monitoring

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/shopspring/decimal"
)

// DashboardManager manages dashboard data and updates
type DashboardManager struct {
	logger        *observability.Logger
	dashboardData *DashboardData
	mu            sync.RWMutex
}

// DashboardData represents the complete dashboard data structure
type DashboardData struct {
	LastUpdated      time.Time                  `json:"last_updated"`
	Overview         *DashboardOverview         `json:"overview"`
	BotSummaries     []*BotSummary             `json:"bot_summaries"`
	PortfolioMetrics *PortfolioMetrics         `json:"portfolio_metrics"`
	SystemStatus     *SystemStatus             `json:"system_status"`
	RecentAlerts     []*Alert                  `json:"recent_alerts"`
	PerformanceChart *PerformanceChartData     `json:"performance_chart"`
	RiskMetrics      *RiskDashboardMetrics     `json:"risk_metrics"`
	TradingActivity  *TradingActivitySummary   `json:"trading_activity"`
}

// DashboardOverview provides high-level overview metrics
type DashboardOverview struct {
	TotalBots        int             `json:"total_bots"`
	ActiveBots       int             `json:"active_bots"`
	PausedBots       int             `json:"paused_bots"`
	ErrorBots        int             `json:"error_bots"`
	TotalValue       decimal.Decimal `json:"total_value"`
	DailyPnL         decimal.Decimal `json:"daily_pnl"`
	TotalReturn      decimal.Decimal `json:"total_return"`
	ActiveAlerts     int             `json:"active_alerts"`
	CriticalAlerts   int             `json:"critical_alerts"`
	SystemHealth     HealthStatus    `json:"system_health"`
}

// BotSummary provides summary information for each bot
type BotSummary struct {
	BotID           string          `json:"bot_id"`
	Name            string          `json:"name"`
	Strategy        string          `json:"strategy"`
	State           string          `json:"state"`
	Health          HealthStatus    `json:"health"`
	PnL             decimal.Decimal `json:"pnl"`
	DailyPnL        decimal.Decimal `json:"daily_pnl"`
	WinRate         decimal.Decimal `json:"win_rate"`
	TotalTrades     int             `json:"total_trades"`
	LastTradeTime   time.Time       `json:"last_trade_time"`
	RiskScore       int             `json:"risk_score"`
	ActiveAlerts    int             `json:"active_alerts"`
	Uptime          time.Duration   `json:"uptime"`
}

// SystemStatus provides system health and performance information
type SystemStatus struct {
	OverallHealth   HealthStatus  `json:"overall_health"`
	CPUUsage        float64       `json:"cpu_usage"`
	MemoryUsage     int64         `json:"memory_usage"`
	DiskUsage       int64         `json:"disk_usage"`
	NetworkIO       int64         `json:"network_io"`
	APIRequestRate  float64       `json:"api_request_rate"`
	ErrorRate       float64       `json:"error_rate"`
	ResponseTime    time.Duration `json:"response_time"`
	Uptime          time.Duration `json:"uptime"`
	LastHealthCheck time.Time     `json:"last_health_check"`
}

// PerformanceChartData provides data for performance charts
type PerformanceChartData struct {
	TimeRange    string                    `json:"time_range"`
	DataPoints   []*PerformanceDataPoint   `json:"data_points"`
	BotSeries    map[string][]*DataPoint   `json:"bot_series"`
	Benchmarks   []*BenchmarkDataPoint     `json:"benchmarks"`
}

// PerformanceDataPoint represents a single point in the performance chart
type PerformanceDataPoint struct {
	Timestamp       time.Time       `json:"timestamp"`
	PortfolioValue  decimal.Decimal `json:"portfolio_value"`
	TotalPnL        decimal.Decimal `json:"total_pnl"`
	DailyPnL        decimal.Decimal `json:"daily_pnl"`
	DrawdownPercent decimal.Decimal `json:"drawdown_percent"`
}

// DataPoint represents a generic data point for charts
type DataPoint struct {
	Timestamp time.Time       `json:"timestamp"`
	Value     decimal.Decimal `json:"value"`
}

// BenchmarkDataPoint represents benchmark comparison data
type BenchmarkDataPoint struct {
	Timestamp time.Time       `json:"timestamp"`
	BTC       decimal.Decimal `json:"btc"`
	ETH       decimal.Decimal `json:"eth"`
	SPY       decimal.Decimal `json:"spy"`
}

// RiskDashboardMetrics provides risk metrics for the dashboard
type RiskDashboardMetrics struct {
	PortfolioVaR95    decimal.Decimal            `json:"portfolio_var_95"`
	PortfolioVaR99    decimal.Decimal            `json:"portfolio_var_99"`
	MaxDrawdown       decimal.Decimal            `json:"max_drawdown"`
	CurrentDrawdown   decimal.Decimal            `json:"current_drawdown"`
	SharpeRatio       decimal.Decimal            `json:"sharpe_ratio"`
	ConcentrationRisk map[string]decimal.Decimal `json:"concentration_risk"`
	CorrelationMatrix map[string]map[string]decimal.Decimal `json:"correlation_matrix"`
	RiskDistribution  map[string]int             `json:"risk_distribution"`
}

// TradingActivitySummary provides trading activity summary
type TradingActivitySummary struct {
	TotalOrders24h    int             `json:"total_orders_24h"`
	FilledOrders24h   int             `json:"filled_orders_24h"`
	CancelledOrders24h int            `json:"cancelled_orders_24h"`
	FailedOrders24h   int             `json:"failed_orders_24h"`
	TotalVolume24h    decimal.Decimal `json:"total_volume_24h"`
	TotalFees24h      decimal.Decimal `json:"total_fees_24h"`
	AvgFillRate       decimal.Decimal `json:"avg_fill_rate"`
	AvgSlippage       decimal.Decimal `json:"avg_slippage"`
	TopTradingPairs   []*TradingPairSummary `json:"top_trading_pairs"`
}

// TradingPairSummary provides summary for a trading pair
type TradingPairSummary struct {
	Symbol      string          `json:"symbol"`
	Volume24h   decimal.Decimal `json:"volume_24h"`
	Trades24h   int             `json:"trades_24h"`
	PnL24h      decimal.Decimal `json:"pnl_24h"`
	LastPrice   decimal.Decimal `json:"last_price"`
	Change24h   decimal.Decimal `json:"change_24h"`
}

// NewDashboardManager creates a new dashboard manager
func NewDashboardManager(logger *observability.Logger) *DashboardManager {
	return &DashboardManager{
		logger: logger,
		dashboardData: &DashboardData{
			LastUpdated:      time.Now(),
			Overview:         &DashboardOverview{},
			BotSummaries:     make([]*BotSummary, 0),
			PortfolioMetrics: &PortfolioMetrics{},
			SystemStatus:     &SystemStatus{},
			RecentAlerts:     make([]*Alert, 0),
			PerformanceChart: &PerformanceChartData{},
			RiskMetrics:      &RiskDashboardMetrics{},
			TradingActivity:  &TradingActivitySummary{},
		},
	}
}

// UpdateDashboard updates the dashboard with latest data
func (dm *DashboardManager) UpdateDashboard(
	ctx context.Context,
	botMetrics map[string]*BotMetrics,
	portfolioMetrics *PortfolioMetrics,
	systemMetrics *SystemMetrics,
	alerts []*Alert,
) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	dm.dashboardData.LastUpdated = time.Now()

	// Update overview
	dm.updateOverview(botMetrics, portfolioMetrics, alerts)

	// Update bot summaries
	dm.updateBotSummaries(botMetrics)

	// Update portfolio metrics
	dm.dashboardData.PortfolioMetrics = portfolioMetrics

	// Update system status
	dm.updateSystemStatus(systemMetrics)

	// Update recent alerts
	dm.updateRecentAlerts(alerts)

	// Update performance chart
	dm.updatePerformanceChart(portfolioMetrics)

	// Update risk metrics
	dm.updateRiskMetrics(botMetrics, portfolioMetrics)

	// Update trading activity
	dm.updateTradingActivity(botMetrics)

	dm.logger.Debug(ctx, "Dashboard updated", map[string]interface{}{
		"total_bots":     dm.dashboardData.Overview.TotalBots,
		"active_bots":    dm.dashboardData.Overview.ActiveBots,
		"active_alerts":  dm.dashboardData.Overview.ActiveAlerts,
		"portfolio_value": dm.dashboardData.Overview.TotalValue.String(),
	})
}

// updateOverview updates the dashboard overview section
func (dm *DashboardManager) updateOverview(botMetrics map[string]*BotMetrics, portfolioMetrics *PortfolioMetrics, alerts []*Alert) {
	overview := dm.dashboardData.Overview

	// Count bots by state
	overview.TotalBots = len(botMetrics)
	overview.ActiveBots = 0
	overview.PausedBots = 0
	overview.ErrorBots = 0

	for _, metrics := range botMetrics {
		switch metrics.State {
		case "running":
			overview.ActiveBots++
		case "paused":
			overview.PausedBots++
		case "error":
			overview.ErrorBots++
		}
	}

	// Portfolio metrics
	overview.TotalValue = portfolioMetrics.TotalValue
	overview.DailyPnL = portfolioMetrics.DailyPnL
	overview.TotalReturn = portfolioMetrics.TotalReturn

	// Alert counts
	overview.ActiveAlerts = 0
	overview.CriticalAlerts = 0
	for _, alert := range alerts {
		if !alert.Resolved {
			overview.ActiveAlerts++
			if alert.Severity == AlertSeverityCritical {
				overview.CriticalAlerts++
			}
		}
	}

	// System health (simplified)
	if overview.CriticalAlerts > 0 {
		overview.SystemHealth = HealthStatusCritical
	} else if overview.ActiveAlerts > 5 {
		overview.SystemHealth = HealthStatusWarning
	} else {
		overview.SystemHealth = HealthStatusHealthy
	}
}

// updateBotSummaries updates bot summary information
func (dm *DashboardManager) updateBotSummaries(botMetrics map[string]*BotMetrics) {
	summaries := make([]*BotSummary, 0, len(botMetrics))

	for botID, metrics := range botMetrics {
		summary := &BotSummary{
			BotID:         botID,
			Name:          botID, // Simplified
			Strategy:      metrics.Strategy,
			State:         metrics.State,
			Health:        metrics.Health.OverallHealth,
			RiskScore:     0,
			ActiveAlerts:  len(metrics.Alerts),
			Uptime:        24 * time.Hour, // Simplified
		}

		// Performance metrics
		if metrics.Performance != nil {
			summary.PnL = metrics.Performance.TotalReturn
			summary.DailyPnL = metrics.Performance.TotalReturn.Div(decimal.NewFromInt(30)) // Simplified daily
			summary.WinRate = metrics.Performance.WinRate
			summary.TotalTrades = metrics.Performance.TotalTrades
		}

		// Risk metrics
		if metrics.Risk != nil {
			summary.RiskScore = metrics.Risk.RiskScore
		}

		// Trading metrics
		if metrics.Trading != nil {
			summary.LastTradeTime = metrics.Trading.LastTradeTime
		}

		summaries = append(summaries, summary)
	}

	dm.dashboardData.BotSummaries = summaries
}

// updateSystemStatus updates system status information
func (dm *DashboardManager) updateSystemStatus(systemMetrics *SystemMetrics) {
	status := dm.dashboardData.SystemStatus

	status.CPUUsage = systemMetrics.CPUUsage
	status.MemoryUsage = systemMetrics.MemoryUsage
	status.DiskUsage = systemMetrics.DiskUsage
	status.NetworkIO = systemMetrics.NetworkIO
	status.APIRequestRate = systemMetrics.APIRequestRate
	status.ErrorRate = systemMetrics.ErrorRate
	status.ResponseTime = systemMetrics.ResponseTime
	status.Uptime = 24 * time.Hour // Simplified
	status.LastHealthCheck = time.Now()

	// Determine overall health
	if status.CPUUsage > 90 || status.ErrorRate > 10 {
		status.OverallHealth = HealthStatusCritical
	} else if status.CPUUsage > 80 || status.ErrorRate > 5 {
		status.OverallHealth = HealthStatusWarning
	} else {
		status.OverallHealth = HealthStatusHealthy
	}
}

// updateRecentAlerts updates recent alerts list
func (dm *DashboardManager) updateRecentAlerts(alerts []*Alert) {
	// Keep only the 10 most recent alerts
	recentAlerts := make([]*Alert, 0)
	for i, alert := range alerts {
		if i >= 10 {
			break
		}
		recentAlerts = append(recentAlerts, alert)
	}

	dm.dashboardData.RecentAlerts = recentAlerts
}

// updatePerformanceChart updates performance chart data
func (dm *DashboardManager) updatePerformanceChart(portfolioMetrics *PortfolioMetrics) {
	chart := dm.dashboardData.PerformanceChart

	// Add new data point
	dataPoint := &PerformanceDataPoint{
		Timestamp:       time.Now(),
		PortfolioValue:  portfolioMetrics.TotalValue,
		TotalPnL:        portfolioMetrics.TotalPnL,
		DailyPnL:        portfolioMetrics.DailyPnL,
		DrawdownPercent: portfolioMetrics.MaxDrawdown,
	}

	if chart.DataPoints == nil {
		chart.DataPoints = make([]*PerformanceDataPoint, 0)
	}

	chart.DataPoints = append(chart.DataPoints, dataPoint)

	// Keep only last 24 hours of data (assuming 1-minute intervals)
	maxPoints := 24 * 60
	if len(chart.DataPoints) > maxPoints {
		chart.DataPoints = chart.DataPoints[len(chart.DataPoints)-maxPoints:]
	}

	chart.TimeRange = "24h"
}

// updateRiskMetrics updates risk metrics for dashboard
func (dm *DashboardManager) updateRiskMetrics(botMetrics map[string]*BotMetrics, portfolioMetrics *PortfolioMetrics) {
	risk := dm.dashboardData.RiskMetrics

	// Portfolio risk metrics
	risk.PortfolioVaR95 = portfolioMetrics.VaR95
	risk.PortfolioVaR99 = portfolioMetrics.VaR95.Mul(decimal.NewFromFloat(1.3)) // Approximate
	risk.MaxDrawdown = portfolioMetrics.MaxDrawdown
	risk.SharpeRatio = portfolioMetrics.SharpeRatio

	// Risk distribution
	riskDistribution := make(map[string]int)
	riskDistribution["low"] = 0
	riskDistribution["medium"] = 0
	riskDistribution["high"] = 0
	riskDistribution["critical"] = 0

	for _, metrics := range botMetrics {
		if metrics.Risk != nil {
			switch {
			case metrics.Risk.RiskScore < 25:
				riskDistribution["low"]++
			case metrics.Risk.RiskScore < 50:
				riskDistribution["medium"]++
			case metrics.Risk.RiskScore < 75:
				riskDistribution["high"]++
			default:
				riskDistribution["critical"]++
			}
		}
	}

	risk.RiskDistribution = riskDistribution
}

// updateTradingActivity updates trading activity summary
func (dm *DashboardManager) updateTradingActivity(botMetrics map[string]*BotMetrics) {
	activity := dm.dashboardData.TradingActivity

	// Aggregate trading metrics
	totalOrders := 0
	filledOrders := 0
	cancelledOrders := 0
	failedOrders := 0
	totalVolume := decimal.Zero
	totalFees := decimal.Zero

	for _, metrics := range botMetrics {
		if metrics.Trading != nil {
			totalOrders += metrics.Trading.OrdersPlaced
			filledOrders += metrics.Trading.OrdersFilled
			cancelledOrders += metrics.Trading.OrdersCancelled
			failedOrders += metrics.Trading.OrdersFailed
			totalVolume = totalVolume.Add(metrics.Trading.TotalVolume)
			totalFees = totalFees.Add(metrics.Trading.TotalFees)
		}
	}

	activity.TotalOrders24h = totalOrders
	activity.FilledOrders24h = filledOrders
	activity.CancelledOrders24h = cancelledOrders
	activity.FailedOrders24h = failedOrders
	activity.TotalVolume24h = totalVolume
	activity.TotalFees24h = totalFees

	// Calculate averages
	if totalOrders > 0 {
		activity.AvgFillRate = decimal.NewFromInt(int64(filledOrders)).Div(decimal.NewFromInt(int64(totalOrders)))
	}

	activity.AvgSlippage = decimal.NewFromFloat(0.001) // Simplified
}

// GetDashboardData returns the current dashboard data
func (dm *DashboardManager) GetDashboardData() *DashboardData {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	return dm.dashboardData
}

// GetDashboardJSON returns dashboard data as JSON
func (dm *DashboardManager) GetDashboardJSON() ([]byte, error) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	return json.Marshal(dm.dashboardData)
}
