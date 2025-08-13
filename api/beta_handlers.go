package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ai-agentic-browser/internal/billing"
	"github.com/gorilla/mux"
)

// BetaHandlers handles beta program related requests
type BetaHandlers struct {
	subscriptionManager *billing.SubscriptionManager
	stripeProcessor     *billing.StripePaymentProcessor
}

// NewBetaHandlers creates new beta handlers
func NewBetaHandlers(
	subscriptionManager *billing.SubscriptionManager,
	stripeProcessor *billing.StripePaymentProcessor,
) *BetaHandlers {
	return &BetaHandlers{
		subscriptionManager: subscriptionManager,
		stripeProcessor:     stripeProcessor,
	}
}

// RegisterRoutes registers beta program routes
func (bh *BetaHandlers) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/beta/signup", bh.BetaSignup).Methods("POST")
	router.HandleFunc("/beta/status", bh.GetBetaStatus).Methods("GET")
	router.HandleFunc("/beta/analytics", bh.GetBetaAnalytics).Methods("GET")
}

// BetaSignupRequest represents a beta signup request
type BetaSignupRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Tier  string `json:"tier"`
}

// BetaSignupResponse represents the response to a beta signup
type BetaSignupResponse struct {
	Success    bool     `json:"success"`
	Message    string   `json:"message"`
	UserID     string   `json:"user_id,omitempty"`
	CustomerID string   `json:"customer_id,omitempty"`
	NextSteps  []string `json:"next_steps,omitempty"`
}

// BetaSignup handles beta program signup
func (bh *BetaHandlers) BetaSignup(w http.ResponseWriter, r *http.Request) {
	var req BetaSignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Name == "" || req.Email == "" || req.Tier == "" {
		http.Error(w, "Name, email, and tier are required", http.StatusBadRequest)
		return
	}

	// Validate tier
	validTiers := map[string]bool{
		"starter":      true,
		"professional": true,
		"enterprise":   true,
	}
	if !validTiers[req.Tier] {
		http.Error(w, "Invalid tier", http.StatusBadRequest)
		return
	}

	// Generate user ID (in real implementation, this would come from user registration)
	userID := fmt.Sprintf("beta_user_%d", time.Now().Unix())

	// Create Stripe customer
	customer, err := bh.stripeProcessor.CreateCustomer(r.Context(), userID, req.Email, req.Name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create customer: %v", err), http.StatusInternalServerError)
		return
	}

	// Get beta pricing (50% off)
	betaPricing := map[string]int64{
		"starter":      2500,  // $25 (50% off $49)
		"professional": 9950,  // $99.50 (50% off $199)
		"enterprise":   49950, // $499.50 (50% off $999)
	}

	// Create subscription with beta pricing
	// Note: In real implementation, you'd create custom Stripe prices for beta
	priceAmount := betaPricing[req.Tier]

	// For now, we'll create a subscription record in our system
	// In production, you'd integrate with Stripe's subscription system
	_, err = bh.subscriptionManager.CreateSubscription(
		r.Context(), userID, req.Tier, "monthly", 7, // 7-day trial
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create subscription: %v", err), http.StatusInternalServerError)
		return
	}

	// Send welcome email (implement email service)
	go bh.sendBetaWelcomeEmail(req.Email, req.Name, req.Tier)

	// Prepare response
	response := BetaSignupResponse{
		Success:    true,
		Message:    "Successfully joined beta program!",
		UserID:     userID,
		CustomerID: customer.ID,
		NextSteps: []string{
			"Check your email for beta access credentials",
			"Join our exclusive Discord community",
			"Schedule your onboarding call",
			"Start trading with AI in 48 hours",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	// Log beta signup for analytics
	go bh.logBetaSignup(userID, req.Email, req.Tier, priceAmount)
}

// GetBetaStatus returns current beta program status
func (bh *BetaHandlers) GetBetaStatus(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"active":          true,
		"spots_available": true,
		"total_signups":   150, // Mock data
		"target_signups":  500,
		"discount":        "50%",
		"expires_at":      "2024-12-31T23:59:59Z",
		"features_enabled": []string{
			"ai_trading",
			"multi_chain",
			"voice_commands",
			"advanced_analytics",
			"performance_fees",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// GetBetaAnalytics returns beta program analytics (admin only)
func (bh *BetaHandlers) GetBetaAnalytics(w http.ResponseWriter, r *http.Request) {
	// Check admin permissions
	if !isAdmin(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	analytics := map[string]interface{}{
		"total_signups": 150,
		"tier_breakdown": map[string]int{
			"starter":      50,
			"professional": 75,
			"enterprise":   25,
		},
		"revenue_projection": map[string]interface{}{
			"monthly": 12500,  // $125/month from beta users
			"annual":  150000, // $150k annual projection
		},
		"conversion_rate": 0.24, // 24% signup to paid conversion
		"churn_rate":      0.02, // 2% monthly churn
		"signup_trend": []map[string]interface{}{
			{"date": "2024-01-01", "signups": 10},
			{"date": "2024-01-02", "signups": 15},
			{"date": "2024-01-03", "signups": 20},
			// Add more trend data
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analytics)
}

// sendBetaWelcomeEmail sends welcome email to beta users
func (bh *BetaHandlers) sendBetaWelcomeEmail(email, name, tier string) {
	// Implement email service integration
	// For now, just log the action
	fmt.Printf("Sending beta welcome email to %s (%s) for tier %s\n", email, name, tier)

	// Email content would include:
	// - Welcome message
	// - Beta access credentials
	// - Discord invite link
	// - Onboarding call scheduling link
	// - Getting started guide
}

// logBetaSignup logs beta signup for analytics
func (bh *BetaHandlers) logBetaSignup(userID, email, tier string, priceAmount int64) {
	// Implement analytics logging
	fmt.Printf("Beta signup: UserID=%s, Email=%s, Tier=%s, Price=$%.2f\n",
		userID, email, tier, float64(priceAmount)/100)

	// This would typically:
	// - Log to analytics database
	// - Send to analytics service (Mixpanel, Amplitude, etc.)
	// - Update metrics dashboard
	// - Trigger marketing automation
}

// BetaMetrics represents beta program metrics
type BetaMetrics struct {
	TotalSignups     int                      `json:"total_signups"`
	TierBreakdown    map[string]int           `json:"tier_breakdown"`
	RevenueProjected float64                  `json:"revenue_projected"`
	ConversionRate   float64                  `json:"conversion_rate"`
	ChurnRate        float64                  `json:"churn_rate"`
	SignupTrend      []map[string]interface{} `json:"signup_trend"`
}

// GetBetaMetrics returns detailed beta metrics
func (bh *BetaHandlers) GetBetaMetrics(ctx context.Context) (*BetaMetrics, error) {
	// In real implementation, this would query the database
	return &BetaMetrics{
		TotalSignups: 150,
		TierBreakdown: map[string]int{
			"starter":      50,
			"professional": 75,
			"enterprise":   25,
		},
		RevenueProjected: 150000, // $150k annual
		ConversionRate:   0.24,   // 24%
		ChurnRate:        0.02,   // 2%
		SignupTrend: []map[string]interface{}{
			{"date": "2024-01-01", "signups": 10},
			{"date": "2024-01-02", "signups": 15},
			{"date": "2024-01-03", "signups": 20},
		},
	}, nil
}
