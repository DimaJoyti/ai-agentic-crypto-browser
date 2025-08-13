package backtesting

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/ai-agentic-browser/internal/strategies/framework"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// BacktestEngine runs backtests on trading strategies using historical data
type BacktestEngine struct {
	logger *observability.Logger
	config BacktestConfig

	// Data and state
	historicalData map[string][]*HistoricalDataPoint
	currentTime    time.Time
	portfolio      *Portfolio

	// Results
	results *BacktestResults
}

// BacktestConfig contains backtesting configuration
type BacktestConfig struct {
	StartTime      time.Time       `json:"start_time"`
	EndTime        time.Time       `json:"end_time"`
	InitialBalance decimal.Decimal `json:"initial_balance"`
	CommissionRate decimal.Decimal `json:"commission_rate"`
	SlippageRate   decimal.Decimal `json:"slippage_rate"`
	DataResolution time.Duration   `json:"data_resolution"`
	Symbols        []string        `json:"symbols"`
	EnableLogging  bool            `json:"enable_logging"`
}

// HistoricalDataPoint represents a single data point in historical data
type HistoricalDataPoint struct {
	Timestamp time.Time       `json:"timestamp"`
	Symbol    string          `json:"symbol"`
	Open      decimal.Decimal `json:"open"`
	High      decimal.Decimal `json:"high"`
	Low       decimal.Decimal `json:"low"`
	Close     decimal.Decimal `json:"close"`
	Volume    decimal.Decimal `json:"volume"`
	BidPrice  decimal.Decimal `json:"bid_price"`
	AskPrice  decimal.Decimal `json:"ask_price"`
	BidVolume decimal.Decimal `json:"bid_volume"`
	AskVolume decimal.Decimal `json:"ask_volume"`
}

// Portfolio tracks the portfolio state during backtesting
type Portfolio struct {
	Cash      decimal.Decimal            `json:"cash"`
	Positions map[string]decimal.Decimal `json:"positions"`
	Value     decimal.Decimal            `json:"value"`
	PnL       decimal.Decimal            `json:"pnl"`
	Trades    []*Trade                   `json:"trades"`
}

// Trade represents a completed trade in backtesting
type Trade struct {
	ID         uuid.UUID       `json:"id"`
	Symbol     string          `json:"symbol"`
	Side       string          `json:"side"`
	Quantity   decimal.Decimal `json:"quantity"`
	Price      decimal.Decimal `json:"price"`
	Commission decimal.Decimal `json:"commission"`
	Slippage   decimal.Decimal `json:"slippage"`
	Timestamp  time.Time       `json:"timestamp"`
	StrategyID uuid.UUID       `json:"strategy_id"`
	SignalID   uuid.UUID       `json:"signal_id"`
}

// BacktestResults contains the results of a backtest
type BacktestResults struct {
	StartTime        time.Time       `json:"start_time"`
	EndTime          time.Time       `json:"end_time"`
	Duration         time.Duration   `json:"duration"`
	InitialBalance   decimal.Decimal `json:"initial_balance"`
	FinalBalance     decimal.Decimal `json:"final_balance"`
	TotalReturn      decimal.Decimal `json:"total_return"`
	TotalReturnPct   decimal.Decimal `json:"total_return_pct"`
	AnnualizedReturn decimal.Decimal `json:"annualized_return"`
	MaxDrawdown      decimal.Decimal `json:"max_drawdown"`
	SharpeRatio      decimal.Decimal `json:"sharpe_ratio"`
	SortinoRatio     decimal.Decimal `json:"sortino_ratio"`
	WinRate          decimal.Decimal `json:"win_rate"`
	ProfitFactor     decimal.Decimal `json:"profit_factor"`
	TotalTrades      int             `json:"total_trades"`
	WinningTrades    int             `json:"winning_trades"`
	LosingTrades     int             `json:"losing_trades"`
	AvgWin           decimal.Decimal `json:"avg_win"`
	AvgLoss          decimal.Decimal `json:"avg_loss"`
	MaxWin           decimal.Decimal `json:"max_win"`
	MaxLoss          decimal.Decimal `json:"max_loss"`
	TotalCommission  decimal.Decimal `json:"total_commission"`
	TotalSlippage    decimal.Decimal `json:"total_slippage"`
	Trades           []*Trade        `json:"trades"`
	EquityCurve      []*EquityPoint  `json:"equity_curve"`
}

// EquityPoint represents a point in the equity curve
type EquityPoint struct {
	Timestamp time.Time       `json:"timestamp"`
	Value     decimal.Decimal `json:"value"`
	Cash      decimal.Decimal `json:"cash"`
	PnL       decimal.Decimal `json:"pnl"`
}

