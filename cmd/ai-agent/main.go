package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ai-agentic-browser/internal/ai"
	"github.com/ai-agentic-browser/internal/browser"
	"github.com/ai-agentic-browser/internal/config"
	"github.com/ai-agentic-browser/pkg/database"
	"github.com/ai-agentic-browser/pkg/middleware"
	"github.com/ai-agentic-browser/pkg/ml"
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

	// Initialize optimized database connections
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

	// Initialize performance monitoring
	perfMonitor := observability.NewPerformanceMonitor(logger)
	defer perfMonitor.Stop()

	// Initialize caching middleware
	cacheMiddleware := middleware.NewCacheMiddleware(redis, logger)

	logger.Info(context.Background(), "Database and caching optimizations initialized", map[string]interface{}{
		"db_max_open_conns":      cfg.Database.MaxOpenConns,
		"db_max_idle_conns":      cfg.Database.MaxIdleConns,
		"redis_pool_size":        cfg.Redis.PoolSize,
		"cache_enabled":          true,
		"performance_monitoring": true,
	})

	// Initialize browser service
	browserService := browser.NewService(db, redis, cfg.Browser, logger)

	// Initialize enhanced AI components
	enhancedAI := ai.NewEnhancedAIService(logger)
	multiModalEngine := ai.NewMultiModalEngine(logger)
	userBehaviorEngine := ai.NewUserBehaviorLearningEngine(logger)
	marketAdaptationEngine := ai.NewMarketAdaptationEngine(logger)
	voiceInterface := ai.NewVoiceInterface(logger, nil, nil, nil)
	conversationalAI := ai.NewConversationalAI(logger, nil, nil, nil)
	cryptoCoinAnalyzer := ai.NewCryptoCoinAnalyzer(logger)

	logger.Info(context.Background(), "AI services initialized", map[string]interface{}{
		"enhanced_ai":       enhancedAI != nil,
		"multimodal_engine": multiModalEngine != nil,
		"voice_interface":   voiceInterface != nil,
		"conversational_ai": conversationalAI != nil,
	})

	// Create HTTP server with performance optimizations
	handler := setupRoutes(browserService, enhancedAI, multiModalEngine, userBehaviorEngine, marketAdaptationEngine, voiceInterface, conversationalAI, cryptoCoinAnalyzer, cfg, logger, db, perfMonitor, cacheMiddleware)

	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%s", cfg.Server.Host, "8082"), // AI Agent port
		Handler:        handler,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		IdleTimeout:    cfg.Server.IdleTimeout,
		MaxHeaderBytes: 1 << 20, // 1MB
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

