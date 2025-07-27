package alerts

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// AlertService manages real-time alerts and notifications
type AlertService struct {
	logger      *observability.Logger
	config      AlertConfig
	channels    map[string]AlertChannel
	rules       []AlertRule
	subscribers map[string][]chan Alert
	history     []Alert
	mu          sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
}

// AlertConfig holds configuration for the alert service
type AlertConfig struct {
	MaxHistorySize  int           `json:"max_history_size"`
	DefaultCooldown time.Duration `json:"default_cooldown"`
	EnableEmail     bool          `json:"enable_email"`
	EnableWebhook   bool          `json:"enable_webhook"`
	EnableSlack     bool          `json:"enable_slack"`
	EnableTelegram  bool          `json:"enable_telegram"`
	EnablePushNotif bool          `json:"enable_push_notifications"`
}

// AlertChannel interface for different notification channels
type AlertChannel interface {
	Send(ctx context.Context, alert Alert) error
	Name() string
	IsEnabled() bool
}

// AlertRule defines conditions for triggering alerts
type AlertRule struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Metric        string                 `json:"metric"`
	Condition     AlertCondition         `json:"condition"`
	Threshold     decimal.Decimal        `json:"threshold"`
	Severity      AlertSeverity          `json:"severity"`
	Cooldown      time.Duration          `json:"cooldown"`
	Enabled       bool                   `json:"enabled"`
	Channels      []string               `json:"channels"`
	Metadata      map[string]interface{} `json:"metadata"`
	LastTriggered *time.Time             `json:"last_triggered,omitempty"`
}

// AlertCondition represents the condition type for alerts
type AlertCondition string

const (
	ConditionGreaterThan   AlertCondition = "greater_than"
	ConditionLessThan      AlertCondition = "less_than"
	ConditionEquals        AlertCondition = "equals"
	ConditionNotEquals     AlertCondition = "not_equals"
	ConditionPercentChange AlertCondition = "percent_change"
	ConditionMovingAverage AlertCondition = "moving_average"
)

// AlertSeverity represents the severity level of an alert
type AlertSeverity string

const (
	SeverityInfo     AlertSeverity = "info"
	SeverityWarning  AlertSeverity = "warning"
	SeverityError    AlertSeverity = "error"
	SeverityCritical AlertSeverity = "critical"
)

// Alert represents a system alert
type Alert struct {
	ID          string                 `json:"id"`
	RuleID      string                 `json:"rule_id"`
	Title       string                 `json:"title"`
	Message     string                 `json:"message"`
	Severity    AlertSeverity          `json:"severity"`
	Metric      string                 `json:"metric"`
	Value       decimal.Decimal        `json:"value"`
	Threshold   decimal.Decimal        `json:"threshold"`
	Timestamp   time.Time              `json:"timestamp"`
	Resolved    bool                   `json:"resolved"`
	ResolvedAt  *time.Time             `json:"resolved_at,omitempty"`
	Channels    []string               `json:"channels"`
	Metadata    map[string]interface{} `json:"metadata"`
	UserID      *uuid.UUID             `json:"user_id,omitempty"`
	PortfolioID *uuid.UUID             `json:"portfolio_id,omitempty"`
}

// EmailChannel implements email notifications
type EmailChannel struct {
	config EmailConfig
	logger *observability.Logger
}

// EmailConfig holds email configuration
type EmailConfig struct {
	SMTPHost    string   `json:"smtp_host"`
	SMTPPort    int      `json:"smtp_port"`
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	FromAddress string   `json:"from_address"`
	ToAddresses []string `json:"to_addresses"`
	Enabled     bool     `json:"enabled"`
}

// WebhookChannel implements webhook notifications
type WebhookChannel struct {
	config WebhookConfig
	logger *observability.Logger
}

// WebhookConfig holds webhook configuration
type WebhookConfig struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Timeout time.Duration     `json:"timeout"`
	Enabled bool              `json:"enabled"`
}

// SlackChannel implements Slack notifications
type SlackChannel struct {
	config SlackConfig
	logger *observability.Logger
}

// SlackConfig holds Slack configuration
type SlackConfig struct {
	WebhookURL string `json:"webhook_url"`
	Channel    string `json:"channel"`
	Username   string `json:"username"`
	IconEmoji  string `json:"icon_emoji"`
	Enabled    bool   `json:"enabled"`
}

