package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/ai-agentic-browser/internal/analytics"
	"github.com/ai-agentic-browser/pkg/observability"
)

// AnalyticsMonitoringDemo demonstrates advanced real-time analytics and monitoring
func main() {
	fmt.Println("🔍 AI-Agentic Crypto Browser - Real-time Analytics & Monitoring Demo")
	fmt.Println("====================================================================")

	ctx := context.Background()
	logger := &observability.Logger{}

	// Demo 1: Real-time Analytics Engine
	fmt.Println("\n📊 Demo 1: Real-time Analytics Engine")
	demoAnalyticsEngine(ctx, logger)

	// Demo 2: Anomaly Detection System
	fmt.Println("\n🚨 Demo 2: Anomaly Detection System")
	demoAnomalyDetection(ctx, logger)

	// Demo 3: Predictive Analytics
	fmt.Println("\n🔮 Demo 3: Predictive Analytics")
	demoPredictiveAnalytics(ctx, logger)

	// Demo 4: Intelligent Alerting
	fmt.Println("\n🔔 Demo 4: Intelligent Alerting")
	demoIntelligentAlerting(ctx, logger)

	// Demo 5: Real-time Dashboards
	fmt.Println("\n📈 Demo 5: Real-time Dashboards")
	demoRealTimeDashboards(ctx, logger)

	// Demo 6: Integrated Analytics Workflow
	fmt.Println("\n🔄 Demo 6: Integrated Analytics Workflow")
	demoIntegratedWorkflow(ctx, logger)

	fmt.Println("\n🎉 Real-time Analytics & Monitoring Demo Complete!")
	fmt.Println("All enterprise-grade analytics and monitoring features are operational.")
}

// demoAnalyticsEngine demonstrates the real-time analytics engine
func demoAnalyticsEngine(ctx context.Context, logger *observability.Logger) {
	fmt.Println("  Creating real-time analytics engine...")

	// Create analytics configuration
	config := &analytics.AnalyticsConfig{
		EnableRealTimeProcessing:    true,
		EnablePredictiveAnalytics:   true,
		EnableAnomalyDetection:      true,
		EnableIntelligentAlerting:   true,
		ProcessingInterval:          100 * time.Millisecond,
		MetricsRetentionPeriod:      24 * time.Hour,
		AnomalyDetectionSensitivity: 0.8,
		PredictionHorizon:           1 * time.Hour,
		MaxConcurrentStreams:        100,
		BufferSize:                  10000,
		EnableDataCompression:       true,
		EnableDataEncryption:        false,
	}

	// Create analytics engine
	engine := analytics.NewRealTimeAnalyticsEngine(logger, config)
	if err := engine.Start(ctx); err != nil {
		log.Printf("    ❌ Error starting analytics engine: %v", err)
		return
	}

	fmt.Printf("    ✅ Analytics engine started with %d max streams\n", config.MaxConcurrentStreams)

	// Create data streams
	tradingStream, err := engine.CreateDataStream(
		"Trading Data",
		"trading_system",
		[]analytics.EventType{analytics.EventTypeTradingActivity, analytics.EventTypeMarketData},
		nil,
	)
	if err != nil {
		fmt.Printf("    ❌ Error creating trading stream: %v\n", err)
		return
	}

	systemStream, err := engine.CreateDataStream(
		"System Metrics",
		"system_monitor",
		[]analytics.EventType{analytics.EventTypeSystemMetric, analytics.EventTypePerformance},
		nil,
	)
	if err != nil {
		fmt.Printf("    ❌ Error creating system stream: %v\n", err)
		return
	}

	fmt.Printf("    ✅ Created data streams: %s, %s\n", tradingStream.Name, systemStream.Name)

	// Subscribe to events
	tradingEvents := engine.Subscribe(analytics.EventTypeTradingActivity, 100)
	systemEvents := engine.Subscribe(analytics.EventTypeSystemMetric, 100)

	fmt.Printf("    ✅ Subscribed to event streams\n")

	// Simulate real-time events
	go simulateEvents(engine)

	// Process events for a short time
	fmt.Printf("    📊 Processing real-time events...\n")
	eventCount := 0
	timeout := time.After(5 * time.Second)

	for {
		select {
		case event := <-tradingEvents:
			eventCount++
			fmt.Printf("    📈 Trading Event: %s (Value: %.2f)\n", event.EventType, event.Metrics["price"])
		case event := <-systemEvents:
			eventCount++
			fmt.Printf("    🖥️  System Event: %s (Value: %.2f)\n", event.EventType, event.Metrics["cpu_usage"])
		case <-timeout:
			fmt.Printf("    ✅ Processed %d real-time events\n", eventCount)
			return
		}

		if eventCount >= 10 {
			fmt.Printf("    ✅ Processed %d real-time events\n", eventCount)
			return
		}
	}
}

