package security

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

// ZeroTrustEngine implements zero-trust security architecture
type ZeroTrustEngine struct {
	logger           *observability.Logger
	config           *ZeroTrustConfig
	deviceRegistry   *DeviceRegistry
	behaviorAnalyzer *ZeroTrustBehaviorAnalyzer
	threatDetector   *ZeroTrustThreatDetector
	policyEngine     *PolicyEngine
	sessionManager   *ZeroTrustSessionManager
	riskCalculator   *RiskCalculator
	mu               sync.RWMutex
}

// ZeroTrustConfig contains zero-trust configuration
type ZeroTrustConfig struct {
	EnableDeviceFingerprinting bool
	EnableBehaviorAnalysis     bool
	EnableThreatDetection      bool
	EnableContinuousAuth       bool
	RiskThreshold              float64
	SessionTimeout             time.Duration
	DeviceTrustDuration        time.Duration
	MaxRiskScore               float64
	RequireMFAForHighRisk      bool
	TTLCalculationStrategy     string // "linear", "exponential", "logarithmic", "stepped", "adaptive"
}

// DeviceRegistry manages trusted devices
type DeviceRegistry struct {
	devices map[string]*TrustedDevice
	logger  *observability.Logger
	mu      sync.RWMutex
}

// TrustedDevice is defined in rate_limiter.go to avoid duplication

// ZeroTrustBehaviorAnalyzer analyzes user behavior patterns
type ZeroTrustBehaviorAnalyzer struct {
	logger          *observability.Logger
	userProfiles    map[uuid.UUID]*UserBehaviorProfile
	anomalyDetector *AnomalyDetector
	baselineBuilder *BaselineBuilder
	mu              sync.RWMutex
}

// UserBehaviorProfile contains user behavior patterns
type UserBehaviorProfile struct {
	UserID              uuid.UUID
	TypicalLoginTimes   []time.Time
	TypicalLocations    []Location
	TypicalDevices      []string
	TypicalActions      map[string]int
	SessionDurations    []time.Duration
	APIUsagePatterns    map[string]int
	RiskScore           float64
	LastUpdated         time.Time
	AnomalyCount        int
	BaselineEstablished bool
}

// Location represents a geographical location
type Location struct {
	Country   string
	Region    string
	City      string
	Latitude  float64
	Longitude float64
	IPRange   string
}

// ZeroTrustThreatDetector detects security threats
type ZeroTrustThreatDetector struct {
	logger          *observability.Logger
	threatIntel     *ThreatIntelligence
	signatureEngine *SignatureEngine
	mlDetector      *MLThreatDetector
	activeThreats   map[string]*ActiveThreat
	mu              sync.RWMutex
}

// ActiveThreat represents an active security threat
type ActiveThreat struct {
	ThreatID    string
	Type        ThreatType
	Severity    ThreatSeverity
	Source      string
	Target      string
	Description string
	Indicators  []string
	FirstSeen   time.Time
	LastSeen    time.Time
	EventCount  int
	Mitigated   bool
	MitigatedAt *time.Time
}

// ThreatType defines types of security threats
type ThreatType string

const (
	ThreatTypeBruteForce          ThreatType = "brute_force"
	ThreatTypeCredStuffing        ThreatType = "credential_stuffing"
	ThreatTypeAccountTakeover     ThreatType = "account_takeover"
	ThreatTypeBotActivity         ThreatType = "bot_activity"
	ThreatTypeDataExfiltration    ThreatType = "data_exfiltration"
	ThreatTypePrivilegeEscalation ThreatType = "privilege_escalation"
	ThreatTypeSuspiciousAPI       ThreatType = "suspicious_api"
)

// ThreatSeverity defines threat severity levels
type ThreatSeverity string

const (
	ThreatSeverityLow      ThreatSeverity = "low"
	ThreatSeverityMedium   ThreatSeverity = "medium"
	ThreatSeverityHigh     ThreatSeverity = "high"
	ThreatSeverityCritical ThreatSeverity = "critical"
)

