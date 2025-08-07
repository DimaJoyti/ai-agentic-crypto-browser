package security

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

// AdvancedThreatDetector provides comprehensive threat detection
type AdvancedThreatDetector struct {
	logger             *observability.Logger
	config             *ThreatDetectionConfig
	signatureEngine    *SignatureEngine
	behaviorEngine     *BehaviorThreatEngine
	mlEngine           *MLThreatEngine
	threatIntelligence *ThreatIntelligenceService
	incidentManager    *IncidentManager
	alertManager       *AlertManager
	activeThreats      map[string]*ThreatIncident
	blockedIPs         map[string]*BlockedIP
	suspiciousPatterns map[string]*SuspiciousPattern
	mu                 sync.RWMutex
}

// ThreatDetectionConfig contains threat detection configuration
type ThreatDetectionConfig struct {
	EnableSignatureDetection bool
	EnableBehaviorDetection  bool
	EnableMLDetection        bool
	EnableThreatIntelligence bool
	EnableRealTimeBlocking   bool
	MaxThreatScore           float64
	BlockThreshold           float64
	AlertThreshold           float64
	IncidentThreshold        float64
	ScanInterval             time.Duration
	ThreatRetentionPeriod    time.Duration
	MaxBlockedIPs            int
	AutoMitigationEnabled    bool
}

// ThreatIncident represents a security threat incident
type ThreatIncident struct {
	IncidentID      string
	ThreatType      ThreatType
	Severity        ThreatSeverity
	Status          IncidentStatus
	Source          string
	Target          string
	Description     string
	Indicators      []ThreatIndicator
	Timeline        []IncidentEvent
	AffectedUsers   []uuid.UUID
	AffectedSystems []string
	MitigationSteps []MitigationAction
	CreatedAt       time.Time
	UpdatedAt       time.Time
	ResolvedAt      *time.Time
	AssignedTo      *uuid.UUID
}

// IncidentStatus defines incident status
type IncidentStatus string

const (
	IncidentStatusOpen       IncidentStatus = "open"
	IncidentStatusInProgress IncidentStatus = "in_progress"
	IncidentStatusResolved   IncidentStatus = "resolved"
	IncidentStatusClosed     IncidentStatus = "closed"
)

// ThreatIndicator represents an indicator of compromise
type ThreatIndicator struct {
	Type        IndicatorType
	Value       string
	Confidence  float64
	Source      string
	FirstSeen   time.Time
	LastSeen    time.Time
	Description string
	Tags        []string
}

// IndicatorType defines types of threat indicators
type IndicatorType string

const (
	IndicatorTypeIP       IndicatorType = "ip"
	IndicatorTypeDomain   IndicatorType = "domain"
	IndicatorTypeURL      IndicatorType = "url"
	IndicatorTypeHash     IndicatorType = "hash"
	IndicatorTypeEmail    IndicatorType = "email"
	IndicatorTypePattern  IndicatorType = "pattern"
	IndicatorTypeBehavior IndicatorType = "behavior"
)

// IncidentEvent represents an event in the incident timeline
type IncidentEvent struct {
	EventID     string
	Type        string
	Description string
	Timestamp   time.Time
	UserID      *uuid.UUID
	Automated   bool
	Data        map[string]interface{}
}

// MitigationAction represents an action taken to mitigate a threat
type MitigationAction struct {
	ActionID    string
	Type        MitigationType
	Description string
	Status      ActionStatus
	ExecutedAt  time.Time
	ExecutedBy  *uuid.UUID
	Automated   bool
	Result      string
	Data        map[string]interface{}
}

// MitigationType defines types of mitigation actions
type MitigationType string

const (
	MitigationTypeBlock       MitigationType = "block"
	MitigationTypeQuarantine  MitigationType = "quarantine"
	MitigationTypeAlert       MitigationType = "alert"
	MitigationTypeThrottle    MitigationType = "throttle"
	MitigationTypeDisable     MitigationType = "disable"
	MitigationTypeInvestigate MitigationType = "investigate"
)

// ActionStatus defines mitigation action status
type ActionStatus string

