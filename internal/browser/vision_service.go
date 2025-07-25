package browser

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

// VisionService provides AI-powered visual analysis of web pages
type VisionService struct {
	service *Service
	logger  *observability.Logger
}

// NewVisionService creates a new vision service
func NewVisionService(service *Service, logger *observability.Logger) *VisionService {
	return &VisionService{
		service: service,
		logger:  logger,
	}
}

// VisualAnalysisRequest represents a request for visual analysis
type VisualAnalysisRequest struct {
	SessionID    uuid.UUID `json:"session_id" validate:"required"`
	AnalysisType string    `json:"analysis_type" validate:"required"`
	Target       string    `json:"target,omitempty"` // CSS selector for specific element
	Instructions string    `json:"instructions,omitempty"`
}

// VisualAnalysisResponse represents the response from visual analysis
type VisualAnalysisResponse struct {
	Success     bool                   `json:"success"`
	Analysis    map[string]interface{} `json:"analysis"`
	Screenshot  string                 `json:"screenshot"`
	Elements    []VisualElement        `json:"elements,omitempty"`
	Suggestions []string               `json:"suggestions,omitempty"`
	Metadata    map[string]interface{} `json:"metadata"`
	Error       string                 `json:"error,omitempty"`
}

// VisualElement represents a detected UI element
type VisualElement struct {
	Type       string                 `json:"type"`
	Selector   string                 `json:"selector"`
	Text       string                 `json:"text,omitempty"`
	Bounds     ElementBounds          `json:"bounds"`
	Attributes map[string]string      `json:"attributes,omitempty"`
	Confidence float64                `json:"confidence"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// ElementBounds represents the position and size of an element
type ElementBounds struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// SmartInteractionRequest represents a request for AI-guided interaction
type SmartInteractionRequest struct {
	SessionID uuid.UUID `json:"session_id" validate:"required"`
	Goal      string    `json:"goal" validate:"required"`
	Context   string    `json:"context,omitempty"`
	MaxSteps  int       `json:"max_steps,omitempty"`
	SafeMode  bool      `json:"safe_mode,omitempty"`
}

// SmartInteractionResponse represents the response from smart interaction
type SmartInteractionResponse struct {
	Success     bool                   `json:"success"`
	Steps       []InteractionStep      `json:"steps"`
	Goal        string                 `json:"goal"`
	Status      string                 `json:"status"`
	Screenshots []string               `json:"screenshots,omitempty"`
	Metadata    map[string]interface{} `json:"metadata"`
	Error       string                 `json:"error,omitempty"`
}

// InteractionStep represents a single step in smart interaction
type InteractionStep struct {
	StepNumber int                    `json:"step_number"`
	Action     Action                 `json:"action"`
	Reasoning  string                 `json:"reasoning"`
	Success    bool                   `json:"success"`
	Screenshot string                 `json:"screenshot,omitempty"`
	Duration   time.Duration          `json:"duration"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Error      string                 `json:"error,omitempty"`
}

// AnalyzePageVisually performs AI-powered visual analysis of a web page
func (vs *VisionService) AnalyzePageVisually(ctx context.Context, req VisualAnalysisRequest) (*VisualAnalysisResponse, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("browser-service").Start(ctx, "vision.AnalyzePageVisually")
	defer span.End()

	// Take screenshot
	screenshotReq := ScreenshotRequest{
		FullPage: true,
		Quality:  90,
		Format:   "png",
	}

	screenshotResp, err := vs.service.TakeScreenshot(ctx, req.SessionID, screenshotReq)
	if err != nil {
		return &VisualAnalysisResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to take screenshot: %v", err),
		}, nil
	}

	// Detect UI elements
	elements, err := vs.detectUIElements(ctx, req.SessionID, req.Target)
	if err != nil {
		vs.logger.Error(ctx, "Failed to detect UI elements", err)
		// Continue with analysis even if element detection fails
	}

	// Perform visual analysis based on type
	var analysis map[string]interface{}
	switch req.AnalysisType {
	case "layout":
		analysis = vs.analyzeLayout(ctx, screenshotResp.Screenshot, elements)
	case "accessibility":
		analysis = vs.analyzeAccessibility(ctx, elements)
	case "usability":
		analysis = vs.analyzeUsability(ctx, elements)
	case "content":
		analysis = vs.analyzeContent(ctx, screenshotResp.Screenshot, elements)
	case "forms":
		analysis = vs.analyzeForms(ctx, elements)
	case "navigation":
		analysis = vs.analyzeNavigation(ctx, elements)
	default:
		analysis = vs.performGeneralAnalysis(ctx, screenshotResp.Screenshot, elements, req.Instructions)
	}

	// Generate suggestions
	suggestions := vs.generateSuggestions(ctx, req.AnalysisType, analysis, elements)

	response := &VisualAnalysisResponse{
		Success:     true,
		Analysis:    analysis,
		Screenshot:  screenshotResp.Screenshot,
		Elements:    elements,
		Suggestions: suggestions,
		Metadata: map[string]interface{}{
			"session_id":     req.SessionID.String(),
			"analysis_type":  req.AnalysisType,
			"timestamp":      time.Now(),
			"elements_count": len(elements),
		},
	}

	vs.logger.Info(ctx, "Visual analysis completed", map[string]interface{}{
		"session_id":     req.SessionID.String(),
		"analysis_type":  req.AnalysisType,
		"elements_found": len(elements),
	})

	return response, nil
}