// SecurityEvent represents a security event
type SecurityEvent struct {
	EventID    string
	Type       string
	Severity   string
	UserID     *uuid.UUID
	DeviceID   string
	IPAddress  string
	UserAgent  string
	Resource   string
	Action     string
	Timestamp  time.Time
	RiskScore  float64
	Indicators []string
	Context    map[string]interface{}
	Blocked    bool
	Mitigated  bool
}

// AccessRequest represents a request for access evaluation
type AccessRequest struct {
	UserID    *uuid.UUID
	DeviceID  string
	IPAddress string
	UserAgent string
	Resource  string
	Action    string
	Timestamp time.Time
	Context   map[string]interface{}
}

// AccessDecision represents the result of access evaluation
type AccessDecision struct {
	Allowed      bool
	RiskScore    float64
	DeviceTrust  float64
	BehaviorRisk float64
	ThreatLevel  float64
	RequiresMFA  bool
	SessionTTL   time.Duration
	Reason       string
	Timestamp    time.Time
}

// RiskFactors contains factors used in risk calculation
type RiskFactors struct {
	DeviceTrust  float64
	BehaviorRisk float64
	ThreatLevel  float64
	UserID       *uuid.UUID
	IPAddress    string
	Resource     string
	Action       string
	Timestamp    time.Time
}

// PolicyDecision represents a policy evaluation result
type PolicyDecision struct {
	Allowed     bool
	RequiresMFA bool
	SessionTTL  time.Duration
	Reason      string
}

// ZeroTrustSessionManager manages user sessions for zero trust
type ZeroTrustSessionManager struct {
	logger *observability.Logger
}

// RiskCalculator calculates risk scores
type RiskCalculator struct {
	logger *observability.Logger
}

// AnomalyDetector detects behavioral anomalies
type AnomalyDetector struct {
	logger *observability.Logger
}

// BaselineBuilder builds behavioral baselines
type BaselineBuilder struct {
	logger *observability.Logger
}

// ThreatIntelligence provides threat intelligence data
type ThreatIntelligence struct {
	logger *observability.Logger
}

// SignatureEngine detects known attack signatures
type SignatureEngine struct {
	logger *observability.Logger
}

// MLThreatDetector uses ML for threat detection
type MLThreatDetector struct {
	logger *observability.Logger
}

// NewZeroTrustEngine creates a new zero-trust security engine
func NewZeroTrustEngine(logger *observability.Logger) *ZeroTrustEngine {
	config := &ZeroTrustConfig{
		EnableDeviceFingerprinting: true,
		EnableBehaviorAnalysis:     true,
		EnableThreatDetection:      true,
		EnableContinuousAuth:       true,
		RiskThreshold:              0.7,
		SessionTimeout:             30 * time.Minute,
		DeviceTrustDuration:        7 * 24 * time.Hour,
		MaxRiskScore:               1.0,
		RequireMFAForHighRisk:      true,
		TTLCalculationStrategy:     "adaptive", // Use adaptive TTL calculation
	}

	return &ZeroTrustEngine{
		logger:           logger,
		config:           config,
		deviceRegistry:   NewDeviceRegistry(logger),
		behaviorAnalyzer: NewZeroTrustBehaviorAnalyzer(logger),
		threatDetector:   NewZeroTrustThreatDetector(logger),
		policyEngine:     NewPolicyEngine(logger),
		sessionManager:   NewZeroTrustSessionManager(logger),
		riskCalculator:   NewRiskCalculator(logger),
	}
}

