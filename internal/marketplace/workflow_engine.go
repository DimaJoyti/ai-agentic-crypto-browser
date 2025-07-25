package marketplace

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/ai-agentic-browser/internal/ai"
	"github.com/ai-agentic-browser/internal/browser"
	"github.com/ai-agentic-browser/internal/web3"
	"github.com/ai-agentic-browser/pkg/database"
	"github.com/ai-agentic-browser/pkg/observability"
)

// WorkflowEngine executes workflows and manages their lifecycle
type WorkflowEngine struct {
	db            *database.DB
	redis         *database.RedisClient
	aiService     *ai.Service
	browserService *browser.Service
	web3Service   *web3.Service
	logger        *observability.Logger
	
	// Execution management
	executions    map[uuid.UUID]*WorkflowExecution
	executionsMux sync.RWMutex
	
	// Step executors
	stepExecutors map[StepType]StepExecutor
}

// StepExecutor interface for different step types
type StepExecutor interface {
	Execute(ctx context.Context, step WorkflowStep, input map[string]interface{}) (map[string]interface{}, error)
	Validate(step WorkflowStep) error
}

// NewWorkflowEngine creates a new workflow engine
func NewWorkflowEngine(
	db *database.DB,
	redis *database.RedisClient,
	aiService *ai.Service,
	browserService *browser.Service,
	web3Service *web3.Service,
	logger *observability.Logger,
) *WorkflowEngine {
	engine := &WorkflowEngine{
		db:             db,
		redis:          redis,
		aiService:      aiService,
		browserService: browserService,
		web3Service:    web3Service,
		logger:         logger,
		executions:     make(map[uuid.UUID]*WorkflowExecution),
		stepExecutors:  make(map[StepType]StepExecutor),
	}

	// Register step executors
	engine.registerStepExecutors()

	return engine
}

func (we *WorkflowEngine) registerStepExecutors() {
	we.stepExecutors[StepTypeAI] = NewAIStepExecutor(we.aiService, we.logger)
	we.stepExecutors[StepTypeBrowser] = NewBrowserStepExecutor(we.browserService, we.logger)
	we.stepExecutors[StepTypeWeb3] = NewWeb3StepExecutor(we.web3Service, we.logger)
	we.stepExecutors[StepTypeAPI] = NewAPIStepExecutor(we.logger)
	we.stepExecutors[StepTypeData] = NewDataStepExecutor(we.logger)
	we.stepExecutors[StepTypeLogic] = NewLogicStepExecutor(we.logger)
	we.stepExecutors[StepTypeNotify] = NewNotifyStepExecutor(we.logger)
	we.stepExecutors[StepTypeWait] = NewWaitStepExecutor(we.logger)
}

// ExecuteWorkflow executes a workflow
func (we *WorkflowEngine) ExecuteWorkflow(ctx context.Context, req WorkflowExecuteRequest) (*WorkflowExecuteResponse, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("workflow-engine").Start(ctx, "workflow.Execute")
	defer span.End()

	// Get workflow definition
	workflow, err := we.getWorkflowDefinition(ctx, req.WorkflowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow: %w", err)
	}

	// Create execution record
	execution := &WorkflowExecution{
		ID:          uuid.New(),
		WorkflowID:  req.WorkflowID,
		Status:      ExecutionStatusPending,
		TriggerType: req.TriggerType,
		Input:       req.Input,
		StartedAt:   time.Now(),
		Steps:       make([]StepExecution, 0),
	}

	// Store execution
	we.executionsMux.Lock()
	we.executions[execution.ID] = execution
	we.executionsMux.Unlock()

	// Save to database
	if err := we.saveExecution(ctx, execution); err != nil {
		we.logger.Error(ctx, "Failed to save execution", err)
	}

	response := &WorkflowExecuteResponse{
		ExecutionID: execution.ID,
		Status:      execution.Status,
		Message:     "Workflow execution started",
	}

	// Execute asynchronously if requested
	if req.Async {
		go we.executeWorkflowAsync(context.Background(), execution, workflow)
		return response, nil
	}

	// Execute synchronously
	output, err := we.executeWorkflowSteps(ctx, execution, workflow)
	if err != nil {
		execution.Status = ExecutionStatusFailed
		execution.ErrorMessage = err.Error()
		response.Status = ExecutionStatusFailed
		response.Message = err.Error()
	} else {
		execution.Status = ExecutionStatusCompleted
		execution.Output = output
		response.Status = ExecutionStatusCompleted
		response.Output = output
		response.Message = "Workflow completed successfully"
	}

	// Update execution
	completedAt := time.Now()
	execution.CompletedAt = &completedAt
	execution.Duration = completedAt.Sub(execution.StartedAt).Milliseconds()

	// Save final state
	if err := we.saveExecution(ctx, execution); err != nil {
		we.logger.Error(ctx, "Failed to save final execution state", err)
	}

	return response, nil
}

