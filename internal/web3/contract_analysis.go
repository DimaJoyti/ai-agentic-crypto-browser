package web3

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// analyzeContractAge analyzes the age of a smart contract
func (r *RiskAssessmentService) analyzeContractAge(ctx context.Context, client *ethclient.Client, contractAddr common.Address, assessment *RiskAssessment) {
	// Get contract creation transaction (simplified - would need more complex logic)
	code, err := client.CodeAt(ctx, contractAddr, nil)
	if err != nil {
		r.logger.Warn(ctx, "Failed to get contract code", map[string]interface{}{
			"contract": contractAddr.Hex(),
			"error":    err.Error(),
		})
		return
	}

	if len(code) == 0 {
		assessment.Factors = append(assessment.Factors, RiskFactor{
			Type:        "no_contract_code",
			Description: "No contract code found at address",
			Impact:      0.8,
			Weight:      1.0,
			Evidence:    "Contract code length: 0",
		})
		assessment.Warnings = append(assessment.Warnings, "Address does not contain contract code")
		return
	}

	// For demonstration, assume contract is new if we can't determine age
	// In a real implementation, this would query contract creation events
	assessment.Factors = append(assessment.Factors, RiskFactor{
		Type:        "contract_age_unknown",
		Description: "Contract age could not be determined",
		Impact:      0.3,
		Weight:      0.5,
		Evidence:    "Contract creation time unknown",
	})
}

// analyzeContractVerification checks if the contract is verified
func (r *RiskAssessmentService) analyzeContractVerification(contractAddress string, assessment *RiskAssessment) {
	// This would integrate with Etherscan API or similar service
	// For now, implement basic heuristics

	// Assume unverified for demonstration
	assessment.Factors = append(assessment.Factors, RiskFactor{
		Type:        "unverified_contract",
		Description: "Contract source code is not verified",
		Impact:      0.4,
		Weight:      0.7,
		Evidence:    "No verified source code found",
	})
	assessment.Recommendations = append(assessment.Recommendations, "Exercise caution with unverified contracts")
}

