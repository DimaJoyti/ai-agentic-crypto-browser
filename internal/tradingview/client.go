package tradingview

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/chromedp/chromedp"
	"github.com/shopspring/decimal"
)

// TradingViewClient provides TradingView integration for technical analysis
type TradingViewClient struct {
	logger  *observability.Logger
	config  Config
	browser context.Context
	cancel  context.CancelFunc

	// Chart data
	charts     map[string]*Chart
	indicators map[string]*Indicator
	signals    map[string]*Signal

	// State management
	isRunning bool
	stopChan  chan struct{}
	wg        sync.WaitGroup
	mu        sync.RWMutex
}

// Config contains TradingView configuration
type Config struct {
	BaseURL          string        `json:"base_url"`
	Username         string        `json:"username"`
	Password         string        `json:"password"`
	Headless         bool          `json:"headless"`
	Timeout          time.Duration `json:"timeout"`
	UpdateInterval   time.Duration `json:"update_interval"`
	EnableSignals    bool          `json:"enable_signals"`
	EnableIndicators bool          `json:"enable_indicators"`
}

// Chart represents a TradingView chart
type Chart struct {
	ID         string                 `json:"id"`
	Symbol     string                 `json:"symbol"`
	Timeframe  string                 `json:"timeframe"`
	URL        string                 `json:"url"`
	LastUpdate time.Time              `json:"last_update"`
	OHLCV      []OHLCV                `json:"ohlcv"`
	Indicators map[string]*Indicator  `json:"indicators"`
	Signals    []Signal               `json:"signals"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// OHLCV represents candlestick data
type OHLCV struct {
	Timestamp time.Time       `json:"timestamp"`
	Open      decimal.Decimal `json:"open"`
	High      decimal.Decimal `json:"high"`
	Low       decimal.Decimal `json:"low"`
	Close     decimal.Decimal `json:"close"`
	Volume    decimal.Decimal `json:"volume"`
}

// Indicator represents a technical indicator
type Indicator struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Type       IndicatorType          `json:"type"`
	Parameters map[string]interface{} `json:"parameters"`
	Values     []IndicatorValue       `json:"values"`
	LastUpdate time.Time              `json:"last_update"`
}

// IndicatorType represents different types of indicators
type IndicatorType string

const (
	IndicatorTypeMA     IndicatorType = "MA"     // Moving Average
	IndicatorTypeEMA    IndicatorType = "EMA"    // Exponential Moving Average
	IndicatorTypeRSI    IndicatorType = "RSI"    // Relative Strength Index
	IndicatorTypeMACD   IndicatorType = "MACD"   // MACD
	IndicatorTypeBB     IndicatorType = "BB"     // Bollinger Bands
	IndicatorTypeStoch  IndicatorType = "STOCH"  // Stochastic
	IndicatorTypeADX    IndicatorType = "ADX"    // Average Directional Index
	IndicatorTypeATR    IndicatorType = "ATR"    // Average True Range
	IndicatorTypeVolume IndicatorType = "VOLUME" // Volume
	IndicatorTypeVWAP   IndicatorType = "VWAP"   // Volume Weighted Average Price
)

// IndicatorValue represents an indicator value at a specific time
type IndicatorValue struct {
	Timestamp time.Time                  `json:"timestamp"`
	Value     decimal.Decimal            `json:"value"`
	Values    map[string]decimal.Decimal `json:"values,omitempty"` // For multi-value indicators like MACD
}

// Signal represents a trading signal from TradingView
type Signal struct {
	ID          string                 `json:"id"`
	Symbol      string                 `json:"symbol"`
	Type        SignalType             `json:"type"`
	Direction   SignalDirection        `json:"direction"`
	Strength    SignalStrength         `json:"strength"`
	Price       decimal.Decimal        `json:"price"`
	Timestamp   time.Time              `json:"timestamp"`
	Timeframe   string                 `json:"timeframe"`
	Indicators  []string               `json:"indicators"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// SignalType represents different types of signals
type SignalType string

const (
	SignalTypeTechnical   SignalType = "TECHNICAL"
	SignalTypeFundamental SignalType = "FUNDAMENTAL"
	SignalTypePattern     SignalType = "PATTERN"
	SignalTypeVolume      SignalType = "VOLUME"
	SignalTypeMomentum    SignalType = "MOMENTUM"
)

// SignalDirection represents signal direction
type SignalDirection string

const (
	SignalDirectionBuy     SignalDirection = "BUY"
	SignalDirectionSell    SignalDirection = "SELL"
	SignalDirectionHold    SignalDirection = "HOLD"
	SignalDirectionNeutral SignalDirection = "NEUTRAL"
)

// SignalStrength represents signal strength
type SignalStrength string

const (
	SignalStrengthWeak       SignalStrength = "WEAK"
	SignalStrengthModerate   SignalStrength = "MODERATE"
	SignalStrengthStrong     SignalStrength = "STRONG"
	SignalStrengthVeryStrong SignalStrength = "VERY_STRONG"
)

// NewTradingViewClient creates a new TradingView client
func NewTradingViewClient(logger *observability.Logger, config Config) *TradingViewClient {
	if config.BaseURL == "" {
		config.BaseURL = "https://www.tradingview.com"
	}

	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	if config.UpdateInterval == 0 {
		config.UpdateInterval = 5 * time.Second
	}

	return &TradingViewClient{
		logger:     logger,
		config:     config,
		charts:     make(map[string]*Chart),
		indicators: make(map[string]*Indicator),
		signals:    make(map[string]*Signal),
		stopChan:   make(chan struct{}),
	}
}

// Start begins the TradingView client
func (tv *TradingViewClient) Start(ctx context.Context) error {
	tv.mu.Lock()
	defer tv.mu.Unlock()

	if tv.isRunning {
		return fmt.Errorf("TradingView client is already running")
	}

	tv.logger.Info(ctx, "Starting TradingView client", map[string]interface{}{
		"base_url":          tv.config.BaseURL,
		"headless":          tv.config.Headless,
		"enable_signals":    tv.config.EnableSignals,
		"enable_indicators": tv.config.EnableIndicators,
	})

	// Create browser context
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", tv.config.Headless),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	allocCtx, _ := chromedp.NewExecAllocator(ctx, opts...)
	tv.browser, tv.cancel = chromedp.NewContext(allocCtx)

	// Test browser connection
	if err := tv.testConnection(ctx); err != nil {
		tv.cancel()
		return fmt.Errorf("failed to test browser connection: %w", err)
	}

	tv.isRunning = true

	// Start monitoring goroutines
	tv.wg.Add(2)
	go tv.updateCharts(ctx)
	go tv.extractSignals(ctx)

	tv.logger.Info(ctx, "TradingView client started successfully", nil)

	return nil
}

// Stop gracefully shuts down the TradingView client
func (tv *TradingViewClient) Stop(ctx context.Context) error {
	tv.mu.Lock()
	defer tv.mu.Unlock()

	if !tv.isRunning {
		return fmt.Errorf("TradingView client is not running")
	}

	tv.logger.Info(ctx, "Stopping TradingView client", nil)

	tv.isRunning = false
	close(tv.stopChan)

	// Cancel browser context
	if tv.cancel != nil {
		tv.cancel()
	}

	tv.wg.Wait()

	tv.logger.Info(ctx, "TradingView client stopped successfully", nil)

	return nil
}

// AddChart adds a chart for monitoring
func (tv *TradingViewClient) AddChart(symbol, timeframe string) (*Chart, error) {
	tv.mu.Lock()
	defer tv.mu.Unlock()

	chartID := fmt.Sprintf("%s_%s", symbol, timeframe)

	if _, exists := tv.charts[chartID]; exists {
		return nil, fmt.Errorf("chart already exists: %s", chartID)
	}

	chart := &Chart{
		ID:         chartID,
		Symbol:     symbol,
		Timeframe:  timeframe,
		URL:        fmt.Sprintf("%s/chart/?symbol=%s", tv.config.BaseURL, symbol),
		LastUpdate: time.Now(),
		OHLCV:      make([]OHLCV, 0),
		Indicators: make(map[string]*Indicator),
		Signals:    make([]Signal, 0),
		Metadata:   make(map[string]interface{}),
	}

	tv.charts[chartID] = chart

	tv.logger.Info(context.Background(), "Chart added", map[string]interface{}{
		"chart_id":  chartID,
		"symbol":    symbol,
		"timeframe": timeframe,
		"url":       chart.URL,
	})

	return chart, nil
}

// GetChart retrieves a chart by ID
func (tv *TradingViewClient) GetChart(chartID string) (*Chart, error) {
	tv.mu.RLock()
	defer tv.mu.RUnlock()

	chart, exists := tv.charts[chartID]
	if !exists {
		return nil, fmt.Errorf("chart not found: %s", chartID)
	}

	return chart, nil
}

// GetSignals retrieves signals for a symbol
func (tv *TradingViewClient) GetSignals(symbol string) []Signal {
	tv.mu.RLock()
	defer tv.mu.RUnlock()

	var signals []Signal
	for _, signal := range tv.signals {
		if signal.Symbol == symbol {
			signals = append(signals, *signal)
		}
	}

	return signals
}

// GetIndicatorValue gets the latest value for an indicator
func (tv *TradingViewClient) GetIndicatorValue(symbol, indicatorType string) (*IndicatorValue, error) {
	tv.mu.RLock()
	defer tv.mu.RUnlock()

	indicatorID := fmt.Sprintf("%s_%s", symbol, indicatorType)
	indicator, exists := tv.indicators[indicatorID]
	if !exists {
		return nil, fmt.Errorf("indicator not found: %s", indicatorID)
	}

	if len(indicator.Values) == 0 {
		return nil, fmt.Errorf("no values available for indicator: %s", indicatorID)
	}

	// Return the latest value
	return &indicator.Values[len(indicator.Values)-1], nil
}

// testConnection tests the browser connection to TradingView
func (tv *TradingViewClient) testConnection(ctx context.Context) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, tv.config.Timeout)
	defer cancel()

	var title string
	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(tv.config.BaseURL),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.Title(&title),
	)

	if err != nil {
		return fmt.Errorf("failed to navigate to TradingView: %w", err)
	}

	if !strings.Contains(strings.ToLower(title), "tradingview") {
		return fmt.Errorf("unexpected page title: %s", title)
	}

	tv.logger.Info(ctx, "TradingView connection test successful", map[string]interface{}{
		"title": title,
	})

	return nil
}