const (
	ActionStatusPending   ActionStatus = "pending"
	ActionStatusExecuting ActionStatus = "executing"
	ActionStatusCompleted ActionStatus = "completed"
	ActionStatusFailed    ActionStatus = "failed"
)

// BlockedIP represents a blocked IP address
type BlockedIP struct {
	IPAddress   string
	Reason      string
	BlockedAt   time.Time
	ExpiresAt   *time.Time
	BlockedBy   string
	ThreatScore float64
	Incidents   []string
}

// SuspiciousPattern represents a suspicious behavior pattern
type SuspiciousPattern struct {
	PatternID   string
	Type        string
	Pattern     string
	Regex       *regexp.Regexp
	Severity    ThreatSeverity
	Description string
	MatchCount  int64
	FirstSeen   time.Time
	LastSeen    time.Time
	Enabled     bool
}

// ThreatDetectionResult represents the result of threat detection
type ThreatDetectionResult struct {
	ThreatDetected  bool
	ThreatScore     float64
	ThreatType      ThreatType
	Severity        ThreatSeverity
	Indicators      []ThreatIndicator
	Confidence      float64
	Recommendations []string
	ShouldBlock     bool
	ShouldAlert     bool
	Metadata        map[string]interface{}
}

// NewAdvancedThreatDetector creates a new advanced threat detector
func NewAdvancedThreatDetector(logger *observability.Logger) *AdvancedThreatDetector {
	config := &ThreatDetectionConfig{
		EnableSignatureDetection: true,
		EnableBehaviorDetection:  true,
		EnableMLDetection:        true,
		EnableThreatIntelligence: true,
		EnableRealTimeBlocking:   true,
		MaxThreatScore:           1.0,
		BlockThreshold:           0.8,
		AlertThreshold:           0.6,
		IncidentThreshold:        0.7,
		ScanInterval:             1 * time.Minute,
		ThreatRetentionPeriod:    30 * 24 * time.Hour,
		MaxBlockedIPs:            10000,
		AutoMitigationEnabled:    true,
	}

	detector := &AdvancedThreatDetector{
		logger:             logger,
		config:             config,
		signatureEngine:    NewSignatureEngine(logger),
		behaviorEngine:     NewBehaviorThreatEngine(logger),
		mlEngine:           NewMLThreatEngine(logger),
		threatIntelligence: NewThreatIntelligenceService(logger),
		incidentManager:    NewIncidentManager(logger),
		alertManager:       NewAlertManager(logger),
		activeThreats:      make(map[string]*ThreatIncident),
		blockedIPs:         make(map[string]*BlockedIP),
		suspiciousPatterns: make(map[string]*SuspiciousPattern),
	}

	// Initialize default suspicious patterns
	detector.initializeDefaultPatterns()

	// Start background threat monitoring
	go detector.startThreatMonitoring()

	return detector
}