// demoAnomalyDetection demonstrates anomaly detection capabilities
func demoAnomalyDetection(ctx context.Context, logger *observability.Logger) {
	fmt.Println("  Creating anomaly detection system...")

	config := &analytics.AnalyticsConfig{
		AnomalyDetectionSensitivity: 0.7,
	}

	detector := analytics.NewAnomalyDetector(logger, config)
	if err := detector.Start(ctx); err != nil {
		log.Printf("    ❌ Error starting anomaly detector: %v", err)
		return
	}

	fmt.Printf("    ✅ Anomaly detector started with %.1f sensitivity\n", config.AnomalyDetectionSensitivity)

	// Register custom detectors
	detector.RegisterMetricDetector("cpu_usage", analytics.DetectionMethodZScore, 0.8, 50)
	detector.RegisterMetricDetector("response_time", analytics.DetectionMethodIQR, 0.7, 30)
	detector.RegisterMetricDetector("trading_volume", analytics.DetectionMethodMovingAverage, 0.6, 20)

	fmt.Printf("    ✅ Registered custom metric detectors\n")

	// Simulate normal data
	fmt.Printf("    📊 Feeding normal data points...\n")
	for i := 0; i < 50; i++ {
		// Normal CPU usage (40-60%)
		cpuUsage := 50 + rand.Float64()*10 - 5
		detector.AddDataPoint("cpu_usage", cpuUsage, map[string]string{"host": "server1"})

		// Normal response time (100-200ms)
		responseTime := 150 + rand.Float64()*50 - 25
		detector.AddDataPoint("response_time", responseTime, map[string]string{"endpoint": "/api/data"})

		// Normal trading volume
		volume := 1000 + rand.Float64()*200 - 100
		detector.AddDataPoint("trading_volume", volume, map[string]string{"symbol": "BTC/USD"})

		time.Sleep(10 * time.Millisecond)
	}

	// Introduce anomalies
	fmt.Printf("    🚨 Introducing anomalous data points...\n")

	// CPU spike
	detector.AddDataPoint("cpu_usage", 95.0, map[string]string{"host": "server1"})
	fmt.Printf("    ⚠️  Injected CPU spike: 95%%\n")

	// Response time spike
	detector.AddDataPoint("response_time", 5000.0, map[string]string{"endpoint": "/api/data"})
	fmt.Printf("    ⚠️  Injected response time spike: 5000ms\n")

	// Trading volume anomaly
	detector.AddDataPoint("trading_volume", 10000.0, map[string]string{"symbol": "BTC/USD"})
	fmt.Printf("    ⚠️  Injected trading volume spike: 10000\n")

	// Wait for anomaly detection
	time.Sleep(1 * time.Second)

	// Check detected anomalies
	anomalies := detector.GetActiveAnomalies()
	fmt.Printf("    ✅ Detected %d active anomalies:\n", len(anomalies))

	for _, anomaly := range anomalies {
		fmt.Printf("      • %s: %.2f (expected: %.2f, deviation: %.2f, severity: %s)\n",
			anomaly.MetricName, anomaly.Value, anomaly.ExpectedValue, anomaly.Deviation, anomaly.Severity)
	}
}

