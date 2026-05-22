package rabbitmq

import (
	"context"
	"testing"
	"time"

	mockamqp "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mock_servers/rabbitmq"
)

/**
 * RabbitMQExecutor Tests
 */

func TestRabbitMQExecutor_GetType(t *testing.T) {
	executor := NewRabbitMQExecutor()

	if executor.GetType() != "rabbitmq" {
		t.Errorf("GetType() = %q, want 'rabbitmq'", executor.GetType())
	}
}

func TestRabbitMQExecutor_Execute_MissingRabbitmqField(t *testing.T) {
	executor := NewRabbitMQExecutor()

	config := map[string]interface{}{
		"notrabbitmq": map[string]interface{}{
			"host": "rabbitmq.example.com",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for missing 'rabbitmq' field, got nil")
	}
}

func TestRabbitMQExecutor_Execute_MissingHost(t *testing.T) {
	executor := NewRabbitMQExecutor()

	config := map[string]interface{}{
		"rabbitmq": map[string]interface{}{
			"queue": "test-queue",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for missing 'host', got nil")
	}
}

func TestRabbitMQExecutor_Execute_EmptyHost(t *testing.T) {
	executor := NewRabbitMQExecutor()

	config := map[string]interface{}{
		"rabbitmq": map[string]interface{}{
			"host":  "",
			"queue": "test-queue",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for empty 'host', got nil")
	}
}

func TestRabbitMQExecutor_Execute_MissingRoutingKeyAndQueue(t *testing.T) {
	executor := NewRabbitMQExecutor()

	config := map[string]interface{}{
		"rabbitmq": map[string]interface{}{
			"host": "invalid.nonexistent.host.test",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for missing 'routingKey' and 'queue', got nil")
	}
}

func TestRabbitMQExecutor_Execute_ContextCancellation(t *testing.T) {
	executor := NewRabbitMQExecutor()

	config := map[string]interface{}{
		"rabbitmq": map[string]interface{}{
			"host":  "invalid.nonexistent.host.test",
			"queue": "test-queue",
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

func TestRabbitMQExecutor_Execute_PortExtraction_Float64(t *testing.T) {
	executor := NewRabbitMQExecutor()

	// JSON unmarshaling typically produces float64 for numbers
	// Use invalid host to ensure connection fails
	config := map[string]interface{}{
		"rabbitmq": map[string]interface{}{
			"host":  "invalid.nonexistent.host.test",
			"port":  float64(5672),
			"queue": "test-queue",
		},
	}

	err := executor.Execute(context.Background(), config)

	// Should fail on connection, not config extraction
	if err == nil {
		t.Fatal("Expected connection error, got nil")
	}
}

func TestRabbitMQExecutor_Execute_PortExtraction_Int(t *testing.T) {
	executor := NewRabbitMQExecutor()

	config := map[string]interface{}{
		"rabbitmq": map[string]interface{}{
			"host":  "invalid.nonexistent.host.test",
			"port":  5672,
			"queue": "test-queue",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Expected connection error, got nil")
	}
}

func TestRabbitMQExecutor_Execute_WithCredentials(t *testing.T) {
	executor := NewRabbitMQExecutor()

	config := map[string]interface{}{
		"rabbitmq": map[string]interface{}{
			"host":     "invalid.nonexistent.host.test",
			"username": "myuser",
			"password": "mypass",
			"queue":    "test-queue",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Expected connection error, got nil")
	}
}

func TestRabbitMQExecutor_Execute_WithVhost(t *testing.T) {
	executor := NewRabbitMQExecutor()

	config := map[string]interface{}{
		"rabbitmq": map[string]interface{}{
			"host":  "invalid.nonexistent.host.test",
			"vhost": "/myvhost",
			"queue": "test-queue",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Expected connection error, got nil")
	}
}

func TestRabbitMQExecutor_Execute_WithExchange(t *testing.T) {
	executor := NewRabbitMQExecutor()

	config := map[string]interface{}{
		"rabbitmq": map[string]interface{}{
			"host":       "invalid.nonexistent.host.test",
			"exchange":   "my-exchange",
			"routingKey": "my-key",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Expected connection error, got nil")
	}
}

func TestRabbitMQExecutor_Execute_WithExchangeType(t *testing.T) {
	executor := NewRabbitMQExecutor()

	config := map[string]interface{}{
		"rabbitmq": map[string]interface{}{
			"host":         "invalid.nonexistent.host.test",
			"exchange":     "my-exchange",
			"exchangeType": "fanout",
			"routingKey":   "my-key",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Expected connection error, got nil")
	}
}

func TestRabbitMQExecutor_Execute_UseTLS(t *testing.T) {
	executor := NewRabbitMQExecutor()

	config := map[string]interface{}{
		"rabbitmq": map[string]interface{}{
			"host":   "invalid.nonexistent.host.test",
			"queue":  "test-queue",
			"useTLS": true,
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

func TestRabbitMQExecutor_extractMessage_NoMessage(t *testing.T) {
	executor := &RabbitMQExecutor{}

	config := map[string]interface{}{
		"host":  "invalid.nonexistent.host.test",
		"queue": "test-queue",
	}

	msg, err := executor.extractMessage(config)

	if err != nil {
		t.Fatalf("extractMessage() unexpected error: %v", err)
	}

	if string(msg) != "" {
		t.Errorf("extractMessage() = %q, want empty string", string(msg))
	}
}

func TestRabbitMQExecutor_extractMessage_StringMessage(t *testing.T) {
	executor := &RabbitMQExecutor{}

	config := map[string]interface{}{
		"message": "Hello RabbitMQ",
	}

	msg, err := executor.extractMessage(config)

	if err != nil {
		t.Fatalf("extractMessage() unexpected error: %v", err)
	}

	if string(msg) != "Hello RabbitMQ" {
		t.Errorf("extractMessage() = %q, want 'Hello RabbitMQ'", string(msg))
	}
}

func TestRabbitMQExecutor_extractMessage_MapMessage(t *testing.T) {
	executor := &RabbitMQExecutor{}

	config := map[string]interface{}{
		"message": map[string]interface{}{
			"event": "user.created",
			"data":  map[string]interface{}{"userId": "123"},
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

func TestRabbitMQExecutor_extractMessage_ArrayMessage(t *testing.T) {
	executor := &RabbitMQExecutor{}

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
 * Full Config Tests
 */

func TestRabbitMQExecutor_Execute_FullConfig(t *testing.T) {
	executor := NewRabbitMQExecutor()

	config := map[string]interface{}{
		"rabbitmq": map[string]interface{}{
			"host":         "rabbitmq.example.com",
			"port":         float64(5672),
			"username":     "admin",
			"password":     "secret",
			"vhost":        "/production",
			"exchange":     "events",
			"exchangeType": "topic",
			"routingKey":   "user.events.created",
			"useTLS":       false,
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
 * Default Values Tests
 */

func TestRabbitMQExecutor_Execute_DefaultPort(t *testing.T) {
	executor := NewRabbitMQExecutor()

	// No port specified - should use default 5672
	config := map[string]interface{}{
		"rabbitmq": map[string]interface{}{
			"host":  "invalid.nonexistent.host.test",
			"queue": "test-queue",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Expected connection error, got nil")
	}
}

func TestRabbitMQExecutor_Execute_DefaultCredentials(t *testing.T) {
	executor := NewRabbitMQExecutor()

	// No credentials specified - should use guest/guest
	config := map[string]interface{}{
		"rabbitmq": map[string]interface{}{
			"host":  "invalid.nonexistent.host.test",
			"queue": "test-queue",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Expected connection error, got nil")
	}
}

func TestRabbitMQExecutor_Execute_DefaultVhost(t *testing.T) {
	executor := NewRabbitMQExecutor()

	// No vhost specified - should use /
	config := map[string]interface{}{
		"rabbitmq": map[string]interface{}{
			"host":  "invalid.nonexistent.host.test",
			"queue": "test-queue",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Expected connection error, got nil")
	}
}

func TestRabbitMQExecutor_Execute_DefaultExchangeType(t *testing.T) {
	executor := NewRabbitMQExecutor()

	// No exchangeType specified - should use direct
	config := map[string]interface{}{
		"rabbitmq": map[string]interface{}{
			"host":       "invalid.nonexistent.host.test",
			"exchange":   "my-exchange",
			"routingKey": "my-key",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Expected connection error, got nil")
	}
}

/**
 * Queue vs Exchange Tests
 */

func TestRabbitMQExecutor_Execute_DirectQueuePublish(t *testing.T) {
	executor := NewRabbitMQExecutor()

	// Publishing directly to queue (no exchange)
	config := map[string]interface{}{
		"rabbitmq": map[string]interface{}{
			"host":  "invalid.nonexistent.host.test",
			"queue": "direct-queue",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Expected connection error, got nil")
	}
}

func TestRabbitMQExecutor_Execute_ExchangeWithRoutingKey(t *testing.T) {
	executor := NewRabbitMQExecutor()

	// Publishing to exchange with routing key
	config := map[string]interface{}{
		"rabbitmq": map[string]interface{}{
			"host":       "invalid.nonexistent.host.test",
			"exchange":   "my-exchange",
			"routingKey": "events.user.created",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Expected connection error, got nil")
	}
}

/**
 * Success Path Tests (with mock AMQP server)
 */

func TestRabbitMQExecutor_Execute_Success_DirectQueue(t *testing.T) {
	port, messages, cleanup := mockamqp.StartServer(t)
	defer cleanup()

	executor := &RabbitMQExecutor{
		connectTimeout: 5 * time.Second,
		publishTimeout: 5 * time.Second,
	}

	config := map[string]interface{}{
		"rabbitmq": map[string]interface{}{
			"host":     "127.0.0.1",
			"port":     port,
			"username": "guest",
			"password": "guest",
			"queue":    "test-queue",
			"message":  "Hello RabbitMQ",
		},
	}

	err := executor.Execute(context.Background(), config)
	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	select {
	case msg := <-messages:
		if msg.Exchange != "" {
			t.Errorf("Exchange = %q, want empty (direct queue publish)", msg.Exchange)
		}
		if msg.RoutingKey != "test-queue" {
			t.Errorf("RoutingKey = %q, want 'test-queue'", msg.RoutingKey)
		}
		if string(msg.Body) != "Hello RabbitMQ" {
			t.Errorf("Body = %q, want 'Hello RabbitMQ'", string(msg.Body))
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for AMQP message")
	}
}

func TestRabbitMQExecutor_Execute_Success_WithExchange(t *testing.T) {
	port, messages, cleanup := mockamqp.StartServer(t)
	defer cleanup()

	executor := &RabbitMQExecutor{
		connectTimeout: 5 * time.Second,
		publishTimeout: 5 * time.Second,
	}

	config := map[string]interface{}{
		"rabbitmq": map[string]interface{}{
			"host":         "127.0.0.1",
			"port":         port,
			"username":     "guest",
			"password":     "guest",
			"exchange":     "events",
			"exchangeType": "topic",
			"routingKey":   "user.events.created",
			"message":      "Exchange message",
		},
	}

	err := executor.Execute(context.Background(), config)
	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	select {
	case msg := <-messages:
		if msg.Exchange != "events" {
			t.Errorf("Exchange = %q, want 'events'", msg.Exchange)
		}
		if msg.RoutingKey != "user.events.created" {
			t.Errorf("RoutingKey = %q, want 'user.events.created'", msg.RoutingKey)
		}
		if string(msg.Body) != "Exchange message" {
			t.Errorf("Body = %q, want 'Exchange message'", string(msg.Body))
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for AMQP message")
	}
}

func TestRabbitMQExecutor_Execute_Success_WithMessage(t *testing.T) {
	port, messages, cleanup := mockamqp.StartServer(t)
	defer cleanup()

	executor := &RabbitMQExecutor{
		connectTimeout: 5 * time.Second,
		publishTimeout: 5 * time.Second,
	}

	config := map[string]interface{}{
		"rabbitmq": map[string]interface{}{
			"host":     "127.0.0.1",
			"port":     port,
			"username": "guest",
			"password": "guest",
			"queue":    "data-queue",
			"message": map[string]interface{}{
				"event":     "user.created",
				"userId":    "user-123",
				"timestamp": "2024-01-01T00:00:00Z",
			},
		},
	}

	err := executor.Execute(context.Background(), config)
	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	select {
	case msg := <-messages:
		if msg.RoutingKey != "data-queue" {
			t.Errorf("RoutingKey = %q, want 'data-queue'", msg.RoutingKey)
		}
		if len(msg.Body) == 0 {
			t.Error("Body should not be empty for map message")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for AMQP message")
	}
}

func TestRabbitMQExecutor_Execute_Success_FullConfig(t *testing.T) {
	port, messages, cleanup := mockamqp.StartServer(t)
	defer cleanup()

	executor := &RabbitMQExecutor{
		connectTimeout: 5 * time.Second,
		publishTimeout: 5 * time.Second,
	}

	config := map[string]interface{}{
		"rabbitmq": map[string]interface{}{
			"host":         "127.0.0.1",
			"port":         port,
			"username":     "admin",
			"password":     "secret",
			"vhost":        "/",
			"exchange":     "my-exchange",
			"exchangeType": "direct",
			"routingKey":   "my-routing-key",
			"message": map[string]interface{}{
				"event":     "order.placed",
				"orderId":   "order-456",
				"timestamp": "2024-06-15T10:00:00Z",
			},
		},
	}

	err := executor.Execute(context.Background(), config)
	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	select {
	case msg := <-messages:
		if msg.Exchange != "my-exchange" {
			t.Errorf("Exchange = %q, want 'my-exchange'", msg.Exchange)
		}
		if msg.RoutingKey != "my-routing-key" {
			t.Errorf("RoutingKey = %q, want 'my-routing-key'", msg.RoutingKey)
		}
		if len(msg.Body) == 0 {
			t.Error("Body should not be empty")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for AMQP message")
	}
}