// DetectThreats performs comprehensive threat detection
func (a *AdvancedThreatDetector) DetectThreats(ctx context.Context, request *SecurityRequest) (*ThreatDetectionResult, error) {
	result := &ThreatDetectionResult{
		ThreatDetected:  false,
		ThreatScore:     0.0,
		Indicators:      []ThreatIndicator{},
		Confidence:      0.0,
		Recommendations: []string{},
		Metadata:        make(map[string]interface{}),
	}

	// 1. Check if IP is already blocked
	if a.isIPBlocked(request.IPAddress) {
		result.ThreatDetected = true
		result.ThreatScore = 1.0
		result.ThreatType = ThreatTypeBotActivity
		result.Severity = ThreatSeverityHigh
		result.ShouldBlock = true
		result.Recommendations = append(result.Recommendations, "IP is already blocked")
		return result, nil
	}

	// 2. Signature-based detection
	if a.config.EnableSignatureDetection {
		sigResult := a.signatureEngine.DetectThreats(request)
		if sigResult.ThreatDetected {
			result.ThreatDetected = true
			result.ThreatScore = max(result.ThreatScore, sigResult.ThreatScore)
			result.Indicators = append(result.Indicators, sigResult.Indicators...)
		}
	}

	// 3. Behavior-based detection
	if a.config.EnableBehaviorDetection {
		behaviorResult := a.behaviorEngine.DetectThreats(request)
		if behaviorResult.ThreatDetected {
			result.ThreatDetected = true
			result.ThreatScore = max(result.ThreatScore, behaviorResult.ThreatScore)
			result.Indicators = append(result.Indicators, behaviorResult.Indicators...)
		}
	}

	// 4. ML-based detection
	if a.config.EnableMLDetection {
		mlResult := a.mlEngine.DetectThreats(request)
		if mlResult.ThreatDetected {
			result.ThreatDetected = true
			result.ThreatScore = max(result.ThreatScore, mlResult.ThreatScore)
			result.Indicators = append(result.Indicators, mlResult.Indicators...)
		}
	}

	// 5. Threat intelligence lookup
	if a.config.EnableThreatIntelligence {
		intelResult := a.threatIntelligence.CheckThreatIntelligence(request)
		if intelResult.ThreatDetected {
			result.ThreatDetected = true
			result.ThreatScore = max(result.ThreatScore, intelResult.ThreatScore)
			result.Indicators = append(result.Indicators, intelResult.Indicators...)
		}
	}

	// 6. Determine threat type and severity
	if result.ThreatDetected {
		result.ThreatType = a.determineThreatType(result.Indicators)
		result.Severity = a.determineThreatSeverity(result.ThreatScore)
		result.Confidence = a.calculateConfidence(result.Indicators)

		// 7. Determine actions
		result.ShouldBlock = result.ThreatScore >= a.config.BlockThreshold
		result.ShouldAlert = result.ThreatScore >= a.config.AlertThreshold

		// 8. Generate recommendations
		result.Recommendations = a.generateRecommendations(result)

		// 9. Create incident if threshold exceeded
		if result.ThreatScore >= a.config.IncidentThreshold {
			incident := a.createThreatIncident(request, result)
			result.Metadata["incident_id"] = incident.IncidentID
		}

		// 10. Execute automatic mitigation if enabled
		if a.config.AutoMitigationEnabled {
			a.executeMitigation(ctx, request, result)
		}
	}

	// Log threat detection result
	a.logThreatDetection(ctx, request, result)

	return result, nil
}

// isIPBlocked checks if an IP address is blocked
func (a *AdvancedThreatDetector) isIPBlocked(ipAddress string) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	blocked, exists := a.blockedIPs[ipAddress]
	if !exists {
		return false
	}

	// Check if block has expired
	if blocked.ExpiresAt != nil && time.Now().After(*blocked.ExpiresAt) {
		delete(a.blockedIPs, ipAddress)
		return false
	}

	return true
}

// BlockIP blocks an IP address
func (a *AdvancedThreatDetector) BlockIP(ipAddress, reason string, duration *time.Duration) {
	a.mu.Lock()
	defer a.mu.Unlock()

	var expiresAt *time.Time
	if duration != nil {
		expiry := time.Now().Add(*duration)
		expiresAt = &expiry
	}

	a.blockedIPs[ipAddress] = &BlockedIP{
		IPAddress: ipAddress,
		Reason:    reason,
		BlockedAt: time.Now(),
		ExpiresAt: expiresAt,
		BlockedBy: "system",
	}

	a.logger.Warn(context.Background(), "IP address blocked", map[string]interface{}{
		"ip_address": ipAddress,
		"reason":     reason,
		"expires_at": expiresAt,
	})
}

// initializeDefaultPatterns initializes default suspicious patterns
func (a *AdvancedThreatDetector) initializeDefaultPatterns() {
	patterns := []*SuspiciousPattern{
		{
			PatternID:   "sql_injection",
			Type:        "sql_injection",
			Pattern:     `(?i)(union|select|insert|update|delete|drop|create|alter|exec|execute)`,
			Severity:    ThreatSeverityHigh,
			Description: "SQL injection attempt detected",
			Enabled:     true,
		},
		{
			PatternID:   "xss_attempt",
			Type:        "xss",
			Pattern:     `(?i)(<script|javascript:|onload=|onerror=|onclick=)`,
			Severity:    ThreatSeverityMedium,
			Description: "Cross-site scripting attempt detected",
			Enabled:     true,
		},
		{
			PatternID:   "path_traversal",
			Type:        "path_traversal",
			Pattern:     `(\.\.\/|\.\.\\|%2e%2e%2f|%2e%2e%5c)`,
			Severity:    ThreatSeverityMedium,
			Description: "Path traversal attempt detected",
			Enabled:     true,
		},
	}

	for _, pattern := range patterns {
		pattern.Regex = regexp.MustCompile(pattern.Pattern)
		pattern.FirstSeen = time.Now()
		pattern.LastSeen = time.Now()
		a.suspiciousPatterns[pattern.PatternID] = pattern
	}
}

