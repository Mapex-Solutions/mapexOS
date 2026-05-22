package http

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

/**
 * HTTPExecutor Tests
 */

func TestHTTPExecutor_GetType(t *testing.T) {
	executor := NewHTTPExecutor()

	if executor.GetType() != "http" {
		t.Errorf("GetType() = %q, want 'http'", executor.GetType())
	}
}

func TestHTTPExecutor_Execute_Success_POST(t *testing.T) {
	// Create a test server
	var receivedMethod string
	var receivedBody []byte
	var receivedContentType string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedContentType = r.Header.Get("Content-Type")
		receivedBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	executor := NewHTTPExecutor()

	config := map[string]interface{}{
		"http": map[string]interface{}{
			"endpoint": server.URL,
			"method":   "POST",
			"body": map[string]interface{}{
				"message": "test",
			},
		},
	}

	err := executor.Execute(context.Background(), config)

	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	if receivedMethod != "POST" {
		t.Errorf("Received method = %q, want 'POST'", receivedMethod)
	}

	if receivedContentType != "application/json" {
		t.Errorf("Content-Type = %q, want 'application/json'", receivedContentType)
	}

	if len(receivedBody) == 0 {
		t.Error("Expected request body, got empty")
	}
}

func TestHTTPExecutor_Execute_Success_GET(t *testing.T) {
	var receivedMethod string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	executor := NewHTTPExecutor()

	config := map[string]interface{}{
		"http": map[string]interface{}{
			"endpoint": server.URL,
			"method":   "GET",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	if receivedMethod != "GET" {
		t.Errorf("Received method = %q, want 'GET'", receivedMethod)
	}
}

func TestHTTPExecutor_Execute_Success_WithHeaders(t *testing.T) {
	var receivedAuth string
	var receivedCustomHeader string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		receivedCustomHeader = r.Header.Get("X-Custom-Header")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	executor := NewHTTPExecutor()

	config := map[string]interface{}{
		"http": map[string]interface{}{
			"endpoint": server.URL,
			"method":   "POST",
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

	if receivedCustomHeader != "custom-value" {
		t.Errorf("X-Custom-Header = %q, want 'custom-value'", receivedCustomHeader)
	}
}

func TestHTTPExecutor_Execute_DefaultMethod(t *testing.T) {
	var receivedMethod string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	executor := NewHTTPExecutor()

	// No method specified - should default to POST
	config := map[string]interface{}{
		"http": map[string]interface{}{
			"endpoint": server.URL,
		},
	}

	err := executor.Execute(context.Background(), config)

	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	if receivedMethod != "POST" {
		t.Errorf("Default method = %q, want 'POST'", receivedMethod)
	}
}

func TestHTTPExecutor_Execute_MissingHttpField(t *testing.T) {
	executor := NewHTTPExecutor()

	config := map[string]interface{}{
		"nothttp": map[string]interface{}{
			"endpoint": "https://example.com",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for missing 'http' field, got nil")
	}
}

func TestHTTPExecutor_Execute_MissingEndpoint(t *testing.T) {
	executor := NewHTTPExecutor()

	config := map[string]interface{}{
		"http": map[string]interface{}{
			"method": "POST",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for missing 'endpoint', got nil")
	}
}

func TestHTTPExecutor_Execute_Non2xxStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "bad request"}`))
	}))
	defer server.Close()

	executor := NewHTTPExecutor()

	config := map[string]interface{}{
		"http": map[string]interface{}{
			"endpoint": server.URL,
			"method":   "POST",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for non-2xx status, got nil")
	}
}

func TestHTTPExecutor_Execute_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	executor := NewHTTPExecutor()

	config := map[string]interface{}{
		"http": map[string]interface{}{
			"endpoint": server.URL,
			"method":   "POST",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for 5xx status, got nil")
	}
}

func TestHTTPExecutor_Execute_InvalidEndpoint(t *testing.T) {
	executor := NewHTTPExecutor()

	config := map[string]interface{}{
		"http": map[string]interface{}{
			"endpoint": "not-a-valid-url",
			"method":   "GET",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for invalid endpoint, got nil")
	}
}

func TestHTTPExecutor_Execute_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Slow response - should be cancelled
		select {
		case <-r.Context().Done():
			return
		}
	}))
	defer server.Close()

	executor := NewHTTPExecutor()

	config := map[string]interface{}{
		"http": map[string]interface{}{
			"endpoint": server.URL,
			"method":   "POST",
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
 * Status Code Tests
 */

func TestHTTPExecutor_Execute_StatusCodes(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		shouldPass bool
	}{
		{"200 OK", 200, true},
		{"201 Created", 201, true},
		{"204 No Content", 204, true},
		{"299 Custom 2xx", 299, true},
		{"301 Redirect", 301, false}, // Redirects without following
		{"400 Bad Request", 400, false},
		{"401 Unauthorized", 401, false},
		{"403 Forbidden", 403, false},
		{"404 Not Found", 404, false},
		{"500 Internal Server Error", 500, false},
		{"502 Bad Gateway", 502, false},
		{"503 Service Unavailable", 503, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			executor := NewHTTPExecutor()

			config := map[string]interface{}{
				"http": map[string]interface{}{
					"endpoint": server.URL,
					"method":   "GET",
				},
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
 * HTTP Methods Tests
 */

func TestHTTPExecutor_Execute_AllMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			var receivedMethod string

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedMethod = r.Method
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			executor := NewHTTPExecutor()

			config := map[string]interface{}{
				"http": map[string]interface{}{
					"endpoint": server.URL,
					"method":   method,
				},
			}

			err := executor.Execute(context.Background(), config)

			if err != nil {
				t.Fatalf("Execute() with method %s unexpected error: %v", method, err)
			}

			if receivedMethod != method {
				t.Errorf("Received method = %q, want %q", receivedMethod, method)
			}
		})
	}
}
