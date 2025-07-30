package analytics

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/ai-agentic-browser/pkg/database"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

// Service provides analytics and business intelligence capabilities
type Service struct {
	db     *database.DB
	redis  *database.RedisClient
	logger *observability.Logger
}

// NewService creates a new analytics service
func NewService(db *database.DB, redis *database.RedisClient, logger *observability.Logger) *Service {
	return &Service{
		db:     db,
		redis:  redis,
		logger: logger,
	}
}

// AnalyticsRequest represents a request for analytics data
type AnalyticsRequest struct {
	UserID      *uuid.UUID             `json:"user_id,omitempty"`
	TeamID      *uuid.UUID             `json:"team_id,omitempty"`
	StartDate   time.Time              `json:"start_date"`
	EndDate     time.Time              `json:"end_date"`
	Granularity string                 `json:"granularity"` // hour, day, week, month
	Metrics     []string               `json:"metrics"`
	Filters     map[string]interface{} `json:"filters"`
}

// AnalyticsResponse represents analytics data response
type AnalyticsResponse struct {
	Period      string                 `json:"period"`
	Granularity string                 `json:"granularity"`
	Metrics     map[string]interface{} `json:"metrics"`
	TimeSeries  []TimeSeriesPoint      `json:"time_series"`
	Insights    []Insight              `json:"insights"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// TimeSeriesPoint represents a point in time series data
type TimeSeriesPoint struct {
	Timestamp time.Time              `json:"timestamp"`
	Values    map[string]interface{} `json:"values"`
}

// Insight represents an analytical insight
type Insight struct {
	Type        string                 `json:"type"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Severity    string                 `json:"severity"`
	Value       interface{}            `json:"value"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// GetAnalytics retrieves analytics data based on the request
func (s *Service) GetAnalytics(ctx context.Context, req AnalyticsRequest) (*AnalyticsResponse, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("analytics-service").Start(ctx, "analytics.GetAnalytics")
	defer span.End()

	// Validate request
	if err := s.validateRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Generate cache key
	cacheKey := s.generateCacheKey(req)

	// Try to get from cache first
	if cached, err := s.getFromCache(ctx, cacheKey); err == nil && cached != nil {
		return cached, nil
	}

	// Calculate metrics
	metrics, err := s.calculateMetrics(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate metrics: %w", err)
	}

	// Generate time series data
	timeSeries, err := s.generateTimeSeries(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate time series: %w", err)
	}

	// Generate insights
	insights, err := s.generateInsights(ctx, req, metrics, timeSeries)
	if err != nil {
		s.logger.Error(ctx, "Failed to generate insights", err)
		insights = []Insight{} // Continue without insights
	}

	response := &AnalyticsResponse{
		Period:      fmt.Sprintf("%s to %s", req.StartDate.Format("2006-01-02"), req.EndDate.Format("2006-01-02")),
		Granularity: req.Granularity,
		Metrics:     metrics,
		TimeSeries:  timeSeries,
		Insights:    insights,
		Metadata: map[string]interface{}{
			"generated_at": time.Now(),
			"cache_key":    cacheKey,
		},
	}

	// Cache the response
	s.cacheResponse(ctx, cacheKey, response)

	return response, nil
}

// calculateMetrics calculates various metrics based on the request
func (s *Service) calculateMetrics(ctx context.Context, req AnalyticsRequest) (map[string]interface{}, error) {
	metrics := make(map[string]interface{})

	for _, metric := range req.Metrics {
		switch metric {
		case "total_executions":
			value, err := s.getTotalExecutions(ctx, req)
			if err != nil {
				return nil, err
			}
			metrics[metric] = value

		case "success_rate":
			value, err := s.getSuccessRate(ctx, req)
			if err != nil {
				return nil, err
			}
			metrics[metric] = value

		case "average_execution_time":
			value, err := s.getAverageExecutionTime(ctx, req)
			if err != nil {
				return nil, err
			}
			metrics[metric] = value

		case "cost_analysis":
			value, err := s.getCostAnalysis(ctx, req)
			if err != nil {
				return nil, err
			}
			metrics[metric] = value

		case "user_engagement":
			value, err := s.getUserEngagement(ctx, req)
			if err != nil {
				return nil, err
			}
			metrics[metric] = value

		case "workflow_performance":
			value, err := s.getWorkflowPerformance(ctx, req)
			if err != nil {
				return nil, err
			}
			metrics[metric] = value

		case "resource_utilization":
			value, err := s.getResourceUtilization(ctx, req)
			if err != nil {
				return nil, err
			}
			metrics[metric] = value

		case "error_analysis":
			value, err := s.getErrorAnalysis(ctx, req)
			if err != nil {
				return nil, err
			}
			metrics[metric] = value
		}
	}

	return metrics, nil
}

// getTotalExecutions calculates total workflow executions
func (s *Service) getTotalExecutions(ctx context.Context, req AnalyticsRequest) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM workflow_executions 
		WHERE started_at >= $1 AND started_at <= $2
	`
	args := []interface{}{req.StartDate, req.EndDate}

	if req.UserID != nil {
		query += " AND user_id = $3"
		args = append(args, *req.UserID)
	} else if req.TeamID != nil {
		query += " AND user_id IN (SELECT user_id FROM team_members WHERE team_id = $3 AND status = 'active')"
		args = append(args, *req.TeamID)
	}

	var count int
	err := s.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

// getSuccessRate calculates workflow success rate
func (s *Service) getSuccessRate(ctx context.Context, req AnalyticsRequest) (float64, error) {
	query := `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as successful
		FROM workflow_executions 
		WHERE started_at >= $1 AND started_at <= $2
	`
	args := []interface{}{req.StartDate, req.EndDate}

	if req.UserID != nil {
		query += " AND user_id = $3"
		args = append(args, *req.UserID)
	} else if req.TeamID != nil {
		query += " AND user_id IN (SELECT user_id FROM team_members WHERE team_id = $3 AND status = 'active')"
		args = append(args, *req.TeamID)
	}

	var total, successful int
	err := s.db.QueryRowContext(ctx, query, args...).Scan(&total, &successful)
	if err != nil {
		return 0, err
	}

	if total == 0 {
		return 0, nil
	}

	return float64(successful) / float64(total) * 100, nil
}

// getAverageExecutionTime calculates average execution time
func (s *Service) getAverageExecutionTime(ctx context.Context, req AnalyticsRequest) (float64, error) {
	query := `
		SELECT AVG(duration) 
		FROM workflow_executions 
		WHERE started_at >= $1 AND started_at <= $2 AND duration > 0
	`
	args := []interface{}{req.StartDate, req.EndDate}

	if req.UserID != nil {
		query += " AND user_id = $3"
		args = append(args, *req.UserID)
	} else if req.TeamID != nil {
		query += " AND user_id IN (SELECT user_id FROM team_members WHERE team_id = $3 AND status = 'active')"
		args = append(args, *req.TeamID)
	}

	var avgDuration float64
	err := s.db.QueryRowContext(ctx, query, args...).Scan(&avgDuration)
	return avgDuration / 1000, err // Convert to seconds
}

// getCostAnalysis calculates cost analysis
func (s *Service) getCostAnalysis(ctx context.Context, req AnalyticsRequest) (map[string]interface{}, error) {
	// This would integrate with billing/cost tracking systems
	// For now, return mock data
	return map[string]interface{}{
		"total_cost":         125.50,
		"ai_cost":            75.30,
		"browser_cost":       25.20,
		"web3_cost":          15.00,
		"storage_cost":       10.00,
		"cost_per_execution": 1.25,
		"cost_trend":         "increasing",
	}, nil
}

// getUserEngagement calculates user engagement metrics
func (s *Service) getUserEngagement(ctx context.Context, req AnalyticsRequest) (map[string]interface{}, error) {
	// Calculate various engagement metrics
	query := `
		SELECT 
			COUNT(DISTINCT user_id) as active_users,
			COUNT(*) as total_sessions,
			AVG(duration) as avg_session_duration
		FROM workflow_executions 
		WHERE started_at >= $1 AND started_at <= $2
	`
	args := []interface{}{req.StartDate, req.EndDate}

	if req.TeamID != nil {
		query += " AND user_id IN (SELECT user_id FROM team_members WHERE team_id = $3 AND status = 'active')"
		args = append(args, *req.TeamID)
	}

	var activeUsers, totalSessions int
	var avgSessionDuration float64
	err := s.db.QueryRowContext(ctx, query, args...).Scan(&activeUsers, &totalSessions, &avgSessionDuration)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"active_users":         activeUsers,
		"total_sessions":       totalSessions,
		"avg_session_duration": avgSessionDuration / 1000, // Convert to seconds
		"sessions_per_user":    float64(totalSessions) / float64(activeUsers),
	}, nil
}

