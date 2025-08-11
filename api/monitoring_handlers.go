package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ai-agentic-browser/internal/trading/monitoring"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/gorilla/mux"
)

// MonitoringHandler handles monitoring and analytics API requests
type MonitoringHandler struct {
	logger  *observability.Logger
	monitor *monitoring.TradingBotMonitor
}

// NewMonitoringHandler creates a new monitoring handler
func NewMonitoringHandler(logger *observability.Logger, monitor *monitoring.TradingBotMonitor) *MonitoringHandler {
	return &MonitoringHandler{
		logger:  logger,
		monitor: monitor,
	}
}

// RegisterRoutes registers monitoring API routes
func (h *MonitoringHandler) RegisterRoutes(router *mux.Router) {
	// Dashboard endpoints
	router.HandleFunc("/api/v1/monitoring/dashboard", h.GetDashboard).Methods("GET")
	router.HandleFunc("/api/v1/monitoring/overview", h.GetOverview).Methods("GET")

	// Bot metrics endpoints
	router.HandleFunc("/api/v1/monitoring/bots", h.GetAllBotMetrics).Methods("GET")
	router.HandleFunc("/api/v1/monitoring/bots/{botId}", h.GetBotMetrics).Methods("GET")
	router.HandleFunc("/api/v1/monitoring/bots/{botId}/performance", h.GetBotPerformance).Methods("GET")
	router.HandleFunc("/api/v1/monitoring/bots/{botId}/health", h.GetBotHealth).Methods("GET")
	router.HandleFunc("/api/v1/monitoring/bots/{botId}/history", h.GetBotHistory).Methods("GET")

	// Portfolio metrics endpoints
	router.HandleFunc("/api/v1/monitoring/portfolio", h.GetPortfolioMetrics).Methods("GET")
	router.HandleFunc("/api/v1/monitoring/portfolio/performance", h.GetPortfolioPerformance).Methods("GET")
	router.HandleFunc("/api/v1/monitoring/portfolio/risk", h.GetPortfolioRisk).Methods("GET")

	// System metrics endpoints
	router.HandleFunc("/api/v1/monitoring/system", h.GetSystemMetrics).Methods("GET")
	router.HandleFunc("/api/v1/monitoring/system/health", h.GetSystemHealth).Methods("GET")

	// Alert endpoints
	router.HandleFunc("/api/v1/monitoring/alerts", h.GetAlerts).Methods("GET")
	router.HandleFunc("/api/v1/monitoring/alerts/{alertId}/acknowledge", h.AcknowledgeAlert).Methods("POST")
	router.HandleFunc("/api/v1/monitoring/alerts/{alertId}/resolve", h.ResolveAlert).Methods("POST")

	// Analytics endpoints
	router.HandleFunc("/api/v1/monitoring/analytics/performance", h.GetPerformanceAnalytics).Methods("GET")
	router.HandleFunc("/api/v1/monitoring/analytics/risk", h.GetRiskAnalytics).Methods("GET")
	router.HandleFunc("/api/v1/monitoring/analytics/trading", h.GetTradingAnalytics).Methods("GET")

	// Real-time endpoints
	router.HandleFunc("/api/v1/monitoring/realtime/metrics", h.GetRealtimeMetrics).Methods("GET")
	router.HandleFunc("/api/v1/monitoring/realtime/alerts", h.GetRealtimeAlerts).Methods("GET")
}

