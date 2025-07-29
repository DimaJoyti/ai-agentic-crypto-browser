package ai

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/ai-agentic-browser/pkg/ml"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

// AdaptiveModelManager manages adaptive models that learn and evolve
type AdaptiveModelManager struct {
	models          map[string]*AdaptiveModel
	learningEngine  *LearningEngine
	logger          *observability.Logger
	config          *AdaptiveModelConfig
	adaptationQueue []AdaptationRequest
}

// AdaptiveModelConfig holds configuration for adaptive models
type AdaptiveModelConfig struct {
	AdaptationInterval    time.Duration `json:"adaptation_interval"`
	PerformanceThreshold  float64       `json:"performance_threshold"`
	MinAdaptationSamples  int           `json:"min_adaptation_samples"`
	MaxAdaptationRate     float64       `json:"max_adaptation_rate"`
	AdaptationDecay       float64       `json:"adaptation_decay"`
	EnableAutoAdaptation  bool          `json:"enable_auto_adaptation"`
	AdaptationStrategies  []string      `json:"adaptation_strategies"`
}

// AdaptationRequest represents a request for model adaptation
type AdaptationRequest struct {
	ModelID     string                 `json:"model_id"`
	Type        string                 `json:"type"` // performance, drift, feedback
	Trigger     string                 `json:"trigger"`
	Data        map[string]interface{} `json:"data"`
	Priority    int                    `json:"priority"`
	RequestedAt time.Time              `json:"requested_at"`
	UserID      uuid.UUID              `json:"user_id,omitempty"`
}

// AdaptationStrategy defines how models should adapt
type AdaptationStrategy interface {
	CanAdapt(model *AdaptiveModel, request *AdaptationRequest) bool
	Adapt(ctx context.Context, model *AdaptiveModel, request *AdaptationRequest) (*AdaptationResult, error)
	GetType() string
}

// AdaptationResult represents the result of an adaptation
type AdaptationResult struct {
	Success         bool                   `json:"success"`
	Type            string                 `json:"type"`
	Description     string                 `json:"description"`
	PerformanceGain float64                `json:"performance_gain"`
	Confidence      float64                `json:"confidence"`
	Changes         []ModelChange          `json:"changes"`
	Metadata        map[string]interface{} `json:"metadata"`
	Timestamp       time.Time              `json:"timestamp"`
}

// ModelChange represents a change made to a model
type ModelChange struct {
	Component   string      `json:"component"` // weights, hyperparameters, architecture
	Parameter   string      `json:"parameter"`
	OldValue    interface{} `json:"old_value"`
	NewValue    interface{} `json:"new_value"`
	Impact      float64     `json:"impact"`
	Reversible  bool        `json:"reversible"`
}

// PerformanceBasedAdaptation adapts models based on performance degradation
type PerformanceBasedAdaptation struct {
	logger *observability.Logger
}

// FeedbackBasedAdaptation adapts models based on user feedback
type FeedbackBasedAdaptation struct {
	logger *observability.Logger
}

// DriftBasedAdaptation adapts models based on concept drift detection
type DriftBasedAdaptation struct {
	logger         *observability.Logger
	driftDetector  *ConceptDriftDetector
}

// ConceptDriftDetector detects concept drift in data
type ConceptDriftDetector struct {
	windowSize      int
	driftThreshold  float64
	recentData      []DataPoint
	referenceData   []DataPoint
	lastCheck       time.Time
	driftHistory    []DriftEvent
}

// DataPoint represents a data point for drift detection
type DataPoint struct {
	Features  map[string]float64 `json:"features"`
	Label     interface{}        `json:"label"`
	Timestamp time.Time          `json:"timestamp"`
	Weight    float64            `json:"weight"`
}

// DriftEvent represents a detected drift event
type DriftEvent struct {
	Type        string    `json:"type"` // gradual, sudden, recurring
	Severity    float64   `json:"severity"`
	Confidence  float64   `json:"confidence"`
	DetectedAt  time.Time `json:"detected_at"`
	Description string    `json:"description"`
	Features    []string  `json:"features"` // which features drifted
}

// OnlineLearningAdapter implements online learning for models
type OnlineLearningAdapter struct {
	logger       *observability.Logger
	learningRate float64
	batchSize    int
	updateFreq   time.Duration
}

