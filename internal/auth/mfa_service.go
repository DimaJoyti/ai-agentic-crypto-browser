package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

// MFAService handles multi-factor authentication
type MFAService struct {
	issuer      string
	backupCodes *BackupCodeService
}

// MFAMethod represents different MFA methods
type MFAMethod string

const (
	MFAMethodTOTP     MFAMethod = "totp"
	MFAMethodSMS      MFAMethod = "sms"
	MFAMethodEmail    MFAMethod = "email"
	MFAMethodBackup   MFAMethod = "backup"
	MFAMethodWebAuthn MFAMethod = "webauthn"
)

// MFASetup contains MFA setup information
type MFASetup struct {
	Secret      string    `json:"secret"`
	QRCodeURL   string    `json:"qr_code_url"`
	QRCodeImage []byte    `json:"qr_code_image"`
	BackupCodes []string  `json:"backup_codes"`
	Method      MFAMethod `json:"method"`
}

// MFAChallenge represents an MFA challenge
type MFAChallenge struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Method      MFAMethod `json:"method"`
	Challenge   string    `json:"challenge"`
	ExpiresAt   time.Time `json:"expires_at"`
	Attempts    int       `json:"attempts"`
	MaxAttempts int       `json:"max_attempts"`
	CreatedAt   time.Time `json:"created_at"`
}

// MFAVerification represents MFA verification request
type MFAVerification struct {
	UserID    uuid.UUID `json:"user_id"`
	Method    MFAMethod `json:"method"`
	Code      string    `json:"code"`
	Challenge string    `json:"challenge,omitempty"`
}

// NewMFAService creates a new MFA service
func NewMFAService(issuer string) *MFAService {
	return &MFAService{
		issuer:      issuer,
		backupCodes: NewBackupCodeService(),
	}
}

// SetupTOTP generates TOTP secret and QR code for user
func (m *MFAService) SetupTOTP(userID uuid.UUID, email string) (*MFASetup, error) {
	// Generate TOTP secret
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      m.issuer,
		AccountName: email,
		SecretSize:  32,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate TOTP key: %w", err)
	}

	// Generate QR code
	qrCode, err := qrcode.Encode(key.URL(), qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Generate backup codes
	backupCodes, err := m.backupCodes.GenerateBackupCodes(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate backup codes: %w", err)
	}

	return &MFASetup{
		Secret:      key.Secret(),
		QRCodeURL:   key.URL(),
		QRCodeImage: qrCode,
		BackupCodes: backupCodes,
		Method:      MFAMethodTOTP,
	}, nil
}

// VerifyTOTP verifies a TOTP code
func (m *MFAService) VerifyTOTP(secret, code string) bool {
	return totp.Validate(code, secret)
}

// VerifyTOTPWithWindow verifies TOTP with time window tolerance
func (m *MFAService) VerifyTOTPWithWindow(secret, code string, window int) bool {
	// Check current time and adjacent windows
	now := time.Now()
	for i := -window; i <= window; i++ {
		testTime := now.Add(time.Duration(i) * 30 * time.Second)
		valid, err := totp.ValidateCustom(code, secret, testTime, totp.ValidateOpts{
			Period:    30,
			Skew:      0,
			Digits:    otp.DigitsSix,
			Algorithm: otp.AlgorithmSHA1,
		})
		if err == nil && valid {
			return true
		}
	}
	return false
}

// CreateSMSChallenge creates an SMS-based MFA challenge
func (m *MFAService) CreateSMSChallenge(userID uuid.UUID, phoneNumber string) (*MFAChallenge, error) {
	// Generate 6-digit code
	code, err := m.generateNumericCode(6)
	if err != nil {
		return nil, fmt.Errorf("failed to generate SMS code: %w", err)
	}

	challenge := &MFAChallenge{
		ID:          uuid.New(),
		UserID:      userID,
		Method:      MFAMethodSMS,
		Challenge:   code,
		ExpiresAt:   time.Now().Add(5 * time.Minute),
		Attempts:    0,
		MaxAttempts: 3,
		CreatedAt:   time.Now(),
	}

	// In a real implementation, send SMS here
	// smsService.SendSMS(phoneNumber, fmt.Sprintf("Your verification code is: %s", code))

	return challenge, nil
}

