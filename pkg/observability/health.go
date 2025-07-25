package observability

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// HealthStatus represents the health status of a component
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
	HealthStatusDegraded  HealthStatus = "degraded"
	HealthStatusUnknown   HealthStatus = "unknown"
)

// HealthCheck represents a health check function
type HealthCheck func(ctx context.Context) HealthCheckResult

// HealthCheckResult represents the result of a health check
type HealthCheckResult struct {
	Status      HealthStatus           `json:"status"`
	Message     string                 `json:"message,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty"`
	Duration    time.Duration          `json:"duration"`
	Timestamp   time.Time              `json:"timestamp"`
	Error       string                 `json:"error,omitempty"`
}

// HealthChecker manages health checks for the application
type HealthChecker struct {
	checks   map[string]HealthCheck
	mu       sync.RWMutex
	timeout  time.Duration
	logger   *Logger
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(logger *Logger) *HealthChecker {
	return &HealthChecker{
		checks:  make(map[string]HealthCheck),
		timeout: 30 * time.Second,
		logger:  logger,
	}
}

// RegisterCheck registers a health check
func (hc *HealthChecker) RegisterCheck(name string, check HealthCheck) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	hc.checks[name] = check
}

// UnregisterCheck removes a health check
func (hc *HealthChecker) UnregisterCheck(name string) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	delete(hc.checks, name)
}

// CheckHealth performs all registered health checks
func (hc *HealthChecker) CheckHealth(ctx context.Context) map[string]HealthCheckResult {
	hc.mu.RLock()
	checks := make(map[string]HealthCheck)
	for name, check := range hc.checks {
		checks[name] = check
	}
	hc.mu.RUnlock()

	results := make(map[string]HealthCheckResult)
	var wg sync.WaitGroup

	for name, check := range checks {
		wg.Add(1)
		go func(name string, check HealthCheck) {
			defer wg.Done()
			
			ctx, cancel := context.WithTimeout(ctx, hc.timeout)
			defer cancel()

			start := time.Now()
			result := hc.executeCheck(ctx, check)
			result.Duration = time.Since(start)
			result.Timestamp = time.Now()

			hc.mu.Lock()
			results[name] = result
			hc.mu.Unlock()
		}(name, check)
	}

	wg.Wait()
	return results
}

// executeCheck executes a single health check with error handling
func (hc *HealthChecker) executeCheck(ctx context.Context, check HealthCheck) HealthCheckResult {
	defer func() {
		if r := recover(); r != nil {
			hc.logger.Error(ctx, "Health check panicked", fmt.Errorf("panic: %v", r))
		}
	}()

	select {
	case <-ctx.Done():
		return HealthCheckResult{
			Status:  HealthStatusUnhealthy,
			Message: "Health check timed out",
			Error:   ctx.Err().Error(),
		}
	default:
		return check(ctx)
	}
}

// GetOverallStatus determines the overall health status
func (hc *HealthChecker) GetOverallStatus(results map[string]HealthCheckResult) HealthStatus {
	if len(results) == 0 {
		return HealthStatusUnknown
	}

	hasUnhealthy := false
	hasDegraded := false

	for _, result := range results {
		switch result.Status {
		case HealthStatusUnhealthy:
			hasUnhealthy = true
		case HealthStatusDegraded:
			hasDegraded = true
		}
	}

	if hasUnhealthy {
		return HealthStatusUnhealthy
	}
	if hasDegraded {
		return HealthStatusDegraded
	}
	return HealthStatusHealthy
}

// HealthResponse represents the complete health check response
type HealthResponse struct {
	Status    HealthStatus                   `json:"status"`
	Timestamp time.Time                      `json:"timestamp"`
	Duration  time.Duration                  `json:"duration"`
	Service   ServiceInfo                    `json:"service"`
	Checks    map[string]HealthCheckResult   `json:"checks"`
	System    SystemInfo                     `json:"system"`
}

// ServiceInfo contains service information
type ServiceInfo struct {
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Environment string    `json:"environment"`
	StartTime   time.Time `json:"start_time"`
	Uptime      string    `json:"uptime"`
}

// SystemInfo contains system information
type SystemInfo struct {
	GoVersion    string  `json:"go_version"`
	NumGoroutine int     `json:"num_goroutine"`
	NumCPU       int     `json:"num_cpu"`
	MemoryUsage  MemInfo `json:"memory_usage"`
}

// MemInfo contains memory information
type MemInfo struct {
	Alloc        uint64 `json:"alloc"`
	TotalAlloc   uint64 `json:"total_alloc"`
	Sys          uint64 `json:"sys"`
	NumGC        uint32 `json:"num_gc"`
	HeapAlloc    uint64 `json:"heap_alloc"`
	HeapSys      uint64 `json:"heap_sys"`
	HeapInuse    uint64 `json:"heap_inuse"`
	HeapReleased uint64 `json:"heap_released"`
}

// HealthServer provides HTTP endpoints for health checks
type HealthServer struct {
	checker   *HealthChecker
	service   ServiceInfo
	startTime time.Time
	logger    *Logger
}

// NewHealthServer creates a new health server
func NewHealthServer(checker *HealthChecker, service ServiceInfo, logger *Logger) *HealthServer {
	return &HealthServer{
		checker:   checker,
		service:   service,
		startTime: time.Now(),
		logger:    logger,
	}
}

// RegisterRoutes registers health check routes
func (hs *HealthServer) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/health", hs.HealthHandler).Methods("GET")
	router.HandleFunc("/health/live", hs.LivenessHandler).Methods("GET")
	router.HandleFunc("/health/ready", hs.ReadinessHandler).Methods("GET")
	router.HandleFunc("/health/startup", hs.StartupHandler).Methods("GET")
}

