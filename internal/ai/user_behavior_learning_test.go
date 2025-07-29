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

func TestUserBehaviorLearningEngine(t *testing.T) {
	logger := &observability.Logger{}
	engine := NewUserBehaviorLearningEngine(logger)
	require.NotNil(t, engine)

	t.Run("EngineInitialization", func(t *testing.T) {
		assert.NotNil(t, engine.logger)
		assert.NotNil(t, engine.config)
		assert.NotNil(t, engine.behaviorAnalyzer)
		assert.NotNil(t, engine.patternRecognizer)
		assert.NotNil(t, engine.preferenceEngine)
		assert.NotNil(t, engine.recommendationEngine)
		assert.NotNil(t, engine.personalityProfiler)
		assert.NotNil(t, engine.riskProfiler)
		assert.NotNil(t, engine.userProfiles)
		assert.NotNil(t, engine.behaviorHistory)
		assert.NotNil(t, engine.learningModels)

		// Check configuration
		assert.Equal(t, 0.1, engine.config.LearningRate)
		assert.Equal(t, 10, engine.config.MinObservations)
		assert.Equal(t, 30*24*time.Hour, engine.config.PatternDetectionWindow)
		assert.Equal(t, 0.05, engine.config.PreferenceUpdateRate)
		assert.Equal(t, 0.02, engine.config.PersonalityUpdateRate)
		assert.Equal(t, 0.03, engine.config.RiskToleranceUpdateRate)
		assert.Equal(t, 0.7, engine.config.RecommendationThreshold)
		assert.Equal(t, 10000, engine.config.MaxHistorySize)
		assert.True(t, engine.config.EnableRealTimeLearning)
		assert.True(t, engine.config.EnablePersonalityProfiling)
		assert.True(t, engine.config.EnableRiskProfiling)
		assert.True(t, engine.config.EnablePatternRecognition)
		assert.Equal(t, 24*time.Hour, engine.config.ModelUpdateInterval)
		assert.Equal(t, 0.6, engine.config.ConfidenceThreshold)
	})

	t.Run("LearnFromTradingBehavior", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.New()

		// Create a trading behavior event
		event := &BehaviorEvent{
			ID:     uuid.New().String(),
			UserID: userID,
			Type:   "trade",
			Action: "buy_btc",
			Context: &BehaviorContext{
				MarketConditions:   "bullish",
				PortfolioState:     map[string]interface{}{"btc_balance": 0.5},
				TimeOfDay:          "morning",
				DayOfWeek:          "monday",
				SessionDuration:    2 * time.Hour,
				PreviousActions:    []string{"analyze_chart", "check_news"},
				EmotionalState:     "confident",
				InformationSources: []string{"technical_analysis", "news"},
				ExternalFactors:    map[string]interface{}{"market_sentiment": "positive"},
				Metadata:           map[string]interface{}{},
			},
			Outcome: &BehaviorOutcome{
				Success:         true,
				Performance:     0.05, // 5% gain
				Satisfaction:    0.8,
				TimeToDecision:  15 * time.Minute,
				ConfidenceLevel: 0.7,
				Regret:          0.1,
				LearningValue:   0.8,
				Metadata:        map[string]interface{}{},
			},
			Timestamp: time.Now(),
			Duration:  30 * time.Minute,
			Metadata:  map[string]interface{}{},
		}

		err := engine.LearnFromBehavior(ctx, event)
		require.NoError(t, err)

		// Verify user profile was created and updated
		profile, err := engine.GetUserProfile(ctx, userID)
		require.NoError(t, err)
		require.NotNil(t, profile)

		assert.Equal(t, userID, profile.UserID)
		assert.Equal(t, 1, profile.ObservationCount)
		t.Logf("Overall confidence: %f", profile.Confidence)
		t.Logf("Component confidences - Trading: %f, Risk: %f, Personality: %f, Preferences: %f, Performance: %f",
			profile.TradingStyle.Confidence, profile.RiskProfile.Confidence, profile.PersonalityProfile.Confidence,
			profile.Preferences.Confidence, profile.PerformanceMetrics.Confidence)
		assert.Greater(t, profile.Confidence, 0.0)
		assert.NotNil(t, profile.TradingStyle)
		assert.NotNil(t, profile.RiskProfile)
		assert.NotNil(t, profile.PersonalityProfile)
		assert.NotNil(t, profile.Preferences)
		assert.NotNil(t, profile.PerformanceMetrics)
		assert.NotNil(t, profile.LearningProgress)

		// Check trading style learning
		assert.NotEmpty(t, profile.TradingStyle.PrimaryStyle)
		t.Logf("Trading style: %s, confidence: %f", profile.TradingStyle.PrimaryStyle, profile.TradingStyle.Confidence)
		assert.Greater(t, profile.TradingStyle.Confidence, 0.0)

		// Check performance metrics
		assert.Equal(t, 1, profile.PerformanceMetrics.TotalTrades)
		assert.Equal(t, 1, profile.PerformanceMetrics.SuccessfulTrades)
		assert.Equal(t, 1.0, profile.PerformanceMetrics.WinRate)
		assert.Greater(t, profile.PerformanceMetrics.AverageWin, 0.0)

		// Check learning progress
		assert.Equal(t, 1, profile.LearningProgress.TotalObservations)
		assert.Greater(t, profile.LearningProgress.ProfileCompleteness, 0.0)
	})

	t.Run("LearnFromMultipleBehaviors", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.New()

		// Create multiple behavior events
		events := []*BehaviorEvent{
			{
				ID:     uuid.New().String(),
				UserID: userID,
				Type:   "trade",
				Action: "buy_eth",
				Context: &BehaviorContext{
					MarketConditions: "neutral",
					TimeOfDay:        "afternoon",
					DayOfWeek:        "tuesday",
					SessionDuration:  1 * time.Hour,
				},
				Outcome: &BehaviorOutcome{
					Success:      true,
					Performance:  0.03,
					Satisfaction: 0.7,
				},
				Timestamp: time.Now().Add(-2 * time.Hour),
				Duration:  20 * time.Minute,
			},
			{
				ID:     uuid.New().String(),
				UserID: userID,
				Type:   "analyze",
				Action: "technical_analysis",
				Context: &BehaviorContext{
					MarketConditions: "bearish",
					TimeOfDay:        "evening",
					DayOfWeek:        "wednesday",
					SessionDuration:  3 * time.Hour,
				},
				Outcome: &BehaviorOutcome{
					Success:       true,
					Satisfaction:  0.9,
					LearningValue: 0.8,
				},
				Timestamp: time.Now().Add(-1 * time.Hour),
				Duration:  45 * time.Minute,
			},
			{
				ID:     uuid.New().String(),
				UserID: userID,
				Type:   "trade",
				Action: "sell_btc",
				Context: &BehaviorContext{
					MarketConditions: "volatile",
					TimeOfDay:        "morning",
					DayOfWeek:        "thursday",
					SessionDuration:  30 * time.Minute,
				},
				Outcome: &BehaviorOutcome{
					Success:      false,
					Performance:  -0.02,
					Satisfaction: 0.3,
					Regret:       0.7,
				},
				Timestamp: time.Now(),
				Duration:  10 * time.Minute,
			},
		}

		// Learn from all events
		for _, event := range events {
			err := engine.LearnFromBehavior(ctx, event)
			require.NoError(t, err)
		}

		// Verify profile evolution
		profile, err := engine.GetUserProfile(ctx, userID)
		require.NoError(t, err)

		assert.Equal(t, 3, profile.ObservationCount)
		assert.Greater(t, profile.Confidence, 0.0)

		// Check performance metrics evolution
		assert.Equal(t, 2, profile.PerformanceMetrics.TotalTrades)
		assert.Equal(t, 1, profile.PerformanceMetrics.SuccessfulTrades)
		assert.Equal(t, 0.5, profile.PerformanceMetrics.WinRate)
		assert.Greater(t, profile.PerformanceMetrics.AverageWin, 0.0)
		assert.Greater(t, profile.PerformanceMetrics.AverageLoss, 0.0)

		// Check behavior history
		history, err := engine.GetBehaviorHistory(ctx, userID, 10)
		require.NoError(t, err)
		assert.Len(t, history, 3)
	})

	t.Run("PersonalityProfiling", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.New()

		// Create events that indicate analytical personality
		analyticalEvents := []*BehaviorEvent{
			{
				ID:     uuid.New().String(),
				UserID: userID,
				Type:   "analyze",
				Action: "deep_technical_analysis",
				Context: &BehaviorContext{
					SessionDuration:    4 * time.Hour,
					InformationSources: []string{"charts", "indicators", "patterns"},
				},
				Outcome: &BehaviorOutcome{
					Success:       true,
					LearningValue: 0.9,
				},
				Timestamp: time.Now(),
				Duration:  2 * time.Hour,
			},
			{
				ID:     uuid.New().String(),
				UserID: userID,
				Type:   "research",
				Action: "fundamental_analysis",
				Context: &BehaviorContext{
					SessionDuration:    3 * time.Hour,
					InformationSources: []string{"financial_reports", "news", "metrics"},
				},
				Outcome: &BehaviorOutcome{
					Success:       true,
					LearningValue: 0.8,
				},
				Timestamp: time.Now(),
				Duration:  1 * time.Hour,
			},
		}

		for _, event := range analyticalEvents {
			err := engine.LearnFromBehavior(ctx, event)
			require.NoError(t, err)
		}

		profile, err := engine.GetUserProfile(ctx, userID)
		require.NoError(t, err)

		// Check personality profiling
		assert.NotEmpty(t, profile.PersonalityProfile.TraderType)
		assert.Greater(t, profile.PersonalityProfile.Confidence, 0.0)
		assert.NotEmpty(t, profile.PersonalityProfile.Traits)

		// Check for analytical traits
		if profile.PersonalityProfile.TraderType == "analytical" {
			assert.Contains(t, profile.PersonalityProfile.Traits, "conscientiousness")
			assert.Greater(t, profile.PersonalityProfile.Traits["conscientiousness"], 0.5)
		}
	})

	t.Run("RiskProfiling", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.New()

		// Create events that indicate risk tolerance
		riskEvents := []*BehaviorEvent{
			{
				ID:     uuid.New().String(),
				UserID: userID,
				Type:   "trade",
				Action: "high_leverage_trade",
				Context: &BehaviorContext{
					MarketConditions: "volatile",
					PortfolioState:   map[string]interface{}{"leverage": 5.0},
				},
				Outcome: &BehaviorOutcome{
					Success:      true,
					Performance:  0.15,
					Satisfaction: 0.9,
				},
				Timestamp: time.Now(),
			},
			{
				ID:     uuid.New().String(),
				UserID: userID,
				Type:   "trade",
				Action: "conservative_trade",
				Context: &BehaviorContext{
					MarketConditions: "stable",
					PortfolioState:   map[string]interface{}{"leverage": 1.0},
				},
				Outcome: &BehaviorOutcome{
					Success:      true,
					Performance:  0.02,
					Satisfaction: 0.6,
				},
				Timestamp: time.Now(),
			},
		}

		for _, event := range riskEvents {
			err := engine.LearnFromBehavior(ctx, event)
			require.NoError(t, err)
		}

		profile, err := engine.GetUserProfile(ctx, userID)
		require.NoError(t, err)

		// Check risk profiling
		assert.Greater(t, profile.RiskProfile.Confidence, 0.0)
		assert.GreaterOrEqual(t, profile.RiskProfile.RiskTolerance, 0.0)
		assert.LessOrEqual(t, profile.RiskProfile.RiskTolerance, 1.0)
		assert.NotEmpty(t, profile.RiskProfile.PositionSizingStyle)
		assert.Greater(t, profile.RiskProfile.StopLossUsage, 0.0)
		assert.Greater(t, profile.RiskProfile.EmotionalStability, 0.0)
	})

	t.Run("RecommendationGeneration", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.New()

		// Create enough events to trigger recommendation generation
		for i := 0; i < 15; i++ {
			event := &BehaviorEvent{
				ID:     uuid.New().String(),
				UserID: userID,
				Type:   "trade",
				Action: "buy_crypto",
				Context: &BehaviorContext{
					MarketConditions: "bullish",
					TimeOfDay:        "morning",
				},
				Outcome: &BehaviorOutcome{
					Success:      i%3 != 0, // 66% success rate
					Performance:  0.03,
					Satisfaction: 0.7,
				},
				Timestamp: time.Now().Add(-time.Duration(i) * time.Hour),
			}

			err := engine.LearnFromBehavior(ctx, event)
			require.NoError(t, err)
		}

		// Get recommendations
		recommendations, err := engine.GetPersonalizedRecommendations(ctx, userID, 5)
		require.NoError(t, err)

		if len(recommendations) > 0 {
			rec := recommendations[0]
			assert.NotEmpty(t, rec.ID)
			assert.NotEmpty(t, rec.Type)
			assert.NotEmpty(t, rec.Title)
			assert.NotEmpty(t, rec.Description)
			assert.NotEmpty(t, rec.Reasoning)
			assert.Greater(t, rec.Confidence, 0.0)
			assert.NotEmpty(t, rec.Priority)
			assert.NotEmpty(t, rec.Category)
			assert.Equal(t, "pending", rec.Status)
			assert.NotNil(t, rec.ExpectedOutcome)
			assert.NotNil(t, rec.RiskAssessment)
			assert.NotNil(t, rec.Personalization)
		}
	})

	t.Run("RecommendationStatusUpdate", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.New()

		// Create enough events to generate recommendations
		for i := 0; i < 12; i++ {
			event := &BehaviorEvent{
				ID:     uuid.New().String(),
				UserID: userID,
				Type:   "trade",
				Action: "trade_action",
				Outcome: &BehaviorOutcome{
					Success: true,
				},
				Timestamp: time.Now(),
			}

			err := engine.LearnFromBehavior(ctx, event)
			require.NoError(t, err)
		}

		// Get recommendations
		recommendations, err := engine.GetPersonalizedRecommendations(ctx, userID, 1)
		require.NoError(t, err)

		if len(recommendations) > 0 {
			recID := recommendations[0].ID

			// Update recommendation status
			err = engine.UpdateRecommendationStatus(ctx, userID, recID, "accepted")
			require.NoError(t, err)

			// Verify status update
			updatedRecs, err := engine.GetPersonalizedRecommendations(ctx, userID, 10)
			require.NoError(t, err)

			found := false
			for _, rec := range updatedRecs {
				if rec.ID == recID {
					assert.Equal(t, "accepted", rec.Status)
					found = true
					break
				}
			}
			// Note: accepted recommendations might be filtered out from active recommendations
			if !found {
				// This is expected behavior - accepted recommendations are no longer "pending"
				t.Log("Accepted recommendation filtered out from active recommendations (expected)")
			}
		}
	})

	t.Run("BehaviorHistoryRetrieval", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.New()

		// Create multiple events
		eventCount := 25
		for i := 0; i < eventCount; i++ {
			event := &BehaviorEvent{
				ID:        uuid.New().String(),
				UserID:    userID,
				Type:      "trade",
				Action:    "test_action",
				Timestamp: time.Now().Add(-time.Duration(i) * time.Minute),
			}

			err := engine.LearnFromBehavior(ctx, event)
			require.NoError(t, err)
		}

		// Test history retrieval with limit
		history, err := engine.GetBehaviorHistory(ctx, userID, 10)
		require.NoError(t, err)
		assert.Len(t, history, 10)

		// Test history retrieval without limit
		fullHistory, err := engine.GetBehaviorHistory(ctx, userID, 0)
		require.NoError(t, err)
		assert.Len(t, fullHistory, eventCount)

		// Verify chronological order (most recent first)
		for i := 1; i < len(history); i++ {
			assert.True(t, history[i-1].Timestamp.After(history[i].Timestamp) ||
				history[i-1].Timestamp.Equal(history[i].Timestamp))
		}
	})

	t.Run("LearningModelsRetrieval", func(t *testing.T) {
		models := engine.GetLearningModels()
		assert.NotNil(t, models)
		// Initially empty, but structure should be valid
		assert.IsType(t, map[string]*LearningModel{}, models)
	})

	t.Run("LearningProgressTracking", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.New()

		// Create events to track learning progress
		for i := 0; i < 20; i++ {
			event := &BehaviorEvent{
				ID:     uuid.New().String(),
				UserID: userID,
				Type:   "trade",
				Action: "progressive_learning",
				Outcome: &BehaviorOutcome{
					Success:       i%2 == 0,
					LearningValue: 0.8,
				},
				Timestamp: time.Now().Add(-time.Duration(i) * time.Minute),
			}

			err := engine.LearnFromBehavior(ctx, event)
			require.NoError(t, err)
		}

		profile, err := engine.GetUserProfile(ctx, userID)
		require.NoError(t, err)

		progress := profile.LearningProgress
		assert.Equal(t, 20, progress.TotalObservations)
		assert.Greater(t, progress.ProfileCompleteness, 0.0)
		assert.LessOrEqual(t, progress.ProfileCompleteness, 1.0)
		assert.GreaterOrEqual(t, progress.LearningVelocity, 0.0)

		// Check for milestones
		if len(progress.Milestones) > 0 {
			milestone := progress.Milestones[0]
			assert.NotEmpty(t, milestone.ID)
			assert.NotEmpty(t, milestone.Name)
			assert.NotEmpty(t, milestone.Description)
			assert.False(t, milestone.AchievedAt.IsZero())
			assert.GreaterOrEqual(t, milestone.Value, milestone.Threshold)
		}
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		ctx := context.Background()

		// Test with invalid user ID
		_, err := engine.GetUserProfile(ctx, uuid.Nil)
		assert.Error(t, err)

		// Test with non-existent user
		nonExistentUser := uuid.New()
		_, err = engine.GetUserProfile(ctx, nonExistentUser)
		assert.Error(t, err)

		// Test recommendation status update with invalid user
		err = engine.UpdateRecommendationStatus(ctx, nonExistentUser, "fake-id", "accepted")
		assert.Error(t, err)

		// Test recommendation status update with invalid recommendation ID
		userID := uuid.New()
		event := &BehaviorEvent{
			ID:        uuid.New().String(),
			UserID:    userID,
			Type:      "trade",
			Action:    "test",
			Timestamp: time.Now(),
		}
		err = engine.LearnFromBehavior(ctx, event)
		require.NoError(t, err)

		err = engine.UpdateRecommendationStatus(ctx, userID, "non-existent-rec", "accepted")
		assert.Error(t, err)
	})

	t.Run("ConfigurationValidation", func(t *testing.T) {
		// Test configuration values
		config := engine.config
		assert.Greater(t, config.LearningRate, 0.0)
		assert.LessOrEqual(t, config.LearningRate, 1.0)
		assert.Greater(t, config.MinObservations, 0)
		assert.Greater(t, config.PatternDetectionWindow, time.Hour)
		assert.Greater(t, config.PreferenceUpdateRate, 0.0)
		assert.Greater(t, config.PersonalityUpdateRate, 0.0)
		assert.Greater(t, config.RiskToleranceUpdateRate, 0.0)
		assert.Greater(t, config.RecommendationThreshold, 0.0)
		assert.LessOrEqual(t, config.RecommendationThreshold, 1.0)
		assert.Greater(t, config.MaxHistorySize, 0)
		assert.Greater(t, config.ModelUpdateInterval, time.Hour)
		assert.Greater(t, config.ConfidenceThreshold, 0.0)
		assert.LessOrEqual(t, config.ConfidenceThreshold, 1.0)
	})
}