func setupRoutes(
	browserService *browser.Service,
	enhancedAI *ai.EnhancedAIService,
	multiModalEngine *ai.MultiModalEngine,
	userBehaviorEngine *ai.UserBehaviorLearningEngine,
	marketAdaptationEngine *ai.MarketAdaptationEngine,
	voiceInterface *ai.VoiceInterface,
	conversationalAI *ai.ConversationalAI,
	cryptoCoinAnalyzer *ai.CryptoCoinAnalyzer,
	cfg *config.Config,
	logger *observability.Logger,
	db *database.DB,
	perfMonitor *observability.PerformanceMonitor,
	cacheMiddleware *middleware.CacheMiddleware,
) http.Handler {
	mux := http.NewServeMux()

	// Apply middleware stack with performance optimizations
	handler := middleware.Recovery(logger)(
		middleware.Logging(logger)(
			middleware.Tracing("ai-agent")(
				cacheMiddleware.Middleware()(
					middleware.CORS(cfg.Security.CORSAllowedOrigins)(
						middleware.RateLimit(cfg.RateLimit)(mux),
					),
				),
			),
		),
	)

	// Health check endpoints with performance metrics
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Check database health
		if err := db.Health(ctx); err != nil {
			http.Error(w, "Database unhealthy", http.StatusServiceUnavailable)
			return
		}

		// Get performance status
		healthStatus := perfMonitor.GetHealthStatus()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(healthStatus)
	})

	// Performance metrics endpoint
	mux.HandleFunc("GET /metrics", func(w http.ResponseWriter, r *http.Request) {
		metrics := perfMonitor.GetMetrics()
		dbMetrics := db.GetMetrics()
		cacheMetrics := cacheMiddleware.GetStats()

		response := map[string]interface{}{
			"performance": metrics,
			"database":    dbMetrics,
			"cache":       cacheMetrics,
			"timestamp":   time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Database metrics endpoint
	mux.HandleFunc("GET /metrics/database", func(w http.ResponseWriter, r *http.Request) {
		metrics := db.GetMetrics()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(metrics)
	})

	// Cache metrics endpoint
	mux.HandleFunc("GET /metrics/cache", func(w http.ResponseWriter, r *http.Request) {
		metrics := cacheMiddleware.GetStats()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(metrics)
	})

	// AI providers health check (simplified for new architecture)
	mux.HandleFunc("GET /health/ai", handleAIHealth(conversationalAI, logger))
	mux.HandleFunc("GET /health/ai/{provider}", handleProviderHealth(conversationalAI, logger))
	mux.HandleFunc("POST /health/ai/{provider}/check", handleProviderHealthCheck(conversationalAI, logger))
	mux.HandleFunc("GET /health/ai/{provider}/models", handleProviderModels(conversationalAI, logger))

	// Protected AI endpoints (enhanced)
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("POST /ai/chat", handleChat(conversationalAI, logger))
	protectedMux.HandleFunc("POST /ai/voice/command", handleVoiceCommandSimple(voiceInterface, logger))
	protectedMux.HandleFunc("POST /ai/conversations/start", handleStartConversationSimple(conversationalAI, logger))

	// Enhanced AI endpoints
	protectedMux.HandleFunc("POST /ai/analyze", handleEnhancedAnalysis(enhancedAI, logger))
	protectedMux.HandleFunc("POST /ai/predict/price", handlePricePrediction(enhancedAI, logger))
	protectedMux.HandleFunc("POST /ai/analyze/sentiment", handleSentimentAnalysis(enhancedAI, logger))
	protectedMux.HandleFunc("POST /ai/analytics/predictive", handlePredictiveAnalytics(enhancedAI, logger))
	protectedMux.HandleFunc("GET /ai/models/status", handleModelStatus(enhancedAI, logger))
	protectedMux.HandleFunc("POST /ai/models/train", handleModelTraining(enhancedAI, logger))
	protectedMux.HandleFunc("POST /ai/models/feedback", handleModelFeedback(enhancedAI, logger))

	// Learning and adaptation endpoints
	protectedMux.HandleFunc("POST /ai/learning/behavior", handleUserBehaviorLearning(enhancedAI, logger))
	protectedMux.HandleFunc("GET /ai/learning/profile", handleGetUserProfile(enhancedAI, logger))
	protectedMux.HandleFunc("GET /ai/learning/patterns", handleGetMarketPatterns(enhancedAI, logger))
	protectedMux.HandleFunc("GET /ai/learning/performance", handleGetPerformanceMetrics(enhancedAI, logger))
	protectedMux.HandleFunc("POST /ai/adaptation/request", handleRequestAdaptation(enhancedAI, logger))
	protectedMux.HandleFunc("GET /ai/adaptation/models", handleGetAdaptiveModels(enhancedAI, logger))
	protectedMux.HandleFunc("GET /ai/adaptation/history/{modelId}", handleGetAdaptationHistory(enhancedAI, logger))

	// Advanced NLP endpoints
	protectedMux.HandleFunc("POST /ai/nlp/analyze", handleAdvancedNLP(enhancedAI, logger))

	// Decision engine endpoints
	protectedMux.HandleFunc("POST /ai/decisions/request", handleDecisionRequest(enhancedAI, logger))
	protectedMux.HandleFunc("GET /ai/decisions/active", handleGetActiveDecisions(enhancedAI, logger))
	protectedMux.HandleFunc("GET /ai/decisions/history", handleGetDecisionHistory(enhancedAI, logger))
	protectedMux.HandleFunc("GET /ai/decisions/performance", handleGetDecisionPerformance(enhancedAI, logger))

	// Multi-Modal AI endpoints
	protectedMux.HandleFunc("POST /ai/multimodal/analyze", handleMultiModalAnalysis(multiModalEngine, logger))
	protectedMux.HandleFunc("POST /ai/multimodal/image", handleImageAnalysis(multiModalEngine, logger))
	protectedMux.HandleFunc("POST /ai/multimodal/document", handleDocumentAnalysis(multiModalEngine, logger))
	protectedMux.HandleFunc("POST /ai/multimodal/audio", handleAudioAnalysis(multiModalEngine, logger))
	protectedMux.HandleFunc("POST /ai/multimodal/chart", handleChartAnalysis(multiModalEngine, logger))
	protectedMux.HandleFunc("GET /ai/multimodal/formats", handleGetSupportedFormats(multiModalEngine, logger))

	// User Behavior Learning endpoints
	protectedMux.HandleFunc("POST /ai/behavior/learn", handleLearnFromBehavior(userBehaviorEngine, logger))
	protectedMux.HandleFunc("GET /ai/behavior/profile", handleGetUserBehaviorProfile(userBehaviorEngine, logger))
	protectedMux.HandleFunc("GET /ai/behavior/recommendations", handleGetRecommendations(userBehaviorEngine, logger))
	protectedMux.HandleFunc("GET /ai/behavior/history", handleGetBehaviorHistory(userBehaviorEngine, logger))
	protectedMux.HandleFunc("PUT /ai/behavior/recommendation/{id}/status", handleUpdateRecommendationStatus(userBehaviorEngine, logger))
	protectedMux.HandleFunc("GET /ai/behavior/models", handleGetLearningModels(userBehaviorEngine, logger))

	// Market Pattern Adaptation endpoints
	protectedMux.HandleFunc("POST /ai/market/patterns/detect", handleDetectMarketPatterns(marketAdaptationEngine, logger))
	protectedMux.HandleFunc("GET /ai/market/patterns", handleGetMarketPatternsAdaptation(marketAdaptationEngine, logger))
	protectedMux.HandleFunc("POST /ai/market/strategies/adapt", handleAdaptStrategies(marketAdaptationEngine, logger))
	protectedMux.HandleFunc("GET /ai/market/strategies", handleGetAdaptiveStrategies(marketAdaptationEngine, logger))
	protectedMux.HandleFunc("POST /ai/market/strategies", handleAddAdaptiveStrategy(marketAdaptationEngine, logger))
	protectedMux.HandleFunc("PUT /ai/market/strategies/{id}/status", handleUpdateStrategyStatus(marketAdaptationEngine, logger))
	protectedMux.HandleFunc("GET /ai/market/adaptation/history", handleGetMarketAdaptationHistory(marketAdaptationEngine, logger))
	protectedMux.HandleFunc("GET /ai/market/performance/{strategy_id}", handleGetStrategyPerformanceMetrics(marketAdaptationEngine, logger))

	// Crypto Coin Analyzer endpoints
	protectedMux.HandleFunc("POST /ai/crypto/analyze/{symbol}", handleCryptoCoinAnalysis(cryptoCoinAnalyzer, logger))
	protectedMux.HandleFunc("GET /ai/crypto/analyze/{symbol}", handleCryptoCoinAnalysis(cryptoCoinAnalyzer, logger))
	protectedMux.HandleFunc("POST /ai/crypto/report/{symbol}", handleCryptoCoinReport(cryptoCoinAnalyzer, logger))
	protectedMux.HandleFunc("GET /ai/crypto/report/{symbol}", handleCryptoCoinReport(cryptoCoinAnalyzer, logger))

	// Apply JWT middleware to protected routes
	mux.Handle("/ai/", middleware.JWT(cfg.JWT.Secret)(protectedMux))

	return handler
}

// Enhanced AI handlers

func handleEnhancedAnalysis(enhancedAI *ai.EnhancedAIService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("user_id").(uuid.UUID)
		if !ok {
			http.Error(w, "User ID not found", http.StatusUnauthorized)
			return
		}

		var req ai.AIRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		req.UserID = userID
		req.RequestID = uuid.New().String()
		req.RequestedAt = time.Now()

		response, err := enhancedAI.ProcessRequest(r.Context(), &req)
		if err != nil {
			logger.Error(r.Context(), "Enhanced AI analysis failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func handlePricePrediction(enhancedAI *ai.EnhancedAIService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("user_id").(uuid.UUID)
		if !ok {
			http.Error(w, "User ID not found", http.StatusUnauthorized)
			return
		}

		var predictionReq ai.PricePredictionRequest
		if err := json.NewDecoder(r.Body).Decode(&predictionReq); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Create AI request
		aiReq := &ai.AIRequest{
			RequestID: uuid.New().String(),
			UserID:    userID,
			Type:      "price_prediction",
			Symbol:    predictionReq.Symbol,
			Data: map[string]interface{}{
				"price_prediction_request": &predictionReq,
			},
			Options: ai.AIRequestOptions{
				IncludePredictions: true,
				TimeHorizon:        predictionReq.Horizon,
			},
			RequestedAt: time.Now(),
		}

		response, err := enhancedAI.ProcessRequest(r.Context(), aiReq)
		if err != nil {
			logger.Error(r.Context(), "Price prediction failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.PricePrediction)
	}
}

func handleSentimentAnalysis(enhancedAI *ai.EnhancedAIService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("user_id").(uuid.UUID)
		if !ok {
			http.Error(w, "User ID not found", http.StatusUnauthorized)
			return
		}

		var sentimentReq ai.SentimentRequest
		if err := json.NewDecoder(r.Body).Decode(&sentimentReq); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Create AI request
		aiReq := &ai.AIRequest{
			RequestID: uuid.New().String(),
			UserID:    userID,
			Type:      "sentiment_analysis",
			Data: map[string]interface{}{
				"sentiment_request": &sentimentReq,
			},
			Options: ai.AIRequestOptions{
				IncludeSentiment: true,
			},
			RequestedAt: time.Now(),
		}

		response, err := enhancedAI.ProcessRequest(r.Context(), aiReq)
		if err != nil {
			logger.Error(r.Context(), "Sentiment analysis failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.SentimentAnalysis)
	}
}

func handleModelStatus(enhancedAI *ai.EnhancedAIService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := enhancedAI.GetModelStatus(r.Context())

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"models":    status,
			"timestamp": time.Now(),
		})
	}
}

func handleModelTraining(enhancedAI *ai.EnhancedAIService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ModelID string          `json:"model_id"`
			Data    ml.TrainingData `json:"data"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err := enhancedAI.TrainModel(r.Context(), req.ModelID, req.Data)
		if err != nil {
			logger.Error(r.Context(), "Model training failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":  true,
			"message":  "Model training started",
			"model_id": req.ModelID,
		})
	}
}

func handleModelFeedback(enhancedAI *ai.EnhancedAIService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ModelID  string                `json:"model_id"`
			Feedback ml.PredictionFeedback `json:"feedback"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err := enhancedAI.ProvideFeedback(r.Context(), req.ModelID, &req.Feedback)
		if err != nil {
			logger.Error(r.Context(), "Model feedback failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":  true,
			"message":  "Feedback processed",
			"model_id": req.ModelID,
		})
	}
}

func handlePredictiveAnalytics(enhancedAI *ai.EnhancedAIService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value("user_id").(uuid.UUID)
		if !ok {
			http.Error(w, "User ID not found", http.StatusUnauthorized)
			return
		}

		var predictiveReq ai.PredictiveRequest
		if err := json.NewDecoder(r.Body).Decode(&predictiveReq); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		predictiveReq.RequestedAt = time.Now()

		response, err := enhancedAI.GeneratePredictiveAnalytics(r.Context(), &predictiveReq)
		if err != nil {
			logger.Error(r.Context(), "Predictive analytics failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// Learning and adaptation handlers

func handleUserBehaviorLearning(enhancedAI *ai.EnhancedAIService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("user_id").(uuid.UUID)
		if !ok {
			http.Error(w, "User ID not found", http.StatusUnauthorized)
			return
		}

		var behaviorData ai.UserBehaviorData
		if err := json.NewDecoder(r.Body).Decode(&behaviorData); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		behaviorData.Timestamp = time.Now()

		err := enhancedAI.LearnFromUserBehavior(r.Context(), userID, &behaviorData)
		if err != nil {
			logger.Error(r.Context(), "User behavior learning failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "User behavior learned successfully",
			"user_id": userID,
		})
	}
}

func handleGetUserProfile(enhancedAI *ai.EnhancedAIService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("user_id").(uuid.UUID)
		if !ok {
			http.Error(w, "User ID not found", http.StatusUnauthorized)
			return
		}

		profile, err := enhancedAI.GetUserProfile(userID)
		if err != nil {
			logger.Error(r.Context(), "Failed to get user profile", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(profile)
	}
}

func handleGetMarketPatterns(enhancedAI *ai.EnhancedAIService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		patterns := enhancedAI.GetMarketPatterns()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"patterns":  patterns,
			"count":     len(patterns),
			"timestamp": time.Now(),
		})
	}
}

func handleGetPerformanceMetrics(enhancedAI *ai.EnhancedAIService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics := enhancedAI.GetPerformanceMetrics()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"metrics":   metrics,
			"count":     len(metrics),
			"timestamp": time.Now(),
		})
	}
}

