package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ai-agentic-browser/internal/config"
	"github.com/ai-agentic-browser/pkg/database"
	"github.com/ai-agentic-browser/pkg/middleware"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/gorilla/websocket"
)

// ServiceEndpoints holds the URLs for all microservices
type ServiceEndpoints struct {
	AuthService    string
	AIAgent        string
	BrowserService string
	Web3Service    string
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// In production, implement proper origin checking
		return true
	},
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize observability
	obsConfig := observability.GetDefaultSimpleConfig()
	obsConfig.ServiceName = "api-gateway"
	obsProvider, err := observability.NewSimpleObservabilityProvider(obsConfig)
	if err != nil {
		log.Fatalf("Failed to initialize observability: %v", err)
	}
	logger := obsProvider.Logger

	// Initialize database connections for health checks
	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	redis, err := database.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redis.Close()

	// Define service endpoints
	endpoints := ServiceEndpoints{
		AuthService:    "http://auth-service:8081",
		AIAgent:        "http://ai-agent:8082",
		BrowserService: "http://browser-service:8083",
		Web3Service:    "http://web3-service:8084",
	}

	// In development, use localhost
	if cfg.Server.Host == "0.0.0.0" || cfg.Server.Host == "localhost" {
		endpoints = ServiceEndpoints{
			AuthService:    "http://localhost:8081",
			AIAgent:        "http://localhost:8082",
			BrowserService: "http://localhost:8083",
			Web3Service:    "http://localhost:8084",
		}
	}

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      setupRoutes(endpoints, cfg, logger, db, redis),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		logger.Info(context.Background(), "Starting API Gateway", map[string]interface{}{
			"addr": server.Addr,
		})
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info(context.Background(), "Shutting down API Gateway...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info(context.Background(), "API Gateway stopped")
}

func setupRoutes(endpoints ServiceEndpoints, cfg *config.Config, logger *observability.Logger, db *database.DB, redis *database.RedisClient) http.Handler {
	mux := http.NewServeMux()

	// Apply middleware
	handler := middleware.Recovery(logger)(
		middleware.Logging(logger)(
			middleware.Tracing("api-gateway")(
				middleware.CORS(cfg.Security.CORSAllowedOrigins)(
					middleware.RateLimit(cfg.RateLimit)(mux),
				),
			),
		),
	)

	// Health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		health := map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now(),
			"services":  make(map[string]string),
		}

		// Check database health
		if err := db.Health(ctx); err != nil {
			health["services"].(map[string]string)["database"] = "unhealthy"
			health["status"] = "degraded"
		} else {
			health["services"].(map[string]string)["database"] = "healthy"
		}

		// Check Redis health
		if err := redis.Health(ctx); err != nil {
			health["services"].(map[string]string)["redis"] = "unhealthy"
			health["status"] = "degraded"
		} else {
			health["services"].(map[string]string)["redis"] = "healthy"
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(health)
	})

	// WebSocket endpoint for real-time communication
	mux.HandleFunc("GET /ws", handleWebSocket(logger))

	// API documentation endpoint
	mux.HandleFunc("GET /api/docs", handleAPIDocs())

	// Service status endpoint
	mux.HandleFunc("GET /api/status", handleServiceStatus(endpoints, logger))

	// Proxy routes to microservices
	setupProxyRoutes(mux, endpoints, logger)

	return handler
}

func setupProxyRoutes(mux *http.ServeMux, endpoints ServiceEndpoints, logger *observability.Logger) {
	// Auth service routes
	authURL, _ := url.Parse(endpoints.AuthService)
	authProxy := httputil.NewSingleHostReverseProxy(authURL)
	mux.Handle("/auth/", createProxyHandler(authProxy, "/auth", logger))

	// AI agent routes
	aiURL, _ := url.Parse(endpoints.AIAgent)
	aiProxy := httputil.NewSingleHostReverseProxy(aiURL)
	mux.Handle("/ai/", createProxyHandler(aiProxy, "/ai", logger))

	// Browser service routes
	browserURL, _ := url.Parse(endpoints.BrowserService)
	browserProxy := httputil.NewSingleHostReverseProxy(browserURL)
	mux.Handle("/browser/", createProxyHandler(browserProxy, "/browser", logger))

	// Web3 service routes
	web3URL, _ := url.Parse(endpoints.Web3Service)
	web3Proxy := httputil.NewSingleHostReverseProxy(web3URL)
	mux.Handle("/web3/", createProxyHandler(web3Proxy, "/web3", logger))
}

func createProxyHandler(proxy *httputil.ReverseProxy, prefix string, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log the proxy request
		logger.Info(r.Context(), "Proxying request", map[string]interface{}{
			"method": r.Method,
			"path":   r.URL.Path,
			"prefix": prefix,
		})

		// Modify the request to remove the prefix if needed
		originalPath := r.URL.Path
		if strings.HasPrefix(r.URL.Path, prefix) {
			r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
			if r.URL.Path == "" {
				r.URL.Path = "/"
			}
		}

		// Set headers for service identification
		r.Header.Set("X-Forwarded-For", r.RemoteAddr)
		r.Header.Set("X-Forwarded-Proto", "http")
		r.Header.Set("X-Gateway", "agentic-browser")

		// Custom error handler
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			logger.Error(r.Context(), "Proxy error", err, map[string]interface{}{
				"original_path": originalPath,
				"target_path":   r.URL.Path,
				"prefix":        prefix,
			})

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadGateway)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error":   "Service unavailable",
				"message": "The requested service is currently unavailable",
				"code":    "SERVICE_UNAVAILABLE",
			})
		}

		proxy.ServeHTTP(w, r)
	}
}

