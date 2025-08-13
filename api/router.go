package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ai-agentic-browser/internal/binance"
	"github.com/ai-agentic-browser/internal/exchanges"
	"github.com/ai-agentic-browser/internal/hft"
	"github.com/ai-agentic-browser/internal/mcp"
	"github.com/ai-agentic-browser/internal/tradingview"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/ai-agentic-browser/pkg/strategies"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

// APIServer provides HTTP API endpoints for the HFT system
type APIServer struct {
	logger *observability.Logger
	config Config
	router *mux.Router
	server *http.Server

	// Core services
	hftEngine          *hft.HFTEngine
	binanceService     *binance.Service
	tradingViewService *tradingview.Service
	mcpService         *mcp.IntegrationService
	strategyEngine     *strategies.StrategyEngine

	// Exchange infrastructure
	exchangeManager   *exchanges.Manager
	orderManager      *exchanges.OrderManager
	marketDataService *exchanges.MarketDataService

	// WebSocket upgrader
	upgrader    websocket.Upgrader
	wsClients   map[*websocket.Conn]bool
	wsBroadcast chan []byte

	// State
	isRunning bool
}

// Config contains API server configuration
type Config struct {
	Host            string        `json:"host"`
	Port            int           `json:"port"`
	ReadTimeout     time.Duration `json:"read_timeout"`
	WriteTimeout    time.Duration `json:"write_timeout"`
	EnableCORS      bool          `json:"enable_cors"`
	EnableWebSocket bool          `json:"enable_websocket"`
	RateLimit       int           `json:"rate_limit"`
}

// Response represents a standard API response
type Response struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

// NewAPIServer creates a new API server
func NewAPIServer(logger *observability.Logger, config Config) *APIServer {
	if config.Host == "" {
		config.Host = "localhost"
	}
	if config.Port == 0 {
		config.Port = 8080
	}
	if config.ReadTimeout == 0 {
		config.ReadTimeout = 30 * time.Second
	}
	if config.WriteTimeout == 0 {
		config.WriteTimeout = 30 * time.Second
	}

	server := &APIServer{
		logger: logger,
		config: config,
		router: mux.NewRouter(),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins in development
			},
		},
		wsClients:   make(map[*websocket.Conn]bool),
		wsBroadcast: make(chan []byte),
	}

	server.setupRoutes()

	return server
}

// SetServices sets the core services
func (s *APIServer) SetServices(
	hftEngine *hft.HFTEngine,
	binanceService *binance.Service,
	tradingViewService *tradingview.Service,
	mcpService *mcp.IntegrationService,
	strategyEngine *strategies.StrategyEngine,
) {
	s.hftEngine = hftEngine
	s.binanceService = binanceService
	s.tradingViewService = tradingViewService
	s.mcpService = mcpService
	s.strategyEngine = strategyEngine
}

// SetExchangeServices sets the exchange infrastructure services
func (s *APIServer) SetExchangeServices(
	exchangeManager *exchanges.Manager,
	orderManager *exchanges.OrderManager,
	marketDataService *exchanges.MarketDataService,
) {
	s.exchangeManager = exchangeManager
	s.orderManager = orderManager
	s.marketDataService = marketDataService
}

// Start begins the API server
func (s *APIServer) Start(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	// Setup CORS if enabled
	var handler http.Handler = s.router
	if s.config.EnableCORS {
		c := cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"*"},
		})
		handler = c.Handler(s.router)
	}

	s.server = &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
	}

	s.logger.Info(ctx, "Starting API server", map[string]interface{}{
		"address":          addr,
		"enable_cors":      s.config.EnableCORS,
		"enable_websocket": s.config.EnableWebSocket,
	})

	// Start WebSocket hub if enabled
	if s.config.EnableWebSocket {
		go s.handleWebSocketHub()
	}

	s.isRunning = true

	// Start server in goroutine
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error(ctx, "API server error", err)
		}
	}()

	s.logger.Info(ctx, "API server started successfully", map[string]interface{}{
		"address": addr,
	})

	return nil
}

// Stop gracefully shuts down the API server
func (s *APIServer) Stop(ctx context.Context) error {
	if !s.isRunning {
		return fmt.Errorf("API server is not running")
	}

	s.logger.Info(ctx, "Stopping API server", nil)

	// Close WebSocket connections
	for client := range s.wsClients {
		client.Close()
	}

	// Shutdown server
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown API server: %w", err)
	}

	s.isRunning = false

	s.logger.Info(ctx, "API server stopped successfully", nil)

	return nil
}