// executeWorkflowAsync executes a workflow asynchronously
func (we *WorkflowEngine) executeWorkflowAsync(ctx context.Context, execution *WorkflowExecution, workflow *WorkflowDefinition) {
	execution.Status = ExecutionStatusRunning
	we.saveExecution(ctx, execution)

	output, err := we.executeWorkflowSteps(ctx, execution, workflow)
	
	completedAt := time.Now()
	execution.CompletedAt = &completedAt
	execution.Duration = completedAt.Sub(execution.StartedAt).Milliseconds()

	if err != nil {
		execution.Status = ExecutionStatusFailed
		execution.ErrorMessage = err.Error()
	} else {
		execution.Status = ExecutionStatusCompleted
		execution.Output = output
	}

	we.saveExecution(ctx, execution)

	// Send notification if configured
	we.sendExecutionNotification(ctx, execution, workflow)
}

// executeWorkflowSteps executes all steps in a workflow
func (we *WorkflowEngine) executeWorkflowSteps(ctx context.Context, execution *WorkflowExecution, workflow *WorkflowDefinition) (map[string]interface{}, error) {
	variables := make(map[string]interface{})
	
	// Initialize with input variables
	for k, v := range execution.Input {
		variables[k] = v
	}
	
	// Initialize with workflow variables
	for k, v := range workflow.Variables {
		variables[k] = v
	}

	// Execute steps in order
	for _, step := range workflow.Steps {
		stepExecution := StepExecution{
			StepID:    step.ID,
			Status:    ExecutionStatusRunning,
			StartedAt: time.Now(),
			Input:     we.prepareStepInput(step, variables),
		}

		// Check conditions
		if !we.evaluateStepConditions(step.Conditions, variables) {
			stepExecution.Status = ExecutionStatusCompleted
			stepExecution.Output = map[string]interface{}{"skipped": true}
			execution.Steps = append(execution.Steps, stepExecution)
			continue
		}

		// Execute step
		output, err := we.executeStep(ctx, step, stepExecution.Input)
		
		completedAt := time.Now()
		stepExecution.CompletedAt = &completedAt
		stepExecution.Duration = completedAt.Sub(stepExecution.StartedAt).Milliseconds()

		if err != nil {
			stepExecution.Status = ExecutionStatusFailed
			stepExecution.ErrorMessage = err.Error()
			execution.Steps = append(execution.Steps, stepExecution)

			// Handle step failure
			if len(step.OnFailure) > 0 {
				// Execute failure handlers
				for _, failureStepID := range step.OnFailure {
					if failureStep := we.findStepByID(workflow.Steps, failureStepID); failureStep != nil {
						we.executeStep(ctx, *failureStep, variables)
					}
				}
			}

			// Check if workflow should continue or fail
			if workflow.Settings.ErrorHandling == "stop" {
				return nil, fmt.Errorf("step %s failed: %w", step.ID, err)
			}
		} else {
			stepExecution.Status = ExecutionStatusCompleted
			stepExecution.Output = output

			// Update variables with step output
			for k, v := range output {
				variables[k] = v
			}

			// Execute success handlers
			if len(step.OnSuccess) > 0 {
				for _, successStepID := range step.OnSuccess {
					if successStep := we.findStepByID(workflow.Steps, successStepID); successStep != nil {
						we.executeStep(ctx, *successStep, variables)
					}
				}
			}
		}

		execution.Steps = append(execution.Steps, stepExecution)

		// Save intermediate state
		we.saveExecution(ctx, execution)
	}

	return variables, nil
}

