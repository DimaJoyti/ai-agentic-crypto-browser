package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// RBACService handles role-based access control
type RBACService struct {
	roles       map[string]*Role
	permissions map[string]*Permission
	policies    map[string]*Policy
}

// Role represents a user role with permissions
type Role struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Permissions []string     `json:"permissions"`
	Inherits    []string     `json:"inherits"`
	IsSystem    bool         `json:"is_system"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// Permission represents a specific permission
type Permission struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Resource    string            `json:"resource"`
	Action      string            `json:"action"`
	Conditions  map[string]string `json:"conditions"`
	IsSystem    bool              `json:"is_system"`
	CreatedAt   time.Time         `json:"created_at"`
}

// Policy represents an access control policy
type Policy struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Rules       []PolicyRule           `json:"rules"`
	Effect      PolicyEffect           `json:"effect"`
	Conditions  map[string]interface{} `json:"conditions"`
	Priority    int                    `json:"priority"`
	IsActive    bool                   `json:"is_active"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// PolicyRule represents a rule within a policy
type PolicyRule struct {
	Resource   string                 `json:"resource"`
	Actions    []string               `json:"actions"`
	Effect     PolicyEffect           `json:"effect"`
	Conditions map[string]interface{} `json:"conditions"`
}

// PolicyEffect represents the effect of a policy
type PolicyEffect string

const (
	PolicyEffectAllow PolicyEffect = "allow"
	PolicyEffectDeny  PolicyEffect = "deny"
)

// AccessRequest represents an access control request
type AccessRequest struct {
	UserID     uuid.UUID              `json:"user_id"`
	Resource   string                 `json:"resource"`
	Action     string                 `json:"action"`
	Context    map[string]interface{} `json:"context"`
	TeamID     *uuid.UUID             `json:"team_id,omitempty"`
	SessionID  string                 `json:"session_id"`
	IPAddress  string                 `json:"ip_address"`
	UserAgent  string                 `json:"user_agent"`
	Timestamp  time.Time              `json:"timestamp"`
}

// AccessDecision represents the result of an access control check
type AccessDecision struct {
	Allowed    bool                   `json:"allowed"`
	Reason     string                 `json:"reason"`
	AppliedPolicies []string          `json:"applied_policies"`
	Context    map[string]interface{} `json:"context"`
	TTL        time.Duration          `json:"ttl"`
}

// NewRBACService creates a new RBAC service
func NewRBACService() *RBACService {
	service := &RBACService{
		roles:       make(map[string]*Role),
		permissions: make(map[string]*Permission),
		policies:    make(map[string]*Policy),
	}

	// Initialize default roles and permissions
	service.initializeDefaults()

	return service
}

// CheckAccess checks if a user has access to perform an action on a resource
func (r *RBACService) CheckAccess(ctx context.Context, req AccessRequest) (*AccessDecision, error) {
	// Get user roles and permissions
	userRoles, err := r.getUserRoles(req.UserID, req.TeamID)
	if err != nil {
		return &AccessDecision{
			Allowed: false,
			Reason:  fmt.Sprintf("failed to get user roles: %v", err),
		}, nil
	}

	userPermissions := r.getUserPermissions(userRoles)

	// Check direct permissions first
	if r.hasDirectPermission(userPermissions, req.Resource, req.Action) {
		return &AccessDecision{
			Allowed: true,
			Reason:  "direct permission granted",
			TTL:     time.Hour,
		}, nil
	}

	// Evaluate policies
	decision := r.evaluatePolicies(req, userRoles, userPermissions)

	// Log access decision for audit
	r.logAccessDecision(ctx, req, decision)

	return decision, nil
}

// HasPermission checks if a user has a specific permission
func (r *RBACService) HasPermission(userID uuid.UUID, teamID *uuid.UUID, permission string) bool {
	userRoles, err := r.getUserRoles(userID, teamID)
	if err != nil {
		return false
	}

	userPermissions := r.getUserPermissions(userRoles)

	for _, perm := range userPermissions {
		if perm == permission {
			return true
		}
	}

	return false
}

