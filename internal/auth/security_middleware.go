package auth

import (
	"crypto/subtle"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

// SecurityMiddleware provides comprehensive security middleware
type SecurityMiddleware struct {
	jwtService     *JWTService
	rbacService    *RBACService
	rateLimiter    *RateLimiter
	ipWhitelist    []string
	ipBlacklist    []string
	trustedProxies []string
}

// SecurityMiddlewareConfig contains security middleware configuration
type SecurityMiddlewareConfig struct {
	JWTService     *JWTService
	RBACService    *RBACService
	RateLimiter    *RateLimiter
	IPWhitelist    []string
	IPBlacklist    []string
	TrustedProxies []string
	CSRFProtection bool
	HSTSEnabled    bool
	HSTSMaxAge     int
}

// NewSecurityMiddleware creates a new security middleware
func NewSecurityMiddleware(config SecurityMiddlewareConfig) *SecurityMiddleware {
	return &SecurityMiddleware{
		jwtService:     config.JWTService,
		rbacService:    config.RBACService,
		rateLimiter:    config.RateLimiter,
		ipWhitelist:    config.IPWhitelist,
		ipBlacklist:    config.IPBlacklist,
		trustedProxies: config.TrustedProxies,
	}
}

// AuthenticationMiddleware validates JWT tokens
func (s *SecurityMiddleware) AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		// Check Bearer token format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		claims, err := s.jwtService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token", "details": err.Error()})
			c.Abort()
			return
		}

		// Validate IP address and User Agent for additional security
		clientIP := s.getClientIP(c)
		userAgent := c.GetHeader("User-Agent")

		if claims.IPAddress != "" && claims.IPAddress != clientIP {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "IP address mismatch"})
			c.Abort()
			return
		}

		if claims.UserAgent != "" && claims.UserAgent != userAgent {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user agent mismatch"})
			c.Abort()
			return
		}

		// Store claims in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("user_permissions", claims.Permissions)
		c.Set("team_id", claims.TeamID)
		c.Set("session_id", claims.SessionID)
		c.Set("mfa_verified", claims.MFAVerified)
		c.Set("token_claims", claims)

		c.Next()
	}
}

// AuthorizationMiddleware checks user permissions
func (s *SecurityMiddleware) AuthorizationMiddleware(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			c.Abort()
			return
		}

		teamID, _ := c.Get("team_id")
		var teamUUID *uuid.UUID
		if teamID != nil {
			if id, ok := teamID.(uuid.UUID); ok {
				teamUUID = &id
			}
		}

		// Create access request
		accessReq := AccessRequest{
			UserID:    userID.(uuid.UUID),
			Resource:  resource,
			Action:    action,
			TeamID:    teamUUID,
			IPAddress: s.getClientIP(c),
			UserAgent: c.GetHeader("User-Agent"),
			Timestamp: time.Now(),
			Context: map[string]interface{}{
				"method": c.Request.Method,
				"path":   c.Request.URL.Path,
			},
		}

		// Check access
		decision, err := s.rbacService.CheckAccess(c.Request.Context(), accessReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "authorization check failed"})
			c.Abort()
			return
		}

		if !decision.Allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"error":  "access denied",
				"reason": decision.Reason,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitingMiddleware applies rate limiting
func (s *SecurityMiddleware) RateLimitingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := s.getClientIP(c)
		userID := s.getUserID(c)

		// Check rate limits
		if !s.rateLimiter.Allow(clientIP, userID) {
			c.Header("X-RateLimit-Limit", strconv.Itoa(s.rateLimiter.GetLimit()))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Hour).Unix(), 10))

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "rate limit exceeded",
				"retry_after": 3600,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// IPFilteringMiddleware filters requests based on IP whitelist/blacklist