// NewAlertService creates a new alert service
func NewAlertService(logger *observability.Logger, config AlertConfig) *AlertService {
	ctx, cancel := context.WithCancel(context.Background())

	return &AlertService{
		logger:      logger,
		config:      config,
		channels:    make(map[string]AlertChannel),
		rules:       make([]AlertRule, 0),
		subscribers: make(map[string][]chan Alert),
		history:     make([]Alert, 0),
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Start begins the alert service
func (a *AlertService) Start() error {
	a.logger.Info(a.ctx, "Starting alert service", map[string]interface{}{
		"max_history_size": a.config.MaxHistorySize,
		"default_cooldown": a.config.DefaultCooldown.String(),
	})

	// Initialize alert channels
	a.initializeChannels()

	// Load default alert rules
	a.loadDefaultRules()

	return nil
}

// Stop stops the alert service
func (a *AlertService) Stop() error {
	a.logger.Info(a.ctx, "Stopping alert service")
	a.cancel()

	// Close all subscriber channels
	a.mu.Lock()
	defer a.mu.Unlock()

	for topic, channels := range a.subscribers {
		for _, ch := range channels {
			close(ch)
		}
		a.logger.Info(a.ctx, "Closed alert subscribers", map[string]interface{}{
			"topic": topic,
			"count": len(channels),
		})
	}

	return nil
}

// SendAlert sends an alert through configured channels
func (a *AlertService) SendAlert(alert Alert) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Add to history
	a.history = append(a.history, alert)
	if len(a.history) > a.config.MaxHistorySize {
		a.history = a.history[1:]
	}

	// Send to subscribers
	a.notifySubscribers(alert)

	// Send through configured channels
	for _, channelName := range alert.Channels {
		if channel, exists := a.channels[channelName]; exists && channel.IsEnabled() {
			go func(ch AlertChannel, al Alert) {
				if err := ch.Send(a.ctx, al); err != nil {
					a.logger.Error(a.ctx, "Failed to send alert", err, map[string]interface{}{
						"alert_id": al.ID,
						"channel":  ch.Name(),
					})
				}
			}(channel, alert)
		}
	}

	a.logger.Info(a.ctx, "Alert sent", map[string]interface{}{
		"alert_id": alert.ID,
		"severity": string(alert.Severity),
		"title":    alert.Title,
		"channels": alert.Channels,
	})

	return nil
}

// CreateAlert creates a new alert
func (a *AlertService) CreateAlert(ruleID, title, message string, severity AlertSeverity, metric string, value, threshold decimal.Decimal, channels []string) Alert {
	alert := Alert{
		ID:        uuid.New().String(),
		RuleID:    ruleID,
		Title:     title,
		Message:   message,
		Severity:  severity,
		Metric:    metric,
		Value:     value,
		Threshold: threshold,
		Timestamp: time.Now(),
		Resolved:  false,
		Channels:  channels,
		Metadata:  make(map[string]interface{}),
	}

	return alert
}

// Subscribe subscribes to alerts for a specific topic
func (a *AlertService) Subscribe(topic string) <-chan Alert {
	a.mu.Lock()
	defer a.mu.Unlock()

	ch := make(chan Alert, 100)

	if a.subscribers[topic] == nil {
		a.subscribers[topic] = make([]chan Alert, 0)
	}

	a.subscribers[topic] = append(a.subscribers[topic], ch)

	a.logger.Info(a.ctx, "New alert subscriber", map[string]interface{}{
		"topic":       topic,
		"subscribers": len(a.subscribers[topic]),
	})

	return ch
}

// GetAlerts returns recent alerts
func (a *AlertService) GetAlerts(limit int) []Alert {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if limit <= 0 || limit > len(a.history) {
		limit = len(a.history)
	}

	// Return most recent alerts
	start := len(a.history) - limit
	return a.history[start:]
}

// GetActiveAlerts returns unresolved alerts
func (a *AlertService) GetActiveAlerts() []Alert {
	a.mu.RLock()
	defer a.mu.RUnlock()

	activeAlerts := make([]Alert, 0)
	for _, alert := range a.history {
		if !alert.Resolved {
			activeAlerts = append(activeAlerts, alert)
		}
	}

	return activeAlerts
}

// ResolveAlert marks an alert as resolved
func (a *AlertService) ResolveAlert(alertID string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	for i, alert := range a.history {
		if alert.ID == alertID {
			now := time.Now()
			a.history[i].Resolved = true
			a.history[i].ResolvedAt = &now

			a.logger.Info(a.ctx, "Alert resolved", map[string]interface{}{
				"alert_id": alertID,
				"duration": now.Sub(alert.Timestamp).String(),
			})

			return nil
		}
	}

	return fmt.Errorf("alert not found: %s", alertID)
}

// AddRule adds a new alert rule
func (a *AlertService) AddRule(rule AlertRule) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.rules = append(a.rules, rule)

	a.logger.Info(a.ctx, "Alert rule added", map[string]interface{}{
		"rule_id":   rule.ID,
		"rule_name": rule.Name,
		"metric":    rule.Metric,
		"condition": string(rule.Condition),
		"threshold": rule.Threshold.String(),
		"severity":  string(rule.Severity),
	})
}

