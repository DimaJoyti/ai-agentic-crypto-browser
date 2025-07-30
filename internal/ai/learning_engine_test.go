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

func TestLearningEngine(t *testing.T) {
	logger := &observability.Logger{}
	engine := NewLearningEngine(logger)
	require.NotNil(t, engine)

	t.Run("EngineInitialization", func(t *testing.T) {
		assert.NotNil(t, engine.logger)
		assert.NotNil(t, engine.config)
		assert.NotNil(t, engine.userProfiles)
		assert.NotNil(t, engine.marketPatterns)
		assert.NotNil(t, engine.performanceTracker)
		assert.NotNil(t, engine.adaptiveModels)
		assert.NotNil(t, engine.feedbackProcessor)

		// Check configuration
		assert.Equal(t, 0.01, engine.config.LearningRate)
		assert.Equal(t, 0.1, engine.config.AdaptationThreshold)
		assert.Equal(t, 50, engine.config.MinDataPoints)
		assert.True(t, engine.config.EnableOnlineLearning)
		assert.True(t, engine.config.EnablePatternLearning)
		assert.True(t, engine.config.EnableUserProfiling)
	})

	t.Run("UserBehaviorLearning", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.New()

		// Test learning from trading behavior
		tradingBehavior := &UserBehaviorData{
			Type:      "trade",
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"symbol":     "BTC",
				"side":       "buy",
				"amount":     1000.0,
				"risk_level": 0.7,
				"decision_factors": map[string]float64{
					"technical_analysis": 0.6,
					"sentiment_analysis": 0.4,
				},
			},
			Context: map[string]interface{}{
				"market_condition": "bullish",
				"volatility":       0.3,
			},
			Outcome:     "success",
			Performance: 0.05, // 5% gain
		}

		err := engine.LearnFromUserBehavior(ctx, userID, tradingBehavior)
		require.NoError(t, err)

		// Verify user profile was created and updated
		profile, err := engine.GetUserProfile(userID)
		require.NoError(t, err)
		require.NotNil(t, profile)

		assert.Equal(t, userID, profile.UserID)
		assert.Greater(t, profile.RiskTolerance, 0.5) // Should have increased from default
		assert.Equal(t, 1, profile.PerformanceMetrics.TotalTrades)
		assert.Greater(t, profile.PerformanceMetrics.WinRate, 0.0)
		assert.NotEmpty(t, profile.LearningHistory)

		// Check learning event was recorded
		lastEvent := profile.LearningHistory[len(profile.LearningHistory)-1]
		assert.Equal(t, "user_behavior", lastEvent.EventType)
		assert.Equal(t, "success", lastEvent.Outcome)
		assert.Greater(t, lastEvent.Impact, 0.0)
	})

	t.Run("MultipleUserBehaviorLearning", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.New()

		// Learn from multiple behaviors
		behaviors := []*UserBehaviorData{
			{
				Type:        "trade",
				Timestamp:   time.Now().Add(-2 * time.Hour),
				Data:        map[string]interface{}{"risk_level": 0.3},
				Outcome:     "success",
				Performance: 0.02,
			},
			{
				Type:        "trade",
				Timestamp:   time.Now().Add(-1 * time.Hour),
				Data:        map[string]interface{}{"risk_level": 0.8},
				Outcome:     "failure",
				Performance: -0.03,
			},
			{
				Type:        "analysis_request",
				Timestamp:   time.Now(),
				Data:        map[string]interface{}{"analysis_type": "technical"},
				Outcome:     "neutral",
				Performance: 0.0,
			},
		}

		for _, behavior := range behaviors {
			err := engine.LearnFromUserBehavior(ctx, userID, behavior)
			require.NoError(t, err)
		}

		// Verify profile evolution
		profile, err := engine.GetUserProfile(userID)
		require.NoError(t, err)

		assert.Equal(t, 3, len(profile.LearningHistory))
		assert.Equal(t, 2, profile.PerformanceMetrics.TotalTrades)
		assert.InDelta(t, 0.5, profile.PerformanceMetrics.WinRate, 0.1) // 1 win, 1 loss
	})

	t.Run("UserProfileCreation", func(t *testing.T) {
		userID := uuid.New()
		profile := engine.createNewUserProfile(userID)

		assert.Equal(t, userID, profile.UserID)
		assert.Equal(t, 0.5, profile.RiskTolerance) // Default moderate
		assert.Equal(t, "moderate", profile.TradingStyle)
		assert.NotNil(t, profile.TradingPatterns)
		assert.NotNil(t, profile.DecisionFactors)
		assert.NotNil(t, profile.PerformanceMetrics)
		assert.NotNil(t, profile.Preferences)
		assert.Equal(t, 0.5, profile.BehaviorScore)

		// Check default trading patterns
		assert.Equal(t, 24*time.Hour, profile.TradingPatterns.AvgHoldingPeriod)
		assert.Equal(t, 1.0, profile.TradingPatterns.TradingFrequency)
		assert.NotNil(t, profile.TradingPatterns.PositionSizing)
		assert.NotNil(t, profile.TradingPatterns.RiskManagement)

		// Check default decision factors
		assert.Contains(t, profile.DecisionFactors, "technical_analysis")
		assert.Contains(t, profile.DecisionFactors, "sentiment_analysis")
		assert.Contains(t, profile.DecisionFactors, "fundamental_analysis")

		// Check default preferences
		assert.True(t, profile.Preferences.NotificationSettings["price_alerts"])
		assert.True(t, profile.Preferences.AnalysisPreferences["technical_analysis"])
	})

	t.Run("LearningImpactCalculation", func(t *testing.T) {
		// Test different behavior types and their impact
		tradeBehavior := &UserBehaviorData{
			Type:        "trade",
			Performance: 0.15, // Significant performance
		}
		tradeImpact := engine.calculateLearningImpact(tradeBehavior)
		assert.Greater(t, tradeImpact, 0.4) // Should be high impact

		analysisBehavior := &UserBehaviorData{
			Type:        "analysis_request",
			Performance: 0.02, // Small performance
		}
		analysisImpact := engine.calculateLearningImpact(analysisBehavior)
		assert.Less(t, analysisImpact, tradeImpact) // Should be lower impact

		feedbackBehavior := &UserBehaviorData{
			Type:        "feedback",
			Performance: 0.0,
		}
		feedbackImpact := engine.calculateLearningImpact(feedbackBehavior)
		assert.Equal(t, 0.1, feedbackImpact) // Base impact only
	})

	t.Run("UpdateWithDecay", func(t *testing.T) {
		// Test exponential moving average update
		current := 0.5
		new := 0.8
		learningRate := 0.1

		updated := engine.updateWithDecay(current, new, learningRate)
		expected := current*(1-learningRate) + new*learningRate
		assert.Equal(t, expected, updated)
		assert.Greater(t, updated, current) // Should move toward new value
		assert.Less(t, updated, new)        // But not reach it completely
	})

	t.Run("GetNonExistentUserProfile", func(t *testing.T) {
		nonExistentUserID := uuid.New()
		profile, err := engine.GetUserProfile(nonExistentUserID)

		assert.Error(t, err)
		assert.Nil(t, profile)
		assert.Contains(t, err.Error(), "user profile not found")
	})

	t.Run("MarketPatternsAccess", func(t *testing.T) {
		patterns := engine.GetMarketPatterns()
		assert.NotNil(t, patterns)
		// Initially empty, but structure should be correct
		assert.IsType(t, map[string]*MarketPattern{}, patterns)
	})

	t.Run("PerformanceMetricsAccess", func(t *testing.T) {
		metrics := engine.GetPerformanceMetrics()
		assert.NotNil(t, metrics)
		// Initially empty, but structure should be correct
		assert.IsType(t, map[string]*ModelPerformance{}, metrics)
	})
}

