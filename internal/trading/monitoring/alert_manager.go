package monitoring

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// AlertManager manages alerts for trading bot monitoring
type AlertManager struct {
	logger       *observability.Logger
	alerts       map[string]*Alert
	botAlerts    map[string][]*BotAlert
	alertHistory []*Alert
	mu           sync.RWMutex
}

// NewAlertManager creates a new alert manager
func NewAlertManager(logger *observability.Logger) *AlertManager {
	return &AlertManager{
		logger:       logger,
		alerts:       make(map[string]*Alert),
		botAlerts:    make(map[string][]*BotAlert),
		alertHistory: make([]*Alert, 0),
	}
}

// ProcessAlerts processes and evaluates alerts for all bots
func (am *AlertManager) ProcessAlerts(ctx context.Context, botMetrics map[string]*BotMetrics, thresholds *PerformanceThresholds) {
	am.mu.Lock()
	defer am.mu.Unlock()

	for botID, metrics := range botMetrics {
		am.evaluateBotAlerts(ctx, botID, metrics, thresholds)
	}

	// Clean up old alerts
	am.cleanupOldAlerts()
}

// evaluateBotAlerts evaluates alerts for a specific bot
func (am *AlertManager) evaluateBotAlerts(ctx context.Context, botID string, metrics *BotMetrics, thresholds *PerformanceThresholds) {
	// Performance alerts
	am.checkPerformanceAlerts(ctx, botID, metrics.Performance, thresholds)

	// Risk alerts
	if metrics.Risk != nil {
		am.checkRiskAlerts(ctx, botID, metrics.Risk, thresholds)
	}

	// Trading alerts
	if metrics.Trading != nil {
		am.checkTradingAlerts(ctx, botID, metrics.Trading, thresholds)
	}

	// System alerts
	if metrics.System != nil {
		am.checkSystemAlerts(ctx, botID, metrics.System, thresholds)
	}

	// Health alerts
	if metrics.Health != nil {
		am.checkHealthAlerts(ctx, botID, metrics.Health)
	}
}

// checkPerformanceAlerts checks performance-related alerts
func (am *AlertManager) checkPerformanceAlerts(ctx context.Context, botID string, performance *BotPerformance, thresholds *PerformanceThresholds) {
	if performance == nil {
		return
	}

	// Win rate alert
	if performance.WinRate.LessThan(thresholds.MinWinRate) {
		am.createAlert(ctx, &Alert{
			Type:     AlertTypePerformance,
			Severity: AlertSeverityWarning,
			Title:    "Low Win Rate",
			Message:  fmt.Sprintf("Bot %s win rate (%s) below threshold (%s)", botID, performance.WinRate.String(), thresholds.MinWinRate.String()),
			BotID:    botID,
			Metadata: map[string]interface{}{
				"current_win_rate": performance.WinRate.String(),
				"threshold":        thresholds.MinWinRate.String(),
				"total_trades":     performance.TotalTrades,
			},
		})
	}

	// Drawdown alert
	if performance.CurrentDrawdown.GreaterThan(thresholds.MaxDrawdown) {
		am.createAlert(ctx, &Alert{
			Type:     AlertTypePerformance,
			Severity: AlertSeverityHigh,
			Title:    "High Drawdown",
			Message:  fmt.Sprintf("Bot %s drawdown (%s) exceeds threshold (%s)", botID, performance.CurrentDrawdown.String(), thresholds.MaxDrawdown.String()),
			BotID:    botID,
			Metadata: map[string]interface{}{
				"current_drawdown": performance.CurrentDrawdown.String(),
				"max_drawdown":     performance.MaxDrawdown.String(),
				"threshold":        thresholds.MaxDrawdown.String(),
			},
		})
	}

	// Sharpe ratio alert
	if performance.SharpeRatio.LessThan(thresholds.MinSharpeRatio) {
		am.createAlert(ctx, &Alert{
			Type:     AlertTypePerformance,
			Severity: AlertSeverityWarning,
			Title:    "Low Sharpe Ratio",
			Message:  fmt.Sprintf("Bot %s Sharpe ratio (%s) below threshold (%s)", botID, performance.SharpeRatio.String(), thresholds.MinSharpeRatio.String()),
			BotID:    botID,
			Metadata: map[string]interface{}{
				"current_sharpe": performance.SharpeRatio.String(),
				"threshold":      thresholds.MinSharpeRatio.String(),
			},
		})
	}
}