// demoPredictiveAnalytics demonstrates predictive analytics capabilities
func demoPredictiveAnalytics(ctx context.Context, logger *observability.Logger) {
	fmt.Println("  Creating predictive analytics system...")

	config := &analytics.AnalyticsConfig{
		PredictionHorizon: 2 * time.Hour,
	}

	analyzer := analytics.NewPredictiveAnalyzer(logger, config)
	if err := analyzer.Start(ctx); err != nil {
		log.Printf("    ❌ Error starting predictive analyzer: %v", err)
		return
	}

	fmt.Printf("    ✅ Predictive analyzer started with %v horizon\n", config.PredictionHorizon)

	// Generate historical data
	fmt.Printf("    📊 Generating historical training data...\n")

	baseTime := time.Now().Add(-24 * time.Hour)
	for i := 0; i < 144; i++ { // 24 hours of 10-minute intervals
		timestamp := baseTime.Add(time.Duration(i) * 10 * time.Minute)

		// Simulate CPU usage with daily pattern
		hour := float64(timestamp.Hour())
		cpuBase := 30 + 20*math.Sin((hour-6)*math.Pi/12) // Peak at 6 PM
		cpuNoise := rand.Float64()*10 - 5
		cpuUsage := math.Max(0, math.Min(100, cpuBase+cpuNoise))

		dataPoint := analytics.DataPoint{
			Timestamp: timestamp,
			Value:     cpuUsage,
			Tags:      map[string]string{"host": "server1"},
		}

		analyzer.AddTrainingData("cpu_usage", dataPoint)

		// Simulate trading volume
		volumeBase := 1000 + 500*math.Sin((hour-10)*math.Pi/8) // Peak at 2 PM
		volumeNoise := rand.Float64()*200 - 100
		volume := math.Max(0, volumeBase+volumeNoise)

		volumePoint := analytics.DataPoint{
			Timestamp: timestamp,
			Value:     volume,
			Tags:      map[string]string{"symbol": "BTC/USD"},
		}

		analyzer.AddTrainingData("trading_volume", volumePoint)
	}

	fmt.Printf("    ✅ Generated 144 historical data points for each metric\n")

	// Create predictive models
	cpuModel, err := analyzer.CreateModel("cpu_usage", analytics.ModelTypeLinearRegression, map[string]float64{})
	if err != nil {
		fmt.Printf("    ❌ Error creating CPU model: %v\n", err)
		return
	}

	volumeModel, err := analyzer.CreateModel("trading_volume", analytics.ModelTypeExponentialSmoothing, map[string]float64{"alpha": 0.3})
	if err != nil {
		fmt.Printf("    ❌ Error creating volume model: %v\n", err)
		return
	}

	fmt.Printf("    ✅ Created predictive models: %s, %s\n", cpuModel.ModelType, volumeModel.ModelType)

	// Wait for model training
	time.Sleep(2 * time.Second)

	// Generate forecasts
	fmt.Printf("    🔮 Generating forecasts...\n")

	cpuForecast, err := analyzer.GenerateForecast(ctx, &analytics.ForecastRequest{
		MetricName: "cpu_usage",
		Horizon:    2 * time.Hour,
		Intervals:  12, // 10-minute intervals
	})
	if err != nil {
		fmt.Printf("    ❌ Error generating CPU forecast: %v\n", err)
		return
	}

	volumeForecast, err := analyzer.GenerateForecast(ctx, &analytics.ForecastRequest{
		MetricName: "trading_volume",
		Horizon:    2 * time.Hour,
		Intervals:  12,
	})
	if err != nil {
		fmt.Printf("    ❌ Error generating volume forecast: %v\n", err)
		return
	}

	fmt.Printf("    ✅ CPU Forecast (next 2 hours, confidence: %.1f%%):\n", cpuForecast.Confidence*100)
	for _, prediction := range cpuForecast.Predictions[:3] { // Show first 3 predictions
		fmt.Printf("      • %s: %.1f%% (trend: %s)\n",
			prediction.TargetTime.Format("15:04"), prediction.PredictedValue, prediction.Trend)
	}

	fmt.Printf("    ✅ Volume Forecast (next 2 hours, confidence: %.1f%%):\n", volumeForecast.Confidence*100)
	for _, prediction := range volumeForecast.Predictions[:3] {
		fmt.Printf("      • %s: %.0f BTC (trend: %s)\n",
			prediction.TargetTime.Format("15:04"), prediction.PredictedValue, prediction.Trend)
	}
}

