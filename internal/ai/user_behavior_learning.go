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

// UserBehaviorLearningEngine provides advanced user behavior learning capabilities
type UserBehaviorLearningEngine struct {
	logger               *observability.Logger
	config               *UserBehaviorConfig
	behaviorAnalyzer     *BehaviorAnalyzer
	patternRecognizer    *PatternRecognizer
	preferenceEngine     *PreferenceEngine
	recommendationEngine *RecommendationEngine
	personalityProfiler  *PersonalityProfiler
	riskProfiler         *RiskProfiler
	userProfiles         map[uuid.UUID]*UserBehaviorProfile
	behaviorHistory      map[uuid.UUID][]*BehaviorEvent
	learningModels       map[string]*LearningModel
	mu                   sync.RWMutex
	lastUpdate           time.Time
}

// UserBehaviorConfig holds configuration for user behavior learning
type UserBehaviorConfig struct {
	LearningRate               float64       `json:"learning_rate"`
	MinObservations            int           `json:"min_observations"`
	PatternDetectionWindow     time.Duration `json:"pattern_detection_window"`
	PreferenceUpdateRate       float64       `json:"preference_update_rate"`
	PersonalityUpdateRate      float64       `json:"personality_update_rate"`
	RiskToleranceUpdateRate    float64       `json:"risk_tolerance_update_rate"`
	RecommendationThreshold    float64       `json:"recommendation_threshold"`
	MaxHistorySize             int           `json:"max_history_size"`
	EnableRealTimeLearning     bool          `json:"enable_real_time_learning"`
	EnablePersonalityProfiling bool          `json:"enable_personality_profiling"`
	EnableRiskProfiling        bool          `json:"enable_risk_profiling"`
	EnablePatternRecognition   bool          `json:"enable_pattern_recognition"`
	ModelUpdateInterval        time.Duration `json:"model_update_interval"`
	ConfidenceThreshold        float64       `json:"confidence_threshold"`
}

// UserBehaviorProfile represents a comprehensive user behavior profile
type UserBehaviorProfile struct {
	UserID             uuid.UUID                     `json:"user_id"`
	CreatedAt          time.Time                     `json:"created_at"`
	LastUpdated        time.Time                     `json:"last_updated"`
	TradingStyle       *TradingStyleProfile          `json:"trading_style"`
	RiskProfile        *UserRiskProfile              `json:"risk_profile"`
	PersonalityProfile *PersonalityProfile           `json:"personality_profile"`
	Preferences        *UserBehaviorPreferences      `json:"preferences"`
	BehaviorPatterns   []*BehaviorPattern            `json:"behavior_patterns"`
	PerformanceMetrics *UserPerformanceProfile       `json:"performance_metrics"`
	LearningProgress   *LearningProgress             `json:"learning_progress"`
	Recommendations    []*PersonalizedRecommendation `json:"recommendations"`
	Confidence         float64                       `json:"confidence"`
	ObservationCount   int                           `json:"observation_count"`
	Metadata           map[string]interface{}        `json:"metadata"`
}

// TradingStyleProfile represents user's trading style characteristics
type TradingStyleProfile struct {
	PrimaryStyle        string            `json:"primary_style"` // scalper, day_trader, swing_trader, position_trader
	SecondaryStyles     []string          `json:"secondary_styles"`
	TradingFrequency    float64           `json:"trading_frequency"` // trades per day
	AverageHoldTime     time.Duration     `json:"average_hold_time"`
	PreferredTimeframes []string          `json:"preferred_timeframes"`
	PreferredAssets     []string          `json:"preferred_assets"`
	TradingHours        *UserTradingHours `json:"trading_hours"`
	DecisionSpeed       string            `json:"decision_speed"` // fast, medium, slow, deliberate
	AnalysisDepth       string            `json:"analysis_depth"` // surface, moderate, deep, exhaustive
	InformationSources  []string          `json:"information_sources"`
	TechnicalFocus      []string          `json:"technical_focus"` // indicators, patterns, volume, etc.
	FundamentalFocus    []string          `json:"fundamental_focus"`
	Confidence          float64           `json:"confidence"`
}

// UserTradingHours represents user's preferred trading hours
type UserTradingHours struct {
	Timezone       string                 `json:"timezone"`
	WeekdayHours   *TimeRange             `json:"weekday_hours"`
	WeekendHours   *TimeRange             `json:"weekend_hours"`
	PeakHours      []*TimeRange           `json:"peak_hours"`
	AvoidanceHours []*TimeRange           `json:"avoidance_hours"`
	SessionPrefs   map[string]float64     `json:"session_preferences"` // asian, european, american
	Metadata       map[string]interface{} `json:"metadata"`
}

// TimeRange represents a time range
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// UserRiskProfile represents user's risk characteristics
type UserRiskProfile struct {
	RiskTolerance        float64 `json:"risk_tolerance"` // 0-1 scale
	RiskCapacity         float64 `json:"risk_capacity"`  // financial capacity
	MaxDrawdownTolerance float64 `json:"max_drawdown_tolerance"`
	PositionSizingStyle  string  `json:"position_sizing_style"` // fixed, percentage, kelly, volatility
	StopLossUsage        float64 `json:"stop_loss_usage"`       // frequency of use
	TakeProfitUsage      float64 `json:"take_profit_usage"`
	LeverageComfort      float64 `json:"leverage_comfort"`
	VolatilityTolerance  float64 `json:"volatility_tolerance"`
	CorrelationAwareness float64 `json:"correlation_awareness"`
	DiversificationLevel float64 `json:"diversification_level"`
	EmotionalStability   float64 `json:"emotional_stability"`
	LossAversion         float64 `json:"loss_aversion"`
	RiskAdjustmentSpeed  float64 `json:"risk_adjustment_speed"`
	Confidence           float64 `json:"confidence"`
}

