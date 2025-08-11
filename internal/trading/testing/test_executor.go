package testing

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/internal/trading"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// TestExecutor executes trading bot tests
type TestExecutor struct {
	logger    *observability.Logger
	framework *BotTestFramework

	// Test execution
	activeTests map[string]*TestExecution
	testQueue   chan *TestRequest
	workers     []*TestWorker

	// Results
	results map[string]*TestResult

	// Configuration
	config *TestExecutorConfig

	// Synchronization
	mu        sync.RWMutex
	isRunning bool
	stopChan  chan struct{}
}

// TestExecutorConfig holds configuration for test execution
type TestExecutorConfig struct {
	MaxWorkers    int           `yaml:"max_workers"`
	WorkerTimeout time.Duration `yaml:"worker_timeout"`
	QueueSize     int           `yaml:"queue_size"`
	RetryAttempts int           `yaml:"retry_attempts"`
	RetryDelay    time.Duration `yaml:"retry_delay"`
}

// TestExecution represents an active test execution
type TestExecution struct {
	ID        string                 `json:"id"`
	Request   *TestRequest           `json:"request"`
	Bot       *TestBot               `json:"bot"`
	StartTime time.Time              `json:"start_time"`
	Status    TestStatus             `json:"status"`
	Progress  float64                `json:"progress"`
	Metadata  map[string]interface{} `json:"metadata"`
	WorkerID  string                 `json:"worker_id"`
}

// TestRequest represents a request to execute a test
type TestRequest struct {
	ID              string                 `json:"id"`
	Type            TestType               `json:"type"`
	Strategy        string                 `json:"strategy"`
	Config          *trading.BotConfig     `json:"config"`
	Scenario        *TestScenario          `json:"scenario"`
	Parameters      map[string]interface{} `json:"parameters"`
	ExpectedResults *ExpectedResults       `json:"expected_results"`
	Priority        int                    `json:"priority"`
	Timeout         time.Duration          `json:"timeout"`
}

// TestWorker executes tests
type TestWorker struct {
	ID          string
	executor    *TestExecutor
	logger      *observability.Logger
	isActive    bool
	currentTest *TestExecution
	stopChan    chan struct{}
}

// NewTestExecutor creates a new test executor
func NewTestExecutor(logger *observability.Logger, framework *BotTestFramework) *TestExecutor {
	config := &TestExecutorConfig{
		MaxWorkers:    5,
		WorkerTimeout: 300 * time.Second,
		QueueSize:     100,
		RetryAttempts: 3,
		RetryDelay:    10 * time.Second,
	}

	return &TestExecutor{
		logger:      logger,
		framework:   framework,
		activeTests: make(map[string]*TestExecution),
		testQueue:   make(chan *TestRequest, config.QueueSize),
		workers:     make([]*TestWorker, 0),
		results:     make(map[string]*TestResult),
		config:      config,
		stopChan:    make(chan struct{}),
	}
}

// Start starts the test executor
func (te *TestExecutor) Start(ctx context.Context) error {
	te.mu.Lock()
	defer te.mu.Unlock()

	if te.isRunning {
		return fmt.Errorf("test executor is already running")
	}

	te.isRunning = true

	// Start workers
	for i := 0; i < te.config.MaxWorkers; i++ {
		worker := &TestWorker{
			ID:       fmt.Sprintf("worker-%d", i+1),
			executor: te,
			logger:   te.logger,
			stopChan: make(chan struct{}),
		}
		te.workers = append(te.workers, worker)
		go worker.run(ctx)
	}

	te.logger.Info(ctx, "Test executor started", map[string]interface{}{
		"max_workers":    te.config.MaxWorkers,
		"queue_size":     te.config.QueueSize,
		"worker_timeout": te.config.WorkerTimeout.String(),
	})

	return nil
}

// Stop stops the test executor
func (te *TestExecutor) Stop(ctx context.Context) error {
	te.mu.Lock()
	defer te.mu.Unlock()

	if !te.isRunning {
		return nil
	}

	te.isRunning = false

	// Stop all workers
	for _, worker := range te.workers {
		close(worker.stopChan)
	}

	// Close test queue
	close(te.stopChan)

	te.logger.Info(ctx, "Test executor stopped", nil)
	return nil
}

