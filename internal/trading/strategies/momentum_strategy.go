package strategies

import (
	"context"
	"fmt"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/shopspring/decimal"
)

// MomentumStrategy implements Momentum trading strategy
type MomentumStrategy struct {
	logger          *observability.Logger
	config          *MomentumConfig
	priceHistory    []decimal.Decimal
	volumeHistory   []decimal.Decimal
	rsiValues       []decimal.Decimal
	currentPosition *Position
	lastSignal      time.Time
}

// MomentumConfig holds configuration for Momentum strategy
type MomentumConfig struct {
	MomentumPeriod    int             `yaml:"momentum_period"`
	RSIThresholdBuy   decimal.Decimal `yaml:"rsi_threshold_buy"`
	RSIThresholdSell  decimal.Decimal `yaml:"rsi_threshold_sell"`
	VolumeThreshold   decimal.Decimal `yaml:"volume_threshold"`
	BreakoutThreshold decimal.Decimal `yaml:"breakout_threshold"`
	TradingPairs      []string        `yaml:"trading_pairs"`
	Exchange          string          `yaml:"exchange"`
	PositionSize      decimal.Decimal `yaml:"position_size"`
	StopLoss          decimal.Decimal `yaml:"stop_loss"`
	TakeProfit        decimal.Decimal `yaml:"take_profit"`
}

// Position represents a trading position
type Position struct {
	ID            string          `json:"id"`
	Symbol        string          `json:"symbol"`
	Side          PositionSide    `json:"side"`
	Size          decimal.Decimal `json:"size"`
	EntryPrice    decimal.Decimal `json:"entry_price"`
	CurrentPrice  decimal.Decimal `json:"current_price"`
	UnrealizedPnL decimal.Decimal `json:"unrealized_pnl"`
	OpenTime      time.Time       `json:"open_time"`
	StopLoss      decimal.Decimal `json:"stop_loss"`
	TakeProfit    decimal.Decimal `json:"take_profit"`
}

// PositionSide represents the side of a position
type PositionSide string

const (
	PositionSideLong  PositionSide = "long"
	PositionSideShort PositionSide = "short"
)

// MomentumSignal represents a momentum trading signal
type MomentumSignal struct {
	*TradingSignal
	RSI           decimal.Decimal `json:"rsi"`
	Momentum      decimal.Decimal `json:"momentum"`
	VolumeRatio   decimal.Decimal `json:"volume_ratio"`
	BreakoutLevel decimal.Decimal `json:"breakout_level"`
	Confidence    decimal.Decimal `json:"confidence"`
}

// NewMomentumStrategy creates a new Momentum strategy instance
func NewMomentumStrategy(logger *observability.Logger, config *MomentumConfig) *MomentumStrategy {
	return &MomentumStrategy{
		logger:        logger,
		config:        config,
		priceHistory:  make([]decimal.Decimal, 0),
		volumeHistory: make([]decimal.Decimal, 0),
		rsiValues:     make([]decimal.Decimal, 0),
	}
}

// Execute executes the Momentum strategy
func (ms *MomentumStrategy) Execute(ctx context.Context, marketData *MarketData) (interface{}, error) {
	// Update price and volume history
	ms.updateHistory(marketData)

	// Check if we have enough data
	if len(ms.priceHistory) < ms.config.MomentumPeriod {
		return nil, nil // Not enough data yet
	}

	// Calculate technical indicators
	rsi := ms.calculateRSI()
	momentum := ms.calculateMomentum()
	volumeRatio := ms.calculateVolumeRatio()

	// Update RSI history
	ms.rsiValues = append(ms.rsiValues, rsi)
	if len(ms.rsiValues) > ms.config.MomentumPeriod {
		ms.rsiValues = ms.rsiValues[1:]
	}

	// Check for trading signals
	signal := ms.generateSignal(ctx, marketData, rsi, momentum, volumeRatio)

	// Update current position if we have one
	if ms.currentPosition != nil {
		ms.updatePosition(marketData.Price)

		// Check for exit conditions
		if exitSignal := ms.checkExitConditions(ctx, marketData); exitSignal != nil {
			return exitSignal, nil
		}
	}

	return signal, nil
}

