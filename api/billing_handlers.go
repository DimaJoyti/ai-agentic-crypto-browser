package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ai-agentic-browser/internal/billing"
	"github.com/gorilla/mux"
)

// BillingHandlers handles billing-related HTTP requests
type BillingHandlers struct {
	subscriptionManager   *billing.SubscriptionManager
	performanceFeeManager *billing.PerformanceFeeManager
	apiUsageManager       *billing.APIUsageManager
}

// NewBillingHandlers creates new billing handlers
func NewBillingHandlers(
	subscriptionManager *billing.SubscriptionManager,
	performanceFeeManager *billing.PerformanceFeeManager,
	apiUsageManager *billing.APIUsageManager,
) *BillingHandlers {
	return &BillingHandlers{
		subscriptionManager:   subscriptionManager,
		performanceFeeManager: performanceFeeManager,
		apiUsageManager:       apiUsageManager,
	}
}

// RegisterRoutes registers billing routes
func (bh *BillingHandlers) RegisterRoutes(router *mux.Router) {
	// Subscription routes
	router.HandleFunc("/billing/subscriptions", bh.CreateSubscription).Methods("POST")
	router.HandleFunc("/billing/subscriptions", bh.GetUserSubscription).Methods("GET")
	router.HandleFunc("/billing/subscriptions/{id}/cancel", bh.CancelSubscription).Methods("POST")
	router.HandleFunc("/billing/subscriptions/tiers", bh.GetSubscriptionTiers).Methods("GET")

	// Performance fee routes
	router.HandleFunc("/billing/performance-fees", bh.GetPerformanceFees).Methods("GET")
	router.HandleFunc("/billing/performance-fees/calculate", bh.CalculatePerformanceFees).Methods("POST")
	router.HandleFunc("/billing/performance-fees/{id}/charge", bh.ChargePerformanceFee).Methods("POST")

	// API usage routes
	router.HandleFunc("/billing/api-usage", bh.GetAPIUsage).Methods("GET")
	router.HandleFunc("/billing/api-usage/summary", bh.GetAPIUsageSummary).Methods("GET")
	router.HandleFunc("/billing/api-bills", bh.GetAPIBills).Methods("GET")
	router.HandleFunc("/billing/api-bills/generate", bh.GenerateAPIBill).Methods("POST")

	// Analytics routes
	router.HandleFunc("/billing/analytics/revenue", bh.GetRevenueAnalytics).Methods("GET")
	router.HandleFunc("/billing/analytics/users", bh.GetUserAnalytics).Methods("GET")
}

// CreateSubscriptionRequest represents subscription creation request
type CreateSubscriptionRequest struct {
	TierID       string `json:"tier_id"`
	BillingCycle string `json:"billing_cycle"` // monthly, annual
	TrialDays    int    `json:"trial_days"`
}

// CreateSubscription creates a new subscription
func (bh *BillingHandlers) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	subscription, err := bh.subscriptionManager.CreateSubscription(
		r.Context(), userID, req.TierID, req.BillingCycle, req.TrialDays,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscription)
}

// GetUserSubscription returns user's current subscription
func (bh *BillingHandlers) GetUserSubscription(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	subscription, err := bh.subscriptionManager.GetUserSubscription(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if subscription == nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"subscription": nil})
	} else {
		json.NewEncoder(w).Encode(subscription)
	}
}

