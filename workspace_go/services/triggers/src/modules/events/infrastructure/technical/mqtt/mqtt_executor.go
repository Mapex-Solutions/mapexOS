package mqtt

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"time"

	"triggers/src/modules/events/application/ports"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// NewMQTTExecutor creates a new MQTT trigger executor adapter.
func NewMQTTExecutor() ports.TriggerExecutor {
	return &MQTTExecutor{
		connectTimeout: 10 * time.Second,
		publishTimeout: 5 * time.Second,
	}
}

// Execute publishes a message to an MQTT broker based on the trigger configuration.
//
// This is the concrete implementation of the ports.TriggerExecutor interface.
//
// Steps:
// 1. Extract MQTT config (broker, port, topic, qos, credentials, message)
// 2. Build MQTT client options
// 3. Connect to broker
// 4. Publish message
// 5. Disconnect
//
// Parameters:
//   - ctx: Context for controlling cancellation and timeouts
//   - config: Trigger configuration (mqtt field) with placeholders already resolved
//
// Returns:
//   - error: If MQTT publish fails
func (e *MQTTExecutor) Execute(ctx context.Context, config map[string]interface{}) error {
	// Extract mqtt config
	mqttConfig, ok := config["mqtt"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("MQTT trigger config missing 'mqtt' field")
	}

	// Extract broker (required)
	broker, ok := mqttConfig["broker"].(string)
	if !ok || broker == "" {
		return fmt.Errorf("MQTT trigger config missing required field 'broker'")
	}

	// Extract port (required)
	port := 1883 // Default MQTT port
	if portVal, exists := mqttConfig["port"]; exists {
		switch p := portVal.(type) {
		case float64:
			port = int(p)
		case int:
			port = p
		}
	}

	// Extract topic (required)
	topic, ok := mqttConfig["topic"].(string)
	if !ok || topic == "" {
		return fmt.Errorf("MQTT trigger config missing required field 'topic'")
	}

	// Extract QoS (optional, default 1)
	qos := byte(1)
	if qosVal, exists := mqttConfig["qos"]; exists {
		switch q := qosVal.(type) {
		case float64:
			qos = byte(q)
		case int:
			qos = byte(q)
		}
	}

	// Extract optional fields
	username, _ := mqttConfig["username"].(string)
	password, _ := mqttConfig["password"].(string)
	clientId, _ := mqttConfig["clientId"].(string)
	if clientId == "" {
		clientId = fmt.Sprintf("triggers-executor-%d", time.Now().UnixNano())
	}

	// Check TLS setting
	useTLS := false
	if tlsVal, exists := mqttConfig["useTLS"]; exists {
		if b, ok := tlsVal.(bool); ok {
			useTLS = b
		}
	}

	// Extract message payload
	message, err := e.extractMessage(mqttConfig)
	if err != nil {
		return fmt.Errorf("failed to extract MQTT message: %w", err)
	}

	// Build broker URL
	protocol := "tcp"
	if useTLS {
		protocol = "ssl"
	}
	brokerURL := fmt.Sprintf("%s://%s:%d", protocol, broker, port)

	// Build MQTT client options
	opts := pahomqtt.NewClientOptions()
	opts.AddBroker(brokerURL)
	opts.SetClientID(clientId)
	opts.SetConnectTimeout(e.connectTimeout)

	if username != "" {
		opts.SetUsername(username)
	}
	if password != "" {
		opts.SetPassword(password)
	}

	if useTLS {
		opts.SetTLSConfig(&tls.Config{
			InsecureSkipVerify: false,
		})
	}

	// Create client
	client := pahomqtt.NewClient(opts)

	logger.Info(fmt.Sprintf("[INFRA:MQTTExecutor] Connecting to %s", brokerURL))

	// Connect with context cancellation check
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	token := client.Connect()
	if !token.WaitTimeout(e.connectTimeout) {
		return fmt.Errorf("MQTT connection timeout")
	}
	if token.Error() != nil {
		logger.Error(token.Error(), fmt.Sprintf("[INFRA:MQTTExecutor] Connection failed: %s", brokerURL))
		return fmt.Errorf("MQTT connection failed: %w", token.Error())
	}

	defer func() {
		client.Disconnect(250)
		logger.Debug(fmt.Sprintf("[INFRA:MQTTExecutor] Disconnected from %s", brokerURL))
	}()

	logger.Info(fmt.Sprintf("[INFRA:MQTTExecutor] Publishing to topic: %s", topic))

	// Publish message
	pubToken := client.Publish(topic, qos, false, message)
	if !pubToken.WaitTimeout(e.publishTimeout) {
		return fmt.Errorf("MQTT publish timeout")
	}
	if pubToken.Error() != nil {
		logger.Error(pubToken.Error(), fmt.Sprintf("[INFRA:MQTTExecutor] Publish failed: %s", topic))
		return fmt.Errorf("MQTT publish failed: %w", pubToken.Error())
	}

	logger.Info(fmt.Sprintf("[INFRA:MQTTExecutor] Published successfully to topic: %s", topic))
	return nil
}

// GetType returns the trigger type this executor handles.
func (e *MQTTExecutor) GetType() string {
	return "mqtt"
}

// extractMessage extracts and serializes the message payload from config.
func (e *MQTTExecutor) extractMessage(config map[string]interface{}) ([]byte, error) {
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
