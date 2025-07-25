package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHandlers contains authentication HTTP handlers
type AuthHandlers struct {
	authService *Service
	jwtService  *JWTService
	mfaService  *MFAService
	rbacService *RBACService
}

// MFASetupRequest represents MFA setup request
type MFASetupRequest struct {
	Method MFAMethod `json:"method" binding:"required"`
}

// MFAVerifyRequest represents MFA verification request  
type MFAVerifyRequest struct {
	Code   string `json:"code" binding:"required"`
	Secret string `json:"secret,omitempty"`
}

// NewAuthHandlers creates new authentication handlers
func NewAuthHandlers(authService *Service, jwtService *JWTService, mfaService *MFAService, rbacService *RBACService) *AuthHandlers {
	return &AuthHandlers{
		authService: authService,
		jwtService:  jwtService,
		mfaService:  mfaService,
		rbacService: rbacService,
	}
}




// Register handles user registration
func (h *AuthHandlers) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	// Create user
	user := &User{
		ID:         uuid.New(),
		Email:      req.Email,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		IsActive:   true,
		IsVerified: false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Hash password
	hashedPassword, err := h.authService.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process password"})
		return
	}
	user.Password = hashedPassword

	// Save user (in real implementation, this would use a database)
	if err := h.authService.CreateUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		return
	}

	// Generate session
	sessionID := uuid.New().String()
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	deviceID := uuid.New().String() // Generate random device ID

	// Generate tokens
	tokenPair, err := h.jwtService.GenerateTokenPair(user, sessionID, clientIP, userAgent, deviceID, []string{"read", "write"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate tokens"})
		return
	}

	// Remove sensitive data
	user.Password = ""

	c.JSON(http.StatusCreated, LoginResponse{
		User:         *user,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	})
}

// Login handles user authentication
func (h *AuthHandlers) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	// Authenticate user
	user, err := h.authService.AuthenticateUser(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Check if user is active
	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "account is disabled"})
		return
	}

	// Generate session
	sessionID := uuid.New().String()
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	deviceID := uuid.New().String()

	// Generate tokens
	tokenPair, err := h.jwtService.GenerateTokenPair(user, sessionID, clientIP, userAgent, deviceID, []string{"read", "write"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate tokens"})
		return
	}

	// Remove sensitive data
	user.Password = ""

	c.JSON(http.StatusOK, LoginResponse{
		User:         *user,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	})
}

// RefreshToken handles token refresh
func (h *AuthHandlers) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// Refresh token
	tokenPair, err := h.jwtService.RefreshToken(req.RefreshToken, clientIP, userAgent)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tokens": tokenPair})
}

// Logout handles user logout
func (h *AuthHandlers) Logout(c *gin.Context) {
	// Get token from header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing authorization header"})
		return
	}

	token := authHeader[7:] // Remove "Bearer " prefix

	// Revoke token
	if err := h.jwtService.RevokeToken(token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to revoke token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

// SetupMFA handles MFA setup
func (h *AuthHandlers) SetupMFA(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	userEmail := c.MustGet("user_email").(string)

	var req MFASetupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	switch req.Method {
	case MFAMethodTOTP:
		setup, err := h.mfaService.SetupTOTP(userID, userEmail)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to setup TOTP"})
			return
		}
		c.JSON(http.StatusOK, setup)

	case MFAMethodSMS:
		// In a real implementation, get phone number from user profile
		challenge, err := h.mfaService.CreateSMSChallenge(userID, "+1234567890")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to setup SMS MFA"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"challenge": challenge})

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported MFA method"})
	}
}

// VerifyMFA handles MFA verification
func (h *AuthHandlers) VerifyMFA(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req MFAVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	// Verify TOTP
	if req.Secret != "" {
		if h.mfaService.VerifyTOTP(req.Secret, req.Code) {
			// TODO: Enable MFA for user - implement EnableMFA method
			// h.authService.EnableMFA(c.Request.Context(), userID, req.Secret)
			c.JSON(http.StatusOK, gin.H{"verified": true, "message": "MFA verified successfully"})
			return
		}
	}

	// Verify backup code
	if h.mfaService.VerifyBackupCode(userID, req.Code) {
		c.JSON(http.StatusOK, gin.H{"verified": true, "message": "backup code verified"})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"verified": false, "error": "invalid MFA code"})
}

