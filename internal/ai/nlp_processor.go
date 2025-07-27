package ai

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ai-agentic-browser/pkg/observability"
)

// NLPProcessor handles natural language processing for voice commands
type NLPProcessor struct {
	logger   *observability.Logger
	patterns map[CommandIntent][]IntentPattern
}

// IntentPattern represents a pattern for intent recognition
type IntentPattern struct {
	Regex      *regexp.Regexp
	Confidence float64
	Entities   []string
}

// NewNLPProcessor creates a new NLP processor
func NewNLPProcessor(logger *observability.Logger) *NLPProcessor {
	processor := &NLPProcessor{
		logger:   logger,
		patterns: make(map[CommandIntent][]IntentPattern),
	}
	
	processor.initializePatterns()
	return processor
}

// ProcessText processes text and extracts intent and entities
func (n *NLPProcessor) ProcessText(ctx context.Context, text string) (CommandIntent, map[string]interface{}, float64, error) {
	text = strings.ToLower(strings.TrimSpace(text))
	
	var bestIntent CommandIntent = IntentUnknown
	var bestConfidence float64 = 0.0
	var bestEntities map[string]interface{}

	// Try to match against all patterns
	for intent, patterns := range n.patterns {
		for _, pattern := range patterns {
			if pattern.Regex.MatchString(text) {
				if pattern.Confidence > bestConfidence {
					bestIntent = intent
					bestConfidence = pattern.Confidence
					bestEntities = n.extractEntities(text, pattern)
				}
			}
		}
	}

	n.logger.Info(ctx, "NLP processing completed", map[string]interface{}{
		"text":       text,
		"intent":     string(bestIntent),
		"confidence": bestConfidence,
		"entities":   bestEntities,
	})

	return bestIntent, bestEntities, bestConfidence, nil
}

// extractEntities extracts entities from text based on pattern
func (n *NLPProcessor) extractEntities(text string, pattern IntentPattern) map[string]interface{} {
	entities := make(map[string]interface{})
	
	// Extract common entities
	entities["amount"] = n.extractAmount(text)
	entities["token"] = n.extractToken(text)
	entities["portfolio_name"] = n.extractPortfolioName(text)
	entities["risk_level"] = n.extractRiskLevel(text)
	entities["strategy"] = n.extractStrategy(text)
	entities["timeframe"] = n.extractTimeframe(text)
	
	return entities
}

// extractAmount extracts monetary amounts from text
func (n *NLPProcessor) extractAmount(text string) string {
	// Pattern for amounts like "$1000", "1000 dollars", "1k", "1.5 eth"
	patterns := []string{
		`\$([0-9,]+(?:\.[0-9]+)?)`,
		`([0-9,]+(?:\.[0-9]+)?)\s*(?:dollars?|usd|bucks?)`,
		`([0-9,]+(?:\.[0-9]+)?)\s*k(?:\s|$)`,
		`([0-9,]+(?:\.[0-9]+)?)\s*(?:eth|bitcoin|btc|ether)`,
	}
	
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(text)
		if len(matches) > 1 {
			amount := strings.ReplaceAll(matches[1], ",", "")
			// Convert k to thousands
			if strings.Contains(text, "k") {
				if val, err := strconv.ParseFloat(amount, 64); err == nil {
					return fmt.Sprintf("%.2f", val*1000)
				}
			}
			return amount
		}
	}
	
	return ""
}

// extractToken extracts token names from text
func (n *NLPProcessor) extractToken(text string) string {
	tokens := map[string]string{
		"bitcoin":   "BTC",
		"btc":       "BTC",
		"ethereum":  "ETH",
		"eth":       "ETH",
		"ether":     "ETH",
		"usdc":      "USDC",
		"usdt":      "USDT",
		"tether":    "USDT",
		"dai":       "DAI",
		"polygon":   "MATIC",
		"matic":     "MATIC",
		"chainlink": "LINK",
		"link":      "LINK",
		"uniswap":   "UNI",
		"uni":       "UNI",
	}
	
	for token, symbol := range tokens {
		if strings.Contains(text, token) {
			return symbol
		}
	}
	
	return ""
}

// extractPortfolioName extracts portfolio names from text
func (n *NLPProcessor) extractPortfolioName(text string) string {
	// Look for quoted names or names after "called" or "named"
	patterns := []string{
		`"([^"]+)"`,
		`'([^']+)'`,
		`(?:called|named)\s+([a-zA-Z0-9\s]+?)(?:\s|$)`,
	}
	
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(text)
		if len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}
	
	return ""
}

