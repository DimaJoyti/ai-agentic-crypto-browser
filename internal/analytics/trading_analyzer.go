package analytics

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/shopspring/decimal"
)

// TradingPerformanceAnalyzer analyzes trading performance metrics
type TradingPerformanceAnalyzer struct {
	logger    *observability.Logger
	config    PerformanceConfig
	trades    []TradeRecord
	metrics   TradingPerformanceMetrics
	mu        sync.RWMutex
	isRunning int32
}

// TradeRecord represents a single trade record
type TradeRecord struct {
	ID           string          `json:"id"`
	Symbol       string          `json:"symbol"`
	Side         string          `json:"side"`
	Quantity     decimal.Decimal `json:"quantity"`
	EntryPrice   decimal.Decimal `json:"entry_price"`
	ExitPrice    decimal.Decimal `json:"exit_price"`
	EntryTime    time.Time       `json:"entry_time"`
	ExitTime     time.Time       `json:"exit_time"`
	PnL          decimal.Decimal `json:"pnl"`
	Commission   decimal.Decimal `json:"commission"`
	Slippage     decimal.Decimal `json:"slippage"`
	Strategy     string          `json:"strategy"`
	Duration     time.Duration   `json:"duration"`
	MaxFavorable decimal.Decimal `json:"max_favorable"`
	MaxAdverse   decimal.Decimal `json:"max_adverse"`
	IsWin        bool            `json:"is_win"`
}

// StrategyPerformance represents performance metrics for a specific strategy
type StrategyPerformance struct {
	StrategyName         string          `json:"strategy_name"`
	TotalTrades          int             `json:"total_trades"`
	WinningTrades        int             `json:"winning_trades"`
	LosingTrades         int             `json:"losing_trades"`
	WinRate              float64         `json:"win_rate"`
	TotalPnL             decimal.Decimal `json:"total_pnl"`
	AverageWin           decimal.Decimal `json:"average_win"`
	AverageLoss          decimal.Decimal `json:"average_loss"`
	ProfitFactor         float64         `json:"profit_factor"`
	SharpeRatio          float64         `json:"sharpe_ratio"`
	MaxDrawdown          decimal.Decimal `json:"max_drawdown"`
	RecoveryFactor       float64         `json:"recovery_factor"`
	Expectancy           decimal.Decimal `json:"expectancy"`
	AverageTradeDuration time.Duration   `json:"average_trade_duration"`
	MaxConsecutiveWins   int             `json:"max_consecutive_wins"`
	MaxConsecutiveLosses int             `json:"max_consecutive_losses"`
}

// NewTradingPerformanceAnalyzer creates a new trading performance analyzer
func NewTradingPerformanceAnalyzer(logger *observability.Logger, config PerformanceConfig) *TradingPerformanceAnalyzer {
	return &TradingPerformanceAnalyzer{
		logger: logger,
		config: config,
		trades: make([]TradeRecord, 0, config.MetricsBufferSize),
	}
}

// Start starts the trading performance analyzer
func (tpa *TradingPerformanceAnalyzer) Start(ctx context.Context) error {
	tpa.logger.Info(ctx, "Starting trading performance analyzer", nil)
	tpa.isRunning = 1
	return nil
}

// Stop stops the trading performance analyzer
func (tpa *TradingPerformanceAnalyzer) Stop(ctx context.Context) error {
	tpa.logger.Info(ctx, "Stopping trading performance analyzer", nil)
	tpa.isRunning = 0
	return nil
}

// AddTrade adds a new trade record for analysis
func (tpa *TradingPerformanceAnalyzer) AddTrade(trade TradeRecord) {
	tpa.mu.Lock()
	defer tpa.mu.Unlock()

	// Add trade to buffer
	tpa.trades = append(tpa.trades, trade)

	// Maintain buffer size
	if len(tpa.trades) > tpa.config.MetricsBufferSize {
		tpa.trades = tpa.trades[1:]
	}

	// Recalculate metrics
	tpa.calculateMetrics()
}

