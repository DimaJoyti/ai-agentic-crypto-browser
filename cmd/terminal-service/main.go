package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ai-agentic-browser/internal/config"
	"github.com/ai-agentic-browser/internal/terminal"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := observability.NewLogger(config.ObservabilityConfig{
		ServiceName: "terminal-service",
		LogLevel:    cfg.Logger.Level,
		LogFormat:   cfg.Logger.Format,
	})

	// Initialize terminal service
	terminalService, err := terminal.NewService(terminal.Config{
		Host:         cfg.Terminal.Host,
		Port:         cfg.Terminal.Port,
		ReadTimeout:  cfg.Terminal.ReadTimeout,
		WriteTimeout: cfg.Terminal.WriteTimeout,
		MaxSessions:  cfg.Terminal.MaxSessions,
		SessionTTL:   cfg.Terminal.SessionTTL,
	}, logger)
	if err != nil {
		logger.Error(context.Background(), "Failed to create terminal service", err)
		os.Exit(1)
	}

	// Setup HTTP router
	router := mux.NewRouter()

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"terminal-service"}`))
	}).Methods("GET")

	// Terminal WebSocket endpoint
	router.HandleFunc("/ws", terminalService.HandleWebSocket).Methods("GET")

	// Terminal API endpoints
	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/sessions", terminalService.HandleCreateSession).Methods("POST")
	api.HandleFunc("/sessions", terminalService.HandleListSessions).Methods("GET")
	api.HandleFunc("/sessions/{sessionId}", terminalService.HandleGetSession).Methods("GET")
	api.HandleFunc("/sessions/{sessionId}", terminalService.HandleDeleteSession).Methods("DELETE")
	api.HandleFunc("/sessions/{sessionId}/history", terminalService.HandleGetHistory).Methods("GET")
	api.HandleFunc("/commands", terminalService.HandleListCommands).Methods("GET")
	api.HandleFunc("/commands/{command}/help", terminalService.HandleGetCommandHelp).Methods("GET")

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Terminal.Host, cfg.Terminal.Port),
		Handler:      handler,
		ReadTimeout:  cfg.Terminal.ReadTimeout,
		WriteTimeout: cfg.Terminal.WriteTimeout,
		IdleTimeout:  cfg.Terminal.IdleTimeout,
	}

	// Start server in goroutine
	go func() {
		logger.Info(context.Background(), "Starting terminal service", map[string]interface{}{
			"host": cfg.Terminal.Host,
			"port": cfg.Terminal.Port,
		})

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(context.Background(), "Failed to start server", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info(context.Background(), "Shutting down terminal service...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error(context.Background(), "Server forced to shutdown", err)
	}

	// Cleanup terminal service
	if err := terminalService.Shutdown(ctx); err != nil {
		logger.Error(context.Background(), "Failed to shutdown terminal service", err)
	}

	logger.Info(context.Background(), "Terminal service stopped")
}
