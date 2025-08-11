# AI-Agentic Crypto Browser Terminal User Guide

## Overview

The AI-Agentic Crypto Browser Terminal is a powerful command-line interface that provides comprehensive access to all platform features including trading, market analysis, AI-powered insights, and Web3 operations.

## Getting Started

### Accessing the Terminal

1. **Web Interface**: Navigate to `/terminal` in your browser
2. **Direct Connection**: Connect via WebSocket to `ws://localhost:8085/ws`
3. **API Access**: Use REST endpoints at `http://localhost:8085/api/v1/`

### Basic Usage

```bash
# Get help
help

# Check system status
status

# View available commands
help

# Get help for specific command
help <command>
```

## Command Categories

### System Commands

#### `status [--verbose] [--json] [--services]`
Display system health and service status.

```bash
status                    # Basic status
status --verbose          # Detailed information
status --json            # JSON format
status --services         # Services only
```

#### `config [get|set|list] [key] [value]`
Manage terminal configuration.

```bash
config list               # Show all settings
config get theme          # Get specific setting
config set theme dark     # Set configuration value
```

#### `alias [name] [command]`
Create and manage command aliases.

```bash
alias                     # List all aliases
alias p "price BTC"       # Create alias
alias status-v "status --verbose"
```

#### `history [--limit 20] [--search pattern]`
View command execution history.

```bash
history                   # Recent commands
history --limit 50        # Show more commands
history --search price    # Search history
```

### Trading Commands

#### `buy <symbol> <amount> [--limit <price>] [--market] [--stop <price>]`
Place buy orders.

```bash
buy BTC 0.1               # Market buy
buy ETH 1.5 --limit 3200  # Limit buy
buy SOL 10 --stop 95      # Stop buy
```

#### `sell <symbol> <amount> [--limit <price>] [--market] [--stop <price>]`
Place sell orders.

```bash
sell BTC 0.05             # Market sell
sell ETH 1.0 --limit 3300 # Limit sell
sell ADA 100 --stop 1.20  # Stop sell
```

#### `portfolio [--detailed]`
View portfolio holdings and performance.

```bash
portfolio                 # Basic portfolio view
portfolio --detailed      # Detailed breakdown
```

#### `orders [--all]`
View active orders.

```bash
orders                    # Summary
orders --all              # All orders with details
```

#### `balance [symbol] [--detailed] [--json]`
Check account balances.

```bash
balance                   # All balances
balance BTC               # Specific asset
balance --detailed        # Detailed view
balance --json            # JSON format
```

### Market Data Commands

#### `price <symbol> [base_currency] [--format json|table|chart] [--watch]`
Get cryptocurrency prices.

```bash
price BTC                 # Bitcoin price
price ETH USD             # Ethereum in USD
price BTC --format json   # JSON format
price BTC --watch         # Live updates
```

#### `chart <symbol> [timeframe] [--ascii]`
Display price charts.

```bash
chart BTC                 # Default chart
chart ETH 4h              # 4-hour timeframe
chart BTC --ascii         # ASCII chart
```

#### `news [symbol] [--limit 10] [--category market|tech|regulation]`
Get latest cryptocurrency news.

```bash
news                      # General news
news BTC                  # Bitcoin-specific news
news --category tech      # Technology news
news --limit 20           # More articles
```

### AI Commands

#### `analyze <symbol> [--timeframe 1h|4h|1d|1w] [--depth basic|detailed|comprehensive]`
AI-powered cryptocurrency analysis.

```bash
analyze BTC               # Basic analysis
analyze ETH --timeframe 1d # Daily analysis
analyze BTC --depth comprehensive
```

#### `predict <symbol> [--horizon 1h|4h|1d|1w|1m] [--model lstm|transformer|ensemble]`
Generate AI price predictions.

```bash
predict BTC               # Default prediction
predict ETH --horizon 1w  # Weekly prediction
predict BTC --model ensemble
```

#### `chat <message> [--model gpt4|claude|local]`
Chat with AI assistant.

```bash
chat "What's the outlook for Bitcoin?"
chat "Explain DeFi yield farming" --model gpt4
```

### Advanced Features

#### `script <file> [--args arg1,arg2] | --list | --create <name>`
Execute and manage scripts.

```bash
script --list             # List available scripts
script --create daily     # Create new script
script daily-check        # Run script
```