func TestAdaptiveModelManager(t *testing.T) {
	logger := &observability.Logger{}
	learningEngine := NewLearningEngine(logger)
	manager := NewAdaptiveModelManager(learningEngine, logger)
	require.NotNil(t, manager)

	t.Run("ManagerInitialization", func(t *testing.T) {
		assert.NotNil(t, manager.models)
		assert.NotNil(t, manager.learningEngine)
		assert.NotNil(t, manager.logger)
		assert.NotNil(t, manager.config)
		assert.NotNil(t, manager.adaptationQueue)

		// Check configuration
		assert.Equal(t, 1*time.Hour, manager.config.AdaptationInterval)
		assert.Equal(t, 0.05, manager.config.PerformanceThreshold)
		assert.Equal(t, 100, manager.config.MinAdaptationSamples)
		assert.True(t, manager.config.EnableAutoAdaptation)
	})

	t.Run("ModelRegistration", func(t *testing.T) {
		// Create a mock model
		pricePrediction := NewPricePredictionModel(logger)

		err := manager.RegisterAdaptiveModel("test_model", pricePrediction)
		require.NoError(t, err)

		// Verify model was registered
		models := manager.GetAdaptiveModels()
		assert.Contains(t, models, "test_model")

		model := models["test_model"]
		assert.Equal(t, "test_model", model.ID)
		assert.Equal(t, "Advanced Price Prediction Model", model.Name)
		assert.NotNil(t, model.BaseModel)
		assert.NotNil(t, model.Performance)
		assert.False(t, model.IsAdapting)
	})

	t.Run("AdaptationRequest", func(t *testing.T) {
		// Register a model first
		pricePrediction := NewPricePredictionModel(logger)
		manager.RegisterAdaptiveModel("adaptation_test_model", pricePrediction)

		request := &AdaptationRequest{
			ModelID:  "adaptation_test_model",
			Type:     "performance",
			Trigger:  "manual_request",
			Priority: 1,
			Data: map[string]interface{}{
				"performance_drop": 0.1,
			},
		}

		err := manager.RequestAdaptation(request)
		require.NoError(t, err)

		// Verify request was queued
		assert.Len(t, manager.adaptationQueue, 1)
		queuedRequest := manager.adaptationQueue[0]
		assert.Equal(t, "adaptation_test_model", queuedRequest.ModelID)
		assert.Equal(t, "performance", queuedRequest.Type)
		assert.Equal(t, "manual_request", queuedRequest.Trigger)
	})

	t.Run("AdaptationStrategySelection", func(t *testing.T) {
		performanceRequest := &AdaptationRequest{Type: "performance"}
		strategy := manager.selectAdaptationStrategy(performanceRequest)
		assert.NotNil(t, strategy)
		assert.Equal(t, "performance", strategy.GetType())

		feedbackRequest := &AdaptationRequest{Type: "feedback"}
		strategy = manager.selectAdaptationStrategy(feedbackRequest)
		assert.NotNil(t, strategy)
		assert.Equal(t, "feedback", strategy.GetType())

		driftRequest := &AdaptationRequest{Type: "drift"}
		strategy = manager.selectAdaptationStrategy(driftRequest)
		assert.NotNil(t, strategy)
		assert.Equal(t, "drift", strategy.GetType())

		unknownRequest := &AdaptationRequest{Type: "unknown"}
		strategy = manager.selectAdaptationStrategy(unknownRequest)
		assert.Nil(t, strategy)
	})

	t.Run("AdaptationHistory", func(t *testing.T) {
		// Register a model
		pricePrediction := NewPricePredictionModel(logger)
		modelID := "history_test_model"
		manager.RegisterAdaptiveModel(modelID, pricePrediction)

		// Initially no history
		history, err := manager.GetAdaptationHistory(modelID)
		require.NoError(t, err)
		assert.Empty(t, history)

		// Add some adaptation history manually
		model := manager.models[modelID]
		adaptation := ModelAdaptation{
			ID:          uuid.New().String(),
			Timestamp:   time.Now(),
			Type:        "test_adaptation",
			Description: "Test adaptation for history",
			Impact:      0.05,
			Success:     true,
		}
		model.Adaptations = append(model.Adaptations, adaptation)

		// Verify history retrieval
		history, err = manager.GetAdaptationHistory(modelID)
		require.NoError(t, err)
		assert.Len(t, history, 1)
		assert.Equal(t, "test_adaptation", history[0].Type)
		assert.True(t, history[0].Success)
	})

	t.Run("NonExistentModelHistory", func(t *testing.T) {
		history, err := manager.GetAdaptationHistory("non_existent_model")
		assert.Error(t, err)
		assert.Nil(t, history)
		assert.Contains(t, err.Error(), "model non_existent_model not found")
	})
}

