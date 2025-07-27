package web3

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// MomentumStrategy implements a momentum-based trading strategy
type MomentumStrategy struct {
	name        string
	description string
	enabled     bool
	parameters  map[string]interface{}
	riskLevel   RiskLevel
}

// NewMomentumStrategy creates a new momentum trading strategy
func NewMomentumStrategy() *MomentumStrategy {
	return &MomentumStrategy{
		name:        "momentum_strategy",
		description: "Trades based on price momentum and technical indicators",
		enabled:     true,
		riskLevel:   RiskLevelMedium,
		parameters: map[string]interface{}{
			"rsi_oversold":    30.0,
			"rsi_overbought":  70.0,
			"momentum_threshold": 0.05, // 5% price change
			"volume_multiplier": 1.5,   // 1.5x average volume
			"confidence_threshold": 0.7,
		},
	}
}

func (m *MomentumStrategy) GetName() string {
	return m.name
}

func (m *MomentumStrategy) GetDescription() string {
	return m.description
}

func (m *MomentumStrategy) GetRiskLevel() RiskLevel {
	return m.riskLevel
}

func (m *MomentumStrategy) IsEnabled() bool {
	return m.enabled
}

func (m *MomentumStrategy) GetParameters() map[string]interface{} {
	return m.parameters
}

func (m *MomentumStrategy) Analyze(ctx context.Context, market *MarketData) (*TradingSignal, error) {
	if market.TechnicalData == nil {
		return &TradingSignal{Action: ActionHold}, nil
	}

	// Get parameters
	rsiOversold := m.parameters["rsi_oversold"].(float64)
	rsiOverbought := m.parameters["rsi_overbought"].(float64)
	momentumThreshold := m.parameters["momentum_threshold"].(float64)
	volumeMultiplier := m.parameters["volume_multiplier"].(float64)

	// Calculate signals
	rsi := market.TechnicalData.RSI.InexactFloat64()
	priceChange := market.PriceChange24h.Div(market.Price).InexactFloat64()
	volumeRatio := market.Volume24h.Div(market.TechnicalData.Volume).InexactFloat64()

	var action TradingAction = ActionHold
	var confidence float64 = 0.0
	var urgency SignalUrgency = UrgencyLow

	// Buy signals
	if rsi < rsiOversold && priceChange > momentumThreshold && volumeRatio > volumeMultiplier {
		action = ActionBuy
		confidence = 0.8
		urgency = UrgencyHigh
	} else if rsi < 40 && priceChange > momentumThreshold/2 {
		action = ActionBuy
		confidence = 0.6
		urgency = UrgencyMedium
	}

	// Sell signals
	if rsi > rsiOverbought && priceChange < -momentumThreshold && volumeRatio > volumeMultiplier {
		action = ActionSell
		confidence = 0.8
		urgency = UrgencyHigh
	} else if rsi > 60 && priceChange < -momentumThreshold/2 {
		action = ActionSell
		confidence = 0.6
		urgency = UrgencyMedium
	}

	// Check confidence threshold
	confidenceThreshold := m.parameters["confidence_threshold"].(float64)
	if confidence < confidenceThreshold {
		action = ActionHold
	}

	signal := &TradingSignal{
		ID:           uuid.New(),
		StrategyName: m.name,
		Action:       action,
		TokenIn:      "USDC", // Simplified - would be dynamic
		TokenOut:     market.TokenAddress,
		Confidence:   confidence,
		Urgency:      urgency,
		ValidUntil:   time.Now().Add(5 * time.Minute),
		CreatedAt:    time.Now(),
		Metadata: map[string]interface{}{
			"rsi":           rsi,
			"price_change":  priceChange,
			"volume_ratio":  volumeRatio,
			"market_data":   market,
		},
	}

	return signal, nil
}

func (m *MomentumStrategy) ValidateSignal(ctx context.Context, signal *TradingSignal) error {
	if signal.Confidence < 0.5 {
		return fmt.Errorf("signal confidence too low: %f", signal.Confidence)
	}

	if time.Now().After(signal.ValidUntil) {
		return fmt.Errorf("signal expired")
	}

	return nil
}

