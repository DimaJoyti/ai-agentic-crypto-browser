package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ai-agentic-browser/internal/hft"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// DashboardHandlers provides HTTP handlers for Real-Time Dashboard
type DashboardHandlers struct {
	dashboard *hft.RealtimeDashboard
	logger    *observability.Logger
}

// NewDashboardHandlers creates new real-time dashboard HTTP handlers
func NewDashboardHandlers(dashboard *hft.RealtimeDashboard, logger *observability.Logger) *DashboardHandlers {
	return &DashboardHandlers{
		dashboard: dashboard,
		logger:    logger,
	}
}

// GetLiveMetrics handles live metrics requests
func (h *DashboardHandlers) GetLiveMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	metrics := h.dashboard.GetLiveMetrics()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		h.logger.Error(ctx, "Failed to encode live metrics", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetAlerts handles alerts requests
func (h *DashboardHandlers) GetAlerts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	severityStr := r.URL.Query().Get("severity")

	limit := 50 // Default limit
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	var severity hft.AlertSeverity
	if severityStr != "" {
		severity = hft.AlertSeverity(severityStr)
	}

	alerts := h.dashboard.GetAlerts(limit, severity)

	response := map[string]interface{}{
		"alerts": alerts,
		"count":  len(alerts),
		"filters": map[string]interface{}{
			"limit":    limit,
			"severity": severityStr,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error(ctx, "Failed to encode alerts", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// AcknowledgeAlert handles alert acknowledgment requests
func (h *DashboardHandlers) AcknowledgeAlert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	alertIDStr := vars["id"]
	if alertIDStr == "" {
		http.Error(w, "Alert ID is required", http.StatusBadRequest)
		return
	}

	alertID, err := uuid.Parse(alertIDStr)
	if err != nil {
		http.Error(w, "Invalid alert ID", http.StatusBadRequest)
		return
	}

	var req struct {
		AcknowledgedBy string `json:"acknowledged_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(ctx, "Failed to decode acknowledge request", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.AcknowledgedBy == "" {
		http.Error(w, "acknowledged_by is required", http.StatusBadRequest)
		return
	}

	if err := h.dashboard.AcknowledgeAlert(ctx, alertID, req.AcknowledgedBy); err != nil {
		h.logger.Error(ctx, "Failed to acknowledge alert", err)
		http.Error(w, "Failed to acknowledge alert: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Alert acknowledged successfully",
	})

	h.logger.Info(ctx, "Alert acknowledged via API", map[string]interface{}{
		"alert_id":        alertID.String(),
		"acknowledged_by": req.AcknowledgedBy,
	})
}

// ResolveAlert handles alert resolution requests
func (h *DashboardHandlers) ResolveAlert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	alertIDStr := vars["id"]
	if alertIDStr == "" {
		http.Error(w, "Alert ID is required", http.StatusBadRequest)
		return
	}

	alertID, err := uuid.Parse(alertIDStr)
	if err != nil {
		http.Error(w, "Invalid alert ID", http.StatusBadRequest)
		return
	}

	if err := h.dashboard.ResolveAlert(ctx, alertID); err != nil {
		h.logger.Error(ctx, "Failed to resolve alert", err)
		http.Error(w, "Failed to resolve alert: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Alert resolved successfully",
	})

	h.logger.Info(ctx, "Alert resolved via API", map[string]interface{}{
		"alert_id": alertID.String(),
	})
}

// CreateSession handles session creation requests
func (h *DashboardHandlers) CreateSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(ctx, "Failed to decode session request", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	session, err := h.dashboard.CreateSession(ctx, req.UserID)
	if err != nil {
		h.logger.Error(ctx, "Failed to create session", err)
		http.Error(w, "Failed to create session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "success",
		"message":    "Session created successfully",
		"session_id": session.ID.String(),
		"session":    session,
	})

	h.logger.Info(ctx, "Dashboard session created via API", map[string]interface{}{
		"session_id": session.ID.String(),
		"user_id":    req.UserID,
	})
}

// GetSession handles session retrieval requests
func (h *DashboardHandlers) GetSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	sessionID := vars["id"]
	if sessionID == "" {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	session := h.dashboard.GetSession(sessionID)
	if session == nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(session); err != nil {
		h.logger.Error(ctx, "Failed to encode session", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// UpdateSession handles session update requests
func (h *DashboardHandlers) UpdateSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	sessionID := vars["id"]
	if sessionID == "" {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		h.logger.Error(ctx, "Failed to decode session updates", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.dashboard.UpdateSession(ctx, sessionID, updates); err != nil {
		h.logger.Error(ctx, "Failed to update session", err)
		http.Error(w, "Failed to update session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Session updated successfully",
	})

	h.logger.Info(ctx, "Dashboard session updated via API", map[string]interface{}{
		"session_id": sessionID,
		"updates":    len(updates),
	})
}

// CloseSession handles session close requests
func (h *DashboardHandlers) CloseSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	sessionID := vars["id"]
	if sessionID == "" {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	if err := h.dashboard.CloseSession(ctx, sessionID); err != nil {
		h.logger.Error(ctx, "Failed to close session", err)
		http.Error(w, "Failed to close session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Session closed successfully",
	})

	h.logger.Info(ctx, "Dashboard session closed via API", map[string]interface{}{
		"session_id": sessionID,
	})
}

// GetWidgets handles widgets requests
func (h *DashboardHandlers) GetWidgets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	widgets := h.dashboard.GetWidgets()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"widgets": widgets,
		"count":   len(widgets),
	}); err != nil {
		h.logger.Error(ctx, "Failed to encode widgets", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetLayouts handles layouts requests
func (h *DashboardHandlers) GetLayouts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	layouts := h.dashboard.GetLayouts()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"layouts": layouts,
		"count":   len(layouts),
	}); err != nil {
		h.logger.Error(ctx, "Failed to encode layouts", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetStatus handles dashboard status requests
func (h *DashboardHandlers) GetStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status := h.dashboard.GetDashboardStatus()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		h.logger.Error(ctx, "Failed to encode dashboard status", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
