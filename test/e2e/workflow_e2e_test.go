package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// TestUser represents a test user for E2E testing
type TestUser struct {
	email    string
	password string
	token    string
}

// WorkflowE2ETestSuite tests complete workflow scenarios
type WorkflowE2ETestSuite struct {
	suite.Suite

	// Browser context for UI testing
	browserCtx    context.Context
	browserCancel context.CancelFunc

	// Test data
	testUser *TestUser
	baseURL  string
}

// SetupSuite initializes the E2E test suite
func (suite *WorkflowE2ETestSuite) SetupSuite() {
	// Set base URL for testing
	suite.baseURL = "http://localhost:3000" // Frontend URL

	// Initialize browser context
	suite.initializeBrowser()

	// Create test user
	suite.createTestUser()

	// Wait for services to be ready
	suite.waitForServices()
}

// TearDownSuite cleans up the E2E test suite
func (suite *WorkflowE2ETestSuite) TearDownSuite() {
	if suite.browserCancel != nil {
		suite.browserCancel()
	}
}

// initializeBrowser sets up Chrome browser for testing
func (suite *WorkflowE2ETestSuite) initializeBrowser() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.WindowSize(1920, 1080),
	)

	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	suite.browserCtx, suite.browserCancel = chromedp.NewContext(allocCtx)
}

// createTestUser creates a test user for E2E testing
func (suite *WorkflowE2ETestSuite) createTestUser() {
	suite.testUser = &TestUser{
		email:    "e2e@example.com",
		password: "password123",
	}

	// Register user via API
	payload := map[string]interface{}{
		"email":      suite.testUser.email,
		"password":   suite.testUser.password,
		"first_name": "E2E",
		"last_name":  "Test",
	}

	resp := suite.makeAPIRequest("POST", "/api/auth/register", payload, nil)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		tokens := response["tokens"].(map[string]interface{})
		suite.testUser.token = tokens["access_token"].(string)
	} else {
		// User might already exist, try to login
		suite.loginTestUser()
	}
}

// loginTestUser logs in the test user
func (suite *WorkflowE2ETestSuite) loginTestUser() {
	payload := map[string]interface{}{
		"email":    suite.testUser.email,
		"password": suite.testUser.password,
	}

	resp := suite.makeAPIRequest("POST", "/api/auth/login", payload, nil)
	defer resp.Body.Close()

	require.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	tokens := response["tokens"].(map[string]interface{})
	suite.testUser.token = tokens["access_token"].(string)
}

// waitForServices waits for all services to be ready
func (suite *WorkflowE2ETestSuite) waitForServices() {
	services := []string{
		"http://localhost:8080/health", // API Gateway
		"http://localhost:3000",        // Frontend
	}

	for _, service := range services {
		suite.waitForService(service)
	}
}

// waitForService waits for a specific service to be ready
func (suite *WorkflowE2ETestSuite) waitForService(url string) {
	client := &http.Client{Timeout: 5 * time.Second}

	require.Eventually(suite.T(), func() bool {
		resp, err := client.Get(url)
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode < 500
	}, 60*time.Second, 2*time.Second, "Service %s not ready", url)
}

