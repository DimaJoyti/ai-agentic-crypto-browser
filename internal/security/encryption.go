package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
)

// EncryptionManager handles all encryption operations
type EncryptionManager struct {
	logger         *observability.Logger
	config         *EncryptionConfig
	keyManager     *KeyManager
	encryptionKeys map[string]*EncryptionKey
	mu             sync.RWMutex
}

// EncryptionConfig contains encryption configuration
type EncryptionConfig struct {
	Algorithm           string        `json:"algorithm"` // AES-256-GCM, RSA-4096
	KeyRotationInterval time.Duration `json:"key_rotation_interval"`
	EnableKeyEscrow     bool          `json:"enable_key_escrow"`
	EnableHSM           bool          `json:"enable_hsm"`
	ComplianceMode      string        `json:"compliance_mode"` // FIPS-140-2, Common Criteria
	EncryptionAtRest    bool          `json:"encryption_at_rest"`
	EncryptionInTransit bool          `json:"encryption_in_transit"`
}

// EncryptionKey represents an encryption key with metadata
type EncryptionKey struct {
	ID          string    `json:"id"`
	Algorithm   string    `json:"algorithm"`
	KeyData     []byte    `json:"-"` // Never serialize key data
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	Purpose     string    `json:"purpose"` // data, transport, signing
	Status      string    `json:"status"`  // active, expired, revoked
	Version     int       `json:"version"`
	Fingerprint string    `json:"fingerprint"`
}

// KeyManager manages encryption keys and rotation
type KeyManager struct {
	logger     *observability.Logger
	config     *EncryptionConfig
	masterKey  []byte
	keyStore   map[string]*EncryptionKey
	rotationCh chan string
	mu         sync.RWMutex
}

// EncryptionResult contains encryption operation result
type EncryptionResult struct {
	EncryptedData []byte            `json:"encrypted_data"`
	KeyID         string            `json:"key_id"`
	Algorithm     string            `json:"algorithm"`
	Metadata      map[string]string `json:"metadata"`
	Timestamp     time.Time         `json:"timestamp"`
}

// DecryptionRequest contains decryption parameters
type DecryptionRequest struct {
	EncryptedData []byte            `json:"encrypted_data"`
	KeyID         string            `json:"key_id"`
	Algorithm     string            `json:"algorithm"`
	Metadata      map[string]string `json:"metadata"`
}

// NewEncryptionManager creates a new encryption manager
func NewEncryptionManager(logger *observability.Logger, config *EncryptionConfig) *EncryptionManager {
	em := &EncryptionManager{
		logger:         logger,
		config:         config,
		encryptionKeys: make(map[string]*EncryptionKey),
	}

	// Initialize key manager
	em.keyManager = NewKeyManager(logger, config)

	return em
}

// Start starts the encryption manager
func (em *EncryptionManager) Start() error {
	em.logger.Info(nil, "Starting encryption manager", map[string]interface{}{
		"algorithm":       em.config.Algorithm,
		"compliance_mode": em.config.ComplianceMode,
		"key_rotation":    em.config.KeyRotationInterval,
	})

	// Start key manager
	if err := em.keyManager.Start(); err != nil {
		return fmt.Errorf("failed to start key manager: %w", err)
	}

	// Generate initial encryption keys
	if err := em.generateInitialKeys(); err != nil {
		return fmt.Errorf("failed to generate initial keys: %w", err)
	}

	return nil
}

// EncryptData encrypts data using the specified algorithm
func (em *EncryptionManager) EncryptData(data []byte, purpose string) (*EncryptionResult, error) {
	em.mu.RLock()
	defer em.mu.RUnlock()

	// Get appropriate key for purpose
	key, err := em.keyManager.GetActiveKey(purpose)
	if err != nil {
		return nil, fmt.Errorf("failed to get encryption key: %w", err)
	}

	switch key.Algorithm {
	case "AES-256-GCM":
		return em.encryptAESGCM(data, key)
	case "RSA-4096":
		return em.encryptRSA(data, key)
	default:
		return nil, fmt.Errorf("unsupported encryption algorithm: %s", key.Algorithm)
	}
}

// DecryptData decrypts data using the specified key
func (em *EncryptionManager) DecryptData(request *DecryptionRequest) ([]byte, error) {
	em.mu.RLock()
	defer em.mu.RUnlock()

	// Get decryption key
	key, err := em.keyManager.GetKey(request.KeyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get decryption key: %w", err)
	}

	switch request.Algorithm {
	case "AES-256-GCM":
		return em.decryptAESGCM(request.EncryptedData, key)
	case "RSA-4096":
		return em.decryptRSA(request.EncryptedData, key)
	default:
		return nil, fmt.Errorf("unsupported decryption algorithm: %s", request.Algorithm)
	}
}

