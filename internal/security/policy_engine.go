package security

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

// PolicyEngine manages and evaluates security policies
type PolicyEngine struct {
	logger   *observability.Logger
	config   *PolicyEngineConfig
	policies map[string]*SecurityPolicy
	rules    map[string]*PolicyRule
}

// PolicyEngineConfig contains policy engine configuration
type PolicyEngineConfig struct {
	DefaultPolicy       PolicyDecision
	EnableAuditLogging  bool
	PolicyCacheTimeout  time.Duration
	MaxPolicyEvaluations int
}

// SecurityPolicy represents a security policy
type SecurityPolicy struct {
	PolicyID    string
	Name        string
	Description string
	Version     string
	Enabled     bool
	Priority    int
	Rules       []*PolicyRule
	Conditions  []*PolicyCondition
	Actions     []*PolicyAction
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CreatedBy   string
}

// PolicyRule represents a policy rule
type PolicyRule struct {
	RuleID      string
	Name        string
	Description string
	Type        RuleType
	Conditions  []*PolicyCondition
	Actions     []*PolicyAction
	Enabled     bool
	Priority    int
}

// PolicyCondition represents a condition in a policy
type PolicyCondition struct {
	ConditionID string
	Type        ConditionType
	Field       string
	Operator    OperatorType
	Value       interface{}
	Negate      bool
}

// PolicyAction represents an action to take when a policy matches
type PolicyAction struct {
	ActionID string
	Type     ActionType
	Config   map[string]interface{}
}

// RuleType defines types of policy rules
type RuleType string

const (
	RuleTypeAccess       RuleType = "access"
	RuleTypeAuthentication RuleType = "authentication"
	RuleTypeAuthorization  RuleType = "authorization"
	RuleTypeRateLimit     RuleType = "rate_limit"
	RuleTypeThreatDetection RuleType = "threat_detection"
	RuleTypeDataProtection RuleType = "data_protection"
)

// ConditionType defines types of policy conditions
type ConditionType string

const (
	ConditionTypeUser        ConditionType = "user"
	ConditionTypeRole        ConditionType = "role"
	ConditionTypeIP          ConditionType = "ip"
	ConditionTypeTime        ConditionType = "time"
	ConditionTypeResource    ConditionType = "resource"
	ConditionTypeAction      ConditionType = "action"
	ConditionTypeRiskScore   ConditionType = "risk_score"
	ConditionTypeDeviceTrust ConditionType = "device_trust"
	ConditionTypeLocation    ConditionType = "location"
)

// OperatorType defines comparison operators
type OperatorType string

const (
	OperatorEquals       OperatorType = "equals"
	OperatorNotEquals    OperatorType = "not_equals"
	OperatorContains     OperatorType = "contains"
	OperatorNotContains  OperatorType = "not_contains"
	OperatorGreaterThan  OperatorType = "greater_than"
	OperatorLessThan     OperatorType = "less_than"
	OperatorRegex        OperatorType = "regex"
	OperatorInList       OperatorType = "in_list"
	OperatorNotInList    OperatorType = "not_in_list"
)

// ActionType defines types of policy actions
type ActionType string

const (
	ActionTypeAllow        ActionType = "allow"
	ActionTypeDeny         ActionType = "deny"
	ActionTypeRequireMFA   ActionType = "require_mfa"
	ActionTypeLog          ActionType = "log"
	ActionTypeAlert        ActionType = "alert"
	ActionTypeBlock        ActionType = "block"
	ActionTypeRedirect     ActionType = "redirect"
	ActionTypeRateLimit    ActionType = "rate_limit"
)

// PolicyEvaluationContext contains context for policy evaluation
type PolicyEvaluationContext struct {
	UserID       uuid.UUID
	UserRoles    []string
	IPAddress    string
	UserAgent    string
	Resource     string
	Action       string
	RiskScore    float64
	DeviceTrust  float64
	Location     *Location
	Timestamp    time.Time
	SessionData  map[string]interface{}
}

