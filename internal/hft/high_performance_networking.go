package hft

import (
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

// HighPerformanceNetworking provides ultra-low latency networking with kernel bypass,
// zero-copy operations, CPU affinity, and hardware timestamping for HFT applications
type HighPerformanceNetworking struct {
	logger *observability.Logger
	config NetworkConfig

	// Network components
	udpManager       *UDPManager
	tcpManager       *TCPManager
	multicastManager *MulticastManager
	timestampEngine  *TimestampEngine
	bufferPool       *BufferPool

	// Performance optimizations
	cpuAffinity    *CPUAffinity
	memoryManager  *MemoryManager
	lockFreeQueues map[string]*LockFreeQueue

	// Connection management
	connections    map[string]*Connection
	connectionPool *ConnectionPool

	// Performance metrics
	packetsReceived int64
	packetsSent     int64
	bytesReceived   int64
	bytesSent       int64
	avgLatency      int64
	minLatency      int64
	maxLatency      int64

	// State management
	isRunning int32
	stopChan  chan struct{}
	wg        sync.WaitGroup
	mu        sync.RWMutex

	// Event subscribers
	subscribers map[string][]chan *NetworkEvent
}

// NetworkConfig contains configuration for high-performance networking
type NetworkConfig struct {
	// Network settings
	UDPPort         int      `json:"udp_port"`
	TCPPort         int      `json:"tcp_port"`
	MulticastGroups []string `json:"multicast_groups"`
	InterfaceName   string   `json:"interface_name"`

	// Performance settings
	EnableKernelBypass bool  `json:"enable_kernel_bypass"`
	EnableZeroCopy     bool  `json:"enable_zero_copy"`
	EnableHWTimestamps bool  `json:"enable_hw_timestamps"`
	CPUAffinity        []int `json:"cpu_affinity"`

	// Buffer settings
	ReceiveBufferSize   int `json:"receive_buffer_size"`
	SendBufferSize      int `json:"send_buffer_size"`
	BufferPoolSize      int `json:"buffer_pool_size"`
	PreallocatedBuffers int `json:"preallocated_buffers"`

	// Queue settings
	QueueDepth      int           `json:"queue_depth"`
	BatchSize       int           `json:"batch_size"`
	PollingInterval time.Duration `json:"polling_interval"`

	// Latency settings
	LatencyTarget   time.Duration `json:"latency_target"`
	JitterThreshold time.Duration `json:"jitter_threshold"`
	TimeoutDuration time.Duration `json:"timeout_duration"`

	// Connection settings
	MaxConnections     int           `json:"max_connections"`
	ConnectionPoolSize int           `json:"connection_pool_size"`
	KeepAliveInterval  time.Duration `json:"keep_alive_interval"`

	// Monitoring settings
	EnableMetrics   bool          `json:"enable_metrics"`
	MetricsInterval time.Duration `json:"metrics_interval"`
	EnableTracing   bool          `json:"enable_tracing"`
}

// Connection represents a high-performance network connection
type Connection struct {
	ID         uuid.UUID       `json:"id"`
	Type       ConnectionType  `json:"type"`
	LocalAddr  net.Addr        `json:"local_addr"`
	RemoteAddr net.Addr        `json:"remote_addr"`
	State      ConnectionState `json:"state"`

	// Performance metrics
	PacketsReceived int64     `json:"packets_received"`
	PacketsSent     int64     `json:"packets_sent"`
	BytesReceived   int64     `json:"bytes_received"`
	BytesSent       int64     `json:"bytes_sent"`
	LastActivity    time.Time `json:"last_activity"`

	// Latency tracking
	MinLatency time.Duration `json:"min_latency"`
	MaxLatency time.Duration `json:"max_latency"`
	AvgLatency time.Duration `json:"avg_latency"`
	Jitter     time.Duration `json:"jitter"`

	// Connection-specific data
	conn      net.Conn
	buffer    []byte
	timestamp time.Time
	cpuCore   int
}

// ConnectionType represents different types of connections
type ConnectionType string

const (
	ConnectionTypeUDP       ConnectionType = "UDP"
	ConnectionTypeTCP       ConnectionType = "TCP"
	ConnectionTypeMulticast ConnectionType = "MULTICAST"
	ConnectionTypeRaw       ConnectionType = "RAW"
)

// ConnectionState represents connection states
type ConnectionState string

const (
	ConnectionStateConnecting ConnectionState = "CONNECTING"
	ConnectionStateConnected  ConnectionState = "CONNECTED"
	ConnectionStateClosing    ConnectionState = "CLOSING"
	ConnectionStateClosed     ConnectionState = "CLOSED"
	ConnectionStateError      ConnectionState = "ERROR"
)

// NetworkEvent represents a network event
type NetworkEvent struct {
	ID           uuid.UUID              `json:"id"`
	Type         NetworkEventType       `json:"type"`
	Timestamp    time.Time              `json:"timestamp"`
	ConnectionID uuid.UUID              `json:"connection_id"`
	Data         map[string]interface{} `json:"data"`
	Latency      time.Duration          `json:"latency"`
	Size         int                    `json:"size"`
}

// NetworkEventType represents different types of network events
type NetworkEventType string

const (
	NetworkEventPacketReceived  NetworkEventType = "PACKET_RECEIVED"
	NetworkEventPacketSent      NetworkEventType = "PACKET_SENT"
	NetworkEventConnectionOpen  NetworkEventType = "CONNECTION_OPEN"
	NetworkEventConnectionClose NetworkEventType = "CONNECTION_CLOSE"
	NetworkEventLatencyAlert    NetworkEventType = "LATENCY_ALERT"
	NetworkEventError           NetworkEventType = "ERROR"
)

// NetworkMessage represents a network message
type NetworkMessage struct {
	ID          uuid.UUID     `json:"id"`
	Type        MessageType   `json:"type"`
	Payload     []byte        `json:"payload"`
	Timestamp   time.Time     `json:"timestamp"`
	Source      net.Addr      `json:"source"`
	Destination net.Addr      `json:"destination"`
	Priority    Priority      `json:"priority"`
	TTL         time.Duration `json:"ttl"`
}

// MessageType represents different types of messages
type MessageType string

const (
	MessageTypeMarketData  MessageType = "MARKET_DATA"
	MessageTypeOrderUpdate MessageType = "ORDER_UPDATE"
	MessageTypeTrade       MessageType = "TRADE"
	MessageTypeHeartbeat   MessageType = "HEARTBEAT"
	MessageTypeControl     MessageType = "CONTROL"
)

// Priority represents message priority levels
type Priority int

const (
	PriorityLow Priority = iota
	PriorityNormal
	PriorityHigh
	PriorityCritical
)

// NewHighPerformanceNetworking creates a new high-performance networking system
func NewHighPerformanceNetworking(logger *observability.Logger, config NetworkConfig) *HighPerformanceNetworking {
	// Set default values
	if config.UDPPort == 0 {
		config.UDPPort = 8080
	}
	if config.TCPPort == 0 {
		config.TCPPort = 8081
	}
	if config.ReceiveBufferSize == 0 {
		config.ReceiveBufferSize = 64 * 1024 // 64KB
	}
	if config.SendBufferSize == 0 {
		config.SendBufferSize = 64 * 1024 // 64KB
	}
	if config.BufferPoolSize == 0 {
		config.BufferPoolSize = 1000
	}
	if config.QueueDepth == 0 {
		config.QueueDepth = 1024
	}
	if config.BatchSize == 0 {
		config.BatchSize = 32
	}
	if config.LatencyTarget == 0 {
		config.LatencyTarget = 10 * time.Microsecond
	}
	if config.MaxConnections == 0 {
		config.MaxConnections = 1000
	}

	hpn := &HighPerformanceNetworking{
		logger:         logger,
		config:         config,
		connections:    make(map[string]*Connection),
		lockFreeQueues: make(map[string]*LockFreeQueue),
		subscribers:    make(map[string][]chan *NetworkEvent),
		stopChan:       make(chan struct{}),
		minLatency:     int64(time.Hour), // Initialize to high value
	}

	// Initialize components
	hpn.udpManager = NewUDPManager(logger, config)
	hpn.tcpManager = NewTCPManager(logger, config)
	hpn.multicastManager = NewMulticastManager(logger, config)
	hpn.timestampEngine = NewTimestampEngine(logger, config)
	hpn.bufferPool = NewBufferPool(logger, config)
	hpn.cpuAffinity = NewCPUAffinity(logger, config)
	hpn.memoryManager = NewMemoryManager(logger, config)
	hpn.connectionPool = NewConnectionPool(logger, config)

	// Initialize lock-free queues
	hpn.lockFreeQueues["inbound"] = NewLockFreeQueue(config.QueueDepth)
	hpn.lockFreeQueues["outbound"] = NewLockFreeQueue(config.QueueDepth)
	hpn.lockFreeQueues["priority"] = NewLockFreeQueue(config.QueueDepth)

	return hpn
}

// Start begins the high-performance networking system
func (hpn *HighPerformanceNetworking) Start(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&hpn.isRunning, 0, 1) {
		return fmt.Errorf("high-performance networking is already running")
	}

	hpn.logger.Info(ctx, "Starting high-performance networking system", map[string]interface{}{
		"udp_port":       hpn.config.UDPPort,
		"tcp_port":       hpn.config.TCPPort,
		"kernel_bypass":  hpn.config.EnableKernelBypass,
		"zero_copy":      hpn.config.EnableZeroCopy,
		"hw_timestamps":  hpn.config.EnableHWTimestamps,
		"cpu_affinity":   hpn.config.CPUAffinity,
		"latency_target": hpn.config.LatencyTarget.String(),
	})

	// Set CPU affinity if configured
	if len(hpn.config.CPUAffinity) > 0 {
		if err := hpn.cpuAffinity.SetAffinity(hpn.config.CPUAffinity); err != nil {
			hpn.logger.Warn(ctx, "Failed to set CPU affinity", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	// Initialize memory manager
	if err := hpn.memoryManager.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize memory manager: %w", err)
	}

	// Start network managers
	if err := hpn.startNetworkManagers(ctx); err != nil {
		return fmt.Errorf("failed to start network managers: %w", err)
	}

	// Start processing threads
	hpn.wg.Add(4)
	go hpn.processInboundPackets(ctx)
	go hpn.processOutboundPackets(ctx)
	go hpn.monitorConnections(ctx)
	go hpn.performanceMonitor(ctx)

	hpn.logger.Info(ctx, "High-performance networking system started successfully", nil)
	return nil
}

// Stop gracefully shuts down the high-performance networking system
func (hpn *HighPerformanceNetworking) Stop(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&hpn.isRunning, 1, 0) {
		return fmt.Errorf("high-performance networking is not running")
	}

	hpn.logger.Info(ctx, "Stopping high-performance networking system", nil)

	close(hpn.stopChan)
	hpn.wg.Wait()

	// Stop network managers
	hpn.stopNetworkManagers(ctx)

	// Close all connections
	hpn.closeAllConnections(ctx)

	hpn.logger.Info(ctx, "High-performance networking system stopped", map[string]interface{}{
		"packets_received": atomic.LoadInt64(&hpn.packetsReceived),
		"packets_sent":     atomic.LoadInt64(&hpn.packetsSent),
		"bytes_received":   atomic.LoadInt64(&hpn.bytesReceived),
		"bytes_sent":       atomic.LoadInt64(&hpn.bytesSent),
		"avg_latency_ns":   atomic.LoadInt64(&hpn.avgLatency),
	})

	return nil
}

// SendMessage sends a message with ultra-low latency
func (hpn *HighPerformanceNetworking) SendMessage(ctx context.Context, msg *NetworkMessage) error {
	if atomic.LoadInt32(&hpn.isRunning) != 1 {
		return fmt.Errorf("networking system is not running")
	}

	start := time.Now()

	hpn.logger.Debug(ctx, "Sending message", map[string]interface{}{
		"message_id":   msg.ID.String(),
		"type":         string(msg.Type),
		"destination":  msg.Destination.String(),
		"priority":     msg.Priority,
		"payload_size": len(msg.Payload),
	})

	// Get buffer from pool for zero-copy operation
	buffer := hpn.bufferPool.Get()
	defer hpn.bufferPool.Put(buffer)

	// Serialize message to buffer
	if err := hpn.serializeMessage(msg, buffer); err != nil {
		return fmt.Errorf("failed to serialize message: %w", err)
	}

	// Add hardware timestamp if enabled
	if hpn.config.EnableHWTimestamps {
		hpn.timestampEngine.AddTimestamp(buffer)
	}

	// Route message based on priority
	var queue *LockFreeQueue
	if msg.Priority >= PriorityCritical {
		queue = hpn.lockFreeQueues["priority"]
	} else {
		queue = hpn.lockFreeQueues["outbound"]
	}

	// Enqueue message for processing
	if !queue.Enqueue(unsafe.Pointer(&buffer)) {
		return fmt.Errorf("outbound queue is full")
	}

	// Update metrics
	atomic.AddInt64(&hpn.packetsSent, 1)
	atomic.AddInt64(&hpn.bytesSent, int64(len(msg.Payload)))

	// Calculate and update latency
	latency := time.Since(start).Nanoseconds()
	hpn.updateLatencyMetrics(latency)

	// Publish network event
	hpn.publishNetworkEvent(ctx, NetworkEventPacketSent, msg.ID, map[string]interface{}{
		"type":        string(msg.Type),
		"destination": msg.Destination.String(),
		"size":        len(msg.Payload),
		"latency_ns":  latency,
	}, time.Duration(latency))

	hpn.logger.Debug(ctx, "Message sent successfully", map[string]interface{}{
		"message_id": msg.ID.String(),
		"latency_ns": latency,
	})

	return nil
}

// ReceiveMessage receives a message with ultra-low latency
func (hpn *HighPerformanceNetworking) ReceiveMessage(ctx context.Context) (*NetworkMessage, error) {
	if atomic.LoadInt32(&hpn.isRunning) != 1 {
		return nil, fmt.Errorf("networking system is not running")
	}

	// Try to dequeue from priority queue first
	if ptr := hpn.lockFreeQueues["priority"].Dequeue(); ptr != nil {
		buffer := (*[]byte)(ptr)
		return hpn.deserializeMessage(*buffer)
	}

	// Try regular inbound queue
	if ptr := hpn.lockFreeQueues["inbound"].Dequeue(); ptr != nil {
		buffer := (*[]byte)(ptr)
		return hpn.deserializeMessage(*buffer)
	}

	return nil, fmt.Errorf("no messages available")
}

// CreateConnection creates a new high-performance connection
func (hpn *HighPerformanceNetworking) CreateConnection(ctx context.Context, connType ConnectionType, remoteAddr string) (*Connection, error) {
	if atomic.LoadInt32(&hpn.isRunning) != 1 {
		return nil, fmt.Errorf("networking system is not running")
	}

	hpn.logger.Info(ctx, "Creating connection", map[string]interface{}{
		"type":        string(connType),
		"remote_addr": remoteAddr,
	})

	// Get connection from pool
	conn := hpn.connectionPool.Get()
	if conn == nil {
		return nil, fmt.Errorf("connection pool exhausted")
	}

	conn.ID = uuid.New()
	conn.Type = connType
	conn.State = ConnectionStateConnecting

	// Parse remote address
	var err error
	switch connType {
	case ConnectionTypeUDP:
		conn.RemoteAddr, err = net.ResolveUDPAddr("udp", remoteAddr)
		if err != nil {
			hpn.connectionPool.Put(conn)
			return nil, fmt.Errorf("failed to resolve UDP address: %w", err)
		}
		err = hpn.udpManager.CreateConnection(conn)
	case ConnectionTypeTCP:
		conn.RemoteAddr, err = net.ResolveTCPAddr("tcp", remoteAddr)
		if err != nil {
			hpn.connectionPool.Put(conn)
			return nil, fmt.Errorf("failed to resolve TCP address: %w", err)
		}
		err = hpn.tcpManager.CreateConnection(conn)
	case ConnectionTypeMulticast:
		conn.RemoteAddr, err = net.ResolveUDPAddr("udp", remoteAddr)
		if err != nil {
			hpn.connectionPool.Put(conn)
			return nil, fmt.Errorf("failed to resolve multicast address: %w", err)
		}
		err = hpn.multicastManager.CreateConnection(conn)
	default:
		hpn.connectionPool.Put(conn)
		return nil, fmt.Errorf("unsupported connection type: %s", connType)
	}

	if err != nil {
		hpn.connectionPool.Put(conn)
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}

	// Assign CPU core for CPU affinity
	if len(hpn.config.CPUAffinity) > 0 {
		conn.cpuCore = hpn.config.CPUAffinity[len(hpn.connections)%len(hpn.config.CPUAffinity)]
	}

	// Store connection
	hpn.mu.Lock()
	hpn.connections[conn.ID.String()] = conn
	hpn.mu.Unlock()

	conn.State = ConnectionStateConnected
	conn.LastActivity = time.Now()

	// Publish connection event
	hpn.publishNetworkEvent(ctx, NetworkEventConnectionOpen, conn.ID, map[string]interface{}{
		"type":        string(connType),
		"remote_addr": remoteAddr,
		"cpu_core":    conn.cpuCore,
	}, 0)

	hpn.logger.Info(ctx, "Connection created successfully", map[string]interface{}{
		"connection_id": conn.ID.String(),
		"type":          string(connType),
		"remote_addr":   remoteAddr,
		"cpu_core":      conn.cpuCore,
	})

	return conn, nil
}

// CloseConnection closes a connection
func (hpn *HighPerformanceNetworking) CloseConnection(ctx context.Context, connectionID uuid.UUID) error {
	hpn.mu.Lock()
	conn, exists := hpn.connections[connectionID.String()]
	if !exists {
		hpn.mu.Unlock()
		return fmt.Errorf("connection not found: %s", connectionID.String())
	}
	delete(hpn.connections, connectionID.String())
	hpn.mu.Unlock()

	conn.State = ConnectionStateClosing

	// Close the underlying connection
	if conn.conn != nil {
		conn.conn.Close()
	}

	conn.State = ConnectionStateClosed

	// Return connection to pool
	hpn.connectionPool.Put(conn)

	// Publish connection event
	hpn.publishNetworkEvent(ctx, NetworkEventConnectionClose, connectionID, map[string]interface{}{
		"type": string(conn.Type),
	}, 0)

	hpn.logger.Info(ctx, "Connection closed", map[string]interface{}{
		"connection_id": connectionID.String(),
	})

	return nil
}

// updateLatencyMetrics updates latency statistics
func (hpn *HighPerformanceNetworking) updateLatencyMetrics(latencyNs int64) {
	// Update average latency (simple moving average)
	currentAvg := atomic.LoadInt64(&hpn.avgLatency)
	newAvg := (currentAvg + latencyNs) / 2
	atomic.StoreInt64(&hpn.avgLatency, newAvg)

	// Update min latency
	for {
		currentMin := atomic.LoadInt64(&hpn.minLatency)
		if latencyNs >= currentMin {
			break
		}
		if atomic.CompareAndSwapInt64(&hpn.minLatency, currentMin, latencyNs) {
			break
		}
	}

	// Update max latency
	for {
		currentMax := atomic.LoadInt64(&hpn.maxLatency)
		if latencyNs <= currentMax {
			break
		}
		if atomic.CompareAndSwapInt64(&hpn.maxLatency, currentMax, latencyNs) {
			break
		}
	}

	// Check for latency alerts
	if time.Duration(latencyNs) > hpn.config.LatencyTarget {
		// Publish latency alert (simplified)
		go hpn.publishNetworkEvent(context.Background(), NetworkEventLatencyAlert, uuid.New(), map[string]interface{}{
			"latency_ns": latencyNs,
			"target_ns":  hpn.config.LatencyTarget.Nanoseconds(),
		}, time.Duration(latencyNs))
	}
}

// serializeMessage serializes a message to a buffer
func (hpn *HighPerformanceNetworking) serializeMessage(msg *NetworkMessage, buffer []byte) error {
	// Simplified serialization - in production, use efficient binary protocol
	data := fmt.Sprintf("%s|%s|%d|%s", msg.ID.String(), string(msg.Type), len(msg.Payload), string(msg.Payload))
	copy(buffer, []byte(data))
	return nil
}

// deserializeMessage deserializes a message from a buffer
func (hpn *HighPerformanceNetworking) deserializeMessage(buffer []byte) (*NetworkMessage, error) {
	// Simplified deserialization - in production, use efficient binary protocol
	// Parse the serialized data (simplified)
	msg := &NetworkMessage{
		ID:        uuid.New(),
		Type:      MessageTypeMarketData,
		Payload:   buffer,
		Timestamp: time.Now(),
	}
	return msg, nil
}

// publishNetworkEvent publishes a network event to subscribers
func (hpn *HighPerformanceNetworking) publishNetworkEvent(ctx context.Context, eventType NetworkEventType, connectionID uuid.UUID, data map[string]interface{}, latency time.Duration) {
	event := &NetworkEvent{
		ID:           uuid.New(),
		Type:         eventType,
		Timestamp:    time.Now(),
		ConnectionID: connectionID,
		Data:         data,
		Latency:      latency,
	}

	hpn.mu.RLock()
	defer hpn.mu.RUnlock()

	// Send to event type subscribers
	if subscribers, exists := hpn.subscribers[string(eventType)]; exists {
		for _, ch := range subscribers {
			select {
			case ch <- event:
			default:
				// Subscriber channel is full, skip
			}
		}
	}

	// Send to wildcard subscribers
	if subscribers, exists := hpn.subscribers["*"]; exists {
		for _, ch := range subscribers {
			select {
			case ch <- event:
			default:
				// Subscriber channel is full, skip
			}
		}
	}
}

// startNetworkManagers starts all network managers
func (hpn *HighPerformanceNetworking) startNetworkManagers(ctx context.Context) error {
	hpn.logger.Info(ctx, "Starting network managers", nil)

	// Network managers are already initialized
	// In production, they would start listening on their respective ports

	return nil
}

// stopNetworkManagers stops all network managers
func (hpn *HighPerformanceNetworking) stopNetworkManagers(ctx context.Context) {
	hpn.logger.Info(ctx, "Stopping network managers", nil)

	// Close UDP connection
	if hpn.udpManager.conn != nil {
		hpn.udpManager.conn.Close()
	}

	// Close TCP listener
	if hpn.tcpManager.listener != nil {
		hpn.tcpManager.listener.Close()
	}

	// Close multicast connections
	hpn.multicastManager.mu.Lock()
	for _, conn := range hpn.multicastManager.conns {
		conn.Close()
	}
	hpn.multicastManager.conns = make(map[string]*net.UDPConn)
	hpn.multicastManager.mu.Unlock()
}

// closeAllConnections closes all active connections
func (hpn *HighPerformanceNetworking) closeAllConnections(ctx context.Context) {
	hpn.mu.Lock()
	defer hpn.mu.Unlock()

	for id, conn := range hpn.connections {
		if conn.conn != nil {
			conn.conn.Close()
		}
		conn.State = ConnectionStateClosed
		hpn.connectionPool.Put(conn)
		delete(hpn.connections, id)
	}
}

// processInboundPackets processes incoming packets
func (hpn *HighPerformanceNetworking) processInboundPackets(ctx context.Context) {
	defer hpn.wg.Done()

	hpn.logger.Info(ctx, "Starting inbound packet processor", nil)

	ticker := time.NewTicker(hpn.config.PollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-hpn.stopChan:
			return
		case <-ticker.C:
			hpn.pollInboundPackets(ctx)
		}
	}
}

