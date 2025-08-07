# Crypto Analyzer CLI

A command-line interface for the AI-powered cryptocurrency analysis tool. This CLI provides comprehensive cryptocurrency analysis using 5 different data sources and generates structured reports following the exact specifications.

## Installation

### Build from Source

```bash
# Clone the repository
git clone https://github.com/DimaJoyti/ai-agentic-crypto-browser.git
cd ai-agentic-crypto-browser

# Build the CLI tool
make build-crypto-analyzer

# The binary will be available at bin/crypto-analyzer
```

### Direct Build

```bash
go build -o crypto-analyzer ./cmd/crypto-analyzer
```

## Usage

### Basic Usage

```bash
# Analyze Bitcoin with markdown output (default)
./crypto-analyzer -symbol BTC

# Analyze Ethereum with JSON output
./crypto-analyzer -symbol ETH -format json

# Verbose analysis with custom timeout
./crypto-analyzer -symbol BTC -verbose -timeout 120s
```

### Command Line Options

| Flag | Description | Default | Required |
|------|-------------|---------|----------|
| `-symbol` | Cryptocurrency symbol to analyze (e.g., BTC, ETH) | - | Yes |
| `-format` | Output format: `markdown` or `json` | `markdown` | No |
| `-verbose` | Enable verbose logging | `false` | No |
| `-timeout` | Analysis timeout duration | `60s` | No |
| `-help` | Show help message | - | No |
| `-version` | Show version information | - | No |

### Examples

#### 1. Basic Bitcoin Analysis

```bash
./crypto-analyzer -symbol BTC
```

Output:
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
- **Bitcoin Reaches New Monthly High** (Aug 4) 📈 - Bitcoin surged...

...
```

#### 2. Ethereum Analysis with JSON Output

```bash
./crypto-analyzer -symbol ETH -format json
```

Output:
```json
{
  "timestamp": "2025-08-05T14:30:25Z",
  "symbol": "ETH",
  "current_data": {
    "price": "2567.89",
    "change_percent_24h": "3.12",
    "market_cap": "300000000000",
    "volume_24h": "8500000000"
  },
  "market_sentiment": {
    "overall_sentiment": "bullish",
    "sentiment_score": "0.6"
  },
  ...
}
```

#### 3. Verbose Analysis with Custom Timeout

```bash
./crypto-analyzer -symbol ADA -verbose -timeout 2m
```

This will show detailed logging information during the analysis process.

#### 4. Quick Analysis

```bash
./crypto-analyzer -symbol SOL -timeout 30s
```

Performs a faster analysis with a 30-second timeout.

## Supported Cryptocurrencies

The CLI supports analysis for any cryptocurrency symbol. Common symbols include:

| Symbol | Name | Symbol | Name |
|--------|------|--------|------|
| BTC | Bitcoin | ETH | Ethereum |
| ADA | Cardano | SOL | Solana |
| MATIC | Polygon | LINK | Chainlink |
| UNI | Uniswap | AAVE | Aave |
| DOT | Polkadot | AVAX | Avalanche |
| ATOM | Cosmos | ALGO | Algorand |
| XTZ | Tezos | FIL | Filecoin |
| ICP | Internet Computer | NEAR | NEAR Protocol |

## Output Formats

### Markdown Format (Default)

The markdown format provides a structured, human-readable report with:

- Current market data with price and percentage changes
- Recent news and developments with impact indicators
- Market sentiment analysis with emojis
- Technical indicators with trend analysis
- Fundamental insights about the project
- Summary and outlook with time-based analysis
- Risk factors and opportunities

### JSON Format

The JSON format provides structured data suitable for:

- Integration with other tools
- Programmatic processing
- Data analysis and visualization
- API consumption

## Error Handling

The CLI includes comprehensive error handling:

### Common Errors

1. **Missing Symbol**
   ```bash
   Error: symbol is required
   ```
   Solution: Provide a valid cryptocurrency symbol using `-symbol`

2. **Invalid Format**
   ```bash
   Error: format must be 'markdown' or 'json'
   ```
   Solution: Use either `markdown` or `json` for the `-format` flag

3. **Analysis Timeout**
   ```bash
   Analysis failed: context deadline exceeded
   ```
   Solution: Increase timeout using `-timeout` flag (e.g., `-timeout 120s`)

4. **Invalid Symbol Format**
   ```bash
   Error: symbol must be between 2 and 10 characters
   ```
   Solution: Provide a valid cryptocurrency symbol

### Troubleshooting

1. **Network Issues**: Ensure internet connectivity for data fetching
2. **Slow Analysis**: Increase timeout or check network speed
3. **Unknown Symbol**: Verify the cryptocurrency symbol exists
4. **Permission Issues**: Ensure the binary has execute permissions

## Performance

- **Analysis Time**: Typically 5-15 seconds
- **Data Sources**: Uses 5 different sources per analysis
- **Timeout**: Default 60 seconds (configurable)
- **Cache**: Results cached for 15 minutes
- **Memory Usage**: Minimal memory footprint

## Integration

### Shell Scripts

```bash
#!/bin/bash
# analyze_portfolio.sh

