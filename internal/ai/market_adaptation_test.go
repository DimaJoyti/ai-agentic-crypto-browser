package ai

import (
	"context"
	"testing"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarketAdaptationEngine(t *testing.T) {
	logger := &observability.Logger{}

	t.Run("EngineInitialization", func(t *testing.T) {
		engine := NewMarketAdaptationEngine(logger)
		require.NotNil(t, engine)
		assert.NotNil(t, engine.logger)
		assert.NotNil(t, engine.config)
		assert.NotNil(t, engine.patternDetector)
		assert.NotNil(t, engine.strategyManager)
		assert.NotNil(t, engine.performanceAnalyzer)
		assert.NotNil(t, engine.adaptationRules)
		assert.NotNil(t, engine.detectedPatterns)
		assert.NotNil(t, engine.adaptiveStrategies)
		assert.NotNil(t, engine.adaptationHistory)
		assert.NotNil(t, engine.performanceMetrics)

		// Check configuration
		assert.Equal(t, 7*24*time.Hour, engine.config.PatternDetectionWindow)
		assert.Equal(t, 0.7, engine.config.AdaptationThreshold)
		assert.Equal(t, 3, engine.config.MinPatternOccurrences)
		assert.Equal(t, 1*time.Hour, engine.config.StrategyUpdateFrequency)
		assert.Equal(t, 24*time.Hour, engine.config.PerformanceEvaluationWindow)
		assert.Equal(t, 0.5, engine.config.RiskAdjustmentSensitivity)
		assert.Equal(t, 0.1, engine.config.AdaptationLearningRate)
		assert.Equal(t, 1000, engine.config.MaxAdaptationHistory)
		assert.True(t, engine.config.EnableRealTimeAdaptation)
		assert.Equal(t, 0.6, engine.config.ConfidenceThreshold)
	})

	t.Run("PatternDetection", func(t *testing.T) {
		engine := NewMarketAdaptationEngine(logger)
		require.NotNil(t, engine)
		ctx := context.Background()

		// Create sample market data
		marketData := map[string]interface{}{
			"prices":     []float64{50000, 50500, 51000, 51500, 52000, 52500, 53000, 53500, 54000, 54500},
			"volumes":    []float64{100, 120, 110, 130, 140, 125, 135, 145, 150, 160},
			"timestamps": []int64{1640995200, 1640998800, 1641002400, 1641006000, 1641009600, 1641013200, 1641016800, 1641020400, 1641024000, 1641027600},
		}

		patterns, err := engine.DetectPatterns(ctx, marketData)
		require.NoError(t, err)
		require.NotEmpty(t, patterns)

		// Verify pattern properties
		pattern := patterns[0]
		assert.NotEmpty(t, pattern.ID)
		assert.Equal(t, "trend", pattern.Type)
		assert.Equal(t, "Upward Trend", pattern.Name)
		assert.NotEmpty(t, pattern.Description)
		assert.Equal(t, "BTC", pattern.Asset)
		assert.Equal(t, "1h", pattern.TimeFrame)
		assert.Greater(t, pattern.Strength, 0.0)
		assert.Greater(t, pattern.Confidence, 0.0)
		assert.Greater(t, pattern.Duration, time.Duration(0))
		assert.NotNil(t, pattern.Characteristics)
		assert.NotEmpty(t, pattern.TriggerConditions)
		assert.NotNil(t, pattern.ExpectedOutcome)
		assert.NotNil(t, pattern.MarketContext)
		assert.Equal(t, 1, pattern.OccurrenceCount)

		// Verify characteristics
		assert.Contains(t, pattern.Characteristics, "slope")
		assert.Contains(t, pattern.Characteristics, "r_squared")
		assert.Contains(t, pattern.Characteristics, "trend_strength")
		assert.Contains(t, pattern.Characteristics, "momentum")

		// Verify trigger conditions
		trigger := pattern.TriggerConditions[0]
		assert.Equal(t, "price", trigger.Type)
		assert.Equal(t, "sma_cross", trigger.Indicator)
		assert.Equal(t, "cross_above", trigger.Operator)
		assert.Greater(t, trigger.Confidence, 0.0)

		// Verify expected outcome
		outcome := pattern.ExpectedOutcome
		assert.Equal(t, "up", outcome.Direction)
		assert.Greater(t, outcome.Magnitude, 0.0)
		assert.Greater(t, outcome.Probability, 0.0)
		assert.Greater(t, outcome.TimeHorizon, time.Duration(0))
		assert.Greater(t, outcome.Confidence, 0.0)
		assert.Greater(t, outcome.RiskReward, 0.0)

		// Verify market context
		context := pattern.MarketContext
		assert.Equal(t, "bull", context.MarketRegime)
		assert.Equal(t, "up", context.TrendDirection)
		assert.Greater(t, context.TrendStrength, 0.0)
		assert.Greater(t, context.MarketSentiment, 0.0)
		assert.NotEmpty(t, context.TechnicalIndicators)
	})

	t.Run("GetDetectedPatterns", func(t *testing.T) {
		engine := NewMarketAdaptationEngine(logger)
		require.NotNil(t, engine)
		ctx := context.Background()

		// First detect some patterns
		marketData := map[string]interface{}{
			"prices": []float64{50000, 50500, 51000, 51500, 52000, 52500, 53000, 53500, 54000, 54500},
		}

		detectedPatterns, err := engine.DetectPatterns(ctx, marketData)
		require.NoError(t, err)
		require.NotEmpty(t, detectedPatterns)

		// Test getting patterns without filters
		patterns, err := engine.GetDetectedPatterns(ctx, map[string]interface{}{})
		require.NoError(t, err)
		assert.NotEmpty(t, patterns)

		// Test getting patterns with asset filter
		filteredPatterns, err := engine.GetDetectedPatterns(ctx, map[string]interface{}{
			"asset": "BTC",
		})
		require.NoError(t, err)
		assert.NotEmpty(t, filteredPatterns)

		// Test getting patterns with type filter
		typeFilteredPatterns, err := engine.GetDetectedPatterns(ctx, map[string]interface{}{
			"type": "trend",
		})
		require.NoError(t, err)
		assert.NotEmpty(t, typeFilteredPatterns)

		// Test getting patterns with confidence filter
		confidenceFilteredPatterns, err := engine.GetDetectedPatterns(ctx, map[string]interface{}{
			"min_confidence": 0.5,
		})
		require.NoError(t, err)
		assert.NotEmpty(t, confidenceFilteredPatterns)

		// Test getting patterns with high confidence filter (should return empty)
		highConfidencePatterns, err := engine.GetDetectedPatterns(ctx, map[string]interface{}{
			"min_confidence": 0.99,
		})
		require.NoError(t, err)
		assert.Empty(t, highConfidencePatterns)
	})

	t.Run("AddAdaptiveStrategy", func(t *testing.T) {
		engine := NewMarketAdaptationEngine(logger)
		require.NotNil(t, engine)
		ctx := context.Background()

		strategy := &AdaptiveStrategy{
			Name:        "Test Trend Following Strategy",
			Description: "A test strategy for trend following",
			Type:        "trend_following",
			BaseParameters: map[string]float64{
				"position_size":   0.05,
				"stop_loss":       0.02,
				"take_profit":     0.04,
				"entry_threshold": 0.7,
			},
			CurrentParameters: map[string]float64{
				"position_size":   0.05,
				"stop_loss":       0.02,
				"take_profit":     0.04,
				"entry_threshold": 0.7,
			},
			PerformanceTargets: &PerformanceTargets{
				TargetReturn:       0.15,
				MaxDrawdown:        0.1,
				MinSharpeRatio:     1.0,
				MinWinRate:         0.6,
				MaxVolatility:      0.3,
				TargetProfitFactor: 1.5,
				EvaluationPeriod:   30 * 24 * time.Hour,
			},
			RiskLimits: &MarketRiskLimits{
				MaxPositionSize:    0.1,
				MaxLeverage:        2.0,
				StopLossPercentage: 0.05,
				TakeProfitRatio:    2.0,
				MaxDailyLoss:       0.02,
				VaRLimit:           0.01,
				ConcentrationLimit: 0.2,
			},
		}

		err := engine.AddAdaptiveStrategy(ctx, strategy)
		require.NoError(t, err)

		// Verify strategy was added
		assert.NotEmpty(t, strategy.ID)
		assert.True(t, strategy.IsActive)
		assert.Equal(t, 0, strategy.AdaptationCount)
		assert.False(t, strategy.LastAdaptation.IsZero())

		// Verify strategy is in the engine
		strategies, err := engine.GetAdaptiveStrategies(ctx)
		require.NoError(t, err)
		assert.Len(t, strategies, 1)
		assert.Equal(t, strategy.ID, strategies[0].ID)
	})

	t.Run("UpdateStrategyStatus", func(t *testing.T) {
		engine := NewMarketAdaptationEngine(logger)
		require.NotNil(t, engine)
		ctx := context.Background()

		// Add a strategy first
		strategy := &AdaptiveStrategy{
			Name: "Test Strategy",
			Type: "trend_following",
			CurrentParameters: map[string]float64{
				"position_size": 0.05,
			},
		}

		err := engine.AddAdaptiveStrategy(ctx, strategy)
		require.NoError(t, err)

		// Update strategy status to inactive
		err = engine.UpdateStrategyStatus(ctx, strategy.ID, false)
		require.NoError(t, err)

		// Verify status was updated
		strategies, err := engine.GetAdaptiveStrategies(ctx)
		require.NoError(t, err)
		assert.Len(t, strategies, 1)
		assert.False(t, strategies[0].IsActive)

		// Update strategy status back to active
		err = engine.UpdateStrategyStatus(ctx, strategy.ID, true)
		require.NoError(t, err)

		// Verify status was updated
		strategies, err = engine.GetAdaptiveStrategies(ctx)
		require.NoError(t, err)
		assert.Len(t, strategies, 1)
		assert.True(t, strategies[0].IsActive)

		// Test updating non-existent strategy
		err = engine.UpdateStrategyStatus(ctx, "non-existent", true)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "strategy not found")
	})

	t.Run("StrategyAdaptation", func(t *testing.T) {
		engine := NewMarketAdaptationEngine(logger)
		require.NotNil(t, engine)
		ctx := context.Background()

		// Add a strategy with performance targets
		strategy := &AdaptiveStrategy{
			Name: "Test Adaptation Strategy",
			Type: "trend_following",
			CurrentParameters: map[string]float64{
				"position_size":   0.05,
				"stop_loss":       0.02,
				"entry_threshold": 0.7,
			},
			PerformanceTargets: &PerformanceTargets{
				MinSharpeRatio: 1.0,
				MaxDrawdown:    0.1,
				MinWinRate:     0.6,
			},
			PerformanceMetrics: &MarketPerformanceMetrics{
				SharpeRatio: 0.5,  // Below target
				MaxDrawdown: 0.15, // Above target
				WinRate:     0.4,  // Below target
			},
		}

		err := engine.AddAdaptiveStrategy(ctx, strategy)
		require.NoError(t, err)

		// Create patterns that should trigger adaptation
		patterns := []*DetectedPattern{
			{
				ID:         uuid.New().String(),
				Type:       "trend",
				Confidence: 0.8, // Above adaptation threshold
				ExpectedOutcome: &ExpectedOutcome{
					Direction: "up",
				},
			},
		}

		// Perform strategy adaptation
		err = engine.AdaptStrategies(ctx, patterns)
		require.NoError(t, err)

		// Verify adaptation occurred
		strategies, err := engine.GetAdaptiveStrategies(ctx)
		require.NoError(t, err)
		assert.Len(t, strategies, 1)

		adaptedStrategy := strategies[0]
		assert.Greater(t, adaptedStrategy.AdaptationCount, 0)
		assert.NotEmpty(t, adaptedStrategy.AdaptationHistory)

		// Verify adaptation history
		history, err := engine.GetAdaptationHistory(ctx, 10)
		require.NoError(t, err)
		assert.NotEmpty(t, history)

		record := history[0]
		assert.NotEmpty(t, record.ID)
		assert.Equal(t, "strategy_adapted", record.Type)
		assert.Equal(t, strategy.ID, record.StrategyID)
		assert.NotEmpty(t, record.Description)
		assert.Greater(t, record.Confidence, 0.0)
		assert.False(t, record.Timestamp.IsZero())
	})

	t.Run("PerformanceMetrics", func(t *testing.T) {
		engine := NewMarketAdaptationEngine(logger)
		require.NotNil(t, engine)
		ctx := context.Background()

		// Add a strategy with performance metrics
		strategy := &AdaptiveStrategy{
			Name: "Performance Test Strategy",
			Type: "momentum",
			CurrentParameters: map[string]float64{
				"position_size": 0.05,
			},
		}

		err := engine.AddAdaptiveStrategy(ctx, strategy)
		require.NoError(t, err)

		// Add performance metrics
		metrics := &MarketPerformanceMetrics{
			StrategyID:       strategy.ID,
			TotalReturn:      0.15,
			AnnualizedReturn: 0.18,
			Volatility:       0.25,
			SharpeRatio:      0.72,
			MaxDrawdown:      0.08,
			WinRate:          0.65,
			ProfitFactor:     1.8,
			TotalTrades:      45,
			WinningTrades:    29,
			LosingTrades:     16,
			LastUpdated:      time.Now(),
		}

		engine.performanceMetrics[strategy.ID] = metrics

		// Test getting performance metrics
		retrievedMetrics, err := engine.GetPerformanceMetrics(ctx, strategy.ID)
		require.NoError(t, err)
		assert.Equal(t, metrics.StrategyID, retrievedMetrics.StrategyID)
		assert.Equal(t, metrics.TotalReturn, retrievedMetrics.TotalReturn)
		assert.Equal(t, metrics.SharpeRatio, retrievedMetrics.SharpeRatio)
		assert.Equal(t, metrics.WinRate, retrievedMetrics.WinRate)

		// Test getting metrics for non-existent strategy
		_, err = engine.GetPerformanceMetrics(ctx, "non-existent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "performance metrics not found")

		// Test with empty strategy ID
		_, err = engine.GetPerformanceMetrics(ctx, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "strategy ID is required")
	})

	t.Run("AdaptationHistoryLimit", func(t *testing.T) {
		engine := NewMarketAdaptationEngine(logger)
		require.NotNil(t, engine)
		ctx := context.Background()

		// Set a small history limit for testing
		originalLimit := engine.config.MaxAdaptationHistory
		engine.config.MaxAdaptationHistory = 3

		// Add a strategy
		strategy := &AdaptiveStrategy{
			Name: "History Test Strategy",
			Type: "trend_following",
			CurrentParameters: map[string]float64{
				"position_size": 0.05,
			},
			PerformanceMetrics: &MarketPerformanceMetrics{
				SharpeRatio: 0.3, // Low to trigger adaptations
			},
			PerformanceTargets: &PerformanceTargets{
				MinSharpeRatio: 1.0,
			},
		}

		err := engine.AddAdaptiveStrategy(ctx, strategy)
		require.NoError(t, err)

		// Create patterns to trigger multiple adaptations
		patterns := []*DetectedPattern{
			{
				ID:         uuid.New().String(),
				Type:       "trend",
				Confidence: 0.8,
			},
		}

		// Perform multiple adaptations
		for i := 0; i < 5; i++ {
			err = engine.AdaptStrategies(ctx, patterns)
			require.NoError(t, err)
		}

		// Verify history is limited
		history, err := engine.GetAdaptationHistory(ctx, 0)
		require.NoError(t, err)
		assert.LessOrEqual(t, len(history), 3)

		// Restore original limit
		engine.config.MaxAdaptationHistory = originalLimit
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		engine := NewMarketAdaptationEngine(logger)
		require.NotNil(t, engine)
		ctx := context.Background()

		// Test pattern detection with invalid data
		invalidData := map[string]interface{}{
			"invalid": "data",
		}

		patterns, err := engine.DetectPatterns(ctx, invalidData)
		require.NoError(t, err) // Should not error, but return empty patterns
		assert.Empty(t, patterns)

		// Test adaptation with empty patterns
		err = engine.AdaptStrategies(ctx, []*DetectedPattern{})
		require.NoError(t, err) // Should not error with empty patterns

		// Test getting strategies when none exist
		strategies, err := engine.GetAdaptiveStrategies(ctx)
		require.NoError(t, err)
		assert.Empty(t, strategies)

		// Test getting history when none exists
		history, err := engine.GetAdaptationHistory(ctx, 10)
		require.NoError(t, err)
		assert.Empty(t, history)
	})
}

func TestPatternDetector(t *testing.T) {
	logger := &observability.Logger{}
	detector := NewPatternDetector(logger)
	require.NotNil(t, detector)

	t.Run("DetectorInitialization", func(t *testing.T) {
		assert.NotNil(t, detector.config)
		assert.NotNil(t, detector.logger)

		// Check configuration
		assert.Equal(t, 5, detector.config.MinPatternLength)
		assert.Equal(t, 50, detector.config.MaxPatternLength)
		assert.Equal(t, 0.8, detector.config.SimilarityThreshold)
		assert.Equal(t, 0.7, detector.config.ConfidenceThreshold)
		assert.Equal(t, 1*time.Hour, detector.config.UpdateFrequency)
	})

	t.Run("PatternDetection", func(t *testing.T) {
		ctx := context.Background()

		// Test with valid price data
		marketData := map[string]interface{}{
			"prices": []float64{50000, 50500, 51000, 51500, 52000, 52500, 53000, 53500, 54000, 54500},
		}

		patterns, err := detector.DetectPatterns(ctx, marketData)
		require.NoError(t, err)
		assert.NotEmpty(t, patterns)

		pattern := patterns[0]
		assert.Equal(t, "trend", pattern.Type)
		assert.Equal(t, "BTC", pattern.Asset)
		assert.Greater(t, pattern.Confidence, 0.0)

		// Test with insufficient data
		insufficientData := map[string]interface{}{
			"prices": []float64{50000, 50500},
		}

		patterns, err = detector.DetectPatterns(ctx, insufficientData)
		require.NoError(t, err)
		assert.Empty(t, patterns)

		// Test with no price data
		noData := map[string]interface{}{
			"volumes": []float64{100, 200, 300},
		}

		patterns, err = detector.DetectPatterns(ctx, noData)
		require.NoError(t, err)
		assert.Empty(t, patterns)
	})
}

func TestStrategyManager(t *testing.T) {
	logger := &observability.Logger{}
	manager := NewStrategyManager(logger)
	require.NotNil(t, manager)

	t.Run("ManagerInitialization", func(t *testing.T) {
		assert.NotNil(t, manager.config)
		assert.NotNil(t, manager.logger)

		// Check configuration
		assert.Equal(t, 0.5, manager.config.AdaptationSensitivity)
		assert.Equal(t, 24*time.Hour, manager.config.MinPerformancePeriod)
		assert.Equal(t, 5, manager.config.MaxAdaptationsPerDay)
		assert.Equal(t, 0.6, manager.config.ConfidenceThreshold)
	})

	t.Run("StrategyAdaptation", func(t *testing.T) {
		ctx := context.Background()

		strategy := &AdaptiveStrategy{
			ID:   uuid.New().String(),
			Name: "Test Strategy",
			Type: "trend_following",
			CurrentParameters: map[string]float64{
				"position_size":   0.05,
				"stop_loss":       0.02,
				"entry_threshold": 0.7,
			},
		}

		patterns := []*DetectedPattern{
			{
				Type: "trend",
				ExpectedOutcome: &ExpectedOutcome{
					Direction: "up",
				},
			},
		}

		// Test adaptation for poor Sharpe ratio
		adaptation, err := manager.AdaptStrategy(ctx, strategy, patterns, "poor_sharpe_ratio")
		require.NoError(t, err)
		assert.NotNil(t, adaptation)
		assert.Equal(t, strategy.ID, adaptation.StrategyID)
		assert.Equal(t, "parameter_adjustment", adaptation.AdaptationType)
		assert.Equal(t, "poor_sharpe_ratio", adaptation.TriggerReason)
		assert.Greater(t, adaptation.Confidence, 0.0)
		assert.True(t, adaptation.Success)

		// Verify position size was reduced
		assert.Less(t, adaptation.NewParameters["position_size"], adaptation.OldParameters["position_size"])

		// Test adaptation for excessive drawdown
		adaptation, err = manager.AdaptStrategy(ctx, strategy, patterns, "excessive_drawdown")
		require.NoError(t, err)
		assert.Less(t, adaptation.NewParameters["stop_loss"], adaptation.OldParameters["stop_loss"])

		// Test adaptation for low win rate
		adaptation, err = manager.AdaptStrategy(ctx, strategy, patterns, "low_win_rate")
		require.NoError(t, err)
		assert.Greater(t, adaptation.NewParameters["entry_threshold"], adaptation.OldParameters["entry_threshold"])

		// Test default adaptation with patterns
		adaptation, err = manager.AdaptStrategy(ctx, strategy, patterns, "pattern_detected_trend")
		require.NoError(t, err)
		assert.Greater(t, adaptation.NewParameters["position_size"], adaptation.OldParameters["position_size"])
	})
}

func TestPerformanceAnalyzer(t *testing.T) {
	logger := &observability.Logger{}
	analyzer := NewPerformanceAnalyzer(logger)
	require.NotNil(t, analyzer)

	t.Run("AnalyzerInitialization", func(t *testing.T) {
		assert.NotNil(t, analyzer.config)
		assert.NotNil(t, analyzer.logger)

		// Check configuration
		assert.Equal(t, 1*time.Hour, analyzer.config.EvaluationFrequency)
		assert.Equal(t, "BTC", analyzer.config.BenchmarkAsset)
		assert.Equal(t, 0.02, analyzer.config.RiskFreeRate)
		assert.Equal(t, 0.95, analyzer.config.ConfidenceLevel)
	})
}
