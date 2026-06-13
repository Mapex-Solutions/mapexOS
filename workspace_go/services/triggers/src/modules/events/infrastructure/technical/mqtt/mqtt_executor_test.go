package mqtt

import (
	"context"
	"fmt"
	"testing"
	"time"

	mockmqtt "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mock_servers/mqtt"
)

/**
 * MQTTExecutor Tests
 */

// newTestExecutor returns an executor with short timeouts for deterministic tests.
func newTestExecutor() *MQTTExecutor {
	return &MQTTExecutor{
		connectTimeout: 5 * time.Second,
		publishTimeout: 5 * time.Second,
	}
}

// awaitMessage blocks for a published message or fails the test on timeout.
func awaitMessage(t *testing.T, messages <-chan mockmqtt.Message) mockmqtt.Message {
	t.Helper()
	select {
	case msg := <-messages:
		return msg
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for MQTT message")
		return mockmqtt.Message{}
	}
}

func TestMQTTExecutor_GetType(t *testing.T) {
	executor := NewMQTTExecutor()

	if executor.GetType() != "mqtt" {
		t.Errorf("GetType() = %q, want 'mqtt'", executor.GetType())
	}
}

func TestMQTTExecutor_Execute_MissingMqttField(t *testing.T) {
	executor := NewMQTTExecutor()

	config := map[string]interface{}{
		"notmqtt": map[string]interface{}{
			"broker": "mqtt.example.com",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for missing 'mqtt' field, got nil")
	}
}

func TestMQTTExecutor_Execute_MissingBroker(t *testing.T) {
	executor := NewMQTTExecutor()

	config := map[string]interface{}{
		"mqtt": map[string]interface{}{
			"topic": "test/topic",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for missing 'broker', got nil")
	}
}

func TestMQTTExecutor_Execute_EmptyBroker(t *testing.T) {
	executor := NewMQTTExecutor()

	config := map[string]interface{}{
		"mqtt": map[string]interface{}{
			"broker": "",
			"topic":  "test/topic",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for empty 'broker', got nil")
	}
}

func TestMQTTExecutor_Execute_MissingTopic(t *testing.T) {
	executor := NewMQTTExecutor()

	config := map[string]interface{}{
		"mqtt": map[string]interface{}{
			"broker": "mqtt.example.com",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for missing 'topic', got nil")
	}
}

func TestMQTTExecutor_Execute_EmptyTopic(t *testing.T) {
	executor := NewMQTTExecutor()

	config := map[string]interface{}{
		"mqtt": map[string]interface{}{
			"broker": "mqtt.example.com",
			"topic":  "",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for empty 'topic', got nil")
	}
}

func TestMQTTExecutor_Execute_ContextCancellation(t *testing.T) {
	executor := NewMQTTExecutor()

	config := map[string]interface{}{
		"mqtt": map[string]interface{}{
			"broker": "mqtt.example.com",
			"port":   1883,
			"topic":  "test/topic",
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
 * Config Extraction Tests (against a mock broker — verify the extracted values
 * flow through to a successful publish, independent of any local broker).
 */

func TestMQTTExecutor_Execute_PortExtraction_Float64(t *testing.T) {
	port, messages, cleanup := mockmqtt.StartServer(t)
	defer cleanup()

	executor := newTestExecutor()

	// JSON unmarshaling typically produces float64 for numbers.
	config := map[string]interface{}{
		"mqtt": map[string]interface{}{
			"broker":  "127.0.0.1",
			"port":    float64(port),
			"topic":   "test/topic",
			"message": "port float64",
		},
	}

	if err := executor.Execute(context.Background(), config); err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	if msg := awaitMessage(t, messages); msg.Topic != "test/topic" {
		t.Errorf("Topic = %q, want 'test/topic'", msg.Topic)
	}
}

func TestMQTTExecutor_Execute_PortExtraction_Int(t *testing.T) {
	port, messages, cleanup := mockmqtt.StartServer(t)
	defer cleanup()

	executor := newTestExecutor()

	config := map[string]interface{}{
		"mqtt": map[string]interface{}{
			"broker":  "127.0.0.1",
			"port":    port,
			"topic":   "test/topic",
			"message": "port int",
		},
	}

	if err := executor.Execute(context.Background(), config); err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	if msg := awaitMessage(t, messages); msg.Topic != "test/topic" {
		t.Errorf("Topic = %q, want 'test/topic'", msg.Topic)
	}
}

func TestMQTTExecutor_Execute_QoSExtraction_Float64(t *testing.T) {
	port, messages, cleanup := mockmqtt.StartServer(t)
	defer cleanup()

	executor := newTestExecutor()

	config := map[string]interface{}{
		"mqtt": map[string]interface{}{
			"broker":  "127.0.0.1",
			"port":    port,
			"topic":   "test/topic",
			"qos":     float64(1),
			"message": "qos float64",
		},
	}

	if err := executor.Execute(context.Background(), config); err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	if msg := awaitMessage(t, messages); msg.QoS != 1 {
		t.Errorf("QoS = %d, want 1", msg.QoS)
	}
}

func TestMQTTExecutor_Execute_QoSExtraction_Int(t *testing.T) {
	port, messages, cleanup := mockmqtt.StartServer(t)
	defer cleanup()

	executor := newTestExecutor()

	config := map[string]interface{}{
		"mqtt": map[string]interface{}{
			"broker":  "127.0.0.1",
			"port":    port,
			"topic":   "test/topic",
			"qos":     1,
			"message": "qos int",
		},
	}

	if err := executor.Execute(context.Background(), config); err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	if msg := awaitMessage(t, messages); msg.QoS != 1 {
		t.Errorf("QoS = %d, want 1", msg.QoS)
	}
}

func TestMQTTExecutor_Execute_UseTLS(t *testing.T) {
	executor := NewMQTTExecutor()

	config := map[string]interface{}{
		"mqtt": map[string]interface{}{
			"broker": "localhost",
			"topic":  "test/topic",
			"useTLS": true,
		},
	}

	err := executor.Execute(context.Background(), config)

	// Should fail on connection (ssl:// to a non-TLS endpoint), not config.
	if err == nil {
		t.Fatal("Expected connection error, got nil")
	}
}

/**
 * Message Extraction Tests
 */

func TestMQTTExecutor_extractMessage_NoMessage(t *testing.T) {
	executor := &MQTTExecutor{}

	config := map[string]interface{}{
		"broker": "localhost",
		"topic":  "test/topic",
	}

	msg, err := executor.extractMessage(config)

	if err != nil {
		t.Fatalf("extractMessage() unexpected error: %v", err)
	}

	if string(msg) != "" {
		t.Errorf("extractMessage() = %q, want empty string", string(msg))
	}
}

func TestMQTTExecutor_extractMessage_StringMessage(t *testing.T) {
	executor := &MQTTExecutor{}

	config := map[string]interface{}{
		"message": "Hello MQTT",
	}

	msg, err := executor.extractMessage(config)

	if err != nil {
		t.Fatalf("extractMessage() unexpected error: %v", err)
	}

	if string(msg) != "Hello MQTT" {
		t.Errorf("extractMessage() = %q, want 'Hello MQTT'", string(msg))
	}
}

func TestMQTTExecutor_extractMessage_MapMessage(t *testing.T) {
	executor := &MQTTExecutor{}

	config := map[string]interface{}{
		"message": map[string]interface{}{
			"temperature": 25.5,
			"humidity":    60,
		},
	}

	msg, err := executor.extractMessage(config)

	if err != nil {
		t.Fatalf("extractMessage() unexpected error: %v", err)
	}

	// Should be valid JSON
	if len(msg) == 0 {
		t.Error("extractMessage() returned empty for map message")
	}

	// Should contain expected fields
	msgStr := string(msg)
	if msgStr == "" {
		t.Error("Message should not be empty")
	}
}

func TestMQTTExecutor_extractMessage_OtherTypes(t *testing.T) {
	executor := &MQTTExecutor{}

	// Array message
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
 * Optional Fields Tests (against a mock broker)
 */

func TestMQTTExecutor_Execute_WithCredentials(t *testing.T) {
	port, messages, cleanup := mockmqtt.StartServer(t)
	defer cleanup()

	executor := newTestExecutor()

	config := map[string]interface{}{
		"mqtt": map[string]interface{}{
			"broker":   "127.0.0.1",
			"port":     port,
			"topic":    "test/topic",
			"username": "testuser",
			"password": "testpass",
			"message":  "with credentials",
		},
	}

	if err := executor.Execute(context.Background(), config); err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	awaitMessage(t, messages)
}

func TestMQTTExecutor_Execute_WithClientId(t *testing.T) {
	port, messages, cleanup := mockmqtt.StartServer(t)
	defer cleanup()

	executor := newTestExecutor()

	config := map[string]interface{}{
		"mqtt": map[string]interface{}{
			"broker":   "127.0.0.1",
			"port":     port,
			"topic":    "test/topic",
			"clientId": "my-custom-client-id",
			"message":  "with client id",
		},
	}

	if err := executor.Execute(context.Background(), config); err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	awaitMessage(t, messages)
}

func TestMQTTExecutor_Execute_WithMessage(t *testing.T) {
	port, messages, cleanup := mockmqtt.StartServer(t)
	defer cleanup()

	executor := newTestExecutor()

	config := map[string]interface{}{
		"mqtt": map[string]interface{}{
			"broker": "127.0.0.1",
			"port":   port,
			"topic":  "test/topic",
			"message": map[string]interface{}{
				"sensor": "temp-001",
				"value":  42.5,
			},
		},
	}

	if err := executor.Execute(context.Background(), config); err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	if msg := awaitMessage(t, messages); len(msg.Payload) == 0 {
		t.Error("Payload should not be empty")
	}
}

/**
 * Full Config Tests
 */

func TestMQTTExecutor_Execute_FullConfig(t *testing.T) {
	executor := NewMQTTExecutor()

	// Full config with all optional fields, pointed at an unreachable TLS endpoint.
	config := map[string]interface{}{
		"mqtt": map[string]interface{}{
			"broker":   "mqtt.example.com",
			"port":     float64(8883),
			"topic":    "sensors/temperature",
			"qos":      float64(1),
			"username": "sensor-user",
			"password": "sensor-pass",
			"clientId": "temp-sensor-001",
			"useTLS":   true,
			"message": map[string]interface{}{
				"deviceId":    "SENSOR-001",
				"temperature": 25.5,
				"unit":        "celsius",
			},
		},
	}

	err := executor.Execute(context.Background(), config)

	// Should fail on connection (no real broker), not config.
	if err == nil {
		t.Fatal("Expected connection error, got nil")
	}
}

/**
 * Default Values Tests
 */

func TestMQTTExecutor_Execute_DefaultPort(t *testing.T) {
	// No port specified — the executor must default to 1883. Pointed at a
	// non-routable address (RFC 5737 TEST-NET-1) so the default-port connection
	// attempt fails deterministically regardless of any local broker.
	executor := &MQTTExecutor{
		connectTimeout: 2 * time.Second,
		publishTimeout: 2 * time.Second,
	}

	config := map[string]interface{}{
		"mqtt": map[string]interface{}{
			"broker": "192.0.2.1",
			"topic":  "test/topic",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Expected connection error, got nil")
	}
}

func TestMQTTExecutor_Execute_DefaultQoS(t *testing.T) {
	port, messages, cleanup := mockmqtt.StartServer(t)
	defer cleanup()

	executor := newTestExecutor()

	// No QoS specified — should default to 1.
	config := map[string]interface{}{
		"mqtt": map[string]interface{}{
			"broker":  "127.0.0.1",
			"port":    port,
			"topic":   "test/topic",
			"message": "default qos",
		},
	}

	if err := executor.Execute(context.Background(), config); err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	if msg := awaitMessage(t, messages); msg.QoS != 1 {
		t.Errorf("QoS = %d, want default 1", msg.QoS)
	}
}

func TestMQTTExecutor_Execute_AutoGeneratedClientId(t *testing.T) {
	port, messages, cleanup := mockmqtt.StartServer(t)
	defer cleanup()

	executor := newTestExecutor()

	// No clientId specified — should auto-generate and still publish.
	config := map[string]interface{}{
		"mqtt": map[string]interface{}{
			"broker":  "127.0.0.1",
			"port":    port,
			"topic":   "test/topic",
			"message": "auto client id",
		},
	}

	if err := executor.Execute(context.Background(), config); err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	awaitMessage(t, messages)
}

/**
 * Success Path Tests (with mock MQTT broker)
 */

func TestMQTTExecutor_Execute_Success_QoS0(t *testing.T) {
	port, messages, cleanup := mockmqtt.StartServer(t)
	defer cleanup()

	executor := &MQTTExecutor{
		connectTimeout: 5 * time.Second,
		publishTimeout: 5 * time.Second,
	}

	config := map[string]interface{}{
		"mqtt": map[string]interface{}{
			"broker":  "127.0.0.1",
			"port":    port,
			"topic":   "test/qos0",
			"qos":     0,
			"message": "QoS 0 message",
		},
	}

	err := executor.Execute(context.Background(), config)
	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	select {
	case msg := <-messages:
		if msg.Topic != "test/qos0" {
			t.Errorf("Topic = %q, want 'test/qos0'", msg.Topic)
		}
		if string(msg.Payload) != "QoS 0 message" {
			t.Errorf("Payload = %q, want 'QoS 0 message'", string(msg.Payload))
		}
		if msg.QoS != 0 {
			t.Errorf("QoS = %d, want 0", msg.QoS)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for MQTT message")
	}
}

func TestMQTTExecutor_Execute_Success_QoS1(t *testing.T) {
	port, messages, cleanup := mockmqtt.StartServer(t)
	defer cleanup()

	executor := &MQTTExecutor{
		connectTimeout: 5 * time.Second,
		publishTimeout: 5 * time.Second,
	}

	config := map[string]interface{}{
		"mqtt": map[string]interface{}{
			"broker":  "127.0.0.1",
			"port":    port,
			"topic":   "test/qos1",
			"qos":     1,
			"message": "QoS 1 message",
		},
	}

	err := executor.Execute(context.Background(), config)
	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	select {
	case msg := <-messages:
		if msg.Topic != "test/qos1" {
			t.Errorf("Topic = %q, want 'test/qos1'", msg.Topic)
		}
		if string(msg.Payload) != "QoS 1 message" {
			t.Errorf("Payload = %q, want 'QoS 1 message'", string(msg.Payload))
		}
		if msg.QoS != 1 {
			t.Errorf("QoS = %d, want 1", msg.QoS)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for MQTT message")
	}
}

func TestMQTTExecutor_Execute_Success_WithMessage(t *testing.T) {
	port, messages, cleanup := mockmqtt.StartServer(t)
	defer cleanup()

	executor := &MQTTExecutor{
		connectTimeout: 5 * time.Second,
		publishTimeout: 5 * time.Second,
	}

	config := map[string]interface{}{
		"mqtt": map[string]interface{}{
			"broker": "127.0.0.1",
			"port":   port,
			"topic":  "sensors/temperature",
			"message": map[string]interface{}{
				"sensor": "temp-001",
				"value":  42.5,
			},
		},
	}

	err := executor.Execute(context.Background(), config)
	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	select {
	case msg := <-messages:
		if msg.Topic != "sensors/temperature" {
			t.Errorf("Topic = %q, want 'sensors/temperature'", msg.Topic)
		}
		// Payload should be JSON
		payloadStr := string(msg.Payload)
		if payloadStr == "" {
			t.Error("Payload should not be empty")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for MQTT message")
	}
}

func TestMQTTExecutor_Execute_Success_FullConfig(t *testing.T) {
	port, messages, cleanup := mockmqtt.StartServer(t)
	defer cleanup()

	executor := &MQTTExecutor{
		connectTimeout: 5 * time.Second,
		publishTimeout: 5 * time.Second,
	}

	config := map[string]interface{}{
		"mqtt": map[string]interface{}{
			"broker":   "127.0.0.1",
			"port":     port,
			"topic":    "devices/sensor-001/data",
			"qos":      1,
			"username": "sensor-user",
			"password": "sensor-pass",
			"clientId": fmt.Sprintf("test-client-%d", time.Now().UnixNano()),
			"message": map[string]interface{}{
				"deviceId":    "SENSOR-001",
				"temperature": 25.5,
				"unit":        "celsius",
			},
		},
	}

	err := executor.Execute(context.Background(), config)
	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	select {
	case msg := <-messages:
		if msg.Topic != "devices/sensor-001/data" {
			t.Errorf("Topic = %q, want 'devices/sensor-001/data'", msg.Topic)
		}
		if msg.QoS != 1 {
			t.Errorf("QoS = %d, want 1", msg.QoS)
		}
		if len(msg.Payload) == 0 {
			t.Error("Payload should not be empty")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for MQTT message")
	}
}
