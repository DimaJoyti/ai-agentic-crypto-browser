package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Portfolio API handlers

// handlePortfolioSummary returns portfolio summary
func (s *APIServer) handlePortfolioSummary(w http.ResponseWriter, r *http.Request) {
	summary := map[string]interface{}{
		"total_value":    "125420.50",
		"cash_balance":   "25000.00",
		"total_pnl":      "15420.50",
		"unrealized_pnl": "387.05",
		"realized_pnl":   "15033.45",
		"day_pnl":        "1234.56",
		"open_positions": 2,
		"total_trades":   1247,
		"win_rate":       68.5,
		"sharpe_ratio":   1.85,
		"max_drawdown":   "2150.30",
		"last_update":    time.Now().Format(time.RFC3339),
	}

	s.sendJSON(w, r, http.StatusOK, summary)
}

// handlePortfolioPositions returns current portfolio positions
func (s *APIServer) handlePortfolioPositions(w http.ResponseWriter, r *http.Request) {
	positions := []map[string]interface{}{
		{
			"symbol":         "BTCUSDT",
			"size":           "0.5",
			"avg_price":      "44500.00",
			"current_price":  "45123.45",
			"unrealized_pnl": "311.73",
			"realized_pnl":   "0.00",
			"commission":     "2.25",
			"open_time":      "2024-01-15T09:30:00Z",
			"update_time":    "2024-01-15T10:30:00Z",
			"exchange":       "binance",
			"strategy_id":    "market_making_1",
		},
		{
			"symbol":         "ETHUSDT",
			"size":           "-2.0",
			"avg_price":      "3050.00",
			"current_price":  "3012.34",
			"unrealized_pnl": "75.32",
			"realized_pnl":   "0.00",
			"commission":     "6.10",
			"open_time":      "2024-01-15T09:45:00Z",
			"update_time":    "2024-01-15T10:30:00Z",
			"exchange":       "binance",
			"strategy_id":    "arbitrage_1",
		},
	}

	s.sendJSON(w, r, http.StatusOK, map[string]interface{}{
		"positions":   positions,
		"count":       len(positions),
		"total_value": "100420.50",
		"total_pnl":   "387.05",
	})
}

// handlePortfolioMetrics returns detailed portfolio metrics
func (s *APIServer) handlePortfolioMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := map[string]interface{}{
		"total_value":            "125420.50",
		"cash_balance":           "25000.00",
		"total_pnl":              "15420.50",
		"unrealized_pnl":         "387.05",
		"realized_pnl":           "15033.45",
		"day_pnl":                "1234.56",
		"max_drawdown":           "2150.30",
		"high_water_mark":        "127570.80",
		"open_positions":         2,
		"total_trades":           1247,
		"winning_trades":         854,
		"losing_trades":          393,
		"win_rate":               68.5,
		"avg_win":                "45.67",
		"avg_loss":               "-23.45",
		"profit_factor":          1.95,
		"sharpe_ratio":           1.85,
		"sortino_ratio":          2.14,
		"calmar_ratio":           7.17,
		"max_consecutive_wins":   12,
		"max_consecutive_losses": 5,
		"last_update":            time.Now().Format(time.RFC3339),
	}

	s.sendJSON(w, r, http.StatusOK, metrics)
}

// handlePortfolioRisk returns portfolio risk metrics
func (s *APIServer) handlePortfolioRisk(w http.ResponseWriter, r *http.Request) {
	risk := map[string]interface{}{
		"var_95":             "-1250.30",
		"var_99":             "-2150.45",
		"expected_shortfall": "-2850.60",
		"sharpe_ratio":       1.85,
		"sortino_ratio":      2.14,
		"max_drawdown":       "-2150.30",
		"beta":               0.85,
		"alpha":              0.12,
		"volatility":         0.18,
		"correlation_btc":    0.75,
		"correlation_eth":    0.68,
		"concentration_risk": 0.35,
		"leverage":           1.2,
		"margin_usage":       0.15,
		"last_update":        time.Now().Format(time.RFC3339),
	}

	s.sendJSON(w, r, http.StatusOK, risk)
}

// Strategy API handlers

// handleStrategies handles strategy management
func (s *APIServer) handleStrategies(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.handleGetStrategies(w, r)
	case "POST":
		s.handleCreateStrategy(w, r)
	}
}

