package ai

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/ml"
	"github.com/ai-agentic-browser/pkg/observability"
)

// RealTimeLearningEngine implements continuous learning from streaming data
type RealTimeLearningEngine struct {
	logger               *observability.Logger
	config               *RealTimeLearningConfig
	models               map[string]*OnlineModel
	dataStream           chan *RealTimeLearningEvent
	feedbackStream       chan *RealTimeFeedbackEvent
	conceptDriftDetector *RealTimeConceptDriftDetector
	performanceMonitor   *OnlinePerformanceMonitor
	adaptationEngine     *ModelAdaptationEngine
	mu                   sync.RWMutex
	isRunning            bool
	stopChan             chan struct{}
}

// RealTimeLearningConfig contains configuration for real-time learning
type RealTimeLearningConfig struct {
	LearningRate         float64
	BatchSize            int
	UpdateFrequency      time.Duration
	DriftThreshold       float64
	PerformanceWindow    int
	AdaptationThreshold  float64
	MaxModelAge          time.Duration
	EnableDriftDetection bool
	EnableAdaptation     bool
}

// OnlineModel represents a model that can learn incrementally
type OnlineModel struct {
	ID                 string
	BaseModel          ml.Model
	LearningRate       float64
	LastUpdate         time.Time
	UpdateCount        int64
	PerformanceHistory []float64
	DriftScore         float64
	IsAdapting         bool
	mu                 sync.RWMutex
}

// RealTimeLearningEvent represents a new data point for learning
type RealTimeLearningEvent struct {
	EventID   string
	ModelID   string
	Features  map[string]interface{}
	Target    interface{}
	Timestamp time.Time
	Weight    float64
	Metadata  map[string]interface{}
}

// RealTimeFeedbackEvent represents feedback on a prediction
type RealTimeFeedbackEvent struct {
	EventID      string
	ModelID      string
	PredictionID string
	Prediction   interface{}
	Actual       interface{}
	Error        float64
	Timestamp    time.Time
	Metadata     map[string]interface{}
}

// RealTimeConceptDriftDetector detects when the data distribution changes
type RealTimeConceptDriftDetector struct {
	logger          *observability.Logger
	config          *DriftDetectionConfig
	referenceWindow []float64
	currentWindow   []float64
	windowSize      int
	driftThreshold  float64
	lastDriftTime   time.Time
	mu              sync.RWMutex
}

// DriftDetectionConfig contains drift detection configuration
type DriftDetectionConfig struct {
	WindowSize      int
	DriftThreshold  float64
	MinSamples      int
	StatisticalTest string // "ks", "chi2", "psi"
}

// OnlinePerformanceMonitor tracks model performance in real-time
type OnlinePerformanceMonitor struct {
	logger            *observability.Logger
	modelMetrics      map[string]*OnlineMetrics
	performanceWindow int
	mu                sync.RWMutex
}

// OnlineMetrics tracks real-time performance metrics
type OnlineMetrics struct {
	Accuracy        float64
	Precision       float64
	Recall          float64
	F1Score         float64
	MAE             float64
	RMSE            float64
	RecentErrors    []float64
	PredictionCount int64
	LastUpdated     time.Time
}

// ModelAdaptationEngine handles automatic model adaptation
type ModelAdaptationEngine struct {
	logger               *observability.Logger
	config               *AdaptationConfig
	adaptationStrategies map[string]RealTimeAdaptationStrategy
	mu                   sync.RWMutex
}

// AdaptationConfig contains adaptation configuration
type AdaptationConfig struct {
	PerformanceThreshold float64
	DriftThreshold       float64
	AdaptationStrategies []string
	CooldownPeriod       time.Duration
}

// RealTimeAdaptationStrategy defines how models should adapt
type RealTimeAdaptationStrategy interface {
	ShouldAdapt(model *OnlineModel, metrics *OnlineMetrics, driftScore float64) bool
	Adapt(ctx context.Context, model *OnlineModel, data []*RealTimeLearningEvent) error
}