// checkRiskAlerts checks risk-related alerts
func (am *AlertManager) checkRiskAlerts(ctx context.Context, botID string, risk *BotRiskMetrics, thresholds *PerformanceThresholds) {
	// High risk score alert
	if risk.RiskScore > 80 {
		am.createAlert(ctx, &Alert{
			Type:     AlertTypeRisk,
			Severity: AlertSeverityHigh,
			Title:    "High Risk Score",
			Message:  fmt.Sprintf("Bot %s risk score (%d) is high", botID, risk.RiskScore),
			BotID:    botID,
			Metadata: map[string]interface{}{
				"risk_score": risk.RiskScore,
				"var_95":     risk.VaR95.String(),
				"exposure":   risk.ExposureRatio.String(),
			},
		})
	}

	// High VaR alert
	if risk.VaR95.GreaterThan(decimal.NewFromFloat(1000)) { // $1000 VaR threshold
		am.createAlert(ctx, &Alert{
			Type:     AlertTypeRisk,
			Severity: AlertSeverityWarning,
			Title:    "High Value at Risk",
			Message:  fmt.Sprintf("Bot %s VaR95 (%s) is elevated", botID, risk.VaR95.String()),
			BotID:    botID,
			Metadata: map[string]interface{}{
				"var_95":     risk.VaR95.String(),
				"var_99":     risk.VaR99.String(),
				"risk_score": risk.RiskScore,
			},
		})
	}
}

// checkTradingAlerts checks trading execution alerts
func (am *AlertManager) checkTradingAlerts(ctx context.Context, botID string, trading *TradingMetrics, thresholds *PerformanceThresholds) {
	// Low fill rate alert
	if trading.FillRate.LessThan(thresholds.MinOrderFillRate) {
		am.createAlert(ctx, &Alert{
			Type:     AlertTypeTrading,
			Severity: AlertSeverityWarning,
			Title:    "Low Order Fill Rate",
			Message:  fmt.Sprintf("Bot %s fill rate (%s) below threshold (%s)", botID, trading.FillRate.String(), thresholds.MinOrderFillRate.String()),
			BotID:    botID,
			Metadata: map[string]interface{}{
				"fill_rate":      trading.FillRate.String(),
				"threshold":      thresholds.MinOrderFillRate.String(),
				"orders_placed":  trading.OrdersPlaced,
				"orders_filled":  trading.OrdersFilled,
				"orders_failed":  trading.OrdersFailed,
			},
		})
	}

	// High slippage alert
	if trading.AvgSlippage.GreaterThan(thresholds.MaxSlippage) {
		am.createAlert(ctx, &Alert{
			Type:     AlertTypeTrading,
			Severity: AlertSeverityWarning,
			Title:    "High Slippage",
			Message:  fmt.Sprintf("Bot %s average slippage (%s) exceeds threshold (%s)", botID, trading.AvgSlippage.String(), thresholds.MaxSlippage.String()),
			BotID:    botID,
			Metadata: map[string]interface{}{
				"avg_slippage": trading.AvgSlippage.String(),
				"threshold":    thresholds.MaxSlippage.String(),
				"total_volume": trading.TotalVolume.String(),
			},
		})
	}

	// Slow execution alert
	if trading.AvgExecutionTime > thresholds.MaxExecutionTime {
		am.createAlert(ctx, &Alert{
			Type:     AlertTypeTrading,
			Severity: AlertSeverityWarning,
			Title:    "Slow Order Execution",
			Message:  fmt.Sprintf("Bot %s average execution time (%s) exceeds threshold (%s)", botID, trading.AvgExecutionTime.String(), thresholds.MaxExecutionTime.String()),
			BotID:    botID,
			Metadata: map[string]interface{}{
				"avg_execution_time": trading.AvgExecutionTime.String(),
				"threshold":          thresholds.MaxExecutionTime.String(),
			},
		})
	}
}

