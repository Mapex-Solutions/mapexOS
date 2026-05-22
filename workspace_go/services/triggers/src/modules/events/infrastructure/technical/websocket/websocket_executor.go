package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"triggers/src/modules/events/application/ports"

	"github.com/gorilla/websocket"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// NewWebSocketExecutor creates a new WebSocket trigger executor adapter.
func NewWebSocketExecutor() ports.TriggerExecutor {
	return &WebSocketExecutor{
		connectTimeout:   10 * time.Second,
		handshakeTimeout: 10 * time.Second,
		writeTimeout:     5 * time.Second,
	}
}

// Execute sends a message to a WebSocket server based on the trigger configuration.
//
// This is the concrete implementation of the ports.TriggerExecutor interface.
//
// Steps:
// 1. Extract WebSocket config (url, headers, message)
// 2. Build connection dialer with headers
// 3. Connect to server
// 4. Send message
// 5. Close connection
//
// Parameters:
//   - ctx: Context for controlling cancellation and timeouts
//   - config: Trigger configuration (websocket field) with placeholders already resolved
//
// Returns:
//   - error: If WebSocket send fails
func (e *WebSocketExecutor) Execute(ctx context.Context, config map[string]interface{}) error {
	// Extract websocket config
	wsConfig, ok := config["websocket"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("WebSocket trigger config missing 'websocket' field")
	}

	// Extract URL (required)
	url, ok := wsConfig["url"].(string)
	if !ok || url == "" {
		return fmt.Errorf("WebSocket trigger config missing required field 'url'")
	}

	// Extract message payload
	message, err := e.extractMessage(wsConfig)
	if err != nil {
		return fmt.Errorf("failed to extract WebSocket message: %w", err)
	}

	// Determine message type
	msgType := websocket.TextMessage
	if typeStr, ok := wsConfig["messageType"].(string); ok && typeStr == "binary" {
		msgType = websocket.BinaryMessage
	}

	// Build request headers
	headers := http.Header{}
	if headerMap, ok := wsConfig["headers"].(map[string]interface{}); ok {
		for key, value := range headerMap {
			if strVal, ok := value.(string); ok {
				headers.Add(key, strVal)
			}
		}
	}

	// Check context cancellation before connecting
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	logger.Info(fmt.Sprintf("[INFRA:WebSocketExecutor] Connecting to %s", url))

	// Create dialer with timeout
	dialer := websocket.Dialer{
		HandshakeTimeout: e.handshakeTimeout,
	}

	// Connect to WebSocket server
	conn, _, err := dialer.DialContext(ctx, url, headers)
	if err != nil {
		logger.Error(err, fmt.Sprintf("[INFRA:WebSocketExecutor] Connection failed: %s", url))
		return fmt.Errorf("WebSocket connection failed: %w", err)
	}
	defer conn.Close()

	// Check context cancellation before writing
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Set write deadline
	if err := conn.SetWriteDeadline(time.Now().Add(e.writeTimeout)); err != nil {
		return fmt.Errorf("failed to set write deadline: %w", err)
	}

	logger.Info(fmt.Sprintf("[INFRA:WebSocketExecutor] Sending message to %s", url))

	// Send message
	if err := conn.WriteMessage(msgType, message); err != nil {
		logger.Error(err, fmt.Sprintf("[INFRA:WebSocketExecutor] Write failed: %s", url))
		return fmt.Errorf("WebSocket write failed: %w", err)
	}

	// Send close frame gracefully
	closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	if err := conn.WriteControl(websocket.CloseMessage, closeMsg, time.Now().Add(time.Second)); err != nil {
		// Log but don't fail - message was already sent
		logger.Debug(fmt.Sprintf("[INFRA:WebSocketExecutor] Close handshake warning: %v", err))
	}

	logger.Info(fmt.Sprintf("[INFRA:WebSocketExecutor] Message sent successfully to %s", url))
	return nil
}

// GetType returns the trigger type this executor handles.
func (e *WebSocketExecutor) GetType() string {
	return "websocket"
}

// extractMessage extracts and serializes the message payload from config.
func (e *WebSocketExecutor) extractMessage(config map[string]interface{}) ([]byte, error) {
	message, exists := config["message"]
	if !exists {
		return []byte(""), nil
	}

	switch m := message.(type) {
	case string:
		return []byte(m), nil
	case map[string]interface{}:
		return json.Marshal(m)
	default:
		return json.Marshal(m)
	}
}
