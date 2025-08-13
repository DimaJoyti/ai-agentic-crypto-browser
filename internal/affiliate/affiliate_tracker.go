package affiliate

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// AffiliateTracker manages affiliate and referral programs
type AffiliateTracker struct {
	db *sql.DB
}

// NewAffiliateTracker creates a new affiliate tracker
func NewAffiliateTracker(db *sql.DB) *AffiliateTracker {
	return &AffiliateTracker{
		db: db,
	}
}

// Affiliate represents an affiliate partner
type Affiliate struct {
	ID                string          `json:"id"`
	UserID            string          `json:"user_id"`
	AffiliateCode     string          `json:"affiliate_code"`
	AffiliateType     string          `json:"affiliate_type"`     // individual, business, influencer
	CommissionRate    decimal.Decimal `json:"commission_rate"`    // 0.20 = 20%
	PaymentMethod     string          `json:"payment_method"`     // stripe, crypto, bank
	PaymentDetails    string          `json:"payment_details"`    // encrypted payment info
	Status            string          `json:"status"`             // active, suspended, pending
	TotalReferrals    int64           `json:"total_referrals"`
	TotalCommissions  decimal.Decimal `json:"total_commissions"`
	UnpaidCommissions decimal.Decimal `json:"unpaid_commissions"`
	LastPayoutAt      *time.Time      `json:"last_payout_at"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	Metadata          map[string]interface{} `json:"metadata"`
}

// Referral represents a referral record
type Referral struct {
	ID               string          `json:"id"`
	AffiliateID      string          `json:"affiliate_id"`
	ReferredUserID   string          `json:"referred_user_id"`
	ReferralCode     string          `json:"referral_code"`
	ReferralSource   string          `json:"referral_source"`   // link, social, email, etc.
	ConversionType   string          `json:"conversion_type"`   // signup, subscription, purchase
	ConversionValue  decimal.Decimal `json:"conversion_value"`  // value of the conversion
	CommissionAmount decimal.Decimal `json:"commission_amount"`
	CommissionRate   decimal.Decimal `json:"commission_rate"`
	Status           string          `json:"status"`            // pending, confirmed, paid, cancelled
	ConvertedAt      time.Time       `json:"converted_at"`
	ConfirmedAt      *time.Time      `json:"confirmed_at"`
	PaidAt           *time.Time      `json:"paid_at"`
	CreatedAt        time.Time       `json:"created_at"`
}

// CommissionTier defines commission rates based on performance
type CommissionTier struct {
	ID                string          `json:"id"`
	TierName          string          `json:"tier_name"`
	MinReferrals      int64           `json:"min_referrals"`
	MaxReferrals      int64           `json:"max_referrals"`
	CommissionRate    decimal.Decimal `json:"commission_rate"`
	BonusRate         decimal.Decimal `json:"bonus_rate"`
	RequiredRevenue   decimal.Decimal `json:"required_revenue"`
	Description       string          `json:"description"`
	IsActive          bool            `json:"is_active"`
}

// AffiliateStats represents affiliate performance statistics
type AffiliateStats struct {
	AffiliateID       string          `json:"affiliate_id"`
	Period            string          `json:"period"`
	TotalClicks       int64           `json:"total_clicks"`
	TotalSignups      int64           `json:"total_signups"`
	TotalConversions  int64           `json:"total_conversions"`
	ConversionRate    decimal.Decimal `json:"conversion_rate"`
	TotalRevenue      decimal.Decimal `json:"total_revenue"`
	TotalCommissions  decimal.Decimal `json:"total_commissions"`
	AverageOrderValue decimal.Decimal `json:"average_order_value"`
	TopReferralSource string          `json:"top_referral_source"`
	LastUpdated       time.Time       `json:"last_updated"`
}

// CreateAffiliate creates a new affiliate account
func (at *AffiliateTracker) CreateAffiliate(ctx context.Context, affiliate *Affiliate) error {
	query := `
		INSERT INTO affiliates (
			id, user_id, affiliate_code, affiliate_type, commission_rate,
			payment_method, payment_details, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := at.db.ExecContext(ctx, query,
		affiliate.ID, affiliate.UserID, affiliate.AffiliateCode, affiliate.AffiliateType,
		affiliate.CommissionRate, affiliate.PaymentMethod, affiliate.PaymentDetails,
		affiliate.Status, affiliate.CreatedAt, affiliate.UpdatedAt,
	)

	return err
}

// TrackReferral records a new referral
func (at *AffiliateTracker) TrackReferral(ctx context.Context, referral *Referral) error {
	// Calculate commission based on affiliate's rate
	affiliate, err := at.GetAffiliateByCode(ctx, referral.ReferralCode)
	if err != nil {
		return fmt.Errorf("failed to get affiliate: %v", err)
	}

	referral.AffiliateID = affiliate.ID
	referral.CommissionRate = affiliate.CommissionRate
	referral.CommissionAmount = referral.ConversionValue.Mul(affiliate.CommissionRate)
	referral.Status = "pending"
	referral.CreatedAt = time.Now()

	query := `
		INSERT INTO referrals (
			id, affiliate_id, referred_user_id, referral_code, referral_source,
			conversion_type, conversion_value, commission_amount, commission_rate,
			status, converted_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err = at.db.ExecContext(ctx, query,
		referral.ID, referral.AffiliateID, referral.ReferredUserID, referral.ReferralCode,
		referral.ReferralSource, referral.ConversionType, referral.ConversionValue,
		referral.CommissionAmount, referral.CommissionRate, referral.Status,
		referral.ConvertedAt, referral.CreatedAt,
	)

	if err != nil {
		return err
	}

	// Update affiliate stats
	return at.updateAffiliateStats(ctx, affiliate.ID, referral)
}

// GetAffiliateByCode retrieves affiliate by referral code
func (at *AffiliateTracker) GetAffiliateByCode(ctx context.Context, code string) (*Affiliate, error) {
	query := `
		SELECT id, user_id, affiliate_code, affiliate_type, commission_rate,
		       payment_method, status, total_referrals, total_commissions,
		       unpaid_commissions, created_at, updated_at
		FROM affiliates 
		WHERE affiliate_code = $1 AND status = 'active'
	`

	affiliate := &Affiliate{}
	err := at.db.QueryRowContext(ctx, query, code).Scan(
		&affiliate.ID, &affiliate.UserID, &affiliate.AffiliateCode, &affiliate.AffiliateType,
		&affiliate.CommissionRate, &affiliate.PaymentMethod, &affiliate.Status,
		&affiliate.TotalReferrals, &affiliate.TotalCommissions, &affiliate.UnpaidCommissions,
		&affiliate.CreatedAt, &affiliate.UpdatedAt,
	)

	return affiliate, err
}

// ConfirmReferral confirms a pending referral
func (at *AffiliateTracker) ConfirmReferral(ctx context.Context, referralID string) error {
	now := time.Now()
	
	query := `
		UPDATE referrals 
		SET status = 'confirmed', confirmed_at = $1
		WHERE id = $2 AND status = 'pending'
	`

	result, err := at.db.ExecContext(ctx, query, now, referralID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("referral not found or already confirmed")
	}

	// Update affiliate unpaid commissions
	return at.updateUnpaidCommissions(ctx, referralID)
}

// GetAffiliateStats returns performance statistics for an affiliate
func (at *AffiliateTracker) GetAffiliateStats(ctx context.Context, affiliateID, period string) (*AffiliateStats, error) {
	query := `
		SELECT affiliate_id, period, total_clicks, total_signups, total_conversions,
		       conversion_rate, total_revenue, total_commissions, average_order_value,
		       top_referral_source, last_updated
		FROM affiliate_stats 
		WHERE affiliate_id = $1 AND period = $2
	`

	stats := &AffiliateStats{}
	err := at.db.QueryRowContext(ctx, query, affiliateID, period).Scan(
		&stats.AffiliateID, &stats.Period, &stats.TotalClicks, &stats.TotalSignups,
		&stats.TotalConversions, &stats.ConversionRate, &stats.TotalRevenue,
		&stats.TotalCommissions, &stats.AverageOrderValue, &stats.TopReferralSource,
		&stats.LastUpdated,
	)

	return stats, err
}

// GetTopAffiliates returns top performing affiliates
func (at *AffiliateTracker) GetTopAffiliates(ctx context.Context, limit int) ([]*Affiliate, error) {
	query := `
		SELECT id, user_id, affiliate_code, affiliate_type, commission_rate,
		       status, total_referrals, total_commissions, unpaid_commissions,
		       created_at, updated_at
		FROM affiliates 
		WHERE status = 'active'
		ORDER BY total_commissions DESC 
		LIMIT $1
	`

	rows, err := at.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var affiliates []*Affiliate
	for rows.Next() {
		affiliate := &Affiliate{}
		err := rows.Scan(
			&affiliate.ID, &affiliate.UserID, &affiliate.AffiliateCode, &affiliate.AffiliateType,
			&affiliate.CommissionRate, &affiliate.Status, &affiliate.TotalReferrals,
			&affiliate.TotalCommissions, &affiliate.UnpaidCommissions,
			&affiliate.CreatedAt, &affiliate.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		affiliates = append(affiliates, affiliate)
	}

	return affiliates, nil
}

// updateAffiliateStats updates affiliate performance statistics
func (at *AffiliateTracker) updateAffiliateStats(ctx context.Context, affiliateID string, referral *Referral) error {
	// Update affiliate totals
	query := `
		UPDATE affiliates 
		SET total_referrals = total_referrals + 1,
		    total_commissions = total_commissions + $1,
		    unpaid_commissions = unpaid_commissions + $1,
		    updated_at = $2
		WHERE id = $3
	`

	_, err := at.db.ExecContext(ctx, query, referral.CommissionAmount, time.Now(), affiliateID)
	return err
}

// updateUnpaidCommissions updates unpaid commission amount
func (at *AffiliateTracker) updateUnpaidCommissions(ctx context.Context, referralID string) error {
	query := `
		UPDATE affiliates 
		SET unpaid_commissions = unpaid_commissions + (
			SELECT commission_amount FROM referrals WHERE id = $1
		)
		WHERE id = (
			SELECT affiliate_id FROM referrals WHERE id = $1
		)
	`

	_, err := at.db.ExecContext(ctx, query, referralID)
	return err
}

// ProcessPayouts processes affiliate commission payouts
func (at *AffiliateTracker) ProcessPayouts(ctx context.Context, minimumPayout decimal.Decimal) error {
	// Get affiliates with unpaid commissions above minimum
	query := `
		SELECT id, user_id, affiliate_code, unpaid_commissions, payment_method
		FROM affiliates 
		WHERE unpaid_commissions >= $1 AND status = 'active'
	`

	rows, err := at.db.QueryContext(ctx, query, minimumPayout)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var affiliateID, userID, affiliateCode, paymentMethod string
		var unpaidCommissions decimal.Decimal

		err := rows.Scan(&affiliateID, &userID, &affiliateCode, &unpaidCommissions, &paymentMethod)
		if err != nil {
			continue
		}

		// Process payout (implementation would integrate with payment provider)
		err = at.processSinglePayout(ctx, affiliateID, unpaidCommissions, paymentMethod)
		if err != nil {
			// Log error but continue with other payouts
			fmt.Printf("Failed to process payout for affiliate %s: %v\n", affiliateCode, err)
			continue
		}
	}

	return nil
}

// processSinglePayout processes a single affiliate payout
func (at *AffiliateTracker) processSinglePayout(ctx context.Context, affiliateID string, amount decimal.Decimal, paymentMethod string) error {
	// Implementation would integrate with Stripe, crypto payments, or bank transfers
	// For now, we'll mark as paid in the database
	
	now := time.Now()
	
	// Update affiliate record
	query := `
		UPDATE affiliates 
		SET unpaid_commissions = 0, last_payout_at = $1, updated_at = $1
		WHERE id = $2
	`

	_, err := at.db.ExecContext(ctx, query, now, affiliateID)
	if err != nil {
		return err
	}

	// Mark referrals as paid
	payQuery := `
		UPDATE referrals 
		SET status = 'paid', paid_at = $1
		WHERE affiliate_id = $2 AND status = 'confirmed'
	`

	_, err = at.db.ExecContext(ctx, payQuery, now, affiliateID)
	return err
}
