package marketplace

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ai-agentic-browser/internal/ai"
	"github.com/ai-agentic-browser/internal/browser"
	"github.com/ai-agentic-browser/internal/web3"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

// AIStepExecutor executes AI-related workflow steps
type AIStepExecutor struct {
	aiService *ai.Service
	logger    *observability.Logger
}

func NewAIStepExecutor(aiService *ai.Service, logger *observability.Logger) *AIStepExecutor {
	return &AIStepExecutor{
		aiService: aiService,
		logger:    logger,
	}
}

func (e *AIStepExecutor) Execute(ctx context.Context, step WorkflowStep, input map[string]interface{}) (map[string]interface{}, error) {
	switch step.Action {
	case "chat":
		return e.executeChat(ctx, step, input)
	case "analyze":
		return e.executeAnalyze(ctx, step, input)
	case "summarize":
		return e.executeSummarize(ctx, step, input)
	case "extract":
		return e.executeExtract(ctx, step, input)
	default:
		return nil, fmt.Errorf("unsupported AI action: %s", step.Action)
	}
}

func (e *AIStepExecutor) Validate(step WorkflowStep) error {
	validActions := []string{"chat", "analyze", "summarize", "extract"}
	for _, action := range validActions {
		if step.Action == action {
			return nil
		}
	}
	return fmt.Errorf("invalid AI action: %s", step.Action)
}

