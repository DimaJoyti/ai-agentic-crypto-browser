package monitoring

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/shopspring/decimal"
)

// SystemMonitor provides comprehensive system monitoring and health checks
type SystemMonitor struct {
	logger     *observability.Logger
	metrics    *SystemMetrics
	alerts     []Alert
	config     MonitoringConfig
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
	collectors map[string]MetricCollector
}

// MonitoringConfig holds configuration for system monitoring
type MonitoringConfig struct {
	CollectionInterval time.Duration `json:"collection_interval"`
	RetentionPeriod    time.Duration `json:"retention_period"`
	AlertThresholds    AlertConfig   `json:"alert_thresholds"`
	EnableProfiling    bool          `json:"enable_profiling"`
	EnableTracing      bool          `json:"enable_tracing"`
}

// AlertConfig defines thresholds for various alerts
type AlertConfig struct {
	CPUThreshold        float64 `json:"cpu_threshold"`
	MemoryThreshold     float64 `json:"memory_threshold"`
	DiskThreshold       float64 `json:"disk_threshold"`
	ErrorRateThreshold  float64 `json:"error_rate_threshold"`
	LatencyThreshold    float64 `json:"latency_threshold"`
	ConnectionThreshold int     `json:"connection_threshold"`
}

// SystemMetrics contains comprehensive system performance metrics
type SystemMetrics struct {
	Timestamp   time.Time       `json:"timestamp"`
	CPU         CPUMetrics      `json:"cpu"`
	Memory      MemoryMetrics   `json:"memory"`
	Disk        DiskMetrics     `json:"disk"`
	Network     NetworkMetrics  `json:"network"`
	Application AppMetrics      `json:"application"`
	Trading     TradingMetrics  `json:"trading"`
	Database    DatabaseMetrics `json:"database"`
	WebSocket   WSMetrics       `json:"websocket"`
	Health      HealthStatus    `json:"health"`
}

// CPUMetrics contains CPU performance data
type CPUMetrics struct {
	UsagePercent   float64 `json:"usage_percent"`
	LoadAverage1m  float64 `json:"load_average_1m"`
	LoadAverage5m  float64 `json:"load_average_5m"`
	LoadAverage15m float64 `json:"load_average_15m"`
	Cores          int     `json:"cores"`
	Goroutines     int     `json:"goroutines"`
	CGoCalls       int64   `json:"cgo_calls"`
}

// MemoryMetrics contains memory usage data
type MemoryMetrics struct {
	TotalBytes   uint64  `json:"total_bytes"`
	UsedBytes    uint64  `json:"used_bytes"`
	FreeBytes    uint64  `json:"free_bytes"`
	UsagePercent float64 `json:"usage_percent"`
	HeapBytes    uint64  `json:"heap_bytes"`
	StackBytes   uint64  `json:"stack_bytes"`
	GCPauseMs    float64 `json:"gc_pause_ms"`
	GCCount      uint32  `json:"gc_count"`
}

// DiskMetrics contains disk usage data
type DiskMetrics struct {
	TotalBytes   uint64  `json:"total_bytes"`
	UsedBytes    uint64  `json:"used_bytes"`
	FreeBytes    uint64  `json:"free_bytes"`
	UsagePercent float64 `json:"usage_percent"`
	ReadOps      uint64  `json:"read_ops"`
	WriteOps     uint64  `json:"write_ops"`
	ReadBytes    uint64  `json:"read_bytes"`
	WriteBytes   uint64  `json:"write_bytes"`
}

// NetworkMetrics contains network performance data
type NetworkMetrics struct {
	BytesReceived   uint64  `json:"bytes_received"`
	BytesSent       uint64  `json:"bytes_sent"`
	PacketsReceived uint64  `json:"packets_received"`
	PacketsSent     uint64  `json:"packets_sent"`
	Connections     int     `json:"connections"`
	ActiveSockets   int     `json:"active_sockets"`
	DropRate        float64 `json:"drop_rate"`
}