func handleRequestAdaptation(enhancedAI *ai.EnhancedAIService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("user_id").(uuid.UUID)
		if !ok {
			http.Error(w, "User ID not found", http.StatusUnauthorized)
			return
		}

		var adaptationReq ai.AdaptationRequest
		if err := json.NewDecoder(r.Body).Decode(&adaptationReq); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		adaptationReq.UserID = userID
		adaptationReq.RequestedAt = time.Now()

		err := enhancedAI.RequestModelAdaptation(&adaptationReq)
		if err != nil {
			logger.Error(r.Context(), "Model adaptation request failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":    true,
			"message":    "Adaptation request submitted",
			"model_id":   adaptationReq.ModelID,
			"request_id": uuid.New().String(),
		})
	}
}

func handleGetAdaptiveModels(enhancedAI *ai.EnhancedAIService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		models := enhancedAI.GetAdaptiveModels()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"models":    models,
			"count":     len(models),
			"timestamp": time.Now(),
		})
	}
}

func handleGetAdaptationHistory(enhancedAI *ai.EnhancedAIService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		modelID := r.PathValue("modelId")
		if modelID == "" {
			http.Error(w, "Model ID is required", http.StatusBadRequest)
			return
		}

		history, err := enhancedAI.GetAdaptationHistory(modelID)
		if err != nil {
			logger.Error(r.Context(), "Failed to get adaptation history", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"model_id":  modelID,
			"history":   history,
			"count":     len(history),
			"timestamp": time.Now(),
		})
	}
}

