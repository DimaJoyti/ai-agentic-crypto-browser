# Crypto Coin Analyzer Agent

The Crypto Coin Analyzer Agent is a comprehensive AI-powered cryptocurrency analysis tool that provides real-time market insights for individual cryptocurrencies. It follows the exact specifications from the agent rules and provides structured analysis reports.

## Features

- **Real-time Market Data**: Current price, market cap, volume, and supply information
- **News Analysis**: Recent news and developments with impact assessment
- **Sentiment Analysis**: Market sentiment from social media and news sources
- **Technical Analysis**: RSI, MACD, support/resistance levels, and trend analysis
- **Fundamental Analysis**: Project status, development activity, and competitive position
- **Structured Reports**: Markdown-formatted reports following exact specifications
- **Multiple Interfaces**: REST API endpoints and CLI tool
- **Caching**: 15-minute cache for improved performance
- **Data Source Tracking**: Tracks all data sources used in analysis

## Architecture

### Core Components

1. **CryptoCoinAnalyzer**: Main service for cryptocurrency analysis
2. **WebSearchService**: Handles web search functionality for data gathering
3. **CryptoAnalysisReportGenerator**: Generates structured markdown reports
4. **CLI Tool**: Command-line interface for direct analysis

### Data Sources

The analyzer uses 5 different data sources as per requirements:

1. **Market Data**: Price, volume, market cap from search results
2. **News Data**: Recent news and developments
3. **Sentiment Data**: Social media and community sentiment
4. **Technical Data**: Technical indicators and chart analysis
5. **Fundamental Data**: Project updates, development activity, tokenomics

## Usage

### REST API

#### Analyze Cryptocurrency (JSON Response)

```bash
# GET or POST request
curl -X GET "http://localhost:8082/ai/crypto/analyze/BTC" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json"
```

#### Generate Structured Report (Markdown)

```bash
# GET or POST request
curl -X GET "http://localhost:8082/ai/crypto/report/BTC" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Get JSON Response with Markdown Content

```bash
curl -X GET "http://localhost:8082/ai/crypto/report/BTC" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Accept: application/json"
```

### CLI Tool

#### Build the CLI Tool

```bash
make build-crypto-analyzer
```

#### Basic Usage

```bash
# Analyze Bitcoin with markdown output
./bin/crypto-analyzer -symbol BTC

# Analyze Ethereum with JSON output
./bin/crypto-analyzer -symbol ETH -format json

# Verbose analysis with custom timeout
./bin/crypto-analyzer -symbol BTC -verbose -timeout 120s
```

#### CLI Options

- `-symbol`: Cryptocurrency symbol to analyze (required)
- `-format`: Output format (markdown, json) - default: markdown
- `-verbose`: Enable verbose logging
- `-timeout`: Analysis timeout - default: 60s
- `-help`: Show help message
- `-version`: Show version information

### Programmatic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/ai-agentic-browser/internal/ai"
    "github.com/ai-agentic-browser/pkg/observability"
)

func main() {
    // Initialize logger
    logger := observability.NewLogger(observability.LoggerConfig{
        Level:  "info",
        Format: "text",
    })
    
    // Create analyzer
    analyzer := ai.NewCryptoCoinAnalyzer(logger)
    
    // Perform analysis
    ctx := context.Background()
    report, err := analyzer.AnalyzeCoin(ctx, "BTC")
    if err != nil {
        log.Fatal(err)
    }
    
    // Generate structured report
    markdown := analyzer.GenerateStructuredReport(report)
    fmt.Println(markdown)
}
```

## Report Format

The analyzer generates structured reports following this exact format:

