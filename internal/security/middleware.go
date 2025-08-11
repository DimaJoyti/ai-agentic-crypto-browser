package security

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ai-agentic-browser/internal/auth"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

// SecurityMiddleware provides comprehensive security middleware
type SecurityMiddleware struct {
	authManager *AuthManager
	jwtService  *auth.JWTService
	logger      *observability.Logger
	config      *SecurityConfig
}

// NewSecurityMiddleware creates a new security middleware
func NewSecurityMiddleware(
	authManager *AuthManager,
	jwtService *auth.JWTService,
	logger *observability.Logger,
	config *SecurityConfig,
) *SecurityMiddleware {
	return &SecurityMiddleware{
		authManager: authManager,
		jwtService:  jwtService,
		logger:      logger,
		config:      config,
	}
}

// AuthenticationMiddleware handles authentication for HTTP requests
func (sm *SecurityMiddleware) AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		
		// Skip authentication for public endpoints
		if sm.isPublicEndpoint(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Extract authentication credentials
		authType, credentials := sm.extractAuthCredentials(r)
		
		var userID uuid.UUID
		var permissions []string
		var securityContext *SecurityContext
		
		switch authType {
		case "Bearer":
			// JWT token authentication
			user, session, err := sm.authenticateJWT(ctx, credentials)
			if err != nil {
				sm.handleAuthError(w, r, "JWT authentication failed", err)
				return
			}
			userID = user.ID
			permissions = user.Permissions
			securityContext = &SecurityContext{
				UserID:      user.ID,
				SessionID:   session.ID,
				AuthType:    "jwt",
				Permissions: permissions,
				Session:     session,
				User:        user,
			}
			
		case "ApiKey":
			// API key authentication
			apiKey, err := sm.authenticateAPIKey(ctx, credentials, sm.getClientIP(r))
			if err != nil {
				sm.handleAuthError(w, r, "API key authentication failed", err)
				return
			}
			userID = apiKey.UserID
			permissions = apiKey.Permissions
			securityContext = &SecurityContext{
				UserID:      apiKey.UserID,
				AuthType:    "api_key",
				Permissions: permissions,
				APIKey:      apiKey,
			}
			
		default:
			sm.handleAuthError(w, r, "No valid authentication provided", fmt.Errorf("missing or invalid authentication"))
			return
		}

		// Add security context to request
		ctx = context.WithValue(ctx, "security_context", securityContext)
		ctx = context.WithValue(ctx, "user_id", userID)
		ctx = context.WithValue(ctx, "permissions", permissions)
		
		// Log successful authentication
		sm.logger.Info(ctx, "Request authenticated", map[string]interface{}{
			"user_id":    userID.String(),
			"auth_type":  authType,
			"endpoint":   r.URL.Path,
			"method":     r.Method,
			"ip_address": sm.getClientIP(r),
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthorizationMiddleware handles authorization for HTTP requests
func (sm *SecurityMiddleware) AuthorizationMiddleware(requiredPermissions []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			
			// Get security context
			securityContext, ok := ctx.Value("security_context").(*SecurityContext)
			if !ok {
				sm.handleAuthError(w, r, "Authorization failed", fmt.Errorf("no security context"))
				return
			}

			// Check permissions
			if !sm.hasRequiredPermissions(securityContext.Permissions, requiredPermissions) {
				sm.handleAuthError(w, r, "Insufficient permissions", fmt.Errorf("missing required permissions"))
				return
			}

			// Additional authorization checks for trading endpoints
			if sm.isTradingEndpoint(r.URL.Path) {
				if err := sm.checkTradingAuthorization(securityContext, r); err != nil {
					sm.handleAuthError(w, r, "Trading authorization failed", err)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RateLimitingMiddleware provides rate limiting
func (sm *SecurityMiddleware) RateLimitingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		clientIP := sm.getClientIP(r)
		
		// Determine rate limit type based on endpoint
		limitType := "api"
		if strings.HasPrefix(r.URL.Path, "/api/v1/auth") {
			limitType = "auth"
		}
		
		// Check rate limit
		if !sm.authManager.rateLimiter.Allow(clientIP, limitType) {
			sm.logger.Warn(ctx, "Rate limit exceeded", map[string]interface{}{
				"ip_address": clientIP,
				"endpoint":   r.URL.Path,
				"limit_type": limitType,
			})
			
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error": "rate_limit_exceeded", "message": "Too many requests"}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// SecurityHeadersMiddleware adds security headers
func (sm *SecurityMiddleware) SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		
		// Remove server information
		w.Header().Set("Server", "")
		
		next.ServeHTTP(w, r)
	})
}

// AuditLoggingMiddleware logs all requests for audit purposes
func (sm *SecurityMiddleware) AuditLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !sm.config.EnableAuditLogging {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		startTime := time.Now()
		
		// Create response writer wrapper to capture status code
		wrappedWriter := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		// Process request
		next.ServeHTTP(wrappedWriter, r)
		
		// Log audit information
		duration := time.Since(startTime)
		
		auditData := map[string]interface{}{
			"timestamp":    startTime,
			"method":       r.Method,
			"path":         r.URL.Path,
			"query":        r.URL.RawQuery,
			"status_code":  wrappedWriter.statusCode,
			"duration_ms":  duration.Milliseconds(),
			"ip_address":   sm.getClientIP(r),
			"user_agent":   r.UserAgent(),
			"content_length": r.ContentLength,
		}
		
		// Add user context if available
		if userID, ok := ctx.Value("user_id").(uuid.UUID); ok {
			auditData["user_id"] = userID.String()
		}
		
		if securityContext, ok := ctx.Value("security_context").(*SecurityContext); ok {
			auditData["auth_type"] = securityContext.AuthType
			if securityContext.Session != nil {
				auditData["session_id"] = securityContext.Session.ID
				auditData["risk_score"] = securityContext.Session.RiskScore
			}
		}
		
		sm.logger.Info(ctx, "Audit log", auditData)
	})
}

// ThreatDetectionMiddleware detects and blocks threats
func (sm *SecurityMiddleware) ThreatDetectionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !sm.config.EnableThreatDetection {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		clientIP := sm.getClientIP(r)
		
		// Check if IP is blocked
		if sm.authManager.isIPBlocked(clientIP) {
			sm.logger.Warn(ctx, "Blocked IP attempted access", map[string]interface{}{
				"ip_address": clientIP,
				"endpoint":   r.URL.Path,
			})
			
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"error": "access_denied", "message": "Access denied"}`))
			return
		}
		
		// Analyze request for threats
		threatLevel := sm.analyzeThreatLevel(r)
		if threatLevel > 80 {
			sm.logger.Warn(ctx, "High threat level detected", map[string]interface{}{
				"ip_address":   clientIP,
				"threat_level": threatLevel,
				"endpoint":     r.URL.Path,
				"user_agent":   r.UserAgent(),
			})
			
			// Block high-threat requests
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"error": "threat_detected", "message": "Suspicious activity detected"}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Helper methods

// SecurityContext holds security information for a request
type SecurityContext struct {
	UserID      uuid.UUID        `json:"user_id"`
	SessionID   string           `json:"session_id,omitempty"`
	AuthType    string           `json:"auth_type"`
	Permissions []string         `json:"permissions"`
	Session     *SecuritySession `json:"session,omitempty"`
	User        *User            `json:"user,omitempty"`
	APIKey      *APIKey          `json:"api_key,omitempty"`
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// isPublicEndpoint checks if an endpoint is public (doesn't require authentication)
func (sm *SecurityMiddleware) isPublicEndpoint(path string) bool {
	publicEndpoints := []string{
		"/health",
		"/metrics",
		"/api/v1/auth/login",
		"/api/v1/auth/register",
		"/api/v1/auth/forgot-password",
		"/api/v1/public/",
	}
	
	for _, endpoint := range publicEndpoints {
		if strings.HasPrefix(path, endpoint) {
			return true
		}
	}
	
	return false
}

// extractAuthCredentials extracts authentication credentials from request
func (sm *SecurityMiddleware) extractAuthCredentials(r *http.Request) (string, string) {
	// Check Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 {
			return parts[0], parts[1]
		}
	}
	
	// Check X-API-Key header
	apiKey := r.Header.Get("X-API-Key")
	if apiKey != "" {
		return "ApiKey", apiKey
	}
	
	return "", ""
}

// authenticateJWT authenticates a JWT token
func (sm *SecurityMiddleware) authenticateJWT(ctx context.Context, token string) (*User, *SecuritySession, error) {
	// Validate JWT token
	claims, err := sm.jwtService.ValidateToken(token)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid JWT token: %w", err)
	}
	
	// Get session information
	sessionID := claims.SessionID
	session, exists := sm.authManager.activeSessions[sessionID]
	if !exists {
		return nil, nil, fmt.Errorf("session not found")
	}
	
	// Check session validity
	if !session.IsActive || time.Now().After(session.ExpiresAt) {
		return nil, nil, fmt.Errorf("session expired")
	}
	
	// Update last activity
	session.LastActivity = time.Now()
	
	// Create user from session
	user := &User{
		ID:          session.UserID,
		Permissions: session.Permissions,
	}
	
	return user, session, nil
}

// authenticateAPIKey authenticates an API key
func (sm *SecurityMiddleware) authenticateAPIKey(ctx context.Context, keyString, ipAddress string) (*APIKey, error) {
	return sm.authManager.AuthenticateAPIKey(ctx, keyString, ipAddress)
}

// getClientIP extracts the client IP address from the request
func (sm *SecurityMiddleware) getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}
	
	// Check X-Real-IP header
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}
	
	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}
	
	return ip
}

// hasRequiredPermissions checks if user has required permissions
func (sm *SecurityMiddleware) hasRequiredPermissions(userPermissions, requiredPermissions []string) bool {
	if len(requiredPermissions) == 0 {
		return true
	}
	
	userPermsMap := make(map[string]bool)
	for _, perm := range userPermissions {
		userPermsMap[perm] = true
	}
	
	for _, required := range requiredPermissions {
		if !userPermsMap[required] {
			return false
		}
	}
	
	return true
}

// isTradingEndpoint checks if an endpoint is trading-related
func (sm *SecurityMiddleware) isTradingEndpoint(path string) bool {
	tradingPaths := []string{
		"/api/v1/trading/",
		"/api/v1/bots/",
		"/api/v1/orders/",
		"/api/v1/positions/",
	}
	
	for _, tradingPath := range tradingPaths {
		if strings.HasPrefix(path, tradingPath) {
			return true
		}
	}
	
	return false
}

// checkTradingAuthorization performs additional checks for trading endpoints
func (sm *SecurityMiddleware) checkTradingAuthorization(securityContext *SecurityContext, r *http.Request) error {
	// Check if trading is enabled for this session/API key
	if securityContext.Session != nil && !securityContext.Session.TradingEnabled {
		return fmt.Errorf("trading disabled for this session")
	}
	
	if securityContext.APIKey != nil && !securityContext.APIKey.TradingEnabled {
		return fmt.Errorf("trading disabled for this API key")
	}
	
	// Additional checks for high-risk sessions
	if securityContext.Session != nil && securityContext.Session.RiskScore > 70 {
		return fmt.Errorf("trading disabled due to high risk score")
	}
	
	return nil
}

// handleAuthError handles authentication/authorization errors
func (sm *SecurityMiddleware) handleAuthError(w http.ResponseWriter, r *http.Request, message string, err error) {
	sm.logger.Warn(r.Context(), message, map[string]interface{}{
		"error":      err.Error(),
		"endpoint":   r.URL.Path,
		"method":     r.Method,
		"ip_address": sm.getClientIP(r),
		"user_agent": r.UserAgent(),
	})
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(fmt.Sprintf(`{"error": "authentication_failed", "message": "%s"}`, message)))
}

// analyzeThreatLevel analyzes the threat level of a request
func (sm *SecurityMiddleware) analyzeThreatLevel(r *http.Request) int {
	threatLevel := 0
	
	// Check user agent
	userAgent := r.UserAgent()
	if userAgent == "" {
		threatLevel += 20
	}
	
	// Check for suspicious patterns in user agent
	suspiciousPatterns := []string{"bot", "crawler", "scanner", "curl", "wget"}
	userAgentLower := strings.ToLower(userAgent)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(userAgentLower, pattern) {
			threatLevel += 30
			break
		}
	}
	
	// Check request frequency (simplified)
	// In a real implementation, this would track request patterns
	
	return threatLevel
}
