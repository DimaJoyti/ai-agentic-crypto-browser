package compliance

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

// AuditTrail manages audit logging for compliance
type AuditTrail struct {
	logger    *observability.Logger
	config    ComplianceConfig
	events    map[string]*AuditEvent
	mu        sync.RWMutex
	isRunning int32
}

// AuditEvent represents an auditable event
type AuditEvent struct {
	ID          uuid.UUID              `json:"id"`
	UserID      string                 `json:"user_id,omitempty"`
	SessionID   string                 `json:"session_id,omitempty"`
	Action      AuditAction            `json:"action"`
	Resource    string                 `json:"resource"`
	ResourceID  string                 `json:"resource_id,omitempty"`
	Details     map[string]interface{} `json:"details"`
	IPAddress   string                 `json:"ip_address,omitempty"`
	UserAgent   string                 `json:"user_agent,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
	Success     bool                   `json:"success"`
	ErrorMsg    string                 `json:"error_message,omitempty"`
	RiskScore   float64                `json:"risk_score"`
	Compliance  AuditCompliance        `json:"compliance"`
}

// AuditAction represents the type of action being audited
type AuditAction string

const (
	// Trading Actions
	AuditActionOrderCreate    AuditAction = "ORDER_CREATE"
	AuditActionOrderModify    AuditAction = "ORDER_MODIFY"
	AuditActionOrderCancel    AuditAction = "ORDER_CANCEL"
	AuditActionOrderExecute   AuditAction = "ORDER_EXECUTE"
	AuditActionPositionOpen   AuditAction = "POSITION_OPEN"
	AuditActionPositionClose  AuditAction = "POSITION_CLOSE"
	AuditActionPositionModify AuditAction = "POSITION_MODIFY"

	// Risk Management Actions
	AuditActionRiskLimitSet    AuditAction = "RISK_LIMIT_SET"
	AuditActionRiskLimitBreach AuditAction = "RISK_LIMIT_BREACH"
	AuditActionEmergencyStop   AuditAction = "EMERGENCY_STOP"
	AuditActionCircuitBreaker  AuditAction = "CIRCUIT_BREAKER"

	// System Actions
	AuditActionLogin          AuditAction = "LOGIN"
	AuditActionLogout         AuditAction = "LOGOUT"
	AuditActionConfigChange   AuditAction = "CONFIG_CHANGE"
	AuditActionSystemStart    AuditAction = "SYSTEM_START"
	AuditActionSystemStop     AuditAction = "SYSTEM_STOP"

	// Compliance Actions
	AuditActionComplianceCheck   AuditAction = "COMPLIANCE_CHECK"
	AuditActionReportGenerate    AuditAction = "REPORT_GENERATE"
	AuditActionViolationDetected AuditAction = "VIOLATION_DETECTED"
	AuditActionViolationResolved AuditAction = "VIOLATION_RESOLVED"

	// Data Actions
	AuditActionDataAccess AuditAction = "DATA_ACCESS"
	AuditActionDataExport AuditAction = "DATA_EXPORT"
	AuditActionDataDelete AuditAction = "DATA_DELETE"
)

// AuditCompliance contains compliance-related audit information
type AuditCompliance struct {
	RequiredByFrameworks []string `json:"required_by_frameworks"`
	RetentionPeriod      int      `json:"retention_period_days"`
	Classification       string   `json:"classification"`
	Sensitive            bool     `json:"sensitive"`
}

// NewAuditTrail creates a new audit trail
func NewAuditTrail(logger *observability.Logger, config ComplianceConfig) *AuditTrail {
	return &AuditTrail{
		logger: logger,
		config: config,
		events: make(map[string]*AuditEvent),
	}
}

// Start starts the audit trail
func (at *AuditTrail) Start(ctx context.Context) error {
	at.logger.Info(ctx, "Starting audit trail", nil)
	at.isRunning = 1
	return nil
}

// Stop stops the audit trail
func (at *AuditTrail) Stop(ctx context.Context) error {
	at.logger.Info(ctx, "Stopping audit trail", nil)
	at.isRunning = 0
	return nil
}

// LogEvent logs an audit event
func (at *AuditTrail) LogEvent(ctx context.Context, event *AuditEvent) error {
	if at.isRunning == 0 {
		return fmt.Errorf("audit trail is not running")
	}

	// Set event ID and timestamp if not provided
	if event.ID == uuid.Nil {
		event.ID = uuid.New()
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// Calculate risk score if not provided
	if event.RiskScore == 0 {
		event.RiskScore = at.calculateRiskScore(event)
	}

	// Set compliance information
	event.Compliance = at.getComplianceInfo(event.Action)

	at.mu.Lock()
	at.events[event.ID.String()] = event
	at.mu.Unlock()

	// Log to structured logger
	at.logger.Info(ctx, "Audit event logged", map[string]interface{}{
		"audit_id":    event.ID,
		"user_id":     event.UserID,
		"action":      event.Action,
		"resource":    event.Resource,
		"success":     event.Success,
		"risk_score":  event.RiskScore,
		"ip_address":  event.IPAddress,
	})

	// Check for high-risk events
	if event.RiskScore >= 80 {
		at.logger.Warn(ctx, "High-risk audit event detected", map[string]interface{}{
			"audit_id":   event.ID,
			"action":     event.Action,
			"risk_score": event.RiskScore,
			"details":    event.Details,
		})
	}

	return nil
}

// LogTradingEvent logs a trading-related audit event
func (at *AuditTrail) LogTradingEvent(ctx context.Context, action AuditAction, userID, symbol string, details map[string]interface{}) error {
	event := &AuditEvent{
		ID:        uuid.New(),
		UserID:    userID,
		Action:    action,
		Resource:  "trading",
		Details:   details,
		Timestamp: time.Now(),
		Success:   true,
	}

	// Add symbol to details if provided
	if symbol != "" {
		if event.Details == nil {
			event.Details = make(map[string]interface{})
		}
		event.Details["symbol"] = symbol
	}

	return at.LogEvent(ctx, event)
}

// LogRiskEvent logs a risk management audit event
func (at *AuditTrail) LogRiskEvent(ctx context.Context, action AuditAction, details map[string]interface{}) error {
	event := &AuditEvent{
		ID:        uuid.New(),
		Action:    action,
		Resource:  "risk_management",
		Details:   details,
		Timestamp: time.Now(),
		Success:   true,
	}

	return at.LogEvent(ctx, event)
}

// LogComplianceEvent logs a compliance-related audit event
func (at *AuditTrail) LogComplianceEvent(ctx context.Context, action AuditAction, frameworkID string, details map[string]interface{}) error {
	event := &AuditEvent{
		ID:         uuid.New(),
		Action:     action,
		Resource:   "compliance",
		ResourceID: frameworkID,
		Details:    details,
		Timestamp:  time.Now(),
		Success:    true,
	}

	return at.LogEvent(ctx, event)
}

// GetEvents retrieves audit events with optional filtering
func (at *AuditTrail) GetEvents(ctx context.Context, filter AuditFilter) ([]*AuditEvent, error) {
	at.mu.RLock()
	defer at.mu.RUnlock()

	var events []*AuditEvent
	for _, event := range at.events {
		if at.matchesFilter(event, filter) {
			events = append(events, event)
		}
	}

	// Sort by timestamp (newest first)
	for i := 0; i < len(events)-1; i++ {
		for j := i + 1; j < len(events); j++ {
			if events[i].Timestamp.Before(events[j].Timestamp) {
				events[i], events[j] = events[j], events[i]
			}
		}
	}

	// Apply limit
	if filter.Limit > 0 && len(events) > filter.Limit {
		events = events[:filter.Limit]
	}

	return events, nil
}

// AuditFilter defines filtering criteria for audit events
type AuditFilter struct {
	UserID     string      `json:"user_id,omitempty"`
	Action     AuditAction `json:"action,omitempty"`
	Resource   string      `json:"resource,omitempty"`
	StartTime  *time.Time  `json:"start_time,omitempty"`
	EndTime    *time.Time  `json:"end_time,omitempty"`
	Success    *bool       `json:"success,omitempty"`
	MinRiskScore float64   `json:"min_risk_score,omitempty"`
	Limit      int         `json:"limit,omitempty"`
}

// matchesFilter checks if an event matches the filter criteria
func (at *AuditTrail) matchesFilter(event *AuditEvent, filter AuditFilter) bool {
	if filter.UserID != "" && event.UserID != filter.UserID {
		return false
	}
	if filter.Action != "" && event.Action != filter.Action {
		return false
	}
	if filter.Resource != "" && event.Resource != filter.Resource {
		return false
	}
	if filter.StartTime != nil && event.Timestamp.Before(*filter.StartTime) {
		return false
	}
	if filter.EndTime != nil && event.Timestamp.After(*filter.EndTime) {
		return false
	}
	if filter.Success != nil && event.Success != *filter.Success {
		return false
	}
	if filter.MinRiskScore > 0 && event.RiskScore < filter.MinRiskScore {
		return false
	}
	return true
}

// calculateRiskScore calculates a risk score for an audit event
func (at *AuditTrail) calculateRiskScore(event *AuditEvent) float64 {
	score := 0.0

	// Base score by action type
	switch event.Action {
	case AuditActionEmergencyStop, AuditActionCircuitBreaker:
		score += 90
	case AuditActionRiskLimitBreach, AuditActionViolationDetected:
		score += 80
	case AuditActionOrderCreate, AuditActionOrderModify, AuditActionPositionOpen:
		score += 30
	case AuditActionLogin, AuditActionConfigChange:
		score += 20
	case AuditActionDataAccess, AuditActionDataExport:
		score += 40
	case AuditActionDataDelete:
		score += 70
	default:
		score += 10
	}

	// Increase score for failures
	if !event.Success {
		score += 30
	}

	// Increase score for sensitive resources
	if event.Resource == "trading" || event.Resource == "risk_management" {
		score += 20
	}

	// Check for suspicious patterns in details
	if event.Details != nil {
		if amount, ok := event.Details["amount"]; ok {
			if amountFloat, ok := amount.(float64); ok && amountFloat > 1000000 {
				score += 25 // Large amounts are higher risk
			}
		}
	}

	// Cap at 100
	if score > 100 {
		score = 100
	}

	return score
}

// getComplianceInfo returns compliance information for an action
func (at *AuditTrail) getComplianceInfo(action AuditAction) AuditCompliance {
	compliance := AuditCompliance{
		RetentionPeriod: at.config.AuditRetentionDays,
		Classification:  "INTERNAL",
		Sensitive:       false,
	}

	// Set frameworks that require this audit
	switch action {
	case AuditActionOrderCreate, AuditActionOrderModify, AuditActionOrderCancel, AuditActionOrderExecute:
		compliance.RequiredByFrameworks = []string{"us_bsa", "eu_amld5", "mifid2"}
		compliance.Classification = "CONFIDENTIAL"
		compliance.Sensitive = true
	case AuditActionRiskLimitBreach, AuditActionEmergencyStop:
		compliance.RequiredByFrameworks = []string{"us_bsa", "eu_amld5", "basel3"}
		compliance.Classification = "CONFIDENTIAL"
		compliance.Sensitive = true
	case AuditActionComplianceCheck, AuditActionReportGenerate:
		compliance.RequiredByFrameworks = []string{"us_bsa", "eu_amld5"}
		compliance.Classification = "RESTRICTED"
	case AuditActionDataAccess, AuditActionDataExport, AuditActionDataDelete:
		compliance.RequiredByFrameworks = []string{"gdpr", "ccpa"}
		compliance.Classification = "CONFIDENTIAL"
		compliance.Sensitive = true
	}

	return compliance
}

// GetAuditSummary returns a summary of audit events
func (at *AuditTrail) GetAuditSummary(ctx context.Context, period time.Duration) (*AuditSummary, error) {
	startTime := time.Now().Add(-period)
	filter := AuditFilter{
		StartTime: &startTime,
	}

	events, err := at.GetEvents(ctx, filter)
	if err != nil {
		return nil, err
	}

	summary := &AuditSummary{
		Period:      period,
		TotalEvents: len(events),
		ActionCounts: make(map[AuditAction]int),
		ResourceCounts: make(map[string]int),
		RiskDistribution: make(map[string]int),
	}

	var totalRiskScore float64
	for _, event := range events {
		// Count by action
		summary.ActionCounts[event.Action]++

		// Count by resource
		summary.ResourceCounts[event.Resource]++

		// Count successful vs failed
		if event.Success {
			summary.SuccessfulEvents++
		} else {
			summary.FailedEvents++
		}

		// Risk distribution
		switch {
		case event.RiskScore >= 80:
			summary.RiskDistribution["HIGH"]++
		case event.RiskScore >= 60:
			summary.RiskDistribution["MEDIUM"]++
		case event.RiskScore >= 40:
			summary.RiskDistribution["LOW"]++
		default:
			summary.RiskDistribution["MINIMAL"]++
		}

		totalRiskScore += event.RiskScore
	}

	if len(events) > 0 {
		summary.AverageRiskScore = totalRiskScore / float64(len(events))
	}

	return summary, nil
}

// AuditSummary provides a summary of audit events
type AuditSummary struct {
	Period           time.Duration            `json:"period"`
	TotalEvents      int                      `json:"total_events"`
	SuccessfulEvents int                      `json:"successful_events"`
	FailedEvents     int                      `json:"failed_events"`
	AverageRiskScore float64                  `json:"average_risk_score"`
	ActionCounts     map[AuditAction]int      `json:"action_counts"`
	ResourceCounts   map[string]int           `json:"resource_counts"`
	RiskDistribution map[string]int           `json:"risk_distribution"`
}

// ExportEvents exports audit events to JSON format
func (at *AuditTrail) ExportEvents(ctx context.Context, filter AuditFilter) ([]byte, error) {
	events, err := at.GetEvents(ctx, filter)
	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(events, "", "  ")
}