// checkSystemAlerts checks system-related alerts
func (am *AlertManager) checkSystemAlerts(ctx context.Context, botID string, system *BotSystemMetrics, thresholds *PerformanceThresholds) {
	// High CPU usage alert
	if system.CPUUsage > thresholds.MaxCPUUsage {
		am.createAlert(ctx, &Alert{
			Type:     AlertTypeSystem,
			Severity: AlertSeverityWarning,
			Title:    "High CPU Usage",
			Message:  fmt.Sprintf("Bot %s CPU usage (%.2f%%) exceeds threshold (%.2f%%)", botID, system.CPUUsage, thresholds.MaxCPUUsage),
			BotID:    botID,
			Metadata: map[string]interface{}{
				"cpu_usage":  system.CPUUsage,
				"threshold":  thresholds.MaxCPUUsage,
				"goroutines": system.GoroutineCount,
			},
		})
	}

	// High memory usage alert
	memoryUsagePercent := float64(system.MemoryUsage) / (1024 * 1024 * 100) * 100 // Convert to percentage
	if memoryUsagePercent > thresholds.MaxMemoryUsage {
		am.createAlert(ctx, &Alert{
			Type:     AlertTypeSystem,
			Severity: AlertSeverityWarning,
			Title:    "High Memory Usage",
			Message:  fmt.Sprintf("Bot %s memory usage (%.2f%%) exceeds threshold (%.2f%%)", botID, memoryUsagePercent, thresholds.MaxMemoryUsage),
			BotID:    botID,
			Metadata: map[string]interface{}{
				"memory_usage_mb": system.MemoryUsage / (1024 * 1024),
				"memory_usage_pct": memoryUsagePercent,
				"threshold":        thresholds.MaxMemoryUsage,
			},
		})
	}

	// High error count alert
	if system.ErrorCount > 10 {
		am.createAlert(ctx, &Alert{
			Type:     AlertTypeSystem,
			Severity: AlertSeverityHigh,
			Title:    "High Error Count",
			Message:  fmt.Sprintf("Bot %s has %d errors", botID, system.ErrorCount),
			BotID:    botID,
			Metadata: map[string]interface{}{
				"error_count":     system.ErrorCount,
				"last_error_time": system.LastErrorTime,
			},
		})
	}
}

// checkHealthAlerts checks health-related alerts
func (am *AlertManager) checkHealthAlerts(ctx context.Context, botID string, health *HealthIndicators) {
	if health.OverallHealth == HealthStatusCritical {
		am.createAlert(ctx, &Alert{
			Type:     AlertTypeHealth,
			Severity: AlertSeverityCritical,
			Title:    "Critical Health Status",
			Message:  fmt.Sprintf("Bot %s overall health is critical", botID),
			BotID:    botID,
			Metadata: map[string]interface{}{
				"overall_health":     string(health.OverallHealth),
				"performance_health": string(health.PerformanceHealth),
				"risk_health":        string(health.RiskHealth),
				"trading_health":     string(health.TradingHealth),
				"system_health":      string(health.SystemHealth),
			},
		})
	}
}

