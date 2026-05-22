package slack

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

// NewSlackExecutor creates a new Slack webhook executor.
func NewSlackExecutor() ports.TriggerExecutor {
	return &SlackExecutor{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Execute sends a message to Slack via webhook.
//
// Parameters:
//   - ctx: Context for controlling cancellation and timeouts
//   - config: Trigger configuration with placeholders already resolved
//
// Returns:
//   - error: If webhook request fails
func (e *SlackExecutor) Execute(ctx context.Context, config map[string]interface{}) error {
	// Unwrap the union-shaped config — the router resolves the trigger
	// config and forwards the full document; the slack-specific fields
	// live under config["slack"], same shape as HTTP / MQTT / Email.
	slackCfg, ok := config["slack"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("Slack trigger config missing 'slack' field")
	}

	// Extract webhook URL (required)
	webhookUrl, ok := slackCfg["webhookUrl"].(string)
	if !ok || webhookUrl == "" {
		return fmt.Errorf("Slack trigger config missing required field 'webhookUrl'")
	}

	// Build Slack message payload
	payload := make(map[string]interface{})

	// Message (required) — the contract names this field `message`;
	// Slack's webhook protocol delivers it under the `text` key.
	message, ok := slackCfg["message"].(string)
	if !ok || message == "" {
		return fmt.Errorf("Slack trigger config missing required field 'message'")
	}
	payload["text"] = message

	// Username (optional)
	if username, exists := slackCfg["username"]; exists {
		payload["username"] = username
	}

	// Icon emoji (optional) — contract spells it iconEmoji.
	if iconEmoji, exists := slackCfg["iconEmoji"]; exists {
		payload["icon_emoji"] = iconEmoji
	}

	// Channel override (optional)
	if channel, exists := slackCfg["channel"]; exists {
		payload["channel"] = channel
	}

	// Blocks for rich formatting (optional)
	if blocks, exists := slackCfg["blocks"]; exists {
		payload["blocks"] = blocks
	}

	// Attachments (optional, legacy)
	if attachments, exists := slackCfg["attachments"]; exists {
		payload["attachments"] = attachments
	}

	// Serialize payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to serialize Slack message: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", webhookUrl, bytes.NewReader(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create Slack webhook request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Execute webhook request
	logger.Info(fmt.Sprintf("[INFRA:SlackExecutor] Sending message to Slack webhook"))

	resp, err := e.client.Do(req)
	if err != nil {
		logger.Error(err, "[INFRA:SlackExecutor] Webhook request failed")
		return fmt.Errorf("Slack webhook request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, _ := io.ReadAll(resp.Body)

	// Check response
	// Slack returns 200 with "ok" for success
	if resp.StatusCode == 200 && string(respBody) == "ok" {
		logger.Info("[INFRA:SlackExecutor] Message sent successfully")
		return nil
	}

	logger.Warn(fmt.Sprintf("[INFRA:SlackExecutor] Failed - Status: %d, Body: %s", resp.StatusCode, string(respBody)))
	return fmt.Errorf("Slack webhook failed: %d - %s", resp.StatusCode, string(respBody))
}

// GetType returns the trigger type this executor handles.
func (e *SlackExecutor) GetType() string {
	return "slack"
}
