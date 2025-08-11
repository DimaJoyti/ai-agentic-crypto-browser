package solana

import (
	"context"
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// DeFiService handles Solana DeFi protocol interactions
type DeFiService struct {
	service        *Service
	jupiterClient  *JupiterClient
	raydiumClient  *RaydiumClient
	orcaClient     *OrcaClient
	marinadeClient *MarinadeClient
}

// DeFiProtocol represents supported DeFi protocols
type DeFiProtocol string

const (
	ProtocolJupiter  DeFiProtocol = "jupiter"
	ProtocolRaydium  DeFiProtocol = "raydium"
	ProtocolOrca     DeFiProtocol = "orca"
	ProtocolMarinade DeFiProtocol = "marinade"
	ProtocolKamino   DeFiProtocol = "kamino"
	ProtocolDrift    DeFiProtocol = "drift"
)

// DeFiPosition represents a DeFi position
type DeFiPosition struct {
	ID            uuid.UUID       `json:"id"`
	WalletID      uuid.UUID       `json:"wallet_id"`
	Protocol      DeFiProtocol    `json:"protocol"`
	PositionType  PositionType    `json:"position_type"`
	PoolAddress   string          `json:"pool_address"`
	TokenA        string          `json:"token_a"`
	TokenB        string          `json:"token_b"`
	AmountA       decimal.Decimal `json:"amount_a"`
	AmountB       decimal.Decimal `json:"amount_b"`
	LPTokens      decimal.Decimal `json:"lp_tokens"`
	APY           decimal.Decimal `json:"apy"`
	EntryPrice    decimal.Decimal `json:"entry_price"`
	CurrentValue  decimal.Decimal `json:"current_value"`
	UnrealizedPnL decimal.Decimal `json:"unrealized_pnl"`
	FeesEarned    decimal.Decimal `json:"fees_earned"`
	RewardsEarned decimal.Decimal `json:"rewards_earned"`
	IsActive      bool            `json:"is_active"`
	OpenedAt      time.Time       `json:"opened_at"`
	ClosedAt      *time.Time      `json:"closed_at,omitempty"`
}

// PositionType represents different types of DeFi positions
type PositionType string

const (
	PositionTypeLiquidity PositionType = "liquidity"
	PositionTypeStake     PositionType = "stake"
	PositionTypeLend      PositionType = "lend"
	PositionTypeBorrow    PositionType = "borrow"
	PositionTypeFarm      PositionType = "farm"
	PositionTypeVault     PositionType = "vault"
)

// SwapRequest represents a token swap request
type SwapRequest struct {
	InputMint     solana.PublicKey `json:"input_mint"`
	OutputMint    solana.PublicKey `json:"output_mint"`
	Amount        decimal.Decimal  `json:"amount"`
	SlippageBps   uint16           `json:"slippage_bps"` // basis points (100 = 1%)
	UserPublicKey solana.PublicKey `json:"user_public_key"`
	Protocol      DeFiProtocol     `json:"protocol,omitempty"`
}

// SwapResult represents the result of a token swap
type SwapResult struct {
	Signature    solana.Signature `json:"signature"`
	InputAmount  decimal.Decimal  `json:"input_amount"`
	OutputAmount decimal.Decimal  `json:"output_amount"`
	PriceImpact  decimal.Decimal  `json:"price_impact"`
	Fee          decimal.Decimal  `json:"fee"`
	Route        []SwapRoute      `json:"route"`
	Success      bool             `json:"success"`
	Error        string           `json:"error,omitempty"`
}

// SwapRoute represents a step in the swap route
type SwapRoute struct {
	Protocol     DeFiProtocol     `json:"protocol"`
	InputMint    solana.PublicKey `json:"input_mint"`
	OutputMint   solana.PublicKey `json:"output_mint"`
	InputAmount  decimal.Decimal  `json:"input_amount"`
	OutputAmount decimal.Decimal  `json:"output_amount"`
	Fee          decimal.Decimal  `json:"fee"`
}

// LiquidityRequest represents a liquidity provision request
type LiquidityRequest struct {
	PoolAddress   solana.PublicKey `json:"pool_address"`
	TokenAMint    solana.PublicKey `json:"token_a_mint"`
	TokenBMint    solana.PublicKey `json:"token_b_mint"`
	AmountA       decimal.Decimal  `json:"amount_a"`
	AmountB       decimal.Decimal  `json:"amount_b"`
	SlippageBps   uint16           `json:"slippage_bps"`
	UserPublicKey solana.PublicKey `json:"user_public_key"`
	Protocol      DeFiProtocol     `json:"protocol"`
}

// LiquidityResult represents the result of liquidity provision
type LiquidityResult struct {
	Signature   solana.Signature `json:"signature"`
	LPTokens    decimal.Decimal  `json:"lp_tokens"`
	AmountAUsed decimal.Decimal  `json:"amount_a_used"`
	AmountBUsed decimal.Decimal  `json:"amount_b_used"`
	PoolShare   decimal.Decimal  `json:"pool_share"`
	Success     bool             `json:"success"`
	Error       string           `json:"error,omitempty"`
}

// SwapInstruction represents a swap instruction for DeFi protocols
type SwapInstruction struct {
	Accounts []AccountMeta `json:"accounts"`
	Data     []byte        `json:"data"`
}

// NewDeFiService creates a new DeFi service
func NewDeFiService(service *Service) *DeFiService {
	return &DeFiService{
		service:        service,
		jupiterClient:  NewJupiterClient(service),
		raydiumClient:  NewRaydiumClient(service),
		orcaClient:     NewOrcaClient(service),
		marinadeClient: NewMarinadeClient(service),
	}
}

// GetBestSwapRoute finds the best swap route across all protocols
func (d *DeFiService) GetBestSwapRoute(ctx context.Context, req SwapRequest) (*SwapResult, error) {
	// Log the operation start
	d.service.logger.Info(ctx, "Getting best swap route", map[string]interface{}{
		"operation":   "GetBestSwapRoute",
		"input_mint":  req.InputMint.String(),
		"output_mint": req.OutputMint.String(),
		"amount":      req.Amount.String(),
	})

	// Get quotes from all protocols
	quotes := make(map[DeFiProtocol]*SwapResult)

	// Jupiter (aggregator - usually best)
	if jupiterQuote, err := d.jupiterClient.GetSwapQuote(ctx, req); err == nil {
		quotes[ProtocolJupiter] = jupiterQuote
	}

	// Raydium
	if raydiumQuote, err := d.raydiumClient.GetSwapQuote(ctx, req); err == nil {
		quotes[ProtocolRaydium] = raydiumQuote
	}

	// Orca
	if orcaQuote, err := d.orcaClient.GetSwapQuote(ctx, req); err == nil {
		quotes[ProtocolOrca] = orcaQuote
	}

	// Find best quote (highest output amount)
	var bestQuote *SwapResult
	var bestProtocol DeFiProtocol

	for protocol, quote := range quotes {
		if bestQuote == nil || quote.OutputAmount.GreaterThan(bestQuote.OutputAmount) {
			bestQuote = quote
			bestProtocol = protocol
		}
	}

	if bestQuote == nil {
		return nil, fmt.Errorf("no valid swap routes found")
	}

	d.service.logger.Info(ctx, "Found best swap route", map[string]interface{}{
		"protocol":      bestProtocol,
		"input_amount":  req.Amount.String(),
		"output_amount": bestQuote.OutputAmount.String(),
		"price_impact":  bestQuote.PriceImpact.String(),
	})

	return bestQuote, nil
}

// ExecuteSwap executes a token swap using the specified protocol
func (d *DeFiService) ExecuteSwap(ctx context.Context, req SwapRequest) (*SwapResult, error) {
	// Log the operation start
	d.service.logger.Info(ctx, "Executing swap", map[string]interface{}{
		"operation": "ExecuteSwap",
		"protocol":  req.Protocol,
		"amount":    req.Amount.String(),
	})

	var result *SwapResult
	var err error

	// If no protocol specified, find the best route
	if req.Protocol == "" {
		bestRoute, err := d.GetBestSwapRoute(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to find best route: %w", err)
		}
		// Use the protocol from the best route
		for _, route := range bestRoute.Route {
			req.Protocol = route.Protocol
			break
		}
	}

	// Execute swap on the specified protocol
	switch req.Protocol {
	case ProtocolJupiter:
		result, err = d.jupiterClient.ExecuteSwap(ctx, req)
	case ProtocolRaydium:
		result, err = d.raydiumClient.ExecuteSwap(ctx, req)
	case ProtocolOrca:
		result, err = d.orcaClient.ExecuteSwap(ctx, req)
	default:
		return nil, fmt.Errorf("unsupported protocol: %s", req.Protocol)
	}

	if err != nil {
		d.service.logger.Error(ctx, "Swap execution failed", err)
		return nil, fmt.Errorf("swap execution failed: %w", err)
	}

	d.service.logger.Info(ctx, "Swap executed successfully", map[string]interface{}{
		"protocol":      req.Protocol,
		"signature":     result.Signature.String(),
		"input_amount":  result.InputAmount.String(),
		"output_amount": result.OutputAmount.String(),
	})

	return result, nil
}

// AddLiquidity adds liquidity to a pool
func (d *DeFiService) AddLiquidity(ctx context.Context, req LiquidityRequest) (*LiquidityResult, error) {
	// Log the operation start
	d.service.logger.Info(ctx, "Adding liquidity", map[string]interface{}{
		"operation": "AddLiquidity",
		"protocol":  req.Protocol,
		"amount_a":  req.AmountA.String(),
		"amount_b":  req.AmountB.String(),
	})

	var result *LiquidityResult
	var err error

	switch req.Protocol {
	case ProtocolRaydium:
		result, err = d.raydiumClient.AddLiquidity(ctx, req)
	case ProtocolOrca:
		result, err = d.orcaClient.AddLiquidity(ctx, req)
	default:
		return nil, fmt.Errorf("unsupported protocol for liquidity: %s", req.Protocol)
	}

	if err != nil {
		d.service.logger.Error(ctx, "Add liquidity failed", err)
		return nil, fmt.Errorf("add liquidity failed: %w", err)
	}

	// Save position to database
	position := &DeFiPosition{
		ID:           uuid.New(),
		Protocol:     req.Protocol,
		PositionType: PositionTypeLiquidity,
		PoolAddress:  req.PoolAddress.String(),
		TokenA:       req.TokenAMint.String(),
		TokenB:       req.TokenBMint.String(),
		AmountA:      result.AmountAUsed,
		AmountB:      result.AmountBUsed,
		LPTokens:     result.LPTokens,
		IsActive:     true,
		OpenedAt:     time.Now(),
	}

	err = d.savePosition(ctx, position)
	if err != nil {
		d.service.logger.Error(ctx, "Failed to save liquidity position", err)
		// Don't return error as the transaction was successful
	}

	d.service.logger.Info(ctx, "Liquidity added successfully", map[string]interface{}{
		"protocol":   req.Protocol,
		"signature":  result.Signature.String(),
		"lp_tokens":  result.LPTokens.String(),
		"pool_share": result.PoolShare.String(),
	})

	return result, nil
}

// GetUserPositions retrieves all DeFi positions for a user
func (d *DeFiService) GetUserPositions(ctx context.Context, walletID uuid.UUID) ([]*DeFiPosition, error) {
	// Log the operation start
	d.service.logger.Info(ctx, "Getting user positions", map[string]interface{}{
		"operation": "GetUserPositions",
		"wallet_id": walletID.String(),
	})

	query := `
		SELECT id, wallet_id, protocol, position_type, pool_address, token_a, token_b,
			   amount_a, amount_b, lp_tokens, apy, entry_price, current_value,
			   unrealized_pnl, fees_earned, rewards_earned, is_active, opened_at, closed_at
		FROM solana_defi_positions 
		WHERE wallet_id = $1 
		ORDER BY opened_at DESC
	`

	rows, err := d.service.db.DB.Query(query, walletID)
	if err != nil {
		d.service.logger.Error(ctx, "Failed to get user positions", err)
		return nil, fmt.Errorf("failed to get user positions: %w", err)
	}
	defer rows.Close()

	var positions []*DeFiPosition
	for rows.Next() {
		var position DeFiPosition
		err := rows.Scan(
			&position.ID,
			&position.WalletID,
			&position.Protocol,
			&position.PositionType,
			&position.PoolAddress,
			&position.TokenA,
			&position.TokenB,
			&position.AmountA,
			&position.AmountB,
			&position.LPTokens,
			&position.APY,
			&position.EntryPrice,
			&position.CurrentValue,
			&position.UnrealizedPnL,
			&position.FeesEarned,
			&position.RewardsEarned,
			&position.IsActive,
			&position.OpenedAt,
			&position.ClosedAt,
		)
		if err != nil {
			d.service.logger.Error(ctx, "Failed to scan position row", err)
			continue
		}

		positions = append(positions, &position)
	}

	return positions, nil
}

// GetProtocolTVL gets the total value locked for a protocol
func (d *DeFiService) GetProtocolTVL(ctx context.Context, protocol DeFiProtocol) (decimal.Decimal, error) {
	// Log the operation start
	d.service.logger.Info(ctx, "Getting protocol TVL", map[string]interface{}{
		"operation": "GetProtocolTVL",
		"protocol":  protocol,
	})

	switch protocol {
	case ProtocolJupiter:
		return d.jupiterClient.GetTVL(ctx)
	case ProtocolRaydium:
		return d.raydiumClient.GetTVL(ctx)
	case ProtocolOrca:
		return d.orcaClient.GetTVL(ctx)
	case ProtocolMarinade:
		return d.marinadeClient.GetTVL(ctx)
	default:
		return decimal.Zero, fmt.Errorf("unsupported protocol: %s", protocol)
	}
}

// Helper methods

func (d *DeFiService) savePosition(ctx context.Context, position *DeFiPosition) error {
	query := `
		INSERT INTO solana_defi_positions (
			id, wallet_id, protocol, position_type, pool_address, token_a, token_b,
			amount_a, amount_b, lp_tokens, apy, entry_price, current_value,
			unrealized_pnl, fees_earned, rewards_earned, is_active, opened_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`

	_, err := d.service.db.DB.Exec(query,
		position.ID,
		position.WalletID,
		position.Protocol,
		position.PositionType,
		position.PoolAddress,
		position.TokenA,
		position.TokenB,
		position.AmountA,
		position.AmountB,
		position.LPTokens,
		position.APY,
		position.EntryPrice,
		position.CurrentValue,
		position.UnrealizedPnL,
		position.FeesEarned,
		position.RewardsEarned,
		position.IsActive,
		position.OpenedAt,
	)

	return err
}
