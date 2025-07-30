package analytics

import (
	"context"
	"runtime"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
)

// SystemPerformanceAnalyzer analyzes system performance metrics
type SystemPerformanceAnalyzer struct {
	logger    *observability.Logger
	config    PerformanceConfig
	metrics   SystemPerformanceMetrics
	history   []SystemPerformanceSnapshot
	mu        sync.RWMutex
	isRunning int32
	stopChan  chan struct{}
}

// SystemPerformanceSnapshot represents a point-in-time system performance snapshot
type SystemPerformanceSnapshot struct {
	Timestamp         time.Time     `json:"timestamp"`
	CPUUsage          float64       `json:"cpu_usage"`
	MemoryUsage       float64       `json:"memory_usage"`
	DiskUsage         float64       `json:"disk_usage"`
	NetworkLatency    time.Duration `json:"network_latency"`
	DatabaseLatency   time.Duration `json:"database_latency"`
	APILatency        time.Duration `json:"api_latency"`
	Throughput        float64       `json:"throughput"`
	ErrorRate         float64       `json:"error_rate"`
	ActiveConnections int           `json:"active_connections"`
	QueueDepth        int           `json:"queue_depth"`
	CacheHitRate      float64       `json:"cache_hit_rate"`
	GCPauseTime       time.Duration `json:"gc_pause_time"`
	GoroutineCount    int           `json:"goroutine_count"`
	HeapSize          int64         `json:"heap_size"`
	AllocRate         float64       `json:"alloc_rate"`
}

// PerformanceTrend represents performance trend analysis
type PerformanceTrend struct {
	Metric        string    `json:"metric"`
	Direction     string    `json:"direction"` // "up", "down", "stable"
	ChangePercent float64   `json:"change_percent"`
	Confidence    float64   `json:"confidence"`
	Period        string    `json:"period"`
	LastUpdated   time.Time `json:"last_updated"`
}

// PerformanceAlert represents a performance alert
type PerformanceAlert struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Severity    string    `json:"severity"`
	Message     string    `json:"message"`
	Metric      string    `json:"metric"`
	Value       float64   `json:"value"`
	Threshold   float64   `json:"threshold"`
	Timestamp   time.Time `json:"timestamp"`
	Acknowledged bool     `json:"acknowledged"`
}

// NewSystemPerformanceAnalyzer creates a new system performance analyzer
func NewSystemPerformanceAnalyzer(logger *observability.Logger, config PerformanceConfig) *SystemPerformanceAnalyzer {
	return &SystemPerformanceAnalyzer{
		logger:   logger,
		config:   config,
		history:  make([]SystemPerformanceSnapshot, 0, config.MetricsBufferSize),
		stopChan: make(chan struct{}),
	}
}

// Start starts the system performance analyzer
func (spa *SystemPerformanceAnalyzer) Start(ctx context.Context) error {
	spa.logger.Info(ctx, "Starting system performance analyzer", nil)
	spa.isRunning = 1

	// Start metrics collection loop
	go spa.metricsCollectionLoop(ctx)

	return nil
}

// Stop stops the system performance analyzer
func (spa *SystemPerformanceAnalyzer) Stop(ctx context.Context) error {
	spa.logger.Info(ctx, "Stopping system performance analyzer", nil)
	spa.isRunning = 0
	close(spa.stopChan)
	return nil
}

// metricsCollectionLoop runs the metrics collection loop
func (spa *SystemPerformanceAnalyzer) metricsCollectionLoop(ctx context.Context) {
	ticker := time.NewTicker(spa.config.AnalysisInterval)
	defer ticker.Stop()

	for {
		select {
		case <-spa.stopChan:
			return
		case <-ticker.C:
			spa.collectMetrics(ctx)
		}
	}
}