// EncryptPII encrypts personally identifiable information
func (em *EncryptionManager) EncryptPII(piiData map[string]interface{}) (map[string]interface{}, error) {
	encryptedData := make(map[string]interface{})

	for field, value := range piiData {
		if em.isPIIField(field) {
			// Convert value to bytes
			valueBytes := []byte(fmt.Sprintf("%v", value))

			// Encrypt the value
			result, err := em.EncryptData(valueBytes, "pii")
			if err != nil {
				return nil, fmt.Errorf("failed to encrypt PII field %s: %w", field, err)
			}

			// Store encrypted value with metadata
			encryptedData[field] = map[string]interface{}{
				"encrypted_data": base64.StdEncoding.EncodeToString(result.EncryptedData),
				"key_id":         result.KeyID,
				"algorithm":      result.Algorithm,
				"encrypted_at":   result.Timestamp,
			}
		} else {
			// Non-PII data remains unencrypted
			encryptedData[field] = value
		}
	}

	return encryptedData, nil
}

// encryptAESGCM encrypts data using AES-256-GCM
func (em *EncryptionManager) encryptAESGCM(data []byte, key *EncryptionKey) (*EncryptionResult, error) {
	block, err := aes.NewCipher(key.KeyData)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt data
	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return &EncryptionResult{
		EncryptedData: ciphertext,
		KeyID:         key.ID,
		Algorithm:     key.Algorithm,
		Metadata: map[string]string{
			"nonce_size": fmt.Sprintf("%d", gcm.NonceSize()),
		},
		Timestamp: time.Now(),
	}, nil
}

// decryptAESGCM decrypts data using AES-256-GCM
func (em *EncryptionManager) decryptAESGCM(encryptedData []byte, key *EncryptionKey) ([]byte, error) {
	block, err := aes.NewCipher(key.KeyData)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, fmt.Errorf("encrypted data too short")
	}

	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// encryptRSA encrypts data using RSA-4096
func (em *EncryptionManager) encryptRSA(data []byte, key *EncryptionKey) (*EncryptionResult, error) {
	// Parse RSA public key
	block, _ := pem.Decode(key.KeyData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}

	// Encrypt data
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPubKey, data, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt with RSA: %w", err)
	}

	return &EncryptionResult{
		EncryptedData: ciphertext,
		KeyID:         key.ID,
		Algorithm:     key.Algorithm,
		Timestamp:     time.Now(),
	}, nil
}

// decryptRSA decrypts data using RSA-4096
func (em *EncryptionManager) decryptRSA(encryptedData []byte, key *EncryptionKey) ([]byte, error) {
	// Parse RSA private key
	block, _ := pem.Decode(key.KeyData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Decrypt data
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privKey, encryptedData, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt with RSA: %w", err)
	}

	return plaintext, nil
}

// generateInitialKeys generates initial encryption keys
func (em *EncryptionManager) generateInitialKeys() error {
	purposes := []string{"data", "pii", "transport", "signing"}

	for _, purpose := range purposes {
		if err := em.keyManager.GenerateKey(purpose, em.config.Algorithm); err != nil {
			return fmt.Errorf("failed to generate %s key: %w", purpose, err)
		}
	}

	return nil
}

// isPIIField checks if a field contains PII data
func (em *EncryptionManager) isPIIField(field string) bool {
	piiFields := map[string]bool{
		"email":          true,
		"phone":          true,
		"ssn":            true,
		"passport":       true,
		"driver_license": true,
		"address":        true,
		"full_name":      true,
		"date_of_birth":  true,
		"credit_card":    true,
		"bank_account":   true,
	}

	return piiFields[field]
}

// RotateKeys rotates encryption keys
func (em *EncryptionManager) RotateKeys() error {
	em.logger.Info(nil, "Starting key rotation", nil)

	purposes := []string{"data", "pii", "transport", "signing"}
	for _, purpose := range purposes {
		if err := em.keyManager.RotateKey(purpose); err != nil {
			em.logger.Error(nil, "Failed to rotate key", err, map[string]interface{}{
				"purpose": purpose,
			})
			return err
		}
	}

	em.logger.Info(nil, "Key rotation completed successfully", nil)
	return nil
}

// GetEncryptionMetrics returns encryption metrics
func (em *EncryptionManager) GetEncryptionMetrics() map[string]interface{} {
	em.mu.RLock()
	defer em.mu.RUnlock()

	activeKeys := 0
	expiredKeys := 0

	for _, key := range em.encryptionKeys {
		if key.Status == "active" {
			activeKeys++
		} else if key.Status == "expired" {
			expiredKeys++
		}
	}

	return map[string]interface{}{
		"active_keys":           activeKeys,
		"expired_keys":          expiredKeys,
		"total_keys":            len(em.encryptionKeys),
		"algorithm":             em.config.Algorithm,
		"key_rotation_interval": em.config.KeyRotationInterval,
		"compliance_mode":       em.config.ComplianceMode,
		"encryption_at_rest":    em.config.EncryptionAtRest,
		"encryption_in_transit": em.config.EncryptionInTransit,
	}
}

