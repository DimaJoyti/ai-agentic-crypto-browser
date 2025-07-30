package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ai-agentic-browser/api"
	"github.com/ai-agentic-browser/internal/binance"
	"github.com/ai-agentic-browser/internal/config"
	"github.com/ai-agentic-browser/internal/hft"
	"github.com/ai-agentic-browser/internal/mcp"
	"github.com/ai-agentic-browser/internal/tradingview"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/ai-agentic-browser/pkg/strategies"
	"github.com/shopspring/decimal"
)

func main() {
	ctx := context.Background()

	// Initialize logger
	logger := observability.NewLogger(config.ObservabilityConfig{
		ServiceName: "ai-agentic-browser",
		LogLevel:    "info",
		LogFormat:   "json",
	})

	logger.Info(ctx, "Starting AI-Agentic Crypto Browser", map[string]interface{}{
		"version": "1.0.0",
		"build":   "development",
	})

	// Initialize HFT Engine
	hftConfig := hft.HFTConfig{
		MaxOrdersPerSecond:   1000,
		LatencyTargetMicros:  1000,
		MaxPositionSize:      decimal.NewFromFloat(1.0),
		RiskLimitPercent:     0.02,
		EnableMarketMaking:   true,
		EnableArbitrage:      true,
		TickerUpdateInterval: time.Second,
	}

	hftEngine := hft.NewHFTEngine(logger, hftConfig)

	// Initialize Binance Service
	binanceServiceConfig := binance.ServiceConfig{
		Binance: binance.Config{
			APIKey:    os.Getenv("BINANCE_API_KEY"),
			SecretKey: os.Getenv("BINANCE_SECRET_KEY"),
			BaseURL:   "https://testnet.binance.vision", // Use testnet for development
			WSBaseURL: "wss://testnet.binance.vision/ws",
			Testnet:   true,
			RateLimit: 1200,
		},
		Symbols:    []string{"BTCUSDT", "ETHUSDT", "BNBUSDT"},
		Streams:    []string{"ticker", "depth", "trade"},
		EnableHFT:  true,
		BufferSize: 10000,
	}

	binanceService := binance.NewService(logger, binanceServiceConfig)

	// Initialize TradingView Service
	tradingViewConfig := tradingview.ServiceConfig{
		TradingView: tradingview.Config{
			BaseURL:          "https://www.tradingview.com",
			Headless:         true,
			Timeout:          30 * time.Second,
			UpdateInterval:   time.Minute,
			EnableSignals:    true,
			EnableIndicators: true,
		},
		SignalProcessor: tradingview.ProcessorConfig{
			QueueSize:           1000,
			ProcessingInterval:  time.Second,
			SignalTimeout:       5 * time.Minute,
			MinConfidence:       0.6,
			MaxSignalsPerSymbol: 10,
			EnableFiltering:     true,
			EnableAggregation:   true,
		},
		DefaultSymbols:    []string{"BINANCE:BTCUSDT", "BINANCE:ETHUSDT"},
		DefaultTimeframes: []string{"1m", "5m", "15m"},
		AutoStart:         true,
	}

	tradingViewService := tradingview.NewService(logger, tradingViewConfig)

	// Initialize MCP Integration Service
	mcpConfig := mcp.Config{
		CryptoAnalysis: mcp.CryptoAnalysisConfig{
			Symbols:           []string{"BTC", "ETH", "BNB"},
			UpdateInterval:    30 * time.Second,
			EnablePredictions: true,
			EnableIndicators:  true,
		},
		UpdateInterval: 10 * time.Second,
		EnableRealtime: true,
		BufferSize:     10000,
	}

	mcpService := mcp.NewIntegrationService(logger, mcpConfig)

	// Initialize Strategy Engine
	strategyConfig := strategies.EngineConfig{
		MaxStrategies:        50,
		ExecutionInterval:    100 * time.Millisecond,
		PerformanceWindow:    24 * time.Hour,
		EnableRiskManagement: true,
		MaxPositionSize:      decimal.NewFromFloat(1.0),
		MaxDailyLoss:         decimal.NewFromFloat(5000.0),
	}

	strategyEngine := strategies.NewStrategyEngine(logger, strategyConfig)

	// Initialize API Server
	apiConfig := api.Config{
		Host:            "localhost",
		Port:            8080,
		ReadTimeout:     30 * time.Second,
		WriteTimeout:    30 * time.Second,
		EnableCORS:      true,
		EnableWebSocket: true,
		RateLimit:       1000,
	}

	apiServer := api.NewAPIServer(logger, apiConfig)
	apiServer.SetServices(hftEngine, binanceService, tradingViewService, mcpService, strategyEngine)

	// Start all services
	logger.Info(ctx, "Starting services...", nil)

	// Start HFT Engine
	if err := hftEngine.Start(ctx); err != nil {
		log.Fatalf("Failed to start HFT engine: %v", err)
	}

	// Start Binance Service
	if err := binanceService.Start(ctx); err != nil {
		log.Fatalf("Failed to start Binance service: %v", err)
	}

	// Start TradingView Service
	if err := tradingViewService.Start(ctx); err != nil {
		log.Fatalf("Failed to start TradingView service: %v", err)
	}

	// Start MCP Service
	if err := mcpService.Start(ctx); err != nil {
		log.Fatalf("Failed to start MCP service: %v", err)
	}

	// Start Strategy Engine
	if err := strategyEngine.Start(ctx); err != nil {
		log.Fatalf("Failed to start strategy engine: %v", err)
	}

	// Start API Server
	if err := apiServer.Start(ctx); err != nil {
		log.Fatalf("Failed to start API server: %v", err)
	}

	logger.Info(ctx, "All services started successfully", map[string]interface{}{
		"api_server": fmt.Sprintf("http://%s:%d", apiConfig.Host, apiConfig.Port),
		"websocket":  fmt.Sprintf("ws://%s:%d/ws/trading", apiConfig.Host, apiConfig.Port),
	})

	// Set up graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	<-sigChan

	logger.Info(ctx, "Shutting down services...", nil)

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Stop all services in reverse order
	if err := apiServer.Stop(shutdownCtx); err != nil {
		logger.Error(shutdownCtx, "Failed to stop API server", err)
	}

	if err := strategyEngine.Stop(shutdownCtx); err != nil {
		logger.Error(shutdownCtx, "Failed to stop strategy engine", err)
	}

	if err := mcpService.Stop(shutdownCtx); err != nil {
		logger.Error(shutdownCtx, "Failed to stop MCP service", err)
	}

	if err := tradingViewService.Stop(shutdownCtx); err != nil {
		logger.Error(shutdownCtx, "Failed to stop TradingView service", err)
	}

	if err := binanceService.Stop(shutdownCtx); err != nil {
		logger.Error(shutdownCtx, "Failed to stop Binance service", err)
	}

	if err := hftEngine.Stop(shutdownCtx); err != nil {
		logger.Error(shutdownCtx, "Failed to stop HFT engine", err)
	}

	logger.Info(ctx, "All services stopped successfully", nil)
}