// EvaluateAccess performs zero-trust access evaluation
func (z *ZeroTrustEngine) EvaluateAccess(ctx context.Context, request *AccessRequest) (*AccessDecision, error) {
	// 1. Device Trust Evaluation
	deviceTrust, err := z.evaluateDeviceTrust(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("device trust evaluation failed: %w", err)
	}

	// 2. Behavior Analysis
	behaviorRisk, err := z.evaluateBehaviorRisk(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("behavior analysis failed: %w", err)
	}

	// 3. Threat Detection
	threatLevel, err := z.evaluateThreatLevel(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("threat detection failed: %w", err)
	}

	// 4. Calculate Overall Risk Score
	riskScore := z.riskCalculator.CalculateRiskScore(&RiskFactors{
		DeviceTrust:  deviceTrust,
		BehaviorRisk: behaviorRisk,
		ThreatLevel:  threatLevel,
		UserID:       request.UserID,
		IPAddress:    request.IPAddress,
		Resource:     request.Resource,
		Action:       request.Action,
		Timestamp:    request.Timestamp,
	})

	// 5. Policy Evaluation
	userID := uuid.Nil
	if request.UserID != nil {
		userID = *request.UserID
	}

	evalCtx := &PolicyEvaluationContext{
		UserID:      userID,
		IPAddress:   request.IPAddress,
		UserAgent:   request.UserAgent,
		Resource:    request.Resource,
		Action:      request.Action,
		RiskScore:   riskScore,
		DeviceTrust: deviceTrust,
		Timestamp:   request.Timestamp,
	}

	policyResult, err := z.policyEngine.EvaluatePolicy(ctx, evalCtx)
	if err != nil {
		z.logger.Error(ctx, "Policy evaluation failed", err)
		return &AccessDecision{
			Allowed: false,
			Reason:  "Policy evaluation failed",
		}, err
	}

	policyDecision := &policyResult.Decision

	// 6. Make Access Decision
	decision := &AccessDecision{
		Allowed:      riskScore <= z.config.RiskThreshold && policyDecision.Allowed,
		RiskScore:    riskScore,
		DeviceTrust:  deviceTrust,
		BehaviorRisk: behaviorRisk,
		ThreatLevel:  threatLevel,
		RequiresMFA:  riskScore > 0.5 || z.config.RequireMFAForHighRisk,
		SessionTTL:   z.calculateSessionTTL(riskScore),
		Reason:       z.generateDecisionReason(riskScore, policyDecision),
		Timestamp:    time.Now(),
	}

	// 7. Log Security Event
	z.logSecurityEvent(ctx, request, decision)

	// 8. Update User Behavior Profile
	if request.UserID != nil {
		z.behaviorAnalyzer.UpdateProfile(*request.UserID, request)
	}

	return decision, nil
}

// evaluateDeviceTrust evaluates device trustworthiness
func (z *ZeroTrustEngine) evaluateDeviceTrust(ctx context.Context, request *AccessRequest) (float64, error) {
	if !z.config.EnableDeviceFingerprinting {
		return 1.0, nil // Default trust if disabled
	}

	deviceID := z.generateDeviceFingerprint(request)
	device := z.deviceRegistry.GetDevice(deviceID)

	if device == nil {
		// New device - register and assign low trust
		device = &TrustedDevice{
			DeviceID:    deviceID,
			UserID:      *request.UserID,
			Fingerprint: deviceID,
			TrustedAt:   time.Now(),
			LastSeen:    time.Now(),
			TrustLevel:  0.3, // Low initial trust
			TrustScore:  30,  // Low initial trust
			IsActive:    true,
			ExpiresAt:   time.Now().Add(24 * time.Hour),
		}
		z.deviceRegistry.RegisterDevice(device)

		z.logger.Warn(ctx, "New device detected", map[string]interface{}{
			"device_id": deviceID,
			"user_id":   request.UserID,
			"ip":        request.IPAddress,
		})

		return device.TrustLevel, nil
	}

	// Update device last seen
	device.LastSeen = time.Now()

	// Calculate trust based on various factors
	trustScore := z.calculateDeviceTrustScore(device, request)
	device.TrustLevel = trustScore

	return trustScore, nil
}

