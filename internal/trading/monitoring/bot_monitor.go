package monitoring

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/internal/trading"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/shopspring/decimal"
)

// TradingBotMonitor provides comprehensive monitoring for trading bots
type TradingBotMonitor struct {
	logger      *observability.Logger
	config      *MonitoringConfig
	botEngine   *trading.TradingBotEngine
	riskManager *trading.BotRiskManager

	// Metrics storage
	botMetrics       map[string]*BotMetrics
	portfolioMetrics *PortfolioMetrics
	systemMetrics    *SystemMetrics

	// Performance tracking
	performanceHistory map[string][]*PerformanceSnapshot
	alertHistory       []*Alert

	// State management
	isRunning bool
	stopChan  chan struct{}
	mu        sync.RWMutex

	// Collectors
	metricsCollector *MetricsCollector
	alertManager     *AlertManager
	dashboardManager *DashboardManager
}

// MonitoringConfig holds configuration for trading bot monitoring
type MonitoringConfig struct {
	// Collection intervals
	MetricsInterval     time.Duration `yaml:"metrics_interval"`
	HealthCheckInterval time.Duration `yaml:"health_check_interval"`
	AlertCheckInterval  time.Duration `yaml:"alert_check_interval"`

	// Data retention
	MetricsRetention  time.Duration `yaml:"metrics_retention"`
	AlertRetention    time.Duration `yaml:"alert_retention"`
	SnapshotRetention time.Duration `yaml:"snapshot_retention"`

	// Performance thresholds
	Thresholds *PerformanceThresholds `yaml:"thresholds"`

	// Features
	EnableRealTimeAlerts bool `yaml:"enable_realtime_alerts"`
	EnableDashboard      bool `yaml:"enable_dashboard"`
	EnableMetricsExport  bool `yaml:"enable_metrics_export"`
	EnableProfiling      bool `yaml:"enable_profiling"`
}

// PerformanceThresholds defines performance alert thresholds
type PerformanceThresholds struct {
	// Bot performance thresholds
	MinWinRate     decimal.Decimal `yaml:"min_win_rate"`
	MaxDrawdown    decimal.Decimal `yaml:"max_drawdown"`
	MaxDailyLoss   decimal.Decimal `yaml:"max_daily_loss"`
	MinSharpeRatio decimal.Decimal `yaml:"min_sharpe_ratio"`

	// System performance thresholds
	MaxCPUUsage     float64       `yaml:"max_cpu_usage"`
	MaxMemoryUsage  float64       `yaml:"max_memory_usage"`
	MaxResponseTime time.Duration `yaml:"max_response_time"`
	MaxErrorRate    float64       `yaml:"max_error_rate"`

	// Trading performance thresholds
	MaxExecutionTime time.Duration   `yaml:"max_execution_time"`
	MinOrderFillRate decimal.Decimal `yaml:"min_order_fill_rate"`
	MaxSlippage      decimal.Decimal `yaml:"max_slippage"`
}

// BotMetrics contains comprehensive metrics for a single trading bot
type BotMetrics struct {
	BotID     string    `json:"bot_id"`
	Strategy  string    `json:"strategy"`
	State     string    `json:"state"`
	Timestamp time.Time `json:"timestamp"`

	// Performance metrics
	Performance *BotPerformance   `json:"performance"`
	Risk        *BotRiskMetrics   `json:"risk"`
	Trading     *TradingMetrics   `json:"trading"`
	System      *BotSystemMetrics `json:"system"`

	// Health indicators
	Health *HealthIndicators `json:"health"`
	Alerts []*BotAlert       `json:"alerts"`
}

