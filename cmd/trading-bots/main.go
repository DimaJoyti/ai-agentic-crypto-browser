package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ai-agentic-browser/api"
	"github.com/ai-agentic-browser/internal/trading"
	"github.com/ai-agentic-browser/internal/trading/monitoring"
	"github.com/ai-agentic-browser/internal/trading/strategies"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gopkg.in/yaml.v3"
)

// TradingBotsConfig holds configuration for the trading bots system
type TradingBotsConfig struct {
	Server struct {
		Host         string        `yaml:"host"`
		Port         int           `yaml:"port"`
		ReadTimeout  time.Duration `yaml:"read_timeout"`
		WriteTimeout time.Duration `yaml:"write_timeout"`
		IdleTimeout  time.Duration `yaml:"idle_timeout"`
	} `yaml:"server"`

	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		SSLMode  string `yaml:"ssl_mode"`
	} `yaml:"database"`

	Redis struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`

	TradingBots struct {
		MaxConcurrentBots         int           `yaml:"max_concurrent_bots"`
		ExecutionInterval         time.Duration `yaml:"execution_interval"`
		OrderTimeout              time.Duration `yaml:"order_timeout"`
		RetryAttempts             int           `yaml:"retry_attempts"`
		PerformanceUpdateInterval time.Duration `yaml:"performance_update_interval"`
		HealthCheckInterval       time.Duration `yaml:"health_check_interval"`
	} `yaml:"trading_bots"`

	Exchanges map[string]ExchangeConfig `yaml:"exchanges"`
}

// ExchangeConfig holds configuration for a cryptocurrency exchange
type ExchangeConfig struct {
	APIURL     string `yaml:"api_url"`
	TestnetURL string `yaml:"testnet_url"`
	RateLimit  int    `yaml:"rate_limit"`
	Sandbox    bool   `yaml:"sandbox"`
	APIKey     string `yaml:"api_key"`
	APISecret  string `yaml:"api_secret"`
	Passphrase string `yaml:"passphrase"`
}

func main() {
	// Load configuration
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	obsConfig := observability.GetDefaultSimpleConfig()
	obsConfig.ServiceName = "trading-bots"
	obsProvider, err := observability.NewSimpleObservabilityProvider(obsConfig)
	if err != nil {
		log.Fatalf("Failed to initialize observability: %v", err)
	}
	logger := obsProvider.Logger

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize risk management system
	riskManager := trading.NewBotRiskManager(logger)
	if err := riskManager.Start(ctx); err != nil {
		log.Fatalf("Failed to start risk manager: %v", err)
	}

	// Initialize trading bot engine
	botEngineConfig := &trading.BotEngineConfig{
		MaxConcurrentBots:         config.TradingBots.MaxConcurrentBots,
		ExecutionInterval:         config.TradingBots.ExecutionInterval,
		OrderTimeout:              config.TradingBots.OrderTimeout,
		RetryAttempts:             config.TradingBots.RetryAttempts,
		PerformanceUpdateInterval: config.TradingBots.PerformanceUpdateInterval,
		HealthCheckInterval:       config.TradingBots.HealthCheckInterval,
	}

	botEngine := trading.NewTradingBotEngine(logger, botEngineConfig)

	// Initialize strategy manager
	strategyManager := strategies.NewStrategyManager(logger)

	// Create default strategies
	if err := strategyManager.CreateDefaultStrategies(); err != nil {
		log.Fatalf("Failed to create default strategies: %v", err)
	}

	// Initialize monitoring system
	monitoringConfig := &monitoring.MonitoringConfig{
		MetricsInterval:      30 * time.Second,
		HealthCheckInterval:  60 * time.Second,
		AlertCheckInterval:   10 * time.Second,
		MetricsRetention:     24 * time.Hour,
		AlertRetention:       7 * 24 * time.Hour,
		SnapshotRetention:    24 * time.Hour,
		EnableRealTimeAlerts: true,
		EnableDashboard:      true,
		EnableMetricsExport:  true,
		EnableProfiling:      false,
	}

	monitor := monitoring.NewTradingBotMonitor(logger, monitoringConfig, botEngine, riskManager)
	if err := monitor.Start(ctx); err != nil {
		log.Fatalf("Failed to start monitoring system: %v", err)
	}

	// Initialize API handlers
	tradingBotHandler := api.NewTradingBotHandler(logger, botEngine, strategyManager)
	riskManagementHandler := api.NewRiskManagementHandler(logger, riskManager)
	monitoringHandler := api.NewMonitoringHandler(logger, monitor)

	// Setup HTTP server
	router := mux.NewRouter()

	// Register API routes
	tradingBotHandler.RegisterRoutes(router)
	riskManagementHandler.RegisterRoutes(router)
	monitoringHandler.RegisterRoutes(router)

	// Add health check endpoint
	router.HandleFunc("/health", healthCheckHandler).Methods("GET")
	router.HandleFunc("/api/v1/health", healthCheckHandler).Methods("GET")

	// Add metrics endpoint
	router.HandleFunc("/metrics", metricsHandler).Methods("GET")

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port),
		Handler:      handler,
		ReadTimeout:  config.Server.ReadTimeout,
		WriteTimeout: config.Server.WriteTimeout,
		IdleTimeout:  config.Server.IdleTimeout,
	}

	// Start trading bot engine
	if err := botEngine.Start(ctx); err != nil {
		log.Fatalf("Failed to start trading bot engine: %v", err)
	}

	// Start HTTP server in a goroutine
	go func() {
		logger.Info(ctx, "Starting trading bots server", map[string]interface{}{
			"address": server.Addr,
		})

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	logger.Info(ctx, "Received shutdown signal", nil)

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Stop trading bot engine
	if err := botEngine.Stop(shutdownCtx); err != nil {
		logger.Error(shutdownCtx, "Failed to stop trading bot engine", err, nil)
	}

	// Stop risk management system
	if err := riskManager.Stop(shutdownCtx); err != nil {
		logger.Error(shutdownCtx, "Failed to stop risk manager", err, nil)
	}

	// Stop monitoring system
	if err := monitor.Stop(shutdownCtx); err != nil {
		logger.Error(shutdownCtx, "Failed to stop monitoring system", err, nil)
	}

	// Stop HTTP server
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error(shutdownCtx, "Failed to shutdown server", err, nil)
	}

	logger.Info(ctx, "Trading bots system shutdown complete", nil)
}

// loadConfig loads configuration from file and environment variables
func loadConfig() (*TradingBotsConfig, error) {
	config := &TradingBotsConfig{}

	// Load from YAML file
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "configs/trading-bots.yaml"
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Override with environment variables
	if host := os.Getenv("SERVER_HOST"); host != "" {
		config.Server.Host = host
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		if p, err := parsePort(port); err == nil {
			config.Server.Port = p
		}
	}

	// Set defaults
	if config.Server.Host == "" {
		config.Server.Host = "0.0.0.0"
	}
	if config.Server.Port == 0 {
		config.Server.Port = 8090
	}
	if config.Server.ReadTimeout == 0 {
		config.Server.ReadTimeout = 30 * time.Second
	}
	if config.Server.WriteTimeout == 0 {
		config.Server.WriteTimeout = 30 * time.Second
	}
	if config.Server.IdleTimeout == 0 {
		config.Server.IdleTimeout = 120 * time.Second
	}

	// Trading bots defaults
	if config.TradingBots.MaxConcurrentBots == 0 {
		config.TradingBots.MaxConcurrentBots = 7
	}
	if config.TradingBots.ExecutionInterval == 0 {
		config.TradingBots.ExecutionInterval = 5 * time.Second
	}
	if config.TradingBots.OrderTimeout == 0 {
		config.TradingBots.OrderTimeout = 30 * time.Second
	}
	if config.TradingBots.RetryAttempts == 0 {
		config.TradingBots.RetryAttempts = 3
	}
	if config.TradingBots.PerformanceUpdateInterval == 0 {
		config.TradingBots.PerformanceUpdateInterval = 1 * time.Minute
	}
	if config.TradingBots.HealthCheckInterval == 0 {
		config.TradingBots.HealthCheckInterval = 30 * time.Second
	}

	return config, nil
}

// parsePort parses a port string to integer
func parsePort(portStr string) (int, error) {
	var port int
	if _, err := fmt.Sscanf(portStr, "%d", &port); err != nil {
		return 0, err
	}
	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("invalid port: %d", port)
	}
	return port, nil
}

// healthCheckHandler handles health check requests
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":    "healthy",
		"service":   "trading-bots",
		"timestamp": time.Now().UTC(),
		"version":   "1.0.0",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// metricsHandler handles metrics requests
func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	// Basic metrics - in production this would integrate with Prometheus
	metrics := `# HELP trading_bots_total Total number of trading bots
# TYPE trading_bots_total gauge
trading_bots_total 7

# HELP trading_bots_active Number of active trading bots
# TYPE trading_bots_active gauge
trading_bots_active 0

# HELP trading_bots_errors_total Total number of trading bot errors
# TYPE trading_bots_errors_total counter
trading_bots_errors_total 0

# HELP trading_bots_trades_total Total number of trades executed
# TYPE trading_bots_trades_total counter
trading_bots_trades_total 0
`

	w.Write([]byte(metrics))
}

// Additional helper functions and middleware can be added here

// corsMiddleware adds CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware logs HTTP requests
func loggingMiddleware(logger *observability.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			logger.Info(r.Context(), "HTTP request", map[string]interface{}{
				"method":   r.Method,
				"path":     r.URL.Path,
				"duration": time.Since(start).String(),
				"remote":   r.RemoteAddr,
			})
		})
	}
}

// recoveryMiddleware recovers from panics
func recoveryMiddleware(logger *observability.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error(r.Context(), "Panic recovered", fmt.Errorf("%v", err), map[string]interface{}{
						"method": r.Method,
						"path":   r.URL.Path,
					})

					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
