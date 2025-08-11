package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ai-agentic-browser/internal/security"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

// SecurityHandler handles security and authentication API requests
type SecurityHandler struct {
	logger      *observability.Logger
	authManager *security.AuthManager
	middleware  *security.SecurityMiddleware
}

// NewSecurityHandler creates a new security handler
func NewSecurityHandler(
	logger *observability.Logger,
	authManager *security.AuthManager,
	middleware *security.SecurityMiddleware,
) *SecurityHandler {
	return &SecurityHandler{
		logger:      logger,
		authManager: authManager,
		middleware:  middleware,
	}
}

// RegisterRoutes registers security API routes
func (h *SecurityHandler) RegisterRoutes(router *mux.Router) {
	// Authentication endpoints
	router.HandleFunc("/api/v1/auth/login", h.Login).Methods("POST")
	router.HandleFunc("/api/v1/auth/logout", h.Logout).Methods("POST")
	router.HandleFunc("/api/v1/auth/refresh", h.RefreshToken).Methods("POST")
	router.HandleFunc("/api/v1/auth/verify", h.VerifyToken).Methods("POST")

	// MFA endpoints
	router.HandleFunc("/api/v1/auth/mfa/setup", h.SetupMFA).Methods("POST")
	router.HandleFunc("/api/v1/auth/mfa/verify", h.VerifyMFA).Methods("POST")
	router.HandleFunc("/api/v1/auth/mfa/disable", h.DisableMFA).Methods("POST")

	// Device management endpoints
	router.HandleFunc("/api/v1/auth/devices", h.GetTrustedDevices).Methods("GET")
	router.HandleFunc("/api/v1/auth/devices/trust", h.TrustDevice).Methods("POST")
	router.HandleFunc("/api/v1/auth/devices/{deviceId}/revoke", h.RevokeDevice).Methods("DELETE")

	// API key management endpoints
	router.HandleFunc("/api/v1/auth/api-keys", h.GetAPIKeys).Methods("GET")
	router.HandleFunc("/api/v1/auth/api-keys", h.CreateAPIKey).Methods("POST")
	router.HandleFunc("/api/v1/auth/api-keys/{keyId}", h.GetAPIKey).Methods("GET")
	router.HandleFunc("/api/v1/auth/api-keys/{keyId}", h.UpdateAPIKey).Methods("PUT")
	router.HandleFunc("/api/v1/auth/api-keys/{keyId}/rotate", h.RotateAPIKey).Methods("POST")
	router.HandleFunc("/api/v1/auth/api-keys/{keyId}/revoke", h.RevokeAPIKey).Methods("DELETE")

	// Session management endpoints
	router.HandleFunc("/api/v1/auth/sessions", h.GetActiveSessions).Methods("GET")
	router.HandleFunc("/api/v1/auth/sessions/{sessionId}/revoke", h.RevokeSession).Methods("DELETE")
	router.HandleFunc("/api/v1/auth/sessions/revoke-all", h.RevokeAllSessions).Methods("POST")

	// Security monitoring endpoints
	router.HandleFunc("/api/v1/security/audit-logs", h.GetAuditLogs).Methods("GET")
	router.HandleFunc("/api/v1/security/suspicious-activities", h.GetSuspiciousActivities).Methods("GET")
	router.HandleFunc("/api/v1/security/blocked-ips", h.GetBlockedIPs).Methods("GET")
	router.HandleFunc("/api/v1/security/blocked-ips/{ip}/unblock", h.UnblockIP).Methods("DELETE")

	// Security settings endpoints
	router.HandleFunc("/api/v1/security/settings", h.GetSecuritySettings).Methods("GET")
	router.HandleFunc("/api/v1/security/settings", h.UpdateSecuritySettings).Methods("PUT")
}

