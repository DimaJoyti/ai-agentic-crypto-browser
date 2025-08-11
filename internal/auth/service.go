package auth

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/ai-agentic-browser/internal/config"
	"github.com/ai-agentic-browser/pkg/database"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
)

// Context keys for storing user information
type contextKey string

const (
	userContextKey   contextKey = "user"
	userIDContextKey contextKey = "user_id"
)

// UserFromContext retrieves the user from the context
func UserFromContext(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(userContextKey).(*User)
	return user, ok
}

// UserIDFromContext retrieves the user ID from the context
func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userIDContextKey).(uuid.UUID)
	return userID, ok
}

// WithUser adds a user to the context
func WithUser(ctx context.Context, user *User) context.Context {
	ctx = context.WithValue(ctx, userContextKey, user)
	ctx = context.WithValue(ctx, userIDContextKey, user.ID)
	return ctx
}

// WithUserID adds a user ID to the context
func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

// Service provides authentication functionality
type Service struct {
	db             *database.DB
	redis          *database.RedisClient
	config         config.JWTConfig
	logger         *observability.Logger
	jwtService     *JWTService
	mfaService     *MFAService
	rbacService    *RBACService
	securityConfig *SecurityConfig
}

// SecurityConfig contains security configuration
type SecurityConfig struct {
	BCryptCost         int
	Argon2Time         uint32
	Argon2Memory       uint32
	Argon2Threads      uint8
	Argon2KeyLen       uint32
	PasswordMinLength  int
	PasswordMaxLength  int
	SessionTimeout     time.Duration
	MaxLoginAttempts   int
	LockoutDuration    time.Duration
	RequireMFA         bool
	AllowedDomains     []string
	PasswordComplexity PasswordComplexityConfig
}

// PasswordComplexityConfig defines password complexity requirements
type PasswordComplexityConfig struct {
	RequireUppercase bool
	RequireLowercase bool
	RequireNumbers   bool
	RequireSymbols   bool
	MinLength        int
	MaxLength        int
	ForbiddenWords   []string
}

// NewService creates a new authentication service
func NewService(db *database.DB, redis *database.RedisClient, cfg config.JWTConfig, logger *observability.Logger) *Service {
	return &Service{
		db:     db,
		redis:  redis,
		config: cfg,
		logger: logger,
	}
}

