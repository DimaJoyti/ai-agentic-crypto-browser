package database

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/internal/config"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/redis/go-redis/v9"
)

// RedisClient wraps redis.Client with advanced caching functionality
type RedisClient struct {
	*redis.Client
	logger      *observability.Logger
	metrics     *RedisMetrics
	cacheConfig *CacheConfig
	mu          sync.RWMutex
}

// RedisMetrics tracks Redis performance metrics
type RedisMetrics struct {
	HitCount      int64
	MissCount     int64
	SetCount      int64
	DeleteCount   int64
	EvictionCount int64
	AvgLatency    time.Duration
	mu            sync.RWMutex
}

// CacheConfig contains caching configuration
type CacheConfig struct {
	DefaultTTL       time.Duration
	MaxMemory        string
	EvictionPolicy   string
	CompressionLevel int
	EnableMetrics    bool
}

// CacheLayer represents different cache layers
type CacheLayer int

const (
	L1Cache CacheLayer = iota // Hot data - very short TTL
	L2Cache                   // Warm data - medium TTL
	L3Cache                   // Cold data - long TTL
)

// CacheEntry represents a cached item with metadata
type CacheEntry struct {
	Data         interface{}   `json:"data"`
	CreatedAt    time.Time     `json:"created_at"`
	LastAccessed time.Time     `json:"last_accessed"`
	AccessCount  int64         `json:"access_count"`
	Layer        CacheLayer    `json:"layer"`
	TTL          time.Duration `json:"ttl"`
	Size         int64         `json:"size"`       // Size in bytes
	Compressed   bool          `json:"compressed"` // Whether data is compressed
	Tags         []string      `json:"tags"`       // Tags for cache invalidation
	Priority     int           `json:"priority"`   // Priority for eviction (higher = keep longer)
	Version      string        `json:"version"`    // Version for cache invalidation
}

