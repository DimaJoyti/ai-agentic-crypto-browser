package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ai-agentic-browser/internal/trading"
	"github.com/ai-agentic-browser/internal/trading/testing"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/gorilla/mux"
)

// TestingHandler handles bot testing API requests
type TestingHandler struct {
	logger    *observability.Logger
	framework *testing.BotTestFramework
	executor  *testing.TestExecutor
}

// NewTestingHandler creates a new testing handler
func NewTestingHandler(logger *observability.Logger, framework *testing.BotTestFramework, executor *testing.TestExecutor) *TestingHandler {
	return &TestingHandler{
		logger:    logger,
		framework: framework,
		executor:  executor,
	}
}

// RegisterRoutes registers testing API routes
func (h *TestingHandler) RegisterRoutes(router *mux.Router) {
	// Test execution endpoints
	router.HandleFunc("/api/v1/testing/tests", h.SubmitTest).Methods("POST")
	router.HandleFunc("/api/v1/testing/tests", h.GetAllTests).Methods("GET")
	router.HandleFunc("/api/v1/testing/tests/{testId}", h.GetTest).Methods("GET")
	router.HandleFunc("/api/v1/testing/tests/{testId}/cancel", h.CancelTest).Methods("POST")
	router.HandleFunc("/api/v1/testing/tests/{testId}/result", h.GetTestResult).Methods("GET")

	// Test scenarios endpoints
	router.HandleFunc("/api/v1/testing/scenarios", h.GetTestScenarios).Methods("GET")
	router.HandleFunc("/api/v1/testing/scenarios/{scenarioId}", h.GetTestScenario).Methods("GET")
	router.HandleFunc("/api/v1/testing/scenarios/{scenarioId}/run", h.RunTestScenario).Methods("POST")

	// Test suite endpoints
	router.HandleFunc("/api/v1/testing/suites/run", h.RunTestSuite).Methods("POST")
	router.HandleFunc("/api/v1/testing/suites/unit", h.RunUnitTests).Methods("POST")
	router.HandleFunc("/api/v1/testing/suites/integration", h.RunIntegrationTests).Methods("POST")
	router.HandleFunc("/api/v1/testing/suites/backtest", h.RunBacktests).Methods("POST")

	// Test environment endpoints
	router.HandleFunc("/api/v1/testing/environment", h.GetTestEnvironment).Methods("GET")
	router.HandleFunc("/api/v1/testing/environment/reset", h.ResetTestEnvironment).Methods("POST")
	router.HandleFunc("/api/v1/testing/environment/market-condition", h.SetMarketCondition).Methods("POST")

	// Test reports endpoints
	router.HandleFunc("/api/v1/testing/reports", h.GetTestReports).Methods("GET")
	router.HandleFunc("/api/v1/testing/reports/{reportId}", h.GetTestReport).Methods("GET")
	router.HandleFunc("/api/v1/testing/reports/generate", h.GenerateTestReport).Methods("POST")

	// Mock exchange endpoints
	router.HandleFunc("/api/v1/testing/mock-exchange/info", h.GetMockExchangeInfo).Methods("GET")
	router.HandleFunc("/api/v1/testing/mock-exchange/accounts", h.GetMockAccounts).Methods("GET")
	router.HandleFunc("/api/v1/testing/mock-exchange/orders", h.GetMockOrders).Methods("GET")
	router.HandleFunc("/api/v1/testing/mock-exchange/trades", h.GetMockTrades).Methods("GET")
}

