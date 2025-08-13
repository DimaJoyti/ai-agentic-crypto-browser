package paper_trading

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/ai-agentic-browser/internal/exchanges"
	"github.com/ai-agentic-browser/internal/exchanges/common"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// PaperTradingService provides paper trading functionality for strategy validation
type PaperTradingService struct {
	logger          *observability.Logger
	exchangeManager *exchanges.Manager
	portfolios      map[uuid.UUID]*PaperPortfolio
	orders          map[uuid.UUID]*PaperOrder
	config          PaperTradingConfig

	// State management
	isRunning bool
	stopChan  chan struct{}
	wg        sync.WaitGroup
	mu        sync.RWMutex
}

// PaperTradingConfig contains paper trading configuration
type PaperTradingConfig struct {
	InitialBalance decimal.Decimal `json:"initial_balance"`
	CommissionRate decimal.Decimal `json:"commission_rate"`
	SlippageRate   decimal.Decimal `json:"slippage_rate"`
	EnableRealism  bool            `json:"enable_realism"`  // Simulate realistic execution
	EnableLatency  bool            `json:"enable_latency"`  // Simulate network latency
	MaxLatency     time.Duration   `json:"max_latency"`     // Maximum simulated latency
	UpdateInterval time.Duration   `json:"update_interval"` // Portfolio update frequency
	EnableLogging  bool            `json:"enable_logging"`  // Enable detailed logging
}

// PaperPortfolio represents a paper trading portfolio
type PaperPortfolio struct {
	ID         uuid.UUID                 `json:"id"`
	StrategyID uuid.UUID                 `json:"strategy_id"`
	Name       string                    `json:"name"`
	Cash       decimal.Decimal           `json:"cash"`
	Positions  map[string]*PaperPosition `json:"positions"`
	Orders     map[uuid.UUID]*PaperOrder `json:"orders"`
	Trades     []*PaperTrade             `json:"trades"`
	Value      decimal.Decimal           `json:"value"`
	PnL        decimal.Decimal           `json:"pnl"`
	Metrics    *PaperTradingMetrics      `json:"metrics"`
	CreatedAt  time.Time                 `json:"created_at"`
	UpdatedAt  time.Time                 `json:"updated_at"`
}

// PaperPosition represents a position in paper trading
type PaperPosition struct {
	Symbol        string          `json:"symbol"`
	Quantity      decimal.Decimal `json:"quantity"`
	AvgPrice      decimal.Decimal `json:"avg_price"`
	MarketValue   decimal.Decimal `json:"market_value"`
	UnrealizedPnL decimal.Decimal `json:"unrealized_pnl"`
	RealizedPnL   decimal.Decimal `json:"realized_pnl"`
	LastPrice     decimal.Decimal `json:"last_price"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// PaperOrder represents an order in paper trading
type PaperOrder struct {
	ID           uuid.UUID          `json:"id"`
	PortfolioID  uuid.UUID          `json:"portfolio_id"`
	StrategyID   uuid.UUID          `json:"strategy_id"`
	Symbol       string             `json:"symbol"`
	Side         common.OrderSide   `json:"side"`
	Type         common.OrderType   `json:"type"`
	Quantity     decimal.Decimal    `json:"quantity"`
	Price        decimal.Decimal    `json:"price"`
	StopPrice    decimal.Decimal    `json:"stop_price,omitempty"`
	TimeInForce  common.TimeInForce `json:"time_in_force"`
	Status       common.OrderStatus `json:"status"`
	FilledQty    decimal.Decimal    `json:"filled_qty"`
	RemainingQty decimal.Decimal    `json:"remaining_qty"`
	AvgFillPrice decimal.Decimal    `json:"avg_fill_price"`
	Commission   decimal.Decimal    `json:"commission"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
	ExecutedAt   time.Time          `json:"executed_at,omitempty"`
	Metadata     map[string]any     `json:"metadata"`
}

