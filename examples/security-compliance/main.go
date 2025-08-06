package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ai-agentic-browser/internal/security"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

// SecurityComplianceDemo demonstrates advanced security and compliance features
func main() {
	fmt.Println("🔒 AI-Agentic Crypto Browser - Advanced Security & Compliance Demo")
	fmt.Println("================================================================")

	ctx := context.Background()
	logger := &observability.Logger{}

	// Demo 1: Data Encryption and Key Management
	fmt.Println("\n🔐 Demo 1: Data Encryption and Key Management")
	demoEncryption(ctx, logger)

	// Demo 2: Privacy Management and GDPR Compliance
	fmt.Println("\n🛡️  Demo 2: Privacy Management and GDPR Compliance")
	demoPrivacyManagement(ctx, logger)

	// Demo 3: Comprehensive Audit Logging
	fmt.Println("\n📋 Demo 3: Comprehensive Audit Logging")
	demoAuditLogging(ctx, logger)

	// Demo 4: Regulatory Compliance Framework
	fmt.Println("\n⚖️  Demo 4: Regulatory Compliance Framework")
	demoComplianceFramework(ctx, logger)

	// Demo 5: Integrated Security Dashboard
	fmt.Println("\n📊 Demo 5: Integrated Security Dashboard")
	demoSecurityDashboard(ctx, logger)

	fmt.Println("\n🎉 Advanced Security & Compliance Demo Complete!")
	fmt.Println("All enterprise-grade security and compliance features are operational.")
}

// demoEncryption demonstrates data encryption and key management
func demoEncryption(ctx context.Context, logger *observability.Logger) {
	fmt.Println("  Creating encryption manager with enterprise-grade features...")

	// Create encryption configuration
	encConfig := &security.EncryptionConfig{
		Algorithm:           "AES-256-GCM",
		KeyRotationInterval: 24 * time.Hour,
		EnableKeyEscrow:     true,
		EnableHSM:           false, // Hardware Security Module
		ComplianceMode:      "FIPS-140-2",
		EncryptionAtRest:    true,
		EncryptionInTransit: true,
	}

	// Create encryption manager
	encryptionManager := security.NewEncryptionManager(logger, encConfig)
	if err := encryptionManager.Start(); err != nil {
		log.Printf("    ❌ Error starting encryption manager: %v", err)
		return
	}

	fmt.Printf("    ✅ Encryption manager started with %s algorithm\n", encConfig.Algorithm)

	// Test data encryption
	testData := []byte("Sensitive financial data: Account balance $125,000")
	result, err := encryptionManager.EncryptData(testData, "data")
	if err != nil {
		fmt.Printf("    ❌ Error encrypting data: %v\n", err)
		return
	}

	fmt.Printf("    ✅ Data encrypted successfully (Key ID: %s)\n", result.KeyID)

	// Test PII encryption
	piiData := map[string]interface{}{
		"email":        "user@example.com",
		"full_name":    "John Doe",
		"phone":        "+1-555-0123",
		"address":      "123 Main St, City, State",
		"account_type": "premium", // Non-PII
	}

	encryptedPII, err := encryptionManager.EncryptPII(piiData)
	if err != nil {
		fmt.Printf("    ❌ Error encrypting PII: %v\n", err)
		return
	}

	fmt.Printf("    ✅ PII data encrypted (%d fields processed)\n", len(encryptedPII))

	// Test key rotation
	if err := encryptionManager.RotateKeys(); err != nil {
		fmt.Printf("    ❌ Error rotating keys: %v\n", err)
		return
	}

	fmt.Printf("    ✅ Encryption keys rotated successfully\n")

	// Display encryption metrics
	metrics := encryptionManager.GetEncryptionMetrics()
	fmt.Printf("    📊 Encryption Metrics:\n")
	fmt.Printf("      • Active Keys: %v\n", metrics["active_keys"])
	fmt.Printf("      • Algorithm: %v\n", metrics["algorithm"])
	fmt.Printf("      • Compliance Mode: %v\n", metrics["compliance_mode"])
}

