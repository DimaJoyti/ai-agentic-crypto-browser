package security

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

// PrivacyManager handles data privacy and protection
type PrivacyManager struct {
	logger            *observability.Logger
	config            *PrivacyConfig
	encryptionManager *EncryptionManager
	consentManager    *ConsentManager
	dataProcessor     *DataProcessor
	retentionManager  *RetentionManager
	anonymizer        *DataAnonymizer
	mu                sync.RWMutex
}

// PrivacyConfig contains privacy configuration
type PrivacyConfig struct {
	EnableGDPRCompliance    bool          `json:"enable_gdpr_compliance"`
	EnableCCPACompliance    bool          `json:"enable_ccpa_compliance"`
	EnableDataMinimization  bool          `json:"enable_data_minimization"`
	EnablePurposeLimitation bool          `json:"enable_purpose_limitation"`
	DefaultRetentionPeriod  time.Duration `json:"default_retention_period"`
	AnonymizationThreshold  time.Duration `json:"anonymization_threshold"`
	ConsentExpirationPeriod time.Duration `json:"consent_expiration_period"`
	EnableRightToErasure    bool          `json:"enable_right_to_erasure"`
	EnableDataPortability   bool          `json:"enable_data_portability"`
	EnableAutomaticDeletion bool          `json:"enable_automatic_deletion"`
}

// ConsentManager manages user consent
type ConsentManager struct {
	logger   *observability.Logger
	consents map[string]*UserConsent
	mu       sync.RWMutex
}

// UserConsent represents user consent information
type UserConsent struct {
	UserID      uuid.UUID                  `json:"user_id"`
	ConsentID   string                     `json:"consent_id"`
	Purposes    map[string]*ConsentPurpose `json:"purposes"`
	GrantedAt   time.Time                  `json:"granted_at"`
	ExpiresAt   time.Time                  `json:"expires_at"`
	WithdrawnAt *time.Time                 `json:"withdrawn_at,omitempty"`
	Version     string                     `json:"version"`
	IPAddress   string                     `json:"ip_address"`
	UserAgent   string                     `json:"user_agent"`
	LegalBasis  string                     `json:"legal_basis"`
	Status      ConsentStatus              `json:"status"`
}

// ConsentPurpose represents a specific purpose for data processing
type ConsentPurpose struct {
	Purpose     string    `json:"purpose"`
	Description string    `json:"description"`
	Granted     bool      `json:"granted"`
	GrantedAt   time.Time `json:"granted_at"`
	Required    bool      `json:"required"`
	Category    string    `json:"category"` // essential, functional, analytics, marketing
}

// ConsentStatus represents consent status
type ConsentStatus string

const (
	ConsentStatusActive    ConsentStatus = "active"
	ConsentStatusExpired   ConsentStatus = "expired"
	ConsentStatusWithdrawn ConsentStatus = "withdrawn"
	ConsentStatusPending   ConsentStatus = "pending"
)

// DataProcessor handles data processing operations
type DataProcessor struct {
	logger            *observability.Logger
	encryptionManager *EncryptionManager
	processingRecords map[string]*ProcessingRecord
	mu                sync.RWMutex
}

// ProcessingRecord tracks data processing activities
type ProcessingRecord struct {
	RecordID        string                 `json:"record_id"`
	UserID          uuid.UUID              `json:"user_id"`
	DataType        string                 `json:"data_type"`
	Purpose         string                 `json:"purpose"`
	LegalBasis      string                 `json:"legal_basis"`
	ProcessedAt     time.Time              `json:"processed_at"`
	ProcessedBy     string                 `json:"processed_by"`
	RetentionPeriod time.Duration          `json:"retention_period"`
	Metadata        map[string]interface{} `json:"metadata"`
	Status          ProcessingStatus       `json:"status"`
}

// ProcessingStatus represents processing status
type ProcessingStatus string

const (
	ProcessingStatusActive     ProcessingStatus = "active"
	ProcessingStatusCompleted  ProcessingStatus = "completed"
	ProcessingStatusDeleted    ProcessingStatus = "deleted"
	ProcessingStatusAnonymized ProcessingStatus = "anonymized"
)

// RetentionManager manages data retention policies
type RetentionManager struct {
	logger             *observability.Logger
	retentionPolicies  map[string]*RetentionPolicy
	scheduledDeletions map[string]*ScheduledDeletion
	mu                 sync.RWMutex
}

