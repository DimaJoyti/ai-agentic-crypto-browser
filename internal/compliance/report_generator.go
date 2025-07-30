package compliance

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ReportGenerator handles compliance report generation
type ReportGenerator struct {
	logger           *observability.Logger
	config           ComplianceConfig
	reports          map[string]*ComplianceReport
	templates        map[string]*ReportTemplate
	scheduledReports map[string]*ScheduledReport
	mu               sync.RWMutex
	isRunning        int32
	stopChan         chan struct{}
}

// ReportTemplate defines a report template
type ReportTemplate struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Type        ReportType        `json:"type"`
	Framework   string            `json:"framework"`
	Sections    []ReportSection   `json:"sections"`
	Parameters  map[string]string `json:"parameters"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// ReportSection defines a section within a report
type ReportSection struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Type       SectionType       `json:"type"`
	Required   bool              `json:"required"`
	DataSource string            `json:"data_source"`
	Query      string            `json:"query,omitempty"`
	Parameters map[string]string `json:"parameters"`
	Order      int               `json:"order"`
}

// SectionType defines types of report sections
type SectionType string

const (
	SectionTypeSummary      SectionType = "SUMMARY"
	SectionTypeMetrics      SectionType = "METRICS"
	SectionTypeTransactions SectionType = "TRANSACTIONS"
	SectionTypePositions    SectionType = "POSITIONS"
	SectionTypeViolations   SectionType = "VIOLATIONS"
	SectionTypeRisk         SectionType = "RISK"
	SectionTypeAudit        SectionType = "AUDIT"
	SectionTypeChart        SectionType = "CHART"
	SectionTypeTable        SectionType = "TABLE"
)

// ScheduledReport defines a scheduled report
type ScheduledReport struct {
	ID         uuid.UUID         `json:"id"`
	TemplateID string            `json:"template_id"`
	Name       string            `json:"name"`
	Schedule   ReportingSchedule `json:"schedule"`
	Parameters map[string]string `json:"parameters"`
	Enabled    bool              `json:"enabled"`
	LastRun    *time.Time        `json:"last_run,omitempty"`
	NextRun    time.Time         `json:"next_run"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

// NewReportGenerator creates a new report generator
func NewReportGenerator(logger *observability.Logger, config ComplianceConfig) *ReportGenerator {
	rg := &ReportGenerator{
		logger:           logger,
		config:           config,
		reports:          make(map[string]*ComplianceReport),
		templates:        make(map[string]*ReportTemplate),
		scheduledReports: make(map[string]*ScheduledReport),
		stopChan:         make(chan struct{}),
	}

	// Initialize default templates
	rg.initializeDefaultTemplates()

	return rg
}

// Start starts the report generator
func (rg *ReportGenerator) Start(ctx context.Context) error {
	rg.logger.Info(ctx, "Starting report generator", nil)
	rg.isRunning = 1

	// Start scheduling goroutine
	go rg.schedulingLoop(ctx)

	return nil
}

// Stop stops the report generator
func (rg *ReportGenerator) Stop(ctx context.Context) error {
	rg.logger.Info(ctx, "Stopping report generator", nil)
	rg.isRunning = 0
	close(rg.stopChan)
	return nil
}

// schedulingLoop runs the report scheduling loop
func (rg *ReportGenerator) schedulingLoop(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute) // Check every minute
	defer ticker.Stop()

	for {
		select {
		case <-rg.stopChan:
			return
		case <-ticker.C:
			rg.checkScheduledReports(ctx)
		}
	}
}

// checkScheduledReports checks for reports that need to be generated
func (rg *ReportGenerator) checkScheduledReports(ctx context.Context) {
	rg.mu.RLock()
	defer rg.mu.RUnlock()

	now := time.Now()
	for _, scheduled := range rg.scheduledReports {
		if scheduled.Enabled && now.After(scheduled.NextRun) {
			go rg.generateScheduledReport(ctx, scheduled)
		}
	}
}