// createAlert creates and stores a new alert
func (am *AlertManager) createAlert(ctx context.Context, alert *Alert) {
	alert.ID = uuid.New().String()
	alert.Timestamp = time.Now()
	alert.Acknowledged = false
	alert.Resolved = false

	// Store in main alerts map
	am.alerts[alert.ID] = alert

	// Store in bot-specific alerts
	if alert.BotID != "" {
		if am.botAlerts[alert.BotID] == nil {
			am.botAlerts[alert.BotID] = make([]*BotAlert, 0)
		}

		botAlert := &BotAlert{
			ID:           alert.ID,
			BotID:        alert.BotID,
			Type:         alert.Type,
			Severity:     alert.Severity,
			Message:      alert.Message,
			Timestamp:    alert.Timestamp,
			Acknowledged: alert.Acknowledged,
			Resolved:     alert.Resolved,
			Metadata:     alert.Metadata,
		}

		am.botAlerts[alert.BotID] = append(am.botAlerts[alert.BotID], botAlert)
	}

	// Add to history
	am.alertHistory = append(am.alertHistory, alert)

	am.logger.Warn(ctx, "Alert created", map[string]interface{}{
		"alert_id":   alert.ID,
		"alert_type": string(alert.Type),
		"severity":   string(alert.Severity),
		"bot_id":     alert.BotID,
		"title":      alert.Title,
		"message":    alert.Message,
	})
}

// GetBotAlerts returns alerts for a specific bot
func (am *AlertManager) GetBotAlerts(botID string) []*BotAlert {
	am.mu.RLock()
	defer am.mu.RUnlock()

	alerts, exists := am.botAlerts[botID]
	if !exists {
		return []*BotAlert{}
	}

	// Return only active alerts
	activeAlerts := make([]*BotAlert, 0)
	for _, alert := range alerts {
		if !alert.Resolved {
			activeAlerts = append(activeAlerts, alert)
		}
	}

	return activeAlerts
}

// GetAllAlerts returns all active alerts
func (am *AlertManager) GetAllAlerts() []*Alert {
	am.mu.RLock()
	defer am.mu.RUnlock()

	activeAlerts := make([]*Alert, 0)
	for _, alert := range am.alerts {
		if !alert.Resolved {
			activeAlerts = append(activeAlerts, alert)
		}
	}

	return activeAlerts
}

// AcknowledgeAlert acknowledges an alert
func (am *AlertManager) AcknowledgeAlert(alertID string) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	alert, exists := am.alerts[alertID]
	if !exists {
		return fmt.Errorf("alert not found: %s", alertID)
	}

	alert.Acknowledged = true
	return nil
}

// ResolveAlert resolves an alert
func (am *AlertManager) ResolveAlert(alertID string) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	alert, exists := am.alerts[alertID]
	if !exists {
		return fmt.Errorf("alert not found: %s", alertID)
	}

	alert.Resolved = true
	alert.Acknowledged = true

	// Also resolve in bot alerts
	if alert.BotID != "" {
		if botAlerts, exists := am.botAlerts[alert.BotID]; exists {
			for _, botAlert := range botAlerts {
				if botAlert.ID == alertID {
					botAlert.Resolved = true
					botAlert.Acknowledged = true
					break
				}
			}
		}
	}

	return nil
}

// cleanupOldAlerts removes old resolved alerts
func (am *AlertManager) cleanupOldAlerts() {
	cutoff := time.Now().Add(-24 * time.Hour) // Keep alerts for 24 hours

	// Clean main alerts
	for id, alert := range am.alerts {
		if alert.Resolved && alert.Timestamp.Before(cutoff) {
			delete(am.alerts, id)
		}
	}

	// Clean bot alerts
	for botID, alerts := range am.botAlerts {
		filteredAlerts := make([]*BotAlert, 0)
		for _, alert := range alerts {
			if !alert.Resolved || alert.Timestamp.After(cutoff) {
				filteredAlerts = append(filteredAlerts, alert)
			}
		}
		am.botAlerts[botID] = filteredAlerts
	}

	// Clean alert history
	filteredHistory := make([]*Alert, 0)
	for _, alert := range am.alertHistory {
		if alert.Timestamp.After(cutoff) {
			filteredHistory = append(filteredHistory, alert)
		}
	}
	am.alertHistory = filteredHistory
}