// collectMetrics collects current system performance metrics
func (spa *SystemPerformanceAnalyzer) collectMetrics(ctx context.Context) {
	spa.mu.Lock()
	defer spa.mu.Unlock()

	snapshot := spa.captureSnapshot()
	
	// Add to history
	spa.history = append(spa.history, snapshot)
	
	// Maintain buffer size
	if len(spa.history) > spa.config.MetricsBufferSize {
		spa.history = spa.history[1:]
	}

	// Update current metrics
	spa.updateCurrentMetrics(snapshot)

	spa.logger.Debug(ctx, "System metrics collected", map[string]interface{}{
		"cpu_usage":      snapshot.CPUUsage,
		"memory_usage":   snapshot.MemoryUsage,
		"goroutines":     snapshot.GoroutineCount,
		"heap_size_mb":   snapshot.HeapSize / 1024 / 1024,
		"api_latency_ms": snapshot.APILatency.Milliseconds(),
	})
}

// captureSnapshot captures a current system performance snapshot
func (spa *SystemPerformanceAnalyzer) captureSnapshot() SystemPerformanceSnapshot {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	snapshot := SystemPerformanceSnapshot{
		Timestamp:         time.Now(),
		CPUUsage:          spa.getCPUUsage(),
		MemoryUsage:       spa.getMemoryUsage(&m),
		DiskUsage:         spa.getDiskUsage(),
		NetworkLatency:    spa.getNetworkLatency(),
		DatabaseLatency:   spa.getDatabaseLatency(),
		APILatency:        spa.getAPILatency(),
		Throughput:        spa.getThroughput(),
		ErrorRate:         spa.getErrorRate(),
		ActiveConnections: spa.getActiveConnections(),
		QueueDepth:        spa.getQueueDepth(),
		CacheHitRate:      spa.getCacheHitRate(),
		GCPauseTime:       time.Duration(m.PauseTotalNs),
		GoroutineCount:    runtime.NumGoroutine(),
		HeapSize:          int64(m.HeapAlloc),
		AllocRate:         spa.getAllocRate(&m),
	}

	return snapshot
}

// getCPUUsage gets current CPU usage percentage
func (spa *SystemPerformanceAnalyzer) getCPUUsage() float64 {
	// Simplified CPU usage calculation
	// In production, use proper CPU monitoring
	return float64(runtime.NumGoroutine()) / 1000.0 * 100
}

// getMemoryUsage gets current memory usage percentage
func (spa *SystemPerformanceAnalyzer) getMemoryUsage(m *runtime.MemStats) float64 {
	// Calculate memory usage as percentage of heap
	if m.HeapSys == 0 {
		return 0
	}
	return float64(m.HeapAlloc) / float64(m.HeapSys) * 100
}

// getDiskUsage gets current disk usage percentage
func (spa *SystemPerformanceAnalyzer) getDiskUsage() float64 {
	// Mock implementation - in production, check actual disk usage
	return 45.0
}

// getNetworkLatency gets current network latency
func (spa *SystemPerformanceAnalyzer) getNetworkLatency() time.Duration {
	// Mock implementation - in production, measure actual network latency
	return 5 * time.Millisecond
}

// getDatabaseLatency gets current database latency
func (spa *SystemPerformanceAnalyzer) getDatabaseLatency() time.Duration {
	// Mock implementation - in production, measure actual database latency
	return 2 * time.Millisecond
}

// getAPILatency gets current API latency
func (spa *SystemPerformanceAnalyzer) getAPILatency() time.Duration {
	// Mock implementation - in production, measure actual API latency
	return 10 * time.Millisecond
}

// getThroughput gets current throughput (requests per second)
func (spa *SystemPerformanceAnalyzer) getThroughput() float64 {
	// Mock implementation - in production, calculate actual throughput
	return 1500.0
}

// getErrorRate gets current error rate percentage
func (spa *SystemPerformanceAnalyzer) getErrorRate() float64 {
	// Mock implementation - in production, calculate actual error rate
	return 0.1
}

// getActiveConnections gets current active connections count
func (spa *SystemPerformanceAnalyzer) getActiveConnections() int {
	// Mock implementation - in production, count actual connections
	return 150
}

// getQueueDepth gets current queue depth
func (spa *SystemPerformanceAnalyzer) getQueueDepth() int {
	// Mock implementation - in production, measure actual queue depth
	return 25
}

// getCacheHitRate gets current cache hit rate percentage
func (spa *SystemPerformanceAnalyzer) getCacheHitRate() float64 {
	// Mock implementation - in production, calculate actual cache hit rate
	return 95.5
}

