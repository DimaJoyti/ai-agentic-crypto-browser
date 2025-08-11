package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ai-agentic-browser/internal/trading"
	"github.com/ai-agentic-browser/internal/trading/strategies"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

// TradingBotHandler handles trading bot API requests
type TradingBotHandler struct {
	logger          *observability.Logger
	botEngine       *trading.TradingBotEngine
	strategyManager *strategies.StrategyManager
}

// NewTradingBotHandler creates a new trading bot handler
func NewTradingBotHandler(logger *observability.Logger, botEngine *trading.TradingBotEngine, strategyManager *strategies.StrategyManager) *TradingBotHandler {
	return &TradingBotHandler{
		logger:          logger,
		botEngine:       botEngine,
		strategyManager: strategyManager,
	}
}

// RegisterRoutes registers trading bot API routes
func (h *TradingBotHandler) RegisterRoutes(router *mux.Router) {
	// Bot management endpoints
	router.HandleFunc("/api/v1/trading-bots", h.ListBots).Methods("GET")
	router.HandleFunc("/api/v1/trading-bots", h.CreateBot).Methods("POST")
	router.HandleFunc("/api/v1/trading-bots/{botId}", h.GetBot).Methods("GET")
	router.HandleFunc("/api/v1/trading-bots/{botId}", h.UpdateBot).Methods("PUT")
	router.HandleFunc("/api/v1/trading-bots/{botId}", h.DeleteBot).Methods("DELETE")

	// Bot control endpoints
	router.HandleFunc("/api/v1/trading-bots/{botId}/start", h.StartBot).Methods("POST")
	router.HandleFunc("/api/v1/trading-bots/{botId}/stop", h.StopBot).Methods("POST")
	router.HandleFunc("/api/v1/trading-bots/{botId}/pause", h.PauseBot).Methods("POST")
	router.HandleFunc("/api/v1/trading-bots/{botId}/resume", h.ResumeBot).Methods("POST")

	// Bot monitoring endpoints
	router.HandleFunc("/api/v1/trading-bots/{botId}/status", h.GetBotStatus).Methods("GET")
	router.HandleFunc("/api/v1/trading-bots/{botId}/performance", h.GetBotPerformance).Methods("GET")
	router.HandleFunc("/api/v1/trading-bots/{botId}/trades", h.GetBotTrades).Methods("GET")

	// Strategy management endpoints
	router.HandleFunc("/api/v1/trading-strategies", h.ListStrategies).Methods("GET")
	router.HandleFunc("/api/v1/trading-strategies/{strategyId}", h.GetStrategy).Methods("GET")
	router.HandleFunc("/api/v1/trading-strategies/{strategyId}/performance", h.GetStrategyPerformance).Methods("GET")

	// Bulk operations
	router.HandleFunc("/api/v1/trading-bots/start-all", h.StartAllBots).Methods("POST")
	router.HandleFunc("/api/v1/trading-bots/stop-all", h.StopAllBots).Methods("POST")
	router.HandleFunc("/api/v1/trading-bots/performance", h.GetAllBotsPerformance).Methods("GET")
}

// CreateBotRequest represents a request to create a new trading bot
type CreateBotRequest struct {
	Name           string                 `json:"name"`
	Strategy       string                 `json:"strategy"`
	TradingPairs   []string               `json:"trading_pairs"`
	Exchange       string                 `json:"exchange"`
	BaseCurrency   string                 `json:"base_currency"`
	StrategyParams map[string]interface{} `json:"strategy_params"`
	RiskProfile    *BotRiskProfileRequest `json:"risk_profile"`
	Capital        *CapitalConfigRequest  `json:"capital"`
	Enabled        bool                   `json:"enabled"`
}

// BotRiskProfileRequest represents risk profile configuration
type BotRiskProfileRequest struct {
	MaxPositionSize decimal.Decimal `json:"max_position_size"`
	StopLoss        decimal.Decimal `json:"stop_loss"`
	TakeProfit      decimal.Decimal `json:"take_profit"`
	MaxDrawdown     decimal.Decimal `json:"max_drawdown"`
}

// CapitalConfigRequest represents capital configuration
type CapitalConfigRequest struct {
	InitialBalance       decimal.Decimal `json:"initial_balance"`
	AllocationPercentage decimal.Decimal `json:"allocation_percentage"`
}

