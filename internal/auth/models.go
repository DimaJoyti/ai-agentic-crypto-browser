package auth

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	Email             string     `json:"email" db:"email"`
	Password          string     `json:"-" db:"password_hash"`
	FirstName         *string    `json:"first_name" db:"first_name"`
	LastName          *string    `json:"last_name" db:"last_name"`
	Role              string     `json:"role" db:"role"`
	Permissions       []string   `json:"permissions" db:"permissions"`
	TeamID            *uuid.UUID `json:"team_id" db:"team_id"`
	MFAEnabled        bool       `json:"mfa_enabled" db:"mfa_enabled"`
	MFASecret         string     `json:"-" db:"mfa_secret"`
	MFAVerified       bool       `json:"mfa_verified" db:"mfa_verified"`
	FailedLoginCount  int        `json:"failed_login_count" db:"failed_login_count"`
	LockedUntil       *time.Time `json:"locked_until" db:"locked_until"`
	LastLoginAt       *time.Time `json:"last_login_at" db:"last_login_at"`
	LastLoginIP       *string    `json:"last_login_ip" db:"last_login_ip"`
	PasswordChangedAt *time.Time `json:"password_changed_at" db:"password_changed_at"`
	IsActive          bool       `json:"is_active" db:"is_active"`
	IsVerified        bool       `json:"is_verified" db:"is_verified"`
	IsEmailVerified   bool       `json:"is_email_verified" db:"is_email_verified"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

// UserSession represents a user session with refresh token
type UserSession struct {
	ID               uuid.UUID `json:"id" db:"id"`
	UserID           uuid.UUID `json:"user_id" db:"user_id"`
	RefreshTokenHash string    `json:"-" db:"refresh_token_hash"`
	ExpiresAt        time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	LastUsedAt       time.Time `json:"last_used_at" db:"last_used_at"`
	UserAgent        *string   `json:"user_agent" db:"user_agent"`
	IPAddress        *string   `json:"ip_address" db:"ip_address"`
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email     string  `json:"email" validate:"required,email"`
	Password  string  `json:"password" validate:"required,min=8"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents a successful login response
type LoginResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// UpdateProfileRequest represents a profile update request
type UpdateProfileRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