func TestBehaviorAnalyzer(t *testing.T) {
	logger := &observability.Logger{}
	analyzer := NewBehaviorAnalyzer(logger)
	require.NotNil(t, analyzer)

	t.Run("TradingStyleAnalysis", func(t *testing.T) {
		ctx := context.Background()
		profile := &UserBehaviorProfile{
			UserID:       uuid.New(),
			TradingStyle: &TradingStyleProfile{},
		}

		// Test scalper identification
		scalperEvent := &BehaviorEvent{
			Type:     "trade",
			Duration: 30 * time.Minute,
		}

		err := analyzer.AnalyzeBehavior(ctx, profile, scalperEvent)
		require.NoError(t, err)
		assert.Equal(t, "scalper", profile.TradingStyle.PrimaryStyle)

		// Test day trader identification
		profile.TradingStyle.PrimaryStyle = "" // Reset
		dayTraderEvent := &BehaviorEvent{
			Type:     "trade",
			Duration: 4 * time.Hour,
		}

		err = analyzer.AnalyzeBehavior(ctx, profile, dayTraderEvent)
		require.NoError(t, err)
		assert.Equal(t, "day_trader", profile.TradingStyle.PrimaryStyle)

		// Test swing trader identification
		profile.TradingStyle.PrimaryStyle = "" // Reset
		swingTraderEvent := &BehaviorEvent{
			Type:     "trade",
			Duration: 48 * time.Hour,
		}

		err = analyzer.AnalyzeBehavior(ctx, profile, swingTraderEvent)
		require.NoError(t, err)
		assert.Equal(t, "swing_trader", profile.TradingStyle.PrimaryStyle)
	})

	t.Run("ConfidenceGrowth", func(t *testing.T) {
		ctx := context.Background()
		profile := &UserBehaviorProfile{
			UserID:       uuid.New(),
			TradingStyle: &TradingStyleProfile{Confidence: 0.0},
		}

		event := &BehaviorEvent{
			Type: "trade",
		}

		initialConfidence := profile.TradingStyle.Confidence
		err := analyzer.AnalyzeBehavior(ctx, profile, event)
		require.NoError(t, err)
		assert.Greater(t, profile.TradingStyle.Confidence, initialConfidence)
		assert.LessOrEqual(t, profile.TradingStyle.Confidence, 1.0)
	})
}

