package funding

import (
	"context"
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

// InvestorRelationsManager manages Series A funding and investor relations
type InvestorRelationsManager struct {
	db                *sql.DB
	pitchDeckManager  *PitchDeckManager
	financialModeling *FinancialModelingEngine
	investorCRM       *InvestorCRM
	dataRoom          *VirtualDataRoom
	complianceTracker *FundingComplianceTracker
}

// NewInvestorRelationsManager creates a new investor relations manager
func NewInvestorRelationsManager(db *sql.DB) *InvestorRelationsManager {
	return &InvestorRelationsManager{
		db:                db,
		pitchDeckManager:  NewPitchDeckManager(),
		financialModeling: NewFinancialModelingEngine(),
		investorCRM:       NewInvestorCRM(),
		dataRoom:          NewVirtualDataRoom(),
		complianceTracker: NewFundingComplianceTracker(),
	}
}

// FundingRound represents a funding round
type FundingRound struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Type            string                 `json:"type"`            // seed, series_a, series_b, series_c
	Status          string                 `json:"status"`          // planning, active, closed, cancelled
	TargetAmount    decimal.Decimal        `json:"target_amount"`
	MinAmount       decimal.Decimal        `json:"min_amount"`
	MaxAmount       decimal.Decimal        `json:"max_amount"`
	RaisedAmount    decimal.Decimal        `json:"raised_amount"`
	Valuation       FundingValuation       `json:"valuation"`
	Terms           FundingTerms           `json:"terms"`
	Timeline        FundingTimeline        `json:"timeline"`
	Investors       []InvestorCommitment   `json:"investors"`
	Documents       []FundingDocument      `json:"documents"`
	Milestones      []FundingMilestone     `json:"milestones"`
	UseOfFunds      UseOfFunds             `json:"use_of_funds"`
	RiskFactors     []RiskFactor           `json:"risk_factors"`
	CompetitiveLandscape CompetitiveLandscape `json:"competitive_landscape"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// FundingValuation represents company valuation
type FundingValuation struct {
	PreMoneyValuation  decimal.Decimal `json:"pre_money_valuation"`
	PostMoneyValuation decimal.Decimal `json:"post_money_valuation"`
	SharePrice         decimal.Decimal `json:"share_price"`
	SharesOutstanding  int64           `json:"shares_outstanding"`
	NewSharesIssued    int64           `json:"new_shares_issued"`
	OptionPool         decimal.Decimal `json:"option_pool"`        // percentage
	FullyDilutedShares int64           `json:"fully_diluted_shares"`
	ValuationMethod    string          `json:"valuation_method"`   // dcf, comparable, revenue_multiple
	Multiples          ValuationMultiples `json:"multiples"`
}

// ValuationMultiples represents valuation multiples
type ValuationMultiples struct {
	RevenueMultiple    decimal.Decimal `json:"revenue_multiple"`
	EBITDAMultiple     decimal.Decimal `json:"ebitda_multiple"`
	UserMultiple       decimal.Decimal `json:"user_multiple"`
	ARRMultiple        decimal.Decimal `json:"arr_multiple"`
	BookValueMultiple  decimal.Decimal `json:"book_value_multiple"`
	ComparableCompanies []ComparableCompany `json:"comparable_companies"`
}

// ComparableCompany represents comparable company data
type ComparableCompany struct {
	Name            string          `json:"name"`
	Ticker          string          `json:"ticker"`
	MarketCap       decimal.Decimal `json:"market_cap"`
	Revenue         decimal.Decimal `json:"revenue"`
	RevenueGrowth   decimal.Decimal `json:"revenue_growth"`
	RevenueMultiple decimal.Decimal `json:"revenue_multiple"`
	Stage           string          `json:"stage"`
	Geography       string          `json:"geography"`
}

// FundingTerms represents funding terms and conditions
type FundingTerms struct {
	SecurityType        string          `json:"security_type"`        // preferred_stock, convertible_note, safe
	LiquidationPreference decimal.Decimal `json:"liquidation_preference"` // 1x, 2x, etc.
	Participation       string          `json:"participation"`        // non_participating, participating, capped
	DividendRate        decimal.Decimal `json:"dividend_rate"`        // annual percentage
	AntiDilution        string          `json:"anti_dilution"`        // weighted_average, full_ratchet, none
	VotingRights        string          `json:"voting_rights"`        // as_converted, separate_class
	BoardSeats          BoardComposition `json:"board_seats"`
	ProtectiveProvisions []string       `json:"protective_provisions"`
	DragAlongRights     bool            `json:"drag_along_rights"`
	TagAlongRights      bool            `json:"tag_along_rights"`
	RedemptionRights    bool            `json:"redemption_rights"`
	ConversionRights    ConversionTerms `json:"conversion_rights"`
}

// BoardComposition represents board composition
type BoardComposition struct {
	TotalSeats      int      `json:"total_seats"`
	CommonSeats     int      `json:"common_seats"`
	PreferredSeats  int      `json:"preferred_seats"`
	IndependentSeats int     `json:"independent_seats"`
	ObserverRights  []string `json:"observer_rights"`
}

// ConversionTerms represents conversion terms
type ConversionTerms struct {
	ConversionRatio    decimal.Decimal `json:"conversion_ratio"`
	ConversionPrice    decimal.Decimal `json:"conversion_price"`
	AutoConversion     bool            `json:"auto_conversion"`
	ConversionTriggers []string        `json:"conversion_triggers"`
}

// FundingTimeline represents funding timeline
type FundingTimeline struct {
	KickoffDate     time.Time `json:"kickoff_date"`
	PitchStartDate  time.Time `json:"pitch_start_date"`
	DueDiligenceDate time.Time `json:"due_diligence_date"`
	TermSheetDate   time.Time `json:"term_sheet_date"`
	ClosingDate     time.Time `json:"closing_date"`
	Milestones      []TimelineMilestone `json:"milestones"`
}

// TimelineMilestone represents timeline milestone
type TimelineMilestone struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	CompletedAt time.Time `json:"completed_at"`
	Status      string    `json:"status"`      // pending, in_progress, completed, delayed
	Owner       string    `json:"owner"`
	Dependencies []string `json:"dependencies"`
}

// InvestorCommitment represents investor commitment
type InvestorCommitment struct {
	InvestorID      string          `json:"investor_id"`
	InvestorName    string          `json:"investor_name"`
	InvestorType    string          `json:"investor_type"`    // vc, angel, strategic, family_office
	CommitmentAmount decimal.Decimal `json:"commitment_amount"`
	CommitmentDate  time.Time       `json:"commitment_date"`
	Status          string          `json:"status"`          // interested, committed, signed, funded
	LeadInvestor    bool            `json:"lead_investor"`
	BoardSeat       bool            `json:"board_seat"`
	ProRataRights   bool            `json:"pro_rata_rights"`
	ReferenceCheck  bool            `json:"reference_check"`
	BackgroundCheck bool            `json:"background_check"`
	Notes           string          `json:"notes"`
}

// FundingDocument represents funding documents
type FundingDocument struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`        // pitch_deck, financial_model, legal_docs, due_diligence
	Category    string    `json:"category"`    // public, confidential, restricted
	Version     string    `json:"version"`
	FileURL     string    `json:"file_url"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	AccessLog   []DocumentAccess `json:"access_log"`
}

// DocumentAccess represents document access log
type DocumentAccess struct {
	UserID      string    `json:"user_id"`
	UserName    string    `json:"user_name"`
	UserType    string    `json:"user_type"`    // investor, advisor, employee
	AccessType  string    `json:"access_type"`  // view, download, share
	AccessedAt  time.Time `json:"accessed_at"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
}