// updateHistory updates price and volume history
func (ms *MomentumStrategy) updateHistory(marketData *MarketData) {
	ms.priceHistory = append(ms.priceHistory, marketData.Price)
	ms.volumeHistory = append(ms.volumeHistory, marketData.Volume)

	// Keep only the required number of periods
	maxPeriods := ms.config.MomentumPeriod + 10 // Keep some extra for calculations
	if len(ms.priceHistory) > maxPeriods {
		ms.priceHistory = ms.priceHistory[len(ms.priceHistory)-maxPeriods:]
	}
	if len(ms.volumeHistory) > maxPeriods {
		ms.volumeHistory = ms.volumeHistory[len(ms.volumeHistory)-maxPeriods:]
	}
}

// calculateRSI calculates the Relative Strength Index
func (ms *MomentumStrategy) calculateRSI() decimal.Decimal {
	if len(ms.priceHistory) < ms.config.MomentumPeriod+1 {
		return decimal.NewFromFloat(50) // Neutral RSI
	}

	gains := decimal.Zero
	losses := decimal.Zero
	period := ms.config.MomentumPeriod

	// Calculate gains and losses
	for i := len(ms.priceHistory) - period; i < len(ms.priceHistory); i++ {
		if i == 0 {
			continue
		}

		change := ms.priceHistory[i].Sub(ms.priceHistory[i-1])
		if change.GreaterThan(decimal.Zero) {
			gains = gains.Add(change)
		} else {
			losses = losses.Add(change.Abs())
		}
	}

	if losses.IsZero() {
		return decimal.NewFromFloat(100)
	}

	avgGain := gains.Div(decimal.NewFromInt(int64(period)))
	avgLoss := losses.Div(decimal.NewFromInt(int64(period)))

	rs := avgGain.Div(avgLoss)
	rsi := decimal.NewFromFloat(100).Sub(decimal.NewFromFloat(100).Div(decimal.NewFromFloat(1).Add(rs)))

	return rsi
}

// calculateMomentum calculates price momentum
func (ms *MomentumStrategy) calculateMomentum() decimal.Decimal {
	if len(ms.priceHistory) < ms.config.MomentumPeriod {
		return decimal.Zero
	}

	currentPrice := ms.priceHistory[len(ms.priceHistory)-1]
	pastPrice := ms.priceHistory[len(ms.priceHistory)-ms.config.MomentumPeriod]

	return currentPrice.Sub(pastPrice).Div(pastPrice).Mul(decimal.NewFromFloat(100))
}

// calculateVolumeRatio calculates volume ratio compared to average
func (ms *MomentumStrategy) calculateVolumeRatio() decimal.Decimal {
	if len(ms.volumeHistory) < ms.config.MomentumPeriod {
		return decimal.NewFromFloat(1)
	}

	currentVolume := ms.volumeHistory[len(ms.volumeHistory)-1]

	// Calculate average volume
	totalVolume := decimal.Zero
	for i := len(ms.volumeHistory) - ms.config.MomentumPeriod; i < len(ms.volumeHistory)-1; i++ {
		totalVolume = totalVolume.Add(ms.volumeHistory[i])
	}

	avgVolume := totalVolume.Div(decimal.NewFromInt(int64(ms.config.MomentumPeriod - 1)))

	if avgVolume.IsZero() {
		return decimal.NewFromFloat(1)
	}

	return currentVolume.Div(avgVolume)
}

