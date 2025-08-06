package ai

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/ml"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

// EnsembleModelManager manages multiple models for improved accuracy
type EnsembleModelManager struct {
	logger             *observability.Logger
	config             *EnsembleConfig
	baseModels         map[string]ml.Model
	metaLearner        *MetaLearner
	votingStrategy     VotingStrategy
	performanceTracker *EnsemblePerformanceTracker
	modelWeights       map[string]float64
	predictionCache    map[string]*EnsemblePrediction
	mu                 sync.RWMutex
	lastUpdate         time.Time
}

// EnsembleConfig contains ensemble model configuration
type EnsembleConfig struct {
	MinModels           int
	MaxModels           int
	VotingStrategy      string
	WeightUpdateRate    float64
	PerformanceWindow   time.Duration
	RetrainingThreshold float64
	EnableMetaLearning  bool
	CacheSize           int
	CacheTTL            time.Duration
}

// VotingStrategy defines how ensemble predictions are combined
type VotingStrategy interface {
	CombinePredictions(predictions []ModelPrediction, weights map[string]float64) (*EnsemblePrediction, error)
}

// MetaLearner learns which models perform best in different conditions
type MetaLearner struct {
	logger           *observability.Logger
	modelPerformance map[string]*ModelPerformanceHistory
	contextFeatures  []string
	learningRate     float64
	mu               sync.RWMutex
}

// EnsemblePerformanceTracker tracks ensemble model performance
type EnsemblePerformanceTracker struct {
	predictions    []EnsemblePredictionRecord
	accuracy       float64
	precision      float64
	recall         float64
	f1Score        float64
	diversityScore float64
	consensusScore float64
	mu             sync.RWMutex
}

// ModelPrediction represents a prediction from a single model
type ModelPrediction struct {
	ModelID    string
	Prediction interface{}
	Confidence float64
	Features   map[string]interface{}
	Timestamp  time.Time
	Metadata   map[string]interface{}
}

// EnsemblePrediction represents the combined prediction from multiple models
type EnsemblePrediction struct {
	PredictionID    string
	FinalPrediction interface{}
	Confidence      float64
	ModelVotes      map[string]ModelPrediction
	Weights         map[string]float64
	Consensus       float64
	Diversity       float64
	Timestamp       time.Time
	Metadata        map[string]interface{}
}

// EnsemblePredictionRecord tracks ensemble prediction outcomes
type EnsemblePredictionRecord struct {
	PredictionID  string
	Prediction    *EnsemblePrediction
	ActualOutcome interface{}
	Correct       bool
	Error         float64
	Timestamp     time.Time
}

// ModelPerformanceHistory tracks individual model performance over time
type ModelPerformanceHistory struct {
	ModelID         string
	AccuracyHistory []float64
	LatencyHistory  []time.Duration
	ContextScores   map[string]float64
	LastUpdated     time.Time
}

// NewEnsembleModelManager creates a new ensemble model manager
func NewEnsembleModelManager(logger *observability.Logger) *EnsembleModelManager {
	config := &EnsembleConfig{
		MinModels:           3,
		MaxModels:           10,
		VotingStrategy:      "weighted_average",
		WeightUpdateRate:    0.1,
		PerformanceWindow:   24 * time.Hour,
		RetrainingThreshold: 0.05,
		EnableMetaLearning:  true,
		CacheSize:           1000,
		CacheTTL:            5 * time.Minute,
	}

	metaLearner := &MetaLearner{
		logger:           logger,
		modelPerformance: make(map[string]*ModelPerformanceHistory),
		contextFeatures:  []string{"market_volatility", "trading_volume", "time_of_day", "market_trend"},
		learningRate:     0.01,
	}

	return &EnsembleModelManager{
		logger:             logger,
		config:             config,
		baseModels:         make(map[string]ml.Model),
		metaLearner:        metaLearner,
		votingStrategy:     NewWeightedAverageVoting(),
		performanceTracker: &EnsemblePerformanceTracker{},
		modelWeights:       make(map[string]float64),
		predictionCache:    make(map[string]*EnsemblePrediction),
	}
}