// evaluateBehaviorRisk evaluates user behavior risk
func (z *ZeroTrustEngine) evaluateBehaviorRisk(ctx context.Context, request *AccessRequest) (float64, error) {
	if !z.config.EnableBehaviorAnalysis || request.UserID == nil {
		return 0.0, nil // No risk if disabled or no user
	}

	profile := z.behaviorAnalyzer.GetProfile(*request.UserID)
	if profile == nil {
		// New user - create profile with baseline risk
		profile = &UserBehaviorProfile{
			UserID:              *request.UserID,
			TypicalLoginTimes:   []time.Time{},
			TypicalLocations:    []Location{},
			TypicalDevices:      []string{},
			TypicalActions:      make(map[string]int),
			SessionDurations:    []time.Duration{},
			APIUsagePatterns:    make(map[string]int),
			RiskScore:           0.5, // Medium baseline risk
			LastUpdated:         time.Now(),
			BaselineEstablished: false,
		}
		z.behaviorAnalyzer.CreateProfile(profile)
		return 0.5, nil
	}

	// Analyze current behavior against profile
	riskScore := z.behaviorAnalyzer.AnalyzeBehaviorProfile(profile, request)

	return riskScore, nil
}

// evaluateThreatLevel evaluates current threat level
func (z *ZeroTrustEngine) evaluateThreatLevel(ctx context.Context, request *AccessRequest) (float64, error) {
	if !z.config.EnableThreatDetection {
		return 0.0, nil
	}

	// Check for active threats
	threatLevel := z.threatDetector.EvaluateThreats(ctx, request)

	return threatLevel, nil
}

// generateDeviceFingerprint creates a unique device fingerprint
func (z *ZeroTrustEngine) generateDeviceFingerprint(request *AccessRequest) string {
	// Combine various device attributes to create fingerprint
	fingerprint := fmt.Sprintf("%s|%s|%s",
		request.IPAddress,
		request.UserAgent,
		request.DeviceID)

	hash := sha256.Sum256([]byte(fingerprint))
	return hex.EncodeToString(hash[:])
}

// extractDeviceAttributes extracts device attributes from request
func (z *ZeroTrustEngine) extractDeviceAttributes(request *AccessRequest) map[string]interface{} {
	return map[string]interface{}{
		"ip_address":  request.IPAddress,
		"user_agent":  request.UserAgent,
		"device_id":   request.DeviceID,
		"timestamp":   request.Timestamp,
		"geolocation": z.getGeolocation(request.IPAddress),
	}
}

// calculateDeviceTrustScore calculates device trust score
func (z *ZeroTrustEngine) calculateDeviceTrustScore(device *TrustedDevice, request *AccessRequest) float64 {
	baseScore := device.TrustLevel

	// Increase trust over time with successful authentications
	daysSinceFirstSeen := time.Since(device.TrustedAt).Hours() / 24
	timeBonus := min(daysSinceFirstSeen/30, 0.3) // Max 0.3 bonus over 30 days

	// Decrease trust if device hasn't been seen recently
	daysSinceLastSeen := time.Since(device.LastSeen).Hours() / 24
	stalePenalty := 0.0
	if daysSinceLastSeen > 30 {
		stalePenalty = min(daysSinceLastSeen/365, 0.5) // Max 0.5 penalty
	}

	// Check for compromise indicators (simplified)
	compromisePenalty := 0.0
	if !device.IsActive {
		compromisePenalty = 0.8
	}

	finalScore := baseScore + timeBonus - stalePenalty - compromisePenalty
	return max(0.0, min(1.0, finalScore))
}

