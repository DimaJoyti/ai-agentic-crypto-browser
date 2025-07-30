//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasicPageLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	baseURL := getEnv("TEST_E2E_BASE_URL", "http://localhost:3000")
	headless := getEnv("TEST_BROWSER_HEADLESS", "true") == "true"

	// Check if the application is running
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(baseURL)
	if err != nil {
		t.Skipf("Application not running at %s, skipping E2E test: %v", baseURL, err)
		return
	}
	resp.Body.Close()

	// Setup Chrome context
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", headless),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Set timeout
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var title string
	err = chromedp.Run(ctx,
		chromedp.Navigate(baseURL),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.Title(&title),
	)

	require.NoError(t, err)
	assert.NotEmpty(t, title, "Page should have a title")

	t.Logf("Page title: %s", title)
}

func TestAPIHealthEndpoint(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	baseURL := getEnv("TEST_E2E_BASE_URL", "http://localhost:3000")
	apiURL := getEnv("TEST_API_BASE_URL", "http://localhost:8080")

	client := &http.Client{Timeout: 10 * time.Second}

	// Test API health endpoint
	resp, err := client.Get(apiURL + "/health")
	if err != nil {
		t.Skipf("API not running at %s, skipping test: %v", apiURL, err)
		return
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test frontend can reach API (if applicable)
	resp, err = client.Get(baseURL + "/api/health")
	if err == nil {
		defer resp.Body.Close()
		// If the frontend has an API proxy, it should work
		assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotFound)
	}
}

func TestNavigationFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	baseURL := getEnv("TEST_E2E_BASE_URL", "http://localhost:3000")
	headless := getEnv("TEST_BROWSER_HEADLESS", "true") == "true"

	// Check if the application is running
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(baseURL)
	if err != nil {
		t.Skipf("Application not running at %s, skipping E2E test: %v", baseURL, err)
		return
	}
	resp.Body.Close()

	// Setup Chrome context
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", headless),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Set timeout
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	err = chromedp.Run(ctx,
		// Navigate to home page
		chromedp.Navigate(baseURL),
		chromedp.WaitVisible("body", chromedp.ByQuery),

		// Wait a bit for the page to fully load
		chromedp.Sleep(2*time.Second),

		// Take a screenshot for debugging if needed
		chromedp.ActionFunc(func(ctx context.Context) error {
			if os.Getenv("TEST_TAKE_SCREENSHOTS") == "true" {
				var buf []byte
				if err := chromedp.CaptureScreenshot(&buf).Do(ctx); err == nil {
					os.WriteFile("screenshot_navigation.png", buf, 0644)
				}
			}
			return nil
		}),
	)

	require.NoError(t, err, "Navigation flow should complete without errors")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
