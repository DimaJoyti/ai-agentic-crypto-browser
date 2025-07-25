package ai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ai-agentic-browser/internal/config"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create test logger
func createLMStudioTestLogger() *observability.Logger {
	return observability.NewLogger(config.ObservabilityConfig{
		ServiceName: "test",
		LogLevel:    "info",
		LogFormat:   "json",
	})
}

func TestLMStudioProvider_GenerateResponse(t *testing.T) {
	// Mock LM Studio server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/chat/completions" {
			response := lmStudioChatResponse{
				ID:      "chatcmpl-123",
				Object:  "chat.completion",
				Created: time.Now().Unix(),
				Model:   "local-model",
				Choices: []lmStudioChoice{
					{
						Index: 0,
						Message: lmStudioMessage{
							Role:    "assistant",
							Content: "Hello! I'm here to help you browse the web and analyze content.",
						},
						FinishReason: "stop",
					},
				},
				Usage: lmStudioUsage{
					PromptTokens:     10,
					CompletionTokens: 15,
					TotalTokens:      25,
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()

	logger := createLMStudioTestLogger()

	config := LMStudioConfig{
		BaseURL:     server.URL,
		Model:       "local-model",
		Temperature: 0.7,
		MaxTokens:   1000,
		TopP:        1.0,
		Timeout:     30 * time.Second,
	}

	provider := NewLMStudioProvider(config, logger)

	messages := []Message{
		{
			Role:    RoleUser,
			Content: "Hello, can you help me browse a website?",
		},
	}

	ctx := context.Background()
	response, err := provider.GenerateResponse(ctx, messages)

	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, RoleAssistant, response.Role)
	assert.Contains(t, response.Content, "help")
	assert.Equal(t, "lmstudio", response.Metadata["provider"])
	assert.Equal(t, 10, response.Metadata["prompt_tokens"])
	assert.Equal(t, 15, response.Metadata["completion_tokens"])
}

func TestLMStudioProvider_AnalyzeContent(t *testing.T) {
	// Mock LM Studio server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/chat/completions" {
			response := lmStudioChatResponse{
				ID:      "chatcmpl-123",
				Object:  "chat.completion",
				Created: time.Now().Unix(),
				Model:   "local-model",
				Choices: []lmStudioChoice{
					{
						Index: 0,
						Message: lmStudioMessage{
							Role:    "assistant",
							Content: `{"sentiment": "positive", "confidence": 0.85, "key_phrases": ["helpful", "analysis"]}`,
						},
						FinishReason: "stop",
					},
				},
				Usage: lmStudioUsage{
					PromptTokens:     20,
					CompletionTokens: 30,
					TotalTokens:      50,
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()

	logger := createLMStudioTestLogger()

	config := LMStudioConfig{
		BaseURL: server.URL,
		Model:   "local-model",
		Timeout: 30 * time.Second,
	}

	provider := NewLMStudioProvider(config, logger)

	ctx := context.Background()
	analysis, err := provider.AnalyzeContent(ctx, "This is a helpful analysis tool", "sentiment")

	require.NoError(t, err)
	assert.NotNil(t, analysis)
	assert.Equal(t, "positive", analysis["sentiment"])
	assert.Equal(t, 0.85, analysis["confidence"])
}

func TestLMStudioProvider_ExtractStructuredData(t *testing.T) {
	// Mock LM Studio server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/chat/completions" {
			response := lmStudioChatResponse{
				ID:      "chatcmpl-123",
				Object:  "chat.completion",
				Created: time.Now().Unix(),
				Model:   "local-model",
				Choices: []lmStudioChoice{
					{
						Index: 0,
						Message: lmStudioMessage{
							Role:    "assistant",
							Content: `{"name": "John Doe", "email": "john@example.com", "age": 30}`,
						},
						FinishReason: "stop",
					},
				},
				Usage: lmStudioUsage{
					PromptTokens:     25,
					CompletionTokens: 20,
					TotalTokens:      45,
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()

	logger := createLMStudioTestLogger()

	config := LMStudioConfig{
		BaseURL: server.URL,
		Model:   "local-model",
		Timeout: 30 * time.Second,
	}

	provider := NewLMStudioProvider(config, logger)

	schema := `{
		"type": "object",
		"properties": {
			"name": {"type": "string"},
			"email": {"type": "string"},
			"age": {"type": "number"}
		}
	}`

	content := "Contact: John Doe, email: john@example.com, age: 30"

	ctx := context.Background()
	data, err := provider.ExtractStructuredData(ctx, content, schema)

	require.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, "John Doe", data["name"])
	assert.Equal(t, "john@example.com", data["email"])
	assert.Equal(t, float64(30), data["age"])
}

func TestLMStudioProvider_SummarizeContent(t *testing.T) {
	// Mock LM Studio server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/chat/completions" {
			response := lmStudioChatResponse{
				ID:      "chatcmpl-123",
				Object:  "chat.completion",
				Created: time.Now().Unix(),
				Model:   "local-model",
				Choices: []lmStudioChoice{
					{
						Index: 0,
						Message: lmStudioMessage{
							Role:    "assistant",
							Content: "This is a concise summary of the provided content.",
						},
						FinishReason: "stop",
					},
				},
				Usage: lmStudioUsage{
					PromptTokens:     50,
					CompletionTokens: 15,
					TotalTokens:      65,
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()

	logger := createLMStudioTestLogger()

	config := LMStudioConfig{
		BaseURL: server.URL,
		Model:   "local-model",
		Timeout: 30 * time.Second,
	}

	provider := NewLMStudioProvider(config, logger)

	content := "This is a long piece of content that needs to be summarized. It contains multiple sentences and paragraphs with various information that should be condensed into a shorter form."

	options := SummarizeOptions{
		Length: "short",
		Style:  "professional",
	}

	ctx := context.Background()
	summary, err := provider.SummarizeContent(ctx, content, options)

	require.NoError(t, err)
	assert.NotEmpty(t, summary)
	assert.Contains(t, summary, "summary")
}

func TestLMStudioProvider_IsHealthy(t *testing.T) {
	// Mock LM Studio server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/models" {
			response := lmStudioModelsResponse{
				Object: "list",
				Data: []lmStudioModelInfo{
					{
						ID:      "local-model",
						Object:  "model",
						Created: time.Now().Unix(),
						OwnedBy: "lmstudio",
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()

	logger := createLMStudioTestLogger()

	config := LMStudioConfig{
		BaseURL: server.URL,
		Model:   "local-model",
		Timeout: 30 * time.Second,
	}

	provider := NewLMStudioProvider(config, logger)

	ctx := context.Background()
	err := provider.IsHealthy(ctx)

	assert.NoError(t, err)
}

func TestLMStudioProvider_ErrorHandling(t *testing.T) {
	// Mock server that returns errors
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
	}))
	defer server.Close()

	logger := createLMStudioTestLogger()

	config := LMStudioConfig{
		BaseURL: server.URL,
		Model:   "local-model",
		Timeout: 30 * time.Second,
	}

	provider := NewLMStudioProvider(config, logger)

	ctx := context.Background()

	// Test GenerateResponse error handling
	messages := []Message{{Role: RoleUser, Content: "test"}}
	_, err := provider.GenerateResponse(ctx, messages)
	assert.Error(t, err)
	assert.Contains(t, strings.ToLower(err.Error()), "failed")

	// Test health check error handling
	err = provider.IsHealthy(ctx)
	assert.Error(t, err)
}
