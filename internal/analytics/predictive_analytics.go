package analytics

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

// PredictiveAnalyzer provides predictive analytics capabilities
type PredictiveAnalyzer struct {
	logger          *observability.Logger
	config          *AnalyticsConfig
	models          map[string]*PredictiveModel
	predictions     map[string]*Prediction
	trainingData    map[string][]DataPoint
	forecastHorizon time.Duration
	updateInterval  time.Duration
	mu              sync.RWMutex
}

// PredictiveModel represents a predictive model
type PredictiveModel struct {
	ModelID        string                 `json:"model_id"`
	MetricName     string                 `json:"metric_name"`
	ModelType      PredictiveModelType    `json:"model_type"`
	Algorithm      string                 `json:"algorithm"`
	Parameters     map[string]float64     `json:"parameters"`
	TrainingData   []DataPoint            `json:"training_data"`
	ValidationData []DataPoint            `json:"validation_data"`
	Accuracy       float64                `json:"accuracy"`
	RMSE           float64                `json:"rmse"`
	MAE            float64                `json:"mae"`
	R2Score        float64                `json:"r2_score"`
	LastTrained    time.Time              `json:"last_trained"`
	LastUpdated    time.Time              `json:"last_updated"`
	Status         ModelStatus            `json:"status"`
	Metadata       map[string]interface{} `json:"metadata"`
	mu             sync.RWMutex           `json:"-"`
}

// PredictiveModelType defines types of predictive models
type PredictiveModelType string

const (
	ModelTypeLinearRegression     PredictiveModelType = "linear_regression"
	ModelTypePolynomialRegression PredictiveModelType = "polynomial_regression"
	ModelTypeMovingAverage        PredictiveModelType = "moving_average"
	ModelTypeExponentialSmoothing PredictiveModelType = "exponential_smoothing"
	ModelTypeARIMA                PredictiveModelType = "arima"
	ModelTypeSeasonal             PredictiveModelType = "seasonal"
	ModelTypeNeuralNetwork        PredictiveModelType = "neural_network"
	ModelTypeEnsemble             PredictiveModelType = "ensemble"
)

// ModelStatus defines model status
type ModelStatus string

const (
	ModelStatusTraining   ModelStatus = "training"
	ModelStatusTrained    ModelStatus = "trained"
	ModelStatusActive     ModelStatus = "active"
	ModelStatusDeprecated ModelStatus = "deprecated"
	ModelStatusError      ModelStatus = "error"
)

// Prediction represents a prediction result
type Prediction struct {
	PredictionID    string                 `json:"prediction_id"`
	ModelID         string                 `json:"model_id"`
	MetricName      string                 `json:"metric_name"`
	PredictedValue  float64                `json:"predicted_value"`
	ConfidenceLevel float64                `json:"confidence_level"`
	PredictionTime  time.Time              `json:"prediction_time"`
	TargetTime      time.Time              `json:"target_time"`
	Horizon         time.Duration          `json:"horizon"`
	UpperBound      float64                `json:"upper_bound"`
	LowerBound      float64                `json:"lower_bound"`
	Trend           TrendDirection         `json:"trend"`
	Seasonality     *SeasonalityInfo       `json:"seasonality,omitempty"`
	Context         map[string]interface{} `json:"context"`
	Accuracy        float64                `json:"accuracy,omitempty"`
	ActualValue     *float64               `json:"actual_value,omitempty"`
	Error           *float64               `json:"error,omitempty"`
}

// TrendDirection defines trend directions
type TrendDirection string

const (
	TrendDirectionUp       TrendDirection = "up"
	TrendDirectionDown     TrendDirection = "down"
	TrendDirectionStable   TrendDirection = "stable"
	TrendDirectionVolatile TrendDirection = "volatile"
)

// SeasonalityInfo contains seasonality information
type SeasonalityInfo struct {
	Period    time.Duration `json:"period"`
	Amplitude float64       `json:"amplitude"`
	Phase     float64       `json:"phase"`
	Strength  float64       `json:"strength"`
}

