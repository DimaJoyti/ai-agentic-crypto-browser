package browser

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/ai-agentic-browser/internal/config"
	"github.com/ai-agentic-browser/pkg/database"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
)

// Service provides browser automation functionality
type Service struct {
	db        *database.DB
	redis     *database.RedisClient
	config    config.BrowserConfig
	logger    *observability.Logger
	instances map[string]*BrowserInstance
}

// NewService creates a new browser service
func NewService(db *database.DB, redis *database.RedisClient, cfg config.BrowserConfig, logger *observability.Logger) *Service {
	return &Service{
		db:        db,
		redis:     redis,
		config:    cfg,
		logger:    logger,
		instances: make(map[string]*BrowserInstance),
	}
}

// CreateSession creates a new browser session
func (s *Service) CreateSession(ctx context.Context, userID uuid.UUID, req SessionCreateRequest) (*BrowserSession, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("browser-service").Start(ctx, "browser.CreateSession")
	defer span.End()

	session := &BrowserSession{
		ID:          uuid.New(),
		UserID:      userID,
		SessionName: req.SessionName,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if session.SessionName == "" {
		session.SessionName = fmt.Sprintf("Session %s", session.ID.String()[:8])
	}

	// Insert session into database
	query := `
		INSERT INTO browser_sessions (id, user_id, session_name, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := s.db.ExecContext(ctx, query, session.ID, session.UserID, session.SessionName, session.IsActive, session.CreatedAt, session.UpdatedAt)
	if err != nil {
		s.logger.Error(ctx, "Failed to create browser session", err)
		return nil, fmt.Errorf("failed to create browser session: %w", err)
	}

	s.logger.Info(ctx, "Browser session created", map[string]interface{}{
		"session_id": session.ID.String(),
		"user_id":    userID.String(),
	})

	return session, nil
}

// Navigate navigates to a URL in a browser context
func (s *Service) Navigate(ctx context.Context, sessionID uuid.UUID, req NavigateRequest) (*NavigateResponse, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("browser-service").Start(ctx, "browser.Navigate")
	defer span.End()

	startTime := time.Now()

	// Create browser context with options
	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", s.config.Headless),
		chromedp.Flag("disable-gpu", s.config.DisableGPU),
		chromedp.Flag("no-sandbox", s.config.NoSandbox),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-background-timer-throttling", false),
		chromedp.Flag("disable-backgrounding-occluded-windows", false),
		chromedp.Flag("disable-renderer-backgrounding", false),
	}

	if req.UserAgent != "" {
		opts = append(opts, chromedp.UserAgent(req.UserAgent))
	}

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	browserCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Set timeout
	timeout := s.config.Timeout
	if req.Timeout > 0 {
		timeout = time.Duration(req.Timeout) * time.Second
	}

	timeoutCtx, cancel := context.WithTimeout(browserCtx, timeout)
	defer cancel()

	var title string
	var screenshot []byte

	// Build tasks
	tasks := []chromedp.Action{
		chromedp.Navigate(req.URL),
	}

	// Add wait condition if specified
	if req.WaitForSelector != "" {
		tasks = append(tasks, chromedp.WaitVisible(req.WaitForSelector))
	} else {
		tasks = append(tasks, chromedp.WaitReady("body"))
	}

	// Get page title
	tasks = append(tasks, chromedp.Title(&title))

	// Take screenshot
	tasks = append(tasks, chromedp.CaptureScreenshot(&screenshot))

	// Execute tasks
	err := chromedp.Run(timeoutCtx, tasks...)
	if err != nil {
		s.logger.Error(ctx, "Navigation failed", err, map[string]interface{}{
			"url":        req.URL,
			"session_id": sessionID.String(),
		})
		return &NavigateResponse{
			Success: false,
			URL:     req.URL,
			Error:   err.Error(),
		}, nil
	}

	loadTime := time.Since(startTime)

	// Encode screenshot
	screenshotB64 := base64.StdEncoding.EncodeToString(screenshot)

	response := &NavigateResponse{
		Success:    true,
		URL:        req.URL,
		Title:      title,
		LoadTime:   loadTime,
		Screenshot: screenshotB64,
		Metadata: map[string]interface{}{
			"session_id": sessionID.String(),
			"timestamp":  time.Now(),
		},
	}

	s.logger.Info(ctx, "Navigation completed", map[string]interface{}{
		"url":        req.URL,
		"title":      title,
		"load_time":  loadTime.Milliseconds(),
		"session_id": sessionID.String(),
	})

	return response, nil
}

// Interact performs interactions with page elements
func (s *Service) Interact(ctx context.Context, sessionID uuid.UUID, req InteractRequest) (*InteractResponse, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("browser-service").Start(ctx, "browser.Interact")
	defer span.End()

	// Create browser context
	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", s.config.Headless),
		chromedp.Flag("disable-gpu", s.config.DisableGPU),
		chromedp.Flag("no-sandbox", s.config.NoSandbox),
	}

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	browserCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	timeoutCtx, cancel := context.WithTimeout(browserCtx, s.config.Timeout)
	defer cancel()

	var results []ActionResult
	var screenshots []string

	for i, action := range req.Actions {
		startTime := time.Now()
		result := ActionResult{
			Action:  action,
			Success: true,
		}

		// Execute action
		err := s.executeAction(timeoutCtx, action)
		if err != nil {
			result.Success = false
			result.Error = err.Error()
			s.logger.Error(ctx, "Action failed", err, map[string]interface{}{
				"action_type":  action.Type,
				"selector":     action.Selector,
				"action_index": i,
			})
		}

		result.Duration = time.Since(startTime)
		results = append(results, result)

		// Take screenshot if requested or after failed action
		if req.Screenshot || !result.Success {
			var screenshot []byte
			if err := chromedp.Run(timeoutCtx, chromedp.CaptureScreenshot(&screenshot)); err == nil {
				screenshotB64 := base64.StdEncoding.EncodeToString(screenshot)
				screenshots = append(screenshots, screenshotB64)
			}
		}

		// Wait between actions if specified
		if req.WaitBetween > 0 && i < len(req.Actions)-1 {
			time.Sleep(time.Duration(req.WaitBetween) * time.Millisecond)
		}
	}

	// Check overall success
	success := true
	for _, result := range results {
		if !result.Success {
			success = false
			break
		}
	}

	response := &InteractResponse{
		Success:     success,
		Results:     results,
		Screenshots: screenshots,
		Metadata: map[string]interface{}{
			"session_id":    sessionID.String(),
			"actions_count": len(req.Actions),
			"timestamp":     time.Now(),
		},
	}

	s.logger.Info(ctx, "Interaction completed", map[string]interface{}{
		"session_id":    sessionID.String(),
		"actions_count": len(req.Actions),
		"success":       success,
	})

	return response, nil
}

// executeAction executes a single browser action
func (s *Service) executeAction(ctx context.Context, action Action) error {
	switch action.Type {
	case ActionTypeClick:
		return chromedp.Run(ctx, chromedp.Click(action.Selector))

	case ActionTypeType:
		return chromedp.Run(ctx, chromedp.SendKeys(action.Selector, action.Value))

	case ActionTypeClear:
		return chromedp.Run(ctx, chromedp.Clear(action.Selector))

	case ActionTypeSelect:
		return chromedp.Run(ctx, chromedp.SetValue(action.Selector, action.Value))

	case ActionTypeScroll:
		if action.Selector != "" {
			return chromedp.Run(ctx, chromedp.ScrollIntoView(action.Selector))
		}
		// Scroll by pixels if value is provided
		if action.Value != "" {
			script := fmt.Sprintf("window.scrollBy(0, %s)", action.Value)
			return chromedp.Run(ctx, chromedp.Evaluate(script, nil))
		}
		return chromedp.Run(ctx, chromedp.Evaluate("window.scrollTo(0, document.body.scrollHeight)", nil))

	case ActionTypeWait:
		if action.Selector != "" {
			return chromedp.Run(ctx, chromedp.WaitVisible(action.Selector))
		}
		if action.Value != "" {
			if duration, err := time.ParseDuration(action.Value); err == nil {
				time.Sleep(duration)
				return nil
			}
		}
		time.Sleep(1 * time.Second)
		return nil

	case ActionTypeHover:
		// Use JavaScript to trigger hover event since chromedp doesn't have MouseOver
		script := fmt.Sprintf(`
			const element = document.querySelector('%s');
			if (element) {
				const event = new MouseEvent('mouseover', {
					view: window,
					bubbles: true,
					cancelable: true
				});
				element.dispatchEvent(event);
			}
		`, action.Selector)
		return chromedp.Run(ctx, chromedp.Evaluate(script, nil))

	case ActionTypeKeyPress:
		return chromedp.Run(ctx, chromedp.KeyEvent(action.Value))

	case ActionTypeSubmit:
		if action.Selector != "" {
			return chromedp.Run(ctx, chromedp.Submit(action.Selector))
		}
		return chromedp.Run(ctx, chromedp.KeyEvent("Enter"))

	case ActionTypeRefresh:
		return chromedp.Run(ctx, chromedp.Reload())

	case ActionTypeGoBack:
		return chromedp.Run(ctx, chromedp.NavigateBack())

	case ActionTypeGoForward:
		return chromedp.Run(ctx, chromedp.NavigateForward())

	default:
		return fmt.Errorf("unsupported action type: %s", action.Type)
	}
}

// Extract extracts content from the current page
func (s *Service) Extract(ctx context.Context, sessionID uuid.UUID, req ExtractRequest) (*ExtractResponse, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("browser-service").Start(ctx, "browser.Extract")
	defer span.End()

	// Create browser context
	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", s.config.Headless),
		chromedp.Flag("disable-gpu", s.config.DisableGPU),
		chromedp.Flag("no-sandbox", s.config.NoSandbox),
	}

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	browserCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	timeoutCtx, cancel := context.WithTimeout(browserCtx, s.config.Timeout)
	defer cancel()

	data := make(map[string]interface{})

	// Extract based on data type
	switch req.DataType {
	case "text":
		data = s.extractText(timeoutCtx, req.Selectors)
	case "links":
		data = s.extractLinks(timeoutCtx, req.Selectors)
	case "images":
		data = s.extractImages(timeoutCtx, req.Selectors)
	case "tables":
		data = s.extractTables(timeoutCtx, req.Selectors)
	case "forms":
		data = s.extractForms(timeoutCtx, req.Selectors)
	default:
		// Extract all types
		data["text"] = s.extractText(timeoutCtx, req.Selectors)
		data["links"] = s.extractLinks(timeoutCtx, req.Selectors)
		data["images"] = s.extractImages(timeoutCtx, req.Selectors)
	}

	// Take screenshot
	var screenshot []byte
	chromedp.Run(timeoutCtx, chromedp.CaptureScreenshot(&screenshot))
	screenshotB64 := base64.StdEncoding.EncodeToString(screenshot)

	response := &ExtractResponse{
		Success:    true,
		Data:       data,
		Screenshot: screenshotB64,
		Metadata: map[string]interface{}{
			"session_id": sessionID.String(),
			"timestamp":  time.Now(),
			"data_type":  req.DataType,
		},
	}

	s.logger.Info(ctx, "Content extraction completed", map[string]interface{}{
		"session_id": sessionID.String(),
		"data_type":  req.DataType,
		"selectors":  len(req.Selectors),
	})

	return response, nil
}

// Helper methods for content extraction
func (s *Service) extractText(ctx context.Context, selectors []string) map[string]interface{} {
	result := make(map[string]interface{})

	if len(selectors) == 0 {
		selectors = []string{"body"}
	}

	for _, selector := range selectors {
		var text string
		if err := chromedp.Run(ctx, chromedp.Text(selector, &text)); err == nil {
			result[selector] = text
		}
	}

	return result
}

func (s *Service) extractLinks(ctx context.Context, selectors []string) map[string]interface{} {
	result := make(map[string]interface{})

	selector := "a[href]"
	if len(selectors) > 0 {
		selector = selectors[0]
	}

	var links []map[string]string
	script := fmt.Sprintf(`
		Array.from(document.querySelectorAll('%s')).map(a => ({
			url: a.href,
			text: a.textContent.trim(),
			title: a.title || ''
		}))
	`, selector)

	if err := chromedp.Run(ctx, chromedp.Evaluate(script, &links)); err == nil {
		result["links"] = links
	}

	return result
}

func (s *Service) extractImages(ctx context.Context, selectors []string) map[string]interface{} {
	result := make(map[string]interface{})

	selector := "img"
	if len(selectors) > 0 {
		selector = selectors[0]
	}

	var images []map[string]interface{}
	script := fmt.Sprintf(`
		Array.from(document.querySelectorAll('%s')).map(img => ({
			url: img.src,
			alt: img.alt || '',
			title: img.title || '',
			width: img.naturalWidth || img.width,
			height: img.naturalHeight || img.height
		}))
	`, selector)

	if err := chromedp.Run(ctx, chromedp.Evaluate(script, &images)); err == nil {
		result["images"] = images
	}

	return result
}

func (s *Service) extractTables(ctx context.Context, selectors []string) map[string]interface{} {
	result := make(map[string]interface{})

	selector := "table"
	if len(selectors) > 0 {
		selector = selectors[0]
	}

	var tables []map[string]interface{}
	script := fmt.Sprintf(`
		Array.from(document.querySelectorAll('%s')).map(table => {
			const headers = Array.from(table.querySelectorAll('th')).map(th => th.textContent.trim());
			const rows = Array.from(table.querySelectorAll('tr')).slice(headers.length > 0 ? 1 : 0).map(tr => 
				Array.from(tr.querySelectorAll('td')).map(td => td.textContent.trim())
			);
			return { headers, rows };
		})
	`, selector)

	if err := chromedp.Run(ctx, chromedp.Evaluate(script, &tables)); err == nil {
		result["tables"] = tables
	}

	return result
}

func (s *Service) extractForms(ctx context.Context, selectors []string) map[string]interface{} {
	result := make(map[string]interface{})

	selector := "form"
	if len(selectors) > 0 {
		selector = selectors[0]
	}

	var forms []map[string]interface{}
	script := fmt.Sprintf(`
		Array.from(document.querySelectorAll('%s')).map(form => ({
			action: form.action || '',
			method: form.method || 'GET',
			fields: Array.from(form.querySelectorAll('input, select, textarea')).map(field => ({
				name: field.name || '',
				type: field.type || field.tagName.toLowerCase(),
				value: field.value || '',
				placeholder: field.placeholder || '',
				required: field.required || false
			}))
		}))
	`, selector)

	if err := chromedp.Run(ctx, chromedp.Evaluate(script, &forms)); err == nil {
		result["forms"] = forms
	}

	return result
}

// TakeScreenshot takes a screenshot of the current page
func (s *Service) TakeScreenshot(ctx context.Context, sessionID uuid.UUID, req ScreenshotRequest) (*ScreenshotResponse, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("browser-service").Start(ctx, "browser.TakeScreenshot")
	defer span.End()

	// Create browser context
	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", s.config.Headless),
		chromedp.Flag("disable-gpu", s.config.DisableGPU),
		chromedp.Flag("no-sandbox", s.config.NoSandbox),
	}

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	browserCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	timeoutCtx, cancel := context.WithTimeout(browserCtx, s.config.Timeout)
	defer cancel()

	var screenshot []byte
	var err error

	if req.Selector != "" {
		// Element screenshot
		err = chromedp.Run(timeoutCtx, chromedp.Screenshot(req.Selector, &screenshot))
	} else if req.FullPage {
		// Full page screenshot
		err = chromedp.Run(timeoutCtx, chromedp.FullScreenshot(&screenshot, req.Quality))
	} else {
		// Viewport screenshot
		err = chromedp.Run(timeoutCtx, chromedp.CaptureScreenshot(&screenshot))
	}

	if err != nil {
		s.logger.Error(ctx, "Screenshot failed", err)
		return &ScreenshotResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	format := req.Format
	if format == "" {
		format = "png"
	}

	screenshotB64 := base64.StdEncoding.EncodeToString(screenshot)

	response := &ScreenshotResponse{
		Success:    true,
		Screenshot: screenshotB64,
		Format:     format,
		Size:       len(screenshot),
		Dimensions: map[string]int{
			"width":  req.Width,
			"height": req.Height,
		},
		Metadata: map[string]interface{}{
			"session_id": sessionID.String(),
			"timestamp":  time.Now(),
		},
	}

	s.logger.Info(ctx, "Screenshot taken", map[string]interface{}{
		"session_id": sessionID.String(),
		"size":       len(screenshot),
		"format":     format,
	})

	return response, nil
}