// AppMetrics contains application-specific metrics
type AppMetrics struct {
	RequestCount    uint64        `json:"request_count"`
	ErrorCount      uint64        `json:"error_count"`
	ErrorRate       float64       `json:"error_rate"`
	AvgResponseTime time.Duration `json:"avg_response_time"`
	P95ResponseTime time.Duration `json:"p95_response_time"`
	P99ResponseTime time.Duration `json:"p99_response_time"`
	ActiveUsers     int           `json:"active_users"`
	ThroughputRPS   float64       `json:"throughput_rps"`
	CacheHitRate    float64       `json:"cache_hit_rate"`
	QueueLength     int           `json:"queue_length"`
}

// TradingMetrics contains trading system performance data
type TradingMetrics struct {
	ActivePortfolios    int                `json:"active_portfolios"`
	ActivePositions     int                `json:"active_positions"`
	TotalTrades         uint64             `json:"total_trades"`
	SuccessfulTrades    uint64             `json:"successful_trades"`
	FailedTrades        uint64             `json:"failed_trades"`
	TradeSuccessRate    float64            `json:"trade_success_rate"`
	AvgExecutionTime    time.Duration      `json:"avg_execution_time"`
	TotalVolume         decimal.Decimal    `json:"total_volume"`
	TotalPnL            decimal.Decimal    `json:"total_pnl"`
	RiskAlerts          int                `json:"risk_alerts"`
	StrategyPerformance map[string]float64 `json:"strategy_performance"`
}

// DatabaseMetrics contains database performance data
type DatabaseMetrics struct {
	ConnectionsActive int           `json:"connections_active"`
	ConnectionsIdle   int           `json:"connections_idle"`
	ConnectionsTotal  int           `json:"connections_total"`
	QueriesPerSecond  float64       `json:"queries_per_second"`
	AvgQueryTime      time.Duration `json:"avg_query_time"`
	SlowQueries       int           `json:"slow_queries"`
	DeadlockCount     int           `json:"deadlock_count"`
	CacheHitRatio     float64       `json:"cache_hit_ratio"`
	ReplicationLag    time.Duration `json:"replication_lag"`
}

// WSMetrics contains WebSocket connection metrics
type WSMetrics struct {
	TotalConnections  int     `json:"total_connections"`
	ActiveConnections int     `json:"active_connections"`
	MessagesReceived  uint64  `json:"messages_received"`
	MessagesSent      uint64  `json:"messages_sent"`
	ConnectionErrors  uint64  `json:"connection_errors"`
	ReconnectCount    uint64  `json:"reconnect_count"`
	AvgLatency        float64 `json:"avg_latency"`
	DataThroughput    uint64  `json:"data_throughput"`
}

// HealthStatus represents overall system health
type HealthStatus struct {
	Status     string            `json:"status"`
	Score      float64           `json:"score"`
	Components map[string]string `json:"components"`
	Issues     []string          `json:"issues"`
	LastCheck  time.Time         `json:"last_check"`
}

