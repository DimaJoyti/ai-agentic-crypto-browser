package web3

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
)

// RiskAssessmentService provides ML-based risk assessment for transactions and contracts
type RiskAssessmentService struct {
	clients   map[int]*ethclient.Client
	logger    *observability.Logger
	mlModels  map[string]*MLRiskModel
	riskCache map[string]*RiskAssessment
	config    RiskConfig
}

// RiskConfig holds configuration for risk assessment
type RiskConfig struct {
	EnableMLModels       bool          `json:"enable_ml_models"`
	CacheTimeout         time.Duration `json:"cache_timeout"`
	MaxRiskScore         int           `json:"max_risk_score"`
	HighRiskThreshold    int           `json:"high_risk_threshold"`
	MediumRiskThreshold  int           `json:"medium_risk_threshold"`
	ContractAgeThreshold time.Duration `json:"contract_age_threshold"`
	MinLiquidityUSD      float64       `json:"min_liquidity_usd"`
}

// RiskLevel represents the risk level of a transaction or contract
type RiskLevel string

const (
	RiskLevelVeryLow  RiskLevel = "very_low"
	RiskLevelLow      RiskLevel = "low"
	RiskLevelMedium   RiskLevel = "medium"
	RiskLevelHigh     RiskLevel = "high"
	RiskLevelCritical RiskLevel = "critical"
)

// SafetyGrade represents A-F safety grading
type SafetyGrade string

const (
	SafetyGradeA SafetyGrade = "A" // Excellent (90-100)
	SafetyGradeB SafetyGrade = "B" // Good (80-89)
	SafetyGradeC SafetyGrade = "C" // Fair (70-79)
	SafetyGradeD SafetyGrade = "D" // Poor (60-69)
	SafetyGradeF SafetyGrade = "F" // Fail (0-59)
)