// PersonalityProfile represents user's trading personality
type PersonalityProfile struct {
	TraderType            string             `json:"trader_type"`            // analytical, intuitive, systematic, discretionary
	DecisionMaking        string             `json:"decision_making"`        // rational, emotional, mixed
	InformationProcessing string             `json:"information_processing"` // sequential, holistic
	PlanningOrientation   string             `json:"planning_orientation"`   // structured, flexible, adaptive
	StressResponse        string             `json:"stress_response"`        // calm, reactive, volatile
	ConfidenceLevel       string             `json:"confidence_level"`       // low, moderate, high, overconfident
	LearningStyle         string             `json:"learning_style"`         // visual, auditory, kinesthetic, reading
	SocialInfluence       float64            `json:"social_influence"`       // susceptibility to others' opinions
	Patience              float64            `json:"patience"`               // 0-1 scale
	Discipline            float64            `json:"discipline"`
	Adaptability          float64            `json:"adaptability"`
	Curiosity             float64            `json:"curiosity"`
	Optimism              float64            `json:"optimism"`
	Traits                map[string]float64 `json:"traits"` // big five personality traits
	Biases                []*CognitiveBias   `json:"biases"`
	Confidence            float64            `json:"confidence"`
}

// CognitiveBias represents a cognitive bias
type CognitiveBias struct {
	Type        string  `json:"type"`      // confirmation, anchoring, overconfidence, etc.
	Strength    float64 `json:"strength"`  // 0-1 scale
	Frequency   float64 `json:"frequency"` // how often it manifests
	Impact      float64 `json:"impact"`    // impact on trading performance
	Description string  `json:"description"`
}

// UserBehaviorPreferences represents user's trading preferences
type UserBehaviorPreferences struct {
	NotificationPrefs *NotificationPreferences `json:"notification_preferences"`
	UIPreferences     *UIPreferences           `json:"ui_preferences"`
	AnalysisPrefs     *AnalysisPreferences     `json:"analysis_preferences"`
	AutomationPrefs   *AutomationPreferences   `json:"automation_preferences"`
	ReportingPrefs    *ReportingPreferences    `json:"reporting_preferences"`
	PrivacyPrefs      *PrivacyPreferences      `json:"privacy_preferences"`
	LanguagePrefs     *LanguagePreferences     `json:"language_preferences"`
	Confidence        float64                  `json:"confidence"`
}

// NotificationPreferences represents notification preferences
type NotificationPreferences struct {
	PriceAlerts        bool                   `json:"price_alerts"`
	TradingSignals     bool                   `json:"trading_signals"`
	PortfolioUpdates   bool                   `json:"portfolio_updates"`
	NewsAlerts         bool                   `json:"news_alerts"`
	RiskWarnings       bool                   `json:"risk_warnings"`
	PerformanceReports bool                   `json:"performance_reports"`
	Frequency          string                 `json:"frequency"` // real_time, hourly, daily, weekly
	Channels           []string               `json:"channels"`  // email, sms, push, in_app
	QuietHours         *TimeRange             `json:"quiet_hours"`
	Urgency            map[string]string      `json:"urgency"` // immediate, normal, low
	Customizations     map[string]interface{} `json:"customizations"`
}

// UIPreferences represents user interface preferences
type UIPreferences struct {
	Theme             string                 `json:"theme"`       // light, dark, auto
	Layout            string                 `json:"layout"`      // compact, standard, spacious
	ChartStyle        string                 `json:"chart_style"` // candlestick, line, bar
	DefaultTimeframe  string                 `json:"default_timeframe"`
	WidgetPreferences map[string]bool        `json:"widget_preferences"`
	DashboardLayout   []string               `json:"dashboard_layout"`
	ColorScheme       string                 `json:"color_scheme"`
	FontSize          string                 `json:"font_size"`
	AnimationsEnabled bool                   `json:"animations_enabled"`
	SoundEnabled      bool                   `json:"sound_enabled"`
	Customizations    map[string]interface{} `json:"customizations"`
}

// AnalysisPreferences represents analysis preferences
type AnalysisPreferences struct {
	PreferredIndicators  []string               `json:"preferred_indicators"`
	AnalysisDepth        string                 `json:"analysis_depth"`
	TimeHorizons         []string               `json:"time_horizons"`
	DataSources          []string               `json:"data_sources"`
	ConfidenceDisplay    bool                   `json:"confidence_display"`
	UncertaintyDisplay   bool                   `json:"uncertainty_display"`
	AlternativeScenarios bool                   `json:"alternative_scenarios"`
	HistoricalComparison bool                   `json:"historical_comparison"`
	Customizations       map[string]interface{} `json:"customizations"`
}

// AutomationPreferences represents automation preferences
type AutomationPreferences struct {
	AutoTrading        bool                   `json:"auto_trading"`
	AutoRebalancing    bool                   `json:"auto_rebalancing"`
	AutoRiskManagement bool                   `json:"auto_risk_management"`
	AutoNotifications  bool                   `json:"auto_notifications"`
	ApprovalRequired   map[string]bool        `json:"approval_required"`
	AutomationLimits   map[string]float64     `json:"automation_limits"`
	SafetyChecks       []string               `json:"safety_checks"`
	Customizations     map[string]interface{} `json:"customizations"`
}

// ReportingPreferences represents reporting preferences
type ReportingPreferences struct {
	ReportFrequency   string                 `json:"report_frequency"`
	ReportTypes       []string               `json:"report_types"`
	MetricsIncluded   []string               `json:"metrics_included"`
	ComparisonPeriods []string               `json:"comparison_periods"`
	DetailLevel       string                 `json:"detail_level"`
	Format            string                 `json:"format"` // pdf, html, json
	DeliveryMethod    string                 `json:"delivery_method"`
	Customizations    map[string]interface{} `json:"customizations"`
}

// PrivacyPreferences represents privacy preferences
type PrivacyPreferences struct {
	DataSharing           bool                   `json:"data_sharing"`
	AnalyticsTracking     bool                   `json:"analytics_tracking"`
	PersonalizationLevel  string                 `json:"personalization_level"`
	DataRetention         string                 `json:"data_retention"`
	ThirdPartyIntegration bool                   `json:"third_party_integration"`
	Anonymization         bool                   `json:"anonymization"`
	Customizations        map[string]interface{} `json:"customizations"`
}

