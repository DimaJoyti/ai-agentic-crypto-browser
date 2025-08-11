package trading

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// BotRiskManager manages risk for individual trading bots and portfolio-level risk
type BotRiskManager struct {
	logger            *observability.Logger
	config            *BotRiskConfig
	botRiskProfiles   map[string]*BotRiskProfile
	portfolioRisk     *PortfolioRisk
	riskLimits        map[string]*RiskLimit
	riskMetrics       map[string]*BotRiskMetrics
	correlationMatrix map[string]map[string]decimal.Decimal
	alertManager      *RiskAlertManager

	// Circuit breakers
	emergencyStop bool
	tradingHalted map[string]bool

	// State management
	isRunning bool
	stopChan  chan struct{}
	mu        sync.RWMutex
}

// BotRiskConfig holds risk management configuration
type BotRiskConfig struct {
	// Portfolio-level limits
	MaxPortfolioExposure decimal.Decimal `yaml:"max_portfolio_exposure"`
	MaxCorrelationLimit  decimal.Decimal `yaml:"max_correlation_limit"`
	VaRLimit             decimal.Decimal `yaml:"var_limit"`
	MaxDrawdownLimit     decimal.Decimal `yaml:"max_drawdown_limit"`

	// Bot-level limits
	MaxBotExposure       decimal.Decimal `yaml:"max_bot_exposure"`
	MaxDailyLossPerBot   decimal.Decimal `yaml:"max_daily_loss_per_bot"`
	MaxConsecutiveLosses int             `yaml:"max_consecutive_losses"`

	// Risk monitoring
	RiskUpdateInterval time.Duration       `yaml:"risk_update_interval"`
	AlertThresholds    *BotAlertThresholds `yaml:"alert_thresholds"`

	// Emergency controls
	EmergencyStopEnabled bool `yaml:"emergency_stop_enabled"`
	AutoHaltOnViolation  bool `yaml:"auto_halt_on_violation"`
}

// BotRiskProfile defines risk parameters for individual bots
type BotRiskProfile struct {
	BotID                string          `json:"bot_id"`
	Strategy             string          `json:"strategy"`
	MaxPositionSize      decimal.Decimal `json:"max_position_size"`
	StopLoss             decimal.Decimal `json:"stop_loss"`
	TakeProfit           decimal.Decimal `json:"take_profit"`
	MaxDrawdown          decimal.Decimal `json:"max_drawdown"`
	MaxDailyLoss         decimal.Decimal `json:"max_daily_loss"`
	MaxConsecutiveLosses int             `json:"max_consecutive_losses"`
	RiskTolerance        RiskTolerance   `json:"risk_tolerance"`
	LastUpdated          time.Time       `json:"last_updated"`
}

// RiskTolerance represents risk tolerance levels
type RiskTolerance string

const (
	RiskToleranceLow        RiskTolerance = "low"
	RiskToleranceMedium     RiskTolerance = "medium"
	RiskToleranceHigh       RiskTolerance = "high"
	RiskToleranceAggressive RiskTolerance = "aggressive"
)

// PortfolioRisk tracks portfolio-level risk metrics
type PortfolioRisk struct {
	TotalExposure     decimal.Decimal            `json:"total_exposure"`
	VaR95             decimal.Decimal            `json:"var_95"`
	VaR99             decimal.Decimal            `json:"var_99"`
	ExpectedShortfall decimal.Decimal            `json:"expected_shortfall"`
	MaxDrawdown       decimal.Decimal            `json:"max_drawdown"`
	CurrentDrawdown   decimal.Decimal            `json:"current_drawdown"`
	SharpeRatio       decimal.Decimal            `json:"sharpe_ratio"`
	ConcentrationRisk map[string]decimal.Decimal `json:"concentration_risk"`
	CorrelationRisk   decimal.Decimal            `json:"correlation_risk"`
	LiquidityRisk     decimal.Decimal            `json:"liquidity_risk"`
	LastUpdated       time.Time                  `json:"last_updated"`
}

