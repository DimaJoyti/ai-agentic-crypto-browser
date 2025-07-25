package browser

import (
	"time"

	"github.com/google/uuid"
)

// BrowserSession represents a browser session
type BrowserSession struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	SessionName string    `json:"session_name" db:"session_name"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Tabs        []Tab     `json:"tabs,omitempty"`
}

// Tab represents a browser tab
type Tab struct {
	ID        uuid.UUID `json:"id" db:"id"`
	SessionID uuid.UUID `json:"session_id" db:"session_id"`
	URL       string    `json:"url" db:"url"`
	Title     string    `json:"title" db:"title"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// NavigateRequest represents a navigation request
type NavigateRequest struct {
	URL             string            `json:"url" validate:"required"`
	WaitForSelector string            `json:"wait_for_selector,omitempty"`
	Timeout         int               `json:"timeout,omitempty"`
	Headers         map[string]string `json:"headers,omitempty"`
	UserAgent       string            `json:"user_agent,omitempty"`
}

// NavigateResponse represents a navigation response
type NavigateResponse struct {
	Success     bool              `json:"success"`
	URL         string            `json:"url"`
	Title       string            `json:"title"`
	StatusCode  int               `json:"status_code,omitempty"`
	LoadTime    time.Duration     `json:"load_time"`
	Screenshot  string            `json:"screenshot,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Error       string            `json:"error,omitempty"`
}

// InteractRequest represents a page interaction request
type InteractRequest struct {
	Actions     []Action `json:"actions" validate:"required"`
	WaitBetween int      `json:"wait_between,omitempty"` // milliseconds
	Screenshot  bool     `json:"screenshot,omitempty"`
}

// Action represents a single browser action
type Action struct {
	Type     ActionType             `json:"type" validate:"required"`
	Selector string                 `json:"selector,omitempty"`
	Value    string                 `json:"value,omitempty"`
	Options  map[string]interface{} `json:"options,omitempty"`
}

// ActionType represents the type of browser action
type ActionType string

const (
	ActionTypeClick       ActionType = "click"
	ActionTypeType        ActionType = "type"
	ActionTypeSelect      ActionType = "select"
	ActionTypeScroll      ActionType = "scroll"
	ActionTypeWait        ActionType = "wait"
	ActionTypeHover       ActionType = "hover"
	ActionTypeScreenshot  ActionType = "screenshot"
	ActionTypeKeyPress    ActionType = "key_press"
	ActionTypeClear       ActionType = "clear"
	ActionTypeSubmit      ActionType = "submit"
	ActionTypeRefresh     ActionType = "refresh"
	ActionTypeGoBack      ActionType = "go_back"
	ActionTypeGoForward   ActionType = "go_forward"
)

// InteractResponse represents a page interaction response
type InteractResponse struct {
	Success     bool                   `json:"success"`
	Results     []ActionResult         `json:"results"`
	Screenshots []string               `json:"screenshots,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Error       string                 `json:"error,omitempty"`
}

// ActionResult represents the result of a single action
type ActionResult struct {
	Action    Action                 `json:"action"`
	Success   bool                   `json:"success"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Error     string                 `json:"error,omitempty"`
	Duration  time.Duration          `json:"duration"`
}

// ExtractRequest represents a content extraction request
type ExtractRequest struct {
	Selectors []string `json:"selectors,omitempty"`
	DataType  string   `json:"data_type,omitempty"` // text, links, images, tables, forms
	Schema    string   `json:"schema,omitempty"`    // JSON schema for structured extraction
	Options   ExtractOptions `json:"options,omitempty"`
}

// ExtractOptions represents options for content extraction
type ExtractOptions struct {
	IncludeAttributes bool     `json:"include_attributes,omitempty"`
	IncludeStyles     bool     `json:"include_styles,omitempty"`
	MaxDepth          int      `json:"max_depth,omitempty"`
	FilterEmpty       bool     `json:"filter_empty,omitempty"`
	AttributeFilter   []string `json:"attribute_filter,omitempty"`
}

// ExtractResponse represents a content extraction response
type ExtractResponse struct {
	Success    bool                   `json:"success"`
	Data       map[string]interface{} `json:"data"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Screenshot string                 `json:"screenshot,omitempty"`
	Error      string                 `json:"error,omitempty"`
}

// ScreenshotRequest represents a screenshot request
type ScreenshotRequest struct {
	Selector   string `json:"selector,omitempty"`   // CSS selector for element screenshot
	FullPage   bool   `json:"full_page,omitempty"`  // Take full page screenshot
	Quality    int    `json:"quality,omitempty"`    // JPEG quality (1-100)
	Format     string `json:"format,omitempty"`     // png, jpeg
	Width      int    `json:"width,omitempty"`      // Viewport width
	Height     int    `json:"height,omitempty"`     // Viewport height
}

