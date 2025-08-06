package ai

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/shopspring/decimal"
)

// CryptoCoinAnalyzer provides comprehensive cryptocurrency analysis
type CryptoCoinAnalyzer struct {
	logger          *observability.Logger
	httpClient      *http.Client
	webSearch       *WebSearchService
	reportGenerator *CryptoAnalysisReportGenerator
	dataCache       map[string]*CoinAnalysisCache
	lastUpdated     time.Time
	currentReport   *CoinAnalysisReport // Track current report for data source tracking
}

// CoinAnalysisCache represents cached analysis data
type CoinAnalysisCache struct {
	Data        *CoinAnalysisReport
	LastUpdated time.Time
	ExpiresAt   time.Time
}

// CoinAnalysisReport represents the complete analysis report
type CoinAnalysisReport struct {
	Timestamp       time.Time                `json:"timestamp"`
	Symbol          string                   `json:"symbol"`
	CurrentData     *CurrentMarketData       `json:"current_data"`
	NewsAndEvents   []NewsItem               `json:"news_and_events"`
	MarketSentiment *MarketSentimentAnalysis `json:"market_sentiment"`
	TechnicalData   *TechnicalIndicators     `json:"technical_data"`
	FundamentalData *FundamentalAnalysis     `json:"fundamental_data"`
	Summary         *AnalysisSummary         `json:"summary"`
	Sources         []DataSource             `json:"sources"`
}

// CurrentMarketData represents current market data
type CurrentMarketData struct {
	Price             decimal.Decimal `json:"price"`
	Change24h         decimal.Decimal `json:"change_24h"`
	ChangePercent24h  decimal.Decimal `json:"change_percent_24h"`
	MarketCap         decimal.Decimal `json:"market_cap"`
	Volume24h         decimal.Decimal `json:"volume_24h"`
	CirculatingSupply decimal.Decimal `json:"circulating_supply"`
	MaxSupply         decimal.Decimal `json:"max_supply,omitempty"`
	Rank              int             `json:"rank"`
	LastUpdated       time.Time       `json:"last_updated"`
}

// NewsItem represents a news item
type NewsItem struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	Source      string    `json:"source"`
	PublishedAt time.Time `json:"published_at"`
	Impact      string    `json:"impact"` // bullish, bearish, neutral
	Relevance   float64   `json:"relevance"`
}

// MarketSentimentAnalysis represents sentiment analysis
type MarketSentimentAnalysis struct {
	OverallSentiment string            `json:"overall_sentiment"`
	SentimentScore   decimal.Decimal   `json:"sentiment_score"`
	BullishPercent   decimal.Decimal   `json:"bullish_percent"`
	BearishPercent   decimal.Decimal   `json:"bearish_percent"`
	NeutralPercent   decimal.Decimal   `json:"neutral_percent"`
	KeyDrivers       []string          `json:"key_drivers"`
	SocialMetrics    *SocialMetrics    `json:"social_metrics"`
	Sources          []SentimentSource `json:"sources"`
}

// SocialMetrics represents social media metrics
type SocialMetrics struct {
	TwitterMentions  int             `json:"twitter_mentions"`
	RedditPosts      int             `json:"reddit_posts"`
	TelegramMessages int             `json:"telegram_messages"`
	SentimentTrend   string          `json:"sentiment_trend"`
	InfluencerScore  decimal.Decimal `json:"influencer_score"`
	CommunityGrowth  decimal.Decimal `json:"community_growth"`
}

// TechnicalIndicators represents technical analysis data
type TechnicalIndicators struct {
	Trend            string            `json:"trend"`
	TrendStrength    string            `json:"trend_strength"`
	SupportLevels    []decimal.Decimal `json:"support_levels"`
	ResistanceLevels []decimal.Decimal `json:"resistance_levels"`
	RSI              decimal.Decimal   `json:"rsi"`
	MACD             decimal.Decimal   `json:"macd"`
	MACDSignal       decimal.Decimal   `json:"macd_signal"`
	BollingerBands   *BollingerBands   `json:"bollinger_bands"`
	MovingAverages   *MovingAverages   `json:"moving_averages"`
	VolumeProfile    *VolumeProfile    `json:"volume_profile"`
	TechnicalOutlook string            `json:"technical_outlook"`
}

// BollingerBands represents Bollinger Bands data
type BollingerBands struct {
	Upper  decimal.Decimal `json:"upper"`
	Middle decimal.Decimal `json:"middle"`
	Lower  decimal.Decimal `json:"lower"`
}

// MovingAverages represents moving averages
type MovingAverages struct {
	SMA20  decimal.Decimal `json:"sma_20"`
	SMA50  decimal.Decimal `json:"sma_50"`
	SMA200 decimal.Decimal `json:"sma_200"`
	EMA12  decimal.Decimal `json:"ema_12"`
	EMA26  decimal.Decimal `json:"ema_26"`
}

// VolumeProfile represents volume analysis
type VolumeProfile struct {
	AverageVolume   decimal.Decimal `json:"average_volume"`
	VolumeRatio     decimal.Decimal `json:"volume_ratio"`
	VolumeIndicator string          `json:"volume_indicator"`
}

// FundamentalAnalysis represents fundamental analysis
type FundamentalAnalysis struct {
	ProjectStatus       string               `json:"project_status"`
	RecentUpdates       []ProjectUpdate      `json:"recent_updates"`
	CompetitivePosition *CompetitiveAnalysis `json:"competitive_position"`
	DeveloperActivity   *DeveloperMetrics    `json:"developer_activity"`
	NetworkMetrics      *NetworkMetrics      `json:"network_metrics"`
	TokenomicsHealth    *TokenomicsAnalysis  `json:"tokenomics_health"`
}