// BotRiskMetrics tracks risk metrics for individual bots
type BotRiskMetrics struct {
	BotID             string          `json:"bot_id"`
	CurrentExposure   decimal.Decimal `json:"current_exposure"`
	DailyPnL          decimal.Decimal `json:"daily_pnl"`
	UnrealizedPnL     decimal.Decimal `json:"unrealized_pnl"`
	MaxDrawdown       decimal.Decimal `json:"max_drawdown"`
	CurrentDrawdown   decimal.Decimal `json:"current_drawdown"`
	ConsecutiveLosses int             `json:"consecutive_losses"`
	ConsecutiveWins   int             `json:"consecutive_wins"`
	VaR95             decimal.Decimal `json:"var_95"`
	Volatility        decimal.Decimal `json:"volatility"`
	Beta              decimal.Decimal `json:"beta"`
	RiskScore         int             `json:"risk_score"` // 0-100
	LastUpdated       time.Time       `json:"last_updated"`
}

// RiskLimit defines specific risk limits
type RiskLimit struct {
	ID          string          `json:"id"`
	Type        RiskLimitType   `json:"type"`
	BotID       string          `json:"bot_id,omitempty"`
	Threshold   decimal.Decimal `json:"threshold"`
	Action      RiskAction      `json:"action"`
	Enabled     bool            `json:"enabled"`
	Description string          `json:"description"`
	CreatedAt   time.Time       `json:"created_at"`
}

// RiskLimitType defines types of risk limits
type RiskLimitType string

const (
	RiskLimitTypePosition      RiskLimitType = "position"
	RiskLimitTypeDrawdown      RiskLimitType = "drawdown"
	RiskLimitTypeDailyLoss     RiskLimitType = "daily_loss"
	RiskLimitTypeVaR           RiskLimitType = "var"
	RiskLimitTypeConcentration RiskLimitType = "concentration"
	RiskLimitTypeCorrelation   RiskLimitType = "correlation"
)

// RiskAction defines actions to take when limits are breached
type RiskAction string

const (
	RiskActionAlert         RiskAction = "alert"
	RiskActionReduceSize    RiskAction = "reduce_size"
	RiskActionHaltBot       RiskAction = "halt_bot"
	RiskActionEmergencyStop RiskAction = "emergency_stop"
)

// BotAlertThresholds defines thresholds for risk alerts
type BotAlertThresholds struct {
	VaRWarning       decimal.Decimal `yaml:"var_warning"`
	VaRCritical      decimal.Decimal `yaml:"var_critical"`
	DrawdownWarning  decimal.Decimal `yaml:"drawdown_warning"`
	DrawdownCritical decimal.Decimal `yaml:"drawdown_critical"`
	LossWarning      decimal.Decimal `yaml:"loss_warning"`
	LossCritical     decimal.Decimal `yaml:"loss_critical"`
}

// NewBotRiskManager creates a new bot risk manager
func NewBotRiskManager(logger *observability.Logger) *BotRiskManager {
	config := &BotRiskConfig{
		MaxPortfolioExposure: decimal.NewFromFloat(0.80), // 80% max exposure
		MaxCorrelationLimit:  decimal.NewFromFloat(0.70), // 70% max correlation
		VaRLimit:             decimal.NewFromFloat(0.05), // 5% VaR limit
		MaxDrawdownLimit:     decimal.NewFromFloat(0.15), // 15% max drawdown
		MaxBotExposure:       decimal.NewFromFloat(0.20), // 20% max per bot
		MaxDailyLossPerBot:   decimal.NewFromFloat(0.05), // 5% daily loss per bot
		MaxConsecutiveLosses: 5,
		RiskUpdateInterval:   time.Minute,
		EmergencyStopEnabled: true,
		AutoHaltOnViolation:  true,
		AlertThresholds: &BotAlertThresholds{
			VaRWarning:       decimal.NewFromFloat(0.03),
			VaRCritical:      decimal.NewFromFloat(0.05),
			DrawdownWarning:  decimal.NewFromFloat(0.10),
			DrawdownCritical: decimal.NewFromFloat(0.15),
			LossWarning:      decimal.NewFromFloat(0.03),
			LossCritical:     decimal.NewFromFloat(0.05),
		},
	}

	return &BotRiskManager{
		logger:            logger,
		config:            config,
		botRiskProfiles:   make(map[string]*BotRiskProfile),
		portfolioRisk:     &PortfolioRisk{ConcentrationRisk: make(map[string]decimal.Decimal)},
		riskLimits:        make(map[string]*RiskLimit),
		riskMetrics:       make(map[string]*BotRiskMetrics),
		correlationMatrix: make(map[string]map[string]decimal.Decimal),
		tradingHalted:     make(map[string]bool),
		alertManager:      NewRiskAlertManager(logger),
		stopChan:          make(chan struct{}),
	}
}