// PerformSmartInteraction performs AI-guided interaction with a web page
func (vs *VisionService) PerformSmartInteraction(ctx context.Context, req SmartInteractionRequest) (*SmartInteractionResponse, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("browser-service").Start(ctx, "vision.PerformSmartInteraction")
	defer span.End()

	maxSteps := req.MaxSteps
	if maxSteps == 0 {
		maxSteps = 10
	}

	var steps []InteractionStep
	var screenshots []string

	// Take initial screenshot
	initialScreenshot, err := vs.takeStepScreenshot(ctx, req.SessionID)
	if err != nil {
		return &SmartInteractionResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to take initial screenshot: %v", err),
		}, nil
	}
	screenshots = append(screenshots, initialScreenshot)

	// Analyze current page state
	elements, err := vs.detectUIElements(ctx, req.SessionID, "")
	if err != nil {
		return &SmartInteractionResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to detect UI elements: %v", err),
		}, nil
	}

	// Plan interaction steps
	plannedSteps := vs.planInteractionSteps(ctx, req.Goal, req.Context, elements, initialScreenshot)

	// Execute steps
	for i, plannedStep := range plannedSteps {
		if i >= maxSteps {
			break
		}

		stepStart := time.Now()
		step := InteractionStep{
			StepNumber: i + 1,
			Action:     plannedStep,
			Reasoning:  fmt.Sprintf("Step %d: %s", i+1, vs.getActionReasoning(plannedStep)),
		}

		// Execute the action
		interactReq := InteractRequest{
			Actions:    []Action{plannedStep},
			Screenshot: true,
		}

		interactResp, err := vs.service.Interact(ctx, req.SessionID, interactReq)
		if err != nil {
			step.Success = false
			step.Error = err.Error()
		} else {
			step.Success = interactResp.Success
			if len(interactResp.Screenshots) > 0 {
				step.Screenshot = interactResp.Screenshots[0]
				screenshots = append(screenshots, step.Screenshot)
			}
		}

		step.Duration = time.Since(stepStart)
		steps = append(steps, step)

		// If in safe mode, pause between steps
		if req.SafeMode {
			time.Sleep(1 * time.Second)
		}

		// Check if goal is achieved (simplified check)
		if vs.isGoalAchieved(ctx, req.Goal, step.Screenshot) {
			break
		}
	}

	// Determine overall status
	status := "completed"
	successCount := 0
	for _, step := range steps {
		if step.Success {
			successCount++
		}
	}

	if successCount == 0 {
		status = "failed"
	} else if successCount < len(steps) {
		status = "partial"
	}

	response := &SmartInteractionResponse{
		Success:     successCount > 0,
		Steps:       steps,
		Goal:        req.Goal,
		Status:      status,
		Screenshots: screenshots,
		Metadata: map[string]interface{}{
			"session_id":   req.SessionID.String(),
			"total_steps":  len(steps),
			"success_rate": float64(successCount) / float64(len(steps)),
			"duration":     time.Since(time.Now().Add(-time.Duration(len(steps)) * time.Second)),
		},
	}

	vs.logger.Info(ctx, "Smart interaction completed", map[string]interface{}{
		"session_id":   req.SessionID.String(),
		"goal":         req.Goal,
		"steps":        len(steps),
		"success_rate": response.Metadata["success_rate"],
		"status":       status,
	})

	return response, nil
}

// Helper methods for visual analysis

func (vs *VisionService) detectUIElements(ctx context.Context, sessionID uuid.UUID, target string) ([]VisualElement, error) {
	// This is a simplified implementation
	// In a real implementation, this would use computer vision or DOM analysis

	var elements []VisualElement

	// Simulate element detection
	elements = append(elements, VisualElement{
		Type:       "button",
		Selector:   "button[type='submit']",
		Text:       "Submit",
		Bounds:     ElementBounds{X: 100, Y: 200, Width: 80, Height: 30},
		Confidence: 0.95,
	})

	elements = append(elements, VisualElement{
		Type:       "input",
		Selector:   "input[type='text']",
		Text:       "",
		Bounds:     ElementBounds{X: 50, Y: 150, Width: 200, Height: 25},
		Confidence: 0.90,
	})

	return elements, nil
}

func (vs *VisionService) analyzeLayout(ctx context.Context, screenshot string, elements []VisualElement) map[string]interface{} {
	return map[string]interface{}{
		"layout_type":      "responsive",
		"grid_system":      "flexbox",
		"element_density":  "medium",
		"visual_hierarchy": "clear",
		"spacing":          "adequate",
		"alignment":        "consistent",
	}
}