// NewAdaptiveModelManager creates a new adaptive model manager
func NewAdaptiveModelManager(learningEngine *LearningEngine, logger *observability.Logger) *AdaptiveModelManager {
	config := &AdaptiveModelConfig{
		AdaptationInterval:    1 * time.Hour,
		PerformanceThreshold:  0.05, // 5% performance drop triggers adaptation
		MinAdaptationSamples:  100,
		MaxAdaptationRate:     0.1,
		AdaptationDecay:       0.95,
		EnableAutoAdaptation:  true,
		AdaptationStrategies:  []string{"performance", "feedback", "drift"},
	}

	manager := &AdaptiveModelManager{
		models:          make(map[string]*AdaptiveModel),
		learningEngine:  learningEngine,
		logger:          logger,
		config:          config,
		adaptationQueue: []AdaptationRequest{},
	}

	// Start adaptation loop
	go manager.startAdaptationLoop()

	return manager
}

// RegisterAdaptiveModel registers a model for adaptive learning
func (m *AdaptiveModelManager) RegisterAdaptiveModel(modelID string, baseModel ml.Model) error {
	adaptiveModel := &AdaptiveModel{
		ID:             modelID,
		Name:           baseModel.GetInfo().Name,
		Type:           string(baseModel.GetInfo().Type),
		BaseModel:      baseModel,
		Adaptations:    []ModelAdaptation{},
		Performance:    &ModelPerformance{ModelID: modelID, ModelType: string(baseModel.GetInfo().Type)},
		LearningRate:   0.01,
		AdaptationRate: 0.05,
		LastAdaptation: time.Now(),
		IsAdapting:     false,
		Metadata:       make(map[string]interface{}),
	}

	m.models[modelID] = adaptiveModel

	m.logger.Info(context.Background(), "Adaptive model registered", map[string]interface{}{
		"model_id":   modelID,
		"model_type": adaptiveModel.Type,
	})

	return nil
}

// RequestAdaptation requests adaptation for a model
func (m *AdaptiveModelManager) RequestAdaptation(request *AdaptationRequest) error {
	request.RequestedAt = time.Now()
	m.adaptationQueue = append(m.adaptationQueue, *request)

	m.logger.Info(context.Background(), "Adaptation requested", map[string]interface{}{
		"model_id": request.ModelID,
		"type":     request.Type,
		"trigger":  request.Trigger,
	})

	return nil
}

// startAdaptationLoop starts the background adaptation loop
func (m *AdaptiveModelManager) startAdaptationLoop() {
	ticker := time.NewTicker(m.config.AdaptationInterval)
	defer ticker.Stop()

	for range ticker.C {
		ctx := context.Background()
		m.processAdaptationQueue(ctx)
		
		if m.config.EnableAutoAdaptation {
			m.checkForAutoAdaptation(ctx)
		}
	}
}

// processAdaptationQueue processes pending adaptation requests
func (m *AdaptiveModelManager) processAdaptationQueue(ctx context.Context) {
	if len(m.adaptationQueue) == 0 {
		return
	}

	// Sort by priority and timestamp
	m.sortAdaptationQueue()

	processed := 0
	for i, request := range m.adaptationQueue {
		if processed >= 5 { // Limit processing per cycle
			break
		}

		if err := m.processAdaptationRequest(ctx, &request); err != nil {
			m.logger.Error(ctx, "Failed to process adaptation request", err, map[string]interface{}{
				"model_id": request.ModelID,
				"type":     request.Type,
			})
		} else {
			processed++
		}

		// Remove processed request
		m.adaptationQueue = append(m.adaptationQueue[:i], m.adaptationQueue[i+1:]...)
	}

	m.logger.Info(ctx, "Processed adaptation requests", map[string]interface{}{
		"processed": processed,
		"remaining": len(m.adaptationQueue),
	})
}

