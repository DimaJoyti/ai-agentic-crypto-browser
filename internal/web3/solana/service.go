package solana

import (
	"context"
	"fmt"
	"time"

	"github.com/ai-agentic-browser/pkg/database"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Service provides Solana blockchain functionality
type Service struct {
	client     *rpc.Client
	wsClient   *ws.Client
	config     SolanaConfig
	logger     *observability.Logger
	db         *database.DB
	redis      *database.RedisClient
	walletMgr  *WalletManager
	txService  *TransactionService
	programMgr *ProgramManager
}

// SolanaConfig holds Solana-specific configuration
type SolanaConfig struct {
	MainnetRPC string        `json:"mainnet_rpc"`
	DevnetRPC  string        `json:"devnet_rpc"`
	TestnetRPC string        `json:"testnet_rpc"`
	WSEndpoint string        `json:"ws_endpoint"`
	Commitment string        `json:"commitment"`
	Timeout    time.Duration `json:"timeout"`
	MaxRetries int           `json:"max_retries"`
	RetryDelay time.Duration `json:"retry_delay"`
}

// WalletConnection represents a connected Solana wallet
type WalletConnection struct {
	ID            uuid.UUID        `json:"id"`
	UserID        uuid.UUID        `json:"user_id"`
	PublicKey     solana.PublicKey `json:"public_key"`
	WalletType    string           `json:"wallet_type"`
	IsActive      bool             `json:"is_active"`
	Balance       decimal.Decimal  `json:"balance"`
	TokenBalances []TokenBalance   `json:"token_balances"`
	ConnectedAt   time.Time        `json:"connected_at"`
}

// TokenBalance represents a token balance in a wallet
type TokenBalance struct {
	Mint     solana.PublicKey `json:"mint"`
	Symbol   string           `json:"symbol"`
	Name     string           `json:"name"`
	Balance  decimal.Decimal  `json:"balance"`
	Decimals uint8            `json:"decimals"`
	USDValue decimal.Decimal  `json:"usd_value"`
}

// TransactionRequest represents a Solana transaction request
type TransactionRequest struct {
	From      solana.PublicKey    `json:"from"`
	To        solana.PublicKey    `json:"to"`
	Amount    decimal.Decimal     `json:"amount"`
	TokenMint *solana.PublicKey   `json:"token_mint,omitempty"`
	Priority  TransactionPriority `json:"priority"`
	MaxFee    decimal.Decimal     `json:"max_fee"`
	Memo      string              `json:"memo,omitempty"`
}

// TransactionPriority defines transaction priority levels
type TransactionPriority string

const (
	PriorityLow    TransactionPriority = "low"
	PriorityMedium TransactionPriority = "medium"
	PriorityHigh   TransactionPriority = "high"
	PriorityMax    TransactionPriority = "max"
)

// TransactionResult represents the result of a transaction
type TransactionResult struct {
	Signature solana.Signature `json:"signature"`
	BlockTime *time.Time       `json:"block_time"`
	Slot      uint64           `json:"slot"`
	Fee       uint64           `json:"fee"`
	Status    string           `json:"status"`
	Error     string           `json:"error,omitempty"`
	Logs      []string         `json:"logs,omitempty"`
}

// NewService creates a new Solana service
func NewService(db *database.DB, redis *database.RedisClient, cfg SolanaConfig, logger *observability.Logger) (*Service, error) {
	// Create RPC client
	client := rpc.New(cfg.MainnetRPC)

	// Create WebSocket client for real-time updates
	wsClient, err := ws.Connect(context.Background(), cfg.WSEndpoint)
	if err != nil {
		logger.Warn(context.Background(), "Failed to connect to Solana WebSocket", map[string]interface{}{
			"endpoint": cfg.WSEndpoint,
			"error":    err.Error(),
		})
		// Continue without WebSocket - it's not critical
	}

	service := &Service{
		client:   client,
		wsClient: wsClient,
		config:   cfg,
		logger:   logger,
		db:       db,
		redis:    redis,
	}

	// Initialize sub-services
	service.walletMgr = NewWalletManager(service)
	service.txService = NewTransactionService(service)
	service.programMgr = NewProgramManager(service)

	return service, nil
}

// GetWalletManager returns the wallet manager instance
func (s *Service) GetWalletManager() *WalletManager {
	return s.walletMgr
}

// GetTransactionService returns the transaction service instance
func (s *Service) GetTransactionService() *TransactionService {
	return s.txService
}

// GetProgramManager returns the program manager instance
func (s *Service) GetProgramManager() *ProgramManager {
	return s.programMgr
}

// GetBalance retrieves the SOL balance for a given public key
func (s *Service) GetBalance(ctx context.Context, pubkey solana.PublicKey) (decimal.Decimal, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("solana").Start(ctx, "solana.GetBalance")
	defer span.End()

	balance, err := s.client.GetBalance(ctx, pubkey, rpc.CommitmentFinalized)
	if err != nil {
		s.logger.Error(ctx, "Failed to get Solana balance", err)
		return decimal.Zero, fmt.Errorf("failed to get balance: %w", err)
	}

	// Convert lamports to SOL (1 SOL = 1e9 lamports)
	solBalance := decimal.NewFromInt(int64(balance.Value)).Div(decimal.NewFromInt(1e9))

	s.logger.Info(ctx, "Retrieved Solana balance", map[string]interface{}{
		"pubkey":  pubkey.String(),
		"balance": solBalance.String(),
	})

	return solBalance, nil
}

// GetTokenBalances retrieves all token balances for a given wallet
func (s *Service) GetTokenBalances(ctx context.Context, pubkey solana.PublicKey) ([]TokenBalance, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("solana").Start(ctx, "solana.GetTokenBalances")
	defer span.End()

	// Get token accounts by owner
	tokenAccounts, err := s.client.GetTokenAccountsByOwner(
		ctx,
		pubkey,
		&rpc.GetTokenAccountsConfig{
			ProgramId: &solana.TokenProgramID,
		},
		&rpc.GetTokenAccountsOpts{
			Commitment: rpc.CommitmentFinalized,
		},
	)
	if err != nil {
		s.logger.Error(ctx, "Failed to get token accounts", err)
		return nil, fmt.Errorf("failed to get token accounts: %w", err)
	}

	var balances []TokenBalance
	for range tokenAccounts.Value {
		// Parse token account data
		// In a real implementation, you would parse the account data properly
		// This is a simplified version for demonstration
		balance := TokenBalance{
			Mint:     solana.PublicKey{}, // Would be parsed from account data
			Symbol:   "UNKNOWN",          // Would be fetched from token metadata
			Name:     "Unknown Token",
			Balance:  decimal.Zero, // Would be parsed from amount
			Decimals: 9,            // Would be fetched from mint info
			USDValue: decimal.Zero, // Would be calculated from price feeds
		}
		balances = append(balances, balance)
	}

	s.logger.Info(ctx, "Retrieved token balances", map[string]interface{}{
		"pubkey":      pubkey.String(),
		"token_count": len(balances),
	})

	return balances, nil
}

// SendTransaction sends a transaction to the Solana network
func (s *Service) SendTransaction(ctx context.Context, req TransactionRequest) (*TransactionResult, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("solana").Start(ctx, "solana.SendTransaction")
	defer span.End()

	// Delegate to transaction service
	return s.txService.SendTransaction(ctx, req)
}

// GetTransaction retrieves transaction details by signature
func (s *Service) GetTransaction(ctx context.Context, signature solana.Signature) (*TransactionResult, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("solana").Start(ctx, "solana.GetTransaction")
	defer span.End()

	tx, err := s.client.GetTransaction(ctx, signature, &rpc.GetTransactionOpts{
		Commitment: rpc.CommitmentFinalized,
	})
	if err != nil {
		s.logger.Error(ctx, "Failed to get transaction", err)
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	result := &TransactionResult{
		Signature: signature,
		Slot:      tx.Slot,
		Fee:       tx.Meta.Fee,
		Status:    "confirmed",
	}

	if tx.BlockTime != nil {
		blockTime := time.Unix(int64(*tx.BlockTime), 0)
		result.BlockTime = &blockTime
	}

	if tx.Meta.Err != nil {
		result.Status = "failed"
		result.Error = fmt.Sprintf("%v", tx.Meta.Err)
	}

	if tx.Meta.LogMessages != nil {
		result.Logs = tx.Meta.LogMessages
	}

	return result, nil
}

// Close closes the service and cleans up resources
func (s *Service) Close() error {
	if s.wsClient != nil {
		s.wsClient.Close()
	}
	return nil
}

// Health checks the health of the Solana service
func (s *Service) Health(ctx context.Context) error {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("solana").Start(ctx, "solana.Health")
	defer span.End()

	// Check if we can get the latest blockhash
	_, err := s.client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		s.logger.Error(ctx, "Solana health check failed", err)
		return fmt.Errorf("solana health check failed: %w", err)
	}

	return nil
}
