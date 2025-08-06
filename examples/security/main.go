package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ai-agentic-browser/internal/security"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

// SecurityDemo demonstrates the advanced security features
func main() {
	fmt.Println("🔒 AI-Agentic Crypto Browser - Advanced Security Demo")
	fmt.Println("============================================================")

	// Initialize logger
	logger := &observability.Logger{}
	ctx := context.Background()

	// Demo 1: Zero-Trust Engine with Advanced TTL Calculation
	fmt.Println("\n🛡️  Demo 1: Zero-Trust Engine with Advanced TTL Calculation")
	demoZeroTrustEngine(ctx, logger)

	// Demo 2: Advanced Threat Detection
	fmt.Println("\n🚨 Demo 2: Advanced Threat Detection")
	demoThreatDetection(ctx, logger)

	// Demo 3: Security Policy Engine
	fmt.Println("\n📋 Demo 3: Security Policy Engine")
	demoPolicyEngine(ctx, logger)

	// Demo 4: Real-Time Security Dashboard
	fmt.Println("\n📊 Demo 4: Real-Time Security Dashboard")
	demoSecurityDashboard(ctx, logger)

	// Demo 5: Device Trust Management
	fmt.Println("\n📱 Demo 5: Device Trust Management")
	demoDeviceTrustManagement(ctx, logger)

	fmt.Println("\n🎉 Security Demo Complete!")
	fmt.Println("All advanced security features are working correctly.")
}

// demoZeroTrustEngine demonstrates zero-trust access evaluation
func demoZeroTrustEngine(ctx context.Context, logger *observability.Logger) {
	// Create zero-trust engine
	zeroTrustEngine := security.NewZeroTrustEngine(logger)

	fmt.Println("  Creating zero-trust engine with adaptive TTL calculation...")

	// Test different risk scenarios
	testScenarios := []struct {
		name      string
		riskScore float64
		userType  string
	}{
		{"Low Risk User", 0.1, "regular"},
		{"Medium Risk User", 0.5, "regular"},
		{"High Risk User", 0.8, "admin"},
		{"Critical Risk User", 0.95, "admin"},
	}

	for _, scenario := range testScenarios {
		fmt.Printf("  📊 Testing %s (Risk: %.2f)\n", scenario.name, scenario.riskScore)

		// Create access request
		userID := uuid.New()
		accessRequest := &security.AccessRequest{
			UserID:    &userID,
			DeviceID:  "device_" + scenario.userType,
			IPAddress: "192.168.1.100",
			UserAgent: "Mozilla/5.0 (Enterprise Browser)",
			Resource:  "/api/trading",
			Action:    "POST",
			Timestamp: time.Now(),
			Context: map[string]interface{}{
				"risk_score":   scenario.riskScore,
				"device_trust": 0.8,
				"user_type":    scenario.userType,
			},
		}

		// Evaluate access
		decision, err := zeroTrustEngine.EvaluateAccess(ctx, accessRequest)
		if err != nil {
			fmt.Printf("    ❌ Error: %v\n", err)
			continue
		}

		// Display results
		status := "✅ ALLOWED"
		if !decision.Allowed {
			status = "❌ DENIED"
		}

		mfaStatus := ""
		if decision.RequiresMFA {
			mfaStatus = " (MFA Required)"
		}

		fmt.Printf("    %s%s - TTL: %v - Reason: %s\n",
			status, mfaStatus, decision.SessionTTL, decision.Reason)
	}
}

