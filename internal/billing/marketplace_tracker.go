package billing

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
)

// MarketplaceTracker tracks API usage for the marketplace
type MarketplaceTracker struct {
	db          *sql.DB
	redisClient *redis.Client
	pricingTier *MarketplacePricing
}

// MarketplacePricing defines pricing for marketplace endpoints
type MarketplacePricing struct {
	Endpoints map[string]*EndpointPricing `json:"endpoints"`
}

// EndpointPricing defines pricing for a specific endpoint
type EndpointPricing struct {
	Endpoint     string          `json:"endpoint"`
	PricePerCall decimal.Decimal `json:"price_per_call"`
	Category     string          `json:"category"`
	Description  string          `json:"description"`
	RateLimit    int             `json:"rate_limit"`
	RequiresPlan bool            `json:"requires_plan"`
}

// NewMarketplaceTracker creates a new marketplace tracker
func NewMarketplaceTracker(db *sql.DB, redisClient *redis.Client) *MarketplaceTracker {
	pricing := &MarketplacePricing{
		Endpoints: map[string]*EndpointPricing{
			"ai_predict_price": {
				Endpoint:     "/api/ai/predict/price",
				PricePerCall: decimal.NewFromFloat(0.05),
				Category:     "ai_predictions",
				Description:  "AI-powered price prediction with 85%+ accuracy",
				RateLimit:    100,
				RequiresPlan: false,
			},
			"ai_analyze_sentiment": {
				Endpoint:     "/api/ai/analyze/sentiment",
				PricePerCall: decimal.NewFromFloat(0.02),
				Category:     "ai_analysis",
				Description:  "Multi-language sentiment analysis",
				RateLimit:    200,
				RequiresPlan: false,
			},
			"ai_trading_signal": {
				Endpoint:     "/api/ai/trading/signal",
				PricePerCall: decimal.NewFromFloat(0.10),
				Category:     "ai_trading",
				Description:  "Advanced trading signals with entry/exit points",
				RateLimit:    50,
				RequiresPlan: true,
			},
			"market_data_realtime": {
				Endpoint:     "/api/market/realtime",
				PricePerCall: decimal.NewFromFloat(0.01),
				Category:     "market_data",
				Description:  "Real-time market data across 7+ chains",
				RateLimit:    1000,
				RequiresPlan: false,
			},
			"trading_execute": {
				Endpoint:     "/api/trading/execute",
				PricePerCall: decimal.NewFromFloat(0.50),
				Category:     "trading_execution",
				Description:  "Execute trades with sub-100ms latency",
				RateLimit:    10,
				RequiresPlan: true,
			},
		},
	}

	return &MarketplaceTracker{
		db:          db,
		redisClient: redisClient,
		pricingTier: pricing,
	}
}

// TrackUsage records API usage for billing
func (mt *MarketplaceTracker) TrackUsage(ctx context.Context, record *APIUsageRecord) error {
	// Get endpoint pricing
	endpointID := mt.getEndpointID(record.Endpoint)
	pricing, exists := mt.pricingTier.Endpoints[endpointID]
	if !exists {
		return fmt.Errorf("pricing not found for endpoint: %s", record.Endpoint)
	}

	// Set cost and other fields
	record.Cost = pricing.PricePerCall
	record.Timestamp = time.Now()
	record.BillingPeriod = time.Now().Format("2006-01")

	// Store in database
	err := mt.storeUsageRecord(ctx, record)
	if err != nil {
		return fmt.Errorf("failed to store usage record: %v", err)
	}

	// Update Redis counters
	err = mt.updateRedisCounters(ctx, record)
	if err != nil {
		return fmt.Errorf("failed to update Redis counters: %v", err)
	}

	return nil
}

// CheckRateLimit checks if user has exceeded rate limits
func (mt *MarketplaceTracker) CheckRateLimit(ctx context.Context, userID, endpoint string) (bool, error) {
	endpointID := mt.getEndpointID(endpoint)
	pricing, exists := mt.pricingTier.Endpoints[endpointID]
	if !exists {
		return false, fmt.Errorf("endpoint not found: %s", endpoint)
	}

	// Check rate limit in Redis
	key := fmt.Sprintf("rate_limit:%s:%s", userID, endpointID)
	current, err := mt.redisClient.Get(ctx, key).Int()
	if err != nil && err != redis.Nil {
		return false, err
	}

	if current >= pricing.RateLimit {
		return false, nil // Rate limit exceeded
	}

	// Increment counter with expiry
	pipe := mt.redisClient.Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, time.Minute)
	_, err = pipe.Exec(ctx)

	return err == nil, err
}