// processAdaptationRequest processes a single adaptation request
func (m *AdaptiveModelManager) processAdaptationRequest(ctx context.Context, request *AdaptationRequest) error {
	model, exists := m.models[request.ModelID]
	if !exists {
		return fmt.Errorf("model %s not found", request.ModelID)
	}

	if model.IsAdapting {
		return fmt.Errorf("model %s is already adapting", request.ModelID)
	}

	model.IsAdapting = true
	defer func() { model.IsAdapting = false }()

	// Select adaptation strategy
	strategy := m.selectAdaptationStrategy(request)
	if strategy == nil {
		return fmt.Errorf("no suitable adaptation strategy for request type %s", request.Type)
	}

	// Check if adaptation is needed
	if !strategy.CanAdapt(model, request) {
		m.logger.Info(ctx, "Adaptation not needed", map[string]interface{}{
			"model_id": request.ModelID,
			"type":     request.Type,
		})
		return nil
	}

	// Perform adaptation
	result, err := strategy.Adapt(ctx, model, request)
	if err != nil {
		return fmt.Errorf("adaptation failed: %w", err)
	}

	// Record adaptation
	adaptation := ModelAdaptation{
		ID:          uuid.New().String(),
		Timestamp:   time.Now(),
		Type:        result.Type,
		Description: result.Description,
		Impact:      result.PerformanceGain,
		Success:     result.Success,
		Metadata:    result.Metadata,
	}

	model.Adaptations = append(model.Adaptations, adaptation)
	model.LastAdaptation = time.Now()

	m.logger.Info(ctx, "Model adaptation completed", map[string]interface{}{
		"model_id":         request.ModelID,
		"adaptation_type":  result.Type,
		"performance_gain": result.PerformanceGain,
		"success":          result.Success,
	})

	return nil
}

// selectAdaptationStrategy selects the appropriate adaptation strategy
func (m *AdaptiveModelManager) selectAdaptationStrategy(request *AdaptationRequest) AdaptationStrategy {
	switch request.Type {
	case "performance":
		return &PerformanceBasedAdaptation{logger: m.logger}
	case "feedback":
		return &FeedbackBasedAdaptation{logger: m.logger}
	case "drift":
		return &DriftBasedAdaptation{
			logger:        m.logger,
			driftDetector: NewConceptDriftDetector(),
		}
	default:
		return nil
	}
}

// checkForAutoAdaptation checks if any models need automatic adaptation
func (m *AdaptiveModelManager) checkForAutoAdaptation(ctx context.Context) {
	for modelID, model := range m.models {
		// Check performance degradation
		if m.shouldAdaptForPerformance(model) {
			request := &AdaptationRequest{
				ModelID:     modelID,
				Type:        "performance",
				Trigger:     "auto_performance_check",
				Priority:    2,
				RequestedAt: time.Now(),
			}
			m.adaptationQueue = append(m.adaptationQueue, *request)
		}

		// Check for concept drift
		if m.shouldAdaptForDrift(model) {
			request := &AdaptationRequest{
				ModelID:     modelID,
				Type:        "drift",
				Trigger:     "auto_drift_check",
				Priority:    3,
				RequestedAt: time.Now(),
			}
			m.adaptationQueue = append(m.adaptationQueue, *request)
		}
	}
}

// shouldAdaptForPerformance checks if model needs adaptation due to performance
func (m *AdaptiveModelManager) shouldAdaptForPerformance(model *AdaptiveModel) bool {
	if model.Performance == nil || len(model.Performance.PerformanceHistory) < 10 {
		return false
	}

	// Check if performance has degraded significantly
	recent := model.Performance.PerformanceHistory[len(model.Performance.PerformanceHistory)-5:]
	historical := model.Performance.PerformanceHistory[:len(model.Performance.PerformanceHistory)-5]

	recentAvg := 0.0
	for _, point := range recent {
		recentAvg += point.Accuracy
	}
	recentAvg /= float64(len(recent))

	historicalAvg := 0.0
	for _, point := range historical {
		historicalAvg += point.Accuracy
	}
	historicalAvg /= float64(len(historical))

	degradation := historicalAvg - recentAvg
	return degradation > m.config.PerformanceThreshold
}

// shouldAdaptForDrift checks if model needs adaptation due to concept drift
func (m *AdaptiveModelManager) shouldAdaptForDrift(model *AdaptiveModel) bool {
	// Simplified drift detection - in practice would use more sophisticated methods
	return time.Since(model.LastAdaptation) > 7*24*time.Hour // Weekly check
}

