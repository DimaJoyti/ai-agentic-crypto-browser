package enterprise

import (
	"context"
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

// SalesPipeline manages enterprise sales processes
type SalesPipeline struct {
	db *sql.DB
}

// NewSalesPipeline creates a new sales pipeline manager
func NewSalesPipeline(db *sql.DB) *SalesPipeline {
	return &SalesPipeline{
		db: db,
	}
}

// Lead represents an enterprise sales lead
type Lead struct {
	ID               string                 `json:"id"`
	CompanyName      string                 `json:"company_name"`
	ContactName      string                 `json:"contact_name"`
	ContactEmail     string                 `json:"contact_email"`
	ContactPhone     string                 `json:"contact_phone"`
	ContactTitle     string                 `json:"contact_title"`
	CompanySize      string                 `json:"company_size"`      // startup, small, medium, large, enterprise
	CompanyType      string                 `json:"company_type"`      // hedge_fund, family_office, prop_trading, crypto_fund, institution
	AUM              decimal.Decimal        `json:"aum"`               // Assets Under Management
	TradingVolume    decimal.Decimal        `json:"trading_volume"`    // Monthly trading volume
	CurrentSolutions []string               `json:"current_solutions"` // Existing trading platforms
	PainPoints       []string               `json:"pain_points"`       // Key challenges
	Budget           decimal.Decimal        `json:"budget"`            // Annual budget
	Timeline         string                 `json:"timeline"`          // Implementation timeline
	DecisionMakers   []string               `json:"decision_makers"`   // Key stakeholders
	Source           string                 `json:"source"`            // lead_source
	Status           string                 `json:"status"`            // new, qualified, proposal, negotiation, closed_won, closed_lost
	Priority         string                 `json:"priority"`          // low, medium, high, critical
	AssignedSalesRep string                 `json:"assigned_sales_rep"`
	EstimatedValue   decimal.Decimal        `json:"estimated_value"`   // Deal size
	ProbabilityScore decimal.Decimal        `json:"probability_score"` // 0-100%
	LastContactDate  time.Time              `json:"last_contact_date"`
	NextFollowUpDate time.Time              `json:"next_follow_up_date"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	Metadata         map[string]interface{} `json:"metadata"`
}

// Deal represents an enterprise deal in progress
type Deal struct {
	ID                string                 `json:"id"`
	LeadID            string                 `json:"lead_id"`
	DealName          string                 `json:"deal_name"`
	CompanyName       string                 `json:"company_name"`
	DealValue         decimal.Decimal        `json:"deal_value"`
	ContractLength    int                    `json:"contract_length"` // months
	Stage             string                 `json:"stage"`           // discovery, demo, proposal, negotiation, contract, closed
	Probability       decimal.Decimal        `json:"probability"`     // 0-100%
	ExpectedCloseDate time.Time              `json:"expected_close_date"`
	ActualCloseDate   *time.Time             `json:"actual_close_date"`
	SalesRep          string                 `json:"sales_rep"`
	SalesEngineer     string                 `json:"sales_engineer"`
	Products          []string               `json:"products"` // List of products/services
	CustomPricing     bool                   `json:"custom_pricing"`
	WhiteLabel        bool                   `json:"white_label"`
	OnPremise         bool                   `json:"on_premise"`
	SLA               string                 `json:"sla"`         // Service level agreement
	Support           string                 `json:"support"`     // Support tier
	LostReason        string                 `json:"lost_reason"` // If deal is lost
	CompetitorInfo    string                 `json:"competitor_info"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	Metadata          map[string]interface{} `json:"metadata"`
}

// Activity represents sales activities and touchpoints
type Activity struct {
	ID           string                 `json:"id"`
	LeadID       string                 `json:"lead_id"`
	DealID       string                 `json:"deal_id"`
	ActivityType string                 `json:"activity_type"` // call, email, meeting, demo, proposal, contract
	Subject      string                 `json:"subject"`
	Description  string                 `json:"description"`
	Outcome      string                 `json:"outcome"`
	NextSteps    string                 `json:"next_steps"`
	SalesRep     string                 `json:"sales_rep"`
	ScheduledAt  time.Time              `json:"scheduled_at"`
	CompletedAt  *time.Time             `json:"completed_at"`
	Duration     int                    `json:"duration"` // minutes
	Attendees    []string               `json:"attendees"`
	CreatedAt    time.Time              `json:"created_at"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// Proposal represents a custom proposal for enterprise clients
type Proposal struct {
	ID             string            `json:"id"`
	DealID         string            `json:"deal_id"`
	ProposalName   string            `json:"proposal_name"`
	Version        int               `json:"version"`
	TotalValue     decimal.Decimal   `json:"total_value"`
	ContractLength int               `json:"contract_length"` // months
	Products       []ProposalProduct `json:"products"`
	CustomFeatures []string          `json:"custom_features"`
	SLA            ProposalSLA       `json:"sla"`
	Pricing        ProposalPricing   `json:"pricing"`
	Terms          ProposalTerms     `json:"terms"`
	Status         string            `json:"status"` // draft, sent, viewed, accepted, rejected
	SentAt         *time.Time        `json:"sent_at"`
	ViewedAt       *time.Time        `json:"viewed_at"`
	RespondedAt    *time.Time        `json:"responded_at"`
	ExpiresAt      time.Time         `json:"expires_at"`
	CreatedBy      string            `json:"created_by"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

// ProposalProduct represents a product in a proposal
type ProposalProduct struct {
	ProductID    string                 `json:"product_id"`
	ProductName  string                 `json:"product_name"`
	Quantity     int                    `json:"quantity"`
	UnitPrice    decimal.Decimal        `json:"unit_price"`
	TotalPrice   decimal.Decimal        `json:"total_price"`
	Discount     decimal.Decimal        `json:"discount"`
	Description  string                 `json:"description"`
	CustomConfig map[string]interface{} `json:"custom_config"`
}

// ProposalSLA represents service level agreement terms
type ProposalSLA struct {
	Uptime           decimal.Decimal `json:"uptime"`          // 99.9%
	ResponseTime     int             `json:"response_time"`   // minutes
	ResolutionTime   int             `json:"resolution_time"` // hours
	SupportHours     string          `json:"support_hours"`   // 24/7, business_hours
	DedicatedSupport bool            `json:"dedicated_support"`
	AccountManager   bool            `json:"account_manager"`
}

// ProposalPricing represents pricing structure
type ProposalPricing struct {
	Model           string           `json:"model"` // subscription, usage_based, hybrid
	SetupFee        decimal.Decimal  `json:"setup_fee"`
	MonthlyFee      decimal.Decimal  `json:"monthly_fee"`
	UsageFee        decimal.Decimal  `json:"usage_fee"`
	PerformanceFee  decimal.Decimal  `json:"performance_fee"`
	VolumeDiscounts []VolumeDiscount `json:"volume_discounts"`
	PaymentTerms    string           `json:"payment_terms"` // net_30, net_60, annual_prepay
}

// VolumeDiscount represents volume-based pricing tiers
type VolumeDiscount struct {
	MinVolume decimal.Decimal `json:"min_volume"`
	MaxVolume decimal.Decimal `json:"max_volume"`
	Discount  decimal.Decimal `json:"discount"`
}

// ProposalTerms represents contract terms
type ProposalTerms struct {
	ContractLength   int      `json:"contract_length"` // months
	AutoRenewal      bool     `json:"auto_renewal"`
	TerminationTerms string   `json:"termination_terms"`
	DataRetention    int      `json:"data_retention"` // days
	Compliance       []string `json:"compliance"`     // SOC2, GDPR, etc.
	Liability        string   `json:"liability"`
	Warranty         string   `json:"warranty"`
}

// CreateLead creates a new enterprise lead
func (sp *SalesPipeline) CreateLead(ctx context.Context, lead *Lead) error {
	query := `
		INSERT INTO enterprise_leads (
			id, company_name, contact_name, contact_email, contact_phone, contact_title,
			company_size, company_type, aum, trading_volume, budget, timeline,
			source, status, priority, assigned_sales_rep, estimated_value,
			probability_score, next_follow_up_date, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
	`

	_, err := sp.db.ExecContext(ctx, query,
		lead.ID, lead.CompanyName, lead.ContactName, lead.ContactEmail, lead.ContactPhone,
		lead.ContactTitle, lead.CompanySize, lead.CompanyType, lead.AUM, lead.TradingVolume,
		lead.Budget, lead.Timeline, lead.Source, lead.Status, lead.Priority,
		lead.AssignedSalesRep, lead.EstimatedValue, lead.ProbabilityScore,
		lead.NextFollowUpDate, lead.CreatedAt, lead.UpdatedAt,
	)

	return err
}

// QualifyLead updates lead status and qualification information
func (sp *SalesPipeline) QualifyLead(ctx context.Context, leadID string, qualification map[string]interface{}) error {
	query := `
		UPDATE enterprise_leads 
		SET status = 'qualified', 
		    probability_score = $1,
		    estimated_value = $2,
		    updated_at = $3
		WHERE id = $4
	`

	probabilityScore := qualification["probability_score"].(decimal.Decimal)
	estimatedValue := qualification["estimated_value"].(decimal.Decimal)

	_, err := sp.db.ExecContext(ctx, query, probabilityScore, estimatedValue, time.Now(), leadID)
	return err
}

// CreateDeal creates a new deal from a qualified lead
func (sp *SalesPipeline) CreateDeal(ctx context.Context, deal *Deal) error {
	query := `
		INSERT INTO enterprise_deals (
			id, lead_id, deal_name, company_name, deal_value, contract_length,
			stage, probability, expected_close_date, sales_rep, sales_engineer,
			custom_pricing, white_label, on_premise, sla, support,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`

	_, err := sp.db.ExecContext(ctx, query,
		deal.ID, deal.LeadID, deal.DealName, deal.CompanyName, deal.DealValue,
		deal.ContractLength, deal.Stage, deal.Probability, deal.ExpectedCloseDate,
		deal.SalesRep, deal.SalesEngineer, deal.CustomPricing, deal.WhiteLabel,
		deal.OnPremise, deal.SLA, deal.Support, deal.CreatedAt, deal.UpdatedAt,
	)

	return err
}

// LogActivity records a sales activity
func (sp *SalesPipeline) LogActivity(ctx context.Context, activity *Activity) error {
	query := `
		INSERT INTO sales_activities (
			id, lead_id, deal_id, activity_type, subject, description,
			outcome, next_steps, sales_rep, scheduled_at, completed_at,
			duration, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := sp.db.ExecContext(ctx, query,
		activity.ID, activity.LeadID, activity.DealID, activity.ActivityType,
		activity.Subject, activity.Description, activity.Outcome, activity.NextSteps,
		activity.SalesRep, activity.ScheduledAt, activity.CompletedAt,
		activity.Duration, activity.CreatedAt,
	)

	return err
}

// GetPipelineMetrics returns sales pipeline metrics
func (sp *SalesPipeline) GetPipelineMetrics(ctx context.Context, salesRep string) (*PipelineMetrics, error) {
	query := `
		SELECT 
			COUNT(*) as total_leads,
			COUNT(CASE WHEN status = 'qualified' THEN 1 END) as qualified_leads,
			COUNT(CASE WHEN status = 'proposal' THEN 1 END) as proposal_stage,
			COUNT(CASE WHEN status = 'negotiation' THEN 1 END) as negotiation_stage,
			COALESCE(SUM(estimated_value), 0) as total_pipeline_value,
			COALESCE(AVG(probability_score), 0) as avg_probability
		FROM enterprise_leads 
		WHERE assigned_sales_rep = $1 OR $1 = ''
	`

	metrics := &PipelineMetrics{}
	err := sp.db.QueryRowContext(ctx, query, salesRep).Scan(
		&metrics.TotalLeads, &metrics.QualifiedLeads, &metrics.ProposalStage,
		&metrics.NegotiationStage, &metrics.TotalPipelineValue, &metrics.AvgProbability,
	)

	return metrics, err
}

// PipelineMetrics represents sales pipeline metrics
type PipelineMetrics struct {
	TotalLeads         int64           `json:"total_leads"`
	QualifiedLeads     int64           `json:"qualified_leads"`
	ProposalStage      int64           `json:"proposal_stage"`
	NegotiationStage   int64           `json:"negotiation_stage"`
	TotalPipelineValue decimal.Decimal `json:"total_pipeline_value"`
	AvgProbability     decimal.Decimal `json:"avg_probability"`
	ConversionRate     decimal.Decimal `json:"conversion_rate"`
	AvgDealSize        decimal.Decimal `json:"avg_deal_size"`
	AvgSalesCycle      int             `json:"avg_sales_cycle"` // days
}

// GetLeadsByStatus returns leads filtered by status
func (sp *SalesPipeline) GetLeadsByStatus(ctx context.Context, status string, limit int) ([]*Lead, error) {
	query := `
		SELECT id, company_name, contact_name, contact_email, company_type,
		       aum, estimated_value, status, priority, assigned_sales_rep,
		       next_follow_up_date, created_at, updated_at
		FROM enterprise_leads 
		WHERE status = $1 OR $1 = ''
		ORDER BY priority DESC, next_follow_up_date ASC
		LIMIT $2
	`

	rows, err := sp.db.QueryContext(ctx, query, status, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leads []*Lead
	for rows.Next() {
		lead := &Lead{}
		err := rows.Scan(
			&lead.ID, &lead.CompanyName, &lead.ContactName, &lead.ContactEmail,
			&lead.CompanyType, &lead.AUM, &lead.EstimatedValue, &lead.Status,
			&lead.Priority, &lead.AssignedSalesRep, &lead.NextFollowUpDate,
			&lead.CreatedAt, &lead.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		leads = append(leads, lead)
	}

	return leads, nil
}

// UpdateDealStage updates deal stage and probability
func (sp *SalesPipeline) UpdateDealStage(ctx context.Context, dealID, stage string, probability decimal.Decimal) error {
	query := `
		UPDATE enterprise_deals 
		SET stage = $1, probability = $2, updated_at = $3
		WHERE id = $4
	`

	_, err := sp.db.ExecContext(ctx, query, stage, probability, time.Now(), dealID)
	return err
}

// CloseDeal marks a deal as won or lost
func (sp *SalesPipeline) CloseDeal(ctx context.Context, dealID string, won bool, actualValue decimal.Decimal, lostReason string) error {
	status := "closed_lost"
	if won {
		status = "closed_won"
	}

	query := `
		UPDATE enterprise_deals 
		SET stage = $1, actual_close_date = $2, deal_value = $3, lost_reason = $4, updated_at = $5
		WHERE id = $6
	`

	_, err := sp.db.ExecContext(ctx, query, status, time.Now(), actualValue, lostReason, time.Now(), dealID)
	return err
}