// LanguagePreferences represents language preferences
type LanguagePreferences struct {
	PrimaryLanguage    string                 `json:"primary_language"`
	SecondaryLanguages []string               `json:"secondary_languages"`
	RegionalSettings   string                 `json:"regional_settings"`
	CurrencyDisplay    string                 `json:"currency_display"`
	DateFormat         string                 `json:"date_format"`
	NumberFormat       string                 `json:"number_format"`
	Customizations     map[string]interface{} `json:"customizations"`
}

// BehaviorPattern represents a detected behavior pattern
type BehaviorPattern struct {
	ID               string                      `json:"id"`
	Type             string                      `json:"type"` // temporal, conditional, sequential, cyclical
	Name             string                      `json:"name"`
	Description      string                      `json:"description"`
	Frequency        float64                     `json:"frequency"` // how often it occurs
	Strength         float64                     `json:"strength"`  // how strong the pattern is
	Confidence       float64                     `json:"confidence"`
	Conditions       []*BehaviorPatternCondition `json:"conditions"`
	Outcomes         []*BehaviorPatternOutcome   `json:"outcomes"`
	Triggers         []string                    `json:"triggers"`
	Context          map[string]interface{}      `json:"context"`
	FirstObserved    time.Time                   `json:"first_observed"`
	LastObserved     time.Time                   `json:"last_observed"`
	ObservationCount int                         `json:"observation_count"`
	Metadata         map[string]interface{}      `json:"metadata"`
}