// handleGetStrategies returns all strategies
func (s *APIServer) handleGetStrategies(w http.ResponseWriter, r *http.Request) {
	strategies := []map[string]interface{}{
		{
			"id":      "market_making_1",
			"name":    "Market Making Strategy",
			"type":    "MARKET_MAKING",
			"enabled": true,
			"symbols": []string{"BTCUSDT", "ETHUSDT"},
			"performance": map[string]interface{}{
				"total_pnl":    "8420.30",
				"total_trades": 856,
				"win_rate":     72.1,
				"sharpe_ratio": 2.15,
				"max_drawdown": "450.20",
			},
			"parameters": map[string]interface{}{
				"spread_bps":    10,
				"order_size":    0.1,
				"max_inventory": 1.0,
			},
			"status":     "RUNNING",
			"created_at": "2024-01-10T08:00:00Z",
			"updated_at": "2024-01-15T10:30:00Z",
		},
		{
			"id":      "arbitrage_1",
			"name":    "Cross-Exchange Arbitrage",
			"type":    "ARBITRAGE",
			"enabled": true,
			"symbols": []string{"BTCUSDT"},
			"performance": map[string]interface{}{
				"total_pnl":    "5200.15",
				"total_trades": 234,
				"win_rate":     89.2,
				"sharpe_ratio": 3.45,
				"max_drawdown": "125.50",
			},
			"parameters": map[string]interface{}{
				"min_profit_bps":    20,
				"max_position_size": 1.0,
				"execution_timeout": 5000,
			},
			"status":     "RUNNING",
			"created_at": "2024-01-12T10:00:00Z",
			"updated_at": "2024-01-15T10:30:00Z",
		},
		{
			"id":      "momentum_1",
			"name":    "Momentum Strategy",
			"type":    "MOMENTUM",
			"enabled": false,
			"symbols": []string{"ETHUSDT", "BNBUSDT"},
			"performance": map[string]interface{}{
				"total_pnl":    "-123.45",
				"total_trades": 157,
				"win_rate":     45.8,
				"sharpe_ratio": -0.25,
				"max_drawdown": "567.80",
			},
			"parameters": map[string]interface{}{
				"lookback_period":    20,
				"momentum_threshold": 0.02,
				"stop_loss":          0.05,
			},
			"status":     "STOPPED",
			"created_at": "2024-01-08T12:00:00Z",
			"updated_at": "2024-01-14T16:00:00Z",
		},
	}

	s.sendJSON(w, r, http.StatusOK, map[string]interface{}{
		"strategies": strategies,
		"count":      len(strategies),
		"active":     2,
		"total_pnl":  "13497.00",
	})
}

// handleCreateStrategy creates a new strategy
func (s *APIServer) handleCreateStrategy(w http.ResponseWriter, r *http.Request) {
	var strategyRequest struct {
		Name       string                 `json:"name"`
		Type       string                 `json:"type"`
		Symbols    []string               `json:"symbols"`
		Parameters map[string]interface{} `json:"parameters"`
		Enabled    bool                   `json:"enabled"`
	}

	if err := json.NewDecoder(r.Body).Decode(&strategyRequest); err != nil {
		s.sendError(w, r, http.StatusBadRequest, "Invalid strategy format")
		return
	}

	// Validate required fields
	if strategyRequest.Name == "" || strategyRequest.Type == "" {
		s.sendError(w, r, http.StatusBadRequest, "Missing required fields")
		return
	}

	// Mock strategy creation
	strategy := map[string]interface{}{
		"id":         "strategy_" + strconv.FormatInt(time.Now().Unix(), 10),
		"name":       strategyRequest.Name,
		"type":       strategyRequest.Type,
		"enabled":    strategyRequest.Enabled,
		"symbols":    strategyRequest.Symbols,
		"parameters": strategyRequest.Parameters,
		"performance": map[string]interface{}{
			"total_pnl":    "0.00",
			"total_trades": 0,
			"win_rate":     0.0,
		},
		"status":     "CREATED",
		"created_at": time.Now().Format(time.RFC3339),
		"updated_at": time.Now().Format(time.RFC3339),
	}

	s.sendJSON(w, r, http.StatusCreated, strategy)

	// Broadcast strategy update
	s.BroadcastMessage("strategy_created", strategy)
}

