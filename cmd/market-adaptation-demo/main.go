package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/ai-agentic-browser/internal/ai"
	"github.com/ai-agentic-browser/pkg/observability"
)

func main() {
	fmt.Println("🚀 AI Agentic Crypto Browser - Market Pattern Adaptation System Demo")
	fmt.Println(strings.Repeat("=", 80))

	// Initialize logger
	logger := &observability.Logger{}

	// Create market adaptation engine
	engine := ai.NewMarketAdaptationEngine(logger)

	ctx := context.Background()

	// Demo 1: Pattern Detection
	fmt.Println("\n📊 Demo 1: Market Pattern Detection")
	fmt.Println(strings.Repeat("-", 40))

	// Simulate market data for different scenarios
	scenarios := []struct {
		name string
		data map[string]interface{}
	}{
		{
			name: "Bullish Trend",
			data: map[string]interface{}{
				"prices":     []float64{50000, 50500, 51000, 51500, 52000, 52500, 53000, 53500, 54000, 54500},
				"volumes":    []float64{100, 120, 110, 130, 140, 125, 135, 145, 150, 160},
				"timestamps": generateTimestamps(10),
			},
		},
		{
			name: "Bearish Trend",
			data: map[string]interface{}{
				"prices":     []float64{54000, 53500, 53000, 52500, 52000, 51500, 51000, 50500, 50000, 49500},
				"volumes":    []float64{150, 160, 170, 180, 190, 200, 210, 220, 230, 240},
				"timestamps": generateTimestamps(10),
			},
		},
		{
			name: "Sideways Movement",
			data: map[string]interface{}{
				"prices":     []float64{52000, 52100, 51900, 52050, 51950, 52000, 52100, 51900, 52000, 52050},
				"volumes":    []float64{100, 105, 95, 110, 90, 100, 105, 95, 100, 110},
				"timestamps": generateTimestamps(10),
			},
		},
	}

	for _, scenario := range scenarios {
		fmt.Printf("\n🔍 Analyzing %s scenario...\n", scenario.name)

		patterns, err := engine.DetectPatterns(ctx, scenario.data)
		if err != nil {
			log.Printf("Error detecting patterns: %v", err)
			continue
		}

		if len(patterns) > 0 {
			pattern := patterns[0]
			fmt.Printf("✅ Pattern detected: %s\n", pattern.Name)
			fmt.Printf("   Type: %s\n", pattern.Type)
			fmt.Printf("   Confidence: %.2f\n", pattern.Confidence)
			fmt.Printf("   Strength: %.2f\n", pattern.Strength)
			fmt.Printf("   Expected Direction: %s\n", pattern.ExpectedOutcome.Direction)
			fmt.Printf("   Expected Magnitude: %.2f%%\n", pattern.ExpectedOutcome.Magnitude*100)
			fmt.Printf("   Success Rate: %.2f%%\n", pattern.ExpectedOutcome.SuccessRate*100)
		} else {
			fmt.Printf("❌ No significant patterns detected\n")
		}
	}

	// Demo 2: Adaptive Strategy Management
	fmt.Println("\n\n🎯 Demo 2: Adaptive Strategy Management")
	fmt.Println(strings.Repeat("-", 40))

	// Create sample adaptive strategies
	strategies := []*ai.AdaptiveStrategy{
		{
			Name:        "Trend Following Strategy",
			Description: "Follows market trends with adaptive position sizing",
			Type:        "trend_following",
			BaseParameters: map[string]float64{
				"position_size":   0.05,
				"stop_loss":       0.02,
				"take_profit":     0.04,
				"entry_threshold": 0.7,
			},
			CurrentParameters: map[string]float64{
				"position_size":   0.05,
				"stop_loss":       0.02,
				"take_profit":     0.04,
				"entry_threshold": 0.7,
			},
			PerformanceTargets: &ai.PerformanceTargets{
				TargetReturn:       0.15,
				MaxDrawdown:        0.1,
				MinSharpeRatio:     1.0,
				MinWinRate:         0.6,
				MaxVolatility:      0.3,
				TargetProfitFactor: 1.5,
				EvaluationPeriod:   30 * 24 * time.Hour,
			},
			RiskLimits: &ai.MarketRiskLimits{
				MaxPositionSize:    0.1,
				MaxLeverage:        2.0,
				StopLossPercentage: 0.05,
				TakeProfitRatio:    2.0,
				MaxDailyLoss:       0.02,
				VaRLimit:           0.01,
				ConcentrationLimit: 0.2,
			},
		},
		{
			Name:        "Mean Reversion Strategy",
			Description: "Exploits price reversions to mean with dynamic thresholds",
			Type:        "mean_reversion",
			BaseParameters: map[string]float64{
				"position_size":       0.03,
				"reversion_threshold": 2.0,
				"hold_time":           24.0,
				"volatility_filter":   0.3,
			},
			CurrentParameters: map[string]float64{
				"position_size":       0.03,
				"reversion_threshold": 2.0,
				"hold_time":           24.0,
				"volatility_filter":   0.3,
			},
			PerformanceTargets: &ai.PerformanceTargets{
				TargetReturn:       0.12,
				MaxDrawdown:        0.08,
				MinSharpeRatio:     1.2,
				MinWinRate:         0.65,
				MaxVolatility:      0.25,
				TargetProfitFactor: 1.8,
				EvaluationPeriod:   30 * 24 * time.Hour,
			},
			RiskLimits: &ai.MarketRiskLimits{
				MaxPositionSize:    0.08,
				MaxLeverage:        1.5,
				StopLossPercentage: 0.03,
				TakeProfitRatio:    1.5,
				MaxDailyLoss:       0.015,
				VaRLimit:           0.008,
				ConcentrationLimit: 0.15,
			},
		},
	}

	// Add strategies to the engine
	for _, strategy := range strategies {
		err := engine.AddAdaptiveStrategy(ctx, strategy)
		if err != nil {
			log.Printf("Error adding strategy: %v", err)
			continue
		}
		fmt.Printf("✅ Added strategy: %s (ID: %s)\n", strategy.Name, strategy.ID)
	}

	// Demo 3: Strategy Adaptation
	fmt.Println("\n\n🔄 Demo 3: Strategy Adaptation Process")
	fmt.Println(strings.Repeat("-", 40))

	// Simulate poor performance scenarios to trigger adaptations
	adaptationScenarios := []struct {
		name     string
		patterns []*ai.DetectedPattern
	}{
		{
			name: "Strong Bullish Signal",
			patterns: []*ai.DetectedPattern{
				{
					Type:       "trend",
					Confidence: 0.85,
					Strength:   0.9,
					ExpectedOutcome: &ai.ExpectedOutcome{
						Direction:   "up",
						Magnitude:   0.12,
						Probability: 0.8,
						TimeHorizon: 24 * time.Hour,
					},
				},
			},
		},
		{
			name: "High Volatility Environment",
			patterns: []*ai.DetectedPattern{
				{
					Type:       "volatility_spike",
					Confidence: 0.75,
					Strength:   0.8,
					ExpectedOutcome: &ai.ExpectedOutcome{
						Direction:   "sideways",
						Magnitude:   0.05,
						Probability: 0.6,
						TimeHorizon: 12 * time.Hour,
					},
				},
			},
		},
	}

	for _, scenario := range adaptationScenarios {
		fmt.Printf("\n🔄 Processing %s...\n", scenario.name)

		err := engine.AdaptStrategies(ctx, scenario.patterns)
		if err != nil {
			log.Printf("Error adapting strategies: %v", err)
			continue
		}

		// Show adaptation results
		adaptedStrategies, err := engine.GetAdaptiveStrategies(ctx)
		if err != nil {
			log.Printf("Error getting strategies: %v", err)
			continue
		}

		for _, strategy := range adaptedStrategies {
			if strategy.AdaptationCount > 0 {
				fmt.Printf("   📈 %s adapted %d times\n", strategy.Name, strategy.AdaptationCount)
				if len(strategy.AdaptationHistory) > 0 {
					lastAdaptation := strategy.AdaptationHistory[len(strategy.AdaptationHistory)-1]
					fmt.Printf("      Last adaptation: %s\n", lastAdaptation.TriggerReason)
					fmt.Printf("      Confidence: %.2f\n", lastAdaptation.Confidence)
				}
			}
		}
	}

	// Demo 4: Performance Metrics and History
	fmt.Println("\n\n📈 Demo 4: Performance Metrics and Adaptation History")
	fmt.Println(strings.Repeat("-", 40))

	// Simulate performance metrics for strategies
	allStrategies, err := engine.GetAdaptiveStrategies(ctx)
	if err != nil {
		log.Printf("Error getting strategies: %v", err)
		return
	}

	for _, strategy := range allStrategies {
		// Simulate performance metrics
		metrics := &ai.MarketPerformanceMetrics{
			StrategyID:       strategy.ID,
			TotalReturn:      0.08 + rand.Float64()*0.15,   // 8-23% return
			AnnualizedReturn: 0.12 + rand.Float64()*0.20,   // 12-32% annualized
			Volatility:       0.15 + rand.Float64()*0.15,   // 15-30% volatility
			SharpeRatio:      0.5 + rand.Float64()*1.0,     // 0.5-1.5 Sharpe
			MaxDrawdown:      0.03 + rand.Float64()*0.12,   // 3-15% drawdown
			WinRate:          0.55 + rand.Float64()*0.25,   // 55-80% win rate
			ProfitFactor:     1.1 + rand.Float64()*0.8,     // 1.1-1.9 profit factor
			TotalTrades:      int(50 + rand.Float64()*100), // 50-150 trades
			LastUpdated:      time.Now(),
		}

		// Add metrics to engine (simulating real performance tracking)
		engine.SetPerformanceMetrics(strategy.ID, metrics)

		fmt.Printf("\n📊 %s Performance:\n", strategy.Name)
		fmt.Printf("   Total Return: %.2f%%\n", metrics.TotalReturn*100)
		fmt.Printf("   Sharpe Ratio: %.2f\n", metrics.SharpeRatio)
		fmt.Printf("   Max Drawdown: %.2f%%\n", metrics.MaxDrawdown*100)
		fmt.Printf("   Win Rate: %.2f%%\n", metrics.WinRate*100)
		fmt.Printf("   Total Trades: %d\n", metrics.TotalTrades)
	}

	// Show adaptation history
	fmt.Println("\n📜 Adaptation History:")
	history, err := engine.GetAdaptationHistory(ctx, 10)
	if err != nil {
		log.Printf("Error getting adaptation history: %v", err)
		return
	}

	if len(history) > 0 {
		for i, record := range history {
			fmt.Printf("   %d. %s - %s (Confidence: %.2f)\n",
				i+1, record.Type, record.Description, record.Confidence)
		}
	} else {
		fmt.Printf("   No adaptation history available\n")
	}

	// Demo 5: Real-time Pattern Monitoring
	fmt.Println("\n\n🔍 Demo 5: Real-time Pattern Monitoring")
	fmt.Println(strings.Repeat("-", 40))

	// Simulate real-time market data updates
	fmt.Println("Simulating real-time market data updates...")

	for i := 0; i < 3; i++ {
		fmt.Printf("\n⏰ Update %d:\n", i+1)

		// Generate random market data
		marketData := generateRandomMarketData()

		patterns, err := engine.DetectPatterns(ctx, marketData)
		if err != nil {
			log.Printf("Error detecting patterns: %v", err)
			continue
		}

		if len(patterns) > 0 {
			fmt.Printf("   🎯 New pattern detected: %s (Confidence: %.2f)\n",
				patterns[0].Name, patterns[0].Confidence)

			// Check if adaptation is needed
			err = engine.AdaptStrategies(ctx, patterns)
			if err != nil {
				log.Printf("Error adapting strategies: %v", err)
				continue
			}

			fmt.Printf("   ✅ Strategies evaluated for adaptation\n")
		} else {
			fmt.Printf("   📊 Market data processed, no significant patterns\n")
		}

		time.Sleep(1 * time.Second) // Simulate real-time delay
	}

	// Demo 6: Export Configuration and Results
	fmt.Println("\n\n💾 Demo 6: System Configuration and Results")
	fmt.Println(strings.Repeat("-", 40))

	// Show current system configuration
	fmt.Println("📋 Current System Configuration:")
	config := map[string]interface{}{
		"pattern_detection_window":      "7 days",
		"adaptation_threshold":          0.7,
		"min_pattern_occurrences":       3,
		"strategy_update_frequency":     "1 hour",
		"performance_evaluation_window": "24 hours",
		"real_time_adaptation":          true,
		"confidence_threshold":          0.6,
	}

	configJSON, _ := json.MarshalIndent(config, "   ", "  ")
	fmt.Printf("   %s\n", string(configJSON))

	// Show summary statistics
	fmt.Println("\n📊 Session Summary:")
	allPatterns, _ := engine.GetDetectedPatterns(ctx, map[string]interface{}{})
	allStrategies, _ = engine.GetAdaptiveStrategies(ctx)
	allHistory, _ := engine.GetAdaptationHistory(ctx, 0)

	fmt.Printf("   Total Patterns Detected: %d\n", len(allPatterns))
	fmt.Printf("   Active Strategies: %d\n", len(allStrategies))
	fmt.Printf("   Adaptation Events: %d\n", len(allHistory))

	activeStrategies := 0
	totalAdaptations := 0
	for _, strategy := range allStrategies {
		if strategy.IsActive {
			activeStrategies++
		}
		totalAdaptations += strategy.AdaptationCount
	}

	fmt.Printf("   Active Strategies: %d\n", activeStrategies)
	fmt.Printf("   Total Adaptations: %d\n", totalAdaptations)

	fmt.Println("\n🎉 Market Pattern Adaptation System Demo Complete!")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("The system successfully demonstrated:")
	fmt.Println("✅ Real-time pattern detection across multiple market scenarios")
	fmt.Println("✅ Adaptive strategy management with dynamic parameter adjustment")
	fmt.Println("✅ Performance monitoring and risk management")
	fmt.Println("✅ Historical analysis and learning from past adaptations")
	fmt.Println("✅ Configurable thresholds and adaptation rules")
	fmt.Println("\nThe AI Agentic Crypto Browser is ready for intelligent market analysis! 🚀")
}

