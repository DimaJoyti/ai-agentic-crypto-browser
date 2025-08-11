package solana

import (
	"context"
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// WalletManager handles Solana wallet operations
type WalletManager struct {
	service *Service
}

// WalletAdapter interface for different wallet types
type WalletAdapter interface {
	Connect(ctx context.Context) (*WalletConnection, error)
	SignTransaction(ctx context.Context, tx *solana.Transaction) (*solana.Transaction, error)
	SignMessage(ctx context.Context, message []byte) ([]byte, error)
	GetPublicKey() solana.PublicKey
	Disconnect(ctx context.Context) error
}

// WalletType represents supported wallet types
type WalletType string

const (
	WalletTypePhantom  WalletType = "phantom"
	WalletTypeSolflare WalletType = "solflare"
	WalletTypeBackpack WalletType = "backpack"
	WalletTypeGlow     WalletType = "glow"
	WalletTypeLedger   WalletType = "ledger"
	WalletTypeTrezor   WalletType = "trezor"
)

// ConnectWalletRequest represents a wallet connection request
type ConnectWalletRequest struct {
	UserID     uuid.UUID  `json:"user_id"`
	WalletType WalletType `json:"wallet_type"`
	PublicKey  string     `json:"public_key"`
	Signature  string     `json:"signature,omitempty"`
}

// ConnectWalletResponse represents a wallet connection response
type ConnectWalletResponse struct {
	Connection *WalletConnection `json:"connection"`
	Success    bool              `json:"success"`
	Error      string            `json:"error,omitempty"`
}

// NewWalletManager creates a new wallet manager
func NewWalletManager(service *Service) *WalletManager {
	return &WalletManager{
		service: service,
	}
}

// ConnectWallet connects a Solana wallet
func (w *WalletManager) ConnectWallet(ctx context.Context, req ConnectWalletRequest) (*ConnectWalletResponse, error) {
	// Log the operation start
	w.service.logger.Info(ctx, "Connecting wallet", map[string]interface{}{
		"operation":   "ConnectWallet",
		"public_key":  req.PublicKey,
		"wallet_type": req.WalletType,
	})

	// Parse public key
	pubkey, err := solana.PublicKeyFromBase58(req.PublicKey)
	if err != nil {
		w.service.logger.Error(ctx, "Invalid public key", err)
		return &ConnectWalletResponse{
			Success: false,
			Error:   "Invalid public key format",
		}, nil
	}

	// Check if wallet already exists
	existingWallet, err := w.getWalletByPublicKey(ctx, pubkey)
	if err != nil {
		w.service.logger.Error(ctx, "Failed to check existing wallet", err)
		return &ConnectWalletResponse{
			Success: false,
			Error:   "Failed to check existing wallet",
		}, nil
	}

	if existingWallet != nil {
		// Update existing wallet
		existingWallet.IsActive = true
		existingWallet.ConnectedAt = time.Now()

		// Get current balance
		balance, err := w.service.GetBalance(ctx, pubkey)
		if err != nil {
			w.service.logger.Warn(ctx, "Failed to get wallet balance", map[string]interface{}{
				"pubkey": pubkey.String(),
				"error":  err.Error(),
			})
		} else {
			existingWallet.Balance = balance
		}

		// Get token balances
		tokenBalances, err := w.service.GetTokenBalances(ctx, pubkey)
		if err != nil {
			w.service.logger.Warn(ctx, "Failed to get token balances", map[string]interface{}{
				"pubkey": pubkey.String(),
				"error":  err.Error(),
			})
		} else {
			existingWallet.TokenBalances = tokenBalances
		}

		// Update in database
		err = w.updateWallet(ctx, existingWallet)
		if err != nil {
			w.service.logger.Error(ctx, "Failed to update wallet", err)
			return &ConnectWalletResponse{
				Success: false,
				Error:   "Failed to update wallet",
			}, nil
		}

		return &ConnectWalletResponse{
			Connection: existingWallet,
			Success:    true,
		}, nil
	}

	// Create new wallet connection
	connection := &WalletConnection{
		ID:          uuid.New(),
		UserID:      req.UserID,
		PublicKey:   pubkey,
		WalletType:  string(req.WalletType),
		IsActive:    true,
		ConnectedAt: time.Now(),
	}

	// Get wallet balance
	balance, err := w.service.GetBalance(ctx, pubkey)
	if err != nil {
		w.service.logger.Warn(ctx, "Failed to get wallet balance", map[string]interface{}{
			"pubkey": pubkey.String(),
			"error":  err.Error(),
		})
		connection.Balance = decimal.Zero
	} else {
		connection.Balance = balance
	}

	// Get token balances
	tokenBalances, err := w.service.GetTokenBalances(ctx, pubkey)
	if err != nil {
		w.service.logger.Warn(ctx, "Failed to get token balances", map[string]interface{}{
			"pubkey": pubkey.String(),
			"error":  err.Error(),
		})
		connection.TokenBalances = []TokenBalance{}
	} else {
		connection.TokenBalances = tokenBalances
	}

	// Save to database
	err = w.saveWallet(ctx, connection)
	if err != nil {
		w.service.logger.Error(ctx, "Failed to save wallet", err)
		return &ConnectWalletResponse{
			Success: false,
			Error:   "Failed to save wallet connection",
		}, nil
	}

	w.service.logger.Info(ctx, "Wallet connected successfully", map[string]interface{}{
		"user_id":     req.UserID.String(),
		"wallet_type": req.WalletType,
		"pubkey":      pubkey.String(),
		"balance":     connection.Balance.String(),
	})

	return &ConnectWalletResponse{
		Connection: connection,
		Success:    true,
	}, nil
}

// DisconnectWallet disconnects a Solana wallet
func (w *WalletManager) DisconnectWallet(ctx context.Context, userID uuid.UUID, pubkey solana.PublicKey) error {
	// Log the operation start
	w.service.logger.Info(ctx, "Disconnecting wallet", map[string]interface{}{
		"operation":  "DisconnectWallet",
		"user_id":    userID.String(),
		"public_key": pubkey.String(),
	})

	// Update wallet status in database
	query := `
		UPDATE solana_wallets 
		SET is_active = false, updated_at = NOW() 
		WHERE user_id = $1 AND public_key = $2
	`

	_, err := w.service.db.DB.Exec(query, userID, pubkey.String())
	if err != nil {
		w.service.logger.Error(ctx, "Failed to disconnect wallet", err)
		return fmt.Errorf("failed to disconnect wallet: %w", err)
	}

	w.service.logger.Info(ctx, "Wallet disconnected", map[string]interface{}{
		"user_id": userID.String(),
		"pubkey":  pubkey.String(),
	})

	return nil
}

// GetUserWallets retrieves all wallets for a user
func (w *WalletManager) GetUserWallets(ctx context.Context, userID uuid.UUID) ([]*WalletConnection, error) {
	// Log the operation start
	w.service.logger.Info(ctx, "Getting user wallets", map[string]interface{}{
		"operation": "GetUserWallets",
		"user_id":   userID.String(),
	})

	query := `
		SELECT id, user_id, public_key, wallet_type, is_active, created_at
		FROM solana_wallets 
		WHERE user_id = $1 
		ORDER BY created_at DESC
	`

	rows, err := w.service.db.DB.Query(query, userID)
	if err != nil {
		w.service.logger.Error(ctx, "Failed to get user wallets", err)
		return nil, fmt.Errorf("failed to get user wallets: %w", err)
	}
	defer rows.Close()

	var wallets []*WalletConnection
	for rows.Next() {
		var wallet WalletConnection
		var pubkeyStr string

		err := rows.Scan(
			&wallet.ID,
			&wallet.UserID,
			&pubkeyStr,
			&wallet.WalletType,
			&wallet.IsActive,
			&wallet.ConnectedAt,
		)
		if err != nil {
			w.service.logger.Error(ctx, "Failed to scan wallet row", err)
			continue
		}

		// Parse public key
		pubkey, err := solana.PublicKeyFromBase58(pubkeyStr)
		if err != nil {
			w.service.logger.Error(ctx, "Failed to parse public key", err)
			continue
		}
		wallet.PublicKey = pubkey

		// Get current balance if wallet is active
		if wallet.IsActive {
			balance, err := w.service.GetBalance(ctx, pubkey)
			if err != nil {
				w.service.logger.Warn(ctx, "Failed to get wallet balance", map[string]interface{}{
					"pubkey": pubkey.String(),
					"error":  err.Error(),
				})
				wallet.Balance = decimal.Zero
			} else {
				wallet.Balance = balance
			}
		}

		wallets = append(wallets, &wallet)
	}

	return wallets, nil
}

// Helper methods

func (w *WalletManager) getWalletByPublicKey(ctx context.Context, pubkey solana.PublicKey) (*WalletConnection, error) {
	query := `
		SELECT id, user_id, public_key, wallet_type, is_active, created_at
		FROM solana_wallets 
		WHERE public_key = $1
	`

	var wallet WalletConnection
	var pubkeyStr string

	err := w.service.db.DB.QueryRow(query, pubkey.String()).Scan(
		&wallet.ID,
		&wallet.UserID,
		&pubkeyStr,
		&wallet.WalletType,
		&wallet.IsActive,
		&wallet.ConnectedAt,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, err
	}

	wallet.PublicKey = pubkey
	return &wallet, nil
}

func (w *WalletManager) saveWallet(ctx context.Context, wallet *WalletConnection) error {
	query := `
		INSERT INTO solana_wallets (id, user_id, public_key, wallet_type, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := w.service.db.DB.Exec(query,
		wallet.ID,
		wallet.UserID,
		wallet.PublicKey.String(),
		wallet.WalletType,
		wallet.IsActive,
		wallet.ConnectedAt,
		wallet.ConnectedAt,
	)

	return err
}

func (w *WalletManager) updateWallet(ctx context.Context, wallet *WalletConnection) error {
	query := `
		UPDATE solana_wallets 
		SET is_active = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := w.service.db.DB.Exec(query,
		wallet.IsActive,
		time.Now(),
		wallet.ID,
	)

	return err
}
