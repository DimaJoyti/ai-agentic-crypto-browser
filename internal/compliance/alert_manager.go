package compliance

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

// AlertManager handles compliance and risk alerts
type AlertManager struct {
	logger      *observability.Logger
	config      ComplianceConfig
	alerts      map[string]*Alert
	rules       map[string]*AlertRule
	channels    map[string]*AlertChannel
	escalations map[string]*EscalationPolicy
	mu          sync.RWMutex
	isRunning   int32
	stopChan    chan struct{}
}

// Alert represents a compliance or risk alert
type Alert struct {
	ID             uuid.UUID              `json:"id"`
	RuleID         string                 `json:"rule_id"`
	Type           AlertType              `json:"type"`
	Severity       AlertSeverity          `json:"severity"`
	Title          string                 `json:"title"`
	Description    string                 `json:"description"`
	Source         string                 `json:"source"`
	Tags           []string               `json:"tags"`
	Metadata       map[string]interface{} `json:"metadata"`
	Timestamp      time.Time              `json:"timestamp"`
	Acknowledged   bool                   `json:"acknowledged"`
	AcknowledgedBy string                 `json:"acknowledged_by,omitempty"`
	AcknowledgedAt *time.Time             `json:"acknowledged_at,omitempty"`
	Resolved       bool                   `json:"resolved"`
	ResolvedBy     string                 `json:"resolved_by,omitempty"`
	ResolvedAt     *time.Time             `json:"resolved_at,omitempty"`
	Escalated      bool                   `json:"escalated"`
	EscalatedAt    *time.Time             `json:"escalated_at,omitempty"`
	Notifications  []AlertNotification    `json:"notifications"`
}

// AlertType defines types of alerts
type AlertType string

const (
	AlertTypeCompliance AlertType = "COMPLIANCE"
	AlertTypeRisk       AlertType = "RISK"
	AlertTypeSecurity   AlertType = "SECURITY"
	AlertTypeSystem     AlertType = "SYSTEM"
	AlertTypeAudit      AlertType = "AUDIT"
)

// AlertRule defines rules for generating alerts
type AlertRule struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Type        AlertType      `json:"type"`
	Severity    AlertSeverity  `json:"severity"`
	Condition   AlertCondition `json:"condition"`
	Actions     []AlertAction  `json:"actions"`
	Enabled     bool           `json:"enabled"`
	Cooldown    time.Duration  `json:"cooldown"`
	LastFired   *time.Time     `json:"last_fired,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// AlertCondition defines when an alert should be triggered
type AlertCondition struct {
	Metric      string            `json:"metric"`
	Operator    ConditionOperator `json:"operator"`
	Threshold   float64           `json:"threshold"`
	Duration    time.Duration     `json:"duration"`
	Aggregation string            `json:"aggregation,omitempty"`
}

// ConditionOperator defines comparison operators for alert conditions
type ConditionOperator string

const (
	OperatorGreaterThan    ConditionOperator = "GT"
	OperatorLessThan       ConditionOperator = "LT"
	OperatorEquals         ConditionOperator = "EQ"
	OperatorNotEquals      ConditionOperator = "NE"
	OperatorGreaterOrEqual ConditionOperator = "GTE"
	OperatorLessOrEqual    ConditionOperator = "LTE"
)

// AlertAction defines actions to take when an alert is triggered
type AlertAction struct {
	Type       ActionType        `json:"type"`
	ChannelID  string            `json:"channel_id"`
	Template   string            `json:"template,omitempty"`
	Parameters map[string]string `json:"parameters,omitempty"`
}

// ActionType defines types of alert actions
type ActionType string

const (
	ActionTypeNotify   ActionType = "NOTIFY"
	ActionTypeEscalate ActionType = "ESCALATE"
	ActionTypeExecute  ActionType = "EXECUTE"
	ActionTypeWebhook  ActionType = "WEBHOOK"
)

// AlertChannel defines notification channels
type AlertChannel struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Type      ChannelType       `json:"type"`
	Config    map[string]string `json:"config"`
	Enabled   bool              `json:"enabled"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// ChannelType defines types of notification channels
type ChannelType string

const (
	ChannelTypeEmail     ChannelType = "EMAIL"
	ChannelTypeSlack     ChannelType = "SLACK"
	ChannelTypeSMS       ChannelType = "SMS"
	ChannelTypeWebhook   ChannelType = "WEBHOOK"
	ChannelTypePagerDuty ChannelType = "PAGERDUTY"
)