func (s *SecurityMiddleware) IPFilteringMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := s.getClientIP(c)

		// Check blacklist first
		if s.isIPBlacklisted(clientIP) {
			c.JSON(http.StatusForbidden, gin.H{"error": "IP address is blacklisted"})
			c.Abort()
			return
		}

		// Check whitelist if configured
		if len(s.ipWhitelist) > 0 && !s.isIPWhitelisted(clientIP) {
			c.JSON(http.StatusForbidden, gin.H{"error": "IP address not whitelisted"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SecurityHeadersMiddleware adds security headers
func (s *SecurityMiddleware) SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// HSTS (HTTP Strict Transport Security)
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		// Content Security Policy
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' https:; connect-src 'self' wss: https:")

		// X-Frame-Options
		c.Header("X-Frame-Options", "DENY")

		// X-Content-Type-Options
		c.Header("X-Content-Type-Options", "nosniff")

		// X-XSS-Protection
		c.Header("X-XSS-Protection", "1; mode=block")

		// Referrer Policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions Policy
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// Remove server information
		c.Header("Server", "")

		c.Next()
	}
}

// CSRFProtectionMiddleware provides CSRF protection
func (s *SecurityMiddleware) CSRFProtectionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip CSRF for safe methods
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// Check CSRF token
		token := c.GetHeader("X-CSRF-Token")
		if token == "" {
			token = c.PostForm("csrf_token")
		}

		if token == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "CSRF token missing"})
			c.Abort()
			return
		}

		// Validate CSRF token
		if !s.validateCSRFToken(token, c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid CSRF token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// MFARequiredMiddleware ensures MFA is verified for sensitive operations
func (s *SecurityMiddleware) MFARequiredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		mfaVerified, exists := c.Get("mfa_verified")
		if !exists || !mfaVerified.(bool) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "MFA verification required",
				"code":  "MFA_REQUIRED",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SessionValidationMiddleware validates session integrity
func (s *SecurityMiddleware) SessionValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, exists := c.Get("session_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
			c.Abort()
			return
		}

		// Validate session in database/cache
		if !s.isSessionValid(sessionID.(string)) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "session expired or invalid"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Helper methods

func (s *SecurityMiddleware) getClientIP(c *gin.Context) string {
	// Check X-Forwarded-For header
	xff := c.GetHeader("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		for _, ip := range ips {
			ip = strings.TrimSpace(ip)
			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}

	// Check X-Real-IP header
	xri := c.GetHeader("X-Real-IP")
	if xri != "" && net.ParseIP(xri) != nil {
		return xri
	}

	// Fall back to remote address
	ip, _, _ := net.SplitHostPort(c.Request.RemoteAddr)
	return ip
}

func (s *SecurityMiddleware) getUserID(c *gin.Context) string {
	if userID, exists := c.Get("user_id"); exists {
		return userID.(uuid.UUID).String()
	}
	return ""
}

func (s *SecurityMiddleware) isIPBlacklisted(ip string) bool {
	for _, blacklistedIP := range s.ipBlacklist {
		if ip == blacklistedIP {
			return true
		}
	}
	return false
}

func (s *SecurityMiddleware) isIPWhitelisted(ip string) bool {
	for _, whitelistedIP := range s.ipWhitelist {
		if ip == whitelistedIP {
			return true
		}
	}
	return false
}

func (s *SecurityMiddleware) validateCSRFToken(token string, c *gin.Context) bool {
	// In a real implementation, this would validate against a stored token
	// For now, implement a simple validation
	sessionID, exists := c.Get("session_id")
	if !exists {
		return false
	}

	expectedToken := s.generateCSRFToken(sessionID.(string))
	return subtle.ConstantTimeCompare([]byte(token), []byte(expectedToken)) == 1
}

func (s *SecurityMiddleware) generateCSRFToken(sessionID string) string {
	// In a real implementation, use a proper CSRF token generation
	return fmt.Sprintf("csrf_%s", sessionID)
}

func (s *SecurityMiddleware) isSessionValid(sessionID string) bool {
	// In a real implementation, check session in database/cache
	return true
}

// RateLimiter provides advanced rate limiting
type RateLimiter struct {
	ipLimiters    map[string]*rate.Limiter
	userLimiters  map[string]*rate.Limiter
	globalLimiter *rate.Limiter
	limit         int
	burst         int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerSecond, burst int) *RateLimiter {
	return &RateLimiter{
		ipLimiters:    make(map[string]*rate.Limiter),
		userLimiters:  make(map[string]*rate.Limiter),
		globalLimiter: rate.NewLimiter(rate.Limit(requestsPerSecond), burst),
		limit:         requestsPerSecond,
		burst:         burst,
	}
}

// Allow checks if a request should be allowed
func (r *RateLimiter) Allow(ip, userID string) bool {
	// Check global rate limit
	if !r.globalLimiter.Allow() {
		return false
	}

	// Check IP-based rate limit
	ipLimiter := r.getIPLimiter(ip)
	if !ipLimiter.Allow() {
		return false
	}

	// Check user-based rate limit if user is authenticated
	if userID != "" {
		userLimiter := r.getUserLimiter(userID)
		if !userLimiter.Allow() {
			return false
		}
	}

	return true
}

// GetLimit returns the rate limit
func (r *RateLimiter) GetLimit() int {
	return r.limit
}

func (r *RateLimiter) getIPLimiter(ip string) *rate.Limiter {
	if limiter, exists := r.ipLimiters[ip]; exists {
		return limiter
	}

	limiter := rate.NewLimiter(rate.Limit(r.limit), r.burst)
	r.ipLimiters[ip] = limiter
	return limiter
}

func (r *RateLimiter) getUserLimiter(userID string) *rate.Limiter {
	if limiter, exists := r.userLimiters[userID]; exists {
		return limiter
	}

	// Users get higher limits than anonymous IPs
	limiter := rate.NewLimiter(rate.Limit(r.limit*2), r.burst*2)
	r.userLimiters[userID] = limiter
	return limiter
}

// CleanupLimiters removes old limiters to prevent memory leaks
func (r *RateLimiter) CleanupLimiters() {
	// In a real implementation, this would run periodically
	// and remove limiters that haven't been used recently
}
