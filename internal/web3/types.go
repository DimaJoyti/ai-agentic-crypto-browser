package web3

import (
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Transaction status constants
const (
	TxStatusPending   = "pending"
	TxStatusConfirmed = "confirmed"
	TxStatusFailed    = "failed"
)

// Supported blockchain networks
var SupportedChains = map[int]string{
	1:     "ethereum",
	137:   "polygon",
	56:    "bsc",
	43114: "avalanche",
	250:   "fantom",
	42161: "arbitrum",
	10:    "optimism",
}

// Wallet represents a cryptocurrency wallet
type Wallet struct {
	ID         uuid.UUID              `json:"id"`
	UserID     uuid.UUID              `json:"user_id"`
	Address    string                 `json:"address"`
	Type       string                 `json:"type"`
	WalletType string                 `json:"wallet_type"`
	Name       string                 `json:"name"`
	ChainID    int                    `json:"chain_id"`
	Balance    decimal.Decimal        `json:"balance"`
	IsActive   bool                   `json:"is_active"`
	IsPrimary  bool                   `json:"is_primary"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// Transaction represents a blockchain transaction
type Transaction struct {
	ID              uuid.UUID              `json:"id"`
	UserID          uuid.UUID              `json:"user_id"`
	WalletID        uuid.UUID              `json:"wallet_id"`
	Hash            string                 `json:"hash"`
	TxHash          string                 `json:"tx_hash"`
	From            string                 `json:"from"`
	FromAddress     string                 `json:"from_address"`
	To              string                 `json:"to"`
	ToAddress       string                 `json:"to_address"`
	Value           *big.Int               `json:"value"`
	Data            string                 `json:"data"`
	GasLimit        uint64                 `json:"gas_limit"`
	GasPrice        *big.Int               `json:"gas_price"`
	GasUsed         uint64                 `json:"gas_used"`
	Status          string                 `json:"status"`
	BlockNumber     uint64                 `json:"block_number"`
	ChainID         int                    `json:"chain_id"`
	TransactionType string                 `json:"transaction_type"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// WalletConnectRequest represents a wallet connection request
type WalletConnectRequest struct {
	WalletType string                 `json:"wallet_type"`
	Address    string                 `json:"address"`
	ChainID    int                    `json:"chain_id"`
	UserID     uuid.UUID              `json:"user_id"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// WalletConnectResponse represents a wallet connection response
type WalletConnectResponse struct {
	WalletID uuid.UUID `json:"wallet_id"`
	Wallet   *Wallet   `json:"wallet"`
	Address  string    `json:"address"`
	ChainID  int       `json:"chain_id"`
	Success  bool      `json:"success"`
	Message  string    `json:"message"`
}

// BalanceRequest represents a balance query request
type BalanceRequest struct {
	WalletID uuid.UUID `json:"wallet_id"`
	Address  string    `json:"address"`
	ChainID  int       `json:"chain_id"`
	Token    string    `json:"token,omitempty"`
}

// BalanceResponse represents a balance query response
type BalanceResponse struct {
	Address       string                 `json:"address"`
	ChainID       int                    `json:"chain_id"`
	Token         string                 `json:"token"`
	Balance       decimal.Decimal        `json:"balance"`
	Symbol        string                 `json:"symbol"`
	Decimals      int                    `json:"decimals"`
	NativeBalance *big.Int               `json:"native_balance"`
	TokenBalances []TokenBalance         `json:"token_balances"`
	TotalUSDValue float64                `json:"total_usd_value"`
	Metadata      map[string]interface{} `json:"metadata"`
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
	WalletID  uuid.UUID              `json:"wallet_id"`
	From      string                 `json:"from"`
	To        string                 `json:"to"`
	ToAddress string                 `json:"to_address"`
	Value     *big.Int               `json:"value"`
	Data      string                 `json:"data"`
	GasLimit  uint64                 `json:"gas_limit"`
	GasPrice  *big.Int               `json:"gas_price"`
	ChainID   int                    `json:"chain_id"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// TransactionResponse represents a transaction creation response
type TransactionResponse struct {
	TransactionID uuid.UUID    `json:"transaction_id"`
	Transaction   *Transaction `json:"transaction"`
	Hash          string       `json:"hash"`
	TxHash        string       `json:"tx_hash"`
	Status        string       `json:"status"`
	Success       bool         `json:"success"`
	Message       string       `json:"message"`
}

// PriceRequest represents a price query request
type PriceRequest struct {
	Token    string `json:"token"`
	ChainID  int    `json:"chain_id"`
	Currency string `json:"currency"`
}

// PriceResponse represents a price query response
type PriceResponse struct {
	Token     string                 `json:"token"`
	ChainID   int                    `json:"chain_id"`
	Price     decimal.Decimal        `json:"price"`
	Currency  string                 `json:"currency"`
	Timestamp time.Time              `json:"timestamp"`
	Prices    map[string]TokenPrice  `json:"prices"`
}

// TokenPrice represents a token price
type TokenPrice struct {
	Token           string    `json:"token"`
	Symbol          string    `json:"symbol"`
	Name            string    `json:"name"`
	Price           float64   `json:"price"`
	PriceChange24h  float64   `json:"price_change_24h"`
	PriceChangePerc float64   `json:"price_change_perc"`
	MarketCap       float64   `json:"market_cap"`
	Volume24h       float64   `json:"volume_24h"`
	Currency        string    `json:"currency"`
	LastUpdated     time.Time `json:"last_updated"`
}

// TokenInfo represents information about a token
type TokenInfo struct {
	Address     string          `json:"address"`
	Symbol      string          `json:"symbol"`
	Name        string          `json:"name"`
	Decimals    int             `json:"decimals"`
	TotalSupply decimal.Decimal `json:"total_supply"`
	ChainID     int             `json:"chain_id"`
}

// SwapRequest represents a token swap request
type SwapRequest struct {
	FromToken string          `json:"from_token"`
	ToToken   string          `json:"to_token"`
	Amount    decimal.Decimal `json:"amount"`
	ChainID   int             `json:"chain_id"`
	Slippage  decimal.Decimal `json:"slippage"`
}

// SwapResponse represents a token swap response
type SwapResponse struct {
	FromToken     string          `json:"from_token"`
	ToToken       string          `json:"to_token"`
	AmountIn      decimal.Decimal `json:"amount_in"`
	AmountOut     decimal.Decimal `json:"amount_out"`
	Price         decimal.Decimal `json:"price"`
	PriceImpact   decimal.Decimal `json:"price_impact"`
	TransactionID uuid.UUID       `json:"transaction_id"`
	Success       bool            `json:"success"`
	Message       string          `json:"message"`
}

// StakeRequest represents a staking request
type StakeRequest struct {
	Token    string          `json:"token"`
	Amount   decimal.Decimal `json:"amount"`
	Protocol string          `json:"protocol"`
	ChainID  int             `json:"chain_id"`
}

// StakeResponse represents a staking response
type StakeResponse struct {
	Token         string          `json:"token"`
	Amount        decimal.Decimal `json:"amount"`
	Protocol      string          `json:"protocol"`
	APY           decimal.Decimal `json:"apy"`
	TransactionID uuid.UUID       `json:"transaction_id"`
	Success       bool            `json:"success"`
	Message       string          `json:"message"`
}

// LiquidityRequest represents a liquidity provision request
type LiquidityRequest struct {
	TokenA   string          `json:"token_a"`
	TokenB   string          `json:"token_b"`
	AmountA  decimal.Decimal `json:"amount_a"`
	AmountB  decimal.Decimal `json:"amount_b"`
	Protocol string          `json:"protocol"`
	ChainID  int             `json:"chain_id"`
}

// LiquidityResponse represents a liquidity provision response
type LiquidityResponse struct {
	TokenA        string          `json:"token_a"`
	TokenB        string          `json:"token_b"`
	AmountA       decimal.Decimal `json:"amount_a"`
	AmountB       decimal.Decimal `json:"amount_b"`
	LPTokens      decimal.Decimal `json:"lp_tokens"`
	Protocol      string          `json:"protocol"`
	TransactionID uuid.UUID       `json:"transaction_id"`
	Success       bool            `json:"success"`
	Message       string          `json:"message"`
}

// NetworkInfo represents blockchain network information
type NetworkInfo struct {
	ChainID     int    `json:"chain_id"`
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	RPC         string `json:"rpc"`
	Explorer    string `json:"explorer"`
	IsTestnet   bool   `json:"is_testnet"`
	IsSupported bool   `json:"is_supported"`
}

// DeFiProtocolRequest represents a DeFi protocol request
type DeFiProtocolRequest struct {
	WalletID uuid.UUID              `json:"wallet_id"`
	Protocol string                 `json:"protocol"`
	Action   string                 `json:"action"`
	ChainID  int                    `json:"chain_id"`
	Amount   decimal.Decimal        `json:"amount"`
	Token    string                 `json:"token"`
	Metadata map[string]interface{} `json:"metadata"`
}

// DeFiProtocolResponse represents a DeFi protocol response
type DeFiProtocolResponse struct {
	Protocol string                 `json:"protocol"`
	Action   string                 `json:"action"`
	Success  bool                   `json:"success"`
	Message  string                 `json:"message"`
	Error    string                 `json:"error,omitempty"`
	TxHash   string                 `json:"tx_hash,omitempty"`
	Position interface{}            `json:"position,omitempty"`
	Data     map[string]interface{} `json:"data"`
	Metadata map[string]interface{} `json:"metadata"`
}

// ContractInfo represents smart contract information
type ContractInfo struct {
	Address    string                 `json:"address"`
	Name       string                 `json:"name"`
	Symbol     string                 `json:"symbol"`
	Type       string                 `json:"type"`
	IsVerified bool                   `json:"is_verified"`
	CreatedAt  time.Time              `json:"created_at"`
	Creator    string                 `json:"creator"`
	ChainID    int                    `json:"chain_id"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// EventLog represents a blockchain event log
type EventLog struct {
	ID          uuid.UUID              `json:"id"`
	Address     string                 `json:"address"`
	Topics      []string               `json:"topics"`
	Data        string                 `json:"data"`
	BlockNumber uint64                 `json:"block_number"`
	TxHash      string                 `json:"tx_hash"`
	TxIndex     uint                   `json:"tx_index"`
	LogIndex    uint                   `json:"log_index"`
	ChainID     int                    `json:"chain_id"`
	Timestamp   time.Time              `json:"timestamp"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// AnalyticsData represents analytics data for portfolios and positions
type AnalyticsData struct {
	PortfolioID   uuid.UUID              `json:"portfolio_id"`
	TotalValue    decimal.Decimal        `json:"total_value"`
	TotalPnL      decimal.Decimal        `json:"total_pnl"`
	DailyPnL      decimal.Decimal        `json:"daily_pnl"`
	WeeklyPnL     decimal.Decimal        `json:"weekly_pnl"`
	MonthlyPnL    decimal.Decimal        `json:"monthly_pnl"`
	Volatility    decimal.Decimal        `json:"volatility"`
	SharpeRatio   decimal.Decimal        `json:"sharpe_ratio"`
	MaxDrawdown   decimal.Decimal        `json:"max_drawdown"`
	WinRate       decimal.Decimal        `json:"win_rate"`
	AvgWin        decimal.Decimal        `json:"avg_win"`
	AvgLoss       decimal.Decimal        `json:"avg_loss"`
	TotalTrades   int                    `json:"total_trades"`
	WinningTrades int                    `json:"winning_trades"`
	LosingTrades  int                    `json:"losing_trades"`
	LastUpdated   time.Time              `json:"last_updated"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// NotificationPreferences represents user notification preferences
type NotificationPreferences struct {
	UserID             uuid.UUID `json:"user_id"`
	EmailNotifications bool      `json:"email_notifications"`
	PushNotifications  bool      `json:"push_notifications"`
	SMSNotifications   bool      `json:"sms_notifications"`
	PriceAlerts        bool      `json:"price_alerts"`
	PortfolioAlerts    bool      `json:"portfolio_alerts"`
	SecurityAlerts     bool      `json:"security_alerts"`
	TradingAlerts      bool      `json:"trading_alerts"`
	NewsAlerts         bool      `json:"news_alerts"`
}

// APIKey represents an API key for external services
type APIKey struct {
	ID          uuid.UUID              `json:"id"`
	UserID      uuid.UUID              `json:"user_id"`
	Name        string                 `json:"name"`
	Key         string                 `json:"key"`
	Secret      string                 `json:"secret"`
	Service     string                 `json:"service"`
	Permissions []string               `json:"permissions"`
	IsActive    bool                   `json:"is_active"`
	ExpiresAt   *time.Time             `json:"expires_at"`
	CreatedAt   time.Time              `json:"created_at"`
	LastUsed    *time.Time             `json:"last_used"`
	Metadata    map[string]interface{} `json:"metadata"`
}