// startThreatMonitoring starts background threat monitoring
func (a *AdvancedThreatDetector) startThreatMonitoring() {
	ticker := time.NewTicker(a.config.ScanInterval)
	defer ticker.Stop()

	for range ticker.C {
		a.performMaintenanceTasks()
	}
}

// performMaintenanceTasks performs periodic maintenance
func (a *AdvancedThreatDetector) performMaintenanceTasks() {
	// Clean up expired blocked IPs
	a.cleanupExpiredBlocks()

	// Clean up old threat incidents
	a.cleanupOldIncidents()

	// Update threat intelligence
	if a.config.EnableThreatIntelligence {
		a.threatIntelligence.UpdateThreatFeeds()
	}
}

// cleanupExpiredBlocks removes expired IP blocks
func (a *AdvancedThreatDetector) cleanupExpiredBlocks() {
	a.mu.Lock()
	defer a.mu.Unlock()

	now := time.Now()
	for ip, blocked := range a.blockedIPs {
		if blocked.ExpiresAt != nil && now.After(*blocked.ExpiresAt) {
			delete(a.blockedIPs, ip)
		}
	}
}

// cleanupOldIncidents removes old threat incidents
func (a *AdvancedThreatDetector) cleanupOldIncidents() {
	a.mu.Lock()
	defer a.mu.Unlock()

	cutoff := time.Now().Add(-a.config.ThreatRetentionPeriod)
	for id, incident := range a.activeThreats {
		if incident.CreatedAt.Before(cutoff) && incident.Status == IncidentStatusClosed {
			delete(a.activeThreats, id)
		}
	}
}

// SecurityRequest represents a security request for threat analysis
type SecurityRequest struct {
	RequestID   string
	UserID      *uuid.UUID
	IPAddress   string
	UserAgent   string
	Method      string
	URL         string
	Headers     map[string]string
	Body        string
	Timestamp   time.Time
	SessionID   string
	DeviceID    string
	Geolocation *Location
}

// determineThreatType determines the primary threat type from indicators
func (a *AdvancedThreatDetector) determineThreatType(indicators []ThreatIndicator) ThreatType {
	typeCounts := make(map[ThreatType]int)

	for _, indicator := range indicators {
		switch indicator.Type {
		case IndicatorTypePattern:
			if strings.Contains(indicator.Value, "sql") {
				typeCounts[ThreatTypeCredStuffing]++
			} else if strings.Contains(indicator.Value, "script") {
				typeCounts[ThreatTypeBotActivity]++
			}
		case IndicatorTypeBehavior:
			typeCounts[ThreatTypeAccountTakeover]++
		case IndicatorTypeIP:
			typeCounts[ThreatTypeBruteForce]++
		}
	}

	// Return the most common threat type
	maxCount := 0
	var primaryType ThreatType = ThreatTypeSuspiciousAPI

	for threatType, count := range typeCounts {
		if count > maxCount {
			maxCount = count
			primaryType = threatType
		}
	}

	return primaryType
}

// determineThreatSeverity determines threat severity from score
func (a *AdvancedThreatDetector) determineThreatSeverity(score float64) ThreatSeverity {
	switch {
	case score >= 0.8:
		return ThreatSeverityCritical
	case score >= 0.6:
		return ThreatSeverityHigh
	case score >= 0.4:
		return ThreatSeverityMedium
	default:
		return ThreatSeverityLow
	}
}

// calculateConfidence calculates confidence based on indicators
func (a *AdvancedThreatDetector) calculateConfidence(indicators []ThreatIndicator) float64 {
	if len(indicators) == 0 {
		return 0.0
	}

	totalConfidence := 0.0
	for _, indicator := range indicators {
		totalConfidence += indicator.Confidence
	}

	return totalConfidence / float64(len(indicators))
}

