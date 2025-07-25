package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
)

// AnthropicProvider implements AIProvider for Anthropic Claude
type AnthropicProvider struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
	logger     *observability.Logger
}

// NewAnthropicProvider creates a new Anthropic provider
func NewAnthropicProvider(apiKey, model string, logger *observability.Logger) *AnthropicProvider {
	if model == "" {
		model = "claude-3-sonnet-20240229"
	}

	return &AnthropicProvider{
		apiKey:  apiKey,
		baseURL: "https://api.anthropic.com/v1",
		model:   model,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		logger: logger,
	}
}

// Anthropic API structures
type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	Messages  []anthropicMessage `json:"messages"`
	System    string             `json:"system,omitempty"`
	Tools     []anthropicTool    `json:"tools,omitempty"`
}

type anthropicMessage struct {
	Role    string                   `json:"role"`
	Content []anthropicContentBlock `json:"content"`
}

type anthropicContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

type anthropicTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"input_schema"`
}

type anthropicResponse struct {
	ID      string                   `json:"id"`
	Type    string                   `json:"type"`
	Role    string                   `json:"role"`
	Content []anthropicContentBlock `json:"content"`
	Model   string                   `json:"model"`
	Usage   anthropicUsage           `json:"usage"`
	Error   *anthropicError          `json:"error,omitempty"`
}

type anthropicUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type anthropicError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// GenerateResponse generates an AI response using Anthropic Claude
func (p *AnthropicProvider) GenerateResponse(ctx context.Context, messages []Message) (*Message, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("ai-service").Start(ctx, "anthropic.GenerateResponse")
	defer span.End()

	// Convert messages to Anthropic format
	anthropicMessages := make([]anthropicMessage, 0, len(messages))
	var systemMessage string

	for _, msg := range messages {
		if msg.Role == RoleSystem {
			systemMessage = msg.Content
			continue
		}

		anthropicMsg := anthropicMessage{
			Role: string(msg.Role),
			Content: []anthropicContentBlock{
				{
					Type: "text",
					Text: msg.Content,
				},
			},
		}
		anthropicMessages = append(anthropicMessages, anthropicMsg)
	}

	// Add default system message if none provided
	if systemMessage == "" {
		systemMessage = `You are an AI assistant that helps users navigate and interact with websites. 
You can perform various tasks like navigating to URLs, extracting content, filling forms, 
taking screenshots, and analyzing web pages. When users ask you to perform web actions, 
you should break down their requests into specific tasks and provide clear instructions.

Always be helpful, accurate, and provide step-by-step guidance when needed.`
	}

	// Create request
	request := anthropicRequest{
		Model:     p.model,
		MaxTokens: 1000,
		Messages:  anthropicMessages,
		System:    systemMessage,
		Tools:     p.getBrowserTools(),
	}

	// Make API call
	response, err := p.makeAPICall(ctx, "/messages", request)
	if err != nil {
		p.logger.Error(ctx, "Anthropic API call failed", err)
		return nil, fmt.Errorf("Anthropic API call failed: %w", err)
	}

	// Extract content from response
	var content string
	if len(response.Content) > 0 && response.Content[0].Type == "text" {
		content = response.Content[0].Text
	}

	// Create response message
	responseMessage := &Message{
		Role:    RoleAssistant,
		Content: content,
		Metadata: map[string]interface{}{
			"model":         response.Model,
			"provider":      "anthropic",
			"input_tokens":  response.Usage.InputTokens,
			"output_tokens": response.Usage.OutputTokens,
		},
	}

	p.logger.Info(ctx, "Anthropic response generated", map[string]interface{}{
		"model":         response.Model,
		"input_tokens":  response.Usage.InputTokens,
		"output_tokens": response.Usage.OutputTokens,
	})

	return responseMessage, nil
}

// AnalyzeContent analyzes content using Anthropic Claude
func (p *AnthropicProvider) AnalyzeContent(ctx context.Context, content string, analysisType string) (map[string]interface{}, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("ai-service").Start(ctx, "anthropic.AnalyzeContent")
	defer span.End()

	prompt := p.buildAnalysisPrompt(content, analysisType)

	messages := []anthropicMessage{
		{
			Role: "user",
			Content: []anthropicContentBlock{
				{
					Type: "text",
					Text: prompt,
				},
			},
		},
	}

	request := anthropicRequest{
		Model:     p.model,
		MaxTokens: 1000,
		Messages:  messages,
		System:    "You are an expert content analyzer. Provide structured analysis in JSON format.",
	}

	response, err := p.makeAPICall(ctx, "/messages", request)
	if err != nil {
		return nil, fmt.Errorf("content analysis failed: %w", err)
	}

	if len(response.Content) == 0 {
		return nil, fmt.Errorf("no analysis response received")
	}

	// Parse JSON response
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(response.Content[0].Text), &result); err != nil {
		// If JSON parsing fails, return the raw content
		result = map[string]interface{}{
			"analysis": response.Content[0].Text,
			"type":     analysisType,
		}
	}

	return result, nil
}

// ExtractStructuredData extracts structured data using Anthropic Claude
func (p *AnthropicProvider) ExtractStructuredData(ctx context.Context, content string, schema string) (map[string]interface{}, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("ai-service").Start(ctx, "anthropic.ExtractStructuredData")
	defer span.End()

	prompt := fmt.Sprintf(`Extract structured data from the following content according to the provided schema.
Return the result as valid JSON.

Schema:
%s

