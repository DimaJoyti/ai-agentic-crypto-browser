package tradingview

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/internal/hft"
	"github.com/ai-agentic-browser/pkg/observability"
)

// Service provides TradingView integration for the HFT system
type Service struct {
	logger          *observability.Logger
	config          ServiceConfig
	client          *TradingViewClient
	signalProcessor *SignalProcessor
	hftEngine       *hft.HFTEngine

	// State management
	isRunning bool
	stopChan  chan struct{}
	wg        sync.WaitGroup
	mu        sync.RWMutex
}

// ServiceConfig contains configuration for the TradingView service
type ServiceConfig struct {
	TradingView         Config          `json:"tradingview"`
	SignalProcessor     ProcessorConfig `json:"signal_processor"`
	DefaultSymbols      []string        `json:"default_symbols"`
	DefaultTimeframes   []string        `json:"default_timeframes"`
	AutoStart           bool            `json:"auto_start"`
	HealthCheckInterval time.Duration   `json:"health_check_interval"`
}

// NewService creates a new TradingView service
func NewService(logger *observability.Logger, config ServiceConfig) *Service {
	if config.HealthCheckInterval == 0 {
		config.HealthCheckInterval = 30 * time.Second
	}

	// Set default symbols if not provided
	if len(config.DefaultSymbols) == 0 {
		config.DefaultSymbols = []string{
			"BINANCE:BTCUSDT", "BINANCE:ETHUSDT", "BINANCE:BNBUSDT",
			"BINANCE:ADAUSDT", "BINANCE:XRPUSDT", "BINANCE:SOLUSDT",
		}
	}

	// Set default timeframes if not provided
	if len(config.DefaultTimeframes) == 0 {
		config.DefaultTimeframes = []string{"1m", "5m", "15m"}
	}

	service := &Service{
		logger:   logger,
		config:   config,
		stopChan: make(chan struct{}),
	}

	// Create TradingView client
	service.client = NewTradingViewClient(logger, config.TradingView)

	// Create signal processor
	service.signalProcessor = NewSignalProcessor(logger, config.SignalProcessor, service.client)

	return service
}

// Start begins the TradingView service
func (s *Service) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return fmt.Errorf("TradingView service is already running")
	}

	s.logger.Info(ctx, "Starting TradingView service", map[string]interface{}{
		"default_symbols":    s.config.DefaultSymbols,
		"default_timeframes": s.config.DefaultTimeframes,
		"auto_start":         s.config.AutoStart,
	})

	// Start TradingView client
	if err := s.client.Start(ctx); err != nil {
		return fmt.Errorf("failed to start TradingView client: %w", err)
	}

	// Start signal processor
	if err := s.signalProcessor.Start(ctx); err != nil {
		return fmt.Errorf("failed to start signal processor: %w", err)
	}

	s.isRunning = true

	// Auto-setup charts if enabled
	if s.config.AutoStart {
		if err := s.setupDefaultCharts(ctx); err != nil {
			s.logger.Warn(ctx, "Failed to setup default charts", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	// Start monitoring goroutines
	s.wg.Add(1)
	go s.healthMonitor(ctx)

	s.logger.Info(ctx, "TradingView service started successfully", nil)

	return nil
}

// Stop gracefully shuts down the TradingView service
func (s *Service) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return fmt.Errorf("TradingView service is not running")
	}

	s.logger.Info(ctx, "Stopping TradingView service", nil)

	s.isRunning = false
	close(s.stopChan)

	// Stop signal processor
	if err := s.signalProcessor.Stop(ctx); err != nil {
		s.logger.Error(ctx, "Failed to stop signal processor", err)
	}

	// Stop TradingView client
	if err := s.client.Stop(ctx); err != nil {
		s.logger.Error(ctx, "Failed to stop TradingView client", err)
	}

	s.wg.Wait()

	s.logger.Info(ctx, "TradingView service stopped successfully", nil)

	return nil
}

// SetHFTEngine sets the HFT engine for signal forwarding
func (s *Service) SetHFTEngine(engine *hft.HFTEngine) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.hftEngine = engine

	// Forward to signal processor
	if s.signalProcessor != nil {
		s.signalProcessor.SetHFTEngine(engine)
	}
}

// GetClient returns the TradingView client
func (s *Service) GetClient() *TradingViewClient {
	return s.client
}

