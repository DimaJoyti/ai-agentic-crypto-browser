package security

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"golang.org/x/time/rate"
)

// RateLimiter provides rate limiting functionality
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	config   *SecurityConfig
	mu       sync.RWMutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(config *SecurityConfig) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		config:   config,
	}
}

// Allow checks if a request is allowed for the given key and type
func (rl *RateLimiter) Allow(key, requestType string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiterKey := key + ":" + requestType
	limiter, exists := rl.limiters[limiterKey]

	if !exists {
		var limit rate.Limit
		var burst int

		switch requestType {
		case "login":
			limit = rate.Limit(rl.config.LoginRateLimit) / rate.Limit(rl.config.RateLimitWindow.Minutes())
			burst = rl.config.LoginRateLimit
		case "api":
			limit = rate.Limit(rl.config.APIRateLimit) / rate.Limit(rl.config.RateLimitWindow.Minutes())
			burst = rl.config.APIRateLimit
		default:
			limit = rate.Limit(100) / rate.Limit(rl.config.RateLimitWindow.Minutes())
			burst = 100
		}

		limiter = rate.NewLimiter(limit, burst)
		rl.limiters[limiterKey] = limiter
	}

	return limiter.Allow()
}

// BruteForceProtection provides brute force attack protection
type BruteForceProtection struct {
	attempts map[string]*AttemptRecord
	config   *SecurityConfig
	mu       sync.RWMutex
}

// AttemptRecord tracks failed login attempts
type AttemptRecord struct {
	Count       int
	LastAttempt time.Time
	LockedUntil time.Time
}

// NewBruteForceProtection creates a new brute force protection instance
func NewBruteForceProtection(config *SecurityConfig) *BruteForceProtection {
	return &BruteForceProtection{
		attempts: make(map[string]*AttemptRecord),
		config:   config,
	}
}

// IsBlocked checks if an IP/username combination is blocked
func (bfp *BruteForceProtection) IsBlocked(ipAddress, username string) bool {
	bfp.mu.RLock()
	defer bfp.mu.RUnlock()

	key := ipAddress + ":" + username
	record, exists := bfp.attempts[key]

	if !exists {
		return false
	}

	// Check if lockout has expired
	if time.Now().After(record.LockedUntil) {
		return false
	}

	return record.Count >= bfp.config.MaxLoginAttempts
}

// RecordFailedAttempt records a failed login attempt
func (bfp *BruteForceProtection) RecordFailedAttempt(ipAddress, username string) {
	bfp.mu.Lock()
	defer bfp.mu.Unlock()

	key := ipAddress + ":" + username
	record, exists := bfp.attempts[key]

	if !exists {
		record = &AttemptRecord{}
		bfp.attempts[key] = record
	}

	record.Count++
	record.LastAttempt = time.Now()

	// Calculate lockout duration (progressive if enabled)
	lockoutDuration := bfp.config.LockoutDuration
	if bfp.config.ProgressiveLockout {
		multiplier := time.Duration(record.Count)
		lockoutDuration = lockoutDuration * multiplier
	}

	record.LockedUntil = time.Now().Add(lockoutDuration)
}

// ClearFailedAttempts clears failed attempts for successful login
func (bfp *BruteForceProtection) ClearFailedAttempts(ipAddress, username string) {
	bfp.mu.Lock()
	defer bfp.mu.Unlock()

	key := ipAddress + ":" + username
	delete(bfp.attempts, key)
}

// DeviceTrustManager manages device trust
type DeviceTrustManager struct {
	logger         *observability.Logger
	trustedDevices map[string]*TrustedDevice
	config         *SecurityConfig
	mu             sync.RWMutex
}

// TrustedDevice represents a trusted device
type TrustedDevice struct {
	UserID      uuid.UUID              `json:"user_id"`
	DeviceID    string                 `json:"device_id"`
	DeviceName  string                 `json:"device_name"`
	Fingerprint string                 `json:"fingerprint"`
	TrustedAt   time.Time              `json:"trusted_at"`
	LastSeen    time.Time              `json:"last_seen"`
	ExpiresAt   time.Time              `json:"expires_at"`
	IsActive    bool                   `json:"is_active"`
	TrustScore  int                    `json:"trust_score"`
	TrustLevel  float64                `json:"trust_level"`
	IPAddresses []string               `json:"ip_addresses"`
	UserAgents  []string               `json:"user_agents"`
	Attributes  map[string]interface{} `json:"attributes"`
	RiskFactors []string               `json:"risk_factors"`
}

