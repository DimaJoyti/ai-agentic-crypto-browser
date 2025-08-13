package billing

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/invoice"
	"github.com/stripe/stripe-go/v76/invoiceitem"
)

// APIBillingService handles API usage billing
type APIBillingService struct {
	db              *sql.DB
	usageTracker    *MarketplaceTracker
	stripeProcessor *StripePaymentProcessor
}

// NewAPIBillingService creates a new API billing service
func NewAPIBillingService(db *sql.DB, usageTracker *MarketplaceTracker, stripeProcessor *StripePaymentProcessor) *APIBillingService {
	return &APIBillingService{
		db:              db,
		usageTracker:    usageTracker,
		stripeProcessor: stripeProcessor,
	}
}

// GenerateMonthlyBill generates monthly API usage bill for a user
func (abs *APIBillingService) GenerateMonthlyBill(ctx context.Context, userID string, billingPeriod time.Time) (*APIBillingRecord, error) {
	// Get usage summary for the billing period
	periodStr := billingPeriod.Format("2006-01-01")
	summary, err := abs.usageTracker.GetUsageSummary(ctx, userID, periodStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get usage summary: %v", err)
	}

	// Calculate volume discount
	discountPercent := abs.calculateVolumeDiscount(summary.TotalCost)
	discountAmount := summary.TotalCost.Mul(decimal.NewFromFloat(discountPercent / 100))
	finalAmount := summary.TotalCost.Sub(discountAmount)

	// Create billing record
	billingRecord := &APIBillingRecord{
		UserID:          userID,
		BillingPeriod:   billingPeriod,
		TotalRequests:   summary.TotalRequests,
		TotalCost:       summary.TotalCost,
		DiscountPercent: decimal.NewFromFloat(discountPercent),
		DiscountAmount:  discountAmount,
		FinalAmount:     finalAmount,
		Currency:        "USD",
		Status:          "pending",
		GeneratedAt:     time.Now(),
		DueDate:         time.Now().AddDate(0, 0, 30), // 30 days to pay
	}

	// Store billing record in database
	err = abs.storeBillingRecord(ctx, billingRecord)
	if err != nil {
		return nil, fmt.Errorf("failed to store billing record: %v", err)
	}

	// Create Stripe invoice if amount > $0
	if finalAmount.GreaterThan(decimal.Zero) {
		stripeInvoiceID, err := abs.createStripeInvoice(ctx, userID, billingRecord)
		if err != nil {
			return nil, fmt.Errorf("failed to create Stripe invoice: %v", err)
		}
		billingRecord.StripeInvoiceID = stripeInvoiceID

		// Update billing record with Stripe invoice ID
		err = abs.updateBillingRecord(ctx, billingRecord)
		if err != nil {
			return nil, fmt.Errorf("failed to update billing record: %v", err)
		}
	}

	return billingRecord, nil
}

// ProcessPendingBills processes all pending bills for the month
func (abs *APIBillingService) ProcessPendingBills(ctx context.Context, billingPeriod time.Time) error {
	// Get all users with API usage in the billing period
	users, err := abs.getUsersWithAPIUsage(ctx, billingPeriod)
	if err != nil {
		return fmt.Errorf("failed to get users with API usage: %v", err)
	}

	for _, userID := range users {
		// Check if bill already exists
		exists, err := abs.billingRecordExists(ctx, userID, billingPeriod)
		if err != nil {
			return fmt.Errorf("failed to check billing record existence: %v", err)
		}

		if !exists {
			// Generate bill for user
			_, err := abs.GenerateMonthlyBill(ctx, userID, billingPeriod)
			if err != nil {
				// Log error but continue with other users
				fmt.Printf("Failed to generate bill for user %s: %v\n", userID, err)
				continue
			}
		}
	}

	return nil
}

// calculateVolumeDiscount calculates volume discount based on total cost
func (abs *APIBillingService) calculateVolumeDiscount(totalCost decimal.Decimal) float64 {
	cost := totalCost.InexactFloat64()

	if cost >= 1000 {
		return 15.0 // 15% discount for $1000+
	} else if cost >= 500 {
		return 10.0 // 10% discount for $500+
	} else if cost >= 100 {
		return 5.0 // 5% discount for $100+
	}

	return 0.0 // No discount
}

// createStripeInvoice creates a Stripe invoice for API usage
func (abs *APIBillingService) createStripeInvoice(ctx context.Context, userID string, billingRecord *APIBillingRecord) (string, error) {
	// Get customer ID (would be stored in user record)
	customerID := fmt.Sprintf("cus_%s", userID) // Mock customer ID

	// Create invoice item for API usage
	invoiceItemParams := &stripe.InvoiceItemParams{
		Customer:    stripe.String(customerID),
		Amount:      stripe.Int64(billingRecord.FinalAmount.Mul(decimal.NewFromInt(100)).IntPart()), // Convert to cents
		Currency:    stripe.String("usd"),
		Description: stripe.String(fmt.Sprintf("API Usage - %s", billingRecord.BillingPeriod.Format("January 2006"))),
		Metadata: map[string]string{
			"user_id":        userID,
			"billing_period": billingRecord.BillingPeriod.Format("2006-01-01"),
			"total_requests": fmt.Sprintf("%d", billingRecord.TotalRequests),
		},
	}

	_, err := invoiceitem.New(invoiceItemParams)
	if err != nil {
		return "", fmt.Errorf("failed to create invoice item: %v", err)
	}

	// Create invoice
	invoiceParams := &stripe.InvoiceParams{
		Customer:    stripe.String(customerID),
		Description: stripe.String(fmt.Sprintf("API Usage Bill - %s", billingRecord.BillingPeriod.Format("January 2006"))),
		DueDate:     stripe.Int64(billingRecord.DueDate.Unix()),
		Metadata: map[string]string{
			"user_id":        userID,
			"billing_period": billingRecord.BillingPeriod.Format("2006-01-01"),
			"service":        "api_marketplace",
		},
	}

	stripeInvoice, err := invoice.New(invoiceParams)
	if err != nil {
		return "", fmt.Errorf("failed to create invoice: %v", err)
	}

	// Finalize invoice to make it payable
	_, err = invoice.FinalizeInvoice(stripeInvoice.ID, nil)
	if err != nil {
		return "", fmt.Errorf("failed to finalize invoice: %v", err)
	}

	return stripeInvoice.ID, nil
}

