package analytics

import (
	"context"
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

// AnomalyDetector detects anomalies in real-time data streams
type AnomalyDetector struct {
	logger          *observability.Logger
	config          *AnalyticsConfig
	detectors       map[string]*MetricDetector
	anomalies       []*Anomaly
	alertThresholds map[string]*AnomalyThreshold
	baselineModels  map[string]*BaselineModel
	mu              sync.RWMutex
}

// MetricDetector detects anomalies for a specific metric
type MetricDetector struct {
	MetricName      string                 `json:"metric_name"`
	DetectionMethod AnomalyDetectionMethod `json:"detection_method"`
	Sensitivity     float64                `json:"sensitivity"`
	WindowSize      int                    `json:"window_size"`
	DataPoints      []DataPoint            `json:"data_points"`
	Statistics      *MetricStatistics      `json:"statistics"`
	LastUpdated     time.Time              `json:"last_updated"`
	mu              sync.RWMutex           `json:"-"`
}

// AnomalyDetectionMethod defines detection methods
type AnomalyDetectionMethod string

const (
	DetectionMethodStatistical     AnomalyDetectionMethod = "statistical"
	DetectionMethodIsolationForest AnomalyDetectionMethod = "isolation_forest"
	DetectionMethodZScore          AnomalyDetectionMethod = "z_score"
	DetectionMethodIQR             AnomalyDetectionMethod = "iqr"
	DetectionMethodMovingAverage   AnomalyDetectionMethod = "moving_average"
	DetectionMethodSeasonal        AnomalyDetectionMethod = "seasonal"
	DetectionMethodML              AnomalyDetectionMethod = "machine_learning"
)

// DataPoint represents a single data point
type DataPoint struct {
	Timestamp time.Time         `json:"timestamp"`
	Value     float64           `json:"value"`
	Tags      map[string]string `json:"tags"`
}

// MetricStatistics contains statistical information about a metric
type MetricStatistics struct {
	Mean              float64   `json:"mean"`
	StandardDeviation float64   `json:"standard_deviation"`
	Variance          float64   `json:"variance"`
	Min               float64   `json:"min"`
	Max               float64   `json:"max"`
	Median            float64   `json:"median"`
	Q1                float64   `json:"q1"`
	Q3                float64   `json:"q3"`
	IQR               float64   `json:"iqr"`
	Skewness          float64   `json:"skewness"`
	Kurtosis          float64   `json:"kurtosis"`
	LastCalculated    time.Time `json:"last_calculated"`
}

// Anomaly represents a detected anomaly
type Anomaly struct {
	AnomalyID       string                 `json:"anomaly_id"`
	MetricName      string                 `json:"metric_name"`
	DetectionMethod AnomalyDetectionMethod `json:"detection_method"`
	Timestamp       time.Time              `json:"timestamp"`
	Value           float64                `json:"value"`
	ExpectedValue   float64                `json:"expected_value"`
	Deviation       float64                `json:"deviation"`
	Severity        AnomalySeverity        `json:"severity"`
	Confidence      float64                `json:"confidence"`
	Description     string                 `json:"description"`
	Context         map[string]interface{} `json:"context"`
	Tags            map[string]string      `json:"tags"`
	Status          AnomalyStatus          `json:"status"`
	ResolvedAt      *time.Time             `json:"resolved_at,omitempty"`
}

// AnomalySeverity defines anomaly severity levels
type AnomalySeverity string

const (
	AnomalySeverityLow      AnomalySeverity = "low"
	AnomalySeverityMedium   AnomalySeverity = "medium"
	AnomalySeverityHigh     AnomalySeverity = "high"
	AnomalySeverityCritical AnomalySeverity = "critical"
)

// AnomalyStatus defines anomaly status
type AnomalyStatus string

const (
	AnomalyStatusActive        AnomalyStatus = "active"
	AnomalyStatusResolved      AnomalyStatus = "resolved"
	AnomalyStatusIgnored       AnomalyStatus = "ignored"
	AnomalyStatusInvestigating AnomalyStatus = "investigating"
)

// AnomalyThreshold defines thresholds for anomaly detection
type AnomalyThreshold struct {
	MetricName        string        `json:"metric_name"`
	LowSeverity       float64       `json:"low_severity"`
	MediumSeverity    float64       `json:"medium_severity"`
	HighSeverity      float64       `json:"high_severity"`
	CriticalSeverity  float64       `json:"critical_severity"`
	EnableAutoResolve bool          `json:"enable_auto_resolve"`
	AutoResolveTime   time.Duration `json:"auto_resolve_time"`
}

// BaselineModel represents a baseline model for anomaly detection
type BaselineModel struct {
	MetricName      string             `json:"metric_name"`
	ModelType       string             `json:"model_type"`
	TrainingData    []DataPoint        `json:"training_data"`
	Parameters      map[string]float64 `json:"parameters"`
	Accuracy        float64            `json:"accuracy"`
	LastTrained     time.Time          `json:"last_trained"`
	PredictionCache map[string]float64 `json:"prediction_cache"`
	mu              sync.RWMutex       `json:"-"`
}

// NewAnomalyDetector creates a new anomaly detector
func NewAnomalyDetector(logger *observability.Logger, config *AnalyticsConfig) *AnomalyDetector {
	return &AnomalyDetector{
		logger:          logger,
		config:          config,
		detectors:       make(map[string]*MetricDetector),
		anomalies:       make([]*Anomaly, 0),
		alertThresholds: make(map[string]*AnomalyThreshold),
		baselineModels:  make(map[string]*BaselineModel),
	}
}

// Start starts the anomaly detector
func (ad *AnomalyDetector) Start(ctx context.Context) error {
	ad.logger.Info(ctx, "Starting anomaly detector", map[string]interface{}{
		"sensitivity": ad.config.AnomalyDetectionSensitivity,
	})

	// Initialize default detectors
	ad.initializeDefaultDetectors()

	// Start background processing
	go ad.processAnomalies(ctx)
	go ad.updateBaselines(ctx)
	go ad.cleanupOldAnomalies(ctx)

	return nil
}

// RegisterMetricDetector registers a new metric detector
func (ad *AnomalyDetector) RegisterMetricDetector(metricName string, method AnomalyDetectionMethod, sensitivity float64, windowSize int) {
	ad.mu.Lock()
	defer ad.mu.Unlock()

	detector := &MetricDetector{
		MetricName:      metricName,
		DetectionMethod: method,
		Sensitivity:     sensitivity,
		WindowSize:      windowSize,
		DataPoints:      make([]DataPoint, 0, windowSize),
		Statistics:      &MetricStatistics{},
		LastUpdated:     time.Now(),
	}

	ad.detectors[metricName] = detector

	ad.logger.Info(context.Background(), "Metric detector registered", map[string]interface{}{
		"metric_name":      metricName,
		"detection_method": method,
		"sensitivity":      sensitivity,
		"window_size":      windowSize,
	})
}

// AddDataPoint adds a data point for anomaly detection
func (ad *AnomalyDetector) AddDataPoint(metricName string, value float64, tags map[string]string) {
	ad.mu.RLock()
	detector, exists := ad.detectors[metricName]
	ad.mu.RUnlock()

	if !exists {
		// Auto-register detector with default settings
		ad.RegisterMetricDetector(metricName, DetectionMethodStatistical, ad.config.AnomalyDetectionSensitivity, 100)
		ad.mu.RLock()
		detector = ad.detectors[metricName]
		ad.mu.RUnlock()
	}

	detector.mu.Lock()
	defer detector.mu.Unlock()

	dataPoint := DataPoint{
		Timestamp: time.Now(),
		Value:     value,
		Tags:      tags,
	}

	// Add data point to window
	detector.DataPoints = append(detector.DataPoints, dataPoint)

	// Maintain window size
	if len(detector.DataPoints) > detector.WindowSize {
		detector.DataPoints = detector.DataPoints[1:]
	}

	// Update statistics
	ad.updateStatistics(detector)

	// Check for anomalies
	if len(detector.DataPoints) >= 10 { // Minimum data points for detection
		if anomaly := ad.detectAnomaly(detector, dataPoint); anomaly != nil {
			ad.mu.Lock()
			ad.anomalies = append(ad.anomalies, anomaly)
			ad.mu.Unlock()

			ad.logger.Warn(context.Background(), "Anomaly detected", map[string]interface{}{
				"anomaly_id":  anomaly.AnomalyID,
				"metric_name": anomaly.MetricName,
				"value":       anomaly.Value,
				"expected":    anomaly.ExpectedValue,
				"deviation":   anomaly.Deviation,
				"severity":    anomaly.Severity,
				"confidence":  anomaly.Confidence,
			})
		}
	}

	detector.LastUpdated = time.Now()
}

// detectAnomaly detects if a data point is anomalous
func (ad *AnomalyDetector) detectAnomaly(detector *MetricDetector, dataPoint DataPoint) *Anomaly {
	switch detector.DetectionMethod {
	case DetectionMethodZScore:
		return ad.detectZScoreAnomaly(detector, dataPoint)
	case DetectionMethodIQR:
		return ad.detectIQRAnomaly(detector, dataPoint)
	case DetectionMethodMovingAverage:
		return ad.detectMovingAverageAnomaly(detector, dataPoint)
	case DetectionMethodStatistical:
		return ad.detectStatisticalAnomaly(detector, dataPoint)
	default:
		return ad.detectStatisticalAnomaly(detector, dataPoint)
	}
}

// detectZScoreAnomaly detects anomalies using Z-score method
func (ad *AnomalyDetector) detectZScoreAnomaly(detector *MetricDetector, dataPoint DataPoint) *Anomaly {
	stats := detector.Statistics
	if stats.StandardDeviation == 0 {
		return nil
	}

	zScore := math.Abs((dataPoint.Value - stats.Mean) / stats.StandardDeviation)
	threshold := 2.0 + (1.0-detector.Sensitivity)*2.0 // Sensitivity affects threshold

	if zScore > threshold {
		severity := ad.calculateSeverity(detector.MetricName, zScore, threshold)
		confidence := math.Min(zScore/threshold, 1.0)

		return &Anomaly{
			AnomalyID:       uuid.New().String(),
			MetricName:      detector.MetricName,
			DetectionMethod: DetectionMethodZScore,
			Timestamp:       dataPoint.Timestamp,
			Value:           dataPoint.Value,
			ExpectedValue:   stats.Mean,
			Deviation:       zScore,
			Severity:        severity,
			Confidence:      confidence,
			Description:     fmt.Sprintf("Z-score anomaly detected: %.2f (threshold: %.2f)", zScore, threshold),
			Context: map[string]interface{}{
				"z_score":            zScore,
				"threshold":          threshold,
				"mean":               stats.Mean,
				"standard_deviation": stats.StandardDeviation,
			},
			Tags:   dataPoint.Tags,
			Status: AnomalyStatusActive,
		}
	}

	return nil
}

// detectIQRAnomaly detects anomalies using Interquartile Range method
func (ad *AnomalyDetector) detectIQRAnomaly(detector *MetricDetector, dataPoint DataPoint) *Anomaly {
	stats := detector.Statistics
	if stats.IQR == 0 {
		return nil
	}

	lowerBound := stats.Q1 - 1.5*stats.IQR
	upperBound := stats.Q3 + 1.5*stats.IQR

	if dataPoint.Value < lowerBound || dataPoint.Value > upperBound {
		deviation := math.Max(lowerBound-dataPoint.Value, dataPoint.Value-upperBound)
		severity := ad.calculateSeverity(detector.MetricName, deviation, stats.IQR)
		confidence := math.Min(deviation/stats.IQR, 1.0)

		return &Anomaly{
			AnomalyID:       uuid.New().String(),
			MetricName:      detector.MetricName,
			DetectionMethod: DetectionMethodIQR,
			Timestamp:       dataPoint.Timestamp,
			Value:           dataPoint.Value,
			ExpectedValue:   (stats.Q1 + stats.Q3) / 2,
			Deviation:       deviation,
			Severity:        severity,
			Confidence:      confidence,
			Description:     fmt.Sprintf("IQR anomaly detected: value %.2f outside bounds [%.2f, %.2f]", dataPoint.Value, lowerBound, upperBound),
			Context: map[string]interface{}{
				"lower_bound": lowerBound,
				"upper_bound": upperBound,
				"q1":          stats.Q1,
				"q3":          stats.Q3,
				"iqr":         stats.IQR,
			},
			Tags:   dataPoint.Tags,
			Status: AnomalyStatusActive,
		}
	}

	return nil
}

// detectMovingAverageAnomaly detects anomalies using moving average method
func (ad *AnomalyDetector) detectMovingAverageAnomaly(detector *MetricDetector, dataPoint DataPoint) *Anomaly {
	if len(detector.DataPoints) < 10 {
		return nil
	}

	// Calculate moving average of last 10 points (excluding current)
	recentPoints := detector.DataPoints[len(detector.DataPoints)-10:]
	sum := 0.0
	for _, point := range recentPoints {
		sum += point.Value
	}
	movingAvg := sum / float64(len(recentPoints))

	// Calculate moving standard deviation
	variance := 0.0
	for _, point := range recentPoints {
		variance += math.Pow(point.Value-movingAvg, 2)
	}
	movingStdDev := math.Sqrt(variance / float64(len(recentPoints)))

	if movingStdDev == 0 {
		return nil
	}

	deviation := math.Abs(dataPoint.Value - movingAvg)
	threshold := movingStdDev * (2.0 + (1.0-detector.Sensitivity)*2.0)

	if deviation > threshold {
		severity := ad.calculateSeverity(detector.MetricName, deviation, threshold)
		confidence := math.Min(deviation/threshold, 1.0)

		return &Anomaly{
			AnomalyID:       uuid.New().String(),
			MetricName:      detector.MetricName,
			DetectionMethod: DetectionMethodMovingAverage,
			Timestamp:       dataPoint.Timestamp,
			Value:           dataPoint.Value,
			ExpectedValue:   movingAvg,
			Deviation:       deviation,
			Severity:        severity,
			Confidence:      confidence,
			Description:     fmt.Sprintf("Moving average anomaly detected: deviation %.2f (threshold: %.2f)", deviation, threshold),
			Context: map[string]interface{}{
				"moving_average": movingAvg,
				"moving_std_dev": movingStdDev,
				"threshold":      threshold,
				"window_size":    len(recentPoints),
			},
			Tags:   dataPoint.Tags,
			Status: AnomalyStatusActive,
		}
	}

	return nil
}

// detectStatisticalAnomaly detects anomalies using statistical methods
func (ad *AnomalyDetector) detectStatisticalAnomaly(detector *MetricDetector, dataPoint DataPoint) *Anomaly {
	// Combine multiple methods for better accuracy
	zScoreAnomaly := ad.detectZScoreAnomaly(detector, dataPoint)
	iqrAnomaly := ad.detectIQRAnomaly(detector, dataPoint)

	// If both methods detect an anomaly, it's likely a true anomaly
	if zScoreAnomaly != nil && iqrAnomaly != nil {
		// Use the higher confidence score
		if zScoreAnomaly.Confidence > iqrAnomaly.Confidence {
			zScoreAnomaly.Description = "Statistical anomaly detected (Z-score + IQR)"
			return zScoreAnomaly
		} else {
			iqrAnomaly.Description = "Statistical anomaly detected (IQR + Z-score)"
			return iqrAnomaly
		}
	}

	// If only one method detects an anomaly, use higher threshold
	if zScoreAnomaly != nil && zScoreAnomaly.Confidence > 0.8 {
		return zScoreAnomaly
	}
	if iqrAnomaly != nil && iqrAnomaly.Confidence > 0.8 {
		return iqrAnomaly
	}

	return nil
}

// updateStatistics updates statistical information for a detector
func (ad *AnomalyDetector) updateStatistics(detector *MetricDetector) {
	if len(detector.DataPoints) == 0 {
		return
	}

	values := make([]float64, len(detector.DataPoints))
	sum := 0.0
	for i, point := range detector.DataPoints {
		values[i] = point.Value
		sum += point.Value
	}

	n := float64(len(values))
	mean := sum / n

	// Calculate variance and standard deviation
	variance := 0.0
	for _, value := range values {
		variance += math.Pow(value-mean, 2)
	}
	variance /= n
	stdDev := math.Sqrt(variance)

	// Sort values for percentile calculations
	sort.Float64s(values)

	// Calculate percentiles
	min := values[0]
	max := values[len(values)-1]
	median := ad.percentile(values, 50)
	q1 := ad.percentile(values, 25)
	q3 := ad.percentile(values, 75)
	iqr := q3 - q1

	// Update statistics
	detector.Statistics.Mean = mean
	detector.Statistics.StandardDeviation = stdDev
	detector.Statistics.Variance = variance
	detector.Statistics.Min = min
	detector.Statistics.Max = max
	detector.Statistics.Median = median
	detector.Statistics.Q1 = q1
	detector.Statistics.Q3 = q3
	detector.Statistics.IQR = iqr
	detector.Statistics.LastCalculated = time.Now()
}

// percentile calculates the percentile of a sorted slice
func (ad *AnomalyDetector) percentile(sortedValues []float64, p float64) float64 {
	if len(sortedValues) == 0 {
		return 0
	}
	if len(sortedValues) == 1 {
		return sortedValues[0]
	}

	index := (p / 100.0) * float64(len(sortedValues)-1)
	lower := int(math.Floor(index))
	upper := int(math.Ceil(index))

	if lower == upper {
		return sortedValues[lower]
	}

	weight := index - float64(lower)
	return sortedValues[lower]*(1-weight) + sortedValues[upper]*weight
}

// calculateSeverity calculates anomaly severity based on deviation
func (ad *AnomalyDetector) calculateSeverity(metricName string, deviation, threshold float64) AnomalySeverity {
	// Check if custom thresholds exist
	if customThreshold, exists := ad.alertThresholds[metricName]; exists {
		ratio := deviation / threshold
		if ratio >= customThreshold.CriticalSeverity {
			return AnomalySeverityCritical
		} else if ratio >= customThreshold.HighSeverity {
			return AnomalySeverityHigh
		} else if ratio >= customThreshold.MediumSeverity {
			return AnomalySeverityMedium
		} else {
			return AnomalySeverityLow
		}
	}

	// Default severity calculation
	ratio := deviation / threshold
	if ratio >= 4.0 {
		return AnomalySeverityCritical
	} else if ratio >= 3.0 {
		return AnomalySeverityHigh
	} else if ratio >= 2.0 {
		return AnomalySeverityMedium
	} else {
		return AnomalySeverityLow
	}
}

// GetActiveAnomalies returns all active anomalies
func (ad *AnomalyDetector) GetActiveAnomalies() []*Anomaly {
	ad.mu.RLock()
	defer ad.mu.RUnlock()

	activeAnomalies := make([]*Anomaly, 0)
	for _, anomaly := range ad.anomalies {
		if anomaly.Status == AnomalyStatusActive {
			activeAnomalies = append(activeAnomalies, anomaly)
		}
	}

	return activeAnomalies
}

// GetAnomaliesByMetric returns anomalies for a specific metric
func (ad *AnomalyDetector) GetAnomaliesByMetric(metricName string) []*Anomaly {
	ad.mu.RLock()
	defer ad.mu.RUnlock()

	metricAnomalies := make([]*Anomaly, 0)
	for _, anomaly := range ad.anomalies {
		if anomaly.MetricName == metricName {
			metricAnomalies = append(metricAnomalies, anomaly)
		}
	}

	return metricAnomalies
}

// ResolveAnomaly marks an anomaly as resolved
func (ad *AnomalyDetector) ResolveAnomaly(anomalyID string) error {
	ad.mu.Lock()
	defer ad.mu.Unlock()

	for _, anomaly := range ad.anomalies {
		if anomaly.AnomalyID == anomalyID {
			anomaly.Status = AnomalyStatusResolved
			now := time.Now()
			anomaly.ResolvedAt = &now
			return nil
		}
	}

	return fmt.Errorf("anomaly not found: %s", anomalyID)
}

// initializeDefaultDetectors initializes default metric detectors
func (ad *AnomalyDetector) initializeDefaultDetectors() {
	defaultMetrics := []struct {
		name   string
		method AnomalyDetectionMethod
	}{
		{"cpu_usage", DetectionMethodStatistical},
		{"memory_usage", DetectionMethodStatistical},
		{"response_time", DetectionMethodZScore},
		{"error_rate", DetectionMethodIQR},
		{"request_count", DetectionMethodMovingAverage},
		{"trading_volume", DetectionMethodStatistical},
		{"price_change", DetectionMethodZScore},
	}

	for _, metric := range defaultMetrics {
		ad.RegisterMetricDetector(metric.name, metric.method, ad.config.AnomalyDetectionSensitivity, 100)
	}
}

// processAnomalies processes anomalies in the background
func (ad *AnomalyDetector) processAnomalies(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ad.processAnomalyBatch()
		}
	}
}

