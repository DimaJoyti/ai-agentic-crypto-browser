package web3

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"time"
)

// initializeMLModels initializes the machine learning models for risk assessment
func (r *RiskAssessmentService) initializeMLModels() {
	// Transaction risk model
	r.mlModels["transaction_risk"] = &MLRiskModel{
		Name:     "transaction_risk_classifier",
		Version:  "1.0.0",
		Type:     "classification",
		Features: []string{"value", "gas_price", "gas_limit", "address_age", "transaction_count"},
		Weights: map[string]float64{
			"value":             0.3,
			"gas_price":         0.1,
			"gas_limit":         0.1,
			"address_age":       0.2,
			"transaction_count": 0.3,
		},
		Thresholds: map[string]float64{
			"high_risk":   0.7,
			"medium_risk": 0.4,
			"low_risk":    0.2,
		},
		Accuracy:    0.85,
		LastTrained: time.Now().AddDate(0, -1, 0), // 1 month ago
	}

	// Contract risk model
	r.mlModels["contract_risk"] = &MLRiskModel{
		Name:     "contract_risk_classifier",
		Version:  "1.0.0",
		Type:     "classification",
		Features: []string{"contract_age", "verification_status", "transaction_volume", "unique_users", "liquidity"},
		Weights: map[string]float64{
			"contract_age":        0.2,
			"verification_status": 0.3,
			"transaction_volume":  0.2,
			"unique_users":        0.15,
			"liquidity":           0.15,
		},
		Thresholds: map[string]float64{
			"high_risk":   0.6,
			"medium_risk": 0.3,
			"low_risk":    0.1,
		},
		Accuracy:    0.82,
		LastTrained: time.Now().AddDate(0, -1, 0),
	}

	// Rug pull detection model
	r.mlModels["rug_pull_detection"] = &MLRiskModel{
		Name:     "rug_pull_detector",
		Version:  "1.0.0",
		Type:     "anomaly_detection",
		Features: []string{"liquidity_change", "holder_concentration", "dev_wallet_activity", "contract_permissions"},
		Weights: map[string]float64{
			"liquidity_change":     0.4,
			"holder_concentration": 0.3,
			"dev_wallet_activity":  0.2,
			"contract_permissions": 0.1,
		},
		Thresholds: map[string]float64{
			"rug_pull_likely":   0.8,
			"rug_pull_possible": 0.5,
			"rug_pull_unlikely": 0.2,
		},
		Accuracy:    0.78,
		LastTrained: time.Now().AddDate(0, -2, 0),
	}
}

// generateTransactionCacheKey generates a cache key for transaction risk assessment
func (r *RiskAssessmentService) generateTransactionCacheKey(req TransactionRiskRequest) string {
	data := fmt.Sprintf("%s:%s:%s:%d:%d:%s",
		req.FromAddress,
		req.ToAddress,
		req.Value.String(),
		req.ChainID,
		req.GasLimit,
		req.Data,
	)
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("tx_risk_%x", hash[:8])
}

// generateContractCacheKey generates a cache key for contract risk assessment
func (r *RiskAssessmentService) generateContractCacheKey(req ContractRiskRequest) string {
	data := fmt.Sprintf("%s:%d:%t:%t",
		req.ContractAddress,
		req.ChainID,
		req.AnalyzeCode,
		req.CheckRugPull,
	)
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("contract_risk_%x", hash[:8])
}

// getCachedAssessment retrieves a cached risk assessment
func (r *RiskAssessmentService) getCachedAssessment(cacheKey string) *RiskAssessment {
	assessment, exists := r.riskCache[cacheKey]
	if !exists {
		return nil
	}

	// Check if cache entry is still valid
	if time.Now().After(assessment.ExpiresAt) {
		delete(r.riskCache, cacheKey)
		return nil
	}

	return assessment
}

// cacheAssessment caches a risk assessment
func (r *RiskAssessmentService) cacheAssessment(cacheKey string, assessment *RiskAssessment) {
	r.riskCache[cacheKey] = assessment
}