// getAllocRate gets current allocation rate
func (spa *SystemPerformanceAnalyzer) getAllocRate(m *runtime.MemStats) float64 {
	// Calculate allocation rate (simplified)
	return float64(m.Mallocs-m.Frees) / 1000.0
}

// updateCurrentMetrics updates the current metrics from snapshot
func (spa *SystemPerformanceAnalyzer) updateCurrentMetrics(snapshot SystemPerformanceSnapshot) {
	// Calculate uptime
	var uptime time.Duration
	if len(spa.history) > 0 {
		uptime = snapshot.Timestamp.Sub(spa.history[0].Timestamp)
	}

	spa.metrics = SystemPerformanceMetrics{
		CPUUsage:          snapshot.CPUUsage,
		MemoryUsage:       snapshot.MemoryUsage,
		DiskUsage:         snapshot.DiskUsage,
		NetworkLatency:    snapshot.NetworkLatency,
		DatabaseLatency:   snapshot.DatabaseLatency,
		APILatency:        snapshot.APILatency,
		Throughput:        snapshot.Throughput,
		ErrorRate:         snapshot.ErrorRate,
		Uptime:            uptime,
		ActiveConnections: snapshot.ActiveConnections,
		QueueDepth:        snapshot.QueueDepth,
		CacheHitRate:      snapshot.CacheHitRate,
		GCPauseTime:       snapshot.GCPauseTime,
		GoroutineCount:    snapshot.GoroutineCount,
		HeapSize:          snapshot.HeapSize,
		AllocRate:         snapshot.AllocRate,
	}
}

// GetMetrics returns current system performance metrics
func (spa *SystemPerformanceAnalyzer) GetMetrics() SystemPerformanceMetrics {
	spa.mu.RLock()
	defer spa.mu.RUnlock()
	return spa.metrics
}

// GetHistory returns system performance history
func (spa *SystemPerformanceAnalyzer) GetHistory(duration time.Duration) []SystemPerformanceSnapshot {
	spa.mu.RLock()
	defer spa.mu.RUnlock()

	cutoff := time.Now().Add(-duration)
	var filtered []SystemPerformanceSnapshot

	for _, snapshot := range spa.history {
		if snapshot.Timestamp.After(cutoff) {
			filtered = append(filtered, snapshot)
		}
	}

	return filtered
}

// GetTrends analyzes performance trends
func (spa *SystemPerformanceAnalyzer) GetTrends(period time.Duration) []PerformanceTrend {
	spa.mu.RLock()
	defer spa.mu.RUnlock()

	if len(spa.history) < 2 {
		return []PerformanceTrend{}
	}

	cutoff := time.Now().Add(-period)
	var recentHistory []SystemPerformanceSnapshot

	for _, snapshot := range spa.history {
		if snapshot.Timestamp.After(cutoff) {
			recentHistory = append(recentHistory, snapshot)
		}
	}

	if len(recentHistory) < 2 {
		return []PerformanceTrend{}
	}

	trends := []PerformanceTrend{
		spa.calculateTrend("cpu_usage", recentHistory, func(s SystemPerformanceSnapshot) float64 { return s.CPUUsage }),
		spa.calculateTrend("memory_usage", recentHistory, func(s SystemPerformanceSnapshot) float64 { return s.MemoryUsage }),
		spa.calculateTrend("disk_usage", recentHistory, func(s SystemPerformanceSnapshot) float64 { return s.DiskUsage }),
		spa.calculateTrend("throughput", recentHistory, func(s SystemPerformanceSnapshot) float64 { return s.Throughput }),
		spa.calculateTrend("error_rate", recentHistory, func(s SystemPerformanceSnapshot) float64 { return s.ErrorRate }),
		spa.calculateTrend("api_latency", recentHistory, func(s SystemPerformanceSnapshot) float64 { return float64(s.APILatency.Milliseconds()) }),
	}

	return trends
}