// executeStep executes a single workflow step
func (we *WorkflowEngine) executeStep(ctx context.Context, step WorkflowStep, input map[string]interface{}) (map[string]interface{}, error) {
	executor, exists := we.stepExecutors[step.Type]
	if !exists {
		return nil, fmt.Errorf("no executor found for step type: %s", step.Type)
	}

	// Apply timeout if specified
	if step.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(step.Timeout)*time.Second)
		defer cancel()
	}

	// Execute with retries
	var lastErr error
	maxRetries := step.Retries
	if maxRetries == 0 {
		maxRetries = 1
	}

	for attempt := 0; attempt < maxRetries; attempt++ {
		output, err := executor.Execute(ctx, step, input)
		if err == nil {
			return output, nil
		}
		
		lastErr = err
		if attempt < maxRetries-1 {
			// Wait before retry
			time.Sleep(time.Duration(attempt+1) * time.Second)
		}
	}

	return nil, fmt.Errorf("step failed after %d attempts: %w", maxRetries, lastErr)
}

// Helper methods

func (we *WorkflowEngine) getWorkflowDefinition(ctx context.Context, workflowID uuid.UUID) (*WorkflowDefinition, error) {
	// This would fetch from database
	// For now, return a mock workflow
	return &WorkflowDefinition{
		Version: "1.0",
		Name:    "Sample Workflow",
		Steps:   []WorkflowStep{},
	}, nil
}

func (we *WorkflowEngine) saveExecution(ctx context.Context, execution *WorkflowExecution) error {
	// Save execution to database
	stepsJSON, _ := json.Marshal(execution.Steps)
	inputJSON, _ := json.Marshal(execution.Input)
	outputJSON, _ := json.Marshal(execution.Output)
	metadataJSON, _ := json.Marshal(execution.Metadata)

	query := `
		INSERT INTO workflow_executions (id, workflow_id, user_id, status, trigger_type, input, output, steps, started_at, completed_at, duration, error_message, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status,
			output = EXCLUDED.output,
			steps = EXCLUDED.steps,
			completed_at = EXCLUDED.completed_at,
			duration = EXCLUDED.duration,
			error_message = EXCLUDED.error_message,
			metadata = EXCLUDED.metadata
	`

	_, err := we.db.ExecContext(ctx, query,
		execution.ID, execution.WorkflowID, execution.UserID, execution.Status,
		execution.TriggerType, inputJSON, outputJSON, stepsJSON,
		execution.StartedAt, execution.CompletedAt, execution.Duration,
		execution.ErrorMessage, metadataJSON,
	)

	return err
}

func (we *WorkflowEngine) prepareStepInput(step WorkflowStep, variables map[string]interface{}) map[string]interface{} {
	input := make(map[string]interface{})
	
	// Copy step parameters
	for k, v := range step.Parameters {
		input[k] = we.interpolateValue(v, variables)
	}
	
	// Add global variables
	input["_variables"] = variables
	
	return input
}

func (we *WorkflowEngine) interpolateValue(value interface{}, variables map[string]interface{}) interface{} {
	switch v := value.(type) {
	case string:
		// Simple variable interpolation: ${variable_name}
		// In a real implementation, this would be more sophisticated
		return v
	default:
		return v
	}
}

func (we *WorkflowEngine) evaluateStepConditions(conditions []StepCondition, variables map[string]interface{}) bool {
	if len(conditions) == 0 {
		return true
	}

	for _, condition := range conditions {
		if !we.evaluateCondition(condition, variables) {
			return false
		}
	}

	return true
}

