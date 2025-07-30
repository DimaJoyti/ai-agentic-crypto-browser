package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ai-agentic-browser/internal/analytics"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/gorilla/mux"
)

// AnalyticsHandlers handles performance analytics API endpoints
type AnalyticsHandlers struct {
	logger            *observability.Logger
	performanceEngine *analytics.PerformanceEngine
}

// NewAnalyticsHandlers creates new analytics handlers
func NewAnalyticsHandlers(logger *observability.Logger, pe *analytics.PerformanceEngine) *AnalyticsHandlers {
	return &AnalyticsHandlers{
		logger:            logger,
		performanceEngine: pe,
	}
}

// RegisterRoutes registers analytics API routes
func (h *AnalyticsHandlers) RegisterRoutes(router *mux.Router) {
	// Performance metrics
	router.HandleFunc("/api/analytics/performance", h.GetPerformanceMetrics).Methods("GET")
	router.HandleFunc("/api/analytics/performance/trading", h.GetTradingMetrics).Methods("GET")
	router.HandleFunc("/api/analytics/performance/system", h.GetSystemMetrics).Methods("GET")
	router.HandleFunc("/api/analytics/performance/portfolio", h.GetPortfolioMetrics).Methods("GET")

	// Performance analysis
	router.HandleFunc("/api/analytics/analysis/overview", h.GetPerformanceOverview).Methods("GET")
	router.HandleFunc("/api/analytics/analysis/trends", h.GetPerformanceTrends).Methods("GET")
	router.HandleFunc("/api/analytics/analysis/comparison", h.GetPerformanceComparison).Methods("GET")

	// Optimization
	router.HandleFunc("/api/analytics/optimization/opportunities", h.GetOptimizationOpportunities).Methods("GET")
	router.HandleFunc("/api/analytics/optimization/recommendations", h.GetOptimizationRecommendations).Methods("GET")
	router.HandleFunc("/api/analytics/optimization/score", h.GetOptimizationScore).Methods("GET")

	// Benchmarking
	router.HandleFunc("/api/analytics/benchmarks", h.GetBenchmarks).Methods("GET")
	router.HandleFunc("/api/analytics/benchmarks/{id}/compare", h.CompareToBenchmark).Methods("GET")

	// Reports
	router.HandleFunc("/api/analytics/reports/performance", h.GeneratePerformanceReport).Methods("POST")
	router.HandleFunc("/api/analytics/reports/{id}", h.GetPerformanceReport).Methods("GET")
	router.HandleFunc("/api/analytics/reports/{id}/export", h.ExportPerformanceReport).Methods("GET")

	// Dashboard
	router.HandleFunc("/api/analytics/dashboard", h.GetAnalyticsDashboard).Methods("GET")
}

// GetPerformanceMetrics returns comprehensive performance metrics
func (h *AnalyticsHandlers) GetPerformanceMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info(ctx, "Getting performance metrics", nil)

	metrics := h.performanceEngine.GetMetrics()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// GetTradingMetrics returns trading performance metrics
func (h *AnalyticsHandlers) GetTradingMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info(ctx, "Getting trading metrics", nil)

	metrics := h.performanceEngine.GetTradingMetrics()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// GetSystemMetrics returns system performance metrics
func (h *AnalyticsHandlers) GetSystemMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info(ctx, "Getting system metrics", nil)

	metrics := h.performanceEngine.GetSystemMetrics()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// GetPortfolioMetrics returns portfolio performance metrics
func (h *AnalyticsHandlers) GetPortfolioMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info(ctx, "Getting portfolio metrics", nil)

	metrics := h.performanceEngine.GetPortfolioMetrics()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// GetPerformanceOverview returns performance overview
