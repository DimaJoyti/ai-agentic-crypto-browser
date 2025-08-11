package terminal

import (
	"context"
	"fmt"
	"time"
)

// NewMockServiceIntegrations creates mock service integrations for testing
func NewMockServiceIntegrations() *ServiceIntegrations {
	return &ServiceIntegrations{
		AI:      &MockAIServiceClient{},
		Trading: &MockTradingServiceClient{},
		Web3:    &MockWeb3ServiceClient{},
		Browser: &MockBrowserServiceClient{},
		Auth:    &MockAuthServiceClient{},
	}
}

// MockAIServiceClient provides mock AI service functionality
type MockAIServiceClient struct{}

func (m *MockAIServiceClient) AnalyzeCryptocurrency(ctx context.Context, symbol string, timeframe string) (*AIAnalysisResult, error) {
	return &AIAnalysisResult{
		Symbol:     symbol,
		Trend:      "bullish",
		Confidence: 0.82,
		Indicators: map[string]interface{}{
			"rsi":  65.4,
			"macd": 0.8,
			"sma":  45200.0,
		},
		Predictions: []PricePrediction{
			{
				Symbol:     symbol,
				Horizon:    "1h",
				Price:      45500.0,
				Confidence: 0.85,
				Range:      PriceRange{Low: 45200.0, High: 45800.0},
				Timestamp:  time.Now(),
			},
		},
		Sentiment: SentimentAnalysis{
			Symbol:    symbol,
			Score:     0.65,
			Label:     "positive",
			Sources:   150,
			Timestamp: time.Now(),
		},
		Timestamp: time.Now(),
	}, nil
}

func (m *MockAIServiceClient) PredictPrice(ctx context.Context, symbol string, horizon string) (*PricePrediction, error) {
	basePrice := 45000.0
	if symbol == "ETH" {
		basePrice = 3200.0
	}
	
	var multiplier float64
	switch horizon {
	case "1h":
		multiplier = 1.002
	case "4h":
		multiplier = 1.015
	case "1d":
		multiplier = 1.035
	case "1w":
		multiplier = 1.08
	default:
		multiplier = 1.0
	}
	
	predictedPrice := basePrice * multiplier
	
	return &PricePrediction{
		Symbol:     symbol,
		Horizon:    horizon,
		Price:      predictedPrice,
		Confidence: 0.78,
		Range: PriceRange{
			Low:  predictedPrice * 0.95,
			High: predictedPrice * 1.05,
		},
		Timestamp: time.Now(),
	}, nil
}

func (m *MockAIServiceClient) ProcessNaturalLanguage(ctx context.Context, query string) (*NLPResponse, error) {
	return &NLPResponse{
		Intent:     "price_query",
		Entities:   map[string]interface{}{"symbol": "BTC"},
		Response:   "Bitcoin is currently trading at $45,000 with a bullish trend.",
		Confidence: 0.92,
	}, nil
}

func (m *MockAIServiceClient) GetSentimentAnalysis(ctx context.Context, symbol string) (*SentimentAnalysis, error) {
	return &SentimentAnalysis{
		Symbol:    symbol,
		Score:     0.65,
		Label:     "positive",
		Sources:   150,
		Timestamp: time.Now(),
	}, nil
}

// MockTradingServiceClient provides mock trading functionality
type MockTradingServiceClient struct{}

func (m *MockTradingServiceClient) PlaceOrder(ctx context.Context, order *OrderRequest) (*OrderResponse, error) {
	return &OrderResponse{
		OrderID:   fmt.Sprintf("ORD-%d", time.Now().Unix()),
		Status:    "filled",
		Symbol:    order.Symbol,
		Side:      order.Side,
		Amount:    order.Amount,
		Price:     order.Price,
		Timestamp: time.Now(),
	}, nil
}

func (m *MockTradingServiceClient) GetPortfolio(ctx context.Context, userID string) (*Portfolio, error) {
	return &Portfolio{
		UserID:     userID,
		TotalValue: 125450.0,
		Holdings: []Holding{
			{Symbol: "BTC", Amount: 2.5, Value: 112500.0, AvgPrice: 42000.0, PnL: 7500.0, PnLPct: 7.14},
			{Symbol: "ETH", Amount: 4.0, Value: 12800.0, AvgPrice: 3000.0, PnL: 800.0, PnLPct: 6.67},
			{Symbol: "ADA", Amount: 150.0, Value: 187.5, AvgPrice: 1.30, PnL: -7.5, PnLPct: -3.85},
		},
		Performance: PerformanceMetrics{
			TotalReturn:   8.5,
			DailyReturn:   1.9,
			WeeklyReturn:  5.2,
			MonthlyReturn: 12.3,
			Volatility:    0.25,
			SharpeRatio:   1.8,
		},
		UpdatedAt: time.Now(),
	}, nil
}