// BotPerformance tracks bot performance metrics
type BotPerformance struct {
	TotalTrades      int             `json:"total_trades"`
	WinningTrades    int             `json:"winning_trades"`
	LosingTrades     int             `json:"losing_trades"`
	WinRate          decimal.Decimal `json:"win_rate"`
	ProfitFactor     decimal.Decimal `json:"profit_factor"`
	SharpeRatio      decimal.Decimal `json:"sharpe_ratio"`
	SortinoRatio     decimal.Decimal `json:"sortino_ratio"`
	MaxDrawdown      decimal.Decimal `json:"max_drawdown"`
	CurrentDrawdown  decimal.Decimal `json:"current_drawdown"`
	TotalReturn      decimal.Decimal `json:"total_return"`
	AnnualizedReturn decimal.Decimal `json:"annualized_return"`
	Volatility       decimal.Decimal `json:"volatility"`
	Beta             decimal.Decimal `json:"beta"`
	Alpha            decimal.Decimal `json:"alpha"`
	CalmarRatio      decimal.Decimal `json:"calmar_ratio"`
}

// BotRiskMetrics tracks bot risk metrics
type BotRiskMetrics struct {
	VaR95             decimal.Decimal `json:"var_95"`
	VaR99             decimal.Decimal `json:"var_99"`
	ExpectedShortfall decimal.Decimal `json:"expected_shortfall"`
	RiskScore         int             `json:"risk_score"`
	ExposureRatio     decimal.Decimal `json:"exposure_ratio"`
	LeverageRatio     decimal.Decimal `json:"leverage_ratio"`
	ConcentrationRisk decimal.Decimal `json:"concentration_risk"`
	CorrelationRisk   decimal.Decimal `json:"correlation_risk"`
	LiquidityRisk     decimal.Decimal `json:"liquidity_risk"`
}

// TradingMetrics tracks trading execution metrics
type TradingMetrics struct {
	OrdersPlaced     int             `json:"orders_placed"`
	OrdersFilled     int             `json:"orders_filled"`
	OrdersCancelled  int             `json:"orders_cancelled"`
	OrdersFailed     int             `json:"orders_failed"`
	FillRate         decimal.Decimal `json:"fill_rate"`
	AvgExecutionTime time.Duration   `json:"avg_execution_time"`
	AvgSlippage      decimal.Decimal `json:"avg_slippage"`
	TotalVolume      decimal.Decimal `json:"total_volume"`
	TotalFees        decimal.Decimal `json:"total_fees"`
	LastTradeTime    time.Time       `json:"last_trade_time"`
}

// BotSystemMetrics tracks system-level metrics for a bot
type BotSystemMetrics struct {
	CPUUsage          float64       `json:"cpu_usage"`
	MemoryUsage       int64         `json:"memory_usage"`
	GoroutineCount    int           `json:"goroutine_count"`
	APICallsPerMinute int           `json:"api_calls_per_minute"`
	ErrorCount        int           `json:"error_count"`
	LastErrorTime     time.Time     `json:"last_error_time"`
	Uptime            time.Duration `json:"uptime"`
}

// HealthIndicators provides health status indicators
type HealthIndicators struct {
	OverallHealth     HealthStatus `json:"overall_health"`
	PerformanceHealth HealthStatus `json:"performance_health"`
	RiskHealth        HealthStatus `json:"risk_health"`
	TradingHealth     HealthStatus `json:"trading_health"`
	SystemHealth      HealthStatus `json:"system_health"`
	LastHealthCheck   time.Time    `json:"last_health_check"`
}

// HealthStatus represents health status levels
type HealthStatus string

const (
	HealthStatusHealthy  HealthStatus = "healthy"
	HealthStatusWarning  HealthStatus = "warning"
	HealthStatusCritical HealthStatus = "critical"
	HealthStatusUnknown  HealthStatus = "unknown"
)