// generateSignal generates trading signals based on momentum indicators
func (ms *MomentumStrategy) generateSignal(ctx context.Context, marketData *MarketData, rsi, momentum, volumeRatio decimal.Decimal) *MomentumSignal {
	// Don't generate signals too frequently
	if time.Since(ms.lastSignal) < time.Minute*5 {
		return nil
	}

	// Check for buy signal
	if ms.currentPosition == nil && ms.shouldBuy(rsi, momentum, volumeRatio) {
		return ms.createBuySignal(ctx, marketData, rsi, momentum, volumeRatio)
	}

	// Check for sell signal
	if ms.currentPosition != nil && ms.currentPosition.Side == PositionSideLong && ms.shouldSell(rsi, momentum, volumeRatio) {
		return ms.createSellSignal(ctx, marketData, rsi, momentum, volumeRatio)
	}

	return nil
}

// shouldBuy determines if we should buy based on momentum indicators
func (ms *MomentumStrategy) shouldBuy(rsi, momentum, volumeRatio decimal.Decimal) bool {
	// Buy conditions:
	// 1. RSI is oversold but starting to recover
	// 2. Positive momentum
	// 3. Volume spike
	// 4. Price breakout

	rsiCondition := rsi.LessThan(ms.config.RSIThresholdBuy) && ms.isRSIRecovering()
	momentumCondition := momentum.GreaterThan(decimal.Zero)
	volumeCondition := volumeRatio.GreaterThan(ms.config.VolumeThreshold)
	breakoutCondition := ms.isBreakout()

	// Require at least 3 out of 4 conditions
	conditionCount := 0
	if rsiCondition {
		conditionCount++
	}
	if momentumCondition {
		conditionCount++
	}
	if volumeCondition {
		conditionCount++
	}
	if breakoutCondition {
		conditionCount++
	}

	return conditionCount >= 3
}

// shouldSell determines if we should sell based on momentum indicators
func (ms *MomentumStrategy) shouldSell(rsi, momentum, volumeRatio decimal.Decimal) bool {
	// Sell conditions:
	// 1. RSI is overbought
	// 2. Negative momentum
	// 3. Volume declining

	rsiCondition := rsi.GreaterThan(ms.config.RSIThresholdSell)
	momentumCondition := momentum.LessThan(decimal.Zero)
	volumeCondition := volumeRatio.LessThan(decimal.NewFromFloat(0.8))

	return rsiCondition || (momentumCondition && volumeCondition)
}

// isRSIRecovering checks if RSI is recovering from oversold levels
func (ms *MomentumStrategy) isRSIRecovering() bool {
	if len(ms.rsiValues) < 3 {
		return false
	}

	recent := ms.rsiValues[len(ms.rsiValues)-1]
	previous := ms.rsiValues[len(ms.rsiValues)-2]

	return recent.GreaterThan(previous)
}

// isBreakout checks if price is breaking out
func (ms *MomentumStrategy) isBreakout() bool {
	if len(ms.priceHistory) < 20 {
		return false
	}

	currentPrice := ms.priceHistory[len(ms.priceHistory)-1]

	// Calculate recent high
	recentHigh := decimal.Zero
	for i := len(ms.priceHistory) - 20; i < len(ms.priceHistory)-1; i++ {
		if ms.priceHistory[i].GreaterThan(recentHigh) {
			recentHigh = ms.priceHistory[i]
		}
	}

	breakoutLevel := recentHigh.Mul(decimal.NewFromFloat(1).Add(ms.config.BreakoutThreshold))
	return currentPrice.GreaterThan(breakoutLevel)
}

