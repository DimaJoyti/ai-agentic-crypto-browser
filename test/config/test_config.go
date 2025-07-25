package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// TestConfig contains all test configuration
type TestConfig struct {
	// Database configuration
	Database DatabaseTestConfig `json:"database"`

	// Redis configuration
	Redis RedisTestConfig `json:"redis"`

	// HTTP configuration
	HTTP HTTPTestConfig `json:"http"`

	// Browser testing configuration
	Browser BrowserTestConfig `json:"browser"`

	// Load testing configuration
	LoadTest LoadTestConfig `json:"load_test"`

	// E2E testing configuration
	E2E E2ETestConfig `json:"e2e"`

	// General test configuration
	General GeneralTestConfig `json:"general"`
}

// DatabaseTestConfig contains database test configuration
type DatabaseTestConfig struct {
	UseTestContainers bool   `json:"use_test_containers"`
	Host              string `json:"host"`
	Port              int    `json:"port"`
	Name              string `json:"name"`
	User              string `json:"user"`
	Password          string `json:"password"`
	SSLMode           string `json:"ssl_mode"`
	MaxConnections    int    `json:"max_connections"`
	CleanupAfterTest  bool   `json:"cleanup_after_test"`
}

// RedisTestConfig contains Redis test configuration
type RedisTestConfig struct {
	UseTestContainers bool   `json:"use_test_containers"`
	Host              string `json:"host"`
	Port              int    `json:"port"`
	Password          string `json:"password"`
	DB                int    `json:"db"`
	CleanupAfterTest  bool   `json:"cleanup_after_test"`
}

// HTTPTestConfig contains HTTP test configuration
type HTTPTestConfig struct {
	BaseURL         string        `json:"base_url"`
	Timeout         time.Duration `json:"timeout"`
	RetryAttempts   int           `json:"retry_attempts"`
	RetryDelay      time.Duration `json:"retry_delay"`
	FollowRedirects bool          `json:"follow_redirects"`
}

// BrowserTestConfig contains browser test configuration
type BrowserTestConfig struct {
	Headless          bool          `json:"headless"`
	WindowWidth       int           `json:"window_width"`
	WindowHeight      int           `json:"window_height"`
	Timeout           time.Duration `json:"timeout"`
	SlowMo            time.Duration `json:"slow_mo"`
	ScreenshotOnError bool          `json:"screenshot_on_error"`
	VideoRecording    bool          `json:"video_recording"`
	BrowserType       string        `json:"browser_type"` // chrome, firefox, safari
}

// LoadTestConfig contains load test configuration
type LoadTestConfig struct {
	Enabled           bool          `json:"enabled"`
	ConcurrentUsers   int           `json:"concurrent_users"`
	TestDuration      time.Duration `json:"test_duration"`
	RampUpTime        time.Duration `json:"ramp_up_time"`
	RequestsPerSecond int           `json:"requests_per_second"`
	MaxResponseTime   time.Duration `json:"max_response_time"`
	MinSuccessRate    float64       `json:"min_success_rate"`
	ReportFormat      string        `json:"report_format"` // json, html, csv
}

// E2ETestConfig contains E2E test configuration
type E2ETestConfig struct {
	Enabled           bool          `json:"enabled"`
	BaseURL           string        `json:"base_url"`
	Timeout           time.Duration `json:"timeout"`
	RetryAttempts     int           `json:"retry_attempts"`
	ParallelExecution bool          `json:"parallel_execution"`
	TestDataCleanup   bool          `json:"test_data_cleanup"`
}

// GeneralTestConfig contains general test configuration
type GeneralTestConfig struct {
	LogLevel          string        `json:"log_level"`
	EnableTracing     bool          `json:"enable_tracing"`
	EnableMetrics     bool          `json:"enable_metrics"`
	TestTimeout       time.Duration `json:"test_timeout"`
	ParallelTests     bool          `json:"parallel_tests"`
	FailFast          bool          `json:"fail_fast"`
	Verbose           bool          `json:"verbose"`
	CoverageThreshold float64       `json:"coverage_threshold"`
	ReportPath        string        `json:"report_path"`
}