// BotAlert represents an alert for a specific bot
type BotAlert struct {
	ID           string                 `json:"id"`
	BotID        string                 `json:"bot_id"`
	Type         AlertType              `json:"type"`
	Severity     AlertSeverity          `json:"severity"`
	Message      string                 `json:"message"`
	Timestamp    time.Time              `json:"timestamp"`
	Acknowledged bool                   `json:"acknowledged"`
	Resolved     bool                   `json:"resolved"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// AlertType defines types of alerts
type AlertType string

const (
	AlertTypePerformance AlertType = "performance"
	AlertTypeRisk        AlertType = "risk"
	AlertTypeTrading     AlertType = "trading"
	AlertTypeSystem      AlertType = "system"
	AlertTypeHealth      AlertType = "health"
)

// AlertSeverity defines alert severity levels
type AlertSeverity string

const (
	AlertSeverityInfo     AlertSeverity = "info"
	AlertSeverityWarning  AlertSeverity = "warning"
	AlertSeverityHigh     AlertSeverity = "high"
	AlertSeverityCritical AlertSeverity = "critical"
)

// PortfolioMetrics tracks portfolio-level metrics
type PortfolioMetrics struct {
	Timestamp         time.Time                  `json:"timestamp"`
	TotalValue        decimal.Decimal            `json:"total_value"`
	TotalPnL          decimal.Decimal            `json:"total_pnl"`
	DailyPnL          decimal.Decimal            `json:"daily_pnl"`
	TotalReturn       decimal.Decimal            `json:"total_return"`
	SharpeRatio       decimal.Decimal            `json:"sharpe_ratio"`
	MaxDrawdown       decimal.Decimal            `json:"max_drawdown"`
	VaR95             decimal.Decimal            `json:"var_95"`
	ActiveBots        int                        `json:"active_bots"`
	TotalBots         int                        `json:"total_bots"`
	StrategyBreakdown map[string]decimal.Decimal `json:"strategy_breakdown"`
	AssetAllocation   map[string]decimal.Decimal `json:"asset_allocation"`
}

// SystemMetrics tracks system-level metrics
type SystemMetrics struct {
	Timestamp      time.Time     `json:"timestamp"`
	CPUUsage       float64       `json:"cpu_usage"`
	MemoryUsage    int64         `json:"memory_usage"`
	DiskUsage      int64         `json:"disk_usage"`
	NetworkIO      int64         `json:"network_io"`
	GoroutineCount int           `json:"goroutine_count"`
	APIRequestRate float64       `json:"api_request_rate"`
	ErrorRate      float64       `json:"error_rate"`
	ResponseTime   time.Duration `json:"response_time"`
}

// PerformanceSnapshot captures a point-in-time performance snapshot
type PerformanceSnapshot struct {
	Timestamp   time.Time         `json:"timestamp"`
	BotID       string            `json:"bot_id"`
	Performance *BotPerformance   `json:"performance"`
	Risk        *BotRiskMetrics   `json:"risk"`
	Trading     *TradingMetrics   `json:"trading"`
	System      *BotSystemMetrics `json:"system"`
}

// Alert represents a monitoring alert
type Alert struct {
	ID           string                 `json:"id"`
	Type         AlertType              `json:"type"`
	Severity     AlertSeverity          `json:"severity"`
	Title        string                 `json:"title"`
	Message      string                 `json:"message"`
	BotID        string                 `json:"bot_id,omitempty"`
	Timestamp    time.Time              `json:"timestamp"`
	Acknowledged bool                   `json:"acknowledged"`
	Resolved     bool                   `json:"resolved"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// NewTradingBotMonitor creates a new trading bot monitor
func NewTradingBotMonitor(
	logger *observability.Logger,
	config *MonitoringConfig,
	botEngine *trading.TradingBotEngine,
	riskManager *trading.BotRiskManager,
) *TradingBotMonitor {
	if config == nil {
		config = getDefaultMonitoringConfig()
	}

	return &TradingBotMonitor{
		logger:             logger,
		config:             config,
		botEngine:          botEngine,
		riskManager:        riskManager,
		botMetrics:         make(map[string]*BotMetrics),
		portfolioMetrics:   &PortfolioMetrics{},
		systemMetrics:      &SystemMetrics{},
		performanceHistory: make(map[string][]*PerformanceSnapshot),
		alertHistory:       make([]*Alert, 0),
		stopChan:           make(chan struct{}),
		metricsCollector:   NewMetricsCollector(logger),
		alertManager:       NewAlertManager(logger),
		dashboardManager:   NewDashboardManager(logger),
	}
}

// Start starts the trading bot monitor
func (tbm *TradingBotMonitor) Start(ctx context.Context) error {
	tbm.mu.Lock()
	defer tbm.mu.Unlock()

	if tbm.isRunning {
		return fmt.Errorf("trading bot monitor is already running")
	}

	tbm.isRunning = true

	// Start monitoring loops
	go tbm.metricsCollectionLoop(ctx)
	go tbm.healthCheckLoop(ctx)
	go tbm.alertProcessingLoop(ctx)

	if tbm.config.EnableDashboard {
		go tbm.dashboardUpdateLoop(ctx)
	}

	tbm.logger.Info(ctx, "Trading bot monitor started", map[string]interface{}{
		"metrics_interval":      tbm.config.MetricsInterval.String(),
		"health_check_interval": tbm.config.HealthCheckInterval.String(),
		"enable_dashboard":      tbm.config.EnableDashboard,
		"enable_alerts":         tbm.config.EnableRealTimeAlerts,
	})

	return nil
}

// Stop stops the trading bot monitor
func (tbm *TradingBotMonitor) Stop(ctx context.Context) error {
	tbm.mu.Lock()
	defer tbm.mu.Unlock()

	if !tbm.isRunning {
		return nil
	}

	tbm.isRunning = false
	close(tbm.stopChan)

	tbm.logger.Info(ctx, "Trading bot monitor stopped", nil)
	return nil
}

// getDefaultMonitoringConfig returns default monitoring configuration
func getDefaultMonitoringConfig() *MonitoringConfig {
	return &MonitoringConfig{
		MetricsInterval:     30 * time.Second,
		HealthCheckInterval: 60 * time.Second,
		AlertCheckInterval:  10 * time.Second,
		MetricsRetention:    24 * time.Hour,
		AlertRetention:      7 * 24 * time.Hour,
		SnapshotRetention:   24 * time.Hour,
		Thresholds: &PerformanceThresholds{
			MinWinRate:       decimal.NewFromFloat(0.40),
			MaxDrawdown:      decimal.NewFromFloat(0.20),
			MaxDailyLoss:     decimal.NewFromFloat(0.05),
			MinSharpeRatio:   decimal.NewFromFloat(0.50),
			MaxCPUUsage:      80.0,
			MaxMemoryUsage:   85.0,
			MaxResponseTime:  2 * time.Second,
			MaxErrorRate:     5.0,
			MaxExecutionTime: 30 * time.Second,
			MinOrderFillRate: decimal.NewFromFloat(0.95),
			MaxSlippage:      decimal.NewFromFloat(0.01),
		},
		EnableRealTimeAlerts: true,
		EnableDashboard:      true,
		EnableMetricsExport:  true,
		EnableProfiling:      false,
	}
}

// metricsCollectionLoop continuously collects metrics from all bots
func (tbm *TradingBotMonitor) metricsCollectionLoop(ctx context.Context) {
	ticker := time.NewTicker(tbm.config.MetricsInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tbm.stopChan:
			return
		case <-ticker.C:
			tbm.collectAllMetrics(ctx)
		}
	}
}

