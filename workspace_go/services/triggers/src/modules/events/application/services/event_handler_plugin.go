package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// executePluginAction executes a fully resolved plugin action pipeline.
// Routes by action.type: http (inline), mqtt/nats/email (via ExecutorRegistry).
// Returns the output handle ("success" or "error"), output data, and any error.
func (s *EventService) executePluginAction(data map[string]interface{}) (string, interface{}, error) {
	action, _ := data["action"].(map[string]interface{})
	if action == nil {
		return "error", nil, fmt.Errorf("plugin action is missing")
	}

	actionType, _ := action["type"].(string)
	logger.Info(fmt.Sprintf("[SERVICE:Event] Executing plugin action: type=%s", actionType))

	switch actionType {
	case "http":
		return s.executePluginHTTP(action)
	case "mqtt", "nats", "email", "rabbitmq", "websocket":
		return s.executePluginViaRegistry(actionType, action)
	default:
		return "error", nil, fmt.Errorf("unsupported plugin action type: %s", actionType)
	}
}

// executePluginHTTP executes a fully resolved HTTP plugin action.
func (s *EventService) executePluginHTTP(action map[string]interface{}) (string, interface{}, error) {
	httpDef, _ := action["http"].(map[string]interface{})
	if httpDef == nil {
		return "error", nil, fmt.Errorf("http action definition is missing")
	}

	url, _ := httpDef["path"].(string)
	method, _ := httpDef["method"].(string)
	if method == "" {
		method = "POST"
	}

	var bodyBytes []byte
	if body := httpDef["body"]; body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return "error", nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	headers := make(map[string]string)
	if h, ok := httpDef["headers"].(map[string]interface{}); ok {
		for k, v := range h {
			headers[k] = fmt.Sprintf("%v", v)
		}
	}
	if _, hasContentType := headers["Content-Type"]; !hasContentType && bodyBytes != nil {
		headers["Content-Type"] = "application/json"
	}

	logger.Debug(fmt.Sprintf("[SERVICE:Event] Plugin HTTP: %s %s", method, url))

	var httpReq *http.Request
	if bodyBytes != nil {
		httpReq, _ = http.NewRequest(method, url, bytes.NewReader(bodyBytes))
	} else {
		httpReq, _ = http.NewRequest(method, url, nil)
	}
	if httpReq == nil {
		return "error", nil, fmt.Errorf("failed to create HTTP request")
	}
	for k, v := range headers {
		httpReq.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "error", nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "error", nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var responseData interface{}
	if err := json.Unmarshal(respBody, &responseData); err != nil {
		responseData = string(respBody)
	}

	outputDef, _ := action["output"].(map[string]interface{})
	if outputDef != nil {
		if dataPath, ok := outputDef["dataPath"].(string); ok && dataPath != "" {
			responseData = navigatePath(responseData, dataPath)
		}
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return "success", responseData, nil
	}
	return "error", responseData, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
}

// executePluginViaRegistry delegates a plugin action to the existing ExecutorRegistry.
// The action's type-specific config (mqtt/nats/email) is passed as the executor config.
func (s *EventService) executePluginViaRegistry(actionType string, action map[string]interface{}) (string, interface{}, error) {
	executor, exists := s.deps.ExecutorRegistry.GetExecutor(actionType)
	if !exists {
		return "error", nil, fmt.Errorf("no executor registered for action type: %s", actionType)
	}

	// Build config map — the executor expects { "<type>": { ...config } }
	typeDef, _ := action[actionType].(map[string]interface{})
	if typeDef == nil {
		return "error", nil, fmt.Errorf("%s action definition is missing", actionType)
	}

	config := map[string]interface{}{
		actionType: typeDef,
	}

	logger.Debug(fmt.Sprintf("[SERVICE:Event] Plugin %s: delegating to executor", actionType))

	if err := executor.Execute(context.Background(), config); err != nil {
		return "error", nil, fmt.Errorf("%s execution failed: %w", actionType, err)
	}

	// Executors don't return response data (fire-and-forget), return success
	return "success", map[string]interface{}{"status": "dispatched"}, nil
}

// executeTriggerEntity executes a registered trigger entity by ID.
// Returns the output handle ("out"), output data, and any error.
func (s *EventService) executeTriggerEntity(data map[string]interface{}) (string, interface{}, error) {
	// TODO: Implement — fetch trigger config by triggerId, resolve, execute, return result
	// For now, return error indicating not yet implemented
	triggerId, _ := data["triggerId"].(string)
	return "error", nil, fmt.Errorf("trigger entity execution not yet implemented (triggerId=%s)", triggerId)
}
