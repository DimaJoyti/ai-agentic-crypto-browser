package compliance

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ComplianceManager handles regulatory compliance and audit requirements
type ComplianceManager struct {
	logger          *observability.Logger
	config          ComplianceConfig
	auditTrail      *AuditTrail
	riskMonitor     *RiskMonitor
	reportGenerator *ReportGenerator
	alertManager    *AlertManager
	frameworks      map[string]*ComplianceFramework
	violations      map[string]*ComplianceViolation
	reports         map[string]*ComplianceReport
	mu              sync.RWMutex
	isRunning       int32
}

// ComplianceConfig contains compliance configuration
type ComplianceConfig struct {
	EnableAuditTrail     bool            `json:"enable_audit_trail"`
	EnableRiskMonitoring bool            `json:"enable_risk_monitoring"`
	EnableReporting      bool            `json:"enable_reporting"`
	AuditRetentionDays   int             `json:"audit_retention_days"`
	ReportingInterval    time.Duration   `json:"reporting_interval"`
	AlertThresholds      AlertThresholds `json:"alert_thresholds"`
	Jurisdictions        []string        `json:"jurisdictions"`
	EnabledFrameworks    []string        `json:"enabled_frameworks"`
}

// AlertThresholds defines compliance alert thresholds
type AlertThresholds struct {
	HighRiskScore      float64         `json:"high_risk_score"`
	CriticalRiskScore  float64         `json:"critical_risk_score"`
	MaxDailyViolations int             `json:"max_daily_violations"`
	MaxPositionSize    decimal.Decimal `json:"max_position_size"`
	MaxDailyVolume     decimal.Decimal `json:"max_daily_volume"`
}

// ComplianceFramework represents a regulatory framework
type ComplianceFramework struct {
	ID           string                  `json:"id"`
	Name         string                  `json:"name"`
	Description  string                  `json:"description"`
	Jurisdiction string                  `json:"jurisdiction"`
	Version      string                  `json:"version"`
	Requirements []ComplianceRequirement `json:"requirements"`
	Schedule     ReportingSchedule       `json:"schedule"`
	Status       ComplianceStatus        `json:"status"`
	LastUpdate   time.Time               `json:"last_update"`
}

// ComplianceRequirement represents a specific compliance requirement
type ComplianceRequirement struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Category    RequirementCategory `json:"category"`
	Mandatory   bool                `json:"mandatory"`
	Status      ComplianceStatus    `json:"status"`
	Controls    []ComplianceControl `json:"controls"`
	RiskLevel   RiskLevel           `json:"risk_level"`
	Deadline    *time.Time          `json:"deadline,omitempty"`
}

// ComplianceControl represents a control mechanism
type ComplianceControl struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	Type          ControlType   `json:"type"`
	Automated     bool          `json:"automated"`
	Frequency     string        `json:"frequency"`
	LastExecution *time.Time    `json:"last_execution,omitempty"`
	NextExecution *time.Time    `json:"next_execution,omitempty"`
	Effectiveness float64       `json:"effectiveness"`
	Status        ControlStatus `json:"status"`
}

// Enums for compliance types
type ComplianceStatus string

const (
	ComplianceStatusCompliant    ComplianceStatus = "COMPLIANT"
	ComplianceStatusNonCompliant ComplianceStatus = "NON_COMPLIANT"
	ComplianceStatusPartial      ComplianceStatus = "PARTIAL"
	ComplianceStatusPending      ComplianceStatus = "PENDING"
	ComplianceStatusUnknown      ComplianceStatus = "UNKNOWN"
)

type RequirementCategory string

const (
	RequirementCategoryKYC       RequirementCategory = "KYC"
	RequirementCategoryAML       RequirementCategory = "AML"
	RequirementCategoryReporting RequirementCategory = "REPORTING"
	RequirementCategoryRisk      RequirementCategory = "RISK"
	RequirementCategorySecurity  RequirementCategory = "SECURITY"
	RequirementCategoryData      RequirementCategory = "DATA"
)

type ControlType string

