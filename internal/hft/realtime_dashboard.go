package hft

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// RealtimeDashboard provides comprehensive real-time monitoring with live metrics,
// alerts, performance tracking, and interactive visualizations for HFT operations
type RealtimeDashboard struct {
	logger *observability.Logger
	config DashboardConfig

	// Data aggregators
	metricsAggregator   *MetricsAggregator
	alertManager        *DashboardAlertManager
	performanceTracker  *PerformanceTracker
	visualizationEngine *VisualizationEngine

	// Real-time data streams
	liveMetrics   *LiveMetrics
	alertStream   chan *DashboardAlert
	metricsStream chan *MetricUpdate
	eventStream   chan *DashboardEvent

	// Dashboard state
	dashboardSessions map[string]*DashboardSession
	widgets           map[string]*Widget
	layouts           map[string]*Layout

	// Performance tracking
	metricsProcessed int64
	alertsGenerated  int64
	sessionsActive   int64
	avgUpdateTime    int64

	// State management
	isRunning int32
	stopChan  chan struct{}
	wg        sync.WaitGroup
	mu        sync.RWMutex

	// WebSocket connections for real-time updates
	wsConnections map[string]*WSConnection
	wsMu          sync.RWMutex
}

// DashboardConfig contains configuration for the real-time dashboard
type DashboardConfig struct {
	// Update intervals
	MetricsUpdateInterval time.Duration `json:"metrics_update_interval"`
	AlertCheckInterval    time.Duration `json:"alert_check_interval"`
	PerformanceInterval   time.Duration `json:"performance_interval"`

	// Data retention
	MetricsRetention time.Duration `json:"metrics_retention"`
	AlertRetention   time.Duration `json:"alert_retention"`
	EventRetention   time.Duration `json:"event_retention"`

	// Display settings
	MaxDataPoints int           `json:"max_data_points"`
	RefreshRate   time.Duration `json:"refresh_rate"`
	AutoRefresh   bool          `json:"auto_refresh"`

	// Alert thresholds
	LatencyThreshold    time.Duration `json:"latency_threshold"`
	ThroughputThreshold int64         `json:"throughput_threshold"`
	ErrorRateThreshold  float64       `json:"error_rate_threshold"`

	// Visualization settings
	ChartTypes       []string `json:"chart_types"`
	ColorScheme      string   `json:"color_scheme"`
	AnimationEnabled bool     `json:"animation_enabled"`

	// WebSocket settings
	WSPort         int           `json:"ws_port"`
	WSPath         string        `json:"ws_path"`
	MaxConnections int           `json:"max_connections"`
	PingInterval   time.Duration `json:"ping_interval"`
}

// LiveMetrics contains current real-time metrics
type LiveMetrics struct {
	// System metrics
	SystemHealth       SystemHealth       `json:"system_health"`
	PerformanceMetrics PerformanceMetrics `json:"performance_metrics"`
	TradingMetrics     TradingMetrics     `json:"trading_metrics"`
	NetworkMetrics     NetworkMetrics     `json:"network_metrics"`
	RiskMetrics        RiskMetrics        `json:"risk_metrics"`

	// Timestamps
	LastUpdate      time.Time     `json:"last_update"`
	UpdateFrequency time.Duration `json:"update_frequency"`
	DataPoints      int           `json:"data_points"`
}

// SystemHealth contains system health indicators
type SystemHealth struct {
	Status            string        `json:"status"`
	Uptime            time.Duration `json:"uptime"`
	CPUUsage          float64       `json:"cpu_usage"`
	MemoryUsage       float64       `json:"memory_usage"`
	DiskUsage         float64       `json:"disk_usage"`
	ActiveConnections int           `json:"active_connections"`
	ErrorRate         float64       `json:"error_rate"`
	LastHealthCheck   time.Time     `json:"last_health_check"`
}