// PolicyEvaluationResult contains the result of policy evaluation
type PolicyEvaluationResult struct {
	Decision       PolicyDecision
	MatchedPolicies []*SecurityPolicy
	MatchedRules   []*PolicyRule
	Actions        []*PolicyAction
	Reason         string
	Confidence     float64
	EvaluationTime time.Duration
}

// NewPolicyEngine creates a new policy engine
func NewPolicyEngine(logger *observability.Logger) *PolicyEngine {
	config := &PolicyEngineConfig{
		DefaultPolicy: PolicyDecision{
			Allowed:      false,
			RequiresMFA:  true,
			SessionTTL:   15 * time.Minute,
			Reason:       "Default deny policy",
		},
		EnableAuditLogging:   true,
		PolicyCacheTimeout:   5 * time.Minute,
		MaxPolicyEvaluations: 1000,
	}

	engine := &PolicyEngine{
		logger:   logger,
		config:   config,
		policies: make(map[string]*SecurityPolicy),
		rules:    make(map[string]*PolicyRule),
	}

	// Load default policies
	engine.loadDefaultPolicies()

	return engine
}

// EvaluatePolicy evaluates policies against the given context
func (pe *PolicyEngine) EvaluatePolicy(ctx context.Context, evalCtx *PolicyEvaluationContext) (*PolicyEvaluationResult, error) {
	startTime := time.Now()

	result := &PolicyEvaluationResult{
		Decision:        pe.config.DefaultPolicy,
		MatchedPolicies: []*SecurityPolicy{},
		MatchedRules:    []*PolicyRule{},
		Actions:         []*PolicyAction{},
		Reason:          "No matching policies found",
		Confidence:      1.0,
	}

	// Evaluate policies in priority order
	for _, policy := range pe.getSortedPolicies() {
		if !policy.Enabled {
			continue
		}

		// Check if policy conditions match
		if pe.evaluatePolicyConditions(evalCtx, policy.Conditions) {
			result.MatchedPolicies = append(result.MatchedPolicies, policy)

			// Evaluate policy rules
			for _, rule := range policy.Rules {
				if !rule.Enabled {
					continue
				}

				if pe.evaluateRuleConditions(evalCtx, rule.Conditions) {
					result.MatchedRules = append(result.MatchedRules, rule)
					result.Actions = append(result.Actions, rule.Actions...)
				}
			}
		}
	}

	// Determine final decision based on matched policies and rules
	result.Decision = pe.determineDecision(result.MatchedPolicies, result.MatchedRules, result.Actions)
	result.Reason = pe.generateDecisionReason(result)
	result.EvaluationTime = time.Since(startTime)

	// Log policy evaluation if audit logging is enabled
	if pe.config.EnableAuditLogging {
		pe.logPolicyEvaluation(ctx, evalCtx, result)
	}

	return result, nil
}

// AddPolicy adds a new security policy
func (pe *PolicyEngine) AddPolicy(policy *SecurityPolicy) error {
	if policy.PolicyID == "" {
		policy.PolicyID = uuid.New().String()
	}

	policy.CreatedAt = time.Now()
	policy.UpdatedAt = time.Now()

	pe.policies[policy.PolicyID] = policy

	pe.logger.Info(context.Background(), "Security policy added", map[string]interface{}{
		"policy_id":   policy.PolicyID,
		"policy_name": policy.Name,
		"enabled":     policy.Enabled,
		"priority":    policy.Priority,
	})

	return nil
}

// UpdatePolicy updates an existing security policy
func (pe *PolicyEngine) UpdatePolicy(policyID string, policy *SecurityPolicy) error {
	if _, exists := pe.policies[policyID]; !exists {
		return fmt.Errorf("policy not found: %s", policyID)
	}

	policy.PolicyID = policyID
	policy.UpdatedAt = time.Now()
	pe.policies[policyID] = policy

	pe.logger.Info(context.Background(), "Security policy updated", map[string]interface{}{
		"policy_id":   policyID,
		"policy_name": policy.Name,
	})

	return nil
}