// NewRealTimeLearningEngine creates a new real-time learning engine
func NewRealTimeLearningEngine(logger *observability.Logger) *RealTimeLearningEngine {
	config := &RealTimeLearningConfig{
		LearningRate:         0.01,
		BatchSize:            32,
		UpdateFrequency:      1 * time.Minute,
		DriftThreshold:       0.05,
		PerformanceWindow:    100,
		AdaptationThreshold:  0.1,
		MaxModelAge:          24 * time.Hour,
		EnableDriftDetection: true,
		EnableAdaptation:     true,
	}

	driftConfig := &DriftDetectionConfig{
		WindowSize:      100,
		DriftThreshold:  0.05,
		MinSamples:      30,
		StatisticalTest: "ks",
	}

	adaptationConfig := &AdaptationConfig{
		PerformanceThreshold: 0.8,
		DriftThreshold:       0.05,
		AdaptationStrategies: []string{"incremental_update", "model_retrain", "ensemble_update"},
		CooldownPeriod:       30 * time.Minute,
	}

	return &RealTimeLearningEngine{
		logger:               logger,
		config:               config,
		models:               make(map[string]*OnlineModel),
		dataStream:           make(chan *RealTimeLearningEvent, 1000),
		feedbackStream:       make(chan *RealTimeFeedbackEvent, 1000),
		conceptDriftDetector: NewRealTimeConceptDriftDetector(logger, driftConfig),
		performanceMonitor:   NewOnlinePerformanceMonitor(logger, config.PerformanceWindow),
		adaptationEngine:     NewModelAdaptationEngine(logger, adaptationConfig),
		stopChan:             make(chan struct{}),
	}
}

// Start starts the real-time learning engine
func (r *RealTimeLearningEngine) Start(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.isRunning {
		return fmt.Errorf("real-time learning engine is already running")
	}

	r.isRunning = true

	// Start processing goroutines
	go r.processDataStream(ctx)
	go r.processFeedbackStream(ctx)
	go r.performPeriodicUpdates(ctx)

	r.logger.Info(ctx, "Real-time learning engine started", map[string]interface{}{
		"learning_rate":      r.config.LearningRate,
		"batch_size":         r.config.BatchSize,
		"update_frequency":   r.config.UpdateFrequency,
		"drift_detection":    r.config.EnableDriftDetection,
		"adaptation_enabled": r.config.EnableAdaptation,
	})

	return nil
}

// Stop stops the real-time learning engine
func (r *RealTimeLearningEngine) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.isRunning {
		return
	}

	close(r.stopChan)
	r.isRunning = false

	r.logger.Info(context.Background(), "Real-time learning engine stopped")
}

// AddModel adds a model for real-time learning
func (r *RealTimeLearningEngine) AddModel(modelID string, baseModel ml.Model) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.models[modelID]; exists {
		return fmt.Errorf("model %s already exists", modelID)
	}

	onlineModel := &OnlineModel{
		ID:                 modelID,
		BaseModel:          baseModel,
		LearningRate:       r.config.LearningRate,
		LastUpdate:         time.Now(),
		UpdateCount:        0,
		PerformanceHistory: make([]float64, 0),
		DriftScore:         0.0,
		IsAdapting:         false,
	}

	r.models[modelID] = onlineModel
	r.performanceMonitor.InitializeModel(modelID)

	r.logger.Info(context.Background(), "Model added for real-time learning", map[string]interface{}{
		"model_id":      modelID,
		"learning_rate": onlineModel.LearningRate,
	})

	return nil
}

// LearnFromData adds new data for learning
func (r *RealTimeLearningEngine) LearnFromData(event *RealTimeLearningEvent) error {
	if !r.isRunning {
		return fmt.Errorf("real-time learning engine is not running")
	}

	select {
	case r.dataStream <- event:
		return nil
	default:
		return fmt.Errorf("data stream buffer is full")
	}
}