func (h *AnalyticsHandlers) GetPerformanceOverview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info(ctx, "Getting performance overview", nil)

	metrics := h.performanceEngine.GetMetrics()

	overview := map[string]interface{}{
		"summary": map[string]interface{}{
			"overall_score":      85.5,
			"trading_score":      metrics.Trading.SharpeRatio * 20,
			"system_score":       95.0 - metrics.System.CPUUsage*0.5,
			"portfolio_score":    metrics.Portfolio.SharpeRatio * 30,
			"optimization_score": metrics.Optimization.OptimizationScore,
		},
		"key_metrics": map[string]interface{}{
			"total_return":     metrics.Portfolio.TotalReturn,
			"sharpe_ratio":     metrics.Trading.SharpeRatio,
			"max_drawdown":     metrics.Portfolio.MaxDrawdown,
			"success_rate":     metrics.Trading.SuccessRate,
			"avg_latency_ms":   metrics.Execution.AverageLatency.Milliseconds(),
			"cpu_usage":        metrics.System.CPUUsage,
			"memory_usage":     metrics.System.MemoryUsage,
		},
		"status": map[string]interface{}{
			"trading_status":  "ACTIVE",
			"system_status":   "HEALTHY",
			"risk_status":     "NORMAL",
			"last_updated":    time.Now().Format(time.RFC3339),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(overview)
}

// GetPerformanceTrends returns performance trends analysis
func (h *AnalyticsHandlers) GetPerformanceTrends(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info(ctx, "Getting performance trends", nil)

	// Parse query parameters
	query := r.URL.Query()
	period := query.Get("period")
	if period == "" {
		period = "24h"
	}

	// Mock trends data - in production, calculate from historical data
	trends := map[string]interface{}{
		"period": period,
		"trends": []map[string]interface{}{
			{
				"metric":         "sharpe_ratio",
				"direction":      "up",
				"change_percent": 5.2,
				"confidence":     0.85,
			},
			{
				"metric":         "total_return",
				"direction":      "up",
				"change_percent": 2.1,
				"confidence":     0.92,
			},
			{
				"metric":         "max_drawdown",
				"direction":      "down",
				"change_percent": -1.5,
				"confidence":     0.78,
			},
			{
				"metric":         "cpu_usage",
				"direction":      "stable",
				"change_percent": 0.3,
				"confidence":     0.95,
			},
		},
		"last_updated": time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trends)
}

// GetPerformanceComparison returns performance comparison data
func (h *AnalyticsHandlers) GetPerformanceComparison(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info(ctx, "Getting performance comparison", nil)

	// Parse query parameters
	query := r.URL.Query()
	compareWith := query.Get("compare_with")
	period := query.Get("period")

	if compareWith == "" {
		compareWith = "benchmark"
	}
	if period == "" {
		period = "30d"
	}

	// Mock comparison data
	comparison := map[string]interface{}{
		"comparison_type": compareWith,
		"period":          period,
		"metrics": map[string]interface{}{
			"portfolio": map[string]interface{}{
				"total_return":   "12.5%",
				"sharpe_ratio":   1.85,
				"max_drawdown":   "8.2%",
				"volatility":     "15.3%",
			},
			"benchmark": map[string]interface{}{
				"total_return":   "8.7%",
				"sharpe_ratio":   1.42,
				"max_drawdown":   "12.1%",
				"volatility":     "18.9%",
			},
		},
		"relative_performance": map[string]interface{}{
			"excess_return":      "3.8%",
			"tracking_error":     "4.2%",
			"information_ratio":  0.90,
			"beta":               0.85,
			"alpha":              "2.1%",
		},
		"last_updated": time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comparison)
}

// GetOptimizationOpportunities returns optimization opportunities
func (h *AnalyticsHandlers) GetOptimizationOpportunities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info(ctx, "Getting optimization opportunities", nil)

	// Mock optimization opportunities
	opportunities := []map[string]interface{}{
		{
			"id":          "opt_001",
			"type":        "latency",
			"priority":    "HIGH",
			"title":       "Reduce API Latency",
			"description": "Optimize database queries and implement caching",
			"impact":      "High",
			"effort":      "Medium",
			"roi":         85.5,
			"category":    "performance",
		},
		{
			"id":          "opt_002",
			"type":        "strategy",
			"priority":    "MEDIUM",
			"title":       "Improve Sharpe Ratio",
			"description": "Optimize strategy parameters and risk management",
			"impact":      "High",
			"effort":      "High",
			"roi":         72.3,
			"category":    "trading",
		},
		{
			"id":          "opt_003",
			"type":        "system",
			"priority":    "LOW",
			"title":       "Memory Optimization",
			"description": "Optimize memory usage and garbage collection",
			"impact":      "Medium",
			"effort":      "Low",
			"roi":         45.8,
			"category":    "system",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"opportunities": opportunities,
		"total":         len(opportunities),
		"last_updated":  time.Now().Format(time.RFC3339),
	})
}

