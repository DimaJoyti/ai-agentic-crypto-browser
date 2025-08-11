# Terminal API Reference

## Overview

The Terminal Service provides both REST API and WebSocket endpoints for programmatic access to terminal functionality.

**Base URL**: `http://localhost:8085`
**WebSocket URL**: `ws://localhost:8085/ws`

## Authentication

All API endpoints require authentication via JWT token in the Authorization header:

```
Authorization: Bearer <jwt_token>
```

## REST API Endpoints

### Health Check

#### GET /health

Check service health status.

**Response:**
```json
{
  "status": "healthy",
  "service": "terminal-service"
}
```

### Session Management

#### POST /api/v1/sessions

Create a new terminal session.

**Request Body:**
```json
{
  "user_id": "string",
  "environment": {
    "key": "value"
  }
}
```

**Response:**
```json
{
  "session_id": "uuid",
  "created_at": "2024-01-15T10:30:00Z"
}
```

#### GET /api/v1/sessions

List user sessions.

**Query Parameters:**
- `user_id` (required): User identifier

**Response:**
```json
[
  {
    "id": "uuid",
    "user_id": "string",
    "created_at": "2024-01-15T10:30:00Z",
    "last_active": "2024-01-15T10:35:00Z",
    "environment": {},
    "state": {
      "current_directory": "/",
      "variables": {},
      "aliases": {},
      "last_command": "status",
      "exit_code": 0
    }
  }
]
```

#### GET /api/v1/sessions/{sessionId}

Get specific session details.

**Response:**
```json
{
  "id": "uuid",
  "user_id": "string",
  "created_at": "2024-01-15T10:30:00Z",
  "last_active": "2024-01-15T10:35:00Z",
  "environment": {},
  "history": [
    {
      "id": "uuid",
      "command": "status",
      "args": [],
      "output": "System Status: ✅ Healthy",
      "exit_code": 0,
      "start_time": "2024-01-15T10:30:00Z",
      "end_time": "2024-01-15T10:30:01Z",
      "duration": 1000
    }
  ],
  "state": {}
}
```

#### DELETE /api/v1/sessions/{sessionId}

Delete a session.

**Response:** `204 No Content`

#### GET /api/v1/sessions/{sessionId}/history

Get command history for a session.

**Response:**
```json
[
  {
    "id": "uuid",
    "command": "status",
    "args": [],
    "output": "System Status: ✅ Healthy",
    "error": "",
    "exit_code": 0,
    "start_time": "2024-01-15T10:30:00Z",
    "end_time": "2024-01-15T10:30:01Z",
    "duration": 1000
  }
]
```

### Command Information

#### GET /api/v1/commands

List all available commands.

**Response:**
```json
[
  {
    "name": "status",
    "description": "Show system status and health information",
    "usage": "status [--verbose] [--json] [--services]",
    "category": "system",
    "examples": [
      "status",
      "status --verbose"
    ]
  }
]
```

#### GET /api/v1/commands/{command}/help

Get help for a specific command.

**Response:**
```json
{
  "name": "status",
  "description": "Show system status and health information",
  "usage": "status [--verbose] [--json] [--services]",
  "category": "system",
  "examples": [
    "status",
    "status --verbose",
    "status --json"
  ]
}
```

## WebSocket API

### Connection

Connect to the WebSocket endpoint:

```javascript
const ws = new WebSocket('ws://localhost:8085/ws');
```

### Message Format

All WebSocket messages use the following format:

