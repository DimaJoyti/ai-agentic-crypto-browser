package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
)

// OllamaProvider implements AIProvider for Ollama
type OllamaProvider struct {
	baseURL    string
	model      string
	httpClient *http.Client
	logger     *observability.Logger
	config     OllamaConfig
}

// OllamaConfig holds Ollama-specific configuration
type OllamaConfig struct {
	BaseURL     string
	Model       string
	Temperature float64
	TopP        float64
	TopK        int
	NumCtx      int
	Timeout     time.Duration
	MaxRetries  int
	RetryDelay  time.Duration
}

// NewOllamaProvider creates a new Ollama provider
func NewOllamaProvider(config OllamaConfig, logger *observability.Logger) *OllamaProvider {
	if config.BaseURL == "" {
		config.BaseURL = "http://localhost:11434"
	}
	if config.Model == "" {
		config.Model = "qwen3"
	}
	if config.Temperature == 0 {
		config.Temperature = 0.7
	}
	if config.Timeout == 0 {
		config.Timeout = 60 * time.Second
	}

	return &OllamaProvider{
		baseURL: config.BaseURL,
		model:   config.Model,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		logger: logger,
		config: config,
	}
}

// Ollama API request/response structures
type ollamaChatRequest struct {
	Model    string                 `json:"model"`
	Messages []ollamaChatMessage    `json:"messages"`
	Stream   bool                   `json:"stream"`
	Options  map[string]interface{} `json:"options,omitempty"`
	Tools    []ollamaToolDefinition `json:"tools,omitempty"`
}

type ollamaChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaChatResponse struct {
	Model     string            `json:"model"`
	CreatedAt time.Time         `json:"created_at"`
	Message   ollamaChatMessage `json:"message"`
	Done      bool              `json:"done"`
	Error     string            `json:"error,omitempty"`
}

type ollamaToolDefinition struct {
	Type     string            `json:"type"`
	Function ollamaFunctionDef `json:"function"`
}

type ollamaFunctionDef struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type ollamaGenerateRequest struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	Stream  bool                   `json:"stream"`
	Options map[string]interface{} `json:"options,omitempty"`
}

type ollamaGenerateResponse struct {
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	Response  string    `json:"response"`
	Done      bool      `json:"done"`
	Error     string    `json:"error,omitempty"`
}

type ollamaModelInfo struct {
	Name       string    `json:"name"`
	ModifiedAt time.Time `json:"modified_at"`
	Size       int64     `json:"size"`
	Digest     string    `json:"digest"`
}

type ollamaModelsResponse struct {
	Models []ollamaModelInfo `json:"models"`
}

// GenerateResponse generates a response using Ollama
func (p *OllamaProvider) GenerateResponse(ctx context.Context, messages []Message) (*Message, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("ai-service").Start(ctx, "ollama.GenerateResponse")
	defer span.End()

	// Convert messages to Ollama format
	ollamaMessages := make([]ollamaChatMessage, len(messages))
	for i, msg := range messages {
		ollamaMessages[i] = ollamaChatMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		}
	}

	// Add system message for browser functions
	systemMessage := ollamaChatMessage{
		Role: "system",
		Content: `You are an AI assistant that can help users browse the web and analyze content.

Available functions:
- navigate_to_url: Navigate to a specific URL
- extract_page_content: Extract content from the current page
- take_screenshot: Take a screenshot of the current page
- click_element: Click on a page element
- fill_form_field: Fill out form fields
- analyze_page_structure: Analyze the structure of a page

Always be helpful, accurate, and provide step-by-step guidance when needed.`,
	}

	allMessages := append([]ollamaChatMessage{systemMessage}, ollamaMessages...)

	// Create request
	request := ollamaChatRequest{
		Model:    p.model,
		Messages: allMessages,
		Stream:   false,
		Options: map[string]interface{}{
			"temperature": p.config.Temperature,
			"top_p":       p.config.TopP,
			"top_k":       p.config.TopK,
			"num_ctx":     p.config.NumCtx,
		},
		Tools: p.getBrowserTools(),
	}

	// Make API call
	response, err := p.makeAPICall(ctx, "/api/chat", request)
	if err != nil {
		p.logger.Error(ctx, "Ollama API call failed", err)
		return nil, fmt.Errorf("Ollama API call failed: %w", err)
	}

	if response.Error != "" {
		return nil, fmt.Errorf("Ollama API error: %s", response.Error)
	}

	// Create response message
	message := &Message{
		Role:      RoleAssistant,
		Content:   response.Message.Content,
		CreatedAt: time.Now(),
		Metadata: map[string]interface{}{
			"model":      response.Model,
			"created_at": response.CreatedAt,
			"provider":   "ollama",
		},
	}

	return message, nil
}

