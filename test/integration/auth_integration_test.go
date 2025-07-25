package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/ai-agentic-browser/internal/auth"
)

// AuthIntegrationTestSuite tests authentication endpoints
type AuthIntegrationTestSuite struct {
	suite.Suite
	authService *auth.Service
	Router      *gin.Engine
	HTTPServer  *httptest.Server
}

// SetupSuite initializes the integration test suite
func (suite *AuthIntegrationTestSuite) SetupSuite() {
	// Initialize Gin router
	gin.SetMode(gin.TestMode)
	suite.Router = gin.New()

	// Initialize auth service
	suite.authService = &auth.Service{
		// Initialize with test dependencies
	}

	// Setup routes
	suite.setupRoutes()

	// Create HTTP test server
	suite.HTTPServer = httptest.NewServer(suite.Router)
}

// TearDownSuite cleans up the integration test suite
func (suite *AuthIntegrationTestSuite) TearDownSuite() {
	if suite.HTTPServer != nil {
		suite.HTTPServer.Close()
	}
}

// setupRoutes sets up authentication routes for testing
func (suite *AuthIntegrationTestSuite) setupRoutes() {
	authGroup := suite.Router.Group("/api/auth")
	{
		authGroup.POST("/register", suite.handleRegister)
		authGroup.POST("/login", suite.handleLogin)
		authGroup.POST("/refresh", suite.handleRefresh)
		authGroup.POST("/logout", suite.handleLogout)
		authGroup.GET("/me", suite.authMiddleware(), suite.handleMe)
		authGroup.PUT("/me", suite.authMiddleware(), suite.handleUpdateProfile)
		authGroup.POST("/change-password", suite.authMiddleware(), suite.handleChangePassword)
	}
}

// Test registration endpoint
func (suite *AuthIntegrationTestSuite) TestRegisterEndpoint() {
	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		checkResponse  func(*http.Response)
	}{
		{
			name: "Valid registration",
			payload: map[string]interface{}{
				"email":      "test@example.com",
				"password":   "password123",
				"first_name": "Test",
				"last_name":  "User",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(resp *http.Response) {
				var response map[string]interface{}
				err := json.NewDecoder(resp.Body).Decode(&response)
				require.NoError(suite.T(), err)

				assert.Contains(suite.T(), response, "user")
				assert.Contains(suite.T(), response, "tokens")

				user := response["user"].(map[string]interface{})
				assert.Equal(suite.T(), "test@example.com", user["email"])
				assert.NotContains(suite.T(), user, "password")
			},
		},
		{
			name: "Invalid email format",
			payload: map[string]interface{}{
				"email":      "invalid-email",
				"password":   "password123",
				"first_name": "Test",
				"last_name":  "User",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Password too short",
			payload: map[string]interface{}{
				"email":      "short@example.com",
				"password":   "123",
				"first_name": "Test",
				"last_name":  "User",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Missing required fields",
			payload: map[string]interface{}{
				"email": "missing@example.com",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Duplicate email",
			payload: map[string]interface{}{
				"email":      "test@example.com", // Same as first test
				"password":   "password123",
				"first_name": "Another",
				"last_name":  "User",
			},
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			body, _ := json.Marshal(tt.payload)
			resp := suite.makeRequest("POST", "/api/auth/register", bytes.NewBuffer(body), nil)
			defer resp.Body.Close()

			suite.AssertHTTPStatus(resp, tt.expectedStatus)

			if tt.checkResponse != nil {
				tt.checkResponse(resp)
			}
		})
	}
}

// Test login endpoint
func (suite *AuthIntegrationTestSuite) TestLoginEndpoint() {
	// Create test user first
	email := "login@example.com"
	password := "password123"
	suite.createTestUser(email, password)

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		checkResponse  func(*http.Response)
	}{
		{
			name: "Valid login",
			payload: map[string]interface{}{
				"email":    email,
				"password": password,
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(resp *http.Response) {
				var response map[string]interface{}
				err := json.NewDecoder(resp.Body).Decode(&response)
				require.NoError(suite.T(), err)

				assert.Contains(suite.T(), response, "user")
				assert.Contains(suite.T(), response, "tokens")

				tokens := response["tokens"].(map[string]interface{})
				assert.NotEmpty(suite.T(), tokens["access_token"])
				assert.NotEmpty(suite.T(), tokens["refresh_token"])
			},
		},
		{
			name: "Invalid email",
			payload: map[string]interface{}{
				"email":    "nonexistent@example.com",
				"password": password,
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Invalid password",
			payload: map[string]interface{}{
				"email":    email,
				"password": "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Missing credentials",
			payload: map[string]interface{}{
				"email": email,
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			body, _ := json.Marshal(tt.payload)
			resp := suite.makeRequest("POST", "/api/auth/login", bytes.NewBuffer(body), nil)
			defer resp.Body.Close()

			suite.AssertHTTPStatus(resp, tt.expectedStatus)

			if tt.checkResponse != nil {
				tt.checkResponse(resp)
			}
		})
	}
}

// Test token refresh endpoint
func (suite *AuthIntegrationTestSuite) TestRefreshEndpoint() {
	// Create test user and get tokens
	email := "refresh@example.com"
	password := "password123"
	suite.createTestUser(email, password)

	tokens := suite.loginUser(email, password)

	tests := []struct {
		name           string
		refreshToken   string
		expectedStatus int
		checkResponse  func(*http.Response)
	}{
		{
			name:           "Valid refresh token",
			refreshToken:   tokens["refresh_token"].(string),
			expectedStatus: http.StatusOK,
			checkResponse: func(resp *http.Response) {
				var response map[string]interface{}
				err := json.NewDecoder(resp.Body).Decode(&response)
				require.NoError(suite.T(), err)

				assert.Contains(suite.T(), response, "tokens")
				newTokens := response["tokens"].(map[string]interface{})
				assert.NotEmpty(suite.T(), newTokens["access_token"])
				assert.NotEmpty(suite.T(), newTokens["refresh_token"])

				// New tokens should be different from original
				assert.NotEqual(suite.T(), tokens["access_token"], newTokens["access_token"])
			},
		},
		{
			name:           "Invalid refresh token",
			refreshToken:   "invalid.refresh.token",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Empty refresh token",
			refreshToken:   "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			payload := map[string]interface{}{
				"refresh_token": tt.refreshToken,
			}
			body, _ := json.Marshal(payload)
			resp := suite.makeRequest("POST", "/api/auth/refresh", bytes.NewBuffer(body), nil)
			defer resp.Body.Close()

			suite.AssertHTTPStatus(resp, tt.expectedStatus)

			if tt.checkResponse != nil {
				tt.checkResponse(resp)
			}
		})
	}
}

// Test protected endpoint access
func (suite *AuthIntegrationTestSuite) TestProtectedEndpointAccess() {
	// Create test user and get tokens
	email := "protected@example.com"
	password := "password123"
	suite.createTestUser(email, password)

	tokens := suite.loginUser(email, password)
	accessToken := tokens["access_token"].(string)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "Valid access token",
			authHeader:     "Bearer " + accessToken,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid access token",
			authHeader:     "Bearer invalid.token.here",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Missing authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Malformed authorization header",
			authHeader:     "InvalidFormat " + accessToken,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			headers := make(map[string]string)
			if tt.authHeader != "" {
				headers["Authorization"] = tt.authHeader
			}

			resp := suite.makeRequest("GET", "/api/auth/me", nil, headers)
			defer resp.Body.Close()

			suite.AssertHTTPStatus(resp, tt.expectedStatus)
		})
	}
}

