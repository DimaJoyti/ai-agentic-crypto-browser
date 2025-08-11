package terminal

import (
	"context"
	"time"
)

// ServiceIntegrations provides access to all platform services
type ServiceIntegrations struct {
	AI      AIServiceClient
	Trading TradingServiceClient
	Web3    Web3ServiceClient
	Browser BrowserServiceClient
	Auth    AuthServiceClient
}

// AIServiceClient interface for AI service integration
type AIServiceClient interface {
	AnalyzeCryptocurrency(ctx context.Context, symbol string, timeframe string) (*AIAnalysisResult, error)
	PredictPrice(ctx context.Context, symbol string, horizon string) (*PricePrediction, error)
	ProcessNaturalLanguage(ctx context.Context, query string) (*NLPResponse, error)
	GetSentimentAnalysis(ctx context.Context, symbol string) (*SentimentAnalysis, error)
}

// TradingServiceClient interface for trading service integration
type TradingServiceClient interface {
	PlaceOrder(ctx context.Context, order *OrderRequest) (*OrderResponse, error)
	GetPortfolio(ctx context.Context, userID string) (*Portfolio, error)
	GetOrderHistory(ctx context.Context, userID string, limit int) ([]*Order, error)
	GetBalance(ctx context.Context, userID string) (*Balance, error)
	CancelOrder(ctx context.Context, orderID string) error
}

// Web3ServiceClient interface for Web3 service integration
type Web3ServiceClient interface {
	GetWalletBalance(ctx context.Context, address string, token string) (*TokenBalance, error)
	TransferTokens(ctx context.Context, request *TransferRequest) (*TransactionResult, error)
	GetDeFiPositions(ctx context.Context, address string) ([]*DeFiPosition, error)
	EstimateGas(ctx context.Context, transaction *Transaction) (*GasEstimate, error)
	GetTokenInfo(ctx context.Context, tokenAddress string) (*TokenInfo, error)
}

// BrowserServiceClient interface for browser automation
type BrowserServiceClient interface {
	NavigateToURL(ctx context.Context, url string) (*NavigationResult, error)
	ExtractData(ctx context.Context, url string, selectors []string) (*ExtractedData, error)
	TakeScreenshot(ctx context.Context, url string) (*Screenshot, error)
	ExecuteScript(ctx context.Context, script string) (*ScriptResult, error)
}

// AuthServiceClient interface for authentication
type AuthServiceClient interface {
	ValidateToken(ctx context.Context, token string) (*UserInfo, error)
	GetUserPermissions(ctx context.Context, userID string) ([]string, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)
}

// Data structures for service responses

// AI Service Types
type AIAnalysisResult struct {
	Symbol       string                 `json:"symbol"`
	Trend        string                 `json:"trend"`
	Confidence   float64                `json:"confidence"`
	Indicators   map[string]interface{} `json:"indicators"`
	Predictions  []PricePrediction      `json:"predictions"`
	Sentiment    SentimentAnalysis      `json:"sentiment"`
	Timestamp    time.Time              `json:"timestamp"`
}

type PricePrediction struct {
	Symbol      string    `json:"symbol"`
	Horizon     string    `json:"horizon"`
	Price       float64   `json:"price"`
	Confidence  float64   `json:"confidence"`
	Range       PriceRange `json:"range"`
	Timestamp   time.Time `json:"timestamp"`
}

type PriceRange struct {
	Low  float64 `json:"low"`
	High float64 `json:"high"`
}

type NLPResponse struct {
	Intent     string                 `json:"intent"`
	Entities   map[string]interface{} `json:"entities"`
	Response   string                 `json:"response"`
	Confidence float64                `json:"confidence"`
}

type SentimentAnalysis struct {
	Symbol    string  `json:"symbol"`
	Score     float64 `json:"score"`
	Label     string  `json:"label"`
	Sources   int     `json:"sources"`
	Timestamp time.Time `json:"timestamp"`
}

// Trading Service Types
type OrderRequest struct {
	Symbol   string  `json:"symbol"`
	Side     string  `json:"side"`
	Type     string  `json:"type"`
	Amount   float64 `json:"amount"`
	Price    float64 `json:"price,omitempty"`
	StopPrice float64 `json:"stop_price,omitempty"`
}