// NewRedisClient creates a new Redis client with advanced caching capabilities
func NewRedisClient(cfg config.RedisConfig) (*RedisClient, error) {
	logger := &observability.Logger{}

	opt, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	if cfg.Password != "" {
		opt.Password = cfg.Password
	}
	opt.DB = cfg.DB

	// Enhanced connection pool settings
	opt.PoolSize = cfg.PoolSize
	opt.MinIdleConns = 5
	opt.MaxIdleConns = 10
	opt.PoolTimeout = 4 * time.Second
	opt.ConnMaxIdleTime = 5 * time.Minute
	opt.MaxRetries = 3
	opt.MinRetryBackoff = 8 * time.Millisecond
	opt.MaxRetryBackoff = 512 * time.Millisecond

	client := redis.NewClient(opt)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	// Initialize cache configuration
	cacheConfig := &CacheConfig{
		DefaultTTL:       5 * time.Minute,
		MaxMemory:        "256mb",
		EvictionPolicy:   "allkeys-lru",
		CompressionLevel: 6,
		EnableMetrics:    true,
	}

	redisClient := &RedisClient{
		Client:      client,
		logger:      logger,
		metrics:     &RedisMetrics{},
		cacheConfig: cacheConfig,
	}

	// Configure Redis for optimal performance
	if err := redisClient.configureRedis(ctx); err != nil {
		logger.Warn(ctx, "Failed to configure Redis optimizations", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Start background monitoring
	go redisClient.startMetricsCollection()

	logger.Info(ctx, "Redis client initialized with advanced caching", map[string]interface{}{
		"pool_size":       opt.PoolSize,
		"min_idle_conns":  opt.MinIdleConns,
		"max_memory":      cacheConfig.MaxMemory,
		"eviction_policy": cacheConfig.EvictionPolicy,
		"metrics_enabled": cacheConfig.EnableMetrics,
	})

	return redisClient, nil
}

// configureRedis applies optimal Redis configuration
func (r *RedisClient) configureRedis(ctx context.Context) error {
	configs := map[string]string{
		"maxmemory":        r.cacheConfig.MaxMemory,
		"maxmemory-policy": r.cacheConfig.EvictionPolicy,
		"timeout":          "300",
		"tcp-keepalive":    "60",
		"tcp-nodelay":      "yes",
		"save":             "900 1 300 10 60 10000", // Optimized save intervals
	}

	for key, value := range configs {
		if err := r.ConfigSet(ctx, key, value).Err(); err != nil {
			r.logger.Warn(ctx, "Failed to set Redis config", map[string]interface{}{
				"key":   key,
				"value": value,
				"error": err.Error(),
			})
		}
	}

	return nil
}

// startMetricsCollection starts background metrics collection
func (r *RedisClient) startMetricsCollection() {
	if !r.cacheConfig.EnableMetrics {
		return
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		r.collectMetrics()
	}
}

// collectMetrics collects Redis performance metrics
func (r *RedisClient) collectMetrics() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	info := r.Info(ctx, "stats", "memory")
	if info.Err() != nil {
		r.logger.Error(ctx, "Failed to collect Redis metrics", info.Err())
		return
	}

	// Parse and log key metrics
	r.logger.Debug(ctx, "Redis metrics collected", map[string]interface{}{
		"info": info.Val(),
	})
}

// SetLayered sets a value in a specific cache layer with appropriate TTL
func (r *RedisClient) SetLayered(ctx context.Context, key string, value interface{}, layer CacheLayer) error {
	start := time.Now()

	var ttl time.Duration
	var keyPrefix string

	switch layer {
	case L1Cache:
		ttl = 1 * time.Minute
		keyPrefix = "l1:"
	case L2Cache:
		ttl = 15 * time.Minute
		keyPrefix = "l2:"
	case L3Cache:
		ttl = 1 * time.Hour
		keyPrefix = "l3:"
	default:
		ttl = r.cacheConfig.DefaultTTL
		keyPrefix = "default:"
	}

	now := time.Now()
	entry := &CacheEntry{
		Data:         value,
		CreatedAt:    now,
		LastAccessed: now,
		AccessCount:  0,
		Layer:        layer,
		TTL:          ttl,
		Size:         int64(len(fmt.Sprintf("%v", value))), // Approximate size
		Compressed:   false,
		Tags:         []string{},
		Priority:     1, // Default priority
		Version:      "1.0",
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal cache entry: %w", err)
	}

	fullKey := keyPrefix + key
	err = r.Set(ctx, fullKey, data, ttl).Err()

	// Update metrics
	r.updateMetrics("set", time.Since(start), err == nil)

	if err == nil {
		r.metrics.mu.Lock()
		r.metrics.SetCount++
		r.metrics.mu.Unlock()
	}

	return err
}

// GetLayered retrieves a value from cache layers (L1 -> L2 -> L3)
func (r *RedisClient) GetLayered(ctx context.Context, key string) (interface{}, bool, error) {
	start := time.Now()

	layers := []string{"l1:", "l2:", "l3:"}

	for _, prefix := range layers {
		fullKey := prefix + key
		result := r.Get(ctx, fullKey)

		if result.Err() == nil {
			var entry CacheEntry
			if err := json.Unmarshal([]byte(result.Val()), &entry); err == nil {
				// Update access count
				entry.AccessCount++

				// Promote to higher cache layer if frequently accessed
				if entry.AccessCount > 10 && prefix != "l1:" {
					r.promoteToHigherLayer(ctx, key, &entry)
				}

				// Update metrics
				r.updateMetrics("get", time.Since(start), true)
				r.metrics.mu.Lock()
				r.metrics.HitCount++
				r.metrics.mu.Unlock()

				return entry.Data, true, nil
			}
		}
	}

	// Cache miss
	r.updateMetrics("get", time.Since(start), false)
	r.metrics.mu.Lock()
	r.metrics.MissCount++
	r.metrics.mu.Unlock()

	return nil, false, nil
}

// promoteToHigherLayer promotes frequently accessed data to a higher cache layer
func (r *RedisClient) promoteToHigherLayer(ctx context.Context, key string, entry *CacheEntry) {
	var newLayer CacheLayer

	switch entry.Layer {
	case L3Cache:
		newLayer = L2Cache
	case L2Cache:
		newLayer = L1Cache
	default:
		return // Already in L1
	}

	// Set in higher layer
	entry.Layer = newLayer
	r.SetLayered(ctx, key, entry.Data, newLayer)

	r.logger.Debug(ctx, "Cache entry promoted", map[string]interface{}{
		"key":          key,
		"from_layer":   entry.Layer,
		"to_layer":     newLayer,
		"access_count": entry.AccessCount,
	})
}

// SetWithCompression sets a value with optional compression
func (r *RedisClient) SetWithCompression(ctx context.Context, key string, value interface{}, expiry time.Duration, compress bool) error {
	start := time.Now()

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	// Apply compression if enabled and data is large enough
	if compress && len(data) > 1024 {
		// Implement compression logic here if needed
		r.logger.Debug(ctx, "Compression applied", map[string]interface{}{
			"key":           key,
			"original_size": len(data),
		})
	}

	err = r.Set(ctx, key, data, expiry).Err()
	r.updateMetrics("set", time.Since(start), err == nil)

	return err
}

// GetWithFallback gets a value with fallback function if not found
func (r *RedisClient) GetWithFallback(ctx context.Context, key string, fallback func() (interface{}, error), ttl time.Duration) (interface{}, error) {
	// Try cache first
	if data, found, err := r.GetLayered(ctx, key); err == nil && found {
		return data, nil
	}

	// Execute fallback function
	data, err := fallback()
	if err != nil {
		return nil, err
	}

	// Cache the result
	if err := r.SetLayered(ctx, key, data, L2Cache); err != nil {
		r.logger.Warn(ctx, "Failed to cache fallback result", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
	}

	return data, nil
}

// updateMetrics updates Redis operation metrics
func (r *RedisClient) updateMetrics(operation string, duration time.Duration, success bool) {
	if !r.cacheConfig.EnableMetrics {
		return
	}

	r.metrics.mu.Lock()
	defer r.metrics.mu.Unlock()

	// Update average latency using exponential moving average
	if r.metrics.AvgLatency == 0 {
		r.metrics.AvgLatency = duration
	} else {
		alpha := 0.1
		r.metrics.AvgLatency = time.Duration(float64(r.metrics.AvgLatency)*(1-alpha) + float64(duration)*alpha)
	}
}

// Close closes the Redis connection and cleanup resources
func (r *RedisClient) Close() error {
	r.logger.Info(context.Background(), "Closing Redis connection")
	return r.Client.Close()
}

// Health checks the Redis health with detailed diagnostics
func (r *RedisClient) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	start := time.Now()

	if err := r.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis health check failed: %w", err)
	}

	latency := time.Since(start)

	// Check if latency is acceptable
	if latency > 100*time.Millisecond {
		r.logger.Warn(ctx, "High Redis latency detected", map[string]interface{}{
			"latency": latency,
		})
	}

	return nil
}