// processAnomalyBatch processes a batch of anomalies
func (ad *AnomalyDetector) processAnomalyBatch() {
	ad.mu.RLock()
	activeAnomalies := ad.GetActiveAnomalies()
	ad.mu.RUnlock()

	for _, anomaly := range activeAnomalies {
		// Check for auto-resolution
		if threshold, exists := ad.alertThresholds[anomaly.MetricName]; exists {
			if threshold.EnableAutoResolve {
				if time.Since(anomaly.Timestamp) > threshold.AutoResolveTime {
					ad.ResolveAnomaly(anomaly.AnomalyID)
					ad.logger.Info(context.Background(), "Anomaly auto-resolved", map[string]interface{}{
						"anomaly_id":  anomaly.AnomalyID,
						"metric_name": anomaly.MetricName,
					})
				}
			}
		}
	}
}

// updateBaselines updates baseline models
func (ad *AnomalyDetector) updateBaselines(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ad.updateBaselineModels()
		}
	}
}

// updateBaselineModels updates all baseline models
func (ad *AnomalyDetector) updateBaselineModels() {
	ad.mu.RLock()
	defer ad.mu.RUnlock()

	for metricName, detector := range ad.detectors {
		if len(detector.DataPoints) >= 100 { // Minimum data for baseline
			ad.updateBaselineModel(metricName, detector.DataPoints)
		}
	}
}

