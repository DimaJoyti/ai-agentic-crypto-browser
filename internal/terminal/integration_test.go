package terminal

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ai-agentic-browser/internal/config"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTerminalServiceIntegration(t *testing.T) {
	// Create test service
	service := createTestTerminalService(t)
	
	// Start service
	ctx := context.Background()
	err := service.Start(ctx)
	require.NoError(t, err)
	defer service.Shutdown(ctx)
	
	// Test service is running
	assert.True(t, service.isRunning)
	
	// Test command registry has commands
	commands := service.commandRegistry.ListCommands()
	assert.Greater(t, len(commands), 0)
	
	// Test essential commands are registered
	essentialCommands := []string{"status", "help", "buy", "sell", "analyze"}
	commandMap := make(map[string]bool)
	for _, cmd := range commands {
		commandMap[cmd.Name] = true
	}
	
	for _, essential := range essentialCommands {
		assert.True(t, commandMap[essential], "Essential command %s not found", essential)
	}
}

func TestSessionWorkflow(t *testing.T) {
	service := createTestTerminalService(t)
	ctx := context.Background()
	
	// Create session
	session, err := service.sessionManager.CreateSession(ctx, "test-user", nil)
	require.NoError(t, err)
	assert.NotEmpty(t, session.ID)
	assert.Equal(t, "test-user", session.UserID)
	
	// Execute commands
	testCommands := []struct {
		command     string
		expectError bool
	}{
		{"status", false},
		{"help", false},
		{"price BTC", false},
		{"analyze BTC", false},
		{"invalid-command", true},
		{"buy", true}, // Should fail due to missing arguments
	}
	
	for _, tc := range testCommands {
		result, err := service.commandRegistry.ExecuteCommand(ctx, tc.command, session)
		require.NoError(t, err)
		
		if tc.expectError {
			assert.NotEqual(t, 0, result.ExitCode, "Command %s should have failed", tc.command)
		} else {
			assert.Equal(t, 0, result.ExitCode, "Command %s should have succeeded", tc.command)
			assert.NotEmpty(t, result.Output, "Command %s should have output", tc.command)
		}
	}
	
	// Check command history
	history, err := service.sessionManager.GetSessionHistory(ctx, session.ID)
	require.NoError(t, err)
	assert.Equal(t, len(testCommands), len(history))
	
	// Delete session
	err = service.sessionManager.DeleteSession(ctx, session.ID)
	require.NoError(t, err)
	
	// Verify session is deleted
	_, err = service.sessionManager.GetSession(ctx, session.ID)
	assert.Error(t, err)
}