// calculateSessionTTL calculates session TTL based on risk score with advanced algorithms
func (z *ZeroTrustEngine) calculateSessionTTL(riskScore float64) time.Duration {
	baseTTL := z.config.SessionTimeout
	minTTL := 5 * time.Minute
	maxTTL := baseTTL

	// Clamp risk score to valid range [0.0, 1.0]
	if riskScore < 0 {
		riskScore = 0
	}
	if riskScore > 1 {
		riskScore = 1
	}

	// Advanced TTL calculation with multiple strategies
	var adjustedTTL time.Duration

	switch z.config.TTLCalculationStrategy {
	case "exponential":
		// Exponential decay for high-risk scenarios
		// Formula: TTL = baseTTL * e^(-k * riskScore)
		k := 2.0 // Decay constant
		multiplier := math.Exp(-k * riskScore)
		adjustedTTL = time.Duration(float64(baseTTL) * multiplier)

	case "logarithmic":
		// Logarithmic scaling for gradual reduction
		// Formula: TTL = baseTTL * (1 - log(1 + riskScore * 9) / log(10))
		if riskScore == 0 {
			adjustedTTL = baseTTL
		} else {
			logFactor := math.Log10(1 + riskScore*9)
			multiplier := 1.0 - logFactor
			adjustedTTL = time.Duration(float64(baseTTL) * multiplier)
		}

	case "stepped":
		// Stepped reduction based on risk thresholds
		switch {
		case riskScore < 0.2:
			adjustedTTL = baseTTL // Low risk: full TTL
		case riskScore < 0.4:
			adjustedTTL = baseTTL * 3 / 4 // Medium-low: 75%
		case riskScore < 0.6:
			adjustedTTL = baseTTL / 2 // Medium: 50%
		case riskScore < 0.8:
			adjustedTTL = baseTTL / 4 // High: 25%
		default:
			adjustedTTL = minTTL // Critical: minimum
		}

	case "adaptive":
		// Adaptive calculation based on historical patterns
		adjustedTTL = z.calculateAdaptiveTTL(riskScore, baseTTL)

	default: // "linear" - original implementation
		// Linear reduction: TTL = baseTTL * (1 - riskScore)
		riskMultiplier := 1.0 - riskScore
		adjustedTTL = time.Duration(float64(baseTTL) * riskMultiplier)
	}

	// Apply time-of-day adjustments for enhanced security
	adjustedTTL = z.applyTimeBasedAdjustments(adjustedTTL, riskScore)

	// Clamp the adjusted TTL to the valid range
	if adjustedTTL < minTTL {
		return minTTL
	}
	if adjustedTTL > maxTTL {
		return maxTTL
	}

	return adjustedTTL
}

// calculateAdaptiveTTL calculates TTL using adaptive algorithms based on historical patterns
func (z *ZeroTrustEngine) calculateAdaptiveTTL(riskScore float64, baseTTL time.Duration) time.Duration {
	// Adaptive calculation considers:
	// 1. Historical risk patterns for this user/device
	// 2. Current system load and threat level
	// 3. Time-based patterns (business hours vs off-hours)

	// Base calculation using exponential decay
	k := 1.5 // Adaptive decay constant
	baseMultiplier := math.Exp(-k * riskScore)

	// Adjust based on time of day (stricter during off-hours)
	timeAdjustment := z.getTimeBasedRiskAdjustment()
	adaptiveMultiplier := baseMultiplier * (1.0 - timeAdjustment*0.2)

	// Consider system threat level
	systemThreatLevel := z.getSystemThreatLevel()
	threatAdjustment := 1.0 - (systemThreatLevel * 0.3)

	finalMultiplier := adaptiveMultiplier * threatAdjustment
	if finalMultiplier < 0.1 {
		finalMultiplier = 0.1 // Minimum 10% of base TTL
	}

	return time.Duration(float64(baseTTL) * finalMultiplier)
}

