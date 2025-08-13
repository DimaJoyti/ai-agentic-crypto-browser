package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ai-agentic-browser/internal/hft"
	"github.com/ai-agentic-browser/pkg/observability"
)

// AnalyticsHandlers provides HTTP handlers for Post-Trade Analytics
type AnalyticsHandlers struct {
	analytics *hft.PostTradeAnalytics
	logger    *observability.Logger
}

// NewAnalyticsHandlers creates new post-trade analytics HTTP handlers
func NewAnalyticsHandlers(analytics *hft.PostTradeAnalytics, logger *observability.Logger) *AnalyticsHandlers {
	return &AnalyticsHandlers{
		analytics: analytics,
		logger:    logger,
	}
}

// GetExecutionMetrics handles execution metrics requests
func (h *AnalyticsHandlers) GetExecutionMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		http.Error(w, "Symbol parameter is required", http.StatusBadRequest)
		return
	}

	metrics := h.analytics.GetExecutionMetrics(symbol)
	if metrics == nil {
		http.Error(w, "No metrics found for symbol", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		h.logger.Error(ctx, "Failed to encode execution metrics", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetPerformanceData handles performance data requests
func (h *AnalyticsHandlers) GetPerformanceData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	strategyID := r.URL.Query().Get("strategy_id")
	if strategyID == "" {
		http.Error(w, "Strategy ID parameter is required", http.StatusBadRequest)
		return
	}

	performance := h.analytics.GetPerformanceData(strategyID)
	if performance == nil {
		http.Error(w, "No performance data found for strategy", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(performance); err != nil {
		h.logger.Error(ctx, "Failed to encode performance data", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetRealtimeMetrics handles real-time metrics requests
func (h *AnalyticsHandlers) GetRealtimeMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	metrics := h.analytics.GetRealtimeMetrics()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		h.logger.Error(ctx, "Failed to encode realtime metrics", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetTradeHistory handles trade history requests
func (h *AnalyticsHandlers) GetTradeHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	symbol := r.URL.Query().Get("symbol")
	strategyID := r.URL.Query().Get("strategy_id")

	limit := 100 // Default limit
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	trades := h.analytics.GetTradeHistory(limit, symbol, strategyID)

	response := map[string]interface{}{
		"trades": trades,
		"count":  len(trades),
		"filters": map[string]interface{}{
			"limit":       limit,
			"symbol":      symbol,
			"strategy_id": strategyID,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error(ctx, "Failed to encode trade history", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// RecordTrade handles trade recording requests
func (h *AnalyticsHandlers) RecordTrade(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var trade hft.TradeRecord
	if err := json.NewDecoder(r.Body).Decode(&trade); err != nil {
		h.logger.Error(ctx, "Failed to decode trade record", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.analytics.RecordTrade(ctx, &trade); err != nil {
		h.logger.Error(ctx, "Failed to record trade", err)
		http.Error(w, "Failed to record trade: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "success",
		"message":  "Trade recorded successfully",
		"trade_id": trade.TradeID.String(),
	})

	h.logger.Info(ctx, "Trade recorded via API", map[string]interface{}{
		"trade_id": trade.TradeID.String(),
		"symbol":   trade.Symbol,
		"side":     string(trade.Side),
		"quantity": trade.Quantity.String(),
	})
}

// GetAnalyticsSummary handles analytics summary requests
func (h *AnalyticsHandlers) GetAnalyticsSummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	realtimeMetrics := h.analytics.GetRealtimeMetrics()
	
	// Get sample execution metrics (first available symbol)
	var sampleExecutionMetrics *hft.ExecutionMetrics
	sampleSymbol := r.URL.Query().Get("sample_symbol")
	if sampleSymbol == "" {
		sampleSymbol = "BTCUSDT" // Default
	}
	sampleExecutionMetrics = h.analytics.GetExecutionMetrics(sampleSymbol)

	// Get sample performance data (first available strategy)
	var samplePerformanceData *hft.PerformanceData
	sampleStrategy := r.URL.Query().Get("sample_strategy")
	if sampleStrategy == "" {
		sampleStrategy = "strategy_1" // Default
	}
	samplePerformanceData = h.analytics.GetPerformanceData(sampleStrategy)

	// Get recent trades
	recentTrades := h.analytics.GetTradeHistory(10, "", "")

	summary := map[string]interface{}{
		"realtime_metrics":     realtimeMetrics,
		"execution_metrics":    sampleExecutionMetrics,
		"performance_data":     samplePerformanceData,
		"recent_trades_count":  len(recentTrades),
		"sample_symbol":        sampleSymbol,
		"sample_strategy":      sampleStrategy,
		"timestamp":            realtimeMetrics.LastUpdate,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(summary); err != nil {
		h.logger.Error(ctx, "Failed to encode analytics summary", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetSlippageAnalysis handles slippage analysis requests
func (h *AnalyticsHandlers) GetSlippageAnalysis(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	symbol := r.URL.Query().Get("symbol")
	limit := 100

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	trades := h.analytics.GetTradeHistory(limit, symbol, "")
	
	// Calculate slippage statistics
	var totalSlippage, minSlippage, maxSlippage float64
	var slippageCount int
	slippageDistribution := make(map[string]int)

	for _, trade := range trades {
		totalSlippage += trade.Slippage
		slippageCount++

		if slippageCount == 1 {
			minSlippage = trade.Slippage
			maxSlippage = trade.Slippage
		} else {
			if trade.Slippage < minSlippage {
				minSlippage = trade.Slippage
			}
			if trade.Slippage > maxSlippage {
				maxSlippage = trade.Slippage
			}
		}

		// Categorize slippage
		if trade.Slippage <= 5 {
			slippageDistribution["0-5 bps"]++
		} else if trade.Slippage <= 10 {
			slippageDistribution["5-10 bps"]++
		} else if trade.Slippage <= 20 {
			slippageDistribution["10-20 bps"]++
		} else {
			slippageDistribution[">20 bps"]++
		}
	}

	var avgSlippage float64
	if slippageCount > 0 {
		avgSlippage = totalSlippage / float64(slippageCount)
	}

	analysis := map[string]interface{}{
		"symbol":               symbol,
		"trade_count":          slippageCount,
		"avg_slippage":         avgSlippage,
		"min_slippage":         minSlippage,
		"max_slippage":         maxSlippage,
		"slippage_distribution": slippageDistribution,
		"analysis_period":      fmt.Sprintf("Last %d trades", limit),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(analysis); err != nil {
		h.logger.Error(ctx, "Failed to encode slippage analysis", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetMarketImpactAnalysis handles market impact analysis requests
func (h *AnalyticsHandlers) GetMarketImpactAnalysis(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	symbol := r.URL.Query().Get("symbol")
	limit := 100

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	trades := h.analytics.GetTradeHistory(limit, symbol, "")
	
	// Calculate market impact statistics
	var totalImpact, minImpact, maxImpact float64
	var impactCount int
	impactDistribution := make(map[string]int)

	for _, trade := range trades {
		totalImpact += trade.MarketImpact
		impactCount++

		if impactCount == 1 {
			minImpact = trade.MarketImpact
			maxImpact = trade.MarketImpact
		} else {
			if trade.MarketImpact < minImpact {
				minImpact = trade.MarketImpact
			}
			if trade.MarketImpact > maxImpact {
				maxImpact = trade.MarketImpact
			}
		}

		// Categorize impact
		if trade.MarketImpact <= 2 {
			impactDistribution["0-2 bps"]++
		} else if trade.MarketImpact <= 5 {
			impactDistribution["2-5 bps"]++
		} else if trade.MarketImpact <= 10 {
			impactDistribution["5-10 bps"]++
		} else {
			impactDistribution[">10 bps"]++
		}
	}

	var avgImpact float64
	if impactCount > 0 {
		avgImpact = totalImpact / float64(impactCount)
	}

	analysis := map[string]interface{}{
		"symbol":               symbol,
		"trade_count":          impactCount,
		"avg_market_impact":    avgImpact,
		"min_market_impact":    minImpact,
		"max_market_impact":    maxImpact,
		"impact_distribution":  impactDistribution,
		"analysis_period":      fmt.Sprintf("Last %d trades", limit),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(analysis); err != nil {
		h.logger.Error(ctx, "Failed to encode market impact analysis", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