// getWorkflowPerformance analyzes workflow performance
func (s *Service) getWorkflowPerformance(ctx context.Context, req AnalyticsRequest) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			w.id,
			w.name,
			COUNT(we.id) as executions,
			AVG(we.duration) as avg_duration,
			COUNT(CASE WHEN we.status = 'completed' THEN 1 END) * 100.0 / COUNT(*) as success_rate
		FROM workflows w
		LEFT JOIN workflow_executions we ON w.id = we.workflow_id 
			AND we.started_at >= $1 AND we.started_at <= $2
		GROUP BY w.id, w.name
		ORDER BY executions DESC
		LIMIT 10
	`
	args := []interface{}{req.StartDate, req.EndDate}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workflows []map[string]interface{}
	for rows.Next() {
		var id uuid.UUID
		var name string
		var executions int
		var avgDuration, successRate float64

		err := rows.Scan(&id, &name, &executions, &avgDuration, &successRate)
		if err != nil {
			continue
		}

		workflows = append(workflows, map[string]interface{}{
			"id":           id,
			"name":         name,
			"executions":   executions,
			"avg_duration": avgDuration / 1000, // Convert to seconds
			"success_rate": successRate,
		})
	}

	return workflows, nil
}

// getResourceUtilization calculates resource utilization
func (s *Service) getResourceUtilization(ctx context.Context, req AnalyticsRequest) (map[string]interface{}, error) {
	// This would integrate with infrastructure monitoring
	// For now, return mock data
	return map[string]interface{}{
		"cpu_utilization":     65.5,
		"memory_utilization":  72.3,
		"storage_utilization": 45.8,
		"network_utilization": 23.1,
		"concurrent_sessions": 15,
		"peak_usage_time":     "14:30",
	}, nil
}

// getErrorAnalysis analyzes errors and failures
func (s *Service) getErrorAnalysis(ctx context.Context, req AnalyticsRequest) (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as total_errors,
			COUNT(*) as total_executions
		FROM workflow_executions 
		WHERE started_at >= $1 AND started_at <= $2
	`
	args := []interface{}{req.StartDate, req.EndDate}

	if req.UserID != nil {
		query += " AND user_id = $3"
		args = append(args, *req.UserID)
	} else if req.TeamID != nil {
		query += " AND user_id IN (SELECT user_id FROM team_members WHERE team_id = $3 AND status = 'active')"
		args = append(args, *req.TeamID)
	}

	var totalErrors, totalExecutions int
	err := s.db.QueryRowContext(ctx, query, args...).Scan(&totalErrors, &totalExecutions)
	if err != nil {
		return nil, err
	}

	errorRate := float64(0)
	if totalExecutions > 0 {
		errorRate = float64(totalErrors) / float64(totalExecutions) * 100
	}

	// Get top error types
	errorTypesQuery := `
		SELECT 
			COALESCE(error_message, 'Unknown') as error_type,
			COUNT(*) as count
		FROM workflow_executions 
		WHERE started_at >= $1 AND started_at <= $2 AND status = 'failed'
		GROUP BY error_message
		ORDER BY count DESC
		LIMIT 5
	`

	rows, err := s.db.QueryContext(ctx, errorTypesQuery, req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var errorTypes []map[string]interface{}
	for rows.Next() {
		var errorType string
		var count int
		if err := rows.Scan(&errorType, &count); err == nil {
			errorTypes = append(errorTypes, map[string]interface{}{
				"type":  errorType,
				"count": count,
			})
		}
	}

	return map[string]interface{}{
		"total_errors": totalErrors,
		"error_rate":   errorRate,
		"error_types":  errorTypes,
		"mttr":         45.5, // Mean Time To Recovery (mock)
		"error_trend":  "decreasing",
	}, nil
}

