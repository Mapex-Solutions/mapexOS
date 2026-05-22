package teams

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

/**
 * TeamsExecutor Tests
 */

func TestTeamsExecutor_GetType(t *testing.T) {
	executor := NewTeamsExecutor()

	if executor.GetType() != "teams" {
		t.Errorf("GetType() = %q, want 'teams'", executor.GetType())
	}
}

func TestTeamsExecutor_Execute_MissingWebhookUrl(t *testing.T) {
	executor := NewTeamsExecutor()

	config := map[string]interface{}{
		"title": "Test Title",
		"text":  "Test Text",
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for missing 'webhookUrl', got nil")
	}
}

func TestTeamsExecutor_Execute_EmptyWebhookUrl(t *testing.T) {
	executor := NewTeamsExecutor()

	config := map[string]interface{}{
		"webhookUrl": "",
		"title":      "Test Title",
		"text":       "Test Text",
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for empty 'webhookUrl', got nil")
	}
}

func TestTeamsExecutor_Execute_Success(t *testing.T) {
	var receivedPayload map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &receivedPayload)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	executor := NewTeamsExecutor()

	config := map[string]interface{}{
		"webhookUrl": server.URL,
		"title":      "Test Alert",
		"text":       "This is a test message",
	}

	err := executor.Execute(context.Background(), config)

	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	// Verify MessageCard format
	if receivedPayload["@type"] != "MessageCard" {
		t.Errorf("Expected @type = 'MessageCard', got %v", receivedPayload["@type"])
	}
	if receivedPayload["@context"] != "https://schema.org/extensions" {
		t.Errorf("Expected @context to be schema.org extensions")
	}
	if receivedPayload["title"] != "Test Alert" {
		t.Errorf("Expected title = 'Test Alert', got %v", receivedPayload["title"])
	}
	if receivedPayload["text"] != "This is a test message" {
		t.Errorf("Expected text = 'This is a test message', got %v", receivedPayload["text"])
	}
}

func TestTeamsExecutor_Execute_WithThemeColor(t *testing.T) {
	var receivedPayload map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &receivedPayload)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	executor := NewTeamsExecutor()

	config := map[string]interface{}{
		"webhookUrl": server.URL,
		"title":      "Alert",
		"text":       "Error occurred",
		"themeColor": "FF0000",
	}

	err := executor.Execute(context.Background(), config)

	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	if receivedPayload["themeColor"] != "FF0000" {
		t.Errorf("Expected themeColor = 'FF0000', got %v", receivedPayload["themeColor"])
	}
}

func TestTeamsExecutor_Execute_WithSections(t *testing.T) {
	var receivedPayload map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &receivedPayload)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	executor := NewTeamsExecutor()

	sections := []interface{}{
		map[string]interface{}{
			"activityTitle": "Details",
			"facts": []interface{}{
				map[string]interface{}{"name": "Server", "value": "prod-01"},
				map[string]interface{}{"name": "Status", "value": "Critical"},
			},
		},
	}

	config := map[string]interface{}{
		"webhookUrl": server.URL,
		"title":      "Alert",
		"text":       "Server alert",
		"sections":   sections,
	}

	err := executor.Execute(context.Background(), config)

	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	if receivedPayload["sections"] == nil {
		t.Error("Expected sections to be included in payload")
	}
}

func TestTeamsExecutor_Execute_ContentType(t *testing.T) {
	var receivedContentType string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedContentType = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	executor := NewTeamsExecutor()

	config := map[string]interface{}{
		"webhookUrl": server.URL,
		"title":      "Test",
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

func TestTeamsExecutor_Execute_Non2xxStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	}))
	defer server.Close()

	executor := NewTeamsExecutor()

	config := map[string]interface{}{
		"webhookUrl": server.URL,
		"title":      "Test",
		"text":       "Test",
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for non-2xx status, got nil")
	}
}

func TestTeamsExecutor_Execute_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	executor := NewTeamsExecutor()

	config := map[string]interface{}{
		"webhookUrl": server.URL,
		"title":      "Test",
		"text":       "Test",
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for 5xx status, got nil")
	}
}

