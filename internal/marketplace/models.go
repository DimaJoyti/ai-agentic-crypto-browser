package marketplace

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// AgentTemplate represents a reusable AI agent template
type AgentTemplate struct {
	ID          uuid.UUID              `json:"id" db:"id"`
	Name        string                 `json:"name" db:"name"`
	Description string                 `json:"description" db:"description"`
	Category    string                 `json:"category" db:"category"`
	Tags        []string               `json:"tags" db:"tags"`
	AuthorID    uuid.UUID              `json:"author_id" db:"author_id"`
	AuthorName  string                 `json:"author_name" db:"author_name"`
	Version     string                 `json:"version" db:"version"`
	IsPublic    bool                   `json:"is_public" db:"is_public"`
	IsFeatured  bool                   `json:"is_featured" db:"is_featured"`
	Price       float64                `json:"price" db:"price"`
	Currency    string                 `json:"currency" db:"currency"`
	Rating      float64                `json:"rating" db:"rating"`
	Downloads   int                    `json:"downloads" db:"downloads"`
	Workflow    WorkflowDefinition     `json:"workflow" db:"workflow"`
	Metadata    map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
}

// WorkflowDefinition defines the structure of an automation workflow
type WorkflowDefinition struct {
	Version     string                 `json:"version"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Triggers    []WorkflowTrigger      `json:"triggers"`
	Steps       []WorkflowStep         `json:"steps"`
	Variables   map[string]interface{} `json:"variables"`
	Settings    WorkflowSettings       `json:"settings"`
}

// WorkflowTrigger defines when a workflow should be executed
type WorkflowTrigger struct {
	Type       TriggerType            `json:"type"`
	Conditions map[string]interface{} `json:"conditions"`
	Schedule   *ScheduleConfig        `json:"schedule,omitempty"`
}

// TriggerType represents different trigger types
type TriggerType string

const (
	TriggerTypeManual    TriggerType = "manual"
	TriggerTypeSchedule  TriggerType = "schedule"
	TriggerTypeWebhook   TriggerType = "webhook"
	TriggerTypeEvent     TriggerType = "event"
	TriggerTypeCondition TriggerType = "condition"
)

// ScheduleConfig defines scheduling parameters
type ScheduleConfig struct {
	Cron     string    `json:"cron,omitempty"`
	Interval string    `json:"interval,omitempty"`
	StartAt  time.Time `json:"start_at,omitempty"`
	EndAt    time.Time `json:"end_at,omitempty"`
}

// WorkflowStep represents a single step in a workflow
type WorkflowStep struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        StepType               `json:"type"`
	Action      string                 `json:"action"`
	Parameters  map[string]interface{} `json:"parameters"`
	Conditions  []StepCondition        `json:"conditions,omitempty"`
	OnSuccess   []string               `json:"on_success,omitempty"`
	OnFailure   []string               `json:"on_failure,omitempty"`
	Timeout     int                    `json:"timeout,omitempty"`
	Retries     int                    `json:"retries,omitempty"`
	Parallel    bool                   `json:"parallel,omitempty"`
}

// StepType represents different types of workflow steps
type StepType string

const (
	StepTypeAI       StepType = "ai"
	StepTypeBrowser  StepType = "browser"
	StepTypeWeb3     StepType = "web3"
	StepTypeAPI      StepType = "api"
	StepTypeData     StepType = "data"
	StepTypeLogic    StepType = "logic"
	StepTypeNotify   StepType = "notify"
	StepTypeWait     StepType = "wait"
)

// StepCondition defines conditional execution
type StepCondition struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

// WorkflowSettings contains workflow execution settings
type WorkflowSettings struct {
	MaxExecutionTime int                    `json:"max_execution_time"`
	MaxRetries       int                    `json:"max_retries"`
	ErrorHandling    string                 `json:"error_handling"`
	Notifications    NotificationSettings   `json:"notifications"`
	Variables        map[string]interface{} `json:"variables"`
}

// NotificationSettings defines notification preferences
type NotificationSettings struct {
	OnSuccess bool     `json:"on_success"`
	OnFailure bool     `json:"on_failure"`
	Channels  []string `json:"channels"`
	Webhook   string   `json:"webhook,omitempty"`
}

// WorkflowExecution represents a workflow execution instance
type WorkflowExecution struct {
	ID           uuid.UUID              `json:"id" db:"id"`
	WorkflowID   uuid.UUID              `json:"workflow_id" db:"workflow_id"`
	UserID       uuid.UUID              `json:"user_id" db:"user_id"`
	Status       ExecutionStatus        `json:"status" db:"status"`
	TriggerType  TriggerType            `json:"trigger_type" db:"trigger_type"`
	Input        map[string]interface{} `json:"input" db:"input"`
	Output       map[string]interface{} `json:"output" db:"output"`
	Steps        []StepExecution        `json:"steps" db:"steps"`
	StartedAt    time.Time              `json:"started_at" db:"started_at"`
	CompletedAt  *time.Time             `json:"completed_at" db:"completed_at"`
	Duration     int64                  `json:"duration" db:"duration"`
	ErrorMessage string                 `json:"error_message" db:"error_message"`
	Metadata     map[string]interface{} `json:"metadata" db:"metadata"`
}

// ExecutionStatus represents the status of a workflow execution
type ExecutionStatus string

const (
	ExecutionStatusPending   ExecutionStatus = "pending"
	ExecutionStatusRunning   ExecutionStatus = "running"
	ExecutionStatusCompleted ExecutionStatus = "completed"
	ExecutionStatusFailed    ExecutionStatus = "failed"
	ExecutionStatusCancelled ExecutionStatus = "cancelled"
	ExecutionStatusPaused    ExecutionStatus = "paused"
)

// StepExecution represents the execution of a single step
type StepExecution struct {
	StepID       string                 `json:"step_id"`
	Status       ExecutionStatus        `json:"status"`
	Input        map[string]interface{} `json:"input"`
	Output       map[string]interface{} `json:"output"`
	StartedAt    time.Time              `json:"started_at"`
	CompletedAt  *time.Time             `json:"completed_at"`
	Duration     int64                  `json:"duration"`
	ErrorMessage string                 `json:"error_message,omitempty"`
	Retries      int                    `json:"retries"`
}

// AgentReview represents a user review of an agent template
type AgentReview struct {
	ID         uuid.UUID `json:"id" db:"id"`
	AgentID    uuid.UUID `json:"agent_id" db:"agent_id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	UserName   string    `json:"user_name" db:"user_name"`
	Rating     int       `json:"rating" db:"rating"`
	Title      string    `json:"title" db:"title"`
	Comment    string    `json:"comment" db:"comment"`
	IsVerified bool      `json:"is_verified" db:"is_verified"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// AgentPurchase represents a purchase of an agent template
type AgentPurchase struct {
	ID          uuid.UUID `json:"id" db:"id"`
	AgentID     uuid.UUID `json:"agent_id" db:"agent_id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Price       float64   `json:"price" db:"price"`
	Currency    string    `json:"currency" db:"currency"`
	PaymentID   string    `json:"payment_id" db:"payment_id"`
	Status      string    `json:"status" db:"status"`
	PurchasedAt time.Time `json:"purchased_at" db:"purchased_at"`
}

