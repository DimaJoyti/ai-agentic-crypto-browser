# Terminal Developer Guide

## Architecture Overview

The AI-Agentic Crypto Browser Terminal is built with a modular, extensible architecture that supports real-time communication, command execution, and service integration.

### Core Components

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Client    │    │  Terminal API   │    │ Command Engine  │
│                 │    │                 │    │                 │
│ - React UI      │◄──►│ - REST API      │◄──►│ - Command Reg.  │
│ - WebSocket     │    │ - WebSocket     │    │ - Session Mgr.  │
│ - Terminal Comp │    │ - Auth Middleware│    │ - Integrations  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
                       │    Services     │
                       │                 │
                       │ - AI Service    │
                       │ - Trading Svc   │
                       │ - Web3 Service  │
                       │ - Browser Svc   │
                       └─────────────────┘
```

## Creating Custom Commands

### Command Interface

All commands must implement the `Command` interface:

```go
type Command interface {
    Name() string
    Description() string
    Usage() string
    Execute(ctx context.Context, args []string, session *Session) (*CommandResult, error)
    Autocomplete(ctx context.Context, args []string) ([]string, error)
}
```

### Example Command Implementation

```go
package commands

import (
    "context"
    "fmt"
    "strings"
    "time"
)

// CustomCommand demonstrates a basic command implementation
type CustomCommand struct {
    integrations *ServiceIntegrations
}

func (c *CustomCommand) Name() string {
    return "custom"
}

func (c *CustomCommand) Description() string {
    return "Example custom command"
}

func (c *CustomCommand) Usage() string {
    return "custom [--option value] <argument>"
}

func (c *CustomCommand) Execute(ctx context.Context, args []string, session *Session) (*CommandResult, error) {
    // Parse arguments
    if len(args) == 0 {
        return &CommandResult{
            Error:    "Usage: " + c.Usage(),
            ExitCode: 1,
        }, nil
    }
    
    // Process command logic
    var output strings.Builder
    output.WriteString("Custom command executed successfully\n")
    output.WriteString(fmt.Sprintf("Arguments: %v\n", args))
    output.WriteString(fmt.Sprintf("Session: %s\n", session.ID))
    output.WriteString(fmt.Sprintf("Timestamp: %s\n", time.Now().Format(time.RFC3339)))
    
    // Use service integrations if needed
    if c.integrations != nil && c.integrations.AI != nil {
        // Call AI service
        // result, err := c.integrations.AI.SomeMethod(ctx, args[0])
    }
    
    return &CommandResult{
        Output:   output.String(),
        ExitCode: 0,
        Metadata: map[string]string{
            "command": c.Name(),
            "args":    strings.Join(args, ","),
        },
    }, nil
}

func (c *CustomCommand) Autocomplete(ctx context.Context, args []string) ([]string, error) {
    if len(args) == 0 {
        return []string{"option1", "option2", "option3"}, nil
    }
    
    return []string{"--option", "--help"}, nil
}
```

### Registering Commands

Register your command in the terminal service:

```go
// In registerDefaultCommands()
s.commandRegistry.RegisterCommand(&CustomCommand{
    integrations: s.integrations,
})
```

## Service Integration

### Creating Service Clients

Implement service client interfaces for external service integration:

```go
type CustomServiceClient interface {
    DoSomething(ctx context.Context, param string) (*Result, error)
    GetData(ctx context.Context, id string) (*Data, error)
}

type customServiceClient struct {
    baseURL string
    client  *http.Client
}

func NewCustomServiceClient(baseURL string) CustomServiceClient {
    return &customServiceClient{
        baseURL: baseURL,
        client:  &http.Client{Timeout: 30 * time.Second},
    }
}