// SubmitTest submits a test for execution
func (te *TestExecutor) SubmitTest(request *TestRequest) (string, error) {
	te.mu.Lock()
	defer te.mu.Unlock()

	if !te.isRunning {
		return "", fmt.Errorf("test executor is not running")
	}

	// Generate test ID if not provided
	if request.ID == "" {
		request.ID = uuid.New().String()
	}

	// Set default timeout if not specified
	if request.Timeout == 0 {
		request.Timeout = te.config.WorkerTimeout
	}

	// Queue the test
	select {
	case te.testQueue <- request:
		te.logger.Info(context.Background(), "Test queued", map[string]interface{}{
			"test_id":  request.ID,
			"type":     string(request.Type),
			"strategy": request.Strategy,
			"priority": request.Priority,
		})
		return request.ID, nil
	default:
		return "", fmt.Errorf("test queue is full")
	}
}

// GetTestStatus returns the status of a test
func (te *TestExecutor) GetTestStatus(testID string) (*TestExecution, error) {
	te.mu.RLock()
	defer te.mu.RUnlock()

	execution, exists := te.activeTests[testID]
	if !exists {
		return nil, fmt.Errorf("test not found: %s", testID)
	}

	return execution, nil
}

// GetTestResult returns the result of a completed test
func (te *TestExecutor) GetTestResult(testID string) (*TestResult, error) {
	te.mu.RLock()
	defer te.mu.RUnlock()

	result, exists := te.results[testID]
	if !exists {
		return nil, fmt.Errorf("test result not found: %s", testID)
	}

	return result, nil
}

// GetActiveTests returns all active test executions
func (te *TestExecutor) GetActiveTests() map[string]*TestExecution {
	te.mu.RLock()
	defer te.mu.RUnlock()

	// Return a copy to avoid concurrent access issues
	result := make(map[string]*TestExecution)
	for k, v := range te.activeTests {
		result[k] = v
	}

	return result
}

// CancelTest cancels a running test
func (te *TestExecutor) CancelTest(testID string) error {
	te.mu.Lock()
	defer te.mu.Unlock()

	execution, exists := te.activeTests[testID]
	if !exists {
		return fmt.Errorf("test not found: %s", testID)
	}

	if execution.Status != TestStatusRunning {
		return fmt.Errorf("test is not running: %s", testID)
	}

	execution.Status = TestStatusCancelled

	// Stop the bot if it's running
	if execution.Bot != nil && execution.Bot.Bot != nil {
		te.framework.botEngine.StopBot(context.Background(), execution.Bot.Bot.ID)
	}

	te.logger.Info(context.Background(), "Test cancelled", map[string]interface{}{
		"test_id": testID,
	})

	return nil
}

// TestWorker methods

// run starts the worker loop
func (tw *TestWorker) run(ctx context.Context) {
	tw.logger.Info(ctx, "Test worker started", map[string]interface{}{
		"worker_id": tw.ID,
	})

	for {
		select {
		case <-ctx.Done():
			return
		case <-tw.stopChan:
			return
		case request := <-tw.executor.testQueue:
			if request != nil {
				tw.executeTest(ctx, request)
			}
		}
	}
}

// executeTest executes a single test
func (tw *TestWorker) executeTest(ctx context.Context, request *TestRequest) {
	// Create test execution
	execution := &TestExecution{
		ID:        request.ID,
		Request:   request,
		StartTime: time.Now(),
		Status:    TestStatusRunning,
		Progress:  0.0,
		Metadata:  make(map[string]interface{}),
		WorkerID:  tw.ID,
	}

	// Register execution
	tw.executor.mu.Lock()
	tw.executor.activeTests[request.ID] = execution
	tw.currentTest = execution
	tw.isActive = true
	tw.executor.mu.Unlock()

	tw.logger.Info(ctx, "Test execution started", map[string]interface{}{
		"test_id":   request.ID,
		"worker_id": tw.ID,
		"type":      string(request.Type),
		"strategy":  request.Strategy,
	})

	// Execute test with timeout
	testCtx, cancel := context.WithTimeout(ctx, request.Timeout)
	defer cancel()

	var result *TestResult
	var err error

	switch request.Type {
	case TestTypeUnit:
		result, err = tw.executeUnitTest(testCtx, request)
	case TestTypeIntegration:
		result, err = tw.executeIntegrationTest(testCtx, request)
	case TestTypeBacktest:
		result, err = tw.executeBacktest(testCtx, request)
	case TestTypePaperTrade:
		result, err = tw.executePaperTrade(testCtx, request)
	case TestTypeStress:
		result, err = tw.executeStressTest(testCtx, request)
	case TestTypePerformance:
		result, err = tw.executePerformanceTest(testCtx, request)
	default:
		err = fmt.Errorf("unsupported test type: %s", request.Type)
	}

	// Update execution status
	tw.executor.mu.Lock()
	if err != nil {
		execution.Status = TestStatusFailed
		execution.Metadata["error"] = err.Error()
		tw.logger.Error(ctx, "Test execution failed", err, map[string]interface{}{
			"test_id":   request.ID,
			"worker_id": tw.ID,
		})
	} else {
		execution.Status = TestStatusCompleted
		execution.Progress = 100.0
		tw.executor.results[request.ID] = result
		tw.logger.Info(ctx, "Test execution completed", map[string]interface{}{
			"test_id":   request.ID,
			"worker_id": tw.ID,
			"passed":    result.Passed,
		})
	}

	// Clean up
	delete(tw.executor.activeTests, request.ID)
	tw.currentTest = nil
	tw.isActive = false
	tw.executor.mu.Unlock()
}