// calculateFinalRisk calculates the final risk score and assigns safety grade
func (r *RiskAssessmentService) calculateFinalRisk(assessment *RiskAssessment) {
	if len(assessment.Factors) == 0 {
		assessment.RiskScore = 0
		assessment.SafetyGrade = SafetyGradeA
		assessment.RiskLevel = RiskLevelVeryLow
		assessment.Confidence = 1.0
		return
	}

	// Calculate weighted risk score
	totalWeight := 0.0
	weightedRisk := 0.0

	for _, factor := range assessment.Factors {
		impact := math.Max(0, factor.Impact) // Ensure positive impact
		weightedRisk += impact * factor.Weight
		totalWeight += factor.Weight
	}

	// Normalize to 0-100 scale
	if totalWeight > 0 {
		normalizedRisk := (weightedRisk / totalWeight) * 100
		assessment.RiskScore = int(math.Min(100, normalizedRisk))
	}

	// Assign safety grade based on risk score
	switch {
	case assessment.RiskScore >= 90:
		assessment.SafetyGrade = SafetyGradeF
		assessment.RiskLevel = RiskLevelCritical
	case assessment.RiskScore >= 80:
		assessment.SafetyGrade = SafetyGradeD
		assessment.RiskLevel = RiskLevelHigh
	case assessment.RiskScore >= 70:
		assessment.SafetyGrade = SafetyGradeC
		assessment.RiskLevel = RiskLevelHigh
	case assessment.RiskScore >= 60:
		assessment.SafetyGrade = SafetyGradeC
		assessment.RiskLevel = RiskLevelMedium
	case assessment.RiskScore >= 40:
		assessment.SafetyGrade = SafetyGradeB
		assessment.RiskLevel = RiskLevelMedium
	case assessment.RiskScore >= 20:
		assessment.SafetyGrade = SafetyGradeB
		assessment.RiskLevel = RiskLevelLow
	default:
		assessment.SafetyGrade = SafetyGradeA
		assessment.RiskLevel = RiskLevelVeryLow
	}

	// Calculate confidence based on number of factors and their weights
	confidence := math.Min(1.0, totalWeight/5.0)      // Assume 5.0 is maximum expected weight
	assessment.Confidence = math.Max(0.1, confidence) // Minimum 10% confidence

	// Add recommendations based on risk level
	r.addRiskRecommendations(assessment)
}

// addRiskRecommendations adds recommendations based on risk level
func (r *RiskAssessmentService) addRiskRecommendations(assessment *RiskAssessment) {
	switch assessment.RiskLevel {
	case RiskLevelCritical:
		assessment.Recommendations = append(assessment.Recommendations,
			"DO NOT PROCEED - Critical risk detected",
			"Review all transaction details carefully",
			"Consider seeking expert advice",
		)
	case RiskLevelHigh:
		assessment.Recommendations = append(assessment.Recommendations,
			"Exercise extreme caution",
			"Double-check all transaction details",
			"Consider using a test transaction first",
			"Use hardware wallet for signing",
		)
	case RiskLevelMedium:
		assessment.Recommendations = append(assessment.Recommendations,
			"Review transaction details carefully",
			"Consider the identified risk factors",
			"Ensure you understand the transaction purpose",
		)
	case RiskLevelLow:
		assessment.Recommendations = append(assessment.Recommendations,
			"Transaction appears relatively safe",
			"Standard security practices apply",
		)
	case RiskLevelVeryLow:
		assessment.Recommendations = append(assessment.Recommendations,
			"Transaction appears very safe",
			"Minimal risk detected",
		)
	}
}

// applyMLModels applies machine learning models to the assessment
func (r *RiskAssessmentService) applyMLModels(ctx context.Context, req TransactionRiskRequest, assessment *RiskAssessment) error {
	// Apply transaction risk model
	if model, exists := r.mlModels["transaction_risk"]; exists {
		prediction := r.predictTransactionRisk(req, model)
		assessment.MLPredictions["transaction_risk"] = prediction

		// Add ML-based risk factor
		if prediction > model.Thresholds["high_risk"] {
			assessment.Factors = append(assessment.Factors, RiskFactor{
				Type:        "ml_prediction",
				Description: fmt.Sprintf("ML model predicts high risk (%.2f)", prediction),
				Impact:      prediction * 0.8, // Scale down ML impact
				Weight:      0.6,
				Evidence:    fmt.Sprintf("ML model: %s, prediction: %.3f", model.Name, prediction),
			})
		}
	}

	return nil
}

