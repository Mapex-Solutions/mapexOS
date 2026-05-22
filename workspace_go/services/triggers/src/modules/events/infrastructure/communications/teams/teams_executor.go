package teams

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

// NewTeamsExecutor creates a new Teams webhook executor.
func NewTeamsExecutor() ports.TriggerExecutor {
	return &TeamsExecutor{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Execute sends a message to Microsoft Teams via webhook.
//
// Parameters:
//   - ctx: Context for controlling cancellation and timeouts
//   - config: Trigger configuration with placeholders already resolved
//
// Returns:
//   - error: If webhook request fails
func (e *TeamsExecutor) Execute(ctx context.Context, config map[string]interface{}) error {
	// Unwrap the union-shaped config — the router resolves the trigger
	// config and forwards the full document; the teams-specific fields
	// live under config["teams"], same shape as HTTP / MQTT / Email.
	teamsCfg, ok := config["teams"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("Teams trigger config missing 'teams' field")
	}

	// Extract webhook URL (required)
	webhookUrl, ok := teamsCfg["webhookUrl"].(string)
	if !ok || webhookUrl == "" {
		return fmt.Errorf("Teams trigger config missing required field 'webhookUrl'")
	}

	// Build MessageCard payload
	messageCard := make(map[string]interface{})
	messageCard["@type"] = "MessageCard"
	messageCard["@context"] = "https://schema.org/extensions"

	// Title (required by contract)
	if title, exists := teamsCfg["title"]; exists {
		messageCard["title"] = title
	}

	// Text (main message body, required by contract)
	if text, exists := teamsCfg["text"]; exists {
		messageCard["text"] = text
	}

	// Theme color (optional)
	if themeColor, exists := teamsCfg["themeColor"]; exists {
		messageCard["themeColor"] = themeColor
	}

	// Sections (optional)
	if sections, exists := teamsCfg["sections"]; exists {
		messageCard["sections"] = sections
	}

	// Serialize payload to JSON
	payloadBytes, err := json.Marshal(messageCard)
	if err != nil {
		return fmt.Errorf("failed to serialize Teams message: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", webhookUrl, bytes.NewReader(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create Teams webhook request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Execute webhook request
	logger.Info(fmt.Sprintf("[INFRA:TeamsExecutor] Sending message to Teams webhook"))

	resp, err := e.client.Do(req)
	if err != nil {
		logger.Error(err, "[INFRA:TeamsExecutor] Webhook request failed")
		return fmt.Errorf("Teams webhook request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, _ := io.ReadAll(resp.Body)

	// Check response status
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		logger.Info(fmt.Sprintf("[INFRA:TeamsExecutor] Message sent successfully - Status: %d", resp.StatusCode))
		return nil
	}

	logger.Warn(fmt.Sprintf("[INFRA:TeamsExecutor] Failed - Status: %d, Body: %s", resp.StatusCode, string(respBody)))
	return fmt.Errorf("Teams webhook returned non-2xx status: %d", resp.StatusCode)
}

// GetType returns the trigger type this executor handles.
func (e *TeamsExecutor) GetType() string {
	return "teams"
}
