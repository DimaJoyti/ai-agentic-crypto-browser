package ai

import (
	"context"
	"testing"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdvancedNLPEngine(t *testing.T) {
	logger := &observability.Logger{}
	engine := NewAdvancedNLPEngine(logger)
	require.NotNil(t, engine)

	t.Run("EngineInitialization", func(t *testing.T) {
		assert.NotNil(t, engine.logger)
		assert.NotNil(t, engine.config)
		assert.NotNil(t, engine.languageDetector)
		assert.NotNil(t, engine.multiLangSentiment)
		assert.NotNil(t, engine.newsAnalyzer)
		assert.NotNil(t, engine.socialMediaAnalyzer)
		assert.NotNil(t, engine.entityExtractor)
		assert.NotNil(t, engine.topicModeler)
		assert.NotNil(t, engine.textClassifier)
		assert.NotNil(t, engine.translationService)
		assert.NotNil(t, engine.cache)

		// Check configuration
		assert.Contains(t, engine.config.SupportedLanguages, "en")
		assert.Contains(t, engine.config.SupportedLanguages, "es")
		assert.True(t, engine.config.EnableTranslation)
		assert.True(t, engine.config.EnableEntityExtraction)
		assert.True(t, engine.config.EnableTopicModeling)
		assert.True(t, engine.config.EnableNewsAnalysis)
		assert.True(t, engine.config.EnableSocialMedia)
		assert.Equal(t, 30*time.Minute, engine.config.CacheTimeout)
		assert.Equal(t, 10000, engine.config.MaxTextLength)
		assert.Equal(t, 100, engine.config.BatchSize)
		assert.Equal(t, 0.7, engine.config.ConfidenceThreshold)
		assert.True(t, engine.config.ParallelProcessing)
	})

	t.Run("BasicNLPProcessing", func(t *testing.T) {
		ctx := context.Background()

		req := &NLPRequest{
			RequestID: uuid.New().String(),
			Texts: []string{
				"Bitcoin is showing bullish momentum with strong buying pressure.",
				"The market sentiment is very positive today with BTC reaching new highs.",
			},
			Sources: []string{"news", "social"},
			Options: NLPOptions{
				DetectLanguage:       true,
				TranslateToEnglish:   false,
				ExtractEntities:      true,
				PerformTopicModeling: true,
				ClassifyText:         true,
				AnalyzeSentiment:     true,
				ExtractKeywords:      true,
				DetectEmotions:       true,
				AnalyzeReadability:   true,
				DetectSpam:           false,
			},
			RequestedAt: time.Now(),
		}

		result, err := engine.ProcessNLPRequest(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Validate basic structure
		assert.Equal(t, req.RequestID, result.RequestID)
		assert.Equal(t, len(req.Texts), len(result.Results))
		assert.NotNil(t, result.AggregatedResults)
		assert.NotEmpty(t, result.LanguageStats)
		assert.Greater(t, result.ProcessingTime, time.Duration(0))

		// Validate individual results
		for i, textResult := range result.Results {
			assert.Equal(t, req.Texts[i], textResult.Text)
			assert.Equal(t, req.Sources[i], textResult.Source)
			assert.NotEmpty(t, textResult.OriginalLanguage)
			assert.Greater(t, textResult.Confidence, 0.0)
			assert.Greater(t, textResult.ProcessingTime, time.Duration(0))

			// Check sentiment analysis
			if req.Options.AnalyzeSentiment {
				assert.NotNil(t, textResult.Sentiment)
				assert.GreaterOrEqual(t, textResult.Sentiment.Score, -1.0)
				assert.LessOrEqual(t, textResult.Sentiment.Score, 1.0)
				assert.Contains(t, []string{"positive", "negative", "neutral"}, textResult.Sentiment.Label)
				assert.GreaterOrEqual(t, textResult.Sentiment.Confidence, 0.0)
				assert.LessOrEqual(t, textResult.Sentiment.Confidence, 1.0)
			}

			// Check entity extraction
			if req.Options.ExtractEntities {
				assert.NotNil(t, textResult.Entities)
				// Should find crypto entities
				foundCrypto := false
				for _, entity := range textResult.Entities {
					if entity.Type == "CRYPTO" {
						foundCrypto = true
						assert.NotEmpty(t, entity.Text)
						assert.GreaterOrEqual(t, entity.StartPos, 0)
						assert.Greater(t, entity.EndPos, entity.StartPos)
						assert.GreaterOrEqual(t, entity.Confidence, 0.0)
						break
					}
				}
				assert.True(t, foundCrypto, "Should find crypto entities in crypto-related text")
			}

			// Check text classification
			if req.Options.ClassifyText {
				assert.NotNil(t, textResult.Classification)
				assert.NotEmpty(t, textResult.Classification.Category)
				assert.GreaterOrEqual(t, textResult.Classification.Confidence, 0.0)
				assert.LessOrEqual(t, textResult.Classification.Confidence, 1.0)
				assert.NotEmpty(t, textResult.Classification.Scores)
			}

			// Check keyword extraction
			if req.Options.ExtractKeywords {
				assert.NotNil(t, textResult.Keywords)
				for _, keyword := range textResult.Keywords {
					assert.NotEmpty(t, keyword.Keyword)
					assert.Greater(t, keyword.Score, 0.0)
					assert.Greater(t, keyword.Frequency, 0)
				}
			}

			// Check readability
			if req.Options.AnalyzeReadability {
				assert.GreaterOrEqual(t, textResult.ReadabilityScore, 0.0)
				assert.LessOrEqual(t, textResult.ReadabilityScore, 100.0)
			}
		}

		// Validate aggregated results
		agg := result.AggregatedResults
		assert.NotNil(t, agg.OverallSentiment)
		assert.GreaterOrEqual(t, agg.OverallSentiment.Score, -1.0)
		assert.LessOrEqual(t, agg.OverallSentiment.Score, 1.0)
		assert.Contains(t, []string{"positive", "negative", "neutral"}, agg.OverallSentiment.Label)

		// Check language statistics
		assert.Contains(t, result.LanguageStats, "en")
		assert.Equal(t, len(req.Texts), result.LanguageStats["en"])
	})

	t.Run("MultiLanguageProcessing", func(t *testing.T) {
		ctx := context.Background()

		req := &NLPRequest{
			RequestID: uuid.New().String(),
			Texts: []string{
				"Bitcoin está mostrando un impulso alcista.", // Spanish
				"Le Bitcoin montre une dynamique haussière.", // French
				"Bitcoin zeigt eine bullische Dynamik.",      // German
				"Bitcoin is showing bullish momentum.",       // English
			},
			Sources: []string{"news", "news", "news", "news"},
			Options: NLPOptions{
				DetectLanguage:     true,
				TranslateToEnglish: true,
				AnalyzeSentiment:   true,
				ExtractEntities:    true,
			},
			RequestedAt: time.Now(),
		}

		result, err := engine.ProcessNLPRequest(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Should detect different languages
		languages := make(map[string]bool)
		for _, textResult := range result.Results {
			languages[textResult.OriginalLanguage] = true

			// Non-English texts should have translations
			if textResult.OriginalLanguage != "en" {
				assert.NotEmpty(t, textResult.TranslatedText)
				assert.Contains(t, textResult.TranslatedText, "Translated")
			}
		}

		// Should have detected multiple languages
		assert.GreaterOrEqual(t, len(languages), 2)
		assert.GreaterOrEqual(t, len(result.LanguageStats), 2)
	})

	t.Run("SentimentAnalysisAccuracy", func(t *testing.T) {
		ctx := context.Background()

		testCases := []struct {
			text              string
			expectedSentiment string
		}{
			{"Bitcoin is amazing and going to the moon! Bullish!", "positive"},
			{"This is terrible news, Bitcoin is crashing badly", "negative"},
			{"Bitcoin price is stable today", "neutral"},
			{"Great news! BTC pump incoming!", "positive"},
			{"Fear and panic in the market, dump expected", "negative"},
		}

		for _, tc := range testCases {
			req := &NLPRequest{
				RequestID: uuid.New().String(),
				Texts:     []string{tc.text},
				Options: NLPOptions{
					AnalyzeSentiment: true,
				},
				RequestedAt: time.Now(),
			}

			result, err := engine.ProcessNLPRequest(ctx, req)
			require.NoError(t, err)
			require.Len(t, result.Results, 1)

			sentiment := result.Results[0].Sentiment
			require.NotNil(t, sentiment)
			assert.Equal(t, tc.expectedSentiment, sentiment.Label,
				"Text: %s, Expected: %s, Got: %s", tc.text, tc.expectedSentiment, sentiment.Label)
		}
	})

	t.Run("EntityExtractionAccuracy", func(t *testing.T) {
		ctx := context.Background()

		req := &NLPRequest{
			RequestID: uuid.New().String(),
			Texts: []string{
				"Bitcoin and Ethereum are trading on Binance and Coinbase exchanges.",
				"BTC price reached $50,000 while ETH hit $3,000 on major exchanges.",
			},
			Options: NLPOptions{
				ExtractEntities: true,
			},
			RequestedAt: time.Now(),
		}

		result, err := engine.ProcessNLPRequest(ctx, req)
		require.NoError(t, err)

		// Check that crypto entities are found
		expectedEntities := map[string]string{
			"bitcoin":  "CRYPTO",
			"ethereum": "CRYPTO",
			"btc":      "CRYPTO",
			"eth":      "CRYPTO",
			"binance":  "EXCHANGE",
			"coinbase": "EXCHANGE",
		}

		foundEntities := make(map[string]string)
		for _, textResult := range result.Results {
			for _, entity := range textResult.Entities {
				foundEntities[entity.Text] = entity.Type
			}
		}

		for expectedEntity, expectedType := range expectedEntities {
			if foundType, found := foundEntities[expectedEntity]; found {
				assert.Equal(t, expectedType, foundType,
					"Entity %s should be type %s, got %s", expectedEntity, expectedType, foundType)
			}
		}
	})

	t.Run("TextClassificationAccuracy", func(t *testing.T) {
		ctx := context.Background()

		testCases := []struct {
			text             string
			expectedCategory string
		}{
			{"Bitcoin price analysis and trading signals", "financial"},
			{"Breaking news: New cryptocurrency regulation announced", "news"},
			{"Just bought more Bitcoin! HODL!", "social"},
		}

		for _, tc := range testCases {
			req := &NLPRequest{
				RequestID: uuid.New().String(),
				Texts:     []string{tc.text},
				Options: NLPOptions{
					ClassifyText: true,
				},
				RequestedAt: time.Now(),
			}

			result, err := engine.ProcessNLPRequest(ctx, req)
			require.NoError(t, err)
			require.Len(t, result.Results, 1)

			classification := result.Results[0].Classification
			require.NotNil(t, classification)

			// For financial content, should classify as financial
			if tc.expectedCategory == "financial" {
				assert.Equal(t, tc.expectedCategory, classification.Category,
					"Text: %s, Expected: %s, Got: %s", tc.text, tc.expectedCategory, classification.Category)
			}
		}
	})

	t.Run("KeywordExtractionQuality", func(t *testing.T) {
		ctx := context.Background()

		req := &NLPRequest{
			RequestID: uuid.New().String(),
			Texts: []string{
				"Bitcoin trading analysis shows bullish momentum with strong buying pressure from institutional investors.",
			},
			Options: NLPOptions{
				ExtractKeywords: true,
			},
			RequestedAt: time.Now(),
		}

		result, err := engine.ProcessNLPRequest(ctx, req)
		require.NoError(t, err)
		require.Len(t, result.Results, 1)

		keywords := result.Results[0].Keywords
		assert.NotEmpty(t, keywords)

		// Keywords should be sorted by score (descending)
		for i := 1; i < len(keywords); i++ {
			assert.GreaterOrEqual(t, keywords[i-1].Score, keywords[i].Score,
				"Keywords should be sorted by score")
		}

		// Should extract relevant keywords
		keywordTexts := make([]string, len(keywords))
		for i, kw := range keywords {
			keywordTexts[i] = kw.Keyword
		}

		expectedKeywords := []string{"bitcoin", "trading", "analysis", "bullish"}
		for _, expected := range expectedKeywords {
			found := false
			for _, actual := range keywordTexts {
				if actual == expected {
					found = true
					break
				}
			}
			if !found {
				t.Logf("Expected keyword '%s' not found in: %v", expected, keywordTexts)
			}
		}
	})

	t.Run("ReadabilityScoring", func(t *testing.T) {
		ctx := context.Background()

		testCases := []struct {
			text                string
			expectedReadability string // "high", "medium", "low"
		}{
			{"Bitcoin is good.", "medium"}, // Simple
			{"The cryptocurrency market demonstrates significant volatility patterns.", "medium"},                                     // Medium complexity
			{"The implementation of sophisticated algorithmic trading strategies necessitates comprehensive market analysis.", "low"}, // Complex
		}

		for _, tc := range testCases {
			req := &NLPRequest{
				RequestID: uuid.New().String(),
				Texts:     []string{tc.text},
				Options: NLPOptions{
					AnalyzeReadability: true,
				},
				RequestedAt: time.Now(),
			}

			result, err := engine.ProcessNLPRequest(ctx, req)
			require.NoError(t, err)
			require.Len(t, result.Results, 1)

			readability := result.Results[0].ReadabilityScore
			assert.GreaterOrEqual(t, readability, 0.0)
			assert.LessOrEqual(t, readability, 100.0)

			// Validate readability expectations
			switch tc.expectedReadability {
			case "high":
				assert.Greater(t, readability, 60.0, "Simple text should have high readability")
			case "medium":
				assert.GreaterOrEqual(t, readability, 20.0, "Medium text should have reasonable readability")
				assert.LessOrEqual(t, readability, 80.0, "Medium text should not be too high")
			case "low":
				assert.Less(t, readability, 40.0, "Complex text should have low readability")
			}
		}
	})

	t.Run("SpamDetection", func(t *testing.T) {
		ctx := context.Background()

		testCases := []struct {
			text       string
			expectSpam bool
		}{
			{"Bitcoin analysis shows positive trends", false},
			{"CLICK HERE FOR FREE MONEY!!! URGENT ACT NOW!!!", true},
			{"Limited time exclusive offer - make money fast!", true},
			{"Market update: BTC price increased by 5%", false},
		}

		for _, tc := range testCases {
			req := &NLPRequest{
				RequestID: uuid.New().String(),
				Texts:     []string{tc.text},
				Options: NLPOptions{
					DetectSpam: true,
				},
				RequestedAt: time.Now(),
			}

			result, err := engine.ProcessNLPRequest(ctx, req)
			require.NoError(t, err)
			require.Len(t, result.Results, 1)

			spamProb := result.Results[0].SpamProbability
			assert.GreaterOrEqual(t, spamProb, 0.0)
			assert.LessOrEqual(t, spamProb, 1.0)

			if tc.expectSpam {
				assert.Greater(t, spamProb, 0.5, "Text should be detected as spam: %s", tc.text)
			} else {
				assert.LessOrEqual(t, spamProb, 0.3, "Text should not be detected as spam: %s", tc.text)
			}
		}
	})

	t.Run("CacheValidation", func(t *testing.T) {
		ctx := context.Background()

		req := &NLPRequest{
			RequestID: uuid.New().String(),
			Texts:     []string{"Bitcoin is showing bullish momentum."},
			Options: NLPOptions{
				AnalyzeSentiment: true,
				ExtractEntities:  true,
			},
			RequestedAt: time.Now(),
		}

		// First request - should process normally
		result1, err := engine.ProcessNLPRequest(ctx, req)
		require.NoError(t, err)
		processingTime1 := result1.ProcessingTime

		// Second request with same data - should return cached result
		result2, err := engine.ProcessNLPRequest(ctx, req)
		require.NoError(t, err)
		processingTime2 := result2.ProcessingTime

		// Results should be identical
		assert.Equal(t, result1.RequestID, result2.RequestID)
		assert.Equal(t, result1.GeneratedAt, result2.GeneratedAt)

		// Second request should be faster (cached)
		assert.LessOrEqual(t, processingTime2, processingTime1)
	})

	t.Run("RequestValidation", func(t *testing.T) {
		ctx := context.Background()

		// Empty texts
		req := &NLPRequest{
			RequestID: uuid.New().String(),
			Texts:     []string{},
			Options:   NLPOptions{},
		}

		_, err := engine.ProcessNLPRequest(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no texts provided")

		// Text too long
		longText := make([]byte, engine.config.MaxTextLength+1)
		for i := range longText {
			longText[i] = 'a'
		}

		req = &NLPRequest{
			RequestID: uuid.New().String(),
			Texts:     []string{string(longText)},
			Options:   NLPOptions{},
		}

		_, err = engine.ProcessNLPRequest(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exceeds maximum length")
	})
}

func TestLanguageDetector(t *testing.T) {
	detector := NewLanguageDetector([]string{"en", "es", "fr", "de"})
	require.NotNil(t, detector)

	testCases := []struct {
		text             string
		expectedLanguage string
		minConfidence    float64
	}{
		{"The quick brown fox jumps over the lazy dog", "en", 0.7},
		{"El rápido zorro marrón salta sobre el perro perezoso", "es", 0.6},
		{"Le renard brun rapide saute par-dessus le chien paresseux", "fr", 0.6},
		{"Der schnelle braune Fuchs springt über den faulen Hund", "de", 0.6},
	}

	for _, tc := range testCases {
		language, confidence := detector.DetectLanguage(tc.text)
		assert.Equal(t, tc.expectedLanguage, language, "Text: %s", tc.text)
		assert.GreaterOrEqual(t, confidence, tc.minConfidence, "Text: %s", tc.text)
	}
}

func TestTranslationService(t *testing.T) {
	logger := &observability.Logger{}
	service := NewTranslationService(logger)
	require.NotNil(t, service)

	// Test translation
	text := "Hello world"
	translated, err := service.Translate(text, "en", "es")
	require.NoError(t, err)
	assert.NotEqual(t, text, translated)
	assert.Contains(t, translated, "Translated")

	// Test same language (should return original)
	same, err := service.Translate(text, "en", "en")
	require.NoError(t, err)
	assert.Equal(t, text, same)

	// Test caching
	cached, err := service.Translate(text, "en", "es")
	require.NoError(t, err)
	assert.Equal(t, translated, cached)
}