func handleAdvancedNLP(enhancedAI *ai.EnhancedAIService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("user_id").(uuid.UUID)
		if !ok {
			http.Error(w, "User ID not found", http.StatusUnauthorized)
			return
		}

		var nlpReq ai.NLPRequest
		if err := json.NewDecoder(r.Body).Decode(&nlpReq); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Set request metadata
		if nlpReq.RequestID == "" {
			nlpReq.RequestID = uuid.New().String()
		}
		nlpReq.RequestedAt = time.Now()

		result, err := enhancedAI.ProcessAdvancedNLP(r.Context(), &nlpReq)
		if err != nil {
			logger.Error(r.Context(), "Advanced NLP processing failed", err, map[string]interface{}{
				"user_id":    userID,
				"request_id": nlpReq.RequestID,
			})
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Learn from user NLP usage
		behaviorData := &ai.UserBehaviorData{
			Type:      "analysis_request",
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"analysis_type": "advanced_nlp",
				"text_count":    len(nlpReq.Texts),
				"sources":       nlpReq.Sources,
				"options":       nlpReq.Options,
			},
			Context: map[string]interface{}{
				"request_id": nlpReq.RequestID,
			},
			Outcome:     "success",
			Performance: 0.0, // Neutral for analysis requests
		}

		// Learn from user behavior (non-blocking)
		go func() {
			if err := enhancedAI.LearnFromUserBehavior(context.Background(), userID, behaviorData); err != nil {
				logger.Warn(context.Background(), "Failed to learn from NLP behavior", map[string]interface{}{
					"error":   err.Error(),
					"user_id": userID,
				})
			}
		}()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

func handleDecisionRequest(enhancedAI *ai.EnhancedAIService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("user_id").(uuid.UUID)
		if !ok {
			http.Error(w, "User ID not found", http.StatusUnauthorized)
			return
		}

		var decisionReq ai.DecisionRequest
		if err := json.NewDecoder(r.Body).Decode(&decisionReq); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Set request metadata
		if decisionReq.RequestID == "" {
			decisionReq.RequestID = uuid.New().String()
		}
		decisionReq.UserID = userID
		decisionReq.RequestedAt = time.Now()

		// Set default expiry if not provided
		if decisionReq.ExpiresAt.IsZero() {
			decisionReq.ExpiresAt = time.Now().Add(24 * time.Hour)
		}

		result, err := enhancedAI.ProcessDecisionRequest(r.Context(), &decisionReq)
		if err != nil {
			logger.Error(r.Context(), "Decision request processing failed", err, map[string]interface{}{
				"user_id":       userID,
				"request_id":    decisionReq.RequestID,
				"decision_type": decisionReq.DecisionType,
			})
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Learn from user decision request
		behaviorData := &ai.UserBehaviorData{
			Type:      "decision_request",
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"decision_type":   decisionReq.DecisionType,
				"confidence":      result.Confidence,
				"auto_executable": result.AutoExecutable,
			},
			Context: map[string]interface{}{
				"request_id":  decisionReq.RequestID,
				"decision_id": result.DecisionID,
			},
			Outcome:     "success",
			Performance: result.Confidence, // Use confidence as performance metric
		}

		// Learn from user behavior (non-blocking)
		go func() {
			if err := enhancedAI.LearnFromUserBehavior(context.Background(), userID, behaviorData); err != nil {
				logger.Warn(context.Background(), "Failed to learn from decision behavior", map[string]interface{}{
					"error":   err.Error(),
					"user_id": userID,
				})
			}
		}()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

func handleGetActiveDecisions(enhancedAI *ai.EnhancedAIService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("user_id").(uuid.UUID)
		if !ok {
			http.Error(w, "User ID not found", http.StatusUnauthorized)
			return
		}

		activeDecisions := enhancedAI.GetActiveDecisions(userID)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"active_decisions": activeDecisions,
			"count":            len(activeDecisions),
			"user_id":          userID,
			"timestamp":        time.Now(),
		})
	}
}

func handleGetDecisionHistory(enhancedAI *ai.EnhancedAIService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("user_id").(uuid.UUID)
		if !ok {
			http.Error(w, "User ID not found", http.StatusUnauthorized)
			return
		}

		// Parse limit parameter
		limit := 50 // Default limit
		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
				limit = parsedLimit
			}
		}

		history := enhancedAI.GetDecisionHistory(userID, limit)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"decision_history": history,
			"count":            len(history),
			"limit":            limit,
			"user_id":          userID,
			"timestamp":        time.Now(),
		})
	}
}

func handleGetDecisionPerformance(enhancedAI *ai.EnhancedAIService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics := enhancedAI.GetDecisionPerformanceMetrics()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"performance_metrics": metrics,
			"timestamp":           time.Now(),
		})
	}
}

