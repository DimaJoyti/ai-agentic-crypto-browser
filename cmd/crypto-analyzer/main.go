package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ai-agentic-browser/internal/ai"
	"github.com/ai-agentic-browser/internal/config"
	"github.com/ai-agentic-browser/pkg/observability"
)

// CLI flags
var (
	symbol  = flag.String("symbol", "", "Cryptocurrency symbol to analyze (e.g., BTC, ETH)")
	format  = flag.String("format", "markdown", "Output format: markdown, json")
	verbose = flag.Bool("verbose", false, "Enable verbose logging")
	timeout = flag.Duration("timeout", 60*time.Second, "Analysis timeout")
	help    = flag.Bool("help", false, "Show help message")
	version = flag.Bool("version", false, "Show version information")
)

const (
	appName    = "crypto-analyzer"
	appVersion = "1.0.0"
	appDesc    = "AI-powered cryptocurrency analysis tool"
)

func main() {
	flag.Parse()

	// Show help
	if *help {
		showHelp()
		return
	}

	// Show version
	if *version {
		showVersion()
		return
	}

	// Validate required flags
	if *symbol == "" {
		fmt.Fprintf(os.Stderr, "Error: symbol is required\n\n")
		showUsage()
		os.Exit(1)
	}

	// Validate format
	if *format != "markdown" && *format != "json" {
		fmt.Fprintf(os.Stderr, "Error: format must be 'markdown' or 'json'\n\n")
		showUsage()
		os.Exit(1)
	}

	// Initialize logger
	loggerConfig := config.ObservabilityConfig{
		ServiceName: "crypto-analyzer",
		LogLevel:    "info",
		LogFormat:   "text",
	}
	if *verbose {
		loggerConfig.LogLevel = "debug"
	}
	logger := observability.NewLogger(loggerConfig)

	// Create crypto coin analyzer
	analyzer := ai.NewCryptoCoinAnalyzer(logger)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	// Normalize symbol
	symbolUpper := strings.ToUpper(strings.TrimSpace(*symbol))

	// Show analysis start message
	if *verbose {
		fmt.Fprintf(os.Stderr, "Starting analysis for %s...\n", symbolUpper)
		fmt.Fprintf(os.Stderr, "Using 5 data sources as per requirements...\n")
	}

	// Perform analysis
	var output string
	var err error

	switch *format {
	case "json":
		report, analysisErr := analyzer.AnalyzeCoin(ctx, symbolUpper)
		if analysisErr != nil {
			err = analysisErr
		} else {
			output, err = formatJSON(report)
		}
	case "markdown":
		output, err = analyzer.AnalyzeCoinWithStructuredReport(ctx, symbolUpper)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Analysis failed: %v\n", err)
		os.Exit(1)
	}

	// Output result
	fmt.Print(output)

	if *verbose {
		fmt.Fprintf(os.Stderr, "\nAnalysis completed successfully!\n")
	}
}

func showHelp() {
	fmt.Printf("%s - %s\n\n", appName, appDesc)
	fmt.Printf("USAGE:\n")
	fmt.Printf("  %s -symbol <SYMBOL> [OPTIONS]\n\n", appName)
	fmt.Printf("DESCRIPTION:\n")
	fmt.Printf("  Performs comprehensive cryptocurrency analysis using AI-powered tools.\n")
	fmt.Printf("  Gathers data from 5 different sources including price data, news,\n")
	fmt.Printf("  sentiment analysis, technical indicators, and fundamental analysis.\n\n")
	fmt.Printf("EXAMPLES:\n")
	fmt.Printf("  # Analyze Bitcoin with markdown output\n")
	fmt.Printf("  %s -symbol BTC\n\n", appName)
	fmt.Printf("  # Analyze Ethereum with JSON output\n")
	fmt.Printf("  %s -symbol ETH -format json\n\n", appName)
	fmt.Printf("  # Verbose analysis with custom timeout\n")
	fmt.Printf("  %s -symbol BTC -verbose -timeout 120s\n\n", appName)
	fmt.Printf("OPTIONS:\n")
	flag.PrintDefaults()
	fmt.Printf("\nSUPPORTED SYMBOLS:\n")
	fmt.Printf("  BTC, ETH, ADA, SOL, MATIC, LINK, UNI, AAVE, and many others\n")
	fmt.Printf("\nOUTPUT FORMATS:\n")
	fmt.Printf("  markdown - Structured markdown report (default)\n")
	fmt.Printf("  json     - JSON format with detailed data\n")
	fmt.Printf("\nNOTE:\n")
	fmt.Printf("  This tool provides analysis for informational purposes only.\n")
	fmt.Printf("  It does not constitute financial advice.\n")
}

