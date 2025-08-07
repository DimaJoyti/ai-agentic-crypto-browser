package observability

import (
	"context"
	"runtime"
	"runtime/debug"
	"sync"
	"time"
)

// PerformanceMonitor tracks system and application performance metrics
type PerformanceMonitor struct {
	logger   *Logger
	metrics  *PerformanceMetrics
	config   *PerformanceConfig
	stopChan chan struct{}
	mu       sync.RWMutex
}

// PerformanceMetrics contains performance data
type PerformanceMetrics struct {
	// System metrics
	CPUUsage       float64
	MemoryUsage    int64
	GoroutineCount int
	GCStats        debug.GCStats

	// Application metrics
	RequestCount  int64
	ResponseTime  time.Duration
	ErrorRate     float64
	ThroughputRPS float64

	// Database metrics
	DBConnections int64
	DBQueryTime   time.Duration
	DBSlowQueries int64

	// Cache metrics
	CacheHitRate   float64
	CacheSize      int64
	CacheEvictions int64

	// Custom metrics
	CustomMetrics map[string]interface{}

	// Timestamps
	LastUpdated time.Time
	mu          sync.RWMutex
}

// PerformanceConfig contains monitoring configuration
type PerformanceConfig struct {
	CollectionInterval time.Duration
	RetentionPeriod    time.Duration
	AlertThresholds    *AlertThresholds
	EnableProfiling    bool
	EnableTracing      bool
}

// AlertThresholds defines performance alert thresholds
type AlertThresholds struct {
	CPUUsageThreshold     float64
	MemoryUsageThreshold  int64
	ResponseTimeThreshold time.Duration
	ErrorRateThreshold    float64
	GoroutineThreshold    int
}

// RequestMetrics tracks individual request performance
type RequestMetrics struct {
	Path       string
	Method     string
	StatusCode int
	Duration   time.Duration
	Size       int64
	UserAgent  string
	IP         string
	Timestamp  time.Time
}

// NewPerformanceMonitor creates a new performance monitor
func NewPerformanceMonitor(logger *Logger) *PerformanceMonitor {
	config := &PerformanceConfig{
		CollectionInterval: 30 * time.Second,
		RetentionPeriod:    24 * time.Hour,
		AlertThresholds: &AlertThresholds{
			CPUUsageThreshold:     80.0,
			MemoryUsageThreshold:  1024 * 1024 * 1024, // 1GB
			ResponseTimeThreshold: 1 * time.Second,
			ErrorRateThreshold:    5.0,
			GoroutineThreshold:    10000,
		},
		EnableProfiling: true,
		EnableTracing:   true,
	}

	pm := &PerformanceMonitor{
		logger:   logger,
		metrics:  &PerformanceMetrics{CustomMetrics: make(map[string]interface{})},
		config:   config,
		stopChan: make(chan struct{}),
	}

	// Start monitoring
	go pm.startMonitoring()

	return pm
}

// startMonitoring begins performance data collection
func (pm *PerformanceMonitor) startMonitoring() {
	ticker := time.NewTicker(pm.config.CollectionInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			pm.collectMetrics()
		case <-pm.stopChan:
			return
		}
	}
}

// collectMetrics gathers current performance metrics
func (pm *PerformanceMonitor) collectMetrics() {
	ctx := context.Background()

	pm.metrics.mu.Lock()
	defer pm.metrics.mu.Unlock()

	// Collect system metrics
	pm.collectSystemMetrics()

	// Update timestamp
	pm.metrics.LastUpdated = time.Now()

	// Check thresholds and alert if necessary
	pm.checkAlertThresholds(ctx)

	// Log metrics periodically
	pm.logger.Debug(ctx, "Performance metrics collected", map[string]interface{}{
		"cpu_usage":       pm.metrics.CPUUsage,
		"memory_usage":    pm.metrics.MemoryUsage,
		"goroutine_count": pm.metrics.GoroutineCount,
		"response_time":   pm.metrics.ResponseTime,
		"error_rate":      pm.metrics.ErrorRate,
		"cache_hit_rate":  pm.metrics.CacheHitRate,
	})
}