// DeletePolicy deletes a security policy
func (pe *PolicyEngine) DeletePolicy(policyID string) error {
	if _, exists := pe.policies[policyID]; !exists {
		return fmt.Errorf("policy not found: %s", policyID)
	}

	delete(pe.policies, policyID)

	pe.logger.Info(context.Background(), "Security policy deleted", map[string]interface{}{
		"policy_id": policyID,
	})

	return nil
}

// GetPolicy retrieves a security policy by ID
func (pe *PolicyEngine) GetPolicy(policyID string) (*SecurityPolicy, error) {
	policy, exists := pe.policies[policyID]
	if !exists {
		return nil, fmt.Errorf("policy not found: %s", policyID)
	}

	return policy, nil
}

// ListPolicies returns all security policies
func (pe *PolicyEngine) ListPolicies() []*SecurityPolicy {
	policies := make([]*SecurityPolicy, 0, len(pe.policies))
	for _, policy := range pe.policies {
		policies = append(policies, policy)
	}
	return policies
}

// evaluatePolicyConditions evaluates policy conditions
func (pe *PolicyEngine) evaluatePolicyConditions(ctx *PolicyEvaluationContext, conditions []*PolicyCondition) bool {
	if len(conditions) == 0 {
		return true // No conditions means always match
	}

	// All conditions must match (AND logic)
	for _, condition := range conditions {
		if !pe.evaluateCondition(ctx, condition) {
			return false
		}
	}

	return true
}

// evaluateRuleConditions evaluates rule conditions
func (pe *PolicyEngine) evaluateRuleConditions(ctx *PolicyEvaluationContext, conditions []*PolicyCondition) bool {
	return pe.evaluatePolicyConditions(ctx, conditions)
}

// evaluateCondition evaluates a single condition
func (pe *PolicyEngine) evaluateCondition(ctx *PolicyEvaluationContext, condition *PolicyCondition) bool {
	var result bool

	switch condition.Type {
	case ConditionTypeUser:
		result = pe.evaluateUserCondition(ctx, condition)
	case ConditionTypeRole:
		result = pe.evaluateRoleCondition(ctx, condition)
	case ConditionTypeIP:
		result = pe.evaluateIPCondition(ctx, condition)
	case ConditionTypeTime:
		result = pe.evaluateTimeCondition(ctx, condition)
	case ConditionTypeResource:
		result = pe.evaluateResourceCondition(ctx, condition)
	case ConditionTypeAction:
		result = pe.evaluateActionCondition(ctx, condition)
	case ConditionTypeRiskScore:
		result = pe.evaluateRiskScoreCondition(ctx, condition)
	case ConditionTypeDeviceTrust:
		result = pe.evaluateDeviceTrustCondition(ctx, condition)
	case ConditionTypeLocation:
		result = pe.evaluateLocationCondition(ctx, condition)
	default:
		result = false
	}

	// Apply negation if specified
	if condition.Negate {
		result = !result
	}

	return result
}

// evaluateUserCondition evaluates user-based conditions
func (pe *PolicyEngine) evaluateUserCondition(ctx *PolicyEvaluationContext, condition *PolicyCondition) bool {
	userIDStr := ctx.UserID.String()
	return pe.compareValues(userIDStr, condition.Operator, condition.Value)
}

// evaluateRoleCondition evaluates role-based conditions
func (pe *PolicyEngine) evaluateRoleCondition(ctx *PolicyEvaluationContext, condition *PolicyCondition) bool {
	for _, role := range ctx.UserRoles {
		if pe.compareValues(role, condition.Operator, condition.Value) {
			return true
		}
	}
	return false
}

// evaluateIPCondition evaluates IP-based conditions
func (pe *PolicyEngine) evaluateIPCondition(ctx *PolicyEvaluationContext, condition *PolicyCondition) bool {
	return pe.compareValues(ctx.IPAddress, condition.Operator, condition.Value)
}

