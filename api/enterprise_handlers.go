package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ai-agentic-browser/internal/enterprise"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

// EnterpriseHandlers handles enterprise sales requests
type EnterpriseHandlers struct {
	salesPipeline *enterprise.SalesPipeline
}

// NewEnterpriseHandlers creates new enterprise handlers
func NewEnterpriseHandlers(salesPipeline *enterprise.SalesPipeline) *EnterpriseHandlers {
	return &EnterpriseHandlers{
		salesPipeline: salesPipeline,
	}
}

// RegisterRoutes registers enterprise sales routes
func (eh *EnterpriseHandlers) RegisterRoutes(router *mux.Router) {
	// Public enterprise routes
	router.HandleFunc("/enterprise/contact", eh.CreateEnterpriseInquiry).Methods("POST")
	router.HandleFunc("/enterprise/pricing", eh.GetEnterprisePricing).Methods("GET")
	router.HandleFunc("/enterprise/demo", eh.RequestDemo).Methods("POST")

	// Sales team routes (authenticated)
	router.HandleFunc("/sales/leads", eh.GetLeads).Methods("GET")
	router.HandleFunc("/sales/leads", eh.CreateLead).Methods("POST")
	router.HandleFunc("/sales/leads/{id}", eh.GetLead).Methods("GET")
	router.HandleFunc("/sales/leads/{id}", eh.UpdateLead).Methods("PUT")
	router.HandleFunc("/sales/leads/{id}/qualify", eh.QualifyLead).Methods("POST")

	router.HandleFunc("/sales/deals", eh.GetDeals).Methods("GET")
	router.HandleFunc("/sales/deals", eh.CreateDeal).Methods("POST")
	router.HandleFunc("/sales/deals/{id}", eh.GetDeal).Methods("GET")
	router.HandleFunc("/sales/deals/{id}", eh.UpdateDeal).Methods("PUT")
	router.HandleFunc("/sales/deals/{id}/close", eh.CloseDeal).Methods("POST")

	router.HandleFunc("/sales/activities", eh.GetActivities).Methods("GET")
	router.HandleFunc("/sales/activities", eh.LogActivity).Methods("POST")

	router.HandleFunc("/sales/proposals", eh.GetProposals).Methods("GET")
	router.HandleFunc("/sales/proposals", eh.CreateProposal).Methods("POST")
	router.HandleFunc("/sales/proposals/{id}", eh.GetProposal).Methods("GET")
	router.HandleFunc("/sales/proposals/{id}/send", eh.SendProposal).Methods("POST")

	// Analytics and reporting
	router.HandleFunc("/sales/pipeline", eh.GetPipelineMetrics).Methods("GET")
	router.HandleFunc("/sales/dashboard", eh.GetSalesDashboard).Methods("GET")
	router.HandleFunc("/sales/forecast", eh.GetSalesForecast).Methods("GET")
}

// EnterpriseInquiryRequest represents an enterprise contact form
type EnterpriseInquiryRequest struct {
	CompanyName      string   `json:"company_name"`
	ContactName      string   `json:"contact_name"`
	ContactEmail     string   `json:"contact_email"`
	ContactPhone     string   `json:"contact_phone"`
	ContactTitle     string   `json:"contact_title"`
	CompanyType      string   `json:"company_type"`
	CompanySize      string   `json:"company_size"`
	AUM              float64  `json:"aum"`
	TradingVolume    float64  `json:"trading_volume"`
	CurrentSolutions []string `json:"current_solutions"`
	PainPoints       []string `json:"pain_points"`
	Budget           float64  `json:"budget"`
	Timeline         string   `json:"timeline"`
	Message          string   `json:"message"`
}