// Start starts the risk manager
func (brm *BotRiskManager) Start(ctx context.Context) error {
	brm.mu.Lock()
	defer brm.mu.Unlock()

	if brm.isRunning {
		return fmt.Errorf("bot risk manager is already running")
	}

	brm.isRunning = true

	// Initialize default risk limits
	brm.initializeDefaultLimits()

	// Start risk monitoring loop
	go brm.riskMonitoringLoop(ctx)

	brm.logger.Info(ctx, "Bot risk manager started", map[string]interface{}{
		"max_portfolio_exposure": brm.config.MaxPortfolioExposure.String(),
		"var_limit":              brm.config.VaRLimit.String(),
		"max_drawdown_limit":     brm.config.MaxDrawdownLimit.String(),
	})

	return nil
}

// Stop stops the risk manager
func (brm *BotRiskManager) Stop(ctx context.Context) error {
	brm.mu.Lock()
	defer brm.mu.Unlock()

	if !brm.isRunning {
		return nil
	}

	brm.isRunning = false
	close(brm.stopChan)

	brm.logger.Info(ctx, "Bot risk manager stopped", nil)
	return nil
}

// RegisterBot registers a bot with the risk manager
func (brm *BotRiskManager) RegisterBot(botID string, strategy string, riskProfile *BotRiskProfile) error {
	brm.mu.Lock()
	defer brm.mu.Unlock()

	if riskProfile == nil {
		riskProfile = brm.createDefaultRiskProfile(botID, strategy)
	}

	riskProfile.BotID = botID
	riskProfile.Strategy = strategy
	riskProfile.LastUpdated = time.Now()

	brm.botRiskProfiles[botID] = riskProfile
	brm.riskMetrics[botID] = &BotRiskMetrics{
		BotID:       botID,
		LastUpdated: time.Now(),
	}

	brm.logger.Info(context.Background(), "Bot registered with risk manager", map[string]interface{}{
		"bot_id":         botID,
		"strategy":       strategy,
		"risk_tolerance": string(riskProfile.RiskTolerance),
	})

	return nil
}

// ValidateOrder validates an order against risk limits
func (brm *BotRiskManager) ValidateOrder(ctx context.Context, botID string, order *OrderRequest) error {
	brm.mu.RLock()
	defer brm.mu.RUnlock()

	// Check if trading is halted
	if brm.emergencyStop {
		return fmt.Errorf("emergency stop activated - all trading halted")
	}

	if brm.tradingHalted[botID] {
		return fmt.Errorf("trading halted for bot %s", botID)
	}

	// Check bot-specific limits
	if err := brm.validateBotLimits(botID, order); err != nil {
		return fmt.Errorf("bot limit violation: %w", err)
	}

	// Check portfolio-level limits
	if err := brm.validatePortfolioLimits(order); err != nil {
		return fmt.Errorf("portfolio limit violation: %w", err)
	}

	// Check correlation limits
	if err := brm.validateCorrelationLimits(botID, order); err != nil {
		return fmt.Errorf("correlation limit violation: %w", err)
	}

	return nil
}

// UpdateBotMetrics updates risk metrics for a bot
func (brm *BotRiskManager) UpdateBotMetrics(botID string, metrics *BotRiskMetrics) error {
	brm.mu.Lock()
	defer brm.mu.Unlock()

	metrics.BotID = botID
	metrics.LastUpdated = time.Now()
	brm.riskMetrics[botID] = metrics

	// Calculate risk score
	metrics.RiskScore = brm.calculateRiskScore(metrics)

	// Update portfolio risk
	brm.updatePortfolioRisk()

	// Check for risk violations
	brm.checkRiskViolations(context.Background(), botID, metrics)

	return nil
}