#### `watch <command> [--interval 5s] [--count 10]`
Monitor commands with real-time updates.

```bash
watch price BTC           # Watch Bitcoin price
watch portfolio --interval 10s
watch "status --services" --count 5
```

#### `export <type> [--format csv|json|pdf] [--output file.ext]`
Export data and reports.

```bash
export portfolio          # Export portfolio
export history --format csv
export trades --output report.pdf
```

## Session Management

### Creating Sessions

Sessions are automatically created when you connect to the terminal. Each session maintains:

- Command history
- Configuration settings
- Environment variables
- Command aliases

### Session Persistence

Sessions persist across browser refreshes and reconnections. Use session IDs to resume previous sessions.

## Configuration

### Environment Variables

Set environment variables for your session:

```bash
config set API_KEY "your-api-key"
config set THEME "dark"
config set AUTO_SAVE "true"
```

### Aliases

Create shortcuts for frequently used commands:

```bash
alias p "price BTC"
alias s "status --verbose"
alias pf "portfolio --detailed"
```

## Tips and Tricks

### Command History Navigation

- Use ↑/↓ arrow keys to navigate command history
- Use `Ctrl+R` to search command history
- Use `history --search <pattern>` to find specific commands

### Autocomplete

- Press `Tab` for command and parameter completion
- Type partial commands and press `Tab` for suggestions
- Use `help <command>` for detailed usage information

### Output Formatting

Many commands support multiple output formats:

```bash
price BTC --format json   # Machine-readable JSON
status --json             # JSON status output
balance --detailed        # Human-readable detailed view
```

### Batch Operations

Use scripts for batch operations:

```bash
# Create a script file
script --create morning-check

# Add commands to the script:
# status --services
# portfolio
# price BTC ETH ADA
# news --limit 5

# Run the script
script morning-check
```

## Keyboard Shortcuts

- `Ctrl+C`: Cancel current command
- `Ctrl+L`: Clear screen (same as `clear`)
- `Ctrl+R`: Search command history
- `Tab`: Autocomplete
- `↑/↓`: Navigate command history
- `Ctrl+A`: Move to beginning of line
- `Ctrl+E`: Move to end of line

## Error Handling

The terminal provides detailed error messages and suggestions:

```bash
$ invalid-command
Command not found: invalid-command
Did you mean: status, config, help?

$ buy
Usage: buy <symbol> <amount> [--limit <price>] [--market] [--stop <price>]
Example: buy BTC 0.1 --limit 45000
```

## Integration with Services

The terminal integrates with all platform services:

- **AI Service**: Powers analysis and predictions
- **Trading Service**: Executes orders and manages portfolio
- **Web3 Service**: Handles blockchain operations
- **Browser Service**: Automates web interactions
- **Auth Service**: Manages authentication and permissions

## Troubleshooting

### Connection Issues

```bash
# Check service status
status --services

# Verify connection
help
```

### Command Errors

```bash
# Get command help
help <command>

# Check command history for syntax
history --search <command>
```

### Performance Issues

```bash
# Check system status
status --verbose

# Clear command history
history --clear
```

## API Reference

For programmatic access, see the [Terminal API Documentation](TERMINAL_API.md).

## Examples

### Daily Trading Routine

```bash
# Morning check
status
portfolio --detailed
news --limit 10

# Market analysis
analyze BTC --timeframe 1d
analyze ETH --timeframe 1d
price BTC ETH ADA SOL

# Place orders based on analysis
buy BTC 0.1 --limit 44500
sell ETH 0.5 --limit 3300

# Monitor positions
watch portfolio --interval 30s
```

### Research Workflow

```bash
# Research a new cryptocurrency
price LINK
analyze LINK --depth comprehensive
news LINK --limit 15
chat "What are the fundamentals of Chainlink?"

# Export research
export analysis --format pdf --output LINK_research.pdf
```

### Portfolio Management

```bash
# Review portfolio
portfolio --detailed
balance --detailed
orders --all

# Rebalance
sell BTC 0.2
buy ETH 1.0
buy ADA 200

# Track performance
export portfolio --format csv
watch portfolio --interval 60s --count 10
```

For more advanced usage and API integration, see the [Terminal Developer Guide](TERMINAL_DEVELOPER_GUIDE.md).