func (e *AIStepExecutor) executeChat(ctx context.Context, step WorkflowStep, input map[string]interface{}) (map[string]interface{}, error) {
	message, ok := input["message"].(string)
	if !ok {
		return nil, fmt.Errorf("message parameter is required for chat action")
	}

	// Get user ID from input
	userIDStr, ok := input["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("user_id parameter is required for chat action")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	// Create AI chat request
	chatReq := ai.ChatRequest{
		Message: message,
	}

	// Set conversation ID if provided
	if conversationIDStr, ok := input["conversation_id"].(string); ok {
		conversationID, err := uuid.Parse(conversationIDStr)
		if err == nil {
			chatReq.ConversationID = &conversationID
		}
	}

	response, err := e.aiService.Chat(ctx, userID, chatReq)
	if err != nil {
		return nil, fmt.Errorf("AI chat failed: %w", err)
	}

	return map[string]interface{}{
		"conversation_id": response.ConversationID.String(),
		"message":         response.Message,
		"tasks":           response.Tasks,
		"suggestions":     response.Suggestions,
	}, nil
}

func (e *AIStepExecutor) executeAnalyze(ctx context.Context, step WorkflowStep, input map[string]interface{}) (map[string]interface{}, error) {
	// Get user ID from input
	userIDStr, ok := input["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("user_id parameter is required for analyze action")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	// Create analyze task request
	taskReq := ai.TaskRequest{
		TaskType:    ai.TaskTypeAnalyze,
		Description: fmt.Sprintf("Analyze content - %s", step.Name),
		InputData:   input,
	}

	// Set conversation ID if provided
	if conversationIDStr, ok := input["conversation_id"].(string); ok {
		conversationID, err := uuid.Parse(conversationIDStr)
		if err == nil {
			taskReq.ConversationID = &conversationID
		}
	}

	response, err := e.aiService.CreateTask(ctx, userID, taskReq)
	if err != nil {
		return nil, fmt.Errorf("content analysis task creation failed: %w", err)
	}

	return map[string]interface{}{
		"task_id":     response.Task.ID.String(),
		"task_status": response.Task.Status,
		"task_type":   response.Task.TaskType,
	}, nil
}

func (e *AIStepExecutor) executeSummarize(ctx context.Context, step WorkflowStep, input map[string]interface{}) (map[string]interface{}, error) {
	// Get user ID from input
	userIDStr, ok := input["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("user_id parameter is required for summarize action")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	// Create summarize task request
	taskReq := ai.TaskRequest{
		TaskType:    ai.TaskTypeSummarize,
		Description: fmt.Sprintf("Summarize content - %s", step.Name),
		InputData:   input,
	}

	// Set conversation ID if provided
	if conversationIDStr, ok := input["conversation_id"].(string); ok {
		conversationID, err := uuid.Parse(conversationIDStr)
		if err == nil {
			taskReq.ConversationID = &conversationID
		}
	}

	response, err := e.aiService.CreateTask(ctx, userID, taskReq)
	if err != nil {
		return nil, fmt.Errorf("content summarization task creation failed: %w", err)
	}

	return map[string]interface{}{
		"task_id":     response.Task.ID.String(),
		"task_status": response.Task.Status,
		"task_type":   response.Task.TaskType,
	}, nil
}

func (e *AIStepExecutor) executeExtract(ctx context.Context, step WorkflowStep, input map[string]interface{}) (map[string]interface{}, error) {
	// Get user ID from input
	userIDStr, ok := input["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("user_id parameter is required for extract action")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	// Create extract task request
	taskReq := ai.TaskRequest{
		TaskType:    ai.TaskTypeExtract,
		Description: fmt.Sprintf("Extract data - %s", step.Name),
		InputData:   input,
	}

	// Set conversation ID if provided
	if conversationIDStr, ok := input["conversation_id"].(string); ok {
		conversationID, err := uuid.Parse(conversationIDStr)
		if err == nil {
			taskReq.ConversationID = &conversationID
		}
	}

	response, err := e.aiService.CreateTask(ctx, userID, taskReq)
	if err != nil {
		return nil, fmt.Errorf("data extraction task creation failed: %w", err)
	}

	return map[string]interface{}{
		"task_id":     response.Task.ID.String(),
		"task_status": response.Task.Status,
		"task_type":   response.Task.TaskType,
	}, nil
}

// BrowserStepExecutor executes browser-related workflow steps
type BrowserStepExecutor struct {
	browserService *browser.Service
	logger         *observability.Logger
}

func NewBrowserStepExecutor(browserService *browser.Service, logger *observability.Logger) *BrowserStepExecutor {
	return &BrowserStepExecutor{
		browserService: browserService,
		logger:         logger,
	}
}

func (e *BrowserStepExecutor) Execute(ctx context.Context, step WorkflowStep, input map[string]interface{}) (map[string]interface{}, error) {
	switch step.Action {
	case "navigate":
		return e.executeNavigate(ctx, step, input)
	case "interact":
		return e.executeInteract(ctx, step, input)
	case "extract":
		return e.executeExtract(ctx, step, input)
	case "screenshot":
		return e.executeScreenshot(ctx, step, input)
	default:
		return nil, fmt.Errorf("unsupported browser action: %s", step.Action)
	}
}

func (e *BrowserStepExecutor) Validate(step WorkflowStep) error {
	validActions := []string{"navigate", "interact", "extract", "screenshot"}
	for _, action := range validActions {
		if step.Action == action {
			return nil
		}
	}
	return fmt.Errorf("invalid browser action: %s", step.Action)
}

func (e *BrowserStepExecutor) executeNavigate(ctx context.Context, step WorkflowStep, input map[string]interface{}) (map[string]interface{}, error) {
	url, ok := input["url"].(string)
	if !ok {
		return nil, fmt.Errorf("url parameter is required for navigate action")
	}

	// Get session ID from input
	sessionIDStr, ok := input["session_id"].(string)
	if !ok {
		return nil, fmt.Errorf("session_id parameter is required for browser actions")
	}

	// Parse session ID
	sessionID, err := parseUUID(sessionIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid session_id: %w", err)
	}

	// Create navigate request
	navReq := browser.NavigateRequest{
		URL:     url,
		Timeout: 30,
	}

	if waitFor, ok := input["wait_for"].(string); ok {
		navReq.WaitForSelector = waitFor
	}

	response, err := e.browserService.Navigate(ctx, sessionID, navReq)
	if err != nil {
		return nil, fmt.Errorf("navigation failed: %w", err)
	}

	return map[string]interface{}{
		"success":   response.Success,
		"url":       response.URL,
		"title":     response.Title,
		"load_time": response.LoadTime,
	}, nil
}

func (e *BrowserStepExecutor) executeInteract(ctx context.Context, step WorkflowStep, input map[string]interface{}) (map[string]interface{}, error) {
	sessionIDStr, ok := input["session_id"].(string)
	if !ok {
		return nil, fmt.Errorf("session_id parameter is required for browser actions")
	}

	sessionID, err := parseUUID(sessionIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid session_id: %w", err)
	}

	// Parse actions
	actionsData, ok := input["actions"]
	if !ok {
		return nil, fmt.Errorf("actions parameter is required for interact action")
	}

	var actions []browser.Action
	if actionsJSON, err := json.Marshal(actionsData); err == nil {
		json.Unmarshal(actionsJSON, &actions)
	}

	interactReq := browser.InteractRequest{
		Actions:    actions,
		Screenshot: true,
	}

	response, err := e.browserService.Interact(ctx, sessionID, interactReq)
	if err != nil {
		return nil, fmt.Errorf("interaction failed: %w", err)
	}

	return map[string]interface{}{
		"success":     response.Success,
		"screenshots": response.Screenshots,
		"results":     response.Results,
	}, nil
}

func (e *BrowserStepExecutor) executeExtract(ctx context.Context, step WorkflowStep, input map[string]interface{}) (map[string]interface{}, error) {
	sessionIDStr, ok := input["session_id"].(string)
	if !ok {
		return nil, fmt.Errorf("session_id parameter is required for browser actions")
	}

	sessionID, err := parseUUID(sessionIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid session_id: %w", err)
	}

	extractReq := browser.ExtractRequest{
		DataType: "text",
	}

	if dataType, ok := input["data_type"].(string); ok {
		extractReq.DataType = dataType
	}

	if selectorsData, ok := input["selectors"]; ok {
		if selectorsJSON, err := json.Marshal(selectorsData); err == nil {
			json.Unmarshal(selectorsJSON, &extractReq.Selectors)
		}
	}

	response, err := e.browserService.Extract(ctx, sessionID, extractReq)
	if err != nil {
		return nil, fmt.Errorf("extraction failed: %w", err)
	}

	return map[string]interface{}{
		"success":  response.Success,
		"data":     response.Data,
		"metadata": response.Metadata,
	}, nil
}

func (e *BrowserStepExecutor) executeScreenshot(ctx context.Context, step WorkflowStep, input map[string]interface{}) (map[string]interface{}, error) {
	sessionIDStr, ok := input["session_id"].(string)
	if !ok {
		return nil, fmt.Errorf("session_id parameter is required for browser actions")
	}

	sessionID, err := parseUUID(sessionIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid session_id: %w", err)
	}

	screenshotReq := browser.ScreenshotRequest{
		FullPage: true,
		Quality:  90,
		Format:   "png",
	}

	if fullPage, ok := input["full_page"].(bool); ok {
		screenshotReq.FullPage = fullPage
	}

	response, err := e.browserService.TakeScreenshot(ctx, sessionID, screenshotReq)
	if err != nil {
		return nil, fmt.Errorf("screenshot failed: %w", err)
	}

	return map[string]interface{}{
		"success":    response.Success,
		"screenshot": response.Screenshot,
		"format":     response.Format,
		"size":       response.Size,
	}, nil
}

// APIStepExecutor executes API-related workflow steps
type APIStepExecutor struct {
	logger *observability.Logger
}

func NewAPIStepExecutor(logger *observability.Logger) *APIStepExecutor {
	return &APIStepExecutor{logger: logger}
}

func (e *APIStepExecutor) Execute(ctx context.Context, step WorkflowStep, input map[string]interface{}) (map[string]interface{}, error) {
	switch step.Action {
	case "http_request":
		return e.executeHTTPRequest(ctx, step, input)
	default:
		return nil, fmt.Errorf("unsupported API action: %s", step.Action)
	}
}

func (e *APIStepExecutor) Validate(step WorkflowStep) error {
	if step.Action == "http_request" {
		return nil
	}
	return fmt.Errorf("invalid API action: %s", step.Action)
}

func (e *APIStepExecutor) executeHTTPRequest(ctx context.Context, step WorkflowStep, input map[string]interface{}) (map[string]interface{}, error) {
	url, ok := input["url"].(string)
	if !ok {
		return nil, fmt.Errorf("url parameter is required for http_request action")
	}

	method, _ := input["method"].(string)
	if method == "" {
		method = "GET"
	}

	// Create HTTP request
	var body io.Reader
	if bodyData, ok := input["body"]; ok {
		if bodyJSON, err := json.Marshal(bodyData); err == nil {
			body = strings.NewReader(string(bodyJSON))
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	if headersData, ok := input["headers"].(map[string]interface{}); ok {
		for key, value := range headersData {
			if valueStr, ok := value.(string); ok {
				req.Header.Set(key, valueStr)
			}
		}
	}

	// Set default content type for POST/PUT
	if (method == "POST" || method == "PUT") && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Execute request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse JSON response if possible
	var responseData interface{}
	if err := json.Unmarshal(responseBody, &responseData); err != nil {
		responseData = string(responseBody)
	}

	return map[string]interface{}{
		"status_code": resp.StatusCode,
		"headers":     resp.Header,
		"body":        responseData,
		"success":     resp.StatusCode >= 200 && resp.StatusCode < 300,
	}, nil
}

// Additional step executors (DataStepExecutor, LogicStepExecutor, etc.) would follow similar patterns

// DataStepExecutor executes data manipulation steps
type DataStepExecutor struct {
	logger *observability.Logger
}

func NewDataStepExecutor(logger *observability.Logger) *DataStepExecutor {
	return &DataStepExecutor{logger: logger}
}

func (e *DataStepExecutor) Execute(ctx context.Context, step WorkflowStep, input map[string]interface{}) (map[string]interface{}, error) {
	switch step.Action {
	case "transform":
		return e.executeTransform(ctx, step, input)
	case "filter":
		return e.executeFilter(ctx, step, input)
	case "aggregate":
		return e.executeAggregate(ctx, step, input)
	default:
		return nil, fmt.Errorf("unsupported data action: %s", step.Action)
	}
}

func (e *DataStepExecutor) Validate(step WorkflowStep) error {
	validActions := []string{"transform", "filter", "aggregate"}
	for _, action := range validActions {
		if step.Action == action {
			return nil
		}
	}
	return fmt.Errorf("invalid data action: %s", step.Action)
}

func (e *DataStepExecutor) executeTransform(ctx context.Context, step WorkflowStep, input map[string]interface{}) (map[string]interface{}, error) {
	// Implement data transformation logic
	return map[string]interface{}{"transformed": true}, nil
}

func (e *DataStepExecutor) executeFilter(ctx context.Context, step WorkflowStep, input map[string]interface{}) (map[string]interface{}, error) {
	// Implement data filtering logic
	return map[string]interface{}{"filtered": true}, nil
}

func (e *DataStepExecutor) executeAggregate(ctx context.Context, step WorkflowStep, input map[string]interface{}) (map[string]interface{}, error) {
	// Implement data aggregation logic
	return map[string]interface{}{"aggregated": true}, nil
}

// LogicStepExecutor executes logic and control flow steps
type LogicStepExecutor struct {
	logger *observability.Logger
}

func NewLogicStepExecutor(logger *observability.Logger) *LogicStepExecutor {
	return &LogicStepExecutor{logger: logger}
}

func (e *LogicStepExecutor) Execute(ctx context.Context, step WorkflowStep, input map[string]interface{}) (map[string]interface{}, error) {
	switch step.Action {
	case "condition":
		return e.executeCondition(ctx, step, input)
	case "loop":
		return e.executeLoop(ctx, step, input)
	default:
		return nil, fmt.Errorf("unsupported logic action: %s", step.Action)
	}
}

func (e *LogicStepExecutor) Validate(step WorkflowStep) error {
	validActions := []string{"condition", "loop"}
	for _, action := range validActions {
		if step.Action == action {
			return nil
		}
	}
	return fmt.Errorf("invalid logic action: %s", step.Action)
}

func (e *LogicStepExecutor) executeCondition(ctx context.Context, step WorkflowStep, input map[string]interface{}) (map[string]interface{}, error) {
	// Implement conditional logic
	return map[string]interface{}{"condition_result": true}, nil
}

func (e *LogicStepExecutor) executeLoop(ctx context.Context, step WorkflowStep, input map[string]interface{}) (map[string]interface{}, error) {
	// Implement loop logic
	return map[string]interface{}{"loop_completed": true}, nil
}

// NotifyStepExecutor and WaitStepExecutor would follow similar patterns
type NotifyStepExecutor struct {
	logger *observability.Logger
}

func NewNotifyStepExecutor(logger *observability.Logger) *NotifyStepExecutor {
	return &NotifyStepExecutor{logger: logger}
}

func (e *NotifyStepExecutor) Execute(ctx context.Context, step WorkflowStep, input map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{"notification_sent": true}, nil
}

func (e *NotifyStepExecutor) Validate(step WorkflowStep) error {
	return nil
}

type WaitStepExecutor struct {
	logger *observability.Logger
}

func NewWaitStepExecutor(logger *observability.Logger) *WaitStepExecutor {
	return &WaitStepExecutor{logger: logger}
}

func (e *WaitStepExecutor) Execute(ctx context.Context, step WorkflowStep, input map[string]interface{}) (map[string]interface{}, error) {
	duration, _ := input["duration"].(float64)
	if duration > 0 {
		time.Sleep(time.Duration(duration) * time.Second)
	}
	return map[string]interface{}{"wait_completed": true}, nil
}

func (e *WaitStepExecutor) Validate(step WorkflowStep) error {
	return nil
}

// Web3StepExecutor executes Web3-related workflow steps
type Web3StepExecutor struct {
	web3Service *web3.Service
	logger      *observability.Logger
}

func NewWeb3StepExecutor(web3Service *web3.Service, logger *observability.Logger) *Web3StepExecutor {
	return &Web3StepExecutor{
		web3Service: web3Service,
		logger:      logger,
	}
}

func (e *Web3StepExecutor) Execute(ctx context.Context, step WorkflowStep, input map[string]interface{}) (map[string]interface{}, error) {
	switch step.Action {
	case "get_balance":
		return e.executeGetBalance(ctx, step, input)
	case "send_transaction":
		return e.executeSendTransaction(ctx, step, input)
	case "defi_interact":
		return e.executeDeFiInteract(ctx, step, input)
	default:
		return nil, fmt.Errorf("unsupported Web3 action: %s", step.Action)
	}
}

func (e *Web3StepExecutor) Validate(step WorkflowStep) error {
	validActions := []string{"get_balance", "send_transaction", "defi_interact"}
	for _, action := range validActions {
		if step.Action == action {
			return nil
		}
	}
	return fmt.Errorf("invalid Web3 action: %s", step.Action)
}

func (e *Web3StepExecutor) executeGetBalance(ctx context.Context, step WorkflowStep, input map[string]interface{}) (map[string]interface{}, error) {
	// Implementation would call web3Service.GetBalance
	return map[string]interface{}{"balance": "1.5", "currency": "ETH"}, nil
}

func (e *Web3StepExecutor) executeSendTransaction(ctx context.Context, step WorkflowStep, input map[string]interface{}) (map[string]interface{}, error) {
	// Implementation would call web3Service.CreateTransaction
	return map[string]interface{}{"tx_hash": "0x123...", "status": "pending"}, nil
}

func (e *Web3StepExecutor) executeDeFiInteract(ctx context.Context, step WorkflowStep, input map[string]interface{}) (map[string]interface{}, error) {
	// Implementation would call web3Service.InteractWithDeFiProtocol
	return map[string]interface{}{"success": true, "protocol": "uniswap"}, nil
}

// Helper function to parse UUID from string
func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// Helper function to convert interface{} to string
func toString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

// Helper function to convert interface{} to int
func toInt(v interface{}) int {
	switch val := v.(type) {
	case int:
		return val
	case float64:
		return int(val)
	case string:
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return 0
}
