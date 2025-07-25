package ai

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Conversation represents an AI conversation
type Conversation struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	SessionID *uuid.UUID `json:"session_id,omitempty" db:"session_id"`
	Title     string     `json:"title" db:"title"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	Messages  []Message  `json:"messages,omitempty"`
}

// Message represents a message in a conversation
type Message struct {
	ID             uuid.UUID              `json:"id" db:"id"`
	ConversationID uuid.UUID              `json:"conversation_id" db:"conversation_id"`
	Role           MessageRole            `json:"role" db:"role"`
	Content        string                 `json:"content" db:"content"`
	Metadata       map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	CreatedAt      time.Time              `json:"created_at" db:"created_at"`
}

// MessageRole represents the role of a message sender
type MessageRole string

const (
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
	RoleSystem    MessageRole = "system"
)

// Task represents an AI task
type Task struct {
	ID             uuid.UUID              `json:"id" db:"id"`
	ConversationID uuid.UUID              `json:"conversation_id" db:"conversation_id"`
	UserID         uuid.UUID              `json:"user_id" db:"user_id"`
	TaskType       TaskType               `json:"task_type" db:"task_type"`
	Description    string                 `json:"description" db:"description"`
	Status         TaskStatus             `json:"status" db:"status"`
	InputData      map[string]interface{} `json:"input_data,omitempty" db:"input_data"`
	OutputData     map[string]interface{} `json:"output_data,omitempty" db:"output_data"`
	ErrorMessage   *string                `json:"error_message,omitempty" db:"error_message"`
	StartedAt      *time.Time             `json:"started_at,omitempty" db:"started_at"`
	CompletedAt    *time.Time             `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt      time.Time              `json:"created_at" db:"created_at"`
}

// TaskType represents the type of AI task
type TaskType string

const (
	TaskTypeNavigate   TaskType = "navigate"
	TaskTypeExtract    TaskType = "extract"
	TaskTypeInteract   TaskType = "interact"
	TaskTypeSummarize  TaskType = "summarize"
	TaskTypeSearch     TaskType = "search"
	TaskTypeFillForm   TaskType = "fill_form"
	TaskTypeScreenshot TaskType = "screenshot"
	TaskTypeAnalyze    TaskType = "analyze"
	TaskTypeCustom     TaskType = "custom"
)

// TaskStatus represents the status of a task
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
	TaskStatusCancelled TaskStatus = "cancelled"
)

// ChatRequest represents a chat message request
type ChatRequest struct {
	ConversationID *uuid.UUID `json:"conversation_id,omitempty"`
	Message        string     `json:"message" validate:"required"`
	SessionID      *uuid.UUID `json:"session_id,omitempty"`
}

// ChatResponse represents a chat response
type ChatResponse struct {
	ConversationID uuid.UUID `json:"conversation_id"`
	Message        Message   `json:"message"`
	Tasks          []Task    `json:"tasks,omitempty"`
	Suggestions    []string  `json:"suggestions,omitempty"`
}

// TaskRequest represents a task creation request
type TaskRequest struct {
	ConversationID *uuid.UUID             `json:"conversation_id,omitempty"`
	TaskType       TaskType               `json:"task_type" validate:"required"`
	Description    string                 `json:"description" validate:"required"`
	InputData      map[string]interface{} `json:"input_data,omitempty"`
}

// TaskResponse represents a task response
type TaskResponse struct {
	Task Task `json:"task"`
}

// NavigateTaskInput represents input for navigation tasks
type NavigateTaskInput struct {
	URL             string            `json:"url" validate:"required"`
	WaitForSelector string            `json:"wait_for_selector,omitempty"`
	Timeout         int               `json:"timeout,omitempty"`
	Headers         map[string]string `json:"headers,omitempty"`
}

// ExtractTaskInput represents input for content extraction tasks
type ExtractTaskInput struct {
	URL       string   `json:"url,omitempty"`
	Selectors []string `json:"selectors,omitempty"`
	DataType  string   `json:"data_type,omitempty"` // text, links, images, tables
	Schema    string   `json:"schema,omitempty"`    // JSON schema for structured extraction
}

