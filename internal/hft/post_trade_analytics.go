package hft

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// PostTradeAnalytics provides comprehensive post-trade analysis with execution quality metrics,
// slippage analysis, market impact measurement, and regulatory reporting capabilities
type PostTradeAnalytics struct {
	logger *observability.Logger
	config PostTradeConfig

	// Analytics engines
	executionAnalyzer   *ExecutionAnalyzer
	slippageCalculator  *SlippageCalculator
	marketImpactTracker *MarketImpactTracker
	performanceEngine   *PerformanceEngine
	reportGenerator     *ReportGenerator

	// Data storage
	tradeHistory     []TradeRecord
	executionMetrics map[string]*ExecutionMetrics
	performanceData  map[string]*PerformanceData

	// Real-time analytics
	realtimeMetrics *RealtimeMetrics
	alertThresholds *AlertThresholds

	// Performance tracking
	tradesAnalyzed   int64
	reportsGenerated int64
	avgAnalysisTime  int64

	// State management
	isRunning int32
	stopChan  chan struct{}
	wg        sync.WaitGroup
	mu        sync.RWMutex

	// Subscribers for analytics events
	subscribers map[string][]chan *AnalyticsEvent
}

// PostTradeConfig contains configuration for post-trade analytics
type PostTradeConfig struct {
	// Analysis settings
	AnalysisInterval  time.Duration `json:"analysis_interval"`  // How often to run analysis
	HistoryRetention  time.Duration `json:"history_retention"`  // How long to keep trade history
	ReportingInterval time.Duration `json:"reporting_interval"` // Report generation frequency

	// Execution quality thresholds
	SlippageThreshold float64       `json:"slippage_threshold"`  // Alert threshold for slippage %
	FillRateThreshold float64       `json:"fill_rate_threshold"` // Minimum acceptable fill rate
	LatencyThreshold  time.Duration `json:"latency_threshold"`   // Maximum acceptable latency

	// Market impact settings
	ImpactMeasureWindow time.Duration `json:"impact_measure_window"` // Window for impact measurement
	BenchmarkType       string        `json:"benchmark_type"`        // TWAP, VWAP, Arrival, etc.

	// Performance measurement
	EnableTCA         bool `json:"enable_tca"`         // Transaction Cost Analysis
	EnableAttribution bool `json:"enable_attribution"` // Performance attribution
	EnableCompliance  bool `json:"enable_compliance"`  // Regulatory compliance

	// Reporting settings
	ReportFormats    []string `json:"report_formats"`    // PDF, CSV, JSON
	AutoReporting    bool     `json:"auto_reporting"`    // Automatic report generation
	ReportRecipients []string `json:"report_recipients"` // Email recipients
}

// TradeRecord represents a complete trade record for analysis
type TradeRecord struct {
	// Trade identification
	TradeID       uuid.UUID `json:"trade_id"`
	OrderID       uuid.UUID `json:"order_id"`
	ParentOrderID uuid.UUID `json:"parent_order_id,omitempty"`
	StrategyID    string    `json:"strategy_id"`

	// Trade details
	Symbol   string          `json:"symbol"`
	Side     OrderSide       `json:"side"`
	Quantity decimal.Decimal `json:"quantity"`
	Price    decimal.Decimal `json:"price"`
	Value    decimal.Decimal `json:"value"`

	// Execution details
	VenueID        string    `json:"venue_id"`
	ExecutionTime  time.Time `json:"execution_time"`
	SettlementTime time.Time `json:"settlement_time"`

	// Quality metrics
	Slippage        float64 `json:"slippage"`         // Basis points
	MarketImpact    float64 `json:"market_impact"`    // Basis points
	TimingCost      float64 `json:"timing_cost"`      // Basis points
	OpportunityCost float64 `json:"opportunity_cost"` // Basis points

	// Fees and costs
	CommissionFee decimal.Decimal `json:"commission_fee"`
	ExchangeFee   decimal.Decimal `json:"exchange_fee"`
	ClearingFee   decimal.Decimal `json:"clearing_fee"`
	TotalCosts    decimal.Decimal `json:"total_costs"`

	// Benchmarks
	ArrivalPrice decimal.Decimal `json:"arrival_price"` // Price when order arrived
	TWAPPrice    decimal.Decimal `json:"twap_price"`    // TWAP during execution
	VWAPPrice    decimal.Decimal `json:"vwap_price"`    // VWAP during execution
	ClosePrice   decimal.Decimal `json:"close_price"`   // Close price

	// Metadata
	Tags  map[string]string `json:"tags,omitempty"`
	Notes string            `json:"notes,omitempty"`
}