// RiskAssessment represents a comprehensive risk assessment
type RiskAssessment struct {
	ID              uuid.UUID              `json:"id"`
	TransactionHash string                 `json:"transaction_hash,omitempty"`
	ContractAddress string                 `json:"contract_address,omitempty"`
	ChainID         int                    `json:"chain_id"`
	RiskScore       int                    `json:"risk_score"`   // 0-100
	SafetyGrade     SafetyGrade            `json:"safety_grade"` // A-F
	RiskLevel       RiskLevel              `json:"risk_level"`
	Confidence      float64                `json:"confidence"` // 0.0-1.0
	Factors         []RiskFactor           `json:"factors"`
	Recommendations []string               `json:"recommendations"`
	Warnings        []string               `json:"warnings"`
	MLPredictions   map[string]float64     `json:"ml_predictions"`
	AssessedAt      time.Time              `json:"assessed_at"`
	ExpiresAt       time.Time              `json:"expires_at"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// RiskFactor represents an individual risk factor
type RiskFactor struct {
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Impact      float64 `json:"impact"` // -1.0 to 1.0 (negative is good, positive is bad)
	Weight      float64 `json:"weight"` // 0.0 to 1.0
	Evidence    string  `json:"evidence"`
}

// TransactionRiskRequest represents a request for transaction risk assessment
type TransactionRiskRequest struct {
	FromAddress     string                 `json:"from_address"`
	ToAddress       string                 `json:"to_address"`
	Value           *big.Int               `json:"value"`
	Data            string                 `json:"data"`
	ChainID         int                    `json:"chain_id"`
	GasLimit        uint64                 `json:"gas_limit"`
	GasPrice        *big.Int               `json:"gas_price"`
	IncludeMLModels bool                   `json:"include_ml_models"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// ContractRiskRequest represents a request for contract risk assessment
type ContractRiskRequest struct {
	ContractAddress string                 `json:"contract_address"`
	ChainID         int                    `json:"chain_id"`
	AnalyzeCode     bool                   `json:"analyze_code"`
	CheckRugPull    bool                   `json:"check_rug_pull"`
	IncludeMLModels bool                   `json:"include_ml_models"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// MLRiskModel represents a machine learning model for risk assessment
type MLRiskModel struct {
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Type        string                 `json:"type"` // "classification", "regression", "anomaly_detection"
	Features    []string               `json:"features"`
	Weights     map[string]float64     `json:"weights"`
	Thresholds  map[string]float64     `json:"thresholds"`
	Accuracy    float64                `json:"accuracy"`
	LastTrained time.Time              `json:"last_trained"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// NewRiskAssessmentService creates a new risk assessment service
func NewRiskAssessmentService(clients map[int]*ethclient.Client, logger *observability.Logger) *RiskAssessmentService {
	config := RiskConfig{
		EnableMLModels:       true,
		CacheTimeout:         15 * time.Minute,
		MaxRiskScore:         100,
		HighRiskThreshold:    70,
		MediumRiskThreshold:  40,
		ContractAgeThreshold: 30 * 24 * time.Hour, // 30 days
		MinLiquidityUSD:      10000,               // $10k minimum liquidity
	}

	service := &RiskAssessmentService{
		clients:   clients,
		logger:    logger,
		mlModels:  make(map[string]*MLRiskModel),
		riskCache: make(map[string]*RiskAssessment),
		config:    config,
	}

	// Initialize ML models
	service.initializeMLModels()

	return service
}

// AssessTransactionRisk performs comprehensive risk assessment for a transaction
func (r *RiskAssessmentService) AssessTransactionRisk(ctx context.Context, req TransactionRiskRequest) (*RiskAssessment, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("risk-assessment").Start(ctx, "risk.AssessTransactionRisk")
	defer span.End()

	// Generate cache key
	cacheKey := r.generateTransactionCacheKey(req)

	// Check cache first
	if cached := r.getCachedAssessment(cacheKey); cached != nil {
		r.logger.Info(ctx, "Risk assessment found in cache", map[string]interface{}{
			"cache_key": cacheKey,
		})
		return cached, nil
	}

	// Create new assessment
	assessment := &RiskAssessment{
		ID:              uuid.New(),
		ChainID:         req.ChainID,
		RiskScore:       0,
		Confidence:      0.0,
		Factors:         []RiskFactor{},
		Recommendations: []string{},
		Warnings:        []string{},
		MLPredictions:   make(map[string]float64),
		AssessedAt:      time.Now(),
		ExpiresAt:       time.Now().Add(r.config.CacheTimeout),
		Metadata:        req.Metadata,
	}

	// Analyze transaction factors
	if err := r.analyzeTransactionFactors(ctx, req, assessment); err != nil {
		return nil, fmt.Errorf("failed to analyze transaction factors: %w", err)
	}

	// Apply ML models if enabled
	if req.IncludeMLModels && r.config.EnableMLModels {
		if err := r.applyMLModels(ctx, req, assessment); err != nil {
			r.logger.Warn(ctx, "ML model application failed", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	// Calculate final risk score and grade
	r.calculateFinalRisk(assessment)

	// Cache the assessment
	r.cacheAssessment(cacheKey, assessment)

	r.logger.Info(ctx, "Transaction risk assessment completed", map[string]interface{}{
		"assessment_id": assessment.ID.String(),
		"risk_score":    assessment.RiskScore,
		"safety_grade":  string(assessment.SafetyGrade),
		"risk_level":    string(assessment.RiskLevel),
	})

	return assessment, nil
}

// AssessContractRisk performs comprehensive risk assessment for a smart contract
func (r *RiskAssessmentService) AssessContractRisk(ctx context.Context, req ContractRiskRequest) (*RiskAssessment, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("risk-assessment").Start(ctx, "risk.AssessContractRisk")
	defer span.End()

	// Generate cache key
	cacheKey := r.generateContractCacheKey(req)

	// Check cache first
	if cached := r.getCachedAssessment(cacheKey); cached != nil {
		r.logger.Info(ctx, "Contract risk assessment found in cache", map[string]interface{}{
			"cache_key": cacheKey,
		})
		return cached, nil
	}

	// Create new assessment
	assessment := &RiskAssessment{
		ID:              uuid.New(),
		ContractAddress: req.ContractAddress,
		ChainID:         req.ChainID,
		RiskScore:       0,
		Confidence:      0.0,
		Factors:         []RiskFactor{},
		Recommendations: []string{},
		Warnings:        []string{},
		MLPredictions:   make(map[string]float64),
		AssessedAt:      time.Now(),
		ExpiresAt:       time.Now().Add(r.config.CacheTimeout),
		Metadata:        req.Metadata,
	}

	// Analyze contract factors
	if err := r.analyzeContractFactors(ctx, req, assessment); err != nil {
		return nil, fmt.Errorf("failed to analyze contract factors: %w", err)
	}

	// Analyze contract code if requested
	if req.AnalyzeCode {
		if err := r.analyzeContractCode(ctx, req, assessment); err != nil {
			r.logger.Warn(ctx, "Contract code analysis failed", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	// Check for rug pull indicators if requested
	if req.CheckRugPull {
		if err := r.checkRugPullIndicators(ctx, req, assessment); err != nil {
			r.logger.Warn(ctx, "Rug pull analysis failed", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	// Apply ML models if enabled
	if req.IncludeMLModels && r.config.EnableMLModels {
		if err := r.applyContractMLModels(ctx, req, assessment); err != nil {
			r.logger.Warn(ctx, "Contract ML model application failed", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	// Calculate final risk score and grade
	r.calculateFinalRisk(assessment)

	// Cache the assessment
	r.cacheAssessment(cacheKey, assessment)

	r.logger.Info(ctx, "Contract risk assessment completed", map[string]interface{}{
		"assessment_id":    assessment.ID.String(),
		"contract_address": req.ContractAddress,
		"risk_score":       assessment.RiskScore,
		"safety_grade":     string(assessment.SafetyGrade),
		"risk_level":       string(assessment.RiskLevel),
	})

	return assessment, nil
}

// analyzeTransactionFactors analyzes various risk factors for a transaction
func (r *RiskAssessmentService) analyzeTransactionFactors(ctx context.Context, req TransactionRiskRequest, assessment *RiskAssessment) error {
	// Check if addresses are known malicious
	r.checkMaliciousAddresses(req.FromAddress, req.ToAddress, assessment)

	// Analyze transaction value
	r.analyzeTransactionValue(req.Value, assessment)

	// Analyze gas settings
	r.analyzeGasSettings(req.GasLimit, req.GasPrice, assessment)

	// Check contract interaction
	if req.Data != "" && req.Data != "0x" {
		r.analyzeContractInteraction(req.ToAddress, req.Data, assessment)
	}

	// Analyze transaction patterns
	r.analyzeTransactionPatterns(req.FromAddress, req.ToAddress, assessment)

	return nil
}

// analyzeContractFactors analyzes various risk factors for a smart contract
func (r *RiskAssessmentService) analyzeContractFactors(ctx context.Context, req ContractRiskRequest, assessment *RiskAssessment) error {
	client, exists := r.clients[req.ChainID]
	if !exists {
		return fmt.Errorf("no client available for chain ID: %d", req.ChainID)
	}

	contractAddr := common.HexToAddress(req.ContractAddress)

	// Check contract age
	r.analyzeContractAge(ctx, client, contractAddr, assessment)

	// Check contract verification status
	r.analyzeContractVerification(req.ContractAddress, assessment)

	// Analyze contract activity
	r.analyzeContractActivity(ctx, client, contractAddr, assessment)

	// Check for proxy patterns
	r.analyzeProxyPatterns(ctx, client, contractAddr, assessment)

	return nil
}

// checkMaliciousAddresses checks if addresses are in known malicious lists
func (r *RiskAssessmentService) checkMaliciousAddresses(fromAddr, toAddr string, assessment *RiskAssessment) {
	// This would integrate with threat intelligence feeds
	// For now, implement basic checks

	maliciousAddresses := map[string]string{
		"0x0000000000000000000000000000000000000000": "Null address",
		// Add more known malicious addresses
	}

	if reason, isMalicious := maliciousAddresses[strings.ToLower(toAddr)]; isMalicious {
		assessment.Factors = append(assessment.Factors, RiskFactor{
			Type:        "malicious_address",
			Description: fmt.Sprintf("Destination address is known malicious: %s", reason),
			Impact:      0.9,
			Weight:      1.0,
			Evidence:    fmt.Sprintf("Address %s flagged as malicious", toAddr),
		})
		assessment.Warnings = append(assessment.Warnings, "Transaction involves known malicious address")
	}
}

// analyzeTransactionValue analyzes the transaction value for risk indicators
func (r *RiskAssessmentService) analyzeTransactionValue(value *big.Int, assessment *RiskAssessment) {
	if value == nil || value.Cmp(big.NewInt(0)) == 0 {
		return
	}

	// Convert to ETH (assuming 18 decimals)
	ethValue := new(big.Float).Quo(new(big.Float).SetInt(value), big.NewFloat(1e18))
	ethFloat, _ := ethValue.Float64()

	// High value transactions carry more risk
	if ethFloat > 100 {
		assessment.Factors = append(assessment.Factors, RiskFactor{
			Type:        "high_value",
			Description: fmt.Sprintf("High value transaction: %.2f ETH", ethFloat),
			Impact:      0.3,
			Weight:      0.7,
			Evidence:    fmt.Sprintf("Transaction value: %.6f ETH", ethFloat),
		})
		assessment.Recommendations = append(assessment.Recommendations, "Consider using hardware wallet for high-value transactions")
	}
}

// analyzeGasSettings analyzes gas limit and price for anomalies
func (r *RiskAssessmentService) analyzeGasSettings(gasLimit uint64, gasPrice *big.Int, assessment *RiskAssessment) {
	// Check for unusually high gas limit
	if gasLimit > 1000000 {
		assessment.Factors = append(assessment.Factors, RiskFactor{
			Type:        "high_gas_limit",
			Description: fmt.Sprintf("Unusually high gas limit: %d", gasLimit),
			Impact:      0.2,
			Weight:      0.5,
			Evidence:    fmt.Sprintf("Gas limit: %d", gasLimit),
		})
	}

	// Check for unusually high gas price
	if gasPrice != nil {
		gweiPrice := new(big.Float).Quo(new(big.Float).SetInt(gasPrice), big.NewFloat(1e9))
		gweiFloat, _ := gweiPrice.Float64()

		if gweiFloat > 200 {
			assessment.Factors = append(assessment.Factors, RiskFactor{
				Type:        "high_gas_price",
				Description: fmt.Sprintf("High gas price: %.2f gwei", gweiFloat),
				Impact:      0.1,
				Weight:      0.3,
				Evidence:    fmt.Sprintf("Gas price: %.2f gwei", gweiFloat),
			})
		}
	}
}

// analyzeContractInteraction analyzes contract interaction data
func (r *RiskAssessmentService) analyzeContractInteraction(toAddr, data string, assessment *RiskAssessment) {
	// Check for common dangerous function signatures
	dangerousFunctions := map[string]string{
		"0xa9059cbb": "transfer(address,uint256) - Token transfer",
		"0x095ea7b3": "approve(address,uint256) - Token approval",
		"0x23b872dd": "transferFrom(address,address,uint256) - Transfer from",
	}

	if len(data) >= 10 {
		funcSig := data[:10]
		if desc, isDangerous := dangerousFunctions[funcSig]; isDangerous {
			assessment.Factors = append(assessment.Factors, RiskFactor{
				Type:        "token_interaction",
				Description: fmt.Sprintf("Token interaction detected: %s", desc),
				Impact:      0.2,
				Weight:      0.6,
				Evidence:    fmt.Sprintf("Function signature: %s", funcSig),
			})
		}
	}
}

// analyzeTransactionPatterns analyzes patterns in transaction history
func (r *RiskAssessmentService) analyzeTransactionPatterns(fromAddr, toAddr string, assessment *RiskAssessment) {
	// This would analyze on-chain transaction patterns
	// For now, implement basic pattern detection

	// Check for common scam patterns
	if strings.HasPrefix(strings.ToLower(toAddr), "0x000000") {
		assessment.Factors = append(assessment.Factors, RiskFactor{
			Type:        "suspicious_pattern",
			Description: "Destination address follows suspicious pattern",
			Impact:      0.4,
			Weight:      0.7,
			Evidence:    fmt.Sprintf("Address pattern: %s", toAddr),
		})
	}
}