// AnalyzeContent analyzes content using Ollama
func (p *OllamaProvider) AnalyzeContent(ctx context.Context, content string, analysisType string) (map[string]interface{}, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("ai-service").Start(ctx, "ollama.AnalyzeContent")
	defer span.End()

	prompt := p.buildAnalysisPrompt(content, analysisType)

	request := ollamaGenerateRequest{
		Model:  p.model,
		Prompt: prompt,
		Stream: false,
		Options: map[string]interface{}{
			"temperature": 0.3,
			"top_p":       0.9,
		},
	}

	response, err := p.makeGenerateAPICall(ctx, "/api/generate", request)
	if err != nil {
		return nil, fmt.Errorf("content analysis failed: %w", err)
	}

	if response.Error != "" {
		return nil, fmt.Errorf("Ollama API error: %s", response.Error)
	}

	// Try to parse as JSON, fallback to plain text
	var analysis map[string]interface{}
	if err := json.Unmarshal([]byte(response.Response), &analysis); err != nil {
		// If not valid JSON, return as plain text analysis
		analysis = map[string]interface{}{
			"summary":        response.Response,
			"analysis_type":  analysisType,
			"content_length": len(content),
		}
	}

	return analysis, nil
}

// ExtractStructuredData extracts structured data using Ollama
func (p *OllamaProvider) ExtractStructuredData(ctx context.Context, content string, schema string) (map[string]interface{}, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("ai-service").Start(ctx, "ollama.ExtractStructuredData")
	defer span.End()

	prompt := fmt.Sprintf(`Extract structured data from the following content according to the provided schema.
Return the result as valid JSON only, no additional text.

Schema:
%s

Content:
%s`, schema, content)

	request := ollamaGenerateRequest{
		Model:  p.model,
		Prompt: prompt,
		Stream: false,
		Options: map[string]interface{}{
			"temperature": 0.1,
			"top_p":       0.9,
		},
	}

	response, err := p.makeGenerateAPICall(ctx, "/api/generate", request)
	if err != nil {
		return nil, fmt.Errorf("data extraction failed: %w", err)
	}

	if response.Error != "" {
		return nil, fmt.Errorf("Ollama API error: %s", response.Error)
	}

	// Parse JSON response
	var extractedData map[string]interface{}
	if err := json.Unmarshal([]byte(response.Response), &extractedData); err != nil {
		return nil, fmt.Errorf("failed to parse extracted data as JSON: %w", err)
	}

	return extractedData, nil
}

// SummarizeContent summarizes content using Ollama
func (p *OllamaProvider) SummarizeContent(ctx context.Context, content string, options SummarizeOptions) (string, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("ai-service").Start(ctx, "ollama.SummarizeContent")
	defer span.End()

	prompt := p.buildSummarizationPrompt(content, options)

	request := ollamaGenerateRequest{
		Model:  p.model,
		Prompt: prompt,
		Stream: false,
		Options: map[string]interface{}{
			"temperature": 0.5,
			"top_p":       0.9,
		},
	}

	response, err := p.makeGenerateAPICall(ctx, "/api/generate", request)
	if err != nil {
		return "", fmt.Errorf("content summarization failed: %w", err)
	}

	if response.Error != "" {
		return "", fmt.Errorf("Ollama API error: %s", response.Error)
	}

	return strings.TrimSpace(response.Response), nil
}

// Helper methods

