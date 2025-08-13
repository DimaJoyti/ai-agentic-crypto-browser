package ml

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

// MLFramework provides machine learning capabilities for trading
type MLFramework struct {
	logger        *observability.Logger
	config        MLConfig
	models        map[string]*Model
	featureStore  *FeatureStore
	dataProcessor *DataProcessor
	modelRegistry *ModelRegistry

	// Training pipeline
	trainingQueue chan *TrainingJob
	trainedModels map[string]*TrainedModel

	// Prediction pipeline
	predictionCache map[string]*PredictionResult
	cacheExpiry     map[string]time.Time

	// Performance tracking
	totalPredictions int64
	totalTraining    int64
	modelAccuracy    map[string]float64

	// State management
	isRunning bool
	stopChan  chan struct{}
	wg        sync.WaitGroup
	mu        sync.RWMutex
}

// MLConfig contains machine learning configuration
type MLConfig struct {
	EnableTraining          bool          `json:"enable_training"`
	EnablePrediction        bool          `json:"enable_prediction"`
	EnableFeatureStore      bool          `json:"enable_feature_store"`
	TrainingInterval        time.Duration `json:"training_interval"`
	PredictionCacheTime     time.Duration `json:"prediction_cache_time"`
	MaxConcurrentTraining   int           `json:"max_concurrent_training"`
	MaxConcurrentPrediction int           `json:"max_concurrent_prediction"`
	FeatureRetention        time.Duration `json:"feature_retention"`
	ModelRetention          time.Duration `json:"model_retention"`
	EnableGPU               bool          `json:"enable_gpu"`
	ModelStoragePath        string        `json:"model_storage_path"`
}