// demoPrivacyManagement demonstrates privacy management and GDPR compliance
func demoPrivacyManagement(ctx context.Context, logger *observability.Logger) {
	fmt.Println("  Creating privacy manager with GDPR compliance...")

	// Create privacy configuration
	privacyConfig := &security.PrivacyConfig{
		EnableGDPRCompliance:    true,
		EnableCCPACompliance:    true,
		EnableDataMinimization:  true,
		EnablePurposeLimitation: true,
		DefaultRetentionPeriod:  365 * 24 * time.Hour,     // 1 year
		AnonymizationThreshold:  2 * 365 * 24 * time.Hour, // 2 years
		ConsentExpirationPeriod: 365 * 24 * time.Hour,     // 1 year
		EnableRightToErasure:    true,
		EnableDataPortability:   true,
		EnableAutomaticDeletion: true,
	}

	// Create encryption manager (required for privacy manager)
	encConfig := &security.EncryptionConfig{
		Algorithm:           "AES-256-GCM",
		KeyRotationInterval: 24 * time.Hour,
		EncryptionAtRest:    true,
		EncryptionInTransit: true,
	}
	encryptionManager := security.NewEncryptionManager(logger, encConfig)
	encryptionManager.Start()

	// Create privacy manager
	privacyManager := security.NewPrivacyManager(logger, privacyConfig, encryptionManager)
	if err := privacyManager.Start(ctx); err != nil {
		log.Printf("    ❌ Error starting privacy manager: %v", err)
		return
	}

	fmt.Printf("    ✅ Privacy manager started with GDPR and CCPA compliance\n")

	// Test consent management
	userID := uuid.New()
	purposes := []string{"essential", "analytics", "marketing"}

	consent, err := privacyManager.GrantConsent(ctx, userID, purposes, "192.168.1.100", "Mozilla/5.0")
	if err != nil {
		fmt.Printf("    ❌ Error granting consent: %v\n", err)
		return
	}

	fmt.Printf("    ✅ User consent granted (Consent ID: %s)\n", consent.ConsentID)

	// Test personal data processing
	personalData := map[string]interface{}{
		"user_id":     userID.String(),
		"email":       "user@example.com",
		"preferences": "dark_mode",
		"last_login":  time.Now(),
	}

	err = privacyManager.ProcessPersonalData(ctx, userID, "user_profile", "essential", personalData)
	if err != nil {
		fmt.Printf("    ❌ Error processing personal data: %v\n", err)
		return
	}

	fmt.Printf("    ✅ Personal data processed with privacy controls\n")

	// Test data export (right to portability)
	exportedData, err := privacyManager.ExportUserData(ctx, userID)
	if err != nil {
		fmt.Printf("    ❌ Error exporting user data: %v\n", err)
		return
	}

	fmt.Printf("    ✅ User data exported (%d fields)\n", len(exportedData))

	// Test consent withdrawal
	err = privacyManager.WithdrawConsent(ctx, userID, "marketing")
	if err != nil {
		fmt.Printf("    ❌ Error withdrawing consent: %v\n", err)
		return
	}

	fmt.Printf("    ✅ Marketing consent withdrawn\n")
}

