package ml

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// PredictiveAnalytics provides price prediction and market analysis capabilities
type PredictiveAnalytics struct {
	logger      *observability.Logger
	config      PredictiveConfig
	mlFramework *MLFramework

	// Market regime detection
	regimeDetector *MarketRegimeDetector

	// Price prediction models
	priceModels map[string]*PricePredictionModel

	// Feature extractors
	featureExtractors map[string]*FeatureExtractor
}

// PredictiveConfig contains predictive analytics configuration
type PredictiveConfig struct {
	EnablePricePrediction   bool            `json:"enable_price_prediction"`
	EnableRegimeDetection   bool            `json:"enable_regime_detection"`
	EnableSentimentAnalysis bool            `json:"enable_sentiment_analysis"`
	PredictionHorizon       time.Duration   `json:"prediction_horizon"`
	UpdateInterval          time.Duration   `json:"update_interval"`
	MinDataPoints           int             `json:"min_data_points"`
	MaxDataPoints           int             `json:"max_data_points"`
	ConfidenceThreshold     decimal.Decimal `json:"confidence_threshold"`
	ModelRetentionPeriod    time.Duration   `json:"model_retention_period"`
}

// PricePredictionModel represents a price prediction model
type PricePredictionModel struct {
	ID              uuid.UUID              `json:"id"`
	Symbol          string                 `json:"symbol"`
	ModelType       PredictionModelType    `json:"model_type"`
	Features        []string               `json:"features"`
	TrainedModel    *TrainedModel          `json:"trained_model"`
	LastPrediction  *PricePrediction       `json:"last_prediction"`
	Accuracy        decimal.Decimal        `json:"accuracy"`
	CreatedAt       time.Time              `json:"created_at"`
	LastUpdated     time.Time              `json:"last_updated"`
	PredictionCount int64                  `json:"prediction_count"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// PredictionModelType represents the type of prediction model
type PredictionModelType string

const (
	PredictionModelTypeLSTM         PredictionModelType = "lstm"
	PredictionModelTypeGRU          PredictionModelType = "gru"
	PredictionModelTypeTransformer  PredictionModelType = "transformer"
	PredictionModelTypeARIMA        PredictionModelType = "arima"
	PredictionModelTypeGARCH        PredictionModelType = "garch"
	PredictionModelTypeRandomForest PredictionModelType = "random_forest"
	PredictionModelTypeEnsemble     PredictionModelType = "ensemble"
)

// PricePrediction represents a price prediction
type PricePrediction struct {
	ID             uuid.UUID              `json:"id"`
	Symbol         string                 `json:"symbol"`
	ModelID        uuid.UUID              `json:"model_id"`
	CurrentPrice   decimal.Decimal        `json:"current_price"`
	PredictedPrice decimal.Decimal        `json:"predicted_price"`
	PriceChange    decimal.Decimal        `json:"price_change"`
	PriceChangePct decimal.Decimal        `json:"price_change_pct"`
	Direction      PriceDirection         `json:"direction"`
	Confidence     decimal.Decimal        `json:"confidence"`
	Horizon        time.Duration          `json:"horizon"`
	Features       map[string]float64     `json:"features"`
	Timestamp      time.Time              `json:"timestamp"`
	ExpiresAt      time.Time              `json:"expires_at"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// PriceDirection represents the predicted price direction
type PriceDirection string

const (
	PriceDirectionUp       PriceDirection = "up"
	PriceDirectionDown     PriceDirection = "down"
	PriceDirectionSideways PriceDirection = "sideways"
)

// MarketRegimeDetector detects market regimes (bull, bear, sideways)
type MarketRegimeDetector struct {
	currentRegime    MarketRegime
	regimeHistory    []*RegimeChange
	lastUpdate       time.Time
	volatilityWindow int
	trendWindow      int
}

// MarketRegime represents the current market regime
type MarketRegime string

const (
	MarketRegimeBull     MarketRegime = "bull"
	MarketRegimeBear     MarketRegime = "bear"
	MarketRegimeSideways MarketRegime = "sideways"
	MarketRegimeVolatile MarketRegime = "volatile"
)

// RegimeChange represents a market regime change
type RegimeChange struct {
	FromRegime MarketRegime       `json:"from_regime"`
	ToRegime   MarketRegime       `json:"to_regime"`
	Timestamp  time.Time          `json:"timestamp"`
	Confidence decimal.Decimal    `json:"confidence"`
	Indicators map[string]float64 `json:"indicators"`
}

// FeatureExtractor extracts features from market data
type FeatureExtractor struct {
	Name       string                 `json:"name"`
	Type       FeatureExtractorType   `json:"type"`
	Config     map[string]interface{} `json:"config"`
	LastUpdate time.Time              `json:"last_update"`
}

// FeatureExtractorType represents the type of feature extractor
type FeatureExtractorType string

const (
	FeatureExtractorTypeTechnical   FeatureExtractorType = "technical"
	FeatureExtractorTypeFundamental FeatureExtractorType = "fundamental"
	FeatureExtractorTypeSentiment   FeatureExtractorType = "sentiment"
	FeatureExtractorTypeMacro       FeatureExtractorType = "macro"
	FeatureExtractorTypeVolume      FeatureExtractorType = "volume"
	FeatureExtractorTypePrice       FeatureExtractorType = "price"
)

// MarketData represents market data for analysis
type MarketData struct {
	Symbol    string          `json:"symbol"`
	Timestamp time.Time       `json:"timestamp"`
	Open      decimal.Decimal `json:"open"`
	High      decimal.Decimal `json:"high"`
	Low       decimal.Decimal `json:"low"`
	Close     decimal.Decimal `json:"close"`
	Volume    decimal.Decimal `json:"volume"`
}

// TechnicalIndicators contains technical analysis indicators
type TechnicalIndicators struct {
	SMA20          decimal.Decimal `json:"sma_20"`
	SMA50          decimal.Decimal `json:"sma_50"`
	SMA200         decimal.Decimal `json:"sma_200"`
	EMA12          decimal.Decimal `json:"ema_12"`
	EMA26          decimal.Decimal `json:"ema_26"`
	RSI            decimal.Decimal `json:"rsi"`
	MACD           decimal.Decimal `json:"macd"`
	MACDSignal     decimal.Decimal `json:"macd_signal"`
	MACDHist       decimal.Decimal `json:"macd_hist"`
	BollingerUpper decimal.Decimal `json:"bollinger_upper"`
	BollingerLower decimal.Decimal `json:"bollinger_lower"`
	ATR            decimal.Decimal `json:"atr"`
	Stochastic     decimal.Decimal `json:"stochastic"`
	Williams       decimal.Decimal `json:"williams"`
	CCI            decimal.Decimal `json:"cci"`
	ADX            decimal.Decimal `json:"adx"`
	OBV            decimal.Decimal `json:"obv"`
	VWAP           decimal.Decimal `json:"vwap"`
}

// SentimentData represents market sentiment data
type SentimentData struct {
	Symbol         string          `json:"symbol"`
	Timestamp      time.Time       `json:"timestamp"`
	SentimentScore decimal.Decimal `json:"sentiment_score"`
	FearGreedIndex decimal.Decimal `json:"fear_greed_index"`
	SocialMentions int64           `json:"social_mentions"`
	NewsScore      decimal.Decimal `json:"news_score"`
	OptionsFlow    decimal.Decimal `json:"options_flow"`
	InsiderTrading decimal.Decimal `json:"insider_trading"`
}

// NewPredictiveAnalytics creates a new predictive analytics engine
func NewPredictiveAnalytics(
	logger *observability.Logger,
	config PredictiveConfig,
	mlFramework *MLFramework,
) *PredictiveAnalytics {
	// Set defaults
	if config.PredictionHorizon == 0 {
		config.PredictionHorizon = 1 * time.Hour
	}
	if config.UpdateInterval == 0 {
		config.UpdateInterval = 5 * time.Minute
	}
	if config.MinDataPoints == 0 {
		config.MinDataPoints = 100
	}
	if config.MaxDataPoints == 0 {
		config.MaxDataPoints = 10000
	}
	if config.ConfidenceThreshold.IsZero() {
		config.ConfidenceThreshold = decimal.NewFromFloat(0.7)
	}
	if config.ModelRetentionPeriod == 0 {
		config.ModelRetentionPeriod = 30 * 24 * time.Hour // 30 days
	}

	return &PredictiveAnalytics{
		logger:      logger,
		config:      config,
		mlFramework: mlFramework,
		regimeDetector: &MarketRegimeDetector{
			currentRegime:    MarketRegimeSideways,
			regimeHistory:    make([]*RegimeChange, 0),
			volatilityWindow: 20,
			trendWindow:      50,
		},
		priceModels:       make(map[string]*PricePredictionModel),
		featureExtractors: make(map[string]*FeatureExtractor),
	}
}

// PredictPrice predicts future price for a symbol
func (pa *PredictiveAnalytics) PredictPrice(ctx context.Context, symbol string, marketData []*MarketData) (*PricePrediction, error) {
	if !pa.config.EnablePricePrediction {
		return nil, fmt.Errorf("price prediction is disabled")
	}

	pa.logger.Info(ctx, "Predicting price", map[string]interface{}{
		"symbol":      symbol,
		"data_points": len(marketData),
		"horizon":     pa.config.PredictionHorizon,
	})

	// Get or create prediction model for symbol
	model, err := pa.getOrCreatePredictionModel(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get prediction model: %w", err)
	}

	// Extract features from market data
	features, err := pa.extractFeatures(ctx, symbol, marketData)
	if err != nil {
		return nil, fmt.Errorf("failed to extract features: %w", err)
	}

	// Make prediction using ML framework
	predictionRequest := &PredictionRequest{
		ModelID:   model.TrainedModel.ModelID,
		Features:  pa.convertFeaturesToFloat64(features),
		Timestamp: time.Now(),
		Metadata: map[string]interface{}{
			"symbol":  symbol,
			"horizon": pa.config.PredictionHorizon.String(),
		},
	}

	predictionResult, err := pa.mlFramework.Predict(ctx, predictionRequest)
	if err != nil {
		return nil, fmt.Errorf("prediction failed: %w", err)
	}

	// Convert prediction result to price prediction
	currentPrice := marketData[len(marketData)-1].Close
	predictedPrice := decimal.NewFromFloat(predictionResult.Predictions[0])
	priceChange := predictedPrice.Sub(currentPrice)
	priceChangePct := priceChange.Div(currentPrice).Mul(decimal.NewFromInt(100))

	// Determine direction
	direction := PriceDirectionSideways
	if priceChangePct.GreaterThan(decimal.NewFromFloat(1.0)) {
		direction = PriceDirectionUp
	} else if priceChangePct.LessThan(decimal.NewFromFloat(-1.0)) {
		direction = PriceDirectionDown
	}

	// Calculate confidence (simplified)
	confidence := decimal.NewFromFloat(0.8) // Placeholder
	if len(predictionResult.Confidence) > 0 {
		confidence = decimal.NewFromFloat(predictionResult.Confidence[0])
	}

	prediction := &PricePrediction{
		ID:             uuid.New(),
		Symbol:         symbol,
		ModelID:        model.ID,
		CurrentPrice:   currentPrice,
		PredictedPrice: predictedPrice,
		PriceChange:    priceChange,
		PriceChangePct: priceChangePct,
		Direction:      direction,
		Confidence:     confidence,
		Horizon:        pa.config.PredictionHorizon,
		Features:       pa.convertFeaturesToMap(features),
		Timestamp:      time.Now(),
		ExpiresAt:      time.Now().Add(pa.config.PredictionHorizon),
		Metadata:       make(map[string]interface{}),
	}

	// Update model
	model.LastPrediction = prediction
	model.PredictionCount++
	model.LastUpdated = time.Now()

	pa.logger.Info(ctx, "Price prediction completed", map[string]interface{}{
		"symbol":           symbol,
		"current_price":    currentPrice.String(),
		"predicted_price":  predictedPrice.String(),
		"price_change_pct": priceChangePct.String(),
		"direction":        string(direction),
		"confidence":       confidence.String(),
	})

	return prediction, nil
}

// DetectMarketRegime detects the current market regime
func (pa *PredictiveAnalytics) DetectMarketRegime(ctx context.Context, marketData []*MarketData) (*MarketRegime, error) {
	if !pa.config.EnableRegimeDetection {
		return nil, fmt.Errorf("regime detection is disabled")
	}

	if len(marketData) < pa.regimeDetector.trendWindow {
		return &pa.regimeDetector.currentRegime, nil
	}

	// Calculate trend indicators
	prices := make([]decimal.Decimal, len(marketData))
	for i, data := range marketData {
		prices[i] = data.Close
	}

	// Calculate moving averages
	shortMA := pa.calculateSMA(prices[len(prices)-20:], 20)
	longMA := pa.calculateSMA(prices[len(prices)-50:], 50)

	// Calculate volatility
	returns := pa.calculateReturns(prices)
	volatility := pa.calculateVolatility(returns[len(returns)-pa.regimeDetector.volatilityWindow:])

	// Determine regime
	newRegime := pa.regimeDetector.currentRegime

	// High volatility regime
	if volatility.GreaterThan(decimal.NewFromFloat(0.05)) { // 5% daily volatility
		newRegime = MarketRegimeVolatile
	} else if shortMA.GreaterThan(longMA.Mul(decimal.NewFromFloat(1.02))) { // 2% above
		newRegime = MarketRegimeBull
	} else if shortMA.LessThan(longMA.Mul(decimal.NewFromFloat(0.98))) { // 2% below
		newRegime = MarketRegimeBear
	} else {
		newRegime = MarketRegimeSideways
	}

	// Check for regime change
	if newRegime != pa.regimeDetector.currentRegime {
		regimeChange := &RegimeChange{
			FromRegime: pa.regimeDetector.currentRegime,
			ToRegime:   newRegime,
			Timestamp:  time.Now(),
			Confidence: decimal.NewFromFloat(0.8), // Simplified confidence
			Indicators: map[string]float64{
				"short_ma":   shortMA.InexactFloat64(),
				"long_ma":    longMA.InexactFloat64(),
				"volatility": volatility.InexactFloat64(),
			},
		}

		pa.regimeDetector.regimeHistory = append(pa.regimeDetector.regimeHistory, regimeChange)
		pa.regimeDetector.currentRegime = newRegime

		pa.logger.Info(ctx, "Market regime change detected", map[string]interface{}{
			"from_regime": string(regimeChange.FromRegime),
			"to_regime":   string(regimeChange.ToRegime),
			"confidence":  regimeChange.Confidence.String(),
		})
	}

	pa.regimeDetector.lastUpdate = time.Now()

	return &newRegime, nil
}

// GetCurrentRegime returns the current market regime
func (pa *PredictiveAnalytics) GetCurrentRegime() MarketRegime {
	return pa.regimeDetector.currentRegime
}

// GetRegimeHistory returns the history of regime changes
func (pa *PredictiveAnalytics) GetRegimeHistory() []*RegimeChange {
	return pa.regimeDetector.regimeHistory
}

// Private helper methods

// getOrCreatePredictionModel gets or creates a prediction model for a symbol
func (pa *PredictiveAnalytics) getOrCreatePredictionModel(ctx context.Context, symbol string) (*PricePredictionModel, error) {
	if model, exists := pa.priceModels[symbol]; exists {
		return model, nil
	}

	// Create new model
	model := &PricePredictionModel{
		ID:              uuid.New(),
		Symbol:          symbol,
		ModelType:       PredictionModelTypeLSTM, // Default to LSTM
		Features:        pa.getDefaultFeatures(),
		Accuracy:        decimal.NewFromFloat(0.5), // Initial accuracy
		CreatedAt:       time.Now(),
		LastUpdated:     time.Now(),
		PredictionCount: 0,
		Metadata:        make(map[string]interface{}),
	}

	// Create and train ML model (simplified)
	mlModel := &Model{
		Name:          fmt.Sprintf("price_prediction_%s", symbol),
		Type:          ModelTypeTimeSeries,
		Algorithm:     AlgorithmLSTM,
		InputFeatures: model.Features,
		OutputTargets: []string{"price"},
		Hyperparams: map[string]interface{}{
			"sequence_length": 60,
			"hidden_units":    50,
			"dropout":         0.2,
		},
	}

	if err := pa.mlFramework.RegisterModel(ctx, mlModel); err != nil {
		return nil, fmt.Errorf("failed to register ML model: %w", err)
	}

	// Create placeholder trained model
	trainedModel := &TrainedModel{
		ModelID:      mlModel.ID,
		Version:      "v1.0",
		TrainedAt:    time.Now(),
		TrainingTime: 5 * time.Minute,
		Status:       ModelStatusTrained,
		Metrics: &ModelMetrics{
			MSE:         0.001,
			RMSE:        0.032,
			MAE:         0.025,
			R2Score:     0.85,
			SharpeRatio: 1.2,
		},
		Metadata: make(map[string]interface{}),
	}

	model.TrainedModel = trainedModel
	pa.priceModels[symbol] = model

	pa.logger.Info(ctx, "Created new prediction model", map[string]interface{}{
		"symbol":     symbol,
		"model_id":   model.ID.String(),
		"model_type": string(model.ModelType),
	})

	return model, nil
}

// extractFeatures extracts features from market data
func (pa *PredictiveAnalytics) extractFeatures(ctx context.Context, symbol string, marketData []*MarketData) (map[string]interface{}, error) {
	if len(marketData) < 20 {
		return nil, fmt.Errorf("insufficient data for feature extraction")
	}

	features := make(map[string]interface{})

	// Extract price features
	prices := make([]decimal.Decimal, len(marketData))
	volumes := make([]decimal.Decimal, len(marketData))
	for i, data := range marketData {
		prices[i] = data.Close
		volumes[i] = data.Volume
	}

	// Calculate technical indicators
	indicators := pa.calculateTechnicalIndicators(marketData)

	// Price-based features
	features["current_price"] = prices[len(prices)-1].InexactFloat64()
	features["price_change_1d"] = pa.calculatePriceChange(prices, 1).InexactFloat64()
	features["price_change_7d"] = pa.calculatePriceChange(prices, 7).InexactFloat64()
	features["price_change_30d"] = pa.calculatePriceChange(prices, 30).InexactFloat64()

	// Technical indicators
	features["sma_20"] = indicators.SMA20.InexactFloat64()
	features["sma_50"] = indicators.SMA50.InexactFloat64()
	features["rsi"] = indicators.RSI.InexactFloat64()
	features["macd"] = indicators.MACD.InexactFloat64()
	features["atr"] = indicators.ATR.InexactFloat64()

	// Volume features
	features["volume_avg_20"] = pa.calculateSMA(volumes[len(volumes)-20:], 20).InexactFloat64()
	features["volume_ratio"] = volumes[len(volumes)-1].Div(pa.calculateSMA(volumes[len(volumes)-20:], 20)).InexactFloat64()

	// Volatility features
	returns := pa.calculateReturns(prices)
	features["volatility_20"] = pa.calculateVolatility(returns[len(returns)-20:]).InexactFloat64()

	// Market regime feature
	regime := pa.regimeDetector.currentRegime
	features["market_regime"] = pa.encodeRegime(regime)

	return features, nil
}

// calculateTechnicalIndicators calculates technical analysis indicators
func (pa *PredictiveAnalytics) calculateTechnicalIndicators(marketData []*MarketData) *TechnicalIndicators {
	prices := make([]decimal.Decimal, len(marketData))
	highs := make([]decimal.Decimal, len(marketData))
	lows := make([]decimal.Decimal, len(marketData))
	volumes := make([]decimal.Decimal, len(marketData))

	for i, data := range marketData {
		prices[i] = data.Close
		highs[i] = data.High
		lows[i] = data.Low
		volumes[i] = data.Volume
	}

	indicators := &TechnicalIndicators{}

	// Simple Moving Averages
	if len(prices) >= 20 {
		indicators.SMA20 = pa.calculateSMA(prices[len(prices)-20:], 20)
	}
	if len(prices) >= 50 {
		indicators.SMA50 = pa.calculateSMA(prices[len(prices)-50:], 50)
	}
	if len(prices) >= 200 {
		indicators.SMA200 = pa.calculateSMA(prices[len(prices)-200:], 200)
	}

	// RSI
	if len(prices) >= 14 {
		indicators.RSI = pa.calculateRSI(prices, 14)
	}

	// MACD
	if len(prices) >= 26 {
		macd, signal, hist := pa.calculateMACD(prices)
		indicators.MACD = macd
		indicators.MACDSignal = signal
		indicators.MACDHist = hist
	}

	// ATR
	if len(prices) >= 14 {
		indicators.ATR = pa.calculateATR(highs, lows, prices, 14)
	}

	return indicators
}

// getDefaultFeatures returns default feature names for price prediction
func (pa *PredictiveAnalytics) getDefaultFeatures() []string {
	return []string{
		"current_price",
		"price_change_1d",
		"price_change_7d",
		"sma_20",
		"sma_50",
		"rsi",
		"macd",
		"atr",
		"volume_ratio",
		"volatility_20",
		"market_regime",
	}
}

// convertFeaturesToFloat64 converts features map to float64 slice
func (pa *PredictiveAnalytics) convertFeaturesToFloat64(features map[string]interface{}) []float64 {
	result := make([]float64, 0, len(features))

	// Convert in consistent order based on default features
	for _, featureName := range pa.getDefaultFeatures() {
		if value, exists := features[featureName]; exists {
			if floatValue, ok := value.(float64); ok {
				result = append(result, floatValue)
			} else {
				result = append(result, 0.0) // Default value
			}
		} else {
			result = append(result, 0.0) // Default value
		}
	}

	return result
}

// calculateSMA calculates Simple Moving Average
func (pa *PredictiveAnalytics) calculateSMA(prices []decimal.Decimal, period int) decimal.Decimal {
	if len(prices) < period {
		return decimal.NewFromInt(0)
	}

	sum := decimal.NewFromInt(0)
	for i := len(prices) - period; i < len(prices); i++ {
		sum = sum.Add(prices[i])
	}

	return sum.Div(decimal.NewFromInt(int64(period)))
}

// calculateReturns calculates price returns
func (pa *PredictiveAnalytics) calculateReturns(prices []decimal.Decimal) []decimal.Decimal {
	if len(prices) < 2 {
		return []decimal.Decimal{}
	}

	returns := make([]decimal.Decimal, len(prices)-1)
	for i := 1; i < len(prices); i++ {
		if prices[i-1].GreaterThan(decimal.NewFromInt(0)) {
			returns[i-1] = prices[i].Sub(prices[i-1]).Div(prices[i-1])
		}
	}

	return returns
}

// calculateVolatility calculates volatility from returns
func (pa *PredictiveAnalytics) calculateVolatility(returns []decimal.Decimal) decimal.Decimal {
	if len(returns) == 0 {
		return decimal.NewFromInt(0)
	}

	// Calculate mean
	sum := decimal.NewFromInt(0)
	for _, ret := range returns {
		sum = sum.Add(ret)
	}
	mean := sum.Div(decimal.NewFromInt(int64(len(returns))))

	// Calculate variance
	sumSquaredDiffs := decimal.NewFromInt(0)
	for _, ret := range returns {
		diff := ret.Sub(mean)
		sumSquaredDiffs = sumSquaredDiffs.Add(diff.Mul(diff))
	}

	variance := sumSquaredDiffs.Div(decimal.NewFromInt(int64(len(returns))))
	return decimal.NewFromFloat(math.Sqrt(variance.InexactFloat64()))
}

// calculatePriceChange calculates price change over specified periods
func (pa *PredictiveAnalytics) calculatePriceChange(prices []decimal.Decimal, periods int) decimal.Decimal {
	if len(prices) <= periods {
		return decimal.NewFromInt(0)
	}

	currentPrice := prices[len(prices)-1]
	pastPrice := prices[len(prices)-1-periods]

	if pastPrice.GreaterThan(decimal.NewFromInt(0)) {
		return currentPrice.Sub(pastPrice).Div(pastPrice)
	}

	return decimal.NewFromInt(0)
}

// encodeRegime encodes market regime as numeric value
func (pa *PredictiveAnalytics) encodeRegime(regime MarketRegime) float64 {
	switch regime {
	case MarketRegimeBull:
		return 1.0
	case MarketRegimeBear:
		return -1.0
	case MarketRegimeVolatile:
		return 0.5
	default: // MarketRegimeSideways
		return 0.0
	}
}

// convertFeaturesToMap converts features to map[string]float64
func (pa *PredictiveAnalytics) convertFeaturesToMap(features map[string]interface{}) map[string]float64 {
	result := make(map[string]float64)

	for key, value := range features {
		if floatValue, ok := value.(float64); ok {
			result[key] = floatValue
		} else {
			result[key] = 0.0 // Default value
		}
	}

	return result
}

// calculateRSI calculates Relative Strength Index
func (pa *PredictiveAnalytics) calculateRSI(prices []decimal.Decimal, period int) decimal.Decimal {
	if len(prices) < period+1 {
		return decimal.NewFromInt(50) // Neutral RSI
	}

	gains := decimal.NewFromInt(0)
	losses := decimal.NewFromInt(0)

	// Calculate initial average gain and loss
	for i := 1; i <= period; i++ {
		change := prices[i].Sub(prices[i-1])
		if change.GreaterThan(decimal.NewFromInt(0)) {
			gains = gains.Add(change)
		} else {
			losses = losses.Add(change.Abs())
		}
	}

	avgGain := gains.Div(decimal.NewFromInt(int64(period)))
	avgLoss := losses.Div(decimal.NewFromInt(int64(period)))

	if avgLoss.IsZero() {
		return decimal.NewFromInt(100)
	}

	rs := avgGain.Div(avgLoss)
	rsi := decimal.NewFromInt(100).Sub(decimal.NewFromInt(100).Div(decimal.NewFromInt(1).Add(rs)))

	return rsi
}

// calculateMACD calculates MACD (Moving Average Convergence Divergence)
func (pa *PredictiveAnalytics) calculateMACD(prices []decimal.Decimal) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	if len(prices) < 26 {
		return decimal.NewFromInt(0), decimal.NewFromInt(0), decimal.NewFromInt(0)
	}

	// Calculate EMAs
	ema12 := pa.calculateEMA(prices, 12)
	ema26 := pa.calculateEMA(prices, 26)

	// MACD line
	macd := ema12.Sub(ema26)

	// Signal line (9-period EMA of MACD)
	// Simplified: use SMA instead of EMA for signal
	signal := macd.Mul(decimal.NewFromFloat(0.9)) // Simplified signal

	// Histogram
	histogram := macd.Sub(signal)

	return macd, signal, histogram
}