// HealthHandler handles comprehensive health checks
func (hs *HealthServer) HealthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	start := time.Now()

	results := hs.checker.CheckHealth(ctx)
	overallStatus := hs.checker.GetOverallStatus(results)

	response := HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Duration:  time.Since(start),
		Service:   hs.getServiceInfo(),
		Checks:    results,
		System:    hs.getSystemInfo(),
	}

	// Set HTTP status code based on health status
	statusCode := http.StatusOK
	switch overallStatus {
	case HealthStatusUnhealthy:
		statusCode = http.StatusServiceUnavailable
	case HealthStatusDegraded:
		statusCode = http.StatusPartialContent
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)

	// Log health check
	hs.logger.Info(ctx, "Health check performed", map[string]interface{}{
		"status":   overallStatus,
		"duration": time.Since(start).Milliseconds(),
		"checks":   len(results),
	})
}

// LivenessHandler handles liveness probes (Kubernetes)
func (hs *HealthServer) LivenessHandler(w http.ResponseWriter, r *http.Request) {
	// Simple liveness check - service is running
	response := map[string]interface{}{
		"status":    "alive",
		"timestamp": time.Now(),
		"service":   hs.service.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ReadinessHandler handles readiness probes (Kubernetes)
func (hs *HealthServer) ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	results := hs.checker.CheckHealth(ctx)
	overallStatus := hs.checker.GetOverallStatus(results)

	response := map[string]interface{}{
		"status":    overallStatus,
		"timestamp": time.Now(),
		"service":   hs.service.Name,
		"ready":     overallStatus == HealthStatusHealthy,
	}

	statusCode := http.StatusOK
	if overallStatus != HealthStatusHealthy {
		statusCode = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// StartupHandler handles startup probes (Kubernetes)
func (hs *HealthServer) StartupHandler(w http.ResponseWriter, r *http.Request) {
	// Check if service has been running for at least 10 seconds
	uptime := time.Since(hs.startTime)
	started := uptime > 10*time.Second

	response := map[string]interface{}{
		"status":    map[bool]string{true: "started", false: "starting"}[started],
		"timestamp": time.Now(),
		"service":   hs.service.Name,
		"uptime":    uptime.String(),
		"started":   started,
	}

	statusCode := http.StatusOK
	if !started {
		statusCode = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// getServiceInfo returns current service information
func (hs *HealthServer) getServiceInfo() ServiceInfo {
	uptime := time.Since(hs.startTime)
	return ServiceInfo{
		Name:        hs.service.Name,
		Version:     hs.service.Version,
		Environment: hs.service.Environment,
		StartTime:   hs.startTime,
		Uptime:      uptime.String(),
	}
}

// getSystemInfo returns current system information
func (hs *HealthServer) getSystemInfo() SystemInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return SystemInfo{
		GoVersion:    runtime.Version(),
		NumGoroutine: runtime.NumGoroutine(),
		NumCPU:       runtime.NumCPU(),
		MemoryUsage: MemInfo{
			Alloc:        m.Alloc,
			TotalAlloc:   m.TotalAlloc,
			Sys:          m.Sys,
			NumGC:        m.NumGC,
			HeapAlloc:    m.HeapAlloc,
			HeapSys:      m.HeapSys,
			HeapInuse:    m.HeapInuse,
			HeapReleased: m.HeapReleased,
		},
	}
}

// Common Health Checks

// DatabaseHealthCheck creates a health check for database connectivity
func DatabaseHealthCheck(pingFunc func(ctx context.Context) error) HealthCheck {
	return func(ctx context.Context) HealthCheckResult {
		if err := pingFunc(ctx); err != nil {
			return HealthCheckResult{
				Status:  HealthStatusUnhealthy,
				Message: "Database connection failed",
				Error:   err.Error(),
			}
		}
		return HealthCheckResult{
			Status:  HealthStatusHealthy,
			Message: "Database connection successful",
		}
	}
}

// RedisHealthCheck creates a health check for Redis connectivity
func RedisHealthCheck(pingFunc func(ctx context.Context) error) HealthCheck {
	return func(ctx context.Context) HealthCheckResult {
		if err := pingFunc(ctx); err != nil {
			return HealthCheckResult{
				Status:  HealthStatusUnhealthy,
				Message: "Redis connection failed",
				Error:   err.Error(),
			}
		}
		return HealthCheckResult{
			Status:  HealthStatusHealthy,
			Message: "Redis connection successful",
		}
	}
}

// HTTPServiceHealthCheck creates a health check for HTTP service dependencies
func HTTPServiceHealthCheck(url string, timeout time.Duration) HealthCheck {
	return func(ctx context.Context) HealthCheckResult {
		client := &http.Client{Timeout: timeout}
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return HealthCheckResult{
				Status:  HealthStatusUnhealthy,
				Message: "Failed to create request",
				Error:   err.Error(),
			}
		}

		resp, err := client.Do(req)
		if err != nil {
			return HealthCheckResult{
				Status:  HealthStatusUnhealthy,
				Message: "HTTP request failed",
				Error:   err.Error(),
			}
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return HealthCheckResult{
				Status:  HealthStatusHealthy,
				Message: "HTTP service is healthy",
				Details: map[string]interface{}{
					"status_code": resp.StatusCode,
					"url":         url,
				},
			}
		}

		return HealthCheckResult{
			Status:  HealthStatusUnhealthy,
			Message: "HTTP service returned error status",
			Details: map[string]interface{}{
				"status_code": resp.StatusCode,
				"url":         url,
			},
		}
	}
}