// APIBillingRecord represents a billing record
type APIBillingRecord struct {
	ID              string          `json:"id"`
	UserID          string          `json:"user_id"`
	BillingPeriod   time.Time       `json:"billing_period"`
	TotalRequests   int64           `json:"total_requests"`
	TotalCost       decimal.Decimal `json:"total_cost"`
	DiscountPercent decimal.Decimal `json:"discount_percent"`
	DiscountAmount  decimal.Decimal `json:"discount_amount"`
	FinalAmount     decimal.Decimal `json:"final_amount"`
	Currency        string          `json:"currency"`
	Status          string          `json:"status"`
	StripeInvoiceID string          `json:"stripe_invoice_id,omitempty"`
	GeneratedAt     time.Time       `json:"generated_at"`
	PaidAt          *time.Time      `json:"paid_at,omitempty"`
	DueDate         time.Time       `json:"due_date"`
}

// storeBillingRecord stores billing record in database
func (abs *APIBillingService) storeBillingRecord(ctx context.Context, record *APIBillingRecord) error {
	query := `
		INSERT INTO api_billing_records (
			user_id, billing_period, total_requests, total_cost, discount_percent,
			discount_amount, final_amount, currency, status, generated_at, due_date
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`

	err := abs.db.QueryRowContext(ctx, query,
		record.UserID, record.BillingPeriod, record.TotalRequests, record.TotalCost,
		record.DiscountPercent, record.DiscountAmount, record.FinalAmount,
		record.Currency, record.Status, record.GeneratedAt, record.DueDate,
	).Scan(&record.ID)

	return err
}

// updateBillingRecord updates billing record in database
func (abs *APIBillingService) updateBillingRecord(ctx context.Context, record *APIBillingRecord) error {
	query := `
		UPDATE api_billing_records 
		SET stripe_invoice_id = $1, status = $2
		WHERE id = $3
	`

	_, err := abs.db.ExecContext(ctx, query, record.StripeInvoiceID, record.Status, record.ID)
	return err
}

// getUsersWithAPIUsage gets all users with API usage in billing period
func (abs *APIBillingService) getUsersWithAPIUsage(ctx context.Context, billingPeriod time.Time) ([]string, error) {
	query := `
		SELECT DISTINCT user_id 
		FROM api_usage_summary 
		WHERE billing_period = $1 AND total_requests > 0
	`

	rows, err := abs.db.QueryContext(ctx, query, billingPeriod.Format("2006-01-01"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		users = append(users, userID)
	}

	return users, nil
}

// billingRecordExists checks if billing record exists for user and period
func (abs *APIBillingService) billingRecordExists(ctx context.Context, userID string, billingPeriod time.Time) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM api_billing_records 
			WHERE user_id = $1 AND billing_period = $2
		)
	`

	var exists bool
	err := abs.db.QueryRowContext(ctx, query, userID, billingPeriod.Format("2006-01-01")).Scan(&exists)
	return exists, err
}

// GetBillingHistory returns billing history for a user
func (abs *APIBillingService) GetBillingHistory(ctx context.Context, userID string, limit int) ([]*APIBillingRecord, error) {
	query := `
		SELECT id, user_id, billing_period, total_requests, total_cost,
		       discount_percent, discount_amount, final_amount, currency,
		       status, stripe_invoice_id, generated_at, paid_at, due_date
		FROM api_billing_records 
		WHERE user_id = $1 
		ORDER BY billing_period DESC 
		LIMIT $2
	`

	rows, err := abs.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*APIBillingRecord
	for rows.Next() {
		record := &APIBillingRecord{}
		var stripeInvoiceID sql.NullString
		var paidAt sql.NullTime

		err := rows.Scan(
			&record.ID, &record.UserID, &record.BillingPeriod, &record.TotalRequests,
			&record.TotalCost, &record.DiscountPercent, &record.DiscountAmount,
			&record.FinalAmount, &record.Currency, &record.Status,
			&stripeInvoiceID, &record.GeneratedAt, &paidAt, &record.DueDate,
		)
		if err != nil {
			return nil, err
		}

		if stripeInvoiceID.Valid {
			record.StripeInvoiceID = stripeInvoiceID.String
		}
		if paidAt.Valid {
			record.PaidAt = &paidAt.Time
		}

		records = append(records, record)
	}

	return records, nil
}