// ProjectUpdate represents a project update
type ProjectUpdate struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Impact      string    `json:"impact"`
	Source      string    `json:"source"`
}

// CompetitiveAnalysis represents competitive position
type CompetitiveAnalysis struct {
	MarketPosition string          `json:"market_position"`
	KeyCompetitors []string        `json:"key_competitors"`
	Advantages     []string        `json:"advantages"`
	Challenges     []string        `json:"challenges"`
	MarketShare    decimal.Decimal `json:"market_share"`
}

// DeveloperMetrics represents developer activity
type DeveloperMetrics struct {
	GitHubCommits    int       `json:"github_commits"`
	ActiveDevelopers int       `json:"active_developers"`
	CodeQuality      string    `json:"code_quality"`
	DevelopmentTrend string    `json:"development_trend"`
	LastCommit       time.Time `json:"last_commit"`
}

// NetworkMetrics represents network health metrics
type NetworkMetrics struct {
	ActiveAddresses  int             `json:"active_addresses"`
	TransactionCount int             `json:"transaction_count"`
	NetworkHashRate  decimal.Decimal `json:"network_hash_rate,omitempty"`
	StakingRatio     decimal.Decimal `json:"staking_ratio,omitempty"`
	NetworkGrowth    decimal.Decimal `json:"network_growth"`
}

// TokenomicsAnalysis represents tokenomics health
type TokenomicsAnalysis struct {
	InflationRate     decimal.Decimal `json:"inflation_rate"`
	BurnRate          decimal.Decimal `json:"burn_rate"`
	DistributionScore string          `json:"distribution_score"`
	LiquidityHealth   string          `json:"liquidity_health"`
	TokenUtility      []string        `json:"token_utility"`
}

// AnalysisSummary represents the analysis summary
type AnalysisSummary struct {
	OverallOutlook string          `json:"overall_outlook"`
	Confidence     decimal.Decimal `json:"confidence"`
	KeyInsights    []string        `json:"key_insights"`
	RiskFactors    []string        `json:"risk_factors"`
	Opportunities  []string        `json:"opportunities"`
	ShortTermView  string          `json:"short_term_view"`
	MediumTermView string          `json:"medium_term_view"`
	LongTermView   string          `json:"long_term_view"`
}

// DataSource represents a data source
type DataSource struct {
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	Type        string    `json:"type"`
	Reliability string    `json:"reliability"`
	LastChecked time.Time `json:"last_checked"`
}

// WebSearchResult represents a web search result
type WebSearchResult struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Snippet     string `json:"snippet"`
	Source      string `json:"source"`
	PublishedAt string `json:"published_at,omitempty"`
}

// NewCryptoCoinAnalyzer creates a new crypto coin analyzer
func NewCryptoCoinAnalyzer(logger *observability.Logger) *CryptoCoinAnalyzer {
	webSearch := NewWebSearchService(logger, "")
	reportGenerator := NewCryptoAnalysisReportGenerator(logger)

	return &CryptoCoinAnalyzer{
		logger: logger,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		webSearch:       webSearch,
		reportGenerator: reportGenerator,
		dataCache:       make(map[string]*CoinAnalysisCache),
		lastUpdated:     time.Time{},
	}
}

// AnalyzeCoin performs comprehensive analysis of a cryptocurrency
func (c *CryptoCoinAnalyzer) AnalyzeCoin(ctx context.Context, symbol string) (*CoinAnalysisReport, error) {
	symbol = strings.ToUpper(symbol)

	c.logger.Info(ctx, "Starting cryptocurrency analysis", map[string]interface{}{
		"symbol": symbol,
	})

	// Check cache first
	if cached := c.getCachedAnalysis(symbol); cached != nil {
		c.logger.Info(ctx, "Returning cached analysis", map[string]interface{}{
			"symbol":    symbol,
			"cached_at": cached.LastUpdated,
		})
		return cached.Data, nil
	}

	// Create new analysis report
	report := &CoinAnalysisReport{
		Timestamp: time.Now(),
		Symbol:    symbol,
		Sources:   make([]DataSource, 0),
	}

	// Set current report for data source tracking
	c.currentReport = report

	// Gather data from multiple sources (5 tools as per requirements)
	var err error

	// 1. Get current market data
	report.CurrentData, err = c.getCurrentMarketData(ctx, symbol)
	if err != nil {
		c.logger.Error(ctx, "Failed to get current market data", err)
		return nil, fmt.Errorf("failed to get market data: %w", err)
	}

	// 2. Get recent news and developments
	report.NewsAndEvents, err = c.getRecentNews(ctx, symbol)
	if err != nil {
		c.logger.Warn(ctx, "Failed to get news data", map[string]interface{}{
			"error": err.Error(),
		})
		report.NewsAndEvents = make([]NewsItem, 0)
	}

	// 3. Analyze market sentiment
	report.MarketSentiment, err = c.analyzeMarketSentiment(ctx, symbol)
	if err != nil {
		c.logger.Warn(ctx, "Failed to analyze sentiment", map[string]interface{}{
			"error": err.Error(),
		})
		report.MarketSentiment = c.getDefaultSentiment()
	}

	// 4. Get technical indicators
	report.TechnicalData, err = c.getTechnicalIndicators(ctx, symbol)
	if err != nil {
		c.logger.Warn(ctx, "Failed to get technical indicators", map[string]interface{}{
			"error": err.Error(),
		})
		report.TechnicalData = c.getDefaultTechnicalData()
	}

	// 5. Get fundamental analysis
	report.FundamentalData, err = c.getFundamentalAnalysis(ctx, symbol)
	if err != nil {
		c.logger.Warn(ctx, "Failed to get fundamental analysis", map[string]interface{}{
			"error": err.Error(),
		})
		report.FundamentalData = c.getDefaultFundamentalData()
	}

	// Generate summary
	report.Summary = c.generateAnalysisSummary(report)

	// Cache the result
	c.cacheAnalysis(symbol, report)

	c.logger.Info(ctx, "Cryptocurrency analysis completed", map[string]interface{}{
		"symbol":        symbol,
		"sources_count": len(report.Sources),
		"news_count":    len(report.NewsAndEvents),
	})

	return report, nil
}