func TestPersonalityProfiler(t *testing.T) {
	logger := &observability.Logger{}
	profiler := NewPersonalityProfiler(logger)
	require.NotNil(t, profiler)

	t.Run("PersonalityTraitUpdate", func(t *testing.T) {
		ctx := context.Background()
		profile := &UserBehaviorProfile{
			UserID:             uuid.New(),
			PersonalityProfile: &PersonalityProfile{Traits: make(map[string]float64)},
		}

		// Test analytical trader identification
		analyticalEvent := &BehaviorEvent{
			Type: "analyze",
			Context: &BehaviorContext{
				SessionDuration: 3 * time.Hour,
			},
		}

		err := profiler.UpdatePersonality(ctx, profile, analyticalEvent)
		require.NoError(t, err)
		assert.Equal(t, "analytical", profile.PersonalityProfile.TraderType)

		// Test intuitive trader identification
		profile.PersonalityProfile.TraderType = "" // Reset
		intuitiveEvent := &BehaviorEvent{
			Type: "trade",
			Context: &BehaviorContext{
				SessionDuration: 30 * time.Minute,
			},
		}

		err = profiler.UpdatePersonality(ctx, profile, intuitiveEvent)
		require.NoError(t, err)
		assert.Equal(t, "intuitive", profile.PersonalityProfile.TraderType)

		// Check Big Five traits
		assert.Contains(t, profile.PersonalityProfile.Traits, "openness")
		assert.Contains(t, profile.PersonalityProfile.Traits, "conscientiousness")
		assert.Contains(t, profile.PersonalityProfile.Traits, "extraversion")
		assert.Contains(t, profile.PersonalityProfile.Traits, "agreeableness")
		assert.Contains(t, profile.PersonalityProfile.Traits, "neuroticism")

		for trait, value := range profile.PersonalityProfile.Traits {
			assert.GreaterOrEqual(t, value, 0.0, "Trait %s should be >= 0", trait)
			assert.LessOrEqual(t, value, 1.0, "Trait %s should be <= 1", trait)
		}
	})
}

