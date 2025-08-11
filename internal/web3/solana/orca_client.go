package solana

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/shopspring/decimal"
)

// OrcaClient handles interactions with Orca DEX
type OrcaClient struct {
	service *Service
}

// OrcaPool represents an Orca liquidity pool
type OrcaPool struct {
	ID             solana.PublicKey `json:"id"`
	TokenAMint     solana.PublicKey `json:"token_a_mint"`
	TokenBMint     solana.PublicKey `json:"token_b_mint"`
	TokenAVault    solana.PublicKey `json:"token_a_vault"`
	TokenBVault    solana.PublicKey `json:"token_b_vault"`
	LPMint         solana.PublicKey `json:"lp_mint"`
	TokenAAmount   decimal.Decimal  `json:"token_a_amount"`
	TokenBAmount   decimal.Decimal  `json:"token_b_amount"`
	LPSupply       decimal.Decimal  `json:"lp_supply"`
	FeeRate        decimal.Decimal  `json:"fee_rate"`
	APY            decimal.Decimal  `json:"apy"`
	Volume24h      decimal.Decimal  `json:"volume_24h"`
	TVL            decimal.Decimal  `json:"tvl"`
	IsStable       bool             `json:"is_stable"`
	IsConcentrated bool             `json:"is_concentrated"`
}

// Well-known Orca pool addresses
var (
	OrcaSOLUSDCPool = solana.MustPublicKeyFromBase58("EGZ7tiLeH62TPV1gL8WwbXGzEPa9zmcpVnnkPKKnrE2U")
	OrcaORCASOLPool = solana.MustPublicKeyFromBase58("2p7nYbtPBgtmY69NsE8DAW6szpRJn7tQvDnqvoEWQvjY")
)

// NewOrcaClient creates a new Orca client
func NewOrcaClient(service *Service) *OrcaClient {
	return &OrcaClient{
		service: service,
	}
}

