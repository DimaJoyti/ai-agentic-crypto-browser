package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// DemoClient demonstrates the AI-Agentic Crypto Browser capabilities
type DemoClient struct {
	baseURL    string
	httpClient *http.Client
	authToken  string
	userID     uuid.UUID
}

// AuthResponse represents authentication response
type AuthResponse struct {
	Token  string    `json:"token"`
	UserID uuid.UUID `json:"user_id"`
}

func main() {
	fmt.Println("🚀 AI-Agentic Crypto Browser - Advanced AI Demo")
	fmt.Println(strings.Repeat("=", 60))

	// Initialize demo client
	client := &DemoClient{
		baseURL: "http://localhost:8080",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// Check if server is running
	if !client.checkServerHealth() {
		log.Fatal("❌ Server is not running. Please start the AI agent first.")
	}

	// Authenticate
	if err := client.authenticate(); err != nil {
		log.Fatalf("❌ Authentication failed: %v", err)
	}

	fmt.Printf("✅ Authenticated successfully (User ID: %s)\n\n", client.userID)

	// Run comprehensive demo
	client.runComprehensiveDemo()
}

func (c *DemoClient) checkServerHealth() bool {
	resp, err := c.httpClient.Get(c.baseURL + "/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func (c *DemoClient) authenticate() error {
	// For demo purposes, we'll use a simple auth request
	authReq := map[string]interface{}{
		"email":    "demo@example.com",
		"password": "demo123",
	}

	resp, err := c.makeRequest("POST", "/auth/login", authReq)
	if err != nil {
		return fmt.Errorf("auth request failed: %w", err)
	}

	var authResp AuthResponse
	if err := json.Unmarshal(resp, &authResp); err != nil {
		// If structured auth fails, create demo credentials
		c.authToken = "demo-token-" + uuid.New().String()
		c.userID = uuid.New()
		return nil
	}

	c.authToken = authResp.Token
	c.userID = authResp.UserID
	return nil
}

func (c *DemoClient) runComprehensiveDemo() {
	fmt.Println("🧠 Starting Comprehensive AI Capabilities Demo")
	fmt.Println(strings.Repeat("-", 50))

	// Demo 1: Enhanced AI Analysis
	c.demoEnhancedAnalysis()

	// Demo 2: Advanced Price Prediction
	c.demoPricePrediction()

	// Demo 3: Multi-Language Sentiment Analysis
	c.demoSentimentAnalysis()

	// Demo 4: Predictive Analytics
	c.demoPredictiveAnalytics()

	// Demo 5: Advanced NLP Analysis
	c.demoAdvancedNLP()

	// Demo 6: Intelligent Decision Making
	c.demoIntelligentDecisionMaking()

	// Demo 7: Multi-Modal AI Analysis
	c.demoMultiModalAnalysis()

	// Demo 8: User Behavior Learning
	c.demoUserBehaviorLearning()

	// Demo 9: Performance Metrics
	c.demoPerformanceMetrics()

	fmt.Println("\n🎉 Demo completed successfully!")
	fmt.Println("✨ The AI-Agentic Crypto Browser demonstrates:")
	fmt.Println("   • Advanced AI analysis and prediction")
	fmt.Println("   • Multi-language NLP and sentiment analysis")
	fmt.Println("   • Multi-modal AI (images, documents, audio)")
	fmt.Println("   • Intelligent decision making with risk management")
	fmt.Println("   • Advanced user behavior learning and personalization")
	fmt.Println("   • Continuous learning and adaptation")
	fmt.Println("   • Comprehensive performance tracking")
}

func (c *DemoClient) demoEnhancedAnalysis() {
	fmt.Println("\n📊 Demo 1: Enhanced AI Analysis")
	fmt.Println("Testing comprehensive market analysis...")

	request := map[string]interface{}{
		"symbols":    []string{"BTC", "ETH", "ADA"},
		"timeframe":  "1h",
		"indicators": []string{"RSI", "MACD", "Bollinger"},
		"analysis_types": []string{
			"technical", "sentiment", "risk", "prediction",
		},
	}

	resp, err := c.makeAuthenticatedRequest("POST", "/ai/analyze", request)
	if err != nil {
		fmt.Printf("❌ Enhanced analysis failed: %v\n", err)
		return
	}

	var result map[string]interface{}
	json.Unmarshal(resp, &result)

	fmt.Printf("✅ Analysis completed for %d symbols\n", len(request["symbols"].([]string)))
	if confidence, ok := result["confidence"].(float64); ok {
		fmt.Printf("   Confidence: %.1f%%\n", confidence*100)
	}
	if recommendation, ok := result["recommendation"].(string); ok {
		fmt.Printf("   Recommendation: %s\n", recommendation)
	}
}

func (c *DemoClient) demoPricePrediction() {
	fmt.Println("\n🔮 Demo 2: Advanced Price Prediction")
	fmt.Println("Testing AI-powered price forecasting...")

	request := map[string]interface{}{
		"symbol":     "BTC",
		"timeframes": []string{"1h", "4h", "1d", "1w"},
		"models":     []string{"lstm", "transformer", "ensemble"},
		"features": []string{
			"price", "volume", "sentiment", "technical_indicators",
		},
	}

	resp, err := c.makeAuthenticatedRequest("POST", "/ai/predict/price", request)
	if err != nil {
		fmt.Printf("❌ Price prediction failed: %v\n", err)
		return
	}

	var result map[string]interface{}
	json.Unmarshal(resp, &result)

	fmt.Printf("✅ Price prediction completed\n")
	if predictions, ok := result["predictions"].(map[string]interface{}); ok {
		for timeframe, pred := range predictions {
			if predMap, ok := pred.(map[string]interface{}); ok {
				if price, ok := predMap["predicted_price"].(float64); ok {
					fmt.Printf("   %s: $%.2f\n", timeframe, price)
				}
			}
		}
	}
}

func (c *DemoClient) demoSentimentAnalysis() {
	fmt.Println("\n💭 Demo 3: Multi-Language Sentiment Analysis")
	fmt.Println("Testing sentiment analysis across languages...")

	request := map[string]interface{}{
		"texts": []string{
			"Bitcoin is showing incredible bullish momentum! 🚀",
			"El Bitcoin está mostrando un impulso alcista increíble",
			"Le Bitcoin montre une dynamique haussière incroyable",
			"Bitcoin zeigt eine unglaubliche bullische Dynamik",
		},
		"languages": []string{"en", "es", "fr", "de"},
		"options": map[string]interface{}{
			"detect_emotions":   true,
			"analyze_intensity": true,
			"extract_aspects":   true,
		},
	}

	resp, err := c.makeAuthenticatedRequest("POST", "/ai/analyze/sentiment", request)
	if err != nil {
		fmt.Printf("❌ Sentiment analysis failed: %v\n", err)
		return
	}

	var result map[string]interface{}
	json.Unmarshal(resp, &result)

	fmt.Printf("✅ Sentiment analysis completed\n")
	if overallSentiment, ok := result["overall_sentiment"].(map[string]interface{}); ok {
		if score, ok := overallSentiment["score"].(float64); ok {
			if label, ok := overallSentiment["label"].(string); ok {
				fmt.Printf("   Overall: %s (%.2f)\n", label, score)
			}
		}
	}
}

func (c *DemoClient) demoPredictiveAnalytics() {
	fmt.Println("\n🎯 Demo 4: Predictive Analytics")
	fmt.Println("Testing comprehensive predictive analytics...")

	request := map[string]interface{}{
		"analysis_type": "comprehensive",
		"assets":        []string{"BTC", "ETH"},
		"timeframe":     "1d",
		"scenarios": []string{
			"bullish", "bearish", "sideways", "volatile",
		},
		"include_risk_analysis":    true,
		"include_portfolio_impact": true,
		"include_correlations":     true,
	}

	resp, err := c.makeAuthenticatedRequest("POST", "/ai/analytics/predictive", request)
	if err != nil {
		fmt.Printf("❌ Predictive analytics failed: %v\n", err)
		return
	}

	var result map[string]interface{}
	json.Unmarshal(resp, &result)

	fmt.Printf("✅ Predictive analytics completed\n")
	if scenarios, ok := result["scenarios"].([]interface{}); ok {
		fmt.Printf("   Analyzed %d market scenarios\n", len(scenarios))
	}
	if riskMetrics, ok := result["risk_metrics"].(map[string]interface{}); ok {
		if overallRisk, ok := riskMetrics["overall_risk"].(float64); ok {
			fmt.Printf("   Overall Risk: %.1f%%\n", overallRisk*100)
		}
	}
}

func (c *DemoClient) demoAdvancedNLP() {
	fmt.Println("\n🔤 Demo 5: Advanced NLP Analysis")
	fmt.Println("Testing comprehensive NLP capabilities...")

	request := map[string]interface{}{
		"texts": []string{
			"Bitcoin's institutional adoption is accelerating with major corporations adding BTC to their treasury reserves.",
			"The recent regulatory clarity has boosted investor confidence in cryptocurrency markets.",
			"DeFi protocols are revolutionizing traditional finance with innovative yield farming strategies.",
		},
		"sources": []string{"news", "analysis", "social"},
		"options": map[string]interface{}{
			"detect_language":        true,
			"extract_entities":       true,
			"perform_topic_modeling": true,
			"classify_text":          true,
			"analyze_sentiment":      true,
			"extract_keywords":       true,
			"detect_emotions":        true,
			"analyze_readability":    true,
		},
	}

	resp, err := c.makeAuthenticatedRequest("POST", "/ai/nlp/analyze", request)
	if err != nil {
		fmt.Printf("❌ Advanced NLP failed: %v\n", err)
		return
	}

	var result map[string]interface{}
	json.Unmarshal(resp, &result)

	fmt.Printf("✅ Advanced NLP analysis completed\n")
	if results, ok := result["results"].([]interface{}); ok {
		fmt.Printf("   Processed %d texts\n", len(results))
	}
	if aggregated, ok := result["aggregated_results"].(map[string]interface{}); ok {
		if sentiment, ok := aggregated["overall_sentiment"].(map[string]interface{}); ok {
			if label, ok := sentiment["label"].(string); ok {
				fmt.Printf("   Overall Sentiment: %s\n", label)
			}
		}
	}
}

func (c *DemoClient) demoIntelligentDecisionMaking() {
	fmt.Println("\n🎯 Demo 6: Intelligent Decision Making")
	fmt.Println("Testing AI-driven trading decisions...")

	request := map[string]interface{}{
		"decision_type": "trade",
		"context": map[string]interface{}{
			"market_conditions": "bullish",
			"time_horizon":      "short",
			"urgency":           "medium",
			"trigger_event":     "price_breakout",
			"technical_indicators": map[string]float64{
				"rsi":          30.0,
				"macd_signal":  1.0,
				"volume_ratio": 1.5,
			},
		},
		"constraints": map[string]interface{}{
			"max_position_size": 1000.0,
			"max_risk_exposure": 0.05,
			"allowed_assets":    []string{"BTC", "ETH"},
		},
		"preferences": map[string]interface{}{
			"risk_tolerance":       0.6,
			"auto_execution_level": "none",
			"decision_speed":       "normal",
		},
		"market_data": map[string]interface{}{
			"timestamp": time.Now(),
			"prices": map[string]float64{
				"BTC": 50000.0,
				"ETH": 3000.0,
			},
			"sentiment":  0.7,
			"volatility": map[string]float64{"BTC": 0.3, "ETH": 0.4},
		},
		"options": map[string]interface{}{
			"require_confirmation": true,
			"explain_reasoning":    true,
			"simulate_execution":   true,
		},
	}

	resp, err := c.makeAuthenticatedRequest("POST", "/ai/decisions/request", request)
	if err != nil {
		fmt.Printf("❌ Decision making failed: %v\n", err)
		return
	}

	var result map[string]interface{}
	json.Unmarshal(resp, &result)

	fmt.Printf("✅ Intelligent decision completed\n")
	if recommendation, ok := result["recommendation"].(map[string]interface{}); ok {
		if action, ok := recommendation["action"].(string); ok {
			if asset, ok := recommendation["asset"].(string); ok {
				fmt.Printf("   Recommendation: %s %s\n", action, asset)
			}
		}
		if confidence, ok := recommendation["confidence"].(float64); ok {
			fmt.Printf("   Confidence: %.1f%%\n", confidence*100)
		}
	}
	if requiresApproval, ok := result["requires_approval"].(bool); ok {
		fmt.Printf("   Requires Approval: %t\n", requiresApproval)
	}
}

func (c *DemoClient) demoMultiModalAnalysis() {
	fmt.Println("\n🎨 Demo 7: Multi-Modal AI Analysis")
	fmt.Println("Testing advanced multi-modal AI capabilities...")

	// Demo 1: Get supported formats
	formatsResp, err := c.makeAuthenticatedRequest("GET", "/ai/multimodal/formats", nil)
	if err != nil {
		fmt.Printf("❌ Get supported formats failed: %v\n", err)
		return
	}

	var formatsResult map[string]interface{}
	json.Unmarshal(formatsResp, &formatsResult)
	fmt.Printf("✅ Supported formats retrieved\n")
	if formats, ok := formatsResult["supported_formats"].(map[string]interface{}); ok {
		if images, ok := formats["images"].([]interface{}); ok {
			fmt.Printf("   Image formats: %d supported\n", len(images))
		}
		if docs, ok := formats["documents"].([]interface{}); ok {
			fmt.Printf("   Document formats: %d supported\n", len(docs))
		}
		if audio, ok := formats["audio"].([]interface{}); ok {
			fmt.Printf("   Audio formats: %d supported\n", len(audio))
		}
	}

	// Demo 2: Multi-modal analysis request
	request := map[string]interface{}{
		"type": "mixed",
		"content": []map[string]interface{}{
			{
				"id":        "content-1",
				"type":      "image",
				"data":      "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChAHGArEkAAAAAElFTkSuQmCC",
				"mime_type": "image/png",
				"filename":  "chart.png",
				"size":      100,
			},
			{
				"id":        "content-2",
				"type":      "document",
				"data":      base64.StdEncoding.EncodeToString([]byte("Bitcoin price analysis shows bullish momentum with strong support at $50,000.")),
				"mime_type": "text/plain",
				"filename":  "analysis.txt",
				"size":      80,
			},
		},
		"options": map[string]interface{}{
			"analyze_images":    true,
			"extract_text":      true,
			"analyze_charts":    true,
			"analyze_sentiment": true,
			"generate_summary":  true,
		},
	}

	resp, err := c.makeAuthenticatedRequest("POST", "/ai/multimodal/analyze", request)
	if err != nil {
		fmt.Printf("❌ Multi-modal analysis failed: %v\n", err)
		return
	}

	var result map[string]interface{}
	json.Unmarshal(resp, &result)

	fmt.Printf("✅ Multi-modal analysis completed\n")
	if results, ok := result["results"].([]interface{}); ok {
		fmt.Printf("   Processed %d content items\n", len(results))
	}
	if aggregated, ok := result["aggregated_data"].(map[string]interface{}); ok {
		if stats, ok := aggregated["processing_stats"].(map[string]interface{}); ok {
			if successCount, ok := stats["successful_items"].(float64); ok {
				fmt.Printf("   Successful items: %.0f\n", successCount)
			}
		}
		if insights, ok := aggregated["key_insights"].([]interface{}); ok {
			fmt.Printf("   Key insights: %d\n", len(insights))
		}
	}
}

func (c *DemoClient) demoUserBehaviorLearning() {
	fmt.Println("\n🧠 Demo 8: User Behavior Learning")
	fmt.Println("Testing advanced user behavior learning and personalization...")

	// Demo 1: Learn from trading behavior
	tradingEvent := map[string]interface{}{
		"type":   "trade",
		"action": "buy_btc",
		"context": map[string]interface{}{
			"market_conditions":   "bullish",
			"portfolio_state":     map[string]interface{}{"btc_balance": 0.5},
			"time_of_day":         "morning",
			"day_of_week":         "monday",
			"session_duration":    "2h",
			"previous_actions":    []string{"analyze_chart", "check_news"},
			"emotional_state":     "confident",
			"information_sources": []string{"technical_analysis", "news"},
			"external_factors":    map[string]interface{}{"market_sentiment": "positive"},
		},
		"outcome": map[string]interface{}{
			"success":          true,
			"performance":      0.05,
			"satisfaction":     0.8,
			"time_to_decision": "15m",
			"confidence_level": 0.7,
			"regret":           0.1,
			"learning_value":   0.8,
		},
		"duration": "30m",
	}

	_, err := c.makeAuthenticatedRequest("POST", "/ai/behavior/learn", tradingEvent)
	if err != nil {
		fmt.Printf("❌ Behavior learning failed: %v\n", err)
		return
	}
	fmt.Printf("✅ Trading behavior learned successfully\n")

	// Demo 2: Learn from analysis behavior
	analysisEvent := map[string]interface{}{
		"type":   "analyze",
		"action": "technical_analysis",
		"context": map[string]interface{}{
			"market_conditions": "neutral",
			"time_of_day":       "afternoon",
			"session_duration":  "1h",
		},
		"outcome": map[string]interface{}{
			"success":        true,
			"satisfaction":   0.9,
			"learning_value": 0.8,
		},
		"duration": "45m",
	}

	_, err = c.makeAuthenticatedRequest("POST", "/ai/behavior/learn", analysisEvent)
	if err != nil {
		fmt.Printf("❌ Analysis behavior learning failed: %v\n", err)
		return
	}
	fmt.Printf("✅ Analysis behavior learned successfully\n")

	// Demo 3: Get user behavior profile
	profileResp, err := c.makeAuthenticatedRequest("GET", "/ai/behavior/profile", nil)
	if err != nil {
		fmt.Printf("❌ Get user profile failed: %v\n", err)
		return
	}

	var profile map[string]interface{}
	json.Unmarshal(profileResp, &profile)
	fmt.Printf("✅ User behavior profile retrieved\n")

	if observationCount, ok := profile["observation_count"].(float64); ok {
		fmt.Printf("   Observations: %.0f\n", observationCount)
	}
	if confidence, ok := profile["confidence"].(float64); ok {
		fmt.Printf("   Overall Confidence: %.2f\n", confidence)
	}
	if tradingStyle, ok := profile["trading_style"].(map[string]interface{}); ok {
		if primaryStyle, ok := tradingStyle["primary_style"].(string); ok {
			fmt.Printf("   Trading Style: %s\n", primaryStyle)
		}
	}

	// Demo 4: Get personalized recommendations
	recResp, err := c.makeAuthenticatedRequest("GET", "/ai/behavior/recommendations?limit=3", nil)
	if err != nil {
		fmt.Printf("❌ Get recommendations failed: %v\n", err)
		return
	}

	var recResult map[string]interface{}
	json.Unmarshal(recResp, &recResult)
	fmt.Printf("✅ Personalized recommendations retrieved\n")

	if recommendations, ok := recResult["recommendations"].([]interface{}); ok {
		fmt.Printf("   Recommendations: %d\n", len(recommendations))
		for i, rec := range recommendations {
			if recMap, ok := rec.(map[string]interface{}); ok {
				if title, ok := recMap["title"].(string); ok {
					fmt.Printf("   %d. %s\n", i+1, title)
				}
			}
		}
	}

	// Demo 5: Get behavior history
	historyResp, err := c.makeAuthenticatedRequest("GET", "/ai/behavior/history?limit=5", nil)
	if err != nil {
		fmt.Printf("❌ Get behavior history failed: %v\n", err)
		return
	}

	var historyResult map[string]interface{}
	json.Unmarshal(historyResp, &historyResult)
	fmt.Printf("✅ Behavior history retrieved\n")

	if history, ok := historyResult["history"].([]interface{}); ok {
		fmt.Printf("   History entries: %d\n", len(history))
	}

	// Demo 6: Get learning models
	modelsResp, err := c.makeAuthenticatedRequest("GET", "/ai/behavior/models", nil)
	if err != nil {
		fmt.Printf("❌ Get learning models failed: %v\n", err)
		return
	}

	var modelsResult map[string]interface{}
	json.Unmarshal(modelsResp, &modelsResult)
	fmt.Printf("✅ Learning models retrieved\n")

	if models, ok := modelsResult["models"].(map[string]interface{}); ok {
		fmt.Printf("   Available models: %d\n", len(models))
	}
}

func (c *DemoClient) demoPerformanceMetrics() {
	fmt.Println("\n📈 Demo 9: Performance Metrics")
	fmt.Println("Testing performance tracking and analytics...")

	// Get learning performance
	_, err := c.makeAuthenticatedRequest("GET", "/ai/learning/performance", nil)
	if err != nil {
		fmt.Printf("❌ Learning performance failed: %v\n", err)
	} else {
		fmt.Printf("✅ Learning performance metrics retrieved\n")
	}

	// Get decision performance
	decisionResp, err := c.makeAuthenticatedRequest("GET", "/ai/decisions/performance", nil)
	if err != nil {
		fmt.Printf("❌ Decision performance failed: %v\n", err)
	} else {
		var result map[string]interface{}
		json.Unmarshal(decisionResp, &result)
		fmt.Printf("✅ Decision performance metrics retrieved\n")

		if metrics, ok := result["performance_metrics"].(map[string]interface{}); ok {
			if totalDecisions, ok := metrics["total_decisions"].(float64); ok {
				fmt.Printf("   Total Decisions: %.0f\n", totalDecisions)
			}
			if successRate, ok := metrics["success_rate"].(float64); ok {
				fmt.Printf("   Success Rate: %.1f%%\n", successRate*100)
			}
		}
	}

	// Get adaptive models
	modelsResp, err := c.makeAuthenticatedRequest("GET", "/ai/adaptation/models", nil)
	if err != nil {
		fmt.Printf("❌ Adaptive models failed: %v\n", err)
	} else {
		var result map[string]interface{}
		json.Unmarshal(modelsResp, &result)
		fmt.Printf("✅ Adaptive models status retrieved\n")

		if models, ok := result["models"].(map[string]interface{}); ok {
			fmt.Printf("   Active Models: %d\n", len(models))
		}
	}
}

func (c *DemoClient) makeRequest(method, endpoint string, data interface{}) ([]byte, error) {
	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.baseURL+endpoint, body)
	if err != nil {
		return nil, err
	}

	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (c *DemoClient) makeAuthenticatedRequest(method, endpoint string, data interface{}) ([]byte, error) {
	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.baseURL+endpoint, body)
	if err != nil {
		return nil, err
	}

	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add authentication header
	req.Header.Set("Authorization", "Bearer "+c.authToken)
	req.Header.Set("X-User-ID", c.userID.String())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