// FundingMilestone represents funding milestones
type FundingMilestone struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`        // revenue, user, product, partnership
	Target      decimal.Decimal `json:"target"`
	Current     decimal.Decimal `json:"current"`
	Unit        string    `json:"unit"`        // dollars, users, percentage
	DueDate     time.Time `json:"due_date"`
	CompletedAt time.Time `json:"completed_at"`
	Status      string    `json:"status"`      // pending, in_progress, completed, missed
	Impact      string    `json:"impact"`      // high, medium, low
}

// UseOfFunds represents use of funds breakdown
type UseOfFunds struct {
	TotalAmount decimal.Decimal    `json:"total_amount"`
	Categories  []FundingCategory  `json:"categories"`
	Timeline    []FundingTimeline  `json:"timeline"`
	Milestones  []FundingMilestone `json:"milestones"`
}

// FundingCategory represents funding category
type FundingCategory struct {
	Name        string          `json:"name"`
	Amount      decimal.Decimal `json:"amount"`
	Percentage  decimal.Decimal `json:"percentage"`
	Description string          `json:"description"`
	Timeline    string          `json:"timeline"`    // immediate, 6_months, 12_months, 18_months
	Subcategories []FundingSubcategory `json:"subcategories"`
}

// FundingSubcategory represents funding subcategory
type FundingSubcategory struct {
	Name        string          `json:"name"`
	Amount      decimal.Decimal `json:"amount"`
	Description string          `json:"description"`
	Justification string        `json:"justification"`
}

// RiskFactor represents risk factors
type RiskFactor struct {
	ID          string `json:"id"`
	Category    string `json:"category"`    // market, technology, regulatory, competitive, financial
	Title       string `json:"title"`
	Description string `json:"description"`
	Impact      string `json:"impact"`      // high, medium, low
	Probability string `json:"probability"` // high, medium, low
	Mitigation  string `json:"mitigation"`
	Status      string `json:"status"`      // active, mitigated, resolved
}

// CompetitiveLandscape represents competitive analysis
type CompetitiveLandscape struct {
	DirectCompetitors   []Competitor `json:"direct_competitors"`
	IndirectCompetitors []Competitor `json:"indirect_competitors"`
	CompetitiveAdvantages []string   `json:"competitive_advantages"`
	MarketPosition      string       `json:"market_position"`
	MarketShare         decimal.Decimal `json:"market_share"`
	CompetitiveMatrix   CompetitiveMatrix `json:"competitive_matrix"`
}

// Competitor represents competitor information
type Competitor struct {
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	Website         string          `json:"website"`
	Funding         decimal.Decimal `json:"funding"`
	Valuation       decimal.Decimal `json:"valuation"`
	Employees       int             `json:"employees"`
	Revenue         decimal.Decimal `json:"revenue"`
	MarketShare     decimal.Decimal `json:"market_share"`
	Strengths       []string        `json:"strengths"`
	Weaknesses      []string        `json:"weaknesses"`
	KeyFeatures     []string        `json:"key_features"`
	PricingModel    string          `json:"pricing_model"`
	TargetMarket    string          `json:"target_market"`
}

// CompetitiveMatrix represents competitive feature matrix
type CompetitiveMatrix struct {
	Features    []string                    `json:"features"`
	Companies   []string                    `json:"companies"`
	Matrix      map[string]map[string]bool  `json:"matrix"` // company -> feature -> has_feature
	Scoring     map[string]decimal.Decimal  `json:"scoring"` // company -> overall_score
}

// Investor represents investor information
type Investor struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Type            string                 `json:"type"`            // vc_fund, angel, strategic, family_office
	Tier            string                 `json:"tier"`            // tier_1, tier_2, tier_3
	FocusStage      []string               `json:"focus_stage"`     // seed, series_a, series_b, growth
	FocusSectors    []string               `json:"focus_sectors"`   // fintech, crypto, ai, saas
	Geography       []string               `json:"geography"`       // us, europe, asia, global
	CheckSize       CheckSizeRange         `json:"check_size"`
	Portfolio       []PortfolioCompany     `json:"portfolio"`
	KeyPersonnel    []InvestorContact      `json:"key_personnel"`
	InvestmentCriteria InvestmentCriteria  `json:"investment_criteria"`
	DecisionProcess DecisionProcess        `json:"decision_process"`
	ReputationScore decimal.Decimal        `json:"reputation_score"`
	NetworkValue    decimal.Decimal        `json:"network_value"`
	AddedValue      []string               `json:"added_value"`
	RecentActivity  []InvestmentActivity   `json:"recent_activity"`
	ContactHistory  []ContactInteraction   `json:"contact_history"`
	Status          string                 `json:"status"`          // prospect, contacted, interested, passed
	Priority        string                 `json:"priority"`        // high, medium, low
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// CheckSizeRange represents investor check size range
type CheckSizeRange struct {
	MinAmount decimal.Decimal `json:"min_amount"`
	MaxAmount decimal.Decimal `json:"max_amount"`
	TypicalAmount decimal.Decimal `json:"typical_amount"`
	Currency  string          `json:"currency"`
}

// PortfolioCompany represents portfolio company
type PortfolioCompany struct {
	Name        string          `json:"name"`
	Sector      string          `json:"sector"`
	Stage       string          `json:"stage"`
	Investment  decimal.Decimal `json:"investment"`
	Valuation   decimal.Decimal `json:"valuation"`
	Status      string          `json:"status"`      // active, exited, failed
	ExitValue   decimal.Decimal `json:"exit_value"`
	ExitType    string          `json:"exit_type"`   // ipo, acquisition, merger
}

// InvestorContact represents investor contact
type InvestorContact struct {
	Name        string   `json:"name"`
	Title       string   `json:"title"`
	Email       string   `json:"email"`
	LinkedIn    string   `json:"linkedin"`
	Role        string   `json:"role"`        // partner, principal, associate, analyst
	Seniority   string   `json:"seniority"`   // senior, junior
	Specialties []string `json:"specialties"`
	Background  string   `json:"background"`
}

// InvestmentCriteria represents investment criteria
type InvestmentCriteria struct {
	MinRevenue      decimal.Decimal `json:"min_revenue"`
	MinGrowthRate   decimal.Decimal `json:"min_growth_rate"`
	MinMarketSize   decimal.Decimal `json:"min_market_size"`
	RequiredMetrics []string        `json:"required_metrics"`
	RedFlags        []string        `json:"red_flags"`
	MustHaves       []string        `json:"must_haves"`
	NiceToHaves     []string        `json:"nice_to_haves"`
}

// DecisionProcess represents investor decision process
type DecisionProcess struct {
	TimelineWeeks   int      `json:"timeline_weeks"`
	DecisionMakers  []string `json:"decision_makers"`
	ProcessSteps    []string `json:"process_steps"`
	RequiredDocs    []string `json:"required_docs"`
	ReferenceChecks bool     `json:"reference_checks"`
	TechnicalDD     bool     `json:"technical_dd"`
	FinancialDD     bool     `json:"financial_dd"`
	LegalDD         bool     `json:"legal_dd"`
}

// InvestmentActivity represents recent investment activity
type InvestmentActivity struct {
	CompanyName   string          `json:"company_name"`
	Sector        string          `json:"sector"`
	Stage         string          `json:"stage"`
	Amount        decimal.Decimal `json:"amount"`
	Date          time.Time       `json:"date"`
	LeadInvestor  bool            `json:"lead_investor"`
	Description   string          `json:"description"`
}

// ContactInteraction represents contact interaction
type ContactInteraction struct {
	Date        time.Time `json:"date"`
	Type        string    `json:"type"`        // email, call, meeting, demo
	Participants []string `json:"participants"`
	Summary     string    `json:"summary"`
	Outcome     string    `json:"outcome"`
	NextSteps   []string  `json:"next_steps"`
	FollowUpDate time.Time `json:"follow_up_date"`
}

// CreateFundingRound creates a new funding round
func (irm *InvestorRelationsManager) CreateFundingRound(ctx context.Context, round *FundingRound) error {
	query := `
		INSERT INTO funding_rounds (
			id, name, type, status, target_amount, min_amount, max_amount,
			valuation, terms, timeline, use_of_funds, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := irm.db.ExecContext(ctx, query,
		round.ID, round.Name, round.Type, round.Status, round.TargetAmount,
		round.MinAmount, round.MaxAmount, round.Valuation, round.Terms,
		round.Timeline, round.UseOfFunds, round.CreatedAt, round.UpdatedAt,
	)

	return err
}

