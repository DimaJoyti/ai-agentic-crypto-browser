package web3

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ai-agentic-browser/internal/config"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRiskAssessmentService(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{})
	clients := make(map[int]*ethclient.Client)

	service := NewRiskAssessmentService(clients, logger)

	t.Run("ServiceInitialization", func(t *testing.T) {
		assert.NotNil(t, service)
		assert.NotNil(t, service.mlModels)
		assert.NotNil(t, service.riskCache)
		assert.Equal(t, 3, len(service.mlModels))
	})

	t.Run("MLModelsInitialized", func(t *testing.T) {
		models := []string{"transaction_risk", "contract_risk", "rug_pull_detection"}
		for _, modelName := range models {
			model, exists := service.mlModels[modelName]
			assert.True(t, exists, "Model %s should exist", modelName)
			assert.NotEmpty(t, model.Name)
			assert.NotEmpty(t, model.Version)
			assert.NotEmpty(t, model.Features)
			assert.NotEmpty(t, model.Weights)
		}
	})
}

func TestTransactionRiskAssessment(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{})
	clients := make(map[int]*ethclient.Client)
	service := NewRiskAssessmentService(clients, logger)

	t.Run("HighValueTransaction", func(t *testing.T) {
		req := TransactionRiskRequest{
			FromAddress:     "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b",
			ToAddress:       "0x1234567890123456789012345678901234567890",
			Value:           big.NewInt(0).Mul(big.NewInt(200), big.NewInt(1e18)), // 200 ETH
			ChainID:         1,
			GasLimit:        21000,
			GasPrice:        big.NewInt(20000000000), // 20 gwei
			IncludeMLModels: false,
		}

		assessment, err := service.AssessTransactionRisk(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, assessment)
		assert.Greater(t, assessment.RiskScore, 0)
		assert.NotEmpty(t, assessment.Factors)

		// Should have high value factor
		hasHighValueFactor := false
		for _, factor := range assessment.Factors {
			if factor.Type == "high_value" {
				hasHighValueFactor = true
				break
			}
		}
		assert.True(t, hasHighValueFactor, "Should detect high value transaction")
	})

	t.Run("MaliciousAddress", func(t *testing.T) {
		req := TransactionRiskRequest{
			FromAddress:     "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b",
			ToAddress:       "0x0000000000000000000000000000000000000000", // Null address
			Value:           big.NewInt(1000000000000000000),              // 1 ETH
			ChainID:         1,
			GasLimit:        21000,
			GasPrice:        big.NewInt(20000000000),
			IncludeMLModels: false,
		}

		assessment, err := service.AssessTransactionRisk(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, assessment)

		// Should detect malicious address
		hasMaliciousFactor := false
		for _, factor := range assessment.Factors {
			if factor.Type == "malicious_address" {
				hasMaliciousFactor = true
				break
			}
		}
		assert.True(t, hasMaliciousFactor, "Should detect malicious address")
		assert.NotEmpty(t, assessment.Warnings)
	})

	t.Run("ContractInteraction", func(t *testing.T) {
		req := TransactionRiskRequest{
			FromAddress:     "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b",
			ToAddress:       "0x1234567890123456789012345678901234567890",
			Value:           big.NewInt(0),
			Data:            "0xa9059cbb000000000000000000000000742d35cc6634c0532925a3b8d4c9db96c4b4d8b0000000000000000000000000000000000000000000000000de0b6b3a7640000", // transfer function
			ChainID:         1,
			GasLimit:        50000,
			GasPrice:        big.NewInt(20000000000),
			IncludeMLModels: false,
		}

		assessment, err := service.AssessTransactionRisk(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, assessment)

		// Should detect token interaction
		hasTokenInteraction := false
		for _, factor := range assessment.Factors {
			if factor.Type == "token_interaction" {
				hasTokenInteraction = true
				break
			}
		}
		assert.True(t, hasTokenInteraction, "Should detect token interaction")
	})

	t.Run("SafetyGradeCalculation", func(t *testing.T) {
		testCases := []struct {
			name          string
			factors       []RiskFactor
			expectedGrade SafetyGrade
		}{
			{
				"ExcellentSafety",
				[]RiskFactor{{Impact: 0.05, Weight: 1.0}}, // Low risk
				SafetyGradeA,
			},
			{
				"GoodSafety",
				[]RiskFactor{{Impact: 0.25, Weight: 1.0}}, // Medium-low risk
				SafetyGradeB,
			},
			{
				"FairSafety",
				[]RiskFactor{{Impact: 0.45, Weight: 1.0}}, // Medium risk
				SafetyGradeB,
			},
			{
				"PoorSafety",
				[]RiskFactor{{Impact: 0.65, Weight: 1.0}}, // High risk
				SafetyGradeC,
			},
			{
				"FailSafety",
				[]RiskFactor{{Impact: 0.95, Weight: 1.0}}, // Critical risk
				SafetyGradeF,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				assessment := &RiskAssessment{
					Factors: tc.factors,
				}
				service.calculateFinalRisk(assessment)
				assert.Equal(t, tc.expectedGrade, assessment.SafetyGrade)
			})
		}
	})
}

