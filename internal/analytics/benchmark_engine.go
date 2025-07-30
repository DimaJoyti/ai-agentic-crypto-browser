package analytics

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/shopspring/decimal"
)

// BenchmarkEngine provides performance benchmarking capabilities
type BenchmarkEngine struct {
	logger     *observability.Logger
	config     PerformanceConfig
	benchmarks map[string]*Benchmark
	metrics    BenchmarkMetrics
	mu         sync.RWMutex
	isRunning  int32
}

// Benchmark represents a performance benchmark
type Benchmark struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Type        BenchmarkType          `json:"type"`
	Data        []BenchmarkDataPoint   `json:"data"`
	Metadata    map[string]interface{} `json:"metadata"`
	LastUpdated time.Time              `json:"last_updated"`
	IsActive    bool                   `json:"is_active"`
}

// BenchmarkType defines types of benchmarks
type BenchmarkType string

const (
	BenchmarkTypeMarket   BenchmarkType = "MARKET"
	BenchmarkTypeStrategy BenchmarkType = "STRATEGY"
	BenchmarkTypePeer     BenchmarkType = "PEER"
	BenchmarkTypeCustom   BenchmarkType = "CUSTOM"
	BenchmarkTypeRiskFree BenchmarkType = "RISK_FREE"
)

// BenchmarkDataPoint represents a single benchmark data point
type BenchmarkDataPoint struct {
	Timestamp time.Time       `json:"timestamp"`
	Value     decimal.Decimal `json:"value"`
	Return    decimal.Decimal `json:"return"`
	Volume    decimal.Decimal `json:"volume,omitempty"`
}

// BenchmarkComparison represents a comparison between portfolio and benchmark
type BenchmarkComparison struct {
	BenchmarkID      string          `json:"benchmark_id"`
	BenchmarkName    string          `json:"benchmark_name"`
	Period           string          `json:"period"`
	PortfolioReturn  decimal.Decimal `json:"portfolio_return"`
	BenchmarkReturn  decimal.Decimal `json:"benchmark_return"`
	ExcessReturn     decimal.Decimal `json:"excess_return"`
	TrackingError    float64         `json:"tracking_error"`
	InformationRatio float64         `json:"information_ratio"`
	Beta             float64         `json:"beta"`
	Alpha            float64         `json:"alpha"`
	Correlation      float64         `json:"correlation"`
	UpCapture        float64         `json:"up_capture"`
	DownCapture      float64         `json:"down_capture"`
	BattingAverage   float64         `json:"batting_average"`
	WinLossRatio     float64         `json:"win_loss_ratio"`
	MaxRelativeDD    decimal.Decimal `json:"max_relative_drawdown"`
	RelativeSharpe   float64         `json:"relative_sharpe"`
}

// PeerAnalysis represents peer comparison analysis
type PeerAnalysis struct {
	PeerGroup          string                `json:"peer_group"`
	TotalPeers         int                   `json:"total_peers"`
	Ranking            int                   `json:"ranking"`
	Percentile         float64               `json:"percentile"`
	Quartile           int                   `json:"quartile"`
	AboveMedian        bool                  `json:"above_median"`
	TopDecile          bool                  `json:"top_decile"`
	TopQuartile        bool                  `json:"top_quartile"`
	Metrics            map[string]PeerMetric `json:"metrics"`
	RelativeStrengths  []string              `json:"relative_strengths"`
	RelativeWeaknesses []string              `json:"relative_weaknesses"`
}

// PeerMetric represents a metric comparison with peers
type PeerMetric struct {
	Value      float64 `json:"value"`
	PeerMedian float64 `json:"peer_median"`
	PeerMean   float64 `json:"peer_mean"`
	Percentile float64 `json:"percentile"`
	ZScore     float64 `json:"z_score"`
	Ranking    int     `json:"ranking"`
}

// NewBenchmarkEngine creates a new benchmark engine
func NewBenchmarkEngine(logger *observability.Logger, config PerformanceConfig) *BenchmarkEngine {
	be := &BenchmarkEngine{
		logger:     logger,
		config:     config,
		benchmarks: make(map[string]*Benchmark),
	}

	// Initialize default benchmarks
	be.initializeDefaultBenchmarks()

	return be
}

// Start starts the benchmark engine
func (be *BenchmarkEngine) Start(ctx context.Context) error {
	be.logger.Info(ctx, "Starting benchmark engine", nil)
	be.isRunning = 1
	return nil
}

// Stop stops the benchmark engine
func (be *BenchmarkEngine) Stop(ctx context.Context) error {
	be.logger.Info(ctx, "Stopping benchmark engine", nil)
	be.isRunning = 0
	return nil
}

