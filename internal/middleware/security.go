package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ai-agentic-browser/internal/security"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

// SecurityMiddleware provides comprehensive security middleware
type SecurityMiddleware struct {
	logger           *observability.Logger
	zeroTrustEngine  *security.ZeroTrustEngine
	threatDetector   *security.AdvancedThreatDetector
	policyEngine     *security.PolicyEngine
	deviceManager    *security.DeviceTrustManager
	behaviorAnalyzer *security.BehaviorAnalyzer
	config           *SecurityConfig
}

// SecurityConfig contains security middleware configuration
type SecurityConfig struct {
	EnableZeroTrust         bool
	EnableThreatDetection   bool
	EnablePolicyEngine      bool
	EnableDeviceTrust       bool
	EnableBehaviorAnalysis  bool
	BlockSuspiciousRequests bool
	LogSecurityEvents       bool
	RequireAuthentication   bool
	AllowedOrigins          []string
	RateLimitEnabled        bool
	CSRFProtection          bool
}

// SecurityContext contains security information for a request
type SecurityContext struct {
	UserID         uuid.UUID
	DeviceID       string
	IPAddress      string
	UserAgent      string
	RiskScore      float64
	DeviceTrust    float64
	ThreatDetected bool
	ThreatScore    float64
	PolicyDecision *security.PolicyDecision
	SessionValid   bool
	RequiresMFA    bool
}

// NewSecurityMiddleware creates a new security middleware
func NewSecurityMiddleware(logger *observability.Logger) *SecurityMiddleware {
	config := &SecurityConfig{
		EnableZeroTrust:         true,
		EnableThreatDetection:   true,
		EnablePolicyEngine:      true,
		EnableDeviceTrust:       true,
		EnableBehaviorAnalysis:  true,
		BlockSuspiciousRequests: true,
		LogSecurityEvents:       true,
		RequireAuthentication:   true,
		AllowedOrigins:          []string{"*"},
		RateLimitEnabled:        true,
		CSRFProtection:          true,
	}

	return &SecurityMiddleware{
		logger:           logger,
		zeroTrustEngine:  security.NewZeroTrustEngine(logger),
		threatDetector:   security.NewAdvancedThreatDetector(logger),
		policyEngine:     security.NewPolicyEngine(logger),
		deviceManager:    security.NewDeviceTrustManager(logger),
		behaviorAnalyzer: security.NewBehaviorAnalyzer(logger),
		config:           config,
	}
}

// SecurityHandler returns the main security middleware handler
func (sm *SecurityMiddleware) SecurityHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			startTime := time.Now()

			// Create security context
			secCtx, err := sm.createSecurityContext(ctx, r)
			if err != nil {
				sm.handleSecurityError(w, r, "Failed to create security context", err)
				return
			}

			// Add security headers
			sm.addSecurityHeaders(w, r)

			// Perform security checks
			if !sm.performSecurityChecks(ctx, w, r, secCtx) {
				return // Request blocked by security checks
			}

			// Add security context to request context
			ctx = context.WithValue(ctx, "security_context", secCtx)
			r = r.WithContext(ctx)

			// Log security event
			if sm.config.LogSecurityEvents {
				sm.logSecurityEvent(ctx, r, secCtx, time.Since(startTime))
			}

			// Continue to next handler
			next.ServeHTTP(w, r)
		})
	}
}

// createSecurityContext creates a security context for the request
func (sm *SecurityMiddleware) createSecurityContext(ctx context.Context, r *http.Request) (*SecurityContext, error) {
	secCtx := &SecurityContext{
		IPAddress: sm.getClientIP(r),
		UserAgent: r.UserAgent(),
	}

	// Extract user ID from authentication context (if available)
	if userID := sm.extractUserID(r); userID != uuid.Nil {
		secCtx.UserID = userID
	}

	// Generate or extract device ID
	secCtx.DeviceID = sm.generateDeviceFingerprint(r)

	return secCtx, nil
}

