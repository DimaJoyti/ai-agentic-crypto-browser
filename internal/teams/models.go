package teams

import (
	"time"

	"github.com/google/uuid"
)

// Team represents a team/organization
type Team struct {
	ID          uuid.UUID              `json:"id" db:"id"`
	Name        string                 `json:"name" db:"name"`
	Description string                 `json:"description" db:"description"`
	OwnerID     uuid.UUID              `json:"owner_id" db:"owner_id"`
	Plan        TeamPlan               `json:"plan" db:"plan"`
	Settings    TeamSettings           `json:"settings" db:"settings"`
	Metadata    map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
}

// TeamPlan represents different team subscription plans
type TeamPlan string

const (
	TeamPlanFree       TeamPlan = "free"
	TeamPlanPro        TeamPlan = "pro"
	TeamPlanEnterprise TeamPlan = "enterprise"
)

// TeamSettings contains team configuration
type TeamSettings struct {
	MaxMembers           int                    `json:"max_members"`
	MaxWorkflows         int                    `json:"max_workflows"`
	MaxExecutionsPerDay  int                    `json:"max_executions_per_day"`
	AllowPublicWorkflows bool                   `json:"allow_public_workflows"`
	RequireApproval      bool                   `json:"require_approval"`
	DefaultPermissions   []Permission           `json:"default_permissions"`
	IntegrationSettings  map[string]interface{} `json:"integration_settings"`
	SecuritySettings     SecuritySettings       `json:"security_settings"`
}

// SecuritySettings contains team security configuration
type SecuritySettings struct {
	RequireMFA           bool     `json:"require_mfa"`
	AllowedDomains       []string `json:"allowed_domains"`
	SessionTimeout       int      `json:"session_timeout"`
	IPWhitelist          []string `json:"ip_whitelist"`
	AuditLogRetention    int      `json:"audit_log_retention"`
	EncryptionRequired   bool     `json:"encryption_required"`
}

// TeamMember represents a team member
type TeamMember struct {
	ID          uuid.UUID   `json:"id" db:"id"`
	TeamID      uuid.UUID   `json:"team_id" db:"team_id"`
	UserID      uuid.UUID   `json:"user_id" db:"user_id"`
	Role        MemberRole  `json:"role" db:"role"`
	Permissions []Permission `json:"permissions" db:"permissions"`
	Status      MemberStatus `json:"status" db:"status"`
	InvitedBy   uuid.UUID   `json:"invited_by" db:"invited_by"`
	JoinedAt    time.Time   `json:"joined_at" db:"joined_at"`
	LastActive  time.Time   `json:"last_active" db:"last_active"`
}

// MemberRole represents different member roles
type MemberRole string

const (
	MemberRoleOwner     MemberRole = "owner"
	MemberRoleAdmin     MemberRole = "admin"
	MemberRoleMember    MemberRole = "member"
	MemberRoleViewer    MemberRole = "viewer"
	MemberRoleGuest     MemberRole = "guest"
)

// MemberStatus represents member status
type MemberStatus string

const (
	MemberStatusActive   MemberStatus = "active"
	MemberStatusInvited  MemberStatus = "invited"
	MemberStatusSuspended MemberStatus = "suspended"
	MemberStatusLeft     MemberStatus = "left"
)

// Permission represents a specific permission
type Permission string

const (
	PermissionViewWorkflows    Permission = "view_workflows"
	PermissionCreateWorkflows  Permission = "create_workflows"
	PermissionEditWorkflows    Permission = "edit_workflows"
	PermissionDeleteWorkflows  Permission = "delete_workflows"
	PermissionExecuteWorkflows Permission = "execute_workflows"
	PermissionViewExecutions   Permission = "view_executions"
	PermissionManageMembers    Permission = "manage_members"
	PermissionManageSettings   Permission = "manage_settings"
	PermissionViewAnalytics    Permission = "view_analytics"
	PermissionManageBilling    Permission = "manage_billing"
	PermissionViewAuditLogs    Permission = "view_audit_logs"
	PermissionManageIntegrations Permission = "manage_integrations"
)

