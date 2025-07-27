package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ai-agentic-browser/internal/ai"
	"github.com/ai-agentic-browser/internal/alerts"
	"github.com/ai-agentic-browser/internal/analytics"
	"github.com/ai-agentic-browser/internal/config"
	"github.com/ai-agentic-browser/internal/monitoring"
	"github.com/ai-agentic-browser/internal/realtime"
	"github.com/ai-agentic-browser/internal/web3"
	"github.com/ai-agentic-browser/pkg/database"
	"github.com/ai-agentic-browser/pkg/middleware"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize observability
	logger := observability.NewLogger(cfg.Observability)
	tracingProvider, err := observability.NewTracingProvider(cfg.Observability)
	if err != nil {
		log.Fatalf("Failed to initialize tracing: %v", err)
	}
	defer tracingProvider.Shutdown(context.Background())

	// Initialize database connections
	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	redis, err := database.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redis.Close()

	// Initialize Web3 service
	web3Service := web3.NewService(db, redis, cfg.Web3, logger)

	// Initialize enhanced Web3 service with autonomous trading
	enhancedService, err := web3.NewEnhancedService(db, redis, cfg.Web3, logger)
	if err != nil {
		log.Fatalf("Failed to initialize enhanced Web3 service: %v", err)
	}

	// Initialize autonomous trading components
	riskAssessment := web3.NewRiskAssessmentService(enhancedService.GetClients(), logger)
	tradingEngine := web3.NewTradingEngine(enhancedService.GetClients(), logger, riskAssessment)
	defiManager := web3.NewDeFiProtocolManager(logger)
	portfolioRebalancer := web3.NewPortfolioRebalancer(logger, tradingEngine, defiManager)

	// Initialize AI components
	voiceInterface := ai.NewVoiceInterface(logger, tradingEngine, defiManager, riskAssessment)
	conversationalAI := ai.NewConversationalAI(logger, tradingEngine, defiManager, riskAssessment)

	// Initialize real-time monitoring components
	marketDataConfig := realtime.MarketDataConfig{
		Exchanges: []realtime.ExchangeConfig{
			{
				Name:     "binance",
				WSUrl:    "wss://stream.binance.com:9443/ws",
				Symbols:  []string{"BTCUSDT", "ETHUSDT", "ADAUSDT"},
				Channels: []string{"ticker", "trade"},
				Enabled:  true,
			},
		},
		ReconnectDelay:  5 * time.Second,
		PingInterval:    30 * time.Second,
		MaxReconnects:   10,
		BufferSize:      1000,
		EnableHeartbeat: true,
	}
	marketDataService := realtime.NewMarketDataService(logger, marketDataConfig)

	// Initialize portfolio analytics
	portfolioAnalytics := analytics.NewPortfolioAnalytics(logger, tradingEngine)

	// Initialize system monitoring
	monitoringConfig := monitoring.MonitoringConfig{
		CollectionInterval: 30 * time.Second,
		RetentionPeriod:    24 * time.Hour,
		AlertThresholds: monitoring.AlertConfig{
			CPUThreshold:        80.0,
			MemoryThreshold:     85.0,
			DiskThreshold:       90.0,
			ErrorRateThreshold:  5.0,
			LatencyThreshold:    1000.0,
			ConnectionThreshold: 1000,
		},
		EnableProfiling: true,
		EnableTracing:   true,
	}
	systemMonitor := monitoring.NewSystemMonitor(logger, monitoringConfig)

	// Initialize alert service
	alertConfig := alerts.AlertConfig{
		MaxHistorySize:  1000,
		DefaultCooldown: 5 * time.Minute,
		EnableEmail:     true,
		EnableWebhook:   true,
		EnableSlack:     true,
		EnableTelegram:  false,
		EnablePushNotif: true,
	}
	alertService := alerts.NewAlertService(logger, alertConfig)

	// Start all services
	go func() {
		if err := tradingEngine.Start(context.Background()); err != nil {
			logger.Error(context.Background(), "Failed to start trading engine", err)
		}
	}()

	go func() {
		if err := marketDataService.Start(); err != nil {
			logger.Error(context.Background(), "Failed to start market data service", err)
		}
	}()

	go func() {
		if err := systemMonitor.Start(); err != nil {
			logger.Error(context.Background(), "Failed to start system monitor", err)
		}
	}()

	go func() {
		if err := alertService.Start(); err != nil {
			logger.Error(context.Background(), "Failed to start alert service", err)
		}
	}()

	// Store components for use in handlers
	_ = portfolioRebalancer // Will be used in handlers

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, "8084"), // Web3 service port
		Handler:      setupRoutes(web3Service, enhancedService, tradingEngine, defiManager, portfolioRebalancer, voiceInterface, conversationalAI, marketDataService, portfolioAnalytics, systemMonitor, alertService, cfg, logger, db),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		logger.Info(context.Background(), "Starting Web3 service", map[string]interface{}{
			"addr": server.Addr,
		})
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info(context.Background(), "Shutting down Web3 service...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info(context.Background(), "Web3 service stopped")
}