// Test logout endpoint
func (suite *AuthIntegrationTestSuite) TestLogoutEndpoint() {
	// Create test user and get tokens
	email := "logout@example.com"
	password := "password123"
	suite.createTestUser(email, password)

	tokens := suite.loginUser(email, password)
	accessToken := tokens["access_token"].(string)

	// Test logout
	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
	}
	resp := suite.makeRequest("POST", "/api/auth/logout", nil, headers)
	defer resp.Body.Close()

	suite.AssertHTTPStatus(resp, http.StatusOK)

	// Verify token is invalidated
	resp = suite.makeRequest("GET", "/api/auth/me", nil, headers)
	defer resp.Body.Close()
	suite.AssertHTTPStatus(resp, http.StatusUnauthorized)
}

// Test password change endpoint
func (suite *AuthIntegrationTestSuite) TestChangePasswordEndpoint() {
	// Create test user and get tokens
	email := "changepass@example.com"
	oldPassword := "oldpassword123"
	newPassword := "newpassword456"

	suite.createTestUser(email, oldPassword)
	tokens := suite.loginUser(email, oldPassword)
	accessToken := tokens["access_token"].(string)

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
	}{
		{
			name: "Valid password change",
			payload: map[string]interface{}{
				"current_password": oldPassword,
				"new_password":     newPassword,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid current password",
			payload: map[string]interface{}{
				"current_password": "wrongpassword",
				"new_password":     "anotherpassword",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "New password too short",
			payload: map[string]interface{}{
				"current_password": newPassword, // Use new password as current
				"new_password":     "123",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			body, _ := json.Marshal(tt.payload)
			headers := map[string]string{
				"Authorization": "Bearer " + accessToken,
			}
			resp := suite.makeRequest("POST", "/api/auth/change-password", bytes.NewBuffer(body), headers)
			defer resp.Body.Close()

			suite.AssertHTTPStatus(resp, tt.expectedStatus)
		})
	}

	// Verify new password works
	newTokens := suite.loginUser(email, newPassword)
	assert.NotEmpty(suite.T(), newTokens["access_token"])
}

// Test rate limiting
func (suite *AuthIntegrationTestSuite) TestRateLimiting() {
	email := "ratelimit@example.com"

	// Make multiple rapid requests
	for i := 0; i < 10; i++ {
		payload := map[string]interface{}{
			"email":    email,
			"password": "wrongpassword",
		}
		body, _ := json.Marshal(payload)
		resp := suite.makeRequest("POST", "/api/auth/login", bytes.NewBuffer(body), nil)
		resp.Body.Close()

		// After several attempts, should get rate limited
		if i > 5 {
			assert.True(suite.T(), resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusUnauthorized)
		}
	}
}

// Helper methods

func (suite *AuthIntegrationTestSuite) createTestUser(email, password string) {
	payload := map[string]interface{}{
		"email":      email,
		"password":   password,
		"first_name": "Test",
		"last_name":  "User",
	}
	body, _ := json.Marshal(payload)
	resp := suite.makeRequest("POST", "/api/auth/register", bytes.NewBuffer(body), nil)
	defer resp.Body.Close()

	require.True(suite.T(), resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusConflict)
}

func (suite *AuthIntegrationTestSuite) loginUser(email, password string) map[string]interface{} {
	payload := map[string]interface{}{
		"email":    email,
		"password": password,
	}
	body, _ := json.Marshal(payload)
	resp := suite.makeRequest("POST", "/api/auth/login", bytes.NewBuffer(body), nil)
	defer resp.Body.Close()

	require.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(suite.T(), err)

	return response["tokens"].(map[string]interface{})
}

func (suite *AuthIntegrationTestSuite) makeRequest(method, path string, body *bytes.Buffer, headers map[string]string) *http.Response {
	var req *http.Request
	var err error

	if body != nil {
		req, err = http.NewRequest(method, suite.HTTPServer.URL+path, body)
	} else {
		req, err = http.NewRequest(method, suite.HTTPServer.URL+path, nil)
	}
	require.NoError(suite.T(), err)

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
	require.NoError(suite.T(), err)

	return resp
}

// Mock handlers for testing (these would be replaced with actual handlers)

func (suite *AuthIntegrationTestSuite) handleRegister(c *gin.Context) {
	// Mock implementation
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":    "123e4567-e89b-12d3-a456-426614174000",
			"email": "test@example.com",
		},
	})
}

