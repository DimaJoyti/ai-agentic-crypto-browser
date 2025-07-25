package testing

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.opentelemetry.io/otel/trace"

	"github.com/ai-agentic-browser/internal/config"
	"github.com/ai-agentic-browser/pkg/observability"
)

// TestSuite provides a base test suite with common testing utilities
type TestSuite struct {
	suite.Suite

	// Test infrastructure
	DB         *sql.DB
	Redis      *redis.Client
	HTTPServer *httptest.Server
	Router     *gin.Engine

	// Test containers
	PostgresContainer testcontainers.Container
	RedisContainer    testcontainers.Container

	// Test configuration
	Config *TestConfig
	Logger *observability.Logger
	Tracer trace.Tracer

	// Test context
	Ctx        context.Context
	CancelFunc context.CancelFunc
}

// TestConfig contains configuration for test setup
type TestConfig struct {
	Database DatabaseTestConfig
	Redis    RedisTestConfig
	Testing  TestingConfig
}

type DatabaseTestConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}

type RedisTestConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type TestingConfig struct {
	UseTestContainers bool
	EnableTracing     bool
	LogLevel          string
}

// SetupSuite initializes the test suite
func (ts *TestSuite) SetupSuite() {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create test context
	ts.Ctx, ts.CancelFunc = context.WithCancel(context.Background())

	// Initialize test configuration
	ts.initializeConfig()

	// Initialize logging
	ts.initializeLogging()

	// Initialize tracing if enabled
	if ts.Config.Testing.EnableTracing {
		ts.initializeTracing()
	}

	// Setup test infrastructure
	ts.setupInfrastructure()

	// Setup HTTP router
	ts.setupHTTPRouter()
}

// TearDownSuite cleans up the test suite
func (ts *TestSuite) TearDownSuite() {
	// Close database connection
	if ts.DB != nil {
		ts.DB.Close()
	}

	// Close Redis connection
	if ts.Redis != nil {
		ts.Redis.Close()
	}

	// Stop HTTP server
	if ts.HTTPServer != nil {
		ts.HTTPServer.Close()
	}

	// Stop test containers
	if ts.PostgresContainer != nil {
		ts.PostgresContainer.Terminate(ts.Ctx)
	}
	if ts.RedisContainer != nil {
		ts.RedisContainer.Terminate(ts.Ctx)
	}

	// Cancel context
	if ts.CancelFunc != nil {
		ts.CancelFunc()
	}
}

// SetupTest runs before each test
func (ts *TestSuite) SetupTest() {
	// Clean database
	ts.cleanDatabase()

	// Clean Redis
	ts.cleanRedis()
}

// TearDownTest runs after each test
func (ts *TestSuite) TearDownTest() {
	// Additional cleanup if needed
}

// initializeConfig sets up test configuration
func (ts *TestSuite) initializeConfig() {
	ts.Config = &TestConfig{
		Database: DatabaseTestConfig{
			Host:     "localhost",
			Port:     5432,
			Name:     "test_db",
			User:     "test_user",
			Password: "test_password",
		},
		Redis: RedisTestConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       1, // Use different DB for tests
		},
		Testing: TestingConfig{
			UseTestContainers: true,
			EnableTracing:     false,
			LogLevel:          "debug",
		},
	}
}

// initializeLogging sets up test logging
func (ts *TestSuite) initializeLogging() {
	logConfig := config.ObservabilityConfig{
		ServiceName: "test-service",
		LogLevel:    ts.Config.Testing.LogLevel,
		LogFormat:   "json",
	}

	ts.Logger = observability.NewLogger(logConfig)
}

// initializeTracing sets up test tracing
func (ts *TestSuite) initializeTracing() {
	// Initialize minimal tracing for tests
	// Implementation would depend on your tracing setup
}

// setupInfrastructure sets up test databases and services
func (ts *TestSuite) setupInfrastructure() {
	if ts.Config.Testing.UseTestContainers {
		ts.setupTestContainers()
	} else {
		ts.setupLocalServices()
	}
}

// setupTestContainers sets up test containers for isolated testing
func (ts *TestSuite) setupTestContainers() {
	// Setup PostgreSQL container
	postgresReq := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       ts.Config.Database.Name,
			"POSTGRES_USER":     ts.Config.Database.User,
			"POSTGRES_PASSWORD": ts.Config.Database.Password,
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	var err error
	ts.PostgresContainer, err = testcontainers.GenericContainer(ts.Ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: postgresReq,
		Started:          true,
	})
	require.NoError(ts.T(), err)

	// Get PostgreSQL connection details
	host, err := ts.PostgresContainer.Host(ts.Ctx)
	require.NoError(ts.T(), err)

	port, err := ts.PostgresContainer.MappedPort(ts.Ctx, "5432")
	require.NoError(ts.T(), err)

	// Setup database connection
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		ts.Config.Database.User,
		ts.Config.Database.Password,
		host,
		port.Port(),
		ts.Config.Database.Name,
	)

	ts.DB, err = sql.Open("postgres", dsn)
	require.NoError(ts.T(), err)

	// Wait for database to be ready
	require.Eventually(ts.T(), func() bool {
		return ts.DB.Ping() == nil
	}, 30*time.Second, 1*time.Second)

	// Setup Redis container
	redisReq := testcontainers.ContainerRequest{
		Image:        "redis:7-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp"),
	}

	ts.RedisContainer, err = testcontainers.GenericContainer(ts.Ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: redisReq,
		Started:          true,
	})
	require.NoError(ts.T(), err)

	// Get Redis connection details
	redisHost, err := ts.RedisContainer.Host(ts.Ctx)
	require.NoError(ts.T(), err)

	redisPort, err := ts.RedisContainer.MappedPort(ts.Ctx, "6379")
	require.NoError(ts.T(), err)

	// Setup Redis connection
	ts.Redis = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort.Port()),
		DB:   ts.Config.Redis.DB,
	})

	// Test Redis connection
	require.Eventually(ts.T(), func() bool {
		return ts.Redis.Ping(ts.Ctx).Err() == nil
	}, 30*time.Second, 1*time.Second)
}