// NewBacktestEngine creates a new backtesting engine
func NewBacktestEngine(logger *observability.Logger, config BacktestConfig) *BacktestEngine {
	return &BacktestEngine{
		logger:         logger,
		config:         config,
		historicalData: make(map[string][]*HistoricalDataPoint),
		portfolio: &Portfolio{
			Cash:      config.InitialBalance,
			Positions: make(map[string]decimal.Decimal),
			Value:     config.InitialBalance,
			PnL:       decimal.NewFromInt(0),
			Trades:    make([]*Trade, 0),
		},
		results: &BacktestResults{
			StartTime:      config.StartTime,
			EndTime:        config.EndTime,
			InitialBalance: config.InitialBalance,
			Trades:         make([]*Trade, 0),
			EquityCurve:    make([]*EquityPoint, 0),
		},
	}
}

// LoadHistoricalData loads historical data for backtesting
func (be *BacktestEngine) LoadHistoricalData(symbol string, data []*HistoricalDataPoint) error {
	if len(data) == 0 {
		return fmt.Errorf("no data provided for symbol %s", symbol)
	}

	// Sort data by timestamp
	sort.Slice(data, func(i, j int) bool {
		return data[i].Timestamp.Before(data[j].Timestamp)
	})

	// Validate data is within backtest period
	firstPoint := data[0]
	lastPoint := data[len(data)-1]

	if firstPoint.Timestamp.After(be.config.EndTime) || lastPoint.Timestamp.Before(be.config.StartTime) {
		return fmt.Errorf("data for symbol %s is outside backtest period", symbol)
	}

	be.historicalData[symbol] = data

	be.logger.Info(context.Background(), "Historical data loaded", map[string]interface{}{
		"symbol":      symbol,
		"data_points": len(data),
		"start_time":  firstPoint.Timestamp,
		"end_time":    lastPoint.Timestamp,
	})

	return nil
}

// RunBacktest runs a backtest for the given strategy
func (be *BacktestEngine) RunBacktest(ctx context.Context, strategy framework.Strategy) (*BacktestResults, error) {
	be.logger.Info(ctx, "Starting backtest", map[string]interface{}{
		"strategy":        strategy.GetName(),
		"start_time":      be.config.StartTime,
		"end_time":        be.config.EndTime,
		"symbols":         be.config.Symbols,
		"initial_balance": be.config.InitialBalance.String(),
	})

	// Initialize strategy
	strategyConfig := framework.StrategyConfig{
		ID:      strategy.GetID(),
		Name:    strategy.GetName(),
		Enabled: true,
		Symbols: be.config.Symbols,
	}

	if err := strategy.Initialize(ctx, strategyConfig); err != nil {
		return nil, fmt.Errorf("failed to initialize strategy: %w", err)
	}

	// Start strategy
	if err := strategy.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start strategy: %w", err)
	}
	defer strategy.Stop(ctx)

	// Run simulation
	be.currentTime = be.config.StartTime
	for be.currentTime.Before(be.config.EndTime) {
		if err := be.processTimeStep(ctx, strategy); err != nil {
			return nil, fmt.Errorf("error processing time step: %w", err)
		}

		be.currentTime = be.currentTime.Add(be.config.DataResolution)
	}

	// Calculate final results
	be.calculateResults()

	be.logger.Info(ctx, "Backtest completed", map[string]interface{}{
		"strategy":      strategy.GetName(),
		"total_trades":  be.results.TotalTrades,
		"final_balance": be.results.FinalBalance.String(),
		"total_return":  be.results.TotalReturnPct.String(),
		"max_drawdown":  be.results.MaxDrawdown.String(),
		"sharpe_ratio":  be.results.SharpeRatio.String(),
	})

	return be.results, nil
}

// processTimeStep processes a single time step in the backtest
func (be *BacktestEngine) processTimeStep(ctx context.Context, strategy framework.Strategy) error {
	// Get market data for current time
	marketData := be.getMarketDataAtTime(be.currentTime)

	// Process market data through strategy
	for _, data := range marketData {
		signals, err := strategy.OnMarketData(ctx, data)
		if err != nil {
			return fmt.Errorf("strategy error processing market data: %w", err)
		}

		// Execute signals
		for _, signal := range signals {
			if err := be.executeSignal(ctx, signal); err != nil {
				be.logger.Error(ctx, "Failed to execute signal", err, map[string]interface{}{
					"signal_id": signal.ID.String(),
					"symbol":    signal.Symbol,
					"side":      string(signal.Side),
				})
			}
		}
	}

	// Update portfolio value
	be.updatePortfolioValue()

	// Record equity curve point
	be.results.EquityCurve = append(be.results.EquityCurve, &EquityPoint{
		Timestamp: be.currentTime,
		Value:     be.portfolio.Value,
		Cash:      be.portfolio.Cash,
		PnL:       be.portfolio.PnL,
	})

	return nil
}

