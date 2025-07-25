package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ai-agentic-browser/internal/ai"
	"github.com/ai-agentic-browser/internal/browser"
	"github.com/ai-agentic-browser/internal/config"
	"github.com/ai-agentic-browser/pkg/database"
	"github.com/ai-agentic-browser/pkg/middleware"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize observability
	logger := observability.NewLogger(cfg.Observability)
	tracingProvider, err := observability.NewTracingProvider(cfg.Observability)
	if err != nil {
		log.Fatalf("Failed to initialize tracing: %v", err)
	}
	defer tracingProvider.Shutdown(context.Background())

	// Initialize database connections
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

	// Initialize browser service
	browserService := browser.NewService(db, redis, cfg.Browser, logger)

	// Initialize AI service
	aiService := ai.NewService(db, redis, cfg.AI, browserService, logger)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, "8082"), // AI Agent port
		Handler:      setupRoutes(aiService, cfg, logger, db),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		logger.Info(context.Background(), "Starting AI agent service", map[string]interface{}{
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

	logger.Info(context.Background(), "Shutting down AI agent service...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info(context.Background(), "AI agent service stopped")
}

func setupRoutes(aiService *ai.Service, cfg *config.Config, logger *observability.Logger, db *database.DB) http.Handler {
	mux := http.NewServeMux()

	// Apply middleware
	handler := middleware.Recovery(logger)(
		middleware.Logging(logger)(
			middleware.Tracing("ai-agent")(
				middleware.CORS(cfg.Security.CORSAllowedOrigins)(
					middleware.RateLimit(cfg.RateLimit)(mux),
				),
			),
		),
	)

	// Health check endpoints
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Check database health
		if err := db.Health(ctx); err != nil {
			http.Error(w, "Database unhealthy", http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	// AI providers health check
	mux.HandleFunc("GET /health/ai", handleAIHealth(aiService, logger))
	mux.HandleFunc("GET /health/ai/{provider}", handleProviderHealth(aiService, logger))
	mux.HandleFunc("POST /health/ai/{provider}/check", handleProviderHealthCheck(aiService, logger))
	mux.HandleFunc("GET /health/ai/{provider}/models", handleProviderModels(aiService, logger))

	// Protected AI endpoints
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("POST /ai/chat", handleChat(aiService, logger))
	protectedMux.HandleFunc("POST /ai/tasks", handleCreateTask(aiService, logger))
	protectedMux.HandleFunc("GET /ai/tasks/{id}", handleGetTask(aiService, logger))
	protectedMux.HandleFunc("GET /ai/conversations", handleListConversations(aiService, logger))
	protectedMux.HandleFunc("GET /ai/conversations/{id}", handleGetConversation(aiService, logger))

	// Apply JWT middleware to protected routes
	mux.Handle("/ai/", middleware.JWT(cfg.JWT.Secret)(protectedMux))

	return handler
}

func handleChat(aiService *ai.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr, ok := middleware.GetUserID(r.Context())
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusInternalServerError)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		var req ai.ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		response, err := aiService.Chat(r.Context(), userID, req)
		if err != nil {
			logger.Error(r.Context(), "Chat request failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func handleCreateTask(aiService *ai.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr, ok := middleware.GetUserID(r.Context())
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusInternalServerError)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		var req ai.TaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		response, err := aiService.CreateTask(r.Context(), userID, req)
		if err != nil {
			logger.Error(r.Context(), "Task creation failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

func handleGetTask(aiService *ai.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskIDStr := r.PathValue("id")
		if taskIDStr == "" {
			http.Error(w, "Task ID is required", http.StatusBadRequest)
			return
		}

		taskID, err := uuid.Parse(taskIDStr)
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		// TODO: Implement GetTask method in AI service
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"task_id": taskID.String(),
			"message": "Task retrieval not implemented yet",
		})
	}
}

func handleListConversations(aiService *ai.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr, ok := middleware.GetUserID(r.Context())
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusInternalServerError)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// TODO: Implement ListConversations method in AI service
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user_id":       userID.String(),
			"conversations": []interface{}{},
			"message":       "Conversation listing not implemented yet",
		})
	}
}

func handleGetConversation(aiService *ai.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conversationIDStr := r.PathValue("id")
		if conversationIDStr == "" {
			http.Error(w, "Conversation ID is required", http.StatusBadRequest)
			return
		}

		conversationID, err := uuid.Parse(conversationIDStr)
		if err != nil {
			http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
			return
		}

		// TODO: Implement GetConversation method in AI service
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"conversation_id": conversationID.String(),
			"message":         "Conversation retrieval not implemented yet",
		})
	}
}

// Health check handlers

func handleAIHealth(aiService *ai.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		healthResponse := aiService.GetHealthStatus()

		// Set appropriate status code based on overall health
		statusCode := http.StatusOK
		if healthResponse.Status == "unhealthy" {
			statusCode = http.StatusServiceUnavailable
		} else if healthResponse.Status == "degraded" {
			statusCode = http.StatusPartialContent
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(healthResponse)
	}
}

func handleProviderHealth(aiService *ai.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		provider := r.PathValue("provider")
		if provider == "" {
			http.Error(w, "Provider name is required", http.StatusBadRequest)
			return
		}

		status, exists := aiService.GetProviderHealth(provider)
		if !exists {
			http.Error(w, "Provider not found", http.StatusNotFound)
			return
		}

		statusCode := http.StatusOK
		if !status.Healthy {
			statusCode = http.StatusServiceUnavailable
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(status)
	}
}

func handleProviderHealthCheck(aiService *ai.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		provider := r.PathValue("provider")
		if provider == "" {
			http.Error(w, "Provider name is required", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		err := aiService.CheckProviderHealth(ctx, provider)
		if err != nil {
			logger.Error(r.Context(), "Provider health check failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Get updated status after check
		status, exists := aiService.GetProviderHealth(provider)
		if !exists {
			http.Error(w, "Provider not found", http.StatusNotFound)
			return
		}

		statusCode := http.StatusOK
		if !status.Healthy {
			statusCode = http.StatusServiceUnavailable
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Health check completed",
			"status":  status,
		})
	}
}

func handleProviderModels(aiService *ai.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		provider := r.PathValue("provider")
		if provider == "" {
			http.Error(w, "Provider name is required", http.StatusBadRequest)
			return
		}

		models, err := aiService.GetProviderModels(provider)
		if err != nil {
			logger.Error(r.Context(), "Failed to get provider models", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"provider": provider,
			"models":   models,
			"count":    len(models),
		})
	}
}