func TestHTTPEndpoints(t *testing.T) {
	service := createTestTerminalService(t)
	
	// Create HTTP router
	router := mux.NewRouter()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})
	router.HandleFunc("/api/v1/sessions", service.HandleCreateSession).Methods("POST")
	router.HandleFunc("/api/v1/sessions", service.HandleListSessions).Methods("GET")
	router.HandleFunc("/api/v1/sessions/{sessionId}", service.HandleGetSession).Methods("GET")
	router.HandleFunc("/api/v1/sessions/{sessionId}", service.HandleDeleteSession).Methods("DELETE")
	router.HandleFunc("/api/v1/sessions/{sessionId}/history", service.HandleGetHistory).Methods("GET")
	router.HandleFunc("/api/v1/commands", service.HandleListCommands).Methods("GET")
	router.HandleFunc("/api/v1/commands/{command}/help", service.HandleGetCommandHelp).Methods("GET")
	
	server := httptest.NewServer(router)
	defer server.Close()
	
	// Test health endpoint
	resp, err := http.Get(server.URL + "/health")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	// Test create session
	createSessionBody := `{"user_id": "test-user", "environment": {}}`
	resp, err = http.Post(server.URL+"/api/v1/sessions", "application/json", strings.NewReader(createSessionBody))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	var createResponse CreateSessionResponse
	err = json.NewDecoder(resp.Body).Decode(&createResponse)
	require.NoError(t, err)
	assert.NotEmpty(t, createResponse.SessionID)
	
	sessionID := createResponse.SessionID
	
	// Test get session
	resp, err = http.Get(server.URL + "/api/v1/sessions/" + sessionID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	// Test list sessions
	resp, err = http.Get(server.URL + "/api/v1/sessions?user_id=test-user")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	// Test list commands
	resp, err = http.Get(server.URL + "/api/v1/commands")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	var commands []CommandInfo
	err = json.NewDecoder(resp.Body).Decode(&commands)
	require.NoError(t, err)
	assert.Greater(t, len(commands), 0)
	
	// Test command help
	resp, err = http.Get(server.URL + "/api/v1/commands/status/help")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	// Test delete session
	req, err := http.NewRequest("DELETE", server.URL+"/api/v1/sessions/"+sessionID, nil)
	require.NoError(t, err)
	
	client := &http.Client{}
	resp, err = client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestWebSocketIntegration(t *testing.T) {
	service := createTestTerminalService(t)
	ctx := context.Background()
	
	// Start service
	err := service.Start(ctx)
	require.NoError(t, err)
	defer service.Shutdown(ctx)
	
	// Create HTTP server for WebSocket
	router := mux.NewRouter()
	router.HandleFunc("/ws", service.HandleWebSocket)
	
	server := httptest.NewServer(router)
	defer server.Close()
	
	// Convert HTTP URL to WebSocket URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	
	// Connect to WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()
	
	// Test session creation
	createMsg := WSMessage{
		Type:      "session_create",
		Data:      map[string]string{"user_id": "test-user"},
		Timestamp: time.Now(),
	}
	
	err = conn.WriteJSON(createMsg)
	require.NoError(t, err)
	
	// Read welcome message
	var welcomeMsg WSMessage
	err = conn.ReadJSON(&welcomeMsg)
	require.NoError(t, err)
	assert.Equal(t, "welcome", welcomeMsg.Type)
	
	// Read session created message
	var sessionMsg WSMessage
	err = conn.ReadJSON(&sessionMsg)
	require.NoError(t, err)
	assert.Equal(t, "session_created", sessionMsg.Type)
	
	// Extract session ID
	sessionData, ok := sessionMsg.Data.(map[string]interface{})
	require.True(t, ok)
	sessionID, ok := sessionData["session_id"].(string)
	require.True(t, ok)
	assert.NotEmpty(t, sessionID)
	
	// Test command execution
	cmdMsg := WSMessage{
		Type: "command",
		Data: map[string]string{
			"command":    "status",
			"session_id": sessionID,
		},
		Timestamp: time.Now(),
	}
	
	err = conn.WriteJSON(cmdMsg)
	require.NoError(t, err)
	
	// Read command output
	var outputMsg WSMessage
	err = conn.ReadJSON(&outputMsg)
	require.NoError(t, err)
	assert.Equal(t, "command_output", outputMsg.Type)
	
	outputData, ok := outputMsg.Data.(map[string]interface{})
	require.True(t, ok)
	assert.Contains(t, outputData, "output")
	assert.Equal(t, float64(0), outputData["exit_code"]) // JSON numbers are float64
}

func TestServiceIntegrations(t *testing.T) {
	service := createTestTerminalService(t)
	ctx := context.Background()
	
	// Create session
	session, err := service.sessionManager.CreateSession(ctx, "test-user", nil)
	require.NoError(t, err)
	
	// Test AI integration through analyze command
	result, err := service.commandRegistry.ExecuteCommand(ctx, "analyze BTC", session)
	require.NoError(t, err)
	assert.Equal(t, 0, result.ExitCode)
	assert.Contains(t, result.Output, "AI Analysis")
	assert.Contains(t, result.Metadata, "symbol")
	assert.Equal(t, "BTC", result.Metadata["symbol"])
	
	// Test trading integration through buy command
	result, err = service.commandRegistry.ExecuteCommand(ctx, "buy BTC 0.1", session)
	require.NoError(t, err)
	assert.Equal(t, 0, result.ExitCode)
	assert.Contains(t, result.Output, "Buy Order Placed")
	assert.Contains(t, result.Metadata, "order_id")
	
	// Test configuration commands
	result, err = service.commandRegistry.ExecuteCommand(ctx, "config set theme dark", session)
	require.NoError(t, err)
	assert.Equal(t, 0, result.ExitCode)
	
	result, err = service.commandRegistry.ExecuteCommand(ctx, "config get theme", session)
	require.NoError(t, err)
	assert.Equal(t, 0, result.ExitCode)
	assert.Contains(t, result.Output, "dark")
	
	// Test alias commands
	result, err = service.commandRegistry.ExecuteCommand(ctx, "alias p \"price BTC\"", session)
	require.NoError(t, err)
	assert.Equal(t, 0, result.ExitCode)
	
	result, err = service.commandRegistry.ExecuteCommand(ctx, "alias", session)
	require.NoError(t, err)
	assert.Equal(t, 0, result.ExitCode)
	assert.Contains(t, result.Output, "p")
}

func TestCommandAutocomplete(t *testing.T) {
	service := createTestTerminalService(t)
	ctx := context.Background()
	
	// Test command name completion
	suggestions, err := service.commandRegistry.GetAutocompleteSuggestions(ctx, "st")
	require.NoError(t, err)
	
	found := false
	for _, suggestion := range suggestions {
		if suggestion == "status" {
			found = true
			break
		}
	}
	assert.True(t, found, "Status command should be in autocomplete suggestions")
	
	// Test argument completion for price command
	suggestions, err = service.commandRegistry.GetAutocompleteSuggestions(ctx, "price ")
	require.NoError(t, err)
	
	found = false
	for _, suggestion := range suggestions {
		if suggestion == "BTC" {
			found = true
			break
		}
	}
	assert.True(t, found, "BTC should be in price command autocomplete suggestions")
}

func TestErrorHandling(t *testing.T) {
	service := createTestTerminalService(t)
	ctx := context.Background()
	
	// Create session
	session, err := service.sessionManager.CreateSession(ctx, "test-user", nil)
	require.NoError(t, err)
	
	// Test invalid command
	result, err := service.commandRegistry.ExecuteCommand(ctx, "invalid-command", session)
	require.NoError(t, err)
	assert.NotEqual(t, 0, result.ExitCode)
	assert.Contains(t, result.Error, "command not found")
	
	// Test command with invalid arguments
	result, err = service.commandRegistry.ExecuteCommand(ctx, "buy", session)
	require.NoError(t, err)
	assert.NotEqual(t, 0, result.ExitCode)
	assert.Contains(t, result.Error, "Usage:")
	
	// Test empty command
	result, err = service.commandRegistry.ExecuteCommand(ctx, "", session)
	require.NoError(t, err)
	assert.Equal(t, 0, result.ExitCode) // Empty command should succeed with no output
}

// Helper function to create test terminal service
func createTestTerminalService(t *testing.T) *Service {
	logger := observability.NewLogger(config.ObservabilityConfig{
		ServiceName: "terminal-test",
		LogLevel:    "info",
		LogFormat:   "json",
	})
	
	config := Config{
		Host:         "localhost",
		Port:         8085,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		MaxSessions:  10,
		SessionTTL:   24 * time.Hour,
	}
	
	service, err := NewService(config, logger)
	require.NoError(t, err)
	
	return service
}