func TestVulnerabilityScanner(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{})
	clients := make(map[int]*ethclient.Client)
	scanner := NewVulnerabilityScanner(clients, logger)

	t.Run("ScannerInitialization", func(t *testing.T) {
		assert.NotNil(t, scanner)
		assert.NotNil(t, scanner.rules)
		assert.Greater(t, len(scanner.rules), 0)
	})

	t.Run("VulnerabilityRules", func(t *testing.T) {
		expectedRules := []string{
			"reentrancy",
			"integer_overflow",
			"unchecked_call",
			"selfdestruct",
			"delegatecall",
			"timestamp_dependence",
			"weak_randomness",
			"unprotected_function",
			"gas_limit_dos",
			"front_running",
		}

		for _, ruleID := range expectedRules {
			rule, exists := scanner.rules[ruleID]
			assert.True(t, exists, "Rule %s should exist", ruleID)
			assert.NotEmpty(t, rule.Name)
			assert.NotEmpty(t, rule.Description)
			assert.NotEmpty(t, rule.Patterns)
			assert.NotEmpty(t, rule.Remediation)
		}
	})

	t.Run("SeverityLevels", func(t *testing.T) {
		severities := []VulnerabilitySeverity{
			SeverityCritical,
			SeverityHigh,
			SeverityMedium,
			SeverityLow,
			SeverityInfo,
		}

		for _, severity := range severities {
			assert.NotEmpty(t, string(severity))
		}
	})

	t.Run("VulnerabilityCategories", func(t *testing.T) {
		categories := []VulnerabilityCategory{
			CategoryReentrancy,
			CategoryOverflow,
			CategoryAccessControl,
			CategoryLogicError,
			CategoryDenialOfService,
			CategoryFrontRunning,
			CategoryTimestamp,
			CategoryRandomness,
		}

		for _, category := range categories {
			assert.NotEmpty(t, string(category))
		}
	})
}

func TestRiskMonitor(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{})
	clients := make(map[int]*ethclient.Client)
	riskAssessment := NewRiskAssessmentService(clients, logger)
	vulnScanner := NewVulnerabilityScanner(clients, logger)

	monitor := NewRiskMonitor(clients, logger, riskAssessment, vulnScanner)

	t.Run("MonitorInitialization", func(t *testing.T) {
		assert.NotNil(t, monitor)
		assert.NotNil(t, monitor.alertRules)
		assert.NotNil(t, monitor.alertChannels)
		assert.NotNil(t, monitor.monitoredAddresses)
		assert.Greater(t, len(monitor.alertRules), 0)
		assert.Greater(t, len(monitor.alertChannels), 0)
	})

	t.Run("AddMonitoredAddress", func(t *testing.T) {
		userID := uuid.New()
		address := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b"
		chainID := 1
		alertRules := []string{"high_risk_score"}

		err := monitor.AddMonitoredAddress(context.Background(), address, chainID, userID, alertRules)
		assert.NoError(t, err)

		key := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b:1"
		monitoredAddr, exists := monitor.monitoredAddresses[key]
		assert.True(t, exists)
		assert.Equal(t, address, monitoredAddr.Address)
		assert.Equal(t, chainID, monitoredAddr.ChainID)
		assert.Equal(t, userID, monitoredAddr.UserID)
		assert.Equal(t, alertRules, monitoredAddr.AlertRules)
	})

	t.Run("RemoveMonitoredAddress", func(t *testing.T) {
		address := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b"
		chainID := 1

		err := monitor.RemoveMonitoredAddress(context.Background(), address, chainID)
		assert.NoError(t, err)

		key := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b:1"
		_, exists := monitor.monitoredAddresses[key]
		assert.False(t, exists)
	})

	t.Run("AlertRules", func(t *testing.T) {
		expectedRules := []string{
			"high_risk_score",
			"critical_risk_score",
			"poor_safety_grade",
		}

		for _, ruleID := range expectedRules {
			rule, exists := monitor.alertRules[ruleID]
			assert.True(t, exists, "Alert rule %s should exist", ruleID)
			assert.NotEmpty(t, rule.Name)
			assert.NotEmpty(t, rule.Description)
			assert.NotEmpty(t, rule.Conditions)
			assert.NotEmpty(t, rule.Actions)
			assert.True(t, rule.Enabled)
		}
	})

	t.Run("AlertChannels", func(t *testing.T) {
		expectedChannels := []string{"log", "webhook", "email"}

		for _, channelType := range expectedChannels {
			channel, exists := monitor.alertChannels[channelType]
			assert.True(t, exists, "Alert channel %s should exist", channelType)
			assert.Equal(t, channelType, channel.GetType())
		}
	})

	t.Run("AlertPriorities", func(t *testing.T) {
		priorities := []AlertPriority{
			AlertPriorityCritical,
			AlertPriorityHigh,
			AlertPriorityMedium,
			AlertPriorityLow,
			AlertPriorityInfo,
		}

		for _, priority := range priorities {
			assert.NotEmpty(t, string(priority))
		}
	})
}

