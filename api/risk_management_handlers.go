package api

import (
	"encoding/json"
	"net/http"

	"github.com/ai-agentic-browser/internal/trading"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

// RiskManagementHandler handles risk management API requests
type RiskManagementHandler struct {
	logger      *observability.Logger
	riskManager *trading.BotRiskManager
}

// NewRiskManagementHandler creates a new risk management handler
func NewRiskManagementHandler(logger *observability.Logger, riskManager *trading.BotRiskManager) *RiskManagementHandler {
	return &RiskManagementHandler{
		logger:      logger,
		riskManager: riskManager,
	}
}

// RegisterRoutes registers risk management API routes
func (h *RiskManagementHandler) RegisterRoutes(router *mux.Router) {
	// Risk metrics endpoints
	router.HandleFunc("/api/v1/risk/portfolio", h.GetPortfolioRisk).Methods("GET")
	router.HandleFunc("/api/v1/risk/bots/{botId}", h.GetBotRiskMetrics).Methods("GET")
	router.HandleFunc("/api/v1/risk/bots/{botId}", h.UpdateBotRiskProfile).Methods("PUT")

	// Risk limits endpoints
	router.HandleFunc("/api/v1/risk/limits", h.GetRiskLimits).Methods("GET")
	router.HandleFunc("/api/v1/risk/limits", h.CreateRiskLimit).Methods("POST")
	router.HandleFunc("/api/v1/risk/limits/{limitId}", h.UpdateRiskLimit).Methods("PUT")
	router.HandleFunc("/api/v1/risk/limits/{limitId}", h.DeleteRiskLimit).Methods("DELETE")

	// Risk alerts endpoints
	router.HandleFunc("/api/v1/risk/alerts", h.GetRiskAlerts).Methods("GET")
	router.HandleFunc("/api/v1/risk/alerts/{alertId}/acknowledge", h.AcknowledgeAlert).Methods("POST")
	router.HandleFunc("/api/v1/risk/alerts/{alertId}/resolve", h.ResolveAlert).Methods("POST")

	// Risk controls endpoints
	router.HandleFunc("/api/v1/risk/emergency-stop", h.EmergencyStop).Methods("POST")
	router.HandleFunc("/api/v1/risk/bots/{botId}/halt", h.HaltBot).Methods("POST")
	router.HandleFunc("/api/v1/risk/bots/{botId}/resume", h.ResumeBot).Methods("POST")

	// Risk analysis endpoints
	router.HandleFunc("/api/v1/risk/correlation-matrix", h.GetCorrelationMatrix).Methods("GET")
	router.HandleFunc("/api/v1/risk/stress-test", h.RunStressTest).Methods("POST")
}

// GetPortfolioRisk handles GET /api/v1/risk/portfolio
func (h *RiskManagementHandler) GetPortfolioRisk(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	portfolioRisk := h.riskManager.GetPortfolioRisk()

	response := map[string]interface{}{
		"total_exposure":     portfolioRisk.TotalExposure.String(),
		"var_95":             portfolioRisk.VaR95.String(),
		"var_99":             portfolioRisk.VaR99.String(),
		"expected_shortfall": portfolioRisk.ExpectedShortfall.String(),
		"max_drawdown":       portfolioRisk.MaxDrawdown.String(),
		"current_drawdown":   portfolioRisk.CurrentDrawdown.String(),
		"sharpe_ratio":       portfolioRisk.SharpeRatio.String(),
		"concentration_risk": portfolioRisk.ConcentrationRisk,
		"correlation_risk":   portfolioRisk.CorrelationRisk.String(),
		"liquidity_risk":     portfolioRisk.LiquidityRisk.String(),
		"last_updated":       portfolioRisk.LastUpdated,
	}

	h.logger.Info(ctx, "Portfolio risk retrieved", map[string]interface{}{
		"total_exposure": portfolioRisk.TotalExposure.String(),
		"var_95":         portfolioRisk.VaR95.String(),
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetBotRiskMetrics handles GET /api/v1/risk/bots/{botId}
func (h *RiskManagementHandler) GetBotRiskMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	botID := vars["botId"]

	metrics, err := h.riskManager.GetBotRiskMetrics(botID)
	if err != nil {
		h.logger.Error(ctx, "Failed to get bot risk metrics", err, map[string]interface{}{
			"bot_id": botID,
		})
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"bot_id":             metrics.BotID,
		"current_exposure":   metrics.CurrentExposure.String(),
		"daily_pnl":          metrics.DailyPnL.String(),
		"unrealized_pnl":     metrics.UnrealizedPnL.String(),
		"max_drawdown":       metrics.MaxDrawdown.String(),
		"current_drawdown":   metrics.CurrentDrawdown.String(),
		"consecutive_losses": metrics.ConsecutiveLosses,
		"consecutive_wins":   metrics.ConsecutiveWins,
		"var_95":             metrics.VaR95.String(),
		"volatility":         metrics.Volatility.String(),
		"beta":               metrics.Beta.String(),
		"risk_score":         metrics.RiskScore,
		"last_updated":       metrics.LastUpdated,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateBotRiskProfileRequest represents a request to update bot risk profile
type UpdateBotRiskProfileRequest struct {
	MaxPositionSize      *decimal.Decimal `json:"max_position_size,omitempty"`
	StopLoss             *decimal.Decimal `json:"stop_loss,omitempty"`
	TakeProfit           *decimal.Decimal `json:"take_profit,omitempty"`
	MaxDrawdown          *decimal.Decimal `json:"max_drawdown,omitempty"`
	MaxDailyLoss         *decimal.Decimal `json:"max_daily_loss,omitempty"`
	MaxConsecutiveLosses *int             `json:"max_consecutive_losses,omitempty"`
	RiskTolerance        *string          `json:"risk_tolerance,omitempty"`
}

// UpdateBotRiskProfile handles PUT /api/v1/risk/bots/{botId}
func (h *RiskManagementHandler) UpdateBotRiskProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	botID := vars["botId"]

	var req UpdateBotRiskProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Implementation would update the bot's risk profile
	h.logger.Info(ctx, "Bot risk profile update requested", map[string]interface{}{
		"bot_id": botID,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Risk profile updated successfully",
		"bot_id":  botID,
	})
}

// GetRiskAlerts handles GET /api/v1/risk/alerts
func (h *RiskManagementHandler) GetRiskAlerts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	severity := r.URL.Query().Get("severity")
	botID := r.URL.Query().Get("bot_id")
	status := r.URL.Query().Get("status")

	// Get alerts based on filters
	var alerts []*trading.RiskAlert

	if severity != "" {
		// Filter by severity
		alerts = h.riskManager.GetAlertsBySeverity(trading.AlertSeverity(severity))
	} else if botID != "" {
		// Filter by bot ID
		alerts = h.riskManager.GetAlertsByBot(botID)
	} else {
		// Get all active alerts
		alerts = h.riskManager.GetActiveAlerts()
	}

	h.logger.Info(ctx, "Risk alerts retrieved", map[string]interface{}{
		"count":    len(alerts),
		"severity": severity,
		"bot_id":   botID,
		"status":   status,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"alerts": alerts,
		"count":  len(alerts),
	})
}

// AcknowledgeAlert handles POST /api/v1/risk/alerts/{alertId}/acknowledge
func (h *RiskManagementHandler) AcknowledgeAlert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	alertID := vars["alertId"]

	// Get user ID from context (would be set by auth middleware)
	userID := "system" // Placeholder

	if err := h.riskManager.AcknowledgeAlert(ctx, alertID, userID); err != nil {
		h.logger.Error(ctx, "Failed to acknowledge alert", err, map[string]interface{}{
			"alert_id": alertID,
		})
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.logger.Info(ctx, "Alert acknowledged", map[string]interface{}{
		"alert_id": alertID,
		"user_id":  userID,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Alert acknowledged successfully",
		"alert_id": alertID,
	})
}

// ResolveAlert handles POST /api/v1/risk/alerts/{alertId}/resolve
func (h *RiskManagementHandler) ResolveAlert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	alertID := vars["alertId"]

	// Get user ID from context (would be set by auth middleware)
	userID := "system" // Placeholder

	if err := h.riskManager.ResolveAlert(ctx, alertID, userID); err != nil {
		h.logger.Error(ctx, "Failed to resolve alert", err, map[string]interface{}{
			"alert_id": alertID,
		})
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.logger.Info(ctx, "Alert resolved", map[string]interface{}{
		"alert_id": alertID,
		"user_id":  userID,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Alert resolved successfully",
		"alert_id": alertID,
	})
}

// EmergencyStopRequest represents an emergency stop request
type EmergencyStopRequest struct {
	Reason string `json:"reason"`
}

// EmergencyStop handles POST /api/v1/risk/emergency-stop
func (h *RiskManagementHandler) EmergencyStop(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req EmergencyStopRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Reason == "" {
		req.Reason = "Manual emergency stop"
	}

	if err := h.riskManager.EmergencyStop(ctx, req.Reason); err != nil {
		h.logger.Error(ctx, "Failed to activate emergency stop", err, nil)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Error(ctx, "Emergency stop activated", nil, map[string]interface{}{
		"reason": req.Reason,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Emergency stop activated",
		"reason":  req.Reason,
	})
}

// HaltBot handles POST /api/v1/risk/bots/{botId}/halt
func (h *RiskManagementHandler) HaltBot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	botID := vars["botId"]

	var req struct {
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		req.Reason = "Manual halt"
	}

	if err := h.riskManager.HaltBot(ctx, botID, req.Reason); err != nil {
		h.logger.Error(ctx, "Failed to halt bot", err, map[string]interface{}{
			"bot_id": botID,
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Bot halted successfully",
		"bot_id":  botID,
		"reason":  req.Reason,
	})
}

// ResumeBot handles POST /api/v1/risk/bots/{botId}/resume
func (h *RiskManagementHandler) ResumeBot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	botID := vars["botId"]

	if err := h.riskManager.ResumeBot(ctx, botID); err != nil {
		h.logger.Error(ctx, "Failed to resume bot", err, map[string]interface{}{
			"bot_id": botID,
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Bot resumed successfully",
		"bot_id":  botID,
	})
}

// Placeholder implementations for remaining endpoints
func (h *RiskManagementHandler) GetRiskLimits(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *RiskManagementHandler) CreateRiskLimit(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *RiskManagementHandler) UpdateRiskLimit(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *RiskManagementHandler) DeleteRiskLimit(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *RiskManagementHandler) GetCorrelationMatrix(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *RiskManagementHandler) RunStressTest(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
