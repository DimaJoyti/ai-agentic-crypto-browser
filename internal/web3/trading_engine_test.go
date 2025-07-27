package web3

import (
	"context"
	"testing"

	"github.com/ai-agentic-browser/internal/config"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestTradingEngine(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{})
	clients := make(map[int]*ethclient.Client)
	riskAssessment := NewRiskAssessmentService(clients, logger)

	engine := NewTradingEngine(clients, logger, riskAssessment)

	t.Run("EngineInitialization", func(t *testing.T) {
		assert.NotNil(t, engine)
		assert.NotNil(t, engine.strategies)
		assert.NotNil(t, engine.activePositions)
		assert.NotNil(t, engine.portfolios)
		assert.Equal(t, 3, len(engine.strategies)) // momentum, mean_reversion, arbitrage
	})

	t.Run("CreatePortfolio", func(t *testing.T) {
		userID := uuid.New()
		initialBalance := decimal.NewFromInt(10000) // $10,000
		riskProfile := RiskProfile{
			Level:                "moderate",
			MaxPositionSize:      decimal.NewFromFloat(0.1),  // 10%
			MaxDailyLoss:         decimal.NewFromFloat(0.05), // 5%
			StopLossPercentage:   decimal.NewFromFloat(0.1),  // 10%
			TakeProfitPercentage: decimal.NewFromFloat(0.2),  // 20%
			AllowedStrategies:    []string{"momentum", "mean_reversion"},
		}

		portfolio, err := engine.CreatePortfolio(context.Background(), userID, "Test Portfolio", initialBalance, riskProfile)
		assert.NoError(t, err)
		assert.NotNil(t, portfolio)
		assert.Equal(t, userID, portfolio.UserID)
		assert.Equal(t, "Test Portfolio", portfolio.Name)
		assert.True(t, portfolio.TotalValue.Equal(initialBalance))
		assert.True(t, portfolio.AvailableBalance.Equal(initialBalance))
		assert.Equal(t, "moderate", portfolio.RiskProfile.Level)
	})

	t.Run("PortfolioManagement", func(t *testing.T) {
		userID := uuid.New()
		portfolio, err := engine.CreatePortfolio(
			context.Background(),
			userID,
			"Management Test",
			decimal.NewFromInt(5000),
			RiskProfile{Level: "conservative"},
		)
		assert.NoError(t, err)

		// Test portfolio retrieval
		retrieved, err := engine.GetPortfolio(portfolio.ID)
		assert.NoError(t, err)
		assert.Equal(t, portfolio.ID, retrieved.ID)

		// Test portfolio value update
		err = engine.UpdatePortfolioValue(context.Background(), portfolio.ID)
		assert.NoError(t, err)
	})
}

func TestTradingStrategies(t *testing.T) {
	t.Run("MomentumStrategy", func(t *testing.T) {
		strategy := NewMomentumStrategy()
		assert.Equal(t, "momentum_strategy", strategy.GetName())
		assert.Equal(t, RiskLevelMedium, strategy.GetRiskLevel())
		assert.True(t, strategy.IsEnabled())
		assert.NotEmpty(t, strategy.GetParameters())

		// Test strategy analysis with mock market data
		marketData := &MarketData{
			TokenAddress:   "0x1234567890123456789012345678901234567890",
			TokenSymbol:    "TEST",
			Price:          decimal.NewFromFloat(100),
			PriceChange24h: decimal.NewFromFloat(6), // 6% increase
			Volume24h:      decimal.NewFromInt(1000000),
			TechnicalData: &TechnicalIndicators{
				RSI:            decimal.NewFromFloat(25), // Oversold
				BollingerUpper: decimal.NewFromFloat(105),
				BollingerLower: decimal.NewFromFloat(95),
				Volume:         decimal.NewFromInt(500000), // 2x volume ratio
			},
		}

		signal, err := strategy.Analyze(context.Background(), marketData)
		assert.NoError(t, err)
		assert.NotNil(t, signal)
		assert.Equal(t, ActionBuy, signal.Action) // Should generate buy signal
		assert.Greater(t, signal.Confidence, 0.7)
	})

	t.Run("MeanReversionStrategy", func(t *testing.T) {
		strategy := NewMeanReversionStrategy()
		assert.Equal(t, "mean_reversion_strategy", strategy.GetName())
		assert.Equal(t, RiskLevelLow, strategy.GetRiskLevel())
		assert.True(t, strategy.IsEnabled())

		// Test with overbought conditions
		marketData := &MarketData{
			TokenAddress: "0x1234567890123456789012345678901234567890",
			TokenSymbol:  "TEST",
			Price:        decimal.NewFromFloat(108), // Above upper Bollinger band
			TechnicalData: &TechnicalIndicators{
				RSI:            decimal.NewFromFloat(85), // Overbought
				BollingerUpper: decimal.NewFromFloat(105),
				BollingerLower: decimal.NewFromFloat(95),
			},
		}

		signal, err := strategy.Analyze(context.Background(), marketData)
		assert.NoError(t, err)
		assert.NotNil(t, signal)
		assert.Equal(t, ActionSell, signal.Action) // Should generate sell signal
	})

	t.Run("ArbitrageStrategy", func(t *testing.T) {
		strategy := NewArbitrageStrategy()
		assert.Equal(t, "arbitrage_strategy", strategy.GetName())
		assert.Equal(t, RiskLevelLow, strategy.GetRiskLevel())
		assert.True(t, strategy.IsEnabled())

		// Arbitrage strategy should return hold for now (not implemented)
		marketData := &MarketData{
			TokenAddress: "0x1234567890123456789012345678901234567890",
			TokenSymbol:  "TEST",
			Price:        decimal.NewFromFloat(100),
		}

		signal, err := strategy.Analyze(context.Background(), marketData)
		assert.NoError(t, err)
		assert.NotNil(t, signal)
		assert.Equal(t, ActionHold, signal.Action)
	})

	t.Run("PositionSizeCalculation", func(t *testing.T) {
		strategy := NewMomentumStrategy()
		portfolio := &Portfolio{
			TotalValue:       decimal.NewFromInt(10000),
			AvailableBalance: decimal.NewFromInt(8000),
			RiskProfile: RiskProfile{
				MaxPositionSize: decimal.NewFromFloat(0.1), // 10%
			},
		}

		signal := &TradingSignal{
			Confidence: 0.8,
		}

		positionSize, err := strategy.CalculatePositionSize(context.Background(), signal, portfolio)
		assert.NoError(t, err)
		assert.True(t, positionSize.GreaterThan(decimal.Zero))
		assert.True(t, positionSize.LessThanOrEqual(portfolio.AvailableBalance))
	})
}