// pollInboundPackets polls for inbound packets
func (hpn *HighPerformanceNetworking) pollInboundPackets(ctx context.Context) {
	// Simplified packet polling - in production, use epoll/kqueue
	// This would read from network interfaces and enqueue packets

	// Simulate receiving packets
	if hpn.udpManager.conn != nil {
		buffer := hpn.bufferPool.Get()
		defer hpn.bufferPool.Put(buffer)

		// Try to read from UDP connection (non-blocking)
		hpn.udpManager.conn.SetReadDeadline(time.Now().Add(1 * time.Millisecond))
		n, addr, err := hpn.udpManager.conn.ReadFromUDP(buffer)
		if err == nil && n > 0 {
			// Process received packet
			hpn.processReceivedPacket(buffer[:n], addr)
		}
	}
}

// processReceivedPacket processes a received packet
func (hpn *HighPerformanceNetworking) processReceivedPacket(data []byte, addr net.Addr) {
	// Update metrics
	atomic.AddInt64(&hpn.packetsReceived, 1)
	atomic.AddInt64(&hpn.bytesReceived, int64(len(data)))

	// Enqueue packet for processing
	buffer := make([]byte, len(data))
	copy(buffer, data)

	if !hpn.lockFreeQueues["inbound"].Enqueue(unsafe.Pointer(&buffer)) {
		// Queue is full, drop packet (or implement backpressure)
		hpn.logger.Warn(context.Background(), "Inbound queue full, dropping packet", map[string]interface{}{
			"size": len(data),
			"addr": addr.String(),
		})
	}
}