// PaperTrade represents a completed trade in paper trading
type PaperTrade struct {
	ID          uuid.UUID        `json:"id"`
	PortfolioID uuid.UUID        `json:"portfolio_id"`
	OrderID     uuid.UUID        `json:"order_id"`
	Symbol      string           `json:"symbol"`
	Side        common.OrderSide `json:"side"`
	Quantity    decimal.Decimal  `json:"quantity"`
	Price       decimal.Decimal  `json:"price"`
	Commission  decimal.Decimal  `json:"commission"`
	PnL         decimal.Decimal  `json:"pnl"`
	Timestamp   time.Time        `json:"timestamp"`
}

// PaperTradingMetrics contains paper trading performance metrics
type PaperTradingMetrics struct {
	TotalTrades     int             `json:"total_trades"`
	WinningTrades   int             `json:"winning_trades"`
	LosingTrades    int             `json:"losing_trades"`
	WinRate         decimal.Decimal `json:"win_rate"`
	TotalPnL        decimal.Decimal `json:"total_pnl"`
	RealizedPnL     decimal.Decimal `json:"realized_pnl"`
	UnrealizedPnL   decimal.Decimal `json:"unrealized_pnl"`
	TotalCommission decimal.Decimal `json:"total_commission"`
	MaxDrawdown     decimal.Decimal `json:"max_drawdown"`
	SharpeRatio     decimal.Decimal `json:"sharpe_ratio"`
	StartTime       time.Time       `json:"start_time"`
	LastTradeTime   time.Time       `json:"last_trade_time"`
}

// NewPaperTradingService creates a new paper trading service
func NewPaperTradingService(
	logger *observability.Logger,
	exchangeManager *exchanges.Manager,
	config PaperTradingConfig,
) *PaperTradingService {
	if config.UpdateInterval == 0 {
		config.UpdateInterval = 1 * time.Second
	}
	if config.MaxLatency == 0 {
		config.MaxLatency = 100 * time.Millisecond
	}

	return &PaperTradingService{
		logger:          logger,
		exchangeManager: exchangeManager,
		portfolios:      make(map[uuid.UUID]*PaperPortfolio),
		orders:          make(map[uuid.UUID]*PaperOrder),
		config:          config,
		stopChan:        make(chan struct{}),
	}
}

// Start starts the paper trading service
func (pts *PaperTradingService) Start(ctx context.Context) error {
	pts.mu.Lock()
	defer pts.mu.Unlock()

	if pts.isRunning {
		return fmt.Errorf("paper trading service is already running")
	}

	pts.logger.Info(ctx, "Starting paper trading service", map[string]interface{}{
		"initial_balance": pts.config.InitialBalance.String(),
		"commission_rate": pts.config.CommissionRate.String(),
		"enable_realism":  pts.config.EnableRealism,
		"enable_latency":  pts.config.EnableLatency,
	})

	pts.isRunning = true

	// Start portfolio update loop
	pts.wg.Add(1)
	go pts.updatePortfolios(ctx)

	// Start order processing loop
	pts.wg.Add(1)
	go pts.processOrders(ctx)

	pts.logger.Info(ctx, "Paper trading service started", nil)

	return nil
}

// Stop stops the paper trading service
func (pts *PaperTradingService) Stop(ctx context.Context) error {
	pts.mu.Lock()
	defer pts.mu.Unlock()

	if !pts.isRunning {
		return fmt.Errorf("paper trading service is not running")
	}

	pts.logger.Info(ctx, "Stopping paper trading service", nil)

	close(pts.stopChan)
	pts.wg.Wait()

	pts.isRunning = false

	pts.logger.Info(ctx, "Paper trading service stopped", nil)

	return nil
}

