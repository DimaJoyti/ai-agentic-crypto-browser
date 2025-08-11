package security

import (
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

// SessionManager manages user sessions
type SessionManager struct {
	logger   *observability.Logger
	sessions map[string]*SecuritySession
	config   *SecurityConfig
	mu       sync.RWMutex
}

// NewSessionManager creates a new session manager
func NewSessionManager(logger *observability.Logger, config *SecurityConfig) *SessionManager {
	return &SessionManager{
		logger:   logger,
		sessions: make(map[string]*SecuritySession),
		config:   config,
	}
}

// CreateSession creates a new session
func (sm *SessionManager) CreateSession(userID uuid.UUID, deviceID, ipAddress, userAgent string) *SecuritySession {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sessionID := uuid.New().String()
	now := time.Now()

	session := &SecuritySession{
		ID:           sessionID,
		UserID:       userID,
		DeviceID:     deviceID,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		CreatedAt:    now,
		LastActivity: now,
		ExpiresAt:    now.Add(sm.config.SessionTimeout),
		IsActive:     true,
		MFAVerified:  false,
		DeviceTrusted: false,
		RiskScore:    50, // Default medium risk
		Permissions:  []string{},
		TradingEnabled: false,
		SecurityFlags: make(map[string]interface{}),
	}

	sm.sessions[sessionID] = session
	return session
}

// GetSession retrieves a session by ID
func (sm *SessionManager) GetSession(sessionID string) (*SecuritySession, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	session, exists := sm.sessions[sessionID]
	return session, exists
}

// UpdateSession updates a session
func (sm *SessionManager) UpdateSession(session *SecuritySession) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session.LastActivity = time.Now()
	sm.sessions[session.ID] = session
}

// RevokeSession revokes a session
func (sm *SessionManager) RevokeSession(sessionID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if session, exists := sm.sessions[sessionID]; exists {
		session.IsActive = false
	}
}

// CleanupExpiredSessions removes expired sessions
func (sm *SessionManager) CleanupExpiredSessions() int {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	now := time.Now()
	cleanedCount := 0

	for sessionID, session := range sm.sessions {
		if now.After(session.ExpiresAt) {
			delete(sm.sessions, sessionID)
			cleanedCount++
		}
	}

	return cleanedCount
}

// GetActiveSessions returns all active sessions
func (sm *SessionManager) GetActiveSessions() []*SecuritySession {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	var activeSessions []*SecuritySession
	now := time.Now()

	for _, session := range sm.sessions {
		if session.IsActive && now.Before(session.ExpiresAt) {
			activeSessions = append(activeSessions, session)
		}
	}

	return activeSessions
}

// GetUserSessions returns all sessions for a user
func (sm *SessionManager) GetUserSessions(userID uuid.UUID) []*SecuritySession {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	var userSessions []*SecuritySession

	for _, session := range sm.sessions {
		if session.UserID == userID {
			userSessions = append(userSessions, session)
		}
	}

	return userSessions
}

// RevokeUserSessions revokes all sessions for a user
func (sm *SessionManager) RevokeUserSessions(userID uuid.UUID) int {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	revokedCount := 0

	for _, session := range sm.sessions {
		if session.UserID == userID && session.IsActive {
			session.IsActive = false
			revokedCount++
		}
	}

	return revokedCount
}

// GetSessionCount returns the total number of sessions
func (sm *SessionManager) GetSessionCount() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return len(sm.sessions)
}

// GetActiveSessionCount returns the number of active sessions
func (sm *SessionManager) GetActiveSessionCount() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	activeCount := 0
	now := time.Now()

	for _, session := range sm.sessions {
		if session.IsActive && now.Before(session.ExpiresAt) {
			activeCount++
		}
	}

	return activeCount
}
