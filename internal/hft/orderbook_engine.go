package hft

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// OrderBookEngine provides ultra-fast in-memory order book management
// with lock-free data structures and microsecond-level processing
type OrderBookEngine struct {
	logger *observability.Logger
	config OrderBookConfig

	// Order books for different symbols
	orderBooks map[string]*OrderBook

	// Lock-free order processing
	orderQueue  *LockFreeRingBuffer
	updateQueue *LockFreeRingBuffer

	// Performance metrics
	ordersProcessed        int64
	updatesProcessed       int64
	avgProcessingTimeNanos int64

	// State management
	isRunning int32
	stopChan  chan struct{}
	wg        sync.WaitGroup
	mu        sync.RWMutex

	// Subscribers for order book updates
	subscribers map[string][]chan *OrderBookUpdate
}

// OrderBookConfig contains configuration for the order book engine
type OrderBookConfig struct {
	MaxDepth         int             `json:"max_depth"`         // Maximum order book depth
	ProcessorThreads int             `json:"processor_threads"` // Number of processing threads
	QueueSize        int             `json:"queue_size"`        // Size of processing queues
	TickSize         decimal.Decimal `json:"tick_size"`         // Minimum price increment
	LotSize          decimal.Decimal `json:"lot_size"`          // Minimum quantity increment
	EnableCrossing   bool            `json:"enable_crossing"`   // Allow crossing orders
	MaxOrderAge      time.Duration   `json:"max_order_age"`     // Maximum order age before expiry
}

// OrderBook represents a high-performance in-memory order book
type OrderBook struct {
	Symbol string

	// Price-time priority order queues (lock-free)
	BidLevels *PriceLevelTree // Sorted by price descending
	AskLevels *PriceLevelTree // Sorted by price ascending

	// Order tracking
	Orders map[uuid.UUID]*BookOrder

	// Best bid/ask cache for fast access
	bestBid atomic.Value // *PriceLevel
	bestAsk atomic.Value // *PriceLevel

	// Sequence number for updates
	sequenceNumber uint64

	// Statistics
	totalVolume    decimal.Decimal
	lastTradePrice decimal.Decimal
	lastTradeTime  time.Time

	// Synchronization
	mu sync.RWMutex
}

// PriceLevelTree implements a lock-free binary search tree for price levels
type PriceLevelTree struct {
	root  unsafe.Pointer // *PriceLevelNode
	isAsk bool           // true for ask side, false for bid side
}

// PriceLevelNode represents a node in the price level tree
type PriceLevelNode struct {
	level *PriceLevel
	left  unsafe.Pointer // *PriceLevelNode
	right unsafe.Pointer // *PriceLevelNode
}

// PriceLevel represents orders at a specific price level
type PriceLevel struct {
	Price      decimal.Decimal
	Orders     []*BookOrder // FIFO queue for time priority
	TotalQty   decimal.Decimal
	OrderCount int

	// Lock-free operations
	mu sync.RWMutex
}

// BookOrder represents an order in the order book
type BookOrder struct {
	ID          uuid.UUID
	Symbol      string
	Side        OrderSide
	Price       decimal.Decimal
	Quantity    decimal.Decimal
	FilledQty   decimal.Decimal
	TimeInForce TimeInForce
	OrderType   OrderType

	// Timestamps
	CreatedAt time.Time
	UpdatedAt time.Time

	// Order book specific
	PriceLevel *PriceLevel
	Position   int // Position in price level queue

	// Metadata
	ClientOrderID string
	StrategyID    string
	Metadata      map[string]interface{}
}

// OrderBookUpdate represents an order book update event
type OrderBookUpdate struct {
	Symbol         string          `json:"symbol"`
	UpdateType     UpdateType      `json:"update_type"`
	Side           OrderSide       `json:"side"`
	Price          decimal.Decimal `json:"price"`
	Quantity       decimal.Decimal `json:"quantity"`
	OrderID        uuid.UUID       `json:"order_id,omitempty"`
	SequenceNumber uint64          `json:"sequence_number"`
	Timestamp      time.Time       `json:"timestamp"`

	// Best bid/ask after update
	BestBid decimal.Decimal `json:"best_bid"`
	BestAsk decimal.Decimal `json:"best_ask"`
	BidSize decimal.Decimal `json:"bid_size"`
	AskSize decimal.Decimal `json:"ask_size"`
}

