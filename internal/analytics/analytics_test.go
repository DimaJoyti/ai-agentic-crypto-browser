package analytics

import (
	"context"
	"testing"
	"time"

	"github.com/ai-agentic-browser/internal/config"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/shopspring/decimal"
)

func TestPerformanceEngine(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{
		ServiceName: "test",
		LogLevel:    "info",
	})

	config := PerformanceConfig{
		AnalysisInterval:  time.Second,
		MetricsBufferSize: 1000,
		AlertThresholds: AlertThresholds{
			LatencyThreshold:     100 * time.Millisecond,
			ThroughputThreshold:  1000,
			ErrorRateThreshold:   1.0,
			SharpeRatioThreshold: 1.0,
			DrawdownThreshold:    0.15,
		},
	}

	engine := NewPerformanceEngine(logger, config)

	ctx := context.Background()
	err := engine.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start performance engine: %v", err)
	}

	// Test getting metrics
	metrics := engine.GetMetrics()
	if metrics.Trading.TotalTrades < 0 {
		t.Error("Expected non-negative total trades")
	}

	err = engine.Stop(ctx)
	if err != nil {
		t.Fatalf("Failed to stop performance engine: %v", err)
	}
}

func TestTradingPerformanceAnalyzer(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{
		ServiceName: "test",
		LogLevel:    "info",
	})

	config := PerformanceConfig{
		MetricsBufferSize: 1000,
	}

	analyzer := NewTradingPerformanceAnalyzer(logger, config)

	ctx := context.Background()
	err := analyzer.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start trading analyzer: %v", err)
	}

	// Add a test trade
	trade := TradeRecord{
		ID:         "test_001",
		Symbol:     "BTCUSD",
		Side:       "BUY",
		Quantity:   decimal.NewFromFloat(0.1),
		EntryPrice: decimal.NewFromFloat(50000),
		ExitPrice:  decimal.NewFromFloat(51000),
		EntryTime:  time.Now().Add(-time.Hour),
		ExitTime:   time.Now(),
		PnL:        decimal.NewFromFloat(100),
		IsWin:      true,
	}

	analyzer.AddTrade(trade)

	// Test getting metrics
	metrics := analyzer.GetMetrics()
	if metrics.TotalTrades != 1 {
		t.Errorf("Expected 1 trade, got %d", metrics.TotalTrades)
	}

	err = analyzer.Stop(ctx)
	if err != nil {
		t.Fatalf("Failed to stop trading analyzer: %v", err)
	}
}

func TestSystemPerformanceAnalyzer(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{
		ServiceName: "test",
		LogLevel:    "info",
	})

	perfConfig := PerformanceConfig{
		AnalysisInterval:  100 * time.Millisecond,
		MetricsBufferSize: 100,
	}

	analyzer := NewSystemPerformanceAnalyzer(logger, perfConfig)

	ctx := context.Background()
	err := analyzer.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start system analyzer: %v", err)
	}

	// Wait for some metrics to be collected
	time.Sleep(200 * time.Millisecond)

	// Test getting metrics
	metrics := analyzer.GetMetrics()
	if metrics.CPUUsage < 0 {
		t.Error("Expected non-negative CPU usage")
	}

	err = analyzer.Stop(ctx)
	if err != nil {
		t.Fatalf("Failed to stop system analyzer: %v", err)
	}
}

func TestPortfolioPerformanceAnalyzer(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{
		ServiceName: "test",
		LogLevel:    "info",
	})

	perfConfig := PerformanceConfig{
		MetricsBufferSize: 1000,
	}

	analyzer := NewPortfolioPerformanceAnalyzer(logger, perfConfig)

	ctx := context.Background()
	err := analyzer.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start portfolio analyzer: %v", err)
	}

	// Add a test position
	position := Position{
		Symbol:       "BTCUSD",
		Quantity:     decimal.NewFromFloat(0.5),
		AveragePrice: decimal.NewFromFloat(50000),
		CurrentPrice: decimal.NewFromFloat(51000),
		MarketValue:  decimal.NewFromFloat(25500),
		Weight:       0.5,
	}

	analyzer.UpdatePosition(position)

	// Test getting positions
	positions := analyzer.GetPositions()
	if len(positions) != 1 {
		t.Errorf("Expected 1 position, got %d", len(positions))
	}

	err = analyzer.Stop(ctx)
	if err != nil {
		t.Fatalf("Failed to stop portfolio analyzer: %v", err)
	}
}

func TestBenchmarkEngine(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{
		ServiceName: "test",
		LogLevel:    "info",
	})

	perfConfig := PerformanceConfig{
		MetricsBufferSize: 1000,
	}

	engine := NewBenchmarkEngine(logger, perfConfig)

	ctx := context.Background()
	err := engine.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start benchmark engine: %v", err)
	}

	// Test getting benchmarks
	benchmarks := engine.GetBenchmarks()
	if len(benchmarks) == 0 {
		t.Error("Expected at least one benchmark")
	}

	// Test comparison
	portfolioReturns := []decimal.Decimal{
		decimal.NewFromFloat(0.01),
		decimal.NewFromFloat(0.02),
		decimal.NewFromFloat(-0.01),
	}

	comparison, err := engine.CompareToBenchmark(portfolioReturns, "btc_index", "test")
	if err != nil {
		t.Fatalf("Failed to compare to benchmark: %v", err)
	}

	if comparison.BenchmarkID != "btc_index" {
		t.Errorf("Expected benchmark ID 'btc_index', got %s", comparison.BenchmarkID)
	}

	err = engine.Stop(ctx)
	if err != nil {
		t.Fatalf("Failed to stop benchmark engine: %v", err)
	}
}

func TestOptimizationEngine(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{
		ServiceName: "test",
		LogLevel:    "info",
	})

	perfConfig := PerformanceConfig{
		MetricsBufferSize: 1000,
	}

	engine := NewOptimizationEngine(logger, perfConfig)

	ctx := context.Background()
	err := engine.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start optimization engine: %v", err)
	}

	// Test getting metrics
	metrics := engine.GetMetrics()
	if metrics.OptimizationScore < 0 || metrics.OptimizationScore > 100 {
		t.Errorf("Expected optimization score between 0-100, got %f", metrics.OptimizationScore)
	}

	// Test getting recommendations
	recommendations := engine.GetRecommendations()
	if recommendations == nil {
		t.Error("Expected recommendations to be initialized")
	}

	err = engine.Stop(ctx)
	if err != nil {
		t.Fatalf("Failed to stop optimization engine: %v", err)
	}
}