// CreateRole creates a new role
func (r *RBACService) CreateRole(role *Role) error {
	if _, exists := r.roles[role.ID]; exists {
		return fmt.Errorf("role %s already exists", role.ID)
	}

	// Validate permissions exist
	for _, permID := range role.Permissions {
		if _, exists := r.permissions[permID]; !exists {
			return fmt.Errorf("permission %s does not exist", permID)
		}
	}

	// Validate inheritance
	for _, inheritID := range role.Inherits {
		if _, exists := r.roles[inheritID]; !exists {
			return fmt.Errorf("inherited role %s does not exist", inheritID)
		}
	}

	role.CreatedAt = time.Now()
	role.UpdatedAt = time.Now()
	r.roles[role.ID] = role

	return nil
}

// CreatePermission creates a new permission
func (r *RBACService) CreatePermission(permission *Permission) error {
	if _, exists := r.permissions[permission.ID]; exists {
		return fmt.Errorf("permission %s already exists", permission.ID)
	}

	permission.CreatedAt = time.Now()
	r.permissions[permission.ID] = permission

	return nil
}

// CreatePolicy creates a new access control policy
func (r *RBACService) CreatePolicy(policy *Policy) error {
	if _, exists := r.policies[policy.ID]; exists {
		return fmt.Errorf("policy %s already exists", policy.ID)
	}

	policy.CreatedAt = time.Now()
	policy.UpdatedAt = time.Now()
	r.policies[policy.ID] = policy

	return nil
}

// AssignRoleToUser assigns a role to a user
func (r *RBACService) AssignRoleToUser(userID uuid.UUID, roleID string, teamID *uuid.UUID) error {
	// In a real implementation, this would update the database
	return nil
}

// RevokeRoleFromUser revokes a role from a user
func (r *RBACService) RevokeRoleFromUser(userID uuid.UUID, roleID string, teamID *uuid.UUID) error {
	// In a real implementation, this would update the database
	return nil
}

// GetUserRoles returns all roles for a user
func (r *RBACService) GetUserRoles(userID uuid.UUID, teamID *uuid.UUID) ([]string, error) {
	return r.getUserRoles(userID, teamID)
}

// GetRolePermissions returns all permissions for a role
func (r *RBACService) GetRolePermissions(roleID string) ([]string, error) {
	role, exists := r.roles[roleID]
	if !exists {
		return nil, fmt.Errorf("role %s not found", roleID)
	}

	permissions := make([]string, len(role.Permissions))
	copy(permissions, role.Permissions)

	// Add inherited permissions
	for _, inheritID := range role.Inherits {
		inheritedPerms, err := r.GetRolePermissions(inheritID)
		if err != nil {
			continue
		}
		permissions = append(permissions, inheritedPerms...)
	}

	// Remove duplicates
	return r.removeDuplicates(permissions), nil
}

// Private methods