func showVersion() {
	fmt.Printf("%s version %s\n", appName, appVersion)
}

func showUsage() {
	fmt.Printf("Usage: %s -symbol <SYMBOL> [OPTIONS]\n", appName)
	fmt.Printf("Use '%s -help' for more information.\n", appName)
}

func formatJSON(report *ai.CoinAnalysisReport) (string, error) {
	// Use a simple JSON marshaling approach
	jsonBytes, err := json.Marshal(report)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return string(jsonBytes), nil
}

// Example usage function for demonstration
func showExamples() {
	examples := []struct {
		description string
		command     string
		explanation string
	}{
		{
			description: "Basic Bitcoin Analysis",
			command:     "./crypto-analyzer -symbol BTC",
			explanation: "Performs comprehensive Bitcoin analysis with markdown output",
		},
		{
			description: "Ethereum Analysis with JSON Output",
			command:     "./crypto-analyzer -symbol ETH -format json",
			explanation: "Analyzes Ethereum and outputs detailed JSON data",
		},
		{
			description: "Verbose Analysis with Custom Timeout",
			command:     "./crypto-analyzer -symbol ADA -verbose -timeout 2m",
			explanation: "Analyzes Cardano with verbose logging and 2-minute timeout",
		},
		{
			description: "Quick Analysis",
			command:     "./crypto-analyzer -symbol SOL -timeout 30s",
			explanation: "Fast Solana analysis with 30-second timeout",
		},
	}

	fmt.Printf("EXAMPLE USAGE:\n\n")
	for i, example := range examples {
		fmt.Printf("%d. %s\n", i+1, example.description)
		fmt.Printf("   Command: %s\n", example.command)
		fmt.Printf("   %s\n\n", example.explanation)
	}
}

// validateSymbol validates cryptocurrency symbol format
func validateSymbol(symbol string) error {
	if len(symbol) < 2 || len(symbol) > 10 {
		return fmt.Errorf("symbol must be between 2 and 10 characters")
	}

	// Check for valid characters (letters and numbers only)
	for _, char := range symbol {
		if !((char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
			return fmt.Errorf("symbol must contain only letters and numbers")
		}
	}

	return nil
}

// Common cryptocurrency symbols for validation
var commonSymbols = map[string]string{
	"BTC":   "Bitcoin",
	"ETH":   "Ethereum",
	"ADA":   "Cardano",
	"SOL":   "Solana",
	"MATIC": "Polygon",
	"LINK":  "Chainlink",
	"UNI":   "Uniswap",
	"AAVE":  "Aave",
	"DOT":   "Polkadot",
	"AVAX":  "Avalanche",
	"ATOM":  "Cosmos",
	"ALGO":  "Algorand",
	"XTZ":   "Tezos",
	"FIL":   "Filecoin",
	"ICP":   "Internet Computer",
	"NEAR":  "NEAR Protocol",
	"FLOW":  "Flow",
	"EGLD":  "MultiversX",
	"SAND":  "The Sandbox",
	"MANA":  "Decentraland",
}

// isKnownSymbol checks if the symbol is in our known list
func isKnownSymbol(symbol string) (string, bool) {
	name, exists := commonSymbols[strings.ToUpper(symbol)]
	return name, exists
}

// showSupportedSymbols displays list of commonly supported symbols
func showSupportedSymbols() {
	fmt.Printf("COMMONLY SUPPORTED CRYPTOCURRENCY SYMBOLS:\n\n")
	for symbol, name := range commonSymbols {
		fmt.Printf("  %-6s - %s\n", symbol, name)
	}
	fmt.Printf("\nNote: Many other symbols are also supported.\n")
}