// createBuySignal creates a buy signal
func (ms *MomentumStrategy) createBuySignal(ctx context.Context, marketData *MarketData, rsi, momentum, volumeRatio decimal.Decimal) *MomentumSignal {
	orderID := fmt.Sprintf("momentum_buy_%d", time.Now().UnixNano())

	signal := &MomentumSignal{
		TradingSignal: &TradingSignal{
			ID:        orderID,
			Symbol:    marketData.Symbol,
			Action:    ActionBuy,
			Amount:    ms.config.PositionSize,
			Price:     marketData.Price,
			OrderType: OrderTypeMarket,
			Timestamp: time.Now(),
			Strategy:  "Momentum",
			Metadata: map[string]interface{}{
				"rsi":          rsi.String(),
				"momentum":     momentum.String(),
				"volume_ratio": volumeRatio.String(),
			},
		},
		RSI:         rsi,
		Momentum:    momentum,
		VolumeRatio: volumeRatio,
		Confidence:  ms.calculateConfidence(rsi, momentum, volumeRatio),
	}

	// Create position
	ms.currentPosition = &Position{
		ID:         orderID,
		Symbol:     marketData.Symbol,
		Side:       PositionSideLong,
		Size:       ms.config.PositionSize,
		EntryPrice: marketData.Price,
		OpenTime:   time.Now(),
		StopLoss:   marketData.Price.Mul(decimal.NewFromFloat(1).Sub(ms.config.StopLoss)),
		TakeProfit: marketData.Price.Mul(decimal.NewFromFloat(1).Add(ms.config.TakeProfit)),
	}

	ms.lastSignal = time.Now()

	ms.logger.Info(ctx, "Momentum buy signal generated", map[string]interface{}{
		"symbol":       signal.Symbol,
		"price":        signal.Price.String(),
		"rsi":          rsi.String(),
		"momentum":     momentum.String(),
		"volume_ratio": volumeRatio.String(),
		"confidence":   signal.Confidence.String(),
	})

	return signal
}

// createSellSignal creates a sell signal
func (ms *MomentumStrategy) createSellSignal(ctx context.Context, marketData *MarketData, rsi, momentum, volumeRatio decimal.Decimal) *MomentumSignal {
	orderID := fmt.Sprintf("momentum_sell_%d", time.Now().UnixNano())

	signal := &MomentumSignal{
		TradingSignal: &TradingSignal{
			ID:        orderID,
			Symbol:    marketData.Symbol,
			Action:    ActionSell,
			Amount:    ms.currentPosition.Size,
			Price:     marketData.Price,
			OrderType: OrderTypeMarket,
			Timestamp: time.Now(),
			Strategy:  "Momentum",
			Metadata: map[string]interface{}{
				"rsi":          rsi.String(),
				"momentum":     momentum.String(),
				"volume_ratio": volumeRatio.String(),
				"position_id":  ms.currentPosition.ID,
			},
		},
		RSI:         rsi,
		Momentum:    momentum,
		VolumeRatio: volumeRatio,
		Confidence:  ms.calculateConfidence(rsi, momentum, volumeRatio),
	}

	// Close position
	ms.currentPosition = nil
	ms.lastSignal = time.Now()

	ms.logger.Info(ctx, "Momentum sell signal generated", map[string]interface{}{
		"symbol":       signal.Symbol,
		"price":        signal.Price.String(),
		"rsi":          rsi.String(),
		"momentum":     momentum.String(),
		"volume_ratio": volumeRatio.String(),
		"confidence":   signal.Confidence.String(),
	})

	return signal
}

// calculateConfidence calculates signal confidence
func (ms *MomentumStrategy) calculateConfidence(rsi, momentum, volumeRatio decimal.Decimal) decimal.Decimal {
	confidence := decimal.Zero

	// RSI contribution (0-30)
	if rsi.LessThan(decimal.NewFromFloat(30)) || rsi.GreaterThan(decimal.NewFromFloat(70)) {
		confidence = confidence.Add(decimal.NewFromFloat(30))
	}

	// Momentum contribution (0-40)
	momentumAbs := momentum.Abs()
	if momentumAbs.GreaterThan(decimal.NewFromFloat(5)) {
		confidence = confidence.Add(decimal.NewFromFloat(40))
	} else if momentumAbs.GreaterThan(decimal.NewFromFloat(2)) {
		confidence = confidence.Add(decimal.NewFromFloat(20))
	}

	// Volume contribution (0-30)
	if volumeRatio.GreaterThan(decimal.NewFromFloat(2)) {
		confidence = confidence.Add(decimal.NewFromFloat(30))
	} else if volumeRatio.GreaterThan(decimal.NewFromFloat(1.5)) {
		confidence = confidence.Add(decimal.NewFromFloat(15))
	}

	return confidence
}

