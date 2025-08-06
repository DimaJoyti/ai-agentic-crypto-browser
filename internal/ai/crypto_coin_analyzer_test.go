package ai

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/ai-agentic-browser/internal/config"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/shopspring/decimal"
)

func TestNewCryptoCoinAnalyzer(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{
		ServiceName: "test",
		LogLevel:    "info",
		LogFormat:   "text",
	})

	analyzer := NewCryptoCoinAnalyzer(logger)

	if analyzer == nil {
		t.Fatal("Expected analyzer to be created, got nil")
	}

	if analyzer.logger == nil {
		t.Error("Expected logger to be set")
	}

	if analyzer.webSearch == nil {
		t.Error("Expected web search service to be initialized")
	}

	if analyzer.reportGenerator == nil {
		t.Error("Expected report generator to be initialized")
	}

	if analyzer.dataCache == nil {
		t.Error("Expected data cache to be initialized")
	}
}

func TestAnalyzeCoin(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{
		ServiceName: "test",
		LogLevel:    "info",
		LogFormat:   "text",
	})

	analyzer := NewCryptoCoinAnalyzer(logger)
	ctx := context.Background()

	tests := []struct {
		name    string
		symbol  string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Valid Bitcoin symbol",
			symbol:  "BTC",
			wantErr: false,
		},
		{
			name:    "Valid Ethereum symbol",
			symbol:  "ETH",
			wantErr: false,
		},
		{
			name:    "Lowercase symbol",
			symbol:  "btc",
			wantErr: false,
		},
		{
			name:    "Valid alternative coin",
			symbol:  "ADA",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report, err := analyzer.AnalyzeCoin(ctx, tt.symbol)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected error message to contain '%s', got '%s'", tt.errMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if report == nil {
				t.Error("Expected report to be generated, got nil")
				return
			}

			// Validate report structure
			if report.Symbol != strings.ToUpper(tt.symbol) {
				t.Errorf("Expected symbol to be %s, got %s", strings.ToUpper(tt.symbol), report.Symbol)
			}

			if report.Timestamp.IsZero() {
				t.Error("Expected timestamp to be set")
			}

			if report.CurrentData == nil {
				t.Error("Expected current data to be set")
			}

			if report.Summary == nil {
				t.Error("Expected summary to be set")
			}
		})
	}
}

func TestGenerateMarkdownReport(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{
		ServiceName: "test",
		LogLevel:    "info",
		LogFormat:   "text",
	})

	analyzer := NewCryptoCoinAnalyzer(logger)

	// Create a sample report
	report := &CoinAnalysisReport{
		Timestamp: time.Now(),
		Symbol:    "BTC",
		CurrentData: &CurrentMarketData{
			Price:             decimal.NewFromFloat(45000.00),
			ChangePercent24h:  decimal.NewFromFloat(2.5),
			MarketCap:         decimal.NewFromFloat(850000000000),
			Volume24h:         decimal.NewFromFloat(15000000000),
			CirculatingSupply: decimal.NewFromFloat(19500000),
			LastUpdated:       time.Now(),
		},
		NewsAndEvents: []NewsItem{
			{
				Title:       "Bitcoin Reaches New High",
				Description: "Bitcoin price surged to new monthly high",
				Source:      "CoinDesk",
				PublishedAt: time.Now().AddDate(0, 0, -1),
				Impact:      "bullish",
			},
		},
		MarketSentiment: &MarketSentimentAnalysis{
			OverallSentiment: "bullish",
			KeyDrivers:       []string{"Institutional adoption", "Positive news"},
		},
		TechnicalData: &TechnicalIndicators{
			Trend:            "uptrend",
			TechnicalOutlook: "Bullish momentum continues",
			SupportLevels:    []decimal.Decimal{decimal.NewFromFloat(43000)},
			ResistanceLevels: []decimal.Decimal{decimal.NewFromFloat(47000)},
		},
		Summary: &AnalysisSummary{
			OverallOutlook: "bullish",
			Confidence:     decimal.NewFromFloat(75),
			KeyInsights:    []string{"Strong institutional interest", "Technical indicators positive"},
			ShortTermView:  "Positive momentum expected to continue",
			MediumTermView: "Fundamentals support higher prices",
			LongTermView:   "Long-term outlook remains strong",
		},
		Sources: []DataSource{
			{
				Name: "Test Source",
				Type: "test",
			},
		},
	}

	markdown := analyzer.GenerateStructuredReport(report)

	// Validate markdown content
	if !strings.Contains(markdown, "# CRYPTOCURRENCY ANALYSIS REPORT") {
		t.Error("Expected markdown to contain main header")
	}

	if !strings.Contains(markdown, "Symbol: BTC") {
		t.Error("Expected markdown to contain symbol")
	}

	if !strings.Contains(markdown, "## CURRENT MARKET DATA") {
		t.Error("Expected markdown to contain market data section")
	}

	if !strings.Contains(markdown, "## RECENT NEWS & DEVELOPMENTS") {
		t.Error("Expected markdown to contain news section")
	}

	if !strings.Contains(markdown, "## MARKET SENTIMENT") {
		t.Error("Expected markdown to contain sentiment section")
	}

	if !strings.Contains(markdown, "## TECHNICAL INDICATORS") {
		t.Error("Expected markdown to contain technical section")
	}

	if !strings.Contains(markdown, "## SUMMARY & OUTLOOK") {
		t.Error("Expected markdown to contain summary section")
	}

	if !strings.Contains(markdown, "$45000.00") {
		t.Error("Expected markdown to contain price")
	}

	if !strings.Contains(markdown, "Bullish") {
		t.Error("Expected markdown to contain sentiment")
	}
}