// ProvideFeedback provides feedback on predictions
func (r *RealTimeLearningEngine) ProvideFeedback(feedback *RealTimeFeedbackEvent) error {
	if !r.isRunning {
		return fmt.Errorf("real-time learning engine is not running")
	}

	select {
	case r.feedbackStream <- feedback:
		return nil
	default:
		return fmt.Errorf("feedback stream buffer is full")
	}
}

// processDataStream processes incoming learning events
func (r *RealTimeLearningEngine) processDataStream(ctx context.Context) {
	batch := make([]*RealTimeLearningEvent, 0, r.config.BatchSize)
	ticker := time.NewTicker(r.config.UpdateFrequency)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-r.stopChan:
			return
		case event := <-r.dataStream:
			batch = append(batch, event)

			// Process batch when full
			if len(batch) >= r.config.BatchSize {
				r.processBatch(ctx, batch)
				batch = batch[:0] // Reset batch
			}
		case <-ticker.C:
			// Process partial batch on timer
			if len(batch) > 0 {
				r.processBatch(ctx, batch)
				batch = batch[:0]
			}
		}
	}
}

// processFeedbackStream processes incoming feedback events
func (r *RealTimeLearningEngine) processFeedbackStream(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-r.stopChan:
			return
		case feedback := <-r.feedbackStream:
			r.processFeedback(ctx, feedback)
		}
	}
}

// processBatch processes a batch of learning events
func (r *RealTimeLearningEngine) processBatch(ctx context.Context, batch []*RealTimeLearningEvent) {
	// Group events by model
	modelBatches := make(map[string][]*RealTimeLearningEvent)
	for _, event := range batch {
		modelBatches[event.ModelID] = append(modelBatches[event.ModelID], event)
	}

	// Process each model's batch
	for modelID, events := range modelBatches {
		r.mu.RLock()
		model, exists := r.models[modelID]
		r.mu.RUnlock()

		if !exists {
			r.logger.Warn(ctx, "Model not found for learning event", map[string]interface{}{
				"model_id": modelID,
			})
			continue
		}

		// Update model with new data
		if err := r.updateModel(ctx, model, events); err != nil {
			r.logger.Error(ctx, "Failed to update model", err, map[string]interface{}{
				"model_id":   modelID,
				"batch_size": len(events),
			})
		}
	}
}

// updateModel updates a model with new learning events
func (r *RealTimeLearningEngine) updateModel(ctx context.Context, model *OnlineModel, events []*RealTimeLearningEvent) error {
	model.mu.Lock()
	defer model.mu.Unlock()

	// Prepare training data
	features := make([]map[string]interface{}, len(events))
	targets := make([]interface{}, len(events))
	weights := make([]float64, len(events))

	for i, event := range events {
		features[i] = event.Features
		targets[i] = event.Target
		weights[i] = event.Weight
	}

	trainingData := ml.TrainingData{
		Features: features,
		Labels:   targets,
		Weights:  weights,
	}

	// Perform incremental learning
	if err := model.BaseModel.Train(ctx, trainingData); err != nil {
		return fmt.Errorf("incremental training failed: %w", err)
	}

	// Update model metadata
	model.LastUpdate = time.Now()
	model.UpdateCount += int64(len(events))

	// Check for concept drift if enabled
	if r.config.EnableDriftDetection {
		driftScore := r.conceptDriftDetector.DetectDrift(events)
		model.DriftScore = driftScore

		if driftScore > r.config.DriftThreshold {
			r.logger.Warn(ctx, "Concept drift detected", map[string]interface{}{
				"model_id":    model.ID,
				"drift_score": driftScore,
				"threshold":   r.config.DriftThreshold,
			})

			// Trigger adaptation if enabled
			if r.config.EnableAdaptation {
				go r.triggerAdaptation(ctx, model, events)
			}
		}
	}

	r.logger.Debug(ctx, "Model updated with new data", map[string]interface{}{
		"model_id":     model.ID,
		"batch_size":   len(events),
		"update_count": model.UpdateCount,
		"drift_score":  model.DriftScore,
	})

	return nil
}