// GetInvestorsByType returns investors by type
func (irm *InvestorRelationsManager) GetInvestorsByType(ctx context.Context, investorType string, limit int) ([]*Investor, error) {
	query := `
		SELECT id, name, type, tier, focus_stage, focus_sectors, geography,
		       reputation_score, network_value, status, priority, created_at
		FROM investors 
		WHERE type = $1 AND status != 'passed'
		ORDER BY reputation_score DESC, network_value DESC
		LIMIT $2
	`

	rows, err := irm.db.QueryContext(ctx, query, investorType, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var investors []*Investor
	for rows.Next() {
		investor := &Investor{}
		err := rows.Scan(
			&investor.ID, &investor.Name, &investor.Type, &investor.Tier,
			&investor.FocusStage, &investor.FocusSectors, &investor.Geography,
			&investor.ReputationScore, &investor.NetworkValue, &investor.Status,
			&investor.Priority, &investor.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		investors = append(investors, investor)
	}

	return investors, nil
}

// Constructor functions
func NewPitchDeckManager() *PitchDeckManager {
	return &PitchDeckManager{}
}

func NewFinancialModelingEngine() *FinancialModelingEngine {
	return &FinancialModelingEngine{}
}

func NewInvestorCRM() *InvestorCRM {
	return &InvestorCRM{}
}

func NewVirtualDataRoom() *VirtualDataRoom {
	return &VirtualDataRoom{}
}

func NewFundingComplianceTracker() *FundingComplianceTracker {
	return &FundingComplianceTracker{}
}

// Placeholder types for other components
type PitchDeckManager struct{}
type FinancialModelingEngine struct{}
type InvestorCRM struct{}
type VirtualDataRoom struct{}
type FundingComplianceTracker struct{}