// updateBaselineModel updates a baseline model for a specific metric
func (ad *AnomalyDetector) updateBaselineModel(metricName string, dataPoints []DataPoint) {
	// Simple baseline model using historical statistics
	model := &BaselineModel{
		MetricName:      metricName,
		ModelType:       "statistical",
		TrainingData:    dataPoints,
		Parameters:      make(map[string]float64),
		LastTrained:     time.Now(),
		PredictionCache: make(map[string]float64),
	}

	// Calculate model parameters
	values := make([]float64, len(dataPoints))
	for i, point := range dataPoints {
		values[i] = point.Value
	}

	sum := 0.0
	for _, value := range values {
		sum += value
	}
	mean := sum / float64(len(values))

	variance := 0.0
	for _, value := range values {
		variance += math.Pow(value-mean, 2)
	}
	variance /= float64(len(values))

	model.Parameters["mean"] = mean
	model.Parameters["variance"] = variance
	model.Parameters["std_dev"] = math.Sqrt(variance)

	// Calculate model accuracy (simplified)
	model.Accuracy = 0.85 // Placeholder

	ad.baselineModels[metricName] = model
}

// cleanupOldAnomalies removes old resolved anomalies
func (ad *AnomalyDetector) cleanupOldAnomalies(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ad.performCleanup()
		}
	}
}

// performCleanup removes old anomalies
func (ad *AnomalyDetector) performCleanup() {
	ad.mu.Lock()
	defer ad.mu.Unlock()

	cutoffTime := time.Now().Add(-7 * 24 * time.Hour) // Keep for 7 days
	filteredAnomalies := make([]*Anomaly, 0)

	for _, anomaly := range ad.anomalies {
		if anomaly.Status == AnomalyStatusActive || anomaly.Timestamp.After(cutoffTime) {
			filteredAnomalies = append(filteredAnomalies, anomaly)
		}
	}

	removed := len(ad.anomalies) - len(filteredAnomalies)
	ad.anomalies = filteredAnomalies

	if removed > 0 {
		ad.logger.Info(context.Background(), "Cleaned up old anomalies", map[string]interface{}{
			"removed_count": removed,
		})
	}
}
