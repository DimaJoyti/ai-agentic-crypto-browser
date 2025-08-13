package billing

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// APIUsageManager handles API usage tracking and billing
type APIUsageManager struct {
	db     *sql.DB
	config *APIBillingConfig
}

// APIBillingConfig defines API pricing structure
type APIBillingConfig struct {
	BasicAPIPrice    decimal.Decimal `json:"basic_api_price"`    // Price per request
	PremiumAPIPrice  decimal.Decimal `json:"premium_api_price"`  // Price per premium request
	FreeRequestLimit int             `json:"free_request_limit"` // Free requests per month
	BillingPeriod    string          `json:"billing_period"`     // monthly, daily
}

// APIUsageRecord tracks API usage
type APIUsageRecord struct {
	ID            string          `json:"id"`
	UserID        string          `json:"user_id"`
	APIKey        string          `json:"api_key"`
	Endpoint      string          `json:"endpoint"`
	Method        string          `json:"method"`
	RequestType   string          `json:"request_type"`  // basic, premium, enterprise
	ResponseTime  int64           `json:"response_time"` // milliseconds
	StatusCode    int             `json:"status_code"`
	RequestSize   int64           `json:"request_size"`  // bytes
	ResponseSize  int64           `json:"response_size"` // bytes
	Cost          decimal.Decimal `json:"cost"`
	Timestamp     time.Time       `json:"timestamp"`
	BillingPeriod string          `json:"billing_period"` // 2024-01, 2024-02, etc.
	Duration      time.Duration   `json:"duration"`
	Success       bool            `json:"success"`
	ErrorCode     string          `json:"error_code,omitempty"`
	Metadata      json.RawMessage `json:"metadata,omitempty"`
}

// APIUsageSummary provides usage summary for billing
type APIUsageSummary struct {
	UserID           string          `json:"user_id"`
	BillingPeriod    string          `json:"billing_period"`
	BasicRequests    int64           `json:"basic_requests"`
	PremiumRequests  int64           `json:"premium_requests"`
	TotalRequests    int64           `json:"total_requests"`
	FreeRequests     int64           `json:"free_requests"`
	BillableRequests int64           `json:"billable_requests"`
	TotalCost        decimal.Decimal `json:"total_cost"`
	AverageLatency   float64         `json:"average_latency"`
	ErrorRate        float64         `json:"error_rate"`
}

// NewAPIUsageManager creates a new API usage manager
func NewAPIUsageManager(db *sql.DB, config *APIBillingConfig) *APIUsageManager {
	return &APIUsageManager{
		db:     db,
		config: config,
	}
}