func (r *RBACService) initializeDefaults() {
	// System permissions
	systemPermissions := []*Permission{
		{ID: "users.read", Name: "Read Users", Resource: "users", Action: "read", IsSystem: true},
		{ID: "users.write", Name: "Write Users", Resource: "users", Action: "write", IsSystem: true},
		{ID: "users.delete", Name: "Delete Users", Resource: "users", Action: "delete", IsSystem: true},
		{ID: "workflows.read", Name: "Read Workflows", Resource: "workflows", Action: "read", IsSystem: true},
		{ID: "workflows.write", Name: "Write Workflows", Resource: "workflows", Action: "write", IsSystem: true},
		{ID: "workflows.execute", Name: "Execute Workflows", Resource: "workflows", Action: "execute", IsSystem: true},
		{ID: "workflows.delete", Name: "Delete Workflows", Resource: "workflows", Action: "delete", IsSystem: true},
		{ID: "teams.read", Name: "Read Teams", Resource: "teams", Action: "read", IsSystem: true},
		{ID: "teams.write", Name: "Write Teams", Resource: "teams", Action: "write", IsSystem: true},
		{ID: "teams.manage", Name: "Manage Teams", Resource: "teams", Action: "manage", IsSystem: true},
		{ID: "analytics.read", Name: "Read Analytics", Resource: "analytics", Action: "read", IsSystem: true},
		{ID: "system.admin", Name: "System Admin", Resource: "system", Action: "admin", IsSystem: true},
	}

	for _, perm := range systemPermissions {
		perm.CreatedAt = time.Now()
		r.permissions[perm.ID] = perm
	}

	// System roles
	systemRoles := []*Role{
		{
			ID:          "super_admin",
			Name:        "Super Administrator",
			Description: "Full system access",
			Permissions: []string{"system.admin"},
			IsSystem:    true,
		},
		{
			ID:          "admin",
			Name:        "Administrator",
			Description: "Administrative access",
			Permissions: []string{"users.read", "users.write", "workflows.read", "workflows.write", "workflows.execute", "teams.read", "teams.write", "analytics.read"},
			IsSystem:    true,
		},
		{
			ID:          "team_owner",
			Name:        "Team Owner",
			Description: "Team ownership privileges",
			Permissions: []string{"teams.manage", "users.read", "users.write", "workflows.read", "workflows.write", "workflows.execute", "analytics.read"},
			IsSystem:    true,
		},
		{
			ID:          "team_admin",
			Name:        "Team Administrator",
			Description: "Team administrative access",
			Permissions: []string{"users.read", "workflows.read", "workflows.write", "workflows.execute", "analytics.read"},
			IsSystem:    true,
		},
		{
			ID:          "user",
			Name:        "User",
			Description: "Standard user access",
			Permissions: []string{"workflows.read", "workflows.execute"},
			IsSystem:    true,
		},
		{
			ID:          "viewer",
			Name:        "Viewer",
			Description: "Read-only access",
			Permissions: []string{"workflows.read"},
			IsSystem:    true,
		},
	}

	for _, role := range systemRoles {
		role.CreatedAt = time.Now()
		role.UpdatedAt = time.Now()
		r.roles[role.ID] = role
	}

	// Default policies
	defaultPolicies := []*Policy{
		{
			ID:          "team_isolation",
			Name:        "Team Isolation Policy",
			Description: "Users can only access resources within their team",
			Effect:      PolicyEffectDeny,
			Rules: []PolicyRule{
				{
					Resource: "workflows",
					Actions:  []string{"read", "write", "execute"},
					Effect:   PolicyEffectDeny,
					Conditions: map[string]interface{}{
						"team_mismatch": true,
					},
				},
			},
			Priority: 100,
			IsActive: true,
		},
		{
			ID:          "rate_limiting",
			Name:        "Rate Limiting Policy",
			Description: "Enforce rate limits on API access",
			Effect:      PolicyEffectDeny,
			Rules: []PolicyRule{
				{
					Resource: "api",
					Actions:  []string{"*"},
					Effect:   PolicyEffectDeny,
					Conditions: map[string]interface{}{
						"rate_limit_exceeded": true,
					},
				},
			},
			Priority: 200,
			IsActive: true,
		},
	}

	for _, policy := range defaultPolicies {
		policy.CreatedAt = time.Now()
		policy.UpdatedAt = time.Now()
		r.policies[policy.ID] = policy
	}
}

func (r *RBACService) getUserRoles(userID uuid.UUID, teamID *uuid.UUID) ([]string, error) {
	// In a real implementation, this would query the database
	// For now, return mock data based on user ID
	if userID.String() == "00000000-0000-0000-0000-000000000001" {
		return []string{"super_admin"}, nil
	}
	return []string{"user"}, nil
}