func TestAdaptationStrategies(t *testing.T) {
	logger := &observability.Logger{}

	t.Run("PerformanceBasedAdaptation", func(t *testing.T) {
		strategy := &PerformanceBasedAdaptation{logger: logger}

		// Test with low performance model
		lowPerfModel := &AdaptiveModel{
			ID: "low_perf_model",
			Performance: &ModelPerformance{
				Accuracy: 0.7, // Below threshold
			},
		}

		request := &AdaptationRequest{Type: "performance"}

		canAdapt := strategy.CanAdapt(lowPerfModel, request)
		assert.True(t, canAdapt)

		// Test adaptation
		ctx := context.Background()
		result, err := strategy.Adapt(ctx, lowPerfModel, request)
		require.NoError(t, err)
		assert.True(t, result.Success)
		assert.Equal(t, "performance_optimization", result.Type)
		assert.Greater(t, result.PerformanceGain, 0.0)
		assert.NotEmpty(t, result.Changes)

		// Test with high performance model
		highPerfModel := &AdaptiveModel{
			ID: "high_perf_model",
			Performance: &ModelPerformance{
				Accuracy: 0.9, // Above threshold
			},
		}

		canAdapt = strategy.CanAdapt(highPerfModel, request)
		assert.False(t, canAdapt)
	})

	t.Run("FeedbackBasedAdaptation", func(t *testing.T) {
		strategy := &FeedbackBasedAdaptation{logger: logger}

		model := &AdaptiveModel{ID: "feedback_model"}

		// Test without feedback data
		requestWithoutFeedback := &AdaptationRequest{
			Type: "feedback",
			Data: map[string]interface{}{},
		}

		canAdapt := strategy.CanAdapt(model, requestWithoutFeedback)
		assert.False(t, canAdapt)

		// Test with feedback data
		requestWithFeedback := &AdaptationRequest{
			Type: "feedback",
			Data: map[string]interface{}{
				"feedback": map[string]interface{}{
					"rating":   4,
					"comments": "Good predictions",
				},
			},
		}

		canAdapt = strategy.CanAdapt(model, requestWithFeedback)
		assert.True(t, canAdapt)

		// Test adaptation
		ctx := context.Background()
		result, err := strategy.Adapt(ctx, model, requestWithFeedback)
		require.NoError(t, err)
		assert.True(t, result.Success)
		assert.Equal(t, "feedback_integration", result.Type)
		assert.Greater(t, result.PerformanceGain, 0.0)
	})

	t.Run("DriftBasedAdaptation", func(t *testing.T) {
		driftDetector := NewConceptDriftDetector()
		strategy := &DriftBasedAdaptation{
			logger:        logger,
			driftDetector: driftDetector,
		}

		model := &AdaptiveModel{ID: "drift_model"}
		request := &AdaptationRequest{Type: "drift"}

		// Initially no drift
		canAdapt := strategy.CanAdapt(model, request)
		assert.False(t, canAdapt)

		// Simulate drift detection
		driftDetector.driftHistory = append(driftDetector.driftHistory, DriftEvent{
			Type:       "gradual",
			DetectedAt: time.Now(),
		})

		canAdapt = strategy.CanAdapt(model, request)
		assert.True(t, canAdapt)

		// Test adaptation
		ctx := context.Background()
		result, err := strategy.Adapt(ctx, model, request)
		require.NoError(t, err)
		assert.True(t, result.Success)
		assert.Equal(t, "drift_compensation", result.Type)
		assert.Greater(t, result.PerformanceGain, 0.0)
	})
}