// InteractTaskInput represents input for page interaction tasks
type InteractTaskInput struct {
	URL         string              `json:"url,omitempty"`
	Actions     []InteractionAction `json:"actions" validate:"required"`
	WaitBetween int                 `json:"wait_between,omitempty"` // milliseconds
}

// InteractionAction represents a single interaction action
type InteractionAction struct {
	Type     ActionType             `json:"type" validate:"required"`
	Selector string                 `json:"selector,omitempty"`
	Value    string                 `json:"value,omitempty"`
	Options  map[string]interface{} `json:"options,omitempty"`
}

// ActionType represents the type of interaction action
type ActionType string

const (
	ActionClick      ActionType = "click"
	ActionType_      ActionType = "type"
	ActionSelect     ActionType = "select"
	ActionScroll     ActionType = "scroll"
	ActionWait       ActionType = "wait"
	ActionHover      ActionType = "hover"
	ActionScreenshot ActionType = "screenshot"
	ActionKeyPress   ActionType = "key_press"
)

// SummarizeTaskInput represents input for content summarization tasks
type SummarizeTaskInput struct {
	URL     string `json:"url,omitempty"`
	Content string `json:"content,omitempty"`
	Length  string `json:"length,omitempty"` // short, medium, long
	Focus   string `json:"focus,omitempty"`  // main_points, technical, business
}

// SearchTaskInput represents input for search tasks
type SearchTaskInput struct {
	Query      string   `json:"query" validate:"required"`
	SearchType string   `json:"search_type,omitempty"` // web, images, news
	MaxResults int      `json:"max_results,omitempty"`
	Filters    []string `json:"filters,omitempty"`
}

// FillFormTaskInput represents input for form filling tasks
type FillFormTaskInput struct {
	URL       string                 `json:"url,omitempty"`
	FormData  map[string]interface{} `json:"form_data" validate:"required"`
	Submit    bool                   `json:"submit,omitempty"`
	Selectors map[string]string      `json:"selectors,omitempty"` // field_name -> selector mapping
}

// AnalyzeTaskInput represents input for content analysis tasks
type AnalyzeTaskInput struct {
	URL          string   `json:"url,omitempty"`
	Content      string   `json:"content,omitempty"`
	AnalysisType string   `json:"analysis_type,omitempty"` // sentiment, keywords, entities, structure
	Criteria     []string `json:"criteria,omitempty"`
}

// TaskOutput represents generic task output
type TaskOutput struct {
	Success     bool                   `json:"success"`
	Data        map[string]interface{} `json:"data,omitempty"`
	Screenshots []string               `json:"screenshots,omitempty"`
	Errors      []string               `json:"errors,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// AIProvider represents an AI service provider
type AIProvider interface {
	GenerateResponse(ctx context.Context, messages []Message) (*Message, error)
	AnalyzeContent(ctx context.Context, content string, analysisType string) (map[string]interface{}, error)
	ExtractStructuredData(ctx context.Context, content string, schema string) (map[string]interface{}, error)
	SummarizeContent(ctx context.Context, content string, options SummarizeOptions) (string, error)
}

// SummarizeOptions represents options for content summarization
type SummarizeOptions struct {
	Length string
	Focus  string
	Style  string
}

// ConversationListRequest represents a request to list conversations
type ConversationListRequest struct {
	UserID uuid.UUID `json:"user_id"`
	Limit  int       `json:"limit,omitempty"`
	Offset int       `json:"offset,omitempty"`
}

// ConversationListResponse represents a response with conversation list
type ConversationListResponse struct {
	Conversations []Conversation `json:"conversations"`
	Total         int            `json:"total"`
	HasMore       bool           `json:"has_more"`
}

// TaskListRequest represents a request to list tasks
type TaskListRequest struct {
	UserID         uuid.UUID   `json:"user_id"`
	ConversationID *uuid.UUID  `json:"conversation_id,omitempty"`
	Status         *TaskStatus `json:"status,omitempty"`
	TaskType       *TaskType   `json:"task_type,omitempty"`
	Limit          int         `json:"limit,omitempty"`
	Offset         int         `json:"offset,omitempty"`
}

// TaskListResponse represents a response with task list
type TaskListResponse struct {
	Tasks   []Task `json:"tasks"`
	Total   int    `json:"total"`
	HasMore bool   `json:"has_more"`
}
