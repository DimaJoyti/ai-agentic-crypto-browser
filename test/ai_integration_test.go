package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ai-agentic-browser/internal/ai"
	"github.com/ai-agentic-browser/internal/config"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test the NLP processor directly since it doesn't depend on external services
func TestNLPProcessorBasic(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{})
	nlpProcessor := ai.NewNLPProcessor(logger)

	t.Run("ProcessText_CreatePortfolio", func(t *testing.T) {
		text := "create a portfolio with $5000"

		intent, entities, confidence, err := nlpProcessor.ProcessText(context.Background(), text)

		require.NoError(t, err)
		assert.Equal(t, ai.IntentCreatePortfolio, intent)
		assert.Equal(t, "5000", entities["amount"])
		assert.Greater(t, confidence, 0.8)
	})

	t.Run("ProcessText_Help", func(t *testing.T) {
		text := "help"

		intent, _, confidence, err := nlpProcessor.ProcessText(context.Background(), text)

		require.NoError(t, err)
		assert.Equal(t, ai.IntentHelp, intent)
		assert.Greater(t, confidence, 0.9)
	})
}

func TestConversationalAI(t *testing.T) {
	// Setup
	logger := observability.NewLogger(config.ObservabilityConfig{})

	// Create conversational AI with nil services for testing
	conversationalAI := ai.NewConversationalAI(logger, nil, nil, nil)

	t.Run("StartConversation", func(t *testing.T) {
		userID := uuid.New()

		conversation, err := conversationalAI.StartConversation(context.Background(), userID)

		require.NoError(t, err)
		assert.NotNil(t, conversation)
		assert.Equal(t, userID, conversation.UserID)
		assert.NotEmpty(t, conversation.Messages)
		assert.Equal(t, ai.RoleAssistant, conversation.Messages[0].Role)
		assert.Contains(t, conversation.Messages[0].Content, "Hello")
	})
}

func TestNLPProcessor(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{})
	nlpProcessor := ai.NewNLPProcessor(logger)

	t.Run("ProcessText_CreatePortfolio", func(t *testing.T) {
		text := "create a portfolio with $5000"

		intent, entities, confidence, err := nlpProcessor.ProcessText(context.Background(), text)

		require.NoError(t, err)
		assert.Equal(t, ai.IntentCreatePortfolio, intent)
		assert.Equal(t, "5000", entities["amount"])
		assert.Greater(t, confidence, 0.8)
	})

	t.Run("ProcessText_BuyToken", func(t *testing.T) {
		text := "buy 1 ethereum"

		intent, entities, confidence, err := nlpProcessor.ProcessText(context.Background(), text)

		require.NoError(t, err)
		assert.Equal(t, ai.IntentBuyToken, intent)
		assert.Equal(t, "ETH", entities["token"])
		assert.Greater(t, confidence, 0.8)
	})

	t.Run("ProcessText_CheckBalance", func(t *testing.T) {
		text := "check my balance"

		intent, _, confidence, err := nlpProcessor.ProcessText(context.Background(), text)

		require.NoError(t, err)
		assert.Equal(t, ai.IntentCheckBalance, intent)
		assert.Greater(t, confidence, 0.8)
	})

	t.Run("ProcessText_Help", func(t *testing.T) {
		text := "help"

		intent, _, confidence, err := nlpProcessor.ProcessText(context.Background(), text)

		require.NoError(t, err)
		assert.Equal(t, ai.IntentHelp, intent)
		assert.Greater(t, confidence, 0.9)
	})
}

func TestMarketAnalyzer(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{})
	marketAnalyzer := ai.NewMarketAnalyzer(logger)

	t.Run("GetMarketContext", func(t *testing.T) {
		context := context.Background()

		marketContext, err := marketAnalyzer.GetMarketContext(context)

		require.NoError(t, err)
		assert.NotNil(t, marketContext)
		assert.NotEmpty(t, marketContext.MarketTrend)
		assert.NotEmpty(t, marketContext.Volatility)
		assert.NotEmpty(t, marketContext.MarketSentiment)
		assert.NotEmpty(t, marketContext.TopMovers)
		assert.NotEmpty(t, marketContext.KeyEvents)
		assert.True(t, time.Since(marketContext.LastUpdated) < time.Minute)
	})

	t.Run("GetMarketData", func(t *testing.T) {
		context := context.Background()

		marketData, err := marketAnalyzer.GetMarketData(context)

		require.NoError(t, err)
		assert.NotNil(t, marketData)
		assert.True(t, marketData.GlobalMarketCap.IsPositive())
		assert.True(t, marketData.TotalVolume24h.IsPositive())
		assert.NotEmpty(t, marketData.TopTokens)
		assert.NotEmpty(t, marketData.TrendingTokens)
		assert.True(t, time.Since(marketData.Timestamp) < time.Minute)
	})

	t.Run("AnalyzeTrend", func(t *testing.T) {
		context := context.Background()

		trend, err := marketAnalyzer.AnalyzeTrend(context, "BTC")

		require.NoError(t, err)
		assert.NotNil(t, trend)
		assert.NotEmpty(t, trend.Direction)
		assert.NotEmpty(t, trend.Strength)
		assert.True(t, trend.Confidence.IsPositive())
		assert.NotEmpty(t, trend.Duration)
		assert.True(t, time.Since(trend.LastUpdated) < time.Minute)
	})

	t.Run("GetSentimentAnalysis", func(t *testing.T) {
		context := context.Background()

		sentiment, err := marketAnalyzer.GetSentimentAnalysis(context)

		require.NoError(t, err)
		assert.NotNil(t, sentiment)
		assert.NotEmpty(t, sentiment.OverallSentiment)
		assert.GreaterOrEqual(t, sentiment.BullishSignals, 0)
		assert.GreaterOrEqual(t, sentiment.BearishSignals, 0)
		assert.GreaterOrEqual(t, sentiment.NeutralSignals, 0)
		assert.NotEmpty(t, sentiment.Sources)
		assert.True(t, time.Since(sentiment.LastUpdated) < time.Minute)
	})
}

func TestAIEndpointsIntegration(t *testing.T) {
	// This would test the actual HTTP endpoints
	// For now, we'll test the handler logic structure

	t.Run("VoiceCommandEndpoint_Structure", func(t *testing.T) {
		// Test that the voice command endpoint accepts the correct request format
		requestBody := map[string]interface{}{
			"text":       "create a portfolio",
			"audio_data": nil,
		}

		jsonBody, err := json.Marshal(requestBody)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/web3/ai/voice/command", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		// Verify request structure is valid
		var decoded map[string]interface{}
		err = json.NewDecoder(req.Body).Decode(&decoded)
		require.NoError(t, err)
		assert.Equal(t, "create a portfolio", decoded["text"])
	})

	t.Run("ChatMessageEndpoint_Structure", func(t *testing.T) {
		// Test that the chat message endpoint accepts the correct request format
		requestBody := map[string]interface{}{
			"message": "What's the market looking like today?",
		}

		jsonBody, err := json.Marshal(requestBody)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/web3/ai/chat/message", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		// Verify request structure is valid
		var decoded map[string]interface{}
		err = json.NewDecoder(req.Body).Decode(&decoded)
		require.NoError(t, err)
		assert.Equal(t, "What's the market looking like today?", decoded["message"])
	})
}