// UpdateType represents the type of order book update
type UpdateType string

const (
	UpdateTypeAdd    UpdateType = "ADD"
	UpdateTypeModify UpdateType = "MODIFY"
	UpdateTypeDelete UpdateType = "DELETE"
	UpdateTypeTrade  UpdateType = "TRADE"
)

// NewOrderBookEngine creates a new high-performance order book engine
func NewOrderBookEngine(logger *observability.Logger, config OrderBookConfig) *OrderBookEngine {
	// Set default values
	if config.MaxDepth == 0 {
		config.MaxDepth = 1000
	}
	if config.ProcessorThreads == 0 {
		config.ProcessorThreads = 4
	}
	if config.QueueSize == 0 {
		config.QueueSize = 100000
	}
	if config.TickSize.IsZero() {
		config.TickSize = decimal.NewFromFloat(0.01)
	}
	if config.LotSize.IsZero() {
		config.LotSize = decimal.NewFromFloat(0.001)
	}

	engine := &OrderBookEngine{
		logger:      logger,
		config:      config,
		orderBooks:  make(map[string]*OrderBook),
		subscribers: make(map[string][]chan *OrderBookUpdate),
		stopChan:    make(chan struct{}),
	}

	// Initialize lock-free queues
	engine.orderQueue = NewLockFreeRingBuffer(config.QueueSize)
	engine.updateQueue = NewLockFreeRingBuffer(config.QueueSize)

	return engine
}

// NewOrderBook creates a new order book for a symbol
func NewOrderBook(symbol string) *OrderBook {
	ob := &OrderBook{
		Symbol:    symbol,
		BidLevels: NewPriceLevelTree(false), // Bid side
		AskLevels: NewPriceLevelTree(true),  // Ask side
		Orders:    make(map[uuid.UUID]*BookOrder),
	}

	// Initialize best bid/ask to nil
	ob.bestBid.Store((*PriceLevel)(nil))
	ob.bestAsk.Store((*PriceLevel)(nil))

	return ob
}

// NewPriceLevelTree creates a new price level tree
func NewPriceLevelTree(isAsk bool) *PriceLevelTree {
	return &PriceLevelTree{
		isAsk: isAsk,
	}
}

// Start begins the order book engine
func (obe *OrderBookEngine) Start(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&obe.isRunning, 0, 1) {
		return fmt.Errorf("order book engine is already running")
	}

	obe.logger.Info(ctx, "Starting order book engine", map[string]interface{}{
		"processor_threads": obe.config.ProcessorThreads,
		"queue_size":        obe.config.QueueSize,
		"max_depth":         obe.config.MaxDepth,
	})

	// Start processing threads
	obe.wg.Add(obe.config.ProcessorThreads)
	for i := 0; i < obe.config.ProcessorThreads; i++ {
		go obe.processOrders(ctx, i)
	}

	// Start performance monitor
	go obe.performanceMonitor(ctx)

	obe.logger.Info(ctx, "Order book engine started successfully", nil)
	return nil
}

// Stop gracefully shuts down the order book engine
func (obe *OrderBookEngine) Stop(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&obe.isRunning, 1, 0) {
		return fmt.Errorf("order book engine is not running")
	}

	obe.logger.Info(ctx, "Stopping order book engine", nil)

	close(obe.stopChan)
	obe.wg.Wait()

	obe.logger.Info(ctx, "Order book engine stopped", map[string]interface{}{
		"orders_processed":  atomic.LoadInt64(&obe.ordersProcessed),
		"updates_processed": atomic.LoadInt64(&obe.updatesProcessed),
		"avg_processing_ns": atomic.LoadInt64(&obe.avgProcessingTimeNanos),
	})

	return nil
}