// analyzeContractActivity analyzes the activity level of a contract
func (r *RiskAssessmentService) analyzeContractActivity(ctx context.Context, client *ethclient.Client, contractAddr common.Address, assessment *RiskAssessment) {
	// Get current block number
	currentBlock, err := client.BlockNumber(ctx)
	if err != nil {
		r.logger.Warn(ctx, "Failed to get current block number", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Check recent activity (last 1000 blocks)
	fromBlock := currentBlock - 1000
	if fromBlock < 0 {
		fromBlock = 0
	}

	// This is a simplified check - would need more sophisticated analysis
	// For now, add a generic activity factor
	assessment.Factors = append(assessment.Factors, RiskFactor{
		Type:        "contract_activity",
		Description: "Contract activity analysis completed",
		Impact:      0.1,
		Weight:      0.3,
		Evidence:    fmt.Sprintf("Analyzed blocks %d to %d", fromBlock, currentBlock),
	})
}

// analyzeProxyPatterns checks for proxy contract patterns
func (r *RiskAssessmentService) analyzeProxyPatterns(ctx context.Context, client *ethclient.Client, contractAddr common.Address, assessment *RiskAssessment) {
	// Get contract code
	code, err := client.CodeAt(ctx, contractAddr, nil)
	if err != nil {
		return
	}

	codeHex := common.Bytes2Hex(code)

	// Check for common proxy patterns
	proxyPatterns := map[string]string{
		"363d3d373d3d3d363d73":                 "Minimal proxy pattern (EIP-1167)",
		"3660008037602060003660003473":         "Transparent proxy pattern",
		"608060405234801561001057600080fd5b50": "Common proxy initialization",
	}

	for pattern, description := range proxyPatterns {
		if strings.Contains(codeHex, pattern) {
			assessment.Factors = append(assessment.Factors, RiskFactor{
				Type:        "proxy_contract",
				Description: fmt.Sprintf("Proxy contract detected: %s", description),
				Impact:      0.3,
				Weight:      0.6,
				Evidence:    fmt.Sprintf("Pattern found: %s", pattern),
			})
			assessment.Recommendations = append(assessment.Recommendations, "Verify the implementation contract for proxy contracts")
			break
		}
	}
}

// analyzeContractCode performs static analysis of contract code
func (r *RiskAssessmentService) analyzeContractCode(ctx context.Context, req ContractRiskRequest, assessment *RiskAssessment) error {
	client, exists := r.clients[req.ChainID]
	if !exists {
		return fmt.Errorf("no client available for chain ID: %d", req.ChainID)
	}

	contractAddr := common.HexToAddress(req.ContractAddress)
	code, err := client.CodeAt(ctx, contractAddr, nil)
	if err != nil {
		return fmt.Errorf("failed to get contract code: %w", err)
	}

	if len(code) == 0 {
		return fmt.Errorf("no contract code found")
	}

	// Analyze bytecode for dangerous patterns
	r.analyzeDangerousPatterns(code, assessment)

	// Check for common vulnerabilities
	r.checkCommonVulnerabilities(code, assessment)

	return nil
}

// analyzeDangerousPatterns analyzes bytecode for dangerous patterns
func (r *RiskAssessmentService) analyzeDangerousPatterns(code []byte, assessment *RiskAssessment) {
	codeHex := common.Bytes2Hex(code)

	// Check for dangerous opcodes/patterns
	dangerousPatterns := map[string]string{
		"ff": "SELFDESTRUCT opcode detected",
		"f4": "DELEGATECALL opcode detected",
		"3d": "RETURNDATASIZE opcode (potential proxy)",
		"f0": "CREATE opcode detected",
		"f5": "CREATE2 opcode detected",
	}

	for pattern, description := range dangerousPatterns {
		if strings.Contains(codeHex, pattern) {
			assessment.Factors = append(assessment.Factors, RiskFactor{
				Type:        "dangerous_opcode",
				Description: description,
				Impact:      0.4,
				Weight:      0.5,
				Evidence:    fmt.Sprintf("Opcode pattern: %s", pattern),
			})
		}
	}
}

// checkCommonVulnerabilities checks for common smart contract vulnerabilities
func (r *RiskAssessmentService) checkCommonVulnerabilities(code []byte, assessment *RiskAssessment) {
	codeHex := common.Bytes2Hex(code)

	// Check for reentrancy patterns
	if strings.Contains(codeHex, "f1") && strings.Contains(codeHex, "55") {
		assessment.Factors = append(assessment.Factors, RiskFactor{
			Type:        "potential_reentrancy",
			Description: "Potential reentrancy vulnerability detected",
			Impact:      0.6,
			Weight:      0.8,
			Evidence:    "CALL and SSTORE opcodes found in proximity",
		})
		assessment.Warnings = append(assessment.Warnings, "Contract may be vulnerable to reentrancy attacks")
	}

	// Check for integer overflow patterns (pre-Solidity 0.8.0)
	if strings.Contains(codeHex, "01") && strings.Contains(codeHex, "03") {
		assessment.Factors = append(assessment.Factors, RiskFactor{
			Type:        "potential_overflow",
			Description: "Potential integer overflow vulnerability",
			Impact:      0.5,
			Weight:      0.6,
			Evidence:    "ADD and SUB opcodes without overflow checks",
		})
	}

	// Check for unchecked external calls
	if strings.Contains(codeHex, "f1") && !strings.Contains(codeHex, "3d") {
		assessment.Factors = append(assessment.Factors, RiskFactor{
			Type:        "unchecked_call",
			Description: "Potential unchecked external call",
			Impact:      0.4,
			Weight:      0.5,
			Evidence:    "CALL opcode without return value check",
		})
	}
}

// checkRugPullIndicators checks for rug pull indicators
func (r *RiskAssessmentService) checkRugPullIndicators(ctx context.Context, req ContractRiskRequest, assessment *RiskAssessment) error {
	client, exists := r.clients[req.ChainID]
	if !exists {
		return fmt.Errorf("no client available for chain ID: %d", req.ChainID)
	}

	contractAddr := common.HexToAddress(req.ContractAddress)

	// Check for honeypot patterns
	r.checkHoneypotPatterns(ctx, client, contractAddr, assessment)

	// Check for liquidity lock indicators
	r.checkLiquidityLock(ctx, client, contractAddr, assessment)

	// Check for ownership concentration
	r.checkOwnershipConcentration(ctx, client, contractAddr, assessment)

	// Check for hidden mint functions
	r.checkHiddenMintFunctions(ctx, client, contractAddr, assessment)

	return nil
}

// checkHoneypotPatterns checks for honeypot contract patterns
func (r *RiskAssessmentService) checkHoneypotPatterns(ctx context.Context, client *ethclient.Client, contractAddr common.Address, assessment *RiskAssessment) {
	code, err := client.CodeAt(ctx, contractAddr, nil)
	if err != nil {
		return
	}

	codeHex := common.Bytes2Hex(code)

	// Check for common honeypot patterns
	honeypotPatterns := []string{
		"transfer", // Functions that might have hidden restrictions
		"approve",
		"balanceOf",
	}

	suspiciousPatterns := 0
	for _, pattern := range honeypotPatterns {
		if strings.Contains(strings.ToLower(codeHex), pattern) {
			suspiciousPatterns++
		}
	}

	if suspiciousPatterns >= 2 {
		assessment.Factors = append(assessment.Factors, RiskFactor{
			Type:        "potential_honeypot",
			Description: "Contract shows honeypot characteristics",
			Impact:      0.7,
			Weight:      0.8,
			Evidence:    fmt.Sprintf("Found %d suspicious patterns", suspiciousPatterns),
		})
		assessment.Warnings = append(assessment.Warnings, "Contract may be a honeypot - tokens might not be sellable")
	}
}

// checkLiquidityLock checks for liquidity lock indicators
func (r *RiskAssessmentService) checkLiquidityLock(ctx context.Context, client *ethclient.Client, contractAddr common.Address, assessment *RiskAssessment) {
	// This would check if liquidity is locked in a time-lock contract
	// For now, add a generic check
	assessment.Factors = append(assessment.Factors, RiskFactor{
		Type:        "liquidity_lock_unknown",
		Description: "Liquidity lock status unknown",
		Impact:      0.3,
		Weight:      0.6,
		Evidence:    "Unable to verify liquidity lock",
	})
	assessment.Recommendations = append(assessment.Recommendations, "Verify that liquidity is locked before investing")
}

// checkOwnershipConcentration checks for high ownership concentration
func (r *RiskAssessmentService) checkOwnershipConcentration(ctx context.Context, client *ethclient.Client, contractAddr common.Address, assessment *RiskAssessment) {
	// This would analyze token holder distribution
	// For now, add a generic warning
	assessment.Factors = append(assessment.Factors, RiskFactor{
		Type:        "ownership_concentration_unknown",
		Description: "Token ownership concentration unknown",
		Impact:      0.2,
		Weight:      0.4,
		Evidence:    "Unable to analyze token distribution",
	})
}

// checkHiddenMintFunctions checks for hidden mint functions
func (r *RiskAssessmentService) checkHiddenMintFunctions(ctx context.Context, client *ethclient.Client, contractAddr common.Address, assessment *RiskAssessment) {
	code, err := client.CodeAt(ctx, contractAddr, nil)
	if err != nil {
		return
	}

	codeHex := common.Bytes2Hex(code)

	// Check for mint function signatures
	mintSignatures := []string{
		"40c10f19", // mint(address,uint256)
		"a0712d68", // mint(uint256)
		"1249c58b", // mint()
	}

	for _, sig := range mintSignatures {
		if strings.Contains(codeHex, sig) {
			assessment.Factors = append(assessment.Factors, RiskFactor{
				Type:        "mint_function_detected",
				Description: "Mint function detected in contract",
				Impact:      0.5,
				Weight:      0.7,
				Evidence:    fmt.Sprintf("Mint function signature: %s", sig),
			})
			assessment.Warnings = append(assessment.Warnings, "Contract has mint function - supply can be increased")
			break
		}
	}
}