// setupRoutes configures all API routes
func (s *APIServer) setupRoutes() {
	// Health check
	s.router.HandleFunc("/health", s.handleHealth).Methods("GET")

	// HFT Engine endpoints
	hftRouter := s.router.PathPrefix("/api/hft").Subrouter()
	hftRouter.HandleFunc("/start", s.handleHFTStart).Methods("POST")
	hftRouter.HandleFunc("/stop", s.handleHFTStop).Methods("POST")
	hftRouter.HandleFunc("/status", s.handleHFTStatus).Methods("GET")
	hftRouter.HandleFunc("/metrics", s.handleHFTMetrics).Methods("GET")
	hftRouter.HandleFunc("/config", s.handleHFTConfig).Methods("GET", "PUT")

	// Market data endpoints
	marketRouter := s.router.PathPrefix("/api/market").Subrouter()
	marketRouter.HandleFunc("/tickers", s.handleMarketTickers).Methods("GET")
	marketRouter.HandleFunc("/ticker/{symbol}", s.handleMarketTicker).Methods("GET")
	marketRouter.HandleFunc("/orderbook/{symbol}", s.handleMarketOrderbook).Methods("GET")
	marketRouter.HandleFunc("/trades/{symbol}", s.handleMarketTrades).Methods("GET")

	// Trading endpoints
	tradingRouter := s.router.PathPrefix("/api/trading").Subrouter()
	tradingRouter.HandleFunc("/orders", s.handleTradingOrders).Methods("GET", "POST")
	tradingRouter.HandleFunc("/orders/{id}", s.handleTradingOrder).Methods("GET", "DELETE")
	tradingRouter.HandleFunc("/positions", s.handleTradingPositions).Methods("GET")
	tradingRouter.HandleFunc("/signals", s.handleTradingSignals).Methods("GET")

	// Portfolio endpoints
	portfolioRouter := s.router.PathPrefix("/api/portfolio").Subrouter()
	portfolioRouter.HandleFunc("/summary", s.handlePortfolioSummary).Methods("GET")
	portfolioRouter.HandleFunc("/positions", s.handlePortfolioPositions).Methods("GET")
	portfolioRouter.HandleFunc("/metrics", s.handlePortfolioMetrics).Methods("GET")
	portfolioRouter.HandleFunc("/risk", s.handlePortfolioRisk).Methods("GET")

	// Strategy endpoints
	strategyRouter := s.router.PathPrefix("/api/strategies").Subrouter()
	strategyRouter.HandleFunc("", s.handleStrategies).Methods("GET", "POST")
	strategyRouter.HandleFunc("/{id}", s.handleStrategy).Methods("GET", "PUT", "DELETE")
	strategyRouter.HandleFunc("/{id}/start", s.handleStrategyStart).Methods("POST")
	strategyRouter.HandleFunc("/{id}/stop", s.handleStrategyStop).Methods("POST")
	strategyRouter.HandleFunc("/{id}/performance", s.handleStrategyPerformance).Methods("GET")

	// Risk management endpoints
	riskRouter := s.router.PathPrefix("/api/risk").Subrouter()
	riskRouter.HandleFunc("/limits", s.handleRiskLimits).Methods("GET", "POST")
	riskRouter.HandleFunc("/limits/{id}", s.handleRiskLimit).Methods("GET", "PUT", "DELETE")
	riskRouter.HandleFunc("/violations", s.handleRiskViolations).Methods("GET")
	riskRouter.HandleFunc("/metrics", s.handleRiskMetrics).Methods("GET")
	riskRouter.HandleFunc("/emergency-stop", s.handleEmergencyStop).Methods("POST")

	// TradingView endpoints
	tvRouter := s.router.PathPrefix("/api/tradingview").Subrouter()
	tvRouter.HandleFunc("/charts", s.handleTradingViewCharts).Methods("GET", "POST")
	tvRouter.HandleFunc("/charts/{id}", s.handleTradingViewChart).Methods("GET", "DELETE")
	tvRouter.HandleFunc("/signals", s.handleTradingViewSignals).Methods("GET")
	tvRouter.HandleFunc("/indicators/{symbol}", s.handleTradingViewIndicators).Methods("GET")

	// MCP integration endpoints
	mcpRouter := s.router.PathPrefix("/api/mcp").Subrouter()
	mcpRouter.HandleFunc("/insights", s.handleMCPInsights).Methods("GET")
	mcpRouter.HandleFunc("/insights/{symbol}", s.handleMCPInsight).Methods("GET")
	mcpRouter.HandleFunc("/sentiment/{symbol}", s.handleMCPSentiment).Methods("GET")
	mcpRouter.HandleFunc("/news/{symbol}", s.handleMCPNews).Methods("GET")

	// Firebase MCP endpoints
	firebaseRouter := s.router.PathPrefix("/api/firebase").Subrouter()

	// Firebase Status
	firebaseRouter.HandleFunc("/status", s.handleFirebaseStatus).Methods("GET")

	// Firebase Authentication
	firebaseRouter.HandleFunc("/auth/users", s.handleFirebaseAuth).Methods("POST")
	firebaseRouter.HandleFunc("/auth/users/{uid}", s.handleFirebaseAuth).Methods("GET")
	firebaseRouter.HandleFunc("/auth/verify", s.handleFirebaseVerifyToken).Methods("POST")

	// Firestore Database
	firebaseRouter.HandleFunc("/firestore/{collection}", s.handleFirestoreDocument).Methods("GET", "POST")
	firebaseRouter.HandleFunc("/firestore/{collection}/{documentId}", s.handleFirestoreDocument).Methods("GET", "PUT", "DELETE")
	firebaseRouter.HandleFunc("/firestore/batch", s.handleFirebaseBatchWrite).Methods("POST")

	// Realtime Database
	firebaseRouter.HandleFunc("/realtime/{path:.*}", s.handleFirebaseRealtimeDB).Methods("GET", "POST", "PUT", "PATCH", "DELETE")

	// Exchange endpoints
	exchangeRouter := s.router.PathPrefix("/api/exchanges").Subrouter()
	exchangeRouter.HandleFunc("", s.handleGetExchanges).Methods("GET")
	exchangeRouter.HandleFunc("/ticker/{symbol}", s.handleGetTicker).Methods("GET")
	exchangeRouter.HandleFunc("/orderbook/{symbol}", s.handleGetOrderBook).Methods("GET")
	exchangeRouter.HandleFunc("/best-price/{symbol}", s.handleGetBestPrice).Methods("GET")
	exchangeRouter.HandleFunc("/orders", s.handleSubmitOrder).Methods("POST")
	exchangeRouter.HandleFunc("/orders/{id}", s.handleGetOrder).Methods("GET")
	exchangeRouter.HandleFunc("/orders/{id}/cancel", s.handleCancelOrder).Methods("POST")

	// WebSocket endpoint
	if s.config.EnableWebSocket {
		s.router.HandleFunc("/ws/trading", s.handleWebSocket)
	}

	// Static file serving (for dashboard)
	s.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/dist/")))
}

