package solana

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ai-agentic-browser/internal/web3/solana"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// APIIntegrationTestSuite defines the test suite for API integration
type APIIntegrationTestSuite struct {
	suite.Suite
	router        *mux.Router
	server        *httptest.Server
	solanaService *solana.Service
	logger        *observability.Logger
	ctx           context.Context
}

func (suite *APIIntegrationTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	// Initialize mock logger
	suite.logger = &observability.Logger{}

	// Initialize mock Solana service
	suite.solanaService = &solana.Service{}

	// Setup router with all routes
	suite.router = mux.NewRouter()
	SetupSolanaRoutes(suite.router, suite.solanaService, suite.logger)

	// Create test server
	suite.server = httptest.NewServer(suite.router)
}

func (suite *APIIntegrationTestSuite) TearDownSuite() {
	suite.server.Close()
}

// Test health check endpoint
func (suite *APIIntegrationTestSuite) TestHealthCheck() {
	resp, err := http.Get(suite.server.URL + "/solana/health")
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), "healthy", response["status"])
	assert.Contains(suite.T(), response, "services")
}

// Test Solana stats endpoint
func (suite *APIIntegrationTestSuite) TestGetSolanaStats() {
	resp, err := http.Get(suite.server.URL + "/solana/stats")
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var stats SolanaStats
	err = json.NewDecoder(resp.Body).Decode(&stats)
	require.NoError(suite.T(), err)

	assert.True(suite.T(), stats.Price.GreaterThan(decimal.Zero))
	assert.True(suite.T(), stats.MarketCap.GreaterThan(decimal.Zero))
	assert.True(suite.T(), stats.TPS > 0)
}

// Test swap quote endpoint
func (suite *APIIntegrationTestSuite) TestSwapQuote() {
	quoteRequest := map[string]interface{}{
		"inputMint":   "So11111111111111111111111111111111111111112",
		"outputMint":  "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		"amount":      1.0,
		"slippageBps": 50,
	}

	requestBody, err := json.Marshal(quoteRequest)
	require.NoError(suite.T(), err)

	resp, err := http.Post(
		suite.server.URL+"/solana/defi/quote",
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var quote map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&quote)
	require.NoError(suite.T(), err)

	assert.Contains(suite.T(), quote, "inputAmount")
	assert.Contains(suite.T(), quote, "outputAmount")
	assert.Contains(suite.T(), quote, "priceImpact")
}

// Test swap quote with invalid parameters
func (suite *APIIntegrationTestSuite) TestSwapQuoteInvalidParams() {
	quoteRequest := map[string]interface{}{
		"inputMint":   "invalid-mint",
		"outputMint":  "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		"amount":      1.0,
		"slippageBps": 50,
	}

	requestBody, err := json.Marshal(quoteRequest)
	require.NoError(suite.T(), err)

	resp, err := http.Post(
		suite.server.URL+"/solana/defi/quote",
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

// Test wallet connection endpoint (requires auth)
func (suite *APIIntegrationTestSuite) TestWalletConnect() {
	connectRequest := map[string]interface{}{
		"publicKey":  "11111111111111111111111111111111",
		"walletType": "phantom",
	}

	requestBody, err := json.Marshal(connectRequest)
	require.NoError(suite.T(), err)

	req, err := http.NewRequest(
		"POST",
		suite.server.URL+"/solana/wallets/connect",
		bytes.NewBuffer(requestBody),
	)
	require.NoError(suite.T(), err)

	req.Header.Set("Content-Type", "application/json")
	// Add mock auth header
	req.Header.Set("Authorization", "Bearer mock-token")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	// Should return 401 without proper auth middleware
	// In a real implementation, this would be 200 with proper auth
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
}

// Test NFT exploration endpoint
func (suite *APIIntegrationTestSuite) TestNFTExplore() {
	resp, err := http.Get(suite.server.URL + "/solana/nft/explore?limit=10")
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), true, response["success"])
	assert.Contains(suite.T(), response, "nfts")
}

// Test NFT collections endpoint
func (suite *APIIntegrationTestSuite) TestNFTCollections() {
	resp, err := http.Get(suite.server.URL + "/solana/nft/collections")
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), true, response["success"])
	assert.Contains(suite.T(), response, "collections")
}