// GetSwapQuote gets a swap quote from Orca
func (o *OrcaClient) GetSwapQuote(ctx context.Context, req SwapRequest) (*SwapResult, error) {
	// Log the operation start
	o.service.logger.Info(ctx, "Getting Orca swap quote", map[string]interface{}{
		"operation":   "GetSwapQuote",
		"input_mint":  req.InputMint.String(),
		"output_mint": req.OutputMint.String(),
		"amount":      req.Amount.String(),
	})

	// Find the best pool for this token pair
	pool, err := o.findBestPool(ctx, req.InputMint, req.OutputMint)
	if err != nil {
		return nil, fmt.Errorf("failed to find pool: %w", err)
	}

	// Calculate swap amounts
	outputAmount, priceImpact, fee := o.calculateSwapOutput(req.Amount, pool, req.InputMint)

	// Apply slippage
	slippageAmount := outputAmount.Mul(decimal.NewFromFloat(float64(req.SlippageBps) / 10000))
	minOutputAmount := outputAmount.Sub(slippageAmount)

	route := []SwapRoute{
		{
			Protocol:     ProtocolOrca,
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

	o.service.logger.Info(ctx, "Orca quote calculated", map[string]interface{}{
		"pool_id":         pool.ID.String(),
		"input_amount":    req.Amount.String(),
		"output_amount":   outputAmount.String(),
		"price_impact":    priceImpact.String(),
		"fee":             fee.String(),
		"is_stable":       pool.IsStable,
		"is_concentrated": pool.IsConcentrated,
	})

	return result, nil
}

// ExecuteSwap executes a swap on Orca
func (o *OrcaClient) ExecuteSwap(ctx context.Context, req SwapRequest) (*SwapResult, error) {
	// Log the operation start
	o.service.logger.Info(ctx, "Executing Orca swap", map[string]interface{}{
		"operation": "ExecuteSwap",
		"amount":    req.Amount.String(),
	})

	// Get quote first
	quote, err := o.GetSwapQuote(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get quote: %w", err)
	}

	// Find the pool
	pool, err := o.findBestPool(ctx, req.InputMint, req.OutputMint)
	if err != nil {
		return nil, fmt.Errorf("failed to find pool: %w", err)
	}

	// Create swap instruction based on pool type
	var instruction *SwapInstruction
	if pool.IsConcentrated {
		instruction, err = o.createConcentratedSwapInstruction(ctx, req, pool)
	} else {
		instruction, err = o.createStandardSwapInstruction(ctx, req, pool)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create swap instruction: %w", err)
	}

	// Execute transaction
	programReq := ProgramInteractionRequest{
		ProgramID:   OrcaProgramID,
		Instruction: "swap",
		Accounts:    instruction.Accounts,
		Data:        instruction.Data,
		Signer:      req.UserPublicKey,
	}

	result, err := o.service.programMgr.InteractWithProgram(ctx, programReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute swap: %w", err)
	}

	if !result.Success {
		return &SwapResult{
			Success: false,
			Error:   result.Error,
		}, nil
	}

	swapResult := &SwapResult{
		Signature:    result.Signature,
		InputAmount:  quote.InputAmount,
		OutputAmount: quote.OutputAmount,
		PriceImpact:  quote.PriceImpact,
		Fee:          quote.Fee,
		Route:        quote.Route,
		Success:      true,
	}

	o.service.logger.Info(ctx, "Orca swap executed", map[string]interface{}{
		"signature":     result.Signature.String(),
		"pool_id":       pool.ID.String(),
		"pool_type":     o.getPoolType(pool),
		"input_amount":  swapResult.InputAmount.String(),
		"output_amount": swapResult.OutputAmount.String(),
	})

	return swapResult, nil
}

// AddLiquidity adds liquidity to an Orca pool
func (o *OrcaClient) AddLiquidity(ctx context.Context, req LiquidityRequest) (*LiquidityResult, error) {
	// Log the operation start
	o.service.logger.Info(ctx, "Adding liquidity to Orca pool", map[string]interface{}{
		"operation": "AddLiquidity",
		"pool":      req.PoolAddress.String(),
		"amount_a":  req.AmountA.String(),
		"amount_b":  req.AmountB.String(),
	})

	// Get pool information
	pool, err := o.getPoolInfo(ctx, req.PoolAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get pool info: %w", err)
	}

	// Calculate optimal amounts and LP tokens
	var optimalAmountA, optimalAmountB, lpTokens decimal.Decimal

	if pool.IsConcentrated {
		// For concentrated liquidity pools, use different calculation
		optimalAmountA, optimalAmountB, lpTokens = o.calculateConcentratedLiquidityAmounts(req.AmountA, req.AmountB, pool)
	} else {
		// For standard pools, use constant product formula
		optimalAmountA, optimalAmountB, lpTokens = o.calculateStandardLiquidityAmounts(req.AmountA, req.AmountB, pool)
	}

	// Create add liquidity instruction
	var instruction *SwapInstruction
	if pool.IsConcentrated {
		instruction, err = o.createConcentratedLiquidityInstruction(ctx, req, pool, optimalAmountA, optimalAmountB)
	} else {
		instruction, err = o.createStandardLiquidityInstruction(ctx, req, pool, optimalAmountA, optimalAmountB)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create add liquidity instruction: %w", err)
	}

	// Execute transaction
	programReq := ProgramInteractionRequest{
		ProgramID:   OrcaProgramID,
		Instruction: "addLiquidity",
		Accounts:    instruction.Accounts,
		Data:        instruction.Data,
		Signer:      req.UserPublicKey,
	}

	result, err := o.service.programMgr.InteractWithProgram(ctx, programReq)
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

	o.service.logger.Info(ctx, "Orca liquidity added", map[string]interface{}{
		"signature":  result.Signature.String(),
		"pool_id":    pool.ID.String(),
		"pool_type":  o.getPoolType(pool),
		"lp_tokens":  lpTokens.String(),
		"pool_share": poolShare.String(),
	})

	return liquidityResult, nil
}

// GetTVL gets Orca's total value locked
func (o *OrcaClient) GetTVL(ctx context.Context) (decimal.Decimal, error) {
	// In a real implementation, this would aggregate TVL from all Orca pools
	return decimal.NewFromInt(1200000000), nil // $1.2B simulated TVL
}

// GetPools gets all active Orca pools
func (o *OrcaClient) GetPools(ctx context.Context) ([]*OrcaPool, error) {
	// Log the operation start
	o.service.logger.Info(ctx, "Getting Orca pools", map[string]interface{}{
		"operation": "GetPools",
	})

	// Example pools (in reality, these would be fetched from the blockchain)
	pools := []*OrcaPool{
		{
			ID:             OrcaSOLUSDCPool,
			TokenAMint:     solana.SolMint,
			TokenBMint:     solana.MustPublicKeyFromBase58("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"), // USDC
			TokenAAmount:   decimal.NewFromInt(80000),                                                      // 80K SOL
			TokenBAmount:   decimal.NewFromInt(16000000),                                                   // 16M USDC
			LPSupply:       decimal.NewFromInt(800000),                                                     // 800K LP tokens
			FeeRate:        decimal.NewFromFloat(0.003),                                                    // 0.3%
			APY:            decimal.NewFromFloat(0.18),                                                     // 18%
			Volume24h:      decimal.NewFromInt(40000000),                                                   // $40M
			TVL:            decimal.NewFromInt(32000000),                                                   // $32M
			IsStable:       false,
			IsConcentrated: true, // Orca Whirlpools are concentrated liquidity
		},
		{
			ID:             OrcaORCASOLPool,
			TokenAMint:     solana.MustPublicKeyFromBase58("orcaEKTdK7LKz57vaAYr9QeNsVEPfiu6QeMU1kektZE"), // ORCA
			TokenBMint:     solana.SolMint,
			TokenAAmount:   decimal.NewFromInt(2000000),  // 2M ORCA
			TokenBAmount:   decimal.NewFromInt(40000),    // 40K SOL
			LPSupply:       decimal.NewFromInt(400000),   // 400K LP tokens
			FeeRate:        decimal.NewFromFloat(0.003),  // 0.3%
			APY:            decimal.NewFromFloat(0.22),   // 22%
			Volume24h:      decimal.NewFromInt(8000000),  // $8M
			TVL:            decimal.NewFromInt(12000000), // $12M
			IsStable:       false,
			IsConcentrated: true,
		},
	}

	return pools, nil
}

// Helper methods

func (o *OrcaClient) findBestPool(ctx context.Context, inputMint, outputMint solana.PublicKey) (*OrcaPool, error) {
	pools, err := o.GetPools(ctx)
	if err != nil {
		return nil, err
	}

	// Find pool that matches the token pair
	for _, pool := range pools {
		if (pool.TokenAMint.Equals(inputMint) && pool.TokenBMint.Equals(outputMint)) ||
			(pool.TokenAMint.Equals(outputMint) && pool.TokenBMint.Equals(inputMint)) {
			return pool, nil
		}
	}

	return nil, fmt.Errorf("no pool found for token pair")
}

func (o *OrcaClient) calculateSwapOutput(inputAmount decimal.Decimal, pool *OrcaPool, inputMint solana.PublicKey) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	var inputReserve, outputReserve decimal.Decimal

	if pool.TokenAMint.Equals(inputMint) {
		inputReserve = pool.TokenAAmount
		outputReserve = pool.TokenBAmount
	} else {
		inputReserve = pool.TokenBAmount
		outputReserve = pool.TokenAAmount
	}

	// Calculate fee
	fee := inputAmount.Mul(pool.FeeRate)
	inputAmountAfterFee := inputAmount.Sub(fee)

	var outputAmount, priceImpact decimal.Decimal

	if pool.IsConcentrated {
		// For concentrated liquidity, use more complex calculation
		// Simplified version here
		outputAmount = o.calculateConcentratedSwapOutput(inputAmountAfterFee, inputReserve, outputReserve)
		priceImpact = inputAmountAfterFee.Div(inputReserve).Mul(decimal.NewFromInt(100))
	} else {
		// Standard constant product formula
		numerator := outputReserve.Mul(inputAmountAfterFee)
		denominator := inputReserve.Add(inputAmountAfterFee)
		outputAmount = numerator.Div(denominator)
		priceImpact = inputAmountAfterFee.Div(inputReserve.Add(inputAmountAfterFee)).Mul(decimal.NewFromInt(100))
	}

	return outputAmount, priceImpact, fee
}

func (o *OrcaClient) calculateConcentratedSwapOutput(inputAmount, inputReserve, outputReserve decimal.Decimal) decimal.Decimal {
	// Simplified concentrated liquidity calculation
	// In reality, this would consider tick ranges and liquidity distribution
	return inputAmount.Mul(outputReserve).Div(inputReserve.Add(inputAmount))
}

func (o *OrcaClient) calculateStandardLiquidityAmounts(amountA, amountB decimal.Decimal, pool *OrcaPool) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	// Standard constant product calculation
	ratio := pool.TokenBAmount.Div(pool.TokenAAmount)
	optimalAmountB := amountA.Mul(ratio)

	if optimalAmountB.LessThanOrEqual(amountB) {
		lpTokens := amountA.Mul(pool.LPSupply).Div(pool.TokenAAmount)
		return amountA, optimalAmountB, lpTokens
	} else {
		optimalAmountA := amountB.Div(ratio)
		lpTokens := optimalAmountA.Mul(pool.LPSupply).Div(pool.TokenAAmount)
		return optimalAmountA, amountB, lpTokens
	}
}

