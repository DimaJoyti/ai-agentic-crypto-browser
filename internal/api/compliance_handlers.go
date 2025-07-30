package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ai-agentic-browser/internal/compliance"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/gorilla/mux"
)

// ComplianceHandlers handles compliance-related API endpoints
type ComplianceHandlers struct {
	logger            *observability.Logger
	complianceManager *compliance.ComplianceManager
}

// NewComplianceHandlers creates new compliance handlers
func NewComplianceHandlers(logger *observability.Logger, cm *compliance.ComplianceManager) *ComplianceHandlers {
	return &ComplianceHandlers{
		logger:            logger,
		complianceManager: cm,
	}
}

// RegisterRoutes registers compliance API routes
func (h *ComplianceHandlers) RegisterRoutes(router *mux.Router) {
	// Compliance frameworks
	router.HandleFunc("/api/compliance/frameworks", h.GetFrameworks).Methods("GET")
	router.HandleFunc("/api/compliance/frameworks/{id}", h.GetFramework).Methods("GET")
	router.HandleFunc("/api/compliance/frameworks/{id}/status", h.GetFrameworkStatus).Methods("GET")

	// Compliance reports
	router.HandleFunc("/api/compliance/reports", h.GetReports).Methods("GET")
	router.HandleFunc("/api/compliance/reports", h.GenerateReport).Methods("POST")
	router.HandleFunc("/api/compliance/reports/{id}", h.GetReport).Methods("GET")
	router.HandleFunc("/api/compliance/reports/{id}/export", h.ExportReport).Methods("GET")

	// Compliance violations
	router.HandleFunc("/api/compliance/violations", h.GetViolations).Methods("GET")
	router.HandleFunc("/api/compliance/violations/{id}", h.GetViolation).Methods("GET")
	router.HandleFunc("/api/compliance/violations/{id}/resolve", h.ResolveViolation).Methods("POST")

	// Risk management
	router.HandleFunc("/api/risk/metrics", h.GetRiskMetrics).Methods("GET")
	router.HandleFunc("/api/risk/alerts", h.GetRiskAlerts).Methods("GET")
	router.HandleFunc("/api/risk/alerts/{id}/acknowledge", h.AcknowledgeAlert).Methods("POST")
	router.HandleFunc("/api/risk/alerts/{id}/resolve", h.ResolveAlert).Methods("POST")
	router.HandleFunc("/api/risk/thresholds", h.GetRiskThresholds).Methods("GET")
	router.HandleFunc("/api/risk/thresholds", h.CreateRiskThreshold).Methods("POST")

	// Audit trail
	router.HandleFunc("/api/audit/events", h.GetAuditEvents).Methods("GET")
	router.HandleFunc("/api/audit/events/export", h.ExportAuditEvents).Methods("GET")
	router.HandleFunc("/api/audit/summary", h.GetAuditSummary).Methods("GET")

	// Alert management
	router.HandleFunc("/api/alerts", h.GetAlerts).Methods("GET")
	router.HandleFunc("/api/alerts/{id}/acknowledge", h.AcknowledgeAlert).Methods("POST")
	router.HandleFunc("/api/alerts/{id}/resolve", h.ResolveAlert).Methods("POST")
	router.HandleFunc("/api/alerts/rules", h.GetAlertRules).Methods("GET")
	router.HandleFunc("/api/alerts/channels", h.GetAlertChannels).Methods("GET")

	// Compliance dashboard
	router.HandleFunc("/api/compliance/dashboard", h.GetComplianceDashboard).Methods("GET")
	router.HandleFunc("/api/compliance/overview", h.GetComplianceOverview).Methods("GET")
}

// GetFrameworks returns all compliance frameworks
func (h *ComplianceHandlers) GetFrameworks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info(ctx, "Getting compliance frameworks", nil)

	// Mock response - in production, get from compliance manager
	frameworks := []map[string]interface{}{
		{
			"id":           "us_bsa",
			"name":         "US Bank Secrecy Act",
			"description":  "US Anti-Money Laundering regulations",
			"jurisdiction": "United States",
			"status":       "COMPLIANT",
			"last_update":  time.Now().Format(time.RFC3339),
		},
		{
			"id":           "eu_amld5",
			"name":         "EU Anti-Money Laundering Directive 5",
			"description":  "European Union AML regulations",
			"jurisdiction": "European Union",
			"status":       "PARTIAL",
			"last_update":  time.Now().Format(time.RFC3339),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"frameworks": frameworks,
		"total":      len(frameworks),
	})
}