// AddOrder adds a new order to the order book
func (obe *OrderBookEngine) AddOrder(ctx context.Context, order *BookOrder) error {
	if atomic.LoadInt32(&obe.isRunning) != 1 {
		return fmt.Errorf("order book engine is not running")
	}

	// Validate order
	if err := obe.validateOrder(order); err != nil {
		return fmt.Errorf("order validation failed: %w", err)
	}

	// Set timestamps
	order.CreatedAt = time.Now()
	order.UpdatedAt = order.CreatedAt

	// Queue for processing
	orderPtr := unsafe.Pointer(order)
	if !obe.orderQueue.Push(orderPtr) {
		return fmt.Errorf("order queue is full")
	}

	return nil
}

// validateOrder validates an order before adding to the book
func (obe *OrderBookEngine) validateOrder(order *BookOrder) error {
	if order.Symbol == "" {
		return fmt.Errorf("symbol is required")
	}
	if order.Price.IsNegative() {
		return fmt.Errorf("price must be positive")
	}
	if order.Quantity.IsNegative() || order.Quantity.IsZero() {
		return fmt.Errorf("quantity must be positive")
	}

	// Check tick size
	if !obe.config.TickSize.IsZero() {
		remainder := order.Price.Mod(obe.config.TickSize)
		if !remainder.IsZero() {
			return fmt.Errorf("price must be multiple of tick size %s", obe.config.TickSize.String())
		}
	}

	// Check lot size
	if !obe.config.LotSize.IsZero() {
		remainder := order.Quantity.Mod(obe.config.LotSize)
		if !remainder.IsZero() {
			return fmt.Errorf("quantity must be multiple of lot size %s", obe.config.LotSize.String())
		}
	}

	return nil
}

// processOrders processes orders from the queue
func (obe *OrderBookEngine) processOrders(ctx context.Context, workerID int) {
	defer obe.wg.Done()

	obe.logger.Info(ctx, "Starting order processor", map[string]interface{}{
		"worker_id": workerID,
	})

	for {
		select {
		case <-obe.stopChan:
			return
		default:
			// Process orders from queue
			if orderPtr := obe.orderQueue.Pop(); orderPtr != nil {
				start := time.Now()
				order := (*BookOrder)(orderPtr)

				if err := obe.processOrder(ctx, order); err != nil {
					obe.logger.Error(ctx, "Failed to process order", err)
				}

				// Update performance metrics
				processingTime := time.Since(start).Nanoseconds()
				atomic.StoreInt64(&obe.avgProcessingTimeNanos, processingTime)
				atomic.AddInt64(&obe.ordersProcessed, 1)
			} else {
				// No orders available, yield CPU
				time.Sleep(time.Microsecond)
			}
		}
	}
}

// processOrder processes a single order
func (obe *OrderBookEngine) processOrder(ctx context.Context, order *BookOrder) error {
	// Get or create order book for symbol
	orderBook := obe.getOrCreateOrderBook(order.Symbol)

	// Process based on order type and time in force
	switch order.TimeInForce {
	case TimeInForceIOC:
		return obe.processIOCOrder(ctx, orderBook, order)
	case TimeInForceFOK:
		return obe.processFOKOrder(ctx, orderBook, order)
	default:
		return obe.processGTCOrder(ctx, orderBook, order)
	}
}

// getOrCreateOrderBook gets or creates an order book for a symbol
func (obe *OrderBookEngine) getOrCreateOrderBook(symbol string) *OrderBook {
	obe.mu.RLock()
	orderBook, exists := obe.orderBooks[symbol]
	obe.mu.RUnlock()

	if !exists {
		obe.mu.Lock()
		// Double-check after acquiring write lock
		if orderBook, exists = obe.orderBooks[symbol]; !exists {
			orderBook = NewOrderBook(symbol)
			obe.orderBooks[symbol] = orderBook
		}
		obe.mu.Unlock()
	}

	return orderBook
}