// collectSystemMetrics gathers system-level performance data
func (pm *PerformanceMonitor) collectSystemMetrics() {
	// Memory statistics
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	pm.metrics.MemoryUsage = int64(memStats.Alloc)

	// Goroutine count
	pm.metrics.GoroutineCount = runtime.NumGoroutine()

	// GC statistics
	debug.ReadGCStats(&pm.metrics.GCStats)

	// CPU usage would require additional system calls or libraries
	// For now, we'll use a placeholder
	pm.metrics.CPUUsage = pm.estimateCPUUsage()
}

// estimateCPUUsage provides a simple CPU usage estimation
func (pm *PerformanceMonitor) estimateCPUUsage() float64 {
	// This is a simplified estimation
	// In production, you'd use proper CPU monitoring
	goroutines := float64(pm.metrics.GoroutineCount)
	if goroutines > 1000 {
		return 50.0 + (goroutines-1000)/100
	}
	return goroutines / 20
}

// RecordRequest records metrics for an HTTP request
func (pm *PerformanceMonitor) RecordRequest(metrics *RequestMetrics) {
	pm.metrics.mu.Lock()
	defer pm.metrics.mu.Unlock()

	// Update request count
	pm.metrics.RequestCount++

	// Update response time (exponential moving average)
	if pm.metrics.ResponseTime == 0 {
		pm.metrics.ResponseTime = metrics.Duration
	} else {
		alpha := 0.1
		pm.metrics.ResponseTime = time.Duration(
			float64(pm.metrics.ResponseTime)*(1-alpha) + float64(metrics.Duration)*alpha,
		)
	}

	// Update error rate
	if metrics.StatusCode >= 400 {
		// Calculate error rate as exponential moving average
		if pm.metrics.ErrorRate == 0 {
			pm.metrics.ErrorRate = 1.0
		} else {
			alpha := 0.1
			pm.metrics.ErrorRate = pm.metrics.ErrorRate*(1-alpha) + alpha
		}
	} else {
		alpha := 0.1
		pm.metrics.ErrorRate = pm.metrics.ErrorRate * (1 - alpha)
	}

	// Calculate throughput (requests per second)
	pm.updateThroughput()
}

// updateThroughput calculates current throughput
func (pm *PerformanceMonitor) updateThroughput() {
	// Simple throughput calculation based on recent activity
	// In production, you'd use a more sophisticated sliding window
	elapsed := time.Since(pm.metrics.LastUpdated)
	if elapsed > 0 {
		pm.metrics.ThroughputRPS = float64(pm.metrics.RequestCount) / elapsed.Seconds()
	}
}

// RecordDatabaseMetrics records database performance metrics
func (pm *PerformanceMonitor) RecordDatabaseMetrics(connections int64, queryTime time.Duration, slowQueries int64) {
	pm.metrics.mu.Lock()
	defer pm.metrics.mu.Unlock()

	pm.metrics.DBConnections = connections

	// Update average query time
	if pm.metrics.DBQueryTime == 0 {
		pm.metrics.DBQueryTime = queryTime
	} else {
		alpha := 0.1
		pm.metrics.DBQueryTime = time.Duration(
			float64(pm.metrics.DBQueryTime)*(1-alpha) + float64(queryTime)*alpha,
		)
	}

	pm.metrics.DBSlowQueries = slowQueries
}

// RecordCacheMetrics records cache performance metrics
func (pm *PerformanceMonitor) RecordCacheMetrics(hitRate float64, size int64, evictions int64) {
	pm.metrics.mu.Lock()
	defer pm.metrics.mu.Unlock()

	pm.metrics.CacheHitRate = hitRate
	pm.metrics.CacheSize = size
	pm.metrics.CacheEvictions = evictions
}

// SetCustomMetric sets a custom performance metric
func (pm *PerformanceMonitor) SetCustomMetric(key string, value interface{}) {
	pm.metrics.mu.Lock()
	defer pm.metrics.mu.Unlock()

	pm.metrics.CustomMetrics[key] = value
}