// generateScheduledReport generates a scheduled report
func (rg *ReportGenerator) generateScheduledReport(ctx context.Context, scheduled *ScheduledReport) {
	rg.logger.Info(ctx, "Generating scheduled report", map[string]interface{}{
		"report_id":   scheduled.ID,
		"template_id": scheduled.TemplateID,
		"name":        scheduled.Name,
	})

	template, exists := rg.templates[scheduled.TemplateID]
	if !exists {
		rg.logger.Info(ctx, "Template not found for scheduled report", map[string]interface{}{
			"template_id": scheduled.TemplateID,
		})
		return
	}

	// Generate report
	period := rg.calculateReportPeriod(scheduled.Schedule.Frequency)
	report, err := rg.GenerateReport(ctx, template, period, scheduled.Parameters)
	if err != nil {
		rg.logger.Info(ctx, "Failed to generate scheduled report", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Update scheduled report
	rg.mu.Lock()
	now := time.Now()
	scheduled.LastRun = &now
	scheduled.NextRun = rg.calculateNextRun(scheduled.Schedule.Frequency, now)
	rg.mu.Unlock()

	rg.logger.Info(ctx, "Scheduled report generated successfully", map[string]interface{}{
		"report_id":    report.ID,
		"scheduled_id": scheduled.ID,
	})
}

// GenerateReport generates a compliance report
func (rg *ReportGenerator) GenerateReport(ctx context.Context, template *ReportTemplate, period ReportPeriod, parameters map[string]string) (*ComplianceReport, error) {
	reportID := uuid.New()

	rg.logger.Info(ctx, "Generating compliance report", map[string]interface{}{
		"report_id":   reportID,
		"template_id": template.ID,
		"type":        template.Type,
		"framework":   template.Framework,
	})

	// Create report structure
	report := &ComplianceReport{
		ID:          reportID,
		FrameworkID: template.Framework,
		Type:        template.Type,
		Period:      period,
		Status:      ReportStatusDraft,
		GeneratedAt: time.Now(),
		Recipients:  []string{}, // Will be populated from framework
	}

	// Generate report data
	reportData, err := rg.generateReportData(ctx, template, period, parameters)
	if err != nil {
		return nil, fmt.Errorf("failed to generate report data: %w", err)
	}

	report.Data = *reportData
	report.Status = ReportStatusGenerated

	// Store report
	rg.mu.Lock()
	rg.reports[reportID.String()] = report
	rg.mu.Unlock()

	rg.logger.Info(ctx, "Compliance report generated successfully", map[string]interface{}{
		"report_id": reportID,
		"type":      template.Type,
		"framework": template.Framework,
	})

	return report, nil
}

// generateReportData generates the actual report data
func (rg *ReportGenerator) generateReportData(ctx context.Context, template *ReportTemplate, period ReportPeriod, parameters map[string]string) (*ComplianceReportData, error) {
	data := &ComplianceReportData{
		Summary:         rg.generateSummary(ctx, template, period),
		Metrics:         rg.generateMetrics(ctx, template, period),
		Violations:      rg.generateViolations(ctx, template, period),
		Findings:        rg.generateFindings(ctx, template, period),
		Recommendations: rg.generateRecommendations(ctx, template, period),
		Attachments:     []ReportAttachment{},
	}

	return data, nil
}

// generateSummary generates the compliance summary
func (rg *ReportGenerator) generateSummary(ctx context.Context, template *ReportTemplate, period ReportPeriod) ComplianceSummary {
	// Mock data - in production, query actual compliance data
	return ComplianceSummary{
		TotalRequirements:        10,
		CompliantRequirements:    8,
		NonCompliantRequirements: 1,
		PendingRequirements:      1,
		OverallComplianceRate:    80.0,
		RiskLevel:                RiskLevelMedium,
	}
}

// generateMetrics generates compliance metrics
func (rg *ReportGenerator) generateMetrics(ctx context.Context, template *ReportTemplate, period ReportPeriod) ComplianceMetrics {
	// Mock data - in production, query actual trading data
	return ComplianceMetrics{
		TotalTrades:           1500,
		TotalVolume:           decimal.NewFromInt(50000000), // $50M
		AverageTradeSize:      decimal.NewFromInt(33333),    // ~$33k
		LargestTrade:          decimal.NewFromInt(500000),   // $500k
		ViolationsCount:       5,
		ResolvedViolations:    4,
		AverageResolutionTime: 2.5, // 2.5 hours
		ComplianceScore:       85.0,
	}
}

// generateViolations generates violation data
func (rg *ReportGenerator) generateViolations(ctx context.Context, template *ReportTemplate, period ReportPeriod) []ComplianceViolation {
	// Mock data - in production, query actual violations
	violations := []ComplianceViolation{
		{
			ID:            uuid.New(),
			FrameworkID:   template.Framework,
			RequirementID: "position_limit",
			Type:          ViolationTypeRisk,
			Severity:      ViolationSeverityMedium,
			Description:   "Position limit exceeded for BTCUSD",
			Details: map[string]interface{}{
				"symbol":        "BTCUSD",
				"position_size": 1500000,
				"limit":         1000000,
			},
			Timestamp: time.Now().Add(-2 * time.Hour),
			Resolved:  true,
		},
	}

	return violations
}

// generateFindings generates compliance findings
func (rg *ReportGenerator) generateFindings(ctx context.Context, template *ReportTemplate, period ReportPeriod) []ComplianceFinding {
	// Mock data - in production, analyze actual compliance data
	findings := []ComplianceFinding{
		{
			ID:          uuid.New(),
			Type:        "risk_management",
			Severity:    RiskLevelMedium,
			Description: "Position concentration risk detected",
			Requirement: "Risk diversification requirements",
			Evidence:    []string{"Position report", "Risk metrics"},
			Remediation: "Implement position size limits and diversification rules",
			Responsible: "Risk Management Team",
		},
	}

	return findings
}

// generateRecommendations generates compliance recommendations
func (rg *ReportGenerator) generateRecommendations(ctx context.Context, template *ReportTemplate, period ReportPeriod) []string {
	recommendations := []string{
		"Implement automated position size monitoring",
		"Enhance real-time risk alerting system",
		"Conduct quarterly compliance training",
		"Review and update risk limits monthly",
		"Implement additional correlation risk controls",
	}

	return recommendations
}

// calculateReportPeriod calculates the report period based on frequency
func (rg *ReportGenerator) calculateReportPeriod(frequency ReportingFrequency) ReportPeriod {
	now := time.Now()
	var startDate time.Time
	var label string

	switch frequency {
	case ReportingFrequencyDaily:
		startDate = now.AddDate(0, 0, -1)
		label = "Daily Report"
	case ReportingFrequencyWeekly:
		startDate = now.AddDate(0, 0, -7)
		label = "Weekly Report"
	case ReportingFrequencyMonthly:
		startDate = now.AddDate(0, -1, 0)
		label = "Monthly Report"
	case ReportingFrequencyQuarterly:
		startDate = now.AddDate(0, -3, 0)
		label = "Quarterly Report"
	case ReportingFrequencyAnnually:
		startDate = now.AddDate(-1, 0, 0)
		label = "Annual Report"
	default:
		startDate = now.AddDate(0, -1, 0)
		label = "Monthly Report"
	}

	return ReportPeriod{
		StartDate: startDate,
		EndDate:   now,
		Label:     label,
	}
}

// calculateNextRun calculates the next run time for a scheduled report
func (rg *ReportGenerator) calculateNextRun(frequency ReportingFrequency, lastRun time.Time) time.Time {
	switch frequency {
	case ReportingFrequencyDaily:
		return lastRun.AddDate(0, 0, 1)
	case ReportingFrequencyWeekly:
		return lastRun.AddDate(0, 0, 7)
	case ReportingFrequencyMonthly:
		return lastRun.AddDate(0, 1, 0)
	case ReportingFrequencyQuarterly:
		return lastRun.AddDate(0, 3, 0)
	case ReportingFrequencyAnnually:
		return lastRun.AddDate(1, 0, 0)
	default:
		return lastRun.AddDate(0, 1, 0)
	}
}

// initializeDefaultTemplates sets up default report templates
func (rg *ReportGenerator) initializeDefaultTemplates() {
	templates := []*ReportTemplate{
		{
			ID:          "us_bsa_monthly",
			Name:        "US BSA Monthly Compliance Report",
			Description: "Monthly compliance report for US Bank Secrecy Act",
			Type:        ReportTypeCompliance,
			Framework:   "us_bsa",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Sections: []ReportSection{
				{
					ID:         "summary",
					Name:       "Executive Summary",
					Type:       SectionTypeSummary,
					Required:   true,
					DataSource: "compliance_metrics",
					Order:      1,
				},
				{
					ID:         "violations",
					Name:       "Compliance Violations",
					Type:       SectionTypeViolations,
					Required:   true,
					DataSource: "violation_log",
					Order:      2,
				},
				{
					ID:         "risk_metrics",
					Name:       "Risk Metrics",
					Type:       SectionTypeRisk,
					Required:   true,
					DataSource: "risk_monitor",
					Order:      3,
				},
			},
		},
		{
			ID:          "eu_amld5_quarterly",
			Name:        "EU AMLD5 Quarterly Report",
			Description: "Quarterly compliance report for EU Anti-Money Laundering Directive 5",
			Type:        ReportTypeCompliance,
			Framework:   "eu_amld5",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Sections: []ReportSection{
				{
					ID:         "summary",
					Name:       "Compliance Overview",
					Type:       SectionTypeSummary,
					Required:   true,
					DataSource: "compliance_metrics",
					Order:      1,
				},
				{
					ID:         "transactions",
					Name:       "Transaction Analysis",
					Type:       SectionTypeTransactions,
					Required:   true,
					DataSource: "transaction_log",
					Order:      2,
				},
			},
		},
		{
			ID:          "risk_daily",
			Name:        "Daily Risk Report",
			Description: "Daily risk monitoring and metrics report",
			Type:        ReportTypeRisk,
			Framework:   "internal",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Sections: []ReportSection{
				{
					ID:         "risk_summary",
					Name:       "Risk Summary",
					Type:       SectionTypeRisk,
					Required:   true,
					DataSource: "risk_monitor",
					Order:      1,
				},
				{
					ID:         "positions",
					Name:       "Position Analysis",
					Type:       SectionTypePositions,
					Required:   true,
					DataSource: "position_manager",
					Order:      2,
				},
			},
		},
	}

	for _, template := range templates {
		rg.templates[template.ID] = template
	}
}

// GetReports returns reports with optional filtering
func (rg *ReportGenerator) GetReports(filter ReportFilter) []*ComplianceReport {
	rg.mu.RLock()
	defer rg.mu.RUnlock()

	var reports []*ComplianceReport
	for _, report := range rg.reports {
		if rg.matchesReportFilter(report, filter) {
			reports = append(reports, report)
		}
	}

	return reports
}

// ReportFilter defines filtering criteria for reports
type ReportFilter struct {
	Type      ReportType   `json:"type,omitempty"`
	Framework string       `json:"framework,omitempty"`
	Status    ReportStatus `json:"status,omitempty"`
	StartDate *time.Time   `json:"start_date,omitempty"`
	EndDate   *time.Time   `json:"end_date,omitempty"`
	Limit     int          `json:"limit,omitempty"`
}

// matchesReportFilter checks if a report matches the filter criteria
func (rg *ReportGenerator) matchesReportFilter(report *ComplianceReport, filter ReportFilter) bool {
	if filter.Type != "" && report.Type != filter.Type {
		return false
	}
	if filter.Framework != "" && report.FrameworkID != filter.Framework {
		return false
	}
	if filter.Status != "" && report.Status != filter.Status {
		return false
	}
	if filter.StartDate != nil && report.GeneratedAt.Before(*filter.StartDate) {
		return false
	}
	if filter.EndDate != nil && report.GeneratedAt.After(*filter.EndDate) {
		return false
	}
	return true
}

// ExportReport exports a report to the specified format
func (rg *ReportGenerator) ExportReport(ctx context.Context, reportID string, format ReportFormat) ([]byte, error) {
	rg.mu.RLock()
	report, exists := rg.reports[reportID]
	rg.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("report not found: %s", reportID)
	}

	switch format {
	case ReportFormatJSON:
		return json.MarshalIndent(report, "", "  ")
	case ReportFormatCSV:
		return rg.exportToCSV(report)
	case ReportFormatPDF:
		return rg.exportToPDF(report)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// exportToCSV exports report to CSV format
func (rg *ReportGenerator) exportToCSV(report *ComplianceReport) ([]byte, error) {
	// Simplified CSV export - in production, use proper CSV library
	csv := "Report ID,Type,Framework,Generated At,Status\n"
	csv += fmt.Sprintf("%s,%s,%s,%s,%s\n",
		report.ID,
		report.Type,
		report.FrameworkID,
		report.GeneratedAt.Format(time.RFC3339),
		report.Status,
	)
	return []byte(csv), nil
}

// exportToPDF exports report to PDF format
func (rg *ReportGenerator) exportToPDF(report *ComplianceReport) ([]byte, error) {
	// Mock PDF export - in production, use proper PDF library
	return []byte("PDF content placeholder"), nil
}
