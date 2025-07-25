//go:build integration

package integration

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseConnection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	dbHost := getEnv("TEST_DB_HOST", "localhost")
	dbPort := getEnv("TEST_DB_PORT", "5432")
	dbName := getEnv("TEST_DB_NAME", "test_db")
	dbUser := getEnv("TEST_DB_USER", "test_user")
	dbPassword := getEnv("TEST_DB_PASSWORD", "test_password")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	assert.NoError(t, err, "Should be able to connect to database")

	// Test basic query
	var result int
	err = db.QueryRowContext(ctx, "SELECT 1").Scan(&result)
	assert.NoError(t, err)
	assert.Equal(t, 1, result)
}

func TestRedisConnection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	redisHost := getEnv("TEST_REDIS_HOST", "localhost")
	redisPort := getEnv("TEST_REDIS_PORT", "6379")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: "", // no password for test
		DB:       0,  // default DB
	})
	defer rdb.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test ping
	pong, err := rdb.Ping(ctx).Result()
	assert.NoError(t, err)
	assert.Equal(t, "PONG", pong)

	// Test set/get
	key := fmt.Sprintf("test_key_%d", time.Now().Unix())
	value := "test_value"

	err = rdb.Set(ctx, key, value, time.Minute).Err()
	assert.NoError(t, err)

	result, err := rdb.Get(ctx, key).Result()
	assert.NoError(t, err)
	assert.Equal(t, value, result)

	// Cleanup
	rdb.Del(ctx, key)
}

func TestHealthEndpoint(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This test assumes the application is running
	// In CI, this would be started by docker-compose
	baseURL := getEnv("TEST_BASE_URL", "http://localhost:8080")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(baseURL + "/health")
	if err != nil {
		t.Skipf("Application not running, skipping health check test: %v", err)
		return
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
