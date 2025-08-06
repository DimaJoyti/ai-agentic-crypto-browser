package middleware

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/database"
	"github.com/ai-agentic-browser/pkg/observability"
)

// CacheMiddleware provides intelligent HTTP response caching
type CacheMiddleware struct {
	redis  *database.RedisClient
	logger *observability.Logger
	config *CacheConfig
	stats  *CacheStats
	mu     sync.RWMutex
}

// CacheConfig contains caching configuration
type CacheConfig struct {
	DefaultTTL       time.Duration
	MaxCacheSize     int64
	EnableGzip       bool
	CacheableStatus  []int
	CacheableMethods []string
	ExcludePaths     []string
	VaryHeaders      []string
}

// CacheStats tracks caching performance
type CacheStats struct {
	Hits      int64
	Misses    int64
	Sets      int64
	Errors    int64
	TotalSize int64
	mu        sync.RWMutex
}

// CachedResponse represents a cached HTTP response
type CachedResponse struct {
	StatusCode int                 `json:"status_code"`
	Headers    map[string][]string `json:"headers"`
	Body       []byte              `json:"body"`
	CreatedAt  time.Time           `json:"created_at"`
	TTL        time.Duration       `json:"ttl"`
	Size       int64               `json:"size"`
}

// cacheResponseWriter wraps http.ResponseWriter to capture response data
type cacheResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
	headers    http.Header
}

// NewCacheMiddleware creates a new cache middleware
func NewCacheMiddleware(redis *database.RedisClient, logger *observability.Logger) *CacheMiddleware {
	config := &CacheConfig{
		DefaultTTL:       5 * time.Minute,
		MaxCacheSize:     100 * 1024 * 1024, // 100MB
		EnableGzip:       true,
		CacheableStatus:  []int{200, 201, 202, 203, 204, 300, 301, 302, 304, 404, 410},
		CacheableMethods: []string{"GET", "HEAD"},
		ExcludePaths:     []string{"/health", "/metrics", "/auth/"},
		VaryHeaders:      []string{"Accept", "Accept-Encoding", "Authorization"},
	}

	return &CacheMiddleware{
		redis:  redis,
		logger: logger,
		config: config,
		stats:  &CacheStats{},
	}
}

// Middleware returns the caching middleware function
func (cm *CacheMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip caching for non-cacheable methods
			if !cm.isCacheableMethod(r.Method) {
				next.ServeHTTP(w, r)
				return
			}

			// Skip caching for excluded paths
			if cm.isExcludedPath(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			// Generate cache key
			cacheKey := cm.generateCacheKey(r)

			// Try to serve from cache
			if cached, found := cm.getFromCache(r.Context(), cacheKey); found {
				cm.serveCachedResponse(w, cached)
				cm.updateStats("hit")
				return
			}

			// Wrap response writer to capture response
			rw := &cacheResponseWriter{
				ResponseWriter: w,
				body:           &bytes.Buffer{},
				headers:        make(http.Header),
			}

			// Call next handler
			next.ServeHTTP(rw, r)

			// Cache the response if cacheable
			if cm.isCacheableResponse(rw.statusCode) {
				cached := &CachedResponse{
					StatusCode: rw.statusCode,
					Headers:    rw.headers,
					Body:       rw.body.Bytes(),
					CreatedAt:  time.Now(),
					TTL:        cm.config.DefaultTTL,
					Size:       int64(len(rw.body.Bytes())),
				}

				if err := cm.setCache(r.Context(), cacheKey, cached); err != nil {
					cm.logger.Error(r.Context(), "Failed to cache response", err, map[string]interface{}{
						"cache_key": cacheKey,
						"path":      r.URL.Path,
					})
					cm.updateStats("error")
				} else {
					cm.updateStats("set")
				}
			}

			cm.updateStats("miss")
		})
	}
}

// generateCacheKey creates a unique cache key for the request
func (cm *CacheMiddleware) generateCacheKey(r *http.Request) string {
	h := md5.New()

	// Include method, path, and query parameters
	h.Write([]byte(r.Method))
	h.Write([]byte(r.URL.Path))
	h.Write([]byte(r.URL.RawQuery))

	// Include vary headers
	for _, header := range cm.config.VaryHeaders {
		if value := r.Header.Get(header); value != "" {
			h.Write([]byte(header + ":" + value))
		}
	}

	return "cache:" + hex.EncodeToString(h.Sum(nil))
}