// healthCheckLoop performs health checks on all bots
func (tbm *TradingBotMonitor) healthCheckLoop(ctx context.Context) {
	ticker := time.NewTicker(tbm.config.HealthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tbm.stopChan:
			return
		case <-ticker.C:
			tbm.performHealthChecks(ctx)
		}
	}
}

// alertProcessingLoop processes and manages alerts
func (tbm *TradingBotMonitor) alertProcessingLoop(ctx context.Context) {
	ticker := time.NewTicker(tbm.config.AlertCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tbm.stopChan:
			return
		case <-ticker.C:
			tbm.processAlerts(ctx)
		}
	}
}

// dashboardUpdateLoop updates dashboard data
func (tbm *TradingBotMonitor) dashboardUpdateLoop(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tbm.stopChan:
			return
		case <-ticker.C:
			tbm.updateDashboard(ctx)
		}
	}
}

// collectAllMetrics collects metrics from all bots and portfolio
func (tbm *TradingBotMonitor) collectAllMetrics(ctx context.Context) {
	tbm.mu.Lock()
	defer tbm.mu.Unlock()

	// Get all bots from the engine
	bots := tbm.botEngine.ListBots()

	// Collect metrics for each bot
	for _, bot := range bots {
		metrics := tbm.collectBotMetrics(ctx, bot)
		tbm.botMetrics[bot.ID] = metrics

		// Create performance snapshot
		snapshot := &PerformanceSnapshot{
			Timestamp:   time.Now(),
			BotID:       bot.ID,
			Performance: metrics.Performance,
			Risk:        metrics.Risk,
			Trading:     metrics.Trading,
			System:      metrics.System,
		}

		// Add to history
		tbm.addPerformanceSnapshot(bot.ID, snapshot)
	}

	// Collect portfolio metrics
	tbm.portfolioMetrics = tbm.collectPortfolioMetrics(ctx, bots)

	// Collect system metrics
	tbm.systemMetrics = tbm.collectSystemMetrics(ctx)

	tbm.logger.Debug(ctx, "Metrics collected for all bots", map[string]interface{}{
		"bot_count":       len(bots),
		"portfolio_value": tbm.portfolioMetrics.TotalValue.String(),
		"system_cpu":      tbm.systemMetrics.CPUUsage,
		"system_memory":   tbm.systemMetrics.MemoryUsage,
	})
}

