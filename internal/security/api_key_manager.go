package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/argon2"
)

// APIKeyManager manages API keys and their lifecycle
type APIKeyManager struct {
	logger    *observability.Logger
	apiKeys   map[string]*APIKey
	userKeys  map[uuid.UUID][]string // userID -> []keyID
	config    *SecurityConfig
	mu        sync.RWMutex
}

// APIKeyRequest represents a request to create an API key
type APIKeyRequest struct {
	UserID          uuid.UUID       `json:"user_id"`
	Name            string          `json:"name"`
	Permissions     []string        `json:"permissions"`
	IPWhitelist     []string        `json:"ip_whitelist,omitempty"`
	ExpiresAt       *time.Time      `json:"expires_at,omitempty"`
	RateLimit       int             `json:"rate_limit,omitempty"`
	TradingEnabled  bool            `json:"trading_enabled"`
	MaxTradeAmount  decimal.Decimal `json:"max_trade_amount,omitempty"`
	AllowedPairs    []string        `json:"allowed_pairs,omitempty"`
	SecurityLevel   SecurityLevel   `json:"security_level"`
}

// APIKeyResponse represents the response when creating an API key
type APIKeyResponse struct {
	KeyID     string    `json:"key_id"`
	KeySecret string    `json:"key_secret"`
	FullKey   string    `json:"full_key"`
	APIKey    *APIKey   `json:"api_key"`
}

// NewAPIKeyManager creates a new API key manager
func NewAPIKeyManager(logger *observability.Logger, config *SecurityConfig) *APIKeyManager {
	return &APIKeyManager{
		logger:   logger,
		apiKeys:  make(map[string]*APIKey),
		userKeys: make(map[uuid.UUID][]string),
		config:   config,
	}
}

// CreateAPIKey creates a new API key
func (akm *APIKeyManager) CreateAPIKey(request *APIKeyRequest) (*APIKeyResponse, error) {
	akm.mu.Lock()
	defer akm.mu.Unlock()

	// Validate request
	if err := akm.validateAPIKeyRequest(request); err != nil {
		return nil, fmt.Errorf("invalid API key request: %w", err)
	}

	// Check user's existing key count
	userKeyCount := len(akm.userKeys[request.UserID])
	if userKeyCount >= 10 { // Maximum 10 keys per user
		return nil, fmt.Errorf("maximum API key limit reached")
	}

	// Generate key ID and secret
	keyID := akm.generateKeyID()
	keySecret := akm.generateKeySecret()
	keyHash := akm.hashKeySecret(keySecret)

	// Set default values
	if request.RateLimit == 0 {
		request.RateLimit = akm.config.APIRateLimit
	}

	if request.MaxTradeAmount.IsZero() {
		request.MaxTradeAmount = decimal.NewFromFloat(1000) // Default $1000 limit
	}

	// Create API key
	apiKey := &APIKey{
		ID:             keyID,
		UserID:         request.UserID,
		Name:           request.Name,
		KeyHash:        keyHash,
		Permissions:    request.Permissions,
		IPWhitelist:    request.IPWhitelist,
		CreatedAt:      time.Now(),
		LastUsed:       time.Time{},
		ExpiresAt:      request.ExpiresAt,
		IsActive:       true,
		RateLimit:      request.RateLimit,
		TradingEnabled: request.TradingEnabled,
		MaxTradeAmount: request.MaxTradeAmount,
		AllowedPairs:   request.AllowedPairs,
		SecurityLevel:  request.SecurityLevel,
		Metadata:       make(map[string]interface{}),
	}

	// Store API key
	akm.apiKeys[keyID] = apiKey
	akm.userKeys[request.UserID] = append(akm.userKeys[request.UserID], keyID)

	// Create full key string
	fullKey := fmt.Sprintf("%s.%s", keyID, keySecret)

	akm.logger.Info(nil, "API key created", map[string]interface{}{
		"key_id":         keyID,
		"user_id":        request.UserID.String(),
		"name":           request.Name,
		"security_level": string(request.SecurityLevel),
		"trading_enabled": request.TradingEnabled,
	})

	return &APIKeyResponse{
		KeyID:     keyID,
		KeySecret: keySecret,
		FullKey:   fullKey,
		APIKey:    apiKey,
	}, nil
}

// GetAPIKey retrieves an API key by ID
func (akm *APIKeyManager) GetAPIKey(keyID string) (*APIKey, error) {
	akm.mu.RLock()
	defer akm.mu.RUnlock()

	apiKey, exists := akm.apiKeys[keyID]
	if !exists {
		return nil, fmt.Errorf("API key not found")
	}

	return apiKey, nil
}