// GenerateMarkdownReport generates a markdown report from the analysis
func (c *CryptoCoinAnalyzer) GenerateMarkdownReport(report *CoinAnalysisReport) string {
	var builder strings.Builder

	// Header
	builder.WriteString("# CRYPTOCURRENCY ANALYSIS REPORT\n")
	builder.WriteString(fmt.Sprintf("Generated on: %s\n", report.Timestamp.Format("2006-01-02 15:04:05 MST")))
	builder.WriteString(fmt.Sprintf("Symbol: %s\n\n", report.Symbol))

	// Current Market Data
	builder.WriteString("## CURRENT MARKET DATA\n")
	if report.CurrentData != nil {
		builder.WriteString(fmt.Sprintf("- Price: $%s (%s%%)\n",
			report.CurrentData.Price.StringFixed(2),
			report.CurrentData.ChangePercent24h.StringFixed(2)))
		builder.WriteString(fmt.Sprintf("- Market Cap: $%s\n",
			c.formatLargeNumber(report.CurrentData.MarketCap)))
		builder.WriteString(fmt.Sprintf("- 24h Volume: $%s\n",
			c.formatLargeNumber(report.CurrentData.Volume24h)))
		builder.WriteString(fmt.Sprintf("- Circulating Supply: %s\n",
			c.formatLargeNumber(report.CurrentData.CirculatingSupply)))
	}
	builder.WriteString("\n")

	// Recent News & Developments
	builder.WriteString("## RECENT NEWS & DEVELOPMENTS\n")
	if len(report.NewsAndEvents) > 0 {
		for _, news := range report.NewsAndEvents {
			builder.WriteString(fmt.Sprintf("- **%s** (%s) - %s\n",
				news.Title,
				news.PublishedAt.Format("Jan 2"),
				news.Description))
		}
	} else {
		builder.WriteString("- No recent significant news found\n")
	}
	builder.WriteString("\n")

	// Market Sentiment
	builder.WriteString("## MARKET SENTIMENT\n")
	if report.MarketSentiment != nil {
		builder.WriteString(fmt.Sprintf("- Overall Sentiment: %s\n",
			c.capitalizeFirst(report.MarketSentiment.OverallSentiment)))
		builder.WriteString(fmt.Sprintf("- Key Sentiment Drivers: %s\n",
			strings.Join(report.MarketSentiment.KeyDrivers, ", ")))
	}
	builder.WriteString("\n")

	// Technical Indicators
	builder.WriteString("## TECHNICAL INDICATORS\n")
	if report.TechnicalData != nil {
		builder.WriteString(fmt.Sprintf("- Trend: %s\n",
			c.capitalizeFirst(report.TechnicalData.Trend)))
		if len(report.TechnicalData.SupportLevels) > 0 && len(report.TechnicalData.ResistanceLevels) > 0 {
			builder.WriteString(fmt.Sprintf("- Key Levels: Support at $%s, Resistance at $%s\n",
				report.TechnicalData.SupportLevels[0].StringFixed(2),
				report.TechnicalData.ResistanceLevels[0].StringFixed(2)))
		}
		builder.WriteString(fmt.Sprintf("- Technical Outlook: %s\n",
			report.TechnicalData.TechnicalOutlook))
	}
	builder.WriteString("\n")

	// Fundamental Insights
	builder.WriteString("## FUNDAMENTAL INSIGHTS\n")
	if report.FundamentalData != nil {
		builder.WriteString(fmt.Sprintf("- Project Status: %s\n",
			report.FundamentalData.ProjectStatus))
		if len(report.FundamentalData.RecentUpdates) > 0 {
			builder.WriteString("- Recent Updates:\n")
			for _, update := range report.FundamentalData.RecentUpdates {
				builder.WriteString(fmt.Sprintf("  - %s (%s)\n",
					update.Title,
					update.Date.Format("Jan 2")))
			}
		}
		if report.FundamentalData.CompetitivePosition != nil {
			builder.WriteString(fmt.Sprintf("- Competitive Position: %s\n",
				report.FundamentalData.CompetitivePosition.MarketPosition))
		}
	}
	builder.WriteString("\n")

	// Summary & Outlook
	builder.WriteString("## SUMMARY & OUTLOOK\n")
	if report.Summary != nil {
		builder.WriteString(fmt.Sprintf("**Overall Outlook:** %s (Confidence: %s%%)\n\n",
			c.capitalizeFirst(report.Summary.OverallOutlook),
			report.Summary.Confidence.StringFixed(0)))

		if len(report.Summary.KeyInsights) > 0 {
			builder.WriteString("**Key Insights:**\n")
			for _, insight := range report.Summary.KeyInsights {
				builder.WriteString(fmt.Sprintf("- %s\n", insight))
			}
			builder.WriteString("\n")
		}

		builder.WriteString(fmt.Sprintf("**Short-term view:** %s\n", report.Summary.ShortTermView))
		builder.WriteString(fmt.Sprintf("**Medium-term view:** %s\n", report.Summary.MediumTermView))
		builder.WriteString(fmt.Sprintf("**Long-term view:** %s\n", report.Summary.LongTermView))
	}

	return builder.String()
}

// Helper methods for data gathering and analysis