func (r *RBACService) getUserPermissions(roles []string) []string {
	var permissions []string

	for _, roleID := range roles {
		rolePerms, err := r.GetRolePermissions(roleID)
		if err != nil {
			continue
		}
		permissions = append(permissions, rolePerms...)
	}

	return r.removeDuplicates(permissions)
}

func (r *RBACService) hasDirectPermission(userPermissions []string, resource, action string) bool {
	requiredPerm := fmt.Sprintf("%s.%s", resource, action)
	wildcardPerm := fmt.Sprintf("%s.*", resource)
	globalPerm := "system.admin"

	for _, perm := range userPermissions {
		if perm == requiredPerm || perm == wildcardPerm || perm == globalPerm {
			return true
		}
	}

	return false
}

func (r *RBACService) evaluatePolicies(req AccessRequest, userRoles, userPermissions []string) *AccessDecision {
	var appliedPolicies []string
	allowed := true
	reason := "no applicable policies"

	// Sort policies by priority (higher priority first)
	policies := r.getSortedPolicies()

	for _, policy := range policies {
		if !policy.IsActive {
			continue
		}

		if r.policyApplies(policy, req, userRoles) {
			appliedPolicies = append(appliedPolicies, policy.ID)

			if policy.Effect == PolicyEffectDeny {
				allowed = false
				reason = fmt.Sprintf("denied by policy: %s", policy.Name)
				break
			}
		}
	}

	return &AccessDecision{
		Allowed:         allowed,
		Reason:          reason,
		AppliedPolicies: appliedPolicies,
		TTL:             time.Minute * 5,
	}
}

func (r *RBACService) policyApplies(policy *Policy, req AccessRequest, userRoles []string) bool {
	for _, rule := range policy.Rules {
		if r.ruleMatches(rule, req) {
			return true
		}
	}
	return false
}

func (r *RBACService) ruleMatches(rule PolicyRule, req AccessRequest) bool {
	// Check resource match
	if rule.Resource != "*" && rule.Resource != req.Resource {
		return false
	}

	// Check action match
	actionMatches := false
	for _, action := range rule.Actions {
		if action == "*" || action == req.Action {
			actionMatches = true
			break
		}
	}

	if !actionMatches {
		return false
	}

	// Check conditions
	return r.evaluateConditions(rule.Conditions, req)
}

func (r *RBACService) evaluateConditions(conditions map[string]interface{}, req AccessRequest) bool {
	for key, value := range conditions {
		switch key {
		case "team_mismatch":
			if value.(bool) && req.TeamID == nil {
				return true
			}
		case "rate_limit_exceeded":
			// In a real implementation, check rate limiting
			return false
		case "ip_whitelist":
			whitelist := value.([]string)
			if !r.ipInWhitelist(req.IPAddress, whitelist) {
				return true
			}
		case "time_restriction":
			// Check time-based restrictions
			return r.checkTimeRestriction(value.(map[string]interface{}))
		}
	}

	return false
}

func (r *RBACService) getSortedPolicies() []*Policy {
	policies := make([]*Policy, 0, len(r.policies))
	for _, policy := range r.policies {
		policies = append(policies, policy)
	}

	// Sort by priority (higher first)
	for i := 0; i < len(policies)-1; i++ {
		for j := i + 1; j < len(policies); j++ {
			if policies[i].Priority < policies[j].Priority {
				policies[i], policies[j] = policies[j], policies[i]
			}
		}
	}

	return policies
}

func (r *RBACService) removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	var result []string

	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}

	return result
}

func (r *RBACService) ipInWhitelist(ip string, whitelist []string) bool {
	for _, whitelistedIP := range whitelist {
		if ip == whitelistedIP {
			return true
		}
	}
	return false
}

func (r *RBACService) checkTimeRestriction(restriction map[string]interface{}) bool {
	// Implement time-based access restrictions
	return false
}

func (r *RBACService) logAccessDecision(ctx context.Context, req AccessRequest, decision *AccessDecision) {
	// In a real implementation, log to audit system
	// auditLogger.LogAccessDecision(ctx, req, decision)
}