// CreateEmailChallenge creates an email-based MFA challenge
func (m *MFAService) CreateEmailChallenge(userID uuid.UUID, email string) (*MFAChallenge, error) {
	// Generate 8-digit code
	code, err := m.generateNumericCode(8)
	if err != nil {
		return nil, fmt.Errorf("failed to generate email code: %w", err)
	}

	challenge := &MFAChallenge{
		ID:          uuid.New(),
		UserID:      userID,
		Method:      MFAMethodEmail,
		Challenge:   code,
		ExpiresAt:   time.Now().Add(10 * time.Minute),
		Attempts:    0,
		MaxAttempts: 5,
		CreatedAt:   time.Now(),
	}

	// In a real implementation, send email here
	// emailService.SendMFACode(email, code)

	return challenge, nil
}

// VerifyChallenge verifies an MFA challenge
func (m *MFAService) VerifyChallenge(challenge *MFAChallenge, code string) (bool, error) {
	// Check if challenge has expired
	if time.Now().After(challenge.ExpiresAt) {
		return false, fmt.Errorf("challenge has expired")
	}

	// Check if max attempts exceeded
	if challenge.Attempts >= challenge.MaxAttempts {
		return false, fmt.Errorf("maximum attempts exceeded")
	}

	// Increment attempts
	challenge.Attempts++

	// Verify code based on method
	switch challenge.Method {
	case MFAMethodSMS, MFAMethodEmail:
		return challenge.Challenge == code, nil
	case MFAMethodBackup:
		return m.backupCodes.VerifyBackupCode(challenge.UserID, code), nil
	default:
		return false, fmt.Errorf("unsupported MFA method: %s", challenge.Method)
	}
}

// VerifyBackupCode verifies a backup code
func (m *MFAService) VerifyBackupCode(userID uuid.UUID, code string) bool {
	return m.backupCodes.VerifyBackupCode(userID, code)
}

// RegenerateBackupCodes generates new backup codes for a user
func (m *MFAService) RegenerateBackupCodes(userID uuid.UUID) ([]string, error) {
	return m.backupCodes.RegenerateBackupCodes(userID)
}

// GetUserMFAMethods returns enabled MFA methods for a user
func (m *MFAService) GetUserMFAMethods(userID uuid.UUID) ([]MFAMethod, error) {
	// In a real implementation, this would query the database
	// For now, return mock data
	return []MFAMethod{MFAMethodTOTP, MFAMethodBackup}, nil
}

// DisableMFA disables MFA for a user
func (m *MFAService) DisableMFA(userID uuid.UUID, method MFAMethod) error {
	// In a real implementation, this would update the database
	// and revoke all active sessions
	return nil
}

// generateNumericCode generates a random numeric code
func (m *MFAService) generateNumericCode(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	code := ""
	for _, b := range bytes {
		code += strconv.Itoa(int(b) % 10)
	}

	return code, nil
}

// BackupCodeService manages backup codes
type BackupCodeService struct {
	// In a real implementation, this would use a database
	userBackupCodes map[uuid.UUID][]string
}

// NewBackupCodeService creates a new backup code service
func NewBackupCodeService() *BackupCodeService {
	return &BackupCodeService{
		userBackupCodes: make(map[uuid.UUID][]string),
	}
}

// GenerateBackupCodes generates backup codes for a user
func (b *BackupCodeService) GenerateBackupCodes(userID uuid.UUID) ([]string, error) {
	codes := make([]string, 10) // Generate 10 backup codes

	for i := 0; i < 10; i++ {
		code, err := b.generateBackupCode()
		if err != nil {
			return nil, fmt.Errorf("failed to generate backup code: %w", err)
		}
		codes[i] = code
	}

	// Store hashed versions
	hashedCodes := make([]string, len(codes))
	for i, code := range codes {
		hashedCodes[i] = b.hashCode(code)
	}

	b.userBackupCodes[userID] = hashedCodes

	return codes, nil
}

// VerifyBackupCode verifies a backup code and marks it as used
func (b *BackupCodeService) VerifyBackupCode(userID uuid.UUID, code string) bool {
	codes, exists := b.userBackupCodes[userID]
	if !exists {
		return false
	}

	hashedCode := b.hashCode(code)

	// Find and remove the code if it exists
	for i, storedCode := range codes {
		if storedCode == hashedCode {
			// Remove the used code
			b.userBackupCodes[userID] = append(codes[:i], codes[i+1:]...)
			return true
		}
	}

	return false
}

// RegenerateBackupCodes generates new backup codes, invalidating old ones
func (b *BackupCodeService) RegenerateBackupCodes(userID uuid.UUID) ([]string, error) {
	// Clear existing codes
	delete(b.userBackupCodes, userID)

	// Generate new codes
	return b.GenerateBackupCodes(userID)
}