func (c *customServiceClient) DoSomething(ctx context.Context, param string) (*Result, error) {
    // Implement HTTP client logic
    url := fmt.Sprintf("%s/api/something?param=%s", c.baseURL, param)
    
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result Result
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    return &result, nil
}
```

### Adding Service to Integrations

```go
// Update ServiceIntegrations struct
type ServiceIntegrations struct {
    AI      AIServiceClient
    Trading TradingServiceClient
    Web3    Web3ServiceClient
    Browser BrowserServiceClient
    Auth    AuthServiceClient
    Custom  CustomServiceClient  // Add your service
}

// Update service initialization
func NewServiceIntegrations() *ServiceIntegrations {
    return &ServiceIntegrations{
        AI:      NewAIServiceClient("http://ai-service:8082"),
        Trading: NewTradingServiceClient("http://trading-service:8083"),
        // ... other services
        Custom:  NewCustomServiceClient("http://custom-service:8090"),
    }
}
```

## WebSocket Communication

### Message Types

The terminal uses structured WebSocket messages:

```go
type WSMessage struct {
    Type      string      `json:"type"`
    SessionID string      `json:"session_id,omitempty"`
    Data      interface{} `json:"data"`
    Timestamp time.Time   `json:"timestamp"`
}
```

### Supported Message Types

- `command`: Execute a command
- `session_create`: Create new session
- `session_join`: Join existing session
- `command_output`: Command execution result
- `error`: Error message

### Custom Message Handlers

Add custom WebSocket message handlers:

```go
// In handleMessage method
case "custom_message":
    var customData CustomMessageData
    if data, err := json.Marshal(message.Data); err == nil {
        json.Unmarshal(data, &customData)
    }
    
    c.handleCustomMessage(ctx, customData)
```

## Testing

### Unit Tests

Create comprehensive unit tests for commands:

```go
func TestCustomCommand(t *testing.T) {
    cmd := &CustomCommand{}
    
    session := &Session{
        ID:     "test-session",
        UserID: "test-user",
        State:  SessionState{},
    }
    
    // Test successful execution
    result, err := cmd.Execute(context.Background(), []string{"arg1"}, session)
    assert.NoError(t, err)
    assert.Equal(t, 0, result.ExitCode)
    assert.Contains(t, result.Output, "Custom command executed")
    
    // Test error cases
    result, err = cmd.Execute(context.Background(), []string{}, session)
    assert.NoError(t, err)
    assert.Equal(t, 1, result.ExitCode)
    assert.Contains(t, result.Error, "Usage:")
}
```

### Integration Tests

Test service integrations:

```go
func TestServiceIntegration(t *testing.T) {
    // Setup mock services
    mockAI := &MockAIServiceClient{}
    integrations := &ServiceIntegrations{AI: mockAI}
    
    cmd := &AnalyzeCommand{integrations: integrations}
    
    // Test with service integration
    result, err := cmd.Execute(context.Background(), []string{"BTC"}, session)
    assert.NoError(t, err)
    assert.Equal(t, 0, result.ExitCode)
}
```

### End-to-End Tests

Test complete workflows:

```go
func TestTerminalWorkflow(t *testing.T) {
    // Create terminal service
    service := createTestService(t)
    
    // Create session
    session := createTestSession(t, service)
    
    // Execute command sequence
    commands := []string{
        "status",
        "price BTC",
        "analyze BTC",
        "portfolio",
    }
    
    for _, cmd := range commands {
        result, err := service.commandRegistry.ExecuteCommand(
            context.Background(), cmd, session)
        assert.NoError(t, err)
        assert.Equal(t, 0, result.ExitCode)
    }
}
```

## Configuration

### Environment Variables

Configure the terminal service using environment variables:

```bash
# Terminal service configuration
TERMINAL_HOST=0.0.0.0
TERMINAL_PORT=8085
TERMINAL_READ_TIMEOUT=15s
TERMINAL_WRITE_TIMEOUT=15s
TERMINAL_MAX_SESSIONS=100
TERMINAL_SESSION_TTL=24h