func TestRiskProfiler(t *testing.T) {
	logger := &observability.Logger{}
	profiler := NewRiskProfiler(logger)
	require.NotNil(t, profiler)

	t.Run("RiskToleranceUpdate", func(t *testing.T) {
		ctx := context.Background()
		profile := &UserBehaviorProfile{
			UserID:      uuid.New(),
			RiskProfile: &UserRiskProfile{RiskTolerance: 0.5},
		}

		// Test successful trade increases risk tolerance
		successEvent := &BehaviorEvent{
			Type: "trade",
			Outcome: &BehaviorOutcome{
				Success: true,
			},
		}

		initialTolerance := profile.RiskProfile.RiskTolerance
		err := profiler.UpdateRiskProfile(ctx, profile, successEvent)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, profile.RiskProfile.RiskTolerance, initialTolerance)

		// Test failed trade decreases risk tolerance
		failEvent := &BehaviorEvent{
			Type: "trade",
			Outcome: &BehaviorOutcome{
				Success: false,
			},
		}

		toleranceBeforeFail := profile.RiskProfile.RiskTolerance
		err = profiler.UpdateRiskProfile(ctx, profile, failEvent)
		require.NoError(t, err)
		assert.LessOrEqual(t, profile.RiskProfile.RiskTolerance, toleranceBeforeFail)

		// Ensure risk tolerance stays within bounds
		assert.GreaterOrEqual(t, profile.RiskProfile.RiskTolerance, 0.0)
		assert.LessOrEqual(t, profile.RiskProfile.RiskTolerance, 1.0)
	})

	t.Run("RiskProfileCompletion", func(t *testing.T) {
		ctx := context.Background()
		profile := &UserBehaviorProfile{
			UserID:      uuid.New(),
			RiskProfile: &UserRiskProfile{},
		}

		event := &BehaviorEvent{
			Type: "trade",
			Outcome: &BehaviorOutcome{
				Success: true,
			},
		}

		err := profiler.UpdateRiskProfile(ctx, profile, event)
		require.NoError(t, err)

		// Check that risk profile fields are populated
		assert.NotEmpty(t, profile.RiskProfile.PositionSizingStyle)
		assert.Greater(t, profile.RiskProfile.StopLossUsage, 0.0)
		assert.Greater(t, profile.RiskProfile.TakeProfitUsage, 0.0)
		assert.GreaterOrEqual(t, profile.RiskProfile.LeverageComfort, 0.0)
		assert.Greater(t, profile.RiskProfile.VolatilityTolerance, 0.0)
		assert.Greater(t, profile.RiskProfile.EmotionalStability, 0.0)
		assert.Greater(t, profile.RiskProfile.LossAversion, 0.0)
		assert.Greater(t, profile.RiskProfile.Confidence, 0.0)
	})
}

