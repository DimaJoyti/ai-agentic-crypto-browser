package ai

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ai-agentic-browser/internal/config"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create test logger
func createTestLogger() *observability.Logger {
	return observability.NewLogger(config.ObservabilityConfig{
		ServiceName: "test",
		LogLevel:    "info",
		LogFormat:   "json",
	})
}

// Mock health checker for testing
type mockHealthChecker struct {
	healthy bool
	models  []string
	err     error
}

func (m *mockHealthChecker) IsHealthy(ctx context.Context) error {
	if m.err != nil {
		return m.err
	}
	if !m.healthy {
		return errors.New("provider is unhealthy")
	}
	return nil
}

func (m *mockHealthChecker) ListModels(ctx context.Context) ([]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.models, nil
}

func TestHealthMonitor_RegisterProvider(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{
		ServiceName: "test",
		LogLevel:    "info",
		LogFormat:   "json",
	})

	monitor := NewHealthMonitor(logger, 1*time.Second)

	mockProvider := &mockHealthChecker{
		healthy: true,
		models:  []string{"model1", "model2"},
	}

	monitor.RegisterProvider("test-provider", mockProvider)

	// Check that provider was registered
	status, exists := monitor.GetStatus("test-provider")
	assert.True(t, exists)
	assert.Equal(t, "test-provider", status.Provider)
	assert.False(t, status.Healthy) // Should be false initially until first check
}

func TestHealthMonitor_CheckProviderNow(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{
		ServiceName: "test",
		LogLevel:    "info",
		LogFormat:   "json",
	})

	monitor := NewHealthMonitor(logger, 1*time.Second)

	mockProvider := &mockHealthChecker{
		healthy: true,
		models:  []string{"model1", "model2"},
	}

	monitor.RegisterProvider("test-provider", mockProvider)

	ctx := context.Background()
	err := monitor.CheckProviderNow(ctx, "test-provider")
	require.NoError(t, err)

	// Check that status was updated
	status, exists := monitor.GetStatus("test-provider")
	assert.True(t, exists)
	assert.True(t, status.Healthy)
	assert.Len(t, status.Models, 2)
	assert.Contains(t, status.Models, "model1")
	assert.Contains(t, status.Models, "model2")
	assert.Empty(t, status.Error)
}

func TestHealthMonitor_CheckUnhealthyProvider(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{
		ServiceName: "test",
		LogLevel:    "info",
		LogFormat:   "json",
	})

	monitor := NewHealthMonitor(logger, 1*time.Second)

	mockProvider := &mockHealthChecker{
		healthy: false,
		err:     errors.New("connection failed"),
	}

	monitor.RegisterProvider("unhealthy-provider", mockProvider)

	ctx := context.Background()
	err := monitor.CheckProviderNow(ctx, "unhealthy-provider")
	require.NoError(t, err) // CheckProviderNow doesn't return provider errors

	// Check that status reflects unhealthy state
	status, exists := monitor.GetStatus("unhealthy-provider")
	assert.True(t, exists)
	assert.False(t, status.Healthy)
	assert.Equal(t, "connection failed", status.Error)
	assert.Empty(t, status.Models)
}

func TestHealthMonitor_GetAllStatuses(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{
		ServiceName: "test",
		LogLevel:    "info",
		LogFormat:   "json",
	})

	monitor := NewHealthMonitor(logger, 1*time.Second)

	// Register multiple providers
	healthyProvider := &mockHealthChecker{
		healthy: true,
		models:  []string{"model1"},
	}

	unhealthyProvider := &mockHealthChecker{
		healthy: false,
		err:     errors.New("failed"),
	}

	monitor.RegisterProvider("healthy", healthyProvider)
	monitor.RegisterProvider("unhealthy", unhealthyProvider)

	ctx := context.Background()
	monitor.CheckProviderNow(ctx, "healthy")
	monitor.CheckProviderNow(ctx, "unhealthy")

	statuses := monitor.GetAllStatuses()
	assert.Len(t, statuses, 2)

	healthyStatus := statuses["healthy"]
	assert.True(t, healthyStatus.Healthy)
	assert.Empty(t, healthyStatus.Error)

	unhealthyStatus := statuses["unhealthy"]
	assert.False(t, unhealthyStatus.Healthy)
	assert.NotEmpty(t, unhealthyStatus.Error)
}

func TestHealthMonitor_IsProviderHealthy(t *testing.T) {
	logger := createTestLogger()

	monitor := NewHealthMonitor(logger, 1*time.Second)

	healthyProvider := &mockHealthChecker{
		healthy: true,
		models:  []string{"model1"},
	}

	monitor.RegisterProvider("test-provider", healthyProvider)

	// Initially should be false
	assert.False(t, monitor.IsProviderHealthy("test-provider"))

	// After health check should be true
	ctx := context.Background()
	monitor.CheckProviderNow(ctx, "test-provider")
	assert.True(t, monitor.IsProviderHealthy("test-provider"))

	// Non-existent provider should be false
	assert.False(t, monitor.IsProviderHealthy("non-existent"))
}