// Middleware for logging and error handling
func (s *APIServer) withLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create request context with logging
		ctx := context.WithValue(r.Context(), "request_id", fmt.Sprintf("%d", start.UnixNano()))
		r = r.WithContext(ctx)

		// Call next handler
		next.ServeHTTP(w, r)

		// Log request
		s.logger.Info(r.Context(), "API request", map[string]interface{}{
			"method":   r.Method,
			"path":     r.URL.Path,
			"duration": time.Since(start).String(),
			"remote":   r.RemoteAddr,
		})
	}
}

// Helper function to send JSON response
func (s *APIServer) sendJSON(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Success:   statusCode < 400,
		Data:      data,
		Timestamp: time.Now(),
		RequestID: r.Context().Value("request_id").(string),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.logger.Error(r.Context(), "Failed to encode JSON response", err)
	}
}

// Helper function to send error response
func (s *APIServer) sendError(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Success:   false,
		Error:     message,
		Timestamp: time.Now(),
		RequestID: r.Context().Value("request_id").(string),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.logger.Error(r.Context(), "Failed to encode error response", err)
	}
}

// Health check endpoint
func (s *APIServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"services": map[string]bool{
			"hft_engine":          s.hftEngine != nil,
			"binance_service":     s.binanceService != nil,
			"tradingview_service": s.tradingViewService != nil,
			"mcp_service":         s.mcpService != nil,
			"strategy_engine":     s.strategyEngine != nil,
			"exchange_manager":    s.exchangeManager != nil,
			"order_manager":       s.orderManager != nil,
			"market_data_service": s.marketDataService != nil,
		},
	}

	s.sendJSON(w, r, http.StatusOK, health)
}

// WebSocket handler
func (s *APIServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Error(r.Context(), "WebSocket upgrade failed", err)
		return
	}
	defer conn.Close()

	// Register client
	s.wsClients[conn] = true

	s.logger.Info(r.Context(), "WebSocket client connected", map[string]interface{}{
		"remote":  r.RemoteAddr,
		"clients": len(s.wsClients),
	})

	// Handle client messages
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			delete(s.wsClients, conn)
			s.logger.Info(r.Context(), "WebSocket client disconnected", map[string]interface{}{
				"remote":  r.RemoteAddr,
				"clients": len(s.wsClients),
			})
			break
		}
	}
}

// WebSocket hub for broadcasting messages
func (s *APIServer) handleWebSocketHub() {
	for {
		select {
		case message := <-s.wsBroadcast:
			for client := range s.wsClients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					client.Close()
					delete(s.wsClients, client)
				}
			}
		}
	}
}

// BroadcastMessage sends a message to all WebSocket clients
func (s *APIServer) BroadcastMessage(messageType string, payload interface{}) {
	if !s.config.EnableWebSocket {
		return
	}

	message := map[string]interface{}{
		"type":      messageType,
		"payload":   payload,
		"timestamp": time.Now(),
	}

	data, err := json.Marshal(message)
	if err != nil {
		s.logger.Error(context.Background(), "Failed to marshal WebSocket message", err)
		return
	}

	select {
	case s.wsBroadcast <- data:
	default:
		// Channel is full, skip message
	}
}