// initializeDefaultBenchmarks sets up default benchmarks
func (be *BenchmarkEngine) initializeDefaultBenchmarks() {
	// Bitcoin benchmark
	btcBenchmark := &Benchmark{
		ID:          "btc_index",
		Name:        "Bitcoin Index",
		Description: "Bitcoin price performance benchmark",
		Type:        BenchmarkTypeMarket,
		Data:        be.generateMockBenchmarkData("BTC", 365),
		Metadata: map[string]interface{}{
			"symbol":      "BTC",
			"asset_class": "cryptocurrency",
			"currency":    "USD",
		},
		LastUpdated: time.Now(),
		IsActive:    true,
	}

	// Crypto market index
	cryptoIndexBenchmark := &Benchmark{
		ID:          "crypto_index",
		Name:        "Crypto Market Index",
		Description: "Weighted cryptocurrency market index",
		Type:        BenchmarkTypeMarket,
		Data:        be.generateMockBenchmarkData("CRYPTO", 365),
		Metadata: map[string]interface{}{
			"composition": []string{"BTC", "ETH", "BNB", "ADA", "SOL"},
			"asset_class": "cryptocurrency",
			"currency":    "USD",
		},
		LastUpdated: time.Now(),
		IsActive:    true,
	}

	// Risk-free rate benchmark
	riskFreeBenchmark := &Benchmark{
		ID:          "risk_free",
		Name:        "Risk-Free Rate",
		Description: "US Treasury 3-month rate",
		Type:        BenchmarkTypeRiskFree,
		Data:        be.generateRiskFreeData(365),
		Metadata: map[string]interface{}{
			"instrument": "US Treasury 3M",
			"currency":   "USD",
			"rate":       0.02,
		},
		LastUpdated: time.Now(),
		IsActive:    true,
	}

	be.benchmarks["btc_index"] = btcBenchmark
	be.benchmarks["crypto_index"] = cryptoIndexBenchmark
	be.benchmarks["risk_free"] = riskFreeBenchmark
}

// generateMockBenchmarkData generates mock benchmark data
func (be *BenchmarkEngine) generateMockBenchmarkData(symbol string, days int) []BenchmarkDataPoint {
	data := make([]BenchmarkDataPoint, days)
	basePrice := decimal.NewFromFloat(50000) // Starting price

	for i := 0; i < days; i++ {
		timestamp := time.Now().AddDate(0, 0, -days+i)

		// Generate random price movement
		change := (math.Sin(float64(i)*0.1) + math.Cos(float64(i)*0.05)) * 0.02
		if symbol == "CRYPTO" {
			change *= 1.5 // Higher volatility for crypto index
		}

		if i == 0 {
			data[i] = BenchmarkDataPoint{
				Timestamp: timestamp,
				Value:     basePrice,
				Return:    decimal.Zero,
				Volume:    decimal.NewFromFloat(1000000),
			}
		} else {
			newPrice := data[i-1].Value.Mul(decimal.NewFromFloat(1 + change))
			returnPct := newPrice.Sub(data[i-1].Value).Div(data[i-1].Value)

			data[i] = BenchmarkDataPoint{
				Timestamp: timestamp,
				Value:     newPrice,
				Return:    returnPct,
				Volume:    decimal.NewFromFloat(1000000 * (1 + math.Abs(change))),
			}
		}
	}

	return data
}

// generateRiskFreeData generates risk-free rate data
func (be *BenchmarkEngine) generateRiskFreeData(days int) []BenchmarkDataPoint {
	data := make([]BenchmarkDataPoint, days)
	annualRate := 0.02 // 2% annual rate
	dailyRate := annualRate / 365

	for i := 0; i < days; i++ {
		timestamp := time.Now().AddDate(0, 0, -days+i)

		data[i] = BenchmarkDataPoint{
			Timestamp: timestamp,
			Value:     decimal.NewFromFloat(1.0 + float64(i)*dailyRate),
			Return:    decimal.NewFromFloat(dailyRate),
		}
	}

	return data
}