// AddModel adds a base model to the ensemble
func (e *EnsembleModelManager) AddModel(modelID string, model ml.Model) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if len(e.baseModels) >= e.config.MaxModels {
		return fmt.Errorf("maximum number of models (%d) reached", e.config.MaxModels)
	}

	e.baseModels[modelID] = model
	e.modelWeights[modelID] = 1.0 / float64(len(e.baseModels)+1) // Equal initial weights

	// Initialize performance history
	e.metaLearner.mu.Lock()
	e.metaLearner.modelPerformance[modelID] = &ModelPerformanceHistory{
		ModelID:         modelID,
		AccuracyHistory: []float64{},
		LatencyHistory:  []time.Duration{},
		ContextScores:   make(map[string]float64),
		LastUpdated:     time.Now(),
	}
	e.metaLearner.mu.Unlock()

	e.logger.Info(context.Background(), "Model added to ensemble", map[string]interface{}{
		"model_id":       modelID,
		"total_models":   len(e.baseModels),
		"initial_weight": e.modelWeights[modelID],
	})

	return nil
}

// Predict generates ensemble predictions using all base models
func (e *EnsembleModelManager) Predict(ctx context.Context, features map[string]interface{}) (*EnsemblePrediction, error) {
	// Check cache first
	cacheKey := e.generateCacheKey(features)
	if cached := e.getCachedPrediction(cacheKey); cached != nil {
		return cached, nil
	}

	e.mu.RLock()
	models := make(map[string]ml.Model)
	for id, model := range e.baseModels {
		models[id] = model
	}
	weights := make(map[string]float64)
	for id, weight := range e.modelWeights {
		weights[id] = weight
	}
	e.mu.RUnlock()

	if len(models) < e.config.MinModels {
		return nil, fmt.Errorf("insufficient models for ensemble prediction: %d < %d", len(models), e.config.MinModels)
	}

	// Collect predictions from all models
	predictions := make([]ModelPrediction, 0, len(models))
	var wg sync.WaitGroup
	predictionChan := make(chan ModelPrediction, len(models))
	errorChan := make(chan error, len(models))

	for modelID, model := range models {
		wg.Add(1)
		go func(id string, m ml.Model) {
			defer wg.Done()

			start := time.Now()
			result, err := m.Predict(ctx, features)
			duration := time.Since(start)

			if err != nil {
				e.logger.Error(ctx, "Model prediction failed", err, map[string]interface{}{
					"model_id": id,
				})
				errorChan <- err
				return
			}

			prediction := ModelPrediction{
				ModelID:    id,
				Prediction: result.Value,
				Confidence: result.Confidence,
				Features:   features,
				Timestamp:  time.Now(),
				Metadata: map[string]interface{}{
					"latency":       duration,
					"model_version": m.GetInfo().Version,
				},
			}

			predictionChan <- prediction
		}(modelID, model)
	}

	// Wait for all predictions to complete
	go func() {
		wg.Wait()
		close(predictionChan)
		close(errorChan)
	}()

	// Collect results
	for prediction := range predictionChan {
		predictions = append(predictions, prediction)
	}

	// Check for errors
	var lastError error
	for err := range errorChan {
		lastError = err
	}

	if len(predictions) < e.config.MinModels {
		return nil, fmt.Errorf("insufficient successful predictions: %d < %d, last error: %v",
			len(predictions), e.config.MinModels, lastError)
	}

	// Combine predictions using voting strategy
	ensemblePrediction, err := e.votingStrategy.CombinePredictions(predictions, weights)
	if err != nil {
		return nil, fmt.Errorf("failed to combine predictions: %w", err)
	}

	// Add ensemble metadata
	ensemblePrediction.PredictionID = uuid.New().String()
	ensemblePrediction.Timestamp = time.Now()
	ensemblePrediction.Weights = weights
	ensemblePrediction.Consensus = e.calculateConsensus(predictions)
	ensemblePrediction.Diversity = e.calculateDiversity(predictions)

	// Cache the prediction
	e.cachePrediction(cacheKey, ensemblePrediction)

	// Update meta-learner with context
	if e.config.EnableMetaLearning {
		e.metaLearner.updateContext(ctx, predictions, features)
	}

	e.logger.Debug(ctx, "Ensemble prediction generated", map[string]interface{}{
		"prediction_id": ensemblePrediction.PredictionID,
		"models_used":   len(predictions),
		"confidence":    ensemblePrediction.Confidence,
		"consensus":     ensemblePrediction.Consensus,
		"diversity":     ensemblePrediction.Diversity,
	})

	return ensemblePrediction, nil
}