// sortAdaptationQueue sorts the adaptation queue by priority and timestamp
func (m *AdaptiveModelManager) sortAdaptationQueue() {
	// Sort by priority (higher first), then by timestamp (older first)
	for i := 0; i < len(m.adaptationQueue)-1; i++ {
		for j := i + 1; j < len(m.adaptationQueue); j++ {
			if m.adaptationQueue[i].Priority < m.adaptationQueue[j].Priority ||
				(m.adaptationQueue[i].Priority == m.adaptationQueue[j].Priority &&
					m.adaptationQueue[i].RequestedAt.After(m.adaptationQueue[j].RequestedAt)) {
				m.adaptationQueue[i], m.adaptationQueue[j] = m.adaptationQueue[j], m.adaptationQueue[i]
			}
		}
	}
}

// Adaptation Strategy Implementations

// CanAdapt checks if performance-based adaptation is needed
func (p *PerformanceBasedAdaptation) CanAdapt(model *AdaptiveModel, request *AdaptationRequest) bool {
	return model.Performance != nil && model.Performance.Accuracy < 0.8
}

// Adapt performs performance-based adaptation
func (p *PerformanceBasedAdaptation) Adapt(ctx context.Context, model *AdaptiveModel, request *AdaptationRequest) (*AdaptationResult, error) {
	p.logger.Info(ctx, "Performing performance-based adaptation", map[string]interface{}{
		"model_id": model.ID,
	})

	// Simulate performance-based adaptation
	// In practice, this would involve:
	// 1. Analyzing performance degradation patterns
	// 2. Adjusting learning rate or model parameters
	// 3. Retraining with recent data
	// 4. Validating improvements

	changes := []ModelChange{
		{
			Component: "hyperparameters",
			Parameter: "learning_rate",
			OldValue:  model.LearningRate,
			NewValue:  model.LearningRate * 1.1, // Increase learning rate
			Impact:    0.05,
			Reversible: true,
		},
	}

	// Update model learning rate
	model.LearningRate *= 1.1

	return &AdaptationResult{
		Success:         true,
		Type:            "performance_optimization",
		Description:     "Adjusted learning rate to improve performance",
		PerformanceGain: 0.05,
		Confidence:      0.7,
		Changes:         changes,
		Timestamp:       time.Now(),
		Metadata: map[string]interface{}{
			"strategy": "learning_rate_adjustment",
		},
	}, nil
}

// GetType returns the adaptation strategy type
func (p *PerformanceBasedAdaptation) GetType() string {
	return "performance"
}

// FeedbackBasedAdaptation implementation
func (f *FeedbackBasedAdaptation) CanAdapt(model *AdaptiveModel, request *AdaptationRequest) bool {
	// Check if there's sufficient feedback data
	feedbackData, exists := request.Data["feedback"]
	return exists && feedbackData != nil
}

func (f *FeedbackBasedAdaptation) Adapt(ctx context.Context, model *AdaptiveModel, request *AdaptationRequest) (*AdaptationResult, error) {
	f.logger.Info(ctx, "Performing feedback-based adaptation", map[string]interface{}{
		"model_id": model.ID,
	})

	// Simulate feedback-based adaptation
	changes := []ModelChange{
		{
			Component: "weights",
			Parameter: "output_layer",
			OldValue:  "original_weights",
			NewValue:  "adjusted_weights",
			Impact:    0.03,
			Reversible: true,
		},
	}

	return &AdaptationResult{
		Success:         true,
		Type:            "feedback_integration",
		Description:     "Integrated user feedback into model weights",
		PerformanceGain: 0.03,
		Confidence:      0.8,
		Changes:         changes,
		Timestamp:       time.Now(),
		Metadata: map[string]interface{}{
			"strategy": "feedback_integration",
		},
	}, nil
}

func (f *FeedbackBasedAdaptation) GetType() string {
	return "feedback"
}

// DriftBasedAdaptation implementation
func (d *DriftBasedAdaptation) CanAdapt(model *AdaptiveModel, request *AdaptationRequest) bool {
	// Check if drift has been detected
	return d.driftDetector.HasDrift()
}