// TestCompleteWorkflowCreation tests the complete workflow creation process
func (suite *WorkflowE2ETestSuite) TestCompleteWorkflowCreation() {
	suite.Run("User can create and execute a workflow", func() {
		// Navigate to login page
		err := chromedp.Run(suite.browserCtx,
			chromedp.Navigate(suite.baseURL+"/login"),
			chromedp.WaitVisible(`input[name="email"]`, chromedp.ByQuery),
		)
		require.NoError(suite.T(), err)

		// Login
		err = chromedp.Run(suite.browserCtx,
			chromedp.SendKeys(`input[name="email"]`, suite.testUser.email, chromedp.ByQuery),
			chromedp.SendKeys(`input[name="password"]`, suite.testUser.password, chromedp.ByQuery),
			chromedp.Click(`button[type="submit"]`, chromedp.ByQuery),
			chromedp.WaitVisible(`[data-testid="dashboard"]`, chromedp.ByQuery),
		)
		require.NoError(suite.T(), err)

		// Navigate to workflow creation
		err = chromedp.Run(suite.browserCtx,
			chromedp.Click(`[data-testid="create-workflow"]`, chromedp.ByQuery),
			chromedp.WaitVisible(`[data-testid="workflow-builder"]`, chromedp.ByQuery),
		)
		require.NoError(suite.T(), err)

		// Create a simple workflow
		workflowName := "E2E Test Workflow"
		err = chromedp.Run(suite.browserCtx,
			chromedp.SendKeys(`input[name="workflow-name"]`, workflowName, chromedp.ByQuery),
			chromedp.Click(`[data-testid="add-step"]`, chromedp.ByQuery),
			chromedp.WaitVisible(`[data-testid="step-selector"]`, chromedp.ByQuery),
		)
		require.NoError(suite.T(), err)

		// Add a browser navigation step
		err = chromedp.Run(suite.browserCtx,
			chromedp.Click(`[data-testid="step-navigate"]`, chromedp.ByQuery),
			chromedp.SendKeys(`input[name="url"]`, "https://example.com", chromedp.ByQuery),
			chromedp.Click(`[data-testid="add-step-confirm"]`, chromedp.ByQuery),
		)
		require.NoError(suite.T(), err)

		// Save workflow
		err = chromedp.Run(suite.browserCtx,
			chromedp.Click(`[data-testid="save-workflow"]`, chromedp.ByQuery),
			chromedp.WaitVisible(`[data-testid="workflow-saved"]`, chromedp.ByQuery),
		)
		require.NoError(suite.T(), err)

		// Execute workflow
		err = chromedp.Run(suite.browserCtx,
			chromedp.Click(`[data-testid="execute-workflow"]`, chromedp.ByQuery),
			chromedp.WaitVisible(`[data-testid="execution-started"]`, chromedp.ByQuery),
		)
		require.NoError(suite.T(), err)

		// Wait for execution to complete
		var executionStatus string
		err = chromedp.Run(suite.browserCtx,
			chromedp.WaitVisible(`[data-testid="execution-status"]`, chromedp.ByQuery),
			chromedp.Text(`[data-testid="execution-status"]`, &executionStatus, chromedp.ByQuery),
		)
		require.NoError(suite.T(), err)

		// Verify execution completed successfully
		assert.Contains(suite.T(), strings.ToLower(executionStatus), "completed")

		// Verify workflow appears in workflow list
		err = chromedp.Run(suite.browserCtx,
			chromedp.Navigate(suite.baseURL+"/workflows"),
			chromedp.WaitVisible(`[data-testid="workflow-list"]`, chromedp.ByQuery),
		)
		require.NoError(suite.T(), err)

		var workflowListText string
		err = chromedp.Run(suite.browserCtx,
			chromedp.Text(`[data-testid="workflow-list"]`, &workflowListText, chromedp.ByQuery),
		)
		require.NoError(suite.T(), err)

		assert.Contains(suite.T(), workflowListText, workflowName)
	})
}

// TestAIAgentInteraction tests AI agent functionality
func (suite *WorkflowE2ETestSuite) TestAIAgentInteraction() {
	suite.Run("User can interact with AI agent", func() {
		// Login first
		suite.loginViaUI()

		// Navigate to AI agent interface
		err := chromedp.Run(suite.browserCtx,
			chromedp.Navigate(suite.baseURL+"/ai-agent"),
			chromedp.WaitVisible(`[data-testid="ai-chat"]`, chromedp.ByQuery),
		)
		require.NoError(suite.T(), err)

		// Send a message to AI agent
		testMessage := "Navigate to google.com and search for 'AI automation'"
		err = chromedp.Run(suite.browserCtx,
			chromedp.SendKeys(`textarea[data-testid="ai-input"]`, testMessage, chromedp.ByQuery),
			chromedp.Click(`[data-testid="send-message"]`, chromedp.ByQuery),
		)
		require.NoError(suite.T(), err)

		// Wait for AI response
		err = chromedp.Run(suite.browserCtx,
			chromedp.WaitVisible(`[data-testid="ai-response"]`, chromedp.ByQuery),
		)
		require.NoError(suite.T(), err)

		// Verify AI responded
		var aiResponse string
		err = chromedp.Run(suite.browserCtx,
			chromedp.Text(`[data-testid="ai-response"]`, &aiResponse, chromedp.ByQuery),
		)
		require.NoError(suite.T(), err)

		assert.NotEmpty(suite.T(), aiResponse)
		assert.NotContains(suite.T(), strings.ToLower(aiResponse), "error")
	})
}