// calculateMetrics calculates comprehensive trading performance metrics
func (tpa *TradingPerformanceAnalyzer) calculateMetrics() {
	if len(tpa.trades) == 0 {
		return
	}

	var totalPnL, totalWins, totalLosses decimal.Decimal
	var winningTrades, losingTrades int64
	var totalVolume decimal.Decimal
	var returns []float64
	var consecutiveWins, consecutiveLosses, maxConsecutiveWins, maxConsecutiveLosses int

	for _, trade := range tpa.trades {
		// Basic counts
		if trade.IsWin {
			winningTrades++
			totalWins = totalWins.Add(trade.PnL)
			consecutiveWins++
			consecutiveLosses = 0
			if consecutiveWins > maxConsecutiveWins {
				maxConsecutiveWins = consecutiveWins
			}
		} else {
			losingTrades++
			totalLosses = totalLosses.Add(trade.PnL.Abs())
			consecutiveLosses++
			consecutiveWins = 0
			if consecutiveLosses > maxConsecutiveLosses {
				maxConsecutiveLosses = consecutiveLosses
			}
		}

		totalPnL = totalPnL.Add(trade.PnL)
		totalVolume = totalVolume.Add(trade.Quantity.Mul(trade.EntryPrice))

		// Calculate return for this trade
		if !trade.EntryPrice.IsZero() {
			returnPct, _ := trade.PnL.Div(trade.Quantity.Mul(trade.EntryPrice)).Float64()
			returns = append(returns, returnPct)
		}
	}

	totalTrades := int64(len(tpa.trades))
	successRate := float64(winningTrades) / float64(totalTrades) * 100

	// Calculate average trade size
	var averageTradeSize decimal.Decimal
	if totalTrades > 0 {
		averageTradeSize = totalVolume.Div(decimal.NewFromInt(totalTrades))
	}

	// Calculate profit factor
	var profitFactor float64
	if !totalLosses.IsZero() {
		profitFactor, _ = totalWins.Div(totalLosses).Float64()
	}

	// Calculate Sharpe ratio
	sharpeRatio := tpa.calculateSharpeRatio(returns)

	// Calculate Sortino ratio
	sortinoRatio := tpa.calculateSortinoRatio(returns)

	// Calculate maximum drawdown
	maxDrawdown := tpa.calculateMaxDrawdown()

	// Calculate current drawdown
	currentDrawdown := tpa.calculateCurrentDrawdown()

	// Calculate average win/loss
	var averageWin, averageLoss decimal.Decimal
	if winningTrades > 0 {
		averageWin = totalWins.Div(decimal.NewFromInt(winningTrades))
	}
	if losingTrades > 0 {
		averageLoss = totalLosses.Div(decimal.NewFromInt(losingTrades))
	}

	// Find largest win/loss
	var largestWin, largestLoss decimal.Decimal
	for _, trade := range tpa.trades {
		if trade.IsWin && trade.PnL.GreaterThan(largestWin) {
			largestWin = trade.PnL
		}
		if !trade.IsWin && trade.PnL.Abs().GreaterThan(largestLoss) {
			largestLoss = trade.PnL.Abs()
		}
	}

	// Calculate trading frequency (trades per day)
	var tradingFrequency float64
	if len(tpa.trades) > 1 {
		firstTrade := tpa.trades[0].EntryTime
		lastTrade := tpa.trades[len(tpa.trades)-1].EntryTime
		daysDiff := lastTrade.Sub(firstTrade).Hours() / 24
		if daysDiff > 0 {
			tradingFrequency = float64(totalTrades) / daysDiff
		}
	}

	// Calculate volatility
	volatility := tpa.calculateVolatility(returns)

	// Calculate beta and alpha (simplified - would need benchmark data)
	beta := 1.0   // Mock value
	alpha := 0.02 // Mock value

	// Calculate information ratio
	informationRatio := tpa.calculateInformationRatio(returns)

	// Calculate Calmar ratio
	calmarRatio := tpa.calculateCalmarRatio(returns, maxDrawdown)

	// Update metrics
	tpa.metrics = TradingPerformanceMetrics{
		TotalTrades:       totalTrades,
		SuccessfulTrades:  winningTrades,
		FailedTrades:      losingTrades,
		SuccessRate:       successRate,
		AverageTradeSize:  averageTradeSize,
		TotalVolume:       totalVolume,
		TotalPnL:          totalPnL,
		WinRate:           float64(winningTrades) / float64(totalTrades) * 100,
		ProfitFactor:      profitFactor,
		SharpeRatio:       sharpeRatio,
		SortinoRatio:      sortinoRatio,
		MaxDrawdown:       maxDrawdown,
		CurrentDrawdown:   currentDrawdown,
		AverageWin:        averageWin,
		AverageLoss:       averageLoss,
		LargestWin:        largestWin,
		LargestLoss:       largestLoss,
		ConsecutiveWins:   maxConsecutiveWins,
		ConsecutiveLosses: maxConsecutiveLosses,
		TradingFrequency:  tradingFrequency,
		Volatility:        volatility,
		Beta:              beta,
		Alpha:             alpha,
		InformationRatio:  informationRatio,
		CalmarRatio:       calmarRatio,
	}
}