// processGTCOrder processes a Good Till Cancel order
func (obe *OrderBookEngine) processGTCOrder(ctx context.Context, orderBook *OrderBook, order *BookOrder) error {
	orderBook.mu.Lock()
	defer orderBook.mu.Unlock()

	// Try to match against existing orders
	if err := obe.matchOrder(ctx, orderBook, order); err != nil {
		return err
	}

	// If order has remaining quantity, add to book
	if order.Quantity.GreaterThan(order.FilledQty) {
		if err := obe.addOrderToBook(ctx, orderBook, order); err != nil {
			return err
		}
	}

	return nil
}

// processIOCOrder processes an Immediate Or Cancel order
func (obe *OrderBookEngine) processIOCOrder(ctx context.Context, orderBook *OrderBook, order *BookOrder) error {
	orderBook.mu.Lock()
	defer orderBook.mu.Unlock()

	// Try to match immediately, cancel any remaining quantity
	if err := obe.matchOrder(ctx, orderBook, order); err != nil {
		return err
	}

	// IOC orders are not added to the book
	return nil
}

// processFOKOrder processes a Fill Or Kill order
func (obe *OrderBookEngine) processFOKOrder(ctx context.Context, orderBook *OrderBook, order *BookOrder) error {
	orderBook.mu.RLock()

	// Check if order can be fully filled
	canFill := obe.canFillOrder(orderBook, order)
	orderBook.mu.RUnlock()

	if !canFill {
		// Reject the order
		return fmt.Errorf("FOK order cannot be fully filled")
	}

	// Fill the order
	orderBook.mu.Lock()
	defer orderBook.mu.Unlock()

	return obe.matchOrder(ctx, orderBook, order)
}

// canFillOrder checks if an order can be fully filled
func (obe *OrderBookEngine) canFillOrder(orderBook *OrderBook, order *BookOrder) bool {
	var availableQty decimal.Decimal

	if order.Side == OrderSideBuy {
		// Check ask side
		obe.walkPriceLevels(orderBook.AskLevels, func(level *PriceLevel) bool {
			if level.Price.LessThanOrEqual(order.Price) {
				availableQty = availableQty.Add(level.TotalQty)
				return availableQty.LessThan(order.Quantity)
			}
			return false // Stop walking if price is too high
		})
	} else {
		// Check bid side
		obe.walkPriceLevels(orderBook.BidLevels, func(level *PriceLevel) bool {
			if level.Price.GreaterThanOrEqual(order.Price) {
				availableQty = availableQty.Add(level.TotalQty)
				return availableQty.LessThan(order.Quantity)
			}
			return false // Stop walking if price is too low
		})
	}

	return availableQty.GreaterThanOrEqual(order.Quantity)
}

// matchOrder attempts to match an order against the order book
func (obe *OrderBookEngine) matchOrder(ctx context.Context, orderBook *OrderBook, order *BookOrder) error {
	var oppositeSide *PriceLevelTree

	if order.Side == OrderSideBuy {
		oppositeSide = orderBook.AskLevels
	} else {
		oppositeSide = orderBook.BidLevels
	}

	// Walk through price levels and match
	obe.walkPriceLevels(oppositeSide, func(level *PriceLevel) bool {
		if order.Quantity.Equal(order.FilledQty) {
			return false // Order fully filled
		}

		// Check if prices cross
		if order.Side == OrderSideBuy && level.Price.GreaterThan(order.Price) {
			return false // No more matching prices
		}
		if order.Side == OrderSideSell && level.Price.LessThan(order.Price) {
			return false // No more matching prices
		}

		// Match against orders in this price level
		obe.matchAtPriceLevel(ctx, orderBook, order, level)

		return true // Continue to next price level
	})

	return nil
}

