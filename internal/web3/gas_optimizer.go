package web3

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"
)

// GasOptimizer provides advanced gas optimization strategies
type GasOptimizer struct {
	clients map[int]*ethclient.Client
	logger  *observability.Logger
}

// GasEstimate represents a gas estimation with optimization
type GasEstimate struct {
	GasLimit             uint64        `json:"gas_limit"`
	GasPrice             *big.Int      `json:"gas_price"`
	MaxFeePerGas         *big.Int      `json:"max_fee_per_gas"`
	MaxPriorityFeePerGas *big.Int      `json:"max_priority_fee_per_gas"`
	EstimatedCost        *big.Int      `json:"estimated_cost"`
	Strategy             string        `json:"strategy"`
	Confidence           float64       `json:"confidence"`
	TimeToConfirm        time.Duration `json:"time_to_confirm"`
}

// GasStrategy represents different gas optimization strategies
type GasStrategy string

const (
	GasStrategyEconomical GasStrategy = "economical"
	GasStrategyStandard   GasStrategy = "standard"
	GasStrategyFast       GasStrategy = "fast"
	GasStrategyInstant    GasStrategy = "instant"
)

// NewGasOptimizer creates a new gas optimizer
func NewGasOptimizer(clients map[int]*ethclient.Client, logger *observability.Logger) *GasOptimizer {
	return &GasOptimizer{
		clients: clients,
		logger:  logger,
	}
}

// EstimateGas provides optimized gas estimation for a transaction
func (g *GasOptimizer) EstimateGas(ctx context.Context, chainID int, tx ethereum.CallMsg, strategy GasStrategy) (*GasEstimate, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("gas-optimizer").Start(ctx, "gas.EstimateGas")
	defer span.End()

	client, exists := g.clients[chainID]
	if !exists {
		return nil, fmt.Errorf("no client available for chain ID: %d", chainID)
	}

	// Get current gas price and network conditions
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	// Estimate gas limit
	gasLimit, err := client.EstimateGas(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas: %w", err)
	}

	// Check if EIP-1559 is supported
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest header: %w", err)
	}

	var estimate *GasEstimate
	if header.BaseFee != nil {
		// EIP-1559 transaction
		estimate, err = g.estimateEIP1559Gas(ctx, client, gasLimit, gasPrice, strategy)
	} else {
		// Legacy transaction
		estimate, err = g.estimateLegacyGas(ctx, gasLimit, gasPrice, strategy)
	}

	if err != nil {
		return nil, err
	}

	// Add buffer for gas limit (10% safety margin)
	estimate.GasLimit = uint64(float64(estimate.GasLimit) * 1.1)

	g.logger.Info(ctx, "Gas estimation completed", map[string]interface{}{
		"chain_id":       chainID,
		"strategy":       string(strategy),
		"gas_limit":      estimate.GasLimit,
		"estimated_cost": estimate.EstimatedCost.String(),
	})

	return estimate, nil
}

// estimateEIP1559Gas estimates gas for EIP-1559 transactions
func (g *GasOptimizer) estimateEIP1559Gas(ctx context.Context, client *ethclient.Client, gasLimit uint64, baseFee *big.Int, strategy GasStrategy) (*GasEstimate, error) {
	// Get fee history for better estimation
	feeHistory, err := client.FeeHistory(ctx, 20, nil, []float64{25, 50, 75})
	if err != nil {
		return nil, fmt.Errorf("failed to get fee history: %w", err)
	}

	// Calculate priority fee based on strategy
	var priorityFeeMultiplier float64
	var timeToConfirm time.Duration
	var confidence float64

	switch strategy {
	case GasStrategyEconomical:
		priorityFeeMultiplier = 1.0
		timeToConfirm = 5 * time.Minute
		confidence = 0.7
	case GasStrategyStandard:
		priorityFeeMultiplier = 1.2
		timeToConfirm = 2 * time.Minute
		confidence = 0.85
	case GasStrategyFast:
		priorityFeeMultiplier = 1.5
		timeToConfirm = 30 * time.Second
		confidence = 0.95
	case GasStrategyInstant:
		priorityFeeMultiplier = 2.0
		timeToConfirm = 15 * time.Second
		confidence = 0.99
	default:
		priorityFeeMultiplier = 1.2
		timeToConfirm = 2 * time.Minute
		confidence = 0.85
	}

	// Calculate average priority fee from fee history
	avgPriorityFee := g.calculateAveragePriorityFee(feeHistory)

	// Apply strategy multiplier
	priorityFee := new(big.Int).Mul(avgPriorityFee, big.NewInt(int64(priorityFeeMultiplier*100)))
	priorityFee.Div(priorityFee, big.NewInt(100))

	// Calculate max fee per gas (base fee + priority fee + buffer)
	maxFeePerGas := new(big.Int).Add(baseFee, priorityFee)
	maxFeePerGas.Mul(maxFeePerGas, big.NewInt(120)) // 20% buffer
	maxFeePerGas.Div(maxFeePerGas, big.NewInt(100))

	// Calculate estimated cost
	estimatedCost := new(big.Int).Mul(maxFeePerGas, big.NewInt(int64(gasLimit)))

	return &GasEstimate{
		GasLimit:             gasLimit,
		MaxFeePerGas:         maxFeePerGas,
		MaxPriorityFeePerGas: priorityFee,
		EstimatedCost:        estimatedCost,
		Strategy:             string(strategy),
		Confidence:           confidence,
		TimeToConfirm:        timeToConfirm,
	}, nil
}