// checkAlertThresholds checks if any metrics exceed alert thresholds
func (pm *PerformanceMonitor) checkAlertThresholds(ctx context.Context) {
	thresholds := pm.config.AlertThresholds

	// Check CPU usage
	if pm.metrics.CPUUsage > thresholds.CPUUsageThreshold {
		pm.logger.Warn(ctx, "High CPU usage detected", map[string]interface{}{
			"current_usage": pm.metrics.CPUUsage,
			"threshold":     thresholds.CPUUsageThreshold,
		})
	}

	// Check memory usage
	if pm.metrics.MemoryUsage > thresholds.MemoryUsageThreshold {
		pm.logger.Warn(ctx, "High memory usage detected", map[string]interface{}{
			"current_usage": pm.metrics.MemoryUsage,
			"threshold":     thresholds.MemoryUsageThreshold,
		})
	}

	// Check response time
	if pm.metrics.ResponseTime > thresholds.ResponseTimeThreshold {
		pm.logger.Warn(ctx, "High response time detected", map[string]interface{}{
			"current_time": pm.metrics.ResponseTime,
			"threshold":    thresholds.ResponseTimeThreshold,
		})
	}

	// Check error rate
	if pm.metrics.ErrorRate > thresholds.ErrorRateThreshold {
		pm.logger.Warn(ctx, "High error rate detected", map[string]interface{}{
			"current_rate": pm.metrics.ErrorRate,
			"threshold":    thresholds.ErrorRateThreshold,
		})
	}

	// Check goroutine count
	if pm.metrics.GoroutineCount > thresholds.GoroutineThreshold {
		pm.logger.Warn(ctx, "High goroutine count detected", map[string]interface{}{
			"current_count": pm.metrics.GoroutineCount,
			"threshold":     thresholds.GoroutineThreshold,
		})
	}
}

// GetMetrics returns current performance metrics
func (pm *PerformanceMonitor) GetMetrics() *PerformanceMetrics {
	pm.metrics.mu.RLock()
	defer pm.metrics.mu.RUnlock()

	// Create a copy without the mutex to avoid race conditions
	customMetrics := make(map[string]interface{})
	for k, v := range pm.metrics.CustomMetrics {
		customMetrics[k] = v
	}

	metrics := &PerformanceMetrics{
		CPUUsage:       pm.metrics.CPUUsage,
		MemoryUsage:    pm.metrics.MemoryUsage,
		GoroutineCount: pm.metrics.GoroutineCount,
		GCStats:        pm.metrics.GCStats,
		RequestCount:   pm.metrics.RequestCount,
		ResponseTime:   pm.metrics.ResponseTime,
		ErrorRate:      pm.metrics.ErrorRate,
		ThroughputRPS:  pm.metrics.ThroughputRPS,
		CustomMetrics:  customMetrics,
		LastUpdated:    pm.metrics.LastUpdated,
	}

	return metrics
}

// Stop stops the performance monitoring
func (pm *PerformanceMonitor) Stop() {
	close(pm.stopChan)
}

// GetHealthStatus returns overall system health status
func (pm *PerformanceMonitor) GetHealthStatus() map[string]interface{} {
	metrics := pm.GetMetrics()
	thresholds := pm.config.AlertThresholds

	status := "healthy"
	issues := []string{}

	if metrics.CPUUsage > thresholds.CPUUsageThreshold {
		status = "warning"
		issues = append(issues, "high_cpu_usage")
	}

	if metrics.MemoryUsage > thresholds.MemoryUsageThreshold {
		status = "warning"
		issues = append(issues, "high_memory_usage")
	}

	if metrics.ResponseTime > thresholds.ResponseTimeThreshold {
		status = "warning"
		issues = append(issues, "high_response_time")
	}

	if metrics.ErrorRate > thresholds.ErrorRateThreshold {
		status = "critical"
		issues = append(issues, "high_error_rate")
	}

	return map[string]interface{}{
		"status":          status,
		"issues":          issues,
		"cpu_usage":       metrics.CPUUsage,
		"memory_usage":    metrics.MemoryUsage,
		"goroutine_count": metrics.GoroutineCount,
		"response_time":   metrics.ResponseTime,
		"error_rate":      metrics.ErrorRate,
		"throughput_rps":  metrics.ThroughputRPS,
		"cache_hit_rate":  metrics.CacheHitRate,
		"last_updated":    metrics.LastUpdated,
	}
}
