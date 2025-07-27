package ai

import (
	"context"
	"fmt"
	"time"

	"github.com/ai-agentic-browser/internal/web3"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// VoiceInterface handles voice command processing and natural language understanding
type VoiceInterface struct {
	logger         *observability.Logger
	tradingEngine  *web3.TradingEngine
	defiManager    *web3.DeFiProtocolManager
	riskAssessment *web3.RiskAssessmentService
	nlpProcessor   *NLPProcessor
	commandHistory []VoiceCommand
	config         VoiceConfig
}

// VoiceConfig holds configuration for voice interface
type VoiceConfig struct {
	Language            string        `json:"language"`
	ConfidenceThreshold float64       `json:"confidence_threshold"`
	MaxCommandHistory   int           `json:"max_command_history"`
	ResponseTimeout     time.Duration `json:"response_timeout"`
	EnableSafetyMode    bool          `json:"enable_safety_mode"`
	RequireConfirmation bool          `json:"require_confirmation"`
}

// VoiceCommand represents a processed voice command
type VoiceCommand struct {
	ID         uuid.UUID              `json:"id"`
	UserID     uuid.UUID              `json:"user_id"`
	RawText    string                 `json:"raw_text"`
	Intent     CommandIntent          `json:"intent"`
	Entities   map[string]interface{} `json:"entities"`
	Confidence float64                `json:"confidence"`
	Status     CommandStatus          `json:"status"`
	Response   string                 `json:"response"`
	ExecutedAt time.Time              `json:"executed_at"`
	Duration   time.Duration          `json:"duration"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// CommandIntent represents the intent of a voice command
type CommandIntent string

const (
	IntentCreatePortfolio    CommandIntent = "create_portfolio"
	IntentBuyToken           CommandIntent = "buy_token"
	IntentSellToken          CommandIntent = "sell_token"
	IntentCheckBalance       CommandIntent = "check_balance"
	IntentCheckPortfolio     CommandIntent = "check_portfolio"
	IntentSetStrategy        CommandIntent = "set_strategy"
	IntentStopTrading        CommandIntent = "stop_trading"
	IntentStartTrading       CommandIntent = "start_trading"
	IntentGetMarketData      CommandIntent = "get_market_data"
	IntentRebalancePortfolio CommandIntent = "rebalance_portfolio"
	IntentCheckRisk          CommandIntent = "check_risk"
	IntentFindYield          CommandIntent = "find_yield"
	IntentHelp               CommandIntent = "help"
	IntentUnknown            CommandIntent = "unknown"
)

// CommandStatus represents the status of command execution
type CommandStatus string

const (
	StatusPending   CommandStatus = "pending"
	StatusExecuting CommandStatus = "executing"
	StatusCompleted CommandStatus = "completed"
	StatusFailed    CommandStatus = "failed"
	StatusCancelled CommandStatus = "cancelled"
)

// VoiceResponse represents a response to a voice command
type VoiceResponse struct {
	Text       string                 `json:"text"`
	AudioURL   string                 `json:"audio_url,omitempty"`
	Data       interface{}            `json:"data,omitempty"`
	Actions    []SuggestedAction      `json:"actions,omitempty"`
	Confidence float64                `json:"confidence"`
	Duration   time.Duration          `json:"duration"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// SuggestedAction represents a suggested follow-up action
type SuggestedAction struct {
	Text        string `json:"text"`
	Command     string `json:"command"`
	Description string `json:"description"`
}

// NewVoiceInterface creates a new voice interface
func NewVoiceInterface(
	logger *observability.Logger,
	tradingEngine *web3.TradingEngine,
	defiManager *web3.DeFiProtocolManager,
	riskAssessment *web3.RiskAssessmentService,
) *VoiceInterface {
	config := VoiceConfig{
		Language:            "en-US",
		ConfidenceThreshold: 0.7,
		MaxCommandHistory:   100,
		ResponseTimeout:     30 * time.Second,
		EnableSafetyMode:    true,
		RequireConfirmation: true,
	}

	return &VoiceInterface{
		logger:         logger,
		tradingEngine:  tradingEngine,
		defiManager:    defiManager,
		riskAssessment: riskAssessment,
		nlpProcessor:   NewNLPProcessor(logger),
		commandHistory: make([]VoiceCommand, 0),
		config:         config,
	}
}

// ProcessVoiceCommand processes a voice command and returns a response
func (v *VoiceInterface) ProcessVoiceCommand(ctx context.Context, userID uuid.UUID, audioData []byte, text string) (*VoiceResponse, error) {
	startTime := time.Now()

	// Create command record
	command := VoiceCommand{
		ID:         uuid.New(),
		UserID:     userID,
		RawText:    text,
		Status:     StatusPending,
		ExecutedAt: startTime,
		Metadata:   make(map[string]interface{}),
	}

	// Process natural language
	intent, entities, confidence, err := v.nlpProcessor.ProcessText(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("NLP processing failed: %w", err)
	}

	command.Intent = intent
	command.Entities = entities
	command.Confidence = confidence
	command.Status = StatusExecuting

	// Check confidence threshold
	if confidence < v.config.ConfidenceThreshold {
		command.Status = StatusFailed
		command.Response = "I'm not confident I understood that correctly. Could you please rephrase?"
		v.addToHistory(command)

		return &VoiceResponse{
			Text:       command.Response,
			Confidence: confidence,
			Duration:   time.Since(startTime),
		}, nil
	}

	// Execute command based on intent
	response, err := v.executeCommand(ctx, command)
	if err != nil {
		command.Status = StatusFailed
		command.Response = fmt.Sprintf("Command execution failed: %s", err.Error())
		v.addToHistory(command)
		return nil, err
	}

	command.Status = StatusCompleted
	command.Response = response.Text
	command.Duration = time.Since(startTime)
	v.addToHistory(command)

	response.Duration = command.Duration
	return response, nil
}

// executeCommand executes a command based on its intent
func (v *VoiceInterface) executeCommand(ctx context.Context, command VoiceCommand) (*VoiceResponse, error) {
	switch command.Intent {
	case IntentCreatePortfolio:
		return v.handleCreatePortfolio(ctx, command)
	case IntentBuyToken:
		return v.handleBuyToken(ctx, command)
	case IntentSellToken:
		return v.handleSellToken(ctx, command)
	case IntentCheckBalance:
		return v.handleCheckBalance(ctx, command)
	case IntentCheckPortfolio:
		return v.handleCheckPortfolio(ctx, command)
	case IntentSetStrategy:
		return v.handleSetStrategy(ctx, command)
	case IntentStartTrading:
		return v.handleStartTrading(ctx, command)
	case IntentStopTrading:
		return v.handleStopTrading(ctx, command)
	case IntentGetMarketData:
		return v.handleGetMarketData(ctx, command)
	case IntentRebalancePortfolio:
		return v.handleRebalancePortfolio(ctx, command)
	case IntentCheckRisk:
		return v.handleCheckRisk(ctx, command)
	case IntentFindYield:
		return v.handleFindYield(ctx, command)
	case IntentHelp:
		return v.handleHelp(ctx, command)
	default:
		return &VoiceResponse{
			Text: "I didn't understand that command. Try saying 'help' to see what I can do.",
			Actions: []SuggestedAction{
				{Text: "Get Help", Command: "help", Description: "Show available commands"},
				{Text: "Check Portfolio", Command: "check my portfolio", Description: "View portfolio status"},
			},
		}, nil
	}
}

// handleCreatePortfolio handles portfolio creation commands
func (v *VoiceInterface) handleCreatePortfolio(ctx context.Context, command VoiceCommand) (*VoiceResponse, error) {
	name, _ := command.Entities["portfolio_name"].(string)
	if name == "" {
		name = "Voice Created Portfolio"
	}

	balanceStr, _ := command.Entities["amount"].(string)
	balance := decimal.NewFromInt(10000) // Default $10k
	if balanceStr != "" {
		if parsed, err := decimal.NewFromString(balanceStr); err == nil {
			balance = parsed
		}
	}

	riskLevel, _ := command.Entities["risk_level"].(string)
	if riskLevel == "" {
		riskLevel = "moderate"
	}

	// Create risk profile based on voice input
	riskProfile := web3.RiskProfile{
		Level:                riskLevel,
		MaxPositionSize:      decimal.NewFromFloat(0.1),
		MaxDailyLoss:         decimal.NewFromFloat(0.05),
		StopLossPercentage:   decimal.NewFromFloat(0.1),
		TakeProfitPercentage: decimal.NewFromFloat(0.2),
		AllowedStrategies:    []string{"momentum", "mean_reversion"},
	}

	portfolio, err := v.tradingEngine.CreatePortfolio(ctx, command.UserID, name, balance, riskProfile)
	if err != nil {
		return nil, fmt.Errorf("failed to create portfolio: %w", err)
	}

	return &VoiceResponse{
		Text: fmt.Sprintf("Successfully created portfolio '%s' with $%s initial balance and %s risk level. Portfolio ID: %s",
			name, balance.String(), riskLevel, portfolio.ID.String()),
		Data: portfolio,
		Actions: []SuggestedAction{
			{Text: "Start Trading", Command: "start trading", Description: "Begin autonomous trading"},
			{Text: "Check Portfolio", Command: "check my portfolio", Description: "View portfolio details"},
		},
		Confidence: command.Confidence,
	}, nil
}

// handleCheckPortfolio handles portfolio status commands
func (v *VoiceInterface) handleCheckPortfolio(ctx context.Context, command VoiceCommand) (*VoiceResponse, error) {
	// For now, we'll need to get the user's portfolios
	// This would typically involve a database query to get user's portfolios

	return &VoiceResponse{
		Text: "Portfolio checking is not fully implemented yet. This would show your current portfolio status, including total value, P&L, and active positions.",
		Actions: []SuggestedAction{
			{Text: "Create Portfolio", Command: "create a portfolio", Description: "Create a new trading portfolio"},
			{Text: "Get Market Data", Command: "what's the price of ethereum", Description: "Check current market prices"},
		},
		Confidence: command.Confidence,
	}, nil
}

// handleHelp handles help commands
func (v *VoiceInterface) handleHelp(ctx context.Context, command VoiceCommand) (*VoiceResponse, error) {
	helpText := `I can help you with cryptocurrency trading and portfolio management. Here are some things you can say:

Portfolio Management:
• "Create a portfolio with $10,000"
• "Check my portfolio"
• "Rebalance my portfolio"

Trading:
• "Buy 1 ETH"
• "Sell 0.5 BTC"
• "Start trading"
• "Stop trading"

Market Data:
• "What's the price of Bitcoin?"
• "Show me yield opportunities"
• "Check risk for this transaction"

Strategies:
• "Set momentum strategy"
• "Use conservative risk level"

Just speak naturally and I'll understand what you want to do!`

	return &VoiceResponse{
		Text: helpText,
		Actions: []SuggestedAction{
			{Text: "Create Portfolio", Command: "create a portfolio", Description: "Start with a new portfolio"},
			{Text: "Check Market", Command: "what's the price of ethereum", Description: "Get current prices"},
			{Text: "Find Yield", Command: "show me yield opportunities", Description: "Discover DeFi yields"},
		},
		Confidence: 1.0,
	}, nil
}

// Placeholder handlers for other intents
func (v *VoiceInterface) handleBuyToken(ctx context.Context, command VoiceCommand) (*VoiceResponse, error) {
	return &VoiceResponse{
		Text:       "Token buying via voice commands is not yet implemented. This would execute a buy order based on your voice command.",
		Confidence: command.Confidence,
	}, nil
}

func (v *VoiceInterface) handleSellToken(ctx context.Context, command VoiceCommand) (*VoiceResponse, error) {
	return &VoiceResponse{
		Text:       "Token selling via voice commands is not yet implemented. This would execute a sell order based on your voice command.",
		Confidence: command.Confidence,
	}, nil
}

func (v *VoiceInterface) handleCheckBalance(ctx context.Context, command VoiceCommand) (*VoiceResponse, error) {
	return &VoiceResponse{
		Text:       "Balance checking is not yet implemented. This would show your current wallet balances across all chains.",
		Confidence: command.Confidence,
	}, nil
}

func (v *VoiceInterface) handleSetStrategy(ctx context.Context, command VoiceCommand) (*VoiceResponse, error) {
	return &VoiceResponse{
		Text:       "Strategy setting is not yet implemented. This would configure your trading strategies based on voice commands.",
		Confidence: command.Confidence,
	}, nil
}

func (v *VoiceInterface) handleStartTrading(ctx context.Context, command VoiceCommand) (*VoiceResponse, error) {
	return &VoiceResponse{
		Text:       "Trading is already running globally. Individual portfolio trading control is not yet implemented.",
		Confidence: command.Confidence,
	}, nil
}

func (v *VoiceInterface) handleStopTrading(ctx context.Context, command VoiceCommand) (*VoiceResponse, error) {
	return &VoiceResponse{
		Text:       "Trading stop functionality is not yet implemented. This would stop autonomous trading for your portfolios.",
		Confidence: command.Confidence,
	}, nil
}

func (v *VoiceInterface) handleGetMarketData(ctx context.Context, command VoiceCommand) (*VoiceResponse, error) {
	return &VoiceResponse{
		Text:       "Market data retrieval is not yet implemented. This would provide current prices and market information.",
		Confidence: command.Confidence,
	}, nil
}

func (v *VoiceInterface) handleRebalancePortfolio(ctx context.Context, command VoiceCommand) (*VoiceResponse, error) {
	return &VoiceResponse{
		Text:       "Portfolio rebalancing via voice is not yet implemented. This would trigger portfolio rebalancing based on your strategy.",
		Confidence: command.Confidence,
	}, nil
}

func (v *VoiceInterface) handleCheckRisk(ctx context.Context, command VoiceCommand) (*VoiceResponse, error) {
	return &VoiceResponse{
		Text:       "Risk checking via voice is not yet implemented. This would analyze transaction or portfolio risk.",
		Confidence: command.Confidence,
	}, nil
}

func (v *VoiceInterface) handleFindYield(ctx context.Context, command VoiceCommand) (*VoiceResponse, error) {
	return &VoiceResponse{
		Text:       "Yield opportunity discovery is not yet implemented. This would show the best DeFi yield opportunities.",
		Confidence: command.Confidence,
	}, nil
}

// addToHistory adds a command to the history
func (v *VoiceInterface) addToHistory(command VoiceCommand) {
	v.commandHistory = append(v.commandHistory, command)

	// Keep history within limits
	if len(v.commandHistory) > v.config.MaxCommandHistory {
		v.commandHistory = v.commandHistory[1:]
	}
}

// GetCommandHistory returns the command history
func (v *VoiceInterface) GetCommandHistory(userID uuid.UUID) []VoiceCommand {
	var userHistory []VoiceCommand
	for _, cmd := range v.commandHistory {
		if cmd.UserID == userID {
			userHistory = append(userHistory, cmd)
		}
	}
	return userHistory
}