// calculateTrend calculates trend for a specific metric
func (spa *SystemPerformanceAnalyzer) calculateTrend(metricName string, history []SystemPerformanceSnapshot, extractor func(SystemPerformanceSnapshot) float64) PerformanceTrend {
	if len(history) < 2 {
		return PerformanceTrend{
			Metric:      metricName,
			Direction:   "stable",
			Confidence:  0,
			LastUpdated: time.Now(),
		}
	}

	// Calculate simple linear trend
	n := len(history)
	firstValue := extractor(history[0])
	lastValue := extractor(history[n-1])

	changePercent := 0.0
	if firstValue != 0 {
		changePercent = (lastValue - firstValue) / firstValue * 100
	}

	direction := "stable"
	if changePercent > 5 {
		direction = "up"
	} else if changePercent < -5 {
		direction = "down"
	}

	// Calculate confidence based on consistency of trend
	confidence := spa.calculateTrendConfidence(history, extractor)

	return PerformanceTrend{
		Metric:        metricName,
		Direction:     direction,
		ChangePercent: changePercent,
		Confidence:    confidence,
		Period:        "recent",
		LastUpdated:   time.Now(),
	}
}

// calculateTrendConfidence calculates confidence in trend analysis
func (spa *SystemPerformanceAnalyzer) calculateTrendConfidence(history []SystemPerformanceSnapshot, extractor func(SystemPerformanceSnapshot) float64) float64 {
	if len(history) < 3 {
		return 0.5
	}

	// Calculate variance to determine confidence
	values := make([]float64, len(history))
	var sum float64
	for i, snapshot := range history {
		values[i] = extractor(snapshot)
		sum += values[i]
	}

	mean := sum / float64(len(values))
	var variance float64
	for _, value := range values {
		variance += (value - mean) * (value - mean)
	}
	variance /= float64(len(values))

	// Lower variance = higher confidence
	// Normalize to 0-1 range
	confidence := 1.0 / (1.0 + variance/mean)
	if confidence > 1.0 {
		confidence = 1.0
	}
	if confidence < 0.0 {
		confidence = 0.0
	}

	return confidence
}

// GetPerformanceScore calculates overall system performance score
func (spa *SystemPerformanceAnalyzer) GetPerformanceScore() float64 {
	spa.mu.RLock()
	defer spa.mu.RUnlock()

	score := 100.0

	// Deduct points for high resource usage
	if spa.metrics.CPUUsage > 80 {
		score -= (spa.metrics.CPUUsage - 80) * 0.5
	}
	if spa.metrics.MemoryUsage > 85 {
		score -= (spa.metrics.MemoryUsage - 85) * 0.5
	}
	if spa.metrics.DiskUsage > 90 {
		score -= (spa.metrics.DiskUsage - 90) * 1.0
	}

	// Deduct points for high latency
	if spa.metrics.APILatency > 100*time.Millisecond {
		score -= float64(spa.metrics.APILatency.Milliseconds()-100) * 0.1
	}

	// Deduct points for high error rate
	if spa.metrics.ErrorRate > 1.0 {
		score -= (spa.metrics.ErrorRate - 1.0) * 10
	}

	// Deduct points for low cache hit rate
	if spa.metrics.CacheHitRate < 90 {
		score -= (90 - spa.metrics.CacheHitRate) * 0.5
	}

	// Ensure score is between 0 and 100
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

// GetResourceUtilization calculates resource utilization efficiency
func (spa *SystemPerformanceAnalyzer) GetResourceUtilization() map[string]float64 {
	spa.mu.RLock()
	defer spa.mu.RUnlock()

	return map[string]float64{
		"cpu":     spa.metrics.CPUUsage,
		"memory":  spa.metrics.MemoryUsage,
		"disk":    spa.metrics.DiskUsage,
		"network": spa.calculateNetworkUtilization(),
		"cache":   spa.metrics.CacheHitRate,
	}
}

// calculateNetworkUtilization calculates network utilization
func (spa *SystemPerformanceAnalyzer) calculateNetworkUtilization() float64 {
	// Simplified calculation based on latency and throughput
	latencyScore := 100.0 - float64(spa.metrics.NetworkLatency.Milliseconds())
	if latencyScore < 0 {
		latencyScore = 0
	}
	
	throughputScore := spa.metrics.Throughput / 2000.0 * 100 // Assume max 2000 RPS
	if throughputScore > 100 {
		throughputScore = 100
	}

	return (latencyScore + throughputScore) / 2
}