func (d *DriftBasedAdaptation) Adapt(ctx context.Context, model *AdaptiveModel, request *AdaptationRequest) (*AdaptationResult, error) {
	d.logger.Info(ctx, "Performing drift-based adaptation", map[string]interface{}{
		"model_id": model.ID,
	})

	// Simulate drift-based adaptation
	changes := []ModelChange{
		{
			Component: "architecture",
			Parameter: "feature_weights",
			OldValue:  "original_features",
			NewValue:  "drift_adjusted_features",
			Impact:    0.08,
			Reversible: false,
		},
	}

	return &AdaptationResult{
		Success:         true,
		Type:            "drift_compensation",
		Description:     "Adjusted model to compensate for concept drift",
		PerformanceGain: 0.08,
		Confidence:      0.75,
		Changes:         changes,
		Timestamp:       time.Now(),
		Metadata: map[string]interface{}{
			"strategy": "drift_compensation",
			"drift_type": "gradual",
		},
	}, nil
}

func (d *DriftBasedAdaptation) GetType() string {
	return "drift"
}

// ConceptDriftDetector implementation
func NewConceptDriftDetector() *ConceptDriftDetector {
	return &ConceptDriftDetector{
		windowSize:     1000,
		driftThreshold: 0.1,
		recentData:     []DataPoint{},
		referenceData:  []DataPoint{},
		lastCheck:      time.Now(),
		driftHistory:   []DriftEvent{},
	}
}

func (c *ConceptDriftDetector) HasDrift() bool {
	// Simplified drift detection
	return len(c.driftHistory) > 0 && time.Since(c.driftHistory[len(c.driftHistory)-1].DetectedAt) < 24*time.Hour
}

func (c *ConceptDriftDetector) AddDataPoint(point DataPoint) {
	c.recentData = append(c.recentData, point)
	
	// Maintain window size
	if len(c.recentData) > c.windowSize {
		c.recentData = c.recentData[1:]
	}
	
	// Check for drift periodically
	if time.Since(c.lastCheck) > 1*time.Hour {
		c.checkForDrift()
		c.lastCheck = time.Now()
	}
}

func (c *ConceptDriftDetector) checkForDrift() {
	if len(c.recentData) < c.windowSize/2 || len(c.referenceData) < c.windowSize/2 {
		return
	}

	// Simplified drift detection using statistical distance
	distance := c.calculateStatisticalDistance(c.recentData, c.referenceData)
	
	if distance > c.driftThreshold {
		event := DriftEvent{
			Type:        "gradual",
			Severity:    distance,
			Confidence:  0.8,
			DetectedAt:  time.Now(),
			Description: fmt.Sprintf("Concept drift detected with distance %.3f", distance),
			Features:    []string{"all"}, // Simplified
		}
		
		c.driftHistory = append(c.driftHistory, event)
		
		// Update reference data
		c.referenceData = make([]DataPoint, len(c.recentData))
		copy(c.referenceData, c.recentData)
	}
}

func (c *ConceptDriftDetector) calculateStatisticalDistance(data1, data2 []DataPoint) float64 {
	// Simplified statistical distance calculation
	// In practice, would use more sophisticated methods like KL divergence, Wasserstein distance, etc.
	
	if len(data1) == 0 || len(data2) == 0 {
		return 0.0
	}
	
	// Calculate mean differences for each feature
	features1 := make(map[string]float64)
	features2 := make(map[string]float64)
	
	for _, point := range data1 {
		for feature, value := range point.Features {
			features1[feature] += value
		}
	}
	
	for _, point := range data2 {
		for feature, value := range point.Features {
			features2[feature] += value
		}
	}
	
	// Normalize by count
	for feature := range features1 {
		features1[feature] /= float64(len(data1))
	}
	for feature := range features2 {
		features2[feature] /= float64(len(data2))
	}
	
	// Calculate Euclidean distance
	distance := 0.0
	for feature := range features1 {
		if val2, exists := features2[feature]; exists {
			diff := features1[feature] - val2
			distance += diff * diff
		}
	}
	
	return math.Sqrt(distance)
}

// GetAdaptiveModels returns all adaptive models
func (m *AdaptiveModelManager) GetAdaptiveModels() map[string]*AdaptiveModel {
	models := make(map[string]*AdaptiveModel)
	for id, model := range m.models {
		models[id] = model
	}
	return models
}

// GetAdaptationHistory returns adaptation history for a model
func (m *AdaptiveModelManager) GetAdaptationHistory(modelID string) ([]ModelAdaptation, error) {
	model, exists := m.models[modelID]
	if !exists {
		return nil, fmt.Errorf("model %s not found", modelID)
	}
	
	return model.Adaptations, nil
}