// collectBotMetrics collects comprehensive metrics for a single bot
func (tbm *TradingBotMonitor) collectBotMetrics(ctx context.Context, bot *trading.TradingBot) *BotMetrics {
	metrics := &BotMetrics{
		BotID:     bot.ID,
		Strategy:  string(bot.Strategy),
		State:     string(bot.State),
		Timestamp: time.Now(),
	}

	// Collect performance metrics
	metrics.Performance = tbm.collectBotPerformance(ctx, bot)

	// Collect risk metrics
	if riskMetrics, err := tbm.riskManager.GetBotRiskMetrics(bot.ID); err == nil {
		metrics.Risk = tbm.convertRiskMetrics(riskMetrics)
	}

	// Collect trading metrics
	metrics.Trading = tbm.collectTradingMetrics(ctx, bot)

	// Collect system metrics
	metrics.System = tbm.collectBotSystemMetrics(ctx, bot)

	// Calculate health indicators
	metrics.Health = tbm.calculateHealthIndicators(metrics)

	// Get active alerts for this bot
	metrics.Alerts = tbm.getBotAlerts(bot.ID)

	return metrics
}

// collectBotPerformance collects performance metrics for a bot
func (tbm *TradingBotMonitor) collectBotPerformance(ctx context.Context, bot *trading.TradingBot) *BotPerformance {
	// Get performance data from bot
	botPerf := bot.Performance
	if botPerf == nil {
		return &BotPerformance{}
	}

	// Calculate additional metrics
	winRate := decimal.Zero
	if botPerf.TotalTrades > 0 {
		winRate = decimal.NewFromInt(int64(botPerf.WinningTrades)).Div(decimal.NewFromInt(int64(botPerf.TotalTrades)))
	}

	profitFactor := decimal.Zero
	if !botPerf.TotalLoss.IsZero() {
		profitFactor = botPerf.TotalProfit.Div(botPerf.TotalLoss.Abs())
	}

	return &BotPerformance{
		TotalTrades:      botPerf.TotalTrades,
		WinningTrades:    botPerf.WinningTrades,
		LosingTrades:     botPerf.LosingTrades,
		WinRate:          winRate,
		ProfitFactor:     profitFactor,
		SharpeRatio:      botPerf.SharpeRatio,
		SortinoRatio:     tbm.calculateSortinoRatio(bot),
		MaxDrawdown:      botPerf.MaxDrawdown,
		CurrentDrawdown:  tbm.calculateCurrentDrawdown(bot),
		TotalReturn:      botPerf.NetProfit,
		AnnualizedReturn: tbm.calculateAnnualizedReturn(bot),
		Volatility:       tbm.calculateVolatility(bot),
		Beta:             tbm.calculateBeta(bot),
		Alpha:            tbm.calculateAlpha(bot),
		CalmarRatio:      tbm.calculateCalmarRatio(bot),
	}
}