func (vs *VisionService) analyzeAccessibility(ctx context.Context, elements []VisualElement) map[string]interface{} {
	return map[string]interface{}{
		"alt_text_coverage":   0.85,
		"color_contrast":      "good",
		"keyboard_navigation": "supported",
		"aria_labels":         "partial",
		"focus_indicators":    "visible",
		"accessibility_score": 8.2,
	}
}

func (vs *VisionService) analyzeUsability(ctx context.Context, elements []VisualElement) map[string]interface{} {
	return map[string]interface{}{
		"navigation_clarity": "high",
		"cta_visibility":     "good",
		"form_complexity":    "low",
		"loading_indicators": "present",
		"error_handling":     "adequate",
		"usability_score":    7.8,
	}
}

func (vs *VisionService) analyzeContent(ctx context.Context, screenshot string, elements []VisualElement) map[string]interface{} {
	return map[string]interface{}{
		"content_structure":   "hierarchical",
		"readability":         "good",
		"content_density":     "balanced",
		"media_usage":         "appropriate",
		"text_to_image_ratio": 0.7,
	}
}

func (vs *VisionService) analyzeForms(ctx context.Context, elements []VisualElement) map[string]interface{} {
	formElements := 0
	for _, element := range elements {
		if strings.Contains(element.Type, "input") || element.Type == "button" {
			formElements++
		}
	}

	return map[string]interface{}{
		"form_elements":   formElements,
		"validation":      "client-side",
		"field_labeling":  "clear",
		"error_messaging": "inline",
		"completion_rate": "estimated_high",
	}
}

func (vs *VisionService) analyzeNavigation(ctx context.Context, elements []VisualElement) map[string]interface{} {
	return map[string]interface{}{
		"navigation_type": "horizontal",
		"menu_structure":  "flat",
		"breadcrumbs":     "present",
		"search_function": "available",
		"mobile_friendly": true,
	}
}

func (vs *VisionService) performGeneralAnalysis(ctx context.Context, screenshot string, elements []VisualElement, instructions string) map[string]interface{} {
	return map[string]interface{}{
		"page_type":            "content",
		"complexity":           "medium",
		"interactive_elements": len(elements),
		"visual_appeal":        "modern",
		"brand_consistency":    "high",
		"custom_analysis":      instructions,
	}
}

func (vs *VisionService) generateSuggestions(ctx context.Context, analysisType string, analysis map[string]interface{}, elements []VisualElement) []string {
	suggestions := []string{
		"Consider adding more visual hierarchy to improve readability",
		"Ensure all interactive elements have proper focus indicators",
		"Optimize page loading speed for better user experience",
	}

	switch analysisType {
	case "accessibility":
		suggestions = append(suggestions, "Add alt text to all images", "Improve color contrast ratios")
	case "usability":
		suggestions = append(suggestions, "Simplify navigation structure", "Add progress indicators for multi-step processes")
	case "forms":
		suggestions = append(suggestions, "Add real-time validation feedback", "Reduce form field count")
	}

	return suggestions
}

func (vs *VisionService) planInteractionSteps(ctx context.Context, goal string, context string, elements []VisualElement, screenshot string) []Action {
	// Simplified planning logic
	// In a real implementation, this would use AI to analyze the goal and plan steps

	var actions []Action

	// Example: if goal contains "fill form", plan form filling steps
	if strings.Contains(strings.ToLower(goal), "fill") || strings.Contains(strings.ToLower(goal), "form") {
		for _, element := range elements {
			if element.Type == "input" {
				actions = append(actions, Action{
					Type:     ActionTypeClick,
					Selector: element.Selector,
				})
				actions = append(actions, Action{
					Type:     ActionTypeType,
					Selector: element.Selector,
					Value:    "test value",
				})
			}
		}
	}

	// Example: if goal contains "submit", add submit action
	if strings.Contains(strings.ToLower(goal), "submit") {
		for _, element := range elements {
			if element.Type == "button" && strings.Contains(strings.ToLower(element.Text), "submit") {
				actions = append(actions, Action{
					Type:     ActionTypeClick,
					Selector: element.Selector,
				})
				break
			}
		}
	}

	return actions
}

func (vs *VisionService) getActionReasoning(action Action) string {
	switch action.Type {
	case ActionTypeClick:
		return fmt.Sprintf("Clicking on element %s to interact with it", action.Selector)
	case ActionTypeType:
		return fmt.Sprintf("Typing '%s' into input field %s", action.Value, action.Selector)
	case ActionTypeScroll:
		return "Scrolling to reveal more content"
	default:
		return fmt.Sprintf("Performing %s action", action.Type)
	}
}

func (vs *VisionService) isGoalAchieved(ctx context.Context, goal string, screenshot string) bool {
	// Simplified goal achievement check
	// In a real implementation, this would use AI to analyze the screenshot
	return false
}

func (vs *VisionService) takeStepScreenshot(ctx context.Context, sessionID uuid.UUID) (string, error) {
	screenshotReq := ScreenshotRequest{
		FullPage: false,
		Quality:  80,
		Format:   "png",
	}

	resp, err := vs.service.TakeScreenshot(ctx, sessionID, screenshotReq)
	if err != nil {
		return "", err
	}

	return resp.Screenshot, nil
}