// GetUserAPIKeys retrieves all API keys for a user
func (akm *APIKeyManager) GetUserAPIKeys(userID uuid.UUID) ([]*APIKey, error) {
	akm.mu.RLock()
	defer akm.mu.RUnlock()

	keyIDs, exists := akm.userKeys[userID]
	if !exists {
		return []*APIKey{}, nil
	}

	apiKeys := make([]*APIKey, 0, len(keyIDs))
	for _, keyID := range keyIDs {
		if apiKey, exists := akm.apiKeys[keyID]; exists {
			apiKeys = append(apiKeys, apiKey)
		}
	}

	return apiKeys, nil
}

// RevokeAPIKey revokes an API key
func (akm *APIKeyManager) RevokeAPIKey(keyID string, userID uuid.UUID) error {
	akm.mu.Lock()
	defer akm.mu.Unlock()

	apiKey, exists := akm.apiKeys[keyID]
	if !exists {
		return fmt.Errorf("API key not found")
	}

	// Verify ownership
	if apiKey.UserID != userID {
		return fmt.Errorf("unauthorized to revoke this API key")
	}

	// Deactivate the key
	apiKey.IsActive = false

	akm.logger.Info(nil, "API key revoked", map[string]interface{}{
		"key_id":  keyID,
		"user_id": userID.String(),
		"name":    apiKey.Name,
	})

	return nil
}

// UpdateAPIKey updates an API key's properties
func (akm *APIKeyManager) UpdateAPIKey(keyID string, userID uuid.UUID, updates *APIKeyUpdateRequest) error {
	akm.mu.Lock()
	defer akm.mu.Unlock()

	apiKey, exists := akm.apiKeys[keyID]
	if !exists {
		return fmt.Errorf("API key not found")
	}

	// Verify ownership
	if apiKey.UserID != userID {
		return fmt.Errorf("unauthorized to update this API key")
	}

	// Apply updates
	if updates.Name != nil {
		apiKey.Name = *updates.Name
	}

	if updates.IPWhitelist != nil {
		apiKey.IPWhitelist = updates.IPWhitelist
	}

	if updates.RateLimit != nil {
		apiKey.RateLimit = *updates.RateLimit
	}

	if updates.TradingEnabled != nil {
		apiKey.TradingEnabled = *updates.TradingEnabled
	}

	if updates.MaxTradeAmount != nil {
		apiKey.MaxTradeAmount = *updates.MaxTradeAmount
	}

	if updates.AllowedPairs != nil {
		apiKey.AllowedPairs = updates.AllowedPairs
	}

	if updates.IsActive != nil {
		apiKey.IsActive = *updates.IsActive
	}

	akm.logger.Info(nil, "API key updated", map[string]interface{}{
		"key_id":  keyID,
		"user_id": userID.String(),
		"name":    apiKey.Name,
	})

	return nil
}

// RotateAPIKey rotates an API key (generates new secret)
func (akm *APIKeyManager) RotateAPIKey(keyID string, userID uuid.UUID) (*APIKeyResponse, error) {
	akm.mu.Lock()
	defer akm.mu.Unlock()

	apiKey, exists := akm.apiKeys[keyID]
	if !exists {
		return nil, fmt.Errorf("API key not found")
	}

	// Verify ownership
	if apiKey.UserID != userID {
		return nil, fmt.Errorf("unauthorized to rotate this API key")
	}

	// Generate new secret
	newKeySecret := akm.generateKeySecret()
	newKeyHash := akm.hashKeySecret(newKeySecret)

	// Update the key
	apiKey.KeyHash = newKeyHash

	// Create full key string
	fullKey := fmt.Sprintf("%s.%s", keyID, newKeySecret)

	akm.logger.Info(nil, "API key rotated", map[string]interface{}{
		"key_id":  keyID,
		"user_id": userID.String(),
		"name":    apiKey.Name,
	})

	return &APIKeyResponse{
		KeyID:     keyID,
		KeySecret: newKeySecret,
		FullKey:   fullKey,
		APIKey:    apiKey,
	}, nil
}

