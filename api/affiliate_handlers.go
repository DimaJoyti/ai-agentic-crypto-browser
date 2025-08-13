package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ai-agentic-browser/internal/affiliate"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// AffiliateHandlers handles affiliate program requests
type AffiliateHandlers struct {
	affiliateTracker *affiliate.AffiliateTracker
}

// NewAffiliateHandlers creates new affiliate handlers
func NewAffiliateHandlers(affiliateTracker *affiliate.AffiliateTracker) *AffiliateHandlers {
	return &AffiliateHandlers{
		affiliateTracker: affiliateTracker,
	}
}

// RegisterRoutes registers affiliate routes
func (ah *AffiliateHandlers) RegisterRoutes(router *mux.Router) {
	// Public affiliate routes
	router.HandleFunc("/affiliate/program", ah.GetProgramInfo).Methods("GET")
	router.HandleFunc("/affiliate/apply", ah.ApplyForProgram).Methods("POST")
	router.HandleFunc("/affiliate/track/{code}", ah.TrackClick).Methods("GET")

	// Authenticated affiliate routes
	router.HandleFunc("/affiliate/dashboard", ah.GetAffiliateDashboard).Methods("GET")
	router.HandleFunc("/affiliate/stats", ah.GetAffiliateStats).Methods("GET")
	router.HandleFunc("/affiliate/referrals", ah.GetReferrals).Methods("GET")
	router.HandleFunc("/affiliate/payouts", ah.GetPayouts).Methods("GET")
	router.HandleFunc("/affiliate/links", ah.GenerateAffiliateLinks).Methods("POST")

	// Admin routes
	router.HandleFunc("/admin/affiliates", ah.ListAffiliates).Methods("GET")
	router.HandleFunc("/admin/affiliates/{id}/approve", ah.ApproveAffiliate).Methods("POST")
	router.HandleFunc("/admin/affiliates/{id}/suspend", ah.SuspendAffiliate).Methods("POST")
	router.HandleFunc("/admin/affiliates/payouts", ah.ProcessPayouts).Methods("POST")
	router.HandleFunc("/admin/affiliates/leaderboard", ah.GetLeaderboard).Methods("GET")
}