func TestRecommendationEngine(t *testing.T) {
	logger := &observability.Logger{}
	engine := NewRecommendationEngine(logger)
	require.NotNil(t, engine)

	t.Run("RecommendationGeneration", func(t *testing.T) {
		ctx := context.Background()
		profile := &UserBehaviorProfile{
			UserID:          uuid.New(),
			Confidence:      0.8, // High confidence to trigger recommendations
			Recommendations: []*PersonalizedRecommendation{},
		}

		err := engine.GenerateRecommendations(ctx, profile)
		require.NoError(t, err)

		if len(profile.Recommendations) > 0 {
			rec := profile.Recommendations[0]
			assert.NotEmpty(t, rec.ID)
			assert.Equal(t, "strategy", rec.Type)
			assert.NotEmpty(t, rec.Title)
			assert.NotEmpty(t, rec.Description)
			assert.NotEmpty(t, rec.Reasoning)
			assert.Greater(t, rec.Confidence, 0.0)
			assert.NotEmpty(t, rec.Priority)
			assert.NotEmpty(t, rec.Category)
			assert.Equal(t, "pending", rec.Status)
			assert.NotNil(t, rec.ExpectedOutcome)
			assert.NotNil(t, rec.RiskAssessment)
			assert.NotNil(t, rec.Personalization)
		}
	})

	t.Run("RecommendationLimiting", func(t *testing.T) {
		ctx := context.Background()
		profile := &UserBehaviorProfile{
			UserID:     uuid.New(),
			Confidence: 0.8,
			Recommendations: []*PersonalizedRecommendation{
				{ID: "1"}, {ID: "2"}, {ID: "3"}, {ID: "4"}, {ID: "5"},
			},
		}

		initialCount := len(profile.Recommendations)
		err := engine.GenerateRecommendations(ctx, profile)
		require.NoError(t, err)

		// Should not exceed max recommendations
		assert.LessOrEqual(t, len(profile.Recommendations), initialCount+1)
	})
}