// updateCharts periodically updates chart data
func (tv *TradingViewClient) updateCharts(ctx context.Context) {
	defer tv.wg.Done()

	ticker := time.NewTicker(tv.config.UpdateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-tv.stopChan:
			return
		case <-ticker.C:
			tv.performChartUpdates(ctx)
		}
	}
}

// extractSignals periodically extracts trading signals
func (tv *TradingViewClient) extractSignals(ctx context.Context) {
	defer tv.wg.Done()

	if !tv.config.EnableSignals {
		return
	}

	ticker := time.NewTicker(tv.config.UpdateInterval * 2)
	defer ticker.Stop()

	for {
		select {
		case <-tv.stopChan:
			return
		case <-ticker.C:
			tv.performSignalExtraction(ctx)
		}
	}
}

// performChartUpdates updates all charts
func (tv *TradingViewClient) performChartUpdates(ctx context.Context) {
	tv.mu.RLock()
	charts := make([]*Chart, 0, len(tv.charts))
	for _, chart := range tv.charts {
		charts = append(charts, chart)
	}
	tv.mu.RUnlock()

	for _, chart := range charts {
		if err := tv.updateChart(ctx, chart); err != nil {
			tv.logger.Error(ctx, "Failed to update chart", err, map[string]interface{}{
				"chart_id": chart.ID,
				"symbol":   chart.Symbol,
			})
		}
	}
}

