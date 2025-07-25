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

	"github.com/google/uuid"
	"github.com/ai-agentic-browser/internal/browser"
	"github.com/ai-agentic-browser/internal/config"
	"github.com/ai-agentic-browser/pkg/database"
	"github.com/ai-agentic-browser/pkg/middleware"
	"github.com/ai-agentic-browser/pkg/observability"
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

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, "8083"), // Browser service port
		Handler:      setupRoutes(browserService, cfg, logger, db),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		logger.Info(context.Background(), "Starting browser service", map[string]interface{}{
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

	logger.Info(context.Background(), "Shutting down browser service...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info(context.Background(), "Browser service stopped")
}

func setupRoutes(browserService *browser.Service, cfg *config.Config, logger *observability.Logger, db *database.DB) http.Handler {
	mux := http.NewServeMux()

	// Apply middleware
	handler := middleware.Recovery(logger)(
		middleware.Logging(logger)(
			middleware.Tracing("browser-service")(
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

		// Check database health
		if err := db.Health(ctx); err != nil {
			http.Error(w, "Database unhealthy", http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	// Protected browser endpoints
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("POST /browser/sessions", handleCreateSession(browserService, logger))
	protectedMux.HandleFunc("GET /browser/sessions", handleListSessions(browserService, logger))
	protectedMux.HandleFunc("POST /browser/navigate", handleNavigate(browserService, logger))
	protectedMux.HandleFunc("POST /browser/interact", handleInteract(browserService, logger))
	protectedMux.HandleFunc("POST /browser/extract", handleExtract(browserService, logger))
	protectedMux.HandleFunc("POST /browser/screenshot", handleScreenshot(browserService, logger))

	// Apply JWT middleware to protected routes
	mux.Handle("/browser/", middleware.JWT(cfg.JWT.Secret)(protectedMux))

	return handler
}

func handleCreateSession(browserService *browser.Service, logger *observability.Logger) http.HandlerFunc {
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

		var req browser.SessionCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		session, err := browserService.CreateSession(r.Context(), userID, req)
		if err != nil {
			logger.Error(r.Context(), "Session creation failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := browser.SessionCreateResponse{Session: *session}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

func handleListSessions(browserService *browser.Service, logger *observability.Logger) http.HandlerFunc {
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

		// TODO: Implement ListSessions method in browser service
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user_id":  userID.String(),
			"sessions": []interface{}{},
			"message":  "Session listing not implemented yet",
		})
	}
}

func handleNavigate(browserService *browser.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionIDStr := r.Header.Get("X-Session-ID")
		if sessionIDStr == "" {
			http.Error(w, "Session ID header required", http.StatusBadRequest)
			return
		}

		sessionID, err := uuid.Parse(sessionIDStr)
		if err != nil {
			http.Error(w, "Invalid session ID", http.StatusBadRequest)
			return
		}

		var req browser.NavigateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		response, err := browserService.Navigate(r.Context(), sessionID, req)
		if err != nil {
			logger.Error(r.Context(), "Navigation failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func handleInteract(browserService *browser.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionIDStr := r.Header.Get("X-Session-ID")
		if sessionIDStr == "" {
			http.Error(w, "Session ID header required", http.StatusBadRequest)
			return
		}

		sessionID, err := uuid.Parse(sessionIDStr)
		if err != nil {
			http.Error(w, "Invalid session ID", http.StatusBadRequest)
			return
		}

		var req browser.InteractRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		response, err := browserService.Interact(r.Context(), sessionID, req)
		if err != nil {
			logger.Error(r.Context(), "Interaction failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func handleExtract(browserService *browser.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionIDStr := r.Header.Get("X-Session-ID")
		if sessionIDStr == "" {
			http.Error(w, "Session ID header required", http.StatusBadRequest)
			return
		}

		sessionID, err := uuid.Parse(sessionIDStr)
		if err != nil {
			http.Error(w, "Invalid session ID", http.StatusBadRequest)
			return
		}

		var req browser.ExtractRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		response, err := browserService.Extract(r.Context(), sessionID, req)
		if err != nil {
			logger.Error(r.Context(), "Content extraction failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func handleScreenshot(browserService *browser.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionIDStr := r.Header.Get("X-Session-ID")
		if sessionIDStr == "" {
			http.Error(w, "Session ID header required", http.StatusBadRequest)
			return
		}

		sessionID, err := uuid.Parse(sessionIDStr)
		if err != nil {
			http.Error(w, "Invalid session ID", http.StatusBadRequest)
			return
		}

		var req browser.ScreenshotRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		response, err := browserService.TakeScreenshot(r.Context(), sessionID, req)
		if err != nil {
			logger.Error(r.Context(), "Screenshot failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