// estimateLegacyGas estimates gas for legacy transactions
func (g *GasOptimizer) estimateLegacyGas(ctx context.Context, gasLimit uint64, gasPrice *big.Int, strategy GasStrategy) (*GasEstimate, error) {
	var multiplier float64
	var timeToConfirm time.Duration
	var confidence float64

	switch strategy {
	case GasStrategyEconomical:
		multiplier = 1.0
		timeToConfirm = 5 * time.Minute
		confidence = 0.7
	case GasStrategyStandard:
		multiplier = 1.2
		timeToConfirm = 2 * time.Minute
		confidence = 0.85
	case GasStrategyFast:
		multiplier = 1.5
		timeToConfirm = 30 * time.Second
		confidence = 0.95
	case GasStrategyInstant:
		multiplier = 2.0
		timeToConfirm = 15 * time.Second
		confidence = 0.99
	default:
		multiplier = 1.2
		timeToConfirm = 2 * time.Minute
		confidence = 0.85
	}

	// Apply strategy multiplier to gas price
	optimizedGasPrice := new(big.Int).Mul(gasPrice, big.NewInt(int64(multiplier*100)))
	optimizedGasPrice.Div(optimizedGasPrice, big.NewInt(100))

	// Calculate estimated cost
	estimatedCost := new(big.Int).Mul(optimizedGasPrice, big.NewInt(int64(gasLimit)))

	return &GasEstimate{
		GasLimit:      gasLimit,
		GasPrice:      optimizedGasPrice,
		EstimatedCost: estimatedCost,
		Strategy:      string(strategy),
		Confidence:    confidence,
		TimeToConfirm: timeToConfirm,
	}, nil
}

// calculateAveragePriorityFee calculates average priority fee from fee history
func (g *GasOptimizer) calculateAveragePriorityFee(feeHistory *ethereum.FeeHistory) *big.Int {
	if len(feeHistory.Reward) == 0 {
		return big.NewInt(2000000000) // 2 gwei default
	}

	total := big.NewInt(0)
	count := 0

	for _, rewards := range feeHistory.Reward {
		if len(rewards) > 1 { // Use 50th percentile
			total.Add(total, rewards[1])
			count++
		}
	}

	if count == 0 {
		return big.NewInt(2000000000) // 2 gwei default
	}

	return total.Div(total, big.NewInt(int64(count)))
}

// OptimizeTransactionBatch optimizes gas for multiple transactions
func (g *GasOptimizer) OptimizeTransactionBatch(ctx context.Context, chainID int, txs []ethereum.CallMsg, strategy GasStrategy) ([]*GasEstimate, error) {
	estimates := make([]*GasEstimate, len(txs))

	for i, tx := range txs {
		estimate, err := g.EstimateGas(ctx, chainID, tx, strategy)
		if err != nil {
			return nil, fmt.Errorf("failed to estimate gas for transaction %d: %w", i, err)
		}
		estimates[i] = estimate
	}

	return estimates, nil
}

// GetNetworkCongestion returns current network congestion level
func (g *GasOptimizer) GetNetworkCongestion(ctx context.Context, chainID int) (float64, error) {
	client, exists := g.clients[chainID]
	if !exists {
		return 0, fmt.Errorf("no client available for chain ID: %d", chainID)
	}

	// Get pending transaction count
	pendingCount, err := client.PendingTransactionCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get pending transaction count: %w", err)
	}

	// Simple congestion calculation (can be enhanced with more sophisticated metrics)
	// Higher pending count = higher congestion
	congestion := float64(pendingCount) / 1000.0 // Normalize to 0-1 scale
	if congestion > 1.0 {
		congestion = 1.0
	}

	return congestion, nil
}
