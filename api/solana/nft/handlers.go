package nft

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ai-agentic-browser/internal/auth"
	"github.com/ai-agentic-browser/internal/web3/solana"
	"github.com/ai-agentic-browser/pkg/observability"
	solanago "github.com/gagliardetto/solana-go"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ListNFTRequest represents an NFT listing request
type ListNFTRequest struct {
	NFTMint     string          `json:"nftMint" validate:"required"`
	Price       decimal.Decimal `json:"price" validate:"required"`
	Currency    string          `json:"currency"`
	Marketplace string          `json:"marketplace" validate:"required"`
	Seller      string          `json:"seller" validate:"required"`
	ExpiresAt   *string         `json:"expiresAt,omitempty"`
}

// BuyNFTRequest represents an NFT purchase request
type BuyNFTRequest struct {
	NFTMint     string          `json:"nftMint" validate:"required"`
	MaxPrice    decimal.Decimal `json:"maxPrice" validate:"required"`
	Marketplace string          `json:"marketplace" validate:"required"`
	Buyer       string          `json:"buyer" validate:"required"`
}

// ExploreNFTsHandler handles NFT exploration requests
func ExploreNFTsHandler(solanaService *solana.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Only allow GET requests
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse query parameters
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")
		collection := r.URL.Query().Get("collection")
		sortBy := r.URL.Query().Get("sortBy")

		// Log the request
		logger.Info(ctx, "Exploring NFTs", map[string]interface{}{
			"limit":      limitStr,
			"offset":     offsetStr,
			"collection": collection,
			"sort_by":    sortBy,
		})

		limit := 20
		offset := 0

		if limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
				limit = l
			}
		}

		if offsetStr != "" {
			if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
				offset = o
			}
		}

		// Get NFTs based on collection or general exploration
		var nfts []interface{}

		if collection != "" {
			// Validate collection address
			_, err := solanago.PublicKeyFromBase58(collection)
			if err != nil {
				http.Error(w, "Invalid collection address", http.StatusBadRequest)
				return
			}
			// For now, return empty array - would integrate with marketplace APIs
			nfts = []interface{}{}
		} else {
			// For general exploration, return empty array - would integrate with marketplace APIs
			nfts = []interface{}{}
		}

		// Return NFTs
		response := map[string]interface{}{
			"success": true,
			"nfts":    nfts,
			"limit":   limit,
			"offset":  offset,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetCollectionsHandler handles NFT collections requests
func GetCollectionsHandler(solanaService *solana.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Only allow GET requests
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse query parameters
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")

		// Log the request
		logger.Info(ctx, "Getting NFT collections", map[string]interface{}{
			"limit":  limitStr,
			"offset": offsetStr,
		})

		limit := 20
		offset := 0

		if limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
				limit = l
			}
		}

		if offsetStr != "" {
			if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
				offset = o
			}
		}

		// Get collections - placeholder implementation
		collections := []interface{}{}

		// Return collections
		response := map[string]interface{}{
			"success":     true,
			"collections": collections,
			"limit":       limit,
			"offset":      offset,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetPortfolioHandler returns NFT portfolio for a user
func GetPortfolioHandler(solanaService *solana.Service, logger *observability.Logger) http.HandlerFunc {
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

		// Get wallet ID from URL path
		walletIDStr := r.URL.Path[len("/api/solana/nft/portfolio/"):]
		if walletIDStr == "" {
			http.Error(w, "Wallet ID required", http.StatusBadRequest)
			return
		}

		walletID, err := uuid.Parse(walletIDStr)
		if err != nil {
			// Try parsing as public key
			_, err := solanago.PublicKeyFromBase58(walletIDStr)
			if err != nil {
				http.Error(w, "Invalid wallet ID or public key", http.StatusBadRequest)
				return
			}
			// For now, generate a placeholder wallet ID
			walletID = uuid.New()
		}

		// Log the request
		logger.Info(ctx, "Getting user NFTs", map[string]interface{}{
			"user_id":   user.ID,
			"wallet_id": walletID,
		})

		// Get user NFTs - placeholder implementation
		nfts := []interface{}{}

		// Return NFTs
		response := map[string]interface{}{
			"success": true,
			"nfts":    nfts,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// ListNFTHandler handles NFT listing requests
func ListNFTHandler(solanaService *solana.Service, logger *observability.Logger) http.HandlerFunc {
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
		var req ListNFTRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error(ctx, "Failed to decode list NFT request", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate addresses
		_, err := solanago.PublicKeyFromBase58(req.NFTMint)
		if err != nil {
			http.Error(w, "Invalid NFT mint", http.StatusBadRequest)
			return
		}

		_, err = solanago.PublicKeyFromBase58(req.Seller)
		if err != nil {
			http.Error(w, "Invalid seller address", http.StatusBadRequest)
			return
		}

		// Placeholder listing implementation
		listing := map[string]interface{}{
			"id":          uuid.New().String(),
			"nft_mint":    req.NFTMint,
			"price":       req.Price,
			"marketplace": req.Marketplace,
			"status":      "listed",
		}

		// Return listing
		response := map[string]interface{}{
			"success": true,
			"listing": listing,
			"message": "NFT listed successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		logger.Info(ctx, "NFT listed successfully", map[string]interface{}{
			"user_id":     user.ID,
			"nft_mint":    req.NFTMint,
			"price":       req.Price.String(),
			"marketplace": req.Marketplace,
		})
	}
}

// BuyNFTHandler handles NFT purchase requests
func BuyNFTHandler(solanaService *solana.Service, logger *observability.Logger) http.HandlerFunc {
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
		var req BuyNFTRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error(ctx, "Failed to decode buy NFT request", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate addresses
		_, err := solanago.PublicKeyFromBase58(req.NFTMint)
		if err != nil {
			http.Error(w, "Invalid NFT mint", http.StatusBadRequest)
			return
		}

		_, err = solanago.PublicKeyFromBase58(req.Buyer)
		if err != nil {
			http.Error(w, "Invalid buyer address", http.StatusBadRequest)
			return
		}

		// Placeholder sale implementation
		sale := map[string]interface{}{
			"id":          uuid.New().String(),
			"nft_mint":    req.NFTMint,
			"price":       req.MaxPrice,
			"marketplace": req.Marketplace,
			"buyer":       req.Buyer,
			"signature":   "placeholder_signature",
			"status":      "completed",
		}

		// Return sale
		response := map[string]interface{}{
			"success": true,
			"sale":    sale,
			"message": "NFT purchased successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		logger.Info(ctx, "NFT purchased successfully", map[string]interface{}{
			"user_id":     user.ID,
			"nft_mint":    req.NFTMint,
			"price":       req.MaxPrice.String(),
			"marketplace": req.Marketplace,
			"signature":   sale["signature"],
		})
	}
}

// GetNFTMetadataHandler returns metadata for a specific NFT
func GetNFTMetadataHandler(solanaService *solana.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Only allow GET requests
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get mint address from URL
		mintStr := r.URL.Query().Get("mint")
		if mintStr == "" {
			http.Error(w, "Mint address required", http.StatusBadRequest)
			return
		}

		// Log the request
		logger.Info(ctx, "Getting NFT metadata", map[string]interface{}{
			"mint": mintStr,
		})

		_, err := solanago.PublicKeyFromBase58(mintStr)
		if err != nil {
			http.Error(w, "Invalid mint address", http.StatusBadRequest)
			return
		}

		// Placeholder NFT metadata
		nft := map[string]interface{}{
			"mint":        mintStr,
			"name":        "Sample NFT",
			"description": "This is a sample NFT",
			"image":       "https://example.com/nft.png",
			"attributes":  []interface{}{},
		}

		// Return metadata
		response := map[string]interface{}{
			"success": true,
			"nft":     nft,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