const (
	ControlTypePreventive   ControlType = "PREVENTIVE"
	ControlTypeDetective    ControlType = "DETECTIVE"
	ControlTypeCorrective   ControlType = "CORRECTIVE"
	ControlTypeCompensating ControlType = "COMPENSATING"
)

type ControlStatus string

const (
	ControlStatusActive   ControlStatus = "ACTIVE"
	ControlStatusInactive ControlStatus = "INACTIVE"
	ControlStatusFailed   ControlStatus = "FAILED"
	ControlStatusPending  ControlStatus = "PENDING"
)

type RiskLevel string

const (
	RiskLevelLow      RiskLevel = "LOW"
	RiskLevelMedium   RiskLevel = "MEDIUM"
	RiskLevelHigh     RiskLevel = "HIGH"
	RiskLevelCritical RiskLevel = "CRITICAL"
)

// ReportingSchedule defines when reports should be generated
type ReportingSchedule struct {
	Frequency     ReportingFrequency `json:"frequency"`
	NextDue       time.Time          `json:"next_due"`
	LastSubmitted *time.Time         `json:"last_submitted,omitempty"`
	Recipients    []string           `json:"recipients"`
	Format        ReportFormat       `json:"format"`
	Automated     bool               `json:"automated"`
}

type ReportingFrequency string

const (
	ReportingFrequencyDaily     ReportingFrequency = "DAILY"
	ReportingFrequencyWeekly    ReportingFrequency = "WEEKLY"
	ReportingFrequencyMonthly   ReportingFrequency = "MONTHLY"
	ReportingFrequencyQuarterly ReportingFrequency = "QUARTERLY"
	ReportingFrequencyAnnually  ReportingFrequency = "ANNUALLY"
	ReportingFrequencyOnDemand  ReportingFrequency = "ON_DEMAND"
)

type ReportFormat string

const (
	ReportFormatPDF  ReportFormat = "PDF"
	ReportFormatCSV  ReportFormat = "CSV"
	ReportFormatJSON ReportFormat = "JSON"
	ReportFormatXML  ReportFormat = "XML"
	ReportFormatXLSX ReportFormat = "XLSX"
)

// ComplianceViolation represents a compliance violation
type ComplianceViolation struct {
	ID            uuid.UUID              `json:"id"`
	FrameworkID   string                 `json:"framework_id"`
	RequirementID string                 `json:"requirement_id"`
	Type          ViolationType          `json:"type"`
	Severity      ViolationSeverity      `json:"severity"`
	Description   string                 `json:"description"`
	Details       map[string]interface{} `json:"details"`
	Timestamp     time.Time              `json:"timestamp"`
	Resolved      bool                   `json:"resolved"`
	ResolvedAt    *time.Time             `json:"resolved_at,omitempty"`
	ResolvedBy    string                 `json:"resolved_by,omitempty"`
	Actions       []RemediationAction    `json:"actions"`
}

type ViolationType string

const (
	ViolationTypeRisk      ViolationType = "RISK"
	ViolationTypeReporting ViolationType = "REPORTING"
	ViolationTypeKYC       ViolationType = "KYC"
	ViolationTypeAML       ViolationType = "AML"
	ViolationTypeSecurity  ViolationType = "SECURITY"
	ViolationTypeData      ViolationType = "DATA"
)

type ViolationSeverity string

const (
	ViolationSeverityLow      ViolationSeverity = "LOW"
	ViolationSeverityMedium   ViolationSeverity = "MEDIUM"
	ViolationSeverityHigh     ViolationSeverity = "HIGH"
	ViolationSeverityCritical ViolationSeverity = "CRITICAL"
)