// getFromCache retrieves a cached response
func (cm *CacheMiddleware) getFromCache(ctx context.Context, key string) (*CachedResponse, bool) {
	data, found, err := cm.redis.GetLayered(ctx, key)
	if err != nil || !found {
		return nil, false
	}

	var cached CachedResponse
	if jsonData, ok := data.(string); ok {
		if err := json.Unmarshal([]byte(jsonData), &cached); err != nil {
			cm.logger.Error(ctx, "Failed to unmarshal cached response", err)
			return nil, false
		}
	} else if err := json.Unmarshal(data.([]byte), &cached); err != nil {
		cm.logger.Error(ctx, "Failed to unmarshal cached response", err)
		return nil, false
	}

	// Check if cache entry is still valid
	if time.Since(cached.CreatedAt) > cached.TTL {
		cm.redis.DeleteKeys(ctx, key)
		return nil, false
	}

	return &cached, true
}

// setCache stores a response in cache
func (cm *CacheMiddleware) setCache(ctx context.Context, key string, cached *CachedResponse) error {
	// Check cache size limits
	if cached.Size > cm.config.MaxCacheSize {
		return fmt.Errorf("response too large to cache: %d bytes", cached.Size)
	}

	data, err := json.Marshal(cached)
	if err != nil {
		return fmt.Errorf("failed to marshal cached response: %w", err)
	}

	// Use L2 cache for most responses
	return cm.redis.SetLayered(ctx, key, string(data), database.L2Cache)
}

// serveCachedResponse serves a cached response
func (cm *CacheMiddleware) serveCachedResponse(w http.ResponseWriter, cached *CachedResponse) {
	// Set headers
	for key, values := range cached.Headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Add cache headers
	w.Header().Set("X-Cache", "HIT")
	w.Header().Set("X-Cache-Date", cached.CreatedAt.Format(time.RFC3339))
	w.Header().Set("Age", strconv.Itoa(int(time.Since(cached.CreatedAt).Seconds())))

	// Set status code and write body
	w.WriteHeader(cached.StatusCode)
	w.Write(cached.Body)
}

// isCacheableMethod checks if the HTTP method is cacheable
func (cm *CacheMiddleware) isCacheableMethod(method string) bool {
	for _, m := range cm.config.CacheableMethods {
		if m == method {
			return true
		}
	}
	return false
}

// isCacheableResponse checks if the response status is cacheable
func (cm *CacheMiddleware) isCacheableResponse(statusCode int) bool {
	for _, code := range cm.config.CacheableStatus {
		if code == statusCode {
			return true
		}
	}
	return false
}

// isExcludedPath checks if the path should be excluded from caching
func (cm *CacheMiddleware) isExcludedPath(path string) bool {
	for _, excluded := range cm.config.ExcludePaths {
		if strings.HasPrefix(path, excluded) {
			return true
		}
	}
	return false
}

// updateStats updates cache statistics
func (cm *CacheMiddleware) updateStats(operation string) {
	cm.stats.mu.Lock()
	defer cm.stats.mu.Unlock()

	switch operation {
	case "hit":
		cm.stats.Hits++
	case "miss":
		cm.stats.Misses++
	case "set":
		cm.stats.Sets++
	case "error":
		cm.stats.Errors++
	}
}

// GetStats returns current cache statistics
func (cm *CacheMiddleware) GetStats() map[string]interface{} {
	cm.stats.mu.RLock()
	defer cm.stats.mu.RUnlock()

	total := cm.stats.Hits + cm.stats.Misses
	hitRate := float64(0)
	if total > 0 {
		hitRate = float64(cm.stats.Hits) / float64(total) * 100
	}

	return map[string]interface{}{
		"hits":       cm.stats.Hits,
		"misses":     cm.stats.Misses,
		"sets":       cm.stats.Sets,
		"errors":     cm.stats.Errors,
		"hit_rate":   hitRate,
		"total_size": cm.stats.TotalSize,
	}
}

// cacheResponseWriter implementation
func (rw *cacheResponseWriter) Write(data []byte) (int, error) {
	rw.body.Write(data)
	return rw.ResponseWriter.Write(data)
}

func (rw *cacheResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode

	// Copy headers
	for key, values := range rw.ResponseWriter.Header() {
		rw.headers[key] = values
	}

	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Header() http.Header {
	return rw.ResponseWriter.Header()
}
