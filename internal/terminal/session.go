package terminal

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

// Session represents a terminal session
type Session struct {
	ID          string            `json:"id"`
	UserID      string            `json:"user_id"`
	CreatedAt   time.Time         `json:"created_at"`
	LastActive  time.Time         `json:"last_active"`
	Environment map[string]string `json:"environment"`
	History     []CommandHistory  `json:"history"`
	State       SessionState      `json:"state"`
	
	// Runtime state
	mu sync.RWMutex
}

// SessionState represents the current state of a session
type SessionState struct {
	CurrentDirectory string            `json:"current_directory"`
	Variables        map[string]string `json:"variables"`
	Aliases          map[string]string `json:"aliases"`
	LastCommand      string            `json:"last_command"`
	ExitCode         int               `json:"exit_code"`
}

// CommandHistory represents a command execution record
type CommandHistory struct {
	ID        string    `json:"id"`
	Command   string    `json:"command"`
	Args      []string  `json:"args"`
	Output    string    `json:"output"`
	Error     string    `json:"error,omitempty"`
	ExitCode  int       `json:"exit_code"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Duration  int64     `json:"duration_ms"`
}

// SessionManager manages terminal sessions
type SessionManager struct {
	logger   *observability.Logger
	config   SessionManagerConfig
	sessions map[string]*Session
	userSessions map[string][]string // userID -> sessionIDs
	mu       sync.RWMutex
	
	// Cleanup ticker
	cleanupTicker *time.Ticker
	stopCleanup   chan struct{}
}

// SessionManagerConfig contains session manager configuration
type SessionManagerConfig struct {
	MaxSessions int           `json:"max_sessions"`
	SessionTTL  time.Duration `json:"session_ttl"`
}

// NewSessionManager creates a new session manager
func NewSessionManager(config SessionManagerConfig, logger *observability.Logger) *SessionManager {
	return &SessionManager{
		logger:       logger,
		config:       config,
		sessions:     make(map[string]*Session),
		userSessions: make(map[string][]string),
		stopCleanup:  make(chan struct{}),
	}
}

// Start starts the session manager
func (sm *SessionManager) Start(ctx context.Context) error {
	sm.logger.Info(ctx, "Starting session manager")
	
	// Start cleanup routine
	sm.cleanupTicker = time.NewTicker(time.Minute)
	go sm.cleanupRoutine()
	
	return nil
}

// Shutdown shuts down the session manager
func (sm *SessionManager) Shutdown(ctx context.Context) error {
	sm.logger.Info(ctx, "Shutting down session manager")
	
	// Stop cleanup routine
	if sm.cleanupTicker != nil {
		sm.cleanupTicker.Stop()
	}
	close(sm.stopCleanup)
	
	// Close all sessions
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	for sessionID := range sm.sessions {
		delete(sm.sessions, sessionID)
	}
	
	return nil
}

// CreateSession creates a new terminal session
func (sm *SessionManager) CreateSession(ctx context.Context, userID string, environment map[string]string) (*Session, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	// Check session limits
	userSessionCount := len(sm.userSessions[userID])
	if userSessionCount >= sm.config.MaxSessions {
		return nil, fmt.Errorf("maximum sessions reached for user %s", userID)
	}
	
	// Create new session
	sessionID := uuid.New().String()
	now := time.Now()
	
	if environment == nil {
		environment = make(map[string]string)
	}
	
	session := &Session{
		ID:          sessionID,
		UserID:      userID,
		CreatedAt:   now,
		LastActive:  now,
		Environment: environment,
		History:     make([]CommandHistory, 0),
		State: SessionState{
			CurrentDirectory: "/",
			Variables:        make(map[string]string),
			Aliases:          make(map[string]string),
			LastCommand:      "",
			ExitCode:         0,
		},
	}
	
	// Store session
	sm.sessions[sessionID] = session
	sm.userSessions[userID] = append(sm.userSessions[userID], sessionID)
	
	sm.logger.Info(ctx, "Created new session", map[string]interface{}{
		"session_id": sessionID,
		"user_id":    userID,
	})
	
	return session, nil
}

// GetSession retrieves a session by ID
func (sm *SessionManager) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	session, exists := sm.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}
	
	// Update last active time
	session.mu.Lock()
	session.LastActive = time.Now()
	session.mu.Unlock()
	
	return session, nil
}

// DeleteSession deletes a session
func (sm *SessionManager) DeleteSession(ctx context.Context, sessionID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	session, exists := sm.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found: %s", sessionID)
	}
	
	// Remove from user sessions
	userSessions := sm.userSessions[session.UserID]
	for i, id := range userSessions {
		if id == sessionID {
			sm.userSessions[session.UserID] = append(userSessions[:i], userSessions[i+1:]...)
			break
		}
	}
	
	// Remove session
	delete(sm.sessions, sessionID)
	
	sm.logger.Info(ctx, "Deleted session", map[string]interface{}{
		"session_id": sessionID,
		"user_id":    session.UserID,
	})
	
	return nil
}

// GetUserSessions returns all sessions for a user
func (sm *SessionManager) GetUserSessions(ctx context.Context, userID string) []*Session {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	sessionIDs := sm.userSessions[userID]
	sessions := make([]*Session, 0, len(sessionIDs))
	
	for _, sessionID := range sessionIDs {
		if session, exists := sm.sessions[sessionID]; exists {
			sessions = append(sessions, session)
		}
	}
	
	return sessions
}

// AddCommandToHistory adds a command execution to session history
func (sm *SessionManager) AddCommandToHistory(ctx context.Context, sessionID string, history CommandHistory) error {
	session, err := sm.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}
	
	session.mu.Lock()
	defer session.mu.Unlock()
	
	// Add to history
	session.History = append(session.History, history)
	
	// Limit history size (keep last 1000 commands)
	if len(session.History) > 1000 {
		session.History = session.History[len(session.History)-1000:]
	}
	
	// Update session state
	session.State.LastCommand = history.Command
	session.State.ExitCode = history.ExitCode
	session.LastActive = time.Now()
	
	return nil
}

// GetSessionHistory returns command history for a session
func (sm *SessionManager) GetSessionHistory(ctx context.Context, sessionID string) ([]CommandHistory, error) {
	session, err := sm.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	
	session.mu.RLock()
	defer session.mu.RUnlock()
	
	// Return copy of history
	history := make([]CommandHistory, len(session.History))
	copy(history, session.History)
	
	return history, nil
}

// UpdateSessionState updates session state
func (sm *SessionManager) UpdateSessionState(ctx context.Context, sessionID string, state SessionState) error {
	session, err := sm.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}
	
	session.mu.Lock()
	defer session.mu.Unlock()
	
	session.State = state
	session.LastActive = time.Now()
	
	return nil
}

// cleanupRoutine periodically cleans up expired sessions
func (sm *SessionManager) cleanupRoutine() {
	for {
		select {
		case <-sm.cleanupTicker.C:
			sm.cleanupExpiredSessions()
		case <-sm.stopCleanup:
			return
		}
	}
}

// cleanupExpiredSessions removes sessions that have exceeded TTL
func (sm *SessionManager) cleanupExpiredSessions() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	now := time.Now()
	expiredSessions := make([]string, 0)
	
	for sessionID, session := range sm.sessions {
		if now.Sub(session.LastActive) > sm.config.SessionTTL {
			expiredSessions = append(expiredSessions, sessionID)
		}
	}
	
	// Remove expired sessions
	for _, sessionID := range expiredSessions {
		session := sm.sessions[sessionID]
		
		// Remove from user sessions
		userSessions := sm.userSessions[session.UserID]
		for i, id := range userSessions {
			if id == sessionID {
				sm.userSessions[session.UserID] = append(userSessions[:i], userSessions[i+1:]...)
				break
			}
		}
		
		// Remove session
		delete(sm.sessions, sessionID)
		
		sm.logger.Info(context.Background(), "Cleaned up expired session", map[string]interface{}{
			"session_id": sessionID,
			"user_id":    session.UserID,
		})
	}
}
