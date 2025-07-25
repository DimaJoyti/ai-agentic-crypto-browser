//go:build load

package performance

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/time/rate"
)

// LoadTestSuite performs load testing on the application
type LoadTestSuite struct {
	suite.Suite

	baseURL    string
	httpClient *http.Client
	testUsers  []TestUser
	authTokens map[string]string
}

// TestUser represents a test user for load testing
type TestUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"-"`
}

// LoadTestConfig contains configuration for load testing
type LoadTestConfig struct {
	BaseURL           string
	ConcurrentUsers   int
	TestDuration      time.Duration
	RequestsPerSecond int
	RampUpTime        time.Duration
}

// PerformanceMetrics tracks performance metrics during load testing
type PerformanceMetrics struct {
	TotalRequests      int64
	SuccessfulRequests int64
	FailedRequests     int64
	TotalResponseTime  time.Duration
	MinResponseTime    time.Duration
	MaxResponseTime    time.Duration
	ResponseTimes      []time.Duration
	ErrorCounts        map[string]int64
	mutex              sync.RWMutex
}

// SetupSuite initializes the load test suite
func (suite *LoadTestSuite) SetupSuite() {
	suite.baseURL = "http://localhost:8080" // API Gateway URL
	suite.httpClient = &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	suite.authTokens = make(map[string]string)

	// Create test users
	suite.createTestUsers(10) // Create 10 test users for load testing (reduced for faster tests)
}

// createTestUsers creates multiple test users for load testing
func (suite *LoadTestSuite) createTestUsers(count int) {
	suite.testUsers = make([]TestUser, count)

	for i := 0; i < count; i++ {
		user := TestUser{
			Email:    fmt.Sprintf("loadtest%d@example.com", i),
			Password: "password123",
		}

		// Register user
		suite.registerUser(user)
		suite.testUsers[i] = user
	}
}

// registerUser registers a single test user
func (suite *LoadTestSuite) registerUser(user TestUser) {
	payload := map[string]interface{}{
		"email":      user.Email,
		"password":   user.Password,
		"first_name": "Load",
		"last_name":  "Test",
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", suite.baseURL+"/api/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(req)
	if err != nil {
		suite.T().Logf("Failed to register user %s: %v", user.Email, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusConflict {
		// User created or already exists, try to login
		suite.loginUser(user)
	}
}

// loginUser logs in a test user and stores the token
func (suite *LoadTestSuite) loginUser(user TestUser) {
	payload := map[string]interface{}{
		"email":    user.Email,
		"password": user.Password,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", suite.baseURL+"/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(req)
	if err != nil {
		suite.T().Logf("Failed to login user %s: %v", user.Email, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		if tokens, ok := response["tokens"].(map[string]interface{}); ok {
			if accessToken, ok := tokens["access_token"].(string); ok {
				suite.authTokens[user.Email] = accessToken
			}
		}
	}
}

// TestAuthenticationLoad tests authentication endpoint under load
func (suite *LoadTestSuite) TestAuthenticationLoad() {
	config := LoadTestConfig{
		BaseURL:           suite.baseURL,
		ConcurrentUsers:   20,
		TestDuration:      60 * time.Second,
		RequestsPerSecond: 50,
		RampUpTime:        10 * time.Second,
	}

	metrics := &PerformanceMetrics{
		ErrorCounts:     make(map[string]int64),
		MinResponseTime: time.Hour, // Initialize to high value
		ResponseTimes:   make([]time.Duration, 0),
	}

	suite.Run("Authentication endpoint load test", func() {
		suite.runLoadTest(config, metrics, suite.authenticationLoadTest)

		// Analyze results
		suite.analyzeResults(metrics, "Authentication")

		// Assert performance requirements
		avgResponseTime := metrics.TotalResponseTime / time.Duration(metrics.TotalRequests)
		assert.Less(suite.T(), avgResponseTime, 500*time.Millisecond, "Average response time should be under 500ms")

		successRate := float64(metrics.SuccessfulRequests) / float64(metrics.TotalRequests) * 100
		assert.Greater(suite.T(), successRate, 95.0, "Success rate should be above 95%")
	})
}

// TestWorkflowAPILoad tests workflow API endpoints under load
func (suite *LoadTestSuite) TestWorkflowAPILoad() {
	config := LoadTestConfig{
		BaseURL:           suite.baseURL,
		ConcurrentUsers:   15,
		TestDuration:      90 * time.Second,
		RequestsPerSecond: 30,
		RampUpTime:        15 * time.Second,
	}

	metrics := &PerformanceMetrics{
		ErrorCounts:     make(map[string]int64),
		MinResponseTime: time.Hour,
		ResponseTimes:   make([]time.Duration, 0),
	}

	suite.Run("Workflow API load test", func() {
		suite.runLoadTest(config, metrics, suite.workflowLoadTest)

		// Analyze results
		suite.analyzeResults(metrics, "Workflow API")

		// Assert performance requirements
		avgResponseTime := metrics.TotalResponseTime / time.Duration(metrics.TotalRequests)
		assert.Less(suite.T(), avgResponseTime, 1*time.Second, "Average response time should be under 1 second")

		successRate := float64(metrics.SuccessfulRequests) / float64(metrics.TotalRequests) * 100
		assert.Greater(suite.T(), successRate, 90.0, "Success rate should be above 90%")
	})
}

// TestBrowserServiceLoad tests browser service under load
func (suite *LoadTestSuite) TestBrowserServiceLoad() {
	config := LoadTestConfig{
		BaseURL:           suite.baseURL,
		ConcurrentUsers:   10, // Lower concurrency for browser operations
		TestDuration:      120 * time.Second,
		RequestsPerSecond: 10,
		RampUpTime:        20 * time.Second,
	}

	metrics := &PerformanceMetrics{
		ErrorCounts:     make(map[string]int64),
		MinResponseTime: time.Hour,
		ResponseTimes:   make([]time.Duration, 0),
	}

	suite.Run("Browser service load test", func() {
		suite.runLoadTest(config, metrics, suite.browserServiceLoadTest)

		// Analyze results
		suite.analyzeResults(metrics, "Browser Service")

		// Assert performance requirements
		avgResponseTime := metrics.TotalResponseTime / time.Duration(metrics.TotalRequests)
		assert.Less(suite.T(), avgResponseTime, 5*time.Second, "Average response time should be under 5 seconds")

		successRate := float64(metrics.SuccessfulRequests) / float64(metrics.TotalRequests) * 100
		assert.Greater(suite.T(), successRate, 85.0, "Success rate should be above 85%")
	})
}

// TestConcurrentWorkflowExecution tests concurrent workflow execution
func (suite *LoadTestSuite) TestConcurrentWorkflowExecution() {
	config := LoadTestConfig{
		BaseURL:           suite.baseURL,
		ConcurrentUsers:   25,
		TestDuration:      180 * time.Second,
		RequestsPerSecond: 5, // Lower rate for workflow execution
		RampUpTime:        30 * time.Second,
	}

	metrics := &PerformanceMetrics{
		ErrorCounts:     make(map[string]int64),
		MinResponseTime: time.Hour,
		ResponseTimes:   make([]time.Duration, 0),
	}

	suite.Run("Concurrent workflow execution test", func() {
		suite.runLoadTest(config, metrics, suite.workflowExecutionLoadTest)

		// Analyze results
		suite.analyzeResults(metrics, "Workflow Execution")

		// Assert performance requirements
		avgResponseTime := metrics.TotalResponseTime / time.Duration(metrics.TotalRequests)
		assert.Less(suite.T(), avgResponseTime, 10*time.Second, "Average response time should be under 10 seconds")

		successRate := float64(metrics.SuccessfulRequests) / float64(metrics.TotalRequests) * 100
		assert.Greater(suite.T(), successRate, 80.0, "Success rate should be above 80%")
	})
}

// runLoadTest executes a load test with the given configuration
func (suite *LoadTestSuite) runLoadTest(config LoadTestConfig, metrics *PerformanceMetrics, testFunc func(string) (time.Duration, error)) {
	ctx, cancel := context.WithTimeout(context.Background(), config.TestDuration+config.RampUpTime)
	defer cancel()

	// Rate limiter for requests per second
	limiter := rate.NewLimiter(rate.Limit(config.RequestsPerSecond), config.RequestsPerSecond)

	var wg sync.WaitGroup
	userChan := make(chan TestUser, len(suite.testUsers))

	// Fill user channel
	for _, user := range suite.testUsers {
		userChan <- user
	}
	close(userChan)

	// Start workers
	for i := 0; i < config.ConcurrentUsers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			// Ramp up delay
			rampUpDelay := time.Duration(workerID) * config.RampUpTime / time.Duration(config.ConcurrentUsers)
			time.Sleep(rampUpDelay)

			for {
				select {
				case <-ctx.Done():
					return
				default:
					// Wait for rate limiter
					if err := limiter.Wait(ctx); err != nil {
						return
					}

					// Get a user (cycle through users)
					userIndex := i % len(suite.testUsers)
					user := suite.testUsers[userIndex]

					// Get auth token
					token, exists := suite.authTokens[user.Email]
					if !exists {
						continue
					}

					// Execute test function
					duration, err := testFunc(token)

					// Record metrics
					suite.recordMetrics(metrics, duration, err)
				}
			}
		}(i)
	}

	wg.Wait()
}

// recordMetrics records performance metrics thread-safely
func (suite *LoadTestSuite) recordMetrics(metrics *PerformanceMetrics, duration time.Duration, err error) {
	metrics.mutex.Lock()
	defer metrics.mutex.Unlock()

	metrics.TotalRequests++
	metrics.ResponseTimes = append(metrics.ResponseTimes, duration)

	if err != nil {
		metrics.FailedRequests++
		metrics.ErrorCounts[err.Error()]++
	} else {
		metrics.SuccessfulRequests++
	}

	metrics.TotalResponseTime += duration

	if duration < metrics.MinResponseTime {
		metrics.MinResponseTime = duration
	}
	if duration > metrics.MaxResponseTime {
		metrics.MaxResponseTime = duration
	}
}

// Load test functions

// authenticationLoadTest performs authentication load testing
func (suite *LoadTestSuite) authenticationLoadTest(token string) (time.Duration, error) {
	start := time.Now()

	req, _ := http.NewRequest("GET", suite.baseURL+"/api/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := suite.httpClient.Do(req)
	if err != nil {
		return time.Since(start), err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return time.Since(start), fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return time.Since(start), nil
}

// workflowLoadTest performs workflow API load testing
func (suite *LoadTestSuite) workflowLoadTest(token string) (time.Duration, error) {
	start := time.Now()

	req, _ := http.NewRequest("GET", suite.baseURL+"/api/workflows", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := suite.httpClient.Do(req)
	if err != nil {
		return time.Since(start), err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return time.Since(start), fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return time.Since(start), nil
}

// browserServiceLoadTest performs browser service load testing
func (suite *LoadTestSuite) browserServiceLoadTest(token string) (time.Duration, error) {
	start := time.Now()

	// Create browser session
	payload := map[string]interface{}{
		"headless": true,
		"timeout":  30,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", suite.baseURL+"/api/browser/sessions", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(req)
	if err != nil {
		return time.Since(start), err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return time.Since(start), fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// Parse session ID and clean up
	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	if sessionID, ok := response["session_id"].(string); ok {
		// Clean up session
		deleteReq, _ := http.NewRequest("DELETE", suite.baseURL+"/api/browser/sessions/"+sessionID, nil)
		deleteReq.Header.Set("Authorization", "Bearer "+token)
		suite.httpClient.Do(deleteReq)
	}

	return time.Since(start), nil
}

// workflowExecutionLoadTest performs workflow execution load testing
func (suite *LoadTestSuite) workflowExecutionLoadTest(token string) (time.Duration, error) {
	start := time.Now()

	// Create a simple workflow
	workflow := map[string]interface{}{
		"name":        "Load Test Workflow",
		"description": "Simple workflow for load testing",
		"steps": []map[string]interface{}{
			{
				"type": "navigate",
				"url":  "https://httpbin.org/get",
			},
		},
	}

	body, _ := json.Marshal(workflow)
	req, _ := http.NewRequest("POST", suite.baseURL+"/api/workflows", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(req)
	if err != nil {
		return time.Since(start), err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return time.Since(start), fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return time.Since(start), nil
}

// analyzeResults analyzes and prints performance test results
func (suite *LoadTestSuite) analyzeResults(metrics *PerformanceMetrics, testName string) {
	metrics.mutex.RLock()
	defer metrics.mutex.RUnlock()

	if metrics.TotalRequests == 0 {
		suite.T().Logf("%s: No requests completed", testName)
		return
	}

	avgResponseTime := metrics.TotalResponseTime / time.Duration(metrics.TotalRequests)
	successRate := float64(metrics.SuccessfulRequests) / float64(metrics.TotalRequests) * 100

	// Calculate percentiles
	p95 := suite.calculatePercentile(metrics.ResponseTimes, 95)
	p99 := suite.calculatePercentile(metrics.ResponseTimes, 99)

	suite.T().Logf(`