// MarketplaceStats represents marketplace statistics
type MarketplaceStats struct {
	TotalAgents     int     `json:"total_agents"`
	TotalDownloads  int     `json:"total_downloads"`
	TotalRevenue    float64 `json:"total_revenue"`
	TopCategories   []CategoryStats `json:"top_categories"`
	FeaturedAgents  []AgentTemplate `json:"featured_agents"`
	RecentAgents    []AgentTemplate `json:"recent_agents"`
	TopRatedAgents  []AgentTemplate `json:"top_rated_agents"`
}

// CategoryStats represents statistics for a category
type CategoryStats struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
	Revenue  float64 `json:"revenue"`
}

// SearchRequest represents a marketplace search request
type SearchRequest struct {
	Query      string   `json:"query,omitempty"`
	Category   string   `json:"category,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	MinRating  float64  `json:"min_rating,omitempty"`
	MaxPrice   float64  `json:"max_price,omitempty"`
	IsFree     *bool    `json:"is_free,omitempty"`
	SortBy     string   `json:"sort_by,omitempty"`
	SortOrder  string   `json:"sort_order,omitempty"`
	Limit      int      `json:"limit,omitempty"`
	Offset     int      `json:"offset,omitempty"`
}

// SearchResponse represents a marketplace search response
type SearchResponse struct {
	Agents     []AgentTemplate `json:"agents"`
	Total      int             `json:"total"`
	HasMore    bool            `json:"has_more"`
	Facets     SearchFacets    `json:"facets"`
}

// SearchFacets represents search facets for filtering
type SearchFacets struct {
	Categories []FacetItem `json:"categories"`
	Tags       []FacetItem `json:"tags"`
	PriceRanges []FacetItem `json:"price_ranges"`
	Ratings    []FacetItem `json:"ratings"`
}

// FacetItem represents a single facet item
type FacetItem struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

// WorkflowCreateRequest represents a request to create a workflow
type WorkflowCreateRequest struct {
	Name        string             `json:"name" validate:"required"`
	Description string             `json:"description"`
	Category    string             `json:"category" validate:"required"`
	Tags        []string           `json:"tags"`
	IsPublic    bool               `json:"is_public"`
	Price       float64            `json:"price"`
	Workflow    WorkflowDefinition `json:"workflow" validate:"required"`
}

// WorkflowExecuteRequest represents a request to execute a workflow
type WorkflowExecuteRequest struct {
	WorkflowID  uuid.UUID              `json:"workflow_id" validate:"required"`
	Input       map[string]interface{} `json:"input"`
	TriggerType TriggerType            `json:"trigger_type"`
	Async       bool                   `json:"async"`
}

// WorkflowExecuteResponse represents a workflow execution response
type WorkflowExecuteResponse struct {
	ExecutionID uuid.UUID              `json:"execution_id"`
	Status      ExecutionStatus        `json:"status"`
	Output      map[string]interface{} `json:"output,omitempty"`
	Message     string                 `json:"message"`
}

// AgentInstallRequest represents a request to install an agent
type AgentInstallRequest struct {
	AgentID uuid.UUID `json:"agent_id" validate:"required"`
	Version string    `json:"version,omitempty"`
}

// AgentInstallResponse represents an agent installation response
type AgentInstallResponse struct {
	Success   bool   `json:"success"`
	AgentID   uuid.UUID `json:"agent_id"`
	Message   string `json:"message"`
	Workflow  WorkflowDefinition `json:"workflow,omitempty"`
}

// Custom JSON marshaling for WorkflowDefinition
func (w WorkflowDefinition) Value() (interface{}, error) {
	return json.Marshal(w)
}

func (w *WorkflowDefinition) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, w)
	case string:
		return json.Unmarshal([]byte(v), w)
	default:
		return nil
	}
}

// Predefined workflow categories
var WorkflowCategories = []string{
	"Web Scraping",
	"Data Entry",
	"Social Media",
	"E-commerce",
	"Finance & Trading",
	"Content Creation",
	"SEO & Marketing",
	"Testing & QA",
	"Monitoring",
	"Productivity",
	"Research",
	"Communication",
}

// Predefined workflow tags
var WorkflowTags = []string{
	"automation", "scraping", "data", "social", "trading", "content",
	"seo", "testing", "monitoring", "productivity", "research", "email",
	"forms", "reports", "analytics", "notifications", "scheduling",
	"integration", "api", "webhook", "csv", "excel", "pdf", "image",
}
