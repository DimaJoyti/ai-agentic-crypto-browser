package solana

import (
	"encoding/json"
	"net/http"

	"github.com/ai-agentic-browser/api/solana/defi"
	"github.com/ai-agentic-browser/api/solana/nft"
	"github.com/ai-agentic-browser/api/solana/wallets"
	"github.com/ai-agentic-browser/internal/web3/solana"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

// SolanaStats represents Solana network statistics
type SolanaStats struct {
	Price         decimal.Decimal `json:"price"`
	PriceChange24h decimal.Decimal `json:"priceChange24h"`
	MarketCap     decimal.Decimal `json:"marketCap"`
	Volume24h     decimal.Decimal `json:"volume24h"`
	TPS           int             `json:"tps"`
	TotalAccounts int             `json:"totalAccounts"`
}

// DashboardStats represents user dashboard statistics
type DashboardStats struct {
	TotalValue      decimal.Decimal `json:"totalValue"`
	TotalPnL        decimal.Decimal `json:"totalPnL"`
	ActivePositions int             `json:"activePositions"`
	NFTCount        int             `json:"nftCount"`
	SOLBalance      decimal.Decimal `json:"solBalance"`
	TokenCount      int             `json:"tokenCount"`
}

// SetupSolanaRoutes sets up all Solana-related API routes
func SetupSolanaRoutes(r *mux.Router, solanaService *solana.Service, logger *observability.Logger) {
	// Create Solana subrouter
	solanaRouter := r.PathPrefix("/solana").Subrouter()

	// Network stats endpoint
	solanaRouter.HandleFunc("/stats", GetSolanaStatsHandler(logger)).Methods("GET")
	
	// Dashboard stats endpoint
	solanaRouter.HandleFunc("/dashboard/stats", GetDashboardStatsHandler(solanaService, logger)).Methods("GET")

	// Wallet routes
	walletRouter := solanaRouter.PathPrefix("/wallets").Subrouter()
	walletRouter.HandleFunc("/connect", wallets.ConnectWalletHandler(solanaService, logger)).Methods("POST")
	walletRouter.HandleFunc("/disconnect", wallets.DisconnectWalletHandler(solanaService, logger)).Methods("POST")
	walletRouter.HandleFunc("", wallets.GetWalletsHandler(solanaService, logger)).Methods("GET")
	walletRouter.HandleFunc("/balance", wallets.GetWalletBalanceHandler(solanaService, logger)).Methods("GET")
	walletRouter.HandleFunc("/refresh", wallets.RefreshBalanceHandler(solanaService, logger)).Methods("POST")

	// DeFi routes
	defiRouter := solanaRouter.PathPrefix("/defi").Subrouter()
	defiRouter.HandleFunc("/quote", defi.GetSwapQuoteHandler(solanaService, logger)).Methods("POST")
	defiRouter.HandleFunc("/swap", defi.ExecuteSwapHandler(solanaService, logger)).Methods("POST")
	defiRouter.HandleFunc("/liquidity", defi.AddLiquidityHandler(solanaService, logger)).Methods("POST")
	defiRouter.HandleFunc("/portfolio/{walletId}", defi.GetPortfolioHandler(solanaService, logger)).Methods("GET")
	defiRouter.HandleFunc("/tvl", defi.GetProtocolTVLHandler(solanaService, logger)).Methods("GET")

	// NFT routes
	nftRouter := solanaRouter.PathPrefix("/nft").Subrouter()
	nftRouter.HandleFunc("/explore", nft.ExploreNFTsHandler(solanaService, logger)).Methods("GET")
	nftRouter.HandleFunc("/collections", nft.GetCollectionsHandler(solanaService, logger)).Methods("GET")
	nftRouter.HandleFunc("/portfolio/{walletId}", nft.GetPortfolioHandler(solanaService, logger)).Methods("GET")
	nftRouter.HandleFunc("/list", nft.ListNFTHandler(solanaService, logger)).Methods("POST")
	nftRouter.HandleFunc("/buy", nft.BuyNFTHandler(solanaService, logger)).Methods("POST")
	nftRouter.HandleFunc("/metadata", nft.GetNFTMetadataHandler(solanaService, logger)).Methods("GET")
}

// GetSolanaStatsHandler returns Solana network statistics
func GetSolanaStatsHandler(logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// In a real implementation, this would fetch live data from:
		// - CoinGecko/CoinMarketCap for price data
		// - Solana RPC for network stats
		// - DeFiLlama for TVL data
		
		// For now, return mock data
		stats := SolanaStats{
			Price:          decimal.NewFromFloat(200.45),
			PriceChange24h: decimal.NewFromFloat(5.2),
			MarketCap:      decimal.NewFromInt(94500000000),  // $94.5B
			Volume24h:      decimal.NewFromInt(2800000000),   // $2.8B
			TPS:            65000,
			TotalAccounts:  180000000,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)

		logger.Info(ctx, "Solana stats retrieved", map[string]interface{}{
			"price":     stats.Price.String(),
			"market_cap": stats.MarketCap.String(),
		})
	}
}

// GetDashboardStatsHandler returns user dashboard statistics
func GetDashboardStatsHandler(solanaService *solana.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get user from context (would be set by auth middleware)
		// For now, return mock data
		stats := DashboardStats{
			TotalValue:      decimal.NewFromFloat(15420.50),
			TotalPnL:        decimal.NewFromFloat(1250.75),
			ActivePositions: 8,
			NFTCount:        12,
			SOLBalance:      decimal.NewFromFloat(75.25),
			TokenCount:      6,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)

		logger.Info(ctx, "Dashboard stats retrieved", map[string]interface{}{
			"total_value":       stats.TotalValue.String(),
			"active_positions":  stats.ActivePositions,
		})
	}
}

// HealthCheckHandler returns the health status of Solana services
func HealthCheckHandler(solanaService *solana.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Check Solana RPC connection
		health := map[string]interface{}{
			"status": "healthy",
			"services": map[string]interface{}{
				"solana_rpc": "connected",
				"database":   "connected",
				"defi":       "operational",
				"nft":        "operational",
			},
			"timestamp": r.Header.Get("X-Request-Time"),
		}

		// In a real implementation, you would:
		// 1. Check Solana RPC connectivity
		// 2. Check database connectivity
		// 3. Check external API availability (Jupiter, Magic Eden, etc.)
		// 4. Return appropriate status codes

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(health)

		logger.Info(ctx, "Health check completed", map[string]interface{}{
			"status": "healthy",
		})
	}
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
}

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// WriteErrorResponse writes an error response to the HTTP response writer
func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := ErrorResponse{
		Success: false,
		Error:   message,
		Code:    code,
	}
	
	json.NewEncoder(w).Encode(response)
}

// WriteSuccessResponse writes a success response to the HTTP response writer
func WriteSuccessResponse(w http.ResponseWriter, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	
	response := SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	}
	
	json.NewEncoder(w).Encode(response)
}

// CORSMiddleware adds CORS headers to responses
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(logger *observability.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			
			logger.Info(ctx, "HTTP request", map[string]interface{}{
				"method": r.Method,
				"path":   r.URL.Path,
				"query":  r.URL.RawQuery,
				"user_agent": r.UserAgent(),
			})

			next.ServeHTTP(w, r)
		})
	}
}

// RateLimitMiddleware implements basic rate limiting
func RateLimitMiddleware(requestsPerMinute int) func(http.Handler) http.Handler {
	// In a real implementation, this would use Redis or in-memory store
	// to track request counts per IP/user
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// For now, just pass through
			// Real implementation would check rate limits here
			next.ServeHTTP(w, r)
		})
	}
}
