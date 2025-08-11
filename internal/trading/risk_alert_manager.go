package trading

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

// RiskAlertManager manages risk alerts and notifications
type RiskAlertManager struct {
	logger   *observability.Logger
	alerts   map[string]*RiskAlert
	channels []AlertChannel
	config   *AlertConfig
	mu       sync.RWMutex
}

// AlertConfig holds configuration for alert management
type AlertConfig struct {
	MaxAlertsPerHour    int           `yaml:"max_alerts_per_hour"`
	AlertRetention      time.Duration `yaml:"alert_retention"`
	EnableEmailAlerts   bool          `yaml:"enable_email_alerts"`
	EnableSlackAlerts   bool          `yaml:"enable_slack_alerts"`
	EnableWebhookAlerts bool          `yaml:"enable_webhook_alerts"`
}

// Use types from advanced_risk_manager.go to avoid duplication

// AlertChannel interface for different alert delivery methods
type AlertChannel interface {
	SendAlert(ctx context.Context, alert *RiskAlert) error
	GetType() string
	IsEnabled() bool
}

// NewRiskAlertManager creates a new risk alert manager
func NewRiskAlertManager(logger *observability.Logger) *RiskAlertManager {
	config := &AlertConfig{
		MaxAlertsPerHour:    100,
		AlertRetention:      24 * time.Hour,
		EnableEmailAlerts:   true,
		EnableSlackAlerts:   false,
		EnableWebhookAlerts: false,
	}

	return &RiskAlertManager{
		logger:   logger,
		alerts:   make(map[string]*RiskAlert),
		channels: make([]AlertChannel, 0),
		config:   config,
	}
}

// SendAlert sends a risk alert
func (ram *RiskAlertManager) SendAlert(ctx context.Context, alert *RiskAlert) error {
	ram.mu.Lock()
	defer ram.mu.Unlock()

	// Generate ID if not provided
	if alert.ID == "" {
		alert.ID = uuid.New().String()
	}

	// Set timestamps
	alert.CreatedAt = time.Now()
	alert.UpdatedAt = time.Now()
	alert.Status = AlertStatusActive

	// Store alert
	ram.alerts[alert.ID] = alert

	// Log alert
	ram.logger.Error(ctx, "Risk alert generated", fmt.Errorf(alert.Message), map[string]interface{}{
		"alert_id":   alert.ID,
		"alert_type": string(alert.Type),
		"severity":   string(alert.Severity),
		"bot_id":     alert.BotID,
		"symbol":     alert.Symbol,
	})

	// Send through all enabled channels
	for _, channel := range ram.channels {
		if channel.IsEnabled() {
			go func(ch AlertChannel) {
				if err := ch.SendAlert(ctx, alert); err != nil {
					ram.logger.Error(ctx, "Failed to send alert through channel", err, map[string]interface{}{
						"alert_id":     alert.ID,
						"channel_type": ch.GetType(),
					})
				}
			}(channel)
		}
	}

	return nil
}

// AcknowledgeAlert acknowledges an alert
func (ram *RiskAlertManager) AcknowledgeAlert(ctx context.Context, alertID string, userID string) error {
	ram.mu.Lock()
	defer ram.mu.Unlock()

	alert, exists := ram.alerts[alertID]
	if !exists {
		return fmt.Errorf("alert not found: %s", alertID)
	}

	alert.Status = AlertStatusAcknowledged
	alert.UpdatedAt = time.Now()

	ram.logger.Info(ctx, "Alert acknowledged", map[string]interface{}{
		"alert_id": alertID,
		"user_id":  userID,
	})

	return nil
}

// ResolveAlert resolves an alert
func (ram *RiskAlertManager) ResolveAlert(ctx context.Context, alertID string, userID string) error {
	ram.mu.Lock()
	defer ram.mu.Unlock()

	alert, exists := ram.alerts[alertID]
	if !exists {
		return fmt.Errorf("alert not found: %s", alertID)
	}

	alert.Status = AlertStatusResolved
	alert.UpdatedAt = time.Now()

	ram.logger.Info(ctx, "Alert resolved", map[string]interface{}{
		"alert_id": alertID,
		"user_id":  userID,
	})

	return nil
}

// GetActiveAlerts returns all active alerts
func (ram *RiskAlertManager) GetActiveAlerts() []*RiskAlert {
	ram.mu.RLock()
	defer ram.mu.RUnlock()

	var activeAlerts []*RiskAlert
	for _, alert := range ram.alerts {
		if alert.Status == AlertStatusActive {
			activeAlerts = append(activeAlerts, alert)
		}
	}

	return activeAlerts
}

// GetAlertsByBot returns alerts for a specific bot
func (ram *RiskAlertManager) GetAlertsByBot(botID string) []*RiskAlert {
	ram.mu.RLock()
	defer ram.mu.RUnlock()

	var botAlerts []*RiskAlert
	for _, alert := range ram.alerts {
		if alert.BotID == botID {
			botAlerts = append(botAlerts, alert)
		}
	}

	return botAlerts
}

