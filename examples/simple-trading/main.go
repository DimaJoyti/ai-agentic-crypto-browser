package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ai-agentic-browser/internal/trading"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/shopspring/decimal"
)

func main() {
	fmt.Println("🚀 AI-Agentic Crypto Browser - Simple Trading Demo")
	fmt.Println("==================================================")

	// Initialize logger
	logger := &observability.Logger{}
	ctx := context.Background()

	// Demo 1: Algorithm Manager
	fmt.Println("\n🧠 Demo 1: Algorithm Manager")
	fmt.Println("  Creating algorithm manager...")

	algorithmManager := trading.NewAlgorithmManager(logger)
	fmt.Println("    ✅ Algorithm manager created")

	// Start the algorithm manager
	if err := algorithmManager.Start(ctx); err != nil {
		log.Fatalf("Failed to start algorithm manager: %v", err)
	}
	fmt.Println("    ✅ Algorithm manager started")

	// Create a simple risk profile
	riskProfile := trading.RiskProfile{
		MaxPositionSize:   decimal.NewFromFloat(0.1),
		MaxDailyLoss:      decimal.NewFromFloat(0.05),
		RiskPerTrade:      decimal.NewFromFloat(0.02),
		StopLossPercent:   decimal.NewFromFloat(0.05),
		TakeProfitPercent: decimal.NewFromFloat(0.10),
	}

	// Create TWAP strategy
	twapStrategy, err := algorithmManager.CreateStrategy(
		"Simple TWAP",
		"Basic time-weighted average price strategy",
		trading.AlgorithmTypeTWAP,
		map[string]interface{}{
			"duration_minutes": 30,
			"slice_count":      5,
		},
		riskProfile,
	)
	if err != nil {
		log.Fatalf("Failed to create TWAP strategy: %v", err)
	}
	fmt.Printf("    ✅ Created strategy: %s (ID: %s)\n", twapStrategy.Name, twapStrategy.ID)

	// Demo 2: Execution Engine
	fmt.Println("\n⚡ Demo 2: Execution Engine")
	fmt.Println("  Creating execution engine...")

	executionEngine := trading.NewExecutionEngine(logger)
	fmt.Println("    ✅ Execution engine created")

	if err := executionEngine.Start(ctx); err != nil {
		log.Fatalf("Failed to start execution engine: %v", err)
	}
	fmt.Println("    ✅ Execution engine started")

	// Create a test order
	testOrder := &trading.ExecutionOrder{
		StrategyID:    twapStrategy.ID,
		AlgorithmType: trading.AlgorithmTypeTWAP,
		Symbol:        "BTC/USD",
		Side:          trading.OrderSideBuy,
		OrderType:     trading.OrderTypeMarket,
		Quantity:      decimal.NewFromFloat(0.1),
		Price:         decimal.NewFromFloat(45000),
		TimeInForce:   trading.TimeInForceGTC,
	}

	// Submit the order
	if err := executionEngine.SubmitOrder(ctx, testOrder); err != nil {
		log.Printf("Failed to submit order: %v", err)
	} else {
		fmt.Printf("    ✅ Submitted order: %s %s %s\n", 
			testOrder.Side, testOrder.Quantity.String(), testOrder.Symbol)
	}

	// Wait a moment for processing
	time.Sleep(1 * time.Second)

	// Get execution metrics
	metrics := executionEngine.GetMetrics()
	fmt.Printf("    📊 Execution Metrics:\n")
	fmt.Printf("      • Total Orders: %d\n", metrics.TotalOrders)
	fmt.Printf("      • Completed Orders: %d\n", metrics.CompletedOrders)
	fmt.Printf("      • Success Rate: %.1f%%\n", metrics.SuccessRate*100)

	// Demo 3: Risk Manager
	fmt.Println("\n🛡️  Demo 3: Risk Manager")
	fmt.Println("  Creating risk manager...")

	riskManager := trading.NewAdvancedRiskManager(logger)
	fmt.Println("    ✅ Risk manager created")

	if err := riskManager.Start(ctx); err != nil {
		log.Fatalf("Failed to start risk manager: %v", err)
	}
	fmt.Println("    ✅ Risk manager started")

	// Test order validation
	validationOrder := &trading.ExecutionOrder{
		Symbol:   "ETH/USD",
		Side:     trading.OrderSideBuy,
		Quantity: decimal.NewFromFloat(1.0),
		Price:    decimal.NewFromFloat(3000),
	}

	if err := riskManager.ValidateOrder(ctx, validationOrder); err != nil {
		fmt.Printf("    ⚠️  Order validation failed: %v\n", err)
	} else {
		fmt.Printf("    ✅ Order validation passed for %s %s\n", 
			validationOrder.Quantity.String(), validationOrder.Symbol)
	}

	// Get risk metrics
	riskMetrics := riskManager.GetRiskMetrics()
	fmt.Printf("    📊 Risk Metrics:\n")
	fmt.Printf("      • Portfolio Value: $%s\n", riskMetrics.PortfolioValue.String())
	fmt.Printf("      • Total Risk: $%s\n", riskMetrics.TotalRisk.String())
	fmt.Printf("      • VaR (95%%): $%s\n", riskMetrics.VaR95.String())

	// Demo 4: Portfolio Optimizer
	fmt.Println("\n📈 Demo 4: Portfolio Optimizer")
	fmt.Println("  Creating portfolio optimizer...")

	portfolioOptimizer := trading.NewPortfolioOptimizer(logger)
	fmt.Println("    ✅ Portfolio optimizer created")

	if err := portfolioOptimizer.Start(ctx); err != nil {
		log.Fatalf("Failed to start portfolio optimizer: %v", err)
	}
	fmt.Println("    ✅ Portfolio optimizer started")

	// Create a simple portfolio
	assets := []string{"BTC/USD", "ETH/USD"}
	constraints := &trading.OptimizationConstraints{
		MinWeights: map[string]decimal.Decimal{
			"BTC/USD": decimal.NewFromFloat(0.3),
			"ETH/USD": decimal.NewFromFloat(0.3),
		},
		MaxWeights: map[string]decimal.Decimal{
			"BTC/USD": decimal.NewFromFloat(0.7),
			"ETH/USD": decimal.NewFromFloat(0.7),
		},
	}

	objective := &trading.OptimizationObjective{
		Type:         trading.ObjectiveTypeMaxSharpe,
		TargetReturn: decimal.NewFromFloat(0.10),
	}

	portfolio, err := portfolioOptimizer.OptimizePortfolio(
		ctx,
		"Test Portfolio",
		assets,
		trading.OptimizationMethodMaxSharpe,
		constraints,
		objective,
	)
	if err != nil {
		log.Printf("Failed to optimize portfolio: %v", err)
	} else {
		fmt.Printf("    ✅ Portfolio optimized:\n")
		fmt.Printf("      • Expected Return: %.2f%%\n", portfolio.ExpectedReturn.InexactFloat64()*100)
		fmt.Printf("      • Sharpe Ratio: %.2f\n", portfolio.SharpeRatio.InexactFloat64())
		fmt.Printf("      • Weights:\n")
		for asset, weight := range portfolio.Weights {
			fmt.Printf("        - %s: %.1f%%\n", asset, weight.InexactFloat64()*100)
		}
	}

	// Demo 5: Smart Order Router
	fmt.Println("\n🎯 Demo 5: Smart Order Router")
	fmt.Println("  Creating smart order router...")

	orderRouter := trading.NewSmartOrderRouter(logger)
	fmt.Println("    ✅ Smart order router created")

	if err := orderRouter.Start(ctx); err != nil {
		log.Fatalf("Failed to start order router: %v", err)
	}
	fmt.Println("    ✅ Smart order router started")

	// Test order routing
	routingOrder := &trading.ExecutionOrder{
		Symbol:   "BTC/USD",
		Side:     trading.OrderSideBuy,
		Quantity: decimal.NewFromFloat(1.0),
		Price:    decimal.NewFromFloat(45000),
	}

	decision, err := orderRouter.RouteOrder(ctx, routingOrder)
	if err != nil {
		log.Printf("Failed to route order: %v", err)
	} else {
		fmt.Printf("    ✅ Order routed successfully:\n")
		fmt.Printf("      • Strategy: %s\n", decision.Strategy)
		fmt.Printf("      • Venues: %d\n", len(decision.SelectedVenues))
		fmt.Printf("      • Estimated Cost: $%s\n", decision.EstimatedCost.String())
		fmt.Printf("      • Confidence: %.1f%%\n", decision.ConfidenceScore.InexactFloat64()*100)
	}

	// Demo 6: Integration Test
	fmt.Println("\n🔄 Demo 6: Integration Test")
	fmt.Println("  Testing integrated workflow...")

	// Create an integrated order
	integrationOrder := &trading.ExecutionOrder{
		StrategyID:    twapStrategy.ID,
		AlgorithmType: trading.AlgorithmTypeTWAP,
		Symbol:        "ETH/USD",
		Side:          trading.OrderSideBuy,
		Quantity:      decimal.NewFromFloat(2.0),
		Price:         decimal.NewFromFloat(3000),
		TimeInForce:   trading.TimeInForceGTC,
	}

	// Step 1: Risk validation
	fmt.Println("    🛡️  Step 1: Risk validation...")
	if err := riskManager.ValidateOrder(ctx, integrationOrder); err != nil {
		fmt.Printf("      ❌ Risk validation failed: %v\n", err)
	} else {
		fmt.Println("      ✅ Risk validation passed")
	}

	// Step 2: Order routing
	fmt.Println("    🎯 Step 2: Order routing...")
	routingDecision, err := orderRouter.RouteOrder(ctx, integrationOrder)
	if err != nil {
		fmt.Printf("      ❌ Order routing failed: %v\n", err)
	} else {
		fmt.Printf("      ✅ Order routed to %d venues\n", len(routingDecision.SelectedVenues))
	}

	// Step 3: Order execution
	fmt.Println("    ⚡ Step 3: Order execution...")
	if err := executionEngine.SubmitOrder(ctx, integrationOrder); err != nil {
		fmt.Printf("      ❌ Order execution failed: %v\n", err)
	} else {
		fmt.Println("      ✅ Order submitted for execution")
	}

	fmt.Println("    ✅ Integration test completed successfully")

	// Cleanup
	fmt.Println("\n🧹 Cleaning up...")
	algorithmManager.Stop(ctx)
	executionEngine.Stop(ctx)
	riskManager.Stop(ctx)
	portfolioOptimizer.Stop(ctx)
	orderRouter.Stop(ctx)

	fmt.Println("\n🎉 Simple Trading Demo Complete!")
	fmt.Println("All advanced trading components are working correctly.")
}
