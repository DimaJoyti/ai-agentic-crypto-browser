package security

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ai-agentic-browser/internal/auth"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/argon2"
)

// AuthManager provides comprehensive authentication and authorization management
type AuthManager struct {
	logger         *observability.Logger
	jwtService     *auth.JWTService
	mfaService     *auth.MFAService
	rbacService    *auth.RBACService
	sessionManager *SessionManager
	apiKeyManager  *APIKeyManager
	rateLimiter    *RateLimiter
	securityConfig *SecurityConfig

	// Security features
	bruteForceProtection *BruteForceProtection
	deviceTrustManager   *DeviceTrustManager
	behaviorAnalyzer     *BehaviorAnalyzer
	threatDetector       *ThreatDetector

	// State management
	activeSessions       map[string]*SecuritySession
	activeAPIKeys        map[string]*APIKey
	blockedIPs           map[string]*IPBlock
	suspiciousActivities map[string]*SuspiciousActivity

	mu sync.RWMutex
}

// SecurityConfig holds comprehensive security configuration
type SecurityConfig struct {
	// Authentication settings
	RequireMFA               bool          `yaml:"require_mfa"`
	SessionTimeout           time.Duration `yaml:"session_timeout"`
	MaxConcurrentSessions    int           `yaml:"max_concurrent_sessions"`
	RequireStrongPasswords   bool          `yaml:"require_strong_passwords"`
	PasswordMinLength        int           `yaml:"password_min_length"`
	PasswordRequireUppercase bool          `yaml:"password_require_uppercase"`
	PasswordRequireLowercase bool          `yaml:"password_require_lowercase"`
	PasswordRequireNumbers   bool          `yaml:"password_require_numbers"`
	PasswordRequireSymbols   bool          `yaml:"password_require_symbols"`

	// Rate limiting
	LoginRateLimit  int           `yaml:"login_rate_limit"`
	APIRateLimit    int           `yaml:"api_rate_limit"`
	RateLimitWindow time.Duration `yaml:"rate_limit_window"`

	// Brute force protection
	MaxLoginAttempts   int           `yaml:"max_login_attempts"`
	LockoutDuration    time.Duration `yaml:"lockout_duration"`
	ProgressiveLockout bool          `yaml:"progressive_lockout"`

	// Device trust
	RequireDeviceRegistration bool          `yaml:"require_device_registration"`
	DeviceTrustDuration       time.Duration `yaml:"device_trust_duration"`
	MaxTrustedDevices         int           `yaml:"max_trusted_devices"`

	// API security
	RequireAPIKeyAuth      bool          `yaml:"require_api_key_auth"`
	APIKeyRotationInterval time.Duration `yaml:"api_key_rotation_interval"`
	AllowedIPRanges        []string      `yaml:"allowed_ip_ranges"`

	// Advanced security
	EnableBehaviorAnalysis bool `yaml:"enable_behavior_analysis"`
	EnableThreatDetection  bool `yaml:"enable_threat_detection"`
	EnableZeroTrust        bool `yaml:"enable_zero_trust"`
	RequireGeoVerification bool `yaml:"require_geo_verification"`

	// Audit and compliance
	EnableAuditLogging   bool          `yaml:"enable_audit_logging"`
	RetainAuditLogs      time.Duration `yaml:"retain_audit_logs"`
	EnableComplianceMode bool          `yaml:"enable_compliance_mode"`
}

// SecuritySession represents an authenticated session with security context
type SecuritySession struct {
	ID             string                 `json:"id"`
	UserID         uuid.UUID              `json:"user_id"`
	DeviceID       string                 `json:"device_id"`
	IPAddress      string                 `json:"ip_address"`
	UserAgent      string                 `json:"user_agent"`
	CreatedAt      time.Time              `json:"created_at"`
	LastActivity   time.Time              `json:"last_activity"`
	ExpiresAt      time.Time              `json:"expires_at"`
	IsActive       bool                   `json:"is_active"`
	MFAVerified    bool                   `json:"mfa_verified"`
	DeviceTrusted  bool                   `json:"device_trusted"`
	RiskScore      int                    `json:"risk_score"`
	Permissions    []string               `json:"permissions"`
	TradingEnabled bool                   `json:"trading_enabled"`
	MaxTradeAmount decimal.Decimal        `json:"max_trade_amount"`
	AllowedPairs   []string               `json:"allowed_pairs"`
	SecurityFlags  map[string]interface{} `json:"security_flags"`
	GeoLocation    *GeoLocation           `json:"geo_location,omitempty"`
}

