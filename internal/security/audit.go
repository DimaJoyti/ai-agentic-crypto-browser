package security

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

// AuditManager handles comprehensive audit logging
type AuditManager struct {
	logger            *observability.Logger
	config            *AuditConfig
	encryptionManager *EncryptionManager
	auditStore        *AuditStore
	eventProcessor    *AuditEventProcessor
	complianceEngine  *ComplianceEngine
	mu                sync.RWMutex
}

// AuditConfig contains audit configuration
type AuditConfig struct {
	EnableAuditLogging     bool          `json:"enable_audit_logging"`
	EnableRealTimeAuditing bool          `json:"enable_real_time_auditing"`
	EnableEncryption       bool          `json:"enable_encryption"`
	RetentionPeriod        time.Duration `json:"retention_period"`
	ComplianceStandards    []string      `json:"compliance_standards"` // SOX, GDPR, HIPAA, PCI-DSS
	AuditLevel             AuditLevel    `json:"audit_level"`
	EnableIntegrityCheck   bool          `json:"enable_integrity_check"`
	EnableTamperDetection  bool          `json:"enable_tamper_detection"`
	MaxAuditLogSize        int64         `json:"max_audit_log_size"`
	ArchiveThreshold       int64         `json:"archive_threshold"`
}

// AuditLevel defines the level of audit logging
type AuditLevel string

const (
	AuditLevelMinimal       AuditLevel = "minimal"       // Critical events only
	AuditLevelStandard      AuditLevel = "standard"      // Standard business events
	AuditLevelDetailed      AuditLevel = "detailed"      // Detailed system events
	AuditLevelComprehensive AuditLevel = "comprehensive" // All events
)

// AuditEvent represents an audit event
type AuditEvent struct {
	EventID       string                 `json:"event_id"`
	Timestamp     time.Time              `json:"timestamp"`
	EventType     AuditEventType         `json:"event_type"`
	Category      AuditCategory          `json:"category"`
	Severity      AuditSeverity          `json:"severity"`
	UserID        *uuid.UUID             `json:"user_id,omitempty"`
	SessionID     string                 `json:"session_id,omitempty"`
	IPAddress     string                 `json:"ip_address,omitempty"`
	UserAgent     string                 `json:"user_agent,omitempty"`
	Resource      string                 `json:"resource,omitempty"`
	Action        string                 `json:"action"`
	Result        AuditResult            `json:"result"`
	Details       map[string]interface{} `json:"details,omitempty"`
	RiskScore     float64                `json:"risk_score,omitempty"`
	ComplianceTag string                 `json:"compliance_tag,omitempty"`
	Hash          string                 `json:"hash"` // For integrity verification
	PreviousHash  string                 `json:"previous_hash,omitempty"`
}

// AuditEventType defines types of audit events
type AuditEventType string

const (
	AuditEventTypeAuthentication   AuditEventType = "authentication"
	AuditEventTypeAuthorization    AuditEventType = "authorization"
	AuditEventTypeDataAccess       AuditEventType = "data_access"
	AuditEventTypeDataModification AuditEventType = "data_modification"
	AuditEventTypeSystemAccess     AuditEventType = "system_access"
	AuditEventTypeConfiguration    AuditEventType = "configuration"
	AuditEventTypeSecurity         AuditEventType = "security"
	AuditEventTypeCompliance       AuditEventType = "compliance"
	AuditEventTypeTrading          AuditEventType = "trading"
	AuditEventTypeFinancial        AuditEventType = "financial"
)

// AuditCategory defines audit categories
type AuditCategory string

const (
	AuditCategorySecurity    AuditCategory = "security"
	AuditCategoryCompliance  AuditCategory = "compliance"
	AuditCategoryBusiness    AuditCategory = "business"
	AuditCategoryTechnical   AuditCategory = "technical"
	AuditCategoryOperational AuditCategory = "operational"
)

// AuditSeverity defines audit severity levels
type AuditSeverity string