// TrackAPIUsage records an API request for billing
func (aum *APIUsageManager) TrackAPIUsage(ctx context.Context, usage *APIUsageRecord) error {
	// Calculate cost based on request type
	var cost decimal.Decimal
	switch usage.RequestType {
	case "basic":
		cost = aum.config.BasicAPIPrice
	case "premium":
		cost = aum.config.PremiumAPIPrice
	case "enterprise":
		// Enterprise pricing is custom, set to zero for now
		cost = decimal.Zero
	default:
		cost = aum.config.BasicAPIPrice
	}

	usage.Cost = cost
	usage.BillingPeriod = time.Now().Format("2006-01")

	// Save usage record
	query := `
		INSERT INTO api_usage_records (
			id, user_id, api_key, endpoint, method, request_type,
			response_time, status_code, request_size, response_size,
			cost, timestamp, billing_period
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := aum.db.ExecContext(ctx, query,
		usage.ID, usage.UserID, usage.APIKey, usage.Endpoint,
		usage.Method, usage.RequestType, usage.ResponseTime,
		usage.StatusCode, usage.RequestSize, usage.ResponseSize,
		usage.Cost, usage.Timestamp, usage.BillingPeriod,
	)

	return err
}

// GetUsageSummary returns usage summary for billing period
func (aum *APIUsageManager) GetUsageSummary(ctx context.Context, userID, billingPeriod string) (*APIUsageSummary, error) {
	query := `
		SELECT 
			COUNT(*) as total_requests,
			COUNT(CASE WHEN request_type = 'basic' THEN 1 END) as basic_requests,
			COUNT(CASE WHEN request_type = 'premium' THEN 1 END) as premium_requests,
			COALESCE(SUM(cost), 0) as total_cost,
			COALESCE(AVG(response_time), 0) as avg_latency,
			COALESCE(AVG(CASE WHEN status_code >= 400 THEN 1.0 ELSE 0.0 END), 0) as error_rate
		FROM api_usage_records 
		WHERE user_id = $1 AND billing_period = $2
	`

	row := aum.db.QueryRowContext(ctx, query, userID, billingPeriod)

	summary := &APIUsageSummary{
		UserID:        userID,
		BillingPeriod: billingPeriod,
	}

	err := row.Scan(
		&summary.TotalRequests,
		&summary.BasicRequests,
		&summary.PremiumRequests,
		&summary.TotalCost,
		&summary.AverageLatency,
		&summary.ErrorRate,
	)
	if err != nil {
		return nil, err
	}

	// Calculate free vs billable requests
	if summary.TotalRequests <= int64(aum.config.FreeRequestLimit) {
		summary.FreeRequests = summary.TotalRequests
		summary.BillableRequests = 0
		summary.TotalCost = decimal.Zero
	} else {
		summary.FreeRequests = int64(aum.config.FreeRequestLimit)
		summary.BillableRequests = summary.TotalRequests - summary.FreeRequests
	}

	return summary, nil
}

// GenerateAPIBill creates a bill for API usage
func (aum *APIUsageManager) GenerateAPIBill(ctx context.Context, userID, billingPeriod string) (*APIBill, error) {
	summary, err := aum.GetUsageSummary(ctx, userID, billingPeriod)
	if err != nil {
		return nil, fmt.Errorf("failed to get usage summary: %w", err)
	}

	bill := &APIBill{
		ID:               generateBillID(),
		UserID:           userID,
		BillingPeriod:    billingPeriod,
		BasicRequests:    summary.BasicRequests,
		PremiumRequests:  summary.PremiumRequests,
		TotalRequests:    summary.TotalRequests,
		FreeRequests:     summary.FreeRequests,
		BillableRequests: summary.BillableRequests,
		BasicCost:        decimal.NewFromInt(summary.BasicRequests).Mul(aum.config.BasicAPIPrice),
		PremiumCost:      decimal.NewFromInt(summary.PremiumRequests).Mul(aum.config.PremiumAPIPrice),
		TotalCost:        summary.TotalCost,
		Status:           "generated",
		GeneratedAt:      time.Now(),
	}

	// Apply free tier discount
	if bill.BillableRequests > 0 {
		freeDiscount := decimal.NewFromInt(summary.FreeRequests).Mul(aum.config.BasicAPIPrice)
		bill.TotalCost = bill.TotalCost.Sub(freeDiscount)
	}

	// Save bill to database
	err = aum.saveAPIBill(ctx, bill)
	if err != nil {
		return nil, fmt.Errorf("failed to save API bill: %w", err)
	}

	return bill, nil
}

// APIBill represents a billing statement for API usage
type APIBill struct {
	ID               string          `json:"id"`
	UserID           string          `json:"user_id"`
	BillingPeriod    string          `json:"billing_period"`
	BasicRequests    int64           `json:"basic_requests"`
	PremiumRequests  int64           `json:"premium_requests"`
	TotalRequests    int64           `json:"total_requests"`
	FreeRequests     int64           `json:"free_requests"`
	BillableRequests int64           `json:"billable_requests"`
	BasicCost        decimal.Decimal `json:"basic_cost"`
	PremiumCost      decimal.Decimal `json:"premium_cost"`
	TotalCost        decimal.Decimal `json:"total_cost"`
	Status           string          `json:"status"` // generated, sent, paid, overdue
	GeneratedAt      time.Time       `json:"generated_at"`
	PaidAt           *time.Time      `json:"paid_at,omitempty"`
}

// GetTopAPIUsers returns users with highest API usage
func (aum *APIUsageManager) GetTopAPIUsers(ctx context.Context, billingPeriod string, limit int) ([]*APIUsageSummary, error) {
	query := `
		SELECT 
			user_id,
			COUNT(*) as total_requests,
			COUNT(CASE WHEN request_type = 'basic' THEN 1 END) as basic_requests,
			COUNT(CASE WHEN request_type = 'premium' THEN 1 END) as premium_requests,
			COALESCE(SUM(cost), 0) as total_cost,
			COALESCE(AVG(response_time), 0) as avg_latency,
			COALESCE(AVG(CASE WHEN status_code >= 400 THEN 1.0 ELSE 0.0 END), 0) as error_rate
		FROM api_usage_records 
		WHERE billing_period = $1
		GROUP BY user_id
		ORDER BY total_cost DESC
		LIMIT $2
	`

	rows, err := aum.db.QueryContext(ctx, query, billingPeriod, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []*APIUsageSummary
	for rows.Next() {
		summary := &APIUsageSummary{
			BillingPeriod: billingPeriod,
		}

		err := rows.Scan(
			&summary.UserID,
			&summary.TotalRequests,
			&summary.BasicRequests,
			&summary.PremiumRequests,
			&summary.TotalCost,
			&summary.AverageLatency,
			&summary.ErrorRate,
		)
		if err != nil {
			return nil, err
		}

		summaries = append(summaries, summary)
	}

	return summaries, nil
}

// Helper functions
func (aum *APIUsageManager) saveAPIBill(ctx context.Context, bill *APIBill) error {
	query := `
		INSERT INTO api_bills (
			id, user_id, billing_period, basic_requests, premium_requests,
			total_requests, free_requests, billable_requests,
			basic_cost, premium_cost, total_cost, status, generated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := aum.db.ExecContext(ctx, query,
		bill.ID, bill.UserID, bill.BillingPeriod,
		bill.BasicRequests, bill.PremiumRequests, bill.TotalRequests,
		bill.FreeRequests, bill.BillableRequests,
		bill.BasicCost, bill.PremiumCost, bill.TotalCost,
		bill.Status, bill.GeneratedAt,
	)

	return err
}

func generateBillID() string {
	return fmt.Sprintf("bill_%d", time.Now().UnixNano())
}