// extractRiskLevel extracts risk levels from text
func (n *NLPProcessor) extractRiskLevel(text string) string {
	riskLevels := map[string]string{
		"conservative": "conservative",
		"safe":         "conservative",
		"low risk":     "conservative",
		"moderate":     "moderate",
		"medium":       "moderate",
		"balanced":     "moderate",
		"aggressive":   "aggressive",
		"high risk":    "aggressive",
		"risky":        "aggressive",
	}
	
	for phrase, level := range riskLevels {
		if strings.Contains(text, phrase) {
			return level
		}
	}
	
	return ""
}

// extractStrategy extracts trading strategies from text
func (n *NLPProcessor) extractStrategy(text string) string {
	strategies := map[string]string{
		"momentum":       "momentum",
		"trend":          "momentum",
		"mean reversion": "mean_reversion",
		"contrarian":     "mean_reversion",
		"arbitrage":      "arbitrage",
		"scalping":       "arbitrage",
	}
	
	for phrase, strategy := range strategies {
		if strings.Contains(text, phrase) {
			return strategy
		}
	}
	
	return ""
}

// extractTimeframe extracts timeframes from text
func (n *NLPProcessor) extractTimeframe(text string) string {
	timeframes := map[string]string{
		"daily":   "1d",
		"weekly":  "1w",
		"monthly": "1m",
		"hourly":  "1h",
		"minute":  "1min",
	}
	
	for phrase, timeframe := range timeframes {
		if strings.Contains(text, phrase) {
			return timeframe
		}
	}
	
	return ""
}

