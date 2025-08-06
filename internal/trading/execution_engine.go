package trading

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ExecutionEngine handles order execution with advanced algorithms
type ExecutionEngine struct {
	logger        *observability.Logger
	orderQueue    chan *ExecutionOrder
	executionPool *ExecutionPool
	venues        map[string]ExecutionVenue
	router        *SmartOrderRouter
	mu            sync.RWMutex
	isRunning     bool
	stopChan      chan struct{}
	metrics       *ExecutionMetrics
}

// ExecutionOrder represents an order for execution
type ExecutionOrder struct {
	ID              string                 `json:"id"`
	ClientOrderID   string                 `json:"client_order_id"`
	StrategyID      string                 `json:"strategy_id"`
	AlgorithmType   AlgorithmType          `json:"algorithm_type"`
	Symbol          string                 `json:"symbol"`
	Side            OrderSide              `json:"side"`
	OrderType       OrderType              `json:"order_type"`
	Quantity        decimal.Decimal        `json:"quantity"`
	Price           decimal.Decimal        `json:"price"`
	TimeInForce     TimeInForce            `json:"time_in_force"`
	Parameters      map[string]interface{} `json:"parameters"`
	RiskLimits      *RiskLimits            `json:"risk_limits"`
	ExecutionStart  time.Time              `json:"execution_start"`
	ExecutionEnd    time.Time              `json:"execution_end"`
	Status          ExecutionStatus        `json:"status"`
	FilledQuantity  decimal.Decimal        `json:"filled_quantity"`
	AveragePrice    decimal.Decimal        `json:"average_price"`
	TotalSlippage   decimal.Decimal        `json:"total_slippage"`
	TotalCommission decimal.Decimal        `json:"total_commission"`
	Executions      []*ChildExecution      `json:"executions"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// ChildExecution represents a child order execution
type ChildExecution struct {
	ID           string          `json:"id"`
	ParentID     string          `json:"parent_id"`
	Venue        string          `json:"venue"`
	Quantity     decimal.Decimal `json:"quantity"`
	Price        decimal.Decimal `json:"price"`
	Commission   decimal.Decimal `json:"commission"`
	Slippage     decimal.Decimal `json:"slippage"`
	Latency      time.Duration   `json:"latency"`
	ExecutedAt   time.Time       `json:"executed_at"`
	Status       ExecutionStatus `json:"status"`
	ErrorMessage string          `json:"error_message,omitempty"`
}

// ExecutionStatus defines execution status
type ExecutionStatus string

const (
	ExecutionStatusPending    ExecutionStatus = "pending"
	ExecutionStatusExecuting  ExecutionStatus = "executing"
	ExecutionStatusCompleted  ExecutionStatus = "completed"
	ExecutionStatusPartial    ExecutionStatus = "partial"
	ExecutionStatusCanceled   ExecutionStatus = "canceled"
	ExecutionStatusFailed     ExecutionStatus = "failed"
	ExecutionStatusRejected   ExecutionStatus = "rejected"
)

// OrderSide defines order side
type OrderSide string

const (
	OrderSideBuy  OrderSide = "buy"
	OrderSideSell OrderSide = "sell"
)

// OrderType defines order type
type OrderType string

const (
	OrderTypeMarket      OrderType = "market"
	OrderTypeLimit       OrderType = "limit"
	OrderTypeStopLoss    OrderType = "stop_loss"
	OrderTypeTakeProfit  OrderType = "take_profit"
	OrderTypeTrailingStop OrderType = "trailing_stop"
)

// TimeInForce defines time in force
type TimeInForce string

const (
	TimeInForceGTC TimeInForce = "gtc" // Good Till Canceled
	TimeInForceIOC TimeInForce = "ioc" // Immediate Or Cancel
	TimeInForceFOK TimeInForce = "fok" // Fill Or Kill
	TimeInForceGTD TimeInForce = "gtd" // Good Till Date
)

// RiskLimits defines risk limits for execution
type RiskLimits struct {
	MaxSlippageBps    int             `json:"max_slippage_bps"`
	MaxLatencyMs      int             `json:"max_latency_ms"`
	MaxPositionSize   decimal.Decimal `json:"max_position_size"`
	MaxNotionalValue  decimal.Decimal `json:"max_notional_value"`
	MinLiquidityRatio decimal.Decimal `json:"min_liquidity_ratio"`
}

// ExecutionVenue represents a trading venue
type ExecutionVenue interface {
	GetName() string
	GetLatency() time.Duration
	GetLiquidity(symbol string) decimal.Decimal
	GetFeeRate() decimal.Decimal
	ExecuteOrder(ctx context.Context, order *ExecutionOrder) (*ChildExecution, error)
	IsAvailable() bool
}

// ExecutionPool manages concurrent execution workers
type ExecutionPool struct {
	workers   int
	workChan  chan *ExecutionOrder
	resultChan chan *ExecutionResult
	stopChan  chan struct{}
	wg        sync.WaitGroup
}

// ExecutionResult represents the result of an execution
type ExecutionResult struct {
	Order    *ExecutionOrder `json:"order"`
	Success  bool            `json:"success"`
	Error    error           `json:"error,omitempty"`
	Latency  time.Duration   `json:"latency"`
	Slippage decimal.Decimal `json:"slippage"`
}

// ExecutionMetrics tracks execution performance
type ExecutionMetrics struct {
	TotalOrders       int64           `json:"total_orders"`
	CompletedOrders   int64           `json:"completed_orders"`
	FailedOrders      int64           `json:"failed_orders"`
	AverageLatency    time.Duration   `json:"average_latency"`
	AverageSlippage   decimal.Decimal `json:"average_slippage"`
	TotalVolume       decimal.Decimal `json:"total_volume"`
	TotalCommissions  decimal.Decimal `json:"total_commissions"`
	FillRate          float64         `json:"fill_rate"`
	SuccessRate       float64         `json:"success_rate"`
	LastUpdated       time.Time       `json:"last_updated"`
}

// NewExecutionEngine creates a new execution engine
func NewExecutionEngine(logger *observability.Logger) *ExecutionEngine {
	return &ExecutionEngine{
		logger:        logger,
		orderQueue:    make(chan *ExecutionOrder, 1000),
		executionPool: NewExecutionPool(10), // 10 workers
		venues:        make(map[string]ExecutionVenue),
		stopChan:      make(chan struct{}),
		metrics: &ExecutionMetrics{
			LastUpdated: time.Now(),
		},
	}
}

// NewExecutionPool creates a new execution pool
func NewExecutionPool(workers int) *ExecutionPool {
	return &ExecutionPool{
		workers:    workers,
		workChan:   make(chan *ExecutionOrder, 1000),
		resultChan: make(chan *ExecutionResult, 1000),
		stopChan:   make(chan struct{}),
	}
}

// Start starts the execution engine
func (ee *ExecutionEngine) Start(ctx context.Context) error {
	ee.mu.Lock()
	defer ee.mu.Unlock()

	if ee.isRunning {
		return fmt.Errorf("execution engine is already running")
	}

	ee.isRunning = true

	// Start execution pool
	ee.executionPool.Start(ctx, ee)

	// Start order processing
	go ee.processOrders(ctx)
	go ee.processResults(ctx)

	ee.logger.Info(ctx, "Execution engine started", map[string]interface{}{
		"workers": ee.executionPool.workers,
		"venues":  len(ee.venues),
	})

	return nil
}

// Stop stops the execution engine
func (ee *ExecutionEngine) Stop(ctx context.Context) error {
	ee.mu.Lock()
	defer ee.mu.Unlock()

	if !ee.isRunning {
		return nil
	}

	ee.isRunning = false
	close(ee.stopChan)

	// Stop execution pool
	ee.executionPool.Stop()

	ee.logger.Info(ctx, "Execution engine stopped", nil)
	return nil
}

// SubmitOrder submits an order for execution
func (ee *ExecutionEngine) SubmitOrder(ctx context.Context, order *ExecutionOrder) error {
	if !ee.isRunning {
		return fmt.Errorf("execution engine is not running")
	}

	if order.ID == "" {
		order.ID = uuid.New().String()
	}

	order.Status = ExecutionStatusPending
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	select {
	case ee.orderQueue <- order:
		ee.logger.Info(ctx, "Order submitted for execution", map[string]interface{}{
			"order_id":       order.ID,
			"strategy_id":    order.StrategyID,
			"algorithm_type": order.AlgorithmType,
			"symbol":         order.Symbol,
			"quantity":       order.Quantity.String(),
		})
		return nil
	default:
		return fmt.Errorf("order queue is full")
	}
}

// RegisterVenue registers a new execution venue
func (ee *ExecutionEngine) RegisterVenue(venue ExecutionVenue) {
	ee.mu.Lock()
	defer ee.mu.Unlock()

	ee.venues[venue.GetName()] = venue

	ee.logger.Info(context.Background(), "Execution venue registered", map[string]interface{}{
		"venue_name": venue.GetName(),
		"latency":    venue.GetLatency(),
		"fee_rate":   venue.GetFeeRate().String(),
	})
}

// processOrders processes orders from the queue
func (ee *ExecutionEngine) processOrders(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-ee.stopChan:
			return
		case order := <-ee.orderQueue:
			// Route order to execution pool
			select {
			case ee.executionPool.workChan <- order:
				ee.metrics.TotalOrders++
			default:
				ee.logger.Warn(ctx, "Execution pool is full", map[string]interface{}{
					"order_id": order.ID,
				})
			}
		}
	}
}

// processResults processes execution results
func (ee *ExecutionEngine) processResults(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-ee.stopChan:
			return
		case result := <-ee.executionPool.resultChan:
			ee.updateMetrics(result)
			ee.logger.Info(ctx, "Order execution completed", map[string]interface{}{
				"order_id": result.Order.ID,
				"success":  result.Success,
				"latency":  result.Latency,
				"slippage": result.Slippage.String(),
			})
		}
	}
}

// updateMetrics updates execution metrics
func (ee *ExecutionEngine) updateMetrics(result *ExecutionResult) {
	ee.mu.Lock()
	defer ee.mu.Unlock()

	if result.Success {
		ee.metrics.CompletedOrders++
	} else {
		ee.metrics.FailedOrders++
	}

	// Update success rate
	ee.metrics.SuccessRate = float64(ee.metrics.CompletedOrders) / float64(ee.metrics.TotalOrders)

	// Update average latency
	if ee.metrics.TotalOrders > 0 {
		totalLatency := time.Duration(ee.metrics.TotalOrders-1) * ee.metrics.AverageLatency
		ee.metrics.AverageLatency = (totalLatency + result.Latency) / time.Duration(ee.metrics.TotalOrders)
	}

	// Update average slippage
	if ee.metrics.TotalOrders > 0 {
		totalSlippage := ee.metrics.AverageSlippage.Mul(decimal.NewFromInt(ee.metrics.TotalOrders - 1))
		ee.metrics.AverageSlippage = totalSlippage.Add(result.Slippage).Div(decimal.NewFromInt(ee.metrics.TotalOrders))
	}

	ee.metrics.LastUpdated = time.Now()
}

// Start starts the execution pool
func (ep *ExecutionPool) Start(ctx context.Context, engine *ExecutionEngine) {
	for i := 0; i < ep.workers; i++ {
		ep.wg.Add(1)
		go ep.worker(ctx, engine)
	}
}

// Stop stops the execution pool
func (ep *ExecutionPool) Stop() {
	close(ep.stopChan)
	ep.wg.Wait()
}

// worker processes orders in the execution pool
func (ep *ExecutionPool) worker(ctx context.Context, engine *ExecutionEngine) {
	defer ep.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ep.stopChan:
			return
		case order := <-ep.workChan:
			result := ep.executeOrder(ctx, engine, order)
			select {
			case ep.resultChan <- result:
			default:
				// Result channel is full, log warning
				engine.logger.Warn(ctx, "Result channel is full", map[string]interface{}{
					"order_id": order.ID,
				})
			}
		}
	}
}

// executeOrder executes an order using the appropriate algorithm
func (ep *ExecutionPool) executeOrder(ctx context.Context, engine *ExecutionEngine, order *ExecutionOrder) *ExecutionResult {
	start := time.Now()
	order.Status = ExecutionStatusExecuting
	order.ExecutionStart = start

	var err error
	var slippage decimal.Decimal

	// Execute based on algorithm type
	switch order.AlgorithmType {
	case AlgorithmTypeTWAP:
		err = ep.executeTWAP(ctx, engine, order)
	case AlgorithmTypeVWAP:
		err = ep.executeVWAP(ctx, engine, order)
	case AlgorithmTypeIceberg:
		err = ep.executeIceberg(ctx, engine, order)
	case AlgorithmTypeSniper:
		err = ep.executeSniper(ctx, engine, order)
	default:
		err = ep.executeMarket(ctx, engine, order)
	}

	latency := time.Since(start)
	order.ExecutionEnd = time.Now()

	if err != nil {
		order.Status = ExecutionStatusFailed
	} else {
		order.Status = ExecutionStatusCompleted
	}

	return &ExecutionResult{
		Order:    order,
		Success:  err == nil,
		Error:    err,
		Latency:  latency,
		Slippage: slippage,
	}
}

// executeTWAP executes a TWAP algorithm
func (ep *ExecutionPool) executeTWAP(ctx context.Context, engine *ExecutionEngine, order *ExecutionOrder) error {
	// Simplified TWAP implementation
	duration := 60 * time.Minute // Default 1 hour
	if d, ok := order.Parameters["duration_minutes"].(int); ok {
		duration = time.Duration(d) * time.Minute
	}

	sliceCount := 10
	if s, ok := order.Parameters["slice_count"].(int); ok {
		sliceCount = s
	}

	sliceSize := order.Quantity.Div(decimal.NewFromInt(int64(sliceCount)))
	interval := duration / time.Duration(sliceCount)

	for i := 0; i < sliceCount; i++ {
		// Execute slice
		childOrder := &ExecutionOrder{
			ID:        uuid.New().String(),
			Symbol:    order.Symbol,
			Side:      order.Side,
			OrderType: OrderTypeMarket,
			Quantity:  sliceSize,
		}

		execution := &ChildExecution{
			ID:         childOrder.ID,
			ParentID:   order.ID,
			Venue:      "default",
			Quantity:   sliceSize,
			Price:      order.Price, // Simplified
			ExecutedAt: time.Now(),
			Status:     ExecutionStatusCompleted,
		}

		order.Executions = append(order.Executions, execution)
		order.FilledQuantity = order.FilledQuantity.Add(sliceSize)

		// Wait for next slice
		if i < sliceCount-1 {
			time.Sleep(interval)
		}
	}

	return nil
}

// executeVWAP executes a VWAP algorithm
func (ep *ExecutionPool) executeVWAP(ctx context.Context, engine *ExecutionEngine, order *ExecutionOrder) error {
	// Simplified VWAP implementation
	return ep.executeTWAP(ctx, engine, order) // Use TWAP for now
}

// executeIceberg executes an Iceberg algorithm
func (ep *ExecutionPool) executeIceberg(ctx context.Context, engine *ExecutionEngine, order *ExecutionOrder) error {
	// Simplified Iceberg implementation
	visibleSize := decimal.NewFromFloat(0.05) // 5% visible
	if v, ok := order.Parameters["visible_size"].(float64); ok {
		visibleSize = decimal.NewFromFloat(v)
	}

	remaining := order.Quantity
	for remaining.GreaterThan(decimal.Zero) {
		sliceSize := order.Quantity.Mul(visibleSize)
		if sliceSize.GreaterThan(remaining) {
			sliceSize = remaining
		}

		execution := &ChildExecution{
			ID:         uuid.New().String(),
			ParentID:   order.ID,
			Venue:      "default",
			Quantity:   sliceSize,
			Price:      order.Price,
			ExecutedAt: time.Now(),
			Status:     ExecutionStatusCompleted,
		}

		order.Executions = append(order.Executions, execution)
		order.FilledQuantity = order.FilledQuantity.Add(sliceSize)
		remaining = remaining.Sub(sliceSize)

		// Small delay between slices
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

// executeSniper executes a Sniper algorithm
func (ep *ExecutionPool) executeSniper(ctx context.Context, engine *ExecutionEngine, order *ExecutionOrder) error {
	// Simplified Sniper implementation - immediate execution
	return ep.executeMarket(ctx, engine, order)
}

// executeMarket executes a market order
func (ep *ExecutionPool) executeMarket(ctx context.Context, engine *ExecutionEngine, order *ExecutionOrder) error {
	execution := &ChildExecution{
		ID:         uuid.New().String(),
		ParentID:   order.ID,
		Venue:      "default",
		Quantity:   order.Quantity,
		Price:      order.Price,
		ExecutedAt: time.Now(),
		Status:     ExecutionStatusCompleted,
	}

	order.Executions = append(order.Executions, execution)
	order.FilledQuantity = order.Quantity
	order.AveragePrice = order.Price

	return nil
}

// GetMetrics returns execution metrics
func (ee *ExecutionEngine) GetMetrics() *ExecutionMetrics {
	ee.mu.RLock()
	defer ee.mu.RUnlock()

	return ee.metrics
}