// TradingMetrics contains trading-specific metrics
type TradingMetrics struct {
	// Order metrics
	OrdersPerSecond float64         `json:"orders_per_second"`
	FillRate        float64         `json:"fill_rate"`
	AvgOrderSize    decimal.Decimal `json:"avg_order_size"`

	// Execution metrics
	AvgExecutionTime time.Duration `json:"avg_execution_time"`
	AvgSlippage      float64       `json:"avg_slippage"`
	SuccessRate      float64       `json:"success_rate"`

	// P&L metrics
	DailyPnL      decimal.Decimal `json:"daily_pnl"`
	UnrealizedPnL decimal.Decimal `json:"unrealized_pnl"`
	TotalVolume   decimal.Decimal `json:"total_volume"`

	// Strategy metrics
	ActiveStrategies int    `json:"active_strategies"`
	TopPerformer     string `json:"top_performer"`
	WorstPerformer   string `json:"worst_performer"`
}

// DashboardAlert represents a dashboard alert
type DashboardAlert struct {
	ID             uuid.UUID              `json:"id"`
	Type           AlertType              `json:"type"`
	Severity       AlertSeverity          `json:"severity"`
	Title          string                 `json:"title"`
	Message        string                 `json:"message"`
	Source         string                 `json:"source"`
	Timestamp      time.Time              `json:"timestamp"`
	Data           map[string]interface{} `json:"data"`
	Acknowledged   bool                   `json:"acknowledged"`
	AcknowledgedBy string                 `json:"acknowledged_by,omitempty"`
	AcknowledgedAt time.Time              `json:"acknowledged_at,omitempty"`
	Resolved       bool                   `json:"resolved"`
	ResolvedAt     time.Time              `json:"resolved_at,omitempty"`
}

// AlertType represents different types of alerts
type AlertType string

const (
	AlertTypePerformance AlertType = "PERFORMANCE"
	AlertTypeLatency     AlertType = "LATENCY"
	AlertTypeThroughput  AlertType = "THROUGHPUT"
	AlertTypeError       AlertType = "ERROR"
	AlertTypeRisk        AlertType = "RISK"
	AlertTypeSystem      AlertType = "SYSTEM"
	AlertTypeTrading     AlertType = "TRADING"
)

// AlertSeverity represents alert severity levels
type AlertSeverity string

const (
	AlertSeverityInfo     AlertSeverity = "INFO"
	AlertSeverityWarning  AlertSeverity = "WARNING"
	AlertSeverityError    AlertSeverity = "ERROR"
	AlertSeverityCritical AlertSeverity = "CRITICAL"
)

// MetricUpdate represents a metric update
type MetricUpdate struct {
	ID         uuid.UUID         `json:"id"`
	MetricName string            `json:"metric_name"`
	Value      interface{}       `json:"value"`
	Timestamp  time.Time         `json:"timestamp"`
	Source     string            `json:"source"`
	Tags       map[string]string `json:"tags"`
}