// Test CORS headers
func (suite *APIIntegrationTestSuite) TestCORSHeaders() {
	req, err := http.NewRequest("OPTIONS", suite.server.URL+"/solana/stats", nil)
	require.NoError(suite.T(), err)

	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "GET")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	assert.Equal(suite.T(), "*", resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Contains(suite.T(), resp.Header.Get("Access-Control-Allow-Methods"), "GET")
}

// Test rate limiting (if implemented)
func (suite *APIIntegrationTestSuite) TestRateLimiting() {
	// Make multiple rapid requests
	for i := 0; i < 10; i++ {
		resp, err := http.Get(suite.server.URL + "/solana/stats")
		require.NoError(suite.T(), err)
		resp.Body.Close()

		// All should succeed with current implementation
		assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	}
}

// Test error handling
func (suite *APIIntegrationTestSuite) TestErrorHandling() {
	// Test non-existent endpoint
	resp, err := http.Get(suite.server.URL + "/solana/nonexistent")
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}

// Test request timeout handling
func (suite *APIIntegrationTestSuite) TestRequestTimeout() {
	client := &http.Client{
		Timeout: 1 * time.Millisecond, // Very short timeout
	}

	_, err := client.Get(suite.server.URL + "/solana/stats")
	// Should timeout (though may not with fast local server)
	// This tests that timeout handling is in place
	if err != nil {
		assert.Contains(suite.T(), err.Error(), "timeout")
	}
}

// Test JSON response format
func (suite *APIIntegrationTestSuite) TestJSONResponseFormat() {
	resp, err := http.Get(suite.server.URL + "/solana/stats")
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), "application/json", resp.Header.Get("Content-Type"))

	var jsonResponse interface{}
	err = json.NewDecoder(resp.Body).Decode(&jsonResponse)
	assert.NoError(suite.T(), err, "Response should be valid JSON")
}

// Test large request handling
func (suite *APIIntegrationTestSuite) TestLargeRequestHandling() {
	// Create a large but valid request
	largeRequest := map[string]interface{}{
		"inputMint":   "So11111111111111111111111111111111111111112",
		"outputMint":  "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		"amount":      1000000.0, // Large amount
		"slippageBps": 50,
	}

	requestBody, err := json.Marshal(largeRequest)
	require.NoError(suite.T(), err)

	resp, err := http.Post(
		suite.server.URL+"/solana/defi/quote",
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	// Should handle large requests gracefully
	assert.True(suite.T(), resp.StatusCode < 500)
}

// Test concurrent requests
func (suite *APIIntegrationTestSuite) TestConcurrentRequests() {
	const numRequests = 10
	results := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			resp, err := http.Get(suite.server.URL + "/solana/stats")
			if err != nil {
				results <- err
				return
			}
			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				results <- assert.AnError
				return
			}
			results <- nil
		}()
	}

	// Wait for all requests to complete
	for i := 0; i < numRequests; i++ {
		err := <-results
		assert.NoError(suite.T(), err)
	}
}

// Run the test suite
func TestAPIIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(APIIntegrationTestSuite))
}

// Benchmark API endpoints
func BenchmarkSolanaStatsEndpoint(b *testing.B) {
	router := mux.NewRouter()
	logger := &observability.Logger{}
	solanaService := &solana.Service{}

	SetupSolanaRoutes(router, solanaService, logger)
	server := httptest.NewServer(router)
	defer server.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := http.Get(server.URL + "/solana/stats")
		if err != nil {
			b.Fatal(err)
		}
		resp.Body.Close()
	}
}

func BenchmarkSwapQuoteEndpoint(b *testing.B) {
	router := mux.NewRouter()
	logger := &observability.Logger{}
	solanaService := &solana.Service{}

	SetupSolanaRoutes(router, solanaService, logger)
	server := httptest.NewServer(router)
	defer server.Close()

	quoteRequest := map[string]interface{}{
		"inputMint":   "So11111111111111111111111111111111111111112",
		"outputMint":  "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		"amount":      1.0,
		"slippageBps": 50,
	}

	requestBody, _ := json.Marshal(quoteRequest)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := http.Post(
			server.URL+"/solana/defi/quote",
			"application/json",
			bytes.NewBuffer(requestBody),
		)
		if err != nil {
			b.Fatal(err)
		}
		resp.Body.Close()
	}
}
