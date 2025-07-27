package web3

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
)

// RiskMonitor provides real-time risk monitoring and alerting
type RiskMonitor struct {
	clients            map[int]*ethclient.Client
	logger             *observability.Logger
	riskAssessment     *RiskAssessmentService
	vulnScanner        *VulnerabilityScanner
	alertChannels      map[string]AlertChannel
	monitoredAddresses map[string]*MonitoredAddress
	alertRules         map[string]*AlertRule
	isRunning          bool
	stopChan           chan struct{}
	mu                 sync.RWMutex
}

// MonitoredAddress represents an address being monitored for risk
type MonitoredAddress struct {
	Address     string                 `json:"address"`
	ChainID     int                    `json:"chain_id"`
	UserID      uuid.UUID              `json:"user_id"`
	AlertRules  []string               `json:"alert_rules"`
	LastChecked time.Time              `json:"last_checked"`
	RiskScore   int                    `json:"risk_score"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   time.Time              `json:"created_at"`
}

// AlertRule defines conditions for triggering alerts
type AlertRule struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Conditions    []AlertCondition       `json:"conditions"`
	Actions       []AlertAction          `json:"actions"`
	Enabled       bool                   `json:"enabled"`
	Priority      AlertPriority          `json:"priority"`
	Cooldown      time.Duration          `json:"cooldown"`
	LastTriggered time.Time              `json:"last_triggered"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// AlertCondition defines a condition that triggers an alert
type AlertCondition struct {
	Type      string      `json:"type"`
	Field     string      `json:"field"`
	Operator  string      `json:"operator"`
	Value     interface{} `json:"value"`
	Threshold float64     `json:"threshold"`
}

// AlertAction defines an action to take when an alert is triggered
type AlertAction struct {
	Type       string                 `json:"type"`
	Target     string                 `json:"target"`
	Template   string                 `json:"template"`
	Parameters map[string]interface{} `json:"parameters"`
}

// AlertPriority represents the priority level of an alert
type AlertPriority string

const (
	AlertPriorityCritical AlertPriority = "critical"
	AlertPriorityHigh     AlertPriority = "high"
	AlertPriorityMedium   AlertPriority = "medium"
	AlertPriorityLow      AlertPriority = "low"
	AlertPriorityInfo     AlertPriority = "info"
)

// Alert represents a triggered alert
type Alert struct {
	ID          uuid.UUID              `json:"id"`
	RuleID      string                 `json:"rule_id"`
	RuleName    string                 `json:"rule_name"`
	Priority    AlertPriority          `json:"priority"`
	Title       string                 `json:"title"`
	Message     string                 `json:"message"`
	Address     string                 `json:"address"`
	ChainID     int                    `json:"chain_id"`
	UserID      uuid.UUID              `json:"user_id"`
	Data        map[string]interface{} `json:"data"`
	TriggeredAt time.Time              `json:"triggered_at"`
	Status      AlertStatus            `json:"status"`
	Actions     []AlertActionResult    `json:"actions"`
}

// AlertStatus represents the status of an alert
type AlertStatus string

const (
	AlertStatusTriggered    AlertStatus = "triggered"
	AlertStatusAcknowledged AlertStatus = "acknowledged"
	AlertStatusResolved     AlertStatus = "resolved"
	AlertStatusSuppressed   AlertStatus = "suppressed"
)

// AlertActionResult represents the result of an alert action
type AlertActionResult struct {
	Type      string    `json:"type"`
	Success   bool      `json:"success"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// AlertChannel interface for different alert delivery methods
type AlertChannel interface {
	SendAlert(ctx context.Context, alert *Alert) error
	GetType() string
	IsEnabled() bool
}

// NewRiskMonitor creates a new risk monitor
func NewRiskMonitor(
	clients map[int]*ethclient.Client,
	logger *observability.Logger,
	riskAssessment *RiskAssessmentService,
	vulnScanner *VulnerabilityScanner,
) *RiskMonitor {
	monitor := &RiskMonitor{
		clients:            clients,
		logger:             logger,
		riskAssessment:     riskAssessment,
		vulnScanner:        vulnScanner,
		alertChannels:      make(map[string]AlertChannel),
		monitoredAddresses: make(map[string]*MonitoredAddress),
		alertRules:         make(map[string]*AlertRule),
		stopChan:           make(chan struct{}),
	}

	// Initialize default alert rules
	monitor.initializeDefaultRules()

	// Initialize alert channels
	monitor.initializeAlertChannels()

	return monitor
}

// Start starts the risk monitoring service
func (r *RiskMonitor) Start(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.isRunning {
		return fmt.Errorf("risk monitor is already running")
	}

	r.isRunning = true

	// Start monitoring goroutine
	go r.monitoringLoop(ctx)

	r.logger.Info(ctx, "Risk monitor started", map[string]interface{}{
		"monitored_addresses": len(r.monitoredAddresses),
		"alert_rules":         len(r.alertRules),
		"alert_channels":      len(r.alertChannels),
	})

	return nil
}

// Stop stops the risk monitoring service
func (r *RiskMonitor) Stop(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.isRunning {
		return fmt.Errorf("risk monitor is not running")
	}

	close(r.stopChan)
	r.isRunning = false

	r.logger.Info(ctx, "Risk monitor stopped", nil)

	return nil
}

// AddMonitoredAddress adds an address to be monitored
func (r *RiskMonitor) AddMonitoredAddress(ctx context.Context, address string, chainID int, userID uuid.UUID, alertRules []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := fmt.Sprintf("%s:%d", address, chainID)

	monitoredAddr := &MonitoredAddress{
		Address:     address,
		ChainID:     chainID,
		UserID:      userID,
		AlertRules:  alertRules,
		LastChecked: time.Time{},
		RiskScore:   0,
		Metadata:    make(map[string]interface{}),
		CreatedAt:   time.Now(),
	}

	r.monitoredAddresses[key] = monitoredAddr

	r.logger.Info(ctx, "Address added to monitoring", map[string]interface{}{
		"address":     address,
		"chain_id":    chainID,
		"user_id":     userID.String(),
		"alert_rules": alertRules,
	})

	return nil
}

// RemoveMonitoredAddress removes an address from monitoring
func (r *RiskMonitor) RemoveMonitoredAddress(ctx context.Context, address string, chainID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := fmt.Sprintf("%s:%d", address, chainID)
	delete(r.monitoredAddresses, key)

	r.logger.Info(ctx, "Address removed from monitoring", map[string]interface{}{
		"address":  address,
		"chain_id": chainID,
	})

	return nil
}

// AddAlertRule adds a new alert rule
func (r *RiskMonitor) AddAlertRule(ctx context.Context, rule *AlertRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.alertRules[rule.ID] = rule

	r.logger.Info(ctx, "Alert rule added", map[string]interface{}{
		"rule_id":   rule.ID,
		"rule_name": rule.Name,
		"enabled":   rule.Enabled,
	})

	return nil
}

// monitoringLoop is the main monitoring loop
func (r *RiskMonitor) monitoringLoop(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second) // Check every 30 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-r.stopChan:
			return
		case <-ticker.C:
			r.performMonitoringCheck(ctx)
		}
	}
}

// performMonitoringCheck performs a monitoring check on all addresses
func (r *RiskMonitor) performMonitoringCheck(ctx context.Context) {
	r.mu.RLock()
	addresses := make([]*MonitoredAddress, 0, len(r.monitoredAddresses))
	for _, addr := range r.monitoredAddresses {
		addresses = append(addresses, addr)
	}
	r.mu.RUnlock()

	for _, addr := range addresses {
		r.checkAddress(ctx, addr)
	}
}

// checkAddress performs risk checks on a specific address
func (r *RiskMonitor) checkAddress(ctx context.Context, addr *MonitoredAddress) {
	// Skip if checked recently (within last 5 minutes)
	if time.Since(addr.LastChecked) < 5*time.Minute {
		return
	}

	// Perform contract risk assessment if it's a contract
	if r.isContract(ctx, addr.Address, addr.ChainID) {
		r.checkContractRisk(ctx, addr)
	}

	// Check for suspicious activity
	r.checkSuspiciousActivity(ctx, addr)

	// Update last checked time
	addr.LastChecked = time.Now()
}

// isContract checks if an address is a contract
func (r *RiskMonitor) isContract(ctx context.Context, address string, chainID int) bool {
	client, exists := r.clients[chainID]
	if !exists {
		return false
	}

	code, err := client.CodeAt(ctx, common.HexToAddress(address), nil)
	if err != nil {
		return false
	}

	return len(code) > 0
}

// checkContractRisk performs contract risk assessment
func (r *RiskMonitor) checkContractRisk(ctx context.Context, addr *MonitoredAddress) {
	req := ContractRiskRequest{
		ContractAddress: addr.Address,
		ChainID:         addr.ChainID,
		AnalyzeCode:     true,
		CheckRugPull:    true,
		IncludeMLModels: true,
	}

	assessment, err := r.riskAssessment.AssessContractRisk(ctx, req)
	if err != nil {
		r.logger.Warn(ctx, "Contract risk assessment failed", map[string]interface{}{
			"address": addr.Address,
			"error":   err.Error(),
		})
		return
	}

	// Update risk score
	addr.RiskScore = assessment.RiskScore

	// Check alert rules
	r.evaluateAlertRules(ctx, addr, assessment)
}

// checkSuspiciousActivity checks for suspicious activity patterns
func (r *RiskMonitor) checkSuspiciousActivity(ctx context.Context, addr *MonitoredAddress) {
	// This would implement real-time activity monitoring
	// For now, add placeholder monitoring
	r.logger.Debug(ctx, "Checking suspicious activity", map[string]interface{}{
		"address": addr.Address,
	})
}

// evaluateAlertRules evaluates alert rules against assessment results
func (r *RiskMonitor) evaluateAlertRules(ctx context.Context, addr *MonitoredAddress, assessment *RiskAssessment) {
	for _, ruleID := range addr.AlertRules {
		rule, exists := r.alertRules[ruleID]
		if !exists || !rule.Enabled {
			continue
		}

		// Check cooldown
		if time.Since(rule.LastTriggered) < rule.Cooldown {
			continue
		}

		// Evaluate conditions
		if r.evaluateConditions(rule.Conditions, assessment) {
			r.triggerAlert(ctx, rule, addr, assessment)
		}
	}
}

// evaluateConditions evaluates alert conditions
func (r *RiskMonitor) evaluateConditions(conditions []AlertCondition, assessment *RiskAssessment) bool {
	for _, condition := range conditions {
		if !r.evaluateCondition(condition, assessment) {
			return false // All conditions must be true
		}
	}
	return true
}

// evaluateCondition evaluates a single alert condition
func (r *RiskMonitor) evaluateCondition(condition AlertCondition, assessment *RiskAssessment) bool {
	switch condition.Type {
	case "risk_score":
		return r.compareValues(float64(assessment.RiskScore), condition.Operator, condition.Threshold)
	case "safety_grade":
		gradeValue := r.safetyGradeToValue(assessment.SafetyGrade)
		return r.compareValues(gradeValue, condition.Operator, condition.Threshold)
	case "vulnerability_count":
		// This would check vulnerability count from scanner results
		return false
	default:
		return false
	}
}

// compareValues compares two values using an operator
func (r *RiskMonitor) compareValues(actual float64, operator string, threshold float64) bool {
	switch operator {
	case "gt", ">":
		return actual > threshold
	case "gte", ">=":
		return actual >= threshold
	case "lt", "<":
		return actual < threshold
	case "lte", "<=":
		return actual <= threshold
	case "eq", "==":
		return actual == threshold
	case "ne", "!=":
		return actual != threshold
	default:
		return false
	}
}

// safetyGradeToValue converts safety grade to numeric value
func (r *RiskMonitor) safetyGradeToValue(grade SafetyGrade) float64 {
	switch grade {
	case SafetyGradeA:
		return 5.0
	case SafetyGradeB:
		return 4.0
	case SafetyGradeC:
		return 3.0
	case SafetyGradeD:
		return 2.0
	case SafetyGradeF:
		return 1.0
	default:
		return 0.0
	}
}

// triggerAlert triggers an alert
func (r *RiskMonitor) triggerAlert(ctx context.Context, rule *AlertRule, addr *MonitoredAddress, assessment *RiskAssessment) {
	alert := &Alert{
		ID:       uuid.New(),
		RuleID:   rule.ID,
		RuleName: rule.Name,
		Priority: rule.Priority,
		Title:    fmt.Sprintf("Risk Alert: %s", rule.Name),
		Message:  r.formatAlertMessage(rule, addr, assessment),
		Address:  addr.Address,
		ChainID:  addr.ChainID,
		UserID:   addr.UserID,
		Data: map[string]interface{}{
			"risk_score":   assessment.RiskScore,
			"safety_grade": assessment.SafetyGrade,
			"risk_level":   assessment.RiskLevel,
		},
		TriggeredAt: time.Now(),
		Status:      AlertStatusTriggered,
		Actions:     []AlertActionResult{},
	}

	// Execute alert actions
	r.executeAlertActions(ctx, alert, rule.Actions)

	// Update rule last triggered time
	rule.LastTriggered = time.Now()

	r.logger.Info(ctx, "Alert triggered", map[string]interface{}{
		"alert_id":   alert.ID.String(),
		"rule_id":    rule.ID,
		"address":    addr.Address,
		"risk_score": assessment.RiskScore,
		"priority":   string(rule.Priority),
	})
}

// formatAlertMessage formats the alert message
func (r *RiskMonitor) formatAlertMessage(rule *AlertRule, addr *MonitoredAddress, assessment *RiskAssessment) string {
	return fmt.Sprintf("Risk alert for address %s: Risk score %d, Safety grade %s",
		addr.Address, assessment.RiskScore, assessment.SafetyGrade)
}

// executeAlertActions executes alert actions
func (r *RiskMonitor) executeAlertActions(ctx context.Context, alert *Alert, actions []AlertAction) {
	for _, action := range actions {
		result := r.executeAlertAction(ctx, alert, action)
		alert.Actions = append(alert.Actions, result)
	}
}

// executeAlertAction executes a single alert action
func (r *RiskMonitor) executeAlertAction(ctx context.Context, alert *Alert, action AlertAction) AlertActionResult {
	result := AlertActionResult{
		Type:      action.Type,
		Timestamp: time.Now(),
	}

	switch action.Type {
	case "email":
		// Send email notification
		result.Success = true
		result.Message = "Email notification sent"
	case "webhook":
		// Send webhook notification
		result.Success = true
		result.Message = "Webhook notification sent"
	case "log":
		// Log the alert
		r.logger.Warn(ctx, "Risk alert", map[string]interface{}{
			"alert_id": alert.ID.String(),
			"message":  alert.Message,
		})
		result.Success = true
		result.Message = "Alert logged"
	default:
		result.Success = false
		result.Message = fmt.Sprintf("Unknown action type: %s", action.Type)
	}

	return result
}

// initializeDefaultRules initializes default alert rules
func (r *RiskMonitor) initializeDefaultRules() {
	// High risk score alert
	r.alertRules["high_risk_score"] = &AlertRule{
		ID:          "high_risk_score",
		Name:        "High Risk Score Alert",
		Description: "Triggers when risk score exceeds 70",
		Conditions: []AlertCondition{
			{
				Type:      "risk_score",
				Operator:  "gt",
				Threshold: 70.0,
			},
		},
		Actions: []AlertAction{
			{
				Type:   "log",
				Target: "system",
			},
		},
		Enabled:  true,
		Priority: AlertPriorityHigh,
		Cooldown: 15 * time.Minute,
	}

	// Critical risk score alert
	r.alertRules["critical_risk_score"] = &AlertRule{
		ID:          "critical_risk_score",
		Name:        "Critical Risk Score Alert",
		Description: "Triggers when risk score exceeds 90",
		Conditions: []AlertCondition{
			{
				Type:      "risk_score",
				Operator:  "gt",
				Threshold: 90.0,
			},
		},
		Actions: []AlertAction{
			{
				Type:   "log",
				Target: "system",
			},
		},
		Enabled:  true,
		Priority: AlertPriorityCritical,
		Cooldown: 5 * time.Minute,
	}

	// Poor safety grade alert
	r.alertRules["poor_safety_grade"] = &AlertRule{
		ID:          "poor_safety_grade",
		Name:        "Poor Safety Grade Alert",
		Description: "Triggers when safety grade is D or F",
		Conditions: []AlertCondition{
			{
				Type:      "safety_grade",
				Operator:  "lte",
				Threshold: 2.0, // D or F grade
			},
		},
		Actions: []AlertAction{
			{
				Type:   "log",
				Target: "system",
			},
		},
		Enabled:  true,
		Priority: AlertPriorityMedium,
		Cooldown: 30 * time.Minute,
	}
}

// initializeAlertChannels initializes alert delivery channels
func (r *RiskMonitor) initializeAlertChannels() {
	// Initialize log channel
	r.alertChannels["log"] = &LogAlertChannel{
		logger:  r.logger,
		enabled: true,
	}

	// Initialize webhook channel (placeholder)
	r.alertChannels["webhook"] = &WebhookAlertChannel{
		enabled: false, // Disabled by default
	}

	// Initialize email channel (placeholder)
	r.alertChannels["email"] = &EmailAlertChannel{
		enabled: false, // Disabled by default
	}
}

// LogAlertChannel implements AlertChannel for logging
type LogAlertChannel struct {
	logger  *observability.Logger
	enabled bool
}

func (l *LogAlertChannel) SendAlert(ctx context.Context, alert *Alert) error {
	if !l.enabled {
		return fmt.Errorf("log alert channel is disabled")
	}

	l.logger.Warn(ctx, "Risk Alert", map[string]interface{}{
		"alert_id":  alert.ID.String(),
		"rule_name": alert.RuleName,
		"priority":  string(alert.Priority),
		"address":   alert.Address,
		"chain_id":  alert.ChainID,
		"message":   alert.Message,
		"data":      alert.Data,
	})

	return nil
}

func (l *LogAlertChannel) GetType() string {
	return "log"
}

func (l *LogAlertChannel) IsEnabled() bool {
	return l.enabled
}

// WebhookAlertChannel implements AlertChannel for webhooks
type WebhookAlertChannel struct {
	webhookURL string
	enabled    bool
}

func (w *WebhookAlertChannel) SendAlert(ctx context.Context, alert *Alert) error {
	if !w.enabled {
		return fmt.Errorf("webhook alert channel is disabled")
	}

	// This would send HTTP POST to webhook URL
	// For now, just return success
	return nil
}

func (w *WebhookAlertChannel) GetType() string {
	return "webhook"
}

func (w *WebhookAlertChannel) IsEnabled() bool {
	return w.enabled
}

// EmailAlertChannel implements AlertChannel for email
type EmailAlertChannel struct {
	smtpConfig map[string]string
	enabled    bool
}

func (e *EmailAlertChannel) SendAlert(ctx context.Context, alert *Alert) error {
	if !e.enabled {
		return fmt.Errorf("email alert channel is disabled")
	}

	// This would send email notification
	// For now, just return success
	return nil
}

func (e *EmailAlertChannel) GetType() string {
	return "email"
}

func (e *EmailAlertChannel) IsEnabled() bool {
	return e.enabled
}
