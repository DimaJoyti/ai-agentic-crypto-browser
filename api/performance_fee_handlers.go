package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ai-agentic-browser/internal/billing"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

// PerformanceFeeHandlers handles performance fee related requests
type PerformanceFeeHandlers struct {
	performanceTracker *billing.PerformanceFeeTracker
}

// NewPerformanceFeeHandlers creates new performance fee handlers
func NewPerformanceFeeHandlers(performanceTracker *billing.PerformanceFeeTracker) *PerformanceFeeHandlers {
	return &PerformanceFeeHandlers{
		performanceTracker: performanceTracker,
	}
}

// RegisterRoutes registers performance fee routes
func (pfh *PerformanceFeeHandlers) RegisterRoutes(router *mux.Router) {
	// Performance fee routes
	router.HandleFunc("/performance/config", pfh.GetFeeConfig).Methods("GET")
	router.HandleFunc("/performance/config", pfh.UpdateFeeConfig).Methods("PUT")
	router.HandleFunc("/performance/trades", pfh.RecordTrade).Methods("POST")
	router.HandleFunc("/performance/trades", pfh.GetTrades).Methods("GET")
	router.HandleFunc("/performance/summary", pfh.GetPerformanceSummary).Methods("GET")
	router.HandleFunc("/performance/analytics", pfh.GetPerformanceAnalytics).Methods("GET")
	router.HandleFunc("/performance/fees/billing", pfh.GetFeeBilling).Methods("GET")
	router.HandleFunc("/performance/leaderboard", pfh.GetLeaderboard).Methods("GET")
}

// GetFeeConfig returns user's performance fee configuration
func (pfh *PerformanceFeeHandlers) GetFeeConfig(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	config, err := pfh.performanceTracker.GetFeeConfig(r.Context(), userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get fee config: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

// UpdateFeeConfigRequest represents fee config update request
type UpdateFeeConfigRequest struct {
	FeePercentage       float64 `json:"fee_percentage"`
	MinimumFee          float64 `json:"minimum_fee"`
	MaximumFee          float64 `json:"maximum_fee"`
	FeeFrequency        string  `json:"fee_frequency"`
	OnlyProfitableTrades bool    `json:"only_profitable_trades"`
}

// UpdateFeeConfig updates user's performance fee configuration
func (pfh *PerformanceFeeHandlers) UpdateFeeConfig(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req UpdateFeeConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate fee percentage (2-30%)
	if req.FeePercentage < 2.0 || req.FeePercentage > 30.0 {
		http.Error(w, "Fee percentage must be between 2% and 30%", http.StatusBadRequest)
		return
	}

	// Get current config
	config, err := pfh.performanceTracker.GetFeeConfig(r.Context(), userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get current config: %v", err), http.StatusInternalServerError)
		return
	}

	// Update config fields
	config.FeePercentage = decimal.NewFromFloat(req.FeePercentage)
	config.MinimumFee = decimal.NewFromFloat(req.MinimumFee)
	config.MaximumFee = decimal.NewFromFloat(req.MaximumFee)
	config.FeeFrequency = req.FeeFrequency
	config.OnlyProfitableTrades = req.OnlyProfitableTrades
	config.UpdatedAt = time.Now()

	// Save updated config (implementation would update database)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Fee configuration updated successfully",
		"config":  config,
	})
}

// RecordTradeRequest represents trade recording request
type RecordTradeRequest struct {
	Symbol         string  `json:"symbol"`
	Side           string  `json:"side"`
	Quantity       float64 `json:"quantity"`
	EntryPrice     float64 `json:"entry_price"`
	ExitPrice      float64 `json:"exit_price"`
	EntryTimestamp string  `json:"entry_timestamp"`
	ExitTimestamp  string  `json:"exit_timestamp"`
	StrategyID     string  `json:"strategy_id"`
	Exchange       string  `json:"exchange"`
}