// CreatePortfolio creates a new paper trading portfolio
func (pts *PaperTradingService) CreatePortfolio(ctx context.Context, strategyID uuid.UUID, name string) (*PaperPortfolio, error) {
	pts.mu.Lock()
	defer pts.mu.Unlock()

	portfolio := &PaperPortfolio{
		ID:         uuid.New(),
		StrategyID: strategyID,
		Name:       name,
		Cash:       pts.config.InitialBalance,
		Positions:  make(map[string]*PaperPosition),
		Orders:     make(map[uuid.UUID]*PaperOrder),
		Trades:     make([]*PaperTrade, 0),
		Value:      pts.config.InitialBalance,
		PnL:        decimal.NewFromInt(0),
		Metrics: &PaperTradingMetrics{
			StartTime: time.Now(),
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	pts.portfolios[portfolio.ID] = portfolio

	pts.logger.Info(ctx, "Paper trading portfolio created", map[string]interface{}{
		"portfolio_id":    portfolio.ID.String(),
		"strategy_id":     strategyID.String(),
		"name":            name,
		"initial_balance": pts.config.InitialBalance.String(),
	})

	return portfolio, nil
}

// SubmitOrder submits an order for paper trading
func (pts *PaperTradingService) SubmitOrder(ctx context.Context, portfolioID uuid.UUID, orderReq *common.OrderRequest) (*PaperOrder, error) {
	pts.mu.Lock()
	defer pts.mu.Unlock()

	portfolio, exists := pts.portfolios[portfolioID]
	if !exists {
		return nil, fmt.Errorf("portfolio not found: %s", portfolioID.String())
	}

	// Create paper order
	order := &PaperOrder{
		ID:           uuid.New(),
		PortfolioID:  portfolioID,
		Symbol:       orderReq.Symbol,
		Side:         orderReq.Side,
		Type:         orderReq.Type,
		Quantity:     orderReq.Quantity,
		Price:        orderReq.Price,
		StopPrice:    orderReq.StopPrice,
		TimeInForce:  orderReq.TimeInForce,
		Status:       common.OrderStatusNew,
		FilledQty:    decimal.NewFromInt(0),
		RemainingQty: orderReq.Quantity,
		AvgFillPrice: decimal.NewFromInt(0),
		Commission:   decimal.NewFromInt(0),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Metadata:     orderReq.Metadata,
	}

	// Extract strategy ID from metadata if available
	if orderReq.Metadata != nil {
		if strategyID, ok := orderReq.Metadata["strategy_id"].(string); ok {
			if id, err := uuid.Parse(strategyID); err == nil {
				order.StrategyID = id
			}
		}
	}

	// Validate order
	if err := pts.validateOrder(portfolio, order); err != nil {
		return nil, fmt.Errorf("order validation failed: %w", err)
	}

	// Store order
	pts.orders[order.ID] = order
	portfolio.Orders[order.ID] = order

	// Process order immediately (latency simulation is handled in the processing loop)
	pts.logger.Info(ctx, "Order placed successfully", map[string]interface{}{
		"order_id": order.ID,
		"symbol":   order.Symbol,
		"side":     order.Side,
		"quantity": order.Quantity.String(),
		"price":    order.Price.String(),
	})

	if pts.config.EnableLogging {
		pts.logger.Info(ctx, "Paper order submitted", map[string]interface{}{
			"order_id":     order.ID.String(),
			"portfolio_id": portfolioID.String(),
			"symbol":       order.Symbol,
			"side":         string(order.Side),
			"quantity":     order.Quantity.String(),
			"price":        order.Price.String(),
		})
	}

	return order, nil
}

// GetPortfolio returns a paper trading portfolio
func (pts *PaperTradingService) GetPortfolio(portfolioID uuid.UUID) (*PaperPortfolio, error) {
	pts.mu.RLock()
	defer pts.mu.RUnlock()

	portfolio, exists := pts.portfolios[portfolioID]
	if !exists {
		return nil, fmt.Errorf("portfolio not found: %s", portfolioID.String())
	}

	return portfolio, nil
}

// GetOrder returns a paper trading order
func (pts *PaperTradingService) GetOrder(orderID uuid.UUID) (*PaperOrder, error) {
	pts.mu.RLock()
	defer pts.mu.RUnlock()

	order, exists := pts.orders[orderID]
	if !exists {
		return nil, fmt.Errorf("order not found: %s", orderID.String())
	}

	return order, nil
}

// CancelOrder cancels a paper trading order
func (pts *PaperTradingService) CancelOrder(ctx context.Context, orderID uuid.UUID) error {
	pts.mu.Lock()
	defer pts.mu.Unlock()

	order, exists := pts.orders[orderID]
	if !exists {
		return fmt.Errorf("order not found: %s", orderID.String())
	}

	if order.Status != common.OrderStatusNew && order.Status != common.OrderStatusPartiallyFilled {
		return fmt.Errorf("order cannot be cancelled: %s", order.Status)
	}

	order.Status = common.OrderStatusCanceled
	order.UpdatedAt = time.Now()

	if pts.config.EnableLogging {
		pts.logger.Info(ctx, "Paper order cancelled", map[string]interface{}{
			"order_id": orderID.String(),
		})
	}

	return nil
}

// Private methods

// validateOrder validates a paper trading order
func (pts *PaperTradingService) validateOrder(portfolio *PaperPortfolio, order *PaperOrder) error {
	if order.Quantity.IsZero() || order.Quantity.IsNegative() {
		return fmt.Errorf("quantity must be positive")
	}

	if order.Side == common.OrderSideBuy {
		// Check if we have enough cash
		if order.Type == common.OrderTypeMarket {
			// For market orders, we need to estimate the cost
			// This is a simplified check
			estimatedCost := order.Quantity.Mul(order.Price)
			if portfolio.Cash.LessThan(estimatedCost) {
				return fmt.Errorf("insufficient cash for buy order")
			}
		} else {
			// For limit orders, check exact cost
			totalCost := order.Quantity.Mul(order.Price)
			if portfolio.Cash.LessThan(totalCost) {
				return fmt.Errorf("insufficient cash for buy order")
			}
		}
	} else {
		// Check if we have enough position to sell
		position, exists := portfolio.Positions[order.Symbol]
		if !exists || position.Quantity.LessThan(order.Quantity) {
			return fmt.Errorf("insufficient position for sell order")
		}
	}

	return nil
}

// updatePortfolios runs the portfolio update loop
func (pts *PaperTradingService) updatePortfolios(ctx context.Context) {
	defer pts.wg.Done()

	ticker := time.NewTicker(5 * time.Second) // Update every 5 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-pts.stopChan:
			return
		case <-ticker.C:
			pts.mu.RLock()
			for userID, portfolio := range pts.portfolios {
				// Update portfolio values based on current market prices
				totalValue := portfolio.Cash

				for _, position := range portfolio.Positions {
					// Get current market price (placeholder - would use real market data)
					currentPrice := position.AvgPrice.Mul(decimal.NewFromFloat(1.0 + (rand.Float64()-0.5)*0.02)) // ±1% random movement
					positionValue := position.Quantity.Mul(currentPrice)
					totalValue = totalValue.Add(positionValue)

					// Update position last price
					position.LastPrice = currentPrice
					position.MarketValue = positionValue
					position.UnrealizedPnL = position.Quantity.Mul(currentPrice.Sub(position.AvgPrice))
					position.UpdatedAt = time.Now()
				}

				portfolio.Value = totalValue
				portfolio.PnL = totalValue.Sub(pts.config.InitialBalance)
				portfolio.UpdatedAt = time.Now()

				pts.logger.Debug(ctx, "Updated portfolio", map[string]interface{}{
					"user_id":     userID,
					"total_value": totalValue.String(),
					"total_pnl":   portfolio.PnL.String(),
				})
			}
			pts.mu.RUnlock()
		}
	}
}