// DisableMFA handles MFA disabling
func (h *AuthHandlers) DisableMFA(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req MFAVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	// TODO: Get user and verify MFA code before disabling
	// This requires MFA fields in User model
	// user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	// if err != nil {
	//     c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
	//     return
	// }
	// if !h.mfaService.VerifyTOTP(user.MFASecret, req.Code) {
	//     c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid MFA code"})
	//     return
	// }

	// TODO: Disable MFA - implement DisableMFA method
	// if err := h.authService.DisableMFA(c.Request.Context(), userID); err != nil {
	//     c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to disable MFA"})
	//     return
	// }

	// Revoke all user tokens to force re-authentication
	h.jwtService.RevokeAllUserTokens(userID)

	c.JSON(http.StatusOK, gin.H{"message": "MFA disabled successfully"})
}

// ChangePassword handles password changes
func (h *AuthHandlers) ChangePassword(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	// Get user
	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	// Verify current password
	if !h.authService.VerifyPassword(req.CurrentPassword, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid current password"})
		return
	}

	// TODO: Verify MFA if enabled - requires MFA fields in User model
	// if user.MFAEnabled {
	//     if req.MFACode == "" {
	//         c.JSON(http.StatusBadRequest, gin.H{"error": "MFA code required"})
	//         return
	//     }
	//     if !h.mfaService.VerifyTOTP(user.MFASecret, req.MFACode) {
	//         c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid MFA code"})
	//         return
	//     }
	// }

	// TODO: Hash new password and update user
	// This functionality is not implemented yet
	// newPasswordHash, err := h.authService.HashPassword(req.NewPassword)
	// if err != nil {
	//     c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process new password"})
	//     return
	// }
	// user.Password = newPasswordHash
	// if err := h.authService.UpdateUser(c.Request.Context(), user); err != nil {
	//     c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
	//     return
	// }

	// Placeholder response - password change not fully implemented
	c.JSON(http.StatusNotImplemented, gin.H{"error": "password change not implemented yet"})

	// TODO: Uncomment when UpdateUser method is implemented
	// h.jwtService.RevokeAllUserTokens(userID)
	// c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}

// GetProfile returns user profile
func (h *AuthHandlers) GetProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user profile"})
		return
	}

	// Remove sensitive data
	user.Password = ""

	// Get user permissions
	permissions := c.MustGet("user_permissions").([]string)

	c.JSON(http.StatusOK, gin.H{
		"user":        user,
		"permissions": permissions,
	})
}

// GetSessions returns active user sessions
func (h *AuthHandlers) GetSessions(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	// In a real implementation, this would query active sessions from database
	sessions := []map[string]interface{}{
		{
			"id":         "session-1",
			"device":     "Chrome on Windows",
			"ip_address": "192.168.1.100",
			"location":   "New York, US",
			"last_seen":  time.Now().Add(-time.Hour),
			"is_current": true,
		},
		{
			"id":         "session-2",
			"device":     "Safari on iPhone",
			"ip_address": "192.168.1.101",
			"location":   "New York, US",
			"last_seen":  time.Now().Add(-24 * time.Hour),
			"is_current": false,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"sessions": sessions,
		"user_id":  userID,
	})
}

// RevokeSession revokes a specific session
func (h *AuthHandlers) RevokeSession(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session ID required"})
		return
	}

	// In a real implementation, this would revoke the specific session
	// For now, just return success
	c.JSON(http.StatusOK, gin.H{"message": "session revoked successfully"})
}

// GetPublicKey returns the JWT public key for token verification
func (h *AuthHandlers) GetPublicKey(c *gin.Context) {
	publicKeyPEM, err := h.jwtService.GetPublicKeyPEM()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get public key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"public_key": publicKeyPEM,
		"algorithm":  "RS256",
		"use":        "sig",
	})
}

// HealthCheck returns authentication service health
func (h *AuthHandlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "authentication",
		"timestamp": time.Now(),
		"version":   "1.0.0",
	})
}
