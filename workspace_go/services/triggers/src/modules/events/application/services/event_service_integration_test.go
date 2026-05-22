package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"triggers/src/modules/events/application/di"
	"triggers/src/modules/events/application/ports"
	"triggers/src/modules/events/infrastructure/registry"
	triggerDtos "triggers/src/modules/triggers/application/dtos"

	triggers "github.com/Mapex-Solutions/MapexOS/contracts/services/triggers/triggers"
	mockamqp "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mock_servers/rabbitmq"
	mockmqtt "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mock_servers/mqtt"
	mocknats "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mock_servers/nats"
	mocksmtp "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mock_servers/smtp"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

/**
 * Integration Test Helpers
 *
 * These helpers create an EventService wired with REAL executors (via ExecutorRegistry)
 * and mock protocol servers. Only TriggerService and NatsBus are mocked.
 */

// createIntegrationEventService creates an EventService with a real ExecutorRegistry.
// TriggerService is mocked (returns predefined triggers without DB).
// NatsBus is mocked (no real NATS needed for event publishing).
func createIntegrationEventService(triggerSvc *MockTriggerServicePort, reg ports.ExecutorRegistry) *EventService {
	return &EventService{
		deps: di.EventServiceDependenciesInjection{
			TriggerService:   triggerSvc,
			ExecutorRegistry: reg,
			NatsBus:          &MockCorePublisher{},
			Metrics:          createTestMetrics(),
		},
		workers: 4,
	}
}

// createTriggerResponseForType creates a TriggerResponse for the given trigger type and config.
func createTriggerResponseForType(triggerType, category string, config *triggers.TriggerConfig) *triggerDtos.TriggerResponse {
	name := fmt.Sprintf("Integration Test %s Trigger", triggerType)
	enabled := true
	return &triggerDtos.TriggerResponse{
		Name:        &name,
		TriggerType: &triggerType,
		Category:    &category,
		Enabled:     &enabled,
		Config:      config,
	}
}

/**
 * Integration Tests: EventService → Real Executor → Mock Server
 *
 * These tests exercise the FULL pipeline without mocking executors:
 *   EventService.ProcessTriggerExecutionBatch()
 *   → TriggerConfigToMap()
 *   → ResolvePlaceholdersInMap()
 *   → Real ExecutorRegistry
 *   → Real Executor (NATS/MQTT/RabbitMQ/Email)
 *   → Mock protocol server
 *   → Verify message received
 */