// ForecastRequest represents a forecast request
type ForecastRequest struct {
	MetricName string                 `json:"metric_name"`
	Horizon    time.Duration          `json:"horizon"`
	Intervals  int                    `json:"intervals"`
	ModelType  *PredictiveModelType   `json:"model_type,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// ForecastResult represents a forecast result
type ForecastResult struct {
	MetricName  string                 `json:"metric_name"`
	Predictions []Prediction           `json:"predictions"`
	ModelUsed   string                 `json:"model_used"`
	Confidence  float64                `json:"confidence"`
	GeneratedAt time.Time              `json:"generated_at"`
	ValidUntil  time.Time              `json:"valid_until"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// NewPredictiveAnalyzer creates a new predictive analyzer
func NewPredictiveAnalyzer(logger *observability.Logger, config *AnalyticsConfig) *PredictiveAnalyzer {
	return &PredictiveAnalyzer{
		logger:          logger,
		config:          config,
		models:          make(map[string]*PredictiveModel),
		predictions:     make(map[string]*Prediction),
		trainingData:    make(map[string][]DataPoint),
		forecastHorizon: config.PredictionHorizon,
		updateInterval:  1 * time.Hour,
	}
}

// Start starts the predictive analyzer
func (pa *PredictiveAnalyzer) Start(ctx context.Context) error {
	pa.logger.Info(ctx, "Starting predictive analyzer", map[string]interface{}{
		"forecast_horizon": pa.forecastHorizon,
		"update_interval":  pa.updateInterval,
	})

	// Initialize default models
	pa.initializeDefaultModels()

	// Start background processes
	go pa.updateModels(ctx)
	go pa.generatePredictions(ctx)
	go pa.validatePredictions(ctx)

	return nil
}

// AddTrainingData adds training data for a metric
func (pa *PredictiveAnalyzer) AddTrainingData(metricName string, dataPoint DataPoint) {
	pa.mu.Lock()
	defer pa.mu.Unlock()

	if _, exists := pa.trainingData[metricName]; !exists {
		pa.trainingData[metricName] = make([]DataPoint, 0)
	}

	pa.trainingData[metricName] = append(pa.trainingData[metricName], dataPoint)

	// Maintain training data window (keep last 1000 points)
	maxPoints := 1000
	if len(pa.trainingData[metricName]) > maxPoints {
		pa.trainingData[metricName] = pa.trainingData[metricName][len(pa.trainingData[metricName])-maxPoints:]
	}

	// Trigger model retraining if enough new data
	if len(pa.trainingData[metricName])%100 == 0 {
		go pa.retrainModel(metricName)
	}
}

// CreateModel creates a new predictive model
func (pa *PredictiveAnalyzer) CreateModel(metricName string, modelType PredictiveModelType, parameters map[string]float64) (*PredictiveModel, error) {
	pa.mu.Lock()
	defer pa.mu.Unlock()

	modelID := uuid.New().String()
	model := &PredictiveModel{
		ModelID:     modelID,
		MetricName:  metricName,
		ModelType:   modelType,
		Algorithm:   string(modelType),
		Parameters:  parameters,
		Status:      ModelStatusTraining,
		LastUpdated: time.Now(),
		Metadata:    make(map[string]interface{}),
	}

	pa.models[modelID] = model

	// Train the model
	go pa.trainModel(model)

	pa.logger.Info(context.Background(), "Predictive model created", map[string]interface{}{
		"model_id":    modelID,
		"metric_name": metricName,
		"model_type":  modelType,
	})

	return model, nil
}

// GenerateForecast generates a forecast for a metric
func (pa *PredictiveAnalyzer) GenerateForecast(ctx context.Context, request *ForecastRequest) (*ForecastResult, error) {
	// Find the best model for the metric
	var bestModel *PredictiveModel
	bestAccuracy := 0.0

	pa.mu.RLock()
	for _, model := range pa.models {
		if model.MetricName == request.MetricName && model.Status == ModelStatusActive {
			if request.ModelType == nil || model.ModelType == *request.ModelType {
				if model.Accuracy > bestAccuracy {
					bestModel = model
					bestAccuracy = model.Accuracy
				}
			}
		}
	}
	pa.mu.RUnlock()

	if bestModel == nil {
		return nil, fmt.Errorf("no suitable model found for metric: %s", request.MetricName)
	}

	// Generate predictions
	predictions := make([]Prediction, 0, request.Intervals)
	currentTime := time.Now()
	intervalDuration := request.Horizon / time.Duration(request.Intervals)

	for i := 0; i < request.Intervals; i++ {
		targetTime := currentTime.Add(time.Duration(i+1) * intervalDuration)
		prediction := pa.generateSinglePrediction(bestModel, targetTime)
		predictions = append(predictions, *prediction)
	}

	result := &ForecastResult{
		MetricName:  request.MetricName,
		Predictions: predictions,
		ModelUsed:   bestModel.ModelID,
		Confidence:  bestModel.Accuracy,
		GeneratedAt: time.Now(),
		ValidUntil:  time.Now().Add(request.Horizon),
		Metadata: map[string]interface{}{
			"model_type": bestModel.ModelType,
			"algorithm":  bestModel.Algorithm,
			"rmse":       bestModel.RMSE,
			"mae":        bestModel.MAE,
			"r2_score":   bestModel.R2Score,
		},
	}

	pa.logger.Info(ctx, "Forecast generated", map[string]interface{}{
		"metric_name":      request.MetricName,
		"model_id":         bestModel.ModelID,
		"prediction_count": len(predictions),
		"confidence":       bestModel.Accuracy,
	})

	return result, nil
}

// generateSinglePrediction generates a single prediction
func (pa *PredictiveAnalyzer) generateSinglePrediction(model *PredictiveModel, targetTime time.Time) *Prediction {
	// Get training data for the model
	pa.mu.RLock()
	trainingData, exists := pa.trainingData[model.MetricName]
	pa.mu.RUnlock()

	if !exists || len(trainingData) == 0 {
		return &Prediction{
			PredictionID:    uuid.New().String(),
			ModelID:         model.ModelID,
			MetricName:      model.MetricName,
			PredictedValue:  0,
			ConfidenceLevel: 0,
			PredictionTime:  time.Now(),
			TargetTime:      targetTime,
			Horizon:         targetTime.Sub(time.Now()),
			Trend:           TrendDirectionStable,
		}
	}

	// Generate prediction based on model type
	var predictedValue float64
	var confidence float64
	var upperBound, lowerBound float64
	var trend TrendDirection

	switch model.ModelType {
	case ModelTypeMovingAverage:
		predictedValue, confidence = pa.predictMovingAverage(trainingData, model.Parameters)
	case ModelTypeLinearRegression:
		predictedValue, confidence = pa.predictLinearRegression(trainingData, model.Parameters)
	case ModelTypeExponentialSmoothing:
		predictedValue, confidence = pa.predictExponentialSmoothing(trainingData, model.Parameters)
	default:
		predictedValue, confidence = pa.predictMovingAverage(trainingData, model.Parameters)
	}

	// Calculate bounds and trend
	variance := model.Parameters["variance"]
	if variance > 0 {
		stdDev := math.Sqrt(variance)
		upperBound = predictedValue + 2*stdDev
		lowerBound = predictedValue - 2*stdDev
	} else {
		upperBound = predictedValue * 1.1
		lowerBound = predictedValue * 0.9
	}

	trend = pa.calculateTrend(trainingData)

	prediction := &Prediction{
		PredictionID:    uuid.New().String(),
		ModelID:         model.ModelID,
		MetricName:      model.MetricName,
		PredictedValue:  predictedValue,
		ConfidenceLevel: confidence,
		PredictionTime:  time.Now(),
		TargetTime:      targetTime,
		Horizon:         targetTime.Sub(time.Now()),
		UpperBound:      upperBound,
		LowerBound:      lowerBound,
		Trend:           trend,
		Context: map[string]interface{}{
			"model_type":      model.ModelType,
			"training_points": len(trainingData),
		},
	}

	// Store prediction
	pa.mu.Lock()
	pa.predictions[prediction.PredictionID] = prediction
	pa.mu.Unlock()

	return prediction
}

// predictMovingAverage predicts using moving average
func (pa *PredictiveAnalyzer) predictMovingAverage(data []DataPoint, parameters map[string]float64) (float64, float64) {
	if len(data) == 0 {
		return 0, 0
	}

	window := int(parameters["window"])
	if window <= 0 || window > len(data) {
		window = min(10, len(data))
	}

	// Calculate moving average
	sum := 0.0
	recentData := data[len(data)-window:]
	for _, point := range recentData {
		sum += point.Value
	}

	average := sum / float64(len(recentData))
	confidence := 0.7 // Base confidence for moving average

	return average, confidence
}

// predictLinearRegression predicts using linear regression
func (pa *PredictiveAnalyzer) predictLinearRegression(data []DataPoint, parameters map[string]float64) (float64, float64) {
	if len(data) < 2 {
		return 0, 0
	}

	// Simple linear regression implementation
	n := float64(len(data))
	sumX, sumY, sumXY, sumX2 := 0.0, 0.0, 0.0, 0.0

	for i, point := range data {
		x := float64(i)
		y := point.Value
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	// Calculate slope and intercept
	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	intercept := (sumY - slope*sumX) / n

	// Predict next value
	nextX := float64(len(data))
	predictedValue := slope*nextX + intercept

	// Calculate R-squared for confidence
	meanY := sumY / n
	ssRes, ssTot := 0.0, 0.0
	for i, point := range data {
		predicted := slope*float64(i) + intercept
		ssRes += math.Pow(point.Value-predicted, 2)
		ssTot += math.Pow(point.Value-meanY, 2)
	}

	r2 := 1 - (ssRes / ssTot)
	confidence := math.Max(0, math.Min(1, r2))

	return predictedValue, confidence
}

// predictExponentialSmoothing predicts using exponential smoothing
func (pa *PredictiveAnalyzer) predictExponentialSmoothing(data []DataPoint, parameters map[string]float64) (float64, float64) {
	if len(data) == 0 {
		return 0, 0
	}

	alpha := parameters["alpha"]
	if alpha <= 0 || alpha > 1 {
		alpha = 0.3 // Default smoothing factor
	}

	// Initialize with first value
	smoothed := data[0].Value

	// Apply exponential smoothing
	for i := 1; i < len(data); i++ {
		smoothed = alpha*data[i].Value + (1-alpha)*smoothed
	}

	confidence := 0.8 // Higher confidence for exponential smoothing

	return smoothed, confidence
}

// calculateTrend calculates the trend direction
func (pa *PredictiveAnalyzer) calculateTrend(data []DataPoint) TrendDirection {
	if len(data) < 2 {
		return TrendDirectionStable
	}

	// Calculate trend over last 10 points or all data if less
	window := min(10, len(data))
	recentData := data[len(data)-window:]

	if len(recentData) < 2 {
		return TrendDirectionStable
	}

	// Simple trend calculation
	first := recentData[0].Value
	last := recentData[len(recentData)-1].Value
	change := (last - first) / first

	if math.Abs(change) < 0.05 { // Less than 5% change
		return TrendDirectionStable
	} else if change > 0 {
		return TrendDirectionUp
	} else {
		return TrendDirectionDown
	}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// initializeDefaultModels initializes default predictive models
func (pa *PredictiveAnalyzer) initializeDefaultModels() {
	defaultModels := []struct {
		metricName string
		modelType  PredictiveModelType
		parameters map[string]float64
	}{
		{"cpu_usage", ModelTypeMovingAverage, map[string]float64{"window": 10}},
		{"memory_usage", ModelTypeLinearRegression, map[string]float64{}},
		{"response_time", ModelTypeExponentialSmoothing, map[string]float64{"alpha": 0.3}},
		{"request_count", ModelTypeMovingAverage, map[string]float64{"window": 15}},
		{"trading_volume", ModelTypeLinearRegression, map[string]float64{}},
		{"price_change", ModelTypeExponentialSmoothing, map[string]float64{"alpha": 0.2}},
	}

	for _, model := range defaultModels {
		pa.CreateModel(model.metricName, model.modelType, model.parameters)
	}
}

// updateModels updates models periodically
func (pa *PredictiveAnalyzer) updateModels(ctx context.Context) {
	ticker := time.NewTicker(pa.updateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			pa.updateAllModels()
		}
	}
}

// updateAllModels updates all models
func (pa *PredictiveAnalyzer) updateAllModels() {
	pa.mu.RLock()
	models := make([]*PredictiveModel, 0, len(pa.models))
	for _, model := range pa.models {
		models = append(models, model)
	}
	pa.mu.RUnlock()

	for _, model := range models {
		if time.Since(model.LastTrained) > 24*time.Hour {
			go pa.retrainModel(model.MetricName)
		}
	}
}

// retrainModel retrains a model for a specific metric
func (pa *PredictiveAnalyzer) retrainModel(metricName string) {
	pa.mu.RLock()
	trainingData, exists := pa.trainingData[metricName]
	if !exists || len(trainingData) < 10 {
		pa.mu.RUnlock()
		return
	}

	// Find models for this metric
	var modelsToRetrain []*PredictiveModel
	for _, model := range pa.models {
		if model.MetricName == metricName {
			modelsToRetrain = append(modelsToRetrain, model)
		}
	}
	pa.mu.RUnlock()

	for _, model := range modelsToRetrain {
		pa.trainModel(model)
	}
}

// trainModel trains a predictive model
func (pa *PredictiveAnalyzer) trainModel(model *PredictiveModel) {
	model.mu.Lock()
	defer model.mu.Unlock()

	pa.logger.Info(context.Background(), "Training predictive model", map[string]interface{}{
		"model_id":    model.ModelID,
		"metric_name": model.MetricName,
		"model_type":  model.ModelType,
	})

	model.Status = ModelStatusTraining

	// Get training data
	pa.mu.RLock()
	trainingData, exists := pa.trainingData[model.MetricName]
	pa.mu.RUnlock()

	if !exists || len(trainingData) < 10 {
		model.Status = ModelStatusError
		return
	}

	// Split data into training and validation
	splitIndex := int(0.8 * float64(len(trainingData)))
	model.TrainingData = trainingData[:splitIndex]
	model.ValidationData = trainingData[splitIndex:]

	// Train based on model type
	switch model.ModelType {
	case ModelTypeMovingAverage:
		pa.trainMovingAverage(model)
	case ModelTypeLinearRegression:
		pa.trainLinearRegression(model)
	case ModelTypeExponentialSmoothing:
		pa.trainExponentialSmoothing(model)
	default:
		pa.trainMovingAverage(model)
	}

	// Validate model
	pa.validateModel(model)

	model.LastTrained = time.Now()
	model.LastUpdated = time.Now()
	model.Status = ModelStatusActive

	pa.logger.Info(context.Background(), "Model training completed", map[string]interface{}{
		"model_id":    model.ModelID,
		"metric_name": model.MetricName,
		"accuracy":    model.Accuracy,
		"rmse":        model.RMSE,
		"mae":         model.MAE,
	})
}

// trainMovingAverage trains a moving average model
func (pa *PredictiveAnalyzer) trainMovingAverage(model *PredictiveModel) {
	// Optimize window size
	bestWindow := 10
	bestError := math.Inf(1)

	for window := 5; window <= 20; window++ {
		error := pa.calculateMovingAverageError(model.TrainingData, window)
		if error < bestError {
			bestError = error
			bestWindow = window
		}
	}

	model.Parameters["window"] = float64(bestWindow)
	model.Parameters["error"] = bestError
}

// trainLinearRegression trains a linear regression model
func (pa *PredictiveAnalyzer) trainLinearRegression(model *PredictiveModel) {
	data := model.TrainingData
	if len(data) < 2 {
		return
	}

	n := float64(len(data))
	sumX, sumY, sumXY, sumX2 := 0.0, 0.0, 0.0, 0.0

	for i, point := range data {
		x := float64(i)
		y := point.Value
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	// Calculate slope and intercept
	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	intercept := (sumY - slope*sumX) / n

	model.Parameters["slope"] = slope
	model.Parameters["intercept"] = intercept

	// Calculate variance
	meanY := sumY / n
	variance := 0.0
	for _, point := range data {
		variance += math.Pow(point.Value-meanY, 2)
	}
	variance /= n
	model.Parameters["variance"] = variance
}

// trainExponentialSmoothing trains an exponential smoothing model
func (pa *PredictiveAnalyzer) trainExponentialSmoothing(model *PredictiveModel) {
	// Optimize alpha parameter
	bestAlpha := 0.3
	bestError := math.Inf(1)

	for alpha := 0.1; alpha <= 0.9; alpha += 0.1 {
		error := pa.calculateExponentialSmoothingError(model.TrainingData, alpha)
		if error < bestError {
			bestError = error
			bestAlpha = alpha
		}
	}

	model.Parameters["alpha"] = bestAlpha
	model.Parameters["error"] = bestError
}

// validateModel validates a trained model
func (pa *PredictiveAnalyzer) validateModel(model *PredictiveModel) {
	if len(model.ValidationData) == 0 {
		model.Accuracy = 0.5 // Default accuracy
		return
	}

	predictions := make([]float64, len(model.ValidationData))
	actuals := make([]float64, len(model.ValidationData))

	// Generate predictions for validation data
	for i, point := range model.ValidationData {
		// Use training data up to this point
		trainingSubset := append(model.TrainingData, model.ValidationData[:i]...)

		var predicted float64
		switch model.ModelType {
		case ModelTypeMovingAverage:
			predicted, _ = pa.predictMovingAverage(trainingSubset, model.Parameters)
		case ModelTypeLinearRegression:
			predicted, _ = pa.predictLinearRegression(trainingSubset, model.Parameters)
		case ModelTypeExponentialSmoothing:
			predicted, _ = pa.predictExponentialSmoothing(trainingSubset, model.Parameters)
		default:
			predicted, _ = pa.predictMovingAverage(trainingSubset, model.Parameters)
		}

		predictions[i] = predicted
		actuals[i] = point.Value
	}

	// Calculate metrics
	model.RMSE = pa.calculateRMSE(predictions, actuals)
	model.MAE = pa.calculateMAE(predictions, actuals)
	model.R2Score = pa.calculateR2(predictions, actuals)
	model.Accuracy = math.Max(0, 1-model.RMSE/pa.calculateMean(actuals))
}

// calculateMovingAverageError calculates error for moving average
func (pa *PredictiveAnalyzer) calculateMovingAverageError(data []DataPoint, window int) float64 {
	if len(data) <= window {
		return math.Inf(1)
	}

	totalError := 0.0
	count := 0

	for i := window; i < len(data); i++ {
		// Calculate moving average
		sum := 0.0
		for j := i - window; j < i; j++ {
			sum += data[j].Value
		}
		predicted := sum / float64(window)
		actual := data[i].Value
		totalError += math.Pow(predicted-actual, 2)
		count++
	}

	return math.Sqrt(totalError / float64(count))
}

// calculateExponentialSmoothingError calculates error for exponential smoothing
func (pa *PredictiveAnalyzer) calculateExponentialSmoothingError(data []DataPoint, alpha float64) float64 {
	if len(data) < 2 {
		return math.Inf(1)
	}

	smoothed := data[0].Value
	totalError := 0.0

	for i := 1; i < len(data); i++ {
		predicted := smoothed
		actual := data[i].Value
		totalError += math.Pow(predicted-actual, 2)
		smoothed = alpha*actual + (1-alpha)*smoothed
	}

	return math.Sqrt(totalError / float64(len(data)-1))
}

// calculateRMSE calculates Root Mean Square Error
func (pa *PredictiveAnalyzer) calculateRMSE(predictions, actuals []float64) float64 {
	if len(predictions) != len(actuals) || len(predictions) == 0 {
		return 0
	}

	sumSquaredError := 0.0
	for i := range predictions {
		error := predictions[i] - actuals[i]
		sumSquaredError += error * error
	}

	return math.Sqrt(sumSquaredError / float64(len(predictions)))
}

// calculateMAE calculates Mean Absolute Error
func (pa *PredictiveAnalyzer) calculateMAE(predictions, actuals []float64) float64 {
	if len(predictions) != len(actuals) || len(predictions) == 0 {
		return 0
	}

	sumAbsoluteError := 0.0
	for i := range predictions {
		error := math.Abs(predictions[i] - actuals[i])
		sumAbsoluteError += error
	}

	return sumAbsoluteError / float64(len(predictions))
}

// calculateR2 calculates R-squared
func (pa *PredictiveAnalyzer) calculateR2(predictions, actuals []float64) float64 {
	if len(predictions) != len(actuals) || len(predictions) == 0 {
		return 0
	}

	meanActual := pa.calculateMean(actuals)
	ssRes := 0.0
	ssTot := 0.0

	for i := range predictions {
		ssRes += math.Pow(actuals[i]-predictions[i], 2)
		ssTot += math.Pow(actuals[i]-meanActual, 2)
	}

	if ssTot == 0 {
		return 0
	}

	return 1 - (ssRes / ssTot)
}

// calculateMean calculates the mean of a slice
func (pa *PredictiveAnalyzer) calculateMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sum := 0.0
	for _, value := range values {
		sum += value
	}

	return sum / float64(len(values))
}

// generatePredictions generates predictions periodically
func (pa *PredictiveAnalyzer) generatePredictions(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			pa.generatePeriodicPredictions()
		}
	}
}

