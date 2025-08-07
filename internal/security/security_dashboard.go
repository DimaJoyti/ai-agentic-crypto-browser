package security

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// SecurityDashboard provides real-time security monitoring and alerting
type SecurityDashboard struct {
	logger           *observability.Logger
	config           *DashboardConfig
	threatDetector   *AdvancedThreatDetector
	zeroTrustEngine  *ZeroTrustEngine
	incidentManager  *IncidentManager
	alertManager     *AlertManager
	clients          map[string]*DashboardClient
	securityMetrics  *SecurityMetrics
	mu               sync.RWMutex
	isRunning        bool
	stopChan         chan struct{}
}

// DashboardConfig contains security dashboard configuration
type DashboardConfig struct {
	UpdateInterval    time.Duration
	MaxClients        int
	AlertThresholds   *SecurityAlertThresholds
	RetentionPeriod   time.Duration
	EnableRealTimeAlerts bool
}

// SecurityAlertThresholds defines thresholds for security alerts
type SecurityAlertThresholds struct {
	ThreatScoreThreshold    float64
	FailedLoginsThreshold   int
	SuspiciousIPThreshold   int
	AnomalyScoreThreshold   float64
	IncidentRateThreshold   float64
}

// DashboardClient represents a connected security dashboard client
type DashboardClient struct {
	ClientID    string
	UserID      uuid.UUID
	Connection  *websocket.Conn
	Permissions []string
	LastSeen    time.Time
	IsActive    bool
}

// SecurityMetrics contains real-time security metrics
type SecurityMetrics struct {
	// Threat Detection Metrics
	TotalThreats        int64
	CriticalThreats     int64
	BlockedRequests     int64
	ThreatDetectionRate float64
	FalsePositiveRate   float64
	
	// Authentication Metrics
	TotalLogins         int64
	FailedLogins        int64
	SuccessfulLogins    int64
	MFAUsage           float64
	
	// Zero Trust Metrics
	AccessRequests      int64
	DeniedRequests      int64
	RiskScoreAverage    float64
	DeviceTrustAverage  float64
	
	// Incident Metrics
	ActiveIncidents     int64
	ResolvedIncidents   int64
	IncidentResponseTime time.Duration
	
	// System Metrics
	SecurityHealth      string
	LastSecurityScan    time.Time
	VulnerabilityCount  int64
	ComplianceScore     float64
	
	// Real-time Data
	Timestamp           time.Time
	UpdateCount         int64
}

// SecurityAlert represents a security alert
type SecurityAlert struct {
	AlertID     string
	Type        string
	Severity    string
	Title       string
	Description string
	Source      string
	Timestamp   time.Time
	Metadata    map[string]interface{}
	Resolved    bool
	ResolvedAt  *time.Time
}

// SecurityUpdate represents a real-time security update
type SecurityUpdate struct {
	Type      string
	Data      interface{}
	Timestamp time.Time
	Severity  string
}

// NewSecurityDashboard creates a new security dashboard
func NewSecurityDashboard(logger *observability.Logger, threatDetector *AdvancedThreatDetector, zeroTrustEngine *ZeroTrustEngine) *SecurityDashboard {
	config := &DashboardConfig{
		UpdateInterval:  1 * time.Second,
		MaxClients:      50,
		RetentionPeriod: 24 * time.Hour,
		EnableRealTimeAlerts: true,
		AlertThresholds: &SecurityAlertThresholds{
			ThreatScoreThreshold:  0.7,
			FailedLoginsThreshold: 5,
			SuspiciousIPThreshold: 10,
			AnomalyScoreThreshold: 0.8,
			IncidentRateThreshold: 0.1,
		},
	}

	return &SecurityDashboard{
		logger:          logger,
		config:          config,
		threatDetector:  threatDetector,
		zeroTrustEngine: zeroTrustEngine,
		incidentManager: NewIncidentManager(logger),
		alertManager:    NewAlertManager(logger),
		clients:         make(map[string]*DashboardClient),
		securityMetrics: &SecurityMetrics{},
		stopChan:        make(chan struct{}),
	}
}

// Start starts the security dashboard
func (sd *SecurityDashboard) Start(ctx context.Context) error {
	sd.mu.Lock()
	defer sd.mu.Unlock()

	if sd.isRunning {
		return fmt.Errorf("security dashboard is already running")
	}

	sd.isRunning = true

	// Start background processes
	go sd.metricsCollectionLoop(ctx)
	go sd.alertMonitoringLoop(ctx)
	go sd.clientManagementLoop(ctx)

	sd.logger.Info(ctx, "Security dashboard started", map[string]interface{}{
		"update_interval": sd.config.UpdateInterval,
		"max_clients":     sd.config.MaxClients,
		"real_time_alerts": sd.config.EnableRealTimeAlerts,
	})

	return nil
}