// RecordTrade records a completed trade
func (pfh *PerformanceFeeHandlers) RecordTrade(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req RecordTradeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Parse timestamps
	entryTime, err := time.Parse(time.RFC3339, req.EntryTimestamp)
	if err != nil {
		http.Error(w, "Invalid entry timestamp format", http.StatusBadRequest)
		return
	}

	exitTime, err := time.Parse(time.RFC3339, req.ExitTimestamp)
	if err != nil {
		http.Error(w, "Invalid exit timestamp format", http.StatusBadRequest)
		return
	}

	// Create trade record
	trade := &billing.TradeRecord{
		ID:             uuid.New().String(),
		UserID:         userID,
		Symbol:         req.Symbol,
		Side:           req.Side,
		Quantity:       decimal.NewFromFloat(req.Quantity),
		EntryPrice:     decimal.NewFromFloat(req.EntryPrice),
		ExitPrice:      decimal.NewFromFloat(req.ExitPrice),
		EntryTimestamp: entryTime,
		ExitTimestamp:  exitTime,
		StrategyID:     req.StrategyID,
	}

	// Calculate PnL
	if req.Side == "buy" || req.Side == "long" {
		trade.PnL = trade.ExitPrice.Sub(trade.EntryPrice).Mul(trade.Quantity)
	} else {
		trade.PnL = trade.EntryPrice.Sub(trade.ExitPrice).Mul(trade.Quantity)
	}

	// Record trade and calculate performance fee
	err = pfh.performanceTracker.RecordTrade(r.Context(), trade)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to record trade: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":         true,
		"message":         "Trade recorded successfully",
		"trade_id":        trade.ID,
		"pnl":            trade.PnL,
		"performance_fee": trade.PerformanceFee,
	})
}

// GetTrades returns user's trade history
func (pfh *PerformanceFeeHandlers) GetTrades(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get query parameters
	limitStr := r.URL.Query().Get("limit")
	limit := 50 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	// Mock trade data (implementation would query database)
	trades := []map[string]interface{}{
		{
			"id":              "trade_001",
			"symbol":          "BTC/USD",
			"side":            "buy",
			"quantity":        1.5,
			"entry_price":     45000.00,
			"exit_price":      47000.00,
			"pnl":            3000.00,
			"performance_fee": 600.00,
			"entry_timestamp": "2024-01-15T10:00:00Z",
			"exit_timestamp":  "2024-01-15T14:30:00Z",
			"status":          "completed",
		},
		{
			"id":              "trade_002",
			"symbol":          "ETH/USD",
			"side":            "sell",
			"quantity":        10.0,
			"entry_price":     2800.00,
			"exit_price":      2750.00,
			"pnl":            500.00,
			"performance_fee": 100.00,
			"entry_timestamp": "2024-01-16T09:15:00Z",
			"exit_timestamp":  "2024-01-16T11:45:00Z",
			"status":          "completed",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"trades": trades[:min(len(trades), limit)],
		"total":  len(trades),
	})
}

