package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ai-agentic-browser/internal/browser"
	"github.com/ai-agentic-browser/internal/config"
	"github.com/ai-agentic-browser/pkg/database"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

// Service provides AI agent functionality
type Service struct {
	db             *database.DB
	redis          *database.RedisClient
	aiProvider     AIProvider
	browserService *browser.Service
	logger         *observability.Logger
	healthMonitor  *HealthMonitor
}

// NewService creates a new AI service
func NewService(db *database.DB, redis *database.RedisClient, cfg config.AIConfig, browserService *browser.Service, logger *observability.Logger) *Service {
	var aiProvider AIProvider

	switch cfg.Provider {
	case "openai":
		aiProvider = NewOpenAIProvider(cfg.OpenAIKey, cfg.ModelName, logger)
	case "anthropic":
		aiProvider = NewAnthropicProvider(cfg.AnthropicKey, cfg.ModelName, logger)
	case "ollama":
		ollamaConfig := OllamaConfig{
			BaseURL:     cfg.OllamaConfig.BaseURL,
			Model:       cfg.OllamaConfig.Model,
			Temperature: cfg.OllamaConfig.Temperature,
			TopP:        cfg.OllamaConfig.TopP,
			TopK:        cfg.OllamaConfig.TopK,
			NumCtx:      cfg.OllamaConfig.NumCtx,
			Timeout:     cfg.OllamaConfig.Timeout,
			MaxRetries:  cfg.OllamaConfig.MaxRetries,
			RetryDelay:  cfg.OllamaConfig.RetryDelay,
		}
		aiProvider = NewOllamaProvider(ollamaConfig, logger)
	case "lmstudio":
		lmStudioConfig := LMStudioConfig{
			BaseURL:     cfg.LMStudioConfig.BaseURL,
			Model:       cfg.LMStudioConfig.Model,
			Temperature: cfg.LMStudioConfig.Temperature,
			MaxTokens:   cfg.LMStudioConfig.MaxTokens,
			TopP:        cfg.LMStudioConfig.TopP,
			Timeout:     cfg.LMStudioConfig.Timeout,
			MaxRetries:  cfg.LMStudioConfig.MaxRetries,
			RetryDelay:  cfg.LMStudioConfig.RetryDelay,
		}
		aiProvider = NewLMStudioProvider(lmStudioConfig, logger)
	default:
		logger.Warn(context.Background(), "Unknown AI provider, falling back to OpenAI", map[string]interface{}{
			"provider": cfg.Provider,
		})
		aiProvider = NewOpenAIProvider(cfg.OpenAIKey, cfg.ModelName, logger)
	}

	// Initialize health monitor
	healthMonitor := NewHealthMonitor(logger, cfg.OllamaConfig.HealthCheckInterval)

	// Register providers for health monitoring
	if healthChecker, ok := aiProvider.(HealthChecker); ok {
		healthMonitor.RegisterProvider(cfg.Provider, healthChecker)
	}

	service := &Service{
		db:             db,
		redis:          redis,
		aiProvider:     aiProvider,
		browserService: browserService,
		logger:         logger,
		healthMonitor:  healthMonitor,
	}

	// Start health monitoring in background
	go healthMonitor.Start(context.Background())

	return service
}

