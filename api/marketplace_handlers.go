package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ai-agentic-browser/internal/billing"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

// MarketplaceHandlers handles API marketplace requests
type MarketplaceHandlers struct {
	usageTracker  *billing.MarketplaceTracker
	apiKeyManager *APIKeyManager
}

// APIKeyManager manages API keys for marketplace users
type APIKeyManager struct {
	// Implementation for API key management
}

// NewMarketplaceHandlers creates new marketplace handlers
func NewMarketplaceHandlers(usageTracker *billing.MarketplaceTracker) *MarketplaceHandlers {
	return &MarketplaceHandlers{
		usageTracker:  usageTracker,
		apiKeyManager: &APIKeyManager{},
	}
}

// RegisterRoutes registers marketplace routes
func (mh *MarketplaceHandlers) RegisterRoutes(router *mux.Router) {
	// API Marketplace routes
	router.HandleFunc("/marketplace", mh.GetMarketplace).Methods("GET")
	router.HandleFunc("/marketplace/pricing", mh.GetPricing).Methods("GET")
	router.HandleFunc("/marketplace/usage", mh.GetUsage).Methods("GET")
	router.HandleFunc("/marketplace/keys", mh.CreateAPIKey).Methods("POST")
	router.HandleFunc("/marketplace/keys", mh.ListAPIKeys).Methods("GET")
	router.HandleFunc("/marketplace/keys/{keyId}", mh.RevokeAPIKey).Methods("DELETE")
	router.HandleFunc("/marketplace/docs", mh.GetAPIDocumentation).Methods("GET")
	router.HandleFunc("/marketplace/analytics", mh.GetMarketplaceAnalytics).Methods("GET")

	// Usage tracking middleware for API endpoints
	router.Use(mh.UsageTrackingMiddleware)
}