// BehaviorPatternCondition represents a condition for a behavior pattern
type BehaviorPatternCondition struct {
	Type       string                 `json:"type"` // market, time, portfolio, emotional
	Field      string                 `json:"field"`
	Operator   string                 `json:"operator"` // eq, gt, lt, in, between
	Value      interface{}            `json:"value"`
	Confidence float64                `json:"confidence"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// BehaviorPatternOutcome represents an outcome of a behavior pattern
type BehaviorPatternOutcome struct {
	Type        string                 `json:"type"` // action, decision, performance
	Action      string                 `json:"action"`
	Probability float64                `json:"probability"`
	Impact      float64                `json:"impact"`
	Confidence  float64                `json:"confidence"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// UserPerformanceProfile represents user's performance characteristics
type UserPerformanceProfile struct {
	OverallReturn         float64            `json:"overall_return"`
	AnnualizedReturn      float64            `json:"annualized_return"`
	Volatility            float64            `json:"volatility"`
	SharpeRatio           float64            `json:"sharpe_ratio"`
	MaxDrawdown           float64            `json:"max_drawdown"`
	WinRate               float64            `json:"win_rate"`
	ProfitFactor          float64            `json:"profit_factor"`
	AverageWin            float64            `json:"average_win"`
	AverageLoss           float64            `json:"average_loss"`
	TotalTrades           int                `json:"total_trades"`
	SuccessfulTrades      int                `json:"successful_trades"`
	ConsecutiveWins       int                `json:"consecutive_wins"`
	ConsecutiveLosses     int                `json:"consecutive_losses"`
	BestTrade             float64            `json:"best_trade"`
	WorstTrade            float64            `json:"worst_trade"`
	AverageTradeTime      time.Duration      `json:"average_trade_time"`
	RiskAdjustedReturn    float64            `json:"risk_adjusted_return"`
	PerformanceTrend      string             `json:"performance_trend"` // improving, declining, stable
	PerformanceByPeriod   map[string]float64 `json:"performance_by_period"`
	PerformanceByAsset    map[string]float64 `json:"performance_by_asset"`
	PerformanceByStrategy map[string]float64 `json:"performance_by_strategy"`
	Confidence            float64            `json:"confidence"`
	LastUpdated           time.Time          `json:"last_updated"`
}

// LearningProgress represents learning progress for a user
type LearningProgress struct {
	TotalObservations   int                    `json:"total_observations"`
	LearningRate        float64                `json:"learning_rate"`
	ModelAccuracy       float64                `json:"model_accuracy"`
	PredictionAccuracy  float64                `json:"prediction_accuracy"`
	PatternRecognition  float64                `json:"pattern_recognition"`
	PreferenceStability float64                `json:"preference_stability"`
	ProfileCompleteness float64                `json:"profile_completeness"`
	LearningVelocity    float64                `json:"learning_velocity"`
	AdaptationRate      float64                `json:"adaptation_rate"`
	ConfidenceGrowth    float64                `json:"confidence_growth"`
	Milestones          []*LearningMilestone   `json:"milestones"`
	LastModelUpdate     time.Time              `json:"last_model_update"`
	NextUpdateDue       time.Time              `json:"next_update_due"`
	Metadata            map[string]interface{} `json:"metadata"`
}

// LearningMilestone represents a learning milestone
type LearningMilestone struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	AchievedAt  time.Time              `json:"achieved_at"`
	Value       float64                `json:"value"`
	Threshold   float64                `json:"threshold"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// PersonalizedRecommendation represents a personalized recommendation
type PersonalizedRecommendation struct {
	ID              string                 `json:"id"`
	Type            string                 `json:"type"` // trade, strategy, education, risk
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	Reasoning       []string               `json:"reasoning"`
	Confidence      float64                `json:"confidence"`
	Priority        string                 `json:"priority"` // low, medium, high, urgent
	Category        string                 `json:"category"`
	ActionRequired  bool                   `json:"action_required"`
	Deadline        *time.Time             `json:"deadline,omitempty"`
	Parameters      map[string]interface{} `json:"parameters"`
	ExpectedOutcome *UserExpectedOutcome   `json:"expected_outcome"`
	RiskAssessment  *RecommendationRisk    `json:"risk_assessment"`
	Personalization *PersonalizationInfo   `json:"personalization"`
	CreatedAt       time.Time              `json:"created_at"`
	ExpiresAt       *time.Time             `json:"expires_at,omitempty"`
	Status          string                 `json:"status"` // pending, accepted, rejected, expired
	Metadata        map[string]interface{} `json:"metadata"`
}

// ExpectedOutcome represents expected outcome of a recommendation
type UserExpectedOutcome struct {
	ProbabilityOfSuccess float64                `json:"probability_of_success"`
	ExpectedReturn       float64                `json:"expected_return"`
	ExpectedRisk         float64                `json:"expected_risk"`
	TimeHorizon          time.Duration          `json:"time_horizon"`
	ConfidenceInterval   *ConfidenceInterval    `json:"confidence_interval"`
	AlternativeOutcomes  []*AlternativeOutcome  `json:"alternative_outcomes"`
	Metadata             map[string]interface{} `json:"metadata"`
}

// ConfidenceInterval represents a confidence interval
type ConfidenceInterval struct {
	Lower      float64 `json:"lower"`
	Upper      float64 `json:"upper"`
	Confidence float64 `json:"confidence"` // e.g., 0.95 for 95% confidence
}

// AlternativeOutcome represents an alternative outcome
type AlternativeOutcome struct {
	Scenario    string  `json:"scenario"`
	Probability float64 `json:"probability"`
	Return      float64 `json:"return"`
	Risk        float64 `json:"risk"`
}

// RecommendationRisk represents risk assessment for a recommendation
type RecommendationRisk struct {
	OverallRisk float64                `json:"overall_risk"`
	RiskFactors []*BehaviorRiskFactor  `json:"risk_factors"`
	Mitigation  []string               `json:"mitigation"`
	MaxLoss     float64                `json:"max_loss"`
	Probability float64                `json:"probability"`
	RiskReward  float64                `json:"risk_reward"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// BehaviorRiskFactor represents a risk factor
type BehaviorRiskFactor struct {
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Impact      float64 `json:"impact"`
	Probability float64 `json:"probability"`
	Severity    string  `json:"severity"`
}

// PersonalizationInfo represents personalization information
type PersonalizationInfo struct {
	PersonalizationScore float64                `json:"personalization_score"`
	UserFactors          []string               `json:"user_factors"`
	BehaviorFactors      []string               `json:"behavior_factors"`
	ContextFactors       []string               `json:"context_factors"`
	Adaptations          []string               `json:"adaptations"`
	Metadata             map[string]interface{} `json:"metadata"`
}

// BehaviorEvent represents a user behavior event
type BehaviorEvent struct {
	ID        string                 `json:"id"`
	UserID    uuid.UUID              `json:"user_id"`
	Type      string                 `json:"type"` // trade, view, search, analyze, etc.
	Action    string                 `json:"action"`
	Context   *BehaviorContext       `json:"context"`
	Outcome   *BehaviorOutcome       `json:"outcome"`
	Timestamp time.Time              `json:"timestamp"`
	Duration  time.Duration          `json:"duration,omitempty"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// BehaviorContext represents the context of a behavior event
type BehaviorContext struct {
	MarketConditions   string                 `json:"market_conditions"`
	PortfolioState     map[string]interface{} `json:"portfolio_state"`
	TimeOfDay          string                 `json:"time_of_day"`
	DayOfWeek          string                 `json:"day_of_week"`
	SessionDuration    time.Duration          `json:"session_duration"`
	PreviousActions    []string               `json:"previous_actions"`
	EmotionalState     string                 `json:"emotional_state,omitempty"`
	InformationSources []string               `json:"information_sources"`
	ExternalFactors    map[string]interface{} `json:"external_factors"`
	Metadata           map[string]interface{} `json:"metadata"`
}

// BehaviorOutcome represents the outcome of a behavior event
type BehaviorOutcome struct {
	Success         bool                   `json:"success"`
	Performance     float64                `json:"performance,omitempty"`
	Satisfaction    float64                `json:"satisfaction,omitempty"`
	TimeToDecision  time.Duration          `json:"time_to_decision,omitempty"`
	ConfidenceLevel float64                `json:"confidence_level,omitempty"`
	Regret          float64                `json:"regret,omitempty"`
	LearningValue   float64                `json:"learning_value"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// LearningModel represents a machine learning model for user behavior
type LearningModel struct {
	ID           string                 `json:"id"`
	Type         string                 `json:"type"`    // neural_network, decision_tree, ensemble
	Purpose      string                 `json:"purpose"` // preference, pattern, performance
	Version      string                 `json:"version"`
	Accuracy     float64                `json:"accuracy"`
	TrainingData int                    `json:"training_data"`
	LastTrained  time.Time              `json:"last_trained"`
	NextTraining time.Time              `json:"next_training"`
	Parameters   map[string]interface{} `json:"parameters"`
	Performance  map[string]float64     `json:"performance"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// Supporting component types
type BehaviorAnalyzer struct {
	config *BehaviorAnalyzerConfig
	logger *observability.Logger
}

type PatternRecognizer struct {
	config *PatternRecognizerConfig
	logger *observability.Logger
}

type PreferenceEngine struct {
	config *PreferenceEngineConfig
	logger *observability.Logger
}

type RecommendationEngine struct {
	config *RecommendationEngineConfig
	logger *observability.Logger
}

type PersonalityProfiler struct {
	config *PersonalityProfilerConfig
	logger *observability.Logger
}

type RiskProfiler struct {
	config *RiskProfilerConfig
	logger *observability.Logger
}

// Configuration types for components
type BehaviorAnalyzerConfig struct {
	AnalysisWindow       time.Duration `json:"analysis_window"`
	MinEventsForAnalysis int           `json:"min_events_for_analysis"`
	ConfidenceThreshold  float64       `json:"confidence_threshold"`
}

type PatternRecognizerConfig struct {
	MinPatternOccurrences    int           `json:"min_pattern_occurrences"`
	PatternStrengthThreshold float64       `json:"pattern_strength_threshold"`
	TemporalWindowSize       time.Duration `json:"temporal_window_size"`
}

type PreferenceEngineConfig struct {
	LearningRate        float64 `json:"learning_rate"`
	DecayRate           float64 `json:"decay_rate"`
	ConfidenceThreshold float64 `json:"confidence_threshold"`
}

type RecommendationEngineConfig struct {
	MaxRecommendations    int     `json:"max_recommendations"`
	MinConfidence         float64 `json:"min_confidence"`
	PersonalizationWeight float64 `json:"personalization_weight"`
}

type PersonalityProfilerConfig struct {
	TraitUpdateRate        float64 `json:"trait_update_rate"`
	BiasDetectionThreshold float64 `json:"bias_detection_threshold"`
	ConfidenceThreshold    float64 `json:"confidence_threshold"`
}

type RiskProfilerConfig struct {
	RiskUpdateRate      float64       `json:"risk_update_rate"`
	VolatilityWindow    time.Duration `json:"volatility_window"`
	ConfidenceThreshold float64       `json:"confidence_threshold"`
}

// NewUserBehaviorLearningEngine creates a new user behavior learning engine
func NewUserBehaviorLearningEngine(logger *observability.Logger) *UserBehaviorLearningEngine {
	config := &UserBehaviorConfig{
		LearningRate:               0.1,
		MinObservations:            10,
		PatternDetectionWindow:     30 * 24 * time.Hour, // 30 days
		PreferenceUpdateRate:       0.05,
		PersonalityUpdateRate:      0.02,
		RiskToleranceUpdateRate:    0.03,
		RecommendationThreshold:    0.7,
		MaxHistorySize:             10000,
		EnableRealTimeLearning:     true,
		EnablePersonalityProfiling: true,
		EnableRiskProfiling:        true,
		EnablePatternRecognition:   true,
		ModelUpdateInterval:        24 * time.Hour,
		ConfidenceThreshold:        0.6,
	}

	engine := &UserBehaviorLearningEngine{
		logger:               logger,
		config:               config,
		behaviorAnalyzer:     NewBehaviorAnalyzer(logger),
		patternRecognizer:    NewPatternRecognizer(logger),
		preferenceEngine:     NewPreferenceEngine(logger),
		recommendationEngine: NewRecommendationEngine(logger),
		personalityProfiler:  NewPersonalityProfiler(logger),
		riskProfiler:         NewRiskProfiler(logger),
		userProfiles:         make(map[uuid.UUID]*UserBehaviorProfile),
		behaviorHistory:      make(map[uuid.UUID][]*BehaviorEvent),
		learningModels:       make(map[string]*LearningModel),
		lastUpdate:           time.Now(),
	}

	logger.Info(context.Background(), "User behavior learning engine initialized", map[string]interface{}{
		"learning_rate":            config.LearningRate,
		"min_observations":         config.MinObservations,
		"pattern_detection_window": config.PatternDetectionWindow.String(),
		"real_time_learning":       config.EnableRealTimeLearning,
		"personality_profiling":    config.EnablePersonalityProfiling,
		"risk_profiling":           config.EnableRiskProfiling,
		"pattern_recognition":      config.EnablePatternRecognition,
	})

	return engine
}

// LearnFromBehavior learns from a user behavior event
func (u *UserBehaviorLearningEngine) LearnFromBehavior(ctx context.Context, event *BehaviorEvent) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.logger.Info(ctx, "Learning from user behavior", map[string]interface{}{
		"user_id":    event.UserID,
		"event_type": event.Type,
		"action":     event.Action,
	})

	// Add event to history
	if err := u.addBehaviorEvent(event); err != nil {
		return fmt.Errorf("failed to add behavior event: %w", err)
	}

	// Get or create user profile
	profile := u.getUserProfile(event.UserID)

	// Update behavior analysis
	if err := u.updateBehaviorAnalysis(ctx, profile, event); err != nil {
		u.logger.Warn(ctx, "Failed to update behavior analysis", map[string]interface{}{
			"error":   err.Error(),
			"user_id": event.UserID,
		})
	}

	// Update pattern recognition
	if u.config.EnablePatternRecognition {
		if err := u.updatePatternRecognition(ctx, profile, event); err != nil {
			u.logger.Warn(ctx, "Failed to update pattern recognition", map[string]interface{}{
				"error":   err.Error(),
				"user_id": event.UserID,
			})
		}
	}

	// Update preferences
	if err := u.updatePreferences(ctx, profile, event); err != nil {
		u.logger.Warn(ctx, "Failed to update preferences", map[string]interface{}{
			"error":   err.Error(),
			"user_id": event.UserID,
		})
	}

	// Update personality profile
	if u.config.EnablePersonalityProfiling {
		if err := u.updatePersonalityProfile(ctx, profile, event); err != nil {
			u.logger.Warn(ctx, "Failed to update personality profile", map[string]interface{}{
				"error":   err.Error(),
				"user_id": event.UserID,
			})
		}
	}

	// Update risk profile
	if u.config.EnableRiskProfiling {
		if err := u.updateRiskProfile(ctx, profile, event); err != nil {
			u.logger.Warn(ctx, "Failed to update risk profile", map[string]interface{}{
				"error":   err.Error(),
				"user_id": event.UserID,
			})
		}
	}

	// Update performance metrics
	if err := u.updatePerformanceMetrics(ctx, profile, event); err != nil {
		u.logger.Warn(ctx, "Failed to update performance metrics", map[string]interface{}{
			"error":   err.Error(),
			"user_id": event.UserID,
		})
	}

	// Update observation count first
	profile.ObservationCount++
	profile.LastUpdated = time.Now()

	// Update learning progress
	u.updateLearningProgress(profile)

	// Generate recommendations if needed
	if profile.ObservationCount >= u.config.MinObservations {
		if err := u.generateRecommendations(ctx, profile); err != nil {
			u.logger.Warn(ctx, "Failed to generate recommendations", map[string]interface{}{
				"error":   err.Error(),
				"user_id": event.UserID,
			})
		}
	}

	u.logger.Info(ctx, "User behavior learning completed", map[string]interface{}{
		"user_id":           event.UserID,
		"observation_count": profile.ObservationCount,
		"confidence":        profile.Confidence,
	})

	return nil
}

// GetUserProfile retrieves a user's behavior profile
func (u *UserBehaviorLearningEngine) GetUserProfile(ctx context.Context, userID uuid.UUID) (*UserBehaviorProfile, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	profile, exists := u.userProfiles[userID]
	if !exists {
		return nil, fmt.Errorf("user profile not found for user %s", userID)
	}

	return profile, nil
}

// GetPersonalizedRecommendations retrieves personalized recommendations for a user
func (u *UserBehaviorLearningEngine) GetPersonalizedRecommendations(ctx context.Context, userID uuid.UUID, limit int) ([]*PersonalizedRecommendation, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	profile, exists := u.userProfiles[userID]
	if !exists {
		return nil, fmt.Errorf("user profile not found for user %s", userID)
	}

	// Filter active recommendations
	var activeRecommendations []*PersonalizedRecommendation
	now := time.Now()

	for _, rec := range profile.Recommendations {
		if rec.Status == "pending" && (rec.ExpiresAt == nil || rec.ExpiresAt.After(now)) {
			activeRecommendations = append(activeRecommendations, rec)
		}
	}

	// Sort by priority and confidence
	sort.Slice(activeRecommendations, func(i, j int) bool {
		// Priority order: urgent > high > medium > low
		priorityOrder := map[string]int{"urgent": 4, "high": 3, "medium": 2, "low": 1}

		if priorityOrder[activeRecommendations[i].Priority] != priorityOrder[activeRecommendations[j].Priority] {
			return priorityOrder[activeRecommendations[i].Priority] > priorityOrder[activeRecommendations[j].Priority]
		}

		return activeRecommendations[i].Confidence > activeRecommendations[j].Confidence
	})

	// Apply limit
	if limit > 0 && len(activeRecommendations) > limit {
		activeRecommendations = activeRecommendations[:limit]
	}

	return activeRecommendations, nil
}

// Helper methods

func (u *UserBehaviorLearningEngine) addBehaviorEvent(event *BehaviorEvent) error {
	// Add to history
	history := u.behaviorHistory[event.UserID]
	history = append(history, event)

	// Maintain max history size
	if len(history) > u.config.MaxHistorySize {
		history = history[len(history)-u.config.MaxHistorySize:]
	}

	u.behaviorHistory[event.UserID] = history
	return nil
}

func (u *UserBehaviorLearningEngine) getUserProfile(userID uuid.UUID) *UserBehaviorProfile {
	profile, exists := u.userProfiles[userID]
	if !exists {
		profile = &UserBehaviorProfile{
			UserID:             userID,
			CreatedAt:          time.Now(),
			LastUpdated:        time.Now(),
			TradingStyle:       &TradingStyleProfile{Confidence: 0.0},
			RiskProfile:        &UserRiskProfile{Confidence: 0.0},
			PersonalityProfile: &PersonalityProfile{Confidence: 0.0, Traits: make(map[string]float64)},
			Preferences:        &UserBehaviorPreferences{Confidence: 0.0},
			BehaviorPatterns:   []*BehaviorPattern{},
			PerformanceMetrics: &UserPerformanceProfile{Confidence: 0.0},
			LearningProgress:   &LearningProgress{Milestones: []*LearningMilestone{}},
			Recommendations:    []*PersonalizedRecommendation{},
			Confidence:         0.0,
			ObservationCount:   0,
			Metadata:           make(map[string]interface{}),
		}
		u.userProfiles[userID] = profile
	}
	return profile
}

func (u *UserBehaviorLearningEngine) updateBehaviorAnalysis(ctx context.Context, profile *UserBehaviorProfile, event *BehaviorEvent) error {
	// Simplified behavior analysis
	return u.behaviorAnalyzer.AnalyzeBehavior(ctx, profile, event)
}

func (u *UserBehaviorLearningEngine) updatePatternRecognition(ctx context.Context, profile *UserBehaviorProfile, event *BehaviorEvent) error {
	// Simplified pattern recognition
	return u.patternRecognizer.RecognizePatterns(ctx, profile, event)
}

func (u *UserBehaviorLearningEngine) updatePreferences(ctx context.Context, profile *UserBehaviorProfile, event *BehaviorEvent) error {
	// Simplified preference learning
	return u.preferenceEngine.UpdatePreferences(ctx, profile, event)
}

func (u *UserBehaviorLearningEngine) updatePersonalityProfile(ctx context.Context, profile *UserBehaviorProfile, event *BehaviorEvent) error {
	// Simplified personality profiling
	return u.personalityProfiler.UpdatePersonality(ctx, profile, event)
}

func (u *UserBehaviorLearningEngine) updateRiskProfile(ctx context.Context, profile *UserBehaviorProfile, event *BehaviorEvent) error {
	// Simplified risk profiling
	return u.riskProfiler.UpdateRiskProfile(ctx, profile, event)
}

func (u *UserBehaviorLearningEngine) updatePerformanceMetrics(ctx context.Context, profile *UserBehaviorProfile, event *BehaviorEvent) error {
	// Update performance metrics based on event outcome
	if event.Outcome != nil && event.Type == "trade" {
		metrics := profile.PerformanceMetrics

		if event.Outcome.Success {
			metrics.SuccessfulTrades++
			if event.Outcome.Performance > 0 {
				metrics.AverageWin = (metrics.AverageWin*float64(metrics.SuccessfulTrades-1) + event.Outcome.Performance) / float64(metrics.SuccessfulTrades)
			}
		} else {
			if event.Outcome.Performance < 0 {
				lossTrades := metrics.TotalTrades - metrics.SuccessfulTrades
				if lossTrades > 0 {
					metrics.AverageLoss = (metrics.AverageLoss*float64(lossTrades-1) + math.Abs(event.Outcome.Performance)) / float64(lossTrades)
				} else {
					metrics.AverageLoss = math.Abs(event.Outcome.Performance)
				}
			}
		}

		metrics.TotalTrades++
		metrics.WinRate = float64(metrics.SuccessfulTrades) / float64(metrics.TotalTrades)

		if metrics.AverageLoss > 0 {
			metrics.ProfitFactor = metrics.AverageWin / metrics.AverageLoss
		}

		metrics.LastUpdated = time.Now()
		metrics.Confidence = math.Min(1.0, float64(metrics.TotalTrades)/100.0) // Confidence grows with trade count
	}

	return nil
}

func (u *UserBehaviorLearningEngine) updateLearningProgress(profile *UserBehaviorProfile) {
	progress := profile.LearningProgress

	// Update basic metrics
	progress.TotalObservations = profile.ObservationCount
	progress.ProfileCompleteness = u.calculateProfileCompleteness(profile)
	progress.LearningVelocity = u.calculateLearningVelocity(profile)

	// Update overall confidence
	totalConfidence := profile.TradingStyle.Confidence +
		profile.RiskProfile.Confidence +
		profile.PersonalityProfile.Confidence +
		profile.Preferences.Confidence +
		profile.PerformanceMetrics.Confidence

	profile.Confidence = totalConfidence / 5.0

	// Check for milestones
	u.checkLearningMilestones(profile)
}

func (u *UserBehaviorLearningEngine) calculateProfileCompleteness(profile *UserBehaviorProfile) float64 {
	// Simplified completeness calculation
	completeness := 0.0

	if profile.TradingStyle.Confidence > 0.05 {
		completeness += 0.2
	}
	if profile.RiskProfile.Confidence > 0.02 {
		completeness += 0.2
	}
	if profile.PersonalityProfile.Confidence > 0.01 {
		completeness += 0.2
	}
	if profile.Preferences.Confidence > 0.03 {
		completeness += 0.2
	}
	if len(profile.BehaviorPatterns) > 0 {
		completeness += 0.2
	}

	return completeness
}

func (u *UserBehaviorLearningEngine) calculateLearningVelocity(profile *UserBehaviorProfile) float64 {
	// Simplified learning velocity calculation
	if profile.ObservationCount < 10 {
		return 0.0
	}

	timeSinceCreation := time.Since(profile.CreatedAt)
	if timeSinceCreation.Hours() < 1 {
		return 0.0
	}

	return float64(profile.ObservationCount) / timeSinceCreation.Hours()
}

func (u *UserBehaviorLearningEngine) checkLearningMilestones(profile *UserBehaviorProfile) {
	milestones := []*LearningMilestone{
		{
			ID:          "first_10_observations",
			Name:        "First 10 Observations",
			Description: "Completed first 10 behavior observations",
			Threshold:   10,
		},
		{
			ID:          "profile_50_complete",
			Name:        "Profile 50% Complete",
			Description: "User profile is 50% complete",
			Threshold:   0.5,
		},
		{
			ID:          "high_confidence",
			Name:        "High Confidence Profile",
			Description: "Achieved high confidence in user profile",
			Threshold:   0.8,
		},
	}

	for _, milestone := range milestones {
		// Check if milestone already achieved
		achieved := false
		for _, existing := range profile.LearningProgress.Milestones {
			if existing.ID == milestone.ID {
				achieved = true
				break
			}
		}

		if !achieved {
			var value float64
			switch milestone.ID {
			case "first_10_observations":
				value = float64(profile.ObservationCount)
			case "profile_50_complete":
				value = profile.LearningProgress.ProfileCompleteness
			case "high_confidence":
				value = profile.Confidence
			}

			if value >= milestone.Threshold {
				milestone.AchievedAt = time.Now()
				milestone.Value = value
				profile.LearningProgress.Milestones = append(profile.LearningProgress.Milestones, milestone)
			}
		}
	}
}

func (u *UserBehaviorLearningEngine) generateRecommendations(ctx context.Context, profile *UserBehaviorProfile) error {
	// Simplified recommendation generation
	return u.recommendationEngine.GenerateRecommendations(ctx, profile)
}

// Component implementations (simplified)

func NewBehaviorAnalyzer(logger *observability.Logger) *BehaviorAnalyzer {
	return &BehaviorAnalyzer{
		config: &BehaviorAnalyzerConfig{
			AnalysisWindow:       7 * 24 * time.Hour,
			MinEventsForAnalysis: 5,
			ConfidenceThreshold:  0.6,
		},
		logger: logger,
	}
}

func (ba *BehaviorAnalyzer) AnalyzeBehavior(ctx context.Context, profile *UserBehaviorProfile, event *BehaviorEvent) error {
	// Simplified behavior analysis
	if event.Type == "trade" {
		// Update trading style
		if profile.TradingStyle.PrimaryStyle == "" {
			if event.Duration < time.Hour {
				profile.TradingStyle.PrimaryStyle = "scalper"
			} else if event.Duration < 24*time.Hour {
				profile.TradingStyle.PrimaryStyle = "day_trader"
			} else {
				profile.TradingStyle.PrimaryStyle = "swing_trader"
			}
		}

		// Update confidence
		profile.TradingStyle.Confidence = math.Min(1.0, profile.TradingStyle.Confidence+0.1)
	}

	return nil
}

func NewPatternRecognizer(logger *observability.Logger) *PatternRecognizer {
	return &PatternRecognizer{
		config: &PatternRecognizerConfig{
			MinPatternOccurrences:    3,
			PatternStrengthThreshold: 0.7,
			TemporalWindowSize:       24 * time.Hour,
		},
		logger: logger,
	}
}

func (pr *PatternRecognizer) RecognizePatterns(ctx context.Context, profile *UserBehaviorProfile, event *BehaviorEvent) error {
	// Simplified pattern recognition
	// This would implement sophisticated pattern detection algorithms
	return nil
}

func NewPreferenceEngine(logger *observability.Logger) *PreferenceEngine {
	return &PreferenceEngine{
		config: &PreferenceEngineConfig{
			LearningRate:        0.1,
			DecayRate:           0.01,
			ConfidenceThreshold: 0.6,
		},
		logger: logger,
	}
}

func (pe *PreferenceEngine) UpdatePreferences(ctx context.Context, profile *UserBehaviorProfile, event *BehaviorEvent) error {
	// Simplified preference learning
	if profile.Preferences.UIPreferences == nil {
		profile.Preferences.UIPreferences = &UIPreferences{
			Theme:             "dark",
			Layout:            "standard",
			ChartStyle:        "candlestick",
			DefaultTimeframe:  "1h",
			WidgetPreferences: make(map[string]bool),
			DashboardLayout:   []string{},
			ColorScheme:       "default",
			FontSize:          "medium",
			AnimationsEnabled: true,
			SoundEnabled:      false,
			Customizations:    make(map[string]interface{}),
		}
	}

	profile.Preferences.Confidence = math.Min(1.0, profile.Preferences.Confidence+0.05)
	return nil
}

func NewRecommendationEngine(logger *observability.Logger) *RecommendationEngine {
	return &RecommendationEngine{
		config: &RecommendationEngineConfig{
			MaxRecommendations:    10,
			MinConfidence:         0.7,
			PersonalizationWeight: 0.8,
		},
		logger: logger,
	}
}

func (re *RecommendationEngine) GenerateRecommendations(ctx context.Context, profile *UserBehaviorProfile) error {
	// Simplified recommendation generation
	if profile.Confidence > 0.6 && len(profile.Recommendations) < 5 {
		recommendation := &PersonalizedRecommendation{
			ID:             uuid.New().String(),
			Type:           "strategy",
			Title:          "Optimize Your Trading Strategy",
			Description:    "Based on your trading patterns, consider adjusting your position sizing strategy",
			Reasoning:      []string{"Your win rate is above average", "Risk-adjusted returns could be improved"},
			Confidence:     0.75,
			Priority:       "medium",
			Category:       "optimization",
			ActionRequired: false,
			Parameters:     make(map[string]interface{}),
			ExpectedOutcome: &UserExpectedOutcome{
				ProbabilityOfSuccess: 0.7,
				ExpectedReturn:       0.05,
				ExpectedRisk:         0.02,
				TimeHorizon:          30 * 24 * time.Hour,
			},
			RiskAssessment: &RecommendationRisk{
				OverallRisk: 0.2,
				RiskFactors: []*BehaviorRiskFactor{},
				Mitigation:  []string{"Start with small position sizes", "Monitor performance closely"},
				MaxLoss:     0.01,
				Probability: 0.1,
				RiskReward:  2.5,
			},
			Personalization: &PersonalizationInfo{
				PersonalizationScore: 0.8,
				UserFactors:          []string{"trading_style", "risk_tolerance"},
				BehaviorFactors:      []string{"win_rate", "position_sizing"},
				ContextFactors:       []string{"market_conditions"},
				Adaptations:          []string{"personalized_thresholds"},
			},
			CreatedAt: time.Now(),
			Status:    "pending",
			Metadata:  make(map[string]interface{}),
		}

		profile.Recommendations = append(profile.Recommendations, recommendation)
	}

	return nil
}

func NewPersonalityProfiler(logger *observability.Logger) *PersonalityProfiler {
	return &PersonalityProfiler{
		config: &PersonalityProfilerConfig{
			TraitUpdateRate:        0.02,
			BiasDetectionThreshold: 0.6,
			ConfidenceThreshold:    0.6,
		},
		logger: logger,
	}
}

func (pp *PersonalityProfiler) UpdatePersonality(ctx context.Context, profile *UserBehaviorProfile, event *BehaviorEvent) error {
	// Simplified personality profiling
	if profile.PersonalityProfile.TraderType == "" {
		if event.Context != nil && event.Context.SessionDuration > 2*time.Hour {
			profile.PersonalityProfile.TraderType = "analytical"
		} else {
			profile.PersonalityProfile.TraderType = "intuitive"
		}
	}

	// Update traits
	if profile.PersonalityProfile.Traits == nil {
		profile.PersonalityProfile.Traits = make(map[string]float64)
	}

	// Big Five traits (simplified)
	profile.PersonalityProfile.Traits["openness"] = 0.7
	profile.PersonalityProfile.Traits["conscientiousness"] = 0.8
	profile.PersonalityProfile.Traits["extraversion"] = 0.5
	profile.PersonalityProfile.Traits["agreeableness"] = 0.6
	profile.PersonalityProfile.Traits["neuroticism"] = 0.3

	profile.PersonalityProfile.Confidence = math.Min(1.0, profile.PersonalityProfile.Confidence+0.02)
	return nil
}

func NewRiskProfiler(logger *observability.Logger) *RiskProfiler {
	return &RiskProfiler{
		config: &RiskProfilerConfig{
			RiskUpdateRate:      0.03,
			VolatilityWindow:    30 * 24 * time.Hour,
			ConfidenceThreshold: 0.6,
		},
		logger: logger,
	}
}

func (rp *RiskProfiler) UpdateRiskProfile(ctx context.Context, profile *UserBehaviorProfile, event *BehaviorEvent) error {
	// Simplified risk profiling
	if event.Type == "trade" && event.Outcome != nil {
		riskProfile := profile.RiskProfile

		// Update risk tolerance based on behavior
		if event.Outcome.Success {
			riskProfile.RiskTolerance = math.Min(1.0, riskProfile.RiskTolerance+0.01)
		} else {
			riskProfile.RiskTolerance = math.Max(0.0, riskProfile.RiskTolerance-0.02)
		}

		// Update other risk metrics
		if riskProfile.PositionSizingStyle == "" {
			riskProfile.PositionSizingStyle = "percentage"
		}

		riskProfile.StopLossUsage = 0.8 // Assume most users use stop losses
		riskProfile.TakeProfitUsage = 0.6
		riskProfile.LeverageComfort = 0.3 // Conservative default
		riskProfile.VolatilityTolerance = 0.5
		riskProfile.EmotionalStability = 0.7
		riskProfile.LossAversion = 0.6

		riskProfile.Confidence = math.Min(1.0, riskProfile.Confidence+0.03)
	}

	return nil
}

// GetBehaviorHistory retrieves behavior history for a user
func (u *UserBehaviorLearningEngine) GetBehaviorHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*BehaviorEvent, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	history, exists := u.behaviorHistory[userID]
	if !exists {
		return []*BehaviorEvent{}, nil
	}

	// Apply limit
	if limit > 0 && len(history) > limit {
		return history[len(history)-limit:], nil
	}

	return history, nil
}

// GetLearningModels retrieves learning models
func (u *UserBehaviorLearningEngine) GetLearningModels() map[string]*LearningModel {
	u.mu.RLock()
	defer u.mu.RUnlock()

	return u.learningModels
}

// UpdateRecommendationStatus updates the status of a recommendation
func (u *UserBehaviorLearningEngine) UpdateRecommendationStatus(ctx context.Context, userID uuid.UUID, recommendationID string, status string) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	profile, exists := u.userProfiles[userID]
	if !exists {
		return fmt.Errorf("user profile not found for user %s", userID)
	}

	for _, rec := range profile.Recommendations {
		if rec.ID == recommendationID {
			rec.Status = status
			u.logger.Info(ctx, "Recommendation status updated", map[string]interface{}{
				"user_id":           userID,
				"recommendation_id": recommendationID,
				"status":            status,
			})
			return nil
		}
	}

	return fmt.Errorf("recommendation not found: %s", recommendationID)
}