// GetAlertsBySeverity returns alerts by severity level
func (ram *RiskAlertManager) GetAlertsBySeverity(severity AlertSeverity) []*RiskAlert {
	ram.mu.RLock()
	defer ram.mu.RUnlock()

	var severityAlerts []*RiskAlert
	for _, alert := range ram.alerts {
		if alert.Severity == severity {
			severityAlerts = append(severityAlerts, alert)
		}
	}

	return severityAlerts
}

// AddAlertChannel adds an alert delivery channel
func (ram *RiskAlertManager) AddAlertChannel(channel AlertChannel) {
	ram.mu.Lock()
	defer ram.mu.Unlock()

	ram.channels = append(ram.channels, channel)

	ram.logger.Info(context.Background(), "Alert channel added", map[string]interface{}{
		"channel_type": channel.GetType(),
		"enabled":      channel.IsEnabled(),
	})
}

// CleanupOldAlerts removes old alerts based on retention policy
func (ram *RiskAlertManager) CleanupOldAlerts() {
	ram.mu.Lock()
	defer ram.mu.Unlock()

	cutoff := time.Now().Add(-ram.config.AlertRetention)

	for id, alert := range ram.alerts {
		if alert.CreatedAt.Before(cutoff) && alert.Status == AlertStatusResolved {
			delete(ram.alerts, id)
		}
	}
}

// GetAlertStatistics returns alert statistics
func (ram *RiskAlertManager) GetAlertStatistics() map[string]interface{} {
	ram.mu.RLock()
	defer ram.mu.RUnlock()

	stats := map[string]interface{}{
		"total_alerts":        len(ram.alerts),
		"active_alerts":       0,
		"acknowledged_alerts": 0,
		"resolved_alerts":     0,
		"by_severity":         make(map[AlertSeverity]int),
		"by_type":             make(map[RiskAlertType]int),
	}

	for _, alert := range ram.alerts {
		switch alert.Status {
		case AlertStatusActive:
			stats["active_alerts"] = stats["active_alerts"].(int) + 1
		case AlertStatusAcknowledged:
			stats["acknowledged_alerts"] = stats["acknowledged_alerts"].(int) + 1
		case AlertStatusResolved:
			stats["resolved_alerts"] = stats["resolved_alerts"].(int) + 1
		}

		severityMap := stats["by_severity"].(map[AlertSeverity]int)
		severityMap[alert.Severity]++

		typeMap := stats["by_type"].(map[RiskAlertType]int)
		typeMap[alert.Type]++
	}

	return stats
}

// EmailAlertChannel implements AlertChannel for email notifications
type EmailAlertChannel struct {
	enabled bool
	config  *EmailConfig
}

// EmailConfig holds email configuration
type EmailConfig struct {
	SMTPHost    string   `yaml:"smtp_host"`
	SMTPPort    int      `yaml:"smtp_port"`
	Username    string   `yaml:"username"`
	Password    string   `yaml:"password"`
	FromAddress string   `yaml:"from_address"`
	ToAddresses []string `yaml:"to_addresses"`
	UseTLS      bool     `yaml:"use_tls"`
}

// NewEmailAlertChannel creates a new email alert channel
func NewEmailAlertChannel(config *EmailConfig) *EmailAlertChannel {
	return &EmailAlertChannel{
		enabled: true,
		config:  config,
	}
}

// SendAlert sends an alert via email
func (eac *EmailAlertChannel) SendAlert(ctx context.Context, alert *RiskAlert) error {
	// Email sending implementation would go here
	// For now, just log that we would send an email
	return nil
}

// GetType returns the channel type
func (eac *EmailAlertChannel) GetType() string {
	return "email"
}

// IsEnabled returns whether the channel is enabled
func (eac *EmailAlertChannel) IsEnabled() bool {
	return eac.enabled
}

// WebhookAlertChannel implements AlertChannel for webhook notifications
type WebhookAlertChannel struct {
	enabled bool
	config  *WebhookConfig
}

// WebhookConfig holds webhook configuration
type WebhookConfig struct {
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers"`
	Timeout time.Duration     `yaml:"timeout"`
}

// NewWebhookAlertChannel creates a new webhook alert channel
func NewWebhookAlertChannel(config *WebhookConfig) *WebhookAlertChannel {
	return &WebhookAlertChannel{
		enabled: true,
		config:  config,
	}
}

// SendAlert sends an alert via webhook
func (wac *WebhookAlertChannel) SendAlert(ctx context.Context, alert *RiskAlert) error {
	// Webhook sending implementation would go here
	// For now, just log that we would send a webhook
	return nil
}

// GetType returns the channel type
func (wac *WebhookAlertChannel) GetType() string {
	return "webhook"
}

// IsEnabled returns whether the channel is enabled
func (wac *WebhookAlertChannel) IsEnabled() bool {
	return wac.enabled
}