// DashboardEvent represents a dashboard event
type DashboardEvent struct {
	ID          uuid.UUID              `json:"id"`
	Type        EventType              `json:"type"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Timestamp   time.Time              `json:"timestamp"`
	Source      string                 `json:"source"`
	Data        map[string]interface{} `json:"data"`
	Impact      EventImpact            `json:"impact"`
}

// EventType represents different types of events
type EventType string

const (
	EventTypeOrderPlaced    EventType = "ORDER_PLACED"
	EventTypeOrderFilled    EventType = "ORDER_FILLED"
	EventTypeOrderCancelled EventType = "ORDER_CANCELLED"
	EventTypeTradeExecuted  EventType = "TRADE_EXECUTED"
	EventTypeRiskViolation  EventType = "RISK_VIOLATION"
	EventTypeSystemEvent    EventType = "SYSTEM_EVENT"
	EventTypeStrategyEvent  EventType = "STRATEGY_EVENT"
)

// EventImpact represents the impact level of an event
type EventImpact string

const (
	EventImpactLow      EventImpact = "LOW"
	EventImpactMedium   EventImpact = "MEDIUM"
	EventImpactHigh     EventImpact = "HIGH"
	EventImpactCritical EventImpact = "CRITICAL"
)

// DashboardSession represents a user dashboard session
type DashboardSession struct {
	ID           uuid.UUID              `json:"id"`
	UserID       string                 `json:"user_id"`
	StartTime    time.Time              `json:"start_time"`
	LastActivity time.Time              `json:"last_activity"`
	Layout       *Layout                `json:"layout"`
	Widgets      []string               `json:"widgets"`
	Preferences  map[string]interface{} `json:"preferences"`
	IsActive     bool                   `json:"is_active"`
}

// Widget represents a dashboard widget
type Widget struct {
	ID              string                 `json:"id"`
	Type            WidgetType             `json:"type"`
	Title           string                 `json:"title"`
	Position        WidgetPosition         `json:"position"`
	Size            WidgetSize             `json:"size"`
	Config          map[string]interface{} `json:"config"`
	DataSource      string                 `json:"data_source"`
	RefreshInterval time.Duration          `json:"refresh_interval"`
	LastUpdate      time.Time              `json:"last_update"`
}

// WidgetType represents different types of widgets
type WidgetType string

const (
	WidgetTypeChart       WidgetType = "CHART"
	WidgetTypeTable       WidgetType = "TABLE"
	WidgetTypeMetric      WidgetType = "METRIC"
	WidgetTypeAlert       WidgetType = "ALERT"
	WidgetTypeLog         WidgetType = "LOG"
	WidgetTypeOrderBook   WidgetType = "ORDER_BOOK"
	WidgetTypePositions   WidgetType = "POSITIONS"
	WidgetTypePerformance WidgetType = "PERFORMANCE"
)

// Layout represents a dashboard layout
type Layout struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Widgets     []string   `json:"widgets"`
	Grid        GridConfig `json:"grid"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// GridConfig represents grid configuration
type GridConfig struct {
	Columns int `json:"columns"`
	Rows    int `json:"rows"`
	Margin  int `json:"margin"`
	Padding int `json:"padding"`
}

// WSConnection represents a WebSocket connection
type WSConnection struct {
	ID               uuid.UUID `json:"id"`
	UserID           string    `json:"user_id"`
	ConnectedAt      time.Time `json:"connected_at"`
	LastPing         time.Time `json:"last_ping"`
	Subscriptions    []string  `json:"subscriptions"`
	IsActive         bool      `json:"is_active"`
	MessagesSent     int64     `json:"messages_sent"`
	MessagesReceived int64     `json:"messages_received"`
}

// NewRealtimeDashboard creates a new real-time dashboard
func NewRealtimeDashboard(logger *observability.Logger, config DashboardConfig) *RealtimeDashboard {
	// Set default values
	if config.MetricsUpdateInterval == 0 {
		config.MetricsUpdateInterval = time.Second
	}
	if config.AlertCheckInterval == 0 {
		config.AlertCheckInterval = 5 * time.Second
	}
	if config.RefreshRate == 0 {
		config.RefreshRate = 100 * time.Millisecond
	}
	if config.MaxDataPoints == 0 {
		config.MaxDataPoints = 1000
	}
	if config.WSPort == 0 {
		config.WSPort = 8090
	}
	if config.WSPath == "" {
		config.WSPath = "/ws"
	}

	rd := &RealtimeDashboard{
		logger:            logger,
		config:            config,
		dashboardSessions: make(map[string]*DashboardSession),
		widgets:           make(map[string]*Widget),
		layouts:           make(map[string]*Layout),
		wsConnections:     make(map[string]*WSConnection),
		alertStream:       make(chan *DashboardAlert, 1000),
		metricsStream:     make(chan *MetricUpdate, 10000),
		eventStream:       make(chan *DashboardEvent, 1000),
		stopChan:          make(chan struct{}),
	}

	// Initialize components
	rd.metricsAggregator = NewMetricsAggregator(logger, config)
	rd.alertManager = NewDashboardAlertManager(logger, config)
	rd.performanceTracker = NewPerformanceTracker(logger, config)
	rd.visualizationEngine = NewVisualizationEngine(logger, config)
	rd.liveMetrics = &LiveMetrics{}

	// Initialize default widgets and layouts
	rd.initializeDefaultWidgets()
	rd.initializeDefaultLayouts()

	return rd
}

// Start begins the real-time dashboard system
func (rd *RealtimeDashboard) Start(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&rd.isRunning, 0, 1) {
		return fmt.Errorf("real-time dashboard is already running")
	}

	rd.logger.Info(ctx, "Starting real-time dashboard system", map[string]interface{}{
		"metrics_interval": rd.config.MetricsUpdateInterval.String(),
		"alert_interval":   rd.config.AlertCheckInterval.String(),
		"refresh_rate":     rd.config.RefreshRate.String(),
		"ws_port":          rd.config.WSPort,
		"max_data_points":  rd.config.MaxDataPoints,
	})

	// Start processing threads
	rd.wg.Add(5)
	go rd.processMetrics(ctx)
	go rd.processAlerts(ctx)
	go rd.processEvents(ctx)
	go rd.updateLiveMetrics(ctx)
	go rd.performanceMonitor(ctx)

	rd.logger.Info(ctx, "Real-time dashboard system started successfully", nil)
	return nil
}

