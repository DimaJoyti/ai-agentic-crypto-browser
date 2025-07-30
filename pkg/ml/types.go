package ml

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

// ModelType represents the type of ML model
type ModelType string

const (
	ModelTypeClassification   ModelType = "classification"
	ModelTypeRegression       ModelType = "regression"
	ModelTypeTimeSeries       ModelType = "time_series"
	ModelTypeAnomalyDetection ModelType = "anomaly_detection"
	ModelTypeNLP              ModelType = "nlp"
	ModelTypeDeepLearning     ModelType = "deep_learning"
)

// ModelStatus represents the status of a model
type ModelStatus string

const (
	ModelStatusTraining   ModelStatus = "training"
	ModelStatusReady      ModelStatus = "ready"
	ModelStatusUpdating   ModelStatus = "updating"
	ModelStatusError      ModelStatus = "error"
	ModelStatusDeprecated ModelStatus = "deprecated"
)

// Model represents a machine learning model
type Model interface {
	// Predict makes a prediction using the model
	Predict(ctx context.Context, features map[string]interface{}) (*Prediction, error)

	// Train trains the model with new data
	Train(ctx context.Context, data TrainingData) error

	// Evaluate evaluates the model performance
	Evaluate(ctx context.Context, testData TrainingData) (*ModelMetrics, error)

	// GetInfo returns model information
	GetInfo() *ModelInfo

	// IsReady returns true if the model is ready for predictions
	IsReady() bool

	// UpdateWeights updates model weights (for online learning)
	UpdateWeights(ctx context.Context, feedback *PredictionFeedback) error
}

// ModelInfo contains metadata about a model
type ModelInfo struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Version      string                 `json:"version"`
	Type         ModelType              `json:"type"`
	Status       ModelStatus            `json:"status"`
	Features     []string               `json:"features"`
	OutputSchema map[string]interface{} `json:"output_schema"`
	Accuracy     float64                `json:"accuracy"`
	Precision    float64                `json:"precision"`
	Recall       float64                `json:"recall"`
	F1Score      float64                `json:"f1_score"`
	LastTrained  time.Time              `json:"last_trained"`
	LastUpdated  time.Time              `json:"last_updated"`
	TrainingSize int                    `json:"training_size"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// Prediction represents a model prediction
type Prediction struct {
	Value       interface{}            `json:"value"`
	Confidence  float64                `json:"confidence"`
	Probability map[string]float64     `json:"probability,omitempty"`
	Features    map[string]interface{} `json:"features"`
	ModelID     string                 `json:"model_id"`
	Timestamp   time.Time              `json:"timestamp"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// TrainingData represents training data for ML models
type TrainingData struct {
	Features []map[string]interface{} `json:"features"`
	Labels   []interface{}            `json:"labels"`
	Weights  []float64                `json:"weights,omitempty"`
	Metadata map[string]interface{}   `json:"metadata"`
}