// Login handles POST /api/v1/auth/login
func (h *SecurityHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var loginRequest struct {
		Username          string `json:"username"`
		Password          string `json:"password"`
		MFACode           string `json:"mfa_code,omitempty"`
		DeviceFingerprint string `json:"device_fingerprint"`
		RememberDevice    bool   `json:"remember_device"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create authentication request
	authRequest := &security.AuthenticationRequest{
		Username:          loginRequest.Username,
		Password:          loginRequest.Password,
		MFACode:           security.NullString{String: loginRequest.MFACode, Valid: loginRequest.MFACode != ""},
		IPAddress:         h.getClientIP(r),
		UserAgent:         r.UserAgent(),
		DeviceFingerprint: loginRequest.DeviceFingerprint,
		RememberDevice:    loginRequest.RememberDevice,
		Metadata:          make(map[string]interface{}),
	}

	// Authenticate user
	response, err := h.authManager.AuthenticateUser(ctx, authRequest)
	if err != nil {
		h.logger.Error(ctx, "Authentication failed", err, map[string]interface{}{
			"username":   loginRequest.Username,
			"ip_address": authRequest.IPAddress,
		})
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	h.logger.Info(ctx, "User authenticated successfully", map[string]interface{}{
		"username":             loginRequest.Username,
		"success":              response.Success,
		"require_mfa":          response.RequireMFA,
		"require_device_setup": response.RequireDeviceSetup,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Logout handles POST /api/v1/auth/logout
func (h *SecurityHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get security context
	securityContext, ok := ctx.Value("security_context").(*security.SecurityContext)
	if !ok {
		http.Error(w, "No active session", http.StatusBadRequest)
		return
	}

	// Revoke session if it exists
	if securityContext.SessionID != "" {
		// Implementation would revoke the session
		h.logger.Info(ctx, "User logged out", map[string]interface{}{
			"user_id":    securityContext.UserID.String(),
			"session_id": securityContext.SessionID,
		})
	}

	response := map[string]interface{}{
		"message": "Logged out successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateAPIKey handles POST /api/v1/auth/api-keys
func (h *SecurityHandler) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user ID from context
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var createRequest struct {
		Name           string                 `json:"name"`
		Permissions    []string               `json:"permissions"`
		IPWhitelist    []string               `json:"ip_whitelist,omitempty"`
		ExpiresAt      *time.Time             `json:"expires_at,omitempty"`
		RateLimit      int                    `json:"rate_limit,omitempty"`
		TradingEnabled bool                   `json:"trading_enabled"`
		MaxTradeAmount decimal.Decimal        `json:"max_trade_amount,omitempty"`
		AllowedPairs   []string               `json:"allowed_pairs,omitempty"`
		SecurityLevel  security.SecurityLevel `json:"security_level"`
	}

	if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create API key request
	apiKeyRequest := &security.APIKeyRequest{
		UserID:         userID,
		Name:           createRequest.Name,
		Permissions:    createRequest.Permissions,
		IPWhitelist:    createRequest.IPWhitelist,
		ExpiresAt:      createRequest.ExpiresAt,
		RateLimit:      createRequest.RateLimit,
		TradingEnabled: createRequest.TradingEnabled,
		MaxTradeAmount: createRequest.MaxTradeAmount,
		AllowedPairs:   createRequest.AllowedPairs,
		SecurityLevel:  createRequest.SecurityLevel,
	}

	// Create API key
	response, err := h.authManager.GetAPIKeyManager().CreateAPIKey(apiKeyRequest)
	if err != nil {
		h.logger.Error(ctx, "Failed to create API key", err, map[string]interface{}{
			"user_id": userID.String(),
			"name":    createRequest.Name,
		})
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Info(ctx, "API key created", map[string]interface{}{
		"user_id": userID.String(),
		"key_id":  response.KeyID,
		"name":    createRequest.Name,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAPIKeys handles GET /api/v1/auth/api-keys
func (h *SecurityHandler) GetAPIKeys(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user ID from context
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Get user's API keys
	apiKeys, err := h.authManager.GetAPIKeyManager().GetUserAPIKeys(userID)
	if err != nil {
		h.logger.Error(ctx, "Failed to get API keys", err, map[string]interface{}{
			"user_id": userID.String(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Remove sensitive information before returning
	sanitizedKeys := make([]map[string]interface{}, len(apiKeys))
	for i, key := range apiKeys {
		sanitizedKeys[i] = map[string]interface{}{
			"id":               key.ID,
			"name":             key.Name,
			"permissions":      key.Permissions,
			"ip_whitelist":     key.IPWhitelist,
			"created_at":       key.CreatedAt,
			"last_used":        key.LastUsed,
			"expires_at":       key.ExpiresAt,
			"is_active":        key.IsActive,
			"rate_limit":       key.RateLimit,
			"trading_enabled":  key.TradingEnabled,
			"max_trade_amount": key.MaxTradeAmount,
			"allowed_pairs":    key.AllowedPairs,
			"security_level":   key.SecurityLevel,
		}
	}

	response := map[string]interface{}{
		"api_keys": sanitizedKeys,
		"count":    len(sanitizedKeys),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RevokeAPIKey handles DELETE /api/v1/auth/api-keys/{keyId}/revoke
func (h *SecurityHandler) RevokeAPIKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	keyID := vars["keyId"]

	// Get user ID from context
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Revoke API key
	err := h.authManager.GetAPIKeyManager().RevokeAPIKey(keyID, userID)
	if err != nil {
		h.logger.Error(ctx, "Failed to revoke API key", err, map[string]interface{}{
			"user_id": userID.String(),
			"key_id":  keyID,
		})
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Info(ctx, "API key revoked", map[string]interface{}{
		"user_id": userID.String(),
		"key_id":  keyID,
	})

	response := map[string]interface{}{
		"message": "API key revoked successfully",
		"key_id":  keyID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetSecuritySettings handles GET /api/v1/security/settings
func (h *SecurityHandler) GetSecuritySettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user ID from context
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Get security context for additional information
	securityContext, _ := ctx.Value("security_context").(*security.SecurityContext)

	settings := map[string]interface{}{
		"user_id":                 userID.String(),
		"mfa_enabled":             true, // Would come from user service
		"device_trust_enabled":    true,
		"geo_verification":        false,
		"trading_enabled":         true,
		"max_concurrent_sessions": 5,
		"session_timeout":         "24h",
		"api_key_count":           0, // Would be calculated
		"trusted_device_count":    0, // Would be calculated
		"last_login":              time.Now(),
		"risk_score":              0,
	}

	if securityContext != nil && securityContext.Session != nil {
		settings["current_session"] = map[string]interface{}{
			"id":             securityContext.Session.ID,
			"created_at":     securityContext.Session.CreatedAt,
			"last_activity":  securityContext.Session.LastActivity,
			"expires_at":     securityContext.Session.ExpiresAt,
			"device_trusted": securityContext.Session.DeviceTrusted,
			"risk_score":     securityContext.Session.RiskScore,
			"ip_address":     securityContext.Session.IPAddress,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

// Helper methods

// getClientIP extracts the client IP address from the request
func (h *SecurityHandler) getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		return xff
	}

	// Check X-Real-IP header
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}

// Placeholder implementations for remaining endpoints
func (h *SecurityHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *SecurityHandler) VerifyToken(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *SecurityHandler) SetupMFA(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *SecurityHandler) VerifyMFA(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *SecurityHandler) DisableMFA(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *SecurityHandler) GetTrustedDevices(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *SecurityHandler) TrustDevice(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *SecurityHandler) RevokeDevice(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *SecurityHandler) GetAPIKey(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *SecurityHandler) UpdateAPIKey(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *SecurityHandler) RotateAPIKey(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *SecurityHandler) GetActiveSessions(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *SecurityHandler) RevokeSession(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *SecurityHandler) RevokeAllSessions(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *SecurityHandler) GetAuditLogs(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *SecurityHandler) GetSuspiciousActivities(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *SecurityHandler) GetBlockedIPs(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *SecurityHandler) UnblockIP(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *SecurityHandler) UpdateSecuritySettings(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