// CompareToBenchmark compares portfolio performance to a benchmark
func (be *BenchmarkEngine) CompareToBenchmark(portfolioReturns []decimal.Decimal, benchmarkID string, period string) (*BenchmarkComparison, error) {
	be.mu.RLock()
	benchmark, exists := be.benchmarks[benchmarkID]
	be.mu.RUnlock()

	if !exists {
		return nil, ErrBenchmarkNotFound
	}

	// Extract benchmark returns for the same period
	benchmarkReturns := be.extractBenchmarkReturns(benchmark, len(portfolioReturns))

	if len(portfolioReturns) != len(benchmarkReturns) {
		return nil, ErrMismatchedDataLength
	}

	// Calculate comparison metrics
	portfolioReturn := be.calculateTotalReturn(portfolioReturns)
	benchmarkReturn := be.calculateTotalReturn(benchmarkReturns)
	excessReturn := portfolioReturn.Sub(benchmarkReturn)

	trackingError := be.calculateTrackingError(portfolioReturns, benchmarkReturns)
	informationRatio := be.calculateInformationRatio(excessReturn, trackingError)
	beta := be.calculateBeta(portfolioReturns, benchmarkReturns)
	alpha := be.calculateAlpha(portfolioReturns, benchmarkReturns, beta)
	correlation := be.calculateCorrelation(portfolioReturns, benchmarkReturns)
	upCapture := be.calculateUpCapture(portfolioReturns, benchmarkReturns)
	downCapture := be.calculateDownCapture(portfolioReturns, benchmarkReturns)
	battingAverage := be.calculateBattingAverage(portfolioReturns, benchmarkReturns)
	winLossRatio := be.calculateWinLossRatio(portfolioReturns, benchmarkReturns)
	maxRelativeDD := be.calculateMaxRelativeDrawdown(portfolioReturns, benchmarkReturns)
	relativeSharpe := be.calculateRelativeSharpe(portfolioReturns, benchmarkReturns)

	return &BenchmarkComparison{
		BenchmarkID:      benchmarkID,
		BenchmarkName:    benchmark.Name,
		Period:           period,
		PortfolioReturn:  portfolioReturn,
		BenchmarkReturn:  benchmarkReturn,
		ExcessReturn:     excessReturn,
		TrackingError:    trackingError,
		InformationRatio: informationRatio,
		Beta:             beta,
		Alpha:            alpha,
		Correlation:      correlation,
		UpCapture:        upCapture,
		DownCapture:      downCapture,
		BattingAverage:   battingAverage,
		WinLossRatio:     winLossRatio,
		MaxRelativeDD:    maxRelativeDD,
		RelativeSharpe:   relativeSharpe,
	}, nil
}

// extractBenchmarkReturns extracts returns from benchmark data
func (be *BenchmarkEngine) extractBenchmarkReturns(benchmark *Benchmark, count int) []decimal.Decimal {
	if len(benchmark.Data) == 0 {
		return []decimal.Decimal{}
	}

	// Get the most recent data points
	start := len(benchmark.Data) - count
	if start < 0 {
		start = 0
	}

	returns := make([]decimal.Decimal, 0, count)
	for i := start; i < len(benchmark.Data); i++ {
		returns = append(returns, benchmark.Data[i].Return)
	}

	return returns
}

// calculateTotalReturn calculates total return from a series of returns
func (be *BenchmarkEngine) calculateTotalReturn(returns []decimal.Decimal) decimal.Decimal {
	if len(returns) == 0 {
		return decimal.Zero
	}

	totalReturn := decimal.NewFromFloat(1.0)
	for _, ret := range returns {
		totalReturn = totalReturn.Mul(decimal.NewFromFloat(1.0).Add(ret))
	}

	return totalReturn.Sub(decimal.NewFromFloat(1.0))
}

// calculateTrackingError calculates tracking error between portfolio and benchmark
func (be *BenchmarkEngine) calculateTrackingError(portfolioReturns, benchmarkReturns []decimal.Decimal) float64 {
	if len(portfolioReturns) != len(benchmarkReturns) || len(portfolioReturns) < 2 {
		return 0
	}

	// Calculate excess returns
	excessReturns := make([]float64, len(portfolioReturns))
	for i := range portfolioReturns {
		excess := portfolioReturns[i].Sub(benchmarkReturns[i])
		excessReturns[i], _ = excess.Float64()
	}

	// Calculate standard deviation of excess returns
	var sum float64
	for _, ret := range excessReturns {
		sum += ret
	}
	mean := sum / float64(len(excessReturns))

	var variance float64
	for _, ret := range excessReturns {
		variance += math.Pow(ret-mean, 2)
	}
	variance /= float64(len(excessReturns) - 1)

	return math.Sqrt(variance) * math.Sqrt(252) // Annualized
}

// calculateInformationRatio calculates information ratio
func (be *BenchmarkEngine) calculateInformationRatio(excessReturn decimal.Decimal, trackingError float64) float64 {
	if trackingError == 0 {
		return 0
	}

	excessReturnFloat, _ := excessReturn.Float64()
	return excessReturnFloat / trackingError
}