// TeamInvitation represents a team invitation
type TeamInvitation struct {
	ID          uuid.UUID   `json:"id" db:"id"`
	TeamID      uuid.UUID   `json:"team_id" db:"team_id"`
	Email       string      `json:"email" db:"email"`
	Role        MemberRole  `json:"role" db:"role"`
	Permissions []Permission `json:"permissions" db:"permissions"`
	InvitedBy   uuid.UUID   `json:"invited_by" db:"invited_by"`
	Token       string      `json:"token" db:"token"`
	ExpiresAt   time.Time   `json:"expires_at" db:"expires_at"`
	AcceptedAt  *time.Time  `json:"accepted_at" db:"accepted_at"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
}

// SharedWorkflow represents a workflow shared within a team
type SharedWorkflow struct {
	ID           uuid.UUID              `json:"id" db:"id"`
	TeamID       uuid.UUID              `json:"team_id" db:"team_id"`
	WorkflowID   uuid.UUID              `json:"workflow_id" db:"workflow_id"`
	SharedBy     uuid.UUID              `json:"shared_by" db:"shared_by"`
	ShareType    ShareType              `json:"share_type" db:"share_type"`
	Permissions  []Permission           `json:"permissions" db:"permissions"`
	AccessLevel  AccessLevel            `json:"access_level" db:"access_level"`
	Metadata     map[string]interface{} `json:"metadata" db:"metadata"`
	SharedAt     time.Time              `json:"shared_at" db:"shared_at"`
}

// ShareType represents how a workflow is shared
type ShareType string

const (
	ShareTypeTeam   ShareType = "team"
	ShareTypePublic ShareType = "public"
	ShareTypeLink   ShareType = "link"
)

// AccessLevel represents access level for shared workflows
type AccessLevel string

const (
	AccessLevelView    AccessLevel = "view"
	AccessLevelExecute AccessLevel = "execute"
	AccessLevelEdit    AccessLevel = "edit"
	AccessLevelAdmin   AccessLevel = "admin"
)

// TeamWorkspace represents a team workspace
type TeamWorkspace struct {
	ID          uuid.UUID              `json:"id" db:"id"`
	TeamID      uuid.UUID              `json:"team_id" db:"team_id"`
	Name        string                 `json:"name" db:"name"`
	Description string                 `json:"description" db:"description"`
	CreatedBy   uuid.UUID              `json:"created_by" db:"created_by"`
	Settings    WorkspaceSettings      `json:"settings" db:"settings"`
	Metadata    map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
}

// WorkspaceSettings contains workspace configuration
type WorkspaceSettings struct {
	DefaultBrowserSettings map[string]interface{} `json:"default_browser_settings"`
	DefaultAISettings      map[string]interface{} `json:"default_ai_settings"`
	DefaultWeb3Settings    map[string]interface{} `json:"default_web3_settings"`
	ResourceLimits         ResourceLimits         `json:"resource_limits"`
	NotificationSettings   NotificationSettings   `json:"notification_settings"`
}

// ResourceLimits defines resource usage limits
type ResourceLimits struct {
	MaxConcurrentExecutions int `json:"max_concurrent_executions"`
	MaxExecutionTime        int `json:"max_execution_time"`
	MaxMemoryUsage          int `json:"max_memory_usage"`
	MaxStorageUsage         int `json:"max_storage_usage"`
}

// NotificationSettings defines notification preferences
type NotificationSettings struct {
	EmailNotifications    bool     `json:"email_notifications"`
	SlackIntegration      bool     `json:"slack_integration"`
	WebhookURL            string   `json:"webhook_url"`
	NotificationChannels  []string `json:"notification_channels"`
	NotificationEvents    []string `json:"notification_events"`
}

// TeamAnalytics represents team usage analytics
type TeamAnalytics struct {
	TeamID              uuid.UUID              `json:"team_id"`
	Period              string                 `json:"period"`
	TotalExecutions     int                    `json:"total_executions"`
	SuccessfulExecutions int                   `json:"successful_executions"`
	FailedExecutions    int                    `json:"failed_executions"`
	TotalExecutionTime  int64                  `json:"total_execution_time"`
	AverageExecutionTime float64               `json:"average_execution_time"`
	TopWorkflows        []WorkflowUsage        `json:"top_workflows"`
	TopUsers            []UserUsage            `json:"top_users"`
	ResourceUsage       ResourceUsage          `json:"resource_usage"`
	CostBreakdown       CostBreakdown          `json:"cost_breakdown"`
	Metadata            map[string]interface{} `json:"metadata"`
	GeneratedAt         time.Time              `json:"generated_at"`
}

// WorkflowUsage represents workflow usage statistics
type WorkflowUsage struct {
	WorkflowID   uuid.UUID `json:"workflow_id"`
	WorkflowName string    `json:"workflow_name"`
	Executions   int       `json:"executions"`
	SuccessRate  float64   `json:"success_rate"`
	AvgDuration  float64   `json:"avg_duration"`
}

// UserUsage represents user usage statistics
type UserUsage struct {
	UserID      uuid.UUID `json:"user_id"`
	UserName    string    `json:"user_name"`
	Executions  int       `json:"executions"`
	SuccessRate float64   `json:"success_rate"`
	LastActive  time.Time `json:"last_active"`
}

// ResourceUsage represents resource usage statistics
type ResourceUsage struct {
	CPUUsage     float64 `json:"cpu_usage"`
	MemoryUsage  float64 `json:"memory_usage"`
	StorageUsage float64 `json:"storage_usage"`
	NetworkUsage float64 `json:"network_usage"`
	APIUsage     int     `json:"api_usage"`
}

// CostBreakdown represents cost breakdown
type CostBreakdown struct {
	TotalCost       float64            `json:"total_cost"`
	AIServiceCost   float64            `json:"ai_service_cost"`
	BrowserCost     float64            `json:"browser_cost"`
	Web3Cost        float64            `json:"web3_cost"`
	StorageCost     float64            `json:"storage_cost"`
	BandwidthCost   float64            `json:"bandwidth_cost"`
	DetailedCosts   map[string]float64 `json:"detailed_costs"`
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID          uuid.UUID              `json:"id" db:"id"`
	TeamID      uuid.UUID              `json:"team_id" db:"team_id"`
	UserID      uuid.UUID              `json:"user_id" db:"user_id"`
	Action      string                 `json:"action" db:"action"`
	Resource    string                 `json:"resource" db:"resource"`
	ResourceID  uuid.UUID              `json:"resource_id" db:"resource_id"`
	Details     map[string]interface{} `json:"details" db:"details"`
	IPAddress   string                 `json:"ip_address" db:"ip_address"`
	UserAgent   string                 `json:"user_agent" db:"user_agent"`
	Timestamp   time.Time              `json:"timestamp" db:"timestamp"`
}

// TeamIntegration represents external integrations
type TeamIntegration struct {
	ID            uuid.UUID              `json:"id" db:"id"`
	TeamID        uuid.UUID              `json:"team_id" db:"team_id"`
	IntegrationType string               `json:"integration_type" db:"integration_type"`
	Name          string                 `json:"name" db:"name"`
	Configuration map[string]interface{} `json:"configuration" db:"configuration"`
	IsActive      bool                   `json:"is_active" db:"is_active"`
	CreatedBy     uuid.UUID              `json:"created_by" db:"created_by"`
	CreatedAt     time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at" db:"updated_at"`
}

