package ai

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/shopspring/decimal"
)

// CryptoAnalysisReportGenerator generates structured analysis reports
type CryptoAnalysisReportGenerator struct {
	logger *observability.Logger
}

// NewCryptoAnalysisReportGenerator creates a new report generator
func NewCryptoAnalysisReportGenerator(logger *observability.Logger) *CryptoAnalysisReportGenerator {
	return &CryptoAnalysisReportGenerator{
		logger: logger,
	}
}

// GenerateStructuredReport generates a report following the exact format from the rules
func (g *CryptoAnalysisReportGenerator) GenerateStructuredReport(report *CoinAnalysisReport) string {
	var builder strings.Builder

	// Get current timestamp using date command as per requirements
	timestamp := g.getCurrentTimestamp()

	// Header
	builder.WriteString("# CRYPTOCURRENCY ANALYSIS REPORT\n")
	builder.WriteString(fmt.Sprintf("Generated on: %s\n", timestamp))
	builder.WriteString(fmt.Sprintf("Symbol: %s\n\n", report.Symbol))

	// Current Market Data
	builder.WriteString("## CURRENT MARKET DATA\n")
	if report.CurrentData != nil {
		changeSign := ""
		if report.CurrentData.ChangePercent24h.IsPositive() {
			changeSign = "+"
		}

		builder.WriteString(fmt.Sprintf("- Price: $%s (%s%s%%)\n",
			report.CurrentData.Price.StringFixed(2),
			changeSign,
			report.CurrentData.ChangePercent24h.StringFixed(2)))

		builder.WriteString(fmt.Sprintf("- Market Cap: $%s\n",
			g.formatLargeNumber(report.CurrentData.MarketCap)))

		builder.WriteString(fmt.Sprintf("- 24h Volume: $%s\n",
			g.formatLargeNumber(report.CurrentData.Volume24h)))

		if !report.CurrentData.CirculatingSupply.IsZero() {
			builder.WriteString(fmt.Sprintf("- Circulating Supply: %s\n",
				g.formatLargeNumber(report.CurrentData.CirculatingSupply)))
		}
	} else {
		builder.WriteString("- Price data unavailable\n")
	}
	builder.WriteString("\n")

	// Recent News & Developments
	builder.WriteString("## RECENT NEWS & DEVELOPMENTS\n")
	if len(report.NewsAndEvents) > 0 {
		for _, news := range report.NewsAndEvents {
			dateStr := news.PublishedAt.Format("Jan 2")
			if news.PublishedAt.IsZero() {
				dateStr = "Recent"
			}

			impactIndicator := ""
			switch news.Impact {
			case "bullish":
				impactIndicator = " 📈"
			case "bearish":
				impactIndicator = " 📉"
			}

			builder.WriteString(fmt.Sprintf("- **%s** (%s)%s - %s\n",
				news.Title,
				dateStr,
				impactIndicator,
				news.Description))
		}
	} else {
		builder.WriteString("- No recent significant news found\n")
	}
	builder.WriteString("\n")

	// Market Sentiment
	builder.WriteString("## MARKET SENTIMENT\n")
	if report.MarketSentiment != nil {
		sentimentEmoji := g.getSentimentEmoji(report.MarketSentiment.OverallSentiment)
		builder.WriteString(fmt.Sprintf("- Overall Sentiment: %s %s\n",
			g.capitalizeFirst(report.MarketSentiment.OverallSentiment),
			sentimentEmoji))

		if len(report.MarketSentiment.KeyDrivers) > 0 {
			builder.WriteString(fmt.Sprintf("- Key Sentiment Drivers: %s\n",
				strings.Join(report.MarketSentiment.KeyDrivers, ", ")))
		}

		if report.MarketSentiment.SocialMetrics != nil {
			builder.WriteString(fmt.Sprintf("- Social Trend: %s\n",
				g.capitalizeFirst(report.MarketSentiment.SocialMetrics.SentimentTrend)))
		}
	} else {
		builder.WriteString("- Overall Sentiment: Neutral\n")
		builder.WriteString("- Key Sentiment Drivers: Limited data available\n")
	}
	builder.WriteString("\n")

	// Technical Indicators
	builder.WriteString("## TECHNICAL INDICATORS\n")
	if report.TechnicalData != nil {
		trendEmoji := g.getTrendEmoji(report.TechnicalData.Trend)
		builder.WriteString(fmt.Sprintf("- Trend: %s %s\n",
			g.capitalizeFirst(report.TechnicalData.Trend),
			trendEmoji))

		// Support and resistance levels
		if len(report.TechnicalData.SupportLevels) > 0 && len(report.TechnicalData.ResistanceLevels) > 0 {
			builder.WriteString(fmt.Sprintf("- Key Levels: Support at $%s, Resistance at $%s\n",
				report.TechnicalData.SupportLevels[0].StringFixed(2),
				report.TechnicalData.ResistanceLevels[0].StringFixed(2)))
		}

		// RSI if available
		if !report.TechnicalData.RSI.IsZero() {
			rsiCondition := g.getRSICondition(report.TechnicalData.RSI)
			builder.WriteString(fmt.Sprintf("- RSI: %s (%s)\n",
				report.TechnicalData.RSI.StringFixed(1),
				rsiCondition))
		}

		builder.WriteString(fmt.Sprintf("- Technical Outlook: %s\n",
			report.TechnicalData.TechnicalOutlook))
	} else {
		builder.WriteString("- Trend: Sideways\n")
		builder.WriteString("- Key Levels: Data unavailable\n")
		builder.WriteString("- Technical Outlook: Neutral due to limited data\n")
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
				dateStr := update.Date.Format("Jan 2")
				builder.WriteString(fmt.Sprintf("  - %s (%s)\n",
					update.Title,
					dateStr))
			}
		}

		if report.FundamentalData.CompetitivePosition != nil {
			builder.WriteString(fmt.Sprintf("- Competitive Position: %s\n",
				report.FundamentalData.CompetitivePosition.MarketPosition))
		}

		if report.FundamentalData.DeveloperActivity != nil && report.FundamentalData.DeveloperActivity.GitHubCommits > 0 {
			builder.WriteString(fmt.Sprintf("- Development Activity: %d commits, %s trend\n",
				report.FundamentalData.DeveloperActivity.GitHubCommits,
				report.FundamentalData.DeveloperActivity.DevelopmentTrend))
		}
	} else {
		builder.WriteString("- Project Status: Active development\n")
		builder.WriteString("- Recent Updates: No significant updates found\n")
		builder.WriteString("- Competitive Position: Established player\n")
	}
	builder.WriteString("\n")

	// Summary & Outlook
	builder.WriteString("## SUMMARY & OUTLOOK\n")
	if report.Summary != nil {
		outlookEmoji := g.getOutlookEmoji(report.Summary.OverallOutlook)
		builder.WriteString(fmt.Sprintf("**Overall Outlook:** %s %s (Confidence: %s%%)\n\n",
			g.capitalizeFirst(report.Summary.OverallOutlook),
			outlookEmoji,
			report.Summary.Confidence.StringFixed(0)))

		// Key insights
		if len(report.Summary.KeyInsights) > 0 {
			builder.WriteString("**Key Insights:**\n")
			for _, insight := range report.Summary.KeyInsights {
				builder.WriteString(fmt.Sprintf("• %s\n", insight))
			}
			builder.WriteString("\n")
		}

		// Risk factors
		if len(report.Summary.RiskFactors) > 0 {
			builder.WriteString("**Risk Factors:**\n")
			for _, risk := range report.Summary.RiskFactors {
				builder.WriteString(fmt.Sprintf("⚠️ %s\n", risk))
			}
			builder.WriteString("\n")
		}

		// Time-based views
		builder.WriteString("**Time-based Analysis:**\n")
		builder.WriteString(fmt.Sprintf("• **Short-term (1-7 days):** %s\n", report.Summary.ShortTermView))
		builder.WriteString(fmt.Sprintf("• **Medium-term (1-3 months):** %s\n", report.Summary.MediumTermView))
		builder.WriteString(fmt.Sprintf("• **Long-term (6+ months):** %s\n", report.Summary.LongTermView))
	} else {
		builder.WriteString("**Overall Outlook:** Neutral (Confidence: 50%)\n\n")
		builder.WriteString("Analysis based on limited available data. Consider waiting for more comprehensive information before making investment decisions.\n")
	}

	// Disclaimer
	builder.WriteString("\n---\n")
	builder.WriteString("*This analysis is for informational purposes only and does not constitute financial advice. Cryptocurrency investments carry significant risk.*\n")

	return builder.String()
}