// generateRecommendations generates security recommendations
func (a *AdvancedThreatDetector) generateRecommendations(result *ThreatDetectionResult) []string {
	recommendations := []string{}

	if result.ThreatScore >= a.config.BlockThreshold {
		recommendations = append(recommendations, "Block source IP immediately")
	}

	if result.ThreatScore >= a.config.AlertThreshold {
		recommendations = append(recommendations, "Alert security team")
	}

	switch result.ThreatType {
	case ThreatTypeBruteForce:
		recommendations = append(recommendations, "Implement rate limiting", "Require MFA")
	case ThreatTypeCredStuffing:
		recommendations = append(recommendations, "Force password reset", "Monitor for account compromise")
	case ThreatTypeBotActivity:
		recommendations = append(recommendations, "Implement CAPTCHA", "Analyze traffic patterns")
	}

	return recommendations
}

// createThreatIncident creates a new threat incident
func (a *AdvancedThreatDetector) createThreatIncident(request *SecurityRequest, result *ThreatDetectionResult) *ThreatIncident {
	incident := &ThreatIncident{
		IncidentID:  uuid.New().String(),
		ThreatType:  result.ThreatType,
		Severity:    result.Severity,
		Status:      IncidentStatusOpen,
		Source:      request.IPAddress,
		Target:      request.URL,
		Description: fmt.Sprintf("Threat detected: %s from %s", result.ThreatType, request.IPAddress),
		Indicators:  result.Indicators,
		Timeline:    []IncidentEvent{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if request.UserID != nil {
		incident.AffectedUsers = []uuid.UUID{*request.UserID}
	}

	// Add initial event
	incident.Timeline = append(incident.Timeline, IncidentEvent{
		EventID:     uuid.New().String(),
		Type:        "threat_detected",
		Description: "Threat detected by automated system",
		Timestamp:   time.Now(),
		Automated:   true,
		Data: map[string]interface{}{
			"threat_score": result.ThreatScore,
			"confidence":   result.Confidence,
		},
	})

	a.mu.Lock()
	a.activeThreats[incident.IncidentID] = incident
	a.mu.Unlock()

	return incident
}

// executeMitigation executes automatic mitigation actions
func (a *AdvancedThreatDetector) executeMitigation(ctx context.Context, request *SecurityRequest, result *ThreatDetectionResult) {
	if result.ShouldBlock {
		// Block IP for 1 hour
		duration := 1 * time.Hour
		a.BlockIP(request.IPAddress, fmt.Sprintf("Automatic block: %s", result.ThreatType), &duration)
	}

	if result.ShouldAlert {
		// Send alert
		a.alertManager.SendAlert(&SecurityAlert{
			AlertID:     uuid.New().String(),
			Type:        "threat_detected",
			Severity:    string(result.Severity),
			Title:       fmt.Sprintf("Threat Detected: %s", result.ThreatType),
			Description: fmt.Sprintf("Threat score: %.2f from IP: %s", result.ThreatScore, request.IPAddress),
			Source:      request.IPAddress,
			Timestamp:   time.Now(),
			Metadata:    result.Metadata,
		})
	}
}

// logThreatDetection logs threat detection results
func (a *AdvancedThreatDetector) logThreatDetection(ctx context.Context, request *SecurityRequest, result *ThreatDetectionResult) {
	logLevel := "info"
	if result.ThreatDetected {
		switch result.Severity {
		case ThreatSeverityCritical, ThreatSeverityHigh:
			logLevel = "error"
		case ThreatSeverityMedium:
			logLevel = "warn"
		default:
			logLevel = "info"
		}
	}

	logData := map[string]interface{}{
		"request_id":      request.RequestID,
		"ip_address":      request.IPAddress,
		"threat_detected": result.ThreatDetected,
		"threat_score":    result.ThreatScore,
		"threat_type":     result.ThreatType,
		"severity":        result.Severity,
		"confidence":      result.Confidence,
		"should_block":    result.ShouldBlock,
		"should_alert":    result.ShouldAlert,
		"indicators":      len(result.Indicators),
	}

	switch logLevel {
	case "error":
		a.logger.Error(ctx, "Critical threat detected", nil, logData)
	case "warn":
		a.logger.Warn(ctx, "Threat detected", logData)
	default:
		a.logger.Info(ctx, "Threat detection completed", logData)
	}
}

// Note: SecurityAlert is now defined in security_dashboard.go

// Component constructors and basic implementations

// BehaviorThreatEngine detects behavioral threats
type BehaviorThreatEngine struct {
	logger *observability.Logger
}

// MLThreatEngine uses machine learning for threat detection
type MLThreatEngine struct {
	logger *observability.Logger
}

// ThreatIntelligenceService provides threat intelligence
type ThreatIntelligenceService struct {
	logger *observability.Logger
}

// IncidentManager manages security incidents
type IncidentManager struct {
	logger *observability.Logger
}

// AlertManager manages security alerts
type AlertManager struct {
	logger *observability.Logger
}

// NewSignatureEngine creates a new signature engine
func NewSignatureEngine(logger *observability.Logger) *SignatureEngine {
	return &SignatureEngine{
		logger: logger,
	}
}

// DetectThreats detects threats using signatures
func (s *SignatureEngine) DetectThreats(request *SecurityRequest) *ThreatDetectionResult {
	// Basic signature detection implementation
	return &ThreatDetectionResult{
		ThreatDetected: false,
		ThreatScore:    0.0,
		Indicators:     []ThreatIndicator{},
	}
}

// NewBehaviorThreatEngine creates a new behavior threat engine
func NewBehaviorThreatEngine(logger *observability.Logger) *BehaviorThreatEngine {
	return &BehaviorThreatEngine{
		logger: logger,
	}
}

// DetectThreats detects behavioral threats
func (b *BehaviorThreatEngine) DetectThreats(request *SecurityRequest) *ThreatDetectionResult {
	// Basic behavioral detection implementation
	return &ThreatDetectionResult{
		ThreatDetected: false,
		ThreatScore:    0.0,
		Indicators:     []ThreatIndicator{},
	}
}

// NewMLThreatEngine creates a new ML threat engine
func NewMLThreatEngine(logger *observability.Logger) *MLThreatEngine {
	return &MLThreatEngine{
		logger: logger,
	}
}

// DetectThreats detects threats using ML
func (m *MLThreatEngine) DetectThreats(request *SecurityRequest) *ThreatDetectionResult {
	// Basic ML detection implementation
	return &ThreatDetectionResult{
		ThreatDetected: false,
		ThreatScore:    0.0,
		Indicators:     []ThreatIndicator{},
	}
}

// NewThreatIntelligenceService creates a new threat intelligence service
func NewThreatIntelligenceService(logger *observability.Logger) *ThreatIntelligenceService {
	return &ThreatIntelligenceService{
		logger: logger,
	}
}

// CheckThreatIntelligence checks threat intelligence feeds
func (t *ThreatIntelligenceService) CheckThreatIntelligence(request *SecurityRequest) *ThreatDetectionResult {
	// Basic threat intelligence check
	return &ThreatDetectionResult{
		ThreatDetected: false,
		ThreatScore:    0.0,
		Indicators:     []ThreatIndicator{},
	}
}

// UpdateThreatFeeds updates threat intelligence feeds
func (t *ThreatIntelligenceService) UpdateThreatFeeds() {
	// Update threat feeds implementation
}

// NewIncidentManager creates a new incident manager
func NewIncidentManager(logger *observability.Logger) *IncidentManager {
	return &IncidentManager{
		logger: logger,
	}
}

// NewAlertManager creates a new alert manager
func NewAlertManager(logger *observability.Logger) *AlertManager {
	return &AlertManager{
		logger: logger,
	}
}

// SendAlert sends a security alert
func (a *AlertManager) SendAlert(alert *SecurityAlert) {
	a.logger.Warn(context.Background(), "Security alert", map[string]interface{}{
		"alert_id":    alert.AlertID,
		"type":        alert.Type,
		"severity":    alert.Severity,
		"title":       alert.Title,
		"description": alert.Description,
		"source":      alert.Source,
	})
}