// calculateEMA calculates Exponential Moving Average
func (pa *PredictiveAnalytics) calculateEMA(prices []decimal.Decimal, period int) decimal.Decimal {
	if len(prices) < period {
		return decimal.NewFromInt(0)
	}

	// Start with SMA for first value
	sma := pa.calculateSMA(prices[:period], period)
	ema := sma

	// Calculate multiplier
	multiplier := decimal.NewFromInt(2).Div(decimal.NewFromInt(int64(period + 1)))

	// Calculate EMA for remaining values
	for i := period; i < len(prices); i++ {
		ema = prices[i].Mul(multiplier).Add(ema.Mul(decimal.NewFromInt(1).Sub(multiplier)))
	}

	return ema
}

// calculateATR calculates Average True Range
func (pa *PredictiveAnalytics) calculateATR(highs, lows, closes []decimal.Decimal, period int) decimal.Decimal {
	if len(highs) < period+1 || len(lows) < period+1 || len(closes) < period+1 {
		return decimal.NewFromInt(0)
	}

	trueRanges := make([]decimal.Decimal, len(highs)-1)

	for i := 1; i < len(highs); i++ {
		// True Range = max(high-low, |high-prevClose|, |low-prevClose|)
		hl := highs[i].Sub(lows[i])
		hc := highs[i].Sub(closes[i-1]).Abs()
		lc := lows[i].Sub(closes[i-1]).Abs()

		tr := hl
		if hc.GreaterThan(tr) {
			tr = hc
		}
		if lc.GreaterThan(tr) {
			tr = lc
		}

		trueRanges[i-1] = tr
	}

	// Calculate ATR as SMA of True Ranges
	if len(trueRanges) >= period {
		return pa.calculateSMA(trueRanges[len(trueRanges)-period:], period)
	}

	return decimal.NewFromInt(0)
}