// applyTimeBasedAdjustments applies time-of-day security adjustments
func (z *ZeroTrustEngine) applyTimeBasedAdjustments(ttl time.Duration, riskScore float64) time.Duration {
	now := time.Now()
	hour := now.Hour()

	// Define business hours (9 AM to 6 PM)
	isBusinessHours := hour >= 9 && hour < 18
	isWeekend := now.Weekday() == time.Saturday || now.Weekday() == time.Sunday

	// Apply stricter policies during off-hours
	if !isBusinessHours || isWeekend {
		// Reduce TTL by 20-50% during off-hours based on risk
		reductionFactor := 0.2 + (riskScore * 0.3) // 20% to 50% reduction
		adjustedTTL := time.Duration(float64(ttl) * (1.0 - reductionFactor))

		// Ensure minimum TTL
		minOffHoursTTL := 3 * time.Minute
		if adjustedTTL < minOffHoursTTL {
			return minOffHoursTTL
		}
		return adjustedTTL
	}

	return ttl
}

// getTimeBasedRiskAdjustment returns risk adjustment based on time patterns
func (z *ZeroTrustEngine) getTimeBasedRiskAdjustment() float64 {
	now := time.Now()
	hour := now.Hour()

	// Higher risk during unusual hours
	switch {
	case hour >= 2 && hour < 6: // Late night/early morning
		return 0.4
	case hour >= 22 || hour < 2: // Late evening/night
		return 0.3
	case hour >= 6 && hour < 9: // Early morning
		return 0.1
	case hour >= 18 && hour < 22: // Evening
		return 0.1
	default: // Business hours
		return 0.0
	}
}

// getSystemThreatLevel returns current system-wide threat level
func (z *ZeroTrustEngine) getSystemThreatLevel() float64 {
	// This would typically query threat intelligence feeds
	// and system-wide security metrics
	// For now, return a baseline threat level
	return 0.1 // 10% baseline threat level
}

// generateDecisionReason generates human-readable decision reason
func (z *ZeroTrustEngine) generateDecisionReason(riskScore float64, policyDecision *PolicyDecision) string {
	if riskScore > z.config.RiskThreshold {
		return fmt.Sprintf("Access denied: Risk score %.2f exceeds threshold %.2f",
			riskScore, z.config.RiskThreshold)
	}

	if !policyDecision.Allowed {
		return fmt.Sprintf("Access denied: %s", policyDecision.Reason)
	}

	return "Access granted: All security checks passed"
}

// logSecurityEvent logs security events for audit and monitoring
func (z *ZeroTrustEngine) logSecurityEvent(ctx context.Context, request *AccessRequest, decision *AccessDecision) {
	event := &SecurityEvent{
		EventID:   uuid.New().String(),
		Type:      "access_evaluation",
		Severity:  z.getSeverityFromRisk(decision.RiskScore),
		UserID:    request.UserID,
		DeviceID:  request.DeviceID,
		IPAddress: request.IPAddress,
		UserAgent: request.UserAgent,
		Resource:  request.Resource,
		Action:    request.Action,
		Timestamp: time.Now(),
		RiskScore: decision.RiskScore,
		Blocked:   !decision.Allowed,
		Context: map[string]interface{}{
			"device_trust":  decision.DeviceTrust,
			"behavior_risk": decision.BehaviorRisk,
			"threat_level":  decision.ThreatLevel,
			"requires_mfa":  decision.RequiresMFA,
			"session_ttl":   decision.SessionTTL,
		},
	}

	z.logger.Info(ctx, "Zero-trust access evaluation", map[string]interface{}{
		"event_id":      event.EventID,
		"allowed":       decision.Allowed,
		"risk_score":    decision.RiskScore,
		"device_trust":  decision.DeviceTrust,
		"behavior_risk": decision.BehaviorRisk,
		"threat_level":  decision.ThreatLevel,
		"requires_mfa":  decision.RequiresMFA,
	})
}

// getSeverityFromRisk converts risk score to severity level
func (z *ZeroTrustEngine) getSeverityFromRisk(riskScore float64) string {
	switch {
	case riskScore >= 0.8:
		return "critical"
	case riskScore >= 0.6:
		return "high"
	case riskScore >= 0.4:
		return "medium"
	default:
		return "low"
	}
}

