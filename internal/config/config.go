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
	Terminal      TerminalConfig
	Observability ObservabilityConfig
	RateLimit     RateLimitConfig
	Security      SecurityConfig
	Logger        LoggerConfig
	Firebase      FirebaseConfig
	MCP           MCPConfig
}

type ServerConfig struct {
	Port         string
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	URL                 string
	MaxOpenConns        int
	MaxIdleConns        int
	ConnMaxLifetime     time.Duration
	ConnMaxIdleTime     time.Duration
	QueryTimeout        time.Duration
	EnableQueryCache    bool
	CacheSize           int
	CacheTTL            time.Duration
	ReadReplicaURL      string
	EnableReadReplica   bool
	HealthCheckInterval time.Duration
}

type RedisConfig struct {
	URL                string
	Password           string
	DB                 int
	PoolSize           int
	MinIdleConns       int
	MaxIdleConns       int
	PoolTimeout        time.Duration
	IdleTimeout        time.Duration
	IdleCheckFrequency time.Duration
	MaxRetries         int
	MinRetryBackoff    time.Duration
	MaxRetryBackoff    time.Duration
	EnableMetrics      bool
	MaxMemory          string
	EvictionPolicy     string
	CompressionLevel   int
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
	EthereumRPC        string
	PolygonRPC         string
	ArbitrumRPC        string
	OptimismRPC        string
	BSCMainnetRPC      string
	BSCTestnetRPC      string
	SepoliaRPC         string
	IPFSNodeURL        string
	IPFSGateway        string
	IPFSMaxFileSize    int64
	GasOptimization    bool
	HardwareWallets    bool
	ENSResolution      bool
	TransactionTimeout time.Duration
	MaxRetries         int
	RetryDelay         time.Duration
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
			URL:                 getEnv("DATABASE_URL", ""),
			MaxOpenConns:        getIntEnv("DB_MAX_OPEN_CONNS", 50),
			MaxIdleConns:        getIntEnv("DB_MAX_IDLE_CONNS", 25),
			ConnMaxLifetime:     getDurationEnv("DB_CONN_MAX_LIFETIME", 5*time.Minute),
			ConnMaxIdleTime:     getDurationEnv("DB_CONN_MAX_IDLE_TIME", 5*time.Minute),
			QueryTimeout:        getDurationEnv("DB_QUERY_TIMEOUT", 30*time.Second),
			EnableQueryCache:    getBoolEnv("DB_ENABLE_QUERY_CACHE", true),
			CacheSize:           getIntEnv("DB_CACHE_SIZE", 1000),
			CacheTTL:            getDurationEnv("DB_CACHE_TTL", 5*time.Minute),
			ReadReplicaURL:      getEnv("DATABASE_READ_REPLICA_URL", ""),
			EnableReadReplica:   getBoolEnv("DB_ENABLE_READ_REPLICA", false),
			HealthCheckInterval: getDurationEnv("DB_HEALTH_CHECK_INTERVAL", 30*time.Second),
		},
		Redis: RedisConfig{
			URL:                getEnv("REDIS_URL", "redis://localhost:6379"),
			Password:           getEnv("REDIS_PASSWORD", ""),
			DB:                 getIntEnv("REDIS_DB", 0),
			PoolSize:           getIntEnv("REDIS_POOL_SIZE", 20),
			MinIdleConns:       getIntEnv("REDIS_MIN_IDLE_CONNS", 5),
			MaxIdleConns:       getIntEnv("REDIS_MAX_IDLE_CONNS", 10),
			PoolTimeout:        getDurationEnv("REDIS_POOL_TIMEOUT", 4*time.Second),
			IdleTimeout:        getDurationEnv("REDIS_IDLE_TIMEOUT", 5*time.Minute),
			IdleCheckFrequency: getDurationEnv("REDIS_IDLE_CHECK_FREQUENCY", 1*time.Minute),
			MaxRetries:         getIntEnv("REDIS_MAX_RETRIES", 3),
			MinRetryBackoff:    getDurationEnv("REDIS_MIN_RETRY_BACKOFF", 8*time.Millisecond),
			MaxRetryBackoff:    getDurationEnv("REDIS_MAX_RETRY_BACKOFF", 512*time.Millisecond),
			EnableMetrics:      getBoolEnv("REDIS_ENABLE_METRICS", true),
			MaxMemory:          getEnv("REDIS_MAX_MEMORY", "256mb"),
			EvictionPolicy:     getEnv("REDIS_EVICTION_POLICY", "allkeys-lru"),
			CompressionLevel:   getIntEnv("REDIS_COMPRESSION_LEVEL", 6),
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
			EthereumRPC:        getEnv("ETHEREUM_RPC_URL", ""),
			PolygonRPC:         getEnv("POLYGON_RPC_URL", ""),
			ArbitrumRPC:        getEnv("ARBITRUM_RPC_URL", ""),
			OptimismRPC:        getEnv("OPTIMISM_RPC_URL", ""),
			BSCMainnetRPC:      getEnv("BSC_MAINNET_RPC_URL", ""),
			BSCTestnetRPC:      getEnv("BSC_TESTNET_RPC_URL", ""),
			SepoliaRPC:         getEnv("SEPOLIA_RPC_URL", ""),
			IPFSNodeURL:        getEnv("IPFS_NODE_URL", "http://localhost:5001"),
			IPFSGateway:        getEnv("IPFS_GATEWAY", "https://ipfs.io"),
			IPFSMaxFileSize:    int64(getIntEnv("IPFS_MAX_FILE_SIZE", 10*1024*1024)), // 10MB default
			GasOptimization:    getBoolEnv("WEB3_GAS_OPTIMIZATION", true),
			HardwareWallets:    getBoolEnv("WEB3_HARDWARE_WALLETS", true),
			ENSResolution:      getBoolEnv("WEB3_ENS_RESOLUTION", true),
			TransactionTimeout: getDurationEnv("WEB3_TRANSACTION_TIMEOUT", 5*time.Minute),
			MaxRetries:         getIntEnv("WEB3_MAX_RETRIES", 3),
			RetryDelay:         getDurationEnv("WEB3_RETRY_DELAY", 2*time.Second),
		},
		Browser: BrowserConfig{
			Headless:   getBoolEnv("CHROME_HEADLESS", true),
			DisableGPU: getBoolEnv("CHROME_DISABLE_GPU", true),
			NoSandbox:  getBoolEnv("CHROME_NO_SANDBOX", true),
			Timeout:    getDurationEnv("BROWSER_TIMEOUT", 30*time.Second),
		},
		Terminal: TerminalConfig{
			Host:         getEnv("TERMINAL_HOST", "0.0.0.0"),
			Port:         getIntEnv("TERMINAL_PORT", 8085),
			ReadTimeout:  getDurationEnv("TERMINAL_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getDurationEnv("TERMINAL_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getDurationEnv("TERMINAL_IDLE_TIMEOUT", 60*time.Second),
			MaxSessions:  getIntEnv("TERMINAL_MAX_SESSIONS", 10),
			SessionTTL:   getDurationEnv("TERMINAL_SESSION_TTL", 24*time.Hour),
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
		Logger: LoggerConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		Firebase: FirebaseConfig{
			ProjectID:          getEnv("FIREBASE_PROJECT_ID", "ai-agentic-crypto-browser"),
			CredentialsPath:    getEnv("FIREBASE_CREDENTIALS_PATH", "./configs/firebase-service-account.json"),
			DatabaseURL:        getEnv("FIREBASE_DATABASE_URL", ""),
			StorageBucket:      getEnv("FIREBASE_STORAGE_BUCKET", ""),
			EnableAuth:         getBoolEnv("FIREBASE_ENABLE_AUTH", true),
			EnableFirestore:    getBoolEnv("FIREBASE_ENABLE_FIRESTORE", true),
			EnableRealtimeDB:   getBoolEnv("FIREBASE_ENABLE_REALTIME_DB", true),
			EnableStorage:      getBoolEnv("FIREBASE_ENABLE_STORAGE", true),
			EnableAnalytics:    getBoolEnv("FIREBASE_ENABLE_ANALYTICS", true),
			EnableMessaging:    getBoolEnv("FIREBASE_ENABLE_MESSAGING", true),
			EnableRemoteConfig: getBoolEnv("FIREBASE_ENABLE_REMOTE_CONFIG", false),
			EnableDynamicLinks: getBoolEnv("FIREBASE_ENABLE_DYNAMIC_LINKS", false),
			EnableMLKit:        getBoolEnv("FIREBASE_ENABLE_ML_KIT", false),
			EnablePerformance:  getBoolEnv("FIREBASE_ENABLE_PERFORMANCE", true),
			EnableCrashlytics:  getBoolEnv("FIREBASE_ENABLE_CRASHLYTICS", true),
			EnableAppCheck:     getBoolEnv("FIREBASE_ENABLE_APP_CHECK", true),
			EnableExtensions:   getBoolEnv("FIREBASE_ENABLE_EXTENSIONS", false),
			EnableHosting:      getBoolEnv("FIREBASE_ENABLE_HOSTING", false),
			EnableFunctions:    getBoolEnv("FIREBASE_ENABLE_FUNCTIONS", true),
			EnableEmulators:    getBoolEnv("FIREBASE_ENABLE_EMULATORS", false),
			EmulatorConfig: EmulatorConfig{
				Host:          getEnv("FIREBASE_EMULATOR_HOST", "localhost"),
				AuthPort:      getIntEnv("FIREBASE_AUTH_EMULATOR_PORT", 9099),
				FirestorePort: getIntEnv("FIREBASE_FIRESTORE_EMULATOR_PORT", 8080),
				DatabasePort:  getIntEnv("FIREBASE_DATABASE_EMULATOR_PORT", 9000),
				StoragePort:   getIntEnv("FIREBASE_STORAGE_EMULATOR_PORT", 9199),
				FunctionsPort: getIntEnv("FIREBASE_FUNCTIONS_EMULATOR_PORT", 5001),
				HostingPort:   getIntEnv("FIREBASE_HOSTING_EMULATOR_PORT", 5000),
				PubSubPort:    getIntEnv("FIREBASE_PUBSUB_EMULATOR_PORT", 8085),
				UIPort:        getIntEnv("FIREBASE_UI_EMULATOR_PORT", 4000),
			},
		},
		MCP: MCPConfig{
			CryptoAnalysis: CryptoAnalysisConfig{
				Enabled:        getBoolEnv("MCP_CRYPTO_ANALYSIS_ENABLED", true),
				UpdateInterval: getDurationEnv("MCP_CRYPTO_ANALYSIS_UPDATE_INTERVAL", 30*time.Second),
				Symbols:        getSliceEnv("MCP_CRYPTO_ANALYSIS_SYMBOLS", []string{"BTCUSDT", "ETHUSDT", "ADAUSDT", "SOLUSDT"}),
			},
			SentimentAnalysis: SentimentConfig{
				Enabled:        getBoolEnv("MCP_SENTIMENT_ANALYSIS_ENABLED", true),
				Sources:        getSliceEnv("MCP_SENTIMENT_SOURCES", []string{"reddit", "twitter", "news"}),
				UpdateInterval: getDurationEnv("MCP_SENTIMENT_UPDATE_INTERVAL", 5*time.Minute),
			},
			BrowserAutomation: BrowserConfig{
				Headless:   getBoolEnv("MCP_BROWSER_HEADLESS", true),
				DisableGPU: getBoolEnv("MCP_BROWSER_DISABLE_GPU", true),
				NoSandbox:  getBoolEnv("MCP_BROWSER_NO_SANDBOX", true),
				Timeout:    getDurationEnv("MCP_BROWSER_TIMEOUT", 30*time.Second),
			},
			Firebase: MCPFirebaseConfig{
				Enabled: getBoolEnv("MCP_FIREBASE_ENABLED", true),
				Collections: map[string]string{
					"trading_signals":       "trading_signals",
					"market_data":           "market_data",
					"user_portfolios":       "user_portfolios",
					"trading_strategies":    "trading_strategies",
					"risk_metrics":          "risk_metrics",
					"performance_analytics": "performance_analytics",
				},
				RealtimePaths: map[string]string{
					"live_prices":    "/live_prices",
					"active_orders":  "/active_orders",
					"user_sessions":  "/user_sessions",
					"system_status":  "/system_status",
				},
			},
			Cloudflare: CloudflareConfig{
				Enabled: getBoolEnv("MCP_CLOUDFLARE_ENABLED", false),
			},
			Search: SearchConfig{
				Enabled:   getBoolEnv("MCP_SEARCH_ENABLED", true),
				Providers: getSliceEnv("MCP_SEARCH_PROVIDERS", []string{"google", "bing"}),
			},
			UpdateInterval: getDurationEnv("MCP_UPDATE_INTERVAL", 1*time.Minute),
			EnableRealtime: getBoolEnv("MCP_ENABLE_REALTIME", true),
			BufferSize:     getIntEnv("MCP_BUFFER_SIZE", 1000),
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

// TerminalConfig contains terminal service configuration
type TerminalConfig struct {
	Host         string        `json:"host"`
	Port         int           `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
	MaxSessions  int           `json:"max_sessions"`
	SessionTTL   time.Duration `json:"session_ttl"`
}

// LoggerConfig contains logger configuration
type LoggerConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
}

// FirebaseConfig contains Firebase configuration
type FirebaseConfig struct {
	ProjectID          string         `json:"project_id"`
	CredentialsPath    string         `json:"credentials_path"`
	DatabaseURL        string         `json:"database_url"`
	StorageBucket      string         `json:"storage_bucket"`
	EnableAuth         bool           `json:"enable_auth"`
	EnableFirestore    bool           `json:"enable_firestore"`
	EnableRealtimeDB   bool           `json:"enable_realtime_db"`
	EnableStorage      bool           `json:"enable_storage"`
	EnableAnalytics    bool           `json:"enable_analytics"`
	EnableMessaging    bool           `json:"enable_messaging"`
	EnableRemoteConfig bool           `json:"enable_remote_config"`
	EnableDynamicLinks bool           `json:"enable_dynamic_links"`
	EnableMLKit        bool           `json:"enable_ml_kit"`
	EnablePerformance  bool           `json:"enable_performance"`
	EnableCrashlytics  bool           `json:"enable_crashlytics"`
	EnableAppCheck     bool           `json:"enable_app_check"`
	EnableExtensions   bool           `json:"enable_extensions"`
	EnableHosting      bool           `json:"enable_hosting"`
	EnableFunctions    bool           `json:"enable_functions"`
	EnableEmulators    bool           `json:"enable_emulators"`
	EmulatorConfig     EmulatorConfig `json:"emulator_config"`
}

// EmulatorConfig contains Firebase emulator configuration
type EmulatorConfig struct {
	AuthPort      int    `json:"auth_port"`
	FirestorePort int    `json:"firestore_port"`
	DatabasePort  int    `json:"database_port"`
	StoragePort   int    `json:"storage_port"`
	FunctionsPort int    `json:"functions_port"`
	HostingPort   int    `json:"hosting_port"`
	PubSubPort    int    `json:"pubsub_port"`
	UIPort        int    `json:"ui_port"`
	Host          string `json:"host"`
}

// MCPConfig contains MCP integration configuration
type MCPConfig struct {
	CryptoAnalysis    CryptoAnalysisConfig `json:"crypto_analysis"`
	SentimentAnalysis SentimentConfig      `json:"sentiment_analysis"`
	BrowserAutomation BrowserConfig        `json:"browser_automation"`
	Firebase          MCPFirebaseConfig    `json:"firebase"`
	Cloudflare        CloudflareConfig     `json:"cloudflare"`
	Search            SearchConfig         `json:"search"`
	UpdateInterval    time.Duration        `json:"update_interval"`
	EnableRealtime    bool                 `json:"enable_realtime"`
	BufferSize        int                  `json:"buffer_size"`
}

// CryptoAnalysisConfig contains crypto analysis configuration
type CryptoAnalysisConfig struct {
	Enabled        bool          `json:"enabled"`
	UpdateInterval time.Duration `json:"update_interval"`
	Symbols        []string      `json:"symbols"`
}

// SentimentConfig contains sentiment analysis configuration
type SentimentConfig struct {
	Enabled        bool          `json:"enabled"`
	Sources        []string      `json:"sources"`
	UpdateInterval time.Duration `json:"update_interval"`
}

// MCPFirebaseConfig contains Firebase-specific MCP configuration
type MCPFirebaseConfig struct {
	Enabled       bool              `json:"enabled"`
	Collections   map[string]string `json:"collections"`
	RealtimePaths map[string]string `json:"realtime_paths"`
}

// CloudflareConfig contains Cloudflare configuration
type CloudflareConfig struct {
	Enabled bool `json:"enabled"`
}

// SearchConfig contains search configuration
type SearchConfig struct {
	Enabled   bool     `json:"enabled"`
	Providers []string `json:"providers"`
}