```markdown
# CRYPTOCURRENCY ANALYSIS REPORT
Generated on: 2025-08-05 14:30:25 UTC
Symbol: BTC

## CURRENT MARKET DATA
- Price: $45,234.56 (+2.34%)
- Market Cap: $850.00B
- 24h Volume: $15.20B
- Circulating Supply: 19.50M

## RECENT NEWS & DEVELOPMENTS
- **Bitcoin Reaches New Monthly High** (Aug 4) 📈 - Bitcoin surged to new monthly high...

## MARKET SENTIMENT
- Overall Sentiment: Bullish 🐂
- Key Sentiment Drivers: Institutional adoption, Positive news coverage

## TECHNICAL INDICATORS
- Trend: Uptrend 📈
- Key Levels: Support at $43,500, Resistance at $47,000
- RSI: 58.0 (Neutral)
- Technical Outlook: Technical indicators suggest continued upward momentum

## FUNDAMENTAL INSIGHTS
- Project Status: Active development with regular updates
- Recent Updates:
  - Lightning Network improvements (Aug 3)
- Competitive Position: Market leader
- Development Activity: 150 commits, Stable trend

## SUMMARY & OUTLOOK
**Overall Outlook:** Bullish 🚀 (Confidence: 75%)

**Key Insights:**
• Strong institutional interest continues
• Technical indicators show positive momentum

**Risk Factors:**
⚠️ Cryptocurrency markets are highly volatile

**Time-based Analysis:**
• **Short-term (1-7 days):** Positive momentum expected to continue
• **Medium-term (1-3 months):** Fundamentals support higher prices
• **Long-term (6+ months):** Long-term outlook remains strong

---
*This analysis is for informational purposes only and does not constitute financial advice.*
```

## Configuration

### Environment Variables

- `LOG_LEVEL`: Logging level (debug, info, warn, error)
- `SEARCH_API_KEY`: API key for web search service (optional)
- `CACHE_DURATION`: Cache duration in minutes (default: 15)

### Supported Cryptocurrencies

The analyzer supports analysis for any cryptocurrency symbol, with optimized support for:

- BTC (Bitcoin)
- ETH (Ethereum)
- ADA (Cardano)
- SOL (Solana)
- MATIC (Polygon)
- LINK (Chainlink)
- UNI (Uniswap)
- AAVE (Aave)
- DOT (Polkadot)
- AVAX (Avalanche)
- And many others...

## Testing

Run the test suite:

```bash
# Run all tests
go test ./internal/ai/...

# Run with coverage
go test -v -coverprofile=coverage.out ./internal/ai/...
go tool cover -html=coverage.out -o coverage.html
```

## Performance

- **Cache Duration**: 15 minutes for analysis results
- **Timeout**: 60 seconds default (configurable)
- **Data Sources**: 5 sources per analysis as per requirements
- **Response Time**: Typically 5-15 seconds for fresh analysis

## Error Handling

The analyzer includes comprehensive error handling:

- Invalid symbol format validation
- Timeout handling for web searches
- Graceful degradation when data sources are unavailable
- Fallback to default values when specific data is missing
- Detailed error logging for debugging

## Security

- JWT authentication required for API endpoints
- Input validation and sanitization
- Rate limiting through middleware
- No sensitive data stored in cache
- Secure HTTP client configuration

## Monitoring

The analyzer integrates with the observability stack:

- Structured logging with correlation IDs
- Performance metrics tracking
- Health check endpoints
- Distributed tracing support
- Error rate monitoring

## Limitations

- Web search results are simulated for demonstration
- Real-time data depends on search result quality
- Analysis confidence varies with data availability
- Cache may serve stale data for up to 15 minutes
- Rate limits may apply to external data sources

## Future Enhancements

- Integration with real cryptocurrency APIs
- Advanced technical analysis indicators
- Historical price analysis
- Portfolio analysis capabilities
- Real-time WebSocket updates
- Machine learning sentiment analysis
- Custom alert configurations

## Support

For issues and questions:

1. Check the logs for detailed error information
2. Verify symbol format and network connectivity
3. Review the API documentation
4. Submit issues through the project repository

## License

This component is part of the AI Agentic Browser project and follows the same licensing terms.