// performSecurityChecks performs comprehensive security checks
func (sm *SecurityMiddleware) performSecurityChecks(ctx context.Context, w http.ResponseWriter, r *http.Request, secCtx *SecurityContext) bool {
	// 1. Threat Detection
	if sm.config.EnableThreatDetection {
		if !sm.performThreatDetection(ctx, w, r, secCtx) {
			return false
		}
	}

	// 2. Device Trust Evaluation
	if sm.config.EnableDeviceTrust {
		if !sm.evaluateDeviceTrust(ctx, w, r, secCtx) {
			return false
		}
	}

	// 3. Behavioral Analysis
	if sm.config.EnableBehaviorAnalysis {
		if !sm.performBehaviorAnalysis(ctx, w, r, secCtx) {
			return false
		}
	}

	// 4. Zero Trust Evaluation
	if sm.config.EnableZeroTrust {
		if !sm.performZeroTrustEvaluation(ctx, w, r, secCtx) {
			return false
		}
	}

	// 5. Policy Engine Evaluation
	if sm.config.EnablePolicyEngine {
		if !sm.evaluatePolicies(ctx, w, r, secCtx) {
			return false
		}
	}

	return true
}

// performThreatDetection performs threat detection on the request
func (sm *SecurityMiddleware) performThreatDetection(ctx context.Context, w http.ResponseWriter, r *http.Request, secCtx *SecurityContext) bool {
	request := &security.SecurityRequest{
		RequestID: uuid.New().String(),
		UserID:    &secCtx.UserID,
		IPAddress: secCtx.IPAddress,
		UserAgent: secCtx.UserAgent,
		Method:    r.Method,
		URL:       r.URL.String(),
		Headers:   sm.extractHeaders(r),
		Body:      sm.extractBody(r),
		Timestamp: time.Now(),
	}

	result, err := sm.threatDetector.DetectThreats(ctx, request)
	if err != nil {
		sm.logger.Error(ctx, "Threat detection failed", err)
		return true // Continue on error (fail open)
	}

	secCtx.ThreatDetected = result.ThreatDetected
	secCtx.ThreatScore = result.ThreatScore

	// Block request if threat detected and blocking is enabled
	if result.ShouldBlock && sm.config.BlockSuspiciousRequests {
		sm.handleThreatBlocked(w, r, result)
		return false
	}

	return true
}

// evaluateDeviceTrust evaluates device trust
func (sm *SecurityMiddleware) evaluateDeviceTrust(ctx context.Context, w http.ResponseWriter, r *http.Request, secCtx *SecurityContext) bool {
	device, err := sm.deviceManager.GetDevice(secCtx.DeviceID)
	if err != nil {
		// New device - register with low trust
		// Convert map[string]string to map[string]interface{}
		attributes := make(map[string]interface{})
		for k, v := range sm.extractDeviceAttributes(r) {
			attributes[k] = v
		}

		device = &security.TrustedDevice{
			DeviceID:    secCtx.DeviceID,
			UserID:      secCtx.UserID,
			TrustLevel:  0.3, // Low initial trust
			Attributes:  attributes,
			RiskFactors: []string{"new_device"},
			LastSeen:    time.Now(),
		}
		sm.deviceManager.RegisterDevice(device)
	} else {
		// Update last seen
		device.LastSeen = time.Now()
		sm.deviceManager.UpdateDevice(device)
	}

	secCtx.DeviceTrust = device.TrustLevel

	// Require additional verification for low trust devices
	if device.TrustLevel < 0.5 {
		secCtx.RequiresMFA = true
	}

	return true
}