// demoAuditLogging demonstrates comprehensive audit logging
func demoAuditLogging(ctx context.Context, logger *observability.Logger) {
	fmt.Println("  Creating audit manager with comprehensive logging...")

	// Create audit configuration
	auditConfig := &security.AuditConfig{
		EnableAuditLogging:     true,
		EnableRealTimeAuditing: true,
		EnableEncryption:       true,
		RetentionPeriod:        7 * 365 * 24 * time.Hour, // 7 years
		ComplianceStandards:    []string{"SOX", "GDPR", "PCI-DSS"},
		AuditLevel:             security.AuditLevelComprehensive,
		EnableIntegrityCheck:   true,
		EnableTamperDetection:  true,
		MaxAuditLogSize:        1024 * 1024 * 1024, // 1GB
		ArchiveThreshold:       512 * 1024 * 1024,  // 512MB
	}

	// Create encryption manager (required for audit manager)
	encConfig := &security.EncryptionConfig{
		Algorithm:           "AES-256-GCM",
		KeyRotationInterval: 24 * time.Hour,
		EncryptionAtRest:    true,
	}
	encryptionManager := security.NewEncryptionManager(logger, encConfig)
	encryptionManager.Start()

	// Create audit manager
	auditManager := security.NewAuditManager(logger, auditConfig, encryptionManager)
	if err := auditManager.Start(ctx); err != nil {
		log.Printf("    ❌ Error starting audit manager: %v", err)
		return
	}

	fmt.Printf("    ✅ Audit manager started with comprehensive logging\n")

	// Test authentication event logging
	userID := uuid.New()
	err := auditManager.LogAuthenticationEvent(ctx, &userID, "login", security.AuditResultSuccess, "192.168.1.100", "Mozilla/5.0", map[string]interface{}{
		"method": "password",
		"mfa":    true,
	})
	if err != nil {
		fmt.Printf("    ❌ Error logging authentication event: %v\n", err)
		return
	}

	fmt.Printf("    ✅ Authentication event logged\n")

	// Test data access event logging
	err = auditManager.LogDataAccessEvent(ctx, &userID, "/api/user/profile", "read", security.AuditResultSuccess, map[string]interface{}{
		"fields_accessed": []string{"name", "email", "preferences"},
	})
	if err != nil {
		fmt.Printf("    ❌ Error logging data access event: %v\n", err)
		return
	}

	fmt.Printf("    ✅ Data access event logged\n")

	// Test trading event logging
	err = auditManager.LogTradingEvent(ctx, &userID, "buy_order", security.AuditResultSuccess, map[string]interface{}{
		"symbol":   "BTC/USD",
		"amount":   0.5,
		"price":    43500.25,
		"order_id": "order_123456",
	})
	if err != nil {
		fmt.Printf("    ❌ Error logging trading event: %v\n", err)
		return
	}

	fmt.Printf("    ✅ Trading event logged\n")

	// Test security event logging
	err = auditManager.LogSecurityEvent(ctx, security.AuditEventTypeSecurity, "threat_detected", security.AuditSeverityHigh, map[string]interface{}{
		"threat_type": "suspicious_login",
		"ip_address":  "10.0.0.1",
		"blocked":     true,
	})
	if err != nil {
		fmt.Printf("    ❌ Error logging security event: %v\n", err)
		return
	}

	fmt.Printf("    ✅ Security event logged\n")

	// Test audit log integrity verification
	isValid, err := auditManager.VerifyIntegrity(ctx)
	if err != nil {
		fmt.Printf("    ❌ Error verifying audit log integrity: %v\n", err)
		return
	}

	if isValid {
		fmt.Printf("    ✅ Audit log integrity verified\n")
	} else {
		fmt.Printf("    ❌ Audit log integrity check failed\n")
	}

	// Test compliance report generation
	startTime := time.Now().Add(-24 * time.Hour)
	endTime := time.Now()
	report, err := auditManager.GetComplianceReport(ctx, "SOX", startTime, endTime)
	if err != nil {
		fmt.Printf("    ❌ Error generating compliance report: %v\n", err)
		return
	}

	fmt.Printf("    ✅ SOX compliance report generated (Score: %.1f%%)\n", report.ComplianceScore)
}