// GetBotRiskMetrics returns risk metrics for a bot
func (brm *BotRiskManager) GetBotRiskMetrics(botID string) (*BotRiskMetrics, error) {
	brm.mu.RLock()
	defer brm.mu.RUnlock()

	metrics, exists := brm.riskMetrics[botID]
	if !exists {
		return nil, fmt.Errorf("risk metrics not found for bot %s", botID)
	}

	return metrics, nil
}

// GetPortfolioRisk returns portfolio-level risk metrics
func (brm *BotRiskManager) GetPortfolioRisk() *PortfolioRisk {
	brm.mu.RLock()
	defer brm.mu.RUnlock()

	return brm.portfolioRisk
}

// HaltBot halts trading for a specific bot
func (brm *BotRiskManager) HaltBot(ctx context.Context, botID string, reason string) error {
	brm.mu.Lock()
	defer brm.mu.Unlock()

	brm.tradingHalted[botID] = true

	brm.logger.Error(ctx, "Bot trading halted", fmt.Errorf(reason), map[string]interface{}{
		"bot_id": botID,
		"reason": reason,
	})

	// Send alert
	brm.alertManager.SendAlert(ctx, &RiskAlert{
		Type:     RiskAlertTypeBotHalted,
		Severity: AlertSeverityCritical,
		BotID:    botID,
		Message:  fmt.Sprintf("Bot %s halted: %s", botID, reason),
	})

	return nil
}

// ResumeBot resumes trading for a specific bot
func (brm *BotRiskManager) ResumeBot(ctx context.Context, botID string) error {
	brm.mu.Lock()
	defer brm.mu.Unlock()

	brm.tradingHalted[botID] = false

	brm.logger.Info(ctx, "Bot trading resumed", map[string]interface{}{
		"bot_id": botID,
	})

	return nil
}

// EmergencyStop activates emergency stop for all bots
func (brm *BotRiskManager) EmergencyStop(ctx context.Context, reason string) error {
	brm.mu.Lock()
	defer brm.mu.Unlock()

	brm.emergencyStop = true

	brm.logger.Error(ctx, "Emergency stop activated", fmt.Errorf(reason), map[string]interface{}{
		"reason": reason,
	})

	// Send critical alert
	brm.alertManager.SendAlert(ctx, &RiskAlert{
		Type:     RiskAlertTypeEmergencyStop,
		Severity: AlertSeverityCritical,
		Message:  fmt.Sprintf("Emergency stop activated: %s", reason),
	})

	return nil
}

// OrderRequest represents an order validation request
type OrderRequest struct {
	Symbol    string          `json:"symbol"`
	Side      string          `json:"side"`
	Amount    decimal.Decimal `json:"amount"`
	Price     decimal.Decimal `json:"price"`
	OrderType string          `json:"order_type"`
}

// createDefaultRiskProfile creates a default risk profile for a bot
func (brm *BotRiskManager) createDefaultRiskProfile(botID, strategy string) *BotRiskProfile {
	var riskTolerance RiskTolerance
	var maxPositionSize, stopLoss, takeProfit, maxDrawdown, maxDailyLoss decimal.Decimal
	var maxConsecutiveLosses int

	// Set defaults based on strategy
	switch strategy {
	case "dca":
		riskTolerance = RiskToleranceLow
		maxPositionSize = decimal.NewFromFloat(0.20)
		stopLoss = decimal.NewFromFloat(0.15)
		takeProfit = decimal.NewFromFloat(0.30)
		maxDrawdown = decimal.NewFromFloat(0.10)
		maxDailyLoss = decimal.NewFromFloat(0.03)
		maxConsecutiveLosses = 3
	case "grid":
		riskTolerance = RiskToleranceMedium
		maxPositionSize = decimal.NewFromFloat(0.15)
		stopLoss = decimal.NewFromFloat(0.20)
		takeProfit = decimal.NewFromFloat(0.25)
		maxDrawdown = decimal.NewFromFloat(0.12)
		maxDailyLoss = decimal.NewFromFloat(0.04)
		maxConsecutiveLosses = 4
	case "momentum":
		riskTolerance = RiskToleranceHigh
		maxPositionSize = decimal.NewFromFloat(0.10)
		stopLoss = decimal.NewFromFloat(0.08)
		takeProfit = decimal.NewFromFloat(0.15)
		maxDrawdown = decimal.NewFromFloat(0.15)
		maxDailyLoss = decimal.NewFromFloat(0.05)
		maxConsecutiveLosses = 5
	case "scalping":
		riskTolerance = RiskToleranceAggressive
		maxPositionSize = decimal.NewFromFloat(0.10)
		stopLoss = decimal.NewFromFloat(0.003)
		takeProfit = decimal.NewFromFloat(0.002)
		maxDrawdown = decimal.NewFromFloat(0.15)
		maxDailyLoss = decimal.NewFromFloat(0.06)
		maxConsecutiveLosses = 8
	default:
		riskTolerance = RiskToleranceMedium
		maxPositionSize = decimal.NewFromFloat(0.15)
		stopLoss = decimal.NewFromFloat(0.10)
		takeProfit = decimal.NewFromFloat(0.20)
		maxDrawdown = decimal.NewFromFloat(0.12)
		maxDailyLoss = decimal.NewFromFloat(0.04)
		maxConsecutiveLosses = 4
	}

	return &BotRiskProfile{
		BotID:                botID,
		Strategy:             strategy,
		MaxPositionSize:      maxPositionSize,
		StopLoss:             stopLoss,
		TakeProfit:           takeProfit,
		MaxDrawdown:          maxDrawdown,
		MaxDailyLoss:         maxDailyLoss,
		MaxConsecutiveLosses: maxConsecutiveLosses,
		RiskTolerance:        riskTolerance,
		LastUpdated:          time.Now(),
	}
}