%s Load Test Results:
- Total Requests: %d
- Successful Requests: %d
- Failed Requests: %d
- Success Rate: %.2f%%
- Average Response Time: %v
- Min Response Time: %v
- Max Response Time: %v
- 95th Percentile: %v
- 99th Percentile: %v
- Requests per Second: %.2f
`,
		testName,
		metrics.TotalRequests,
		metrics.SuccessfulRequests,
		metrics.FailedRequests,
		successRate,
		avgResponseTime,
		metrics.MinResponseTime,
		metrics.MaxResponseTime,
		p95,
		p99,
		float64(metrics.TotalRequests)/60.0, // Assuming 60 second test duration
	)

	// Log error details
	if len(metrics.ErrorCounts) > 0 {
		suite.T().Logf("Error breakdown:")
		for errorMsg, count := range metrics.ErrorCounts {
			suite.T().Logf("  %s: %d", errorMsg, count)
		}
	}
}

// calculatePercentile calculates the nth percentile of response times
func (suite *LoadTestSuite) calculatePercentile(times []time.Duration, percentile int) time.Duration {
	if len(times) == 0 {
		return 0
	}

	// Sort times (simple bubble sort for small datasets)
	sorted := make([]time.Duration, len(times))
	copy(sorted, times)

	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	index := (percentile * len(sorted)) / 100
	if index >= len(sorted) {
		index = len(sorted) - 1
	}

	return sorted[index]
}

// Run the load test suite
func TestLoadTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load tests in short mode")
	}

	suite.Run(t, new(LoadTestSuite))
}