// EscalationPolicy defines escalation rules
type EscalationPolicy struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Steps       []EscalationStep `json:"steps"`
	Enabled     bool             `json:"enabled"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// EscalationStep defines a step in an escalation policy
type EscalationStep struct {
	Level      int           `json:"level"`
	Delay      time.Duration `json:"delay"`
	ChannelID  string        `json:"channel_id"`
	Recipients []string      `json:"recipients"`
}

// AlertNotification tracks notifications sent for an alert
type AlertNotification struct {
	ID        uuid.UUID          `json:"id"`
	ChannelID string             `json:"channel_id"`
	Recipient string             `json:"recipient"`
	Status    NotificationStatus `json:"status"`
	SentAt    time.Time          `json:"sent_at"`
	Error     string             `json:"error,omitempty"`
}

// NotificationStatus defines notification delivery status
type NotificationStatus string

const (
	NotificationStatusPending   NotificationStatus = "PENDING"
	NotificationStatusSent      NotificationStatus = "SENT"
	NotificationStatusFailed    NotificationStatus = "FAILED"
	NotificationStatusDelivered NotificationStatus = "DELIVERED"
)

// NewAlertManager creates a new alert manager
func NewAlertManager(logger *observability.Logger, config ComplianceConfig) *AlertManager {
	am := &AlertManager{
		logger:      logger,
		config:      config,
		alerts:      make(map[string]*Alert),
		rules:       make(map[string]*AlertRule),
		channels:    make(map[string]*AlertChannel),
		escalations: make(map[string]*EscalationPolicy),
		stopChan:    make(chan struct{}),
	}

	// Initialize default rules and channels
	am.initializeDefaults()

	return am
}

// Start starts the alert manager
func (am *AlertManager) Start(ctx context.Context) error {
	am.logger.Info(ctx, "Starting alert manager", nil)
	am.isRunning = 1

	// Start monitoring goroutine
	go am.monitoringLoop(ctx)

	return nil
}

// Stop stops the alert manager
func (am *AlertManager) Stop(ctx context.Context) error {
	am.logger.Info(ctx, "Stopping alert manager", nil)
	am.isRunning = 0
	close(am.stopChan)
	return nil
}

// monitoringLoop runs the main monitoring loop
func (am *AlertManager) monitoringLoop(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second) // Check every 10 seconds
	defer ticker.Stop()

	for {
		select {
		case <-am.stopChan:
			return
		case <-ticker.C:
			am.checkEscalations(ctx)
		}
	}
}

// TriggerAlert triggers a new alert
func (am *AlertManager) TriggerAlert(ctx context.Context, ruleID string, metadata map[string]interface{}) error {
	am.mu.RLock()
	rule, exists := am.rules[ruleID]
	am.mu.RUnlock()

	if !exists {
		return fmt.Errorf("alert rule not found: %s", ruleID)
	}

	if !rule.Enabled {
		return nil // Rule is disabled
	}

	// Check cooldown period
	if rule.LastFired != nil && time.Since(*rule.LastFired) < rule.Cooldown {
		return nil // Still in cooldown
	}

	// Create alert
	alert := &Alert{
		ID:            uuid.New(),
		RuleID:        ruleID,
		Type:          rule.Type,
		Severity:      rule.Severity,
		Title:         rule.Name,
		Description:   rule.Description,
		Source:        "compliance_system",
		Metadata:      metadata,
		Timestamp:     time.Now(),
		Notifications: []AlertNotification{},
	}

	// Store alert
	am.mu.Lock()
	am.alerts[alert.ID.String()] = alert
	rule.LastFired = &alert.Timestamp
	am.mu.Unlock()

	// Execute alert actions
	for _, action := range rule.Actions {
		go am.executeAction(ctx, alert, action)
	}

	am.logger.Warn(ctx, "Alert triggered", map[string]interface{}{
		"alert_id": alert.ID,
		"rule_id":  ruleID,
		"type":     alert.Type,
		"severity": alert.Severity,
		"title":    alert.Title,
	})

	return nil
}

// executeAction executes an alert action
func (am *AlertManager) executeAction(ctx context.Context, alert *Alert, action AlertAction) {
	switch action.Type {
	case ActionTypeNotify:
		am.sendNotification(ctx, alert, action.ChannelID)
	case ActionTypeEscalate:
		am.escalateAlert(ctx, alert)
	case ActionTypeWebhook:
		am.sendWebhook(ctx, alert, action.Parameters)
	}
}

// sendNotification sends a notification for an alert
func (am *AlertManager) sendNotification(ctx context.Context, alert *Alert, channelID string) {
	am.mu.RLock()
	channel, exists := am.channels[channelID]
	am.mu.RUnlock()

	if !exists || !channel.Enabled {
		am.logger.Info(ctx, "Alert channel not found or disabled", map[string]interface{}{
			"channel_id": channelID,
		})
		return
	}

	notification := AlertNotification{
		ID:        uuid.New(),
		ChannelID: channelID,
		Status:    NotificationStatusPending,
		SentAt:    time.Now(),
	}

	// Send notification based on channel type
	var err error
	switch channel.Type {
	case ChannelTypeEmail:
		err = am.sendEmail(ctx, alert, channel)
	case ChannelTypeSlack:
		err = am.sendSlack(ctx, alert, channel)
	case ChannelTypeSMS:
		err = am.sendSMS(ctx, alert, channel)
	case ChannelTypeWebhook:
		err = am.sendWebhookNotification(ctx, alert, channel)
	}

	if err != nil {
		notification.Status = NotificationStatusFailed
		notification.Error = err.Error()
		am.logger.Info(ctx, "Failed to send alert notification", map[string]interface{}{
			"alert_id":   alert.ID,
			"channel_id": channelID,
			"error":      err.Error(),
		})
	} else {
		notification.Status = NotificationStatusSent
	}

	// Update alert with notification
	am.mu.Lock()
	alert.Notifications = append(alert.Notifications, notification)
	am.mu.Unlock()
}

// sendEmail sends an email notification (mock implementation)
func (am *AlertManager) sendEmail(ctx context.Context, alert *Alert, channel *AlertChannel) error {
	am.logger.Info(ctx, "Sending email notification", map[string]interface{}{
		"alert_id": alert.ID,
		"to":       channel.Config["to"],
		"subject":  fmt.Sprintf("[%s] %s", alert.Severity, alert.Title),
	})
	// Mock implementation - in production, integrate with email service
	return nil
}

// sendSlack sends a Slack notification (mock implementation)
func (am *AlertManager) sendSlack(ctx context.Context, alert *Alert, channel *AlertChannel) error {
	am.logger.Info(ctx, "Sending Slack notification", map[string]interface{}{
		"alert_id": alert.ID,
		"webhook":  channel.Config["webhook_url"],
	})
	// Mock implementation - in production, integrate with Slack API
	return nil
}

// sendSMS sends an SMS notification (mock implementation)
func (am *AlertManager) sendSMS(ctx context.Context, alert *Alert, channel *AlertChannel) error {
	am.logger.Info(ctx, "Sending SMS notification", map[string]interface{}{
		"alert_id": alert.ID,
		"to":       channel.Config["phone"],
	})
	// Mock implementation - in production, integrate with SMS service
	return nil
}

// sendWebhookNotification sends a webhook notification (mock implementation)
func (am *AlertManager) sendWebhookNotification(ctx context.Context, alert *Alert, channel *AlertChannel) error {
	am.logger.Info(ctx, "Sending webhook notification", map[string]interface{}{
		"alert_id": alert.ID,
		"url":      channel.Config["url"],
	})
	// Mock implementation - in production, make HTTP request
	return nil
}

// sendWebhook sends a webhook for alert actions
func (am *AlertManager) sendWebhook(ctx context.Context, alert *Alert, parameters map[string]string) {
	url := parameters["url"]
	if url == "" {
		return
	}

	am.logger.Info(ctx, "Sending alert webhook", map[string]interface{}{
		"alert_id": alert.ID,
		"url":      url,
	})
	// Mock implementation - in production, make HTTP request
}

// escalateAlert escalates an alert
func (am *AlertManager) escalateAlert(ctx context.Context, alert *Alert) {
	am.mu.Lock()
	alert.Escalated = true
	now := time.Now()
	alert.EscalatedAt = &now
	am.mu.Unlock()

	am.logger.Warn(ctx, "Alert escalated", map[string]interface{}{
		"alert_id": alert.ID,
		"title":    alert.Title,
	})
}

// checkEscalations checks for alerts that need escalation
func (am *AlertManager) checkEscalations(ctx context.Context) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	now := time.Now()
	for _, alert := range am.alerts {
		if alert.Resolved || alert.Escalated {
			continue
		}

		// Check if alert should be escalated (e.g., after 30 minutes)
		if now.Sub(alert.Timestamp) > 30*time.Minute {
			go am.escalateAlert(ctx, alert)
		}
	}
}

// AcknowledgeAlert acknowledges an alert
func (am *AlertManager) AcknowledgeAlert(ctx context.Context, alertID, userID string) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	alert, exists := am.alerts[alertID]
	if !exists {
		return fmt.Errorf("alert not found: %s", alertID)
	}

	if alert.Acknowledged {
		return fmt.Errorf("alert already acknowledged")
	}

	now := time.Now()
	alert.Acknowledged = true
	alert.AcknowledgedBy = userID
	alert.AcknowledgedAt = &now

	am.logger.Info(ctx, "Alert acknowledged", map[string]interface{}{
		"alert_id": alertID,
		"user_id":  userID,
	})

	return nil
}

// ResolveAlert resolves an alert
func (am *AlertManager) ResolveAlert(ctx context.Context, alertID, userID string) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	alert, exists := am.alerts[alertID]
	if !exists {
		return fmt.Errorf("alert not found: %s", alertID)
	}

	if alert.Resolved {
		return fmt.Errorf("alert already resolved")
	}

	now := time.Now()
	alert.Resolved = true
	alert.ResolvedBy = userID
	alert.ResolvedAt = &now

	am.logger.Info(ctx, "Alert resolved", map[string]interface{}{
		"alert_id": alertID,
		"user_id":  userID,
	})

	return nil
}

// GetAlerts returns alerts with optional filtering
func (am *AlertManager) GetAlerts(filter AlertFilter) []*Alert {
	am.mu.RLock()
	defer am.mu.RUnlock()

	var alerts []*Alert
	for _, alert := range am.alerts {
		if am.matchesAlertFilter(alert, filter) {
			alerts = append(alerts, alert)
		}
	}

	return alerts
}

// AlertFilter defines filtering criteria for alerts
type AlertFilter struct {
	Type         AlertType     `json:"type,omitempty"`
	Severity     AlertSeverity `json:"severity,omitempty"`
	Acknowledged *bool         `json:"acknowledged,omitempty"`
	Resolved     *bool         `json:"resolved,omitempty"`
	StartTime    *time.Time    `json:"start_time,omitempty"`
	EndTime      *time.Time    `json:"end_time,omitempty"`
	Limit        int           `json:"limit,omitempty"`
}

// matchesAlertFilter checks if an alert matches the filter criteria
func (am *AlertManager) matchesAlertFilter(alert *Alert, filter AlertFilter) bool {
	if filter.Type != "" && alert.Type != filter.Type {
		return false
	}
	if filter.Severity != "" && alert.Severity != filter.Severity {
		return false
	}
	if filter.Acknowledged != nil && alert.Acknowledged != *filter.Acknowledged {
		return false
	}
	if filter.Resolved != nil && alert.Resolved != *filter.Resolved {
		return false
	}
	if filter.StartTime != nil && alert.Timestamp.Before(*filter.StartTime) {
		return false
	}
	if filter.EndTime != nil && alert.Timestamp.After(*filter.EndTime) {
		return false
	}
	return true
}

// initializeDefaults sets up default alert rules and channels
func (am *AlertManager) initializeDefaults() {
	// Default alert rules
	rules := []*AlertRule{
		{
			ID:          "risk_limit_breach",
			Name:        "Risk Limit Breach",
			Description: "Triggered when a risk limit is breached",
			Type:        AlertTypeRisk,
			Severity:    AlertSeverityError,
			Enabled:     true,
			Cooldown:    5 * time.Minute,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Actions: []AlertAction{
				{
					Type:      ActionTypeNotify,
					ChannelID: "email_risk_team",
				},
			},
		},
		{
			ID:          "compliance_violation",
			Name:        "Compliance Violation",
			Description: "Triggered when a compliance violation is detected",
			Type:        AlertTypeCompliance,
			Severity:    AlertSeverityCritical,
			Enabled:     true,
			Cooldown:    1 * time.Minute,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Actions: []AlertAction{
				{
					Type:      ActionTypeNotify,
					ChannelID: "email_compliance_team",
				},
				{
					Type: ActionTypeEscalate,
				},
			},
		},
	}

	// Default alert channels
	channels := []*AlertChannel{
		{
			ID:   "email_risk_team",
			Name: "Risk Team Email",
			Type: ChannelTypeEmail,
			Config: map[string]string{
				"to":      "risk-team@example.com",
				"from":    "alerts@example.com",
				"subject": "Risk Alert",
			},
			Enabled:   true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:   "email_compliance_team",
			Name: "Compliance Team Email",
			Type: ChannelTypeEmail,
			Config: map[string]string{
				"to":      "compliance-team@example.com",
				"from":    "alerts@example.com",
				"subject": "Compliance Alert",
			},
			Enabled:   true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, rule := range rules {
		am.rules[rule.ID] = rule
	}

	for _, channel := range channels {
		am.channels[channel.ID] = channel
	}
}