// GetTestConfig returns the test configuration based on environment variables
func GetTestConfig() *TestConfig {
	return &TestConfig{
		Database: DatabaseTestConfig{
			UseTestContainers: getBoolEnv("TEST_USE_CONTAINERS", true),
			Host:              getStringEnv("TEST_DB_HOST", "localhost"),
			Port:              getIntEnv("TEST_DB_PORT", 5432),
			Name:              getStringEnv("TEST_DB_NAME", "test_db"),
			User:              getStringEnv("TEST_DB_USER", "test_user"),
			Password:          getStringEnv("TEST_DB_PASSWORD", "test_password"),
			SSLMode:           getStringEnv("TEST_DB_SSL_MODE", "disable"),
			MaxConnections:    getIntEnv("TEST_DB_MAX_CONNECTIONS", 10),
			CleanupAfterTest:  getBoolEnv("TEST_DB_CLEANUP", true),
		},
		Redis: RedisTestConfig{
			UseTestContainers: getBoolEnv("TEST_USE_CONTAINERS", true),
			Host:              getStringEnv("TEST_REDIS_HOST", "localhost"),
			Port:              getIntEnv("TEST_REDIS_PORT", 6379),
			Password:          getStringEnv("TEST_REDIS_PASSWORD", ""),
			DB:                getIntEnv("TEST_REDIS_DB", 1),
			CleanupAfterTest:  getBoolEnv("TEST_REDIS_CLEANUP", true),
		},
		HTTP: HTTPTestConfig{
			BaseURL:         getStringEnv("TEST_HTTP_BASE_URL", "http://localhost:8080"),
			Timeout:         getDurationEnv("TEST_HTTP_TIMEOUT", 30*time.Second),
			RetryAttempts:   getIntEnv("TEST_HTTP_RETRY_ATTEMPTS", 3),
			RetryDelay:      getDurationEnv("TEST_HTTP_RETRY_DELAY", 1*time.Second),
			FollowRedirects: getBoolEnv("TEST_HTTP_FOLLOW_REDIRECTS", true),
		},
		Browser: BrowserTestConfig{
			Headless:          getBoolEnv("TEST_BROWSER_HEADLESS", true),
			WindowWidth:       getIntEnv("TEST_BROWSER_WIDTH", 1920),
			WindowHeight:      getIntEnv("TEST_BROWSER_HEIGHT", 1080),
			Timeout:           getDurationEnv("TEST_BROWSER_TIMEOUT", 30*time.Second),
			SlowMo:            getDurationEnv("TEST_BROWSER_SLOW_MO", 0),
			ScreenshotOnError: getBoolEnv("TEST_BROWSER_SCREENSHOT_ON_ERROR", true),
			VideoRecording:    getBoolEnv("TEST_BROWSER_VIDEO_RECORDING", false),
			BrowserType:       getStringEnv("TEST_BROWSER_TYPE", "chrome"),
		},
		LoadTest: LoadTestConfig{
			Enabled:           getBoolEnv("TEST_LOAD_ENABLED", false),
			ConcurrentUsers:   getIntEnv("TEST_LOAD_CONCURRENT_USERS", 10),
			TestDuration:      getDurationEnv("TEST_LOAD_DURATION", 60*time.Second),
			RampUpTime:        getDurationEnv("TEST_LOAD_RAMP_UP", 10*time.Second),
			RequestsPerSecond: getIntEnv("TEST_LOAD_RPS", 10),
			MaxResponseTime:   getDurationEnv("TEST_LOAD_MAX_RESPONSE_TIME", 5*time.Second),
			MinSuccessRate:    getFloatEnv("TEST_LOAD_MIN_SUCCESS_RATE", 95.0),
			ReportFormat:      getStringEnv("TEST_LOAD_REPORT_FORMAT", "json"),
		},
		E2E: E2ETestConfig{
			Enabled:           getBoolEnv("TEST_E2E_ENABLED", true),
			BaseURL:           getStringEnv("TEST_E2E_BASE_URL", "http://localhost:3000"),
			Timeout:           getDurationEnv("TEST_E2E_TIMEOUT", 60*time.Second),
			RetryAttempts:     getIntEnv("TEST_E2E_RETRY_ATTEMPTS", 2),
			ParallelExecution: getBoolEnv("TEST_E2E_PARALLEL", false),
			TestDataCleanup:   getBoolEnv("TEST_E2E_CLEANUP", true),
		},
		General: GeneralTestConfig{
			LogLevel:          getStringEnv("TEST_LOG_LEVEL", "info"),
			EnableTracing:     getBoolEnv("TEST_ENABLE_TRACING", false),
			EnableMetrics:     getBoolEnv("TEST_ENABLE_METRICS", false),
			TestTimeout:       getDurationEnv("TEST_TIMEOUT", 300*time.Second),
			ParallelTests:     getBoolEnv("TEST_PARALLEL", true),
			FailFast:          getBoolEnv("TEST_FAIL_FAST", false),
			Verbose:           getBoolEnv("TEST_VERBOSE", false),
			CoverageThreshold: getFloatEnv("TEST_COVERAGE_THRESHOLD", 80.0),
			ReportPath:        getStringEnv("TEST_REPORT_PATH", "./test-reports"),
		},
	}
}

// Environment variable helper functions

func getStringEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getFloatEnv(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// TestEnvironment represents different test environments
type TestEnvironment string

const (
	TestEnvLocal      TestEnvironment = "local"
	TestEnvCI         TestEnvironment = "ci"
	TestEnvStaging    TestEnvironment = "staging"
	TestEnvProduction TestEnvironment = "production"
)

// GetTestEnvironment returns the current test environment
func GetTestEnvironment() TestEnvironment {
	env := getStringEnv("TEST_ENVIRONMENT", "local")
	return TestEnvironment(env)
}

// IsCI returns true if running in CI environment
func IsCI() bool {
	return getBoolEnv("CI", false) ||
		getBoolEnv("GITHUB_ACTIONS", false) ||
		getBoolEnv("GITLAB_CI", false) ||
		getBoolEnv("JENKINS_URL", false) != false
}

// GetCIProvider returns the CI provider name
func GetCIProvider() string {
	if getBoolEnv("GITHUB_ACTIONS", false) {
		return "github"
	}
	if getBoolEnv("GITLAB_CI", false) {
		return "gitlab"
	}
	if getStringEnv("JENKINS_URL", "") != "" {
		return "jenkins"
	}
	if getStringEnv("CIRCLECI", "") != "" {
		return "circleci"
	}
	return "unknown"
}

// TestCategories defines test categories for selective execution
type TestCategory string

const (
	CategoryUnit        TestCategory = "unit"
	CategoryIntegration TestCategory = "integration"
	CategoryE2E         TestCategory = "e2e"
	CategoryLoad        TestCategory = "load"
	CategorySecurity    TestCategory = "security"
	CategorySmoke       TestCategory = "smoke"
)

// ShouldRunCategory determines if a test category should run based on environment
func ShouldRunCategory(category TestCategory) bool {
	env := GetTestEnvironment()

	switch env {
	case TestEnvLocal:
		// Run all tests locally except load tests by default
		return category != CategoryLoad
	case TestEnvCI:
		// Run unit, integration, and smoke tests in CI
		return category == CategoryUnit ||
			category == CategoryIntegration ||
			category == CategorySmoke
	case TestEnvStaging:
		// Run all tests in staging
		return true
	case TestEnvProduction:
		// Only run smoke tests in production
		return category == CategorySmoke
	default:
		return true
	}
}

// TestTags contains test tags for categorization
type TestTags struct {
	Category  TestCategory `json:"category"`
	Service   string       `json:"service"`
	Component string       `json:"component"`
	Priority  string       `json:"priority"` // high, medium, low
	Flaky     bool         `json:"flaky"`
	Slow      bool         `json:"slow"`
	External  bool         `json:"external"` // requires external services
}

// DatabaseTestHelper provides database testing utilities
type DatabaseTestHelper struct {
	Config *DatabaseTestConfig
}

// NewDatabaseTestHelper creates a new database test helper
func NewDatabaseTestHelper(config *DatabaseTestConfig) *DatabaseTestHelper {
	return &DatabaseTestHelper{
		Config: config,
	}
}

// GetConnectionString returns the database connection string for testing
func (h *DatabaseTestHelper) GetConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		h.Config.User,
		h.Config.Password,
		h.Config.Host,
		h.Config.Port,
		h.Config.Name,
		h.Config.SSLMode,
	)
}

// RedisTestHelper provides Redis testing utilities
type RedisTestHelper struct {
	Config *RedisTestConfig
}

// NewRedisTestHelper creates a new Redis test helper
func NewRedisTestHelper(config *RedisTestConfig) *RedisTestHelper {
	return &RedisTestHelper{
		Config: config,
	}
}

// GetConnectionString returns the Redis connection string for testing
func (h *RedisTestHelper) GetConnectionString() string {
	if h.Config.Password != "" {
		return fmt.Sprintf("redis://:%s@%s:%d/%d",
			h.Config.Password,
			h.Config.Host,
			h.Config.Port,
			h.Config.DB,
		)
	}
	return fmt.Sprintf("redis://%s:%d/%d",
		h.Config.Host,
		h.Config.Port,
		h.Config.DB,
	)
}

// TestReporter handles test result reporting
type TestReporter struct {
	Config *GeneralTestConfig
}

// NewTestReporter creates a new test reporter
func NewTestReporter(config *GeneralTestConfig) *TestReporter {
	return &TestReporter{
		Config: config,
	}
}

// ReportTestResults reports test results in the configured format
func (r *TestReporter) ReportTestResults(results interface{}) error {
	// Implementation would depend on the reporting format
	// Could generate JUnit XML, HTML reports, JSON reports, etc.
	return nil
}

// TestDataManager handles test data management
type TestDataManager struct {
	Config *TestConfig
}

// NewTestDataManager creates a new test data manager
func NewTestDataManager(config *TestConfig) *TestDataManager {
	return &TestDataManager{
		Config: config,
	}
}

// SetupTestData sets up test data for testing
func (m *TestDataManager) SetupTestData() error {
	// Implementation would set up test fixtures, seed data, etc.
	return nil
}

// CleanupTestData cleans up test data after testing
func (m *TestDataManager) CleanupTestData() error {
	// Implementation would clean up test data
	return nil
}