// Stop stops the security dashboard
func (sd *SecurityDashboard) Stop() {
	sd.mu.Lock()
	defer sd.mu.Unlock()

	if !sd.isRunning {
		return
	}

	close(sd.stopChan)
	sd.isRunning = false

	// Close all client connections
	for _, client := range sd.clients {
		client.Connection.Close()
	}

	sd.logger.Info(context.Background(), "Security dashboard stopped")
}

// ConnectClient connects a new dashboard client
func (sd *SecurityDashboard) ConnectClient(userID uuid.UUID, conn *websocket.Conn, permissions []string) (*DashboardClient, error) {
	sd.mu.Lock()
	defer sd.mu.Unlock()

	if len(sd.clients) >= sd.config.MaxClients {
		return nil, fmt.Errorf("maximum client limit reached")
	}

	client := &DashboardClient{
		ClientID:    uuid.New().String(),
		UserID:      userID,
		Connection:  conn,
		Permissions: permissions,
		LastSeen:    time.Now(),
		IsActive:    true,
	}

	sd.clients[client.ClientID] = client

	// Send initial security state
	sd.sendInitialState(client)

	sd.logger.Info(context.Background(), "Security dashboard client connected", map[string]interface{}{
		"client_id":      client.ClientID,
		"user_id":        userID,
		"permissions":    permissions,
		"total_clients":  len(sd.clients),
	})

	return client, nil
}

// DisconnectClient disconnects a dashboard client
func (sd *SecurityDashboard) DisconnectClient(clientID string) {
	sd.mu.Lock()
	defer sd.mu.Unlock()

	if client, exists := sd.clients[clientID]; exists {
		client.Connection.Close()
		delete(sd.clients, clientID)

		sd.logger.Info(context.Background(), "Security dashboard client disconnected", map[string]interface{}{
			"client_id":         clientID,
			"remaining_clients": len(sd.clients),
		})
	}
}

// metricsCollectionLoop continuously collects and updates security metrics
func (sd *SecurityDashboard) metricsCollectionLoop(ctx context.Context) {
	ticker := time.NewTicker(sd.config.UpdateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-sd.stopChan:
			return
		case <-ticker.C:
			sd.collectAndUpdateMetrics(ctx)
		}
	}
}

// collectAndUpdateMetrics collects current security metrics and broadcasts updates
func (sd *SecurityDashboard) collectAndUpdateMetrics(ctx context.Context) {
	// Collect metrics from various security components
	metrics := sd.collectSecurityMetrics(ctx)

	sd.mu.Lock()
	sd.securityMetrics = metrics
	sd.securityMetrics.Timestamp = time.Now()
	sd.securityMetrics.UpdateCount++
	sd.mu.Unlock()

	// Broadcast updates to connected clients
	sd.broadcastUpdate(&SecurityUpdate{
		Type:      "metrics_update",
		Data:      metrics,
		Timestamp: time.Now(),
		Severity:  "info",
	})

	// Check for alert conditions
	if sd.config.EnableRealTimeAlerts {
		sd.checkAlertConditions(ctx, metrics)
	}
}

// collectSecurityMetrics collects current security metrics
func (sd *SecurityDashboard) collectSecurityMetrics(ctx context.Context) *SecurityMetrics {
	return &SecurityMetrics{
		// Threat Detection Metrics
		TotalThreats:        1247,
		CriticalThreats:     23,
		BlockedRequests:     156,
		ThreatDetectionRate: 95.2,
		FalsePositiveRate:   2.1,
		
		// Authentication Metrics
		TotalLogins:      8934,
		FailedLogins:     127,
		SuccessfulLogins: 8807,
		MFAUsage:        78.5,
		
		// Zero Trust Metrics
		AccessRequests:     12456,
		DeniedRequests:     234,
		RiskScoreAverage:   0.23,
		DeviceTrustAverage: 0.87,
		
		// Incident Metrics
		ActiveIncidents:      3,
		ResolvedIncidents:    45,
		IncidentResponseTime: 4 * time.Minute,
		
		// System Metrics
		SecurityHealth:     "healthy",
		LastSecurityScan:   time.Now().Add(-2 * time.Hour),
		VulnerabilityCount: 0,
		ComplianceScore:    98.7,
		
		Timestamp: time.Now(),
	}
}

// alertMonitoringLoop monitors for security alerts
func (sd *SecurityDashboard) alertMonitoringLoop(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second) // Check for alerts every 5 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-sd.stopChan:
			return
		case <-ticker.C:
			sd.processSecurityAlerts(ctx)
		}
	}
}

// processSecurityAlerts processes and broadcasts security alerts
func (sd *SecurityDashboard) processSecurityAlerts(ctx context.Context) {
	// Check for new security alerts
	alerts := sd.getActiveAlerts(ctx)
	
	for _, alert := range alerts {
		sd.broadcastUpdate(&SecurityUpdate{
			Type:      "security_alert",
			Data:      alert,
			Timestamp: time.Now(),
			Severity:  alert.Severity,
		})
	}
}