func (o *OrcaClient) calculateConcentratedLiquidityAmounts(amountA, amountB decimal.Decimal, pool *OrcaPool) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	// Concentrated liquidity calculation (simplified)
	// In reality, this would consider price ranges and tick spacing
	return o.calculateStandardLiquidityAmounts(amountA, amountB, pool)
}

func (o *OrcaClient) getPoolInfo(ctx context.Context, poolAddress solana.PublicKey) (*OrcaPool, error) {
	// Mock pool info - in reality, fetch from blockchain
	return &OrcaPool{
		ID:             poolAddress,
		TokenAMint:     solana.SolMint,
		TokenBMint:     solana.MustPublicKeyFromBase58("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"),
		TokenAAmount:   decimal.NewFromInt(80000),
		TokenBAmount:   decimal.NewFromInt(16000000),
		LPSupply:       decimal.NewFromInt(800000),
		FeeRate:        decimal.NewFromFloat(0.003),
		IsConcentrated: true,
	}, nil
}

func (o *OrcaClient) getPoolType(pool *OrcaPool) string {
	if pool.IsConcentrated {
		return "concentrated"
	} else if pool.IsStable {
		return "stable"
	}
	return "standard"
}

func (o *OrcaClient) createStandardSwapInstruction(ctx context.Context, req SwapRequest, pool *OrcaPool) (*SwapInstruction, error) {
	return &SwapInstruction{
		Accounts: []AccountMeta{
			{PublicKey: req.UserPublicKey, IsSigner: true, IsWritable: false},
			{PublicKey: pool.ID, IsSigner: false, IsWritable: true},
		},
		Data: []byte{},
	}, nil
}

