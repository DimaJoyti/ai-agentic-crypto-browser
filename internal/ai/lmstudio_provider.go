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

// LMStudioProvider implements AIProvider for LM Studio
type LMStudioProvider struct {
	baseURL    string
	model      string
	httpClient *http.Client
	logger     *observability.Logger
	config     LMStudioConfig
}

// LMStudioConfig holds LM Studio-specific configuration
type LMStudioConfig struct {
	BaseURL     string
	Model       string
	Temperature float64
	MaxTokens   int
	TopP        float64
	Timeout     time.Duration
	MaxRetries  int
	RetryDelay  time.Duration
}

// NewLMStudioProvider creates a new LM Studio provider
func NewLMStudioProvider(config LMStudioConfig, logger *observability.Logger) *LMStudioProvider {
	if config.BaseURL == "" {
		config.BaseURL = "http://localhost:1234/v1"
	}
	if config.Model == "" {
		config.Model = "local-model"
	}
	if config.Temperature == 0 {
		config.Temperature = 0.7
	}
	if config.MaxTokens == 0 {
		config.MaxTokens = 1000
	}
	if config.Timeout == 0 {
		config.Timeout = 60 * time.Second
	}

	return &LMStudioProvider{
		baseURL: config.BaseURL,
		model:   config.Model,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		logger: logger,
		config: config,
	}
}

// LM Studio API request/response structures (OpenAI-compatible)
type lmStudioChatRequest struct {
	Model        string                `json:"model"`
	Messages     []lmStudioMessage     `json:"messages"`
	Temperature  float64               `json:"temperature,omitempty"`
	MaxTokens    int                   `json:"max_tokens,omitempty"`
	TopP         float64               `json:"top_p,omitempty"`
	Stream       bool                  `json:"stream,omitempty"`
	Functions    []lmStudioFunctionDef `json:"functions,omitempty"`
	FunctionCall interface{}           `json:"function_call,omitempty"`
}

type lmStudioMessage struct {
	Role         string                `json:"role"`
	Content      string                `json:"content"`
	FunctionCall *lmStudioFunctionCall `json:"function_call,omitempty"`
}

type lmStudioFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type lmStudioFunctionDef struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type lmStudioChatResponse struct {
	ID      string           `json:"id"`
	Object  string           `json:"object"`
	Created int64            `json:"created"`
	Model   string           `json:"model"`
	Choices []lmStudioChoice `json:"choices"`
	Usage   lmStudioUsage    `json:"usage"`
	Error   *lmStudioError   `json:"error,omitempty"`
}

type lmStudioChoice struct {
	Index        int             `json:"index"`
	Message      lmStudioMessage `json:"message"`
	FinishReason string          `json:"finish_reason"`
}

type lmStudioUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type lmStudioError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

type lmStudioModelsResponse struct {
	Object string              `json:"object"`
	Data   []lmStudioModelInfo `json:"data"`
}

type lmStudioModelInfo struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

