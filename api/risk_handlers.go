package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ai-agentic-browser/internal/hft"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/gorilla/mux"
)

// RiskHandlers provides HTTP handlers for Advanced Risk Management
type RiskHandlers struct {
	riskManager *hft.AdvancedRiskManager
	logger      *observability.Logger
}

// NewRiskHandlers creates new risk management HTTP handlers
func NewRiskHandlers(riskManager *hft.AdvancedRiskManager, logger *observability.Logger) *RiskHandlers {
	return &RiskHandlers{
		riskManager: riskManager,
		logger:      logger,
	}
}

// Risk Management API handlers

// handleRiskLimits handles risk limits management
func (s *APIServer) handleRiskLimits(w http.ResponseWriter, r *http.Request) {
	if s.advancedRiskManager == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Advanced Risk Manager not available")
		return
	}

	riskHandlers := NewRiskHandlers(s.advancedRiskManager, s.logger)

	switch r.Method {
	case "GET":
		riskHandlers.GetLimits(w, r)
	case "POST", "PUT":
		riskHandlers.UpdateLimits(w, r)
	}
}

// GetLimits handles risk limits requests
func (h *RiskHandlers) GetLimits(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limits := h.riskManager.GetLimits()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(limits); err != nil {
		h.logger.Error(ctx, "Failed to encode limits", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// UpdateLimits handles risk limits update requests
func (h *RiskHandlers) UpdateLimits(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var limits hft.RiskLimits
	if err := json.NewDecoder(r.Body).Decode(&limits); err != nil {
		h.logger.Error(ctx, "Failed to decode limits update request", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.riskManager.UpdateLimits(ctx, &limits); err != nil {
		h.logger.Error(ctx, "Failed to update risk limits", err)
		http.Error(w, "Failed to update limits: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Risk limits updated successfully",
	})
}

// GetMetrics handles risk metrics requests
func (h *RiskHandlers) GetMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	metrics := h.riskManager.GetMetrics()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		h.logger.Error(ctx, "Failed to encode risk metrics", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetViolations handles risk violations requests
func (h *RiskHandlers) GetViolations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse limit parameter
	limit := 50 // Default limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	violations := h.riskManager.GetViolations(limit)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(violations); err != nil {
		h.logger.Error(ctx, "Failed to encode violations", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// ValidateOrder handles order validation requests
func (h *RiskHandlers) ValidateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var orderReq hft.OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
		h.logger.Error(ctx, "Failed to decode order validation request", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the order
	err := h.riskManager.ValidateOrder(ctx, &orderReq)

	response := map[string]interface{}{
		"order_id": orderReq.ID.String(),
		"symbol":   orderReq.Symbol,
		"valid":    err == nil,
	}

	if err != nil {
		response["error"] = err.Error()
		response["reason"] = "Risk validation failed"
	} else {
		response["message"] = "Order passed risk validation"
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error(ctx, "Failed to encode validation response", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// EmergencyStop handles emergency stop requests
func (h *RiskHandlers) EmergencyStop(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Reason string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(ctx, "Failed to decode emergency stop request", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Reason == "" {
		req.Reason = "Manual emergency stop via API"
	}

	if err := h.riskManager.EmergencyStop(ctx, req.Reason); err != nil {
		h.logger.Error(ctx, "Failed to trigger emergency stop", err)
		http.Error(w, "Failed to trigger emergency stop: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Emergency stop activated",
		"reason":  req.Reason,
	})
}

// handleGetRiskLimits returns all risk limits
func (s *APIServer) handleGetRiskLimits(w http.ResponseWriter, r *http.Request) {
	limits := []map[string]interface{}{
		{
			"id":         "daily_loss_limit",
			"name":       "Daily Loss Limit",
			"type":       "LOSS_LIMIT",
			"current":    2500.0,
			"limit":      5000.0,
			"percentage": 50.0,
			"status":     "OK",
			"enabled":    true,
			"created_at": "2024-01-10T08:00:00Z",
			"updated_at": "2024-01-15T10:30:00Z",
		},
		{
			"id":         "position_size_limit",
			"name":       "Position Size Limit",
			"type":       "POSITION_LIMIT",
			"current":    0.8,
			"limit":      1.0,
			"percentage": 80.0,
			"status":     "OK",
			"enabled":    true,
			"created_at": "2024-01-10T08:00:00Z",
			"updated_at": "2024-01-15T10:30:00Z",
		},
		{
			"id":         "max_drawdown_limit",
			"name":       "Max Drawdown Limit",
			"type":       "DRAWDOWN_LIMIT",
			"current":    1200.0,
			"limit":      2000.0,
			"percentage": 60.0,
			"status":     "WARNING",
			"enabled":    true,
			"created_at": "2024-01-10T08:00:00Z",
			"updated_at": "2024-01-15T10:30:00Z",
		},
		{
			"id":         "order_rate_limit",
			"name":       "Order Rate Limit",
			"type":       "RATE_LIMIT",
			"current":    45.0,
			"limit":      60.0,
			"percentage": 75.0,
			"status":     "OK",
			"enabled":    true,
			"created_at": "2024-01-10T08:00:00Z",
			"updated_at": "2024-01-15T10:30:00Z",
		},
	}

	s.sendJSON(w, r, http.StatusOK, map[string]interface{}{
		"limits": limits,
		"count":  len(limits),
		"summary": map[string]interface{}{
			"total_limits":    len(limits),
			"ok_limits":       3,
			"warning_limits":  1,
			"critical_limits": 0,
		},
	})
}

// handleCreateRiskLimit creates a new risk limit
func (s *APIServer) handleCreateRiskLimit(w http.ResponseWriter, r *http.Request) {
	var limitRequest struct {
		Name    string  `json:"name"`
		Type    string  `json:"type"`
		Limit   float64 `json:"limit"`
		Enabled bool    `json:"enabled"`
	}

	if err := json.NewDecoder(r.Body).Decode(&limitRequest); err != nil {
		s.sendError(w, r, http.StatusBadRequest, "Invalid limit format")
		return
	}

	// Validate required fields
	if limitRequest.Name == "" || limitRequest.Type == "" || limitRequest.Limit <= 0 {
		s.sendError(w, r, http.StatusBadRequest, "Missing or invalid required fields")
		return
	}

	// Mock limit creation
	limit := map[string]interface{}{
		"id":         "limit_" + time.Now().Format("20060102150405"),
		"name":       limitRequest.Name,
		"type":       limitRequest.Type,
		"current":    0.0,
		"limit":      limitRequest.Limit,
		"percentage": 0.0,
		"status":     "OK",
		"enabled":    limitRequest.Enabled,
		"created_at": time.Now().Format(time.RFC3339),
		"updated_at": time.Now().Format(time.RFC3339),
	}

	s.sendJSON(w, r, http.StatusCreated, limit)

	// Broadcast risk update
	s.BroadcastMessage("risk_limit_created", limit)
}

// handleRiskLimit handles individual risk limit operations
func (s *APIServer) handleRiskLimit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	limitID := vars["id"]

	switch r.Method {
	case "GET":
		// Mock limit data
		limit := map[string]interface{}{
			"id":         limitID,
			"name":       "Daily Loss Limit",
			"type":       "LOSS_LIMIT",
			"current":    2500.0,
			"limit":      5000.0,
			"percentage": 50.0,
			"status":     "OK",
			"enabled":    true,
			"created_at": "2024-01-10T08:00:00Z",
			"updated_at": "2024-01-15T10:30:00Z",
			"history": []map[string]interface{}{
				{
					"timestamp": "2024-01-15T10:30:00Z",
					"value":     2500.0,
					"status":    "OK",
				},
				{
					"timestamp": "2024-01-15T10:00:00Z",
					"value":     2300.0,
					"status":    "OK",
				},
			},
		}

		s.sendJSON(w, r, http.StatusOK, limit)

	case "PUT":
		var updateRequest struct {
			Name    string   `json:"name,omitempty"`
			Limit   *float64 `json:"limit,omitempty"`
			Enabled *bool    `json:"enabled,omitempty"`
		}

		if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
			s.sendError(w, r, http.StatusBadRequest, "Invalid update format")
			return
		}

		// Mock limit update
		s.sendJSON(w, r, http.StatusOK, map[string]interface{}{
			"id":         limitID,
			"message":    "Risk limit updated successfully",
			"updated_at": time.Now().Format(time.RFC3339),
		})

		// Broadcast risk update
		s.BroadcastMessage("risk_limit_updated", map[string]interface{}{
			"id":         limitID,
			"updated_at": time.Now().Format(time.RFC3339),
		})

	case "DELETE":
		// Mock limit deletion
		s.sendJSON(w, r, http.StatusOK, map[string]interface{}{
			"id":      limitID,
			"message": "Risk limit deleted successfully",
		})

		// Broadcast risk update
		s.BroadcastMessage("risk_limit_deleted", map[string]interface{}{
			"id": limitID,
		})
	}
}

// Legacy handleRiskViolations - kept for compatibility
func (s *APIServer) handleLegacyRiskViolations(w http.ResponseWriter, r *http.Request) {
	violations := []map[string]interface{}{
		{
			"id":           "violation_1",
			"type":         "Position Size",
			"symbol":       "BTCUSDT",
			"severity":     "HIGH",
			"message":      "Position size exceeded 80% of limit",
			"current":      0.85,
			"limit":        1.0,
			"percentage":   85.0,
			"timestamp":    "2024-01-15T14:30:25Z",
			"acknowledged": false,
			"resolved":     false,
		},
		{
			"id":           "violation_2",
			"type":         "Daily Loss",
			"symbol":       "ALL",
			"severity":     "MEDIUM",
			"message":      "Daily loss approaching 60% of limit",
			"current":      2800.0,
			"limit":        5000.0,
			"percentage":   56.0,
			"timestamp":    "2024-01-15T13:45:10Z",
			"acknowledged": true,
			"resolved":     false,
		},
	}

	s.sendJSON(w, r, http.StatusOK, map[string]interface{}{
		"violations": violations,
		"count":      len(violations),
		"summary": map[string]interface{}{
			"total":        len(violations),
			"high":         1,
			"medium":       1,
			"low":          0,
			"unresolved":   2,
			"acknowledged": 1,
		},
	})
}

// Legacy handleRiskMetrics - kept for compatibility
func (s *APIServer) handleLegacyRiskMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := map[string]interface{}{
		"overall_risk_score": 3.2,
		"risk_level":         "MODERATE",
		"var_metrics": map[string]interface{}{
			"var_95":             "-1250.30",
			"var_99":             "-2150.45",
			"expected_shortfall": "-2850.60",
			"confidence_level":   0.95,
		},
		"portfolio_metrics": map[string]interface{}{
			"beta":               0.85,
			"alpha":              0.12,
			"sharpe_ratio":       1.85,
			"sortino_ratio":      2.14,
			"max_drawdown":       "-2150.30",
			"volatility":         0.18,
			"correlation_market": 0.75,
		},
		"concentration_risk": map[string]interface{}{
			"largest_position":    0.35,
			"top_3_positions":     0.68,
			"top_5_positions":     0.82,
			"herfindahl_index":    0.28,
			"concentration_score": "MODERATE",
		},
		"liquidity_risk": map[string]interface{}{
			"avg_bid_ask_spread":  0.0025,
			"market_impact_score": 2.1,
			"liquidity_score":     "GOOD",
		},
		"operational_risk": map[string]interface{}{
			"system_uptime":     99.95,
			"avg_latency_ms":    2.3,
			"error_rate":        0.001,
			"operational_score": "EXCELLENT",
		},
		"stress_test": map[string]interface{}{
			"market_crash_scenario": "-15420.50",
			"liquidity_crisis":      "-8750.25",
			"volatility_spike":      "-5230.80",
			"worst_case_scenario":   "-18950.75",
		},
		"last_update": time.Now().Format(time.RFC3339),
	}

	s.sendJSON(w, r, http.StatusOK, metrics)
}

// handleEmergencyStop triggers emergency stop procedures
func (s *APIServer) handleEmergencyStop(w http.ResponseWriter, r *http.Request) {
	var stopRequest struct {
		Reason         string `json:"reason"`
		StopTrading    bool   `json:"stop_trading"`
		CancelOrders   bool   `json:"cancel_orders"`
		ClosePositions bool   `json:"close_positions"`
	}

	if err := json.NewDecoder(r.Body).Decode(&stopRequest); err != nil {
		s.sendError(w, r, http.StatusBadRequest, "Invalid emergency stop request")
		return
	}

	// Mock emergency stop execution
	actions := []string{}
	if stopRequest.StopTrading {
		actions = append(actions, "Trading halted")
	}
	if stopRequest.CancelOrders {
		actions = append(actions, "All orders canceled")
	}
	if stopRequest.ClosePositions {
		actions = append(actions, "All positions closed")
	}

	response := map[string]interface{}{
		"status":    "EMERGENCY_STOP_ACTIVATED",
		"reason":    stopRequest.Reason,
		"actions":   actions,
		"timestamp": time.Now().Format(time.RFC3339),
		"operator":  "system", // In real implementation, would be the authenticated user
	}

	s.sendJSON(w, r, http.StatusOK, response)

	// Broadcast emergency stop
	s.BroadcastMessage("emergency_stop", response)

	// Log emergency stop
	s.logger.Warn(r.Context(), "Emergency stop activated", map[string]interface{}{
		"reason":          stopRequest.Reason,
		"stop_trading":    stopRequest.StopTrading,
		"cancel_orders":   stopRequest.CancelOrders,
		"close_positions": stopRequest.ClosePositions,
	})
}

// TradingView API handlers

// handleTradingViewCharts handles TradingView charts management
func (s *APIServer) handleTradingViewCharts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.handleGetTradingViewCharts(w, r)
	case "POST":
		s.handleCreateTradingViewChart(w, r)
	}
}

// handleGetTradingViewCharts returns all TradingView charts
func (s *APIServer) handleGetTradingViewCharts(w http.ResponseWriter, r *http.Request) {
	charts := []map[string]interface{}{
		{
			"id":          "chart_btcusdt_1m",
			"symbol":      "BTCUSDT",
			"timeframe":   "1m",
			"url":         "https://www.tradingview.com/chart/?symbol=BINANCE:BTCUSDT",
			"last_update": "2024-01-15T10:30:00Z",
			"indicators":  []string{"RSI", "MACD", "BB"},
			"signals":     3,
		},
		{
			"id":          "chart_ethusdt_5m",
			"symbol":      "ETHUSDT",
			"timeframe":   "5m",
			"url":         "https://www.tradingview.com/chart/?symbol=BINANCE:ETHUSDT",
			"last_update": "2024-01-15T10:29:00Z",
			"indicators":  []string{"RSI", "EMA"},
			"signals":     1,
		},
	}

	s.sendJSON(w, r, http.StatusOK, map[string]interface{}{
		"charts": charts,
		"count":  len(charts),
	})
}

// handleCreateTradingViewChart creates a new TradingView chart
func (s *APIServer) handleCreateTradingViewChart(w http.ResponseWriter, r *http.Request) {
	var chartRequest struct {
		Symbol    string `json:"symbol"`
		Timeframe string `json:"timeframe"`
	}

	if err := json.NewDecoder(r.Body).Decode(&chartRequest); err != nil {
		s.sendError(w, r, http.StatusBadRequest, "Invalid chart request")
		return
	}

	// Validate required fields
	if chartRequest.Symbol == "" || chartRequest.Timeframe == "" {
		s.sendError(w, r, http.StatusBadRequest, "Missing required fields")
		return
	}

	// Mock chart creation
	chart := map[string]interface{}{
		"id":          "chart_" + chartRequest.Symbol + "_" + chartRequest.Timeframe,
		"symbol":      chartRequest.Symbol,
		"timeframe":   chartRequest.Timeframe,
		"url":         "https://www.tradingview.com/chart/?symbol=BINANCE:" + chartRequest.Symbol,
		"last_update": time.Now().Format(time.RFC3339),
		"indicators":  []string{},
		"signals":     0,
	}

	s.sendJSON(w, r, http.StatusCreated, chart)

	// Broadcast chart update
	s.BroadcastMessage("tradingview_chart_created", chart)
}

// handleTradingViewChart handles individual chart operations
func (s *APIServer) handleTradingViewChart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chartID := vars["id"]

	switch r.Method {
	case "GET":
		// Mock chart data
		chart := map[string]interface{}{
			"id":          chartID,
			"symbol":      "BTCUSDT",
			"timeframe":   "1m",
			"url":         "https://www.tradingview.com/chart/?symbol=BINANCE:BTCUSDT",
			"last_update": "2024-01-15T10:30:00Z",
			"indicators":  []string{"RSI", "MACD", "BB"},
			"signals":     3,
			"ohlcv": []map[string]interface{}{
				{
					"timestamp": "2024-01-15T10:30:00Z",
					"open":      "45120.00",
					"high":      "45135.50",
					"low":       "45110.25",
					"close":     "45123.45",
					"volume":    "12.345",
				},
			},
		}

		s.sendJSON(w, r, http.StatusOK, chart)

	case "DELETE":
		// Mock chart deletion
		s.sendJSON(w, r, http.StatusOK, map[string]interface{}{
			"id":      chartID,
			"message": "Chart deleted successfully",
		})

		// Broadcast chart update
		s.BroadcastMessage("tradingview_chart_deleted", map[string]interface{}{
			"id": chartID,
		})
	}
}

// handleTradingViewSignals returns TradingView signals
func (s *APIServer) handleTradingViewSignals(w http.ResponseWriter, r *http.Request) {
	signals := []map[string]interface{}{
		{
			"id":          "tv_signal_1",
			"symbol":      "BTCUSDT",
			"type":        "TECHNICAL",
			"direction":   "BUY",
			"strength":    "STRONG",
			"price":       "45123.45",
			"timestamp":   "2024-01-15T10:30:00Z",
			"timeframe":   "1m",
			"indicators":  []string{"RSI", "MACD"},
			"description": "RSI oversold signal with MACD bullish crossover",
			"confidence":  0.85,
		},
		{
			"id":          "tv_signal_2",
			"symbol":      "ETHUSDT",
			"type":        "PATTERN",
			"direction":   "SELL",
			"strength":    "MODERATE",
			"price":       "3012.34",
			"timestamp":   "2024-01-15T10:28:00Z",
			"timeframe":   "5m",
			"indicators":  []string{"BB"},
			"description": "Bollinger Bands squeeze breakout to downside",
			"confidence":  0.72,
		},
	}

	s.sendJSON(w, r, http.StatusOK, map[string]interface{}{
		"signals": signals,
		"count":   len(signals),
	})
}

// handleTradingViewIndicators returns indicator values for a symbol
func (s *APIServer) handleTradingViewIndicators(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	symbol := vars["symbol"]

	indicators := map[string]interface{}{
		"symbol":    symbol,
		"timestamp": time.Now().Format(time.RFC3339),
		"indicators": map[string]interface{}{
			"RSI": map[string]interface{}{
				"value":     45.67,
				"signal":    "NEUTRAL",
				"timestamp": time.Now().Format(time.RFC3339),
			},
			"MACD": map[string]interface{}{
				"macd":      0.5,
				"signal":    0.3,
				"histogram": 0.2,
				"direction": "BULLISH",
				"timestamp": time.Now().Format(time.RFC3339),
			},
			"BB": map[string]interface{}{
				"upper":     "45200.00",
				"middle":    "45123.45",
				"lower":     "45050.00",
				"position":  "MIDDLE",
				"timestamp": time.Now().Format(time.RFC3339),
			},
		},
	}

	s.sendJSON(w, r, http.StatusOK, indicators)
}