SYMBOLS=("BTC" "ETH" "ADA" "SOL")

for symbol in "${SYMBOLS[@]}"; do
    echo "Analyzing $symbol..."
    ./crypto-analyzer -symbol "$symbol" > "reports/${symbol}_analysis.md"
done
```

### Automation

```bash
# Cron job for daily analysis
0 9 * * * /path/to/crypto-analyzer -symbol BTC > /path/to/daily_btc_report.md
```

### CI/CD Integration

```yaml
# GitHub Actions example
- name: Analyze Cryptocurrency
  run: |
    ./crypto-analyzer -symbol BTC -format json > btc_analysis.json
    # Process the analysis results
```

## Development

### Building

```bash
# Build for current platform
go build -o crypto-analyzer ./cmd/crypto-analyzer

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o crypto-analyzer-linux ./cmd/crypto-analyzer
GOOS=windows GOARCH=amd64 go build -o crypto-analyzer.exe ./cmd/crypto-analyzer
GOOS=darwin GOARCH=amd64 go build -o crypto-analyzer-mac ./cmd/crypto-analyzer
```

### Testing

```bash
# Run tests
go test ./cmd/crypto-analyzer/...

# Test the binary
./crypto-analyzer -symbol BTC -timeout 30s
```

## Configuration

### Environment Variables

- `LOG_LEVEL`: Set logging level (debug, info, warn, error)
- `CRYPTO_ANALYZER_TIMEOUT`: Default timeout duration
- `CRYPTO_ANALYZER_FORMAT`: Default output format

### Example Configuration

```bash
export LOG_LEVEL=debug
export CRYPTO_ANALYZER_TIMEOUT=120s
export CRYPTO_ANALYZER_FORMAT=json

./crypto-analyzer -symbol BTC
```

## Help and Support

### Getting Help

```bash
# Show help message
./crypto-analyzer -help

# Show version
./crypto-analyzer -version

# Show supported symbols (if implemented)
./crypto-analyzer -list-symbols
```

### Common Use Cases

1. **Daily Market Analysis**: Automated daily reports for portfolio tracking
2. **Research**: In-depth analysis for investment research
3. **Monitoring**: Regular monitoring of specific cryptocurrencies
4. **Integration**: Data source for other applications and tools

## License

This CLI tool is part of the AI Agentic Browser project and follows the same licensing terms.

## Contributing

Contributions are welcome! Please see the main project repository for contribution guidelines.

## Changelog

### v1.0.0
- Initial release
- Support for comprehensive cryptocurrency analysis
- Markdown and JSON output formats
- Configurable timeout and verbose logging
- Error handling and validation