const (
	AuditSeverityLow      AuditSeverity = "low"
	AuditSeverityMedium   AuditSeverity = "medium"
	AuditSeverityHigh     AuditSeverity = "high"
	AuditSeverityCritical AuditSeverity = "critical"
)

// AuditResult defines audit result
type AuditResult string

const (
	AuditResultSuccess AuditResult = "success"
	AuditResultFailure AuditResult = "failure"
	AuditResultDenied  AuditResult = "denied"
	AuditResultError   AuditResult = "error"
)

// AuditStore handles audit log storage
type AuditStore struct {
	logger     *observability.Logger
	config     *AuditConfig
	events     []AuditEvent
	eventIndex map[string]*AuditEvent
	lastHash   string
	totalSize  int64
	mu         sync.RWMutex
}

// AuditEventProcessor processes audit events
type AuditEventProcessor struct {
	logger    *observability.Logger
	config    *AuditConfig
	filters   []AuditFilter
	enrichers []AuditEnricher
	mu        sync.RWMutex
}

// AuditFilter filters audit events
type AuditFilter interface {
	ShouldAudit(event *AuditEvent) bool
}

// AuditEnricher enriches audit events
type AuditEnricher interface {
	Enrich(event *AuditEvent) error
}

// ComplianceEngine handles compliance-specific audit requirements
type ComplianceEngine struct {
	logger     *observability.Logger
	config     *AuditConfig
	standards  map[string]*ComplianceStandard
	violations []ComplianceViolation
	mu         sync.RWMutex
}

// ComplianceStandard defines a compliance standard
type ComplianceStandard struct {
	Name         string                  `json:"name"`
	Requirements []ComplianceRequirement `json:"requirements"`
	Enabled      bool                    `json:"enabled"`
}

// ComplianceRequirement defines a compliance requirement
type ComplianceRequirement struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	EventTypes  []string `json:"event_types"`
	Mandatory   bool     `json:"mandatory"`
}

// ComplianceViolation represents a compliance violation
type ComplianceViolation struct {
	ViolationID string    `json:"violation_id"`
	Standard    string    `json:"standard"`
	Requirement string    `json:"requirement"`
	EventID     string    `json:"event_id"`
	Timestamp   time.Time `json:"timestamp"`
	Severity    string    `json:"severity"`
	Description string    `json:"description"`
}

// NewAuditManager creates a new audit manager
func NewAuditManager(logger *observability.Logger, config *AuditConfig, encryptionManager *EncryptionManager) *AuditManager {
	am := &AuditManager{
		logger:            logger,
		config:            config,
		encryptionManager: encryptionManager,
	}

	// Initialize components
	am.auditStore = NewAuditStore(logger, config)
	am.eventProcessor = NewAuditEventProcessor(logger, config)
	am.complianceEngine = NewComplianceEngine(logger, config)

	return am
}

// Start starts the audit manager
func (am *AuditManager) Start(ctx context.Context) error {
	am.logger.Info(ctx, "Starting audit manager", map[string]interface{}{
		"audit_level":          am.config.AuditLevel,
		"compliance_standards": am.config.ComplianceStandards,
		"encryption_enabled":   am.config.EnableEncryption,
	})

	// Initialize compliance standards
	if err := am.complianceEngine.InitializeStandards(); err != nil {
		return fmt.Errorf("failed to initialize compliance standards: %w", err)
	}

	return nil
}