// CheckRules evaluates alert rules against current metrics
func (a *AlertService) CheckRules(metrics map[string]decimal.Decimal) {
	a.mu.RLock()
	rules := make([]AlertRule, len(a.rules))
	copy(rules, a.rules)
	a.mu.RUnlock()

	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}

		// Check cooldown
		if rule.LastTriggered != nil && time.Since(*rule.LastTriggered) < rule.Cooldown {
			continue
		}

		// Get metric value
		value, exists := metrics[rule.Metric]
		if !exists {
			continue
		}

		// Evaluate condition
		if a.evaluateCondition(rule.Condition, value, rule.Threshold) {
			// Create and send alert
			alert := a.CreateAlert(
				rule.ID,
				rule.Name,
				fmt.Sprintf("%s: %s %s %s", rule.Description, value.String(), string(rule.Condition), rule.Threshold.String()),
				rule.Severity,
				rule.Metric,
				value,
				rule.Threshold,
				rule.Channels,
			)

			a.SendAlert(alert)

			// Update last triggered time
			a.mu.Lock()
			for i := range a.rules {
				if a.rules[i].ID == rule.ID {
					now := time.Now()
					a.rules[i].LastTriggered = &now
					break
				}
			}
			a.mu.Unlock()
		}
	}
}

// evaluateCondition evaluates an alert condition
func (a *AlertService) evaluateCondition(condition AlertCondition, value, threshold decimal.Decimal) bool {
	switch condition {
	case ConditionGreaterThan:
		return value.GreaterThan(threshold)
	case ConditionLessThan:
		return value.LessThan(threshold)
	case ConditionEquals:
		return value.Equal(threshold)
	case ConditionNotEquals:
		return !value.Equal(threshold)
	case ConditionPercentChange:
		// Simplified percent change evaluation
		return value.Sub(threshold).Div(threshold).Abs().GreaterThan(decimal.NewFromFloat(0.1))
	default:
		return false
	}
}

// notifySubscribers notifies all subscribers of an alert
func (a *AlertService) notifySubscribers(alert Alert) {
	// Notify general subscribers
	if subscribers, exists := a.subscribers["all"]; exists {
		for _, ch := range subscribers {
			select {
			case ch <- alert:
			default:
				// Channel is full, skip
			}
		}
	}

	// Notify severity-specific subscribers
	severityTopic := fmt.Sprintf("severity_%s", string(alert.Severity))
	if subscribers, exists := a.subscribers[severityTopic]; exists {
		for _, ch := range subscribers {
			select {
			case ch <- alert:
			default:
				// Channel is full, skip
			}
		}
	}

	// Notify metric-specific subscribers
	metricTopic := fmt.Sprintf("metric_%s", alert.Metric)
	if subscribers, exists := a.subscribers[metricTopic]; exists {
		for _, ch := range subscribers {
			select {
			case ch <- alert:
			default:
				// Channel is full, skip
			}
		}
	}
}

// initializeChannels initializes alert channels
func (a *AlertService) initializeChannels() {
	// Initialize email channel
	if a.config.EnableEmail {
		emailConfig := EmailConfig{
			SMTPHost:    "smtp.gmail.com",
			SMTPPort:    587,
			FromAddress: "alerts@crypto-browser.com",
			Enabled:     true,
		}
		a.channels["email"] = NewEmailChannel(emailConfig, a.logger)
	}

	// Initialize webhook channel
	if a.config.EnableWebhook {
		webhookConfig := WebhookConfig{
			URL:     "https://api.example.com/webhooks/alerts",
			Method:  "POST",
			Headers: map[string]string{"Content-Type": "application/json"},
			Timeout: 10 * time.Second,
			Enabled: true,
		}
		a.channels["webhook"] = NewWebhookChannel(webhookConfig, a.logger)
	}

	// Initialize Slack channel
	if a.config.EnableSlack {
		slackConfig := SlackConfig{
			WebhookURL: "https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK",
			Channel:    "#alerts",
			Username:   "CryptoBrowser",
			IconEmoji:  ":warning:",
			Enabled:    true,
		}
		a.channels["slack"] = NewSlackChannel(slackConfig, a.logger)
	}
}