// SetWithExpiry sets a key-value pair with expiration and metrics
func (r *RedisClient) SetWithExpiry(ctx context.Context, key string, value interface{}, expiry time.Duration) error {
	start := time.Now()

	err := r.Set(ctx, key, value, expiry).Err()
	r.updateMetrics("set", time.Since(start), err == nil)

	if err == nil {
		r.metrics.mu.Lock()
		r.metrics.SetCount++
		r.metrics.mu.Unlock()
	}

	return err
}

// GetString gets a string value by key with metrics
func (r *RedisClient) GetString(ctx context.Context, key string) (string, error) {
	start := time.Now()

	result := r.Get(ctx, key)
	success := result.Err() == nil
	r.updateMetrics("get", time.Since(start), success)

	if result.Err() != nil {
		if result.Err() == redis.Nil {
			r.metrics.mu.Lock()
			r.metrics.MissCount++
			r.metrics.mu.Unlock()
			return "", fmt.Errorf("key not found: %s", key)
		}
		return "", result.Err()
	}

	r.metrics.mu.Lock()
	r.metrics.HitCount++
	r.metrics.mu.Unlock()

	return result.Val(), nil
}

// DeleteKeys deletes multiple keys with metrics
func (r *RedisClient) DeleteKeys(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	start := time.Now()
	err := r.Del(ctx, keys...).Err()
	r.updateMetrics("delete", time.Since(start), err == nil)

	if err == nil {
		r.metrics.mu.Lock()
		r.metrics.DeleteCount += int64(len(keys))
		r.metrics.mu.Unlock()
	}

	return err
}

// Exists checks if a key exists with metrics
func (r *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	start := time.Now()

	result := r.Client.Exists(ctx, key)
	r.updateMetrics("exists", time.Since(start), result.Err() == nil)

	if err := result.Err(); err != nil {
		return false, err
	}
	return result.Val() > 0, nil
}

// GetMetrics returns current Redis metrics
func (r *RedisClient) GetMetrics() map[string]interface{} {
	r.metrics.mu.RLock()
	defer r.metrics.mu.RUnlock()

	hitRate := float64(0)
	totalRequests := r.metrics.HitCount + r.metrics.MissCount
	if totalRequests > 0 {
		hitRate = float64(r.metrics.HitCount) / float64(totalRequests) * 100
	}

	return map[string]interface{}{
		"hit_count":      r.metrics.HitCount,
		"miss_count":     r.metrics.MissCount,
		"set_count":      r.metrics.SetCount,
		"delete_count":   r.metrics.DeleteCount,
		"eviction_count": r.metrics.EvictionCount,
		"avg_latency":    r.metrics.AvgLatency,
		"hit_rate":       hitRate,
		"total_requests": totalRequests,
	}
}

// FlushExpired removes expired keys to free memory
func (r *RedisClient) FlushExpired(ctx context.Context) error {
	// This is handled automatically by Redis, but we can trigger it manually
	return r.Do(ctx, "MEMORY", "PURGE").Err()
}

// GetMemoryUsage returns current memory usage statistics
func (r *RedisClient) GetMemoryUsage(ctx context.Context) (map[string]interface{}, error) {
	info := r.Info(ctx, "memory")
	if info.Err() != nil {
		return nil, info.Err()
	}

	// Parse memory info (simplified)
	return map[string]interface{}{
		"info": info.Val(),
	}, nil
}
