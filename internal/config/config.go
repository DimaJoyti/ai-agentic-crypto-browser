package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server        ServerConfig
	Database      DatabaseConfig
	Redis         RedisConfig
	JWT           JWTConfig
	AI            AIConfig
	Web3          Web3Config
	Browser       BrowserConfig
	Observability ObservabilityConfig
	RateLimit     RateLimitConfig
	Security      SecurityConfig
}

type ServerConfig struct {
	Port         string
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	URL      string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret             string
	Expiry             time.Duration
	RefreshTokenExpiry time.Duration
}

type AIConfig struct {
	Provider       string
	OpenAIKey      string
	AnthropicKey   string
	ModelName      string
	OllamaConfig   OllamaConfig
	LMStudioConfig LMStudioConfig
}

type OllamaConfig struct {
	BaseURL             string
	Model               string
	Temperature         float64
	TopP                float64
	TopK                int
	NumCtx              int
	Timeout             time.Duration
	MaxRetries          int
	RetryDelay          time.Duration
	HealthCheckInterval time.Duration
}

type LMStudioConfig struct {
	BaseURL             string
	Model               string
	Temperature         float64
	MaxTokens           int
	TopP                float64
	Timeout             time.Duration
	MaxRetries          int
	RetryDelay          time.Duration
	HealthCheckInterval time.Duration
}

type Web3Config struct {
	EthereumRPC string
	PolygonRPC  string
	ArbitrumRPC string
	OptimismRPC string
}

type BrowserConfig struct {
	Headless   bool
	DisableGPU bool
	NoSandbox  bool
	Timeout    time.Duration
}

type ObservabilityConfig struct {
	JaegerEndpoint string
	ServiceName    string
	LogLevel       string
	LogFormat      string
}

type RateLimitConfig struct {
	RequestsPerMinute int
	Burst             int
}

type SecurityConfig struct {
	CORSAllowedOrigins []string
	BCryptCost         int
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			Host:         getEnv("HOST", "0.0.0.0"),
			ReadTimeout:  getDurationEnv("READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getDurationEnv("WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getDurationEnv("IDLE_TIMEOUT", 60*time.Second),
		},
		Database: DatabaseConfig{
			URL:             getEnv("DATABASE_URL", ""),
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		Redis: RedisConfig{
			URL:      getEnv("REDIS_URL", "redis://localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getIntEnv("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:             getEnv("JWT_SECRET", ""),
			Expiry:             getDurationEnv("JWT_EXPIRY", 24*time.Hour),
			RefreshTokenExpiry: getDurationEnv("REFRESH_TOKEN_EXPIRY", 168*time.Hour),
		},
		AI: AIConfig{
			Provider:     getEnv("AI_MODEL_PROVIDER", "openai"),
			OpenAIKey:    getEnv("OPENAI_API_KEY", ""),
			AnthropicKey: getEnv("ANTHROPIC_API_KEY", ""),
			ModelName:    getEnv("AI_MODEL_NAME", "gpt-4-turbo-preview"),
			OllamaConfig: OllamaConfig{
				BaseURL:             getEnv("OLLAMA_BASE_URL", "http://localhost:11434"),
				Model:               getEnv("OLLAMA_MODEL", "qwen3"),
				Temperature:         getFloatEnv("OLLAMA_TEMPERATURE", 0.7),
				TopP:                getFloatEnv("OLLAMA_TOP_P", 1.0),
				TopK:                getIntEnv("OLLAMA_TOP_K", 40),
				NumCtx:              getIntEnv("OLLAMA_NUM_CTX", 2048),
				Timeout:             getDurationEnv("OLLAMA_TIMEOUT", 300*time.Second),
				MaxRetries:          getIntEnv("OLLAMA_MAX_RETRIES", 3),
				RetryDelay:          getDurationEnv("OLLAMA_RETRY_DELAY", 2*time.Second),
				HealthCheckInterval: getDurationEnv("OLLAMA_HEALTH_CHECK_INTERVAL", 30*time.Second),
			},
			LMStudioConfig: LMStudioConfig{
				BaseURL:             getEnv("LMSTUDIO_BASE_URL", "http://localhost:1234/v1"),
				Model:               getEnv("LMSTUDIO_MODEL", "local-model"),
				Temperature:         getFloatEnv("LMSTUDIO_TEMPERATURE", 0.7),
				MaxTokens:           getIntEnv("LMSTUDIO_MAX_TOKENS", 4000),
				TopP:                getFloatEnv("LMSTUDIO_TOP_P", 1.0),
				Timeout:             getDurationEnv("LMSTUDIO_TIMEOUT", 300*time.Second),
				MaxRetries:          getIntEnv("LMSTUDIO_MAX_RETRIES", 3),
				RetryDelay:          getDurationEnv("LMSTUDIO_RETRY_DELAY", 2*time.Second),
				HealthCheckInterval: getDurationEnv("LMSTUDIO_HEALTH_CHECK_INTERVAL", 30*time.Second),
			},
		},
		Web3: Web3Config{
			EthereumRPC: getEnv("ETHEREUM_RPC_URL", ""),
			PolygonRPC:  getEnv("POLYGON_RPC_URL", ""),
			ArbitrumRPC: getEnv("ARBITRUM_RPC_URL", ""),
			OptimismRPC: getEnv("OPTIMISM_RPC_URL", ""),
		},
		Browser: BrowserConfig{
			Headless:   getBoolEnv("CHROME_HEADLESS", true),
			DisableGPU: getBoolEnv("CHROME_DISABLE_GPU", true),
			NoSandbox:  getBoolEnv("CHROME_NO_SANDBOX", true),
			Timeout:    getDurationEnv("BROWSER_TIMEOUT", 30*time.Second),
		},
		Observability: ObservabilityConfig{
			JaegerEndpoint: getEnv("JAEGER_ENDPOINT", "http://localhost:14268/api/traces"),
			ServiceName:    getEnv("OTEL_SERVICE_NAME", "agentic-browser"),
			LogLevel:       getEnv("LOG_LEVEL", "info"),
			LogFormat:      getEnv("LOG_FORMAT", "json"),
		},
		RateLimit: RateLimitConfig{
			RequestsPerMinute: getIntEnv("RATE_LIMIT_REQUESTS_PER_MINUTE", 100),
			Burst:             getIntEnv("RATE_LIMIT_BURST", 20),
		},
		Security: SecurityConfig{
			CORSAllowedOrigins: getSliceEnv("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
			BCryptCost:         getIntEnv("BCRYPT_COST", 12),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.Database.URL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	return nil
}

// Helper functions for environment variable parsing
func getEnv(key, defaultValue string) string {
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

func getSliceEnv(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Simple comma-separated parsing
		result := []string{}
		for _, item := range []string{value} {
			if item != "" {
				result = append(result, item)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultValue
}
