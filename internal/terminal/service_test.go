package terminal

import (
	"context"
	"testing"
	"time"

	"github.com/ai-agentic-browser/internal/config"
	"github.com/ai-agentic-browser/pkg/observability"
)

func TestTerminalService(t *testing.T) {
	// Create logger
	logger := observability.NewLogger(config.ObservabilityConfig{
		ServiceName: "terminal-test",
		LogLevel:    "info",
		LogFormat:   "json",
	})

	// Create service config
	config := Config{
		Host:         "localhost",
		Port:         8085,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		MaxSessions:  10,
		SessionTTL:   24 * time.Hour,
	}

	// Create service
	service, err := NewService(config, logger)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	// Test command registry
	commands := service.commandRegistry.ListCommands()
	if len(commands) == 0 {
		t.Error("No commands registered")
	}

	// Check for essential commands
	essentialCommands := []string{"status", "help", "clear", "exit", "buy", "sell", "price", "analyze"}
	commandMap := make(map[string]bool)
	for _, cmd := range commands {
		commandMap[cmd.Name] = true
	}

	for _, essential := range essentialCommands {
		if !commandMap[essential] {
			t.Errorf("Essential command '%s' not registered", essential)
		}
	}
}

func TestCommandExecution(t *testing.T) {
	// Create logger
	logger := observability.NewLogger(config.ObservabilityConfig{
		ServiceName: "terminal-test",
		LogLevel:    "info",
		LogFormat:   "json",
	})

	// Create service
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
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	// Create test session
	session := &Session{
		ID:          "test-session",
		UserID:      "test-user",
		CreatedAt:   time.Now(),
		LastActive:  time.Now(),
		Environment: make(map[string]string),
		History:     make([]CommandHistory, 0),
		State: SessionState{
			CurrentDirectory: "/",
			Variables:        make(map[string]string),
			Aliases:          make(map[string]string),
			LastCommand:      "",
			ExitCode:         0,
		},
	}

	// Test status command
	result, err := service.commandRegistry.ExecuteCommand(context.Background(), "status", session)
	if err != nil {
		t.Errorf("Failed to execute status command: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("Status command failed with exit code %d", result.ExitCode)
	}

	if result.Output == "" {
		t.Error("Status command returned empty output")
	}

	// Test help command
	result, err = service.commandRegistry.ExecuteCommand(context.Background(), "help", session)
	if err != nil {
		t.Errorf("Failed to execute help command: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("Help command failed with exit code %d", result.ExitCode)
	}

	// Test price command
	result, err = service.commandRegistry.ExecuteCommand(context.Background(), "price BTC", session)
	if err != nil {
		t.Errorf("Failed to execute price command: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("Price command failed with exit code %d", result.ExitCode)
	}

	// Test invalid command
	result, err = service.commandRegistry.ExecuteCommand(context.Background(), "invalid-command", session)
	if err != nil {
		t.Errorf("Unexpected error for invalid command: %v", err)
	}

	if result.ExitCode == 0 {
		t.Error("Invalid command should return non-zero exit code")
	}
}

func TestSessionManager(t *testing.T) {
	// Create logger
	logger := observability.NewLogger(config.ObservabilityConfig{
		ServiceName: "terminal-test",
		LogLevel:    "info",
		LogFormat:   "json",
	})

	// Create session manager
	smConfig := SessionManagerConfig{
		MaxSessions: 5,
		SessionTTL:  time.Hour,
	}

	manager := NewSessionManager(smConfig, logger)

	ctx := context.Background()

	// Test session creation
	session, err := manager.CreateSession(ctx, "test-user", nil)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	if session.ID == "" {
		t.Error("Session ID is empty")
	}

	if session.UserID != "test-user" {
		t.Errorf("Expected user ID 'test-user', got '%s'", session.UserID)
	}

	// Test session retrieval
	retrievedSession, err := manager.GetSession(ctx, session.ID)
	if err != nil {
		t.Errorf("Failed to retrieve session: %v", err)
	}

	if retrievedSession.ID != session.ID {
		t.Errorf("Retrieved session ID mismatch: expected %s, got %s", session.ID, retrievedSession.ID)
	}

	// Test session deletion
	err = manager.DeleteSession(ctx, session.ID)
	if err != nil {
		t.Errorf("Failed to delete session: %v", err)
	}

	// Verify session is deleted
	_, err = manager.GetSession(ctx, session.ID)
	if err == nil {
		t.Error("Session should be deleted but still exists")
	}
}

func TestCommandRegistry(t *testing.T) {
	// Create logger
	logger := observability.NewLogger(config.ObservabilityConfig{
		ServiceName: "terminal-test",
		LogLevel:    "info",
		LogFormat:   "json",
	})

	registry := NewCommandRegistry(logger)

	// Test command registration
	statusCmd := &StatusCommand{}
	registry.RegisterCommand(statusCmd)

	// Test command retrieval
	cmd, exists := registry.GetCommand("status")
	if !exists {
		t.Error("Status command not found after registration")
	}

	if cmd.Name() != "status" {
		t.Errorf("Expected command name 'status', got '%s'", cmd.Name())
	}

	// Test command listing
	commands := registry.ListCommands()
	if len(commands) != 1 {
		t.Errorf("Expected 1 command, got %d", len(commands))
	}

	// Test command unregistration
	registry.UnregisterCommand("status")
	_, exists = registry.GetCommand("status")
	if exists {
		t.Error("Status command should be unregistered")
	}
}

func TestAutocompleteSuggestions(t *testing.T) {
	// Create logger
	logger := observability.NewLogger(config.ObservabilityConfig{
		ServiceName: "terminal-test",
		LogLevel:    "info",
		LogFormat:   "json",
	})

	registry := NewCommandRegistry(logger)
	registry.RegisterCommand(&StatusCommand{})
	registry.RegisterCommand(&PriceCommand{})

	ctx := context.Background()

	// Test command name completion
	suggestions, err := registry.GetAutocompleteSuggestions(ctx, "st")
	if err != nil {
		t.Errorf("Failed to get autocomplete suggestions: %v", err)
	}

	found := false
	for _, suggestion := range suggestions {
		if suggestion == "status" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Status command not found in autocomplete suggestions")
	}

	// Test argument completion
	suggestions, err = registry.GetAutocompleteSuggestions(ctx, "price ")
	if err != nil {
		t.Errorf("Failed to get autocomplete suggestions for price command: %v", err)
	}

	// Should include cryptocurrency symbols
	found = false
	for _, suggestion := range suggestions {
		if suggestion == "BTC" {
			found = true
			break
		}
	}

	if !found {
		t.Error("BTC not found in price command autocomplete suggestions")
	}
}