// CreateEnterpriseInquiry handles enterprise contact form submissions
func (eh *EnterpriseHandlers) CreateEnterpriseInquiry(w http.ResponseWriter, r *http.Request) {
	var req EnterpriseInquiryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.CompanyName == "" || req.ContactName == "" || req.ContactEmail == "" {
		http.Error(w, "Company name, contact name, and email are required", http.StatusBadRequest)
		return
	}

	// Create lead from inquiry
	lead := &enterprise.Lead{
		ID:               uuid.New().String(),
		CompanyName:      req.CompanyName,
		ContactName:      req.ContactName,
		ContactEmail:     req.ContactEmail,
		ContactPhone:     req.ContactPhone,
		ContactTitle:     req.ContactTitle,
		CompanySize:      req.CompanySize,
		CompanyType:      req.CompanyType,
		AUM:              decimal.NewFromFloat(req.AUM),
		TradingVolume:    decimal.NewFromFloat(req.TradingVolume),
		CurrentSolutions: req.CurrentSolutions,
		PainPoints:       req.PainPoints,
		Budget:           decimal.NewFromFloat(req.Budget),
		Timeline:         req.Timeline,
		Source:           "website",
		Status:           "new",
		Priority:         eh.calculatePriority(req),
		EstimatedValue:   decimal.NewFromFloat(req.Budget),
		NextFollowUpDate: time.Now().Add(24 * time.Hour),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Auto-assign to sales rep based on company type/size
	lead.AssignedSalesRep = eh.assignSalesRep(req.CompanyType, req.CompanySize)

	err := eh.salesPipeline.CreateLead(r.Context(), lead)
	if err != nil {
		http.Error(w, "Failed to create lead", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"lead_id": lead.ID,
		"message": "Thank you for your interest! Our enterprise team will contact you within 24 hours.",
		"next_steps": []string{
			"Our sales team will review your requirements",
			"We'll schedule a discovery call within 24 hours",
			"Custom demo and proposal will be prepared",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetEnterprisePricing returns enterprise pricing information
func (eh *EnterpriseHandlers) GetEnterprisePricing(w http.ResponseWriter, r *http.Request) {
	pricing := map[string]interface{}{
		"pricing_model": "Custom pricing based on AUM and trading volume",
		"tiers": []map[string]interface{}{
			{
				"name":        "Startup Fund",
				"min_aum":     0,
				"monthly_fee": "$5,000",
				"setup_fee":   "$2,500",
				"description": "For emerging crypto funds and prop trading firms",
				"features": []string{
					"AI trading signals with 85%+ accuracy",
					"Multi-exchange connectivity",
					"Risk management tools",
					"Basic reporting and analytics",
					"Email support",
				},
			},
			{
				"name":        "Growth Fund",
				"min_aum":     10000000,
				"monthly_fee": "$15,000",
				"setup_fee":   "$7,500",
				"description": "For established funds with $10M+ AUM",
				"features": []string{
					"Everything in Startup Fund",
					"Advanced portfolio optimization",
					"Custom trading strategies",
					"Real-time risk monitoring",
					"Phone and chat support",
					"Dedicated account manager",
				},
			},
			{
				"name":        "Institutional",
				"min_aum":     100000000,
				"monthly_fee": "$50,000",
				"setup_fee":   "$25,000",
				"description": "For large institutions with $100M+ AUM",
				"features": []string{
					"Everything in Growth Fund",
					"White-label solutions",
					"On-premise deployment options",
					"Custom integrations",
					"24/7 priority support",
					"SLA guarantees",
					"Compliance reporting",
				},
			},
			{
				"name":        "Enterprise",
				"min_aum":     1000000000,
				"monthly_fee": "$150,000+",
				"setup_fee":   "$75,000+",
				"description": "For major institutions with $1B+ AUM",
				"features": []string{
					"Everything in Institutional",
					"Fully customized solutions",
					"Dedicated infrastructure",
					"Co-location options",
					"Custom AI model training",
					"Regulatory compliance support",
					"Executive support team",
				},
			},
		},
		"additional_fees": map[string]interface{}{
			"performance_fee": "2-20% of profits (negotiable)",
			"api_usage":       "$0.001-$0.01 per request (volume discounts)",
			"data_feeds":      "$1,000-$10,000/month (depends on sources)",
			"training":        "$5,000-$25,000 (one-time)",
		},
		"volume_discounts": []map[string]interface{}{
			{"min_aum": 100000000, "discount": "10%"},
			{"min_aum": 500000000, "discount": "20%"},
			{"min_aum": 1000000000, "discount": "30%"},
		},
		"payment_terms": []string{
			"Net 30 payment terms",
			"Annual prepayment discounts available",
			"Crypto payment options",
			"Escrow arrangements for large contracts",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pricing)
}

// DemoRequest represents a demo request
type DemoRequest struct {
	CompanyName   string  `json:"company_name"`
	ContactName   string  `json:"contact_name"`
	ContactEmail  string  `json:"contact_email"`
	ContactPhone  string  `json:"contact_phone"`
	PreferredTime string  `json:"preferred_time"`
	UseCase       string  `json:"use_case"`
	AUM           float64 `json:"aum"`
}

// RequestDemo handles demo requests
func (eh *EnterpriseHandlers) RequestDemo(w http.ResponseWriter, r *http.Request) {
	var req DemoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create demo activity
	activityID := uuid.New().String()

	response := map[string]interface{}{
		"success": true,
		"demo_id": activityID,
		"message": "Demo request received! We'll contact you to schedule.",
		"next_steps": []string{
			"Our sales engineer will contact you within 4 hours",
			"We'll schedule a 45-minute custom demo",
			"Demo will be tailored to your specific use case",
		},
		"what_to_expect": []string{
			"Live trading demonstration with real market data",
			"AI prediction accuracy showcase",
			"Custom strategy configuration",
			"Q&A with our technical team",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetLeads returns sales leads with filtering
func (eh *EnterpriseHandlers) GetLeads(w http.ResponseWriter, r *http.Request) {
	if !isSalesTeam(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	status := r.URL.Query().Get("status")
	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	leads, err := eh.salesPipeline.GetLeadsByStatus(r.Context(), status, limit)
	if err != nil {
		http.Error(w, "Failed to get leads", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"leads": leads,
		"total": len(leads),
	})
}

// CreateLeadRequest represents lead creation request
type CreateLeadRequest struct {
	CompanyName      string   `json:"company_name"`
	ContactName      string   `json:"contact_name"`
	ContactEmail     string   `json:"contact_email"`
	ContactPhone     string   `json:"contact_phone"`
	ContactTitle     string   `json:"contact_title"`
	CompanySize      string   `json:"company_size"`
	CompanyType      string   `json:"company_type"`
	AUM              float64  `json:"aum"`
	TradingVolume    float64  `json:"trading_volume"`
	CurrentSolutions []string `json:"current_solutions"`
	PainPoints       []string `json:"pain_points"`
	Budget           float64  `json:"budget"`
	Timeline         string   `json:"timeline"`
	Source           string   `json:"source"`
	Priority         string   `json:"priority"`
}

// CreateLead creates a new sales lead
func (eh *EnterpriseHandlers) CreateLead(w http.ResponseWriter, r *http.Request) {
	if !isSalesTeam(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req CreateLeadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	lead := &enterprise.Lead{
		ID:               uuid.New().String(),
		CompanyName:      req.CompanyName,
		ContactName:      req.ContactName,
		ContactEmail:     req.ContactEmail,
		ContactPhone:     req.ContactPhone,
		ContactTitle:     req.ContactTitle,
		CompanySize:      req.CompanySize,
		CompanyType:      req.CompanyType,
		AUM:              decimal.NewFromFloat(req.AUM),
		TradingVolume:    decimal.NewFromFloat(req.TradingVolume),
		CurrentSolutions: req.CurrentSolutions,
		PainPoints:       req.PainPoints,
		Budget:           decimal.NewFromFloat(req.Budget),
		Timeline:         req.Timeline,
		Source:           req.Source,
		Status:           "new",
		Priority:         req.Priority,
		AssignedSalesRep: getUserIDFromContext(r.Context()),
		EstimatedValue:   decimal.NewFromFloat(req.Budget),
		NextFollowUpDate: time.Now().Add(24 * time.Hour),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := eh.salesPipeline.CreateLead(r.Context(), lead)
	if err != nil {
		http.Error(w, "Failed to create lead", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"lead_id": lead.ID,
		"message": "Lead created successfully",
	})
}

// GetPipelineMetrics returns sales pipeline metrics
func (eh *EnterpriseHandlers) GetPipelineMetrics(w http.ResponseWriter, r *http.Request) {
	if !isSalesTeam(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	salesRep := r.URL.Query().Get("sales_rep")
	if salesRep == "" {
		salesRep = getUserIDFromContext(r.Context())
	}

	metrics, err := eh.salesPipeline.GetPipelineMetrics(r.Context(), salesRep)
	if err != nil {
		http.Error(w, "Failed to get metrics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// GetSalesDashboard returns comprehensive sales dashboard data
func (eh *EnterpriseHandlers) GetSalesDashboard(w http.ResponseWriter, r *http.Request) {
	if !isSalesTeam(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Mock dashboard data (implementation would query database)
	dashboard := map[string]interface{}{
		"overview": map[string]interface{}{
			"total_leads":       156,
			"qualified_leads":   89,
			"active_deals":      34,
			"pipeline_value":    2850000.00,
			"closed_this_month": 450000.00,
			"quota_achievement": "78%",
		},
		"pipeline_by_stage": []map[string]interface{}{
			{"stage": "Discovery", "count": 12, "value": 850000.00},
			{"stage": "Demo", "count": 8, "value": 620000.00},
			{"stage": "Proposal", "count": 6, "value": 780000.00},
			{"stage": "Negotiation", "count": 4, "value": 450000.00},
			{"stage": "Contract", "count": 2, "value": 150000.00},
		},
		"top_opportunities": []map[string]interface{}{
			{
				"company":     "Crypto Capital Partners",
				"value":       250000.00,
				"stage":       "Negotiation",
				"probability": 75,
				"close_date":  "2024-02-15",
			},
			{
				"company":     "Digital Asset Fund",
				"value":       180000.00,
				"stage":       "Proposal",
				"probability": 60,
				"close_date":  "2024-02-28",
			},
		},
		"recent_activities": []map[string]interface{}{
			{
				"type":       "Demo",
				"company":    "Blockchain Ventures",
				"date":       "2024-01-28",
				"outcome":    "Positive",
				"next_steps": "Send proposal",
			},
			{
				"type":       "Call",
				"company":    "Crypto Capital Partners",
				"date":       "2024-01-27",
				"outcome":    "Negotiation",
				"next_steps": "Contract review",
			},
		},
		"performance_metrics": map[string]interface{}{
			"conversion_rate":      "18.5%",
			"avg_deal_size":        "$125,000",
			"avg_sales_cycle":      "45 days",
			"activities_this_week": 23,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dashboard)
}

// Helper functions
func (eh *EnterpriseHandlers) calculatePriority(req EnterpriseInquiryRequest) string {
	if req.AUM >= 1000000000 { // $1B+
		return "critical"
	} else if req.AUM >= 100000000 { // $100M+
		return "high"
	} else if req.AUM >= 10000000 { // $10M+
		return "medium"
	}
	return "low"
}

func (eh *EnterpriseHandlers) assignSalesRep(companyType, companySize string) string {
	// Simple assignment logic (would be more sophisticated in real implementation)
	if companyType == "hedge_fund" || companyType == "institution" {
		return "enterprise_sales_rep"
	}
	return "general_sales_rep"
}

func isSalesTeam(ctx context.Context) bool {
	// Implementation would check sales team permissions
	return true
}

// GetLead returns a specific lead
func (eh *EnterpriseHandlers) GetLead(w http.ResponseWriter, r *http.Request) {
	if !isSalesTeam(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	leadID := vars["id"]

	// Mock lead data
	lead := map[string]interface{}{
		"id":           leadID,
		"company_name": "Sample Company",
		"contact_name": "John Doe",
		"status":       "qualified",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lead)
}

// UpdateLead updates a lead
func (eh *EnterpriseHandlers) UpdateLead(w http.ResponseWriter, r *http.Request) {
	if !isSalesTeam(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	leadID := vars["id"]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"lead_id": leadID,
		"message": "Lead updated successfully",
	})
}

// QualifyLead qualifies a lead
func (eh *EnterpriseHandlers) QualifyLead(w http.ResponseWriter, r *http.Request) {
	if !isSalesTeam(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	leadID := vars["id"]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"lead_id": leadID,
		"message": "Lead qualified successfully",
	})
}

// GetDeals returns deals
func (eh *EnterpriseHandlers) GetDeals(w http.ResponseWriter, r *http.Request) {
	if !isSalesTeam(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	deals := []map[string]interface{}{
		{
			"id":           "deal_001",
			"company_name": "Crypto Capital",
			"deal_value":   250000,
			"stage":        "negotiation",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"deals": deals,
		"total": len(deals),
	})
}

// CreateDeal creates a new deal
func (eh *EnterpriseHandlers) CreateDeal(w http.ResponseWriter, r *http.Request) {
	if !isSalesTeam(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	dealID := uuid.New().String()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"deal_id": dealID,
		"message": "Deal created successfully",
	})
}

// GetDeal returns a specific deal
func (eh *EnterpriseHandlers) GetDeal(w http.ResponseWriter, r *http.Request) {
	if !isSalesTeam(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	dealID := vars["id"]

	deal := map[string]interface{}{
		"id":           dealID,
		"company_name": "Sample Company",
		"deal_value":   150000,
		"stage":        "proposal",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deal)
}

// UpdateDeal updates a deal
func (eh *EnterpriseHandlers) UpdateDeal(w http.ResponseWriter, r *http.Request) {
	if !isSalesTeam(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	dealID := vars["id"]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"deal_id": dealID,
		"message": "Deal updated successfully",
	})
}

// CloseDeal closes a deal
func (eh *EnterpriseHandlers) CloseDeal(w http.ResponseWriter, r *http.Request) {
	if !isSalesTeam(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	dealID := vars["id"]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"deal_id": dealID,
		"message": "Deal closed successfully",
	})
}

// GetActivities returns sales activities
func (eh *EnterpriseHandlers) GetActivities(w http.ResponseWriter, r *http.Request) {
	if !isSalesTeam(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	activities := []map[string]interface{}{
		{
			"id":      "activity_001",
			"type":    "demo",
			"company": "Crypto Capital",
			"date":    "2024-01-28",
			"outcome": "positive",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"activities": activities,
		"total":      len(activities),
	})
}

// LogActivity logs a sales activity
func (eh *EnterpriseHandlers) LogActivity(w http.ResponseWriter, r *http.Request) {
	if !isSalesTeam(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	activityID := uuid.New().String()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"activity_id": activityID,
		"message":     "Activity logged successfully",
	})
}

// GetProposals returns proposals
func (eh *EnterpriseHandlers) GetProposals(w http.ResponseWriter, r *http.Request) {
	if !isSalesTeam(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	proposals := []map[string]interface{}{
		{
			"id":          "proposal_001",
			"deal_name":   "Crypto Capital Deal",
			"total_value": 250000,
			"status":      "sent",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"proposals": proposals,
		"total":     len(proposals),
	})
}

// CreateProposal creates a new proposal
func (eh *EnterpriseHandlers) CreateProposal(w http.ResponseWriter, r *http.Request) {
	if !isSalesTeam(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	proposalID := uuid.New().String()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"proposal_id": proposalID,
		"message":     "Proposal created successfully",
	})
}

// GetProposal returns a specific proposal
func (eh *EnterpriseHandlers) GetProposal(w http.ResponseWriter, r *http.Request) {
	if !isSalesTeam(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	proposalID := vars["id"]

	proposal := map[string]interface{}{
		"id":          proposalID,
		"deal_name":   "Sample Deal",
		"total_value": 150000,
		"status":      "draft",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(proposal)
}

// SendProposal sends a proposal
func (eh *EnterpriseHandlers) SendProposal(w http.ResponseWriter, r *http.Request) {
	if !isSalesTeam(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	proposalID := vars["id"]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"proposal_id": proposalID,
		"message":     "Proposal sent successfully",
	})
}

// GetSalesForecast returns sales forecast
func (eh *EnterpriseHandlers) GetSalesForecast(w http.ResponseWriter, r *http.Request) {
	if !isSalesTeam(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	forecast := map[string]interface{}{
		"current_quarter": map[string]interface{}{
			"target":      2000000,
			"actual":      1250000,
			"forecast":    1800000,
			"achievement": "90%",
		},
		"next_quarter": map[string]interface{}{
			"target":   2500000,
			"forecast": 2200000,
			"pipeline": 3500000,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(forecast)
}
