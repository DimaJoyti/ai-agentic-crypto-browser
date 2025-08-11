package solana

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/shopspring/decimal"
)

// RaydiumClient handles interactions with Raydium AMM
type RaydiumClient struct {
	service *Service
}

// RaydiumPool represents a Raydium liquidity pool
type RaydiumPool struct {
	ID           solana.PublicKey `json:"id"`
	BaseMint     solana.PublicKey `json:"base_mint"`
	QuoteMint    solana.PublicKey `json:"quote_mint"`
	BaseVault    solana.PublicKey `json:"base_vault"`
	QuoteVault   solana.PublicKey `json:"quote_vault"`
	LPMint       solana.PublicKey `json:"lp_mint"`
	BaseReserve  decimal.Decimal  `json:"base_reserve"`
	QuoteReserve decimal.Decimal  `json:"quote_reserve"`
	LPSupply     decimal.Decimal  `json:"lp_supply"`
	Fee          decimal.Decimal  `json:"fee"`
	APY          decimal.Decimal  `json:"apy"`
	Volume24h    decimal.Decimal  `json:"volume_24h"`
	TVL          decimal.Decimal  `json:"tvl"`
	IsActive     bool             `json:"is_active"`
}

// Well-known Raydium pool addresses (examples)
var (
	RaydiumSOLUSDCPool = solana.MustPublicKeyFromBase58("58oQChx4yWmvKdwLLZzBi4ChoCc2fqCUWBkwMihLYQo2")
	RaydiumRAYSOLPool  = solana.MustPublicKeyFromBase58("AVs9TA4nWDzfPJE9gGVNJMVhcQy3V9PGazuz33BfG2RA")
)

// NewRaydiumClient creates a new Raydium client
func NewRaydiumClient(service *Service) *RaydiumClient {
	return &RaydiumClient{
		service: service,
	}
}

// GetSwapQuote gets a swap quote from Raydium
func (r *RaydiumClient) GetSwapQuote(ctx context.Context, req SwapRequest) (*SwapResult, error) {
	// Log the operation start
	r.service.logger.Info(ctx, "Getting Raydium swap quote", map[string]interface{}{
		"operation":   "GetSwapQuote",
		"input_mint":  req.InputMint.String(),
		"output_mint": req.OutputMint.String(),
		"amount":      req.Amount.String(),
	})

	// Find the best pool for this token pair
	pool, err := r.findBestPool(ctx, req.InputMint, req.OutputMint)
	if err != nil {
		return nil, fmt.Errorf("failed to find pool: %w", err)
	}

	// Calculate swap amounts using constant product formula (x * y = k)
	outputAmount, priceImpact, fee := r.calculateSwapOutput(req.Amount, pool, req.InputMint)

	// Check slippage
	slippageAmount := outputAmount.Mul(decimal.NewFromFloat(float64(req.SlippageBps) / 10000))
	minOutputAmount := outputAmount.Sub(slippageAmount)

	route := []SwapRoute{
		{
			Protocol:     ProtocolRaydium,
			InputMint:    req.InputMint,
			OutputMint:   req.OutputMint,
			InputAmount:  req.Amount,
			OutputAmount: outputAmount,
			Fee:          fee,
		},
	}

	result := &SwapResult{
		InputAmount:  req.Amount,
		OutputAmount: minOutputAmount,
		PriceImpact:  priceImpact,
		Fee:          fee,
		Route:        route,
		Success:      true,
	}

	r.service.logger.Info(ctx, "Raydium quote calculated", map[string]interface{}{
		"pool_id":       pool.ID.String(),
		"input_amount":  req.Amount.String(),
		"output_amount": outputAmount.String(),
		"price_impact":  priceImpact.String(),
		"fee":           fee.String(),
	})

	return result, nil
}

// ExecuteSwap executes a swap on Raydium
func (r *RaydiumClient) ExecuteSwap(ctx context.Context, req SwapRequest) (*SwapResult, error) {
	// Log the operation start
	r.service.logger.Info(ctx, "Executing Raydium swap", map[string]interface{}{
		"operation": "ExecuteSwap",
		"amount":    req.Amount.String(),
	})

	// Get quote first
	quote, err := r.GetSwapQuote(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get quote: %w", err)
	}

	// Find the pool
	pool, err := r.findBestPool(ctx, req.InputMint, req.OutputMint)
	if err != nil {
		return nil, fmt.Errorf("failed to find pool: %w", err)
	}

	// Create swap instruction
	instruction, err := r.createSwapInstruction(ctx, req, pool)
	if err != nil {
		return nil, fmt.Errorf("failed to create swap instruction: %w", err)
	}

	// Execute transaction through program manager
	programReq := ProgramInteractionRequest{
		ProgramID:   RaydiumAMMProgramID,
		Instruction: "swap",
		Accounts:    instruction.Accounts,
		Data:        instruction.Data,
		Signer:      req.UserPublicKey,
	}

	result, err := r.service.programMgr.InteractWithProgram(ctx, programReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute swap: %w", err)
	}

	if !result.Success {
		return &SwapResult{
			Success: false,
			Error:   result.Error,
		}, nil
	}

	// Return successful result
	swapResult := &SwapResult{
		Signature:    result.Signature,
		InputAmount:  quote.InputAmount,
		OutputAmount: quote.OutputAmount,
		PriceImpact:  quote.PriceImpact,
		Fee:          quote.Fee,
		Route:        quote.Route,
		Success:      true,
	}

	r.service.logger.Info(ctx, "Raydium swap executed", map[string]interface{}{
		"signature":     result.Signature.String(),
		"pool_id":       pool.ID.String(),
		"input_amount":  swapResult.InputAmount.String(),
		"output_amount": swapResult.OutputAmount.String(),
	})

	return swapResult, nil
}