// GetSignalProcessor returns the signal processor
func (s *Service) GetSignalProcessor() *SignalProcessor {
	return s.signalProcessor
}

// AddChart adds a chart for monitoring
func (s *Service) AddChart(ctx context.Context, symbol, timeframe string) (*Chart, error) {
	if s.client == nil {
		return nil, fmt.Errorf("TradingView client not initialized")
	}

	chart, err := s.client.AddChart(symbol, timeframe)
	if err != nil {
		return nil, fmt.Errorf("failed to add chart: %w", err)
	}

	s.logger.Info(ctx, "Chart added to TradingView service", map[string]interface{}{
		"chart_id":  chart.ID,
		"symbol":    symbol,
		"timeframe": timeframe,
	})

	return chart, nil
}

// GetChart retrieves a chart by ID
func (s *Service) GetChart(chartID string) (*Chart, error) {
	if s.client == nil {
		return nil, fmt.Errorf("TradingView client not initialized")
	}

	return s.client.GetChart(chartID)
}

// GetSignals retrieves signals for a symbol
func (s *Service) GetSignals(symbol string) []Signal {
	if s.client == nil {
		return nil
	}

	return s.client.GetSignals(symbol)
}

// GetProcessedSignals returns all processed signals
func (s *Service) GetProcessedSignals() map[string]*ProcessedSignal {
	if s.signalProcessor == nil {
		return nil
	}

	return s.signalProcessor.GetProcessedSignals()
}

// GetIndicatorValue gets the latest value for an indicator
func (s *Service) GetIndicatorValue(symbol, indicatorType string) (*IndicatorValue, error) {
	if s.client == nil {
		return nil, fmt.Errorf("TradingView client not initialized")
	}

	return s.client.GetIndicatorValue(symbol, indicatorType)
}

// AddSignalFilter adds a custom signal filter
func (s *Service) AddSignalFilter(filter SignalFilter) error {
	if s.signalProcessor == nil {
		return fmt.Errorf("signal processor not initialized")
	}

	s.signalProcessor.AddSignalFilter(filter)

	s.logger.Info(context.Background(), "Signal filter added to TradingView service", map[string]interface{}{
		"filter_id":   filter.ID,
		"filter_name": filter.Name,
		"filter_type": string(filter.Type),
	})

	return nil
}

// setupDefaultCharts sets up default charts for monitoring
func (s *Service) setupDefaultCharts(ctx context.Context) error {
	for _, symbol := range s.config.DefaultSymbols {
		for _, timeframe := range s.config.DefaultTimeframes {
			if _, err := s.AddChart(ctx, symbol, timeframe); err != nil {
				s.logger.Warn(ctx, "Failed to add default chart", map[string]interface{}{
					"symbol":    symbol,
					"timeframe": timeframe,
					"error":     err.Error(),
				})
			}
		}
	}

	s.logger.Info(ctx, "Default charts setup completed", map[string]interface{}{
		"symbols_count":    len(s.config.DefaultSymbols),
		"timeframes_count": len(s.config.DefaultTimeframes),
		"total_charts":     len(s.config.DefaultSymbols) * len(s.config.DefaultTimeframes),
	})

	return nil
}

// healthMonitor monitors the health of the TradingView service
func (s *Service) healthMonitor(ctx context.Context) {
	defer s.wg.Done()

	ticker := time.NewTicker(s.config.HealthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			s.performHealthCheck(ctx)
		}
	}
}