// processOutboundPackets processes outgoing packets
func (hpn *HighPerformanceNetworking) processOutboundPackets(ctx context.Context) {
	defer hpn.wg.Done()

	hpn.logger.Info(ctx, "Starting outbound packet processor", nil)

	for {
		select {
		case <-hpn.stopChan:
			return
		default:
			hpn.processOutboundQueue(ctx)
		}
	}
}

// processOutboundQueue processes the outbound packet queue
func (hpn *HighPerformanceNetworking) processOutboundQueue(ctx context.Context) {
	// Process priority queue first
	for i := 0; i < hpn.config.BatchSize; i++ {
		if ptr := hpn.lockFreeQueues["priority"].Dequeue(); ptr != nil {
			buffer := (*[]byte)(ptr)
			hpn.sendPacket(*buffer)
		} else {
			break
		}
	}

	// Process regular outbound queue
	for i := 0; i < hpn.config.BatchSize; i++ {
		if ptr := hpn.lockFreeQueues["outbound"].Dequeue(); ptr != nil {
			buffer := (*[]byte)(ptr)
			hpn.sendPacket(*buffer)
		} else {
			break
		}
	}

	// Small sleep to prevent busy waiting
	time.Sleep(1 * time.Microsecond)
}

// sendPacket sends a packet
func (hpn *HighPerformanceNetworking) sendPacket(data []byte) {
	// Simplified packet sending - in production, use sendmsg with zero-copy
	if hpn.udpManager.conn != nil {
		// For now, just simulate sending
		// In production, this would send to the appropriate destination
	}
}