// APIKey represents an API key with security metadata
type APIKey struct {
	ID             string                 `json:"id"`
	UserID         uuid.UUID              `json:"user_id"`
	Name           string                 `json:"name"`
	KeyHash        string                 `json:"key_hash"`
	Permissions    []string               `json:"permissions"`
	IPWhitelist    []string               `json:"ip_whitelist"`
	CreatedAt      time.Time              `json:"created_at"`
	LastUsed       time.Time              `json:"last_used"`
	ExpiresAt      *time.Time             `json:"expires_at,omitempty"`
	IsActive       bool                   `json:"is_active"`
	RateLimit      int                    `json:"rate_limit"`
	TradingEnabled bool                   `json:"trading_enabled"`
	MaxTradeAmount decimal.Decimal        `json:"max_trade_amount"`
	AllowedPairs   []string               `json:"allowed_pairs"`
	SecurityLevel  SecurityLevel          `json:"security_level"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// SecurityLevel defines different security levels for API keys
type SecurityLevel string

const (
	SecurityLevelReadOnly SecurityLevel = "read_only"
	SecurityLevelTrading  SecurityLevel = "trading"
	SecurityLevelAdmin    SecurityLevel = "admin"
)

// GeoLocation represents geographical location information
type GeoLocation struct {
	Country     string  `json:"country"`
	Region      string  `json:"region"`
	City        string  `json:"city"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	ISP         string  `json:"isp"`
	Timezone    string  `json:"timezone"`
	IsVPN       bool    `json:"is_vpn"`
	IsTor       bool    `json:"is_tor"`
	ThreatLevel int     `json:"threat_level"`
}

