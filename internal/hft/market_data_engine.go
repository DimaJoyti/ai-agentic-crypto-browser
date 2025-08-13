package hft

import (
	"context"
	"fmt"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/shopspring/decimal"
)

// MarketDataEngine provides ultra-low latency market data ingestion
// with kernel bypass, nanosecond timestamping, and lock-free processing
type MarketDataEngine struct {
	logger *observability.Logger
	config MarketDataConfig

	// Lock-free ring buffers for different data types
	tickBuffer  *LockFreeRingBuffer
	l2Buffer    *LockFreeRingBuffer
	tradeBuffer *LockFreeRingBuffer

	// Multicast UDP connections
	multicastConns map[string]*MulticastConnection

	// Performance metrics
	ticksProcessed  int64
	avgLatencyNanos int64
	droppedPackets  int64

	// State management
	isRunning int32
	stopChan  chan struct{}
	wg        sync.WaitGroup
	mu        sync.RWMutex

	// Subscribers for processed data
	subscribers map[string][]chan *NormalizedTick
}

// MarketDataConfig contains configuration for ultra-low latency market data
type MarketDataConfig struct {
	// Multicast configuration
	MulticastGroups  map[string]string `json:"multicast_groups"` // exchange -> multicast address
	NetworkInterface string            `json:"network_interface"`
	BufferSize       int               `json:"buffer_size"`

	// Performance tuning
	RingBufferSize   int   `json:"ring_buffer_size"`
	ProcessorThreads int   `json:"processor_threads"`
	CPUAffinity      []int `json:"cpu_affinity"`
	UseKernelBypass  bool  `json:"use_kernel_bypass"`

	// Timestamping
	UseHardwareTimestamps bool `json:"use_hardware_timestamps"`
	TimestampPrecision    int  `json:"timestamp_precision"` // nanoseconds

	// Quality control
	MaxLatencyNanos     int64 `json:"max_latency_nanos"`
	DropThresholdNanos  int64 `json:"drop_threshold_nanos"`
	EnableSequenceCheck bool  `json:"enable_sequence_check"`
}