// Helper functions

func generateTimestamps(count int) []int64 {
	timestamps := make([]int64, count)
	baseTime := time.Now().Unix() - int64(count*3600) // Start count hours ago

	for i := 0; i < count; i++ {
		timestamps[i] = baseTime + int64(i*3600) // 1 hour intervals
	}

	return timestamps
}

func generateRandomMarketData() map[string]interface{} {
	count := 10
	basePrice := 50000.0 + rand.Float64()*10000 // 50k-60k base price

	prices := make([]float64, count)
	volumes := make([]float64, count)

	for i := 0; i < count; i++ {
		// Generate price with some trend and noise
		trend := (rand.Float64() - 0.5) * 0.02 // ±1% trend per step
		noise := (rand.Float64() - 0.5) * 0.01 // ±0.5% noise

		if i == 0 {
			prices[i] = basePrice
		} else {
			prices[i] = prices[i-1] * (1 + trend + noise)
		}

		// Generate volume with some correlation to price movement
		baseVolume := 100.0 + rand.Float64()*50
		if i > 0 {
			priceChange := math.Abs(prices[i]-prices[i-1]) / prices[i-1]
			volumeMultiplier := 1.0 + priceChange*10 // Higher volume on bigger moves
			volumes[i] = baseVolume * volumeMultiplier
		} else {
			volumes[i] = baseVolume
		}
	}

	return map[string]interface{}{
		"prices":     prices,
		"volumes":    volumes,
		"timestamps": generateTimestamps(count),
	}
}