// initializeDefaultLimits sets up default risk limits
func (brm *BotRiskManager) initializeDefaultLimits() {
	limits := []*RiskLimit{
		{
			ID:          uuid.New().String(),
			Type:        RiskLimitTypeDrawdown,
			Threshold:   brm.config.MaxDrawdownLimit,
			Action:      RiskActionEmergencyStop,
			Enabled:     true,
			Description: "Portfolio maximum drawdown limit",
			CreatedAt:   time.Now(),
		},
		{
			ID:          uuid.New().String(),
			Type:        RiskLimitTypeVaR,
			Threshold:   brm.config.VaRLimit,
			Action:      RiskActionAlert,
			Enabled:     true,
			Description: "Portfolio VaR limit",
			CreatedAt:   time.Now(),
		},
		{
			ID:          uuid.New().String(),
			Type:        RiskLimitTypeConcentration,
			Threshold:   brm.config.MaxBotExposure,
			Action:      RiskActionReduceSize,
			Enabled:     true,
			Description: "Maximum bot exposure limit",
			CreatedAt:   time.Now(),
		},
	}

	for _, limit := range limits {
		brm.riskLimits[limit.ID] = limit
	}
}

// validateBotLimits validates order against bot-specific limits
func (brm *BotRiskManager) validateBotLimits(botID string, order *OrderRequest) error {
	profile, exists := brm.botRiskProfiles[botID]
	if !exists {
		return fmt.Errorf("risk profile not found for bot %s", botID)
	}

	metrics, exists := brm.riskMetrics[botID]
	if !exists {
		return fmt.Errorf("risk metrics not found for bot %s", botID)
	}

	// Check position size limit
	orderValue := order.Amount.Mul(order.Price)
	if orderValue.GreaterThan(profile.MaxPositionSize.Mul(brm.portfolioRisk.TotalExposure)) {
		return fmt.Errorf("order exceeds max position size: %s > %s",
			orderValue.String(), profile.MaxPositionSize.String())
	}

	// Check daily loss limit
	if metrics.DailyPnL.LessThan(profile.MaxDailyLoss.Neg()) {
		return fmt.Errorf("daily loss limit exceeded: %s < %s",
			metrics.DailyPnL.String(), profile.MaxDailyLoss.Neg().String())
	}

	// Check consecutive losses
	if metrics.ConsecutiveLosses >= profile.MaxConsecutiveLosses {
		return fmt.Errorf("consecutive losses limit exceeded: %d >= %d",
			metrics.ConsecutiveLosses, profile.MaxConsecutiveLosses)
	}

	// Check drawdown limit
	if metrics.CurrentDrawdown.GreaterThan(profile.MaxDrawdown) {
		return fmt.Errorf("drawdown limit exceeded: %s > %s",
			metrics.CurrentDrawdown.String(), profile.MaxDrawdown.String())
	}

	return nil
}