// ScreenshotResponse represents a screenshot response
type ScreenshotResponse struct {
	Success    bool              `json:"success"`
	Screenshot string            `json:"screenshot"` // Base64 encoded image
	Format     string            `json:"format"`
	Size       int               `json:"size"`
	Dimensions map[string]int    `json:"dimensions"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Error      string            `json:"error,omitempty"`
}

// WaitRequest represents a wait request
type WaitRequest struct {
	Type     WaitType `json:"type" validate:"required"`
	Selector string   `json:"selector,omitempty"`
	Text     string   `json:"text,omitempty"`
	Timeout  int      `json:"timeout,omitempty"` // milliseconds
	Visible  bool     `json:"visible,omitempty"`
}

// WaitType represents the type of wait condition
type WaitType string

const (
	WaitTypeSelector    WaitType = "selector"
	WaitTypeText        WaitType = "text"
	WaitTypeNavigation  WaitType = "navigation"
	WaitTypeTimeout     WaitType = "timeout"
	WaitTypeNetworkIdle WaitType = "network_idle"
)

// WaitResponse represents a wait response
type WaitResponse struct {
	Success   bool              `json:"success"`
	Condition string            `json:"condition"`
	Duration  time.Duration     `json:"duration"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Error     string            `json:"error,omitempty"`
}

// PageInfo represents information about the current page
type PageInfo struct {
	URL         string                 `json:"url"`
	Title       string                 `json:"title"`
	Description string                 `json:"description,omitempty"`
	Keywords    []string               `json:"keywords,omitempty"`
	Language    string                 `json:"language,omitempty"`
	Charset     string                 `json:"charset,omitempty"`
	StatusCode  int                    `json:"status_code,omitempty"`
	LoadTime    time.Duration          `json:"load_time"`
	Size        int                    `json:"size,omitempty"`
	Links       []Link                 `json:"links,omitempty"`
	Images      []Image                `json:"images,omitempty"`
	Forms       []Form                 `json:"forms,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// Link represents a link on the page
type Link struct {
	URL    string `json:"url"`
	Text   string `json:"text"`
	Title  string `json:"title,omitempty"`
	Target string `json:"target,omitempty"`
}

// Image represents an image on the page
type Image struct {
	URL    string `json:"url"`
	Alt    string `json:"alt,omitempty"`
	Title  string `json:"title,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

// Form represents a form on the page
type Form struct {
	Action string  `json:"action,omitempty"`
	Method string  `json:"method,omitempty"`
	Fields []Field `json:"fields,omitempty"`
}

// Field represents a form field
type Field struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Value       string `json:"value,omitempty"`
	Placeholder string `json:"placeholder,omitempty"`
	Required    bool   `json:"required,omitempty"`
	Options     []Option `json:"options,omitempty"` // For select fields
}

// Option represents a select option
type Option struct {
	Value    string `json:"value"`
	Text     string `json:"text"`
	Selected bool   `json:"selected,omitempty"`
}

// SessionCreateRequest represents a request to create a browser session
type SessionCreateRequest struct {
	SessionName string `json:"session_name,omitempty"`
	UserAgent   string `json:"user_agent,omitempty"`
	Viewport    Viewport `json:"viewport,omitempty"`
}

// Viewport represents browser viewport settings
type Viewport struct {
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`
}

// SessionCreateResponse represents a response for session creation
type SessionCreateResponse struct {
	Session BrowserSession `json:"session"`
}

// SessionListRequest represents a request to list browser sessions
type SessionListRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	IsActive *bool     `json:"is_active,omitempty"`
	Limit    int       `json:"limit,omitempty"`
	Offset   int       `json:"offset,omitempty"`
}

// SessionListResponse represents a response with session list
type SessionListResponse struct {
	Sessions []BrowserSession `json:"sessions"`
	Total    int              `json:"total"`
	HasMore  bool             `json:"has_more"`
}

// TabCreateRequest represents a request to create a new tab
type TabCreateRequest struct {
	SessionID uuid.UUID `json:"session_id"`
	URL       string    `json:"url,omitempty"`
}

// TabCreateResponse represents a response for tab creation
type TabCreateResponse struct {
	Tab Tab `json:"tab"`
}

// BrowserConfig represents browser configuration
type BrowserConfig struct {
	Headless         bool              `json:"headless"`
	DisableGPU       bool              `json:"disable_gpu"`
	NoSandbox        bool              `json:"no_sandbox"`
	DisableImages    bool              `json:"disable_images,omitempty"`
	DisableJavaScript bool             `json:"disable_javascript,omitempty"`
	UserAgent        string            `json:"user_agent,omitempty"`
	Viewport         Viewport          `json:"viewport,omitempty"`
	Timeout          time.Duration     `json:"timeout,omitempty"`
	Headers          map[string]string `json:"headers,omitempty"`
	Proxy            string            `json:"proxy,omitempty"`
}

// BrowserInstance represents a browser instance
type BrowserInstance struct {
	ID        string        `json:"id"`
	SessionID uuid.UUID     `json:"session_id"`
	Config    BrowserConfig `json:"config"`
	CreatedAt time.Time     `json:"created_at"`
	LastUsed  time.Time     `json:"last_used"`
	IsActive  bool          `json:"is_active"`
}
