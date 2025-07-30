package tradingview

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/internal/hft"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// SignalProcessor processes TradingView signals and converts them to HFT trading signals
type SignalProcessor struct {
	logger    *observability.Logger
	config    ProcessorConfig
	tvClient  *TradingViewClient
	hftEngine *hft.HFTEngine

	// Signal processing
	signalQueue      chan Signal
	processedSignals map[string]*ProcessedSignal
	signalFilters    []SignalFilter

	// Performance tracking
	signalsProcessed int64
	signalsGenerated int64
	signalsFiltered  int64

	// State management
	isRunning bool
	stopChan  chan struct{}
	wg        sync.WaitGroup
	mu        sync.RWMutex
}

// ProcessorConfig contains configuration for the signal processor
type ProcessorConfig struct {
	QueueSize           int           `json:"queue_size"`
	ProcessingInterval  time.Duration `json:"processing_interval"`
	SignalTimeout       time.Duration `json:"signal_timeout"`
	MinConfidence       float64       `json:"min_confidence"`
	MaxSignalsPerSymbol int           `json:"max_signals_per_symbol"`
	EnableFiltering     bool          `json:"enable_filtering"`
	EnableAggregation   bool          `json:"enable_aggregation"`
}

// ProcessedSignal represents a processed TradingView signal ready for HFT
type ProcessedSignal struct {
	ID             uuid.UUID              `json:"id"`
	OriginalSignal Signal                 `json:"original_signal"`
	HFTSignal      hft.TradingSignal      `json:"hft_signal"`
	Confidence     float64                `json:"confidence"`
	ProcessedAt    time.Time              `json:"processed_at"`
	Status         ProcessedSignalStatus  `json:"status"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// ProcessedSignalStatus represents the status of a processed signal
type ProcessedSignalStatus string

const (
	ProcessedSignalStatusPending  ProcessedSignalStatus = "PENDING"
	ProcessedSignalStatusApproved ProcessedSignalStatus = "APPROVED"
	ProcessedSignalStatusRejected ProcessedSignalStatus = "REJECTED"
	ProcessedSignalStatusExecuted ProcessedSignalStatus = "EXECUTED"
	ProcessedSignalStatusExpired  ProcessedSignalStatus = "EXPIRED"
)

// SignalFilter defines a filter for signals
type SignalFilter struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Type       SignalFilterType       `json:"type"`
	Parameters map[string]interface{} `json:"parameters"`
	Enabled    bool                   `json:"enabled"`
	FilterFunc func(Signal) bool      `json:"-"`
}

// SignalFilterType represents different types of signal filters
type SignalFilterType string

const (
	SignalFilterTypeStrength  SignalFilterType = "STRENGTH"
	SignalFilterTypeTimeframe SignalFilterType = "TIMEFRAME"
	SignalFilterTypeSymbol    SignalFilterType = "SYMBOL"
	SignalFilterTypeIndicator SignalFilterType = "INDICATOR"
	SignalFilterTypeVolume    SignalFilterType = "VOLUME"
	SignalFilterTypePrice     SignalFilterType = "PRICE"
	SignalFilterTypeCustom    SignalFilterType = "CUSTOM"
)

// NewSignalProcessor creates a new signal processor
func NewSignalProcessor(logger *observability.Logger, config ProcessorConfig, tvClient *TradingViewClient) *SignalProcessor {
	if config.QueueSize == 0 {
		config.QueueSize = 1000
	}

	if config.ProcessingInterval == 0 {
		config.ProcessingInterval = time.Second
	}

	if config.SignalTimeout == 0 {
		config.SignalTimeout = 5 * time.Minute
	}

	if config.MinConfidence == 0 {
		config.MinConfidence = 0.6
	}

	if config.MaxSignalsPerSymbol == 0 {
		config.MaxSignalsPerSymbol = 10
	}

	processor := &SignalProcessor{
		logger:           logger,
		config:           config,
		tvClient:         tvClient,
		signalQueue:      make(chan Signal, config.QueueSize),
		processedSignals: make(map[string]*ProcessedSignal),
		signalFilters:    make([]SignalFilter, 0),
		stopChan:         make(chan struct{}),
	}

	// Initialize default filters
	processor.initializeDefaultFilters()

	return processor
}

// Start begins the signal processor
func (sp *SignalProcessor) Start(ctx context.Context) error {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	if sp.isRunning {
		return fmt.Errorf("signal processor is already running")
	}

	sp.logger.Info(ctx, "Starting TradingView signal processor", map[string]interface{}{
		"queue_size":             sp.config.QueueSize,
		"processing_interval":    sp.config.ProcessingInterval.String(),
		"min_confidence":         sp.config.MinConfidence,
		"max_signals_per_symbol": sp.config.MaxSignalsPerSymbol,
		"filters_count":          len(sp.signalFilters),
	})

	sp.isRunning = true

	// Start processing goroutines
	sp.wg.Add(3)
	go sp.processSignals(ctx)
	go sp.collectSignals(ctx)
	go sp.cleanupExpiredSignals(ctx)

	sp.logger.Info(ctx, "TradingView signal processor started successfully", nil)

	return nil
}

// Stop gracefully shuts down the signal processor
func (sp *SignalProcessor) Stop(ctx context.Context) error {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	if !sp.isRunning {
		return fmt.Errorf("signal processor is not running")
	}

	sp.logger.Info(ctx, "Stopping TradingView signal processor", nil)

	sp.isRunning = false
	close(sp.stopChan)

	sp.wg.Wait()

	sp.logger.Info(ctx, "TradingView signal processor stopped", map[string]interface{}{
		"signals_processed": sp.signalsProcessed,
		"signals_generated": sp.signalsGenerated,
		"signals_filtered":  sp.signalsFiltered,
	})

	return nil
}

// SetHFTEngine sets the HFT engine for signal forwarding
func (sp *SignalProcessor) SetHFTEngine(engine *hft.HFTEngine) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	sp.hftEngine = engine
}

// SubmitSignal submits a TradingView signal for processing
func (sp *SignalProcessor) SubmitSignal(signal Signal) error {
	select {
	case sp.signalQueue <- signal:
		return nil
	default:
		return fmt.Errorf("signal queue is full")
	}
}

// GetProcessedSignals returns all processed signals
func (sp *SignalProcessor) GetProcessedSignals() map[string]*ProcessedSignal {
	sp.mu.RLock()
	defer sp.mu.RUnlock()

	signals := make(map[string]*ProcessedSignal)
	for id, signal := range sp.processedSignals {
		signals[id] = signal
	}

	return signals
}

// AddSignalFilter adds a custom signal filter
func (sp *SignalProcessor) AddSignalFilter(filter SignalFilter) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	sp.signalFilters = append(sp.signalFilters, filter)

	sp.logger.Info(context.Background(), "Signal filter added", map[string]interface{}{
		"filter_id":   filter.ID,
		"filter_name": filter.Name,
		"filter_type": string(filter.Type),
	})
}

// collectSignals collects signals from TradingView client
func (sp *SignalProcessor) collectSignals(ctx context.Context) {
	defer sp.wg.Done()

	ticker := time.NewTicker(sp.config.ProcessingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-sp.stopChan:
			return
		case <-ticker.C:
			sp.collectFromTradingView(ctx)
		}
	}
}

// collectFromTradingView collects signals from TradingView client
func (sp *SignalProcessor) collectFromTradingView(ctx context.Context) {
	if sp.tvClient == nil {
		return
	}

	// Get all charts and their signals
	charts := sp.tvClient.charts

	sp.mu.RLock()
	chartsCopy := make(map[string]*Chart)
	for id, chart := range charts {
		chartsCopy[id] = chart
	}
	sp.mu.RUnlock()

	for _, chart := range chartsCopy {
		signals := sp.tvClient.GetSignals(chart.Symbol)

		for _, signal := range signals {
			// Check if signal is recent
			if time.Since(signal.Timestamp) < sp.config.SignalTimeout {
				if err := sp.SubmitSignal(signal); err != nil {
					sp.logger.Warn(ctx, "Failed to submit signal", map[string]interface{}{
						"signal_id": signal.ID,
						"symbol":    signal.Symbol,
						"error":     err.Error(),
					})
				}
			}
		}
	}
}

// processSignals processes signals from the queue
func (sp *SignalProcessor) processSignals(ctx context.Context) {
	defer sp.wg.Done()

	for {
		select {
		case <-sp.stopChan:
			return
		case signal := <-sp.signalQueue:
			if err := sp.processSignal(ctx, signal); err != nil {
				sp.logger.Error(ctx, "Failed to process signal", err, map[string]interface{}{
					"signal_id": signal.ID,
					"symbol":    signal.Symbol,
				})
			}
		}
	}
}

// processSignal processes a single signal
func (sp *SignalProcessor) processSignal(ctx context.Context, signal Signal) error {
	sp.signalsProcessed++

	// Apply filters
	if sp.config.EnableFiltering && !sp.applyFilters(signal) {
		sp.signalsFiltered++
		sp.logger.Debug(ctx, "Signal filtered out", map[string]interface{}{
			"signal_id": signal.ID,
			"symbol":    signal.Symbol,
		})
		return nil
	}

	// Convert to HFT signal
	hftSignal, confidence, err := sp.convertToHFTSignal(signal)
	if err != nil {
		return fmt.Errorf("failed to convert signal: %w", err)
	}

	// Check confidence threshold
	if confidence < sp.config.MinConfidence {
		sp.signalsFiltered++
		sp.logger.Debug(ctx, "Signal below confidence threshold", map[string]interface{}{
			"signal_id":  signal.ID,
			"confidence": confidence,
			"threshold":  sp.config.MinConfidence,
		})
		return nil
	}

	// Create processed signal
	processedSignal := &ProcessedSignal{
		ID:             uuid.New(),
		OriginalSignal: signal,
		HFTSignal:      hftSignal,
		Confidence:     confidence,
		ProcessedAt:    time.Now(),
		Status:         ProcessedSignalStatusPending,
		Metadata: map[string]interface{}{
			"processor_version": "1.0",
			"filters_applied":   len(sp.signalFilters),
		},
	}

	// Store processed signal
	sp.mu.Lock()
	sp.processedSignals[processedSignal.ID.String()] = processedSignal
	sp.mu.Unlock()

	// Forward to HFT engine if available
	if sp.hftEngine != nil {
		// Note: This would require implementing a method to submit signals to HFT engine
		sp.logger.Info(ctx, "Signal forwarded to HFT engine", map[string]interface{}{
			"signal_id":  processedSignal.ID.String(),
			"symbol":     signal.Symbol,
			"direction":  string(signal.Direction),
			"confidence": confidence,
		})

		processedSignal.Status = ProcessedSignalStatusApproved
	}

	sp.signalsGenerated++

	return nil
}

// applyFilters applies all enabled filters to a signal
func (sp *SignalProcessor) applyFilters(signal Signal) bool {
	for _, filter := range sp.signalFilters {
		if !filter.Enabled {
			continue
		}

		if filter.FilterFunc != nil && !filter.FilterFunc(signal) {
			return false
		}
	}

	return true
}

// convertToHFTSignal converts a TradingView signal to an HFT trading signal
func (sp *SignalProcessor) convertToHFTSignal(signal Signal) (hft.TradingSignal, float64, error) {
	// Convert signal direction to HFT order side
	var side hft.OrderSide
	switch signal.Direction {
	case SignalDirectionBuy:
		side = hft.OrderSideBuy
	case SignalDirectionSell:
		side = hft.OrderSideSell
	default:
		return hft.TradingSignal{}, 0, fmt.Errorf("unsupported signal direction: %s", signal.Direction)
	}

	// Calculate confidence based on signal strength
	confidence := sp.calculateConfidence(signal)

	// Determine order type (default to market for HFT)
	orderType := hft.OrderTypeMarket

	// Calculate quantity (simplified - would use position sizing in real implementation)
	quantity := decimal.NewFromFloat(0.01) // Default small quantity

	// Create HFT signal
	hftSignal := hft.TradingSignal{
		ID:         uuid.New(),
		Symbol:     signal.Symbol,
		Side:       side,
		OrderType:  orderType,
		Quantity:   quantity,
		Price:      signal.Price,
		Confidence: confidence,
		StrategyID: "tradingview_signals",
		Timestamp:  time.Now(),
		Metadata: map[string]interface{}{
			"original_signal_id": signal.ID,
			"signal_type":        string(signal.Type),
			"signal_strength":    string(signal.Strength),
			"timeframe":          signal.Timeframe,
			"indicators":         signal.Indicators,
		},
	}

	return hftSignal, confidence, nil
}

// calculateConfidence calculates confidence score for a signal
func (sp *SignalProcessor) calculateConfidence(signal Signal) float64 {
	baseConfidence := 0.5

	// Adjust based on signal strength
	switch signal.Strength {
	case SignalStrengthVeryStrong:
		baseConfidence = 0.9
	case SignalStrengthStrong:
		baseConfidence = 0.8
	case SignalStrengthModerate:
		baseConfidence = 0.6
	case SignalStrengthWeak:
		baseConfidence = 0.4
	}

	// Adjust based on signal type
	switch signal.Type {
	case SignalTypeTechnical:
		baseConfidence += 0.1
	case SignalTypePattern:
		baseConfidence += 0.05
	case SignalTypeMomentum:
		baseConfidence += 0.05
	}

	// Adjust based on number of indicators
	if len(signal.Indicators) > 1 {
		baseConfidence += float64(len(signal.Indicators)-1) * 0.05
	}

	// Ensure confidence is within bounds
	if baseConfidence > 1.0 {
		baseConfidence = 1.0
	}
	if baseConfidence < 0.0 {
		baseConfidence = 0.0
	}

	return baseConfidence
}

// cleanupExpiredSignals removes expired processed signals
func (sp *SignalProcessor) cleanupExpiredSignals(ctx context.Context) {
	defer sp.wg.Done()

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-sp.stopChan:
			return
		case <-ticker.C:
			sp.performCleanup(ctx)
		}
	}
}

// performCleanup removes expired signals
func (sp *SignalProcessor) performCleanup(ctx context.Context) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	expiredCount := 0
	now := time.Now()

	for id, signal := range sp.processedSignals {
		if now.Sub(signal.ProcessedAt) > sp.config.SignalTimeout {
			if signal.Status == ProcessedSignalStatusPending {
				signal.Status = ProcessedSignalStatusExpired
			}

			// Remove signals older than 1 hour
			if now.Sub(signal.ProcessedAt) > time.Hour {
				delete(sp.processedSignals, id)
				expiredCount++
			}
		}
	}

	if expiredCount > 0 {
		sp.logger.Info(ctx, "Cleaned up expired signals", map[string]interface{}{
			"expired_count": expiredCount,
		})
	}
}

// initializeDefaultFilters sets up default signal filters
func (sp *SignalProcessor) initializeDefaultFilters() {
	// Strength filter
	sp.signalFilters = append(sp.signalFilters, SignalFilter{
		ID:      "strength_filter",
		Name:    "Minimum Strength Filter",
		Type:    SignalFilterTypeStrength,
		Enabled: true,
		FilterFunc: func(signal Signal) bool {
			return signal.Strength == SignalStrengthStrong || signal.Strength == SignalStrengthVeryStrong
		},
	})

	// Timeframe filter (prefer shorter timeframes for HFT)
	sp.signalFilters = append(sp.signalFilters, SignalFilter{
		ID:      "timeframe_filter",
		Name:    "HFT Timeframe Filter",
		Type:    SignalFilterTypeTimeframe,
		Enabled: true,
		FilterFunc: func(signal Signal) bool {
			// Prefer 1m, 5m, 15m timeframes for HFT
			return signal.Timeframe == "1m" || signal.Timeframe == "5m" || signal.Timeframe == "15m"
		},
	})

	// Technical signal filter
	sp.signalFilters = append(sp.signalFilters, SignalFilter{
		ID:      "technical_filter",
		Name:    "Technical Signal Filter",
		Type:    SignalFilterTypeIndicator,
		Enabled: true,
		FilterFunc: func(signal Signal) bool {
			return signal.Type == SignalTypeTechnical || signal.Type == SignalTypeMomentum
		},
	})
}

// GetMetrics returns processor metrics
func (sp *SignalProcessor) GetMetrics() SignalProcessorMetrics {
	sp.mu.RLock()
	defer sp.mu.RUnlock()

	return SignalProcessorMetrics{
		IsRunning:          sp.isRunning,
		SignalsProcessed:   sp.signalsProcessed,
		SignalsGenerated:   sp.signalsGenerated,
		SignalsFiltered:    sp.signalsFiltered,
		QueueSize:          len(sp.signalQueue),
		QueueCapacity:      cap(sp.signalQueue),
		ProcessedSignals:   len(sp.processedSignals),
		FiltersCount:       len(sp.signalFilters),
		HFTEngineConnected: sp.hftEngine != nil,
	}
}

// SignalProcessorMetrics contains processor performance metrics
type SignalProcessorMetrics struct {
	IsRunning          bool  `json:"is_running"`
	SignalsProcessed   int64 `json:"signals_processed"`
	SignalsGenerated   int64 `json:"signals_generated"`
	SignalsFiltered    int64 `json:"signals_filtered"`
	QueueSize          int   `json:"queue_size"`
	QueueCapacity      int   `json:"queue_capacity"`
	ProcessedSignals   int   `json:"processed_signals"`
	FiltersCount       int   `json:"filters_count"`
	HFTEngineConnected bool  `json:"hft_engine_connected"`
}