type OrderResponse struct {
	OrderID   string    `json:"order_id"`
	Status    string    `json:"status"`
	Symbol    string    `json:"symbol"`
	Side      string    `json:"side"`
	Amount    float64   `json:"amount"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}

type Portfolio struct {
	UserID      string            `json:"user_id"`
	TotalValue  float64           `json:"total_value"`
	Holdings    []Holding         `json:"holdings"`
	Performance PerformanceMetrics `json:"performance"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type Holding struct {
	Symbol   string  `json:"symbol"`
	Amount   float64 `json:"amount"`
	Value    float64 `json:"value"`
	AvgPrice float64 `json:"avg_price"`
	PnL      float64 `json:"pnl"`
	PnLPct   float64 `json:"pnl_pct"`
}

type PerformanceMetrics struct {
	TotalReturn    float64 `json:"total_return"`
	DailyReturn    float64 `json:"daily_return"`
	WeeklyReturn   float64 `json:"weekly_return"`
	MonthlyReturn  float64 `json:"monthly_return"`
	Volatility     float64 `json:"volatility"`
	SharpeRatio    float64 `json:"sharpe_ratio"`
}

type Order struct {
	ID        string    `json:"id"`
	Symbol    string    `json:"symbol"`
	Side      string    `json:"side"`
	Type      string    `json:"type"`
	Amount    float64   `json:"amount"`
	Price     float64   `json:"price"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Balance struct {
	UserID    string          `json:"user_id"`
	Balances  []AssetBalance  `json:"balances"`
	Total     float64         `json:"total"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type AssetBalance struct {
	Asset     string  `json:"asset"`
	Available float64 `json:"available"`
	Locked    float64 `json:"locked"`
	Total     float64 `json:"total"`
	Value     float64 `json:"value"`
}

// Web3 Service Types
type TokenBalance struct {
	Address   string  `json:"address"`
	Token     string  `json:"token"`
	Balance   float64 `json:"balance"`
	Decimals  int     `json:"decimals"`
	Value     float64 `json:"value"`
}

type TransferRequest struct {
	From     string  `json:"from"`
	To       string  `json:"to"`
	Token    string  `json:"token"`
	Amount   float64 `json:"amount"`
	GasLimit uint64  `json:"gas_limit,omitempty"`
	GasPrice uint64  `json:"gas_price,omitempty"`
}

type TransactionResult struct {
	Hash      string    `json:"hash"`
	Status    string    `json:"status"`
	GasUsed   uint64    `json:"gas_used"`
	Fee       float64   `json:"fee"`
	Timestamp time.Time `json:"timestamp"`
}

type DeFiPosition struct {
	Protocol string  `json:"protocol"`
	Type     string  `json:"type"`
	Token    string  `json:"token"`
	Amount   float64 `json:"amount"`
	Value    float64 `json:"value"`
	APY      float64 `json:"apy"`
}

type GasEstimate struct {
	GasLimit uint64  `json:"gas_limit"`
	GasPrice uint64  `json:"gas_price"`
	Fee      float64 `json:"fee"`
}

type TokenInfo struct {
	Address  string `json:"address"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
	Supply   string `json:"supply"`
}

type Transaction struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Value    string `json:"value"`
	Data     string `json:"data"`
	GasLimit uint64 `json:"gas_limit"`
}

// Browser Service Types
type NavigationResult struct {
	URL       string    `json:"url"`
	Title     string    `json:"title"`
	Status    int       `json:"status"`
	LoadTime  int64     `json:"load_time"`
	Timestamp time.Time `json:"timestamp"`
}

type ExtractedData struct {
	URL       string                 `json:"url"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

type Screenshot struct {
	URL       string    `json:"url"`
	ImageData []byte    `json:"image_data"`
	Format    string    `json:"format"`
	Timestamp time.Time `json:"timestamp"`
}

type ScriptResult struct {
	Result    interface{} `json:"result"`
	Error     string      `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// Auth Service Types
type UserInfo struct {
	UserID      string    `json:"user_id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Roles       []string  `json:"roles"`
	Permissions []string  `json:"permissions"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}
