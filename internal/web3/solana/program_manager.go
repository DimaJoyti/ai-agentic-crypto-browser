package solana

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// ProgramManager handles Solana program (smart contract) interactions
type ProgramManager struct {
	service *Service
}

// ProgramInfo represents information about a Solana program
type ProgramInfo struct {
	ID            solana.PublicKey  `json:"id"`
	Name          string            `json:"name"`
	Version       string            `json:"version"`
	Description   string            `json:"description"`
	IsUpgradeable bool              `json:"is_upgradeable"`
	Authority     *solana.PublicKey `json:"authority,omitempty"`
}

// ProgramInteractionRequest represents a request to interact with a program
type ProgramInteractionRequest struct {
	ProgramID    solana.PublicKey `json:"program_id"`
	Instruction  string           `json:"instruction"`
	Accounts     []AccountMeta    `json:"accounts"`
	Data         []byte           `json:"data"`
	Signer       solana.PublicKey `json:"signer"`
	ComputeLimit *uint32          `json:"compute_limit,omitempty"`
}

// AccountMeta represents account metadata for program instructions
type AccountMeta struct {
	PublicKey  solana.PublicKey `json:"public_key"`
	IsSigner   bool             `json:"is_signer"`
	IsWritable bool             `json:"is_writable"`
}

// ProgramInteractionResult represents the result of a program interaction
type ProgramInteractionResult struct {
	Signature   solana.Signature `json:"signature"`
	Success     bool             `json:"success"`
	Error       string           `json:"error,omitempty"`
	Logs        []string         `json:"logs,omitempty"`
	ComputeUsed *uint64          `json:"compute_used,omitempty"`
}

// Well-known Solana program IDs
var (
	SystemProgramID               = solana.SystemProgramID
	TokenProgramID                = solana.TokenProgramID
	AssociatedTokenProgramID      = solana.SPLAssociatedTokenAccountProgramID
	StakeProgramID                = solana.StakeProgramID
	VoteProgramID                 = solana.VoteProgramID
	BPFLoaderProgramID            = solana.BPFLoaderProgramID
	BPFLoaderUpgradeableProgramID = solana.BPFLoaderUpgradeableProgramID
)

// Common DeFi program IDs (these would be the actual addresses in production)
var (
	// Jupiter Aggregator
	JupiterProgramID = solana.MustPublicKeyFromBase58("JUP6LkbZbjS1jKKwapdHNy74zcZ3tLUZoi5QNyVTaV4")

	// Raydium AMM
	RaydiumAMMProgramID = solana.MustPublicKeyFromBase58("675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8")

	// Orca
	OrcaProgramID = solana.MustPublicKeyFromBase58("9W959DqEETiGZocYWCQPaJ6sBmUzgfxXfqGeTEdp3aQP")

	// Marinade Finance
	MarinadeProgramID = solana.MustPublicKeyFromBase58("MarBmsSgKXdrN1egZf5sqe1TMai9K1rChYNDJgjq7aD")

	// Metaplex Token Metadata
	MetaplexTokenMetadataProgramID = solana.MustPublicKeyFromBase58("metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s")
)

// NewProgramManager creates a new program manager
func NewProgramManager(service *Service) *ProgramManager {
	return &ProgramManager{
		service: service,
	}
}