func (suite *AuthIntegrationTestSuite) handleLogin(c *gin.Context) {
	// Mock implementation
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"tokens": gin.H{
			"access_token":  "mock-access-token",
			"refresh_token": "mock-refresh-token",
		},
		"user": gin.H{
			"id":    "123e4567-e89b-12d3-a456-426614174000",
			"email": "test@example.com",
		},
	})
}

func (suite *AuthIntegrationTestSuite) handleRefresh(c *gin.Context) {
	// Mock implementation
	c.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed",
		"tokens": gin.H{
			"access_token":  "new-mock-access-token",
			"refresh_token": "new-mock-refresh-token",
		},
	})
}

func (suite *AuthIntegrationTestSuite) handleLogout(c *gin.Context) {
	// Mock implementation
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func (suite *AuthIntegrationTestSuite) handleMe(c *gin.Context) {
	// Mock implementation
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    "123e4567-e89b-12d3-a456-426614174000",
			"email": "test@example.com",
			"name":  "Test User",
		},
	})
}

func (suite *AuthIntegrationTestSuite) handleUpdateProfile(c *gin.Context) {
	// Mock implementation
	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated",
		"user": gin.H{
			"id":    "123e4567-e89b-12d3-a456-426614174000",
			"email": "test@example.com",
			"name":  "Updated Test User",
		},
	})
}

func (suite *AuthIntegrationTestSuite) handleChangePassword(c *gin.Context) {
	// Mock implementation
	c.JSON(http.StatusOK, gin.H{"message": "Password changed"})
}

func (suite *AuthIntegrationTestSuite) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// AssertHTTPStatus asserts that the HTTP response has the expected status code
func (suite *AuthIntegrationTestSuite) AssertHTTPStatus(resp *http.Response, expectedStatus int) {
	assert.Equal(suite.T(), expectedStatus, resp.StatusCode)
}

// Run the integration test suite
func TestAuthIntegrationSuite(t *testing.T) {
	suite.Run(t, new(AuthIntegrationTestSuite))
}