// Request/Response types

// CreateTeamRequest represents a request to create a team
type CreateTeamRequest struct {
	Name        string                 `json:"name" validate:"required"`
	Description string                 `json:"description"`
	Plan        TeamPlan               `json:"plan"`
	Settings    TeamSettings           `json:"settings"`
}

// InviteMemberRequest represents a request to invite a team member
type InviteMemberRequest struct {
	Email       string      `json:"email" validate:"required,email"`
	Role        MemberRole  `json:"role" validate:"required"`
	Permissions []Permission `json:"permissions"`
	Message     string      `json:"message"`
}

// UpdateMemberRequest represents a request to update a team member
type UpdateMemberRequest struct {
	Role        *MemberRole  `json:"role,omitempty"`
	Permissions []Permission `json:"permissions,omitempty"`
	Status      *MemberStatus `json:"status,omitempty"`
}

// ShareWorkflowRequest represents a request to share a workflow
type ShareWorkflowRequest struct {
	WorkflowID  uuid.UUID   `json:"workflow_id" validate:"required"`
	ShareType   ShareType   `json:"share_type" validate:"required"`
	AccessLevel AccessLevel `json:"access_level" validate:"required"`
	Permissions []Permission `json:"permissions"`
	ExpiresAt   *time.Time  `json:"expires_at,omitempty"`
}