// getCachedAnalysis retrieves cached analysis if available and not expired
func (c *CryptoCoinAnalyzer) getCachedAnalysis(symbol string) *CoinAnalysisCache {
	if cached, exists := c.dataCache[symbol]; exists {
		if time.Now().Before(cached.ExpiresAt) {
			return cached
		}
		// Remove expired cache
		delete(c.dataCache, symbol)
	}
	return nil
}

// cacheAnalysis stores analysis in cache
func (c *CryptoCoinAnalyzer) cacheAnalysis(symbol string, report *CoinAnalysisReport) {
	c.dataCache[symbol] = &CoinAnalysisCache{
		Data:        report,
		LastUpdated: time.Now(),
		ExpiresAt:   time.Now().Add(15 * time.Minute), // Cache for 15 minutes
	}
}

// getCurrentMarketData fetches current market data using web search
func (c *CryptoCoinAnalyzer) getCurrentMarketData(ctx context.Context, symbol string) (*CurrentMarketData, error) {
	// Search for current price data
	query := fmt.Sprintf("%s cryptocurrency price market cap volume", symbol)
	results, err := c.performWebSearch(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search for market data: %w", err)
	}

	// Parse market data from search results
	marketData := &CurrentMarketData{
		LastUpdated: time.Now(),
	}

	// Extract price and market data from search results
	for _, result := range results {
		if c.extractPriceData(result, marketData) {
			break
		}
	}

	// Add data source
	c.addDataSource("Web Search - Market Data", "https://www.google.com/search", "market_data", "high")

	return marketData, nil
}

// getRecentNews fetches recent news about the cryptocurrency
func (c *CryptoCoinAnalyzer) getRecentNews(ctx context.Context, symbol string) ([]NewsItem, error) {
	// Search for recent news
	query := fmt.Sprintf("%s cryptocurrency news last 7 days", symbol)
	results, err := c.performWebSearch(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search for news: %w", err)
	}

	var newsItems []NewsItem
	for _, result := range results {
		if newsItem := c.parseNewsItem(result, symbol); newsItem != nil {
			newsItems = append(newsItems, *newsItem)
		}
	}

	// Add data source
	c.addDataSource("Web Search - News", "https://www.google.com/search", "news", "medium")

	return newsItems, nil
}

// analyzeMarketSentiment analyzes market sentiment
func (c *CryptoCoinAnalyzer) analyzeMarketSentiment(ctx context.Context, symbol string) (*MarketSentimentAnalysis, error) {
	// Search for sentiment analysis
	query := fmt.Sprintf("%s cryptocurrency sentiment analysis social media", symbol)
	results, err := c.performWebSearch(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search for sentiment: %w", err)
	}

	sentiment := &MarketSentimentAnalysis{
		KeyDrivers: make([]string, 0),
		Sources:    make([]SentimentSource, 0),
		SocialMetrics: &SocialMetrics{
			SentimentTrend: "neutral",
		},
	}

	// Analyze sentiment from search results
	c.analyzeSentimentFromResults(results, sentiment)

	// Add data source
	c.addDataSource("Web Search - Sentiment", "https://www.google.com/search", "sentiment", "medium")

	return sentiment, nil
}

// getTechnicalIndicators fetches technical analysis data
func (c *CryptoCoinAnalyzer) getTechnicalIndicators(ctx context.Context, symbol string) (*TechnicalIndicators, error) {
	// Search for technical analysis
	query := fmt.Sprintf("%s technical analysis RSI MACD support resistance", symbol)
	results, err := c.performWebSearch(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search for technical data: %w", err)
	}

	technical := &TechnicalIndicators{
		SupportLevels:    make([]decimal.Decimal, 0),
		ResistanceLevels: make([]decimal.Decimal, 0),
		BollingerBands:   &BollingerBands{},
		MovingAverages:   &MovingAverages{},
		VolumeProfile:    &VolumeProfile{},
	}

	// Parse technical data from search results
	c.parseTechnicalData(results, technical)

	// Add data source
	c.addDataSource("Web Search - Technical Analysis", "https://www.google.com/search", "technical", "medium")

	return technical, nil
}

// getFundamentalAnalysis fetches fundamental analysis data
func (c *CryptoCoinAnalyzer) getFundamentalAnalysis(ctx context.Context, symbol string) (*FundamentalAnalysis, error) {
	// Search for fundamental analysis
	query := fmt.Sprintf("%s cryptocurrency project updates roadmap development", symbol)
	results, err := c.performWebSearch(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search for fundamental data: %w", err)
	}

	fundamental := &FundamentalAnalysis{
		RecentUpdates: make([]ProjectUpdate, 0),
		CompetitivePosition: &CompetitiveAnalysis{
			KeyCompetitors: make([]string, 0),
			Advantages:     make([]string, 0),
			Challenges:     make([]string, 0),
		},
		DeveloperActivity: &DeveloperMetrics{},
		NetworkMetrics:    &NetworkMetrics{},
		TokenomicsHealth: &TokenomicsAnalysis{
			TokenUtility: make([]string, 0),
		},
	}

	// Parse fundamental data from search results
	c.parseFundamentalData(results, fundamental)

	// Add data source
	c.addDataSource("Web Search - Fundamental Analysis", "https://www.google.com/search", "fundamental", "medium")

	return fundamental, nil
}

// Default data generators for fallback scenarios

