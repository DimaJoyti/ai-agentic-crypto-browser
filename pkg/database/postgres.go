package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ai-agentic-browser/internal/config"
	"github.com/ai-agentic-browser/pkg/observability"
	_ "github.com/lib/pq"
)

// DB wraps sql.DB with additional functionality and performance optimizations
type DB struct {
	*sql.DB
	logger      *observability.Logger
	metrics     *DatabaseMetrics
	queryCache  *QueryCache
	connPool    *ConnectionPool
	readReplica *sql.DB
	mu          sync.RWMutex
}

// DatabaseMetrics tracks database performance metrics
type DatabaseMetrics struct {
	QueryCount      int64
	SlowQueryCount  int64
	ConnectionCount int64
	CacheHitCount   int64
	CacheMissCount  int64
	AvgQueryTime    time.Duration
	mu              sync.RWMutex
}

// QueryCache provides intelligent query result caching
type QueryCache struct {
	cache   map[string]*QueryCacheEntry
	maxSize int
	ttl     time.Duration
	mu      sync.RWMutex
}

// QueryCacheEntry represents a cached query result
type QueryCacheEntry struct {
	Data      interface{}
	ExpiresAt time.Time
	HitCount  int64
}

// ConnectionPool manages database connections with advanced pooling
type ConnectionPool struct {
	primary     *sql.DB
	readReplica *sql.DB
	config      *PoolConfig
	metrics     *PoolMetrics
	mu          sync.RWMutex
}

// PoolConfig contains connection pool configuration
type PoolConfig struct {
	MaxOpenConns        int
	MaxIdleConns        int
	ConnMaxLifetime     time.Duration
	ConnMaxIdleTime     time.Duration
	ReadReplicaEnabled  bool
	ReadWriteRatio      float64
	HealthCheckInterval time.Duration
}

// PoolMetrics tracks connection pool performance
type PoolMetrics struct {
	ActiveConnections int64
	IdleConnections   int64
	WaitCount         int64
	WaitDuration      time.Duration
	mu                sync.RWMutex
}

// NewPostgresDB creates a new PostgreSQL database connection with advanced optimizations
func NewPostgresDB(cfg config.DatabaseConfig) (*DB, error) {
	logger := &observability.Logger{}

	// Create primary database connection
	primary, err := sql.Open("postgres", cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to open primary database: %w", err)
	}

	// Configure advanced connection pool settings
	poolConfig := &PoolConfig{
		MaxOpenConns:        cfg.MaxOpenConns,
		MaxIdleConns:        cfg.MaxIdleConns,
		ConnMaxLifetime:     cfg.ConnMaxLifetime,
		ConnMaxIdleTime:     5 * time.Minute,
		ReadReplicaEnabled:  false, // Will be enabled based on environment
		ReadWriteRatio:      0.7,   // 70% reads, 30% writes
		HealthCheckInterval: 30 * time.Second,
	}

	// Apply optimized connection pool settings
	primary.SetMaxOpenConns(poolConfig.MaxOpenConns)
	primary.SetMaxIdleConns(poolConfig.MaxIdleConns)
	primary.SetConnMaxLifetime(poolConfig.ConnMaxLifetime)
	primary.SetConnMaxIdleTime(poolConfig.ConnMaxIdleTime)

	// Test the primary connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := primary.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping primary database: %w", err)
	}

	// Initialize components
	metrics := &DatabaseMetrics{}
	queryCache := NewQueryCache(1000, 5*time.Minute) // 1000 entries, 5min TTL
	connPool := &ConnectionPool{
		primary: primary,
		config:  poolConfig,
		metrics: &PoolMetrics{},
	}

	db := &DB{
		DB:         primary,
		logger:     logger,
		metrics:    metrics,
		queryCache: queryCache,
		connPool:   connPool,
	}

	// Start background health monitoring
	go db.startHealthMonitoring()

	logger.Info(context.Background(), "Database connection established with optimizations", map[string]interface{}{
		"max_open_conns":    poolConfig.MaxOpenConns,
		"max_idle_conns":    poolConfig.MaxIdleConns,
		"conn_max_lifetime": poolConfig.ConnMaxLifetime,
		"cache_enabled":     true,
		"health_monitoring": true,
	})

	return db, nil
}