// Register creates a new user account
func (s *Service) Register(ctx context.Context, req RegisterRequest) (*User, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("auth-service").Start(ctx, "auth.Register")
	defer span.End()

	// Check if user already exists
	existingUser, err := s.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error(ctx, "Failed to hash password", err)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &User{
		ID:        uuid.New(),
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert user into database
	query := `
		INSERT INTO users (id, email, password_hash, first_name, last_name, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = s.db.ExecContext(ctx, query, user.ID, user.Email, user.Password, user.FirstName, user.LastName, user.IsActive, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		s.logger.Error(ctx, "Failed to create user", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	s.logger.Info(ctx, "User registered successfully", map[string]interface{}{
		"user_id": user.ID.String(),
		"email":   user.Email,
	})

	// Clear password before returning
	user.Password = ""
	return user, nil
}

// Login authenticates a user and returns tokens
func (s *Service) Login(ctx context.Context, req LoginRequest, userAgent, ipAddress string) (*LoginResponse, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("auth-service").Start(ctx, "auth.Login")
	defer span.End()

	// Get user by email
	user, err := s.GetUserByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Warn(ctx, "Login attempt with non-existent email", map[string]interface{}{
			"email": req.Email,
		})
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		s.logger.Warn(ctx, "Login attempt with invalid password", map[string]interface{}{
			"user_id": user.ID.String(),
			"email":   user.Email,
		})
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		s.logger.Error(ctx, "Failed to generate access token", err)
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		s.logger.Error(ctx, "Failed to generate refresh token", err)
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store refresh token session
	session := &UserSession{
		ID:               uuid.New(),
		UserID:           user.ID,
		RefreshTokenHash: s.hashToken(refreshToken),
		ExpiresAt:        time.Now().Add(s.config.RefreshTokenExpiry),
		CreatedAt:        time.Now(),
		LastUsedAt:       time.Now(),
		UserAgent:        &userAgent,
		IPAddress:        &ipAddress,
	}

	err = s.createSession(ctx, session)
	if err != nil {
		s.logger.Error(ctx, "Failed to create session", err)
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	s.logger.Info(ctx, "User logged in successfully", map[string]interface{}{
		"user_id":    user.ID.String(),
		"email":      user.Email,
		"session_id": session.ID.String(),
	})

	// Clear password before returning
	user.Password = ""

	return &LoginResponse{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.config.Expiry.Seconds()),
	}, nil
}

// RefreshToken generates new tokens using a refresh token
func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("auth-service").Start(ctx, "auth.RefreshToken")
	defer span.End()

	// Find session by refresh token hash
	tokenHash := s.hashToken(refreshToken)
	session, err := s.getSessionByRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		s.deleteSession(ctx, session.ID)
		return nil, fmt.Errorf("refresh token expired")
	}

	// Get user
	user, err := s.GetUserByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check if user is still active
	if !user.IsActive {
		s.deleteSession(ctx, session.ID)
		return nil, fmt.Errorf("account is deactivated")
	}

	// Generate new access token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		s.logger.Error(ctx, "Failed to generate access token", err)
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Update session last used time
	session.LastUsedAt = time.Now()
	s.updateSessionLastUsed(ctx, session.ID, session.LastUsedAt)

	// Clear password before returning
	user.Password = ""

	return &LoginResponse{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.config.Expiry.Seconds()),
	}, nil
}

// Logout invalidates a refresh token
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("auth-service").Start(ctx, "auth.Logout")
	defer span.End()

	tokenHash := s.hashToken(refreshToken)
	session, err := s.getSessionByRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil // Already logged out or invalid token
	}

	err = s.deleteSession(ctx, session.ID)
	if err != nil {
		s.logger.Error(ctx, "Failed to delete session", err)
		return fmt.Errorf("failed to logout: %w", err)
	}

	s.logger.Info(ctx, "User logged out successfully", map[string]interface{}{
		"user_id":    session.UserID.String(),
		"session_id": session.ID.String(),
	})

	return nil
}

// GetUserByID retrieves a user by ID
func (s *Service) GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, is_active, is_verified, created_at, updated_at
		FROM users WHERE id = $1
	`
	user := &User{}
	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.IsActive, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *Service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, is_active, is_verified, created_at, updated_at
		FROM users WHERE email = $1
	`
	user := &User{}
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.IsActive, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

// generateAccessToken creates a new JWT access token
func (s *Service) generateAccessToken(user *User) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"iat":     now.Unix(),
		"exp":     now.Add(s.config.Expiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.Secret))
}

// generateRefreshToken creates a new refresh token
func (s *Service) generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// hashToken creates a hash of a token for storage
func (s *Service) hashToken(token string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	return string(hash)
}

// Helper methods for session management
func (s *Service) createSession(ctx context.Context, session *UserSession) error {
	query := `
		INSERT INTO user_sessions (id, user_id, refresh_token_hash, expires_at, created_at, last_used_at, user_agent, ip_address)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := s.db.ExecContext(ctx, query, session.ID, session.UserID, session.RefreshTokenHash,
		session.ExpiresAt, session.CreatedAt, session.LastUsedAt, session.UserAgent, session.IPAddress)
	return err
}

func (s *Service) getSessionByRefreshToken(ctx context.Context, tokenHash string) (*UserSession, error) {
	query := `
		SELECT id, user_id, refresh_token_hash, expires_at, created_at, last_used_at, user_agent, ip_address
		FROM user_sessions WHERE refresh_token_hash = $1
	`
	session := &UserSession{}
	err := s.db.QueryRowContext(ctx, query, tokenHash).Scan(
		&session.ID, &session.UserID, &session.RefreshTokenHash, &session.ExpiresAt,
		&session.CreatedAt, &session.LastUsedAt, &session.UserAgent, &session.IPAddress,
	)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *Service) updateSessionLastUsed(ctx context.Context, sessionID uuid.UUID, lastUsed time.Time) error {
	query := `UPDATE user_sessions SET last_used_at = $1 WHERE id = $2`
	_, err := s.db.ExecContext(ctx, query, lastUsed, sessionID)
	return err
}

func (s *Service) deleteSession(ctx context.Context, sessionID uuid.UUID) error {
	query := `DELETE FROM user_sessions WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, sessionID)
	return err
}

// Advanced Security Methods

// HashPassword hashes a password using Argon2id for enhanced security
func (s *Service) HashPassword(password string) (string, error) {
	// Validate password complexity
	if err := s.validatePasswordComplexity(password); err != nil {
		return "", err
	}

	// Generate salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Use Argon2id for password hashing (more secure than bcrypt)
	if s.securityConfig != nil {
		hash := argon2.IDKey([]byte(password), salt, s.securityConfig.Argon2Time,
			s.securityConfig.Argon2Memory, s.securityConfig.Argon2Threads, s.securityConfig.Argon2KeyLen)

		encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
			argon2.Version, s.securityConfig.Argon2Memory, s.securityConfig.Argon2Time,
			s.securityConfig.Argon2Threads, hex.EncodeToString(salt), hex.EncodeToString(hash))

		return encoded, nil
	}

	// Fallback to bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// VerifyPassword verifies a password against its hash (supports both Argon2id and bcrypt)
func (s *Service) VerifyPassword(password, hash string) bool {
	if len(hash) > 0 && hash[0] == '$' {
		if strings.HasPrefix(hash, "$argon2id") {
			return s.verifyArgon2Password(password, hash)
		} else if strings.HasPrefix(hash, "$2a$") || strings.HasPrefix(hash, "$2b$") || strings.HasPrefix(hash, "$2y$") {
			return s.verifyBcryptPassword(password, hash)
		}
	}
	return false
}

// CreateUser creates a new user with enhanced security
func (s *Service) CreateUser(ctx context.Context, user *User) error {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("auth-service").Start(ctx, "auth.CreateUser")
	defer span.End()

	// Check if user already exists
	existingUser, err := s.GetUserByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return fmt.Errorf("user with email %s already exists", user.Email)
	}

	// Validate email domain if restrictions are configured
	if s.securityConfig != nil && len(s.securityConfig.AllowedDomains) > 0 {
		if !s.isEmailDomainAllowed(user.Email) {
			return fmt.Errorf("email domain not allowed")
		}
	}

	// Set default values
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Insert user into database
	query := `
		INSERT INTO users (id, email, password_hash, first_name, last_name, role, permissions, team_id,
						  mfa_enabled, mfa_secret, is_active, is_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err = s.db.ExecContext(ctx, query,
		user.ID, user.Email, user.Password, user.FirstName, user.LastName,
		user.Role, user.Permissions, user.TeamID, user.MFAEnabled, user.MFASecret,
		user.IsActive, user.IsVerified, user.CreatedAt, user.UpdatedAt,
	)

	if err != nil {
		s.logger.Error(ctx, "Failed to create user", err)
		return fmt.Errorf("failed to create user: %w", err)
	}

	s.logger.Info(ctx, "User created successfully", map[string]interface{}{
		"user_id": user.ID.String(),
		"email":   user.Email,
	})

	return nil
}

// AuthenticateUser authenticates a user with enhanced security checks
func (s *Service) AuthenticateUser(ctx context.Context, email, password string) (*User, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("auth-service").Start(ctx, "auth.AuthenticateUser")
	defer span.End()

	// Get user by email
	user, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		s.logLoginAttempt(ctx, email, "", "", false, "user not found")
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if account is locked
	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		s.logLoginAttempt(ctx, email, "", "", false, "account locked")
		return nil, fmt.Errorf("account is locked")
	}

	// Check if account is active
	if !user.IsActive {
		s.logLoginAttempt(ctx, email, "", "", false, "account inactive")
		return nil, fmt.Errorf("account is inactive")
	}

	// Verify password
	if !s.VerifyPassword(password, user.Password) {
		// Increment failed login count
		s.incrementFailedLoginCount(ctx, user.ID)
		s.logLoginAttempt(ctx, email, "", "", false, "invalid password")
		return nil, fmt.Errorf("invalid credentials")
	}

	// Reset failed login count on successful authentication
	s.resetFailedLoginCount(ctx, user.ID)

	// Update last login
	s.updateLastLogin(ctx, user.ID, "")

	s.logLoginAttempt(ctx, email, "", "", true, "success")

	return user, nil
}

// Helper methods for enhanced security

func (s *Service) validatePasswordComplexity(password string) error {
	if s.securityConfig == nil {
		return nil
	}

	config := s.securityConfig.PasswordComplexity

	if len(password) < config.MinLength {
		return fmt.Errorf("password must be at least %d characters long", config.MinLength)
	}

	if len(password) > config.MaxLength {
		return fmt.Errorf("password must be no more than %d characters long", config.MaxLength)
	}

	if config.RequireUppercase && !containsUppercase(password) {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}

	if config.RequireLowercase && !containsLowercase(password) {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}

	if config.RequireNumbers && !containsNumber(password) {
		return fmt.Errorf("password must contain at least one number")
	}

	if config.RequireSymbols && !containsSymbol(password) {
		return fmt.Errorf("password must contain at least one symbol")
	}

	// Check forbidden words
	for _, word := range config.ForbiddenWords {
		if strings.Contains(strings.ToLower(password), strings.ToLower(word)) {
			return fmt.Errorf("password contains forbidden word")
		}
	}

	return nil
}

func (s *Service) verifyArgon2Password(password, hash string) bool {
	// Parse Argon2id hash format: $argon2id$v=19$m=65536,t=1,p=4$salt$hash
	parts := strings.Split(hash, "$")
	if len(parts) != 6 {
		return false
	}

	// Extract parameters (simplified implementation)
	saltBytes, err := hex.DecodeString(parts[4])
	if err != nil {
		return false
	}

	hashBytes, err := hex.DecodeString(parts[5])
	if err != nil {
		return false
	}

	// Use default parameters if securityConfig is not available
	time := uint32(1)
	memory := uint32(64 * 1024)
	threads := uint8(4)
	keyLen := uint32(32)

	if s.securityConfig != nil {
		time = s.securityConfig.Argon2Time
		memory = s.securityConfig.Argon2Memory
		threads = s.securityConfig.Argon2Threads
		keyLen = s.securityConfig.Argon2KeyLen
	}

	// Generate hash with same parameters
	computedHash := argon2.IDKey([]byte(password), saltBytes, time, memory, threads, keyLen)

	// Use constant time comparison to prevent timing attacks
	return subtle.ConstantTimeCompare(hashBytes, computedHash) == 1
}

func (s *Service) verifyBcryptPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *Service) isEmailDomainAllowed(email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	domain := parts[1]

	for _, allowedDomain := range s.securityConfig.AllowedDomains {
		if domain == allowedDomain {
			return true
		}
	}
	return false
}