// GetOptimizationRecommendations returns optimization recommendations
func (h *AnalyticsHandlers) GetOptimizationRecommendations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info(ctx, "Getting optimization recommendations", nil)

	// Mock recommendations
	recommendations := []map[string]interface{}{
		{
			"id":          "rec_001",
			"priority":    "HIGH",
			"title":       "Implement Redis Caching",
			"description": "Add Redis caching layer to reduce database load",
			"impact":      "High",
			"effort":      "Medium",
			"timeline":    "2-3 weeks",
			"resources":   []string{"Backend Developer", "DevOps Engineer"},
			"status":      "PENDING",
		},
		{
			"id":          "rec_002",
			"priority":    "MEDIUM",
			"title":       "Optimize Trading Algorithms",
			"description": "Review and optimize core trading algorithms",
			"impact":      "High",
			"effort":      "High",
			"timeline":    "6-8 weeks",
			"resources":   []string{"Quantitative Developer", "Trading Engineer"},
			"status":      "IN_PROGRESS",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"recommendations": recommendations,
		"total":           len(recommendations),
		"last_updated":    time.Now().Format(time.RFC3339),
	})
}

// GetOptimizationScore returns optimization score
func (h *AnalyticsHandlers) GetOptimizationScore(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info(ctx, "Getting optimization score", nil)

	metrics := h.performanceEngine.GetMetrics()

	score := map[string]interface{}{
		"overall_score":         metrics.Optimization.OptimizationScore,
		"efficiency_rating":     metrics.Optimization.EfficiencyRating,
		"resource_utilization":  metrics.Optimization.ResourceUtilization,
		"performance_gap":       metrics.Optimization.PerformanceGap,
		"potential_improvement": metrics.Optimization.PotentialImprovement,
		"implementation_cost":   metrics.Optimization.ImplementationCost,
		"roi":                   metrics.Optimization.ROI,
		"breakdown": map[string]interface{}{
			"system_score":    85.0,
			"trading_score":   78.5,
			"portfolio_score": 82.3,
			"risk_score":      88.7,
		},
		"last_updated": time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(score)
}

// GetBenchmarks returns available benchmarks
func (h *AnalyticsHandlers) GetBenchmarks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info(ctx, "Getting benchmarks", nil)

	// Mock benchmarks
	benchmarks := []map[string]interface{}{
		{
			"id":          "btc_index",
			"name":        "Bitcoin Index",
			"description": "Bitcoin price performance benchmark",
			"type":        "MARKET",
			"active":      true,
		},
		{
			"id":          "crypto_index",
			"name":        "Crypto Market Index",
			"description": "Weighted cryptocurrency market index",
			"type":        "MARKET",
			"active":      true,
		},
		{
			"id":          "risk_free",
			"name":        "Risk-Free Rate",
			"description": "US Treasury 3-month rate",
			"type":        "RISK_FREE",
			"active":      true,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"benchmarks": benchmarks,
		"total":      len(benchmarks),
	})
}

// CompareToBenchmark compares performance to a benchmark
func (h *AnalyticsHandlers) CompareToBenchmark(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	benchmarkID := vars["id"]

	h.logger.Info(ctx, "Comparing to benchmark", map[string]interface{}{
		"benchmark_id": benchmarkID,
	})

	// Mock comparison data
	comparison := map[string]interface{}{
		"benchmark_id":     benchmarkID,
		"benchmark_name":   "Bitcoin Index",
		"period":           "30d",
		"portfolio_return": "12.5%",
		"benchmark_return": "8.7%",
		"excess_return":    "3.8%",
		"tracking_error":   "4.2%",
		"information_ratio": 0.90,
		"beta":             0.85,
		"alpha":            "2.1%",
		"correlation":      0.78,
		"up_capture":       "105.2%",
		"down_capture":     "82.3%",
		"batting_average":  "65.5%",
		"last_updated":     time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comparison)
}

