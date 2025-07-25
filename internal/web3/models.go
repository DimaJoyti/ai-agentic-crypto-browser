package web3

import (
	"math/big"
	"time"

	"github.com/google/uuid"
)

// Wallet represents a connected cryptocurrency wallet
type Wallet struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Address   string    `json:"address" db:"address"`
	ChainID   int       `json:"chain_id" db:"chain_id"`
	WalletType string   `json:"wallet_type" db:"wallet_type"`
	IsPrimary bool      `json:"is_primary" db:"is_primary"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Transaction represents a blockchain transaction
type Transaction struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	UserID          uuid.UUID  `json:"user_id" db:"user_id"`
	WalletID        uuid.UUID  `json:"wallet_id" db:"wallet_id"`
	TxHash          string     `json:"tx_hash" db:"tx_hash"`
	ChainID         int        `json:"chain_id" db:"chain_id"`
	FromAddress     string     `json:"from_address" db:"from_address"`
	ToAddress       *string    `json:"to_address" db:"to_address"`
	Value           *big.Int   `json:"value" db:"value"`
	GasUsed         *big.Int   `json:"gas_used" db:"gas_used"`
	GasPrice        *big.Int   `json:"gas_price" db:"gas_price"`
	Status          TxStatus   `json:"status" db:"status"`
	BlockNumber     *big.Int   `json:"block_number" db:"block_number"`
	TransactionType string     `json:"transaction_type" db:"transaction_type"`
	Metadata        map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

// TxStatus represents transaction status
type TxStatus string

const (
	TxStatusPending   TxStatus = "pending"
	TxStatusConfirmed TxStatus = "confirmed"
	TxStatusFailed    TxStatus = "failed"
)

// DeFiPosition represents a DeFi position
type DeFiPosition struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	WalletID     uuid.UUID `json:"wallet_id" db:"wallet_id"`
	ProtocolName string    `json:"protocol_name" db:"protocol_name"`
	PositionType string    `json:"position_type" db:"position_type"`
	TokenAddress *string   `json:"token_address" db:"token_address"`
	TokenSymbol  *string   `json:"token_symbol" db:"token_symbol"`
	Amount       *big.Int  `json:"amount" db:"amount"`
	USDValue     *float64  `json:"usd_value" db:"usd_value"`
	APY          *float64  `json:"apy" db:"apy"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// WalletConnectRequest represents a wallet connection request
type WalletConnectRequest struct {
	Address    string `json:"address" validate:"required"`
	ChainID    int    `json:"chain_id" validate:"required"`
	WalletType string `json:"wallet_type" validate:"required"`
	Signature  string `json:"signature,omitempty"`
	Message    string `json:"message,omitempty"`
}

// WalletConnectResponse represents a wallet connection response
type WalletConnectResponse struct {
	Wallet  Wallet `json:"wallet"`
	Message string `json:"message"`
}

// BalanceRequest represents a balance query request
type BalanceRequest struct {
	WalletID     *uuid.UUID `json:"wallet_id,omitempty"`
	Address      *string    `json:"address,omitempty"`
	ChainID      *int       `json:"chain_id,omitempty"`
	TokenAddress *string    `json:"token_address,omitempty"`
}

