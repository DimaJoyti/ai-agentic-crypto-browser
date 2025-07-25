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

// OpenAIProvider implements AIProvider for OpenAI
type OpenAIProvider struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
	logger     *observability.Logger
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(apiKey, model string, logger *observability.Logger) *OpenAIProvider {
	return &OpenAIProvider{
		apiKey:  apiKey,
		baseURL: "https://api.openai.com/v1",
		model:   model,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// OpenAI API structures
type openAIRequest struct {
	Model       string              `json:"model"`
	Messages    []openAIMessage     `json:"messages"`
	Temperature float64             `json:"temperature,omitempty"`
	MaxTokens   int                 `json:"max_tokens,omitempty"`
	Functions   []openAIFunction    `json:"functions,omitempty"`
	FunctionCall interface{}        `json:"function_call,omitempty"`
}

type openAIMessage struct {
	Role         string                 `json:"role"`
	Content      string                 `json:"content"`
	FunctionCall *openAIFunctionCall    `json:"function_call,omitempty"`
	Name         string                 `json:"name,omitempty"`
}

type openAIFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type openAIFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type openAIResponse struct {
	ID      string           `json:"id"`
	Object  string           `json:"object"`
	Created int64            `json:"created"`
	Model   string           `json:"model"`
	Choices []openAIChoice   `json:"choices"`
	Usage   openAIUsage      `json:"usage"`
	Error   *openAIError     `json:"error,omitempty"`
}

type openAIChoice struct {
	Index        int           `json:"index"`
	Message      openAIMessage `json:"message"`
	FinishReason string        `json:"finish_reason"`
}

type openAIUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type openAIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

// GenerateResponse generates an AI response for the given messages
func (p *OpenAIProvider) GenerateResponse(ctx context.Context, messages []Message) (*Message, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("ai-service").Start(ctx, "openai.GenerateResponse")
	defer span.End()

	// Convert messages to OpenAI format
	openAIMessages := make([]openAIMessage, len(messages))
	for i, msg := range messages {
		openAIMessages[i] = openAIMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		}
	}

	// Add system message for browser agent context
	systemMessage := openAIMessage{
		Role: "system",
		Content: `You are an AI assistant that helps users navigate and interact with websites. 
You can perform various tasks like navigating to URLs, extracting content, filling forms, 
taking screenshots, and analyzing web pages. When users ask you to perform web actions, 
you should break down their requests into specific tasks and provide clear instructions.

Available task types:
- navigate: Go to a specific URL
- extract: Extract content from a page
- interact: Click, type, or interact with page elements
- summarize: Summarize page content
- search: Search for information
- fill_form: Fill out forms on websites
- screenshot: Take screenshots of pages
- analyze: Analyze page content for specific criteria

Always be helpful, accurate, and provide step-by-step guidance when needed.`,
	}

	allMessages := append([]openAIMessage{systemMessage}, openAIMessages...)

	// Create request
	request := openAIRequest{
		Model:       p.model,
		Messages:    allMessages,
		Temperature: 0.7,
		MaxTokens:   1000,
		Functions:   p.getBrowserFunctions(),
		FunctionCall: "auto",
	}

	// Make API call
	response, err := p.makeAPICall(ctx, "/chat/completions", request)
	if err != nil {
		p.logger.Error(ctx, "OpenAI API call failed", err)
		return nil, fmt.Errorf("OpenAI API call failed: %w", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no response choices returned from OpenAI")
	}

	choice := response.Choices[0]
	
	// Create response message
	responseMessage := &Message{
		Role:    RoleAssistant,
		Content: choice.Message.Content,
		Metadata: map[string]interface{}{
			"model":        response.Model,
			"finish_reason": choice.FinishReason,
			"usage":        response.Usage,
		},
	}

	// Handle function calls if present
	if choice.Message.FunctionCall != nil {
		responseMessage.Metadata["function_call"] = choice.Message.FunctionCall
	}

	p.logger.Info(ctx, "OpenAI response generated", map[string]interface{}{
		"model":         response.Model,
		"total_tokens":  response.Usage.TotalTokens,
		"finish_reason": choice.FinishReason,
	})

	return responseMessage, nil
}

// AnalyzeContent analyzes content using OpenAI
func (p *OpenAIProvider) AnalyzeContent(ctx context.Context, content string, analysisType string) (map[string]interface{}, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("ai-service").Start(ctx, "openai.AnalyzeContent")
	defer span.End()

	prompt := p.buildAnalysisPrompt(content, analysisType)
	
	messages := []openAIMessage{
		{
			Role:    "system",
			Content: "You are an expert content analyzer. Provide structured analysis in JSON format.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	request := openAIRequest{
		Model:       p.model,
		Messages:    messages,
		Temperature: 0.3,
		MaxTokens:   1000,
	}

	response, err := p.makeAPICall(ctx, "/chat/completions", request)
	if err != nil {
		return nil, fmt.Errorf("content analysis failed: %w", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no analysis response received")
	}

	// Parse JSON response
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(response.Choices[0].Message.Content), &result); err != nil {
		// If JSON parsing fails, return the raw content
		result = map[string]interface{}{
			"analysis": response.Choices[0].Message.Content,
			"type":     analysisType,
		}
	}

	return result, nil
}

// ExtractStructuredData extracts structured data using OpenAI
func (p *OpenAIProvider) ExtractStructuredData(ctx context.Context, content string, schema string) (map[string]interface{}, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("ai-service").Start(ctx, "openai.ExtractStructuredData")
	defer span.End()

	prompt := fmt.Sprintf(`Extract structured data from the following content according to the provided schema.
Return the result as valid JSON.

Schema:
%s

Content:
%s`, schema, content)

	messages := []openAIMessage{
		{
			Role:    "system",
			Content: "You are a data extraction expert. Extract structured data according to the provided schema and return valid JSON.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	request := openAIRequest{
		Model:       p.model,
		Messages:    messages,
		Temperature: 0.1,
		MaxTokens:   1500,
	}

	response, err := p.makeAPICall(ctx, "/chat/completions", request)
	if err != nil {
		return nil, fmt.Errorf("data extraction failed: %w", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no extraction response received")
	}

	// Parse JSON response
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(response.Choices[0].Message.Content), &result); err != nil {
		return nil, fmt.Errorf("failed to parse extracted data as JSON: %w", err)
	}

	return result, nil
}

// SummarizeContent summarizes content using OpenAI
func (p *OpenAIProvider) SummarizeContent(ctx context.Context, content string, options SummarizeOptions) (string, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("ai-service").Start(ctx, "openai.SummarizeContent")
	defer span.End()

	prompt := p.buildSummarizationPrompt(content, options)

	messages := []openAIMessage{
		{
			Role:    "system",
			Content: "You are an expert content summarizer. Provide clear, concise summaries based on the specified requirements.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	request := openAIRequest{
		Model:       p.model,
		Messages:    messages,
		Temperature: 0.5,
		MaxTokens:   800,
	}

	response, err := p.makeAPICall(ctx, "/chat/completions", request)
	if err != nil {
		return "", fmt.Errorf("content summarization failed: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no summarization response received")
	}

	return response.Choices[0].Message.Content, nil
}

// makeAPICall makes an HTTP request to the OpenAI API
func (p *OpenAIProvider) makeAPICall(ctx context.Context, endpoint string, request interface{}) (*openAIResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	var response openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("OpenAI API error: %s", response.Error.Message)
	}

	return &response, nil
}

// getBrowserFunctions returns function definitions for browser automation
func (p *OpenAIProvider) getBrowserFunctions() []openAIFunction {
	return []openAIFunction{
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
			Parameters: map[string]interface{}{
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
			Parameters: map[string]interface{}{
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

// Helper methods for building prompts
func (p *OpenAIProvider) buildAnalysisPrompt(content string, analysisType string) string {
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

func (p *OpenAIProvider) buildSummarizationPrompt(content string, options SummarizeOptions) string {
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