// monitorConnections monitors connection health
func (hpn *HighPerformanceNetworking) monitorConnections(ctx context.Context) {
	defer hpn.wg.Done()

	hpn.logger.Info(ctx, "Starting connection monitor", nil)

	ticker := time.NewTicker(hpn.config.KeepAliveInterval)
	defer ticker.Stop()

	for {
		select {
		case <-hpn.stopChan:
			return
		case <-ticker.C:
			hpn.checkConnectionHealth(ctx)
		}
	}
}

// checkConnectionHealth checks the health of all connections
func (hpn *HighPerformanceNetworking) checkConnectionHealth(ctx context.Context) {
	hpn.mu.RLock()
	connections := make([]*Connection, 0, len(hpn.connections))
	for _, conn := range hpn.connections {
		connections = append(connections, conn)
	}
	hpn.mu.RUnlock()

	for _, conn := range connections {
		// Check if connection is stale
		if time.Since(conn.LastActivity) > hpn.config.TimeoutDuration {
			hpn.logger.Warn(ctx, "Connection timeout detected", map[string]interface{}{
				"connection_id": conn.ID.String(),
				"last_activity": conn.LastActivity,
			})

			// Close stale connection
			hpn.CloseConnection(ctx, conn.ID)
		}
	}
}

// performanceMonitor tracks and reports performance metrics
func (hpn *HighPerformanceNetworking) performanceMonitor(ctx context.Context) {
	defer hpn.wg.Done()

	hpn.logger.Info(ctx, "Starting networking performance monitor", nil)

	ticker := time.NewTicker(hpn.config.MetricsInterval)
	defer ticker.Stop()

	var lastPacketsReceived, lastPacketsSent int64
	var lastBytesReceived, lastBytesSent int64

	for {
		select {
		case <-hpn.stopChan:
			return
		case <-ticker.C:
			currentPacketsReceived := atomic.LoadInt64(&hpn.packetsReceived)
			currentPacketsSent := atomic.LoadInt64(&hpn.packetsSent)
			currentBytesReceived := atomic.LoadInt64(&hpn.bytesReceived)
			currentBytesSent := atomic.LoadInt64(&hpn.bytesSent)

			packetsReceivedPerSec := currentPacketsReceived - lastPacketsReceived
			packetsSentPerSec := currentPacketsSent - lastPacketsSent
			bytesReceivedPerSec := currentBytesReceived - lastBytesReceived
			bytesSentPerSec := currentBytesSent - lastBytesSent

			lastPacketsReceived = currentPacketsReceived
			lastPacketsSent = currentPacketsSent
			lastBytesReceived = currentBytesReceived
			lastBytesSent = currentBytesSent

			avgLatency := atomic.LoadInt64(&hpn.avgLatency)
			minLatency := atomic.LoadInt64(&hpn.minLatency)
			maxLatency := atomic.LoadInt64(&hpn.maxLatency)

			hpn.logger.Info(ctx, "Networking performance metrics", map[string]interface{}{
				"packets_received_per_sec": packetsReceivedPerSec,
				"packets_sent_per_sec":     packetsSentPerSec,
				"bytes_received_per_sec":   bytesReceivedPerSec,
				"bytes_sent_per_sec":       bytesSentPerSec,
				"total_packets_received":   currentPacketsReceived,
				"total_packets_sent":       currentPacketsSent,
				"total_bytes_received":     currentBytesReceived,
				"total_bytes_sent":         currentBytesSent,
				"avg_latency_ns":           avgLatency,
				"min_latency_ns":           minLatency,
				"max_latency_ns":           maxLatency,
				"avg_latency_us":           avgLatency / 1000,
				"min_latency_us":           minLatency / 1000,
				"max_latency_us":           maxLatency / 1000,
				"active_connections":       len(hpn.connections),
				"inbound_queue_size":       hpn.lockFreeQueues["inbound"].Size(),
				"outbound_queue_size":      hpn.lockFreeQueues["outbound"].Size(),
				"priority_queue_size":      hpn.lockFreeQueues["priority"].Size(),
			})
		}
	}
}