// GetFramework returns a specific compliance framework
func (h *ComplianceHandlers) GetFramework(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	frameworkID := vars["id"]

	h.logger.Info(ctx, "Getting compliance framework", map[string]interface{}{
		"framework_id": frameworkID,
	})

	// Mock response - in production, get from compliance manager
	framework := map[string]interface{}{
		"id":           frameworkID,
		"name":         "US Bank Secrecy Act",
		"description":  "US Anti-Money Laundering regulations",
		"jurisdiction": "United States",
		"version":      "2023.1",
		"status":       "COMPLIANT",
		"requirements": []map[string]interface{}{
			{
				"id":          "customer_identification",
				"name":        "Customer Identification Program",
				"description": "Verify customer identity before account opening",
				"category":    "KYC",
				"mandatory":   true,
				"status":      "COMPLIANT",
				"risk_level":  "HIGH",
			},
		},
		"last_update": time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(framework)
}

// GetFrameworkStatus returns compliance status for a framework
func (h *ComplianceHandlers) GetFrameworkStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	frameworkID := vars["id"]

	h.logger.Info(ctx, "Getting framework status", map[string]interface{}{
		"framework_id": frameworkID,
	})

	status := map[string]interface{}{
		"framework_id":               frameworkID,
		"overall_status":             "COMPLIANT",
		"compliance_rate":            85.0,
		"total_requirements":         10,
		"compliant_requirements":     8,
		"non_compliant_requirements": 1,
		"pending_requirements":       1,
		"last_assessment":            time.Now().Format(time.RFC3339),
		"next_review":                time.Now().AddDate(0, 3, 0).Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// GetReports returns compliance reports
func (h *ComplianceHandlers) GetReports(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info(ctx, "Getting compliance reports", nil)

	// Parse query parameters
	query := r.URL.Query()
	reportType := query.Get("type")
	framework := query.Get("framework")
	status := query.Get("status")

	// Mock response - in production, get from report generator
	reports := []map[string]interface{}{
		{
			"id":           "report_001",
			"type":         "COMPLIANCE",
			"framework":    "us_bsa",
			"status":       "GENERATED",
			"generated_at": time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
			"period": map[string]interface{}{
				"start_date": time.Now().AddDate(0, -1, 0).Format(time.RFC3339),
				"end_date":   time.Now().Format(time.RFC3339),
				"label":      "Monthly Report",
			},
		},
		{
			"id":           "report_002",
			"type":         "RISK",
			"framework":    "internal",
			"status":       "GENERATED",
			"generated_at": time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
			"period": map[string]interface{}{
				"start_date": time.Now().AddDate(0, 0, -1).Format(time.RFC3339),
				"end_date":   time.Now().Format(time.RFC3339),
				"label":      "Daily Report",
			},
		},
	}

	// Apply filters (simplified)
	filteredReports := []map[string]interface{}{}
	for _, report := range reports {
		if reportType != "" && report["type"] != reportType {
			continue
		}
		if framework != "" && report["framework"] != framework {
			continue
		}
		if status != "" && report["status"] != status {
			continue
		}
		filteredReports = append(filteredReports, report)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"reports": filteredReports,
		"total":   len(filteredReports),
	})
}

// GenerateReport generates a new compliance report
func (h *ComplianceHandlers) GenerateReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info(ctx, "Generating compliance report", nil)

	var request struct {
		TemplateID string            `json:"template_id"`
		Framework  string            `json:"framework"`
		Type       string            `json:"type"`
		Parameters map[string]string `json:"parameters"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Mock report generation - in production, use report generator
	reportID := "report_" + strconv.FormatInt(time.Now().Unix(), 10)

	report := map[string]interface{}{
		"id":           reportID,
		"type":         request.Type,
		"framework":    request.Framework,
		"status":       "GENERATED",
		"generated_at": time.Now().Format(time.RFC3339),
		"data": map[string]interface{}{
			"summary": map[string]interface{}{
				"total_requirements":         10,
				"compliant_requirements":     8,
				"non_compliant_requirements": 1,
				"pending_requirements":       1,
				"overall_compliance_rate":    80.0,
				"risk_level":                 "MEDIUM",
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(report)
}

// GetReport returns a specific compliance report
func (h *ComplianceHandlers) GetReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	reportID := vars["id"]

	h.logger.Info(ctx, "Getting compliance report", map[string]interface{}{
		"report_id": reportID,
	})

	// Mock response - in production, get from report generator
	report := map[string]interface{}{
		"id":           reportID,
		"type":         "COMPLIANCE",
		"framework":    "us_bsa",
		"status":       "GENERATED",
		"generated_at": time.Now().Format(time.RFC3339),
		"data": map[string]interface{}{
			"summary": map[string]interface{}{
				"total_requirements":         10,
				"compliant_requirements":     8,
				"non_compliant_requirements": 1,
				"pending_requirements":       1,
				"overall_compliance_rate":    80.0,
				"risk_level":                 "MEDIUM",
			},
			"metrics": map[string]interface{}{
				"total_trades":        1500,
				"total_volume":        "50000000",
				"violations_count":    5,
				"resolved_violations": 4,
				"compliance_score":    85.0,
			},
			"violations": []map[string]interface{}{
				{
					"id":          "violation_001",
					"type":        "RISK",
					"severity":    "MEDIUM",
					"description": "Position limit exceeded",
					"timestamp":   time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
					"resolved":    true,
				},
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// ExportReport exports a compliance report
func (h *ComplianceHandlers) ExportReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	reportID := vars["id"]
	format := r.URL.Query().Get("format")

	if format == "" {
		format = "json"
	}

	h.logger.Info(ctx, "Exporting compliance report", map[string]interface{}{
		"report_id": reportID,
		"format":    format,
	})

	switch format {
	case "json":
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", "attachment; filename=report_"+reportID+".json")
	case "csv":
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=report_"+reportID+".csv")
	case "pdf":
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "attachment; filename=report_"+reportID+".pdf")
	default:
		http.Error(w, "Unsupported format", http.StatusBadRequest)
		return
	}

	// Mock export - in production, use report generator
	w.Write([]byte("Mock exported report content"))
}

// GetRiskMetrics returns current risk metrics
func (h *ComplianceHandlers) GetRiskMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info(ctx, "Getting risk metrics", nil)

	// Mock response - in production, get from risk monitor
	metrics := map[string]interface{}{
		"total_exposure":     "5000000",
		"max_drawdown":       "250000",
		"current_drawdown":   "50000",
		"var_95":             "150000",
		"var_99":             "200000",
		"portfolio_value":    "10000000",
		"leverage_ratio":     2.5,
		"concentration_risk": 0.3,
		"correlation_risk":   0.25,
		"liquidity_risk":     0.15,
		"volatility_risk":    0.18,
		"last_updated":       time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// GetRiskAlerts returns current risk alerts
func (h *ComplianceHandlers) GetRiskAlerts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info(ctx, "Getting risk alerts", nil)

	// Parse query parameters
	query := r.URL.Query()
	alertType := query.Get("type")
	severity := query.Get("severity")
	resolved := query.Get("resolved")

	// Mock response - in production, get from risk monitor
	alerts := []map[string]interface{}{
		{
			"id":           "alert_001",
			"type":         "POSITION_LIMIT",
			"severity":     "WARNING",
			"title":        "Position limit approaching",
			"description":  "BTCUSD position approaching limit",
			"value":        "950000",
			"threshold":    "1000000",
			"timestamp":    time.Now().Add(-30 * time.Minute).Format(time.RFC3339),
			"acknowledged": false,
			"resolved":     false,
		},
		{
			"id":           "alert_002",
			"type":         "DAILY_LOSS",
			"severity":     "ERROR",
			"title":        "Daily loss limit exceeded",
			"description":  "Daily loss limit has been exceeded",
			"value":        "25000",
			"threshold":    "20000",
			"timestamp":    time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
			"acknowledged": true,
			"resolved":     false,
		},
	}

	// Apply filters (simplified)
	filteredAlerts := []map[string]interface{}{}
	for _, alert := range alerts {
		if alertType != "" && alert["type"] != alertType {
			continue
		}
		if severity != "" && alert["severity"] != severity {
			continue
		}
		if resolved != "" {
			isResolved := alert["resolved"].(bool)
			if (resolved == "true" && !isResolved) || (resolved == "false" && isResolved) {
				continue
			}
		}
		filteredAlerts = append(filteredAlerts, alert)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"alerts": filteredAlerts,
		"total":  len(filteredAlerts),
	})
}

// GetComplianceDashboard returns compliance dashboard data
func (h *ComplianceHandlers) GetComplianceDashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info(ctx, "Getting compliance dashboard", nil)

	dashboard := map[string]interface{}{
		"overview": map[string]interface{}{
			"total_frameworks":        2,
			"compliant_frameworks":    1,
			"total_requirements":      20,
			"compliant_requirements":  16,
			"overall_compliance_rate": 80.0,
			"risk_level":              "MEDIUM",
		},
		"recent_violations": []map[string]interface{}{
			{
				"id":          "violation_001",
				"type":        "RISK",
				"severity":    "MEDIUM",
				"description": "Position limit exceeded",
				"timestamp":   time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
			},
		},
		"active_alerts": []map[string]interface{}{
			{
				"id":        "alert_001",
				"type":      "POSITION_LIMIT",
				"severity":  "WARNING",
				"title":     "Position limit approaching",
				"timestamp": time.Now().Add(-30 * time.Minute).Format(time.RFC3339),
			},
		},
		"upcoming_reports": []map[string]interface{}{
			{
				"id":       "scheduled_001",
				"name":     "Monthly BSA Report",
				"due_date": time.Now().AddDate(0, 0, 5).Format(time.RFC3339),
				"status":   "PENDING",
			},
		},
		"metrics": map[string]interface{}{
			"compliance_score":      85.0,
			"risk_score":            25.0,
			"violations_this_month": 5,
			"resolved_violations":   4,
			"pending_violations":    1,
			"reports_generated":     12,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dashboard)
}

// GetComplianceOverview returns high-level compliance overview
func (h *ComplianceHandlers) GetComplianceOverview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info(ctx, "Getting compliance overview", nil)

	overview := map[string]interface{}{
		"status":                  "COMPLIANT",
		"overall_compliance_rate": 85.0,
		"risk_level":              "MEDIUM",
		"total_frameworks":        2,
		"compliant_frameworks":    1,
		"active_violations":       1,
		"pending_reports":         2,
		"last_assessment":         time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
		"next_review":             time.Now().AddDate(0, 1, 0).Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(overview)
}

// GetViolations returns compliance violations
func (h *ComplianceHandlers) GetViolations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info(ctx, "Getting compliance violations", nil)

	// Parse query parameters
	query := r.URL.Query()
	violationType := query.Get("type")
	severity := query.Get("severity")
	resolved := query.Get("resolved")

	// Mock violations data
	violations := []map[string]interface{}{
		{
			"id":          "violation_001",
			"type":        "RISK",
			"severity":    "MEDIUM",
			"description": "Position limit exceeded for BTCUSD",
			"framework":   "us_bsa",
			"timestamp":   time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
			"resolved":    true,
			"resolved_at": time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
		},
		{
			"id":          "violation_002",
			"type":        "COMPLIANCE",
			"severity":    "HIGH",
			"description": "KYC documentation incomplete",
			"framework":   "eu_amld5",
			"timestamp":   time.Now().Add(-4 * time.Hour).Format(time.RFC3339),
			"resolved":    false,
		},
	}

	// Apply filters
	filteredViolations := []map[string]interface{}{}
	for _, violation := range violations {
		if violationType != "" && violation["type"] != violationType {
			continue
		}
		if severity != "" && violation["severity"] != severity {
			continue
		}
		if resolved != "" {
			isResolved := violation["resolved"].(bool)
			if (resolved == "true" && !isResolved) || (resolved == "false" && isResolved) {
				continue
			}
		}
		filteredViolations = append(filteredViolations, violation)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"violations": filteredViolations,
		"total":      len(filteredViolations),
	})
}

// GetViolation returns a specific compliance violation
func (h *ComplianceHandlers) GetViolation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	violationID := vars["id"]

	h.logger.Info(ctx, "Getting compliance violation", map[string]interface{}{
		"violation_id": violationID,
	})

	// Mock violation data
	violation := map[string]interface{}{
		"id":          violationID,
		"type":        "RISK",
		"severity":    "MEDIUM",
		"description": "Position limit exceeded for BTCUSD",
		"framework":   "us_bsa",
		"requirement": "position_limit",
		"details": map[string]interface{}{
			"symbol":        "BTCUSD",
			"position_size": 1500000,
			"limit":         1000000,
			"excess":        500000,
		},
		"timestamp":   time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
		"resolved":    false,
		"remediation": "Reduce position size to comply with limits",
		"responsible": "Risk Management Team",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(violation)
}

// ResolveViolation resolves a compliance violation
func (h *ComplianceHandlers) ResolveViolation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	violationID := vars["id"]

	h.logger.Info(ctx, "Resolving compliance violation", map[string]interface{}{
		"violation_id": violationID,
	})

	var request struct {
		Resolution string `json:"resolution"`
		Notes      string `json:"notes"`
		ResolvedBy string `json:"resolved_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Mock resolution
	response := map[string]interface{}{
		"id":          violationID,
		"resolved":    true,
		"resolved_at": time.Now().Format(time.RFC3339),
		"resolved_by": request.ResolvedBy,
		"resolution":  request.Resolution,
		"notes":       request.Notes,
		"status":      "RESOLVED",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetRiskThresholds returns risk thresholds
func (h *ComplianceHandlers) GetRiskThresholds(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info(ctx, "Getting risk thresholds", nil)

	// Mock thresholds data
	thresholds := []map[string]interface{}{
		{
			"id":             "threshold_001",
			"type":           "DAILY_LOSS",
			"warning_level":  "10000",
			"error_level":    "25000",
			"critical_level": "50000",
			"enabled":        true,
		},
		{
			"id":             "threshold_002",
			"type":           "LEVERAGE",
			"warning_level":  "2.0",
			"error_level":    "5.0",
			"critical_level": "10.0",
			"enabled":        true,
		},
		{
			"id":             "threshold_003",
			"type":           "CONCENTRATION",
			"warning_level":  "0.3",
			"error_level":    "0.5",
			"critical_level": "0.7",
			"enabled":        true,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"thresholds": thresholds,
		"total":      len(thresholds),
	})
}

// CreateRiskThreshold creates a new risk threshold
func (h *ComplianceHandlers) CreateRiskThreshold(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info(ctx, "Creating risk threshold", nil)

	var request struct {
		Type          string `json:"type"`
		Symbol        string `json:"symbol,omitempty"`
		WarningLevel  string `json:"warning_level"`
		ErrorLevel    string `json:"error_level"`
		CriticalLevel string `json:"critical_level"`
		Enabled       bool   `json:"enabled"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Mock threshold creation
	thresholdID := "threshold_" + strconv.FormatInt(time.Now().Unix(), 10)

	threshold := map[string]interface{}{
		"id":             thresholdID,
		"type":           request.Type,
		"symbol":         request.Symbol,
		"warning_level":  request.WarningLevel,
		"error_level":    request.ErrorLevel,
		"critical_level": request.CriticalLevel,
		"enabled":        request.Enabled,
		"created_at":     time.Now().Format(time.RFC3339),
		"updated_at":     time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(threshold)
}

// GetAuditEvents returns audit events
func (h *ComplianceHandlers) GetAuditEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info(ctx, "Getting audit events", nil)

	// Parse query parameters
	query := r.URL.Query()
	eventType := query.Get("type")
	riskLevel := query.Get("risk_level")
	limit := query.Get("limit")

	limitInt := 50
	if limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			limitInt = l
		}
	}

	// Mock audit events
	events := []map[string]interface{}{
		{
			"id":         "event_001",
			"type":       "TRADE_EXECUTION",
			"risk_level": "LOW",
			"user_id":    "user_123",
			"action":     "BUY",
			"details": map[string]interface{}{
				"symbol":   "BTCUSD",
				"quantity": "0.5",
				"price":    "45000",
			},
			"timestamp":  time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
			"ip_address": "192.168.1.100",
			"framework":  "us_bsa",
		},
		{
			"id":         "event_002",
			"type":       "COMPLIANCE_CHECK",
			"risk_level": "MEDIUM",
			"user_id":    "system",
			"action":     "KYC_VERIFICATION",
			"details": map[string]interface{}{
				"customer_id": "cust_456",
				"status":      "PENDING",
			},
			"timestamp":  time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
			"ip_address": "internal",
			"framework":  "eu_amld5",
		},
		{
			"id":         "event_003",
			"type":       "RISK_ALERT",
			"risk_level": "HIGH",
			"user_id":    "system",
			"action":     "POSITION_LIMIT_BREACH",
			"details": map[string]interface{}{
				"symbol":        "ETHUSD",
				"position_size": "1500000",
				"limit":         "1000000",
			},
			"timestamp":  time.Now().Add(-30 * time.Minute).Format(time.RFC3339),
			"ip_address": "internal",
			"framework":  "internal",
		},
	}

	// Apply filters
	filteredEvents := []map[string]interface{}{}
	for _, event := range events {
		if eventType != "" && event["type"] != eventType {
			continue
		}
		if riskLevel != "" && event["risk_level"] != riskLevel {
			continue
		}
		filteredEvents = append(filteredEvents, event)
		if len(filteredEvents) >= limitInt {
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"events": filteredEvents,
		"total":  len(filteredEvents),
		"limit":  limitInt,
	})
}

// ExportAuditEvents exports audit events
func (h *ComplianceHandlers) ExportAuditEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info(ctx, "Exporting audit events", nil)

	format := r.URL.Query().Get("format")
	if format == "" {
		format = "csv"
	}

	switch format {
	case "csv":
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=audit_events.csv")
		w.Write([]byte("ID,Type,Risk Level,User ID,Action,Timestamp\n"))
		w.Write([]byte("event_001,TRADE_EXECUTION,LOW,user_123,BUY," + time.Now().Format(time.RFC3339) + "\n"))
	case "json":
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", "attachment; filename=audit_events.json")
		w.Write([]byte(`{"events": [{"id": "event_001", "type": "TRADE_EXECUTION"}]}`))
	default:
		http.Error(w, "Unsupported format", http.StatusBadRequest)
		return
	}
}

// GetAuditSummary returns audit summary
func (h *ComplianceHandlers) GetAuditSummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info(ctx, "Getting audit summary", nil)

	summary := map[string]interface{}{
		"total_events":      1247,
		"events_today":      89,
		"high_risk_events":  12,
		"failed_logins":     3,
		"compliance_checks": 45,
		"risk_alerts":       8,
		"data_exports":      2,
		"by_type": map[string]interface{}{
			"TRADE_EXECUTION":  856,
			"COMPLIANCE_CHECK": 234,
			"RISK_ALERT":       89,
			"USER_LOGIN":       45,
			"DATA_EXPORT":      23,
		},
		"by_risk_level": map[string]interface{}{
			"LOW":      987,
			"MEDIUM":   189,
			"HIGH":     56,
			"CRITICAL": 15,
		},
		"last_updated": time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// GetAlerts returns alerts
func (h *ComplianceHandlers) GetAlerts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info(ctx, "Getting alerts", nil)

	// Parse query parameters
	query := r.URL.Query()
	alertType := query.Get("type")
	severity := query.Get("severity")
	acknowledged := query.Get("acknowledged")

	// Mock alerts data
	alerts := []map[string]interface{}{
		{
			"id":           "alert_001",
			"type":         "COMPLIANCE",
			"severity":     "HIGH",
			"title":        "KYC Documentation Missing",
			"description":  "Customer KYC documentation is incomplete",
			"timestamp":    time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
			"acknowledged": false,
			"resolved":     false,
		},
		{
			"id":           "alert_002",
			"type":         "RISK",
			"severity":     "MEDIUM",
			"title":        "Position Limit Approaching",
			"description":  "BTCUSD position approaching limit",
			"timestamp":    time.Now().Add(-30 * time.Minute).Format(time.RFC3339),
			"acknowledged": true,
			"resolved":     false,
		},
	}

	// Apply filters
	filteredAlerts := []map[string]interface{}{}
	for _, alert := range alerts {
		if alertType != "" && alert["type"] != alertType {
			continue
		}
		if severity != "" && alert["severity"] != severity {
			continue
		}
		if acknowledged != "" {
			isAcknowledged := alert["acknowledged"].(bool)
			if (acknowledged == "true" && !isAcknowledged) || (acknowledged == "false" && isAcknowledged) {
				continue
			}
		}
		filteredAlerts = append(filteredAlerts, alert)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"alerts": filteredAlerts,
		"total":  len(filteredAlerts),
	})
}

// AcknowledgeAlert acknowledges an alert
func (h *ComplianceHandlers) AcknowledgeAlert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	alertID := vars["id"]

	h.logger.Info(ctx, "Acknowledging alert", map[string]interface{}{
		"alert_id": alertID,
	})

	var request struct {
		AcknowledgedBy string `json:"acknowledged_by"`
		Notes          string `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"id":              alertID,
		"acknowledged":    true,
		"acknowledged_at": time.Now().Format(time.RFC3339),
		"acknowledged_by": request.AcknowledgedBy,
		"notes":           request.Notes,
		"status":          "ACKNOWLEDGED",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ResolveAlert resolves an alert
func (h *ComplianceHandlers) ResolveAlert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	alertID := vars["id"]

	h.logger.Info(ctx, "Resolving alert", map[string]interface{}{
		"alert_id": alertID,
	})

	var request struct {
		Resolution string `json:"resolution"`
		ResolvedBy string `json:"resolved_by"`
		Notes      string `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"id":          alertID,
		"resolved":    true,
		"resolved_at": time.Now().Format(time.RFC3339),
		"resolved_by": request.ResolvedBy,
		"resolution":  request.Resolution,
		"notes":       request.Notes,
		"status":      "RESOLVED",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAlertRules returns alert rules
func (h *ComplianceHandlers) GetAlertRules(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info(ctx, "Getting alert rules", nil)

	rules := []map[string]interface{}{
		{
			"id":          "rule_001",
			"name":        "High Risk Transaction",
			"description": "Alert on transactions above $100k",
			"type":        "TRANSACTION",
			"enabled":     true,
			"conditions": map[string]interface{}{
				"amount_threshold": "100000",
				"currency":         "USD",
			},
		},
		{
			"id":          "rule_002",
			"name":        "KYC Expiration",
			"description": "Alert when KYC documents expire",
			"type":        "COMPLIANCE",
			"enabled":     true,
			"conditions": map[string]interface{}{
				"days_before_expiry": "30",
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"rules": rules,
		"total": len(rules),
	})
}

// GetAlertChannels returns alert channels
func (h *ComplianceHandlers) GetAlertChannels(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Info(ctx, "Getting alert channels", nil)

	channels := []map[string]interface{}{
		{
			"id":      "channel_001",
			"name":    "Email - Compliance Team",
			"type":    "EMAIL",
			"enabled": true,
			"config": map[string]interface{}{
				"recipients": []string{"compliance@example.com"},
			},
		},
		{
			"id":      "channel_002",
			"name":    "Slack - Risk Channel",
			"type":    "SLACK",
			"enabled": true,
			"config": map[string]interface{}{
				"webhook_url": "https://hooks.slack.com/...",
				"channel":     "#risk-alerts",
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"channels": channels,
		"total":    len(channels),
	})
}