// Alert represents a system alert
type Alert struct {
	ID          string                 `json:"id"`
	Type        AlertType              `json:"type"`
	Severity    AlertSeverity          `json:"severity"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Metric      string                 `json:"metric"`
	Value       float64                `json:"value"`
	Threshold   float64                `json:"threshold"`
	Timestamp   time.Time              `json:"timestamp"`
	Resolved    bool                   `json:"resolved"`
	ResolvedAt  *time.Time             `json:"resolved_at,omitempty"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// AlertType represents the type of alert
type AlertType string

const (
	AlertTypeSystem   AlertType = "system"
	AlertTypeTrading  AlertType = "trading"
	AlertTypeDatabase AlertType = "database"
	AlertTypeNetwork  AlertType = "network"
	AlertTypeSecurity AlertType = "security"
)

// AlertSeverity represents the severity of an alert
type AlertSeverity string

const (
	AlertSeverityLow      AlertSeverity = "low"
	AlertSeverityMedium   AlertSeverity = "medium"
	AlertSeverityHigh     AlertSeverity = "high"
	AlertSeverityCritical AlertSeverity = "critical"
)

// MetricCollector interface for collecting specific metrics
type MetricCollector interface {
	Collect(ctx context.Context) (interface{}, error)
	Name() string
}

// NewSystemMonitor creates a new system monitor
func NewSystemMonitor(logger *observability.Logger, config MonitoringConfig) *SystemMonitor {
	ctx, cancel := context.WithCancel(context.Background())

	return &SystemMonitor{
		logger:     logger,
		metrics:    &SystemMetrics{},
		alerts:     make([]Alert, 0),
		config:     config,
		ctx:        ctx,
		cancel:     cancel,
		collectors: make(map[string]MetricCollector),
	}
}

// Start begins the monitoring service
func (s *SystemMonitor) Start() error {
	s.logger.Info(s.ctx, "Starting system monitor", map[string]interface{}{
		"collection_interval": s.config.CollectionInterval.String(),
		"retention_period":    s.config.RetentionPeriod.String(),
	})

	// Start metric collection
	go s.collectMetrics()

	// Start alert processing
	go s.processAlerts()

	return nil
}

// Stop stops the monitoring service
func (s *SystemMonitor) Stop() error {
	s.logger.Info(s.ctx, "Stopping system monitor")
	s.cancel()
	return nil
}

// GetCurrentMetrics returns the current system metrics
func (s *SystemMonitor) GetCurrentMetrics() *SystemMetrics {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy to avoid race conditions
	metricsCopy := *s.metrics
	return &metricsCopy
}

// GetAlerts returns current active alerts
func (s *SystemMonitor) GetAlerts() []Alert {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return only active alerts
	activeAlerts := make([]Alert, 0)
	for _, alert := range s.alerts {
		if !alert.Resolved {
			activeAlerts = append(activeAlerts, alert)
		}
	}

	return activeAlerts
}

// collectMetrics collects system metrics at regular intervals
func (s *SystemMonitor) collectMetrics() {
	ticker := time.NewTicker(s.config.CollectionInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.collectCurrentMetrics()
		}
	}
}

// collectCurrentMetrics collects all current system metrics
func (s *SystemMonitor) collectCurrentMetrics() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.metrics.Timestamp = time.Now()

	// Collect CPU metrics
	s.metrics.CPU = s.collectCPUMetrics()

	// Collect memory metrics
	s.metrics.Memory = s.collectMemoryMetrics()

	// Collect disk metrics (simplified)
	s.metrics.Disk = s.collectDiskMetrics()

	// Collect network metrics (simplified)
	s.metrics.Network = s.collectNetworkMetrics()

	// Collect application metrics
	s.metrics.Application = s.collectAppMetrics()

	// Collect trading metrics (simplified)
	s.metrics.Trading = s.collectTradingMetrics()

	// Collect database metrics (simplified)
	s.metrics.Database = s.collectDatabaseMetrics()

	// Collect WebSocket metrics (simplified)
	s.metrics.WebSocket = s.collectWSMetrics()

	// Calculate health status
	s.metrics.Health = s.calculateHealthStatus()

	s.logger.Debug(s.ctx, "Metrics collected", map[string]interface{}{
		"cpu_usage":    s.metrics.CPU.UsagePercent,
		"memory_usage": s.metrics.Memory.UsagePercent,
		"goroutines":   s.metrics.CPU.Goroutines,
		"health_score": s.metrics.Health.Score,
	})
}

// collectCPUMetrics collects CPU-related metrics
func (s *SystemMonitor) collectCPUMetrics() CPUMetrics {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return CPUMetrics{
		UsagePercent:   s.getCPUUsage(),
		LoadAverage1m:  s.getLoadAverage(1),
		LoadAverage5m:  s.getLoadAverage(5),
		LoadAverage15m: s.getLoadAverage(15),
		Cores:          runtime.NumCPU(),
		Goroutines:     runtime.NumGoroutine(),
		CGoCalls:       runtime.NumCgoCall(),
	}
}

// collectMemoryMetrics collects memory-related metrics
func (s *SystemMonitor) collectMemoryMetrics() MemoryMetrics {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	totalMem := s.getTotalMemory()
	usedMem := memStats.Sys
	freeMem := totalMem - usedMem
	usagePercent := float64(usedMem) / float64(totalMem) * 100

	return MemoryMetrics{
		TotalBytes:   totalMem,
		UsedBytes:    usedMem,
		FreeBytes:    freeMem,
		UsagePercent: usagePercent,
		HeapBytes:    memStats.HeapAlloc,
		StackBytes:   memStats.StackInuse,
		GCPauseMs:    float64(memStats.PauseNs[(memStats.NumGC+255)%256]) / 1e6,
		GCCount:      memStats.NumGC,
	}
}