// NewDeviceTrustManager creates a new device trust manager
func NewDeviceTrustManager(logger *observability.Logger, config *SecurityConfig) *DeviceTrustManager {
	return &DeviceTrustManager{
		logger:         logger,
		trustedDevices: make(map[string]*TrustedDevice),
		config:         config,
	}
}

// IsDeviceTrusted checks if a device is trusted for a user
func (dtm *DeviceTrustManager) IsDeviceTrusted(userID uuid.UUID, deviceFingerprint string) bool {
	dtm.mu.RLock()
	defer dtm.mu.RUnlock()

	key := userID.String() + ":" + deviceFingerprint
	device, exists := dtm.trustedDevices[key]

	if !exists {
		return false
	}

	// Check if device trust has expired
	if time.Now().After(device.ExpiresAt) {
		return false
	}

	return device.IsActive
}

// TrustDevice marks a device as trusted
func (dtm *DeviceTrustManager) TrustDevice(userID uuid.UUID, deviceFingerprint, deviceName, ipAddress, userAgent string) error {
	dtm.mu.Lock()
	defer dtm.mu.Unlock()

	key := userID.String() + ":" + deviceFingerprint

	device := &TrustedDevice{
		UserID:      userID,
		DeviceID:    key,
		DeviceName:  deviceName,
		Fingerprint: deviceFingerprint,
		TrustedAt:   time.Now(),
		LastSeen:    time.Now(),
		ExpiresAt:   time.Now().Add(dtm.config.DeviceTrustDuration),
		IsActive:    true,
		TrustScore:  100,
		TrustLevel:  1.0,
		IPAddresses: []string{ipAddress},
		UserAgents:  []string{userAgent},
	}

	dtm.trustedDevices[key] = device
	return nil
}

// RegisterDevice registers a new trusted device
func (dtm *DeviceTrustManager) RegisterDevice(device *TrustedDevice) error {
	dtm.mu.Lock()
	defer dtm.mu.Unlock()

	key := device.UserID.String() + ":" + device.Fingerprint
	dtm.trustedDevices[key] = device
	return nil
}

// GetDevice retrieves a device by ID
func (dtm *DeviceTrustManager) GetDevice(deviceID string) (*TrustedDevice, error) {
	dtm.mu.RLock()
	defer dtm.mu.RUnlock()

	for _, device := range dtm.trustedDevices {
		if device.DeviceID == deviceID {
			return device, nil
		}
	}

	return nil, fmt.Errorf("device not found: %s", deviceID)
}

// UpdateDevice updates a device
func (dtm *DeviceTrustManager) UpdateDevice(device *TrustedDevice) error {
	dtm.mu.Lock()
	defer dtm.mu.Unlock()

	key := device.UserID.String() + ":" + device.Fingerprint
	dtm.trustedDevices[key] = device
	return nil
}

// UpdateTrustLevel updates the trust level of a device
func (dtm *DeviceTrustManager) UpdateTrustLevel(deviceID string, trustLevel float64) error {
	dtm.mu.Lock()
	defer dtm.mu.Unlock()

	for _, device := range dtm.trustedDevices {
		if device.DeviceID == deviceID {
			device.TrustLevel = trustLevel
			device.TrustScore = int(trustLevel * 100)
			device.LastSeen = time.Now()
			return nil
		}
	}

	return fmt.Errorf("device not found: %s", deviceID)
}

// BehaviorAnalyzer analyzes user behavior patterns
type BehaviorAnalyzer struct {
	logger   *observability.Logger
	patterns map[string]*BehaviorPattern
	config   *SecurityConfig
	mu       sync.RWMutex
}

// BehaviorPattern represents a user's behavior pattern
type BehaviorPattern struct {
	UserID           uuid.UUID          `json:"user_id"`
	TypicalLocations []string           `json:"typical_locations"`
	TypicalTimes     []TimeRange        `json:"typical_times"`
	TypicalDevices   []string           `json:"typical_devices"`
	TradingPatterns  *TradingBehavior   `json:"trading_patterns"`
	LastUpdated      time.Time          `json:"last_updated"`
	Anomalies        []*BehaviorAnomaly `json:"anomalies"`
}