// getGeolocation gets geolocation from IP address
func (z *ZeroTrustEngine) getGeolocation(ipAddress string) *Location {
	// In production, use a geolocation service
	// For now, return a placeholder
	return &Location{
		Country: "Unknown",
		Region:  "Unknown",
		City:    "Unknown",
		IPRange: z.getIPRange(ipAddress),
	}
}

// getIPRange gets IP range for the given IP address
func (z *ZeroTrustEngine) getIPRange(ipAddress string) string {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return "unknown"
	}

	if ip.To4() != nil {
		// IPv4 - return /24 subnet
		return fmt.Sprintf("%s/24", ip.Mask(net.CIDRMask(24, 32)))
	}

	// IPv6 - return /64 subnet
	return fmt.Sprintf("%s/64", ip.Mask(net.CIDRMask(64, 128)))
}

// NewDeviceRegistry creates a new device registry
func NewDeviceRegistry(logger *observability.Logger) *DeviceRegistry {
	return &DeviceRegistry{
		devices: make(map[string]*TrustedDevice),
		logger:  logger,
	}
}

// GetDevice retrieves a device by ID
func (d *DeviceRegistry) GetDevice(deviceID string) *TrustedDevice {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.devices[deviceID]
}

// RegisterDevice registers a new device
func (d *DeviceRegistry) RegisterDevice(device *TrustedDevice) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.devices[device.DeviceID] = device
}

// NewZeroTrustBehaviorAnalyzer creates a new zero trust behavior analyzer
func NewZeroTrustBehaviorAnalyzer(logger *observability.Logger) *ZeroTrustBehaviorAnalyzer {
	return &ZeroTrustBehaviorAnalyzer{
		logger:       logger,
		userProfiles: make(map[uuid.UUID]*UserBehaviorProfile),
	}
}

// GetProfile gets a user behavior profile
func (ba *ZeroTrustBehaviorAnalyzer) GetProfile(userID uuid.UUID) *UserBehaviorProfile {
	ba.mu.RLock()
	defer ba.mu.RUnlock()
	return ba.userProfiles[userID]
}

// CreateProfile creates a new user behavior profile
func (ba *ZeroTrustBehaviorAnalyzer) CreateProfile(profile *UserBehaviorProfile) {
	ba.mu.Lock()
	defer ba.mu.Unlock()
	ba.userProfiles[profile.UserID] = profile
}

// UpdateProfile updates a user behavior profile
func (ba *ZeroTrustBehaviorAnalyzer) UpdateProfile(userID uuid.UUID, request *AccessRequest) {
	ba.mu.Lock()
	defer ba.mu.Unlock()

	profile := ba.userProfiles[userID]
	if profile == nil {
		return
	}

	// Update profile based on request
	profile.LastUpdated = time.Now()
	// Add more profile update logic here
}

// AnalyzeBehaviorProfile analyzes current behavior against profile
func (ba *ZeroTrustBehaviorAnalyzer) AnalyzeBehaviorProfile(profile *UserBehaviorProfile, request *AccessRequest) float64 {
	riskScore := 0.0

	// Check time-based anomalies
	currentHour := time.Now().Hour()
	if currentHour < 6 || currentHour > 22 {
		riskScore += 0.2
	}

	// Check location-based anomalies
	if request.IPAddress != "" {
		// Simplified location check
		riskScore += 0.1
	}

	// Ensure risk score is within bounds
	if riskScore > 1.0 {
		riskScore = 1.0
	}

	return riskScore
}