func TestDeFiProtocolManager(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{})
	manager := NewDeFiProtocolManager(logger)

	t.Run("ManagerInitialization", func(t *testing.T) {
		assert.NotNil(t, manager)
		protocols := manager.GetProtocols()
		assert.Greater(t, len(protocols), 0)

		// Check specific protocols
		uniswap, err := manager.GetProtocol("uniswap_v3")
		assert.NoError(t, err)
		assert.Equal(t, "Uniswap V3", uniswap.Name)
		assert.Equal(t, ProtocolTypeDEX, uniswap.Type)

		compound, err := manager.GetProtocol("compound")
		assert.NoError(t, err)
		assert.Equal(t, "Compound", compound.Name)
		assert.Equal(t, ProtocolTypeLending, compound.Type)

		aave, err := manager.GetProtocol("aave")
		assert.NoError(t, err)
		assert.Equal(t, "Aave", aave.Name)
		assert.Equal(t, ProtocolTypeLending, aave.Type)
	})

	t.Run("YieldOpportunities", func(t *testing.T) {
		minAPY := decimal.NewFromFloat(0.05) // 5% minimum APY
		maxRisk := RiskLevelMedium

		opportunities, err := manager.GetBestYieldOpportunities(context.Background(), minAPY, maxRisk)
		assert.NoError(t, err)
		assert.Greater(t, len(opportunities), 0)

		// Check that all opportunities meet criteria
		for _, opp := range opportunities {
			assert.True(t, opp.APY.GreaterThanOrEqual(minAPY))
			assert.True(t, manager.isRiskHigher(maxRisk, opp.RiskLevel) || opp.RiskLevel == maxRisk)
		}

		// Check sorting (should be sorted by APY descending)
		for i := 1; i < len(opportunities); i++ {
			assert.True(t, opportunities[i-1].APY.GreaterThanOrEqual(opportunities[i].APY))
		}
	})

	t.Run("ProtocolTypes", func(t *testing.T) {
		types := []ProtocolType{
			ProtocolTypeDEX,
			ProtocolTypeLending,
			ProtocolTypeYieldFarm,
			ProtocolTypeStaking,
			ProtocolTypeLiquidStaking,
			ProtocolTypeInsurance,
			ProtocolTypeSynthetics,
		}

		for _, protocolType := range types {
			assert.NotEmpty(t, string(protocolType))
		}
	})

	t.Run("PositionTypes", func(t *testing.T) {
		types := []PositionType{
			PositionTypeLiquidity,
			PositionTypeLending,
			PositionTypeBorrowing,
			PositionTypeStaking,
			PositionTypeFarming,
		}

		for _, positionType := range types {
			assert.NotEmpty(t, string(positionType))
		}
	})
}