// IPBlock represents a blocked IP address
type IPBlock struct {
	IPAddress   string    `json:"ip_address"`
	Reason      string    `json:"reason"`
	BlockedAt   time.Time `json:"blocked_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	Attempts    int       `json:"attempts"`
	Severity    string    `json:"severity"`
	IsAutomatic bool      `json:"is_automatic"`
}

// SuspiciousActivity represents detected suspicious activity
type SuspiciousActivity struct {
	ID           string                 `json:"id"`
	UserID       *uuid.UUID             `json:"user_id,omitempty"`
	IPAddress    string                 `json:"ip_address"`
	ActivityType string                 `json:"activity_type"`
	Description  string                 `json:"description"`
	RiskScore    int                    `json:"risk_score"`
	DetectedAt   time.Time              `json:"detected_at"`
	Metadata     map[string]interface{} `json:"metadata"`
	Resolved     bool                   `json:"resolved"`
	ResolvedAt   *time.Time             `json:"resolved_at,omitempty"`
}

// NewAuthManager creates a new authentication manager
func NewAuthManager(
	logger *observability.Logger,
	jwtService *auth.JWTService,
	mfaService *auth.MFAService,
	rbacService *auth.RBACService,
	config *SecurityConfig,
) *AuthManager {
	if config == nil {
		config = getDefaultSecurityConfig()
	}

	am := &AuthManager{
		logger:               logger,
		jwtService:           jwtService,
		mfaService:           mfaService,
		rbacService:          rbacService,
		securityConfig:       config,
		activeSessions:       make(map[string]*SecuritySession),
		activeAPIKeys:        make(map[string]*APIKey),
		blockedIPs:           make(map[string]*IPBlock),
		suspiciousActivities: make(map[string]*SuspiciousActivity),
	}

	// Initialize security components
	am.sessionManager = NewSessionManager(logger, config)
	am.apiKeyManager = NewAPIKeyManager(logger, config)
	am.rateLimiter = NewRateLimiter(config)
	am.bruteForceProtection = NewBruteForceProtection(config)
	am.deviceTrustManager = NewDeviceTrustManager(logger, config)
	am.behaviorAnalyzer = NewBehaviorAnalyzer(logger, config)
	am.threatDetector = NewThreatDetector(logger, config)

	return am
}

// AuthenticateUser authenticates a user with comprehensive security checks
func (am *AuthManager) AuthenticateUser(ctx context.Context, req *AuthenticationRequest) (*AuthenticationResponse, error) {
	am.mu.Lock()
	defer am.mu.Unlock()

	// 1. Rate limiting check
	if !am.rateLimiter.Allow(req.IPAddress, "login") {
		am.logSecurityEvent(ctx, "rate_limit_exceeded", req.IPAddress, nil, map[string]interface{}{
			"type": "login_attempt",
		})
		return nil, fmt.Errorf("rate limit exceeded")
	}

	// 2. IP blocking check
	if am.isIPBlocked(req.IPAddress) {
		am.logSecurityEvent(ctx, "blocked_ip_access", req.IPAddress, nil, map[string]interface{}{
			"ip": req.IPAddress,
		})
		return nil, fmt.Errorf("IP address is blocked")
	}

	// 3. Brute force protection check
	if am.bruteForceProtection.IsBlocked(req.IPAddress, req.Username) {
		am.logSecurityEvent(ctx, "brute_force_blocked", req.IPAddress, nil, map[string]interface{}{
			"username": req.Username,
		})
		return nil, fmt.Errorf("account temporarily locked due to failed attempts")
	}

	// 4. Threat detection
	if am.securityConfig.EnableThreatDetection {
		threatLevel := am.threatDetector.AnalyzeRequest(ctx, req)
		if threatLevel > 80 {
			am.logSecurityEvent(ctx, "high_threat_detected", req.IPAddress, nil, map[string]interface{}{
				"threat_level": threatLevel,
				"username":     req.Username,
			})
			return nil, fmt.Errorf("suspicious activity detected")
		}
	}

	// 5. Validate credentials (this would integrate with your user service)
	user, err := am.validateCredentials(ctx, req.Username, req.Password)
	if err != nil {
		// Record failed attempt
		am.bruteForceProtection.RecordFailedAttempt(req.IPAddress, req.Username)

		am.logSecurityEvent(ctx, "authentication_failed", req.IPAddress, nil, map[string]interface{}{
			"username": req.Username,
			"reason":   "invalid_credentials",
		})
		return nil, fmt.Errorf("invalid credentials")
	}

	// 6. Device trust evaluation
	deviceTrusted := false
	if am.securityConfig.RequireDeviceRegistration {
		deviceTrusted = am.deviceTrustManager.IsDeviceTrusted(user.ID, req.DeviceFingerprint)
		if !deviceTrusted && am.securityConfig.RequireDeviceRegistration {
			// Initiate device registration process
			return &AuthenticationResponse{
				Success:                 false,
				RequireDeviceSetup:      true,
				DeviceRegistrationToken: am.generateDeviceRegistrationToken(user.ID, req.DeviceFingerprint),
			}, nil
		}
	}

	// 7. MFA check
	mfaRequired := am.securityConfig.RequireMFA || user.MFAEnabled
	if mfaRequired && !req.MFACode.Valid {
		// Generate MFA challenge (simplified for now)
		challenge := &auth.MFAChallenge{
			ID:        uuid.New(),
			UserID:    user.ID,
			Method:    "totp",
			ExpiresAt: time.Now().Add(5 * time.Minute),
		}

		return &AuthenticationResponse{
			Success:      false,
			RequireMFA:   true,
			MFAChallenge: challenge,
		}, nil
	}

	// 8. Verify MFA if provided (simplified verification)
	if mfaRequired && req.MFACode.Valid {
		// In a real implementation, this would verify the TOTP code
		if req.MFACode.String == "" || len(req.MFACode.String) != 6 {
			am.bruteForceProtection.RecordFailedAttempt(req.IPAddress, req.Username)
			return nil, fmt.Errorf("invalid MFA code")
		}
	}

	// 9. Geo-location verification
	geoLocation := am.getGeoLocation(req.IPAddress)
	if am.securityConfig.RequireGeoVerification {
		if err := am.verifyGeoLocation(user.ID, geoLocation); err != nil {
			am.logSecurityEvent(ctx, "geo_verification_failed", req.IPAddress, &user.ID, map[string]interface{}{
				"location": geoLocation,
			})
			return nil, fmt.Errorf("geo-location verification failed")
		}
	}

	// 10. Create secure session
	session, err := am.createSecureSession(ctx, user, req, deviceTrusted, geoLocation)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// 11. Generate tokens
	authUser := &auth.User{
		ID:          user.ID,
		Email:       user.Email,
		Role:        user.Role,
		Permissions: user.Permissions,
		TeamID:      user.TeamID,
		MFAEnabled:  user.MFAEnabled,
		MFAVerified: user.MFAVerified,
		IsActive:    user.IsActive,
		CreatedAt:   user.CreatedAt,
	}
	tokens, err := am.jwtService.GenerateTokenPair(authUser, session.ID, req.IPAddress, req.UserAgent, req.DeviceFingerprint, []string{"trading"})
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// 12. Clear failed attempts on successful login
	am.bruteForceProtection.ClearFailedAttempts(req.IPAddress, req.Username)

	// 13. Log successful authentication
	am.logSecurityEvent(ctx, "authentication_success", req.IPAddress, &user.ID, map[string]interface{}{
		"session_id":     session.ID,
		"device_trusted": deviceTrusted,
		"mfa_verified":   mfaRequired,
	})

	return &AuthenticationResponse{
		Success:      true,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    int(tokens.ExpiresIn),
		Session:      session,
		User:         user,
	}, nil
}

// AuthenticateAPIKey authenticates an API key request
func (am *AuthManager) AuthenticateAPIKey(ctx context.Context, keyString, ipAddress string) (*APIKey, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	// Rate limiting for API key requests
	if !am.rateLimiter.Allow(ipAddress, "api") {
		return nil, fmt.Errorf("API rate limit exceeded")
	}

	// Validate API key format and extract key ID
	keyID, keySecret, err := am.parseAPIKey(keyString)
	if err != nil {
		return nil, fmt.Errorf("invalid API key format")
	}

	// Get API key from storage
	apiKey, exists := am.activeAPIKeys[keyID]
	if !exists {
		am.logSecurityEvent(ctx, "invalid_api_key", ipAddress, nil, map[string]interface{}{
			"key_id": keyID,
		})
		return nil, fmt.Errorf("invalid API key")
	}

	// Verify API key secret
	if !am.verifyAPIKeySecret(apiKey.KeyHash, keySecret) {
		am.logSecurityEvent(ctx, "api_key_verification_failed", ipAddress, &apiKey.UserID, map[string]interface{}{
			"key_id": keyID,
		})
		return nil, fmt.Errorf("API key verification failed")
	}

	// Check if API key is active and not expired
	if !apiKey.IsActive {
		return nil, fmt.Errorf("API key is disabled")
	}

	if apiKey.ExpiresAt != nil && time.Now().After(*apiKey.ExpiresAt) {
		return nil, fmt.Errorf("API key has expired")
	}

	// Check IP whitelist
	if len(apiKey.IPWhitelist) > 0 && !am.isIPWhitelisted(ipAddress, apiKey.IPWhitelist) {
		am.logSecurityEvent(ctx, "api_key_ip_not_whitelisted", ipAddress, &apiKey.UserID, map[string]interface{}{
			"key_id": keyID,
		})
		return nil, fmt.Errorf("IP address not whitelisted for this API key")
	}

	// Update last used timestamp
	apiKey.LastUsed = time.Now()

	am.logSecurityEvent(ctx, "api_key_authentication_success", ipAddress, &apiKey.UserID, map[string]interface{}{
		"key_id": keyID,
	})

	return apiKey, nil
}

// GetAPIKeyManager returns the API key manager
func (am *AuthManager) GetAPIKeyManager() *APIKeyManager {
	return am.apiKeyManager
}

// GetSessionManager returns the session manager
func (am *AuthManager) GetSessionManager() *SessionManager {
	return am.sessionManager
}

// GetRateLimiter returns the rate limiter
func (am *AuthManager) GetRateLimiter() *RateLimiter {
	return am.rateLimiter
}

// GetBruteForceProtection returns the brute force protection
func (am *AuthManager) GetBruteForceProtection() *BruteForceProtection {
	return am.bruteForceProtection
}

// GetDeviceTrustManager returns the device trust manager
func (am *AuthManager) GetDeviceTrustManager() *DeviceTrustManager {
	return am.deviceTrustManager
}

// GetBehaviorAnalyzer returns the behavior analyzer
func (am *AuthManager) GetBehaviorAnalyzer() *BehaviorAnalyzer {
	return am.behaviorAnalyzer
}

// GetThreatDetector returns the threat detector
func (am *AuthManager) GetThreatDetector() *ThreatDetector {
	return am.threatDetector
}

// GetSecurityConfig returns the security configuration
func (am *AuthManager) GetSecurityConfig() *SecurityConfig {
	return am.securityConfig
}

// GetActiveSessions returns all active sessions
func (am *AuthManager) GetActiveSessions() map[string]*SecuritySession {
	am.mu.RLock()
	defer am.mu.RUnlock()

	// Return a copy to avoid concurrent access issues
	result := make(map[string]*SecuritySession)
	for k, v := range am.activeSessions {
		result[k] = v
	}

	return result
}

// GetActiveAPIKeys returns all active API keys
func (am *AuthManager) GetActiveAPIKeys() map[string]*APIKey {
	am.mu.RLock()
	defer am.mu.RUnlock()

	// Return a copy to avoid concurrent access issues
	result := make(map[string]*APIKey)
	for k, v := range am.activeAPIKeys {
		result[k] = v
	}

	return result
}

// GetBlockedIPs returns all blocked IP addresses
func (am *AuthManager) GetBlockedIPs() map[string]*IPBlock {
	am.mu.RLock()
	defer am.mu.RUnlock()

	// Return a copy to avoid concurrent access issues
	result := make(map[string]*IPBlock)
	for k, v := range am.blockedIPs {
		result[k] = v
	}

	return result
}

// GetSuspiciousActivities returns all suspicious activities
func (am *AuthManager) GetSuspiciousActivities() map[string]*SuspiciousActivity {
	am.mu.RLock()
	defer am.mu.RUnlock()

	// Return a copy to avoid concurrent access issues
	result := make(map[string]*SuspiciousActivity)
	for k, v := range am.suspiciousActivities {
		result[k] = v
	}

	return result
}

// BlockIP blocks an IP address
func (am *AuthManager) BlockIP(ipAddress, reason string, duration time.Duration) {
	am.mu.Lock()
	defer am.mu.Unlock()

	block := &IPBlock{
		IPAddress:   ipAddress,
		Reason:      reason,
		BlockedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(duration),
		Attempts:    1,
		Severity:    "medium",
		IsAutomatic: true,
	}

	am.blockedIPs[ipAddress] = block

	am.logger.Info(context.Background(), "IP address blocked", map[string]interface{}{
		"ip_address": ipAddress,
		"reason":     reason,
		"duration":   duration.String(),
	})
}

// UnblockIP unblocks an IP address
func (am *AuthManager) UnblockIP(ipAddress string) {
	am.mu.Lock()
	defer am.mu.Unlock()

	delete(am.blockedIPs, ipAddress)

	am.logger.Info(context.Background(), "IP address unblocked", map[string]interface{}{
		"ip_address": ipAddress,
	})
}

// RecordSuspiciousActivity records a suspicious activity
func (am *AuthManager) RecordSuspiciousActivity(userID *uuid.UUID, ipAddress, activityType, description string, riskScore int, metadata map[string]interface{}) {
	am.mu.Lock()
	defer am.mu.Unlock()

	activity := &SuspiciousActivity{
		ID:           uuid.New().String(),
		UserID:       userID,
		IPAddress:    ipAddress,
		ActivityType: activityType,
		Description:  description,
		RiskScore:    riskScore,
		DetectedAt:   time.Now(),
		Metadata:     metadata,
		Resolved:     false,
	}

	am.suspiciousActivities[activity.ID] = activity

	logData := map[string]interface{}{
		"activity_id":   activity.ID,
		"ip_address":    ipAddress,
		"activity_type": activityType,
		"risk_score":    riskScore,
	}

	if userID != nil {
		logData["user_id"] = userID.String()
	}

	am.logger.Warn(context.Background(), "Suspicious activity detected", logData)
}

// ValidateSession validates a session and returns it if valid
func (am *AuthManager) ValidateSession(sessionID string) (*SecuritySession, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	session, exists := am.activeSessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found")
	}

	// Check if session is active and not expired
	if !session.IsActive {
		return nil, fmt.Errorf("session is not active")
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, fmt.Errorf("session has expired")
	}

	// Update last activity
	session.LastActivity = time.Now()

	return session, nil
}

// RevokeSession revokes a session
func (am *AuthManager) RevokeSession(sessionID string) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	session, exists := am.activeSessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found")
	}

	session.IsActive = false

	am.logger.Info(context.Background(), "Session revoked", map[string]interface{}{
		"session_id": sessionID,
		"user_id":    session.UserID.String(),
	})

	return nil
}

// RevokeAllUserSessions revokes all sessions for a user
func (am *AuthManager) RevokeAllUserSessions(userID uuid.UUID) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	revokedCount := 0
	for _, session := range am.activeSessions {
		if session.UserID == userID && session.IsActive {
			session.IsActive = false
			revokedCount++
		}
	}

	am.logger.Info(context.Background(), "All user sessions revoked", map[string]interface{}{
		"user_id":       userID.String(),
		"revoked_count": revokedCount,
	})

	return nil
}

// CleanupExpiredSessions removes expired sessions
func (am *AuthManager) CleanupExpiredSessions() {
	am.mu.Lock()
	defer am.mu.Unlock()

	now := time.Now()
	cleanedCount := 0

	for sessionID, session := range am.activeSessions {
		if now.After(session.ExpiresAt) {
			delete(am.activeSessions, sessionID)
			cleanedCount++
		}
	}

	if cleanedCount > 0 {
		am.logger.Info(context.Background(), "Expired sessions cleaned up", map[string]interface{}{
			"cleaned_count": cleanedCount,
		})
	}
}

// GetSecurityMetrics returns security metrics
func (am *AuthManager) GetSecurityMetrics() map[string]interface{} {
	am.mu.RLock()
	defer am.mu.RUnlock()

	activeSessionCount := 0
	for _, session := range am.activeSessions {
		if session.IsActive && time.Now().Before(session.ExpiresAt) {
			activeSessionCount++
		}
	}

	activeAPIKeyCount := 0
	for _, apiKey := range am.activeAPIKeys {
		if apiKey.IsActive {
			activeAPIKeyCount++
		}
	}

	return map[string]interface{}{
		"active_sessions":       activeSessionCount,
		"total_sessions":        len(am.activeSessions),
		"active_api_keys":       activeAPIKeyCount,
		"total_api_keys":        len(am.activeAPIKeys),
		"blocked_ips":           len(am.blockedIPs),
		"suspicious_activities": len(am.suspiciousActivities),
		"timestamp":             time.Now(),
	}
}

// AuthenticationRequest represents an authentication request
type AuthenticationRequest struct {
	Username          string                 `json:"username"`
	Password          string                 `json:"password"`
	MFACode           NullString             `json:"mfa_code,omitempty"`
	IPAddress         string                 `json:"ip_address"`
	UserAgent         string                 `json:"user_agent"`
	DeviceFingerprint string                 `json:"device_fingerprint"`
	RememberDevice    bool                   `json:"remember_device"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

// AuthenticationResponse represents an authentication response
type AuthenticationResponse struct {
	Success                 bool                   `json:"success"`
	AccessToken             string                 `json:"access_token,omitempty"`
	RefreshToken            string                 `json:"refresh_token,omitempty"`
	ExpiresIn               int                    `json:"expires_in,omitempty"`
	Session                 *SecuritySession       `json:"session,omitempty"`
	User                    *User                  `json:"user,omitempty"`
	RequireMFA              bool                   `json:"require_mfa,omitempty"`
	MFAChallenge            *auth.MFAChallenge     `json:"mfa_challenge,omitempty"`
	RequireDeviceSetup      bool                   `json:"require_device_setup,omitempty"`
	DeviceRegistrationToken string                 `json:"device_registration_token,omitempty"`
	SecurityWarnings        []string               `json:"security_warnings,omitempty"`
	Metadata                map[string]interface{} `json:"metadata,omitempty"`
}

// User represents a user for authentication purposes
type User struct {
	ID          uuid.UUID  `json:"id"`
	Email       string     `json:"email"`
	Username    string     `json:"username"`
	Role        string     `json:"role"`
	Permissions []string   `json:"permissions"`
	TeamID      *uuid.UUID `json:"team_id,omitempty"`
	MFAEnabled  bool       `json:"mfa_enabled"`
	MFAVerified bool       `json:"mfa_verified"`
	IsActive    bool       `json:"is_active"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// NullString represents a nullable string
type NullString struct {
	String string
	Valid  bool
}

// Helper methods implementation

// validateCredentials validates user credentials (placeholder - integrate with your user service)
func (am *AuthManager) validateCredentials(ctx context.Context, username, password string) (*User, error) {
	// This would integrate with your actual user service/database
	// For now, return a mock user for demonstration
	return &User{
		ID:          uuid.New(),
		Email:       username,
		Username:    username,
		Role:        "user",
		Permissions: []string{"trading:read", "trading:write"},
		MFAEnabled:  true,
		MFAVerified: false,
		IsActive:    true,
		CreatedAt:   time.Now(),
	}, nil
}

// createSecureSession creates a new secure session
func (am *AuthManager) createSecureSession(ctx context.Context, user *User, req *AuthenticationRequest, deviceTrusted bool, geoLocation *GeoLocation) (*SecuritySession, error) {
	sessionID := uuid.New().String()
	now := time.Now()

	// Calculate risk score based on various factors
	riskScore := am.calculateRiskScore(req, deviceTrusted, geoLocation)

	session := &SecuritySession{
		ID:             sessionID,
		UserID:         user.ID,
		DeviceID:       req.DeviceFingerprint,
		IPAddress:      req.IPAddress,
		UserAgent:      req.UserAgent,
		CreatedAt:      now,
		LastActivity:   now,
		ExpiresAt:      now.Add(am.securityConfig.SessionTimeout),
		IsActive:       true,
		MFAVerified:    am.securityConfig.RequireMFA,
		DeviceTrusted:  deviceTrusted,
		RiskScore:      riskScore,
		Permissions:    user.Permissions,
		TradingEnabled: am.isTradingEnabled(user, riskScore),
		MaxTradeAmount: am.calculateMaxTradeAmount(user, riskScore),
		AllowedPairs:   am.getAllowedTradingPairs(user),
		SecurityFlags:  make(map[string]interface{}),
		GeoLocation:    geoLocation,
	}

	// Set security flags based on risk assessment
	if riskScore > 50 {
		session.SecurityFlags["high_risk"] = true
		session.SecurityFlags["require_additional_verification"] = true
	}

	if !deviceTrusted {
		session.SecurityFlags["untrusted_device"] = true
		session.TradingEnabled = false
	}

	// Store session
	am.activeSessions[sessionID] = session

	return session, nil
}

// calculateRiskScore calculates a risk score for the authentication request
func (am *AuthManager) calculateRiskScore(req *AuthenticationRequest, deviceTrusted bool, geoLocation *GeoLocation) int {
	score := 0

	// Device trust factor
	if !deviceTrusted {
		score += 30
	}

	// Geo-location factors
	if geoLocation != nil {
		if geoLocation.IsVPN {
			score += 20
		}
		if geoLocation.IsTor {
			score += 40
		}
		score += geoLocation.ThreatLevel
	}

	// Time-based factors (unusual login times)
	hour := time.Now().Hour()
	if hour < 6 || hour > 22 {
		score += 10
	}

	// Ensure score is within bounds
	if score > 100 {
		score = 100
	}

	return score
}

// isIPBlocked checks if an IP address is blocked
func (am *AuthManager) isIPBlocked(ipAddress string) bool {
	block, exists := am.blockedIPs[ipAddress]
	if !exists {
		return false
	}

	// Check if block has expired
	if time.Now().After(block.ExpiresAt) {
		delete(am.blockedIPs, ipAddress)
		return false
	}

	return true
}

// parseAPIKey parses an API key string into ID and secret
func (am *AuthManager) parseAPIKey(keyString string) (string, string, error) {
	parts := strings.Split(keyString, ".")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid API key format")
	}

	keyID := parts[0]
	keySecret := parts[1]

	if len(keyID) == 0 || len(keySecret) == 0 {
		return "", "", fmt.Errorf("invalid API key components")
	}

	return keyID, keySecret, nil
}

// verifyAPIKeySecret verifies an API key secret against its hash
func (am *AuthManager) verifyAPIKeySecret(keyHash, keySecret string) bool {
	// Decode the stored hash
	storedHash, err := base64.StdEncoding.DecodeString(keyHash)
	if err != nil {
		return false
	}

	// Extract salt and hash from stored hash
	if len(storedHash) < 32 {
		return false
	}

	salt := storedHash[:16]
	hash := storedHash[16:]

	// Hash the provided secret with the same salt
	providedHash := argon2.IDKey([]byte(keySecret), salt, 1, 64*1024, 4, 32)

	// Compare hashes using constant-time comparison
	return subtle.ConstantTimeCompare(hash, providedHash) == 1
}

// isIPWhitelisted checks if an IP is in the whitelist
func (am *AuthManager) isIPWhitelisted(ipAddress string, whitelist []string) bool {
	for _, whitelistedIP := range whitelist {
		if ipAddress == whitelistedIP {
			return true
		}
		// Could add CIDR range checking here
	}
	return false
}

// getGeoLocation gets geographical location for an IP address
func (am *AuthManager) getGeoLocation(ipAddress string) *GeoLocation {
	// This would integrate with a geo-location service
	// For now, return a mock location
	return &GeoLocation{
		Country:     "US",
		Region:      "CA",
		City:        "San Francisco",
		Latitude:    37.7749,
		Longitude:   -122.4194,
		ISP:         "Example ISP",
		Timezone:    "America/Los_Angeles",
		IsVPN:       false,
		IsTor:       false,
		ThreatLevel: 0,
	}
}

// verifyGeoLocation verifies if the geo-location is acceptable
func (am *AuthManager) verifyGeoLocation(userID uuid.UUID, location *GeoLocation) error {
	// Implement geo-location verification logic
	// For example, check against user's known locations
	if location.ThreatLevel > 50 {
		return fmt.Errorf("high threat level location")
	}

	if location.IsTor {
		return fmt.Errorf("Tor access not allowed")
	}

	return nil
}

// generateDeviceRegistrationToken generates a token for device registration
func (am *AuthManager) generateDeviceRegistrationToken(userID uuid.UUID, deviceFingerprint string) string {
	// Generate a secure random token
	tokenBytes := make([]byte, 32)
	rand.Read(tokenBytes)
	return base64.URLEncoding.EncodeToString(tokenBytes)
}

// isTradingEnabled determines if trading should be enabled for a user/session
func (am *AuthManager) isTradingEnabled(user *User, riskScore int) bool {
	// Disable trading for high-risk sessions
	if riskScore > 70 {
		return false
	}

	// Check user permissions
	for _, perm := range user.Permissions {
		if perm == "trading:write" {
			return true
		}
	}

	return false
}

// calculateMaxTradeAmount calculates maximum trade amount based on risk
func (am *AuthManager) calculateMaxTradeAmount(user *User, riskScore int) decimal.Decimal {
	baseAmount := decimal.NewFromFloat(10000) // $10,000 base limit

	// Reduce limit based on risk score
	if riskScore > 50 {
		reduction := decimal.NewFromFloat(float64(riskScore-50) / 100)
		baseAmount = baseAmount.Mul(decimal.NewFromFloat(1).Sub(reduction))
	}

	return baseAmount
}

// getAllowedTradingPairs gets allowed trading pairs for a user
func (am *AuthManager) getAllowedTradingPairs(user *User) []string {
	// Default allowed pairs
	return []string{"BTC/USDT", "ETH/USDT", "BNB/USDT"}
}

// logSecurityEvent logs a security event for audit purposes
func (am *AuthManager) logSecurityEvent(ctx context.Context, eventType, ipAddress string, userID *uuid.UUID, metadata map[string]interface{}) {
	if !am.securityConfig.EnableAuditLogging {
		return
	}

	logData := map[string]interface{}{
		"event_type": eventType,
		"ip_address": ipAddress,
		"timestamp":  time.Now(),
		"metadata":   metadata,
	}

	if userID != nil {
		logData["user_id"] = userID.String()
	}

	am.logger.Info(ctx, "Security event", logData)
}

// getDefaultSecurityConfig returns default security configuration
func getDefaultSecurityConfig() *SecurityConfig {
	return &SecurityConfig{
		RequireMFA:                true,
		SessionTimeout:            24 * time.Hour,
		MaxConcurrentSessions:     5,
		RequireStrongPasswords:    true,
		PasswordMinLength:         12,
		PasswordRequireUppercase:  true,
		PasswordRequireLowercase:  true,
		PasswordRequireNumbers:    true,
		PasswordRequireSymbols:    true,
		LoginRateLimit:            5,
		APIRateLimit:              1000,
		RateLimitWindow:           time.Minute,
		MaxLoginAttempts:          5,
		LockoutDuration:           15 * time.Minute,
		ProgressiveLockout:        true,
		RequireDeviceRegistration: true,
		DeviceTrustDuration:       30 * 24 * time.Hour,
		MaxTrustedDevices:         5,
		RequireAPIKeyAuth:         true,
		APIKeyRotationInterval:    90 * 24 * time.Hour,
		EnableBehaviorAnalysis:    true,
		EnableThreatDetection:     true,
		EnableZeroTrust:           true,
		RequireGeoVerification:    false,
		EnableAuditLogging:        true,
		RetainAuditLogs:           365 * 24 * time.Hour,
		EnableComplianceMode:      true,
	}
}
