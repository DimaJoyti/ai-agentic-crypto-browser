package web3

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ai-agentic-browser/internal/config"
	"github.com/ai-agentic-browser/pkg/database"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Service provides Web3 and cryptocurrency functionality
type Service struct {
	db        *database.DB
	redis     *database.RedisClient
	config    config.Web3Config
	logger    *observability.Logger
	providers map[int]*ChainProvider
}

// ChainProvider represents a blockchain provider
type ChainProvider struct {
	ChainID int
	RpcURL  string
	Client  interface{} // Would be ethclient.Client in real implementation
}

// NewService creates a new Web3 service
func NewService(db *database.DB, redis *database.RedisClient, cfg config.Web3Config, logger *observability.Logger) *Service {
	providers := make(map[int]*ChainProvider)

	// Initialize providers for supported chains
	if cfg.EthereumRPC != "" {
		providers[1] = &ChainProvider{ChainID: 1, RpcURL: cfg.EthereumRPC}
	}
	if cfg.PolygonRPC != "" {
		providers[137] = &ChainProvider{ChainID: 137, RpcURL: cfg.PolygonRPC}
	}
	if cfg.ArbitrumRPC != "" {
		providers[42161] = &ChainProvider{ChainID: 42161, RpcURL: cfg.ArbitrumRPC}
	}
	if cfg.OptimismRPC != "" {
		providers[10] = &ChainProvider{ChainID: 10, RpcURL: cfg.OptimismRPC}
	}

	return &Service{
		db:        db,
		redis:     redis,
		config:    cfg,
		logger:    logger,
		providers: providers,
	}
}