func (s *Service) incrementFailedLoginCount(ctx context.Context, userID uuid.UUID) {
	if s.securityConfig == nil {
		return
	}

	query := `
		UPDATE users
		SET failed_login_count = failed_login_count + 1,
			locked_until = CASE
				WHEN failed_login_count + 1 >= $1 THEN $2
				ELSE locked_until
			END,
			updated_at = $3
		WHERE id = $4
	`
	lockUntil := time.Now().Add(s.securityConfig.LockoutDuration)
	s.db.ExecContext(ctx, query, s.securityConfig.MaxLoginAttempts, lockUntil, time.Now(), userID)
}

func (s *Service) resetFailedLoginCount(ctx context.Context, userID uuid.UUID) {
	query := `UPDATE users SET failed_login_count = 0, locked_until = NULL, updated_at = $1 WHERE id = $2`
	s.db.ExecContext(ctx, query, time.Now(), userID)
}

func (s *Service) updateLastLogin(ctx context.Context, userID uuid.UUID, ipAddress string) {
	query := `UPDATE users SET last_login_at = $1, last_login_ip = $2, updated_at = $3 WHERE id = $4`
	s.db.ExecContext(ctx, query, time.Now(), ipAddress, time.Now(), userID)
}

func (s *Service) logLoginAttempt(ctx context.Context, email, ipAddress, userAgent string, success bool, reason string) {
	// In a real implementation, this would log to an audit table
	s.logger.Info(ctx, "Login attempt", map[string]interface{}{
		"email":      email,
		"ip_address": ipAddress,
		"user_agent": userAgent,
		"success":    success,
		"reason":     reason,
	})
}

// Password validation helper functions
func containsUppercase(s string) bool {
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			return true
		}
	}
	return false
}

func containsLowercase(s string) bool {
	for _, r := range s {
		if r >= 'a' && r <= 'z' {
			return true
		}
	}
	return false
}

func containsNumber(s string) bool {
	for _, r := range s {
		if r >= '0' && r <= '9' {
			return true
		}
	}
	return false
}

func containsSymbol(s string) bool {
	symbols := "!@#$%^&*()_+-=[]{}|;:,.<>?"
	for _, r := range s {
		for _, sym := range symbols {
			if r == sym {
				return true
			}
		}
	}
	return false
}