// NewQueryCache creates a new query cache
func NewQueryCache(maxSize int, ttl time.Duration) *QueryCache {
	return &QueryCache{
		cache:   make(map[string]*QueryCacheEntry),
		maxSize: maxSize,
		ttl:     ttl,
	}
}

// Get retrieves a cached query result
func (qc *QueryCache) Get(key string) (interface{}, bool) {
	qc.mu.RLock()
	defer qc.mu.RUnlock()

	entry, exists := qc.cache[key]
	if !exists || time.Now().After(entry.ExpiresAt) {
		return nil, false
	}

	entry.HitCount++
	return entry.Data, true
}

// Set stores a query result in cache
func (qc *QueryCache) Set(key string, data interface{}) {
	qc.mu.Lock()
	defer qc.mu.Unlock()

	// Evict expired entries and enforce size limit
	if len(qc.cache) >= qc.maxSize {
		qc.evictOldest()
	}

	qc.cache[key] = &QueryCacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(qc.ttl),
		HitCount:  0,
	}
}

// evictOldest removes the oldest cache entry
func (qc *QueryCache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	for key, entry := range qc.cache {
		if oldestKey == "" || entry.ExpiresAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.ExpiresAt
		}
	}

	if oldestKey != "" {
		delete(qc.cache, oldestKey)
	}
}

// Clear removes all cached entries
func (qc *QueryCache) Clear() {
	qc.mu.Lock()
	defer qc.mu.Unlock()
	qc.cache = make(map[string]*QueryCacheEntry)
}

// Stats returns cache statistics
func (qc *QueryCache) Stats() map[string]interface{} {
	qc.mu.RLock()
	defer qc.mu.RUnlock()

	totalHits := int64(0)
	for _, entry := range qc.cache {
		totalHits += entry.HitCount
	}

	return map[string]interface{}{
		"size":       len(qc.cache),
		"max_size":   qc.maxSize,
		"total_hits": totalHits,
		"ttl":        qc.ttl,
	}
}

// QueryWithCache executes a query with intelligent caching
func (db *DB) QueryWithCache(ctx context.Context, cacheKey string, query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()

	// Try cache first for read queries
	if isReadQuery(query) {
		if cached, found := db.queryCache.Get(cacheKey); found {
			db.metrics.mu.Lock()
			db.metrics.CacheHitCount++
			db.metrics.mu.Unlock()

			db.logger.Debug(ctx, "Query cache hit", map[string]interface{}{
				"cache_key": cacheKey,
				"duration":  time.Since(start),
			})

			return cached.(*sql.Rows), nil
		}

		db.metrics.mu.Lock()
		db.metrics.CacheMissCount++
		db.metrics.mu.Unlock()
	}

	// Execute query
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	// Cache read query results
	if isReadQuery(query) && err == nil {
		db.queryCache.Set(cacheKey, rows)
	}

	// Update metrics
	duration := time.Since(start)
	db.updateMetrics(duration, query)

	return rows, nil
}

// ExecWithMetrics executes a query with performance tracking
func (db *DB) ExecWithMetrics(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()

	result, err := db.ExecContext(ctx, query, args...)

	duration := time.Since(start)
	db.updateMetrics(duration, query)

	if duration > 100*time.Millisecond {
		db.logger.Warn(ctx, "Slow query detected", map[string]interface{}{
			"query":    query,
			"duration": duration,
			"args":     args,
		})

		db.metrics.mu.Lock()
		db.metrics.SlowQueryCount++
		db.metrics.mu.Unlock()
	}

	return result, err
}

// GetReadConnection returns a connection optimized for read operations
func (db *DB) GetReadConnection() *sql.DB {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// Use read replica if available, otherwise use primary
	if db.readReplica != nil {
		return db.readReplica
	}
	return db.DB
}