// makeAPICall makes an HTTP request to Ollama chat API
func (p *OllamaProvider) makeAPICall(ctx context.Context, endpoint string, request ollamaChatRequest) (*ollamaChatResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var response ollamaChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// makeGenerateAPICall makes an HTTP request to Ollama generate API
func (p *OllamaProvider) makeGenerateAPICall(ctx context.Context, endpoint string, request ollamaGenerateRequest) (*ollamaGenerateResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var response ollamaGenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// getBrowserTools returns browser function definitions for Ollama
func (p *OllamaProvider) getBrowserTools() []ollamaToolDefinition {
	return []ollamaToolDefinition{
		{
			Type: "function",
			Function: ollamaFunctionDef{
				Name:        "navigate_to_url",
				Description: "Navigate to a specific URL",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"url": map[string]interface{}{
							"type":        "string",
							"description": "The URL to navigate to",
						},
					},
					"required": []string{"url"},
				},
			},
		},
		{
			Type: "function",
			Function: ollamaFunctionDef{
				Name:        "extract_page_content",
				Description: "Extract content from the current page",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"selector": map[string]interface{}{
							"type":        "string",
							"description": "CSS selector to extract specific content (optional)",
						},
					},
				},
			},
		},
		{
			Type: "function",
			Function: ollamaFunctionDef{
				Name:        "take_screenshot",
				Description: "Take a screenshot of the current page",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"full_page": map[string]interface{}{
							"type":        "boolean",
							"description": "Whether to capture the full page or just the viewport",
						},
					},
				},
			},
		},
		{
			Type: "function",
			Function: ollamaFunctionDef{
				Name:        "click_element",
				Description: "Click on a page element",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"selector": map[string]interface{}{
							"type":        "string",
							"description": "CSS selector for the element to click",
						},
					},
					"required": []string{"selector"},
				},
			},
		},
	}
}

// buildAnalysisPrompt builds a prompt for content analysis
func (p *OllamaProvider) buildAnalysisPrompt(content string, analysisType string) string {
	switch analysisType {
	case "sentiment":
		return fmt.Sprintf(`Analyze the sentiment of the following content. Return a JSON object with sentiment (positive/negative/neutral), confidence score (0-1), and key phrases.

Content:
%s`, content)
	case "summary":
		return fmt.Sprintf(`Provide a concise summary of the following content. Return a JSON object with summary, key_points (array), and word_count.

Content:
%s`, content)
	case "keywords":
		return fmt.Sprintf(`Extract keywords and key phrases from the following content. Return a JSON object with keywords (array) and topics (array).

Content:
%s`, content)
	default:
		return fmt.Sprintf(`Analyze the following content and provide insights. Return a JSON object with analysis, key_insights (array), and content_type.

Content:
%s`, content)
	}
}

// buildSummarizationPrompt builds a prompt for content summarization
func (p *OllamaProvider) buildSummarizationPrompt(content string, options SummarizeOptions) string {
	lengthInstruction := "Provide a summary"
	if options.Length != "" {
		switch options.Length {
		case "short":
			lengthInstruction = "Provide a brief, concise summary (1-2 sentences)"
		case "medium":
			lengthInstruction = "Provide a moderate summary (3-5 sentences)"
		case "long":
			lengthInstruction = "Provide a detailed summary (1-2 paragraphs)"
		}
	}

	focusInstruction := ""
	if options.Focus != "" {
		focusInstruction = fmt.Sprintf(" Focus specifically on: %s.", options.Focus)
	}

	styleInstruction := ""
	if options.Style != "" {
		styleInstruction = fmt.Sprintf(" Use a %s style.", options.Style)
	}

	return fmt.Sprintf(`%s of the following content.%s%s

Content:
%s`, lengthInstruction, focusInstruction, styleInstruction, content)
}

// Health check methods

// IsHealthy checks if the Ollama service is available
func (p *OllamaProvider) IsHealthy(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", p.baseURL+"/api/tags", nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status %d", resp.StatusCode)
	}

	return nil
}

// ListModels returns available models from Ollama
func (p *OllamaProvider) ListModels(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", p.baseURL+"/api/tags", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list models failed with status %d", resp.StatusCode)
	}

	var modelsResp ollamaModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&modelsResp); err != nil {
		return nil, fmt.Errorf("failed to decode models response: %w", err)
	}

	models := make([]string, len(modelsResp.Models))
	for i, model := range modelsResp.Models {
		models[i] = model.Name
	}

	return models, nil
}