// validatePortfolioLimits validates order against portfolio limits
func (brm *BotRiskManager) validatePortfolioLimits(order *OrderRequest) error {
	orderValue := order.Amount.Mul(order.Price)

	// Check total exposure limit
	newExposure := brm.portfolioRisk.TotalExposure.Add(orderValue)
	if newExposure.GreaterThan(brm.config.MaxPortfolioExposure.Mul(brm.portfolioRisk.TotalExposure)) {
		return fmt.Errorf("portfolio exposure limit exceeded")
	}

	// Check VaR limit
	if brm.portfolioRisk.VaR95.GreaterThan(brm.config.VaRLimit.Mul(brm.portfolioRisk.TotalExposure)) {
		return fmt.Errorf("VaR limit exceeded")
	}

	// Check drawdown limit
	if brm.portfolioRisk.CurrentDrawdown.GreaterThan(brm.config.MaxDrawdownLimit) {
		return fmt.Errorf("portfolio drawdown limit exceeded")
	}

	return nil
}

// validateCorrelationLimits validates order against correlation limits
func (brm *BotRiskManager) validateCorrelationLimits(botID string, order *OrderRequest) error {
	// Check correlation with other positions
	correlations, exists := brm.correlationMatrix[order.Symbol]
	if !exists {
		return nil // No correlation data available
	}

	for symbol, correlation := range correlations {
		if correlation.GreaterThan(brm.config.MaxCorrelationLimit) {
			// Check if we have exposure to the correlated symbol
			for _, metrics := range brm.riskMetrics {
				if metrics.CurrentExposure.GreaterThan(decimal.Zero) {
					return fmt.Errorf("correlation limit exceeded with %s: %s > %s",
						symbol, correlation.String(), brm.config.MaxCorrelationLimit.String())
				}
			}
		}
	}

	return nil
}

// riskMonitoringLoop continuously monitors risk metrics
func (brm *BotRiskManager) riskMonitoringLoop(ctx context.Context) {
	ticker := time.NewTicker(brm.config.RiskUpdateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-brm.stopChan:
			return
		case <-ticker.C:
			brm.updatePortfolioRisk()
			brm.checkAllRiskViolations(ctx)
		}
	}
}

// updatePortfolioRisk updates portfolio-level risk metrics
func (brm *BotRiskManager) updatePortfolioRisk() {
	brm.mu.Lock()
	defer brm.mu.Unlock()

	var totalExposure decimal.Decimal
	var totalVaR95 decimal.Decimal
	var maxDrawdown decimal.Decimal

	// Aggregate metrics from all bots
	for _, metrics := range brm.riskMetrics {
		totalExposure = totalExposure.Add(metrics.CurrentExposure)
		totalVaR95 = totalVaR95.Add(metrics.VaR95)

		if metrics.CurrentDrawdown.GreaterThan(maxDrawdown) {
			maxDrawdown = metrics.CurrentDrawdown
		}
	}

	// Update portfolio risk
	brm.portfolioRisk.TotalExposure = totalExposure
	brm.portfolioRisk.VaR95 = totalVaR95
	brm.portfolioRisk.VaR99 = totalVaR95.Mul(decimal.NewFromFloat(1.3)) // Approximate scaling
	brm.portfolioRisk.CurrentDrawdown = maxDrawdown
	brm.portfolioRisk.LastUpdated = time.Now()

	// Calculate concentration risk
	brm.calculateConcentrationRisk()

	// Calculate correlation risk
	brm.calculateCorrelationRisk()
}

// calculateConcentrationRisk calculates concentration risk
func (brm *BotRiskManager) calculateConcentrationRisk() {
	if brm.portfolioRisk.TotalExposure.IsZero() {
		return
	}

	concentrationRisk := make(map[string]decimal.Decimal)

	for botID, metrics := range brm.riskMetrics {
		concentration := metrics.CurrentExposure.Div(brm.portfolioRisk.TotalExposure)
		concentrationRisk[botID] = concentration
	}

	brm.portfolioRisk.ConcentrationRisk = concentrationRisk
}