func (m *MomentumStrategy) CalculatePositionSize(ctx context.Context, signal *TradingSignal, portfolio *Portfolio) (decimal.Decimal, error) {
	// Base position size on confidence and risk profile
	baseSize := portfolio.RiskProfile.MaxPositionSize.Mul(portfolio.TotalValue)
	confidenceMultiplier := decimal.NewFromFloat(signal.Confidence)
	
	positionSize := baseSize.Mul(confidenceMultiplier)
	
	// Ensure position size doesn't exceed available balance
	if positionSize.GreaterThan(portfolio.AvailableBalance) {
		positionSize = portfolio.AvailableBalance.Mul(decimal.NewFromFloat(0.9)) // 90% of available
	}

	return positionSize, nil
}

// MeanReversionStrategy implements a mean reversion trading strategy
type MeanReversionStrategy struct {
	name        string
	description string
	enabled     bool
	parameters  map[string]interface{}
	riskLevel   RiskLevel
}

// NewMeanReversionStrategy creates a new mean reversion trading strategy
func NewMeanReversionStrategy() *MeanReversionStrategy {
	return &MeanReversionStrategy{
		name:        "mean_reversion_strategy",
		description: "Trades based on mean reversion using Bollinger Bands and RSI",
		enabled:     true,
		riskLevel:   RiskLevelLow,
		parameters: map[string]interface{}{
			"bollinger_threshold": 0.02, // 2% outside bands
			"rsi_extreme_low":     20.0,
			"rsi_extreme_high":    80.0,
			"reversion_confidence": 0.75,
		},
	}
}

func (mr *MeanReversionStrategy) GetName() string {
	return mr.name
}

func (mr *MeanReversionStrategy) GetDescription() string {
	return mr.description
}

func (mr *MeanReversionStrategy) GetRiskLevel() RiskLevel {
	return mr.riskLevel
}

func (mr *MeanReversionStrategy) IsEnabled() bool {
	return mr.enabled
}

func (mr *MeanReversionStrategy) GetParameters() map[string]interface{} {
	return mr.parameters
}

func (mr *MeanReversionStrategy) Analyze(ctx context.Context, market *MarketData) (*TradingSignal, error) {
	if market.TechnicalData == nil {
		return &TradingSignal{Action: ActionHold}, nil
	}

	// Get parameters
	bollingerThreshold := mr.parameters["bollinger_threshold"].(float64)
	rsiExtremeLow := mr.parameters["rsi_extreme_low"].(float64)
	rsiExtremeHigh := mr.parameters["rsi_extreme_high"].(float64)

	// Calculate Bollinger Band position
	price := market.Price
	upperBand := market.TechnicalData.BollingerUpper
	lowerBand := market.TechnicalData.BollingerLower
	rsi := market.TechnicalData.RSI.InexactFloat64()

	var action TradingAction = ActionHold
	var confidence float64 = 0.0
	var urgency SignalUrgency = UrgencyLow

	// Check for oversold conditions (buy signal)
	if price.LessThan(lowerBand.Mul(decimal.NewFromFloat(1-bollingerThreshold))) && rsi < rsiExtremeLow {
		action = ActionBuy
		confidence = 0.8
		urgency = UrgencyMedium
	}

	// Check for overbought conditions (sell signal)
	if price.GreaterThan(upperBand.Mul(decimal.NewFromFloat(1+bollingerThreshold))) && rsi > rsiExtremeHigh {
		action = ActionSell
		confidence = 0.8
		urgency = UrgencyMedium
	}

	signal := &TradingSignal{
		ID:           uuid.New(),
		StrategyName: mr.name,
		Action:       action,
		TokenIn:      "USDC",
		TokenOut:     market.TokenAddress,
		Confidence:   confidence,
		Urgency:      urgency,
		ValidUntil:   time.Now().Add(10 * time.Minute),
		CreatedAt:    time.Now(),
		Metadata: map[string]interface{}{
			"rsi":                rsi,
			"price":              price.String(),
			"bollinger_upper":    upperBand.String(),
			"bollinger_lower":    lowerBand.String(),
			"bollinger_position": price.Sub(lowerBand).Div(upperBand.Sub(lowerBand)).String(),
		},
	}

	return signal, nil
}

