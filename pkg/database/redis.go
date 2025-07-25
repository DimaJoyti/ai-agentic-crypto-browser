package database

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/ai-agentic-browser/internal/config"
)

// RedisClient wraps redis.Client with additional functionality
type RedisClient struct {
	*redis.Client
}

// NewRedisClient creates a new Redis client
func NewRedisClient(cfg config.RedisConfig) (*RedisClient, error) {
	opt, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	if cfg.Password != "" {
		opt.Password = cfg.Password
	}
	opt.DB = cfg.DB

	client := redis.NewClient(opt)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	return &RedisClient{Client: client}, nil
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
	return r.Client.Close()
}

// Health checks the Redis health
func (r *RedisClient) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := r.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis health check failed: %w", err)
	}

	return nil
}

// SetWithExpiry sets a key-value pair with expiration
func (r *RedisClient) SetWithExpiry(ctx context.Context, key string, value interface{}, expiry time.Duration) error {
	return r.Set(ctx, key, value, expiry).Err()
}

// GetString gets a string value by key
func (r *RedisClient) GetString(ctx context.Context, key string) (string, error) {
	result := r.Get(ctx, key)
	if err := result.Err(); err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("key not found: %s", key)
		}
		return "", err
	}
	return result.Val(), nil
}

// DeleteKeys deletes multiple keys
func (r *RedisClient) DeleteKeys(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return r.Del(ctx, keys...).Err()
}

// Exists checks if a key exists
func (r *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	result := r.Client.Exists(ctx, key)
	if err := result.Err(); err != nil {
		return false, err
	}
	return result.Val() > 0, nil
}
