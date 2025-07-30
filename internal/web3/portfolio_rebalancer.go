package web3

import (
	"context"
	"fmt"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// PortfolioRebalancer handles automated portfolio rebalancing
type PortfolioRebalancer struct {
	logger         *observability.Logger
	tradingEngine  *TradingEngine
	defiManager    *DeFiProtocolManager
	rebalanceRules map[uuid.UUID]*RebalanceStrategy
	config         RebalancerConfig
}

// RebalancerConfig holds configuration for portfolio rebalancing
type RebalancerConfig struct {
	RebalanceInterval     time.Duration   `json:"rebalance_interval"`
	DriftThreshold        decimal.Decimal `json:"drift_threshold"`       // % deviation from target
	MinRebalanceAmount    decimal.Decimal `json:"min_rebalance_amount"`  // Minimum amount to trigger rebalance
	MaxTransactionCost    decimal.Decimal `json:"max_transaction_cost"`  // Max cost as % of rebalance amount
	VolatilityWindow      time.Duration   `json:"volatility_window"`     // Window for volatility calculation
	CorrelationThreshold  decimal.Decimal `json:"correlation_threshold"` // Asset correlation threshold
	EnableTaxOptimization bool            `json:"enable_tax_optimization"`
	TaxLossHarvestingMin  decimal.Decimal `json:"tax_loss_harvesting_min"`
}

// RebalanceStrategy defines how a portfolio should be rebalanced
type RebalanceStrategy struct {
	ID                uuid.UUID                  `json:"id"`
	PortfolioID       uuid.UUID                  `json:"portfolio_id"`
	Name              string                     `json:"name"`
	Type              RebalanceType              `json:"type"`
	TargetAllocations map[string]decimal.Decimal `json:"target_allocations"` // token -> percentage
	Constraints       []AllocationConstraint     `json:"constraints"`
	TriggerConditions []RebalanceTrigger         `json:"trigger_conditions"`
	IsActive          bool                       `json:"is_active"`
	LastRebalance     time.Time                  `json:"last_rebalance"`
	CreatedAt         time.Time                  `json:"created_at"`
	Metadata          map[string]interface{}     `json:"metadata"`
}

// RebalanceType represents different rebalancing strategies
type RebalanceType string

const (
	RebalanceTypeFixed      RebalanceType = "fixed"       // Fixed percentage allocation
	RebalanceTypeDynamic    RebalanceType = "dynamic"     // Dynamic based on market conditions
	RebalanceTypeRiskParity RebalanceType = "risk_parity" // Risk-weighted allocation
	RebalanceTypeMomentum   RebalanceType = "momentum"    // Momentum-based allocation
	RebalanceTypeMeanRevert RebalanceType = "mean_revert" // Mean reversion allocation
)

// AllocationConstraint defines constraints on asset allocation
type AllocationConstraint struct {
	Asset       string          `json:"asset"`
	MinWeight   decimal.Decimal `json:"min_weight"`
	MaxWeight   decimal.Decimal `json:"max_weight"`
	MaxDrawdown decimal.Decimal `json:"max_drawdown"`
}

// RebalanceTrigger defines conditions that trigger rebalancing
type RebalanceTrigger struct {
	Type      TriggerType     `json:"type"`
	Threshold decimal.Decimal `json:"threshold"`
	Condition string          `json:"condition"` // "greater_than", "less_than", "deviation"
}

// TriggerType represents different trigger types
type TriggerType string

const (
	TriggerTypeDrift       TriggerType = "drift"       // Allocation drift
	TriggerTypeVolatility  TriggerType = "volatility"  // Volatility spike
	TriggerTypeCorrelation TriggerType = "correlation" // Correlation change
	TriggerTypeTime        TriggerType = "time"        // Time-based
	TriggerTypeDrawdown    TriggerType = "drawdown"    // Portfolio drawdown
	TriggerTypeProfit      TriggerType = "profit"      // Profit taking
)

// RebalanceAction represents an action to take during rebalancing
type RebalanceAction struct {
	ID           uuid.UUID       `json:"id"`
	PortfolioID  uuid.UUID       `json:"portfolio_id"`
	ActionType   ActionType      `json:"action_type"`
	FromAsset    string          `json:"from_asset"`
	ToAsset      string          `json:"to_asset"`
	Amount       decimal.Decimal `json:"amount"`
	ExpectedCost decimal.Decimal `json:"expected_cost"`
	Priority     int             `json:"priority"`
	Reason       string          `json:"reason"`
	CreatedAt    time.Time       `json:"created_at"`
}

// ActionType represents different rebalancing actions
type ActionType string

const (
	ActionTypeBuy     ActionType = "buy"
	ActionTypeSell    ActionType = "sell"
	ActionTypeSwap    ActionType = "swap"
	ActionTypeStake   ActionType = "stake"
	ActionTypeUnstake ActionType = "unstake"
)

// NewPortfolioRebalancer creates a new portfolio rebalancer
func NewPortfolioRebalancer(
	logger *observability.Logger,
	tradingEngine *TradingEngine,
	defiManager *DeFiProtocolManager,
) *PortfolioRebalancer {
	config := RebalancerConfig{
		RebalanceInterval:     6 * time.Hour,              // Rebalance every 6 hours
		DriftThreshold:        decimal.NewFromFloat(0.05), // 5% drift threshold
		MinRebalanceAmount:    decimal.NewFromInt(100),    // $100 minimum
		MaxTransactionCost:    decimal.NewFromFloat(0.02), // 2% max transaction cost
		VolatilityWindow:      24 * time.Hour,             // 24-hour volatility window
		CorrelationThreshold:  decimal.NewFromFloat(0.8),  // 80% correlation threshold
		EnableTaxOptimization: true,
		TaxLossHarvestingMin:  decimal.NewFromFloat(0.03), // 3% minimum loss for harvesting
	}

	return &PortfolioRebalancer{
		logger:         logger,
		tradingEngine:  tradingEngine,
		defiManager:    defiManager,
		rebalanceRules: make(map[uuid.UUID]*RebalanceStrategy),
		config:         config,
	}
}

// CreateRebalanceStrategy creates a new rebalancing strategy
func (r *PortfolioRebalancer) CreateRebalanceStrategy(
	ctx context.Context,
	portfolioID uuid.UUID,
	name string,
	strategyType RebalanceType,
	targetAllocations map[string]decimal.Decimal,
) (*RebalanceStrategy, error) {

	// Validate target allocations sum to 100%
	totalAllocation := decimal.Zero
	for _, allocation := range targetAllocations {
		totalAllocation = totalAllocation.Add(allocation)
	}

	if !totalAllocation.Equal(decimal.NewFromInt(1)) {
		return nil, fmt.Errorf("target allocations must sum to 100%%, got %s", totalAllocation.Mul(decimal.NewFromInt(100)).String())
	}

	strategy := &RebalanceStrategy{
		ID:                uuid.New(),
		PortfolioID:       portfolioID,
		Name:              name,
		Type:              strategyType,
		TargetAllocations: targetAllocations,
		Constraints:       []AllocationConstraint{},
		TriggerConditions: r.getDefaultTriggers(strategyType),
		IsActive:          true,
		LastRebalance:     time.Time{},
		CreatedAt:         time.Now(),
		Metadata:          make(map[string]interface{}),
	}

	r.rebalanceRules[portfolioID] = strategy

	r.logger.Info(ctx, "Rebalance strategy created", map[string]interface{}{
		"strategy_id":   strategy.ID.String(),
		"portfolio_id":  portfolioID.String(),
		"strategy_type": string(strategyType),
		"allocations":   targetAllocations,
	})

	return strategy, nil
}

// getDefaultTriggers returns default triggers for a strategy type
func (r *PortfolioRebalancer) getDefaultTriggers(strategyType RebalanceType) []RebalanceTrigger {
	switch strategyType {
	case RebalanceTypeFixed:
		return []RebalanceTrigger{
			{Type: TriggerTypeDrift, Threshold: decimal.NewFromFloat(0.05), Condition: "deviation"},
			{Type: TriggerTypeTime, Threshold: decimal.NewFromFloat(24), Condition: "hours"},
		}
	case RebalanceTypeDynamic:
		return []RebalanceTrigger{
			{Type: TriggerTypeDrift, Threshold: decimal.NewFromFloat(0.03), Condition: "deviation"},
			{Type: TriggerTypeVolatility, Threshold: decimal.NewFromFloat(0.2), Condition: "greater_than"},
			{Type: TriggerTypeCorrelation, Threshold: decimal.NewFromFloat(0.8), Condition: "greater_than"},
		}
	case RebalanceTypeRiskParity:
		return []RebalanceTrigger{
			{Type: TriggerTypeDrift, Threshold: decimal.NewFromFloat(0.04), Condition: "deviation"},
			{Type: TriggerTypeVolatility, Threshold: decimal.NewFromFloat(0.15), Condition: "greater_than"},
		}
	case RebalanceTypeMomentum:
		return []RebalanceTrigger{
			{Type: TriggerTypeDrift, Threshold: decimal.NewFromFloat(0.1), Condition: "deviation"},
			{Type: TriggerTypeProfit, Threshold: decimal.NewFromFloat(0.1), Condition: "greater_than"},
		}
	case RebalanceTypeMeanRevert:
		return []RebalanceTrigger{
			{Type: TriggerTypeDrift, Threshold: decimal.NewFromFloat(0.08), Condition: "deviation"},
			{Type: TriggerTypeDrawdown, Threshold: decimal.NewFromFloat(0.1), Condition: "greater_than"},
		}
	default:
		return []RebalanceTrigger{
			{Type: TriggerTypeDrift, Threshold: decimal.NewFromFloat(0.05), Condition: "deviation"},
		}
	}
}

// RebalancePortfolio performs portfolio rebalancing
func (r *PortfolioRebalancer) RebalancePortfolio(ctx context.Context, portfolioID uuid.UUID) error {
	strategy, exists := r.rebalanceRules[portfolioID]
	if !exists {
		return fmt.Errorf("no rebalance strategy found for portfolio: %s", portfolioID.String())
	}

	if !strategy.IsActive {
		return fmt.Errorf("rebalance strategy is not active")
	}

	// Get current portfolio
	portfolio, err := r.tradingEngine.GetPortfolio(portfolioID)
	if err != nil {
		return fmt.Errorf("failed to get portfolio: %w", err)
	}

	// Check if rebalancing is needed
	shouldRebalance, triggers := r.shouldRebalance(ctx, portfolio, strategy)
	if !shouldRebalance {
		return nil
	}

	r.logger.Info(ctx, "Starting portfolio rebalance", map[string]interface{}{
		"portfolio_id": portfolioID.String(),
		"triggers":     triggers,
	})

	// Calculate current allocations
	currentAllocations := r.calculateCurrentAllocations(portfolio)

	// Generate rebalance actions
	actions := r.generateRebalanceActions(ctx, portfolio, strategy, currentAllocations)

	// Execute rebalance actions
	for _, action := range actions {
		if err := r.executeRebalanceAction(ctx, action); err != nil {
			r.logger.Error(ctx, "Failed to execute rebalance action", err)
			continue
		}
	}

	// Update last rebalance time
	strategy.LastRebalance = time.Now()

	r.logger.Info(ctx, "Portfolio rebalance completed", map[string]interface{}{
		"portfolio_id":     portfolioID.String(),
		"actions_executed": len(actions),
	})

	return nil
}

// shouldRebalance determines if portfolio should be rebalanced
func (r *PortfolioRebalancer) shouldRebalance(ctx context.Context, portfolio *Portfolio, strategy *RebalanceStrategy) (bool, []string) {
	var triggeredConditions []string

	currentAllocations := r.calculateCurrentAllocations(portfolio)

	for _, trigger := range strategy.TriggerConditions {
		triggered := false

		switch trigger.Type {
		case TriggerTypeDrift:
			maxDrift := r.calculateMaxDrift(currentAllocations, strategy.TargetAllocations)
			if maxDrift.GreaterThan(trigger.Threshold) {
				triggered = true
				triggeredConditions = append(triggeredConditions, fmt.Sprintf("drift: %s", maxDrift.String()))
			}

		case TriggerTypeTime:
			hoursSinceRebalance := time.Since(strategy.LastRebalance).Hours()
			if decimal.NewFromFloat(hoursSinceRebalance).GreaterThan(trigger.Threshold) {
				triggered = true
				triggeredConditions = append(triggeredConditions, fmt.Sprintf("time: %f hours", hoursSinceRebalance))
			}

		case TriggerTypeVolatility:
			// This would calculate portfolio volatility
			// For now, use a placeholder
			volatility := decimal.NewFromFloat(0.1) // 10% volatility
			if volatility.GreaterThan(trigger.Threshold) {
				triggered = true
				triggeredConditions = append(triggeredConditions, fmt.Sprintf("volatility: %s", volatility.String()))
			}

		case TriggerTypeDrawdown:
			// Calculate portfolio drawdown
			drawdown := r.calculateDrawdown(portfolio)
			if drawdown.GreaterThan(trigger.Threshold) {
				triggered = true
				triggeredConditions = append(triggeredConditions, fmt.Sprintf("drawdown: %s", drawdown.String()))
			}
		}

		if triggered {
			return true, triggeredConditions
		}
	}

	return false, triggeredConditions
}

// calculateCurrentAllocations calculates current portfolio allocations
func (r *PortfolioRebalancer) calculateCurrentAllocations(portfolio *Portfolio) map[string]decimal.Decimal {
	allocations := make(map[string]decimal.Decimal)

	for asset, holding := range portfolio.Holdings {
		allocation := holding.Value.Div(portfolio.TotalValue)
		allocations[asset] = allocation
	}

	return allocations
}

// calculateMaxDrift calculates maximum drift from target allocations
func (r *PortfolioRebalancer) calculateMaxDrift(current, target map[string]decimal.Decimal) decimal.Decimal {
	maxDrift := decimal.Zero

	for asset, targetAllocation := range target {
		currentAllocation := current[asset]
		if currentAllocation.IsZero() {
			currentAllocation = decimal.Zero
		}

		drift := currentAllocation.Sub(targetAllocation).Abs()
		if drift.GreaterThan(maxDrift) {
			maxDrift = drift
		}
	}

	return maxDrift
}

// calculateDrawdown calculates portfolio drawdown
func (r *PortfolioRebalancer) calculateDrawdown(portfolio *Portfolio) decimal.Decimal {
	// This would calculate the maximum drawdown from peak
	// For now, use a simplified calculation based on daily P&L
	if portfolio.DailyPnL.IsNegative() {
		return portfolio.DailyPnL.Abs().Div(portfolio.TotalValue)
	}
	return decimal.Zero
}

// generateRebalanceActions generates actions needed to rebalance portfolio
func (r *PortfolioRebalancer) generateRebalanceActions(
	ctx context.Context,
	portfolio *Portfolio,
	strategy *RebalanceStrategy,
	currentAllocations map[string]decimal.Decimal,
) []*RebalanceAction {
	var actions []*RebalanceAction

	for asset, targetAllocation := range strategy.TargetAllocations {
		currentAllocation := currentAllocations[asset]
		if currentAllocation.IsZero() {
			currentAllocation = decimal.Zero
		}

		targetValue := targetAllocation.Mul(portfolio.TotalValue)
		currentValue := currentAllocation.Mul(portfolio.TotalValue)
		difference := targetValue.Sub(currentValue)

		// Skip small differences
		if difference.Abs().LessThan(r.config.MinRebalanceAmount) {
			continue
		}

		var actionType ActionType
		var amount decimal.Decimal

		if difference.IsPositive() {
			// Need to buy more of this asset
			actionType = ActionTypeBuy
			amount = difference
		} else {
			// Need to sell some of this asset
			actionType = ActionTypeSell
			amount = difference.Abs()
		}

		action := &RebalanceAction{
			ID:          uuid.New(),
			PortfolioID: portfolio.ID,
			ActionType:  actionType,
			ToAsset:     asset,
			Amount:      amount,
			Priority:    1,
			Reason:      fmt.Sprintf("Rebalance to target allocation: %s", targetAllocation.String()),
			CreatedAt:   time.Now(),
		}

		actions = append(actions, action)
	}

	return actions
}

// executeRebalanceAction executes a single rebalance action
func (r *PortfolioRebalancer) executeRebalanceAction(ctx context.Context, action *RebalanceAction) error {
	// This would execute the actual trade
	// For now, just log the action
	r.logger.Info(ctx, "Executing rebalance action", map[string]interface{}{
		"action_id":   action.ID.String(),
		"action_type": string(action.ActionType),
		"asset":       action.ToAsset,
		"amount":      action.Amount.String(),
		"reason":      action.Reason,
	})

	return nil
}