// performBehaviorAnalysis performs behavioral analysis
func (sm *SecurityMiddleware) performBehaviorAnalysis(ctx context.Context, w http.ResponseWriter, r *http.Request, secCtx *SecurityContext) bool {
	request := &security.SecurityRequest{
		RequestID: uuid.New().String(),
		UserID:    &secCtx.UserID,
		IPAddress: secCtx.IPAddress,
		UserAgent: secCtx.UserAgent,
		Method:    r.Method,
		URL:       r.URL.String(),
		Timestamp: time.Now(),
	}

	// Create user behavior profile for analysis
	profile := &security.UserBehaviorProfile{
		UserID: secCtx.UserID,
	}

	// Convert SecurityRequest to AccessRequest
	accessRequest := &security.AccessRequest{
		UserID:    &secCtx.UserID,
		DeviceID:  secCtx.DeviceID,
		IPAddress: secCtx.IPAddress,
		UserAgent: secCtx.UserAgent,
		Resource:  request.URL,
		Action:    request.Method,
		Timestamp: request.Timestamp,
	}

	riskScore := sm.behaviorAnalyzer.AnalyzeBehavior(profile, accessRequest)
	if riskScore < 0 {
		sm.logger.Error(ctx, "Behavior analysis failed", nil)
		return true // Continue on error
	}

	// Update risk score based on behavioral analysis
	secCtx.RiskScore = riskScore

	return true
}

// performZeroTrustEvaluation performs zero trust evaluation
func (sm *SecurityMiddleware) performZeroTrustEvaluation(ctx context.Context, w http.ResponseWriter, r *http.Request, secCtx *SecurityContext) bool {
	accessRequest := &security.AccessRequest{
		UserID:    &secCtx.UserID,
		DeviceID:  secCtx.DeviceID,
		IPAddress: secCtx.IPAddress,
		UserAgent: secCtx.UserAgent,
		Resource:  r.URL.Path,
		Action:    r.Method,
		Timestamp: time.Now(),
	}

	decision, err := sm.zeroTrustEngine.EvaluateAccess(ctx, accessRequest)
	if err != nil {
		sm.logger.Error(ctx, "Zero trust evaluation failed", err)
		return true // Continue on error
	}

	// Convert AccessDecision to PolicyDecision
	policyDecision := &security.PolicyDecision{
		Allowed:     decision.Allowed,
		RequiresMFA: decision.RequiresMFA,
		SessionTTL:  decision.SessionTTL,
		Reason:      decision.Reason,
	}
	secCtx.PolicyDecision = policyDecision

	// Block access if not allowed
	if !decision.Allowed {
		sm.handleAccessDenied(w, r, policyDecision)
		return false
	}

	// Set MFA requirement
	if decision.RequiresMFA {
		secCtx.RequiresMFA = true
	}

	return true
}

// evaluatePolicies evaluates security policies
func (sm *SecurityMiddleware) evaluatePolicies(ctx context.Context, w http.ResponseWriter, r *http.Request, secCtx *SecurityContext) bool {
	evalCtx := &security.PolicyEvaluationContext{
		UserID:      secCtx.UserID,
		UserRoles:   sm.extractUserRoles(r),
		IPAddress:   secCtx.IPAddress,
		UserAgent:   secCtx.UserAgent,
		Resource:    r.URL.Path,
		Action:      r.Method,
		RiskScore:   secCtx.RiskScore,
		DeviceTrust: secCtx.DeviceTrust,
		Timestamp:   time.Now(),
	}

	result, err := sm.policyEngine.EvaluatePolicy(ctx, evalCtx)
	if err != nil {
		sm.logger.Error(ctx, "Policy evaluation failed", err)
		return true // Continue on error
	}

	// Block access if policy denies
	if !result.Decision.Allowed {
		sm.handlePolicyDenied(w, r, result)
		return false
	}

	return true
}

// addSecurityHeaders adds security headers to the response
func (sm *SecurityMiddleware) addSecurityHeaders(w http.ResponseWriter, r *http.Request) {
	// CORS headers
	if len(sm.config.AllowedOrigins) > 0 {
		origin := r.Header.Get("Origin")
		if sm.isAllowedOrigin(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
	}

	// Security headers
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	w.Header().Set("Content-Security-Policy", "default-src 'self'")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
	w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
}

// Helper methods

func (sm *SecurityMiddleware) getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}
	return ip
}

