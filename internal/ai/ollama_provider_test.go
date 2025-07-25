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
func createOllamaTestLogger() *observability.Logger {
	return observability.NewLogger(config.ObservabilityConfig{
		ServiceName: "test",
		LogLevel:    "info",
		LogFormat:   "json",
	})
}

func TestOllamaProvider_GenerateResponse(t *testing.T) {
	// Mock Ollama server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/chat" {
			response := ollamaChatResponse{
				Model:     "qwen3",
				CreatedAt: time.Now(),
				Message: ollamaChatMessage{
					Role:    "assistant",
					Content: "Hello! I'm here to help you browse the web and analyze content.",
				},
				Done: true,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()

	logger := observability.NewLogger(config.ObservabilityConfig{
		ServiceName: "test",
		LogLevel:    "info",
		LogFormat:   "json",
	})

	config := OllamaConfig{
		BaseURL:     server.URL,
		Model:       "qwen3",
		Temperature: 0.7,
		TopP:        1.0,
		TopK:        40,
		NumCtx:      2048,
		Timeout:     30 * time.Second,
	}

	provider := NewOllamaProvider(config, logger)

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
	assert.Equal(t, "ollama", response.Metadata["provider"])
}

func TestOllamaProvider_AnalyzeContent(t *testing.T) {
	// Mock Ollama server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/generate" {
			response := ollamaGenerateResponse{
				Model:     "qwen3",
				CreatedAt: time.Now(),
				Response:  `{"sentiment": "positive", "confidence": 0.85, "key_phrases": ["helpful", "analysis"]}`,
				Done:      true,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()

	logger := createOllamaTestLogger()

	config := OllamaConfig{
		BaseURL: server.URL,
		Model:   "qwen3",
		Timeout: 30 * time.Second,
	}

	provider := NewOllamaProvider(config, logger)

	ctx := context.Background()
	analysis, err := provider.AnalyzeContent(ctx, "This is a helpful analysis tool", "sentiment")

	require.NoError(t, err)
	assert.NotNil(t, analysis)
	assert.Equal(t, "positive", analysis["sentiment"])
	assert.Equal(t, 0.85, analysis["confidence"])
}

func TestOllamaProvider_ExtractStructuredData(t *testing.T) {
	// Mock Ollama server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/generate" {
			response := ollamaGenerateResponse{
				Model:     "qwen3",
				CreatedAt: time.Now(),
				Response:  `{"name": "John Doe", "email": "john@example.com", "age": 30}`,
				Done:      true,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()

	logger := createOllamaTestLogger()

	config := OllamaConfig{
		BaseURL: server.URL,
		Model:   "qwen3",
		Timeout: 30 * time.Second,
	}

	provider := NewOllamaProvider(config, logger)

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

func TestOllamaProvider_SummarizeContent(t *testing.T) {
	// Mock Ollama server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/generate" {
			response := ollamaGenerateResponse{
				Model:     "qwen3",
				CreatedAt: time.Now(),
				Response:  "This is a concise summary of the provided content.",
				Done:      true,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()

	logger := createOllamaTestLogger()

	config := OllamaConfig{
		BaseURL: server.URL,
		Model:   "qwen3",
		Timeout: 30 * time.Second,
	}

	provider := NewOllamaProvider(config, logger)

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

func TestOllamaProvider_IsHealthy(t *testing.T) {
	// Mock Ollama server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/tags" {
			response := ollamaModelsResponse{
				Models: []ollamaModelInfo{
					{
						Name:       "qwen3",
						ModifiedAt: time.Now(),
						Size:       1000000,
						Digest:     "sha256:abc123",
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()

	logger := createOllamaTestLogger()

	config := OllamaConfig{
		BaseURL: server.URL,
		Model:   "qwen3",
		Timeout: 30 * time.Second,
	}

	provider := NewOllamaProvider(config, logger)

	ctx := context.Background()
	err := provider.IsHealthy(ctx)

	assert.NoError(t, err)
}

func TestOllamaProvider_ListModels(t *testing.T) {
	// Mock Ollama server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/tags" {
			response := ollamaModelsResponse{
				Models: []ollamaModelInfo{
					{Name: "qwen3"},
					{Name: "codellama"},
					{Name: "mistral"},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()

	logger := createOllamaTestLogger()

	config := OllamaConfig{
		BaseURL: server.URL,
		Model:   "qwen3",
		Timeout: 30 * time.Second,
	}

	provider := NewOllamaProvider(config, logger)

	ctx := context.Background()
	models, err := provider.ListModels(ctx)

	require.NoError(t, err)
	assert.Len(t, models, 3)
	assert.Contains(t, models, "qwen3")
	assert.Contains(t, models, "codellama")
	assert.Contains(t, models, "mistral")
}

func TestOllamaProvider_ErrorHandling(t *testing.T) {
	// Mock server that returns errors
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
	}))
	defer server.Close()

	logger := createOllamaTestLogger()

	config := OllamaConfig{
		BaseURL: server.URL,
		Model:   "qwen3",
		Timeout: 30 * time.Second,
	}

	provider := NewOllamaProvider(config, logger)

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
