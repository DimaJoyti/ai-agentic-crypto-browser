package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTService handles JWT token operations with advanced security features
type JWTService struct {
	privateKey       *rsa.PrivateKey
	publicKey        *rsa.PublicKey
	issuer           string
	accessTokenTTL   time.Duration
	refreshTokenTTL  time.Duration
	blacklistService *TokenBlacklistService
}

// TokenClaims represents JWT token claims with enhanced security
type TokenClaims struct {
	UserID      uuid.UUID `json:"user_id"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	Permissions []string  `json:"permissions"`
	TeamID      *uuid.UUID `json:"team_id,omitempty"`
	SessionID   string    `json:"session_id"`
	TokenType   string    `json:"token_type"` // access, refresh, mfa
	MFAVerified bool      `json:"mfa_verified"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	DeviceID    string    `json:"device_id"`
	Scope       []string  `json:"scope"`
	jwt.RegisteredClaims
}

// TokenPair represents access and refresh token pair
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int64     `json:"expires_in"`
	ExpiresAt    time.Time `json:"expires_at"`
	Scope        []string  `json:"scope"`
}

// NewJWTService creates a new JWT service with RSA key pair
func NewJWTService(issuer string, accessTTL, refreshTTL time.Duration, blacklistService *TokenBlacklistService) (*JWTService, error) {
	// Generate RSA key pair for enhanced security
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key pair: %w", err)
	}

	return &JWTService{
		privateKey:       privateKey,
		publicKey:        &privateKey.PublicKey,
		issuer:           issuer,
		accessTokenTTL:   accessTTL,
		refreshTokenTTL:  refreshTTL,
		blacklistService: blacklistService,
	}, nil
}

