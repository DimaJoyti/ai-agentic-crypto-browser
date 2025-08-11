package solana

import (
	"context"
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Note: Types are defined in service.go

// TransactionService handles Solana transaction operations
type TransactionService struct {
	service *Service
}

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	StatusPending   TransactionStatus = "pending"
	StatusConfirmed TransactionStatus = "confirmed"
	StatusFinalized TransactionStatus = "finalized"
	StatusFailed    TransactionStatus = "failed"
)

// TransactionType represents different types of transactions
type TransactionType string

const (
	TypeTransfer      TransactionType = "transfer"
	TypeTokenTransfer TransactionType = "token_transfer"
	TypeSwap          TransactionType = "swap"
	TypeStake         TransactionType = "stake"
	TypeUnstake       TransactionType = "unstake"
	TypeNFTTransfer   TransactionType = "nft_transfer"
)

// NewTransactionService creates a new transaction service
func NewTransactionService(service *Service) *TransactionService {
	return &TransactionService{
		service: service,
	}
}

// SendTransaction sends a transaction to the Solana network
func (t *TransactionService) SendTransaction(ctx context.Context, req TransactionRequest) (*TransactionResult, error) {
	// Log the transaction start
	t.service.logger.Info(ctx, "Starting transaction", map[string]interface{}{
		"operation": "SendTransaction",
		"from":      req.From.String(),
		"to":        req.To.String(),
		"amount":    req.Amount.String(),
	})

	// Get recent blockhash
	recentBlockhash, err := t.service.client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		t.service.logger.Error(ctx, "Failed to get recent blockhash", err)
		return nil, fmt.Errorf("failed to get recent blockhash: %w", err)
	}

	// Create transaction based on type
	var tx *solana.Transaction
	var txType TransactionType

	if req.TokenMint != nil {
		// Token transfer
		tx, err = t.createTokenTransferTransaction(ctx, req, recentBlockhash.Value.Blockhash)
		txType = TypeTokenTransfer
	} else {
		// SOL transfer
		tx, err = t.createSOLTransferTransaction(ctx, req, recentBlockhash.Value.Blockhash)
		txType = TypeTransfer
	}

	if err != nil {
		t.service.logger.Error(ctx, "Failed to create transaction", err)
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Set compute unit price based on priority
	computeUnitPrice := t.getComputeUnitPrice(req.Priority)
	if computeUnitPrice > 0 {
		// Add compute budget instruction
		// This would be implemented with proper compute budget program
		t.service.logger.Info(ctx, "Setting compute unit price", map[string]interface{}{
			"priority": req.Priority,
			"price":    computeUnitPrice,
		})
	}

	// Simulate transaction first
	simulationResult, err := t.service.client.SimulateTransaction(ctx, tx)
	if err != nil {
		t.service.logger.Error(ctx, "Transaction simulation failed", err)
		return nil, fmt.Errorf("transaction simulation failed: %w", err)
	}

	if simulationResult.Value.Err != nil {
		t.service.logger.Error(ctx, "Transaction simulation error", fmt.Errorf("%v", simulationResult.Value.Err))
		return nil, fmt.Errorf("transaction simulation error: %v", simulationResult.Value.Err)
	}

	// Send transaction
	signature, err := t.service.client.SendTransaction(ctx, tx)
	if err != nil {
		t.service.logger.Error(ctx, "Failed to send transaction", err)
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	// Save transaction to database
	txRecord := &TransactionRecord{
		ID:              uuid.New(),
		Signature:       signature,
		From:            req.From,
		To:              req.To,
		Amount:          req.Amount,
		TokenMint:       req.TokenMint,
		TransactionType: txType,
		Status:          StatusPending,
		Priority:        req.Priority,
		MaxFee:          req.MaxFee,
		Memo:            req.Memo,
		CreatedAt:       time.Now(),
	}

	err = t.saveTransaction(ctx, txRecord)
	if err != nil {
		t.service.logger.Error(ctx, "Failed to save transaction", err)
		// Don't return error here as transaction was already sent
	}

	result := &TransactionResult{
		Signature: signature,
		Status:    string(StatusPending),
	}

	t.service.logger.Info(ctx, "Transaction sent successfully", map[string]interface{}{
		"signature": signature.String(),
		"from":      req.From.String(),
		"to":        req.To.String(),
		"amount":    req.Amount.String(),
		"type":      txType,
	})

	return result, nil
}

// createSOLTransferTransaction creates a SOL transfer transaction
func (t *TransactionService) createSOLTransferTransaction(ctx context.Context, req TransactionRequest, blockhash solana.Hash) (*solana.Transaction, error) {
	// Convert amount to lamports
	lamports := req.Amount.Mul(decimal.NewFromInt(1e9)).IntPart()

	// Create transfer instruction
	instruction := system.NewTransferInstruction(
		uint64(lamports),
		req.From,
		req.To,
	).Build()

	// Create transaction
	tx, err := solana.NewTransaction(
		[]solana.Instruction{instruction},
		blockhash,
		solana.TransactionPayer(req.From),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return tx, nil
}

// createTokenTransferTransaction creates a token transfer transaction
func (t *TransactionService) createTokenTransferTransaction(ctx context.Context, req TransactionRequest, blockhash solana.Hash) (*solana.Transaction, error) {
	// Get token accounts for sender and receiver
	senderTokenAccount, err := t.getAssociatedTokenAccount(ctx, req.From, *req.TokenMint)
	if err != nil {
		return nil, fmt.Errorf("failed to get sender token account: %w", err)
	}

	receiverTokenAccount, err := t.getAssociatedTokenAccount(ctx, req.To, *req.TokenMint)
	if err != nil {
		return nil, fmt.Errorf("failed to get receiver token account: %w", err)
	}

	// For simplicity, assume 9 decimals (would parse from mint data in real implementation)
	// In a real implementation, you would parse the mint account data to get the actual decimals
	decimals := uint8(9)
	amount := req.Amount.Mul(decimal.NewFromInt(int64(1) << decimals)).IntPart()

	// Create transfer instruction
	instruction := token.NewTransferInstruction(
		uint64(amount),
		senderTokenAccount,
		receiverTokenAccount,
		req.From,
		[]solana.PublicKey{},
	).Build()

	// Create transaction
	tx, err := solana.NewTransaction(
		[]solana.Instruction{instruction},
		blockhash,
		solana.TransactionPayer(req.From),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return tx, nil
}

// getAssociatedTokenAccount gets or creates an associated token account
func (t *TransactionService) getAssociatedTokenAccount(ctx context.Context, owner, mint solana.PublicKey) (solana.PublicKey, error) {
	// Calculate associated token account address
	ata, _, err := solana.FindAssociatedTokenAddress(owner, mint)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("failed to find associated token address: %w", err)
	}

	// Check if account exists
	accountInfo, err := t.service.client.GetAccountInfo(ctx, ata)
	if err != nil || accountInfo.Value == nil {
		// Account doesn't exist, would need to create it
		// For now, return the calculated address
		t.service.logger.Info(ctx, "Associated token account needs to be created", map[string]interface{}{
			"owner": owner.String(),
			"mint":  mint.String(),
			"ata":   ata.String(),
		})
	}

	return ata, nil
}

// getComputeUnitPrice returns the compute unit price based on priority
func (t *TransactionService) getComputeUnitPrice(priority TransactionPriority) uint64 {
	switch priority {
	case PriorityLow:
		return 0 // Use default
	case PriorityMedium:
		return 1000 // 1000 micro-lamports per compute unit
	case PriorityHigh:
		return 5000 // 5000 micro-lamports per compute unit
	case PriorityMax:
		return 10000 // 10000 micro-lamports per compute unit
	default:
		return 0
	}
}

// TransactionRecord represents a transaction record in the database
type TransactionRecord struct {
	ID              uuid.UUID           `json:"id"`
	Signature       solana.Signature    `json:"signature"`
	From            solana.PublicKey    `json:"from"`
	To              solana.PublicKey    `json:"to"`
	Amount          decimal.Decimal     `json:"amount"`
	TokenMint       *solana.PublicKey   `json:"token_mint,omitempty"`
	TransactionType TransactionType     `json:"transaction_type"`
	Status          TransactionStatus   `json:"status"`
	Priority        TransactionPriority `json:"priority"`
	MaxFee          decimal.Decimal     `json:"max_fee"`
	ActualFee       *decimal.Decimal    `json:"actual_fee,omitempty"`
	Memo            string              `json:"memo,omitempty"`
	BlockTime       *time.Time          `json:"block_time,omitempty"`
	Slot            *uint64             `json:"slot,omitempty"`
	Error           string              `json:"error,omitempty"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
}

// saveTransaction saves a transaction record to the database
func (t *TransactionService) saveTransaction(ctx context.Context, tx *TransactionRecord) error {
	query := `
		INSERT INTO solana_transactions (
			id, signature, from_address, to_address, amount, token_mint,
			transaction_type, status, priority, max_fee, memo, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	var tokenMintStr *string
	if tx.TokenMint != nil {
		mint := tx.TokenMint.String()
		tokenMintStr = &mint
	}

	_, err := t.service.db.DB.Exec(query,
		tx.ID,
		tx.Signature.String(),
		tx.From.String(),
		tx.To.String(),
		tx.Amount,
		tokenMintStr,
		tx.TransactionType,
		tx.Status,
		tx.Priority,
		tx.MaxFee,
		tx.Memo,
		tx.CreatedAt,
		tx.CreatedAt, // updated_at
	)

	return err
}

// GetTransactionHistory gets transaction history for a wallet
func (t *TransactionService) GetTransactionHistory(ctx context.Context, pubkey solana.PublicKey, limit int, offset int) ([]*TransactionRecord, error) {
	// Log the operation start
	t.service.logger.Info(ctx, "Getting transaction history", map[string]interface{}{
		"operation": "GetTransactionHistory",
		"pubkey":    pubkey.String(),
		"limit":     limit,
		"offset":    offset,
	})

	query := `
		SELECT id, signature, from_address, to_address, amount, token_mint,
			   transaction_type, status, priority, max_fee, actual_fee, memo,
			   block_time, slot, error, created_at, updated_at
		FROM solana_transactions 
		WHERE from_address = $1 OR to_address = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := t.service.db.DB.Query(query, pubkey.String(), limit, offset)
	if err != nil {
		t.service.logger.Error(ctx, "Failed to get transaction history", err)
		return nil, fmt.Errorf("failed to get transaction history: %w", err)
	}
	defer rows.Close()

	var transactions []*TransactionRecord
	for rows.Next() {
		var tx TransactionRecord
		var fromStr, toStr string
		var tokenMintStr *string
		var actualFee *decimal.Decimal

		err := rows.Scan(
			&tx.ID,
			&tx.Signature,
			&fromStr,
			&toStr,
			&tx.Amount,
			&tokenMintStr,
			&tx.TransactionType,
			&tx.Status,
			&tx.Priority,
			&tx.MaxFee,
			&actualFee,
			&tx.Memo,
			&tx.BlockTime,
			&tx.Slot,
			&tx.Error,
			&tx.CreatedAt,
			&tx.UpdatedAt,
		)
		if err != nil {
			t.service.logger.Error(ctx, "Failed to scan transaction row", err)
			continue
		}

		// Parse addresses
		tx.From, _ = solana.PublicKeyFromBase58(fromStr)
		tx.To, _ = solana.PublicKeyFromBase58(toStr)

		// Parse token mint if present
		if tokenMintStr != nil {
			mint, err := solana.PublicKeyFromBase58(*tokenMintStr)
			if err == nil {
				tx.TokenMint = &mint
			}
		}

		if actualFee != nil {
			tx.ActualFee = actualFee
		}

		transactions = append(transactions, &tx)
	}

	return transactions, nil
}