// RetentionPolicy defines data retention rules
type RetentionPolicy struct {
	PolicyID        string          `json:"policy_id"`
	DataType        string          `json:"data_type"`
	Purpose         string          `json:"purpose"`
	RetentionPeriod time.Duration   `json:"retention_period"`
	Action          RetentionAction `json:"action"` // delete, anonymize, archive
	LegalBasis      string          `json:"legal_basis"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	Active          bool            `json:"active"`
}

// RetentionAction defines what to do when retention period expires
type RetentionAction string

const (
	RetentionActionDelete    RetentionAction = "delete"
	RetentionActionAnonymize RetentionAction = "anonymize"
	RetentionActionArchive   RetentionAction = "archive"
)

// ScheduledDeletion represents a scheduled data deletion
type ScheduledDeletion struct {
	DeletionID   string    `json:"deletion_id"`
	UserID       uuid.UUID `json:"user_id"`
	DataType     string    `json:"data_type"`
	ScheduledFor time.Time `json:"scheduled_for"`
	Reason       string    `json:"reason"`
	Status       string    `json:"status"`
}

// DataAnonymizer handles data anonymization
type DataAnonymizer struct {
	logger *observability.Logger
	config *PrivacyConfig
}

// NewPrivacyManager creates a new privacy manager
func NewPrivacyManager(logger *observability.Logger, config *PrivacyConfig, encryptionManager *EncryptionManager) *PrivacyManager {
	pm := &PrivacyManager{
		logger:            logger,
		config:            config,
		encryptionManager: encryptionManager,
	}

	// Initialize components
	pm.consentManager = NewConsentManager(logger)
	pm.dataProcessor = NewDataProcessor(logger, encryptionManager)
	pm.retentionManager = NewRetentionManager(logger)
	pm.anonymizer = NewDataAnonymizer(logger, config)

	return pm
}

// Start starts the privacy manager
func (pm *PrivacyManager) Start(ctx context.Context) error {
	pm.logger.Info(ctx, "Starting privacy manager", map[string]interface{}{
		"gdpr_compliance": pm.config.EnableGDPRCompliance,
		"ccpa_compliance": pm.config.EnableCCPACompliance,
	})

	// Start retention scheduler
	go pm.retentionScheduler(ctx)

	return nil
}

// ProcessPersonalData processes personal data with privacy controls
func (pm *PrivacyManager) ProcessPersonalData(ctx context.Context, userID uuid.UUID, dataType, purpose string, data map[string]interface{}) error {
	// Check consent
	if !pm.consentManager.HasValidConsent(userID, purpose) {
		return fmt.Errorf("no valid consent for purpose: %s", purpose)
	}

	// Check purpose limitation
	if pm.config.EnablePurposeLimitation && !pm.isPurposeAllowed(purpose) {
		return fmt.Errorf("purpose not allowed: %s", purpose)
	}

	// Encrypt PII data
	encryptedData, err := pm.encryptionManager.EncryptPII(data)
	if err != nil {
		return fmt.Errorf("failed to encrypt PII: %w", err)
	}

	// Record processing activity
	record := &ProcessingRecord{
		RecordID:        uuid.New().String(),
		UserID:          userID,
		DataType:        dataType,
		Purpose:         purpose,
		LegalBasis:      "consent",
		ProcessedAt:     time.Now(),
		ProcessedBy:     "system",
		RetentionPeriod: pm.config.DefaultRetentionPeriod,
		Metadata:        encryptedData,
		Status:          ProcessingStatusActive,
	}

	pm.dataProcessor.RecordProcessing(record)

	pm.logger.Info(ctx, "Personal data processed", map[string]interface{}{
		"user_id":   userID,
		"data_type": dataType,
		"purpose":   purpose,
		"record_id": record.RecordID,
	})

	return nil
}

// GrantConsent grants user consent for data processing
func (pm *PrivacyManager) GrantConsent(ctx context.Context, userID uuid.UUID, purposes []string, ipAddress, userAgent string) (*UserConsent, error) {
	consent := &UserConsent{
		UserID:     userID,
		ConsentID:  uuid.New().String(),
		Purposes:   make(map[string]*ConsentPurpose),
		GrantedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(pm.config.ConsentExpirationPeriod),
		Version:    "1.0",
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		LegalBasis: "consent",
		Status:     ConsentStatusActive,
	}

	// Add purposes
	for _, purpose := range purposes {
		consent.Purposes[purpose] = &ConsentPurpose{
			Purpose:     purpose,
			Description: pm.getPurposeDescription(purpose),
			Granted:     true,
			GrantedAt:   time.Now(),
			Required:    pm.isPurposeRequired(purpose),
			Category:    pm.getPurposeCategory(purpose),
		}
	}

	pm.consentManager.StoreConsent(consent)

	pm.logger.Info(ctx, "Consent granted", map[string]interface{}{
		"user_id":    userID,
		"consent_id": consent.ConsentID,
		"purposes":   purposes,
	})

	return consent, nil
}

// WithdrawConsent withdraws user consent
func (pm *PrivacyManager) WithdrawConsent(ctx context.Context, userID uuid.UUID, purpose string) error {
	if err := pm.consentManager.WithdrawConsent(userID, purpose); err != nil {
		return fmt.Errorf("failed to withdraw consent: %w", err)
	}

	// Schedule data deletion if required
	if pm.config.EnableRightToErasure {
		pm.scheduleDataDeletion(ctx, userID, purpose, "consent_withdrawn")
	}

	pm.logger.Info(ctx, "Consent withdrawn", map[string]interface{}{
		"user_id": userID,
		"purpose": purpose,
	})

	return nil
}

// ExportUserData exports user data for portability
func (pm *PrivacyManager) ExportUserData(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	if !pm.config.EnableDataPortability {
		return nil, fmt.Errorf("data portability not enabled")
	}

	// Get all user data
	userData := pm.dataProcessor.GetUserData(userID)

	// Decrypt PII data for export
	decryptedData := make(map[string]interface{})
	for key, value := range userData {
		if encryptedValue, ok := value.(map[string]interface{}); ok {
			if encryptedData, exists := encryptedValue["encrypted_data"]; exists {
				// This is encrypted PII data - decrypt it
				// Implementation would decrypt the data here
				decryptedData[key] = encryptedData
			} else {
				decryptedData[key] = value
			}
		} else {
			decryptedData[key] = value
		}
	}

	pm.logger.Info(ctx, "User data exported", map[string]interface{}{
		"user_id": userID,
	})

	return decryptedData, nil
}

// DeleteUserData deletes all user data (right to erasure)
func (pm *PrivacyManager) DeleteUserData(ctx context.Context, userID uuid.UUID) error {
	if !pm.config.EnableRightToErasure {
		return fmt.Errorf("right to erasure not enabled")
	}

	// Delete all user data
	if err := pm.dataProcessor.DeleteUserData(userID); err != nil {
		return fmt.Errorf("failed to delete user data: %w", err)
	}

	// Delete consent records
	if err := pm.consentManager.DeleteUserConsents(userID); err != nil {
		return fmt.Errorf("failed to delete consent records: %w", err)
	}

	pm.logger.Info(ctx, "User data deleted", map[string]interface{}{
		"user_id": userID,
	})

	return nil
}

// retentionScheduler runs data retention policies
func (pm *PrivacyManager) retentionScheduler(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour) // Run daily
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			pm.logger.Info(ctx, "Running data retention policies", nil)
			pm.retentionManager.ProcessRetentionPolicies(ctx)
		}
	}
}

// scheduleDataDeletion schedules data for deletion
func (pm *PrivacyManager) scheduleDataDeletion(ctx context.Context, userID uuid.UUID, dataType, reason string) {
	deletion := &ScheduledDeletion{
		DeletionID:   uuid.New().String(),
		UserID:       userID,
		DataType:     dataType,
		ScheduledFor: time.Now().Add(30 * 24 * time.Hour), // 30 days grace period
		Reason:       reason,
		Status:       "scheduled",
	}

	pm.retentionManager.ScheduleDeletion(deletion)
}

// Helper functions
func (pm *PrivacyManager) isPurposeAllowed(purpose string) bool {
	allowedPurposes := map[string]bool{
		"essential":  true,
		"functional": true,
		"analytics":  true,
		"marketing":  true,
		"trading":    true,
		"compliance": true,
	}
	return allowedPurposes[purpose]
}

func (pm *PrivacyManager) getPurposeDescription(purpose string) string {
	descriptions := map[string]string{
		"essential":  "Essential for platform operation",
		"functional": "Enhance user experience",
		"analytics":  "Analyze usage patterns",
		"marketing":  "Personalized marketing",
		"trading":    "Execute trading operations",
		"compliance": "Regulatory compliance",
	}
	return descriptions[purpose]
}

func (pm *PrivacyManager) isPurposeRequired(purpose string) bool {
	requiredPurposes := map[string]bool{
		"essential":  true,
		"compliance": true,
	}
	return requiredPurposes[purpose]
}

func (pm *PrivacyManager) getPurposeCategory(purpose string) string {
	categories := map[string]string{
		"essential":  "essential",
		"functional": "functional",
		"analytics":  "analytics",
		"marketing":  "marketing",
		"trading":    "functional",
		"compliance": "essential",
	}
	return categories[purpose]
}

// NewConsentManager creates a new consent manager
func NewConsentManager(logger *observability.Logger) *ConsentManager {
	return &ConsentManager{
		logger:   logger,
		consents: make(map[string]*UserConsent),
	}
}

// StoreConsent stores user consent
func (cm *ConsentManager) StoreConsent(consent *UserConsent) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	key := fmt.Sprintf("%s_%s", consent.UserID.String(), consent.ConsentID)
	cm.consents[key] = consent
}

// HasValidConsent checks if user has valid consent for purpose
func (cm *ConsentManager) HasValidConsent(userID uuid.UUID, purpose string) bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	for _, consent := range cm.consents {
		if consent.UserID == userID && consent.Status == ConsentStatusActive {
			if purposeConsent, exists := consent.Purposes[purpose]; exists {
				return purposeConsent.Granted && time.Now().Before(consent.ExpiresAt)
			}
		}
	}
	return false
}

// WithdrawConsent withdraws consent for a specific purpose
func (cm *ConsentManager) WithdrawConsent(userID uuid.UUID, purpose string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for _, consent := range cm.consents {
		if consent.UserID == userID && consent.Status == ConsentStatusActive {
			if purposeConsent, exists := consent.Purposes[purpose]; exists {
				purposeConsent.Granted = false
				now := time.Now()
				consent.WithdrawnAt = &now
				return nil
			}
		}
	}
	return fmt.Errorf("consent not found for user %s and purpose %s", userID, purpose)
}

// DeleteUserConsents deletes all consents for a user
func (cm *ConsentManager) DeleteUserConsents(userID uuid.UUID) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for key, consent := range cm.consents {
		if consent.UserID == userID {
			delete(cm.consents, key)
		}
	}
	return nil
}

// NewDataProcessor creates a new data processor
func NewDataProcessor(logger *observability.Logger, encryptionManager *EncryptionManager) *DataProcessor {
	return &DataProcessor{
		logger:            logger,
		encryptionManager: encryptionManager,
		processingRecords: make(map[string]*ProcessingRecord),
	}
}

// RecordProcessing records a data processing activity
func (dp *DataProcessor) RecordProcessing(record *ProcessingRecord) {
	dp.mu.Lock()
	defer dp.mu.Unlock()

	dp.processingRecords[record.RecordID] = record
}

// GetUserData gets all data for a user
func (dp *DataProcessor) GetUserData(userID uuid.UUID) map[string]interface{} {
	dp.mu.RLock()
	defer dp.mu.RUnlock()

	userData := make(map[string]interface{})
	for _, record := range dp.processingRecords {
		if record.UserID == userID && record.Status == ProcessingStatusActive {
			userData[record.DataType] = record.Metadata
		}
	}
	return userData
}

// DeleteUserData deletes all data for a user
func (dp *DataProcessor) DeleteUserData(userID uuid.UUID) error {
	dp.mu.Lock()
	defer dp.mu.Unlock()

	for _, record := range dp.processingRecords {
		if record.UserID == userID {
			record.Status = ProcessingStatusDeleted
		}
	}
	return nil
}

// NewRetentionManager creates a new retention manager
func NewRetentionManager(logger *observability.Logger) *RetentionManager {
	return &RetentionManager{
		logger:             logger,
		retentionPolicies:  make(map[string]*RetentionPolicy),
		scheduledDeletions: make(map[string]*ScheduledDeletion),
	}
}

// ProcessRetentionPolicies processes retention policies
func (rm *RetentionManager) ProcessRetentionPolicies(ctx context.Context) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	for _, deletion := range rm.scheduledDeletions {
		if time.Now().After(deletion.ScheduledFor) && deletion.Status == "scheduled" {
			rm.logger.Info(ctx, "Processing scheduled deletion", map[string]interface{}{
				"deletion_id": deletion.DeletionID,
				"user_id":     deletion.UserID,
				"data_type":   deletion.DataType,
			})
			deletion.Status = "completed"
		}
	}
}

// ScheduleDeletion schedules a data deletion
func (rm *RetentionManager) ScheduleDeletion(deletion *ScheduledDeletion) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rm.scheduledDeletions[deletion.DeletionID] = deletion
}

// NewDataAnonymizer creates a new data anonymizer
func NewDataAnonymizer(logger *observability.Logger, config *PrivacyConfig) *DataAnonymizer {
	return &DataAnonymizer{
		logger: logger,
		config: config,
	}
}