// GenerateTokenPair creates a new access and refresh token pair
func (j *JWTService) GenerateTokenPair(user *User, sessionID, ipAddress, userAgent, deviceID string, scope []string) (*TokenPair, error) {
	now := time.Now()
	accessExpiresAt := now.Add(j.accessTokenTTL)
	refreshExpiresAt := now.Add(j.refreshTokenTTL)

	// Create access token claims
	accessClaims := &TokenClaims{
		UserID:      user.ID,
		Email:       user.Email,
		Role:        user.Role,
		Permissions: user.Permissions,
		TeamID:      user.TeamID,
		SessionID:   sessionID,
		TokenType:   "access",
		MFAVerified: user.MFAEnabled && user.MFAVerified,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		DeviceID:    deviceID,
		Scope:       scope,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Issuer:    j.issuer,
			Subject:   user.ID.String(),
			Audience:  []string{"agentic-browser"},
			ExpiresAt: jwt.NewNumericDate(accessExpiresAt),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	// Create refresh token claims
	refreshClaims := &TokenClaims{
		UserID:    user.ID,
		Email:     user.Email,
		SessionID: sessionID,
		TokenType: "refresh",
		IPAddress: ipAddress,
		UserAgent: userAgent,
		DeviceID:  deviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Issuer:    j.issuer,
			Subject:   user.ID.String(),
			Audience:  []string{"agentic-browser"},
			ExpiresAt: jwt.NewNumericDate(refreshExpiresAt),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	// Sign tokens
	accessToken, err := j.signToken(accessClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	refreshToken, err := j.signToken(refreshClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(j.accessTokenTTL.Seconds()),
		ExpiresAt:    accessExpiresAt,
		Scope:        scope,
	}, nil
}

// ValidateToken validates and parses a JWT token
func (j *JWTService) ValidateToken(tokenString string) (*TokenClaims, error) {
	// Check if token is blacklisted
	if j.blacklistService.IsBlacklisted(tokenString) {
		return nil, fmt.Errorf("token is blacklisted")
	}

	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Additional validation
	if err := j.validateClaims(claims); err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	return claims, nil
}

// RefreshToken creates a new access token from a valid refresh token
func (j *JWTService) RefreshToken(refreshTokenString, ipAddress, userAgent string) (*TokenPair, error) {
	claims, err := j.ValidateToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	if claims.TokenType != "refresh" {
		return nil, fmt.Errorf("token is not a refresh token")
	}

	// Validate IP and User Agent for security
	if claims.IPAddress != ipAddress {
		return nil, fmt.Errorf("IP address mismatch")
	}

	// Blacklist the old refresh token
	j.blacklistService.BlacklistToken(refreshTokenString, claims.ExpiresAt.Time)

	// Get user details (this would typically come from a user service)
	user := &User{
		ID:          claims.UserID,
		Email:       claims.Email,
		Role:        claims.Role,
		Permissions: claims.Permissions,
		TeamID:      claims.TeamID,
	}

	// Generate new token pair
	return j.GenerateTokenPair(user, claims.SessionID, ipAddress, userAgent, claims.DeviceID, claims.Scope)
}

// RevokeToken blacklists a token
func (j *JWTService) RevokeToken(tokenString string) error {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return fmt.Errorf("cannot revoke invalid token: %w", err)
	}

	j.blacklistService.BlacklistToken(tokenString, claims.ExpiresAt.Time)
	return nil
}

// RevokeAllUserTokens blacklists all tokens for a user
func (j *JWTService) RevokeAllUserTokens(userID uuid.UUID) error {
	return j.blacklistService.BlacklistAllUserTokens(userID)
}

// GenerateMFAToken creates a temporary token for MFA verification
func (j *JWTService) GenerateMFAToken(userID uuid.UUID, email, ipAddress, userAgent string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(5 * time.Minute) // MFA tokens expire in 5 minutes

	claims := &TokenClaims{
		UserID:    userID,
		Email:     email,
		TokenType: "mfa",
		IPAddress: ipAddress,
		UserAgent: userAgent,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Issuer:    j.issuer,
			Subject:   userID.String(),
			Audience:  []string{"agentic-browser-mfa"},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	return j.signToken(claims)
}

// GetPublicKey returns the RSA public key for token verification
func (j *JWTService) GetPublicKey() *rsa.PublicKey {
	return j.publicKey
}

// GetPublicKeyPEM returns the public key in PEM format
func (j *JWTService) GetPublicKeyPEM() (string, error) {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(j.publicKey)
	if err != nil {
		return "", fmt.Errorf("failed to marshal public key: %w", err)
	}

	pubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	})

	return string(pubKeyPEM), nil
}

// signToken signs a token with the private key
func (j *JWTService) signToken(claims *TokenClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	
	// Add key ID to header for key rotation support
	token.Header["kid"] = j.getKeyID()
	
	return token.SignedString(j.privateKey)
}

// validateClaims performs additional claim validation
func (j *JWTService) validateClaims(claims *TokenClaims) error {
	now := time.Now()

	// Check expiration
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(now) {
		return fmt.Errorf("token has expired")
	}

	// Check not before
	if claims.NotBefore != nil && claims.NotBefore.After(now) {
		return fmt.Errorf("token not yet valid")
	}

	// Check issuer
	if claims.Issuer != j.issuer {
		return fmt.Errorf("invalid issuer")
	}

	// Check audience
	validAudience := false
	for _, aud := range claims.Audience {
		if aud == "agentic-browser" || aud == "agentic-browser-mfa" {
			validAudience = true
			break
		}
	}
	if !validAudience {
		return fmt.Errorf("invalid audience")
	}

	return nil
}

// getKeyID generates a key ID for the current key
func (j *JWTService) getKeyID() string {
	pubKeyBytes, _ := x509.MarshalPKIXPublicKey(j.publicKey)
	return base64.URLEncoding.EncodeToString(pubKeyBytes[:8])
}

// TokenBlacklistService manages blacklisted tokens
type TokenBlacklistService struct {
	// This would typically use Redis or a database
	blacklistedTokens map[string]time.Time
}

// NewTokenBlacklistService creates a new token blacklist service
func NewTokenBlacklistService() *TokenBlacklistService {
	return &TokenBlacklistService{
		blacklistedTokens: make(map[string]time.Time),
	}
}

// BlacklistToken adds a token to the blacklist
func (t *TokenBlacklistService) BlacklistToken(token string, expiresAt time.Time) {
	t.blacklistedTokens[token] = expiresAt
}

// IsBlacklisted checks if a token is blacklisted
func (t *TokenBlacklistService) IsBlacklisted(token string) bool {
	expiresAt, exists := t.blacklistedTokens[token]
	if !exists {
		return false
	}

	// Clean up expired blacklisted tokens
	if time.Now().After(expiresAt) {
		delete(t.blacklistedTokens, token)
		return false
	}

	return true
}

// BlacklistAllUserTokens blacklists all tokens for a user
func (t *TokenBlacklistService) BlacklistAllUserTokens(userID uuid.UUID) error {
	// In a real implementation, this would query the database for all user tokens
	// and add them to the blacklist
	return nil
}

// CleanupExpiredTokens removes expired tokens from the blacklist
func (t *TokenBlacklistService) CleanupExpiredTokens() {
	now := time.Now()
	for token, expiresAt := range t.blacklistedTokens {
		if now.After(expiresAt) {
			delete(t.blacklistedTokens, token)
		}
	}
}
