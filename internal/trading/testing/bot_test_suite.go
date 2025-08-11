package testing

import (
	"context"
	"testing"
	"time"

	"github.com/ai-agentic-browser/internal/trading"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// BotTestSuite provides a comprehensive test suite for trading bots
type BotTestSuite struct {
	suite.Suite
	framework *BotTestFramework
	executor  *TestExecutor
}

// SetupSuite initializes the test suite
func (bts *BotTestSuite) SetupSuite() {
	bts.framework = NewBotTestFramework()
	bts.framework.SetupSuite()

	bts.executor = NewTestExecutor(bts.framework.logger, bts.framework)
	bts.executor.Start(context.Background())
}

// TearDownSuite cleans up the test suite
func (bts *BotTestSuite) TearDownSuite() {
	if bts.executor != nil {
		bts.executor.Stop(context.Background())
	}
	if bts.framework != nil {
		bts.framework.TearDownSuite()
	}
}

// TestBotCreation tests bot creation functionality
func (bts *BotTestSuite) TestBotCreation() {
	t := bts.T()

	testCases := []struct {
		name       string
		strategy   string
		config     *trading.BotConfig
		shouldPass bool
	}{
		{
			name:     "Valid DCA Bot",
			strategy: "dca",
			config: &trading.BotConfig{
				TradingPairs: []string{"BTC/USDT"},
				Exchange:     "binance",
				BaseCurrency: "USDT",
				StrategyParams: map[string]interface{}{
					"interval":      "1h",
					"amount":        100.0,
					"max_positions": 5,
				},
				Enabled: true,
			},
			shouldPass: true,
		},
		{
			name:     "Valid Grid Bot",
			strategy: "grid",
			config: &trading.BotConfig{
				TradingPairs: []string{"ETH/USDT"},
				Exchange:     "binance",
				BaseCurrency: "USDT",
				StrategyParams: map[string]interface{}{
					"grid_levels":  10,
					"grid_spacing": 0.5,
				},
				Enabled: true,
			},
			shouldPass: true,
		},
		{
			name:     "Invalid Strategy",
			strategy: "invalid",
			config: &trading.BotConfig{
				TradingPairs: []string{"BTC/USDT"},
				Exchange:     "binance",
				BaseCurrency: "USDT",
				StrategyParams: map[string]interface{}{
					"invalid_param": "value",
				},
				Enabled: true,
			},
			shouldPass: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := &TestRequest{
				Type:     TestTypeUnit,
				Strategy: tc.strategy,
				Config:   tc.config,
				Timeout:  30 * time.Second,
			}

			testID, err := bts.executor.SubmitTest(request)
			require.NoError(t, err)
			require.NotEmpty(t, testID)

			// Wait for test completion
			time.Sleep(5 * time.Second)

			result, err := bts.executor.GetTestResult(testID)
			require.NoError(t, err)
			require.NotNil(t, result)

			if tc.shouldPass {
				assert.True(t, result.Passed, "Test should pass for valid configuration")
				assert.Empty(t, result.FailureReasons, "Should have no failure reasons")
			} else {
				assert.False(t, result.Passed, "Test should fail for invalid configuration")
				assert.NotEmpty(t, result.FailureReasons, "Should have failure reasons")
			}
		})
	}
}

// TestBotIntegration tests bot integration with exchange and market data
func (bts *BotTestSuite) TestBotIntegration() {
	t := bts.T()

	strategies := []string{"dca", "grid", "momentum"}

	for _, strategy := range strategies {
		t.Run(strategy, func(t *testing.T) {
			config := &trading.BotConfig{
				TradingPairs: []string{"BTC/USDT"},
				Exchange:     "binance",
				BaseCurrency: "USDT",
				StrategyParams: map[string]interface{}{
					"strategy": strategy,
				},
				Enabled: true,
			}

			request := &TestRequest{
				Type:     TestTypeIntegration,
				Strategy: strategy,
				Config:   config,
				Timeout:  60 * time.Second,
			}

			testID, err := bts.executor.SubmitTest(request)
			require.NoError(t, err)

			// Wait for test completion
			time.Sleep(10 * time.Second)

			result, err := bts.executor.GetTestResult(testID)
			require.NoError(t, err)
			require.NotNil(t, result)

			assert.True(t, result.Passed, "Integration test should pass")
			assert.Equal(t, TestStatusCompleted, result.Status)
			assert.Greater(t, result.Duration, time.Duration(0))
		})
	}
}

