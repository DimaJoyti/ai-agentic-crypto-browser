package analytics

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

// AlertManager manages intelligent alerting
type AlertManager struct {
	logger        *observability.Logger
	config        *AnalyticsConfig
	alertRules    map[string]*AlertRule
	activeAlerts  map[string]*Alert
	alertHistory  []*Alert
	notifications chan *AlertNotification
	escalations   map[string]*EscalationPolicy
	suppressions  map[string]*AlertSuppression
	mu            sync.RWMutex
}

// AlertRule defines an alert rule
type AlertRule struct {
	RuleID           string            `json:"rule_id"`
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	MetricName       string            `json:"metric_name"`
	Condition        AlertCondition    `json:"condition"`
	Threshold        float64           `json:"threshold"`
	Severity         AlertSeverity     `json:"severity"`
	Duration         time.Duration     `json:"duration"`
	EvaluationWindow time.Duration     `json:"evaluation_window"`
	Enabled          bool              `json:"enabled"`
	Tags             map[string]string `json:"tags"`
	Actions          []AlertAction     `json:"actions"`
	EscalationPolicy string            `json:"escalation_policy,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	LastTriggered    *time.Time        `json:"last_triggered,omitempty"`
}

// AlertCondition defines alert conditions
type AlertCondition string

const (
	ConditionGreaterThan      AlertCondition = "greater_than"
	ConditionLessThan         AlertCondition = "less_than"
	ConditionEquals           AlertCondition = "equals"
	ConditionNotEquals        AlertCondition = "not_equals"
	ConditionAnomalyDetected  AlertCondition = "anomaly_detected"
	ConditionTrendChange      AlertCondition = "trend_change"
	ConditionPredictionFailed AlertCondition = "prediction_failed"
)

// AlertSeverity defines alert severity levels
type AlertSeverity string

const (
	SeverityInfo     AlertSeverity = "info"
	SeverityWarning  AlertSeverity = "warning"
	SeverityError    AlertSeverity = "error"
	SeverityCritical AlertSeverity = "critical"
)

// AlertAction defines actions to take when alert triggers
type AlertAction struct {
	ActionType AlertActionType        `json:"action_type"`
	Target     string                 `json:"target"`
	Parameters map[string]interface{} `json:"parameters"`
	Enabled    bool                   `json:"enabled"`
}

// AlertActionType defines types of alert actions
type AlertActionType string

const (
	ActionTypeEmail         AlertActionType = "email"
	ActionTypeSlack         AlertActionType = "slack"
	ActionTypeWebhook       AlertActionType = "webhook"
	ActionTypeSMS           AlertActionType = "sms"
	ActionTypePagerDuty     AlertActionType = "pagerduty"
	ActionTypeAutoRemediate AlertActionType = "auto_remediate"
)

// Alert represents an active alert
type Alert struct {
	AlertID           string                 `json:"alert_id"`
	RuleID            string                 `json:"rule_id"`
	RuleName          string                 `json:"rule_name"`
	MetricName        string                 `json:"metric_name"`
	Severity          AlertSeverity          `json:"severity"`
	Status            AlertStatus            `json:"status"`
	Message           string                 `json:"message"`
	Value             float64                `json:"value"`
	Threshold         float64                `json:"threshold"`
	TriggeredAt       time.Time              `json:"triggered_at"`
	ResolvedAt        *time.Time             `json:"resolved_at,omitempty"`
	AcknowledgedAt    *time.Time             `json:"acknowledged_at,omitempty"`
	AcknowledgedBy    string                 `json:"acknowledged_by,omitempty"`
	Context           map[string]interface{} `json:"context"`
	Tags              map[string]string      `json:"tags"`
	Escalated         bool                   `json:"escalated"`
	EscalationLevel   int                    `json:"escalation_level"`
	NotificationsSent []AlertNotification    `json:"notifications_sent"`
}

// AlertStatus defines alert status
type AlertStatus string

const (
	AlertStatusActive       AlertStatus = "active"
	AlertStatusAcknowledged AlertStatus = "acknowledged"
	AlertStatusResolved     AlertStatus = "resolved"
	AlertStatusSuppressed   AlertStatus = "suppressed"
)

// AlertNotification represents a notification to be sent
type AlertNotification struct {
	NotificationID string                 `json:"notification_id"`
	AlertID        string                 `json:"alert_id"`
	ActionType     AlertActionType        `json:"action_type"`
	Target         string                 `json:"target"`
	Message        string                 `json:"message"`
	Priority       NotificationPriority   `json:"priority"`
	ScheduledAt    time.Time              `json:"scheduled_at"`
	SentAt         *time.Time             `json:"sent_at,omitempty"`
	Status         NotificationStatus     `json:"status"`
	Retries        int                    `json:"retries"`
	MaxRetries     int                    `json:"max_retries"`
	Error          string                 `json:"error,omitempty"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// NotificationPriority defines notification priority
type NotificationPriority string

const (
	PriorityLow      NotificationPriority = "low"
	PriorityMedium   NotificationPriority = "medium"
	PriorityHigh     NotificationPriority = "high"
	PriorityCritical NotificationPriority = "critical"
)

// NotificationStatus defines notification status
type NotificationStatus string

const (
	NotificationStatusPending  NotificationStatus = "pending"
	NotificationStatusSent     NotificationStatus = "sent"
	NotificationStatusFailed   NotificationStatus = "failed"
	NotificationStatusRetrying NotificationStatus = "retrying"
)

// EscalationPolicy defines escalation rules
type EscalationPolicy struct {
	PolicyID    string           `json:"policy_id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Steps       []EscalationStep `json:"steps"`
	Enabled     bool             `json:"enabled"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// EscalationStep defines a step in escalation
type EscalationStep struct {
	StepNumber int           `json:"step_number"`
	Delay      time.Duration `json:"delay"`
	Actions    []AlertAction `json:"actions"`
	Condition  string        `json:"condition,omitempty"`
}

// AlertSuppression defines alert suppression rules
type AlertSuppression struct {
	SuppressionID string                 `json:"suppression_id"`
	Name          string                 `json:"name"`
	RuleIDs       []string               `json:"rule_ids"`
	MetricNames   []string               `json:"metric_names"`
	Tags          map[string]string      `json:"tags"`
	StartTime     time.Time              `json:"start_time"`
	EndTime       time.Time              `json:"end_time"`
	Reason        string                 `json:"reason"`
	CreatedBy     string                 `json:"created_by"`
	Active        bool                   `json:"active"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// NewAlertManager creates a new alert manager
func NewAlertManager(logger *observability.Logger, config *AnalyticsConfig) *AlertManager {
	return &AlertManager{
		logger:        logger,
		config:        config,
		alertRules:    make(map[string]*AlertRule),
		activeAlerts:  make(map[string]*Alert),
		alertHistory:  make([]*Alert, 0),
		notifications: make(chan *AlertNotification, 1000),
		escalations:   make(map[string]*EscalationPolicy),
		suppressions:  make(map[string]*AlertSuppression),
	}
}

// Start starts the alert manager
func (am *AlertManager) Start(ctx context.Context) error {
	am.logger.Info(ctx, "Starting alert manager", nil)

	// Initialize default alert rules
	am.initializeDefaultRules()

	// Start background processes
	go am.processNotifications(ctx)
	go am.processEscalations(ctx)
	go am.cleanupOldAlerts(ctx)

	return nil
}

// CreateAlertRule creates a new alert rule
func (am *AlertManager) CreateAlertRule(rule *AlertRule) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	if rule.RuleID == "" {
		rule.RuleID = uuid.New().String()
	}

	rule.CreatedAt = time.Now()
	rule.UpdatedAt = time.Now()

	am.alertRules[rule.RuleID] = rule

	am.logger.Info(context.Background(), "Alert rule created", map[string]interface{}{
		"rule_id":     rule.RuleID,
		"name":        rule.Name,
		"metric_name": rule.MetricName,
		"condition":   rule.Condition,
		"threshold":   rule.Threshold,
		"severity":    rule.Severity,
	})

	return nil
}

// EvaluateMetric evaluates a metric against all applicable alert rules
func (am *AlertManager) EvaluateMetric(metricName string, value float64, tags map[string]string) {
	am.mu.RLock()
	applicableRules := make([]*AlertRule, 0)
	for _, rule := range am.alertRules {
		if rule.Enabled && rule.MetricName == metricName {
			if am.matchesTags(rule.Tags, tags) {
				applicableRules = append(applicableRules, rule)
			}
		}
	}
	am.mu.RUnlock()

	for _, rule := range applicableRules {
		am.evaluateRule(rule, value, tags)
	}
}

// evaluateRule evaluates a single rule
func (am *AlertManager) evaluateRule(rule *AlertRule, value float64, tags map[string]string) {
	triggered := false

	switch rule.Condition {
	case ConditionGreaterThan:
		triggered = value > rule.Threshold
	case ConditionLessThan:
		triggered = value < rule.Threshold
	case ConditionEquals:
		triggered = value == rule.Threshold
	case ConditionNotEquals:
		triggered = value != rule.Threshold
	}

	if triggered {
		am.triggerAlert(rule, value, tags)
	} else {
		am.resolveAlert(rule.RuleID)
	}
}

// triggerAlert triggers an alert
func (am *AlertManager) triggerAlert(rule *AlertRule, value float64, tags map[string]string) {
	am.mu.Lock()
	defer am.mu.Unlock()

	// Check if alert already exists
	alertKey := fmt.Sprintf("%s_%s", rule.RuleID, rule.MetricName)
	if existingAlert, exists := am.activeAlerts[alertKey]; exists {
		// Update existing alert
		existingAlert.Value = value
		existingAlert.Context["last_updated"] = time.Now()
		return
	}

	// Check if suppressed
	if am.isAlertSuppressed(rule, tags) {
		return
	}

	// Create new alert
	alert := &Alert{
		AlertID:     uuid.New().String(),
		RuleID:      rule.RuleID,
		RuleName:    rule.Name,
		MetricName:  rule.MetricName,
		Severity:    rule.Severity,
		Status:      AlertStatusActive,
		Message:     am.generateAlertMessage(rule, value),
		Value:       value,
		Threshold:   rule.Threshold,
		TriggeredAt: time.Now(),
		Context: map[string]interface{}{
			"rule_description":  rule.Description,
			"condition":         rule.Condition,
			"evaluation_window": rule.EvaluationWindow,
		},
		Tags:              tags,
		Escalated:         false,
		EscalationLevel:   0,
		NotificationsSent: make([]AlertNotification, 0),
	}

	am.activeAlerts[alertKey] = alert
	am.alertHistory = append(am.alertHistory, alert)

	// Update rule last triggered time
	now := time.Now()
	rule.LastTriggered = &now

	// Send notifications
	am.sendAlertNotifications(alert, rule)

	am.logger.Warn(context.Background(), "Alert triggered", map[string]interface{}{
		"alert_id":    alert.AlertID,
		"rule_name":   rule.Name,
		"metric_name": rule.MetricName,
		"value":       value,
		"threshold":   rule.Threshold,
		"severity":    rule.Severity,
	})
}

// resolveAlert resolves an alert
func (am *AlertManager) resolveAlert(ruleID string) {
	am.mu.Lock()
	defer am.mu.Unlock()

	for key, alert := range am.activeAlerts {
		if alert.RuleID == ruleID && alert.Status == AlertStatusActive {
			alert.Status = AlertStatusResolved
			now := time.Now()
			alert.ResolvedAt = &now
			delete(am.activeAlerts, key)

			am.logger.Info(context.Background(), "Alert resolved", map[string]interface{}{
				"alert_id":  alert.AlertID,
				"rule_name": alert.RuleName,
				"duration":  alert.ResolvedAt.Sub(alert.TriggeredAt),
			})
			break
		}
	}
}

// sendAlertNotifications sends notifications for an alert
func (am *AlertManager) sendAlertNotifications(alert *Alert, rule *AlertRule) {
	for _, action := range rule.Actions {
		if !action.Enabled {
			continue
		}

		notification := &AlertNotification{
			NotificationID: uuid.New().String(),
			AlertID:        alert.AlertID,
			ActionType:     action.ActionType,
			Target:         action.Target,
			Message:        am.formatNotificationMessage(alert, action),
			Priority:       am.severityToPriority(alert.Severity),
			ScheduledAt:    time.Now(),
			Status:         NotificationStatusPending,
			MaxRetries:     3,
			Metadata:       action.Parameters,
		}

		select {
		case am.notifications <- notification:
			alert.NotificationsSent = append(alert.NotificationsSent, *notification)
		default:
			am.logger.Warn(context.Background(), "Notification queue full", map[string]interface{}{
				"alert_id":    alert.AlertID,
				"action_type": action.ActionType,
			})
		}
	}
}

// generateAlertMessage generates an alert message
func (am *AlertManager) generateAlertMessage(rule *AlertRule, value float64) string {
	return fmt.Sprintf("Alert: %s - %s %s %.2f (threshold: %.2f)",
		rule.Name, rule.MetricName, rule.Condition, value, rule.Threshold)
}

// formatNotificationMessage formats a notification message
func (am *AlertManager) formatNotificationMessage(alert *Alert, action AlertAction) string {
	return fmt.Sprintf("🚨 %s Alert: %s\n\nMetric: %s\nValue: %.2f\nThreshold: %.2f\nTriggered: %s",
		alert.Severity, alert.RuleName, alert.MetricName, alert.Value, alert.Threshold, alert.TriggeredAt.Format(time.RFC3339))
}

// severityToPriority converts alert severity to notification priority
func (am *AlertManager) severityToPriority(severity AlertSeverity) NotificationPriority {
	switch severity {
	case SeverityInfo:
		return PriorityLow
	case SeverityWarning:
		return PriorityMedium
	case SeverityError:
		return PriorityHigh
	case SeverityCritical:
		return PriorityCritical
	default:
		return PriorityMedium
	}
}

// matchesTags checks if rule tags match metric tags
func (am *AlertManager) matchesTags(ruleTags, metricTags map[string]string) bool {
	if len(ruleTags) == 0 {
		return true // No tag filter
	}

	for key, value := range ruleTags {
		if metricValue, exists := metricTags[key]; !exists || metricValue != value {
			return false
		}
	}

	return true
}

// isAlertSuppressed checks if an alert is suppressed
func (am *AlertManager) isAlertSuppressed(rule *AlertRule, tags map[string]string) bool {
	now := time.Now()

	for _, suppression := range am.suppressions {
		if !suppression.Active || now.Before(suppression.StartTime) || now.After(suppression.EndTime) {
			continue
		}

		// Check rule IDs
		for _, ruleID := range suppression.RuleIDs {
			if ruleID == rule.RuleID {
				return true
			}
		}

		// Check metric names
		for _, metricName := range suppression.MetricNames {
			if metricName == rule.MetricName {
				return true
			}
		}

		// Check tags
		if am.matchesTags(suppression.Tags, tags) {
			return true
		}
	}

	return false
}

// initializeDefaultRules initializes default alert rules
func (am *AlertManager) initializeDefaultRules() {
	defaultRules := []*AlertRule{
		{
			Name:             "High CPU Usage",
			Description:      "CPU usage is above 80%",
			MetricName:       "cpu_usage",
			Condition:        ConditionGreaterThan,
			Threshold:        80.0,
			Severity:         SeverityWarning,
			Duration:         5 * time.Minute,
			EvaluationWindow: 1 * time.Minute,
			Enabled:          true,
			Actions: []AlertAction{
				{
					ActionType: ActionTypeEmail,
					Target:     "admin@example.com",
					Enabled:    true,
				},
			},
		},
		{
			Name:             "High Memory Usage",
			Description:      "Memory usage is above 90%",
			MetricName:       "memory_usage",
			Condition:        ConditionGreaterThan,
			Threshold:        90.0,
			Severity:         SeverityError,
			Duration:         3 * time.Minute,
			EvaluationWindow: 1 * time.Minute,
			Enabled:          true,
			Actions: []AlertAction{
				{
					ActionType: ActionTypeEmail,
					Target:     "admin@example.com",
					Enabled:    true,
				},
			},
		},
		{
			Name:             "High Error Rate",
			Description:      "Error rate is above 5%",
			MetricName:       "error_rate",
			Condition:        ConditionGreaterThan,
			Threshold:        5.0,
			Severity:         SeverityCritical,
			Duration:         2 * time.Minute,
			EvaluationWindow: 30 * time.Second,
			Enabled:          true,
			Actions: []AlertAction{
				{
					ActionType: ActionTypeEmail,
					Target:     "admin@example.com",
					Enabled:    true,
				},
			},
		},
	}

	for _, rule := range defaultRules {
		am.CreateAlertRule(rule)
	}
}

// processNotifications processes notification queue
func (am *AlertManager) processNotifications(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case notification := <-am.notifications:
			am.sendNotification(notification)
		}
	}
}

// sendNotification sends a notification
func (am *AlertManager) sendNotification(notification *AlertNotification) {
	am.logger.Info(context.Background(), "Sending notification", map[string]interface{}{
		"notification_id": notification.NotificationID,
		"alert_id":        notification.AlertID,
		"action_type":     notification.ActionType,
		"target":          notification.Target,
		"priority":        notification.Priority,
	})

	// Simulate sending notification
	switch notification.ActionType {
	case ActionTypeEmail:
		am.sendEmailNotification(notification)
	case ActionTypeSlack:
		am.sendSlackNotification(notification)
	case ActionTypeWebhook:
		am.sendWebhookNotification(notification)
	case ActionTypeSMS:
		am.sendSMSNotification(notification)
	default:
		am.logger.Warn(context.Background(), "Unsupported notification type", map[string]interface{}{
			"action_type": notification.ActionType,
		})
		notification.Status = NotificationStatusFailed
		notification.Error = "Unsupported notification type"
		return
	}

	notification.Status = NotificationStatusSent
	now := time.Now()
	notification.SentAt = &now
}

// sendEmailNotification sends an email notification
func (am *AlertManager) sendEmailNotification(notification *AlertNotification) {
	// Simulate email sending
	am.logger.Info(context.Background(), "Email notification sent", map[string]interface{}{
		"target":  notification.Target,
		"message": notification.Message,
	})
}

// sendSlackNotification sends a Slack notification
func (am *AlertManager) sendSlackNotification(notification *AlertNotification) {
	// Simulate Slack sending
	am.logger.Info(context.Background(), "Slack notification sent", map[string]interface{}{
		"target":  notification.Target,
		"message": notification.Message,
	})
}

// sendWebhookNotification sends a webhook notification
func (am *AlertManager) sendWebhookNotification(notification *AlertNotification) {
	// Simulate webhook sending
	am.logger.Info(context.Background(), "Webhook notification sent", map[string]interface{}{
		"target":  notification.Target,
		"message": notification.Message,
	})
}

// sendSMSNotification sends an SMS notification
func (am *AlertManager) sendSMSNotification(notification *AlertNotification) {
	// Simulate SMS sending
	am.logger.Info(context.Background(), "SMS notification sent", map[string]interface{}{
		"target":  notification.Target,
		"message": notification.Message,
	})
}

// processEscalations processes alert escalations
func (am *AlertManager) processEscalations(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			am.checkEscalations()
		}
	}
}

// checkEscalations checks for alerts that need escalation
func (am *AlertManager) checkEscalations() {
	am.mu.RLock()
	alertsToEscalate := make([]*Alert, 0)
	for _, alert := range am.activeAlerts {
		if alert.Status == AlertStatusActive && !alert.Escalated {
			// Check if alert should be escalated (e.g., after 15 minutes)
			if time.Since(alert.TriggeredAt) > 15*time.Minute {
				alertsToEscalate = append(alertsToEscalate, alert)
			}
		}
	}
	am.mu.RUnlock()

	for _, alert := range alertsToEscalate {
		am.escalateAlert(alert)
	}
}

// escalateAlert escalates an alert
func (am *AlertManager) escalateAlert(alert *Alert) {
	am.mu.Lock()
	defer am.mu.Unlock()

	alert.Escalated = true
	alert.EscalationLevel++

	am.logger.Warn(context.Background(), "Alert escalated", map[string]interface{}{
		"alert_id":         alert.AlertID,
		"rule_name":        alert.RuleName,
		"escalation_level": alert.EscalationLevel,
	})

	// Send escalation notifications
	escalationNotification := &AlertNotification{
		NotificationID: uuid.New().String(),
		AlertID:        alert.AlertID,
		ActionType:     ActionTypeEmail,
		Target:         "escalation@example.com",
		Message:        fmt.Sprintf("🚨 ESCALATED: %s (Level %d)", alert.Message, alert.EscalationLevel),
		Priority:       PriorityCritical,
		ScheduledAt:    time.Now(),
		Status:         NotificationStatusPending,
		MaxRetries:     3,
	}

	select {
	case am.notifications <- escalationNotification:
		alert.NotificationsSent = append(alert.NotificationsSent, *escalationNotification)
	default:
		am.logger.Warn(context.Background(), "Escalation notification queue full", map[string]interface{}{
			"alert_id": alert.AlertID,
		})
	}
}

// cleanupOldAlerts cleans up old resolved alerts
func (am *AlertManager) cleanupOldAlerts(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			am.performAlertCleanup()
		}
	}
}

// performAlertCleanup removes old alerts from history
func (am *AlertManager) performAlertCleanup() {
	am.mu.Lock()
	defer am.mu.Unlock()

	cutoffTime := time.Now().Add(-7 * 24 * time.Hour) // Keep for 7 days
	filteredHistory := make([]*Alert, 0)

	for _, alert := range am.alertHistory {
		if alert.Status == AlertStatusActive ||
			(alert.ResolvedAt != nil && alert.ResolvedAt.After(cutoffTime)) ||
			alert.TriggeredAt.After(cutoffTime) {
			filteredHistory = append(filteredHistory, alert)
		}
	}

	removed := len(am.alertHistory) - len(filteredHistory)
	am.alertHistory = filteredHistory

	if removed > 0 {
		am.logger.Info(context.Background(), "Cleaned up old alerts", map[string]interface{}{
			"removed_count": removed,
		})
	}
}

// GetActiveAlerts returns all active alerts
func (am *AlertManager) GetActiveAlerts() []*Alert {
	am.mu.RLock()
	defer am.mu.RUnlock()

	alerts := make([]*Alert, 0, len(am.activeAlerts))
	for _, alert := range am.activeAlerts {
		alerts = append(alerts, alert)
	}

	return alerts
}

// GetAlertHistory returns alert history
func (am *AlertManager) GetAlertHistory(limit int) []*Alert {
	am.mu.RLock()
	defer am.mu.RUnlock()

	if limit <= 0 || limit > len(am.alertHistory) {
		limit = len(am.alertHistory)
	}

	// Return most recent alerts
	start := len(am.alertHistory) - limit
	return am.alertHistory[start:]
}

// AcknowledgeAlert acknowledges an alert
func (am *AlertManager) AcknowledgeAlert(alertID, acknowledgedBy string) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	for _, alert := range am.activeAlerts {
		if alert.AlertID == alertID {
			alert.Status = AlertStatusAcknowledged
			now := time.Now()
			alert.AcknowledgedAt = &now
			alert.AcknowledgedBy = acknowledgedBy

			am.logger.Info(context.Background(), "Alert acknowledged", map[string]interface{}{
				"alert_id":        alertID,
				"acknowledged_by": acknowledgedBy,
			})

			return nil
		}
	}

	return fmt.Errorf("alert not found: %s", alertID)
}

// CreateSuppression creates an alert suppression
func (am *AlertManager) CreateSuppression(suppression *AlertSuppression) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	if suppression.SuppressionID == "" {
		suppression.SuppressionID = uuid.New().String()
	}

	am.suppressions[suppression.SuppressionID] = suppression

	am.logger.Info(context.Background(), "Alert suppression created", map[string]interface{}{
		"suppression_id": suppression.SuppressionID,
		"name":           suppression.Name,
		"start_time":     suppression.StartTime,
		"end_time":       suppression.EndTime,
		"reason":         suppression.Reason,
	})

	return nil
}

// RemoveSuppression removes an alert suppression
func (am *AlertManager) RemoveSuppression(suppressionID string) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	if _, exists := am.suppressions[suppressionID]; !exists {
		return fmt.Errorf("suppression not found: %s", suppressionID)
	}

	delete(am.suppressions, suppressionID)

	am.logger.Info(context.Background(), "Alert suppression removed", map[string]interface{}{
		"suppression_id": suppressionID,
	})

	return nil
}