// generateTimeSeries generates time series data
func (s *Service) generateTimeSeries(ctx context.Context, req AnalyticsRequest) ([]TimeSeriesPoint, error) {
	var points []TimeSeriesPoint

	// Determine time intervals based on granularity
	interval := s.getTimeInterval(req.Granularity)
	current := req.StartDate

	for current.Before(req.EndDate) {
		next := current.Add(interval)
		if next.After(req.EndDate) {
			next = req.EndDate
		}

		// Get metrics for this time period
		periodReq := req
		periodReq.StartDate = current
		periodReq.EndDate = next

		values := make(map[string]interface{})

		// Calculate metrics for this period
		if contains(req.Metrics, "total_executions") {
			if count, err := s.getTotalExecutions(ctx, periodReq); err == nil {
				values["executions"] = count
			}
		}

		if contains(req.Metrics, "success_rate") {
			if rate, err := s.getSuccessRate(ctx, periodReq); err == nil {
				values["success_rate"] = rate
			}
		}

		points = append(points, TimeSeriesPoint{
			Timestamp: current,
			Values:    values,
		})

		current = next
	}

	return points, nil
}

// generateInsights generates analytical insights
func (s *Service) generateInsights(ctx context.Context, req AnalyticsRequest, metrics map[string]interface{}, timeSeries []TimeSeriesPoint) ([]Insight, error) {
	var insights []Insight

	// Success rate insight
	if successRate, ok := metrics["success_rate"].(float64); ok {
		severity := "info"
		if successRate < 80 {
			severity = "warning"
		}
		if successRate < 60 {
			severity = "critical"
		}

		insights = append(insights, Insight{
			Type:        "success_rate",
			Title:       "Workflow Success Rate",
			Description: fmt.Sprintf("Current success rate is %.1f%%", successRate),
			Severity:    severity,
			Value:       successRate,
		})
	}

	// Performance trend insight
	if len(timeSeries) > 1 {
		trend := s.calculateTrend(timeSeries, "executions")
		insights = append(insights, Insight{
			Type:        "performance_trend",
			Title:       "Execution Trend",
			Description: fmt.Sprintf("Execution volume is %s", trend),
			Severity:    "info",
			Value:       trend,
		})
	}

	// Cost optimization insight
	if costData, ok := metrics["cost_analysis"].(map[string]interface{}); ok {
		if costPerExecution, ok := costData["cost_per_execution"].(float64); ok && costPerExecution > 2.0 {
			insights = append(insights, Insight{
				Type:        "cost_optimization",
				Title:       "Cost Optimization Opportunity",
				Description: fmt.Sprintf("Cost per execution (%.2f) is above recommended threshold", costPerExecution),
				Severity:    "warning",
				Value:       costPerExecution,
			})
		}
	}

	return insights, nil
}