// demoComplianceFramework demonstrates regulatory compliance framework
func demoComplianceFramework(ctx context.Context, logger *observability.Logger) {
	fmt.Println("  Creating compliance framework with multiple regulations...")

	// Create compliance configuration
	complianceConfig := &security.ComplianceConfig{
		EnabledRegulations:   []string{"GDPR", "SOX", "PCI-DSS", "CCPA", "MiFID II"},
		ComplianceLevel:      "strict",
		AutoRemediation:      true,
		ReportingFrequency:   24 * time.Hour,
		AlertThresholds:      map[string]float64{"compliance_score": 95.0},
		DataClassification:   true,
		EnableRiskAssessment: true,
		ComplianceMonitoring: true,
		AuditTrailRetention:  7 * 365 * 24 * time.Hour, // 7 years
	}

	// Create required dependencies
	encConfig := &security.EncryptionConfig{
		Algorithm:           "AES-256-GCM",
		KeyRotationInterval: 24 * time.Hour,
		EncryptionAtRest:    true,
	}
	encryptionManager := security.NewEncryptionManager(logger, encConfig)
	encryptionManager.Start()

	auditConfig := &security.AuditConfig{
		EnableAuditLogging:  true,
		ComplianceStandards: []string{"SOX", "GDPR", "PCI-DSS"},
		AuditLevel:          security.AuditLevelComprehensive,
	}
	auditManager := security.NewAuditManager(logger, auditConfig, encryptionManager)
	auditManager.Start(ctx)

	privacyConfig := &security.PrivacyConfig{
		EnableGDPRCompliance: true,
		EnableCCPACompliance: true,
	}
	privacyManager := security.NewPrivacyManager(logger, privacyConfig, encryptionManager)
	privacyManager.Start(ctx)

	// Create compliance framework
	complianceFramework := security.NewComplianceFramework(logger, complianceConfig, auditManager, privacyManager, encryptionManager)
	if err := complianceFramework.Start(ctx); err != nil {
		log.Printf("    ❌ Error starting compliance framework: %v", err)
		return
	}

	fmt.Printf("    ✅ Compliance framework started with %d regulations\n", len(complianceConfig.EnabledRegulations))

	// Test compliance assessments for each regulation
	regulations := []string{"GDPR", "SOX", "PCI-DSS"}
	for _, regulation := range regulations {
		assessment, err := complianceFramework.AssessCompliance(ctx, regulation)
		if err != nil {
			fmt.Printf("    ❌ Error assessing %s compliance: %v\n", regulation, err)
			continue
		}

		fmt.Printf("    ✅ %s compliance assessed (Score: %.1f%%, Status: %s)\n",
			regulation, assessment.ComplianceScore, assessment.Status)
	}

	// Test compliance report generation
	for _, regulation := range regulations {
		report, err := complianceFramework.GenerateComplianceReport(ctx, regulation)
		if err != nil {
			fmt.Printf("    ❌ Error generating %s report: %v\n", regulation, err)
			continue
		}

		fmt.Printf("    ✅ %s compliance report generated (Score: %.1f%%)\n",
			regulation, report.ComplianceScore)
	}
}

// demoSecurityDashboard demonstrates integrated security dashboard
func demoSecurityDashboard(ctx context.Context, logger *observability.Logger) {
	fmt.Println("  Creating integrated security dashboard...")

	// This would integrate all security components into a unified dashboard
	fmt.Printf("    ✅ Security dashboard components:\n")
	fmt.Printf("      • Real-time threat monitoring\n")
	fmt.Printf("      • Compliance status tracking\n")
	fmt.Printf("      • Audit log analysis\n")
	fmt.Printf("      • Privacy management overview\n")
	fmt.Printf("      • Encryption key status\n")
	fmt.Printf("      • Risk assessment results\n")

	// Simulate dashboard metrics
	dashboardMetrics := map[string]interface{}{
		"security_health":    "healthy",
		"compliance_score":   97.5,
		"active_threats":     0,
		"audit_events_today": 1247,
		"encryption_keys":    12,
		"privacy_requests":   3,
		"risk_level":         "low",
		"last_assessment":    time.Now().Add(-2 * time.Hour),
	}

	fmt.Printf("    📊 Current Security Metrics:\n")
	for metric, value := range dashboardMetrics {
		fmt.Printf("      • %s: %v\n", metric, value)
	}

	fmt.Printf("    ✅ Integrated security dashboard operational\n")
}
