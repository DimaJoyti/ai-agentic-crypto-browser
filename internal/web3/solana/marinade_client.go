package solana

import (
	"context"
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// MarinadeClient handles interactions with Marinade Finance (liquid staking)
type MarinadeClient struct {
	service *Service
}

// MarinadeStakeRequest represents a liquid staking request
type MarinadeStakeRequest struct {
	Amount        decimal.Decimal  `json:"amount"`
	UserPublicKey solana.PublicKey `json:"user_public_key"`
}

// MarinadeUnstakeRequest represents an unstaking request
type MarinadeUnstakeRequest struct {
	Amount        decimal.Decimal  `json:"amount"`
	UserPublicKey solana.PublicKey `json:"user_public_key"`
	IsDelayed     bool             `json:"is_delayed"` // true for delayed unstake, false for immediate
}

// MarinadeStakeResult represents the result of a staking operation
type MarinadeStakeResult struct {
	Signature    solana.Signature `json:"signature"`
	SOLAmount    decimal.Decimal  `json:"sol_amount"`
	MSOLAmount   decimal.Decimal  `json:"msol_amount"`
	ExchangeRate decimal.Decimal  `json:"exchange_rate"`
	Success      bool             `json:"success"`
	Error        string           `json:"error,omitempty"`
}

// MarinadeUnstakeResult represents the result of an unstaking operation
type MarinadeUnstakeResult struct {
	Signature    solana.Signature `json:"signature"`
	MSOLAmount   decimal.Decimal  `json:"msol_amount"`
	SOLAmount    decimal.Decimal  `json:"sol_amount"`
	ExchangeRate decimal.Decimal  `json:"exchange_rate"`
	IsDelayed    bool             `json:"is_delayed"`
	UnlockEpoch  *uint64          `json:"unlock_epoch,omitempty"`
	Success      bool             `json:"success"`
	Error        string           `json:"error,omitempty"`
}

// MarinadeState represents the current state of Marinade protocol
type MarinadeState struct {
	TotalSOLStaked       decimal.Decimal `json:"total_sol_staked"`
	TotalMSOLSupply      decimal.Decimal `json:"total_msol_supply"`
	ExchangeRate         decimal.Decimal `json:"exchange_rate"`
	APY                  decimal.Decimal `json:"apy"`
	ValidatorCount       int             `json:"validator_count"`
	LiquidityPoolBalance decimal.Decimal `json:"liquidity_pool_balance"`
	LastUpdateEpoch      uint64          `json:"last_update_epoch"`
}

// MarinadeValidator represents a validator in the Marinade network
type MarinadeValidator struct {
	VoteAccount   solana.PublicKey `json:"vote_account"`
	ValidatorName string           `json:"validator_name"`
	Commission    decimal.Decimal  `json:"commission"`
	APY           decimal.Decimal  `json:"apy"`
	StakedAmount  decimal.Decimal  `json:"staked_amount"`
	Score         decimal.Decimal  `json:"score"`
	IsActive      bool             `json:"is_active"`
	IsDelinquent  bool             `json:"is_delinquent"`
}

// Well-known Marinade addresses
var (
	MarinadeStateAddress = solana.MustPublicKeyFromBase58("8szGkuLTAux9XMgZ2vtY39jVSowEcpBfFfD8hXSEqdGC")
	MSOLMintAddress      = solana.MustPublicKeyFromBase58("mSoLzYCxHdYgdzU16g5QSh3i5K3z3KZK7ytfqcJm7So")
)

// NewMarinadeClient creates a new Marinade client
func NewMarinadeClient(service *Service) *MarinadeClient {
	return &MarinadeClient{
		service: service,
	}
}

// StakeSOL stakes SOL and receives mSOL (liquid staking tokens)
func (m *MarinadeClient) StakeSOL(ctx context.Context, req MarinadeStakeRequest) (*MarinadeStakeResult, error) {
	// Log the operation start
	m.service.logger.Info(ctx, "Staking SOL with Marinade", map[string]interface{}{
		"operation": "StakeSOL",
		"amount":    req.Amount.String(),
		"wallet":    req.UserPublicKey.String(),
	})

	// Get current Marinade state
	state, err := m.GetMarinadeState(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Marinade state: %w", err)
	}

	// Calculate mSOL amount based on exchange rate
	msolAmount := req.Amount.Div(state.ExchangeRate)

	// Create stake instruction
	instruction, err := m.createStakeInstruction(ctx, req, msolAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to create stake instruction: %w", err)
	}

	// Execute transaction
	programReq := ProgramInteractionRequest{
		ProgramID:   MarinadeProgramID,
		Instruction: "deposit",
		Accounts:    instruction.Accounts,
		Data:        instruction.Data,
		Signer:      req.UserPublicKey,
	}

	result, err := m.service.programMgr.InteractWithProgram(ctx, programReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute stake: %w", err)
	}

	if !result.Success {
		return &MarinadeStakeResult{
			Success: false,
			Error:   result.Error,
		}, nil
	}

	// Save staking position
	err = m.saveStakingPosition(ctx, req.UserPublicKey, req.Amount, msolAmount, state.ExchangeRate)
	if err != nil {
		m.service.logger.Error(ctx, "Failed to save staking position", err)
		// Don't return error as the transaction was successful
	}

	stakeResult := &MarinadeStakeResult{
		Signature:    result.Signature,
		SOLAmount:    req.Amount,
		MSOLAmount:   msolAmount,
		ExchangeRate: state.ExchangeRate,
		Success:      true,
	}

	m.service.logger.Info(ctx, "SOL staked with Marinade", map[string]interface{}{
		"signature":     result.Signature.String(),
		"sol_amount":    req.Amount.String(),
		"msol_amount":   msolAmount.String(),
		"exchange_rate": state.ExchangeRate.String(),
	})

	return stakeResult, nil
}

// UnstakeSOL unstakes mSOL to receive SOL
func (m *MarinadeClient) UnstakeSOL(ctx context.Context, req MarinadeUnstakeRequest) (*MarinadeUnstakeResult, error) {
	// Log the operation start
	m.service.logger.Info(ctx, "Unstaking mSOL with Marinade", map[string]interface{}{
		"operation": "UnstakeSOL",
		"amount":    req.Amount.String(),
	})

	// Get current Marinade state
	state, err := m.GetMarinadeState(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Marinade state: %w", err)
	}

	// Calculate SOL amount based on exchange rate
	solAmount := req.Amount.Mul(state.ExchangeRate)

	var instruction *SwapInstruction
	var unlockEpoch *uint64

	if req.IsDelayed {
		// Delayed unstake - better exchange rate but need to wait
		instruction, err = m.createDelayedUnstakeInstruction(ctx, req, solAmount)
		if err != nil {
			return nil, fmt.Errorf("failed to create delayed unstake instruction: %w", err)
		}
		// Calculate unlock epoch (current epoch + unstake delay)
		currentEpoch := uint64(time.Now().Unix() / (24 * 60 * 60 * 2)) // Simplified epoch calculation
		epoch := currentEpoch + 2                                      // 2 epochs delay
		unlockEpoch = &epoch
	} else {
		// Immediate unstake - uses liquidity pool with potential fee
		instruction, err = m.createImmediateUnstakeInstruction(ctx, req, solAmount)
		if err != nil {
			return nil, fmt.Errorf("failed to create immediate unstake instruction: %w", err)
		}
		// Apply liquidity pool fee (typically 0.3-0.9%)
		fee := solAmount.Mul(decimal.NewFromFloat(0.005)) // 0.5% fee
		solAmount = solAmount.Sub(fee)
	}

	// Execute transaction
	programReq := ProgramInteractionRequest{
		ProgramID:   MarinadeProgramID,
		Instruction: "withdraw",
		Accounts:    instruction.Accounts,
		Data:        instruction.Data,
		Signer:      req.UserPublicKey,
	}

	result, err := m.service.programMgr.InteractWithProgram(ctx, programReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute unstake: %w", err)
	}

	if !result.Success {
		return &MarinadeUnstakeResult{
			Success: false,
			Error:   result.Error,
		}, nil
	}

	unstakeResult := &MarinadeUnstakeResult{
		Signature:    result.Signature,
		MSOLAmount:   req.Amount,
		SOLAmount:    solAmount,
		ExchangeRate: state.ExchangeRate,
		IsDelayed:    req.IsDelayed,
		UnlockEpoch:  unlockEpoch,
		Success:      true,
	}

	m.service.logger.Info(ctx, "mSOL unstaked with Marinade", map[string]interface{}{
		"signature":    result.Signature.String(),
		"msol_amount":  req.Amount.String(),
		"sol_amount":   solAmount.String(),
		"is_delayed":   req.IsDelayed,
		"unlock_epoch": unlockEpoch,
	})

	return unstakeResult, nil
}

// GetMarinadeState gets the current state of Marinade protocol
func (m *MarinadeClient) GetMarinadeState(ctx context.Context) (*MarinadeState, error) {
	// Log the operation start
	m.service.logger.Info(ctx, "Getting Marinade state", map[string]interface{}{
		"operation": "GetMarinadeState",
	})

	// In a real implementation, this would fetch the actual state from the blockchain
	// For now, return simulated state
	state := &MarinadeState{
		TotalSOLStaked:       decimal.NewFromInt(6500000),    // 6.5M SOL
		TotalMSOLSupply:      decimal.NewFromInt(6200000),    // 6.2M mSOL
		ExchangeRate:         decimal.NewFromFloat(1.048387), // 1 mSOL = 1.048387 SOL
		APY:                  decimal.NewFromFloat(0.074),    // 7.4% APY
		ValidatorCount:       450,
		LiquidityPoolBalance: decimal.NewFromInt(50000), // 50K SOL in liquidity pool
		LastUpdateEpoch:      uint64(time.Now().Unix() / (24 * 60 * 60 * 2)),
	}

	m.service.logger.Info(ctx, "Retrieved Marinade state", map[string]interface{}{
		"total_sol_staked": state.TotalSOLStaked.String(),
		"exchange_rate":    state.ExchangeRate.String(),
		"apy":              state.APY.String(),
		"validator_count":  state.ValidatorCount,
	})

	return state, nil
}

// GetValidators gets the list of validators in the Marinade network
func (m *MarinadeClient) GetValidators(ctx context.Context) ([]*MarinadeValidator, error) {
	// Log the operation start
	m.service.logger.Info(ctx, "Getting Marinade validators", map[string]interface{}{
		"operation": "GetValidators",
	})

	// In a real implementation, this would fetch validators from the blockchain
	validators := []*MarinadeValidator{
		{
			VoteAccount:   solana.MustPublicKeyFromBase58("7K8DVxtNJGnMtUY1CQJT5jcs8sFGSZTDiG7kowvFpECh"),
			ValidatorName: "Marinade Validator 1",
			Commission:    decimal.NewFromFloat(0.05),  // 5%
			APY:           decimal.NewFromFloat(0.075), // 7.5%
			StakedAmount:  decimal.NewFromInt(100000),  // 100K SOL
			Score:         decimal.NewFromFloat(0.95),  // 95% score
			IsActive:      true,
			IsDelinquent:  false,
		},
		{
			VoteAccount:   solana.MustPublicKeyFromBase58("9QU2QSxhb24FUX3Tu2FpczXjpK3VYrvRudywSZaM29mF"),
			ValidatorName: "Marinade Validator 2",
			Commission:    decimal.NewFromFloat(0.07),  // 7%
			APY:           decimal.NewFromFloat(0.073), // 7.3%
			StakedAmount:  decimal.NewFromInt(85000),   // 85K SOL
			Score:         decimal.NewFromFloat(0.92),  // 92% score
			IsActive:      true,
			IsDelinquent:  false,
		},
	}

	return validators, nil
}

// GetTVL gets Marinade's total value locked
func (m *MarinadeClient) GetTVL(ctx context.Context) (decimal.Decimal, error) {
	state, err := m.GetMarinadeState(ctx)
	if err != nil {
		return decimal.Zero, err
	}

	// TVL is the total SOL staked multiplied by SOL price
	// Assuming SOL price of $200 for calculation
	solPrice := decimal.NewFromInt(200)
	tvl := state.TotalSOLStaked.Mul(solPrice)

	return tvl, nil
}

// Helper methods

func (m *MarinadeClient) createStakeInstruction(ctx context.Context, req MarinadeStakeRequest, msolAmount decimal.Decimal) (*SwapInstruction, error) {
	return &SwapInstruction{
		Accounts: []AccountMeta{
			{PublicKey: req.UserPublicKey, IsSigner: true, IsWritable: false},
			{PublicKey: MarinadeStateAddress, IsSigner: false, IsWritable: true},
			{PublicKey: MSOLMintAddress, IsSigner: false, IsWritable: true},
		},
		Data: []byte{}, // Instruction data would be encoded here
	}, nil
}

func (m *MarinadeClient) createDelayedUnstakeInstruction(ctx context.Context, req MarinadeUnstakeRequest, solAmount decimal.Decimal) (*SwapInstruction, error) {
	return &SwapInstruction{
		Accounts: []AccountMeta{
			{PublicKey: req.UserPublicKey, IsSigner: true, IsWritable: false},
			{PublicKey: MarinadeStateAddress, IsSigner: false, IsWritable: true},
			{PublicKey: MSOLMintAddress, IsSigner: false, IsWritable: true},
		},
		Data: []byte{}, // Instruction data for delayed unstake
	}, nil
}

func (m *MarinadeClient) createImmediateUnstakeInstruction(ctx context.Context, req MarinadeUnstakeRequest, solAmount decimal.Decimal) (*SwapInstruction, error) {
	return &SwapInstruction{
		Accounts: []AccountMeta{
			{PublicKey: req.UserPublicKey, IsSigner: true, IsWritable: false},
			{PublicKey: MarinadeStateAddress, IsSigner: false, IsWritable: true},
			{PublicKey: MSOLMintAddress, IsSigner: false, IsWritable: true},
		},
		Data: []byte{}, // Instruction data for immediate unstake
	}, nil
}

func (m *MarinadeClient) saveStakingPosition(ctx context.Context, userPubkey solana.PublicKey, solAmount, msolAmount, exchangeRate decimal.Decimal) error {
	// Find wallet ID
	query := `SELECT id FROM solana_wallets WHERE public_key = $1 AND is_active = true LIMIT 1`
	var walletID uuid.UUID
	err := m.service.db.DB.QueryRow(query, userPubkey.String()).Scan(&walletID)
	if err != nil {
		return fmt.Errorf("failed to find wallet: %w", err)
	}

	// Insert staking position
	insertQuery := `
		INSERT INTO solana_staking_positions (
			wallet_id, stake_account, validator_address, validator_name,
			staked_amount, rewards_earned, apy, status, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err = m.service.db.DB.Exec(insertQuery,
		walletID,
		"marinade_liquid_stake", // Placeholder for stake account
		MarinadeProgramID.String(),
		"Marinade Finance",
		solAmount,
		decimal.Zero,                // No rewards yet
		decimal.NewFromFloat(0.074), // 7.4% APY
		"active",
		time.Now(),
	)

	return err
}