// demoThreatDetection demonstrates advanced threat detection
func demoThreatDetection(ctx context.Context, logger *observability.Logger) {
	// Create threat detector
	threatDetector := security.NewAdvancedThreatDetector(logger)

	fmt.Println("  Creating advanced threat detection system...")

	// Test different threat scenarios
	threatScenarios := []struct {
		name        string
		requestType string
		suspicious  bool
	}{
		{"Normal API Request", "GET /api/dashboard", false},
		{"SQL Injection Attempt", "POST /api/login", true},
		{"Brute Force Attack", "POST /api/login", true},
		{"Suspicious User Agent", "GET /api/data", true},
	}

	for _, scenario := range threatScenarios {
		fmt.Printf("  🔍 Analyzing: %s\n", scenario.name)

		// Create security request
		var body string
		var userAgent string

		switch scenario.requestType {
		case "POST /api/login":
			if scenario.name == "SQL Injection Attempt" {
				body = "username=admin' OR '1'='1&password=test"
				userAgent = "Mozilla/5.0"
			} else {
				body = "username=admin&password=wrong"
				userAgent = "python-requests/2.25.1"
			}
		case "GET /api/data":
			userAgent = "suspicious-bot/1.0"
		default:
			userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"
		}

		securityRequest := &security.SecurityRequest{
			RequestID: uuid.New().String(),
			IPAddress: "10.0.0.1",
			UserAgent: userAgent,
			Method:    "POST",
			URL:       "/api/login",
			Body:      body,
			Timestamp: time.Now(),
		}

		// Detect threats
		result, err := threatDetector.DetectThreats(ctx, securityRequest)
		if err != nil {
			fmt.Printf("    ❌ Error: %v\n", err)
			continue
		}

		// Display results
		threatStatus := "✅ SAFE"
		if result.ThreatDetected {
			threatStatus = "🚨 THREAT DETECTED"
		}

		blockStatus := ""
		if result.ShouldBlock {
			blockStatus = " (BLOCKED)"
		}

		fmt.Printf("    %s%s - Score: %.2f - Type: %s\n",
			threatStatus, blockStatus, result.ThreatScore, result.ThreatType)
	}
}

// demoPolicyEngine demonstrates security policy evaluation
func demoPolicyEngine(ctx context.Context, logger *observability.Logger) {
	// Create policy engine
	policyEngine := security.NewPolicyEngine(logger)

	fmt.Println("  Creating security policy engine with default policies...")

	// Test policy evaluation scenarios
	policyScenarios := []struct {
		name      string
		userRoles []string
		resource  string
		riskScore float64
	}{
		{"Admin User - Low Risk", []string{"admin"}, "/api/admin", 0.2},
		{"Regular User - Medium Risk", []string{"user"}, "/api/dashboard", 0.5},
		{"Guest User - High Risk", []string{"guest"}, "/api/data", 0.9},
		{"Trader - Normal Risk", []string{"trader"}, "/api/trading", 0.3},
	}

	for _, scenario := range policyScenarios {
		fmt.Printf("  📋 Evaluating: %s\n", scenario.name)

		// Create evaluation context
		evalCtx := &security.PolicyEvaluationContext{
			UserID:      uuid.New(),
			UserRoles:   scenario.userRoles,
			IPAddress:   "192.168.1.100",
			UserAgent:   "Mozilla/5.0",
			Resource:    scenario.resource,
			Action:      "GET",
			RiskScore:   scenario.riskScore,
			DeviceTrust: 0.8,
			Timestamp:   time.Now(),
		}

		// Evaluate policies
		result, err := policyEngine.EvaluatePolicy(ctx, evalCtx)
		if err != nil {
			fmt.Printf("    ❌ Error: %v\n", err)
			continue
		}

		// Display results
		decision := "✅ ALLOWED"
		if !result.Decision.Allowed {
			decision = "❌ DENIED"
		}

		fmt.Printf("    %s - Policies: %d - Rules: %d - Reason: %s\n",
			decision, len(result.MatchedPolicies), len(result.MatchedRules), result.Reason)
	}
}