// GetProgramInfo returns affiliate program information
func (ah *AffiliateHandlers) GetProgramInfo(w http.ResponseWriter, r *http.Request) {
	programInfo := map[string]interface{}{
		"program_name": "AI-Agentic Crypto Browser Affiliate Program",
		"description":  "Earn up to 35% commission by referring traders to our AI-powered platform",
		"commission_structure": map[string]interface{}{
			"bronze":   map[string]interface{}{"rate": "15%", "min_referrals": 0},
			"silver":   map[string]interface{}{"rate": "20%", "min_referrals": 10},
			"gold":     map[string]interface{}{"rate": "25%", "min_referrals": 50},
			"platinum": map[string]interface{}{"rate": "30%", "min_referrals": 100},
			"diamond":  map[string]interface{}{"rate": "35%", "min_referrals": 500},
		},
		"payment_schedule": "Monthly payouts on the 1st of each month",
		"minimum_payout":   "$100",
		"payment_methods":  []string{"Stripe", "Crypto", "Bank Transfer", "PayPal"},
		"cookie_duration":  "90 days",
		"benefits": []string{
			"High conversion rates with 85%+ AI accuracy",
			"Recurring commissions on subscriptions",
			"Performance fee commissions (2-20% of profits)",
			"API usage commissions",
			"Real-time tracking and analytics",
			"Dedicated affiliate support",
			"Marketing materials and resources",
		},
		"requirements": []string{
			"Active in crypto/trading community",
			"Quality traffic sources",
			"Compliance with promotional guidelines",
			"Minimum 18 years old",
		},
		"application_process": map[string]interface{}{
			"steps": []string{
				"Submit application with details",
				"Review by affiliate team (1-3 business days)",
				"Approval and affiliate code generation",
				"Access to affiliate dashboard and materials",
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(programInfo)
}

// AffiliateApplicationRequest represents affiliate application
type AffiliateApplicationRequest struct {
	BusinessName        string            `json:"business_name"`
	WebsiteURL          string            `json:"website_url"`
	SocialMediaLinks    map[string]string `json:"social_media_links"`
	MarketingExperience string            `json:"marketing_experience"`
	TargetAudience      string            `json:"target_audience"`
	PromotionalMethods  string            `json:"promotional_methods"`
	ExpectedReferrals   int               `json:"expected_referrals"`
}

// ApplyForProgram handles affiliate program applications
func (ah *AffiliateHandlers) ApplyForProgram(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req AffiliateApplicationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.MarketingExperience == "" || req.TargetAudience == "" {
		http.Error(w, "Marketing experience and target audience are required", http.StatusBadRequest)
		return
	}

	// Create affiliate application (implementation would store in database)
	applicationID := uuid.New().String()

	response := map[string]interface{}{
		"success":        true,
		"application_id": applicationID,
		"message":        "Application submitted successfully",
		"next_steps": []string{
			"Your application will be reviewed within 1-3 business days",
			"You'll receive an email notification with the decision",
			"If approved, you'll get access to your affiliate dashboard",
		},
		"estimated_review_time": "1-3 business days",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// TrackClick tracks affiliate link clicks
func (ah *AffiliateHandlers) TrackClick(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	affiliateCode := vars["code"]

	if affiliateCode == "" {
		http.Error(w, "Invalid affiliate code", http.StatusBadRequest)
		return
	}

	// Get affiliate by code
	affiliate, err := ah.affiliateTracker.GetAffiliateByCode(r.Context(), affiliateCode)
	if err != nil {
		http.Error(w, "Invalid affiliate code", http.StatusNotFound)
		return
	}

	// Track click (implementation would store click data)
	clickData := map[string]interface{}{
		"affiliate_id":   affiliate.ID,
		"affiliate_code": affiliateCode,
		"ip_address":     r.RemoteAddr,
		"user_agent":     r.UserAgent(),
		"referer":        r.Referer(),
		"timestamp":      time.Now(),
	}

	// Redirect to main site with tracking parameters
	redirectURL := fmt.Sprintf("/?ref=%s&utm_source=affiliate&utm_medium=referral&utm_campaign=%s",
		affiliateCode, affiliate.AffiliateType)

	// Set tracking cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "affiliate_ref",
		Value:    affiliateCode,
		Expires:  time.Now().Add(90 * 24 * time.Hour), // 90 days
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	// Log click for analytics
	fmt.Printf("Affiliate click tracked: %+v\n", clickData)

	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

// GetAffiliateDashboard returns affiliate dashboard data
func (ah *AffiliateHandlers) GetAffiliateDashboard(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Mock affiliate data (implementation would query database)
	dashboard := map[string]interface{}{
		"affiliate_info": map[string]interface{}{
			"affiliate_code":  "CRYPTO_TRADER_123",
			"status":          "active",
			"tier":            "Silver",
			"commission_rate": "20%",
			"member_since":    "2024-01-15",
		},
		"performance_summary": map[string]interface{}{
			"total_clicks":        1250,
			"total_signups":       89,
			"total_conversions":   34,
			"conversion_rate":     "2.72%",
			"total_commissions":   2850.00,
			"unpaid_commissions":  450.00,
			"this_month_earnings": 450.00,
			"last_payout":         "2024-01-01",
			"next_payout":         "2024-02-01",
		},
		"recent_referrals": []map[string]interface{}{
			{
				"user_id":          "user_001",
				"conversion_type":  "subscription",
				"conversion_value": 199.00,
				"commission":       39.80,
				"status":           "confirmed",
				"date":             "2024-01-28",
			},
			{
				"user_id":          "user_002",
				"conversion_type":  "api_usage",
				"conversion_value": 50.00,
				"commission":       10.00,
				"status":           "pending",
				"date":             "2024-01-27",
			},
		},
		"marketing_materials": []map[string]interface{}{
			{
				"type":        "banner",
				"title":       "AI Trading Platform - 728x90",
				"url":         "/assets/banners/ai-trading-728x90.png",
				"description": "High-converting banner for websites",
			},
			{
				"type":        "text_link",
				"title":       "Try AI Trading with 85% Accuracy",
				"url":         "https://ai-crypto-browser.com/ref/CRYPTO_TRADER_123",
				"description": "Simple text link for social media",
			},
		},
		"affiliate_links": map[string]interface{}{
			"main_site":   "https://ai-crypto-browser.com/ref/CRYPTO_TRADER_123",
			"beta_signup": "https://ai-crypto-browser.com/beta-signup?ref=CRYPTO_TRADER_123",
			"api_docs":    "https://ai-crypto-browser.com/api?ref=CRYPTO_TRADER_123",
			"pricing":     "https://ai-crypto-browser.com/pricing?ref=CRYPTO_TRADER_123",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dashboard)
}

// GetAffiliateStats returns detailed affiliate statistics
func (ah *AffiliateHandlers) GetAffiliateStats(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "monthly"
	}

	// Mock stats data
	stats := map[string]interface{}{
		"period": period,
		"metrics": map[string]interface{}{
			"clicks":              1250,
			"signups":             89,
			"conversions":         34,
			"conversion_rate":     2.72,
			"total_revenue":       6780.00,
			"total_commissions":   1356.00,
			"average_order_value": 199.41,
		},
		"trends": map[string]interface{}{
			"clicks_trend":     "+15%",
			"conversion_trend": "+8%",
			"revenue_trend":    "+22%",
		},
		"top_sources": []map[string]interface{}{
			{"source": "social_media", "clicks": 450, "conversions": 15},
			{"source": "website", "clicks": 380, "conversions": 12},
			{"source": "email", "clicks": 250, "conversions": 5},
			{"source": "direct", "clicks": 170, "conversions": 2},
		},
		"conversion_types": map[string]interface{}{
			"subscription":     28,
			"api_usage":        4,
			"performance_fees": 2,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// GetReferrals returns affiliate's referral history
func (ah *AffiliateHandlers) GetReferrals(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	// Mock referral data
	referrals := []map[string]interface{}{
		{
			"id":               "ref_001",
			"referred_user_id": "user_001",
			"conversion_type":  "subscription",
			"conversion_value": 199.00,
			"commission":       39.80,
			"status":           "confirmed",
			"source":           "social_media",
			"converted_at":     "2024-01-28T10:30:00Z",
			"confirmed_at":     "2024-01-29T09:15:00Z",
		},
		{
			"id":               "ref_002",
			"referred_user_id": "user_002",
			"conversion_type":  "api_usage",
			"conversion_value": 50.00,
			"commission":       10.00,
			"status":           "pending",
			"source":           "website",
			"converted_at":     "2024-01-27T14:20:00Z",
		},
	}

	response := map[string]interface{}{
		"referrals": referrals[:min(len(referrals), limit)],
		"total":     len(referrals),
		"summary": map[string]interface{}{
			"total_commissions":     49.80,
			"pending_commissions":   10.00,
			"confirmed_commissions": 39.80,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetPayouts returns affiliate payout history
func (ah *AffiliateHandlers) GetPayouts(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	payouts := []map[string]interface{}{
		{
			"id":             "payout_001",
			"period":         "2024-01",
			"amount":         1250.00,
			"fees":           12.50,
			"net_amount":     1237.50,
			"payment_method": "stripe",
			"status":         "completed",
			"processed_at":   "2024-02-01T09:00:00Z",
			"completed_at":   "2024-02-01T09:05:00Z",
		},
		{
			"id":             "payout_002",
			"period":         "2023-12",
			"amount":         890.00,
			"fees":           8.90,
			"net_amount":     881.10,
			"payment_method": "stripe",
			"status":         "completed",
			"processed_at":   "2024-01-01T09:00:00Z",
			"completed_at":   "2024-01-01T09:03:00Z",
		},
	}

	response := map[string]interface{}{
		"payouts":    payouts,
		"total_paid": 2118.60,
		"next_payout": map[string]interface{}{
			"estimated_amount": 450.00,
			"payout_date":      "2024-03-01",
			"status":           "pending",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GenerateAffiliateLinkRequest represents link generation request
type GenerateAffiliateLinkRequest struct {
	Page     string            `json:"page"`
	Campaign string            `json:"campaign"`
	Source   string            `json:"source"`
	Medium   string            `json:"medium"`
	Custom   map[string]string `json:"custom"`
}

// GenerateAffiliateLinks generates custom affiliate links
func (ah *AffiliateHandlers) GenerateAffiliateLinks(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req GenerateAffiliateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Mock affiliate code (would be retrieved from database)
	affiliateCode := "CRYPTO_TRADER_123"

	baseURL := "https://ai-crypto-browser.com"
	if req.Page != "" {
		baseURL += "/" + req.Page
	}

	// Build tracking parameters
	params := fmt.Sprintf("?ref=%s", affiliateCode)
	if req.Campaign != "" {
		params += fmt.Sprintf("&utm_campaign=%s", req.Campaign)
	}
	if req.Source != "" {
		params += fmt.Sprintf("&utm_source=%s", req.Source)
	}
	if req.Medium != "" {
		params += fmt.Sprintf("&utm_medium=%s", req.Medium)
	}

	affiliateLink := baseURL + params

	response := map[string]interface{}{
		"success":        true,
		"affiliate_link": affiliateLink,
		"short_link":     "https://ai-crypto.ly/abc123", // Mock short link
		"qr_code":        "https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=" + affiliateLink,
		"tracking_info": map[string]interface{}{
			"affiliate_code": affiliateCode,
			"campaign":       req.Campaign,
			"source":         req.Source,
			"medium":         req.Medium,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListAffiliates returns list of affiliates (admin only)
func (ah *AffiliateHandlers) ListAffiliates(w http.ResponseWriter, r *http.Request) {
	// Check admin permissions
	if !isAdmin(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	affiliates, err := ah.affiliateTracker.GetTopAffiliates(r.Context(), 100)
	if err != nil {
		http.Error(w, "Failed to get affiliates", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"affiliates": affiliates,
		"total":      len(affiliates),
	})
}

// ApproveAffiliate approves an affiliate application (admin only)
func (ah *AffiliateHandlers) ApproveAffiliate(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	affiliateID := vars["id"]

	// Implementation would update affiliate status to active
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Affiliate %s approved", affiliateID),
	})
}

// SuspendAffiliate suspends an affiliate (admin only)
func (ah *AffiliateHandlers) SuspendAffiliate(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	affiliateID := vars["id"]

	// Implementation would update affiliate status to suspended
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Affiliate %s suspended", affiliateID),
	})
}

// ProcessPayouts processes affiliate payouts (admin only)
func (ah *AffiliateHandlers) ProcessPayouts(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Implementation would process pending payouts
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":         true,
		"message":         "Payouts processed successfully",
		"processed_count": 25,
		"total_amount":    15750.00,
	})
}

// GetLeaderboard returns affiliate leaderboard (admin only)
func (ah *AffiliateHandlers) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	leaderboard := []map[string]interface{}{
		{
			"rank":              1,
			"affiliate_code":    "CRYPTO_MASTER_001",
			"total_referrals":   156,
			"total_commissions": 8750.00,
			"conversion_rate":   "3.2%",
			"tier":              "Diamond",
		},
		{
			"rank":              2,
			"affiliate_code":    "TRADING_PRO_002",
			"total_referrals":   89,
			"total_commissions": 5200.00,
			"conversion_rate":   "2.8%",
			"tier":              "Platinum",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"leaderboard": leaderboard,
		"period":      "all_time",
	})
}