// ConnectWallet connects a cryptocurrency wallet
func (s *Service) ConnectWallet(ctx context.Context, userID uuid.UUID, req WalletConnectRequest) (*WalletConnectResponse, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("web3-service").Start(ctx, "web3.ConnectWallet")
	defer span.End()

	// Validate chain support
	if _, exists := SupportedChains[req.ChainID]; !exists {
		return nil, fmt.Errorf("unsupported chain ID: %d", req.ChainID)
	}

	// Check if wallet already exists
	existingWallet, err := s.getWalletByAddress(ctx, userID, req.Address, req.ChainID)
	if err == nil && existingWallet != nil {
		return &WalletConnectResponse{
			Wallet:  existingWallet,
			Message: "Wallet already connected",
		}, nil
	}

	// Create new wallet
	wallet := &Wallet{
		ID:         uuid.New(),
		UserID:     userID,
		Address:    req.Address,
		ChainID:    req.ChainID,
		WalletType: req.WalletType,
		IsPrimary:  false, // Will be set to true if it's the first wallet
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Check if this is the user's first wallet
	walletCount, err := s.getUserWalletCount(ctx, userID)
	if err != nil {
		s.logger.Error(ctx, "Failed to get user wallet count", err)
	} else if walletCount == 0 {
		wallet.IsPrimary = true
	}

	// Save wallet to database
	if err := s.saveWallet(ctx, wallet); err != nil {
		s.logger.Error(ctx, "Failed to save wallet", err)
		return nil, fmt.Errorf("failed to save wallet: %w", err)
	}

	s.logger.Info(ctx, "Wallet connected successfully", map[string]interface{}{
		"wallet_id":   wallet.ID.String(),
		"user_id":     userID.String(),
		"address":     req.Address,
		"chain_id":    req.ChainID,
		"wallet_type": req.WalletType,
	})

	return &WalletConnectResponse{
		Wallet:  wallet,
		Message: "Wallet connected successfully",
	}, nil
}

// GetBalance retrieves wallet balance information
func (s *Service) GetBalance(ctx context.Context, userID uuid.UUID, req BalanceRequest) (*BalanceResponse, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("web3-service").Start(ctx, "web3.GetBalance")
	defer span.End()

	var address string
	var chainID int

	// Determine address and chain ID
	if req.WalletID != uuid.Nil {
		wallet, err := s.getWalletByID(ctx, req.WalletID)
		if err != nil {
			return nil, fmt.Errorf("wallet not found: %w", err)
		}
		if wallet.UserID != userID {
			return nil, fmt.Errorf("wallet does not belong to user")
		}
		address = wallet.Address
		chainID = wallet.ChainID
	} else if req.Address != "" && req.ChainID != 0 {
		address = req.Address
		chainID = req.ChainID
	} else {
		return nil, fmt.Errorf("either wallet_id or address+chain_id must be provided")
	}

	// Get provider for chain
	provider, exists := s.providers[chainID]
	if !exists {
		return nil, fmt.Errorf("no provider configured for chain ID: %d", chainID)
	}

	// For demo purposes, return mock data
	// In a real implementation, this would query the blockchain
	nativeBalance := big.NewInt(1000000000000000000) // 1 ETH in wei

	tokenBalances := []TokenBalance{
		{
			TokenAddress: "0xA0b86a33E6441b8C4505E2E8E3C3C5C8E6441b8C",
			TokenSymbol:  "USDC",
			TokenName:    "USD Coin",
			Balance:      big.NewInt(1000000000), // 1000 USDC (6 decimals)
			Decimals:     6,
			USDValue:     1000.0,
		},
		{
			TokenAddress: "0xB0b86a33E6441b8C4505E2E8E3C3C5C8E6441b8C",
			TokenSymbol:  "USDT",
			TokenName:    "Tether USD",
			Balance:      big.NewInt(500000000), // 500 USDT (6 decimals)
			Decimals:     6,
			USDValue:     500.0,
		},
	}

	response := &BalanceResponse{
		Address:       address,
		ChainID:       chainID,
		NativeBalance: nativeBalance,
		TokenBalances: tokenBalances,
		TotalUSDValue: 3500.0, // Mock total value
		Metadata: map[string]interface{}{
			"provider":   provider.RpcURL,
			"chain_name": SupportedChains[chainID],
			"timestamp":  time.Now(),
		},
	}

	s.logger.Info(ctx, "Balance retrieved", map[string]interface{}{
		"address":         address,
		"chain_id":        chainID,
		"total_usd_value": response.TotalUSDValue,
	})

	return response, nil
}

// CreateTransaction creates a new blockchain transaction
func (s *Service) CreateTransaction(ctx context.Context, userID uuid.UUID, req TransactionRequest) (*TransactionResponse, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("web3-service").Start(ctx, "web3.CreateTransaction")
	defer span.End()

	// Get wallet
	wallet, err := s.getWalletByID(ctx, req.WalletID)
	if err != nil {
		return nil, fmt.Errorf("wallet not found: %w", err)
	}
	if wallet.UserID != userID {
		return nil, fmt.Errorf("wallet does not belong to user")
	}

	// Validate chain support
	if _, exists := s.providers[wallet.ChainID]; !exists {
		return nil, fmt.Errorf("no provider configured for chain ID: %d", wallet.ChainID)
	}

	// Create transaction record
	transaction := &Transaction{
		ID:              uuid.New(),
		UserID:          userID,
		WalletID:        req.WalletID,
		TxHash:          fmt.Sprintf("0x%x", time.Now().UnixNano()), // Mock hash
		ChainID:         wallet.ChainID,
		FromAddress:     wallet.Address,
		ToAddress:       req.ToAddress,
		Value:           req.Value,
		Status:          TxStatusPending,
		TransactionType: "transfer",
		Metadata:        req.Metadata,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Save transaction to database
	if err := s.saveTransaction(ctx, transaction); err != nil {
		s.logger.Error(ctx, "Failed to save transaction", err)
		return nil, fmt.Errorf("failed to save transaction: %w", err)
	}

	// In a real implementation, this would broadcast the transaction to the network
	// For demo purposes, we'll simulate a successful transaction
	go s.simulateTransactionConfirmation(context.Background(), transaction)

	response := &TransactionResponse{
		Transaction: transaction,
		TxHash:      transaction.TxHash,
		Status:      string(transaction.Status),
	}

	s.logger.Info(ctx, "Transaction created", map[string]interface{}{
		"tx_id":    transaction.ID.String(),
		"tx_hash":  transaction.TxHash,
		"from":     transaction.FromAddress,
		"to":       req.ToAddress,
		"chain_id": wallet.ChainID,
	})

	return response, nil
}

// GetPrices retrieves cryptocurrency prices
func (s *Service) GetPrices(ctx context.Context, req PriceRequest) (*PriceResponse, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("web3-service").Start(ctx, "web3.GetPrices")
	defer span.End()

	// For demo purposes, return mock price data
	// In a real implementation, this would query a price API like CoinGecko
	prices := map[string]TokenPrice{
		"ethereum": {
			Symbol:          "ETH",
			Name:            "Ethereum",
			Price:           2500.0,
			PriceChange24h:  50.0,
			PriceChangePerc: 2.04,
			MarketCap:       300000000000,
			Volume24h:       15000000000,
			LastUpdated:     time.Now(),
		},
		"bitcoin": {
			Symbol:          "BTC",
			Name:            "Bitcoin",
			Price:           45000.0,
			PriceChange24h:  1000.0,
			PriceChangePerc: 2.27,
			MarketCap:       850000000000,
			Volume24h:       25000000000,
			LastUpdated:     time.Now(),
		},
		"polygon": {
			Symbol:          "MATIC",
			Name:            "Polygon",
			Price:           0.85,
			PriceChange24h:  0.05,
			PriceChangePerc: 6.25,
			MarketCap:       8000000000,
			Volume24h:       500000000,
			LastUpdated:     time.Now(),
		},
	}

	currency := req.Currency
	if currency == "" {
		currency = "USD"
	}

	response := &PriceResponse{
		Prices:    prices,
		Currency:  currency,
		Timestamp: time.Now(),
	}

	s.logger.Info(ctx, "Prices retrieved", map[string]interface{}{
		"currency":    currency,
		"price_count": len(prices),
	})

	return response, nil
}

// InteractWithDeFiProtocol interacts with DeFi protocols
func (s *Service) InteractWithDeFiProtocol(ctx context.Context, userID uuid.UUID, req DeFiProtocolRequest) (*DeFiProtocolResponse, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("web3-service").Start(ctx, "web3.InteractWithDeFiProtocol")
	defer span.End()

	// Get wallet
	wallet, err := s.getWalletByID(ctx, req.WalletID)
	if err != nil {
		return nil, fmt.Errorf("wallet not found: %w", err)
	}
	if wallet.UserID != userID {
		return nil, fmt.Errorf("wallet does not belong to user")
	}

	// For demo purposes, simulate DeFi interaction
	// In a real implementation, this would interact with smart contracts

	var position *DeFiPosition
	var txHash string

	switch req.Protocol {
	case "uniswap":
		txHash, position, err = s.simulateUniswapInteraction(ctx, wallet, req)
	case "aave":
		txHash, position, err = s.simulateAaveInteraction(ctx, wallet, req)
	case "compound":
		txHash, position, err = s.simulateCompoundInteraction(ctx, wallet, req)
	default:
		return &DeFiProtocolResponse{
			Success: false,
			Error:   fmt.Sprintf("unsupported protocol: %s", req.Protocol),
		}, nil
	}

	if err != nil {
		s.logger.Error(ctx, "DeFi interaction failed", err)
		return &DeFiProtocolResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	response := &DeFiProtocolResponse{
		Success:  true,
		TxHash:   txHash,
		Position: position,
		Metadata: map[string]interface{}{
			"protocol":  req.Protocol,
			"action":    req.Action,
			"timestamp": time.Now(),
		},
	}

	s.logger.Info(ctx, "DeFi interaction completed", map[string]interface{}{
		"protocol":  req.Protocol,
		"action":    req.Action,
		"tx_hash":   txHash,
		"wallet_id": req.WalletID.String(),
	})

	return response, nil
}

// Helper methods

// getWalletByAddress retrieves a wallet by address and chain ID
func (s *Service) getWalletByAddress(ctx context.Context, userID uuid.UUID, address string, chainID int) (*Wallet, error) {
	query := `
		SELECT id, user_id, address, chain_id, wallet_type, is_primary, created_at, updated_at
		FROM web3_wallets WHERE user_id = $1 AND address = $2 AND chain_id = $3
	`
	wallet := &Wallet{}
	err := s.db.QueryRowContext(ctx, query, userID, address, chainID).Scan(
		&wallet.ID, &wallet.UserID, &wallet.Address, &wallet.ChainID,
		&wallet.WalletType, &wallet.IsPrimary, &wallet.CreatedAt, &wallet.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

// getWalletByID retrieves a wallet by ID
func (s *Service) getWalletByID(ctx context.Context, walletID uuid.UUID) (*Wallet, error) {
	query := `
		SELECT id, user_id, address, chain_id, wallet_type, is_primary, created_at, updated_at
		FROM web3_wallets WHERE id = $1
	`
	wallet := &Wallet{}
	err := s.db.QueryRowContext(ctx, query, walletID).Scan(
		&wallet.ID, &wallet.UserID, &wallet.Address, &wallet.ChainID,
		&wallet.WalletType, &wallet.IsPrimary, &wallet.CreatedAt, &wallet.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

// getUserWalletCount gets the number of wallets for a user
func (s *Service) getUserWalletCount(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM web3_wallets WHERE user_id = $1`
	var count int
	err := s.db.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}

// saveWallet saves a wallet to the database
func (s *Service) saveWallet(ctx context.Context, wallet *Wallet) error {
	query := `
		INSERT INTO web3_wallets (id, user_id, address, chain_id, wallet_type, is_primary, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := s.db.ExecContext(ctx, query, wallet.ID, wallet.UserID, wallet.Address, wallet.ChainID,
		wallet.WalletType, wallet.IsPrimary, wallet.CreatedAt, wallet.UpdatedAt)
	return err
}

// saveTransaction saves a transaction to the database
func (s *Service) saveTransaction(ctx context.Context, tx *Transaction) error {
	metadataJSON, _ := json.Marshal(tx.Metadata)
	query := `
		INSERT INTO web3_transactions (id, user_id, wallet_id, tx_hash, chain_id, from_address, to_address, value, status, transaction_type, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := s.db.ExecContext(ctx, query, tx.ID, tx.UserID, tx.WalletID, tx.TxHash, tx.ChainID,
		tx.FromAddress, tx.ToAddress, tx.Value, tx.Status, tx.TransactionType, metadataJSON, tx.CreatedAt, tx.UpdatedAt)
	return err
}

// simulateTransactionConfirmation simulates transaction confirmation
func (s *Service) simulateTransactionConfirmation(ctx context.Context, tx *Transaction) {
	// Simulate network delay
	time.Sleep(5 * time.Second)

	// Update transaction status
	tx.Status = TxStatusConfirmed
	tx.BlockNumber = 18500000 // Mock block number
	tx.GasUsed = 21000        // Mock gas used
	tx.UpdatedAt = time.Now()

	// In a real implementation, this would update the database
	s.logger.Info(ctx, "Transaction confirmed", map[string]interface{}{
		"tx_hash":      tx.TxHash,
		"block_number": fmt.Sprintf("%d", tx.BlockNumber),
	})
}

// Simulate DeFi protocol interactions
func (s *Service) simulateUniswapInteraction(ctx context.Context, wallet *Wallet, req DeFiProtocolRequest) (string, *DeFiPosition, error) {
	txHash := fmt.Sprintf("0x%x", time.Now().UnixNano())

	position := &DeFiPosition{
		ID:           uuid.New(),
		UserID:       wallet.UserID,
		WalletID:     wallet.ID,
		ProtocolName: "uniswap",
		PositionType: "liquidity_pool",
		TokenSymbol:  "ETH-USDC",
		Amount:       req.Amount,
		USDValue:     decimal.NewFromFloat(5000.0),
		APY:          decimal.NewFromFloat(12.5),
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return txHash, position, nil
}

func (s *Service) simulateAaveInteraction(ctx context.Context, wallet *Wallet, req DeFiProtocolRequest) (string, *DeFiPosition, error) {
	txHash := fmt.Sprintf("0x%x", time.Now().UnixNano())

	position := &DeFiPosition{
		ID:           uuid.New(),
		UserID:       wallet.UserID,
		WalletID:     wallet.ID,
		ProtocolName: "aave",
		PositionType: "lending",
		TokenSymbol:  "USDC",
		Amount:       req.Amount,
		USDValue:     decimal.NewFromFloat(10000.0),
		APY:          decimal.NewFromFloat(4.2),
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return txHash, position, nil
}

func (s *Service) simulateCompoundInteraction(ctx context.Context, wallet *Wallet, req DeFiProtocolRequest) (string, *DeFiPosition, error) {
	txHash := fmt.Sprintf("0x%x", time.Now().UnixNano())

	position := &DeFiPosition{
		ID:           uuid.New(),
		UserID:       wallet.UserID,
		WalletID:     wallet.ID,
		ProtocolName: "compound",
		PositionType: "lending",
		TokenSymbol:  "DAI",
		Amount:       req.Amount,
		USDValue:     decimal.NewFromFloat(7500.0),
		APY:          decimal.NewFromFloat(3.8),
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return txHash, position, nil
}