// demoIntelligentAlerting demonstrates intelligent alerting capabilities
func demoIntelligentAlerting(ctx context.Context, logger *observability.Logger) {
	fmt.Println("  Creating intelligent alerting system...")

	config := &analytics.AnalyticsConfig{
		EnableIntelligentAlerting: true,
	}

	alertManager := analytics.NewAlertManager(logger, config)
	if err := alertManager.Start(ctx); err != nil {
		log.Printf("    ❌ Error starting alert manager: %v", err)
		return
	}

	fmt.Printf("    ✅ Alert manager started with intelligent alerting\n")

	// Create custom alert rules
	cpuRule := &analytics.AlertRule{
		Name:             "High CPU Usage Alert",
		Description:      "Triggers when CPU usage exceeds 80%",
		MetricName:       "cpu_usage",
		Condition:        analytics.ConditionGreaterThan,
		Threshold:        80.0,
		Severity:         analytics.SeverityWarning,
		Duration:         2 * time.Minute,
		EvaluationWindow: 30 * time.Second,
		Enabled:          true,
		Actions: []analytics.AlertAction{
			{
				ActionType: analytics.ActionTypeEmail,
				Target:     "admin@example.com",
				Enabled:    true,
			},
			{
				ActionType: analytics.ActionTypeSlack,
				Target:     "#alerts",
				Enabled:    true,
			},
		},
	}

	responseRule := &analytics.AlertRule{
		Name:             "High Response Time Alert",
		Description:      "Triggers when response time exceeds 1000ms",
		MetricName:       "response_time",
		Condition:        analytics.ConditionGreaterThan,
		Threshold:        1000.0,
		Severity:         analytics.SeverityCritical,
		Duration:         1 * time.Minute,
		EvaluationWindow: 15 * time.Second,
		Enabled:          true,
		Actions: []analytics.AlertAction{
			{
				ActionType: analytics.ActionTypeEmail,
				Target:     "oncall@example.com",
				Enabled:    true,
			},
			{
				ActionType: analytics.ActionTypePagerDuty,
				Target:     "incident-response",
				Enabled:    true,
			},
		},
	}

	// Create alert rules
	if err := alertManager.CreateAlertRule(cpuRule); err != nil {
		fmt.Printf("    ❌ Error creating CPU alert rule: %v\n", err)
		return
	}

	if err := alertManager.CreateAlertRule(responseRule); err != nil {
		fmt.Printf("    ❌ Error creating response time alert rule: %v\n", err)
		return
	}

	fmt.Printf("    ✅ Created %d custom alert rules\n", 2)

	// Simulate metric values that trigger alerts
	fmt.Printf("    🚨 Simulating alert-triggering conditions...\n")

	// Trigger CPU alert
	alertManager.EvaluateMetric("cpu_usage", 85.0, map[string]string{"host": "server1"})
	fmt.Printf("    ⚠️  CPU usage: 85%% (threshold: 80%%)\n")

	// Trigger response time alert
	alertManager.EvaluateMetric("response_time", 1500.0, map[string]string{"endpoint": "/api/data"})
	fmt.Printf("    ⚠️  Response time: 1500ms (threshold: 1000ms)\n")

	// Wait for alert processing
	time.Sleep(1 * time.Second)

	// Check active alerts
	activeAlerts := alertManager.GetActiveAlerts()
	fmt.Printf("    ✅ Active alerts: %d\n", len(activeAlerts))

	for _, alert := range activeAlerts {
		fmt.Printf("      • %s: %s (severity: %s, value: %.1f)\n",
			alert.RuleName, alert.Message, alert.Severity, alert.Value)
	}

	// Demonstrate alert acknowledgment
	if len(activeAlerts) > 0 {
		alertID := activeAlerts[0].AlertID
		if err := alertManager.AcknowledgeAlert(alertID, "admin"); err != nil {
			fmt.Printf("    ❌ Error acknowledging alert: %v\n", err)
		} else {
			fmt.Printf("    ✅ Alert acknowledged by admin\n")
		}
	}
}