// GeneratePerformanceReport generates a performance report
func (h *AnalyticsHandlers) GeneratePerformanceReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info(ctx, "Generating performance report", nil)

	var request struct {
		Type       string `json:"type"`
		Period     string `json:"period"`
		Sections   []string `json:"sections"`
		Format     string `json:"format"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Mock report generation
	reportID := "report_" + strconv.FormatInt(time.Now().Unix(), 10)
	
	report := map[string]interface{}{
		"id":           reportID,
		"type":         request.Type,
		"period":       request.Period,
		"format":       request.Format,
		"status":       "GENERATED",
		"generated_at": time.Now().Format(time.RFC3339),
		"sections":     request.Sections,
		"download_url": "/api/analytics/reports/" + reportID + "/export",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(report)
}

// GetPerformanceReport returns a specific performance report
func (h *AnalyticsHandlers) GetPerformanceReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	reportID := vars["id"]

	h.logger.Info(ctx, "Getting performance report", map[string]interface{}{
		"report_id": reportID,
	})

	// Mock report data
	report := map[string]interface{}{
		"id":           reportID,
		"type":         "COMPREHENSIVE",
		"period":       "30d",
		"format":       "PDF",
		"status":       "GENERATED",
		"generated_at": time.Now().Format(time.RFC3339),
		"data": map[string]interface{}{
			"summary": map[string]interface{}{
				"total_return":   "12.5%",
				"sharpe_ratio":   1.85,
				"max_drawdown":   "8.2%",
				"success_rate":   "96.8%",
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// ExportPerformanceReport exports a performance report
func (h *AnalyticsHandlers) ExportPerformanceReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	reportID := vars["id"]
	format := r.URL.Query().Get("format")

	if format == "" {
		format = "pdf"
	}

	h.logger.Info(ctx, "Exporting performance report", map[string]interface{}{
		"report_id": reportID,
		"format":    format,
	})

	switch format {
	case "pdf":
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "attachment; filename=performance_report_"+reportID+".pdf")
	case "csv":
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=performance_report_"+reportID+".csv")
	case "json":
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", "attachment; filename=performance_report_"+reportID+".json")
	default:
		http.Error(w, "Unsupported format", http.StatusBadRequest)
		return
	}

	// Mock export content
	w.Write([]byte("Mock performance report content"))
}

// GetAnalyticsDashboard returns analytics dashboard data
func (h *AnalyticsHandlers) GetAnalyticsDashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info(ctx, "Getting analytics dashboard", nil)

	metrics := h.performanceEngine.GetMetrics()

	dashboard := map[string]interface{}{
		"overview": map[string]interface{}{
			"overall_score":      85.5,
			"trading_score":      metrics.Trading.SharpeRatio * 20,
			"system_score":       95.0 - metrics.System.CPUUsage*0.5,
			"portfolio_score":    metrics.Portfolio.SharpeRatio * 30,
			"optimization_score": metrics.Optimization.OptimizationScore,
		},
		"key_metrics": map[string]interface{}{
			"total_return":     metrics.Portfolio.TotalReturn,
			"sharpe_ratio":     metrics.Trading.SharpeRatio,
			"max_drawdown":     metrics.Portfolio.MaxDrawdown,
			"success_rate":     metrics.Trading.SuccessRate,
			"avg_latency_ms":   metrics.Execution.AverageLatency.Milliseconds(),
		},
		"system_health": map[string]interface{}{
			"cpu_usage":      metrics.System.CPUUsage,
			"memory_usage":   metrics.System.MemoryUsage,
			"disk_usage":     metrics.System.DiskUsage,
			"error_rate":     metrics.System.ErrorRate,
			"uptime":         metrics.System.Uptime.String(),
		},
		"optimization": map[string]interface{}{
			"score":                  metrics.Optimization.OptimizationScore,
			"opportunities_count":    len(metrics.Optimization.OptimizationOpportunities),
			"potential_improvement":  metrics.Optimization.PotentialImprovement,
			"recommended_actions":    metrics.Optimization.RecommendedActions,
		},
		"alerts": []map[string]interface{}{
			{
				"type":     "performance",
				"severity": "medium",
				"message":  "CPU usage above 80%",
			},
		},
		"last_updated": time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dashboard)
}