// GetProgramInfo retrieves information about a Solana program
func (p *ProgramManager) GetProgramInfo(ctx context.Context, programID solana.PublicKey) (*ProgramInfo, error) {
	// Log the operation start
	p.service.logger.Info(ctx, "Getting program info", map[string]interface{}{
		"operation":  "GetProgramInfo",
		"program_id": programID.String(),
	})

	// Get program account info
	accountInfo, err := p.service.client.GetAccountInfo(ctx, programID)
	if err != nil {
		p.service.logger.Error(ctx, "Failed to get program account info", err)
		return nil, fmt.Errorf("failed to get program account info: %w", err)
	}

	if accountInfo.Value == nil {
		return nil, fmt.Errorf("program not found: %s", programID.String())
	}

	info := &ProgramInfo{
		ID:            programID,
		Name:          p.getProgramName(programID),
		IsUpgradeable: accountInfo.Value.Owner.Equals(BPFLoaderUpgradeableProgramID),
	}

	// If it's an upgradeable program, get the program data account
	if info.IsUpgradeable {
		programDataAddress, err := p.getProgramDataAddress(programID)
		if err == nil {
			programDataInfo, err := p.service.client.GetAccountInfo(ctx, programDataAddress)
			if err == nil && programDataInfo.Value != nil {
				// Parse program data to get authority (simplified)
				// In a real implementation, you would properly parse the program data
				p.service.logger.Info(ctx, "Found upgradeable program data", map[string]interface{}{
					"program_id":   programID.String(),
					"program_data": programDataAddress.String(),
					"data_length":  len(programDataInfo.Value.Data.GetBinary()),
				})
			}
		}
	}

	p.service.logger.Info(ctx, "Retrieved program info", map[string]interface{}{
		"program_id":     programID.String(),
		"name":           info.Name,
		"is_upgradeable": info.IsUpgradeable,
	})

	return info, nil
}

// InteractWithProgram executes an instruction on a Solana program
func (p *ProgramManager) InteractWithProgram(ctx context.Context, req ProgramInteractionRequest) (*ProgramInteractionResult, error) {
	// Log the operation start
	p.service.logger.Info(ctx, "Interacting with program", map[string]interface{}{
		"operation":   "InteractWithProgram",
		"program_id":  req.ProgramID.String(),
		"instruction": req.Instruction,
	})

	// Convert AccountMeta to solana-go format
	var accountSlice solana.AccountMetaSlice
	for _, acc := range req.Accounts {
		accountSlice = append(accountSlice, &solana.AccountMeta{
			PublicKey:  acc.PublicKey,
			IsSigner:   acc.IsSigner,
			IsWritable: acc.IsWritable,
		})
	}

	// Create instruction
	instruction := solana.NewInstruction(
		req.ProgramID,
		accountSlice,
		req.Data,
	)

	// Get recent blockhash
	recentBlockhash, err := p.service.client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		p.service.logger.Error(ctx, "Failed to get recent blockhash", err)
		return nil, fmt.Errorf("failed to get recent blockhash: %w", err)
	}

	// Create transaction
	tx, err := solana.NewTransaction(
		[]solana.Instruction{instruction},
		recentBlockhash.Value.Blockhash,
		solana.TransactionPayer(req.Signer),
	)
	if err != nil {
		p.service.logger.Error(ctx, "Failed to create transaction", err)
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Simulate transaction first
	simulationResult, err := p.service.client.SimulateTransaction(ctx, tx)
	if err != nil {
		p.service.logger.Error(ctx, "Program interaction simulation failed", err)
		return &ProgramInteractionResult{
			Success: false,
			Error:   fmt.Sprintf("simulation failed: %v", err),
		}, nil
	}

	if simulationResult.Value.Err != nil {
		p.service.logger.Error(ctx, "Program interaction simulation error", fmt.Errorf("%v", simulationResult.Value.Err))
		return &ProgramInteractionResult{
			Success: false,
			Error:   fmt.Sprintf("simulation error: %v", simulationResult.Value.Err),
			Logs:    simulationResult.Value.Logs,
		}, nil
	}

	// Send transaction
	signature, err := p.service.client.SendTransaction(ctx, tx)
	if err != nil {
		p.service.logger.Error(ctx, "Failed to send program interaction", err)
		return &ProgramInteractionResult{
			Success: false,
			Error:   fmt.Sprintf("failed to send transaction: %v", err),
		}, nil
	}

	result := &ProgramInteractionResult{
		Signature: signature,
		Success:   true,
		Logs:      simulationResult.Value.Logs,
	}

	if simulationResult.Value.UnitsConsumed != nil {
		result.ComputeUsed = simulationResult.Value.UnitsConsumed
	}

	p.service.logger.Info(ctx, "Program interaction successful", map[string]interface{}{
		"program_id":   req.ProgramID.String(),
		"instruction":  req.Instruction,
		"signature":    signature.String(),
		"compute_used": result.ComputeUsed,
	})

	return result, nil
}