// TimeRange represents a time range
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// TradingBehavior represents trading behavior patterns
type TradingBehavior struct {
	TypicalVolume    decimal.Decimal `json:"typical_volume"`
	TypicalPairs     []string        `json:"typical_pairs"`
	TypicalFrequency int             `json:"typical_frequency"`
	RiskTolerance    string          `json:"risk_tolerance"`
}

// BehaviorAnomaly represents a detected behavior anomaly
type BehaviorAnomaly struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Severity    int                    `json:"severity"`
	DetectedAt  time.Time              `json:"detected_at"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// NewBehaviorAnalyzer creates a new behavior analyzer
func NewBehaviorAnalyzer(logger *observability.Logger, config *SecurityConfig) *BehaviorAnalyzer {
	return &BehaviorAnalyzer{
		logger:   logger,
		patterns: make(map[string]*BehaviorPattern),
		config:   config,
	}
}

// AnalyzeBehavior analyzes user behavior and returns anomaly score
func (ba *BehaviorAnalyzer) AnalyzeBehavior(userID uuid.UUID, request *AuthenticationRequest) int {
	ba.mu.RLock()
	defer ba.mu.RUnlock()

	pattern, exists := ba.patterns[userID.String()]
	if !exists {
		// No pattern established yet, return moderate score
		return 30
	}

	anomalyScore := 0

	// Check location anomaly
	if !ba.isTypicalLocation(pattern, request.IPAddress) {
		anomalyScore += 25
	}

	// Check time anomaly
	if !ba.isTypicalTime(pattern, time.Now()) {
		anomalyScore += 15
	}

	// Check device anomaly
	if !ba.isTypicalDevice(pattern, request.DeviceFingerprint) {
		anomalyScore += 20
	}

	return anomalyScore
}

// Helper methods for behavior analysis
func (ba *BehaviorAnalyzer) isTypicalLocation(pattern *BehaviorPattern, ipAddress string) bool {
	// Simplified location check
	for _, location := range pattern.TypicalLocations {
		if location == ipAddress {
			return true
		}
	}
	return false
}

func (ba *BehaviorAnalyzer) isTypicalTime(pattern *BehaviorPattern, timestamp time.Time) bool {
	// Simplified time check
	hour := timestamp.Hour()
	return hour >= 8 && hour <= 22 // Typical business hours
}

func (ba *BehaviorAnalyzer) isTypicalDevice(pattern *BehaviorPattern, deviceFingerprint string) bool {
	// Simplified device check
	for _, device := range pattern.TypicalDevices {
		if device == deviceFingerprint {
			return true
		}
	}
	return false
}

// ThreatDetector detects various security threats
type ThreatDetector struct {
	logger *observability.Logger
	config *SecurityConfig
}

// NewThreatDetector creates a new threat detector
func NewThreatDetector(logger *observability.Logger, config *SecurityConfig) *ThreatDetector {
	return &ThreatDetector{
		logger: logger,
		config: config,
	}
}

// AnalyzeRequest analyzes a request for threats and returns threat level (0-100)
func (td *ThreatDetector) AnalyzeRequest(ctx context.Context, request *AuthenticationRequest) int {
	threatLevel := 0

	// Check for suspicious user agents
	if td.isSuspiciousUserAgent(request.UserAgent) {
		threatLevel += 30
	}

	// Check for suspicious IP patterns
	if td.isSuspiciousIP(request.IPAddress) {
		threatLevel += 40
	}

	// Check for automation indicators
	if td.isAutomatedRequest(request) {
		threatLevel += 25
	}

	// Ensure threat level is within bounds
	if threatLevel > 100 {
		threatLevel = 100
	}

	return threatLevel
}

// Helper methods for threat detection
func (td *ThreatDetector) isSuspiciousUserAgent(userAgent string) bool {
	suspiciousPatterns := []string{
		"bot", "crawler", "spider", "scraper",
		"curl", "wget", "python", "go-http-client",
	}

	userAgentLower := strings.ToLower(userAgent)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(userAgentLower, pattern) {
			return true
		}
	}

	return false
}

func (td *ThreatDetector) isSuspiciousIP(ipAddress string) bool {
	// This would integrate with threat intelligence feeds
	// For now, just check for localhost and private ranges
	return ipAddress == "127.0.0.1" || strings.HasPrefix(ipAddress, "192.168.")
}

func (td *ThreatDetector) isAutomatedRequest(request *AuthenticationRequest) bool {
	// Check for signs of automation
	if request.UserAgent == "" {
		return true
	}

	// Check for rapid-fire requests (would need request timing data)
	return false
}