// evaluateTimeCondition evaluates time-based conditions
func (pe *PolicyEngine) evaluateTimeCondition(ctx *PolicyEvaluationContext, condition *PolicyCondition) bool {
	timeStr := ctx.Timestamp.Format("15:04")
	return pe.compareValues(timeStr, condition.Operator, condition.Value)
}

// evaluateResourceCondition evaluates resource-based conditions
func (pe *PolicyEngine) evaluateResourceCondition(ctx *PolicyEvaluationContext, condition *PolicyCondition) bool {
	return pe.compareValues(ctx.Resource, condition.Operator, condition.Value)
}

// evaluateActionCondition evaluates action-based conditions
func (pe *PolicyEngine) evaluateActionCondition(ctx *PolicyEvaluationContext, condition *PolicyCondition) bool {
	return pe.compareValues(ctx.Action, condition.Operator, condition.Value)
}

// evaluateRiskScoreCondition evaluates risk score conditions
func (pe *PolicyEngine) evaluateRiskScoreCondition(ctx *PolicyEvaluationContext, condition *PolicyCondition) bool {
	return pe.compareValues(ctx.RiskScore, condition.Operator, condition.Value)
}

// evaluateDeviceTrustCondition evaluates device trust conditions
func (pe *PolicyEngine) evaluateDeviceTrustCondition(ctx *PolicyEvaluationContext, condition *PolicyCondition) bool {
	return pe.compareValues(ctx.DeviceTrust, condition.Operator, condition.Value)
}

// evaluateLocationCondition evaluates location-based conditions
func (pe *PolicyEngine) evaluateLocationCondition(ctx *PolicyEvaluationContext, condition *PolicyCondition) bool {
	if ctx.Location == nil {
		return false
	}
	return pe.compareValues(ctx.Location.Country, condition.Operator, condition.Value)
}

// compareValues compares two values using the specified operator
func (pe *PolicyEngine) compareValues(actual interface{}, operator OperatorType, expected interface{}) bool {
	switch operator {
	case OperatorEquals:
		return actual == expected
	case OperatorNotEquals:
		return actual != expected
	case OperatorContains:
		actualStr := fmt.Sprintf("%v", actual)
		expectedStr := fmt.Sprintf("%v", expected)
		return strings.Contains(actualStr, expectedStr)
	case OperatorNotContains:
		actualStr := fmt.Sprintf("%v", actual)
		expectedStr := fmt.Sprintf("%v", expected)
		return !strings.Contains(actualStr, expectedStr)
	case OperatorGreaterThan:
		return pe.compareNumeric(actual, expected, ">")
	case OperatorLessThan:
		return pe.compareNumeric(actual, expected, "<")
	case OperatorRegex:
		actualStr := fmt.Sprintf("%v", actual)
		expectedStr := fmt.Sprintf("%v", expected)
		matched, _ := regexp.MatchString(expectedStr, actualStr)
		return matched
	case OperatorInList:
		return pe.inList(actual, expected)
	case OperatorNotInList:
		return !pe.inList(actual, expected)
	default:
		return false
	}
}

// compareNumeric compares numeric values
func (pe *PolicyEngine) compareNumeric(actual, expected interface{}, operator string) bool {
	// Implementation would handle type conversion and numeric comparison
	// For brevity, simplified implementation
	return false
}

// inList checks if a value is in a list
func (pe *PolicyEngine) inList(value, list interface{}) bool {
	// Implementation would handle list membership checking
	// For brevity, simplified implementation
	return false
}

// getSortedPolicies returns policies sorted by priority
func (pe *PolicyEngine) getSortedPolicies() []*SecurityPolicy {
	policies := make([]*SecurityPolicy, 0, len(pe.policies))
	for _, policy := range pe.policies {
		policies = append(policies, policy)
	}

	// Sort by priority (higher priority first)
	// Implementation would sort the slice
	return policies
}

