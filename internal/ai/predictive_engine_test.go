package ai

import (
	"context"
	"testing"
	"time"

	"github.com/ai-agentic-browser/pkg/ml"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPredictiveEngine(t *testing.T) {
	logger := &observability.Logger{}
	pricePrediction := NewPricePredictionModel(logger)
	sentimentAnalyzer := NewSentimentAnalyzer(logger)

	engine := NewPredictiveEngine(logger, pricePrediction, sentimentAnalyzer)
	require.NotNil(t, engine)

	t.Run("EngineInitialization", func(t *testing.T) {
		assert.NotNil(t, engine.logger)
		assert.NotNil(t, engine.pricePrediction)
		assert.NotNil(t, engine.sentimentAnalyzer)
		assert.NotNil(t, engine.config)
		assert.NotNil(t, engine.cache)
		assert.Equal(t, 15*time.Minute, engine.config.UpdateInterval)
		assert.Equal(t, 30*time.Minute, engine.config.CacheTimeout)
		assert.Equal(t, 100, engine.config.MinDataPoints)
	})

	t.Run("TrendForecast", func(t *testing.T) {
		ctx := context.Background()

		// Create sample historical data
		historicalData := make(map[string][]ml.PriceData)
		symbols := []string{"BTC", "ETH"}

		for _, symbol := range symbols {
			data := make([]ml.PriceData, 200)
			basePrice := decimal.NewFromFloat(50000.0)
			if symbol == "ETH" {
				basePrice = decimal.NewFromFloat(3000.0)
			}

			for i := 0; i < 200; i++ {
				price := basePrice.Add(decimal.NewFromFloat(float64(i) * 10))
				data[i] = ml.PriceData{
					Symbol:    symbol,
					Timestamp: time.Now().Add(time.Duration(-200+i) * time.Hour),
					Open:      price,
					High:      price.Mul(decimal.NewFromFloat(1.02)),
					Low:       price.Mul(decimal.NewFromFloat(0.98)),
					Close:     price,
					Volume:    decimal.NewFromFloat(1000000),
				}
			}
			historicalData[symbol] = data
		}

		req := &PredictiveRequest{
			Symbols:        symbols,
			TimeHorizon:    24,
			AnalysisType:   "trend",
			HistoricalData: historicalData,
			Options: PredictiveOptions{
				IncludeTrendAnalysis: true,
			},
			RequestedAt: time.Now(),
		}

		result, err := engine.GeneratePredictiveAnalytics(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Validate trend analysis
		assert.NotNil(t, result.TrendAnalysis)
		assert.Equal(t, symbols, result.Symbols)
		assert.Equal(t, 24, result.TimeHorizon)
		assert.GreaterOrEqual(t, result.Confidence, 0.0)
		assert.LessOrEqual(t, result.Confidence, 1.0)

		// Check trend predictions for each symbol
		for _, symbol := range symbols {
			trend, exists := result.TrendAnalysis.Trends[symbol]
			assert.True(t, exists)
			assert.Equal(t, symbol, trend.Symbol)
			assert.Contains(t, []string{"up", "down", "sideways"}, trend.Direction)
			assert.GreaterOrEqual(t, trend.Strength, 0.0)
			assert.LessOrEqual(t, trend.Strength, 1.0)
			assert.GreaterOrEqual(t, trend.Probability, 0.0)
			assert.LessOrEqual(t, trend.Probability, 1.0)
			assert.NotEmpty(t, trend.Duration)
			assert.NotEmpty(t, trend.PriceTargets)
		}

		// Check market trend
		assert.Contains(t, []string{"bullish", "bearish", "sideways"}, result.TrendAnalysis.MarketTrend)
		assert.GreaterOrEqual(t, result.TrendAnalysis.TrendStrength, 0.0)
		assert.LessOrEqual(t, result.TrendAnalysis.TrendStrength, 1.0)

		// Check support and resistance levels
		assert.NotEmpty(t, result.TrendAnalysis.SupportLevels)
		assert.NotEmpty(t, result.TrendAnalysis.ResistanceLevels)
	})

	t.Run("VolatilityForecast", func(t *testing.T) {
		ctx := context.Background()

		// Create sample data with varying volatility
		historicalData := make(map[string][]ml.PriceData)
		symbols := []string{"BTC"}

		data := make([]ml.PriceData, 100)
		basePrice := 50000.0

		for i := 0; i < 100; i++ {
			// Add some volatility to the price
			volatility := 0.02 * (1 + 0.5*float64(i%10)/10.0)
			priceChange := volatility * (float64(i%2)*2 - 1) // Alternating up/down
			price := decimal.NewFromFloat(basePrice * (1 + priceChange))

			data[i] = ml.PriceData{
				Symbol:    "BTC",
				Timestamp: time.Now().Add(time.Duration(-100+i) * time.Hour),
				Close:     price,
				Volume:    decimal.NewFromFloat(1000000),
			}
		}
		historicalData["BTC"] = data

		req := &PredictiveRequest{
			Symbols:        symbols,
			TimeHorizon:    12,
			AnalysisType:   "volatility",
			HistoricalData: historicalData,
			Options: PredictiveOptions{
				IncludeVolatilityForecast: true,
			},
			RequestedAt: time.Now(),
		}

		result, err := engine.GeneratePredictiveAnalytics(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Validate volatility forecast
		assert.NotNil(t, result.VolatilityForecast)

		// Check volatility predictions for BTC
		btcForecast, exists := result.VolatilityForecast.Forecasts["BTC"]
		assert.True(t, exists)
		assert.Equal(t, "BTC", btcForecast.Symbol)
		assert.Greater(t, btcForecast.CurrentVolatility, 0.0)
		assert.NotEmpty(t, btcForecast.PredictedVolatility)
		assert.Contains(t, []string{"low", "normal", "high", "extreme"}, btcForecast.VolatilityRegime)
		assert.GreaterOrEqual(t, btcForecast.Confidence, 0.0)
		assert.LessOrEqual(t, btcForecast.Confidence, 1.0)

		// Check market volatility assessment
		assert.Contains(t, []string{"low", "medium", "high", "extreme"}, result.VolatilityForecast.MarketVolatility)
		assert.Contains(t, []string{"increasing", "decreasing", "stable"}, result.VolatilityForecast.VolatilityTrend)
		assert.Greater(t, result.VolatilityForecast.VIXEquivalent, 0.0)
	})

	t.Run("CorrelationMatrix", func(t *testing.T) {
		ctx := context.Background()

		// Create correlated data for BTC and ETH
		historicalData := make(map[string][]ml.PriceData)
		symbols := []string{"BTC", "ETH"}

		btcData := make([]ml.PriceData, 150)
		ethData := make([]ml.PriceData, 150)

		for i := 0; i < 150; i++ {
			// Create correlated price movements
			btcPrice := 50000.0 + float64(i)*100
			ethPrice := 3000.0 + float64(i)*6 // Correlated with BTC

			timestamp := time.Now().Add(time.Duration(-150+i) * time.Hour)

			btcData[i] = ml.PriceData{
				Symbol:    "BTC",
				Timestamp: timestamp,
				Close:     decimal.NewFromFloat(btcPrice),
				Volume:    decimal.NewFromFloat(1000000),
			}

			ethData[i] = ml.PriceData{
				Symbol:    "ETH",
				Timestamp: timestamp,
				Close:     decimal.NewFromFloat(ethPrice),
				Volume:    decimal.NewFromFloat(500000),
			}
		}

		historicalData["BTC"] = btcData
		historicalData["ETH"] = ethData

		req := &PredictiveRequest{
			Symbols:        symbols,
			TimeHorizon:    6,
			AnalysisType:   "correlation",
			HistoricalData: historicalData,
			Options: PredictiveOptions{
				IncludeCorrelationMatrix: true,
			},
			RequestedAt: time.Now(),
		}

		result, err := engine.GeneratePredictiveAnalytics(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Validate correlation matrix
		assert.NotNil(t, result.CorrelationMatrix)

		// Check matrix structure
		matrix := result.CorrelationMatrix.Matrix
		assert.Contains(t, matrix, "BTC")
		assert.Contains(t, matrix, "ETH")

		// Check self-correlations
		assert.Equal(t, 1.0, matrix["BTC"]["BTC"])
		assert.Equal(t, 1.0, matrix["ETH"]["ETH"])

		// Check cross-correlations
		btcEthCorr := matrix["BTC"]["ETH"]
		ethBtcCorr := matrix["ETH"]["BTC"]
		assert.Equal(t, btcEthCorr, ethBtcCorr) // Should be symmetric
		assert.GreaterOrEqual(t, btcEthCorr, -1.0)
		assert.LessOrEqual(t, btcEthCorr, 1.0)

		// Check market correlation and diversification score
		assert.GreaterOrEqual(t, result.CorrelationMatrix.MarketCorrelation, 0.0)
		assert.LessOrEqual(t, result.CorrelationMatrix.MarketCorrelation, 1.0)
		assert.GreaterOrEqual(t, result.CorrelationMatrix.DiversificationScore, 0.0)
		assert.LessOrEqual(t, result.CorrelationMatrix.DiversificationScore, 1.0)
	})

	t.Run("PortfolioOptimization", func(t *testing.T) {
		ctx := context.Background()

		// Create portfolio data
		portfolioData := &PortfolioData{
			Holdings: []Holding{
				{
					Symbol: "BTC",
					Weight: 0.6,
					Value:  decimal.NewFromFloat(60000),
					Risk:   0.3,
				},
				{
					Symbol: "ETH",
					Weight: 0.4,
					Value:  decimal.NewFromFloat(40000),
					Risk:   0.35,
				},
			},
			TotalValue: decimal.NewFromFloat(100000),
			Cash:       decimal.NewFromFloat(5000),
			Leverage:   1.0,
			Constraints: map[string]interface{}{
				"max_weight": 0.7,
				"min_weight": 0.1,
			},
		}

		// Create historical data
		historicalData := make(map[string][]ml.PriceData)
		symbols := []string{"BTC", "ETH"}

		for _, symbol := range symbols {
			data := make([]ml.PriceData, 120)
			basePrice := 50000.0
			if symbol == "ETH" {
				basePrice = 3000.0
			}

			for i := 0; i < 120; i++ {
				price := basePrice * (1 + 0.001*float64(i)) // Slight upward trend
				data[i] = ml.PriceData{
					Symbol:    symbol,
					Timestamp: time.Now().Add(time.Duration(-120+i) * time.Hour),
					Close:     decimal.NewFromFloat(price),
					Volume:    decimal.NewFromFloat(1000000),
				}
			}
			historicalData[symbol] = data
		}

		req := &PredictiveRequest{
			Symbols:        symbols,
			TimeHorizon:    168, // 1 week
			AnalysisType:   "portfolio",
			HistoricalData: historicalData,
			PortfolioData:  portfolioData,
			Options: PredictiveOptions{
				IncludePortfolioOptimization: true,
				RiskTolerance:                "moderate",
				OptimizationObjective:        "sharpe",
			},
			RequestedAt: time.Now(),
		}

		result, err := engine.GeneratePredictiveAnalytics(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Validate portfolio optimization
		assert.NotNil(t, result.PortfolioOptimization)

		opt := result.PortfolioOptimization

		// Check optimal weights
		assert.NotEmpty(t, opt.OptimalWeights)
		totalWeight := 0.0
		for symbol, weight := range opt.OptimalWeights {
			assert.Contains(t, symbols, symbol)
			assert.GreaterOrEqual(t, weight, 0.0)
			assert.LessOrEqual(t, weight, 1.0)
			totalWeight += weight
		}
		assert.InDelta(t, 1.0, totalWeight, 0.01) // Should sum to ~1.0

		// Check portfolio metrics
		assert.NotZero(t, opt.ExpectedReturn)
		assert.Greater(t, opt.ExpectedRisk, 0.0)
		assert.NotZero(t, opt.SharpeRatio)
		assert.Greater(t, opt.VaR95, 0.0)
		assert.Greater(t, opt.CVaR95, 0.0)

		// Check efficient frontier
		assert.NotEmpty(t, opt.EfficientFrontier)
		for _, point := range opt.EfficientFrontier {
			assert.Greater(t, point.Risk, 0.0)
			assert.NotZero(t, point.Return)
			assert.NotZero(t, point.Sharpe)
			assert.NotEmpty(t, point.Weights)
		}

		// Check rebalancing strategy
		assert.NotNil(t, opt.Rebalancing)
		assert.NotEmpty(t, opt.Rebalancing.Frequency)
		assert.Greater(t, opt.Rebalancing.Threshold, 0.0)
		assert.NotEmpty(t, opt.Rebalancing.Method)
	})

	t.Run("RiskMetrics", func(t *testing.T) {
		ctx := context.Background()

		// Create sample data
		historicalData := make(map[string][]ml.PriceData)
		symbols := []string{"BTC", "ETH"}

		for _, symbol := range symbols {
			data := make([]ml.PriceData, 100)
			basePrice := 50000.0
			if symbol == "ETH" {
				basePrice = 3000.0
			}

			for i := 0; i < 100; i++ {
				// Add some volatility
				volatility := 0.02 * (1 + float64(i%5)/10.0)
				priceChange := volatility * (float64(i%3) - 1) // Some randomness
				price := basePrice * (1 + priceChange)

				data[i] = ml.PriceData{
					Symbol:    symbol,
					Timestamp: time.Now().Add(time.Duration(-100+i) * time.Hour),
					Close:     decimal.NewFromFloat(price),
					Volume:    decimal.NewFromFloat(1000000),
				}
			}
			historicalData[symbol] = data
		}

		req := &PredictiveRequest{
			Symbols:        symbols,
			TimeHorizon:    24,
			AnalysisType:   "risk",
			HistoricalData: historicalData,
			Options: PredictiveOptions{
				IncludeRiskMetrics: true,
			},
			RequestedAt: time.Now(),
		}

		result, err := engine.GeneratePredictiveAnalytics(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Validate risk metrics
		assert.NotNil(t, result.RiskMetrics)

		risk := result.RiskMetrics

		// Check portfolio risk
		assert.NotNil(t, risk.PortfolioRisk)
		assert.Greater(t, risk.PortfolioRisk.TotalRisk, 0.0)
		assert.Greater(t, risk.PortfolioRisk.SystematicRisk, 0.0)
		assert.Greater(t, risk.PortfolioRisk.IdiosyncraticRisk, 0.0)
		assert.NotZero(t, risk.PortfolioRisk.Beta)

		// Check individual asset risks
		assert.NotEmpty(t, risk.IndividualRisks)
		for _, symbol := range symbols {
			assetRisk, exists := risk.IndividualRisks[symbol]
			assert.True(t, exists)
			assert.Equal(t, symbol, assetRisk.Symbol)
			assert.Greater(t, assetRisk.Volatility, 0.0)
			assert.Greater(t, assetRisk.VaR95, 0.0)
			assert.Greater(t, assetRisk.CVaR95, 0.0)
		}

		// Check market risk
		assert.NotNil(t, risk.MarketRisk)
		assert.NotZero(t, risk.MarketRisk.MarketBeta)
		assert.GreaterOrEqual(t, risk.MarketRisk.MarketCorrelation, 0.0)
		assert.LessOrEqual(t, risk.MarketRisk.MarketCorrelation, 1.0)

		// Check liquidity risk
		assert.NotNil(t, risk.LiquidityRisk)
		assert.GreaterOrEqual(t, risk.LiquidityRisk.LiquidityScore, 0.0)
		assert.LessOrEqual(t, risk.LiquidityRisk.LiquidityScore, 1.0)

		// Check concentration risk
		assert.NotNil(t, risk.ConcentrationRisk)
		assert.GreaterOrEqual(t, risk.ConcentrationRisk.HerfindahlIndex, 0.0)
		assert.LessOrEqual(t, risk.ConcentrationRisk.HerfindahlIndex, 1.0)

		// Check tail risk
		assert.NotNil(t, risk.TailRisk)
		assert.Greater(t, risk.TailRisk.VaR99, 0.0)
		assert.Greater(t, risk.TailRisk.CVaR99, 0.0)
		assert.NotEmpty(t, risk.TailRisk.ExtremeEvents)
	})

	t.Run("ScenarioAnalysis", func(t *testing.T) {
		ctx := context.Background()

		// Create sample data
		historicalData := make(map[string][]ml.PriceData)
		symbols := []string{"BTC"}

		data := make([]ml.PriceData, 100)
		for i := 0; i < 100; i++ {
			price := 50000.0 + float64(i)*100
			data[i] = ml.PriceData{
				Symbol:    "BTC",
				Timestamp: time.Now().Add(time.Duration(-100+i) * time.Hour),
				Close:     decimal.NewFromFloat(price),
				Volume:    decimal.NewFromFloat(1000000),
			}
		}
		historicalData["BTC"] = data

		req := &PredictiveRequest{
			Symbols:        symbols,
			TimeHorizon:    48,
			AnalysisType:   "scenario",
			HistoricalData: historicalData,
			Options: PredictiveOptions{
				IncludeScenarioAnalysis: true,
			},
			RequestedAt: time.Now(),
		}

		result, err := engine.GeneratePredictiveAnalytics(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Validate scenario analysis
		assert.NotNil(t, result.ScenarioAnalysis)

		scenario := result.ScenarioAnalysis

		// Check scenarios
		assert.NotEmpty(t, scenario.Scenarios)
		for _, s := range scenario.Scenarios {
			assert.NotEmpty(t, s.Name)
			assert.NotEmpty(t, s.Description)
			assert.GreaterOrEqual(t, s.Probability, 0.0)
			assert.LessOrEqual(t, s.Probability, 1.0)
			assert.NotEmpty(t, s.Duration)
			assert.NotEmpty(t, s.Impact)
			assert.NotEmpty(t, s.Triggers)
		}

		// Check stress tests
		assert.NotEmpty(t, scenario.StressTests)
		for _, test := range scenario.StressTests {
			assert.NotEmpty(t, test.Name)
			assert.Contains(t, []string{"historical", "hypothetical", "monte_carlo"}, test.Type)
			assert.Contains(t, []string{"mild", "moderate", "severe", "extreme"}, test.Severity)
			assert.NotEmpty(t, test.MarketShock)
			assert.GreaterOrEqual(t, test.Probability, 0.0)
			assert.LessOrEqual(t, test.Probability, 1.0)
		}

		// Check Monte Carlo simulations
		assert.NotNil(t, scenario.MonteCarloSims)
		mc := scenario.MonteCarloSims
		assert.Greater(t, mc.Simulations, 0)
		assert.Equal(t, req.TimeHorizon, mc.TimeHorizon)
		assert.NotEmpty(t, mc.Returns)
		assert.NotEmpty(t, mc.Percentiles)
		assert.GreaterOrEqual(t, mc.ProbabilityOfLoss, 0.0)
		assert.LessOrEqual(t, mc.ProbabilityOfLoss, 1.0)

		// Check sensitivity analysis
		assert.NotNil(t, scenario.SensitivityAnalysis)
		sens := scenario.SensitivityAnalysis
		assert.NotEmpty(t, sens.Sensitivities)
		assert.NotEmpty(t, sens.KeyDrivers)
		assert.NotEmpty(t, sens.RiskFactors)
		assert.GreaterOrEqual(t, sens.Confidence, 0.0)
		assert.LessOrEqual(t, sens.Confidence, 1.0)
	})

	t.Run("MarketRegime", func(t *testing.T) {
		ctx := context.Background()

		// Create sample data
		historicalData := make(map[string][]ml.PriceData)
		symbols := []string{"BTC"}

		data := make([]ml.PriceData, 100)
		for i := 0; i < 100; i++ {
			price := 50000.0 + float64(i)*100
			data[i] = ml.PriceData{
				Symbol:    "BTC",
				Timestamp: time.Now().Add(time.Duration(-100+i) * time.Hour),
				Close:     decimal.NewFromFloat(price),
				Volume:    decimal.NewFromFloat(1000000),
			}
		}
		historicalData["BTC"] = data

		// Create market data
		marketData := &ml.MarketData{
			Timestamp:      time.Now(),
			FearGreedIndex: 65,
			Volatility:     0.4,
		}

		req := &PredictiveRequest{
			Symbols:        symbols,
			TimeHorizon:    24,
			AnalysisType:   "regime",
			HistoricalData: historicalData,
			MarketData:     marketData,
			Options:        PredictiveOptions{},
			RequestedAt:    time.Now(),
		}

		result, err := engine.GeneratePredictiveAnalytics(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Validate market regime
		assert.NotNil(t, result.MarketRegime)

		regime := result.MarketRegime

		// Check current regime
		assert.Contains(t, []string{"bull", "bear", "sideways", "volatile"}, regime.CurrentRegime)
		assert.GreaterOrEqual(t, regime.RegimeProbability, 0.0)
		assert.LessOrEqual(t, regime.RegimeProbability, 1.0)

		// Check regime history
		assert.NotEmpty(t, regime.RegimeHistory)
		for _, change := range regime.RegimeHistory {
			assert.NotEmpty(t, change.FromRegime)
			assert.NotEmpty(t, change.ToRegime)
			assert.NotEmpty(t, change.Duration)
			assert.NotEmpty(t, change.Trigger)
			assert.GreaterOrEqual(t, change.Confidence, 0.0)
			assert.LessOrEqual(t, change.Confidence, 1.0)
		}

		// Check transition matrix
		assert.NotEmpty(t, regime.TransitionMatrix)
		for fromRegime, transitions := range regime.TransitionMatrix {
			assert.Contains(t, []string{"bull", "bear", "sideways", "volatile"}, fromRegime)
			totalProb := 0.0
			for toRegime, prob := range transitions {
				assert.Contains(t, []string{"bull", "bear", "sideways", "volatile"}, toRegime)
				assert.GreaterOrEqual(t, prob, 0.0)
				assert.LessOrEqual(t, prob, 1.0)
				totalProb += prob
			}
			assert.InDelta(t, 1.0, totalProb, 0.01) // Should sum to ~1.0
		}

		// Check expected duration
		assert.NotEmpty(t, regime.ExpectedDuration)

		// Check indicators
		assert.NotEmpty(t, regime.Indicators)
		for _, indicator := range regime.Indicators {
			assert.NotEmpty(t, indicator.Name)
			assert.Contains(t, []string{"bullish", "bearish", "neutral"}, indicator.Signal)
			assert.GreaterOrEqual(t, indicator.Weight, 0.0)
			assert.LessOrEqual(t, indicator.Weight, 1.0)
			assert.GreaterOrEqual(t, indicator.Confidence, 0.0)
			assert.LessOrEqual(t, indicator.Confidence, 1.0)
		}
	})

	t.Run("CacheValidation", func(t *testing.T) {
		ctx := context.Background()

		// Create simple request
		historicalData := make(map[string][]ml.PriceData)
		symbols := []string{"BTC"}

		data := make([]ml.PriceData, 100)
		for i := 0; i < 100; i++ {
			price := 50000.0 + float64(i)*10
			data[i] = ml.PriceData{
				Symbol:    "BTC",
				Timestamp: time.Now().Add(time.Duration(-100+i) * time.Hour),
				Close:     decimal.NewFromFloat(price),
				Volume:    decimal.NewFromFloat(1000000),
			}
		}
		historicalData["BTC"] = data

		req := &PredictiveRequest{
			Symbols:        symbols,
			TimeHorizon:    12,
			AnalysisType:   "trend",
			HistoricalData: historicalData,
			Options: PredictiveOptions{
				IncludeTrendAnalysis: true,
			},
			RequestedAt: time.Now(),
		}

		// First request - should generate new result
		result1, err := engine.GeneratePredictiveAnalytics(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, result1)

		// Second request - should return cached result
		result2, err := engine.GeneratePredictiveAnalytics(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, result2)

		// Results should be identical (from cache)
		assert.Equal(t, result1.GeneratedAt, result2.GeneratedAt)
		assert.Equal(t, result1.Confidence, result2.Confidence)
	})
}