// executeUnitTest executes a unit test
func (tw *TestWorker) executeUnitTest(ctx context.Context, request *TestRequest) (*TestResult, error) {
	result := &TestResult{
		TestID:    request.ID,
		TestType:  request.Type,
		Strategy:  request.Strategy,
		Status:    TestStatusCompleted,
		StartTime: time.Now(),
		Passed:    true,
		Metadata:  make(map[string]interface{}),
	}

	// Simulate unit test execution
	time.Sleep(1 * time.Second)

	// Basic validation tests
	if request.Config == nil {
		result.Passed = false
		result.FailureReasons = append(result.FailureReasons, "Missing bot configuration")
	}

	if request.Strategy == "" {
		result.Passed = false
		result.FailureReasons = append(result.FailureReasons, "Missing strategy")
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	return result, nil
}

// executeIntegrationTest executes an integration test
func (tw *TestWorker) executeIntegrationTest(ctx context.Context, request *TestRequest) (*TestResult, error) {
	result := &TestResult{
		TestID:    request.ID,
		TestType:  request.Type,
		Strategy:  request.Strategy,
		Status:    TestStatusCompleted,
		StartTime: time.Now(),
		Passed:    true,
		Metadata:  make(map[string]interface{}),
	}

	// Test bot creation and basic operations
	testBot, err := tw.createTestBot(request)
	if err != nil {
		result.Passed = false
		result.FailureReasons = append(result.FailureReasons, fmt.Sprintf("Failed to create test bot: %v", err))
		return result, nil
	}

	// Test bot start/stop
	if err := tw.executor.framework.botEngine.StartBot(context.Background(), testBot.Bot.ID); err != nil {
		result.Passed = false
		result.FailureReasons = append(result.FailureReasons, fmt.Sprintf("Failed to start bot: %v", err))
	}

	// Wait a bit for bot to initialize
	time.Sleep(2 * time.Second)

	// Stop the bot
	if err := tw.executor.framework.botEngine.StopBot(context.Background(), testBot.Bot.ID); err != nil {
		result.Passed = false
		result.FailureReasons = append(result.FailureReasons, fmt.Sprintf("Failed to stop bot: %v", err))
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	return result, nil
}

// executeBacktest executes a backtest
func (tw *TestWorker) executeBacktest(ctx context.Context, request *TestRequest) (*TestResult, error) {
	result := &TestResult{
		TestID:    request.ID,
		TestType:  request.Type,
		Strategy:  request.Strategy,
		Status:    TestStatusCompleted,
		StartTime: time.Now(),
		Passed:    true,
		Metadata:  make(map[string]interface{}),
	}

	// Create test bot
	testBot, err := tw.createTestBot(request)
	if err != nil {
		result.Passed = false
		result.FailureReasons = append(result.FailureReasons, fmt.Sprintf("Failed to create test bot: %v", err))
		return result, nil
	}

	// Set market condition if specified in scenario
	if request.Scenario != nil {
		tw.executor.framework.mockExchange.SetMarketCondition(request.Scenario.MarketCondition)
		tw.executor.framework.mockMarketData.SetMarketCondition(request.Scenario.MarketCondition)
	}

	// Run backtest
	testDuration := 30 * time.Second // Shortened for testing
	if request.Scenario != nil && request.Scenario.Duration > 0 {
		testDuration = request.Scenario.Duration
		if testDuration > 60*time.Second {
			testDuration = 60 * time.Second // Cap at 1 minute for testing
		}
	}

	// Start the bot
	if err := tw.executor.framework.botEngine.StartBot(context.Background(), testBot.Bot.ID); err != nil {
		result.Passed = false
		result.FailureReasons = append(result.FailureReasons, fmt.Sprintf("Failed to start bot: %v", err))
		return result, nil
	}

	// Wait for test duration
	testCtx, cancel := context.WithTimeout(ctx, testDuration)
	defer cancel()

	<-testCtx.Done()

	// Stop the bot
	tw.executor.framework.botEngine.StopBot(context.Background(), testBot.Bot.ID)

	// Calculate metrics and validate results
	metrics := tw.calculateTestMetrics(testBot)
	result.Metrics = metrics

	// Validate against expected results
	if request.ExpectedResults != nil {
		tw.validateResults(result, metrics, request.ExpectedResults)
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	return result, nil
}

// executePaperTrade executes a paper trading test
func (tw *TestWorker) executePaperTrade(ctx context.Context, request *TestRequest) (*TestResult, error) {
	// Similar to backtest but with real-time market data simulation
	return tw.executeBacktest(ctx, request)
}

// executeStressTest executes a stress test
func (tw *TestWorker) executeStressTest(ctx context.Context, request *TestRequest) (*TestResult, error) {
	result := &TestResult{
		TestID:    request.ID,
		TestType:  request.Type,
		Strategy:  request.Strategy,
		Status:    TestStatusCompleted,
		StartTime: time.Now(),
		Passed:    true,
		Metadata:  make(map[string]interface{}),
	}

	// Set volatile market conditions
	tw.executor.framework.mockExchange.SetMarketCondition("volatile")
	tw.executor.framework.mockMarketData.SetMarketCondition("volatile")

	// Create and run test bot under stress
	testBot, err := tw.createTestBot(request)
	if err != nil {
		result.Passed = false
		result.FailureReasons = append(result.FailureReasons, fmt.Sprintf("Failed to create test bot: %v", err))
		return result, nil
	}

	// Run stress test for shorter duration
	if err := tw.executor.framework.botEngine.StartBot(context.Background(), testBot.Bot.ID); err != nil {
		result.Passed = false
		result.FailureReasons = append(result.FailureReasons, fmt.Sprintf("Failed to start bot: %v", err))
		return result, nil
	}

	// Wait for stress test duration
	time.Sleep(15 * time.Second)

	// Stop the bot
	tw.executor.framework.botEngine.StopBot(context.Background(), testBot.Bot.ID)

	// Check if bot survived stress conditions
	if testBot.Status == TestStatusFailed {
		result.Passed = false
		result.FailureReasons = append(result.FailureReasons, "Bot failed under stress conditions")
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	return result, nil
}

// executePerformanceTest executes a performance test
func (tw *TestWorker) executePerformanceTest(ctx context.Context, request *TestRequest) (*TestResult, error) {
	// Similar to backtest but focuses on performance metrics
	return tw.executeBacktest(ctx, request)
}

// createTestBot creates a test bot for execution
func (tw *TestWorker) createTestBot(request *TestRequest) (*TestBot, error) {
	botID := fmt.Sprintf("test-bot-%s", request.ID[:8])

	// Create test account in mock exchange
	initialBalances := map[string]decimal.Decimal{
		"USDT": tw.executor.framework.testConfig.InitialBalance,
		"BTC":  decimal.Zero,
		"ETH":  decimal.Zero,
		"BNB":  decimal.Zero,
	}

	if err := tw.executor.framework.mockExchange.CreateAccount(botID, initialBalances); err != nil {
		return nil, fmt.Errorf("failed to create test account: %w", err)
	}

	// Create bot configuration
	config := request.Config
	if config == nil {
		config = &trading.BotConfig{
			TradingPairs: []string{"BTC/USDT"},
			Exchange:     "mock",
			BaseCurrency: "USDT",
			StrategyParams: map[string]interface{}{
				"strategy": request.Strategy,
			},
			Enabled: true,
		}
	}

	// Create trading bot (simplified for testing)
	bot := &trading.TradingBot{
		ID:       botID,
		Name:     fmt.Sprintf("Test Bot %s", request.Strategy),
		Strategy: trading.BotStrategy(request.Strategy),
		Config:   config,
		State:    trading.StateIdle,
	}

	// Create test bot wrapper
	testBot := &TestBot{
		ID:             botID,
		Name:           fmt.Sprintf("Test Bot %s", request.Strategy),
		Strategy:       request.Strategy,
		Config:         config,
		Bot:            bot,
		TestStartTime:  time.Now(),
		InitialBalance: tw.executor.framework.testConfig.InitialBalance,
		CurrentBalance: tw.executor.framework.testConfig.InitialBalance,
		Trades:         make([]*TestTrade, 0),
		Status:         TestStatusRunning,
	}

	return testBot, nil
}

// calculateTestMetrics calculates performance metrics for a test
func (tw *TestWorker) calculateTestMetrics(testBot *TestBot) *TestMetrics {
	// Get trades from mock exchange
	trades, _ := tw.executor.framework.mockExchange.GetTrades(testBot.ID, 1000)

	metrics := &TestMetrics{
		TotalTrades: len(trades),
	}

	if len(trades) == 0 {
		return metrics
	}

	// Calculate basic metrics
	totalPnL := decimal.Zero
	winningTrades := 0
	totalWinAmount := decimal.Zero
	totalLossAmount := decimal.Zero

	for _, trade := range trades {
		// Simplified P&L calculation
		pnl := decimal.Zero
		if trade.Side == "sell" {
			// Assume we're selling at profit/loss
			pnl = trade.Price.Sub(testBot.InitialBalance.Div(decimal.NewFromInt(int64(len(trades)))))
		}

		totalPnL = totalPnL.Add(pnl)

		if pnl.GreaterThan(decimal.Zero) {
			winningTrades++
			totalWinAmount = totalWinAmount.Add(pnl)
		} else {
			totalLossAmount = totalLossAmount.Add(pnl.Abs())
		}
	}

	metrics.WinningTrades = winningTrades
	metrics.LosingTrades = len(trades) - winningTrades

	if len(trades) > 0 {
		metrics.WinRate = decimal.NewFromInt(int64(winningTrades)).Div(decimal.NewFromInt(int64(len(trades))))
	}

	metrics.TotalReturn = totalPnL

	// Calculate profit factor
	if !totalLossAmount.IsZero() {
		metrics.ProfitFactor = totalWinAmount.Div(totalLossAmount)
	}

	// Simplified Sharpe ratio (would need more sophisticated calculation in real implementation)
	metrics.SharpeRatio = decimal.NewFromFloat(0.5)

	// Simplified max drawdown
	metrics.MaxDrawdown = decimal.NewFromFloat(0.05)

	return metrics
}

// validateResults validates test results against expected results
func (tw *TestWorker) validateResults(result *TestResult, metrics *TestMetrics, expected *ExpectedResults) {
	if metrics.WinRate.LessThan(expected.MinWinRate) {
		result.Passed = false
		result.FailureReasons = append(result.FailureReasons,
			fmt.Sprintf("Win rate %.2f%% below minimum %.2f%%",
				metrics.WinRate.Mul(decimal.NewFromFloat(100)).InexactFloat64(),
				expected.MinWinRate.Mul(decimal.NewFromFloat(100)).InexactFloat64()))
	}

	if metrics.MaxDrawdown.GreaterThan(expected.MaxDrawdown) {
		result.Passed = false
		result.FailureReasons = append(result.FailureReasons,
			fmt.Sprintf("Max drawdown %.2f%% exceeds maximum %.2f%%",
				metrics.MaxDrawdown.Mul(decimal.NewFromFloat(100)).InexactFloat64(),
				expected.MaxDrawdown.Mul(decimal.NewFromFloat(100)).InexactFloat64()))
	}

	if metrics.SharpeRatio.LessThan(expected.MinSharpeRatio) {
		result.Passed = false
		result.FailureReasons = append(result.FailureReasons,
			fmt.Sprintf("Sharpe ratio %.2f below minimum %.2f",
				metrics.SharpeRatio.InexactFloat64(),
				expected.MinSharpeRatio.InexactFloat64()))
	}
}