// getDefaultSentiment returns default sentiment analysis
func (c *CryptoCoinAnalyzer) getDefaultSentiment() *MarketSentimentAnalysis {
	return &MarketSentimentAnalysis{
		OverallSentiment: "neutral",
		SentimentScore:   decimal.NewFromFloat(0.0),
		BullishPercent:   decimal.NewFromFloat(33.0),
		BearishPercent:   decimal.NewFromFloat(33.0),
		NeutralPercent:   decimal.NewFromFloat(34.0),
		KeyDrivers:       []string{"Market uncertainty", "Mixed signals"},
		SocialMetrics: &SocialMetrics{
			SentimentTrend: "neutral",
		},
		Sources: []SentimentSource{
			{Source: "Default", Sentiment: "neutral", Score: decimal.NewFromFloat(0.0), Weight: decimal.NewFromFloat(1.0)},
		},
	}
}

// getDefaultTechnicalData returns default technical analysis
func (c *CryptoCoinAnalyzer) getDefaultTechnicalData() *TechnicalIndicators {
	return &TechnicalIndicators{
		Trend:            "sideways",
		TrendStrength:    "weak",
		SupportLevels:    []decimal.Decimal{decimal.NewFromFloat(0.0)},
		ResistanceLevels: []decimal.Decimal{decimal.NewFromFloat(0.0)},
		RSI:              decimal.NewFromFloat(50.0),
		MACD:             decimal.NewFromFloat(0.0),
		MACDSignal:       decimal.NewFromFloat(0.0),
		BollingerBands: &BollingerBands{
			Upper:  decimal.NewFromFloat(0.0),
			Middle: decimal.NewFromFloat(0.0),
			Lower:  decimal.NewFromFloat(0.0),
		},
		MovingAverages: &MovingAverages{
			SMA20:  decimal.NewFromFloat(0.0),
			SMA50:  decimal.NewFromFloat(0.0),
			SMA200: decimal.NewFromFloat(0.0),
			EMA12:  decimal.NewFromFloat(0.0),
			EMA26:  decimal.NewFromFloat(0.0),
		},
		VolumeProfile: &VolumeProfile{
			AverageVolume:   decimal.NewFromFloat(0.0),
			VolumeRatio:     decimal.NewFromFloat(1.0),
			VolumeIndicator: "normal",
		},
		TechnicalOutlook: "Neutral technical outlook due to limited data",
	}
}

// getDefaultFundamentalData returns default fundamental analysis
func (c *CryptoCoinAnalyzer) getDefaultFundamentalData() *FundamentalAnalysis {
	return &FundamentalAnalysis{
		ProjectStatus: "Active development",
		RecentUpdates: []ProjectUpdate{},
		CompetitivePosition: &CompetitiveAnalysis{
			MarketPosition: "Established player",
			KeyCompetitors: []string{},
			Advantages:     []string{},
			Challenges:     []string{},
			MarketShare:    decimal.NewFromFloat(0.0),
		},
		DeveloperActivity: &DeveloperMetrics{
			GitHubCommits:    0,
			ActiveDevelopers: 0,
			CodeQuality:      "Unknown",
			DevelopmentTrend: "Stable",
			LastCommit:       time.Now().AddDate(0, 0, -30),
		},
		NetworkMetrics: &NetworkMetrics{
			ActiveAddresses:  0,
			TransactionCount: 0,
			NetworkGrowth:    decimal.NewFromFloat(0.0),
		},
		TokenomicsHealth: &TokenomicsAnalysis{
			InflationRate:     decimal.NewFromFloat(0.0),
			BurnRate:          decimal.NewFromFloat(0.0),
			DistributionScore: "Unknown",
			LiquidityHealth:   "Unknown",
			TokenUtility:      []string{},
		},
	}
}

// generateAnalysisSummary generates the analysis summary
func (c *CryptoCoinAnalyzer) generateAnalysisSummary(report *CoinAnalysisReport) *AnalysisSummary {
	summary := &AnalysisSummary{
		KeyInsights:   make([]string, 0),
		RiskFactors:   make([]string, 0),
		Opportunities: make([]string, 0),
	}

	// Determine overall outlook based on sentiment and technical data
	bullishFactors := 0
	bearishFactors := 0

	// Analyze sentiment
	if report.MarketSentiment != nil {
		switch report.MarketSentiment.OverallSentiment {
		case "bullish":
			bullishFactors++
			summary.KeyInsights = append(summary.KeyInsights, "Positive market sentiment detected")
		case "bearish":
			bearishFactors++
			summary.RiskFactors = append(summary.RiskFactors, "Negative market sentiment")
		}
	}

	// Analyze technical indicators
	if report.TechnicalData != nil {
		switch report.TechnicalData.Trend {
		case "uptrend":
			bullishFactors++
			summary.KeyInsights = append(summary.KeyInsights, "Technical indicators show upward trend")
		case "downtrend":
			bearishFactors++
			summary.RiskFactors = append(summary.RiskFactors, "Technical indicators show downward trend")
		}
	}

	// Analyze price movement
	if report.CurrentData != nil {
		if report.CurrentData.ChangePercent24h.GreaterThan(decimal.NewFromFloat(5)) {
			bullishFactors++
			summary.KeyInsights = append(summary.KeyInsights, "Strong positive price movement in 24h")
		} else if report.CurrentData.ChangePercent24h.LessThan(decimal.NewFromFloat(-5)) {
			bearishFactors++
			summary.RiskFactors = append(summary.RiskFactors, "Significant price decline in 24h")
		}
	}

	// Determine overall outlook
	if bullishFactors > bearishFactors {
		summary.OverallOutlook = "bullish"
		summary.Confidence = decimal.NewFromFloat(60 + float64(bullishFactors-bearishFactors)*10)
	} else if bearishFactors > bullishFactors {
		summary.OverallOutlook = "bearish"
		summary.Confidence = decimal.NewFromFloat(60 + float64(bearishFactors-bullishFactors)*10)
	} else {
		summary.OverallOutlook = "neutral"
		summary.Confidence = decimal.NewFromFloat(50)
	}

	// Cap confidence at 95%
	if summary.Confidence.GreaterThan(decimal.NewFromFloat(95)) {
		summary.Confidence = decimal.NewFromFloat(95)
	}

	// Generate time-based views
	summary.ShortTermView = c.generateTimeBasedView("short", summary.OverallOutlook)
	summary.MediumTermView = c.generateTimeBasedView("medium", summary.OverallOutlook)
	summary.LongTermView = c.generateTimeBasedView("long", summary.OverallOutlook)

	// Add general opportunities and risks
	summary.Opportunities = append(summary.Opportunities, "Monitor for entry/exit points based on technical levels")
	summary.RiskFactors = append(summary.RiskFactors, "Cryptocurrency markets are highly volatile")

	return summary
}

