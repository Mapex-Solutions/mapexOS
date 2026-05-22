package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"triggers/src/modules/events/application/ports"

	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// NewHTTPExecutor creates a new HTTP trigger executor adapter.
func NewHTTPExecutor() ports.TriggerExecutor {
	return &HTTPExecutor{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Execute performs an HTTP request based on the trigger configuration.
//
// This is the concrete implementation of the ports.TriggerExecutor interface.
//
// Steps:
// Extract config fields (endpoint, method, headers, body, timeout)
// Build HTTP request with context
// Execute request
// Log response
// Return error if request fails or status code is not 2xx
//
// Parameters:
//   - ctx: Context for controlling cancellation and timeouts
//   - config: Trigger configuration (http field) with placeholders already resolved
//
// Returns:
//   - error: If HTTP request fails or returns non-2xx status
func (e *HTTPExecutor) Execute(ctx context.Context, config map[string]interface{}) error {
	// Extract http config (the application already extracted this from trigger.config.http)
	httpConfig, ok := config["http"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("HTTP trigger config missing 'http' field")
	}

	// Extract endpoint (required)
	endpoint, ok := httpConfig["endpoint"].(string)
	if !ok || endpoint == "" {
		return fmt.Errorf("HTTP trigger config missing required field 'endpoint'")
	}

	// Extract method (required)
	method, ok := httpConfig["method"].(string)
	if !ok || method == "" {
		method = "POST" // Default to POST if not specified
	}

	// Extract body (optional)
	var bodyReader io.Reader
	if bodyData, exists := httpConfig["body"]; exists && bodyData != nil {
		// Serialize body to JSON
		bodyBytes, err := json.Marshal(bodyData)
		if err != nil {
			return fmt.Errorf("failed to serialize HTTP body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Build HTTP request
	req, err := http.NewRequestWithContext(ctx, method, endpoint, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Add headers (optional)
	if headersData, exists := httpConfig["headers"]; exists {
		if headers, ok := headersData.(map[string]interface{}); ok {
			for key, value := range headers {
				if valueStr, ok := value.(string); ok {
					req.Header.Set(key, valueStr)
				}
			}
		}
	}

	// Set default Content-Type if body exists and Content-Type not set
	if bodyReader != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Apply custom timeout if specified (in milliseconds, convert to seconds)
	if timeoutVal, exists := httpConfig["timeout"]; exists {
		if timeoutFloat, ok := timeoutVal.(float64); ok {
			// Timeout is in milliseconds in config, convert to seconds
			e.client.Timeout = time.Duration(timeoutFloat) * time.Millisecond
		}
	}

	// Execute HTTP request
	logger.Info(fmt.Sprintf("[INFRA:HTTPExecutor] Executing %s %s", method, endpoint))

	resp, err := e.client.Do(req)
	if err != nil {
		logger.Error(err, fmt.Sprintf("[INFRA:HTTPExecutor] Request failed: %s %s", method, endpoint))
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body for logging
	respBody, _ := io.ReadAll(resp.Body)

	// Check response status
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		logger.Info(fmt.Sprintf("[INFRA:HTTPExecutor] Success: %s %s - Status: %d", method, endpoint, resp.StatusCode))
		logger.Debug(fmt.Sprintf("[INFRA:HTTPExecutor] Response: %s", string(respBody)))
		return nil
	}

	// Request failed with non-2xx status
	logger.Warn(fmt.Sprintf("[INFRA:HTTPExecutor] Failed: %s %s - Status: %d, Body: %s", method, endpoint, resp.StatusCode, string(respBody)))
	return fmt.Errorf("HTTP request returned non-2xx status: %d", resp.StatusCode)
}

// GetType returns the trigger type this executor handles.
func (e *HTTPExecutor) GetType() string {
	return "http"
}
