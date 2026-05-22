package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	mocknats "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mock_servers/nats"
)

/**
 * NATSExecutor Tests
 */

func TestNATSExecutor_GetType(t *testing.T) {
	executor := NewNATSExecutor()

	if executor.GetType() != "nats" {
		t.Errorf("GetType() = %q, want 'nats'", executor.GetType())
	}
}

func TestNATSExecutor_Execute_MissingNatsField(t *testing.T) {
	executor := NewNATSExecutor()

	config := map[string]interface{}{
		"notnats": map[string]interface{}{
			"url": "nats://localhost:4222",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for missing 'nats' field, got nil")
	}
}

func TestNATSExecutor_Execute_MissingUrl(t *testing.T) {
	executor := NewNATSExecutor()

	config := map[string]interface{}{
		"nats": map[string]interface{}{
			"subject": "test.subject",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for missing 'url', got nil")
	}
}

func TestNATSExecutor_Execute_EmptyUrl(t *testing.T) {
	executor := NewNATSExecutor()

	config := map[string]interface{}{
		"nats": map[string]interface{}{
			"url":     "",
			"subject": "test.subject",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for empty 'url', got nil")
	}
}

func TestNATSExecutor_Execute_MissingSubject(t *testing.T) {
	executor := NewNATSExecutor()

	config := map[string]interface{}{
		"nats": map[string]interface{}{
			"url": "nats://localhost:4222",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for missing 'subject', got nil")
	}
}

func TestNATSExecutor_Execute_EmptySubject(t *testing.T) {
	executor := NewNATSExecutor()

	config := map[string]interface{}{
		"nats": map[string]interface{}{
			"url":     "nats://localhost:4222",
			"subject": "",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for empty 'subject', got nil")
	}
}

func TestNATSExecutor_Execute_ContextCancellation(t *testing.T) {
	executor := NewNATSExecutor()

	config := map[string]interface{}{
		"nats": map[string]interface{}{
			"url":     "nats://invalid.nonexistent.host.test:4222",
			"subject": "test.subject",
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
 * Config Extraction Tests
 */

func TestNATSExecutor_Execute_WithUserCredentials(t *testing.T) {
	executor := NewNATSExecutor()

	config := map[string]interface{}{
		"nats": map[string]interface{}{
			"url":      "nats://invalid.nonexistent.host.test:4222",
			"subject":  "test.subject",
			"username": "testuser",
			"password": "testpass",
		},
	}

	err := executor.Execute(context.Background(), config)

	// Should fail on connection, not config extraction
	if err == nil {
		t.Fatal("Expected connection error, got nil")
	}
}

func TestNATSExecutor_Execute_WithToken(t *testing.T) {
	executor := NewNATSExecutor()

	config := map[string]interface{}{
		"nats": map[string]interface{}{
			"url":     "nats://invalid.nonexistent.host.test:4222",
			"subject": "test.subject",
			"token":   "myauthtoken",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Expected connection error, got nil")
	}
}

func TestNATSExecutor_Execute_WithCredsFile(t *testing.T) {
	executor := NewNATSExecutor()

	config := map[string]interface{}{
		"nats": map[string]interface{}{
			"url":       "nats://invalid.nonexistent.host.test:4222",
			"subject":   "test.subject",
			"credsFile": "/nonexistent/path/to/creds",
		},
	}

	err := executor.Execute(context.Background(), config)

	// Should fail (either on creds file or connection)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestNATSExecutor_Execute_WithHeaders(t *testing.T) {
	executor := NewNATSExecutor()

	config := map[string]interface{}{
		"nats": map[string]interface{}{
			"url":     "nats://invalid.nonexistent.host.test:4222",
			"subject": "test.subject",
			"headers": map[string]interface{}{
				"X-Custom-Header": "custom-value",
				"X-Request-ID":    "req-123",
			},
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Expected connection error, got nil")
	}
}

/**
 * Message Extraction Tests
 */

func TestNATSExecutor_extractMessage_NoMessage(t *testing.T) {
	executor := &NATSExecutor{}

	config := map[string]interface{}{
		"url":     "nats://localhost:4222",
		"subject": "test.subject",
	}

	msg, err := executor.extractMessage(config)

	if err != nil {
		t.Fatalf("extractMessage() unexpected error: %v", err)
	}

	if string(msg) != "" {
		t.Errorf("extractMessage() = %q, want empty string", string(msg))
	}
}

func TestNATSExecutor_extractMessage_StringMessage(t *testing.T) {
	executor := &NATSExecutor{}

	config := map[string]interface{}{
		"message": "Hello NATS",
	}

	msg, err := executor.extractMessage(config)

	if err != nil {
		t.Fatalf("extractMessage() unexpected error: %v", err)
	}

	if string(msg) != "Hello NATS" {
		t.Errorf("extractMessage() = %q, want 'Hello NATS'", string(msg))
	}
}

func TestNATSExecutor_extractMessage_MapMessage(t *testing.T) {
	executor := &NATSExecutor{}

	config := map[string]interface{}{
		"message": map[string]interface{}{
			"event": "user.created",
			"data": map[string]interface{}{
				"userId": "123",
			},
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

func TestNATSExecutor_extractMessage_ArrayMessage(t *testing.T) {
	executor := &NATSExecutor{}

	config := map[string]interface{}{
		"message": []interface{}{"item1", "item2", "item3"},
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
 * Full Config Tests
 */

func TestNATSExecutor_Execute_FullConfig(t *testing.T) {
	executor := NewNATSExecutor()

	config := map[string]interface{}{
		"nats": map[string]interface{}{
			"url":      "nats://nats.example.com:4222",
			"subject":  "events.user.created",
			"username": "admin",
			"password": "secret",
			"headers": map[string]interface{}{
				"X-Source":     "triggers-service",
				"X-Request-ID": "req-456",
			},
			"message": map[string]interface{}{
				"event":     "user.created",
				"userId":    "user-123",
				"timestamp": "2024-01-01T00:00:00Z",
			},
		},
	}

	err := executor.Execute(context.Background(), config)

	// Should fail on connection, not config
	if err == nil {
		t.Fatal("Expected connection error, got nil")
	}
}

/**
 * URL Format Tests
 */

func TestNATSExecutor_Execute_MultipleUrls(t *testing.T) {
	executor := NewNATSExecutor()

	// NATS supports comma-separated URLs for cluster
	config := map[string]interface{}{
		"nats": map[string]interface{}{
			"url":     "nats://invalid1.test:4222,nats://invalid2.test:4222",
			"subject": "test.subject",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Expected connection error, got nil")
	}
}

func TestNATSExecutor_Execute_TLSUrl(t *testing.T) {
	executor := NewNATSExecutor()

	config := map[string]interface{}{
		"nats": map[string]interface{}{
			"url":     "tls://invalid.nonexistent.host.test:4222",
			"subject": "test.subject",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Expected connection error, got nil")
	}
}

/**
 * Subject Pattern Tests
 */

func TestNATSExecutor_Execute_HierarchicalSubject(t *testing.T) {
	executor := NewNATSExecutor()

	config := map[string]interface{}{
		"nats": map[string]interface{}{
			"url":     "nats://invalid.nonexistent.host.test:4222",
			"subject": "events.user.profile.updated",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Expected connection error, got nil")
	}
}

func TestNATSExecutor_Execute_WildcardSubject(t *testing.T) {
	executor := NewNATSExecutor()

	// Note: Wildcards are typically for subscriptions, not publishing
	// But the executor should still accept them for flexibility
	config := map[string]interface{}{
		"nats": map[string]interface{}{
			"url":     "nats://invalid.nonexistent.host.test:4222",
			"subject": "events.>",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Expected connection error, got nil")
	}
}

/**
 * Success Path Tests (with mock NATS server)
 */

func TestNATSExecutor_Execute_Success(t *testing.T) {
	port, messages, cleanup := mocknats.StartServer(t)
	defer cleanup()

	executor := NewNATSExecutor()
	config := map[string]interface{}{
		"nats": map[string]interface{}{
			"url":     fmt.Sprintf("nats://127.0.0.1:%d", port),
			"subject": "test.subject",
			"message": "Hello NATS",
		},
	}

	err := executor.Execute(context.Background(), config)
	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	select {
	case msg := <-messages:
		if msg.Subject != "test.subject" {
			t.Errorf("Subject = %q, want 'test.subject'", msg.Subject)
		}
		if string(msg.Data) != "Hello NATS" {
			t.Errorf("Data = %q, want 'Hello NATS'", string(msg.Data))
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for NATS message")
	}
}

func TestNATSExecutor_Execute_Success_WithHeaders(t *testing.T) {
	port, messages, cleanup := mocknats.StartServer(t)
	defer cleanup()

	executor := NewNATSExecutor()
	config := map[string]interface{}{
		"nats": map[string]interface{}{
			"url":     fmt.Sprintf("nats://127.0.0.1:%d", port),
			"subject": "test.headers",
			"message": "With headers",
			"headers": map[string]interface{}{
				"X-Custom":     "custom-value",
				"X-Request-ID": "req-123",
			},
		},
	}

	err := executor.Execute(context.Background(), config)
	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	select {
	case msg := <-messages:
		if msg.Subject != "test.headers" {
			t.Errorf("Subject = %q, want 'test.headers'", msg.Subject)
		}
		if string(msg.Data) != "With headers" {
			t.Errorf("Data = %q, want 'With headers'", string(msg.Data))
		}
		if msg.Headers["X-Custom"] != "custom-value" {
			t.Errorf("Header X-Custom = %q, want 'custom-value'", msg.Headers["X-Custom"])
		}
		if msg.Headers["X-Request-ID"] != "req-123" {
			t.Errorf("Header X-Request-ID = %q, want 'req-123'", msg.Headers["X-Request-ID"])
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for NATS message")
	}
}

func TestNATSExecutor_Execute_Success_MapMessage(t *testing.T) {
	port, messages, cleanup := mocknats.StartServer(t)
	defer cleanup()

	executor := NewNATSExecutor()
	config := map[string]interface{}{
		"nats": map[string]interface{}{
			"url":     fmt.Sprintf("nats://127.0.0.1:%d", port),
			"subject": "test.map",
			"message": map[string]interface{}{
				"event":  "user.created",
				"userId": "user-123",
			},
		},
	}

	err := executor.Execute(context.Background(), config)
	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	select {
	case msg := <-messages:
		if msg.Subject != "test.map" {
			t.Errorf("Subject = %q, want 'test.map'", msg.Subject)
		}
		var payload map[string]interface{}
		if err := json.Unmarshal(msg.Data, &payload); err != nil {
			t.Fatalf("Failed to unmarshal message: %v", err)
		}
		if payload["event"] != "user.created" {
			t.Errorf("event = %v, want 'user.created'", payload["event"])
		}
		if payload["userId"] != "user-123" {
			t.Errorf("userId = %v, want 'user-123'", payload["userId"])
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for NATS message")
	}
}

func TestNATSExecutor_Execute_Success_WithCredentials(t *testing.T) {
	port, messages, cleanup := mocknats.StartServer(t)
	defer cleanup()

	executor := NewNATSExecutor()
	config := map[string]interface{}{
		"nats": map[string]interface{}{
			"url":      fmt.Sprintf("nats://127.0.0.1:%d", port),
			"subject":  "test.auth",
			"message":  "Authenticated message",
			"username": "testuser",
			"password": "testpass",
		},
	}

	err := executor.Execute(context.Background(), config)
	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	select {
	case msg := <-messages:
		if msg.Subject != "test.auth" {
			t.Errorf("Subject = %q, want 'test.auth'", msg.Subject)
		}
		if string(msg.Data) != "Authenticated message" {
			t.Errorf("Data = %q, want 'Authenticated message'", string(msg.Data))
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for NATS message")
	}
}