func TestHealthMonitor_GetHealthCheckResponse(t *testing.T) {
	logger := createTestLogger()

	monitor := NewHealthMonitor(logger, 1*time.Second)

	// Register providers with different health states
	healthyProvider1 := &mockHealthChecker{healthy: true, models: []string{"model1"}}
	healthyProvider2 := &mockHealthChecker{healthy: true, models: []string{"model2"}}
	unhealthyProvider := &mockHealthChecker{healthy: false, err: errors.New("failed")}

	monitor.RegisterProvider("healthy1", healthyProvider1)
	monitor.RegisterProvider("healthy2", healthyProvider2)
	monitor.RegisterProvider("unhealthy", unhealthyProvider)

	ctx := context.Background()
	monitor.CheckProviderNow(ctx, "healthy1")
	monitor.CheckProviderNow(ctx, "healthy2")
	monitor.CheckProviderNow(ctx, "unhealthy")

	response := monitor.GetHealthCheckResponse()

	assert.Equal(t, "degraded", response.Status) // Some healthy, some unhealthy
	assert.Equal(t, 3, response.Summary.Total)
	assert.Equal(t, 2, response.Summary.Healthy)
	assert.Equal(t, 1, response.Summary.Unhealthy)
	assert.Len(t, response.Providers, 3)
}

func TestHealthMonitor_GetHealthCheckResponse_AllHealthy(t *testing.T) {
	logger := createTestLogger()

	monitor := NewHealthMonitor(logger, 1*time.Second)

	healthyProvider := &mockHealthChecker{healthy: true, models: []string{"model1"}}
	monitor.RegisterProvider("healthy", healthyProvider)

	ctx := context.Background()
	monitor.CheckProviderNow(ctx, "healthy")

	response := monitor.GetHealthCheckResponse()

	assert.Equal(t, "healthy", response.Status)
	assert.Equal(t, 1, response.Summary.Total)
	assert.Equal(t, 1, response.Summary.Healthy)
	assert.Equal(t, 0, response.Summary.Unhealthy)
}

func TestHealthMonitor_GetHealthCheckResponse_AllUnhealthy(t *testing.T) {
	logger := createTestLogger()

	monitor := NewHealthMonitor(logger, 1*time.Second)

	unhealthyProvider := &mockHealthChecker{healthy: false, err: errors.New("failed")}
	monitor.RegisterProvider("unhealthy", unhealthyProvider)

	ctx := context.Background()
	monitor.CheckProviderNow(ctx, "unhealthy")

	response := monitor.GetHealthCheckResponse()

	assert.Equal(t, "unhealthy", response.Status)
	assert.Equal(t, 1, response.Summary.Total)
	assert.Equal(t, 0, response.Summary.Healthy)
	assert.Equal(t, 1, response.Summary.Unhealthy)
}

func TestHealthMonitor_GetProviderModels(t *testing.T) {
	logger := createTestLogger()

	monitor := NewHealthMonitor(logger, 1*time.Second)

	healthyProvider := &mockHealthChecker{
		healthy: true,
		models:  []string{"model1", "model2", "model3"},
	}

	monitor.RegisterProvider("test-provider", healthyProvider)

	ctx := context.Background()
	monitor.CheckProviderNow(ctx, "test-provider")

	models, err := monitor.GetProviderModels("test-provider")
	require.NoError(t, err)
	assert.Len(t, models, 3)
	assert.Contains(t, models, "model1")
	assert.Contains(t, models, "model2")
	assert.Contains(t, models, "model3")
}

func TestHealthMonitor_GetProviderModels_UnhealthyProvider(t *testing.T) {
	logger := createTestLogger()

	monitor := NewHealthMonitor(logger, 1*time.Second)

	unhealthyProvider := &mockHealthChecker{
		healthy: false,
		err:     errors.New("connection failed"),
	}

	monitor.RegisterProvider("unhealthy-provider", unhealthyProvider)

	ctx := context.Background()
	monitor.CheckProviderNow(ctx, "unhealthy-provider")

	models, err := monitor.GetProviderModels("unhealthy-provider")
	assert.Error(t, err)
	assert.Nil(t, models)
	assert.Contains(t, err.Error(), "not healthy")
}

func TestHealthMonitor_GetProviderModels_NonExistentProvider(t *testing.T) {
	logger := createTestLogger()

	monitor := NewHealthMonitor(logger, 1*time.Second)

	models, err := monitor.GetProviderModels("non-existent")
	assert.Error(t, err)
	assert.Nil(t, models)
	assert.Contains(t, err.Error(), "not found")
}

func TestHealthMonitor_ConcurrentAccess(t *testing.T) {
	logger := createTestLogger()

	monitor := NewHealthMonitor(logger, 1*time.Second)

	healthyProvider := &mockHealthChecker{
		healthy: true,
		models:  []string{"model1"},
	}

	monitor.RegisterProvider("test-provider", healthyProvider)

	ctx := context.Background()

	// Test concurrent access to health monitor
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()

			// Perform various operations concurrently
			monitor.CheckProviderNow(ctx, "test-provider")
			monitor.GetStatus("test-provider")
			monitor.IsProviderHealthy("test-provider")
			monitor.GetAllStatuses()
			monitor.GetHealthCheckResponse()
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify final state is consistent
	status, exists := monitor.GetStatus("test-provider")
	assert.True(t, exists)
	assert.True(t, status.Healthy)
}