// performHealthCheck performs a health check of the service
func (s *Service) performHealthCheck(ctx context.Context) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.isRunning {
		return
	}

	// Check client health
	clientMetrics := s.client.GetMetrics()

	// Check signal processor health
	processorMetrics := s.signalProcessor.GetMetrics()

	// Overall health assessment
	isHealthy := clientMetrics.IsRunning &&
		clientMetrics.BrowserActive &&
		processorMetrics.IsRunning &&
		processorMetrics.QueueSize < processorMetrics.QueueCapacity

	s.logger.Info(ctx, "TradingView service health check", map[string]interface{}{
		"is_healthy":               isHealthy,
		"client_running":           clientMetrics.IsRunning,
		"client_browser_active":    clientMetrics.BrowserActive,
		"client_charts_count":      clientMetrics.ChartsCount,
		"client_signals_count":     clientMetrics.SignalsCount,
		"processor_running":        processorMetrics.IsRunning,
		"processor_queue_size":     processorMetrics.QueueSize,
		"processor_queue_capacity": processorMetrics.QueueCapacity,
		"signals_processed":        processorMetrics.SignalsProcessed,
		"signals_generated":        processorMetrics.SignalsGenerated,
		"hft_engine_connected":     processorMetrics.HFTEngineConnected,
	})

	// Alert if unhealthy
	if !isHealthy {
		s.logger.Warn(ctx, "TradingView service health issues detected", map[string]interface{}{
			"client_issues":    !clientMetrics.IsRunning || !clientMetrics.BrowserActive,
			"processor_issues": !processorMetrics.IsRunning || processorMetrics.QueueSize >= processorMetrics.QueueCapacity,
		})
	}
}

// GetMetrics returns service metrics
func (s *Service) GetMetrics() ServiceMetrics {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var clientMetrics TradingViewMetrics
	var processorMetrics SignalProcessorMetrics

	if s.client != nil {
		clientMetrics = s.client.GetMetrics()
	}

	if s.signalProcessor != nil {
		processorMetrics = s.signalProcessor.GetMetrics()
	}

	return ServiceMetrics{
		IsRunning:          s.isRunning,
		ClientMetrics:      clientMetrics,
		ProcessorMetrics:   processorMetrics,
		HFTEngineConnected: s.hftEngine != nil,
		DefaultChartsCount: len(s.config.DefaultSymbols) * len(s.config.DefaultTimeframes),
	}
}

// ServiceMetrics contains service performance metrics
type ServiceMetrics struct {
	IsRunning          bool                   `json:"is_running"`
	ClientMetrics      TradingViewMetrics     `json:"client_metrics"`
	ProcessorMetrics   SignalProcessorMetrics `json:"processor_metrics"`
	HFTEngineConnected bool                   `json:"hft_engine_connected"`
	DefaultChartsCount int                    `json:"default_charts_count"`
}

// IsHealthy returns whether the service is healthy
func (s *Service) IsHealthy() bool {
	metrics := s.GetMetrics()
	return metrics.IsRunning &&
		metrics.ClientMetrics.IsRunning &&
		metrics.ClientMetrics.BrowserActive &&
		metrics.ProcessorMetrics.IsRunning &&
		metrics.ProcessorMetrics.QueueSize < metrics.ProcessorMetrics.QueueCapacity
}

// GetSignalSummary returns a summary of recent signals
func (s *Service) GetSignalSummary() map[string]interface{} {
	processedSignals := s.GetProcessedSignals()

	summary := map[string]interface{}{
		"total_signals":  len(processedSignals),
		"by_status":      make(map[string]int),
		"by_symbol":      make(map[string]int),
		"by_direction":   make(map[string]int),
		"recent_signals": make([]map[string]interface{}, 0),
	}

	// Count by status
	statusCounts := make(map[string]int)
	symbolCounts := make(map[string]int)
	directionCounts := make(map[string]int)

	// Recent signals (last 10)
	recentSignals := make([]*ProcessedSignal, 0)

	for _, signal := range processedSignals {
		statusCounts[string(signal.Status)]++
		symbolCounts[signal.OriginalSignal.Symbol]++
		directionCounts[string(signal.OriginalSignal.Direction)]++

		recentSignals = append(recentSignals, signal)
	}

	// Sort recent signals by timestamp (simplified)
	if len(recentSignals) > 10 {
		recentSignals = recentSignals[:10]
	}

	// Convert to summary format
	recentSummary := make([]map[string]interface{}, len(recentSignals))
	for i, signal := range recentSignals {
		recentSummary[i] = map[string]interface{}{
			"id":         signal.ID.String(),
			"symbol":     signal.OriginalSignal.Symbol,
			"direction":  string(signal.OriginalSignal.Direction),
			"strength":   string(signal.OriginalSignal.Strength),
			"confidence": signal.Confidence,
			"status":     string(signal.Status),
			"timestamp":  signal.ProcessedAt,
		}
	}

	summary["by_status"] = statusCounts
	summary["by_symbol"] = symbolCounts
	summary["by_direction"] = directionCounts
	summary["recent_signals"] = recentSummary

	return summary
}