func handleChat(conversationalAI *ai.ConversationalAI, logger *observability.Logger) http.HandlerFunc {
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

		var req struct {
			Message string `json:"message"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		response, err := conversationalAI.ProcessMessage(r.Context(), userID, req.Message)
		if err != nil {
			logger.Error(r.Context(), "Chat request failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func handleVoiceCommandSimple(voiceInterface *ai.VoiceInterface, logger *observability.Logger) http.HandlerFunc {
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

		var req struct {
			Text      string `json:"text"`
			AudioData []byte `json:"audio_data,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		response, err := voiceInterface.ProcessVoiceCommand(r.Context(), userID, req.AudioData, req.Text)
		if err != nil {
			logger.Error(r.Context(), "Voice command processing failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func handleStartConversationSimple(conversationalAI *ai.ConversationalAI, logger *observability.Logger) http.HandlerFunc {
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

		conversation, err := conversationalAI.StartConversation(r.Context(), userID)
		if err != nil {
			logger.Error(r.Context(), "Conversation start failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(conversation)
	}
}

// Health check handlers (simplified)

func handleAIHealth(conversationalAI *ai.ConversationalAI, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "healthy",
			"service":   "ai-agent",
			"timestamp": time.Now(),
		})
	}
}

func handleProviderHealth(conversationalAI *ai.ConversationalAI, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		provider := r.PathValue("provider")
		if provider == "" {
			http.Error(w, "Provider name is required", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"provider":  provider,
			"status":    "healthy",
			"timestamp": time.Now(),
		})
	}
}

func handleProviderHealthCheck(conversationalAI *ai.ConversationalAI, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		provider := r.PathValue("provider")
		if provider == "" {
			http.Error(w, "Provider name is required", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":   "Health check completed",
			"provider":  provider,
			"status":    "healthy",
			"timestamp": time.Now(),
		})
	}
}

func handleProviderModels(conversationalAI *ai.ConversationalAI, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		provider := r.PathValue("provider")
		if provider == "" {
			http.Error(w, "Provider name is required", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"provider": provider,
			"models":   []string{"gpt-3.5-turbo", "gpt-4"},
			"count":    2,
		})
	}
}

// Multi-Modal AI handlers