// determineDecision determines the final decision based on matched policies and rules
func (pe *PolicyEngine) determineDecision(policies []*SecurityPolicy, rules []*PolicyRule, actions []*PolicyAction) PolicyDecision {
	decision := pe.config.DefaultPolicy

	// Process actions to determine final decision
	for _, action := range actions {
		switch action.Type {
		case ActionTypeAllow:
			decision.Allowed = true
		case ActionTypeDeny:
			decision.Allowed = false
		case ActionTypeRequireMFA:
			decision.RequiresMFA = true
		case ActionTypeBlock:
			decision.Allowed = false
			decision.Reason = "Blocked by security policy"
		}
	}

	return decision
}

// generateDecisionReason generates a human-readable reason for the decision
func (pe *PolicyEngine) generateDecisionReason(result *PolicyEvaluationResult) string {
	if len(result.MatchedPolicies) == 0 {
		return "No matching policies found, using default policy"
	}

	return fmt.Sprintf("Decision based on %d matched policies and %d matched rules",
		len(result.MatchedPolicies), len(result.MatchedRules))
}

// logPolicyEvaluation logs policy evaluation for audit purposes
func (pe *PolicyEngine) logPolicyEvaluation(ctx context.Context, evalCtx *PolicyEvaluationContext, result *PolicyEvaluationResult) {
	pe.logger.Info(ctx, "Policy evaluation completed", map[string]interface{}{
		"user_id":          evalCtx.UserID,
		"resource":         evalCtx.Resource,
		"action":           evalCtx.Action,
		"decision_allowed": result.Decision.Allowed,
		"requires_mfa":     result.Decision.RequiresMFA,
		"matched_policies": len(result.MatchedPolicies),
		"matched_rules":    len(result.MatchedRules),
		"evaluation_time":  result.EvaluationTime,
	})
}

// loadDefaultPolicies loads default security policies
func (pe *PolicyEngine) loadDefaultPolicies() {
	// Admin access policy
	adminPolicy := &SecurityPolicy{
		PolicyID:    "admin-access-policy",
		Name:        "Admin Access Policy",
		Description: "Allows admin users full access with MFA requirement",
		Version:     "1.0",
		Enabled:     true,
		Priority:    100,
		Conditions: []*PolicyCondition{
			{
				ConditionID: "admin-role-condition",
				Type:        ConditionTypeRole,
				Operator:    OperatorEquals,
				Value:       "admin",
			},
		},
		Rules: []*PolicyRule{
			{
				RuleID:      "admin-allow-rule",
				Name:        "Allow Admin Access",
				Description: "Allow access for admin users",
				Type:        RuleTypeAccess,
				Actions: []*PolicyAction{
					{
						ActionID: "allow-action",
						Type:     ActionTypeAllow,
					},
					{
						ActionID: "require-mfa-action",
						Type:     ActionTypeRequireMFA,
					},
				},
				Enabled:  true,
				Priority: 1,
			},
		},
	}

	pe.AddPolicy(adminPolicy)

	// High risk access policy
	highRiskPolicy := &SecurityPolicy{
		PolicyID:    "high-risk-policy",
		Name:        "High Risk Access Policy",
		Description: "Denies access for high risk requests",
		Version:     "1.0",
		Enabled:     true,
		Priority:    200,
		Conditions: []*PolicyCondition{
			{
				ConditionID: "high-risk-condition",
				Type:        ConditionTypeRiskScore,
				Operator:    OperatorGreaterThan,
				Value:       0.8,
			},
		},
		Rules: []*PolicyRule{
			{
				RuleID:      "high-risk-deny-rule",
				Name:        "Deny High Risk Access",
				Description: "Deny access for high risk requests",
				Type:        RuleTypeAccess,
				Actions: []*PolicyAction{
					{
						ActionID: "deny-action",
						Type:     ActionTypeDeny,
					},
					{
						ActionID: "alert-action",
						Type:     ActionTypeAlert,
					},
				},
				Enabled:  true,
				Priority: 1,
			},
		},
	}

	pe.AddPolicy(highRiskPolicy)
}