func setupRoutes(
	web3Service *web3.Service,
	enhancedService *web3.EnhancedService,
	tradingEngine *web3.TradingEngine,
	defiManager *web3.DeFiProtocolManager,
	portfolioRebalancer *web3.PortfolioRebalancer,
	voiceInterface *ai.VoiceInterface,
	conversationalAI *ai.ConversationalAI,
	marketDataService *realtime.MarketDataService,
	portfolioAnalytics *analytics.PortfolioAnalytics,
	systemMonitor *monitoring.SystemMonitor,
	alertService *alerts.AlertService,
	cfg *config.Config,
	logger *observability.Logger,
	db *database.DB,
) http.Handler {
	mux := http.NewServeMux()

	// Apply middleware
	handler := middleware.Recovery(logger)(
		middleware.Logging(logger)(
			middleware.Tracing("web3-service")(
				middleware.CORS(cfg.Security.CORSAllowedOrigins)(
					middleware.RateLimit(cfg.RateLimit)(mux),
				),
			),
		),
	)

	// Health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Check database health
		if err := db.Health(ctx); err != nil {
			http.Error(w, "Database unhealthy", http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	// Protected Web3 endpoints
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("POST /web3/connect-wallet", handleConnectWallet(web3Service, logger))
	protectedMux.HandleFunc("GET /web3/wallets", handleListWallets(web3Service, logger))
	protectedMux.HandleFunc("GET /web3/balance", handleGetBalance(web3Service, logger))
	protectedMux.HandleFunc("POST /web3/transaction", handleCreateTransaction(web3Service, logger))
	protectedMux.HandleFunc("GET /web3/transactions", handleListTransactions(web3Service, logger))
	protectedMux.HandleFunc("GET /web3/prices", handleGetPrices(web3Service, logger))
	protectedMux.HandleFunc("POST /web3/defi/interact", handleDeFiInteraction(web3Service, logger))
	protectedMux.HandleFunc("GET /web3/defi/positions", handleListDeFiPositions(web3Service, logger))
	protectedMux.HandleFunc("GET /web3/chains", handleGetSupportedChains(web3Service, logger))

	// Enhanced Web3 endpoints
	protectedMux.HandleFunc("POST /web3/enhanced/transaction", handleEnhancedTransaction(enhancedService, logger))

	// Autonomous Trading endpoints
	protectedMux.HandleFunc("POST /web3/trading/portfolio", handleCreatePortfolio(tradingEngine, logger))
	protectedMux.HandleFunc("GET /web3/trading/portfolio/{id}", handleGetPortfolio(tradingEngine, logger))
	protectedMux.HandleFunc("POST /web3/trading/portfolio/{id}/start", handleStartTrading(tradingEngine, logger))
	protectedMux.HandleFunc("POST /web3/trading/portfolio/{id}/stop", handleStopTrading(tradingEngine, logger))
	protectedMux.HandleFunc("GET /web3/trading/positions/{portfolio_id}", handleGetPositions(tradingEngine, logger))
	protectedMux.HandleFunc("POST /web3/trading/positions/{id}/close", handleClosePosition(tradingEngine, logger))

	// DeFi Protocol endpoints
	protectedMux.HandleFunc("GET /web3/defi/protocols", handleGetProtocols(defiManager, logger))
	protectedMux.HandleFunc("GET /web3/defi/protocols/{id}", handleGetProtocol(defiManager, logger))
	protectedMux.HandleFunc("GET /web3/defi/opportunities", handleGetYieldOpportunities(defiManager, logger))

	// Portfolio Rebalancing endpoints
	protectedMux.HandleFunc("POST /web3/rebalance/strategy", handleCreateRebalanceStrategy(portfolioRebalancer, logger))
	protectedMux.HandleFunc("GET /web3/rebalance/strategy/{portfolio_id}", handleGetRebalanceStrategy(portfolioRebalancer, logger))
	protectedMux.HandleFunc("POST /web3/rebalance/execute/{portfolio_id}", handleExecuteRebalancing(portfolioRebalancer, logger))

	// AI Voice Interface endpoints
	protectedMux.HandleFunc("POST /web3/ai/voice/command", handleVoiceCommand(voiceInterface, logger))
	protectedMux.HandleFunc("GET /web3/ai/voice/history", handleVoiceHistory(voiceInterface, logger))

	// Conversational AI endpoints
	protectedMux.HandleFunc("POST /web3/ai/chat/message", handleChatMessage(conversationalAI, logger))
	protectedMux.HandleFunc("POST /web3/ai/chat/start", handleStartConversation(conversationalAI, logger))
	protectedMux.HandleFunc("GET /web3/ai/market/analysis", handleMarketAnalysis(conversationalAI, logger))

	// Real-time Market Data endpoints
	protectedMux.HandleFunc("GET /web3/realtime/market/status", handleMarketDataStatus(marketDataService, logger))
	protectedMux.HandleFunc("GET /web3/realtime/market/subscribe/{symbol}", handleMarketDataSubscribe(marketDataService, logger))

	// Portfolio Analytics endpoints
	protectedMux.HandleFunc("GET /web3/analytics/portfolio/{portfolio_id}", handlePortfolioAnalytics(portfolioAnalytics, logger))
	protectedMux.HandleFunc("GET /web3/analytics/portfolio/{portfolio_id}/performance", handlePortfolioPerformance(portfolioAnalytics, logger))
	protectedMux.HandleFunc("GET /web3/analytics/portfolio/compare", handlePortfolioComparison(portfolioAnalytics, logger))

	// System Monitoring endpoints
	protectedMux.HandleFunc("GET /web3/monitoring/health", handleSystemHealth(systemMonitor, logger))
	protectedMux.HandleFunc("GET /web3/monitoring/metrics", handleSystemMetrics(systemMonitor, logger))
	protectedMux.HandleFunc("GET /web3/monitoring/status", handleSystemStatus(systemMonitor, logger))

	// Alert Management endpoints
	protectedMux.HandleFunc("GET /web3/alerts", handleGetAlerts(alertService, logger))
	protectedMux.HandleFunc("GET /web3/alerts/active", handleGetActiveAlerts(alertService, logger))
	protectedMux.HandleFunc("POST /web3/alerts/{alert_id}/resolve", handleResolveAlert(alertService, logger))
	protectedMux.HandleFunc("GET /web3/alerts/subscribe/{topic}", handleAlertSubscribe(alertService, logger))

	// Apply JWT middleware to protected routes
	mux.Handle("/web3/", middleware.JWT(cfg.JWT.Secret)(protectedMux))

	return handler
}

func handleConnectWallet(web3Service *web3.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr, ok := middleware.GetUserID(r.Context())
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusInternalServerError)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		var req web3.WalletConnectRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		response, err := web3Service.ConnectWallet(r.Context(), userID, req)
		if err != nil {
			logger.Error(r.Context(), "Wallet connection failed", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

func handleListWallets(web3Service *web3.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr, ok := middleware.GetUserID(r.Context())
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusInternalServerError)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// TODO: Implement ListWallets method in Web3 service
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user_id": userID.String(),
			"wallets": []interface{}{},
			"message": "Wallet listing not implemented yet",
		})
	}
}