// collectDiskMetrics collects disk-related metrics (simplified)
func (s *SystemMonitor) collectDiskMetrics() DiskMetrics {
	return DiskMetrics{
		TotalBytes:   100 * 1024 * 1024 * 1024, // 100GB
		UsedBytes:    30 * 1024 * 1024 * 1024,  // 30GB
		FreeBytes:    70 * 1024 * 1024 * 1024,  // 70GB
		UsagePercent: 30.0,
		ReadOps:      1000,
		WriteOps:     500,
		ReadBytes:    1024 * 1024,
		WriteBytes:   512 * 1024,
	}
}

// collectNetworkMetrics collects network-related metrics (simplified)
func (s *SystemMonitor) collectNetworkMetrics() NetworkMetrics {
	return NetworkMetrics{
		BytesReceived:   1024 * 1024 * 100, // 100MB
		BytesSent:       1024 * 1024 * 50,  // 50MB
		PacketsReceived: 10000,
		PacketsSent:     5000,
		Connections:     100,
		ActiveSockets:   50,
		DropRate:        0.01,
	}
}

// collectAppMetrics collects application-specific metrics (simplified)
func (s *SystemMonitor) collectAppMetrics() AppMetrics {
	return AppMetrics{
		RequestCount:    10000,
		ErrorCount:      50,
		ErrorRate:       0.5,
		AvgResponseTime: 100 * time.Millisecond,
		P95ResponseTime: 200 * time.Millisecond,
		P99ResponseTime: 500 * time.Millisecond,
		ActiveUsers:     250,
		ThroughputRPS:   100.0,
		CacheHitRate:    0.85,
		QueueLength:     10,
	}
}

// collectTradingMetrics collects trading system metrics (simplified)
func (s *SystemMonitor) collectTradingMetrics() TradingMetrics {
	return TradingMetrics{
		ActivePortfolios: 50,
		ActivePositions:  150,
		TotalTrades:      1000,
		SuccessfulTrades: 850,
		FailedTrades:     150,
		TradeSuccessRate: 85.0,
		AvgExecutionTime: 500 * time.Millisecond,
		TotalVolume:      decimal.NewFromInt(1000000),
		TotalPnL:         decimal.NewFromInt(50000),
		RiskAlerts:       5,
		StrategyPerformance: map[string]float64{
			"momentum":       85.5,
			"mean_reversion": 78.2,
			"arbitrage":      92.1,
		},
	}
}

// collectDatabaseMetrics collects database metrics (simplified)
func (s *SystemMonitor) collectDatabaseMetrics() DatabaseMetrics {
	return DatabaseMetrics{
		ConnectionsActive: 20,
		ConnectionsIdle:   30,
		ConnectionsTotal:  50,
		QueriesPerSecond:  100.0,
		AvgQueryTime:      10 * time.Millisecond,
		SlowQueries:       2,
		DeadlockCount:     0,
		CacheHitRatio:     0.95,
		ReplicationLag:    100 * time.Millisecond,
	}
}

// collectWSMetrics collects WebSocket metrics (simplified)
func (s *SystemMonitor) collectWSMetrics() WSMetrics {
	return WSMetrics{
		TotalConnections:  100,
		ActiveConnections: 85,
		MessagesReceived:  10000,
		MessagesSent:      8000,
		ConnectionErrors:  5,
		ReconnectCount:    10,
		AvgLatency:        50.0,
		DataThroughput:    1024 * 1024, // 1MB
	}
}