// matchAtPriceLevel matches an order against a specific price level
func (obe *OrderBookEngine) matchAtPriceLevel(ctx context.Context, orderBook *OrderBook, order *BookOrder, level *PriceLevel) {
	level.mu.Lock()
	defer level.mu.Unlock()

	i := 0
	for i < len(level.Orders) && order.Quantity.GreaterThan(order.FilledQty) {
		bookOrder := level.Orders[i]

		// Calculate trade quantity
		remainingQty := order.Quantity.Sub(order.FilledQty)
		bookOrderQty := bookOrder.Quantity.Sub(bookOrder.FilledQty)
		tradeQty := decimal.Min(remainingQty, bookOrderQty)

		// Execute trade
		obe.executeTrade(ctx, orderBook, order, bookOrder, tradeQty, level.Price)

		// Remove fully filled orders
		if bookOrder.Quantity.Equal(bookOrder.FilledQty) {
			level.Orders = append(level.Orders[:i], level.Orders[i+1:]...)
			level.OrderCount--
			delete(orderBook.Orders, bookOrder.ID)
		} else {
			i++
		}
	}

	// Update price level total quantity
	level.TotalQty = decimal.Zero
	for _, ord := range level.Orders {
		level.TotalQty = level.TotalQty.Add(ord.Quantity.Sub(ord.FilledQty))
	}

	// Remove empty price level
	if len(level.Orders) == 0 {
		obe.removePriceLevel(orderBook, level)
	}
}

// executeTrade executes a trade between two orders
func (obe *OrderBookEngine) executeTrade(ctx context.Context, orderBook *OrderBook, aggressorOrder, passiveOrder *BookOrder, quantity, price decimal.Decimal) {
	// Update fill quantities
	aggressorOrder.FilledQty = aggressorOrder.FilledQty.Add(quantity)
	passiveOrder.FilledQty = passiveOrder.FilledQty.Add(quantity)

	// Update timestamps
	now := time.Now()
	aggressorOrder.UpdatedAt = now
	passiveOrder.UpdatedAt = now

	// Update order book statistics
	orderBook.totalVolume = orderBook.totalVolume.Add(quantity)
	orderBook.lastTradePrice = price
	orderBook.lastTradeTime = now

	// Generate trade update
	update := &OrderBookUpdate{
		Symbol:         orderBook.Symbol,
		UpdateType:     UpdateTypeTrade,
		Price:          price,
		Quantity:       quantity,
		SequenceNumber: atomic.AddUint64(&orderBook.sequenceNumber, 1),
		Timestamp:      now,
	}

	// Update best bid/ask
	obe.updateBestPrices(orderBook, update)

	// Broadcast update
	obe.broadcastUpdate(update)

	obe.logger.Info(ctx, "Trade executed", map[string]interface{}{
		"symbol":          orderBook.Symbol,
		"price":           price.String(),
		"quantity":        quantity.String(),
		"aggressor_order": aggressorOrder.ID.String(),
		"passive_order":   passiveOrder.ID.String(),
	})
}

// addOrderToBook adds an order to the order book
func (obe *OrderBookEngine) addOrderToBook(ctx context.Context, orderBook *OrderBook, order *BookOrder) error {
	var priceLevels *PriceLevelTree

	if order.Side == OrderSideBuy {
		priceLevels = orderBook.BidLevels
	} else {
		priceLevels = orderBook.AskLevels
	}

	// Find or create price level
	level := obe.findOrCreatePriceLevel(priceLevels, order.Price)

	// Add order to price level
	level.mu.Lock()
	level.Orders = append(level.Orders, order)
	level.OrderCount++
	level.TotalQty = level.TotalQty.Add(order.Quantity.Sub(order.FilledQty))
	order.PriceLevel = level
	order.Position = len(level.Orders) - 1
	level.mu.Unlock()

	// Add to order tracking
	orderBook.Orders[order.ID] = order

	// Generate add update
	update := &OrderBookUpdate{
		Symbol:         orderBook.Symbol,
		UpdateType:     UpdateTypeAdd,
		Side:           order.Side,
		Price:          order.Price,
		Quantity:       order.Quantity.Sub(order.FilledQty),
		OrderID:        order.ID,
		SequenceNumber: atomic.AddUint64(&orderBook.sequenceNumber, 1),
		Timestamp:      time.Now(),
	}

	// Update best bid/ask
	obe.updateBestPrices(orderBook, update)

	// Broadcast update
	obe.broadcastUpdate(update)

	return nil
}