func (we *WorkflowEngine) evaluateCondition(condition StepCondition, variables map[string]interface{}) bool {
	fieldValue, exists := variables[condition.Field]
	if !exists {
		return false
	}

	switch condition.Operator {
	case "equals":
		return fieldValue == condition.Value
	case "not_equals":
		return fieldValue != condition.Value
	case "greater_than":
		if fv, ok := fieldValue.(float64); ok {
			if cv, ok := condition.Value.(float64); ok {
				return fv > cv
			}
		}
	case "less_than":
		if fv, ok := fieldValue.(float64); ok {
			if cv, ok := condition.Value.(float64); ok {
				return fv < cv
			}
		}
	case "contains":
		if fv, ok := fieldValue.(string); ok {
			if cv, ok := condition.Value.(string); ok {
				return len(fv) > 0 && len(cv) > 0
			}
		}
	}

	return false
}

func (we *WorkflowEngine) findStepByID(steps []WorkflowStep, stepID string) *WorkflowStep {
	for _, step := range steps {
		if step.ID == stepID {
			return &step
		}
	}
	return nil
}

func (we *WorkflowEngine) sendExecutionNotification(ctx context.Context, execution *WorkflowExecution, workflow *WorkflowDefinition) {
	// Send notifications based on workflow settings
	// This would integrate with email, Slack, webhooks, etc.
	we.logger.Info(ctx, "Workflow execution completed", map[string]interface{}{
		"execution_id": execution.ID.String(),
		"status":       execution.Status,
		"duration":     execution.Duration,
	})
}

// GetExecution retrieves a workflow execution
func (we *WorkflowEngine) GetExecution(ctx context.Context, executionID uuid.UUID) (*WorkflowExecution, error) {
	we.executionsMux.RLock()
	execution, exists := we.executions[executionID]
	we.executionsMux.RUnlock()

	if exists {
		return execution, nil
	}

	// Load from database
	return we.loadExecutionFromDB(ctx, executionID)
}

func (we *WorkflowEngine) loadExecutionFromDB(ctx context.Context, executionID uuid.UUID) (*WorkflowExecution, error) {
	query := `
		SELECT id, workflow_id, user_id, status, trigger_type, input, output, steps,
			   started_at, completed_at, duration, error_message, metadata
		FROM workflow_executions WHERE id = $1
	`

	execution := &WorkflowExecution{}
	var inputJSON, outputJSON, stepsJSON, metadataJSON []byte

	err := we.db.QueryRowContext(ctx, query, executionID).Scan(
		&execution.ID, &execution.WorkflowID, &execution.UserID, &execution.Status,
		&execution.TriggerType, &inputJSON, &outputJSON, &stepsJSON,
		&execution.StartedAt, &execution.CompletedAt, &execution.Duration,
		&execution.ErrorMessage, &metadataJSON,
	)

	if err != nil {
		return nil, err
	}

	// Unmarshal JSON fields
	json.Unmarshal(inputJSON, &execution.Input)
	json.Unmarshal(outputJSON, &execution.Output)
	json.Unmarshal(stepsJSON, &execution.Steps)
	json.Unmarshal(metadataJSON, &execution.Metadata)

	return execution, nil
}

// CancelExecution cancels a running workflow execution
func (we *WorkflowEngine) CancelExecution(ctx context.Context, executionID uuid.UUID) error {
	we.executionsMux.Lock()
	defer we.executionsMux.Unlock()

	execution, exists := we.executions[executionID]
	if !exists {
		return fmt.Errorf("execution not found: %s", executionID)
	}

	if execution.Status != ExecutionStatusRunning {
		return fmt.Errorf("execution is not running: %s", execution.Status)
	}

	execution.Status = ExecutionStatusCancelled
	completedAt := time.Now()
	execution.CompletedAt = &completedAt
	execution.Duration = completedAt.Sub(execution.StartedAt).Milliseconds()

	return we.saveExecution(ctx, execution)
}
