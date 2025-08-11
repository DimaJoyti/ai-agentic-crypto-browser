package billing

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// SubscriptionManager handles subscription billing and management
type SubscriptionManager struct {
	db     *sql.DB
	config *SubscriptionConfig
}

// SubscriptionConfig defines subscription tiers and pricing
type SubscriptionConfig struct {
	Tiers map[string]*SubscriptionTier `json:"tiers"`
}

// SubscriptionTier defines a subscription tier
type SubscriptionTier struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Price        decimal.Decimal `json:"price"`        // Monthly price
	AnnualPrice  decimal.Decimal `json:"annual_price"` // Annual price (with discount)
	Features     []string        `json:"features"`
	Limits       *TierLimits     `json:"limits"`
	Description  string          `json:"description"`
	Popular      bool            `json:"popular"`
	Enterprise   bool            `json:"enterprise"`
}

// TierLimits defines usage limits for each tier
type TierLimits struct {
	MaxStrategies    int   `json:"max_strategies"`
	MaxPortfolios    int   `json:"max_portfolios"`
	MaxAPIRequests   int   `json:"max_api_requests"`   // Per month
	MaxChains        int   `json:"max_chains"`
	AdvancedAI       bool  `json:"advanced_ai"`
	VoiceCommands    bool  `json:"voice_commands"`
	CustomStrategies bool  `json:"custom_strategies"`
	PrioritySupport  bool  `json:"priority_support"`
	WhiteLabel       bool  `json:"white_label"`
}