func handleGetBalance(web3Service *web3.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr, ok := middleware.GetUserID(r.Context())
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusInternalServerError)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		var req web3.BalanceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		response, err := web3Service.GetBalance(r.Context(), userID, req)
		if err != nil {
			logger.Error(r.Context(), "Balance retrieval failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func handleCreateTransaction(web3Service *web3.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr, ok := middleware.GetUserID(r.Context())
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusInternalServerError)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		var req web3.TransactionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		response, err := web3Service.CreateTransaction(r.Context(), userID, req)
		if err != nil {
			logger.Error(r.Context(), "Transaction creation failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

func handleListTransactions(web3Service *web3.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr, ok := middleware.GetUserID(r.Context())
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusInternalServerError)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// TODO: Implement ListTransactions method in Web3 service
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user_id":      userID.String(),
			"transactions": []interface{}{},
			"message":      "Transaction listing not implemented yet",
		})
	}
}

func handleGetPrices(web3Service *web3.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req web3.PriceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			// If no body, use default request
			req = web3.PriceRequest{Currency: "USD"}
		}

		response, err := web3Service.GetPrices(r.Context(), req)
		if err != nil {
			logger.Error(r.Context(), "Price retrieval failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func handleDeFiInteraction(web3Service *web3.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr, ok := middleware.GetUserID(r.Context())
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusInternalServerError)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		var req web3.DeFiProtocolRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		response, err := web3Service.InteractWithDeFiProtocol(r.Context(), userID, req)
		if err != nil {
			logger.Error(r.Context(), "DeFi interaction failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func handleListDeFiPositions(web3Service *web3.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr, ok := middleware.GetUserID(r.Context())
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusInternalServerError)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// TODO: Implement ListDeFiPositions method in Web3 service
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user_id":   userID.String(),
			"positions": []interface{}{},
			"message":   "DeFi position listing not implemented yet",
		})
	}
}