func TestRiskFactorAnalysis(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{})
	clients := make(map[int]*ethclient.Client)
	service := NewRiskAssessmentService(clients, logger)

	t.Run("HighGasLimit", func(t *testing.T) {
		assessment := &RiskAssessment{
			Factors: []RiskFactor{},
		}

		service.analyzeGasSettings(2000000, big.NewInt(20000000000), assessment)

		hasHighGasLimitFactor := false
		for _, factor := range assessment.Factors {
			if factor.Type == "high_gas_limit" {
				hasHighGasLimitFactor = true
				break
			}
		}
		assert.True(t, hasHighGasLimitFactor, "Should detect high gas limit")
	})

	t.Run("HighGasPrice", func(t *testing.T) {
		assessment := &RiskAssessment{
			Factors: []RiskFactor{},
		}

		highGasPrice := big.NewInt(0).Mul(big.NewInt(300), big.NewInt(1e9)) // 300 gwei
		service.analyzeGasSettings(21000, highGasPrice, assessment)

		hasHighGasPriceFactor := false
		for _, factor := range assessment.Factors {
			if factor.Type == "high_gas_price" {
				hasHighGasPriceFactor = true
				break
			}
		}
		assert.True(t, hasHighGasPriceFactor, "Should detect high gas price")
	})

	t.Run("SuspiciousPattern", func(t *testing.T) {
		assessment := &RiskAssessment{
			Factors: []RiskFactor{},
		}

		service.analyzeTransactionPatterns(
			"0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b",
			"0x000000123456789012345678901234567890",
			assessment,
		)

		hasSuspiciousPattern := false
		for _, factor := range assessment.Factors {
			if factor.Type == "suspicious_pattern" {
				hasSuspiciousPattern = true
				break
			}
		}
		assert.True(t, hasSuspiciousPattern, "Should detect suspicious pattern")
	})
}

func TestCacheManagement(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{})
	clients := make(map[int]*ethclient.Client)
	service := NewRiskAssessmentService(clients, logger)

	t.Run("CacheKeyGeneration", func(t *testing.T) {
		req := TransactionRiskRequest{
			FromAddress: "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b",
			ToAddress:   "0x1234567890123456789012345678901234567890",
			Value:       big.NewInt(1000000000000000000),
			ChainID:     1,
			GasLimit:    21000,
			Data:        "0x",
		}

		key1 := service.generateTransactionCacheKey(req)
		key2 := service.generateTransactionCacheKey(req)
		assert.Equal(t, key1, key2, "Same request should generate same cache key")

		req.Value = big.NewInt(2000000000000000000)
		key3 := service.generateTransactionCacheKey(req)
		assert.NotEqual(t, key1, key3, "Different request should generate different cache key")
	})

	t.Run("CacheExpiration", func(t *testing.T) {
		assessment := &RiskAssessment{
			ID:        uuid.New(),
			ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired
		}

		service.cacheAssessment("test_key", assessment)
		cached := service.getCachedAssessment("test_key")
		assert.Nil(t, cached, "Expired cache entry should return nil")
	})
}