// initializePatterns initializes intent recognition patterns
func (n *NLPProcessor) initializePatterns() {
	// Portfolio creation patterns
	n.patterns[IntentCreatePortfolio] = []IntentPattern{
		{
			Regex:      regexp.MustCompile(`create\s+(?:a\s+)?(?:new\s+)?portfolio`),
			Confidence: 0.9,
		},
		{
			Regex:      regexp.MustCompile(`(?:make|start)\s+(?:a\s+)?(?:new\s+)?portfolio`),
			Confidence: 0.85,
		},
		{
			Regex:      regexp.MustCompile(`set\s+up\s+(?:a\s+)?portfolio`),
			Confidence: 0.8,
		},
	}

	// Buy token patterns
	n.patterns[IntentBuyToken] = []IntentPattern{
		{
			Regex:      regexp.MustCompile(`buy\s+(?:[0-9.]+\s+)?(?:bitcoin|btc|ethereum|eth|usdc)`),
			Confidence: 0.9,
		},
		{
			Regex:      regexp.MustCompile(`purchase\s+(?:[0-9.]+\s+)?(?:bitcoin|btc|ethereum|eth|usdc)`),
			Confidence: 0.85,
		},
		{
			Regex:      regexp.MustCompile(`get\s+(?:some\s+)?(?:[0-9.]+\s+)?(?:bitcoin|btc|ethereum|eth|usdc)`),
			Confidence: 0.7,
		},
	}

	// Sell token patterns
	n.patterns[IntentSellToken] = []IntentPattern{
		{
			Regex:      regexp.MustCompile(`sell\s+(?:[0-9.]+\s+)?(?:bitcoin|btc|ethereum|eth|usdc)`),
			Confidence: 0.9,
		},
		{
			Regex:      regexp.MustCompile(`(?:dispose|liquidate)\s+(?:[0-9.]+\s+)?(?:bitcoin|btc|ethereum|eth|usdc)`),
			Confidence: 0.8,
		},
	}

	// Check balance patterns
	n.patterns[IntentCheckBalance] = []IntentPattern{
		{
			Regex:      regexp.MustCompile(`(?:check|show|what.s)\s+(?:my\s+)?balance`),
			Confidence: 0.9,
		},
		{
			Regex:      regexp.MustCompile(`how\s+much\s+(?:money\s+)?(?:do\s+)?i\s+have`),
			Confidence: 0.85,
		},
		{
			Regex:      regexp.MustCompile(`(?:my\s+)?wallet\s+balance`),
			Confidence: 0.8,
		},
	}

	// Check portfolio patterns
	n.patterns[IntentCheckPortfolio] = []IntentPattern{
		{
			Regex:      regexp.MustCompile(`(?:check|show|view)\s+(?:my\s+)?portfolio`),
			Confidence: 0.9,
		},
		{
			Regex:      regexp.MustCompile(`portfolio\s+(?:status|summary|overview)`),
			Confidence: 0.85,
		},
		{
			Regex:      regexp.MustCompile(`how\s+(?:is\s+)?(?:my\s+)?portfolio\s+(?:doing|performing)`),
			Confidence: 0.8,
		},
	}

	// Trading control patterns
	n.patterns[IntentStartTrading] = []IntentPattern{
		{
			Regex:      regexp.MustCompile(`start\s+trading`),
			Confidence: 0.9,
		},
		{
			Regex:      regexp.MustCompile(`begin\s+(?:auto\s+)?trading`),
			Confidence: 0.85,
		},
		{
			Regex:      regexp.MustCompile(`turn\s+on\s+trading`),
			Confidence: 0.8,
		},
	}

	n.patterns[IntentStopTrading] = []IntentPattern{
		{
			Regex:      regexp.MustCompile(`stop\s+trading`),
			Confidence: 0.9,
		},
		{
			Regex:      regexp.MustCompile(`halt\s+(?:all\s+)?trading`),
			Confidence: 0.85,
		},
		{
			Regex:      regexp.MustCompile(`turn\s+off\s+trading`),
			Confidence: 0.8,
		},
	}

	// Market data patterns
	n.patterns[IntentGetMarketData] = []IntentPattern{
		{
			Regex:      regexp.MustCompile(`(?:what.s\s+the\s+price\s+of|price\s+of)\s+(?:bitcoin|btc|ethereum|eth)`),
			Confidence: 0.9,
		},
		{
			Regex:      regexp.MustCompile(`(?:show|get)\s+(?:market\s+)?(?:data|prices)`),
			Confidence: 0.8,
		},
		{
			Regex:      regexp.MustCompile(`how\s+much\s+is\s+(?:bitcoin|btc|ethereum|eth)`),
			Confidence: 0.85,
		},
	}

	// Strategy patterns
	n.patterns[IntentSetStrategy] = []IntentPattern{
		{
			Regex:      regexp.MustCompile(`(?:set|use|enable)\s+(?:the\s+)?(?:momentum|mean\s+reversion|arbitrage)\s+strategy`),
			Confidence: 0.9,
		},
		{
			Regex:      regexp.MustCompile(`switch\s+to\s+(?:momentum|mean\s+reversion|arbitrage)`),
			Confidence: 0.85,
		},
	}

	// Rebalancing patterns
	n.patterns[IntentRebalancePortfolio] = []IntentPattern{
		{
			Regex:      regexp.MustCompile(`rebalance\s+(?:my\s+)?portfolio`),
			Confidence: 0.9,
		},
		{
			Regex:      regexp.MustCompile(`(?:adjust|optimize)\s+(?:my\s+)?portfolio`),
			Confidence: 0.8,
		},
	}

	// Risk checking patterns
	n.patterns[IntentCheckRisk] = []IntentPattern{
		{
			Regex:      regexp.MustCompile(`(?:check|assess|analyze)\s+(?:the\s+)?risk`),
			Confidence: 0.9,
		},
		{
			Regex:      regexp.MustCompile(`how\s+(?:risky|safe)\s+is\s+this`),
			Confidence: 0.85,
		},
	}

	// Yield finding patterns
	n.patterns[IntentFindYield] = []IntentPattern{
		{
			Regex:      regexp.MustCompile(`(?:find|show|get)\s+(?:yield\s+)?(?:opportunities|farming)`),
			Confidence: 0.9,
		},
		{
			Regex:      regexp.MustCompile(`(?:best\s+)?(?:apy|yield|returns?)`),
			Confidence: 0.8,
		},
		{
			Regex:      regexp.MustCompile(`defi\s+(?:opportunities|yields?)`),
			Confidence: 0.85,
		},
	}

	// Help patterns
	n.patterns[IntentHelp] = []IntentPattern{
		{
			Regex:      regexp.MustCompile(`help`),
			Confidence: 0.95,
		},
		{
			Regex:      regexp.MustCompile(`what\s+can\s+(?:you\s+)?(?:do|i\s+say)`),
			Confidence: 0.9,
		},
		{
			Regex:      regexp.MustCompile(`(?:show\s+)?(?:commands|options)`),
			Confidence: 0.85,
		},
	}
}