// TestBotBacktesting tests backtesting functionality
func (bts *BotTestSuite) TestBotBacktesting() {
	t := bts.T()

	testScenarios := []struct {
		name               string
		strategy           string
		marketCondition    string
		expectedMinWinRate decimal.Decimal
	}{
		{
			name:               "DCA Bull Market",
			strategy:           "dca",
			marketCondition:    "bull",
			expectedMinWinRate: decimal.NewFromFloat(0.60),
		},
		{
			name:               "Grid Sideways Market",
			strategy:           "grid",
			marketCondition:    "sideways",
			expectedMinWinRate: decimal.NewFromFloat(0.70),
		},
		{
			name:               "Momentum Trending Market",
			strategy:           "momentum",
			marketCondition:    "bull",
			expectedMinWinRate: decimal.NewFromFloat(0.55),
		},
	}

	for _, scenario := range testScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			config := &trading.BotConfig{
				TradingPairs: []string{"BTC/USDT"},
				Exchange:     "binance",
				BaseCurrency: "USDT",
				StrategyParams: map[string]interface{}{
					"strategy": scenario.strategy,
				},
				Enabled: true,
			}

			testScenario := &TestScenario{
				ID:              scenario.name,
				Name:            scenario.name,
				Type:            TestTypeBacktest,
				Strategy:        scenario.strategy,
				MarketCondition: scenario.marketCondition,
				Duration:        30 * time.Second,
				ExpectedResults: &ExpectedResults{
					MinWinRate:     scenario.expectedMinWinRate,
					MaxDrawdown:    decimal.NewFromFloat(0.20),
					MinSharpeRatio: decimal.NewFromFloat(0.30),
					MaxRiskScore:   80,
				},
			}

			request := &TestRequest{
				Type:            TestTypeBacktest,
				Strategy:        scenario.strategy,
				Config:          config,
				Scenario:        testScenario,
				ExpectedResults: testScenario.ExpectedResults,
				Timeout:         120 * time.Second,
			}

			testID, err := bts.executor.SubmitTest(request)
			require.NoError(t, err)

			// Wait for backtest completion
			time.Sleep(45 * time.Second)

			result, err := bts.executor.GetTestResult(testID)
			require.NoError(t, err)
			require.NotNil(t, result)

			assert.Equal(t, TestStatusCompleted, result.Status)
			assert.NotNil(t, result.Metrics)

			// Log results for analysis
			t.Logf("Backtest Results for %s:", scenario.name)
			t.Logf("  Total Trades: %d", result.Metrics.TotalTrades)
			t.Logf("  Win Rate: %.2f%%", result.Metrics.WinRate.Mul(decimal.NewFromFloat(100)).InexactFloat64())
			t.Logf("  Total Return: %s", result.Metrics.TotalReturn.String())
			t.Logf("  Sharpe Ratio: %s", result.Metrics.SharpeRatio.String())
			t.Logf("  Max Drawdown: %.2f%%", result.Metrics.MaxDrawdown.Mul(decimal.NewFromFloat(100)).InexactFloat64())

			if !result.Passed {
				t.Logf("  Failure Reasons: %v", result.FailureReasons)
			}
		})
	}
}

// TestBotStressTesting tests bot behavior under stress conditions
func (bts *BotTestSuite) TestBotStressTesting() {
	t := bts.T()

	strategies := []string{"dca", "grid", "momentum"}

	for _, strategy := range strategies {
		t.Run(strategy+"_stress", func(t *testing.T) {
			config := &trading.BotConfig{
				TradingPairs: []string{"BTC/USDT", "ETH/USDT"},
				Exchange:     "binance",
				BaseCurrency: "USDT",
				StrategyParams: map[string]interface{}{
					"strategy": strategy,
				},
				Enabled: true,
			}

			request := &TestRequest{
				Type:     TestTypeStress,
				Strategy: strategy,
				Config:   config,
				Timeout:  60 * time.Second,
			}

			testID, err := bts.executor.SubmitTest(request)
			require.NoError(t, err)

			// Wait for stress test completion
			time.Sleep(25 * time.Second)

			result, err := bts.executor.GetTestResult(testID)
			require.NoError(t, err)
			require.NotNil(t, result)

			assert.Equal(t, TestStatusCompleted, result.Status)

			// Stress tests should generally pass (bot should survive)
			if !result.Passed {
				t.Logf("Stress test failed for %s: %v", strategy, result.FailureReasons)
			}
		})
	}
}

