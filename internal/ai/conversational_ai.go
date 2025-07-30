package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ai-agentic-browser/internal/web3"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ConversationalAI provides intelligent market analysis and investment insights
type ConversationalAI struct {
	logger         *observability.Logger
	tradingEngine  *web3.TradingEngine
	defiManager    *web3.DeFiProtocolManager
	riskAssessment *web3.RiskAssessmentService
	marketAnalyzer *MarketAnalyzer
	conversations  map[uuid.UUID]*Conversation
	config         ConversationalConfig
}

// ConversationalConfig holds configuration for conversational AI
type ConversationalConfig struct {
	MaxConversationHistory int           `json:"max_conversation_history"`
	ContextWindow          int           `json:"context_window"`
	ResponseTimeout        time.Duration `json:"response_timeout"`
	EnablePersonalization  bool          `json:"enable_personalization"`
	EnableMarketInsights   bool          `json:"enable_market_insights"`
	EnableRiskWarnings     bool          `json:"enable_risk_warnings"`
}

// Conversation represents an ongoing conversation with a user
type Conversation struct {
	ID         uuid.UUID              `json:"id"`
	UserID     uuid.UUID              `json:"user_id"`
	Messages   []ConversationMessage  `json:"messages"`
	Context    ConversationContext    `json:"context"`
	StartedAt  time.Time              `json:"started_at"`
	LastActive time.Time              `json:"last_active"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// ConversationMessage represents a message in a conversation
type ConversationMessage struct {
	ID        uuid.UUID       `json:"id"`
	Role      MessageRole     `json:"role"`
	Content   string          `json:"content"`
	Timestamp time.Time       `json:"timestamp"`
	Metadata  json.RawMessage `json:"metadata,omitempty"`
}

// MessageRole represents the role of a message sender
type MessageRole string

const (
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
	RoleSystem    MessageRole = "system"
)

// ConversationContext holds context for the conversation
type ConversationContext struct {
	UserPreferences  UserPreferences        `json:"user_preferences"`
	CurrentPortfolio *web3.Portfolio        `json:"current_portfolio,omitempty"`
	MarketContext    MarketContext          `json:"market_context"`
	RecentActions    []string               `json:"recent_actions"`
	Topics           []string               `json:"topics"`
	Metadata         map[string]interface{} `json:"metadata"`
}

// UserPreferences represents user preferences for AI interactions
type UserPreferences struct {
	RiskTolerance     string   `json:"risk_tolerance"`
	InvestmentGoals   []string `json:"investment_goals"`
	PreferredTokens   []string `json:"preferred_tokens"`
	TradingStyle      string   `json:"trading_style"`
	NotificationLevel string   `json:"notification_level"`
}

// MarketContext represents current market conditions
type MarketContext struct {
	MarketTrend     string          `json:"market_trend"`
	Volatility      string          `json:"volatility"`
	TopMovers       []TokenMovement `json:"top_movers"`
	MarketSentiment string          `json:"market_sentiment"`
	KeyEvents       []MarketEvent   `json:"key_events"`
	LastUpdated     time.Time       `json:"last_updated"`
}

// TokenMovement represents price movement data
type TokenMovement struct {
	Symbol     string          `json:"symbol"`
	Price      decimal.Decimal `json:"price"`
	Change24h  decimal.Decimal `json:"change_24h"`
	ChangePerc decimal.Decimal `json:"change_perc"`
	Volume24h  decimal.Decimal `json:"volume_24h"`
}

// MarketEvent represents a significant market event
type MarketEvent struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Impact      string    `json:"impact"`
	Timestamp   time.Time `json:"timestamp"`
	Source      string    `json:"source"`
}

// ConversationalResponse represents an AI response
type ConversationalResponse struct {
	Content     string                 `json:"content"`
	Insights    []MarketInsight        `json:"insights,omitempty"`
	Suggestions []ActionSuggestion     `json:"suggestions,omitempty"`
	Warnings    []RiskWarning          `json:"warnings,omitempty"`
	Data        interface{}            `json:"data,omitempty"`
	Confidence  float64                `json:"confidence"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// MarketInsight represents an AI-generated market insight
type MarketInsight struct {
	Type        string          `json:"type"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Confidence  float64         `json:"confidence"`
	Impact      string          `json:"impact"`
	Timeframe   string          `json:"timeframe"`
	Data        json.RawMessage `json:"data,omitempty"`
}

// ActionSuggestion represents a suggested action
type ActionSuggestion struct {
	Action      string          `json:"action"`
	Description string          `json:"description"`
	Reasoning   string          `json:"reasoning"`
	Risk        string          `json:"risk"`
	Potential   decimal.Decimal `json:"potential,omitempty"`
	Command     string          `json:"command,omitempty"`
}

// RiskWarning represents a risk warning
type RiskWarning struct {
	Level       string `json:"level"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Mitigation  string `json:"mitigation"`
}

// NewConversationalAI creates a new conversational AI service
func NewConversationalAI(
	logger *observability.Logger,
	tradingEngine *web3.TradingEngine,
	defiManager *web3.DeFiProtocolManager,
	riskAssessment *web3.RiskAssessmentService,
) *ConversationalAI {
	config := ConversationalConfig{
		MaxConversationHistory: 50,
		ContextWindow:          10,
		ResponseTimeout:        30 * time.Second,
		EnablePersonalization:  true,
		EnableMarketInsights:   true,
		EnableRiskWarnings:     true,
	}

	return &ConversationalAI{
		logger:         logger,
		tradingEngine:  tradingEngine,
		defiManager:    defiManager,
		riskAssessment: riskAssessment,
		marketAnalyzer: NewMarketAnalyzer(logger),
		conversations:  make(map[uuid.UUID]*Conversation),
		config:         config,
	}
}

// StartConversation starts a new conversation with a user
func (c *ConversationalAI) StartConversation(ctx context.Context, userID uuid.UUID) (*Conversation, error) {
	conversation := &Conversation{
		ID:         uuid.New(),
		UserID:     userID,
		Messages:   make([]ConversationMessage, 0),
		Context:    c.initializeContext(ctx, userID),
		StartedAt:  time.Now(),
		LastActive: time.Now(),
		Metadata:   make(map[string]interface{}),
	}

	c.conversations[userID] = conversation

	// Add welcome message
	welcomeMsg := c.generateWelcomeMessage(ctx, conversation)
	c.addMessage(conversation, RoleAssistant, welcomeMsg)

	c.logger.Info(ctx, "Conversation started", map[string]interface{}{
		"conversation_id": conversation.ID.String(),
		"user_id":         userID.String(),
	})

	return conversation, nil
}

// ProcessMessage processes a user message and generates a response
func (c *ConversationalAI) ProcessMessage(ctx context.Context, userID uuid.UUID, message string) (*ConversationalResponse, error) {
	conversation, exists := c.conversations[userID]
	if !exists {
		var err error
		conversation, err = c.StartConversation(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to start conversation: %w", err)
		}
	}

	// Add user message
	c.addMessage(conversation, RoleUser, message)
	conversation.LastActive = time.Now()

	// Update context based on message
	c.updateContext(ctx, conversation, message)

	// Generate response
	response, err := c.generateResponse(ctx, conversation, message)
	if err != nil {
		return nil, fmt.Errorf("failed to generate response: %w", err)
	}

	// Add assistant response
	c.addMessage(conversation, RoleAssistant, response.Content)

	return response, nil
}

// generateResponse generates an AI response based on the conversation context
func (c *ConversationalAI) generateResponse(ctx context.Context, conversation *Conversation, message string) (*ConversationalResponse, error) {
	// Analyze the message intent and context
	intent := c.analyzeIntent(message)

	// Get market context
	marketContext, err := c.marketAnalyzer.GetMarketContext(ctx)
	if err != nil {
		c.logger.Warn(ctx, "Failed to get market context", map[string]interface{}{"error": err.Error()})
	}

	response := &ConversationalResponse{
		Insights:    make([]MarketInsight, 0),
		Suggestions: make([]ActionSuggestion, 0),
		Warnings:    make([]RiskWarning, 0),
		Confidence:  0.8,
		Metadata:    make(map[string]interface{}),
	}

	// Generate response based on intent
	switch intent {
	case "market_analysis":
		response.Content = c.generateMarketAnalysis(ctx, conversation, marketContext)
		response.Insights = c.generateMarketInsights(ctx, marketContext)
	case "portfolio_advice":
		response.Content = c.generatePortfolioAdvice(ctx, conversation)
		response.Suggestions = c.generatePortfolioSuggestions(ctx, conversation)
	case "risk_assessment":
		response.Content = c.generateRiskAssessment(ctx, conversation, message)
		response.Warnings = c.generateRiskWarnings(ctx, conversation)
	case "yield_opportunities":
		response.Content = c.generateYieldAnalysis(ctx, conversation)
		response.Suggestions = c.generateYieldSuggestions(ctx)
	case "general_question":
		response.Content = c.generateGeneralResponse(ctx, conversation, message)
	default:
		response.Content = c.generateDefaultResponse(ctx, conversation, message)
	}

	return response, nil
}

// generateMarketAnalysis generates market analysis content
func (c *ConversationalAI) generateMarketAnalysis(ctx context.Context, conversation *Conversation, marketContext *MarketContext) string {
	if marketContext == nil {
		return "I'm currently unable to access real-time market data, but I can help you with portfolio management and trading strategies. What specific aspect of the market are you interested in?"
	}

	analysis := fmt.Sprintf(`Based on current market conditions:

📈 **Market Trend**: %s
📊 **Volatility**: %s  
💭 **Sentiment**: %s

**Key Observations:**
• The market is showing %s characteristics
• Volatility levels are %s, which suggests %s
• Current sentiment indicates %s

**What this means for your portfolio:**
• Consider %s strategies in this environment
• Risk management is %s important right now
• Opportunities may exist in %s sectors

Would you like me to analyze specific tokens or discuss portfolio adjustments?`,
		marketContext.MarketTrend,
		marketContext.Volatility,
		marketContext.MarketSentiment,
		strings.ToLower(marketContext.MarketTrend),
		strings.ToLower(marketContext.Volatility),
		c.getVolatilityImplication(marketContext.Volatility),
		strings.ToLower(marketContext.MarketSentiment),
		c.getRecommendedStrategy(marketContext.MarketTrend),
		c.getRiskImportance(marketContext.Volatility),
		c.getOpportunitySectors(marketContext.MarketTrend))

	return analysis
}

// generatePortfolioAdvice generates portfolio advice
func (c *ConversationalAI) generatePortfolioAdvice(ctx context.Context, conversation *Conversation) string {
	if conversation.Context.CurrentPortfolio == nil {
		return `I'd love to help you with portfolio advice! However, I don't see an active portfolio in our conversation. 

Here's what I can help you with:

🎯 **Portfolio Creation**: I can guide you through creating a diversified portfolio
📊 **Risk Assessment**: Analyze your risk tolerance and suggest appropriate allocations  
⚖️ **Rebalancing**: Help optimize your current holdings
📈 **Strategy Selection**: Choose trading strategies that match your goals

Would you like to start by creating a portfolio or discussing your investment goals?`
	}

	portfolio := conversation.Context.CurrentPortfolio
	return fmt.Sprintf(`Let me analyze your current portfolio:

💼 **Portfolio Overview**:
• Total Value: $%s
• Available Balance: $%s  
• Total P&L: $%s (%.2f%%)
• Active Positions: %d

**Performance Analysis**:
%s

**Recommendations**:
%s

Would you like me to suggest any specific adjustments or analyze particular positions?`,
		portfolio.TotalValue.String(),
		portfolio.AvailableBalance.String(),
		portfolio.TotalPnL.String(),
		portfolio.TotalPnL.Div(portfolio.InvestedAmount).Mul(decimal.NewFromInt(100)).InexactFloat64(),
		len(portfolio.ActivePositions),
		c.analyzePortfolioPerformance(portfolio),
		c.generatePortfolioRecommendations(portfolio))
}

// Helper methods for generating contextual responses
func (c *ConversationalAI) getVolatilityImplication(volatility string) string {
	switch strings.ToLower(volatility) {
	case "high":
		return "both higher risks and potential opportunities"
	case "low":
		return "more stable conditions with limited price swings"
	default:
		return "moderate price movements"
	}
}

func (c *ConversationalAI) getRecommendedStrategy(trend string) string {
	switch strings.ToLower(trend) {
	case "bullish":
		return "momentum-based"
	case "bearish":
		return "defensive or mean-reversion"
	default:
		return "balanced"
	}
}

func (c *ConversationalAI) getRiskImportance(volatility string) string {
	switch strings.ToLower(volatility) {
	case "high":
		return "especially"
	case "low":
		return "still"
	default:
		return "particularly"
	}
}

func (c *ConversationalAI) getOpportunitySectors(trend string) string {
	switch strings.ToLower(trend) {
	case "bullish":
		return "growth and momentum"
	case "bearish":
		return "defensive and value"
	default:
		return "balanced"
	}
}

// Placeholder methods for complex analysis
func (c *ConversationalAI) analyzeIntent(message string) string {
	message = strings.ToLower(message)

	if strings.Contains(message, "market") || strings.Contains(message, "price") || strings.Contains(message, "trend") {
		return "market_analysis"
	}
	if strings.Contains(message, "portfolio") || strings.Contains(message, "holdings") {
		return "portfolio_advice"
	}
	if strings.Contains(message, "risk") || strings.Contains(message, "safe") {
		return "risk_assessment"
	}
	if strings.Contains(message, "yield") || strings.Contains(message, "apy") || strings.Contains(message, "defi") {
		return "yield_opportunities"
	}

	return "general_question"
}

func (c *ConversationalAI) initializeContext(ctx context.Context, userID uuid.UUID) ConversationContext {
	return ConversationContext{
		UserPreferences: UserPreferences{
			RiskTolerance:     "moderate",
			InvestmentGoals:   []string{"growth"},
			PreferredTokens:   []string{"ETH", "BTC"},
			TradingStyle:      "balanced",
			NotificationLevel: "normal",
		},
		MarketContext: MarketContext{LastUpdated: time.Now()},
		RecentActions: make([]string, 0),
		Topics:        make([]string, 0),
		Metadata:      make(map[string]interface{}),
	}
}

func (c *ConversationalAI) generateWelcomeMessage(ctx context.Context, conversation *Conversation) string {
	return `👋 Hello! I'm your AI crypto assistant. I can help you with:

🎯 **Portfolio Management** - Create, analyze, and optimize your portfolios
📊 **Market Analysis** - Get insights on market trends and token performance  
⚖️ **Risk Assessment** - Evaluate risks for transactions and strategies
🏦 **DeFi Opportunities** - Find the best yield farming and staking options
🤖 **Autonomous Trading** - Set up and manage automated trading strategies

What would you like to explore today?`
}

func (c *ConversationalAI) addMessage(conversation *Conversation, role MessageRole, content string) {
	message := ConversationMessage{
		ID:        uuid.New(),
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
	}

	conversation.Messages = append(conversation.Messages, message)

	// Keep conversation history within limits
	if len(conversation.Messages) > c.config.MaxConversationHistory {
		conversation.Messages = conversation.Messages[1:]
	}
}

func (c *ConversationalAI) updateContext(ctx context.Context, conversation *Conversation, message string) {
	// Update topics and recent actions based on message content
	// This would be more sophisticated in a real implementation
	conversation.Context.RecentActions = append(conversation.Context.RecentActions, message)
	if len(conversation.Context.RecentActions) > 5 {
		conversation.Context.RecentActions = conversation.Context.RecentActions[1:]
	}
}

// Placeholder methods for complex generation
func (c *ConversationalAI) generateMarketInsights(ctx context.Context, marketContext *MarketContext) []MarketInsight {
	return []MarketInsight{}
}

func (c *ConversationalAI) generatePortfolioSuggestions(ctx context.Context, conversation *Conversation) []ActionSuggestion {
	return []ActionSuggestion{}
}

func (c *ConversationalAI) generateRiskAssessment(ctx context.Context, conversation *Conversation, message string) string {
	return "Risk assessment functionality is not yet fully implemented. I can help you understand general risk principles and portfolio safety measures."
}

func (c *ConversationalAI) generateRiskWarnings(ctx context.Context, conversation *Conversation) []RiskWarning {
	return []RiskWarning{}
}

func (c *ConversationalAI) generateYieldAnalysis(ctx context.Context, conversation *Conversation) string {
	return "I can help you find yield opportunities across DeFi protocols. The system supports Uniswap V3, Compound, and Aave with real-time APY tracking."
}

func (c *ConversationalAI) generateYieldSuggestions(ctx context.Context) []ActionSuggestion {
	return []ActionSuggestion{}
}

func (c *ConversationalAI) generateGeneralResponse(ctx context.Context, conversation *Conversation, message string) string {
	return "I understand you have a question about cryptocurrency and trading. Could you be more specific about what you'd like to know? I can help with portfolio management, market analysis, risk assessment, and DeFi opportunities."
}

func (c *ConversationalAI) generateDefaultResponse(ctx context.Context, conversation *Conversation, message string) string {
	return "I'm here to help with your cryptocurrency needs. You can ask me about market trends, portfolio management, trading strategies, or DeFi opportunities. What would you like to explore?"
}

func (c *ConversationalAI) analyzePortfolioPerformance(portfolio *web3.Portfolio) string {
	if portfolio.TotalPnL.IsPositive() {
		return "Your portfolio is performing well with positive returns."
	} else if portfolio.TotalPnL.IsNegative() {
		return "Your portfolio is currently showing losses. Consider reviewing your strategy."
	}
	return "Your portfolio is at break-even."
}

func (c *ConversationalAI) generatePortfolioRecommendations(portfolio *web3.Portfolio) string {
	return "Consider diversifying across different asset classes and maintaining appropriate risk management."
}