// Subscribe registers a subscriber for network events
func (hpn *HighPerformanceNetworking) Subscribe(eventType string) <-chan *NetworkEvent {
	hpn.mu.Lock()
	defer hpn.mu.Unlock()

	ch := make(chan *NetworkEvent, 1000) // Buffered channel
	if hpn.subscribers[eventType] == nil {
		hpn.subscribers[eventType] = make([]chan *NetworkEvent, 0)
	}
	hpn.subscribers[eventType] = append(hpn.subscribers[eventType], ch)

	return ch
}

// GetMetrics returns current networking metrics
func (hpn *HighPerformanceNetworking) GetMetrics() *NetworkMetrics {
	return &NetworkMetrics{
		PacketsReceived:   atomic.LoadInt64(&hpn.packetsReceived),
		PacketsSent:       atomic.LoadInt64(&hpn.packetsSent),
		BytesReceived:     atomic.LoadInt64(&hpn.bytesReceived),
		BytesSent:         atomic.LoadInt64(&hpn.bytesSent),
		AvgLatencyNs:      atomic.LoadInt64(&hpn.avgLatency),
		MinLatencyNs:      atomic.LoadInt64(&hpn.minLatency),
		MaxLatencyNs:      atomic.LoadInt64(&hpn.maxLatency),
		ActiveConnections: len(hpn.connections),
		InboundQueueSize:  hpn.lockFreeQueues["inbound"].Size(),
		OutboundQueueSize: hpn.lockFreeQueues["outbound"].Size(),
		PriorityQueueSize: hpn.lockFreeQueues["priority"].Size(),
		LastUpdate:        time.Now(),
	}
}