// GetWriteConnection returns a connection for write operations
func (db *DB) GetWriteConnection() *sql.DB {
	return db.DB // Always use primary for writes
}

// updateMetrics updates database performance metrics
func (db *DB) updateMetrics(duration time.Duration, query string) {
	db.metrics.mu.Lock()
	defer db.metrics.mu.Unlock()

	db.metrics.QueryCount++

	// Update average query time using exponential moving average
	if db.metrics.AvgQueryTime == 0 {
		db.metrics.AvgQueryTime = duration
	} else {
		alpha := 0.1 // Smoothing factor
		db.metrics.AvgQueryTime = time.Duration(float64(db.metrics.AvgQueryTime)*(1-alpha) + float64(duration)*alpha)
	}
}

// isReadQuery determines if a query is a read operation
func isReadQuery(query string) bool {
	query = strings.ToUpper(strings.TrimSpace(query))
	return strings.HasPrefix(query, "SELECT") ||
		strings.HasPrefix(query, "WITH") ||
		strings.HasPrefix(query, "SHOW") ||
		strings.HasPrefix(query, "EXPLAIN")
}

// startHealthMonitoring starts background health monitoring
func (db *DB) startHealthMonitoring() {
	ticker := time.NewTicker(db.connPool.config.HealthCheckInterval)
	defer ticker.Stop()

	for range ticker.C {
		db.performHealthCheck()
	}
}

// performHealthCheck checks database health and updates metrics
func (db *DB) performHealthCheck() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check primary connection
	if err := db.DB.PingContext(ctx); err != nil {
		db.logger.Error(ctx, "Primary database health check failed", err)
		return
	}

	// Update connection pool metrics
	stats := db.DB.Stats()
	db.connPool.metrics.mu.Lock()
	db.connPool.metrics.ActiveConnections = int64(stats.OpenConnections)
	db.connPool.metrics.IdleConnections = int64(stats.Idle)
	db.connPool.metrics.WaitCount = stats.WaitCount
	db.connPool.metrics.WaitDuration = stats.WaitDuration
	db.connPool.metrics.mu.Unlock()

	// Log health status
	db.logger.Debug(ctx, "Database health check completed", map[string]interface{}{
		"open_connections": stats.OpenConnections,
		"idle_connections": stats.Idle,
		"wait_count":       stats.WaitCount,
		"wait_duration":    stats.WaitDuration,
	})
}

// GetMetrics returns current database metrics
func (db *DB) GetMetrics() map[string]interface{} {
	db.metrics.mu.RLock()
	defer db.metrics.mu.RUnlock()

	db.connPool.metrics.mu.RLock()
	defer db.connPool.metrics.mu.RUnlock()

	return map[string]interface{}{
		"query_count":        db.metrics.QueryCount,
		"slow_query_count":   db.metrics.SlowQueryCount,
		"cache_hit_count":    db.metrics.CacheHitCount,
		"cache_miss_count":   db.metrics.CacheMissCount,
		"avg_query_time":     db.metrics.AvgQueryTime,
		"active_connections": db.connPool.metrics.ActiveConnections,
		"idle_connections":   db.connPool.metrics.IdleConnections,
		"wait_count":         db.connPool.metrics.WaitCount,
		"wait_duration":      db.connPool.metrics.WaitDuration,
		"cache_stats":        db.queryCache.Stats(),
	}
}

// Close closes the database connection and cleanup resources
func (db *DB) Close() error {
	db.logger.Info(context.Background(), "Closing database connections")

	// Clear cache
	db.queryCache.Clear()

	// Close read replica if exists
	if db.readReplica != nil {
		if err := db.readReplica.Close(); err != nil {
			db.logger.Error(context.Background(), "Failed to close read replica", err)
		}
	}

	// Close primary connection
	return db.DB.Close()
}

// Health checks the database health
func (db *DB) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

// Transaction executes a function within a database transaction
func (db *DB) Transaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}