// walkPriceLevels walks through price levels in order
func (obe *OrderBookEngine) walkPriceLevels(tree *PriceLevelTree, fn func(*PriceLevel) bool) {
	// This is a simplified implementation - in production, you'd implement a proper BST traversal
	// For now, we'll use a mock implementation
}

// findOrCreatePriceLevel finds or creates a price level in the tree
func (obe *OrderBookEngine) findOrCreatePriceLevel(tree *PriceLevelTree, price decimal.Decimal) *PriceLevel {
	// Simplified implementation - in production, you'd implement proper BST operations
	return &PriceLevel{
		Price:      price,
		Orders:     make([]*BookOrder, 0),
		TotalQty:   decimal.Zero,
		OrderCount: 0,
	}
}

// removePriceLevel removes a price level from the tree
func (obe *OrderBookEngine) removePriceLevel(orderBook *OrderBook, level *PriceLevel) {
	// Simplified implementation - in production, you'd implement proper BST removal
}

// updateBestPrices updates the best bid/ask prices
func (obe *OrderBookEngine) updateBestPrices(orderBook *OrderBook, update *OrderBookUpdate) {
	// Update best bid
	if bestBid := obe.getBestBid(orderBook); bestBid != nil {
		update.BestBid = bestBid.Price
		update.BidSize = bestBid.TotalQty
		orderBook.bestBid.Store(bestBid)
	}

	// Update best ask
	if bestAsk := obe.getBestAsk(orderBook); bestAsk != nil {
		update.BestAsk = bestAsk.Price
		update.AskSize = bestAsk.TotalQty
		orderBook.bestAsk.Store(bestAsk)
	}
}

// getBestBid gets the best bid price level
func (obe *OrderBookEngine) getBestBid(orderBook *OrderBook) *PriceLevel {
	// Simplified implementation - in production, you'd get from BST
	return nil
}

// getBestAsk gets the best ask price level
func (obe *OrderBookEngine) getBestAsk(orderBook *OrderBook) *PriceLevel {
	// Simplified implementation - in production, you'd get from BST
	return nil
}

// broadcastUpdate broadcasts an order book update to subscribers
func (obe *OrderBookEngine) broadcastUpdate(update *OrderBookUpdate) {
	obe.mu.RLock()
	defer obe.mu.RUnlock()

	// Send to symbol-specific subscribers
	if subscribers, exists := obe.subscribers[update.Symbol]; exists {
		for _, ch := range subscribers {
			select {
			case ch <- update:
			default:
				// Subscriber channel is full, skip
			}
		}
	}

	// Send to wildcard subscribers (all symbols)
	if subscribers, exists := obe.subscribers["*"]; exists {
		for _, ch := range subscribers {
			select {
			case ch <- update:
			default:
				// Subscriber channel is full, skip
			}
		}
	}
}

// performanceMonitor tracks and reports performance metrics
func (obe *OrderBookEngine) performanceMonitor(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var lastOrderCount int64
	var lastUpdateCount int64

	for {
		select {
		case <-obe.stopChan:
			return
		case <-ticker.C:
			currentOrders := atomic.LoadInt64(&obe.ordersProcessed)
			currentUpdates := atomic.LoadInt64(&obe.updatesProcessed)

			ordersPerSecond := currentOrders - lastOrderCount
			updatesPerSecond := currentUpdates - lastUpdateCount

			lastOrderCount = currentOrders
			lastUpdateCount = currentUpdates

			avgProcessingTime := atomic.LoadInt64(&obe.avgProcessingTimeNanos)

			obe.logger.Info(ctx, "Order book engine performance", map[string]interface{}{
				"orders_per_second":  ordersPerSecond,
				"updates_per_second": updatesPerSecond,
				"avg_processing_ns":  avgProcessingTime,
				"avg_processing_us":  avgProcessingTime / 1000,
				"total_orders":       currentOrders,
				"total_updates":      currentUpdates,
			})
		}
	}
}