func handleGetSupportedChains(web3Service *web3.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"chains": web3.SupportedChains,
		})
	}
}

// Enhanced Web3 handlers
func handleEnhancedTransaction(enhancedService *web3.EnhancedService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr, ok := middleware.GetUserID(r.Context())
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusInternalServerError)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		var req web3.EnhancedTransactionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		response, err := enhancedService.CreateEnhancedTransaction(r.Context(), userID, req)
		if err != nil {
			logger.Error(r.Context(), "Enhanced transaction creation failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

// Trading Engine handlers
func handleCreatePortfolio(tradingEngine *web3.TradingEngine, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr, ok := middleware.GetUserID(r.Context())
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusInternalServerError)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		var req struct {
			Name           string           `json:"name"`
			InitialBalance string           `json:"initial_balance"`
			RiskProfile    web3.RiskProfile `json:"risk_profile"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Parse initial balance
		initialBalance, err := decimal.NewFromString(req.InitialBalance)
		if err != nil {
			http.Error(w, "Invalid initial balance", http.StatusBadRequest)
			return
		}

		portfolio, err := tradingEngine.CreatePortfolio(r.Context(), userID, req.Name, initialBalance, req.RiskProfile)
		if err != nil {
			logger.Error(r.Context(), "Portfolio creation failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(portfolio)
	}
}

func handleGetPortfolio(tradingEngine *web3.TradingEngine, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		portfolioIDStr := strings.TrimPrefix(r.URL.Path, "/web3/trading/portfolio/")
		portfolioID, err := uuid.Parse(portfolioIDStr)
		if err != nil {
			http.Error(w, "Invalid portfolio ID", http.StatusBadRequest)
			return
		}

		portfolio, err := tradingEngine.GetPortfolio(portfolioID)
		if err != nil {
			logger.Error(r.Context(), "Portfolio retrieval failed", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(portfolio)
	}
}

func handleStartTrading(tradingEngine *web3.TradingEngine, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Trading engine is already running globally",
			"status":  "success",
		})
	}
}

func handleStopTrading(tradingEngine *web3.TradingEngine, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := tradingEngine.Stop(r.Context())
		if err != nil {
			logger.Error(r.Context(), "Failed to stop trading", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Trading stopped successfully",
			"status":  "success",
		})
	}
}

func handleGetPositions(tradingEngine *web3.TradingEngine, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		portfolioIDStr := strings.TrimPrefix(r.URL.Path, "/web3/trading/positions/")
		portfolioID, err := uuid.Parse(portfolioIDStr)
		if err != nil {
			http.Error(w, "Invalid portfolio ID", http.StatusBadRequest)
			return
		}

		positions, err := tradingEngine.GetActivePositions(portfolioID)
		if err != nil {
			logger.Error(r.Context(), "Positions retrieval failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"portfolio_id": portfolioID.String(),
			"positions":    positions,
		})
	}
}

func handleClosePosition(tradingEngine *web3.TradingEngine, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		positionIDStr := strings.TrimPrefix(r.URL.Path, "/web3/trading/positions/")
		positionIDStr = strings.TrimSuffix(positionIDStr, "/close")
		positionID, err := uuid.Parse(positionIDStr)
		if err != nil {
			http.Error(w, "Invalid position ID", http.StatusBadRequest)
			return
		}

		var req struct {
			Reason string `json:"reason"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			req.Reason = "Manual close"
		}

		err = tradingEngine.ClosePosition(r.Context(), positionID, req.Reason)
		if err != nil {
			logger.Error(r.Context(), "Position close failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":     "Position closed successfully",
			"position_id": positionID.String(),
			"reason":      req.Reason,
		})
	}
}

