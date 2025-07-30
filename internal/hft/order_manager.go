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

// OrderManager handles high-frequency order execution and management
type OrderManager struct {
	logger      *observability.Logger
	config      HFTConfig
	orders      map[uuid.UUID]*Order
	orderQueue  chan Order
	cancelQueue chan uuid.UUID
	updateChan  chan OrderUpdate

	// Performance tracking
	ordersSubmitted int64
	ordersExecuted  int64
	ordersCanceled  int64
	ordersRejected  int64

	// State management
	isRunning int32
	stopChan  chan struct{}
	wg        sync.WaitGroup
	mu        sync.RWMutex

	// Exchange connections
	exchangeClients map[string]ExchangeClient
}

// Order represents a trading order
type Order struct {
	ID            uuid.UUID              `json:"id"`
	ClientOrderID string                 `json:"client_order_id"`
	Symbol        string                 `json:"symbol"`
	Side          OrderSide              `json:"side"`
	Type          OrderType              `json:"type"`
	Quantity      decimal.Decimal        `json:"quantity"`
	Price         decimal.Decimal        `json:"price"`
	StopPrice     decimal.Decimal        `json:"stop_price,omitempty"`
	TimeInForce   TimeInForce            `json:"time_in_force"`
	Status        OrderStatus            `json:"status"`
	FilledQty     decimal.Decimal        `json:"filled_qty"`
	AvgFillPrice  decimal.Decimal        `json:"avg_fill_price"`
	Commission    decimal.Decimal        `json:"commission"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	Exchange      string                 `json:"exchange"`
	StrategyID    string                 `json:"strategy_id"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// TimeInForce represents order time in force
type TimeInForce string

const (
	TimeInForceGTC TimeInForce = "GTC" // Good Till Canceled
	TimeInForceIOC TimeInForce = "IOC" // Immediate Or Cancel
	TimeInForceFOK TimeInForce = "FOK" // Fill Or Kill
	TimeInForceGTD TimeInForce = "GTD" // Good Till Date
)

// ExchangeClient interface for exchange-specific order execution
type ExchangeClient interface {
	SubmitOrder(ctx context.Context, order Order) (*OrderResponse, error)
	CancelOrder(ctx context.Context, orderID uuid.UUID) error
	GetOrderStatus(ctx context.Context, orderID uuid.UUID) (*Order, error)
	GetOpenOrders(ctx context.Context, symbol string) ([]*Order, error)
}

// OrderResponse represents the response from order submission
type OrderResponse struct {
	OrderID       uuid.UUID   `json:"order_id"`
	ExchangeID    string      `json:"exchange_id"`
	Status        OrderStatus `json:"status"`
	Message       string      `json:"message"`
	LatencyMicros int64       `json:"latency_micros"`
}

// TradingSignal represents a signal to execute a trade
type TradingSignal struct {
	ID         uuid.UUID              `json:"id"`
	Symbol     string                 `json:"symbol"`
	Side       OrderSide              `json:"side"`
	OrderType  OrderType              `json:"order_type"`
	Quantity   decimal.Decimal        `json:"quantity"`
	Price      decimal.Decimal        `json:"price"`
	StopPrice  decimal.Decimal        `json:"stop_price,omitempty"`
	Confidence float64                `json:"confidence"`
	StrategyID string                 `json:"strategy_id"`
	Timestamp  time.Time              `json:"timestamp"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// NewOrderManager creates a new order manager
func NewOrderManager(logger *observability.Logger, config HFTConfig) *OrderManager {
	return &OrderManager{
		logger:          logger,
		config:          config,
		orders:          make(map[uuid.UUID]*Order),
		orderQueue:      make(chan Order, 10000),
		cancelQueue:     make(chan uuid.UUID, 1000),
		updateChan:      make(chan OrderUpdate, 10000),
		stopChan:        make(chan struct{}),
		exchangeClients: make(map[string]ExchangeClient),
	}
}

// Start begins the order manager
func (om *OrderManager) Start(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&om.isRunning, 0, 1) {
		return fmt.Errorf("order manager is already running")
	}

	om.logger.Info(ctx, "Starting order manager", map[string]interface{}{
		"queue_size": cap(om.orderQueue),
		"exchanges":  len(om.exchangeClients),
	})

	// Start processing goroutines
	om.wg.Add(3)
	go om.processOrders(ctx)
	go om.processCancellations(ctx)
	go om.monitorOrders(ctx)

	return nil
}

// Stop gracefully shuts down the order manager
func (om *OrderManager) Stop(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&om.isRunning, 1, 0) {
		return fmt.Errorf("order manager is not running")
	}

	om.logger.Info(ctx, "Stopping order manager", nil)

	close(om.stopChan)
	om.wg.Wait()

	// Cancel all open orders
	om.cancelAllOpenOrders(ctx)

	om.logger.Info(ctx, "Order manager stopped", map[string]interface{}{
		"orders_submitted": atomic.LoadInt64(&om.ordersSubmitted),
		"orders_executed":  atomic.LoadInt64(&om.ordersExecuted),
		"orders_canceled":  atomic.LoadInt64(&om.ordersCanceled),
		"orders_rejected":  atomic.LoadInt64(&om.ordersRejected),
	})

	return nil
}

// SubmitOrder submits a new order for execution
func (om *OrderManager) SubmitOrder(ctx context.Context, order Order) error {
	if atomic.LoadInt32(&om.isRunning) != 1 {
		return fmt.Errorf("order manager is not running")
	}

	// Set default values
	if order.ID == uuid.Nil {
		order.ID = uuid.New()
	}
	if order.ClientOrderID == "" {
		order.ClientOrderID = fmt.Sprintf("HFT_%d", time.Now().UnixNano())
	}
	if order.TimeInForce == "" {
		order.TimeInForce = TimeInForceGTC
	}
	if order.CreatedAt.IsZero() {
		order.CreatedAt = time.Now()
	}

	order.Status = OrderStatusNew
	order.UpdatedAt = time.Now()

	// Store order
	om.mu.Lock()
	om.orders[order.ID] = &order
	om.mu.Unlock()

	// Queue for processing
	select {
	case om.orderQueue <- order:
		atomic.AddInt64(&om.ordersSubmitted, 1)
		return nil
	default:
		return fmt.Errorf("order queue is full")
	}
}

// CancelOrder cancels an existing order
func (om *OrderManager) CancelOrder(ctx context.Context, orderID uuid.UUID) error {
	if atomic.LoadInt32(&om.isRunning) != 1 {
		return fmt.Errorf("order manager is not running")
	}

	select {
	case om.cancelQueue <- orderID:
		return nil
	default:
		return fmt.Errorf("cancel queue is full")
	}
}

// GetOrder retrieves an order by ID
func (om *OrderManager) GetOrder(orderID uuid.UUID) (*Order, error) {
	om.mu.RLock()
	defer om.mu.RUnlock()

	order, exists := om.orders[orderID]
	if !exists {
		return nil, fmt.Errorf("order not found: %s", orderID.String())
	}

	return order, nil
}

// GetOpenOrders returns all open orders
func (om *OrderManager) GetOpenOrders() []*Order {
	om.mu.RLock()
	defer om.mu.RUnlock()

	var openOrders []*Order
	for _, order := range om.orders {
		if order.Status == OrderStatusNew || order.Status == OrderStatusPartialFill {
			openOrders = append(openOrders, order)
		}
	}

	return openOrders
}

// HandleUpdate processes order updates from exchanges
func (om *OrderManager) HandleUpdate(update OrderUpdate) {
	om.mu.Lock()
	defer om.mu.Unlock()

	order, exists := om.orders[update.OrderID]
	if !exists {
		return
	}

	// Update order status
	order.Status = update.Status
	order.FilledQty = update.FilledQty
	order.AvgFillPrice = update.AvgFillPrice
	order.UpdatedAt = update.Timestamp

	// Update metrics
	switch update.Status {
	case OrderStatusFilled:
		atomic.AddInt64(&om.ordersExecuted, 1)
	case OrderStatusCanceled:
		atomic.AddInt64(&om.ordersCanceled, 1)
	case OrderStatusRejected:
		atomic.AddInt64(&om.ordersRejected, 1)
	}
}

// processOrders processes orders from the queue
func (om *OrderManager) processOrders(ctx context.Context) {
	defer om.wg.Done()

	for {
		select {
		case <-om.stopChan:
			return
		case order := <-om.orderQueue:
			if err := om.executeOrder(ctx, order); err != nil {
				om.logger.Error(ctx, "Failed to execute order", err, map[string]interface{}{
					"order_id": order.ID.String(),
					"symbol":   order.Symbol,
					"side":     string(order.Side),
					"quantity": order.Quantity.String(),
				})
			}
		}
	}
}

// processCancellations processes order cancellations
func (om *OrderManager) processCancellations(ctx context.Context) {
	defer om.wg.Done()

	for {
		select {
		case <-om.stopChan:
			return
		case orderID := <-om.cancelQueue:
			if err := om.cancelOrderInternal(ctx, orderID); err != nil {
				om.logger.Error(ctx, "Failed to cancel order", err, map[string]interface{}{
					"order_id": orderID.String(),
				})
			}
		}
	}
}

// monitorOrders monitors order timeouts and status
func (om *OrderManager) monitorOrders(ctx context.Context) {
	defer om.wg.Done()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-om.stopChan:
			return
		case <-ticker.C:
			om.checkOrderTimeouts(ctx)
		}
	}
}

// executeOrder executes an order on the appropriate exchange
func (om *OrderManager) executeOrder(ctx context.Context, order Order) error {
	start := time.Now()

	// Get exchange client
	client, exists := om.exchangeClients[order.Exchange]
	if !exists {
		return fmt.Errorf("exchange client not found: %s", order.Exchange)
	}

	// Submit order to exchange
	response, err := client.SubmitOrder(ctx, order)
	if err != nil {
		return fmt.Errorf("failed to submit order to exchange: %w", err)
	}

	// Track latency
	latency := time.Since(start).Microseconds()

	om.logger.Info(ctx, "Order submitted to exchange", map[string]interface{}{
		"order_id":       order.ID.String(),
		"exchange_id":    response.ExchangeID,
		"status":         string(response.Status),
		"latency_micros": latency,
	})

	return nil
}

// cancelOrderInternal cancels an order internally
func (om *OrderManager) cancelOrderInternal(ctx context.Context, orderID uuid.UUID) error {
	om.mu.RLock()
	order, exists := om.orders[orderID]
	om.mu.RUnlock()

	if !exists {
		return fmt.Errorf("order not found: %s", orderID.String())
	}

	// Get exchange client
	client, exists := om.exchangeClients[order.Exchange]
	if !exists {
		return fmt.Errorf("exchange client not found: %s", order.Exchange)
	}

	// Cancel order on exchange
	return client.CancelOrder(ctx, orderID)
}

// checkOrderTimeouts checks for and handles order timeouts
func (om *OrderManager) checkOrderTimeouts(ctx context.Context) {
	om.mu.RLock()
	defer om.mu.RUnlock()

	timeout := time.Duration(om.config.OrderTimeoutMs) * time.Millisecond
	now := time.Now()

	for _, order := range om.orders {
		if order.Status == OrderStatusNew || order.Status == OrderStatusPartialFill {
			if now.Sub(order.CreatedAt) > timeout {
				go func(orderID uuid.UUID) {
					if err := om.CancelOrder(ctx, orderID); err != nil {
						om.logger.Error(ctx, "Failed to cancel timed out order", err)
					}
				}(order.ID)
			}
		}
	}
}

// cancelAllOpenOrders cancels all open orders
func (om *OrderManager) cancelAllOpenOrders(ctx context.Context) {
	openOrders := om.GetOpenOrders()

	for _, order := range openOrders {
		if err := om.CancelOrder(ctx, order.ID); err != nil {
			om.logger.Error(ctx, "Failed to cancel order during shutdown", err)
		}
	}
}

// AddExchangeClient adds an exchange client
func (om *OrderManager) AddExchangeClient(exchange string, client ExchangeClient) {
	om.exchangeClients[exchange] = client
}

// GetMetrics returns order manager metrics
func (om *OrderManager) GetMetrics() OrderManagerMetrics {
	return OrderManagerMetrics{
		OrdersSubmitted: atomic.LoadInt64(&om.ordersSubmitted),
		OrdersExecuted:  atomic.LoadInt64(&om.ordersExecuted),
		OrdersCanceled:  atomic.LoadInt64(&om.ordersCanceled),
		OrdersRejected:  atomic.LoadInt64(&om.ordersRejected),
		OpenOrders:      len(om.GetOpenOrders()),
		IsRunning:       atomic.LoadInt32(&om.isRunning) == 1,
	}
}

// OrderManagerMetrics contains order manager performance metrics
type OrderManagerMetrics struct {
	OrdersSubmitted int64 `json:"orders_submitted"`
	OrdersExecuted  int64 `json:"orders_executed"`
	OrdersCanceled  int64 `json:"orders_canceled"`
	OrdersRejected  int64 `json:"orders_rejected"`
	OpenOrders      int   `json:"open_orders"`
	IsRunning       bool  `json:"is_running"`
}