// updateChart updates a specific chart
func (tv *TradingViewClient) updateChart(ctx context.Context, chart *Chart) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, tv.config.Timeout)
	defer cancel()

	// Navigate to chart
	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(chart.URL),
		chromedp.WaitVisible("div[data-name='legend-source-item']", chromedp.ByQuery),
		chromedp.Sleep(2*time.Second), // Wait for chart to load
	)

	if err != nil {
		return fmt.Errorf("failed to navigate to chart: %w", err)
	}

	// Extract price data
	if err := tv.extractPriceData(timeoutCtx, chart); err != nil {
		tv.logger.Warn(ctx, "Failed to extract price data", map[string]interface{}{
			"chart_id": chart.ID,
			"error":    err.Error(),
		})
	}

	// Extract indicators if enabled
	if tv.config.EnableIndicators {
		if err := tv.extractIndicators(timeoutCtx, chart); err != nil {
			tv.logger.Warn(ctx, "Failed to extract indicators", map[string]interface{}{
				"chart_id": chart.ID,
				"error":    err.Error(),
			})
		}
	}

	chart.LastUpdate = time.Now()

	return nil
}

// extractPriceData extracts OHLCV data from the chart
func (tv *TradingViewClient) extractPriceData(ctx context.Context, chart *Chart) error {
	// This is a simplified implementation
	// In a real implementation, you would extract actual OHLCV data from the chart

	var priceText string
	err := chromedp.Run(ctx,
		chromedp.Text("div[data-name='legend-source-item'] > div:nth-child(2)", &priceText, chromedp.ByQuery),
	)

	if err != nil {
		return fmt.Errorf("failed to extract price text: %w", err)
	}

	// Parse price (simplified)
	if price, err := decimal.NewFromString(strings.TrimSpace(priceText)); err == nil {
		ohlcv := OHLCV{
			Timestamp: time.Now(),
			Open:      price,
			High:      price,
			Low:       price,
			Close:     price,
			Volume:    decimal.Zero,
		}

		chart.OHLCV = append(chart.OHLCV, ohlcv)

		// Keep only last 1000 candles
		if len(chart.OHLCV) > 1000 {
			chart.OHLCV = chart.OHLCV[len(chart.OHLCV)-1000:]
		}
	}

	return nil
}