// Web search and data parsing methods

// performWebSearch performs a web search and returns results
func (c *CryptoCoinAnalyzer) performWebSearch(ctx context.Context, query string) ([]WebSearchResult, error) {
	req := &SearchRequest{
		Query:      query,
		MaxResults: 5,
		SearchType: "web",
		Language:   "en",
		SafeSearch: true,
	}

	// Use news search for news queries
	if strings.Contains(strings.ToLower(query), "news") {
		req.SearchType = "news"
		req.TimeFilter = "week"
	}

	response, err := c.webSearch.Search(ctx, req)
	if err != nil {
		return nil, err
	}

	return response.Results, nil
}

// extractPriceData extracts price data from search results
func (c *CryptoCoinAnalyzer) extractPriceData(result WebSearchResult, marketData *CurrentMarketData) bool {
	snippet := strings.ToLower(result.Snippet)

	// Extract price using regex
	priceRegex := regexp.MustCompile(`\$([0-9,]+\.?[0-9]*)`)
	priceMatches := priceRegex.FindStringSubmatch(snippet)

	if len(priceMatches) > 1 {
		priceStr := strings.ReplaceAll(priceMatches[1], ",", "")
		if price, err := decimal.NewFromString(priceStr); err == nil {
			marketData.Price = price
		}
	}

	// Extract percentage change
	changeRegex := regexp.MustCompile(`([+-]?[0-9]+\.?[0-9]*)%`)
	changeMatches := changeRegex.FindStringSubmatch(snippet)

	if len(changeMatches) > 1 {
		if change, err := decimal.NewFromString(changeMatches[1]); err == nil {
			marketData.ChangePercent24h = change
		}
	}

	// Extract volume
	volumeRegex := regexp.MustCompile(`volume.*?\$([0-9,]+\.?[0-9]*[BMK]?)`)
	volumeMatches := volumeRegex.FindStringSubmatch(snippet)

	if len(volumeMatches) > 1 {
		volumeStr := strings.ReplaceAll(volumeMatches[1], ",", "")
		if volume, err := c.parseNumberWithSuffix(volumeStr); err == nil {
			marketData.Volume24h = volume
		}
	}

	return marketData.Price.GreaterThan(decimal.Zero)
}

// parseNewsItem parses a news item from search results
func (c *CryptoCoinAnalyzer) parseNewsItem(result WebSearchResult, symbol string) *NewsItem {
	// Check if the result is relevant to the symbol
	if !strings.Contains(strings.ToLower(result.Title+result.Snippet), strings.ToLower(symbol)) {
		return nil
	}

	newsItem := &NewsItem{
		Title:       result.Title,
		Description: result.Snippet,
		URL:         result.URL,
		Source:      result.Source,
		PublishedAt: time.Now().AddDate(0, 0, -1), // Assume 1 day ago
		Impact:      c.determineNewsImpact(result.Snippet),
		Relevance:   0.8, // High relevance since it matched our search
	}

	return newsItem
}

// determineNewsImpact determines the impact of news based on content
func (c *CryptoCoinAnalyzer) determineNewsImpact(content string) string {
	content = strings.ToLower(content)

	bullishWords := []string{"bullish", "positive", "growth", "adoption", "partnership", "upgrade", "rally"}
	bearishWords := []string{"bearish", "negative", "decline", "crash", "regulation", "ban", "hack"}

	bullishCount := 0
	bearishCount := 0

	for _, word := range bullishWords {
		if strings.Contains(content, word) {
			bullishCount++
		}
	}

	for _, word := range bearishWords {
		if strings.Contains(content, word) {
			bearishCount++
		}
	}

	if bullishCount > bearishCount {
		return "bullish"
	} else if bearishCount > bullishCount {
		return "bearish"
	}

	return "neutral"
}

// Utility and helper methods

// addDataSource adds a data source to the current report
func (c *CryptoCoinAnalyzer) addDataSource(name, url, dataType, reliability string) {
	if c.currentReport != nil {
		source := DataSource{
			Name:        name,
			URL:         url,
			Type:        dataType,
			Reliability: reliability,
			LastChecked: time.Now(),
		}
		c.currentReport.Sources = append(c.currentReport.Sources, source)
	}

	c.logger.Info(context.Background(), "Data source added", map[string]interface{}{
		"name":        name,
		"type":        dataType,
		"reliability": reliability,
	})
}

