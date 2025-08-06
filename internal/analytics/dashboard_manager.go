package analytics

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

// DashboardManager manages real-time dashboards
type DashboardManager struct {
	logger     *observability.Logger
	config     *AnalyticsConfig
	dashboards map[string]*Dashboard
	widgets    map[string]*Widget
	layouts    map[string]*DashboardLayout
	themes     map[string]*DashboardTheme
	mu         sync.RWMutex
}

// Dashboard represents a real-time dashboard
type Dashboard struct {
	DashboardID string                 `json:"dashboard_id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Category    DashboardCategory      `json:"category"`
	Layout      *DashboardLayout       `json:"layout"`
	Widgets     []*Widget              `json:"widgets"`
	Theme       string                 `json:"theme"`
	RefreshRate time.Duration          `json:"refresh_rate"`
	AutoRefresh bool                   `json:"auto_refresh"`
	Permissions *DashboardPermissions  `json:"permissions"`
	Tags        []string               `json:"tags"`
	CreatedBy   string                 `json:"created_by"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	LastViewed  *time.Time             `json:"last_viewed,omitempty"`
	ViewCount   int64                  `json:"view_count"`
	IsPublic    bool                   `json:"is_public"`
	IsTemplate  bool                   `json:"is_template"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// DashboardCategory defines dashboard categories
type DashboardCategory string

const (
	CategoryOverview    DashboardCategory = "overview"
	CategoryPerformance DashboardCategory = "performance"
	CategorySecurity    DashboardCategory = "security"
	CategoryBusiness    DashboardCategory = "business"
	CategoryTrading     DashboardCategory = "trading"
	CategoryCompliance  DashboardCategory = "compliance"
	CategoryCustom      DashboardCategory = "custom"
)

// DashboardLayout defines dashboard layout
type DashboardLayout struct {
	Type        LayoutType     `json:"type"`
	Columns     int            `json:"columns"`
	Rows        int            `json:"rows"`
	GridSize    GridSize       `json:"grid_size"`
	Responsive  bool           `json:"responsive"`
	Breakpoints map[string]int `json:"breakpoints"`
}

// LayoutType defines layout types
type LayoutType string

const (
	LayoutTypeGrid    LayoutType = "grid"
	LayoutTypeFlex    LayoutType = "flex"
	LayoutTypeMasonry LayoutType = "masonry"
	LayoutTypeCustom  LayoutType = "custom"
)

// GridSize defines grid dimensions
type GridSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Widget represents a dashboard widget
type Widget struct {
	WidgetID      string                 `json:"widget_id"`
	Name          string                 `json:"name"`
	Type          WidgetType             `json:"type"`
	Position      WidgetPosition         `json:"position"`
	Size          WidgetSize             `json:"size"`
	Configuration *WidgetConfiguration   `json:"configuration"`
	DataSource    *WidgetDataSource      `json:"data_source"`
	Visualization *WidgetVisualization   `json:"visualization"`
	Filters       []WidgetFilter         `json:"filters"`
	RefreshRate   time.Duration          `json:"refresh_rate"`
	LastUpdated   time.Time              `json:"last_updated"`
	IsVisible     bool                   `json:"is_visible"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// WidgetType defines widget types
type WidgetType string

const (
	WidgetTypeMetric WidgetType = "metric"
	WidgetTypeChart  WidgetType = "chart"
	WidgetTypeTable  WidgetType = "table"
	WidgetTypeGauge  WidgetType = "gauge"
	WidgetTypeMap    WidgetType = "map"
	WidgetTypeText   WidgetType = "text"
	WidgetTypeImage  WidgetType = "image"
	WidgetTypeAlert  WidgetType = "alert"
	WidgetTypeLog    WidgetType = "log"
	WidgetTypeCustom WidgetType = "custom"
)

// WidgetPosition defines widget position
type WidgetPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// WidgetSize defines widget size
type WidgetSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// WidgetConfiguration contains widget configuration
type WidgetConfiguration struct {
	Title       string                 `json:"title"`
	Subtitle    string                 `json:"subtitle,omitempty"`
	ShowTitle   bool                   `json:"show_title"`
	ShowLegend  bool                   `json:"show_legend"`
	ShowTooltip bool                   `json:"show_tooltip"`
	Colors      []string               `json:"colors,omitempty"`
	Thresholds  []WidgetThreshold      `json:"thresholds,omitempty"`
	Format      string                 `json:"format,omitempty"`
	Unit        string                 `json:"unit,omitempty"`
	Decimals    int                    `json:"decimals"`
	Options     map[string]interface{} `json:"options,omitempty"`
}

// WidgetThreshold defines widget thresholds
type WidgetThreshold struct {
	Value float64 `json:"value"`
	Color string  `json:"color"`
	Label string  `json:"label,omitempty"`
}

// WidgetDataSource defines widget data source
type WidgetDataSource struct {
	Type        DataSourceType         `json:"type"`
	MetricName  string                 `json:"metric_name,omitempty"`
	Query       string                 `json:"query,omitempty"`
	TimeRange   TimeRange              `json:"time_range"`
	Aggregation string                 `json:"aggregation,omitempty"`
	GroupBy     []string               `json:"group_by,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

// DataSourceType defines data source types
type DataSourceType string

const (
	DataSourceTypeMetrics     DataSourceType = "metrics"
	DataSourceTypeLogs        DataSourceType = "logs"
	DataSourceTypeEvents      DataSourceType = "events"
	DataSourceTypeAlerts      DataSourceType = "alerts"
	DataSourceTypeAnomalies   DataSourceType = "anomalies"
	DataSourceTypePredictions DataSourceType = "predictions"
	DataSourceTypeCustom      DataSourceType = "custom"
)

// TimeRange defines time range for data
type TimeRange struct {
	From     time.Time `json:"from"`
	To       time.Time `json:"to"`
	Relative string    `json:"relative,omitempty"` // e.g., "1h", "24h", "7d"
}

// WidgetVisualization defines widget visualization
type WidgetVisualization struct {
	ChartType  ChartType              `json:"chart_type,omitempty"`
	SeriesType SeriesType             `json:"series_type,omitempty"`
	XAxis      AxisConfiguration      `json:"x_axis,omitempty"`
	YAxis      AxisConfiguration      `json:"y_axis,omitempty"`
	Animation  bool                   `json:"animation"`
	Stacked    bool                   `json:"stacked,omitempty"`
	Smooth     bool                   `json:"smooth,omitempty"`
	FillArea   bool                   `json:"fill_area,omitempty"`
	Options    map[string]interface{} `json:"options,omitempty"`
}

// ChartType defines chart types
type ChartType string

const (
	ChartTypeLine      ChartType = "line"
	ChartTypeArea      ChartType = "area"
	ChartTypeBar       ChartType = "bar"
	ChartTypeColumn    ChartType = "column"
	ChartTypePie       ChartType = "pie"
	ChartTypeDoughnut  ChartType = "doughnut"
	ChartTypeScatter   ChartType = "scatter"
	ChartTypeHeatmap   ChartType = "heatmap"
	ChartTypeGauge     ChartType = "gauge"
	ChartTypeSparkline ChartType = "sparkline"
)

// SeriesType defines series types
type SeriesType string

const (
	SeriesTypeSingle   SeriesType = "single"
	SeriesTypeMultiple SeriesType = "multiple"
	SeriesTypeStacked  SeriesType = "stacked"
)

// AxisConfiguration defines axis configuration
type AxisConfiguration struct {
	Label    string   `json:"label,omitempty"`
	Min      *float64 `json:"min,omitempty"`
	Max      *float64 `json:"max,omitempty"`
	LogScale bool     `json:"log_scale,omitempty"`
	Format   string   `json:"format,omitempty"`
	Unit     string   `json:"unit,omitempty"`
}

// WidgetFilter defines widget filters
type WidgetFilter struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
	Enabled  bool        `json:"enabled"`
}

// DashboardPermissions defines dashboard permissions
type DashboardPermissions struct {
	Owner   string   `json:"owner"`
	Viewers []string `json:"viewers"`
	Editors []string `json:"editors"`
	Public  bool     `json:"public"`
}

// DashboardTheme defines dashboard theme
type DashboardTheme struct {
	ThemeID      string          `json:"theme_id"`
	Name         string          `json:"name"`
	Colors       ThemeColors     `json:"colors"`
	Typography   ThemeTypography `json:"typography"`
	Spacing      ThemeSpacing    `json:"spacing"`
	BorderRadius int             `json:"border_radius"`
	Shadows      bool            `json:"shadows"`
	CustomCSS    string          `json:"custom_css,omitempty"`
}

// ThemeColors defines theme colors
type ThemeColors struct {
	Primary       string `json:"primary"`
	Secondary     string `json:"secondary"`
	Background    string `json:"background"`
	Surface       string `json:"surface"`
	Text          string `json:"text"`
	TextSecondary string `json:"text_secondary"`
	Border        string `json:"border"`
	Success       string `json:"success"`
	Warning       string `json:"warning"`
	Error         string `json:"error"`
	Info          string `json:"info"`
}

// ThemeTypography defines theme typography
type ThemeTypography struct {
	FontFamily string  `json:"font_family"`
	FontSize   int     `json:"font_size"`
	LineHeight float64 `json:"line_height"`
}

// ThemeSpacing defines theme spacing
type ThemeSpacing struct {
	Small  int `json:"small"`
	Medium int `json:"medium"`
	Large  int `json:"large"`
}

// NewDashboardManager creates a new dashboard manager
func NewDashboardManager(logger *observability.Logger, config *AnalyticsConfig) *DashboardManager {
	dm := &DashboardManager{
		logger:     logger,
		config:     config,
		dashboards: make(map[string]*Dashboard),
		widgets:    make(map[string]*Widget),
		layouts:    make(map[string]*DashboardLayout),
		themes:     make(map[string]*DashboardTheme),
	}

	// Initialize default themes
	dm.initializeDefaultThemes()

	return dm
}

// Start starts the dashboard manager
func (dm *DashboardManager) Start(ctx context.Context) error {
	dm.logger.Info(ctx, "Starting dashboard manager", nil)

	// Initialize default dashboards
	dm.initializeDefaultDashboards()

	// Start background processes
	go dm.updateDashboards(ctx)

	return nil
}

// CreateDashboard creates a new dashboard
func (dm *DashboardManager) CreateDashboard(dashboard *Dashboard) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if dashboard.DashboardID == "" {
		dashboard.DashboardID = uuid.New().String()
	}

	dashboard.CreatedAt = time.Now()
	dashboard.UpdatedAt = time.Now()

	dm.dashboards[dashboard.DashboardID] = dashboard

	dm.logger.Info(context.Background(), "Dashboard created", map[string]interface{}{
		"dashboard_id": dashboard.DashboardID,
		"name":         dashboard.Name,
		"category":     dashboard.Category,
		"widget_count": len(dashboard.Widgets),
	})

	return nil
}

// GetDashboard retrieves a dashboard by ID
func (dm *DashboardManager) GetDashboard(dashboardID string) (*Dashboard, error) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	dashboard, exists := dm.dashboards[dashboardID]
	if !exists {
		return nil, fmt.Errorf("dashboard not found: %s", dashboardID)
	}

	// Update view statistics
	dashboard.ViewCount++
	now := time.Now()
	dashboard.LastViewed = &now

	return dashboard, nil
}

// GetDashboards retrieves all dashboards
func (dm *DashboardManager) GetDashboards() []*Dashboard {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	dashboards := make([]*Dashboard, 0, len(dm.dashboards))
	for _, dashboard := range dm.dashboards {
		dashboards = append(dashboards, dashboard)
	}

	return dashboards
}

// GetDashboardsByCategory retrieves dashboards by category
func (dm *DashboardManager) GetDashboardsByCategory(category DashboardCategory) []*Dashboard {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	dashboards := make([]*Dashboard, 0)
	for _, dashboard := range dm.dashboards {
		if dashboard.Category == category {
			dashboards = append(dashboards, dashboard)
		}
	}

	return dashboards
}

// UpdateDashboard updates a dashboard
func (dm *DashboardManager) UpdateDashboard(dashboard *Dashboard) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if _, exists := dm.dashboards[dashboard.DashboardID]; !exists {
		return fmt.Errorf("dashboard not found: %s", dashboard.DashboardID)
	}

	dashboard.UpdatedAt = time.Now()
	dm.dashboards[dashboard.DashboardID] = dashboard

	dm.logger.Info(context.Background(), "Dashboard updated", map[string]interface{}{
		"dashboard_id": dashboard.DashboardID,
		"name":         dashboard.Name,
	})

	return nil
}

// DeleteDashboard deletes a dashboard
func (dm *DashboardManager) DeleteDashboard(dashboardID string) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if _, exists := dm.dashboards[dashboardID]; !exists {
		return fmt.Errorf("dashboard not found: %s", dashboardID)
	}

	delete(dm.dashboards, dashboardID)

	dm.logger.Info(context.Background(), "Dashboard deleted", map[string]interface{}{
		"dashboard_id": dashboardID,
	})

	return nil
}

// initializeDefaultThemes initializes default dashboard themes
func (dm *DashboardManager) initializeDefaultThemes() {
	// Light theme
	lightTheme := &DashboardTheme{
		ThemeID: "light",
		Name:    "Light Theme",
		Colors: ThemeColors{
			Primary:       "#1976d2",
			Secondary:     "#dc004e",
			Background:    "#ffffff",
			Surface:       "#f5f5f5",
			Text:          "#212121",
			TextSecondary: "#757575",
			Border:        "#e0e0e0",
			Success:       "#4caf50",
			Warning:       "#ff9800",
			Error:         "#f44336",
			Info:          "#2196f3",
		},
		Typography: ThemeTypography{
			FontFamily: "Roboto, sans-serif",
			FontSize:   14,
			LineHeight: 1.5,
		},
		Spacing: ThemeSpacing{
			Small:  8,
			Medium: 16,
			Large:  24,
		},
		BorderRadius: 4,
		Shadows:      true,
	}

	// Dark theme
	darkTheme := &DashboardTheme{
		ThemeID: "dark",
		Name:    "Dark Theme",
		Colors: ThemeColors{
			Primary:       "#90caf9",
			Secondary:     "#f48fb1",
			Background:    "#121212",
			Surface:       "#1e1e1e",
			Text:          "#ffffff",
			TextSecondary: "#b0b0b0",
			Border:        "#333333",
			Success:       "#81c784",
			Warning:       "#ffb74d",
			Error:         "#e57373",
			Info:          "#64b5f6",
		},
		Typography: ThemeTypography{
			FontFamily: "Roboto, sans-serif",
			FontSize:   14,
			LineHeight: 1.5,
		},
		Spacing: ThemeSpacing{
			Small:  8,
			Medium: 16,
			Large:  24,
		},
		BorderRadius: 4,
		Shadows:      true,
	}

	dm.themes["light"] = lightTheme
	dm.themes["dark"] = darkTheme
}

// initializeDefaultDashboards initializes default dashboards
func (dm *DashboardManager) initializeDefaultDashboards() {
	// System Overview Dashboard
	systemDashboard := &Dashboard{
		Name:        "System Overview",
		Description: "Real-time system performance and health metrics",
		Category:    CategoryOverview,
		Layout: &DashboardLayout{
			Type:       LayoutTypeGrid,
			Columns:    12,
			Rows:       8,
			GridSize:   GridSize{Width: 100, Height: 100},
			Responsive: true,
			Breakpoints: map[string]int{
				"sm": 576,
				"md": 768,
				"lg": 992,
				"xl": 1200,
			},
		},
		Theme:       "light",
		RefreshRate: 30 * time.Second,
		AutoRefresh: true,
		Permissions: &DashboardPermissions{
			Owner:  "system",
			Public: true,
		},
		Tags:      []string{"system", "overview", "performance"},
		CreatedBy: "system",
		IsPublic:  true,
	}

	// Add widgets to system dashboard
	systemDashboard.Widgets = []*Widget{
		{
			Name:     "CPU Usage",
			Type:     WidgetTypeGauge,
			Position: WidgetPosition{X: 0, Y: 0},
			Size:     WidgetSize{Width: 3, Height: 3},
			Configuration: &WidgetConfiguration{
				Title:     "CPU Usage",
				ShowTitle: true,
				Unit:      "%",
				Decimals:  1,
				Thresholds: []WidgetThreshold{
					{Value: 70, Color: "#ff9800", Label: "Warning"},
					{Value: 90, Color: "#f44336", Label: "Critical"},
				},
			},
			DataSource: &WidgetDataSource{
				Type:        DataSourceTypeMetrics,
				MetricName:  "cpu_usage",
				TimeRange:   TimeRange{Relative: "5m"},
				Aggregation: "avg",
			},
			Visualization: &WidgetVisualization{
				ChartType: ChartTypeGauge,
				Animation: true,
			},
			RefreshRate: 10 * time.Second,
			IsVisible:   true,
		},
		{
			Name:     "Memory Usage",
			Type:     WidgetTypeGauge,
			Position: WidgetPosition{X: 3, Y: 0},
			Size:     WidgetSize{Width: 3, Height: 3},
			Configuration: &WidgetConfiguration{
				Title:     "Memory Usage",
				ShowTitle: true,
				Unit:      "%",
				Decimals:  1,
				Thresholds: []WidgetThreshold{
					{Value: 80, Color: "#ff9800", Label: "Warning"},
					{Value: 95, Color: "#f44336", Label: "Critical"},
				},
			},
			DataSource: &WidgetDataSource{
				Type:        DataSourceTypeMetrics,
				MetricName:  "memory_usage",
				TimeRange:   TimeRange{Relative: "5m"},
				Aggregation: "avg",
			},
			Visualization: &WidgetVisualization{
				ChartType: ChartTypeGauge,
				Animation: true,
			},
			RefreshRate: 10 * time.Second,
			IsVisible:   true,
		},
	}

	// Create dashboards
	dm.CreateDashboard(systemDashboard)
}

// updateDashboards updates dashboard data periodically
func (dm *DashboardManager) updateDashboards(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			dm.refreshDashboardData()
		}
	}
}

// refreshDashboardData refreshes data for all dashboards
func (dm *DashboardManager) refreshDashboardData() {
	dm.mu.RLock()
	dashboards := make([]*Dashboard, 0, len(dm.dashboards))
	for _, dashboard := range dm.dashboards {
		if dashboard.AutoRefresh {
			dashboards = append(dashboards, dashboard)
		}
	}
	dm.mu.RUnlock()

	for _, dashboard := range dashboards {
		dm.refreshDashboard(dashboard)
	}
}

// refreshDashboard refreshes data for a specific dashboard
func (dm *DashboardManager) refreshDashboard(dashboard *Dashboard) {
	for _, widget := range dashboard.Widgets {
		if widget.IsVisible {
			dm.refreshWidget(widget)
		}
	}
}

// refreshWidget refreshes data for a specific widget
func (dm *DashboardManager) refreshWidget(widget *Widget) {
	// Simulate data refresh
	widget.LastUpdated = time.Now()

	dm.logger.Debug(context.Background(), "Widget data refreshed", map[string]interface{}{
		"widget_id":   widget.WidgetID,
		"widget_name": widget.Name,
		"widget_type": widget.Type,
	})
}

// ExportDashboard exports a dashboard configuration
func (dm *DashboardManager) ExportDashboard(dashboardID string) ([]byte, error) {
	dashboard, err := dm.GetDashboard(dashboardID)
	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(dashboard, "", "  ")
}

// ImportDashboard imports a dashboard configuration
func (dm *DashboardManager) ImportDashboard(data []byte) (*Dashboard, error) {
	var dashboard Dashboard
	if err := json.Unmarshal(data, &dashboard); err != nil {
		return nil, fmt.Errorf("failed to parse dashboard: %w", err)
	}

	// Generate new ID to avoid conflicts
	dashboard.DashboardID = uuid.New().String()

	if err := dm.CreateDashboard(&dashboard); err != nil {
		return nil, err
	}

	return &dashboard, nil
}

// GetDashboardMetrics returns dashboard usage metrics
func (dm *DashboardManager) GetDashboardMetrics() map[string]interface{} {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	totalViews := int64(0)
	categoryCount := make(map[DashboardCategory]int)

	for _, dashboard := range dm.dashboards {
		totalViews += dashboard.ViewCount
		categoryCount[dashboard.Category]++
	}

	return map[string]interface{}{
		"total_dashboards": len(dm.dashboards),
		"total_views":      totalViews,
		"category_count":   categoryCount,
		"total_widgets":    len(dm.widgets),
		"total_themes":     len(dm.themes),
	}
}