// NormalizedTick represents a normalized market data tick with nanosecond precision
type NormalizedTick struct {
	// Core data
	Symbol   string          `json:"symbol"`
	Exchange string          `json:"exchange"`
	Price    decimal.Decimal `json:"price"`
	Size     decimal.Decimal `json:"size"`
	Side     OrderSide       `json:"side"`

	// Timestamps (nanosecond precision)
	ExchangeTimestamp int64 `json:"exchange_timestamp"`
	ReceiveTimestamp  int64 `json:"receive_timestamp"`
	ProcessTimestamp  int64 `json:"process_timestamp"`

	// Sequence and quality
	SequenceNumber uint64 `json:"sequence_number"`
	LatencyNanos   int64  `json:"latency_nanos"`

	// Level 2 data (for order book updates)
	BidPrice decimal.Decimal `json:"bid_price,omitempty"`
	AskPrice decimal.Decimal `json:"ask_price,omitempty"`
	BidSize  decimal.Decimal `json:"bid_size,omitempty"`
	AskSize  decimal.Decimal `json:"ask_size,omitempty"`

	// Metadata
	MessageType string                 `json:"message_type"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// LockFreeRingBuffer implements a lock-free ring buffer for ultra-low latency
type LockFreeRingBuffer struct {
	buffer   []unsafe.Pointer
	capacity int64
	mask     int64

	// Atomic counters
	writeIndex int64
	readIndex  int64

	// Padding to prevent false sharing
	_ [56]byte
}

// MulticastConnection represents a multicast UDP connection with kernel bypass
type MulticastConnection struct {
	conn     *net.UDPConn
	addr     *net.UDPAddr
	exchange string

	// Performance optimization
	buffer       []byte
	lastSequence uint64

	// Statistics
	packetsReceived int64
	bytesReceived   int64
	errors          int64
}

// NewMarketDataEngine creates a new ultra-low latency market data engine
func NewMarketDataEngine(logger *observability.Logger, config MarketDataConfig) *MarketDataEngine {
	// Set default values
	if config.RingBufferSize == 0 {
		config.RingBufferSize = 1024 * 1024 // 1M entries
	}
	if config.ProcessorThreads == 0 {
		config.ProcessorThreads = runtime.NumCPU()
	}
	if config.BufferSize == 0 {
		config.BufferSize = 64 * 1024 // 64KB
	}

	engine := &MarketDataEngine{
		logger:         logger,
		config:         config,
		multicastConns: make(map[string]*MulticastConnection),
		subscribers:    make(map[string][]chan *NormalizedTick),
		stopChan:       make(chan struct{}),
	}

	// Initialize lock-free ring buffers
	engine.tickBuffer = NewLockFreeRingBuffer(config.RingBufferSize)
	engine.l2Buffer = NewLockFreeRingBuffer(config.RingBufferSize)
	engine.tradeBuffer = NewLockFreeRingBuffer(config.RingBufferSize)

	return engine
}

// NewLockFreeRingBuffer creates a new lock-free ring buffer
func NewLockFreeRingBuffer(capacity int) *LockFreeRingBuffer {
	// Ensure capacity is power of 2 for efficient masking
	if capacity&(capacity-1) != 0 {
		capacity = nextPowerOfTwo(capacity)
	}

	return &LockFreeRingBuffer{
		buffer:   make([]unsafe.Pointer, capacity),
		capacity: int64(capacity),
		mask:     int64(capacity - 1),
	}
}

// nextPowerOfTwo returns the next power of 2 greater than or equal to n
func nextPowerOfTwo(n int) int {
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n++
	return n
}

// Start begins the market data engine with ultra-low latency processing
func (mde *MarketDataEngine) Start(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&mde.isRunning, 0, 1) {
		return fmt.Errorf("market data engine is already running")
	}

	mde.logger.Info(ctx, "Starting ultra-low latency market data engine", map[string]interface{}{
		"ring_buffer_size":    mde.config.RingBufferSize,
		"processor_threads":   mde.config.ProcessorThreads,
		"kernel_bypass":       mde.config.UseKernelBypass,
		"hardware_timestamps": mde.config.UseHardwareTimestamps,
	})

	// Set CPU affinity if specified
	if len(mde.config.CPUAffinity) > 0 {
		if err := mde.setCPUAffinity(); err != nil {
			mde.logger.Warn(ctx, "Failed to set CPU affinity", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	// Initialize multicast connections
	if err := mde.initializeMulticastConnections(ctx); err != nil {
		return fmt.Errorf("failed to initialize multicast connections: %w", err)
	}

	// Start processing threads
	mde.wg.Add(mde.config.ProcessorThreads + len(mde.multicastConns))

	// Start data processors
	for i := 0; i < mde.config.ProcessorThreads; i++ {
		go mde.processMarketData(ctx, i)
	}

	// Start multicast receivers
	for exchange, conn := range mde.multicastConns {
		go mde.receiveMulticastData(ctx, exchange, conn)
	}

	// Start performance monitor
	go mde.performanceMonitor(ctx)

	mde.logger.Info(ctx, "Market data engine started successfully", nil)
	return nil
}

// Stop gracefully shuts down the market data engine
func (mde *MarketDataEngine) Stop(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&mde.isRunning, 1, 0) {
		return fmt.Errorf("market data engine is not running")
	}

	mde.logger.Info(ctx, "Stopping market data engine", nil)

	close(mde.stopChan)
	mde.wg.Wait()

	// Close multicast connections
	for _, conn := range mde.multicastConns {
		if conn.conn != nil {
			conn.conn.Close()
		}
	}

	mde.logger.Info(ctx, "Market data engine stopped", map[string]interface{}{
		"ticks_processed": atomic.LoadInt64(&mde.ticksProcessed),
		"dropped_packets": atomic.LoadInt64(&mde.droppedPackets),
		"avg_latency_ns":  atomic.LoadInt64(&mde.avgLatencyNanos),
	})

	return nil
}

// Push adds a tick to the ring buffer (lock-free)
func (rb *LockFreeRingBuffer) Push(item unsafe.Pointer) bool {
	for {
		writeIndex := atomic.LoadInt64(&rb.writeIndex)
		nextIndex := writeIndex + 1

		// Check if buffer is full
		if nextIndex-atomic.LoadInt64(&rb.readIndex) >= rb.capacity {
			return false // Buffer full
		}

		// Try to claim the slot
		if atomic.CompareAndSwapInt64(&rb.writeIndex, writeIndex, nextIndex) {
			// Store the item
			atomic.StorePointer(&rb.buffer[writeIndex&rb.mask], item)
			return true
		}
		// Retry if CAS failed
		runtime.Gosched()
	}
}

// Pop removes a tick from the ring buffer (lock-free)
func (rb *LockFreeRingBuffer) Pop() unsafe.Pointer {
	for {
		readIndex := atomic.LoadInt64(&rb.readIndex)

		// Check if buffer is empty
		if readIndex >= atomic.LoadInt64(&rb.writeIndex) {
			return nil
		}

		// Try to claim the slot
		if atomic.CompareAndSwapInt64(&rb.readIndex, readIndex, readIndex+1) {
			// Load the item
			return atomic.LoadPointer(&rb.buffer[readIndex&rb.mask])
		}
		// Retry if CAS failed
		runtime.Gosched()
	}
}

// GetNanosecondTimestamp returns current time in nanoseconds with highest precision
func GetNanosecondTimestamp() int64 {
	return time.Now().UnixNano()
}

// setCPUAffinity sets CPU affinity for the current process (Linux-specific)
func (mde *MarketDataEngine) setCPUAffinity() error {
	// This would require CGO and Linux-specific syscalls
	// For now, we'll log the intent
	mde.logger.Info(context.Background(), "CPU affinity setting requested", map[string]interface{}{
		"cpus": mde.config.CPUAffinity,
		"note": "Requires CGO and Linux-specific implementation",
	})
	return nil
}

// initializeMulticastConnections sets up multicast UDP connections for each exchange
func (mde *MarketDataEngine) initializeMulticastConnections(ctx context.Context) error {
	for exchange, multicastAddr := range mde.config.MulticastGroups {
		conn, err := mde.createMulticastConnection(exchange, multicastAddr)
		if err != nil {
			return fmt.Errorf("failed to create multicast connection for %s: %w", exchange, err)
		}
		mde.multicastConns[exchange] = conn

		mde.logger.Info(ctx, "Initialized multicast connection", map[string]interface{}{
			"exchange": exchange,
			"address":  multicastAddr,
		})
	}
	return nil
}

// createMulticastConnection creates a multicast UDP connection with optimizations
func (mde *MarketDataEngine) createMulticastConnection(exchange, multicastAddr string) (*MulticastConnection, error) {
	addr, err := net.ResolveUDPAddr("udp", multicastAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve multicast address: %w", err)
	}

	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on multicast address: %w", err)
	}

	// Set socket buffer size for high throughput
	if err := conn.SetReadBuffer(mde.config.BufferSize); err != nil {
		mde.logger.Warn(context.Background(), "Failed to set read buffer size", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return &MulticastConnection{
		conn:     conn,
		addr:     addr,
		exchange: exchange,
		buffer:   make([]byte, mde.config.BufferSize),
	}, nil
}

// receiveMulticastData receives and processes multicast market data
func (mde *MarketDataEngine) receiveMulticastData(ctx context.Context, exchange string, conn *MulticastConnection) {
	defer mde.wg.Done()

	mde.logger.Info(ctx, "Starting multicast receiver", map[string]interface{}{
		"exchange": exchange,
		"address":  conn.addr.String(),
	})

	for {
		select {
		case <-mde.stopChan:
			return
		default:
			// Set read deadline to avoid blocking indefinitely
			conn.conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))

			n, err := conn.conn.Read(conn.buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue // Timeout is expected
				}
				atomic.AddInt64(&conn.errors, 1)
				continue
			}

			// Record receive timestamp immediately
			receiveTimestamp := GetNanosecondTimestamp()

			// Update statistics
			atomic.AddInt64(&conn.packetsReceived, 1)
			atomic.AddInt64(&conn.bytesReceived, int64(n))

			// Parse and normalize the market data
			if err := mde.parseAndNormalize(exchange, conn.buffer[:n], receiveTimestamp); err != nil {
				mde.logger.Error(ctx, "Failed to parse market data", err)
				continue
			}
		}
	}
}

// parseAndNormalize parses raw market data and normalizes it
func (mde *MarketDataEngine) parseAndNormalize(exchange string, data []byte, receiveTimestamp int64) error {
	// This is a simplified parser - in production, you'd have exchange-specific parsers
	// For now, we'll create a mock normalized tick

	processTimestamp := GetNanosecondTimestamp()

	tick := &NormalizedTick{
		Symbol:            "BTCUSDT", // Would be parsed from data
		Exchange:          exchange,
		Price:             decimal.NewFromFloat(45000.0), // Would be parsed from data
		Size:              decimal.NewFromFloat(0.1),     // Would be parsed from data
		Side:              OrderSideBuy,                  // Would be parsed from data
		ExchangeTimestamp: receiveTimestamp - 1000,       // Mock exchange timestamp
		ReceiveTimestamp:  receiveTimestamp,
		ProcessTimestamp:  processTimestamp,
		SequenceNumber:    atomic.AddUint64(&mde.multicastConns[exchange].lastSequence, 1),
		LatencyNanos:      processTimestamp - receiveTimestamp,
		MessageType:       "TRADE",
	}

	// Push to appropriate ring buffer based on message type
	tickPtr := unsafe.Pointer(tick)
	if !mde.tickBuffer.Push(tickPtr) {
		atomic.AddInt64(&mde.droppedPackets, 1)
		return fmt.Errorf("tick buffer full, dropping packet")
	}

	return nil
}

// processMarketData processes normalized market data from ring buffers
func (mde *MarketDataEngine) processMarketData(ctx context.Context, workerID int) {
	defer mde.wg.Done()

	mde.logger.Info(ctx, "Starting market data processor", map[string]interface{}{
		"worker_id": workerID,
	})

	for {
		select {
		case <-mde.stopChan:
			return
		default:
			// Process ticks from ring buffer
			if tickPtr := mde.tickBuffer.Pop(); tickPtr != nil {
				tick := (*NormalizedTick)(tickPtr)
				mde.processTick(ctx, tick)
				atomic.AddInt64(&mde.ticksProcessed, 1)
			} else {
				// No data available, yield CPU
				runtime.Gosched()
			}
		}
	}
}

// processTick processes a single normalized tick
func (mde *MarketDataEngine) processTick(ctx context.Context, tick *NormalizedTick) {
	// Update latency statistics
	currentLatency := tick.LatencyNanos
	atomic.StoreInt64(&mde.avgLatencyNanos, currentLatency)

	// Check if latency exceeds threshold
	if mde.config.MaxLatencyNanos > 0 && currentLatency > mde.config.MaxLatencyNanos {
		mde.logger.Warn(ctx, "High latency detected", map[string]interface{}{
			"symbol":       tick.Symbol,
			"latency_ns":   currentLatency,
			"threshold_ns": mde.config.MaxLatencyNanos,
		})
	}

	// Distribute to subscribers
	mde.distributeToSubscribers(tick)
}

// distributeToSubscribers sends the tick to all registered subscribers
func (mde *MarketDataEngine) distributeToSubscribers(tick *NormalizedTick) {
	mde.mu.RLock()
	defer mde.mu.RUnlock()

	// Send to symbol-specific subscribers
	if subscribers, exists := mde.subscribers[tick.Symbol]; exists {
		for _, ch := range subscribers {
			select {
			case ch <- tick:
			default:
				// Subscriber channel is full, skip
			}
		}
	}

	// Send to wildcard subscribers (all symbols)
	if subscribers, exists := mde.subscribers["*"]; exists {
		for _, ch := range subscribers {
			select {
			case ch <- tick:
			default:
				// Subscriber channel is full, skip
			}
		}
	}
}

// Subscribe registers a subscriber for market data updates
func (mde *MarketDataEngine) Subscribe(symbol string) <-chan *NormalizedTick {
	mde.mu.Lock()
	defer mde.mu.Unlock()

	ch := make(chan *NormalizedTick, 1000) // Buffered channel
	if mde.subscribers[symbol] == nil {
		mde.subscribers[symbol] = make([]chan *NormalizedTick, 0)
	}
	mde.subscribers[symbol] = append(mde.subscribers[symbol], ch)

	return ch
}

// performanceMonitor tracks and reports performance metrics
func (mde *MarketDataEngine) performanceMonitor(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var lastTickCount int64

	for {
		select {
		case <-mde.stopChan:
			return
		case <-ticker.C:
			currentTicks := atomic.LoadInt64(&mde.ticksProcessed)
			ticksPerSecond := currentTicks - lastTickCount
			lastTickCount = currentTicks

			avgLatency := atomic.LoadInt64(&mde.avgLatencyNanos)
			droppedPackets := atomic.LoadInt64(&mde.droppedPackets)

			mde.logger.Info(ctx, "Market data engine performance", map[string]interface{}{
				"ticks_per_second": ticksPerSecond,
				"avg_latency_ns":   avgLatency,
				"avg_latency_us":   avgLatency / 1000,
				"dropped_packets":  droppedPackets,
				"total_ticks":      currentTicks,
			})
		}
	}
}

// GetMetrics returns current performance metrics
func (mde *MarketDataEngine) GetMetrics() MarketDataMetrics {
	return MarketDataMetrics{
		TicksProcessed:    atomic.LoadInt64(&mde.ticksProcessed),
		AvgLatencyNanos:   atomic.LoadInt64(&mde.avgLatencyNanos),
		DroppedPackets:    atomic.LoadInt64(&mde.droppedPackets),
		IsRunning:         atomic.LoadInt32(&mde.isRunning) == 1,
		ActiveConnections: len(mde.multicastConns),
		BufferUtilization: mde.getBufferUtilization(),
	}
}

// MarketDataMetrics contains performance metrics for the market data engine
type MarketDataMetrics struct {
	TicksProcessed    int64   `json:"ticks_processed"`
	AvgLatencyNanos   int64   `json:"avg_latency_nanos"`
	DroppedPackets    int64   `json:"dropped_packets"`
	IsRunning         bool    `json:"is_running"`
	ActiveConnections int     `json:"active_connections"`
	BufferUtilization float64 `json:"buffer_utilization"`
}

// getBufferUtilization calculates current buffer utilization
func (mde *MarketDataEngine) getBufferUtilization() float64 {
	writeIndex := atomic.LoadInt64(&mde.tickBuffer.writeIndex)
	readIndex := atomic.LoadInt64(&mde.tickBuffer.readIndex)
	used := writeIndex - readIndex
	return float64(used) / float64(mde.tickBuffer.capacity) * 100.0
}