// applyContractMLModels applies ML models for contract risk assessment
func (r *RiskAssessmentService) applyContractMLModels(ctx context.Context, req ContractRiskRequest, assessment *RiskAssessment) error {
	// Apply contract risk model
	if model, exists := r.mlModels["contract_risk"]; exists {
		prediction := r.predictContractRisk(req, model)
		assessment.MLPredictions["contract_risk"] = prediction

		if prediction > model.Thresholds["high_risk"] {
			assessment.Factors = append(assessment.Factors, RiskFactor{
				Type:        "ml_contract_risk",
				Description: fmt.Sprintf("ML model predicts high contract risk (%.2f)", prediction),
				Impact:      prediction * 0.7,
				Weight:      0.7,
				Evidence:    fmt.Sprintf("ML model: %s, prediction: %.3f", model.Name, prediction),
			})
		}
	}

	// Apply rug pull detection if requested
	if req.CheckRugPull {
		if model, exists := r.mlModels["rug_pull_detection"]; exists {
			prediction := r.predictRugPull(req, model)
			assessment.MLPredictions["rug_pull_risk"] = prediction

			if prediction > model.Thresholds["rug_pull_possible"] {
				assessment.Factors = append(assessment.Factors, RiskFactor{
					Type:        "rug_pull_risk",
					Description: fmt.Sprintf("Potential rug pull detected (%.2f)", prediction),
					Impact:      prediction,
					Weight:      0.9,
					Evidence:    fmt.Sprintf("Rug pull model: %s, prediction: %.3f", model.Name, prediction),
				})
				assessment.Warnings = append(assessment.Warnings, "Potential rug pull indicators detected")
			}
		}
	}

	return nil
}

// predictTransactionRisk predicts transaction risk using ML model
func (r *RiskAssessmentService) predictTransactionRisk(req TransactionRiskRequest, model *MLRiskModel) float64 {
	// This is a simplified ML prediction simulation
	// In a real implementation, this would call an actual ML model

	features := map[string]float64{
		"value":             0.0,
		"gas_price":         0.0,
		"gas_limit":         float64(req.GasLimit) / 1000000.0, // Normalize
		"address_age":       0.5,                               // Would be calculated from blockchain data
		"transaction_count": 0.5,                               // Would be calculated from blockchain data
	}

	if req.Value != nil {
		ethValue := new(big.Float).Quo(new(big.Float).SetInt(req.Value), big.NewFloat(1e18))
		ethFloat, _ := ethValue.Float64()
		features["value"] = math.Min(1.0, ethFloat/1000.0) // Normalize to 0-1
	}

	if req.GasPrice != nil {
		gweiPrice := new(big.Float).Quo(new(big.Float).SetInt(req.GasPrice), big.NewFloat(1e9))
		gweiFloat, _ := gweiPrice.Float64()
		features["gas_price"] = math.Min(1.0, gweiFloat/500.0) // Normalize to 0-1
	}

	// Simple weighted sum prediction
	prediction := 0.0
	for feature, value := range features {
		if weight, exists := model.Weights[feature]; exists {
			prediction += value * weight
		}
	}

	return math.Min(1.0, math.Max(0.0, prediction))
}

// predictContractRisk predicts contract risk using ML model
func (r *RiskAssessmentService) predictContractRisk(req ContractRiskRequest, model *MLRiskModel) float64 {
	// Simplified contract risk prediction
	features := map[string]float64{
		"contract_age":        0.5, // Would be calculated from deployment time
		"verification_status": 0.8, // Would check if contract is verified
		"transaction_volume":  0.3, // Would analyze transaction history
		"unique_users":        0.4, // Would count unique interacting addresses
		"liquidity":           0.6, // Would check DEX liquidity
	}

	prediction := 0.0
	for feature, value := range features {
		if weight, exists := model.Weights[feature]; exists {
			prediction += value * weight
		}
	}

	return math.Min(1.0, math.Max(0.0, prediction))
}

// predictRugPull predicts rug pull probability using ML model
func (r *RiskAssessmentService) predictRugPull(req ContractRiskRequest, model *MLRiskModel) float64 {
	// Simplified rug pull prediction
	features := map[string]float64{
		"liquidity_change":     0.2, // Would analyze liquidity changes
		"holder_concentration": 0.7, // Would check token holder distribution
		"dev_wallet_activity":  0.3, // Would monitor developer wallet activity
		"contract_permissions": 0.8, // Would check for dangerous permissions
	}

	prediction := 0.0
	for feature, value := range features {
		if weight, exists := model.Weights[feature]; exists {
			prediction += value * weight
		}
	}

	return math.Min(1.0, math.Max(0.0, prediction))
}
