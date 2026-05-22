package slack

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

/**
 * SlackExecutor Tests
 */

func TestSlackExecutor_GetType(t *testing.T) {
	executor := NewSlackExecutor()

	if executor.GetType() != "slack" {
		t.Errorf("GetType() = %q, want 'slack'", executor.GetType())
	}
}

func TestSlackExecutor_Execute_MissingWebhookUrl(t *testing.T) {
	executor := NewSlackExecutor()

	config := map[string]interface{}{
		"text": "Test message",
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for missing 'webhookUrl', got nil")
	}
}

func TestSlackExecutor_Execute_EmptyWebhookUrl(t *testing.T) {
	executor := NewSlackExecutor()

	config := map[string]interface{}{
		"webhookUrl": "",
		"text":       "Test message",
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for empty 'webhookUrl', got nil")
	}
}

func TestSlackExecutor_Execute_MissingText(t *testing.T) {
	executor := NewSlackExecutor()

	config := map[string]interface{}{
		"webhookUrl": "https://hooks.slack.com/test",
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for missing 'text', got nil")
	}
}

func TestSlackExecutor_Execute_Success(t *testing.T) {
	var receivedPayload map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &receivedPayload)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer server.Close()

	executor := NewSlackExecutor()

	config := map[string]interface{}{
		"webhookUrl": server.URL,
		"text":       "Hello from test!",
	}

	err := executor.Execute(context.Background(), config)

	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	if receivedPayload["text"] != "Hello from test!" {
		t.Errorf("Expected text = 'Hello from test!', got %v", receivedPayload["text"])
	}
}

func TestSlackExecutor_Execute_WithUsername(t *testing.T) {
	var receivedPayload map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &receivedPayload)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer server.Close()

	executor := NewSlackExecutor()

	config := map[string]interface{}{
		"webhookUrl": server.URL,
		"text":       "Test message",
		"username":   "Alert Bot",
	}

	err := executor.Execute(context.Background(), config)

	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	if receivedPayload["username"] != "Alert Bot" {
		t.Errorf("Expected username = 'Alert Bot', got %v", receivedPayload["username"])
	}
}

func TestSlackExecutor_Execute_WithIconEmoji(t *testing.T) {
	var receivedPayload map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &receivedPayload)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer server.Close()

	executor := NewSlackExecutor()

	config := map[string]interface{}{
		"webhookUrl": server.URL,
		"text":       "Test message",
		"icon_emoji": ":warning:",
	}

	err := executor.Execute(context.Background(), config)

	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	if receivedPayload["icon_emoji"] != ":warning:" {
		t.Errorf("Expected icon_emoji = ':warning:', got %v", receivedPayload["icon_emoji"])
	}
}

func TestSlackExecutor_Execute_WithChannel(t *testing.T) {
	var receivedPayload map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &receivedPayload)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer server.Close()

	executor := NewSlackExecutor()

	config := map[string]interface{}{
		"webhookUrl": server.URL,
		"text":       "Test message",
		"channel":    "#alerts",
	}

	err := executor.Execute(context.Background(), config)

	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	if receivedPayload["channel"] != "#alerts" {
		t.Errorf("Expected channel = '#alerts', got %v", receivedPayload["channel"])
	}
}

func TestSlackExecutor_Execute_WithBlocks(t *testing.T) {
	var receivedPayload map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &receivedPayload)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer server.Close()

	executor := NewSlackExecutor()

	blocks := []interface{}{
		map[string]interface{}{
			"type": "section",
			"text": map[string]interface{}{
				"type": "mrkdwn",
				"text": "*Alert*: Server is down",
			},
		},
	}

	config := map[string]interface{}{
		"webhookUrl": server.URL,
		"text":       "Fallback text",
		"blocks":     blocks,
	}

	err := executor.Execute(context.Background(), config)

	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	if receivedPayload["blocks"] == nil {
		t.Error("Expected blocks to be included in payload")
	}
}

func TestSlackExecutor_Execute_WithAttachments(t *testing.T) {
	var receivedPayload map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &receivedPayload)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer server.Close()

	executor := NewSlackExecutor()

	attachments := []interface{}{
		map[string]interface{}{
			"color":  "danger",
			"title":  "Error",
			"text":   "Something went wrong",
			"fields": []interface{}{},
		},
	}

	config := map[string]interface{}{
		"webhookUrl":  server.URL,
		"text":        "Alert",
		"attachments": attachments,
	}

	err := executor.Execute(context.Background(), config)

	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	if receivedPayload["attachments"] == nil {
		t.Error("Expected attachments to be included in payload")
	}
}