// GetTokenMetadata retrieves metadata for a token using Metaplex
func (p *ProgramManager) GetTokenMetadata(ctx context.Context, mintAddress solana.PublicKey) (*TokenMetadata, error) {
	// Log the operation start
	p.service.logger.Info(ctx, "Getting token metadata", map[string]interface{}{
		"operation":    "GetTokenMetadata",
		"mint_address": mintAddress.String(),
	})

	// Calculate metadata account address
	metadataAddress, err := p.getMetadataAddress(mintAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate metadata address: %w", err)
	}

	// Get metadata account
	accountInfo, err := p.service.client.GetAccountInfo(ctx, metadataAddress)
	if err != nil {
		p.service.logger.Error(ctx, "Failed to get metadata account", err)
		return nil, fmt.Errorf("failed to get metadata account: %w", err)
	}

	if accountInfo.Value == nil {
		return nil, fmt.Errorf("metadata not found for mint: %s", mintAddress.String())
	}

	// Parse metadata (simplified - in reality you'd use proper Metaplex parsing)
	metadata := &TokenMetadata{
		Mint:        mintAddress,
		Name:        "Unknown Token",
		Symbol:      "UNK",
		Description: "Token metadata parsing not fully implemented",
		Image:       "",
		Attributes:  []TokenAttribute{},
	}

	p.service.logger.Info(ctx, "Retrieved token metadata", map[string]interface{}{
		"mint":     mintAddress.String(),
		"metadata": metadataAddress.String(),
		"name":     metadata.Name,
	})

	return metadata, nil
}

// TokenMetadata represents token metadata
type TokenMetadata struct {
	Mint        solana.PublicKey  `json:"mint"`
	Name        string            `json:"name"`
	Symbol      string            `json:"symbol"`
	Description string            `json:"description"`
	Image       string            `json:"image"`
	Attributes  []TokenAttribute  `json:"attributes"`
	Collection  *solana.PublicKey `json:"collection,omitempty"`
	Creators    []Creator         `json:"creators,omitempty"`
}

// TokenAttribute represents a token attribute
type TokenAttribute struct {
	TraitType string `json:"trait_type"`
	Value     string `json:"value"`
}

// Creator represents a token creator
type Creator struct {
	Address  solana.PublicKey `json:"address"`
	Verified bool             `json:"verified"`
	Share    uint8            `json:"share"`
}

// Helper methods

func (p *ProgramManager) getProgramName(programID solana.PublicKey) string {
	switch programID {
	case SystemProgramID:
		return "System Program"
	case TokenProgramID:
		return "Token Program"
	case AssociatedTokenProgramID:
		return "Associated Token Program"
	case StakeProgramID:
		return "Stake Program"
	case VoteProgramID:
		return "Vote Program"
	case JupiterProgramID:
		return "Jupiter Aggregator"
	case RaydiumAMMProgramID:
		return "Raydium AMM"
	case OrcaProgramID:
		return "Orca"
	case MarinadeProgramID:
		return "Marinade Finance"
	case MetaplexTokenMetadataProgramID:
		return "Metaplex Token Metadata"
	default:
		return "Unknown Program"
	}
}

func (p *ProgramManager) getProgramDataAddress(programID solana.PublicKey) (solana.PublicKey, error) {
	// Calculate program data address for upgradeable programs
	seeds := [][]byte{programID.Bytes()}
	programDataAddress, _, err := solana.FindProgramAddress(seeds, BPFLoaderUpgradeableProgramID)
	return programDataAddress, err
}

func (p *ProgramManager) getMetadataAddress(mintAddress solana.PublicKey) (solana.PublicKey, error) {
	// Calculate metadata account address using Metaplex standard
	seeds := [][]byte{
		[]byte("metadata"),
		MetaplexTokenMetadataProgramID.Bytes(),
		mintAddress.Bytes(),
	}
	metadataAddress, _, err := solana.FindProgramAddress(seeds, MetaplexTokenMetadataProgramID)
	return metadataAddress, err
}