// RemediationAction represents an action taken to resolve a violation
type RemediationAction struct {
	ID          uuid.UUID `json:"id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	TakenBy     string    `json:"taken_by"`
	TakenAt     time.Time `json:"taken_at"`
	Effective   bool      `json:"effective"`
}

// ComplianceReport represents a compliance report
type ComplianceReport struct {
	ID          uuid.UUID            `json:"id"`
	FrameworkID string               `json:"framework_id"`
	Type        ReportType           `json:"type"`
	Period      ReportPeriod         `json:"period"`
	Status      ReportStatus         `json:"status"`
	Data        ComplianceReportData `json:"data"`
	GeneratedAt time.Time            `json:"generated_at"`
	SubmittedAt *time.Time           `json:"submitted_at,omitempty"`
	Recipients  []string             `json:"recipients"`
	FilePath    string               `json:"file_path,omitempty"`
}

type ReportType string

const (
	ReportTypeCompliance  ReportType = "COMPLIANCE"
	ReportTypeRisk        ReportType = "RISK"
	ReportTypeAudit       ReportType = "AUDIT"
	ReportTypePosition    ReportType = "POSITION"
	ReportTypeTransaction ReportType = "TRANSACTION"
)

type ReportStatus string

const (
	ReportStatusDraft     ReportStatus = "DRAFT"
	ReportStatusGenerated ReportStatus = "GENERATED"
	ReportStatusSubmitted ReportStatus = "SUBMITTED"
	ReportStatusApproved  ReportStatus = "APPROVED"
	ReportStatusRejected  ReportStatus = "REJECTED"
)

// ReportPeriod defines the time period for a report
type ReportPeriod struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Label     string    `json:"label"`
}

// ComplianceReportData contains the actual report data
type ComplianceReportData struct {
	Summary         ComplianceSummary     `json:"summary"`
	Metrics         ComplianceMetrics     `json:"metrics"`
	Violations      []ComplianceViolation `json:"violations"`
	Findings        []ComplianceFinding   `json:"findings"`
	Recommendations []string              `json:"recommendations"`
	Attachments     []ReportAttachment    `json:"attachments"`
}

// ComplianceSummary provides high-level compliance overview
type ComplianceSummary struct {
	TotalRequirements        int       `json:"total_requirements"`
	CompliantRequirements    int       `json:"compliant_requirements"`
	NonCompliantRequirements int       `json:"non_compliant_requirements"`
	PendingRequirements      int       `json:"pending_requirements"`
	OverallComplianceRate    float64   `json:"overall_compliance_rate"`
	RiskLevel                RiskLevel `json:"risk_level"`
}

// ComplianceMetrics contains detailed compliance metrics
type ComplianceMetrics struct {
	TotalTrades           int             `json:"total_trades"`
	TotalVolume           decimal.Decimal `json:"total_volume"`
	AverageTradeSize      decimal.Decimal `json:"average_trade_size"`
	LargestTrade          decimal.Decimal `json:"largest_trade"`
	ViolationsCount       int             `json:"violations_count"`
	ResolvedViolations    int             `json:"resolved_violations"`
	AverageResolutionTime float64         `json:"average_resolution_time"` // hours
	ComplianceScore       float64         `json:"compliance_score"`
}

// ComplianceFinding represents a compliance finding
type ComplianceFinding struct {
	ID          uuid.UUID  `json:"id"`
	Type        string     `json:"type"`
	Severity    RiskLevel  `json:"severity"`
	Description string     `json:"description"`
	Requirement string     `json:"requirement"`
	Evidence    []string   `json:"evidence"`
	Remediation string     `json:"remediation"`
	Deadline    *time.Time `json:"deadline,omitempty"`
	Responsible string     `json:"responsible"`
}

// ReportAttachment represents a file attachment to a report
type ReportAttachment struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Size       int64     `json:"size"`
	Path       string    `json:"path"`
	Hash       string    `json:"hash"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// NewComplianceManager creates a new compliance manager
func NewComplianceManager(logger *observability.Logger, config ComplianceConfig) *ComplianceManager {
	cm := &ComplianceManager{
		logger:     logger,
		config:     config,
		frameworks: make(map[string]*ComplianceFramework),
		violations: make(map[string]*ComplianceViolation),
		reports:    make(map[string]*ComplianceReport),
	}

	// Initialize components
	cm.auditTrail = NewAuditTrail(logger, config)
	cm.riskMonitor = NewRiskMonitor(logger, config)
	cm.reportGenerator = NewReportGenerator(logger, config)
	cm.alertManager = NewAlertManager(logger, config)

	// Initialize default frameworks
	cm.initializeDefaultFrameworks()

	return cm
}

// Start starts the compliance manager
func (cm *ComplianceManager) Start(ctx context.Context) error {
	cm.logger.Info(ctx, "Starting compliance manager", nil)

	// Start components
	if err := cm.auditTrail.Start(ctx); err != nil {
		return fmt.Errorf("failed to start audit trail: %w", err)
	}

	if err := cm.riskMonitor.Start(ctx); err != nil {
		return fmt.Errorf("failed to start risk monitor: %w", err)
	}

	if err := cm.reportGenerator.Start(ctx); err != nil {
		return fmt.Errorf("failed to start report generator: %w", err)
	}

	if err := cm.alertManager.Start(ctx); err != nil {
		return fmt.Errorf("failed to start alert manager: %w", err)
	}

	cm.isRunning = 1
	cm.logger.Info(ctx, "Compliance manager started successfully", nil)
	return nil
}

// Stop stops the compliance manager
func (cm *ComplianceManager) Stop(ctx context.Context) error {
	cm.logger.Info(ctx, "Stopping compliance manager", nil)

	// Stop components
	cm.auditTrail.Stop(ctx)
	cm.riskMonitor.Stop(ctx)
	cm.reportGenerator.Stop(ctx)
	cm.alertManager.Stop(ctx)

	cm.isRunning = 0
	cm.logger.Info(ctx, "Compliance manager stopped", nil)
	return nil
}

// initializeDefaultFrameworks sets up default compliance frameworks
func (cm *ComplianceManager) initializeDefaultFrameworks() {
	// US BSA Framework
	usBSA := &ComplianceFramework{
		ID:           "us_bsa",
		Name:         "US Bank Secrecy Act",
		Description:  "US Anti-Money Laundering regulations",
		Jurisdiction: "United States",
		Version:      "2023.1",
		Status:       ComplianceStatusCompliant,
		LastUpdate:   time.Now(),
		Schedule: ReportingSchedule{
			Frequency:  ReportingFrequencyMonthly,
			NextDue:    time.Now().AddDate(0, 1, 0),
			Recipients: []string{"compliance@example.com"},
			Format:     ReportFormatPDF,
			Automated:  true,
		},
		Requirements: []ComplianceRequirement{
			{
				ID:          "customer_identification",
				Name:        "Customer Identification Program",
				Description: "Verify customer identity before account opening",
				Category:    RequirementCategoryKYC,
				Mandatory:   true,
				Status:      ComplianceStatusCompliant,
				RiskLevel:   RiskLevelHigh,
				Controls: []ComplianceControl{
					{
						ID:            "kyc_verification",
						Name:          "KYC Verification Process",
						Description:   "Automated KYC verification for all customers",
						Type:          ControlTypePreventive,
						Automated:     true,
						Frequency:     "per_transaction",
						Effectiveness: 0.95,
						Status:        ControlStatusActive,
					},
				},
			},
		},
	}

	// EU AMLD5 Framework
	euAMLD5 := &ComplianceFramework{
		ID:           "eu_amld5",
		Name:         "EU Anti-Money Laundering Directive 5",
		Description:  "European Union AML regulations",
		Jurisdiction: "European Union",
		Version:      "2020.1",
		Status:       ComplianceStatusPartial,
		LastUpdate:   time.Now(),
		Schedule: ReportingSchedule{
			Frequency:  ReportingFrequencyQuarterly,
			NextDue:    time.Now().AddDate(0, 3, 0),
			Recipients: []string{"eu-compliance@example.com"},
			Format:     ReportFormatXML,
			Automated:  false,
		},
		Requirements: []ComplianceRequirement{
			{
				ID:          "customer_due_diligence",
				Name:        "Customer Due Diligence",
				Description: "Perform enhanced due diligence on high-risk customers",
				Category:    RequirementCategoryKYC,
				Mandatory:   true,
				Status:      ComplianceStatusPartial,
				RiskLevel:   RiskLevelHigh,
			},
		},
	}

	cm.frameworks["us_bsa"] = usBSA
	cm.frameworks["eu_amld5"] = euAMLD5
}