func (m *MockTradingServiceClient) GetOrderHistory(ctx context.Context, userID string, limit int) ([]*Order, error) {
	return []*Order{
		{
			ID:        "ORD-001",
			Symbol:    "BTC",
			Side:      "buy",
			Type:      "market",
			Amount:    0.1,
			Price:     44500.0,
			Status:    "filled",
			CreatedAt: time.Now().Add(-2 * time.Hour),
			UpdatedAt: time.Now().Add(-2 * time.Hour),
		},
		{
			ID:        "ORD-002",
			Symbol:    "ETH",
			Side:      "sell",
			Type:      "limit",
			Amount:    1.0,
			Price:     3300.0,
			Status:    "pending",
			CreatedAt: time.Now().Add(-1 * time.Hour),
			UpdatedAt: time.Now().Add(-1 * time.Hour),
		},
	}, nil
}

func (m *MockTradingServiceClient) GetBalance(ctx context.Context, userID string) (*Balance, error) {
	return &Balance{
		UserID: userID,
		Balances: []AssetBalance{
			{Asset: "USD", Available: 25000.0, Locked: 5000.0, Total: 30000.0, Value: 30000.0},
			{Asset: "BTC", Available: 2.5, Locked: 0.5, Total: 3.0, Value: 135000.0},
			{Asset: "ETH", Available: 8.0, Locked: 2.0, Total: 10.0, Value: 32000.0},
		},
		Total:     197000.0,
		UpdatedAt: time.Now(),
	}, nil
}

func (m *MockTradingServiceClient) CancelOrder(ctx context.Context, orderID string) error {
	return nil
}

// MockWeb3ServiceClient provides mock Web3 functionality
type MockWeb3ServiceClient struct{}

func (m *MockWeb3ServiceClient) GetWalletBalance(ctx context.Context, address string, token string) (*TokenBalance, error) {
	return &TokenBalance{
		Address:  address,
		Token:    token,
		Balance:  1.5,
		Decimals: 18,
		Value:    4800.0,
	}, nil
}

func (m *MockWeb3ServiceClient) TransferTokens(ctx context.Context, request *TransferRequest) (*TransactionResult, error) {
	return &TransactionResult{
		Hash:      "0x1234567890abcdef",
		Status:    "confirmed",
		GasUsed:   21000,
		Fee:       0.002,
		Timestamp: time.Now(),
	}, nil
}

func (m *MockWeb3ServiceClient) GetDeFiPositions(ctx context.Context, address string) ([]*DeFiPosition, error) {
	return []*DeFiPosition{
		{Protocol: "Uniswap", Type: "liquidity", Token: "ETH/USDC", Amount: 1000.0, Value: 1050.0, APY: 12.5},
		{Protocol: "Compound", Type: "lending", Token: "USDC", Amount: 5000.0, Value: 5125.0, APY: 8.2},
	}, nil
}

func (m *MockWeb3ServiceClient) EstimateGas(ctx context.Context, transaction *Transaction) (*GasEstimate, error) {
	return &GasEstimate{
		GasLimit: 21000,
		GasPrice: 20000000000, // 20 gwei
		Fee:      0.0042,
	}, nil
}

func (m *MockWeb3ServiceClient) GetTokenInfo(ctx context.Context, tokenAddress string) (*TokenInfo, error) {
	return &TokenInfo{
		Address:  tokenAddress,
		Name:     "Mock Token",
		Symbol:   "MOCK",
		Decimals: 18,
		Supply:   "1000000000000000000000000",
	}, nil
}

// MockBrowserServiceClient provides mock browser functionality
type MockBrowserServiceClient struct{}

func (m *MockBrowserServiceClient) NavigateToURL(ctx context.Context, url string) (*NavigationResult, error) {
	return &NavigationResult{
		URL:       url,
		Title:     "Mock Page Title",
		Status:    200,
		LoadTime:  1500,
		Timestamp: time.Now(),
	}, nil
}

func (m *MockBrowserServiceClient) ExtractData(ctx context.Context, url string, selectors []string) (*ExtractedData, error) {
	return &ExtractedData{
		URL:       url,
		Data:      map[string]interface{}{"price": "$45,000", "volume": "$28.5B"},
		Timestamp: time.Now(),
	}, nil
}

func (m *MockBrowserServiceClient) TakeScreenshot(ctx context.Context, url string) (*Screenshot, error) {
	return &Screenshot{
		URL:       url,
		ImageData: []byte("mock-image-data"),
		Format:    "png",
		Timestamp: time.Now(),
	}, nil
}

func (m *MockBrowserServiceClient) ExecuteScript(ctx context.Context, script string) (*ScriptResult, error) {
	return &ScriptResult{
		Result:    "Script executed successfully",
		Timestamp: time.Now(),
	}, nil
}

// MockAuthServiceClient provides mock authentication functionality
type MockAuthServiceClient struct{}

func (m *MockAuthServiceClient) ValidateToken(ctx context.Context, token string) (*UserInfo, error) {
	return &UserInfo{
		UserID:      "user-123",
		Username:    "testuser",
		Email:       "test@example.com",
		Roles:       []string{"user", "trader"},
		Permissions: []string{"trade", "view_portfolio", "use_terminal"},
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}, nil
}

func (m *MockAuthServiceClient) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	return []string{"trade", "view_portfolio", "use_terminal", "admin"}, nil
}

func (m *MockAuthServiceClient) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	return &TokenResponse{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}, nil
}