// generatePeriodicPredictions generates predictions for all active models
func (pa *PredictiveAnalyzer) generatePeriodicPredictions() {
	pa.mu.RLock()
	activeModels := make([]*PredictiveModel, 0)
	for _, model := range pa.models {
		if model.Status == ModelStatusActive {
			activeModels = append(activeModels, model)
		}
	}
	pa.mu.RUnlock()

	for _, model := range activeModels {
		// Generate short-term prediction
		targetTime := time.Now().Add(15 * time.Minute)
		prediction := pa.generateSinglePrediction(model, targetTime)

		pa.logger.Debug(context.Background(), "Periodic prediction generated", map[string]interface{}{
			"model_id":        model.ModelID,
			"metric_name":     model.MetricName,
			"predicted_value": prediction.PredictedValue,
			"confidence":      prediction.ConfidenceLevel,
		})
	}
}

// validatePredictions validates predictions against actual values
func (pa *PredictiveAnalyzer) validatePredictions(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			pa.validatePastPredictions()
		}
	}
}

// validatePastPredictions validates past predictions
func (pa *PredictiveAnalyzer) validatePastPredictions() {
	pa.mu.Lock()
	defer pa.mu.Unlock()

	now := time.Now()
	for _, prediction := range pa.predictions {
		// Check if prediction target time has passed and we have actual data
		if now.After(prediction.TargetTime) && prediction.ActualValue == nil {
			// Try to find actual value from training data
			if trainingData, exists := pa.trainingData[prediction.MetricName]; exists {
				for _, point := range trainingData {
					// Find data point close to target time (within 5 minutes)
					if math.Abs(point.Timestamp.Sub(prediction.TargetTime).Minutes()) < 5 {
						actualValue := point.Value
						prediction.ActualValue = &actualValue
						error := math.Abs(prediction.PredictedValue - actualValue)
						prediction.Error = &error
						prediction.Accuracy = 1 - (error / math.Max(prediction.PredictedValue, actualValue))
						break
					}
				}
			}
		}
	}
}