// handleStrategy handles individual strategy operations
func (s *APIServer) handleStrategy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["id"]

	switch r.Method {
	case "GET":
		// Mock strategy data
		strategy := map[string]interface{}{
			"id":      strategyID,
			"name":    "Market Making Strategy",
			"type":    "MARKET_MAKING",
			"enabled": true,
			"symbols": []string{"BTCUSDT", "ETHUSDT"},
			"performance": map[string]interface{}{
				"total_pnl":    "8420.30",
				"total_trades": 856,
				"win_rate":     72.1,
				"sharpe_ratio": 2.15,
				"max_drawdown": "450.20",
			},
			"parameters": map[string]interface{}{
				"spread_bps":    10,
				"order_size":    0.1,
				"max_inventory": 1.0,
			},
			"status":     "RUNNING",
			"created_at": "2024-01-10T08:00:00Z",
			"updated_at": "2024-01-15T10:30:00Z",
		}

		s.sendJSON(w, r, http.StatusOK, strategy)

	case "PUT":
		var updateRequest struct {
			Name       string                 `json:"name,omitempty"`
			Enabled    *bool                  `json:"enabled,omitempty"`
			Symbols    []string               `json:"symbols,omitempty"`
			Parameters map[string]interface{} `json:"parameters,omitempty"`
		}

		if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
			s.sendError(w, r, http.StatusBadRequest, "Invalid update format")
			return
		}

		// Mock strategy update
		s.sendJSON(w, r, http.StatusOK, map[string]interface{}{
			"id":         strategyID,
			"message":    "Strategy updated successfully",
			"updated_at": time.Now().Format(time.RFC3339),
		})

		// Broadcast strategy update
		s.BroadcastMessage("strategy_updated", map[string]interface{}{
			"id":         strategyID,
			"updated_at": time.Now().Format(time.RFC3339),
		})

	case "DELETE":
		// Mock strategy deletion
		s.sendJSON(w, r, http.StatusOK, map[string]interface{}{
			"id":      strategyID,
			"message": "Strategy deleted successfully",
		})

		// Broadcast strategy update
		s.BroadcastMessage("strategy_deleted", map[string]interface{}{
			"id": strategyID,
		})
	}
}

// handleStrategyStart starts a strategy
func (s *APIServer) handleStrategyStart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["id"]

	s.sendJSON(w, r, http.StatusOK, map[string]interface{}{
		"id":      strategyID,
		"status":  "RUNNING",
		"message": "Strategy started successfully",
	})

	// Broadcast strategy update
	s.BroadcastMessage("strategy_started", map[string]interface{}{
		"id":     strategyID,
		"status": "RUNNING",
	})
}

// handleStrategyStop stops a strategy
func (s *APIServer) handleStrategyStop(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["id"]

	s.sendJSON(w, r, http.StatusOK, map[string]interface{}{
		"id":      strategyID,
		"status":  "STOPPED",
		"message": "Strategy stopped successfully",
	})

	// Broadcast strategy update
	s.BroadcastMessage("strategy_stopped", map[string]interface{}{
		"id":     strategyID,
		"status": "STOPPED",
	})
}

// handleStrategyPerformance returns strategy performance metrics
func (s *APIServer) handleStrategyPerformance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["id"]

	performance := map[string]interface{}{
		"strategy_id":            strategyID,
		"total_trades":           856,
		"winning_trades":         617,
		"losing_trades":          239,
		"win_rate":               72.1,
		"total_pnl":              "8420.30",
		"realized_pnl":           "8033.25",
		"unrealized_pnl":         "387.05",
		"max_drawdown":           "450.20",
		"sharpe_ratio":           2.15,
		"sortino_ratio":          2.85,
		"average_win":            "25.67",
		"average_loss":           "-12.34",
		"profit_factor":          2.08,
		"max_consecutive_wins":   8,
		"max_consecutive_losses": 3,
		"last_update":            time.Now().Format(time.RFC3339),
		"daily_performance": []map[string]interface{}{
			{
				"date":     "2024-01-15",
				"pnl":      "234.56",
				"trades":   45,
				"win_rate": 73.3,
			},
			{
				"date":     "2024-01-14",
				"pnl":      "189.23",
				"trades":   38,
				"win_rate": 71.1,
			},
		},
	}

	s.sendJSON(w, r, http.StatusOK, performance)
}
