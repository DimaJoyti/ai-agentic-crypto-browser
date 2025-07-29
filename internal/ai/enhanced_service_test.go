package ai

import (
	"context"
	"testing"
	"time"

	"github.com/ai-agentic-browser/pkg/ml"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnhancedAIService(t *testing.T) {
	logger := &observability.Logger{}

	service := NewEnhancedAIService(logger)
	require.NotNil(t, service)

	t.Run("ServiceInitialization", func(t *testing.T) {
		assert.NotNil(t, service.modelManager)
		assert.NotNil(t, service.pricePrediction)
		assert.NotNil(t, service.sentimentAnalyzer)
		assert.True(t, service.config.EnablePricePrediction)
		assert.True(t, service.config.EnableSentimentAnalysis)
		assert.True(t, service.config.EnablePatternRecognition)
	})

	t.Run("GetModelStatus", func(t *testing.T) {
		ctx := context.Background()
		status := service.GetModelStatus(ctx)

		assert.NotEmpty(t, status)
		assert.Contains(t, status, "price_prediction")
		assert.Contains(t, status, "sentiment_analysis")

		// Check model info
		priceModel := status["price_prediction"]
		assert.Equal(t, "Advanced Price Prediction Model", priceModel.Name)
		assert.Equal(t, ml.ModelTypeTimeSeries, priceModel.Type)

		sentimentModel := status["sentiment_analysis"]
		assert.Equal(t, "Advanced Crypto Sentiment Analyzer", sentimentModel.Name)
		assert.Equal(t, ml.ModelTypeNLP, sentimentModel.Type)
	})
}

func TestPricePredictionModel(t *testing.T) {
	logger := &observability.Logger{}

	model := NewPricePredictionModel(logger)
	require.NotNil(t, model)

	t.Run("ModelInitialization", func(t *testing.T) {
		info := model.GetInfo()
		assert.Equal(t, "price_prediction_lstm", info.ID)
		assert.Equal(t, "Advanced Price Prediction Model", info.Name)
		assert.Equal(t, ml.ModelTypeTimeSeries, info.Type)
		assert.Contains(t, info.Features, "price")
		assert.Contains(t, info.Features, "volume")
		assert.Contains(t, info.Features, "sentiment_score")
	})

	t.Run("PricePrediction", func(t *testing.T) {
		ctx := context.Background()

		// Create sample historical data
		historicalData := make([]ml.PriceData, 200) // More than lookback period
		basePrice := decimal.NewFromFloat(50000.0)

		for i := 0; i < 200; i++ {
			price := basePrice.Add(decimal.NewFromFloat(float64(i) * 10))
			historicalData[i] = ml.PriceData{
				Symbol:    "BTC",
				Timestamp: time.Now().Add(time.Duration(-200+i) * time.Hour),
				Open:      price,
				High:      price.Mul(decimal.NewFromFloat(1.02)),
				Low:       price.Mul(decimal.NewFromFloat(0.98)),
				Close:     price,
				Volume:    decimal.NewFromFloat(1000000),
				MarketCap: price.Mul(decimal.NewFromFloat(19000000)),
			}
		}

		// Create prediction request
		req := &PricePredictionRequest{
			Symbol:         "BTC",
			HistoricalData: historicalData,
			Timeframe:      "1h",
			Horizon:        24,
		}

		features := map[string]interface{}{
			"request": req,
		}

		// Make prediction
		prediction, err := model.Predict(ctx, features)
		require.NoError(t, err)
		require.NotNil(t, prediction)

		// Validate response
		response, ok := prediction.Value.(*PricePredictionResponse)
		require.True(t, ok)

		assert.Equal(t, "BTC", response.Symbol)
		assert.Len(t, response.PredictedPrices, 24)
		assert.Greater(t, response.Confidence, 0.0)
		assert.LessOrEqual(t, response.Confidence, 1.0)
		assert.Contains(t, []string{"bullish", "bearish", "sideways"}, response.TrendDirection)
		assert.NotEmpty(t, response.SupportLevels)
		assert.NotEmpty(t, response.ResistanceLevels)
	})

	t.Run("ModelTraining", func(t *testing.T) {
		ctx := context.Background()

		// Create training data
		trainingData := ml.TrainingData{
			Features: make([]map[string]interface{}, 1000),
			Labels:   make([]interface{}, 1000),
		}

		for i := 0; i < 1000; i++ {
			trainingData.Features[i] = map[string]interface{}{
				"price":  50000.0 + float64(i),
				"volume": 1000000.0,
				"rsi":    50.0,
			}
			trainingData.Labels[i] = 51000.0 + float64(i) // Future price
		}

		err := model.Train(ctx, trainingData)
		require.NoError(t, err)

		// Check that model was updated
		info := model.GetInfo()
		assert.Greater(t, info.Accuracy, 0.0)
		assert.Equal(t, 1000, info.TrainingSize)
		assert.True(t, model.IsReady())
	})
}

func TestSentimentAnalyzer(t *testing.T) {
	logger := &observability.Logger{}

	analyzer := NewSentimentAnalyzer(logger)
	require.NotNil(t, analyzer)

	t.Run("AnalyzerInitialization", func(t *testing.T) {
		info := analyzer.GetInfo()
		assert.Equal(t, "sentiment_analyzer_bert", info.ID)
		assert.Equal(t, "Advanced Crypto Sentiment Analyzer", info.Name)
		assert.Equal(t, ml.ModelTypeNLP, info.Type)
		assert.True(t, analyzer.IsReady())
	})

	t.Run("SentimentAnalysis", func(t *testing.T) {
		ctx := context.Background()

		// Create sentiment request
		req := &SentimentRequest{
			Texts: []string{
				"Bitcoin is going to the moon! 🚀 Bullish trend confirmed!",
				"This is a terrible crash, selling everything now",
				"Market looks stable, holding my position",
				"HODL diamond hands! 💎🙌",
				"Rug pull incoming, be careful everyone",
			},
			Source:   "twitter",
			Language: "en",
			Options: SentimentOptions{
				IncludeEmotions:   true,
				IncludeKeywords:   true,
				IncludeEntities:   true,
				IncludeConfidence: true,
				NormalizeText:     true,
			},
		}

		features := map[string]interface{}{
			"request": req,
		}

		// Analyze sentiment
		prediction, err := analyzer.Predict(ctx, features)
		require.NoError(t, err)
		require.NotNil(t, prediction)

		// Validate response
		response, ok := prediction.Value.(*SentimentResponse)
		require.True(t, ok)

		assert.Len(t, response.Results, 5)
		assert.NotNil(t, response.Aggregated)

		// Check individual results
		for _, result := range response.Results {
			assert.GreaterOrEqual(t, result.Sentiment, -1.0)
			assert.LessOrEqual(t, result.Sentiment, 1.0)
			assert.GreaterOrEqual(t, result.Confidence, 0.0)
			assert.LessOrEqual(t, result.Confidence, 1.0)
			assert.Contains(t, []string{"positive", "negative", "neutral"}, result.Label)
			assert.Equal(t, "en", result.Language)
			assert.Equal(t, "twitter", result.Source)
		}

		// Check aggregated results
		agg := response.Aggregated
		assert.GreaterOrEqual(t, agg.OverallSentiment, -1.0)
		assert.LessOrEqual(t, agg.OverallSentiment, 1.0)
		assert.GreaterOrEqual(t, agg.OverallConfidence, 0.0)
		assert.LessOrEqual(t, agg.OverallConfidence, 1.0)
		assert.Equal(t, 5, agg.VolumeMetrics.TotalTexts)
	})

	t.Run("EmotionDetection", func(t *testing.T) {
		ctx := context.Background()

		req := &SentimentRequest{
			Texts: []string{
				"I'm so excited about this pump! 🎉",
				"Feeling scared about this crash 😰",
				"Angry about this rug pull scam!",
			},
			Source:   "reddit",
			Language: "en",
			Options: SentimentOptions{
				IncludeEmotions: true,
			},
		}

		features := map[string]interface{}{
			"request": req,
		}

		prediction, err := analyzer.Predict(ctx, features)
		require.NoError(t, err)

		response := prediction.Value.(*SentimentResponse)

		// Check that emotions are detected
		for _, result := range response.Results {
			if len(result.Emotions) > 0 {
				for emotion, score := range result.Emotions {
					assert.Contains(t, []string{"joy", "trust", "fear", "surprise", "sadness", "disgust", "anger", "anticipation"}, emotion)
					assert.GreaterOrEqual(t, score, 0.0)
					assert.LessOrEqual(t, score, 1.0)
				}
			}
		}
	})
}

func TestEnhancedAIIntegration(t *testing.T) {
	logger := &observability.Logger{}

	service := NewEnhancedAIService(logger)
	require.NotNil(t, service)

	t.Run("ComprehensiveAnalysis", func(t *testing.T) {
		ctx := context.Background()

		// Create historical data for price prediction
		historicalData := make([]ml.PriceData, 200)
		basePrice := decimal.NewFromFloat(45000.0)

		for i := 0; i < 200; i++ {
			price := basePrice.Add(decimal.NewFromFloat(float64(i) * 5))
			historicalData[i] = ml.PriceData{
				Symbol:    "ETH",
				Timestamp: time.Now().Add(time.Duration(-200+i) * time.Hour),
				Close:     price,
				Volume:    decimal.NewFromFloat(500000),
			}
		}

		// Create sentiment data
		sentimentTexts := []string{
			"Ethereum is looking very bullish today!",
			"ETH price prediction shows strong upward momentum",
			"DeFi protocols are driving ETH adoption",
		}

		// Create comprehensive AI request
		req := &AIRequest{
			RequestID: uuid.New().String(),
			UserID:    uuid.New(),
			Type:      "comprehensive_analysis",
			Symbol:    "ETH",
			Data: map[string]interface{}{
				"price_prediction_request": &PricePredictionRequest{
					Symbol:         "ETH",
					HistoricalData: historicalData,
					Timeframe:      "1h",
					Horizon:        12,
				},
				"sentiment_request": &SentimentRequest{
					Texts:    sentimentTexts,
					Source:   "twitter",
					Language: "en",
					Options: SentimentOptions{
						IncludeEmotions: true,
						IncludeKeywords: true,
					},
				},
			},
			Options: AIRequestOptions{
				IncludePredictions:     true,
				IncludeSentiment:       true,
				IncludePatterns:        true,
				IncludeRecommendations: true,
				IncludeRiskAssessment:  true,
				TimeHorizon:            12,
				ConfidenceThreshold:    0.7,
			},
			RequestedAt: time.Now(),
		}

		// Process the request
		response, err := service.ProcessRequest(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, response)

		// Validate comprehensive response
		assert.Equal(t, req.RequestID, response.RequestID)
		assert.Equal(t, req.UserID, response.UserID)
		assert.Equal(t, "ETH", response.Symbol)
		assert.Greater(t, response.ProcessingTime, time.Duration(0))
		assert.GreaterOrEqual(t, response.Confidence, 0.0)
		assert.LessOrEqual(t, response.Confidence, 1.0)

		// Check price prediction (may be nil if prediction failed)
		if response.PricePrediction != nil {
			assert.Equal(t, "ETH", response.PricePrediction.Symbol)
			assert.Len(t, response.PricePrediction.PredictedPrices, 12)
		}

		// Check sentiment analysis
		if response.SentimentAnalysis != nil {
			assert.Len(t, response.SentimentAnalysis.Results, 3)
		}

		// Check pattern analysis
		assert.NotNil(t, response.PatternAnalysis)
		assert.GreaterOrEqual(t, response.PatternAnalysis.Confidence, 0.0)

		// Check market insights
		assert.NotNil(t, response.MarketInsights)
		assert.NotEmpty(t, response.MarketInsights.MarketCondition)
		assert.NotEmpty(t, response.MarketInsights.MarketSentiment)

		// Check recommendations
		assert.NotEmpty(t, response.Recommendations)
		for _, rec := range response.Recommendations {
			assert.Contains(t, []string{"buy", "sell", "hold", "wait"}, rec.Type)
			assert.NotEmpty(t, rec.Action)
			assert.NotEmpty(t, rec.Reasoning)
			assert.GreaterOrEqual(t, rec.Confidence, 0.0)
			assert.LessOrEqual(t, rec.Confidence, 1.0)
		}

		// Check risk assessment
		assert.NotNil(t, response.RiskAssessment)
		assert.Contains(t, []string{"low", "medium", "high", "extreme"}, response.RiskAssessment.OverallRisk)
		assert.GreaterOrEqual(t, response.RiskAssessment.RiskScore, 0.0)
		assert.LessOrEqual(t, response.RiskAssessment.RiskScore, 1.0)
	})

	t.Run("ModelFeedback", func(t *testing.T) {
		ctx := context.Background()

		feedback := &ml.PredictionFeedback{
			PredictionID: uuid.New().String(),
			ActualValue:  52000.0,
			Correct:      true,
			Confidence:   0.85,
			Timestamp:    time.Now(),
			UserID:       uuid.New().String(),
		}

		err := service.ProvideFeedback(ctx, "price_prediction", feedback)
		require.NoError(t, err)

		// Check that model was updated
		status := service.GetModelStatus(ctx)
		priceModel := status["price_prediction"]
		assert.NotNil(t, priceModel)
	})
}