// Stop gracefully shuts down the real-time dashboard system
func (rd *RealtimeDashboard) Stop(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&rd.isRunning, 1, 0) {
		return fmt.Errorf("real-time dashboard is not running")
	}

	rd.logger.Info(ctx, "Stopping real-time dashboard system", nil)

	close(rd.stopChan)
	rd.wg.Wait()

	// Close all WebSocket connections
	rd.closeAllConnections()

	rd.logger.Info(ctx, "Real-time dashboard system stopped", map[string]interface{}{
		"metrics_processed": atomic.LoadInt64(&rd.metricsProcessed),
		"alerts_generated":  atomic.LoadInt64(&rd.alertsGenerated),
		"sessions_active":   atomic.LoadInt64(&rd.sessionsActive),
		"avg_update_time":   atomic.LoadInt64(&rd.avgUpdateTime),
	})

	return nil
}

// UpdateMetric updates a metric in real-time
func (rd *RealtimeDashboard) UpdateMetric(ctx context.Context, metricName string, value interface{}, source string, tags map[string]string) error {
	if atomic.LoadInt32(&rd.isRunning) != 1 {
		return fmt.Errorf("dashboard is not running")
	}

	update := &MetricUpdate{
		ID:         uuid.New(),
		MetricName: metricName,
		Value:      value,
		Timestamp:  time.Now(),
		Source:     source,
		Tags:       tags,
	}

	select {
	case rd.metricsStream <- update:
		atomic.AddInt64(&rd.metricsProcessed, 1)
		return nil
	default:
		return fmt.Errorf("metrics stream is full")
	}
}

// TriggerAlert triggers a dashboard alert
func (rd *RealtimeDashboard) TriggerAlert(ctx context.Context, alertType AlertType, severity AlertSeverity, title, message, source string, data map[string]interface{}) error {
	if atomic.LoadInt32(&rd.isRunning) != 1 {
		return fmt.Errorf("dashboard is not running")
	}

	alert := &DashboardAlert{
		ID:        uuid.New(),
		Type:      alertType,
		Severity:  severity,
		Title:     title,
		Message:   message,
		Source:    source,
		Timestamp: time.Now(),
		Data:      data,
	}

	select {
	case rd.alertStream <- alert:
		atomic.AddInt64(&rd.alertsGenerated, 1)
		return nil
	default:
		return fmt.Errorf("alert stream is full")
	}
}

// RecordEvent records a dashboard event
func (rd *RealtimeDashboard) RecordEvent(ctx context.Context, eventType EventType, title, description, source string, data map[string]interface{}, impact EventImpact) error {
	if atomic.LoadInt32(&rd.isRunning) != 1 {
		return fmt.Errorf("dashboard is not running")
	}

	event := &DashboardEvent{
		ID:          uuid.New(),
		Type:        eventType,
		Title:       title,
		Description: description,
		Timestamp:   time.Now(),
		Source:      source,
		Data:        data,
		Impact:      impact,
	}

	select {
	case rd.eventStream <- event:
		return nil
	default:
		return fmt.Errorf("event stream is full")
	}
}