func handleMultiModalAnalysis(engine *ai.MultiModalEngine, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var req ai.MultiModalRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Get user ID from context
		userID, err := getUserIDFromContext(ctx)
		if err != nil {
			http.Error(w, "User ID required", http.StatusUnauthorized)
			return
		}
		req.UserID = userID

		// Set request ID if not provided
		if req.RequestID == "" {
			req.RequestID = uuid.New().String()
		}

		result, err := engine.ProcessMultiModalRequest(ctx, &req)
		if err != nil {
			logger.Error(ctx, "Multi-modal analysis failed", err, map[string]interface{}{
				"request_id": req.RequestID,
			})
			http.Error(w, "Analysis failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)

		logger.Info(ctx, "Multi-modal analysis completed", map[string]interface{}{
			"request_id":      req.RequestID,
			"content_count":   len(req.Content),
			"processing_time": result.ProcessingTime.Milliseconds(),
		})
	}
}

func handleImageAnalysis(engine *ai.MultiModalEngine, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Parse multipart form
		err := r.ParseMultipartForm(10 << 20) // 10MB max
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Image file required", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Get user ID from context
		userID, err := getUserIDFromContext(ctx)
		if err != nil {
			http.Error(w, "User ID required", http.StatusUnauthorized)
			return
		}

		// Validate image format
		if !engine.ValidateImageFormat(header.Filename) {
			http.Error(w, "Unsupported image format", http.StatusBadRequest)
			return
		}

		// Parse options
		options := ai.MultiModalOptions{
			AnalyzeImages:   true,
			ExtractText:     r.FormValue("extract_text") == "true",
			AnalyzeCharts:   r.FormValue("analyze_charts") == "true",
			DetectObjects:   r.FormValue("detect_objects") == "true",
			GenerateSummary: r.FormValue("generate_summary") == "true",
		}

		result, err := engine.ProcessImageFile(ctx, userID, file, header, options)
		if err != nil {
			logger.Error(ctx, "Image analysis failed", err, map[string]interface{}{
				"filename": header.Filename,
			})
			http.Error(w, "Image analysis failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)

		logger.Info(ctx, "Image analysis completed", map[string]interface{}{
			"filename":        header.Filename,
			"processing_time": result.ProcessingTime.Milliseconds(),
		})
	}
}

func handleDocumentAnalysis(engine *ai.MultiModalEngine, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Parse multipart form
		err := r.ParseMultipartForm(50 << 20) // 50MB max
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("document")
		if err != nil {
			http.Error(w, "Document file required", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Get user ID from context
		userID, err := getUserIDFromContext(ctx)
		if err != nil {
			http.Error(w, "User ID required", http.StatusUnauthorized)
			return
		}

		// Read file data and create request
		data, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}

		encodedData := base64.StdEncoding.EncodeToString(data)

		req := &ai.MultiModalRequest{
			RequestID: uuid.New().String(),
			UserID:    userID,
			Type:      "document",
			Content: []ai.MultiModalContent{
				{
					ID:       uuid.New().String(),
					Type:     "document",
					Data:     encodedData,
					MimeType: header.Header.Get("Content-Type"),
					Filename: header.Filename,
					Size:     header.Size,
				},
			},
			Options: ai.MultiModalOptions{
				ExtractText:      true,
				AnalyzeSentiment: r.FormValue("analyze_sentiment") == "true",
				ExtractEntities:  r.FormValue("extract_entities") == "true",
				GenerateSummary:  r.FormValue("generate_summary") == "true",
			},
			RequestedAt: time.Now(),
		}

		result, err := engine.ProcessMultiModalRequest(ctx, req)
		if err != nil {
			logger.Error(ctx, "Document analysis failed", err, map[string]interface{}{
				"filename": header.Filename,
			})
			http.Error(w, "Document analysis failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)

		logger.Info(ctx, "Document analysis completed", map[string]interface{}{
			"filename":        header.Filename,
			"processing_time": result.ProcessingTime.Milliseconds(),
		})
	}
}

func handleAudioAnalysis(engine *ai.MultiModalEngine, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Parse multipart form
		err := r.ParseMultipartForm(50 << 20) // 50MB max
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("audio")
		if err != nil {
			http.Error(w, "Audio file required", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Get user ID from context
		userID, err := getUserIDFromContext(ctx)
		if err != nil {
			http.Error(w, "User ID required", http.StatusUnauthorized)
			return
		}

		// Read file data and create request
		data, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}

		encodedData := base64.StdEncoding.EncodeToString(data)

		req := &ai.MultiModalRequest{
			RequestID: uuid.New().String(),
			UserID:    userID,
			Type:      "audio",
			Content: []ai.MultiModalContent{
				{
					ID:       uuid.New().String(),
					Type:     "audio",
					Data:     encodedData,
					MimeType: header.Header.Get("Content-Type"),
					Filename: header.Filename,
					Size:     header.Size,
				},
			},
			Options: ai.MultiModalOptions{
				ProcessAudio:     true,
				AnalyzeSentiment: r.FormValue("analyze_sentiment") == "true",
				ExtractEntities:  r.FormValue("extract_entities") == "true",
			},
			RequestedAt: time.Now(),
		}

		result, err := engine.ProcessMultiModalRequest(ctx, req)
		if err != nil {
			logger.Error(ctx, "Audio analysis failed", err, map[string]interface{}{
				"filename": header.Filename,
			})
			http.Error(w, "Audio analysis failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)

		logger.Info(ctx, "Audio analysis completed", map[string]interface{}{
			"filename":        header.Filename,
			"processing_time": result.ProcessingTime.Milliseconds(),
		})
	}
}

func handleChartAnalysis(engine *ai.MultiModalEngine, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Parse multipart form
		err := r.ParseMultipartForm(10 << 20) // 10MB max
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("chart")
		if err != nil {
			http.Error(w, "Chart image required", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Get user ID from context
		userID, err := getUserIDFromContext(ctx)
		if err != nil {
			http.Error(w, "User ID required", http.StatusUnauthorized)
			return
		}

		// Validate image format
		if !engine.ValidateImageFormat(header.Filename) {
			http.Error(w, "Unsupported image format", http.StatusBadRequest)
			return
		}

		// Parse options with chart-specific settings
		options := ai.MultiModalOptions{
			AnalyzeImages: true,
			AnalyzeCharts: true,
			ExtractText:   true,
			DetectObjects: true,
		}

		result, err := engine.ProcessImageFile(ctx, userID, file, header, options)
		if err != nil {
			logger.Error(ctx, "Chart analysis failed", err, map[string]interface{}{
				"filename": header.Filename,
			})
			http.Error(w, "Chart analysis failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)

		logger.Info(ctx, "Chart analysis completed", map[string]interface{}{
			"filename":        header.Filename,
			"processing_time": result.ProcessingTime.Milliseconds(),
		})
	}
}

func handleGetSupportedFormats(engine *ai.MultiModalEngine, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		formats := engine.GetSupportedFormats()

		response := map[string]interface{}{
			"supported_formats": formats,
			"timestamp":         time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		logger.Info(ctx, "Supported formats retrieved", map[string]interface{}{
			"timestamp": time.Now(),
		})
	}
}

func getUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userIDStr, ok := middleware.GetUserID(ctx)
	if !ok {
		return uuid.Nil, fmt.Errorf("user ID not found in context")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	return userID, nil
}

// User Behavior Learning handlers

func handleLearnFromBehavior(engine *ai.UserBehaviorLearningEngine, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var event ai.BehaviorEvent
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Get user ID from context
		userID, err := getUserIDFromContext(ctx)
		if err != nil {
			http.Error(w, "User ID required", http.StatusUnauthorized)
			return
		}
		event.UserID = userID

		// Set event ID if not provided
		if event.ID == "" {
			event.ID = uuid.New().String()
		}

		// Set timestamp if not provided
		if event.Timestamp.IsZero() {
			event.Timestamp = time.Now()
		}

		err = engine.LearnFromBehavior(ctx, &event)
		if err != nil {
			logger.Error(ctx, "Failed to learn from behavior", err, map[string]interface{}{
				"user_id":    userID,
				"event_type": event.Type,
			})
			http.Error(w, "Learning failed", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"success":   true,
			"event_id":  event.ID,
			"timestamp": time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		logger.Info(ctx, "Behavior learning completed", map[string]interface{}{
			"user_id":    userID,
			"event_id":   event.ID,
			"event_type": event.Type,
		})
	}
}

func handleGetUserBehaviorProfile(engine *ai.UserBehaviorLearningEngine, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get user ID from context
		userID, err := getUserIDFromContext(ctx)
		if err != nil {
			http.Error(w, "User ID required", http.StatusUnauthorized)
			return
		}

		profile, err := engine.GetUserProfile(ctx, userID)
		if err != nil {
			logger.Error(ctx, "Failed to get user profile", err, map[string]interface{}{
				"user_id": userID,
			})
			http.Error(w, "Profile not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(profile)

		logger.Info(ctx, "User profile retrieved", map[string]interface{}{
			"user_id":           userID,
			"observation_count": profile.ObservationCount,
			"confidence":        profile.Confidence,
		})
	}
}

func handleGetRecommendations(engine *ai.UserBehaviorLearningEngine, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get user ID from context
		userID, err := getUserIDFromContext(ctx)
		if err != nil {
			http.Error(w, "User ID required", http.StatusUnauthorized)
			return
		}

		// Parse limit parameter
		limitStr := r.URL.Query().Get("limit")
		limit := 10 // default
		if limitStr != "" {
			if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
				limit = parsedLimit
			}
		}

		recommendations, err := engine.GetPersonalizedRecommendations(ctx, userID, limit)
		if err != nil {
			logger.Error(ctx, "Failed to get recommendations", err, map[string]interface{}{
				"user_id": userID,
			})
			http.Error(w, "Failed to get recommendations", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"recommendations": recommendations,
			"count":           len(recommendations),
			"user_id":         userID,
			"timestamp":       time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		logger.Info(ctx, "Recommendations retrieved", map[string]interface{}{
			"user_id": userID,
			"count":   len(recommendations),
		})
	}
}

func handleGetBehaviorHistory(engine *ai.UserBehaviorLearningEngine, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get user ID from context
		userID, err := getUserIDFromContext(ctx)
		if err != nil {
			http.Error(w, "User ID required", http.StatusUnauthorized)
			return
		}

		// Parse limit parameter
		limitStr := r.URL.Query().Get("limit")
		limit := 50 // default
		if limitStr != "" {
			if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
				limit = parsedLimit
			}
		}

		history, err := engine.GetBehaviorHistory(ctx, userID, limit)
		if err != nil {
			logger.Error(ctx, "Failed to get behavior history", err, map[string]interface{}{
				"user_id": userID,
			})
			http.Error(w, "Failed to get behavior history", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"history":   history,
			"count":     len(history),
			"user_id":   userID,
			"timestamp": time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		logger.Info(ctx, "Behavior history retrieved", map[string]interface{}{
			"user_id": userID,
			"count":   len(history),
		})
	}
}

func handleUpdateRecommendationStatus(engine *ai.UserBehaviorLearningEngine, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get user ID from context
		userID, err := getUserIDFromContext(ctx)
		if err != nil {
			http.Error(w, "User ID required", http.StatusUnauthorized)
			return
		}

		// Get recommendation ID from path
		recommendationID := r.PathValue("id")
		if recommendationID == "" {
			http.Error(w, "Recommendation ID required", http.StatusBadRequest)
			return
		}

		// Parse request body
		var req struct {
			Status string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate status
		validStatuses := map[string]bool{
			"pending":  true,
			"accepted": true,
			"rejected": true,
			"expired":  true,
		}
		if !validStatuses[req.Status] {
			http.Error(w, "Invalid status", http.StatusBadRequest)
			return
		}

		err = engine.UpdateRecommendationStatus(ctx, userID, recommendationID, req.Status)
		if err != nil {
			logger.Error(ctx, "Failed to update recommendation status", err, map[string]interface{}{
				"user_id":           userID,
				"recommendation_id": recommendationID,
				"status":            req.Status,
			})
			http.Error(w, "Failed to update recommendation status", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"success":           true,
			"recommendation_id": recommendationID,
			"status":            req.Status,
			"timestamp":         time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		logger.Info(ctx, "Recommendation status updated", map[string]interface{}{
			"user_id":           userID,
			"recommendation_id": recommendationID,
			"status":            req.Status,
		})
	}
}

func handleGetLearningModels(engine *ai.UserBehaviorLearningEngine, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		models := engine.GetLearningModels()

		response := map[string]interface{}{
			"models":    models,
			"count":     len(models),
			"timestamp": time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		logger.Info(ctx, "Learning models retrieved", map[string]interface{}{
			"count": len(models),
		})
	}
}

// Market Pattern Adaptation handlers

func handleDetectMarketPatterns(engine *ai.MarketAdaptationEngine, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var marketData map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&marketData); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		patterns, err := engine.DetectPatterns(ctx, marketData)
		if err != nil {
			logger.Error(ctx, "Failed to detect market patterns", err, map[string]interface{}{
				"data_points": len(marketData),
			})
			http.Error(w, "Pattern detection failed", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"patterns":  patterns,
			"count":     len(patterns),
			"timestamp": time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		logger.Info(ctx, "Market patterns detected", map[string]interface{}{
			"patterns_count": len(patterns),
		})
	}
}

func handleGetMarketPatternsAdaptation(engine *ai.MarketAdaptationEngine, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Parse query parameters for filters
		filters := make(map[string]interface{})
		if asset := r.URL.Query().Get("asset"); asset != "" {
			filters["asset"] = asset
		}
		if patternType := r.URL.Query().Get("type"); patternType != "" {
			filters["type"] = patternType
		}
		if minConfidenceStr := r.URL.Query().Get("min_confidence"); minConfidenceStr != "" {
			if minConfidence, err := strconv.ParseFloat(minConfidenceStr, 64); err == nil {
				filters["min_confidence"] = minConfidence
			}
		}

		patterns, err := engine.GetDetectedPatterns(ctx, filters)
		if err != nil {
			logger.Error(ctx, "Failed to get market patterns", err, map[string]interface{}{
				"filters": filters,
			})
			http.Error(w, "Failed to get market patterns", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"patterns":  patterns,
			"count":     len(patterns),
			"filters":   filters,
			"timestamp": time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		logger.Info(ctx, "Market patterns retrieved", map[string]interface{}{
			"count":   len(patterns),
			"filters": filters,
		})
	}
}

func handleAdaptStrategies(engine *ai.MarketAdaptationEngine, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var request struct {
			Patterns []*ai.DetectedPattern `json:"patterns"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err := engine.AdaptStrategies(ctx, request.Patterns)
		if err != nil {
			logger.Error(ctx, "Failed to adapt strategies", err, map[string]interface{}{
				"patterns_count": len(request.Patterns),
			})
			http.Error(w, "Strategy adaptation failed", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"success":        true,
			"patterns_count": len(request.Patterns),
			"timestamp":      time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		logger.Info(ctx, "Strategies adapted successfully", map[string]interface{}{
			"patterns_count": len(request.Patterns),
		})
	}
}

func handleGetAdaptiveStrategies(engine *ai.MarketAdaptationEngine, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		strategies, err := engine.GetAdaptiveStrategies(ctx)
		if err != nil {
			logger.Error(ctx, "Failed to get adaptive strategies", err, nil)
			http.Error(w, "Failed to get adaptive strategies", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"strategies": strategies,
			"count":      len(strategies),
			"timestamp":  time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		logger.Info(ctx, "Adaptive strategies retrieved", map[string]interface{}{
			"count": len(strategies),
		})
	}
}

func handleAddAdaptiveStrategy(engine *ai.MarketAdaptationEngine, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var strategy ai.AdaptiveStrategy
		if err := json.NewDecoder(r.Body).Decode(&strategy); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err := engine.AddAdaptiveStrategy(ctx, &strategy)
		if err != nil {
			logger.Error(ctx, "Failed to add adaptive strategy", err, map[string]interface{}{
				"strategy_name": strategy.Name,
				"strategy_type": strategy.Type,
			})
			http.Error(w, "Failed to add adaptive strategy", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"success":     true,
			"strategy_id": strategy.ID,
			"timestamp":   time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		logger.Info(ctx, "Adaptive strategy added", map[string]interface{}{
			"strategy_id":   strategy.ID,
			"strategy_name": strategy.Name,
			"strategy_type": strategy.Type,
		})
	}
}

func handleUpdateStrategyStatus(engine *ai.MarketAdaptationEngine, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get strategy ID from path
		strategyID := r.PathValue("id")
		if strategyID == "" {
			http.Error(w, "Strategy ID required", http.StatusBadRequest)
			return
		}

		// Parse request body
		var req struct {
			IsActive bool `json:"is_active"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err := engine.UpdateStrategyStatus(ctx, strategyID, req.IsActive)
		if err != nil {
			logger.Error(ctx, "Failed to update strategy status", err, map[string]interface{}{
				"strategy_id": strategyID,
				"is_active":   req.IsActive,
			})
			http.Error(w, "Failed to update strategy status", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"success":     true,
			"strategy_id": strategyID,
			"is_active":   req.IsActive,
			"timestamp":   time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		logger.Info(ctx, "Strategy status updated", map[string]interface{}{
			"strategy_id": strategyID,
			"is_active":   req.IsActive,
		})
	}
}

func handleGetMarketAdaptationHistory(engine *ai.MarketAdaptationEngine, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Parse limit parameter
		limitStr := r.URL.Query().Get("limit")
		limit := 50 // default
		if limitStr != "" {
			if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
				limit = parsedLimit
			}
		}

		history, err := engine.GetAdaptationHistory(ctx, limit)
		if err != nil {
			logger.Error(ctx, "Failed to get adaptation history", err, map[string]interface{}{
				"limit": limit,
			})
			http.Error(w, "Failed to get adaptation history", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"history":   history,
			"count":     len(history),
			"limit":     limit,
			"timestamp": time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		logger.Info(ctx, "Adaptation history retrieved", map[string]interface{}{
			"count": len(history),
			"limit": limit,
		})
	}
}

func handleGetStrategyPerformanceMetrics(engine *ai.MarketAdaptationEngine, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get strategy ID from path
		strategyID := r.PathValue("strategy_id")
		if strategyID == "" {
			http.Error(w, "Strategy ID required", http.StatusBadRequest)
			return
		}

		metrics, err := engine.GetPerformanceMetrics(ctx, strategyID)
		if err != nil {
			logger.Error(ctx, "Failed to get performance metrics", err, map[string]interface{}{
				"strategy_id": strategyID,
			})
			http.Error(w, "Failed to get performance metrics", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(metrics)

		logger.Info(ctx, "Performance metrics retrieved", map[string]interface{}{
			"strategy_id": strategyID,
		})
	}
}

// Crypto Coin Analyzer handlers

func handleCryptoCoinAnalysis(analyzer *ai.CryptoCoinAnalyzer, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get symbol from path
		symbol := r.PathValue("symbol")
		if symbol == "" {
			http.Error(w, "Symbol is required", http.StatusBadRequest)
			return
		}

		// Validate symbol format (basic validation)
		symbol = strings.ToUpper(strings.TrimSpace(symbol))
		if len(symbol) < 2 || len(symbol) > 10 {
			http.Error(w, "Invalid symbol format", http.StatusBadRequest)
			return
		}

		logger.Info(ctx, "Starting crypto coin analysis", map[string]interface{}{
			"symbol": symbol,
			"method": r.Method,
		})

		// Perform analysis
		report, err := analyzer.AnalyzeCoin(ctx, symbol)
		if err != nil {
			logger.Error(ctx, "Crypto coin analysis failed", err, map[string]interface{}{
				"symbol": symbol,
			})
			http.Error(w, fmt.Sprintf("Analysis failed: %v", err), http.StatusInternalServerError)
			return
		}

		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(report); err != nil {
			logger.Error(ctx, "Failed to encode response", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}

		logger.Info(ctx, "Crypto coin analysis completed", map[string]interface{}{
			"symbol":        symbol,
			"sources_count": len(report.Sources),
			"news_count":    len(report.NewsAndEvents),
		})
	}
}

func handleCryptoCoinReport(analyzer *ai.CryptoCoinAnalyzer, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get symbol from path
		symbol := r.PathValue("symbol")
		if symbol == "" {
			http.Error(w, "Symbol is required", http.StatusBadRequest)
			return
		}

		// Validate symbol format
		symbol = strings.ToUpper(strings.TrimSpace(symbol))
		if len(symbol) < 2 || len(symbol) > 10 {
			http.Error(w, "Invalid symbol format", http.StatusBadRequest)
			return
		}

		logger.Info(ctx, "Generating crypto coin report", map[string]interface{}{
			"symbol": symbol,
			"method": r.Method,
		})

		// Generate structured report
		reportMarkdown, err := analyzer.AnalyzeCoinWithStructuredReport(ctx, symbol)
		if err != nil {
			logger.Error(ctx, "Crypto coin report generation failed", err, map[string]interface{}{
				"symbol": symbol,
			})
			http.Error(w, fmt.Sprintf("Report generation failed: %v", err), http.StatusInternalServerError)
			return
		}

		// Check if client wants JSON or markdown
		acceptHeader := r.Header.Get("Accept")
		if strings.Contains(acceptHeader, "application/json") {
			// Return JSON with markdown content
			response := map[string]interface{}{
				"symbol":    symbol,
				"report":    reportMarkdown,
				"format":    "markdown",
				"timestamp": time.Now(),
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else {
			// Return raw markdown
			w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
			w.Write([]byte(reportMarkdown))
		}

		logger.Info(ctx, "Crypto coin report generated", map[string]interface{}{
			"symbol":      symbol,
			"report_size": len(reportMarkdown),
		})
	}
}