// ValidateAPIKey validates an API key and returns the key if valid
func (akm *APIKeyManager) ValidateAPIKey(fullKey, ipAddress string) (*APIKey, error) {
	// Parse the key
	parts := strings.Split(fullKey, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid API key format")
	}

	keyID := parts[0]
	keySecret := parts[1]

	akm.mu.RLock()
	defer akm.mu.RUnlock()

	// Get the API key
	apiKey, exists := akm.apiKeys[keyID]
	if !exists {
		return nil, fmt.Errorf("API key not found")
	}

	// Check if key is active
	if !apiKey.IsActive {
		return nil, fmt.Errorf("API key is disabled")
	}

	// Check expiration
	if apiKey.ExpiresAt != nil && time.Now().After(*apiKey.ExpiresAt) {
		return nil, fmt.Errorf("API key has expired")
	}

	// Verify the secret
	if !akm.verifyKeySecret(apiKey.KeyHash, keySecret) {
		return nil, fmt.Errorf("invalid API key secret")
	}

	// Check IP whitelist
	if len(apiKey.IPWhitelist) > 0 && !akm.isIPWhitelisted(ipAddress, apiKey.IPWhitelist) {
		return nil, fmt.Errorf("IP address not whitelisted")
	}

	// Update last used timestamp
	apiKey.LastUsed = time.Now()

	return apiKey, nil
}

// Helper methods

// validateAPIKeyRequest validates an API key creation request
func (akm *APIKeyManager) validateAPIKeyRequest(request *APIKeyRequest) error {
	if request.UserID == uuid.Nil {
		return fmt.Errorf("user ID is required")
	}

	if request.Name == "" {
		return fmt.Errorf("API key name is required")
	}

	if len(request.Name) > 100 {
		return fmt.Errorf("API key name too long")
	}

	if len(request.Permissions) == 0 {
		return fmt.Errorf("at least one permission is required")
	}

	// Validate permissions
	validPermissions := map[string]bool{
		"trading:read":     true,
		"trading:write":    true,
		"account:read":     true,
		"account:write":    true,
		"monitoring:read":  true,
		"monitoring:write": true,
		"admin:read":       true,
		"admin:write":      true,
	}

	for _, perm := range request.Permissions {
		if !validPermissions[perm] {
			return fmt.Errorf("invalid permission: %s", perm)
		}
	}

	// Validate security level
	switch request.SecurityLevel {
	case SecurityLevelReadOnly, SecurityLevelTrading, SecurityLevelAdmin:
		// Valid
	default:
		return fmt.Errorf("invalid security level")
	}

	return nil
}

// generateKeyID generates a unique key ID
func (akm *APIKeyManager) generateKeyID() string {
	// Generate a random key ID
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)[:22] // Remove padding
}

// generateKeySecret generates a secure key secret
func (akm *APIKeyManager) generateKeySecret() string {
	// Generate a random secret
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)[:43] // Remove padding
}

// hashKeySecret hashes a key secret using Argon2
func (akm *APIKeyManager) hashKeySecret(secret string) string {
	// Generate salt
	salt := make([]byte, 16)
	rand.Read(salt)

	// Hash the secret
	hash := argon2.IDKey([]byte(secret), salt, 1, 64*1024, 4, 32)

	// Combine salt and hash
	combined := append(salt, hash...)

	// Encode to base64
	return base64.StdEncoding.EncodeToString(combined)
}

// verifyKeySecret verifies a key secret against its hash
func (akm *APIKeyManager) verifyKeySecret(keyHash, secret string) bool {
	// Decode the stored hash
	storedHash, err := base64.StdEncoding.DecodeString(keyHash)
	if err != nil {
		return false
	}

	// Extract salt and hash
	if len(storedHash) < 32 {
		return false
	}

	salt := storedHash[:16]
	hash := storedHash[16:]

	// Hash the provided secret with the same salt
	providedHash := argon2.IDKey([]byte(secret), salt, 1, 64*1024, 4, 32)

	// Compare hashes using constant-time comparison
	return subtle.ConstantTimeCompare(hash, providedHash) == 1
}

// isIPWhitelisted checks if an IP is in the whitelist
func (akm *APIKeyManager) isIPWhitelisted(ipAddress string, whitelist []string) bool {
	for _, whitelistedIP := range whitelist {
		if ipAddress == whitelistedIP {
			return true
		}
		// Could add CIDR range checking here
	}
	return false
}

// APIKeyUpdateRequest represents a request to update an API key
type APIKeyUpdateRequest struct {
	Name           *string          `json:"name,omitempty"`
	IPWhitelist    []string         `json:"ip_whitelist,omitempty"`
	RateLimit      *int             `json:"rate_limit,omitempty"`
	TradingEnabled *bool            `json:"trading_enabled,omitempty"`
	MaxTradeAmount *decimal.Decimal `json:"max_trade_amount,omitempty"`
	AllowedPairs   []string         `json:"allowed_pairs,omitempty"`
	IsActive       *bool            `json:"is_active,omitempty"`
}