// extractIndicators extracts technical indicators from the chart
func (tv *TradingViewClient) extractIndicators(ctx context.Context, chart *Chart) error {
	// This is a simplified implementation
	// In a real implementation, you would extract actual indicator values

	// Try to extract RSI value
	var rsiText string
	err := chromedp.Run(ctx,
		chromedp.Text("div[data-name='legend-source-item']:contains('RSI')", &rsiText, chromedp.ByQuery),
	)

	if err == nil && rsiText != "" {
		// Parse RSI value (simplified)
		parts := strings.Fields(rsiText)
		for _, part := range parts {
			if val, err := strconv.ParseFloat(part, 64); err == nil && val >= 0 && val <= 100 {
				indicatorID := fmt.Sprintf("%s_RSI", chart.Symbol)

				tv.mu.Lock()
				if tv.indicators[indicatorID] == nil {
					tv.indicators[indicatorID] = &Indicator{
						ID:     indicatorID,
						Name:   "RSI",
						Type:   IndicatorTypeRSI,
						Values: make([]IndicatorValue, 0),
					}
				}

				indicator := tv.indicators[indicatorID]
				indicator.Values = append(indicator.Values, IndicatorValue{
					Timestamp: time.Now(),
					Value:     decimal.NewFromFloat(val),
				})

				// Keep only last 1000 values
				if len(indicator.Values) > 1000 {
					indicator.Values = indicator.Values[len(indicator.Values)-1000:]
				}

				indicator.LastUpdate = time.Now()
				tv.mu.Unlock()

				break
			}
		}
	}

	return nil
}