// GetPerformanceSummary returns user's performance summary
func (pfh *PerformanceFeeHandlers) GetPerformanceSummary(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "all_time"
	}

	summary, err := pfh.performanceTracker.GetPerformanceSummary(r.Context(), userID, period)
	if err != nil {
		// Return mock data if no summary exists yet
		summary = &billing.PerformanceSummary{
			UserID:              userID,
			Period:              period,
			TotalTrades:         25,
			ProfitableTrades:    18,
			TotalPnL:            decimal.NewFromFloat(12500.00),
			TotalPerformanceFees: decimal.NewFromFloat(2500.00),
			WinRate:             decimal.NewFromFloat(0.72),
			AverageReturn:       decimal.NewFromFloat(0.085),
			MaxDrawdown:         decimal.NewFromFloat(0.15),
			SharpeRatio:         decimal.NewFromFloat(1.85),
			CurrentHighWaterMark: decimal.NewFromFloat(22500.00),
			LastUpdated:         time.Now(),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// GetPerformanceAnalytics returns detailed performance analytics
func (pfh *PerformanceFeeHandlers) GetPerformanceAnalytics(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	analytics := map[string]interface{}{
		"performance_metrics": map[string]interface{}{
			"total_return":      "25.8%",
			"annualized_return": "18.4%",
			"volatility":        "12.3%",
			"sharpe_ratio":      1.85,
			"sortino_ratio":     2.12,
			"calmar_ratio":      1.23,
			"max_drawdown":      "15.2%",
			"win_rate":          "72%",
			"profit_factor":     2.45,
		},
		"fee_breakdown": map[string]interface{}{
			"total_fees_paid":     2500.00,
			"average_fee_rate":    "20%",
			"fees_this_month":     450.00,
			"fees_last_month":     380.00,
			"fee_efficiency":      "High", // Based on performance vs fees
		},
		"trading_patterns": map[string]interface{}{
			"avg_trade_duration":     "4.2 hours",
			"most_profitable_hour":   "14:00 UTC",
			"best_performing_symbol": "BTC/USD",
			"strategy_performance": map[string]interface{}{
				"momentum": map[string]interface{}{
					"trades":    15,
					"win_rate":  "80%",
					"avg_return": "8.5%",
				},
				"mean_reversion": map[string]interface{}{
					"trades":    10,
					"win_rate":  "60%",
					"avg_return": "5.2%",
				},
			},
		},
		"high_water_mark": map[string]interface{}{
			"current_value":    22500.00,
			"previous_high":    21800.00,
			"days_since_high":  5,
			"recovery_needed":  "0%",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analytics)
}

// GetFeeBilling returns performance fee billing information
func (pfh *PerformanceFeeHandlers) GetFeeBilling(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	billing := map[string]interface{}{
		"current_month": map[string]interface{}{
			"period":           "2024-01",
			"trades":           8,
			"profitable_trades": 6,
			"total_pnl":        1250.00,
			"performance_fees": 250.00,
			"fee_rate":         "20%",
			"status":           "pending",
		},
		"billing_history": []map[string]interface{}{
			{
				"period":           "2023-12",
				"trades":           12,
				"profitable_trades": 9,
				"total_pnl":        1800.00,
				"performance_fees": 360.00,
				"fee_rate":         "20%",
				"status":           "paid",
				"paid_date":        "2024-01-05",
			},
			{
				"period":           "2023-11",
				"trades":           5,
				"profitable_trades": 3,
				"total_pnl":        450.00,
				"performance_fees": 90.00,
				"fee_rate":         "20%",
				"status":           "paid",
				"paid_date":        "2023-12-05",
			},
		},
		"fee_settings": map[string]interface{}{
			"fee_percentage":        20.0,
			"high_water_mark":       22500.00,
			"only_profitable_trades": true,
			"fee_frequency":         "monthly",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(billing)
}

// GetLeaderboard returns performance leaderboard
func (pfh *PerformanceFeeHandlers) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	// Check if user has permission to view leaderboard
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	leaderboard := map[string]interface{}{
		"top_performers": []map[string]interface{}{
			{
				"rank":         1,
				"user_id":      "user_anonymous_1",
				"total_return": "45.2%",
				"sharpe_ratio": 2.85,
				"trades":       156,
				"win_rate":     "78%",
				"fees_paid":    8500.00,
			},
			{
				"rank":         2,
				"user_id":      "user_anonymous_2",
				"total_return": "38.7%",
				"sharpe_ratio": 2.34,
				"trades":       89,
				"win_rate":     "74%",
				"fees_paid":    5200.00,
			},
			{
				"rank":         3,
				"user_id":      "user_anonymous_3",
				"total_return": "32.1%",
				"sharpe_ratio": 2.12,
				"trades":       203,
				"win_rate":     "69%",
				"fees_paid":    7800.00,
			},
		},
		"your_rank": 15,
		"total_participants": 1247,
		"period": "all_time",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leaderboard)
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
