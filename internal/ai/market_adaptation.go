package ai

import (
	"context"
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

// MarketAdaptationEngine provides adaptive learning from market patterns
type MarketAdaptationEngine struct {
	logger              *observability.Logger
	config              *MarketAdaptationConfig
	patternDetector     *PatternDetector
	strategyManager     *StrategyManager
	performanceAnalyzer *PerformanceAnalyzer
	adaptationRules     []*AdaptationRule
	detectedPatterns    []*DetectedPattern
	adaptiveStrategies  []*AdaptiveStrategy
	adaptationHistory   []*AdaptationRecord
	performanceMetrics  map[string]*MarketPerformanceMetrics
	mu                  sync.RWMutex
	lastUpdate          time.Time
}

// MarketAdaptationConfig holds configuration for market adaptation
type MarketAdaptationConfig struct {
	PatternDetectionWindow      time.Duration `json:"pattern_detection_window"`
	AdaptationThreshold         float64       `json:"adaptation_threshold"`
	MinPatternOccurrences       int           `json:"min_pattern_occurrences"`
	StrategyUpdateFrequency     time.Duration `json:"strategy_update_frequency"`
	PerformanceEvaluationWindow time.Duration `json:"performance_evaluation_window"`
	RiskAdjustmentSensitivity   float64       `json:"risk_adjustment_sensitivity"`
	AdaptationLearningRate      float64       `json:"adaptation_learning_rate"`
	MaxAdaptationHistory        int           `json:"max_adaptation_history"`
	EnableRealTimeAdaptation    bool          `json:"enable_real_time_adaptation"`
	ConfidenceThreshold         float64       `json:"confidence_threshold"`
}

// DetectedPattern represents a detected market pattern
type DetectedPattern struct {
	ID                string                 `json:"id"`
	Type              string                 `json:"type"` // trend, reversal, breakout, consolidation
	Name              string                 `json:"name"`
	Description       string                 `json:"description"`
	Asset             string                 `json:"asset"`
	TimeFrame         string                 `json:"timeframe"`
	Strength          float64                `json:"strength"`
	Confidence        float64                `json:"confidence"`
	Duration          time.Duration          `json:"duration"`
	Characteristics   map[string]float64     `json:"characteristics"`
	TriggerConditions []*TriggerCondition    `json:"trigger_conditions"`
	ExpectedOutcome   *ExpectedOutcome       `json:"expected_outcome"`
	MarketContext     *MarketContextInfo     `json:"market_context"`
	FirstDetected     time.Time              `json:"first_detected"`
	LastSeen          time.Time              `json:"last_seen"`
	OccurrenceCount   int                    `json:"occurrence_count"`
	SuccessRate       float64                `json:"success_rate"`
	AverageReturn     float64                `json:"average_return"`
	Metadata          map[string]interface{} `json:"metadata"`
}

// TriggerCondition represents a condition that triggers pattern recognition
type TriggerCondition struct {
	Type       string                 `json:"type"` // price, volume, indicator, time
	Indicator  string                 `json:"indicator"`
	Operator   string                 `json:"operator"` // gt, lt, eq, cross_above, cross_below
	Value      float64                `json:"value"`
	Timeframe  string                 `json:"timeframe"`
	Confidence float64                `json:"confidence"`
	Weight     float64                `json:"weight"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// ExpectedOutcome represents expected outcome of a pattern
type ExpectedOutcome struct {
	Direction     string        `json:"direction"` // up, down, sideways
	Magnitude     float64       `json:"magnitude"`
	Probability   float64       `json:"probability"`
	TimeHorizon   time.Duration `json:"time_horizon"`
	Confidence    float64       `json:"confidence"`
	RiskReward    float64       `json:"risk_reward"`
	StopLoss      float64       `json:"stop_loss"`
	TakeProfit    float64       `json:"take_profit"`
	SuccessRate   float64       `json:"success_rate"`
	AverageReturn float64       `json:"average_return"`
	MaxDrawdown   float64       `json:"max_drawdown"`
}

// MarketContextInfo represents market context information
type MarketContextInfo struct {
	MarketRegime        string                 `json:"market_regime"`     // bull, bear, sideways
	VolatilityRegime    string                 `json:"volatility_regime"` // low, medium, high
	TrendDirection      string                 `json:"trend_direction"`
	TrendStrength       float64                `json:"trend_strength"`
	MarketSentiment     float64                `json:"market_sentiment"`
	TechnicalIndicators map[string]float64     `json:"technical_indicators"`
	FundamentalFactors  map[string]float64     `json:"fundamental_factors"`
	Metadata            map[string]interface{} `json:"metadata"`
}

// AdaptiveStrategy represents a strategy that adapts to market conditions
type AdaptiveStrategy struct {
	ID                 string                      `json:"id"`
	Name               string                      `json:"name"`
	Description        string                      `json:"description"`
	Type               string                      `json:"type"` // trend_following, mean_reversion, momentum
	BaseParameters     map[string]float64          `json:"base_parameters"`
	CurrentParameters  map[string]float64          `json:"current_parameters"`
	AdaptationRules    []*AdaptationRule           `json:"adaptation_rules"`
	PerformanceTargets *PerformanceTargets         `json:"performance_targets"`
	RiskLimits         *MarketRiskLimits           `json:"risk_limits"`
	AdaptationHistory  []*MarketStrategyAdaptation `json:"adaptation_history"`
	PerformanceMetrics *MarketPerformanceMetrics   `json:"performance_metrics"`
	LastAdaptation     time.Time                   `json:"last_adaptation"`
	AdaptationCount    int                         `json:"adaptation_count"`
	IsActive           bool                        `json:"is_active"`
	Confidence         float64                     `json:"confidence"`
	Metadata           map[string]interface{}      `json:"metadata"`
}

// AdaptationRule represents a rule for strategy adaptation
type AdaptationRule struct {
	ID               string                 `json:"id"`
	Name             string                 `json:"name"`
	Description      string                 `json:"description"`
	TriggerType      string                 `json:"trigger_type"` // performance, pattern, market_condition
	TriggerCondition *TriggerCondition      `json:"trigger_condition"`
	AdaptationType   string                 `json:"adaptation_type"` // parameter_adjustment, strategy_switch
	AdaptationAction *AdaptationAction      `json:"adaptation_action"`
	Priority         int                    `json:"priority"`
	IsEnabled        bool                   `json:"is_enabled"`
	Confidence       float64                `json:"confidence"`
	Metadata         map[string]interface{} `json:"metadata"`
}

// AdaptationAction represents an action to take when adapting
type AdaptationAction struct {
	Type            string                 `json:"type"` // adjust_parameter, switch_strategy
	ParameterName   string                 `json:"parameter_name"`
	AdjustmentType  string                 `json:"adjustment_type"` // absolute, relative, multiplicative
	AdjustmentValue float64                `json:"adjustment_value"`
	MinValue        float64                `json:"min_value"`
	MaxValue        float64                `json:"max_value"`
	TargetStrategy  string                 `json:"target_strategy"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// PerformanceTargets represents performance targets for a strategy
type PerformanceTargets struct {
	TargetReturn       float64       `json:"target_return"`
	MaxDrawdown        float64       `json:"max_drawdown"`
	MinSharpeRatio     float64       `json:"min_sharpe_ratio"`
	MinWinRate         float64       `json:"min_win_rate"`
	MaxVolatility      float64       `json:"max_volatility"`
	TargetProfitFactor float64       `json:"target_profit_factor"`
	EvaluationPeriod   time.Duration `json:"evaluation_period"`
}

// MarketRiskLimits represents risk limits for a strategy
type MarketRiskLimits struct {
	MaxPositionSize    float64 `json:"max_position_size"`
	MaxLeverage        float64 `json:"max_leverage"`
	StopLossPercentage float64 `json:"stop_loss_percentage"`
	TakeProfitRatio    float64 `json:"take_profit_ratio"`
	MaxDailyLoss       float64 `json:"max_daily_loss"`
	VaRLimit           float64 `json:"var_limit"`
	ConcentrationLimit float64 `json:"concentration_limit"`
}

// MarketStrategyAdaptation represents a strategy adaptation event
type MarketStrategyAdaptation struct {
	ID                string                 `json:"id"`
	StrategyID        string                 `json:"strategy_id"`
	AdaptationType    string                 `json:"adaptation_type"`
	TriggerReason     string                 `json:"trigger_reason"`
	OldParameters     map[string]float64     `json:"old_parameters"`
	NewParameters     map[string]float64     `json:"new_parameters"`
	PerformanceBefore *PerformanceSnapshot   `json:"performance_before"`
	PerformanceAfter  *PerformanceSnapshot   `json:"performance_after"`
	MarketContext     *MarketContextInfo     `json:"market_context"`
	Confidence        float64                `json:"confidence"`
	Success           bool                   `json:"success"`
	Timestamp         time.Time              `json:"timestamp"`
	Metadata          map[string]interface{} `json:"metadata"`
}

// PerformanceSnapshot represents a performance snapshot
type PerformanceSnapshot struct {
	Return       float64   `json:"return"`
	Volatility   float64   `json:"volatility"`
	SharpeRatio  float64   `json:"sharpe_ratio"`
	MaxDrawdown  float64   `json:"max_drawdown"`
	WinRate      float64   `json:"win_rate"`
	ProfitFactor float64   `json:"profit_factor"`
	TotalTrades  int       `json:"total_trades"`
	Timestamp    time.Time `json:"timestamp"`
}

// PerformanceMetrics represents comprehensive performance metrics
type MarketPerformanceMetrics struct {
	StrategyID         string                 `json:"strategy_id"`
	TotalReturn        float64                `json:"total_return"`
	AnnualizedReturn   float64                `json:"annualized_return"`
	Volatility         float64                `json:"volatility"`
	SharpeRatio        float64                `json:"sharpe_ratio"`
	SortinoRatio       float64                `json:"sortino_ratio"`
	MaxDrawdown        float64                `json:"max_drawdown"`
	WinRate            float64                `json:"win_rate"`
	ProfitFactor       float64                `json:"profit_factor"`
	AverageWin         float64                `json:"average_win"`
	AverageLoss        float64                `json:"average_loss"`
	TotalTrades        int                    `json:"total_trades"`
	WinningTrades      int                    `json:"winning_trades"`
	LosingTrades       int                    `json:"losing_trades"`
	AverageHoldTime    time.Duration          `json:"average_hold_time"`
	BestTrade          float64                `json:"best_trade"`
	WorstTrade         float64                `json:"worst_trade"`
	ConsecutiveWins    int                    `json:"consecutive_wins"`
	ConsecutiveLosses  int                    `json:"consecutive_losses"`
	RiskAdjustedReturn float64                `json:"risk_adjusted_return"`
	AdaptationImpact   *AdaptationImpact      `json:"adaptation_impact"`
	LastUpdated        time.Time              `json:"last_updated"`
	Metadata           map[string]interface{} `json:"metadata"`
}

// AdaptationImpact represents the impact of adaptations on performance
type AdaptationImpact struct {
	TotalAdaptations        int                    `json:"total_adaptations"`
	SuccessfulAdaptations   int                    `json:"successful_adaptations"`
	AdaptationSuccessRate   float64                `json:"adaptation_success_rate"`
	PerformanceImprovement  float64                `json:"performance_improvement"`
	RiskReduction           float64                `json:"risk_reduction"`
	AdaptationFrequency     float64                `json:"adaptation_frequency"`
	AverageAdaptationImpact float64                `json:"average_adaptation_impact"`
	AdaptationsByType       map[string]int         `json:"adaptations_by_type"`
	AdaptationsByTrigger    map[string]int         `json:"adaptations_by_trigger"`
	Metadata                map[string]interface{} `json:"metadata"`
}

// AdaptationRecord represents an adaptation record
type AdaptationRecord struct {
	ID             string                 `json:"id"`
	Type           string                 `json:"type"` // pattern_detected, strategy_adapted, performance_evaluated
	Description    string                 `json:"description"`
	StrategyID     string                 `json:"strategy_id"`
	PatternID      string                 `json:"pattern_id"`
	AdaptationData map[string]interface{} `json:"adaptation_data"`
	Impact         *AdaptationImpact      `json:"impact"`
	Confidence     float64                `json:"confidence"`
	Timestamp      time.Time              `json:"timestamp"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// Supporting component types
type PatternDetector struct {
	config *PatternDetectorConfig
	logger *observability.Logger
}

type StrategyManager struct {
	config *StrategyManagerConfig
	logger *observability.Logger
}

type PerformanceAnalyzer struct {
	config *PerformanceAnalyzerConfig
	logger *observability.Logger
}

// Configuration types for components
type PatternDetectorConfig struct {
	MinPatternLength    int           `json:"min_pattern_length"`
	MaxPatternLength    int           `json:"max_pattern_length"`
	SimilarityThreshold float64       `json:"similarity_threshold"`
	ConfidenceThreshold float64       `json:"confidence_threshold"`
	UpdateFrequency     time.Duration `json:"update_frequency"`
}

type StrategyManagerConfig struct {
	AdaptationSensitivity float64       `json:"adaptation_sensitivity"`
	MinPerformancePeriod  time.Duration `json:"min_performance_period"`
	MaxAdaptationsPerDay  int           `json:"max_adaptations_per_day"`
	ConfidenceThreshold   float64       `json:"confidence_threshold"`
}

type PerformanceAnalyzerConfig struct {
	EvaluationFrequency time.Duration `json:"evaluation_frequency"`
	BenchmarkAsset      string        `json:"benchmark_asset"`
	RiskFreeRate        float64       `json:"risk_free_rate"`
	ConfidenceLevel     float64       `json:"confidence_level"`
}

// NewMarketAdaptationEngine creates a new market adaptation engine
func NewMarketAdaptationEngine(logger *observability.Logger) *MarketAdaptationEngine {
	config := &MarketAdaptationConfig{
		PatternDetectionWindow:      7 * 24 * time.Hour, // 7 days
		AdaptationThreshold:         0.7,
		MinPatternOccurrences:       3,
		StrategyUpdateFrequency:     1 * time.Hour,
		PerformanceEvaluationWindow: 24 * time.Hour,
		RiskAdjustmentSensitivity:   0.5,
		AdaptationLearningRate:      0.1,
		MaxAdaptationHistory:        1000,
		EnableRealTimeAdaptation:    true,
		ConfidenceThreshold:         0.6,
	}

	engine := &MarketAdaptationEngine{
		logger:              logger,
		config:              config,
		patternDetector:     NewPatternDetector(logger),
		strategyManager:     NewStrategyManager(logger),
		performanceAnalyzer: NewPerformanceAnalyzer(logger),
		adaptationRules:     []*AdaptationRule{},
		detectedPatterns:    []*DetectedPattern{},
		adaptiveStrategies:  []*AdaptiveStrategy{},
		adaptationHistory:   []*AdaptationRecord{},
		performanceMetrics:  make(map[string]*MarketPerformanceMetrics),
		lastUpdate:          time.Now(),
	}

	logger.Info(context.Background(), "Market adaptation engine initialized", map[string]interface{}{
		"pattern_detection_window":      config.PatternDetectionWindow.String(),
		"adaptation_threshold":          config.AdaptationThreshold,
		"min_pattern_occurrences":       config.MinPatternOccurrences,
		"strategy_update_frequency":     config.StrategyUpdateFrequency.String(),
		"performance_evaluation_window": config.PerformanceEvaluationWindow.String(),
		"real_time_adaptation":          config.EnableRealTimeAdaptation,
	})

	return engine
}

// DetectPatterns detects patterns in market data
func (m *MarketAdaptationEngine) DetectPatterns(ctx context.Context, marketData map[string]interface{}) ([]*DetectedPattern, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logger.Info(ctx, "Detecting market patterns", map[string]interface{}{
		"data_points": len(marketData),
	})

	// Use pattern detector to identify patterns
	patterns, err := m.patternDetector.DetectPatterns(ctx, marketData)
	if err != nil {
		return nil, fmt.Errorf("failed to detect patterns: %w", err)
	}

	// Update pattern database
	for _, pattern := range patterns {
		// Check if pattern already exists
		found := false
		for i, existingPattern := range m.detectedPatterns {
			if existingPattern.Type == pattern.Type && existingPattern.Asset == pattern.Asset {
				// Update existing pattern
				m.detectedPatterns[i].LastSeen = time.Now()
				m.detectedPatterns[i].OccurrenceCount++
				m.detectedPatterns[i].Confidence = (existingPattern.Confidence + pattern.Confidence) / 2
				found = true
				break
			}
		}

		if !found {
			// Add new pattern
			pattern.FirstDetected = time.Now()
			pattern.LastSeen = time.Now()
			pattern.OccurrenceCount = 1
			m.detectedPatterns = append(m.detectedPatterns, pattern)
		}
	}

	m.logger.Info(ctx, "Market pattern detection completed", map[string]interface{}{
		"patterns_detected": len(patterns),
		"total_patterns":    len(m.detectedPatterns),
	})

	return patterns, nil
}

// AdaptStrategies adapts trading strategies based on detected patterns
func (m *MarketAdaptationEngine) AdaptStrategies(ctx context.Context, patterns []*DetectedPattern) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logger.Info(ctx, "Adapting strategies", map[string]interface{}{
		"patterns":   len(patterns),
		"strategies": len(m.adaptiveStrategies),
	})

	adaptationCount := 0

	for _, strategy := range m.adaptiveStrategies {
		if !strategy.IsActive {
			continue
		}

		// Check if strategy needs adaptation
		needsAdaptation, reason := m.evaluateAdaptationNeed(ctx, strategy, patterns)
		if !needsAdaptation {
			continue
		}

		// Perform adaptation
		adaptation, err := m.strategyManager.AdaptStrategy(ctx, strategy, patterns, reason)
		if err != nil {
			m.logger.Warn(ctx, "Failed to adapt strategy", map[string]interface{}{
				"strategy_id": strategy.ID,
				"error":       err.Error(),
			})
			continue
		}

		// Apply adaptation
		if err := m.applyAdaptation(ctx, strategy, adaptation); err != nil {
			m.logger.Warn(ctx, "Failed to apply adaptation", map[string]interface{}{
				"strategy_id":   strategy.ID,
				"adaptation_id": adaptation.ID,
				"error":         err.Error(),
			})
			continue
		}

		// Record adaptation event
		record := &AdaptationRecord{
			ID:             uuid.New().String(),
			Type:           "strategy_adapted",
			Description:    fmt.Sprintf("Strategy %s adapted due to %s", strategy.Name, reason),
			StrategyID:     strategy.ID,
			AdaptationData: map[string]interface{}{"adaptation": adaptation},
			Confidence:     adaptation.Confidence,
			Timestamp:      time.Now(),
			Metadata:       map[string]interface{}{},
		}

		m.adaptationHistory = append(m.adaptationHistory, record)
		adaptationCount++

		// Maintain history size
		if len(m.adaptationHistory) > m.config.MaxAdaptationHistory {
			m.adaptationHistory = m.adaptationHistory[len(m.adaptationHistory)-m.config.MaxAdaptationHistory:]
		}

		strategy.LastAdaptation = time.Now()
		strategy.AdaptationCount++

		m.logger.Info(ctx, "Strategy adapted successfully", map[string]interface{}{
			"strategy_id":   strategy.ID,
			"adaptation_id": adaptation.ID,
			"reason":        reason,
			"confidence":    adaptation.Confidence,
		})
	}

	m.logger.Info(ctx, "Strategy adaptation completed", map[string]interface{}{
		"adaptations_made": adaptationCount,
		"total_history":    len(m.adaptationHistory),
	})

	return nil
}

// Helper methods

func (m *MarketAdaptationEngine) evaluateAdaptationNeed(ctx context.Context, strategy *AdaptiveStrategy, patterns []*DetectedPattern) (bool, string) {
	// Check performance-based adaptation needs
	if strategy.PerformanceMetrics != nil {
		if strategy.PerformanceTargets != nil {
			if strategy.PerformanceMetrics.SharpeRatio < strategy.PerformanceTargets.MinSharpeRatio {
				return true, "poor_sharpe_ratio"
			}
			if strategy.PerformanceMetrics.MaxDrawdown > strategy.PerformanceTargets.MaxDrawdown {
				return true, "excessive_drawdown"
			}
			if strategy.PerformanceMetrics.WinRate < strategy.PerformanceTargets.MinWinRate {
				return true, "low_win_rate"
			}
		}
	}

	// Check pattern-based adaptation needs
	for _, pattern := range patterns {
		if pattern.Confidence > m.config.AdaptationThreshold {
			return true, fmt.Sprintf("pattern_detected_%s", pattern.Type)
		}
	}

	// Check time-based adaptation needs
	if time.Since(strategy.LastAdaptation) > m.config.StrategyUpdateFrequency {
		return true, "scheduled_update"
	}

	return false, ""
}

func (m *MarketAdaptationEngine) applyAdaptation(ctx context.Context, strategy *AdaptiveStrategy, adaptation *MarketStrategyAdaptation) error {
	// Store old parameters
	oldParams := make(map[string]float64)
	for k, v := range strategy.CurrentParameters {
		oldParams[k] = v
	}

	// Apply new parameters
	for k, v := range adaptation.NewParameters {
		strategy.CurrentParameters[k] = v
	}

	// Update adaptation history
	strategy.AdaptationHistory = append(strategy.AdaptationHistory, adaptation)

	// Maintain adaptation history size
	if len(strategy.AdaptationHistory) > 100 {
		strategy.AdaptationHistory = strategy.AdaptationHistory[len(strategy.AdaptationHistory)-100:]
	}

	return nil
}

// Component implementations

func NewPatternDetector(logger *observability.Logger) *PatternDetector {
	return &PatternDetector{
		config: &PatternDetectorConfig{
			MinPatternLength:    5,
			MaxPatternLength:    50,
			SimilarityThreshold: 0.8,
			ConfidenceThreshold: 0.7,
			UpdateFrequency:     1 * time.Hour,
		},
		logger: logger,
	}
}

func (pd *PatternDetector) DetectPatterns(ctx context.Context, marketData map[string]interface{}) ([]*DetectedPattern, error) {
	// Simplified pattern detection
	patterns := []*DetectedPattern{}

	// Simulate trend pattern detection
	if priceData, ok := marketData["prices"].([]float64); ok && len(priceData) >= 10 {
		trendPattern := &DetectedPattern{
			ID:          uuid.New().String(),
			Type:        "trend",
			Name:        "Upward Trend",
			Description: "Detected upward price trend",
			Asset:       "BTC",
			TimeFrame:   "1h",
			Strength:    0.8,
			Confidence:  0.75,
			Duration:    4 * time.Hour,
			Characteristics: map[string]float64{
				"slope":          0.05,
				"r_squared":      0.85,
				"trend_strength": 0.8,
				"momentum":       0.7,
			},
			TriggerConditions: []*TriggerCondition{
				{
					Type:       "price",
					Indicator:  "sma_cross",
					Operator:   "cross_above",
					Value:      50000.0,
					Confidence: 0.8,
				},
			},
			ExpectedOutcome: &ExpectedOutcome{
				Direction:     "up",
				Magnitude:     0.1,
				Probability:   0.7,
				TimeHorizon:   24 * time.Hour,
				Confidence:    0.75,
				RiskReward:    2.5,
				StopLoss:      0.02,
				TakeProfit:    0.05,
				SuccessRate:   0.65,
				AverageReturn: 0.08,
				MaxDrawdown:   0.03,
			},
			MarketContext: &MarketContextInfo{
				MarketRegime:    "bull",
				TrendDirection:  "up",
				TrendStrength:   0.8,
				MarketSentiment: 0.7,
				TechnicalIndicators: map[string]float64{
					"rsi":  65.0,
					"macd": 0.5,
					"sma":  50000.0,
				},
			},
			Metadata: map[string]interface{}{},
		}
		patterns = append(patterns, trendPattern)
	}

	return patterns, nil
}

func NewStrategyManager(logger *observability.Logger) *StrategyManager {
	return &StrategyManager{
		config: &StrategyManagerConfig{
			AdaptationSensitivity: 0.5,
			MinPerformancePeriod:  24 * time.Hour,
			MaxAdaptationsPerDay:  5,
			ConfidenceThreshold:   0.6,
		},
		logger: logger,
	}
}

func (sm *StrategyManager) AdaptStrategy(ctx context.Context, strategy *AdaptiveStrategy, patterns []*DetectedPattern, reason string) (*MarketStrategyAdaptation, error) {
	// Simplified strategy adaptation
	adaptation := &MarketStrategyAdaptation{
		ID:             uuid.New().String(),
		StrategyID:     strategy.ID,
		AdaptationType: "parameter_adjustment",
		TriggerReason:  reason,
		OldParameters:  make(map[string]float64),
		NewParameters:  make(map[string]float64),
		Confidence:     0.75,
		Success:        true,
		Timestamp:      time.Now(),
		Metadata:       map[string]interface{}{},
	}

	// Copy old parameters
	for k, v := range strategy.CurrentParameters {
		adaptation.OldParameters[k] = v
		adaptation.NewParameters[k] = v
	}

	// Adapt based on reason
	switch reason {
	case "poor_sharpe_ratio":
		// Reduce position size to improve risk-adjusted returns
		if posSize, exists := adaptation.NewParameters["position_size"]; exists {
			adaptation.NewParameters["position_size"] = posSize * 0.8
		}
	case "excessive_drawdown":
		// Tighten stop losses
		if stopLoss, exists := adaptation.NewParameters["stop_loss"]; exists {
			adaptation.NewParameters["stop_loss"] = math.Max(stopLoss*0.8, 0.01)
		}
	case "low_win_rate":
		// Adjust entry conditions to be more selective
		if entryThreshold, exists := adaptation.NewParameters["entry_threshold"]; exists {
			adaptation.NewParameters["entry_threshold"] = entryThreshold * 1.2
		}
	default:
		// Default adaptation based on patterns
		if len(patterns) > 0 {
			pattern := patterns[0]
			if pattern.Type == "trend" && pattern.ExpectedOutcome != nil && pattern.ExpectedOutcome.Direction == "up" {
				// Increase position size for upward trends
				if posSize, exists := adaptation.NewParameters["position_size"]; exists {
					adaptation.NewParameters["position_size"] = math.Min(posSize*1.1, 0.1) // Cap at 10%
				}
			}
		}
	}

	return adaptation, nil
}

func NewPerformanceAnalyzer(logger *observability.Logger) *PerformanceAnalyzer {
	return &PerformanceAnalyzer{
		config: &PerformanceAnalyzerConfig{
			EvaluationFrequency: 1 * time.Hour,
			BenchmarkAsset:      "BTC",
			RiskFreeRate:        0.02,
			ConfidenceLevel:     0.95,
		},
		logger: logger,
	}
}

// Public API methods

// GetDetectedPatterns retrieves detected market patterns
func (m *MarketAdaptationEngine) GetDetectedPatterns(ctx context.Context, filters map[string]interface{}) ([]*DetectedPattern, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	patterns := []*DetectedPattern{}
	for _, pattern := range m.detectedPatterns {
		// Apply filters if provided
		if assetFilter, ok := filters["asset"].(string); ok {
			if pattern.Asset != assetFilter {
				continue
			}
		}
		if typeFilter, ok := filters["type"].(string); ok {
			if pattern.Type != typeFilter {
				continue
			}
		}
		if minConfidence, ok := filters["min_confidence"].(float64); ok {
			if pattern.Confidence < minConfidence {
				continue
			}
		}

		patterns = append(patterns, pattern)
	}

	// Sort by confidence descending
	sort.Slice(patterns, func(i, j int) bool {
		return patterns[i].Confidence > patterns[j].Confidence
	})

	return patterns, nil
}

// GetAdaptiveStrategies retrieves adaptive strategies
func (m *MarketAdaptationEngine) GetAdaptiveStrategies(ctx context.Context) ([]*AdaptiveStrategy, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.adaptiveStrategies, nil
}

// GetAdaptationHistory retrieves adaptation history
func (m *MarketAdaptationEngine) GetAdaptationHistory(ctx context.Context, limit int) ([]*AdaptationRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	history := m.adaptationHistory
	if limit > 0 && len(history) > limit {
		history = history[len(history)-limit:]
	}

	return history, nil
}

// GetPerformanceMetrics retrieves performance metrics for strategies
func (m *MarketAdaptationEngine) GetPerformanceMetrics(ctx context.Context, strategyID string) (*MarketPerformanceMetrics, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if strategyID == "" {
		return nil, fmt.Errorf("strategy ID is required")
	}

	metrics, exists := m.performanceMetrics[strategyID]
	if !exists {
		return nil, fmt.Errorf("performance metrics not found for strategy %s", strategyID)
	}

	return metrics, nil
}

// AddAdaptiveStrategy adds a new adaptive strategy
func (m *MarketAdaptationEngine) AddAdaptiveStrategy(ctx context.Context, strategy *AdaptiveStrategy) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if strategy.ID == "" {
		strategy.ID = uuid.New().String()
	}

	strategy.LastAdaptation = time.Now()
	strategy.AdaptationCount = 0
	strategy.IsActive = true

	m.adaptiveStrategies = append(m.adaptiveStrategies, strategy)

	m.logger.Info(ctx, "Adaptive strategy added", map[string]interface{}{
		"strategy_id":   strategy.ID,
		"strategy_name": strategy.Name,
		"strategy_type": strategy.Type,
	})

	return nil
}

// UpdateStrategyStatus updates the status of an adaptive strategy
func (m *MarketAdaptationEngine) UpdateStrategyStatus(ctx context.Context, strategyID string, isActive bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, strategy := range m.adaptiveStrategies {
		if strategy.ID == strategyID {
			strategy.IsActive = isActive

			m.logger.Info(ctx, "Strategy status updated", map[string]interface{}{
				"strategy_id": strategyID,
				"is_active":   isActive,
			})

			return nil
		}
	}

	return fmt.Errorf("strategy not found: %s", strategyID)
}

// SetPerformanceMetrics sets performance metrics for a strategy (for demo purposes)
func (m *MarketAdaptationEngine) SetPerformanceMetrics(strategyID string, metrics *MarketPerformanceMetrics) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.performanceMetrics[strategyID] = metrics
}
