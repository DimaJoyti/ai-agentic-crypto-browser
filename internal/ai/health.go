package ai

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
)

// HealthChecker interface for AI provider health checks
type HealthChecker interface {
	IsHealthy(ctx context.Context) error
	ListModels(ctx context.Context) ([]string, error)
}

// HealthStatus represents the health status of an AI provider
type HealthStatus struct {
	Provider     string        `json:"provider"`
	Healthy      bool          `json:"healthy"`
	LastChecked  time.Time     `json:"last_checked"`
	Error        string        `json:"error,omitempty"`
	Models       []string      `json:"models,omitempty"`
	ResponseTime time.Duration `json:"response_time"`
}

// HealthMonitor monitors the health of AI providers
type HealthMonitor struct {
	providers map[string]HealthChecker
	statuses  map[string]*HealthStatus
	mutex     sync.RWMutex
	logger    *observability.Logger
	stopCh    chan struct{}
	interval  time.Duration
}

// NewHealthMonitor creates a new health monitor
func NewHealthMonitor(logger *observability.Logger, interval time.Duration) *HealthMonitor {
	if interval == 0 {
		interval = 30 * time.Second
	}

	return &HealthMonitor{
		providers: make(map[string]HealthChecker),
		statuses:  make(map[string]*HealthStatus),
		logger:    logger,
		stopCh:    make(chan struct{}),
		interval:  interval,
	}
}

// RegisterProvider registers a provider for health monitoring
func (hm *HealthMonitor) RegisterProvider(name string, provider HealthChecker) {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()

	hm.providers[name] = provider
	hm.statuses[name] = &HealthStatus{
		Provider:    name,
		Healthy:     false,
		LastChecked: time.Time{},
	}
}

// Start begins health monitoring
func (hm *HealthMonitor) Start(ctx context.Context) {
	ticker := time.NewTicker(hm.interval)
	defer ticker.Stop()

	// Initial health check
	hm.checkAllProviders(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-hm.stopCh:
			return
		case <-ticker.C:
			hm.checkAllProviders(ctx)
		}
	}
}

// Stop stops health monitoring
func (hm *HealthMonitor) Stop() {
	close(hm.stopCh)
}

// GetStatus returns the health status of a specific provider
func (hm *HealthMonitor) GetStatus(provider string) (*HealthStatus, bool) {
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()

	status, exists := hm.statuses[provider]
	if !exists {
		return nil, false
	}

	// Return a copy to avoid race conditions
	statusCopy := *status
	return &statusCopy, true
}

// GetAllStatuses returns the health status of all providers
func (hm *HealthMonitor) GetAllStatuses() map[string]*HealthStatus {
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()

	result := make(map[string]*HealthStatus)
	for name, status := range hm.statuses {
		statusCopy := *status
		result[name] = &statusCopy
	}

	return result
}

// IsProviderHealthy checks if a specific provider is healthy
func (hm *HealthMonitor) IsProviderHealthy(provider string) bool {
	status, exists := hm.GetStatus(provider)
	return exists && status.Healthy
}

// checkAllProviders performs health checks on all registered providers
func (hm *HealthMonitor) checkAllProviders(ctx context.Context) {
	hm.mutex.RLock()
	providers := make(map[string]HealthChecker)
	for name, provider := range hm.providers {
		providers[name] = provider
	}
	hm.mutex.RUnlock()

	var wg sync.WaitGroup
	for name, provider := range providers {
		wg.Add(1)
		go func(providerName string, p HealthChecker) {
			defer wg.Done()
			hm.checkProvider(ctx, providerName, p)
		}(name, provider)
	}

	wg.Wait()
}

// checkProvider performs a health check on a single provider
func (hm *HealthMonitor) checkProvider(ctx context.Context, name string, provider HealthChecker) {
	start := time.Now()

	// Create a timeout context for the health check
	checkCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var err error
	var models []string

	// Check if provider is healthy
	err = provider.IsHealthy(checkCtx)
	responseTime := time.Since(start)

	// If healthy, try to get models
	if err == nil {
		models, _ = provider.ListModels(checkCtx)
	}

	// Update status
	hm.mutex.Lock()
	status := hm.statuses[name]
	status.Healthy = err == nil
	status.LastChecked = time.Now()
	status.ResponseTime = responseTime
	status.Models = models
	if err != nil {
		status.Error = err.Error()
	} else {
		status.Error = ""
	}
	hm.mutex.Unlock()

	// Log health check result
	logFields := map[string]interface{}{
		"provider":      name,
		"healthy":       err == nil,
		"response_time": responseTime,
		"models_count":  len(models),
	}

	if err != nil {
		logFields["error"] = err.Error()
		hm.logger.Warn(ctx, "Provider health check failed", logFields)
	} else {
		hm.logger.Debug(ctx, "Provider health check successful", logFields)
	}
}

// HealthCheckResponse represents the response for health check endpoints
type HealthCheckResponse struct {
	Status    string                   `json:"status"`
	Timestamp time.Time                `json:"timestamp"`
	Providers map[string]*HealthStatus `json:"providers"`
	Summary   HealthSummary            `json:"summary"`
}

// HealthSummary provides a summary of provider health
type HealthSummary struct {
	Total     int `json:"total"`
	Healthy   int `json:"healthy"`
	Unhealthy int `json:"unhealthy"`
}

// GetHealthCheckResponse returns a comprehensive health check response
func (hm *HealthMonitor) GetHealthCheckResponse() *HealthCheckResponse {
	statuses := hm.GetAllStatuses()

	summary := HealthSummary{
		Total: len(statuses),
	}

	for _, status := range statuses {
		if status.Healthy {
			summary.Healthy++
		} else {
			summary.Unhealthy++
		}
	}

	overallStatus := "healthy"
	if summary.Unhealthy > 0 {
		if summary.Healthy == 0 {
			overallStatus = "unhealthy"
		} else {
			overallStatus = "degraded"
		}
	}

	return &HealthCheckResponse{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Providers: statuses,
		Summary:   summary,
	}
}

// CheckProviderNow performs an immediate health check on a specific provider
func (hm *HealthMonitor) CheckProviderNow(ctx context.Context, providerName string) error {
	hm.mutex.RLock()
	provider, exists := hm.providers[providerName]
	hm.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("provider %s not found", providerName)
	}

	hm.checkProvider(ctx, providerName, provider)
	return nil
}

// GetProviderModels returns the available models for a specific provider
func (hm *HealthMonitor) GetProviderModels(providerName string) ([]string, error) {
	status, exists := hm.GetStatus(providerName)
	if !exists {
		return nil, fmt.Errorf("provider %s not found", providerName)
	}

	if !status.Healthy {
		return nil, fmt.Errorf("provider %s is not healthy: %s", providerName, status.Error)
	}

	return status.Models, nil
}