func (mr *MeanReversionStrategy) ValidateSignal(ctx context.Context, signal *TradingSignal) error {
	if signal.Confidence < 0.6 {
		return fmt.Errorf("signal confidence too low: %f", signal.Confidence)
	}

	if time.Now().After(signal.ValidUntil) {
		return fmt.Errorf("signal expired")
	}

	return nil
}

func (mr *MeanReversionStrategy) CalculatePositionSize(ctx context.Context, signal *TradingSignal, portfolio *Portfolio) (decimal.Decimal, error) {
	// Conservative position sizing for mean reversion
	baseSize := portfolio.RiskProfile.MaxPositionSize.Mul(decimal.NewFromFloat(0.5)).Mul(portfolio.TotalValue)
	confidenceMultiplier := decimal.NewFromFloat(signal.Confidence)
	
	positionSize := baseSize.Mul(confidenceMultiplier)
	
	// Ensure position size doesn't exceed available balance
	if positionSize.GreaterThan(portfolio.AvailableBalance) {
		positionSize = portfolio.AvailableBalance.Mul(decimal.NewFromFloat(0.8)) // 80% of available
	}

	return positionSize, nil
}

// ArbitrageStrategy implements a cross-DEX arbitrage strategy
type ArbitrageStrategy struct {
	name        string
	description string
	enabled     bool
	parameters  map[string]interface{}
	riskLevel   RiskLevel
}

// NewArbitrageStrategy creates a new arbitrage trading strategy
func NewArbitrageStrategy() *ArbitrageStrategy {
	return &ArbitrageStrategy{
		name:        "arbitrage_strategy",
		description: "Exploits price differences across different DEXs",
		enabled:     true,
		riskLevel:   RiskLevelLow,
		parameters: map[string]interface{}{
			"min_profit_threshold": 0.005, // 0.5% minimum profit
			"max_gas_cost_ratio":   0.3,   // Max 30% of profit for gas
			"slippage_buffer":      0.002, // 0.2% slippage buffer
		},
	}
}

func (a *ArbitrageStrategy) GetName() string {
	return a.name
}

func (a *ArbitrageStrategy) GetDescription() string {
	return a.description
}

func (a *ArbitrageStrategy) GetRiskLevel() RiskLevel {
	return a.riskLevel
}

func (a *ArbitrageStrategy) IsEnabled() bool {
	return a.enabled
}

func (a *ArbitrageStrategy) GetParameters() map[string]interface{} {
	return a.parameters
}

func (a *ArbitrageStrategy) Analyze(ctx context.Context, market *MarketData) (*TradingSignal, error) {
	// This would analyze prices across multiple DEXs
	// For now, return hold signal as this requires complex multi-DEX price comparison
	
	signal := &TradingSignal{
		ID:           uuid.New(),
		StrategyName: a.name,
		Action:       ActionHold,
		TokenIn:      "USDC",
		TokenOut:     market.TokenAddress,
		Confidence:   0.0,
		Urgency:      UrgencyLow,
		ValidUntil:   time.Now().Add(1 * time.Minute),
		CreatedAt:    time.Now(),
		Metadata: map[string]interface{}{
			"note": "Arbitrage opportunities require multi-DEX price analysis",
		},
	}

	return signal, nil
}

func (a *ArbitrageStrategy) ValidateSignal(ctx context.Context, signal *TradingSignal) error {
	return nil // Arbitrage signals are always valid if they exist
}

func (a *ArbitrageStrategy) CalculatePositionSize(ctx context.Context, signal *TradingSignal, portfolio *Portfolio) (decimal.Decimal, error) {
	// Arbitrage can use larger position sizes due to lower risk
	baseSize := portfolio.RiskProfile.MaxPositionSize.Mul(decimal.NewFromFloat(2.0)).Mul(portfolio.TotalValue)
	
	// Ensure position size doesn't exceed available balance
	if baseSize.GreaterThan(portfolio.AvailableBalance) {
		baseSize = portfolio.AvailableBalance.Mul(decimal.NewFromFloat(0.95)) // 95% of available
	}

	return baseSize, nil
}