func TestConceptDriftDetector(t *testing.T) {
	detector := NewConceptDriftDetector()
	require.NotNil(t, detector)

	t.Run("DetectorInitialization", func(t *testing.T) {
		assert.Equal(t, 1000, detector.windowSize)
		assert.Equal(t, 0.1, detector.driftThreshold)
		assert.Empty(t, detector.recentData)
		assert.Empty(t, detector.referenceData)
		assert.Empty(t, detector.driftHistory)
	})

	t.Run("DataPointAddition", func(t *testing.T) {
		point := DataPoint{
			Features: map[string]float64{
				"price":     50000.0,
				"volume":    1000000.0,
				"sentiment": 0.6,
			},
			Timestamp: time.Now(),
			Weight:    1.0,
		}

		detector.AddDataPoint(point)
		assert.Len(t, detector.recentData, 1)
		assert.Equal(t, point.Features["price"], detector.recentData[0].Features["price"])
	})

	t.Run("WindowSizeMaintenance", func(t *testing.T) {
		// Add more points than window size
		for i := 0; i < 1200; i++ {
			point := DataPoint{
				Features: map[string]float64{
					"price": 50000.0 + float64(i),
				},
				Timestamp: time.Now().Add(time.Duration(i) * time.Minute),
				Weight:    1.0,
			}
			detector.AddDataPoint(point)
		}

		// Should maintain window size
		assert.Equal(t, detector.windowSize, len(detector.recentData))
	})

	t.Run("StatisticalDistanceCalculation", func(t *testing.T) {
		data1 := []DataPoint{
			{Features: map[string]float64{"price": 50000.0, "volume": 1000.0}},
			{Features: map[string]float64{"price": 51000.0, "volume": 1100.0}},
		}

		data2 := []DataPoint{
			{Features: map[string]float64{"price": 50000.0, "volume": 1000.0}},
			{Features: map[string]float64{"price": 51000.0, "volume": 1100.0}},
		}

		// Identical data should have zero distance
		distance := detector.calculateStatisticalDistance(data1, data2)
		assert.Equal(t, 0.0, distance)

		// Different data should have non-zero distance
		data3 := []DataPoint{
			{Features: map[string]float64{"price": 60000.0, "volume": 2000.0}},
			{Features: map[string]float64{"price": 61000.0, "volume": 2100.0}},
		}

		distance = detector.calculateStatisticalDistance(data1, data3)
		assert.Greater(t, distance, 0.0)
	})

	t.Run("DriftDetection", func(t *testing.T) {
		// Initially no drift
		assert.False(t, detector.HasDrift())

		// Add drift event
		detector.driftHistory = append(detector.driftHistory, DriftEvent{
			Type:       "sudden",
			DetectedAt: time.Now(),
			Severity:   0.8,
		})

		// Should detect drift
		assert.True(t, detector.HasDrift())

		// Old drift should not be detected
		detector.driftHistory = []DriftEvent{
			{
				Type:       "old_drift",
				DetectedAt: time.Now().Add(-48 * time.Hour), // 2 days ago
				Severity:   0.8,
			},
		}

		assert.False(t, detector.HasDrift())
	})
}
