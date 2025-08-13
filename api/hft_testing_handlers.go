package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ai-agentic-browser/internal/hft"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// HFTTestingHandlers provides HTTP handlers for HFT Testing & Simulation Framework
type HFTTestingHandlers struct {
	framework *hft.TestingSimulationFramework
	logger    *observability.Logger
}

// NewHFTTestingHandlers creates new HFT testing framework HTTP handlers
func NewHFTTestingHandlers(framework *hft.TestingSimulationFramework, logger *observability.Logger) *HFTTestingHandlers {
	return &HFTTestingHandlers{
		framework: framework,
		logger:    logger,
	}
}

// RunTest handles test execution requests
func (h *HFTTestingHandlers) RunTest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name            string         `json:"name"`
		Description     string         `json:"description"`
		Type            string         `json:"type"`
		EnvironmentID   string         `json:"environment_id"`
		Steps           []hft.TestStep `json:"steps"`
		Preconditions   []string       `json:"preconditions"`
		Postconditions  []string       `json:"postconditions"`
		ExpectedResults []string       `json:"expected_results"`
		Tags            []string       `json:"tags"`
		Priority        int            `json:"priority"`
		Timeout         int            `json:"timeout_seconds"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(ctx, "Failed to decode test request", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Type == "" || req.EnvironmentID == "" {
		http.Error(w, "Name, type, and environment_id are required", http.StatusBadRequest)
		return
	}

	// Create test scenario
	scenario := &hft.TestScenario{
		ID:              uuid.New(),
		Name:            req.Name,
		Description:     req.Description,
		Type:            hft.TestType(req.Type),
		Steps:           req.Steps,
		Preconditions:   req.Preconditions,
		Postconditions:  req.Postconditions,
		ExpectedResults: req.ExpectedResults,
		Tags:            req.Tags,
		Priority:        req.Priority,
	}

	// Execute test
	result, err := h.framework.RunTest(ctx, scenario, req.EnvironmentID)
	if err != nil {
		h.logger.Error(ctx, "Failed to run test", err)
		http.Error(w, "Failed to run test: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Test executed successfully",
		"result":  result,
	})

	h.logger.Info(ctx, "Test executed via API", map[string]interface{}{
		"test_id":   result.ID.String(),
		"test_name": req.Name,
		"success":   result.Success,
		"duration":  result.Duration.String(),
	})
}

// GetTestResults handles test results requests
func (h *HFTTestingHandlers) GetTestResults(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	testTypeStr := r.URL.Query().Get("type")

	limit := 50 // Default limit
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	var testType hft.TestType
	if testTypeStr != "" {
		testType = hft.TestType(testTypeStr)
	}

	results := h.framework.GetTestResults(limit, testType)

	response := map[string]interface{}{
		"results": results,
		"count":   len(results),
		"filters": map[string]interface{}{
			"limit": limit,
			"type":  testTypeStr,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error(ctx, "Failed to encode test results", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// RunSimulation handles simulation execution requests
func (h *HFTTestingHandlers) RunSimulation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Type       string                 `json:"type"`
		Parameters map[string]interface{} `json:"parameters"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(ctx, "Failed to decode simulation request", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Type == "" {
		http.Error(w, "Type is required", http.StatusBadRequest)
		return
	}

	if req.Parameters == nil {
		req.Parameters = make(map[string]interface{})
	}

	// Execute simulation
	session, err := h.framework.RunSimulation(ctx, hft.SimulationType(req.Type), req.Parameters)
	if err != nil {
		h.logger.Error(ctx, "Failed to run simulation", err)
		http.Error(w, "Failed to run simulation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "success",
		"message":    "Simulation started successfully",
		"session_id": session.ID.String(),
		"session":    session,
	})

	h.logger.Info(ctx, "Simulation started via API", map[string]interface{}{
		"session_id": session.ID.String(),
		"type":       req.Type,
		"parameters": len(req.Parameters),
	})
}

// GetSimulations handles simulation list requests
func (h *HFTTestingHandlers) GetSimulations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	simTypeStr := r.URL.Query().Get("type")

	limit := 50 // Default limit
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	var simType hft.SimulationType
	if simTypeStr != "" {
		simType = hft.SimulationType(simTypeStr)
	}

	results := h.framework.GetSimulationResults(limit, simType)

	response := map[string]interface{}{
		"results": results,
		"count":   len(results),
		"filters": map[string]interface{}{
			"limit": limit,
			"type":  simTypeStr,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error(ctx, "Failed to encode simulation results", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetSimulation handles single simulation requests
func (h *HFTTestingHandlers) GetSimulation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	sessionID := vars["id"]
	if sessionID == "" {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	session := h.framework.GetSimulationSession(sessionID)
	if session == nil {
		http.Error(w, "Simulation session not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(session); err != nil {
		h.logger.Error(ctx, "Failed to encode simulation session", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// CreateEnvironment handles test environment creation requests
func (h *HFTTestingHandlers) CreateEnvironment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name   string                 `json:"name"`
		Type   string                 `json:"type"`
		Config map[string]interface{} `json:"config"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(ctx, "Failed to decode environment request", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Type == "" {
		http.Error(w, "Name and type are required", http.StatusBadRequest)
		return
	}

	if req.Config == nil {
		req.Config = make(map[string]interface{})
	}

	// Create environment
	environment, err := h.framework.CreateTestEnvironment(ctx, req.Name, hft.EnvironmentType(req.Type), req.Config)
	if err != nil {
		h.logger.Error(ctx, "Failed to create test environment", err)
		http.Error(w, "Failed to create environment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":         "success",
		"message":        "Test environment created successfully",
		"environment_id": environment.ID.String(),
		"environment":    environment,
	})

	h.logger.Info(ctx, "Test environment created via API", map[string]interface{}{
		"environment_id": environment.ID.String(),
		"name":           req.Name,
		"type":           req.Type,
	})
}

// GetEnvironments handles test environments list requests
func (h *HFTTestingHandlers) GetEnvironments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	environments := h.framework.GetTestEnvironments()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"environments": environments,
		"count":        len(environments),
	}); err != nil {
		h.logger.Error(ctx, "Failed to encode test environments", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetEnvironment handles single environment requests
func (h *HFTTestingHandlers) GetEnvironment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	environmentID := vars["id"]
	if environmentID == "" {
		http.Error(w, "Environment ID is required", http.StatusBadRequest)
		return
	}

	environment := h.framework.GetTestEnvironment(environmentID)
	if environment == nil {
		http.Error(w, "Test environment not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(environment); err != nil {
		h.logger.Error(ctx, "Failed to encode test environment", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetStatus handles framework status requests
func (h *HFTTestingHandlers) GetStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status := h.framework.GetFrameworkStatus()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		h.logger.Error(ctx, "Failed to encode framework status", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