// Chat handles chat messages and generates responses
func (s *Service) Chat(ctx context.Context, userID uuid.UUID, req ChatRequest) (*ChatResponse, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("ai-service").Start(ctx, "ai.Chat")
	defer span.End()

	// Get or create conversation
	conversation, err := s.getOrCreateConversation(ctx, userID, req.ConversationID, req.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	// Add user message
	userMessage := &Message{
		ID:             uuid.New(),
		ConversationID: conversation.ID,
		Role:           RoleUser,
		Content:        req.Message,
		CreatedAt:      time.Now(),
	}

	if err := s.saveMessage(ctx, userMessage); err != nil {
		return nil, fmt.Errorf("failed to save user message: %w", err)
	}

	// Get conversation history
	messages, err := s.getConversationMessages(ctx, conversation.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation messages: %w", err)
	}

	// Generate AI response
	aiMessage, err := s.aiProvider.GenerateResponse(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("failed to generate AI response: %w", err)
	}

	// Set message metadata
	aiMessage.ID = uuid.New()
	aiMessage.ConversationID = conversation.ID
	aiMessage.CreatedAt = time.Now()

	// Save AI message
	if err := s.saveMessage(ctx, aiMessage); err != nil {
		return nil, fmt.Errorf("failed to save AI message: %w", err)
	}

	// Parse function calls and create tasks if needed
	var tasks []Task
	if functionCall, exists := aiMessage.Metadata["function_call"]; exists {
		task, err := s.processFunctionCall(ctx, userID, conversation.ID, functionCall)
		if err != nil {
			s.logger.Error(ctx, "Failed to process function call", err)
		} else if task != nil {
			tasks = append(tasks, *task)
		}
	}

	// Generate suggestions
	suggestions := s.generateSuggestions(ctx, req.Message, aiMessage.Content)

	response := &ChatResponse{
		ConversationID: conversation.ID,
		Message:        *aiMessage,
		Tasks:          tasks,
		Suggestions:    suggestions,
	}

	s.logger.Info(ctx, "Chat response generated", map[string]interface{}{
		"conversation_id": conversation.ID.String(),
		"user_id":         userID.String(),
		"tasks_created":   len(tasks),
	})

	return response, nil
}

// CreateTask creates a new AI task
func (s *Service) CreateTask(ctx context.Context, userID uuid.UUID, req TaskRequest) (*TaskResponse, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("ai-service").Start(ctx, "ai.CreateTask")
	defer span.End()

	// Get or create conversation
	conversation, err := s.getOrCreateConversation(ctx, userID, req.ConversationID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	// Create task
	task := &Task{
		ID:             uuid.New(),
		ConversationID: conversation.ID,
		UserID:         userID,
		TaskType:       req.TaskType,
		Description:    req.Description,
		Status:         TaskStatusPending,
		InputData:      req.InputData,
		CreatedAt:      time.Now(),
	}

	// Save task
	if err := s.saveTask(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to save task: %w", err)
	}

	// Execute task asynchronously
	go s.executeTask(context.Background(), task)

	response := &TaskResponse{
		Task: *task,
	}

	s.logger.Info(ctx, "Task created", map[string]interface{}{
		"task_id":   task.ID.String(),
		"task_type": task.TaskType,
		"user_id":   userID.String(),
	})

	return response, nil
}

// executeTask executes a task based on its type
func (s *Service) executeTask(ctx context.Context, task *Task) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("ai-service").Start(ctx, "ai.ExecuteTask")
	defer span.End()

	// Update task status to running
	task.Status = TaskStatusRunning
	task.StartedAt = &[]time.Time{time.Now()}[0]
	s.updateTaskStatus(ctx, task)

	var err error
	var outputData map[string]interface{}

	switch task.TaskType {
	case TaskTypeNavigate:
		outputData, err = s.executeNavigateTask(ctx, task)
	case TaskTypeExtract:
		outputData, err = s.executeExtractTask(ctx, task)
	case TaskTypeInteract:
		outputData, err = s.executeInteractTask(ctx, task)
	case TaskTypeSummarize:
		outputData, err = s.executeSummarizeTask(ctx, task)
	case TaskTypeScreenshot:
		outputData, err = s.executeScreenshotTask(ctx, task)
	case TaskTypeAnalyze:
		outputData, err = s.executeAnalyzeTask(ctx, task)
	default:
		err = fmt.Errorf("unsupported task type: %s", task.TaskType)
	}

	// Update task with results
	now := time.Now()
	task.CompletedAt = &now
	task.OutputData = outputData

	if err != nil {
		task.Status = TaskStatusFailed
		task.ErrorMessage = &[]string{err.Error()}[0]
		s.logger.Error(ctx, "Task execution failed", err, map[string]interface{}{
			"task_id":   task.ID.String(),
			"task_type": task.TaskType,
		})
	} else {
		task.Status = TaskStatusCompleted
		s.logger.Info(ctx, "Task execution completed", map[string]interface{}{
			"task_id":   task.ID.String(),
			"task_type": task.TaskType,
		})
	}

	s.updateTaskStatus(ctx, task)
}

// executeNavigateTask executes a navigation task
func (s *Service) executeNavigateTask(ctx context.Context, task *Task) (map[string]interface{}, error) {
	var input NavigateTaskInput
	if err := s.parseTaskInput(task.InputData, &input); err != nil {
		return nil, fmt.Errorf("invalid navigate task input: %w", err)
	}

	// Create browser session for this task
	session, err := s.browserService.CreateSession(ctx, task.UserID, browser.SessionCreateRequest{
		SessionName: fmt.Sprintf("Task %s", task.ID.String()[:8]),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create browser session: %w", err)
	}

	// Navigate to URL
	navReq := browser.NavigateRequest{
		URL:             input.URL,
		WaitForSelector: input.WaitForSelector,
		Timeout:         input.Timeout,
		Headers:         input.Headers,
	}

	response, err := s.browserService.Navigate(ctx, session.ID, navReq)
	if err != nil {
		return nil, fmt.Errorf("navigation failed: %w", err)
	}

	return map[string]interface{}{
		"navigation_response": response,
		"session_id":          session.ID.String(),
	}, nil
}

// executeExtractTask executes a content extraction task
func (s *Service) executeExtractTask(ctx context.Context, task *Task) (map[string]interface{}, error) {
	var input ExtractTaskInput
	if err := s.parseTaskInput(task.InputData, &input); err != nil {
		return nil, fmt.Errorf("invalid extract task input: %w", err)
	}

	// Create browser session
	session, err := s.browserService.CreateSession(ctx, task.UserID, browser.SessionCreateRequest{
		SessionName: fmt.Sprintf("Extract Task %s", task.ID.String()[:8]),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create browser session: %w", err)
	}

	// Navigate to URL if provided
	if input.URL != "" {
		navReq := browser.NavigateRequest{URL: input.URL}
		_, err := s.browserService.Navigate(ctx, session.ID, navReq)
		if err != nil {
			return nil, fmt.Errorf("navigation failed: %w", err)
		}
	}

	// Extract content
	extractReq := browser.ExtractRequest{
		Selectors: input.Selectors,
		DataType:  input.DataType,
		Schema:    input.Schema,
	}

	response, err := s.browserService.Extract(ctx, session.ID, extractReq)
	if err != nil {
		return nil, fmt.Errorf("extraction failed: %w", err)
	}

	return map[string]interface{}{
		"extraction_response": response,
		"session_id":          session.ID.String(),
	}, nil
}

// executeInteractTask executes a page interaction task
func (s *Service) executeInteractTask(ctx context.Context, task *Task) (map[string]interface{}, error) {
	var input InteractTaskInput
	if err := s.parseTaskInput(task.InputData, &input); err != nil {
		return nil, fmt.Errorf("invalid interact task input: %w", err)
	}

	// Create browser session
	session, err := s.browserService.CreateSession(ctx, task.UserID, browser.SessionCreateRequest{
		SessionName: fmt.Sprintf("Interact Task %s", task.ID.String()[:8]),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create browser session: %w", err)
	}

	// Navigate to URL if provided
	if input.URL != "" {
		navReq := browser.NavigateRequest{URL: input.URL}
		_, err := s.browserService.Navigate(ctx, session.ID, navReq)
		if err != nil {
			return nil, fmt.Errorf("navigation failed: %w", err)
		}
	}

	// Convert AI actions to browser actions
	var browserActions []browser.Action
	for _, action := range input.Actions {
		browserAction := browser.Action{
			Type:     browser.ActionType(action.Type),
			Selector: action.Selector,
			Value:    action.Value,
			Options:  action.Options,
		}
		browserActions = append(browserActions, browserAction)
	}

	// Perform interactions
	interactReq := browser.InteractRequest{
		Actions:     browserActions,
		WaitBetween: input.WaitBetween,
		Screenshot:  true,
	}

	response, err := s.browserService.Interact(ctx, session.ID, interactReq)
	if err != nil {
		return nil, fmt.Errorf("interaction failed: %w", err)
	}

	return map[string]interface{}{
		"interaction_response": response,
		"session_id":           session.ID.String(),
	}, nil
}

// executeSummarizeTask executes a content summarization task
func (s *Service) executeSummarizeTask(ctx context.Context, task *Task) (map[string]interface{}, error) {
	var input SummarizeTaskInput
	if err := s.parseTaskInput(task.InputData, &input); err != nil {
		return nil, fmt.Errorf("invalid summarize task input: %w", err)
	}

	content := input.Content

	// If URL is provided, extract content first
	if input.URL != "" && content == "" {
		session, err := s.browserService.CreateSession(ctx, task.UserID, browser.SessionCreateRequest{})
		if err != nil {
			return nil, fmt.Errorf("failed to create browser session: %w", err)
		}

		// Navigate and extract content
		navReq := browser.NavigateRequest{URL: input.URL}
		_, err = s.browserService.Navigate(ctx, session.ID, navReq)
		if err != nil {
			return nil, fmt.Errorf("navigation failed: %w", err)
		}

		extractReq := browser.ExtractRequest{DataType: "text"}
		extractResp, err := s.browserService.Extract(ctx, session.ID, extractReq)
		if err != nil {
			return nil, fmt.Errorf("content extraction failed: %w", err)
		}

		if textData, ok := extractResp.Data["text"].(map[string]interface{}); ok {
			if bodyText, ok := textData["body"].(string); ok {
				content = bodyText
			}
		}
	}

	if content == "" {
		return nil, fmt.Errorf("no content to summarize")
	}

	// Summarize content
	options := SummarizeOptions{
		Length: input.Length,
		Focus:  input.Focus,
	}

	summary, err := s.aiProvider.SummarizeContent(ctx, content, options)
	if err != nil {
		return nil, fmt.Errorf("summarization failed: %w", err)
	}

	return map[string]interface{}{
		"summary":        summary,
		"options":        options,
		"content_length": len(content),
	}, nil
}

// executeScreenshotTask executes a screenshot task
func (s *Service) executeScreenshotTask(ctx context.Context, task *Task) (map[string]interface{}, error) {
	// Create browser session
	session, err := s.browserService.CreateSession(ctx, task.UserID, browser.SessionCreateRequest{
		SessionName: fmt.Sprintf("Screenshot Task %s", task.ID.String()[:8]),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create browser session: %w", err)
	}

	// Take screenshot
	screenshotReq := browser.ScreenshotRequest{
		FullPage: true,
		Quality:  90,
		Format:   "png",
	}

	response, err := s.browserService.TakeScreenshot(ctx, session.ID, screenshotReq)
	if err != nil {
		return nil, fmt.Errorf("screenshot failed: %w", err)
	}

	return map[string]interface{}{
		"screenshot_response": response,
		"session_id":          session.ID.String(),
	}, nil
}

// executeAnalyzeTask executes a content analysis task
func (s *Service) executeAnalyzeTask(ctx context.Context, task *Task) (map[string]interface{}, error) {
	var input AnalyzeTaskInput
	if err := s.parseTaskInput(task.InputData, &input); err != nil {
		return nil, fmt.Errorf("invalid analyze task input: %w", err)
	}

	content := input.Content

	// If URL is provided, extract content first
	if input.URL != "" && content == "" {
		session, err := s.browserService.CreateSession(ctx, task.UserID, browser.SessionCreateRequest{})
		if err != nil {
			return nil, fmt.Errorf("failed to create browser session: %w", err)
		}

		navReq := browser.NavigateRequest{URL: input.URL}
		_, err = s.browserService.Navigate(ctx, session.ID, navReq)
		if err != nil {
			return nil, fmt.Errorf("navigation failed: %w", err)
		}

		extractReq := browser.ExtractRequest{DataType: "text"}
		extractResp, err := s.browserService.Extract(ctx, session.ID, extractReq)
		if err != nil {
			return nil, fmt.Errorf("content extraction failed: %w", err)
		}

		if textData, ok := extractResp.Data["text"].(map[string]interface{}); ok {
			if bodyText, ok := textData["body"].(string); ok {
				content = bodyText
			}
		}
	}

	if content == "" {
		return nil, fmt.Errorf("no content to analyze")
	}

	// Analyze content
	analysisType := input.AnalysisType
	if analysisType == "" {
		analysisType = "general"
	}

	analysis, err := s.aiProvider.AnalyzeContent(ctx, content, analysisType)
	if err != nil {
		return nil, fmt.Errorf("analysis failed: %w", err)
	}

	return map[string]interface{}{
		"analysis":       analysis,
		"analysis_type":  analysisType,
		"content_length": len(content),
		"criteria":       input.Criteria,
	}, nil
}

// Health check methods

// GetHealthStatus returns the health status of all AI providers
func (s *Service) GetHealthStatus() *HealthCheckResponse {
	return s.healthMonitor.GetHealthCheckResponse()
}

// GetProviderHealth returns the health status of a specific provider
func (s *Service) GetProviderHealth(provider string) (*HealthStatus, bool) {
	return s.healthMonitor.GetStatus(provider)
}

// CheckProviderHealth performs an immediate health check on a provider
func (s *Service) CheckProviderHealth(ctx context.Context, provider string) error {
	return s.healthMonitor.CheckProviderNow(ctx, provider)
}

// IsProviderHealthy checks if a provider is currently healthy
func (s *Service) IsProviderHealthy(provider string) bool {
	return s.healthMonitor.IsProviderHealthy(provider)
}

// GetProviderModels returns available models for a provider
func (s *Service) GetProviderModels(provider string) ([]string, error) {
	return s.healthMonitor.GetProviderModels(provider)
}

// StopHealthMonitoring stops the health monitoring
func (s *Service) StopHealthMonitoring() {
	if s.healthMonitor != nil {
		s.healthMonitor.Stop()
	}
}

// Helper methods

// getOrCreateConversation gets an existing conversation or creates a new one
func (s *Service) getOrCreateConversation(ctx context.Context, userID uuid.UUID, conversationID *uuid.UUID, sessionID *uuid.UUID) (*Conversation, error) {
	if conversationID != nil {
		return s.getConversation(ctx, *conversationID)
	}

	// Create new conversation
	conversation := &Conversation{
		ID:        uuid.New(),
		UserID:    userID,
		SessionID: sessionID,
		Title:     "New Conversation",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := `
		INSERT INTO ai_conversations (id, user_id, session_id, title, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := s.db.ExecContext(ctx, query, conversation.ID, conversation.UserID, conversation.SessionID, conversation.Title, conversation.CreatedAt, conversation.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	return conversation, nil
}

// getConversation retrieves a conversation by ID
func (s *Service) getConversation(ctx context.Context, conversationID uuid.UUID) (*Conversation, error) {
	query := `
		SELECT id, user_id, session_id, title, created_at, updated_at
		FROM ai_conversations WHERE id = $1
	`
	conversation := &Conversation{}
	err := s.db.QueryRowContext(ctx, query, conversationID).Scan(
		&conversation.ID, &conversation.UserID, &conversation.SessionID,
		&conversation.Title, &conversation.CreatedAt, &conversation.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("conversation not found: %w", err)
	}
	return conversation, nil
}

// saveMessage saves a message to the database
func (s *Service) saveMessage(ctx context.Context, message *Message) error {
	metadataJSON, _ := json.Marshal(message.Metadata)
	query := `
		INSERT INTO ai_messages (id, conversation_id, role, content, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := s.db.ExecContext(ctx, query, message.ID, message.ConversationID, message.Role, message.Content, metadataJSON, message.CreatedAt)
	return err
}

// getConversationMessages retrieves all messages for a conversation
func (s *Service) getConversationMessages(ctx context.Context, conversationID uuid.UUID) ([]Message, error) {
	query := `
		SELECT id, conversation_id, role, content, metadata, created_at
		FROM ai_messages WHERE conversation_id = $1 ORDER BY created_at ASC
	`
	rows, err := s.db.QueryContext(ctx, query, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		var metadataJSON []byte
		err := rows.Scan(&message.ID, &message.ConversationID, &message.Role, &message.Content, &metadataJSON, &message.CreatedAt)
		if err != nil {
			return nil, err
		}
		if len(metadataJSON) > 0 {
			json.Unmarshal(metadataJSON, &message.Metadata)
		}
		messages = append(messages, message)
	}
	return messages, nil
}

// saveTask saves a task to the database
func (s *Service) saveTask(ctx context.Context, task *Task) error {
	inputDataJSON, _ := json.Marshal(task.InputData)
	outputDataJSON, _ := json.Marshal(task.OutputData)

	query := `
		INSERT INTO ai_tasks (id, conversation_id, user_id, task_type, description, status, input_data, output_data, error_message, started_at, completed_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := s.db.ExecContext(ctx, query, task.ID, task.ConversationID, task.UserID, task.TaskType, task.Description, task.Status, inputDataJSON, outputDataJSON, task.ErrorMessage, task.StartedAt, task.CompletedAt, task.CreatedAt)
	return err
}

// updateTaskStatus updates a task's status and related fields
func (s *Service) updateTaskStatus(ctx context.Context, task *Task) error {
	outputDataJSON, _ := json.Marshal(task.OutputData)

	query := `
		UPDATE ai_tasks
		SET status = $1, output_data = $2, error_message = $3, started_at = $4, completed_at = $5
		WHERE id = $6
	`
	_, err := s.db.ExecContext(ctx, query, task.Status, outputDataJSON, task.ErrorMessage, task.StartedAt, task.CompletedAt, task.ID)
	return err
}

// parseTaskInput parses task input data into a specific struct
func (s *Service) parseTaskInput(inputData map[string]interface{}, target interface{}) error {
	jsonData, err := json.Marshal(inputData)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, target)
}

// processFunctionCall processes a function call from the AI and creates a task
func (s *Service) processFunctionCall(ctx context.Context, userID, conversationID uuid.UUID, functionCall interface{}) (*Task, error) {
	// This would parse the function call and create appropriate tasks
	// For now, return nil to indicate no task was created
	return nil, nil
}

// generateSuggestions generates conversation suggestions
func (s *Service) generateSuggestions(ctx context.Context, userMessage, aiResponse string) []string {
	// Simple suggestions based on common patterns
	suggestions := []string{
		"Take a screenshot of this page",
		"Extract all links from this page",
		"Summarize the main content",
		"Navigate to another page",
	}
	return suggestions
}