// PredictionFeedback represents feedback on a prediction
type PredictionFeedback struct {
	PredictionID string                 `json:"prediction_id"`
	ActualValue  interface{}            `json:"actual_value"`
	Correct      bool                   `json:"correct"`
	Confidence   float64                `json:"confidence"`
	Timestamp    time.Time              `json:"timestamp"`
	UserID       string                 `json:"user_id,omitempty"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// ModelMetrics represents model performance metrics
type ModelMetrics struct {
	Accuracy             float64                `json:"accuracy"`
	Precision            float64                `json:"precision"`
	Recall               float64                `json:"recall"`
	F1Score              float64                `json:"f1_score"`
	AUC                  float64                `json:"auc,omitempty"`
	MAE                  float64                `json:"mae,omitempty"`  // Mean Absolute Error
	MSE                  float64                `json:"mse,omitempty"`  // Mean Squared Error
	RMSE                 float64                `json:"rmse,omitempty"` // Root Mean Squared Error
	ConfusionMatrix      [][]int                `json:"confusion_matrix,omitempty"`
	ClassificationReport map[string]interface{} `json:"classification_report,omitempty"`
	FeatureImportance    map[string]float64     `json:"feature_importance,omitempty"`
	TestSize             int                    `json:"test_size"`
	EvaluatedAt          time.Time              `json:"evaluated_at"`
}

// TimeSeriesData represents time series data for prediction
type TimeSeriesData struct {
	Timestamp time.Time              `json:"timestamp"`
	Values    map[string]interface{} `json:"values"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// PriceData represents cryptocurrency price data
type PriceData struct {
	Symbol    string                 `json:"symbol"`
	Timestamp time.Time              `json:"timestamp"`
	Open      decimal.Decimal        `json:"open"`
	High      decimal.Decimal        `json:"high"`
	Low       decimal.Decimal        `json:"low"`
	Close     decimal.Decimal        `json:"close"`
	Volume    decimal.Decimal        `json:"volume"`
	MarketCap decimal.Decimal        `json:"market_cap,omitempty"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// SentimentData represents sentiment analysis data
type SentimentData struct {
	Text       string                 `json:"text"`
	Source     string                 `json:"source"` // twitter, reddit, news, etc.
	Timestamp  time.Time              `json:"timestamp"`
	Symbol     string                 `json:"symbol,omitempty"`
	Sentiment  float64                `json:"sentiment"` // -1.0 to 1.0
	Confidence float64                `json:"confidence"`
	Emotions   map[string]float64     `json:"emotions,omitempty"`
	Keywords   []string               `json:"keywords,omitempty"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// MarketData represents comprehensive market data
type MarketData struct {
	Timestamp      time.Time              `json:"timestamp"`
	Prices         []PriceData            `json:"prices"`
	Sentiment      []SentimentData        `json:"sentiment"`
	Volume         decimal.Decimal        `json:"volume"`
	MarketCap      decimal.Decimal        `json:"market_cap"`
	Dominance      map[string]float64     `json:"dominance"`
	FearGreedIndex int                    `json:"fear_greed_index"`
	Volatility     float64                `json:"volatility"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// ModelConfig represents configuration for ML models
type ModelConfig struct {
	ModelType        ModelType              `json:"model_type"`
	HyperParameters  map[string]interface{} `json:"hyperparameters"`
	TrainingConfig   TrainingConfig         `json:"training_config"`
	ValidationConfig ValidationConfig       `json:"validation_config"`
	DeploymentConfig DeploymentConfig       `json:"deployment_config"`
}

// TrainingConfig represents training configuration
type TrainingConfig struct {
	BatchSize       int           `json:"batch_size"`
	Epochs          int           `json:"epochs"`
	LearningRate    float64       `json:"learning_rate"`
	ValidationSplit float64       `json:"validation_split"`
	EarlyStopping   bool          `json:"early_stopping"`
	Patience        int           `json:"patience"`
	Timeout         time.Duration `json:"timeout"`
	SaveCheckpoints bool          `json:"save_checkpoints"`
}

// ValidationConfig represents validation configuration
type ValidationConfig struct {
	Method     string  `json:"method"` // cross_validation, holdout, time_series_split
	Folds      int     `json:"folds,omitempty"`
	TestSize   float64 `json:"test_size,omitempty"`
	Shuffle    bool    `json:"shuffle"`
	RandomSeed int     `json:"random_seed"`
}

// DeploymentConfig represents deployment configuration
type DeploymentConfig struct {
	Environment      string            `json:"environment"` // development, staging, production
	ScalingPolicy    string            `json:"scaling_policy"`
	ResourceLimits   map[string]string `json:"resource_limits"`
	HealthCheck      HealthCheckConfig `json:"health_check"`
	MonitoringConfig MonitoringConfig  `json:"monitoring"`
}

// HealthCheckConfig represents health check configuration
type HealthCheckConfig struct {
	Enabled          bool          `json:"enabled"`
	Interval         time.Duration `json:"interval"`
	Timeout          time.Duration `json:"timeout"`
	FailureThreshold int           `json:"failure_threshold"`
	SuccessThreshold int           `json:"success_threshold"`
}

// MonitoringConfig represents monitoring configuration
type MonitoringConfig struct {
	MetricsEnabled    bool               `json:"metrics_enabled"`
	LoggingEnabled    bool               `json:"logging_enabled"`
	TracingEnabled    bool               `json:"tracing_enabled"`
	AlertingEnabled   bool               `json:"alerting_enabled"`
	MetricsInterval   time.Duration      `json:"metrics_interval"`
	PerformanceAlerts map[string]float64 `json:"performance_alerts"`
}