// GenerateResponse generates a response using LM Studio
func (p *LMStudioProvider) GenerateResponse(ctx context.Context, messages []Message) (*Message, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("ai-service").Start(ctx, "lmstudio.GenerateResponse")
	defer span.End()

	// Convert messages to LM Studio format
	lmStudioMessages := make([]lmStudioMessage, len(messages))
	for i, msg := range messages {
		lmStudioMessages[i] = lmStudioMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		}
	}

	// Add system message for browser functions
	systemMessage := lmStudioMessage{
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

	allMessages := append([]lmStudioMessage{systemMessage}, lmStudioMessages...)

	// Create request
	request := lmStudioChatRequest{
		Model:        p.model,
		Messages:     allMessages,
		Temperature:  p.config.Temperature,
		MaxTokens:    p.config.MaxTokens,
		TopP:         p.config.TopP,
		Stream:       false,
		Functions:    p.getBrowserFunctions(),
		FunctionCall: "auto",
	}

	// Make API call
	response, err := p.makeAPICall(ctx, "/chat/completions", request)
	if err != nil {
		p.logger.Error(ctx, "LM Studio API call failed", err)
		return nil, fmt.Errorf("LM Studio API call failed: %w", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("LM Studio API error: %s", response.Error.Message)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no response choices returned from LM Studio")
	}

	choice := response.Choices[0]

	// Create response message
	message := &Message{
		Role:      RoleAssistant,
		Content:   choice.Message.Content,
		CreatedAt: time.Now(),
		Metadata: map[string]interface{}{
			"model":             response.Model,
			"finish_reason":     choice.FinishReason,
			"prompt_tokens":     response.Usage.PromptTokens,
			"completion_tokens": response.Usage.CompletionTokens,
			"total_tokens":      response.Usage.TotalTokens,
			"provider":          "lmstudio",
		},
	}

	// Handle function calls
	if choice.Message.FunctionCall != nil {
		message.Metadata["function_call"] = map[string]interface{}{
			"name":      choice.Message.FunctionCall.Name,
			"arguments": choice.Message.FunctionCall.Arguments,
		}
	}

	return message, nil
}

// AnalyzeContent analyzes content using LM Studio
func (p *LMStudioProvider) AnalyzeContent(ctx context.Context, content string, analysisType string) (map[string]interface{}, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("ai-service").Start(ctx, "lmstudio.AnalyzeContent")
	defer span.End()

	prompt := p.buildAnalysisPrompt(content, analysisType)

	messages := []lmStudioMessage{
		{
			Role:    "system",
			Content: "You are an expert content analyzer. Provide structured analysis in JSON format.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	request := lmStudioChatRequest{
		Model:       p.model,
		Messages:    messages,
		Temperature: 0.3,
		MaxTokens:   1000,
		TopP:        0.9,
	}

	response, err := p.makeAPICall(ctx, "/chat/completions", request)
	if err != nil {
		return nil, fmt.Errorf("content analysis failed: %w", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("LM Studio API error: %s", response.Error.Message)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no analysis response received")
	}

	content = response.Choices[0].Message.Content

	// Try to parse as JSON, fallback to plain text
	var analysis map[string]interface{}
	if err := json.Unmarshal([]byte(content), &analysis); err != nil {
		// If not valid JSON, return as plain text analysis
		analysis = map[string]interface{}{
			"summary":        content,
			"analysis_type":  analysisType,
			"content_length": len(content),
		}
	}

	return analysis, nil
}

// ExtractStructuredData extracts structured data using LM Studio
func (p *LMStudioProvider) ExtractStructuredData(ctx context.Context, content string, schema string) (map[string]interface{}, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("ai-service").Start(ctx, "lmstudio.ExtractStructuredData")
	defer span.End()

	prompt := fmt.Sprintf(`Extract structured data from the following content according to the provided schema.
Return the result as valid JSON only, no additional text.

Schema:
%s

Content:
%s`, schema, content)

	messages := []lmStudioMessage{
		{
			Role:    "system",
			Content: "You are a data extraction expert. Extract structured data according to the provided schema and return valid JSON.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	request := lmStudioChatRequest{
		Model:       p.model,
		Messages:    messages,
		Temperature: 0.1,
		MaxTokens:   1500,
		TopP:        0.9,
	}

	response, err := p.makeAPICall(ctx, "/chat/completions", request)
	if err != nil {
		return nil, fmt.Errorf("data extraction failed: %w", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("LM Studio API error: %s", response.Error.Message)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no extraction response received")
	}

	responseContent := response.Choices[0].Message.Content

	// Parse JSON response
	var extractedData map[string]interface{}
	if err := json.Unmarshal([]byte(responseContent), &extractedData); err != nil {
		return nil, fmt.Errorf("failed to parse extracted data as JSON: %w", err)
	}

	return extractedData, nil
}

// SummarizeContent summarizes content using LM Studio
func (p *LMStudioProvider) SummarizeContent(ctx context.Context, content string, options SummarizeOptions) (string, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("ai-service").Start(ctx, "lmstudio.SummarizeContent")
	defer span.End()

	prompt := p.buildSummarizationPrompt(content, options)

	messages := []lmStudioMessage{
		{
			Role:    "system",
			Content: "You are an expert content summarizer. Provide clear, concise summaries based on the specified requirements.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	request := lmStudioChatRequest{
		Model:       p.model,
		Messages:    messages,
		Temperature: 0.5,
		MaxTokens:   800,
		TopP:        0.9,
	}

	response, err := p.makeAPICall(ctx, "/chat/completions", request)
	if err != nil {
		return "", fmt.Errorf("content summarization failed: %w", err)
	}

	if response.Error != nil {
		return "", fmt.Errorf("LM Studio API error: %s", response.Error.Message)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no summarization response received")
	}

	return strings.TrimSpace(response.Choices[0].Message.Content), nil
}

// Helper methods

// makeAPICall makes an HTTP request to LM Studio API
func (p *LMStudioProvider) makeAPICall(ctx context.Context, endpoint string, request lmStudioChatRequest) (*lmStudioChatResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer lm-studio") // LM Studio doesn't require real auth

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var response lmStudioChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// getBrowserFunctions returns browser function definitions for LM Studio
func (p *LMStudioProvider) getBrowserFunctions() []lmStudioFunctionDef {
	return []lmStudioFunctionDef{
		{
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
		{
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
		{
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
		{
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
		{
			Name:        "fill_form_field",
			Description: "Fill out a form field",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"selector": map[string]interface{}{
						"type":        "string",
						"description": "CSS selector for the form field",
					},
					"value": map[string]interface{}{
						"type":        "string",
						"description": "Value to fill in the field",
					},
				},
				"required": []string{"selector", "value"},
			},
		},
	}
}

// buildAnalysisPrompt builds a prompt for content analysis
func (p *LMStudioProvider) buildAnalysisPrompt(content string, analysisType string) string {
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
	case "structure":
		return fmt.Sprintf(`Analyze the structure and organization of the following content. Return a JSON object with structure_type, sections (array), and hierarchy.

Content:
%s`, content)
	default:
		return fmt.Sprintf(`Analyze the following content and provide insights. Return a JSON object with analysis, key_insights (array), and content_type.

Content:
%s`, content)
	}
}

// buildSummarizationPrompt builds a prompt for content summarization
func (p *LMStudioProvider) buildSummarizationPrompt(content string, options SummarizeOptions) string {
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

// IsHealthy checks if the LM Studio service is available
func (p *LMStudioProvider) IsHealthy(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", p.baseURL+"/models", nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer lm-studio")

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

// ListModels returns available models from LM Studio
func (p *LMStudioProvider) ListModels(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", p.baseURL+"/models", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer lm-studio")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list models failed with status %d", resp.StatusCode)
	}

	var modelsResp lmStudioModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&modelsResp); err != nil {
		return nil, fmt.Errorf("failed to decode models response: %w", err)
	}

	models := make([]string, len(modelsResp.Data))
	for i, model := range modelsResp.Data {
		models[i] = model.ID
	}

	return models, nil
}