// calculateBeta calculates beta relative to benchmark
func (be *BenchmarkEngine) calculateBeta(portfolioReturns, benchmarkReturns []decimal.Decimal) float64 {
	if len(portfolioReturns) != len(benchmarkReturns) || len(portfolioReturns) < 2 {
		return 1.0
	}

	// Convert to float64 arrays
	portRets := make([]float64, len(portfolioReturns))
	benchRets := make([]float64, len(benchmarkReturns))

	for i := range portfolioReturns {
		portRets[i], _ = portfolioReturns[i].Float64()
		benchRets[i], _ = benchmarkReturns[i].Float64()
	}

	// Calculate means
	var portSum, benchSum float64
	for i := range portRets {
		portSum += portRets[i]
		benchSum += benchRets[i]
	}
	portMean := portSum / float64(len(portRets))
	benchMean := benchSum / float64(len(benchRets))

	// Calculate covariance and variance
	var covariance, benchVariance float64
	for i := range portRets {
		covariance += (portRets[i] - portMean) * (benchRets[i] - benchMean)
		benchVariance += math.Pow(benchRets[i]-benchMean, 2)
	}

	if benchVariance == 0 {
		return 1.0
	}

	return covariance / benchVariance
}

// calculateAlpha calculates alpha relative to benchmark
func (be *BenchmarkEngine) calculateAlpha(portfolioReturns, benchmarkReturns []decimal.Decimal, beta float64) float64 {
	if len(portfolioReturns) == 0 || len(benchmarkReturns) == 0 {
		return 0
	}

	portfolioReturn := be.calculateTotalReturn(portfolioReturns)
	benchmarkReturn := be.calculateTotalReturn(benchmarkReturns)
	riskFreeRate := 0.02 / 252 * float64(len(portfolioReturns)) // Approximate risk-free return

	portRetFloat, _ := portfolioReturn.Float64()
	benchRetFloat, _ := benchmarkReturn.Float64()

	return portRetFloat - riskFreeRate - beta*(benchRetFloat-riskFreeRate)
}

// calculateCorrelation calculates correlation between portfolio and benchmark
func (be *BenchmarkEngine) calculateCorrelation(portfolioReturns, benchmarkReturns []decimal.Decimal) float64 {
	if len(portfolioReturns) != len(benchmarkReturns) || len(portfolioReturns) < 2 {
		return 0
	}

	// Convert to float64 arrays
	portRets := make([]float64, len(portfolioReturns))
	benchRets := make([]float64, len(benchmarkReturns))

	for i := range portfolioReturns {
		portRets[i], _ = portfolioReturns[i].Float64()
		benchRets[i], _ = benchmarkReturns[i].Float64()
	}

	// Calculate means
	var portSum, benchSum float64
	for i := range portRets {
		portSum += portRets[i]
		benchSum += benchRets[i]
	}
	portMean := portSum / float64(len(portRets))
	benchMean := benchSum / float64(len(benchRets))

	// Calculate correlation
	var numerator, portSumSq, benchSumSq float64
	for i := range portRets {
		portDiff := portRets[i] - portMean
		benchDiff := benchRets[i] - benchMean
		numerator += portDiff * benchDiff
		portSumSq += portDiff * portDiff
		benchSumSq += benchDiff * benchDiff
	}

	denominator := math.Sqrt(portSumSq * benchSumSq)
	if denominator == 0 {
		return 0
	}

	return numerator / denominator
}

// calculateUpCapture calculates up capture ratio
func (be *BenchmarkEngine) calculateUpCapture(portfolioReturns, benchmarkReturns []decimal.Decimal) float64 {
	if len(portfolioReturns) != len(benchmarkReturns) {
		return 0
	}

	var upPortfolio, upBenchmark decimal.Decimal
	upPeriods := 0

	for i := range benchmarkReturns {
		if benchmarkReturns[i].GreaterThan(decimal.Zero) {
			upPortfolio = upPortfolio.Add(portfolioReturns[i])
			upBenchmark = upBenchmark.Add(benchmarkReturns[i])
			upPeriods++
		}
	}

	if upPeriods == 0 || upBenchmark.IsZero() {
		return 0
	}

	upCaptureRatio := upPortfolio.Div(upBenchmark)
	result, _ := upCaptureRatio.Float64()
	return result * 100 // Return as percentage
}