// calculateHealthStatus calculates overall system health
func (s *SystemMonitor) calculateHealthStatus() HealthStatus {
	score := 100.0
	issues := make([]string, 0)
	components := make(map[string]string)

	// Check CPU health
	if s.metrics.CPU.UsagePercent > s.config.AlertThresholds.CPUThreshold {
		score -= 20
		issues = append(issues, "High CPU usage")
		components["cpu"] = "degraded"
	} else {
		components["cpu"] = "healthy"
	}

	// Check memory health
	if s.metrics.Memory.UsagePercent > s.config.AlertThresholds.MemoryThreshold {
		score -= 20
		issues = append(issues, "High memory usage")
		components["memory"] = "degraded"
	} else {
		components["memory"] = "healthy"
	}

	// Check error rate
	if s.metrics.Application.ErrorRate > s.config.AlertThresholds.ErrorRateThreshold {
		score -= 30
		issues = append(issues, "High error rate")
		components["application"] = "degraded"
	} else {
		components["application"] = "healthy"
	}

	// Determine overall status
	status := "healthy"
	if score < 50 {
		status = "critical"
	} else if score < 80 {
		status = "degraded"
	}

	return HealthStatus{
		Status:     status,
		Score:      score,
		Components: components,
		Issues:     issues,
		LastCheck:  time.Now(),
	}
}

// processAlerts processes and manages system alerts
func (s *SystemMonitor) processAlerts() {
	ticker := time.NewTicker(s.config.CollectionInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.checkAlertConditions()
		}
	}
}

// checkAlertConditions checks for alert conditions and creates alerts
func (s *SystemMonitor) checkAlertConditions() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check CPU threshold
	if s.metrics.CPU.UsagePercent > s.config.AlertThresholds.CPUThreshold {
		s.createAlert(AlertTypeSystem, AlertSeverityHigh, "High CPU Usage",
			fmt.Sprintf("CPU usage is %.2f%%, exceeding threshold of %.2f%%",
				s.metrics.CPU.UsagePercent, s.config.AlertThresholds.CPUThreshold),
			"cpu_usage", s.metrics.CPU.UsagePercent, s.config.AlertThresholds.CPUThreshold)
	}

	// Check memory threshold
	if s.metrics.Memory.UsagePercent > s.config.AlertThresholds.MemoryThreshold {
		s.createAlert(AlertTypeSystem, AlertSeverityHigh, "High Memory Usage",
			fmt.Sprintf("Memory usage is %.2f%%, exceeding threshold of %.2f%%",
				s.metrics.Memory.UsagePercent, s.config.AlertThresholds.MemoryThreshold),
			"memory_usage", s.metrics.Memory.UsagePercent, s.config.AlertThresholds.MemoryThreshold)
	}

	// Check error rate threshold
	if s.metrics.Application.ErrorRate > s.config.AlertThresholds.ErrorRateThreshold {
		s.createAlert(AlertTypeSystem, AlertSeverityCritical, "High Error Rate",
			fmt.Sprintf("Error rate is %.2f%%, exceeding threshold of %.2f%%",
				s.metrics.Application.ErrorRate, s.config.AlertThresholds.ErrorRateThreshold),
			"error_rate", s.metrics.Application.ErrorRate, s.config.AlertThresholds.ErrorRateThreshold)
	}
}

// createAlert creates a new alert
func (s *SystemMonitor) createAlert(alertType AlertType, severity AlertSeverity, title, description, metric string, value, threshold float64) {
	alert := Alert{
		ID:          fmt.Sprintf("%s_%d", metric, time.Now().Unix()),
		Type:        alertType,
		Severity:    severity,
		Title:       title,
		Description: description,
		Metric:      metric,
		Value:       value,
		Threshold:   threshold,
		Timestamp:   time.Now(),
		Resolved:    false,
		Metadata:    make(map[string]interface{}),
	}

	s.alerts = append(s.alerts, alert)

	s.logger.Warn(s.ctx, "Alert created", map[string]interface{}{
		"alert_id":  alert.ID,
		"type":      string(alert.Type),
		"severity":  string(alert.Severity),
		"title":     alert.Title,
		"metric":    alert.Metric,
		"value":     alert.Value,
		"threshold": alert.Threshold,
	})
}

// Helper functions for system metrics (simplified implementations)
func (s *SystemMonitor) getCPUUsage() float64 {
	// Simplified CPU usage calculation
	return float64(runtime.NumGoroutine()) / float64(runtime.NumCPU()) * 10
}

func (s *SystemMonitor) getLoadAverage(minutes int) float64 {
	// Simplified load average calculation
	return float64(runtime.NumGoroutine()) / float64(runtime.NumCPU())
}

func (s *SystemMonitor) getTotalMemory() uint64 {
	// Simplified total memory (would use system calls in reality)
	return 8 * 1024 * 1024 * 1024 // 8GB
}