// calculateSharpeRatio calculates the Sharpe ratio
func (tpa *TradingPerformanceAnalyzer) calculateSharpeRatio(returns []float64) float64 {
	if len(returns) < 2 {
		return 0
	}

	// Calculate mean return
	var sum float64
	for _, ret := range returns {
		sum += ret
	}
	mean := sum / float64(len(returns))

	// Calculate standard deviation
	var variance float64
	for _, ret := range returns {
		variance += math.Pow(ret-mean, 2)
	}
	variance /= float64(len(returns) - 1)
	stdDev := math.Sqrt(variance)

	if stdDev == 0 {
		return 0
	}

	// Assume risk-free rate of 2% annually (0.02/252 daily)
	riskFreeRate := 0.02 / 252

	return (mean - riskFreeRate) / stdDev * math.Sqrt(252) // Annualized
}

// calculateSortinoRatio calculates the Sortino ratio
func (tpa *TradingPerformanceAnalyzer) calculateSortinoRatio(returns []float64) float64 {
	if len(returns) < 2 {
		return 0
	}

	// Calculate mean return
	var sum float64
	for _, ret := range returns {
		sum += ret
	}
	mean := sum / float64(len(returns))

	// Calculate downside deviation
	var downsideVariance float64
	downsideCount := 0
	for _, ret := range returns {
		if ret < 0 {
			downsideVariance += math.Pow(ret, 2)
			downsideCount++
		}
	}

	if downsideCount == 0 {
		return math.Inf(1) // No downside risk
	}

	downsideVariance /= float64(downsideCount)
	downsideStdDev := math.Sqrt(downsideVariance)

	if downsideStdDev == 0 {
		return 0
	}

	riskFreeRate := 0.02 / 252
	return (mean - riskFreeRate) / downsideStdDev * math.Sqrt(252)
}

// calculateMaxDrawdown calculates the maximum drawdown
func (tpa *TradingPerformanceAnalyzer) calculateMaxDrawdown() decimal.Decimal {
	if len(tpa.trades) == 0 {
		return decimal.Zero
	}

	var runningPnL, peak, maxDrawdown decimal.Decimal
	peak = decimal.Zero

	for _, trade := range tpa.trades {
		runningPnL = runningPnL.Add(trade.PnL)

		if runningPnL.GreaterThan(peak) {
			peak = runningPnL
		}

		drawdown := peak.Sub(runningPnL)
		if drawdown.GreaterThan(maxDrawdown) {
			maxDrawdown = drawdown
		}
	}

	return maxDrawdown
}

