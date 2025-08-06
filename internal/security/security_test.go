package security

import (
	"context"
	"testing"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestZeroTrustEngine_EvaluateAccess(t *testing.T) {
	logger := &observability.Logger{}
	engine := NewZeroTrustEngine(logger)

	tests := []struct {
		name           string
		request        *AccessRequest
		expectedResult bool
		expectedReason string
	}{
		{
			name: "Low risk access allowed",
			request: &AccessRequest{
				UserID:    &[]uuid.UUID{uuid.New()}[0],
				IPAddress: "192.168.1.100",
				Resource:  "/api/dashboard",
				Action:    "GET",
				Timestamp: time.Now(),
			},
			expectedResult: true,
		},
		{
			name: "High risk access denied",
			request: &AccessRequest{
				UserID:    &[]uuid.UUID{uuid.New()}[0],
				IPAddress: "10.0.0.1", // Suspicious IP
				Resource:  "/api/admin",
				Action:    "DELETE",
				Timestamp: time.Now(),
			},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decision, err := engine.EvaluateAccess(context.Background(), tt.request)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedResult, decision.Allowed)
		})
	}
}

func TestZeroTrustEngine_CalculateSessionTTL(t *testing.T) {
	logger := &observability.Logger{}
	engine := NewZeroTrustEngine(logger)

	tests := []struct {
		name      string
		riskScore float64
		expected  time.Duration
	}{
		{
			name:      "Zero risk - full TTL",
			riskScore: 0.0,
			expected:  30 * time.Minute, // baseTTL
		},
		{
			name:      "Medium risk - reduced TTL",
			riskScore: 0.5,
			expected:  15 * time.Minute, // 50% of baseTTL
		},
		{
			name:      "High risk - minimum TTL",
			riskScore: 0.9,
			expected:  5 * time.Minute, // minimum TTL
		},
		{
			name:      "Maximum risk - minimum TTL",
			riskScore: 1.0,
			expected:  5 * time.Minute, // minimum TTL
		},
		{
			name:      "Invalid high risk - minimum TTL",
			riskScore: 1.5,
			expected:  5 * time.Minute, // clamped to minimum
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := engine.calculateSessionTTL(tt.riskScore)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDeviceTrustManager_RegisterDevice(t *testing.T) {
	logger := &observability.Logger{}
	manager := NewDeviceTrustManager(logger)

	userID := uuid.New()
	deviceID := "test-device-123"

	device := &TrustedDevice{
		DeviceID:   deviceID,
		UserID:     userID,
		TrustLevel: 0.5,
		Attributes: map[string]interface{}{
			"user_agent": "Mozilla/5.0",
			"screen":     "1920x1080",
		},
		RiskFactors: []string{"new_device"},
		LastSeen:    time.Now(),
	}

	err := manager.RegisterDevice(device)
	require.NoError(t, err)

	// Verify device was registered
	retrievedDevice, err := manager.GetDevice(deviceID)
	require.NoError(t, err)
	assert.Equal(t, deviceID, retrievedDevice.DeviceID)
	assert.Equal(t, userID, retrievedDevice.UserID)
}

func TestDeviceTrustManager_UpdateTrustLevel(t *testing.T) {
	logger := &observability.Logger{}
	manager := NewDeviceTrustManager(logger)

	userID := uuid.New()
	deviceID := "test-device-456"

	// Register device
	device := &TrustedDevice{
		DeviceID:   deviceID,
		UserID:     userID,
		TrustLevel: 0.3,
		LastSeen:   time.Now(),
	}

	err := manager.RegisterDevice(device)
	require.NoError(t, err)

	// Update trust level
	newTrustLevel := 0.8
	err = manager.UpdateTrustLevel(deviceID, newTrustLevel)
	require.NoError(t, err)

	// Verify trust level was updated
	updatedDevice, err := manager.GetDevice(deviceID)
	require.NoError(t, err)
	assert.Equal(t, newTrustLevel, updatedDevice.TrustLevel)
}

func TestAdvancedThreatDetector_DetectThreats(t *testing.T) {
	logger := &observability.Logger{}
	detector := NewAdvancedThreatDetector(logger)

	tests := []struct {
		name           string
		request        *SecurityRequest
		expectedThreat bool
		expectedScore  float64
		expectedType   ThreatType
	}{
		{
			name: "Normal request - no threat",
			request: &SecurityRequest{
				RequestID: "req-123",
				IPAddress: "192.168.1.100",
				UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
				Method:    "GET",
				URL:       "/api/dashboard",
				Timestamp: time.Now(),
			},
			expectedThreat: false,
			expectedScore:  0.1,
		},
		{
			name: "SQL injection attempt",
			request: &SecurityRequest{
				RequestID: "req-456",
				IPAddress: "10.0.0.1",
				UserAgent: "curl/7.68.0",
				Method:    "POST",
				URL:       "/api/login",
				Body:      "username=admin' OR '1'='1&password=test",
				Timestamp: time.Now(),
			},
			expectedThreat: true,
			expectedScore:  0.9,
			expectedType:   ThreatTypeSuspiciousAPI,
		},
		{
			name: "Brute force attempt",
			request: &SecurityRequest{
				RequestID: "req-789",
				IPAddress: "10.0.0.2",
				UserAgent: "python-requests/2.25.1",
				Method:    "POST",
				URL:       "/api/login",
				Timestamp: time.Now(),
			},
			expectedThreat: true,
			expectedScore:  0.7,
			expectedType:   ThreatTypeBruteForce,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := detector.DetectThreats(context.Background(), tt.request)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedThreat, result.ThreatDetected)
			if tt.expectedThreat {
				assert.GreaterOrEqual(t, result.ThreatScore, tt.expectedScore-0.1)
				assert.LessOrEqual(t, result.ThreatScore, tt.expectedScore+0.1)
			}
		})
	}
}

func TestBehaviorAnalyzer_AnalyzeBehavior(t *testing.T) {
	logger := &observability.Logger{}
	analyzer := NewBehaviorAnalyzer(logger)

	userID := uuid.New()

	// Create a user behavior profile
	profile := &UserBehaviorProfile{
		UserID: userID,
		TypicalLoginTimes: []time.Time{
			time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
			time.Date(2024, 1, 1, 14, 0, 0, 0, time.UTC),
		},
		TypicalLocations: []Location{
			{Country: "US", City: "New York", IPRange: "192.168.1.0/24"},
		},
		TypicalDevices: []string{"device-123"},
		LastUpdated:    time.Now(),
	}

	// Test normal behavior
	normalRequest := &AccessRequest{
		UserID:    &userID,
		IPAddress: "192.168.1.100",
		Timestamp: time.Date(2024, 1, 1, 9, 30, 0, 0, time.UTC),
	}

	riskScore := analyzer.AnalyzeBehavior(profile, normalRequest)
	assert.LessOrEqual(t, riskScore, 0.3) // Should be low risk

	// Test anomalous behavior
	anomalousRequest := &AccessRequest{
		UserID:    &userID,
		IPAddress: "10.0.0.1",                                  // Different IP
		Timestamp: time.Date(2024, 1, 1, 3, 0, 0, 0, time.UTC), // Unusual time
	}

	riskScore = analyzer.AnalyzeBehavior(profile, anomalousRequest)
	assert.Greater(t, riskScore, 0.5) // Should be high risk
}

func TestRiskCalculator_CalculateRisk(t *testing.T) {
	logger := &observability.Logger{}
	calculator := NewRiskCalculator(logger)

	factors := &RiskFactors{
		DeviceTrust:  0.8,
		BehaviorRisk: 0.3,
		ThreatLevel:  0.2,
	}

	riskScore := calculator.CalculateRiskScore(factors)

	// Risk score should be reasonable based on factors
	assert.GreaterOrEqual(t, riskScore, 0.0)
	assert.LessOrEqual(t, riskScore, 1.0)
}

func BenchmarkZeroTrustEngine_EvaluateAccess(b *testing.B) {
	logger := &observability.Logger{}
	engine := NewZeroTrustEngine(logger)

	request := &AccessRequest{
		UserID:    &[]uuid.UUID{uuid.New()}[0],
		IPAddress: "192.168.1.100",
		Resource:  "/api/dashboard",
		Action:    "GET",
		Timestamp: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.EvaluateAccess(context.Background(), request)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAdvancedThreatDetector_DetectThreats(b *testing.B) {
	logger := &observability.Logger{}
	detector := NewAdvancedThreatDetector(logger)

	request := &SecurityRequest{
		RequestID: "bench-req",
		IPAddress: "192.168.1.100",
		UserAgent: "Mozilla/5.0",
		Method:    "GET",
		URL:       "/api/dashboard",
		Timestamp: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := detector.DetectThreats(context.Background(), request)
		if err != nil {
			b.Fatal(err)
		}
	}
}