func (o *OrcaClient) createConcentratedSwapInstruction(ctx context.Context, req SwapRequest, pool *OrcaPool) (*SwapInstruction, error) {
	return &SwapInstruction{
		Accounts: []AccountMeta{
			{PublicKey: req.UserPublicKey, IsSigner: true, IsWritable: false},
			{PublicKey: pool.ID, IsSigner: false, IsWritable: true},
		},
		Data: []byte{},
	}, nil
}

func (o *OrcaClient) createStandardLiquidityInstruction(ctx context.Context, req LiquidityRequest, pool *OrcaPool, amountA, amountB decimal.Decimal) (*SwapInstruction, error) {
	return &SwapInstruction{
		Accounts: []AccountMeta{
			{PublicKey: req.UserPublicKey, IsSigner: true, IsWritable: false},
			{PublicKey: pool.ID, IsSigner: false, IsWritable: true},
		},
		Data: []byte{},
	}, nil
}

func (o *OrcaClient) createConcentratedLiquidityInstruction(ctx context.Context, req LiquidityRequest, pool *OrcaPool, amountA, amountB decimal.Decimal) (*SwapInstruction, error) {
	return &SwapInstruction{
		Accounts: []AccountMeta{
			{PublicKey: req.UserPublicKey, IsSigner: true, IsWritable: false},
			{PublicKey: pool.ID, IsSigner: false, IsWritable: true},
		},
		Data: []byte{},
	}, nil
}