// calculateCurrentDrawdown calculates the current drawdown
func (tpa *TradingPerformanceAnalyzer) calculateCurrentDrawdown() decimal.Decimal {
	if len(tpa.trades) == 0 {
		return decimal.Zero
	}

	var runningPnL, peak decimal.Decimal
	peak = decimal.Zero

	for _, trade := range tpa.trades {
		runningPnL = runningPnL.Add(trade.PnL)
		if runningPnL.GreaterThan(peak) {
			peak = runningPnL
		}
	}

	return peak.Sub(runningPnL)
}

// calculateVolatility calculates the volatility of returns
func (tpa *TradingPerformanceAnalyzer) calculateVolatility(returns []float64) float64 {
	if len(returns) < 2 {
		return 0
	}

	var sum float64
	for _, ret := range returns {
		sum += ret
	}
	mean := sum / float64(len(returns))

	var variance float64
	for _, ret := range returns {
		variance += math.Pow(ret-mean, 2)
	}
	variance /= float64(len(returns) - 1)

	return math.Sqrt(variance) * math.Sqrt(252) // Annualized
}

// calculateInformationRatio calculates the information ratio
func (tpa *TradingPerformanceAnalyzer) calculateInformationRatio(returns []float64) float64 {
	if len(returns) < 2 {
		return 0
	}

	// Simplified calculation - would need benchmark returns in practice
	var sum float64
	for _, ret := range returns {
		sum += ret
	}
	mean := sum / float64(len(returns))

	var variance float64
	for _, ret := range returns {
		variance += math.Pow(ret-mean, 2)
	}
	variance /= float64(len(returns) - 1)
	trackingError := math.Sqrt(variance)

	if trackingError == 0 {
		return 0
	}

	// Assume benchmark return of 0 for simplification
	return mean / trackingError
}

// calculateCalmarRatio calculates the Calmar ratio
func (tpa *TradingPerformanceAnalyzer) calculateCalmarRatio(returns []float64, maxDrawdown decimal.Decimal) float64 {
	if len(returns) == 0 || maxDrawdown.IsZero() {
		return 0
	}

	var sum float64
	for _, ret := range returns {
		sum += ret
	}
	annualizedReturn := sum / float64(len(returns)) * 252

	maxDrawdownFloat, _ := maxDrawdown.Float64()
	if maxDrawdownFloat == 0 {
		return math.Inf(1)
	}

	return annualizedReturn / maxDrawdownFloat
}

// GetMetrics returns current trading performance metrics
func (tpa *TradingPerformanceAnalyzer) GetMetrics() TradingPerformanceMetrics {
	tpa.mu.RLock()
	defer tpa.mu.RUnlock()
	return tpa.metrics
}

// GetStrategyPerformance returns performance metrics for a specific strategy
func (tpa *TradingPerformanceAnalyzer) GetStrategyPerformance(strategyName string) StrategyPerformance {
	tpa.mu.RLock()
	defer tpa.mu.RUnlock()

	var strategyTrades []TradeRecord
	for _, trade := range tpa.trades {
		if trade.Strategy == strategyName {
			strategyTrades = append(strategyTrades, trade)
		}
	}

	return tpa.calculateStrategyMetrics(strategyName, strategyTrades)
}

