package hft

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// PortfolioManager manages real-time portfolio positions and P&L
type PortfolioManager struct {
	logger     *observability.Logger
	config     HFTConfig
	positions  map[string]*Position
	portfolio  *Portfolio
	priceCache map[string]decimal.Decimal

	// Performance tracking
	totalTrades   int64
	totalPnL      decimal.Decimal
	unrealizedPnL decimal.Decimal
	realizedPnL   decimal.Decimal

	// State management
	isRunning int32
	stopChan  chan struct{}
	wg        sync.WaitGroup
	mu        sync.RWMutex

	// Real-time updates
	updateChan chan PositionUpdate
}

// Position represents a trading position
type Position struct {
	ID            uuid.UUID              `json:"id"`
	Symbol        string                 `json:"symbol"`
	Size          decimal.Decimal        `json:"size"` // Positive for long, negative for short
	AvgPrice      decimal.Decimal        `json:"avg_price"`
	CurrentPrice  decimal.Decimal        `json:"current_price"`
	UnrealizedPnL decimal.Decimal        `json:"unrealized_pnl"`
	RealizedPnL   decimal.Decimal        `json:"realized_pnl"`
	Commission    decimal.Decimal        `json:"commission"`
	OpenTime      time.Time              `json:"open_time"`
	UpdateTime    time.Time              `json:"update_time"`
	Exchange      string                 `json:"exchange"`
	StrategyID    string                 `json:"strategy_id"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// Portfolio represents the overall portfolio
type Portfolio struct {
	ID            uuid.UUID                  `json:"id"`
	TotalValue    decimal.Decimal            `json:"total_value"`
	CashBalance   decimal.Decimal            `json:"cash_balance"`
	TotalPnL      decimal.Decimal            `json:"total_pnl"`
	UnrealizedPnL decimal.Decimal            `json:"unrealized_pnl"`
	RealizedPnL   decimal.Decimal            `json:"realized_pnl"`
	DayPnL        decimal.Decimal            `json:"day_pnl"`
	MaxDrawdown   decimal.Decimal            `json:"max_drawdown"`
	HighWaterMark decimal.Decimal            `json:"high_water_mark"`
	Positions     map[string]*Position       `json:"positions"`
	Allocations   map[string]decimal.Decimal `json:"allocations"`
	RiskMetrics   *PortfolioRiskMetrics      `json:"risk_metrics"`
	LastUpdate    time.Time                  `json:"last_update"`
}

// PortfolioRiskMetrics contains portfolio risk metrics
type PortfolioRiskMetrics struct {
	VaR95             decimal.Decimal    `json:"var_95"`             // Value at Risk 95%
	VaR99             decimal.Decimal    `json:"var_99"`             // Value at Risk 99%
	ExpectedShortfall decimal.Decimal    `json:"expected_shortfall"` // Conditional VaR
	SharpeRatio       float64            `json:"sharpe_ratio"`
	SortinoRatio      float64            `json:"sortino_ratio"`
	MaxDrawdown       decimal.Decimal    `json:"max_drawdown"`
	Beta              float64            `json:"beta"`
	Alpha             float64            `json:"alpha"`
	Volatility        float64            `json:"volatility"`
	Correlation       map[string]float64 `json:"correlation"`
}

// Trade represents a completed trade
type Trade struct {
	ID         uuid.UUID       `json:"id"`
	Symbol     string          `json:"symbol"`
	Side       OrderSide       `json:"side"`
	Quantity   decimal.Decimal `json:"quantity"`
	Price      decimal.Decimal `json:"price"`
	Commission decimal.Decimal `json:"commission"`
	PnL        decimal.Decimal `json:"pnl"`
	Timestamp  time.Time       `json:"timestamp"`
	Exchange   string          `json:"exchange"`
	StrategyID string          `json:"strategy_id"`
	OrderID    uuid.UUID       `json:"order_id"`
}

// NewPortfolioManager creates a new portfolio manager
func NewPortfolioManager(logger *observability.Logger, config HFTConfig) *PortfolioManager {
	portfolio := &Portfolio{
		ID:          uuid.New(),
		TotalValue:  decimal.NewFromFloat(100000), // Default $100k
		CashBalance: decimal.NewFromFloat(100000),
		Positions:   make(map[string]*Position),
		Allocations: make(map[string]decimal.Decimal),
		RiskMetrics: &PortfolioRiskMetrics{},
		LastUpdate:  time.Now(),
	}

	return &PortfolioManager{
		logger:     logger,
		config:     config,
		positions:  make(map[string]*Position),
		portfolio:  portfolio,
		priceCache: make(map[string]decimal.Decimal),
		stopChan:   make(chan struct{}),
		updateChan: make(chan PositionUpdate, 1000),
	}
}

// Start begins the portfolio manager
func (pm *PortfolioManager) Start(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&pm.isRunning, 0, 1) {
		return fmt.Errorf("portfolio manager is already running")
	}

	pm.logger.Info(ctx, "Starting portfolio manager", map[string]interface{}{
		"portfolio_id":  pm.portfolio.ID.String(),
		"initial_value": pm.portfolio.TotalValue.String(),
		"cash_balance":  pm.portfolio.CashBalance.String(),
	})

	// Start processing goroutines
	pm.wg.Add(2)
	go pm.processUpdates(ctx)
	go pm.calculateMetrics(ctx)

	return nil
}

// Stop gracefully shuts down the portfolio manager
func (pm *PortfolioManager) Stop(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&pm.isRunning, 1, 0) {
		return fmt.Errorf("portfolio manager is not running")
	}

	pm.logger.Info(ctx, "Stopping portfolio manager", nil)

	close(pm.stopChan)
	pm.wg.Wait()

	pm.logger.Info(ctx, "Portfolio manager stopped", map[string]interface{}{
		"total_trades":   atomic.LoadInt64(&pm.totalTrades),
		"total_pnl":      pm.totalPnL.String(),
		"unrealized_pnl": pm.unrealizedPnL.String(),
		"realized_pnl":   pm.realizedPnL.String(),
	})

	return nil
}

// UpdatePrice updates the current price for a symbol
func (pm *PortfolioManager) UpdatePrice(symbol string, price decimal.Decimal) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.priceCache[symbol] = price

	// Update position if exists
	if position, exists := pm.positions[symbol]; exists {
		position.CurrentPrice = price
		position.UpdateTime = time.Now()

		// Calculate unrealized P&L
		if !position.Size.IsZero() {
			priceDiff := price.Sub(position.AvgPrice)
			position.UnrealizedPnL = position.Size.Mul(priceDiff)
		}
	}
}

// HandleOrderUpdate processes order updates and updates positions
func (pm *PortfolioManager) HandleOrderUpdate(update OrderUpdate) {
	if update.Status != OrderStatusFilled && update.Status != OrderStatusPartialFill {
		return
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Get or create position
	position, exists := pm.positions[update.Symbol]
	if !exists {
		position = &Position{
			ID:         uuid.New(),
			Symbol:     update.Symbol,
			Size:       decimal.Zero,
			AvgPrice:   decimal.Zero,
			OpenTime:   update.Timestamp,
			UpdateTime: update.Timestamp,
		}
		pm.positions[update.Symbol] = position
	}

	// Calculate trade details
	tradeSize := update.FilledQty
	if update.Side == OrderSideSell {
		tradeSize = tradeSize.Neg()
	}

	// Update position
	pm.updatePosition(position, tradeSize, update.AvgFillPrice, update.Timestamp)

	// Update portfolio
	pm.updatePortfolio()

	atomic.AddInt64(&pm.totalTrades, 1)

	pm.logger.Info(context.Background(), "Position updated", map[string]interface{}{
		"symbol":         update.Symbol,
		"trade_size":     tradeSize.String(),
		"trade_price":    update.AvgFillPrice.String(),
		"position_size":  position.Size.String(),
		"avg_price":      position.AvgPrice.String(),
		"unrealized_pnl": position.UnrealizedPnL.String(),
	})
}

// HandlePositionUpdate processes external position updates
func (pm *PortfolioManager) HandlePositionUpdate(update PositionUpdate) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	position, exists := pm.positions[update.Symbol]
	if !exists {
		position = &Position{
			ID:       uuid.New(),
			Symbol:   update.Symbol,
			OpenTime: update.Timestamp,
		}
		pm.positions[update.Symbol] = position
	}

	position.Size = update.Size
	position.AvgPrice = update.AvgPrice
	position.UnrealizedPnL = update.UnrealizedPnL
	position.RealizedPnL = update.RealizedPnL
	position.UpdateTime = update.Timestamp

	pm.updatePortfolio()
}

// GetPosition returns a position by symbol
func (pm *PortfolioManager) GetPosition(symbol string) (*Position, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	position, exists := pm.positions[symbol]
	if !exists {
		return nil, fmt.Errorf("position not found for symbol: %s", symbol)
	}

	return position, nil
}

// GetAllPositions returns all positions
func (pm *PortfolioManager) GetAllPositions() map[string]*Position {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	positions := make(map[string]*Position)
	for symbol, position := range pm.positions {
		positions[symbol] = position
	}

	return positions
}

// GetPortfolio returns the current portfolio
func (pm *PortfolioManager) GetPortfolio() *Portfolio {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	return pm.portfolio
}

// GetPositionCount returns the number of open positions
func (pm *PortfolioManager) GetPositionCount() int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	count := 0
	for _, position := range pm.positions {
		if !position.Size.IsZero() {
			count++
		}
	}

	return count
}

// GetTotalPnL returns the total P&L
func (pm *PortfolioManager) GetTotalPnL() decimal.Decimal {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	return pm.portfolio.TotalPnL
}

// updatePosition updates a position with a new trade
func (pm *PortfolioManager) updatePosition(position *Position, tradeSize, tradePrice decimal.Decimal, timestamp time.Time) {
	if position.Size.IsZero() {
		// Opening new position
		position.Size = tradeSize
		position.AvgPrice = tradePrice
		position.OpenTime = timestamp
	} else if position.Size.Sign() == tradeSize.Sign() {
		// Adding to existing position
		totalValue := position.Size.Mul(position.AvgPrice).Add(tradeSize.Mul(tradePrice))
		position.Size = position.Size.Add(tradeSize)
		if !position.Size.IsZero() {
			position.AvgPrice = totalValue.Div(position.Size)
		}
	} else {
		// Reducing or closing position
		if tradeSize.Abs().GreaterThan(position.Size.Abs()) {
			// Reversing position
			closingSize := position.Size.Neg()
			remainingSize := tradeSize.Add(closingSize)

			// Realize P&L for closed portion
			pnl := closingSize.Mul(tradePrice.Sub(position.AvgPrice))
			position.RealizedPnL = position.RealizedPnL.Add(pnl)
			pm.realizedPnL = pm.realizedPnL.Add(pnl)

			// Set new position
			position.Size = remainingSize
			position.AvgPrice = tradePrice
		} else {
			// Partial or full close
			pnl := tradeSize.Neg().Mul(tradePrice.Sub(position.AvgPrice))
			position.RealizedPnL = position.RealizedPnL.Add(pnl)
			pm.realizedPnL = pm.realizedPnL.Add(pnl)

			position.Size = position.Size.Add(tradeSize)
		}
	}

	position.UpdateTime = timestamp

	// Update unrealized P&L if we have current price
	if currentPrice, exists := pm.priceCache[position.Symbol]; exists {
		position.CurrentPrice = currentPrice
		if !position.Size.IsZero() {
			priceDiff := currentPrice.Sub(position.AvgPrice)
			position.UnrealizedPnL = position.Size.Mul(priceDiff)
		}
	}
}

// updatePortfolio recalculates portfolio metrics
func (pm *PortfolioManager) updatePortfolio() {
	totalUnrealizedPnL := decimal.Zero
	totalRealizedPnL := decimal.Zero
	totalValue := pm.portfolio.CashBalance

	for _, position := range pm.positions {
		totalUnrealizedPnL = totalUnrealizedPnL.Add(position.UnrealizedPnL)
		totalRealizedPnL = totalRealizedPnL.Add(position.RealizedPnL)

		// Add position value to total
		if !position.Size.IsZero() && !position.CurrentPrice.IsZero() {
			positionValue := position.Size.Abs().Mul(position.CurrentPrice)
			totalValue = totalValue.Add(positionValue)
		}
	}

	pm.portfolio.UnrealizedPnL = totalUnrealizedPnL
	pm.portfolio.RealizedPnL = totalRealizedPnL
	pm.portfolio.TotalPnL = totalUnrealizedPnL.Add(totalRealizedPnL)
	pm.portfolio.TotalValue = totalValue.Add(pm.portfolio.TotalPnL)
	pm.portfolio.LastUpdate = time.Now()

	// Update high water mark and drawdown
	if pm.portfolio.TotalValue.GreaterThan(pm.portfolio.HighWaterMark) {
		pm.portfolio.HighWaterMark = pm.portfolio.TotalValue
	}

	drawdown := pm.portfolio.HighWaterMark.Sub(pm.portfolio.TotalValue)
	if drawdown.GreaterThan(pm.portfolio.MaxDrawdown) {
		pm.portfolio.MaxDrawdown = drawdown
	}

	pm.totalPnL = pm.portfolio.TotalPnL
	pm.unrealizedPnL = totalUnrealizedPnL
}

// processUpdates processes position updates
func (pm *PortfolioManager) processUpdates(ctx context.Context) {
	defer pm.wg.Done()

	for {
		select {
		case <-pm.stopChan:
			return
		case update := <-pm.updateChan:
			pm.HandlePositionUpdate(update)
		}
	}
}

// calculateMetrics periodically calculates portfolio metrics
func (pm *PortfolioManager) calculateMetrics(ctx context.Context) {
	defer pm.wg.Done()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-pm.stopChan:
			return
		case <-ticker.C:
			pm.calculateRiskMetrics()
		}
	}
}

// calculateRiskMetrics calculates portfolio risk metrics
func (pm *PortfolioManager) calculateRiskMetrics() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// This is a simplified implementation
	// In production, you would use historical returns and more sophisticated calculations

	if pm.portfolio.RiskMetrics == nil {
		pm.portfolio.RiskMetrics = &PortfolioRiskMetrics{}
	}

	// Calculate basic metrics
	totalValue := pm.portfolio.TotalValue
	if totalValue.GreaterThan(decimal.Zero) {
		drawdownPercent := pm.portfolio.MaxDrawdown.Div(pm.portfolio.HighWaterMark).Mul(decimal.NewFromInt(100))
		pm.portfolio.RiskMetrics.MaxDrawdown = drawdownPercent
	}

	// TODO: Implement more sophisticated risk calculations
	// - VaR calculations using historical simulation or Monte Carlo
	// - Sharpe and Sortino ratios using return history
	// - Beta and Alpha calculations against benchmark
	// - Correlation analysis
}

// GetMetrics returns portfolio manager metrics
func (pm *PortfolioManager) GetMetrics() PortfolioMetrics {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	return PortfolioMetrics{
		TotalTrades:   atomic.LoadInt64(&pm.totalTrades),
		TotalPnL:      pm.totalPnL,
		UnrealizedPnL: pm.unrealizedPnL,
		RealizedPnL:   pm.realizedPnL,
		TotalValue:    pm.portfolio.TotalValue,
		CashBalance:   pm.portfolio.CashBalance,
		OpenPositions: pm.GetPositionCount(),
		MaxDrawdown:   pm.portfolio.MaxDrawdown,
		HighWaterMark: pm.portfolio.HighWaterMark,
		IsRunning:     atomic.LoadInt32(&pm.isRunning) == 1,
	}
}

// PortfolioMetrics contains portfolio performance metrics
type PortfolioMetrics struct {
	TotalTrades   int64           `json:"total_trades"`
	TotalPnL      decimal.Decimal `json:"total_pnl"`
	UnrealizedPnL decimal.Decimal `json:"unrealized_pnl"`
	RealizedPnL   decimal.Decimal `json:"realized_pnl"`
	TotalValue    decimal.Decimal `json:"total_value"`
	CashBalance   decimal.Decimal `json:"cash_balance"`
	OpenPositions int             `json:"open_positions"`
	MaxDrawdown   decimal.Decimal `json:"max_drawdown"`
	HighWaterMark decimal.Decimal `json:"high_water_mark"`
	IsRunning     bool            `json:"is_running"`
}
