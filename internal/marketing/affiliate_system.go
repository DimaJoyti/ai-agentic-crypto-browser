package marketing

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// AffiliateSystem manages affiliate marketing and referrals
type AffiliateSystem struct {
	db     *sql.DB
	config *AffiliateConfig
}

// AffiliateConfig defines affiliate program settings
type AffiliateConfig struct {
	DefaultCommissionRate decimal.Decimal `json:"default_commission_rate"` // e.g., 0.20 for 20%
	CookieDuration        time.Duration   `json:"cookie_duration"`         // e.g., 30 days
	MinPayoutAmount       decimal.Decimal `json:"min_payout_amount"`       // e.g., $100
	PayoutSchedule        string          `json:"payout_schedule"`         // monthly, weekly
	TierCommissions       map[string]decimal.Decimal `json:"tier_commissions"`
}

// Affiliate represents an affiliate partner
type Affiliate struct {
	ID                string          `json:"id"`
	UserID            string          `json:"user_id"`
	Code              string          `json:"code"`              // Unique affiliate code
	Name              string          `json:"name"`
	Email             string          `json:"email"`
	Status            string          `json:"status"`            // active, suspended, pending
	CommissionRate    decimal.Decimal `json:"commission_rate"`
	TotalEarnings     decimal.Decimal `json:"total_earnings"`
	PendingEarnings   decimal.Decimal `json:"pending_earnings"`
	PaidEarnings      decimal.Decimal `json:"paid_earnings"`
	TotalReferrals    int             `json:"total_referrals"`
	ActiveReferrals   int             `json:"active_referrals"`
	ConversionRate    float64         `json:"conversion_rate"`
	PaymentMethod     string          `json:"payment_method"`
	PaymentDetails    string          `json:"payment_details"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

// Referral represents a referral tracking record
type Referral struct {
	ID               string          `json:"id"`
	AffiliateID      string          `json:"affiliate_id"`
	AffiliateCode    string          `json:"affiliate_code"`
	ReferredUserID   string          `json:"referred_user_id"`
	ReferredEmail    string          `json:"referred_email"`
	Status           string          `json:"status"`           // pending, converted, cancelled
	SubscriptionID   string          `json:"subscription_id"`
	CommissionAmount decimal.Decimal `json:"commission_amount"`
	CommissionRate   decimal.Decimal `json:"commission_rate"`
	ConvertedAt      *time.Time      `json:"converted_at,omitempty"`
	FirstPurchase    decimal.Decimal `json:"first_purchase"`
	LifetimeValue    decimal.Decimal `json:"lifetime_value"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

// Commission represents a commission payment
type Commission struct {
	ID          string          `json:"id"`
	AffiliateID string          `json:"affiliate_id"`
	ReferralID  string          `json:"referral_id"`
	Amount      decimal.Decimal `json:"amount"`
	Type        string          `json:"type"`        // signup, subscription, performance
	Status      string          `json:"status"`      // pending, paid, cancelled
	PaidAt      *time.Time      `json:"paid_at,omitempty"`
	PayoutID    string          `json:"payout_id"`
	CreatedAt   time.Time       `json:"created_at"`
}

// NewAffiliateSystem creates a new affiliate system
func NewAffiliateSystem(db *sql.DB) *AffiliateSystem {
	config := &AffiliateConfig{
		DefaultCommissionRate: decimal.NewFromFloat(0.20), // 20%
		CookieDuration:        30 * 24 * time.Hour,        // 30 days
		MinPayoutAmount:       decimal.NewFromInt(100),    // $100
		PayoutSchedule:        "monthly",
		TierCommissions: map[string]decimal.Decimal{
			"starter":      decimal.NewFromFloat(0.15), // 15%
			"professional": decimal.NewFromFloat(0.20), // 20%
			"enterprise":   decimal.NewFromFloat(0.25), // 25%
		},
	}

	return &AffiliateSystem{
		db:     db,
		config: config,
	}
}

// CreateAffiliate creates a new affiliate account
func (as *AffiliateSystem) CreateAffiliate(ctx context.Context, userID, name, email string) (*Affiliate, error) {
	// Generate unique affiliate code
	code, err := as.generateAffiliateCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate affiliate code: %w", err)
	}

	affiliate := &Affiliate{
		ID:             generateID(),
		UserID:         userID,
		Code:           code,
		Name:           name,
		Email:          email,
		Status:         "pending",
		CommissionRate: as.config.DefaultCommissionRate,
		TotalEarnings:  decimal.Zero,
		PendingEarnings: decimal.Zero,
		PaidEarnings:   decimal.Zero,
		TotalReferrals: 0,
		ActiveReferrals: 0,
		ConversionRate: 0.0,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Save to database
	err = as.saveAffiliate(ctx, affiliate)
	if err != nil {
		return nil, fmt.Errorf("failed to save affiliate: %w", err)
	}

	return affiliate, nil
}

// TrackReferral tracks a new referral click
func (as *AffiliateSystem) TrackReferral(ctx context.Context, affiliateCode, referredEmail, ipAddress string) (*Referral, error) {
	// Get affiliate by code
	affiliate, err := as.getAffiliateByCode(ctx, affiliateCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get affiliate: %w", err)
	}

	if affiliate.Status != "active" {
		return nil, fmt.Errorf("affiliate is not active")
	}

	// Check for existing referral
	existing, err := as.getReferralByEmail(ctx, referredEmail)
	if err == nil && existing != nil {
		return existing, nil // Return existing referral
	}

	referral := &Referral{
		ID:            generateID(),
		AffiliateID:   affiliate.ID,
		AffiliateCode: affiliateCode,
		ReferredEmail: referredEmail,
		Status:        "pending",
		CommissionRate: affiliate.CommissionRate,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Save referral
	err = as.saveReferral(ctx, referral)
	if err != nil {
		return nil, fmt.Errorf("failed to save referral: %w", err)
	}

	return referral, nil
}

// ConvertReferral converts a referral when user subscribes
func (as *AffiliateSystem) ConvertReferral(ctx context.Context, userID, subscriptionID string, subscriptionAmount decimal.Decimal, tierID string) error {
	// Get referral by user ID
	referral, err := as.getReferralByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get referral: %w", err)
	}

	if referral == nil || referral.Status != "pending" {
		return nil // No pending referral to convert
	}

	// Get commission rate for tier
	commissionRate := as.config.DefaultCommissionRate
	if tierRate, exists := as.config.TierCommissions[tierID]; exists {
		commissionRate = tierRate
	}

	// Calculate commission
	commissionAmount := subscriptionAmount.Mul(commissionRate)

	// Update referral
	now := time.Now()
	referral.Status = "converted"
	referral.ReferredUserID = userID
	referral.SubscriptionID = subscriptionID
	referral.CommissionAmount = commissionAmount
	referral.CommissionRate = commissionRate
	referral.ConvertedAt = &now
	referral.FirstPurchase = subscriptionAmount
	referral.UpdatedAt = now

	err = as.updateReferral(ctx, referral)
	if err != nil {
		return fmt.Errorf("failed to update referral: %w", err)
	}

	// Create commission record
	commission := &Commission{
		ID:          generateID(),
		AffiliateID: referral.AffiliateID,
		ReferralID:  referral.ID,
		Amount:      commissionAmount,
		Type:        "subscription",
		Status:      "pending",
		CreatedAt:   now,
	}

	err = as.saveCommission(ctx, commission)
	if err != nil {
		return fmt.Errorf("failed to save commission: %w", err)
	}

	// Update affiliate stats
	err = as.updateAffiliateStats(ctx, referral.AffiliateID)
	if err != nil {
		return fmt.Errorf("failed to update affiliate stats: %w", err)
	}

	return nil
}

// GetAffiliateStats returns affiliate performance statistics
func (as *AffiliateSystem) GetAffiliateStats(ctx context.Context, affiliateID string) (*AffiliateStats, error) {
	stats := &AffiliateStats{}

	// Get basic affiliate info
	affiliate, err := as.getAffiliate(ctx, affiliateID)
	if err != nil {
		return nil, err
	}

	stats.Affiliate = affiliate

	// Get referral stats
	query := `
		SELECT 
			COUNT(*) as total_referrals,
			COUNT(CASE WHEN status = 'converted' THEN 1 END) as converted_referrals,
			COALESCE(SUM(CASE WHEN status = 'converted' THEN commission_amount ELSE 0 END), 0) as total_commissions,
			COALESCE(AVG(CASE WHEN status = 'converted' THEN first_purchase ELSE NULL END), 0) as avg_order_value
		FROM referrals 
		WHERE affiliate_id = $1
	`

	row := as.db.QueryRowContext(ctx, query, affiliateID)
	err = row.Scan(
		&stats.TotalReferrals,
		&stats.ConvertedReferrals,
		&stats.TotalCommissions,
		&stats.AverageOrderValue,
	)
	if err != nil {
		return nil, err
	}

	// Calculate conversion rate
	if stats.TotalReferrals > 0 {
		stats.ConversionRate = float64(stats.ConvertedReferrals) / float64(stats.TotalReferrals)
	}

	// Get monthly stats
	monthlyStats, err := as.getMonthlyStats(ctx, affiliateID)
	if err != nil {
		return nil, err
	}
	stats.MonthlyStats = monthlyStats

	return stats, nil
}

// AffiliateStats represents affiliate performance statistics
type AffiliateStats struct {
	Affiliate          *Affiliate                 `json:"affiliate"`
	TotalReferrals     int                        `json:"total_referrals"`
	ConvertedReferrals int                        `json:"converted_referrals"`
	ConversionRate     float64                    `json:"conversion_rate"`
	TotalCommissions   decimal.Decimal            `json:"total_commissions"`
	AverageOrderValue  decimal.Decimal            `json:"average_order_value"`
	MonthlyStats       map[string]*MonthlyStats   `json:"monthly_stats"`
}

// MonthlyStats represents monthly performance
type MonthlyStats struct {
	Month       string          `json:"month"`
	Referrals   int             `json:"referrals"`
	Conversions int             `json:"conversions"`
	Commissions decimal.Decimal `json:"commissions"`
}

// Helper functions
func (as *AffiliateSystem) generateAffiliateCode() (string, error) {
	bytes := make([]byte, 4)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (as *AffiliateSystem) saveAffiliate(ctx context.Context, affiliate *Affiliate) error {
	query := `
		INSERT INTO affiliates (
			id, user_id, code, name, email, status, commission_rate,
			total_earnings, pending_earnings, paid_earnings,
			total_referrals, active_referrals, conversion_rate,
			payment_method, payment_details, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`

	_, err := as.db.ExecContext(ctx, query,
		affiliate.ID, affiliate.UserID, affiliate.Code, affiliate.Name,
		affiliate.Email, affiliate.Status, affiliate.CommissionRate,
		affiliate.TotalEarnings, affiliate.PendingEarnings, affiliate.PaidEarnings,
		affiliate.TotalReferrals, affiliate.ActiveReferrals, affiliate.ConversionRate,
		affiliate.PaymentMethod, affiliate.PaymentDetails,
		affiliate.CreatedAt, affiliate.UpdatedAt,
	)

	return err
}

func (as *AffiliateSystem) getAffiliateByCode(ctx context.Context, code string) (*Affiliate, error) {
	// Implementation to get affiliate by code
	return nil, nil
}

func (as *AffiliateSystem) getReferralByEmail(ctx context.Context, email string) (*Referral, error) {
	// Implementation to get referral by email
	return nil, nil
}

func (as *AffiliateSystem) getReferralByUserID(ctx context.Context, userID string) (*Referral, error) {
	// Implementation to get referral by user ID
	return nil, nil
}

func (as *AffiliateSystem) saveReferral(ctx context.Context, referral *Referral) error {
	// Implementation to save referral
	return nil
}

func (as *AffiliateSystem) updateReferral(ctx context.Context, referral *Referral) error {
	// Implementation to update referral
	return nil
}

func (as *AffiliateSystem) saveCommission(ctx context.Context, commission *Commission) error {
	// Implementation to save commission
	return nil
}

func (as *AffiliateSystem) updateAffiliateStats(ctx context.Context, affiliateID string) error {
	// Implementation to update affiliate statistics
	return nil
}

func (as *AffiliateSystem) getAffiliate(ctx context.Context, affiliateID string) (*Affiliate, error) {
	// Implementation to get affiliate
	return nil, nil
}

func (as *AffiliateSystem) getMonthlyStats(ctx context.Context, affiliateID string) (map[string]*MonthlyStats, error) {
	// Implementation to get monthly statistics
	return nil, nil
}

func generateID() string {
	return fmt.Sprintf("aff_%d", time.Now().UnixNano())
}