// calculateCorrelationRisk calculates correlation risk
func (brm *BotRiskManager) calculateCorrelationRisk() {
	// Simplified correlation risk calculation
	// In practice, this would use historical price correlations
	var correlationRisk decimal.Decimal

	for _, correlations := range brm.correlationMatrix {
		for _, correlation := range correlations {
			if correlation.GreaterThan(brm.config.MaxCorrelationLimit) {
				correlationRisk = correlationRisk.Add(correlation.Sub(brm.config.MaxCorrelationLimit))
			}
		}
	}

	brm.portfolioRisk.CorrelationRisk = correlationRisk
}

// calculateRiskScore calculates a risk score for a bot (0-100)
func (brm *BotRiskManager) calculateRiskScore(metrics *BotRiskMetrics) int {
	score := 0

	// Drawdown component (0-30 points)
	if !metrics.CurrentDrawdown.IsZero() {
		drawdownScore := int(metrics.CurrentDrawdown.Mul(decimal.NewFromFloat(100)).IntPart())
		if drawdownScore > 30 {
			drawdownScore = 30
		}
		score += drawdownScore
	}

	// Volatility component (0-25 points)
	if !metrics.Volatility.IsZero() {
		volatilityScore := int(metrics.Volatility.Mul(decimal.NewFromFloat(50)).IntPart())
		if volatilityScore > 25 {
			volatilityScore = 25
		}
		score += volatilityScore
	}

	// Consecutive losses component (0-20 points)
	lossScore := metrics.ConsecutiveLosses * 4
	if lossScore > 20 {
		lossScore = 20
	}
	score += lossScore

	// VaR component (0-25 points)
	if !metrics.VaR95.IsZero() && !metrics.CurrentExposure.IsZero() {
		varRatio := metrics.VaR95.Div(metrics.CurrentExposure)
		varScore := int(varRatio.Mul(decimal.NewFromFloat(100)).IntPart())
		if varScore > 25 {
			varScore = 25
		}
		score += varScore
	}

	if score > 100 {
		score = 100
	}

	return score
}

// checkAllRiskViolations checks for risk violations across all bots
func (brm *BotRiskManager) checkAllRiskViolations(ctx context.Context) {
	brm.mu.RLock()
	defer brm.mu.RUnlock()

	// Check portfolio-level violations
	brm.checkPortfolioViolations(ctx)

	// Check individual bot violations
	for botID, metrics := range brm.riskMetrics {
		brm.checkRiskViolations(ctx, botID, metrics)
	}
}

// checkPortfolioViolations checks for portfolio-level risk violations
func (brm *BotRiskManager) checkPortfolioViolations(ctx context.Context) {
	// Check VaR limit
	if brm.portfolioRisk.VaR95.GreaterThan(brm.config.VaRLimit.Mul(brm.portfolioRisk.TotalExposure)) {
		brm.alertManager.SendAlert(ctx, &RiskAlert{
			Type:        RiskAlertTypeVaR,
			Severity:    AlertSeverityCritical,
			Message:     "Portfolio VaR limit exceeded",
			Threshold:   brm.config.VaRLimit,
			ActualValue: brm.portfolioRisk.VaR95.Div(brm.portfolioRisk.TotalExposure),
		})

		if brm.config.AutoHaltOnViolation {
			brm.EmergencyStop(ctx, "Portfolio VaR limit exceeded")
		}
	}

	// Check drawdown limit
	if brm.portfolioRisk.CurrentDrawdown.GreaterThan(brm.config.MaxDrawdownLimit) {
		brm.alertManager.SendAlert(ctx, &RiskAlert{
			Type:        RiskAlertTypeDrawdown,
			Severity:    AlertSeverityCritical,
			Message:     "Portfolio drawdown limit exceeded",
			Threshold:   brm.config.MaxDrawdownLimit,
			ActualValue: brm.portfolioRisk.CurrentDrawdown,
		})

		if brm.config.AutoHaltOnViolation {
			brm.EmergencyStop(ctx, "Portfolio drawdown limit exceeded")
		}
	}
}