func (sm *SecurityMiddleware) extractUserID(r *http.Request) uuid.UUID {
	// Extract user ID from JWT token or session
	// Implementation would depend on authentication mechanism
	return uuid.Nil
}

func (sm *SecurityMiddleware) generateDeviceFingerprint(r *http.Request) string {
	// Generate device fingerprint based on headers and other attributes
	fingerprint := fmt.Sprintf("%s-%s", r.UserAgent(), r.Header.Get("Accept-Language"))
	return fmt.Sprintf("device-%x", []byte(fingerprint))
}

func (sm *SecurityMiddleware) extractHeaders(r *http.Request) map[string]string {
	headers := make(map[string]string)
	for name, values := range r.Header {
		if len(values) > 0 {
			headers[name] = values[0]
		}
	}
	return headers
}

func (sm *SecurityMiddleware) extractBody(r *http.Request) string {
	// Extract request body for analysis (with size limits)
	// Implementation would read and restore the body
	return ""
}

func (sm *SecurityMiddleware) extractDeviceAttributes(r *http.Request) map[string]string {
	return map[string]string{
		"user_agent":      r.UserAgent(),
		"accept_language": r.Header.Get("Accept-Language"),
		"accept_encoding": r.Header.Get("Accept-Encoding"),
	}
}

func (sm *SecurityMiddleware) extractUserRoles(r *http.Request) []string {
	// Extract user roles from authentication context
	return []string{"user"}
}

func (sm *SecurityMiddleware) isAllowedOrigin(origin string) bool {
	for _, allowed := range sm.config.AllowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
	}
	return false
}

// Error handlers

func (sm *SecurityMiddleware) handleSecurityError(w http.ResponseWriter, r *http.Request, message string, err error) {
	sm.logger.Error(r.Context(), message, err)
	http.Error(w, "Security check failed", http.StatusInternalServerError)
}

func (sm *SecurityMiddleware) handleThreatBlocked(w http.ResponseWriter, r *http.Request, result *security.ThreatDetectionResult) {
	sm.logger.Warn(r.Context(), "Request blocked due to threat detection", map[string]interface{}{
		"threat_type":  result.ThreatType,
		"threat_score": result.ThreatScore,
		"ip_address":   sm.getClientIP(r),
	})
	http.Error(w, "Request blocked for security reasons", http.StatusForbidden)
}

func (sm *SecurityMiddleware) handleAccessDenied(w http.ResponseWriter, r *http.Request, decision *security.PolicyDecision) {
	sm.logger.Warn(r.Context(), "Access denied by zero trust policy", map[string]interface{}{
		"reason":     decision.Reason,
		"ip_address": sm.getClientIP(r),
	})
	http.Error(w, "Access denied", http.StatusForbidden)
}

func (sm *SecurityMiddleware) handlePolicyDenied(w http.ResponseWriter, r *http.Request, result *security.PolicyEvaluationResult) {
	sm.logger.Warn(r.Context(), "Access denied by security policy", map[string]interface{}{
		"reason":     result.Reason,
		"ip_address": sm.getClientIP(r),
	})
	http.Error(w, "Access denied by policy", http.StatusForbidden)
}

func (sm *SecurityMiddleware) logSecurityEvent(ctx context.Context, r *http.Request, secCtx *SecurityContext, duration time.Duration) {
	sm.logger.Info(ctx, "Security evaluation completed", map[string]interface{}{
		"user_id":         secCtx.UserID,
		"device_id":       secCtx.DeviceID,
		"ip_address":      secCtx.IPAddress,
		"method":          r.Method,
		"url":             r.URL.String(),
		"risk_score":      secCtx.RiskScore,
		"device_trust":    secCtx.DeviceTrust,
		"threat_detected": secCtx.ThreatDetected,
		"threat_score":    secCtx.ThreatScore,
		"requires_mfa":    secCtx.RequiresMFA,
		"duration":        duration,
	})
}