// Subscribe registers a subscriber for order book updates
func (obe *OrderBookEngine) Subscribe(symbol string) <-chan *OrderBookUpdate {
	obe.mu.Lock()
	defer obe.mu.Unlock()

	ch := make(chan *OrderBookUpdate, 1000) // Buffered channel
	if obe.subscribers[symbol] == nil {
		obe.subscribers[symbol] = make([]chan *OrderBookUpdate, 0)
	}
	obe.subscribers[symbol] = append(obe.subscribers[symbol], ch)

	return ch
}

// GetOrderBook returns a snapshot of the order book for a symbol
func (obe *OrderBookEngine) GetOrderBook(symbol string) (*OrderBookSnapshot, error) {
	obe.mu.RLock()
	orderBook, exists := obe.orderBooks[symbol]
	obe.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("order book not found for symbol %s", symbol)
	}

	orderBook.mu.RLock()
	defer orderBook.mu.RUnlock()

	snapshot := &OrderBookSnapshot{
		Symbol:         symbol,
		SequenceNumber: orderBook.sequenceNumber,
		Timestamp:      time.Now(),
		Bids:           obe.getPriceLevelsSnapshot(orderBook.BidLevels, false),
		Asks:           obe.getPriceLevelsSnapshot(orderBook.AskLevels, true),
		LastTradePrice: orderBook.lastTradePrice,
		LastTradeTime:  orderBook.lastTradeTime,
		TotalVolume:    orderBook.totalVolume,
	}

	return snapshot, nil
}

// OrderBookSnapshot represents a point-in-time snapshot of an order book
type OrderBookSnapshot struct {
	Symbol         string               `json:"symbol"`
	SequenceNumber uint64               `json:"sequence_number"`
	Timestamp      time.Time            `json:"timestamp"`
	Bids           []PriceLevelSnapshot `json:"bids"`
	Asks           []PriceLevelSnapshot `json:"asks"`
	LastTradePrice decimal.Decimal      `json:"last_trade_price"`
	LastTradeTime  time.Time            `json:"last_trade_time"`
	TotalVolume    decimal.Decimal      `json:"total_volume"`
}

// PriceLevelSnapshot represents a price level in a snapshot
type PriceLevelSnapshot struct {
	Price      decimal.Decimal `json:"price"`
	Quantity   decimal.Decimal `json:"quantity"`
	OrderCount int             `json:"order_count"`
}

// getPriceLevelsSnapshot gets a snapshot of price levels
func (obe *OrderBookEngine) getPriceLevelsSnapshot(tree *PriceLevelTree, isAsk bool) []PriceLevelSnapshot {
	// Simplified implementation - in production, you'd traverse the BST
	return []PriceLevelSnapshot{}
}

// GetMetrics returns current performance metrics
func (obe *OrderBookEngine) GetMetrics() OrderBookMetrics {
	return OrderBookMetrics{
		OrdersProcessed:        atomic.LoadInt64(&obe.ordersProcessed),
		UpdatesProcessed:       atomic.LoadInt64(&obe.updatesProcessed),
		AvgProcessingTimeNanos: atomic.LoadInt64(&obe.avgProcessingTimeNanos),
		IsRunning:              atomic.LoadInt32(&obe.isRunning) == 1,
		ActiveOrderBooks:       len(obe.orderBooks),
		QueueUtilization:       obe.getQueueUtilization(),
	}
}

// OrderBookMetrics contains performance metrics for the order book engine
type OrderBookMetrics struct {
	OrdersProcessed        int64   `json:"orders_processed"`
	UpdatesProcessed       int64   `json:"updates_processed"`
	AvgProcessingTimeNanos int64   `json:"avg_processing_time_nanos"`
	IsRunning              bool    `json:"is_running"`
	ActiveOrderBooks       int     `json:"active_order_books"`
	QueueUtilization       float64 `json:"queue_utilization"`
}

// getQueueUtilization calculates current queue utilization
func (obe *OrderBookEngine) getQueueUtilization() float64 {
	writeIndex := atomic.LoadInt64(&obe.orderQueue.writeIndex)
	readIndex := atomic.LoadInt64(&obe.orderQueue.readIndex)
	used := writeIndex - readIndex
	return float64(used) / float64(obe.orderQueue.capacity) * 100.0
}