// addPerformanceSnapshot adds a performance snapshot to history
func (tbm *TradingBotMonitor) addPerformanceSnapshot(botID string, snapshot *PerformanceSnapshot) {
	if tbm.performanceHistory[botID] == nil {
		tbm.performanceHistory[botID] = make([]*PerformanceSnapshot, 0)
	}

	tbm.performanceHistory[botID] = append(tbm.performanceHistory[botID], snapshot)

	// Maintain retention limit
	maxSnapshots := int(tbm.config.SnapshotRetention / tbm.config.MetricsInterval)
	if len(tbm.performanceHistory[botID]) > maxSnapshots {
		tbm.performanceHistory[botID] = tbm.performanceHistory[botID][1:]
	}
}

// performHealthChecks performs health checks on all bots
func (tbm *TradingBotMonitor) performHealthChecks(ctx context.Context) {
	tbm.mu.RLock()
	defer tbm.mu.RUnlock()

	for botID, metrics := range tbm.botMetrics {
		health := tbm.calculateHealthIndicators(metrics)
		metrics.Health = health

		tbm.logger.Debug(ctx, "Health check completed", map[string]interface{}{
			"bot_id":         botID,
			"overall_health": string(health.OverallHealth),
			"risk_health":    string(health.RiskHealth),
			"trading_health": string(health.TradingHealth),
			"system_health":  string(health.SystemHealth),
		})
	}
}

// processAlerts processes and manages alerts
func (tbm *TradingBotMonitor) processAlerts(ctx context.Context) {
	if !tbm.config.EnableRealTimeAlerts {
		return
	}

	tbm.alertManager.ProcessAlerts(ctx, tbm.botMetrics, tbm.config.Thresholds)
}

// updateDashboard updates dashboard data
func (tbm *TradingBotMonitor) updateDashboard(ctx context.Context) {
	if !tbm.config.EnableDashboard {
		return
	}

	tbm.mu.RLock()
	botMetrics := tbm.botMetrics
	portfolioMetrics := tbm.portfolioMetrics
	systemMetrics := tbm.systemMetrics
	tbm.mu.RUnlock()

	alerts := tbm.alertManager.GetAllAlerts()

	tbm.dashboardManager.UpdateDashboard(ctx, botMetrics, portfolioMetrics, systemMetrics, alerts)
}

// collectPortfolioMetrics collects portfolio-level metrics
func (tbm *TradingBotMonitor) collectPortfolioMetrics(ctx context.Context, bots []*trading.TradingBot) *PortfolioMetrics {
	return tbm.metricsCollector.CollectPortfolioMetrics(ctx, bots)
}

// collectSystemMetrics collects system-level metrics
func (tbm *TradingBotMonitor) collectSystemMetrics(ctx context.Context) *SystemMetrics {
	return tbm.metricsCollector.CollectSystemMetrics(ctx)
}

// collectTradingMetrics collects trading metrics for a bot
func (tbm *TradingBotMonitor) collectTradingMetrics(ctx context.Context, bot *trading.TradingBot) *TradingMetrics {
	return tbm.metricsCollector.CollectTradingMetrics(ctx, bot)
}

// collectBotSystemMetrics collects system metrics for a bot
func (tbm *TradingBotMonitor) collectBotSystemMetrics(ctx context.Context, bot *trading.TradingBot) *BotSystemMetrics {
	return tbm.metricsCollector.CollectBotSystemMetrics(ctx, bot)
}

// convertRiskMetrics converts risk manager metrics to monitoring metrics
func (tbm *TradingBotMonitor) convertRiskMetrics(riskMetrics *trading.BotRiskMetrics) *BotRiskMetrics {
	return tbm.metricsCollector.ConvertRiskMetrics(riskMetrics)
}