// GetDashboard handles GET /api/v1/monitoring/dashboard
func (h *MonitoringHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dashboardData := h.monitor.GetDashboardData()

	h.logger.Info(ctx, "Dashboard data retrieved", map[string]interface{}{
		"total_bots":      dashboardData.Overview.TotalBots,
		"active_bots":     dashboardData.Overview.ActiveBots,
		"active_alerts":   dashboardData.Overview.ActiveAlerts,
		"portfolio_value": dashboardData.Overview.TotalValue.String(),
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dashboardData)
}

// GetOverview handles GET /api/v1/monitoring/overview
func (h *MonitoringHandler) GetOverview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dashboardData := h.monitor.GetDashboardData()

	response := map[string]interface{}{
		"overview":          dashboardData.Overview,
		"portfolio_metrics": dashboardData.PortfolioMetrics,
		"system_status":     dashboardData.SystemStatus,
		"recent_alerts":     dashboardData.RecentAlerts,
		"last_updated":      dashboardData.LastUpdated,
	}

	h.logger.Info(ctx, "Overview data retrieved", nil)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAllBotMetrics handles GET /api/v1/monitoring/bots
func (h *MonitoringHandler) GetAllBotMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	botMetrics := h.monitor.GetAllBotMetrics()

	response := map[string]interface{}{
		"bots":      botMetrics,
		"count":     len(botMetrics),
		"timestamp": time.Now(),
	}

	h.logger.Info(ctx, "All bot metrics retrieved", map[string]interface{}{
		"bot_count": len(botMetrics),
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetBotMetrics handles GET /api/v1/monitoring/bots/{botId}
func (h *MonitoringHandler) GetBotMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	botID := vars["botId"]

	metrics, err := h.monitor.GetBotMetrics(botID)
	if err != nil {
		h.logger.Error(ctx, "Failed to get bot metrics", err, map[string]interface{}{
			"bot_id": botID,
		})
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.logger.Info(ctx, "Bot metrics retrieved", map[string]interface{}{
		"bot_id":   botID,
		"strategy": metrics.Strategy,
		"state":    metrics.State,
		"health":   string(metrics.Health.OverallHealth),
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// GetBotPerformance handles GET /api/v1/monitoring/bots/{botId}/performance
func (h *MonitoringHandler) GetBotPerformance(w http.ResponseWriter, r *http.Request) {
	_ = r.Context()
	vars := mux.Vars(r)
	botID := vars["botId"]

	metrics, err := h.monitor.GetBotMetrics(botID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"bot_id":      botID,
		"performance": metrics.Performance,
		"risk":        metrics.Risk,
		"timestamp":   metrics.Timestamp,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetBotHealth handles GET /api/v1/monitoring/bots/{botId}/health
func (h *MonitoringHandler) GetBotHealth(w http.ResponseWriter, r *http.Request) {
	_ = r.Context()
	vars := mux.Vars(r)
	botID := vars["botId"]

	metrics, err := h.monitor.GetBotMetrics(botID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"bot_id":    botID,
		"health":    metrics.Health,
		"alerts":    metrics.Alerts,
		"system":    metrics.System,
		"timestamp": metrics.Timestamp,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetBotHistory handles GET /api/v1/monitoring/bots/{botId}/history
func (h *MonitoringHandler) GetBotHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	botID := vars["botId"]

	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	limit := 100 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	history, err := h.monitor.GetPerformanceHistory(botID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Limit results
	if len(history) > limit {
		history = history[len(history)-limit:]
	}

	response := map[string]interface{}{
		"bot_id":  botID,
		"history": history,
		"count":   len(history),
		"limit":   limit,
	}

	h.logger.Info(ctx, "Bot history retrieved", map[string]interface{}{
		"bot_id": botID,
		"count":  len(history),
		"limit":  limit,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetPortfolioMetrics handles GET /api/v1/monitoring/portfolio
func (h *MonitoringHandler) GetPortfolioMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	portfolioMetrics := h.monitor.GetPortfolioMetrics()

	h.logger.Info(ctx, "Portfolio metrics retrieved", map[string]interface{}{
		"total_value": portfolioMetrics.TotalValue.String(),
		"total_pnl":   portfolioMetrics.TotalPnL.String(),
		"active_bots": portfolioMetrics.ActiveBots,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(portfolioMetrics)
}

// GetPortfolioPerformance handles GET /api/v1/monitoring/portfolio/performance
func (h *MonitoringHandler) GetPortfolioPerformance(w http.ResponseWriter, r *http.Request) {
	_ = r.Context()

	portfolioMetrics := h.monitor.GetPortfolioMetrics()

	response := map[string]interface{}{
		"total_value":        portfolioMetrics.TotalValue,
		"total_pnl":          portfolioMetrics.TotalPnL,
		"daily_pnl":          portfolioMetrics.DailyPnL,
		"total_return":       portfolioMetrics.TotalReturn,
		"sharpe_ratio":       portfolioMetrics.SharpeRatio,
		"max_drawdown":       portfolioMetrics.MaxDrawdown,
		"strategy_breakdown": portfolioMetrics.StrategyBreakdown,
		"asset_allocation":   portfolioMetrics.AssetAllocation,
		"timestamp":          portfolioMetrics.Timestamp,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetPortfolioRisk handles GET /api/v1/monitoring/portfolio/risk
func (h *MonitoringHandler) GetPortfolioRisk(w http.ResponseWriter, r *http.Request) {
	_ = r.Context()

	portfolioMetrics := h.monitor.GetPortfolioMetrics()

	response := map[string]interface{}{
		"var_95":             portfolioMetrics.VaR95,
		"max_drawdown":       portfolioMetrics.MaxDrawdown,
		"strategy_breakdown": portfolioMetrics.StrategyBreakdown,
		"asset_allocation":   portfolioMetrics.AssetAllocation,
		"active_bots":        portfolioMetrics.ActiveBots,
		"total_bots":         portfolioMetrics.TotalBots,
		"timestamp":          portfolioMetrics.Timestamp,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetSystemMetrics handles GET /api/v1/monitoring/system
func (h *MonitoringHandler) GetSystemMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	systemMetrics := h.monitor.GetSystemMetrics()

	h.logger.Info(ctx, "System metrics retrieved", map[string]interface{}{
		"cpu_usage":    systemMetrics.CPUUsage,
		"memory_usage": systemMetrics.MemoryUsage,
		"error_rate":   systemMetrics.ErrorRate,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(systemMetrics)
}

// GetSystemHealth handles GET /api/v1/monitoring/system/health
func (h *MonitoringHandler) GetSystemHealth(w http.ResponseWriter, r *http.Request) {
	_ = r.Context()

	dashboardData := h.monitor.GetDashboardData()

	response := map[string]interface{}{
		"system_status":   dashboardData.SystemStatus,
		"overview":        dashboardData.Overview,
		"active_alerts":   len(dashboardData.RecentAlerts),
		"critical_alerts": dashboardData.Overview.CriticalAlerts,
		"timestamp":       time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAlerts handles GET /api/v1/monitoring/alerts
func (h *MonitoringHandler) GetAlerts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	severity := r.URL.Query().Get("severity")
	botID := r.URL.Query().Get("bot_id")

	alerts := h.monitor.GetActiveAlerts()

	// Filter alerts based on query parameters
	filteredAlerts := make([]*monitoring.Alert, 0)
	for _, alert := range alerts {
		if severity != "" && string(alert.Severity) != severity {
			continue
		}
		if botID != "" && alert.BotID != botID {
			continue
		}
		filteredAlerts = append(filteredAlerts, alert)
	}

	response := map[string]interface{}{
		"alerts": filteredAlerts,
		"count":  len(filteredAlerts),
		"total":  len(alerts),
		"filters": map[string]string{
			"severity": severity,
			"bot_id":   botID,
		},
		"timestamp": time.Now(),
	}

	h.logger.Info(ctx, "Alerts retrieved", map[string]interface{}{
		"total_alerts":    len(alerts),
		"filtered_alerts": len(filteredAlerts),
		"severity_filter": severity,
		"bot_id_filter":   botID,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// AcknowledgeAlert handles POST /api/v1/monitoring/alerts/{alertId}/acknowledge
func (h *MonitoringHandler) AcknowledgeAlert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	alertID := vars["alertId"]

	if err := h.monitor.AcknowledgeAlert(alertID); err != nil {
		h.logger.Error(ctx, "Failed to acknowledge alert", err, map[string]interface{}{
			"alert_id": alertID,
		})
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.logger.Info(ctx, "Alert acknowledged", map[string]interface{}{
		"alert_id": alertID,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Alert acknowledged successfully",
		"alert_id": alertID,
	})
}

// ResolveAlert handles POST /api/v1/monitoring/alerts/{alertId}/resolve
func (h *MonitoringHandler) ResolveAlert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	alertID := vars["alertId"]

	if err := h.monitor.ResolveAlert(alertID); err != nil {
		h.logger.Error(ctx, "Failed to resolve alert", err, map[string]interface{}{
			"alert_id": alertID,
		})
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.logger.Info(ctx, "Alert resolved", map[string]interface{}{
		"alert_id": alertID,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Alert resolved successfully",
		"alert_id": alertID,
	})
}

// Placeholder implementations for analytics endpoints
func (h *MonitoringHandler) GetPerformanceAnalytics(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *MonitoringHandler) GetRiskAnalytics(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *MonitoringHandler) GetTradingAnalytics(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *MonitoringHandler) GetRealtimeMetrics(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *MonitoringHandler) GetRealtimeAlerts(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