Content:
%s`, schema, content)

	messages := []anthropicMessage{
		{
			Role: "user",
			Content: []anthropicContentBlock{
				{
					Type: "text",
					Text: prompt,
				},
			},
		},
	}

	request := anthropicRequest{
		Model:     p.model,
		MaxTokens: 1500,
		Messages:  messages,
		System:    "You are a data extraction expert. Extract structured data according to the provided schema and return valid JSON.",
	}

	response, err := p.makeAPICall(ctx, "/messages", request)
	if err != nil {
		return nil, fmt.Errorf("data extraction failed: %w", err)
	}

	if len(response.Content) == 0 {
		return nil, fmt.Errorf("no extraction response received")
	}

	// Parse JSON response
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(response.Content[0].Text), &result); err != nil {
		return nil, fmt.Errorf("failed to parse extracted data as JSON: %w", err)
	}

	return result, nil
}

// SummarizeContent summarizes content using Anthropic Claude
func (p *AnthropicProvider) SummarizeContent(ctx context.Context, content string, options SummarizeOptions) (string, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("ai-service").Start(ctx, "anthropic.SummarizeContent")
	defer span.End()

	prompt := p.buildSummarizationPrompt(content, options)

	messages := []anthropicMessage{
		{
			Role: "user",
			Content: []anthropicContentBlock{
				{
					Type: "text",
					Text: prompt,
				},
			},
		},
	}

	request := anthropicRequest{
		Model:     p.model,
		MaxTokens: 800,
		Messages:  messages,
		System:    "You are an expert content summarizer. Provide clear, concise summaries based on the specified requirements.",
	}

	response, err := p.makeAPICall(ctx, "/messages", request)
	if err != nil {
		return "", fmt.Errorf("content summarization failed: %w", err)
	}

	if len(response.Content) == 0 {
		return "", fmt.Errorf("no summarization response received")
	}

	return response.Content[0].Text, nil
}

// makeAPICall makes an HTTP request to the Anthropic API
func (p *AnthropicProvider) makeAPICall(ctx context.Context, endpoint string, request interface{}) (*anthropicResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	var response anthropicResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("Anthropic API error: %s", response.Error.Message)
	}

	return &response, nil
}

// getBrowserTools returns tool definitions for browser automation
func (p *AnthropicProvider) getBrowserTools() []anthropicTool {
	return []anthropicTool{
		{
			Name:        "navigate_to_url",
			Description: "Navigate to a specific URL",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"url": map[string]interface{}{
						"type":        "string",
						"description": "The URL to navigate to",
					},
					"wait_for_selector": map[string]interface{}{
						"type":        "string",
						"description": "CSS selector to wait for after navigation",
					},
				},
				"required": []string{"url"},
			},
		},
		{
			Name:        "extract_content",
			Description: "Extract content from the current page",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"selectors": map[string]interface{}{
						"type":        "array",
						"items":       map[string]interface{}{"type": "string"},
						"description": "CSS selectors for content to extract",
					},
					"data_type": map[string]interface{}{
						"type":        "string",
						"description": "Type of data to extract (text, links, images, tables)",
					},
				},
			},
		},
		{
			Name:        "interact_with_page",
			Description: "Interact with page elements (click, type, etc.)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"actions": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"type": map[string]interface{}{
									"type": "string",
									"enum": []string{"click", "type", "select", "scroll", "wait"},
								},
								"selector": map[string]interface{}{
									"type": "string",
								},
								"value": map[string]interface{}{
									"type": "string",
								},
							},
							"required": []string{"type"},
						},
					},
				},
				"required": []string{"actions"},
			},
		},
	}
}

// Helper methods for building prompts (similar to OpenAI provider)
func (p *AnthropicProvider) buildAnalysisPrompt(content string, analysisType string) string {
	switch analysisType {
	case "sentiment":
		return fmt.Sprintf("Analyze the sentiment of this content and return a JSON object with 'sentiment' (positive/negative/neutral), 'confidence' (0-1), and 'key_phrases':\n\n%s", content)
	case "keywords":
		return fmt.Sprintf("Extract the main keywords and key phrases from this content. Return a JSON object with 'keywords' array and 'topics' array:\n\n%s", content)
	case "entities":
		return fmt.Sprintf("Extract named entities from this content. Return a JSON object with entities categorized by type (person, organization, location, etc.):\n\n%s", content)
	case "structure":
		return fmt.Sprintf("Analyze the structure and organization of this content. Return a JSON object with 'sections', 'headings', 'main_topics', and 'content_type':\n\n%s", content)
	default:
		return fmt.Sprintf("Provide a comprehensive analysis of this content. Return a JSON object with relevant insights:\n\n%s", content)
	}
}

func (p *AnthropicProvider) buildSummarizationPrompt(content string, options SummarizeOptions) string {
	lengthInstruction := ""
	switch options.Length {
	case "short":
		lengthInstruction = "Provide a brief summary in 1-2 sentences."
	case "medium":
		lengthInstruction = "Provide a moderate summary in 3-5 sentences."
	case "long":
		lengthInstruction = "Provide a detailed summary in 1-2 paragraphs."
	default:
		lengthInstruction = "Provide an appropriate length summary."
	}

	focusInstruction := ""
	switch options.Focus {
	case "main_points":
		focusInstruction = "Focus on the main points and key takeaways."
	case "technical":
		focusInstruction = "Focus on technical details and specifications."
	case "business":
		focusInstruction = "Focus on business implications and strategic insights."
	default:
		focusInstruction = "Focus on the most important information."
	}

	return fmt.Sprintf("%s %s\n\nContent to summarize:\n%s", lengthInstruction, focusInstruction, content)
}