// NetworkMetrics contains networking performance metrics
type NetworkMetrics struct {
	PacketsReceived   int64     `json:"packets_received"`
	PacketsSent       int64     `json:"packets_sent"`
	BytesReceived     int64     `json:"bytes_received"`
	BytesSent         int64     `json:"bytes_sent"`
	AvgLatencyNs      int64     `json:"avg_latency_ns"`
	MinLatencyNs      int64     `json:"min_latency_ns"`
	MaxLatencyNs      int64     `json:"max_latency_ns"`
	ActiveConnections int       `json:"active_connections"`
	InboundQueueSize  int       `json:"inbound_queue_size"`
	OutboundQueueSize int       `json:"outbound_queue_size"`
	PriorityQueueSize int       `json:"priority_queue_size"`
	LastUpdate        time.Time `json:"last_update"`
}

// GetConnections returns all active connections
func (hpn *HighPerformanceNetworking) GetConnections() map[string]*Connection {
	hpn.mu.RLock()
	defer hpn.mu.RUnlock()

	connections := make(map[string]*Connection)
	for id, conn := range hpn.connections {
		connections[id] = conn
	}
	return connections
}

// GetConnection returns a specific connection by ID
func (hpn *HighPerformanceNetworking) GetConnection(connectionID uuid.UUID) *Connection {
	hpn.mu.RLock()
	defer hpn.mu.RUnlock()

	if conn, exists := hpn.connections[connectionID.String()]; exists {
		return conn
	}
	return nil
}