// TestBotPerformanceMetrics tests performance metric calculations
func (bts *BotTestSuite) TestBotPerformanceMetrics() {
	t := bts.T()

	config := &trading.BotConfig{
		TradingPairs: []string{"BTC/USDT"},
		Exchange:     "binance",
		BaseCurrency: "USDT",
		StrategyParams: map[string]interface{}{
			"strategy": "dca",
		},
		Enabled: true,
	}

	request := &TestRequest{
		Type:     TestTypePerformance,
		Strategy: "dca",
		Config:   config,
		Timeout:  60 * time.Second,
	}

	testID, err := bts.executor.SubmitTest(request)
	require.NoError(t, err)

	// Wait for performance test completion
	time.Sleep(35 * time.Second)

	result, err := bts.executor.GetTestResult(testID)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.Metrics)

	metrics := result.Metrics

	// Validate metric ranges
	assert.GreaterOrEqual(t, metrics.WinRate.InexactFloat64(), 0.0)
	assert.LessOrEqual(t, metrics.WinRate.InexactFloat64(), 1.0)

	assert.GreaterOrEqual(t, metrics.MaxDrawdown.InexactFloat64(), 0.0)
	assert.LessOrEqual(t, metrics.MaxDrawdown.InexactFloat64(), 1.0)

	assert.GreaterOrEqual(t, metrics.RiskScore, 0)
	assert.LessOrEqual(t, metrics.RiskScore, 100)

	// Log metrics for analysis
	t.Logf("Performance Metrics:")
	t.Logf("  Total Trades: %d", metrics.TotalTrades)
	t.Logf("  Winning Trades: %d", metrics.WinningTrades)
	t.Logf("  Losing Trades: %d", metrics.LosingTrades)
	t.Logf("  Win Rate: %.2f%%", metrics.WinRate.Mul(decimal.NewFromFloat(100)).InexactFloat64())
	t.Logf("  Total Return: %s", metrics.TotalReturn.String())
	t.Logf("  Sharpe Ratio: %s", metrics.SharpeRatio.String())
	t.Logf("  Max Drawdown: %.2f%%", metrics.MaxDrawdown.Mul(decimal.NewFromFloat(100)).InexactFloat64())
	t.Logf("  Profit Factor: %s", metrics.ProfitFactor.String())
}

// TestBotRiskManagement tests risk management integration
func (bts *BotTestSuite) TestBotRiskManagement() {
	t := bts.T()

	// Test with high risk configuration
	config := &trading.BotConfig{
		TradingPairs: []string{"BTC/USDT"},
		Exchange:     "binance",
		BaseCurrency: "USDT",
		StrategyParams: map[string]interface{}{
			"strategy":      "momentum",
			"max_positions": 20, // High position count
		},
		Enabled: true,
	}

	request := &TestRequest{
		Type:     TestTypeIntegration,
		Strategy: "momentum",
		Config:   config,
		Timeout:  60 * time.Second,
	}

	testID, err := bts.executor.SubmitTest(request)
	require.NoError(t, err)

	// Wait for test completion
	time.Sleep(10 * time.Second)

	result, err := bts.executor.GetTestResult(testID)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Risk management should prevent excessive risk
	assert.Equal(t, TestStatusCompleted, result.Status)
}

// TestConcurrentBots tests multiple bots running concurrently
func (bts *BotTestSuite) TestConcurrentBots() {
	t := bts.T()

	strategies := []string{"dca", "grid", "momentum"}
	testIDs := make([]string, 0)

	// Submit multiple tests concurrently
	for i, strategy := range strategies {
		config := &trading.BotConfig{
			TradingPairs: []string{"BTC/USDT"},
			Exchange:     "binance",
			BaseCurrency: "USDT",
			StrategyParams: map[string]interface{}{
				"strategy": strategy,
			},
			Enabled: true,
		}

		request := &TestRequest{
			Type:     TestTypeIntegration,
			Strategy: strategy,
			Config:   config,
			Priority: i + 1,
			Timeout:  60 * time.Second,
		}

		testID, err := bts.executor.SubmitTest(request)
		require.NoError(t, err)
		testIDs = append(testIDs, testID)
	}

	// Wait for all tests to complete
	time.Sleep(15 * time.Second)

	// Check all test results
	for i, testID := range testIDs {
		result, err := bts.executor.GetTestResult(testID)
		require.NoError(t, err, "Test %d should complete", i)
		require.NotNil(t, result)

		assert.Equal(t, TestStatusCompleted, result.Status, "Test %d should complete successfully", i)
		assert.True(t, result.Passed, "Test %d should pass", i)
	}
}

// TestBotErrorHandling tests error handling and recovery
func (bts *BotTestSuite) TestBotErrorHandling() {
	t := bts.T()

	// Test with invalid configuration
	config := &trading.BotConfig{
		TradingPairs: []string{}, // Empty trading pairs
		Exchange:     "binance",
		BaseCurrency: "USDT",
		StrategyParams: map[string]interface{}{
			"strategy": "invalid_strategy",
		},
		Enabled: true,
	}

	request := &TestRequest{
		Type:     TestTypeUnit,
		Strategy: "invalid_strategy",
		Config:   config,
		Timeout:  30 * time.Second,
	}

	testID, err := bts.executor.SubmitTest(request)
	require.NoError(t, err)

	// Wait for test completion
	time.Sleep(5 * time.Second)

	result, err := bts.executor.GetTestResult(testID)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Should fail gracefully with proper error messages
	assert.False(t, result.Passed, "Test should fail for invalid configuration")
	assert.NotEmpty(t, result.FailureReasons, "Should have failure reasons")
	assert.Equal(t, TestStatusCompleted, result.Status, "Should complete even with errors")
}

// TestBotTestSuite runs the complete bot test suite
func TestBotTestSuite(t *testing.T) {
	suite.Run(t, new(BotTestSuite))
}