// CreateSession creates a new dashboard session
func (rd *RealtimeDashboard) CreateSession(ctx context.Context, userID string) (*DashboardSession, error) {
	session := &DashboardSession{
		ID:           uuid.New(),
		UserID:       userID,
		StartTime:    time.Now(),
		LastActivity: time.Now(),
		Layout:       rd.getDefaultLayout(),
		Widgets:      rd.getDefaultWidgets(),
		Preferences:  make(map[string]interface{}),
		IsActive:     true,
	}

	rd.mu.Lock()
	rd.dashboardSessions[session.ID.String()] = session
	rd.mu.Unlock()

	atomic.AddInt64(&rd.sessionsActive, 1)

	rd.logger.Info(ctx, "Dashboard session created", map[string]interface{}{
		"session_id": session.ID.String(),
		"user_id":    userID,
	})

	return session, nil
}

// GetSession retrieves a dashboard session
func (rd *RealtimeDashboard) GetSession(sessionID string) *DashboardSession {
	rd.mu.RLock()
	defer rd.mu.RUnlock()

	if session, exists := rd.dashboardSessions[sessionID]; exists {
		return session
	}
	return nil
}

// UpdateSession updates a dashboard session
func (rd *RealtimeDashboard) UpdateSession(ctx context.Context, sessionID string, updates map[string]interface{}) error {
	rd.mu.Lock()
	defer rd.mu.Unlock()

	session, exists := rd.dashboardSessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	session.LastActivity = time.Now()

	// Apply updates
	for key, value := range updates {
		switch key {
		case "layout":
			if layout, ok := value.(*Layout); ok {
				session.Layout = layout
			}
		case "widgets":
			if widgets, ok := value.([]string); ok {
				session.Widgets = widgets
			}
		case "preferences":
			if prefs, ok := value.(map[string]interface{}); ok {
				session.Preferences = prefs
			}
		}
	}

	rd.logger.Debug(ctx, "Dashboard session updated", map[string]interface{}{
		"session_id": sessionID,
		"updates":    len(updates),
	})

	return nil
}

// CloseSession closes a dashboard session
func (rd *RealtimeDashboard) CloseSession(ctx context.Context, sessionID string) error {
	rd.mu.Lock()
	defer rd.mu.Unlock()

	session, exists := rd.dashboardSessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	session.IsActive = false
	delete(rd.dashboardSessions, sessionID)

	atomic.AddInt64(&rd.sessionsActive, -1)

	rd.logger.Info(ctx, "Dashboard session closed", map[string]interface{}{
		"session_id": sessionID,
		"duration":   time.Since(session.StartTime).String(),
	})

	return nil
}

// GetLiveMetrics returns current live metrics
func (rd *RealtimeDashboard) GetLiveMetrics() *LiveMetrics {
	rd.mu.RLock()
	defer rd.mu.RUnlock()

	// Return a copy of current metrics
	metrics := *rd.liveMetrics
	return &metrics
}

// GetAlerts returns recent alerts
func (rd *RealtimeDashboard) GetAlerts(limit int, severity AlertSeverity) []*DashboardAlert {
	return rd.alertManager.GetAlerts(limit, severity)
}

// AcknowledgeAlert acknowledges an alert
func (rd *RealtimeDashboard) AcknowledgeAlert(ctx context.Context, alertID uuid.UUID, acknowledgedBy string) error {
	return rd.alertManager.AcknowledgeAlert(alertID, acknowledgedBy)
}

// ResolveAlert resolves an alert
func (rd *RealtimeDashboard) ResolveAlert(ctx context.Context, alertID uuid.UUID) error {
	return rd.alertManager.ResolveAlert(alertID)
}

// processMetrics processes metric updates
func (rd *RealtimeDashboard) processMetrics(ctx context.Context) {
	defer rd.wg.Done()

	rd.logger.Info(ctx, "Starting metrics processor", nil)

	for {
		select {
		case <-rd.stopChan:
			return
		case update := <-rd.metricsStream:
			rd.metricsAggregator.ProcessMetric(update)
			rd.broadcastMetricUpdate(update)
		}
	}
}