// AddLiquidity adds liquidity to a Raydium pool
func (r *RaydiumClient) AddLiquidity(ctx context.Context, req LiquidityRequest) (*LiquidityResult, error) {
	// Log the operation start
	r.service.logger.Info(ctx, "Adding liquidity to Raydium pool", map[string]interface{}{
		"operation": "AddLiquidity",
		"pool":      req.PoolAddress.String(),
		"amount_a":  req.AmountA.String(),
		"amount_b":  req.AmountB.String(),
	})

	// Get pool information
	pool, err := r.getPoolInfo(ctx, req.PoolAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get pool info: %w", err)
	}

	// Calculate optimal amounts and LP tokens
	optimalAmountA, optimalAmountB, lpTokens := r.calculateLiquidityAmounts(req.AmountA, req.AmountB, pool)

	// Create add liquidity instruction
	instruction, err := r.createAddLiquidityInstruction(ctx, req, pool, optimalAmountA, optimalAmountB)
	if err != nil {
		return nil, fmt.Errorf("failed to create add liquidity instruction: %w", err)
	}

	// Execute transaction
	programReq := ProgramInteractionRequest{
		ProgramID:   RaydiumAMMProgramID,
		Instruction: "addLiquidity",
		Accounts:    instruction.Accounts,
		Data:        instruction.Data,
		Signer:      req.UserPublicKey,
	}

	result, err := r.service.programMgr.InteractWithProgram(ctx, programReq)
	if err != nil {
		return nil, fmt.Errorf("failed to add liquidity: %w", err)
	}

	if !result.Success {
		return &LiquidityResult{
			Success: false,
			Error:   result.Error,
		}, nil
	}

	// Calculate pool share
	poolShare := lpTokens.Div(pool.LPSupply.Add(lpTokens)).Mul(decimal.NewFromInt(100))

	liquidityResult := &LiquidityResult{
		Signature:   result.Signature,
		LPTokens:    lpTokens,
		AmountAUsed: optimalAmountA,
		AmountBUsed: optimalAmountB,
		PoolShare:   poolShare,
		Success:     true,
	}

	r.service.logger.Info(ctx, "Raydium liquidity added", map[string]interface{}{
		"signature":  result.Signature.String(),
		"pool_id":    pool.ID.String(),
		"lp_tokens":  lpTokens.String(),
		"pool_share": poolShare.String(),
	})

	return liquidityResult, nil
}

// GetTVL gets Raydium's total value locked
func (r *RaydiumClient) GetTVL(ctx context.Context) (decimal.Decimal, error) {
	// In a real implementation, this would aggregate TVL from all Raydium pools
	// For now, return a simulated value
	return decimal.NewFromInt(1500000000), nil // $1.5B simulated TVL
}

// GetPools gets all active Raydium pools
func (r *RaydiumClient) GetPools(ctx context.Context) ([]*RaydiumPool, error) {
	// Log the operation start
	r.service.logger.Info(ctx, "Getting Raydium pools", map[string]interface{}{
		"operation": "GetPools",
	})

	// In a real implementation, this would fetch all pools from the Raydium program
	// For now, return some example pools
	pools := []*RaydiumPool{
		{
			ID:           RaydiumSOLUSDCPool,
			BaseMint:     solana.SolMint,
			QuoteMint:    solana.MustPublicKeyFromBase58("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"), // USDC
			BaseReserve:  decimal.NewFromInt(100000),                                                     // 100K SOL
			QuoteReserve: decimal.NewFromInt(20000000),                                                   // 20M USDC
			LPSupply:     decimal.NewFromInt(1000000),                                                    // 1M LP tokens
			Fee:          decimal.NewFromFloat(0.0025),                                                   // 0.25%
			APY:          decimal.NewFromFloat(0.15),                                                     // 15%
			Volume24h:    decimal.NewFromInt(50000000),                                                   // $50M
			TVL:          decimal.NewFromInt(40000000),                                                   // $40M
			IsActive:     true,
		},
		{
			ID:           RaydiumRAYSOLPool,
			BaseMint:     solana.MustPublicKeyFromBase58("4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R"), // RAY
			QuoteMint:    solana.SolMint,
			BaseReserve:  decimal.NewFromInt(5000000),  // 5M RAY
			QuoteReserve: decimal.NewFromInt(50000),    // 50K SOL
			LPSupply:     decimal.NewFromInt(500000),   // 500K LP tokens
			Fee:          decimal.NewFromFloat(0.0025), // 0.25%
			APY:          decimal.NewFromFloat(0.25),   // 25%
			Volume24h:    decimal.NewFromInt(10000000), // $10M
			TVL:          decimal.NewFromInt(15000000), // $15M
			IsActive:     true,
		},
	}

	return pools, nil
}