// demoRealTimeDashboards demonstrates real-time dashboard capabilities
func demoRealTimeDashboards(ctx context.Context, logger *observability.Logger) {
	fmt.Println("  Creating real-time dashboard system...")

	config := &analytics.AnalyticsConfig{}
	dashboardManager := analytics.NewDashboardManager(logger, config)
	if err := dashboardManager.Start(ctx); err != nil {
		log.Printf("    ❌ Error starting dashboard manager: %v", err)
		return
	}

	fmt.Printf("    ✅ Dashboard manager started\n")

	// Get default dashboards
	dashboards := dashboardManager.GetDashboards()
	fmt.Printf("    ✅ Available dashboards: %d\n", len(dashboards))

	for _, dashboard := range dashboards {
		fmt.Printf("      • %s (%s) - %d widgets\n",
			dashboard.Name, dashboard.Category, len(dashboard.Widgets))
	}

	// Create custom trading dashboard
	tradingDashboard := &analytics.DashboardConfig{
		Name:        "Custom Trading Dashboard",
		Description: "Real-time trading analytics and market insights",
		Category:    analytics.CategoryTrading,
		Layout: &analytics.DashboardLayout{
			Type:       analytics.LayoutTypeGrid,
			Columns:    12,
			Rows:       8,
			GridSize:   analytics.GridSize{Width: 100, Height: 100},
			Responsive: true,
		},
		Theme:       "dark",
		RefreshRate: 10 * time.Second,
		AutoRefresh: true,
		Permissions: &analytics.DashboardPermissions{
			Owner:  "trader",
			Public: false,
		},
		Tags:      []string{"trading", "custom", "real-time"},
		CreatedBy: "demo",
		IsPublic:  false,
	}

	// Add widgets to trading dashboard
	tradingDashboard.Widgets = []*analytics.DashboardWidget{
		{
			Name:     "BTC Price",
			Type:     analytics.WidgetTypeMetric,
			Position: analytics.WidgetPosition{X: 0, Y: 0},
			Size:     analytics.WidgetSize{Width: 3, Height: 2},
			Configuration: &analytics.WidgetConfiguration{
				Title:     "BTC Price",
				ShowTitle: true,
				Unit:      "USD",
				Decimals:  2,
				Format:    "currency",
			},
			DataSource: &analytics.WidgetDataSource{
				Type:       analytics.DataSourceTypeMetrics,
				MetricName: "btc_price",
				TimeRange: analytics.TimeRange{
					Start: time.Now().Add(-time.Minute),
					End:   time.Now(),
				},
			},
			RefreshRate: 5 * time.Second,
			IsVisible:   true,
		},
	}

	// Create the dashboard
	if err := dashboardManager.CreateDashboard(tradingDashboard); err != nil {
		fmt.Printf("    ❌ Error creating trading dashboard: %v\n", err)
		return
	}

	fmt.Printf("    ✅ Created custom trading dashboard with %d widgets\n", len(tradingDashboard.Widgets))

	// Export dashboard configuration
	exportData, err := dashboardManager.ExportDashboard(tradingDashboard.DashboardID)
	if err != nil {
		fmt.Printf("    ❌ Error exporting dashboard: %v\n", err)
	} else {
		fmt.Printf("    ✅ Dashboard exported (%d bytes)\n", len(exportData))
	}

	// Get dashboard metrics
	metrics := dashboardManager.GetDashboardMetrics()
	fmt.Printf("    📊 Dashboard Metrics:\n")
	fmt.Printf("      • Total Dashboards: %v\n", metrics["total_dashboards"])
	fmt.Printf("      • Total Views: %v\n", metrics["total_views"])
	fmt.Printf("      • Total Themes: %v\n", metrics["total_themes"])
}

