package terminal

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
)

// Command interface defines the contract for terminal commands
type Command interface {
	Name() string
	Description() string
	Usage() string
	Execute(ctx context.Context, args []string, session *Session) (*CommandResult, error)
	Autocomplete(ctx context.Context, args []string) ([]string, error)
}

// CommandResult represents the result of command execution
type CommandResult struct {
	Output    string            `json:"output"`
	Error     string            `json:"error,omitempty"`
	ExitCode  int               `json:"exit_code"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	Streaming bool              `json:"streaming"`
}

// CommandInfo provides information about a command
type CommandInfo struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Usage       string   `json:"usage"`
	Category    string   `json:"category"`
	Examples    []string `json:"examples,omitempty"`
}

// CommandRegistry manages available terminal commands
type CommandRegistry struct {
	logger   *observability.Logger
	commands map[string]Command
	mu       sync.RWMutex
}

// NewCommandRegistry creates a new command registry
func NewCommandRegistry(logger *observability.Logger) *CommandRegistry {
	return &CommandRegistry{
		logger:   logger,
		commands: make(map[string]Command),
	}
}

// RegisterCommand registers a new command
func (cr *CommandRegistry) RegisterCommand(cmd Command) {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	cr.commands[cmd.Name()] = cmd
	cr.logger.Info(context.Background(), "Registered command", map[string]interface{}{
		"command": cmd.Name(),
	})
}

// UnregisterCommand removes a command from the registry
func (cr *CommandRegistry) UnregisterCommand(name string) {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	delete(cr.commands, name)
	cr.logger.Info(context.Background(), "Unregistered command", map[string]interface{}{
		"command": name,
	})
}

// GetCommand retrieves a command by name
func (cr *CommandRegistry) GetCommand(name string) (Command, bool) {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	cmd, exists := cr.commands[name]
	return cmd, exists
}

// ListCommands returns information about all registered commands
func (cr *CommandRegistry) ListCommands() []CommandInfo {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	commands := make([]CommandInfo, 0, len(cr.commands))
	for _, cmd := range cr.commands {
		commands = append(commands, CommandInfo{
			Name:        cmd.Name(),
			Description: cmd.Description(),
			Usage:       cmd.Usage(),
			Category:    getCommandCategory(cmd.Name()),
		})
	}

	return commands
}

// GetCommandHelp returns help information for a specific command
func (cr *CommandRegistry) GetCommandHelp(ctx context.Context, name string) (*CommandInfo, error) {
	cmd, exists := cr.GetCommand(name)
	if !exists {
		return nil, fmt.Errorf("command not found: %s", name)
	}

	return &CommandInfo{
		Name:        cmd.Name(),
		Description: cmd.Description(),
		Usage:       cmd.Usage(),
		Category:    getCommandCategory(cmd.Name()),
		Examples:    getCommandExamples(cmd.Name()),
	}, nil
}

// ExecuteCommand executes a command with the given arguments
func (cr *CommandRegistry) ExecuteCommand(ctx context.Context, commandLine string, session *Session) (*CommandResult, error) {
	// Parse command line
	parts := strings.Fields(commandLine)
	if len(parts) == 0 {
		return &CommandResult{
			Output:   "",
			ExitCode: 0,
		}, nil
	}

	commandName := parts[0]
	args := parts[1:]

	// Get command
	cmd, exists := cr.GetCommand(commandName)
	if !exists {
		// Create history entry for invalid command
		historyEntry := CommandHistory{
			ID:        fmt.Sprintf("cmd-%d", time.Now().UnixNano()),
			Command:   commandName,
			Args:      args,
			StartTime: time.Now(),
			EndTime:   time.Now(),
			Duration:  0,
			Error:     fmt.Sprintf("command not found: %s", commandName),
			ExitCode:  127,
		}

		// Add to session history
		session.History = append(session.History, historyEntry)

		return &CommandResult{
			Error:    fmt.Sprintf("command not found: %s", commandName),
			ExitCode: 127,
		}, nil
	}

	// Execute command
	startTime := time.Now()
	result, err := cmd.Execute(ctx, args, session)
	endTime := time.Now()
	duration := endTime.Sub(startTime)

	// Create history entry
	historyEntry := CommandHistory{
		ID:        fmt.Sprintf("cmd-%d", time.Now().UnixNano()),
		Command:   commandName,
		Args:      args,
		StartTime: startTime,
		EndTime:   endTime,
		Duration:  duration.Milliseconds(),
	}

	if err != nil {
		historyEntry.Error = err.Error()
		historyEntry.ExitCode = 1

		cr.logger.Error(ctx, "Command execution failed", err, map[string]interface{}{
			"command":  commandName,
			"args":     args,
			"duration": duration.Milliseconds(),
		})

		// Add to session history
		session.History = append(session.History, historyEntry)

		return &CommandResult{
			Error:    err.Error(),
			ExitCode: 1,
		}, nil
	}

	// Update history entry with result
	historyEntry.Output = result.Output
	historyEntry.Error = result.Error
	historyEntry.ExitCode = result.ExitCode

	// Add to session history
	session.History = append(session.History, historyEntry)

	cr.logger.Info(ctx, "Command executed successfully", map[string]interface{}{
		"command":   commandName,
		"args":      args,
		"duration":  duration.Milliseconds(),
		"exit_code": result.ExitCode,
	})

	return result, nil
}

// GetAutocompleteSuggestions returns autocomplete suggestions for a command line
func (cr *CommandRegistry) GetAutocompleteSuggestions(ctx context.Context, commandLine string) ([]string, error) {
	parts := strings.Fields(commandLine)
	if len(parts) == 0 {
		// Return all command names
		return cr.getAllCommandNames(), nil
	}

	commandName := parts[0]

	// If we're still typing the command name
	if len(parts) == 1 && !strings.HasSuffix(commandLine, " ") {
		return cr.getCommandNameSuggestions(commandName), nil
	}

	// Get command-specific autocomplete
	cmd, exists := cr.GetCommand(commandName)
	if !exists {
		return nil, nil
	}

	args := parts[1:]
	return cmd.Autocomplete(ctx, args)
}

// Helper functions

func (cr *CommandRegistry) getAllCommandNames() []string {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	names := make([]string, 0, len(cr.commands))
	for name := range cr.commands {
		names = append(names, name)
	}

	return names
}

func (cr *CommandRegistry) getCommandNameSuggestions(prefix string) []string {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	suggestions := make([]string, 0)
	for name := range cr.commands {
		if strings.HasPrefix(name, prefix) {
			suggestions = append(suggestions, name)
		}
	}

	return suggestions
}

func getCommandCategory(name string) string {
	categories := map[string]string{
		"status":    "system",
		"help":      "system",
		"clear":     "system",
		"exit":      "system",
		"config":    "system",
		"logs":      "system",
		"health":    "system",
		"version":   "system",
		"buy":       "trading",
		"sell":      "trading",
		"portfolio": "trading",
		"orders":    "trading",
		"history":   "trading",
		"balance":   "trading",
		"positions": "trading",
		"price":     "market",
		"chart":     "market",
		"news":      "market",
		"analysis":  "market",
		"alerts":    "market",
		"analyze":   "ai",
		"predict":   "ai",
		"sentiment": "ai",
		"chat":      "ai",
		"learn":     "ai",
		"wallet":    "web3",
		"connect":   "web3",
		"transfer":  "web3",
		"defi":      "web3",
	}

	if category, exists := categories[name]; exists {
		return category
	}

	return "misc"
}

func getCommandExamples(name string) []string {
	examples := map[string][]string{
		"status": {
			"status",
			"status --verbose",
		},
		"help": {
			"help",
			"help status",
			"help trading",
		},
		"price": {
			"price BTC",
			"price ETH USD",
			"price BTC --format json",
		},
		"buy": {
			"buy BTC 0.1",
			"buy ETH 1.5 --limit 2000",
		},
		"sell": {
			"sell BTC 0.05",
			"sell ETH 1.0 --market",
		},
		"portfolio": {
			"portfolio",
			"portfolio --detailed",
		},
		"analyze": {
			"analyze BTC",
			"analyze ETH --timeframe 1h",
		},
	}

	if exampleList, exists := examples[name]; exists {
		return exampleList
	}

	return []string{}
}