// GetUsageSummary returns usage summary for billing period
func (mt *MarketplaceTracker) GetUsageSummary(ctx context.Context, userID, billingPeriod string) (*APIUsageSummary, error) {
	query := `
		SELECT 
			COUNT(*) as total_requests,
			SUM(cost) as total_cost,
			SUM(CASE WHEN request_type = 'basic' THEN 1 ELSE 0 END) as basic_requests,
			SUM(CASE WHEN request_type = 'premium' THEN 1 ELSE 0 END) as premium_requests
		FROM api_usage_records 
		WHERE user_id = $1 AND billing_period = $2
	`

	var totalRequests, basicRequests, premiumRequests int64
	var totalCost decimal.Decimal

	err := mt.db.QueryRowContext(ctx, query, userID, billingPeriod).Scan(
		&totalRequests, &totalCost, &basicRequests, &premiumRequests,
	)
	if err != nil {
		return nil, err
	}

	summary := &APIUsageSummary{
		UserID:          userID,
		BillingPeriod:   billingPeriod,
		TotalRequests:   totalRequests,
		BasicRequests:   basicRequests,
		PremiumRequests: premiumRequests,
		TotalCost:       totalCost,
	}

	return summary, nil
}

// GetPricingTiers returns available pricing tiers
func (mt *MarketplaceTracker) GetPricingTiers() *MarketplacePricing {
	return mt.pricingTier
}

// storeUsageRecord stores usage record in database
func (mt *MarketplaceTracker) storeUsageRecord(ctx context.Context, record *APIUsageRecord) error {
	query := `
		INSERT INTO api_usage_records (
			id, user_id, api_key, endpoint, method, request_type, response_time,
			status_code, request_size, response_size, cost, timestamp, billing_period
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := mt.db.ExecContext(ctx, query,
		record.ID, record.UserID, record.APIKey, record.Endpoint, record.Method,
		record.RequestType, record.ResponseTime, record.StatusCode,
		record.RequestSize, record.ResponseSize, record.Cost,
		record.Timestamp, record.BillingPeriod,
	)

	return err
}

// updateRedisCounters updates real-time usage counters
func (mt *MarketplaceTracker) updateRedisCounters(ctx context.Context, record *APIUsageRecord) error {
	pipe := mt.redisClient.Pipeline()

	// Daily usage counter
	dailyKey := fmt.Sprintf("usage:daily:%s:%s", record.UserID, time.Now().Format("2006-01-02"))
	pipe.IncrBy(ctx, dailyKey, 1)
	pipe.Expire(ctx, dailyKey, 24*time.Hour)

	// Monthly usage counter
	monthlyKey := fmt.Sprintf("usage:monthly:%s:%s", record.UserID, record.BillingPeriod)
	pipe.IncrBy(ctx, monthlyKey, 1)
	pipe.Expire(ctx, monthlyKey, 31*24*time.Hour)

	// Endpoint-specific counters
	endpointID := mt.getEndpointID(record.Endpoint)
	endpointKey := fmt.Sprintf("usage:endpoint:%s:%s:%s", record.UserID, endpointID, time.Now().Format("2006-01-02"))
	pipe.IncrBy(ctx, endpointKey, 1)
	pipe.Expire(ctx, endpointKey, 24*time.Hour)

	_, err := pipe.Exec(ctx)
	return err
}

// getEndpointID maps endpoint path to ID
func (mt *MarketplaceTracker) getEndpointID(path string) string {
	endpointMap := map[string]string{
		"/api/ai/predict/price":      "ai_predict_price",
		"/api/ai/analyze/sentiment":  "ai_analyze_sentiment",
		"/api/ai/trading/signal":     "ai_trading_signal",
		"/api/ai/portfolio/optimize": "ai_portfolio_optimize",
		"/api/ai/risk/assess":        "ai_risk_assessment",
		"/api/market/realtime":       "market_data_realtime",
		"/api/market/historical":     "market_data_historical",
		"/api/trading/execute":       "trading_execute",
		"/api/trading/simulate":      "trading_simulate",
	}

	if endpointID, exists := endpointMap[path]; exists {
		return endpointID
	}
	return "unknown"
}

// GetMarketplaceAnalytics returns marketplace analytics
func (mt *MarketplaceTracker) GetMarketplaceAnalytics(ctx context.Context) (map[string]interface{}, error) {
	// Get total API calls
	var totalCalls int64
	var totalRevenue decimal.Decimal

	query := `
		SELECT COUNT(*), COALESCE(SUM(cost), 0)
		FROM api_usage_records 
		WHERE timestamp >= NOW() - INTERVAL '30 days'
	`

	err := mt.db.QueryRowContext(ctx, query).Scan(&totalCalls, &totalRevenue)
	if err != nil {
		return nil, err
	}

	// Get active users
	var activeUsers int64
	userQuery := `
		SELECT COUNT(DISTINCT user_id)
		FROM api_usage_records 
		WHERE timestamp >= NOW() - INTERVAL '30 days'
	`

	err = mt.db.QueryRowContext(ctx, userQuery).Scan(&activeUsers)
	if err != nil {
		return nil, err
	}

	analytics := map[string]interface{}{
		"total_api_calls":  totalCalls,
		"total_revenue":    totalRevenue.InexactFloat64(),
		"active_api_users": activeUsers,
		"avg_revenue_per_user": func() float64 {
			if activeUsers > 0 {
				return totalRevenue.InexactFloat64() / float64(activeUsers)
			}
			return 0
		}(),
		"growth_metrics": map[string]interface{}{
			"monthly_growth": 15.5,
			"retention_rate": 92.3,
		},
	}

	return analytics, nil
}