// calculateHealthIndicators calculates health indicators for a bot
func (tbm *TradingBotMonitor) calculateHealthIndicators(metrics *BotMetrics) *HealthIndicators {
	health := &HealthIndicators{
		LastHealthCheck: time.Now(),
	}

	// Performance health
	if metrics.Performance != nil {
		if metrics.Performance.WinRate.GreaterThan(decimal.NewFromFloat(0.6)) {
			health.PerformanceHealth = HealthStatusHealthy
		} else if metrics.Performance.WinRate.GreaterThan(decimal.NewFromFloat(0.4)) {
			health.PerformanceHealth = HealthStatusWarning
		} else {
			health.PerformanceHealth = HealthStatusCritical
		}
	} else {
		health.PerformanceHealth = HealthStatusUnknown
	}

	// Risk health
	if metrics.Risk != nil {
		if metrics.Risk.RiskScore < 50 {
			health.RiskHealth = HealthStatusHealthy
		} else if metrics.Risk.RiskScore < 75 {
			health.RiskHealth = HealthStatusWarning
		} else {
			health.RiskHealth = HealthStatusCritical
		}
	} else {
		health.RiskHealth = HealthStatusUnknown
	}

	// Trading health
	if metrics.Trading != nil {
		if metrics.Trading.FillRate.GreaterThan(decimal.NewFromFloat(0.9)) {
			health.TradingHealth = HealthStatusHealthy
		} else if metrics.Trading.FillRate.GreaterThan(decimal.NewFromFloat(0.8)) {
			health.TradingHealth = HealthStatusWarning
		} else {
			health.TradingHealth = HealthStatusCritical
		}
	} else {
		health.TradingHealth = HealthStatusUnknown
	}

	// System health
	if metrics.System != nil {
		if metrics.System.CPUUsage < 70 && metrics.System.ErrorCount < 5 {
			health.SystemHealth = HealthStatusHealthy
		} else if metrics.System.CPUUsage < 85 && metrics.System.ErrorCount < 10 {
			health.SystemHealth = HealthStatusWarning
		} else {
			health.SystemHealth = HealthStatusCritical
		}
	} else {
		health.SystemHealth = HealthStatusUnknown
	}

	// Overall health (worst of all components)
	healths := []HealthStatus{
		health.PerformanceHealth,
		health.RiskHealth,
		health.TradingHealth,
		health.SystemHealth,
	}

	health.OverallHealth = HealthStatusHealthy
	for _, h := range healths {
		if h == HealthStatusCritical {
			health.OverallHealth = HealthStatusCritical
			break
		} else if h == HealthStatusWarning && health.OverallHealth != HealthStatusCritical {
			health.OverallHealth = HealthStatusWarning
		}
	}

	return health
}

// getBotAlerts returns active alerts for a bot
func (tbm *TradingBotMonitor) getBotAlerts(botID string) []*BotAlert {
	return tbm.alertManager.GetBotAlerts(botID)
}

// Performance calculation methods

// calculateSortinoRatio calculates the Sortino ratio for a bot
func (tbm *TradingBotMonitor) calculateSortinoRatio(bot *trading.TradingBot) decimal.Decimal {
	// Simplified Sortino ratio calculation
	if bot.Performance == nil {
		return decimal.Zero
	}

	// Use Sharpe ratio as approximation for now
	return bot.Performance.SharpeRatio.Mul(decimal.NewFromFloat(1.1))
}

// calculateCurrentDrawdown calculates current drawdown for a bot
func (tbm *TradingBotMonitor) calculateCurrentDrawdown(bot *trading.TradingBot) decimal.Decimal {
	// Simplified current drawdown calculation
	if bot.Performance == nil {
		return decimal.Zero
	}

	// Use max drawdown as approximation for current
	return bot.Performance.MaxDrawdown.Mul(decimal.NewFromFloat(0.7))
}

// calculateAnnualizedReturn calculates annualized return for a bot
func (tbm *TradingBotMonitor) calculateAnnualizedReturn(bot *trading.TradingBot) decimal.Decimal {
	if bot.Performance == nil {
		return decimal.Zero
	}

	// Simplified annualized return (assume 30 days of operation)
	dailyReturn := bot.Performance.NetProfit.Div(decimal.NewFromInt(30))
	return dailyReturn.Mul(decimal.NewFromInt(365))
}

// calculateVolatility calculates volatility for a bot
func (tbm *TradingBotMonitor) calculateVolatility(bot *trading.TradingBot) decimal.Decimal {
	// Simplified volatility calculation
	if bot.Performance == nil {
		return decimal.Zero
	}

	// Estimate volatility based on drawdown
	return bot.Performance.MaxDrawdown.Mul(decimal.NewFromFloat(2.0))
}