func TestAnalyzeCoinWithStructuredReport(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{
		ServiceName: "test",
		LogLevel:    "info",
		LogFormat:   "text",
	})

	analyzer := NewCryptoCoinAnalyzer(logger)
	ctx := context.Background()

	report, err := analyzer.AnalyzeCoinWithStructuredReport(ctx, "BTC")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if report == "" {
		t.Error("Expected structured report to be generated")
		return
	}

	// Validate structured report format
	if !strings.Contains(report, "# CRYPTOCURRENCY ANALYSIS REPORT") {
		t.Error("Expected structured report to contain main header")
	}

	if !strings.Contains(report, "Generated on:") {
		t.Error("Expected structured report to contain timestamp")
	}

	if !strings.Contains(report, "Symbol: BTC") {
		t.Error("Expected structured report to contain symbol")
	}
}

func TestCacheAnalysis(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{
		ServiceName: "test",
		LogLevel:    "info",
		LogFormat:   "text",
	})

	analyzer := NewCryptoCoinAnalyzer(logger)

	// Create a sample report
	report := &CoinAnalysisReport{
		Timestamp: time.Now(),
		Symbol:    "BTC",
		CurrentData: &CurrentMarketData{
			Price: decimal.NewFromFloat(45000),
		},
	}

	// Cache the report
	analyzer.cacheAnalysis("BTC", report)

	// Retrieve from cache
	cached := analyzer.getCachedAnalysis("BTC")
	if cached == nil {
		t.Error("Expected cached analysis to be retrieved")
		return
	}

	if cached.Data.Symbol != "BTC" {
		t.Errorf("Expected cached symbol to be BTC, got %s", cached.Data.Symbol)
	}

	// Test cache expiration
	analyzer.dataCache["BTC"].ExpiresAt = time.Now().Add(-1 * time.Hour)
	expired := analyzer.getCachedAnalysis("BTC")
	if expired != nil {
		t.Error("Expected expired cache to return nil")
	}
}

func TestFormatLargeNumber(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{
		ServiceName: "test",
		LogLevel:    "info",
		LogFormat:   "text",
	})

	analyzer := NewCryptoCoinAnalyzer(logger)

	tests := []struct {
		name     string
		input    decimal.Decimal
		expected string
	}{
		{
			name:     "Trillion",
			input:    decimal.NewFromFloat(1200000000000),
			expected: "1.20T",
		},
		{
			name:     "Billion",
			input:    decimal.NewFromFloat(850000000000),
			expected: "850.00B",
		},
		{
			name:     "Million",
			input:    decimal.NewFromFloat(15000000),
			expected: "15.00M",
		},
		{
			name:     "Thousand",
			input:    decimal.NewFromFloat(5000),
			expected: "5.00K",
		},
		{
			name:     "Small number",
			input:    decimal.NewFromFloat(123.45),
			expected: "123.45",
		},
		{
			name:     "Zero",
			input:    decimal.Zero,
			expected: "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.formatLargeNumber(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetCurrentTimestamp(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{
		ServiceName: "test",
		LogLevel:    "info",
		LogFormat:   "text",
	})

	analyzer := NewCryptoCoinAnalyzer(logger)
	timestamp := analyzer.GetCurrentTimestamp()

	if timestamp == "" {
		t.Error("Expected timestamp to be generated")
	}

	// Basic format validation
	if !strings.Contains(timestamp, "2025") {
		t.Error("Expected timestamp to contain current year")
	}
}