// LogEvent logs an audit event
func (am *AuditManager) LogEvent(ctx context.Context, event *AuditEvent) error {
	// Set event ID and timestamp if not provided
	if event.EventID == "" {
		event.EventID = uuid.New().String()
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// Process the event
	if err := am.eventProcessor.ProcessEvent(event); err != nil {
		return fmt.Errorf("failed to process audit event: %w", err)
	}

	// Calculate hash for integrity
	event.Hash = am.calculateEventHash(event)
	event.PreviousHash = am.auditStore.GetLastHash()

	// Store the event
	if err := am.auditStore.StoreEvent(event); err != nil {
		return fmt.Errorf("failed to store audit event: %w", err)
	}

	// Check compliance
	if err := am.complianceEngine.CheckCompliance(event); err != nil {
		am.logger.Warn(ctx, "Compliance check failed", map[string]interface{}{
			"event_id": event.EventID,
			"error":    err.Error(),
		})
	}

	return nil
}

// LogAuthenticationEvent logs an authentication event
func (am *AuditManager) LogAuthenticationEvent(ctx context.Context, userID *uuid.UUID, action string, result AuditResult, ipAddress, userAgent string, details map[string]interface{}) error {
	event := &AuditEvent{
		EventType:     AuditEventTypeAuthentication,
		Category:      AuditCategorySecurity,
		Severity:      am.getSeverityForResult(result),
		UserID:        userID,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		Action:        action,
		Result:        result,
		Details:       details,
		ComplianceTag: "AUTH",
	}

	return am.LogEvent(ctx, event)
}

// LogDataAccessEvent logs a data access event
func (am *AuditManager) LogDataAccessEvent(ctx context.Context, userID *uuid.UUID, resource, action string, result AuditResult, details map[string]interface{}) error {
	event := &AuditEvent{
		EventType:     AuditEventTypeDataAccess,
		Category:      AuditCategoryBusiness,
		Severity:      am.getSeverityForDataAccess(action),
		UserID:        userID,
		Resource:      resource,
		Action:        action,
		Result:        result,
		Details:       details,
		ComplianceTag: "DATA",
	}

	return am.LogEvent(ctx, event)
}

// LogTradingEvent logs a trading event
func (am *AuditManager) LogTradingEvent(ctx context.Context, userID *uuid.UUID, action string, result AuditResult, details map[string]interface{}) error {
	event := &AuditEvent{
		EventType:     AuditEventTypeTrading,
		Category:      AuditCategoryBusiness,
		Severity:      AuditSeverityHigh, // Trading events are always high severity
		UserID:        userID,
		Action:        action,
		Result:        result,
		Details:       details,
		ComplianceTag: "TRADE",
	}

	return am.LogEvent(ctx, event)
}

// LogSecurityEvent logs a security event
func (am *AuditManager) LogSecurityEvent(ctx context.Context, eventType AuditEventType, action string, severity AuditSeverity, details map[string]interface{}) error {
	event := &AuditEvent{
		EventType:     eventType,
		Category:      AuditCategorySecurity,
		Severity:      severity,
		Action:        action,
		Result:        AuditResultSuccess, // Security events are informational
		Details:       details,
		ComplianceTag: "SEC",
	}

	return am.LogEvent(ctx, event)
}

// GetAuditEvents retrieves audit events with filtering
func (am *AuditManager) GetAuditEvents(ctx context.Context, filter AuditEventFilter) ([]AuditEvent, error) {
	return am.auditStore.GetEvents(filter)
}

// GetComplianceReport generates a compliance report
func (am *AuditManager) GetComplianceReport(ctx context.Context, standard string, startTime, endTime time.Time) (*ComplianceReport, error) {
	return am.complianceEngine.GenerateReport(standard, startTime, endTime)
}

// VerifyIntegrity verifies audit log integrity
func (am *AuditManager) VerifyIntegrity(ctx context.Context) (bool, error) {
	if !am.config.EnableIntegrityCheck {
		return true, nil
	}

	return am.auditStore.VerifyIntegrity()
}

// calculateEventHash calculates hash for event integrity
func (am *AuditManager) calculateEventHash(event *AuditEvent) string {
	// Create a copy without hash fields
	eventCopy := *event
	eventCopy.Hash = ""
	eventCopy.PreviousHash = ""

	// Serialize to JSON
	data, _ := json.Marshal(eventCopy)

	// Calculate SHA-256 hash
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

// getSeverityForResult determines severity based on result
func (am *AuditManager) getSeverityForResult(result AuditResult) AuditSeverity {
	switch result {
	case AuditResultSuccess:
		return AuditSeverityLow
	case AuditResultFailure:
		return AuditSeverityMedium
	case AuditResultDenied:
		return AuditSeverityHigh
	case AuditResultError:
		return AuditSeverityCritical
	default:
		return AuditSeverityMedium
	}
}

// getSeverityForDataAccess determines severity for data access
func (am *AuditManager) getSeverityForDataAccess(action string) AuditSeverity {
	highRiskActions := map[string]bool{
		"delete": true,
		"export": true,
		"modify": true,
	}

	if highRiskActions[action] {
		return AuditSeverityHigh
	}
	return AuditSeverityMedium
}

// AuditEventFilter defines filters for audit events
type AuditEventFilter struct {
	StartTime  *time.Time       `json:"start_time,omitempty"`
	EndTime    *time.Time       `json:"end_time,omitempty"`
	EventTypes []AuditEventType `json:"event_types,omitempty"`
	Categories []AuditCategory  `json:"categories,omitempty"`
	Severities []AuditSeverity  `json:"severities,omitempty"`
	UserID     *uuid.UUID       `json:"user_id,omitempty"`
	IPAddress  string           `json:"ip_address,omitempty"`
	Resource   string           `json:"resource,omitempty"`
	Action     string           `json:"action,omitempty"`
	Result     *AuditResult     `json:"result,omitempty"`
	Limit      int              `json:"limit,omitempty"`
	Offset     int              `json:"offset,omitempty"`
}

// ComplianceReport represents a compliance report
type ComplianceReport struct {
	Standard        string                  `json:"standard"`
	Period          string                  `json:"period"`
	TotalEvents     int                     `json:"total_events"`
	Violations      []ComplianceViolation   `json:"violations"`
	ComplianceScore float64                 `json:"compliance_score"`
	Requirements    []RequirementCompliance `json:"requirements"`
	GeneratedAt     time.Time               `json:"generated_at"`
}

// RequirementCompliance represents compliance for a specific requirement
type RequirementCompliance struct {
	RequirementID string  `json:"requirement_id"`
	Description   string  `json:"description"`
	Compliant     bool    `json:"compliant"`
	EventCount    int     `json:"event_count"`
	Score         float64 `json:"score"`
}

// NewAuditStore creates a new audit store
func NewAuditStore(logger *observability.Logger, config *AuditConfig) *AuditStore {
	return &AuditStore{
		logger:     logger,
		config:     config,
		events:     make([]AuditEvent, 0),
		eventIndex: make(map[string]*AuditEvent),
	}
}

// StoreEvent stores an audit event
func (as *AuditStore) StoreEvent(event *AuditEvent) error {
	as.mu.Lock()
	defer as.mu.Unlock()

	// Store event
	as.events = append(as.events, *event)
	as.eventIndex[event.EventID] = &as.events[len(as.events)-1]
	as.lastHash = event.Hash

	// Update total size
	eventSize := int64(len(fmt.Sprintf("%+v", event)))
	as.totalSize += eventSize

	// Check if archiving is needed
	if as.totalSize > as.config.ArchiveThreshold {
		go as.archiveOldEvents()
	}

	return nil
}

// GetEvents retrieves events with filtering
func (as *AuditStore) GetEvents(filter AuditEventFilter) ([]AuditEvent, error) {
	as.mu.RLock()
	defer as.mu.RUnlock()

	var filteredEvents []AuditEvent

	for _, event := range as.events {
		if as.matchesFilter(&event, filter) {
			filteredEvents = append(filteredEvents, event)
		}
	}

	// Apply limit and offset
	if filter.Offset > 0 && filter.Offset < len(filteredEvents) {
		filteredEvents = filteredEvents[filter.Offset:]
	}
	if filter.Limit > 0 && filter.Limit < len(filteredEvents) {
		filteredEvents = filteredEvents[:filter.Limit]
	}

	return filteredEvents, nil
}

// GetLastHash returns the last event hash
func (as *AuditStore) GetLastHash() string {
	as.mu.RLock()
	defer as.mu.RUnlock()
	return as.lastHash
}

// VerifyIntegrity verifies the integrity of the audit log
func (as *AuditStore) VerifyIntegrity() (bool, error) {
	as.mu.RLock()
	defer as.mu.RUnlock()

	for i, event := range as.events {
		// Verify hash
		expectedHash := as.calculateHash(&event)
		if event.Hash != expectedHash {
			return false, fmt.Errorf("hash mismatch for event %s", event.EventID)
		}

		// Verify chain
		if i > 0 {
			previousEvent := as.events[i-1]
			if event.PreviousHash != previousEvent.Hash {
				return false, fmt.Errorf("chain broken at event %s", event.EventID)
			}
		}
	}

	return true, nil
}

// matchesFilter checks if an event matches the filter
func (as *AuditStore) matchesFilter(event *AuditEvent, filter AuditEventFilter) bool {
	// Time range filter
	if filter.StartTime != nil && event.Timestamp.Before(*filter.StartTime) {
		return false
	}
	if filter.EndTime != nil && event.Timestamp.After(*filter.EndTime) {
		return false
	}

	// Event type filter
	if len(filter.EventTypes) > 0 {
		found := false
		for _, eventType := range filter.EventTypes {
			if event.EventType == eventType {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// User ID filter
	if filter.UserID != nil && (event.UserID == nil || *event.UserID != *filter.UserID) {
		return false
	}

	// IP address filter
	if filter.IPAddress != "" && event.IPAddress != filter.IPAddress {
		return false
	}

	// Resource filter
	if filter.Resource != "" && event.Resource != filter.Resource {
		return false
	}

	// Action filter
	if filter.Action != "" && event.Action != filter.Action {
		return false
	}

	// Result filter
	if filter.Result != nil && event.Result != *filter.Result {
		return false
	}

	return true
}

// calculateHash calculates the hash for an event
func (as *AuditStore) calculateHash(event *AuditEvent) string {
	eventCopy := *event
	eventCopy.Hash = ""
	eventCopy.PreviousHash = ""

	data, _ := json.Marshal(eventCopy)
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

// archiveOldEvents archives old events
func (as *AuditStore) archiveOldEvents() {
	as.mu.Lock()
	defer as.mu.Unlock()

	// Archive events older than retention period
	cutoffTime := time.Now().Add(-as.config.RetentionPeriod)
	var activeEvents []AuditEvent

	for _, event := range as.events {
		if event.Timestamp.After(cutoffTime) {
			activeEvents = append(activeEvents, event)
		}
	}

	as.events = activeEvents
	as.eventIndex = make(map[string]*AuditEvent)

	// Rebuild index
	for i := range as.events {
		as.eventIndex[as.events[i].EventID] = &as.events[i]
	}
}

// NewAuditEventProcessor creates a new audit event processor
func NewAuditEventProcessor(logger *observability.Logger, config *AuditConfig) *AuditEventProcessor {
	return &AuditEventProcessor{
		logger:    logger,
		config:    config,
		filters:   make([]AuditFilter, 0),
		enrichers: make([]AuditEnricher, 0),
	}
}

// ProcessEvent processes an audit event
func (aep *AuditEventProcessor) ProcessEvent(event *AuditEvent) error {
	// Apply filters
	for _, filter := range aep.filters {
		if !filter.ShouldAudit(event) {
			return nil // Event filtered out
		}
	}

	// Apply enrichers
	for _, enricher := range aep.enrichers {
		if err := enricher.Enrich(event); err != nil {
			aep.logger.Warn(nil, "Failed to enrich audit event", map[string]interface{}{
				"event_id": event.EventID,
				"error":    err.Error(),
			})
		}
	}

	return nil
}

// NewComplianceEngine creates a new compliance engine
func NewComplianceEngine(logger *observability.Logger, config *AuditConfig) *ComplianceEngine {
	return &ComplianceEngine{
		logger:     logger,
		config:     config,
		standards:  make(map[string]*ComplianceStandard),
		violations: make([]ComplianceViolation, 0),
	}
}

// InitializeStandards initializes compliance standards
func (ce *ComplianceEngine) InitializeStandards() error {
	// Initialize SOX compliance
	if ce.containsStandard("SOX") {
		ce.standards["SOX"] = &ComplianceStandard{
			Name: "Sarbanes-Oxley Act",
			Requirements: []ComplianceRequirement{
				{
					ID:          "SOX-404",
					Description: "Internal controls over financial reporting",
					EventTypes:  []string{"financial", "trading"},
					Mandatory:   true,
				},
			},
			Enabled: true,
		}
	}

	// Initialize GDPR compliance
	if ce.containsStandard("GDPR") {
		ce.standards["GDPR"] = &ComplianceStandard{
			Name: "General Data Protection Regulation",
			Requirements: []ComplianceRequirement{
				{
					ID:          "GDPR-32",
					Description: "Security of processing",
					EventTypes:  []string{"data_access", "data_modification"},
					Mandatory:   true,
				},
			},
			Enabled: true,
		}
	}

	return nil
}

// CheckCompliance checks event compliance
func (ce *ComplianceEngine) CheckCompliance(event *AuditEvent) error {
	ce.mu.Lock()
	defer ce.mu.Unlock()

	for _, standard := range ce.standards {
		if !standard.Enabled {
			continue
		}

		for _, requirement := range standard.Requirements {
			if ce.eventMatchesRequirement(event, &requirement) {
				// Event matches requirement - check compliance
				if !ce.isCompliant(event, &requirement) {
					violation := ComplianceViolation{
						ViolationID: uuid.New().String(),
						Standard:    standard.Name,
						Requirement: requirement.ID,
						EventID:     event.EventID,
						Timestamp:   time.Now(),
						Severity:    string(event.Severity),
						Description: fmt.Sprintf("Event does not meet %s requirement", requirement.ID),
					}
					ce.violations = append(ce.violations, violation)
				}
			}
		}
	}

	return nil
}

// GenerateReport generates a compliance report
func (ce *ComplianceEngine) GenerateReport(standard string, startTime, endTime time.Time) (*ComplianceReport, error) {
	ce.mu.RLock()
	defer ce.mu.RUnlock()

	report := &ComplianceReport{
		Standard:    standard,
		Period:      fmt.Sprintf("%s to %s", startTime.Format("2006-01-02"), endTime.Format("2006-01-02")),
		GeneratedAt: time.Now(),
	}

	// Count violations in period
	var periodViolations []ComplianceViolation
	for _, violation := range ce.violations {
		if violation.Timestamp.After(startTime) && violation.Timestamp.Before(endTime) {
			periodViolations = append(periodViolations, violation)
		}
	}

	report.Violations = periodViolations
	report.ComplianceScore = ce.calculateComplianceScore(standard, periodViolations)

	return report, nil
}

// Helper methods
func (ce *ComplianceEngine) containsStandard(standard string) bool {
	for _, s := range ce.config.ComplianceStandards {
		if s == standard {
			return true
		}
	}
	return false
}

func (ce *ComplianceEngine) eventMatchesRequirement(event *AuditEvent, requirement *ComplianceRequirement) bool {
	for _, eventType := range requirement.EventTypes {
		if string(event.EventType) == eventType {
			return true
		}
	}
	return false
}

func (ce *ComplianceEngine) isCompliant(event *AuditEvent, requirement *ComplianceRequirement) bool {
	// Basic compliance check - can be extended
	return event.Result == AuditResultSuccess
}

func (ce *ComplianceEngine) calculateComplianceScore(standard string, violations []ComplianceViolation) float64 {
	if len(violations) == 0 {
		return 100.0
	}
	// Simple calculation - can be made more sophisticated
	return math.Max(0, 100.0-float64(len(violations))*5.0)
}