// processAlerts processes alert updates
func (rd *RealtimeDashboard) processAlerts(ctx context.Context) {
	defer rd.wg.Done()

	rd.logger.Info(ctx, "Starting alerts processor", nil)

	for {
		select {
		case <-rd.stopChan:
			return
		case alert := <-rd.alertStream:
			rd.alertManager.AddAlert(alert)
			rd.broadcastAlert(alert)
		}
	}
}

// processEvents processes dashboard events
func (rd *RealtimeDashboard) processEvents(ctx context.Context) {
	defer rd.wg.Done()

	rd.logger.Info(ctx, "Starting events processor", nil)

	for {
		select {
		case <-rd.stopChan:
			return
		case event := <-rd.eventStream:
			rd.broadcastEvent(event)
		}
	}
}

// updateLiveMetrics updates live metrics periodically
func (rd *RealtimeDashboard) updateLiveMetrics(ctx context.Context) {
	defer rd.wg.Done()

	rd.logger.Info(ctx, "Starting live metrics updater", nil)

	ticker := time.NewTicker(rd.config.MetricsUpdateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-rd.stopChan:
			return
		case <-ticker.C:
			rd.refreshLiveMetrics()
		}
	}
}

// refreshLiveMetrics refreshes the live metrics data
func (rd *RealtimeDashboard) refreshLiveMetrics() {
	rd.mu.Lock()
	defer rd.mu.Unlock()

	// Update system health
	rd.liveMetrics.SystemHealth = SystemHealth{
		Status:            "HEALTHY",
		Uptime:            time.Since(time.Now().Add(-24 * time.Hour)), // Mock uptime
		CPUUsage:          45.5,
		MemoryUsage:       67.2,
		DiskUsage:         23.8,
		ActiveConnections: len(rd.wsConnections),
		ErrorRate:         0.1,
		LastHealthCheck:   time.Now(),
	}

	// Update performance metrics
	rd.liveMetrics.PerformanceMetrics = *rd.performanceTracker.GetPerformanceMetrics()

	// Update trading metrics (mock data)
	rd.liveMetrics.TradingMetrics = TradingMetrics{
		OrdersPerSecond:  125.5,
		FillRate:         98.5,
		AvgOrderSize:     decimal.NewFromFloat(1000.0),
		AvgExecutionTime: 5 * time.Millisecond,
		AvgSlippage:      2.3,
		SuccessRate:      99.2,
		DailyPnL:         decimal.NewFromFloat(12345.67),
		UnrealizedPnL:    decimal.NewFromFloat(2345.89),
		TotalVolume:      decimal.NewFromFloat(5000000.0),
		ActiveStrategies: 5,
		TopPerformer:     "Strategy_A",
		WorstPerformer:   "Strategy_C",
	}

	// Update network metrics (mock data)
	rd.liveMetrics.NetworkMetrics = NetworkMetrics{
		PacketsReceived:   12345,
		PacketsSent:       11234,
		BytesReceived:     1234567,
		BytesSent:         1123456,
		AvgLatencyNs:      5000000,  // 5ms
		MinLatencyNs:      1000000,  // 1ms
		MaxLatencyNs:      10000000, // 10ms
		ActiveConnections: len(rd.wsConnections),
		LastUpdate:        time.Now(),
	}

	// Update risk metrics (mock data)
	rd.liveMetrics.RiskMetrics = RiskMetrics{
		TotalValue:    decimal.NewFromFloat(10000000.0),
		TotalExposure: decimal.NewFromFloat(8500000.0),
		Concentration: 15.5,
		VaR95:         decimal.NewFromFloat(50000.0),
		MaxDrawdown:   5.2,
		LastUpdate:    time.Now(),
	}

	rd.liveMetrics.LastUpdate = time.Now()
	rd.liveMetrics.UpdateFrequency = rd.config.MetricsUpdateInterval
}