// updatePosition updates the current position with new price
func (ms *MomentumStrategy) updatePosition(currentPrice decimal.Decimal) {
	if ms.currentPosition == nil {
		return
	}

	ms.currentPosition.CurrentPrice = currentPrice

	if ms.currentPosition.Side == PositionSideLong {
		ms.currentPosition.UnrealizedPnL = currentPrice.Sub(ms.currentPosition.EntryPrice).Mul(ms.currentPosition.Size)
	}
}

// checkExitConditions checks if we should exit the current position
func (ms *MomentumStrategy) checkExitConditions(ctx context.Context, marketData *MarketData) *MomentumSignal {
	if ms.currentPosition == nil {
		return nil
	}

	currentPrice := marketData.Price

	// Check stop loss
	if currentPrice.LessThanOrEqual(ms.currentPosition.StopLoss) {
		return ms.createExitSignal(ctx, marketData, "stop_loss")
	}

	// Check take profit
	if currentPrice.GreaterThanOrEqual(ms.currentPosition.TakeProfit) {
		return ms.createExitSignal(ctx, marketData, "take_profit")
	}

	return nil
}

// createExitSignal creates an exit signal
func (ms *MomentumStrategy) createExitSignal(ctx context.Context, marketData *MarketData, reason string) *MomentumSignal {
	orderID := fmt.Sprintf("momentum_exit_%d", time.Now().UnixNano())

	signal := &MomentumSignal{
		TradingSignal: &TradingSignal{
			ID:        orderID,
			Symbol:    marketData.Symbol,
			Action:    ActionSell,
			Amount:    ms.currentPosition.Size,
			Price:     marketData.Price,
			OrderType: OrderTypeMarket,
			Timestamp: time.Now(),
			Strategy:  "Momentum",
			Metadata: map[string]interface{}{
				"exit_reason": reason,
				"position_id": ms.currentPosition.ID,
				"entry_price": ms.currentPosition.EntryPrice.String(),
				"pnl":         ms.currentPosition.UnrealizedPnL.String(),
			},
		},
	}

	ms.logger.Info(ctx, "Momentum exit signal generated", map[string]interface{}{
		"symbol":      signal.Symbol,
		"price":       signal.Price.String(),
		"exit_reason": reason,
		"pnl":         ms.currentPosition.UnrealizedPnL.String(),
	})

	// Close position
	ms.currentPosition = nil

	return signal
}

// GetPerformance returns performance metrics for the Momentum strategy
func (ms *MomentumStrategy) GetPerformance() *StrategyPerformance {
	// Implementation for performance metrics
	return &StrategyPerformance{}
}

// Reset resets the strategy state
func (ms *MomentumStrategy) Reset() {
	ms.priceHistory = make([]decimal.Decimal, 0)
	ms.volumeHistory = make([]decimal.Decimal, 0)
	ms.rsiValues = make([]decimal.Decimal, 0)
	ms.currentPosition = nil
	ms.lastSignal = time.Time{}
}

// GetName returns the strategy name
func (ms *MomentumStrategy) GetName() string {
	return "Momentum Trading"
}

// GetType returns the strategy type
func (ms *MomentumStrategy) GetType() string {
	return "Momentum"
}

// Validate validates the strategy configuration
func (ms *MomentumStrategy) Validate() error {
	if ms.config.MomentumPeriod < 5 {
		return fmt.Errorf("momentum period must be at least 5")
	}

	if ms.config.RSIThresholdBuy.LessThan(decimal.Zero) || ms.config.RSIThresholdBuy.GreaterThan(decimal.NewFromFloat(100)) {
		return fmt.Errorf("RSI buy threshold must be between 0 and 100")
	}

	if ms.config.RSIThresholdSell.LessThan(decimal.Zero) || ms.config.RSIThresholdSell.GreaterThan(decimal.NewFromFloat(100)) {
		return fmt.Errorf("RSI sell threshold must be between 0 and 100")
	}

	if ms.config.PositionSize.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("position size must be positive")
	}

	return nil
}