func TestIntegration_EventService_NATS_Success(t *testing.T) {
	// Start mock NATS server
	port, messages, cleanup := mocknats.StartServer(t)
	defer cleanup()

	// Mock TriggerService returns a NATS trigger pointing to mock server.
	// NatsConfig.Server (json:"server") → TriggerConfigToMap → {"nats": {"server": "nats://..."}}
	// The NATS executor accepts both "url" and "server" field names.
	natsURL := fmt.Sprintf("nats://127.0.0.1:%d", port)
	mockTriggerSvc := &MockTriggerServicePort{
		GetTriggerByIdFunc: func(ctx context.Context, triggerId *string) (*triggerDtos.TriggerResponse, error) {
			return createTriggerResponseForType("nats", "technical", &triggers.TriggerConfig{
				Nats: &triggers.NatsConfig{
					Server:  natsURL,
					Subject: "integration.test.subject",
					Message: map[string]interface{}{
						"source": "integration-test",
					},
				},
			}), nil
		},
	}

	// Use real ExecutorRegistry (all real executors)
	realRegistry := registry.NewExecutorRegistry()
	service := createIntegrationEventService(mockTriggerSvc, realRegistry)

	// Build trigger execution event
	eventData := createTriggerExecuteEventJSON("trigger-nats-001", map[string]interface{}{
		"message": "Hello from integration test",
	})
	msg, tracker := createTestMessage(eventData, 0)

	// Execute
	err := service.ProcessTriggerExecutionBatch([]*natsModel.Message{msg})
	if err != nil {
		t.Fatalf("ProcessTriggerExecutionBatch() unexpected error: %v", err)
	}

	// Verify message was ACKed
	if !tracker.AckCalled {
		t.Error("Message should be ACKed on success")
	}
	if tracker.NackCalled || tracker.RejectCalled {
		t.Errorf("Message should only be ACKed, got Nack=%v Reject=%v (reason=%q)",
			tracker.NackCalled, tracker.RejectCalled, tracker.RejectReason)
	}

	// Verify mock NATS server received the message
	select {
	case received := <-messages:
		if received.Subject != "integration.test.subject" {
			t.Errorf("Subject = %q, want 'integration.test.subject'", received.Subject)
		}
		if len(received.Data) == 0 {
			t.Error("Message data should not be empty")
		}
		// Verify message payload contains expected content
		var payload map[string]interface{}
		if err := json.Unmarshal(received.Data, &payload); err != nil {
			t.Fatalf("Failed to unmarshal received message: %v", err)
		}
		if payload["source"] != "integration-test" {
			t.Errorf("payload.source = %v, want 'integration-test'", payload["source"])
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for NATS message on mock server")
	}
}

func TestIntegration_EventService_MQTT_Success(t *testing.T) {
	// Start mock MQTT broker
	port, messages, cleanup := mockmqtt.StartServer(t)
	defer cleanup()

	// Mock TriggerService returns an MQTT trigger
	mockTriggerSvc := &MockTriggerServicePort{
		GetTriggerByIdFunc: func(ctx context.Context, triggerId *string) (*triggerDtos.TriggerResponse, error) {
			return createTriggerResponseForType("mqtt", "technical", &triggers.TriggerConfig{
				Mqtt: &triggers.MqttConfig{
					Broker: "127.0.0.1",
					Port:   port,
					Topic:  "sensors/temperature/data",
					Qos:    1,
					Message: map[string]interface{}{
						"sensorId":    "TEMP-001",
						"temperature": 42.5,
						"unit":        "celsius",
					},
				},
			}), nil
		},
	}

	realRegistry := registry.NewExecutorRegistry()
	service := createIntegrationEventService(mockTriggerSvc, realRegistry)

	eventData := createTriggerExecuteEventJSON("trigger-mqtt-001", map[string]interface{}{
		"deviceId": "sensor-001",
	})
	msg, tracker := createTestMessage(eventData, 0)

	err := service.ProcessTriggerExecutionBatch([]*natsModel.Message{msg})
	if err != nil {
		t.Fatalf("ProcessTriggerExecutionBatch() unexpected error: %v", err)
	}

	if !tracker.AckCalled {
		t.Error("Message should be ACKed on success")
	}
	if tracker.NackCalled || tracker.RejectCalled {
		t.Errorf("Message should only be ACKed, got Nack=%v Reject=%v (reason=%q)",
			tracker.NackCalled, tracker.RejectCalled, tracker.RejectReason)
	}

	// Verify mock MQTT broker received the message
	select {
	case received := <-messages:
		if received.Topic != "sensors/temperature/data" {
			t.Errorf("Topic = %q, want 'sensors/temperature/data'", received.Topic)
		}
		if len(received.Payload) == 0 {
			t.Error("Payload should not be empty")
		}
		// Verify payload contains sensor data
		var payload map[string]interface{}
		if err := json.Unmarshal(received.Payload, &payload); err != nil {
			t.Fatalf("Failed to unmarshal MQTT payload: %v", err)
		}
		if payload["sensorId"] != "TEMP-001" {
			t.Errorf("payload.sensorId = %v, want 'TEMP-001'", payload["sensorId"])
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for MQTT message on mock broker")
	}
}

func TestIntegration_EventService_RabbitMQ_Success(t *testing.T) {
	// Start mock AMQP server
	port, messages, cleanup := mockamqp.StartServer(t)
	defer cleanup()

	// Mock TriggerService returns a RabbitMQ trigger
	exchangeName := "integration-events"
	exchangeType := "topic"
	routingKey := "user.events.created"
	mockTriggerSvc := &MockTriggerServicePort{
		GetTriggerByIdFunc: func(ctx context.Context, triggerId *string) (*triggerDtos.TriggerResponse, error) {
			return createTriggerResponseForType("rabbitmq", "technical", &triggers.TriggerConfig{
				Rabbitmq: &triggers.RabbitmqConfig{
					Host:         "127.0.0.1",
					Port:         port,
					Username:     "guest",
					Password:     "guest",
					PublishMode:  "exchange",
					Exchange:     &exchangeName,
					ExchangeType: &exchangeType,
					RoutingKey:   &routingKey,
					Message: map[string]interface{}{
						"event":  "user.created",
						"userId": "user-integration-123",
					},
				},
			}), nil
		},
	}

	realRegistry := registry.NewExecutorRegistry()
	service := createIntegrationEventService(mockTriggerSvc, realRegistry)

	eventData := createTriggerExecuteEventJSON("trigger-rabbitmq-001", map[string]interface{}{
		"userId": "user-123",
	})
	msg, tracker := createTestMessage(eventData, 0)

	err := service.ProcessTriggerExecutionBatch([]*natsModel.Message{msg})
	if err != nil {
		t.Fatalf("ProcessTriggerExecutionBatch() unexpected error: %v", err)
	}

	if !tracker.AckCalled {
		t.Error("Message should be ACKed on success")
	}
	if tracker.NackCalled || tracker.RejectCalled {
		t.Errorf("Message should only be ACKed, got Nack=%v Reject=%v (reason=%q)",
			tracker.NackCalled, tracker.RejectCalled, tracker.RejectReason)
	}

	// Verify mock AMQP server received the message
	select {
	case received := <-messages:
		if received.Exchange != "integration-events" {
			t.Errorf("Exchange = %q, want 'integration-events'", received.Exchange)
		}
		if received.RoutingKey != "user.events.created" {
			t.Errorf("RoutingKey = %q, want 'user.events.created'", received.RoutingKey)
		}
		if len(received.Body) == 0 {
			t.Error("Body should not be empty")
		}
		var payload map[string]interface{}
		if err := json.Unmarshal(received.Body, &payload); err != nil {
			t.Fatalf("Failed to unmarshal AMQP body: %v", err)
		}
		if payload["event"] != "user.created" {
			t.Errorf("payload.event = %v, want 'user.created'", payload["event"])
		}
		if payload["userId"] != "user-integration-123" {
			t.Errorf("payload.userId = %v, want 'user-integration-123'", payload["userId"])
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for AMQP message on mock server")
	}
}

func TestIntegration_EventService_Email_Success(t *testing.T) {
	// Start mock SMTP server
	smtpPort, messages, cleanup := mocksmtp.StartServer(t)
	defer cleanup()

	// SMTP config is now part of the trigger document (config.email),
	// same pattern as HTTP (endpoint) and MQTT (broker).
	username := "testuser"
	password := "testpass"
	mockTriggerSvc := &MockTriggerServicePort{
		GetTriggerByIdFunc: func(ctx context.Context, triggerId *string) (*triggerDtos.TriggerResponse, error) {
			subject := "Integration Test Alert"
			body := "Sensor TEMP-001 reported critical temperature"
			return createTriggerResponseForType("email", "communication", &triggers.TriggerConfig{
				Email: &triggers.EmailConfig{
					SmtpHost: "127.0.0.1",
					SmtpPort: smtpPort,
					Username: &username,
					Password: &password,
					FromAddr: "alerts@test.com",
					To:       "admin@example.com",
					Subject:  subject,
					Body:     &body,
				},
			}), nil
		},
	}

	service := createIntegrationEventService(mockTriggerSvc, registry.NewExecutorRegistry())

	eventData := createTriggerExecuteEventJSON("trigger-email-001", map[string]interface{}{
		"alertType": "temperature",
	})
	msg, tracker := createTestMessage(eventData, 0)

	err := service.ProcessTriggerExecutionBatch([]*natsModel.Message{msg})
	if err != nil {
		t.Fatalf("ProcessTriggerExecutionBatch() unexpected error: %v", err)
	}

	if !tracker.AckCalled {
		t.Error("Message should be ACKed on success")
	}
	if tracker.NackCalled || tracker.RejectCalled {
		t.Errorf("Message should only be ACKed, got Nack=%v Reject=%v (reason=%q)",
			tracker.NackCalled, tracker.RejectCalled, tracker.RejectReason)
	}

	// Verify mock SMTP server received the email
	select {
	case received := <-messages:
		if received.From != "alerts@test.com" {
			t.Errorf("From = %q, want 'alerts@test.com'", received.From)
		}
		if len(received.Recipients) != 1 || received.Recipients[0] != "admin@example.com" {
			t.Errorf("Recipients = %v, want ['admin@example.com']", received.Recipients)
		}
		if !strings.Contains(received.Data, "Subject: Integration Test Alert") {
			t.Error("Email should contain the subject header")
		}
		if !strings.Contains(received.Data, "Sensor TEMP-001 reported critical temperature") {
			t.Error("Email should contain the body text")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for email on mock SMTP server")
	}
}

func TestIntegration_EventService_PlaceholderResolution_NATS(t *testing.T) {
	// Start mock NATS server
	port, messages, cleanup := mocknats.StartServer(t)
	defer cleanup()

	// Config with placeholders in subject and message
	natsURL := fmt.Sprintf("nats://127.0.0.1:%d", port)
	mockTriggerSvc := &MockTriggerServicePort{
		GetTriggerByIdFunc: func(ctx context.Context, triggerId *string) (*triggerDtos.TriggerResponse, error) {
			return createTriggerResponseForType("nats", "technical", &triggers.TriggerConfig{
				Nats: &triggers.NatsConfig{
					Server:  natsURL,
					Subject: "alerts.{{severity}}",
					Message: map[string]interface{}{
						"text": "Sensor {{sensorId}} at {{temperature}}°C",
					},
				},
			}), nil
		},
	}

	realRegistry := registry.NewExecutorRegistry()
	service := createIntegrationEventService(mockTriggerSvc, realRegistry)

	// Payload with placeholder values
	eventData := createTriggerExecuteEventJSON("trigger-nats-placeholder", map[string]interface{}{
		"severity":    "high",
		"sensorId":    "TEMP-001",
		"temperature": 42.5,
	})
	msg, tracker := createTestMessage(eventData, 0)

	err := service.ProcessTriggerExecutionBatch([]*natsModel.Message{msg})
	if err != nil {
		t.Fatalf("ProcessTriggerExecutionBatch() unexpected error: %v", err)
	}

	if !tracker.AckCalled {
		t.Error("Message should be ACKed on success")
	}

	// Verify placeholders were resolved in the received message
	select {
	case received := <-messages:
		// Subject placeholder: "alerts.{{severity}}" → "alerts.high"
		if received.Subject != "alerts.high" {
			t.Errorf("Subject = %q, want 'alerts.high' (placeholder not resolved)", received.Subject)
		}
		// Message placeholder: "Sensor {{sensorId}} at {{temperature}}°C"
		var payload map[string]interface{}
		if err := json.Unmarshal(received.Data, &payload); err != nil {
			t.Fatalf("Failed to unmarshal received message: %v", err)
		}
		text, ok := payload["text"].(string)
		if !ok {
			t.Fatal("payload.text should be a string")
		}
		expectedText := "Sensor TEMP-001 at 42.5°C"
		if text != expectedText {
			t.Errorf("payload.text = %q, want %q", text, expectedText)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for NATS message on mock server")
	}
}

func TestIntegration_EventService_BatchMixed(t *testing.T) {
	// Start mock NATS server for the one successful message
	port, natsMessages, cleanup := mocknats.StartServer(t)
	defer cleanup()

	natsURL := fmt.Sprintf("nats://127.0.0.1:%d", port)

	mockTriggerSvc := &MockTriggerServicePort{
		GetTriggerByIdFunc: func(ctx context.Context, triggerId *string) (*triggerDtos.TriggerResponse, error) {
			switch *triggerId {
			case "trigger-nats-ok":
				return createTriggerResponseForType("nats", "technical", &triggers.TriggerConfig{
					Nats: &triggers.NatsConfig{
						Server:  natsURL,
						Subject: "batch.test",
						Message: map[string]interface{}{
							"source": "batch-integration",
						},
					},
				}), nil
			case "trigger-disabled":
				resp := createTriggerResponseForType("nats", "technical", nil)
				enabled := false
				resp.Enabled = &enabled
				return resp, nil
			default:
				return nil, fmt.Errorf("trigger not found: %s", *triggerId)
			}
		},
	}

	realRegistry := registry.NewExecutorRegistry()
	service := createIntegrationEventService(mockTriggerSvc, realRegistry)

	// Message 0: NATS success → ACK
	eventData0 := createTriggerExecuteEventJSON("trigger-nats-ok", map[string]interface{}{})
	msg0, tracker0 := createTestMessage(eventData0, 0)

	// Message 1: Trigger disabled → ACK (skip)
	eventData1 := createTriggerExecuteEventJSON("trigger-disabled", map[string]interface{}{})
	msg1, tracker1 := createTestMessage(eventData1, 1)

	// Message 2: Invalid JSON → Reject
	msg2, tracker2 := createTestMessage([]byte("not-valid-json!!!"), 2)

	err := service.ProcessTriggerExecutionBatch([]*natsModel.Message{msg0, msg1, msg2})
	if err != nil {
		t.Fatalf("ProcessTriggerExecutionBatch() unexpected error: %v", err)
	}

	// Message 0: NATS success → ACK
	if !tracker0.AckCalled {
		t.Error("Message 0 (NATS success) should be ACKed")
	}
	if tracker0.NackCalled || tracker0.RejectCalled {
		t.Error("Message 0 should only be ACKed")
	}

	// Message 1: Trigger disabled → ACK (skip silently)
	if !tracker1.AckCalled {
		t.Error("Message 1 (disabled trigger) should be ACKed")
	}
	if tracker1.NackCalled || tracker1.RejectCalled {
		t.Error("Message 1 should only be ACKed")
	}

	// Message 2: Invalid JSON → Reject
	if !tracker2.RejectCalled {
		t.Error("Message 2 (invalid JSON) should be Rejected")
	}
	if tracker2.AckCalled || tracker2.NackCalled {
		t.Error("Message 2 should only be Rejected")
	}

	// Verify mock NATS server received exactly 1 message (from message 0)
	select {
	case received := <-natsMessages:
		if received.Subject != "batch.test" {
			t.Errorf("Subject = %q, want 'batch.test'", received.Subject)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for NATS message from batch")
	}

	// Verify no more messages were received
	select {
	case extra := <-natsMessages:
		t.Errorf("Mock server should have received only 1 message, got extra: subject=%q", extra.Subject)
	case <-time.After(200 * time.Millisecond):
		// Expected: no more messages
	}
}