```json
{
  "type": "string",
  "session_id": "string",
  "data": {},
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Message Types

#### Client to Server

##### `session_create`

Create a new session.

```json
{
  "type": "session_create",
  "data": {
    "user_id": "string",
    "environment": {}
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

##### `session_join`

Join an existing session.

```json
{
  "type": "session_join",
  "data": {
    "session_id": "uuid"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

##### `command`

Execute a command.

```json
{
  "type": "command",
  "data": {
    "command": "status --verbose",
    "session_id": "uuid"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

#### Server to Client

##### `welcome`

Connection established.

```json
{
  "type": "welcome",
  "data": {
    "message": "Connected to terminal service"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

##### `session_created`

Session created successfully.

```json
{
  "type": "session_created",
  "data": {
    "session_id": "uuid",
    "created_at": "2024-01-15T10:30:00Z"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

##### `session_joined`

Joined session successfully.

```json
{
  "type": "session_joined",
  "data": {
    "session_id": "uuid",
    "session": {}
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

##### `command_output`

Command execution result.

```json
{
  "type": "command_output",
  "data": {
    "output": "System Status: ✅ Healthy\n...",
    "error": "",
    "exit_code": 0,
    "session_id": "uuid",
    "streaming": false
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

##### `error`

Error message.

```json
{
  "type": "error",
  "data": {
    "message": "Session not found"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## Client Libraries

### JavaScript/TypeScript

```typescript
class TerminalClient {
  private ws: WebSocket;
  private sessionId: string | null = null;

  constructor(url: string) {
    this.ws = new WebSocket(url);
    this.setupEventHandlers();
  }

  private setupEventHandlers() {
    this.ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      this.handleMessage(message);
    };
  }

  createSession(userId: string): Promise<string> {
    return new Promise((resolve, reject) => {
      const message = {
        type: 'session_create',
        data: { user_id: userId },
        timestamp: new Date().toISOString()
      };
      
      this.ws.send(JSON.stringify(message));
      
      // Handle response...
    });
  }

  executeCommand(command: string): Promise<CommandResult> {
    return new Promise((resolve, reject) => {
      if (!this.sessionId) {
        reject(new Error('No active session'));
        return;
      }

      const message = {
        type: 'command',
        data: {
          command,
          session_id: this.sessionId
        },
        timestamp: new Date().toISOString()
      };

      this.ws.send(JSON.stringify(message));
      
      // Handle response...
    });
  }
}
```

### Python

```python
import asyncio
import json
import websockets
from typing import Dict, Any, Optional

class TerminalClient:
    def __init__(self, url: str):
        self.url = url
        self.ws = None
        self.session_id: Optional[str] = None

    async def connect(self):
        self.ws = await websockets.connect(self.url)

    async def create_session(self, user_id: str) -> str:
        message = {
            "type": "session_create",
            "data": {"user_id": user_id},
            "timestamp": datetime.utcnow().isoformat() + "Z"
        }
        
        await self.ws.send(json.dumps(message))
        
        # Wait for response
        response = await self.ws.recv()
        data = json.loads(response)
        
        if data["type"] == "session_created":
            self.session_id = data["data"]["session_id"]
            return self.session_id
        else:
            raise Exception(f"Failed to create session: {data}")

    async def execute_command(self, command: str) -> Dict[str, Any]:
        if not self.session_id:
            raise Exception("No active session")

        message = {
            "type": "command",
            "data": {
                "command": command,
                "session_id": self.session_id
            },
            "timestamp": datetime.utcnow().isoformat() + "Z"
        }

        await self.ws.send(json.dumps(message))
        
        # Wait for response
        response = await self.ws.recv()
        return json.loads(response)
```

### Go

```go
package main

import (
    "encoding/json"
    "fmt"
    "github.com/gorilla/websocket"
    "time"
)

type TerminalClient struct {
    conn      *websocket.Conn
    sessionID string
}

type WSMessage struct {
    Type      string      `json:"type"`
    SessionID string      `json:"session_id,omitempty"`
    Data      interface{} `json:"data"`
    Timestamp time.Time   `json:"timestamp"`
}

func NewTerminalClient(url string) (*TerminalClient, error) {
    conn, _, err := websocket.DefaultDialer.Dial(url, nil)
    if err != nil {
        return nil, err
    }

    return &TerminalClient{conn: conn}, nil
}

func (c *TerminalClient) CreateSession(userID string) error {
    message := WSMessage{
        Type: "session_create",
        Data: map[string]string{"user_id": userID},
        Timestamp: time.Now(),
    }

    return c.conn.WriteJSON(message)
}

func (c *TerminalClient) ExecuteCommand(command string) error {
    message := WSMessage{
        Type: "command",
        Data: map[string]string{
            "command":    command,
            "session_id": c.sessionID,
        },
        Timestamp: time.Now(),
    }

    return c.conn.WriteJSON(message)
}
```

## Error Codes

| Code | Description |
|------|-------------|
| 400  | Bad Request - Invalid request format |
| 401  | Unauthorized - Invalid or missing authentication |
| 403  | Forbidden - Insufficient permissions |
| 404  | Not Found - Session or command not found |
| 429  | Too Many Requests - Rate limit exceeded |
| 500  | Internal Server Error - Server error |

## Rate Limiting

API endpoints are rate limited:

- **REST API**: 100 requests per minute per user
- **WebSocket**: 50 commands per minute per session
- **Command Execution**: 10 concurrent commands per session

Rate limit headers are included in responses:

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1642248600
```

## Examples

### Complete Session Workflow

```bash
# 1. Create session
curl -X POST http://localhost:8085/api/v1/sessions \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user123"}'

# 2. Execute commands via WebSocket
# (See client library examples above)

# 3. Get session history
curl -X GET http://localhost:8085/api/v1/sessions/$SESSION_ID/history \
  -H "Authorization: Bearer $JWT_TOKEN"

# 4. Delete session
curl -X DELETE http://localhost:8085/api/v1/sessions/$SESSION_ID \
  -H "Authorization: Bearer $JWT_TOKEN"
```

For more examples and integration guides, see the [Developer Guide](TERMINAL_DEVELOPER_GUIDE.md).