func TestTeamsExecutor_Execute_InvalidUrl(t *testing.T) {
	executor := NewTeamsExecutor()

	config := map[string]interface{}{
		"webhookUrl": "http://invalid.nonexistent.host.test/webhook",
		"title":      "Test",
		"text":       "Test",
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for invalid URL, got nil")
	}
}

func TestTeamsExecutor_Execute_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Slow response
		select {
		case <-r.Context().Done():
			return
		}
	}))
	defer server.Close()

	executor := NewTeamsExecutor()

	config := map[string]interface{}{
		"webhookUrl": server.URL,
		"title":      "Test",
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
 * Status Code Tests
 */

func TestTeamsExecutor_Execute_StatusCodes(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		shouldPass bool
	}{
		{"200 OK", 200, true},
		{"201 Created", 201, true},
		{"202 Accepted", 202, true},
		{"204 No Content", 204, true},
		{"400 Bad Request", 400, false},
		{"401 Unauthorized", 401, false},
		{"403 Forbidden", 403, false},
		{"404 Not Found", 404, false},
		{"500 Internal Server Error", 500, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			executor := NewTeamsExecutor()

			config := map[string]interface{}{
				"webhookUrl": server.URL,
				"title":      "Test",
				"text":       "Test",
			}

			err := executor.Execute(context.Background(), config)

			if tt.shouldPass && err != nil {
				t.Errorf("Execute() with status %d unexpected error: %v", tt.statusCode, err)
			}
			if !tt.shouldPass && err == nil {
				t.Errorf("Execute() with status %d expected error, got nil", tt.statusCode)
			}
		})
	}
}

/**
 * Full Config Tests
 */

func TestTeamsExecutor_Execute_FullConfig(t *testing.T) {
	var receivedPayload map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &receivedPayload)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	executor := NewTeamsExecutor()

	config := map[string]interface{}{
		"webhookUrl": server.URL,
		"title":      "Critical Alert",
		"text":       "Server CPU usage exceeded 90%",
		"themeColor": "FF0000",
		"sections": []interface{}{
			map[string]interface{}{
				"activityTitle": "Server Details",
				"facts": []interface{}{
					map[string]interface{}{"name": "Server", "value": "prod-web-01"},
					map[string]interface{}{"name": "CPU", "value": "92%"},
					map[string]interface{}{"name": "Memory", "value": "78%"},
				},
			},
		},
	}

	err := executor.Execute(context.Background(), config)

	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	// Verify all fields
	if receivedPayload["@type"] != "MessageCard" {
		t.Error("Missing MessageCard @type")
	}
	if receivedPayload["title"] != "Critical Alert" {
		t.Error("Missing or incorrect title")
	}
	if receivedPayload["text"] != "Server CPU usage exceeded 90%" {
		t.Error("Missing or incorrect text")
	}
	if receivedPayload["themeColor"] != "FF0000" {
		t.Error("Missing or incorrect themeColor")
	}
	if receivedPayload["sections"] == nil {
		t.Error("Missing sections")
	}
}

/**
 * Optional Fields Tests
 */

func TestTeamsExecutor_Execute_TitleOnly(t *testing.T) {
	var receivedPayload map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &receivedPayload)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	executor := NewTeamsExecutor()

	config := map[string]interface{}{
		"webhookUrl": server.URL,
		"title":      "Just a title",
	}

	err := executor.Execute(context.Background(), config)

	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	if receivedPayload["title"] != "Just a title" {
		t.Error("Title should be included")
	}
	if receivedPayload["text"] != nil {
		t.Error("Text should not be included when not provided")
	}
}

func TestTeamsExecutor_Execute_TextOnly(t *testing.T) {
	var receivedPayload map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &receivedPayload)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	executor := NewTeamsExecutor()

	config := map[string]interface{}{
		"webhookUrl": server.URL,
		"text":       "Just some text",
	}

	err := executor.Execute(context.Background(), config)

	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	if receivedPayload["text"] != "Just some text" {
		t.Error("Text should be included")
	}
}