// setupLocalServices connects to local services for testing
func (ts *TestSuite) setupLocalServices() {
	// Setup local database connection
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		ts.Config.Database.User,
		ts.Config.Database.Password,
		ts.Config.Database.Host,
		ts.Config.Database.Port,
		ts.Config.Database.Name,
	)

	var err error
	ts.DB, err = sql.Open("postgres", dsn)
	require.NoError(ts.T(), err)

	// Setup local Redis connection
	ts.Redis = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", ts.Config.Redis.Host, ts.Config.Redis.Port),
		DB:   ts.Config.Redis.DB,
	})
}

// setupHTTPRouter sets up test HTTP router
func (ts *TestSuite) setupHTTPRouter() {
	ts.Router = gin.New()

	// Add test middleware
	ts.Router.Use(gin.Recovery())
	ts.Router.Use(func(c *gin.Context) {
		c.Set("test_context", ts.Ctx)
		c.Set("test_db", ts.DB)
		c.Set("test_redis", ts.Redis)
		c.Set("test_logger", ts.Logger)
		c.Next()
	})

	// Create test server
	ts.HTTPServer = httptest.NewServer(ts.Router)
}

// cleanDatabase cleans the test database
func (ts *TestSuite) cleanDatabase() {
	if ts.DB == nil {
		return
	}

	// Get all table names
	rows, err := ts.DB.Query(`
		SELECT tablename FROM pg_tables 
		WHERE schemaname = 'public' 
		AND tablename NOT LIKE 'pg_%'
	`)
	require.NoError(ts.T(), err)
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		require.NoError(ts.T(), err)
		tables = append(tables, table)
	}

	// Truncate all tables
	for _, table := range tables {
		_, err := ts.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		require.NoError(ts.T(), err)
	}
}

// cleanRedis cleans the test Redis database
func (ts *TestSuite) cleanRedis() {
	if ts.Redis == nil {
		return
	}

	err := ts.Redis.FlushDB(ts.Ctx).Err()
	require.NoError(ts.T(), err)
}

// Helper methods for common test operations

// CreateTestUser creates a test user in the database
func (ts *TestSuite) CreateTestUser(email, password string) uuid.UUID {
	userID := uuid.New()

	query := `
		INSERT INTO users (id, email, password_hash, first_name, last_name, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := ts.DB.Exec(query,
		userID, email, password, "Test", "User", true, time.Now(), time.Now(),
	)
	require.NoError(ts.T(), err)

	return userID
}

// CreateTestWorkflow creates a test workflow in the database
func (ts *TestSuite) CreateTestWorkflow(userID uuid.UUID, name string) uuid.UUID {
	workflowID := uuid.New()

	query := `
		INSERT INTO workflows (id, user_id, name, description, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := ts.DB.Exec(query,
		workflowID, userID, name, "Test workflow", "active", time.Now(), time.Now(),
	)
	require.NoError(ts.T(), err)

	return workflowID
}

// MakeHTTPRequest makes an HTTP request to the test server
func (ts *TestSuite) MakeHTTPRequest(method, path string, body interface{}, headers map[string]string) *http.Response {
	req := ts.createHTTPRequest(method, path, body)

	// Add headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(ts.T(), err)

	return resp
}

// createHTTPRequest creates an HTTP request for testing
func (ts *TestSuite) createHTTPRequest(method, path string, body interface{}) *http.Request {
	var req *http.Request
	var err error

	if body != nil {
		// Implementation would depend on body type (JSON, form data, etc.)
		req, err = http.NewRequest(method, ts.HTTPServer.URL+path, nil)
	} else {
		req, err = http.NewRequest(method, ts.HTTPServer.URL+path, nil)
	}

	require.NoError(ts.T(), err)
	return req
}

// AssertHTTPStatus asserts HTTP response status
func (ts *TestSuite) AssertHTTPStatus(resp *http.Response, expectedStatus int) {
	assert.Equal(ts.T(), expectedStatus, resp.StatusCode)
}

// AssertDatabaseRecord asserts that a database record exists
func (ts *TestSuite) AssertDatabaseRecord(table string, conditions map[string]interface{}) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE ", table)

	var args []interface{}
	var whereClauses []string
	i := 1

	for column, value := range conditions {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", column, i))
		args = append(args, value)
		i++
	}

	query += fmt.Sprintf("%s", whereClauses[0])
	for _, clause := range whereClauses[1:] {
		query += " AND " + clause
	}

	var count int
	err := ts.DB.QueryRow(query, args...).Scan(&count)
	require.NoError(ts.T(), err)
	assert.Greater(ts.T(), count, 0, "Expected database record not found")
}

// AssertRedisKey asserts that a Redis key exists
func (ts *TestSuite) AssertRedisKey(key string) {
	exists, err := ts.Redis.Exists(ts.Ctx, key).Result()
	require.NoError(ts.T(), err)
	assert.Equal(ts.T(), int64(1), exists, "Expected Redis key not found")
}
