package ai

import (
	"context"
	"math/rand"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/shopspring/decimal"
)

// MarketAnalyzer provides market analysis and insights
type MarketAnalyzer struct {
	logger      *observability.Logger
	dataCache   map[string]interface{}
	lastUpdated time.Time
}

// MarketData represents comprehensive market data
type MarketData struct {
	Timestamp       time.Time              `json:"timestamp"`
	GlobalMarketCap decimal.Decimal        `json:"global_market_cap"`
	TotalVolume24h  decimal.Decimal        `json:"total_volume_24h"`
	BTCDominance    decimal.Decimal        `json:"btc_dominance"`
	ETHDominance    decimal.Decimal        `json:"eth_dominance"`
	FearGreedIndex  int                    `json:"fear_greed_index"`
	TopTokens       []TokenData            `json:"top_tokens"`
	TrendingTokens  []TokenData            `json:"trending_tokens"`
	MarketSentiment MarketSentimentData    `json:"market_sentiment"`
	TechnicalData   TechnicalAnalysisData  `json:"technical_data"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// TokenData represents individual token data
type TokenData struct {
	Symbol      string          `json:"symbol"`
	Name        string          `json:"name"`
	Price       decimal.Decimal `json:"price"`
	MarketCap   decimal.Decimal `json:"market_cap"`
	Volume24h   decimal.Decimal `json:"volume_24h"`
	Change1h    decimal.Decimal `json:"change_1h"`
	Change24h   decimal.Decimal `json:"change_24h"`
	Change7d    decimal.Decimal `json:"change_7d"`
	Rank        int             `json:"rank"`
	LastUpdated time.Time       `json:"last_updated"`
}

// MarketSentimentData represents market sentiment analysis
type MarketSentimentData struct {
	OverallSentiment string            `json:"overall_sentiment"`
	BullishSignals   int               `json:"bullish_signals"`
	BearishSignals   int               `json:"bearish_signals"`
	NeutralSignals   int               `json:"neutral_signals"`
	SentimentScore   decimal.Decimal   `json:"sentiment_score"`
	Sources          []SentimentSource `json:"sources"`
	LastUpdated      time.Time         `json:"last_updated"`
}

// SentimentSource represents a source of sentiment data
type SentimentSource struct {
	Source    string          `json:"source"`
	Sentiment string          `json:"sentiment"`
	Score     decimal.Decimal `json:"score"`
	Weight    decimal.Decimal `json:"weight"`
}

// TechnicalAnalysisData represents technical analysis indicators
type TechnicalAnalysisData struct {
	RSI            decimal.Decimal `json:"rsi"`
	MACD           decimal.Decimal `json:"macd"`
	MACDSignal     decimal.Decimal `json:"macd_signal"`
	BollingerUpper decimal.Decimal `json:"bollinger_upper"`
	BollingerLower decimal.Decimal `json:"bollinger_lower"`
	SMA20          decimal.Decimal `json:"sma_20"`
	SMA50          decimal.Decimal `json:"sma_50"`
	SMA200         decimal.Decimal `json:"sma_200"`
	EMA12          decimal.Decimal `json:"ema_12"`
	EMA26          decimal.Decimal `json:"ema_26"`
	StochasticK    decimal.Decimal `json:"stochastic_k"`
	StochasticD    decimal.Decimal `json:"stochastic_d"`
	WilliamsR      decimal.Decimal `json:"williams_r"`
	CCI            decimal.Decimal `json:"cci"`
	ADX            decimal.Decimal `json:"adx"`
	LastUpdated    time.Time       `json:"last_updated"`
}

// MarketTrend represents market trend analysis
type MarketTrend struct {
	Direction   string          `json:"direction"`
	Strength    string          `json:"strength"`
	Confidence  decimal.Decimal `json:"confidence"`
	Duration    string          `json:"duration"`
	Resistance  decimal.Decimal `json:"resistance,omitempty"`
	Support     decimal.Decimal `json:"support,omitempty"`
	Target      decimal.Decimal `json:"target,omitempty"`
	LastUpdated time.Time       `json:"last_updated"`
}

// NewMarketAnalyzer creates a new market analyzer
func NewMarketAnalyzer(logger *observability.Logger) *MarketAnalyzer {
	return &MarketAnalyzer{
		logger:      logger,
		dataCache:   make(map[string]interface{}),
		lastUpdated: time.Time{},
	}
}

// GetMarketContext returns current market context
func (m *MarketAnalyzer) GetMarketContext(ctx context.Context) (*MarketContext, error) {
	// In a real implementation, this would fetch data from multiple sources
	// For now, we'll generate mock data that represents realistic market conditions

	marketContext := &MarketContext{
		MarketTrend:     m.generateMarketTrend(),
		Volatility:      m.generateVolatility(),
		TopMovers:       m.generateTopMovers(),
		MarketSentiment: m.generateMarketSentiment(),
		KeyEvents:       m.generateKeyEvents(),
		LastUpdated:     time.Now(),
	}

	m.logger.Info(ctx, "Market context generated", map[string]interface{}{
		"trend":      marketContext.MarketTrend,
		"volatility": marketContext.Volatility,
		"sentiment":  marketContext.MarketSentiment,
	})

	return marketContext, nil
}

// GetMarketData returns comprehensive market data
func (m *MarketAnalyzer) GetMarketData(ctx context.Context) (*MarketData, error) {
	// Check cache first
	if time.Since(m.lastUpdated) < 5*time.Minute {
		if cached, exists := m.dataCache["market_data"]; exists {
			return cached.(*MarketData), nil
		}
	}

	// Generate fresh market data
	marketData := &MarketData{
		Timestamp:       time.Now(),
		GlobalMarketCap: decimal.NewFromFloat(1200000000000), // $1.2T
		TotalVolume24h:  decimal.NewFromFloat(45000000000),   // $45B
		BTCDominance:    decimal.NewFromFloat(42.5),          // 42.5%
		ETHDominance:    decimal.NewFromFloat(18.3),          // 18.3%
		FearGreedIndex:  m.generateFearGreedIndex(),
		TopTokens:       m.generateTopTokens(),
		TrendingTokens:  m.generateTrendingTokens(),
		MarketSentiment: m.generateMarketSentimentData(),
		TechnicalData:   m.generateTechnicalData(),
		Metadata:        make(map[string]interface{}),
	}

	// Cache the data
	m.dataCache["market_data"] = marketData
	m.lastUpdated = time.Now()

	return marketData, nil
}

// AnalyzeTrend analyzes market trend for a specific token
func (m *MarketAnalyzer) AnalyzeTrend(ctx context.Context, symbol string) (*MarketTrend, error) {
	// In a real implementation, this would analyze historical price data
	trend := &MarketTrend{
		Direction:   m.generateTrendDirection(),
		Strength:    m.generateTrendStrength(),
		Confidence:  decimal.NewFromFloat(0.75 + rand.Float64()*0.2), // 75-95%
		Duration:    m.generateTrendDuration(),
		LastUpdated: time.Now(),
	}

	// Add support/resistance levels for major tokens
	if symbol == "BTC" || symbol == "ETH" {
		trend.Support = decimal.NewFromFloat(40000 + rand.Float64()*5000)
		trend.Resistance = decimal.NewFromFloat(50000 + rand.Float64()*10000)
		trend.Target = decimal.NewFromFloat(55000 + rand.Float64()*15000)
	}

	return trend, nil
}

// GetSentimentAnalysis returns sentiment analysis for the market
func (m *MarketAnalyzer) GetSentimentAnalysis(ctx context.Context) (*MarketSentimentData, error) {
	sentiment := m.generateMarketSentimentData()
	return &sentiment, nil
}

// Helper methods for generating mock data

func (m *MarketAnalyzer) generateMarketTrend() string {
	trends := []string{"bullish", "bearish", "sideways", "volatile"}
	return trends[rand.Intn(len(trends))]
}

func (m *MarketAnalyzer) generateVolatility() string {
	volatilities := []string{"low", "medium", "high"}
	return volatilities[rand.Intn(len(volatilities))]
}

func (m *MarketAnalyzer) generateMarketSentiment() string {
	sentiments := []string{"bullish", "bearish", "neutral", "mixed"}
	return sentiments[rand.Intn(len(sentiments))]
}

func (m *MarketAnalyzer) generateTopMovers() []TokenMovement {
	tokens := []string{"BTC", "ETH", "ADA", "SOL", "MATIC", "LINK", "UNI", "AAVE"}
	movers := make([]TokenMovement, 0, 5)

	for i := 0; i < 5; i++ {
		symbol := tokens[rand.Intn(len(tokens))]
		basePrice := 100.0
		if symbol == "BTC" {
			basePrice = 45000.0
		} else if symbol == "ETH" {
			basePrice = 2500.0
		}

		change := (rand.Float64() - 0.5) * 20 // -10% to +10%
		price := basePrice * (1 + change/100)

		movers = append(movers, TokenMovement{
			Symbol:     symbol,
			Price:      decimal.NewFromFloat(price),
			Change24h:  decimal.NewFromFloat(change),
			ChangePerc: decimal.NewFromFloat(change),
			Volume24h:  decimal.NewFromFloat(rand.Float64() * 1000000000), // Up to $1B
		})
	}

	return movers
}

func (m *MarketAnalyzer) generateKeyEvents() []MarketEvent {
	events := []MarketEvent{
		{
			Title:       "Federal Reserve Interest Rate Decision",
			Description: "The Fed announced a 0.25% rate increase, impacting crypto markets",
			Impact:      "bearish",
			Timestamp:   time.Now().Add(-2 * time.Hour),
			Source:      "Federal Reserve",
		},
		{
			Title:       "Major DeFi Protocol Upgrade",
			Description: "Uniswap V4 announcement drives positive sentiment",
			Impact:      "bullish",
			Timestamp:   time.Now().Add(-6 * time.Hour),
			Source:      "Uniswap Labs",
		},
		{
			Title:       "Institutional Bitcoin Purchase",
			Description: "MicroStrategy adds $100M worth of Bitcoin to treasury",
			Impact:      "bullish",
			Timestamp:   time.Now().Add(-12 * time.Hour),
			Source:      "MicroStrategy",
		},
	}

	return events
}

func (m *MarketAnalyzer) generateFearGreedIndex() int {
	return 25 + rand.Intn(50) // 25-75 range
}

func (m *MarketAnalyzer) generateTopTokens() []TokenData {
	tokens := []TokenData{
		{
			Symbol:      "BTC",
			Name:        "Bitcoin",
			Price:       decimal.NewFromFloat(45000 + rand.Float64()*5000),
			MarketCap:   decimal.NewFromFloat(850000000000),
			Volume24h:   decimal.NewFromFloat(15000000000),
			Change24h:   decimal.NewFromFloat((rand.Float64() - 0.5) * 10),
			Rank:        1,
			LastUpdated: time.Now(),
		},
		{
			Symbol:      "ETH",
			Name:        "Ethereum",
			Price:       decimal.NewFromFloat(2500 + rand.Float64()*500),
			MarketCap:   decimal.NewFromFloat(300000000000),
			Volume24h:   decimal.NewFromFloat(8000000000),
			Change24h:   decimal.NewFromFloat((rand.Float64() - 0.5) * 12),
			Rank:        2,
			LastUpdated: time.Now(),
		},
		{
			Symbol:      "USDT",
			Name:        "Tether",
			Price:       decimal.NewFromFloat(1.0),
			MarketCap:   decimal.NewFromFloat(85000000000),
			Volume24h:   decimal.NewFromFloat(25000000000),
			Change24h:   decimal.NewFromFloat((rand.Float64() - 0.5) * 0.5),
			Rank:        3,
			LastUpdated: time.Now(),
		},
	}

	return tokens
}

func (m *MarketAnalyzer) generateTrendingTokens() []TokenData {
	// Similar to top tokens but with higher volatility
	return m.generateTopTokens()
}

func (m *MarketAnalyzer) generateMarketSentimentData() MarketSentimentData {
	bullish := rand.Intn(50) + 10
	bearish := rand.Intn(50) + 10
	neutral := 100 - bullish - bearish

	overall := "neutral"
	if bullish > bearish+10 {
		overall = "bullish"
	} else if bearish > bullish+10 {
		overall = "bearish"
	}

	return MarketSentimentData{
		OverallSentiment: overall,
		BullishSignals:   bullish,
		BearishSignals:   bearish,
		NeutralSignals:   neutral,
		SentimentScore:   decimal.NewFromFloat(float64(bullish-bearish) / 100.0),
		Sources: []SentimentSource{
			{Source: "Twitter", Sentiment: overall, Score: decimal.NewFromFloat(0.6), Weight: decimal.NewFromFloat(0.3)},
			{Source: "Reddit", Sentiment: overall, Score: decimal.NewFromFloat(0.7), Weight: decimal.NewFromFloat(0.2)},
			{Source: "News", Sentiment: overall, Score: decimal.NewFromFloat(0.8), Weight: decimal.NewFromFloat(0.5)},
		},
		LastUpdated: time.Now(),
	}
}

func (m *MarketAnalyzer) generateTechnicalData() TechnicalAnalysisData {
	return TechnicalAnalysisData{
		RSI:            decimal.NewFromFloat(30 + rand.Float64()*40), // 30-70
		MACD:           decimal.NewFromFloat((rand.Float64() - 0.5) * 100),
		MACDSignal:     decimal.NewFromFloat((rand.Float64() - 0.5) * 100),
		BollingerUpper: decimal.NewFromFloat(2600),
		BollingerLower: decimal.NewFromFloat(2400),
		SMA20:          decimal.NewFromFloat(2500),
		SMA50:          decimal.NewFromFloat(2480),
		SMA200:         decimal.NewFromFloat(2450),
		EMA12:          decimal.NewFromFloat(2510),
		EMA26:          decimal.NewFromFloat(2490),
		StochasticK:    decimal.NewFromFloat(rand.Float64() * 100),
		StochasticD:    decimal.NewFromFloat(rand.Float64() * 100),
		WilliamsR:      decimal.NewFromFloat(-rand.Float64() * 100),
		CCI:            decimal.NewFromFloat((rand.Float64() - 0.5) * 200),
		ADX:            decimal.NewFromFloat(20 + rand.Float64()*60),
		LastUpdated:    time.Now(),
	}
}

func (m *MarketAnalyzer) generateTrendDirection() string {
	directions := []string{"uptrend", "downtrend", "sideways"}
	return directions[rand.Intn(len(directions))]
}

func (m *MarketAnalyzer) generateTrendStrength() string {
	strengths := []string{"weak", "moderate", "strong"}
	return strengths[rand.Intn(len(strengths))]
}

func (m *MarketAnalyzer) generateTrendDuration() string {
	durations := []string{"short-term", "medium-term", "long-term"}
	return durations[rand.Intn(len(durations))]
}