// ExecutionMetrics contains execution quality metrics
type ExecutionMetrics struct {
	Symbol string    `json:"symbol"`
	Period time.Time `json:"period"`

	// Volume metrics
	TotalVolume  decimal.Decimal `json:"total_volume"`
	TradeCount   int             `json:"trade_count"`
	AvgTradeSize decimal.Decimal `json:"avg_trade_size"`

	// Execution quality
	FillRate    float64       `json:"fill_rate"`    // % of orders filled
	AvgSlippage float64       `json:"avg_slippage"` // Average slippage in bps
	AvgLatency  time.Duration `json:"avg_latency"`  // Average execution latency

	// Cost analysis
	AvgCommission decimal.Decimal `json:"avg_commission"` // Average commission per trade
	TotalCosts    decimal.Decimal `json:"total_costs"`    // Total transaction costs
	CostPerShare  decimal.Decimal `json:"cost_per_share"` // Cost per unit traded

	// Market impact
	AvgMarketImpact  float64 `json:"avg_market_impact"` // Average market impact in bps
	ImpactVolatility float64 `json:"impact_volatility"` // Volatility of market impact

	// Performance vs benchmarks
	VsArrival float64 `json:"vs_arrival"` // Performance vs arrival price
	VsTWAP    float64 `json:"vs_twap"`    // Performance vs TWAP
	VsVWAP    float64 `json:"vs_vwap"`    // Performance vs VWAP
	VsClose   float64 `json:"vs_close"`   // Performance vs close
}

// PerformanceData contains strategy performance data
type PerformanceData struct {
	StrategyID string    `json:"strategy_id"`
	Period     time.Time `json:"period"`

	// P&L metrics
	GrossPnL   decimal.Decimal `json:"gross_pnl"`
	NetPnL     decimal.Decimal `json:"net_pnl"`
	TotalCosts decimal.Decimal `json:"total_costs"`

	// Risk metrics
	Volatility   float64 `json:"volatility"`
	MaxDrawdown  float64 `json:"max_drawdown"`
	SharpeRatio  float64 `json:"sharpe_ratio"`
	SortinoRatio float64 `json:"sortino_ratio"`

	// Execution metrics
	TotalTrades  int             `json:"total_trades"`
	WinRate      float64         `json:"win_rate"`
	AvgWin       decimal.Decimal `json:"avg_win"`
	AvgLoss      decimal.Decimal `json:"avg_loss"`
	ProfitFactor float64         `json:"profit_factor"`

	// Attribution
	AlphaReturn  float64 `json:"alpha_return"`
	BetaExposure float64 `json:"beta_exposure"`
	MarketReturn float64 `json:"market_return"`
}

// RealtimeMetrics contains real-time analytics metrics
type RealtimeMetrics struct {
	// Current session metrics
	SessionVolume decimal.Decimal `json:"session_volume"`
	SessionTrades int             `json:"session_trades"`
	SessionPnL    decimal.Decimal `json:"session_pnl"`

	// Real-time execution quality
	CurrentSlippage float64       `json:"current_slippage"`
	CurrentLatency  time.Duration `json:"current_latency"`
	FillRateToday   float64       `json:"fill_rate_today"`

	// Market conditions
	MarketVolatility float64 `json:"market_volatility"`
	LiquidityScore   float64 `json:"liquidity_score"`

	// Alerts
	ActiveAlerts  int       `json:"active_alerts"`
	LastAlertTime time.Time `json:"last_alert_time"`

	// Performance
	LastUpdate      time.Time     `json:"last_update"`
	UpdateFrequency time.Duration `json:"update_frequency"`
}

// AlertThresholds defines thresholds for generating alerts
type AlertThresholds struct {
	MaxSlippage float64       `json:"max_slippage"`  // Maximum acceptable slippage
	MinFillRate float64       `json:"min_fill_rate"` // Minimum fill rate
	MaxLatency  time.Duration `json:"max_latency"`   // Maximum latency
	MaxDrawdown float64       `json:"max_drawdown"`  // Maximum drawdown
	MinSharpe   float64       `json:"min_sharpe"`    // Minimum Sharpe ratio
}

