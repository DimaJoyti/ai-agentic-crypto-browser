package exchanges

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/internal/exchanges/common"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// OrderManager manages orders across multiple exchanges with advanced routing
type OrderManager struct {
	logger     *observability.Logger
	manager    *Manager
	orders     map[uuid.UUID]*ManagedOrder
	strategies map[string]ExecutionStrategy
	config     OrderManagerConfig

	// Performance tracking
	ordersSubmitted int64
	ordersExecuted  int64
	ordersCanceled  int64
	ordersRejected  int64

	// State management
	isRunning bool
	stopChan  chan struct{}
	wg        sync.WaitGroup
	mu        sync.RWMutex
}

// OrderManagerConfig contains configuration for the order manager
type OrderManagerConfig struct {
	MaxOrdersPerSecond   int             `json:"max_orders_per_second"`
	OrderTimeoutSeconds  int             `json:"order_timeout_seconds"`
	EnableSmartRouting   bool            `json:"enable_smart_routing"`
	EnableOrderSplitting bool            `json:"enable_order_splitting"`
	MaxOrderSize         decimal.Decimal `json:"max_order_size"`
	MinOrderSize         decimal.Decimal `json:"min_order_size"`
	DefaultStrategy      string          `json:"default_strategy"`
}

// ManagedOrder represents an order managed by the order manager
type ManagedOrder struct {
	ID              uuid.UUID              `json:"id"`
	OriginalRequest *common.OrderRequest   `json:"original_request"`
	SubOrders       []*SubOrder            `json:"sub_orders"`
	Strategy        string                 `json:"strategy"`
	Status          OrderStatus            `json:"status"`
	TotalFilled     decimal.Decimal        `json:"total_filled"`
	AvgFillPrice    decimal.Decimal        `json:"avg_fill_price"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// SubOrder represents a sub-order sent to a specific exchange
type SubOrder struct {
	ID         uuid.UUID             `json:"id"`
	Exchange   string                `json:"exchange"`
	ExchangeID string                `json:"exchange_id"`
	Request    *common.OrderRequest  `json:"request"`
	Response   *common.OrderResponse `json:"response"`
	Status     OrderStatus           `json:"status"`
	CreatedAt  time.Time             `json:"created_at"`
	UpdatedAt  time.Time             `json:"updated_at"`
}

// OrderStatus represents the status of a managed order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusSubmitted OrderStatus = "SUBMITTED"
	OrderStatusPartial   OrderStatus = "PARTIAL"
	OrderStatusCompleted OrderStatus = "COMPLETED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
	OrderStatusFailed    OrderStatus = "FAILED"
)

// ExecutionStrategy defines how orders should be executed
type ExecutionStrategy interface {
	Execute(ctx context.Context, order *ManagedOrder, manager *Manager) error
	GetName() string
}

// NewOrderManager creates a new order manager
func NewOrderManager(logger *observability.Logger, manager *Manager, config OrderManagerConfig) *OrderManager {
	om := &OrderManager{
		logger:     logger,
		manager:    manager,
		orders:     make(map[uuid.UUID]*ManagedOrder),
		strategies: make(map[string]ExecutionStrategy),
		config:     config,
		stopChan:   make(chan struct{}),
	}

	// Register default strategies
	om.registerDefaultStrategies()

	return om
}

// Start starts the order manager
func (om *OrderManager) Start(ctx context.Context) error {
	om.mu.Lock()
	defer om.mu.Unlock()

	if om.isRunning {
		return nil
	}

	om.logger.Info(ctx, "Starting order manager", map[string]interface{}{
		"max_orders_per_second": om.config.MaxOrdersPerSecond,
		"enable_smart_routing":  om.config.EnableSmartRouting,
		"default_strategy":      om.config.DefaultStrategy,
	})

	om.isRunning = true

	// Start order monitoring
	om.wg.Add(1)
	go om.monitorOrders(ctx)

	om.logger.Info(ctx, "Order manager started", nil)

	return nil
}

// Stop stops the order manager
func (om *OrderManager) Stop(ctx context.Context) error {
	om.mu.Lock()
	defer om.mu.Unlock()

	if !om.isRunning {
		return nil
	}

	om.logger.Info(ctx, "Stopping order manager", nil)

	close(om.stopChan)
	om.wg.Wait()

	om.isRunning = false

	om.logger.Info(ctx, "Order manager stopped", nil)

	return nil
}

// SubmitOrder submits an order for execution
func (om *OrderManager) SubmitOrder(ctx context.Context, request *common.OrderRequest) (*ManagedOrder, error) {
	om.mu.Lock()
	defer om.mu.Unlock()

	if !om.isRunning {
		return nil, fmt.Errorf("order manager is not running")
	}

	// Validate order
	if err := om.validateOrder(request); err != nil {
		om.ordersRejected++
		return nil, fmt.Errorf("order validation failed: %w", err)
	}

	// Create managed order
	managedOrder := &ManagedOrder{
		ID:              uuid.New(),
		OriginalRequest: request,
		SubOrders:       make([]*SubOrder, 0),
		Strategy:        om.getStrategy(request),
		Status:          OrderStatusPending,
		TotalFilled:     decimal.NewFromInt(0),
		AvgFillPrice:    decimal.NewFromInt(0),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Metadata:        make(map[string]interface{}),
	}

	// Store order
	om.orders[managedOrder.ID] = managedOrder
	om.ordersSubmitted++

	// Execute order asynchronously
	go om.executeOrder(ctx, managedOrder)

	om.logger.Info(ctx, "Order submitted", map[string]interface{}{
		"order_id": managedOrder.ID.String(),
		"symbol":   request.Symbol,
		"side":     string(request.Side),
		"quantity": request.Quantity.String(),
		"strategy": managedOrder.Strategy,
	})

	return managedOrder, nil
}

// CancelOrder cancels a managed order
func (om *OrderManager) CancelOrder(ctx context.Context, orderID uuid.UUID) error {
	om.mu.Lock()
	defer om.mu.Unlock()

	order, exists := om.orders[orderID]
	if !exists {
		return fmt.Errorf("order not found: %s", orderID.String())
	}

	if order.Status == OrderStatusCompleted || order.Status == OrderStatusCancelled {
		return fmt.Errorf("order cannot be cancelled: %s", order.Status)
	}

	// Cancel all sub-orders
	for _, subOrder := range order.SubOrders {
		if subOrder.Status == OrderStatusSubmitted || subOrder.Status == OrderStatusPartial {
			exchange, err := om.manager.GetExchange(subOrder.Exchange)
			if err != nil {
				om.logger.Error(ctx, "Failed to get exchange for cancellation", err, map[string]interface{}{
					"exchange":     subOrder.Exchange,
					"sub_order_id": subOrder.ID.String(),
				})
				continue
			}

			_, err = exchange.CancelOrder(ctx, subOrder.Request.Symbol, subOrder.ExchangeID)
			if err != nil {
				om.logger.Error(ctx, "Failed to cancel sub-order", err, map[string]interface{}{
					"exchange":     subOrder.Exchange,
					"exchange_id":  subOrder.ExchangeID,
					"sub_order_id": subOrder.ID.String(),
				})
			} else {
				subOrder.Status = OrderStatusCancelled
				subOrder.UpdatedAt = time.Now()
			}
		}
	}

	order.Status = OrderStatusCancelled
	order.UpdatedAt = time.Now()
	om.ordersCanceled++

	om.logger.Info(ctx, "Order cancelled", map[string]interface{}{
		"order_id": orderID.String(),
	})

	return nil
}

// GetOrder gets a managed order by ID
func (om *OrderManager) GetOrder(orderID uuid.UUID) (*ManagedOrder, error) {
	om.mu.RLock()
	defer om.mu.RUnlock()

	order, exists := om.orders[orderID]
	if !exists {
		return nil, fmt.Errorf("order not found: %s", orderID.String())
	}

	return order, nil
}

// GetOrders gets all managed orders
func (om *OrderManager) GetOrders() []*ManagedOrder {
	om.mu.RLock()
	defer om.mu.RUnlock()

	orders := make([]*ManagedOrder, 0, len(om.orders))
	for _, order := range om.orders {
		orders = append(orders, order)
	}

	return orders
}

// GetOrdersByStatus gets orders by status
func (om *OrderManager) GetOrdersByStatus(status OrderStatus) []*ManagedOrder {
	om.mu.RLock()
	defer om.mu.RUnlock()

	orders := make([]*ManagedOrder, 0)
	for _, order := range om.orders {
		if order.Status == status {
			orders = append(orders, order)
		}
	}

	return orders
}

// GetMetrics returns order manager metrics
func (om *OrderManager) GetMetrics() OrderManagerMetrics {
	om.mu.RLock()
	defer om.mu.RUnlock()

	return OrderManagerMetrics{
		OrdersSubmitted: om.ordersSubmitted,
		OrdersExecuted:  om.ordersExecuted,
		OrdersCanceled:  om.ordersCanceled,
		OrdersRejected:  om.ordersRejected,
		ActiveOrders:    int64(len(om.orders)),
	}
}

// OrderManagerMetrics contains order manager performance metrics
type OrderManagerMetrics struct {
	OrdersSubmitted int64 `json:"orders_submitted"`
	OrdersExecuted  int64 `json:"orders_executed"`
	OrdersCanceled  int64 `json:"orders_canceled"`
	OrdersRejected  int64 `json:"orders_rejected"`
	ActiveOrders    int64 `json:"active_orders"`
}

// Private methods

// validateOrder validates an order request
func (om *OrderManager) validateOrder(request *common.OrderRequest) error {
	if request.Symbol == "" {
		return fmt.Errorf("symbol is required")
	}

	if request.Quantity.IsZero() || request.Quantity.IsNegative() {
		return fmt.Errorf("quantity must be positive")
	}

	if request.Quantity.GreaterThan(om.config.MaxOrderSize) {
		return fmt.Errorf("quantity exceeds maximum order size")
	}

	if request.Quantity.LessThan(om.config.MinOrderSize) {
		return fmt.Errorf("quantity below minimum order size")
	}

	return nil
}

// getStrategy determines the execution strategy for an order
func (om *OrderManager) getStrategy(request *common.OrderRequest) string {
	// Check if strategy is specified in metadata
	if request.Metadata != nil {
		if strategy, exists := request.Metadata["strategy"]; exists {
			if strategyStr, ok := strategy.(string); ok {
				if _, exists := om.strategies[strategyStr]; exists {
					return strategyStr
				}
			}
		}
	}

	return om.config.DefaultStrategy
}

// executeOrder executes a managed order using the specified strategy
func (om *OrderManager) executeOrder(ctx context.Context, order *ManagedOrder) {
	strategy, exists := om.strategies[order.Strategy]
	if !exists {
		om.logger.Error(ctx, "Strategy not found", nil, map[string]interface{}{
			"order_id": order.ID.String(),
			"strategy": order.Strategy,
		})
		order.Status = OrderStatusFailed
		order.UpdatedAt = time.Now()
		return
	}

	order.Status = OrderStatusSubmitted
	order.UpdatedAt = time.Now()

	if err := strategy.Execute(ctx, order, om.manager); err != nil {
		om.logger.Error(ctx, "Order execution failed", err, map[string]interface{}{
			"order_id": order.ID.String(),
			"strategy": order.Strategy,
		})
		order.Status = OrderStatusFailed
		order.UpdatedAt = time.Now()
		return
	}

	om.ordersExecuted++
}

// monitorOrders monitors order status and updates
func (om *OrderManager) monitorOrders(ctx context.Context) {
	defer om.wg.Done()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-om.stopChan:
			return
		case <-ticker.C:
			om.updateOrderStatuses(ctx)
		}
	}
}

// updateOrderStatuses updates the status of all active orders
func (om *OrderManager) updateOrderStatuses(ctx context.Context) {
	om.mu.Lock()
	defer om.mu.Unlock()

	for _, order := range om.orders {
		if order.Status == OrderStatusSubmitted || order.Status == OrderStatusPartial {
			om.updateOrderStatus(ctx, order)
		}
	}
}

// updateOrderStatus updates the status of a single order
func (om *OrderManager) updateOrderStatus(ctx context.Context, order *ManagedOrder) {
	totalFilled := decimal.NewFromInt(0)
	totalValue := decimal.NewFromInt(0)
	allCompleted := true

	for _, subOrder := range order.SubOrders {
		if subOrder.Status == OrderStatusSubmitted || subOrder.Status == OrderStatusPartial {
			// Query exchange for order status
			exchange, err := om.manager.GetExchange(subOrder.Exchange)
			if err != nil {
				continue
			}

			response, err := exchange.GetOrder(ctx, subOrder.Request.Symbol, subOrder.ExchangeID)
			if err != nil {
				continue
			}

			// Update sub-order status
			subOrder.Response = response
			subOrder.UpdatedAt = time.Now()

			switch response.Status {
			case common.OrderStatusFilled:
				subOrder.Status = OrderStatusCompleted
			case common.OrderStatusPartiallyFilled:
				subOrder.Status = OrderStatusPartial
				allCompleted = false
			case common.OrderStatusCanceled:
				subOrder.Status = OrderStatusCancelled
			case common.OrderStatusRejected:
				subOrder.Status = OrderStatusFailed
			default:
				allCompleted = false
			}
		}

		if subOrder.Response != nil {
			totalFilled = totalFilled.Add(subOrder.Response.FilledQty)
			if !subOrder.Response.FilledQty.IsZero() && !subOrder.Response.AvgFillPrice.IsZero() {
				totalValue = totalValue.Add(subOrder.Response.FilledQty.Mul(subOrder.Response.AvgFillPrice))
			}
		}

		if subOrder.Status != OrderStatusCompleted && subOrder.Status != OrderStatusCancelled && subOrder.Status != OrderStatusFailed {
			allCompleted = false
		}
	}

	// Update order totals
	order.TotalFilled = totalFilled
	if !totalFilled.IsZero() {
		order.AvgFillPrice = totalValue.Div(totalFilled)
	}

	// Update order status
	if allCompleted {
		if totalFilled.Equal(order.OriginalRequest.Quantity) {
			order.Status = OrderStatusCompleted
		} else if totalFilled.GreaterThan(decimal.NewFromInt(0)) {
			order.Status = OrderStatusPartial
		} else {
			order.Status = OrderStatusCancelled
		}
	} else if totalFilled.GreaterThan(decimal.NewFromInt(0)) {
		order.Status = OrderStatusPartial
	}

	order.UpdatedAt = time.Now()
}

// registerDefaultStrategies registers default execution strategies
func (om *OrderManager) registerDefaultStrategies() {
	om.strategies["simple"] = &SimpleStrategy{}
	om.strategies["smart_routing"] = &SmartRoutingStrategy{}
	om.strategies["twap"] = &TWAPStrategy{}
}

// Execution Strategies

// SimpleStrategy executes orders on the default exchange
type SimpleStrategy struct{}

func (s *SimpleStrategy) GetName() string {
	return "simple"
}

func (s *SimpleStrategy) Execute(ctx context.Context, order *ManagedOrder, manager *Manager) error {
	exchange, err := manager.GetDefaultExchange()
	if err != nil {
		return fmt.Errorf("failed to get default exchange: %w", err)
	}

	subOrder := &SubOrder{
		ID:        uuid.New(),
		Exchange:  exchange.GetExchangeName(),
		Request:   order.OriginalRequest,
		Status:    OrderStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	response, err := exchange.PlaceOrder(ctx, order.OriginalRequest)
	if err != nil {
		subOrder.Status = OrderStatusFailed
		order.SubOrders = append(order.SubOrders, subOrder)
		return fmt.Errorf("failed to place order: %w", err)
	}

	subOrder.ExchangeID = response.OrderID
	subOrder.Response = response
	subOrder.Status = OrderStatusSubmitted
	subOrder.UpdatedAt = time.Now()

	order.SubOrders = append(order.SubOrders, subOrder)

	return nil
}

// SmartRoutingStrategy routes orders to the exchange with the best price
type SmartRoutingStrategy struct{}

func (s *SmartRoutingStrategy) GetName() string {
	return "smart_routing"
}

func (s *SmartRoutingStrategy) Execute(ctx context.Context, order *ManagedOrder, manager *Manager) error {
	bestPrice, err := manager.GetBestPrice(ctx, order.OriginalRequest.Symbol, order.OriginalRequest.Side)
	if err != nil {
		return fmt.Errorf("failed to get best price: %w", err)
	}

	exchange, err := manager.GetExchange(bestPrice.Exchange)
	if err != nil {
		return fmt.Errorf("failed to get exchange: %w", err)
	}

	subOrder := &SubOrder{
		ID:        uuid.New(),
		Exchange:  bestPrice.Exchange,
		Request:   order.OriginalRequest,
		Status:    OrderStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	response, err := exchange.PlaceOrder(ctx, order.OriginalRequest)
	if err != nil {
		subOrder.Status = OrderStatusFailed
		order.SubOrders = append(order.SubOrders, subOrder)
		return fmt.Errorf("failed to place order: %w", err)
	}

	subOrder.ExchangeID = response.OrderID
	subOrder.Response = response
	subOrder.Status = OrderStatusSubmitted
	subOrder.UpdatedAt = time.Now()

	order.SubOrders = append(order.SubOrders, subOrder)

	return nil
}

// TWAPStrategy implements Time-Weighted Average Price execution
type TWAPStrategy struct{}

func (s *TWAPStrategy) GetName() string {
	return "twap"
}

func (s *TWAPStrategy) Execute(ctx context.Context, order *ManagedOrder, manager *Manager) error {
	// This is a simplified TWAP implementation
	// In practice, this would split the order into smaller chunks over time

	exchange, err := manager.GetDefaultExchange()
	if err != nil {
		return fmt.Errorf("failed to get default exchange: %w", err)
	}

	// For now, just execute as a single order
	// TODO: Implement proper TWAP logic with time-based splitting
	subOrder := &SubOrder{
		ID:        uuid.New(),
		Exchange:  exchange.GetExchangeName(),
		Request:   order.OriginalRequest,
		Status:    OrderStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	response, err := exchange.PlaceOrder(ctx, order.OriginalRequest)
	if err != nil {
		subOrder.Status = OrderStatusFailed
		order.SubOrders = append(order.SubOrders, subOrder)
		return fmt.Errorf("failed to place order: %w", err)
	}

	subOrder.ExchangeID = response.OrderID
	subOrder.Response = response
	subOrder.Status = OrderStatusSubmitted
	subOrder.UpdatedAt = time.Now()

	order.SubOrders = append(order.SubOrders, subOrder)

	return nil
}