// BotResponse represents a trading bot response
type BotResponse struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	Strategy    string                  `json:"strategy"`
	State       string                  `json:"state"`
	Config      *BotConfigResponse      `json:"config"`
	Performance *BotPerformanceResponse `json:"performance"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
}

// BotConfigResponse represents bot configuration response
type BotConfigResponse struct {
	TradingPairs   []string               `json:"trading_pairs"`
	Exchange       string                 `json:"exchange"`
	BaseCurrency   string                 `json:"base_currency"`
	StrategyParams map[string]interface{} `json:"strategy_params"`
	Capital        *CapitalConfigRequest  `json:"capital"`
	Enabled        bool                   `json:"enabled"`
}

// BotPerformanceResponse represents bot performance response
type BotPerformanceResponse struct {
	TotalTrades   int             `json:"total_trades"`
	WinningTrades int             `json:"winning_trades"`
	LosingTrades  int             `json:"losing_trades"`
	WinRate       decimal.Decimal `json:"win_rate"`
	TotalProfit   decimal.Decimal `json:"total_profit"`
	TotalLoss     decimal.Decimal `json:"total_loss"`
	NetProfit     decimal.Decimal `json:"net_profit"`
	MaxDrawdown   decimal.Decimal `json:"max_drawdown"`
	SharpeRatio   decimal.Decimal `json:"sharpe_ratio"`
	LastUpdated   time.Time       `json:"last_updated"`
}

// ListBots handles GET /api/v1/trading-bots
func (h *TradingBotHandler) ListBots(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	bots := h.botEngine.ListBots()

	var response []BotResponse
	for _, bot := range bots {
		response = append(response, h.convertBotToResponse(bot))
	}

	h.logger.Info(ctx, "Listed trading bots", map[string]interface{}{
		"count": len(response),
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"bots":  response,
		"count": len(response),
	})
}

// CreateBot handles POST /api/v1/trading-bots
func (h *TradingBotHandler) CreateBot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateBotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(ctx, "Failed to decode create bot request", err, nil)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := h.validateCreateBotRequest(&req); err != nil {
		h.logger.Error(ctx, "Invalid create bot request", err, nil)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create bot configuration
	botConfig := &trading.BotConfig{
		TradingPairs:   req.TradingPairs,
		Exchange:       req.Exchange,
		BaseCurrency:   req.BaseCurrency,
		StrategyParams: req.StrategyParams,
		Capital: &trading.CapitalConfig{
			InitialBalance:       req.Capital.InitialBalance,
			AllocationPercentage: req.Capital.AllocationPercentage,
		},
		Enabled: req.Enabled,
	}

	// Register bot with engine
	bot, err := h.botEngine.RegisterBot(ctx, botConfig, trading.BotStrategy(req.Strategy))
	if err != nil {
		h.logger.Error(ctx, "Failed to register bot", err, map[string]interface{}{
			"strategy": req.Strategy,
		})
		http.Error(w, "Failed to create bot", http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, "Trading bot created", map[string]interface{}{
		"bot_id":   bot.ID,
		"strategy": req.Strategy,
		"pairs":    req.TradingPairs,
	})

	response := h.convertBotToResponse(bot)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetBot handles GET /api/v1/trading-bots/{botId}
func (h *TradingBotHandler) GetBot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	botID := vars["botId"]

	bot, err := h.botEngine.GetBot(botID)
	if err != nil {
		h.logger.Error(ctx, "Bot not found", err, map[string]interface{}{
			"bot_id": botID,
		})
		http.Error(w, "Bot not found", http.StatusNotFound)
		return
	}

	response := h.convertBotToResponse(bot)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// StartBot handles POST /api/v1/trading-bots/{botId}/start
func (h *TradingBotHandler) StartBot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	botID := vars["botId"]

	if err := h.botEngine.StartBot(ctx, botID); err != nil {
		h.logger.Error(ctx, "Failed to start bot", err, map[string]interface{}{
			"bot_id": botID,
		})
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Info(ctx, "Bot started", map[string]interface{}{
		"bot_id": botID,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Bot started successfully",
		"bot_id":  botID,
	})
}

// StopBot handles POST /api/v1/trading-bots/{botId}/stop
func (h *TradingBotHandler) StopBot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	botID := vars["botId"]

	if err := h.botEngine.StopBot(ctx, botID); err != nil {
		h.logger.Error(ctx, "Failed to stop bot", err, map[string]interface{}{
			"bot_id": botID,
		})
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Info(ctx, "Bot stopped", map[string]interface{}{
		"bot_id": botID,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Bot stopped successfully",
		"bot_id":  botID,
	})
}

// GetBotStatus handles GET /api/v1/trading-bots/{botId}/status
func (h *TradingBotHandler) GetBotStatus(w http.ResponseWriter, r *http.Request) {
	_ = r.Context()
	vars := mux.Vars(r)
	botID := vars["botId"]

	bot, err := h.botEngine.GetBot(botID)
	if err != nil {
		http.Error(w, "Bot not found", http.StatusNotFound)
		return
	}

	status := map[string]interface{}{
		"bot_id":       bot.ID,
		"name":         bot.Name,
		"state":        string(bot.State),
		"strategy":     string(bot.Strategy),
		"is_active":    bot.State == trading.StateRunning,
		"last_updated": time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// GetBotPerformance handles GET /api/v1/trading-bots/{botId}/performance
func (h *TradingBotHandler) GetBotPerformance(w http.ResponseWriter, r *http.Request) {
	_ = r.Context()
	vars := mux.Vars(r)
	botID := vars["botId"]

	bot, err := h.botEngine.GetBot(botID)
	if err != nil {
		http.Error(w, "Bot not found", http.StatusNotFound)
		return
	}

	performance := h.convertPerformanceToResponse(bot.Performance)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(performance)
}

// Helper methods

// validateCreateBotRequest validates the create bot request
func (h *TradingBotHandler) validateCreateBotRequest(req *CreateBotRequest) error {
	if req.Name == "" {
		return fmt.Errorf("bot name is required")
	}

	if req.Strategy == "" {
		return fmt.Errorf("strategy is required")
	}

	if len(req.TradingPairs) == 0 {
		return fmt.Errorf("at least one trading pair is required")
	}

	if req.Exchange == "" {
		return fmt.Errorf("exchange is required")
	}

	if req.Capital == nil {
		return fmt.Errorf("capital configuration is required")
	}

	if req.Capital.InitialBalance.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("initial balance must be positive")
	}

	return nil
}

// convertBotToResponse converts a trading bot to response format
func (h *TradingBotHandler) convertBotToResponse(bot *trading.TradingBot) BotResponse {
	return BotResponse{
		ID:       bot.ID,
		Name:     bot.Name,
		Strategy: string(bot.Strategy),
		State:    string(bot.State),
		Config: &BotConfigResponse{
			TradingPairs:   bot.Config.TradingPairs,
			Exchange:       bot.Config.Exchange,
			BaseCurrency:   bot.Config.BaseCurrency,
			StrategyParams: bot.Config.StrategyParams,
			Capital: &CapitalConfigRequest{
				InitialBalance:       bot.Config.Capital.InitialBalance,
				AllocationPercentage: bot.Config.Capital.AllocationPercentage,
			},
			Enabled: bot.Config.Enabled,
		},
		Performance: h.convertPerformanceToResponse(bot.Performance),
		CreatedAt:   time.Now(), // This should come from the bot
		UpdatedAt:   time.Now(), // This should come from the bot
	}
}

// convertPerformanceToResponse converts performance to response format
func (h *TradingBotHandler) convertPerformanceToResponse(perf *trading.BotPerformance) *BotPerformanceResponse {
	if perf == nil {
		return &BotPerformanceResponse{}
	}

	return &BotPerformanceResponse{
		TotalTrades:   perf.TotalTrades,
		WinningTrades: perf.WinningTrades,
		LosingTrades:  perf.LosingTrades,
		WinRate:       perf.WinRate,
		TotalProfit:   perf.TotalProfit,
		TotalLoss:     perf.TotalLoss,
		NetProfit:     perf.NetProfit,
		MaxDrawdown:   perf.MaxDrawdown,
		SharpeRatio:   perf.SharpeRatio,
		LastUpdated:   perf.LastUpdated,
	}
}

// Placeholder implementations for remaining endpoints
func (h *TradingBotHandler) UpdateBot(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *TradingBotHandler) DeleteBot(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *TradingBotHandler) PauseBot(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *TradingBotHandler) ResumeBot(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *TradingBotHandler) GetBotTrades(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *TradingBotHandler) ListStrategies(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *TradingBotHandler) GetStrategy(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *TradingBotHandler) GetStrategyPerformance(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *TradingBotHandler) StartAllBots(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *TradingBotHandler) StopAllBots(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *TradingBotHandler) GetAllBotsPerformance(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