# Service endpoints
AI_SERVICE_URL=http://ai-service:8082
TRADING_SERVICE_URL=http://trading-service:8083
WEB3_SERVICE_URL=http://web3-service:8084
```

### Configuration Structure

```go
type Config struct {
    Host         string        `json:"host"`
    Port         int           `json:"port"`
    ReadTimeout  time.Duration `json:"read_timeout"`
    WriteTimeout time.Duration `json:"write_timeout"`
    MaxSessions  int           `json:"max_sessions"`
    SessionTTL   time.Duration `json:"session_ttl"`
}
```

## Deployment

### Docker Configuration

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o terminal-service ./cmd/terminal-service

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/terminal-service .
EXPOSE 8085
CMD ["./terminal-service"]
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: terminal-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: terminal-service
  template:
    metadata:
      labels:
        app: terminal-service
    spec:
      containers:
      - name: terminal-service
        image: terminal-service:latest
        ports:
        - containerPort: 8085
        env:
        - name: TERMINAL_HOST
          value: "0.0.0.0"
        - name: TERMINAL_PORT
          value: "8085"
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
```

## Monitoring and Observability

### Metrics

The terminal service exposes metrics for monitoring:

- Command execution count
- Command execution duration
- Active sessions count
- WebSocket connections
- Error rates

### Logging

Structured logging with OpenTelemetry:

```go
logger.Info(ctx, "Command executed", map[string]interface{}{
    "command":   commandName,
    "user_id":   session.UserID,
    "session_id": session.ID,
    "duration":  duration.Milliseconds(),
    "exit_code": result.ExitCode,
})
```

### Tracing

Distributed tracing for command execution:

```go
func (c *CustomCommand) Execute(ctx context.Context, args []string, session *Session) (*CommandResult, error) {
    ctx, span := otel.Tracer("terminal").Start(ctx, "custom_command")
    defer span.End()
    
    span.SetAttributes(
        attribute.String("command", c.Name()),
        attribute.String("session_id", session.ID),
        attribute.StringSlice("args", args),
    )
    
    // Command logic...
    
    span.SetAttributes(attribute.Int("exit_code", result.ExitCode))
    return result, nil
}
```

## Security Considerations

### Authentication

Commands can check user permissions:

```go
func (c *AdminCommand) Execute(ctx context.Context, args []string, session *Session) (*CommandResult, error) {
    // Check user permissions
    if !hasPermission(session.UserID, "admin") {
        return &CommandResult{
            Error:    "Permission denied: admin access required",
            ExitCode: 403,
        }, nil
    }
    
    // Execute admin command...
}
```

### Input Validation

Always validate and sanitize user input:

```go
func validateSymbol(symbol string) error {
    if len(symbol) < 2 || len(symbol) > 10 {
        return fmt.Errorf("invalid symbol length")
    }
    
    if !regexp.MustCompile(`^[A-Z0-9]+$`).MatchString(symbol) {
        return fmt.Errorf("invalid symbol format")
    }
    
    return nil
}
```

### Rate Limiting

Implement rate limiting for commands:

```go
type RateLimitedCommand struct {
    Command
    limiter *rate.Limiter
}

func (c *RateLimitedCommand) Execute(ctx context.Context, args []string, session *Session) (*CommandResult, error) {
    if !c.limiter.Allow() {
        return &CommandResult{
            Error:    "Rate limit exceeded",
            ExitCode: 429,
        }, nil
    }
    
    return c.Command.Execute(ctx, args, session)
}
```

## Best Practices

1. **Error Handling**: Always provide clear, actionable error messages
2. **Performance**: Use context for timeouts and cancellation
3. **Logging**: Log important events with structured data
4. **Testing**: Write comprehensive tests for all commands
5. **Documentation**: Document command usage and examples
6. **Security**: Validate input and check permissions
7. **Monitoring**: Add metrics and tracing for observability

For more information, see the [Terminal API Reference](TERMINAL_API.md) and [User Guide](TERMINAL_USER_GUIDE.md).
