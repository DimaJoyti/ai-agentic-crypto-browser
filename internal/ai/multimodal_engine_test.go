package ai

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMultiModalEngine(t *testing.T) {
	logger := &observability.Logger{}
	engine := NewMultiModalEngine(logger)
	require.NotNil(t, engine)

	t.Run("EngineInitialization", func(t *testing.T) {
		assert.NotNil(t, engine.logger)
		assert.NotNil(t, engine.config)
		assert.NotNil(t, engine.imageProcessor)
		assert.NotNil(t, engine.documentAnalyzer)
		assert.NotNil(t, engine.voiceProcessor)
		assert.NotNil(t, engine.chartAnalyzer)
		assert.NotNil(t, engine.ocrEngine)
		assert.NotNil(t, engine.cache)

		// Check configuration
		assert.Equal(t, int64(10*1024*1024), engine.config.MaxImageSize)
		assert.Equal(t, int64(50*1024*1024), engine.config.MaxDocumentSize)
		assert.Equal(t, 10*time.Minute, engine.config.MaxAudioDuration)
		assert.True(t, engine.config.EnableOCR)
		assert.True(t, engine.config.EnableChartAnalysis)
		assert.True(t, engine.config.EnableVoiceCommands)
		assert.Equal(t, 30*time.Minute, engine.config.CacheTimeout)
		assert.Equal(t, 5*time.Minute, engine.config.ProcessingTimeout)
		assert.True(t, engine.config.ParallelProcessing)

		// Check supported formats
		assert.Contains(t, engine.config.SupportedImageTypes, "jpg")
		assert.Contains(t, engine.config.SupportedImageTypes, "png")
		assert.Contains(t, engine.config.SupportedDocTypes, "pdf")
		assert.Contains(t, engine.config.SupportedAudioTypes, "mp3")
	})

	t.Run("ImageAnalysis", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.New()

		// Create a simple base64 encoded image (1x1 pixel PNG)
		imageData := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChAHGArEkAAAAAElFTkSuQmCC"

		req := &MultiModalRequest{
			RequestID: uuid.New().String(),
			UserID:    userID,
			Type:      "image",
			Content: []MultiModalContent{
				{
					ID:       uuid.New().String(),
					Type:     "image",
					Data:     imageData,
					MimeType: "image/png",
					Filename: "test.png",
					Size:     100,
				},
			},
			Options: MultiModalOptions{
				AnalyzeImages: true,
				ExtractText:   true,
				AnalyzeCharts: true,
				DetectObjects: true,
			},
			RequestedAt: time.Now(),
		}

		result, err := engine.ProcessMultiModalRequest(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Validate basic structure
		assert.Equal(t, req.RequestID, result.RequestID)
		assert.Equal(t, userID, result.UserID)
		assert.Equal(t, "image", result.Type)
		assert.Len(t, result.Results, 1)

		// Validate image analysis result
		imageResult := result.Results[0]
		assert.Equal(t, req.Content[0].ID, imageResult.ContentID)
		assert.Equal(t, "image", imageResult.Type)
		assert.NotNil(t, imageResult.ImageAnalysis)
		assert.NotNil(t, imageResult.ChartAnalysis)
		assert.NotNil(t, imageResult.OCRResult)
		assert.GreaterOrEqual(t, imageResult.Confidence, 0.0)
		assert.LessOrEqual(t, imageResult.Confidence, 1.0)

		// Validate aggregated data
		assert.NotNil(t, result.AggregatedData)
		assert.Equal(t, 1, result.AggregatedData.ProcessingStats.TotalItems)
		assert.Equal(t, 1, result.AggregatedData.ProcessingStats.SuccessfulItems)
		assert.Equal(t, 0, result.AggregatedData.ProcessingStats.FailedItems)
		assert.GreaterOrEqual(t, result.AggregatedData.QualityMetrics.AverageConfidence, 0.0)
	})

	t.Run("DocumentAnalysis", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.New()

		// Create a simple text document
		documentText := "This is a sample financial document with revenue of $1.2B and profit margin of 15%."
		documentData := base64.StdEncoding.EncodeToString([]byte(documentText))

		req := &MultiModalRequest{
			RequestID: uuid.New().String(),
			UserID:    userID,
			Type:      "document",
			Content: []MultiModalContent{
				{
					ID:       uuid.New().String(),
					Type:     "document",
					Data:     documentData,
					MimeType: "text/plain",
					Filename: "financial_report.txt",
					Size:     int64(len(documentText)),
				},
			},
			Options: MultiModalOptions{
				ExtractText:      true,
				AnalyzeSentiment: true,
				ExtractEntities:  true,
				GenerateSummary:  true,
			},
			RequestedAt: time.Now(),
		}

		result, err := engine.ProcessMultiModalRequest(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Validate basic structure
		assert.Equal(t, req.RequestID, result.RequestID)
		assert.Equal(t, userID, result.UserID)
		assert.Equal(t, "document", result.Type)
		assert.Len(t, result.Results, 1)

		// Validate document analysis result
		docResult := result.Results[0]
		assert.Equal(t, req.Content[0].ID, docResult.ContentID)
		assert.Equal(t, "document", docResult.Type)
		assert.NotNil(t, docResult.DocumentAnalysis)
		assert.NotEmpty(t, docResult.ExtractedText)
		assert.GreaterOrEqual(t, docResult.Confidence, 0.0)

		// Validate document analysis details
		docAnalysis := docResult.DocumentAnalysis
		assert.Equal(t, "financial_report", docAnalysis.DocumentType)
		assert.NotEmpty(t, docAnalysis.ExtractedText)
		assert.NotNil(t, docAnalysis.Structure)
		assert.NotEmpty(t, docAnalysis.KeyInformation)
		assert.NotEmpty(t, docAnalysis.TradingInsights)
	})

	t.Run("AudioAnalysis", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.New()

		// Create a simple audio data (mock)
		audioData := base64.StdEncoding.EncodeToString([]byte("mock audio data"))

		req := &MultiModalRequest{
			RequestID: uuid.New().String(),
			UserID:    userID,
			Type:      "audio",
			Content: []MultiModalContent{
				{
					ID:       uuid.New().String(),
					Type:     "audio",
					Data:     audioData,
					MimeType: "audio/wav",
					Filename: "voice_command.wav",
					Size:     1000,
				},
			},
			Options: MultiModalOptions{
				ProcessAudio:     true,
				AnalyzeSentiment: true,
				ExtractEntities:  true,
			},
			RequestedAt: time.Now(),
		}

		result, err := engine.ProcessMultiModalRequest(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Validate basic structure
		assert.Equal(t, req.RequestID, result.RequestID)
		assert.Equal(t, userID, result.UserID)
		assert.Equal(t, "audio", result.Type)
		assert.Len(t, result.Results, 1)

		// Validate audio analysis result
		audioResult := result.Results[0]
		assert.Equal(t, req.Content[0].ID, audioResult.ContentID)
		assert.Equal(t, "audio", audioResult.Type)
		assert.NotNil(t, audioResult.AudioAnalysis)
		assert.GreaterOrEqual(t, audioResult.Confidence, 0.0)

		// Validate audio analysis details
		audioAnalysis := audioResult.AudioAnalysis
		assert.NotNil(t, audioAnalysis.Transcription)
		assert.NotEmpty(t, audioAnalysis.Transcription.Text)
		assert.NotEmpty(t, audioAnalysis.TradingCommands)
		assert.NotNil(t, audioAnalysis.AudioQuality)
	})

	t.Run("MultiContentAnalysis", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.New()

		// Create mixed content
		imageData := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChAHGArEkAAAAAElFTkSuQmCC"
		textData := base64.StdEncoding.EncodeToString([]byte("Financial analysis document"))

		req := &MultiModalRequest{
			RequestID: uuid.New().String(),
			UserID:    userID,
			Type:      "mixed",
			Content: []MultiModalContent{
				{
					ID:       uuid.New().String(),
					Type:     "image",
					Data:     imageData,
					MimeType: "image/png",
					Filename: "chart.png",
					Size:     100,
				},
				{
					ID:       uuid.New().String(),
					Type:     "document",
					Data:     textData,
					MimeType: "text/plain",
					Filename: "analysis.txt",
					Size:     50,
				},
			},
			Options: MultiModalOptions{
				AnalyzeImages:    true,
				ExtractText:      true,
				AnalyzeCharts:    true,
				AnalyzeSentiment: true,
				GenerateSummary:  true,
			},
			RequestedAt: time.Now(),
		}

		result, err := engine.ProcessMultiModalRequest(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Validate basic structure
		assert.Equal(t, req.RequestID, result.RequestID)
		assert.Equal(t, userID, result.UserID)
		assert.Equal(t, "mixed", result.Type)
		assert.Len(t, result.Results, 2)

		// Validate aggregated data
		assert.NotNil(t, result.AggregatedData)
		assert.Equal(t, 2, result.AggregatedData.ProcessingStats.TotalItems)
		assert.GreaterOrEqual(t, result.AggregatedData.ProcessingStats.SuccessfulItems, 1)
		assert.NotEmpty(t, result.AggregatedData.ContentTypes)
		assert.Contains(t, result.AggregatedData.ContentTypes, "image")
		assert.Contains(t, result.AggregatedData.ContentTypes, "document")
	})

	t.Run("RequestValidation", func(t *testing.T) {
		ctx := context.Background()

		// Test missing request ID
		req := &MultiModalRequest{
			UserID: uuid.New(),
			Type:   "image",
			Content: []MultiModalContent{
				{
					ID:   uuid.New().String(),
					Type: "image",
					Data: "test",
					Size: 100,
				},
			},
		}

		_, err := engine.ProcessMultiModalRequest(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "request ID is required")

		// Test missing user ID
		req = &MultiModalRequest{
			RequestID: uuid.New().String(),
			Type:      "image",
			Content: []MultiModalContent{
				{
					ID:   uuid.New().String(),
					Type: "image",
					Data: "test",
					Size: 100,
				},
			},
		}

		_, err = engine.ProcessMultiModalRequest(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user ID is required")

		// Test empty content
		req = &MultiModalRequest{
			RequestID: uuid.New().String(),
			UserID:    uuid.New(),
			Type:      "image",
			Content:   []MultiModalContent{},
		}

		_, err = engine.ProcessMultiModalRequest(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no content provided")

		// Test content with missing type
		req = &MultiModalRequest{
			RequestID: uuid.New().String(),
			UserID:    uuid.New(),
			Type:      "image",
			Content: []MultiModalContent{
				{
					ID:   uuid.New().String(),
					Data: "test",
					Size: 100,
				},
			},
		}

		_, err = engine.ProcessMultiModalRequest(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "type is required")

		// Test content with missing data
		req = &MultiModalRequest{
			RequestID: uuid.New().String(),
			UserID:    uuid.New(),
			Type:      "image",
			Content: []MultiModalContent{
				{
					ID:   uuid.New().String(),
					Type: "image",
					Size: 100,
				},
			},
		}

		_, err = engine.ProcessMultiModalRequest(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "data is required")
	})

	t.Run("ImageFormatValidation", func(t *testing.T) {
		// Test supported formats
		supportedFormats := []string{"test.jpg", "test.jpeg", "test.png", "test.gif", "test.webp"}
		for _, filename := range supportedFormats {
			assert.True(t, engine.ValidateImageFormat(filename), "Format should be supported: %s", filename)
		}

		// Test unsupported formats
		unsupportedFormats := []string{"test.bmp", "test.tiff", "test.svg", "test.pdf", "test.txt"}
		for _, filename := range unsupportedFormats {
			assert.False(t, engine.ValidateImageFormat(filename), "Format should not be supported: %s", filename)
		}
	})

	t.Run("SupportedFormats", func(t *testing.T) {
		formats := engine.GetSupportedFormats()
		require.NotNil(t, formats)

		// Check that all format categories are present
		assert.Contains(t, formats, "images")
		assert.Contains(t, formats, "documents")
		assert.Contains(t, formats, "audio")

		// Check that formats contain expected values
		imageFormats := formats["images"]
		assert.Contains(t, imageFormats, "jpg")
		assert.Contains(t, imageFormats, "png")

		docFormats := formats["documents"]
		assert.Contains(t, docFormats, "pdf")
		assert.Contains(t, docFormats, "txt")

		audioFormats := formats["audio"]
		assert.Contains(t, audioFormats, "mp3")
		assert.Contains(t, audioFormats, "wav")
	})

	t.Run("CacheOperations", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.New()

		imageData := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChAHGArEkAAAAAElFTkSuQmCC"

		req := &MultiModalRequest{
			RequestID: "cache-test-123",
			UserID:    userID,
			Type:      "image",
			Content: []MultiModalContent{
				{
					ID:       uuid.New().String(),
					Type:     "image",
					Data:     imageData,
					MimeType: "image/png",
					Filename: "test.png",
					Size:     100,
				},
			},
			Options: MultiModalOptions{
				AnalyzeImages: true,
			},
			RequestedAt: time.Now(),
		}

		// First request should process normally
		result1, err := engine.ProcessMultiModalRequest(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, result1)

		// Second identical request should return cached result
		result2, err := engine.ProcessMultiModalRequest(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, result2)

		// Results should be identical (from cache)
		assert.Equal(t, result1.RequestID, result2.RequestID)
		assert.Equal(t, result1.UserID, result2.UserID)
		assert.Equal(t, result1.Type, result2.Type)
	})

	t.Run("ParallelProcessing", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.New()

		// Create multiple content items
		imageData := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChAHGArEkAAAAAElFTkSuQmCC"
		textData := base64.StdEncoding.EncodeToString([]byte("Test document"))

		req := &MultiModalRequest{
			RequestID: uuid.New().String(),
			UserID:    userID,
			Type:      "mixed",
			Content: []MultiModalContent{
				{
					ID:       uuid.New().String(),
					Type:     "image",
					Data:     imageData,
					MimeType: "image/png",
					Filename: "test1.png",
					Size:     100,
				},
				{
					ID:       uuid.New().String(),
					Type:     "image",
					Data:     imageData,
					MimeType: "image/png",
					Filename: "test2.png",
					Size:     100,
				},
				{
					ID:       uuid.New().String(),
					Type:     "document",
					Data:     textData,
					MimeType: "text/plain",
					Filename: "test.txt",
					Size:     50,
				},
			},
			Options: MultiModalOptions{
				AnalyzeImages: true,
				ExtractText:   true,
			},
			RequestedAt: time.Now(),
		}

		// Test with parallel processing enabled
		engine.config.ParallelProcessing = true
		result, err := engine.ProcessMultiModalRequest(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.Len(t, result.Results, 3)
		assert.Equal(t, 3, result.AggregatedData.ProcessingStats.TotalItems)

		// Test with parallel processing disabled
		engine.config.ParallelProcessing = false
		result2, err := engine.ProcessMultiModalRequest(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, result2)

		assert.Len(t, result2.Results, 3)
		assert.Equal(t, 3, result2.AggregatedData.ProcessingStats.TotalItems)

		// Reset to default
		engine.config.ParallelProcessing = true
	})
}

func TestImageProcessor(t *testing.T) {
	logger := &observability.Logger{}
	processor := NewImageProcessor(logger)
	require.NotNil(t, processor)

	t.Run("ImageAnalysis", func(t *testing.T) {
		ctx := context.Background()
		content := MultiModalContent{
			ID:       uuid.New().String(),
			Type:     "image",
			Data:     "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChAHGArEkAAAAAElFTkSuQmCC",
			MimeType: "image/png",
			Filename: "test.png",
			Size:     100,
		}

		result, err := processor.AnalyzeImage(ctx, content)
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.NotEmpty(t, result.Objects)
		assert.NotEmpty(t, result.Scenes)
		assert.NotEmpty(t, result.Colors)
		assert.NotNil(t, result.Quality)
		assert.GreaterOrEqual(t, result.Confidence, 0.0)
		assert.LessOrEqual(t, result.Confidence, 1.0)

		// Check object detection
		obj := result.Objects[0]
		assert.NotEmpty(t, obj.Label)
		assert.GreaterOrEqual(t, obj.Confidence, 0.0)
		assert.LessOrEqual(t, obj.Confidence, 1.0)

		// Check scene detection
		scene := result.Scenes[0]
		assert.NotEmpty(t, scene.Label)
		assert.GreaterOrEqual(t, scene.Confidence, 0.0)

		// Check color analysis
		color := result.Colors[0]
		assert.NotEmpty(t, color.Color)
		assert.GreaterOrEqual(t, color.Percentage, 0.0)
		assert.LessOrEqual(t, color.Percentage, 100.0)
		assert.Len(t, color.RGB, 3)

		// Check quality metrics
		assert.GreaterOrEqual(t, result.Quality.Overall, 0.0)
		assert.LessOrEqual(t, result.Quality.Overall, 1.0)
	})
}

func TestDocumentAnalyzer(t *testing.T) {
	logger := &observability.Logger{}
	analyzer := NewDocumentAnalyzer(logger)
	require.NotNil(t, analyzer)

	t.Run("DocumentAnalysis", func(t *testing.T) {
		ctx := context.Background()
		content := MultiModalContent{
			ID:       uuid.New().String(),
			Type:     "document",
			Data:     base64.StdEncoding.EncodeToString([]byte("Sample financial document")),
			MimeType: "text/plain",
			Filename: "report.txt",
			Size:     100,
		}

		options := MultiModalOptions{
			ExtractText:      true,
			AnalyzeSentiment: true,
			ExtractEntities:  true,
			GenerateSummary:  true,
		}

		result, err := analyzer.AnalyzeDocument(ctx, content, options)
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.NotEmpty(t, result.DocumentType)
		assert.NotEmpty(t, result.ExtractedText)
		assert.NotNil(t, result.Structure)
		assert.NotEmpty(t, result.KeyInformation)
		assert.NotEmpty(t, result.TradingInsights)
		assert.GreaterOrEqual(t, result.Confidence, 0.0)

		// Check structure
		assert.NotEmpty(t, result.Structure.Title)
		assert.GreaterOrEqual(t, result.Structure.PageCount, 1)
		assert.GreaterOrEqual(t, result.Structure.WordCount, 1)

		// Check key information
		keyInfo := result.KeyInformation[0]
		assert.NotEmpty(t, keyInfo.Key)
		assert.NotEmpty(t, keyInfo.Value)
		assert.GreaterOrEqual(t, keyInfo.Confidence, 0.0)

		// Check trading insights
		insight := result.TradingInsights[0]
		assert.NotEmpty(t, insight.Type)
		assert.NotEmpty(t, insight.Insight)
		assert.GreaterOrEqual(t, insight.Confidence, 0.0)
	})
}

func TestVoiceProcessor(t *testing.T) {
	logger := &observability.Logger{}
	processor := NewVoiceProcessor(logger)
	require.NotNil(t, processor)

	t.Run("AudioProcessing", func(t *testing.T) {
		ctx := context.Background()
		content := MultiModalContent{
			ID:       uuid.New().String(),
			Type:     "audio",
			Data:     base64.StdEncoding.EncodeToString([]byte("mock audio data")),
			MimeType: "audio/wav",
			Filename: "command.wav",
			Size:     1000,
		}

		options := MultiModalOptions{
			ProcessAudio:     true,
			AnalyzeSentiment: true,
		}

		result, err := processor.ProcessAudio(ctx, content, options)
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.NotNil(t, result.Transcription)
		assert.NotEmpty(t, result.TradingCommands)
		assert.NotNil(t, result.AudioQuality)
		assert.GreaterOrEqual(t, result.Confidence, 0.0)

		// Check transcription
		assert.NotEmpty(t, result.Transcription.Text)
		assert.NotEmpty(t, result.Transcription.Language)
		assert.GreaterOrEqual(t, result.Transcription.Confidence, 0.0)
		assert.NotEmpty(t, result.Transcription.Segments)

		// Check trading commands
		command := result.TradingCommands[0]
		assert.NotEmpty(t, command.Action)
		assert.NotEmpty(t, command.Asset)
		assert.GreaterOrEqual(t, command.Confidence, 0.0)

		// Check audio quality
		assert.GreaterOrEqual(t, result.AudioQuality.Overall, 0.0)
		assert.LessOrEqual(t, result.AudioQuality.Overall, 1.0)
	})
}

func TestChartAnalyzer(t *testing.T) {
	logger := &observability.Logger{}
	analyzer := NewChartAnalyzer(logger)
	require.NotNil(t, analyzer)

	t.Run("ChartAnalysis", func(t *testing.T) {
		ctx := context.Background()
		content := MultiModalContent{
			ID:       uuid.New().String(),
			Type:     "image",
			Data:     "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChAHGArEkAAAAAElFTkSuQmCC",
			MimeType: "image/png",
			Filename: "btc_chart.png",
			Size:     100,
		}

		result, err := analyzer.AnalyzeChart(ctx, content)
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.NotEmpty(t, result.ChartType)
		assert.NotEmpty(t, result.Asset)
		assert.NotEmpty(t, result.TechnicalSignals)
		assert.NotEmpty(t, result.Patterns)
		assert.NotNil(t, result.TrendAnalysis)
		assert.NotNil(t, result.Recommendation)
		assert.GreaterOrEqual(t, result.Confidence, 0.0)

		// Check technical signals
		signal := result.TechnicalSignals[0]
		assert.NotEmpty(t, signal.Indicator)
		assert.NotEmpty(t, signal.Signal)
		assert.GreaterOrEqual(t, signal.Confidence, 0.0)

		// Check patterns
		pattern := result.Patterns[0]
		assert.NotEmpty(t, pattern.Pattern)
		assert.NotEmpty(t, pattern.Type)
		assert.GreaterOrEqual(t, pattern.Confidence, 0.0)

		// Check trend analysis
		assert.NotEmpty(t, result.TrendAnalysis.Direction)
		assert.GreaterOrEqual(t, result.TrendAnalysis.Strength, 0.0)
		assert.GreaterOrEqual(t, result.TrendAnalysis.Confidence, 0.0)

		// Check recommendation
		assert.NotEmpty(t, result.Recommendation.Action)
		assert.GreaterOrEqual(t, result.Recommendation.Confidence, 0.0)
		assert.NotEmpty(t, result.Recommendation.Reasoning)
	})
}

func TestOCREngine(t *testing.T) {
	logger := &observability.Logger{}
	engine := NewOCREngine(logger)
	require.NotNil(t, engine)

	t.Run("TextExtraction", func(t *testing.T) {
		ctx := context.Background()
		content := MultiModalContent{
			ID:       uuid.New().String(),
			Type:     "image",
			Data:     "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChAHGArEkAAAAAElFTkSuQmCC",
			MimeType: "image/png",
			Filename: "text_image.png",
			Size:     100,
		}

		result, err := engine.ExtractText(ctx, content)
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.NotEmpty(t, result.ExtractedText)
		assert.GreaterOrEqual(t, result.Confidence, 0.0)
		assert.NotEmpty(t, result.Language)
		assert.NotEmpty(t, result.TextBlocks)

		// Check text blocks
		block := result.TextBlocks[0]
		assert.NotEmpty(t, block.Text)
		assert.GreaterOrEqual(t, block.Confidence, 0.0)
		assert.GreaterOrEqual(t, block.BoundingBox.Width, 0.0)
		assert.GreaterOrEqual(t, block.BoundingBox.Height, 0.0)
	})
}