// Helper methods

func (r *RaydiumClient) findBestPool(ctx context.Context, inputMint, outputMint solana.PublicKey) (*RaydiumPool, error) {
	pools, err := r.GetPools(ctx)
	if err != nil {
		return nil, err
	}

	// Find pool that matches the token pair
	for _, pool := range pools {
		if (pool.BaseMint.Equals(inputMint) && pool.QuoteMint.Equals(outputMint)) ||
			(pool.BaseMint.Equals(outputMint) && pool.QuoteMint.Equals(inputMint)) {
			return pool, nil
		}
	}

	return nil, fmt.Errorf("no pool found for token pair")
}

func (r *RaydiumClient) calculateSwapOutput(inputAmount decimal.Decimal, pool *RaydiumPool, inputMint solana.PublicKey) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	var inputReserve, outputReserve decimal.Decimal

	// Determine which reserve is input and which is output
	if pool.BaseMint.Equals(inputMint) {
		inputReserve = pool.BaseReserve
		outputReserve = pool.QuoteReserve
	} else {
		inputReserve = pool.QuoteReserve
		outputReserve = pool.BaseReserve
	}

	// Calculate fee
	fee := inputAmount.Mul(pool.Fee)
	inputAmountAfterFee := inputAmount.Sub(fee)

	// Constant product formula: (x + Δx) * (y - Δy) = x * y
	// Δy = (y * Δx) / (x + Δx)
	numerator := outputReserve.Mul(inputAmountAfterFee)
	denominator := inputReserve.Add(inputAmountAfterFee)
	outputAmount := numerator.Div(denominator)

	// Calculate price impact
	priceImpact := inputAmountAfterFee.Div(inputReserve.Add(inputAmountAfterFee)).Mul(decimal.NewFromInt(100))

	return outputAmount, priceImpact, fee
}

func (r *RaydiumClient) getPoolInfo(ctx context.Context, poolAddress solana.PublicKey) (*RaydiumPool, error) {
	// In a real implementation, this would fetch pool data from the blockchain
	// For now, return a mock pool
	return &RaydiumPool{
		ID:           poolAddress,
		BaseMint:     solana.SolMint,
		QuoteMint:    solana.MustPublicKeyFromBase58("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"),
		BaseReserve:  decimal.NewFromInt(100000),
		QuoteReserve: decimal.NewFromInt(20000000),
		LPSupply:     decimal.NewFromInt(1000000),
		Fee:          decimal.NewFromFloat(0.0025),
		IsActive:     true,
	}, nil
}

func (r *RaydiumClient) calculateLiquidityAmounts(amountA, amountB decimal.Decimal, pool *RaydiumPool) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	// Calculate optimal ratio based on pool reserves
	ratio := pool.QuoteReserve.Div(pool.BaseReserve)

	// Adjust amounts to maintain pool ratio
	optimalAmountB := amountA.Mul(ratio)
	if optimalAmountB.LessThanOrEqual(amountB) {
		// Use amountA and calculated amountB
		lpTokens := amountA.Mul(pool.LPSupply).Div(pool.BaseReserve)
		return amountA, optimalAmountB, lpTokens
	} else {
		// Use amountB and calculate amountA
		optimalAmountA := amountB.Div(ratio)
		lpTokens := optimalAmountA.Mul(pool.LPSupply).Div(pool.BaseReserve)
		return optimalAmountA, amountB, lpTokens
	}
}

func (r *RaydiumClient) createSwapInstruction(ctx context.Context, req SwapRequest, pool *RaydiumPool) (*SwapInstruction, error) {
	// This would create the actual Raydium swap instruction
	// For now, return a mock instruction structure
	return &SwapInstruction{
		Accounts: []AccountMeta{
			{PublicKey: req.UserPublicKey, IsSigner: true, IsWritable: false},
			{PublicKey: pool.ID, IsSigner: false, IsWritable: true},
			// Add other required accounts...
		},
		Data: []byte{}, // Instruction data would be encoded here
	}, nil
}

func (r *RaydiumClient) createAddLiquidityInstruction(ctx context.Context, req LiquidityRequest, pool *RaydiumPool, amountA, amountB decimal.Decimal) (*SwapInstruction, error) {
	// This would create the actual Raydium add liquidity instruction
	return &SwapInstruction{
		Accounts: []AccountMeta{
			{PublicKey: req.UserPublicKey, IsSigner: true, IsWritable: false},
			{PublicKey: pool.ID, IsSigner: false, IsWritable: true},
			// Add other required accounts...
		},
		Data: []byte{}, // Instruction data would be encoded here
	}, nil
}