// loadDefaultRules loads default alert rules
func (a *AlertService) loadDefaultRules() {
	defaultRules := []AlertRule{
		{
			ID:          "high_cpu_usage",
			Name:        "High CPU Usage",
			Description: "CPU usage exceeds threshold",
			Metric:      "cpu_usage_percent",
			Condition:   ConditionGreaterThan,
			Threshold:   decimal.NewFromFloat(80),
			Severity:    SeverityWarning,
			Cooldown:    5 * time.Minute,
			Enabled:     true,
			Channels:    []string{"email", "slack"},
		},
		{
			ID:          "high_memory_usage",
			Name:        "High Memory Usage",
			Description: "Memory usage exceeds threshold",
			Metric:      "memory_usage_percent",
			Condition:   ConditionGreaterThan,
			Threshold:   decimal.NewFromFloat(85),
			Severity:    SeverityWarning,
			Cooldown:    5 * time.Minute,
			Enabled:     true,
			Channels:    []string{"email", "slack"},
		},
		{
			ID:          "high_error_rate",
			Name:        "High Error Rate",
			Description: "Application error rate exceeds threshold",
			Metric:      "error_rate_percent",
			Condition:   ConditionGreaterThan,
			Threshold:   decimal.NewFromFloat(5),
			Severity:    SeverityCritical,
			Cooldown:    2 * time.Minute,
			Enabled:     true,
			Channels:    []string{"email", "slack", "webhook"},
		},
		{
			ID:          "portfolio_loss",
			Name:        "Portfolio Loss Alert",
			Description: "Portfolio loss exceeds threshold",
			Metric:      "portfolio_loss_percent",
			Condition:   ConditionGreaterThan,
			Threshold:   decimal.NewFromFloat(10),
			Severity:    SeverityError,
			Cooldown:    10 * time.Minute,
			Enabled:     true,
			Channels:    []string{"email", "slack"},
		},
	}

	for _, rule := range defaultRules {
		a.AddRule(rule)
	}
}

// Channel implementations

// NewEmailChannel creates a new email channel
func NewEmailChannel(config EmailConfig, logger *observability.Logger) *EmailChannel {
	return &EmailChannel{
		config: config,
		logger: logger,
	}
}

func (e *EmailChannel) Send(ctx context.Context, alert Alert) error {
	// Simplified email sending (would use actual SMTP in production)
	e.logger.Info(ctx, "Email alert sent", map[string]interface{}{
		"alert_id": alert.ID,
		"title":    alert.Title,
		"to":       e.config.ToAddresses,
	})
	return nil
}

func (e *EmailChannel) Name() string {
	return "email"
}

func (e *EmailChannel) IsEnabled() bool {
	return e.config.Enabled
}

// NewWebhookChannel creates a new webhook channel
func NewWebhookChannel(config WebhookConfig, logger *observability.Logger) *WebhookChannel {
	return &WebhookChannel{
		config: config,
		logger: logger,
	}
}

func (w *WebhookChannel) Send(ctx context.Context, alert Alert) error {
	// Simplified webhook sending
	payload, _ := json.Marshal(alert)

	w.logger.Info(ctx, "Webhook alert sent", map[string]interface{}{
		"alert_id": alert.ID,
		"url":      w.config.URL,
		"payload":  string(payload),
	})

	// In production, would make actual HTTP request
	return nil
}

func (w *WebhookChannel) Name() string {
	return "webhook"
}

func (w *WebhookChannel) IsEnabled() bool {
	return w.config.Enabled
}

// NewSlackChannel creates a new Slack channel
func NewSlackChannel(config SlackConfig, logger *observability.Logger) *SlackChannel {
	return &SlackChannel{
		config: config,
		logger: logger,
	}
}

func (s *SlackChannel) Send(ctx context.Context, alert Alert) error {
	// Simplified Slack message formatting
	color := "warning"
	switch alert.Severity {
	case SeverityCritical:
		color = "danger"
	case SeverityError:
		color = "danger"
	case SeverityWarning:
		color = "warning"
	case SeverityInfo:
		color = "good"
	}

	message := map[string]interface{}{
		"channel":    s.config.Channel,
		"username":   s.config.Username,
		"icon_emoji": s.config.IconEmoji,
		"attachments": []map[string]interface{}{
			{
				"color": color,
				"title": alert.Title,
				"text":  alert.Message,
				"fields": []map[string]interface{}{
					{"title": "Severity", "value": string(alert.Severity), "short": true},
					{"title": "Metric", "value": alert.Metric, "short": true},
					{"title": "Value", "value": alert.Value.String(), "short": true},
					{"title": "Threshold", "value": alert.Threshold.String(), "short": true},
				},
				"ts": alert.Timestamp.Unix(),
			},
		},
	}

	payload, _ := json.Marshal(message)

	s.logger.Info(ctx, "Slack alert sent", map[string]interface{}{
		"alert_id": alert.ID,
		"channel":  s.config.Channel,
		"payload":  string(payload),
	})

	// In production, would make actual HTTP request to Slack webhook
	return nil
}

func (s *SlackChannel) Name() string {
	return "slack"
}

func (s *SlackChannel) IsEnabled() bool {
	return s.config.Enabled
}