// calculateStrategyMetrics calculates performance metrics for a specific strategy
func (tpa *TradingPerformanceAnalyzer) calculateStrategyMetrics(strategyName string, trades []TradeRecord) StrategyPerformance {
	if len(trades) == 0 {
		return StrategyPerformance{StrategyName: strategyName}
	}

	var totalPnL, totalWins, totalLosses decimal.Decimal
	var winningTrades, losingTrades int
	var totalDuration time.Duration
	var consecutiveWins, consecutiveLosses, maxConsecutiveWins, maxConsecutiveLosses int
	var returns []float64

	for _, trade := range trades {
		totalPnL = totalPnL.Add(trade.PnL)
		totalDuration += trade.Duration

		if trade.IsWin {
			winningTrades++
			totalWins = totalWins.Add(trade.PnL)
			consecutiveWins++
			consecutiveLosses = 0
			if consecutiveWins > maxConsecutiveWins {
				maxConsecutiveWins = consecutiveWins
			}
		} else {
			losingTrades++
			totalLosses = totalLosses.Add(trade.PnL.Abs())
			consecutiveLosses++
			consecutiveWins = 0
			if consecutiveLosses > maxConsecutiveLosses {
				maxConsecutiveLosses = consecutiveLosses
			}
		}

		// Calculate return for Sharpe ratio
		if !trade.EntryPrice.IsZero() {
			returnPct, _ := trade.PnL.Div(trade.Quantity.Mul(trade.EntryPrice)).Float64()
			returns = append(returns, returnPct)
		}
	}

	totalTrades := len(trades)
	winRate := float64(winningTrades) / float64(totalTrades) * 100

	var averageWin, averageLoss decimal.Decimal
	if winningTrades > 0 {
		averageWin = totalWins.Div(decimal.NewFromInt(int64(winningTrades)))
	}
	if losingTrades > 0 {
		averageLoss = totalLosses.Div(decimal.NewFromInt(int64(losingTrades)))
	}

	var profitFactor float64
	if !totalLosses.IsZero() {
		profitFactor, _ = totalWins.Div(totalLosses).Float64()
	}

	sharpeRatio := tpa.calculateSharpeRatio(returns)
	maxDrawdown := tpa.calculateStrategyMaxDrawdown(trades)

	var recoveryFactor float64
	if !maxDrawdown.IsZero() {
		totalPnLFloat, _ := totalPnL.Float64()
		maxDrawdownFloat, _ := maxDrawdown.Float64()
		recoveryFactor = totalPnLFloat / maxDrawdownFloat
	}

	var expectancy decimal.Decimal
	if totalTrades > 0 {
		expectancy = totalPnL.Div(decimal.NewFromInt(int64(totalTrades)))
	}

	var averageTradeDuration time.Duration
	if totalTrades > 0 {
		averageTradeDuration = totalDuration / time.Duration(totalTrades)
	}

	return StrategyPerformance{
		StrategyName:         strategyName,
		TotalTrades:          totalTrades,
		WinningTrades:        winningTrades,
		LosingTrades:         losingTrades,
		WinRate:              winRate,
		TotalPnL:             totalPnL,
		AverageWin:           averageWin,
		AverageLoss:          averageLoss,
		ProfitFactor:         profitFactor,
		SharpeRatio:          sharpeRatio,
		MaxDrawdown:          maxDrawdown,
		RecoveryFactor:       recoveryFactor,
		Expectancy:           expectancy,
		AverageTradeDuration: averageTradeDuration,
		MaxConsecutiveWins:   maxConsecutiveWins,
		MaxConsecutiveLosses: maxConsecutiveLosses,
	}
}

// calculateStrategyMaxDrawdown calculates max drawdown for a specific strategy
func (tpa *TradingPerformanceAnalyzer) calculateStrategyMaxDrawdown(trades []TradeRecord) decimal.Decimal {
	if len(trades) == 0 {
		return decimal.Zero
	}

	var runningPnL, peak, maxDrawdown decimal.Decimal
	peak = decimal.Zero

	for _, trade := range trades {
		runningPnL = runningPnL.Add(trade.PnL)

		if runningPnL.GreaterThan(peak) {
			peak = runningPnL
		}

		drawdown := peak.Sub(runningPnL)
		if drawdown.GreaterThan(maxDrawdown) {
			maxDrawdown = drawdown
		}
	}

	return maxDrawdown
}