// Model represents a machine learning model
type Model struct {
	ID            uuid.UUID              `json:"id"`
	Name          string                 `json:"name"`
	Type          ModelType              `json:"type"`
	Algorithm     Algorithm              `json:"algorithm"`
	Version       string                 `json:"version"`
	Description   string                 `json:"description"`
	InputFeatures []string               `json:"input_features"`
	OutputTargets []string               `json:"output_targets"`
	Hyperparams   map[string]interface{} `json:"hyperparams"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	Status        ModelStatus            `json:"status"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// ModelType represents the type of ML model
type ModelType string

const (
	ModelTypeRegression     ModelType = "regression"
	ModelTypeClassification ModelType = "classification"
	ModelTypeTimeSeries     ModelType = "time_series"
	ModelTypeClustering     ModelType = "clustering"
	ModelTypeReinforcement  ModelType = "reinforcement"
	ModelTypeEnsemble       ModelType = "ensemble"
	ModelTypeDeepLearning   ModelType = "deep_learning"
)

// Algorithm represents the ML algorithm
type Algorithm string

const (
	AlgorithmLinearRegression   Algorithm = "linear_regression"
	AlgorithmLogisticRegression Algorithm = "logistic_regression"
	AlgorithmRandomForest       Algorithm = "random_forest"
	AlgorithmGradientBoosting   Algorithm = "gradient_boosting"
	AlgorithmSVM                Algorithm = "svm"
	AlgorithmNeuralNetwork      Algorithm = "neural_network"
	AlgorithmLSTM               Algorithm = "lstm"
	AlgorithmGRU                Algorithm = "gru"
	AlgorithmTransformer        Algorithm = "transformer"
	AlgorithmReinforcement      Algorithm = "reinforcement"
	AlgorithmEnsemble           Algorithm = "ensemble"
)

// ModelStatus represents the status of a model
type ModelStatus string

const (
	ModelStatusCreated    ModelStatus = "created"
	ModelStatusTraining   ModelStatus = "training"
	ModelStatusTrained    ModelStatus = "trained"
	ModelStatusDeployed   ModelStatus = "deployed"
	ModelStatusFailed     ModelStatus = "failed"
	ModelStatusDeprecated ModelStatus = "deprecated"
)

// TrainedModel represents a trained ML model
type TrainedModel struct {
	ModelID        uuid.UUID              `json:"model_id"`
	Version        string                 `json:"version"`
	TrainingData   *TrainingDataset       `json:"training_data"`
	ValidationData *ValidationDataset     `json:"validation_data"`
	Metrics        *ModelMetrics          `json:"metrics"`
	Weights        []byte                 `json:"weights,omitempty"`
	TrainedAt      time.Time              `json:"trained_at"`
	TrainingTime   time.Duration          `json:"training_time"`
	Status         ModelStatus            `json:"status"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// TrainingJob represents a model training job
type TrainingJob struct {
	ID          uuid.UUID              `json:"id"`
	ModelID     uuid.UUID              `json:"model_id"`
	Dataset     *TrainingDataset       `json:"dataset"`
	Config      *TrainingConfig        `json:"config"`
	Priority    int                    `json:"priority"`
	CreatedAt   time.Time              `json:"created_at"`
	StartedAt   time.Time              `json:"started_at,omitempty"`
	CompletedAt time.Time              `json:"completed_at,omitempty"`
	Status      TrainingStatus         `json:"status"`
	Progress    float64                `json:"progress"`
	Error       string                 `json:"error,omitempty"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// TrainingStatus represents the status of a training job
type TrainingStatus string

const (
	TrainingStatusQueued    TrainingStatus = "queued"
	TrainingStatusRunning   TrainingStatus = "running"
	TrainingStatusCompleted TrainingStatus = "completed"
	TrainingStatusFailed    TrainingStatus = "failed"
	TrainingStatusCancelled TrainingStatus = "cancelled"
)

// TrainingDataset represents training data
type TrainingDataset struct {
	ID           uuid.UUID              `json:"id"`
	Name         string                 `json:"name"`
	Features     [][]float64            `json:"features"`
	Targets      [][]float64            `json:"targets"`
	FeatureNames []string               `json:"feature_names"`
	TargetNames  []string               `json:"target_names"`
	Size         int                    `json:"size"`
	CreatedAt    time.Time              `json:"created_at"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// ValidationDataset represents validation data
type ValidationDataset struct {
	ID       uuid.UUID              `json:"id"`
	Features [][]float64            `json:"features"`
	Targets  [][]float64            `json:"targets"`
	Size     int                    `json:"size"`
	Metadata map[string]interface{} `json:"metadata"`
}

// TrainingConfig contains training configuration
type TrainingConfig struct {
	Epochs          int                    `json:"epochs"`
	BatchSize       int                    `json:"batch_size"`
	LearningRate    float64                `json:"learning_rate"`
	ValidationSplit float64                `json:"validation_split"`
	EarlyStopping   bool                   `json:"early_stopping"`
	Patience        int                    `json:"patience"`
	Regularization  map[string]interface{} `json:"regularization"`
	Optimizer       string                 `json:"optimizer"`
	LossFunction    string                 `json:"loss_function"`
	Metrics         []string               `json:"metrics"`
	Callbacks       []string               `json:"callbacks"`
}

// ModelMetrics contains model performance metrics
type ModelMetrics struct {
	Accuracy      float64            `json:"accuracy,omitempty"`
	Precision     float64            `json:"precision,omitempty"`
	Recall        float64            `json:"recall,omitempty"`
	F1Score       float64            `json:"f1_score,omitempty"`
	MSE           float64            `json:"mse,omitempty"`
	RMSE          float64            `json:"rmse,omitempty"`
	MAE           float64            `json:"mae,omitempty"`
	R2Score       float64            `json:"r2_score,omitempty"`
	AUC           float64            `json:"auc,omitempty"`
	LogLoss       float64            `json:"log_loss,omitempty"`
	SharpeRatio   float64            `json:"sharpe_ratio,omitempty"`
	MaxDrawdown   float64            `json:"max_drawdown,omitempty"`
	WinRate       float64            `json:"win_rate,omitempty"`
	ProfitFactor  float64            `json:"profit_factor,omitempty"`
	CustomMetrics map[string]float64 `json:"custom_metrics,omitempty"`
}

// PredictionRequest represents a prediction request
type PredictionRequest struct {
	ModelID   uuid.UUID              `json:"model_id"`
	Features  []float64              `json:"features"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// PredictionResult represents a prediction result
type PredictionResult struct {
	ID          uuid.UUID              `json:"id"`
	ModelID     uuid.UUID              `json:"model_id"`
	Predictions []float64              `json:"predictions"`
	Confidence  []float64              `json:"confidence,omitempty"`
	Probability []float64              `json:"probability,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
	Latency     time.Duration          `json:"latency"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// FeatureStore manages feature engineering and storage
type FeatureStore struct {
	features   map[string]*Feature
	timeSeries map[string]*TimeSeries
	mu         sync.RWMutex
}

// Feature represents a single feature
type Feature struct {
	Name      string                 `json:"name"`
	Type      FeatureType            `json:"type"`
	Value     interface{}            `json:"value"`
	Timestamp time.Time              `json:"timestamp"`
	Source    string                 `json:"source"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// FeatureType represents the type of feature
type FeatureType string

const (
	FeatureTypeNumeric     FeatureType = "numeric"
	FeatureTypeCategorical FeatureType = "categorical"
	FeatureTypeBoolean     FeatureType = "boolean"
	FeatureTypeText        FeatureType = "text"
	FeatureTypeTimestamp   FeatureType = "timestamp"
	FeatureTypeVector      FeatureType = "vector"
)

// TimeSeries represents time series data
type TimeSeries struct {
	Name     string                 `json:"name"`
	Values   []TimeSeriesPoint      `json:"values"`
	Metadata map[string]interface{} `json:"metadata"`
}

// TimeSeriesPoint represents a single point in time series
type TimeSeriesPoint struct {
	Timestamp time.Time   `json:"timestamp"`
	Value     interface{} `json:"value"`
}

// DataProcessor handles data preprocessing and feature engineering
type DataProcessor struct {
	transformers map[string]*DataTransformer
	mu           sync.RWMutex
}

// DataTransformer represents a data transformation
type DataTransformer struct {
	Name       string                 `json:"name"`
	Type       TransformerType        `json:"type"`
	Config     map[string]interface{} `json:"config"`
	Fitted     bool                   `json:"fitted"`
	Parameters map[string]interface{} `json:"parameters"`
}

// TransformerType represents the type of data transformer
type TransformerType string

const (
	TransformerTypeScaler     TransformerType = "scaler"
	TransformerTypeNormalizer TransformerType = "normalizer"
	TransformerTypeEncoder    TransformerType = "encoder"
	TransformerTypeImputer    TransformerType = "imputer"
	TransformerTypeSelector   TransformerType = "selector"
	TransformerTypeReducer    TransformerType = "reducer"
	TransformerTypeAggregator TransformerType = "aggregator"
	TransformerTypeIndicator  TransformerType = "indicator"
)

// ModelRegistry manages model versions and deployment
type ModelRegistry struct {
	models   map[uuid.UUID]*Model
	versions map[uuid.UUID]map[string]*TrainedModel
	deployed map[uuid.UUID]*TrainedModel
	mu       sync.RWMutex
}

// NewMLFramework creates a new machine learning framework
func NewMLFramework(logger *observability.Logger, config MLConfig) *MLFramework {
	// Set defaults
	if config.TrainingInterval == 0 {
		config.TrainingInterval = 1 * time.Hour
	}
	if config.PredictionCacheTime == 0 {
		config.PredictionCacheTime = 5 * time.Minute
	}
	if config.MaxConcurrentTraining == 0 {
		config.MaxConcurrentTraining = 2
	}
	if config.MaxConcurrentPrediction == 0 {
		config.MaxConcurrentPrediction = 10
	}
	if config.FeatureRetention == 0 {
		config.FeatureRetention = 30 * 24 * time.Hour // 30 days
	}
	if config.ModelRetention == 0 {
		config.ModelRetention = 90 * 24 * time.Hour // 90 days
	}

	return &MLFramework{
		logger: logger,
		config: config,
		models: make(map[string]*Model),
		featureStore: &FeatureStore{
			features:   make(map[string]*Feature),
			timeSeries: make(map[string]*TimeSeries),
		},
		dataProcessor: &DataProcessor{
			transformers: make(map[string]*DataTransformer),
		},
		modelRegistry: &ModelRegistry{
			models:   make(map[uuid.UUID]*Model),
			versions: make(map[uuid.UUID]map[string]*TrainedModel),
			deployed: make(map[uuid.UUID]*TrainedModel),
		},
		trainingQueue:   make(chan *TrainingJob, 100),
		trainedModels:   make(map[string]*TrainedModel),
		predictionCache: make(map[string]*PredictionResult),
		cacheExpiry:     make(map[string]time.Time),
		modelAccuracy:   make(map[string]float64),
		stopChan:        make(chan struct{}),
	}
}

// Start starts the ML framework
func (mlf *MLFramework) Start(ctx context.Context) error {
	mlf.mu.Lock()
	defer mlf.mu.Unlock()

	if mlf.isRunning {
		return fmt.Errorf("ML framework is already running")
	}

	mlf.logger.Info(ctx, "Starting ML framework", map[string]interface{}{
		"enable_training":           mlf.config.EnableTraining,
		"enable_prediction":         mlf.config.EnablePrediction,
		"enable_feature_store":      mlf.config.EnableFeatureStore,
		"max_concurrent_training":   mlf.config.MaxConcurrentTraining,
		"max_concurrent_prediction": mlf.config.MaxConcurrentPrediction,
		"enable_gpu":                mlf.config.EnableGPU,
	})

	mlf.isRunning = true

	// Start training workers
	if mlf.config.EnableTraining {
		for i := 0; i < mlf.config.MaxConcurrentTraining; i++ {
			mlf.wg.Add(1)
			go mlf.trainingWorker(ctx, i)
		}
	}

	// Start cache cleanup
	mlf.wg.Add(1)
	go mlf.cacheCleanup(ctx)

	// Start feature store cleanup
	if mlf.config.EnableFeatureStore {
		mlf.wg.Add(1)
		go mlf.featureStoreCleanup(ctx)
	}

	mlf.logger.Info(ctx, "ML framework started", map[string]interface{}{
		"registered_models": len(mlf.models),
	})

	return nil
}

// Stop stops the ML framework
func (mlf *MLFramework) Stop(ctx context.Context) error {
	mlf.mu.Lock()
	defer mlf.mu.Unlock()

	if !mlf.isRunning {
		return fmt.Errorf("ML framework is not running")
	}

	mlf.logger.Info(ctx, "Stopping ML framework", nil)

	close(mlf.stopChan)
	mlf.wg.Wait()

	mlf.isRunning = false

	mlf.logger.Info(ctx, "ML framework stopped", nil)

	return nil
}

// RegisterModel registers a new ML model
func (mlf *MLFramework) RegisterModel(ctx context.Context, model *Model) error {
	mlf.mu.Lock()
	defer mlf.mu.Unlock()

	if _, exists := mlf.models[model.Name]; exists {
		return fmt.Errorf("model already exists: %s", model.Name)
	}

	model.ID = uuid.New()
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()
	model.Status = ModelStatusCreated

	mlf.models[model.Name] = model
	mlf.modelRegistry.models[model.ID] = model
	mlf.modelRegistry.versions[model.ID] = make(map[string]*TrainedModel)

	mlf.logger.Info(ctx, "Model registered", map[string]interface{}{
		"model_id":   model.ID.String(),
		"model_name": model.Name,
		"model_type": string(model.Type),
		"algorithm":  string(model.Algorithm),
	})

	return nil
}

// SubmitTrainingJob submits a training job
func (mlf *MLFramework) SubmitTrainingJob(ctx context.Context, job *TrainingJob) error {
	if !mlf.config.EnableTraining {
		return fmt.Errorf("training is disabled")
	}

	job.ID = uuid.New()
	job.CreatedAt = time.Now()
	job.Status = TrainingStatusQueued

	select {
	case mlf.trainingQueue <- job:
		mlf.logger.Info(ctx, "Training job submitted", map[string]interface{}{
			"job_id":   job.ID.String(),
			"model_id": job.ModelID.String(),
			"priority": job.Priority,
		})
		return nil
	default:
		return fmt.Errorf("training queue is full")
	}
}

// Predict makes a prediction using a trained model
func (mlf *MLFramework) Predict(ctx context.Context, request *PredictionRequest) (*PredictionResult, error) {
	if !mlf.config.EnablePrediction {
		return nil, fmt.Errorf("prediction is disabled")
	}

	mlf.totalPredictions++

	// Check cache first
	cacheKey := mlf.generateCacheKey(request)
	if result := mlf.getCachedPrediction(cacheKey); result != nil {
		return result, nil
	}

	// Get deployed model
	mlf.modelRegistry.mu.RLock()
	trainedModel, exists := mlf.modelRegistry.deployed[request.ModelID]
	mlf.modelRegistry.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no deployed model found for ID: %s", request.ModelID.String())
	}

	startTime := time.Now()

	// Perform prediction (placeholder implementation)
	predictions := mlf.performPrediction(trainedModel, request.Features)

	result := &PredictionResult{
		ID:          uuid.New(),
		ModelID:     request.ModelID,
		Predictions: predictions,
		Timestamp:   time.Now(),
		Latency:     time.Since(startTime),
		Metadata:    request.Metadata,
	}

	// Cache result
	mlf.cachePrediction(cacheKey, result)

	mlf.logger.Debug(ctx, "Prediction completed", map[string]interface{}{
		"model_id":    request.ModelID.String(),
		"latency_ms":  result.Latency.Milliseconds(),
		"predictions": len(result.Predictions),
	})

	return result, nil
}

// AddFeature adds a feature to the feature store
func (mlf *MLFramework) AddFeature(ctx context.Context, feature *Feature) error {
	if !mlf.config.EnableFeatureStore {
		return fmt.Errorf("feature store is disabled")
	}

	mlf.featureStore.mu.Lock()
	defer mlf.featureStore.mu.Unlock()

	feature.Timestamp = time.Now()
	mlf.featureStore.features[feature.Name] = feature

	return nil
}

// GetFeature retrieves a feature from the feature store
func (mlf *MLFramework) GetFeature(name string) (*Feature, error) {
	if !mlf.config.EnableFeatureStore {
		return nil, fmt.Errorf("feature store is disabled")
	}

	mlf.featureStore.mu.RLock()
	defer mlf.featureStore.mu.RUnlock()

	feature, exists := mlf.featureStore.features[name]
	if !exists {
		return nil, fmt.Errorf("feature not found: %s", name)
	}

	return feature, nil
}

// GetFrameworkMetrics returns ML framework metrics
func (mlf *MLFramework) GetFrameworkMetrics() *MLFrameworkMetrics {
	mlf.mu.RLock()
	defer mlf.mu.RUnlock()

	return &MLFrameworkMetrics{
		TotalPredictions:  mlf.totalPredictions,
		TotalTraining:     mlf.totalTraining,
		RegisteredModels:  len(mlf.models),
		DeployedModels:    len(mlf.modelRegistry.deployed),
		CachedPredictions: len(mlf.predictionCache),
		IsRunning:         mlf.isRunning,
	}
}

// MLFrameworkMetrics contains ML framework metrics
type MLFrameworkMetrics struct {
	TotalPredictions  int64 `json:"total_predictions"`
	TotalTraining     int64 `json:"total_training"`
	RegisteredModels  int   `json:"registered_models"`
	DeployedModels    int   `json:"deployed_models"`
	CachedPredictions int   `json:"cached_predictions"`
	IsRunning         bool  `json:"is_running"`
}

// Private methods

// trainingWorker processes training jobs
func (mlf *MLFramework) trainingWorker(ctx context.Context, workerID int) {
	defer mlf.wg.Done()

	mlf.logger.Info(ctx, "Training worker started", map[string]interface{}{
		"worker_id": workerID,
	})

	for {
		select {
		case <-mlf.stopChan:
			return
		case job := <-mlf.trainingQueue:
			mlf.processTrainingJob(ctx, job, workerID)
		}
	}
}

// processTrainingJob processes a single training job
func (mlf *MLFramework) processTrainingJob(ctx context.Context, job *TrainingJob, workerID int) {
	mlf.totalTraining++

	job.Status = TrainingStatusRunning
	job.StartedAt = time.Now()

	mlf.logger.Info(ctx, "Starting training job", map[string]interface{}{
		"job_id":    job.ID.String(),
		"model_id":  job.ModelID.String(),
		"worker_id": workerID,
	})

	// Get model
	mlf.modelRegistry.mu.RLock()
	model, exists := mlf.modelRegistry.models[job.ModelID]
	mlf.modelRegistry.mu.RUnlock()

	if !exists {
		job.Status = TrainingStatusFailed
		job.Error = "model not found"
		job.CompletedAt = time.Now()
		return
	}

	// Perform training (placeholder implementation)
	trainedModel, err := mlf.performTraining(ctx, model, job)
	if err != nil {
		job.Status = TrainingStatusFailed
		job.Error = err.Error()
		job.CompletedAt = time.Now()

		mlf.logger.Error(ctx, "Training job failed", err, map[string]interface{}{
			"job_id":   job.ID.String(),
			"model_id": job.ModelID.String(),
		})
		return
	}

	// Store trained model
	mlf.modelRegistry.mu.Lock()
	if mlf.modelRegistry.versions[job.ModelID] == nil {
		mlf.modelRegistry.versions[job.ModelID] = make(map[string]*TrainedModel)
	}
	mlf.modelRegistry.versions[job.ModelID][trainedModel.Version] = trainedModel
	mlf.modelRegistry.mu.Unlock()

	job.Status = TrainingStatusCompleted
	job.CompletedAt = time.Now()
	job.Progress = 100.0

	mlf.logger.Info(ctx, "Training job completed", map[string]interface{}{
		"job_id":        job.ID.String(),
		"model_id":      job.ModelID.String(),
		"training_time": time.Since(job.StartedAt),
		"accuracy":      trainedModel.Metrics.Accuracy,
	})
}

// performTraining performs the actual model training (placeholder)
func (mlf *MLFramework) performTraining(ctx context.Context, model *Model, job *TrainingJob) (*TrainedModel, error) {
	// This is a placeholder implementation
	// In a real system, this would integrate with ML libraries like TensorFlow, PyTorch, etc.

	startTime := time.Now()

	// Simulate training time
	time.Sleep(100 * time.Millisecond)

	// Create mock metrics based on algorithm
	metrics := &ModelMetrics{}

	switch model.Type {
	case ModelTypeRegression:
		metrics.MSE = 0.05
		metrics.RMSE = math.Sqrt(metrics.MSE)
		metrics.MAE = 0.03
		metrics.R2Score = 0.85
	case ModelTypeClassification:
		metrics.Accuracy = 0.88
		metrics.Precision = 0.86
		metrics.Recall = 0.84
		metrics.F1Score = 0.85
		metrics.AUC = 0.92
	case ModelTypeTimeSeries:
		metrics.MSE = 0.02
		metrics.RMSE = math.Sqrt(metrics.MSE)
		metrics.SharpeRatio = 1.5
		metrics.MaxDrawdown = 0.15
		metrics.WinRate = 0.65
	}

	trainedModel := &TrainedModel{
		ModelID:      model.ID,
		Version:      fmt.Sprintf("v%d", time.Now().Unix()),
		TrainingData: job.Dataset,
		Metrics:      metrics,
		TrainedAt:    time.Now(),
		TrainingTime: time.Since(startTime),
		Status:       ModelStatusTrained,
		Metadata:     make(map[string]interface{}),
	}

	return trainedModel, nil
}

// performPrediction performs model prediction (placeholder)
func (mlf *MLFramework) performPrediction(trainedModel *TrainedModel, features []float64) []float64 {
	// This is a placeholder implementation
	// In a real system, this would use the actual trained model

	predictions := make([]float64, 1)

	// Simple mock prediction based on features
	if len(features) > 0 {
		sum := 0.0
		for _, feature := range features {
			sum += feature
		}
		predictions[0] = sum / float64(len(features))
	}

	return predictions
}

// cacheCleanup periodically cleans up expired cache entries
func (mlf *MLFramework) cacheCleanup(ctx context.Context) {
	defer mlf.wg.Done()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-mlf.stopChan:
			return
		case <-ticker.C:
			mlf.cleanupExpiredCache()
		}
	}
}

// cleanupExpiredCache removes expired cache entries
func (mlf *MLFramework) cleanupExpiredCache() {
	mlf.mu.Lock()
	defer mlf.mu.Unlock()

	now := time.Now()
	for key, expiry := range mlf.cacheExpiry {
		if now.After(expiry) {
			delete(mlf.predictionCache, key)
			delete(mlf.cacheExpiry, key)
		}
	}
}

// featureStoreCleanup periodically cleans up old features
func (mlf *MLFramework) featureStoreCleanup(ctx context.Context) {
	defer mlf.wg.Done()

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-mlf.stopChan:
			return
		case <-ticker.C:
			mlf.cleanupOldFeatures()
		}
	}
}

// cleanupOldFeatures removes old features from the feature store
func (mlf *MLFramework) cleanupOldFeatures() {
	mlf.featureStore.mu.Lock()
	defer mlf.featureStore.mu.Unlock()

	cutoff := time.Now().Add(-mlf.config.FeatureRetention)
	for name, feature := range mlf.featureStore.features {
		if feature.Timestamp.Before(cutoff) {
			delete(mlf.featureStore.features, name)
		}
	}
}

// generateCacheKey generates a cache key for prediction requests
func (mlf *MLFramework) generateCacheKey(request *PredictionRequest) string {
	return fmt.Sprintf("%s_%d", request.ModelID.String(), len(request.Features))
}

// getCachedPrediction retrieves a cached prediction result
func (mlf *MLFramework) getCachedPrediction(key string) *PredictionResult {
	mlf.mu.RLock()
	defer mlf.mu.RUnlock()

	// Check if cached and not expired
	if expiry, exists := mlf.cacheExpiry[key]; exists {
		if time.Now().Before(expiry) {
			if result, exists := mlf.predictionCache[key]; exists {
				return result
			}
		}
	}

	return nil
}

// cachePrediction caches a prediction result
func (mlf *MLFramework) cachePrediction(key string, result *PredictionResult) {
	mlf.mu.Lock()
	defer mlf.mu.Unlock()

	mlf.predictionCache[key] = result
	mlf.cacheExpiry[key] = time.Now().Add(mlf.config.PredictionCacheTime)
}