// processFeedback processes feedback events
func (r *RealTimeLearningEngine) processFeedback(ctx context.Context, feedback *RealTimeFeedbackEvent) {
	// Update performance metrics
	r.performanceMonitor.UpdateMetrics(feedback.ModelID, feedback.Error, feedback.Actual, feedback.Prediction)

	// Use feedback for model improvement
	r.mu.RLock()
	model, exists := r.models[feedback.ModelID]
	r.mu.RUnlock()

	if !exists {
		r.logger.Warn(ctx, "Model not found for feedback", map[string]interface{}{
			"model_id": feedback.ModelID,
		})
		return
	}

	// Create feedback for model weight updates
	modelFeedback := &ml.PredictionFeedback{
		PredictionID: feedback.PredictionID,
		Correct:      feedback.Error < 0.1, // Threshold for correctness
		Confidence:   1.0 - feedback.Error,
		ActualValue:  feedback.Actual,
		Metadata:     feedback.Metadata,
	}

	if err := model.BaseModel.UpdateWeights(ctx, modelFeedback); err != nil {
		r.logger.Error(ctx, "Failed to update model weights", err, map[string]interface{}{
			"model_id":      feedback.ModelID,
			"prediction_id": feedback.PredictionID,
		})
	}
}

// performPeriodicUpdates performs periodic maintenance and optimization
func (r *RealTimeLearningEngine) performPeriodicUpdates(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute) // Periodic maintenance every 5 minutes
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-r.stopChan:
			return
		case <-ticker.C:
			r.performMaintenance(ctx)
		}
	}
}

// performMaintenance performs periodic maintenance tasks
func (r *RealTimeLearningEngine) performMaintenance(ctx context.Context) {
	r.mu.RLock()
	models := make(map[string]*OnlineModel)
	for id, model := range r.models {
		models[id] = model
	}
	r.mu.RUnlock()

	for modelID, model := range models {
		// Check if model needs adaptation
		metrics := r.performanceMonitor.GetMetrics(modelID)
		if metrics != nil && r.adaptationEngine.ShouldAdapt(model, metrics) {
			r.logger.Info(ctx, "Triggering model adaptation", map[string]interface{}{
				"model_id":    modelID,
				"accuracy":    metrics.Accuracy,
				"drift_score": model.DriftScore,
			})

			go r.triggerAdaptation(ctx, model, nil)
		}

		// Clean up old performance data
		r.performanceMonitor.CleanupOldData(modelID)
	}
}

// triggerAdaptation triggers model adaptation
func (r *RealTimeLearningEngine) triggerAdaptation(ctx context.Context, model *OnlineModel, recentData []*RealTimeLearningEvent) {
	model.mu.Lock()
	if model.IsAdapting {
		model.mu.Unlock()
		return // Already adapting
	}
	model.IsAdapting = true
	model.mu.Unlock()

	defer func() {
		model.mu.Lock()
		model.IsAdapting = false
		model.mu.Unlock()
	}()

	if err := r.adaptationEngine.AdaptModel(ctx, model, recentData); err != nil {
		r.logger.Error(ctx, "Model adaptation failed", err, map[string]interface{}{
			"model_id": model.ID,
		})
	} else {
		r.logger.Info(ctx, "Model adaptation completed", map[string]interface{}{
			"model_id": model.ID,
		})
	}
}

