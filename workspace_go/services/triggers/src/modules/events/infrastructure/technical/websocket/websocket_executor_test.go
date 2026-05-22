package websocket

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

/**
 * WebSocketExecutor Tests
 */

func TestWebSocketExecutor_GetType(t *testing.T) {
	executor := NewWebSocketExecutor()

	if executor.GetType() != "websocket" {
		t.Errorf("GetType() = %q, want 'websocket'", executor.GetType())
	}
}

func TestWebSocketExecutor_Execute_MissingWebsocketField(t *testing.T) {
	executor := NewWebSocketExecutor()

	config := map[string]interface{}{
		"notwebsocket": map[string]interface{}{
			"url": "ws://localhost:8080/ws",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for missing 'websocket' field, got nil")
	}
}

func TestWebSocketExecutor_Execute_MissingUrl(t *testing.T) {
	executor := NewWebSocketExecutor()

	config := map[string]interface{}{
		"websocket": map[string]interface{}{
			"message": "test",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for missing 'url', got nil")
	}
}

func TestWebSocketExecutor_Execute_EmptyUrl(t *testing.T) {
	executor := NewWebSocketExecutor()

	config := map[string]interface{}{
		"websocket": map[string]interface{}{
			"url":     "",
			"message": "test",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for empty 'url', got nil")
	}
}

func TestWebSocketExecutor_Execute_ContextCancellation(t *testing.T) {
	executor := NewWebSocketExecutor()

	config := map[string]interface{}{
		"websocket": map[string]interface{}{
			"url":     "ws://invalid.nonexistent.host.test:8080/ws",
			"message": "test",
		},
	}

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := executor.Execute(ctx, config)

	if err == nil {
		t.Fatal("Execute() expected error for cancelled context, got nil")
	}
}

/**
 * Integration Tests with Test Server
 */

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func TestWebSocketExecutor_Execute_Success_TextMessage(t *testing.T) {
	var receivedMessage []byte
	var receivedMsgType int
	done := make(chan struct{})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Logf("Failed to upgrade: %v", err)
			return
		}
		defer conn.Close()

		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			close(done)
			return
		}
		receivedMsgType = msgType
		receivedMessage = msg
		close(done)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	executor := NewWebSocketExecutor()

	config := map[string]interface{}{
		"websocket": map[string]interface{}{
			"url":     wsURL,
			"message": "Hello WebSocket",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	// Wait for server to receive the message
	<-done

	if receivedMsgType != websocket.TextMessage {
		t.Errorf("Expected text message, got type %d", receivedMsgType)
	}

	if string(receivedMessage) != "Hello WebSocket" {
		t.Errorf("Received message = %q, want 'Hello WebSocket'", string(receivedMessage))
	}
}

func TestWebSocketExecutor_Execute_Success_BinaryMessage(t *testing.T) {
	var receivedMsgType int

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		msgType, _, err := conn.ReadMessage()
		if err != nil {
			return
		}
		receivedMsgType = msgType
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	executor := NewWebSocketExecutor()

	config := map[string]interface{}{
		"websocket": map[string]interface{}{
			"url":         wsURL,
			"message":     "binary data",
			"messageType": "binary",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	if receivedMsgType != websocket.BinaryMessage {
		t.Errorf("Expected binary message, got type %d", receivedMsgType)
	}
}

func TestWebSocketExecutor_Execute_Success_WithHeaders(t *testing.T) {
	var receivedAuth string
	var receivedCustom string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		receivedCustom = r.Header.Get("X-Custom-Header")

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		conn.ReadMessage()
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	executor := NewWebSocketExecutor()

	config := map[string]interface{}{
		"websocket": map[string]interface{}{
			"url":     wsURL,
			"message": "test",
			"headers": map[string]interface{}{
				"Authorization":   "Bearer token123",
				"X-Custom-Header": "custom-value",
			},
		},
	}

	err := executor.Execute(context.Background(), config)

	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	if receivedAuth != "Bearer token123" {
		t.Errorf("Authorization header = %q, want 'Bearer token123'", receivedAuth)
	}

	if receivedCustom != "custom-value" {
		t.Errorf("X-Custom-Header = %q, want 'custom-value'", receivedCustom)
	}
}

func TestWebSocketExecutor_Execute_Success_JSONMessage(t *testing.T) {
	var receivedMessage []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		receivedMessage = msg
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	executor := NewWebSocketExecutor()

	config := map[string]interface{}{
		"websocket": map[string]interface{}{
			"url": wsURL,
			"message": map[string]interface{}{
				"type":   "event",
				"action": "user.created",
				"data": map[string]interface{}{
					"userId": "123",
				},
			},
		},
	}

	err := executor.Execute(context.Background(), config)

	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	if len(receivedMessage) == 0 {
		t.Error("Expected JSON message, got empty")
	}

	// Should contain expected fields (basic check)
	msgStr := string(receivedMessage)
	if !strings.Contains(msgStr, "user.created") {
		t.Errorf("Message should contain 'user.created': %s", msgStr)
	}
}

/**
 * Message Extraction Tests
 */

func TestWebSocketExecutor_extractMessage_NoMessage(t *testing.T) {
	executor := &WebSocketExecutor{}

	config := map[string]interface{}{
		"url": "ws://localhost:8080/ws",
	}

	msg, err := executor.extractMessage(config)

	if err != nil {
		t.Fatalf("extractMessage() unexpected error: %v", err)
	}

	if string(msg) != "" {
		t.Errorf("extractMessage() = %q, want empty string", string(msg))
	}
}

func TestWebSocketExecutor_extractMessage_StringMessage(t *testing.T) {
	executor := &WebSocketExecutor{}

	config := map[string]interface{}{
		"message": "Hello WebSocket",
	}

	msg, err := executor.extractMessage(config)

	if err != nil {
		t.Fatalf("extractMessage() unexpected error: %v", err)
	}

	if string(msg) != "Hello WebSocket" {
		t.Errorf("extractMessage() = %q, want 'Hello WebSocket'", string(msg))
	}
}

func TestWebSocketExecutor_extractMessage_MapMessage(t *testing.T) {
	executor := &WebSocketExecutor{}

	config := map[string]interface{}{
		"message": map[string]interface{}{
			"type": "event",
			"data": "value",
		},
	}

	msg, err := executor.extractMessage(config)

	if err != nil {
		t.Fatalf("extractMessage() unexpected error: %v", err)
	}

	if len(msg) == 0 {
		t.Error("extractMessage() returned empty for map message")
	}
}

func TestWebSocketExecutor_extractMessage_ArrayMessage(t *testing.T) {
	executor := &WebSocketExecutor{}

	config := map[string]interface{}{
		"message": []interface{}{"item1", "item2"},
	}

	msg, err := executor.extractMessage(config)

	if err != nil {
		t.Fatalf("extractMessage() unexpected error: %v", err)
	}

	if len(msg) == 0 {
		t.Error("extractMessage() returned empty for array message")
	}
}

/**
 * Connection Error Tests
 */

func TestWebSocketExecutor_Execute_InvalidUrl(t *testing.T) {
	executor := NewWebSocketExecutor()

	config := map[string]interface{}{
		"websocket": map[string]interface{}{
			"url":     "ws://invalid.nonexistent.host.test:8080/ws",
			"message": "test",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected connection error, got nil")
	}
}

func TestWebSocketExecutor_Execute_WSSUrl(t *testing.T) {
	executor := NewWebSocketExecutor()

	// wss:// URL to a non-existent host
	config := map[string]interface{}{
		"websocket": map[string]interface{}{
			"url":     "wss://invalid.nonexistent.host.test:443/ws",
			"message": "test",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected connection error, got nil")
	}
}

/**
 * Full Config Tests
 */

func TestWebSocketExecutor_Execute_FullConfig(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		conn.ReadMessage()
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	executor := NewWebSocketExecutor()

	config := map[string]interface{}{
		"websocket": map[string]interface{}{
			"url": wsURL,
			"headers": map[string]interface{}{
				"Authorization": "Bearer mytoken",
				"X-Request-ID":  "req-123",
			},
			"messageType": "text",
			"message": map[string]interface{}{
				"event":     "notification",
				"channel":   "user.123",
				"data":      "Hello!",
				"timestamp": "2024-01-01T00:00:00Z",
			},
		},
	}

	err := executor.Execute(context.Background(), config)

	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}
}

/**
 * Default Values Tests
 */

func TestWebSocketExecutor_Execute_DefaultMessageType(t *testing.T) {
	var receivedMsgType int

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		msgType, _, _ := conn.ReadMessage()
		receivedMsgType = msgType
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	executor := NewWebSocketExecutor()

	// No messageType specified - should default to text
	config := map[string]interface{}{
		"websocket": map[string]interface{}{
			"url":     wsURL,
			"message": "test",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	if receivedMsgType != websocket.TextMessage {
		t.Errorf("Default message type = %d, want TextMessage (%d)", receivedMsgType, websocket.TextMessage)
	}
}

func TestWebSocketExecutor_Execute_EmptyMessage(t *testing.T) {
	var receivedMessage []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		_, msg, _ := conn.ReadMessage()
		receivedMessage = msg
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	executor := NewWebSocketExecutor()

	// No message specified
	config := map[string]interface{}{
		"websocket": map[string]interface{}{
			"url": wsURL,
		},
	}

	err := executor.Execute(context.Background(), config)

	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	if string(receivedMessage) != "" {
		t.Errorf("Expected empty message, got %q", string(receivedMessage))
	}
}