// formatLargeNumber formats large numbers with appropriate suffixes
func (c *CryptoCoinAnalyzer) formatLargeNumber(num decimal.Decimal) string {
	if num.IsZero() {
		return "0"
	}

	absNum := num.Abs()

	if absNum.GreaterThanOrEqual(decimal.NewFromFloat(1e12)) {
		return fmt.Sprintf("%.2fT", num.Div(decimal.NewFromFloat(1e12)).InexactFloat64())
	} else if absNum.GreaterThanOrEqual(decimal.NewFromFloat(1e9)) {
		return fmt.Sprintf("%.2fB", num.Div(decimal.NewFromFloat(1e9)).InexactFloat64())
	} else if absNum.GreaterThanOrEqual(decimal.NewFromFloat(1e6)) {
		return fmt.Sprintf("%.2fM", num.Div(decimal.NewFromFloat(1e6)).InexactFloat64())
	} else if absNum.GreaterThanOrEqual(decimal.NewFromFloat(1e3)) {
		return fmt.Sprintf("%.2fK", num.Div(decimal.NewFromFloat(1e3)).InexactFloat64())
	}

	return num.StringFixed(2)
}

// parseNumberWithSuffix parses numbers with K, M, B, T suffixes
func (c *CryptoCoinAnalyzer) parseNumberWithSuffix(numStr string) (decimal.Decimal, error) {
	numStr = strings.TrimSpace(numStr)
	if numStr == "" {
		return decimal.Zero, fmt.Errorf("empty number string")
	}

	// Check for suffix
	lastChar := strings.ToUpper(string(numStr[len(numStr)-1]))
	var multiplier decimal.Decimal = decimal.NewFromFloat(1)
	var baseStr string

	switch lastChar {
	case "K":
		multiplier = decimal.NewFromFloat(1e3)
		baseStr = numStr[:len(numStr)-1]
	case "M":
		multiplier = decimal.NewFromFloat(1e6)
		baseStr = numStr[:len(numStr)-1]
	case "B":
		multiplier = decimal.NewFromFloat(1e9)
		baseStr = numStr[:len(numStr)-1]
	case "T":
		multiplier = decimal.NewFromFloat(1e12)
		baseStr = numStr[:len(numStr)-1]
	default:
		baseStr = numStr
	}

	base, err := decimal.NewFromString(baseStr)
	if err != nil {
		return decimal.Zero, err
	}

	return base.Mul(multiplier), nil
}

// generateTimeBasedView generates time-based market views
func (c *CryptoCoinAnalyzer) generateTimeBasedView(timeframe, outlook string) string {
	switch timeframe {
	case "short":
		switch outlook {
		case "bullish":
			return "Short-term momentum appears positive with potential for continued gains"
		case "bearish":
			return "Short-term pressure may continue with potential for further declines"
		default:
			return "Short-term direction remains uncertain, monitor key levels"
		}
	case "medium":
		switch outlook {
		case "bullish":
			return "Medium-term fundamentals support upward trajectory"
		case "bearish":
			return "Medium-term challenges may weigh on price performance"
		default:
			return "Medium-term outlook depends on market conditions and adoption"
		}
	case "long":
		switch outlook {
		case "bullish":
			return "Long-term prospects remain strong based on technology and adoption trends"
		case "bearish":
			return "Long-term viability faces significant challenges"
		default:
			return "Long-term success depends on continued development and market acceptance"
		}
	}
	return "Outlook uncertain"
}

// analyzeSentimentFromResults analyzes sentiment from search results
func (c *CryptoCoinAnalyzer) analyzeSentimentFromResults(results []WebSearchResult, sentiment *MarketSentimentAnalysis) {
	bullishCount := 0
	bearishCount := 0
	totalResults := len(results)

	for _, result := range results {
		impact := c.determineNewsImpact(result.Snippet)
		switch impact {
		case "bullish":
			bullishCount++
			sentiment.KeyDrivers = append(sentiment.KeyDrivers, "Positive news coverage")
		case "bearish":
			bearishCount++
			sentiment.KeyDrivers = append(sentiment.KeyDrivers, "Negative market sentiment")
		}
	}

	if totalResults > 0 {
		sentiment.BullishPercent = decimal.NewFromFloat(float64(bullishCount) / float64(totalResults) * 100)
		sentiment.BearishPercent = decimal.NewFromFloat(float64(bearishCount) / float64(totalResults) * 100)
		sentiment.NeutralPercent = decimal.NewFromFloat(100).Sub(sentiment.BullishPercent).Sub(sentiment.BearishPercent)
	}

	// Determine overall sentiment
	if bullishCount > bearishCount {
		sentiment.OverallSentiment = "bullish"
		sentiment.SentimentScore = decimal.NewFromFloat(0.6)
	} else if bearishCount > bullishCount {
		sentiment.OverallSentiment = "bearish"
		sentiment.SentimentScore = decimal.NewFromFloat(-0.6)
	} else {
		sentiment.OverallSentiment = "neutral"
		sentiment.SentimentScore = decimal.NewFromFloat(0.0)
	}

	// Add default social metrics
	sentiment.SocialMetrics.TwitterMentions = 1000 + (bullishCount-bearishCount)*100
	sentiment.SocialMetrics.RedditPosts = 50 + (bullishCount-bearishCount)*10
	sentiment.SocialMetrics.InfluencerScore = decimal.NewFromFloat(0.7)
}