// GetModelMetrics returns real-time metrics for a model
func (r *RealTimeLearningEngine) GetModelMetrics(modelID string) map[string]interface{} {
	r.mu.RLock()
	model, exists := r.models[modelID]
	r.mu.RUnlock()

	if !exists {
		return nil
	}

	metrics := r.performanceMonitor.GetMetrics(modelID)

	result := map[string]interface{}{
		"model_id":      model.ID,
		"last_update":   model.LastUpdate,
		"update_count":  model.UpdateCount,
		"drift_score":   model.DriftScore,
		"is_adapting":   model.IsAdapting,
		"learning_rate": model.LearningRate,
	}

	if metrics != nil {
		result["accuracy"] = metrics.Accuracy
		result["precision"] = metrics.Precision
		result["recall"] = metrics.Recall
		result["f1_score"] = metrics.F1Score
		result["mae"] = metrics.MAE
		result["rmse"] = metrics.RMSE
		result["prediction_count"] = metrics.PredictionCount
	}

	return result
}

// GetSystemMetrics returns overall system metrics
func (r *RealTimeLearningEngine) GetSystemMetrics() map[string]interface{} {
	r.mu.RLock()
	modelCount := len(r.models)
	r.mu.RUnlock()

	return map[string]interface{}{
		"is_running":           r.isRunning,
		"model_count":          modelCount,
		"data_stream_size":     len(r.dataStream),
		"feedback_stream_size": len(r.feedbackStream),
		"learning_rate":        r.config.LearningRate,
		"batch_size":           r.config.BatchSize,
		"update_frequency":     r.config.UpdateFrequency,
		"drift_detection":      r.config.EnableDriftDetection,
		"adaptation_enabled":   r.config.EnableAdaptation,
	}
}

// NewRealTimeConceptDriftDetector creates a new concept drift detector
func NewRealTimeConceptDriftDetector(logger *observability.Logger, config *DriftDetectionConfig) *RealTimeConceptDriftDetector {
	return &RealTimeConceptDriftDetector{
		logger:          logger,
		config:          config,
		referenceWindow: make([]float64, 0, config.WindowSize),
		currentWindow:   make([]float64, 0, config.WindowSize),
		windowSize:      config.WindowSize,
		driftThreshold:  config.DriftThreshold,
		lastDriftTime:   time.Now(),
	}
}

// DetectDrift detects concept drift in the data
func (c *RealTimeConceptDriftDetector) DetectDrift(events []*RealTimeLearningEvent) float64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Extract numeric features for drift detection
	for _, event := range events {
		if value, ok := event.Features["target_value"].(float64); ok {
			c.currentWindow = append(c.currentWindow, value)

			// Maintain window size
			if len(c.currentWindow) > c.windowSize {
				c.currentWindow = c.currentWindow[1:]
			}
		}
	}

	// Need sufficient data for comparison
	if len(c.referenceWindow) < c.config.MinSamples || len(c.currentWindow) < c.config.MinSamples {
		return 0.0
	}

	// Perform statistical test
	var driftScore float64
	switch c.config.StatisticalTest {
	case "ks":
		driftScore = c.kolmogorovSmirnovTest()
	case "chi2":
		driftScore = c.chiSquareTest()
	case "psi":
		driftScore = c.populationStabilityIndex()
	default:
		driftScore = c.kolmogorovSmirnovTest()
	}

	// Update reference window if significant drift detected
	if driftScore > c.driftThreshold {
		c.referenceWindow = make([]float64, len(c.currentWindow))
		copy(c.referenceWindow, c.currentWindow)
		c.lastDriftTime = time.Now()
	}

	return driftScore
}