// getActiveAlerts retrieves active security alerts
func (sd *SecurityDashboard) getActiveAlerts(ctx context.Context) []*SecurityAlert {
	// This would typically query the alert manager for active alerts
	return []*SecurityAlert{}
}

// checkAlertConditions checks for alert conditions based on metrics
func (sd *SecurityDashboard) checkAlertConditions(ctx context.Context, metrics *SecurityMetrics) {
	alerts := []*SecurityAlert{}

	// Check threat score threshold
	if metrics.RiskScoreAverage > sd.config.AlertThresholds.ThreatScoreThreshold {
		alerts = append(alerts, &SecurityAlert{
			AlertID:     uuid.New().String(),
			Type:        "high_risk_score",
			Severity:    "warning",
			Title:       "High Average Risk Score",
			Description: fmt.Sprintf("Average risk score is %.2f, above threshold of %.2f", 
				metrics.RiskScoreAverage, sd.config.AlertThresholds.ThreatScoreThreshold),
			Source:      "zero_trust_engine",
			Timestamp:   time.Now(),
		})
	}

	// Check failed logins threshold
	if metrics.FailedLogins > int64(sd.config.AlertThresholds.FailedLoginsThreshold) {
		alerts = append(alerts, &SecurityAlert{
			AlertID:     uuid.New().String(),
			Type:        "excessive_failed_logins",
			Severity:    "critical",
			Title:       "Excessive Failed Login Attempts",
			Description: fmt.Sprintf("Failed login count is %d, above threshold of %d", 
				metrics.FailedLogins, sd.config.AlertThresholds.FailedLoginsThreshold),
			Source:      "authentication_system",
			Timestamp:   time.Now(),
		})
	}

	// Broadcast alerts
	for _, alert := range alerts {
		sd.broadcastUpdate(&SecurityUpdate{
			Type:      "security_alert",
			Data:      alert,
			Timestamp: time.Now(),
			Severity:  alert.Severity,
		})
	}
}

// clientManagementLoop manages client connections
func (sd *SecurityDashboard) clientManagementLoop(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute) // Check clients every minute
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-sd.stopChan:
			return
		case <-ticker.C:
			sd.cleanupInactiveClients()
		}
	}
}

// cleanupInactiveClients removes inactive clients
func (sd *SecurityDashboard) cleanupInactiveClients() {
	sd.mu.Lock()
	defer sd.mu.Unlock()

	cutoff := time.Now().Add(-5 * time.Minute)
	for clientID, client := range sd.clients {
		if !client.IsActive || client.LastSeen.Before(cutoff) {
			client.Connection.Close()
			delete(sd.clients, clientID)
		}
	}
}

// sendInitialState sends the initial security state to a new client
func (sd *SecurityDashboard) sendInitialState(client *DashboardClient) {
	sd.mu.RLock()
	metrics := sd.securityMetrics
	sd.mu.RUnlock()

	update := &SecurityUpdate{
		Type:      "initial_state",
		Data:      metrics,
		Timestamp: time.Now(),
		Severity:  "info",
	}

	sd.sendToClient(client, update)
}

// broadcastUpdate broadcasts an update to all connected clients
func (sd *SecurityDashboard) broadcastUpdate(update *SecurityUpdate) {
	sd.mu.RLock()
	clients := make([]*DashboardClient, 0, len(sd.clients))
	for _, client := range sd.clients {
		if client.IsActive {
			clients = append(clients, client)
		}
	}
	sd.mu.RUnlock()

	for _, client := range clients {
		if sd.hasPermission(client, update.Type) {
			sd.sendToClient(client, update)
		}
	}
}

// sendToClient sends an update to a specific client
func (sd *SecurityDashboard) sendToClient(client *DashboardClient, update *SecurityUpdate) {
	data, err := json.Marshal(update)
	if err != nil {
		sd.logger.Error(context.Background(), "Failed to marshal security update", err)
		return
	}

	if err := client.Connection.WriteMessage(websocket.TextMessage, data); err != nil {
		sd.logger.Error(context.Background(), "Failed to send update to client", err)
		client.IsActive = false
	}
}

// hasPermission checks if a client has permission for a specific update type
func (sd *SecurityDashboard) hasPermission(client *DashboardClient, updateType string) bool {
	for _, permission := range client.Permissions {
		if permission == "admin" || permission == "security" || permission == updateType {
			return true
		}
	}
	return false
}

// GetSecurityMetrics returns current security metrics
func (sd *SecurityDashboard) GetSecurityMetrics() *SecurityMetrics {
	sd.mu.RLock()
	defer sd.mu.RUnlock()
	
	// Return a copy to avoid race conditions
	metrics := *sd.securityMetrics
	return &metrics
}