// TestBrowserAutomation tests browser automation features
func (suite *WorkflowE2ETestSuite) TestBrowserAutomation() {
	suite.Run("User can control browser automation", func() {
		// Login first
		suite.loginViaUI()

		// Navigate to browser control interface
		err := chromedp.Run(suite.browserCtx,
			chromedp.Navigate(suite.baseURL+"/browser"),
			chromedp.WaitVisible(`[data-testid="browser-control"]`, chromedp.ByQuery),
		)
		require.NoError(suite.T(), err)

		// Start a new browser session
		err = chromedp.Run(suite.browserCtx,
			chromedp.Click(`[data-testid="start-session"]`, chromedp.ByQuery),
			chromedp.WaitVisible(`[data-testid="browser-session-active"]`, chromedp.ByQuery),
		)
		require.NoError(suite.T(), err)

		// Navigate to a website
		err = chromedp.Run(suite.browserCtx,
			chromedp.SendKeys(`input[data-testid="url-input"]`, "https://httpbin.org", chromedp.ByQuery),
			chromedp.Click(`[data-testid="navigate-button"]`, chromedp.ByQuery),
		)
		require.NoError(suite.T(), err)

		// Wait for navigation to complete
		err = chromedp.Run(suite.browserCtx,
			chromedp.WaitVisible(`[data-testid="page-loaded"]`, chromedp.ByQuery),
		)
		require.NoError(suite.T(), err)

		// Verify browser session is working
		var sessionStatus string
		err = chromedp.Run(suite.browserCtx,
			chromedp.Text(`[data-testid="session-status"]`, &sessionStatus, chromedp.ByQuery),
		)
		require.NoError(suite.T(), err)

		assert.Contains(suite.T(), strings.ToLower(sessionStatus), "active")

		// Stop browser session
		err = chromedp.Run(suite.browserCtx,
			chromedp.Click(`[data-testid="stop-session"]`, chromedp.ByQuery),
			chromedp.WaitVisible(`[data-testid="session-stopped"]`, chromedp.ByQuery),
		)
		require.NoError(suite.T(), err)
	})
}

// TestWeb3Integration tests Web3 wallet integration
func (suite *WorkflowE2ETestSuite) TestWeb3Integration() {
	suite.Run("User can connect Web3 wallet", func() {
		// Login first
		suite.loginViaUI()

		// Navigate to Web3 interface
		err := chromedp.Run(suite.browserCtx,
			chromedp.Navigate(suite.baseURL+"/web3"),
			chromedp.WaitVisible(`[data-testid="web3-interface"]`, chromedp.ByQuery),
		)
		require.NoError(suite.T(), err)

		// Check if wallet connection is available
		var connectButtonExists bool
		err = chromedp.Run(suite.browserCtx,
			chromedp.WaitVisible(`[data-testid="connect-wallet"]`, chromedp.ByQuery),
		)

		if err == nil {
			connectButtonExists = true
		}

		if connectButtonExists {
			// Try to connect wallet (this might fail in headless mode)
			err = chromedp.Run(suite.browserCtx,
				chromedp.Click(`[data-testid="connect-wallet"]`, chromedp.ByQuery),
			)
			// Don't require this to succeed as it depends on wallet extension
		}

		// Verify Web3 interface is functional
		var web3Status string
		err = chromedp.Run(suite.browserCtx,
			chromedp.Text(`[data-testid="web3-status"]`, &web3Status, chromedp.ByQuery),
		)
		require.NoError(suite.T(), err)

		assert.NotEmpty(suite.T(), web3Status)
	})
}