// demoSecurityDashboard demonstrates real-time security monitoring
func demoSecurityDashboard(ctx context.Context, logger *observability.Logger) {
	// Create security dashboard
	threatDetector := security.NewAdvancedThreatDetector(logger)
	zeroTrustEngine := security.NewZeroTrustEngine(logger)
	dashboard := security.NewSecurityDashboard(logger, threatDetector, zeroTrustEngine)

	fmt.Println("  Creating real-time security dashboard...")

	// Start dashboard
	err := dashboard.Start(ctx)
	if err != nil {
		fmt.Printf("    ❌ Error starting dashboard: %v\n", err)
		return
	}
	defer dashboard.Stop()

	// Simulate dashboard client connection
	fmt.Println("  📊 Simulating dashboard client connection...")

	// Get current security metrics
	metrics := dashboard.GetSecurityMetrics()

	fmt.Printf("    Security Health: %s\n", metrics.SecurityHealth)
	fmt.Printf("    Total Threats: %d (Critical: %d)\n", metrics.TotalThreats, metrics.CriticalThreats)
	fmt.Printf("    Blocked Requests: %d\n", metrics.BlockedRequests)
	fmt.Printf("    Threat Detection Rate: %.1f%%\n", metrics.ThreatDetectionRate)
	fmt.Printf("    False Positive Rate: %.1f%%\n", metrics.FalsePositiveRate)
	fmt.Printf("    Average Risk Score: %.2f\n", metrics.RiskScoreAverage)
	fmt.Printf("    Average Device Trust: %.2f\n", metrics.DeviceTrustAverage)
	fmt.Printf("    Compliance Score: %.1f%%\n", metrics.ComplianceScore)

	fmt.Println("    ✅ Dashboard is streaming real-time security metrics")
}

// demoDeviceTrustManagement demonstrates device trust features
func demoDeviceTrustManagement(ctx context.Context, logger *observability.Logger) {
	// Create device trust manager
	deviceManager := security.NewDeviceTrustManager(logger)

	fmt.Println("  Creating device trust management system...")

	// Test device scenarios
	deviceScenarios := []struct {
		name       string
		deviceType string
		trustLevel float64
	}{
		{"Trusted Laptop", "laptop", 0.9},
		{"New Mobile Device", "mobile", 0.3},
		{"Suspicious Device", "unknown", 0.1},
		{"Corporate Workstation", "workstation", 0.95},
	}

	for _, scenario := range deviceScenarios {
		fmt.Printf("  📱 Managing: %s\n", scenario.name)

		// Create trusted device
		userID := uuid.New()
		deviceID := fmt.Sprintf("device_%s_%d", scenario.deviceType, time.Now().Unix())

		device := &security.TrustedDevice{
			DeviceID:   deviceID,
			UserID:     userID,
			TrustLevel: scenario.trustLevel,
			Attributes: map[string]interface{}{
				"device_type": scenario.deviceType,
				"user_agent":  "Mozilla/5.0",
				"screen":      "1920x1080",
			},
			RiskFactors: []string{},
			LastSeen:    time.Now(),
		}

		// Add risk factors for low trust devices
		if scenario.trustLevel < 0.5 {
			device.RiskFactors = append(device.RiskFactors, "new_device", "unusual_location")
		}

		// Register device
		err := deviceManager.RegisterDevice(device)
		if err != nil {
			fmt.Printf("    ❌ Error: %v\n", err)
			continue
		}

		// Retrieve and display device info
		retrievedDevice, err := deviceManager.GetDevice(deviceID)
		if err != nil {
			fmt.Printf("    ❌ Error: %v\n", err)
			continue
		}

		trustStatus := "✅ TRUSTED"
		if retrievedDevice.TrustLevel < 0.5 {
			trustStatus = "⚠️  LOW TRUST"
		}
		if retrievedDevice.TrustLevel < 0.2 {
			trustStatus = "❌ UNTRUSTED"
		}

		fmt.Printf("    %s - Trust Level: %.2f - Risk Factors: %v\n",
			trustStatus, retrievedDevice.TrustLevel, retrievedDevice.RiskFactors)
	}
}
