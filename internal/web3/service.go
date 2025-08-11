package web3

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ai-agentic-browser/internal/config"
	"github.com/ai-agentic-browser/pkg/database"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Service provides Web3 and cryptocurrency functionality
type Service struct {
	db         *database.DB
	redis      *database.RedisClient
	config     config.Web3Config
	logger     *observability.Logger
	providers  map[int]*ChainProvider
	walletRepo WalletRepository
	txRepo     TransactionRepository
}

// ChainProvider represents a blockchain provider
type ChainProvider struct {
	ChainID int
	RpcURL  string
	Client  interface{} // lazily set to *ethclient.Client when first used
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

	walletRepo := NewPostgresWalletRepository(db)
	txRepo := NewPostgresTransactionRepository(db)

	return &Service{
		db:         db,
		redis:      redis,
		config:     cfg,
		logger:     logger,
		providers:  providers,
		walletRepo: walletRepo,
		txRepo:     txRepo,
	}
}

// ConnectWallet connects a cryptocurrency wallet
func (s *Service) ConnectWallet(ctx context.Context, userID uuid.UUID, req WalletConnectRequest) (*WalletConnectResponse, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("web3-service").Start(ctx, "web3.ConnectWallet")
	defer span.End()

	// Validate input
	if req.Address == "" || len(req.Address) < 4 {
		return nil, fmt.Errorf("invalid address")
	}
	if _, exists := SupportedChains[req.ChainID]; !exists {
		return nil, fmt.Errorf("unsupported chain ID: %d", req.ChainID)
	}

	// Check if wallet already exists
	existingWallet, err := s.walletRepo.GetByAddress(ctx, userID, req.Address, req.ChainID)
	if err == nil && existingWallet != nil {
		return &WalletConnectResponse{Wallet: existingWallet, Message: "Wallet already connected"}, nil
	}

	// Create new wallet
	wallet := &Wallet{
		ID:         uuid.New(),
		UserID:     userID,
		Address:    strings.ToLower(req.Address),
		ChainID:    req.ChainID,
		WalletType: req.WalletType,
		IsPrimary:  false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Check if this is the user's first wallet
	walletCount, err := s.walletRepo.CountByUser(ctx, userID)
	if err != nil {
		s.logger.Error(ctx, "Failed to get user wallet count", err)
	} else if walletCount == 0 {
		wallet.IsPrimary = true
	}

	// Save wallet to database
	if err := s.walletRepo.Save(ctx, wallet); err != nil {
		s.logger.Error(ctx, "Failed to save wallet", err)
		return nil, fmt.Errorf("failed to save wallet: %w", err)
	}

	s.logger.Info(ctx, "Wallet connected successfully", map[string]any{
		"wallet_id":   wallet.ID.String(),
		"user_id":     userID.String(),
		"address":     wallet.Address,
		"chain_id":    req.ChainID,
		"wallet_type": req.WalletType,
	})

	return &WalletConnectResponse{Wallet: wallet, Message: "Wallet connected successfully"}, nil
}

// GetBalance retrieves wallet balance information
func (s *Service) GetBalance(ctx context.Context, userID uuid.UUID, req BalanceRequest) (*BalanceResponse, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("web3-service").Start(ctx, "web3.GetBalance")
	defer span.End()

	var address string
	var chainID int

	// Determine address and chain ID
	if req.WalletID != uuid.Nil {
		wallet, err := s.walletRepo.GetByID(ctx, req.WalletID)
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

	// Fetch native balance
	nativeBalance, err := s.getNativeBalance(ctx, chainID, address)
	if err != nil {
		s.logger.Warn(ctx, "Failed to fetch native balance", map[string]any{"error": err.Error(), "address": address, "chain_id": chainID})
		nativeBalance = big.NewInt(0)
	}

	// Fetch common ERC-20 token balances with caching
	var tokenBalances []TokenBalance
	if tokens, ok := CommonERC20Tokens[chainID]; ok {
		for _, t := range tokens {
			decimals, derr := s.getERC20Decimals(ctx, chainID, t.Address)
			if derr != nil {
				s.logger.Warn(ctx, "Failed to read token decimals", map[string]any{"error": derr.Error(), "token": t.Symbol, "address": t.Address, "chain_id": chainID})
				continue
			}
			bal, berr := s.getERC20Balance(ctx, chainID, t.Address, address)
			if berr != nil {
				s.logger.Warn(ctx, "Failed to read token balance", map[string]any{"error": berr.Error(), "token": t.Symbol, "address": t.Address, "chain_id": chainID})
				continue
			}
			tokenBalances = append(tokenBalances, TokenBalance{
				TokenAddress: t.Address,
				TokenSymbol:  t.Symbol,
				TokenName:    t.Name,
				Balance:      bal,
				Decimals:     decimals,
				USDValue:     0, // priced elsewhere
			})
		}
	}

	// Price native and tokens via CoinGecko (USD)
	cg := NewCoinGeckoClient(s.redis)
	nativeID := NativeCoinGeckoIDByChain[chainID]
	priceIDs := []string{}
	if nativeID != "" {
		priceIDs = append(priceIDs, nativeID)
	}
	for _, t := range tokenBalances {
		// look up CG id from CommonERC20Tokens map by address
		for _, meta := range CommonERC20Tokens[chainID] {
			if strings.EqualFold(meta.Address, t.TokenAddress) && meta.CoinGeckoID != "" {
				priceIDs = append(priceIDs, meta.CoinGeckoID)
			}
		}
	}
	prices := map[string]TokenPrice{}
	if len(priceIDs) > 0 {
		if p, err := cg.GetPrices(ctx, "USD", priceIDs); err == nil {
			prices = p
		} else {
			s.logger.Warn(ctx, "Price fetch failed", map[string]any{"error": err.Error()})
		}
	}

	// Compute USD values
	totalUSD := 0.0
	// native
	if nativeID != "" {
		if pt, ok := prices[nativeID]; ok {
			// native decimals are 18 for ETH; Polygon MATIC also 18 (commonly). Keep it simple: 18.
			usd := new(big.Float).Quo(new(big.Float).SetInt(nativeBalance), new(big.Float).SetFloat64(1e18))
			v, _ := new(big.Float).Mul(usd, big.NewFloat(pt.Price)).Float64()
			totalUSD += v
		}
	}
	// tokens
	for i := range tokenBalances {
		metaID := ""
		for _, meta := range CommonERC20Tokens[chainID] {
			if strings.EqualFold(meta.Address, tokenBalances[i].TokenAddress) {
				metaID = meta.CoinGeckoID
				break
			}
		}
		if metaID == "" {
			continue
		}
		pt, ok := prices[metaID]
		if !ok {
			continue
		}
		decPow := new(big.Float).SetFloat64(1.0)
		for j := 0; j < tokenBalances[i].Decimals; j++ {
			decPow = new(big.Float).Mul(decPow, big.NewFloat(10))
		}
		amt := new(big.Float).Quo(new(big.Float).SetInt(tokenBalances[i].Balance), decPow)
		v, _ := new(big.Float).Mul(amt, big.NewFloat(pt.Price)).Float64()
		tokenBalances[i].USDValue = v
		totalUSD += v
	}

	response := &BalanceResponse{
		Address:       address,
		ChainID:       chainID,
		NativeBalance: nativeBalance,
		TokenBalances: tokenBalances,
		TotalUSDValue: totalUSD,
		Metadata: map[string]any{
			"provider":   provider.RpcURL,
			"chain_name": SupportedChains[chainID],
			"timestamp":  time.Now(),
		},
	}

	s.logger.Info(ctx, "Balance retrieved", map[string]any{
		"address":   address,
		"chain_id":  chainID,
		"tokens":    len(tokenBalances),
		"total_usd": totalUSD,
	})

	return response, nil
}

// CreateTransaction creates a new blockchain transaction
func (s *Service) CreateTransaction(ctx context.Context, userID uuid.UUID, req TransactionRequest) (*TransactionResponse, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("web3-service").Start(ctx, "web3.CreateTransaction")
	defer span.End()

	// Get wallet
	wallet, err := s.walletRepo.GetByID(ctx, req.WalletID)
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
	if err := s.txRepo.Save(ctx, transaction); err != nil {
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

	s.logger.Info(ctx, "Transaction created", map[string]any{
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

	// Normalize currency and token IDs
	currency := req.Currency
	if currency == "" {
		currency = "USD"
	}
	// Default tokens if none provided
	ids := []string{"ethereum", "bitcoin", "polygon"}
	if strings.TrimSpace(req.Token) != "" {
		ids = []string{strings.ToLower(req.Token)}
	}

	cg := NewCoinGeckoClient(s.redis)
	prices, err := cg.GetPrices(ctx, currency, ids)
	if err != nil {
		s.logger.Error(ctx, "CoinGecko price fetch failed", err)
		return nil, fmt.Errorf("failed to fetch prices: %w", err)
	}

	response := &PriceResponse{
		Prices:    prices,
		Currency:  strings.ToUpper(currency),
		Timestamp: time.Now(),
	}

	s.logger.Info(ctx, "Prices retrieved", map[string]any{
		"currency":    response.Currency,
		"price_count": len(prices),
	})

	return response, nil
}

// ListWallets returns user's wallets with filters and pagination
func (s *Service) ListWallets(ctx context.Context, userID uuid.UUID, filter WalletListFilter) ([]*Wallet, Pagination, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 || filter.PageSize > 100 {
		filter.PageSize = 20
	}
	return s.walletRepo.ListByUser(ctx, userID, filter)
}

// ListTransactions returns user's transactions with filters and pagination
func (s *Service) ListTransactions(ctx context.Context, userID uuid.UUID, filter TransactionListFilter) ([]*Transaction, Pagination, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 || filter.PageSize > 100 {
		filter.PageSize = 20
	}
	return s.txRepo.ListByUser(ctx, userID, filter)
}

// InteractWithDeFiProtocol interacts with DeFi protocols
func (s *Service) InteractWithDeFiProtocol(ctx context.Context, userID uuid.UUID, req DeFiProtocolRequest) (*DeFiProtocolResponse, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("web3-service").Start(ctx, "web3.InteractWithDeFiProtocol")
	defer span.End()

	// Get wallet
	wallet, err := s.walletRepo.GetByID(ctx, req.WalletID)
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
		Metadata: map[string]any{
			"protocol":  req.Protocol,
			"action":    req.Action,
			"timestamp": time.Now(),
		},
	}

	s.logger.Info(ctx, "DeFi interaction completed", map[string]any{
		"protocol":  req.Protocol,
		"action":    req.Action,
		"tx_hash":   txHash,
		"wallet_id": req.WalletID.String(),
	})

	return response, nil
}

// Repository-backed helper methods (kept minimal)

// simulateTransactionConfirmation simulates transaction confirmation
func (s *Service) simulateTransactionConfirmation(ctx context.Context, tx *Transaction) {
	// Simulate network delay
	time.Sleep(5 * time.Second)

	// Update transaction status in DB
	_ = s.txRepo.UpdateStatus(ctx, tx.ID, TxStatusConfirmed)

	s.logger.Info(ctx, "Transaction confirmed", map[string]any{
		"tx_hash": tx.TxHash,
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