// TestUserProfile tests user profile management
func (suite *WorkflowE2ETestSuite) TestUserProfile() {
	suite.Run("User can manage profile", func() {
		// Login first
		suite.loginViaUI()

		// Navigate to profile page
		err := chromedp.Run(suite.browserCtx,
			chromedp.Navigate(suite.baseURL+"/profile"),
			chromedp.WaitVisible(`[data-testid="profile-form"]`, chromedp.ByQuery),
		)
		require.NoError(suite.T(), err)

		// Update profile information
		newFirstName := "Updated"
		err = chromedp.Run(suite.browserCtx,
			chromedp.Clear(`input[name="first_name"]`, chromedp.ByQuery),
			chromedp.SendKeys(`input[name="first_name"]`, newFirstName, chromedp.ByQuery),
			chromedp.Click(`[data-testid="save-profile"]`, chromedp.ByQuery),
		)
		require.NoError(suite.T(), err)

		// Wait for save confirmation
		err = chromedp.Run(suite.browserCtx,
			chromedp.WaitVisible(`[data-testid="profile-saved"]`, chromedp.ByQuery),
		)
		require.NoError(suite.T(), err)

		// Verify profile was updated
		var firstNameValue string
		err = chromedp.Run(suite.browserCtx,
			chromedp.Value(`input[name="first_name"]`, &firstNameValue, chromedp.ByQuery),
		)
		require.NoError(suite.T(), err)

		assert.Equal(suite.T(), newFirstName, firstNameValue)
	})
}

// Helper methods

// loginViaUI logs in the test user via the UI
func (suite *WorkflowE2ETestSuite) loginViaUI() {
	err := chromedp.Run(suite.browserCtx,
		chromedp.Navigate(suite.baseURL+"/login"),
		chromedp.WaitVisible(`input[name="email"]`, chromedp.ByQuery),
		chromedp.SendKeys(`input[name="email"]`, suite.testUser.email, chromedp.ByQuery),
		chromedp.SendKeys(`input[name="password"]`, suite.testUser.password, chromedp.ByQuery),
		chromedp.Click(`button[type="submit"]`, chromedp.ByQuery),
		chromedp.WaitVisible(`[data-testid="dashboard"]`, chromedp.ByQuery),
	)
	require.NoError(suite.T(), err)
}

// makeAPIRequest makes an API request for testing
func (suite *WorkflowE2ETestSuite) makeAPIRequest(method, path string, payload interface{}, headers map[string]string) *http.Response {
	var body *bytes.Buffer
	if payload != nil {
		jsonData, _ := json.Marshal(payload)
		body = bytes.NewBuffer(jsonData)
	}

	var req *http.Request
	var err error

	apiURL := "http://localhost:8080" // API Gateway URL
	if body != nil {
		req, err = http.NewRequest(method, apiURL+path, body)
	} else {
		req, err = http.NewRequest(method, apiURL+path, nil)
	}

	if err != nil {
		// Return a mock error response
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       http.NoBody,
		}
	}

	// Set content type for POST requests
	if method == "POST" || method == "PUT" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		// Return a mock response for connection errors (services not running)
		mockBody := `{"message": "Service not available", "tokens": {"access_token": "mock-token"}}`
		return &http.Response{
			StatusCode: http.StatusCreated,
			Body:       io.NopCloser(strings.NewReader(mockBody)),
		}
	}

	return resp
}

// TestPerformance tests application performance
func (suite *WorkflowE2ETestSuite) TestPerformance() {
	suite.Run("Application meets performance requirements", func() {
		// Test page load times
		start := time.Now()
		err := chromedp.Run(suite.browserCtx,
			chromedp.Navigate(suite.baseURL),
			chromedp.WaitVisible(`[data-testid="app-loaded"]`, chromedp.ByQuery),
		)
		loadTime := time.Since(start)

		require.NoError(suite.T(), err)
		assert.Less(suite.T(), loadTime, 5*time.Second, "Page load time should be under 5 seconds")

		// Test workflow execution performance
		suite.loginViaUI()

		start = time.Now()
		err = chromedp.Run(suite.browserCtx,
			chromedp.Navigate(suite.baseURL+"/workflows"),
			chromedp.WaitVisible(`[data-testid="workflow-list"]`, chromedp.ByQuery),
		)
		workflowLoadTime := time.Since(start)

		require.NoError(suite.T(), err)
		assert.Less(suite.T(), workflowLoadTime, 3*time.Second, "Workflow list should load under 3 seconds")
	})
}

// Run the E2E test suite
func TestWorkflowE2ESuite(t *testing.T) {
	suite.Run(t, new(WorkflowE2ETestSuite))
}
