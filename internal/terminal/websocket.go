package terminal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/gorilla/websocket"
)

// WebSocketManager manages WebSocket connections for terminal sessions
type WebSocketManager struct {
	logger          *observability.Logger
	upgrader        websocket.Upgrader
	clients         map[*websocket.Conn]*Client
	broadcast       chan []byte
	register        chan *Client
	unregister      chan *Client
	mu              sync.RWMutex
	commandRegistry *CommandRegistry
	sessionManager  *SessionManager
}

// Client represents a WebSocket client connection
type Client struct {
	conn      *websocket.Conn
	send      chan []byte
	sessionID string
	userID    string
	manager   *WebSocketManager
}

// WSMessage represents a WebSocket message
type WSMessage struct {
	Type      string      `json:"type"`
	SessionID string      `json:"session_id,omitempty"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// WSCommandMessage represents a command execution message
type WSCommandMessage struct {
	Command   string `json:"command"`
	SessionID string `json:"session_id"`
}

// WSOutputMessage represents command output message
type WSOutputMessage struct {
	Output    string `json:"output"`
	Error     string `json:"error,omitempty"`
	ExitCode  int    `json:"exit_code"`
	SessionID string `json:"session_id"`
	Streaming bool   `json:"streaming"`
}

// WSSessionMessage represents session-related messages
type WSSessionMessage struct {
	SessionID string   `json:"session_id"`
	Action    string   `json:"action"` // create, delete, update
	Session   *Session `json:"session,omitempty"`
}

// NewWebSocketManager creates a new WebSocket manager
func NewWebSocketManager(logger *observability.Logger, commandRegistry *CommandRegistry, sessionManager *SessionManager) *WebSocketManager {
	return &WebSocketManager{
		logger: logger,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// TODO: Implement proper origin checking for production
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		clients:         make(map[*websocket.Conn]*Client),
		broadcast:       make(chan []byte, 256),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
		commandRegistry: commandRegistry,
		sessionManager:  sessionManager,
	}
}

// Start starts the WebSocket manager
func (wsm *WebSocketManager) Start(ctx context.Context) error {
	wsm.logger.Info(ctx, "Starting WebSocket manager")

	go wsm.run()

	return nil
}

// Shutdown shuts down the WebSocket manager
func (wsm *WebSocketManager) Shutdown(ctx context.Context) error {
	wsm.logger.Info(ctx, "Shutting down WebSocket manager")

	wsm.mu.Lock()
	defer wsm.mu.Unlock()

	// Close all client connections
	for conn, client := range wsm.clients {
		close(client.send)
		conn.Close()
		delete(wsm.clients, conn)
	}

	return nil
}

// UpgradeConnection upgrades an HTTP connection to WebSocket
func (wsm *WebSocketManager) UpgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := wsm.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to upgrade connection: %w", err)
	}

	return conn, nil
}

// HandleConnection handles a new WebSocket connection
func (wsm *WebSocketManager) HandleConnection(ctx context.Context, conn *websocket.Conn) {
	// Extract user info from context or query params
	userID := "anonymous" // TODO: Extract from JWT token

	client := &Client{
		conn:    conn,
		send:    make(chan []byte, 256),
		userID:  userID,
		manager: wsm,
	}

	wsm.register <- client

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()
}

// BroadcastToSession sends a message to all clients in a session
func (wsm *WebSocketManager) BroadcastToSession(sessionID string, message WSMessage) {
	wsm.mu.RLock()
	defer wsm.mu.RUnlock()

	data, err := json.Marshal(message)
	if err != nil {
		wsm.logger.Error(context.Background(), "Failed to marshal message", err)
		return
	}

	for _, client := range wsm.clients {
		if client.sessionID == sessionID {
			select {
			case client.send <- data:
			default:
				close(client.send)
				delete(wsm.clients, client.conn)
			}
		}
	}
}

// BroadcastToUser sends a message to all clients for a user
func (wsm *WebSocketManager) BroadcastToUser(userID string, message WSMessage) {
	wsm.mu.RLock()
	defer wsm.mu.RUnlock()

	data, err := json.Marshal(message)
	if err != nil {
		wsm.logger.Error(context.Background(), "Failed to marshal message", err)
		return
	}

	for _, client := range wsm.clients {
		if client.userID == userID {
			select {
			case client.send <- data:
			default:
				close(client.send)
				delete(wsm.clients, client.conn)
			}
		}
	}
}

// run handles the main WebSocket manager loop
func (wsm *WebSocketManager) run() {
	for {
		select {
		case client := <-wsm.register:
			wsm.mu.Lock()
			wsm.clients[client.conn] = client
			wsm.mu.Unlock()

			wsm.logger.Info(context.Background(), "Client connected", map[string]interface{}{
				"user_id": client.userID,
			})

			// Send welcome message
			welcome := WSMessage{
				Type:      "welcome",
				Data:      map[string]string{"message": "Connected to terminal service"},
				Timestamp: time.Now(),
			}

			data, _ := json.Marshal(welcome)
			select {
			case client.send <- data:
			default:
				close(client.send)
				delete(wsm.clients, client.conn)
			}

		case client := <-wsm.unregister:
			wsm.mu.Lock()
			if _, ok := wsm.clients[client.conn]; ok {
				delete(wsm.clients, client.conn)
				close(client.send)
			}
			wsm.mu.Unlock()

			wsm.logger.Info(context.Background(), "Client disconnected", map[string]interface{}{
				"user_id":    client.userID,
				"session_id": client.sessionID,
			})

		case message := <-wsm.broadcast:
			wsm.mu.RLock()
			for _, client := range wsm.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(wsm.clients, client.conn)
				}
			}
			wsm.mu.RUnlock()
		}
	}
}

// Client methods

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

// readPump pumps messages from the WebSocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.manager.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.manager.logger.Error(context.Background(), "WebSocket error", err)
			}
			break
		}

		// Parse message
		var message WSMessage
		if err := json.Unmarshal(messageBytes, &message); err != nil {
			c.manager.logger.Error(context.Background(), "Failed to parse message", err)
			continue
		}

		// Handle message based on type
		c.handleMessage(message)
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage handles incoming WebSocket messages
func (c *Client) handleMessage(message WSMessage) {
	ctx := context.Background()

	switch message.Type {
	case "command":
		// Handle command execution
		var cmdMsg WSCommandMessage
		if data, err := json.Marshal(message.Data); err == nil {
			json.Unmarshal(data, &cmdMsg)
		}

		c.sessionID = cmdMsg.SessionID
		c.executeCommand(ctx, cmdMsg.Command, cmdMsg.SessionID)

	case "session_create":
		// Handle session creation
		c.createSession(ctx)

	case "session_join":
		// Handle joining existing session
		var sessionMsg WSSessionMessage
		if data, err := json.Marshal(message.Data); err == nil {
			json.Unmarshal(data, &sessionMsg)
		}

		c.sessionID = sessionMsg.SessionID
		c.joinSession(ctx, sessionMsg.SessionID)

	default:
		c.manager.logger.Warn(ctx, "Unknown message type", map[string]interface{}{
			"type": message.Type,
		})
	}
}

// executeCommand executes a command and sends the result back
func (c *Client) executeCommand(ctx context.Context, command, sessionID string) {
	// Get session
	session, err := c.manager.sessionManager.GetSession(ctx, sessionID)
	if err != nil {
		c.sendError(fmt.Sprintf("Session not found: %s", sessionID))
		return
	}

	// Execute command
	result, err := c.manager.commandRegistry.ExecuteCommand(ctx, command, session)
	if err != nil {
		c.sendError(fmt.Sprintf("Command execution failed: %v", err))
		return
	}

	// Send result back to client
	response := WSOutputMessage{
		Output:    result.Output,
		Error:     result.Error,
		ExitCode:  result.ExitCode,
		SessionID: sessionID,
		Streaming: result.Streaming,
	}

	c.sendMessage("command_output", response)
}

// createSession creates a new session
func (c *Client) createSession(ctx context.Context) {
	session, err := c.manager.sessionManager.CreateSession(ctx, c.userID, nil)
	if err != nil {
		c.sendError(fmt.Sprintf("Failed to create session: %v", err))
		return
	}

	c.sessionID = session.ID

	response := WSSessionMessage{
		SessionID: session.ID,
		Action:    "create",
		Session:   session,
	}

	c.sendMessage("session_created", response)
}

// joinSession joins an existing session
func (c *Client) joinSession(ctx context.Context, sessionID string) {
	session, err := c.manager.sessionManager.GetSession(ctx, sessionID)
	if err != nil {
		c.sendError(fmt.Sprintf("Session not found: %s", sessionID))
		return
	}

	c.sessionID = sessionID

	response := WSSessionMessage{
		SessionID: sessionID,
		Action:    "join",
		Session:   session,
	}

	c.sendMessage("session_joined", response)
}

// sendMessage sends a message to the client
func (c *Client) sendMessage(messageType string, data interface{}) {
	message := WSMessage{
		Type:      messageType,
		Data:      data,
		Timestamp: time.Now(),
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		c.manager.logger.Error(context.Background(), "Failed to marshal message", err)
		return
	}

	select {
	case c.send <- messageBytes:
	default:
		close(c.send)
		delete(c.manager.clients, c.conn)
	}
}

// sendError sends an error message to the client
func (c *Client) sendError(errorMsg string) {
	c.sendMessage("error", map[string]string{"message": errorMsg})
}