// AnalyzeBehaviorProfile analyzes current behavior against profile
func (b *BehaviorAnalyzer) AnalyzeBehaviorProfile(profile *UserBehaviorProfile, request *AccessRequest) float64 {
	riskScore := 0.0

	// Check time-based anomalies
	currentHour := request.Timestamp.Hour()
	typicalHours := make(map[int]bool)
	for _, loginTime := range profile.TypicalLoginTimes {
		typicalHours[loginTime.Hour()] = true
	}

	if !typicalHours[currentHour] && len(profile.TypicalLoginTimes) > 10 {
		riskScore += 0.2 // Unusual time
	}

	// Check action frequency anomalies
	action := fmt.Sprintf("%s:%s", request.Resource, request.Action)
	if profile.TypicalActions[action] == 0 && len(profile.TypicalActions) > 20 {
		riskScore += 0.3 // Unusual action
	}

	return min(1.0, riskScore)
}

// NewZeroTrustThreatDetector creates a new zero trust threat detector
func NewZeroTrustThreatDetector(logger *observability.Logger) *ZeroTrustThreatDetector {
	return &ZeroTrustThreatDetector{
		logger:        logger,
		activeThreats: make(map[string]*ActiveThreat),
	}
}

// EvaluateThreats evaluates threats for a request
func (td *ZeroTrustThreatDetector) EvaluateThreats(ctx context.Context, request *AccessRequest) float64 {
	threatLevel := 0.0

	// Check for suspicious IP patterns
	if request.IPAddress == "127.0.0.1" || strings.HasPrefix(request.IPAddress, "192.168.") {
		threatLevel += 0.1
	}

	// Check for suspicious user agents
	if request.UserAgent == "" {
		threatLevel += 0.2
	}

	// Ensure threat level is within bounds
	if threatLevel > 1.0 {
		threatLevel = 1.0
	}

	return threatLevel
}

// EvaluateThreats evaluates threats for the given request
func (t *ThreatDetector) EvaluateThreats(request *AccessRequest) float64 {
	threatLevel := 0.0

	// Check for brute force patterns
	if t.detectBruteForce(request) {
		threatLevel += 0.8
	}

	// Check for suspicious API usage
	if t.detectSuspiciousAPI(request) {
		threatLevel += 0.6
	}

	// Check for bot activity
	if t.detectBotActivity(request) {
		threatLevel += 0.7
	}

	return min(1.0, threatLevel)
}

// detectBruteForce detects brute force attack patterns
func (t *ThreatDetector) detectBruteForce(request *AccessRequest) bool {
	// Simple implementation - in production, use more sophisticated detection
	return false
}

// detectSuspiciousAPI detects suspicious API usage patterns
func (t *ThreatDetector) detectSuspiciousAPI(request *AccessRequest) bool {
	// Check for rapid API calls, unusual endpoints, etc.
	return false
}

// detectBotActivity detects bot activity
func (t *ThreatDetector) detectBotActivity(request *AccessRequest) bool {
	// Check user agent patterns, request timing, etc.
	return false
}

// Note: PolicyEngine is now defined in policy_engine.go

// DeviceTrustManager is defined in rate_limiter.go to avoid duplication

// NewZeroTrustSessionManager creates a new zero trust session manager
func NewZeroTrustSessionManager(logger *observability.Logger) *ZeroTrustSessionManager {
	return &ZeroTrustSessionManager{
		logger: logger,
	}
}

// NewRiskCalculator creates a new risk calculator
func NewRiskCalculator(logger *observability.Logger) *RiskCalculator {
	return &RiskCalculator{
		logger: logger,
	}
}

// CalculateRiskScore calculates overall risk score
func (r *RiskCalculator) CalculateRiskScore(factors *RiskFactors) float64 {
	// Weighted risk calculation
	deviceWeight := 0.3
	behaviorWeight := 0.4
	threatWeight := 0.3

	riskScore := (1.0-factors.DeviceTrust)*deviceWeight +
		factors.BehaviorRisk*behaviorWeight +
		factors.ThreatLevel*threatWeight

	return min(1.0, riskScore)
}

// Helper functions
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