// GetSubscriptionTiers returns available subscription tiers
func (bh *BillingHandlers) GetSubscriptionTiers(w http.ResponseWriter, r *http.Request) {
	// This would typically come from the subscription manager
	tiers := map[string]interface{}{
		"starter": map[string]interface{}{
			"id":           "starter",
			"name":         "Starter",
			"price":        49,
			"annual_price": 490,
			"features": []string{
				"Basic AI Trading",
				"3 Trading Strategies",
				"Single Chain Support",
				"Basic Analytics",
				"Email Support",
			},
			"popular": false,
		},
		"professional": map[string]interface{}{
			"id":           "professional",
			"name":         "Professional",
			"price":        199,
			"annual_price": 1990,
			"features": []string{
				"Advanced AI Trading",
				"10+ Trading Strategies",
				"Multi-Chain Support",
				"Advanced Analytics",
				"DeFi Integration",
				"Voice Commands",
				"Priority Support",
			},
			"popular": true,
		},
		"enterprise": map[string]interface{}{
			"id":           "enterprise",
			"name":         "Enterprise",
			"price":        999,
			"annual_price": 9990,
			"features": []string{
				"Full Platform Access",
				"Unlimited Strategies",
				"All Chains Supported",
				"Custom AI Models",
				"White-Label Solution",
				"Dedicated Support",
				"Custom Integrations",
				"SLA Guarantee",
			},
			"popular": false,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tiers)
}

// GetPerformanceFees returns user's performance fee history
func (bh *BillingHandlers) GetPerformanceFees(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 50 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	fees, err := bh.performanceFeeManager.GetPerformanceFeeHistory(r.Context(), userID, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fees)
}

// CalculatePerformanceFeesRequest represents performance fee calculation request
type CalculatePerformanceFeesRequest struct {
	StrategyID  string    `json:"strategy_id"`
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
}

// CalculatePerformanceFees calculates performance fees for a period
func (bh *BillingHandlers) CalculatePerformanceFees(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CalculatePerformanceFeesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	record, err := bh.performanceFeeManager.CalculatePerformanceFee(
		r.Context(), userID, req.StrategyID, req.PeriodStart, req.PeriodEnd,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(record)
}

// ChargePerformanceFee processes a performance fee charge
func (bh *BillingHandlers) ChargePerformanceFee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recordID := vars["id"]

	err := bh.performanceFeeManager.ChargePerformanceFee(r.Context(), recordID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "charged"})
}

// GetAPIUsageSummary returns API usage summary
func (bh *BillingHandlers) GetAPIUsageSummary(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	billingPeriod := r.URL.Query().Get("period")
	if billingPeriod == "" {
		billingPeriod = time.Now().Format("2006-01") // Current month
	}

	summary, err := bh.apiUsageManager.GetUsageSummary(r.Context(), userID, billingPeriod)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// GetAPIUsage returns detailed API usage records
func (bh *BillingHandlers) GetAPIUsage(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Implementation would fetch detailed usage records
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "API usage details"})
}

// GetAPIBills returns API billing history
func (bh *BillingHandlers) GetAPIBills(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Implementation would fetch billing history
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "API bills"})
}

// GenerateAPIBill generates a new API bill
func (bh *BillingHandlers) GenerateAPIBill(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	billingPeriod := r.URL.Query().Get("period")
	if billingPeriod == "" {
		billingPeriod = time.Now().AddDate(0, -1, 0).Format("2006-01") // Previous month
	}

	bill, err := bh.apiUsageManager.GenerateAPIBill(r.Context(), userID, billingPeriod)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bill)
}

// GetRevenueAnalytics returns revenue analytics (admin only)
func (bh *BillingHandlers) GetRevenueAnalytics(w http.ResponseWriter, r *http.Request) {
	// Check admin permissions
	if !isAdmin(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Mock revenue analytics
	analytics := map[string]interface{}{
		"total_revenue":        1250000,
		"monthly_revenue":      125000,
		"subscription_revenue": 100000,
		"performance_fees":     20000,
		"api_revenue":          5000,
		"growth_rate":          0.15,
		"churn_rate":           0.05,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analytics)
}

// GetUserAnalytics returns user analytics (admin only)
func (bh *BillingHandlers) GetUserAnalytics(w http.ResponseWriter, r *http.Request) {
	// Check admin permissions
	if !isAdmin(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Mock user analytics
	analytics := map[string]interface{}{
		"total_users":        5000,
		"active_subscribers": 1200,
		"trial_users":        300,
		"enterprise_clients": 50,
		"conversion_rate":    0.24,
		"ltv":                2400,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analytics)
}

// Helper functions
func getUserIDFromContext(ctx context.Context) string {
	// Implementation would extract user ID from JWT token
	return "user123" // Mock for now
}

func isAdmin(ctx context.Context) bool {
	// Implementation would check admin role
	return true // Mock for now
}

// CancelSubscription cancels a user's subscription
func (bh *BillingHandlers) CancelSubscription(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Implementation would cancel the subscription
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "cancelled"})
}