// demoIntegratedWorkflow demonstrates integrated analytics workflow
func demoIntegratedWorkflow(ctx context.Context, logger *observability.Logger) {
	fmt.Println("  Demonstrating integrated analytics workflow...")

	// Create integrated analytics system
	config := &analytics.AnalyticsConfig{
		EnableRealTimeProcessing:    true,
		EnablePredictiveAnalytics:   true,
		EnableAnomalyDetection:      true,
		EnableIntelligentAlerting:   true,
		ProcessingInterval:          100 * time.Millisecond,
		AnomalyDetectionSensitivity: 0.8,
		PredictionHorizon:           1 * time.Hour,
		MaxConcurrentStreams:        50,
		BufferSize:                  5000,
	}

	// Initialize all components
	engine := analytics.NewRealTimeAnalyticsEngine(logger, config)
	detector := analytics.NewAnomalyDetector(logger, config)
	analyzer := analytics.NewPredictiveAnalyzer(logger, config)
	alertManager := analytics.NewAlertManager(logger, config)
	dashboardManager := analytics.NewDashboardManager(logger, config)

	// Start all components
	components := []struct {
		name    string
		starter interface{ Start(context.Context) error }
	}{
		{"Analytics Engine", engine},
		{"Anomaly Detector", detector},
		{"Predictive Analyzer", analyzer},
		{"Alert Manager", alertManager},
		{"Dashboard Manager", dashboardManager},
	}

	for _, component := range components {
		if err := component.starter.Start(ctx); err != nil {
			fmt.Printf("    ❌ Error starting %s: %v\n", component.name, err)
			return
		}
	}

	fmt.Printf("    ✅ Started %d integrated components\n", len(components))

	// Create data stream
	stream, err := engine.CreateDataStream(
		"Integrated Workflow",
		"demo_system",
		[]analytics.EventType{
			analytics.EventTypeSystemMetric,
			analytics.EventTypeTradingActivity,
			analytics.EventTypePerformance,
		},
		nil,
	)
	if err != nil {
		fmt.Printf("    ❌ Error creating integrated stream: %v\n", err)
		return
	}

	fmt.Printf("    ✅ Created integrated data stream: %s\n", stream.Name)

	// Register anomaly detectors
	detector.RegisterMetricDetector("integrated_metric", analytics.DetectionMethodStatistical, 0.8, 30)

	// Create alert rule
	alertRule := &analytics.AlertRule{
		Name:        "Integrated Alert",
		Description: "Integrated workflow alert",
		MetricName:  "integrated_metric",
		Condition:   analytics.ConditionGreaterThan,
		Threshold:   100.0,
		Severity:    analytics.SeverityWarning,
		Enabled:     true,
		Actions: []analytics.AlertAction{
			{
				ActionType: analytics.ActionTypeEmail,
				Target:     "integrated@example.com",
				Enabled:    true,
			},
		},
	}
	alertManager.CreateAlertRule(alertRule)

	// Simulate integrated workflow
	fmt.Printf("    🔄 Running integrated analytics workflow...\n")

	for i := 0; i < 20; i++ {
		// Generate metric value
		value := 50 + rand.Float64()*100
		if i > 15 { // Introduce anomaly
			value = 150 + rand.Float64()*50
		}

		// Publish to analytics engine
		event := &analytics.AnalyticsEvent{
			EventType: analytics.EventTypeSystemMetric,
			Source:    "integrated_demo",
			Timestamp: time.Now(),
			Data:      map[string]interface{}{"metric": "integrated_metric"},
			Metrics:   map[string]float64{"integrated_metric": value},
			Tags:      []string{"demo", "integrated"},
			Priority:  analytics.EventPriorityMedium,
		}
		engine.PublishEvent(event)

		// Add to anomaly detector
		detector.AddDataPoint("integrated_metric", value, map[string]string{"source": "demo"})

		// Add to predictive analyzer
		dataPoint := analytics.DataPoint{
			Timestamp: time.Now(),
			Value:     value,
			Tags:      map[string]string{"source": "demo"},
		}
		analyzer.AddTrainingData("integrated_metric", dataPoint)

		// Evaluate alerts
		alertManager.EvaluateMetric("integrated_metric", value, map[string]string{"source": "demo"})

		time.Sleep(100 * time.Millisecond)
	}

	// Wait for processing
	time.Sleep(2 * time.Second)

	// Show results
	anomalies := detector.GetActiveAnomalies()
	alerts := alertManager.GetActiveAlerts()
	dashboards := dashboardManager.GetDashboards()

	fmt.Printf("    ✅ Integrated Workflow Results:\n")
	fmt.Printf("      • Active Anomalies: %d\n", len(anomalies))
	fmt.Printf("      • Active Alerts: %d\n", len(alerts))
	fmt.Printf("      • Available Dashboards: %d\n", len(dashboards))

	if len(anomalies) > 0 {
		fmt.Printf("      • Latest Anomaly: %s (severity: %s)\n",
			anomalies[0].MetricName, anomalies[0].Severity)
	}

	if len(alerts) > 0 {
		fmt.Printf("      • Latest Alert: %s (severity: %s)\n",
			alerts[0].RuleName, alerts[0].Severity)
	}

	fmt.Printf("    ✅ Integrated analytics workflow completed successfully\n")
}

// simulateEvents simulates real-time events for the analytics engine
func simulateEvents(engine *analytics.RealTimeAnalyticsEngine) {
	for i := 0; i < 50; i++ {
		// Trading event
		tradingEvent := &analytics.AnalyticsEvent{
			EventType: analytics.EventTypeTradingActivity,
			Source:    "trading_simulator",
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"symbol":   "BTC/USD",
				"side":     "buy",
				"quantity": rand.Float64() * 10,
			},
			Metrics: map[string]float64{
				"price":  43000 + rand.Float64()*1000,
				"volume": 100 + rand.Float64()*900,
			},
			Tags:     []string{"trading", "btc"},
			Priority: analytics.EventPriorityMedium,
		}

		// System event
		systemEvent := &analytics.AnalyticsEvent{
			EventType: analytics.EventTypeSystemMetric,
			Source:    "system_monitor",
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"host":   "server1",
				"metric": "cpu_usage",
			},
			Metrics: map[string]float64{
				"cpu_usage":    30 + rand.Float64()*40,
				"memory_usage": 40 + rand.Float64()*30,
				"disk_usage":   20 + rand.Float64()*20,
			},
			Tags:     []string{"system", "monitoring"},
			Priority: analytics.EventPriorityLow,
		}

		engine.PublishEvent(tradingEvent)
		engine.PublishEvent(systemEvent)

		time.Sleep(200 * time.Millisecond)
	}
}