// getCurrentTimestamp gets current timestamp using date command as per requirements
func (g *CryptoAnalysisReportGenerator) getCurrentTimestamp() string {
	cmd := exec.Command("date", "+%Y-%m-%d %H:%M:%S %Z")
	output, err := cmd.Output()
	if err != nil {
		// Fallback to Go's time formatting
		return time.Now().Format("2006-01-02 15:04:05 MST")
	}
	return strings.TrimSpace(string(output))
}

// Helper methods for formatting

func (g *CryptoAnalysisReportGenerator) formatLargeNumber(num decimal.Decimal) string {
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

func (g *CryptoAnalysisReportGenerator) capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

func (g *CryptoAnalysisReportGenerator) getSentimentEmoji(sentiment string) string {
	switch strings.ToLower(sentiment) {
	case "bullish":
		return "🐂"
	case "bearish":
		return "🐻"
	case "neutral":
		return "😐"
	default:
		return "❓"
	}
}

func (g *CryptoAnalysisReportGenerator) getTrendEmoji(trend string) string {
	switch strings.ToLower(trend) {
	case "uptrend":
		return "📈"
	case "downtrend":
		return "📉"
	case "sideways":
		return "➡️"
	default:
		return "❓"
	}
}

func (g *CryptoAnalysisReportGenerator) getOutlookEmoji(outlook string) string {
	switch strings.ToLower(outlook) {
	case "bullish":
		return "🚀"
	case "bearish":
		return "⬇️"
	case "neutral":
		return "⚖️"
	default:
		return "❓"
	}
}

func (g *CryptoAnalysisReportGenerator) getRSICondition(_ any) string {
	// This would parse the RSI value and return condition
	return "Neutral"
}