// UpdateModelWeights updates model weights based on performance
func (e *EnsembleModelManager) UpdateModelWeights(ctx context.Context, feedback *EnsembleFeedback) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Update individual model performance
	for modelID, _ := range feedback.ModelPredictions {
		if history, exists := e.metaLearner.modelPerformance[modelID]; exists {
			accuracy := 1.0
			if !feedback.Correct {
				accuracy = 0.0
			}

			history.AccuracyHistory = append(history.AccuracyHistory, accuracy)

			// Keep only recent history
			if len(history.AccuracyHistory) > 100 {
				history.AccuracyHistory = history.AccuracyHistory[1:]
			}

			history.LastUpdated = time.Now()
		}
	}

	// Recalculate weights based on recent performance
	e.recalculateWeights()

	// Update ensemble performance tracker
	e.performanceTracker.mu.Lock()
	record := EnsemblePredictionRecord{
		PredictionID:  feedback.PredictionID,
		Prediction:    feedback.Prediction,
		ActualOutcome: feedback.ActualOutcome,
		Correct:       feedback.Correct,
		Error:         feedback.Error,
		Timestamp:     time.Now(),
	}
	e.performanceTracker.predictions = append(e.performanceTracker.predictions, record)
	e.performanceTracker.mu.Unlock()

	e.logger.Info(ctx, "Model weights updated", map[string]interface{}{
		"prediction_id": feedback.PredictionID,
		"correct":       feedback.Correct,
		"error":         feedback.Error,
	})

	return nil
}

// EnsembleFeedback contains feedback for ensemble predictions
type EnsembleFeedback struct {
	PredictionID     string
	Prediction       *EnsemblePrediction
	ActualOutcome    interface{}
	Correct          bool
	Error            float64
	ModelPredictions map[string]ModelPrediction
}

// recalculateWeights recalculates model weights based on performance
func (e *EnsembleModelManager) recalculateWeights() {
	totalWeight := 0.0
	newWeights := make(map[string]float64)

	for modelID := range e.baseModels {
		if history, exists := e.metaLearner.modelPerformance[modelID]; exists {
			// Calculate recent accuracy
			recentAccuracy := e.calculateRecentAccuracy(history.AccuracyHistory)

			// Weight based on accuracy with minimum threshold
			weight := math.Max(0.1, recentAccuracy)
			newWeights[modelID] = weight
			totalWeight += weight
		}
	}

	// Normalize weights
	for modelID := range newWeights {
		e.modelWeights[modelID] = newWeights[modelID] / totalWeight
	}
}

// calculateRecentAccuracy calculates recent accuracy for a model
func (e *EnsembleModelManager) calculateRecentAccuracy(accuracyHistory []float64) float64 {
	if len(accuracyHistory) == 0 {
		return 0.5 // Default for new models
	}

	// Use exponential moving average for recent accuracy
	alpha := 0.1
	ema := accuracyHistory[0]
	for i := 1; i < len(accuracyHistory); i++ {
		ema = alpha*accuracyHistory[i] + (1-alpha)*ema
	}

	return ema
}