// Helper methods

func (s *Service) validateRequest(req AnalyticsRequest) error {
	if req.EndDate.Before(req.StartDate) {
		return fmt.Errorf("end date must be after start date")
	}

	if len(req.Metrics) == 0 {
		return fmt.Errorf("at least one metric must be specified")
	}

	validGranularities := []string{"hour", "day", "week", "month"}
	if !contains(validGranularities, req.Granularity) {
		return fmt.Errorf("invalid granularity: %s", req.Granularity)
	}

	return nil
}

func (s *Service) generateCacheKey(req AnalyticsRequest) string {
	data, _ := json.Marshal(req)
	return fmt.Sprintf("analytics:%x", data)
}

func (s *Service) getFromCache(ctx context.Context, key string) (*AnalyticsResponse, error) {
	data, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var response AnalyticsResponse
	if err := json.Unmarshal([]byte(data), &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *Service) cacheResponse(ctx context.Context, key string, response *AnalyticsResponse) {
	data, _ := json.Marshal(response)
	s.redis.Set(ctx, key, data, 15*time.Minute) // Cache for 15 minutes
}

func (s *Service) getTimeInterval(granularity string) time.Duration {
	switch granularity {
	case "hour":
		return time.Hour
	case "day":
		return 24 * time.Hour
	case "week":
		return 7 * 24 * time.Hour
	case "month":
		return 30 * 24 * time.Hour
	default:
		return 24 * time.Hour
	}
}

func (s *Service) calculateTrend(timeSeries []TimeSeriesPoint, metric string) string {
	if len(timeSeries) < 2 {
		return "stable"
	}

	first := timeSeries[0].Values[metric]
	last := timeSeries[len(timeSeries)-1].Values[metric]

	firstVal, ok1 := first.(float64)
	lastVal, ok2 := last.(float64)

	if !ok1 || !ok2 {
		return "stable"
	}

	change := (lastVal - firstVal) / firstVal * 100

	if math.Abs(change) < 5 {
		return "stable"
	} else if change > 0 {
		return "increasing"
	} else {
		return "decreasing"
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