// CreateWorkspaceRequest represents a request to create a workspace
type CreateWorkspaceRequest struct {
	Name        string            `json:"name" validate:"required"`
	Description string            `json:"description"`
	Settings    WorkspaceSettings `json:"settings"`
}

// TeamStatsResponse represents team statistics
type TeamStatsResponse struct {
	TotalMembers     int                `json:"total_members"`
	ActiveMembers    int                `json:"active_members"`
	TotalWorkflows   int                `json:"total_workflows"`
	SharedWorkflows  int                `json:"shared_workflows"`
	TotalExecutions  int                `json:"total_executions"`
	RecentActivity   []ActivityItem     `json:"recent_activity"`
	ResourceUsage    ResourceUsage      `json:"resource_usage"`
	PlanLimits       TeamPlanLimits     `json:"plan_limits"`
}

// ActivityItem represents a recent activity item
type ActivityItem struct {
	ID          uuid.UUID              `json:"id"`
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	UserID      uuid.UUID              `json:"user_id"`
	UserName    string                 `json:"user_name"`
	Metadata    map[string]interface{} `json:"metadata"`
	Timestamp   time.Time              `json:"timestamp"`
}

// TeamPlanLimits represents limits for different team plans
type TeamPlanLimits struct {
	MaxMembers          int `json:"max_members"`
	MaxWorkflows        int `json:"max_workflows"`
	MaxExecutionsPerDay int `json:"max_executions_per_day"`
	MaxStorageGB        int `json:"max_storage_gb"`
	MaxConcurrentRuns   int `json:"max_concurrent_runs"`
}

// Predefined team plan limits
var TeamPlanLimitsMap = map[TeamPlan]TeamPlanLimits{
	TeamPlanFree: {
		MaxMembers:          5,
		MaxWorkflows:        10,
		MaxExecutionsPerDay: 100,
		MaxStorageGB:        1,
		MaxConcurrentRuns:   2,
	},
	TeamPlanPro: {
		MaxMembers:          25,
		MaxWorkflows:        100,
		MaxExecutionsPerDay: 1000,
		MaxStorageGB:        10,
		MaxConcurrentRuns:   10,
	},
	TeamPlanEnterprise: {
		MaxMembers:          -1, // Unlimited
		MaxWorkflows:        -1, // Unlimited
		MaxExecutionsPerDay: -1, // Unlimited
		MaxStorageGB:        100,
		MaxConcurrentRuns:   50,
	},
}

// Role permissions mapping
var RolePermissions = map[MemberRole][]Permission{
	MemberRoleOwner: {
		PermissionViewWorkflows, PermissionCreateWorkflows, PermissionEditWorkflows,
		PermissionDeleteWorkflows, PermissionExecuteWorkflows, PermissionViewExecutions,
		PermissionManageMembers, PermissionManageSettings, PermissionViewAnalytics,
		PermissionManageBilling, PermissionViewAuditLogs, PermissionManageIntegrations,
	},
	MemberRoleAdmin: {
		PermissionViewWorkflows, PermissionCreateWorkflows, PermissionEditWorkflows,
		PermissionDeleteWorkflows, PermissionExecuteWorkflows, PermissionViewExecutions,
		PermissionManageMembers, PermissionManageSettings, PermissionViewAnalytics,
		PermissionViewAuditLogs, PermissionManageIntegrations,
	},
	MemberRoleMember: {
		PermissionViewWorkflows, PermissionCreateWorkflows, PermissionEditWorkflows,
		PermissionExecuteWorkflows, PermissionViewExecutions,
	},
	MemberRoleViewer: {
		PermissionViewWorkflows, PermissionViewExecutions,
	},
	MemberRoleGuest: {
		PermissionViewWorkflows,
	},
}