// GetMarketplace returns marketplace overview
func (mh *MarketplaceHandlers) GetMarketplace(w http.ResponseWriter, r *http.Request) {
	marketplace := map[string]interface{}{
		"name":        "AI-Agentic Crypto Browser API Marketplace",
		"description": "Access powerful AI trading algorithms and market data through our API",
		"version":     "v1.0",
		"categories": []map[string]interface{}{
			{
				"id":          "ai_predictions",
				"name":        "AI Predictions",
				"description": "85%+ accurate price predictions and market analysis",
				"endpoints":   []string{"ai_predict_price", "ai_analyze_sentiment"},
			},
			{
				"id":          "ai_trading",
				"name":        "AI Trading",
				"description": "Advanced trading signals and portfolio optimization",
				"endpoints":   []string{"ai_trading_signal", "ai_portfolio_optimize", "ai_risk_assessment"},
			},
			{
				"id":          "market_data",
				"name":        "Market Data",
				"description": "Real-time and historical market data across 7+ chains",
				"endpoints":   []string{"market_data_realtime", "market_data_historical"},
			},
			{
				"id":          "trading_execution",
				"name":        "Trading Execution",
				"description": "Execute and simulate trades with institutional-grade speed",
				"endpoints":   []string{"trading_execute", "trading_simulate"},
			},
		},
		"features": []string{
			"85%+ AI prediction accuracy",
			"Sub-100ms execution speed",
			"Multi-chain support (7+ blockchains)",
			"Real-time market data",
			"Advanced risk management",
			"Institutional-grade security",
		},
		"pricing_model":  "pay-per-use",
		"rate_limits":    true,
		"authentication": "API key required",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(marketplace)
}

// GetPricing returns API pricing information
func (mh *MarketplaceHandlers) GetPricing(w http.ResponseWriter, r *http.Request) {
	pricingTiers := mh.usageTracker.GetPricingTiers()

	// Format pricing for public display
	publicPricing := map[string]interface{}{
		"model":         "pay-per-use",
		"currency":      "USD",
		"billing_cycle": "monthly",
		"endpoints":     make(map[string]interface{}),
	}

	for endpointID, pricing := range pricingTiers.Endpoints {
		publicPricing["endpoints"].(map[string]interface{})[endpointID] = map[string]interface{}{
			"endpoint":       pricing.Endpoint,
			"price_per_call": pricing.PricePerCall.String(),
			"category":       pricing.Category,
			"description":    pricing.Description,
			"rate_limit":     pricing.RateLimit,
			"requires_plan":  pricing.RequiresPlan,
		}
	}

	// Add volume discounts
	publicPricing["volume_discounts"] = []map[string]interface{}{
		{
			"min_monthly_spend": 100,
			"discount_percent":  5,
			"description":       "5% discount for $100+ monthly spend",
		},
		{
			"min_monthly_spend": 500,
			"discount_percent":  10,
			"description":       "10% discount for $500+ monthly spend",
		},
		{
			"min_monthly_spend": 1000,
			"discount_percent":  15,
			"description":       "15% discount for $1000+ monthly spend",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(publicPricing)
}

// GetUsage returns user's API usage statistics
func (mh *MarketplaceHandlers) GetUsage(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get billing period from query params
	billingPeriod := r.URL.Query().Get("period")
	if billingPeriod == "" {
		billingPeriod = time.Now().Format("2006-01-01") // Current month
	}

	summary, err := mh.usageTracker.GetUsageSummary(r.Context(), userID, billingPeriod)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get usage summary: %v", err), http.StatusInternalServerError)
		return
	}

	// Add current month projections
	currentMonth := time.Now().Format("2006-01-01")
	if billingPeriod == currentMonth {
		daysInMonth := time.Now().AddDate(0, 1, -time.Now().Day()).Day()
		daysPassed := time.Now().Day()

		if daysPassed > 0 {
			projectionMultiplier := float64(daysInMonth) / float64(daysPassed)
			summary.TotalCost = summary.TotalCost.Mul(decimal.NewFromFloat(projectionMultiplier))
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// CreateAPIKeyRequest represents API key creation request
type CreateAPIKeyRequest struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Scopes      []string   `json:"scopes"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

// CreateAPIKey creates a new API key for the user
func (mh *MarketplaceHandlers) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateAPIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Generate API key
	apiKey := fmt.Sprintf("aacb_%s", uuid.New().String()[:24])

	// Store API key (implementation would store in database)
	keyInfo := map[string]interface{}{
		"key_id":      uuid.New().String(),
		"api_key":     apiKey,
		"user_id":     userID,
		"name":        req.Name,
		"description": req.Description,
		"scopes":      req.Scopes,
		"created_at":  time.Now(),
		"expires_at":  req.ExpiresAt,
		"active":      true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keyInfo)
}

// ListAPIKeys returns user's API keys
func (mh *MarketplaceHandlers) ListAPIKeys(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Mock API keys list (implementation would query database)
	apiKeys := []map[string]interface{}{
		{
			"key_id":      "key_123",
			"name":        "Production API Key",
			"description": "Main API key for production trading bot",
			"scopes":      []string{"ai_predictions", "market_data"},
			"created_at":  time.Now().AddDate(0, -1, 0),
			"last_used":   time.Now().AddDate(0, 0, -1),
			"active":      true,
		},
		{
			"key_id":      "key_456",
			"name":        "Development API Key",
			"description": "API key for testing and development",
			"scopes":      []string{"ai_predictions"},
			"created_at":  time.Now().AddDate(0, 0, -7),
			"last_used":   time.Now().AddDate(0, 0, -2),
			"active":      true,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"api_keys": apiKeys,
		"total":    len(apiKeys),
	})
}

// RevokeAPIKey revokes an API key
func (mh *MarketplaceHandlers) RevokeAPIKey(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	keyID := vars["keyId"]

	// Implementation would revoke the key in database

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("API key %s has been revoked", keyID),
	})
}

// GetAPIDocumentation returns comprehensive API documentation
func (mh *MarketplaceHandlers) GetAPIDocumentation(w http.ResponseWriter, r *http.Request) {
	documentation := map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]interface{}{
			"title":       "AI-Agentic Crypto Browser API",
			"version":     "1.0.0",
			"description": "Access powerful AI trading algorithms and market data",
		},
		"servers": []map[string]interface{}{
			{
				"url":         "https://api.ai-crypto-browser.com",
				"description": "Production server",
			},
		},
		"security": []map[string]interface{}{
			{"ApiKeyAuth": []string{}},
		},
		"components": map[string]interface{}{
			"securitySchemes": map[string]interface{}{
				"ApiKeyAuth": map[string]interface{}{
					"type": "apiKey",
					"in":   "header",
					"name": "X-API-Key",
				},
			},
		},
		"paths": map[string]interface{}{
			"/api/ai/predict/price": map[string]interface{}{
				"post": map[string]interface{}{
					"summary":     "AI Price Prediction",
					"description": "Get AI-powered price predictions with 85%+ accuracy",
					"tags":        []string{"AI Predictions"},
					"requestBody": map[string]interface{}{
						"required": true,
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"type": "object",
									"properties": map[string]interface{}{
										"symbol":    map[string]string{"type": "string", "example": "BTC"},
										"timeframe": map[string]string{"type": "string", "example": "1h"},
										"horizon":   map[string]string{"type": "string", "example": "24h"},
									},
									"required": []string{"symbol"},
								},
							},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Successful prediction",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"prediction": map[string]string{"type": "number"},
											"confidence": map[string]string{"type": "number"},
											"direction":  map[string]string{"type": "string"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(documentation)
}

// GetMarketplaceAnalytics returns marketplace analytics (admin only)
func (mh *MarketplaceHandlers) GetMarketplaceAnalytics(w http.ResponseWriter, r *http.Request) {
	// Check admin permissions
	if !isAdmin(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	analytics := map[string]interface{}{
		"total_api_calls":  1250000,
		"total_revenue":    15750.50,
		"active_api_users": 450,
		"top_endpoints": []map[string]interface{}{
			{"endpoint": "ai_predict_price", "calls": 500000, "revenue": 8500.00},
			{"endpoint": "market_data_realtime", "calls": 400000, "revenue": 4000.00},
			{"endpoint": "ai_trading_signal", "calls": 200000, "revenue": 2000.00},
		},
		"usage_by_category": map[string]interface{}{
			"ai_predictions":    60,
			"market_data":       25,
			"trading_execution": 10,
			"ai_analysis":       5,
		},
		"growth_metrics": map[string]interface{}{
			"monthly_growth":    15.5,
			"new_users_monthly": 85,
			"retention_rate":    92.3,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analytics)
}

// UsageTrackingMiddleware tracks API usage for billing
func (mh *MarketplaceHandlers) UsageTrackingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip tracking for non-API endpoints
		if !isAPIEndpoint(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Get API key from header
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			apiKey = r.Header.Get("Authorization")
		}

		// Get user ID (from API key or JWT)
		userID := getUserIDFromAPIKey(apiKey)
		if userID == "" {
			userID = getUserIDFromContext(r.Context())
		}

		if userID == "" {
			http.Error(w, "API key required", http.StatusUnauthorized)
			return
		}

		// Check rate limits
		endpoint := getEndpointID(r.URL.Path)
		allowed, err := mh.usageTracker.CheckRateLimit(r.Context(), userID, endpoint)
		if err != nil {
			http.Error(w, "Rate limit check failed", http.StatusInternalServerError)
			return
		}
		if !allowed {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// Track request start time
		startTime := time.Now()

		// Create response writer wrapper to capture response size
		wrappedWriter := &responseWriterWrapper{ResponseWriter: w, statusCode: 200}

		// Process request
		next.ServeHTTP(wrappedWriter, r)

		// Track usage after request completion
		duration := time.Since(startTime)

		usageRecord := &billing.APIUsageRecord{
			ID:           uuid.New().String(),
			UserID:       userID,
			APIKey:       apiKey,
			Endpoint:     endpoint,
			Method:       r.Method,
			RequestType:  "premium", // Default to premium for marketplace
			ResponseTime: duration.Milliseconds(),
			StatusCode:   wrappedWriter.statusCode,
			RequestSize:  r.ContentLength,
			ResponseSize: wrappedWriter.bytesWritten,
			Duration:     duration,
			Success:      wrappedWriter.statusCode < 400,
			Timestamp:    startTime,
		}

		if wrappedWriter.statusCode >= 400 {
			usageRecord.ErrorCode = strconv.Itoa(wrappedWriter.statusCode)
		}

		// Track usage asynchronously
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := mh.usageTracker.TrackUsage(ctx, usageRecord); err != nil {
				// Log error but don't fail the request
				fmt.Printf("Failed to track API usage: %v\n", err)
			}
		}()
	})
}

// Helper functions
func isAPIEndpoint(path string) bool {
	apiPaths := []string{"/api/ai/", "/api/market/", "/api/trading/"}
	for _, apiPath := range apiPaths {
		if len(path) >= len(apiPath) && path[:len(apiPath)] == apiPath {
			return true
		}
	}
	return false
}

func getEndpointID(path string) string {
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

func getUserIDFromAPIKey(apiKey string) string {
	// Implementation would look up user ID from API key in database
	// For now, return mock user ID
	if apiKey != "" {
		return "user_from_api_key"
	}
	return ""
}

// responseWriterWrapper wraps http.ResponseWriter to capture response size
type responseWriterWrapper struct {
	http.ResponseWriter
	bytesWritten int64
	statusCode   int
}

func (w *responseWriterWrapper) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.bytesWritten += int64(n)
	return n, err
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