// NewKeyManager creates a new key manager
func NewKeyManager(logger *observability.Logger, config *EncryptionConfig) *KeyManager {
	return &KeyManager{
		logger:     logger,
		config:     config,
		keyStore:   make(map[string]*EncryptionKey),
		rotationCh: make(chan string, 10),
	}
}

// Start starts the key manager
func (km *KeyManager) Start() error {
	km.logger.Info(nil, "Starting key manager", nil)

	// Start key rotation scheduler
	go km.keyRotationScheduler()

	return nil
}

// GenerateKey generates a new encryption key
func (km *KeyManager) GenerateKey(purpose, algorithm string) error {
	km.mu.Lock()
	defer km.mu.Unlock()

	var keyData []byte

	switch algorithm {
	case "AES-256-GCM":
		keyData = make([]byte, 32) // 256 bits
		if _, err := rand.Read(keyData); err != nil {
			return fmt.Errorf("failed to generate AES key: %w", err)
		}
	case "RSA-4096":
		// Generate RSA key pair
		privKey, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return fmt.Errorf("failed to generate RSA key: %w", err)
		}
		keyData = x509.MarshalPKCS1PrivateKey(privKey)
	default:
		return fmt.Errorf("unsupported algorithm: %s", algorithm)
	}

	// Create key metadata
	key := &EncryptionKey{
		ID:          fmt.Sprintf("%s_%s_%d", purpose, algorithm, time.Now().Unix()),
		Algorithm:   algorithm,
		KeyData:     keyData,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(km.config.KeyRotationInterval),
		Purpose:     purpose,
		Status:      "active",
		Version:     1,
		Fingerprint: km.generateFingerprint(keyData),
	}

	km.keyStore[key.ID] = key

	km.logger.Info(nil, "Generated new encryption key", map[string]interface{}{
		"key_id":      key.ID,
		"purpose":     purpose,
		"algorithm":   algorithm,
		"fingerprint": key.Fingerprint,
	})

	return nil
}

// GetActiveKey returns the active key for a purpose
func (km *KeyManager) GetActiveKey(purpose string) (*EncryptionKey, error) {
	km.mu.RLock()
	defer km.mu.RUnlock()

	for _, key := range km.keyStore {
		if key.Purpose == purpose && key.Status == "active" && time.Now().Before(key.ExpiresAt) {
			return key, nil
		}
	}

	return nil, fmt.Errorf("no active key found for purpose: %s", purpose)
}

// GetKey returns a key by ID
func (km *KeyManager) GetKey(keyID string) (*EncryptionKey, error) {
	km.mu.RLock()
	defer km.mu.RUnlock()

	key, exists := km.keyStore[keyID]
	if !exists {
		return nil, fmt.Errorf("key not found: %s", keyID)
	}

	return key, nil
}

// RotateKey rotates a key for a specific purpose
func (km *KeyManager) RotateKey(purpose string) error {
	// Mark old key as expired
	km.mu.Lock()
	for _, key := range km.keyStore {
		if key.Purpose == purpose && key.Status == "active" {
			key.Status = "expired"
		}
	}
	km.mu.Unlock()

	// Generate new key
	return km.GenerateKey(purpose, km.config.Algorithm)
}

// keyRotationScheduler schedules automatic key rotation
func (km *KeyManager) keyRotationScheduler() {
	ticker := time.NewTicker(km.config.KeyRotationInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			km.logger.Info(nil, "Automatic key rotation triggered", nil)
			purposes := []string{"data", "pii", "transport", "signing"}
			for _, purpose := range purposes {
				if err := km.RotateKey(purpose); err != nil {
					km.logger.Error(nil, "Failed to rotate key", err, map[string]interface{}{
						"purpose": purpose,
					})
				}
			}
		case purpose := <-km.rotationCh:
			km.logger.Info(nil, "Manual key rotation triggered", map[string]interface{}{
				"purpose": purpose,
			})
			if err := km.RotateKey(purpose); err != nil {
				km.logger.Error(nil, "Failed to rotate key", err, map[string]interface{}{
					"purpose": purpose,
				})
			}
		}
	}
}

// generateFingerprint generates a fingerprint for a key
func (km *KeyManager) generateFingerprint(keyData []byte) string {
	hash := sha256.Sum256(keyData)
	return fmt.Sprintf("%x", hash[:8]) // First 8 bytes as hex
}
