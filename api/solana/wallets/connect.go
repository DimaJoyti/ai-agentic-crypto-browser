package wallets

import (
	"encoding/json"
	"net/http"

	"github.com/ai-agentic-browser/internal/auth"
	"github.com/ai-agentic-browser/internal/web3/solana"
	"github.com/ai-agentic-browser/pkg/observability"
	solanago "github.com/gagliardetto/solana-go"
	"github.com/google/uuid"
)

// ConnectWalletRequest represents a wallet connection request
type ConnectWalletRequest struct {
	PublicKey  string `json:"publicKey" validate:"required"`
	WalletType string `json:"walletType" validate:"required"`
}

// ConnectWalletResponse represents a wallet connection response
type ConnectWalletResponse struct {
	Success    bool        `json:"success"`
	Connection interface{} `json:"connection,omitempty"`
	Message    string      `json:"message,omitempty"`
}

// ConnectWalletHandler handles Solana wallet connection requests
func ConnectWalletHandler(solanaService *solana.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Only allow POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get user from context
		user, ok := auth.UserFromContext(ctx)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Parse request body
		var req ConnectWalletRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error(ctx, "Failed to decode request", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate public key
		_, err := solanago.PublicKeyFromBase58(req.PublicKey)
		if err != nil {
			logger.Error(ctx, "Invalid public key", err)
			http.Error(w, "Invalid public key format", http.StatusBadRequest)
			return
		}

		// Placeholder wallet connection
		connection := map[string]interface{}{
			"id":          uuid.New().String(),
			"user_id":     user.ID,
			"public_key":  req.PublicKey,
			"wallet_type": req.WalletType,
			"is_active":   true,
		}

		// Return success response
		response := ConnectWalletResponse{
			Success:    true,
			Connection: connection,
			Message:    "Wallet connected successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		logger.Info(ctx, "Wallet connected successfully", map[string]interface{}{
			"user_id":     user.ID,
			"public_key":  req.PublicKey,
			"wallet_type": req.WalletType,
		})
	}
}

// DisconnectWalletHandler handles wallet disconnection requests
func DisconnectWalletHandler(solanaService *solana.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Only allow POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get user from context
		user, ok := auth.UserFromContext(ctx)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Get wallet ID from URL or request body
		walletIDStr := r.URL.Query().Get("wallet_id")
		if walletIDStr == "" {
			var req struct {
				WalletID string `json:"wallet_id"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}
			walletIDStr = req.WalletID
		}

		walletID, err := uuid.Parse(walletIDStr)
		if err != nil {
			http.Error(w, "Invalid wallet ID", http.StatusBadRequest)
			return
		}

		// Placeholder wallet disconnection
		// In a real implementation, this would remove the wallet from the database

		// Return success response
		response := map[string]interface{}{
			"success": true,
			"message": "Wallet disconnected successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		logger.Info(ctx, "Wallet disconnected successfully", map[string]interface{}{
			"user_id":   user.ID,
			"wallet_id": walletID,
		})
	}
}

// GetWalletsHandler returns all connected wallets for a user
func GetWalletsHandler(solanaService *solana.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Only allow GET requests
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get user from context
		user, ok := auth.UserFromContext(ctx)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Get user's wallets from the wallet manager
		userWallets, err := solanaService.GetWalletManager().GetUserWallets(ctx, user.ID)
		if err != nil {
			http.Error(w, "Failed to get wallets", http.StatusInternalServerError)
			return
		}

		// Convert to response format
		wallets := make([]interface{}, len(userWallets))
		for i, wallet := range userWallets {
			wallets[i] = map[string]interface{}{
				"id":           wallet.ID,
				"public_key":   wallet.PublicKey,
				"wallet_type":  wallet.WalletType,
				"is_active":    wallet.IsActive,
				"connected_at": wallet.ConnectedAt,
				"balance":      wallet.Balance,
			}
		}

		// Return wallets
		response := map[string]interface{}{
			"success": true,
			"wallets": wallets,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetWalletBalanceHandler returns balance for a specific wallet
func GetWalletBalanceHandler(solanaService *solana.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Only allow GET requests
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get user from context
		user, ok := auth.UserFromContext(ctx)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Get wallet ID from URL
		walletIDStr := r.URL.Query().Get("wallet_id")
		walletID, err := uuid.Parse(walletIDStr)
		if err != nil {
			http.Error(w, "Invalid wallet ID", http.StatusBadRequest)
			return
		}

		// Get user's wallets to verify ownership
		userWallets, err := solanaService.GetWalletManager().GetUserWallets(ctx, user.ID)
		if err != nil {
			http.Error(w, "Failed to get wallets", http.StatusInternalServerError)
			return
		}

		// Find the requested wallet
		var targetWallet *solana.WalletConnection
		for _, wallet := range userWallets {
			if wallet.ID == walletID {
				targetWallet = wallet
				break
			}
		}

		if targetWallet == nil {
			http.Error(w, "Wallet not found", http.StatusNotFound)
			return
		}

		// Get balance from the service
		solBalance, err := solanaService.GetBalance(ctx, targetWallet.PublicKey)
		if err != nil {
			logger.Error(ctx, "Failed to get balance", err)
			http.Error(w, "Failed to get balance", http.StatusInternalServerError)
			return
		}

		// Get token balances
		tokenBalances, err := solanaService.GetTokenBalances(ctx, targetWallet.PublicKey)
		if err != nil {
			logger.Error(ctx, "Failed to get token balances", err)
			// Don't fail the request, just log the error
			tokenBalances = []solana.TokenBalance{}
		}

		// Format response
		balance := map[string]interface{}{
			"sol_balance":    solBalance.String(),
			"token_balances": tokenBalances,
		}

		// Return balance
		response := map[string]interface{}{
			"success": true,
			"balance": balance,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// RefreshBalanceHandler refreshes balance for a wallet
func RefreshBalanceHandler(solanaService *solana.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Only allow POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get user from context
		user, ok := auth.UserFromContext(ctx)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Get wallet ID from request
		var req struct {
			WalletID string `json:"wallet_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		walletID, err := uuid.Parse(req.WalletID)
		if err != nil {
			http.Error(w, "Invalid wallet ID", http.StatusBadRequest)
			return
		}

		// Placeholder refresh balance
		balance := map[string]interface{}{
			"sol_balance":    "0.0",
			"token_balances": []interface{}{},
		}

		// Return updated balance
		response := map[string]interface{}{
			"success": true,
			"balance": balance,
			"message": "Balance refreshed successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		logger.Info(ctx, "Wallet balance refreshed", map[string]interface{}{
			"user_id":   user.ID,
			"wallet_id": walletID,
		})
	}
}