// checkRiskViolations checks for risk violations for a specific bot
func (brm *BotRiskManager) checkRiskViolations(ctx context.Context, botID string, metrics *BotRiskMetrics) {
	profile, exists := brm.botRiskProfiles[botID]
	if !exists {
		return
	}

	// Check drawdown violation
	if metrics.CurrentDrawdown.GreaterThan(profile.MaxDrawdown) {
		brm.alertManager.SendAlert(ctx, &RiskAlert{
			Type:        RiskAlertTypeDrawdown,
			Severity:    AlertSeverityHigh,
			BotID:       botID,
			Message:     fmt.Sprintf("Bot %s drawdown limit exceeded", botID),
			Threshold:   profile.MaxDrawdown,
			ActualValue: metrics.CurrentDrawdown,
		})

		if brm.config.AutoHaltOnViolation {
			brm.HaltBot(ctx, botID, "Drawdown limit exceeded")
		}
	}

	// Check daily loss violation
	if metrics.DailyPnL.LessThan(profile.MaxDailyLoss.Neg()) {
		brm.alertManager.SendAlert(ctx, &RiskAlert{
			Type:        RiskAlertTypeDailyLoss,
			Severity:    AlertSeverityHigh,
			BotID:       botID,
			Message:     fmt.Sprintf("Bot %s daily loss limit exceeded", botID),
			Threshold:   profile.MaxDailyLoss,
			ActualValue: metrics.DailyPnL.Abs(),
		})

		if brm.config.AutoHaltOnViolation {
			brm.HaltBot(ctx, botID, "Daily loss limit exceeded")
		}
	}

	// Check consecutive losses
	if metrics.ConsecutiveLosses >= profile.MaxConsecutiveLosses {
		brm.alertManager.SendAlert(ctx, &RiskAlert{
			Type:     RiskAlertTypePosition,
			Severity: AlertSeverityWarning,
			BotID:    botID,
			Message:  fmt.Sprintf("Bot %s has %d consecutive losses", botID, metrics.ConsecutiveLosses),
		})
	}
}

// GetAlertsBySeverity returns alerts by severity level
func (brm *BotRiskManager) GetAlertsBySeverity(severity AlertSeverity) []*RiskAlert {
	return brm.alertManager.GetAlertsBySeverity(severity)
}

// GetAlertsByBot returns alerts for a specific bot
func (brm *BotRiskManager) GetAlertsByBot(botID string) []*RiskAlert {
	return brm.alertManager.GetAlertsByBot(botID)
}

// GetActiveAlerts returns all active alerts
func (brm *BotRiskManager) GetActiveAlerts() []*RiskAlert {
	return brm.alertManager.GetActiveAlerts()
}

// AcknowledgeAlert acknowledges an alert
func (brm *BotRiskManager) AcknowledgeAlert(ctx context.Context, alertID string, userID string) error {
	return brm.alertManager.AcknowledgeAlert(ctx, alertID, userID)
}

// ResolveAlert resolves an alert
func (brm *BotRiskManager) ResolveAlert(ctx context.Context, alertID string, userID string) error {
	return brm.alertManager.ResolveAlert(ctx, alertID, userID)
}

// GetRiskStatistics returns comprehensive risk statistics
func (brm *BotRiskManager) GetRiskStatistics() map[string]interface{} {
	brm.mu.RLock()
	defer brm.mu.RUnlock()

	stats := map[string]interface{}{
		"portfolio_risk":   brm.portfolioRisk,
		"total_bots":       len(brm.botRiskProfiles),
		"active_bots":      0,
		"halted_bots":      0,
		"emergency_stop":   brm.emergencyStop,
		"alert_statistics": brm.alertManager.GetAlertStatistics(),
	}

	// Count active and halted bots
	for botID := range brm.botRiskProfiles {
		if brm.tradingHalted[botID] {
			stats["halted_bots"] = stats["halted_bots"].(int) + 1
		} else {
			stats["active_bots"] = stats["active_bots"].(int) + 1
		}
	}

	return stats
}

// UpdateCorrelationMatrix updates the correlation matrix
func (brm *BotRiskManager) UpdateCorrelationMatrix(correlations map[string]map[string]decimal.Decimal) {
	brm.mu.Lock()
	defer brm.mu.Unlock()

	brm.correlationMatrix = correlations
}

// GetCorrelationMatrix returns the current correlation matrix
func (brm *BotRiskManager) GetCorrelationMatrix() map[string]map[string]decimal.Decimal {
	brm.mu.RLock()
	defer brm.mu.RUnlock()

	return brm.correlationMatrix
}