// generateBackupCode generates a single backup code
func (b *BackupCodeService) generateBackupCode() (string, error) {
	bytes := make([]byte, 5)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Convert to base32 and format
	code := base32.StdEncoding.EncodeToString(bytes)
	code = strings.ToLower(code)
	code = strings.TrimRight(code, "=")

	// Format as XXXX-XXXX
	if len(code) >= 8 {
		code = code[:4] + "-" + code[4:8]
	}

	return code, nil
}

// hashCode creates a hash of the backup code
func (b *BackupCodeService) hashCode(code string) string {
	hash := sha256.Sum256([]byte(code))
	return hex.EncodeToString(hash[:])
}

// WebAuthnService handles WebAuthn/FIDO2 authentication
type WebAuthnService struct {
	rpID     string
	rpName   string
	rpOrigin string
}

// NewWebAuthnService creates a new WebAuthn service
func NewWebAuthnService(rpID, rpName, rpOrigin string) *WebAuthnService {
	return &WebAuthnService{
		rpID:     rpID,
		rpName:   rpName,
		rpOrigin: rpOrigin,
	}
}

// WebAuthnCredential represents a WebAuthn credential
type WebAuthnCredential struct {
	ID              []byte        `json:"id"`
	PublicKey       []byte        `json:"public_key"`
	AttestationType string        `json:"attestation_type"`
	Transport       []string      `json:"transport"`
	Flags           byte          `json:"flags"`
	Authenticator   Authenticator `json:"authenticator"`
	CreatedAt       time.Time     `json:"created_at"`
}

// Authenticator represents authenticator information
type Authenticator struct {
	AAGUID       []byte `json:"aaguid"`
	SignCount    uint32 `json:"sign_count"`
	CloneWarning bool   `json:"clone_warning"`
}

// BeginRegistration starts WebAuthn registration process
func (w *WebAuthnService) BeginRegistration(userID uuid.UUID, username string) ([]byte, error) {
	// In a real implementation, this would use the webauthn library
	// to create a credential creation options
	return []byte("mock-registration-challenge"), nil
}

// FinishRegistration completes WebAuthn registration
func (w *WebAuthnService) FinishRegistration(userID uuid.UUID, response []byte) (*WebAuthnCredential, error) {
	// In a real implementation, this would verify the attestation response
	// and store the credential
	return &WebAuthnCredential{
		ID:              []byte("mock-credential-id"),
		PublicKey:       []byte("mock-public-key"),
		AttestationType: "none",
		Transport:       []string{"usb", "nfc"},
		CreatedAt:       time.Now(),
	}, nil
}

// BeginAuthentication starts WebAuthn authentication process
func (w *WebAuthnService) BeginAuthentication(userID uuid.UUID) ([]byte, error) {
	// In a real implementation, this would create assertion options
	return []byte("mock-authentication-challenge"), nil
}

// FinishAuthentication completes WebAuthn authentication
func (w *WebAuthnService) FinishAuthentication(userID uuid.UUID, response []byte) (bool, error) {
	// In a real implementation, this would verify the assertion response
	return true, nil
}

// MFARecoveryService handles MFA recovery scenarios
type MFARecoveryService struct {
	backupCodes *BackupCodeService
}

// NewMFARecoveryService creates a new MFA recovery service
func NewMFARecoveryService(backupCodes *BackupCodeService) *MFARecoveryService {
	return &MFARecoveryService{
		backupCodes: backupCodes,
	}
}

// RecoveryMethod represents different recovery methods
type RecoveryMethod string

const (
	RecoveryMethodBackupCode RecoveryMethod = "backup_code"
	RecoveryMethodAdminReset RecoveryMethod = "admin_reset"
	RecoveryMethodIdentity   RecoveryMethod = "identity_verification"
)

// InitiateRecovery starts the MFA recovery process
func (r *MFARecoveryService) InitiateRecovery(userID uuid.UUID, method RecoveryMethod) error {
	switch method {
	case RecoveryMethodBackupCode:
		// User will provide backup code
		return nil
	case RecoveryMethodAdminReset:
		// Admin intervention required
		return r.requestAdminReset(userID)
	case RecoveryMethodIdentity:
		// Identity verification process
		return r.initiateIdentityVerification(userID)
	default:
		return fmt.Errorf("unsupported recovery method: %s", method)
	}
}

// requestAdminReset requests admin intervention for MFA reset
func (r *MFARecoveryService) requestAdminReset(userID uuid.UUID) error {
	// In a real implementation, this would create a support ticket
	// or notify administrators
	return nil
}

// initiateIdentityVerification starts identity verification process
func (r *MFARecoveryService) initiateIdentityVerification(userID uuid.UUID) error {
	// In a real implementation, this would integrate with identity
	// verification services like Jumio, Onfido, etc.
	return nil
}