// kolmogorovSmirnovTest performs Kolmogorov-Smirnov test
func (c *RealTimeConceptDriftDetector) kolmogorovSmirnovTest() float64 {
	// Simplified KS test implementation
	// In production, use a proper statistical library

	// Sort both samples
	ref := make([]float64, len(c.referenceWindow))
	curr := make([]float64, len(c.currentWindow))
	copy(ref, c.referenceWindow)
	copy(curr, c.currentWindow)

	// Calculate empirical CDFs and find maximum difference
	maxDiff := 0.0
	for i := 0; i < len(ref); i++ {
		refCDF := float64(i) / float64(len(ref))

		// Find corresponding position in current window
		currPos := 0
		for j := 0; j < len(curr); j++ {
			if curr[j] <= ref[i] {
				currPos = j + 1
			}
		}
		currCDF := float64(currPos) / float64(len(curr))

		diff := math.Abs(refCDF - currCDF)
		if diff > maxDiff {
			maxDiff = diff
		}
	}

	return maxDiff
}

// chiSquareTest performs Chi-square test
func (c *RealTimeConceptDriftDetector) chiSquareTest() float64 {
	// Simplified Chi-square test
	// Bin the data and compare distributions
	bins := 10
	refHist := make([]int, bins)
	currHist := make([]int, bins)

	// Find min/max for binning
	minVal, maxVal := math.Inf(1), math.Inf(-1)
	for _, val := range c.referenceWindow {
		if val < minVal {
			minVal = val
		}
		if val > maxVal {
			maxVal = val
		}
	}
	for _, val := range c.currentWindow {
		if val < minVal {
			minVal = val
		}
		if val > maxVal {
			maxVal = val
		}
	}

	binWidth := (maxVal - minVal) / float64(bins)
	if binWidth == 0 {
		return 0.0
	}

	// Create histograms
	for _, val := range c.referenceWindow {
		bin := int((val - minVal) / binWidth)
		if bin >= bins {
			bin = bins - 1
		}
		refHist[bin]++
	}

	for _, val := range c.currentWindow {
		bin := int((val - minVal) / binWidth)
		if bin >= bins {
			bin = bins - 1
		}
		currHist[bin]++
	}

	// Calculate chi-square statistic
	chiSquare := 0.0
	for i := 0; i < bins; i++ {
		expected := float64(refHist[i])
		observed := float64(currHist[i])
		if expected > 0 {
			chiSquare += math.Pow(observed-expected, 2) / expected
		}
	}

	// Normalize to [0, 1]
	return math.Min(1.0, chiSquare/float64(bins))
}

// populationStabilityIndex calculates Population Stability Index
func (c *RealTimeConceptDriftDetector) populationStabilityIndex() float64 {
	// Simplified PSI calculation
	bins := 10
	refHist := make([]float64, bins)
	currHist := make([]float64, bins)

	// Create normalized histograms
	// ... (similar binning logic as chi-square)

	psi := 0.0
	for i := 0; i < bins; i++ {
		if refHist[i] > 0 && currHist[i] > 0 {
			psi += (currHist[i] - refHist[i]) * math.Log(currHist[i]/refHist[i])
		}
	}

	return math.Min(1.0, psi)
}

// NewOnlinePerformanceMonitor creates a new online performance monitor
func NewOnlinePerformanceMonitor(logger *observability.Logger, windowSize int) *OnlinePerformanceMonitor {
	return &OnlinePerformanceMonitor{
		logger:            logger,
		modelMetrics:      make(map[string]*OnlineMetrics),
		performanceWindow: windowSize,
	}
}

// InitializeModel initializes metrics for a new model
func (o *OnlinePerformanceMonitor) InitializeModel(modelID string) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.modelMetrics[modelID] = &OnlineMetrics{
		RecentErrors:    make([]float64, 0, o.performanceWindow),
		PredictionCount: 0,
		LastUpdated:     time.Now(),
	}
}