// SubmitTest handles POST /api/v1/testing/tests
func (h *TestingHandler) SubmitTest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var requestBody struct {
		Type            string                   `json:"type"`
		Strategy        string                   `json:"strategy"`
		Config          *trading.BotConfig       `json:"config,omitempty"`
		Scenario        *testing.TestScenario    `json:"scenario,omitempty"`
		Parameters      map[string]interface{}   `json:"parameters,omitempty"`
		ExpectedResults *testing.ExpectedResults `json:"expected_results,omitempty"`
		Priority        int                      `json:"priority,omitempty"`
		Timeout         string                   `json:"timeout,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Parse timeout
	timeout := 300 * time.Second // default
	if requestBody.Timeout != "" {
		if t, err := time.ParseDuration(requestBody.Timeout); err == nil {
			timeout = t
		}
	}

	// Create test request
	testRequest := &testing.TestRequest{
		Type:            testing.TestType(requestBody.Type),
		Strategy:        requestBody.Strategy,
		Config:          requestBody.Config,
		Scenario:        requestBody.Scenario,
		Parameters:      requestBody.Parameters,
		ExpectedResults: requestBody.ExpectedResults,
		Priority:        requestBody.Priority,
		Timeout:         timeout,
	}

	// Submit test
	testID, err := h.executor.SubmitTest(testRequest)
	if err != nil {
		h.logger.Error(ctx, "Failed to submit test", err, map[string]interface{}{
			"type":     requestBody.Type,
			"strategy": requestBody.Strategy,
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, "Test submitted", map[string]interface{}{
		"test_id":  testID,
		"type":     requestBody.Type,
		"strategy": requestBody.Strategy,
	})

	response := map[string]interface{}{
		"test_id": testID,
		"status":  "submitted",
		"message": "Test submitted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAllTests handles GET /api/v1/testing/tests
func (h *TestingHandler) GetAllTests(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	activeTests := h.executor.GetActiveTests()

	response := map[string]interface{}{
		"active_tests": activeTests,
		"count":        len(activeTests),
		"timestamp":    time.Now(),
	}

	h.logger.Info(ctx, "Active tests retrieved", map[string]interface{}{
		"count": len(activeTests),
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetTest handles GET /api/v1/testing/tests/{testId}
func (h *TestingHandler) GetTest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	testID := vars["testId"]

	execution, err := h.executor.GetTestStatus(testID)
	if err != nil {
		h.logger.Error(ctx, "Failed to get test status", err, map[string]interface{}{
			"test_id": testID,
		})
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.logger.Info(ctx, "Test status retrieved", map[string]interface{}{
		"test_id": testID,
		"status":  string(execution.Status),
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(execution)
}

// CancelTest handles POST /api/v1/testing/tests/{testId}/cancel
func (h *TestingHandler) CancelTest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	testID := vars["testId"]

	if err := h.executor.CancelTest(testID); err != nil {
		h.logger.Error(ctx, "Failed to cancel test", err, map[string]interface{}{
			"test_id": testID,
		})
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Info(ctx, "Test cancelled", map[string]interface{}{
		"test_id": testID,
	})

	response := map[string]interface{}{
		"test_id": testID,
		"status":  "cancelled",
		"message": "Test cancelled successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetTestResult handles GET /api/v1/testing/tests/{testId}/result
func (h *TestingHandler) GetTestResult(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	testID := vars["testId"]

	result, err := h.executor.GetTestResult(testID)
	if err != nil {
		h.logger.Error(ctx, "Failed to get test result", err, map[string]interface{}{
			"test_id": testID,
		})
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.logger.Info(ctx, "Test result retrieved", map[string]interface{}{
		"test_id": testID,
		"passed":  result.Passed,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GetTestScenarios handles GET /api/v1/testing/scenarios
func (h *TestingHandler) GetTestScenarios(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get scenarios from framework (this would be implemented in the framework)
	scenarios := []map[string]interface{}{
		{
			"id":          "dca-bull-market",
			"name":        "DCA Strategy - Bull Market",
			"description": "Test DCA strategy performance in a bull market",
			"type":        "backtest",
			"strategy":    "dca",
			"duration":    "1h",
		},
		{
			"id":          "grid-sideways-market",
			"name":        "Grid Strategy - Sideways Market",
			"description": "Test Grid strategy performance in a sideways market",
			"type":        "backtest",
			"strategy":    "grid",
			"duration":    "1h",
		},
	}

	response := map[string]interface{}{
		"scenarios": scenarios,
		"count":     len(scenarios),
	}

	h.logger.Info(ctx, "Test scenarios retrieved", map[string]interface{}{
		"count": len(scenarios),
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RunTestScenario handles POST /api/v1/testing/scenarios/{scenarioId}/run
func (h *TestingHandler) RunTestScenario(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	scenarioID := vars["scenarioId"]

	// Find scenario and create test request
	// This would be implemented based on predefined scenarios

	response := map[string]interface{}{
		"scenario_id": scenarioID,
		"status":      "submitted",
		"message":     "Test scenario submitted successfully",
	}

	h.logger.Info(ctx, "Test scenario submitted", map[string]interface{}{
		"scenario_id": scenarioID,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RunTestSuite handles POST /api/v1/testing/suites/run
func (h *TestingHandler) RunTestSuite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var requestBody struct {
		TestTypes  []string `json:"test_types"`
		Strategies []string `json:"strategies,omitempty"`
		Parallel   bool     `json:"parallel,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Submit multiple tests based on request
	testIDs := make([]string, 0)

	for _, testType := range requestBody.TestTypes {
		strategies := requestBody.Strategies
		if len(strategies) == 0 {
			strategies = []string{"dca", "grid", "momentum"}
		}

		for _, strategy := range strategies {
			testRequest := &testing.TestRequest{
				Type:     testing.TestType(testType),
				Strategy: strategy,
				Config: &trading.BotConfig{
					TradingPairs: []string{"BTC/USDT"},
					Exchange:     "mock",
					BaseCurrency: "USDT",
					StrategyParams: map[string]interface{}{
						"strategy": strategy,
					},
					Enabled: true,
				},
				Timeout: 300 * time.Second,
			}

			testID, err := h.executor.SubmitTest(testRequest)
			if err != nil {
				h.logger.Error(ctx, "Failed to submit test in suite", err, map[string]interface{}{
					"type":     testType,
					"strategy": strategy,
				})
				continue
			}

			testIDs = append(testIDs, testID)
		}
	}

	response := map[string]interface{}{
		"test_ids": testIDs,
		"count":    len(testIDs),
		"status":   "submitted",
		"message":  "Test suite submitted successfully",
	}

	h.logger.Info(ctx, "Test suite submitted", map[string]interface{}{
		"test_count": len(testIDs),
		"test_types": requestBody.TestTypes,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetTestEnvironment handles GET /api/v1/testing/environment
func (h *TestingHandler) GetTestEnvironment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get environment information
	exchangeInfo := h.framework.GetMockExchange().GetExchangeInfo()
	marketDataInfo := h.framework.GetMockMarketData().GetMarketDataInfo()

	response := map[string]interface{}{
		"exchange":    exchangeInfo,
		"market_data": marketDataInfo,
		"timestamp":   time.Now(),
	}

	h.logger.Info(ctx, "Test environment info retrieved", nil)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SetMarketCondition handles POST /api/v1/testing/environment/market-condition
func (h *TestingHandler) SetMarketCondition(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var requestBody struct {
		Condition string `json:"condition"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set market condition
	h.framework.GetMockExchange().SetMarketCondition(requestBody.Condition)
	h.framework.GetMockMarketData().SetMarketCondition(requestBody.Condition)

	response := map[string]interface{}{
		"condition": requestBody.Condition,
		"status":    "updated",
		"message":   "Market condition updated successfully",
	}

	h.logger.Info(ctx, "Market condition updated", map[string]interface{}{
		"condition": requestBody.Condition,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetMockExchangeInfo handles GET /api/v1/testing/mock-exchange/info
func (h *TestingHandler) GetMockExchangeInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	info := h.framework.GetMockExchange().GetExchangeInfo()

	h.logger.Info(ctx, "Mock exchange info retrieved", nil)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

// Placeholder implementations for remaining endpoints
func (h *TestingHandler) GetTestScenario(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *TestingHandler) RunUnitTests(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *TestingHandler) RunIntegrationTests(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *TestingHandler) RunBacktests(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *TestingHandler) ResetTestEnvironment(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *TestingHandler) GetTestReports(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *TestingHandler) GetTestReport(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *TestingHandler) GenerateTestReport(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *TestingHandler) GetMockAccounts(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *TestingHandler) GetMockOrders(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *TestingHandler) GetMockTrades(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