// getMarketDataAtTime gets market data for all symbols at the specified time
func (be *BacktestEngine) getMarketDataAtTime(timestamp time.Time) []*framework.MarketData {
	var marketData []*framework.MarketData

	for symbol, data := range be.historicalData {
		// Find the data point closest to the current time
		dataPoint := be.findDataPointAtTime(data, timestamp)
		if dataPoint == nil {
			continue
		}

		// Convert to market data format
		marketData = append(marketData, &framework.MarketData{
			Type:      framework.MarketDataTypeTicker,
			Symbol:    symbol,
			Exchange:  "backtest",
			Timestamp: timestamp,
			Data:      dataPoint,
		})
	}

	return marketData
}

// findDataPointAtTime finds the data point at or before the specified time
func (be *BacktestEngine) findDataPointAtTime(data []*HistoricalDataPoint, timestamp time.Time) *HistoricalDataPoint {
	// Binary search for efficiency
	left, right := 0, len(data)-1
	var result *HistoricalDataPoint

	for left <= right {
		mid := (left + right) / 2
		if data[mid].Timestamp.Equal(timestamp) {
			return data[mid]
		} else if data[mid].Timestamp.Before(timestamp) {
			result = data[mid]
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return result
}

// executeSignal executes a trading signal in the backtest
func (be *BacktestEngine) executeSignal(ctx context.Context, signal *framework.Signal) error {
	// Get current market data for the symbol
	dataPoint := be.getCurrentDataPoint(signal.Symbol)
	if dataPoint == nil {
		return fmt.Errorf("no market data available for symbol %s", signal.Symbol)
	}

	// Calculate execution price with slippage
	executionPrice := be.calculateExecutionPrice(signal, dataPoint)

	// Calculate commission
	notionalValue := signal.Quantity.Mul(executionPrice)
	commission := notionalValue.Mul(be.config.CommissionRate)

	// Check if we have enough cash/position for the trade
	if !be.canExecuteTrade(signal, executionPrice, commission) {
		return fmt.Errorf("insufficient funds or position for trade")
	}

	// Execute the trade
	trade := &Trade{
		ID:         uuid.New(),
		Symbol:     signal.Symbol,
		Side:       string(signal.Side),
		Quantity:   signal.Quantity,
		Price:      executionPrice,
		Commission: commission,
		Slippage:   executionPrice.Sub(signal.Price).Abs(),
		Timestamp:  be.currentTime,
		StrategyID: signal.StrategyID,
		SignalID:   signal.ID,
	}

	be.executeTrade(trade)

	if be.config.EnableLogging {
		be.logger.Info(ctx, "Trade executed", map[string]interface{}{
			"trade_id":   trade.ID.String(),
			"symbol":     trade.Symbol,
			"side":       trade.Side,
			"quantity":   trade.Quantity.String(),
			"price":      trade.Price.String(),
			"commission": trade.Commission.String(),
		})
	}

	return nil
}

// Private helper methods

// getCurrentDataPoint gets the current data point for a symbol
func (be *BacktestEngine) getCurrentDataPoint(symbol string) *HistoricalDataPoint {
	data, exists := be.historicalData[symbol]
	if !exists {
		return nil
	}

	return be.findDataPointAtTime(data, be.currentTime)
}

// calculateExecutionPrice calculates the execution price with slippage
func (be *BacktestEngine) calculateExecutionPrice(signal *framework.Signal, dataPoint *HistoricalDataPoint) decimal.Decimal {
	var basePrice decimal.Decimal

	// Use bid/ask prices if available, otherwise use close price
	if signal.Side == "BUY" {
		if !dataPoint.AskPrice.IsZero() {
			basePrice = dataPoint.AskPrice
		} else {
			basePrice = dataPoint.Close
		}
	} else {
		if !dataPoint.BidPrice.IsZero() {
			basePrice = dataPoint.BidPrice
		} else {
			basePrice = dataPoint.Close
		}
	}

	// Apply slippage
	slippage := basePrice.Mul(be.config.SlippageRate)
	if signal.Side == "BUY" {
		return basePrice.Add(slippage)
	} else {
		return basePrice.Sub(slippage)
	}
}

// canExecuteTrade checks if a trade can be executed
func (be *BacktestEngine) canExecuteTrade(signal *framework.Signal, price, commission decimal.Decimal) bool {
	if signal.Side == "BUY" {
		// Check if we have enough cash
		totalCost := signal.Quantity.Mul(price).Add(commission)
		return be.portfolio.Cash.GreaterThanOrEqual(totalCost)
	} else {
		// Check if we have enough position
		currentPosition := be.portfolio.Positions[signal.Symbol]
		return currentPosition.GreaterThanOrEqual(signal.Quantity)
	}
}

// executeTrade executes a trade and updates the portfolio
func (be *BacktestEngine) executeTrade(trade *Trade) {
	if trade.Side == "BUY" {
		// Deduct cash and add position
		totalCost := trade.Quantity.Mul(trade.Price).Add(trade.Commission)
		be.portfolio.Cash = be.portfolio.Cash.Sub(totalCost)

		currentPosition := be.portfolio.Positions[trade.Symbol]
		be.portfolio.Positions[trade.Symbol] = currentPosition.Add(trade.Quantity)
	} else {
		// Add cash and reduce position
		totalReceived := trade.Quantity.Mul(trade.Price).Sub(trade.Commission)
		be.portfolio.Cash = be.portfolio.Cash.Add(totalReceived)

		currentPosition := be.portfolio.Positions[trade.Symbol]
		be.portfolio.Positions[trade.Symbol] = currentPosition.Sub(trade.Quantity)
	}

	// Add trade to portfolio and results
	be.portfolio.Trades = append(be.portfolio.Trades, trade)
	be.results.Trades = append(be.results.Trades, trade)
}

// updatePortfolioValue updates the total portfolio value
func (be *BacktestEngine) updatePortfolioValue() {
	totalValue := be.portfolio.Cash

	// Add value of all positions
	for symbol, quantity := range be.portfolio.Positions {
		if quantity.IsZero() {
			continue
		}

		dataPoint := be.getCurrentDataPoint(symbol)
		if dataPoint != nil {
			positionValue := quantity.Mul(dataPoint.Close)
			totalValue = totalValue.Add(positionValue)
		}
	}

	be.portfolio.Value = totalValue
	be.portfolio.PnL = totalValue.Sub(be.config.InitialBalance)
}

// calculateResults calculates final backtest results
func (be *BacktestEngine) calculateResults() {
	// Calculate total return percentage
	if be.config.InitialBalance.IsPositive() {
		be.results.TotalReturnPct = be.portfolio.PnL.Div(be.config.InitialBalance).Mul(decimal.NewFromInt(100))
	}

	// Set final balance
	be.results.FinalBalance = be.portfolio.Value

	// Calculate win rate
	if be.results.TotalTrades > 0 {
		be.results.WinRate = decimal.NewFromInt(int64(be.results.WinningTrades)).Div(decimal.NewFromInt(int64(be.results.TotalTrades))).Mul(decimal.NewFromInt(100))
	}

	// Calculate basic trade statistics
	if be.results.TotalTrades > 0 {
		// Calculate average win and loss
		totalWins := decimal.Zero
		totalLosses := decimal.Zero

		// For now, use simplified calculation based on total return
		if be.results.WinningTrades > 0 {
			be.results.AvgWin = be.results.TotalReturn.Div(decimal.NewFromInt(int64(be.results.WinningTrades)))
		}

		if be.results.LosingTrades > 0 {
			be.results.AvgLoss = be.results.TotalReturn.Div(decimal.NewFromInt(int64(be.results.LosingTrades))).Neg()
		}

		// Calculate profit factor
		if be.results.AvgLoss.IsPositive() && be.results.WinningTrades > 0 && be.results.LosingTrades > 0 {
			totalWins = be.results.AvgWin.Mul(decimal.NewFromInt(int64(be.results.WinningTrades)))
			totalLosses = be.results.AvgLoss.Mul(decimal.NewFromInt(int64(be.results.LosingTrades)))
			if totalLosses.IsPositive() {
				be.results.ProfitFactor = totalWins.Div(totalLosses)
			}
		}
	}

	// Calculate simplified Sharpe ratio based on total return and volatility estimate
	if be.results.TotalReturnPct.IsPositive() && be.results.TotalTrades > 1 {
		// Simplified volatility estimate (this would be more accurate with actual price data)
		volatilityEstimate := be.results.MaxDrawdown.Mul(decimal.NewFromFloat(2.0))
		if volatilityEstimate.IsPositive() {
			be.results.SharpeRatio = be.results.TotalReturnPct.Div(volatilityEstimate)
		}
	}
}