// UpdateMetrics updates performance metrics for a model
func (o *OnlinePerformanceMonitor) UpdateMetrics(modelID string, error float64, actual, predicted interface{}) {
	o.mu.Lock()
	defer o.mu.Unlock()

	metrics, exists := o.modelMetrics[modelID]
	if !exists {
		return
	}

	// Update error tracking
	metrics.RecentErrors = append(metrics.RecentErrors, error)
	if len(metrics.RecentErrors) > o.performanceWindow {
		metrics.RecentErrors = metrics.RecentErrors[1:]
	}

	// Calculate running metrics
	metrics.PredictionCount++

	// Calculate accuracy (for classification) or MAE/RMSE (for regression)
	if len(metrics.RecentErrors) > 0 {
		totalError := 0.0
		squaredError := 0.0
		correctPredictions := 0

		for _, err := range metrics.RecentErrors {
			totalError += math.Abs(err)
			squaredError += err * err
			if err < 0.1 { // Threshold for correctness
				correctPredictions++
			}
		}

		metrics.Accuracy = float64(correctPredictions) / float64(len(metrics.RecentErrors))
		metrics.MAE = totalError / float64(len(metrics.RecentErrors))
		metrics.RMSE = math.Sqrt(squaredError / float64(len(metrics.RecentErrors)))
	}

	metrics.LastUpdated = time.Now()
}

// GetMetrics returns metrics for a model
func (o *OnlinePerformanceMonitor) GetMetrics(modelID string) *OnlineMetrics {
	o.mu.RLock()
	defer o.mu.RUnlock()

	if metrics, exists := o.modelMetrics[modelID]; exists {
		// Return a copy to avoid race conditions
		return &OnlineMetrics{
			Accuracy:        metrics.Accuracy,
			Precision:       metrics.Precision,
			Recall:          metrics.Recall,
			F1Score:         metrics.F1Score,
			MAE:             metrics.MAE,
			RMSE:            metrics.RMSE,
			PredictionCount: metrics.PredictionCount,
			LastUpdated:     metrics.LastUpdated,
		}
	}

	return nil
}

// CleanupOldData removes old performance data
func (o *OnlinePerformanceMonitor) CleanupOldData(modelID string) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if metrics, exists := o.modelMetrics[modelID]; exists {
		// Keep only recent errors within the window
		if len(metrics.RecentErrors) > o.performanceWindow {
			metrics.RecentErrors = metrics.RecentErrors[len(metrics.RecentErrors)-o.performanceWindow:]
		}
	}
}

// NewModelAdaptationEngine creates a new model adaptation engine
func NewModelAdaptationEngine(logger *observability.Logger, config *AdaptationConfig) *ModelAdaptationEngine {
	return &ModelAdaptationEngine{
		logger:               logger,
		config:               config,
		adaptationStrategies: make(map[string]RealTimeAdaptationStrategy),
	}
}

// ShouldAdapt determines if a model should be adapted
func (m *ModelAdaptationEngine) ShouldAdapt(model *OnlineModel, metrics *OnlineMetrics) bool {
	// Check performance threshold
	if metrics.Accuracy < m.config.PerformanceThreshold {
		return true
	}

	// Check drift threshold
	if model.DriftScore > m.config.DriftThreshold {
		return true
	}

	// Check if enough time has passed since last adaptation
	if time.Since(model.LastUpdate) < m.config.CooldownPeriod {
		return false
	}

	return false
}

// AdaptModel adapts a model based on current conditions
func (m *ModelAdaptationEngine) AdaptModel(ctx context.Context, model *OnlineModel, recentData []*RealTimeLearningEvent) error {
	// Simple adaptation: increase learning rate if performance is poor
	model.mu.Lock()
	defer model.mu.Unlock()

	if model.DriftScore > m.config.DriftThreshold {
		// Increase learning rate for faster adaptation to drift
		model.LearningRate = math.Min(0.1, model.LearningRate*1.5)
	} else {
		// Decrease learning rate for stability
		model.LearningRate = math.Max(0.001, model.LearningRate*0.9)
	}

	m.logger.Info(ctx, "Model adapted", map[string]interface{}{
		"model_id":    model.ID,
		"new_lr":      model.LearningRate,
		"drift_score": model.DriftScore,
	})

	return nil
}