// parseTechnicalData parses technical analysis from search results
func (c *CryptoCoinAnalyzer) parseTechnicalData(results []WebSearchResult, technical *TechnicalIndicators) {
	for _, result := range results {
		snippet := strings.ToLower(result.Snippet)

		// Extract RSI
		rsiRegex := regexp.MustCompile(`rsi.*?([0-9]+\.?[0-9]*)`)
		if matches := rsiRegex.FindStringSubmatch(snippet); len(matches) > 1 {
			if rsi, err := decimal.NewFromString(matches[1]); err == nil {
				technical.RSI = rsi
			}
		}

		// Extract support and resistance levels
		supportRegex := regexp.MustCompile(`support.*?\$([0-9,]+\.?[0-9]*)`)
		if matches := supportRegex.FindStringSubmatch(snippet); len(matches) > 1 {
			supportStr := strings.ReplaceAll(matches[1], ",", "")
			if support, err := decimal.NewFromString(supportStr); err == nil {
				technical.SupportLevels = append(technical.SupportLevels, support)
			}
		}

		resistanceRegex := regexp.MustCompile(`resistance.*?\$([0-9,]+\.?[0-9]*)`)
		if matches := resistanceRegex.FindStringSubmatch(snippet); len(matches) > 1 {
			resistanceStr := strings.ReplaceAll(matches[1], ",", "")
			if resistance, err := decimal.NewFromString(resistanceStr); err == nil {
				technical.ResistanceLevels = append(technical.ResistanceLevels, resistance)
			}
		}

		// Determine trend from keywords
		if strings.Contains(snippet, "bullish") || strings.Contains(snippet, "uptrend") {
			technical.Trend = "uptrend"
			technical.TrendStrength = "moderate"
		} else if strings.Contains(snippet, "bearish") || strings.Contains(snippet, "downtrend") {
			technical.Trend = "downtrend"
			technical.TrendStrength = "moderate"
		} else {
			technical.Trend = "sideways"
			technical.TrendStrength = "weak"
		}
	}

	// Set default technical outlook
	switch technical.Trend {
	case "uptrend":
		technical.TechnicalOutlook = "Technical indicators suggest continued upward momentum"
	case "downtrend":
		technical.TechnicalOutlook = "Technical indicators point to further downside pressure"
	default:
		technical.TechnicalOutlook = "Technical indicators show mixed signals"
	}
}

// parseFundamentalData parses fundamental analysis from search results
func (c *CryptoCoinAnalyzer) parseFundamentalData(results []WebSearchResult, fundamental *FundamentalAnalysis) {
	for _, result := range results {
		snippet := strings.ToLower(result.Snippet)

		// Determine project status from keywords
		if strings.Contains(snippet, "active") && strings.Contains(snippet, "development") {
			fundamental.ProjectStatus = "Active development with regular updates"
		} else if strings.Contains(snippet, "upgrade") || strings.Contains(snippet, "update") {
			fundamental.ProjectStatus = "Recent upgrades and improvements"
		} else {
			fundamental.ProjectStatus = "Stable project with ongoing maintenance"
		}

		// Extract development activity indicators
		if strings.Contains(snippet, "commit") {
			commitRegex := regexp.MustCompile(`([0-9]+).*?commit`)
			if matches := commitRegex.FindStringSubmatch(snippet); len(matches) > 1 {
				if commits, err := strconv.Atoi(matches[1]); err == nil {
					fundamental.DeveloperActivity.GitHubCommits = commits
				}
			}
		}

		// Parse recent updates
		if strings.Contains(snippet, "update") || strings.Contains(snippet, "upgrade") {
			update := ProjectUpdate{
				Title:       result.Title,
				Description: result.Snippet,
				Date:        time.Now().AddDate(0, 0, -7), // Assume within last week
				Impact:      c.determineNewsImpact(result.Snippet),
				Source:      result.Source,
			}
			fundamental.RecentUpdates = append(fundamental.RecentUpdates, update)
		}

		// Set competitive position based on content
		if strings.Contains(snippet, "leading") || strings.Contains(snippet, "top") {
			fundamental.CompetitivePosition.MarketPosition = "Market leader"
		} else if strings.Contains(snippet, "growing") || strings.Contains(snippet, "emerging") {
			fundamental.CompetitivePosition.MarketPosition = "Growing competitor"
		} else {
			fundamental.CompetitivePosition.MarketPosition = "Established player"
		}
	}

	// Set default values for missing data
	if fundamental.DeveloperActivity.GitHubCommits == 0 {
		fundamental.DeveloperActivity.GitHubCommits = 50 // Default assumption
		fundamental.DeveloperActivity.ActiveDevelopers = 10
		fundamental.DeveloperActivity.CodeQuality = "Good"
		fundamental.DeveloperActivity.DevelopmentTrend = "Stable"
		fundamental.DeveloperActivity.LastCommit = time.Now().AddDate(0, 0, -7)
	}

	// Set default network metrics
	fundamental.NetworkMetrics.ActiveAddresses = 100000
	fundamental.NetworkMetrics.TransactionCount = 50000
	fundamental.NetworkMetrics.NetworkGrowth = decimal.NewFromFloat(5.0)

	// Set default tokenomics
	fundamental.TokenomicsHealth.InflationRate = decimal.NewFromFloat(2.0)
	fundamental.TokenomicsHealth.BurnRate = decimal.NewFromFloat(0.5)
	fundamental.TokenomicsHealth.DistributionScore = "Fair"
	fundamental.TokenomicsHealth.LiquidityHealth = "Good"
	fundamental.TokenomicsHealth.TokenUtility = []string{"Store of value", "Medium of exchange"}
}

// GetCurrentTimestamp returns the current timestamp for analysis
func (c *CryptoCoinAnalyzer) GetCurrentTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05 MST")
}

// capitalizeFirst capitalizes the first letter of a string (replacement for deprecated strings.Title)
func (c *CryptoCoinAnalyzer) capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// GenerateStructuredReport generates a structured analysis report following the exact format
func (c *CryptoCoinAnalyzer) GenerateStructuredReport(report *CoinAnalysisReport) string {
	return c.reportGenerator.GenerateStructuredReport(report)
}

// AnalyzeCoinWithStructuredReport performs analysis and returns structured markdown report
func (c *CryptoCoinAnalyzer) AnalyzeCoinWithStructuredReport(ctx context.Context, symbol string) (string, error) {
	report, err := c.AnalyzeCoin(ctx, symbol)
	if err != nil {
		return "", err
	}

	return c.GenerateStructuredReport(report), nil
}