// calculateBeta calculates beta for a bot
func (tbm *TradingBotMonitor) calculateBeta(bot *trading.TradingBot) decimal.Decimal {
	// Simplified beta calculation (correlation with market)
	// For crypto bots, assume moderate correlation with BTC
	return decimal.NewFromFloat(0.7)
}

// calculateAlpha calculates alpha for a bot
func (tbm *TradingBotMonitor) calculateAlpha(bot *trading.TradingBot) decimal.Decimal {
	if bot.Performance == nil {
		return decimal.Zero
	}

	// Simplified alpha calculation
	// Alpha = Return - (Beta * Market Return)
	beta := tbm.calculateBeta(bot)
	marketReturn := decimal.NewFromFloat(0.10) // Assume 10% market return
	expectedReturn := beta.Mul(marketReturn)

	actualReturn := bot.Performance.NetProfit.Div(decimal.NewFromFloat(10000)) // Assume 10k base
	return actualReturn.Sub(expectedReturn)
}

// calculateCalmarRatio calculates Calmar ratio for a bot
func (tbm *TradingBotMonitor) calculateCalmarRatio(bot *trading.TradingBot) decimal.Decimal {
	if bot.Performance == nil || bot.Performance.MaxDrawdown.IsZero() {
		return decimal.Zero
	}

	annualizedReturn := tbm.calculateAnnualizedReturn(bot)
	return annualizedReturn.Div(bot.Performance.MaxDrawdown)
}

// Public API methods

// GetBotMetrics returns metrics for a specific bot
func (tbm *TradingBotMonitor) GetBotMetrics(botID string) (*BotMetrics, error) {
	tbm.mu.RLock()
	defer tbm.mu.RUnlock()

	metrics, exists := tbm.botMetrics[botID]
	if !exists {
		return nil, fmt.Errorf("metrics not found for bot %s", botID)
	}

	return metrics, nil
}

// GetAllBotMetrics returns metrics for all bots
func (tbm *TradingBotMonitor) GetAllBotMetrics() map[string]*BotMetrics {
	tbm.mu.RLock()
	defer tbm.mu.RUnlock()

	// Return a copy to avoid concurrent access issues
	result := make(map[string]*BotMetrics)
	for k, v := range tbm.botMetrics {
		result[k] = v
	}

	return result
}

// GetPortfolioMetrics returns portfolio metrics
func (tbm *TradingBotMonitor) GetPortfolioMetrics() *PortfolioMetrics {
	tbm.mu.RLock()
	defer tbm.mu.RUnlock()

	return tbm.portfolioMetrics
}

// GetSystemMetrics returns system metrics
func (tbm *TradingBotMonitor) GetSystemMetrics() *SystemMetrics {
	tbm.mu.RLock()
	defer tbm.mu.RUnlock()

	return tbm.systemMetrics
}

// GetPerformanceHistory returns performance history for a bot
func (tbm *TradingBotMonitor) GetPerformanceHistory(botID string) ([]*PerformanceSnapshot, error) {
	tbm.mu.RLock()
	defer tbm.mu.RUnlock()

	history, exists := tbm.performanceHistory[botID]
	if !exists {
		return nil, fmt.Errorf("performance history not found for bot %s", botID)
	}

	return history, nil
}

// GetDashboardData returns dashboard data
func (tbm *TradingBotMonitor) GetDashboardData() *DashboardData {
	return tbm.dashboardManager.GetDashboardData()
}

// GetActiveAlerts returns all active alerts
func (tbm *TradingBotMonitor) GetActiveAlerts() []*Alert {
	return tbm.alertManager.GetAllAlerts()
}

// AcknowledgeAlert acknowledges an alert
func (tbm *TradingBotMonitor) AcknowledgeAlert(alertID string) error {
	return tbm.alertManager.AcknowledgeAlert(alertID)
}

// ResolveAlert resolves an alert
func (tbm *TradingBotMonitor) ResolveAlert(alertID string) error {
	return tbm.alertManager.ResolveAlert(alertID)
}