// processOrders runs the order processing loop
func (pts *PaperTradingService) processOrders(ctx context.Context) {
	defer pts.wg.Done()

	ticker := time.NewTicker(1 * time.Second) // Process orders every second
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-pts.stopChan:
			return
		case <-ticker.C:
			pts.mu.Lock()
			for orderID, order := range pts.orders {
				if order.Status == "pending" {
					// Simulate order latency
					if time.Since(order.CreatedAt) >= pts.simulateOrderLatency() {
						// Process the order
						if err := pts.processOrder(ctx, order.PortfolioID, order); err != nil {
							pts.logger.Error(ctx, "Failed to process order", err, map[string]interface{}{
								"portfolio_id": order.PortfolioID,
								"order_id":     orderID,
							})
							order.Status = "failed"
							order.UpdatedAt = time.Now()
						}
					}
				}
			}
			pts.mu.Unlock()
		}
	}
}

// simulateOrderLatency returns a simulated order processing latency
func (pts *PaperTradingService) simulateOrderLatency() time.Duration {
	// Simulate 100ms to 500ms latency
	latencyMs := 100 + rand.Intn(400)
	return time.Duration(latencyMs) * time.Millisecond
}

// processOrder processes a single order
func (pts *PaperTradingService) processOrder(ctx context.Context, portfolioID uuid.UUID, order *PaperOrder) error {
	portfolio, exists := pts.portfolios[portfolioID]
	if !exists {
		return fmt.Errorf("portfolio not found for ID %s", portfolioID)
	}

	// Get current market price (placeholder - would use real market data)
	marketPrice := order.Price.Mul(decimal.NewFromFloat(1.0 + (rand.Float64()-0.5)*0.001)) // ±0.1% slippage

	// Calculate total cost including fees
	totalCost := order.Quantity.Mul(marketPrice)
	fee := totalCost.Mul(pts.config.CommissionRate)

	if order.Side == "buy" {
		// Check if we have enough cash
		totalRequired := totalCost.Add(fee)
		if portfolio.Cash.LessThan(totalRequired) {
			return fmt.Errorf("insufficient cash for buy order")
		}

		// Execute buy order
		portfolio.Cash = portfolio.Cash.Sub(totalRequired)

		// Update position
		position, exists := portfolio.Positions[order.Symbol]
		if !exists {
			position = &PaperPosition{
				Symbol:        order.Symbol,
				Quantity:      decimal.Zero,
				AvgPrice:      decimal.Zero,
				MarketValue:   decimal.Zero,
				UnrealizedPnL: decimal.Zero,
				RealizedPnL:   decimal.Zero,
				LastPrice:     marketPrice,
				UpdatedAt:     time.Now(),
			}
			portfolio.Positions[order.Symbol] = position
		}

		// Calculate new average price
		totalQuantity := position.Quantity.Add(order.Quantity)
		totalValue := position.Quantity.Mul(position.AvgPrice).Add(order.Quantity.Mul(marketPrice))
		position.AvgPrice = totalValue.Div(totalQuantity)
		position.Quantity = totalQuantity
		position.LastPrice = marketPrice
		position.MarketValue = position.Quantity.Mul(marketPrice)
		position.UnrealizedPnL = position.Quantity.Mul(marketPrice.Sub(position.AvgPrice))
		position.UpdatedAt = time.Now()

	} else { // sell
		// Check if we have enough position
		position, exists := portfolio.Positions[order.Symbol]
		if !exists || position.Quantity.LessThan(order.Quantity) {
			return fmt.Errorf("insufficient position for sell order")
		}

		// Execute sell order
		proceeds := totalCost.Sub(fee)
		portfolio.Cash = portfolio.Cash.Add(proceeds)

		// Update position
		position.Quantity = position.Quantity.Sub(order.Quantity)
		if position.Quantity.IsZero() {
			delete(portfolio.Positions, order.Symbol)
		} else {
			position.LastPrice = marketPrice
			position.MarketValue = position.Quantity.Mul(marketPrice)
			position.UnrealizedPnL = position.Quantity.Mul(marketPrice.Sub(position.AvgPrice))
			position.UpdatedAt = time.Now()
		}
	}

	// Update order status
	order.Status = common.OrderStatusFilled
	order.FilledQty = order.Quantity
	order.RemainingQty = decimal.Zero
	order.AvgFillPrice = marketPrice
	order.Commission = fee
	order.UpdatedAt = time.Now()
	order.ExecutedAt = time.Now()

	// Update portfolio
	portfolio.UpdatedAt = time.Now()

	pts.logger.Info(ctx, "Order processed successfully", map[string]interface{}{
		"portfolio_id": portfolioID,
		"order_id":     order.ID,
		"symbol":       order.Symbol,
		"side":         order.Side,
		"quantity":     order.Quantity.String(),
		"filled_price": marketPrice.String(),
		"commission":   fee.String(),
	})

	return nil
}