// BalanceResponse represents a balance query response
type BalanceResponse struct {
	Address      string             `json:"address"`
	ChainID      int                `json:"chain_id"`
	NativeBalance *big.Int          `json:"native_balance"`
	TokenBalances []TokenBalance    `json:"token_balances,omitempty"`
	TotalUSDValue float64           `json:"total_usd_value"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// TokenBalance represents a token balance
type TokenBalance struct {
	TokenAddress string   `json:"token_address"`
	TokenSymbol  string   `json:"token_symbol"`
	TokenName    string   `json:"token_name"`
	Balance      *big.Int `json:"balance"`
	Decimals     int      `json:"decimals"`
	USDValue     float64  `json:"usd_value"`
}

// TransactionRequest represents a transaction creation request
type TransactionRequest struct {
	WalletID    uuid.UUID              `json:"wallet_id" validate:"required"`
	ToAddress   string                 `json:"to_address" validate:"required"`
	Value       *big.Int               `json:"value,omitempty"`
	Data        string                 `json:"data,omitempty"`
	GasLimit    *big.Int               `json:"gas_limit,omitempty"`
	GasPrice    *big.Int               `json:"gas_price,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// TransactionResponse represents a transaction response
type TransactionResponse struct {
	Transaction Transaction `json:"transaction"`
	TxHash      string      `json:"tx_hash"`
	Status      string      `json:"status"`
}

// DeFiProtocolRequest represents a DeFi protocol interaction request
type DeFiProtocolRequest struct {
	WalletID     uuid.UUID              `json:"wallet_id" validate:"required"`
	Protocol     string                 `json:"protocol" validate:"required"`
	Action       string                 `json:"action" validate:"required"`
	TokenAddress string                 `json:"token_address,omitempty"`
	Amount       *big.Int               `json:"amount,omitempty"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
}

// DeFiProtocolResponse represents a DeFi protocol interaction response
type DeFiProtocolResponse struct {
	Success     bool                   `json:"success"`
	TxHash      string                 `json:"tx_hash,omitempty"`
	Position    *DeFiPosition          `json:"position,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Error       string                 `json:"error,omitempty"`
}

// PriceRequest represents a price query request
type PriceRequest struct {
	TokenAddresses []string `json:"token_addresses,omitempty"`
	TokenSymbols   []string `json:"token_symbols,omitempty"`
	Currency       string   `json:"currency,omitempty"` // USD, EUR, etc.
}

// PriceResponse represents a price query response
type PriceResponse struct {
	Prices    map[string]TokenPrice `json:"prices"`
	Currency  string                `json:"currency"`
	Timestamp time.Time             `json:"timestamp"`
}

// TokenPrice represents a token price
type TokenPrice struct {
	Symbol           string    `json:"symbol"`
	Name             string    `json:"name"`
	Price            float64   `json:"price"`
	PriceChange24h   float64   `json:"price_change_24h"`
	PriceChangePerc  float64   `json:"price_change_percentage_24h"`
	MarketCap        float64   `json:"market_cap"`
	Volume24h        float64   `json:"volume_24h"`
	LastUpdated      time.Time `json:"last_updated"`
}

// NFTRequest represents an NFT query request
type NFTRequest struct {
	WalletID        *uuid.UUID `json:"wallet_id,omitempty"`
	Address         *string    `json:"address,omitempty"`
	ChainID         *int       `json:"chain_id,omitempty"`
	ContractAddress *string    `json:"contract_address,omitempty"`
	TokenID         *string    `json:"token_id,omitempty"`
}

// NFTResponse represents an NFT query response
type NFTResponse struct {
	NFTs     []NFT                  `json:"nfts"`
	Total    int                    `json:"total"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// NFT represents a non-fungible token
type NFT struct {
	ContractAddress string                 `json:"contract_address"`
	TokenID         string                 `json:"token_id"`
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	Image           string                 `json:"image"`
	Attributes      []NFTAttribute         `json:"attributes,omitempty"`
	Collection      string                 `json:"collection,omitempty"`
	Owner           string                 `json:"owner"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// NFTAttribute represents an NFT attribute
type NFTAttribute struct {
	TraitType string      `json:"trait_type"`
	Value     interface{} `json:"value"`
}

// ChainInfo represents blockchain network information
type ChainInfo struct {
	ChainID     int    `json:"chain_id"`
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	RpcURL      string `json:"rpc_url"`
	ExplorerURL string `json:"explorer_url"`
	IsTestnet   bool   `json:"is_testnet"`
}

// SupportedChains represents supported blockchain networks
var SupportedChains = map[int]ChainInfo{
	1: {
		ChainID:     1,
		Name:        "Ethereum Mainnet",
		Symbol:      "ETH",
		RpcURL:      "https://mainnet.infura.io/v3/",
		ExplorerURL: "https://etherscan.io",
		IsTestnet:   false,
	},
	137: {
		ChainID:     137,
		Name:        "Polygon Mainnet",
		Symbol:      "MATIC",
		RpcURL:      "https://polygon-mainnet.infura.io/v3/",
		ExplorerURL: "https://polygonscan.com",
		IsTestnet:   false,
	},
	42161: {
		ChainID:     42161,
		Name:        "Arbitrum One",
		Symbol:      "ETH",
		RpcURL:      "https://arbitrum-mainnet.infura.io/v3/",
		ExplorerURL: "https://arbiscan.io",
		IsTestnet:   false,
	},
	10: {
		ChainID:     10,
		Name:        "Optimism",
		Symbol:      "ETH",
		RpcURL:      "https://optimism-mainnet.infura.io/v3/",
		ExplorerURL: "https://optimistic.etherscan.io",
		IsTestnet:   false,
	},
	11155111: {
		ChainID:     11155111,
		Name:        "Sepolia Testnet",
		Symbol:      "ETH",
		RpcURL:      "https://sepolia.infura.io/v3/",
		ExplorerURL: "https://sepolia.etherscan.io",
		IsTestnet:   true,
	},
}

// WalletListRequest represents a request to list wallets
type WalletListRequest struct {
	UserID  uuid.UUID `json:"user_id"`
	ChainID *int      `json:"chain_id,omitempty"`
	Limit   int       `json:"limit,omitempty"`
	Offset  int       `json:"offset,omitempty"`
}

// WalletListResponse represents a response with wallet list
type WalletListResponse struct {
	Wallets []Wallet `json:"wallets"`
	Total   int      `json:"total"`
	HasMore bool     `json:"has_more"`
}

// TransactionListRequest represents a request to list transactions
type TransactionListRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	WalletID *uuid.UUID `json:"wallet_id,omitempty"`
	ChainID  *int       `json:"chain_id,omitempty"`
	Status   *TxStatus  `json:"status,omitempty"`
	Limit    int        `json:"limit,omitempty"`
	Offset   int        `json:"offset,omitempty"`
}

// TransactionListResponse represents a response with transaction list
type TransactionListResponse struct {
	Transactions []Transaction `json:"transactions"`
	Total        int           `json:"total"`
	HasMore      bool          `json:"has_more"`
}

// DeFiPositionListRequest represents a request to list DeFi positions
type DeFiPositionListRequest struct {
	UserID       uuid.UUID `json:"user_id"`
	WalletID     *uuid.UUID `json:"wallet_id,omitempty"`
	ProtocolName *string    `json:"protocol_name,omitempty"`
	IsActive     *bool      `json:"is_active,omitempty"`
	Limit        int        `json:"limit,omitempty"`
	Offset       int        `json:"offset,omitempty"`
}

// DeFiPositionListResponse represents a response with DeFi position list
type DeFiPositionListResponse struct {
	Positions []DeFiPosition `json:"positions"`
	Total     int            `json:"total"`
	HasMore   bool           `json:"has_more"`
}