// calculateConsensus calculates how much models agree
func (e *EnsembleModelManager) calculateConsensus(predictions []ModelPrediction) float64 {
	if len(predictions) < 2 {
		return 1.0
	}

	// For numeric predictions, calculate variance
	values := make([]float64, 0, len(predictions))
	for _, pred := range predictions {
		if val, ok := pred.Prediction.(float64); ok {
			values = append(values, val)
		}
	}

	if len(values) < 2 {
		return 1.0
	}

	// Calculate coefficient of variation (lower = higher consensus)
	mean := 0.0
	for _, val := range values {
		mean += val
	}
	mean /= float64(len(values))

	variance := 0.0
	for _, val := range values {
		variance += math.Pow(val-mean, 2)
	}
	variance /= float64(len(values))

	if mean == 0 {
		return 1.0
	}

	cv := math.Sqrt(variance) / math.Abs(mean)
	return math.Max(0.0, 1.0-cv) // Higher consensus = lower coefficient of variation
}

// calculateDiversity calculates prediction diversity
func (e *EnsembleModelManager) calculateDiversity(predictions []ModelPrediction) float64 {
	// Diversity is the opposite of consensus
	consensus := e.calculateConsensus(predictions)
	return 1.0 - consensus
}

// generateCacheKey generates a cache key for features
func (e *EnsembleModelManager) generateCacheKey(features map[string]interface{}) string {
	// Simple hash-based cache key generation
	// In production, use a proper hash function
	key := "ensemble:"
	for k, v := range features {
		key += fmt.Sprintf("%s:%v:", k, v)
	}
	return key
}

// getCachedPrediction retrieves a cached prediction
func (e *EnsembleModelManager) getCachedPrediction(key string) *EnsemblePrediction {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if prediction, exists := e.predictionCache[key]; exists {
		if time.Since(prediction.Timestamp) < e.config.CacheTTL {
			return prediction
		}
		// Remove expired cache entry
		delete(e.predictionCache, key)
	}

	return nil
}

// cachePrediction caches a prediction
func (e *EnsembleModelManager) cachePrediction(key string, prediction *EnsemblePrediction) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Enforce cache size limit
	if len(e.predictionCache) >= e.config.CacheSize {
		// Remove oldest entry (simple FIFO)
		for k := range e.predictionCache {
			delete(e.predictionCache, k)
			break
		}
	}

	e.predictionCache[key] = prediction
}

// GetPerformanceMetrics returns ensemble performance metrics
func (e *EnsembleModelManager) GetPerformanceMetrics() map[string]interface{} {
	e.performanceTracker.mu.RLock()
	defer e.performanceTracker.mu.RUnlock()

	totalPredictions := len(e.performanceTracker.predictions)
	if totalPredictions == 0 {
		return map[string]interface{}{
			"total_predictions": 0,
			"accuracy":          0.0,
		}
	}

	correct := 0
	totalError := 0.0
	for _, record := range e.performanceTracker.predictions {
		if record.Correct {
			correct++
		}
		totalError += record.Error
	}

	accuracy := float64(correct) / float64(totalPredictions)
	avgError := totalError / float64(totalPredictions)

	return map[string]interface{}{
		"total_predictions": totalPredictions,
		"accuracy":          accuracy,
		"average_error":     avgError,
		"model_weights":     e.modelWeights,
		"cache_size":        len(e.predictionCache),
	}
}

// updateContext updates meta-learner context
func (m *MetaLearner) updateContext(ctx context.Context, predictions []ModelPrediction, features map[string]interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Extract context features and update model performance scores
	for _, pred := range predictions {
		if history, exists := m.modelPerformance[pred.ModelID]; exists {
			// Update context-specific scores based on current market conditions
			for _, feature := range m.contextFeatures {
				if value, ok := features[feature]; ok {
					if score, exists := history.ContextScores[feature]; exists {
						// Update with exponential moving average
						if numVal, ok := value.(float64); ok {
							history.ContextScores[feature] = score*(1-m.learningRate) + numVal*m.learningRate
						}
					} else {
						if numVal, ok := value.(float64); ok {
							history.ContextScores[feature] = numVal
						}
					}
				}
			}
		}
	}
}