// DeFi Protocol handlers
func handleGetProtocols(defiManager *web3.DeFiProtocolManager, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		protocols := defiManager.GetProtocols()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"protocols": protocols,
		})
	}
}

func handleGetProtocol(defiManager *web3.DeFiProtocolManager, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		protocolID := strings.TrimPrefix(r.URL.Path, "/web3/defi/protocols/")

		protocol, err := defiManager.GetProtocol(protocolID)
		if err != nil {
			logger.Error(r.Context(), "Protocol retrieval failed", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(protocol)
	}
}

func handleGetYieldOpportunities(defiManager *web3.DeFiProtocolManager, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		minAPYStr := r.URL.Query().Get("min_apy")
		maxRiskStr := r.URL.Query().Get("max_risk")

		minAPY := decimal.NewFromFloat(0.01) // Default 1%
		if minAPYStr != "" {
			if parsed, err := decimal.NewFromString(minAPYStr); err == nil {
				minAPY = parsed
			}
		}

		maxRisk := web3.RiskLevelMedium // Default medium risk
		if maxRiskStr != "" {
			switch maxRiskStr {
			case "very_low":
				maxRisk = web3.RiskLevelVeryLow
			case "low":
				maxRisk = web3.RiskLevelLow
			case "medium":
				maxRisk = web3.RiskLevelMedium
			case "high":
				maxRisk = web3.RiskLevelHigh
			case "critical":
				maxRisk = web3.RiskLevelCritical
			}
		}

		opportunities, err := defiManager.GetBestYieldOpportunities(r.Context(), minAPY, maxRisk)
		if err != nil {
			logger.Error(r.Context(), "Yield opportunities retrieval failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"opportunities": opportunities,
			"filters": map[string]interface{}{
				"min_apy":  minAPY.String(),
				"max_risk": string(maxRisk),
			},
		})
	}
}

// Portfolio Rebalancing handlers
func handleCreateRebalanceStrategy(portfolioRebalancer *web3.PortfolioRebalancer, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			PortfolioID       string                     `json:"portfolio_id"`
			Name              string                     `json:"name"`
			Type              web3.RebalanceType         `json:"type"`
			TargetAllocations map[string]decimal.Decimal `json:"target_allocations"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		portfolioID, err := uuid.Parse(req.PortfolioID)
		if err != nil {
			http.Error(w, "Invalid portfolio ID", http.StatusBadRequest)
			return
		}

		strategy, err := portfolioRebalancer.CreateRebalanceStrategy(
			r.Context(), portfolioID, req.Name, req.Type, req.TargetAllocations)
		if err != nil {
			logger.Error(r.Context(), "Rebalance strategy creation failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(strategy)
	}
}

func handleGetRebalanceStrategy(portfolioRebalancer *web3.PortfolioRebalancer, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		portfolioIDStr := strings.TrimPrefix(r.URL.Path, "/web3/rebalance/strategy/")
		portfolioID, err := uuid.Parse(portfolioIDStr)
		if err != nil {
			http.Error(w, "Invalid portfolio ID", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"portfolio_id": portfolioID.String(),
			"message":      "Rebalance strategy retrieval not implemented yet",
		})
	}
}

func handleExecuteRebalancing(portfolioRebalancer *web3.PortfolioRebalancer, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		portfolioIDStr := strings.TrimPrefix(r.URL.Path, "/web3/rebalance/execute/")
		portfolioID, err := uuid.Parse(portfolioIDStr)
		if err != nil {
			http.Error(w, "Invalid portfolio ID", http.StatusBadRequest)
			return
		}

		err = portfolioRebalancer.RebalancePortfolio(r.Context(), portfolioID)
		if err != nil {
			logger.Error(r.Context(), "Portfolio rebalancing failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":      "Portfolio rebalanced successfully",
			"portfolio_id": portfolioID.String(),
		})
	}
}

// AI Voice Interface handlers
func handleVoiceCommand(voiceInterface *ai.VoiceInterface, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr, ok := middleware.GetUserID(r.Context())
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusInternalServerError)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		var req struct {
			Text      string `json:"text"`
			AudioData []byte `json:"audio_data,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		response, err := voiceInterface.ProcessVoiceCommand(r.Context(), userID, req.AudioData, req.Text)
		if err != nil {
			logger.Error(r.Context(), "Voice command processing failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func handleVoiceHistory(voiceInterface *ai.VoiceInterface, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr, ok := middleware.GetUserID(r.Context())
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusInternalServerError)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		history := voiceInterface.GetCommandHistory(userID)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user_id": userID.String(),
			"history": history,
		})
	}
}