// AnalyticsEvent represents an analytics event
type AnalyticsEvent struct {
	ID        uuid.UUID              `json:"id"`
	Type      AnalyticsEventType     `json:"type"`
	Severity  Severity               `json:"severity"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
	Source    string                 `json:"source"`
}

// AnalyticsEventType represents different types of analytics events
type AnalyticsEventType string

const (
	AnalyticsEventTrade       AnalyticsEventType = "TRADE_EXECUTED"
	AnalyticsEventSlippage    AnalyticsEventType = "HIGH_SLIPPAGE"
	AnalyticsEventLatency     AnalyticsEventType = "HIGH_LATENCY"
	AnalyticsEventFillRate    AnalyticsEventType = "LOW_FILL_RATE"
	AnalyticsEventPerformance AnalyticsEventType = "PERFORMANCE_ALERT"
	AnalyticsEventReport      AnalyticsEventType = "REPORT_GENERATED"
)

// NewPostTradeAnalytics creates a new post-trade analytics system
func NewPostTradeAnalytics(logger *observability.Logger, config PostTradeConfig) *PostTradeAnalytics {
	// Set default values
	if config.AnalysisInterval == 0 {
		config.AnalysisInterval = time.Minute
	}
	if config.HistoryRetention == 0 {
		config.HistoryRetention = 30 * 24 * time.Hour // 30 days
	}
	if config.ReportingInterval == 0 {
		config.ReportingInterval = 24 * time.Hour // Daily
	}
	if config.SlippageThreshold == 0 {
		config.SlippageThreshold = 10.0 // 10 basis points
	}
	if config.FillRateThreshold == 0 {
		config.FillRateThreshold = 0.95 // 95%
	}
	if config.LatencyThreshold == 0 {
		config.LatencyThreshold = 100 * time.Millisecond
	}

	pta := &PostTradeAnalytics{
		logger:           logger,
		config:           config,
		tradeHistory:     make([]TradeRecord, 0),
		executionMetrics: make(map[string]*ExecutionMetrics),
		performanceData:  make(map[string]*PerformanceData),
		realtimeMetrics:  &RealtimeMetrics{},
		alertThresholds: &AlertThresholds{
			MaxSlippage: config.SlippageThreshold,
			MinFillRate: config.FillRateThreshold,
			MaxLatency:  config.LatencyThreshold,
			MaxDrawdown: 20.0, // 20%
			MinSharpe:   1.0,  // Minimum Sharpe ratio
		},
		subscribers: make(map[string][]chan *AnalyticsEvent),
		stopChan:    make(chan struct{}),
	}

	// Initialize components
	pta.executionAnalyzer = NewExecutionAnalyzer(logger, config)
	pta.slippageCalculator = NewSlippageCalculator(logger, config)
	pta.marketImpactTracker = NewMarketImpactTracker(logger, config)
	pta.performanceEngine = NewPerformanceEngine(logger, config)
	pta.reportGenerator = NewReportGenerator(logger, config)

	return pta
}

// Start begins the post-trade analytics system
func (pta *PostTradeAnalytics) Start(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&pta.isRunning, 0, 1) {
		return fmt.Errorf("post-trade analytics is already running")
	}

	pta.logger.Info(ctx, "Starting post-trade analytics system", map[string]interface{}{
		"analysis_interval":  pta.config.AnalysisInterval.String(),
		"history_retention":  pta.config.HistoryRetention.String(),
		"reporting_interval": pta.config.ReportingInterval.String(),
		"enable_tca":         pta.config.EnableTCA,
		"enable_attribution": pta.config.EnableAttribution,
	})

	// Start processing threads
	pta.wg.Add(4)
	go pta.analyzeExecutions(ctx)
	go pta.generateReports(ctx)
	go pta.monitorPerformance(ctx)
	go pta.performanceMonitor(ctx)

	pta.logger.Info(ctx, "Post-trade analytics system started successfully", nil)
	return nil
}

// Stop gracefully shuts down the post-trade analytics system
func (pta *PostTradeAnalytics) Stop(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&pta.isRunning, 1, 0) {
		return fmt.Errorf("post-trade analytics is not running")
	}

	pta.logger.Info(ctx, "Stopping post-trade analytics system", nil)

	close(pta.stopChan)
	pta.wg.Wait()

	pta.logger.Info(ctx, "Post-trade analytics system stopped", map[string]interface{}{
		"trades_analyzed":   atomic.LoadInt64(&pta.tradesAnalyzed),
		"reports_generated": atomic.LoadInt64(&pta.reportsGenerated),
		"avg_analysis_time": atomic.LoadInt64(&pta.avgAnalysisTime),
	})

	return nil
}

// RecordTrade records a trade for post-trade analysis
func (pta *PostTradeAnalytics) RecordTrade(ctx context.Context, trade *TradeRecord) error {
	if atomic.LoadInt32(&pta.isRunning) != 1 {
		return fmt.Errorf("post-trade analytics is not running")
	}

	start := time.Now()

	pta.logger.Debug(ctx, "Recording trade for analysis", map[string]interface{}{
		"trade_id": trade.TradeID.String(),
		"symbol":   trade.Symbol,
		"side":     string(trade.Side),
		"quantity": trade.Quantity.String(),
		"price":    trade.Price.String(),
		"venue":    trade.VenueID,
	})

	// Calculate execution metrics
	if err := pta.calculateExecutionMetrics(ctx, trade); err != nil {
		pta.logger.Error(ctx, "Failed to calculate execution metrics", err)
	}

	// Calculate slippage
	slippage := pta.slippageCalculator.CalculateSlippage(trade)
	trade.Slippage = slippage

	// Calculate market impact
	marketImpact := pta.marketImpactTracker.CalculateImpact(trade)
	trade.MarketImpact = marketImpact

	// Store trade record
	pta.mu.Lock()
	pta.tradeHistory = append(pta.tradeHistory, *trade)
	pta.mu.Unlock()

	// Update real-time metrics
	pta.updateRealtimeMetrics(trade)

	// Check for alerts
	pta.checkAlerts(ctx, trade)

	// Publish analytics event
	pta.publishAnalyticsEvent(ctx, AnalyticsEventTrade, SeverityLow, "Trade recorded", map[string]interface{}{
		"trade_id": trade.TradeID.String(),
		"symbol":   trade.Symbol,
		"slippage": slippage,
		"impact":   marketImpact,
	})

	// Update performance metrics
	analysisTime := time.Since(start).Nanoseconds()
	atomic.StoreInt64(&pta.avgAnalysisTime, analysisTime)
	atomic.AddInt64(&pta.tradesAnalyzed, 1)

	pta.logger.Debug(ctx, "Trade recorded successfully", map[string]interface{}{
		"trade_id":      trade.TradeID.String(),
		"analysis_time": analysisTime,
		"slippage":      slippage,
		"market_impact": marketImpact,
	})

	return nil
}

// calculateExecutionMetrics calculates execution quality metrics for a trade
func (pta *PostTradeAnalytics) calculateExecutionMetrics(ctx context.Context, trade *TradeRecord) error {
	// Calculate timing cost (simplified)
	if !trade.ArrivalPrice.IsZero() {
		priceDiff := trade.Price.Sub(trade.ArrivalPrice)
		if trade.Side == OrderSideSell {
			priceDiff = priceDiff.Neg()
		}
		timingCost := priceDiff.Div(trade.ArrivalPrice).InexactFloat64() * 10000 // Convert to basis points
		trade.TimingCost = timingCost
	}

	// Calculate opportunity cost (simplified)
	if !trade.ClosePrice.IsZero() {
		priceDiff := trade.ClosePrice.Sub(trade.Price)
		if trade.Side == OrderSideSell {
			priceDiff = priceDiff.Neg()
		}
		opportunityCost := priceDiff.Div(trade.Price).InexactFloat64() * 10000 // Convert to basis points
		trade.OpportunityCost = opportunityCost
	}

	// Calculate total costs
	trade.TotalCosts = trade.CommissionFee.Add(trade.ExchangeFee).Add(trade.ClearingFee)

	return nil
}

// updateRealtimeMetrics updates real-time analytics metrics
func (pta *PostTradeAnalytics) updateRealtimeMetrics(trade *TradeRecord) {
	pta.mu.Lock()
	defer pta.mu.Unlock()

	// Update session metrics
	pta.realtimeMetrics.SessionVolume = pta.realtimeMetrics.SessionVolume.Add(trade.Value)
	pta.realtimeMetrics.SessionTrades++

	// Update current execution quality
	pta.realtimeMetrics.CurrentSlippage = trade.Slippage
	pta.realtimeMetrics.CurrentLatency = time.Since(trade.ExecutionTime)

	// Calculate today's fill rate (simplified)
	pta.realtimeMetrics.FillRateToday = 0.95 // Would calculate from actual data

	pta.realtimeMetrics.LastUpdate = time.Now()
}

// checkAlerts checks if any alerts should be triggered
func (pta *PostTradeAnalytics) checkAlerts(ctx context.Context, trade *TradeRecord) {
	// Check slippage alert
	if trade.Slippage > pta.alertThresholds.MaxSlippage {
		pta.publishAnalyticsEvent(ctx, AnalyticsEventSlippage, SeverityMedium,
			fmt.Sprintf("High slippage detected: %.2f bps", trade.Slippage), map[string]interface{}{
				"trade_id":  trade.TradeID.String(),
				"symbol":    trade.Symbol,
				"slippage":  trade.Slippage,
				"threshold": pta.alertThresholds.MaxSlippage,
			})
	}

	// Check latency alert
	executionLatency := time.Since(trade.ExecutionTime)
	if executionLatency > pta.alertThresholds.MaxLatency {
		pta.publishAnalyticsEvent(ctx, AnalyticsEventLatency, SeverityMedium,
			fmt.Sprintf("High execution latency: %s", executionLatency.String()), map[string]interface{}{
				"trade_id":  trade.TradeID.String(),
				"symbol":    trade.Symbol,
				"latency":   executionLatency.String(),
				"threshold": pta.alertThresholds.MaxLatency.String(),
			})
	}
}

// publishAnalyticsEvent publishes an analytics event to subscribers
func (pta *PostTradeAnalytics) publishAnalyticsEvent(ctx context.Context, eventType AnalyticsEventType, severity Severity, message string, data map[string]interface{}) {
	event := &AnalyticsEvent{
		ID:        uuid.New(),
		Type:      eventType,
		Severity:  severity,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
		Source:    "PostTradeAnalytics",
	}

	pta.mu.RLock()
	defer pta.mu.RUnlock()

	// Send to event type subscribers
	if subscribers, exists := pta.subscribers[string(eventType)]; exists {
		for _, ch := range subscribers {
			select {
			case ch <- event:
			default:
				// Subscriber channel is full, skip
			}
		}
	}

	// Send to wildcard subscribers
	if subscribers, exists := pta.subscribers["*"]; exists {
		for _, ch := range subscribers {
			select {
			case ch <- event:
			default:
				// Subscriber channel is full, skip
			}
		}
	}
}

// analyzeExecutions continuously analyzes trade executions
func (pta *PostTradeAnalytics) analyzeExecutions(ctx context.Context) {
	defer pta.wg.Done()

	pta.logger.Info(ctx, "Starting execution analyzer", nil)

	ticker := time.NewTicker(pta.config.AnalysisInterval)
	defer ticker.Stop()

	for {
		select {
		case <-pta.stopChan:
			return
		case <-ticker.C:
			pta.performExecutionAnalysis(ctx)
		}
	}
}

// performExecutionAnalysis performs periodic execution analysis
func (pta *PostTradeAnalytics) performExecutionAnalysis(ctx context.Context) {
	pta.mu.RLock()
	trades := make([]TradeRecord, len(pta.tradeHistory))
	copy(trades, pta.tradeHistory)
	pta.mu.RUnlock()

	if len(trades) == 0 {
		return
	}

	// Group trades by symbol for analysis
	tradesBySymbol := make(map[string][]TradeRecord)
	for _, trade := range trades {
		tradesBySymbol[trade.Symbol] = append(tradesBySymbol[trade.Symbol], trade)
	}

	// Analyze each symbol
	for symbol, symbolTrades := range tradesBySymbol {
		metrics := pta.calculateExecutionMetricsForSymbol(symbol, symbolTrades)

		pta.mu.Lock()
		pta.executionMetrics[symbol] = metrics
		pta.mu.Unlock()

		pta.logger.Debug(ctx, "Execution metrics calculated", map[string]interface{}{
			"symbol":       symbol,
			"trade_count":  metrics.TradeCount,
			"avg_slippage": metrics.AvgSlippage,
			"fill_rate":    metrics.FillRate,
		})
	}
}

// calculateExecutionMetricsForSymbol calculates execution metrics for a symbol
func (pta *PostTradeAnalytics) calculateExecutionMetricsForSymbol(symbol string, trades []TradeRecord) *ExecutionMetrics {
	if len(trades) == 0 {
		return &ExecutionMetrics{Symbol: symbol, Period: time.Now()}
	}

	metrics := &ExecutionMetrics{
		Symbol:     symbol,
		Period:     time.Now(),
		TradeCount: len(trades),
	}

	var totalVolume decimal.Decimal
	var totalSlippage, totalMarketImpact float64
	var totalLatency time.Duration
	var totalCosts decimal.Decimal

	for _, trade := range trades {
		totalVolume = totalVolume.Add(trade.Value)
		totalSlippage += trade.Slippage
		totalMarketImpact += trade.MarketImpact
		totalLatency += time.Since(trade.ExecutionTime)
		totalCosts = totalCosts.Add(trade.TotalCosts)
	}

	metrics.TotalVolume = totalVolume
	metrics.AvgTradeSize = totalVolume.Div(decimal.NewFromInt(int64(len(trades))))
	metrics.AvgSlippage = totalSlippage / float64(len(trades))
	metrics.AvgMarketImpact = totalMarketImpact / float64(len(trades))
	metrics.AvgLatency = totalLatency / time.Duration(len(trades))
	metrics.TotalCosts = totalCosts
	metrics.CostPerShare = totalCosts.Div(totalVolume)

	// Calculate fill rate (simplified - assume all trades are filled)
	metrics.FillRate = 100.0

	return metrics
}

// generateReports continuously generates reports
func (pta *PostTradeAnalytics) generateReports(ctx context.Context) {
	defer pta.wg.Done()

	pta.logger.Info(ctx, "Starting report generator", nil)

	ticker := time.NewTicker(pta.config.ReportingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-pta.stopChan:
			return
		case <-ticker.C:
			if pta.config.AutoReporting {
				pta.generateDailyReport(ctx)
			}
		}
	}
}

// generateDailyReport generates a daily execution report
func (pta *PostTradeAnalytics) generateDailyReport(ctx context.Context) {
	pta.logger.Info(ctx, "Generating daily execution report", nil)

	report := pta.reportGenerator.GenerateDailyReport(pta.tradeHistory, pta.executionMetrics)

	// Publish report event
	pta.publishAnalyticsEvent(ctx, AnalyticsEventReport, SeverityLow, "Daily report generated", map[string]interface{}{
		"report_type": "daily_execution",
		"trade_count": len(pta.tradeHistory),
		"symbols":     len(pta.executionMetrics),
	})

	atomic.AddInt64(&pta.reportsGenerated, 1)

	pta.logger.Info(ctx, "Daily report generated successfully", map[string]interface{}{
		"report_id": report.ID,
	})
}

// monitorPerformance continuously monitors performance
func (pta *PostTradeAnalytics) monitorPerformance(ctx context.Context) {
	defer pta.wg.Done()

	pta.logger.Info(ctx, "Starting performance monitor", nil)

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-pta.stopChan:
			return
		case <-ticker.C:
			pta.updatePerformanceMetrics(ctx)
		}
	}
}

// updatePerformanceMetrics updates performance metrics for all strategies
func (pta *PostTradeAnalytics) updatePerformanceMetrics(ctx context.Context) {
	pta.mu.RLock()
	trades := make([]TradeRecord, len(pta.tradeHistory))
	copy(trades, pta.tradeHistory)
	pta.mu.RUnlock()

	// Group trades by strategy
	tradesByStrategy := make(map[string][]TradeRecord)
	for _, trade := range trades {
		if trade.StrategyID != "" {
			tradesByStrategy[trade.StrategyID] = append(tradesByStrategy[trade.StrategyID], trade)
		}
	}

	// Calculate performance for each strategy
	for strategyID, strategyTrades := range tradesByStrategy {
		performance := pta.performanceEngine.CalculatePerformance(strategyTrades, strategyID)

		pta.mu.Lock()
		pta.performanceData[strategyID] = performance
		pta.mu.Unlock()

		// Check for performance alerts
		pta.checkPerformanceAlerts(ctx, performance)
	}
}

// checkPerformanceAlerts checks for performance-related alerts
func (pta *PostTradeAnalytics) checkPerformanceAlerts(ctx context.Context, performance *PerformanceData) {
	// Check Sharpe ratio
	if performance.SharpeRatio < pta.alertThresholds.MinSharpe {
		pta.publishAnalyticsEvent(ctx, AnalyticsEventPerformance, SeverityMedium,
			fmt.Sprintf("Low Sharpe ratio for strategy %s: %.2f", performance.StrategyID, performance.SharpeRatio),
			map[string]interface{}{
				"strategy_id":  performance.StrategyID,
				"sharpe_ratio": performance.SharpeRatio,
				"threshold":    pta.alertThresholds.MinSharpe,
			})
	}

	// Check drawdown
	if performance.MaxDrawdown > pta.alertThresholds.MaxDrawdown {
		pta.publishAnalyticsEvent(ctx, AnalyticsEventPerformance, SeverityHigh,
			fmt.Sprintf("High drawdown for strategy %s: %.2f%%", performance.StrategyID, performance.MaxDrawdown),
			map[string]interface{}{
				"strategy_id":  performance.StrategyID,
				"max_drawdown": performance.MaxDrawdown,
				"threshold":    pta.alertThresholds.MaxDrawdown,
			})
	}
}

// performanceMonitor tracks and reports performance metrics
func (pta *PostTradeAnalytics) performanceMonitor(ctx context.Context) {
	defer pta.wg.Done()

	pta.logger.Info(ctx, "Starting analytics performance monitor", nil)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	var lastTradeCount int64

	for {
		select {
		case <-pta.stopChan:
			return
		case <-ticker.C:
			currentTrades := atomic.LoadInt64(&pta.tradesAnalyzed)
			tradesPerSecond := (currentTrades - lastTradeCount) / 10
			lastTradeCount = currentTrades

			reports := atomic.LoadInt64(&pta.reportsGenerated)
			avgAnalysisTime := atomic.LoadInt64(&pta.avgAnalysisTime)

			pta.logger.Info(ctx, "Post-trade analytics performance", map[string]interface{}{
				"trades_per_second": tradesPerSecond,
				"total_trades":      currentTrades,
				"reports_generated": reports,
				"avg_analysis_time": avgAnalysisTime,
				"avg_analysis_us":   avgAnalysisTime / 1000,
				"active_symbols":    len(pta.executionMetrics),
				"active_strategies": len(pta.performanceData),
			})
		}
	}
}

// Subscribe registers a subscriber for analytics events
func (pta *PostTradeAnalytics) Subscribe(eventType string) <-chan *AnalyticsEvent {
	pta.mu.Lock()
	defer pta.mu.Unlock()

	ch := make(chan *AnalyticsEvent, 1000) // Buffered channel
	if pta.subscribers[eventType] == nil {
		pta.subscribers[eventType] = make([]chan *AnalyticsEvent, 0)
	}
	pta.subscribers[eventType] = append(pta.subscribers[eventType], ch)

	return ch
}

// GetExecutionMetrics returns execution metrics for a symbol
func (pta *PostTradeAnalytics) GetExecutionMetrics(symbol string) *ExecutionMetrics {
	pta.mu.RLock()
	defer pta.mu.RUnlock()

	if metrics, exists := pta.executionMetrics[symbol]; exists {
		return metrics
	}
	return nil
}

// GetPerformanceData returns performance data for a strategy
func (pta *PostTradeAnalytics) GetPerformanceData(strategyID string) *PerformanceData {
	pta.mu.RLock()
	defer pta.mu.RUnlock()

	if performance, exists := pta.performanceData[strategyID]; exists {
		return performance
	}
	return nil
}

// GetRealtimeMetrics returns current real-time metrics
func (pta *PostTradeAnalytics) GetRealtimeMetrics() *RealtimeMetrics {
	pta.mu.RLock()
	defer pta.mu.RUnlock()

	// Return a copy of current metrics
	metrics := *pta.realtimeMetrics
	return &metrics
}

// GetTradeHistory returns trade history with optional filtering
func (pta *PostTradeAnalytics) GetTradeHistory(limit int, symbol string, strategyID string) []TradeRecord {
	pta.mu.RLock()
	defer pta.mu.RUnlock()

	trades := make([]TradeRecord, 0)

	for _, trade := range pta.tradeHistory {
		// Apply filters
		if symbol != "" && trade.Symbol != symbol {
			continue
		}
		if strategyID != "" && trade.StrategyID != strategyID {
			continue
		}

		trades = append(trades, trade)
	}

	// Sort by execution time (most recent first)
	sort.Slice(trades, func(i, j int) bool {
		return trades[i].ExecutionTime.After(trades[j].ExecutionTime)
	})

	// Apply limit
	if limit > 0 && len(trades) > limit {
		trades = trades[:limit]
	}

	return trades
}