// performanceMonitor tracks dashboard performance
func (rd *RealtimeDashboard) performanceMonitor(ctx context.Context) {
	defer rd.wg.Done()

	rd.logger.Info(ctx, "Starting dashboard performance monitor", nil)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-rd.stopChan:
			return
		case <-ticker.C:
			metricsProcessed := atomic.LoadInt64(&rd.metricsProcessed)
			alertsGenerated := atomic.LoadInt64(&rd.alertsGenerated)
			sessionsActive := atomic.LoadInt64(&rd.sessionsActive)
			avgUpdateTime := atomic.LoadInt64(&rd.avgUpdateTime)

			rd.logger.Info(ctx, "Dashboard performance metrics", map[string]interface{}{
				"metrics_processed": metricsProcessed,
				"alerts_generated":  alertsGenerated,
				"sessions_active":   sessionsActive,
				"avg_update_time":   avgUpdateTime,
				"ws_connections":    len(rd.wsConnections),
				"active_widgets":    len(rd.widgets),
				"active_layouts":    len(rd.layouts),
			})
		}
	}
}

// broadcastMetricUpdate broadcasts metric updates to WebSocket connections
func (rd *RealtimeDashboard) broadcastMetricUpdate(update *MetricUpdate) {
	message := map[string]interface{}{
		"type":    "metric_update",
		"payload": update,
	}

	rd.broadcastToSubscribers("metrics", message)
}

// broadcastAlert broadcasts alerts to WebSocket connections
func (rd *RealtimeDashboard) broadcastAlert(alert *DashboardAlert) {
	message := map[string]interface{}{
		"type":    "alert",
		"payload": alert,
	}

	rd.broadcastToSubscribers("alerts", message)
}

// broadcastEvent broadcasts events to WebSocket connections
func (rd *RealtimeDashboard) broadcastEvent(event *DashboardEvent) {
	message := map[string]interface{}{
		"type":    "event",
		"payload": event,
	}

	rd.broadcastToSubscribers("events", message)
}

// broadcastToSubscribers broadcasts a message to subscribers
func (rd *RealtimeDashboard) broadcastToSubscribers(channel string, message map[string]interface{}) {
	rd.wsMu.RLock()
	defer rd.wsMu.RUnlock()

	_, err := json.Marshal(message)
	if err != nil {
		rd.logger.Error(context.Background(), "Failed to marshal broadcast message", err)
		return
	}

	for _, conn := range rd.wsConnections {
		// Check if connection is subscribed to this channel
		subscribed := false
		for _, sub := range conn.Subscriptions {
			if sub == channel || sub == "*" {
				subscribed = true
				break
			}
		}

		if subscribed && conn.IsActive {
			// In production, this would send via WebSocket
			// For now, just increment message count
			atomic.AddInt64(&conn.MessagesSent, 1)
		}
	}
}

// closeAllConnections closes all WebSocket connections
func (rd *RealtimeDashboard) closeAllConnections() {
	rd.wsMu.Lock()
	defer rd.wsMu.Unlock()

	for id, conn := range rd.wsConnections {
		conn.IsActive = false
		delete(rd.wsConnections, id)
	}
}

// GetWidgets returns all widgets
func (rd *RealtimeDashboard) GetWidgets() map[string]*Widget {
	rd.mu.RLock()
	defer rd.mu.RUnlock()

	widgets := make(map[string]*Widget)
	for id, widget := range rd.widgets {
		widgets[id] = widget
	}
	return widgets
}

// GetLayouts returns all layouts
func (rd *RealtimeDashboard) GetLayouts() map[string]*Layout {
	rd.mu.RLock()
	defer rd.mu.RUnlock()

	layouts := make(map[string]*Layout)
	for id, layout := range rd.layouts {
		layouts[id] = layout
	}
	return layouts
}

// GetDashboardStatus returns dashboard status
func (rd *RealtimeDashboard) GetDashboardStatus() map[string]interface{} {
	return map[string]interface{}{
		"status":            "operational",
		"uptime":            time.Since(time.Now().Add(-24 * time.Hour)).String(),
		"metrics_processed": atomic.LoadInt64(&rd.metricsProcessed),
		"alerts_generated":  atomic.LoadInt64(&rd.alertsGenerated),
		"sessions_active":   atomic.LoadInt64(&rd.sessionsActive),
		"ws_connections":    len(rd.wsConnections),
		"widgets":           len(rd.widgets),
		"layouts":           len(rd.layouts),
		"last_update":       time.Now(),
	}
}