func handleWebSocket(logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Error(r.Context(), "WebSocket upgrade failed", err)
			return
		}
		defer conn.Close()

		logger.Info(r.Context(), "WebSocket connection established", map[string]interface{}{
			"remote_addr": r.RemoteAddr,
		})

		// Handle WebSocket messages
		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logger.Error(r.Context(), "WebSocket error", err)
				}
				break
			}

			// Echo the message back (in a real implementation, this would handle different message types)
			response := map[string]interface{}{
				"type":      "echo",
				"message":   string(message),
				"timestamp": time.Now(),
			}

			responseData, _ := json.Marshal(response)
			if err := conn.WriteMessage(messageType, responseData); err != nil {
				logger.Error(r.Context(), "WebSocket write error", err)
				break
			}
		}

		logger.Info(r.Context(), "WebSocket connection closed", map[string]interface{}{
			"remote_addr": r.RemoteAddr,
		})
	}
}

func handleAPIDocs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		docs := map[string]interface{}{
			"title":       "AI Agentic Browser API",
			"version":     "1.0.0",
			"description": "API documentation for the AI-Powered Agentic Crypto Browser",
			"endpoints": map[string]interface{}{
				"authentication": map[string]interface{}{
					"base_url": "/auth",
					"endpoints": []map[string]string{
						{"method": "POST", "path": "/auth/register", "description": "Register a new user"},
						{"method": "POST", "path": "/auth/login", "description": "Login user"},
						{"method": "POST", "path": "/auth/refresh", "description": "Refresh access token"},
						{"method": "POST", "path": "/auth/logout", "description": "Logout user"},
						{"method": "GET", "path": "/auth/me", "description": "Get user profile"},
					},
				},
				"ai_agent": map[string]interface{}{
					"base_url": "/ai",
					"endpoints": []map[string]string{
						{"method": "POST", "path": "/ai/chat", "description": "Send message to AI agent"},
						{"method": "POST", "path": "/ai/tasks", "description": "Create AI task"},
						{"method": "GET", "path": "/ai/tasks/{id}", "description": "Get task status"},
						{"method": "GET", "path": "/ai/conversations", "description": "List conversations"},
					},
				},
				"browser": map[string]interface{}{
					"base_url": "/browser",
					"endpoints": []map[string]string{
						{"method": "POST", "path": "/browser/sessions", "description": "Create browser session"},
						{"method": "POST", "path": "/browser/navigate", "description": "Navigate to URL"},
						{"method": "POST", "path": "/browser/interact", "description": "Interact with page"},
						{"method": "POST", "path": "/browser/extract", "description": "Extract content"},
						{"method": "POST", "path": "/browser/screenshot", "description": "Take screenshot"},
					},
				},
				"web3": map[string]interface{}{
					"base_url": "/web3",
					"endpoints": []map[string]string{
						{"method": "POST", "path": "/web3/connect-wallet", "description": "Connect wallet"},
						{"method": "GET", "path": "/web3/balance", "description": "Get wallet balance"},
						{"method": "POST", "path": "/web3/transaction", "description": "Create transaction"},
						{"method": "GET", "path": "/web3/prices", "description": "Get crypto prices"},
						{"method": "POST", "path": "/web3/defi/interact", "description": "DeFi interaction"},
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(docs)
	}
}

func handleServiceStatus(endpoints ServiceEndpoints, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		status := map[string]interface{}{
			"gateway": map[string]interface{}{
				"status":    "healthy",
				"timestamp": time.Now(),
			},
			"services": make(map[string]interface{}),
		}

		// Check each service health
		services := map[string]string{
			"auth":    endpoints.AuthService + "/health",
			"ai":      endpoints.AIAgent + "/health",
			"browser": endpoints.BrowserService + "/health",
			"web3":    endpoints.Web3Service + "/health",
		}

		for name, healthURL := range services {
			serviceStatus := checkServiceHealth(ctx, healthURL)
			status["services"].(map[string]interface{})[name] = serviceStatus
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	}
}

func checkServiceHealth(ctx context.Context, healthURL string) map[string]interface{} {
	client := &http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequestWithContext(ctx, "GET", healthURL, nil)
	if err != nil {
		return map[string]interface{}{
			"status": "error",
			"error":  err.Error(),
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return map[string]interface{}{
			"status": "unhealthy",
			"error":  err.Error(),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return map[string]interface{}{
			"status":      "healthy",
			"status_code": resp.StatusCode,
		}
	}

	return map[string]interface{}{
		"status":      "unhealthy",
		"status_code": resp.StatusCode,
	}
}