// calculateDownCapture calculates down capture ratio
func (be *BenchmarkEngine) calculateDownCapture(portfolioReturns, benchmarkReturns []decimal.Decimal) float64 {
	if len(portfolioReturns) != len(benchmarkReturns) {
		return 0
	}

	var downPortfolio, downBenchmark decimal.Decimal
	downPeriods := 0

	for i := range benchmarkReturns {
		if benchmarkReturns[i].LessThan(decimal.Zero) {
			downPortfolio = downPortfolio.Add(portfolioReturns[i])
			downBenchmark = downBenchmark.Add(benchmarkReturns[i])
			downPeriods++
		}
	}

	if downPeriods == 0 || downBenchmark.IsZero() {
		return 0
	}

	downCaptureRatio := downPortfolio.Div(downBenchmark)
	result, _ := downCaptureRatio.Float64()
	return result * 100 // Return as percentage
}

// calculateBattingAverage calculates batting average (% of periods outperforming benchmark)
func (be *BenchmarkEngine) calculateBattingAverage(portfolioReturns, benchmarkReturns []decimal.Decimal) float64 {
	if len(portfolioReturns) != len(benchmarkReturns) || len(portfolioReturns) == 0 {
		return 0
	}

	outperformingPeriods := 0
	for i := range portfolioReturns {
		if portfolioReturns[i].GreaterThan(benchmarkReturns[i]) {
			outperformingPeriods++
		}
	}

	return float64(outperformingPeriods) / float64(len(portfolioReturns)) * 100
}

// calculateWinLossRatio calculates win/loss ratio
func (be *BenchmarkEngine) calculateWinLossRatio(portfolioReturns, benchmarkReturns []decimal.Decimal) float64 {
	if len(portfolioReturns) != len(benchmarkReturns) {
		return 0
	}

	var totalWins, totalLosses decimal.Decimal

	for i := range portfolioReturns {
		excess := portfolioReturns[i].Sub(benchmarkReturns[i])
		if excess.GreaterThan(decimal.Zero) {
			totalWins = totalWins.Add(excess)
		} else if excess.LessThan(decimal.Zero) {
			totalLosses = totalLosses.Add(excess.Abs())
		}
	}

	if totalLosses.IsZero() {
		return math.Inf(1)
	}

	ratio := totalWins.Div(totalLosses)
	result, _ := ratio.Float64()
	return result
}

// calculateMaxRelativeDrawdown calculates maximum relative drawdown vs benchmark
func (be *BenchmarkEngine) calculateMaxRelativeDrawdown(portfolioReturns, benchmarkReturns []decimal.Decimal) decimal.Decimal {
	if len(portfolioReturns) != len(benchmarkReturns) {
		return decimal.Zero
	}

	var portfolioCumulative, benchmarkCumulative decimal.Decimal = decimal.NewFromFloat(1.0), decimal.NewFromFloat(1.0)
	var maxRelativeValue, maxRelativeDD decimal.Decimal

	for i := range portfolioReturns {
		portfolioCumulative = portfolioCumulative.Mul(decimal.NewFromFloat(1.0).Add(portfolioReturns[i]))
		benchmarkCumulative = benchmarkCumulative.Mul(decimal.NewFromFloat(1.0).Add(benchmarkReturns[i]))

		relativeValue := portfolioCumulative.Div(benchmarkCumulative)

		if relativeValue.GreaterThan(maxRelativeValue) {
			maxRelativeValue = relativeValue
		}

		relativeDD := maxRelativeValue.Sub(relativeValue).Div(maxRelativeValue)
		if relativeDD.GreaterThan(maxRelativeDD) {
			maxRelativeDD = relativeDD
		}
	}

	return maxRelativeDD
}

// calculateRelativeSharpe calculates relative Sharpe ratio
func (be *BenchmarkEngine) calculateRelativeSharpe(portfolioReturns, benchmarkReturns []decimal.Decimal) float64 {
	// Simplified calculation - would need proper Sharpe ratio calculation for both
	return 0.5 // Mock value
}

// GetMetrics returns current benchmark metrics
func (be *BenchmarkEngine) GetMetrics() BenchmarkMetrics {
	be.mu.RLock()
	defer be.mu.RUnlock()
	return be.metrics
}

// GetBenchmarks returns available benchmarks
func (be *BenchmarkEngine) GetBenchmarks() map[string]*Benchmark {
	be.mu.RLock()
	defer be.mu.RUnlock()

	benchmarks := make(map[string]*Benchmark)
	for k, v := range be.benchmarks {
		benchmarks[k] = v
	}
	return benchmarks
}

// Error definitions
var (
	ErrBenchmarkNotFound    = fmt.Errorf("benchmark not found")
	ErrMismatchedDataLength = fmt.Errorf("mismatched data length")
)