func TestSlackExecutor_Execute_ContentType(t *testing.T) {
	var receivedContentType string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedContentType = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer server.Close()

	executor := NewSlackExecutor()

	config := map[string]interface{}{
		"webhookUrl": server.URL,
		"text":       "Test",
	}

	executor.Execute(context.Background(), config)

	if receivedContentType != "application/json" {
		t.Errorf("Expected Content-Type = 'application/json', got %q", receivedContentType)
	}
}

/**
 * Error Response Tests
 */

func TestSlackExecutor_Execute_NotOkResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid_payload"))
	}))
	defer server.Close()

	executor := NewSlackExecutor()

	config := map[string]interface{}{
		"webhookUrl": server.URL,
		"text":       "Test",
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for non-'ok' response, got nil")
	}
}

func TestSlackExecutor_Execute_Non200Status(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid_payload"))
	}))
	defer server.Close()

	executor := NewSlackExecutor()

	config := map[string]interface{}{
		"webhookUrl": server.URL,
		"text":       "Test",
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for non-200 status, got nil")
	}
}

func TestSlackExecutor_Execute_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	executor := NewSlackExecutor()

	config := map[string]interface{}{
		"webhookUrl": server.URL,
		"text":       "Test",
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for 5xx status, got nil")
	}
}

func TestSlackExecutor_Execute_InvalidUrl(t *testing.T) {
	executor := NewSlackExecutor()

	config := map[string]interface{}{
		"webhookUrl": "http://invalid.nonexistent.host.test/webhook",
		"text":       "Test",
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for invalid URL, got nil")
	}
}

func TestSlackExecutor_Execute_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			return
		}
	}))
	defer server.Close()

	executor := NewSlackExecutor()

	config := map[string]interface{}{
		"webhookUrl": server.URL,
		"text":       "Test",
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := executor.Execute(ctx, config)

	if err == nil {
		t.Fatal("Execute() expected error for cancelled context, got nil")
	}
}

/**
 * Full Config Tests
 */

func TestSlackExecutor_Execute_FullConfig(t *testing.T) {
	var receivedPayload map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &receivedPayload)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer server.Close()

	executor := NewSlackExecutor()

	config := map[string]interface{}{
		"webhookUrl": server.URL,
		"text":       "Critical: Server down",
		"username":   "Alert Bot",
		"icon_emoji": ":rotating_light:",
		"channel":    "#incidents",
		"blocks": []interface{}{
			map[string]interface{}{
				"type": "section",
				"text": map[string]interface{}{
					"type": "mrkdwn",
					"text": "*Server Alert*\nproduction-web-01 is not responding",
				},
			},
		},
	}

	err := executor.Execute(context.Background(), config)

	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	// Verify all fields
	if receivedPayload["text"] != "Critical: Server down" {
		t.Error("Missing or incorrect text")
	}
	if receivedPayload["username"] != "Alert Bot" {
		t.Error("Missing or incorrect username")
	}
	if receivedPayload["icon_emoji"] != ":rotating_light:" {
		t.Error("Missing or incorrect icon_emoji")
	}
	if receivedPayload["channel"] != "#incidents" {
		t.Error("Missing or incorrect channel")
	}
	if receivedPayload["blocks"] == nil {
		t.Error("Missing blocks")
	}
}

/**
 * Slack Response Tests
 */

func TestSlackExecutor_Execute_SlackErrorResponses(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		shouldPass bool
	}{
		{"Success - ok", 200, "ok", true},
		{"Error - invalid_payload", 200, "invalid_payload", false},
		{"Error - channel_not_found", 200, "channel_not_found", false},
		{"Error - action_prohibited", 200, "action_prohibited", false},
		{"Error - 400 Bad Request", 400, "Bad Request", false},
		{"Error - 404 Not Found", 404, "Webhook not found", false},
		{"Error - 500 Server Error", 500, "Internal error", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.body))
			}))
			defer server.Close()

			executor := NewSlackExecutor()

			config := map[string]interface{}{
				"webhookUrl": server.URL,
				"text":       "Test",
			}

			err := executor.Execute(context.Background(), config)

			if tt.shouldPass && err != nil {
				t.Errorf("Execute() unexpected error: %v", err)
			}
			if !tt.shouldPass && err == nil {
				t.Error("Execute() expected error, got nil")
			}
		})
	}
}