// Conversational AI handlers
func handleChatMessage(conversationalAI *ai.ConversationalAI, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr, ok := middleware.GetUserID(r.Context())
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusInternalServerError)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		var req struct {
			Message string `json:"message"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		response, err := conversationalAI.ProcessMessage(r.Context(), userID, req.Message)
		if err != nil {
			logger.Error(r.Context(), "Chat message processing failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func handleStartConversation(conversationalAI *ai.ConversationalAI, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr, ok := middleware.GetUserID(r.Context())
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusInternalServerError)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		conversation, err := conversationalAI.StartConversation(r.Context(), userID)
		if err != nil {
			logger.Error(r.Context(), "Conversation start failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(conversation)
	}
}

func handleMarketAnalysis(conversationalAI *ai.ConversationalAI, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr, ok := middleware.GetUserID(r.Context())
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusInternalServerError)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Generate market analysis using conversational AI
		response, err := conversationalAI.ProcessMessage(r.Context(), userID, "Give me a comprehensive market analysis")
		if err != nil {
			logger.Error(r.Context(), "Market analysis failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// Real-time Market Data handlers
func handleMarketDataStatus(marketDataService *realtime.MarketDataService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := marketDataService.GetConnectionStatus()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":      "active",
			"connections": status,
			"timestamp":   time.Now(),
		})
	}
}

func handleMarketDataSubscribe(marketDataService *realtime.MarketDataService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		symbol := strings.TrimPrefix(r.URL.Path, "/web3/realtime/market/subscribe/")

		// Subscribe to market data updates
		updateChan := marketDataService.Subscribe(symbol)

		// Set up Server-Sent Events
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Send initial connection message
		fmt.Fprintf(w, "data: {\"type\":\"connected\",\"symbol\":\"%s\"}\n\n", symbol)
		w.(http.Flusher).Flush()

		// Stream market data updates
		for {
			select {
			case update := <-updateChan:
				data, _ := json.Marshal(update)
				fmt.Fprintf(w, "data: %s\n\n", data)
				w.(http.Flusher).Flush()
			case <-r.Context().Done():
				marketDataService.Unsubscribe(symbol, updateChan)
				return
			}
		}
	}
}

// Portfolio Analytics handlers
func handlePortfolioAnalytics(portfolioAnalytics *analytics.PortfolioAnalytics, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		portfolioIDStr := strings.TrimPrefix(r.URL.Path, "/web3/analytics/portfolio/")
		portfolioID, err := uuid.Parse(portfolioIDStr)
		if err != nil {
			http.Error(w, "Invalid portfolio ID", http.StatusBadRequest)
			return
		}

		metrics, err := portfolioAnalytics.GetPortfolioMetrics(r.Context(), portfolioID)
		if err != nil {
			logger.Error(r.Context(), "Portfolio analytics retrieval failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(metrics)
	}
}

func handlePortfolioPerformance(portfolioAnalytics *analytics.PortfolioAnalytics, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		portfolioIDStr := strings.TrimPrefix(r.URL.Path, "/web3/analytics/portfolio/")
		portfolioIDStr = strings.TrimSuffix(portfolioIDStr, "/performance")
		portfolioID, err := uuid.Parse(portfolioIDStr)
		if err != nil {
			http.Error(w, "Invalid portfolio ID", http.StatusBadRequest)
			return
		}

		metrics, err := portfolioAnalytics.GetPortfolioMetrics(r.Context(), portfolioID)
		if err != nil {
			logger.Error(r.Context(), "Portfolio performance retrieval failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"portfolio_id":  portfolioID.String(),
			"performance":   metrics.Performance,
			"risk_metrics":  metrics.RiskMetrics,
			"sharpe_ratio":  metrics.SharpeRatio,
			"sortino_ratio": metrics.SortinoRatio,
			"max_drawdown":  metrics.MaxDrawdown,
			"volatility":    metrics.Volatility,
		})
	}
}

func handlePortfolioComparison(portfolioAnalytics *analytics.PortfolioAnalytics, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		portfolioIDsStr := r.URL.Query().Get("portfolio_ids")
		if portfolioIDsStr == "" {
			http.Error(w, "portfolio_ids parameter required", http.StatusBadRequest)
			return
		}

		portfolioIDStrs := strings.Split(portfolioIDsStr, ",")
		portfolioIDs := make([]uuid.UUID, 0, len(portfolioIDStrs))

		for _, idStr := range portfolioIDStrs {
			id, err := uuid.Parse(strings.TrimSpace(idStr))
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid portfolio ID: %s", idStr), http.StatusBadRequest)
				return
			}
			portfolioIDs = append(portfolioIDs, id)
		}

		comparison, err := portfolioAnalytics.GetPortfolioComparison(r.Context(), portfolioIDs)
		if err != nil {
			logger.Error(r.Context(), "Portfolio comparison failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comparison)
	}
}

// System Monitoring handlers
func handleSystemHealth(systemMonitor *monitoring.SystemMonitor, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics := systemMonitor.GetCurrentMetrics()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":     metrics.Health.Status,
			"score":      metrics.Health.Score,
			"components": metrics.Health.Components,
			"issues":     metrics.Health.Issues,
			"last_check": metrics.Health.LastCheck,
		})
	}
}

func handleSystemMetrics(systemMonitor *monitoring.SystemMonitor, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics := systemMonitor.GetCurrentMetrics()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(metrics)
	}
}

func handleSystemStatus(systemMonitor *monitoring.SystemMonitor, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics := systemMonitor.GetCurrentMetrics()
		alerts := systemMonitor.GetAlerts()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"health":    metrics.Health,
			"alerts":    alerts,
			"timestamp": time.Now(),
		})
	}
}

// Alert Management handlers
func handleGetAlerts(alertService *alerts.AlertService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limitStr := r.URL.Query().Get("limit")
		limit := 50 // Default limit
		if limitStr != "" {
			if parsed, err := fmt.Sscanf(limitStr, "%d", &limit); err != nil || parsed != 1 {
				limit = 50
			}
		}

		alertList := alertService.GetAlerts(limit)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"alerts": alertList,
			"count":  len(alertList),
			"limit":  limit,
		})
	}
}

func handleGetActiveAlerts(alertService *alerts.AlertService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		activeAlerts := alertService.GetActiveAlerts()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"alerts": activeAlerts,
			"count":  len(activeAlerts),
		})
	}
}

func handleResolveAlert(alertService *alerts.AlertService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alertID := strings.TrimPrefix(r.URL.Path, "/web3/alerts/")
		alertID = strings.TrimSuffix(alertID, "/resolve")

		err := alertService.ResolveAlert(alertID)
		if err != nil {
			logger.Error(r.Context(), "Alert resolution failed", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":  "Alert resolved successfully",
			"alert_id": alertID,
		})
	}
}

func handleAlertSubscribe(alertService *alerts.AlertService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		topic := strings.TrimPrefix(r.URL.Path, "/web3/alerts/subscribe/")

		// Subscribe to alerts
		alertChan := alertService.Subscribe(topic)

		// Set up Server-Sent Events
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Send initial connection message
		fmt.Fprintf(w, "data: {\"type\":\"connected\",\"topic\":\"%s\"}\n\n", topic)
		w.(http.Flusher).Flush()

		// Stream alert updates
		for {
			select {
			case alert := <-alertChan:
				data, _ := json.Marshal(alert)
				fmt.Fprintf(w, "data: %s\n\n", data)
				w.(http.Flusher).Flush()
			case <-r.Context().Done():
				return
			}
		}
	}
}
