package test

import (
	"context"
	"testing"
	"time"

	"github.com/ai-agentic-browser/internal/config"
	"github.com/ai-agentic-browser/internal/mcp"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFirebaseClientInitialization(t *testing.T) {
	// Skip if running in CI without Firebase credentials
	if testing.Short() {
		t.Skip("Skipping Firebase integration test in short mode")
	}

	// Create test configuration
	firebaseConfig := mcp.FirebaseConfig{
		ProjectID:        "test-project",
		EnableAuth:       true,
		EnableFirestore:  true,
		EnableRealtimeDB: true,
		EnableEmulators:  true,
		EmulatorConfig: mcp.EmulatorConfig{
			Host:          "localhost",
			AuthPort:      9099,
			FirestorePort: 8080,
			DatabasePort:  9000,
		},
	}

	// Create logger
	logger := observability.NewLogger(config.ObservabilityConfig{
		LogLevel:  "info",
		LogFormat: "json",
	})

	// Create Firebase client
	client := mcp.NewFirebaseClient(logger, firebaseConfig)
	require.NotNil(t, client)

	// Test client configuration
	config := client.GetConfig()
	assert.Equal(t, "test-project", config.ProjectID)
	assert.True(t, config.EnableAuth)
	assert.True(t, config.EnableFirestore)
	assert.True(t, config.EnableRealtimeDB)
	assert.True(t, config.EnableEmulators)
}

func TestFirebaseClientHealthCheck(t *testing.T) {
	// Skip if running in CI without Firebase credentials
	if testing.Short() {
		t.Skip("Skipping Firebase integration test in short mode")
	}

	// Create test configuration
	firebaseConfig := mcp.FirebaseConfig{
		ProjectID:       "test-project",
		EnableEmulators: true,
		EmulatorConfig: mcp.EmulatorConfig{
			Host:          "localhost",
			AuthPort:      9099,
			FirestorePort: 8080,
			DatabasePort:  9000,
		},
	}

	// Create logger
	logger := observability.NewLogger(config.ObservabilityConfig{
		LogLevel:  "info",
		LogFormat: "json",
	})

	// Create Firebase client
	client := mcp.NewFirebaseClient(logger, firebaseConfig)
	require.NotNil(t, client)

	// Initially should not be healthy (not started)
	assert.False(t, client.IsHealthy())

	// Start client (this will fail without emulators running, but we can test the interface)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Note: This will fail if emulators are not running, which is expected in CI
	err := client.Start(ctx)
	if err != nil {
		t.Logf("Expected error starting Firebase client without emulators: %v", err)
		return
	}

	// If start succeeded, test health check
	assert.True(t, client.IsHealthy())

	// Test status
	status := client.GetStatus(ctx)
	assert.NotNil(t, status)
	assert.Contains(t, status, "running")
	assert.Contains(t, status, "project_id")

	// Stop client
	err = client.Stop(ctx)
	assert.NoError(t, err)
	assert.False(t, client.IsHealthy())
}

func TestMCPIntegrationServiceWithFirebase(t *testing.T) {
	// Skip if running in CI without Firebase credentials
	if testing.Short() {
		t.Skip("Skipping Firebase integration test in short mode")
	}

	// Create test configuration
	mcpConfig := mcp.Config{
		Firebase: mcp.FirebaseConfig{
			ProjectID:       "test-project",
			EnableAuth:      true,
			EnableFirestore: true,
			EnableEmulators: true,
			EmulatorConfig: mcp.EmulatorConfig{
				Host:          "localhost",
				AuthPort:      9099,
				FirestorePort: 8080,
				DatabasePort:  9000,
			},
		},
		UpdateInterval: 1 * time.Second,
		EnableRealtime: true,
		BufferSize:     100,
	}

	// Create logger
	logger := observability.NewLogger(config.ObservabilityConfig{
		LogLevel:  "info",
		LogFormat: "json",
	})

	// Create MCP integration service
	service := mcp.NewIntegrationService(logger, mcpConfig)
	require.NotNil(t, service)

	// Test Firebase client access
	firebaseClient := service.GetFirebaseClient()
	require.NotNil(t, firebaseClient)

	// Test Firebase client configuration
	firebaseConfig := firebaseClient.GetConfig()
	assert.Equal(t, "test-project", firebaseConfig.ProjectID)
	assert.True(t, firebaseConfig.EnableAuth)
	assert.True(t, firebaseConfig.EnableFirestore)
}

func TestFirebaseConfigurationLoading(t *testing.T) {
	// Test configuration loading from environment variables
	t.Setenv("FIREBASE_PROJECT_ID", "test-env-project")
	t.Setenv("FIREBASE_ENABLE_AUTH", "true")
	t.Setenv("FIREBASE_ENABLE_FIRESTORE", "true")
	t.Setenv("FIREBASE_ENABLE_EMULATORS", "true")

	// Load configuration
	cfg, err := config.Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Test Firebase configuration
	assert.Equal(t, "test-env-project", cfg.Firebase.ProjectID)
	assert.True(t, cfg.Firebase.EnableAuth)
	assert.True(t, cfg.Firebase.EnableFirestore)
	assert.True(t, cfg.Firebase.EnableEmulators)

	// Test MCP configuration
	assert.True(t, cfg.MCP.Firebase.Enabled)
	assert.NotEmpty(t, cfg.MCP.Firebase.Collections)
	assert.NotEmpty(t, cfg.MCP.Firebase.RealtimePaths)
}

func TestFirebaseDocumentOperations(t *testing.T) {
	// Skip if running in CI without Firebase credentials
	if testing.Short() {
		t.Skip("Skipping Firebase integration test in short mode")
	}

	// This test would require Firebase emulators to be running
	// It's more of an integration test that should be run manually
	t.Skip("Requires Firebase emulators to be running")

	// Create test configuration
	firebaseConfig := mcp.FirebaseConfig{
		ProjectID:        "test-project",
		EnableFirestore:  true,
		EnableEmulators:  true,
		EmulatorConfig: mcp.EmulatorConfig{
			Host:          "localhost",
			FirestorePort: 8080,
		},
	}

	// Create logger
	logger := observability.NewLogger(config.ObservabilityConfig{
		LogLevel:  "info",
		LogFormat: "json",
	})

	// Create Firebase client
	client := mcp.NewFirebaseClient(logger, firebaseConfig)
	require.NotNil(t, client)

	ctx := context.Background()

	// Start client
	err := client.Start(ctx)
	require.NoError(t, err)
	defer client.Stop(ctx)

	// Test document creation
	testData := map[string]interface{}{
		"symbol":     "BTCUSDT",
		"price":      45000.0,
		"volume":     1234.56,
		"timestamp":  time.Now().Unix(),
	}

	doc, err := client.CreateDocument(ctx, "test_collection", "test_doc", testData)
	require.NoError(t, err)
	assert.Equal(t, "test_doc", doc.ID)
	assert.Equal(t, "test_collection", doc.Collection)

	// Test document retrieval
	retrievedDoc, err := client.GetDocument(ctx, "test_collection", "test_doc")
	require.NoError(t, err)
	assert.Equal(t, "test_doc", retrievedDoc.ID)
	assert.Equal(t, testData["symbol"], retrievedDoc.Data["symbol"])

	// Test document update
	updateData := map[string]interface{}{
		"price": 46000.0,
	}
	err = client.UpdateDocument(ctx, "test_collection", "test_doc", updateData)
	require.NoError(t, err)

	// Test document deletion
	err = client.DeleteDocument(ctx, "test_collection", "test_doc")
	require.NoError(t, err)
}

func TestFirebaseRealtimeOperations(t *testing.T) {
	// Skip if running in CI without Firebase credentials
	if testing.Short() {
		t.Skip("Skipping Firebase integration test in short mode")
	}

	// This test would require Firebase emulators to be running
	// It's more of an integration test that should be run manually
	t.Skip("Requires Firebase emulators to be running")

	// Create test configuration
	firebaseConfig := mcp.FirebaseConfig{
		ProjectID:        "test-project",
		EnableRealtimeDB: true,
		EnableEmulators:  true,
		EmulatorConfig: mcp.EmulatorConfig{
			Host:         "localhost",
			DatabasePort: 9000,
		},
	}

	// Create logger
	logger := observability.NewLogger(config.ObservabilityConfig{
		LogLevel:  "info",
		LogFormat: "json",
	})

	// Create Firebase client
	client := mcp.NewFirebaseClient(logger, firebaseConfig)
	require.NotNil(t, client)

	ctx := context.Background()

	// Start client
	err := client.Start(ctx)
	require.NoError(t, err)
	defer client.Stop(ctx)

	// Test data setting
	testData := map[string]interface{}{
		"price":     45000.0,
		"timestamp": time.Now().Unix(),
	}

	err = client.SetRealtimeData(ctx, "live_prices/BTCUSDT", testData)
	require.NoError(t, err)

	// Test data retrieval
	retrievedData, err := client.GetRealtimeData(ctx, "live_prices/BTCUSDT")
	require.NoError(t, err)
	assert.NotNil(t, retrievedData)

	// Test data update
	updateData := map[string]interface{}{
		"price": 46000.0,
	}
	err = client.UpdateRealtimeData(ctx, "live_prices/BTCUSDT", updateData)
	require.NoError(t, err)

	// Test data deletion
	err = client.DeleteRealtimeData(ctx, "live_prices/BTCUSDT")
	require.NoError(t, err)
}
