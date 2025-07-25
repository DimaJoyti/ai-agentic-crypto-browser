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

	"github.com/ai-agentic-browser/internal/auth"
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

	// Initialize auth service
	authService := auth.NewService(db, redis, cfg.JWT, logger)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      setupRoutes(authService, cfg, logger, db),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		logger.Info(context.Background(), "Starting auth service", map[string]interface{}{
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

	logger.Info(context.Background(), "Shutting down auth service...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info(context.Background(), "Auth service stopped")
}

func setupRoutes(authService *auth.Service, cfg *config.Config, logger *observability.Logger, db *database.DB) http.Handler {
	mux := http.NewServeMux()

	// Apply middleware
	handler := middleware.Recovery(logger)(
		middleware.Logging(logger)(
			middleware.Tracing("auth-service")(
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

	// Authentication endpoints
	mux.HandleFunc("POST /auth/register", handleRegister(authService, logger))
	mux.HandleFunc("POST /auth/login", handleLogin(authService, logger))
	mux.HandleFunc("POST /auth/refresh", handleRefreshToken(authService, logger))
	mux.HandleFunc("POST /auth/logout", handleLogout(authService, logger))

	// Protected endpoints
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("GET /auth/me", handleGetProfile(authService, logger))
	protectedMux.HandleFunc("PUT /auth/me", handleUpdateProfile(authService, logger))
	protectedMux.HandleFunc("POST /auth/change-password", handleChangePassword(authService, logger))

	// Apply JWT middleware to protected routes
	mux.Handle("/auth/me", middleware.JWT(cfg.JWT.Secret)(protectedMux))
	mux.Handle("/auth/change-password", middleware.JWT(cfg.JWT.Secret)(protectedMux))

	return handler
}

func handleRegister(authService *auth.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req auth.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		user, err := authService.Register(r.Context(), req)
		if err != nil {
			logger.Error(r.Context(), "Registration failed", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}

func handleLogin(authService *auth.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req auth.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		userAgent := r.Header.Get("User-Agent")
		ipAddress := r.RemoteAddr

		response, err := authService.Login(r.Context(), req, userAgent, ipAddress)
		if err != nil {
			logger.Error(r.Context(), "Login failed", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func handleRefreshToken(authService *auth.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req auth.RefreshTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		response, err := authService.RefreshToken(r.Context(), req.RefreshToken)
		if err != nil {
			logger.Error(r.Context(), "Token refresh failed", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func handleLogout(authService *auth.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req auth.RefreshTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err := authService.Logout(r.Context(), req.RefreshToken)
		if err != nil {
			logger.Error(r.Context(), "Logout failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
	}
}

func handleGetProfile(authService *auth.Service, logger *observability.Logger) http.HandlerFunc {
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

		user, err := authService.GetUserByID(r.Context(), userID)
		if err != nil {
			logger.Error(r.Context(), "Failed to get user profile", err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func handleUpdateProfile(authService *auth.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implementation for profile update
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Profile update not implemented yet"})
	}
}

func handleChangePassword(authService *auth.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implementation for password change
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Password change not implemented yet"})
	}
}