// Subscription represents a user's subscription
type Subscription struct {
	ID                string          `json:"id"`
	UserID            string          `json:"user_id"`
	TierID            string          `json:"tier_id"`
	Status            string          `json:"status"` // active, cancelled, expired, trial
	BillingCycle      string          `json:"billing_cycle"` // monthly, annual
	Price             decimal.Decimal `json:"price"`
	Currency          string          `json:"currency"`
	StartDate         time.Time       `json:"start_date"`
	EndDate           time.Time       `json:"end_date"`
	TrialEndDate      *time.Time      `json:"trial_end_date,omitempty"`
	CancelledAt       *time.Time      `json:"cancelled_at,omitempty"`
	CancellationReason string         `json:"cancellation_reason,omitempty"`
	PaymentMethodID   string          `json:"payment_method_id"`
	NextBillingDate   time.Time       `json:"next_billing_date"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

// NewSubscriptionManager creates a new subscription manager
func NewSubscriptionManager(db *sql.DB) *SubscriptionManager {
	config := &SubscriptionConfig{
		Tiers: map[string]*SubscriptionTier{
			"starter": {
				ID:          "starter",
				Name:        "Starter",
				Price:       decimal.NewFromInt(49),
				AnnualPrice: decimal.NewFromInt(490), // 2 months free
				Features: []string{
					"Basic AI Trading",
					"3 Trading Strategies",
					"Single Chain Support",
					"Basic Analytics",
					"Email Support",
				},
				Limits: &TierLimits{
					MaxStrategies:    3,
					MaxPortfolios:    1,
					MaxAPIRequests:   1000,
					MaxChains:        1,
					AdvancedAI:       false,
					VoiceCommands:    false,
					CustomStrategies: false,
					PrioritySupport:  false,
					WhiteLabel:       false,
				},
				Description: "Perfect for beginners getting started with AI trading",
				Popular:     false,
				Enterprise:  false,
			},
			"professional": {
				ID:          "professional",
				Name:        "Professional",
				Price:       decimal.NewFromInt(199),
				AnnualPrice: decimal.NewFromInt(1990), // 2 months free
				Features: []string{
					"Advanced AI Trading",
					"10+ Trading Strategies",
					"Multi-Chain Support",
					"Advanced Analytics",
					"DeFi Integration",
					"Voice Commands",
					"Priority Support",
				},
				Limits: &TierLimits{
					MaxStrategies:    15,
					MaxPortfolios:    5,
					MaxAPIRequests:   10000,
					MaxChains:        5,
					AdvancedAI:       true,
					VoiceCommands:    true,
					CustomStrategies: true,
					PrioritySupport:  true,
					WhiteLabel:       false,
				},
				Description: "For serious traders who want advanced features",
				Popular:     true,
				Enterprise:  false,
			},
			"enterprise": {
				ID:          "enterprise",
				Name:        "Enterprise",
				Price:       decimal.NewFromInt(999),
				AnnualPrice: decimal.NewFromInt(9990), // 2 months free
				Features: []string{
					"Full Platform Access",
					"Unlimited Strategies",
					"All Chains Supported",
					"Custom AI Models",
					"White-Label Solution",
					"Dedicated Support",
					"Custom Integrations",
					"SLA Guarantee",
				},
				Limits: &TierLimits{
					MaxStrategies:    -1, // Unlimited
					MaxPortfolios:    -1, // Unlimited
					MaxAPIRequests:   100000,
					MaxChains:        -1, // All chains
					AdvancedAI:       true,
					VoiceCommands:    true,
					CustomStrategies: true,
					PrioritySupport:  true,
					WhiteLabel:       true,
				},
				Description: "Complete solution for enterprises and institutions",
				Popular:     false,
				Enterprise:  true,
			},
		},
	}

	return &SubscriptionManager{
		db:     db,
		config: config,
	}
}

// CreateSubscription creates a new subscription for a user
func (sm *SubscriptionManager) CreateSubscription(ctx context.Context, userID, tierID, billingCycle string, trialDays int) (*Subscription, error) {
	tier, exists := sm.config.Tiers[tierID]
	if !exists {
		return nil, fmt.Errorf("invalid tier ID: %s", tierID)
	}

	now := time.Now()
	var price decimal.Decimal
	var endDate time.Time

	// Set price based on billing cycle
	if billingCycle == "annual" {
		price = tier.AnnualPrice
		endDate = now.AddDate(1, 0, 0) // 1 year
	} else {
		price = tier.Price
		endDate = now.AddDate(0, 1, 0) // 1 month
	}

	subscription := &Subscription{
		ID:              generateSubscriptionID(),
		UserID:          userID,
		TierID:          tierID,
		Status:          "active",
		BillingCycle:    billingCycle,
		Price:           price,
		Currency:        "USD",
		StartDate:       now,
		EndDate:         endDate,
		NextBillingDate: endDate,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	// Add trial period if specified
	if trialDays > 0 {
		trialEnd := now.AddDate(0, 0, trialDays)
		subscription.TrialEndDate = &trialEnd
		subscription.Status = "trial"
		subscription.NextBillingDate = trialEnd
	}

	// Save to database
	err := sm.saveSubscription(ctx, subscription)
	if err != nil {
		return nil, fmt.Errorf("failed to save subscription: %w", err)
	}

	return subscription, nil
}

// GetUserSubscription returns the active subscription for a user
func (sm *SubscriptionManager) GetUserSubscription(ctx context.Context, userID string) (*Subscription, error) {
	query := `
		SELECT id, user_id, tier_id, status, billing_cycle, price, currency,
		       start_date, end_date, trial_end_date, cancelled_at, 
		       cancellation_reason, payment_method_id, next_billing_date,
		       created_at, updated_at
		FROM subscriptions 
		WHERE user_id = $1 AND status IN ('active', 'trial')
		ORDER BY created_at DESC
		LIMIT 1
	`

	row := sm.db.QueryRowContext(ctx, query, userID)
	
	subscription := &Subscription{}
	err := row.Scan(
		&subscription.ID, &subscription.UserID, &subscription.TierID,
		&subscription.Status, &subscription.BillingCycle, &subscription.Price,
		&subscription.Currency, &subscription.StartDate, &subscription.EndDate,
		&subscription.TrialEndDate, &subscription.CancelledAt,
		&subscription.CancellationReason, &subscription.PaymentMethodID,
		&subscription.NextBillingDate, &subscription.CreatedAt, &subscription.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil // No active subscription
	}
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

// CheckUserAccess verifies if user has access to a feature
func (sm *SubscriptionManager) CheckUserAccess(ctx context.Context, userID string, feature string) (bool, error) {
	subscription, err := sm.GetUserSubscription(ctx, userID)
	if err != nil {
		return false, err
	}

	if subscription == nil {
		return false, nil // No subscription = no access
	}

	tier, exists := sm.config.Tiers[subscription.TierID]
	if !exists {
		return false, fmt.Errorf("invalid tier ID: %s", subscription.TierID)
	}

	// Check feature access based on tier limits
	switch feature {
	case "advanced_ai":
		return tier.Limits.AdvancedAI, nil
	case "voice_commands":
		return tier.Limits.VoiceCommands, nil
	case "custom_strategies":
		return tier.Limits.CustomStrategies, nil
	case "priority_support":
		return tier.Limits.PrioritySupport, nil
	case "white_label":
		return tier.Limits.WhiteLabel, nil
	default:
		return true, nil // Basic features available to all
	}
}

// GetTierLimits returns the limits for a user's subscription tier
func (sm *SubscriptionManager) GetTierLimits(ctx context.Context, userID string) (*TierLimits, error) {
	subscription, err := sm.GetUserSubscription(ctx, userID)
	if err != nil {
		return nil, err
	}

	if subscription == nil {
		// Return free tier limits
		return &TierLimits{
			MaxStrategies:    1,
			MaxPortfolios:    1,
			MaxAPIRequests:   100,
			MaxChains:        1,
			AdvancedAI:       false,
			VoiceCommands:    false,
			CustomStrategies: false,
			PrioritySupport:  false,
			WhiteLabel:       false,
		}, nil
	}

	tier, exists := sm.config.Tiers[subscription.TierID]
	if !exists {
		return nil, fmt.Errorf("invalid tier ID: %s", subscription.TierID)
	}

	return tier.Limits, nil
}

// Helper functions
func (sm *SubscriptionManager) saveSubscription(ctx context.Context, subscription *Subscription) error {
	query := `
		INSERT INTO subscriptions (
			id, user_id, tier_id, status, billing_cycle, price, currency,
			start_date, end_date, trial_end_date, payment_method_id,
			next_billing_date, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := sm.db.ExecContext(ctx, query,
		subscription.ID, subscription.UserID, subscription.TierID,
		subscription.Status, subscription.BillingCycle, subscription.Price,
		subscription.Currency, subscription.StartDate, subscription.EndDate,
		subscription.TrialEndDate, subscription.PaymentMethodID,
		subscription.NextBillingDate, subscription.CreatedAt, subscription.UpdatedAt,
	)

	return err
}

func generateSubscriptionID() string {
	return fmt.Sprintf("sub_%d", time.Now().UnixNano())
}
