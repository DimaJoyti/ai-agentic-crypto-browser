package terminal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/gorilla/mux"
)

// Service provides terminal functionality
type Service struct {
	logger          *observability.Logger
	config          Config
	sessionManager  *SessionManager
	commandRegistry *CommandRegistry
	wsManager       *WebSocketManager
	integrations    *ServiceIntegrations

	// State
	isRunning bool
	mu        sync.RWMutex
}

// Config contains terminal service configuration
type Config struct {
	Host         string        `json:"host"`
	Port         int           `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
	MaxSessions  int           `json:"max_sessions"`
	SessionTTL   time.Duration `json:"session_ttl"`
}

// NewService creates a new terminal service
func NewService(config Config, logger *observability.Logger) (*Service, error) {
	// Initialize session manager
	sessionManager := NewSessionManager(SessionManagerConfig{
		MaxSessions: config.MaxSessions,
		SessionTTL:  config.SessionTTL,
	}, logger)

	// Initialize command registry
	commandRegistry := NewCommandRegistry(logger)

	// Initialize service integrations (using mock for now)
	integrations := NewMockServiceIntegrations()

	// Initialize WebSocket manager
	wsManager := NewWebSocketManager(logger, commandRegistry, sessionManager)

	service := &Service{
		logger:          logger,
		config:          config,
		sessionManager:  sessionManager,
		commandRegistry: commandRegistry,
		wsManager:       wsManager,
		integrations:    integrations,
		isRunning:       false,
	}

	// Register default commands
	if err := service.registerDefaultCommands(); err != nil {
		return nil, fmt.Errorf("failed to register default commands: %w", err)
	}

	return service, nil
}

// Start starts the terminal service
func (s *Service) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return fmt.Errorf("service is already running")
	}

	s.logger.Info(ctx, "Starting terminal service")

	// Start session manager
	if err := s.sessionManager.Start(ctx); err != nil {
		return fmt.Errorf("failed to start session manager: %w", err)
	}

	// Start WebSocket manager
	if err := s.wsManager.Start(ctx); err != nil {
		return fmt.Errorf("failed to start WebSocket manager: %w", err)
	}

	s.isRunning = true
	s.logger.Info(ctx, "Terminal service started successfully")

	return nil
}

// Shutdown gracefully shuts down the terminal service
func (s *Service) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return nil
	}

	s.logger.Info(ctx, "Shutting down terminal service")

	// Shutdown WebSocket manager
	if err := s.wsManager.Shutdown(ctx); err != nil {
		s.logger.Error(ctx, "Failed to shutdown WebSocket manager", err)
	}

	// Shutdown session manager
	if err := s.sessionManager.Shutdown(ctx); err != nil {
		s.logger.Error(ctx, "Failed to shutdown session manager", err)
	}

	s.isRunning = false
	s.logger.Info(ctx, "Terminal service shutdown complete")

	return nil
}

// HandleWebSocket handles WebSocket connections for terminal sessions
func (s *Service) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Upgrade connection to WebSocket
	conn, err := s.wsManager.UpgradeConnection(w, r)
	if err != nil {
		s.logger.Error(ctx, "Failed to upgrade WebSocket connection", err)
		http.Error(w, "Failed to upgrade connection", http.StatusBadRequest)
		return
	}

	// Handle the WebSocket connection
	s.wsManager.HandleConnection(ctx, conn)
}

// HandleCreateSession creates a new terminal session
func (s *Service) HandleCreateSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Error(ctx, "Failed to decode create session request", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create new session
	session, err := s.sessionManager.CreateSession(ctx, req.UserID, req.Environment)
	if err != nil {
		s.logger.Error(ctx, "Failed to create session", err)
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Return session info
	response := CreateSessionResponse{
		SessionID: session.ID,
		CreatedAt: session.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleListSessions lists all active sessions for a user
func (s *Service) HandleListSessions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := r.URL.Query().Get("user_id")

	if userID == "" {
		http.Error(w, "user_id parameter is required", http.StatusBadRequest)
		return
	}

	sessions := s.sessionManager.GetUserSessions(ctx, userID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessions)
}

// HandleGetSession gets a specific session
func (s *Service) HandleGetSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]

	session, err := s.sessionManager.GetSession(ctx, sessionID)
	if err != nil {
		s.logger.Error(ctx, "Failed to get session", err)
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

// HandleDeleteSession deletes a session
func (s *Service) HandleDeleteSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]

	if err := s.sessionManager.DeleteSession(ctx, sessionID); err != nil {
		s.logger.Error(ctx, "Failed to delete session", err)
		http.Error(w, "Failed to delete session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleGetHistory gets command history for a session
func (s *Service) HandleGetHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]

	history, err := s.sessionManager.GetSessionHistory(ctx, sessionID)
	if err != nil {
		s.logger.Error(ctx, "Failed to get session history", err)
		http.Error(w, "Failed to get history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

// HandleListCommands lists all available commands
func (s *Service) HandleListCommands(w http.ResponseWriter, r *http.Request) {
	commands := s.commandRegistry.ListCommands()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(commands)
}

// HandleGetCommandHelp gets help for a specific command
func (s *Service) HandleGetCommandHelp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	commandName := vars["command"]

	help, err := s.commandRegistry.GetCommandHelp(ctx, commandName)
	if err != nil {
		s.logger.Error(ctx, "Failed to get command help", err)
		http.Error(w, "Command not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(help)
}

// registerDefaultCommands registers the default set of commands
func (s *Service) registerDefaultCommands() error {
	// System commands
	s.commandRegistry.RegisterCommand(&StatusCommand{})
	s.commandRegistry.RegisterCommand(&HelpCommand{registry: s.commandRegistry})
	s.commandRegistry.RegisterCommand(&ClearCommand{})
	s.commandRegistry.RegisterCommand(&ExitCommand{})
	s.commandRegistry.RegisterCommand(&ConfigCommand{})
	s.commandRegistry.RegisterCommand(&AliasCommand{})
	s.commandRegistry.RegisterCommand(&HistoryCommand{})

	// Advanced features
	s.commandRegistry.RegisterCommand(&ScriptCommand{})
	s.commandRegistry.RegisterCommand(&WatchCommand{})
	s.commandRegistry.RegisterCommand(&ExportCommand{})

	// Trading commands
	s.commandRegistry.RegisterCommand(&BuyCommand{})
	s.commandRegistry.RegisterCommand(&SellCommand{})
	s.commandRegistry.RegisterCommand(&PortfolioCommand{})
	s.commandRegistry.RegisterCommand(&OrdersCommand{})

	// Market commands
	s.commandRegistry.RegisterCommand(&PriceCommand{integrations: s.integrations})
	s.commandRegistry.RegisterCommand(&ChartCommand{})

	// AI commands
	s.commandRegistry.RegisterCommand(&AnalyzeCommand{integrations: s.integrations})

	return nil
}

// Request/Response types
type CreateSessionRequest struct {
	UserID      string            `json:"user_id"`
	Environment map[string]string `json:"environment,omitempty"`
}

type CreateSessionResponse struct {
	SessionID string    `json:"session_id"`
	CreatedAt time.Time `json:"created_at"`
}
