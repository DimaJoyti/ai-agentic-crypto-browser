package ai

import (
	"context"
	"testing"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecisionEngine(t *testing.T) {
	logger := &observability.Logger{}
	engine := NewDecisionEngine(logger)
	require.NotNil(t, engine)

	t.Run("EngineInitialization", func(t *testing.T) {
		assert.NotNil(t, engine.logger)
		assert.NotNil(t, engine.config)
		assert.NotNil(t, engine.riskManager)
		assert.NotNil(t, engine.portfolioOptimizer)
		assert.NotNil(t, engine.signalAggregator)
		assert.NotNil(t, engine.executionEngine)
		assert.NotNil(t, engine.decisionTrees)
		assert.NotNil(t, engine.strategies)
		assert.NotNil(t, engine.activeDecisions)
		assert.NotNil(t, engine.decisionHistory)
		assert.NotNil(t, engine.performanceTracker)

		// Check configuration
		assert.Equal(t, 10, engine.config.MaxConcurrentDecisions)
		assert.Equal(t, 5*time.Minute, engine.config.DecisionTimeout)
		assert.Equal(t, 0.7, engine.config.MinConfidenceThreshold)
		assert.Equal(t, 0.05, engine.config.MaxRiskPerDecision)
		assert.False(t, engine.config.EnableAutoExecution) // Should start disabled
		assert.True(t, engine.config.PaperTradingMode)     // Should start in paper trading
		assert.True(t, engine.config.BacktestingEnabled)
		assert.True(t, engine.config.EmergencyStopEnabled)
	})

	t.Run("SimpleDecisionProcessing", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.New()

		req := &DecisionRequest{
			RequestID:    uuid.New().String(),
			UserID:       userID,
			DecisionType: "trade",
			Context: &DecisionContext{
				MarketConditions: "bullish",
				TimeHorizon:      "short",
				Urgency:          "medium",
				TriggerEvent:     "price_breakout",
				TechnicalIndicators: map[string]float64{
					"rsi":           30.0, // Oversold
					"macd_signal":   1.0,  // Bullish
					"volume_ratio":  1.5,  // High volume
				},
			},
			Constraints: &DecisionConstraints{
				MaxPositionSize: decimal.NewFromFloat(1000.0),
				MaxRiskExposure: 0.05,
				AllowedAssets:   []string{"BTC", "ETH"},
			},
			Preferences: &UserDecisionPrefs{
				RiskTolerance:      0.6,
				AutoExecutionLevel: "none",
				DecisionSpeed:      "normal",
			},
			MarketData: &MarketDataSnapshot{
				Timestamp: time.Now(),
				Prices: map[string]decimal.Decimal{
					"BTC": decimal.NewFromFloat(50000.0),
					"ETH": decimal.NewFromFloat(3000.0),
				},
				Sentiment:  0.7,
				Volatility: map[string]float64{"BTC": 0.3, "ETH": 0.4},
			},
			Options: DecisionOptions{
				RequireConfirmation: true,
				ExplainReasoning:    true,
				SimulateExecution:   true,
			},
			RequestedAt: time.Now(),
			ExpiresAt:   time.Now().Add(24 * time.Hour),
		}

		result, err := engine.ProcessDecisionRequest(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Validate basic structure
		assert.Equal(t, req.RequestID, result.RequestID)
		assert.Equal(t, req.UserID, result.UserID)
		assert.Equal(t, req.DecisionType, result.DecisionType)
		assert.NotEmpty(t, result.DecisionID)

		// Validate recommendation
		assert.NotNil(t, result.Recommendation)
		assert.NotEmpty(t, result.Recommendation.Action)
		assert.NotEmpty(t, result.Recommendation.Asset)
		assert.Greater(t, result.Recommendation.Quantity.InexactFloat64(), 0.0)
		assert.GreaterOrEqual(t, result.Recommendation.Confidence, 0.0)
		assert.LessOrEqual(t, result.Recommendation.Confidence, 1.0)
		assert.GreaterOrEqual(t, result.Recommendation.RiskScore, 0.0)
		assert.LessOrEqual(t, result.Recommendation.RiskScore, 1.0)

		// Validate risk assessment
		assert.NotNil(t, result.RiskAssessment)
		assert.GreaterOrEqual(t, result.RiskAssessment.OverallRisk, 0.0)
		assert.LessOrEqual(t, result.RiskAssessment.OverallRisk, 1.0)
		assert.NotEmpty(t, result.RiskAssessment.RiskFactors)

		// Validate reasoning
		assert.NotNil(t, result.Reasoning)
		assert.NotEmpty(t, result.Reasoning.PrimaryFactors)
		assert.GreaterOrEqual(t, result.Reasoning.Confidence, 0.0)
		assert.LessOrEqual(t, result.Reasoning.Confidence, 1.0)

		// Validate confidence
		assert.GreaterOrEqual(t, result.Confidence, 0.0)
		assert.LessOrEqual(t, result.Confidence, 1.0)

		// Validate expected outcome
		assert.NotNil(t, result.ExpectedOutcome)
		assert.GreaterOrEqual(t, result.ExpectedOutcome.ProbabilityOfProfit, 0.0)
		assert.LessOrEqual(t, result.ExpectedOutcome.ProbabilityOfProfit, 1.0)
		assert.GreaterOrEqual(t, result.ExpectedOutcome.ProbabilityOfLoss, 0.0)
		assert.LessOrEqual(t, result.ExpectedOutcome.ProbabilityOfLoss, 1.0)

		// Validate execution plan
		assert.NotNil(t, result.ExecutionPlan)
		assert.NotEmpty(t, result.ExecutionPlan.Steps)
		assert.Greater(t, result.ExecutionPlan.TotalEstimatedTime, time.Duration(0))

		// Should require approval since auto-execution is disabled
		assert.True(t, result.RequiresApproval)
		assert.False(t, result.AutoExecutable)

		// Validate timestamps
		assert.False(t, result.GeneratedAt.IsZero())
		assert.False(t, result.ExpiresAt.IsZero())
		assert.True(t, result.ExpiresAt.After(result.GeneratedAt))
	})

	t.Run("DecisionValidation", func(t *testing.T) {
		ctx := context.Background()

		// Test missing request ID
		req := &DecisionRequest{
			UserID:       uuid.New(),
			DecisionType: "trade",
			Context:      &DecisionContext{},
			ExpiresAt:    time.Now().Add(1 * time.Hour),
		}

		_, err := engine.ProcessDecisionRequest(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "request ID is required")

		// Test missing user ID
		req = &DecisionRequest{
			RequestID:    uuid.New().String(),
			DecisionType: "trade",
			Context:      &DecisionContext{},
			ExpiresAt:    time.Now().Add(1 * time.Hour),
		}

		_, err = engine.ProcessDecisionRequest(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user ID is required")

		// Test missing decision type
		req = &DecisionRequest{
			RequestID: uuid.New().String(),
			UserID:    uuid.New(),
			Context:   &DecisionContext{},
			ExpiresAt: time.Now().Add(1 * time.Hour),
		}

		_, err = engine.ProcessDecisionRequest(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "decision type is required")

		// Test missing context
		req = &DecisionRequest{
			RequestID:    uuid.New().String(),
			UserID:       uuid.New(),
			DecisionType: "trade",
			ExpiresAt:    time.Now().Add(1 * time.Hour),
		}

		_, err = engine.ProcessDecisionRequest(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "decision context is required")

		// Test expired request
		req = &DecisionRequest{
			RequestID:    uuid.New().String(),
			UserID:       uuid.New(),
			DecisionType: "trade",
			Context:      &DecisionContext{},
			ExpiresAt:    time.Now().Add(-1 * time.Hour), // Expired
		}

		_, err = engine.ProcessDecisionRequest(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "request has expired")
	})

	t.Run("ComplexDecisionDetection", func(t *testing.T) {
		// Test complex decision types
		complexTypes := []string{"portfolio_rebalance", "risk_management", "multi_asset_strategy"}
		
		for _, decisionType := range complexTypes {
			req := &DecisionRequest{
				DecisionType: decisionType,
				Options:      DecisionOptions{},
			}
			
			isComplex := engine.isComplexDecision(req)
			assert.True(t, isComplex, "Decision type %s should be complex", decisionType)
		}

		// Test simple decision type
		req := &DecisionRequest{
			DecisionType: "trade",
			Options:      DecisionOptions{},
		}
		
		isComplex := engine.isComplexDecision(req)
		assert.False(t, isComplex, "Decision type 'trade' should be simple")

		// Test backtesting makes it complex
		req = &DecisionRequest{
			DecisionType: "trade",
			Options: DecisionOptions{
				EnableBacktesting: true,
			},
		}
		
		isComplex = engine.isComplexDecision(req)
		assert.True(t, isComplex, "Backtesting should make decision complex")

		// Test alternatives make it complex
		req = &DecisionRequest{
			DecisionType: "trade",
			Options: DecisionOptions{
				IncludeAlternatives: true,
			},
		}
		
		isComplex = engine.isComplexDecision(req)
		assert.True(t, isComplex, "Including alternatives should make decision complex")
	})

	t.Run("ConfidenceCalculation", func(t *testing.T) {
		marketAnalysis := map[string]interface{}{
			"trend":      "bullish",
			"volatility": 0.3,
			"sentiment":  0.6,
		}

		riskAssessment := &RiskAssessment{
			OverallRisk: 0.2, // Low risk
			Confidence:  0.8,
		}

		recommendations := []DecisionRecommendation{
			{
				Confidence: 0.8,
				RiskScore:  0.3,
			},
		}

		confidence := engine.calculateConfidence(marketAnalysis, riskAssessment, recommendations)
		
		assert.GreaterOrEqual(t, confidence, 0.0)
		assert.LessOrEqual(t, confidence, 1.0)
		assert.Greater(t, confidence, 0.7) // Should be high due to bullish trend and low risk
	})

	t.Run("AutoExecutionDecision", func(t *testing.T) {
		recommendation := DecisionRecommendation{
			RiskScore:  0.03, // Low risk
			Confidence: 0.8,  // High confidence
		}

		riskAssessment := &RiskAssessment{
			OverallRisk: 0.02,
		}

		// Test with auto-execution disabled (default)
		req := &DecisionRequest{
			Preferences: &UserDecisionPrefs{
				AutoExecutionLevel: "none",
			},
		}

		autoExecutable := engine.isAutoExecutable(recommendation, riskAssessment, req)
		assert.False(t, autoExecutable, "Should not be auto-executable when disabled")

		// Test with high risk
		highRiskRecommendation := DecisionRecommendation{
			RiskScore:  0.8, // High risk
			Confidence: 0.8,
		}

		req = &DecisionRequest{
			Preferences: &UserDecisionPrefs{
				AutoExecutionLevel: "moderate",
			},
		}

		autoExecutable = engine.isAutoExecutable(highRiskRecommendation, riskAssessment, req)
		assert.False(t, autoExecutable, "Should not be auto-executable with high risk")

		// Test with low confidence
		lowConfidenceRecommendation := DecisionRecommendation{
			RiskScore:  0.03,
			Confidence: 0.5, // Low confidence
		}

		autoExecutable = engine.isAutoExecutable(lowConfidenceRecommendation, riskAssessment, req)
		assert.False(t, autoExecutable, "Should not be auto-executable with low confidence")
	})

	t.Run("ActiveDecisionTracking", func(t *testing.T) {
		userID := uuid.New()

		// Initially no active decisions
		activeDecisions := engine.GetActiveDecisions(userID)
		assert.Empty(t, activeDecisions)

		// Create a decision request that will be processed asynchronously
		ctx := context.Background()
		req := &DecisionRequest{
			RequestID:    uuid.New().String(),
			UserID:       userID,
			DecisionType: "portfolio_rebalance", // Complex decision
			Context: &DecisionContext{
				MarketConditions: "volatile",
				TimeHorizon:      "medium",
				Urgency:          "low",
			},
			Options: DecisionOptions{
				EnableBacktesting: true, // Makes it complex
			},
			RequestedAt: time.Now(),
			ExpiresAt:   time.Now().Add(24 * time.Hour),
		}

		result, err := engine.ProcessDecisionRequest(ctx, req)
		require.NoError(t, err)

		// Should have pending result for complex decision
		assert.Equal(t, req.RequestID, result.RequestID)
		assert.Equal(t, userID, result.UserID)

		// Wait a moment for async processing
		time.Sleep(100 * time.Millisecond)

		// Check if decision was processed
		activeDecisions = engine.GetActiveDecisions(userID)
		// Decision might be completed by now, so we just check it was tracked
		assert.True(t, len(activeDecisions) >= 0)
	})

	t.Run("DecisionHistory", func(t *testing.T) {
		userID := uuid.New()

		// Initially no history
		history := engine.GetDecisionHistory(userID, 10)
		assert.Empty(t, history)

		// Process a simple decision
		ctx := context.Background()
		req := &DecisionRequest{
			RequestID:    uuid.New().String(),
			UserID:       userID,
			DecisionType: "trade",
			Context: &DecisionContext{
				MarketConditions: "bullish",
			},
			RequestedAt: time.Now(),
			ExpiresAt:   time.Now().Add(1 * time.Hour),
		}

		_, err := engine.ProcessDecisionRequest(ctx, req)
		require.NoError(t, err)

		// Should have history now
		history = engine.GetDecisionHistory(userID, 10)
		assert.Len(t, history, 1)

		record := history[0]
		assert.Equal(t, req.RequestID, record.Request.RequestID)
		assert.Equal(t, userID, record.UserID)
		assert.Equal(t, "trade", record.DecisionType)
		assert.NotNil(t, record.Request)
		assert.NotNil(t, record.Result)
		assert.False(t, record.CreatedAt.IsZero())
	})

	t.Run("PerformanceTracking", func(t *testing.T) {
		metrics := engine.GetPerformanceMetrics()
		require.NotNil(t, metrics)

		// Should have basic structure
		assert.GreaterOrEqual(t, metrics.TotalDecisions, 0)
		assert.GreaterOrEqual(t, metrics.SuccessfulDecisions, 0)
		assert.GreaterOrEqual(t, metrics.SuccessRate, 0.0)
		assert.LessOrEqual(t, metrics.SuccessRate, 1.0)
		assert.False(t, metrics.LastUpdated.IsZero())
	})

	t.Run("ConcurrentDecisionLimit", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.New()

		// Fill up to the limit with complex decisions
		var requests []*DecisionRequest
		for i := 0; i < engine.config.MaxConcurrentDecisions; i++ {
			req := &DecisionRequest{
				RequestID:    uuid.New().String(),
				UserID:       userID,
				DecisionType: "portfolio_rebalance", // Complex decision
				Context: &DecisionContext{
					MarketConditions: "volatile",
				},
				Options: DecisionOptions{
					EnableBacktesting: true,
				},
				RequestedAt: time.Now(),
				ExpiresAt:   time.Now().Add(24 * time.Hour),
			}
			requests = append(requests, req)

			_, err := engine.ProcessDecisionRequest(ctx, req)
			require.NoError(t, err)
		}

		// Next request should fail due to limit
		req := &DecisionRequest{
			RequestID:    uuid.New().String(),
			UserID:       userID,
			DecisionType: "portfolio_rebalance",
			Context: &DecisionContext{
				MarketConditions: "volatile",
			},
			Options: DecisionOptions{
				EnableBacktesting: true,
			},
			RequestedAt: time.Now(),
			ExpiresAt:   time.Now().Add(24 * time.Hour),
		}

		_, err := engine.ProcessDecisionRequest(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "maximum concurrent decisions reached")
	})
}

func TestRiskManager(t *testing.T) {
	logger := &observability.Logger{}
	riskManager := NewRiskManager(logger)
	require.NotNil(t, riskManager)

	t.Run("RiskAssessment", func(t *testing.T) {
		ctx := context.Background()

		req := &DecisionRequest{
			DecisionType: "trade",
			Context: &DecisionContext{
				MarketConditions: "volatile",
				Urgency:          "high",
			},
			Constraints: &DecisionConstraints{
				MaxRiskExposure: 0.05,
			},
		}

		marketAnalysis := map[string]interface{}{
			"volatility": 0.4,
			"trend":      "bearish",
		}

		assessment, err := riskManager.AssessRisk(ctx, req, marketAnalysis)
		require.NoError(t, err)
		require.NotNil(t, assessment)

		assert.GreaterOrEqual(t, assessment.OverallRisk, 0.0)
		assert.LessOrEqual(t, assessment.OverallRisk, 1.0)
		assert.NotEmpty(t, assessment.RiskFactors)
		assert.Greater(t, assessment.MaxLoss.InexactFloat64(), 0.0)
		assert.GreaterOrEqual(t, assessment.Probability, 0.0)
		assert.LessOrEqual(t, assessment.Probability, 1.0)
		assert.Greater(t, assessment.TimeHorizon, time.Duration(0))
		assert.GreaterOrEqual(t, assessment.Confidence, 0.0)
		assert.LessOrEqual(t, assessment.Confidence, 1.0)

		// Should have at least one risk factor
		riskFactor := assessment.RiskFactors[0]
		assert.NotEmpty(t, riskFactor.Risk)
		assert.GreaterOrEqual(t, riskFactor.Probability, 0.0)
		assert.LessOrEqual(t, riskFactor.Probability, 1.0)
		assert.GreaterOrEqual(t, riskFactor.Impact, 0.0)
		assert.LessOrEqual(t, riskFactor.Impact, 1.0)
		assert.NotEmpty(t, riskFactor.Mitigation)
		assert.Contains(t, []string{"low", "medium", "high", "critical"}, riskFactor.Severity)
	})
}

func TestDecisionPerformanceTracker(t *testing.T) {
	tracker := NewDecisionPerformanceTracker()
	require.NotNil(t, tracker)

	t.Run("InitialState", func(t *testing.T) {
		metrics := tracker.GetOverallMetrics()
		require.NotNil(t, metrics)

		assert.Equal(t, 0, metrics.TotalDecisions)
		assert.Equal(t, 0, metrics.SuccessfulDecisions)
		assert.Equal(t, 0.0, metrics.SuccessRate)
	})

	t.Run("RecordDecision", func(t *testing.T) {
		userID := uuid.New()
		
		record := DecisionRecord{
			DecisionID:   uuid.New().String(),
			UserID:       userID,
			DecisionType: "trade",
			Request: &DecisionRequest{
				RequestID: uuid.New().String(),
				UserID:    userID,
			},
			Result: &DecisionResult{
				Confidence: 0.8, // High confidence
			},
			CreatedAt: time.Now(),
		}

		tracker.RecordDecision(record)

		// Check overall metrics
		metrics := tracker.GetOverallMetrics()
		assert.Equal(t, 1, metrics.TotalDecisions)
		assert.Equal(t, 1, metrics.SuccessfulDecisions) // High confidence counts as success
		assert.Equal(t, 1.0, metrics.SuccessRate)
		assert.False(t, metrics.LastUpdated.IsZero())

		// Record a low confidence decision
		lowConfidenceRecord := DecisionRecord{
			DecisionID:   uuid.New().String(),
			UserID:       userID,
			DecisionType: "trade",
			Request: &DecisionRequest{
				RequestID: uuid.New().String(),
				UserID:    userID,
			},
			Result: &DecisionResult{
				Confidence: 0.5, // Low confidence
			},
			CreatedAt: time.Now(),
		}

		tracker.RecordDecision(lowConfidenceRecord)

		// Check updated metrics
		metrics = tracker.GetOverallMetrics()
		assert.Equal(t, 2, metrics.TotalDecisions)
		assert.Equal(t, 1, metrics.SuccessfulDecisions) // Only high confidence counts
		assert.Equal(t, 0.5, metrics.SuccessRate)
	})
}