// performSignalExtraction extracts trading signals
func (tv *TradingViewClient) performSignalExtraction(ctx context.Context) {
	tv.mu.RLock()
	charts := make([]*Chart, 0, len(tv.charts))
	for _, chart := range tv.charts {
		charts = append(charts, chart)
	}
	tv.mu.RUnlock()

	for _, chart := range charts {
		signals := tv.generateSignalsFromChart(chart)

		tv.mu.Lock()
		for _, signal := range signals {
			signalID := fmt.Sprintf("%s_%d", signal.Symbol, signal.Timestamp.Unix())
			tv.signals[signalID] = &signal
		}
		tv.mu.Unlock()
	}
}

// generateSignalsFromChart generates trading signals from chart data
func (tv *TradingViewClient) generateSignalsFromChart(chart *Chart) []Signal {
	var signals []Signal

	// Simple RSI-based signal generation
	rsiIndicatorID := fmt.Sprintf("%s_RSI", chart.Symbol)

	tv.mu.RLock()
	rsiIndicator, exists := tv.indicators[rsiIndicatorID]
	tv.mu.RUnlock()

	if exists && len(rsiIndicator.Values) > 0 {
		latestRSI := rsiIndicator.Values[len(rsiIndicator.Values)-1]

		var direction SignalDirection
		var strength SignalStrength

		rsiValue := latestRSI.Value

		if rsiValue.LessThan(decimal.NewFromInt(30)) {
			direction = SignalDirectionBuy
			strength = SignalStrengthStrong
		} else if rsiValue.GreaterThan(decimal.NewFromInt(70)) {
			direction = SignalDirectionSell
			strength = SignalStrengthStrong
		} else if rsiValue.LessThan(decimal.NewFromInt(40)) {
			direction = SignalDirectionBuy
			strength = SignalStrengthModerate
		} else if rsiValue.GreaterThan(decimal.NewFromInt(60)) {
			direction = SignalDirectionSell
			strength = SignalStrengthModerate
		} else {
			direction = SignalDirectionHold
			strength = SignalStrengthWeak
		}

		if direction != SignalDirectionHold {
			signal := Signal{
				ID:          fmt.Sprintf("%s_RSI_%d", chart.Symbol, time.Now().Unix()),
				Symbol:      chart.Symbol,
				Type:        SignalTypeTechnical,
				Direction:   direction,
				Strength:    strength,
				Price:       decimal.Zero, // Would be extracted from chart
				Timestamp:   time.Now(),
				Timeframe:   chart.Timeframe,
				Indicators:  []string{"RSI"},
				Description: fmt.Sprintf("RSI %s signal (RSI: %s)", direction, rsiValue.String()),
				Metadata: map[string]interface{}{
					"rsi_value": rsiValue.String(),
				},
			}

			signals = append(signals, signal)
		}
	}

	return signals
}

// GetMetrics returns client metrics
func (tv *TradingViewClient) GetMetrics() TradingViewMetrics {
	tv.mu.RLock()
	defer tv.mu.RUnlock()

	return TradingViewMetrics{
		IsRunning:       tv.isRunning,
		ChartsCount:     len(tv.charts),
		IndicatorsCount: len(tv.indicators),
		SignalsCount:    len(tv.signals),
		BrowserActive:   tv.browser != nil,
	}
}

// TradingViewMetrics contains client performance metrics
type TradingViewMetrics struct {
	IsRunning       bool `json:"is_running"`
	ChartsCount     int  `json:"charts_count"`
	IndicatorsCount int  `json:"indicators_count"`
	SignalsCount    int  `json:"signals_count"`
	BrowserActive   bool `json:"browser_active"`
}