func TestPortfolioRebalancer(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{})
	clients := make(map[int]*ethclient.Client)
	riskAssessment := NewRiskAssessmentService(clients, logger)
	tradingEngine := NewTradingEngine(clients, logger, riskAssessment)
	defiManager := NewDeFiProtocolManager(logger)

	rebalancer := NewPortfolioRebalancer(logger, tradingEngine, defiManager)

	t.Run("RebalancerInitialization", func(t *testing.T) {
		assert.NotNil(t, rebalancer)
		assert.NotNil(t, rebalancer.rebalanceRules)
		assert.NotNil(t, rebalancer.config)
	})

	t.Run("CreateRebalanceStrategy", func(t *testing.T) {
		portfolioID := uuid.New()
		targetAllocations := map[string]decimal.Decimal{
			"ETH":  decimal.NewFromFloat(0.5), // 50%
			"USDC": decimal.NewFromFloat(0.3), // 30%
			"BTC":  decimal.NewFromFloat(0.2), // 20%
		}

		strategy, err := rebalancer.CreateRebalanceStrategy(
			context.Background(),
			portfolioID,
			"Balanced Strategy",
			RebalanceTypeFixed,
			targetAllocations,
		)

		assert.NoError(t, err)
		assert.NotNil(t, strategy)
		assert.Equal(t, portfolioID, strategy.PortfolioID)
		assert.Equal(t, "Balanced Strategy", strategy.Name)
		assert.Equal(t, RebalanceTypeFixed, strategy.Type)
		assert.Equal(t, targetAllocations, strategy.TargetAllocations)
		assert.True(t, strategy.IsActive)
	})

	t.Run("InvalidAllocations", func(t *testing.T) {
		portfolioID := uuid.New()
		invalidAllocations := map[string]decimal.Decimal{
			"ETH":  decimal.NewFromFloat(0.6), // 60%
			"USDC": decimal.NewFromFloat(0.3), // 30%
			// Total = 90%, should fail
		}

		_, err := rebalancer.CreateRebalanceStrategy(
			context.Background(),
			portfolioID,
			"Invalid Strategy",
			RebalanceTypeFixed,
			invalidAllocations,
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must sum to 100%")
	})

	t.Run("RebalanceTypes", func(t *testing.T) {
		types := []RebalanceType{
			RebalanceTypeFixed,
			RebalanceTypeDynamic,
			RebalanceTypeRiskParity,
			RebalanceTypeMomentum,
			RebalanceTypeMeanRevert,
		}

		for _, rebalanceType := range types {
			assert.NotEmpty(t, string(rebalanceType))
		}
	})

	t.Run("TriggerTypes", func(t *testing.T) {
		types := []TriggerType{
			TriggerTypeDrift,
			TriggerTypeVolatility,
			TriggerTypeCorrelation,
			TriggerTypeTime,
			TriggerTypeDrawdown,
			TriggerTypeProfit,
		}

		for _, triggerType := range types {
			assert.NotEmpty(t, string(triggerType))
		}
	})
}

func TestTradingActions(t *testing.T) {
	t.Run("TradingActions", func(t *testing.T) {
		actions := []TradingAction{
			ActionBuy,
			ActionSell,
			ActionHold,
			ActionSwap,
			ActionStake,
			ActionUnstake,
		}

		for _, action := range actions {
			assert.NotEmpty(t, string(action))
		}
	})

	t.Run("SignalUrgency", func(t *testing.T) {
		urgencies := []SignalUrgency{
			UrgencyLow,
			UrgencyMedium,
			UrgencyHigh,
			UrgencyCritical,
		}

		for _, urgency := range urgencies {
			assert.NotEmpty(t, string(urgency))
		}
	})

	t.Run("PositionStatus", func(t *testing.T) {
		statuses := []PositionStatus{
			PositionStatusOpen,
			PositionStatusClosed,
			PositionStatusPending,
		}

		for _, status := range statuses {
			assert.NotEmpty(t, string(status))
		}
	})
}

func TestRiskProfiles(t *testing.T) {
	t.Run("ConservativeProfile", func(t *testing.T) {
		profile := RiskProfile{
			Level:                "conservative",
			MaxPositionSize:      decimal.NewFromFloat(0.05), // 5%
			MaxDailyLoss:         decimal.NewFromFloat(0.02), // 2%
			StopLossPercentage:   decimal.NewFromFloat(0.05), // 5%
			TakeProfitPercentage: decimal.NewFromFloat(0.1),  // 10%
			AllowedStrategies:    []string{"mean_reversion"},
		}

		assert.Equal(t, "conservative", profile.Level)
		assert.True(t, profile.MaxPositionSize.LessThan(decimal.NewFromFloat(0.1)))
		assert.Contains(t, profile.AllowedStrategies, "mean_reversion")
	})

	t.Run("AggressiveProfile", func(t *testing.T) {
		profile := RiskProfile{
			Level:                "aggressive",
			MaxPositionSize:      decimal.NewFromFloat(0.2),  // 20%
			MaxDailyLoss:         decimal.NewFromFloat(0.1),  // 10%
			StopLossPercentage:   decimal.NewFromFloat(0.15), // 15%
			TakeProfitPercentage: decimal.NewFromFloat(0.3),  // 30%
			AllowedStrategies:    []string{"momentum", "arbitrage"},
		}

		assert.Equal(t, "aggressive", profile.Level)
		assert.True(t, profile.MaxPositionSize.GreaterThan(decimal.NewFromFloat(0.1)))
		assert.Contains(t, profile.AllowedStrategies, "momentum")
	})
}
